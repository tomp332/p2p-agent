package src

import (
	"context"
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
	ID             string
	Context        context.Context
	wg             sync.WaitGroup
	ConnectedPeers map[string]pb.FileSharingClient
	server         *FileSharingServer
}

func (n *Node) StartServer() {
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", n.ServerPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := &FileSharingServer{}
		s.Server = grpc.NewServer()
		pb.RegisterFileSharingServer(s.Server, s)
		log.Printf("gRPC server listening on port %d", n.ServerPort)
		if err := s.Server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (n *Node) StopServer() {
	log.Println("Stopping gRPC server...")
	n.server.Server.GracefulStop()
	n.wg.Wait()
	log.Println("gRPC server stopped.")
}

func (n *Node) handleBootstrap() {
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
				c, err := n.connectToNode(&addr)
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
		log.Fatalf("Bootstraping process failed, expected %d bootstrapped nodes,"+
			" ended up with %d", len(n.BootstrapPeerAddrs), len(n.ConnectedPeers))
	}
}

func NewNode(options *NodeOptions) *Node {
	node := &Node{
		ID:             GenerateRandomID(),
		Context:        context.Background(),
		NodeOptions:    *options,
		wg:             sync.WaitGroup{},
		ConnectedPeers: make(map[string]pb.FileSharingClient),
	}
	if options.BootstrapPeerAddrs != nil {
		node.handleBootstrap()
	}
	log.Println("Created new Node with ID: " + node.ID)
	return node
}

func (n *Node) connectToNode(address *string) (pb.FileSharingClient, error) {
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
		return nil, fmt.Errorf("failed to validate health for node server: %s", *address)
	}
	return c, nil
}
