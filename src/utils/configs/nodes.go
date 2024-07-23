package configs

import "time"

type P2PNodeConfig struct {
	Type                 string                 `json:"type"`
	ID                   string                 `json:"id"`
	BootstrapPeerAddrs   []string               `json:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration          `json:"bootstrap_node_timeout"`
	ExtraConfig          map[string]interface{} `json:"extra_config"`
}

type FileSystemNodeConfig struct {
	LocalStorageConfig
}
