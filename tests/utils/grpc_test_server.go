package utils

import (
	"context"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

const (
	TestServerHostname = "bufnet"
)

type TestGRPCable struct {
	src.GRPCServer
	t        *testing.T
	listener *bufconn.Listener
}

func NewTestAgentServer(t *testing.T) *TestGRPCable {
	testListener := bufconn.Listen(1024 * 1024)
	s := &TestGRPCable{
		t: t,
		GRPCServer: src.GRPCServer{
			BaseServer: grpc.NewServer(),
			Address:    TestServerHostname,
			Listener:   testListener,
		},
	}
	s.listener = testListener
	return s
}

func (s *TestGRPCable) ClientConnection(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if address == TestServerHostname {
		// To test non-valid addresses, needed to patch this.
		opts = append(opts, grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.listener.Dial()
		}))
		resolver.SetDefaultScheme("passthrough")
	} else {
		opts = append(opts, grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return bufconn.Listen(1024 * 1024).Dial()
		}))
	}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		utils.Logger.Warn().
			Err(err).
			Str("address", address).
			Msg("Failed to create client connection for TestGrpcServer.")
		return nil, err
	}
	return conn, nil
}
