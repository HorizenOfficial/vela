package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/horizen-pes/pkg/communication"

	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/wasm"
)

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)


	// Create the executor configuration
	config, err := executor.LoadConfigFromFile()
	if err != nil {
		// Use a temporary logger for fatal error
		log := logger.NewLogger(&logger.Config{Kind: "zerolog", ConsoleLevel: "info", Console: true})
		log.Fatal("Failed to load configuration: %v", err)
	}

	// Create a logger from config
	log := logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      config.LogConsole,
		ConsoleLevel: config.LogConsoleLevel,
		ConsoleColor: config.LogConsoleColor,
		FileName:     config.LogFileName,
		FileLevel:    config.LogFileLevel,
	})
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing logger: %v\n", err)
		}
	}()


	// Create the WASM runtime
	runtime := wasm.NewWasmtimeRuntime(log)

	// Create the appropriate server based on configuration
	var factory communication.ConnectionFactory
	var adminFactory communication.ConnectionFactory
	switch config.ServerType {
	case "tcp":
		factory = communication.NewTCPConnectionFactory(config.ServerAddr)
		adminFactory = communication.NewTCPConnectionFactory(config.AdminServerAddr)
	case "vsock":
		factory = communication.NewVSockConnectionFactory(config.ServerCid, config.ServerPort)
		adminFactory = communication.NewVSockConnectionFactory(config.ServerCid, config.AdminServerPort)
	default:
		log.Error("Unsupported server type: %s", config.ServerType)
		return
	}

	server := communication.NewServer(factory, log)
	adminServer := communication.NewAdminServer(adminFactory, log)

	// Create the executor
	exec, err := executor.NewStatelessExecutor(config, runtime, server, adminServer, log)
	if err != nil {
		log.Error("Error creating executor: %v", err)
		return
	}

	// Start the executor
	log.Info("Starting executor service...")
	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	if err := exec.Start(ctx); err != nil {
		log.Error("Error starting executor: %v", err)
		return
	}
	log.Info("Executor started")

	// Wait for shutdown signal
	<-sigChan
	signal.Stop(sigChan)
	// Handle shutdown signal (Ctrl+C or SIGTERM)
	log.Info("Received shutdown signal. Shutting down gracefully...")
	
	// Stop the executor
	cancel()

	log.Info("Stopping executor service...")
	if err := exec.Close(); err != nil {
		log.Error("Error stopping executor: %v", err)
	}
	log.Info("Executor service stopped")
}
