package system

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/horizen-pes/pkg/wasm"
	"github.com/stretchr/testify/require"
)



type SystemTestSuite struct {
	t                  *testing.T
	manager            manager.Manager
	executor           executor.Executor
	blockchainClient   *blockchain.BlockChainClient
	dataLayer          *mockdb.MockDataLayer
	eventChannel       chan interface{}
	ctx                context.Context
	cancel             context.CancelFunc
	executorCommKey    *cryptotypes.PrivateKeyP521      // Executor's communication key for testing
	executorSigningKey *cryptotypes.PrivateKeySecp256k1 // Executor's signing key for testing
	simNode			   *blockchain.SimTestHelper
}

type MockBlockchainClient struct {
	*blockchain.BlockChainClient
}

func (c *MockBlockchainClient) Connect(ctx context.Context) error {
	fmt.Println("MockBlockchainClient: Simulating connection to blockchain")
	return nil
}


func NewSystemTestSuite(t *testing.T, appType string) *SystemTestSuite {
	ctx, cancel := context.WithCancel(context.Background())

	// Create the simulated blockchain node, with all the needed contracts
	execConfig := executor.DefaultConfig()
	execAddress := crypto.PubkeyToAddress(execConfig.SignatureKey.PrivateKey.PublicKey)
	sim := blockchain.NewSimTestHelper(t, true, false, &execAddress)
	// Create mock components
	//blockchainClient := blockchain.NewMockClient()
	blockchainClient := sim.SetupNewBlockChainClient()
	dataLayer := mockdb.NewMockDataLayer()

	// Create an executor client (TCP for testing)
	factory := communication.NewTCPConnectionFactory("localhost:8080")
	executorClient := communication.NewClient(factory)

	// Create manager
	config := manager.DefaultConfig()
	config.ExecutorConnectionType = "tcp"
	config.ExecutorConnectionParams = map[string]string{"url": "http://localhost:8080"}
	mgr := manager.NewSecureProcessorManager(config, &MockBlockchainClient{blockchainClient}, dataLayer, executorClient)

	// Create executor

	execConfig.ServerType = "tcp"
	execConfig.ServerAddr = "localhost:8080"

	server := communication.NewServer(factory)
	var runtime executor.Runtime
	switch appType {
	case "wasmtime-payment":
		runtime = wasm.NewWasmtimeRuntime()
	case "mock-runtime":
		runtime = executor.NewMockRuntime()
	default:
		t.Fatalf("Unknown app type: %s", appType)
	}
	exec := executor.NewStatelessExecutor(execConfig, runtime, server)

	// Create event channel
	eventChannel := make(chan interface{}, 100)
//TODO	blockchainClient.SubscribeToEvents(ctx, eventChannel)

	return &SystemTestSuite{
		t:                  t,
		manager:            mgr,
		executor:           exec,
		blockchainClient:   blockchainClient,
		dataLayer:          dataLayer,
		eventChannel:       eventChannel,
		ctx:                ctx,
		cancel:             cancel,
		executorCommKey:    execConfig.CommunicationKey, // Store the executor's communication key
		executorSigningKey: execConfig.SignatureKey,     // Store the executor's communication key
		simNode: 		  sim,
	}
}

func (s *SystemTestSuite) StartManager() error {
	go func() {
		if err := s.manager.Start(s.ctx); err != nil {
			fmt.Println("Manager failed: ", err)
			s.t.Errorf("Manager failed: %v", err)
		}
		fmt.Println("Manager started")
	}()

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) StartExecutor() error {
	go func() {
		if err := s.executor.Start(s.ctx); err != nil {
			s.t.Errorf("Executor failed: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) AddUserKeys(sender *bind.TransactOpts, publicKey []byte) error {
	// Register in a blockchain client
	tx, err := s.simNode.RegisterUserKey(sender, publicKey)
	if err != nil {	
		return err
	}

	s.simNode.WaitMined(tx)	
	
	// Register in data layer
	return s.dataLayer.StoreUserKey(s.ctx, sender.From.Hex(), publicKey)
}

func (s *SystemTestSuite) SubmitRequest(req *common.Request) (string, error) {

	return s.SubmitRequestFromUser(req, s.simNode.Submitter)
}

func (s *SystemTestSuite) SubmitRequestFromUser(req *common.Request, sender *bind.TransactOpts) (string, error) {
	appId, ok := common.StringToBigInt(req.ApplicationID)
	require.True(s.t, ok, "invalid application ID")
	tx := s.simNode.SubmitRequestFromUser(appId, req.RequestType,req.Payload, new(big.Int).SetUint64(req.Value), sender)

	s.simNode.WaitMined(tx)
	
	event := s.simNode.GetRequestSubmittedEvent(tx)


	return event.RequestId.String(), nil
}


func (s *SystemTestSuite) WaitForAppStateInDB(appID string, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			state, err := s.dataLayer.GetApplicationState(s.ctx, appID)
			if err == nil {
				return state, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %s", appID)

		}
	}
}

func (s *SystemTestSuite) WaitForAppStateInBlockchain(appID string, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
		state := &common.ApplicationState{ApplicationID: appID, StateRoot: s.simNode.GetStateRoot()}
		return state, nil
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %s in blockchain", appID)

		}
	}
}

func (s *SystemTestSuite) AssertRequestCompleted(requestID string, timeout time.Duration) (blockchain.RequestStatus, error) {
	return s.simNode.WaitForRequestCompletion(requestID, timeout)
}

// WaitForEvent waits for a specific event to be published for a user
func (s *SystemTestSuite) WaitForEvent(userID string, timeout time.Duration) (*common.Event, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case event := <-s.eventChannel:
			log.Printf("TESTING: Received event: %+v", event)
			if evt, ok := event.(common.Event); ok && evt.UserID == userID {
				return &evt, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for event for user %s", userID)
		}
	}
}

// WaitForDeanonymizationReport waits for a deanonymization report to be generated
func (s *SystemTestSuite) WaitForDeanonymizationReport(reportID string, timeout time.Duration) (*common.DeanonymizationReport, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			//TODO
			// // Check if deanonymization report exists in blockchain
			// report, err := s.blockchainClient.GetDeanonymizationReport(s.ctx, reportID)
			// if err == nil && report != nil {
			// 	return report, nil
			// }
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for deanonymization report %s", reportID)
		}
	}
}

// WaitForWithdrawal waits for a withdrawal to be processed
func (s *SystemTestSuite) WaitForWithdrawal(appID string, timeout time.Duration) (*common.Withdrawal, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			//TODO
			// // Check if withdrawal exists in blockchain
			// withdrawals, err := s.blockchainClient.GetWithdrawals(s.ctx, appID)
			// if err == nil && withdrawals != nil && len(*withdrawals) > 0 {
			// 	return &(*withdrawals)[0], nil
			// }
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for withdrawal for app %s", appID)
		}
	}
}


// GetExecutorCommunicationKey returns the executor's communication public key for encryption
func (s *SystemTestSuite) GetExecutorCommunicationKey() (*cryptotypes.PublicKeyP521, error) {
	if s.executorCommKey == nil {
		return nil, fmt.Errorf("executor communication key not initialized")
	}
	return s.executorCommKey.PublicKey(), nil
}

// GetExecutorSigningKey returns the executor's signing public key for encryption
func (s *SystemTestSuite) GetExecutorSigningKey() (*cryptotypes.PublicKeySecp256k1, error) {
	if s.executorSigningKey == nil {
		return nil, fmt.Errorf("executor signing key not initialized")
	}
	return s.executorSigningKey.PublicKey(), nil
}

func (s *SystemTestSuite) LoadWasmModule(t *testing.T, moduleFilename string) []byte {
	// Load the compiled WASM module
	wasmPath := filepath.Join("wasm", moduleFilename)
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")
	return wasmBytes
}

func (s *SystemTestSuite) Cleanup() error {
	fmt.Println("Cleaning up the test suite")
	s.cancel()

	if s.manager != nil {
		s.manager.Stop()
	}

	if s.executor != nil {
		s.executor.Close()
	}

	if s.simNode != nil {
		s.simNode.Close()	
	}

	return nil
}

func (s *SystemTestSuite) TransferFunds(sender *bind.TransactOpts, toAddress ethCommon.Address, value *big.Int) {
	
	tx := s.simNode.TransferFunds(sender, toAddress, value)

	s.simNode.WaitMined(tx)	

	receipt, err := s.simNode.GetTxReceipt( tx)
	require.NoError(s.t, err, "error getting transaction receipt")
	require.Equal(s.t, uint64(1), receipt.Status, "Transaction failed")

}
