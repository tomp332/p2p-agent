package main

import (
	"github.com/tomp332/p2p-agent/src/node"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	err := configs.LoadConfig("/Users/tompaz/Documents/git/p2p_test-agent/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	utils.SetupLogger()

	// Initialize nodes and register services
	err = node.InitializeP2PNodes()
	if err != nil {
		log.Fatalf("failed to initialize nodes: %v", err)
	}

	err = node.MainAgentServer.Start()
	if err != nil {
		utils.Logger.Error().Msgf("failed to start agent server: %v", err)
		return
	}
	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	utils.Logger.Debug().Msgf("Received signal %s, exiting gracefully...", sig)

}
