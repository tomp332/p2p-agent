package src

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	MainConfig *Config
)

func init() {
	if err := LoadConfig("config.json"); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	ServerPort           int32         `json:"server_port"`
	BootstrapPeerAddrs   []string      `json:"bootstrap_peer_addrs"`
	BootstrapNodeTimeout time.Duration `json:"bootstrap_node_timeout"`
	LogLevel             string        `json:"log_level"`
	LoggerMode           string        `json:"logger_mode"`
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	decoder := json.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}
	MainConfig = config
	return nil
}
