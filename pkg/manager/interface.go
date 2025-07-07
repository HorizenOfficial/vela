package manager

import (
	"context"
)

// Manager defines the interface for the Secure Processor Manager
type Manager interface {
	// Start starts the manager
	Start(ctx context.Context) error
	// Stop stops the manager
	Stop() error
}

// Config defines the configuration for the Secure Processor Manager
type Config struct {
	// BlockchainPollingInterval is the interval at which to poll the blockchain for new requests
	BlockchainPollingInterval int64
	// ExecutorConnectionType is the type of connection to use for the executor
	ExecutorConnectionType string
	// ExecutorConnectionParams are the parameters for the executor connection
	ExecutorConnectionParams map[string]string

	// TODO: blockchain connection parameters
	// TODO: dataLayer connection parameters
}

// DefaultConfig returns the default configuration for the Secure Processor Manager
func DefaultConfig() *Config {
	return &Config{
		BlockchainPollingInterval: 5,     // 5 seconds
		ExecutorConnectionType:    "tcp", // or "vsock"
		ExecutorConnectionParams: map[string]string{
			"url": "localhost:10000",
		},
	}
}
