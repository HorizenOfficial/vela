package communication

import (
	"context"
	"crypto/sha256"
	"os"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/common"
	apperrors "github.com/horizen-pes/pkg/common/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRequestHandler is a mock implementation of the RequestHandler interface for testing
type MockRequestHandler struct {
	ProcessRequestFunc                func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure)
	DeployAppFunc                     func(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure)
	GenerateDeanonymizationReportFunc func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, *apperrors.RequestFailure)
}

func (m *MockRequestHandler) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure) {
	if m.ProcessRequestFunc != nil {
		return m.ProcessRequestFunc(ctx, req, appState, wasmModule)
	}
	newStateRoot := sha256.Sum256([]byte("new-state-root"))
	return &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			RequestID:     req.RequestID,
			PrevStateRoot: appState.StateRoot,
			NewStateRoot:  newStateRoot,
			Events:        []common.Event{{ApplicationID: req.ApplicationID, EncryptedData: []byte("test-event")}},
			Withdrawals:   []common.Withdrawal{{DestinationAddress: "test-address", Amount: 100}},
			Signature:     []byte("test-signature"),
		},
		&common.ApplicationState{
			ApplicationID:  req.ApplicationID,
			StateRoot:      newStateRoot,
			EncryptedState: []byte("test-encrypted-state"),
		},
		nil
}

func (m *MockRequestHandler) HandleDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure) {
	if m.DeployAppFunc != nil {
		return m.DeployAppFunc(ctx, req)
	}
	newStateRoot := sha256.Sum256([]byte("new-state-root"))
	return &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			RequestID:     req.RequestID,
			PrevStateRoot: [32]byte{},
			NewStateRoot:  newStateRoot,
			Signature:     []byte("test-signature"),
		},
		&common.ApplicationState{
			ApplicationID:  req.ApplicationID,
			StateRoot:      newStateRoot,
			EncryptedState: []byte("test-encrypted-state"),
		},
		nil
}

func (m *MockRequestHandler) HandleGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, *apperrors.RequestFailure) {
	if m.GenerateDeanonymizationReportFunc != nil {
		return m.GenerateDeanonymizationReportFunc(ctx, req, appState, wasmModule)
	}
	return &common.DeanonymizationReport{
			ApplicationID:   req.ApplicationID,
			ReportID:        "test-report-id",
			EncryptedReport: []byte("test-encrypted-report"),
		},
		nil
}

// MockClientRequestHandler is a mock implementation for testing the new client
type MockClientRequestHandler struct {
	GetUserKeysFunc func(ctx context.Context, users []string) (map[string][]byte, *apperrors.RequestFailure)
}

func TestTCPClientServer_ClientToServerRequest(t *testing.T) {
	serverHandler := &MockRequestHandler{}
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8083")
	server := NewServer(factory)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx)
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory)

	// Test connecting to the server
	err = client.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	// Test HandleProcessRequest
	req := &common.Request{
		ProtocolVersion: "1.0",
		ApplicationID:   "test-app",
		RequestID:       "test-request-id",
		RequestType:     common.Process,
		Payload:         []byte("test-encrypted-action"),
		Timestamp:       time.Now().Unix(),
		Sender:          "test-sender",
		//Value:           0,
	}
	appState := &common.ApplicationState{
		ApplicationID:  "test-app",
		StateRoot:      sha256.Sum256([]byte("test-state-root")),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")

	updatePayload, _, failure := client.SendProcessRequest(ctx, req, appState, wasmModule)
	require.Nil(t, failure)
	assert.Equal(t, req.ApplicationID, updatePayload.ApplicationID)
	assert.Equal(t, appState.StateRoot, updatePayload.PrevStateRoot)
	assert.Equal(t, sha256.Sum256([]byte("new-state-root")), updatePayload.NewStateRoot)
	assert.Len(t, updatePayload.Events, 1)
	assert.Equal(t, req.ApplicationID, updatePayload.Events[0].ApplicationID)
	assert.Equal(t, []byte("test-event"), updatePayload.Events[0].EncryptedData)

	// Test HandleDeployApp
	updatedState, appState2, failure := client.SendDeployApp(ctx, req)
	require.Nil(t, failure)
	assert.Equal(t, req.ApplicationID, updatedState.ApplicationID)
	assert.Equal(t, sha256.Sum256([]byte("new-state-root")), updatedState.NewStateRoot)
	assert.Equal(t, req.ApplicationID, appState2.ApplicationID)
	assert.Equal(t, sha256.Sum256([]byte("new-state-root")), appState2.StateRoot)

	// Test HandleGenerateDeanonymizationReport
	report, failure := client.SendGenerateDeanonymizationReport(ctx, req, appState, wasmModule)
	require.Nil(t, failure)
	assert.Equal(t, req.ApplicationID, report.ApplicationID)
	assert.Equal(t, "test-report-id", report.ReportID)
	assert.Equal(t, []byte("test-encrypted-report"), report.EncryptedReport)
}

func TestTCPClientServer_MultipleSequentialRequests(t *testing.T) {
	// This test verifies that multiple sequential bi-directional requests work correctly

	// Create a mock request handler for the server
	serverHandler := &MockRequestHandler{}

	// Create a mock client request handler
	clientHandler := &MockClientRequestHandler{}

	// Create a context
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8086")
	server := NewServer(factory)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx)
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory)
	client.SetClientRequestHandler(clientHandler)

	// Test connecting to the server
	err = client.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	// Perform multiple sequential operations
	for i := 0; i < 5; i++ {
		// Client-initiated request
		req := &common.Request{
			ProtocolVersion: "1.0",
			ApplicationID:   "test-app",
			RequestID:       "test-request-id",
			RequestType:     common.Process,
			Payload:         []byte("test-encrypted-action"),
			Timestamp:       time.Now().Unix(),
			Sender:          "test-sender",
			//Value:           0,
		}
		appState := &common.ApplicationState{
			ApplicationID:  "test-app",
			StateRoot:      sha256.Sum256([]byte("test-state-root")),
			EncryptedState: []byte("test-encrypted-state"),
		}
		wasmModule := []byte("test-wasm-module")

		_, _, failure := client.SendProcessRequest(ctx, req, appState, wasmModule)
		require.Nil(t, failure, "Client request %d should succeed", i)

	}
}

func TestTCPClientServer_ConnectionHandling(t *testing.T) {
	// Test connection lifecycle and error handling

	// Create a mock request handler for the server
	serverHandler := &MockRequestHandler{}

	// Create a context
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8087")
	server := NewServer(factory)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx)
	require.NoError(t, err)
	defer server.Stop()

	// Test connecting and disconnecting multiple times
	for i := 0; i < 3; i++ {
		client := NewClient(factory)

		err = client.Connect(ctx)
		require.NoError(t, err, "Connection %d should succeed", i)

		// Give some time for the connection to be established
		time.Sleep(50 * time.Millisecond)

		// Perform a simple request to verify connection works
		req := &common.Request{
			ProtocolVersion: "1.0",
			ApplicationID:   "test-app",
			RequestID:       "test-request-id",
			RequestType:     common.Deploy,
			Payload:         []byte("test-encrypted-action"),
			Timestamp:       time.Now().Unix(),
			Sender:          "test-sender",
			//Value:           0,
		}

		_, appState, failure := client.SendDeployApp(ctx, req)
		require.Nil(t, failure)
		assert.Equal(t, []byte("test-encrypted-state"), appState.EncryptedState)

		err = client.Close()
		require.NoError(t, err, "Close %d should succeed", i)
	}
}

func TestTCPClientServer_ErrorHandling(t *testing.T) {
	// Create a mock request handler that returns errors
	serverHandler := &MockRequestHandler{
		ProcessRequestFunc: func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure) {
			return nil, nil, apperrors.New(apperrors.CodeInternalFallback, "handler error", assert.AnError)
		},
	}

	// Create a mock client request handler that returns errors
	clientHandler := &MockClientRequestHandler{
		GetUserKeysFunc: func(ctx context.Context, users []string) (map[string][]byte, *apperrors.RequestFailure) {
			return nil, apperrors.New(apperrors.CodeInternalFallback, "handler error", assert.AnError)
		},
	}

	// Create a context
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8088")
	server := NewServer(factory)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx)
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory)
	client.SetClientRequestHandler(clientHandler)

	// Test connecting to the server
	err = client.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	// Test client request error handling
	req := &common.Request{
		ProtocolVersion: "1.0",
		ApplicationID:   "test-app",
		RequestID:       "test-request-id",
		RequestType:     common.Process,
		Payload:         []byte("test-encrypted-action"),
		Timestamp:       time.Now().Unix(),
		Sender:          "test-sender",
		//Value:           0,
	}
	appState := &common.ApplicationState{
		ApplicationID:  "test-app",
		StateRoot:      sha256.Sum256([]byte("test-state-root")),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")

	_, _, failure := client.SendProcessRequest(ctx, req, appState, wasmModule)
	t.Logf("failure: %#v", failure)

	require.NotNil(t, failure, "Client request should return failure")
	assert.Contains(t, failure.Error(), "handler error", "Failure should indicate handler error")

}

func TestTCPClientServer_ServerTimeout(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	ctx := context.Background()
	// Create a mock request handler that simulates slow processing
	serverHandler := &MockRequestHandler{
		ProcessRequestFunc: func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure) {
			// Simulate slow processing that exceeds timeout
			// check is performed each 5 seconds, and timeout is 30 seconds, so 35 is the worst case
			time.Sleep(35 * time.Second)
			return &common.UpdatePayload{
				ApplicationID: req.ApplicationID,
				PrevStateRoot: appState.StateRoot,
				NewStateRoot:  sha256.Sum256([]byte("new-state-root")),
				Events:        []common.Event{{ApplicationID: req.ApplicationID, EncryptedData: []byte("test-event")}},
				Withdrawals:   []common.Withdrawal{{DestinationAddress: "test-address", Amount: 100}},
				Signature:     []byte("test-signature"),
			}, appState, nil
		},
	}

	// Create a server
	factory := NewTCPConnectionFactory(":8089")
	server := NewServer(factory)
	server.SetRequestHandler(serverHandler)
	err := server.Start(context.Background())
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory)
	clientHandler := &MockClientRequestHandler{}
	client.SetClientRequestHandler(clientHandler)

	err = client.Connect(context.Background())
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	req := &common.Request{
		ProtocolVersion: "1.0",
		ApplicationID:   "test-app",
		RequestID:       "test-request-id",
		RequestType:     common.Process,
		Payload:         []byte("test-encrypted-action"),
		Timestamp:       time.Now().Unix(),
		Sender:          "test-sender",
	}
	appState := &common.ApplicationState{
		ApplicationID:  "test-app",
		StateRoot:      sha256.Sum256([]byte("test-state-root")),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")

	var failure *apperrors.RequestFailure
	start := time.Now()
	_, _, failure = client.SendProcessRequest(ctx, req, appState, wasmModule)
	elapsed := time.Since(start)

	// Should timeout and return a failure
	require.NotNil(t, failure)
	assert.Contains(t, failure.Error(), "failed to send process request")
	assert.Greater(t, elapsed, 30*time.Second, "Should timeout at least after 30 seconds")
}
