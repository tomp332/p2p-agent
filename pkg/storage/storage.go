package storage

import (
	"context"
	"github.com/tomp332/p2p-agent/pkg/utils/types"
	"sync"
)

type Storage interface {
	Initialize() error
	Put(ctx context.Context, fileName string, fileDataChan <-chan *types.TransferChunkData, wg *sync.WaitGroup) error
	Delete(ctx context.Context, fileName string) error
	Get(ctx context.Context, fileName string) (<-chan *types.TransferChunkData, error)
	Exists(fileName string) bool
}
