package nodes

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/server"
	"github.com/tomp332/p2p-agent/pkg/utils"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"
)

type P2PNoder interface {
	Register(server *grpc.Server)
	Options() *configs.NodeConfig
	ConnectToBootstrapPeers(server *server.GRPCServer) error
}

type P2PNodeClienter interface{}

type P2PNodeConnection interface{}

type NodeConnection struct {
	P2PNodeConnection
	Address        string
	GrpcConnection *grpc.ClientConn
	NodeClient     interface{}
}

type BaseNode struct {
	configs.NodeConfig
	ConnectedPeers []P2PNodeConnection
	Context        context.Context
	wg             sync.WaitGroup
	cancelFunc     context.CancelFunc
}

func NewBaseNode(options *configs.NodeConfig) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	var id string
	if id = options.ID; id == "" {
		options.ID = utils.GenerateRandomID()
	}
	n := &BaseNode{
		NodeConfig:     *options,
		ConnectedPeers: make([]P2PNodeConnection, 0),
		Context:        ctx,
		wg:             sync.WaitGroup{},
		cancelFunc:     cancel,
	}
	return n
}

func (n *BaseNode) Register(_ *grpc.Server) {}

func (n *BaseNode) Options() *configs.NodeConfig {
	return &n.NodeConfig
}

func (n *BaseNode) Terminate() error {
	n.cancelFunc()
	for _, conn := range n.ConnectedPeers {
		if dynamicField, ok := conn.(NodeConnection); ok {
			dynamicField.GrpcConnection.Close()
		}
	}
	return nil
}

func (n *BaseNode) ConnectToBootstrapPeers(server *grpc.Server) error {
	for _, address := range n.BootstrapPeerAddrs {
		conn, err := n.ConnectToPeer(server, address, n.BootstrapNodeTimeout)
		if err != nil {
			return err
		}
		n.ConnectedPeers = append(n.ConnectedPeers, NodeConnection{
			Address:        address,
			GrpcConnection: conn,
		})
	}
	return nil
}

func (n *BaseNode) ConnectToPeer(server *grpc.Server, address string, timeout time.Duration) (*grpc.ClientConn, error) {
	connection, err := clientConnection(address)
	if err != nil {
		log.Warn().Err(err).Str("Address", address).Msg("Failed to connect to peers")
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	log.Debug().Str("nodeType", n.Type.ToString()).Msg("Connecting to bootstrap peer")
	healthClient := grpc_health_v1.NewHealthClient(connection)
	res, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return nil, err
	}
	if res.Status == grpc_health_v1.HealthCheckResponse_SERVING {
		log.Debug().Str("healthStatus", res.Status.String()).Msgf("Bootstrap peer nodes is healthy.")
	} else {
		log.Debug().Str("healthStatus", res.Status.String()).Msgf("Bootstrap peer nodes is not healthy.")
		return nil, fmt.Errorf("bootstrap peer nodes is not healthy")
	}
	return connection, nil
}

func clientConnection(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		log.Warn().Err(err).Str("address", address).Msg("Failed to connect to peers")
		return nil, err
	}
	return conn, nil
}
