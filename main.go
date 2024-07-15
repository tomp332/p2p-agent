package main

import (
	"github.com/tomp332/p2fs/src"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	nodeOptions := &src.NodeOptions{
		Storage:              map[string][]byte{},
		BootstrapPeerAddrs:   src.MainConfig.BootstrapPeerAddrs,
		ServerPort:           src.MainConfig.ServerPort,
		BootstrapNodeTimeout: src.MainConfig.BootstrapNodeTimeout,
	}
	node := src.NewNode(nodeOptions)
	if err := node.StartServer(); err != nil {
		src.Logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	src.Logger.Debug().Msg("Received shutdown signal, exiting gracefully...")
	if err := node.Terminate(); err != nil {
		src.Logger.Fatal().Err(err)
	}
}
