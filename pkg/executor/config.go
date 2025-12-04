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
	// ServerType is the type of server to use (tcp / vsock)
	ServerType string
	// ServerAddr is the address for the TCP server
	ServerAddr string
	// ServerCid is the cid for the v-socket server
	ServerCid uint32
	// ServerPort is the port for the v-socket server
	ServerPort uint32
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
	serverType := os.Getenv("EXECUTOR_SERVER_TYPE")
	if serverType == "" {
		serverType = "tcp"
	}

	serverAddress := os.Getenv("EXECUTOR_IP_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost"
	}

	serverCid, err := strconv.ParseUint(os.Getenv("EXECUTOR_VSOCK_CID"), 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert EXECUTOR_VSOCK_CID for error %v, using default value\n", err)
		serverCid = 0
	}

	serverPort, err := strconv.ParseUint(os.Getenv("EXECUTOR_IP_PORT"), 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert EXECUTOR_IP_PORT for error %v, using default value\n", err)
		serverPort = 8080
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
		ServerType:         serverType,
		ServerAddr:         serverAddress + ":" + strconv.FormatUint(uint64(serverPort), 10),
		ServerCid:          uint32(serverCid),
		ServerPort:         uint32(serverPort),
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

	return &Config{
		ServerType:         config.MustGetString("ServerType"),
		ServerAddr:         config.MustGetString("ServerAddr"),
		ServerCid:          config.GetUint32("ServerCid", 2),
		ServerPort:         config.GetUint32("ServerPort", 54321),
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
