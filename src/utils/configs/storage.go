package configs

type LocalStorageConfig struct {
	RootDirectory  string  `json:"root_directory"`
	MaxFileSize    float64 `json:"max_file_size"`
	MaxStorageSize float64 `json:"max_storage_size"`
}
