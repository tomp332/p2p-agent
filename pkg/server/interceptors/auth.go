package interceptors

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/server/managers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"slices"
)

type AuthInterceptor struct {
	jwtManager           *managers.JWTManager
	nodeAccessibleRoutes []string
}

func NewAuthInterceptor(jwtManager *managers.JWTManager, nodeAccessibleRoutes []string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, nodeAccessibleRoutes}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	if !slices.Contains(interceptor.nodeAccessibleRoutes, method) {
		// everyone can access
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	log.Debug().Str("nodeType", claims.NodeType.ToString()).Msgf("Granted access to node route.")
	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
