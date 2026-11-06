package communication

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	apperrors "github.com/HorizenOfficial/vela/pkg/common/apperrors"
	"github.com/HorizenOfficial/vela/pkg/common/testutil"
	ethCommon "github.com/ethereum/go-ethereum/common"
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
	ProcessRequestFunc      func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error)
	BatchProcessRequestFunc func(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasmModule []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error)
	DeployAppFunc           func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	AdminCommandFunc        func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error)
}

func (m *MockRequestHandler) HandleBatchProcessRequest(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasmModule []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
	if m.BatchProcessRequestFunc != nil {
		return m.BatchProcessRequestFunc(ctx, requests, appState, wasmModule)
	}
	return nil, nil, nil, nil, fmt.Errorf("batch process request not supported in mock")
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

func (m *MockRequestHandler) HandleDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	if m.DeployAppFunc != nil {
		return m.DeployAppFunc(ctx, req, appState, wasmModule)
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
		TokenAddress:    velacommon.ETH_TOKEN,
		AssetAmount:     common.NewBig(0),
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
	deployWasm := []byte("deploy-wasm-module")
	updatedState, appState2, failure := client.SendDeployApp(ctx, req, nil, deployWasm)
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
			TokenAddress:    velacommon.ETH_TOKEN,
			AssetAmount:     common.NewBig(0),
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
			TokenAddress:    velacommon.ETH_TOKEN,
			AssetAmount:     common.NewBig(0),
			MaxFeeValue:     common.NewBig(100),
		}

		_, appState, failure := client.SendDeployApp(ctx, req, nil, []byte("deploy-wasm-module"))
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
		TokenAddress:    velacommon.ETH_TOKEN,
		AssetAmount:     common.NewBig(0),
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

func TestTCPClientServer_AdminCommandRequest(t *testing.T) {
	expectedResponse := json.RawMessage(`"test-attestation-document-bytes"`)

	serverHandler := &MockRequestHandler{
		AdminCommandFunc: func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "key_attestation", cmdType)
			return expectedResponse, nil
		},
	}

	ctx := context.Background()

	factory := NewTCPConnectionFactory(":8091")
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx, "Server")
	require.NoError(t, err)
	defer server.Stop()

	client := NewClient(factory, commParams, testLogger)
	err = client.Connect(ctx, "Client")
	require.NoError(t, err)
	defer client.Close()

	time.Sleep(100 * time.Millisecond)

	// Test success case
	result, err := client.ForwardAdminCommand(ctx, "key_attestation", nil)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, result)

	// Test error case
	serverHandler.AdminCommandFunc = func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
		return nil, assert.AnError
	}

	result, err = client.ForwardAdminCommand(ctx, "key_attestation", nil)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "COMMAND_ERROR")
}

// TestTCPClientServer_BatchProcessRequest exercises the batch message end to end
// over a real connection: client.SendBatchProcessRequest -> server routing ->
// RequestHandler -> BatchProcessResponseData -> client. The direct handler tests
// in pkg/executor bypass the protocol layer, so this is the only coverage for
// BatchProcessRequestData.Validate and for the response fields surviving JSON.
func TestTCPClientServer_BatchProcessRequest(t *testing.T) {
	newBatchRequest := func() *common.Request {
		return &common.Request{
			ProtocolVersion: 0,
			ApplicationID:   ApplicationId,
			RequestID:       testutil.GenerateRandomRequestID(),
			RequestType:     common.Process,
			Payload:         []byte("test-encrypted-action"),
			Timestamp:       common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
			Sender:          senderAddress,
			TokenAddress:    velacommon.ETH_TOKEN,
			AssetAmount:     common.NewBig(0),
			MaxFeeValue:     common.NewBig(100),
		}
	}

	appState := &common.ApplicationState{
		ApplicationID:  ApplicationId,
		StateRoot:      sha256.Sum256([]byte("test-state-root")),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")
	batchSignature := bytes.Repeat([]byte{0xAB}, 65)

	// chainPayloads builds `count` unsigned payloads whose state roots chain from
	// appState, mirroring what the real executor returns for a batch.
	chainPayloads := func(requests []*common.Request, count int) ([]*common.UpdatePayload, [32]byte) {
		prev := appState.StateRoot
		payloads := make([]*common.UpdatePayload, 0, count)
		for i := 0; i < count; i++ {
			next := sha256.Sum256([]byte(fmt.Sprintf("batch-state-root-%d", i)))
			payloads = append(payloads, &common.UpdatePayload{
				ApplicationID:  requests[i].ApplicationID,
				RequestID:      requests[i].RequestID,
				PrevStateRoot:  prev,
				NewStateRoot:   next,
				Events:         []common.Event{{ApplicationID: requests[i].ApplicationID, EncryptedData: []byte(fmt.Sprintf("batch-event-%d", i))}},
				RefundAmount:   requests[i].MaxFeeValue,
				ApplicationFee: common.NewBig(100),
			})
			prev = next
		}
		return payloads, prev
	}

	var (
		gotRequests []*common.Request
		gotAppState *common.ApplicationState
		gotWasm     []byte
	)

	// Default handler: every request succeeds.
	serverHandler := &MockRequestHandler{
		BatchProcessRequestFunc: func(ctx context.Context, requests []*common.Request, state *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			gotRequests, gotAppState, gotWasm = requests, state, wasm

			payloads, finalRoot := chainPayloads(requests, len(requests))
			finalState := &common.ApplicationState{
				ApplicationID:  state.ApplicationID,
				StateRoot:      finalRoot,
				EncryptedState: []byte("final-encrypted-state"),
			}
			reports := []*common.DeanonymizationReport{{
				ApplicationID:   requests[0].ApplicationID,
				ReportID:        requests[0].RequestID,
				EncryptedReport: []byte("encrypted-report"),
				Authority:       senderAddress,
			}}
			return payloads, batchSignature, finalState, reports, nil
		},
	}

	ctx := context.Background()

	factory := NewTCPConnectionFactory(reserveTCPAddress(t))
	server := NewServer(factory, commParams, testLogger)
	server.SetRequestHandler(serverHandler)
	err := server.Start(ctx, "Server")
	require.NoError(t, err)
	defer server.Stop()

	client := NewClient(factory, commParams, testLogger)
	err = client.Connect(ctx, "Client")
	require.NoError(t, err)
	defer client.Close()

	// Give some time for the connection to be established
	time.Sleep(100 * time.Millisecond)

	t.Run("full batch round-trip", func(t *testing.T) {
		requests := []*common.Request{newBatchRequest(), newBatchRequest(), newBatchRequest()}

		payloads, sig, finalState, reports, err := client.SendBatchProcessRequest(ctx, requests, appState, wasmModule)
		require.NoError(t, err)

		// The request side arrived intact
		require.Len(t, gotRequests, len(requests))
		for i, r := range requests {
			assert.Equal(t, r.RequestID, gotRequests[i].RequestID)
			assert.Equal(t, r.ApplicationID, gotRequests[i].ApplicationID)
		}
		require.NotNil(t, gotAppState)
		assert.Equal(t, appState.StateRoot, gotAppState.StateRoot)
		assert.Equal(t, wasmModule, gotWasm)

		// The response side survived JSON
		require.Len(t, payloads, len(requests))
		assert.Equal(t, batchSignature, sig)
		assert.Equal(t, appState.StateRoot, payloads[0].PrevStateRoot)
		for i := 1; i < len(payloads); i++ {
			assert.Equal(t, payloads[i-1].NewStateRoot, payloads[i].PrevStateRoot)
		}
		for i, p := range payloads {
			assert.Equal(t, requests[i].RequestID, p.RequestID)
			assert.Empty(t, p.Signature, "batch payloads must not be individually signed")
			require.Len(t, p.Events, 1)
			assert.Equal(t, []byte(fmt.Sprintf("batch-event-%d", i)), p.Events[0].EncryptedData)
		}
		require.NotNil(t, finalState)
		assert.Equal(t, payloads[len(payloads)-1].NewStateRoot, finalState.StateRoot)
		assert.Equal(t, []byte("final-encrypted-state"), finalState.EncryptedState)
		require.Len(t, reports, 1)
		assert.Equal(t, requests[0].RequestID, reports[0].ReportID)
	})

	t.Run("partial batch returns fewer payloads than requests", func(t *testing.T) {
		requests := []*common.Request{newBatchRequest(), newBatchRequest(), newBatchRequest()}

		// A hard failure stopped the batch after the first request: fewer payloads
		// than requests, and no error.
		serverHandler.BatchProcessRequestFunc = func(ctx context.Context, requests []*common.Request, state *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			payloads, finalRoot := chainPayloads(requests, 1)
			finalState := &common.ApplicationState{
				ApplicationID:  state.ApplicationID,
				StateRoot:      finalRoot,
				EncryptedState: []byte("final-encrypted-state"),
			}
			return payloads, batchSignature, finalState, nil, nil
		}

		payloads, sig, finalState, reports, err := client.SendBatchProcessRequest(ctx, requests, appState, wasmModule)
		require.NoError(t, err)
		require.Len(t, payloads, 1)
		assert.Less(t, len(payloads), len(requests), "manager must be able to see the batch was cut short")
		assert.Equal(t, batchSignature, sig)
		assert.Equal(t, payloads[0].NewStateRoot, finalState.StateRoot)
		assert.Empty(t, reports)
	})

	t.Run("nil ApplicationState is rejected without killing the server", func(t *testing.T) {
		// The manager passes a nil state when it has none for the application.
		// This must come back as a protocol error; before the Validate nil check
		// it dereferenced the nil state and panicked inside the executor process.
		serverHandler.BatchProcessRequestFunc = func(ctx context.Context, requests []*common.Request, state *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			t.Error("handler must not be reached for an invalid batch")
			return nil, nil, nil, nil, nil
		}

		payloads, sig, finalState, reports, err := client.SendBatchProcessRequest(ctx, []*common.Request{newBatchRequest()}, nil, wasmModule)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ApplicationState is required")
		assert.Nil(t, payloads)
		assert.Nil(t, sig)
		assert.Nil(t, finalState)
		assert.Nil(t, reports)

		// The connection must still be usable: a panic in the read loop would
		// have taken the server down with it.
		serverHandler.BatchProcessRequestFunc = func(ctx context.Context, requests []*common.Request, state *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			payloads, _ := chainPayloads(requests, len(requests))
			return payloads, batchSignature, state, nil, nil
		}

		payloads, _, _, _, err = client.SendBatchProcessRequest(ctx, []*common.Request{newBatchRequest()}, appState, wasmModule)
		require.NoError(t, err, "server must survive an invalid batch")
		require.Len(t, payloads, 1)
	})

	t.Run("mismatched applicationId is rejected", func(t *testing.T) {
		serverHandler.BatchProcessRequestFunc = func(ctx context.Context, requests []*common.Request, state *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			t.Error("handler must not be reached for an invalid batch")
			return nil, nil, nil, nil, nil
		}

		foreign := newBatchRequest()
		foreign.ApplicationID = common.NewApplicationId(999)

		_, _, _, _, err := client.SendBatchProcessRequest(ctx, []*common.Request{newBatchRequest(), foreign}, appState, wasmModule)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ApplicationID mismatch")
	})

	t.Run("handler error is propagated", func(t *testing.T) {
		serverHandler.BatchProcessRequestFunc = func(ctx context.Context, requests []*common.Request, state *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			return nil, nil, nil, nil, fmt.Errorf("wasm fingerprint mismatch for application 1")
		}

		payloads, sig, finalState, reports, err := client.SendBatchProcessRequest(ctx, []*common.Request{newBatchRequest()}, appState, wasmModule)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "wasm fingerprint mismatch")
		assert.Nil(t, payloads)
		assert.Nil(t, sig)
		assert.Nil(t, finalState)
		assert.Nil(t, reports)
	})
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
			_, err := conn.GetKeysetRecovery(ctx)
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
	factory := NewTCPConnectionFactory(reserveTCPAddress(t))
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
		TokenAddress:    velacommon.ETH_TOKEN,
		AssetAmount:     common.NewBig(0),
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

func reserveTCPAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	return listener.Addr().String()
}
