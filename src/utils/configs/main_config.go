package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

var (
	MainConfig *Config
)

func LoadConfig(file string) {
	dir := filepath.Dir(file)
	viper.AddConfigPath(dir)
	viper.SetConfigName("config") // Register config file name (no extension)
	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if file != "" {
		// If a config file was specified but not found, return an error
		log.Fatalf("Config file not found: %s", file)
	}
	err := viper.Unmarshal(&MainConfig)
	if err != nil {
		log.Fatalf("Unable to decode into config struct, %v", err)
	}
}

type Config struct {
	ServerConfig ServerConfig  `mapstructure:"server"`
	Nodes        []NodeConfigs `mapstructure:"nodes"`
	LoggerMode   string        `mapstructure:"logger_mode"`
	LogLevel     string        `mapstructure:"log_level"`
}
