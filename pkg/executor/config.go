// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"fmt"
	"math/big"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/magiconair/properties"
)

// Config defines the configuration for the executor application
type Config struct {
	// ChannelType is the type of communication channel to use between manager and executor (tcp / vsock)
	ChannelType string
	// ChannelParams are the parameters for the connection server
	ChannelParams common.ChannelConnectionParams

	// AdminChannelParams are the parameters for the admin server
	AdminChannelParams common.ChannelConnectionParams

	// KeySetRecoveryType is the type of recovery mechanism to use for the keyset
	// Type 0: Unsafe (development only) - masterKey stored in plaintext
	// Type 1: AWS KMS - masterKey encrypted with KMS using Nitro attestation
	KeySetRecoveryType int

	// KMS Configuration (used when KeySetRecoveryType == 1)
	// KMSEnabled enables Type 1 KMS-based key management.
	// When true, KeySetRecoveryType is automatically set to 1.
	KMSEnabled bool
	// KMSKeyARN is the ARN of the AWS KMS key used for envelope encryption.
	// Required when KMSEnabled is true.
	KMSKeyARN string
	// KMSRegion is the AWS region where the KMS key is located.
	// Defaults to "us-east-1" if not specified.
	KMSRegion string
	// KMSProxyPort is the vsock port where the KMS proxy listens on the parent EC2.
	// Inside a Nitro Enclave, KMS requests are forwarded through this proxy.
	// Defaults to 8000 if not specified.
	KMSProxyPort uint32
	// FuelPricePerUnit is the price of fuel per unit
	FuelPricePerUnit *big.Int
	// MinFeePerRquest is the minimum fee to be paid for each request
	MinFeePerRequest *big.Int

	// LogKind is the type of logger to use (e.g., "zeronetwork", "tcplog", "zerolog").
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

	// LogChannelParams are the parameters for sending logging to a remote server.
	LogChannelParams common.ChannelConnectionParams

	// LogNetworkLevel is the level of logging for the network
	LogNetworkLevel string

	// CommunicationParams holds parameters for communication between manager and executor
	CommunicationParams common.CommunicationParams

	// AdminCommunicationParams holds parameters for communication towards the admin server
	AdminCommunicationParams common.CommunicationParams
}

const confFileName = "executor.conf"

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

	var channelServerConnectionParams common.ChannelConnectionParams
	executorServerPort := common.GetConfigVarInt64("EXECUTOR_PORT", 4000, fileProperties)
	var logClientConnectionParams common.ChannelConnectionParams
	logServerPort := common.GetConfigVarInt64("LOG_SERVER_PORT", 5000, fileProperties)
	var adminChannelConnectionParams common.ChannelConnectionParams
	adminServerPort := common.GetConfigVarInt64("EXECUTOR_ADMIN_PORT", 4001, fileProperties)
	if channelType == "vsock" {
		// CID 3 is reserved for the parent EC2 instance (where manager runs), CID >= 16 are available for enclaves (where executor runs)
		// CID is not used actually when creating a listening server
		channelServerConnectionParams = common.VSockChannelConnectionParams{Port: uint32(executorServerPort)}
		adminChannelConnectionParams = common.VSockChannelConnectionParams{Port: uint32(adminServerPort)}
		// CID and port are both used when connecting to a server
		managerCid := common.GetConfigVarInt64("MANAGER_VSOCK_CID", 3, fileProperties)
		logClientConnectionParams = common.VSockChannelConnectionParams{CID: uint32(managerCid), Port: uint32(logServerPort)}
	} else {
		executorIpHost := common.GetConfigVar("EXECUTOR_IP_HOST", "localhost", fileProperties)
		channelServerConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpHost, Port: uint32(executorServerPort)}
		adminChannelConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpHost, Port: uint32(adminServerPort)}
		logServerIpHost := common.GetConfigVar("LOG_SERVER_IP_HOST", "localhost", fileProperties)
		logClientConnectionParams = common.TcpChannelConnectionParams{Ip: logServerIpHost, Port: uint32(logServerPort)}
	}

	
	communicationParams := common.CommunicationParams{
		RequestTimeoutSec: time.Duration(common.GetConfigVarInt64("EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC", 30, fileProperties)),
	}
	
	adminCommunicationParams := common.CommunicationParams{
		RequestTimeoutSec: time.Duration(common.GetConfigVarInt64("ADMIN_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC", 30, fileProperties)),
	}
	
	
	
	// KMS Configuration
	kmsEnabled := common.GetConfigVarBool("EXECUTOR_KMS_ENABLED", false, fileProperties)
	kmsKeyARN := common.GetConfigVar("EXECUTOR_KMS_KEY_ARN", "", fileProperties)
	kmsRegion := common.GetConfigVar("EXECUTOR_KMS_REGION", "us-east-1", fileProperties)
	kmsProxyPort := uint32(common.GetConfigVarInt64("EXECUTOR_KMS_PROXY_PORT", 8000, fileProperties))

	// Determine KeySetRecoveryType: if KMS is enabled, force Type 1
	keySetRecoveryType := int(common.GetConfigVarInt64("EXECUTOR_KEYSET_RECOVERY_TYPE", 0, fileProperties))
	if kmsEnabled {
		keySetRecoveryType = 1
	}

	return &Config{
		ChannelType:              channelType,
		ChannelParams:            channelServerConnectionParams,
		AdminChannelParams:       adminChannelConnectionParams,
		KeySetRecoveryType:       keySetRecoveryType,
		KMSEnabled:               kmsEnabled,
		KMSKeyARN:                kmsKeyARN,
		KMSRegion:                kmsRegion,
		KMSProxyPort:             kmsProxyPort,
		FuelPricePerUnit:         big.NewInt(common.GetConfigVarInt64("EXECUTOR_FUEL_PRICE_PER_UNIT", 1, fileProperties)),
		MinFeePerRequest:         big.NewInt(common.GetConfigVarInt64("EXECUTOR_MIN_FEE_PER_REQUEST", 10, fileProperties)),
		LogKind:                  common.GetConfigVar("EXECUTOR_LOG_KIND", "zeronetwork", fileProperties),
		LogConsole:               common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE", true, fileProperties),
		LogConsoleLevel:          common.GetConfigVar("EXECUTOR_LOG_CONSOLE_LEVEL", "info", fileProperties),
		LogConsoleColor:          common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE_COLOR", false, fileProperties),
		LogFileName:              common.GetConfigVar("EXECUTOR_LOG_FILE_NAME", "", fileProperties),
		LogFileLevel:             common.GetConfigVar("EXECUTOR_LOG_FILE_LEVEL", "info", fileProperties),
		LogChannelParams:         logClientConnectionParams,
		LogNetworkLevel:          common.GetConfigVar("EXECUTOR_LOG_NETWORK_LEVEL", "info", fileProperties),
		CommunicationParams:      communicationParams,
		AdminCommunicationParams: adminCommunicationParams,
	}, nil
}

// ValidateKMSConfig validates KMS-related configuration.
// Returns an error if KMS is enabled but required settings are missing.
func (c *Config) ValidateKMSConfig() error {
	if c.KeySetRecoveryType == 1 || c.KMSEnabled {
		if c.KMSKeyARN == "" {
			return fmt.Errorf("KMS key ARN is required when KMS is enabled (EXECUTOR_KMS_KEY_ARN)")
		}
		if c.KMSRegion == "" {
			return fmt.Errorf("KMS region is required when KMS is enabled (EXECUTOR_KMS_REGION)")
		}
	}
	return nil
}
