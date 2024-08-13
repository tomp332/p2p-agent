package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"os"
	"os/signal"
	"syscall"
)

var cfgFile string

func init() {
	// Define a persistent flag for the config file path
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is $CWD/config.yaml)")
	// Define other flags and bind them to viper
	defineAndBindFlags(rootCmd)
}

func defineAndBindFlags(cmd *cobra.Command) {
	cmd.Flags().String("host", "localhost", "Host address of the server")
	cmd.Flags().Int("port", 8080, "Port on which the server runs")
	cmd.Flags().String("storage-path", "./data", "Path to store node data")
	cmd.Flags().String("logger-mode", "production", "Logger mode (e.g., development, production)")
	cmd.Flags().String("log-level", "info", "Log level (e.g., debug, info, warn, error)")
	cmd.Flags().StringSlice("nodes", []string{}, "List of node configurations (JSON encoded)")

	// Bind the flags directly to viper
	viper.BindPFlag("server.host", cmd.Flags().Lookup("host"))
	viper.BindPFlag("server.port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("storage.storage_path", cmd.Flags().Lookup("storage-path"))
	viper.BindPFlag("logger_mode", cmd.Flags().Lookup("logger-mode"))
	viper.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))
	viper.BindPFlag("nodes", cmd.Flags().Lookup("nodes"))
}

// rootCmd is the main command for the CLI
var rootCmd = &cobra.Command{
	Use:   "p2p-agent",
	Short: "P2P Agent Node",
	Long:  `P2P Agent Node - A CLI to configure and start the P2P agent node.`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.LoadConfig(cfgFile)
		utils.SetupLogger()
		grpcServer := src.NewP2pAgentServer()
		_, err := node.InitializeNodes(grpcServer, &configs.MainConfig.Nodes)
		if err != nil {
			return
		}
		// Signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		utils.Logger.Debug().Msgf("Received signal %s, exiting gracefully...", sig)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
