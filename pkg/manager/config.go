package manager

import (
	"os"
	"strconv"

	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/crypto"
	"github.com/magiconair/properties"
)

const confFileName = "manager.conf"

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

	// LogKind is the type of logger to use (e.g., "zeronetwork", "tcplog")
	LogKind string

	// LogConsole is true if we want output on console
	LogConsole bool
	// LogConsoleLevel is the level of logging for the console
	LogConsoleLevel string
	// LogConsoleColor is true if the log output should be colored
	LogConsoleColor bool

	// LogFileName is the path to the log file. An empty string means no output on log file
	LogFileName string
	// LogFileLevel is the level of logging for the console
	LogFileLevel string

	// RemoteLogAddress is the address for remote logging (e.g., "127.0.0.1:12345"). An empty string means no remote logging.
	RemoteLogAddress string
	// RemoteLogNetwork is the network for remote logging (e.g., "tcp", "vsock").
	RemoteLogNetwork string
	// NetworkLevel is the level of logging for the network
	NetworkLevel string

	// LogServerTCPAddress is the address where the manager's log server will listen for incoming TCP log messages.
	LogServerTCPAddress string
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

	logKind := os.Getenv("MANAGER_LOG_KIND")
	if logKind == "" {
		logKind = "zerolog"
	}

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
		reorgTimeout = 180
	}

	handshakeTimeoutEnvVar := os.Getenv("HANDSHAKE_TIMEOUT")
	handshakeTimeout, err := strconv.ParseInt(handshakeTimeoutEnvVar, 10, 32)
	if err != nil {
		handshakeTimeout = 5
	}

	blockchainPollingIntervalEnvVar := os.Getenv("BLOCKCHAIN_POLLING_INTERVAL")
	blockchainPollingInterval, err := strconv.ParseInt(blockchainPollingIntervalEnvVar, 10, 32)
	if err != nil {
		blockchainPollingInterval = 5
	}

	logConsole, err := strconv.ParseBool(os.Getenv("MANAGER_LOG_CONSOLE"))
	if err != nil {
		logConsole = true
	}

	logConsoleLevel := os.Getenv("MANAGER_LOG_CONSOLE_LEVEL")
	if logConsoleLevel == "" {
		logConsoleLevel = "info"
	}

	logFileName := os.Getenv("MANAGER_LOG_FILE_NAME")
	if logFileName == "" {
		logFileName = ""
	}

	logFileLevel := os.Getenv("MANAGER_LOG_FILE_LEVEL")
	if logFileLevel == "" {
		if logFileName != "" {
			logFileLevel = "info"
		}
	}

	logConsoleColor, err := strconv.ParseBool(os.Getenv("MANAGER_LOG_CONSOLE_COLOR"))
	if err != nil {
		logConsoleColor = false
	}

	remoteLogAddress := os.Getenv("MANAGER_LOG_REMOTE_ADDRESS")
	remoteLogNetwork := os.Getenv("MANAGER_LOG_REMOTE_NETWORK")
	networkLevel := os.Getenv("MANAGER_LOG_NETWORK_LEVEL")
	if networkLevel == "" {
		networkLevel = "info"
	}

	logServerTCPAddress := os.Getenv("MANAGER_LOG_SERVER_TCP_ADDRESS")

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

		LogKind:         logKind,
		LogConsole:      logConsole,
		LogConsoleLevel: logConsoleLevel,
		LogConsoleColor: logConsoleColor,

		LogFileName:         logFileName,
		LogFileLevel:        logFileLevel,
		RemoteLogAddress:    remoteLogAddress,
		RemoteLogNetwork:    remoteLogNetwork,
		NetworkLevel:        networkLevel,
		LogServerTCPAddress: logServerTCPAddress,
	}
}

// LoadConfigFromFile loads the configuration from the given file
func LoadConfigFromFile() (*Config, error) {
	if !common.FileExists(confFileName) {
		return DefaultConfig(), nil
	}
	// Load properties from file
	config, err := properties.LoadFile(confFileName, properties.UTF8)
	if err != nil {
		return nil, err
	}
	PrivateKey, err := crypto.ImportPrivateKeySecp256k1FromHex(config.MustGetString("PrivateKey"))
	if err != nil {
		return nil, err
	}
	return &Config{
		ReorgTimeout:              config.GetInt64("ReorgTimeout", 180), // 3 minutes
		HandshakeTimeout:          config.GetInt64("HandshakeTimeout", 5),
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
		LogKind:                   config.GetString("LogKind", "zerolog"),
		LogConsole:                config.GetBool("LogConsole", true),
		LogConsoleLevel:           config.GetString("LogConsoleLevel", "info"),
		LogConsoleColor:           config.GetBool("LogColor", true),
		LogFileName:               config.GetString("LogFileName", ""),
		LogFileLevel:              config.GetString("LogFileLevel", "info"),
		RemoteLogAddress:          config.GetString("LogRemoteAddress", ""),
		RemoteLogNetwork:          config.GetString("LogRemoteNetwork", ""),
		NetworkLevel:              config.GetString("LogNetworkLevel", "info"),
		LogServerTCPAddress:       config.GetString("LogServerTCPAddress", ""),
	}, nil
}
