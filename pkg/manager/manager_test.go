package manager

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/communication"
	cryptos "github.com/horizen-pes/pkg/crypto"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/stretchr/testify/require"
)

type MockExecutorClient struct {
	*testutil.MockFunctions
}

// HandShake implements communication.ExecutorClient.
func (m *MockExecutorClient) HandShake(ctx context.Context, message string) (string, error) {
	panic("unimplemented")
}

func NewMockExecutorClient() *MockExecutorClient {
	return &MockExecutorClient{MockFunctions: testutil.NewMockFunctions()}
}

func (m *MockExecutorClient) Connect(ctx context.Context, tag string) error {
	if f, ok := m.GetMockedFunc("Connect"); ok {
		return f.(func() error)()
	}
	return nil
}

func (m *MockExecutorClient) Close() error {
	if f, ok := m.GetMockedFunc("Close"); ok {
		return f.(func() error)()
	}
	return nil
}

func (m *MockExecutorClient) SendDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, error) {
	if f, ok := m.GetMockedFunc("SendDeployApp"); ok {
		return f.(func(context.Context, *common.Request) (*common.UpdatePayload, *common.ApplicationState, error))(ctx, req)
	}
	stateRoot := m.generateRandomStateRoot()
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, NewStateRoot: stateRoot}, &common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
}

func (m *MockExecutorClient) SendGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error) {
	if f, ok := m.GetMockedFunc("SendGenerateDeanonymizationReport"); ok {
		return f.(func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.DeanonymizationReport, error))(ctx, req, appState, wasmModule)
	}
	return &common.DeanonymizationReport{ApplicationID: req.ApplicationID, ReportID: req.RequestID}, nil
}

func (m *MockExecutorClient) SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	if f, ok := m.GetMockedFunc("SendProcessRequest"); ok {
		return f.(func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.UpdatePayload, *common.ApplicationState, error))(ctx, req, appState, wasmModule)
	}

	stateRoot := m.generateRandomStateRoot()
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, PrevStateRoot: appState.StateRoot, NewStateRoot: stateRoot},
		&common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
}

func (m *MockExecutorClient) SetClientRequestHandler(handler communication.ClientRequestHandler) {

}

// GenerateRandomID generates a random ID
func (m *MockExecutorClient) generateRandomStateRoot() [32]byte {
	var b [32]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic("failed to generate random ID")
	}

	return b
}

const (
	sender = "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199"
)

func createRequest(requestType common.RequestType, appID string) *common.Request {
	requestId, err := blockchain.GenerateRandomID()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random ID: %v", err))
	}
	request := &common.Request{ProtocolVersion: "1.0", ApplicationID: appID, RequestID: requestId, RequestType: requestType, Sender: sender}
	return request
}

func createRequestWithPayload(requestType common.RequestType, appID string, payload []byte) *common.Request {
	requestId, err := blockchain.GenerateRandomID()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random ID: %v", err))
	}
	request := &common.Request{ProtocolVersion: "1.0", ApplicationID: appID, RequestID: requestId, RequestType: requestType, Sender: sender, Payload: payload}
	return request
}

func TestStart(t *testing.T) {
	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := NewMockExecutorClient()
	key, _ := cryptos.GeneratePrivateKeySecp256k1()
	manager := NewSecureProcessorManager(&Config{HandshakeTimeout: 10, BlockchainPollingInterval: 10, PrivateKey: *key}, bcClient, mockDataLayer, execClient)
	require.False(t, manager.isRunning, "Manager should not be running initially")

	// Start the manager but execClient fails to connect
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func() error {
		return fmt.Errorf("Connect failed")
	})
	err := manager.Start(context.Background())
	require.Error(t, err, "Failed to connect to executor, should return error")
	require.False(t, manager.isRunning, "Manager should not be running after failed start")

	// Reset the executor client
	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("Connect")
	// Start the manager but blockchainClient fails to connect
	bcClient.AddMockedFunc("Connect", func(context.Context) error {
		return fmt.Errorf("failed to connect blockchain client")
	})

	// Mock successful executor client connection and handshake completion
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func() error {
		go manager.completeExecutorHandshake(nil)
		return nil
	})

	err = manager.Start(context.Background())
	require.Error(t, err, "failed to connect to blockchain, should return error")
	require.False(t, manager.isRunning, "Manager should not be running after failed start")

	// Reset the blockchain client
	bcClient.RemoveMockedFunc("Connect")
	ctx, cancel := context.WithCancel(context.Background())
	err = manager.Start(ctx)
	require.NoError(t, err, "Failed to start manager")

	require.True(t, manager.isRunning, "Manager should be running after start")
	err = manager.Start(context.Background())
	require.Error(t, err, "Manager is already started, should return error")

	// Stopping the polling goroutine
	cancel()
	manager.wg.Wait()
	t.Log("TestStart completed")

}

func TestStop(t *testing.T) {

	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := NewMockExecutorClient()
	key, _ := cryptos.GeneratePrivateKeySecp256k1()
	manager := NewSecureProcessorManager(&Config{HandshakeTimeout: 10, BlockchainPollingInterval: 10, PrivateKey: *key}, bcClient, mockDataLayer, execClient)
	require.False(t, manager.isRunning, "Manager should not be running initially")

	// Stop a manager that is not running
	err := manager.Stop()
	require.NoError(t, err, "Stopping a non-running manager should not return error")

	ctx, cancel := context.WithCancel(context.Background())

	// Mock successful executor client connection and handshake completion
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func() error {
		go manager.completeExecutorHandshake(nil)
		return nil
	})

	err = manager.Start(ctx)
	require.NoError(t, err, "Failed to start manager")
	// Stopping the polling goroutine, otherwise Stop() will block forever
	cancel()

	// Stop the manager but execClient fails to stop
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Close", func() error {
		return fmt.Errorf("Close failed")
	})
	err = manager.Stop()
	require.Error(t, err, "Failed to stop executor, should return error")
	require.True(t, manager.isRunning, "Manager should be running after failed stop")

	// Reset the executor client
	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("Close")

	// Stop the manager but blockchainClient fails to stop
	bcClient.AddMockedFunc("Close", func() error { return fmt.Errorf("failed to close blockchain client") })

	err = manager.Stop()
	require.Error(t, err, "Failed to stop executor, should return error")
	require.True(t, manager.isRunning, "Manager should be running after failed stop")

	// Reset the blockchain client
	bcClient.RemoveMockedFunc("Close")

	// Stop the manager but DataLayer fails to stop
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Close", func() error { return fmt.Errorf("failed to close data layer") })

	err = manager.Stop()
	require.Error(t, err, "Failed to stop executor, should return error")
	require.True(t, manager.isRunning, "Manager should be running after failed stop")

	// Reset the blockchain client
	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Close")

	err = manager.Stop()
	require.NoError(t, err, "Failed to stop manager")
}

func TestProcessRequestFromChain(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Deploy request
	request := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 0, len(failedRequests), "expected 0 failed request")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	require.Equal(t, request.RequestID, completedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[0].RequestType, "Wrong RequestType")

	// Process request
	request = createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	failedRequests = mockBCClient.GetFailedRequests()
	require.Equal(t, 0, len(failedRequests), "expected 0 failed request")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed requests")
	require.Equal(t, request.RequestID, completedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[1].RequestType, "Wrong RequestType")

	// Deanonymize request
	request = createRequest(common.Deanonymize, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 0 completed request")

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	failedRequests = mockBCClient.GetFailedRequests()
	require.Equal(t, 0, len(failedRequests), "expected 0 failed request")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 3, len(completedRequests), "expected 3 completed requests")
	require.Equal(t, request.RequestID, completedRequests[2].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[2].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[2].RequestType, "Wrong RequestType")

}

func TestMarkRequestFailed(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Deploy request
	request := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, error) {
		return nil, nil, fmt.Errorf("failed to deploy app")
	})
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, failedRequests[0].RequestType, "Wrong RequestType")

	// reset all
	mockBCClient.ClearAllData()
	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendDeployApp")

	// Process request
	request = createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	// Should fail because there is no appplication already deployed
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	failedRequests = mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, failedRequests[0].RequestType, "Wrong RequestType")

	// reset all
	mockBCClient.ClearAllData()

	// Deanonymize request
	request = createRequest(common.Deanonymize, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	// Should fail because there is no appplication already deployed
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	failedRequests = mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, failedRequests[0].RequestType, "Wrong RequestType")

	// Invalid request type
	// reset all
	mockBCClient.ClearAllData()
	request = createRequest("invalidType", "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	failedRequests = mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, failedRequests[0].RequestType, "Wrong RequestType")

}

func TestMarkRequestFailedWithError(t *testing.T) {
	mockBCClient, manager := setupTest()

	request := createRequest("invalidType", "1")
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	mockBCClient.AddMockedFunc("MarkRequestFailed", func(ctx context.Context, requestID string) error {
		return fmt.Errorf("failed to mark request as failed")
	})

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expect request still pending since MarkRequestFailed failed")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 0, len(failedRequests), "expected 0 failed request")

}

func TestProcessRequestsFromChainMixed(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Prepare different requests

	// Failure expected
	requestInvalid := createRequest("invalidType", "1")
	err := mockBCClient.SendRequestToChain(context.Background(), requestInvalid)
	require.NoError(t, err)

	requestDeploy := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), requestDeploy)
	require.NoError(t, err)

	requestReport := createRequest(common.Deanonymize, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), requestReport)
	require.NoError(t, err)

	// redeploy the same appId (failure expected)
	requestReDeploy := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), requestReDeploy)
	require.NoError(t, err)

	// deploy an app with an appID other than "1" (failure expected)
	// TODO it will change in future
	requestDeployWrongId := createRequestWithPayload(common.Deploy, "33", []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), requestDeployWrongId)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 5, len(pendingRequests), "expected 5 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	for i := 0; i < 5; i++ {
		err = manager.processRequestFromChain(context.Background())
		require.NoError(t, err)
	}

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")

	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 3, len(failedRequests), "expected 3 failed request")

	require.Equal(t, requestInvalid.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, requestInvalid.ApplicationID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestInvalid.RequestType, failedRequests[0].RequestType, "Wrong RequestType")

	require.Equal(t, requestReDeploy.RequestID, failedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, requestReDeploy.ApplicationID, failedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestReDeploy.RequestType, failedRequests[1].RequestType, "Wrong RequestType")

	require.Equal(t, requestDeployWrongId.RequestID, failedRequests[2].RequestID, "Wrong requestID")
	require.Equal(t, requestDeployWrongId.ApplicationID, failedRequests[2].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestDeployWrongId.RequestType, failedRequests[2].RequestType, "Wrong RequestType")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 5, len(completedRequests), "expected 5 completed request")

	// They should be in the same order of insertion
	require.Equal(t, requestInvalid.RequestID, completedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, requestInvalid.ApplicationID, completedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestInvalid.RequestType, completedRequests[0].RequestType, "Wrong RequestType")

	require.Equal(t, requestDeploy.RequestID, completedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, requestDeploy.ApplicationID, completedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestDeploy.RequestType, completedRequests[1].RequestType, "Wrong RequestType")

	require.Equal(t, requestReport.RequestID, completedRequests[2].RequestID, "Wrong requestID")
	require.Equal(t, requestReport.ApplicationID, completedRequests[2].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestReport.RequestType, completedRequests[2].RequestType, "Wrong RequestType")

	require.Equal(t, requestReDeploy.RequestID, completedRequests[3].RequestID, "Wrong requestID")
	require.Equal(t, requestReDeploy.ApplicationID, completedRequests[3].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestReDeploy.RequestType, completedRequests[3].RequestType, "Wrong RequestType")

	require.Equal(t, requestDeployWrongId.RequestID, completedRequests[4].RequestID, "Wrong requestID")
	require.Equal(t, requestDeployWrongId.ApplicationID, completedRequests[4].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestDeployWrongId.RequestType, completedRequests[4].RequestType, "Wrong RequestType")

}

func TestProcessDeployAppWithErrors(t *testing.T) {
	mockBCClient, manager := setupTest()

	request := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Test executor failure
	expectedError := "failed to deploy app"
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, error) {
		return nil, nil, fmt.Errorf("%s", expectedError)
	})

	err = manager.processDeployApp(context.Background(), request)
	require.ErrorContains(t, err, expectedError)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendDeployApp")

	// Test data layer failure. In this case, it shouldn't return an error (son MarkFailed is not called) but it shouldn't call stateUpdate on chain either
	expectedError = "failed to store state"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) error {
		return fmt.Errorf("%s", expectedError)
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

	// Test blockchain failure for errors that can be due to reorgs. In this case the request shouldn't be marked as failed.
	// The errors that can be due to reorgs are:
	// - InvalidRequestId
	// - InvalidStateRoot
	// - InvalidApplicationId
	// - NonceTooLow
	// The local db should be reverted to the previous state
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return blockchain.ReorgError{}
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)
	// Check that the local db has been reverted to the initial state
	_, err = manager.dataLayer.LastVersionID()
	require.Error(t, err)
	dbErr, ok := err.(*storageErrors.Error)
	require.True(t, ok && dbErr.Code == storageErrors.NoVersionInDb)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	// Test blockchain failure for any other errors but reorgs. In this case the request should be marked as failed.
	// The local db should be reverted to the previous state
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return fmt.Errorf("some other error")
	})

	err = manager.processDeployApp(context.Background(), request)
	require.Error(t, err)
	// Check that the local db has been reverted to the initial state
	_, err = manager.dataLayer.LastVersionID()
	require.Error(t, err)
	dbErr, ok = err.(*storageErrors.Error)
	require.True(t, ok && dbErr.Code == storageErrors.NoVersionInDb)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

}

func TestProcessProcessRequestWithErrors(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Deploy the application first
	deployRequest := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	oldDbVersion, err := manager.dataLayer.LastVersionID()
	require.NoError(t, err)

	request := createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Failure in GetApplicationState. If the application wasn't already deployed, the request should be marked as failed
	request.ApplicationID = "invalid app"
	err = manager.processProcessRequest(context.Background(), request)
	require.Error(t, err)
	require.Contains(t, err.Error(), "application state not found")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	//Other failures in GetApplicationState, it may a temp error => the request shouldn't be marked as failed
	request.ApplicationID = deployRequest.ApplicationID

	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, string) (*common.ApplicationState, error) {
		return nil, fmt.Errorf("error")
	})
	err = manager.processProcessRequest(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetApplicationState")

	// Failure in GetWasmCode. In this case any error should be treated as temp error => the request shouldn't be marked as failed
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetWASMBytecode", func(context.Context, string) ([]byte, error) {
		return nil, fmt.Errorf("wasm bytecode not found for application")
	})
	err = manager.processProcessRequest(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetWASMBytecode")

	// Test failure in executor
	expectedError := "failed to execute app"
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendProcessRequest", func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		return nil, nil, fmt.Errorf("%s", expectedError)
	})

	err = manager.processProcessRequest(context.Background(), request)
	require.ErrorContains(t, err, expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendProcessRequest")

	// Test data layer failure. In this case, it shouldn't return an error (son MarkFailed is not called) but it shouldn't call stateUpdate on chain either
	expectedError = "failed to store state"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) error {
		return fmt.Errorf("%s", expectedError)
	})

	err = manager.processProcessRequest(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

	// Test blockchain failure for errors that can be due to reorgs. In this case the request shouldn't be marked as failed.
	// The errors that can be due to reorgs are:
	// - InvalidRequestId
	// - InvalidStateRoot
	// - InvalidApplicationId
	// - NonceTooLow
	// The local db should be reverted to the previous state
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return blockchain.ReorgError{}
	})

	err = manager.processProcessRequest(context.Background(), request)
	require.NoError(t, err)
	// Check that the local db has been reverted to the initial state
	newDbVersion, err := manager.dataLayer.LastVersionID()
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	// Test blockchain failure for any other errors but reorgs. In this case the request should be marked as failed.
	// The local db should be reverted to the previous state
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return fmt.Errorf("some other error")
	})

	err = manager.processProcessRequest(context.Background(), request)
	require.Error(t, err)
	// Check that the local db has been reverted to the initial state
	newDbVersion, err = manager.dataLayer.LastVersionID()
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

}

func TestProcessProcessDeanonymization(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Deploy the application first
	deployRequest := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	request := createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Failure in GetApplicationState. If the application wasn't already deployed, the request should be marked as failed
	request.ApplicationID = "invalid app"
	err = manager.processDeanonymization(context.Background(), request)
	require.Error(t, err)
	require.Contains(t, err.Error(), "application state not found")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	//Other failures in GetApplicationState, it may a temp error => the request shouldn't be marked as failed
	request.ApplicationID = deployRequest.ApplicationID

	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, string) (*common.ApplicationState, error) {
		return nil, fmt.Errorf("error")
	})
	err = manager.processDeanonymization(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetApplicationState")

	// Failure in GetWasmCode. In this case any error should be treated as temp error => the request shouldn't be marked as failed
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetWASMBytecode", func(context.Context, string) ([]byte, error) {
		return nil, fmt.Errorf("wasm bytecode not found for application")
	})
	err = manager.processDeanonymization(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetWASMBytecode")

	// Test failure in executor
	expectedError := "failed to execute app"
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendGenerateDeanonymizationReport", func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.DeanonymizationReport, error) {
		return nil, fmt.Errorf("%s", expectedError)
	})

	err = manager.processDeanonymization(context.Background(), request)
	require.ErrorContains(t, err, expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendGenerateDeanonymizationReport")

	// Test data layer failure. In this case, it shouldn't return an error (son MarkFailed is not called) but it shouldn't call stateUpdate on chain either
	expectedError = "failed to store report"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("StoreDeanonymizationReport", func(context.Context, *common.DeanonymizationReport) error {
		return fmt.Errorf("%s", expectedError)
	})

	err = manager.processDeanonymization(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("StoreDeanonymizationReport")

	// Test blockchain failures. For deanonymization, any error shouldn't mark the request as failed
	mockBCClient.AddMockedFunc("SubmitDeanonymizationReport", func(context.Context, *common.DeanonymizationReport) error {
		return fmt.Errorf("some other error")
	})

	err = manager.processDeanonymization(context.Background(), request)
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

}

func TestProcessRequestFromChainWithReorgs(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Prepare initial state in the database
	_, initialStateRootOnChain, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)

	// Execute some requests just to have different versions in the DB
	request1 := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), request1)
	require.NoError(t, err)

	request2 := createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request2)
	require.NoError(t, err)

	request3 := createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request3)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 3, len(pendingRequests), "expected 3 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 2, len(pendingRequests), "expected 2 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	db_version, err := manager.dataLayer.LastVersionID()
	require.NoError(t, err)

	nextPendingReq, stateRootOnChain1, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRootOnChain1[:], db_version), "State root in DB should be equal to state root on chain")
	require.Equal(t, request2.RequestID, nextPendingReq.RequestID)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed requests")

	db_version, err = manager.dataLayer.LastVersionID()
	require.NoError(t, err)

	nextPendingReq, stateRootOnChain2, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRootOnChain2[:], db_version), "State root in DB should be equal to state root on chain")
	require.Equal(t, request3.RequestID, nextPendingReq.RequestID)

	// Now simulate a reorg on chain, by making GetNextPendingRequest to always return the first request and initial state root

	mockedGetNextPendingRequest := func(context.Context) (*common.Request, [32]byte, error) {
		return request1, initialStateRootOnChain, nil
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedGetNextPendingRequest)

	// SubmitStateUpdate should not be called in case of reorg
	mockedSubmitStateUpdatePanics := func(context.Context, *common.UpdatePayload) error {
		panic("SubmitStateUpdate should not be called in case of reorg")
	}
	mockBCClient.AddMockedFunc("SubmitStateUpdate", mockedSubmitStateUpdatePanics)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	require.False(t, manager.endReorgTime.IsZero(), "endReorgTime should be set when reorg is detected")
	currentEndReorgTime := manager.endReorgTime

	// Try with the second request and state root
	mockedGetNextPendingRequest = func(context.Context) (*common.Request, [32]byte, error) {
		return request2, stateRootOnChain1, nil
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedGetNextPendingRequest)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	require.Equal(t, currentEndReorgTime, manager.endReorgTime, "endReorgTime should not change if reorg is not yet resolved")

	// Solve the reorg and process the last request
	mockBCClient.RemoveMockedFunc("GetNextPendingRequest")
	mockBCClient.RemoveMockedFunc("SubmitStateUpdate")

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	require.True(t, manager.endReorgTime.IsZero(), "endReorgTime should be reset when reorg is solved")

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 3, len(completedRequests), "expected 3 completed requests")

	db_version, err = manager.dataLayer.LastVersionID()
	require.NoError(t, err)

	_, stateRootOnChain3, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRootOnChain3[:], db_version), "State root in DB should be equal to state root on chain")

	// test unrecoverable disalignment between DB and chain
	request4 := createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request4)
	require.NoError(t, err)

	mockedGetNextPendingRequest = func(context.Context) (*common.Request, [32]byte, error) {
		return request4, [32]byte{0x11, 0x66}, nil
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedGetNextPendingRequest)

	err = manager.processRequestFromChain(context.Background())
	require.Error(t, err, "Should return error due to unrecoverable disalignment between DB and chain")

	// test reorg not solved within timeout. The local db is reverted to the same state of the chain and the request is executed
	//Remove all old requests and reset to the initial state root

	mockBCClient.ClearAllData()
	//Re-submit old request
	err = mockBCClient.SendRequestToChain(context.Background(), request1)
	require.NoError(t, err)

	manager.config.ReorgTimeout = 1 // 1 second

	// First processRequestFromChain sets the reorg timeout
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	// wait for more than reorg timeout
	// Instead of sleeping, we will simulate the time.Sleep by manipulating the endReorgTime
	manager.endReorgTime = manager.endReorgTime.Add(-2 * time.Second) // go back in time by 2 seconds

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	db_version, err = manager.dataLayer.LastVersionID()
	require.NoError(t, err)

	_, stateRootOnChain, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRootOnChain[:], db_version), "State root in DB should be equal to state root on chain")

}

func TestProcessRequestFromChainWithErrors(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Setup the application
	request := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	//**********************
	// Check that if GetNextPendingRequest returns an error, processRequestFromChain doesn't execute the request and doesn't return an error
	//**********************

	mockedGetNextPendingRequest := func(context.Context) (*common.Request, [32]byte, error) {
		return nil, [32]byte{}, fmt.Errorf("GetNextPendingRequest error")
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedGetNextPendingRequest)

	// SubmitStateUpdate should not be called in case of reorg
	mockedSubmitStateUpdatePanics := func(context.Context, *common.UpdatePayload) error {
		panic("SubmitStateUpdate should not be called in case of reorg")
	}
	mockBCClient.AddMockedFunc("SubmitStateUpdate", mockedSubmitStateUpdatePanics)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err, "processRequestFromChain should not return an error if GetNextPendingRequest fails")

	request1 := createRequest(common.Process, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request1)
	require.NoError(t, err)
	mockBCClient.RemoveMockedFunc("GetNextPendingRequest")

	//**********************
	// Check that if LastVersionID returns an error, processRequestFromChain doesn't execute the request and doesn't return an error
	//**********************

	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("LastVersionID", func() ([]byte, error) {
		return nil, fmt.Errorf("LastVersionID error")
	})

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err, "processRequestFromChain should not return an error if LastVersionID fails")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("LastVersionID")

	//**********************
	// Check that if ListVersions returns an error, processRequestFromChain doesn't execute the request and doesn't return an error
	//**********************
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("ListVersions", func() ([][]byte, error) {
		return nil, fmt.Errorf("ListVersions error")
	})

	// Setup a fake reorg situation
	mockedGetNextPendingRequest = func(context.Context) (*common.Request, [32]byte, error) {
		return request1, [32]byte{0x11, 0x66}, nil
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedGetNextPendingRequest)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err, "processRequestFromChain should not return an error if ListVersions fails")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("ListVersions")
	mockBCClient.RemoveMockedFunc("GetNextPendingRequest")

}

func setupTest() (*blockchain.MockClient, *SecureProcessorManager) {
	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := NewMockExecutorClient()

	processor := &SecureProcessorManager{
		config:           &Config{ReorgTimeout: 60},
		executorClient:   execClient,
		blockchainClient: bcClient,
		dataLayer:        mockDataLayer,
		isRunning:        true}

	return bcClient, processor
}

func TestProcessDeanonymizationWithReportSaving(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Create a temporary directory for the reports
	tempDir, err := os.MkdirTemp("", "reports")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Deploy the application first
	deployRequest := createRequestWithPayload(common.Deploy, "1", []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	// Case 1: DeanonymizationReportPath is not set, so the report should not be saved to the filesystem
	// Create a deanonymization request
	request := createRequest(common.Deanonymize, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)
	manager.config.DeanonymizationReportPath = ""
	err = manager.processDeanonymization(context.Background(), request)
	require.NoError(t, err)
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")
	// Check that the report file does not exist
	reportFilePath := filepath.Join(tempDir, request.ApplicationID+"_"+request.RequestID)
	_, err = os.Stat(reportFilePath)
	require.True(t, os.IsNotExist(err), "Report file should not exist when DeanonymizationReportPath is not set")
	// check we have it in the data layer
	storedReport, err := manager.dataLayer.GetDeanonymizationReport(context.Background(), request.RequestID)
	require.NoError(t, err)
	require.Equal(t, storedReport.ReportID, request.RequestID)

	// Case 2: DeanonymizationReportPath is set, so the report should be saved to the filesystem
	// Create a deanonymization request
	request = createRequest(common.Deanonymize, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)
	manager.config.DeanonymizationReportPath = tempDir
	err = manager.processDeanonymization(context.Background(), request)
	require.NoError(t, err)
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 3, len(completedRequests), "expected 3 completed request")
	// Check that the report file exists
	reportFilePath = filepath.Join(tempDir, request.ApplicationID+"_"+request.RequestID)
	_, err = os.Stat(reportFilePath)
	require.NoError(t, err, "Report file should exist when DeanonymizationReportPath is set")
	// Read the report from the filesystem and verify its contents
	reportBytes, err := os.ReadFile(reportFilePath)
	require.NoError(t, err)
	var report common.DeanonymizationReport
	err = json.Unmarshal(reportBytes, &report)
	require.NoError(t, err)
	require.Equal(t, request.RequestID, report.ReportID, "Report ID should match the request ID")
	require.Equal(t, request.ApplicationID, report.ApplicationID, "Report App ID should match the request App ID")
	// check we have it also in the data layer
	storedReport, err = manager.dataLayer.GetDeanonymizationReport(context.Background(), request.RequestID)
	require.NoError(t, err)
	require.Equal(t, storedReport.ReportID, request.RequestID)

	// Case 3: Error creating the directory
	// Create a deanonymization request
	request = createRequest(common.Deanonymize, "1")
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)
	// Set the path to a read-only directory to simulate an error
	readOnlyDir := filepath.Join(tempDir, "readonly")
	err = os.Mkdir(readOnlyDir, 0555)
	require.NoError(t, err)
	manager.config.DeanonymizationReportPath = filepath.Join(readOnlyDir, "reports")
	err = manager.processDeanonymization(context.Background(), request)
	require.NoError(t, err, "processDeanonymization should not return an error even if it fails to create the directory")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 4, len(completedRequests), "expected 4 completed request")
	// check we have it also in the data layer
	storedReport, err = manager.dataLayer.GetDeanonymizationReport(context.Background(), request.RequestID)
	require.NoError(t, err)
	require.Equal(t, storedReport.ReportID, request.RequestID)
}
