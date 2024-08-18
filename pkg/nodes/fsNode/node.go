package fsNode

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
	"github.com/tomp332/p2p-agent/pkg/utils/types"
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

func NewP2PFilesNode(baseNode *nodes.BaseNode, storage storage.Storage) *FileNode {
	baseNode.ProtectedRoutes = []string{
		"/p2p_agent.FilesNodeService/DownloadFile",
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
	transferChan := make(chan types.TransferChunkData)
	ctx := stream.Context()

	go func() {
		defer close(transferChan)

		for {
			req, err := stream.Recv()
			if err == io.EOF {
				// End of stream, stop receiving
				return
			}
			if err != nil {
				// Handle the error if it occurs during reception
				log.Error().Err(err).Msg("Error receiving chunk from stream")
				return
			}
			chunkData := req.GetChunkData()
			log.Debug().Int("chunkSize", len(chunkData)).Msg("Received upload file data chunk")
			buffer.Write(chunkData)
			transferChan <- types.TransferChunkData{
				ChunkData: chunkData,
				ChunkSize: int64(len(chunkData)),
			}
		}
	}()

	// Buffer all data before proceeding
	fileID, _ := createUniqueFileID(&buffer)

	// Put a file into storage and handle context cancellation
	fileSize, storageErr := n.Storage.Put(ctx, fileID, transferChan)
	if storageErr != nil {
		// If the error is due to context cancellation, handle it gracefully
		if errors.Is(storageErr, ctx.Err()) {
			log.Info().Str("fileId", fileID).Msg("Upload was cancelled")
		} else {
			// Handle other types of errors
			log.Error().Err(storageErr).Msg("Failed to store file")
		}

		// Send an error response to the client
		return stream.SendAndClose(&pb.UploadFileResponse{
			FileId:  fileID,
			Success: false,
			Message: fmt.Sprintf("Failed to store file: %v", storageErr),
		})
	}

	// Successfully stored file
	log.Info().Str("fileId", fileID).Msg("File uploaded successfully.")
	return stream.SendAndClose(&pb.UploadFileResponse{
		FileId:   fileID,
		FileSize: fileSize,
		Success:  true,
		Message:  "File uploaded successfully",
	})
}

func (n *FileNode) DownloadFile(req *pb.DownloadFileRequest, stream pb.FilesNodeService_DownloadFileServer) error {
	ctx := stream.Context()

	// Attempt to retrieve the file from local storage first
	dataChan, storageErr := n.Storage.Get(ctx, req.FileId)
	if storageErr == nil {
		log.Debug().Str("fileId", req.FileId).Msg("File found locally, downloading from local storage")
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					// Data channel is closed, end of data transfer
					log.Debug().Str("fileId", req.FileId).Msg("Data channel closed, end of file transfer")
					return nil // End the stream by returning from the function
				}
				// Send the data to the DownloadFile stream
				if err := stream.Send(&pb.DownloadFileResponse{
					FileId:    req.FileId,
					Exists:    true,
					Chunk:     data.ChunkData,
					ChunkSize: data.ChunkSize,
				}); err != nil {
					log.Error().Err(err).Str("fileId", req.FileId).Msg("Error sending file data to stream")
					return err
				}
				log.Debug().Str("fileId", req.FileId).Int64("chunkSize", data.ChunkSize).Msg("Sent file data to stream")
			case <-ctx.Done():
				// Handle context cancellation
				log.Info().Str("fileId", req.FileId).Msg("Context cancelled, stopping file download")
				return ctx.Err()
			}
		}
	} else {
		log.Debug().
			Str("fileId", req.GetFileId()).
			Int("connectedPeers", len(n.ConnectedPeers)).
			Msg("File not found locally, searching in peer network")
		return n.searchPeers(ctx, req.GetFileId(), stream)
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
			authenticate, err := client.Authenticate(context.Background(), &pb.AuthenticateRequest{
				Username: connectionInfo.Username,
				Password: connectionInfo.Password,
			})
			if err != nil {
				log.Error().Err(err).Str("address", fmt.Sprintf("%s:%v", connectionInfo.Host, connectionInfo.Port)).Msg("")
				errChan <- err
				return
			}
			n.ConnectedPeers = append(n.ConnectedPeers, nodes.PeerConnection{
				ConnectionInfo: connectionInfo,
				GrpcConnection: con,
				NodeClient:     NewFileNodeClient(client, n.BootstrapNodeTimeout),
				Token:          authenticate.Token,
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

func (n *FileNode) searchPeers(ctx context.Context, fileId string, stream pb.FilesNodeService_DownloadFileServer) error {
	var wg sync.WaitGroup

	// Iterate over connected peers to broadcast a search query
	if len(n.ConnectedPeers) == 0 {
		log.Debug().Str("fileId", fileId).Msg("No connected peers available, file not found")
		return errors.New("file was not found locally and there are no connected peers")
	}
	for _, peer := range n.ConnectedPeers {
		wg.Add(1)
		go func(peer *nodes.PeerConnection) {
			defer wg.Done()
			log.Debug().
				Str("peerAddress", fmt.Sprintf("%s:%v", peer.ConnectionInfo.Host, peer.ConnectionInfo.Port)).
				Msg("Sending file search request to peer")
			dynamicPeerClient, ok := peer.NodeClient.(*FileNodeClient)
			if !ok {
				log.Warn().
					Str("peer", fmt.Sprintf("%s:%v", peer.ConnectionInfo.Host, peer.ConnectionInfo.Port)).
					Msg("Failed to use the current peer as a FileNodeClient.")
				return
			}
			peerDataChan, downloadErrChan := dynamicPeerClient.DownloadFile(ctx, fileId)

			// Read from the peer in loop
			for {
				select {
				case data, ok := <-peerDataChan:
					if !ok {
						log.Debug().Str("peerAddress", fmt.Sprintf("%s:%v", peer.ConnectionInfo.Host, peer.ConnectionInfo.Port)).Msg("Peer data channel closed, ending reception")
						return
					}
					log.Debug().
						Str("peerAddress", fmt.Sprintf("%s:%v", peer.ConnectionInfo.Host, peer.ConnectionInfo.Port)).
						Msg("Received chunk from remote peer")
					if sendErr := stream.Send(&pb.DownloadFileResponse{
						FileId:    fileId,
						Exists:    true,
						Chunk:     data.ChunkData,
						ChunkSize: data.ChunkSize,
					}); sendErr != nil {
						log.Error().Err(sendErr).Str("fileId", fileId).Msg("Error sending file data to stream")
						return
					}
				case err, ok := <-downloadErrChan:
					if !ok && err != nil {
						log.Error().
							Str("peerAddress", fmt.Sprintf("%s:%v", peer.ConnectionInfo.Host, peer.ConnectionInfo.Port)).
							Msg("Failed to download file from remote peer")
					}
					if err != nil {
						log.Error().Err(err).
							Str("peerAddress", fmt.Sprintf("%s:%v", peer.ConnectionInfo.Host, peer.ConnectionInfo.Port)).
							Msg("Failed to download file from remote peer")
						return
					}
				case <-ctx.Done():
					log.Info().Str("fileId", fileId).Msg("Context done, cancelling search")
					return
				}
			}
		}(&peer)
	}

	// Wait for all goroutines to finish once we have a result or a timeout
	wg.Wait()
	return nil
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
