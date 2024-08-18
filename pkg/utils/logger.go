package utils

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func SetupLogger() {
	var logger zerolog.Logger
	loggerMode := configs.MainConfig.LoggerMode
	if loggerMode != "" {
		switch loggerMode {
		case "dev":
			logger = setupDevLogger()
		default:
			logger = setupProdLogger()
		}
	}
	log.Logger = logger
}

func setupDevLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("|%-4s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if configs.MainConfig.LogLevel != "" {
		switch strings.ToUpper(configs.MainConfig.LogLevel) {
		case "DEBUG":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "TRACE":
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		case "INFO":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "WARNING":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	}
	return zerolog.New(output).With().Timestamp().Logger()
}

func setupProdLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	return log.With().Caller().Logger()
}
