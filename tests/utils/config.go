package utils

import "github.com/tomp332/p2p-agent/src/utils/configs"

func SetupConfig() {
	configs.MainConfig = &configs.Config{
		LoggerMode: "dev",
	}
}
