package node

import (
	"github.com/tomp332/p2p-agent/src/utils"
)

func InitializeNodes() error {
	for _, conf := range utils.MainConfig.Nodes {
		nodeOptions := &NodeOptions{
			Storage:              make(map[string][]byte),
			BootstrapPeerAddrs:   conf.BootstrapPeerAddrs,
			BootstrapNodeTimeout: conf.BootstrapNodeTimeout,
		}

		var n P2PNode
		switch conf.Type {
		case "file_system":
			n = NewFileSystemNode(nodeOptions)
		default:
			utils.Logger.Error().Str("nodeType", conf.Type).Msgf("Unkown node type specified in configuration.")
			return nil
		}
		utils.Logger.Info().Str("nodeId", n.GetID()).Str("type", n.GetType().String()).Msgf("Created new node")
		n.Register(MainAgentServer.BaseServer)
	}
	return nil
}
