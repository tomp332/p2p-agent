package fsNode

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"github.com/tomp332/p2p-agent/pkg/utils/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"slices"
	"time"
)

type FileNodeClientType interface {
	DownloadFile(ctx context.Context, fileId string) (<-chan []byte, <-chan error)
	Authenticate(ctx context.Context, username string, password string) (*pb.AuthenticateResponse, error)
	UploadFile(ctx context.Context, filePath string) <-chan error
}

type FileNodeClient struct {
	client          pb.FilesNodeServiceClient
	authToken       *string
	protectedRoutes *[]string
	requestTimeout  time.Duration
}

func NewFileNodeClient(client pb.FilesNodeServiceClient, authToken *string, requestTimeout time.Duration, protectedRoutes *[]string) *FileNodeClient {
	return &FileNodeClient{
		client:          client,
		authToken:       authToken,
		requestTimeout:  requestTimeout,
		protectedRoutes: protectedRoutes,
	}
}

func (fc *FileNodeClient) SearchFile(searchCtx context.Context, fileName string) (*pb.SearchFileResponse, error) {
	return fc.client.SearchFile(searchCtx, &pb.SearchFileRequest{
		FileName: fileName,
	})
}

func (fc *FileNodeClient) DownloadFile(searchCtx context.Context, fileName string) (<-chan types.TransferChunkData, <-chan error) {
	searchAuthCtx := fc.appendTokenToCall(searchCtx) // Ensure token is attached
	dataChan := make(chan types.TransferChunkData)
	errChan := make(chan error, 1) // Buffered channel to avoid goroutine leak
	// Create a goroutine to handle the download
	go func() {
		defer close(dataChan)
		defer close(errChan) // Ensure errChan is closed when goroutine finishes

		req := &pb.DownloadFileRequest{FileName: fileName}
		stream, err := fc.client.DownloadFile(searchAuthCtx, req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to download file from remote peer.")
			errChan <- err // Send the error to the error channel
			return
		}

		// Main goroutine to receive data from the stream
		for {
			select {
			case <-searchAuthCtx.Done():
				log.Info().Str("fileName", fileName).Msg("Streaming context canceled, stopping download")
				return
			default:
				res, streamErr := stream.Recv()
				if streamErr == io.EOF {
					// End of the stream
					log.Debug().Msg("End of stream for file download from remote peer")
					return
				}
				if streamErr != nil {
					log.Error().Err(streamErr).Msg("Failed to receive file chunk from remote peer")
					errChan <- streamErr // Send the stream error to the error channel
					return
				}
				// Send the received chunk to the data channel
				dataChan <- types.TransferChunkData{
					ChunkSize: int(res.GetChunkSize()),
					ChunkData: res.GetChunk(),
				}
			}
		}
	}()

	return dataChan, errChan // Return both channels to the caller
}

func (fc *FileNodeClient) Authenticate(ctx context.Context, username string, password string) (*pb.AuthenticateResponse, error) {
	req := &pb.AuthenticateRequest{Username: username, Password: password}
	res, err := fc.client.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (fc *FileNodeClient) UploadFile(ctx context.Context, filePath string) <-chan error {
	authCtx := fc.appendTokenToCall(ctx)
	errChan := make(chan error, 1) // Buffered channel to send any errors

	go func() {
		defer close(errChan)

		// Open the file for reading
		file, openErr := os.Open(filePath)
		if openErr != nil {
			log.Error().Msgf("Failed to open file: %v", openErr)
			errChan <- openErr
			return
		}
		defer file.Close()

		// Create a new upload stream
		stream, streamErr := fc.client.UploadFile(authCtx)
		if streamErr != nil {
			log.Error().Msgf("Failed to create upload stream: %v", streamErr)
			errChan <- streamErr
			return
		}

		buffer := make([]byte, 1024) // Buffer to hold file chunks
		for {
			// Read a chunk from the file
			bytesRead, readErr := file.Read(buffer)
			if readErr == io.EOF {
				// End of file reached, close the stream
				response, closeErr := stream.CloseAndRecv()
				if closeErr != nil {
					log.Error().Msgf("Failed to close stream and receive response: %v", closeErr)
					errChan <- closeErr
					return
				}

				if response.GetFileHash() != "" {
					errChan <- fmt.Errorf("file upload failed")
					return
				}

				log.Debug().
					Str("fileName", response.GetFileName()).
					Uint64("fileSize", response.GetFileSize()).
					Str("fileHash", response.GetFileHash()).
					Msgf("File uploaded successfully")
				return
			}

			if readErr != nil {
				log.Error().Msgf("Failed to read file: %v", readErr)
				errChan <- readErr
				return
			}

			// Send the chunk to the server
			sendErr := stream.Send(&pb.UploadFileRequest{
				ChunkData: buffer[:bytesRead],
			})
			if sendErr != nil {
				log.Error().Msgf("Failed to send file chunk: %v", sendErr)
				errChan <- sendErr
				return
			}
		}
	}()

	return errChan
}

func (fc *FileNodeClient) appendTokenToCall(parentCtx context.Context) context.Context {
	method, _ := grpc.Method(parentCtx)
	if fc.protectedRoutes == nil || fc.authToken == nil {
		log.Trace().Msg("No auth token or protected routes provided for current fsNode client, skipping auth attachment.")
		return parentCtx
	}
	if !slices.Contains(*fc.protectedRoutes, method) {
		// Method is not protected and does not need to attach a token.
		return nil
	}
	return metadata.AppendToOutgoingContext(parentCtx, "authorization", "Bearer "+*fc.authToken)
}
