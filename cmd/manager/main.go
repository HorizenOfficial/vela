package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"os"
	"os/signal"
	"syscall"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/mockdb"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

func createDataLayer(config *manager.Config) (storage.DataLayer, error) {
	// first of all if we are using mockdb do not even care for file and other configs
	if config.DataLayerType == "mockdb" {
		return mockdb.NewMockDataLayer(), nil
	}

	if strings.TrimSpace(config.DataLayerDBPath) == "" {
		return nil, fmt.Errorf("data layer path is empty")
	}

	var dl storage.DataLayer
	var err error

	switch config.DataLayerType {
	case "versioned_leveldb":
		cfg := versionedDb.VersionedLevelDBConfig{
			DBPath:         config.DataLayerDBPath,
			VersionsToKeep: config.DataLayerNumOfVersions,
		}
		dl, err = versionedDb.NewVersionedLevelDBDataLayer(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create VersionedLevelDBDataLayer: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown data layer type: %s", config.DataLayerType)
	}

	return dl, nil
}

func createBlockchainClient(config *manager.Config) (blockchain.Client, error) {
	if config.MockBlockChainClient {
		return blockchain.NewMockClient(), nil
	}

	if strings.TrimSpace(config.ProcessorAddress) == "" {
		return nil, fmt.Errorf("processor address is empty")
	}
	if strings.TrimSpace(config.RpcURL) == "" {
		return nil, fmt.Errorf("rpc url is empty")
	}

	if !ethCommon.IsHexAddress(config.ProcessorAddress) {
		return nil, fmt.Errorf("processor address is not a valid hex address")
	}

	if !ethCommon.IsHexAddress(config.TeeAuthAddress) {
		return nil, fmt.Errorf("teeauthenticator address is not a valid hex address")
	}
	bcClient := blockchain.NewBlockChainClient(
		ethCommon.HexToAddress(config.ProcessorAddress),
		ethCommon.HexToAddress(config.TeeAuthAddress),
		config.RpcURL,
		&config.PrivateKey)

	return bcClient, nil
}

func startLogServer(ctx context.Context, listenAddr string, appLogger logger.Logger) {
	if strings.TrimSpace(listenAddr) == "" {
		appLogger.Info("Log server address is empty, not starting log server.")
		return
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		appLogger.Error("Failed to start log server on %s: %v", listenAddr, err)
		return
	}
	appLogger.Info("Log server listening on %s", listenAddr)

	go func() {
		<-ctx.Done() // Wait for context cancellation
		appLogger.Info("Shutting down log server...")
		if err := listener.Close(); err != nil {
			appLogger.Error("Error closing log server listener: %v", err)
		}
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return // Context canceled, listener is likely closed
				default:
					appLogger.Error("Log server accept error: %v", err)
				}
				continue
			}

			go handleLogConnection(conn, appLogger)
		}
	}()
}

func handleLogConnection(conn net.Conn, appLogger logger.Logger) {
	defer func() {
		if err := conn.Close(); err != nil {
			appLogger.Error("Error closing log connection: %v", err)
		}
	}()

	appLogger.Info("New log connection from %s", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				appLogger.Error("Error reading from log connection %s: %v", conn.RemoteAddr(), err)
			}
			break
		}
		// Log the received message using the manager's logger
		appLogger.Info("Remote Log: %s", strings.TrimSpace(message))
	}
	appLogger.Info("Log connection from %s closed", conn.RemoteAddr())
}

func main() {
	// Load configuration
	config, err := manager.LoadConfigFromFile()
	if err != nil {
		// Use a temporary logger for fatal error
		log := logger.NewLogger(&logger.Config{Kind: "zerolog", ConsoleLevel: "info", Console: true})
		log.Fatal("Failed to load configuration: %v", err)
	}

	// Create a logger from config
	log := logger.NewLogger(&logger.Config{
		Kind:                config.LogKind,
		Console:             config.LogConsole,
		ConsoleLevel:        config.LogConsoleLevel,
		ConsoleColor:        config.LogConsoleColor,
		FileName:            config.LogFileName,
		FileLevel:           config.LogFileLevel,
		RemoteLogAddress: config.RemoteLogAddress,
		RemoteLogNetwork: config.RemoteLogNetwork,
		NetworkLevel:     config.NetworkLevel,
	})
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing logger: %v\n", err)
		}
	}()

	log.Info("Starting manager...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the log server if configured
	startLogServer(ctx, config.LogServerTCPAddress, log)

	// Create the blockchain client
	blockchainClient, err := createBlockchainClient(config)
	if err != nil {
		log.Error("Failed to create blockchain client: %v", err)
		return
	}

	// Create the data layer
	dataLayer, err := createDataLayer(config)
	if err != nil {
		log.Error("Failed to create data layer: %v", err)
		return
	}

	// Create the executor client
	var executorClient communication.ExecutorClient
	switch config.ExecutorConnectionType {
	case "tcp":
		if strings.TrimSpace(config.ExecutorConnectionParams["url"]) == "" {
			log.Error("Tcp url is empty")
			return
		}
		factory := communication.NewTCPConnectionFactory(config.ExecutorConnectionParams["url"])
		executorClient = communication.NewClient(factory, log)
	case "vsock":
		cidStr, err := strconv.ParseUint(config.ExecutorConnectionParams["cid"], 10, 32)
		if err != nil {
			log.Error("Failed to parse port: %v", err)
			return
		}
		cid := uint32(cidStr)

		portStr, err := strconv.ParseUint(config.ExecutorConnectionParams["port"], 10, 32)
		if err != nil {
			log.Error("Failed to parse executor connection parameters: %v", err)
			return
		}
		port := uint32(portStr)

		factory := communication.NewVSockConnectionFactory(cid, port)
		executorClient = communication.NewClient(factory, log)
	default:
		log.Error("Unsupported executor connection type: %s", config.ExecutorConnectionType)
		return
	}

	// Create the manager
	secureProcessorManager := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient, log)
	log.Info("Starting manager...")
	// Start the manager
	if err := secureProcessorManager.Start(ctx); err != nil {
		log.Error("Failed to start manager: %v", err)
	}
	log.Info("Manager started")

	// Wait for shutdown signal
	<-sigChan
	signal.Stop(sigChan)
	// Handle shutdown signal (Ctrl+C or SIGTERM)
	log.Info("Received shutdown signal. Shutting down gracefully...")

	// Stop the manager
	cancel()
	if err := secureProcessorManager.Stop(); err != nil {
		log.Error("Failed to stop manager: %v", err)
	}
	log.Info("Manager stopped")
}
