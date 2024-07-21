package p2p

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/utils"
	"google.golang.org/grpc"
	"io"
)

type FileSystemNodeConfig struct {
	// File size in MB
	MaxFileSize float64 `json:"max_file_size"`
}

type FileSystemNode struct {
	*BaseNode
	Storage map[string][]byte
	pb.UnimplementedFileSystemNodeServiceServer
	config *FileSystemNodeConfig
}

func NewFileSystemNode(nodeOptions *P2PNodeConfig) *FileSystemNode {
	return &FileSystemNode{BaseNode: NewBaseNode(nodeOptions), Storage: make(map[string][]byte)}
}

func (n *FileSystemNode) ParseNodeConfig(specificConfig map[string]interface{}) error {
	config, err := utils.MapToStruct[FileSystemNodeConfig](specificConfig)
	if err != nil {
		return err
	}
	n.config = config
	return nil
}

func (n *FileSystemNode) Stop() error {
	return n.Terminate()
}

func (n *FileSystemNode) Register(server *grpc.Server) {
	pb.RegisterFileSystemNodeServiceServer(server, n)
	utils.Logger.Info().Str("nodeType", n.NodeType.String()).Msg("Node registered")
}

// UploadFile handles client-streaming RPC for file upload
func (n *FileSystemNode) UploadFile(stream pb.FileSystemNodeService_UploadFileServer) error {
	var buffer bytes.Buffer
	var fileID string
	totalSize := float64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// End of stream, store the accumulated file in memory
			fileHash := sha256.Sum256(buffer.Bytes())
			fileID = hex.EncodeToString(fileHash[:])
			n.Storage[fileID] = buffer.Bytes()

			return stream.SendAndClose(&pb.UploadFileResponse{
				FileId:  fileID,
				Success: true,
				Message: "File uploaded successfully",
			})
		}
		if err != nil {
			return err
		}

		// Check file size limit in megabytes
		totalSize += float64(len(req.ChunkData)) / (1024 * 1024)
		if totalSize > n.config.MaxFileSize {
			return stream.SendAndClose(&pb.UploadFileResponse{
				Success: false,
				Message: fmt.Sprintf("File size exceeds the maximum allowed size of %d bytes", n.config.MaxFileSize),
			})
		}

		// Write chunk to buffer
		if _, err := buffer.Write(req.ChunkData); err != nil {
			return fmt.Errorf("failed to write chunk data: %v", err)
		}
	}
}

func (n *FileSystemNode) DownloadFile(ctx context.Context, request *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return &pb.DownloadFileResponse{
		FileId: request.FileId,
		Chunk:  n.Storage[request.FileId],
	}, nil
}

func (n *FileSystemNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	//TODO implement me
	panic("implement me")
}
