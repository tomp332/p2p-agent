package p2p_file_node

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/storage"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"google.golang.org/grpc"
	"io"
	"sync"
)

type FilesNode struct {
	*p2p.BaseNode
	Storage *storage.LocalStorage
	pb.UnimplementedFilesNodeServiceServer
	config *configs.P2PFilesNodeConfig
}

type FileNodeConnection struct {
	p2p.NodeConnection
	NodeClient FileNodeClient
}

func NewP2PFilesNode(baseNode *p2p.BaseNode, nodeOptions *configs.P2PFilesNodeConfig) *FilesNode {
	n := &FilesNode{BaseNode: baseNode, Storage: storage.NewLocalStorage(&nodeOptions.Storage)}
	return n
}

func (n *FilesNode) UploadFile(stream pb.FilesNodeService_UploadFileServer) error {
	var buffer bytes.Buffer
	dataChan := make(chan []byte, 1024) // Buffered channel for file data

	ctx := stream.Context()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fileID, err := createUniqueFileID(&buffer)
			close(dataChan)

			fileSize, err := n.Storage.Put(ctx, fileID, dataChan)
			if err != nil {
				return stream.SendAndClose(&pb.UploadFileResponse{
					FileId:  fileID,
					Success: false,
					Message: fmt.Sprintf("Failed to store file: %v", err),
				})
			}

			return stream.SendAndClose(&pb.UploadFileResponse{
				FileId:   fileID,
				FileSize: fileSize,
				Success:  true,
				Message:  "File uploaded successfully",
			})
		}
		if err != nil {
			return err
		}
		chunk := req.GetChunkData()
		buffer.Write(chunk)
		dataChan <- chunk
	}
}

func (n *FilesNode) DownloadFile(req *pb.DownloadFileRequest, stream pb.FilesNodeService_DownloadFileServer) error {
	ctx := stream.Context()
	var mainDataChan <-chan []byte
	var errChan <-chan error // Added error channel

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		foundFileNode, err := n.searchFileInNetwork(req.FileId)
		if err != nil {
			// File not found in the network, retrieve from storage
			dataChan, storageErr := n.Storage.Get(ctx, req.FileId)
			if storageErr != nil {
				return storageErr
			}
			mainDataChan = dataChan
		} else {
			// Found the required file in the network, fetch from remote
			mainDataChan, errChan = foundFileNode.NodeClient.DownloadFile(ctx, req.GetFileId())
		}

		for {
			select {
			case chunk, ok := <-mainDataChan:
				if !ok {
					// Channel closed, end of data
					return nil
				}
				if streamErr := stream.Send(&pb.DownloadFileResponse{
					FileId: req.FileId,
					Exists: true,
					Chunk:  chunk,
				}); streamErr != nil {
					return streamErr
				}
			case downloadErr := <-errChan:
				// Handle error from errChan (if any)
				if downloadErr != nil {
					return downloadErr
				}
			case <-ctx.Done():
				// Handle context cancellation
				return ctx.Err()
			}
		}
	}
}

func (n *FilesNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	err := n.Storage.Delete(ctx, request.GetFileId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{FileId: request.GetFileId()}, nil
}

func (n *FilesNode) Register(server *grpc.Server) {
	pb.RegisterFilesNodeServiceServer(server, n)
	utils.Logger.Info().Str("nodeType", n.Type).Msg("Node registered")
}

func (n *FilesNode) ConnectToBootstrapPeers(server src.AgentGRPCServer) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(n.BootstrapPeerAddrs))

	for _, address := range n.BootstrapPeerAddrs {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			con, err := n.ConnectToPeer(server, address, n.NodeConfigs.BootstrapNodeTimeout)
			if err != nil {
				utils.Logger.Error().Err(err).Msg("Failed to connect to bootstrap peer")
				errChan <- err
				return
			}
			client := pb.NewFilesNodeServiceClient(con)
			n.ConnectedPeers = append(n.ConnectedPeers, FileNodeConnection{
				NodeConnection: p2p.NodeConnection{
					Address:        address,
					GrpcConnection: con,
				},
				NodeClient: *NewFileNodeClient(client, n.BootstrapNodeTimeout),
			})
			utils.Logger.Debug().Str("address", address).Msgf("Successfully connected to bootstrap node.")
		}(address)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (n *FilesNode) SearchFile(_ context.Context, request *pb.SearchFileRequest) (*pb.SearchFileResponse, error) {
	response := &pb.SearchFileResponse{
		FileId: request.GetFileId(),
		Exists: true,
	}
	if n.Storage.Search(request.GetFileId()) {
		response.Exists = true
	} else {
		response.Exists = false
	}
	return response, nil
}

func createUniqueFileID(buffer *bytes.Buffer) (string, error) {
	fileHash := md5.Sum(buffer.Bytes())
	return hex.EncodeToString(fileHash[:]) + "-" + utils.GenerateRandomID(), nil
}

func (n *FilesNode) searchFileInNetwork(fileId string) (FileNodeConnection, error) {
	// Check if there are no connected peers
	if len(n.ConnectedPeers) == 0 {
		return FileNodeConnection{}, errors.New("no connected peers")
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileFoundChan := make(chan FileNodeConnection, 1) // Buffer size 1 to avoid blocking

	// Send search RPC call to all connected peers, to see if the network already contains this file.
	for _, peer := range n.ConnectedPeers {
		wg.Add(1)
		go func(peer src.P2PNodeConnection) {
			defer wg.Done()
			dynamicPeer := peer.(FileNodeConnection)
			result, err := dynamicPeer.NodeClient.SearchFile(ctx, fileId)
			if err != nil {
				utils.Logger.Debug().Str("peer", dynamicPeer.Address).Msgf("Failed to send SearchFile RPC call for peer.")
				return
			}
			if result {
				select {
				case fileFoundChan <- dynamicPeer:
					cancel() // Cancel the context to stop other goroutines
				default:
				}
			}
		}(peer)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(fileFoundChan)
	}()

	// Check for results
	select {
	case fileNode := <-fileFoundChan:
		return fileNode, nil
	case <-ctx.Done(): // If the context is done (canceled), return error
		return FileNodeConnection{}, errors.New("file was not found in network")
	}
}
