package src

import (
	"context"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"google.golang.org/grpc"
)

type P2PNoder interface {
	Register(server *grpc.Server)
	Options() *configs.NodeConfig
	ConnectToBootstrapPeers(server AgentGRPCServer) error
}

type P2PNodeClienter interface{}

type P2PNodeConnection interface{}

type Storage interface {
	Initialize() error
	Put(ctx context.Context, fileID string, dataChan <-chan []byte) (float64, error)
	Delete(ctx context.Context, fileId string) error
	Get(ctx context.Context, fileID string) (<-chan []byte, error)
	Search(fileID string) bool
}

type AgentGRPCServer interface {
	Start() error
	Terminate() error
	ClientConnection(address string) (*grpc.ClientConn, error)
	ServerObj() *grpc.Server
}
