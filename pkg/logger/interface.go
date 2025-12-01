package logger

type Logger interface {
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)
	SetLevel(level string) error
	Close() error
}

// Config holds the configuration for the logger.
type Config struct {
	Kind             string
	Console          bool
	ConsoleLevel     string
	ConsoleColor     bool
	FileName         string
	FileLevel        string
	RemoteLogAddress string // Address for remote logging (e.g., "127.0.0.1:12345")
	RemoteLogNetwork string // Network for remote logging (e.g., "tcp", "vsock")
	NetworkLevel     string // Level for remote logging
}

// minimal cfg
func DefaultLogConfig(kind string) Config {
	return Config{
		Kind:         kind,
		Console:      true,
		ConsoleLevel: "info",
	}
}

// Factory function: creates a Logger based on config name.
func NewLogger(cfg *Config) Logger {
	switch cfg.Kind {
	case "zerolog":
		return NewZeroLogger(cfg)
	case "tcplog":
		return NewTCPLogger(cfg)
	case "zeronetwork":
		return NewZeroNetworkLogger(cfg)
	default:
		return NewPrintfLogger(cfg) // default fallback
	}
}
