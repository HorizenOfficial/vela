package logserver

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
	"github.com/mdlayher/vsock"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogServer handles writing log messages from remote clients to a file and/or console.
type LogServer struct {
	logger         logger.Logger
	logWriter      io.WriteCloser // can be *os.File or *lumberjack.Logger
	fileMutex      sync.Mutex
	consoleEnabled bool
}

// LogServerConfig holds configuration for starting the log server.
type LogServerConfig struct {
	TCPAddr        common.TcpChannelConnectionParams
	VSockAddr      common.VSockChannelConnectionParams
	LogFilePath    string
	ConsoleEnabled bool
	// Log rotation settings (only used when LogFilePath is set)
	RotationEnabled bool // Enable log rotation using lumberjack
	MaxSizeMB       int  // Max size in megabytes before rotation (default if negative or zero: 100)
	MaxBackups      int  // Max number of old log files to retain, 0 = retain all (default if negative: 3)
	MaxAgeDays      int  // Max days to retain old log files, 0 = no limit (default if negative: 28)
	Compress        bool // Compress rotated log files with gzip (default: false)
}

// StartLogServer starts server to receive log messages from remote clients via TCP and VSOCK.
// It performs synchronous validation of inputs and file access before launching background routines.
func StartLogServer(ctx context.Context, cfg LogServerConfig) error {
	tcpAddrStr := ""
	if strings.TrimSpace(cfg.TCPAddr.Ip) != "" && cfg.TCPAddr.Port != 0 {
		tcpAddrStr = cfg.TCPAddr.Url()
	}

	if strings.TrimSpace(tcpAddrStr) == "" && cfg.VSockAddr.Port == 0 {
		return fmt.Errorf("log server TCP and VSOCK addresses are both empty, not starting log server")
	}

	if cfg.LogFilePath == "" && !cfg.ConsoleEnabled {
		return fmt.Errorf("log server file path is empty and console output is disabled, not starting log server")
	}

	var logWriter io.WriteCloser
	// Variables for rotation settings (used in log message later)
	var maxSize, maxBackups, maxAge int
	if cfg.LogFilePath != "" {
		if cfg.RotationEnabled {
			// Use lumberjack for log rotation
			maxSize = cfg.MaxSizeMB
			if maxSize <= 0 { // <= 0 (not < 0): lumberjack also treats 0 as "use default 100MB", so there's no useful zero-value to preserve
				maxSize = 100 // default 100MB
			}
			maxBackups = cfg.MaxBackups
			if maxBackups < 0 {
				maxBackups = 3 // default 3 backups
			}
			maxAge = cfg.MaxAgeDays
			if maxAge < 0 {
				maxAge = 28 // default 28 days
			}
			// Validate upper bounds to prevent misconfiguration from exhausting disk space
			if maxSize > 1024 {
				return fmt.Errorf("MaxSizeMB=%d exceeds maximum allowed value of 1024 (1GB)", maxSize)
			}
			if maxBackups > 100 {
				return fmt.Errorf("MaxBackups=%d exceeds maximum allowed value of 100", maxBackups)
			}
			if maxAge > 365 {
				return fmt.Errorf("MaxAgeDays=%d exceeds maximum allowed value of 365", maxAge)
			}
			// Validate file access before configuring lumberjack (which defers file creation to first Write)
			f, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to open log file %s: %w", cfg.LogFilePath, err)
			}
			f.Close()
			logWriter = &lumberjack.Logger{
				Filename:   cfg.LogFilePath,
				MaxSize:    maxSize,
				MaxBackups: maxBackups,
				MaxAge:     maxAge,
				Compress:   cfg.Compress,
			}
		} else {
			// Use plain file without rotation
			var err error
			logWriter, err = os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to open log file %s: %w", cfg.LogFilePath, err)
			}
		}
	}

	logServerLogger := logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      true,
		ConsoleLevel: "trace",
	})

	logServer := &LogServer{
		logger:         logServerLogger,
		logWriter:      logWriter,
		consoleEnabled: cfg.ConsoleEnabled,
	}

	if logWriter != nil {
		if cfg.RotationEnabled {
			logServer.logger.Info("Remote logs will be written to file %s (rotation enabled: maxSize=%dMB, maxBackups=%d, maxAge=%d days, compress=%v)",
				cfg.LogFilePath, maxSize, maxBackups, maxAge, cfg.Compress)
		} else {
			logServer.logger.Info("Remote logs will be written to file %s (rotation disabled)", cfg.LogFilePath)
		}
	}
	if cfg.ConsoleEnabled {
		if logWriter == nil {
			logServer.logger.Warn("No log file configured, logs will only be written to console")
		}
		logServer.logger.Info("Remote logs will be written to console")
	}

	go logServer.run(ctx, tcpAddrStr, cfg.VSockAddr)

	return nil
}

// run handles the lifecycle of the listeners and shutdown
func (ls *LogServer) run(ctx context.Context, tcpAddrStr string, vsockAddr common.VSockChannelConnectionParams) {
	// Ensure logger is closed when the server stops completely
	defer ls.logger.Close()

	var wg sync.WaitGroup

	// TCP listener
	if tcpAddrStr != "" {
		tcpListener, err := net.Listen("tcp", tcpAddrStr)
		if err != nil {
			ls.logger.Error("Failed to start log server on %s: %v", tcpAddrStr, err)
			// We don't panic here to allow the other listener (VSOCK) to potentially work
			// TODO: or we could decide to shutdown all
		} else {
			ls.logger.Info("Log server listening on TCP %s", tcpAddrStr)
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Pass &wg to track individual connections
				ls.acceptConnections(ctx, tcpListener, &wg)
			}()
		}
	}

	// VSOCK listener
	if vsockAddr.Port != 0 {
		vsockListener, err := vsock.Listen(vsockAddr.Port, nil)
		if err != nil {
			ls.logger.Error("Failed to start log server on vsock port:%d: %v", vsockAddr.Port, err)
		} else {
			ls.logger.Info("Log server listening on %s %s", vsockListener.Addr().Network(), vsockListener.Addr().String())
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Pass &wg to track individual connections
				ls.acceptConnections(ctx, vsockListener, &wg)
			}()
		}
	}

	// shutdown handler is part of lifecycle
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		ls.logger.Info("Shutting down log server...")

		ls.fileMutex.Lock()
		defer ls.fileMutex.Unlock()

		if ls.logWriter != nil {
			if err := ls.logWriter.Close(); err != nil {
				ls.logger.Error("Error closing log writer: %v", err)
			}
			ls.logWriter = nil
		}
	}()

	wg.Wait()
}

func (ls *LogServer) acceptConnections(ctx context.Context, listener net.Listener, wg *sync.WaitGroup) {
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return // Context canceled, listener is likely closed
			default:
				ls.logger.Error("Log server accept error: %v", err)
			}
			continue
		}

		// Track individual connection so we don't exit while handling a log
		wg.Add(1)
		go func() {
			defer wg.Done()
			ls.handleLogConnection(conn)
		}()
	}
}

func (ls *LogServer) handleLogConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			ls.logger.Error("Error closing log connection: %v", err)
		}
	}()

	ls.logger.Info("New log connection from %s", conn.RemoteAddr())

	// handle oversized lines gracefully. ReadLine() returns isPrefix=true when a line
	// exceeds the buffer, allowing us to truncate it to the first chunk and continue reading.
	// Buffer size of 64KB is generous for structured zerolog JSON entries (typically 150-300 bytes).
	const maxLogLineSize = 64 * 1024 // 64KB
	reader := bufio.NewReaderSize(conn, maxLogLineSize)

	for {
		lineBytes, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				ls.logger.Info("Log connection from %s closed gracefully", conn.RemoteAddr())
			} else {
				ls.logger.Error("Error reading from log connection %s: %v", conn.RemoteAddr(), err)
			}
			break
		}

		// Check isPrefix before TrimSpace to ensure oversized lines are handled
		// even if their first chunk is all whitespace (otherwise we'd continue past
		// the truncation block, leaving the rest of the oversized line in the buffer).
		if isPrefix {
			ls.logger.Warn("Oversized log line from %s (>%d bytes), truncating", conn.RemoteAddr(), maxLogLineSize)
			// Copy the first chunk before draining — lineBytes references the reader's
			// internal buffer which gets overwritten by subsequent ReadLine calls.
			truncated := make([]byte, len(lineBytes))
			copy(truncated, lineBytes)
			// Drain remaining chunks of the oversized line
			for isPrefix {
				_, isPrefix, err = reader.ReadLine()
				if err != nil {
					break
				}
			}
			// Re-emit as a valid zerolog entry with the truncated content as the message
			var buf bytes.Buffer
			zl := zerolog.New(&buf).With().Timestamp().Logger()
			zl.Warn().Bool("_truncated", true).Msg(string(truncated))
			lineBytes = buf.Bytes()
		}

		trimmed := bytes.TrimSpace(lineBytes)
		if len(trimmed) == 0 {
			continue
		}

		// Write to all configured outputs (no filtering — level control is at the source logger).
		if err := ls.writeToOutputs(trimmed); err != nil {
			ls.logger.Error("Error processing remote log: %v", err)
		}
	}
}

// writeToOutputs writes an incoming log entry to console and/or file unconditionally.
func (ls *LogServer) writeToOutputs(jsonMsg []byte) error {
	if ls.consoleEnabled {
		fmt.Println(string(jsonMsg))
	}

	ls.fileMutex.Lock()
	defer ls.fileMutex.Unlock()

	if ls.logWriter != nil {
		if _, err := ls.logWriter.Write(append(jsonMsg, '\n')); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
