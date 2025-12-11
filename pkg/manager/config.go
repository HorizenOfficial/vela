package manager

import (
	"os"

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
	// ChannelType is the type of communication channel between manager and executor
	ChannelType string
	// ChannelParams are the parameters for the connection with the executor
	ChannelParams common.ChannelConnectionParams

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

	// LogKind is the type of logger to use (e.g., "zerolog", "tcplog")
	LogKind string

	// LogConsole is true if we want output on console
	LogConsole bool
	// LogConsoleLevel is the level of logging for the console
	LogConsoleLevel string
	// LogConsoleColor is true if the log output should be colored
	LogConsoleColor bool

	// LogFileName is the path to the log file. An empty string means no output on log file
	LogFileName string
	// LogFIleLevel is the level of logging for the console
	LogFileLevel string
}

func LoadConfig() (*Config, error) {
	fileProperties := properties.NewProperties()
	if common.FileExists(confFileName) {
		filePropertiesFromFile, err := properties.LoadFile(confFileName, properties.UTF8)
		if err != nil {
			return nil, err
		}
		fileProperties = filePropertiesFromFile
	}

	var channelType = common.GetConfigVar("CHANNEL_TYPE", "vsock", fileProperties)
	var channelConnectionParams common.ChannelConnectionParams
	executorServerPort := common.GetConfigVarInt64("EXECUTOR_PORT", 4000, fileProperties)
	if channelType == "vsock" {
		executorServerCid := common.GetConfigVarInt64("CHANNEL_VSOCK_CID", 20, fileProperties)
		channelConnectionParams = common.VSockChannelConnectionParams{CID: uint32(executorServerCid), Port: uint32(executorServerPort)}
	} else {
		executorIpAddress := common.GetConfigVar("EXECUTOR_IP_ADDRESS", "localhost", fileProperties)
		channelConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpAddress, Port: uint32(executorServerPort)}
	}

	var privateKey *cryptotypes.PrivateKeySecp256k1
	privateKeyFromEnv := os.Getenv("MANAGER_KEY_SECP256")
	if privateKeyFromEnv == "" {
		privateKey, _ = crypto.GeneratePrivateKeySecp256k1()
	} else {
		privateKey, _ = crypto.ImportPrivateKeySecp256k1FromHex(privateKeyFromEnv)
	}

	return &Config{
		ChannelType:   channelType,
		ChannelParams: channelConnectionParams,

		ReorgTimeout:              common.GetConfigVarInt64("REORG_TIMEOUT", 180, fileProperties), // 3 minutes
		HandshakeTimeout:          common.GetConfigVarInt64("HANDSHAKE_TIMEOUT", 5, fileProperties),
		BlockchainPollingInterval: common.GetConfigVarInt64("BLOCKCHAIN_POLLING_INTERVAL", 5, fileProperties),

		RpcURL: common.GetConfigVar("CHAIN_RPC_PROTOCOL", "http", fileProperties) + "://" +
			common.GetConfigVar("CHAIN_RPC_ADDRESS", "127.0.0.1", fileProperties) + ":" +
			common.GetConfigVar("CHAIN_RPC_PORT", "8545", fileProperties),

		PrivateKey:           *privateKey,
		ProcessorAddress:     common.GetConfigVar("CHAIN_PROCESSOR_ADDRESS", "", fileProperties),
		TeeAuthAddress:       common.GetConfigVar("CHAIN_TEEAUTHENTICATOR_ADDRESS", "", fileProperties),
		MockBlockChainClient: false,
		// Data layer configuration
		DataLayerType:             "versioned_leveldb",
		DataLayerDBPath:           common.GetConfigVar("MANAGER_DATA_FOLDER", "", fileProperties),
		DataLayerNumOfVersions:    10,
		DeanonymizationReportPath: common.GetConfigVar("MANAGER_REPORTS_FOLDER", "/tmp/horizen-pes-data/manager_reports", fileProperties),
		InputWasmPath:             common.GetConfigVar("MANAGER_INPUT_WASMS", "", fileProperties),
		LogKind:                   common.GetConfigVar("MANAGER_LOG_KIND", "zerolog", fileProperties),
		LogConsole:                common.GetConfigVarBool("MANAGER_LOG_CONSOLE", true, fileProperties),
		LogConsoleLevel:           common.GetConfigVar("MANAGER_LOG_CONSOLE_LEVEL", "info", fileProperties),
		LogConsoleColor:           common.GetConfigVarBool("MANAGER_LOG_CONSOLE_COLOR", false, fileProperties),
		LogFileName:               common.GetConfigVar("MANAGER_LOG_FILE_NAME", "", fileProperties),
		LogFileLevel:              common.GetConfigVar("MANAGER_LOG_FILE_LEVEL", "info", fileProperties),
	}, nil
}
