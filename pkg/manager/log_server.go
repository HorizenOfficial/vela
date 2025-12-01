package manager

import (
	"bufio"
	"context"

	"io"
	"net"
	"strings"

	"github.com/horizen-pes/pkg/logger"
)

// StartLogServer starts a TCP server to receive log messages from remote clients.
func StartLogServer(ctx context.Context, listenAddr string, appLogger logger.Logger) {
	if strings.TrimSpace(listenAddr) == "" {
		appLogger.Info("Log server address is empty, not starting log server.")
		return
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		appLogger.Error("Failed to start log server on %s: %v", listenAddr, err)
		return
	}
	appLogger.Info("Log server listening on %s", listenAddr)

	go func() {
		<-ctx.Done() // Wait for context cancellation
		appLogger.Info("Shutting down log server...")
		if err := listener.Close(); err != nil {
			appLogger.Error("Error closing log server listener: %v", err)
		}
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return // Context canceled, listener is likely closed
				default:
					appLogger.Error("Log server accept error: %v", err)
				}
			continue
			}

			go handleLogConnection(conn, appLogger)
		}
	}()
}

func handleLogConnection(conn net.Conn, appLogger logger.Logger) {
	defer func() {
		if err := conn.Close(); err != nil {
			appLogger.Error("Error closing log connection: %v", err)
		}
	}()

	appLogger.Info("New log connection from %s", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				appLogger.Error("Error reading from log connection %s: %v", conn.RemoteAddr(), err)
			}
			break
		}
		// Log the received message using the manager's logger
		appLogger.Info("Remote Log: %s", strings.TrimSpace(message))
	}
	appLogger.Info("Log connection from %s closed", conn.RemoteAddr())
}
