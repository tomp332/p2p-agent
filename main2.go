package main

import (
	"github.com/tomp332/p2fs/src/util"
	"github.com/tomp332/p2fs/temp"
	"github.com/tomp332/p2fs/temp/fsP2p"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := util.LoadConfig(&util.ConfigSource{FilePath: "config.json"}); err != nil {
		log.Fatal(err)
	}
	util.SetupLogger()
}

func main() {
	nodeOptions := &temp.NodeOptions{
		Storage:            map[string][]byte{},
		BootstrapPeerAddrs: util.MainConfig.BootstrapPeerAddrs,
		NodeServerOptions: temp.NodeServerOptions{
			ServerPort: util.MainConfig.ServerPort,
			ServerHost: util.MainConfig.ServerHost,
		},
		BootstrapNodeTimeout: util.MainConfig.BootstrapNodeTimeout,
	}
	node := fsP2p.NewFsNode(nodeOptions)
	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	node.Terminate()
}
