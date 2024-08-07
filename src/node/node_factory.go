package node

import (
	"errors"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/node/p2p/file_node"
	"github.com/tomp332/p2p-agent/src/storages"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"reflect"
)

func InitializeNodes(server src.AgentGRPCServer, nodeConfigs *configs.AllNodesMap) ([]src.P2PNoder, error) {
	var initializedNodes []src.P2PNoder
	val := reflect.ValueOf(nodeConfigs)
	// Dereference the pointer to access the struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// Get the number of fields in the struct
	numFields := val.NumField()
	for i := 0; i < numFields; i++ {
		field := val.Type().Field(i)
		value := val.Field(i)
		nodeConfig := value.Interface().(configs.NodeConfig)
		nodeConfig.Type = configs.StrToNodeType(field.Name)
		if fileNode, err := InitializeNode(server, &nodeConfig); err != nil {
			utils.Logger.Error().Err(err).Str("nodeType", nodeConfigs.FileNode.Type.ToString()).Msgf("Failed to initialize node")
		} else {
			initializedNodes = append(initializedNodes, fileNode)
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

func InitializeNode(server src.AgentGRPCServer, config *configs.NodeConfig) (src.P2PNoder, error) {
	var n src.P2PNoder
	baseNode := p2p.NewBaseNode(config)
	switch config.Type {
	case configs.FilesNodeType:
		storage := storages.NewLocalStorage(&config.Storage)
		n = file_node.NewP2PFilesNode(baseNode, storage)
	default:
		return nil, errors.New("invalid node type")
	}
	n.Register(server.ServerObj())
	err := n.ConnectToBootstrapPeers(server)
	if err != nil {
		utils.Logger.Warn().Msgf("Bootstrapping process failed")
	}
	options := n.Options()
	utils.Logger.Info().
		Str("id", options.ID).
		Str("type", options.Type.ToString()).
		Msg("Initialized node successfully.")
	return n, nil
}
