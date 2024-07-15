package fsNode

import (
	"context"
	"fmt"
	"github.com/tomp332/p2fs/src/common"
	pb "github.com/tomp332/p2fs/src/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
	"time"
)

type ServerOptions struct {
	ServerPort int32  `json:"server_port"`
	ServerHost string `json:"server_host"`
}

// FileSharingServer implements the FileSharing gRPC service.
type FileSharingServer struct {
	pb.UnimplementedFileSharingServer
	grpcServer *grpc.Server
	wg         sync.WaitGroup
	node       *Node
}

func NewFileSharingServer(n *Node) *FileSharingServer {
	return &FileSharingServer{
		node: n,
	}
}

func (s *FileSharingServer) StartServer() error {
	errorChan := make(chan error, 1)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.node.NodeOptions.ServerPort))
		if err != nil {
			errorChan <- err
			return
		}
		s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(common.UnaryServerInterceptor()),
			grpc.StreamInterceptor(common.StreamServerInterceptor()))
		pb.RegisterFileSharingServer(s.grpcServer, s)
		reflection.Register(s.grpcServer)
		common.Logger.Info().Msgf("gRPC server listening on port %d", s.node.NodeOptions.ServerPort)
		if err := s.grpcServer.Serve(lis); err != nil {
			errorChan <- fmt.Errorf("failed to serve: %v", err)
		}
	}()

	select {
	case err := <-errorChan:
		return err
	case <-time.After(time.Second * 1): // Adjust timeout as needed
		return nil
	}
}

func (s *FileSharingServer) StopServer() {
	if s.grpcServer != nil {
		common.Logger.Debug().Msg("Stopping gRPC Node server.")
		s.grpcServer.Stop()
		s.wg.Wait()
		common.Logger.Info().Msg("gRPC server stopped successfully.")
	}
}

func (s *FileSharingServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	return s.node.UploadFile(ctx, req)
}

func (s *FileSharingServer) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	return s.node.DownloadFile(ctx, req)
}

func (s *FileSharingServer) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	return s.node.DeleteFile(ctx, req)
}

func (s *FileSharingServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthResponse, error) {
	return s.node.HealthCheck(ctx, req)
}
