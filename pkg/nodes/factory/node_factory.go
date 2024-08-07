package factory

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/nodes"
	"github.com/tomp332/p2p-agent/pkg/nodes/file_node"
	"github.com/tomp332/p2p-agent/pkg/storage"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"google.golang.org/grpc"
	"reflect"
)

func InitializeNodes(server *grpc.Server, nodeConfigs *configs.AllNodesMap) ([]nodes.P2PNoder, error) {
	var initializedNodes []nodes.P2PNoder
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
			log.Error().Err(err).Str("nodeType", nodeConfigs.FileNode.Type.ToString()).Msgf("Failed to initialize nodes")
		} else {
			initializedNodes = append(initializedNodes, fileNode)
		}
	}
	if initializedNodes == nil {
		return nil, errors.New("failed to initialize nodes")
	}
	return initializedNodes, nil
}

func InitializeNode(server *grpc.Server, config *configs.NodeConfig) (nodes.P2PNoder, error) {
	var n nodes.P2PNoder
	baseNode := nodes.NewBaseNode(config)
	switch config.Type {
	case configs.FilesNodeType:
		localStorage := storage.NewLocalStorage(&config.Storage)
		n = file_node.NewP2PFilesNode(baseNode, localStorage)
	default:
		return nil, errors.New("invalid nodes type")
	}
	n.Register(server)
	err := n.ConnectToBootstrapPeers(server)
	if err != nil {
		log.Warn().Msgf("Bootstrapping process failed")
	}
	options := n.Options()
	log.Info().
		Str("id", options.ID).
		Str("type", options.Type.ToString()).
		Msg("Initialized nodes successfully.")
	return n, nil
}
