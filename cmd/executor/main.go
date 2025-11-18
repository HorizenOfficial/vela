package main

import (
	"context"

	"github.com/horizen-pes/pkg/communication"

	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/wasm"
)

func main() {
	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Create the WASM runtime
	runtime := wasm.NewWasmtimeRuntime(log)

	// Create the appropriate server based on configuration
	var server communication.ExecutorServer
	switch config.ServerType {
	case "tcp":
		factory := communication.NewTCPConnectionFactory(config.ServerAddr)
		server = communication.NewServer(factory, log)
	case "vsock":
		factory := communication.NewVSockConnectionFactory(config.ServerCid, config.ServerPort)
		server = communication.NewServer(factory, log)
	default:
		log.Fatal("Unsupported server type: %s", config.ServerType)
	}

	// Create the executor
	exec, err := executor.NewStatelessExecutor(config, runtime, server, log)
	if err != nil {
		log.Fatal("Error creating executor: %v", err)
	}

	// Start the executor
	log.Info("Starting executor service...")
	if err := exec.Start(ctx); err != nil {
		log.Fatal("Error starting executor: %v", err)
	}
	log.Info("Executor started")

	// Wait for the context to be canceled
	<-ctx.Done()

	// Stop the executor
	log.Info("Stopping executor service...")
	if err := exec.Close(); err != nil {
		log.Error("Error stopping executor: %v", err)
	}
	log.Info("Executor service stopped")
}
