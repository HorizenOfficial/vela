package logger

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// TODO
// UDP/TCP split based on log severity
// Support for vsock connection
// Async write mode
// - Push the log message into an in-memory queue (channel, ring buffer, etc.)
// - Have a dedicated goroutine consume the queue and write to the network
// Buffered memory queue on disconnection

func init() {
	// Zerolog's default internal skip is usually 2.
	// By setting it to 3, we are adding 1 extra skip for the wrapper function.
	zerolog.CallerSkipFrameCount = 3
}

// ReconnectingWriter is a resilient io.Writer that automatically reconnects.
type ReconnectingWriter struct {
	mu             sync.Mutex
	conn           net.Conn
	cfg            *Config
	fallbackLogger Logger
	retryDelay     time.Duration
	maxWait        time.Duration
}

// NewReconnectingWriter creates a new writer instance
func NewReconnectingWriter(cfg *Config, fallback Logger) *ReconnectingWriter {
	return &ReconnectingWriter{
		cfg:            cfg,
		fallbackLogger: fallback,
		retryDelay:     2 * time.Second,
		maxWait:        30 * time.Second,
	}
}

// Write implements the io.Writer interface. It handles connection and writing.
func (w *ReconnectingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// If not connected, try to connect
	if w.conn == nil {
		if err := w.connectWithRetry(); err != nil {
			// Connection failed, use fallback
			w.fallbackLogger.Error("Failed to connect to remote log server, using fallback: %v", err)
			return w.fallbackLogger.(io.Writer).Write(p)
		}
	}

	// Write to the connection
	n, err = w.conn.Write(p)
	if err != nil {
		w.fallbackLogger.Error("Failed to write to remote log server: %v. Will attempt to reconnect on next log.", err)
		// Mark as disconnected, so the next write triggers a reconnect
		w.conn.Close()
		w.conn = nil
	}

	return n, err
}

// connectWithRetry tries to connect multiple times until maxWait
func (w *ReconnectingWriter) connectWithRetry() error {
	deadline := time.Now().Add(w.maxWait)

	for {
		conn, err := net.Dial(w.cfg.RemoteLogNetwork, w.cfg.RemoteLogAddress)
		if err == nil {
			w.conn = conn
			fmt.Printf("[zeronetwork] connected to %s\n", w.cfg.RemoteLogAddress)
			return nil
		}

		if time.Now().After(deadline) {
			return errors.New("could not connect to remote logger in time")
		}

		fmt.Println("[zeronetwork] server not ready, retrying...")
		time.Sleep(w.retryDelay)
	}
}

// Close closes the connection
func (w *ReconnectingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}

// ZeroNetworkLogger implements the Logger interface.
type ZeroNetworkLogger struct {
	logger *zerolog.Logger
	writer *ReconnectingWriter
}

// NewZeroNetworkLogger creates a new logger instance
func NewZeroNetworkLogger(cfg *Config) *ZeroNetworkLogger {
	if cfg.RemoteLogAddress == "" {
		panic("RemoteLogAddress cannot be empty")
	}
	if cfg.RemoteLogNetwork == "" {
		cfg.RemoteLogNetwork = "tcp"
	}

	fallback := NewPrintfLogger(&Config{Console: true, ConsoleLevel: "error"})
	writer := NewReconnectingWriter(cfg, fallback)

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Caller(). // see init() func above
		Logger()

	return &ZeroNetworkLogger{
		logger: &logger,
		writer: writer,
	}
}

// Logging methods just pass through to the underlying zerolog instance
func (z *ZeroNetworkLogger) Trace(msg string, args ...any) { z.logger.Trace().Msgf(msg, args...) }
func (z *ZeroNetworkLogger) Debug(msg string, args ...any) { z.logger.Debug().Msgf(msg, args...) }
func (z *ZeroNetworkLogger) Info(msg string, args ...any)  { z.logger.Info().Msgf(msg, args...) }
func (z *ZeroNetworkLogger) Warn(msg string, args ...any)  { z.logger.Warn().Msgf(msg, args...) }
func (z *ZeroNetworkLogger) Error(msg string, args ...any) { z.logger.Error().Msgf(msg, args...) }
func (z *ZeroNetworkLogger) Fatal(msg string, args ...any) { z.logger.Fatal().Msgf(msg, args...) }
func (z *ZeroNetworkLogger) Panic(msg string, args ...any) { z.logger.Panic().Msgf(msg, args...) }

// Close safely closes the network connection
func (z *ZeroNetworkLogger) Close() error {
	return z.writer.Close()
}
