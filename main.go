package main

import (
	"github.com/tomp332/p2fs/src/common"
	"github.com/tomp332/p2fs/src/fsNode"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	nodeOptions := &fsNode.NodeOptions{
		Storage:            map[string][]byte{},
		BootstrapPeerAddrs: common.MainConfig.BootstrapPeerAddrs,
		ServerOptions: fsNode.ServerOptions{
			ServerPort: common.MainConfig.ServerPort,
			ServerHost: common.MainConfig.ServerHost,
		},
		BootstrapNodeTimeout: common.MainConfig.BootstrapNodeTimeout,
	}
	node := fsNode.NewNode(nodeOptions)
	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	node.Terminate()
}
