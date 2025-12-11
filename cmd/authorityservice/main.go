package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/horizen-pes/pkg/authorityservice"
	"github.com/horizen-pes/pkg/logger"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := authorityservice.LoadConfig()
	if err != nil {
		// Use a temporary logger for fatal error
		log := logger.NewLogger(&logger.Config{Kind: "zerolog", ConsoleLevel: "info", Console: true})
		log.Fatal("Failed to load configuration: %v", err)
	}

	// Create a logger from config
	log := logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      cfg.LogConsole,
		ConsoleLevel: cfg.LogConsoleLevel,
		ConsoleColor: cfg.LogConsoleColor,
		FileName:     cfg.LogFileName,
		FileLevel:    cfg.LogFileLevel,
	})
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing logger: %v\n", err)
		}
	}()

	if strings.TrimSpace(cfg.ReportsPath) == "" {
		log.Error("Reports path is not configured")
		return
	}

	svc, err := authorityservice.NewAuthorityService(cfg.ChainID, time.Duration(cfg.NonceTTLSeconds)*time.Second, cfg.ReportsPath, log)
	if err != nil {
		log.Error("Failed to create authority service: %v", err)
		return
	}

	server := authorityservice.NewHTTPServer(cfg, svc.Handler())

	go func() {
		if err := server.Start(); err != nil {
			log.Error("Failed to start HTTP server: %v", err)
			return
		}
	}()

	log.Info("Authority service listening on %s", cfg.ListenAddress)

	<-sigChan
	signal.Stop(sigChan)
	cancel()
	_ = server.Stop(ctx)
	log.Info("Authority service stopped")
}
