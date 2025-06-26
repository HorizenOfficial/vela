package main

import (
	"context"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"log"
)

func main() {
	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the blockchain client
	blockchainClient := blockchain.NewMockClient()

	// Create the data layer
	dataLayer := storage.NewMockDataLayer()

	// Create the manager configuration
	config := manager.DefaultConfig()

	// Create the manager
	secureProcessorManager, err := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// Start the manager
	if err := secureProcessorManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start manager: %v", err)
	}
	log.Println("Manager started")

	// Wait for the context to be canceled
	<-ctx.Done()

	// Stop the manager
	if err := secureProcessorManager.Stop(); err != nil {
		log.Fatalf("Failed to stop manager: %v", err)
	}
	log.Println("Manager stopped")
}
