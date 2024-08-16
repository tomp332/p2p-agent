package configs

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
}
