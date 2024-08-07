package configs

import (
	"time"
)

const (
	FilesNodeType NodeType = "FileNode"
)

type NodeType string

// ToString method with a value receiver.
func (n NodeType) ToString() string {
	return string(n)
}

func StrToNodeType(str string) NodeType {
	return NodeType(str)
}

type AllNodesMap struct {
	FileNode NodeConfig `mapstructure:"file-nodes"`
}

type NodeConfig struct {
	ID                   string        `yaml:"id"`
	BootstrapPeerAddrs   []string      `mapstructure:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration `mapstructure:"bootstrap_node_timeout"`
	Type                 NodeType

	// All nodes type configurations
	FilesNodeConfig `mapstructure:"file_config"`
}

type FilesNodeConfig struct {
	Storage LocalStorageConfig `mapstructure:"storage"`
}
