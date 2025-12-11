// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"math/big"

	"github.com/horizen-pes/pkg/common"
	"github.com/magiconair/properties"
)

// Config defines the configuration for the executor application
type Config struct {
	// ChannelType is the type of communication channel to use between manager and executor (tcp / vsock)
	ChannelType string
	// ChannelParams are the parameters for the connection server
	ChannelParams common.ChannelConnectionParams

	// KeySetRecoveryType is the type of recovery mechanism to use for the keyset
	KeySetRecoveryType int
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
	var channelConnectionParams common.ChannelConnectionParams
	executorServerPort := common.GetConfigVarInt64("EXECUTOR_PORT", 4000, fileProperties)
	var logConnectionParams common.ChannelConnectionParams
	logExecutorServerPort := common.GetConfigVarInt64("EXECUTOR_LOG_PORT", 5000, fileProperties)
	if channelType == "vsock" {
		//note: CID is always 3 inside AWS-Nitro, and since we are not planning to use other TEE it is hard-coded
		channelConnectionParams = common.VSockChannelConnectionParams{CID: 3, Port: uint32(executorServerPort)}
		logConnectionParams = common.VSockChannelConnectionParams{CID: 3, Port: uint32(logExecutorServerPort)}
	} else {
		executorIpAddress := common.GetConfigVar("EXECUTOR_IP_ADDRESS", "localhost", fileProperties)
		channelConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpAddress, Port: uint32(executorServerPort)}
		logConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpAddress, Port: uint32(logExecutorServerPort)}
	}

	return &Config{
		ChannelType:        channelType,
		ChannelParams:      channelConnectionParams,
		KeySetRecoveryType: int(common.GetConfigVarInt64("EXECUTOR_KEYSET_RECOVERY_TYPE", 0, fileProperties)),
		FuelPricePerUnit:   big.NewInt(common.GetConfigVarInt64("EXECUTOR_FUEL_PRICE_PER_UNIT", 1, fileProperties)),
		MinFeePerRequest:   big.NewInt(common.GetConfigVarInt64("EXECUTOR_MIN_FEE_PER_REQUEST", 10, fileProperties)),
		LogKind:            common.GetConfigVar("EXECUTOR_LOG_KIND", "zeronetwork", fileProperties),
		LogConsole:         common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE", true, fileProperties),
		LogConsoleLevel:    common.GetConfigVar("EXECUTOR_LOG_CONSOLE_LEVEL", "info", fileProperties),
		LogConsoleColor:    common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE_COLOR", false, fileProperties),
		LogFileName:        common.GetConfigVar("EXECUTOR_LOG_FILE_NAME", "", fileProperties),
		LogFileLevel:       common.GetConfigVar("EXECUTOR_LOG_FILE_LEVEL", "info", fileProperties),
		LogChannelParams:   logConnectionParams,
		LogNetworkLevel:    common.GetConfigVar("EXECUTOR_LOG_NETWORK_LEVEL", "info", fileProperties),
	}, nil
}
