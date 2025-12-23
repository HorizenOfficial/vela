package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	mu            sync.RWMutex
	fileLogger    *zerolog.Logger
	consoleLogger *zerolog.Logger
	logFile       *os.File
}

func init() {
	// Zerolog's default internal skip is usually 2.
	// By setting it to 3, we are adding 1 extra skip for the wrapper function.
	zerolog.CallerSkipFrameCount = 3
	// print timestamp in human readable format with milliseconds precision
	//zerolog.TimeFieldFormat = time.StampMilli
	zerolog.TimeFieldFormat = TimeStampFormatMs

}

func NewZeroLogger(cfg *Config) *ZeroLogger {
	var fileLogger *zerolog.Logger
	var consoleLogger *zerolog.Logger
	var logFile *os.File

	if cfg.FileName != "" {
		var err error
		logFile, err = os.OpenFile(cfg.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		fileLevel, err := zerolog.ParseLevel(cfg.FileLevel)
		if err != nil {
			fileLevel = zerolog.InfoLevel
		}

		l := zerolog.New(logFile).
			With().
			Timestamp().
			Caller().
			Logger().
			Level(fileLevel)
		fileLogger = &l
	}

	if cfg.Console {
		consoleLevel, err := zerolog.ParseLevel(cfg.ConsoleLevel)
		if err != nil {
			consoleLevel = zerolog.InfoLevel
		}

		writer := zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !cfg.ConsoleColor,
			// Zerolog ConsoleWriter does not use the global zerolog.TimeFieldFormat, fix it
			TimeFormat: TimeStampFormatMs,
		}

		l := zerolog.New(writer).
			With().
			Timestamp().
			Caller().
			Logger().
			Level(consoleLevel)
		consoleLogger = &l
	}

	return &ZeroLogger{fileLogger: fileLogger, consoleLogger: consoleLogger, logFile: logFile}
}

// we have two loggers here actually, for the time being we set both
func (z *ZeroLogger) SetLevel(level string) error {
	z.mu.Lock()
	defer z.mu.Unlock()
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	if z.fileLogger != nil {
		lgr := z.fileLogger.Level(lvl)
		z.fileLogger = &lgr
	}
	if z.consoleLogger != nil {
		lgr := z.consoleLogger.Level(lvl)
		z.consoleLogger = &lgr
	}
	return nil
}

func (z *ZeroLogger) Trace(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.consoleLogger != nil {
		z.consoleLogger.Trace().Msgf(msg, args...)
	}
	if z.fileLogger != nil {
		z.fileLogger.Trace().Msgf(msg, args...)
	}
}
func (z *ZeroLogger) Debug(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.consoleLogger != nil {
		z.consoleLogger.Debug().Msgf(msg, args...)
	}
	if z.fileLogger != nil {
		z.fileLogger.Debug().Msgf(msg, args...)
	}
}
func (z *ZeroLogger) Info(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.consoleLogger != nil {
		z.consoleLogger.Info().Msgf(msg, args...)
	}
	if z.fileLogger != nil {
		z.fileLogger.Info().Msgf(msg, args...)
	}
}
func (z *ZeroLogger) Warn(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.consoleLogger != nil {
		z.consoleLogger.Warn().Msgf(msg, args...)
	}
	if z.fileLogger != nil {
		z.fileLogger.Warn().Msgf(msg, args...)
	}
}
func (z *ZeroLogger) Error(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.consoleLogger != nil {
		z.consoleLogger.Error().Msgf(msg, args...)
	}
	if z.fileLogger != nil {
		z.fileLogger.Error().Msgf(msg, args...)
	}
}
func (z *ZeroLogger) Fatal(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.consoleLogger != nil {
		z.consoleLogger.Fatal().Stack().Msgf(msg, args...)
	}
	if z.fileLogger != nil {
		z.fileLogger.Fatal().Stack().Msgf(msg, args...)
	}
}
func (z *ZeroLogger) Panic(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.consoleLogger != nil {
		z.consoleLogger.Panic().Stack().Msgf(msg, args...)
	}
	if z.fileLogger != nil {
		z.fileLogger.Panic().Stack().Msgf(msg, args...)
	}
}
func (z *ZeroLogger) Close() error {
	if z.logFile != nil {
		return z.logFile.Close()
	}
	return nil
}
