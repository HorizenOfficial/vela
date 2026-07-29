package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/crypto"
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
	// BlockchainConnectTimeout is the max time to wait for the dial and initial RPC handshake (ChainID) when connecting (in seconds)
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
	// DataLayerNumOfVersions specifies how many historical versions to keep per application.
	// Each app's version history is pruned independently when it exceeds this limit.
	// Only used by "versioned_leveldb".
	DataLayerNumOfVersions int

	// DeanonymizationReportPath is the path to a folder where to store deanonymization reports.
	DeanonymizationReportPath string

	// ArtifactsPath is the shared filesystem path where deploy artifacts are stored.
	ArtifactsPath string

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
	executorServerPort := common.GetConfigVarUint32("EXECUTOR_PORT", 4000, fileProperties)
	logServerPort := common.GetConfigVarUint32("LOG_SERVER_PORT", 5000, fileProperties)

	// channel communication between manager and executor
	var channelConnectionParams common.ChannelConnectionParams

	// log server Vsock address, initialized only in case channel type is Vsock
	var logServerVsockAddress common.VSockChannelConnectionParams

	if channelType == "vsock" {
		// CID >= 16 are available values for EC2 enclaves (where executor runs)
		// CID and port are both used when connecting to a server
		executorServerCid := common.GetConfigVarUint32("EXECUTOR_VSOCK_CID", 20, fileProperties)
		channelConnectionParams = common.VSockChannelConnectionParams{CID: executorServerCid, Port: executorServerPort}
		// if channel is vsock it means we also have a vsock connection used by executor for logging, we use of course a separate port
		// CID is not used actually when creating a lstening server
		logServerVsockAddress = common.VSockChannelConnectionParams{Port: logServerPort}
	} else {
		executorIpHost := common.GetConfigVar("EXECUTOR_IP_HOST", "localhost", fileProperties)
		channelConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpHost, Port: executorServerPort}
	}

	// TCP connection on log server, typically used for manager logs
	logServerTcpHost := common.GetConfigVar("LOG_SERVER_IP_HOST", "localhost", fileProperties)
	logServerTcpAddress := common.TcpChannelConnectionParams{Ip: logServerTcpHost, Port: logServerPort}

	var privateKey *cryptotypes.PrivateKeySecp256k1
	privateKeyFromEnv := os.Getenv("MANAGER_KEY_SECP256")
	var keyErr error
	if privateKeyFromEnv == "" {
		privateKey, keyErr = crypto.GeneratePrivateKeySecp256k1()
	} else {
		privateKey, keyErr = crypto.ImportPrivateKeySecp256k1FromHex(privateKeyFromEnv)
	}
	if keyErr != nil {
		return nil, fmt.Errorf("failed to load MANAGER_KEY_SECP256: %w", keyErr)
	}

	communicationParams := common.CommunicationParams{
		RequestTimeoutSec: time.Duration(common.GetConfigVarInt64("MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC", 30, fileProperties)),
	}

	// Admin command server configuration
	adminServerPort := common.GetConfigVarUint32("MANAGER_ADMIN_PORT", 4002, fileProperties)
	var adminChannelConnectionParams common.ChannelConnectionParams
	// Admin server always uses TCP for external access.
	// MANAGER_ADMIN_SERVER_HOST controls the bind address for the admin server.
	// Defaults to "0.0.0.0" so admincli can connect from outside the container.
	adminServerHost := common.GetConfigVar("MANAGER_ADMIN_SERVER_HOST", "0.0.0.0", fileProperties)
	adminChannelConnectionParams = common.TcpChannelConnectionParams{Ip: adminServerHost, Port: adminServerPort}

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
		DeanonymizationReportPath: common.GetConfigVar("MANAGER_REPORTS_FOLDER", "/tmp/vela-data/manager_reports", fileProperties),
		ArtifactsPath:             common.GetConfigVar("MANAGER_ARTIFACTS_PATH", "", fileProperties),
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

	cfg.ArtifactsPath = strings.TrimSpace(cfg.ArtifactsPath)
	if cfg.ArtifactsPath == "" {
		return nil, fmt.Errorf("MANAGER_ARTIFACTS_PATH not configured")
	}

	return cfg, nil
}

// Validate checks the manager configuration for resource collisions that would
// cause silent failures, data corruption, or cryptic errors at runtime.
// It is scoped to resources owned by the manager process only.
func (c *Config) Validate() error {
	var errs []string

	// --- Channel type ---
	// Only "tcp" and "vsock" are supported. Validate() is called before any
	// resources are created, so catching this early avoids starting the log
	// server, blockchain client, and other heavy initialization only to fail
	// in the ChannelType switch.
	if c.ChannelType != "tcp" && c.ChannelType != "vsock" {
		errs = append(errs, fmt.Sprintf(
			"CHANNEL_TYPE must be \"tcp\" or \"vsock\", got %q", c.ChannelType))
	}

	// --- TCP port uniqueness ---
	// The admin server (MANAGER_ADMIN_PORT) and the log server (LOG_SERVER_PORT)
	// both create TCP listeners on the manager host. If they bind the same
	// address, the second net.Listen call gets "address already in use".
	// Depending on startup order, either the admin CLI becomes unreachable or the
	// executor loses its log sink — with only a warning in the logs that is easy
	// to miss.
	adminAddr, adminOk := c.AdminChannelParams.(common.TcpChannelConnectionParams)
	logServerAddr := c.LogServerTCPAddress
	if adminOk && adminAddr.Url() == logServerAddr.Url() {
		errs = append(errs, fmt.Sprintf(
			"MANAGER_ADMIN_PORT and LOG_SERVER_PORT resolve to the same address (%s); they must be different",
			adminAddr.Url()))
	}

	// --- Timeout values ---
	// HandshakeTimeout feeds time.After(Duration * time.Second). A zero or
	// negative value fires the timer immediately, so the executor handshake
	// always times out before it can complete.
	if c.HandshakeTimeout <= 0 {
		errs = append(errs, fmt.Sprintf(
			"HANDSHAKE_TIMEOUT must be > 0 (seconds), got %d", c.HandshakeTimeout))
	}

	// BlockchainPollingInterval feeds time.NewTicker(Duration * time.Second).
	// A zero or negative value panics the Go runtime ("non-positive interval
	// for NewTicker"), crashing the process.
	if c.BlockchainPollingInterval <= 0 {
		errs = append(errs, fmt.Sprintf(
			"BLOCKCHAIN_POLLING_INTERVAL must be > 0 (seconds), got %d", c.BlockchainPollingInterval))
	}

	// ReorgTimeout feeds time.Now().Add(Duration * time.Second). A zero value
	// means the reorg tolerance window expires instantly — any chain
	// reorganization triggers an immediate failure instead of waiting for
	// resolution. A negative value has the same effect.
	if c.ReorgTimeout <= 0 {
		errs = append(errs, fmt.Sprintf(
			"REORG_TIMEOUT must be > 0 (seconds), got %d", c.ReorgTimeout))
	}

	// BlockchainConnectTimeout == 0 is legal: the blockchain client falls back
	// to its own defaultConnectTimeout (10s). However, a negative value converts
	// to a negative time.Duration, producing a context that is already expired
	// the moment it is created — the ChainID RPC call always fails.
	if c.BlockchainConnectTimeout < 0 {
		errs = append(errs, fmt.Sprintf(
			"BLOCKCHAIN_CONNECT_TIMEOUT must be >= 0 (seconds, 0 = default 10s), got %d", c.BlockchainConnectTimeout))
	}

	// RequestTimeoutSec feeds time.Now().Add(Duration * time.Second) for pending
	// request deadlines. A zero value means every request's deadline is already
	// in the past at creation time — the response-routing loop evicts it
	// immediately, so every request fails with a timeout error.
	if c.CommunicationParams.RequestTimeoutSec <= 0 {
		errs = append(errs, fmt.Sprintf(
			"MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC must be > 0 (seconds), got %d",
			c.CommunicationParams.RequestTimeoutSec))
	}
	if c.AdminCommunicationParams.RequestTimeoutSec <= 0 {
		errs = append(errs, fmt.Sprintf(
			"MANAGER_ADMIN_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC must be > 0 (seconds), got %d",
			c.AdminCommunicationParams.RequestTimeoutSec))
	}

	// --- Directory uniqueness ---

	// LevelDB (MANAGER_DATA_FOLDER) manages its directory exclusively, creating
	// MANIFEST, *.ldb, and LOCK files. Deanonymization reports (MANAGER_REPORTS_FOLDER)
	// are written as .json files. If both point to the same directory, cleanup scripts
	// or backup tooling targeting one set of files can destroy the other — for example,
	// wiping the "reports folder" would delete the LevelDB database.
	if c.DataLayerDBPath != "" && c.DeanonymizationReportPath != "" {
		dbAbs, err1 := filepath.Abs(c.DataLayerDBPath)
		reportsAbs, err2 := filepath.Abs(c.DeanonymizationReportPath)
		if err := errors.Join(err1, err2); err != nil {
			errs = append(errs, fmt.Sprintf(
				"failed to resolve absolute paths for MANAGER_DATA_FOLDER / MANAGER_REPORTS_FOLDER: %v", err))
		} else if dbAbs == reportsAbs {
			errs = append(errs, fmt.Sprintf(
				"MANAGER_DATA_FOLDER and MANAGER_REPORTS_FOLDER resolve to the same path (%s); "+
					"LevelDB files and report JSON files must not share a directory", dbAbs))
		}
	}

	// Artifact files (MANAGER_ARTIFACTS_PATH) mixed with LevelDB files carry a
	// similar risk: any directory-level operation (rm, rsync, backup) on one becomes
	// hazardous to the other.
	if c.DataLayerDBPath != "" && c.ArtifactsPath != "" {
		dbAbs, err1 := filepath.Abs(c.DataLayerDBPath)
		artifactsAbs, err2 := filepath.Abs(c.ArtifactsPath)
		if err := errors.Join(err1, err2); err != nil {
			errs = append(errs, fmt.Sprintf(
				"failed to resolve absolute paths for MANAGER_DATA_FOLDER / MANAGER_ARTIFACTS_PATH: %v", err))
		} else if dbAbs == artifactsAbs {
			errs = append(errs, fmt.Sprintf(
				"MANAGER_DATA_FOLDER and MANAGER_ARTIFACTS_PATH resolve to the same path (%s); "+
					"LevelDB files and artifact files must not share a directory", dbAbs))
		}
	}

	// --- Log file uniqueness ---
	// The manager log (MANAGER_LOG_FILE_NAME) and the log server
	// (LOG_SERVER_FILE_NAME) use two independent writers — zerolog and the log
	// server's writer (optionally with lumberjack rotation). If both open the same
	// file, each does independent buffered writes and flushes, producing interleaved
	// bytes and corrupted output. When rotation is enabled, lumberjack
	// renames/truncates the file while zerolog still holds an open fd, causing
	// further data loss.
	if c.LogFileName != "" && c.LogServerLogFile != "" {
		logAbs, err1 := filepath.Abs(c.LogFileName)
		logServerAbs, err2 := filepath.Abs(c.LogServerLogFile)
		if err := errors.Join(err1, err2); err != nil {
			errs = append(errs, fmt.Sprintf(
				"failed to resolve absolute paths for MANAGER_LOG_FILE_NAME / LOG_SERVER_FILE_NAME: %v", err))
		} else if logAbs == logServerAbs {
			errs = append(errs, fmt.Sprintf(
				"MANAGER_LOG_FILE_NAME and LOG_SERVER_FILE_NAME resolve to the same file (%s); "+
					"concurrent writers will corrupt output", logAbs))
		}
	}

	// --- Log rotation without log file ---
	// If LOG_SERVER_FILE_ROTATION is enabled but LOG_SERVER_FILE_NAME is empty,
	// the rotation settings are silently ignored — the log server skips file
	// creation entirely. This likely indicates the operator intended to write
	// rotated logs but forgot to set the file path.
	if c.LogServerRotationEnabled && c.LogServerLogFile == "" {
		errs = append(errs, "LOG_SERVER_FILE_ROTATION is enabled but LOG_SERVER_FILE_NAME is empty; "+
			"rotation has no effect without a log file path")
	}

	if len(errs) > 0 {
		return fmt.Errorf("configuration conflicts:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}
