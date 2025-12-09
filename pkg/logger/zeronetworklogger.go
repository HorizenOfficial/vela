package logger

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mdlayher/vsock"
	"github.com/rs/zerolog"
)

// TODO
// UDP/TCP split based on log severity
// We would have to implement application-level UDP-like protocol over VSOCK, since VSOCK is stream-oriented and does not support UDP natively.

const (
	defaultBuffer      = 1000
	defaultRetryDelay  = 2 * time.Second
	defaultMaxWait     = 30 * time.Second
	defaultDialTimeout = 5 * time.Second
)

// LogConnectionFactory abstracts the network dialling logic for different protocols.
type LogConnectionFactory interface {
	Dial(timeout time.Duration) (net.Conn, error)
}

// tcpLogConnectionFactory implements LogConnectionFactory for TCP connections.
type tcpLogConnectionFactory struct {
	network string
	address string
}

// NewTCPLogConnectionFactory creates a new TCP log connection factory.
func NewTCPLogConnectionFactory(network, address string) *tcpLogConnectionFactory {
	return &tcpLogConnectionFactory{network: network, address: address}
}

// Dial creates a TCP client connection with a timeout.
func (f *tcpLogConnectionFactory) Dial(timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(f.network, f.address, timeout)
}

// vsockLogConnectionFactory implements LogConnectionFactory for v-sock connections.
type vsockLogConnectionFactory struct {
	cid  uint32
	port uint32
}

// NewVSockLogConnectionFactory creates a new v-sock log connection factory.
func NewVSockLogConnectionFactory(cid, port uint32) *vsockLogConnectionFactory {
	return &vsockLogConnectionFactory{cid: cid, port: port}
}

// Dial establishes a VSOCK connection.
// NOTE: VSOCK has no built-in connect timeout; this call will block until a peer is available.
func (f *vsockLogConnectionFactory) Dial(_ time.Duration) (net.Conn, error) {
	conn, err := vsock.Dial(f.cid, f.port, nil)
	if err != nil {
		fmt.Printf("[zeronetwork] VSOCK connection failed: %v\n", err)
	}
	return conn, err
}

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
	logFactory     LogConnectionFactory
}

// NewAsyncWriter creates a new writer that handles buffering and reconnecting.
func NewAsyncWriter(cfg *Config, fallback io.Writer, factory LogConnectionFactory) *AsyncWriter {
	writer := &AsyncWriter{
		cfg:            cfg,
		fallbackWriter: fallback,
		logBuffer:      make(chan []byte, defaultBuffer),
		stopChan:       make(chan struct{}),
		logFactory:     factory,
	}

	writer.wg.Add(1)
	go writer.run() // Start the async worker

	return writer
}

// Write implements the io.Writer interface. It is non-blocking.
func (w *AsyncWriter) Write(p []byte) (n int, err error) {
	select {
	case <-w.stopChan:
		// writer is closing, drop message safely
		return len(p), nil
	default:
	}

	// Immediately write to fallback for real-time local logging
	if w.conn == nil {

		if _, err := w.fallbackWriter.Write(p); err != nil {
			fmt.Printf("[zeronetwork] Fallback logger failed: %v\n", err)
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
			// Shutdown signal received, buffer will be drained in Close()
			return
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

// processBuffer reads from the buffer and writes to the network.
// Returns false if connection is broken.
func (w *AsyncWriter) processBuffer() bool {
	for {
		select {
		case msg := <-w.logBuffer:
			if msg == nil {
				return true // Channel closed
			}
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
			return true // Shutdown
		default:
			return true // No more messages
		}
	}
}

// connectWithRetry attempts to establish a connection repeatedly until success or max wait time.
func (w *AsyncWriter) connectWithRetry() error {
	deadline := time.Now().Add(defaultMaxWait)
	var lastErr error

	for time.Now().Before(deadline) {
		fmt.Println("[zeronetwork] attempting to connect...")
		conn, err := w.logFactory.Dial(defaultDialTimeout)
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
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn == nil {
		fmt.Println("[zeronetwork] cannot drain buffer, no connection.")
		return
	}
	close(w.logBuffer) // Close channel to range over remaining items
	fmt.Println("[zeronetwork] draining log buffer before shutdown...")
	for {
		select {
		case msg := <-w.logBuffer:
			if msg == nil {
				return
			}
		if _, err := w.conn.Write(msg); err != nil {
			fmt.Printf("[zeronetwork] failed to write during drain: %v\n", err)
			}
		default:
			return // buffer empty
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
	w.closeConn()     // Close active network connection
	w.drainBuffer()   // Drain remaining logs safely
	return nil
}

// ZeroNetworkLogger wraps zerolog with an AsyncWriter.
type ZeroNetworkLogger struct {
	mu     sync.RWMutex
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

	var factory LogConnectionFactory
	switch cfg.RemoteLogNetwork {
	case "tcp":
		factory = NewTCPLogConnectionFactory("tcp", cfg.RemoteLogAddress)
	case "vsock":
		// RemoteLogAddress for vsock should be in "cid:port" format
		parts := splitVsockAddr(cfg.RemoteLogAddress)
		if len(parts) != 2 {
			panic(fmt.Sprintf("Invalid VSock RemoteLogAddress format: %s. Expected 'cid:port'", cfg.RemoteLogAddress))
		}
		cid, err := strconv.ParseUint(parts[0], 10, 32)
		if err != nil {
			panic(fmt.Sprintf("Invalid VSock CID in RemoteLogAddress: %s", parts[0]))
		}
		port, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			panic(fmt.Sprintf("Invalid VSock port in RemoteLogAddress: %s", parts[1]))
		}
		factory = NewVSockLogConnectionFactory(uint32(cid), uint32(port))
	default:
		panic(fmt.Sprintf("Unsupported RemoteLogNetwork: %s", cfg.RemoteLogNetwork))
	}

	writer := NewAsyncWriter(cfg, fallback, factory)

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Caller(). // see init() func above
		Logger()

	if cfg.NetworkLevel != "" {
		logLevel, err := zerolog.ParseLevel(cfg.NetworkLevel)
		if err != nil {
			logLevel = zerolog.InfoLevel
		}
		logger = logger.Level(logLevel)
	}

	return &ZeroNetworkLogger{
		logger: &logger,
		writer: writer,
	}
}

// splitVsockAddr splits a "cid:port" string.
func splitVsockAddr(addr string) []string {
	return strings.Split(strings.TrimSpace(addr), ":")
}

// SetLevel updates the logging level.
func (z *ZeroNetworkLogger) SetLevel(level string) error {
	z.mu.Lock()
	defer z.mu.Unlock()
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	*z.logger = z.logger.Level(lvl)
	return nil
}

// Logging methods just pass through to the underlying zerolog instance
func (z *ZeroNetworkLogger) Trace(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Trace().Msgf(msg, args...)
}
func (z *ZeroNetworkLogger) Debug(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Debug().Msgf(msg, args...)
}
func (z *ZeroNetworkLogger) Info(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Info().Msgf(msg, args...)
}
func (z *ZeroNetworkLogger) Warn(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Warn().Msgf(msg, args...)
}
func (z *ZeroNetworkLogger) Error(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Error().Msgf(msg, args...)
}
func (z *ZeroNetworkLogger) Fatal(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Fatal().Msgf(msg, args...)
}
func (z *ZeroNetworkLogger) Panic(msg string, args ...any) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.logger.Panic().Msgf(msg, args...)
}

// Close safely closes the network connection
func (z *ZeroNetworkLogger) Close() error {
	return z.writer.Close()
}
