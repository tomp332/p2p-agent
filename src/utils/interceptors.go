package utils

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

// UnaryServerInterceptor logs unary RPCs.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		h, err := handler(ctx, req)
		duration := time.Since(start)

		st, _ := status.FromError(err)
		Logger.Info().Msgf("RPC: %s, Duration: %d, Status: %s", info.FullMethod, duration, st.Code())
		if err != nil {
			Logger.Error().Err(err)
		}
		return h, err
	}
}

// StreamServerInterceptor logs stream RPCs.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		st, _ := status.FromError(err)
		Logger.Info().Msgf("RPC: %s, Duration: %d, Status: %s", info.FullMethod, duration, st.Code())
		if err != nil {
			Logger.Error().Err(err)
		}
		return err
	}
}
