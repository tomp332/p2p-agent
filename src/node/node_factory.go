package node

import (
	"errors"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/node/p2p/file_node"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
)

func InitializeNodes(server src.AgentGRPCServer, configs *[]configs.NodeConfigs) ([]src.P2PNoder, error) {
	var initializedNodes []src.P2PNoder
	for _, nodeConfig := range *configs {
		if initNode, err := InitializeNode(server, &nodeConfig); err != nil {
			utils.Logger.Error().Err(err).Str("nodeType", nodeConfig.BaseNodeConfigs.Type).Msgf("Failed to initialize node")
			continue
		} else {
			initializedNodes = append(initializedNodes, initNode)
		}
	}
	if initializedNodes == nil {
		return nil, errors.New("failed to initialize nodes")
	}
	err := server.Start()
	if err != nil {
		return nil, err
	}
	return initializedNodes, nil
}

func InitializeNode(server src.AgentGRPCServer, config *configs.NodeConfigs) (src.P2PNoder, error) {
	var n src.P2PNoder
	baseNode := p2p.NewBaseNode(server, &config.BaseNodeConfigs)
	switch config.BaseNodeConfigs.Type {
	case "base":
		n = baseNode
	case "file":
		n = file_node.NewP2PFilesNode(baseNode, &config.FilesNodeConfigs)
	default:
		return nil, errors.New("invalid node type")
	}
	n.Register()
	err := n.ConnectToBootstrapPeers()
	if err != nil {
		utils.Logger.Warn().Msgf("Bootstrapping process failed")
	}
	return n, nil
}
