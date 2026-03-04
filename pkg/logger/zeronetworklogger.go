package logger

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/mdlayher/vsock"
	"github.com/rs/zerolog"
)

// We could consider implementing a UDP/TCP split based on log severity, but that would be feasible
// only for tcp connections, but we would have to implement application-level UDP-like
// protocol over VSOCK, since VSOCK is stream-oriented and does not support UDP natively.

const (
	defaultBuffer        = 1000
	defaultRetryDelay    = 2 * time.Second
	defaultMaxWait       = 30 * time.Second
	defaultDialTimeout   = 5 * time.Second
	defaultShutdownGrace = 200 * time.Millisecond
)

// LogConnectionFactory abstracts the network dialling logic for different protocols.
type LogConnectionFactory interface {
	Dial(timeout time.Duration) (net.Conn, error)
	GetAddress() string
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

func (f *tcpLogConnectionFactory) GetAddress() string {
	return fmt.Sprintf("%s:%s", f.network, f.address)
}

// vsockLogConnectionFactory implements LogConnectionFactory for v-sock connections.
type vsockLogConnectionFactory struct {
	cid  uint32
	port uint32
}

func (f *vsockLogConnectionFactory) GetAddress() string {
	return fmt.Sprintf("%d:%d", f.cid, f.port)
}

// NewVSockLogConnectionFactory creates a new v-sock log connection factory.
func NewVSockLogConnectionFactory(cid, port uint32) *vsockLogConnectionFactory {
	return &vsockLogConnectionFactory{cid: cid, port: port}
}

// Dial establishes a VSOCK connection.
// NOTE: VSOCK has no built-in connect timeout; this call will block until a peer is available.
func (f *vsockLogConnectionFactory) Dial(_ time.Duration) (net.Conn, error) {
	return vsock.Dial(f.cid, f.port, nil)
}

// AsyncWriter is a resilient, non-blocking io.Writer.
// It writes to a fallback writer immediately and buffers logs for an async worker.
type AsyncWriter struct {
	mu             sync.RWMutex
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

// getConn provides a thread-safe snapshot of the current connection.
func (w *AsyncWriter) getConn() net.Conn {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.conn
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
	if w.getConn() == nil {
		// we also store it in the buffer too, but if it gets full we lose it completely
		b := append([]byte("[zeronetwork conn down] "), p...)
		if _, err := w.fallbackWriter.Write(b); err != nil {
			fmt.Printf("[zeronetwork] Fallback logger failed: %v\n", err)
		}
	}

	// Create a copy, as the slice is reused by zerolog
	buf := make([]byte, len(p))
	copy(buf, p)

	// Send to buffer without blocking
	w.sendToLogBuffer(buf)

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
			// attempt to connect until success or shutdown
			if w.getConn() == nil {
				if err := w.connectWithRetry(); err != nil {
					// Sleep before retrying connection
					time.Sleep(defaultRetryDelay)
					continue
				}
			}

			// Connection is live, start processing buffer.
			w.processBuffer()
		}
	}
}

// processBuffer reads from the buffer and writes to the network.
// Returns false if connection is broken.
func (w *AsyncWriter) processBuffer() {
	for {
		select {
		// blocks until a msg is available or the chan is closed
		case msg, ok := <-w.logBuffer:
			if !ok {
				return // channel closed
			}
			// Capture a snapshot of the connection to prevent race during Write
			conn := w.getConn()
			if conn == nil {
				// put the msg back in the buffer
				w.sendToLogBuffer(msg)
				return
			}

			if _, err := conn.Write(msg); err != nil {
				errMsg := fmt.Sprintf("[zeronetwork] write failed: %v. Reconnecting...\n", err)
				w.fallbackWriter.Write([]byte(errMsg))
				w.closeConn()
				// put the message back in the buffer, but non-blocking
				// in case the buffer is full (should be rare)
				w.sendToLogBuffer(msg)
				return // Signal connection is broken
			}
			// go on with the for loop and consume any other msg if any
			//return true
		case <-w.stopChan:
			return // Shutdown
		}
	}
}

// sendToLogBuffer attempts to put a message in the writer log buffer
func (w *AsyncWriter) sendToLogBuffer(msg []byte) {
	select {
	case w.logBuffer <- msg:
	default:
		// If buffer is full, we must drop it to prevent deadlock
		w.fallbackWriter.Write([]byte("Remote log buffer is full. Dropping message:\n"))
	}
}

// connectWithRetry attempts to establish a connection repeatedly until success or max wait time.
func (w *AsyncWriter) connectWithRetry() error {
	deadline := time.Now().Add(defaultMaxWait)
	var lastErr error

	for time.Now().Before(deadline) {
		w.fallbackWriter.Write(fmt.Appendf(nil, "[zeronetwork] attempting to connect to %s...\n", w.logFactory.GetAddress()))
		conn, err := w.logFactory.Dial(defaultDialTimeout)
		if err == nil {
			w.mu.Lock()
			w.conn = conn
			w.mu.Unlock()
			logMsg := fmt.Sprintf("[zeronetwork] connected on %s\n", w.cfg.RemoteLogNetwork)
			w.fallbackWriter.Write([]byte(logMsg))
			return nil
		} else {
			errMsg := fmt.Sprintf("[zeronetwork] connection failed: %v\n", err)
			w.fallbackWriter.Write([]byte(errMsg))
		}
		lastErr = err

		select {
		case <-w.stopChan:
			return errors.New("shutdown requested")
		case <-time.After(defaultRetryDelay):
			// continue
		}
	}
	errMsg := fmt.Sprintf("[zeronetwork] could not connect in time: %v\n", lastErr)
	w.fallbackWriter.Write([]byte(errMsg))
	return fmt.Errorf("could not connect to remote logger in time: %w", lastErr)
}

// drainBuffer writes any remaining logs in the queue before shutdown.
// Messages are written to both the network connection and the fallback
// (console) writer to guarantee visibility even if the log server shuts
// down before it processes the TCP data.
func (w *AsyncWriter) drainBuffer() {
	conn := w.getConn()
	if conn == nil {
		w.fallbackWriter.Write([]byte("[zeronetwork] cannot drain buffer, no connection.\n"))
		return
	}
	// do not close logBuffer, Write() might still attempt to send into it
	//close(w.logBuffer)
	w.fallbackWriter.Write([]byte("[zeronetwork] draining log buffer before shutdown...\n"))
	for {
		select {
		case msg := <-w.logBuffer:
			if _, err := conn.Write(msg); err != nil {
				errMsg := fmt.Sprintf("[zeronetwork] failed to write during drain: %v\n", err)
				w.fallbackWriter.Write([]byte(errMsg))
			}
			// Also write to fallback so the message is visible on the
			// console even if the log server is torn down before it
			// reads these bytes from the TCP buffer.
			w.fallbackWriter.Write(msg)
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
	w.drainBuffer()   // Drain remaining logs safely

	// Brief grace period: the worker goroutine may have already written
	// messages to the TCP socket before stopping. The log server needs a
	// moment to read those bytes from its receive buffer and write them
	// to file/console before we close the connection and the caller tears
	// down the server via context cancellation.
	time.Sleep(defaultShutdownGrace)

	w.closeConn() // Close active network connection
	return nil
}

// ZeroNetworkLogger wraps zerolog with an AsyncWriter.
type ZeroNetworkLogger struct {
	mu     sync.RWMutex
	logger *zerolog.Logger
	writer *AsyncWriter
}

type TimestampWriter struct {
	out io.Writer
}

func (tw *TimestampWriter) Write(p []byte) (int, error) {
	stamp := time.Now().Format(TimeStampFormatMs)
	line := fmt.Sprintf("%s %s", stamp, p)
	return tw.out.Write([]byte(line))
}

// NewZeroNetworkLogger creates a new logger instance
func NewZeroNetworkLogger(cfg *Config) *ZeroNetworkLogger {
	if cfg.RemoteLogNetwork == "" {
		cfg.RemoteLogNetwork = "tcp"
	}

	//	fallback := NewPrintfLogger(&Config{Console: true, ConsoleLevel: "error"})
	//fallback := os.Stdout
	fallback := &TimestampWriter{os.Stdout}

	var factory LogConnectionFactory
	switch cfg.RemoteLogNetwork {
	case "tcp":
		params, ok := cfg.RemoteLogParams.(common.TcpChannelConnectionParams)
		if !ok {
			panic(fmt.Sprintf("Invalid RemoteLogParams for tcp: %T", cfg.RemoteLogParams))
		}
		factory = NewTCPLogConnectionFactory("tcp", params.Url())
	case "vsock":
		params, ok := cfg.RemoteLogParams.(common.VSockChannelConnectionParams)
		if !ok {
			panic(fmt.Sprintf("Invalid RemoteLogParams for vsock: %T", cfg.RemoteLogParams))
		}
		factory = NewVSockLogConnectionFactory(params.CID, params.Port)
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

func (z *ZeroNetworkLogger) SetLevel(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}

	z.mu.Lock()
	defer z.mu.Unlock()

	// Create a new branched logger instance
	newLogger := z.logger.Level(lvl)
	// Update the pointer to the new instance
	z.logger = &newLogger

	return nil
}

func (z *ZeroNetworkLogger) GetLevel() string {
	z.mu.RLock()
	defer z.mu.RUnlock()
	return z.logger.GetLevel().String()
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
