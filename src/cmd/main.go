package main

import (
	"github.com/tomp332/p2p-agent/src"
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
	err := configs.LoadConfig("/Users/tompaz/Documents/git/p2p-agent/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	utils.SetupLogger()
	grpcServer := src.NewP2pAgentServer()
	nodes, err := node.InitializeNodes(grpcServer, &configs.MainConfig.Nodes)
	for _, n := range nodes {
		err = n.ConnectToBootstrapPeers()
		if err != nil {
			utils.Logger.Warn().Err(err).Msg("failed to connect to bootstrap node")
		}
	}
	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	utils.Logger.Debug().Msgf("Received signal %s, exiting gracefully...", sig)

}
