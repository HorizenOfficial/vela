package manager

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

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
	stateRoot := m.generateRandomStateRoot()
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, NewStateRoot: stateRoot}, &common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
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
	stateRoot := m.generateRandomStateRoot()
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, NewStateRoot: stateRoot}, &common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
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

func createRequest(requestType common.RequestType) *common.Request {
	requestId, err := blockchain.GenerateRandomID()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random ID: %v", err))
	}
	request := &common.Request{ProtocolVersion: "1.0", ApplicationID: "app1", RequestID: requestId, RequestType: requestType, Sender: sender}
	return request
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
	require.Error(t, err, "Manager his already started, should return error")

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
	request := createRequest(common.Deploy)
	err := mockBCClient.SubmitRequest(context.Background(), request)
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
	request = createRequest(common.Process)
	err = mockBCClient.SubmitRequest(context.Background(), request)
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
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")
	require.Equal(t, request.RequestID, completedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[1].RequestType, "Wrong RequestType")

	// Deanonymize request
	request = createRequest(common.Deanonymize)
	err = mockBCClient.SubmitRequest(context.Background(), request)
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
	require.Equal(t, 3, len(completedRequests), "expected 3 completed request")
	require.Equal(t, request.RequestID, completedRequests[2].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[2].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, request.RequestType, completedRequests[2].RequestType, "Wrong RequestType")

}

func TestProcessRequestsFromChainWithFailure(t *testing.T) {
	mockBCClient, manager := setupTest()

	manager.executorClient = &MockExecutorClient{failure: true}

	// Deploy request
	request := createRequest(common.Deploy)
	err := mockBCClient.SubmitRequest(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

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

	// Process request
	request = createRequest(common.Process)
	err = mockBCClient.SubmitRequest(context.Background(), request)
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

	// reset all
	mockBCClient.ClearAllData()

	// Deanonymize request
	request = createRequest(common.Deanonymize)
	err = mockBCClient.SubmitRequest(context.Background(), request)
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

	// Invalid request type
	// reset all
	mockBCClient.ClearAllData()
	request = createRequest("invalidType")
	err = mockBCClient.SubmitRequest(context.Background(), request)
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

func TestProcessRequestsFromChainWithReorgs(t *testing.T) {
	mockBCClient, manager := setupTest()

	// Prepare initial state in the database
	_, initialStateRootOnChain, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)

	// Execute some requests just to have different versions in the DB
	request1 := createRequest(common.Deploy)
	err = mockBCClient.SubmitRequest(context.Background(), request1)
	require.NoError(t, err)

	request2 := createRequest(common.Process)
	err = mockBCClient.SubmitRequest(context.Background(), request2)
	require.NoError(t, err)

	request3 := createRequest(common.Process)
	err = mockBCClient.SubmitRequest(context.Background(), request3)
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
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

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
	mockedSubmitStateUpdate := func(context.Context, *common.UpdatePayload) error {
		panic("SubmitStateUpdate should not be called in case of reorg")
	}
	mockBCClient.AddMockedFunc("SubmitStateUpdate", mockedSubmitStateUpdate)

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
	require.Equal(t, 3, len(completedRequests), "expected 3 completed request")

	db_version, err = manager.dataLayer.LastVersionID()
	require.NoError(t, err)

	_, stateRootOnChain3, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRootOnChain3[:], db_version), "State root in DB should be equal to state root on chain")

	// test unrecoverable disalignment between DB and chain
	request4 := createRequest(common.Process)
	err = mockBCClient.SubmitRequest(context.Background(), request4)
	require.NoError(t, err)

	mockedGetNextPendingRequest = func(context.Context) (*common.Request, [32]byte, error) {
		return request4, [32]byte{0x11, 0x66}, nil
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedGetNextPendingRequest)

	err = manager.processRequestFromChain(context.Background())
	require.Error(t, err, "Should return error due to unrecoverable disalignment between DB and chain")

	// test reorg not solved within timeout
	manager.config.ReorgTimeout = 1 // 1 second

	mockedGetNextPendingRequest = func(context.Context) (*common.Request, [32]byte, error) {
		return request1, initialStateRootOnChain, nil
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedGetNextPendingRequest)

	// SubmitStateUpdate should not be called in case of reorg
	mockedSubmitStateUpdate = func(context.Context, *common.UpdatePayload) error {
		panic("SubmitStateUpdate should not be called in case of reorg")
	}
	mockBCClient.AddMockedFunc("SubmitStateUpdate", mockedSubmitStateUpdate)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	// wait for more than reorg timeout
	// Instead of sleeping, we will simulate the time.Sleep by manipulating the endReorgTime
	manager.endReorgTime = manager.endReorgTime.Add(-2 * time.Second) // go back in time by 2 seconds
	err = manager.processRequestFromChain(context.Background())
	require.Error(t, err, "Should return error due to reorg not solved within timeout")

	// Check that, even if the timeout has been reached, if the state roots match, the reorg is resolved
	// Solve the reorg and process the last request
	mockBCClient.RemoveMockedFunc("GetNextPendingRequest")
	mockBCClient.RemoveMockedFunc("SubmitStateUpdate")
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRootOnChain3[:], db_version), "State root in DB should be equal to state root on chain")

}

func setupTest() (*blockchain.MockClient, *SecureProcessorManager) {
	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := &MockExecutorClient{failure: false}

	processor := &SecureProcessorManager{
		config: &Config{ReorgTimeout: 60},
		executorClient: execClient,
		blockchainClient: bcClient,
		dataLayer:        mockDataLayer,
		isRunning:        true}

	return bcClient, processor
}
