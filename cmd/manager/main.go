package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/mockdb"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

func createDataLayer(config *manager.Config) (storage.ApplicationStateStore, error) {
	// first of all if we are using mockdb do not even care for file and other configs
	if config.DataLayerType == "mockdb" {
		return mockdb.NewMockDataLayer(), nil
	}

	if strings.TrimSpace(config.DataLayerDBPath) == "" {
		return nil, fmt.Errorf("data layer path is empty")
	}

	var dl storage.ApplicationStateStore
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
	case "boltdb":
		return nil, fmt.Errorf("boltdb data layer creation is not yet implemented")
	default:
		return nil, fmt.Errorf("unknown data layer type: %s", config.DataLayerType)
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
	dataLayer, err := createDataLayer(config)
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
