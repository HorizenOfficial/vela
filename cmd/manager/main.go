package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

func createDataLayer() (storage.ApplicationStateStore, error) {
	// Create a unique temporary directory. TODO use os.MkdirAll() in production
	testVersionedLevelDBBaseDir := "/tmp/"
	tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-data-layer-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Configure Versioned LevelDB to use a file within the temporary directory.
	cfg := versionedDb.VersionedLevelDBConfig{
		Path:           filepath.Join(tempDir, "manager.db"), // The actual .db file
		VersionsToKeep: 5,                                    // Keep a small number of versions for testing
	}
	dl, err := versionedDb.NewVersionedLevelDBDataLayer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create VersionedLevelDBDataLayer: %w", err)
	}
	return dl, nil
}

func main() {
	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the manager configuration
	config := manager.DefaultConfig()

	// Create the blockchain client
	blockchainClient := blockchain.NewMockClient()

	// Create the data layer
	dataLayer, err := createDataLayer()
	if err != nil {
		log.Fatalf("Failed to create data layer: %v", err)
	}

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

	// Wait for the context to be canceled
	<-ctx.Done()

	// Stop the manager
	if err := secureProcessorManager.Stop(); err != nil {
		log.Fatalf("Failed to stop manager: %v", err)
	}
	log.Println("Manager stopped")
}
