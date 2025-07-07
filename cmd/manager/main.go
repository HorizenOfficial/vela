package main

import (
	"context"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"log"
	"os"
    "os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the manager configuration
	config := manager.DefaultConfig()

	// Create the blockchain client
	blockchainClient := blockchain.NewMockClient()

	// Create the data layer
	dataLayer := storage.NewMockDataLayer()

	// Create the executor client
	var executorClient communication.ExecutorClient
	switch config.ExecutorConnectionType {
	case "tcp":
		executorClient = communication.NewTCPClient(config.ExecutorConnectionParams["url"])
	case "vsock":
		executorClient = communication.NewVSockClient(config.ExecutorConnectionParams["cid"], config.ExecutorConnectionParams["port"])
	default:
		log.Fatalf("Unsupported executor connection type: %s", config.ExecutorConnectionType)
	}

	// Create the manager
	secureProcessorManager, err := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// Start the manager
	if err := secureProcessorManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start manager: %v", err)
	}
	log.Println("Manager started")

	// Wait for shutdown signal
    <-sigChan
	signal.Stop(sigChan)
	// Handle shutdown signal (Ctrl+C or SIGTERM)
	log.Println("Received shutdown signal. Shutting down gracefully...")

	// Stop the manager
	cancel()
	if err := secureProcessorManager.Stop(); err != nil {
		log.Fatalf("Failed to stop manager: %v", err)
	}
	log.Println("Manager stopped")
}
