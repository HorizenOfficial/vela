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

	// HandshakeTimeout is the max time interval for waiting for the executor handshake to be completed (in seconds)
	HandshakeTimeout int64

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

	// DeanonymizationReportPath is the path to a folder where to store deanonymization reports.
	DeanonymizationReportPath string

	// InputWasmPath is the path where the wasm bytecode to be deployed is retrieved if not found in the payload
	// (we may need to load it externally for GAS limitation)
	InputWasmPath string
}

// DefaultConfig returns the default configuration for the Secure Processor Manager  (possibly overridden by env variables)
func DefaultConfig() *Config {
	executorServerAddress := os.Getenv("EXECUTOR_IP_ADDRESS")
	if executorServerAddress == "" {
		executorServerAddress = "localhost"
	}
	executorServerPort := os.Getenv("EXECUTOR_IP_PORT")
	if executorServerPort == "" {
		executorServerPort = "8080"
	}
	var privateKey *cryptotypes.PrivateKeySecp256k1
	privateKeyFromEnv := os.Getenv("MANAGER_KEY_SECP256")
	if privateKeyFromEnv == "" {
		privateKey, _ = crypto.GeneratePrivateKeySecp256k1()
	} else {
		privateKey, _ = crypto.ImportPrivateKeySecp256k1FromHex(privateKeyFromEnv)
	}
	dataPath := os.Getenv("MANAGER_DATA_FOLDER")
	if dataPath == "" {
		dataPath = "/tmp/horizen-pes-data/manager_db"
	}
	inputWasmPath := os.Getenv("MANAGER_INPUT_WASMS")
	nodeProtocol := os.Getenv("CHAIN_RPC_PROTOCOL")
	if nodeProtocol == "" {
		nodeProtocol = "http"
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

	reportsPath := os.Getenv("MANAGER_REPORTS_FOLDER")

	reorgTimeoutEnvVar := os.Getenv("REORG_TIMEOUT")
	reorgTimeout, err := strconv.ParseInt(reorgTimeoutEnvVar, 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert REORG_TIMEOUT for error %v, using default value\n", err)
		reorgTimeout = 180
	}

	handshakeTimeoutEnvVar := os.Getenv("HANDSHAKE_TIMEOUT")
	handshakeTimeout, err := strconv.ParseInt(handshakeTimeoutEnvVar, 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert HANDSHAKE_TIMEOUT for error %v, using default value\n", err)
		handshakeTimeout = 30
	}

	blockchainPollingIntervalEnvVar := os.Getenv("BLOCKCHAIN_POLLING_INTERVAL")
	blockchainPollingInterval, err := strconv.ParseInt(blockchainPollingIntervalEnvVar, 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert BLOCKCHAIN_POLLING_INTERVAL for error %v, using default value\n", err)
		blockchainPollingInterval = 5
	}

	return &Config{
		ReorgTimeout:              reorgTimeout,
		HandshakeTimeout:          handshakeTimeout,
		BlockchainPollingInterval: blockchainPollingInterval,
		ExecutorConnectionType:    "tcp", // or "vsock"
		ExecutorConnectionParams: map[string]string{
			"url": executorServerAddress + ":" + executorServerPort,
		},
		RpcURL:           nodeProtocol + "://" + nodeUrl + ":" + nodePort,
		PrivateKey:       *privateKey,
		ProcessorAddress: processorAddress,
		TeeAuthAddress:   teeAuthAddress,

		MockBlockChainClient: false,
		// Data layer configuration
		DataLayerType:             "versioned_leveldb",
		DataLayerDBPath:           dataPath,
		DataLayerNumOfVersions:    10,          // useful only for versioned leveldb
		DeanonymizationReportPath: reportsPath, // optional, default to not-there semantic
		InputWasmPath:             inputWasmPath,
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
		HandshakeTimeout:          config.GetInt64("HandshakeTimeout", 30),
		BlockchainPollingInterval: config.GetInt64("BlockchainPollingInterval", 1),
		ExecutorConnectionType:    config.MustGetString("ExecutorConnectionType"),
		ExecutorConnectionParams: map[string]string{
			"url": config.MustGetString("ExecutorConnectionUrl"),
		},
		RpcURL:               config.MustGetString("RpcUrl"),
		PrivateKey:           *PrivateKey,
		ProcessorAddress:     config.MustGetString("ProcessorAddress"),
		TeeAuthAddress:       config.MustGetString("TeeAuthenticatorAddress"),
		MockBlockChainClient: config.MustGetBool("MockBlockChainClient"),
		// Data layer configuration
		DataLayerType:             config.MustGetString("DataLayerType"),
		DataLayerDBPath:           config.MustGetString("DataLayerDBPath"),
		DataLayerNumOfVersions:    config.MustGetInt("DataLayerNumOfVersions"),
		DeanonymizationReportPath: config.GetString("DeanonymizationReportPath", ""),
		InputWasmPath:             config.GetString("InputWasmPath", ""),
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
