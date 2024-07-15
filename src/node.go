package src

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/tomp332/p2fs/src/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"sync"
	"time"
)

type newPeer struct {
	address string
	client  pb.FileSharingClient
}

type NodeOptions struct {
	Storage              map[string][]byte
	ServerPort           int32
	BootstrapPeerAddrs   []string
	BootstrapNodeTimeout time.Duration
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
		ID:             GenerateRandomID(),
		Context:        ctx,
		NodeOptions:    *options,
		wg:             sync.WaitGroup{},
		ConnectedPeers: make(map[string]pb.FileSharingClient),
		cancelFunc:     cancel,
	}
	if options.BootstrapPeerAddrs != nil {
		if err := node.handleBootstrap(); err != nil {
			Logger.Error().Err(err).Msg("handleBootstrap failed")
			log.Fatal(err)
		}
	}
	Logger.Info().Str("nodeID", node.ID).Msg("Created new Node.")
	return node
}

// StartServer starts the gRPC server.
func (n *Node) StartServer() error {
	errorChan := make(chan error, 1)

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", n.ServerPort))
		if err != nil {
			errorChan <- err
			return
		}
		n.fileSharingServer = &FileSharingServer{}
		n.fileSharingServer.grpcServer = grpc.NewServer()
		pb.RegisterFileSharingServer(n.fileSharingServer.grpcServer, n.fileSharingServer)
		Logger.Info().Int32("serverPort", n.ServerPort).Msg("Started gRPC Node server.")
		if err := n.fileSharingServer.grpcServer.Serve(lis); err != nil {
			errorChan <- fmt.Errorf("failed to serve: %v", err)
		}
	}()

	select {
	case err := <-errorChan:
		return err
	case <-time.After(time.Second * 1): // Adjust timeout as needed
		return nil
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
	// Perform a health check to determine if the boostrap node is valid.
	healthCtx, healthCancel := context.WithTimeout(context.Background(), n.BootstrapNodeTimeout)
	defer healthCancel()

	_, err = c.HealthCheck(healthCtx, &pb.HealthCheckRequest{})
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		} // Close the connection if the Healthcheck fails
		return nil, fmt.Errorf("failed to validate health for node fileSharingServer: %s", *address)
	}
	return c, nil
}

func (n *Node) Terminate() error {
	if err := n.stopServer(); err != nil {
		return err
	}
	n.cancelFunc()
	return nil
}

func (n *Node) stopServer() error {
	if n.fileSharingServer != nil {
		Logger.Debug().Msg("Stopping gRPC Node server.")
		n.fileSharingServer.grpcServer.GracefulStop()
		n.wg.Wait()
		Logger.Info().Msg("gRPC gRPC server stopped.")
		return nil
	}
	return errors.New("file sharing server has been terminated")
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
					log.Printf("failed to connect to bootstrap node %s: %v", addr, err)
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
