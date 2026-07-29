package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/blockchain"
	"github.com/HorizenOfficial/vela/pkg/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/communication"
	"github.com/HorizenOfficial/vela/pkg/executor"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/logserver"
	"github.com/HorizenOfficial/vela/pkg/manager"
	"github.com/HorizenOfficial/vela/pkg/storage"
	"github.com/HorizenOfficial/vela/pkg/storage/mockdb"
	"github.com/HorizenOfficial/vela/pkg/storage/versioned_leveldb"
	"github.com/HorizenOfficial/vela/pkg/wasm"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var commParams = common.CommunicationParams{RequestTimeoutSec: 30}

// TestSuiteCore contains the shared infrastructure for both SystemTestSuite (mock)
// and FullStackSystemTestSuite (real chain). Both suite types embed this struct
// and inherit its methods via Go's method promotion.
//
// Exported so that the fullstack package (pkg/testutil/fullstack) can embed it
// without introducing a transitive simulated-backend dependency into this package.
type TestSuiteCore struct {
	t                  *testing.T
	manager            manager.Manager
	executor           executor.Executor
	dataLayer          storage.DataLayer
	eventChannel       chan interface{}
	ctx                context.Context
	cancel             context.CancelFunc
	executorCommKey    *cryptotypes.PrivateKeyP521      // Executor's communication key for testing
	executorSigningKey *cryptotypes.PrivateKeySecp256k1 // Executor's signing key for testing
	dbPath             string
	reportsPath        string
	artifactsPath      string
	log                logger.Logger

	// Deps stashed to support RestartAll — rebuilding manager + executor
	// requires the original configs, transport factory, blockchain client, and
	// loggers. The dataLayer field above is also preserved across restart so
	// the keyset-recovery handshake can serve the stored recovery data to the
	// fresh executor.
	appType          string
	mgrConfig        *manager.Config
	execConfig       *executor.Config
	excFactory       communication.ConnectionFactory
	blockchainClient blockchain.Client
	mgrLog           logger.Logger
	excLog           logger.Logger
}

// EnableReports marks the manager config so NewTestSuiteCore will provision a
// temporary deanonymization-report directory during construction. Call this
// BEFORE NewTestSuiteCore. The actual path is then available via
// GetReportsPath() on the returned core.
//
// The mock suite does not call this; the fullstack suite calls it so the
// in-process authority service and manager can share the same reports dir.
func EnableReports(mgrConfig *manager.Config) {
	if mgrConfig.DeanonymizationReportPath == "" {
		mgrConfig.DeanonymizationReportPath = "fullstack-reports-placeholder"
	}
}

// NewTestSuiteCore builds the infrastructure shared by both mock and fullstack suites.
// The caller provides the blockchain client implementation (MockClient or real BlockChainClient).
func NewTestSuiteCore(
	t *testing.T,
	appType string,
	mgrConfig *manager.Config,
	execConfig *executor.Config,
	blockchainClient blockchain.Client,
	keySet *executor.EnclaveKeySet,
	recoveryData *common.EnclaveKeySetRecovery,
	mgrLogCfg *logger.Config,
	excLogCfg *logger.Config,
) *TestSuiteCore {
	ctx, cancel := context.WithCancel(context.Background())

	// Assign ephemeral ports for executor, admin, and log server
	// to avoid cross-test collisions.
	executorPort := mustGetFreeTCPPort(t)
	adminPort := mustGetFreeTCPPort(t)
	logServerPort := mustGetFreeTCPPort(t)

	tcpParams := common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: executorPort}
	execConfig.ChannelParams = tcpParams
	execConfig.ChannelType = "tcp"
	mgrConfig.ChannelParams = tcpParams
	mgrConfig.ChannelType = "tcp"
	logServerAddr := common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: logServerPort}
	execConfig.LogChannelParams = logServerAddr
	mgrConfig.LogServerTCPAddress = logServerAddr
	mgrConfig.AdminChannelParams = common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: adminPort}

	// Inject the log-server address into the logger configs so that
	// zeronetwork loggers connect to the correct ephemeral port.
	// Console/file loggers ignore RemoteLogParams harmlessly.
	mgrLogCfg.RemoteLogParams = logServerAddr
	mgrLogCfg.RemoteLogNetwork = "tcp"
	excLogCfg.RemoteLogParams = logServerAddr
	excLogCfg.RemoteLogNetwork = "tcp"
	mgrLog := logger.NewLogger(mgrLogCfg)
	excLog := logger.NewLogger(excLogCfg)

	// Create an executor client (TCP for testing)
	factory := communication.NewTCPConnectionFactory(mgrConfig.ChannelParams.(common.TcpChannelConnectionParams).Url())
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
	dbPath, err := os.MkdirTemp("", "vela-test-db")
	require.NoError(t, err)

	// Create a temporary directory for deploy artifacts and force manager to use it.
	artifactsPath, err := os.MkdirTemp("", "horizen-pes-test-artifacts")
	require.NoError(t, err)
	mgrConfig.ArtifactsPath = artifactsPath

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
		runtime = wasm.NewWasmtimeRuntime(excLog, 0)
	}

	// Create the executor (nil KMS dependencies for Type 0 testing)
	exec, err := executor.NewStatelessExecutor(execConfig, runtime, server, excLog, nil, nil)
	require.NoError(t, err)

	if keySet != nil && recoveryData != nil {
		err := dataLayer.StoreEnclaveKeySetRecovery(ctx, recoveryData)
		require.NoError(t, err)
	}

	// Create event channel (populated by each suite's specific mechanism)
	eventChannel := make(chan interface{}, 100)

	core := &TestSuiteCore{
		t:                t,
		manager:          mgr,
		executor:         exec,
		dataLayer:        dataLayer,
		eventChannel:     eventChannel,
		ctx:              ctx,
		cancel:           cancel,
		dbPath:           dbPath,
		reportsPath:      reportsPath,
		artifactsPath:    artifactsPath,
		log:              mgrLog,
		appType:          appType,
		mgrConfig:        mgrConfig,
		execConfig:       execConfig,
		excFactory:       factory,
		blockchainClient: blockchainClient,
		mgrLog:           mgrLog,
		excLog:           excLog,
	}

	if keySet != nil {
		core.executorCommKey = &keySet.CommunicationKey
		core.executorSigningKey = &keySet.SigningKey
	}

	return core
}

func mustGetFreeTCPPort(t *testing.T) uint32 {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	return uint32(listener.Addr().(*net.TCPAddr).Port)
}

// --- Methods shared by both suites (promoted via embedding) ---

func (s *TestSuiteCore) StartManager() error {
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

func (s *TestSuiteCore) StartExecutor() error {
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

// StopExecutor closes the current executor (its TCP server + runtime). Safe to
// call when executor is already nil; the field is cleared so callers know it
// needs rebuilding before the next Start.
func (s *TestSuiteCore) StopExecutor() error {
	if s.executor == nil {
		return nil
	}
	err := s.executor.Close()
	s.executor = nil
	return err
}

// StopManager stops the current manager's polling loop and admin server. Safe
// to call when manager is already nil.
func (s *TestSuiteCore) StopManager() error {
	if s.manager == nil {
		return nil
	}
	err := s.manager.Stop()
	s.manager = nil
	return err
}

// RestartCore stops both manager and executor and rebuilds them from
// scratch. Used by each suite's public RestartAll() to handle the parts that
// are independent of suite type; the caller must pass a freshly-constructed
// blockchainClient because manager.Stop() closes the one the previous
// manager owned (production-realistic behavior).
//
// The preserved on-disk data — specifically the LevelDB dbPath — carries
// the keyset-recovery blob stored at initial startup. RestartCore re-opens
// the data layer at the same path; the fresh manager serves the stored
// recovery blob during the fresh executor's Type-0 handshake, and the
// executor restores its original keyset. Used by tests that need to prove
// state survives a coupled process restart (both containers dying
// together, as in a host reboot).
//
// Why both at once rather than executor-only: the manager has no
// auto-reconnect logic for the executor channel. A solo executor restart
// in the current design would leave the manager with a dead connection
// forever. When that gap is fixed, a StopExecutor + StartExecutor flow
// without the manager side becomes meaningful as a separate, stronger
// test.
func (s *TestSuiteCore) RestartCore(freshBlockchainClient blockchain.Client) error {
	if err := s.StopManager(); err != nil {
		return fmt.Errorf("stop manager: %w", err)
	}
	if err := s.StopExecutor(); err != nil {
		return fmt.Errorf("stop executor: %w", err)
	}
	// Brief pause so any in-flight sockets settle before we bind the new server.
	time.Sleep(100 * time.Millisecond)

	// Re-open the data layer at the same dbPath. manager.Stop() closed the
	// previous handle; the on-disk LevelDB files (including the
	// keyset-recovery blob) are intact. For the mockdb suite this would
	// allocate a fresh in-memory store and lose the recovery data — not a
	// useful test with that backend, so we don't attempt it.
	if s.mgrConfig.DataLayerType == "mockdb" {
		return fmt.Errorf("RestartCore is not supported with DataLayerType=mockdb (no on-disk persistence)")
	}
	cfg := versioned_leveldb.VersionedLevelDBConfig{
		DBPath:         s.dbPath,
		VersionsToKeep: s.mgrConfig.DataLayerNumOfVersions,
	}
	newDataLayer, err := versioned_leveldb.NewVersionedLevelDBDataLayer(cfg)
	if err != nil {
		return fmt.Errorf("reopen data layer: %w", err)
	}
	s.dataLayer = newDataLayer
	s.blockchainClient = freshBlockchainClient

	// Rebuild executor: fresh comm server (shutdownChan is single-use in
	// Server.Stop, so reuse isn't possible), fresh runtime (Wasmtime closes
	// on executor.Close and isn't reusable), fresh executor over the SAME
	// config and factory so it lands on the same port the manager will dial.
	server := communication.NewServer(s.excFactory, commParams, s.excLog)
	var runtime executor.Runtime
	switch s.appType {
	case "mock-runtime":
		runtime = executor.NewMockRuntime(s.excLog)
	default:
		runtime = wasm.NewWasmtimeRuntime(s.excLog, 0)
	}
	exec, err := executor.NewStatelessExecutor(s.execConfig, runtime, server, s.excLog, nil, nil)
	if err != nil {
		return fmt.Errorf("new executor: %w", err)
	}
	s.executor = exec

	// Rebuild manager: fresh executor client (its handshake state is
	// sync.Once-latched so a reuse can't handshake again), fresh data layer
	// handle, caller-supplied fresh blockchain client.
	executorClient := communication.NewClient(s.excFactory, commParams, s.mgrLog)
	mgr := manager.NewSecureProcessorManager(s.mgrConfig, s.blockchainClient, s.dataLayer, executorClient, nil, s.mgrLog)
	s.manager = mgr

	if err := s.StartExecutor(); err != nil {
		return fmt.Errorf("start executor: %w", err)
	}
	if err := s.StartManager(); err != nil {
		return fmt.Errorf("start manager: %w", err)
	}
	return nil
}

func (s *TestSuiteCore) WaitForAppStateInDB(appID common.ApplicationIdType, timeout time.Duration) (*common.ApplicationState, error) {
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
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %d", appID)

		}
	}
}

// WaitForEvent waits for a specific event to be published for a user.
// If eventSubType is the zero value, any subtype is accepted.
func (s *TestSuiteCore) WaitForEvent(userID ethCommon.Address, eventSubType [32]byte, timeout time.Duration) (*common.Event, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)
	zero := [32]byte{}

	for {
		select {
		case event := <-s.eventChannel:
			if evt, ok := event.(common.Event); ok && evt.UserID == userID && (eventSubType == zero || evt.EventSubType == eventSubType) {
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

// WaitForEventBySubtypes waits for an event for userID whose EventSubType is in validSubtypes.
// Used for privacy-preserving subtype retrieval where the exact subtype is not known in advance.
func (s *TestSuiteCore) WaitForEventBySubtypes(userID ethCommon.Address, validSubtypes [][32]byte, timeout time.Duration) (*common.Event, error) {
	subtypeSet := make(map[[32]byte]struct{}, len(validSubtypes))
	for _, st := range validSubtypes {
		subtypeSet[st] = struct{}{}
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case event := <-s.eventChannel:
			if evt, ok := event.(common.Event); ok && evt.UserID == userID {
				if _, matched := subtypeSet[evt.EventSubType]; matched {
					s.log.Info("TESTING: Received event: %v", evt)
					return &evt, nil
				}
			}
			s.log.Info("TESTING: Received unexpected event: %v", event)
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for event for user %s", userID)
		}
	}
}

// WaitForDeanonymizationReport waits for a deanonymization report to be generated and saved to the filesystem
func (s *TestSuiteCore) WaitForDeanonymizationReport(reportID common.RequestIdType, timeout time.Duration) (*common.DeanonymizationReport, error) {
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
func (s *TestSuiteCore) GetReportsPath() string {
	return s.reportsPath
}

func (s *TestSuiteCore) GetArtifactsPath() string {
	return s.artifactsPath
}

// GetExecutorCommunicationKey returns the executor's communication public key for encryption
func (s *TestSuiteCore) GetExecutorCommunicationKey() (*cryptotypes.PublicKeyP521, error) {
	if s.executorCommKey == nil {
		return nil, fmt.Errorf("executor communication key not initialized")
	}
	return s.executorCommKey.PublicKey(), nil
}

func (s *TestSuiteCore) GetExecutorSigningKey() (*cryptotypes.PublicKeySecp256k1, error) {
	if s.executorSigningKey == nil {
		return nil, fmt.Errorf("executor signing key not initialized")
	}
	return s.executorSigningKey.PublicKey(), nil
}

func (s *TestSuiteCore) GetDataLayer() storage.DataLayer {
	return s.dataLayer
}

func (s *TestSuiteCore) LoadWasmModule(t *testing.T, moduleFilename string) []byte {
	wasmBytes, err := os.ReadFile(moduleFilename)
	require.NoError(t, err, "Failed to read WASM file")
	return wasmBytes
}

// ManagerLogger returns the logger the manager writes to, so callers can route
// their own instrumentation into the same sink.
//
// Only valid once NewTestSuiteCore has returned: it is what injects the ephemeral
// log-server address into the manager's logger config (RemoteLogParams), which a
// "zeronetwork" logger requires. Building a logger from that config any earlier
// panics — use this instead of constructing a second one.
func (s *TestSuiteCore) ManagerLogger() logger.Logger {
	return s.mgrLog
}

// SetEventChannel replaces the event channel. Used by the fullstack suite to
// inject its own channel connected to the event-broadcasting wrapper.
func (s *TestSuiteCore) SetEventChannel(ch chan interface{}) {
	s.eventChannel = ch
}

// Context returns the suite's root context.
func (s *TestSuiteCore) Context() context.Context {
	return s.ctx
}

// CleanupCore performs cleanup common to both suites.
// Each suite's Cleanup() calls this plus any suite-specific teardown.
func (s *TestSuiteCore) CleanupCore() {
	s.cancel()

	if s.manager != nil {
		s.manager.Stop()
	}

	if s.executor != nil {
		s.executor.Close()
	}

	// Remove the temporary database directory
	if s.dbPath != "" {
		os.RemoveAll(s.dbPath)
	}

	// Remove the temporary reports directory
	if s.reportsPath != "" {
		os.RemoveAll(s.reportsPath)
	}
	if s.artifactsPath != "" {
		os.RemoveAll(s.artifactsPath)
	}
}

func PrettyPrintJSON(data interface{}) string {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "Invalid data: could not marshal into a json"
	}
	return string(prettyJSON)
}
