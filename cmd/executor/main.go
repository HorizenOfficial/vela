package main

import (
	"context"
	"github.com/horizen-pes/pkg/communication"
	"log"

	"github.com/horizen-pes/pkg/executor"
)

func main() {
	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the executor configuration
	config := executor.DefaultConfig()

	// Create the WASM runtime
	runtime := executor.NewMockRuntime()

	// Create the appropriate server based on configuration
	var server communication.ExecutorServer
	switch config.ServerType {
	case "tcp":
		factory := communication.NewTCPConnectionFactory(config.ServerAddr)
		server = communication.NewServer(factory)
	case "vsock":
		factory := communication.NewVSockConnectionFactory(config.ServerCid, config.ServerPort)
		server = communication.NewServer(factory)
	default:
		log.Fatalf("Unsupported server type: %s", config.ServerType)
	}

	// Create the executor
	exec, err := executor.NewStatelessExecutor(config, runtime, server)
	if err != nil {
		log.Fatalf("Error creating executor: %v", err)
	}

	// Start the executor
	log.Printf("Starting executor service...")
	if err := exec.Start(ctx); err != nil {
		log.Fatalf("Error starting executor: %v", err)
	}
	log.Println("Executor started")

	// Wait for the context to be canceled
	<-ctx.Done()

	// Stop the executor
	log.Printf("Stopping executor service...")
	if err := exec.Close(); err != nil {
		log.Printf("Error stopping executor: %v", err)
	}
	log.Printf("Executor service stopped")
}
