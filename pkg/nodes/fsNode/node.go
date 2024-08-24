package fsNode

import (
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
		"/p2pAgent.FilesNodeService/DownloadFile",
		"/p2pAgent.FilesNodeService/UploadFile",
		"/p2pAgent.FilesNodeService/DeleteFile",
	}
	n := &FileNode{
		BaseNode: baseNode,
		Storage:  storage,
	}
	n.createAuthInformation()
	return n
}

func (n *FileNode) InterceptorsRegister(server *server.GRPCServer) {
	fileNodeAuthInterceptor := interceptors.NewAuthInterceptor(n.AuthManager, n.ProtectedRoutes)
	server.AddUnaryInterceptors(fileNodeAuthInterceptor.Unary())
	server.AddStreamInterceptors(fileNodeAuthInterceptor.Stream())
}

func (n *FileNode) UploadFile(stream pb.FilesNodeService_UploadFileServer) error {
	ctx := stream.Context()
	var fileSize uint64
	fileHash := md5.New()
	var fileName string

	// Channels
	fileDataChan := make(chan *types.TransferChunkData)
	fileNameChan := make(chan string, 1)
	errChan := make(chan error, 1)
	defer close(fileNameChan)

	var wg sync.WaitGroup

	// Start goroutine to consume stream
	wg.Add(1)
	go func() {
		defer close(fileDataChan)
		defer wg.Done()

		for {
			req, err := stream.Recv()
			if err == io.EOF {
				log.Debug().Str("fileName", fileName).Msg("Done processing file upload stream.")
				return // Exit normally on EOF
			}
			if err != nil {
				log.Error().Err(err).Msg("Error receiving file chunk")
				select {
				case errChan <- err:
				default:
				}
				return
			}
			if req.FileName == "" {
				log.Error().Msg("No file name was specified in the upload request")
				select {
				case errChan <- errors.New("no file name was specified in the upload request"):
				default:
				}
				return
			}
			if fileName == "" {
				fileName = req.FileName
				fileNameChan <- req.FileName
			}

			fileSize += uint64(len(req.GetChunkData()))
			if _, hashErr := fileHash.Write(req.GetChunkData()); hashErr != nil {
				log.Error().Err(hashErr).Msg("Error writing new file chunk to md5 hash")
				select {
				case errChan <- hashErr:
				default:
				}
				return
			}
			log.Debug().Str("fileName", req.FileName).Int("chunkSize", len(req.GetChunkData())).Msg("Received file chunk")

			// Pass the file data to the processing channel
			fileDataChan <- &types.TransferChunkData{
				ChunkSize: len(req.GetChunkData()),
				ChunkData: req.GetChunkData(),
			}
		}
	}()

	// Start goroutine to save file to storage
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case requestedFile := <-fileNameChan:
			err := n.Storage.Put(ctx, requestedFile, fileDataChan, &wg)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		case <-ctx.Done():
			// Handle context cancellation
			select {
			case errChan <- ctx.Err():
			default:
			}
			return
		}
	}()

	// Use a select statement to immediately handle the first error encountered
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-func() chan struct{} {
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()
		return done
	}():
		// All goroutines completed without error
	}

	log.Info().Str("fileName", fileName).Uint64("fileSize", fileSize).Msg("Uploaded file successfully")
	err := stream.SendAndClose(&pb.UploadFileResponse{
		FileName: fileName,
		FileSize: fileSize,
		FileHash: hex.EncodeToString(fileHash.Sum(nil)),
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *FileNode) DownloadFile(req *pb.DownloadFileRequest, stream pb.FilesNodeService_DownloadFileServer) error {
	streamCtx := stream.Context()

	// Error and response channels
	errChan := make(chan error, 1)
	respChan := make(chan *pb.DownloadFileResponse)

	// Start a goroutine to listen for context cancellation
	go func() {
		defer close(errChan)
		<-streamCtx.Done()
		if err := streamCtx.Err(); err != nil {
			select {
			case errChan <- err:
			default:
				// Avoid sending to a closed channel
			}
		}
	}()

	// Search for the file
	searchResponse, searchErr := n.SearchFile(streamCtx, &pb.SearchFileRequest{
		FileName: req.FileName,
	})

	if searchErr != nil {
		log.Error().Err(searchErr).Msg("Failed to search for file in network peer")
		return searchErr
	}

	// Handle file download based on search result
	go func() {
		defer close(respChan) // Close respChan when done
		if searchResponse.NodeId == configs.MainConfig.ID {
			n.downloadLocal(req.FileName, streamCtx, respChan, errChan)
		} else {
			n.downloadRemote(searchResponse.NodeId, req.FileName, streamCtx, respChan, errChan)
		}
	}()

	// Handle sending chunks to the client
	for {
		select {
		case data, ok := <-respChan:
			if !ok {
				// Response channel is closed, exit the loop
				log.Debug().Str("fileName", req.FileName).Msg("Finished handling download file request")
				return nil
			}
			log.Debug().Int32("chunkSize", data.ChunkSize).Str("fileName", req.FileName).Msg("Sending file chunk to client")
			if err := stream.Send(data); err != nil {
				log.Error().Err(err).Str("fileName", req.FileName).Msg("Unable to send download file stream data to client.")
				return err
			}
		case err := <-errChan:
			if err != nil {
				return err
			}
		}
	}
}

func (n *FileNode) SearchFile(ctx context.Context, req *pb.SearchFileRequest) (*pb.SearchFileResponse, error) {
	log.Debug().Str("fileName", req.GetFileName()).Msg("Starting to search for file")
	exists := n.Storage.Exists(req.GetFileName())
	if !exists {
		return n.searchPeers(ctx, req.GetFileName())
	}
	return &pb.SearchFileResponse{
		FileName: req.FileName,
		NodeId:   configs.MainConfig.ID,
	}, nil
}

func (n *FileNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	err := n.Storage.Delete(ctx, request.GetFileName())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{FileName: request.GetFileName()}, nil
}

func (n *FileNode) ServiceRegister(server *server.GRPCServer) {
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
			fsNodeClient := NewFileNodeClient(pb.NewFilesNodeServiceClient(con), &n.ProtectedRoutes, connectionInfo, n.BootstrapNodeTimeout)
			nodeId, err := fsNodeClient.Connect()
			if err != nil {
				log.Error().Err(err).Str("address", fmt.Sprintf("%s:%v", connectionInfo.Host, connectionInfo.Port)).Msg("")
				errChan <- err
				return
			}
			n.ConnectedPeers[nodeId] = nodes.PeerConnection{
				ID:             nodeId,
				ConnectionInfo: connectionInfo,
				NodeClient:     fsNodeClient,
			}
			log.Debug().
				Str("address", fmt.Sprintf("%s:%v", connectionInfo.Host, connectionInfo.Port)).
				Str("nodeId", nodeId).
				Msgf("Successfully connected to bootstrap node.")
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
	return &pb.AuthenticateResponse{Token: token, NodeId: configs.MainConfig.ID}, nil
}

func (n *FileNode) searchPeers(parentCtx context.Context, fileName string) (*pb.SearchFileResponse, error) {
	var wg sync.WaitGroup
	searchCtx, cancelSearch := context.WithCancel(parentCtx)
	defer cancelSearch() // Ensure the search context is canceled when the function exits

	if len(n.ConnectedPeers) == 0 {
		log.Debug().Str("fileName", fileName).Msg("No connected peers available, file not found")
		return nil, errors.New("file was not found locally, and there are no connected peers")
	}

	// Channel to receive the first successful search response
	resultChan := make(chan *pb.SearchFileResponse, 1)
	errorChan := make(chan error, 1) // Channel to handle errors

	for _, peer := range n.ConnectedPeers {
		wg.Add(1)
		go func(peer *nodes.PeerConnection) {
			defer wg.Done()
			dynamicPeerClient, ok := peer.NodeClient.(*FileNodeClient)
			if !ok {
				log.Warn().Str("peer", fmt.Sprintf("%s:%v", peer.ConnectionInfo.Host, peer.ConnectionInfo.Port)).Msg("Failed to use the current peer as a FileNodeClient.")
				return
			}

			searchResp, searchErr := dynamicPeerClient.SearchFile(searchCtx, fileName)
			if searchErr != nil {
				errorChan <- searchErr
				return
			}

			if searchResp != nil {
				select {
				case resultChan <- searchResp:
					// Successfully sent result to channel
					cancelSearch() // Cancel other searches
				case <-searchCtx.Done():
					// Context canceled before sending result
					log.Info().Msg("Search context canceled before sending result")
				}
			}
		}(&peer)
	}

	// Wait for either result or context cancellation
	select {
	case foundFile := <-resultChan:
		// Successfully received a search response
		remoteNode := n.ConnectedPeers[foundFile.NodeId]
		log.Debug().
			Str("nodeId", foundFile.NodeId).
			Str("remotePeer", fmt.Sprintf("%s:%v", remoteNode.ConnectionInfo.Host, remoteNode.ConnectionInfo.Port)).
			Msg("Found file in peer network")
		return foundFile, nil
	case err := <-errorChan:
		// Error occurred during search
		log.Error().Err(err).Msg("Error during file search")
		return nil, err
	case <-searchCtx.Done():
		// Search context was canceled
		log.Debug().Msg("Search context canceled or no file found in peer network")
	}

	// Ensure all goroutines have completed
	wg.Wait()
	return nil, errors.New("file not found in peer network")
}

func (n *FileNode) downloadLocal(fileName string, streamCtx context.Context, respChan chan *pb.DownloadFileResponse, errChan chan error) {
	log.Debug().Str("fileId", fileName).Msg("File found locally, downloading from local storage")
	dataChan, localFileErr := n.Storage.Get(streamCtx, fileName)
	if localFileErr != nil {
		log.Error().Err(localFileErr).Msg("Failed to download local file")
		select {
		case errChan <- localFileErr:
		default:
			// Avoid sending to a closed channel
		}
		return
	}

	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				// Data channel is closed, end of data transfer
				log.Debug().Str("fileName", fileName).Msg("Data channel closed, end of file transfer")
				return
			}
			log.Debug().
				Str("fileName", fileName).
				Int("chunkSize", data.ChunkSize).
				Msg("Received file chunk from storage")
			respChan <- &pb.DownloadFileResponse{
				FileName:  fileName,
				Chunk:     data.ChunkData,
				ChunkSize: int32(data.ChunkSize),
			}
		case <-streamCtx.Done():
			log.Debug().Str("fileName", fileName).Msg("Stream context canceled during file download")
			return
		}
	}
}

func (n *FileNode) downloadRemote(nodeId string, fileName string, streamCtx context.Context, respChan chan *pb.DownloadFileResponse, errChan chan error) {
	// Download file from remote peer
	remoteNode, exists := n.ConnectedPeers[nodeId]
	if !exists {
		log.Error().Str("nodeId", nodeId).Msg("Failed to fetch connected peer by provided nodeId")
		errChan <- fmt.Errorf("failed to fetch connected peer by provided nodeId %s", nodeId)
		return
	}
	dynamicPeerClient, ok := remoteNode.NodeClient.(*FileNodeClient)
	if !ok {
		log.Warn().
			Str("peer", fmt.Sprintf("%s:%v", remoteNode.ConnectionInfo.Host, remoteNode.ConnectionInfo.Port)).
			Msg("Failed to use the current peer as a FileNodeClient.")
		errChan <- fmt.Errorf("failed to use the current peer as a FileNodeClient")
		return
	}

	peerDataChan, downloadErrChan := dynamicPeerClient.DownloadFile(streamCtx, fileName)

	for {
		select {
		case data, ok := <-peerDataChan:
			if !ok {
				log.Debug().Str("peerAddress",
					fmt.Sprintf("%s:%v", remoteNode.ConnectionInfo.Host, remoteNode.ConnectionInfo.Port)).
					Msg("Peer data channel closed, ending reception")
				return // Channel closed
			}

			log.Debug().
				Int("chunkSize", data.ChunkSize).
				Str("fileName", fileName).
				Str("peerAddress", fmt.Sprintf("%s:%v", remoteNode.ConnectionInfo.Host, remoteNode.ConnectionInfo.Port)).
				Msg("Received file chunk from network peer")

			respChan <- &pb.DownloadFileResponse{
				FileName:  fileName,
				Chunk:     data.ChunkData,
				ChunkSize: int32(data.ChunkSize),
			}

		case err := <-downloadErrChan:
			if err != nil {
				log.Error().
					Err(err).
					Str("peerAddress", fmt.Sprintf("%s:%v", remoteNode.ConnectionInfo.Host, remoteNode.ConnectionInfo.Port)).
					Msg("Failed to download file from remote peer")
				errChan <- err
				return
			}
		}
	}
}

func (n *FileNode) createAuthInformation() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(n.Auth.Password), bcrypt.DefaultCost)
	// Make sure its stored decoded.
	n.Auth.Password = string(hashedPassword)
	n.AuthManager = managers.NewJWTManager(n.Auth.JwtSecret, 24*time.Hour)
}
