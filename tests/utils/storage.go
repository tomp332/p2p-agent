package utils

import (
	"github.com/tomp332/p2p-agent/src/storage"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"os"
	"testing"
)

func CreateTempStorage(t *testing.T) (*storage.LocalStorage, func()) {
	tempDir, err := os.MkdirTemp("", "storage_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	storageConfig := &configs.LocalStorageConfig{
		RootDirectory: tempDir,
	}
	localStorage := storage.NewLocalStorage(storageConfig)

	cleanup := func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	return localStorage, cleanup
}
