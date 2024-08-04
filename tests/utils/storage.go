package utils

import (
	"github.com/tomp332/p2p-agent/src/storage"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"os"
)

func CreateTempStorage() (*storage.LocalStorage, func()) {
	tempDir, err := os.MkdirTemp("", "storage_test")
	if err != nil {
		panic(err)
	}

	storageConfig := &configs.LocalStorageConfig{
		RootDirectory: tempDir,
	}
	localStorage := storage.NewLocalStorage(storageConfig)

	cleanup := func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			panic(err)
		}
	}

	return localStorage, cleanup
}
