package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/horizen-pes/pkg/authorityservice"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/mockdb"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

func createDataLayer(config *authorityservice.Config) (storage.DataLayer, error) {
	// Mock data layer is used only for tests; for production we read real data.
	if config.DataLayerType == "mockdb" {
		return mockdb.NewMockDataLayer(), nil
	}

	if strings.TrimSpace(config.DataLayerDBPath) == "" {
		return nil, fmt.Errorf("data layer path is empty")
	}

	switch config.DataLayerType {
	case "versioned_leveldb":
		cfg := versionedDb.VersionedLevelDBConfig{
			DBPath:         config.DataLayerDBPath,
			VersionsToKeep: config.DataLayerNumOfVersions,
		}
		return versionedDb.NewVersionedLevelDBDataLayer(cfg)
	default:
		return nil, fmt.Errorf("unknown data layer type: %s", config.DataLayerType)
	}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := authorityservice.ReadConfig()

	dl, err := createDataLayer(cfg)
	if err != nil {
		log.Fatalf("Failed to create data layer: %v", err)
	}
	defer dl.Close()

	svc, err := authorityservice.NewAuthorityService(dl, cfg.ChainID, time.Duration(cfg.NonceTTLSeconds)*time.Second)
	if err != nil {
		log.Fatalf("Failed to create authority service: %v", err)
	}

	server := authorityservice.NewHTTPServer(cfg, svc.Handler())

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	log.Printf("Authority service listening on %s", cfg.ListenAddress)

	<-sigChan
	signal.Stop(sigChan)
	cancel()
	_ = server.Stop(ctx)
	log.Println("Authority service stopped")
}
