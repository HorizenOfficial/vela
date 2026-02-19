package manager

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	// BlockchainPollingInterval is the interval at which to poll the blockchain for new requests (in seconds)
	BlockchainPollingInterval int64
	// BlockchainConnectTimeout is the max time to wait for the initial RPC handshake (ChainID) when connecting (in seconds)
	BlockchainConnectTimeout int64
	// ChannelType is the type of communication channel between manager and executor
	ChannelType string
	// ChannelParams are the parameters for the connection with the executor
	ChannelParams common.ChannelConnectionParams

	// AdminChannelParams are the parameters for the admin command server
	AdminChannelParams common.ChannelConnectionParams
	// AdminCommunicationParams holds parameters for communication towards the admin server
	AdminCommunicationParams common.CommunicationParams

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

	// Manager logging
	//--------------------------
	// LogKind is the type of logger to use (e.g., "zeronetwork", "zerolog")
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
	// LogNetworkLevel is the level of logging for the network
	LogNetworkLevel string

	// Log Server
	//-------------------------
	// LogServerTCPAddress contains the TCP address where the log server will listen for incoming log messages.
	LogServerTCPAddress common.TcpChannelConnectionParams
	// LogServerVSockAddress is the address where the log server will listen for incoming VSOCK log messages.
	LogServerVSockAddress common.VSockChannelConnectionParams
	// LogServerLogFile is the file path where the manager's log server will write incoming log messages.
	LogServerLogFile string
	// LogServerConsole is true if we want output on console
	LogServerConsole bool
	// LogServerRotationEnabled enables log rotation using lumberjack (only when LogServerLogFile is set)
	LogServerRotationEnabled bool
	// LogServerMaxSizeMB is the max size in megabytes before rotation (default: 100)
	LogServerMaxSizeMB int
	// LogServerMaxBackups is the max number of old log files to retain (default: 3)
	LogServerMaxBackups int
	// LogServerMaxAgeDays is the max days to retain old log files, 0 = no limit (default: 28)
	LogServerMaxAgeDays int
	// LogServerCompress enables gzip compression for rotated log files (default: true)
	LogServerCompress bool

	// CommunicationParams holds parameters for communication between manager and executor
	CommunicationParams common.CommunicationParams
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

	// We use the same port value for TCP or Vsock channel types
	executorServerPort := common.GetConfigVarInt64("EXECUTOR_PORT", 4000, fileProperties)
	logServerPort := common.GetConfigVarInt64("LOG_SERVER_PORT", 5000, fileProperties)

	// channel communication between manager and executor
	var channelConnectionParams common.ChannelConnectionParams

	// log server Vsock address, initialized only in case channel type is Vsock
	var logServerVsockAddress common.VSockChannelConnectionParams

	if channelType == "vsock" {
		// CID >= 16 are available values for EC2 enclaves (where executor runs)
		// CID and port are both used when connecting to a server
		executorServerCid := common.GetConfigVarInt64("EXECUTOR_VSOCK_CID", 20, fileProperties)
		channelConnectionParams = common.VSockChannelConnectionParams{CID: uint32(executorServerCid), Port: uint32(executorServerPort)}
		// if channel is vsock it means we also have a vsock connection used by executor for logging, we use of course a separate port
		// CID is not used actually when creating a lstening server
		logServerVsockAddress = common.VSockChannelConnectionParams{Port: uint32(logServerPort)}
	} else {
		executorIpHost := common.GetConfigVar("EXECUTOR_IP_HOST", "localhost", fileProperties)
		channelConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpHost, Port: uint32(executorServerPort)}
	}

	// TCP connection on log server, typically used for manager logs
	logServerTcpHost := common.GetConfigVar("LOG_SERVER_IP_HOST", "localhost", fileProperties)
	logServerTcpAddress := common.TcpChannelConnectionParams{Ip: logServerTcpHost, Port: uint32(logServerPort)}

	var privateKey *cryptotypes.PrivateKeySecp256k1
	privateKeyFromEnv := os.Getenv("MANAGER_KEY_SECP256")
	if privateKeyFromEnv == "" {
		privateKey, _ = crypto.GeneratePrivateKeySecp256k1()
	} else {
		privateKey, _ = crypto.ImportPrivateKeySecp256k1FromHex(privateKeyFromEnv)
	}

	communicationParams := common.CommunicationParams{
		RequestTimeoutSec: time.Duration(common.GetConfigVarInt64("MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC", 30, fileProperties)),
	}

	// Admin command server configuration
	adminServerPort := common.GetConfigVarInt64("MANAGER_ADMIN_PORT", 4002, fileProperties)
	var adminChannelConnectionParams common.ChannelConnectionParams
	// Admin server always uses TCP for external access
	adminServerHost := common.GetConfigVar("MANAGER_IP_HOST", "localhost", fileProperties)
	adminChannelConnectionParams = common.TcpChannelConnectionParams{Ip: adminServerHost, Port: uint32(adminServerPort)}

	adminCommunicationParams := common.CommunicationParams{
		RequestTimeoutSec: time.Duration(common.GetConfigVarInt64("MANAGER_ADMIN_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC", 30, fileProperties)),
	}

	cfg := &Config{
		ChannelType:              channelType,
		ChannelParams:            channelConnectionParams,
		AdminChannelParams:       adminChannelConnectionParams,
		AdminCommunicationParams: adminCommunicationParams,

		ReorgTimeout:              common.GetConfigVarInt64("REORG_TIMEOUT", 180, fileProperties), // 3 minutes
		HandshakeTimeout:          common.GetConfigVarInt64("HANDSHAKE_TIMEOUT", 5, fileProperties),
		BlockchainPollingInterval: common.GetConfigVarInt64("BLOCKCHAIN_POLLING_INTERVAL", 5, fileProperties),
		BlockchainConnectTimeout: common.GetConfigVarInt64("BLOCKCHAIN_CONNECT_TIMEOUT", 10, fileProperties),

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
		LogKind:                   common.GetConfigVar("MANAGER_LOG_KIND", "zeronetwork", fileProperties),
		LogConsole:                common.GetConfigVarBool("MANAGER_LOG_CONSOLE", true, fileProperties),
		LogConsoleLevel:           common.GetConfigVar("MANAGER_LOG_CONSOLE_LEVEL", "info", fileProperties),
		LogConsoleColor:           common.GetConfigVarBool("MANAGER_LOG_CONSOLE_COLOR", false, fileProperties),
		LogFileName:               common.GetConfigVar("MANAGER_LOG_FILE_NAME", "", fileProperties),
		LogFileLevel:              common.GetConfigVar("MANAGER_LOG_FILE_LEVEL", "info", fileProperties),
		LogNetworkLevel:           common.GetConfigVar("MANAGER_LOG_NETWORK_LEVEL", "info", fileProperties),
		LogServerTCPAddress:       logServerTcpAddress,
		LogServerVSockAddress:     logServerVsockAddress,
		LogServerLogFile:          common.GetConfigVar("LOG_SERVER_FILE_NAME", "", fileProperties),
		LogServerConsole:          common.GetConfigVarBool("LOG_SERVER_CONSOLE", true, fileProperties),
		LogServerRotationEnabled:  common.GetConfigVarBool("LOG_SERVER_FILE_ROTATION", false, fileProperties),
		LogServerMaxSizeMB:        int(common.GetConfigVarInt64("LOG_SERVER_FILE_MAX_SIZE_MB", 100, fileProperties)),
		LogServerMaxBackups:       int(common.GetConfigVarInt64("LOG_SERVER_FILE_MAX_BACKUPS", 3, fileProperties)),
		LogServerMaxAgeDays:       int(common.GetConfigVarInt64("LOG_SERVER_FILE_MAX_AGE_DAYS", 28, fileProperties)),
		LogServerCompress:         common.GetConfigVarBool("LOG_SERVER_FILE_COMPRESS", true, fileProperties),
		CommunicationParams:       communicationParams,
	}

	if strings.TrimSpace(cfg.DeanonymizationReportPath) == "" {
		return nil, fmt.Errorf("MANAGER_REPORTS_FOLDER (DeanonymizationReportPath) not configured")
	}

	return cfg, nil
}
