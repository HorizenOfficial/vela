package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/authorityservice"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/logger"
)

func createBlockchainClient(cfg *authorityservice.Config) (blockchain.Client, error) {
	if strings.TrimSpace(cfg.ProcessorAddress) == "" {
		return nil, fmt.Errorf("processor address is empty")
	}
	if strings.TrimSpace(cfg.RpcURL) == "" {
		return nil, fmt.Errorf("rpc url is empty")
	}

	if !ethCommon.IsHexAddress(cfg.ProcessorAddress) {
		return nil, fmt.Errorf("processor address is not a valid hex address")
	}

	return blockchain.NewReadOnlyBlockChainClient(
		ethCommon.HexToAddress(cfg.ProcessorAddress),
		cfg.RpcURL,
	), nil
}

func ensureChainID(ctx context.Context, bc blockchain.Client, expected uint64, rpcURL string) error {
	connectedChainID, err := bc.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve blockchain chain ID: %w", err)
	}
	if connectedChainID == nil || connectedChainID.Uint64() != expected {
		return fmt.Errorf("blockchain chain ID mismatch: expected %d got %v (rpc %s)", expected, connectedChainID, rpcURL)
	}
	return nil
}

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

	bc, err := createBlockchainClient(cfg)
	if err != nil {
		log.Error("Failed to create blockchain client: %v", err)
		return
	}
	if err := bc.Connect(ctx); err != nil {
		log.Error("Failed to connect blockchain client: %v", err)
		return
	}
	if err := ensureChainID(ctx, bc, cfg.ChainID, cfg.RpcURL); err != nil {
		log.Error("%v", err)
		return
	}
	defer func() {
		if err := bc.Close(); err != nil {
			log.Error("Failed to close blockchain client: %v", err)
		}
	}()

	svc, err := authorityservice.NewAuthorityService(cfg.ChainID, time.Duration(cfg.NonceTTLSeconds)*time.Second, cfg.ReportsPath, bc, log)
	if err != nil {
		log.Error("Failed to create authority service: %v", err)
		return
	}

	server := authorityservice.NewHTTPServer(cfg, svc.Handler())

	serverErr := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	log.Info("Authority service listening on %s", cfg.ListenAddress)

	select {
	case sig := <-sigChan:
		log.Info("Received signal %s, shutting down", sig)
	case err := <-serverErr:
		log.Error("HTTP server exited with error: %v", err)
	}
	signal.Stop(sigChan)
	cancel()
	if err := server.Stop(ctx); err != nil {
		log.Error("Failed to stop HTTP server: %v", err)
	}
	log.Info("Authority service stopped")
}
