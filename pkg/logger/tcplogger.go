package logger

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// TCPLogger implements the Logger interface for sending logs over TCP.
type TCPLogger struct {
	conn           net.Conn
	fallbackLogger Logger // Used for internal errors of the TCPLogger
	logLevel       map[string]int
	minLevel       int
}

// NewTCPLogger creates a new TCPLogger instance.
func NewTCPLogger(cfg *Config) *TCPLogger {
	if cfg.RemoteLogAddress == "" {
		panic("RemoteLogAddress cannot be empty for tcplog kind")
	}

	conn, err := net.Dial("tcp", cfg.RemoteLogAddress)
	if err != nil {
		// If TCP connection fails, panic as it's a critical setup error for this logger type
		panic(fmt.Sprintf("failed to connect to remote TCP logger at %s: %v", cfg.RemoteLogAddress, err))
	}

	// Create a fallback logger for internal errors of the TCPLogger itself
	fallback := NewPrintfLogger(&Config{Console: true, ConsoleLevel: "error"})

	logLevel := map[string]int{
		"trace": 0,
		"debug": 1,
		"info":  2,
		"warn":  3,
		"error": 4,
		"fatal": 5,
		"panic": 6,
	}

	minLevel := logLevel["info"] // Default minimum level
	if level, ok := logLevel[cfg.ConsoleLevel]; ok {
		minLevel = level
	}

	return &TCPLogger{
		conn:           conn,
		fallbackLogger: fallback,
		logLevel:       logLevel,
		minLevel:       minLevel,
	}
}

func (t *TCPLogger) sendLog(level string, msg string, args ...any) {
	if t.logLevel[level] < t.minLevel {
		return
	}

	formattedMsg := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, fmt.Sprintf(msg, args...))
	_, err := io.WriteString(t.conn, formattedMsg)
	if err != nil {
		t.fallbackLogger.Error("failed to send remote log via TCP: %v", err)
	}
}

func (t *TCPLogger) Trace(msg string, args ...any) { t.sendLog("TRACE", msg, args...) }
func (t *TCPLogger) Debug(msg string, args ...any) { t.sendLog("DEBUG", msg, args...) }
func (t *TCPLogger) Info(msg string, args ...any)  { t.sendLog("INFO", msg, args...) }
func (t *TCPLogger) Warn(msg string, args ...any)  { t.sendLog("WARN", msg, args...) }
func (t *TCPLogger) Error(msg string, args ...any) { t.sendLog("ERROR", msg, args...) }
func (t *TCPLogger) Fatal(msg string, args ...any) {
	t.sendLog("FATAL", msg, args...)
	log.Fatalf(msg, args...) // Also cause local application to exit
}
func (t *TCPLogger) Panic(msg string, args ...any) {
	t.sendLog("PANIC", msg, args...)
	log.Panicf(msg, args...) // Also cause local application to panic
}

// Close closes the TCP connection.
func (t *TCPLogger) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
