package types

type TransferChunkData struct {
	ChunkSize int
	ChunkData []byte
	Metadata  interface{}
}
