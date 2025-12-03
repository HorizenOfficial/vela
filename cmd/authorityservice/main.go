package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/horizen-pes/pkg/authorityservice"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/factory"
)

func createDataLayer(config *authorityservice.Config) (storage.DataLayer, error) {
	return factory.NewDataLayer(factory.DataLayerConfig{
		Type:        config.DataLayerType,
		DBPath:      config.DataLayerDBPath,
		NumVersions: config.DataLayerNumOfVersions,
	})
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
