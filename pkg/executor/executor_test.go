package executor

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hf/nsm/request"
	"github.com/hf/nsm/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/admin"
	"github.com/horizen-pes/pkg/logger"
)


// mockNsmSession implements the NsmSession interface for testing.
type mockNsmSession struct {
	sendResponse response.Response
	sendError    error
	closeError   error
	sendCalled   bool
	closeCalled  bool
}

func (m *mockNsmSession) Send(req request.Request) (response.Response, error) {
	m.sendCalled = true
	if m.sendError != nil {
		return response.Response{}, m.sendError
	}
	return m.sendResponse, nil
}

func (m *mockNsmSession) Close() error {
	m.closeCalled = true
	return m.closeError
}

func TestCreateKeyAttestationInternal(t *testing.T) {

	log :=	 logger.NewLogger(&logger.Config{
				Kind:         "zerolog",
				Console:      true,
				ConsoleLevel: "trace",
			})


	// A basic executor for testing
	getTestExecutor := func() *StatelessExecutor {
		executor := &StatelessExecutor{log: log}
		return executor
	}

	// Helper to create a valid key set for testing
	newTestKeySet := func() *EnclaveKeySet {
		keySet, err := CreateNewKeySet()
		assert.NoError(t, err)
		return keySet
	}

	t.Run("success", func(t *testing.T) {
		executor := getTestExecutor()
		executor.keySet = newTestKeySet()

		mockDoc := []byte("attestation-document")
		mockSession := &mockNsmSession{
			sendResponse: response.Response{
				Attestation: &response.Attestation{
					Document: mockDoc,
				},
			},
		}
		
		nsmSessionOpener := func() (NsmSession, error) {
			return mockSession, nil
		}

		doc, err := executor.createKeyAttestationInternal(context.Background(), nsmSessionOpener)

		assert.NoError(t, err)
		assert.Equal(t, mockDoc, doc)
		assert.True(t, mockSession.sendCalled)
		assert.True(t, mockSession.closeCalled)
	})

	t.Run("nil keyset", func(t *testing.T) {
		executor := getTestExecutor()
		executor.keySet = nil // Ensure it's nil

		nsmSessionOpener := func() (NsmSession, error) {
			return &mockNsmSession{}, nil
		}

		doc, err := executor.createKeyAttestationInternal(context.Background(), nsmSessionOpener)

		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Contains(t, err.Error(), "keyset is empty")
	})

	t.Run("session open error", func(t *testing.T) {
		executor := getTestExecutor()
		executor.keySet = newTestKeySet()

		expectedErr := errors.New("failed to open session")
		nsmSessionOpener := func() (NsmSession, error) {
			return nil, expectedErr
		}

		doc, err := executor.createKeyAttestationInternal(context.Background(), nsmSessionOpener)

		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Contains(t, err.Error(), "failed to generate attestation")
	})

	t.Run("send error", func(t *testing.T) {
		executor := getTestExecutor()
		executor.keySet = newTestKeySet()

		expectedErr := errors.New("send failed")
		mockSession := &mockNsmSession{
			sendError: expectedErr,
		}

		nsmSessionOpener := func() (NsmSession, error) {
			return mockSession, nil
		}

		doc, err := executor.createKeyAttestationInternal(context.Background(), nsmSessionOpener)

		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Contains(t, err.Error(), "failed to generate attestation")
		assert.True(t, mockSession.closeCalled)
	})

	t.Run("error in response", func(t *testing.T) {

		executor := getTestExecutor()
		executor.keySet = newTestKeySet()

		mockSession := &mockNsmSession{
			sendResponse: response.Response{
				Error: response.ECInternalError,
			},
		}			
		nsmSessionOpener := func() (NsmSession, error) {
			return mockSession, nil
		}

		doc, err := executor.createKeyAttestationInternal(context.Background(), nsmSessionOpener)

		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Contains(t, err.Error(), "failed to generate attestation")
		assert.True(t, mockSession.closeCalled)

	})

	t.Run("nil attestation in response", func(t *testing.T) {
		executor := getTestExecutor()
		executor.keySet = newTestKeySet()

		mockSession := &mockNsmSession{
			sendResponse: response.Response{
				Attestation: nil, // Simulate nil attestation
			},
		}

		nsmSessionOpener := func() (NsmSession, error) {
			return mockSession, nil
		}

		doc, err := executor.createKeyAttestationInternal(context.Background(), nsmSessionOpener)

		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Contains(t, err.Error(), "failed to generate attestation")
		assert.True(t, mockSession.closeCalled)
	})
}

// TestExecuteCommand_SetLogLevel_NonZeroNetworkLogger verifies that SetLogLevel and GetLogLevel
// return errors when the executor uses a non-ZeroNetworkLogger (e.g., zerolog).
func TestExecuteCommand_SetLogLevel_NonZeroNetworkLogger(t *testing.T) {
	log := logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      true,
		ConsoleLevel: "trace",
	})
	executor := &StatelessExecutor{log: log}
	ctx := context.Background()

	// SetLogLevel should fail
	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "error"})
	require.NoError(t, err)
	setMsg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err := executor.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Nil(t, result)

	// GetLogLevel should fail
	getMsg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage}
	result, err = executor.ExecuteCommand(ctx, getMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Empty(t, result)
}

// TestExecuteCommand_TargetValidation verifies that SetLogLevel and GetLogLevel
// reject invalid targets on the executor.
func TestExecuteCommand_TargetValidation(t *testing.T) {
	log := logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      true,
		ConsoleLevel: "trace",
	})
	executor := &StatelessExecutor{log: log}
	ctx := context.Background()

	// SetLogLevel with target "manager" should be rejected
	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "debug", Target: "manager"})
	require.NoError(t, err)
	setMsg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err := executor.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "this is the executor admin server")
	require.Contains(t, err.Error(), "target 'manager' is not supported")
	require.Nil(t, result)

	// SetLogLevel with unknown target should be rejected
	setData, err = json.Marshal(admin.SetLogLevelRequest{Level: "debug", Target: "unknown"})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err = executor.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown target 'unknown'")
	require.Nil(t, result)

	// SetLogLevel with target "executor" should be accepted (fails on logger type, not target)
	setData, err = json.Marshal(admin.SetLogLevelRequest{Level: "debug", Target: "executor"})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err = executor.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Nil(t, result)

	// GetLogLevel with target "manager" should be rejected
	getData, err := json.Marshal(admin.GetLogLevelRequest{Target: "manager"})
	require.NoError(t, err)
	getMsg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Data: getData}
	result, err = executor.ExecuteCommand(ctx, getMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "this is the executor admin server")
	require.Nil(t, result)

	// GetLogLevel with empty target should be accepted (fails on logger type, not target)
	getMsg = admin.AdminMessage{Type: admin.GetLogLevelRequestMessage}
	result, err = executor.ExecuteCommand(ctx, getMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Empty(t, result)

	// SetLogLevel with target "all" should be accepted (fails on logger type, not target)
	setData, err = json.Marshal(admin.SetLogLevelRequest{Level: "debug", Target: "all"})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err = executor.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Nil(t, result)

	// GetLogLevel with target "all" should be accepted (fails on logger type, not target)
	getData, err = json.Marshal(admin.GetLogLevelRequest{Target: "all"})
	require.NoError(t, err)
	getMsg = admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Data: getData}
	result, err = executor.ExecuteCommand(ctx, getMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Empty(t, result)
}