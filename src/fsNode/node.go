package fsNode

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tomp332/p2fs/src/common"
	pb "github.com/tomp332/p2fs/src/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

type newPeer struct {
	address string
	client  pb.FileSharingClient
}

type NodeOptions struct {
	Storage              map[string][]byte
	BootstrapPeerAddrs   []string
	BootstrapNodeTimeout time.Duration
	ServerOptions
}

type Node struct {
	NodeOptions
	ID                string
	Context           context.Context
	wg                sync.WaitGroup
	ConnectedPeers    map[string]pb.FileSharingClient
	fileSharingServer *FileSharingServer
	cancelFunc        context.CancelFunc
}

func NewNode(options *NodeOptions) *Node {
	ctx, cancel := context.WithCancel(context.Background())
	node := &Node{
		ID:             common.GenerateRandomID(),
		Context:        ctx,
		NodeOptions:    *options,
		wg:             sync.WaitGroup{},
		ConnectedPeers: make(map[string]pb.FileSharingClient),
		cancelFunc:     cancel,
	}
	node.fileSharingServer = NewFileSharingServer(node)
	err := node.fileSharingServer.StartServer()
	if err != nil {
		common.Logger.Fatal().Msg(err.Error())
	}
	if options.BootstrapPeerAddrs != nil {
		if err := node.handleBootstrap(); err != nil {
			common.Logger.Error().Err(err).Msg("handleBootstrap failed")
			log.Fatal(err)
		}
	}
	common.Logger.Info().Str("nodeID", node.ID).Msg("Created new Node.")
	return node
}

func (n *Node) Terminate() {
	if n.fileSharingServer != nil {
		n.fileSharingServer.StopServer()
		n.cancelFunc()
	}
}

func (n *Node) ConnectToNode(address *string) (pb.FileSharingClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(*address, opts...)
	if err != nil {
		return nil, err
	}
	c := pb.NewFileSharingClient(conn)
	// Perform a health check to determine if the boostrap fsNode is valid.
	healthCtx, healthCancel := context.WithTimeout(context.Background(), n.BootstrapNodeTimeout)
	defer healthCancel()

	_, err = c.HealthCheck(healthCtx, &pb.HealthCheckRequest{})
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		} // Close the connection if the Healthcheck fails
		return nil, fmt.Errorf("failed to validate health for fsNode fileSharingServer: %s", *address)
	}
	return c, nil
}

func (n *Node) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	hash := sha256.Sum256(req.Payload)
	fileId := hex.EncodeToString(hash[:])
	n.Storage[fileId] = req.Payload
	return &pb.UploadFileResponse{
		FileId: fileId,
	}, nil
}

func (n *Node) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	chunk, exists := n.Storage[req.FileId]
	return &pb.DownloadFileResponse{
		Chunk:  chunk,
		FileId: req.FileId,
		Exists: exists,
	}, nil
}

func (n *Node) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	delete(n.Storage, req.FileId)
	return &pb.DeleteFileResponse{
		FileId: req.FileId,
	}, nil
}

func (n *Node) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status: pb.HealthStatus_ALIVE,
	}, nil
}

func (n *Node) handleBootstrap() error {
	newPeerUpdates := make(chan newPeer)
	go func() {
		for update := range newPeerUpdates {
			n.ConnectedPeers[update.address] = update.client
		}
	}()
	// Connect to bootstrap nodes
	if n.BootstrapPeerAddrs != nil {
		for _, address := range n.BootstrapPeerAddrs {
			n.wg.Add(1)
			go func(addr string) {
				defer n.wg.Done()
				c, err := n.ConnectToNode(&addr)
				if err != nil {
					log.Printf("failed to connect to bootstrap fsNode %s: %v", addr, err)
					return
				}
				newPeerUpdates <- newPeer{address: addr, client: c}
			}(address)
		}
	}
	n.wg.Wait()
	close(newPeerUpdates)
	if len(n.ConnectedPeers) != len(n.BootstrapPeerAddrs) {
		return errors.New(fmt.Sprintf("Bootstraping process failed, expected %d bootstrapped nodes,ended up with %d", len(n.BootstrapPeerAddrs), len(n.ConnectedPeers)))
	}
	return nil
}
