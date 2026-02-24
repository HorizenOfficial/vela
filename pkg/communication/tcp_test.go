package communication

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	apperrors "github.com/horizen-pes/pkg/common/apperrors"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ApplicationId      = common.NewApplicationId(1)
	senderAddress      = ethCommon.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	destinationAddress = ethCommon.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12")
)

// MockRequestHandler is a mock implementation of the RequestHandler interface for testing
type MockRequestHandler struct {
	ProcessRequestFunc func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error)
	DeployAppFunc      func(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error)
	HelloFunc          func(ctx context.Context, message string) (string, error)
	AdminCommandFunc   func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error)
}

func (m *MockRequestHandler) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
	if m.ProcessRequestFunc != nil {
		return m.ProcessRequestFunc(ctx, req, appState, wasmModule)
	}
	newStateRoot := sha256.Sum256([]byte("new-state-root"))
	return &common.UpdatePayload{
			ApplicationID:  req.ApplicationID,
			RequestID:      req.RequestID,
			PrevStateRoot:  appState.StateRoot,
			NewStateRoot:   newStateRoot,
			Events:         []common.Event{{ApplicationID: req.ApplicationID, EncryptedData: []byte("test-event")}},
			Withdrawals:    []common.Withdrawal{{DestinationAddress: destinationAddress, Amount: common.NewBig(100)}},
			Signature:      []byte("test-signature"),
			RefundAmount:   req.MaxFeeValue,
			ApplicationFee: common.NewBig(100),
		},
		&common.ApplicationState{
			ApplicationID:  req.ApplicationID,
			StateRoot:      newStateRoot,
			EncryptedState: []byte("test-encrypted-state"),
		},
		nil,
		nil
}

func (m *MockRequestHandler) HandleAdminCommand(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
	if m.AdminCommandFunc != nil {
		return m.AdminCommandFunc(ctx, cmdType, data)
	}
	return nil, fmt.Errorf("admin command not supported in mock")
}

func (m *MockRequestHandler) HandleDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error) {
	if m.DeployAppFunc != nil {
		return m.DeployAppFunc(ctx, req, appState)
	}
	newStateRoot := sha256.Sum256([]byte("new-state-root"))
	return &common.UpdatePayload{
			ApplicationID:  req.ApplicationID,
			RequestID:      req.RequestID,
			PrevStateRoot:  [32]byte{},
			NewStateRoot:   newStateRoot,
			Signature:      []byte("test-signature"),
			RefundAmount:   req.MaxFeeValue,
			ApplicationFee: common.NewBig(100),
		},
		&common.ApplicationState{
			ApplicationID:  req.ApplicationID,
			StateRoot:      newStateRoot,
			EncryptedState: []byte("test-encrypted-state"),
		},
		nil
}

// MockClientRequestHandler is a mock implementation for testing the new client
type MockClientRequestHandler struct {
	GetUserKeysFunc       func(ctx context.Context, users []string) (map[string][]byte, *apperrors.RequestFailure)
	SetKeysetRecoveryFunc func(ctx context.Context, recv *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr string) error
	GetKeysetRecoveryFunc func(ctx context.Context) (*common.EnclaveKeySetRecovery, error)
}

// HandleKeysetRecoveryResult implements ClientRequestHandler.
func (m *MockClientRequestHandler) HandleKeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr string) error {
	// For testing purposes, we can just return nil or log something.
	return nil
}

// HandleSetKeysetRecoveryRequest implements ClientRequestHandler.
func (m *MockClientRequestHandler) HandleSetKeysetRecoveryRequest(ctx context.Context, recv *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr string) error {
	if m.SetKeysetRecoveryFunc != nil {
		return m.SetKeysetRecoveryFunc(ctx, recv, commPubKey, signingKeyAddr)
	}
	return nil
}

func (m *MockClientRequestHandler) HandleGetKeysetRecoveryRequest(ctx context.Context) (*common.EnclaveKeySetRecovery, error) {
	if m.GetKeysetRecoveryFunc != nil {
		return m.GetKeysetRecoveryFunc(ctx)
	}
	recv := &common.EnclaveKeySetRecovery{
		RecoveryType:       common.RecoveryTypeKMS,
		KeySetCiphertext:   []byte{0x01, 0x02, 0x03},
		RecoveryCiphertext: []byte{0x04, 0x05, 0x06},
	}
	return recv, nil
}

func TestTCPClientServer_ClientToServerRequest(t *testing.T) {
	serverHandler := &MockRequestHandler{}
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8083")
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx, "Server")
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory, commParams, testLogger)

	// Test connecting to the server
	err = client.Connect(ctx, "Client")
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	// Test HandleProcessRequest
	req := &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   ApplicationId,
		RequestID:       testutil.GenerateRandomRequestID(),
		RequestType:     common.Process,
		Payload:         []byte("test-encrypted-action"),
		Timestamp:       common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		Sender:          senderAddress,
		DepositAmount:   common.NewBig(0),
		MaxFeeValue:     common.NewBig(100),
	}
	appState := &common.ApplicationState{
		ApplicationID:  ApplicationId,
		StateRoot:      sha256.Sum256([]byte("test-state-root")),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")

	updatePayload, _, _, failure := client.SendProcessRequest(ctx, req, appState, wasmModule)
	require.Nil(t, failure)
	assert.Equal(t, req.ApplicationID, updatePayload.ApplicationID)
	assert.Equal(t, appState.StateRoot, updatePayload.PrevStateRoot)
	assert.Equal(t, sha256.Sum256([]byte("new-state-root")), updatePayload.NewStateRoot)
	assert.Len(t, updatePayload.Events, 1)
	assert.Equal(t, req.ApplicationID, updatePayload.Events[0].ApplicationID)
	assert.Equal(t, []byte("test-event"), updatePayload.Events[0].EncryptedData)

	// Test HandleDeployApp
	updatedState, appState2, failure := client.SendDeployApp(ctx, req, nil)
	require.Nil(t, failure)
	assert.Equal(t, req.ApplicationID, updatedState.ApplicationID)
	assert.Equal(t, sha256.Sum256([]byte("new-state-root")), updatedState.NewStateRoot)
	assert.Equal(t, req.ApplicationID, appState2.ApplicationID)
	assert.Equal(t, sha256.Sum256([]byte("new-state-root")), appState2.StateRoot)
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
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx, "Server")
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory, commParams, testLogger)
	client.SetClientRequestHandler(clientHandler)

	// Test connecting to the server
	err = client.Connect(ctx, "Client")
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	// Perform multiple sequential operations
	for i := 0; i < 5; i++ {
		// Client-initiated request
		req := &common.Request{
			ProtocolVersion: 0,
			ApplicationID:   ApplicationId,
			RequestID:       testutil.GenerateRandomRequestID(),
			RequestType:     common.Process,
			Payload:         []byte("test-encrypted-action"),
			Timestamp:       common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
			Sender:          senderAddress,
			DepositAmount:   common.NewBig(0),
			MaxFeeValue:     common.NewBig(100),
		}
		appState := &common.ApplicationState{
			ApplicationID:  ApplicationId,
			StateRoot:      sha256.Sum256([]byte("test-state-root")),
			EncryptedState: []byte("test-encrypted-state"),
		}
		wasmModule := []byte("test-wasm-module")

		_, _, _, failure := client.SendProcessRequest(ctx, req, appState, wasmModule)
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
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx, "Server")
	require.NoError(t, err)
	defer server.Stop()

	// Test connecting and disconnecting multiple times
	for i := 0; i < 3; i++ {
		client := NewClient(factory, commParams, testLogger)

		err = client.Connect(ctx, "Client")
		require.NoError(t, err, "Connection %d should succeed", i)

		// Give some time for the connection to be established
		time.Sleep(50 * time.Millisecond)

		// Perform a simple request to verify connection works
		req := &common.Request{
			ProtocolVersion: 0,
			ApplicationID:   ApplicationId,
			RequestID:       testutil.GenerateRandomRequestID(),
			RequestType:     common.Deploy,
			Payload:         []byte("test-encrypted-action"),
			Timestamp:       common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
			Sender:          senderAddress,
			DepositAmount:   common.NewBig(0),
			MaxFeeValue:     common.NewBig(100),
		}

		_, appState, failure := client.SendDeployApp(ctx, req, nil)
		require.Nil(t, failure)
		assert.Equal(t, []byte("test-encrypted-state"), appState.EncryptedState)

		err = client.Close()
		require.NoError(t, err, "Close %d should succeed", i)
	}
}

func TestTCPClientServer_ErrorHandling(t *testing.T) {
	// Create a mock request handler that returns errors
	serverHandler := &MockRequestHandler{
		ProcessRequestFunc: func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
			return nil, nil, nil, apperrors.New(apperrors.CodeInternalFallback, "handler error")
		},
	}

	// Create a mock client request handler that returns errors
	clientHandler := &MockClientRequestHandler{
		GetUserKeysFunc: func(ctx context.Context, users []string) (map[string][]byte, *apperrors.RequestFailure) {
			return nil, apperrors.New(apperrors.CodeInternalFallback, "handler error")
		},
	}

	// Create a context
	ctx := context.Background()

	// Create a server
	factory := NewTCPConnectionFactory(":8088")
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx, "Server")
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory, commParams, testLogger)
	client.SetClientRequestHandler(clientHandler)

	// Test connecting to the server
	err = client.Connect(ctx, "Client")
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	// Test client request error handling
	req := &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   ApplicationId,
		RequestID:       testutil.GenerateRandomRequestID(),
		RequestType:     common.Process,
		Payload:         []byte("test-encrypted-action"),
		Timestamp:       common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		Sender:          senderAddress,
		DepositAmount:   common.NewBig(0),
		MaxFeeValue:     common.NewBig(100),
	}
	appState := &common.ApplicationState{
		ApplicationID:  ApplicationId,
		StateRoot:      sha256.Sum256([]byte("test-state-root")),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")

	_, _, _, failure := client.SendProcessRequest(ctx, req, appState, wasmModule)
	t.Logf("failure: %#v", failure)

	require.NotNil(t, failure, "Client request should return failure")
	assert.Contains(t, failure.Error(), "handler error", "Failure should indicate handler error")

}

func TestTCPClientServer_ServerToClientRequest(t *testing.T) {
	// This test verifies that the server can send a request to the client
	// and receive a response.

	// Create a mock request handler for the server
	serverHandler := &MockRequestHandler{}

	// Use a channel to signal when the hello message has been handled
	helloHandled := make(chan bool, 1)

	// Create a mock client request handler
	clientHandler := &MockClientRequestHandler{
		GetKeysetRecoveryFunc: func(ctx context.Context) (*common.EnclaveKeySetRecovery, error) {
			t.Log("GetKeysetRecoveryFunc called")
			helloHandled <- true
			t.Log("helloHandled sent true")
			return &common.EnclaveKeySetRecovery{},
				nil
		},
	}

	// Create a context
	ctx := context.Background()

	// Add a waitgroup to ensure the server-side goroutine completes
	var wg sync.WaitGroup

	// Create a server
	factory := NewTCPConnectionFactory(":8090")
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	server.SetConnectionHandler(func(ctx context.Context, conn ServerConnection) {
		wg.Add(1)
		// When a client connects, send a request to it.
		go func() {
			defer wg.Done()
			_, _, err := conn.GetKeysetRecovery(ctx)
			require.NoError(t, err)
		}()
	})
	err := server.Start(ctx, "Server")
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory, commParams, testLogger)
	client.SetClientRequestHandler(clientHandler)

	// Test connecting to the server
	err = client.Connect(ctx, "Client")
	require.NoError(t, err)
	defer client.Close()

	// Wait for the hello message to be handled or timeout
	select {
	case <-helloHandled:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for hello message to be handled")
	}

	// Wait for the server-side goroutine to finish
	wg.Wait()
}

func TestTCPClientServer_ServerTimeout(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	ctx := context.Background()
	// Create a mock request handler that simulates slow processing
	serverHandler := &MockRequestHandler{
		ProcessRequestFunc: func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
			// Simulate slow processing that exceeds timeout
			// check is performed each 5 seconds, and timeout is 30 seconds, so 35 is the worst case
			time.Sleep(35 * time.Second)
			return &common.UpdatePayload{
				ApplicationID:  req.ApplicationID,
				PrevStateRoot:  appState.StateRoot,
				NewStateRoot:   sha256.Sum256([]byte("new-state-root")),
				Events:         []common.Event{{ApplicationID: req.ApplicationID, EncryptedData: []byte("test-event")}},
				Withdrawals:    []common.Withdrawal{{DestinationAddress: destinationAddress, Amount: common.NewBig(100)}},
				Signature:      []byte("test-signature"),
				RefundAmount:   req.MaxFeeValue,
				ApplicationFee: common.NewBig(100),
			}, appState, nil, nil
		},
	}

	// Create a server
	factory := NewTCPConnectionFactory(":8089")
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	err := server.Start(context.Background(), "Server")
	require.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewClient(factory, commParams, testLogger)
	clientHandler := &MockClientRequestHandler{}
	client.SetClientRequestHandler(clientHandler)

	err = client.Connect(context.Background(), "Client")
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	req := &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   1,
		RequestID:       testutil.GenerateRandomRequestID(),
		RequestType:     common.Process,
		Payload:         []byte("test-encrypted-action"),
		Timestamp:       common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		Sender:          senderAddress,
		DepositAmount:   common.NewBig(0),
		MaxFeeValue:     common.NewBig(100),
	}
	appState := &common.ApplicationState{
		ApplicationID:  1,
		StateRoot:      sha256.Sum256([]byte("test-state-root")),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")

	var failure error
	start := time.Now()
	_, _, _, failure = client.SendProcessRequest(ctx, req, appState, wasmModule)
	elapsed := time.Since(start)

	// Should timeout and return a failure
	require.NotNil(t, failure)
	assert.Contains(t, failure.Error(), "failed to send process request")
	assert.Greater(t, elapsed, 30*time.Second, "Should timeout at least after 30 seconds")
}
