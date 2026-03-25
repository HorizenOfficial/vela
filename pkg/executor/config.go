// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/magiconair/properties"
)

// Config defines the configuration for the executor application
type Config struct {
	// ChannelType is the type of communication channel to use between manager and executor (tcp / vsock)
	ChannelType string
	// ChannelParams are the parameters for the connection server
	ChannelParams common.ChannelConnectionParams

	// KeySetRecoveryType is the type of recovery mechanism to use for the keyset
	// Type 0: Unsafe (development only) - masterKey stored in plaintext
	// Type 1: AWS KMS - masterKey encrypted with KMS using Nitro attestation
	KeySetRecoveryType common.RecoveryType

	// KMS Configuration (used when KeySetRecoveryType == RecoveryTypeKMS)
	// KMSKeyARN is the ARN of the AWS KMS key used for envelope encryption.
	// Required when KeySetRecoveryType is RecoveryTypeKMS.
	KMSKeyARN string
	// KMSRegion is the AWS region where the KMS key is located.
	// Defaults to "eu-west-1" if not specified.
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
	if channelType == "vsock" {
		// CID 3 is reserved for the parent EC2 instance (where manager runs), CID >= 16 are available for enclaves (where executor runs)
		// CID is not used actually when creating a listening server
		channelServerConnectionParams = common.VSockChannelConnectionParams{Port: uint32(executorServerPort)}
		// CID and port are both used when connecting to a server
		managerCid := common.GetConfigVarInt64("MANAGER_VSOCK_CID", 3, fileProperties)
		logClientConnectionParams = common.VSockChannelConnectionParams{CID: uint32(managerCid), Port: uint32(logServerPort)}
	} else {
		executorIpHost := common.GetConfigVar("EXECUTOR_IP_HOST", "localhost", fileProperties)
		channelServerConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpHost, Port: uint32(executorServerPort)}
		logServerIpHost := common.GetConfigVar("LOG_SERVER_IP_HOST", "localhost", fileProperties)
		logClientConnectionParams = common.TcpChannelConnectionParams{Ip: logServerIpHost, Port: uint32(logServerPort)}
	}

	communicationParams := common.CommunicationParams{
		RequestTimeoutSec: time.Duration(common.GetConfigVarInt64("EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC", 30, fileProperties)),
	}

	// KMS Configuration
	kmsKeyARN := common.GetConfigVar("EXECUTOR_KMS_KEY_ARN", "", fileProperties)
	kmsRegion := common.GetConfigVar("EXECUTOR_KMS_REGION", "eu-west-1", fileProperties)
	kmsProxyPort := uint32(common.GetConfigVarInt64("EXECUTOR_KMS_PROXY_PORT", 8000, fileProperties))

	keySetRecoveryType := common.RecoveryType(common.GetConfigVarInt64("EXECUTOR_KEYSET_RECOVERY_TYPE", int64(common.RecoveryTypeKMS), fileProperties))
	if keySetRecoveryType != common.RecoveryTypeUnsafe && keySetRecoveryType != common.RecoveryTypeKMS {
		return nil, fmt.Errorf("invalid EXECUTOR_KEYSET_RECOVERY_TYPE: %d", keySetRecoveryType)
	}

	return &Config{
		ChannelType:         channelType,
		ChannelParams:       channelServerConnectionParams,
		KeySetRecoveryType:  keySetRecoveryType,
		KMSKeyARN:           kmsKeyARN,
		KMSRegion:           kmsRegion,
		KMSProxyPort:        kmsProxyPort,
		FuelPricePerUnit:    big.NewInt(common.GetConfigVarInt64("EXECUTOR_FUEL_PRICE_PER_UNIT", 1, fileProperties)),
		MinFeePerRequest:    big.NewInt(common.GetConfigVarInt64("EXECUTOR_MIN_FEE_PER_REQUEST", 10, fileProperties)),
		LogKind:             common.GetConfigVar("EXECUTOR_LOG_KIND", "zeronetwork", fileProperties),
		LogConsole:          common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE", true, fileProperties),
		LogConsoleLevel:     common.GetConfigVar("EXECUTOR_LOG_CONSOLE_LEVEL", "info", fileProperties),
		LogConsoleColor:     common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE_COLOR", false, fileProperties),
		LogFileName:         common.GetConfigVar("EXECUTOR_LOG_FILE_NAME", "", fileProperties),
		LogFileLevel:        common.GetConfigVar("EXECUTOR_LOG_FILE_LEVEL", "info", fileProperties),
		LogChannelParams:    logClientConnectionParams,
		LogNetworkLevel:     common.GetConfigVar("EXECUTOR_LOG_NETWORK_LEVEL", "info", fileProperties),
		CommunicationParams: communicationParams,
	}, nil
}

// Validate checks the executor configuration for invalid or conflicting settings.
func (c *Config) Validate() error {
	var errs []string

	// --- Channel type ---
	// Only "tcp" and "vsock" are supported. Validate() is called before the
	// logger or WASM runtime are created, so catching this early avoids wasted
	// startup effort and prevents the unguarded type-assertions in the
	// ChannelType switch from panicking.
	if c.ChannelType != "tcp" && c.ChannelType != "vsock" {
		errs = append(errs, fmt.Sprintf(
			"CHANNEL_TYPE must be \"tcp\" or \"vsock\", got %q", c.ChannelType))
	}

	// --- Fee parameters ---
	// FuelPricePerUnit is used in Mul operations (e.g. totalFuel * FuelPricePerUnit)
	// to compute application fees. A nil value would panic. A zero value effectively
	// disables fuel metering — every fee computation returns 0 and then falls back
	// to MinFeePerRequest, meaning all requests cost the same flat fee regardless
	// of how much computation they consume. A negative value would produce negative
	// fees, which breaks the economic model.
	if c.FuelPricePerUnit == nil || c.FuelPricePerUnit.Sign() <= 0 {
		errs = append(errs, "EXECUTOR_FUEL_PRICE_PER_UNIT must be > 0")
	}

	// MinFeePerRequest is compared against req.MaxFeeValue to reject underfunded
	// requests, and used as a floor for the application fee. A nil value would
	// panic. A negative value means the fee check always passes regardless of what
	// the user sends — a security concern since the system would process requests
	// without collecting fees. Zero is legal: it means there is no minimum fee.
	if c.MinFeePerRequest == nil || c.MinFeePerRequest.Sign() < 0 {
		errs = append(errs, "EXECUTOR_MIN_FEE_PER_REQUEST must be >= 0")
	}

	// --- Communication timeout ---
	// RequestTimeoutSec feeds time.Now().Add(Duration * time.Second) for pending
	// request deadlines. A zero value means every request's deadline is already
	// in the past at creation time — the response-routing loop evicts it
	// immediately, so every request fails with a timeout error.
	if c.CommunicationParams.RequestTimeoutSec <= 0 {
		errs = append(errs, fmt.Sprintf(
			"EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC must be > 0 (seconds), got %d",
			c.CommunicationParams.RequestTimeoutSec))
	}

	// --- KMS configuration ---
	errs = append(errs, c.validateKMSConfig()...)

	if len(errs) > 0 {
		return fmt.Errorf("configuration errors:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

// validateKMSConfig validates KMS-related configuration.
// Returns all errors found so the caller can present them at once.
func (c *Config) validateKMSConfig() []string {
	if c.KeySetRecoveryType != common.RecoveryTypeKMS {
		return nil
	}
	var errs []string
	if c.KMSKeyARN == "" {
		errs = append(errs, "KMS key ARN is required when KMS is enabled (EXECUTOR_KMS_KEY_ARN)")
	}
	if c.KMSRegion == "" {
		errs = append(errs, "KMS region is required when KMS is enabled (EXECUTOR_KMS_REGION)")
	}
	return errs
}
