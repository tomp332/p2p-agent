package utils

import (
	"context"
	"errors"
	"github.com/tomp332/p2p-agent/src/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

const BUFSIZE = 1024 * 1024

func Dialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func StartTestGRPCServer(t *testing.T, filesNode pb.FilesNodeServiceServer) (*grpc.ClientConn, func()) {
	lis := bufconn.Listen(BUFSIZE)
	s := grpc.NewServer()
	pb.RegisterFilesNodeServiceServer(s, filesNode)
	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			t.Fatalf("Server exited with error: %v", err)
		}
	}()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithContextDialer(Dialer(lis)))
	conn, err := grpc.NewClient("passthrough://bufnet", opts...)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	cleanup := func() {
		s.Stop()
		lis.Close()
		conn.Close()
	}

	return conn, cleanup
}
