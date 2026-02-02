package logserver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
	"github.com/mdlayher/vsock"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Struct to hold console and file log levels atomically
type logLevels struct {
	console int
	file    int
}

// LogServer handles writing log messages from remote clients to a file.
type LogServer struct {
	logger         logger.Logger
	logWriter      io.WriteCloser // can be *os.File or *lumberjack.Logger
	fileMutex      sync.Mutex
	consoleEnabled bool
	levels         atomic.Value // stores logLevels
}

// LogServerConfig holds configuration for starting the log server.
type LogServerConfig struct {
	TCPAddr        common.TcpChannelConnectionParams
	VSockAddr      common.VSockChannelConnectionParams
	LogFilePath    string
	ConsoleEnabled bool
	ConsoleLevel   string
	FileLevel      string
	// Log rotation settings (only used when LogFilePath is set)
	RotationEnabled bool // Enable log rotation using lumberjack
	MaxSizeMB       int  // Max size in megabytes before rotation (default: 100)
	MaxBackups      int  // Max number of old log files to retain (default: 3)
	MaxAgeDays      int  // Max days to retain old log files, 0 = no limit (default: 28)
	Compress        bool // Compress rotated log files with gzip (default: true)
}

// Static log level priority map
var levelPriority = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"fatal": 5,
	"panic": 6,
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

	// Validate log levels
	if cfg.ConsoleEnabled {
		if _, ok := levelPriority[strings.ToLower(cfg.ConsoleLevel)]; !ok {
			return fmt.Errorf("invalid console log level: %s", cfg.ConsoleLevel)
		}
	}
	if cfg.LogFilePath != "" {
		if _, ok := levelPriority[strings.ToLower(cfg.FileLevel)]; !ok {
			return fmt.Errorf("invalid file log level: %s", cfg.FileLevel)
		}
	}

	var logWriter io.WriteCloser
	if cfg.LogFilePath != "" {
		if cfg.RotationEnabled {
			// Use lumberjack for log rotation
			maxSize := cfg.MaxSizeMB
			if maxSize <= 0 {
				maxSize = 100 // default 100MB
			}
			maxBackups := cfg.MaxBackups
			if maxBackups <= 0 {
				maxBackups = 3 // default 3 backups
			}
			maxAge := cfg.MaxAgeDays
			if maxAge < 0 {
				maxAge = 28 // default 28 days
			}
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

	// Initialize atomic logLevels
	logServer.levels.Store(logLevels{
		console: levelPriority[strings.ToLower(cfg.ConsoleLevel)],
		file:    levelPriority[strings.ToLower(cfg.FileLevel)],
	})

	if logWriter != nil {
		if cfg.RotationEnabled {
			logServer.logger.Info("Remote logs will be written to file %s with level [%s] (rotation enabled: maxSize=%dMB, maxBackups=%d, maxAge=%d days, compress=%v)",
				cfg.LogFilePath, cfg.FileLevel, cfg.MaxSizeMB, cfg.MaxBackups, cfg.MaxAgeDays, cfg.Compress)
		} else {
			logServer.logger.Info("Remote logs will be written to file %s with level [%s] (rotation disabled)", cfg.LogFilePath, cfg.FileLevel)
		}
	}
	if cfg.ConsoleEnabled {
		if logWriter == nil {
			logServer.logger.Warn("No log file configured, logs will only be written to console")
		}
		logServer.logger.Info("Remote logs will be written to console with level [%s]", cfg.ConsoleLevel)
	}

	go logServer.run(ctx, tcpAddrStr, cfg.VSockAddr)

	return nil
}

// UpdateLogLevels updates console and file log levels at runtime
func (ls *LogServer) UpdateLogLevels(consoleLevel, fileLevel string) error {
	console, ok := levelPriority[strings.ToLower(consoleLevel)]
	if !ok {
		return fmt.Errorf("invalid console log level: %s", consoleLevel)
	}
	file, ok := levelPriority[strings.ToLower(fileLevel)]
	if !ok {
		return fmt.Errorf("invalid file log level: %s", fileLevel)
	}

	ls.levels.Store(logLevels{
		console: console,
		file:    file,
	})

	ls.logger.Info("Log levels updated at runtime (console=%s, file=%s)", consoleLevel, fileLevel)
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
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				ls.logger.Info("Log connection from %s closed gracefully", conn.RemoteAddr())
			} else {
				ls.logger.Error("Error reading from log connection %s: %v", conn.RemoteAddr(), err)
			}
			break
		}

		trimmed := bytes.TrimSpace([]byte(message))
		if len(trimmed) == 0 {
			continue
		}

		// Apply per-level filtering for console/file outputs.
		if err := ls.filterAndWrite(trimmed); err != nil {
			ls.logger.Error("Error processing remote log: %v", err)
		}
	}
}

// filterAndWrite routes an incoming JSON log entry to console and/or file depending on configured log levels
func (ls *LogServer) filterAndWrite(jsonMsg []byte) error {
	// TODO: do something faster than unmarshalling for optimize this
	var entry struct {
		Level string `json:"level"`
	}
	if err := json.Unmarshal(jsonMsg, &entry); err != nil {
		return fmt.Errorf("invalid JSON log: %w", err)
	}

	msgLevel := strings.ToLower(strings.TrimSpace(entry.Level))
	entryPriority, ok := levelPriority[msgLevel]
	if !ok {
		entryPriority = levelPriority["info"]
	}

	levels := ls.levels.Load().(logLevels)

	if ls.consoleEnabled && entryPriority >= levels.console {
		fmt.Println(string(jsonMsg))
	}

	ls.fileMutex.Lock()
	defer ls.fileMutex.Unlock()

	if ls.logWriter != nil && entryPriority >= levels.file {
		if _, err := ls.logWriter.Write(append(jsonMsg, '\n')); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
