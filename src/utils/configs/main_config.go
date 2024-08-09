package configs

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
