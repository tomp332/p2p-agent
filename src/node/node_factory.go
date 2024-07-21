package node

import (
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/utils"
)

func InitializeP2PNodes() error {
	for _, conf := range utils.MainConfig.Nodes {
		var n p2p.P2PNode
		nodeOptions, err := utils.MapToStruct[p2p.P2PNodeConfig](conf)
		if err != nil {
			utils.Logger.Warn().Msgf("Error parsing node config: %v", err)
			continue
		}
		switch nodeOptions.Type {
		case "file_system":
			n = p2p.NewFileSystemNode(nodeOptions)
		default:
			utils.Logger.Error().Str("nodeType", nodeOptions.Type).Msgf("Unkown node type specified in configuration.")
			return nil
		}
		if err = n.ParseNodeConfig(nodeOptions.ExtraConfig); err != nil {
			utils.Logger.Warn().Msgf("Error parsing node config: %v", err)
			continue
		}
		utils.Logger.Info().Str("nodeId", n.GetID()).Str("type", n.GetType().String()).Msgf("Created new node")
		n.Register(MainAgentServer.BaseServer)
	}
	return nil
}
