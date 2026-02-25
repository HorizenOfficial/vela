package communication

import (
	"context"
	"net"

	"github.com/horizen-pes/pkg/common"
)

// ExecutorClient defines the interface for communication with the WASM Executor.
// Client means it will connect to the listener created by the ExecutorServer.
// The connection is bi-directional, meaning the client can send requests to the server
// and the server can send requests to the client.
type ExecutorClient interface {
	// Connect establishes a connection to the executor
	Connect(ctx context.Context, identityLogTag string) error
	// Close closes the connection to the executor
	Close() error
	// SendProcessRequest sends a request to the executor and returns the response
	// The response includes an optional deanonymization report if the request type was Deanonymize
	SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error)
	// SendDeployApp deploys a new application to the executor
	SendDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error)
	// SendKeyAttestationRequest requests a key attestation from the executor
	SendKeyAttestationRequest(ctx context.Context) ([]byte, error)
	// SetClientRequestHandler sets the handler for incoming requests from server
	SetClientRequestHandler(handler ClientRequestHandler)
}

// ExecutorServer defines the interface for communication with the Manager.
// Server means it will listen for connections from the ExecutorClient.
// The connection is bi-directional, meaning the client can send requests to the server
// and the server can send requests to the client.
type ExecutorServer interface {
	// Start starts the server
	Start(ctx context.Context, identityLogTag string) error
	// Stop stops the server
	Stop() error
	// SetRequestHandler sets the handler for incoming requests
	SetRequestHandler(handler RequestHandler)
	// SetConnectionHandler sets the handler for new client connections.
	SetConnectionHandler(handler ConnectionHandler)
}

// ServerConnection defines the interface for a server-side client connection,
// allowing the server to send requests to the client.
type ServerConnection interface {
	GetKeysetRecovery(ctx context.Context) (bool, *common.EnclaveKeySetRecovery, error)
	SetKeysetRecovery(ctx context.Context, recovery *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr string) error
	KeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr string) error
	Close()
}

// ConnectionHandler is a function that handles a new client connection.
type ConnectionHandler func(ctx context.Context, conn ServerConnection)

// ClientRequestHandler defines the interface for handling requests from server (the executor) to client (the manager)
type ClientRequestHandler interface {
	HandleGetKeysetRecoveryRequest(ctx context.Context) (*common.EnclaveKeySetRecovery, error)
	HandleSetKeysetRecoveryRequest(ctx context.Context, recv *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr string) error
	HandleKeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr string) error
}

// RequestHandler defines the interface for handling requests in the WASM Executor
type RequestHandler interface {
	// HandleProcessRequest processes a request and returns the response
	// The response includes an optional deanonymization report if the request type was Deanonymize
	HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error)
	// HandleDeployApp deploys a new application
	HandleDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error)
	// HandleKeyAttestationRequest creates a key attestation document
	HandleKeyAttestationRequest(ctx context.Context) ([]byte, error)
}

type ConnectionFactory interface {
	CreateClientConnection() (net.Conn, error)
	CreateServerListener() (net.Listener, error)
}
