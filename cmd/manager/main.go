package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"strconv"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/crypto"
	"os"
    "os/signal"
	"syscall"
	boltDb "github.com/horizen-pes/pkg/storage/boltdb"
	"github.com/horizen-pes/pkg/storage/mockdb"

	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
	commonEth "github.com/ethereum/go-ethereum/common"
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
		cfg := boltDb.BoltDBConfig{
			Path:    config.DataLayerDBPath,
			Timeout: 1 * time.Second,
		}
		dl, err = boltDb.NewBoltDBDataLayer(cfg)
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
	if strings.TrimSpace(config.KeyRegistryAddress) == "" {
		return nil, fmt.Errorf("key registry address is empty")
	}
	if strings.TrimSpace(config.RpcURL) == "" {
		return nil, fmt.Errorf("rpc url is empty")
	}
	if strings.TrimSpace(config.PrivateKey) == "" {
		return nil, fmt.Errorf("private key is empty")
	}

	privKey, err := crypto.ImportPrivateKeySecp256k1FromHex(config.PrivateKey)
    if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	} 

	if !commonEth.IsHexAddress(config.ProcessorAddress){
		return nil, fmt.Errorf("processor address is not a valid hex address")
	}

	if !commonEth.IsHexAddress(config.KeyRegistryAddress){
		return nil, fmt.Errorf("keyregistry address is not a valid hex address")
	}

	bcClient := blockchain.NewBlockChainClient(
		commonEth.HexToAddress(config.ProcessorAddress), 
		commonEth.HexToAddress(config.KeyRegistryAddress), 
		config.RpcURL, 
		privKey)

	return bcClient, nil
}	

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the manager configuration
	config := manager.DefaultConfig()

	// Create the blockchain client
	blockchainClient, err := createBlockchainClient(config)
	if err != nil {
		log.Fatalf("Failed to create blockchain client: %v", err)
	}

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
