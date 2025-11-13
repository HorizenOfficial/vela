package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	logger zerolog.Logger
}

func NewZeroLogger() *ZeroLogger {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &ZeroLogger{logger: logger}
}

func (z *ZeroLogger) Debug(msg string, args ...any)  { z.logger.Debug().Msgf(msg, args...) }
func (z *ZeroLogger) Info(msg string, args ...any)   { z.logger.Info().Msgf(msg, args...) }
func (z *ZeroLogger) Warn(msg string, args ...any)   { z.logger.Warn().Msgf(msg, args...) }
func (z *ZeroLogger) Error(msg string, args ...any)  { z.logger.Error().Msgf(msg, args...) }
func (z *ZeroLogger) Fatalf(msg string, args ...any) { z.logger.Fatal().Msgf(msg, args...) }
