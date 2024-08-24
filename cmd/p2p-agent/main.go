package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tomp332/p2p-agent/pkg/nodes/factory"
	"github.com/tomp332/p2p-agent/pkg/server"
	"github.com/tomp332/p2p-agent/pkg/utils"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"os"
	"os/signal"
	"syscall"
)

var cfgFile string

func init() {
	// Define a persistent flag for the config file path
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file path (default is $CWD/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&configs.MainConfig.LogLevel, "log-level", "l", "info", "Set log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVarP(&configs.MainConfig.LoggerMode, "logger-mode", "m", "console", "Logger mode (dev, production)")
	rootCmd.PersistentFlags().StringVarP(&configs.MainConfig.ServerConfig.Host, "host", "", "localhost", "Hostname to host the server on")
	rootCmd.PersistentFlags().Int32VarP(&configs.MainConfig.ServerConfig.Port, "port", "", 8080, "Port to listen on")
	rootCmd.PersistentFlags().StringVarP(&configs.MainConfig.ID, "node-id", "i", utils.GenerateRandomID(), "Node ID")
}

// rootCmd is the main command for the CLI
var rootCmd = &cobra.Command{
	Use:   "p2p-agent",
	Short: "P2P Agent Node",
	Long:  `P2P Agent Node - A CLI to configure and start the P2P agent nodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.LoadConfig(cfgFile)
		utils.SetupLogger()
		grpcServer := server.NewP2pAgentServer()
		_, err := factory.InitializeNodes(grpcServer, &configs.MainConfig.Nodes)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize nodes")
			return
		}
		err = grpcServer.Setup()
		err = grpcServer.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to star grpc server")
		}
		// Signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		log.Debug().Msgf("Received signal %s, exiting gracefully...", sig)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
