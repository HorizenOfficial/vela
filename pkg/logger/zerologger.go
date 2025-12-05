package logger

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	mu      sync.RWMutex
	logger  zerolog.Logger
	logFile *os.File
}

func NewZeroLogger(cfg *Config) *ZeroLogger {
	writers := []io.Writer{}

	var logFile *os.File
	if cfg.FileName != "" {
		var err error
		logFile, err = os.OpenFile(cfg.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		writers = append(writers, zerolog.ConsoleWriter{
			Out:     logFile,
			NoColor: true,
		})
	}

	if cfg.Console {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    !cfg.ConsoleColor,
			TimeFormat: "2006-Dec-02 15:04:05.000",
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
		Caller(). // see init() func above
		Logger().
		Level(logLevel)

	return &ZeroLogger{logger: logger, logFile: logFile}
}

func (z *ZeroLogger) SetLevel(level string) error {
	z.mu.Lock()
	defer z.mu.Unlock()
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	z.logger = z.logger.Level(lvl)
	return nil
}

func (z *ZeroLogger) Trace(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Trace().Msgf(msg, args...)
}
func (z *ZeroLogger) Debug(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Debug().Msgf(msg, args...)
}
func (z *ZeroLogger) Info(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Info().Msgf(msg, args...)
}
func (z *ZeroLogger) Warn(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Warn().Msgf(msg, args...)
}
func (z *ZeroLogger) Error(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Error().Msgf(msg, args...)
}
func (z *ZeroLogger) Fatal(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Fatal().Stack().Msgf(msg, args...)
}
func (z *ZeroLogger) Panic(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Panic().Stack().Msgf(msg, args...)
}
func (z *ZeroLogger) Close() error {
	if z.logFile != nil {
		return z.logFile.Close()
	}
	return nil
}
