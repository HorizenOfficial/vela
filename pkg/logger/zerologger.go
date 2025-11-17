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

func NewZeroLogger(level, format string) *ZeroLogger {
	var writer io.Writer
	if format == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	} else {
		writer = os.Stderr
	}

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Logger().
		Level(logLevel)

	return &ZeroLogger{logger: logger}
}

func (z *ZeroLogger) Debug(msg string, args ...any) { z.logger.Debug().Msgf(msg, args...) }
func (z *ZeroLogger) Info(msg string, args ...any)  { z.logger.Info().Msgf(msg, args...) }
func (z *ZeroLogger) Warn(msg string, args ...any)  { z.logger.Warn().Msgf(msg, args...) }
func (z *ZeroLogger) Error(msg string, args ...any) { z.logger.Error().Msgf(msg, args...) }
func (z *ZeroLogger) Fatal(msg string, args ...any) { z.logger.Fatal().Stack().Msgf(msg, args...) }
