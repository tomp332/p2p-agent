package configs

import (
	"time"
)

type NodeConfigs struct {
	BaseNodeConfigs P2PNodeBaseConfig `json:"node_config"`
	// All node type configurations
	FilesNodeConfigs P2PFilesNodeConfig `json:"files_node_config"`
}

type P2PNodeBaseConfig struct {
	Type                 string        `json:"type" default:"base" validate:"type-validate"`
	ID                   string        `json:"id"`
	BootstrapPeerAddrs   []string      `json:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration `json:"bootstrap_node_timeout" default:"10"`
}

type P2PFilesNodeConfig struct {
	Storage LocalStorageConfig
}
