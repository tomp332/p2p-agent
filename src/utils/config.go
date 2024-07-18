package utils

import (
	"encoding/json"
	"os"
	"time"
)

var (
	MainConfig *Config
)

type NodeConfig struct {
	Type                 string        `json:"type"`
	BootstrapPeerAddrs   []string      `json:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration `json:"bootstrap_node_timeout"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type Config struct {
	ServerConfig ServerConfig `json:"server_config"`
	Nodes        []NodeConfig `json:"nodes"`
	LoggerMode   string       `json:"logger_mode"`
	LogLevel     string       `json:"log_level"`
}

func LoadConfig(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&MainConfig)
	return err
}
