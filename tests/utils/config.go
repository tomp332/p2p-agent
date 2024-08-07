package utils

import (
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"testing"
)

func BaseNodeConfig() *configs.NodeConfigs {
	return &configs.NodeConfigs{}
}

func FileNodeConfig(t *testing.T) *configs.NodeConfigs {
	tempDir := t.TempDir()
	return &configs.NodeConfigs{
		Type:                 "file",
		BootstrapNodeTimeout: 10,
		P2PFilesNodeConfig: configs.P2PFilesNodeConfig{
			Storage: configs.LocalStorageConfig{
				RootDirectory: tempDir,
			},
		},
	}
}
