package main

import (
	"context"
	"github.com/horizen-pes/pkg/executor"
)

func main() {
	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the executor configuration
	_ = executor.DefaultConfig()

	// Create the executor

	// Start the executor

	// Wait for the context to be canceled
	<-ctx.Done()

	// Stop the executor

}
