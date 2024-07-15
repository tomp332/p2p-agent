package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomp332/p2p-agent/src/node"
	"github.com/tomp332/p2p-agent/src/util"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	conf, err := util.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(util.UnaryServerInterceptor()),
		grpc.StreamInterceptor(util.StreamServerInterceptor()),
	)

	// Initialize nodes and register services
	nodes, err := node.InitializeNodes(conf.Nodes, grpcServer)
	if err != nil {
		log.Fatalf("failed to initialize nodes: %v", err)
	}

	// Start gRPC server
	address := conf.ServerConfig.Host + ":" + conf.ServerConfig.Port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Printf("gRPC server listening on %s", address)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Received signal %s, exiting gracefully...", sig)

	// Stop all nodes
	for _, n := range nodes {
		if err := n.Stop(); err != nil {
			log.Printf("failed to stop node: %v", err)
		}
	}

	grpcServer.GracefulStop()
	log.Println("gRPC server stopped.")
}
