package logger

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

// Factory function: creates a Logger based on config name.
func NewLogger(kind string) Logger {
	switch kind {
	case "zerolog":
		return NewZeroLogger()
	default:
		return NewPrintfLogger() // default fallback
	}
}
