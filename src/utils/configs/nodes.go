package configs

import (
	"time"
)

type NodeConfigs struct {
	Type                 string        `yaml:"type"`
	ID                   string        `yaml:"id"`
	BootstrapPeerAddrs   []string      `mapstructure:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration `mapstructure:"bootstrap_node_timeout"`
	// All node type configurations
	P2PFilesNodeConfig `mapstructure:"file_config"`
}

type P2PFilesNodeConfig struct {
	Storage LocalStorageConfig `mapstructure:"storage"`
}
