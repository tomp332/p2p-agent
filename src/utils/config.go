package utils

import (
	"encoding/json"
	"os"
)

var (
	MainConfig *Config
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type Config struct {
	ServerConfig ServerConfig             `json:"server_config"`
	Nodes        []map[string]interface{} `json:"nodes"`
	LoggerMode   string                   `json:"logger_mode"`
	LogLevel     string                   `json:"log_level"`
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
