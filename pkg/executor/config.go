// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

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

const confFileName = "executor.conf"

// DefaultConfig returns the default configuration (possibly overridden by env variables)
func DefaultConfig() *Config {
	channelType := os.Getenv("CHANNEL_TYPE")
	if channelType == "" {
		channelType = "vsock"
	}
	executorServerPort, err := strconv.ParseUint(os.Getenv("EXECUTOR_PORT"), 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert EXECUTOR_PORT for error %v, using default value\n", err)
		executorServerPort = 4000
	}

	var channelConnectionParams common.ChannelConnectionParams
	if channelType == "vsock" {
		//note: CID is always 3 inside AWS-Nitro, and since we are not planning to use other TEE it is hard-coded
		channelConnectionParams = common.VSockChannelConnectionParams{CID: 3, Port: uint32(executorServerPort)}
	} else {
		executorIpAddress := os.Getenv("EXECUTOR_IP_ADDRESS")
		if executorIpAddress == "" {
			executorIpAddress = "localhost"
		}
		channelConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpAddress, Port: uint32(executorServerPort)}
	}

	recType, err := strconv.Atoi(os.Getenv("EXECUTOR_KEYSET_RECOVERY_TYPE"))
	if err != nil {
		fmt.Printf("Failed to convert EXECUTOR_KEYSET_RECOVERY_TYPE for error %v, using default value\n", err)
		recType = 0
	}

	fuelPrice := big.NewInt(1)
	minFeePerRequest := big.NewInt(5)

	if fuelStr := os.Getenv("EXECUTOR_FUEL_PRICE_PER_UNIT"); fuelStr != "" {
		if v, ok := new(big.Int).SetString(fuelStr, 10); ok && v.Sign() >= 0 {
			fuelPrice = v
		} else {
			fmt.Printf("Failed to parse EXECUTOR_FUEL_PRICE_PER_UNIT=%q, using default %s\n", fuelStr, fuelPrice.String())
		}
	}

	if minFeeStr := os.Getenv("EXECUTOR_MIN_FEE_PER_REQUEST"); minFeeStr != "" {
		if v, ok := new(big.Int).SetString(minFeeStr, 10); ok && v.Sign() >= 0 {
			minFeePerRequest = v
		} else {
			fmt.Printf("Failed to parse EXECUTOR_MIN_FEE_PER_REQUEST=%q, using default %s\n", minFeeStr, minFeePerRequest.String())
		}
	}

	logConsole, err := strconv.ParseBool(os.Getenv("EXECUTOR_LOG_CONSOLE"))
	if err != nil {
		logConsole = true
	}

	logConsoleLevel := os.Getenv("EXECUTOR_LOG_CONSOLE_LEVEL")
	if logConsoleLevel == "" {
		logConsoleLevel = "info"
	}

	logFileName := os.Getenv("EXECUTOR_LOG_FILE_NAME")
	if logFileName == "" {
		logFileName = ""
	}

	logFileLevel := os.Getenv("EXECUTOR_LOG_FILE_LEVEL")
	if logFileLevel == "" {
		if logFileName != "" {
			logFileLevel = "info"
		}
	}

	logConsoleColor, err := strconv.ParseBool(os.Getenv("EXECUTOR_LOG_CONSOLE_COLOR"))
	if err != nil {
		logConsoleColor = false
	}

	return &Config{
		ChannelType:        channelType,
		ChannelParams:      channelConnectionParams,
		KeySetRecoveryType: recType,
		FuelPricePerUnit:   fuelPrice,
		MinFeePerRequest:   minFeePerRequest,
		LogConsole:         logConsole,
		LogConsoleLevel:    logConsoleLevel,
		LogConsoleColor:    logConsoleColor,
		LogFileName:        logFileName,
		LogFileLevel:       logFileLevel,
	}
}

func LoadConfigFromFile() (*Config, error) {

	if !common.FileExists(confFileName) {
		log.Printf("File %s not found. Using default configuration for Executor.\n", confFileName)
		return DefaultConfig(), nil
	}
	// Load properties from file
	config, err := properties.LoadFile(confFileName, properties.UTF8)
	if err != nil {
		return nil, err
	}

	var channelType = config.MustGetString("ChannelType")
	var channelConnectionParams common.ChannelConnectionParams
	executorServerPort, err := strconv.ParseUint(config.GetString("ExecutorPort", "4000"), 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert ExecutorPort for error %v, using default value\n", err)
		executorServerPort = 4000
	}
	if channelType == "vsock" {
		channelConnectionParams = common.VSockChannelConnectionParams{CID: 3, Port: uint32(executorServerPort)}
	} else {
		executorIpAddress := config.GetString("ExecutorIpAddress", "localhost")
		channelConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpAddress, Port: uint32(executorServerPort)}
	}

	return &Config{
		ChannelType:        config.MustGetString("ChannelType"),
		ChannelParams:      channelConnectionParams,
		KeySetRecoveryType: config.GetInt("KeySetRecoveryType", 0),
		FuelPricePerUnit:   big.NewInt(int64(config.GetInt("FuelPricePerUnit", 1))),
		MinFeePerRequest:   big.NewInt(int64(config.GetInt("MinFeePerRequest", 5))),
		LogConsole:         config.GetBool("LogConsole", true),
		LogConsoleLevel:    config.GetString("LogConsoleLevel", "info"),
		LogConsoleColor:    config.GetBool("LogConsoleColor", true),
		LogFileName:        config.GetString("LogFileName", ""),
		LogFileLevel:       config.GetString("LogFileLevel", "info"),
	}, nil
}
