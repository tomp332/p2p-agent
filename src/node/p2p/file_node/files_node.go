package file_node

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/storage"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
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

// UploadFile handles client-streaming RPC for file upload
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
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		foundFileNode := n.searchFileInNetwork(req.FileId)
		if foundFileNode.Address != "" {
			// Found file in the network fetch from remote
			remoteDataChan, errChan := foundFileNode.NodeClient.DownloadFile(ctx, req.GetFileId())
			if errChan != nil {
				return <-errChan
			}
			mainDataChan = remoteDataChan
		} else {
			dataChan, err := n.Storage.Get(ctx, req.FileId)
			if err != nil {
				return err
			}
			mainDataChan = dataChan
		}
		for chunk := range mainDataChan {
			err := stream.Send(&pb.DownloadFileResponse{
				FileId: req.FileId,
				Exists: true,
				Chunk:  chunk,
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (n *FilesNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	err := n.Storage.Delete(ctx, request.GetFileId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{FileId: request.GetFileId()}, nil
}

func (n *FilesNode) Register() {
	pb.RegisterFilesNodeServiceServer(n.Server.ServerObj(), n)
	utils.Logger.Info().Str("nodeType", n.Type).Msg("Node registered")
}

func (n *FilesNode) ConnectToBootstrapPeers() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(n.BootstrapPeerAddrs))

	for _, address := range n.BootstrapPeerAddrs {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			con, err := n.ConnectToPeer(address, n.P2PNodeBaseConfig.BootstrapNodeTimeout)
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

func (n *FilesNode) searchFileInNetwork(fileId string) FileNodeConnection {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileFoundChan := make(chan FileNodeConnection, 1) // Buffer size 1 to avoid blocking
	// Send search RPC call to all connected peers,to see
	//if the network already contains this file.
	for _, peer := range n.ConnectedPeers {
		wg.Add(1)
		go func(peer src.P2PNodeConnection) {
			dynamicPeer := peer.(FileNodeConnection)
			defer wg.Done()
			result, err := dynamicPeer.NodeClient.SearchFile(ctx, fileId)
			if err != nil {
				// Handle error (for example, log it)
				utils.Logger.Debug().Str("peer", dynamicPeer.Address).Msgf("Failed to send SearchFile rpc call for peer.")
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
		return fileNode
	case <-ctx.Done(): // If the context is done (canceled), return false
		return FileNodeConnection{}
	}
}
