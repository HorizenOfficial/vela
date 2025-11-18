package logger

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

// Config holds the configuration for the logger.
type Config struct {
	Kind         string
	Console      bool
	ConsoleLevel string
	ConsoleColor bool
	FileName     string
	FileLevel    string
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
	default:
		return NewPrintfLogger(cfg) // default fallback
	}
}
