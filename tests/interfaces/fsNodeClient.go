package interfaces

import (
	"github.com/tomp332/p2p-agent/pkg/pb"
	"google.golang.org/grpc"
)

// UploadFileServer defines the server-side stream interface
type UploadFileServer interface {
	SendAndClose(*pb.UploadFileResponse) error
	Recv() (*pb.UploadFileRequest, error)
	grpc.ServerStream // Embeds ServerStream for other methods
}

type DownloadFileServer interface {
	Send(*pb.DownloadFileResponse) error
	grpc.ServerStream // Embeds ServerStream for other methods
}

type DownloadFileClient interface {
	SendAndClose(response *pb.DownloadFileResponse) error
	Recv() (*pb.DownloadFileRequest, error)
	grpc.ServerStream // Embeds ServerStream for other methods
}
