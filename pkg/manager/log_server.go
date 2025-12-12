package manager

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

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
	"github.com/mdlayher/vsock"
)

// LogServer handles writing log messages from remote clients to a file.
type LogServer struct {
	logger         logger.Logger
	logFile        *os.File
	fileMutex      sync.Mutex
	consoleEnabled bool
	consoleLevel   string
	fileLevel      string
}

// Static log level priority map
var levelPriority = map[string]int{
	"trace": 0, "debug": 1, "info": 2, "warn": 3, "error": 4, "fatal": 5, "panic": 6,
}

// StartLogServer starts servers to receive log messages from remote clients via TCP and VSOCK.
// Console and file output behavior is controlled via configuration parameters so we can
// enable/disable console logging and choose levels for console/file outputs.
// TODO: have a dynamic levels configuration via setLevel() func etc...
func StartLogServer(
	ctx context.Context,
	tcpAddr common.TcpChannelConnectionParams,
	vsockAddr common.VSockChannelConnectionParams,
	logFilePath string,
	consoleEnabled bool,
	consoleLevel string,
	fileLevel string,
) {
	go func() {
		tcpAddrStr := ""
		if strings.TrimSpace(tcpAddr.Ip) != "" && tcpAddr.Port != 0 {
			tcpAddrStr = tcpAddr.Url()
		}

		// Internal logger for the log server itself, to avoid circular dependencies.
		logServerLogger := logger.NewLogger(&logger.Config{
			Kind:         "zerolog",
			Console:      true,
			ConsoleLevel: "trace",
		})
		defer logServerLogger.Close()

		if strings.TrimSpace(tcpAddrStr) == "" && vsockAddr.Port == 0 {
			logServerLogger.Error("Log server TCP and VSOCK addresses are both empty, not starting log server.")
			panic(fmt.Errorf("log server TCP and VSOCK addresses are both empty, not starting log server"))
		}

		var logFile *os.File
		var err error
		if logFilePath != "" {
			logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logServerLogger.Error("Failed to open log file %s: %v", logFilePath, err)
				panic(err)
			}
			logServerLogger.Info("Remote logs will be written to %s", logFilePath)
		} else {
			logServerLogger.Warn("Log server file path empty, writing logs to console output.")
		}

		logServer := &LogServer{
			logger:         logServerLogger,
			logFile:        logFile,
			consoleEnabled: consoleEnabled,
			consoleLevel:   consoleLevel,
			fileLevel:      fileLevel,
		}

		var wg sync.WaitGroup

		// TCP listener
		if tcpAddrStr != "" {
			tcpListener, err := net.Listen("tcp", tcpAddrStr)
			if err != nil {
				logServerLogger.Error("Failed to start log server on %s: %v", tcpAddrStr, err)
				panic(err)
			} else {
				logServerLogger.Info("Log server listening on TCP %s", tcpAddrStr)
				wg.Add(1)
				go func() {
					defer wg.Done()
					logServer.acceptConnections(ctx, tcpListener)
				}()
			}
		} else {
			logServerLogger.Info("Log server TCP address empty.")
		}

		// VSOCK listener
		if vsockAddr.Port != 0 {
			vsockListener, err := vsock.ListenContextID(vsockAddr.CID, vsockAddr.Port, nil)
			if err != nil {
				logServerLogger.Error("Failed to start log server on vsock %d:%d: %v", vsockAddr.CID, vsockAddr.Port, err)
				panic(err)
			} else {
				logServerLogger.Info("Log server listening on VSOCK %d:%d", vsockAddr.CID, vsockAddr.Port)
				wg.Add(1)
				go func() {
					defer wg.Done()
					logServer.acceptConnections(ctx, vsockListener)
				}()
			}
		} else {
			logServerLogger.Info("Log server VSOCK address empty.")
		}

		go func() {
			<-ctx.Done()
			logServerLogger.Info("Shutting down log server...")
			if logServer.logFile != nil {
				if err := logServer.logFile.Close(); err != nil {
					logServerLogger.Error("Error closing log file: %v", err)
				}
			}
		}()

		wg.Wait()
	}()
}

func (ls *LogServer) acceptConnections(ctx context.Context, listener net.Listener) {
	go func() {
		<-ctx.Done()
		listener.Close()
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
		go ls.handleLogConnection(conn)
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
			if err != io.EOF {
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
	ls.logger.Info("Log connection from %s closed", conn.RemoteAddr())
}

// filterAndWrite routes an incoming JSON log entry to console and/or file depending
// on configured log levels. Entries missing or with unknown levels default to "info".
// File writes are protected by fileMutex to ensure thread-safe concurrent access.
func (ls *LogServer) filterAndWrite(jsonMsg []byte) error {
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

	consolePriority := levelPriority[strings.ToLower(ls.consoleLevel)]
	filePriority := levelPriority[strings.ToLower(ls.fileLevel)]

	if ls.consoleEnabled && entryPriority >= consolePriority {
		fmt.Println(string(jsonMsg))
	}

	if ls.logFile != nil && entryPriority >= filePriority {
		ls.fileMutex.Lock()
		defer ls.fileMutex.Unlock()
		if _, err := ls.logFile.Write(append(jsonMsg, '\n')); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
