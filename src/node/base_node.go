package node

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeOptions struct {
	Storage              map[string][]byte
	BootstrapPeerAddrs   []string
	BootstrapNodeTimeout time.Duration
}

type Node interface {
	Start() error
	Stop() error
	Register(server *grpc.Server)
}

type BaseNode struct {
	NodeOptions
	ID             string
	Context        context.Context
	wg             sync.WaitGroup
	ConnectedPeers map[string]*grpc.ClientConn
	cancelFunc     context.CancelFunc
}

func NewBaseNode(options *NodeOptions) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	node := &BaseNode{
		ID:             util.GenerateRandomID(),
		Context:        ctx,
		NodeOptions:    *options,
		wg:             sync.WaitGroup{},
		ConnectedPeers: make(map[string]*grpc.ClientConn),
		cancelFunc:     cancel,
	}
	return node
}

func (n *BaseNode) ConnectToNode(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewNodeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to connect to node at %s: %w", address, err)
	}
	n.ConnectedPeers[address] = conn
	return conn, nil
}

func (n *BaseNode) Terminate() error {
	n.cancelFunc()
	for _, conn := range n.ConnectedPeers {
		conn.Close()
	}
	return nil
}
