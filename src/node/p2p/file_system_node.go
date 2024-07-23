package p2p

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/storage"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"google.golang.org/grpc"
	"io"
)

type FileSystemNode struct {
	*BaseNode
	Storage *storage.LocalStorage
	pb.UnimplementedFileSystemNodeServiceServer
	config *configs.FileSystemNodeConfig
}

func NewFileSystemNode(storageOptions *configs.LocalStorageConfig, nodeOptions *configs.P2PNodeConfig) *FileSystemNode {
	return &FileSystemNode{BaseNode: NewBaseNode(nodeOptions), Storage: storage.NewLocalStorage(storageOptions)}
}

func (n *FileSystemNode) ParseNodeConfig(specificConfig map[string]interface{}) error {
	config, err := utils.MapToStruct[configs.FileSystemNodeConfig](specificConfig)
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
	dataChan := make(chan []byte, 1024) // Buffered channel for file data

	ctx := stream.Context()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fileID, err := createUniqueFileID(&buffer)
			close(dataChan)

			fileSize, err := n.Storage.Put(ctx, fileID, dataChan)
			if err != nil {
				return stream.SendAndClose(&pb.UploadFileResponse{
					FileId:  fileID,
					Success: false,
					Message: fmt.Sprintf("Failed to store file: %v", err),
				})
			}

			return stream.SendAndClose(&pb.UploadFileResponse{
				FileId:   fileID,
				FileSize: fileSize,
				Success:  true,
				Message:  "File uploaded successfully",
			})
		}
		if err != nil {
			return err
		}
		chunk := req.GetChunkData()
		buffer.Write(chunk)
		dataChan <- chunk
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
		Chunk:  make([]byte, 1),
	}, nil
}

func (n *FileSystemNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func createUniqueFileID(buffer *bytes.Buffer) (string, error) {
	fileHash := md5.Sum(buffer.Bytes())
	return hex.EncodeToString(fileHash[:]) + "-" + utils.GenerateRandomID(), nil
}
