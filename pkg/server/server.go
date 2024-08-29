package server

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/server/interceptors"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

type ServiceRegistrar func(server *grpc.Server)

type AgentServer interface {
	Terminate(context.Context) error
	Setup() error
	Start() error
}

type GRPCServer struct {
	BaseServer         *grpc.Server
	Listener           net.Listener
	Address            string
	ServiceRegistrars  map[configs.NodeType]ServiceRegistrar
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func NewP2pAgentServer() *GRPCServer {
	s := &GRPCServer{
		ServiceRegistrars:  make(map[configs.NodeType]ServiceRegistrar),
		unaryInterceptors:  make([]grpc.UnaryServerInterceptor, 0),
		streamInterceptors: make([]grpc.StreamServerInterceptor, 0),
	}
	address := fmt.Sprintf("%s:%d", configs.MainConfig.ServerConfig.Host, configs.MainConfig.ServerConfig.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).
			Str("address", configs.MainConfig.ServerConfig.Host).
			Int32("port", configs.MainConfig.ServerConfig.Port).
			Msg("Agent server failed to start listening on configured address and port.")
	}
	s.Listener = lis
	s.Address = address
	return s
}

func (s *GRPCServer) Terminate() error {
	s.BaseServer.GracefulStop()
	return nil
}

func (s *GRPCServer) Setup() error {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors.UnaryServerInterceptor())
	s.streamInterceptors = append(s.streamInterceptors, interceptors.StreamServerInterceptor())
	s.BaseServer = grpc.NewServer(
		grpc.UnaryInterceptor(chainUnaryInterceptors(s.unaryInterceptors...)),
		grpc.StreamInterceptor(chainStreamInterceptors(s.streamInterceptors...)),
	)
	return nil
}

func (s *GRPCServer) Start() error {
	// ServiceRegister all services
	for nodeType, registrar := range s.ServiceRegistrars {
		registrar(s.BaseServer)
		log.Debug().Str("service", nodeType.ToString()).Msg("Registered service")
	}

	if s.BaseServer == nil {
		log.Fatal().Msg("server setup() has not been called yet")
	}
	// ServiceRegister Health service.
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

func (s *GRPCServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) error {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
	return nil
}

func (s *GRPCServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) error {
	s.streamInterceptors = append(s.streamInterceptors, interceptors...)
	return nil
}

func chainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Wrap the handler with all interceptors
		chainedHandler := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			currentInterceptor := interceptors[i]
			nextHandler := chainedHandler // store the current handler to pass to the next interceptor
			chainedHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
				return currentInterceptor(ctx, req, info, nextHandler)
			}
		}
		// Call the final handler
		return chainedHandler(ctx, req)
	}
}

func chainStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Define a handler function that calls each interceptor in the chain
		currentHandler := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := currentHandler
			currentHandler = func(currentSrv interface{}, currentStream grpc.ServerStream) error {
				return interceptor(currentSrv, currentStream, info, next)
			}
		}
		// Call the final handler
		return currentHandler(srv, ss)
	}
}
