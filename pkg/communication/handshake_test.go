package communication

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/logger"
	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
)

var testLogger logger.Logger

var commParams = common.CommunicationParams{RequestTimeoutSec: 30}

func TestMain(m *testing.M) {
	// Initialize once
	//	testLogger = logger.NewLogger(&logger.Config{Kind: "printf"})
	testLogger = logger.NewLogger(
		&logger.Config{
			Kind:         "zerolog",
			ConsoleColor: false, // colors can print escape chars on tty
			Console:      true,
			ConsoleLevel: "trace",
			//FileName:     "qqq.log",
			//FileLevel:    "info",
		},
	)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// mockExecutor is a mock implementation of the server-side (executor) handshake logic
type mockExecutor struct {
	t              *testing.T
	conn           ServerConnection
	handshakeDone  chan struct{}
	handshakeError error

	// reportedTimeoutMs records the RequestTimeoutMs the manager sent in its
	// GetKeysetRecovery reply, so tests can check the field is populated.
	reportedTimeoutMs int64
}

func (m *mockExecutor) handleConnection(ctx context.Context, conn ServerConnection) {
	m.conn = conn
	go func() {
		err := m.performHandshake(ctx)
		if err != nil {
			m.handshakeError = err
		}
		close(m.handshakeDone)
	}()
}

func (m *mockExecutor) performHandshake(ctx context.Context) error {
	testLogger.Info("MockExecutor: entering %s", common.FnName())
	handshake, err := m.conn.GetKeysetRecovery(ctx)
	if err != nil {
		return err
	}
	m.reportedTimeoutMs = handshake.RequestTimeoutMs

	if handshake.DataFound {
		recoveryData := handshake.KeySetRecovery
		// Simulate restoring keyset
		if recoveryData == nil {
			return nil
		}
		if bytes.Equal(recoveryData.KeySetCiphertext, []byte("corrupted-keyset")) {
			return fmt.Errorf("simulated restore error")
		}
		testLogger.Info("MockExecutor: simulating restoring keyset")
		return m.conn.KeysetRecoveryResult(ctx, nil, "mock-comm-pub-key", "mock-signing-addr")
	} else {
		// Simulate generating new keyset
		newRecoveryData := &common.EnclaveKeySetRecovery{
			RecoveryType:       common.RecoveryTypeKMS,
			KeySetCiphertext:   []byte("new-keyset"),
			RecoveryCiphertext: []byte("new-recovery"),
		}
		testLogger.Info("MockExecutor: simulating new keyset")
		return m.conn.SetKeysetRecovery(ctx, newRecoveryData, "mock-comm-pub-key", "mock-signing-addr")
	}
}

// mockManager is a mock implementation of the client-side (manager) handshake logic
type mockManager struct {
	t                     *testing.T
	recoveryData          *common.EnclaveKeySetRecovery
	getRecoveryError      error
	setRecoveryError      error
	handshakeSuccess      bool
	handshakeSuccessMutex sync.Mutex
}

func (m *mockManager) HandleGetKeysetRecoveryRequest(ctx context.Context) (*common.EnclaveKeySetRecovery, error) {
	testLogger.Info("MockExecutor: entering %s", common.FnName())
	if m.getRecoveryError != nil {
		return nil, m.getRecoveryError
	}
	if m.recoveryData == nil {
		m.recoveryData = &common.EnclaveKeySetRecovery{
			RecoveryType:       common.RecoveryTypeKMS,
			KeySetCiphertext:   []byte("existing-keyset"),
			RecoveryCiphertext: []byte("existing-recovery"),
		}
	}
	return m.recoveryData, nil
}

func (m *mockManager) HandleSetKeysetRecoveryRequest(ctx context.Context, recv *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr string) error {
	testLogger.Info("MockExecutor: entering %s", common.FnName())
	if m.setRecoveryError != nil {
		return m.setRecoveryError
	}
	m.recoveryData = recv
	return nil
}

func (m *mockManager) HandleKeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr string) error {
	testLogger.Info("MockManager: entering %s", common.FnName())
	if result == nil {
		m.handshakeSuccessMutex.Lock()
		m.handshakeSuccess = true
		m.handshakeSuccessMutex.Unlock()
	}
	return nil
}

func (m *mockManager) wasHandshakeSuccessful() bool {
	testLogger.Info("MockExecutor: entering %s", common.FnName())
	m.handshakeSuccessMutex.Lock()
	defer m.handshakeSuccessMutex.Unlock()
	return m.handshakeSuccess
}

func setupHandshakeTest(t *testing.T) (context.Context, *Client, *Server, *mockExecutor, *mockManager) {
	testLogger.Info("MockExecutor: entering %s", common.FnName())
	ctx := context.Background()
	//factory := NewTCPConnectionFactory(":0") // Use :0 for random port
	factory := NewTCPConnectionFactory("localhost:1234") // Use :0 for random port

	// Setup server (executor)
	server := NewServer(factory, commParams, testLogger)
	executor := &mockExecutor{t: t, handshakeDone: make(chan struct{})}
	server.SetConnectionHandler(executor.handleConnection)
	err := server.Start(ctx, "Executor")
	require.NoError(t, err)

	// Setup client (manager)
	// We need to get the random port from the server
	//addr := server.Addr().String()
	//clientFactory := NewTCPConnectionFactory(addr)
	clientFactory := NewTCPConnectionFactory("localhost:1234")
	client := NewClient(clientFactory, commParams, testLogger)
	manager := &mockManager{t: t}
	client.SetClientRequestHandler(manager)

	return ctx, client, server, executor, manager
}

func TestHandshake_FirstConnection(t *testing.T) {
	ctx, client, server, executor, manager := setupHandshakeTest(t)
	defer client.Close()
	defer server.Stop()

	// Manager returns error on get because this is the very first connection
	manager.getRecoveryError = storageErrors.NewError(storageErrors.NotFound, "Test Error")

	err := client.Connect(ctx, "Manager")
	require.NoError(t, err)

	// Wait for handshake to complete
	<-executor.handshakeDone

	require.NoError(t, executor.handshakeError)
	require.NotNil(t, manager.recoveryData)
	require.Equal(t, []byte("new-keyset"), manager.recoveryData.KeySetCiphertext)

	// The manager's reply must report its own request timeout: it is what the
	// executor validates its guest execution bound against.
	require.Equal(t, int64(commParams.RequestTimeoutSec)*1000, executor.reportedTimeoutMs,
		"GetKeysetRecovery reply must carry the manager's request timeout in ms")
}

func TestHandshake_Reconnection(t *testing.T) {
	ctx, client, server, executor, manager := setupHandshakeTest(t)
	defer client.Close()
	defer server.Stop()

	// Pre-set recovery data in manager
	manager.recoveryData = &common.EnclaveKeySetRecovery{
		RecoveryType:       common.RecoveryTypeKMS,
		KeySetCiphertext:   []byte("existing-keyset"),
		RecoveryCiphertext: []byte("existing-recovery"),
	}

	err := client.Connect(ctx, "Manager")
	require.NoError(t, err)

	// Wait for handshake to complete
	<-executor.handshakeDone

	require.NoError(t, executor.handshakeError)

	time.Sleep(1 * time.Second)
	require.True(t, manager.wasHandshakeSuccessful())
}

func TestHandshake_ManagerGetError(t *testing.T) {
	ctx, client, server, executor, manager := setupHandshakeTest(t)
	defer client.Close()
	defer server.Stop()

	// Manager returns error on get
	manager.getRecoveryError = fmt.Errorf("Test Get Error")

	err := client.Connect(ctx, "Manager")
	require.NoError(t, err)

	// Wait for handshake to complete
	<-executor.handshakeDone

	// Executor should see an error, because it is not a NotFound error
	require.Error(t, executor.handshakeError)
	require.Contains(t, executor.handshakeError.Error(), "Test Get Error")
	require.Nil(t, manager.recoveryData)
}

func TestHandshake_ManagerSetError(t *testing.T) {
	ctx, client, server, executor, manager := setupHandshakeTest(t)
	defer client.Close()
	defer server.Stop()

	// Manager returns error on get, this is OK, the executor will create a new keyset
	manager.getRecoveryError = storageErrors.NewError(storageErrors.NotFound, "Test Error")

	// Manager returns error on set
	manager.setRecoveryError = fmt.Errorf("Test Set Error")

	err := client.Connect(ctx, "Manager")
	require.NoError(t, err)

	// Wait for handshake to complete
	<-executor.handshakeDone

	require.Error(t, executor.handshakeError)
	require.Contains(t, executor.handshakeError.Error(), "Test Set Error")
	require.Nil(t, manager.recoveryData)
}

func TestHandshake_ExecutorRestoreFailure(t *testing.T) {
	ctx, client, server, executor, manager := setupHandshakeTest(t)
	defer client.Close()
	defer server.Stop()

	// Pre-set corrupted recovery data in manager
	manager.recoveryData = &common.EnclaveKeySetRecovery{
		RecoveryType:       common.RecoveryTypeKMS,
		KeySetCiphertext:   []byte("corrupted-keyset"),
		RecoveryCiphertext: []byte("corrupted-recovery"),
	}

	err := client.Connect(ctx, "Manager")
	require.NoError(t, err)

	// Wait for handshake to complete
	<-executor.handshakeDone

	// Executor should see an error
	require.Error(t, executor.handshakeError)
	require.Contains(t, executor.handshakeError.Error(), "simulated restore error")
}
