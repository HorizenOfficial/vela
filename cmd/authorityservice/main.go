package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/horizen-pes/pkg/authorityservice"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := authorityservice.ReadConfig()

	if strings.TrimSpace(cfg.ReportsPath) == "" {
		log.Fatalf("Reports path is not configured")
	}

	svc, err := authorityservice.NewAuthorityService(cfg.ChainID, time.Duration(cfg.NonceTTLSeconds)*time.Second, cfg.ReportsPath)
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
