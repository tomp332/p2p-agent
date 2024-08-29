package storage

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg"
	"github.com/tomp332/p2p-agent/pkg/utils"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"github.com/tomp332/p2p-agent/pkg/utils/types"
	"io"
	"math"
	"os"
	"path/filepath"
	"sync"
)

type LocalStorage struct {
	options          *configs.LocalStorageConfig
	maxFileSizeBytes int64
	mainPath         string
}

func NewLocalStorage(options *configs.LocalStorageConfig) *LocalStorage {
	s := &LocalStorage{
		options: options,
	}
	err := s.Initialize()
	if err != nil {
		log.Fatal().Msgf("Failed to initialize local storage, %s", err.Error())
		return nil
	}
	return s
}

func (s *LocalStorage) Exists(fileName string) bool {
	path := filepath.Join(s.mainPath, fileName)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Debug().Str("fileName", fileName).Msg("Local file does not exist")
		return false
	}
	log.Debug().Str("fileName", fileName).Msg("Local file exists")
	return true
}

func (s *LocalStorage) Initialize() error {
	s.mainPath = filepath.Join(pkg.LocalStorageDefaultDir, s.options.RootDirectory)
	if _, err := os.Stat(s.mainPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(s.mainPath, os.ModePerm); err != nil {
			log.Error().Msgf("Failed to create default storage directory, %s", err.Error())
			return err
		}
	}
	if s.options.MaxStorageSize == 0 {
		s.options.MaxStorageSize = math.MaxInt64
	}
	if s.options.MaxFileSize == 0 {
		s.maxFileSizeBytes = 4 * 1024 * 1024 * 1024
	} else {
		s.maxFileSizeBytes = int64(s.options.MaxFileSize * 1024 * 1024 * 1024)
	}
	log.Debug().
		Str("storageDir", s.mainPath).
		Int64("maxFileSizeBytes", s.maxFileSizeBytes).
		Float64("MaxStorageSize", s.options.MaxStorageSize).
		Msg("Initialized local storage successfully.")
	return nil
}

func (s *LocalStorage) Put(ctx context.Context, fileName string, fileDataChan <-chan *types.TransferChunkData, wg *sync.WaitGroup) error {
	//defer wg.Done()
	if fileName == "" {
		fileName = utils.GenerateRandomID()
		log.Warn().Str("fileName", fileName).Msg("No file name specified to local storage write, generated random.")
	}
	log.Debug().Str("fileName", fileName).Msg("Started receiving file data on local storage")
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		file := &File{FilePath: filepath.Join(s.mainPath, fileName)}
		defer file.Close()
		for data := range fileDataChan {
			log.Debug().Int("chunkSize", data.ChunkSize).Str("fileName", fileName).Msg("Received file data on local storage writer")
			err := file.Write(data.ChunkData)
			if err != nil {
				log.Error().Err(err).Str("fileName", fileName).Msg("Failed to write to local file")
				return err
			}
		}
	}
	log.Debug().Str("fileName", fileName).Msg("Done writing data to local storage.")
	return nil
}

func (s *LocalStorage) Get(ctx context.Context, fileName string) (<-chan *types.TransferChunkData, error) {
	filePath := filepath.Join(s.mainPath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Str("fileName", fileName).Msgf("File does not exist locally.")
		}
		return nil, err
	}
	dataChan := make(chan *types.TransferChunkData)
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
					dataChan <- &types.TransferChunkData{
						ChunkData: readBytes[:numBytesRead],
						ChunkSize: numBytesRead,
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

func (s *LocalStorage) Delete(ctx context.Context, fileName string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	file := &File{FilePath: filepath.Join(s.mainPath, fileName)}
	log.Debug().Str("fileName", fileName).Msg("Deleting file from storage")

	err := file.Delete()
	if err != nil {
		log.Error().Err(err).Str("fileName", fileName).Msg("Failed to delete local file")
		return err
	}

	log.Info().Str("fileName", fileName).Msg("File removed successfully from storage")
	return nil
}
