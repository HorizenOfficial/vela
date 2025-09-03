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

	// Blockchain client parameters
	// MockBlockChainClient specifies if the mock BlockChainClient should be used. Only for testing and development.
	MockBlockChainClient bool
	// Rpc address of the blockchain node
	RpcURL               string
	// TODO this is the priv key of the account used to sign transactions. For production, it should be managed by a secure vault.
	PrivateKey           string
	// Address of the ProcessorEndpoint contract
	ProcessorAddress    string
	// Address of the KeyRegistry contract
	KeyRegistryAddress  string

	// DataLayerType specifies the database implementation to use. Supported values: "versioned_leveldb", "mockdb".
	DataLayerType string
	// DataLayerDBPath is the full file path for the database.
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
		RpcURL: 			   "http://127.0.0.1:8545",
		PrivateKey: 			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a", // well known private key for local development
		ProcessorAddress: 		"0xCfEB869F69431e42cd62bB3aB5De8C3fB4d92dB", // well known address for local development
		KeyRegistryAddress: 	"0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1", // well known address for local development
		MockBlockChainClient:    false,
		// Data layer configuration
		DataLayerType:          "versioned_leveldb",
		DataLayerDBPath:        "/tmp/horizen-pes-data/manager.db",
		DataLayerNumOfVersions: 10, // useful only for versioned leveldb
	}
}
