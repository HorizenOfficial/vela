package communication

import (
	"context"
	"github.com/horizen-pes/pkg/common"
	"net"
)

// ClientRequestHandler defines the interface for handling requests from server to client
type ClientRequestHandler interface {
	// GetUserKeys retrieves user keys for the specified application and users
	GetUserKeys(ctx context.Context, users []string) (map[string][]byte, error)
}

// ExecutorClient defines the interface for communication with the WASM Executor
type ExecutorClient interface {
	// Connect establishes a connection to the executor
	Connect(ctx context.Context) error
	// Close closes the connection to the executor
	Close() error
	// ProcessRequest sends a request to the executor and returns the response
	ProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	// DeployApp deploys a new application to the executor
	DeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, []byte, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.DeanonymizationReport, error)
	// SetClientRequestHandler sets the handler for incoming requests from server
	SetClientRequestHandler(handler ClientRequestHandler)
}

// ExecutorServer defines the interface for the WASM Executor server
type ExecutorServer interface {
	// Start starts the server
	Start(ctx context.Context) error
	// Stop stops the server
	Stop() error
	// SetRequestHandler sets the handler for incoming requests
	SetRequestHandler(handler RequestHandler)
	// GetUserKeys requests user keys from the connected client
	GetUserKeys(ctx context.Context, users []string) (map[string][]byte, error)
}

// RequestHandler defines the interface for handling requests in the WASM Executor
type RequestHandler interface {
	// ProcessRequest processes a request and returns the response
	ProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	// DeployApp deploys a new application
	DeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, []byte, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.DeanonymizationReport, error)
}

type ConnectionFactory interface {
	CreateClientConnection() (net.Conn, error)
	CreateServerListener() (net.Listener, error)
}
