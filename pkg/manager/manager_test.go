package manager

import (
	"context"
	"fmt"
	"testing"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/stretchr/testify/require"
)

type MockExecutorClient struct {
	failure bool
}

func (m *MockExecutorClient) Connect(ctx context.Context) error {
	if m.failure {
		return fmt.Errorf("failed to connect to executor")
	}

	return nil
}

func (m *MockExecutorClient) Close() error {
	if m.failure {
		return fmt.Errorf("failed to close executor")
	}
	return nil
}

func (m *MockExecutorClient) SendDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, error) {
	if m.failure {
		return nil, nil, fmt.Errorf("failed to deploy app")
	}
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID}, &common.ApplicationState{ApplicationID: req.ApplicationID}, nil
}

func (m *MockExecutorClient) SendGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error) {
	if m.failure {
		return nil, fmt.Errorf("failed to generate report")
	}
	return &common.DeanonymizationReport{ReportID: req.RequestID}, nil
}

func (m *MockExecutorClient) SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	if m.failure {
		return nil, nil, fmt.Errorf("failed to process request")
	}
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID}, appState, nil
}

func (m *MockExecutorClient) SetClientRequestHandler(handler communication.ClientRequestHandler) {

}

const (
	sender = "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199"
)

type MockBCClientWithError struct {
	blockchain.Client
}

func (c *MockBCClientWithError) Connect(ctx context.Context) error {
	return fmt.Errorf("failed to connect blockchain client")
}

func (c *MockBCClientWithError) Close() error {
	return fmt.Errorf("failed to close blockchain client")
}

type MockDataLayerWithFailure struct {
	storage.DataLayer
}

func (d *MockDataLayerWithFailure) Close() error {
	return fmt.Errorf("failed to close")
}

func TestStart(t *testing.T) {

	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := &MockExecutorClient{failure: false}

	manager := NewSecureProcessorManager(&Config{BlockchainPollingInterval: 10}, bcClient, mockDataLayer, execClient)
	require.False(t, manager.isRunning, "Manager should not be running initially")

	// Start the manager but execClient fails to connect
	manager.executorClient = &MockExecutorClient{failure: true}
	err := manager.Start(context.Background())
	require.Error(t, err, "Failed to connect to executor, should return error")
	require.False(t, manager.isRunning, "Manager should not be running after failed start")

	// Reset the executor client
	manager.executorClient = &MockExecutorClient{failure: false}

	// Start the manager but blockchainClient fails to connect
	manager.blockchainClient = &MockBCClientWithError{manager.blockchainClient}
	err = manager.Start(context.Background())
	require.Error(t, err, "failed to connect to blockchain, should return error")
	require.False(t, manager.isRunning, "Manager should not be running after failed start")

	// Reset the blockchain client
	manager.blockchainClient = blockchain.NewMockClient()

	ctx, cancel := context.WithCancel(context.Background())
	err = manager.Start(ctx)
	require.NoError(t, err, "Failed to start manager")

	require.True(t, manager.isRunning, "Manager should be running after start")

	err = manager.Start(context.Background())
	require.Error(t, err, "Manager is already started, should return error")

	// Stopping the polling goroutine
	cancel()
	manager.wg.Wait()

}

func TestStop(t *testing.T) {

	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := &MockExecutorClient{failure: false}

	manager := NewSecureProcessorManager(&Config{BlockchainPollingInterval: 10}, bcClient, mockDataLayer, execClient)
	require.False(t, manager.isRunning, "Manager should not be running initially")

	// Stop a manager that is not running
	err := manager.Stop()
	require.NoError(t, err, "Stopping a non-running manager should not return error")

	ctx, cancel := context.WithCancel(context.Background())

	err = manager.Start(ctx)
	require.NoError(t, err, "Failed to start manager")
	// Stopping the polling goroutine, otherwise Stop() will block forever
	cancel()

	// Stop the manager but execClient fails to stop
	manager.executorClient = &MockExecutorClient{failure: true}
	err = manager.Stop()
	require.Error(t, err, "Failed to stop executor, should return error")
	require.True(t, manager.isRunning, "Manager should be running after failed stop")

	// Reset the executor client
	manager.executorClient = &MockExecutorClient{failure: false}

	// Stop the manager but blockchainClient fails to stop
	manager.blockchainClient = &MockBCClientWithError{manager.blockchainClient}
	err = manager.Stop()
	require.Error(t, err, "Failed to stop executor, should return error")
	require.True(t, manager.isRunning, "Manager should be running after failed stop")

	// Reset the blockchain client
	manager.blockchainClient = blockchain.NewMockClient()

	// Stop the manager but DataLayer fails to stop
	manager.dataLayer = &MockDataLayerWithFailure{manager.dataLayer}
	err = manager.Stop()
	require.Error(t, err, "Failed to stop executor, should return error")
	require.True(t, manager.isRunning, "Manager should be running after failed stop")

	// Reset the blockchain client
	manager.dataLayer = mockdb.NewMockDataLayer()

	err = manager.Stop()
	require.NoError(t, err, "Failed to stop manager")
}

func TestProcessRequestsFromChain(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Deploy request
	request := createRequest(common.Deploy, "1")
	err := mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 0, len(failedRequests), "expected 0 failed request")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	require.Equal(t, request.RequestID, completedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[0].RequestType, "Wrong RequestType")

	// reset all
	mockBCClient.ClearAllData()

	// Process request
	request = createRequest(common.Process, "1")
	err = mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	failedRequests = mockBCClient.GetFailedRequests()
	require.Equal(t, 0, len(failedRequests), "expected 0 failed request")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	require.Equal(t, request.RequestID, completedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[0].RequestType, "Wrong RequestType")

	// reset all
	mockBCClient.ClearAllData()

	// Deanonymize request
	request = createRequest(common.Deanonymize, "1")
	err = mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	failedRequests = mockBCClient.GetFailedRequests()
	require.Equal(t, 0, len(failedRequests), "expected 0 failed request")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	require.Equal(t, request.RequestID, completedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[0].RequestType, "Wrong RequestType")

}

func createRequest(requestType common.RequestType, appID string) *common.Request {
	requestId, err := blockchain.GenerateRandomID()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random ID: %v", err))
	}
	request := &common.Request{ProtocolVersion: "1.0", ApplicationID: appID, RequestID: requestId, RequestType: requestType, Sender: sender}
	return request
}

func TestProcessRequestsFromChainWithFailure(t *testing.T) {
	mockBCClient, manager := setupTest()

	manager.executorClient = &MockExecutorClient{failure: true}

	// Deploy request
	request := createRequest(common.Deploy, "1")
	err := mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

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

	// Process request
	request = createRequest(common.Process, "1")
	err = mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

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
	err = mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

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
	err = mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

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

func TestProcessRequestsFromChainMixed(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Prepare different requests

	// Failure expected
	requestInvalid := createRequest("invalidType", "1")
	err := mockBCClient.SubmitRequest(context.Background(), requestInvalid)
	require.NoError(t, err)

	requestDeploy := createRequest(common.Deploy, "1")
	err = mockBCClient.SubmitRequest(context.Background(), requestDeploy)
	require.NoError(t, err)

	requestReport := createRequest(common.Deanonymize, "1")
	err = mockBCClient.SubmitRequest(context.Background(), requestReport)
	require.NoError(t, err)

	// redeploy the same appId (failure expected)
	requestReDeploy := createRequest(common.Deploy, "1")
	err = mockBCClient.SubmitRequest(context.Background(), requestReDeploy)
	require.NoError(t, err)

	// deploy an app with an appID other than "1" (failure expected)
	// TODO it will change in future
	requestDeployWrongId := createRequest(common.Deploy, "33")
	err = mockBCClient.SubmitRequest(context.Background(), requestDeployWrongId)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 5, len(pendingRequests), "expected 5 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.processRequestsFromChain(context.Background())

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

func setupTest() (*blockchain.MockClient, *SecureProcessorManager) {
	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := &MockExecutorClient{failure: false}

	processor := &SecureProcessorManager{executorClient: execClient,
		blockchainClient: bcClient,
		dataLayer:        mockDataLayer,
		isRunning:        true}

	return bcClient, processor
}
