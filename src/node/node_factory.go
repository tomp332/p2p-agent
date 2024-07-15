package node

import (
	"fmt"

	"github.com/tomp332/p2p-agent/src/util"
	"google.golang.org/grpc"
)

func InitializeNodes(configs []util.NodeConfig, server *grpc.Server) ([]Node, error) {
	var nodes []Node
	for _, conf := range configs {
		nodeOptions := &NodeOptions{
			Storage:              conf.Storage,
			BootstrapPeerAddrs:   conf.BootstrapPeerAddrs,
			BootstrapNodeTimeout: conf.BootstrapNodeTimeout,
		}

		var n Node
		switch conf.Type {
		case "file_system":
			n = NewFileSystemNode(nodeOptions)
		default:
			return nil, fmt.Errorf("unknown node type: %s", conf.Type)
		}

		n.Register(server)
		nodes = append(nodes, n)
	}

	return nodes, nil
}
