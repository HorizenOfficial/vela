package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/hf/nsm/request"
	"github.com/hf/nsm/response"
	"github.com/stretchr/testify/assert"

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

