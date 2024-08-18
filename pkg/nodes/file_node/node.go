package file_node

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/nodes"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"github.com/tomp332/p2p-agent/pkg/server"
	"github.com/tomp332/p2p-agent/pkg/server/interceptors"
	"github.com/tomp332/p2p-agent/pkg/server/managers"
	"github.com/tomp332/p2p-agent/pkg/storage"
	"github.com/tomp332/p2p-agent/pkg/utils"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

type FileNode struct {
	*nodes.BaseNode
	Storage storage.Storage
	pb.UnimplementedFilesNodeServiceServer
}

type FileNodeConnection struct {
	nodes.NodeConnection
	NodeClient FileNodeClient
}

func NewP2PFilesNode(baseNode *nodes.BaseNode, storage storage.Storage) *FileNode {
	baseNode.ProtectedRoutes = []string{
		"/p2p_agent.FilesNodeService/DirectDownloadFile",
		"/p2p_agent.FilesNodeService/UploadFile",
	}
	n := &FileNode{
		BaseNode: baseNode,
		Storage:  storage,
	}
	n.createAuthInformation()
	fileNodeAuthInterceptor := interceptors.NewAuthInterceptor(n.AuthManager, baseNode.ProtectedRoutes)
	n.AddUnaryInterceptors(fileNodeAuthInterceptor.Unary())
	n.AddStreamInterceptors(fileNodeAuthInterceptor.Stream())
	return n
}

func (n *FileNode) UploadFile(stream pb.FilesNodeService_UploadFileServer) error {
	var buffer bytes.Buffer
	dataChan := make(chan []byte, 1024) // Buffered channel for file data

	ctx := stream.Context()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fileID, _ := createUniqueFileID(&buffer)
			close(dataChan)
			fileSize, storageErr := n.Storage.Put(ctx, fileID, dataChan)
			if storageErr != nil {
				return stream.SendAndClose(&pb.UploadFileResponse{
					FileId:  fileID,
					Success: false,
					Message: fmt.Sprintf("Failed to store file: %v", storageErr),
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

func (n *FileNode) DownloadFile(req *pb.DownloadFileRequest, stream pb.FilesNodeService_DownloadFileServer) error {
	ctx := stream.Context()
	var mainDataChan <-chan []byte
	var errChan <-chan error

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		foundFileNode, err := n.SearchFileInNetwork(req.FileId)
		if err != nil {
			log.Debug().Str("fileId", req.FileId).Msg("File not found in network, downloading from local")
			// File not found in the network, retrieve from storage
			dataChan, storageErr := n.Storage.Get(ctx, req.FileId)
			if storageErr != nil {
				return storageErr
			}
			mainDataChan = dataChan
		} else {
			log.Debug().Str("fileId", req.FileId).
				Str("peerAddress", fmt.Sprintf("%s:%v", foundFileNode.ConnectionInfo.Host, foundFileNode.ConnectionInfo.Port)).
				Msg("File found in network, downloading from remote peer")
			// Found the required file in the network, fetch from remote
			mainDataChan, errChan = foundFileNode.NodeClient.DirectDownloadFile(ctx, req.GetFileId(), n.ID)
		}
		// Use the common handler function
		return n.handleFileDownload(ctx, req.FileId, mainDataChan, errChan, stream)
	}
}

func (n *FileNode) DirectDownloadFile(req *pb.DirectDownloadFileRequest, stream pb.FilesNodeService_DirectDownloadFileServer) error {
	ctx := stream.Context()
	var mainDataChan <-chan []byte
	var errChan <-chan error

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Debug().
			Str("fileId", req.GetFileId()).
			Str("peerNodeId", req.GetNodeId()).
			Msg("Direct download request received, downloading from local")
		// Retrieve file from storage
		dataChan, storageErr := n.Storage.Get(ctx, req.GetFileId())
		if storageErr != nil {
			return storageErr
		}
		mainDataChan = dataChan

		// Use the common handler function
		return n.handleFileDownload(ctx, req.GetFileId(), mainDataChan, errChan, stream)
	}
}

func (n *FileNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	err := n.Storage.Delete(ctx, request.GetFileId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{FileId: request.GetFileId()}, nil
}

func (n *FileNode) Register(server *server.GRPCServer) {
	server.AddUnaryInterceptors(n.UnaryInterceptors...)
	server.AddStreamInterceptors(n.StreamInterceptors...)
	server.ServiceRegistrars[n.Type] = func(server *grpc.Server) {
		pb.RegisterFilesNodeServiceServer(server, n)
	}
}

func (n *FileNode) ConnectToBootstrapPeers() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(n.BootstrapPeerAddrs))

	for _, connectionInfo := range n.BootstrapPeerAddrs {
		wg.Add(1)
		go func(connectionInfo *configs.BootStrapNodeConnection) {
			defer wg.Done()
			con, err := n.ConnectToPeer(connectionInfo, n.NodeConfig.BootstrapNodeTimeout)
			if err != nil {
				log.Error().Err(err).Str("peer", fmt.Sprintf("%s:%v", connectionInfo.Host, connectionInfo.Port)).Msg("Failed to connect to bootstrap peer")
				errChan <- err
				return
			}
			client := pb.NewFilesNodeServiceClient(con)
			n.ConnectedPeers = append(n.ConnectedPeers, FileNodeConnection{
				NodeConnection: nodes.NodeConnection{
					ConnectionInfo: connectionInfo,
					GrpcConnection: con,
				},
				NodeClient: *NewFileNodeClient(client, n.BootstrapNodeTimeout),
			})
			log.Debug().
				Str("address", fmt.Sprintf("%s:%v", connectionInfo.Host, connectionInfo.Port)).
				Msgf("Successfully connected to bootstrap nodes.")
		}(&connectionInfo)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (n *FileNode) SearchFile(_ context.Context, request *pb.SearchFileRequest) (*pb.SearchFileResponse, error) {
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

func (n *FileNode) Authenticate(_ context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if req.GetUsername() != n.Auth.Username {
		log.Warn().Str("username", req.GetUsername()).Msg("Invalid username provided for node authentication")
		return nil, errors.New("invalid username/password")
	}
	if !utils.BcryptCompare(n.Auth.Password, req.GetPassword()) {
		log.Warn().Str("password", req.GetPassword()).Msg("Invalid password provided for node authentication")
		return nil, errors.New("invalid username/password")
	}
	token, err := n.AuthManager.Generate(n.Auth.Username, configs.FilesNodeType)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("Successfully authenticated peer")
	return &pb.AuthenticateResponse{Token: token}, nil
}

func (n *FileNode) SearchFileInNetwork(fileId string) (FileNodeConnection, error) {
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
		go func(peer nodes.P2PNodeConnection) {
			defer wg.Done()
			dynamicPeer := peer.(FileNodeConnection)
			fileFound, err := dynamicPeer.NodeClient.SearchFile(ctx, fileId)
			if err != nil {
				log.Debug().
					Str("peerAddress", fmt.Sprintf("%s:%v", dynamicPeer.ConnectionInfo.Host, dynamicPeer.ConnectionInfo.Port)).
					Msgf("Failed to send SearchFile RPC call for peer.")
				return
			}
			if fileFound {
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
		if fileNode.ConnectionInfo != nil {
			return fileNode, nil
		}
	case <-ctx.Done(): // If the context is done (canceled), return error
		return FileNodeConnection{}, errors.New("file was not found in network")
	}
	return FileNodeConnection{}, errors.New("file was not found in network")
}

func (n *FileNode) handleFileDownload(ctx context.Context, fileId string, mainDataChan <-chan []byte, errChan <-chan error, stream interface{}) error {
	for {
		select {
		case chunk, ok := <-mainDataChan:
			if !ok {
				// Channel closed, end of data
				return nil
			}

			// Use type switch to handle the specific stream type
			switch s := stream.(type) {
			case pb.FilesNodeService_DownloadFileServer:
				if streamErr := s.Send(&pb.DownloadFileResponse{
					FileId: fileId,
					Exists: true,
					Chunk:  chunk,
				}); streamErr != nil {
					return streamErr
				}
			case pb.FilesNodeService_DirectDownloadFileServer:
				if streamErr := s.Send(&pb.DirectDownloadFileResponse{
					FileId: fileId,
					Exists: true,
					Chunk:  chunk,
				}); streamErr != nil {
					return streamErr
				}
			default:
				return fmt.Errorf("unsupported stream type")
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

func createUniqueFileID(buffer *bytes.Buffer) (string, error) {
	fileHash := md5.Sum(buffer.Bytes())
	return hex.EncodeToString(fileHash[:]) + "-" + utils.GenerateRandomID(), nil
}

func (n *FileNode) createAuthInformation() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(n.Auth.Password), bcrypt.DefaultCost)
	// Make sure its stored decoded.
	n.Auth.Password = string(hashedPassword)
	n.AuthManager = managers.NewJWTManager(n.Auth.JwtSecret, 24*time.Hour)
}
