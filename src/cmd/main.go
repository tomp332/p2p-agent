package main

import (
	"encoding/json"
	"github.com/mcuadros/go-defaults"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func LoadConfig(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&configs.MainConfig)
	defaults.SetDefaults(configs.MainConfig)
	return err
}

func main() {
	// Load configuration
	err := LoadConfig("/Users/tompaz/Documents/git/p2p-agent/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	utils.SetupLogger()
	grpcServer := src.NewP2pAgentServer()
	_, err = node.InitializeNodes(grpcServer, &configs.MainConfig.Nodes)
	if err != nil {
		return
	}
	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	utils.Logger.Debug().Msgf("Received signal %s, exiting gracefully...", sig)

}
