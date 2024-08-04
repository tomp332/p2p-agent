package node

import (
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
)

func InitializeP2PNodes() error {
	// Iterate over `Nodes` section in the main configuration and handle each logic.
	for _, conf := range configs.MainConfig.Nodes {
		var n p2p.P2PNode
		baseNode := p2p.NewBaseNode(&conf.BaseNodeConfigs)
		switch conf.BaseNodeConfigs.Type {
		case "files":
			n = p2p.NewP2PFilesNode(baseNode, &conf.FilesNodeConfigs)
		default:
			utils.Logger.Error().Str("nodeType", conf.BaseNodeConfigs.Type).Msgf("Unkown node type specified in configuration.")
			return nil
		}
		baseNode.ConnectToBootstrapNodes()
		n.ConnectToPeers()
		utils.Logger.Info().Str("nodeId", n.GetID()).Str("type", n.GetType().String()).Msgf("Created new node")
		n.Register(MainAgentServer.BaseServer)
	}
	return nil
}
