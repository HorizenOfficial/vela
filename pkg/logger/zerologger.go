package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	logger  zerolog.Logger
	logFile *os.File
}

func init() {
	// Zerolog's default internal skip is usually 2.
	// By setting it to 3, we are adding 1 extra skip for the wrapper function.
	zerolog.CallerSkipFrameCount = 3
}

func NewZeroLogger(cfg *Config) *ZeroLogger {
	writers := []io.Writer{}

	var logFile *os.File
	if cfg.FileName != "" {
		fileLevel, err := zerolog.ParseLevel(cfg.FileLevel)
		if err != nil {
			fileLevel = zerolog.InfoLevel
		}

		logFile, err = os.OpenFile(cfg.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		fileLogger := zerolog.New(logFile).With().Logger().Level(fileLevel)
		writers = append(writers, fileLogger)
	}

	if cfg.Console {
		consoleLevel, err := zerolog.ParseLevel(cfg.ConsoleLevel)
		if err != nil {
			consoleLevel = zerolog.InfoLevel
		}
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    !cfg.ConsoleColor,
		}
		consoleLogger := zerolog.New(consoleWriter).With().Logger().Level(consoleLevel)
		writers = append(writers, consoleLogger)
	}

	var writer io.Writer
	if len(writers) > 0 {
		writer = zerolog.MultiLevelWriter(writers...)
	} else {
		// default to stderr if no output is specified
		writer = os.Stderr
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Caller(). // see init() func above
		Logger().
		Level(zerolog.TraceLevel) // Set global level to trace, filtering is done by individual writers

	return &ZeroLogger{logger: logger, logFile: logFile}
}

func (z *ZeroLogger) Trace(msg string, args ...any) { z.logger.Trace().Msgf(msg, args...) }
func (z *ZeroLogger) Debug(msg string, args ...any) { z.logger.Debug().Msgf(msg, args...) }
func (z *ZeroLogger) Info(msg string, args ...any)  { z.logger.Info().Msgf(msg, args...) }
func (z *ZeroLogger) Warn(msg string, args ...any)  { z.logger.Warn().Msgf(msg, args...) }
func (z *ZeroLogger) Error(msg string, args ...any) { z.logger.Error().Msgf(msg, args...) }
func (z *ZeroLogger) Fatal(msg string, args ...any) { z.logger.Fatal().Stack().Msgf(msg, args...) }
func (z *ZeroLogger) Panic(msg string, args ...any) { z.logger.Panic().Stack().Msgf(msg, args...) }
func (z *ZeroLogger) Close() error {
	if z.logFile != nil {
		return z.logFile.Close()
	}
	return nil
}
