package node

import (
	"errors"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/node/p2p/p2p_file_node"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
)

func InitializeNodes(server src.AgentGRPCServer, configs *[]configs.NodeConfigs) ([]src.P2PNoder, error) {
	var initializedNodes []src.P2PNoder
	for _, nodeConfig := range *configs {
		if checkTypeAlreadyInitialized(&initializedNodes, &nodeConfig) {
			utils.Logger.Error().Str("nodeType", nodeConfig.Type).Msg("Duplicate node type detected.")
			return nil, errors.New("duplicate node type detected")
		}
		if initNode, err := InitializeNode(server, &nodeConfig); err != nil {
			utils.Logger.Error().Err(err).Str("nodeType", nodeConfig.Type).Msgf("Failed to initialize node")
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
	baseNode := p2p.NewBaseNode(config)
	switch config.Type {
	case "base":
		n = baseNode
	case "file":
		n = p2p_file_node.NewP2PFilesNode(baseNode, &config.P2PFilesNodeConfig)
	default:
		return nil, errors.New("invalid node type")
	}
	n.Register(server.ServerObj())
	err := n.ConnectToBootstrapPeers(server)
	if err != nil {
		utils.Logger.Warn().Msgf("Bootstrapping process failed")
	}
	return n, nil
}

func checkTypeAlreadyInitialized(initializedNodes *[]src.P2PNoder, config *configs.NodeConfigs) bool {
	for _, node := range *initializedNodes {
		if node.Options().Type == config.Type {
			return true
		}
	}
	return false
}
