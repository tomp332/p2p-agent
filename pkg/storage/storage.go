package storage

import (
	"context"
	"github.com/tomp332/p2p-agent/pkg/utils/types"
)

type Storage interface {
	Initialize() error
	Put(ctx context.Context, fileID string, dataChan <-chan types.TransferChunkData) (int64, error)
	Delete(ctx context.Context, fileId string) error
	Get(ctx context.Context, fileID string) (<-chan types.TransferChunkData, error)
}
