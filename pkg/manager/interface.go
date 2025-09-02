package manager

import (
	"context"
	"github.com/horizen-pes/pkg/communication"
)

// Manager defines the interface for the Secure Processor Manager
type Manager interface {
	// Start starts the manager
	Start(ctx context.Context) error
	// Stop stops the manager
	Stop() error
	communication.ClientRequestHandler
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

	// DataLayerType specifies the database implementation to use. Supported values: "versioned_leveldb", "mockdb".
	DataLayerType string
	// DataLayerDBPath is the path for the database. For "versioned_leveldb", this is a base directory.
	DataLayerDBPath string
	// DataLayerNumOfVersions specifies how many historical versions to keep. Only used by "versioned_leveldb".
	DataLayerNumOfVersions int
}

// TODO maybe it is best using a config file
// DefaultConfig returns the default configuration for the Secure Processor Manager
func DefaultConfig() *Config {
	return &Config{
		BlockchainPollingInterval: 5,     // 5 seconds
		ExecutorConnectionType:    "tcp", // or "vsock"
		ExecutorConnectionParams: map[string]string{
			"url": "localhost:8080",
		},
		DataLayerType:          "versioned_leveldb",
		DataLayerDBPath:        "/tmp/horizen-pes-data/manager_db",
		DataLayerNumOfVersions: 10, // useful only for versioned leveldb
	}
}
