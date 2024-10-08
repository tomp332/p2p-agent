package nodes

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/server"
	"github.com/tomp332/p2p-agent/pkg/server/managers"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

type P2PNoder interface {
	ServiceRegister(server *server.GRPCServer)
	InterceptorsRegister(server *server.GRPCServer)
	Options() *configs.NodeConfig
	ConnectToBootstrapPeers()
	Terminate() error
}

type PeerConnection struct {
	ID             string
	ConnectionInfo *configs.BootStrapNodeConnection
	NodeClient     NodeClient
	Token          string
}

type BaseNode struct {
	configs.NodeConfig
	ConnectedPeers     map[string]PeerConnection
	ProtectedRoutes    []string
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
	AuthManager        managers.AuthenticationManager
	MainContext        context.Context
	MainContextCancel  context.CancelFunc
}

func NewBaseNode(options *configs.NodeConfig, mainCtx context.Context) *BaseNode {
	ctx, cancel := context.WithCancel(mainCtx)
	n := &BaseNode{
		NodeConfig:        *options,
		MainContext:       ctx,
		MainContextCancel: cancel,
		ConnectedPeers:    make(map[string]PeerConnection),
	}
	return n
}

func (n *BaseNode) Register(_ *server.GRPCServer) {}

func (n *BaseNode) InterceptorsRegister(_ *server.GRPCServer) {}

func (n *BaseNode) Options() *configs.NodeConfig {
	return &n.NodeConfig
}

func (n *BaseNode) Terminate() error {
	n.MainContextCancel()
	return nil
}

func (n *BaseNode) ConnectToBootstrapPeers() error {
	return nil
}

func (n *BaseNode) ConnectToPeer(connectionInfo *configs.BootStrapNodeConnection, timeout time.Duration) (*grpc.ClientConn, error) {
	connection, err := clientConnection(connectionInfo)
	if err != nil {
		log.Warn().Err(err).Str("Address", fmt.Sprintf("%s:%v", connectionInfo.Host, connectionInfo.Port)).Msg("Failed to connect to peers")
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

func clientConnection(connectionInfo *configs.BootStrapNodeConnection) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%v", connectionInfo.Host, connectionInfo.Port), opts...)
	if err != nil {
		log.Warn().Err(err).Str("address", connectionInfo.Host).Int64("port", connectionInfo.Port).Msg("Failed to connect to peers")
		return nil, err
	}
	return conn, nil
}
