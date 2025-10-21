package manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/crypto"
	"github.com/magiconair/properties"
)

const confFileName = "manager.conf"

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
	// ReorgTimeout is the max time interval for waiting for a reorg to be resolved (in seconds)
	ReorgTimeout int64

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
	RpcURL string
	// TODO this is the priv key of the account used to sign transactions. For production, it should be managed by a secure vault.
	PrivateKey cryptotypes.PrivateKeySecp256k1

	// Address of the ProcessorEndpoint contract
	ProcessorAddress string
	// Address of the TeeAuthenticator contract
	TeeAuthAddress string

	// DataLayerType specifies the database implementation to use. Supported values: "versioned_leveldb", "mockdb".
	DataLayerType string
	// DataLayerDBPath is the path for the database. For "versioned_leveldb", this is a base directory.
	DataLayerDBPath string
	// DataLayerNumOfVersions specifies how many historical versions to keep. Only used by "versioned_leveldb".
	DataLayerNumOfVersions int
}

// DefaultConfig returns the default configuration for the Secure Processor Manager
func DefaultConfig() *Config {
	executorServerAddress := os.Getenv("EXECUTOR_IP_ADDRESS")
	if executorServerAddress == "" {
		executorServerAddress = "localhost"
	}
	executorServerPort := os.Getenv("EXECUTOR_IP_PORT")
	if executorServerPort == "" {
		executorServerPort = "8080"
	}
	PrivateKey, _ := crypto.ImportPrivateKeySecp256k1FromHex(os.Getenv("MANAGER_PRIVATE_KEY")) // well known private key for local development
	dataPath := os.Getenv("MANAGER_DATA_FOLDER")
	if dataPath == "" {
		dataPath = "/tmp/horizen-pes-data/manager_db"
	}
	nodeUrl := os.Getenv("CHAIN_RPC_ADDRESS")
	if nodeUrl == "" {
		nodeUrl = "127.0.0.1"
	}
	nodePort := os.Getenv("CHAIN_RPC_PORT")
	if nodePort == "" {
		nodePort = "8545"
	}
	processorAddress := os.Getenv("CHAIN_PROCESSOR_ADDRESS")
	teeAuthAddress := os.Getenv("CHAIN_TEEAUTHENTICATOR_ADDRESS")

	reorgTimeoutEnvVar := os.Getenv("REORG_TIMEOUT")
	reorgTimeout, err := strconv.ParseInt(reorgTimeoutEnvVar, 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert REORG_TIMEOUT for error %v, using default value", err)
		reorgTimeout = 180
	}

	blockchainPollingIntervalEnvVar := os.Getenv("BLOCKCHAIN_POLLING_INTERVAL")
	blockchainPollingInterval, err := strconv.ParseInt(blockchainPollingIntervalEnvVar, 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert BLOCKCHAIN_POLLING_INTERVAL for error %v, using default value", err)
		blockchainPollingInterval = 5
	}

	return &Config{
		ReorgTimeout:              reorgTimeout, 
		BlockchainPollingInterval: blockchainPollingInterval,     
		ExecutorConnectionType:    "tcp", // or "vsock"
		ExecutorConnectionParams: map[string]string{
			"url": executorServerAddress + ":" + executorServerPort,
		},
		RpcURL:               "http://" + nodeUrl + ":" + nodePort,
		PrivateKey:           *PrivateKey,
		ProcessorAddress:     processorAddress,
		TeeAuthAddress: 	  teeAuthAddress,

		MockBlockChainClient: false,
		// Data layer configuration
		DataLayerType:          "versioned_leveldb",
		DataLayerDBPath:        dataPath,
		DataLayerNumOfVersions: 10, // useful only for versioned leveldb
	}
}

func ReadConfig() *Config {
	if !fileExists(confFileName) {
		log.Printf("File %s not found. Using default configuration for Manager.\n", confFileName)
		return DefaultConfig()
	}
	// Load properties from file
	config, err := properties.LoadFile(confFileName, properties.UTF8)
	if err != nil {
		panic(err)
	}
	PrivateKey, err := crypto.ImportPrivateKeySecp256k1FromHex(config.MustGetString("PrivateKey"))
	if err != nil {
		log.Printf("Error loading PrivateKey from config file: %v.", err)
		panic(err)
	}
	return &Config{
		ReorgTimeout:              config.GetInt64("ReorgTimeout", 180), // 3 minutes
		BlockchainPollingInterval: config.GetInt64("BlockchainPollingInterval", 1),
		ExecutorConnectionType:    config.MustGetString("ExecutorConnectionType"),
		ExecutorConnectionParams: map[string]string{
			"url": config.MustGetString("ExecutorConnectionUrl"),
		},
		RpcURL:               config.MustGetString("RpcUrl"),
		PrivateKey:           *PrivateKey,
		ProcessorAddress:     config.MustGetString("ProcessorAddress"),
		TeeAuthAddress:	      config.MustGetString("TeeAuthenticatorAddress"),
		MockBlockChainClient: config.MustGetBool("MockBlockChainClient"),
		// Data layer configuration
		DataLayerType:          config.MustGetString("DataLayerType"),
		DataLayerDBPath:        config.MustGetString("DataLayerDBPath"),
		DataLayerNumOfVersions: config.MustGetInt("DataLayerNumOfVersions"),
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
