package utils

import "fmt"

// Logger for TinyGo WASI guest applications.
// The host parses these prefixes and routes to appropriate log levels.
//
// Note: fmt.Println writes to stdout, which is captured by the WASI pipe.

// LogTrace writes a trace-level log message (captured by WASI pipe)
func LogTrace(format string, args ...any) {
	logf("TRC", format, args...)
}

// LogDebug writes a debug-level log message (captured by WASI pipe)
func LogDebug(format string, args ...any) {
	logf("DBG", format, args...)
}

// LogInfo writes an info-level log message (captured by WASI pipe)
func LogInfo(format string, args ...any) {
	logf("INF", format, args...)
}

// LogWarn writes a warning-level log message (captured by WASI pipe)
func LogWarn(format string, args ...any) {
	logf("WRN", format, args...)
}

// LogError writes an error-level log message (captured by WASI pipe)
func LogError(format string, args ...any) {
	logf("ERR", format, args...)
}

func logf(prefix, format string, args ...any) {
	if len(args) == 0 {
		fmt.Println(prefix, format)
	} else {
		fmt.Println(prefix, fmt.Sprintf(format, args...))
	}
}
