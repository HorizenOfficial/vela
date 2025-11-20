package logger

import (
	"io"
	"net"
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
		var err error
		logFile, err = os.OpenFile(cfg.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
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
		Caller(). // see init() func above
		Logger().
		Level(logLevel)

	return &ZeroLogger{logger: logger, logFile: logFile}
}
func initConnection() {
	// Connect to the remote TCP logging server
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create a Zerolog logger writing to the TCP connection
	logger := zerolog.New(conn).With().Timestamp().Logger()

	logger.Info().Msg("Hello from Zerolog over TCP!")
	logger.Warn().Str("component", "auth").Msg("Something looks suspicious")
	logger.Error().Err(err).Msg("Example error log")

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
