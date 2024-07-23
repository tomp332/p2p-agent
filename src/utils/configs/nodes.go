package configs

import "time"

type NodeConfigs struct {
	BaseNodeConfigs P2PNodeBaseConfig `json:"node_configs"`
	// All node type configurations
	FilesNodeConfigs P2PFilesNodeConfig `json:"files_node_config"`
}

type P2PNodeBaseConfig struct {
	Type                 string        `json:"type"`
	ID                   string        `json:"id"`
	BootstrapPeerAddrs   []string      `json:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration `json:"bootstrap_node_timeout"`
}

type P2PFilesNodeConfig struct {
	Storage LocalStorageConfig
}
