package src

import (
	"fmt"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

type AgentGRPCServer interface {
	Start() error
	Terminate() error
	ClientConnection(address string) (*grpc.ClientConn, error)
	ServerObj() *grpc.Server
}

type GRPCServer struct {
	BaseServer *grpc.Server
	Listener   net.Listener
	Address    string
	wg         sync.WaitGroup
}

func NewP2pAgentServer() *GRPCServer {
	g := &GRPCServer{
		wg: sync.WaitGroup{},
		BaseServer: grpc.NewServer(
			grpc.UnaryInterceptor(utils.UnaryServerInterceptor()),
			grpc.StreamInterceptor(utils.StreamServerInterceptor()),
		),
	}
	address := fmt.Sprintf("%s:%d", configs.MainConfig.ServerConfig.Host, configs.MainConfig.ServerConfig.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		utils.Logger.Fatal().Err(err).
			Str("address", configs.MainConfig.ServerConfig.Host).
			Int32("port", configs.MainConfig.ServerConfig.Port).
			Msg("Agent server failed to start listening on configured address and port.")
	}
	g.Listener = lis
	g.Address = address
	return g
}

func (s *GRPCServer) Terminate() error {
	s.BaseServer.GracefulStop()
	return nil
}

func (s *GRPCServer) ClientConnection(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		utils.Logger.Warn().Err(err).Str("address", address).Msg("Failed to connect to peers")
		return nil, err
	}
	return conn, nil
}

func (s *GRPCServer) ServerObj() *grpc.Server {
	return s.BaseServer
}

func (s *GRPCServer) Start() error {
	// Register health server
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.BaseServer, healthServer)
	reflection.Register(s.BaseServer)
	go func() {
		utils.Logger.Info().Msgf("gRPC server listening on %s", s.Address)
		if err := s.BaseServer.Serve(s.Listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}
