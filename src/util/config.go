package util

import (
	"encoding/json"
	"os"
	"time"
)

type NodeConfig struct {
	Type                 string            `json:"type"`
	Storage              map[string][]byte `json:"storage"`
	BootstrapPeerAddrs   []string          `json:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration     `json:"bootstrap_node_timeout"`
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

func LoadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&config)
	return &config, err
}
