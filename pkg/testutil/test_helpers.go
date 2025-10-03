package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/horizen-pes/pkg/wasm"
	appCommon "github.com/horizen-pes/pkg/wasm/common"
	"github.com/stretchr/testify/require"
)

type SystemTestSuite struct {
	t                  *testing.T
	manager            manager.Manager
	executor           executor.Executor
	blockchainClient   *blockchain.MockClient
	dataLayer          *mockdb.MockDataLayer
	eventChannel       chan interface{}
	ctx                context.Context
	cancel             context.CancelFunc
	executorCommKey    *cryptotypes.PrivateKeyP521      // Executor's communication key for testing
	executorSigningKey *cryptotypes.PrivateKeySecp256k1 // Executor's signing key for testing
}

func NewSystemTestSuite(t *testing.T, appType string) *SystemTestSuite {
	ctx, cancel := context.WithCancel(context.Background())

	// Create mock components
	blockchainClient := blockchain.NewMockClient()
	dataLayer := mockdb.NewMockDataLayer()

	// Create an executor client (TCP for testing)
	factory := communication.NewTCPConnectionFactory("localhost:8080")
	executorClient := communication.NewClient(factory)

	// Create manager
	config := manager.ReadConfig()
	mgr := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient)

	// Create executor
	execConfig := executor.DefaultConfig() // just to generate keys

	server := communication.NewServer(factory)
	var runtime executor.Runtime
	switch appType {
	case "mock-runtime":
		t.Log("mock app type: ", appType)
		runtime = executor.NewMockRuntime()
	default:
		t.Log("wasm app type: ", appType)
		runtime = wasm.NewWasmtimeRuntime()
	}
	exec, err := executor.NewStatelessExecutor(execConfig, runtime, server)
	require.NoError(t, err)

	// Create event channel
	eventChannel := make(chan interface{}, 100)
	blockchainClient.SubscribeToEvents(ctx, eventChannel)

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
		executorSigningKey: execConfig.SignatureKey,     // Store the executor's signing key
	}
}

func (s *SystemTestSuite) StartManager() error {
	go func() {
		if err := s.manager.Start(s.ctx); err != nil {
			s.t.Errorf("Manager failed: %v", err)
		}
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

func (s *SystemTestSuite) SubmitRequest(req *common.Request) error {
	return s.blockchainClient.SubmitRequest(s.ctx, req)
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
			state, err := s.blockchainClient.GetApplicationState(s.ctx, appID)
			if err == nil {
				return state, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %s in blockchain", appID)

		}
	}
}

func (s *SystemTestSuite) AssertRequestCompleted(requestID string, timeout time.Duration) error {
	return s.blockchainClient.WaitForRequestCompletion(requestID, timeout)
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
			// Check if deanonymization report exists in blockchain
			report, err := s.blockchainClient.GetDeanonymizationReport(s.ctx, reportID)
			if err == nil && report != nil {
				return report, nil
			}
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
			// Check if withdrawal exists in blockchain
			withdrawals, err := s.blockchainClient.GetWithdrawals(s.ctx, appID)
			if err == nil && withdrawals != nil && len(*withdrawals) > 0 {
				return &(*withdrawals)[0], nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for withdrawal for app %s", appID)
		}
	}
}

func (s *SystemTestSuite) GetRequestUpdatePayload(reqId string) (*common.UpdatePayload, error) {
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

// GetExecutorSigningKey returns the executor's signing public key for encryption
func (s *SystemTestSuite) GetExecutorSigningKey() (*cryptotypes.PublicKeySecp256k1, error) {
	if s.executorSigningKey == nil {
		return nil, fmt.Errorf("executor signing key not initialized")
	}
	return s.executorSigningKey.PublicKey(), nil
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

	return nil
}

func ExecTestAppFullSystemFlow(t *testing.T, suite *SystemTestSuite, bytecode []byte) {
	const appId = "1"
	user1 := fmt.Sprintf("0xadd%037x", 1)
	user2 := fmt.Sprintf("0xadd%037x", 2)
	const auditor = "auditor"

	cryptoHelper := NewCryptoHelper()

	t.Log("Step 0: Starting system components and deploying app")

	var err error

	err = suite.StartExecutor()
	require.NoError(t, err)
	err = suite.StartManager()
	require.NoError(t, err)

	// Get executor's communication key for encryption, for now get from the test suite
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// Get executor's signing key for signature verification
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)

	RequestID := "2133"
	// Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appId,
		RequestID:     RequestID,
		Payload:       bytecode,
		Sender:        user1,
		Timestamp:     time.Now().Unix(),
	}
	err = suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// Wait for app to be deployed
	appState, err := suite.WaitForAppStateInDB(appId, 100*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	appState, err = suite.WaitForAppStateInBlockchain(appId, 100*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	err = suite.AssertRequestCompleted(RequestID, 100*time.Second)
	require.NoError(t, err)

	// Verify updatePayload signature
	payload, err := suite.GetRequestUpdatePayload(RequestID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

	t.Log("Step 1: Setup user keys for encryption/decryption")

	// Generate user and auditor keys
	user1Key, err := cryptoHelper.GenerateUserKey(user1)
	require.NoError(t, err)
	user2Key, err := cryptoHelper.GenerateUserKey(user2)
	require.NoError(t, err)
	auditorKey, err := cryptoHelper.GenerateUserKey(auditor)
	require.NoError(t, err)

	//register key 1
	RequestID = "2130"
	associateKey1Req, err := cryptoHelper.CreateAssociateKeyRequest(appId, RequestID, user1, user1Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, 100*time.Second)
	require.NoError(t, err)

	//register key 2
	RequestID = "2131"
	associateKey2Req, err := cryptoHelper.CreateAssociateKeyRequest(appId, RequestID, user2, user2Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey2Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, 100*time.Second)
	require.NoError(t, err)

	//register key 3
	RequestID = "2132"
	associateKey3Req, err := cryptoHelper.CreateAssociateKeyRequest(appId, RequestID, auditor, auditorKey.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey3Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, 100*time.Second)
	require.NoError(t, err)

	t.Log("Step 2: Sending deposit request")

	RequestID = "2134"
	depositAmount := uint64(2000000000000000000)
	depositReq, err := cryptoHelper.CreateDepositRequest(
		appId,
		RequestID,
		user1,
		depositAmount,
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(depositReq)
	require.NoError(t, err)

	// Wait for deposit to be processed
	err = suite.AssertRequestCompleted(RequestID, 100*time.Second)
	require.NoError(t, err)

	// Wait for deposit event
	depositEvent, err := suite.WaitForEvent(user1, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, depositEvent)

	// Decrypt and verify deposit event
	decryptedDepositData, err := cryptoHelper.DecryptEvent(user1, depositEvent, executorPubKey)
	require.NoError(t, err)

	var depositEventData appCommon.DepositEvent
	err = json.Unmarshal(decryptedDepositData, &depositEventData)
	require.NoError(t, err)
	require.Equal(t, "deposit", depositEventData.Type)
	require.Equal(t, depositAmount, depositEventData.Amount)

	// Verify updatePayload signature
	payload, err = suite.GetRequestUpdatePayload(RequestID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

	t.Log("Step 3: Sending transfer request")

	RequestID = "2135"
	sentAmount := uint64(500000000000000000) // 0.5 ETH
	transferReq, err := cryptoHelper.CreateTransferRequest(
		appId,
		RequestID,
		user1,
		user2,
		sentAmount,
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(transferReq)
	require.NoError(t, err)

	// Wait for transfer to be processed
	err = suite.AssertRequestCompleted(RequestID, 10*time.Second)
	require.NoError(t, err)

	// Wait for transfer events for both users
	senderEvent, err := suite.WaitForEvent(user1, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, senderEvent)

	recipientEvent, err := suite.WaitForEvent(user2, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, recipientEvent)

	// Decrypt and verify sender event
	decryptedSenderData, err := cryptoHelper.DecryptEvent(user1, senderEvent, executorPubKey)
	require.NoError(t, err)

	var senderEventData appCommon.SenderEvent
	err = json.Unmarshal(decryptedSenderData, &senderEventData)
	require.NoError(t, err)
	require.Equal(t, "transfer_sent", senderEventData.Type)
	require.Equal(t, user2, senderEventData.To)
	require.Equal(t, sentAmount, senderEventData.Amount)

	// Decrypt and verify recipient event
	decryptedRecipientData, err := cryptoHelper.DecryptEvent(user2, recipientEvent, executorPubKey)
	require.NoError(t, err)

	var recipientEventData appCommon.RecipientEvent
	err = json.Unmarshal(decryptedRecipientData, &recipientEventData)
	require.NoError(t, err)
	require.Equal(t, "transfer_received", recipientEventData.Type)
	require.Equal(t, user1, recipientEventData.From)
	require.Equal(t, sentAmount, recipientEventData.Amount)

	// Verify updatePayload signature
	payload, err = suite.GetRequestUpdatePayload(RequestID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

	t.Log("Step 4: Sending deanonymization request as auditor")

	RequestID = "2136"

	deanonReq, err := cryptoHelper.CreateDeanonymizationRequest(
		appId,
		RequestID,
		auditor,
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(deanonReq)
	require.NoError(t, err)

	// Wait for deanonymization request to be processed
	err = suite.AssertRequestCompleted(RequestID, 10*time.Second)
	require.NoError(t, err)

	// Wait for deanonymization report
	deanonReport, err := suite.WaitForDeanonymizationReport(RequestID, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, deanonReport)

	// Decrypt and verify deanonymization report
	decryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditor, deanonReport, executorPubKey)
	require.NoError(t, err)

	var reportData appCommon.UnencryptedDeanonymizationReportData
	err = json.Unmarshal(decryptedReport, &reportData)
	require.NoError(t, err)
	require.Equal(t, appId, reportData.ApplicationID)
	require.Equal(t, RequestID, reportData.RequestID)

	// Verify account information in the report
	accounts := reportData.Accounts
	require.Contains(t, accounts, user1)
	require.Equal(t, uint64(1500000000000000000), accounts[user1].Balance)

	require.Contains(t, accounts, user2)
	require.Equal(t, uint64(500000000000000000), accounts[user2].Balance)

	// Deanon report does not contain signature for now, possibly add later

	// Step 5: As another user, send withdrawal request
	t.Log("Step 5: Sending withdrawal request as user2")

	RequestID = "2137"
	withdrawAmount := uint64(500000000000000000) // 0.5 ETH
	withdrawalReq, err := cryptoHelper.CreateWithdrawalRequest(
		appId,
		RequestID,
		user2,
		"0x1234567890123456789012345678901234567890",
		withdrawAmount,
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(withdrawalReq)
	require.NoError(t, err)

	// Wait for withdrawal to be processed
	err = suite.AssertRequestCompleted(RequestID, 10*time.Second)
	require.NoError(t, err)

	// Wait for withdrawal event
	withdrawalEvent, err := suite.WaitForEvent(user2, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, withdrawalEvent)

	// Decrypt and verify withdrawal event
	decryptedWithdrawalData, err := cryptoHelper.DecryptEvent(user2, withdrawalEvent, executorPubKey)
	require.NoError(t, err)

	var withdrawalEventData appCommon.WithdrawalEvent
	err = json.Unmarshal(decryptedWithdrawalData, &withdrawalEventData)
	require.NoError(t, err)
	require.Equal(t, "withdrawal", withdrawalEventData.Type)
	require.Equal(t, "0x1234567890123456789012345678901234567890", withdrawalEventData.To)
	require.Equal(t, withdrawAmount, withdrawalEventData.Amount)

	// Wait for actual withdrawal to be recorded
	withdrawal, err := suite.WaitForWithdrawal(appId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, withdrawal)
	require.Equal(t, "0x1234567890123456789012345678901234567890", withdrawal.DestinationAddress)
	require.Equal(t, withdrawAmount, withdrawal.Amount)

	// Verify updatePayload signature
	payload, err = suite.GetRequestUpdatePayload(RequestID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

	t.Log("system test completed successfully!")

}
