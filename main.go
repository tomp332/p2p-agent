package main

import (
	"github.com/tomp332/p2fs/src"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var bootstrapAddresses = []string{
	"192.168.6.41:5001",
	"192.168.6.49:5002",
}

func main() {
	nodeOptions := &src.NodeOptions{
		Storage:              map[string][]byte{},
		BootstrapPeerAddrs:   bootstrapAddresses,
		ServerPort:           5001,
		BootstrapNodeTimeout: time.Second * 6,
	}
	node := src.NewNode(nodeOptions)
	node.StartServer()
	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Received shutdown signal, exiting gracefully...")
	node.StopServer()
}
