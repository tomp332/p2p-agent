package src

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Logger zerolog.Logger
)

func init() {
	SetupLogger()
}

func SetupLogger() {
	if MainConfig.LoggerMode == "dev" {
		setupDevLogger()
	} else if MainConfig.LoggerMode == "prod" {
		setupProdLogger()
	} else {
		log.Fatal().Msg("Invalid logger mode specified.")
	}
}

func setupDevLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("|%-4s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if MainConfig.LogLevel != "" {
		if strings.ToUpper(MainConfig.LogLevel) == "DEBUG" {
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
