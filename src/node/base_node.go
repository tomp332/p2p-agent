package node

import (
	"context"
	"fmt"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/utils"
	"sync"
	"time"

	"github.com/tomp332/p2p-agent/src/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type P2PNode interface {
	Stop() error
	Register(server *grpc.Server)
	GetID() string
	GetType() src.NodeType
}

type NodeOptions struct {
	Storage              map[string][]byte
	BootstrapPeerAddrs   []string
	BootstrapNodeTimeout time.Duration
}

type BaseNode struct {
	NodeOptions
	ID             string
	NodeType       src.NodeType
	Context        context.Context
	wg             sync.WaitGroup
	ConnectedPeers map[string]*grpc.ClientConn
	cancelFunc     context.CancelFunc
}

func NewBaseNode(options *NodeOptions) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	node := &BaseNode{
		ID:             utils.GenerateRandomID(),
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

func (n *BaseNode) GetID() string {
	return n.ID
}

func (n *BaseNode) GetType() src.NodeType {
	return n.NodeType
}
