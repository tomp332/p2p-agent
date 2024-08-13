package p2p

import (
	"context"
	"fmt"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"
)

type NodeConnection struct {
	src.P2PNodeConnection
	Address        string
	GrpcConnection *grpc.ClientConn
	NodeClient     interface{}
}

type BaseNode struct {
	configs.NodeConfigs
	Server         src.AgentGRPCServer
	ConnectedPeers []src.P2PNodeConnection
	Context        context.Context
	wg             sync.WaitGroup
	cancelFunc     context.CancelFunc
}

func NewBaseNode(options *configs.NodeConfigs) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	var id string
	if id = options.ID; id == "" {
		options.ID = utils.GenerateRandomID()
	}
	n := &BaseNode{
		NodeConfigs:    *options,
		ConnectedPeers: make([]src.P2PNodeConnection, 0),
		Context:        ctx,
		wg:             sync.WaitGroup{},
		cancelFunc:     cancel,
	}
	return n
}

func (n *BaseNode) Register(_ *grpc.Server) {
	utils.Logger.Info().Str("nodeType", n.Type).Msg("Node registered")
}

func (n *BaseNode) Options() *configs.NodeConfigs {
	return &n.NodeConfigs
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

func (n *BaseNode) ConnectToBootstrapPeers(server src.AgentGRPCServer) error {
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

func (n *BaseNode) ConnectToPeer(server src.AgentGRPCServer, address string, timeout time.Duration) (*grpc.ClientConn, error) {
	connection, err := server.ClientConnection(address)
	if err != nil {
		utils.Logger.Warn().Err(err).Str("Address", address).Msg("Failed to connect to peers")
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	utils.Logger.Debug().Str("nodeType", n.Type).Msg("Connecting to bootstrap peer")
	healthClient := grpc_health_v1.NewHealthClient(connection)
	res, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	if res.Status == grpc_health_v1.HealthCheckResponse_SERVING {
		utils.Logger.Debug().Str("healthStatus", res.Status.String()).Msgf("Bootstrap peer node is healthy.")
	} else {
		utils.Logger.Debug().Str("healthStatus", res.Status.String()).Msgf("Bootstrap peer node is not healthy.")
		return nil, fmt.Errorf("bootstrap peer node is not healthy")
	}
	return connection, nil
}
