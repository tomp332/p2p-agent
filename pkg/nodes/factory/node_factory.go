package factory

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/nodes"
	"github.com/tomp332/p2p-agent/pkg/nodes/fsNode"
	"github.com/tomp332/p2p-agent/pkg/server"
	"github.com/tomp332/p2p-agent/pkg/storage"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"reflect"
)

func InitializeNodes(ctx context.Context, server *server.GRPCServer, nodeConfigs *configs.AllNodesMap) ([]nodes.P2PNoder, error) {
	var initializedNodes []nodes.P2PNoder
	val := reflect.ValueOf(nodeConfigs)
	// Dereference the pointer to access the struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// Get the number of fields in the struct
	numFields := val.NumField()

	for i := 0; i < numFields; i++ {
		// Check if the context has been canceled
		select {
		case <-ctx.Done():
			log.Warn().Msg("Initialization cancelled by context")
			return initializedNodes, ctx.Err()
		default:
			// Continue initialization if context is still active
		}

		field := val.Type().Field(i)
		value := val.Field(i)
		nodeConfig := value.Interface().(*configs.NodeConfig)
		if nodeConfig == nil {
			continue
		}

		nodeConfig.Type = configs.StrToNodeType(field.Name)
		if fileNode, err := InitializeNode(ctx, server, nodeConfig); err != nil {
			log.Error().Err(err).Str("nodeType", nodeConfig.Type.ToString()).Msgf("Failed to initialize node")
		} else {
			initializedNodes = append(initializedNodes, fileNode)
		}
	}

	if initializedNodes == nil {
		return nil, errors.New("no nodes have been configured")
	}
	return initializedNodes, nil
}

func InitializeNode(mainCtx context.Context, server *server.GRPCServer, config *configs.NodeConfig) (nodes.P2PNoder, error) {
	var n nodes.P2PNoder
	baseNode := nodes.NewBaseNode(config, mainCtx)
	switch config.Type {
	case configs.FilesNodeType:
		if config.Storage.RootDirectory == "" {
			config.Storage.RootDirectory = configs.MainConfig.ID
		}
		localStorage := storage.NewLocalStorage(&config.Storage)
		n = fsNode.NewP2PFilesNode(baseNode, localStorage)
	default:
		return nil, errors.New("invalid nodes type")
	}
	n.InterceptorsRegister(server)
	n.ServiceRegister(server)
	n.ConnectToBootstrapPeers()
	options := n.Options()
	log.Info().
		Str("id", configs.MainConfig.ID).
		Str("type", options.Type.ToString()).
		Msg("Initialized nodes successfully.")
	return n, nil
}
