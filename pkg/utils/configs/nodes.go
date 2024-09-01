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

type NodeAuthConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	JwtSecret string `yaml:"jwt_secret"`
}

type AllNodesMap struct {
	FileNode *NodeConfig `mapstructure:"file-node"`
}

type BootStrapNodeConnection struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type NodeConfig struct {
	BootstrapPeerAddrs   []BootStrapNodeConnection `mapstructure:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration             `mapstructure:"bootstrap_node_timeout"`
	Type                 NodeType
	Auth                 NodeAuthConfig `mapstructure:"auth"`
	// All nodes type configurations
	FilesNodeConfig `mapstructure:"file_config"`
}

type FilesNodeConfig struct {
	Storage LocalStorageConfig `mapstructure:"storage"`
}

type PeerHealthCheckConfig struct {
	HealthCheckInterval       time.Duration `mapstructure:"health_check_interval"`
	FailedHealthCheckInterval time.Duration `mapstructure:"failed_health_check_interval"`
}
