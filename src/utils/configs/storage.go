package configs

type LocalStorageConfig struct {
	RootDirectory  string  `mapstructure:"root_directory"`
	MaxFileSize    float64 `mapstructure:"max_file_size"`
	MaxStorageSize float64 `mapstructure:"max_size"`
}
