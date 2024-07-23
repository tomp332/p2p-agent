package p2p

import (
	"context"
	"fmt"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

type P2PNode interface {
	Stop() error
	Register(server *grpc.Server)
	GetID() string
	GetType() src.NodeType
	ConnectToPeers()
}

type NodeConnection struct {
	address        string
	grpcConnection *grpc.ClientConn
	client         any
}

type BaseNode struct {
	configs.P2PNodeBaseConfig
	ID             string
	NodeType       src.NodeType
	Context        context.Context
	wg             sync.WaitGroup
	ConnectedPeers []NodeConnection
	cancelFunc     context.CancelFunc
}

func NewBaseNode(options *configs.P2PNodeBaseConfig) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	var id string
	if id = options.ID; id == "" {
		id = utils.GenerateRandomID()
	}
	node := &BaseNode{
		ID:                id,
		Context:           ctx,
		P2PNodeBaseConfig: *options,
		wg:                sync.WaitGroup{},
		ConnectedPeers:    make([]NodeConnection, 0),
		cancelFunc:        cancel,
	}
	return node
}

func (n *BaseNode) ConnectToBootstrapNodes() {
	for _, nodeAddress := range n.BootstrapPeerAddrs {
		_, err := n.connectToNode(nodeAddress)
		if err != nil {
			utils.Logger.Warn().Str("nodeAddress", nodeAddress).Err(err).Msg("")
		}
	}
}

func (n *BaseNode) Terminate() error {
	n.cancelFunc()
	for _, conn := range n.ConnectedPeers {
		conn.grpcConnection.Close()
	}
	return nil
}

func (n *BaseNode) GetID() string {
	return n.ID
}

func (n *BaseNode) GetType() src.NodeType {
	return n.NodeType
}

func (n *BaseNode) connectToNode(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewNodeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.P2PNodeBaseConfig.BootstrapNodeTimeout*time.Second)
	defer cancel()
	_, err = client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to connect to bootstrap node")
	}
	n.ConnectedPeers = append(n.ConnectedPeers, NodeConnection{
		address:        address,
		grpcConnection: conn,
	})
	return conn, nil
}
