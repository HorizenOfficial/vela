package manager

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/admin"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/communication"
	cryptos "github.com/horizen-pes/pkg/crypto"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/logserver"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/stretchr/testify/require"
)

var testLogger logger.Logger

func TestMain(m *testing.M) {
	// Initialize once, by default it writes on stderr
	//testLogger = logger.NewLogger(&logger.Config{Kind: "printf"})
	//cfg := logger.DefaultLogConfig("zerolog")
	//testLogger = logger.NewLogger(&cfg)
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
	/*
		testLogger = logger.NewLogger(
			&logger.Config{
				Kind:         "zeronetwork",
				ConsoleLevel: "trace",
				// use a non-default port otherwise we can have bind failures if running tests concurrently
				RemoteLogParams:  common.TcpChannelConnectionParams{Ip: "localhost", Port: 5001},
				RemoteLogNetwork: "tcp",
				NetworkLevel:     "trace"},
		)
	*/

	// Run tests
	code := m.Run()
	os.Exit(code)
}

var (
	ApplicationId = common.NewApplicationId(1)
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
		return f.(func(context.Context, string) error)(ctx, tag)
	}
	return nil
}

func (m *MockExecutorClient) Close() error {
	if f, ok := m.GetMockedFunc("Close"); ok {
		return f.(func() error)()
	}
	return nil
}

func (m *MockExecutorClient) SendDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error) {
	if f, ok := m.GetMockedFunc("SendDeployApp"); ok {
		return f.(func(context.Context, *common.Request, *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error))(ctx, req, appState)
	}

	if req.ApplicationID != ApplicationId {	
		return nil, nil, fmt.Errorf("application id %s is not admitted", req.ApplicationID)
	}

	if req.RequestType != common.Deploy  {
		return nil, nil,fmt.Errorf("wrong request type: %d", req.RequestType)

	}

	if appState != nil {	
		failurePayload := &common.UpdatePayload{
			ApplicationID: req.ApplicationID, 
			RequestID: req.RequestID, 
			PrevStateRoot: appState.StateRoot, 
			NewStateRoot: appState.StateRoot,
			ErrorCode: uint8(apperrors.CodeApplicationAlreadyDeployed.Category.Category),
			ErrorMsg: fmt.Sprintf("application %s was already deployed", req.ApplicationID),
		}	

		return failurePayload, nil, nil
	}
	stateRoot := m.generateRandomStateRoot()
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, NewStateRoot: stateRoot}, &common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
}

func (m *MockExecutorClient) SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
	if f, ok := m.GetMockedFunc("SendProcessRequest"); ok {
		return f.(func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error))(ctx, req, appState, wasmModule)
	}

	if req.ApplicationID != ApplicationId {	
		return nil, nil, nil, fmt.Errorf("application id %s is not admitted", req.ApplicationID)
	}

	if req.RequestType != common.Process && req.RequestType != common.Deanonymize && req.RequestType != common.AssociateKey {
		return nil, nil, nil, fmt.Errorf("unsupported request type: %d", req.RequestType)

	}

	if appState == nil {	
		failurePayload := &common.UpdatePayload{
			ApplicationID: req.ApplicationID, 
			RequestID: req.RequestID, 
			PrevStateRoot: [32]byte{}, 
			NewStateRoot: [32]byte{},
			ErrorCode: uint8(apperrors.CodeAppStateNotFound.Category.Category),
			ErrorMsg: "application state not found",
			}	
		return failurePayload, nil, nil, nil
	}

	if string(req.Payload) == "invalid" {
		failurePayload := &common.UpdatePayload{
			ApplicationID: req.ApplicationID, 
			RequestID: req.RequestID, 
			PrevStateRoot: [32]byte{}, 
			NewStateRoot: [32]byte{},
			ErrorCode: uint8(apperrors.CategoryInternalMeta.Category),
			ErrorMsg: "invalid request",
			}	
		return failurePayload, nil, nil, nil
	}

	stateRoot := m.generateRandomStateRoot()

	// If this is a deanonymization request, return a deanonymization report
	var report *common.DeanonymizationReport
	if req.RequestType == common.Deanonymize {
		report = &common.DeanonymizationReport{
			ApplicationID:   req.ApplicationID,
			ReportID:        req.RequestID,
			Authority:       req.Sender,
			EncryptedReport: []byte(`{"accounts":{},"nonce":0}`),
		}
	}

	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, PrevStateRoot: appState.StateRoot, NewStateRoot: stateRoot},
		&common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, report, nil
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

var (
	sender = ethCommon.HexToAddress("0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199")
)

func createRequest(requestType common.RequestType, appID common.ApplicationIdType) *common.Request {
	requestId := testutil.GenerateRandomRequestID()

	request := &common.Request{ProtocolVersion: 0, ApplicationID: appID, RequestID: requestId, RequestType: requestType, Sender: sender, MaxFeeValue: common.NewBig(100)}
	return request
}

func createRequestWithPayload(requestType common.RequestType, appID common.ApplicationIdType, payload []byte) *common.Request {
	requestId := testutil.GenerateRandomRequestID()

	request := &common.Request{ProtocolVersion: 0, ApplicationID: appID, RequestID: requestId, RequestType: requestType, Sender: sender, Payload: payload, MaxFeeValue: common.NewBig(100)}
	return request
}

func TestStart(t *testing.T) {

	key, _ := cryptos.GeneratePrivateKeySecp256k1()
	config := &Config{
		HandshakeTimeout:          10,
		BlockchainPollingInterval: 10,
		PrivateKey:                *key,
		LogServerTCPAddress:       common.TcpChannelConnectionParams{Ip: "localhost", Port: 5001},
		LogServerLogFile:          "/tmp/temp.log",
	}
	stopChan := make(chan struct{})
	executorHandShake := ExecutorHandShake{
		isComplete: make(chan struct{}),
	}
	ctx := context.Background()
	mgrAlreadyStarted := false
	startLogServer := true
	bcClient, manager := setupTestWithConfig(t, ctx, *config, mgrAlreadyStarted, &executorHandShake, stopChan, startLogServer)
	require.False(t, manager.isRunning, "Manager should not be running initially")

	// Start the manager but execClient fails to connect
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func(context.Context, string) error {
		return fmt.Errorf("Connect failed")
	})
	err := manager.Start(ctx)
	require.Error(t, err, "Failed to connect to executor, should return error")
	require.False(t, manager.isRunning, "Manager should not be running after failed start")

	// Reset the executor client
	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("Connect")
	// Start the manager but blockchainClient fails to connect
	bcClient.AddMockedFunc("Connect", func(context.Context) error {
		return fmt.Errorf("failed to connect blockchain client")
	})

	// Mock successful executor client connection and handshake completion
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func(context.Context, string) error {
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

	time.Sleep(1 * time.Second)
	// Stopping the polling goroutine
	cancel()
	manager.wg.Wait()
	t.Log("TestStart completed")

}

func TestStop(t *testing.T) {
	key, _ := cryptos.GeneratePrivateKeySecp256k1()
	config := &Config{
		HandshakeTimeout:          10,
		BlockchainPollingInterval: 10,
		PrivateKey:                *key,
		LogServerTCPAddress:       common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	stopChan := make(chan struct{})
	executorHandShake := ExecutorHandShake{
		isComplete: make(chan struct{}),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bcClient, manager := setupTestWithConfig(t, ctx, *config, false, &executorHandShake, stopChan, false)
	require.False(t, manager.isRunning, "Manager should not be running initially")

	// Stop a manager that is not running
	err := manager.Stop()
	require.NoError(t, err, "Stopping a non-running manager should not return error")

	// Mock successful executor client connection and handshake completion
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func(context.Context, string) error {
		go manager.completeExecutorHandshake(nil)
		return nil
	})

	err = manager.Start(ctx)
	require.NoError(t, err, "Failed to start manager")

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
	mockBCClient, manager := setupTest(t)

	// Deploy request
	request := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
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
	request = createRequest(common.Process, ApplicationId)
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
	request = createRequest(common.Deanonymize, ApplicationId)
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

func TestProcessRequestsFromChainMixed(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Prepare different requests

	// Failure expected
	requestDeploy := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), requestDeploy)
	require.NoError(t, err)

	requestInvalid :=  createRequestWithPayload(common.Process, ApplicationId, []byte("invalid"))
	err = mockBCClient.SendRequestToChain(context.Background(), requestInvalid)
	require.NoError(t, err)

	requestReport := createRequest(common.Deanonymize, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), requestReport)
	require.NoError(t, err)

	// redeploy the same appId (failure expected)
	requestReDeploy := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), requestReDeploy)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 4, len(pendingRequests), "expected 4 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	for i := 0; i < 4; i++ {
		err = manager.processRequestFromChain(context.Background())
		require.NoError(t, err)
	}

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")

	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 2, len(failedRequests), "expected 2 failed request")

	require.Equal(t, requestInvalid.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, requestInvalid.ApplicationID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestInvalid.RequestType, failedRequests[0].RequestType, "Wrong RequestType")

	require.Equal(t, requestReDeploy.RequestID, failedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, requestReDeploy.ApplicationID, failedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestReDeploy.RequestType, failedRequests[1].RequestType, "Wrong RequestType")

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 4, len(completedRequests), "expected 4 completed request")

	// They should be in the same order of insertion
	require.Equal(t, requestDeploy.RequestID, completedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, requestDeploy.ApplicationID, completedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestDeploy.RequestType, completedRequests[0].RequestType, "Wrong RequestType")

	require.Equal(t, requestInvalid.RequestID, completedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, requestInvalid.ApplicationID, completedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestInvalid.RequestType, completedRequests[1].RequestType, "Wrong RequestType")

	require.Equal(t, requestReport.RequestID, completedRequests[2].RequestID, "Wrong requestID")
	require.Equal(t, requestReport.ApplicationID, completedRequests[2].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestReport.RequestType, completedRequests[2].RequestType, "Wrong RequestType")

	require.Equal(t, requestReDeploy.RequestID, completedRequests[3].RequestID, "Wrong requestID")
	require.Equal(t, requestReDeploy.ApplicationID, completedRequests[3].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, requestReDeploy.RequestType, completedRequests[3].RequestType, "Wrong RequestType")


}

func TestProcessDeployAppWithFailure(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	request := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Test that if it is a failure payload returned by the executor, submitStateOnChain is called but the state is not stored in the data layer 
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) error {
		t.Fatal("Store should not be called if the executor returned a failure payload")
		return nil
	})
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error) {
		return &common.UpdatePayload{ApplicationID: ApplicationId, 
			RequestID: req.RequestID, 
			ErrorCode: 1, 
			ErrorMsg: "error"}, nil, nil
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	
	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")
	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendDeployApp")

}

func TestProcessDeployAppWithErrors(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	request := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Test executor failure
	expectedError := "failed to deploy app"
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error) {
		return nil, nil, fmt.Errorf("%s", expectedError)
	})

	err = manager.processDeployApp(context.Background(), request)
	require.ErrorContains(t, err, expectedError)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendDeployApp")

	// Test data layer failure. In this case, it shouldn't call stateUpdate on chain and it returns the error  
	expectedError = "failed to store state"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) error {
		return fmt.Errorf("%s", expectedError)
	})

	failure := manager.processDeployApp(context.Background(), request)
	require.Error(t, failure)
	require.ErrorContains(t, failure, expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

	// Test blockchain failure for errors that can be due to reorgs. 
	// The errors that can be due to reorgs are:
	// - InvalidRequestId
	// - InvalidStateRoot
	// - InvalidApplicationId
	// - NonceTooLow
	// The local db should be reverted to the previous state
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return blockchain.ReorgError{}
	})

	failure = manager.processDeployApp(context.Background(), request)
	require.Nil(t, failure)
	// Check that the local db has been reverted to the initial state
	_, err = manager.dataLayer.LastVersionID()
	require.Error(t, err)
	dbErr, ok := err.(*storageErrors.Error)
	require.True(t, ok && dbErr.Code == storageErrors.NoVersionInDb)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	// Test blockchain failure for any other errors but reorgs. 
	// The local db should be reverted to the previous state
	expectedError = "some other error"
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return errors.New(expectedError)
	})

	err = manager.processDeployApp(context.Background(), request)
	require.Error(t, err)
	require.ErrorContains(t, err, expectedError)
	// Check that the local db has been reverted to the initial state
	_, err = manager.dataLayer.LastVersionID()
	require.Error(t, err)
	dbErr, ok = err.(*storageErrors.Error)
	require.True(t, ok && dbErr.Code == storageErrors.NoVersionInDb)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

}

func TestProcessProcessRequestWithFailure(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy the application first
	deployRequest := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	request := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Test that if it is a failure payload returned by the executor, submitStateOnChain is called but the state is not stored in the data layer 
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) error {
		t.Fatal("Store should not be called if the executor returned a failure payload")
		return nil
	})

	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendProcessRequest", 
	func(_ context.Context, req *common.Request, _ *common.ApplicationState,_ []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
		return &common.UpdatePayload{ApplicationID: ApplicationId, 
			RequestID: req.RequestID, 
			ErrorCode: 1, 
			ErrorMsg: "error"}, nil, nil, nil	})


	err = manager.processProcessRequest(context.Background(), request)
	require.NoError(t, err)
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")
	
	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")
	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendProcessRequest")

}

func TestProcessProcessRequestWithErrors(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy the application first
	deployRequest := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	oldDbVersion, err := manager.dataLayer.LastVersionID()
	require.NoError(t, err)

	// Simulate application state not found. In this case, it should call SendProcessRequest and return a failure payload, then submitStateOnChain is called but the state is not stored in the data layer
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, common.ApplicationIdType) (*common.ApplicationState, error) {
		return nil, storageErrors.ErrNotFound("application state not found")
	})
	request := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Failure in GetApplicationState. If the application wasn't already deployed, SendProcessRequest is called
	failure := manager.processProcessRequest(context.Background(), request)
	require.NoError(t, failure)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")
	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetApplicationState")

	//Other failures in GetApplicationState, stop processing and return the error
	request = createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	expectedError := "error"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, common.ApplicationIdType) (*common.ApplicationState, error) {
		return nil, errors.New(expectedError)
	})
	failure = manager.processProcessRequest(context.Background(), request)
	require.Error(t, failure)
	require.ErrorContains(t, failure, expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetApplicationState")

	// Failure in GetWasmCode, stop processing and return the error 
	expectedError = "wasm bytecode not found for application"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetWASMBytecode", func(context.Context, common.ApplicationIdType) ([]byte, error) {
		return nil, errors.New(expectedError)
	})
	failure = manager.processProcessRequest(context.Background(), request)
	require.Error(t, failure)
	require.ErrorContains(t, failure, expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetWASMBytecode")

	// Test failure in executor
	expectedError = "failed to execute app"
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendProcessRequest", func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
		return nil, nil, nil, errors.New(expectedError)
	})

	failure = manager.processProcessRequest(context.Background(), request)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendProcessRequest")

	// Test data layer failure, stop processing and return the error
	expectedError = "failed to store state"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) error {
		return errors.New(expectedError)
	})

	failure = manager.processProcessRequest(context.Background(), request)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

	// Test blockchain failure for errors that can be due to reorgs. 
	// The errors that can be due to reorgs are:
	// - InvalidRequestId
	// - InvalidStateRoot
	// - InvalidApplicationId
	// - NonceTooLow
	// The local db should be reverted to the previous state
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return blockchain.ReorgError{}
	})

	failure = manager.processProcessRequest(context.Background(), request)
	require.Nil(t, failure)
	// Check that the local db has been reverted to the initial state
	newDbVersion, err := manager.dataLayer.LastVersionID()
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed requests")

	// Test blockchain failure for any other errors but reorgs. 
	// The local db should be reverted to the previous state
	expectedError = "some other error"
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(ctx context.Context, payload *common.UpdatePayload) error {
		return errors.New(expectedError)
	})

	failure = manager.processProcessRequest(context.Background(), request)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	// Check that the local db has been reverted to the initial state
	newDbVersion, err = manager.dataLayer.LastVersionID()
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

}

// TestProcessDeanonymizationViaProcessRequest tests the deanonymization flow which is now
// handled through processProcessRequest with RequestType = Deanonymize
func TestProcessDeanonymizationViaProcessRequest(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy the application first
	deployRequest := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	// Create a Deanonymize request (now handled via processProcessRequest)
	request := createRequest(common.Deanonymize, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Process the deanonymization request via processRequestFromChain
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed requests")
	require.Equal(t, request.RequestID, completedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, common.Deanonymize, completedRequests[1].RequestType, "Wrong RequestType")
}

func TestProcessRequestFromChainWithReorgs(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Prepare initial state in the database
	_, initialStateRootOnChain, err := mockBCClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)

	// Execute some requests just to have different versions in the DB
	request1 := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), request1)
	require.NoError(t, err)

	request2 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request2)
	require.NoError(t, err)

	request3 := createRequest(common.Process, ApplicationId)
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
		t.Fatal("SubmitStateUpdate should not be called in case of reorg")
		return nil
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
	request4 := createRequest(common.Process, ApplicationId)
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
	mockBCClient, manager := setupTest(t)

	// Setup the application
	request := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
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
		t.Fatal("SubmitStateUpdate should not be called in case of reorg")
		return nil
	}
	mockBCClient.AddMockedFunc("SubmitStateUpdate", mockedSubmitStateUpdatePanics)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err, "processRequestFromChain should not return an error if GetNextPendingRequest fails")

	request1 := createRequest(common.Process, ApplicationId)
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

	//**********************
	// Check that if processRequest returns an error, processRequestFromChain doesn't execute the request and doesn't return an error
	//**********************

	// Setup a fake situation to make processRequest return an error
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, common.ApplicationIdType) (*common.ApplicationState, error) {
		return nil, fmt.Errorf("error")
	})

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err, "processRequestFromChain should not return an error if processRequest fails")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetApplicationState")

}

func setupTest(t *testing.T) (*blockchain.MockClient, *SecureProcessorManager) {
	config := Config{
		ReorgTimeout:        60,
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	return setupTestWithConfig(t, context.Background(), config, true, &ExecutorHandShake{}, nil, false)
}

func setupTestWithConfig(
	t *testing.T,
	ctx context.Context,
	config Config,
	managerIsRunning bool,
	executorHandShake *ExecutorHandShake,
	stopChan chan struct{},
	startLogServer bool,
) (*blockchain.MockClient, *SecureProcessorManager) {
	mockDataLayer := mockdb.NewMockDataLayer()
	bcClient := blockchain.NewMockClient()
	execClient := NewMockExecutorClient()
	tmpDir, err := os.MkdirTemp("", "reports")
	require.NoError(t, err)
	config.DeanonymizationReportPath = tmpDir

	processor := &SecureProcessorManager{
		config:            &config,
		executorClient:    execClient,
		blockchainClient:  bcClient,
		dataLayer:         mockDataLayer,
		isRunning:         managerIsRunning,
		executorHandShake: executorHandShake,
		stopChan:          stopChan,
		log:               testLogger,
	}

	if startLogServer {
		logserver.StartLogServer(
			ctx,
			logserver.LogServerConfig{
				TCPAddr:        config.LogServerTCPAddress,
				VSockAddr:      config.LogServerVSockAddress,
				LogFilePath:    config.LogServerLogFile,
				ConsoleEnabled: config.LogServerConsole,
				ConsoleLevel:   config.LogServerConsoleLevel,
				FileLevel:      config.LogServerFileLevel,
			},
		)
		time.Sleep(500 * time.Millisecond)
	}

	return bcClient, processor
}

func TestProcessDeanonymizationWithReportSaving(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Create a temporary directory for the reports
	tempDir, err := os.MkdirTemp("", "reports")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Deploy the application first
	deployRequest := createRequestWithPayload(common.Deploy, ApplicationId, []byte{0x01})
	err = mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	// DeanonymizationReportPath is set, so the report should be saved to the filesystem
	// Create a deanonymization request
	request := createRequest(common.Deanonymize, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)
	manager.config.DeanonymizationReportPath = tempDir
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")
	// Check that the report file exists
	reportFilePath := filepath.Join(tempDir, common.ReportFilename(request.ApplicationID, request.RequestID))
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
	require.Equal(t, sender, report.Authority, "Report authority should match the request sender")
}

func TestGetAndSetLogLevel(t *testing.T) {
	_, manager := setupTest(t)
	ctx := context.Background()

	// 1. GetLogLevel - get current level (should be "trace" as configured in TestMain)
	getMsg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage}
	result, err := manager.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	initialLevel := result.(string)
	require.Equal(t, "trace", initialLevel, "Initial log level should be trace")

	// 2. SetLogLevel - change it to "error"
	setData, err := json.Marshal(struct {
		Level string `json:"level"`
	}{Level: "error"})
	require.NoError(t, err)
	setMsg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err = manager.ExecuteCommand(ctx, setMsg)
	require.NoError(t, err)
	// Verify the set response
	respBytes, err := json.Marshal(result)
	require.NoError(t, err)
	var setResp struct {
		Success bool   `json:"success"`
		Level   string `json:"level"`
	}
	err = json.Unmarshal(respBytes, &setResp)
	require.NoError(t, err)
	require.True(t, setResp.Success)
	require.Equal(t, "error", setResp.Level)

	// 3. GetLogLevel again - verify it changed to "error"
	result, err = manager.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	newLevel := result.(string)
	require.Equal(t, "error", newLevel, "Log level should be error after SetLogLevel")

	// 4. SetLogLevel with empty string - should fail
	setData, err = json.Marshal(struct {
		Level string `json:"level"`
	}{Level: ""})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err = manager.ExecuteCommand(ctx, setMsg)
	require.Error(t, err, "Empty log level should return an error")
	require.Nil(t, result)
	require.Contains(t, err.Error(), "must not be empty")

	// Verify level is still "error" (unchanged after failed set)
	result, err = manager.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	require.Equal(t, "error", result.(string), "Log level should remain error after failed SetLogLevel")

	// 5. SetLogLevel with invalid level - should fail
	setData, err = json.Marshal(struct {
		Level string `json:"level"`
	}{Level: "bogus"})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	result, err = manager.ExecuteCommand(ctx, setMsg)
	require.Error(t, err, "Invalid log level should return an error")
	require.Nil(t, result)

	// Restore original level for other tests
	setData, err = json.Marshal(struct {
		Level string `json:"level"`
	}{Level: initialLevel})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Data: setData}
	_, err = manager.ExecuteCommand(ctx, setMsg)
	require.NoError(t, err)
}
