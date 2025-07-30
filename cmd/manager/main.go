package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/communication"
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
		factory := communication.NewTCPConnectionFactory(config.ExecutorConnectionParams["url"])
		executorClient = communication.NewClient(factory)
	case "vsock":
		cidStr, err := strconv.ParseUint(config.ExecutorConnectionParams["cid"], 10, 32)
		if err != nil {
			log.Fatalf("Failed to parse port: %v", err)
		}
		cid := uint32(cidStr)

		portStr, err := strconv.ParseUint(config.ExecutorConnectionParams["port"], 10, 32)
		if err != nil {
			log.Fatalf("Failed to parse executor connection parameters: %v", err)
		}
		port := uint32(portStr)

		factory := communication.NewVSockConnectionFactory(cid, port)
		executorClient = communication.NewClient(factory)
	default:
		log.Fatalf("Unsupported executor connection type: %s", config.ExecutorConnectionType)
	}

	// Create the manager
	secureProcessorManager := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient)

	// Start the manager
	log.Println("Starting manager...")
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