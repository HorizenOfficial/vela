package communication

import (
	"context"
	"github.com/horizen-pes/pkg/common"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockRequestHandler is a mock implementation of the RequestHandler interface for testing
type MockRequestHandler struct {
	ProcessRequestFunc                func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	DeployAppFunc                     func(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error)
	GenerateDeanonymizationReportFunc func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error)
}

func (m MockRequestHandler) ProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	return m.ProcessRequestFunc(ctx, req, appState, wasmModule)
}

func (m MockRequestHandler) DeployApp(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error) {
	return m.DeployAppFunc(ctx, req)
}

func (m MockRequestHandler) GenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error) {
	return m.GenerateDeanonymizationReportFunc(ctx, req, appState, wasmModule)
}

func TestTCPClientServer(t *testing.T) {
	// Create a mock request handler
	handler := &MockRequestHandler{
		ProcessRequestFunc: func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
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
		},
		DeployAppFunc: func(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error) {
			return &common.ApplicationState{
					ApplicationID:  req.ApplicationID,
					StateRoot:      []byte("new-state-root"),
					EncryptedState: []byte("test-encrypted-state"),
				},
				[]byte("test-wasm-module"),
				nil
		},
		GenerateDeanonymizationReportFunc: func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error) {
			return &common.DeanonymizationReport{
				ApplicationID:   req.ApplicationID,
				ReportID:        "test-report-id",
				EncryptedReport: []byte("test-encrypted-report"),
			}, nil
		},
	}

	// Create a context
	ctx := context.Background()

	// Create a server
	server := NewTCPServer(":8081")
	server.SetRequestHandler(handler)
	err := server.Start(ctx)
	assert.NoError(t, err)
	defer server.Stop()

	// Create a client
	client := NewTCPClient("127.0.0.1:8081")

	// Test connecting to the server
	err = client.Connect(ctx)
	assert.NoError(t, err)
	defer client.Close()

	// Test processing a request
	req := &common.Request{
		ProtocolVersion: "1.0",
		ApplicationID:   "test-app",
		RequestType:     common.Process,
		Payload:         []byte("test-payload"),
		Timestamp:       time.Now().Unix(),
		Sender:          "test-sender",
		Signature:       []byte("test-signature"),
	}
	appState := &common.ApplicationState{
		ApplicationID:  "test-app",
		StateRoot:      []byte("test-state-root"),
		EncryptedState: []byte("test-encrypted-state"),
	}
	wasmModule := []byte("test-wasm-module")
	updatePayload, updatedState, err := client.ProcessRequest(ctx, req, appState, wasmModule)
	assert.NoError(t, err)
	assert.Equal(t, req.ApplicationID, updatePayload.ApplicationID)
	assert.Equal(t, appState.StateRoot, updatePayload.PrevStateRoot)
	assert.Equal(t, []byte("new-state-root"), updatePayload.NewStateRoot)
	assert.Len(t, updatePayload.Events, 1)
	assert.Equal(t, req.ApplicationID, updatePayload.Events[0].ApplicationID)
	assert.Equal(t, []byte("test-event"), updatePayload.Events[0].EncryptedData)
	assert.Len(t, updatePayload.Withdrawals, 1)
	assert.Equal(t, "test-address", updatePayload.Withdrawals[0].DestinationAddress)
	assert.Equal(t, "100", updatePayload.Withdrawals[0].Amount)
	assert.Equal(t, []byte("test-signature"), updatePayload.Signature)
	// Check the updated application state
	assert.Equal(t, req.ApplicationID, updatedState.ApplicationID)
	assert.Equal(t, []byte("new-state-root"), updatedState.StateRoot)
	assert.Equal(t, []byte("test-encrypted-state"), updatedState.EncryptedState)

	// Test deploying an application
	appState, wasmBytes, err := client.DeployApp(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.ApplicationID, appState.ApplicationID)
	assert.Equal(t, []byte("new-state-root"), appState.StateRoot)
	assert.Equal(t, []byte("test-encrypted-state"), appState.EncryptedState)
	assert.Equal(t, []byte("test-wasm-module"), wasmBytes)

	// Test generating a deanonymization report
	report, err := client.GenerateDeanonymizationReport(ctx, req, appState, wasmModule)
	assert.NoError(t, err)
	assert.Equal(t, req.ApplicationID, report.ApplicationID)
	assert.Equal(t, "test-report-id", report.ReportID)
	assert.Equal(t, []byte("test-encrypted-report"), report.EncryptedReport)

	client.Close()
	server.Stop()
}
