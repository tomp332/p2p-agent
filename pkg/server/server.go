package server

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/utils"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

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
		log.Fatal().Err(err).
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

func (s *GRPCServer) ServerObj() *grpc.Server {
	return s.BaseServer
}

func (s *GRPCServer) Start() error {
	// Register health server
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.BaseServer, healthServer)
	reflection.Register(s.BaseServer)
	go func() {
		log.Info().Msgf("gRPC server listening on %s", s.Address)
		if err := s.BaseServer.Serve(s.Listener); err != nil {
			log.Fatal().Msgf("failed to serve: %v", err)
		}
	}()
	return nil
}
