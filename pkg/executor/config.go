// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"fmt"
	"log"
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
	serverAddress := os.Getenv("EXECUTOR_IP_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost"
	}
	serverPort := os.Getenv("EXECUTOR_IP_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	recTypeVar := os.Getenv("EXECUTOR_KEYSET_RECOVERY_TYPE")
	recType, err := strconv.Atoi(recTypeVar)
	if err != nil {
		fmt.Printf("Failed to convert EXECUTOR_KEYSET_RECOVERY_TYPE for error %v, using default value\n", err)
		recType = 0
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
		ServerType:         "tcp",
		ServerAddr:         serverAddress + ":" + serverPort,
		KeySetRecoveryType: recType,
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
		LogConsole:         config.GetBool("LogConsole", true),
		LogConsoleLevel:    config.GetString("LogConsoleLevel", "info"),
		LogConsoleColor:    config.GetBool("LogColor", true),
		LogFileName:        config.GetString("LogFileName", ""),
		LogFileLevel:       config.GetString("LogFileLevel", "info"),
	}, nil
}
