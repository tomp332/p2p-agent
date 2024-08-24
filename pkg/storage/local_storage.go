package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"github.com/tomp332/p2p-agent/pkg/utils/types"
	"io"
	"math"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	options          *configs.LocalStorageConfig
	maxFileSizeBytes int64
}

func NewLocalStorage(options *configs.LocalStorageConfig) *LocalStorage {
	s := &LocalStorage{
		options: options,
	}
	err := s.Initialize()
	if err != nil {
		log.Error().Msgf("Failed to initialize local storage, %s", err.Error())
		return nil
	}
	return s
}

func (s *LocalStorage) Initialize() error {
	mainPath := filepath.Join(pkg.LocalStorageDefaultDir, s.options.RootDirectory)
	if _, err := os.Stat(mainPath); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(mainPath, os.ModePerm); err != nil {
			log.Error().Msgf("Failed to create default storage directory, %s", err.Error())
			return err
		}
	}
	s.options.RootDirectory = mainPath
	if s.options.MaxStorageSize == 0 {
		s.options.MaxStorageSize = math.MaxInt64
	}
	if s.options.MaxFileSize == 0 {
		s.maxFileSizeBytes = 4 * 1024 * 1024 * 1024
	} else {
		s.maxFileSizeBytes = int64(s.options.MaxFileSize * 1024 * 1024 * 1024)
	}
	if err := os.MkdirAll(mainPath, os.ModePerm); err != nil {
		log.Error().
			Str("rootDirectory", s.options.RootDirectory).
			Msgf("Error initializing local file system storage. Failed to create root directory.")
		return err
	}
	log.Debug().
		Str("rootDirectory", s.options.RootDirectory).
		Int64("maxFileSizeBytes", s.maxFileSizeBytes).
		Float64("MaxStorageSize", s.options.MaxStorageSize).
		Msg("Initialized local storage successfully.")
	return nil
}

func (s *LocalStorage) Put(ctx context.Context, fileID string, dataChan <-chan types.TransferChunkData) (int64, error) {
	filePath := filepath.Join(s.options.RootDirectory, fileID)
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var totalSize int64
	done := make(chan error, 1) // Buffer of 1 to handle the single error case

	go func() {
		defer close(done)
		for data := range dataChan {
			n, writeErr := file.Write(data.ChunkData)
			if writeErr != nil {
				log.Error().Err(writeErr).Msg("Error writing chunk to file")
				done <- writeErr
				return
			}
			totalSize += int64(n)
			if totalSize > s.maxFileSizeBytes {
				done <- fmt.Errorf("file size exceeds the maximum allowed size of %.2f GB", s.options.MaxFileSize)
				return
			}
		}
		done <- nil
	}()

	select {
	case <-ctx.Done():
		// Context was canceled, but do not delete the file
		return totalSize, ctx.Err()
	case uploadErr := <-done:
		if uploadErr != nil {
			log.Error().Err(uploadErr).Msg("Error putting file in local storage")
			os.Remove(filePath) // Attempt to remove the partially written file
			return 0, uploadErr
		}
	}

	log.Info().Str("fileId", fileID).Msg("File uploaded successfully.")
	return totalSize, nil
}

func (s *LocalStorage) Get(ctx context.Context, fileID string) (<-chan types.TransferChunkData, error) {
	filePath := filepath.Join(s.options.RootDirectory, fileID)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Msgf("File %s does not exist.", fileID)
		}
		return nil, err
	}
	dataChan := make(chan types.TransferChunkData)
	go func() {
		defer file.Close()
		defer close(dataChan)

		readBytes := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Read file from storage canceled.")
				return
			default:
				numBytesRead, readErr := file.Read(readBytes)
				if readErr != nil {
					if readErr != io.EOF {
						log.Error().Err(readErr).Msg("Failed to read data from file")
					}
					return
				}
				if numBytesRead > 0 {
					dataChan <- types.TransferChunkData{
						ChunkData: readBytes[:numBytesRead],
						ChunkSize: int64(numBytesRead),
					}
				}
				if numBytesRead < len(readBytes) {
					log.Debug().Msg("Read file from storage successfully.")
					return
				}
			}
		}
	}()
	return dataChan, nil
}

func (s *LocalStorage) Delete(ctx context.Context, fileID string) error {
	filePath := filepath.Join(s.options.RootDirectory, fileID)
	log.Debug().Str("fileId", fileID).Msgf("Deleteing file from storage")
	err := os.Remove(filePath)
	if err != nil {
		return ctx.Err()
	}
	log.Info().Str("fileId", fileID).Msgf("File removed successfully from storage")
	return nil
}
