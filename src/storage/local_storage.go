package storage

import (
	"context"
	"fmt"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"math"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	options          *configs.LocalStorageConfig
	maxFileSizeBytes float64
}

func NewLocalStorage(options *configs.LocalStorageConfig) *LocalStorage {
	s := &LocalStorage{
		options: options,
	}
	err := s.Initialize()
	if err != nil {
		utils.Logger.Error().Msgf("Failed to initialize local storage, %s", err.Error())
		return nil
	}
	return s
}

func (s *LocalStorage) Initialize() error {
	newPath := filepath.Join(s.options.RootDirectory, src.LocalStorageDefaultDir)
	s.options.RootDirectory = newPath
	if s.options.MaxStorageSize == 0 {
		s.options.MaxStorageSize = math.MaxInt64
	}
	if s.options.MaxFileSize == 0 {
		s.maxFileSizeBytes = 4 * 1024 * 1024 * 1024
	} else {
		s.maxFileSizeBytes = s.options.MaxFileSize * 1024 * 1024 * 1024
	}
	if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
		utils.Logger.Error().
			Str("rootDirectory", s.options.RootDirectory).
			Msgf("Error initializing local file system storage. Failed to create root directory.")
		return err
	}
	utils.Logger.Debug().
		Str("rootDirectory", s.options.RootDirectory).
		Float64("maxFileSizeBytes", s.maxFileSizeBytes).
		Float64("MaxStorageSize", s.options.MaxStorageSize).
		Msg("Initialized local storage successfully.")
	return nil
}

func (s *LocalStorage) Put(ctx context.Context, fileID string, dataChan <-chan []byte) (float64, error) {
	filePath := filepath.Join(s.options.RootDirectory, fileID)
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var totalSize float64
	done := make(chan error)
	go func() {
		for chunk := range dataChan {
			n, err := file.Write(chunk)
			if err != nil {
				done <- err
				return
			}
			totalSize += float64(n)
			if totalSize > s.maxFileSizeBytes {
				done <- fmt.Errorf("file size exceeds the maximum allowed size of %.2f GB", s.options.MaxFileSize)
				return
			}
		}
		done <- nil
	}()

	select {
	case <-ctx.Done():
		os.Remove(filePath)
		return 0, ctx.Err()
	case err := <-done:
		if err != nil {
			utils.Logger.Error().Msgf("Error uploading file: %v", err)
			os.Remove(filePath)
			return 0, err
		}
	}
	return totalSize, nil
}
