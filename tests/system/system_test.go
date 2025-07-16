package system

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
)

func TestMockAppFullSystemFlow(t *testing.T) {
	suite := NewSystemTestSuite(t)

	const appId = "payment-app"
	const user1 = "user1"
	const user2 = "user2"
	const auditor = "auditor"

	cryptoHelper := NewCryptoHelper()

	t.Log("Step 0: Setup user keys for encryption/decryption")

	// Generate user and auditor keys
	user1Key, err := cryptoHelper.GenerateUserKey(user1)
	require.NoError(t, err)
	user2Key, err := cryptoHelper.GenerateUserKey(user2)
	require.NoError(t, err)
	auditorKey, err := cryptoHelper.GenerateUserKey(auditor)
	require.NoError(t, err)

	// Register keys in the system
	err = suite.AddUserKeys(user1, user1Key.PublicKey().Bytes())
	require.NoError(t, err)
	err = suite.AddUserKeys(user2, user2Key.PublicKey().Bytes())
	require.NoError(t, err)
	err = suite.AddUserKeys(auditor, auditorKey.PublicKey().Bytes())
	require.NoError(t, err)

	t.Log("Step 1: Starting system components and deploying app")

	err = suite.StartExecutor()
	require.NoError(t, err)
	err = suite.StartManager()
	require.NoError(t, err)

	// Get executor's communication key for encryption, for now get from the test suite
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appId,
		RequestID:     "deploy-1",
		Payload:       []byte("payment-app-wasm-bytecode"),
		Sender:        user1,
		Timestamp:     time.Now().Unix(),
	}
	err = suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// Wait for app to be deployed
	appState, err := suite.WaitForAppStateInDB(appId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	appState, err = suite.WaitForAppStateInBlockchain(appId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	err = suite.AssertRequestCompleted("deploy-1", 10*time.Second)
	require.NoError(t, err)

	t.Log("Step 2: Sending deposit request")

	depositReq, err := cryptoHelper.CreateDepositRequest(
		appId,
		"deposit-1",
		user1,
		2000000000000000000, // 2 ETH
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(depositReq)
	require.NoError(t, err)

	// Wait for deposit to be processed
	err = suite.AssertRequestCompleted("deposit-1", 10*time.Second)
	require.NoError(t, err)

	// Wait for deposit event
	depositEvent, err := suite.WaitForEvent(user1, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, depositEvent)

	// Decrypt and verify deposit event
	decryptedDepositData, err := cryptoHelper.DecryptEvent(user1, depositEvent, executorPubKey)
	require.NoError(t, err)

	var depositEventData map[string]interface{}
	err = json.Unmarshal(decryptedDepositData, &depositEventData)
	require.NoError(t, err)
	require.Equal(t, "deposit", depositEventData["type"])
	require.Equal(t, float64(2000000000000000000), depositEventData["amount"])

	t.Log("Step 3: Sending transfer request")

	transferReq, err := cryptoHelper.CreateTransferRequest(
		appId,
		"transfer-1",
		user1,
		user2,
		500000000000000000, // 0.5 ETH
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(transferReq)
	require.NoError(t, err)

	// Wait for transfer to be processed
	err = suite.AssertRequestCompleted("transfer-1", 10*time.Second)
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

	var senderEventData map[string]interface{}
	err = json.Unmarshal(decryptedSenderData, &senderEventData)
	require.NoError(t, err)
	require.Equal(t, "transfer_sent", senderEventData["type"])
	require.Equal(t, user2, senderEventData["to"])
	require.Equal(t, float64(500000000000000000), senderEventData["amount"])

	// Decrypt and verify recipient event
	decryptedRecipientData, err := cryptoHelper.DecryptEvent(user2, recipientEvent, executorPubKey)
	require.NoError(t, err)

	var recipientEventData map[string]interface{}
	err = json.Unmarshal(decryptedRecipientData, &recipientEventData)
	require.NoError(t, err)
	require.Equal(t, "transfer_received", recipientEventData["type"])
	require.Equal(t, user1, recipientEventData["from"])
	require.Equal(t, float64(500000000000000000), recipientEventData["amount"])

	t.Log("Step 4: Sending deanonymization request as auditor")

	deanonReq, err := cryptoHelper.CreateDeanonymizationRequest(
		appId,
		"deanon-1",
		auditor,
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(deanonReq)
	require.NoError(t, err)

	// Wait for deanonymization request to be processed
	err = suite.AssertRequestCompleted("deanon-1", 10*time.Second)
	require.NoError(t, err)

	// Wait for deanonymization report
	deanonReport, err := suite.WaitForDeanonymizationReport("deanon-1", 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, deanonReport)

	// Decrypt and verify deanonymization report
	decryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditor, deanonReport, executorPubKey)
	require.NoError(t, err)

	var reportData map[string]interface{}
	err = json.Unmarshal(decryptedReport, &reportData)
	require.NoError(t, err)
	require.Equal(t, appId, reportData["applicationId"])
	require.Equal(t, "deanon-1", reportData["requestId"])
	require.Contains(t, reportData, "accounts")

	// Verify account information in the report
	accounts, ok := reportData["accounts"].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, accounts, user1)
	require.Equal(t, 1500000000000000000.0, accounts[user1].(map[string]interface{})["balance"])
	require.Contains(t, accounts, user2)
	require.Equal(t, 500000000000000000.0, accounts[user2].(map[string]interface{})["balance"])

	// Step 5: As another user, send withdrawal request
	t.Log("Step 5: Sending withdrawal request as user2")

	withdrawalReq, err := cryptoHelper.CreateWithdrawalRequest(
		appId,
		"withdraw-1",
		user2,
		"0x1234567890123456789012345678901234567890",
		500000000000000000, // 0.5 ETH
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(withdrawalReq)
	require.NoError(t, err)

	// Wait for withdrawal to be processed
	err = suite.AssertRequestCompleted("withdraw-1", 10*time.Second)
	require.NoError(t, err)

	// Wait for withdrawal event
	withdrawalEvent, err := suite.WaitForEvent(user2, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, withdrawalEvent)

	// Decrypt and verify withdrawal event
	decryptedWithdrawalData, err := cryptoHelper.DecryptEvent(user2, withdrawalEvent, executorPubKey)
	require.NoError(t, err)

	var withdrawalEventData map[string]interface{}
	err = json.Unmarshal(decryptedWithdrawalData, &withdrawalEventData)
	require.NoError(t, err)
	require.Equal(t, "withdrawal", withdrawalEventData["type"])
	require.Equal(t, "0x1234567890123456789012345678901234567890", withdrawalEventData["to"])
	require.Equal(t, float64(500000000000000000), withdrawalEventData["amount"])

	// Wait for actual withdrawal to be recorded
	withdrawal, err := suite.WaitForWithdrawal(appId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, withdrawal)
	require.Equal(t, "0x1234567890123456789012345678901234567890", withdrawal.DestinationAddress)
	require.Equal(t, "500000000000000000", withdrawal.Amount)

	t.Log("system test completed successfully!")

	suite.Cleanup()
}
