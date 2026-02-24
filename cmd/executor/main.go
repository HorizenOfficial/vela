package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/executor/kms"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/version"
	"github.com/horizen-pes/pkg/wasm"
)

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create the executor configuration
	config, err := executor.LoadConfig()
	if err != nil {
		// Use a temporary logger for fatal error
		log := logger.NewLogger(&logger.Config{Kind: "zerolog", ConsoleLevel: "info", Console: true})
		log.Fatal("Failed to load configuration: %v", err)
	}

	// Create context first so defer cancel() is registered before defer log.Close().
	// Defers run LIFO, so log.Close() will drain the async buffer to the log server
	// before cancel() shuts down context-dependent resources.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a logger from config
	log := logger.NewLogger(&logger.Config{
		Kind:             config.LogKind,
		Console:          config.LogConsole,
		ConsoleLevel:     config.LogConsoleLevel,
		ConsoleColor:     config.LogConsoleColor,
		FileName:         config.LogFileName,
		FileLevel:        config.LogFileLevel,
		RemoteLogNetwork: config.ChannelType,
		RemoteLogParams:  config.LogChannelParams,
		NetworkLevel:     config.LogNetworkLevel,
	})
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing logger: %v\n", err)
		}
	}()

	// Create the WASM runtime
	runtime := wasm.NewWasmtimeRuntime(log)

	// Create the appropriate server based on configuration
	var server communication.ExecutorServer
	switch config.ChannelType {
	case "tcp":
		factory := communication.NewTCPConnectionFactory(config.ChannelParams.(common.TcpChannelConnectionParams).Url())
		server = communication.NewServer(factory, config.CommunicationParams, log)
	case "vsock":
		factory := communication.NewVSockConnectionFactory(
			config.ChannelParams.(common.VSockChannelConnectionParams).CID,
			config.ChannelParams.(common.VSockChannelConnectionParams).Port,
		)
		server = communication.NewServer(factory, config.CommunicationParams, log)
	default:
		log.Error("Unsupported channel type: %s", config.ChannelType)
		return
	}


	// Initialize KMS dependencies if Type 1 is configured
	var kmsClient kms.KMSClient
	var enclaveHandle kms.EnclaveHandle

	if config.KeySetRecoveryType == common.RecoveryTypeKMS {
		// Validate KMS configuration
		if err := config.ValidateKMSConfig(); err != nil {
			log.Fatal("Invalid KMS configuration: %v", err)
		}

		log.Info("Initializing Type 1 (KMS) key recovery with Nitro Enclave attestation...")

		// Initialize enclave handle first - this will fail if not running in a Nitro Enclave
		nitroEnclave, err := kms.NewNitroEnclaveHandle()
		if err != nil {
			log.Fatal("Failed to initialize Nitro Enclave handle: %v. "+
				"Type 1 recovery requires running inside a Nitro Enclave. "+
				"For local development, use Type 0 (EXECUTOR_KEYSET_RECOVERY_TYPE=0).", err)
		}
		enclaveHandle = nitroEnclave
		log.Info("Nitro Enclave handle initialized successfully")

		// Initialize KMS client with vsock proxy
		nitroKMS, err := kms.NewNitroKMSClient(
			context.Background(),
			config.KMSRegion,
			config.KMSKeyARN,
			config.KMSProxyPort,
		)
		if err != nil {
			log.Fatal("Failed to initialize Nitro KMS client: %v", err)
		}
		kmsClient = nitroKMS
		log.Info("Nitro KMS client initialized (region=%s, proxy_port=%d)", config.KMSRegion, config.KMSProxyPort)
	}
	// For Type 0, kmsClient and enclaveHandle remain nil (which is valid)

	// Create the executor
	exec, err := executor.NewStatelessExecutor(config, runtime, server, log, kmsClient, enclaveHandle)
	if err != nil {
		log.Error("Error creating executor: %v", err)
		return
	}

	log.Info("Executor version: %s", version.Version)

	// Start the executor
	log.Info("Starting executor service...")
	if err := exec.Start(ctx); err != nil {
		log.Error("Error starting executor: %v", err)
		return
	}
	log.Info("Executor started")

	// Wait for shutdown signal
	<-sigChan
	signal.Stop(sigChan)
	// Handle shutdown signal (Ctrl+C or SIGTERM)
	log.Info("Received shutdown signal. Shutting down gracefully...")
	
	// Stop the executor
	cancel()

	log.Info("Stopping executor service...")
	if err := exec.Close(); err != nil {
		log.Error("Error stopping executor: %v", err)
	}
	log.Info("Executor service stopped")
}
