package logger

import (
	"github.com/rs/zerolog"
)

// init is a special Go function that is executed automatically when this package
// is first imported. It's used here to set up global configurations for the
// underlying zerolog library.
//
// Any zerolog instance created within this application will automatically inherit
// these settings, ensuring consistent logging behavior across all logger types
// (like ZeroLogger and ZeroNetworkLogger).
func init() {
	// By setting this global variable, all zerolog timestamps will have
	// nanosecond precision and be formatted as a human-readable string.
	//zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimeFieldFormat = "2006-Dec-02 15:04:05.000"

	// Zerolog's default internal skip is usually 2.
	// By setting it to 3, we are adding 1 extra skip for our wrapper functions,
	// ensuring the reported caller is correct.
	zerolog.CallerSkipFrameCount = 3
}
