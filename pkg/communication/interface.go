package communication

import (
	"context"
	"io"

	"github.com/horizen-pes/pkg/common"
)

// ExecutorClient defines the interface for communication with the WASM Executor
type ExecutorClient interface {
	// Connect establishes a connection to the executor
	Connect(ctx context.Context) error
	// Close closes the connection to the executor
	Close() error
	// ProcessRequest sends a request to the executor and returns the response
	ProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	// DeployApp deploys a new application to the executor
	DeployApp(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error)
}

// ExecutorServer defines the interface for the WASM Executor server
type ExecutorServer interface {
	// Start starts the server
	Start(ctx context.Context) error
	// Stop stops the server
	Stop() error
	// SetRequestHandler sets the handler for incoming requests
	SetRequestHandler(handler RequestHandler)
}

// RequestHandler defines the interface for handling requests in the WASM Executor
type RequestHandler interface {
	// ProcessRequest processes a request and returns the response
	ProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	// DeployApp deploys a new application
	DeployApp(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error)
}

// VSockTransport defines the interface for v-socket communication
type VSockTransport interface {
	// Connect establishes a v-socket connection
	Connect(ctx context.Context, cid uint32, port uint32) error
	// Close closes the v-socket connection
	Close() error
	// Read reads data from the v-socket
	Read(p []byte) (n int, err error)
	// Write writes data to the v-socket
	Write(p []byte) (n int, err error)
}

// HTTPTransport defines the interface for HTTP communication
type HTTPTransport interface {
	// Connect establishes an HTTP connection
	Connect(ctx context.Context, url string) error
	// Close closes the HTTP connection
	Close() error
	// Send sends a request and returns the response
	Send(ctx context.Context, method string, path string, body io.Reader) ([]byte, error)
}
