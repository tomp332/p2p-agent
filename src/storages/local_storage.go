package storages

import (
	"context"
	"errors"
	"fmt"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/utils"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"io"
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
		utils.Logger.Error().Msgf("Failed to initialize local storages, %s", err.Error())
		return nil
	}
	return s
}

func (s *LocalStorage) Initialize() error {
	newPath := filepath.Join(s.options.RootDirectory, src.LocalStorageDefaultDir)
	if _, err := os.Stat(newPath); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(newPath, os.ModePerm); err != nil {
			utils.Logger.Error().Msgf("Failed to create default storages directory, %s", err.Error())
			return err
		}
	}
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
			Msgf("Error initializing local file system storages. Failed to create root directory.")
		return err
	}
	utils.Logger.Debug().
		Str("rootDirectory", s.options.RootDirectory).
		Float64("maxFileSizeBytes", s.maxFileSizeBytes).
		Float64("MaxStorageSize", s.options.MaxStorageSize).
		Msg("Initialized local storages successfully.")
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
	utils.Logger.Info().Str("fileId", fileID).Msg("File uploaded successfully.")
	return totalSize, nil
}

func (s *LocalStorage) Search(fileID string) bool {
	filePath := filepath.Join(s.options.RootDirectory, fileID)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (s *LocalStorage) Get(ctx context.Context, fileID string) (<-chan []byte, error) {
	filePath := filepath.Join(s.options.RootDirectory, fileID)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			utils.Logger.Warn().Msgf("File %s does not exist.", fileID)
		}
		return nil, err
	}
	dataChan := make(chan []byte)
	go func() {
		defer file.Close()
		defer close(dataChan)

		readBytes := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				utils.Logger.Info().Msg("Read file from storages canceled.")
				return
			default:
				numBytesRead, err := file.Read(readBytes)
				if err != nil {
					if err != io.EOF {
						utils.Logger.Error().Err(err).Msg("Failed to read data from file")
					}
					return
				}
				if numBytesRead > 0 {
					dataChan <- readBytes[:numBytesRead]
				}
				if numBytesRead < len(readBytes) {
					utils.Logger.Debug().Msg("Read file from storages successfully.")
					return
				}
			}
		}
	}()
	return dataChan, nil
}

func (s *LocalStorage) Delete(ctx context.Context, fileID string) error {
	filePath := filepath.Join(s.options.RootDirectory, fileID)
	utils.Logger.Debug().Str("fileId", fileID).Msgf("Deleteing file from storages")
	err := os.Remove(filePath)
	if err != nil {
		return ctx.Err()
	}
	utils.Logger.Info().Str("fileId", fileID).Msgf("File removed successfully from storages")
	return nil
}
