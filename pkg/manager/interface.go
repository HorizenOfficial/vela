package manager

import (
	"context"

	"github.com/HorizenOfficial/vela/pkg/communication"
)

// Manager defines the interface for the Secure Processor Manager
type Manager interface {
	// Start starts the manager
	Start(ctx context.Context) error
	// Stop stops the manager
	Stop() error
	communication.ClientRequestHandler
}
