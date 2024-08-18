package types

type TransferChunkData struct {
	ChunkSize int64
	ChunkData []byte
	TotalSize int64
}
