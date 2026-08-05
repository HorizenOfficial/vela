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
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
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
	stateRoot := generateRandomStateRoot()
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
		// Mirrors the real executor: app existence is validated on-chain, so a
		// nil state is a hard failure, not a signed error payload.
		return nil, nil, nil, fmt.Errorf("state not found for application %d", req.ApplicationID)
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

	stateRoot := generateRandomStateRoot()

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

func (m *MockExecutorClient) SendBatchProcessRequest(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasmModule []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
	if f, ok := m.GetMockedFunc("SendBatchProcessRequest"); ok {
		return f.(func(context.Context, []*common.Request, *common.ApplicationState, []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error))(ctx, requests, appState, wasmModule)
	}
	if appState == nil {
		return nil, nil, nil, nil, fmt.Errorf("state not found for application %d", requests[0].ApplicationID)
	}
	prev := appState.StateRoot
	results := make([]*common.UpdatePayload, 0, len(requests))
	last := prev
	var reports []*common.DeanonymizationReport
	
	for _, req := range requests {
		if string(req.Payload) == "invalid" {
			failurePayload := &common.UpdatePayload{
				ApplicationID: req.ApplicationID,
				RequestID:     req.RequestID,
				PrevStateRoot: prev,
				NewStateRoot:  prev,
				ErrorCode:     uint8(apperrors.CategoryInternalMeta.Category),
				ErrorMsg:      "invalid request",
			}
			results = append(results, failurePayload)
			break
		}
		nr := generateRandomStateRoot()
		results = append(results, &common.UpdatePayload{
			ApplicationID: req.ApplicationID,
			RequestID:     req.RequestID,
			PrevStateRoot: prev,
			NewStateRoot:  nr,
		})
		if req.RequestType == common.Deanonymize {
			reports = append(reports, &common.DeanonymizationReport{
				ApplicationID:   req.ApplicationID,
				ReportID:        req.RequestID,
				Authority:       req.Sender,
				EncryptedReport: []byte(`{"accounts":{}}`),
			})
		}
		prev = nr
		last = nr
	}
	final := &common.ApplicationState{ApplicationID: appState.ApplicationID, StateRoot: last, EncryptedState: []byte("enc")}
	return results, []byte("batch-sig"), final, reports, nil
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
func generateRandomStateRoot() [32]byte {
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

func TestProcessBatchFromChain(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy request
	request := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 0, len(completedRequests), "expected 0 completed request")

	err = manager.processBatchFromChain(context.Background())
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

	err = manager.processBatchFromChain(context.Background())
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

	err = manager.processBatchFromChain(context.Background())
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
	manager.config.MaxBatchSize = 1

	for i := 0; i < 4; i++ {
		err = manager.processBatchFromChain(context.Background())
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

func TestProcessBatchWithFailure(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	seedDeployedApp(t, mockBCClient, manager, ApplicationId)

	request := createRequestWithPayload(common.Process, ApplicationId, []byte("invalid"))
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Test that if it is a failure payload returned by the executor, submitStateOnChain is called but the state is not stored in the data layer
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, *common.ApplicationState, *common.WASMData) error {
		t.Fatal("Store should not be called if the executor returned a failure payload")
		return nil
	})

	appState, _ := manager.dataLayer.GetApplicationState(context.Background(), ApplicationId)

	err = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, appState.StateRoot)
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed request")

	failedRequests := mockBCClient.GetFailedRequests()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, request.RequestID, failedRequests[0].RequestID, "Wrong requestID")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

}

func TestProcessBatchWithErrors(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy the application first
	seedDeployedApp(t, mockBCClient, manager, ApplicationId)

	oldDbVersion, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)

	appState, err := manager.dataLayer.GetApplicationState(context.Background(), ApplicationId)
	require.NoError(t, err)
	initialStateRoot := appState.StateRoot
	// Simulate application state not found. SendProcessRequest is still called with a nil
	// state, but the executor rejects it with a hard error: the request stays pending.
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, common.ApplicationIdType) (*common.ApplicationState, error) {
		return nil, storageErrors.ErrNotFound("application state not found")
	})
	request := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	failure := manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	
	require.Error(t, failure)
	require.ErrorContains(t, failure, "state not found for application")

	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	require.Empty(t, mockBCClient.GetFailedRequests())
	pendingRequests, _ := mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")

	expectedError := "error"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, common.ApplicationIdType) (*common.ApplicationState, error) {
		return nil, errors.New(expectedError)
	})
	failure = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	require.Error(t, failure)
	require.ErrorContains(t, failure, expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetApplicationState")

	// Failure in GetWasmCode, stop processing and return the error
	expectedError = "wasm bytecode not found for application"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetWASMBytecode", func(context.Context, common.ApplicationIdType) ([]byte, error) {
		return nil, errors.New(expectedError)
	})
	failure = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	require.Error(t, failure)
	require.ErrorContains(t, failure, expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetWASMBytecode")

	// Test failure in executor
	expectedError = "failed to execute app"
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendBatchProcessRequest",
		func(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			return nil, nil, nil, nil, errors.New(expectedError)
		})

	failure = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")
	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")

	manager.executorClient.(*MockExecutorClient).RemoveMockedFunc("SendBatchProcessRequest")

	// Test data layer failure, stop processing and return the error
	expectedError = "failed to store state"
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("Store", func(context.Context, []byte, *common.ApplicationState, *common.WASMData) error {
		return errors.New(expectedError)
	})

	failure = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

	// Test blockchain failure for errors that can be due to reorgs.
	// The errors that can be due to reorgs are:
	// - InvalidRequestId
	// - InvalidStateRoot
	// - InvalidApplicationId
	// - NonceTooLow
	// The local db should be reverted to the previous state

	mockBCClient.AddMockedFunc("SubmitBatchStateUpdate", func(ctx context.Context, updates []*common.UpdatePayload, sig []byte) error {
		return blockchain.ReorgError{}
	})
	failure = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	require.Nil(t, failure)
	// Check that the local db has been reverted to the initial state
	newDbVersion, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	// Test blockchain failure for any other errors but reorgs.
	// The local db should be reverted to the previous state
	expectedError = "some other error"
	mockBCClient.AddMockedFunc("SubmitBatchStateUpdate", func(ctx context.Context, updates []*common.UpdatePayload, sig []byte) error {
		return errors.New(expectedError)
	})	

	failure = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	// Check that the local db has been reverted to the initial state
	newDbVersion, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	// Same test but with an error payload for the SubmitBatchStateUpdate
	request.Payload = []byte("invalid")
	failure = manager.processBatch(context.Background(), ApplicationId, []*common.Request{request}, initialStateRoot)
	require.Error(t, failure)
	require.Contains(t, failure.Error(), expectedError)

	// Check that the local db has been reverted to the initial state
	newDbVersion, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.Equal(t, oldDbVersion, newDbVersion)

	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("Store")

}

// TestProcessDeanonymizationViaProcessRequest tests the deanonymization flow which is now
// handled through processBatchFromChain with RequestType = Deanonymize
func TestProcessDeanonymizationViaProcessRequest(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Deploy the application first
	seedDeployedApp(t, mockBCClient, manager, ApplicationId)
	// Create a Deanonymize request (now handled via processBatch)
	request := createRequest(common.Deanonymize, ApplicationId)
	err := mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)

	// Process the deanonymization request via processBatchFromChain
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	completedRequests := mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed requests")
	require.Equal(t, request.RequestID, completedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, request.ApplicationID, completedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, common.Deanonymize, completedRequests[1].RequestType, "Wrong RequestType")
}

func TestProcessBatchRequestFromChainWithReorgs(t *testing.T) {
	mockBCClient, manager := setupTest(t)

	// Prepare initial state in the database
	_, _, initialStateRootOnChain, err := mockBCClient.GetPendingRequestsWithStateRoot(context.Background(), 1)
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

	manager.config.MaxBatchSize = 1
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 2, len(pendingRequests), "expected 2 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completedRequests), "expected 1 completed request")

	db_version, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)

	appId, nextPendingReq, stateRootOnChain1, err := mockBCClient.GetPendingRequestsWithStateRoot(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, appId, ApplicationId)
	require.True(t, bytes.Equal(stateRootOnChain1[:], db_version), "State root in DB should be equal to state root on chain")
	require.Equal(t, len(nextPendingReq), 1)
	require.Equal(t, ApplicationId, nextPendingReq[0].ApplicationID)
	require.Equal(t, request2.RequestID, nextPendingReq[0].RequestID)

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	pendingRequests, _ = mockBCClient.GetPendingRequests(context.Background())
	require.Equal(t, 1, len(pendingRequests), "expected 1 pending request")
	completedRequests = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completedRequests), "expected 2 completed requests")

	db_version, err = manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)

	appId, nextPendingReq, stateRootOnChain2, err := mockBCClient.GetPendingRequestsWithStateRoot(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, appId, ApplicationId)
	require.True(t, bytes.Equal(stateRootOnChain2[:], db_version), "State root in DB should be equal to state root on chain")
	require.Equal(t, len(nextPendingReq), 1)
	require.Equal(t, ApplicationId, nextPendingReq[0].ApplicationID)
	require.Equal(t, request3.RequestID, nextPendingReq[0].RequestID)

	// Now simulate a reorg on chain, by making GetNextPendingRequest to always return the first request and initial state root

	mockedGetPendingRequestsWithStateRoot := func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return request1.ApplicationID, []*common.Request{request1}, initialStateRootOnChain, nil
	}

	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", mockedGetPendingRequestsWithStateRoot)

	// SubmitBatchStateUpdate should not be called in case of reorg
	mockedSubmitBatchStateUpdatePanics := func(context.Context, []*common.UpdatePayload, []byte) error {
		t.Fatal("SubmitBatchStateUpdate should not be called in case of reorg")
		return nil
	}
	mockBCClient.AddMockedFunc("SubmitBatchStateUpdate", mockedSubmitBatchStateUpdatePanics)

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)
	require.False(t, manager.endReorgTime.IsZero(), "endReorgTime should be set when reorg is detected")
	currentEndReorgTime := manager.endReorgTime

	// Try with the second request and state root
	mockedGetPendingRequestsWithStateRoot = func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return request2.ApplicationID, []*common.Request{request2}, stateRootOnChain1, nil
	}
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", mockedGetPendingRequestsWithStateRoot)

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)
	require.Equal(t, currentEndReorgTime, manager.endReorgTime, "endReorgTime should not change if reorg is not yet resolved")

	// Solve the reorg and process the last request
	mockBCClient.RemoveMockedFunc("GetPendingRequestsWithStateRoot")
	mockBCClient.RemoveMockedFunc("SubmitBatchStateUpdate")

	err = manager.processBatchFromChain(context.Background())
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

	mockedGetPendingRequestsWithStateRoot = func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return request4.ApplicationID, []*common.Request{request4}, [32]byte{0x11, 0x66}, nil
	}
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", mockedGetPendingRequestsWithStateRoot)

	err = manager.processBatchFromChain(context.Background())
	require.Error(t, err, "Should return error due to unrecoverable disalignment between DB and chain")

	// test reorg not solved within timeout. The local db is reverted to the same state of the chain and the request is executed
	//Remove all old requests and reset to the initial state root

	mockBCClient.ClearAllData()
	//Re-submit old request
	err = mockBCClient.SendRequestToChain(context.Background(), request1)
	require.NoError(t, err)

	manager.config.ReorgTimeout = 1 // 1 second

	// First processRequestFromChain sets the reorg timeout
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	// wait for more than reorg timeout
	// Instead of sleeping, we will simulate the time.Sleep by manipulating endReorgTime
	manager.endReorgTime = manager.endReorgTime.Add(-2 * time.Second) // go back in time by 2 seconds

	err = manager.processBatchFromChain(context.Background())
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
	seedDeployedApp(t, mockBCClient, manager, ApplicationId)

	//**********************
	// Check that if GetPendingRequestsWithStateRoot returns an error, processRequestFromChain doesn't execute the request and doesn't return an error
	//**********************

	mockedGetPendingRequestsWithStateRoot := func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return 0, nil, [32]byte{}, fmt.Errorf("GetPendingRequestsWithStateRoot error")
	}
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", mockedGetPendingRequestsWithStateRoot)

	// SubmitBatchStateUpdate should not be called in case of an error in GetPendingRequestsWithStateRoot
	mockedSubmitBatchStateUpdatePanics := func(context.Context, []*common.UpdatePayload, []byte) error {
		t.Fatal("SubmitBatchStateUpdate should not be called in case of error")
		return nil
	}
	mockBCClient.AddMockedFunc("SubmitBatchStateUpdate", mockedSubmitBatchStateUpdatePanics)

	err := manager.processBatchFromChain(context.Background())
	require.NoError(t, err, "processBatchFromChain should not return an error if GetPendingRequestsWithStateRoot fails")

	request1 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request1)
	require.NoError(t, err)
	mockBCClient.RemoveMockedFunc("GetPendingRequestsWithStateRoot")

	//**********************
	// Check that if LastVersionID returns an error, processBatchFromChain doesn't execute the request and doesn't return an error
	//**********************

	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("LastVersionID", func(common.ApplicationIdType) ([]byte, error) {
		return nil, fmt.Errorf("LastVersionID error")
	})

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err, "processBatchFromChain should not return an error if LastVersionID fails")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("LastVersionID")

	//**********************
	// Check that if ListVersions returns an error, processBatchFromChain doesn't execute the request and doesn't return an error
	//**********************
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("ListVersions", func(common.ApplicationIdType) ([][]byte, error) {
		return nil, fmt.Errorf("ListVersions error")
	})

	// Setup a fake reorg situation
	mockedGetPendingRequestsWithStateRoot = func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return request1.ApplicationID, []*common.Request{request1}, [32]byte{0x11, 0x66}, nil
	}
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", mockedGetPendingRequestsWithStateRoot)

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err, "processBatchFromChain should not return an error if ListVersions fails")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("ListVersions")
	mockBCClient.RemoveMockedFunc("GetPendingRequestsWithStateRoot")

	//**********************
	// Check that if processBatch returns an error, processBatchFromChain doesn't execute the request and doesn't return an error
	//**********************

	// Setup a fake situation to make processRequest return an error
	manager.dataLayer.(*mockdb.MockDataLayer).AddMockedFunc("GetApplicationState", func(context.Context, common.ApplicationIdType) (*common.ApplicationState, error) {
		return nil, fmt.Errorf("error")
	})

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err, "processBatchFromChain should not return an error if processBatch fails")

	manager.dataLayer.(*mockdb.MockDataLayer).RemoveMockedFunc("GetApplicationState")

}

func setupTest(t *testing.T) (*blockchain.MockClient, *SecureProcessorManager) {
	config := Config{
		ReorgTimeout:        60,
		MaxBatchSize:        5,
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
	seedDeployedApp(t, mockBCClient, manager, ApplicationId)

	// DeanonymizationReportPath is set, so the report should be saved to the filesystem
	// Create a deanonymization request
	request := createRequest(common.Deanonymize, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), request)
	require.NoError(t, err)
	manager.config.DeanonymizationReportPath = tempDir
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)
	completedRequests := mockBCClient.GetCompletedRequests()
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

	// Deploy App 2
	deployReq2 := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	err = mockBCClient.SendRequestToChain(context.Background(), deployReq2)
	require.NoError(t, err)

	//Check that every deploy request is executed separately and independently
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)
	completed := mockBCClient.GetCompletedRequests()
	require.Equal(t, 1, len(completed), "App 1 deploy should complete")

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)
	completed = mockBCClient.GetCompletedRequests()
	require.Equal(t, 2, len(completed), "App 2 deploy should complete")

	// Verify independent state roots
	stateRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	stateRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.False(t, bytes.Equal(stateRoot1, stateRoot2), "App 1 and App 2 should have different state roots")

	// Create a request for App 1
	processReq1 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq1)
	require.NoError(t, err)

	// Create a request for App 2
	processReq2 := createRequest(common.Process, ApplicationId2)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq2)
	require.NoError(t, err)

	// Prcess a batch of requests (should process App 1 first)
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	// Verify App 1 state root changed but App 2 unchanged
	newStateRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.False(t, bytes.Equal(stateRoot1, newStateRoot1), "App 1 state root should change after processing")

	unchangedStateRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRoot2, unchangedStateRoot2), "App 2 state root should be unchanged")

	// Prcess a batch of requests (should process App 2)
	err = manager.processBatchFromChain(context.Background())
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
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	processReq1 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq1)
	require.NoError(t, err)
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	stateRoot1, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)

	// Deploy App 2 — no prior state
	deployReq2 := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	err = mockBCClient.SendRequestToChain(context.Background(), deployReq2)
	require.NoError(t, err)
	err = manager.processBatchFromChain(context.Background())
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
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	deployReq2 := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	err = mockBCClient.SendRequestToChain(context.Background(), deployReq2)
	require.NoError(t, err)
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	// Process a request for App 1 (creates a second version)
	processReq1 := createRequest(common.Process, ApplicationId)
	err = mockBCClient.SendRequestToChain(context.Background(), processReq1)
	require.NoError(t, err)
	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	stateRoot1AfterProcess, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	stateRoot2AfterDeploy, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)

	// Simulate reorg for App 1 — chain returns old stateRoot (zero, as if rolled back to before deploy)
	mockedReorg := func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) { 
		return processReq1.ApplicationID, []*common.Request{processReq1}, [32]byte{}, nil
	}
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", mockedReorg)
	mockBCClient.AddMockedFunc("SubmitBatchStateUpdate", func(context.Context, []*common.UpdatePayload, []byte) error {
		t.Fatal("SubmitBatchStateUpdate should not be called during reorg")
		return nil
	})

	err = manager.processBatchFromChain(context.Background())
	require.NoError(t, err)

	// Reorg timer should be set (chain-level, triggered by App 1's mismatch)
	require.False(t, manager.endReorgTime.IsZero(), "reorg timer should be set")

	// App 2 state should be untouched
	unchangedRoot2, err := manager.dataLayer.LastVersionID(ApplicationId2)
	require.NoError(t, err)
	require.True(t, bytes.Equal(stateRoot2AfterDeploy, unchangedRoot2), "App 2 state should be unchanged during App 1 reorg")

	// Resolve reorg — remove mock, next poll returns App 1's correct state
	mockBCClient.RemoveMockedFunc("GetPendingRequestsWithStateRoot")
	mockBCClient.RemoveMockedFunc("SubmitBatchStateUpdate")

	// Restore App 1's request to pending and set correct state root
	err = mockBCClient.SendRequestToChain(context.Background(), createRequest(common.Process, ApplicationId))
	require.NoError(t, err)

	err = manager.processBatchFromChain(context.Background())
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
	require.NoError(t, manager.processBatchFromChain(ctx))

	processA := createRequest(common.Process, ApplicationId)
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, processA))
	require.NoError(t, manager.processBatchFromChain(ctx))

	deployB := createDeployRequestWithWASM(t, manager, ApplicationId2, []byte{0x02})
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, deployB))
	require.NoError(t, manager.processBatchFromChain(ctx))

	processBReq := createRequest(common.Process, ApplicationId2)
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, processBReq))
	require.NoError(t, manager.processBatchFromChain(ctx))

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
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return processA.ApplicationID, []*common.Request{processA}, reorgStateRootA, nil
	})
	mockBCClient.AddMockedFunc("SubmitBatchStateUpdate", func(context.Context, []*common.UpdatePayload, []byte) error {
		t.Fatal("SubmitBatchStateUpdate should not be called during reorg")
		return nil
	})

	require.NoError(t, manager.processBatchFromChain(ctx))
	require.False(t, manager.endReorgTime.IsZero(), "reorg timer should be set after App A mismatch")

	// Second poll: App B sees reorg (chain returns deploy root).
	// The shared timer is already running — it should NOT be reset.
	var reorgStateRootB [32]byte
	copy(reorgStateRootB[:], stateRootBAfterDeploy)
	timerBefore := manager.endReorgTime
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return processBReq.ApplicationID, []*common.Request{processBReq}, reorgStateRootB, nil
	})

	require.NoError(t, manager.processBatchFromChain(ctx))
	require.Equal(t, timerBefore, manager.endReorgTime, "shared timer should not change on second mismatch")

	// ── Step 3: Reorg #1 resolves — only App A gets a request ──

	mockBCClient.RemoveMockedFunc("GetPendingRequestsWithStateRoot")
	mockBCClient.RemoveMockedFunc("SubmitBatchStateUpdate")

	// App A gets a new request; state roots will match → shared timer cleared.
	newReqA := createRequest(common.Process, ApplicationId)
	require.NoError(t, mockBCClient.SendRequestToChain(ctx, newReqA))
	require.NoError(t, manager.processBatchFromChain(ctx))

	require.True(t, manager.endReorgTime.IsZero(),
		"shared timer should be cleared when App A's state roots match (reorg resolved)")

	// ── Step 4: Reorg #2 — App B sees a NEW mismatch ──

	// The chain has reorged again: it reports the deploy root for App B.
	processB2 := createRequest(common.Process, ApplicationId2)
	mockBCClient.AddMockedFunc("GetPendingRequestsWithStateRoot", func(context.Context, uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
		return processB2.ApplicationID, []*common.Request{processB2}, reorgStateRootB, nil
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

	require.NoError(t, manager.processBatchFromChain(ctx))

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

// ---------------------------------------------------------------------------
// Batch orchestration tests (processBatchFromChain / processBatch)
// ---------------------------------------------------------------------------


// seedDeployedApp deploys an application through the batch poll path so that its
// state and WASM exist in the data layer, and returns the resulting state root.
func seedDeployedApp(t *testing.T, mockBC *blockchain.MockClient, manager *SecureProcessorManager, appID common.ApplicationIdType) [32]byte {
	t.Helper()
	dep := createDeployRequestWithWASM(t, manager, appID, []byte{0x01})
	require.NoError(t, mockBC.SendRequestToChain(context.Background(), dep))
	require.NoError(t, manager.processBatchFromChain(context.Background()))
	root, err := manager.dataLayer.LastVersionID(appID)
	require.NoError(t, err)
	var out [32]byte
	copy(out[:], root)
	return out
}

// mockBatchSuccess installs a SendBatchProcessRequest mock that produces a success
// payload per request, chaining state roots from the current app state.
func mockBatchSuccess(manager *SecureProcessorManager, t *testing.T) {
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendBatchProcessRequest",
		func(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			require.NotNil(t, appState, "batch executor requires a non-nil app state")
			prev := appState.StateRoot
			results := make([]*common.UpdatePayload, 0, len(requests))
			last := prev
			var reports []*common.DeanonymizationReport
			for _, req := range requests {
				nr := generateRandomStateRoot()
				results = append(results, &common.UpdatePayload{
					ApplicationID: req.ApplicationID,
					RequestID:     req.RequestID,
					PrevStateRoot: prev,
					NewStateRoot:  nr,
				})
				if req.RequestType == common.Deanonymize {
					reports = append(reports, &common.DeanonymizationReport{
						ApplicationID:   req.ApplicationID,
						ReportID:        req.RequestID,
						Authority:       req.Sender,
						EncryptedReport: []byte(`{"accounts":{}}`),
					})
				}
				prev = nr
				last = nr
			}
			final := &common.ApplicationState{ApplicationID: appState.ApplicationID, StateRoot: last, EncryptedState: []byte("enc")}
			return results, []byte("batch-sig"), final, reports, nil
		})
}

func TestProcessBatchFromChain_HappyPath(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()

	seedDeployedApp(t, mockBC, manager, ApplicationId)
	mockBatchSuccess(manager, t)

	reqs := []*common.Request{
		createRequest(common.Process, ApplicationId),
		createRequest(common.Process, ApplicationId),
		createRequest(common.Process, ApplicationId),
	}
	for _, r := range reqs {
		require.NoError(t, mockBC.SendRequestToChain(ctx, r))
	}

	require.NoError(t, manager.processBatchFromChain(ctx))

	pending, _ := mockBC.GetPendingRequests(ctx)
	require.Equal(t, 0, len(pending), "all batched requests should be completed in one poll")
	require.Equal(t, 0, len(mockBC.GetFailedRequests()), "no failures expected")
	// deploy + 3 process requests
	require.Equal(t, 4, len(mockBC.GetCompletedRequests()))
}

func TestProcessBatchFromChain_SingleRequestGoesThroughBatchPath(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()

	seedDeployedApp(t, mockBC, manager, ApplicationId)

	called := false
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendBatchProcessRequest",
		func(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			called = true
			require.Len(t, requests, 1, "a single pending request is still routed through the batch path")
			nr := generateRandomStateRoot()
			res := []*common.UpdatePayload{{ApplicationID: requests[0].ApplicationID, RequestID: requests[0].RequestID, PrevStateRoot: appState.StateRoot, NewStateRoot: nr}}
			final := &common.ApplicationState{ApplicationID: appState.ApplicationID, StateRoot: nr, EncryptedState: []byte("enc")}
			return res, []byte("batch-sig"), final, nil, nil
		})

	req := createRequest(common.Process, ApplicationId)
	require.NoError(t, mockBC.SendRequestToChain(ctx, req))
	require.NoError(t, manager.processBatchFromChain(ctx))

	require.True(t, called, "single non-deploy request must go through SendBatchProcessRequest")
	pending, _ := mockBC.GetPendingRequests(ctx)
	require.Equal(t, 0, len(pending))
}

func TestProcessBatchFromChain_DeployRoutedIndividually(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()

	// A batch executor call must NOT happen for a deploy.
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendBatchProcessRequest",
		func(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			t.Fatal("deploy must not be routed through the batch path")
			return nil, nil, nil, nil, nil
		})

	dep := createDeployRequestWithWASM(t, manager, ApplicationId, []byte{0x01})
	require.NoError(t, mockBC.SendRequestToChain(ctx, dep))
	require.NoError(t, manager.processBatchFromChain(ctx))

	pending, _ := mockBC.GetPendingRequests(ctx)
	require.Equal(t, 0, len(pending))
	require.Equal(t, 1, len(mockBC.GetCompletedRequests()))
}

func TestProcessBatchFromChain_EmptyQueue(t *testing.T) {
	mockBC, manager := setupTest(t)
	require.NoError(t, manager.processBatchFromChain(context.Background()))
	pending, _ := mockBC.GetPendingRequests(context.Background())
	require.Equal(t, 0, len(pending))
	require.Equal(t, 0, len(mockBC.GetCompletedRequests()))
}

func TestProcessBatchFromChain_QueueLargerThanMaxBatchSize(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()
	manager.config.MaxBatchSize = 2

	seedDeployedApp(t, mockBC, manager, ApplicationId)
	mockBatchSuccess(manager, t)

	for i := 0; i < 3; i++ {
		require.NoError(t, mockBC.SendRequestToChain(ctx, createRequest(common.Process, ApplicationId)))
	}

	require.NoError(t, manager.processBatchFromChain(ctx))

	pending, _ := mockBC.GetPendingRequests(ctx)
	require.Equal(t, 1, len(pending), "only MaxBatchSize (2) of 3 requests should be processed this poll")
}

func TestProcessBatchFromChain_BatchRevertRollsBack(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()

	deployRoot := seedDeployedApp(t, mockBC, manager, ApplicationId)
	mockBatchSuccess(manager, t)

	// The batch tx reverts.
	mockBC.AddMockedFunc("SubmitBatchStateUpdate", func(ctx context.Context, updates []*common.UpdatePayload, sig []byte) error {
		return fmt.Errorf("tx reverted")
	})

	for i := 0; i < 2; i++ {
		require.NoError(t, mockBC.SendRequestToChain(ctx, createRequest(common.Process, ApplicationId)))
	}

	require.NoError(t, manager.processBatchFromChain(ctx))

	// Requests remain pending and local state is rolled back to the pre-batch root.
	pending, _ := mockBC.GetPendingRequests(ctx)
	require.Equal(t, 2, len(pending), "requests stay pending when the batch tx reverts")
	root, err := manager.dataLayer.LastVersionID(ApplicationId)
	require.NoError(t, err)
	require.Equal(t, deployRoot[:], root, "local state must be rolled back to the pre-batch root")
}

func TestProcessBatchFromChain_HardFailureOnFirstRequest(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()

	seedDeployedApp(t, mockBC, manager, ApplicationId)

	// No payloads returned: nothing to submit.
	submitted := false
	mockBC.AddMockedFunc("SubmitBatchStateUpdate", func(ctx context.Context, updates []*common.UpdatePayload, sig []byte) error {
		submitted = true
		return nil
	})
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendBatchProcessRequest",
		func(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			return nil, nil, nil, nil, nil
		})

	for i := 0; i < 2; i++ {
		require.NoError(t, mockBC.SendRequestToChain(ctx, createRequest(common.Process, ApplicationId)))
	}

	require.NoError(t, manager.processBatchFromChain(ctx))

	require.False(t, submitted, "nothing should be submitted when the executor returns no payloads")
	pending, _ := mockBC.GetPendingRequests(ctx)
	require.Equal(t, 2, len(pending), "all requests stay pending on a hard failure on the first request")
}

func TestProcessBatchFromChain_HardFailureMidBatch(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()

	seedDeployedApp(t, mockBC, manager, ApplicationId)

	// 3 requests submitted, but the executor stops after 2 (hard failure on the 3rd).
	manager.executorClient.(*MockExecutorClient).AddMockedFunc("SendBatchProcessRequest",
		func(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasm []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
			require.GreaterOrEqual(t, len(requests), 3)
			prev := appState.StateRoot
			results := make([]*common.UpdatePayload, 0, 2)
			last := prev
			for _, req := range requests[:2] {
				nr := generateRandomStateRoot()
				results = append(results, &common.UpdatePayload{ApplicationID: req.ApplicationID, RequestID: req.RequestID, PrevStateRoot: prev, NewStateRoot: nr})
				prev = nr
				last = nr
			}
			final := &common.ApplicationState{ApplicationID: appState.ApplicationID, StateRoot: last, EncryptedState: []byte("enc")}
			return results, []byte("batch-sig"), final, nil, nil
		})

	for i := 0; i < 3; i++ {
		require.NoError(t, mockBC.SendRequestToChain(ctx, createRequest(common.Process, ApplicationId)))
	}

	require.NoError(t, manager.processBatchFromChain(ctx))

	pending, _ := mockBC.GetPendingRequests(ctx)
	require.Equal(t, 1, len(pending), "the request that caused the hard stop stays pending")
	require.Equal(t, 3, len(mockBC.GetCompletedRequests()), "deploy + 2 processed requests")
}

func TestProcessBatchFromChain_SavesDeanonymizationReports(t *testing.T) {
	mockBC, manager := setupTest(t)
	ctx := context.Background()

	seedDeployedApp(t, mockBC, manager, ApplicationId)
	mockBatchSuccess(manager, t)

	deanonReq := createRequest(common.Deanonymize, ApplicationId)
	require.NoError(t, mockBC.SendRequestToChain(ctx, createRequest(common.Process, ApplicationId)))
	require.NoError(t, mockBC.SendRequestToChain(ctx, deanonReq))

	require.NoError(t, manager.processBatchFromChain(ctx))

	reportFilePath := filepath.Join(manager.config.DeanonymizationReportPath, common.ReportFilename(deanonReq.ApplicationID, deanonReq.RequestID))
	_, err := os.Stat(reportFilePath)
	require.NoError(t, err, "deanonymization report from the batch should be saved to disk")
}
