package communication

import (
	"context"
	"net"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
)

// ExecutorClient defines the interface for communication with the WASM Executor.
// Client means it will connect to the listener created by the ExecutorServer.
// The connection is bi-directional, meaning the client can send requests to the server
// and the server can send requests to the client.
type ExecutorClient interface {
	// Connect establishes a connection to the executor
	Connect(ctx context.Context) error
	// Close closes the connection to the executor
	Close() error
	// SendProcessRequest sends a request to the executor and returns the response
	SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure)
	// SendDeployApp deploys a new application to the executor
	SendDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure)
	// SendGenerateDeanonymizationReport generates a deanonymization report
	SendGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, *apperrors.RequestFailure)
	// SetClientRequestHandler sets the handler for incoming requests from server
	SetClientRequestHandler(handler ClientRequestHandler)
}

// ExecutorServer defines the interface for communication with the Manager.
// Server means it will listen for connections from the ExecutorClient.
// The connection is bi-directional, meaning the client can send requests to the server
// and the server can send requests to the client.
type ExecutorServer interface {
	// Start starts the server
	Start(ctx context.Context) error
	// Stop stops the server
	Stop() error
	// SetRequestHandler sets the handler for incoming requests
	SetRequestHandler(handler RequestHandler)
}

// ClientRequestHandler defines the interface for handling requests from server (the executor) to client (the manager)
type ClientRequestHandler interface {
	//TBD: add here any request needed
}

// RequestHandler defines the interface for handling requests in the WASM Executor
type RequestHandler interface {
	// HandleProcessRequest processes a request and returns the response
	HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure)
	// HandleDeployApp deploys a new application
	HandleDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure)
	// HandleGenerateDeanonymizationReport generates a deanonymization report
	HandleGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, *apperrors.RequestFailure)
}

type ConnectionFactory interface {
	CreateClientConnection() (net.Conn, error)
	CreateServerListener() (net.Listener, error)
}
