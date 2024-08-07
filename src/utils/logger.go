package utils

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Logger zerolog.Logger
)

func SetupLogger() {
	loggerMode := configs.MainConfig.LoggerMode
	if loggerMode != "" {
		switch loggerMode {
		case "dev":
			setupDevLogger()
		default:
			setupProdLogger()
		}
	}
}

func setupDevLogger() {
	output := zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("|%-4s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if configs.MainConfig.LogLevel != "" {
		if strings.ToUpper(configs.MainConfig.LogLevel) == "DEBUG" {
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
