package configs

type LocalStorageConfig struct {
	RootDirectory  string  `json:"root_directory" default:"/tmp"`
	MaxFileSize    float64 `json:"max_file_size" default:"99999999"`
	MaxStorageSize float64 `json:"max_storage_size" default:"10000000"`
}
