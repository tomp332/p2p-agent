package p2p

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

type ConfigurableP2PNode interface {
	ParseNodeConfig(specificConfig map[string]interface{}) error
}

type P2PNode interface {
	ConfigurableP2PNode
	Stop() error
	Register(server *grpc.Server)
	GetID() string
	GetType() src.NodeType
}

type P2PNodeConfig struct {
	Type                 string                 `json:"type"`
	BootstrapPeerAddrs   []string               `json:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration          `json:"bootstrap_node_timeout"`
	ExtraConfig          map[string]interface{} `json:"extra_config"`
}

type BaseNode struct {
	P2PNodeConfig
	ID             string
	NodeType       src.NodeType
	Context        context.Context
	wg             sync.WaitGroup
	ConnectedPeers map[string]*grpc.ClientConn
	cancelFunc     context.CancelFunc
}

func NewBaseNode(options *P2PNodeConfig) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	node := &BaseNode{
		ID:             utils.GenerateRandomID(),
		Context:        ctx,
		P2PNodeConfig:  *options,
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
