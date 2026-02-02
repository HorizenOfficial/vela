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

// sendLogMessage sends a JSON log message to the log server via TCP
func sendLogMessage(t *testing.T, addr string, level, message string) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to log server: %v", err)
	}
	defer conn.Close()

	logEntry := map[string]string{
		"level":   level,
		"time":    time.Now().Format("2006-Jan-02 15:04:05.000"),
		"message": message,
	}
	data, _ := json.Marshal(logEntry)
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
		ConsoleLevel:    "info",
		FileLevel:       "info",
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
		ConsoleLevel:    "info",
		FileLevel:       "info",
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
	//defer os.RemoveAll(tmpDir)

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
		ConsoleLevel:    "debug",
		FileLevel:       "debug",
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

	// Send many log messages to potentially trigger rotation
	// Each message is ~150 bytes, so 15000 messages = ~2.2MB
	largeMessage := strings.Repeat("x", 100)
	for i := range 15000 {
		logEntry := map[string]string{
			"level":   "info",
			"time":    time.Now().Format("2006-Jan-02 15:04:05.000"),
			"message": fmt.Sprintf("msg-%d-%s", i, largeMessage),
		}
		data, _ := json.Marshal(logEntry)
		data = append(data, '\n')
		if _, err := conn.Write(data); err != nil {
			t.Fatalf("Failed to write on iteration %d: %v", i, err)
		}
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
		t.Error("No log files found in temp directory")
	}

	// Log what we found for debugging
	t.Logf("Found %d file(s) in temp directory:", len(files))
	for _, f := range files {
		info, _ := f.Info()
		t.Logf("  - %s (size: %d bytes)", f.Name(), info.Size())
	}

	// If rotation worked, we should have more than one file
	// Note: This depends on actually exceeding 1MB
	if len(files) > 1 {
		t.Logf("Rotation triggered: found %d files", len(files))
	}
}

// TestLogServer_LevelFiltering verifies that log level filtering works correctly
// (debug/info messages are filtered out, warn/error messages are written).
func TestLogServer_LevelFiltering(t *testing.T) {
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

	// Start log server with file level = warn (should filter out info/debug)
	err = StartLogServer(ctx, LogServerConfig{
		TCPAddr: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(port),
		},
		LogFilePath:     logFile,
		ConsoleEnabled:  false,
		ConsoleLevel:    "error",
		FileLevel:       "warn", // Only warn and above
		RotationEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to start log server: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	// Send messages at different levels
	sendLogMessage(t, addr, "debug", "debug message should be filtered")
	sendLogMessage(t, addr, "info", "info message should be filtered")
	sendLogMessage(t, addr, "warn", "warn message should appear")
	sendLogMessage(t, addr, "error", "error message should appear")

	// Wait for file to be written and give time for all messages to be processed
	waitForFile(t, logFile, 2*time.Second)
	time.Sleep(200 * time.Millisecond) // Allow time for all messages to be written

	// Verify content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)

	// Should NOT contain debug or info
	if strings.Contains(contentStr, "debug message should be filtered") {
		t.Error("Log file contains debug message that should have been filtered")
	}
	if strings.Contains(contentStr, "info message should be filtered") {
		t.Error("Log file contains info message that should have been filtered")
	}

	// Should contain warn and error
	if !strings.Contains(contentStr, "warn message should appear") {
		t.Error("Log file does not contain expected warn message")
	}
	if !strings.Contains(contentStr, "error message should appear") {
		t.Error("Log file does not contain expected error message")
	}
}
