package interfaces

import (
	"github.com/tomp332/p2p-agent/pkg/pb"
	"google.golang.org/grpc"
)

// UploadFileServerStream defines the server-side stream interface
type UploadFileServerStream interface {
	SendAndClose(*pb.UploadFileResponse) error
	Recv() (*pb.UploadFileRequest, error)
	grpc.ServerStream // Embeds ServerStream for other methods
}

type DownloadFileClientStream interface {
	Send(*pb.DownloadFileResponse) error
	Recv() (*pb.DownloadFileResponse, error)
	grpc.ClientStream // Embeds ServerStream for other methods
	grpc.ServerStream // Embeds ServerStream for other methods
}

type DownloadFileServerStream interface {
	Send(*pb.DownloadFileResponse) error
	SendAndClose(response *pb.DownloadFileResponse) error
	Recv() (*pb.DownloadFileRequest, error)
	grpc.ServerStream // Embeds ServerStream for other methods
}
