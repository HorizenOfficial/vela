package system

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
	appCommon "github.com/horizen-pes/pkg/wasm/common"
)

func TestWasmtimePaymentAppFullSystemFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	suite := NewSystemTestSuite(t, "wasmtime-payment")
	defer suite.Cleanup()
	// Load wasm bytecode for the payment app
	wasmBytecode := suite.LoadWasmModule(t, "payment_app.wasm")

	testPaymentAppFullSystemFlow(t, suite, wasmBytecode)
}

func TestMockRuntimePaymentAppFullSystemFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	suite := NewSystemTestSuite(t, "mock-runtime")
	defer suite.Cleanup()
	// Load wasm bytecode for the payment app
	wasmBytecode := []byte("mock-runtime-payment-app-bytecode")

	testPaymentAppFullSystemFlow(t, suite, wasmBytecode)
}

func testPaymentAppFullSystemFlow(t *testing.T, suite *SystemTestSuite, bytecode []byte) {
	const appId = "1"
	user1 := fmt.Sprintf("0xadd%037x", 1)
	user2 := fmt.Sprintf("0xadd%037x", 2)
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

	t.Log("Step 2: Sending deposit request")

	RequestID = "87654"
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


	RequestID = "4367"
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

	RequestID = "76333"

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

	RequestID =  "1654332"
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
