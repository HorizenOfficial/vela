package logger

import (
	"bytes"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	// Start a dummy listener for network-based loggers to connect to
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	go func() {
		conn, _ := ln.Accept()
		if conn != nil {
			conn.Close()
		}
	}()

	testCases := []struct {
		name       string
		kind       string
		expected   interface{}
		shouldPanic bool
		remoteAddr string
	}{
		{"zerolog", "zerolog", (*ZeroLogger)(nil), false, ""},
		{"zeronetwork", "zeronetwork", (*ZeroNetworkLogger)(nil), false, ln.Addr().String()},
		{"tcplog", "tcplog", (*TCPLogger)(nil), false, ln.Addr().String()},
		{"default", "unknown", (*PrintfLogger)(nil), false, ""},
		{"printf", "printf", (*PrintfLogger)(nil), false, ""},
		{"tcplog_panic_no_address", "tcplog", nil, true, ""},
		{"zeronetwork_panic_no_address", "zeronetwork", nil, true, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				Kind:             tc.kind,
				Console:          true,
				ConsoleLevel:     "info",
				RemoteLogAddress: tc.remoteAddr,
			}

			if tc.shouldPanic {
				require.Panics(t, func() {
					NewLogger(cfg)
				})
				return
			}

			logger := NewLogger(cfg)
			require.NotNil(t, logger)
			assert.IsType(t, tc.expected, logger)
			logger.Close()
		})
	}
}

// captureOutput captures everything written to os.Stderr during the execution of f.
func captureOutput(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	originalStderr := os.Stderr
	os.Stderr = w

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		out <- buf.String()
	}()

	f()

	os.Stderr = originalStderr
	w.Close()

	return <-out
}

func TestZeroLogger_SetLevel(t *testing.T) {
	cfg := &Config{
		Kind:         "zerolog",
		Console:      true,
		ConsoleLevel: "info",
		ConsoleColor: false,
	}

	r, w, err := os.Pipe()
	require.NoError(t, err)

	originalStderr := os.Stderr
	os.Stderr = w
	defer func() {
		os.Stderr = originalStderr
	}()

	logger := NewLogger(cfg)

	// Test initial level
	logger.Debug("this should not appear")

	// Change level and test again
	err = logger.SetLevel("debug")
	require.NoError(t, err)

	logger.Debug("this should appear now")

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NotContains(t, output, "this should not appear")
	assert.Contains(t, output, "this should appear now")

	// Test invalid level
	err = logger.SetLevel("invalid")
	require.Error(t, err)

	logger.Close()
}

func TestZeroNetworkLogger_SetLevel(t *testing.T) {
	addr := "127.0.0.1:0"
	ln, err := net.Listen("tcp", addr)
	require.NoError(t, err)
	defer ln.Close()

	msgChan := make(chan string, 10)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		reader := bytes.NewBuffer(nil)
		for {
			_, err := io.CopyN(reader, conn, 1)
			if err != nil {
				return
			}
			if reader.Bytes()[reader.Len()-1] == '\n' {
				msgChan <- reader.String()
				reader.Reset()
			}
		}
	}()

	cfg := &Config{
		Kind:             "zeronetwork",
		Console:          false,
		NetworkLevel:     "info",
		RemoteLogAddress: ln.Addr().String(),
		RemoteLogNetwork: "tcp",
	}

	logger := NewLogger(cfg)
	defer logger.Close()

	// Wait for the async logger to connect
	time.Sleep(100 * time.Millisecond)

	logger.Debug("should not be sent")
	logger.Info("should be sent")

	select {
	case msg := <-msgChan:
		assert.Contains(t, msg, "should be sent")
		assert.Contains(t, msg, "\"level\":\"info\"")
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for info message")
	}
	assert.Len(t, msgChan, 0, "should have only received one message")

	// Update level
	err = logger.SetLevel("debug")
	require.NoError(t, err)

	logger.Debug("should now be sent")

	select {
	case msg := <-msgChan:
		assert.Contains(t, msg, "should now be sent")
		assert.Contains(t, msg, "\"level\":\"debug\"")
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for debug message")
	}
}

func TestPrintfLogger_SetLevel(t *testing.T) {
	cfg := &Config{Kind: "printf"}
	logger := NewLogger(cfg)

	// SetLevel on PrintfLogger is a no-op, should not error
	err := logger.SetLevel("any-level")
	require.NoError(t, err)
	logger.Close()
}
