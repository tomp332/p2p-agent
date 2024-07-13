package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	pb "github.com/tomp332/p2fs/src/protocol"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedFileSharingServer
	storage map[string][]byte
}

func (s *server) StoreChunk(ctx context.Context, req *pb.StoreChunkRequest) (*pb.StoreChunkResponse, error) {
	hash := hashChunk(req.GetChunk())
	s.storage[hash] = req.GetChunk()
	return &pb.StoreChunkResponse{Hash: hash}, nil
}

func (s *server) RetrieveChunk(ctx context.Context, req *pb.RetrieveChunkRequest) (*pb.RetrieveChunkResponse, error) {
	chunk, exists := s.storage[req.GetHash()]
	return &pb.RetrieveChunkResponse{Chunk: chunk, Exists: exists}, nil
}

func hashChunk(chunk []byte) string {
	hash := sha256.Sum256(chunk)
	return hex.EncodeToString(hash[:])
}

func startGRPCServer(port string, storage map[string][]byte) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileSharingServer(s, &server{storage: storage})
	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
