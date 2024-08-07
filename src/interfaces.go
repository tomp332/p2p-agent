package src

import (
	"context"
)

type P2PNoder interface {
	Register()
	GetID() string
	GetType() NodeType
	ConnectToBootstrapPeers() error
}

type Storage interface {
	Initialize() error
	Put(ctx context.Context, fileID string, dataChan <-chan []byte) error
	Delete(ctx context.Context, fileId string) error
	Get(ctx *context.Context, fileId string, dataChan chan<- []byte) error
}
