package configs

import (
	"encoding/json"
	"os"
)

var (
	MainConfig *Config
)

type Config struct {
	ServerConfig  ServerConfig       `json:"server"`
	Nodes         []NodeConfigs      `json:"nodes"`
	StorageConfig LocalStorageConfig `json:"storage"`
	LoggerMode    string             `json:"logger_mode"`
	LogLevel      string             `json:"log_level"`
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
