package logserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/common"
)

// sendLogMessage sends a JSON log message to the log server via TCP.
// Creates a new connection for each call - use writeLogEntry for bulk writes.
func sendLogMessage(t *testing.T, addr string, level, message string) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to log server: %v", err)
	}
	defer conn.Close()

	writeLogEntry(t, conn, level, message)
}

// writeLogEntry writes a JSON log entry to an existing connection.
// Use this for bulk writes to avoid connection overhead.
func writeLogEntry(t *testing.T, conn net.Conn, level, message string) {
	logEntry := map[string]string{
		"level":   level,
		"time":    time.Now().Format("2006-Jan-02 15:04:05.000"),
		"message": message,
	}
	data, err := json.Marshal(logEntry)
	if err != nil {
		t.Fatalf("Failed to marshal log entry: %v", err)
	}
	data = append(data, '\n')

	_, err = conn.Write(data)
	if err != nil {
		t.Fatalf("Failed to write log message: %v", err)
	}
}

// waitForFile waits for a file to exist and have content
func waitForFile(t *testing.T, path string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if info, err := os.Stat(path); err == nil && info.Size() > 0 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("File %s did not appear or remained empty within timeout", path)
}

// TestLogServer_ValidationErrors verifies that StartLogServer returns appropriate errors
// for invalid configurations (table-driven).
func TestLogServer_ValidationErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	tests := []struct {
		name      string
		cfg       LogServerConfig
		wantError string
	}{
		{
			name:      "both addresses empty",
			cfg:       LogServerConfig{LogFilePath: logFile},
			wantError: "TCP and VSOCK addresses are both empty",
		},
		{
			name: "no file path and console disabled",
			cfg: LogServerConfig{
				TCPAddr:        common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: 9999},
				ConsoleEnabled: false,
			},
			wantError: "file path is empty and console output is disabled",
		},
		{
			name: "MaxSizeMB exceeds upper bound",
			cfg: LogServerConfig{
				TCPAddr:         common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: 9999},
				LogFilePath:     logFile,
				RotationEnabled: true,
				MaxSizeMB:       2000,
			},
			wantError: "MaxSizeMB=2000 exceeds maximum allowed value of 1024",
		},
		{
			name: "MaxBackups exceeds upper bound",
			cfg: LogServerConfig{
				TCPAddr:         common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: 9999},
				LogFilePath:     logFile,
				RotationEnabled: true,
				MaxSizeMB:       10,
				MaxBackups:      200,
			},
			wantError: "MaxBackups=200 exceeds maximum allowed value of 100",
		},
		{
			name: "MaxAgeDays exceeds upper bound",
			cfg: LogServerConfig{
				TCPAddr:         common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: 9999},
				LogFilePath:     logFile,
				RotationEnabled: true,
				MaxSizeMB:       10,
				MaxBackups:      3,
				MaxAgeDays:      500,
			},
			wantError: "MaxAgeDays=500 exceeds maximum allowed value of 365",
		},
		{
			name: "invalid file path with rotation enabled",
			cfg: LogServerConfig{
				TCPAddr:         common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: 9999},
				LogFilePath:     "/nonexistent/path/test.log",
				RotationEnabled: true,
				MaxSizeMB:       10,
			},
			wantError: "failed to open log file",
		},
		{
			name: "invalid file path without rotation",
			cfg: LogServerConfig{
				TCPAddr:     common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: 9999},
				LogFilePath: "/nonexistent/path/test.log",
			},
			wantError: "failed to open log file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := StartLogServer(ctx, tc.cfg)
			if err == nil {
				t.Fatal("Expected error but got nil")
			}
			if !strings.Contains(err.Error(), tc.wantError) {
				t.Errorf("Expected error containing %q, got: %v", tc.wantError, err)
			}
		})
	}
}

// TestLogServer_OversizedLineTruncated verifies that a log line exceeding the 64KB buffer
// is truncated (first chunk written) without killing the connection, and subsequent normal-sized lines are still processed.
func TestLogServer_OversizedLineTruncated(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send a normal message first
	writeLogEntry(t, conn, "info", "before oversized")

	// Send an oversized line (>64KB without newline, then newline)
	oversized := `{"level":"info","message":"` + strings.Repeat("x", 70*1024) + `"}` + "\n"
	_, err = conn.Write([]byte(oversized))
	if err != nil {
		t.Fatalf("Failed to write oversized message: %v", err)
	}

	// Send another normal message after the oversized one
	writeLogEntry(t, conn, "info", "after oversized")

	// Wait for messages to be processed
	waitForFile(t, logFile, 2*time.Second)
	time.Sleep(300 * time.Millisecond)

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)

	// The message before the oversized line should be present
	if !strings.Contains(contentStr, "before oversized") {
		t.Error("Log file missing 'before oversized' message")
	}

	// The message after the oversized line should also be present (connection survived)
	if !strings.Contains(contentStr, "after oversized") {
		t.Error("Log file missing 'after oversized' message — connection did not survive oversized line")
	}

	// The truncated entry should be re-emitted as a valid zerolog warn entry
	if !strings.Contains(contentStr, `"_truncated":true`) {
		t.Error("Log file missing _truncated marker for the oversized entry")
	}
	if !strings.Contains(contentStr, `"level":"warn"`) {
		t.Error("Log file missing warn-level truncated entry")
	}
	// The full oversized content (70KB of x's) should NOT be in the file —
	// only the first 64KB chunk is captured in the truncated entry's message.
	if strings.Contains(contentStr, strings.Repeat("x", 70*1024)) {
		t.Error("Log file contains the full oversized content that should have been truncated")
	}
}

// TestLogServer_UnknownLogLevel verifies that a message with an unrecognized log level
// is still written to the file (log server is a transparent relay).
func TestLogServer_UnknownLogLevel(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	// Send message with unknown level — should be treated as info and written
	sendLogMessage(t, addr, "banana", "unknown level message")

	waitForFile(t, logFile, 2*time.Second)

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "unknown level message") {
		t.Error("Message with unknown log level was not written (should default to info priority)")
	}
}

// TestLogServer_InvalidJSON verifies that invalid JSON messages are handled gracefully
// without crashing the connection, and subsequent valid messages are still processed.
func TestLogServer_InvalidJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send invalid JSON
	_, err = conn.Write([]byte("this is not json\n"))
	if err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	// Send a valid message after the invalid one
	writeLogEntry(t, conn, "info", "valid after invalid")

	waitForFile(t, logFile, 2*time.Second)
	time.Sleep(200 * time.Millisecond)

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// The valid message after invalid JSON should still be processed
	if !strings.Contains(string(content), "valid after invalid") {
		t.Error("Valid message after invalid JSON was not written — connection did not survive")
	}
}

// TestLogServer_WritesToFile_NoRotation verifies basic file logging without rotation.
func TestLogServer_WritesToFile_NoRotation(t *testing.T) {
	// Create temp directory for log file
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start log server without rotation
	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Send a log message
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	sendLogMessage(t, addr, "info", "test message without rotation")

	// Wait for file to be written
	waitForFile(t, logFile, 2*time.Second)

	// Verify content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "test message without rotation") {
		t.Errorf("Log file does not contain expected message. Content: %s", content)
	}
}

// TestLogServer_WritesToFile_WithRotation verifies file logging with lumberjack rotation enabled.
func TestLogServer_WritesToFile_WithRotation(t *testing.T) {
	// Create temp directory for log file
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start log server with rotation enabled
	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: true,
		MaxSizeMB:       1, // 1MB for testing
		MaxBackups:      2,
		MaxAgeDays:      1,
		Compress:        false, // Disable compression for easier testing
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Send a log message
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	sendLogMessage(t, addr, "info", "test message with rotation enabled")

	// Wait for file to be written
	waitForFile(t, logFile, 2*time.Second)

	// Verify content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "test message with rotation enabled") {
		t.Errorf("Log file does not contain expected message. Content: %s", content)
	}
}

// TestLogServer_RotationTriggered writes ~2.5MB of logs to trigger rotation and
// confirms multiple rotated files are created.
func TestLogServer_RotationTriggered(t *testing.T) {
	// Create temp directory for log file
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start log server with very small max size to trigger rotation
	// Note: lumberjack's minimum effective size is 1MB, but we can test the setup
	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: true,
		MaxSizeMB:       1, // 1MB minimum
		MaxBackups:      3,
		MaxAgeDays:      0, // No age limit
		Compress:        true,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	// Connect once and reuse the connection for all messages
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send many log messages to trigger rotation
	// Each message is ~150 bytes, so 15000 messages = ~2.2MB (exceeds 1MB threshold twice)
	largeMessage := strings.Repeat("x", 100)
	for i := range 15000 {
		writeLogEntry(t, conn, "info", fmt.Sprintf("msg-%d-%s", i, largeMessage))
	}

	// Give time for writes to complete
	time.Sleep(500 * time.Millisecond)

	// Check that files exist in the directory
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	// We should have at least the main log file
	if len(files) == 0 {
		t.Fatal("No log files found in temp directory")
	}

	// Log what we found for debugging
	t.Logf("Found %d file(s) in temp directory:", len(files))
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			t.Errorf("Failed to get info for file %s: %v", f.Name(), err)
			continue
		}
		t.Logf("  - %s (size: %d bytes)", f.Name(), info.Size())
	}

	// Rotation should have triggered since we wrote ~2.2MB with a 1MB threshold
	if len(files) < 2 {
		t.Errorf("Expected rotation to trigger (at least 2 files), but found only %d file(s)", len(files))
	}
}

// TestLogServer_NoLevelFiltering verifies that the log server acts as a transparent relay:
// messages at ALL levels are written to the file regardless of their level.
func TestLogServer_NoLevelFiltering(t *testing.T) {
	// Create temp directory for log file
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	// Send messages at all levels — ALL should appear in the file (no filtering)
	sendLogMessage(t, addr, "trace", "trace message should appear")
	sendLogMessage(t, addr, "debug", "debug message should appear")
	sendLogMessage(t, addr, "info", "info message should appear")
	sendLogMessage(t, addr, "warn", "warn message should appear")
	sendLogMessage(t, addr, "error", "error message should appear")

	// Wait for file to be written and give time for all messages to be processed
	waitForFile(t, logFile, 2*time.Second)
	time.Sleep(200 * time.Millisecond)

	// Verify content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)

	// ALL messages should be present (no filtering at the log server)
	for _, expected := range []string{
		"trace message should appear",
		"debug message should appear",
		"info message should appear",
		"warn message should appear",
		"error message should appear",
	} {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Log file missing expected message: %s", expected)
		}
	}
}

// TestLogServer_GracefulShutdown verifies that context cancellation properly closes the log writer and flushes all pending data to the file.
// The test sends 5 messages, cancels the context to trigger shutdown, then verifies all messages are present in the log file and that new connections are refused.
// It verifies:
//  1. Data flushing - All log messages written before shutdown are properly flushed to the file
//  2. Connection rejection - After context cancellation, the server no longer accepts new TCP connections
func TestLogServer_GracefulShutdown(t *testing.T) {
	// Create temp directory for log file
	tmpDir, err := os.MkdirTemp("", "logserver_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")

	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Start log server
	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		RotationEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	// Send multiple log messages before shutdown
	for i := range 5 {
		sendLogMessage(t, addr, "info", fmt.Sprintf("pre-shutdown message %d", i))
	}

	// Wait for messages to be written
	waitForFile(t, logFile, 2*time.Second)
	time.Sleep(100 * time.Millisecond)

	// Trigger graceful shutdown by canceling context
	cancel()

	// Give shutdown time to complete
	time.Sleep(200 * time.Millisecond)

	// Verify all messages were flushed to file before shutdown
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	for i := range 5 {
		expectedMsg := fmt.Sprintf("pre-shutdown message %d", i)
		if !strings.Contains(contentStr, expectedMsg) {
			t.Errorf("Log file missing message: %s", expectedMsg)
		}
	}

	// Verify server is no longer accepting connections
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err == nil {
		conn.Close()
		t.Error("Expected connection to be refused after shutdown, but connection succeeded")
	}
}
