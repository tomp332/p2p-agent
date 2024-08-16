package storage

import "context"

type Storage interface {
	Initialize() error
	Put(ctx context.Context, fileID string, dataChan <-chan []byte) (float64, error)
	Delete(ctx context.Context, fileId string) error
	Get(ctx context.Context, fileID string) (<-chan []byte, error)
	Search(fileID string) bool
}
