package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	logger zerolog.Logger
}

func NewZeroLogger(cfg *Config) *ZeroLogger {
	writers := []io.Writer{}

	if cfg.FileName != "" {
		logFile, err := os.OpenFile(cfg.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		writers = append(writers, logFile)
	}

	if cfg.Console {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    !cfg.ConsoleColor,
		})
	}

	var writer io.Writer
	if len(writers) > 0 {
		writer = zerolog.MultiLevelWriter(writers...)
	} else {
		// default to stderr if no output is specified
		writer = os.Stderr
	}

	logLevel, err := zerolog.ParseLevel(cfg.ConsoleLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(logLevel)

	return &ZeroLogger{logger: logger}
}

func (z *ZeroLogger) Debug(msg string, args ...any) { z.logger.Debug().Msgf(msg, args...) }
func (z *ZeroLogger) Info(msg string, args ...any)  { z.logger.Info().Msgf(msg, args...) }
func (z *ZeroLogger) Warn(msg string, args ...any)  { z.logger.Warn().Msgf(msg, args...) }
func (z *ZeroLogger) Error(msg string, args ...any) { z.logger.Error().Msgf(msg, args...) }
func (z *ZeroLogger) Fatal(msg string, args ...any) { z.logger.Fatal().Stack().Msgf(msg, args...) }
