package main

import (
	"context"
	"fmt"
	"strings"

	"os"
	"os/signal"
	"syscall"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/admin"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/logserver"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/factory"
)

func createDataLayer(config *manager.Config) (storage.DataLayer, error) {
	return factory.NewDataLayer(factory.DataLayerConfig{
		Type:        config.DataLayerType,
		DBPath:      config.DataLayerDBPath,
		NumVersions: config.DataLayerNumOfVersions,
	})
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

func main() {
	// Use a temporary logger for fatal error, used only during early startup
	logTmp := logger.NewLogger(&logger.Config{Kind: "zerolog", ConsoleLevel: "info", Console: true})
	defer logTmp.Close()

	// Load configuration
	config, err := manager.LoadConfig()
	if err != nil {
		logTmp.Fatal("Failed to load configuration: %v", err)
	}

	// Create a logger from config
	log := logger.NewLogger(&logger.Config{
		Kind:             config.LogKind,
		Console:          config.LogConsole,
		ConsoleLevel:     config.LogConsoleLevel,
		ConsoleColor:     config.LogConsoleColor,
		FileName:         config.LogFileName,
		FileLevel:        config.LogFileLevel,
		RemoteLogNetwork: "tcp",
		RemoteLogParams:  config.LogServerTCPAddress,
		NetworkLevel:     config.LogNetworkLevel,
	})
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing logger: %v\n", err)
		}
	}()

	log.Warn("Initializing manager...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the log server if configured
	err = logserver.StartLogServer(
		ctx,
		logserver.LogServerConfig{
			TCPAddr:        config.LogServerTCPAddress,
			VSockAddr:      config.LogServerVSockAddress,
			LogFilePath:    config.LogServerLogFile,
			ConsoleEnabled: config.LogServerConsole,
			ConsoleLevel:   config.LogServerConsoleLevel,
			FileLevel:      config.LogServerFileLevel,
		},
	)
	if err != nil {
		logTmp.Error("Failed to start the log server: %v", err)
		return
	}

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
	switch config.ChannelType {
	case "tcp":
		channelParams, ok := config.ChannelParams.(common.TcpChannelConnectionParams)
		if !ok {
			log.Error("Invalid channel parameters for tcp channel")
			return
		}
		factory := communication.NewTCPConnectionFactory(channelParams.Url())
		executorClient = communication.NewClient(factory, config.CommunicationParams, log)
	case "vsock":
		channelParams, ok := config.ChannelParams.(common.VSockChannelConnectionParams)
		if !ok {
			log.Error("Invalid channel parameters for vsock channel")
			return
		}
		factory := communication.NewVSockConnectionFactory(
			channelParams.CID,
			channelParams.Port,
		)
		executorClient = communication.NewClient(factory, config.CommunicationParams, log)
	default:
		log.Error("Unsupported channel type: %s", config.ChannelType)
		return
	}

	// Create the admin command server
	adminChannelParams := config.AdminChannelParams.(common.TcpChannelConnectionParams)
	adminFactory := communication.NewTCPConnectionFactory(adminChannelParams.Url())
	adminServer := admin.NewAdminServer(adminFactory, config.AdminCommunicationParams, log)

	// Create the manager
	secureProcessorManager := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient, adminServer, log)
	log.Info("Starting manager...")
	// Start the manager
	if err := secureProcessorManager.Start(ctx); err != nil {
		log.Error("Failed to start manager: %v", err)
		return
	}
	log.Info("Manager started")

	// Wait for shutdown signal or fatal error
	select {
	case <-sigChan:
		signal.Stop(sigChan)
		log.Info("Received shutdown signal. Shutting down gracefully...")
	case err := <-secureProcessorManager.FatalErrChan():
		signal.Stop(sigChan)
		log.Error("Fatal error occurred: %v. Shutting down gracefully...", err)
	}

	// Stop the manager
	cancel()
	if err := secureProcessorManager.Stop(); err != nil {
		log.Error("Failed to stop manager: %v", err)
	}
	log.Info("Manager stopped")
}
