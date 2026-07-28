package manager

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/HorizenOfficial/vela/pkg/admin"
	"github.com/HorizenOfficial/vela/pkg/blockchain"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	"github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/communication"
	cryptos "github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/logserver"
	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
	"github.com/HorizenOfficial/vela/pkg/storage/mockdb"
	"github.com/HorizenOfficial/vela/pkg/version"
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

func (m *MockExecutorClient) IsConnected() bool {
	if f, ok := m.GetMockedFunc("IsConnected"); ok {
		return f.(func() bool)()
	}
	return true
}

func (m *MockExecutorClient) SendDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	if f, ok := m.GetMockedFunc("SendDeployApp"); ok {
		return f.(func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.UpdatePayload, *common.ApplicationState, error))(ctx, req, appState, wasmModule)
	}

	if req.RequestType != common.Deploy {
		return nil, nil, fmt.Errorf("wrong request type: %d", req.RequestType)

	}

	if appState != nil {
		failurePayload := &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			RequestID:     req.RequestID,
			PrevStateRoot: appState.StateRoot,
			NewStateRoot:  appState.StateRoot,
			ErrorCode:     uint8(apperrors.CodeApplicationAlreadyDeployed.Category.Category),
			ErrorMsg:      fmt.Sprintf("application %s was already deployed", req.ApplicationID),
		}

		return failurePayload, nil, nil
	}
	if len(wasmModule) == 0 {
		return &common.UpdatePayload{
			ApplicationID:  req.ApplicationID,
			RequestID:      req.RequestID,
			PrevStateRoot:  [32]byte{},
			NewStateRoot:   [32]byte{},
			ErrorCode:      uint8(apperrors.CodeFailedLoadingOrGettingModule.Category.Category),
			ErrorMsg:       "failed to load or get module",
			RefundAmount:   req.MaxFeeValue,
			ApplicationFee: common.NewBig(0),
		}, nil, nil
	}
	stateRoot := m.generateRandomStateRoot()
	return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, NewStateRoot: stateRoot}, &common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
}

func (m *MockExecutorClient) SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
	if f, ok := m.GetMockedFunc("SendProcessRequest"); ok {
		return f.(func(context.Context, *common.Request, *common.ApplicationState, []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error))(ctx, req, appState, wasmModule)
	}

	if req.RequestType != common.Process && req.RequestType != common.Deanonymize && req.RequestType != common.AssociateKey {
		return nil, nil, nil, fmt.Errorf("unsupported request type: %d", req.RequestType)

	}

	if appState == nil {
		failurePayload := &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			RequestID:     req.RequestID,
			PrevStateRoot: [32]byte{},
			NewStateRoot:  [32]byte{},
			ErrorCode:     uint8(apperrors.CodeAppStateNotFound.Category.Category),
			ErrorMsg:      "application state not found",
		}
		return failurePayload, nil, nil, nil
	}

	if string(req.Payload) == "invalid" {
		failurePayload := &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			RequestID:     req.RequestID,
			PrevStateRoot: [32]byte{},
			NewStateRoot:  [32]byte{},
			ErrorCode:     uint8(apperrors.CategoryInternalMeta.Category),
			ErrorMsg:      "invalid request",
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

func (m *MockExecutorClient) ForwardAdminCommand(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
	if f, ok := m.GetMockedFunc("ForwardAdminCommand"); ok {
		return f.(func(context.Context, string, json.RawMessage) (json.RawMessage, error))(ctx, cmdType, data)
	}
	return nil, fmt.Errorf("ForwardAdminCommand not mocked")
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

func createDeployRequestWithWASM(t *testing.T, manager *SecureProcessorManager, appID common.ApplicationIdType, wasm []byte) *common.Request {
	t.Helper()

	sum := sha256.Sum256(wasm)
	wasmSHA := hex.EncodeToString(sum[:])
	descriptorPayload := createDeployDescriptorPayload(t, wasmSHA)

	artifactBlobPath := filepath.Join(manager.config.ArtifactsPath, artifactBlobsFolder, wasmSHA+".wasm")
	require.NoError(t, os.MkdirAll(filepath.Dir(artifactBlobPath), 0o755))
	require.NoError(t, os.WriteFile(artifactBlobPath, wasm, 0o600))

	return createRequestWithPayload(common.Deploy, appID, descriptorPayload)
}

func createDeployDescriptorPayload(t *testing.T, wasmSHA string) []byte {
	t.Helper()

	artifactID, err := common.BuildArtifactID(wasmSHA)
	require.NoError(t, err)

	descriptor := common.DeployDescriptor{
		Mode:       common.DeployModeArtifactRef,
		ArtifactID: artifactID,
		WasmSHA256: wasmSHA,
	}
	payload, err := json.Marshal(descriptor)
	require.NoError(t, err)

	return payload
}

func writeArtifactBlob(t *testing.T, manager *SecureProcessorManager, wasmSHA string, wasm []byte) {
	t.Helper()
	artifactBlobPath := filepath.Join(manager.config.ArtifactsPath, artifactBlobsFolder, wasmSHA+".wasm")
	require.NoError(t, os.MkdirAll(filepath.Dir(artifactBlobPath), 0o755))
	require.NoError(t, os.WriteFile(artifactBlobPath, wasm, 0o600))
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
	// Start the manager when blockchainClient fails to connect — Start should
	// succeed anyway (non-fatal). The polling loop retries the connection.
	bcClient.AddMockedFunc("Connect", func(context.Context) error {
		return fmt.Errorf("failed to connect blockchain client")
	})

	// Mock successful executor client connection and handshake completion
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func(context.Context, string) error {
		go manager.completeExecutorHandshake(nil)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	err = manager.Start(ctx)
	require.NoError(t, err, "Start should succeed even if blockchain connect fails initially")
	require.True(t, manager.isRunning, "Manager should be running (blockchain will retry during polling)")

	// Stop so we can test a fresh start with successful blockchain connect
	cancel()
	manager.wg.Wait()
	manager.isRunning = false

	// Reset the blockchain client
	bcClient.RemoveMockedFunc("Connect")
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("Connect", func(context.Context, string) error {
		go manager.completeExecutorHandshake(nil)
		return nil
	})
	ctx, cancel = context.WithCancel(context.Background())
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
	request := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
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
	requestDeploy := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), requestDeploy)
	require.NoError(t, err)

	requestInvalid := createRequestWithPayload(common.Process, ApplicationId, []byte("invalid"))
	err = mockBCClient.SendRequestToChain(context.Background(), requestInvalid)
	require.NoError(t, err)

	requestReport := createRequest(common.Deanonymize, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), requestReport)
	require.NoError(t, err)

	// redeploy the same appId (failure expected)
	requestReDeploy := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
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

	request := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Test that if it is a failure payload returned by the executor, submitStateOnChain is called but the state is not stored in the data layer
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, *common.ApplicationState, *common.WASMData) error {
		t.Fatal("Store should not be called if the executor returned a failure payload")
		return nil
	})
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		return &common.UpdatePayload{ApplicationID: ApplicationId,
			RequestID: req.RequestID,
			ErrorCode: 1,
			ErrorMsg:  "error"}, nil, nil
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

	request := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Test executor failure
	expectedError := "failed to deploy app"
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		return nil, nil, fmt.Errorf("%s", expectedError)
	})

	err = manager.processDeployApp(context.Background(), request)
	require.ErrorContains(t, err, expectedError)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendDeployApp")

	// Test data layer failure. In this case, it shouldn't call stateUpdate on chain and it returns the error
	expectedError = "failed to store state"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, *common.ApplicationState, *common.WASMData) error {
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
	_, err = manager.dataLayer.LastVersionID(ApplicationId)
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
	_, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.Error(t, err)
	dbErr, ok = err.(*storageErrors.Error)
	require.True(t, ok && dbErr.Code == storageErrors.NoVersionInDb)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

}

func TestProcessDeployApp_InvalidDescriptorIsForwardedToExecutor(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	request := createRequestWithPayload(common.Deploy, ApplicationId, []byte("not-json"))
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	var capturedReq *common.Request
	var capturedWASM []byte
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		capturedReq = req
		capturedWASM = wasmModule
		return &common.UpdatePayload{
			ApplicationID:  req.ApplicationID,
			RequestID:      req.RequestID,
			PrevStateRoot:  [32]byte{},
			NewStateRoot:   [32]byte{},
			ErrorCode:      uint8(apperrors.CodeInternalFallback.Category.Category),
			ErrorMsg:       "failed to deploy application",
			RefundAmount:   req.MaxFeeValue,
			ApplicationFee: common.NewBig(0),
		}, nil, nil
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)

	require.NotNil(t, capturedReq)
	require.Equal(t, request.Payload, capturedReq.Payload)
	require.Nil(t, capturedWASM)

	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
}

func TestProcessDeployApp_IgnoresRequestSenderForAuthorization(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	request := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	request.Sender = ethCommon.HexToAddress("0x1111111111111111111111111111111111111111")
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		require.Equal(t, request.Sender, req.Sender)
		require.NotEmpty(t, wasmModule)
		stateRoot := sha256.Sum256([]byte("sender-agnostic-state"))
		return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, NewStateRoot: stateRoot}, &common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)

	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	require.Equal(t, 0, len(mockBCClient.GetFailedRequests()), "expected 0 failed requests")
}

func TestProcessDeployApp_MissingArtifactStaysPending(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	missingSHA := strings.Repeat("a", 64)
	request := createRequestWithPayload(common.Deploy, ApplicationId, createDeployDescriptorPayload(t, missingSHA))
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	var capturedWASM []byte
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		capturedWASM = wasmModule
		return &common.UpdatePayload{
			ApplicationID:  req.ApplicationID,
			RequestID:      req.RequestID,
			PrevStateRoot:  [32]byte{},
			NewStateRoot:   [32]byte{},
			ErrorCode:      uint8(apperrors.CodeFailedLoadingOrGettingModule.Category.Category),
			ErrorMsg:       "failed to load or get module",
			RefundAmount:   req.MaxFeeValue,
			ApplicationFee: common.NewBig(0),
		}, nil, nil
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)
	require.Nil(t, capturedWASM)
	require.Equal(t, 1, len(mockBCClient.GetCompletedRequests()))
	require.Equal(t, 1, len(mockBCClient.GetFailedRequests()))
}

func TestProcessDeployApp_ArtifactReadErrorSendsNilWASMToExecutor(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Simulate a non-not-found I/O read error by creating a directory where a blob file is expected.
	wasmSHA := strings.Repeat("c", 64)
	artifactBlobPath := filepath.Join(manager.config.ArtifactsPath, artifactBlobsFolder, wasmSHA+".wasm")
	require.NoError(t, os.MkdirAll(artifactBlobPath, 0o755))

	request := createRequestWithPayload(common.Deploy, ApplicationId, createDeployDescriptorPayload(t, wasmSHA))
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	var capturedWASM []byte
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		capturedWASM = wasmModule
		return &common.UpdatePayload{
			ApplicationID:  req.ApplicationID,
			RequestID:      req.RequestID,
			PrevStateRoot:  [32]byte{},
			NewStateRoot:   [32]byte{},
			ErrorCode:      uint8(apperrors.CodeFailedLoadingOrGettingModule.Category.Category),
			ErrorMsg:       "failed to load or get module",
			RefundAmount:   req.MaxFeeValue,
			ApplicationFee: common.NewBig(0),
		}, nil, nil
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)
	require.Nil(t, capturedWASM)
	require.Equal(t, 1, len(mockBCClient.GetCompletedRequests()))
	require.Equal(t, 1, len(mockBCClient.GetFailedRequests()))
}

func TestProcessDeployApp_AvailableArtifactForwardsDescriptorAndWASM(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	wasm := []byte("available-wasm")
	sum := sha256.Sum256(wasm)
	wasmSHA := hex.EncodeToString(sum[:])
	request := createRequestWithPayload(common.Deploy, ApplicationId, createDeployDescriptorPayload(t, wasmSHA))
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	writeArtifactBlob(t, manager, wasmSHA, wasm)

	var capturedReq *common.Request
	var capturedWASM []byte
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		capturedReq = req
		capturedWASM = append([]byte(nil), wasmModule...)
		stateRoot := sha256.Sum256([]byte("artifact-forward-success"))
		return &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, NewStateRoot: stateRoot}, &common.ApplicationState{ApplicationID: req.ApplicationID, StateRoot: stateRoot}, nil
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	require.Equal(t, request.Payload, capturedReq.Payload)
	require.Equal(t, wasm, capturedWASM)

	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	require.Equal(t, 0, len(mockBCClient.GetFailedRequests()), "expected 0 failed requests")
}

func TestProcessDeployApp_HashMismatchIsHandledByExecutor(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	wasm := []byte("real-content")
	actualHash := sha256.Sum256(wasm)
	descriptorHash := strings.Repeat("b", 64)
	writeArtifactBlob(t, manager, descriptorHash, wasm)

	request := createRequestWithPayload(common.Deploy, ApplicationId, createDeployDescriptorPayload(t, descriptorHash))
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	var capturedReq *common.Request
	var capturedWASM []byte
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendDeployApp", func(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
		capturedReq = req
		capturedWASM = append([]byte(nil), wasmModule...)
		return &common.UpdatePayload{
			ApplicationID:  req.ApplicationID,
			RequestID:      req.RequestID,
			PrevStateRoot:  [32]byte{},
			NewStateRoot:   [32]byte{},
			ErrorCode:      uint8(apperrors.CodeFailedLoadingOrGettingModule.Category.Category),
			ErrorMsg:       "failed to load or get module",
			RefundAmount:   req.MaxFeeValue,
			ApplicationFee: common.NewBig(0),
		}, nil, nil
	})

	err = manager.processDeployApp(context.Background(), request)
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	require.Equal(t, request.Payload, capturedReq.Payload)
	require.Equal(t, wasm, capturedWASM)
	require.NotEqual(t, descriptorHash, hex.EncodeToString(actualHash[:]))
	require.Equal(t, 1, len(mockBCClient.GetCompletedRequests()))
	require.Equal(t, 1, len(mockBCClient.GetFailedRequests()))
}

func TestProcessProcessRequestWithFailure(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy the application first
	deployRequest := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
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
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, *common.ApplicationState, *common.WASMData) error {
		t.Fatal("Store should not be called if the executor returned a failure payload")
		return nil
	})

	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendProcessRequest",
		func(_ context.Context, req *common.Request, _ *common.ApplicationState, _ []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
			return &common.UpdatePayload{ApplicationID: ApplicationId,
				RequestID: req.RequestID,
				ErrorCode: 1,
				ErrorMsg:  "error"}, nil, nil, nil
		})

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
	deployRequest := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployRequest)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	oldDbVersion, err := manager.dataLayer.LastVersionID(ApplicationId)
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
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, *common.ApplicationState, *common.WASMData) error {
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
	newDbVersion, err := manager.dataLayer.LastVersionID(ApplicationId)
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
	newDbVersion, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	// Same test but with an error payload for the SubmitStateUpdate
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendProcessRequest",
		func(_ context.Context, req *common.Request, _ *common.ApplicationState, _ []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
			return &common.UpdatePayload{ApplicationID: ApplicationId,
				RequestID: req.RequestID,
				ErrorCode: 1,
				ErrorMsg:  "error"}, nil, nil, nil
		})

	failure = manager.processProcessRequest(context.Background(), request)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	// Check that the local db has been reverted to the initial state
	newDbVersion, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendProcessRequest")
	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

}

// TestProcessDeanonymizationViaProcessRequest tests the deanonymization flow which is now
// handled through processProcessRequest with RequestType = Deanonymize
func TestProcessDeanonymizationViaProcessRequest(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy the application first
	deployRequest := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
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
	request1 := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
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

	db_version, err := manager.dataLayer.LastVersionID(ApplicationId)
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

	db_version, err = manager.dataLayer.LastVersionID(ApplicationId)
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

	db_version, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)

	chainState3, err := mockBCClient.GetApplicationState(context.Background(), ApplicationId)
	require.NoError(t, err)
	require.True(t, bytes.Equal(chainState3.StateRoot[:], db_version), "State root in DB should be equal to state root on chain")

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
	// Instead of sleeping, we will simulate the time.Sleep by manipulating endReorgTime
	manager.endReorgTime = manager.endReorgTime.Add(-2 * time.Second) // go back in time by 2 seconds

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pendingRequests), "expected 0 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	db_version, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)

	chainState, err := mockBCClient.GetApplicationState(context.Background(), ApplicationId)
	require.NoError(t, err)
	require.True(t, bytes.Equal(chainState.StateRoot[:], db_version), "State root in DB should be equal to state root on chain")

}

func TestProcessRequestFromChainWithErrors(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Setup the application
	request := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
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

	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("LastVersionID", func(common.ApplicationIdType) ([]byte, error) {
		return nil, fmt.Errorf("LastVersionID error")
	})

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err, "processRequestFromChain should not return an error if LastVersionID fails")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("LastVersionID")

	//**********************
	// Check that if ListVersions returns an error, processRequestFromChain doesn't execute the request and doesn't return an error
	//**********************
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("ListVersions", func(common.ApplicationIdType) ([][]byte, error) {
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
	if config.ArtifactsPath == "" {
		config.ArtifactsPath = t.TempDir()
	}

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
	deployRequest := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
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

// TestGetAndSetLogLevel_NonZeroNetworkLogger verifies that SetLogLevel and GetLogLevel
// return errors when the manager uses a non-ZeroNetworkLogger (e.g., zerolog).
func TestGetAndSetLogLevel_NonZeroNetworkLogger(t *testing.T) {
	_, manager := setupTest(t)
	ctx := context.Background()

	// With target "manager", GetLogLevel should fail directly.
	getMsg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Target: "manager"}
	result, err := manager.ExecuteCommand(ctx, getMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Empty(t, result)

	// With target "manager", SetLogLevel should fail directly.
	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "error"})
	require.NoError(t, err)
	setMsg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "manager", Data: setData}
	result, err = manager.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Nil(t, result)
}

func TestExecuteCommand_SetLogLevel_NilData(t *testing.T) {
	_, mgr := setupTest(t)

	// No Data at all
	msg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "manager"}
	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing request data for set_log_level")
	require.Nil(t, result)

	// Explicit null Data
	msg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "manager", Data: json.RawMessage("null")}
	result, err = mgr.ExecuteCommand(context.Background(), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing request data for set_log_level")
	require.Nil(t, result)
}

// TestSetLogLevel_TargetValidation verifies that SetLogLevel and GetLogLevel
// reject invalid targets on the manager.
func TestSetLogLevel_TargetValidation(t *testing.T) {
	_, manager := setupTest(t)
	ctx := context.Background()

	// SetLogLevel with target "executor" should be forwarded to executor only.
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdSetLogLevel, cmdType)
			resp, _ := json.Marshal(admin.SetLogLevelResponse{Level: "debug"})
			return resp, nil
		},
	)
	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	setMsg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "executor", Data: setData}
	result, err := manager.ExecuteCommand(ctx, setMsg)
	require.NoError(t, err)
	execResp, ok := result.(*admin.SetLogLevelResponse)
	require.True(t, ok)
	require.Equal(t, "debug", execResp.Level)

	// SetLogLevel with unknown target should be rejected
	setData, err = json.Marshal(admin.SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "unknown", Data: setData}
	result, err = manager.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown target 'unknown'")
	require.Nil(t, result)

	// SetLogLevel with target "manager" should be accepted (fails on logger type, not target)
	setData, err = json.Marshal(admin.SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "manager", Data: setData}
	result, err = manager.ExecuteCommand(ctx, setMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supported with the ZeroNetworkLogger")
	require.Nil(t, result)

	// GetLogLevel with target "executor" should be forwarded to executor only.
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdGetLogLevel, cmdType)
			resp, _ := json.Marshal("debug")
			return resp, nil
		},
	)
	getMsg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Target: "executor"}
	result, err = manager.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	require.Equal(t, "debug", result)

	// GetLogLevel with empty target defaults to "all".
	// Manager fails (non-ZeroNetworkLogger), executor mock still active → partial success.
	getMsg = admin.AdminMessage{Type: admin.GetLogLevelRequestMessage}
	result, err = manager.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	aggGetResp, ok := result.(admin.AggregatedGetLogLevelResponse)
	require.True(t, ok)
	require.Contains(t, aggGetResp.ManagerError, "only supported with the ZeroNetworkLogger")
	require.Equal(t, "debug", aggGetResp.Executor)

	// SetLogLevel with target "all": manager fails (non-ZeroNetworkLogger),
	// executor mock returns success → partial success.
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			resp, _ := json.Marshal(admin.SetLogLevelResponse{Level: "debug"})
			return resp, nil
		},
	)
	setData, err = json.Marshal(admin.SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	setMsg = admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "all", Data: setData}
	result, err = manager.ExecuteCommand(ctx, setMsg)
	require.NoError(t, err)
	aggSetResp, ok := result.(admin.AggregatedSetLogLevelResponse)
	require.True(t, ok)
	require.Contains(t, aggSetResp.ManagerError, "only supported with the ZeroNetworkLogger")
	require.Equal(t, "debug", aggSetResp.Executor)

	// GetLogLevel with target "all": manager fails, executor mock still active → partial success.
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			resp, _ := json.Marshal("debug")
			return resp, nil
		},
	)
	getMsg = admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Target: "all"}
	result, err = manager.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	aggGetResp, ok = result.(admin.AggregatedGetLogLevelResponse)
	require.True(t, ok)
	require.Contains(t, aggGetResp.ManagerError, "only supported with the ZeroNetworkLogger")
	require.Equal(t, "debug", aggGetResp.Executor)
}

// TestGetAndSetLogLevel_WithZeroNetworkLogger verifies the positive SetLogLevel/GetLogLevel
// round-trip when the manager uses a real ZeroNetworkLogger.
func TestGetAndSetLogLevel_WithZeroNetworkLogger(t *testing.T) {
	// Start a dummy TCP listener to act as the log sink.
	logSink, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer logSink.Close()
	go func() {
		for {
			conn, err := logSink.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				buf := make([]byte, 4096)
				for {
					if _, err := c.Read(buf); err != nil {
						c.Close()
						return
					}
				}
			}(conn)
		}
	}()

	logSinkAddr := logSink.Addr().(*net.TCPAddr)
	znl := logger.NewZeroNetworkLogger(&logger.Config{
		RemoteLogNetwork: "tcp",
		RemoteLogParams: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(logSinkAddr.Port),
		},
		NetworkLevel: "info",
	})
	defer znl.Close()

	// Create a manager with the ZeroNetworkLogger.
	_, mgr := setupTest(t)
	mgr.log = znl
	ctx := context.Background()

	// GetLogLevel with explicit target "manager" should return "info" (the initial level).
	getMsg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Target: "manager"}
	result, err := mgr.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	require.Equal(t, "info", result)

	// SetLogLevel to "debug" with explicit target "manager".
	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	setMsg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "manager", Data: setData}
	result, err = mgr.ExecuteCommand(ctx, setMsg)
	require.NoError(t, err)
	resp, ok := result.(admin.SetLogLevelResponse)
	require.True(t, ok)
	require.Equal(t, "debug", resp.Level)

	// GetLogLevel should now return "debug".
	result, err = mgr.ExecuteCommand(ctx, getMsg)
	require.NoError(t, err)
	require.Equal(t, "debug", result)

	// The underlying logger should also reflect the change.
	require.Equal(t, "debug", znl.GetLevel())
}

// --- Proxy tests (ForwardAdminCommand via communication channel) ---

func TestExecuteCommand_SetLogLevel_TargetAll(t *testing.T) {
	// Start a dummy TCP listener to act as the log sink.
	logSink, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer logSink.Close()
	go func() {
		for {
			conn, err := logSink.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				buf := make([]byte, 4096)
				for {
					if _, err := c.Read(buf); err != nil {
						c.Close()
						return
					}
				}
			}(conn)
		}
	}()

	logSinkAddr := logSink.Addr().(*net.TCPAddr)
	znl := logger.NewZeroNetworkLogger(&logger.Config{
		RemoteLogNetwork: "tcp",
		RemoteLogParams: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(logSinkAddr.Port),
		},
		NetworkLevel: "info",
	})
	defer znl.Close()

	_, mgr := setupTest(t)
	mgr.log = znl

	// Mock ForwardAdminCommand on the executor client
	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdSetLogLevel, cmdType)
			resp, _ := json.Marshal(admin.SetLogLevelResponse{Level: "debug"})
			return resp, nil
		},
	)

	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	msg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "all", Data: setData}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err)

	aggResp, ok := result.(admin.AggregatedSetLogLevelResponse)
	require.True(t, ok)
	require.Equal(t, "debug", aggResp.Manager)
	require.Equal(t, "debug", aggResp.Executor)

	// Verify the manager's logger was actually updated
	require.Equal(t, "debug", znl.GetLevel())
}

func TestExecuteCommand_GetLogLevel_TargetAll(t *testing.T) {
	logSink, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer logSink.Close()
	go func() {
		for {
			conn, err := logSink.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				buf := make([]byte, 4096)
				for {
					if _, err := c.Read(buf); err != nil {
						c.Close()
						return
					}
				}
			}(conn)
		}
	}()

	logSinkAddr := logSink.Addr().(*net.TCPAddr)
	znl := logger.NewZeroNetworkLogger(&logger.Config{
		RemoteLogNetwork: "tcp",
		RemoteLogParams: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(logSinkAddr.Port),
		},
		NetworkLevel: "info",
	})
	defer znl.Close()

	_, mgr := setupTest(t)
	mgr.log = znl

	// Mock ForwardAdminCommand on the executor client
	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdGetLogLevel, cmdType)
			resp, _ := json.Marshal("debug")
			return resp, nil
		},
	)

	msg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Target: "all"}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err)

	aggResp, ok := result.(admin.AggregatedGetLogLevelResponse)
	require.True(t, ok)
	require.Equal(t, "info", aggResp.Manager)
	require.Equal(t, "debug", aggResp.Executor)
}

func TestExecuteCommand_SetLogLevel_TargetAll_ExecutorFails(t *testing.T) {
	logSink, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer logSink.Close()
	go func() {
		for {
			conn, err := logSink.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				buf := make([]byte, 4096)
				for {
					if _, err := c.Read(buf); err != nil {
						c.Close()
						return
					}
				}
			}(conn)
		}
	}()

	logSinkAddr := logSink.Addr().(*net.TCPAddr)
	znl := logger.NewZeroNetworkLogger(&logger.Config{
		RemoteLogNetwork: "tcp",
		RemoteLogParams: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(logSinkAddr.Port),
		},
		NetworkLevel: "info",
	})
	defer znl.Close()

	_, mgr := setupTest(t)
	mgr.log = znl

	// Mock ForwardAdminCommand to fail (simulating executor unreachable)
	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			return nil, fmt.Errorf("connection refused")
		},
	)

	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "warn"})
	require.NoError(t, err)
	msg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "all", Data: setData}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err) // partial success is not an error

	aggResp, ok := result.(admin.AggregatedSetLogLevelResponse)
	require.True(t, ok)
	require.Equal(t, "warn", aggResp.Manager)
	require.Empty(t, aggResp.ManagerError)
	require.Contains(t, aggResp.ExecutorError, "connection refused")

	// Verify the manager's log level WAS changed (even though executor failed)
	require.Equal(t, "warn", znl.GetLevel())
}

func TestExecuteCommand_KeyAttestation_ForwardSuccess(t *testing.T) {
	_, mgr := setupTest(t)

	expectedAttestation := []byte(`"base64-encoded-attestation-doc"`)
	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdKeyAttestation, cmdType)
			return expectedAttestation, nil
		},
	)

	msg := admin.AdminMessage{Type: admin.KeyAttestationRequestMessage}
	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err)
	require.Equal(t, json.RawMessage(expectedAttestation), result)
}

func TestExecuteCommand_KeyAttestation_ForwardError(t *testing.T) {
	_, mgr := setupTest(t)

	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			return nil, fmt.Errorf("executor unreachable")
		},
	)

	msg := admin.AdminMessage{Type: admin.KeyAttestationRequestMessage}
	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "executor unreachable")
	require.Nil(t, result)
}

func TestExecuteCommand_KeyAttestation_RejectsTargetManager(t *testing.T) {
	_, mgr := setupTest(t)

	msg := admin.AdminMessage{Type: admin.KeyAttestationRequestMessage, Target: "manager"}
	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "key_attestation is only supported on the executor")
	require.Nil(t, result)
}

func TestExecuteCommand_GetVersion_TargetManager(t *testing.T) {
	_, mgr := setupTest(t)

	msg := admin.AdminMessage{Type: admin.GetVersionRequestMessage, Target: "manager"}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err)
	require.Equal(t, version.Version, result)
}

func TestExecuteCommand_GetVersion_TargetExecutor(t *testing.T) {
	_, mgr := setupTest(t)

	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdGetVersion, cmdType)
			resp, _ := json.Marshal("v1.0.0-executor")
			return resp, nil
		},
	)

	msg := admin.AdminMessage{Type: admin.GetVersionRequestMessage, Target: "executor"}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err)
	require.Equal(t, "v1.0.0-executor", result)
}

func TestExecuteCommand_GetVersion_TargetAll(t *testing.T) {
	_, mgr := setupTest(t)

	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdGetVersion, cmdType)
			resp, _ := json.Marshal("v2.0.0-executor")
			return resp, nil
		},
	)

	msg := admin.AdminMessage{Type: admin.GetVersionRequestMessage, Target: "all"}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err)

	aggResp, ok := result.(admin.AggregatedGetVersionResponse)
	require.True(t, ok)
	require.Equal(t, version.Version, aggResp.Manager)
	require.Equal(t, "v2.0.0-executor", aggResp.Executor)
}

func TestExecuteCommand_GetVersion_DefaultTarget(t *testing.T) {
	_, mgr := setupTest(t)

	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdGetVersion, cmdType)
			resp, _ := json.Marshal("v3.0.0-executor")
			return resp, nil
		},
	)

	// No data (null) — should default to target="all"
	msg := admin.AdminMessage{Type: admin.GetVersionRequestMessage}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err)

	aggResp, ok := result.(admin.AggregatedGetVersionResponse)
	require.True(t, ok)
	require.Equal(t, version.Version, aggResp.Manager)
	require.Equal(t, "v3.0.0-executor", aggResp.Executor)
}

func TestExecuteCommand_GetVersion_TargetAll_ExecutorFails(t *testing.T) {
	_, mgr := setupTest(t)

	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			return nil, fmt.Errorf("connection refused")
		},
	)

	msg := admin.AdminMessage{Type: admin.GetVersionRequestMessage, Target: "all"}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err) // partial success is not an error

	aggResp, ok := result.(admin.AggregatedGetVersionResponse)
	require.True(t, ok)
	require.Equal(t, version.Version, aggResp.Manager)
	require.Empty(t, aggResp.ManagerError)
	require.Empty(t, aggResp.Executor)
	require.Contains(t, aggResp.ExecutorError, "connection refused")
}

func TestExecuteCommand_SetLogLevel_TargetAll_BothFail(t *testing.T) {
	// Manager uses the default test logger (not ZeroNetworkLogger), so local
	// SetLogLevel will fail. The executor mock also returns an error.
	_, mgr := setupTest(t)

	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			return nil, fmt.Errorf("executor unreachable")
		},
	)

	setData, err := json.Marshal(admin.SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	msg := admin.AdminMessage{Type: admin.SetLogLevelRequestMessage, Target: "all", Data: setData}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err) // target="all" always returns aggregated response, never an error

	aggResp, ok := result.(admin.AggregatedSetLogLevelResponse)
	require.True(t, ok)
	require.Empty(t, aggResp.Manager)
	require.Contains(t, aggResp.ManagerError, "only supported with the ZeroNetworkLogger")
	require.Empty(t, aggResp.Executor)
	require.Contains(t, aggResp.ExecutorError, "executor unreachable")
}

func TestExecuteCommand_GetLogLevel_TargetAll_BothFail(t *testing.T) {
	// Manager uses the default test logger (not ZeroNetworkLogger), so local
	// GetLogLevel will fail. The executor mock also returns an error.
	_, mgr := setupTest(t)

	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			return nil, fmt.Errorf("executor unreachable")
		},
	)

	msg := admin.AdminMessage{Type: admin.GetLogLevelRequestMessage, Target: "all"}

	result, err := mgr.ExecuteCommand(context.Background(), msg)
	require.NoError(t, err) // target="all" always returns aggregated response, never an error

	aggResp, ok := result.(admin.AggregatedGetLogLevelResponse)
	require.True(t, ok)
	require.Empty(t, aggResp.Manager)
	require.Contains(t, aggResp.ManagerError, "only supported with the ZeroNetworkLogger")
	require.Empty(t, aggResp.Executor)
	require.Contains(t, aggResp.ExecutorError, "executor unreachable")
}

func TestExecuteCommand_UnsupportedCommand(t *testing.T) {
	_, mgr := setupTest(t)

	unknownMsg := admin.AdminMessage{Type: "unknown_type"}
	result, err := mgr.ExecuteCommand(context.Background(), unknownMsg)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "unsupported command type")
}

// --- Multi-App Tests ---

var ApplicationId2 = common.NewApplicationId(2)

func TestMultiAppDeployAndProcess(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy App 1
	deployReq1 := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployReq1)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	completed := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completed), "App 1 deploy should complete")

	// Deploy App 2
	deployReq2 := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	err = mockBCClient.SendRequestToChain(context.Background(), deployReq2)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	completed = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completed), "App 2 deploy should complete")

	// Verify independent state roots
	stateRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	stateRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.False(t, bytes.Equal(stateRoot1, stateRoot2), "App 1 and App 2 should have different state roots")

	// Process a request for App 1
	processReq1 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq1)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	// Verify App 1 state root changed but App 2 unchanged
	newStateRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.False(t, bytes.Equal(stateRoot1, newStateRoot1), "App 1 state root should change after processing")

	unchangedStateRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRoot2, unchangedStateRoot2), "App 2 state root should be unchanged")

	// Process a request for App 2
	processReq2 := createRequest(common.Process, ApplicationId2)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq2)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	completed = mockBCClient.GetCompletedRequests()
	require.Equal(t, 4, len(completed), "All 4 requests should be completed")

	// Verify App 2 state root changed
	newStateRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.False(t, bytes.Equal(stateRoot2, newStateRoot2), "App 2 state root should change after processing")
}

func TestMultiAppDeployNewApp(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy and process App 1
	deployReq1 := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployReq1)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	processReq1 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq1)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	stateRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)

	// Deploy App 2 — no prior state
	deployReq2 := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	err = mockBCClient.SendRequestToChain(context.Background(), deployReq2)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	completed := mockBCClient.GetCompletedRequests()
	require.Equal(t, 3, len(completed), "All 3 requests should be completed")

	// App 1 state unchanged
	unchangedRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRoot1, unchangedRoot1), "App 1 state should be unchanged after App 2 deploy")

	// App 2 now has state
	stateRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.NotEqual(t, make([]byte, 32), stateRoot2, "App 2 should have a non-zero state root after deploy")
}

func TestMultiAppReorgIsolation(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy App 1 and App 2
	deployReq1 := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), deployReq1)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	deployReq2 := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	err = mockBCClient.SendRequestToChain(context.Background(), deployReq2)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	// Process a request for App 1 (creates a second version)
	processReq1 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq1)
	require.NoError(t, err)
	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	stateRoot1AfterProcess, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	stateRoot2AfterDeploy, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)

	// Simulate reorg for App 1 — chain returns old stateRoot (zero, as if rolled back to before deploy)
	mockedReorg := func(context.Context) (*common.Request, [32]byte, error) {
		return processReq1, [32]byte{}, nil
	}
	mockBCClient.AddMockedFunc("GetNextPendingRequest", mockedReorg)
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(context.Context, *common.UpdatePayload) error {
		t.Fatal("SubmitStateUpdate should not be called during reorg")
		return nil
	})

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	// Reorg timer should be set (chain-level, triggered by App 1's mismatch)
	require.False(t, manager.endReorgTime.IsZero(), "reorg timer should be set")

	// App 2 state should be untouched
	unchangedRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRoot2AfterDeploy, unchangedRoot2), "App 2 state should be unchanged during App 1 reorg")

	// Resolve reorg — remove mock, next poll returns App 1's correct state
	mockBCClient.RemoveMockedFunc("GetNextPendingRequest")
	mockBCClient.RemoveMockedFunc("SubmitStateUpdate")

	// Restore App 1's request to pending and set correct state root
	err = mockBCClient.SendRequestToChain(context.Background(), createRequest(common.Process, ApplicationId))
	require.NoError(t, err)

	err = manager.processRequestFromChain(context.Background())
	require.NoError(t, err)

	// Reorg timer should be cleared (state roots matched)
	require.True(t, manager.endReorgTime.IsZero(), "reorg timer should be cleared after resolution")

	// App 1 should have processed the new request
	newRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.False(t, bytes.Equal(stateRoot1AfterProcess, newRoot1), "App 1 should have new state after processing")
}

// TestSingleReorgTimerPreventsStaleRollback verifies that the chain-level reorg
// timer prevents the stale-timer bug that existed with per-app timers.
//
// With per-app timers, an app that received no pending requests between two
// reorgs would keep a stale expired timer, causing the second reorg to skip the
// wait and trigger an immediate (dangerous) rollback.
//
// The single chain-level timer fixes this: when ANY app resolves (state roots
// match), the shared timer is cleared — so no stale timer can survive.
//
// Timeline:
//  1. Setup: Deploy App A and App B, process one request each so each app has
//     ≥2 DB versions (needed for checkIfReorg to find a matching old root).
//  2. Reorg #1: both apps see a state-root mismatch → shared timer starts.
//  3. Reorg #1 resolves: App A gets a new pending request → state roots match
//     → the shared timer is cleared. App B gets NO pending request, but the
//     shared timer is already gone — no stale state can accumulate.
//  4. Reorg #2: App B sees a new state-root mismatch. The timer is zero, so a
//     fresh timeout is started. No immediate rollback.
func TestSingleReorgTimerPreventsStaleRollback(t *testing.T) {
	mockBCClient, manager := setupTest(t)
	ctx := context.Background()

	// Use a long reorg timeout so we can clearly tell if the wait was skipped.
	manager.config.ReorgTimeout = 300

	// ── Step 1: Deploy and process one request for App A and App B ──

	deployA := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, deployA))
	require.NoError(t, manager.processRequestFromChain(ctx))

	processA := createRequest(common.Process, ApplicationId)
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, processA))
	require.NoError(t, manager.processRequestFromChain(ctx))

	deployB := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, deployB))
	require.NoError(t, manager.processRequestFromChain(ctx))

	processBReq := createRequest(common.Process, ApplicationId2)
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, processBReq))
	require.NoError(t, manager.processRequestFromChain(ctx))

	// Grab the deploy-time (older) state roots for both apps.
	// ListVersions returns LIFO: [latest, ..., oldest].
	// The deploy root is the second entry (index 1).
	appAVersions, err := manager.dataLayer.ListVersions(ApplicationId)
	require.NoError(t, err)
	require.True(t, len(appAVersions) >= 2, "App A should have ≥2 versions")
	stateRootAAfterDeploy := appAVersions[1] // older version (deploy root)

	appBVersions, err := manager.dataLayer.ListVersions(ApplicationId2)
	require.NoError(t, err)
	require.True(t, len(appBVersions) >= 2, "App B should have ≥2 versions")
	stateRootBAfterDeploy := appBVersions[1] // older version (deploy root)

	// ── Step 2: Reorg #1 — both apps see a mismatch ──

	// First poll: App A sees reorg (chain returns deploy root, but local is at process root)
	var reorgStateRootA [32]byte
	copy(reorgStateRootA[:], stateRootAAfterDeploy)
	mockBCClient.AddMockedFunc("GetNextPendingRequest", func(context.Context) (*common.Request, [32]byte, error) {
		return processA, reorgStateRootA, nil
	})
	mockBCClient.AddMockedFunc("SubmitStateUpdate", func(context.Context, *common.UpdatePayload) error {
		t.Fatal("SubmitStateUpdate should not be called during reorg")
		return nil
	})

	require.NoError(t, manager.processRequestFromChain(ctx))
	require.False(t, manager.endReorgTime.IsZero(), "reorg timer should be set after App A mismatch")

	// Second poll: App B sees reorg (chain returns deploy root).
	// The shared timer is already running — it should NOT be reset.
	var reorgStateRootB [32]byte
	copy(reorgStateRootB[:], stateRootBAfterDeploy)
	timerBefore := manager.endReorgTime
	mockBCClient.AddMockedFunc("GetNextPendingRequest", func(context.Context) (*common.Request, [32]byte, error) {
		return processBReq, reorgStateRootB, nil
	})

	require.NoError(t, manager.processRequestFromChain(ctx))
	require.Equal(t, timerBefore, manager.endReorgTime, "shared timer should not change on second mismatch")

	// ── Step 3: Reorg #1 resolves — only App A gets a request ──

	mockBCClient.RemoveMockedFunc("GetNextPendingRequest")
	mockBCClient.RemoveMockedFunc("SubmitStateUpdate")

	// App A gets a new request; state roots will match → shared timer cleared.
	newReqA := createRequest(common.Process, ApplicationId)
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, newReqA))
	require.NoError(t, manager.processRequestFromChain(ctx))

	require.True(t, manager.endReorgTime.IsZero(),
		"shared timer should be cleared when App A's state roots match (reorg resolved)")

	// ── Step 4: Reorg #2 — App B sees a NEW mismatch ──

	// The chain has reorged again: it reports the deploy root for App B.
	processB2 := createRequest(common.Process, ApplicationId2)
	mockBCClient.AddMockedFunc("GetNextPendingRequest", func(context.Context) (*common.Request, [32]byte, error) {
		return processB2, reorgStateRootB, nil
	})

	// Track whether Rollback is called for App B during this poll.
	rollbackCalled := false
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Rollback",
		func(appID common.ApplicationIdType, versionID []byte) error {
			if appID == ApplicationId2 {
				rollbackCalled = true
			}
			return nil
		})

	require.NoError(t, manager.processRequestFromChain(ctx))

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Rollback")

	// The shared timer was zero before this poll, so a fresh timeout must be
	// started — no rollback should occur.
	require.False(t, rollbackCalled,
		"Rollback should not be called — the new reorg should get a fresh timeout")
	require.False(t, manager.endReorgTime.IsZero(),
		"reorg timer should be set with a fresh deadline for the new reorg")
	require.True(t, manager.endReorgTime.After(time.Now()),
		"reorg timer should be in the future (fresh timeout)")
}

// TestHandleKeysetRecovery_SetsRunningPcr0 verifies the manager records the
// PCR0 (running image) the executor reports during the handshake, for both the
// recovery-result and set-recovery paths.
func TestHandleKeysetRecovery_SetsRunningPcr0(t *testing.T) {
	ctx := context.Background()
	config := Config{
		ReorgTimeout:        60,
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}

	// Recovery-result path (reconnection / restore).
	_, mgr := setupTestWithConfig(t, ctx, config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, nil, false)
	require.Empty(t, mgr.RunningPcr0(), "PCR0 should be empty before handshake")

	err := mgr.HandleKeysetRecoveryResult(ctx, nil, "comm-pub-key", "0xsigner", "deadbeef")
	require.NoError(t, err)
	require.Equal(t, "deadbeef", mgr.RunningPcr0())

	// Set-recovery path (first connection / generate).
	_, mgr2 := setupTestWithConfig(t, ctx, config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, nil, false)
	err = mgr2.HandleSetKeysetRecoveryRequest(ctx, &common.EnclaveKeySetRecovery{
		RecoveryType: common.RecoveryTypeKMS,
	}, "comm-pub-key", "0xsigner", "cafebabe")
	require.NoError(t, err)
	require.Equal(t, "cafebabe", mgr2.RunningPcr0())
}

// TestHandleSetKeysetRecoveryRequest_StoreExistsFailsHandshake verifies that
// when the data layer refuses to overwrite an existing keyset recovery blob
// (ErrRecoveryDataExists — the R2 key-continuity guard), the error is surfaced
// by HandleSetKeysetRecoveryRequest and the handshake is failed loudly with the
// same typed error rather than silently completing.
func TestHandleSetKeysetRecoveryRequest_StoreExistsFailsHandshake(t *testing.T) {
	ctx := context.Background()
	config := Config{
		ReorgTimeout:        60,
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	_, mgr := setupTestWithConfig(t, ctx, config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, nil, false)

	storeErr := storageErrors.ErrRecoveryDataExists(
		"refusing to overwrite existing enclave keyset recovery data")
	mgr.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("StoreEnclaveKeySetRecovery",
		func(context.Context, *common.EnclaveKeySetRecovery) error {
			return storeErr
		})

	err := mgr.HandleSetKeysetRecoveryRequest(ctx, &common.EnclaveKeySetRecovery{
		RecoveryType: common.RecoveryTypeKMS,
	}, "comm-pub-key", "0xsigner", "deadbeef")

	// The typed store error is surfaced to the caller.
	require.ErrorIs(t, err, storeErr)
	require.True(t, storageErrors.IsRecoveryDataExists(err),
		"the returned error must be the RecoveryDataExists guard error")

	// The handshake was completed loudly with the same error — not silently succeeded.
	require.ErrorIs(t, mgr.executorHandShake.err, storeErr,
		"the handshake must fail with the store error")
	require.ErrorIs(t, mgr.waitForExecutorHandshake(), storeErr,
		"waitForExecutorHandshake must surface the store error")
}

// TestHandleGetKeysetRecoveryRequest_IncompatibleProtocol verifies that an
// incompatible executor protocol version fails the handshake with a typed
// error, before any keyset-recovery data is read or exchanged.
func TestHandleGetKeysetRecoveryRequest_IncompatibleProtocol(t *testing.T) {
	ctx := context.Background()
	config := Config{
		ReorgTimeout:        60,
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	_, mgr := setupTestWithConfig(t, ctx, config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, nil, false)

	incompatible := communication.WireProtocolVersion + 1
	recv, err := mgr.HandleGetKeysetRecoveryRequest(ctx, incompatible)

	require.Nil(t, recv, "no recovery data should be returned for an incompatible peer")
	var incompatErr *communication.IncompatibleProtocolError
	require.ErrorAs(t, err, &incompatErr)
	require.Equal(t, incompatible, incompatErr.Peer)

	// The handshake was completed with the typed error — a clean abort, not a timeout.
	require.ErrorAs(t, mgr.executorHandShake.err, &incompatErr)
}

// TestHandleGetKeysetRecoveryRequest_CompatibleProceeds verifies a compatible
// peer is not rejected: the handler proceeds to read recovery data (NotFound on
// a fresh mock) and does not fail the handshake.
func TestHandleGetKeysetRecoveryRequest_CompatibleProceeds(t *testing.T) {
	ctx := context.Background()
	config := Config{
		ReorgTimeout:        60,
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	_, mgr := setupTestWithConfig(t, ctx, config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, nil, false)

	recv, err := mgr.HandleGetKeysetRecoveryRequest(ctx, communication.WireProtocolVersion)

	require.Nil(t, recv)
	require.Error(t, err, "fresh mock has no recovery data (NotFound)")
	var incompatErr *communication.IncompatibleProtocolError
	require.False(t, errors.As(err, &incompatErr), "compatible peer must not yield an incompatibility error")
	require.NoError(t, mgr.executorHandShake.err, "handshake must not be failed for a compatible peer")
}

// TestForwardShutdown_ToleratesDisconnect verifies the manager does not treat a
// missing shutdown ack (executor already exiting / disconnected) as an error —
// the drain sequence must not turn this into a fatal error.
func TestForwardShutdown_ToleratesDisconnect(t *testing.T) {
	_, mgr := setupTest(t)
	ctx := context.Background()

	// Executor drops the connection instead of acking.
	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			require.Equal(t, admin.AdminCmdShutdown, cmdType)
			return nil, fmt.Errorf("connection closed")
		},
	)
	require.NoError(t, mgr.forwardShutdown(ctx), "disconnect on shutdown must be swallowed")

	// Executor acks cleanly.
	mgr.executorClient.(*MockExecutorClient).AddMockedFunc("ForwardAdminCommand",
		func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
			resp, _ := json.Marshal(admin.ShutdownResponse{Stopping: true})
			return resp, nil
		},
	)
	require.NoError(t, mgr.forwardShutdown(ctx))
}

// --- Task 5: swap observation / drain / reconnect / signer continuity ---

// newSwapTestManager builds a manager with an initialized handshake and returns
// it plus its mock blockchain client, for TEE-swap tests.
func newSwapTestManager(t *testing.T) (*blockchain.MockClient, *SecureProcessorManager) {
	config := Config{
		ReorgTimeout:        60,
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	return setupTestWithConfig(t, context.Background(), config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, nil, false)
}

// pcr0AndImage returns a hex-encoded PCR0 and its keccak256 (the activeImage).
func pcr0AndImage(seed byte) (string, [32]byte) {
	raw := make([]byte, 48)
	for i := range raw {
		raw[i] = seed + byte(i)
	}
	return hex.EncodeToString(raw), [32]byte(ethCrypto.Keccak256Hash(raw))
}

func TestMaintainExecutor_MatchDispatches(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	pcr0, image := pcr0AndImage(0x10)
	mgr.setExecutorIdentity(pcr0, "0xabc")
	bc.SetActiveImage(image)
	// teeSigner unset (zero) -> signer check deferred, not fatal.

	skip, err := mgr.maintainExecutorForActiveImage(ctx)
	require.NoError(t, err)
	require.False(t, skip, "matching image must not skip dispatch")
	require.False(t, mgr.isDraining())
}

func TestMaintainExecutor_MismatchDrains(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	pcr0, _ := pcr0AndImage(0x20)
	_, otherImage := pcr0AndImage(0x99) // different image on-chain
	mgr.setExecutorIdentity(pcr0, "0xabc")
	bc.SetActiveImage(otherImage)

	var shutdownSent, closed bool
	mec := mgr.executorClient.(*MockExecutorClient)
	mec.AddMockedFunc("ForwardAdminCommand", func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
		if cmdType == admin.AdminCmdShutdown {
			shutdownSent = true
		}
		return nil, fmt.Errorf("connection closed")
	})
	mec.AddMockedFunc("Close", func() error { closed = true; return nil })

	skip, err := mgr.maintainExecutorForActiveImage(ctx)
	require.NoError(t, err, "drain is not a fatal error")
	require.True(t, skip, "mismatch must skip dispatch")
	require.True(t, mgr.isDraining())
	require.True(t, shutdownSent, "shutdown must be signaled on drain")
	require.True(t, closed, "channel must be closed on drain")
}

func TestMaintainExecutor_DevMarkerProceeds(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	// Executor reported no PCR0 (dev/TCP). Any activeImage is ignored.
	mgr.setExecutorIdentity("", "")
	bc.SetActiveImage([32]byte{0x01})

	skip, err := mgr.maintainExecutorForActiveImage(ctx)
	require.NoError(t, err)
	require.False(t, skip)
	require.False(t, mgr.isDraining())
}

func TestVerifySignerContinuity(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	signer := ethCommon.HexToAddress("0x1111111111111111111111111111111111111111")
	mgr.setExecutorIdentity("aa", signer.Hex())

	// Bootstrap: teeSigner unset -> deferred, dispatch proceeds, not verified.
	skip, err := mgr.verifySignerContinuity(ctx)
	require.NoError(t, err)
	require.False(t, skip)
	mgr.identityMu.RLock()
	require.False(t, mgr.signerVerified)
	mgr.identityMu.RUnlock()

	// Matching non-zero teeSigner -> passes and is marked verified.
	bc.SetTeeSigner(signer)
	skip, err = mgr.verifySignerContinuity(ctx)
	require.NoError(t, err)
	require.False(t, skip)
	mgr.identityMu.RLock()
	require.True(t, mgr.signerVerified)
	mgr.identityMu.RUnlock()
}

func TestVerifySignerContinuity_MismatchFatal(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	mgr.setExecutorIdentity("aa", "0x1111111111111111111111111111111111111111")
	bc.SetTeeSigner(ethCommon.HexToAddress("0x2222222222222222222222222222222222222222"))

	skip, err := mgr.verifySignerContinuity(ctx)
	require.False(t, skip)
	var sce *SignerContinuityError
	require.ErrorAs(t, err, &sce)
}

// TestVerifySignerContinuity_TransientErrorSkips verifies a transient teeSigner
// read error skips dispatch (retry next tick) instead of dispatching unverified,
// consistent with the activeImage read.
func TestVerifySignerContinuity_TransientErrorSkips(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	mgr.setExecutorIdentity("aa", "0x1111111111111111111111111111111111111111")
	bc.AddMockedFunc("GetTeeSigner", func(context.Context) (ethCommon.Address, error) {
		return ethCommon.Address{}, fmt.Errorf("rpc unavailable")
	})

	skip, err := mgr.verifySignerContinuity(ctx)
	require.NoError(t, err, "transient read error is not fatal")
	require.True(t, skip, "transient read error must skip dispatch, not proceed unverified")
	mgr.identityMu.RLock()
	require.False(t, mgr.signerVerified)
	mgr.identityMu.RUnlock()
}

// TestReconnect_ConcurrentStaleHandshakeCompletion exercises the handShakeMu
// guard (Fix 1): a detached handler goroutine from a drained connection may call
// completeExecutorHandshake concurrently with reconnectExecutor swapping in a
// fresh handshake tracker. Under -race this must not report a data race on the
// executorHandShake pointer.
func TestReconnect_ConcurrentStaleHandshakeCompletion(t *testing.T) {
	_, mgr := newSwapTestManager(t)
	ctx := context.Background()

	mec := mgr.executorClient.(*MockExecutorClient)
	mec.AddMockedFunc("Connect", func(ctx context.Context, tag string) error {
		mgr.completeExecutorHandshake(nil) // the new connection's handshake
		return nil
	})

	// A stale detached handler keeps trying to complete the handshake while the
	// reconnect swaps the tracker pointer underneath it.
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				mgr.completeExecutorHandshake(nil)
			}
		}
	}()

	mgr.identityMu.Lock()
	mgr.draining = true
	mgr.identityMu.Unlock()

	require.NoError(t, mgr.reconnectExecutor(ctx))
	close(stop)
	wg.Wait()
}

// TestReconnect_HandshakeFailureClosesAndRetries exercises Fix #2: when the
// re-handshake fails (Connect succeeds but the handshake never completes and
// times out), reconnectExecutor must close the half-open client so the next
// tick can re-dial — instead of wedging forever on "already connected".
func TestReconnect_HandshakeFailureClosesAndRetries(t *testing.T) {
	config := Config{
		ReorgTimeout:        60,
		HandshakeTimeout:    0, // time.After(0) => the handshake wait times out at once
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	_, mgr := setupTestWithConfig(t, context.Background(), config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, nil, false)
	ctx := context.Background()

	// Model the real client's connect/close contract: a second Connect without an
	// intervening Close returns "already connected" (the wedge Fix #2 prevents).
	var connectCalls, closeCalls int
	connected := false
	mec := mgr.executorClient.(*MockExecutorClient)
	mec.AddMockedFunc("Connect", func(ctx context.Context, tag string) error {
		connectCalls++
		if connected {
			return fmt.Errorf("already connected")
		}
		connected = true
		return nil // succeeds, but never completes the handshake -> wait times out
	})
	mec.AddMockedFunc("Close", func() error {
		closeCalls++
		connected = false
		return nil
	})

	mgr.identityMu.Lock()
	mgr.draining = true
	mgr.identityMu.Unlock()

	// Tick 1: dial succeeds, handshake times out -> reconnect fails and closes.
	err := mgr.reconnectExecutor(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "handshake failed", "Connect should succeed; the handshake times out")
	require.Equal(t, 1, connectCalls)
	require.Equal(t, 1, closeCalls, "half-open client must be closed on handshake failure")
	require.True(t, mgr.isDraining(), "still draining after a failed reconnect")

	// Tick 2: must re-dial (Close reset the connected state) and again fail on the
	// handshake — NOT wedge on "already connected".
	err = mgr.reconnectExecutor(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "handshake failed", "next tick must re-dial, not wedge on 'already connected'")
	require.Equal(t, 2, connectCalls, "next tick must attempt Connect again")
	require.Equal(t, 2, closeCalls)
}

// TestMaintainExecutor_ReconnectsOnConnectionLoss exercises the plain
// connection-loss path: the executor crashed/restarted with the SAME image, so
// no swap occurred and draining is never set. maintainExecutorForActiveImage
// must still detect the dropped channel, re-dial, and re-handshake before
// dispatching — rather than retrying into a dead connection forever.
func TestMaintainExecutor_ReconnectsOnConnectionLoss(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	pcr0, image := pcr0AndImage(0x40)
	mgr.setExecutorIdentity(pcr0, "")
	bc.SetActiveImage(image)

	mec := mgr.executorClient.(*MockExecutorClient)
	// The channel has dropped: the reader-loop teardown flipped connected off.
	connected := false
	mec.AddMockedFunc("IsConnected", func() bool { return connected })
	// Reconnect models the relaunched enclave (same image) re-handshaking.
	var connectCalls int
	mec.AddMockedFunc("Connect", func(ctx context.Context, tag string) error {
		connectCalls++
		connected = true
		mgr.setExecutorIdentity(pcr0, "")
		mgr.completeExecutorHandshake(nil)
		return nil
	})

	// Not draining -> the only reconnect trigger is the dropped connection.
	require.False(t, mgr.isDraining())

	skip, err := mgr.maintainExecutorForActiveImage(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, connectCalls, "a dropped connection must trigger a re-dial")
	require.False(t, skip, "after reconnecting to the matching image, dispatch resumes")
	require.False(t, mgr.isDraining())
}

// TestStop_DoesNotDeadlockWithInFlightTick guards the Stop lock-ordering fix.
// processRequestFromChain does its slow work (blockchain reconnect,
// maintainExecutorForActiveImage → GetActiveImage / reconnect handshake) BEFORE
// taking m.mu, then acquires m.mu only at dispatch. Stop must flip isRunning and
// close stopChan under m.mu, then RELEASE m.mu before wg.Wait(). Otherwise a tick
// already parked in that pre-lock section blocks forever at m.mu.Lock() while Stop
// blocks forever in wg.Wait() holding m.mu — a lock-ordering deadlock.
func TestStop_DoesNotDeadlockWithInFlightTick(t *testing.T) {
	config := Config{
		ReorgTimeout:        60,
		LogServerTCPAddress: common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
	}
	bc, mgr := setupTestWithConfig(t, context.Background(), config, true,
		&ExecutorHandShake{isComplete: make(chan struct{})}, make(chan struct{}), false)
	ctx := context.Background()

	// Running a known image that matches on-chain, so maintain does NOT skip and
	// the tick proceeds all the way to m.mu.Lock() at dispatch (the skip path would
	// early-return before the lock and never reproduce the deadlock).
	pcr0, image := pcr0AndImage(0x50)
	mgr.setExecutorIdentity(pcr0, "")
	bc.SetActiveImage(image)

	// Block the pre-lock section deterministically inside GetActiveImage. The mock
	// holds only the MockClient's own lock here, never the manager's m.mu.
	entered := make(chan struct{})
	release := make(chan struct{})
	bc.AddMockedFunc("GetActiveImage", func(context.Context) ([32]byte, error) {
		close(entered)
		<-release
		return image, nil
	})

	// Register a poll-like goroutine in m.wg, exactly as pollBlockchain does.
	mgr.wg.Add(1)
	go func() {
		defer mgr.wg.Done()
		_ = mgr.processRequestFromChain(ctx)
	}()

	<-entered // the tick is now parked in the pre-lock section

	stopped := make(chan error, 1)
	go func() { stopped <- mgr.Stop() }()

	// Give Stop time to take m.mu, flip isRunning, close stopChan, and reach
	// wg.Wait(); then unblock the tick so it advances to m.mu.Lock() at dispatch.
	time.Sleep(100 * time.Millisecond)
	close(release)

	select {
	case err := <-stopped:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Stop deadlocked: held m.mu across wg.Wait() while an in-flight tick needed it")
	}

	// The tick must have exited via the !isRunning guard, not dispatched.
	mgr.mu.Lock()
	running := mgr.isRunning
	mgr.mu.Unlock()
	require.False(t, running, "Stop must leave the manager not running")
}

// TestSwapDrainReconnect_Converges walks the full swap cycle: an activeImage
// change drains the old enclave, a reconnect re-handshakes the relaunched
// enclave whose PCR0 matches, and dispatch resumes. Also covers rollback (the
// activeImage is simply re-pointed) and level-based convergence.
func TestSwapDrainReconnect_Converges(t *testing.T) {
	bc, mgr := newSwapTestManager(t)
	ctx := context.Background()

	oldPcr0, _ := pcr0AndImage(0x30)
	newPcr0, newImage := pcr0AndImage(0x31)

	// Running the old image; chain now points activeImage at the new image.
	mgr.setExecutorIdentity(oldPcr0, "")
	bc.SetActiveImage(newImage)

	mec := mgr.executorClient.(*MockExecutorClient)
	mec.AddMockedFunc("ForwardAdminCommand", func(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
		return nil, fmt.Errorf("connection closed")
	})
	mec.AddMockedFunc("Close", func() error { return nil })
	// Reconnect simulates the relaunched (new-image) enclave handshaking.
	mec.AddMockedFunc("Connect", func(ctx context.Context, tag string) error {
		mgr.setExecutorIdentity(newPcr0, "")
		mgr.completeExecutorHandshake(nil)
		return nil
	})

	// Tick 1: mismatch -> drain.
	skip, err := mgr.maintainExecutorForActiveImage(ctx)
	require.NoError(t, err)
	require.True(t, skip)
	require.True(t, mgr.isDraining())

	// Tick 2: draining -> reconnect to new enclave, whose PCR0 matches -> resume.
	skip, err = mgr.maintainExecutorForActiveImage(ctx)
	require.NoError(t, err)
	require.False(t, skip, "after reconnect to the matching image, dispatch resumes")
	require.False(t, mgr.isDraining())
	require.Equal(t, newPcr0, mgr.RunningPcr0())
}
