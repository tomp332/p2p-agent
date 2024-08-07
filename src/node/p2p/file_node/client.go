package file_node

import (
	"context"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/utils"
	"io"
	"time"
)

type FileNodeClient struct {
	src.P2PNodeClienter
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
		utils.Logger.Error().Err(err).Msg("could not search file on current node.")
		return false, err
	}
	return res.Exists, nil
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
			res, err := stream.Recv()
			if err == io.EOF {
				// End of the stream
				errChan <- nil
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			// Send the received chunk to the data channel
			dataChan <- res.GetChunk()
		}
	}()

	return dataChan, errChan
}
