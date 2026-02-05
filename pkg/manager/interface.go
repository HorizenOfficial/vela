package manager

import (
	"context"

	"github.com/horizen-pes/pkg/communication"
)

// Manager defines the interface for the Secure Processor Manager
type Manager interface {
	// Start starts the manager
	Start(ctx context.Context) error
	// Stop stops the manager
	Stop() error
	communication.ClientRequestHandler
}
