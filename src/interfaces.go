package src

import (
	"context"
	"github.com/tomp332/p2p-agent/src/utils/configs"
)

type P2PNoder interface {
	Register()
	Options() *configs.P2PNodeBaseConfig
	ConnectToBootstrapPeers() error
}

type P2PNodeClienter interface{}

type P2PNodeConnection interface{}

type Storage interface {
	Initialize() error
	Put(ctx context.Context, fileID string, dataChan <-chan []byte) error
	Delete(ctx context.Context, fileId string) error
	Get(ctx *context.Context, fileId string, dataChan chan<- []byte) error
}
