package src

import (
	pb "github.com/tomp332/p2fs/src/protocol"
	"google.golang.org/grpc"
)

type Options struct {
	Port int32
}

type FileSharingServer struct {
	pb.FileSharingServer
	Server *grpc.Server
}
