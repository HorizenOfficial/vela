package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/logserver"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/horizen-pes/pkg/storage/versioned_leveldb"
	"github.com/horizen-pes/pkg/wasm"
	"github.com/stretchr/testify/require"
)

var commParams = common.CommunicationParams{RequestTimeoutSec: 30}

type SystemTestSuite struct {
	t                  *testing.T
	manager            manager.Manager
	executor           executor.Executor
	blockchainClient   *blockchain.MockClient
	dataLayer          storage.DataLayer
	eventChannel       chan interface{}
	ctx                context.Context
	cancel             context.CancelFunc
	executorCommKey    *cryptotypes.PrivateKeyP521      // Executor's communication key for testing
	executorSigningKey *cryptotypes.PrivateKeySecp256k1 // Executor's signing key for testing
	dbPath             string
	reportsPath        string
	log                logger.Logger
}

func NewSystemTestSuite(t *testing.T, appType string, mgrLog logger.Logger, excLog logger.Logger) *SystemTestSuite {
	// log is passed from outside, the log settings in the manager configuration does not affect it.
	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	keySet, newRecoveryData, err := executor.GenerateEnclaveKeySet(execConfig.KeySetRecoveryType)
	require.NoError(t, err)
	return NewSystemTestSuiteWithConfigs(t, appType, mgrConfig, execConfig, keySet, newRecoveryData, mgrLog, excLog)
}

func NewSystemTestSuiteWithConfigs(
	t *testing.T,
	appType string,
	mgrConfig *manager.Config,
	execConfig *executor.Config,
	keySet *executor.EnclaveKeySet,
	recoveryData *common.EnclaveKeySetRecovery,
	mgrLog logger.Logger,
	excLog logger.Logger,
) *SystemTestSuite {
	ctx, cancel := context.WithCancel(context.Background())

	// Normalize channel params for tests: default configs may use vsock, but tests run over TCP.
	var tcpParams common.TcpChannelConnectionParams
	switch p := execConfig.ChannelParams.(type) {
	case common.TcpChannelConnectionParams:
		tcpParams = p
	case common.VSockChannelConnectionParams:
		tcpParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: p.Port}
		execConfig.ChannelParams = tcpParams
		execConfig.ChannelType = "tcp"
	default:
		t.Fatal("unsupported executor channel params type")
	}
	switch p := mgrConfig.ChannelParams.(type) {
	case common.TcpChannelConnectionParams:
		// keep as is
	case common.VSockChannelConnectionParams:
		mgrConfig.ChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: p.Port}
		mgrConfig.ChannelType = "tcp"
	default:
		t.Fatal("unsupported manager channel params type")
	}

	// Create mock components
	blockchainClient := blockchain.NewMockClient()
	// Create an executor client (TCP for testing)
	factory := communication.NewTCPConnectionFactory(tcpParams.Url())
	executorClient := communication.NewClient(factory, commParams, mgrLog)

	// Create manager
	var err error
	var reportsPath string = ""
	if mgrConfig.DeanonymizationReportPath != "" {
		// Create a temporary directory for reports, we overwrite this optional setting
		// because this is a test environment
		reportsPath, err = os.MkdirTemp("", "test-reports")
		require.NoError(t, err)
		mgrConfig.DeanonymizationReportPath = reportsPath
	}

	// Create a temporary directory for the database
	dbPath, err := os.MkdirTemp("", "horizen-pes-test-db")
	require.NoError(t, err)

	cfg := versioned_leveldb.VersionedLevelDBConfig{
		DBPath:         dbPath,
		VersionsToKeep: mgrConfig.DataLayerNumOfVersions,
	}

	var dataLayer storage.DataLayer = nil
	if mgrConfig.DataLayerType == "mockdb" {
		dataLayer = mockdb.NewMockDataLayer()
	} else {
		dataLayer, err = versioned_leveldb.NewVersionedLevelDBDataLayer(cfg)
		require.NoError(t, err)
	}

	mgr := manager.NewSecureProcessorManager(mgrConfig, blockchainClient, dataLayer, executorClient, nil, mgrLog)

	logserver.StartLogServer(
		ctx,
		logserver.LogServerConfig{
			TCPAddr:        mgrConfig.LogServerTCPAddress,
			VSockAddr:      mgrConfig.LogServerVSockAddress,
			LogFilePath:    mgrConfig.LogServerLogFile,
			ConsoleEnabled: mgrConfig.LogServerConsole,
			ConsoleLevel:   mgrConfig.LogServerConsoleLevel,
			FileLevel:      mgrConfig.LogServerFileLevel,
		},
	)

	// Create executor
	server := communication.NewServer(factory, commParams, excLog)
	var runtime executor.Runtime
	switch appType {
	case "mock-runtime":
		t.Log("mock app type: ", appType)
		runtime = executor.NewMockRuntime(excLog)
	default:
		t.Log("wasm app type: ", appType)
		runtime = wasm.NewWasmtimeRuntime(excLog)
	}

	// Create the executor
	exec, err := executor.NewStatelessExecutor(execConfig, runtime, server, excLog)
	require.NoError(t, err)

	if keySet != nil && recoveryData != nil {
		err := dataLayer.StoreEnclaveKeySetRecovery(ctx, recoveryData)
		require.NoError(t, err)
	}

	// Create event channel
	eventChannel := make(chan interface{}, 100)
	blockchainClient.SubscribeToEvents(ctx, eventChannel)

	suite := &SystemTestSuite{
		t:                t,
		manager:          mgr,
		executor:         exec,
		blockchainClient: blockchainClient,
		dataLayer:        dataLayer,
		eventChannel:     eventChannel,
		ctx:              ctx,
		cancel:           cancel,
		dbPath:           dbPath,
		reportsPath:      reportsPath,
		log:              mgrLog,
	}

	if keySet != nil {
		suite.executorCommKey = &keySet.CommunicationKey
		suite.executorSigningKey = &keySet.SigningKey
	}

	return suite
}

func (s *SystemTestSuite) StartManager() error {
	errChan := make(chan error, 1)

	go func() {
		if err := s.manager.Start(s.ctx); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	// Wait for a result from the goroutine
	if err := <-errChan; err != nil {
		s.log.Info("Manager failed to start: %v", err)
		return err
	}

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) StartExecutor() error {
	errChan := make(chan error, 1)

	go func() {
		if err := s.executor.Start(s.ctx); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	// Wait for a result from the goroutine
	if err := <-errChan; err != nil {
		s.t.Fatalf("Executor failed to start: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) SubmitRequest(req *common.Request) error {
	return s.blockchainClient.SendRequestToChain(s.ctx, req) //use test function in mock_client
}

func (s *SystemTestSuite) WaitForAppStateInDB(appID common.ApplicationIdType, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			state, err := s.dataLayer.GetApplicationState(s.ctx, appID)
			if err == nil {
				s.t.Log("AppID found on data layer", appID)
				return state, nil
			}
			// this will print every tick
			//s.t.Log("err: ", err.Error())
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %d", appID)

		}
	}
}

func (s *SystemTestSuite) GetFailedRequest() []*common.Request {
	failedReq := s.blockchainClient.GetFailedRequests()
	return failedReq
}

func (s *SystemTestSuite) WaitForAppStateInBlockchain(appID common.ApplicationIdType, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			state, err := s.blockchainClient.GetApplicationState(s.ctx, appID)
			if err == nil {
				return state, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %d in blockchain", appID)

		}
	}
}

func (s *SystemTestSuite) AssertRequestCompleted(requestID common.RequestIdType, timeout time.Duration) error {
	return s.blockchainClient.WaitForRequestCompletion(requestID, timeout)
}

// WaitForEvent waits for a specific event to be published for a user.
// If eventSubType is empty, any subtype is accepted.
func (s *SystemTestSuite) WaitForEvent(userID ethCommon.Address, eventSubType string, timeout time.Duration) (*common.Event, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case event := <-s.eventChannel:
			if evt, ok := event.(common.Event); ok && evt.UserID == userID && (eventSubType == "" || evt.EventSubType == eventSubType) {
				s.log.Info("TESTING: Received event: %v", evt)
				return &evt, nil
			} else {
				s.log.Info("TESTING: Received unexpected event: %v", event)
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for event for user %s", userID)
		}
	}
}

// WaitForDeanonymizationReport waits for a deanonymization report to be generated and saved to the filesystem
func (s *SystemTestSuite) WaitForDeanonymizationReport(reportID common.RequestIdType, timeout time.Duration) (*common.DeanonymizationReport, error) {
	if s.reportsPath == "" {
		return nil, fmt.Errorf("reports path not configured")
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			// Check if deanonymization report exists in filesystem
			// We need to find the report file by iterating through possible app IDs
			// since we don't have the app ID in this function
			files, err := os.ReadDir(s.reportsPath)
			if err != nil {
				continue
			}
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				// Report filename format: report_<appID>_<requestID>.json
				// Check if filename contains the requestID
				if !(strings.Contains(f.Name(), reportID.String())) {
					continue
				}
				reportPath := s.reportsPath + "/" + f.Name()
				data, err := os.ReadFile(reportPath)
				if err != nil {
					continue
				}
				var report common.DeanonymizationReport
				if err := json.Unmarshal(data, &report); err != nil {
					continue
				}
				return &report, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for deanonymization report %s", reportID)
		}
	}
}

// GetReportsPath returns the path where deanonymization reports are saved
func (s *SystemTestSuite) GetReportsPath() string {
	return s.reportsPath
}

// WaitForWithdrawal waits for a withdrawal to be processed
func (s *SystemTestSuite) WaitForWithdrawal(appID common.ApplicationIdType, timeout time.Duration) (*common.Withdrawal, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			// Check if withdrawal exists in blockchain
			withdrawals, err := s.blockchainClient.GetWithdrawals(s.ctx, appID)
			if err == nil && withdrawals != nil && len(*withdrawals) > 0 {
				return &(*withdrawals)[0], nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for withdrawal for app %d", appID)
		}
	}
}

func (s *SystemTestSuite) GetRequestUpdatePayload(reqId common.RequestIdType) (*common.UpdatePayload, error) {
	// Get the update payload for the request
	return s.blockchainClient.GetRequestUpdatePayload(s.ctx, reqId)
}

// GetExecutorCommunicationKey returns the executor's communication public key for encryption
func (s *SystemTestSuite) GetExecutorCommunicationKey() (*cryptotypes.PublicKeyP521, error) {
	if s.executorCommKey == nil {
		return nil, fmt.Errorf("executor communication key not initialized")
	}
	return s.executorCommKey.PublicKey(), nil
}

func (s *SystemTestSuite) GetExecutorSigningKey() (*cryptotypes.PublicKeySecp256k1, error) {
	if s.executorSigningKey == nil {
		return nil, fmt.Errorf("executor signing key not initialized")
	}
	return s.executorSigningKey.PublicKey(), nil
}

func (s *SystemTestSuite) GetDataLayer() storage.DataLayer {
	return s.dataLayer
}

func (s *SystemTestSuite) LoadWasmModule(t *testing.T, moduleFilename string) []byte {
	wasmBytes, err := os.ReadFile(moduleFilename)
	require.NoError(t, err, "Failed to read WASM file")
	return wasmBytes
}

func (s *SystemTestSuite) Cleanup() error {
	s.cancel()

	if s.manager != nil {
		s.manager.Stop()
	}

	if s.executor != nil {
		s.executor.Close()
	}

	s.blockchainClient.ClearAllData()

	// Remove the temporary database directory
	if s.dbPath != "" {
		os.RemoveAll(s.dbPath)
	}

	// Remove the temporary reports directory
	if s.reportsPath != "" {
		os.RemoveAll(s.reportsPath)
	}

	return nil
}

func PrettyPrintJSON(data interface{}) string {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "Invalid data: could not marshal into a json"
	}
	return string(prettyJSON)
}
