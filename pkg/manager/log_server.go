package manager

import (
	"bufio"
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
)

// LogServer handles writing log messages from remote clients to a file.
type LogServer struct {
	logger    logger.Logger
	logFile   *os.File
	fileMutex sync.Mutex
}

// StartLogServer starts servers to receive log messages from remote clients via TCP and VSOCK.
func StartLogServer(ctx context.Context, tcpAddr common.TcpChannelConnectionParams, vsockAddr, logFilePath string) {
	go func() {
		tcpAddrStr := ""
		if strings.TrimSpace(tcpAddr.Ip) != "" && tcpAddr.Port != 0 {
			tcpAddrStr = tcpAddr.Url()
		}

		// Internal logger for the log server itself, to avoid circular dependencies.
		logServerLogger := logger.NewLogger(&logger.Config{
			Kind:         "zerolog",
			Console:      true,
			ConsoleLevel: "info",
		})
		defer logServerLogger.Close()

		if strings.TrimSpace(tcpAddrStr) == "" && strings.TrimSpace(vsockAddr) == "" {
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
			logger:  logServerLogger,
			logFile: logFile,
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
		if vsockAddr != "" {
			cid, port, err := parseVSockAddr(vsockAddr)
			if err != nil {
				logServerLogger.Error("Invalid VSOCK address: %v", err)
				panic(err)
			} else {
				vsockListener, err := vsock.ListenContextID(cid, port, nil)
				if err != nil {
					logServerLogger.Error("Failed to start log server on vsock %s: %v", vsockAddr, err)
					panic(err)
				} else {
					logServerLogger.Info("Log server listening on VSOCK %s", vsockAddr)
					wg.Add(1)
					go func() {
						defer wg.Done()
						logServer.acceptConnections(ctx, vsockListener)
					}()
				}
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

		// Print to console
		// TODO: we always log on console, shall we render it configurable? Or as an alternative to using a file?
		fmt.Print(message)

		// Write to file if configured
		if ls.logFile != nil {
			ls.fileMutex.Lock()
			if _, err := ls.logFile.WriteString(message); err != nil {
				ls.logger.Error("Error writing to log file: %v", err)
			}
			ls.fileMutex.Unlock()
		}
	}
	ls.logger.Info("Log connection from %s closed", conn.RemoteAddr())
}

func parseVSockAddr(addr string) (uint32, uint32, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid vsock address format")
	}
	cid, err := Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid vsock cid: %v", err)
	}
	port, err := Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid vsock port: %v", err)
	}
	return uint32(cid), uint32(port), nil
}

func Atoi(s string) (int, error) {
	var n int
	_, err := fmt.Sscan(s, &n)
	return n, err
}
