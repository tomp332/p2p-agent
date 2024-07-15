package util

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Logger zerolog.Logger
)

func SetupLogger(config *Config) {
	if config.LoggerMode == "dev" {
		setupDevLogger(config)
	} else if config.LoggerMode == "prod" {
		setupProdLogger()
	} else {
		log.Fatal().Msg("Invalid logger mode specified.")
	}
}

func setupDevLogger(config *Config) {
	output := zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("|%-4s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.LogLevel != "" {
		if strings.ToUpper(config.LogLevel) == "DEBUG" {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	}
	Logger = zerolog.New(output).With().Timestamp().Logger()
}

func setupProdLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	Logger = log.With().Caller().Logger()
}
