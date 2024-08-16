package file_node

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/nodes"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"io"
	"os"
	"time"
)

type FileNodeClient struct {
	nodes.P2PNodeClienter
	client         pb.FilesNodeServiceClient
	requestTimeout time.Duration
}

func NewFileNodeClient(client pb.FilesNodeServiceClient, requestTimeout time.Duration) *FileNodeClient {
	return &FileNodeClient{
		client:         client,
		requestTimeout: requestTimeout,
	}
}

func (fc *FileNodeClient) SearchFile(ctx context.Context, fileId string) (bool, error) {
	req := &pb.SearchFileRequest{FileId: fileId}
	res, err := fc.client.SearchFile(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("could not search file on current node.")
		return false, err
	}
	return res.Exists, nil
}

func (fc *FileNodeClient) DirectDownloadFile(ctx context.Context, fileId string, nodeId string) (<-chan []byte, <-chan error) {
	dataChan := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errChan)

		req := &pb.DirectDownloadFileRequest{NodeId: nodeId, FileId: fileId}
		stream, err := fc.client.DirectDownloadFile(ctx, req)
		if err != nil {
			errChan <- err
			return
		}

		for {
			res, streamErr := stream.Recv()
			if streamErr == io.EOF {
				// End of the stream
				errChan <- nil
				return
			}
			if streamErr != nil {
				errChan <- streamErr
				return
			}
			// Send the received chunk to the data channel
			dataChan <- res.GetChunk()
		}
	}()

	return dataChan, errChan
}

func (fc *FileNodeClient) DownloadFile(ctx context.Context, fileId string) (<-chan []byte, <-chan error) {
	dataChan := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errChan)

		req := &pb.DownloadFileRequest{FileId: fileId}
		stream, err := fc.client.DownloadFile(ctx, req)
		if err != nil {
			errChan <- err
			return
		}

		for {
			res, streamErr := stream.Recv()
			if streamErr == io.EOF {
				// End of the stream
				errChan <- nil
				return
			}
			if streamErr != nil {
				errChan <- streamErr
				return
			}
			// Send the received chunk to the data channel
			dataChan <- res.GetChunk()
		}
	}()

	return dataChan, errChan
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
					Float64("fileSize", response.GetFileSize()).
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
