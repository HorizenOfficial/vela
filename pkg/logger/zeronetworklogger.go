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

const (
	defaultBuffer      = 1000
	defaultRetryDelay  = 2 * time.Second
	defaultMaxWait     = 30 * time.Second
	defaultDialTimeout = 5 * time.Second
)

// AsyncWriter is a resilient, non-blocking io.Writer.
// It writes to a fallback writer immediately and buffers logs for an async worker.
type AsyncWriter struct {
	mu             sync.Mutex
	conn           net.Conn
	cfg            *Config
	fallbackWriter io.Writer
	logBuffer      chan []byte
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

// NewAsyncWriter creates a new writer that handles buffering and reconnecting.
func NewAsyncWriter(cfg *Config, fallback io.Writer) *AsyncWriter {
	writer := &AsyncWriter{
		cfg:            cfg,
		fallbackWriter: fallback,
		logBuffer:      make(chan []byte, defaultBuffer),
		stopChan:       make(chan struct{}),
	}

	writer.wg.Add(1)
	go writer.run() // Start the async worker

	return writer
}

// Write implements the io.Writer interface. It is non-blocking.
func (w *AsyncWriter) Write(p []byte) (n int, err error) {
	// Immediately write to fallback for real-time local logging
	if w.conn == nil {

		if _, err := w.fallbackWriter.Write(p); err != nil {
			fmt.Printf("Fallback logger failed: %v\n", err)
		}
	}

	// Create a copy, as the slice is reused by zerolog
	buf := make([]byte, len(p))
	copy(buf, p)

	// Send to buffer without blocking
	select {
	case w.logBuffer <- buf:
	default:
		// Buffer is full, drop the message but log it to fallback
		w.fallbackWriter.Write([]byte("Remote log buffer is full. Dropping message.\n"))
	}

	return len(p), nil
}

// run is the background worker for connecting and sending logs.
func (w *AsyncWriter) run() {
	defer w.wg.Done()
	for {
		select {
		case <-w.stopChan:
			return // Shutdown signal received
		default:
			// If not connected, block until we are.
			if w.conn == nil {
				if err := w.connectWithRetry(); err != nil {
					// Sleep before retrying connection
					time.Sleep(defaultRetryDelay)
					continue
				}
			}

			// Connection is live, start processing buffer.
			if !w.processBuffer() {
				// processBuffer returned false, meaning connection is dead.
				// The loop will now attempt to reconnect.
				continue
			}
		}
	}
}

// processBuffer reads from the channel and writes to the network.
// It returns `false` if the connection is broken.
func (w *AsyncWriter) processBuffer() bool {
	for {
		select {
		case msg := <-w.logBuffer:
			if _, err := w.conn.Write(msg); err != nil {
				fmt.Printf("[zeronetwork] write failed: %v. Reconnecting...\n", err)
				w.closeConn()
				// put the message back in the buffer, but non-blocking
				// in case the buffer is full (should be rare)
				select {
				case w.logBuffer <- msg:
				default:
				}
				return false // Signal connection is broken
			}
		case <-w.stopChan:
			// Drain remaining buffer before shutting down
			w.drainBuffer()
			return true // Signal shutdown
		}
	}
}

// connectWithRetry tries to connect multiple times.
func (w *AsyncWriter) connectWithRetry() error {
	deadline := time.Now().Add(defaultMaxWait)
	var lastErr error

	for time.Now().Before(deadline) {
		fmt.Println("[zeronetwork] attempting to connect...")
		conn, err := net.DialTimeout(w.cfg.RemoteLogNetwork, w.cfg.RemoteLogAddress, defaultDialTimeout)
		if err == nil {
			w.mu.Lock()
			w.conn = conn
			w.mu.Unlock()
			fmt.Printf("[zeronetwork] connected to %s\n", w.cfg.RemoteLogAddress)
			return nil
		}
		lastErr = err

		select {
		case <-w.stopChan:
			return errors.New("shutdown requested")
		case <-time.After(defaultRetryDelay):
			// continue
		}
	}
	fmt.Printf("[zeronetwork] could not connect in time: %v\n", lastErr)
	return fmt.Errorf("could not connect to remote logger in time: %w", lastErr)
}

// drainBuffer writes any remaining logs in the queue before shutdown.
func (w *AsyncWriter) drainBuffer() {
	if w.conn == nil {
		fmt.Println("[zeronetwork] cannot drain buffer, no connection.")
		return
	}
	close(w.logBuffer) // Close channel to range over remaining items
	fmt.Println("[zeronetwork] draining log buffer before shutdown...")
	for msg := range w.logBuffer {
		if _, err := w.conn.Write(msg); err != nil {
			fmt.Printf("[zeronetwork] failed to write during drain: %v\n", err)
		}
	}
}

// closeConn safely closes the current connection.
func (w *AsyncWriter) closeConn() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
}

// Close gracefully shuts down the writer.
func (w *AsyncWriter) Close() error {
	close(w.stopChan) // Signal worker to stop
	w.wg.Wait()       // Wait for worker to finish
	w.closeConn()
	return nil
}

// ZeroNetworkLogger implements the Logger interface.
type ZeroNetworkLogger struct {
	logger *zerolog.Logger
	writer *AsyncWriter
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
	writer := NewAsyncWriter(cfg, fallback)

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
