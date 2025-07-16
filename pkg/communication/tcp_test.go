package communication

import (
	"context"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRequestHandler is a mock implementation of the RequestHandler interface for testing
type MockRequestHandler struct {
	ProcessRequestFunc                func(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	DeployAppFunc                     func(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, []byte, error)
	GenerateDeanonymizationReportFunc func(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.DeanonymizationReport, error)
}

func (m *MockRequestHandler) ProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	if m.ProcessRequestFunc != nil {
		return m.ProcessRequestFunc(ctx, req, appState, senderKey, wasmModule)
	}
	return &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			PrevStateRoot: appState.StateRoot,
			NewStateRoot:  []byte("new-state-root"),
			Events:        []common.Event{{ApplicationID: req.ApplicationID, EncryptedData: []byte("test-event")}},
			Withdrawals:   []common.Withdrawal{{DestinationAddress: "test-address", Amount: "100"}},
			Signature:     []byte("test-signature"),
		},
		&common.ApplicationState{
			ApplicationID:  req.ApplicationID,
			StateRoot:      []byte("new-state-root"),
			EncryptedState: []byte("test-encrypted-state"),
		},
		nil
}

func (m *MockRequestHandler) DeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, []byte, error) {
	if m.DeployAppFunc != nil {
		return m.DeployAppFunc(ctx, req)
	}
	return &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			PrevStateRoot: nil,
			NewStateRoot:  []byte("new-state-root"),
			Signature:     []byte("test-signature"),
		},
		&common.ApplicationState{
			ApplicationID:  req.ApplicationID,
			StateRoot:      []byte("new-state-root"),
			EncryptedState: []byte("test-encrypted-state"),
		},
		[]byte("test-wasm-module"),
		nil
}

func (m *MockRequestHandler) GenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.DeanonymizationReport, error) {
	if m.GenerateDeanonymizationReportFunc != nil {
		return m.GenerateDeanonymizationReportFunc(ctx, req, appState, senderKey, wasmModule)
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
	GetUserKeysFunc func(ctx context.Context, users []string) (map[string][]byte, error)
}

func (m *MockClientRequestHandler) GetUserKeys(ctx context.Context, users []string) (map[string][]byte, error) {
	if m.GetUserKeysFunc != nil {
		return m.GetUserKeysFunc(ctx, users)
	}
	userKeys := make(map[string][]byte)
	for _, user := range users {
		userKeys[user] = []byte("public-key-for-" + user)
	}
	return userKeys, nil
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

	// Test ProcessRequest
	req := &common.Request{
		ProtocolVersion: "1.0",
		ApplicationID:   "test-app",
		RequestID:       "test-request-id",
		RequestType:     common.Process,
		Payload:         []byte("test-encrypted-action"),
		Timestamp:       time.Now().Unix(),
		Sender:          "test-sender",
		Signature:       []byte("test-signature"),
		Value:           0,
	}
	appState := &common.ApplicationState{
		ApplicationID:  "test-app",
		StateRoot:      []byte("test-state-root"),
		EncryptedState: []byte("test-encrypted-state"),
	}
	senderKey := []byte("test-sender-key")
	wasmModule := []byte("test-wasm-module")

	updatePayload, _, err := client.ProcessRequest(ctx, req, appState, senderKey, wasmModule)
	require.NoError(t, err)
	assert.Equal(t, req.ApplicationID, updatePayload.ApplicationID)
	assert.Equal(t, appState.StateRoot, updatePayload.PrevStateRoot)
	assert.Equal(t, []byte("new-state-root"), updatePayload.NewStateRoot)
	assert.Len(t, updatePayload.Events, 1)
	assert.Equal(t, req.ApplicationID, updatePayload.Events[0].ApplicationID)
	assert.Equal(t, []byte("test-event"), updatePayload.Events[0].EncryptedData)

	// Test DeployApp
	updatedState, appState2, wasmBytes, err := client.DeployApp(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, req.ApplicationID, updatedState.ApplicationID)
	assert.Equal(t, []byte("new-state-root"), updatedState.NewStateRoot)
	assert.Equal(t, req.ApplicationID, appState2.ApplicationID)
	assert.Equal(t, []byte("new-state-root"), appState2.StateRoot)
	assert.Equal(t, []byte("test-wasm-module"), wasmBytes)

	// Test GenerateDeanonymizationReport
	report, err := client.GenerateDeanonymizationReport(ctx, req, appState, senderKey, wasmModule)
	require.NoError(t, err)
	assert.Equal(t, req.ApplicationID, report.ApplicationID)
	assert.Equal(t, "test-report-id", report.ReportID)
	assert.Equal(t, []byte("test-encrypted-report"), report.EncryptedReport)
}

func TestTCPClientServer_ServerToClientRequest(t *testing.T) {
	// Create a mock request handler for the server
	serverHandler := &MockRequestHandler{}

	// Create a mock client request handler
	clientHandler := &MockClientRequestHandler{}

	// Create a context
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8084")
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

	// Test server requesting user keys from client
	users := []string{"user1", "user2", "user3"}
	userKeys, err := server.GetUserKeys(ctx, users)
	require.NoError(t, err)
	assert.Len(t, userKeys, 3)
	assert.Equal(t, []byte("public-key-for-user1"), userKeys["user1"])
	assert.Equal(t, []byte("public-key-for-user2"), userKeys["user2"])
	assert.Equal(t, []byte("public-key-for-user3"), userKeys["user3"])
}

func TestTCPClientServer_ConcurrentBidirectionalCommunication(t *testing.T) {
	// This test verifies that bi-directional communication works correctly
	// even when client requests and server requests happen concurrently

	// Create a mock request handler for the server
	serverHandler := &MockRequestHandler{}

	// Create a mock client request handler that simulates some processing time
	clientHandler := &MockClientRequestHandler{
		GetUserKeysFunc: func(ctx context.Context, users []string) (map[string][]byte, error) {
			// Simulate some processing time
			time.Sleep(50 * time.Millisecond)
			userKeys := make(map[string][]byte)
			for _, user := range users {
				userKeys[user] = []byte("public-key-for-" + user)
			}
			return userKeys, nil
		},
	}

	// Create a context
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8085")
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

	// Create channels to coordinate concurrent operations
	clientDone := make(chan error, 1)
	serverDone := make(chan error, 1)

	// Start a client request in a goroutine
	go func() {
		req := &common.Request{
			ProtocolVersion: "1.0",
			ApplicationID:   "test-app",
			RequestID:       "test-request-id",
			RequestType:     common.Process,
			Payload:         []byte("test-encrypted-action"),
			Timestamp:       time.Now().Unix(),
			Sender:          "test-sender",
			Signature:       []byte("test-signature"),
			Value:           0,
		}
		appState := &common.ApplicationState{
			ApplicationID:  "test-app",
			StateRoot:      []byte("test-state-root"),
			EncryptedState: []byte("test-encrypted-state"),
		}
		senderKey := []byte("test-sender-key")
		wasmModule := []byte("test-wasm-module")

		_, _, err := client.ProcessRequest(ctx, req, appState, senderKey, wasmModule)
		clientDone <- err
	}()

	// Start a server request in a goroutine (with a small delay to ensure overlap)
	go func() {
		time.Sleep(25 * time.Millisecond) // Start server request while client request is in progress
		users := []string{"user1", "user2"}
		_, err := server.GetUserKeys(ctx, users)
		serverDone <- err
	}()

	// Wait for both operations to complete
	clientErr := <-clientDone
	serverErr := <-serverDone

	// Both operations should succeed without race conditions
	assert.NoError(t, clientErr, "Client request should succeed")
	assert.NoError(t, serverErr, "Server request should succeed")
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
			Signature:       []byte("test-signature"),
			Value:           0,
		}
		appState := &common.ApplicationState{
			ApplicationID:  "test-app",
			StateRoot:      []byte("test-state-root"),
			EncryptedState: []byte("test-encrypted-state"),
		}
		senderKey := []byte("test-sender-key")
		wasmModule := []byte("test-wasm-module")

		_, _, err := client.ProcessRequest(ctx, req, appState, senderKey, wasmModule)
		require.NoError(t, err, "Client request %d should succeed", i)

		// Server-initiated request
		users := []string{"user1", "user2"}
		userKeys, err := server.GetUserKeys(ctx, users)
		require.NoError(t, err, "Server request %d should succeed", i)
		assert.Len(t, userKeys, 2, "Should get keys for 2 users")
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
			Signature:       []byte("test-signature"),
			Value:           0,
		}

		_, _, wasmBytes, err := client.DeployApp(ctx, req)
		require.NoError(t, err, "Deploy request %d should succeed", i)
		assert.Equal(t, []byte("test-wasm-module"), wasmBytes)

		err = client.Close()
		require.NoError(t, err, "Close %d should succeed", i)
	}
}

func TestTCPClientServer_ErrorHandling(t *testing.T) {
	// Create a mock request handler that returns errors
	serverHandler := &MockRequestHandler{
		ProcessRequestFunc: func(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
			return nil, nil, assert.AnError
		},
	}

	// Create a mock client request handler that returns errors
	clientHandler := &MockClientRequestHandler{
		GetUserKeysFunc: func(ctx context.Context, users []string) (map[string][]byte, error) {
			return nil, assert.AnError
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
		Signature:       []byte("test-signature"),
		Value:           0,
	}
	appState := &common.ApplicationState{
		ApplicationID:  "test-app",
		StateRoot:      []byte("test-state-root"),
		EncryptedState: []byte("test-encrypted-state"),
	}
	senderKey := []byte("test-sender-key")
	wasmModule := []byte("test-wasm-module")

	_, _, err = client.ProcessRequest(ctx, req, appState, senderKey, wasmModule)
	assert.Error(t, err, "Client request should return error")
	assert.Contains(t, err.Error(), "server error", "Error should indicate server error")

	// Test server request error handling
	users := []string{"user1", "user2"}
	_, err = server.GetUserKeys(ctx, users)
	assert.Error(t, err, "Server request should return error")
	assert.Contains(t, err.Error(), "client error", "Error should indicate client error")
}
