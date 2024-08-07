package utils

import (
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"testing"
)

func BaseNodeConfig(t *testing.T) *configs.NodeConfigs {
	return &configs.NodeConfigs{
		BaseNodeConfigs: configs.P2PNodeBaseConfig{},
	}
}

func FileNodeConfig(t *testing.T) *configs.NodeConfigs {
	tempDir := t.TempDir()
	return &configs.NodeConfigs{
		BaseNodeConfigs: configs.P2PNodeBaseConfig{
			Type:                 "file",
			BootstrapNodeTimeout: 10,
		},
		FilesNodeConfigs: configs.P2PFilesNodeConfig{
			Storage: configs.LocalStorageConfig{
				RootDirectory: tempDir,
			},
		},
	}
}
