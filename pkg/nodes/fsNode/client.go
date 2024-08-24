package fsNode

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"github.com/tomp332/p2p-agent/pkg/utils/types"
	"io"
	"os"
	"time"
)

type FileNodeClientType interface {
	DownloadFile(ctx context.Context, fileId string) (<-chan []byte, <-chan error)
	Authenticate(ctx context.Context, username string, password string) (*pb.AuthenticateResponse, error)
	UploadFile(ctx context.Context, filePath string) <-chan error
}

type FileNodeClient struct {
	client         pb.FilesNodeServiceClient
	requestTimeout time.Duration
}

func NewFileNodeClient(client pb.FilesNodeServiceClient, requestTimeout time.Duration) *FileNodeClient {
	return &FileNodeClient{
		client:         client,
		requestTimeout: requestTimeout,
	}
}

func (fc *FileNodeClient) DownloadFile(ctx context.Context, fileId string) (<-chan types.TransferChunkData, <-chan error) {
	dataChan := make(chan types.TransferChunkData)
	errChan := make(chan error, 1) // buffered channel to avoid goroutine leak

	// Create a goroutine to handle the download.
	go func() {
		defer close(dataChan)
		defer close(errChan) // Ensure errChan is closed when goroutine finishes

		req := &pb.DownloadFileRequest{FileId: fileId}
		stream, err := fc.client.DownloadFile(ctx, req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to download file from remote peer.")
			errChan <- err // Send the error to the error channel
			return
		}

		for {
			select {
			case <-ctx.Done():
				// Handle context cancellation
				log.Info().Str("fileId", fileId).Msg("Download cancelled by context")
				errChan <- ctx.Err() // Send the context error to the error channel
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
				log.Debug().Str("fileId", fileId).Int64("size", res.GetChunkSize()).Msg("Received file chunk from remote peer")
				dataChan <- types.TransferChunkData{
					ChunkSize: res.GetChunkSize(),
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
		stream, streamErr := fc.client.UploadFile(ctx)
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

				if !response.GetSuccess() {
					errChan <- fmt.Errorf("file upload failed: %s", response.GetMessage())
					return
				}

				log.Debug().
					Str("fileId", response.GetFileId()).
					Int64("fileSize", response.GetFileSize()).
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
