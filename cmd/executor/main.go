package main

import (
	"context"
	"log"

	"github.com/horizen-pes/pkg/communication"

	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/wasm"
)

func main() {
	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the executor configuration
	config := executor.ReadConfig()

	// Create the WASM runtime
	runtime := wasm.NewWasmtimeRuntime()

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

	// TODO temp
	qqq, _ := executor.CreateNewKeySet()
	// Create the executor
	exec, err := executor.NewStatelessExecutor(config, runtime, server, qqq)
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
