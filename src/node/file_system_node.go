package node

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/utils"
	"google.golang.org/grpc"
)

type FileSystemNode struct {
	*BaseNode
	pb.UnimplementedFsNodeServiceServer
}

func NewFileSystemNode(nodeOptions *NodeOptions) *FileSystemNode {
	return &FileSystemNode{BaseNode: NewBaseNode(nodeOptions)}
}

func (n *FileSystemNode) Stop() error {
	return n.Terminate()
}

func (n *FileSystemNode) Register(server *grpc.Server) {
	pb.RegisterFsNodeServiceServer(server, &pb.UnimplementedFsNodeServiceServer{})
	utils.Logger.Info().Str("nodeType", n.NodeType.String()).Msg("Node registered")
}

func (n *FileSystemNode) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	hash := sha256.Sum256(req.Payload)
	fileId := hex.EncodeToString(hash[:])
	n.Storage[fileId] = req.Payload
	return &pb.UploadFileResponse{
		FileId: fileId,
	}, nil
}
func (n *FileSystemNode) DownloadFile(ctx context.Context, request *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (n *FileSystemNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	//TODO implement me
	panic("implement me")
}
