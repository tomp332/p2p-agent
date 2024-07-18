package node

import (
	"fmt"
	"github.com/tomp332/p2p-agent/src/utils"
	"google.golang.org/grpc"
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
	address := fmt.Sprintf("%s:%d", utils.MainConfig.ServerConfig.Host, utils.MainConfig.ServerConfig.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		utils.Logger.Info().Msgf("gRPC server listening on %s", address)
		if err := MainAgentServer.BaseServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}

func (s *GrpcAgentServer) Terminate(nodes []P2PNode) error {
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
