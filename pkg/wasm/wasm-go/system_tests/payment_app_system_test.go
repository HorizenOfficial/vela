package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	wasmCommon "github.com/horizen-pes/pkg/wasm/common"

	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestWasmtimePaymentAppFullSystemFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	suite := NewSystemTestSuite(t, "wasmtime-payment")
	wasmBytecode := suite.LoadWasmModule(t, "payment_app.wasm")
	testPaymentAppFullSystemFlow(t, suite, wasmBytecode)
}

func TestMockRuntimePaymentAppFullSystemFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	suite := NewSystemTestSuite(t, "mock-runtime")
	wasmBytecode := []byte("mock-runtime-payment-app-bytecode")
	testPaymentAppFullSystemFlow(t, suite, wasmBytecode)
}

func testPaymentAppFullSystemFlow(t *testing.T, suite *SystemTestSuite, bytecode []byte) {
	const appId = "payment-app"
	user1 := fmt.Sprintf("0xadd%037x", 1)
	user2 := fmt.Sprintf("0xadd%037x", 2)
	const auditor = "auditor"

	cryptoHelper := NewCryptoHelper()

	// keys
	user1Key, err := cryptoHelper.GenerateUserKey(user1)
	require.NoError(t, err)
	user2Key, err := cryptoHelper.GenerateUserKey(user2)
	require.NoError(t, err)
	auditorKey, err := cryptoHelper.GenerateUserKey(auditor)
	require.NoError(t, err)

	require.NoError(t, suite.AddUserKeys(user1, user1Key.PublicKey().Bytes()))
	require.NoError(t, suite.AddUserKeys(user2, user2Key.PublicKey().Bytes()))
	require.NoError(t, suite.AddUserKeys(auditor, auditorKey.PublicKey().Bytes()))

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)

	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appId,
		RequestID:     "deploy-1",
		Payload:       bytecode,
		Sender:        user1,
		Timestamp:     time.Now().Unix(),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))

	_, err = suite.WaitForAppStateInDB(appId, 100*time.Second)
	require.NoError(t, err)
	_, err = suite.WaitForAppStateInBlockchain(appId, 100*time.Second)
	require.NoError(t, err)
	require.NoError(t, suite.AssertRequestCompleted("deploy-1", 100*time.Second))
	payload, err := suite.GetRequestUpdatePayload("deploy-1")
	require.NoError(t, err)
	require.NoError(t, cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey))

	depositAmount := uint64(2000000000000000000)
	depositReq, err := cryptoHelper.CreateDepositRequest(appId, "deposit-1", user1, depositAmount, executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(depositReq))
	require.NoError(t, suite.AssertRequestCompleted("deposit-1", 10*time.Second))
	depositEvent, err := suite.WaitForEvent(user1, 10*time.Second)
	require.NoError(t, err)
	decryptedDepositData, err := cryptoHelper.DecryptEvent(user1, depositEvent, executorPubKey)
	require.NoError(t, err)
	var depositEventData wasmCommon.DepositEvent
	require.NoError(t, json.Unmarshal(decryptedDepositData, &depositEventData))
	require.Equal(t, "deposit", depositEventData.Type)
	require.Equal(t, depositAmount, depositEventData.Amount)
	payload, err = suite.GetRequestUpdatePayload("deposit-1")
	require.NoError(t, err)
	require.NoError(t, cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey))

	sentAmount := uint64(500000000000000000)
	transferReq, err := cryptoHelper.CreateTransferRequest(appId, "transfer-1", user1, user2, sentAmount, executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(transferReq))
	require.NoError(t, suite.AssertRequestCompleted("transfer-1", 10*time.Second))
	senderEvent, err := suite.WaitForEvent(user1, 10*time.Second)
	require.NoError(t, err)
	recipientEvent, err := suite.WaitForEvent(user2, 10*time.Second)
	require.NoError(t, err)
	decryptedSenderData, err := cryptoHelper.DecryptEvent(user1, senderEvent, executorPubKey)
	require.NoError(t, err)
	var senderEventData wasmCommon.SenderEvent
	require.NoError(t, json.Unmarshal(decryptedSenderData, &senderEventData))
	require.Equal(t, "transfer_sent", senderEventData.Type)
	require.Equal(t, user2, senderEventData.To)
	require.Equal(t, sentAmount, senderEventData.Amount)
	decryptedRecipientData, err := cryptoHelper.DecryptEvent(user2, recipientEvent, executorPubKey)
	require.NoError(t, err)
	var recipientEventData wasmCommon.RecipientEvent
	require.NoError(t, json.Unmarshal(decryptedRecipientData, &recipientEventData))
	require.Equal(t, "transfer_received", recipientEventData.Type)
	require.Equal(t, user1, recipientEventData.From)
	require.Equal(t, sentAmount, recipientEventData.Amount)
	payload, err = suite.GetRequestUpdatePayload("transfer-1")
	require.NoError(t, err)
	require.NoError(t, cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey))

	deanonReq, err := cryptoHelper.CreateDeanonymizationRequest(appId, "deanon-1", auditor, executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(deanonReq))
	require.NoError(t, suite.AssertRequestCompleted("deanon-1", 10*time.Second))
	deanonReport, err := suite.WaitForDeanonymizationReport("deanon-1", 10*time.Second)
	require.NoError(t, err)
	decryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditor, deanonReport, executorPubKey)
	require.NoError(t, err)
	// Local struct to avoid depending on exact common type name
	type reportStruct struct {
		ApplicationID string                              `json:"applicationId"`
		RequestID     string                              `json:"requestId"`
		Accounts      map[string]*wasmCommon.AccountState `json:"accounts"`
		Nonce         uint64                              `json:"nonce"`
	}
	var reportData reportStruct
	require.NoError(t, json.Unmarshal(decryptedReport, &reportData))
	require.Equal(t, appId, reportData.ApplicationID)
	require.Equal(t, "deanon-1", reportData.RequestID)
	accounts := reportData.Accounts
	require.Contains(t, accounts, user1)
	require.Equal(t, uint64(1500000000000000000), accounts[user1].Balance)
	require.Contains(t, accounts, user2)
	require.Equal(t, uint64(500000000000000000), accounts[user2].Balance)

	withdrawAmount := uint64(500000000000000000)
	withdrawalReq, err := cryptoHelper.CreateWithdrawalRequest(appId, "withdraw-1", user2, "0x1234567890123456789012345678901234567890", withdrawAmount, executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(withdrawalReq))
	require.NoError(t, suite.AssertRequestCompleted("withdraw-1", 10*time.Second))
	withdrawalEvent, err := suite.WaitForEvent(user2, 10*time.Second)
	require.NoError(t, err)
	decryptedWithdrawalData, err := cryptoHelper.DecryptEvent(user2, withdrawalEvent, executorPubKey)
	require.NoError(t, err)
	var withdrawalEventData wasmCommon.WithdrawalEvent
	require.NoError(t, json.Unmarshal(decryptedWithdrawalData, &withdrawalEventData))
	require.Equal(t, "withdrawal", withdrawalEventData.Type)
	require.Equal(t, "0x1234567890123456789012345678901234567890", withdrawalEventData.To)
	require.Equal(t, withdrawAmount, withdrawalEventData.Amount)
	withdrawal, err := suite.WaitForWithdrawal(appId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, withdrawal)
	require.Equal(t, "0x1234567890123456789012345678901234567890", withdrawal.DestinationAddress)
	require.Equal(t, withdrawAmount, withdrawal.Amount)
	payload, err = suite.GetRequestUpdatePayload("withdraw-1")
	require.NoError(t, err)
	require.NoError(t, cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey))

	suite.Cleanup()
}
