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
	"io"
)

type FilesNode struct {
	*BaseNode
	Storage *storage.LocalStorage
	pb.UnimplementedFilesNodeServiceServer
	config *configs.P2PFilesNodeConfig
}

func NewP2PFilesNode(baseNode *BaseNode, nodeOptions *configs.P2PFilesNodeConfig) *FilesNode {
	n := &FilesNode{BaseNode: baseNode, Storage: storage.NewLocalStorage(&nodeOptions.Storage)}
	n.Register()
	return n
}

// UploadFile handles client-streaming RPC for file upload
func (n *FilesNode) UploadFile(stream pb.FilesNodeService_UploadFileServer) error {
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

func (n *FilesNode) DownloadFile(req *pb.DownloadFileRequest, stream pb.FilesNodeService_DownloadFileServer) error {
	ctx := stream.Context()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		dataChan, err := n.Storage.Get(ctx, req.FileId)
		if err != nil {
			return err
		}
		for chunk := range dataChan {
			err := stream.Send(&pb.DownloadFileResponse{
				FileId: req.FileId,
				Exists: true,
				Chunk:  chunk,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *FilesNode) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	err := n.Storage.Delete(ctx, request.GetFileId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{FileId: request.GetFileId()}, nil
}

func (n *FilesNode) Register() {
	pb.RegisterFilesNodeServiceServer(n.Server.ServerObj(), n)
	utils.Logger.Info().Str("nodeType", n.NodeType.String()).Msg("Node registered")
}

func createUniqueFileID(buffer *bytes.Buffer) (string, error) {
	fileHash := md5.Sum(buffer.Bytes())
	return hex.EncodeToString(fileHash[:]) + "-" + utils.GenerateRandomID(), nil
}
