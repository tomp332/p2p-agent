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
	address        string
	grpcConnection *grpc.ClientConn
	nodeClients    *[]interface{}
}

type BaseNode struct {
	configs.P2PNodeBaseConfig
	ID             string
	NodeType       src.NodeType
	Context        context.Context
	ConnectedPeers []NodeConnection
	Server         src.AgentGRPCServer
	wg             sync.WaitGroup
	cancelFunc     context.CancelFunc
}

func NewBaseNode(agentServer src.AgentGRPCServer, options *configs.P2PNodeBaseConfig) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	var id string
	if id = options.ID; id == "" {
		id = utils.GenerateRandomID()
	}
	n := &BaseNode{
		ID:                id,
		Context:           ctx,
		P2PNodeBaseConfig: *options,
		ConnectedPeers:    make([]NodeConnection, 0),
		Server:            agentServer,
		wg:                sync.WaitGroup{},
		cancelFunc:        cancel,
	}
	return n
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

func (n *BaseNode) ConnectToBootstrapPeers() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, address := range n.BootstrapPeerAddrs {
		connection, err := n.Server.ClientConnection(address)
		if err != nil {
			utils.Logger.Warn().Err(err).Str("address", address).Msg("Failed to connect to peers")
			continue
		}
		utils.Logger.Debug().Str("nodeType", n.NodeType.String()).Msg("Connecting to bootstrap peer")
		healthClient := grpc_health_v1.NewHealthClient(connection)
		res, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		if res.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			utils.Logger.Debug().Str("healthStatus", res.Status.String()).Msgf("Bootstrap peer node is healthy.")
		} else {
			utils.Logger.Debug().Str("healthStatus", res.Status.String()).Msgf("Bootstrap peer node is not healthy.")
			continue
		}
		n.ConnectedPeers = append(n.ConnectedPeers, NodeConnection{
			address:        address,
			grpcConnection: connection,
		})
		utils.Logger.Debug().Str("address", address).Msg("Successfully connected to bootstrap peer node")
	}
	return nil
}
