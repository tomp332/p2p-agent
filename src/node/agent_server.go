package node

import (
	"fmt"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

var (
	MainAgentServer *GrpcAgentServer
)

type P2pAgentServer interface {
	Start() error
	Terminate() error
}

type GrpcAgentServer struct {
	wg         sync.WaitGroup
	BaseServer *grpc.Server
}

func init() {
	// Create gRPC server
	MainAgentServer = NewP2pAgentServer()
	MainAgentServer.BaseServer = grpc.NewServer(
		grpc.UnaryInterceptor(utils.UnaryServerInterceptor()),
		grpc.StreamInterceptor(utils.StreamServerInterceptor()),
	)
}

func NewP2pAgentServer() *GrpcAgentServer {
	return &GrpcAgentServer{
		wg: sync.WaitGroup{},
	}
}

func (s *GrpcAgentServer) Start() error {
	// Start gRPC server
	address := fmt.Sprintf("%s:%d", configs.MainConfig.ServerConfig.Host, configs.MainConfig.ServerConfig.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		utils.Logger.Fatal().Err(err).
			Str("address", configs.MainConfig.ServerConfig.Host).
			Int32("port", configs.MainConfig.ServerConfig.Port).
			Msg("Agent server failed to start listening on configured address and port.")
	}
	reflection.Register(s.BaseServer)
	go func() {
		utils.Logger.Info().Msgf("gRPC server listening on %s", address)
		if err := MainAgentServer.BaseServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}

func (s *GrpcAgentServer) Terminate(nodes []p2p.P2PNode) error {
	// Stop all nodes
	for _, n := range nodes {
		if err := n.Stop(); err != nil {
			utils.Logger.Warn().Msgf("failed to stop node: %v", err)
		}
	}
	MainAgentServer.BaseServer.GracefulStop()
	utils.Logger.Info().Msg("gRPC server stopped.")
	return nil
}
