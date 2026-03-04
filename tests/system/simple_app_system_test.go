package system

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	commontestutil "github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/testutil"
)

// host-side event types for test validation (app-specific, not framework types)
type hostDepositEvent struct {
	Type    string      `json:"type"`
	Amount  *common.Big `json:"amount"`
	Balance *common.Big `json:"balance"`
	Nonce   uint64      `json:"nonce"`
}

type hostWithdrawalEvent struct {
	Type    string            `json:"type"`
	To      ethCommon.Address `json:"to"`
	Amount  *common.Big       `json:"amount"`
	Balance *common.Big       `json:"balance"`
	Nonce   uint64            `json:"nonce"`
}

// getTestLogger creates a new logger instance for every test
func getTestLogger(t *testing.T, useNetwork bool) logger.Logger {
	if !useNetwork {
		return logger.NewLogger(&logger.Config{
			Kind:         "zerolog",
			ConsoleLevel: "info",
			Console:      true,
			ConsoleColor: false,
		})
	}

	// Network logger config
	return logger.NewLogger(&logger.Config{
		Kind:             "zeronetwork",
		ConsoleLevel:     "trace",
		RemoteLogNetwork: "tcp",
		RemoteLogParams:  common.TcpChannelConnectionParams{Ip: "127.0.0.1", Port: 5000},
		NetworkLevel:     "trace",
	})
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}

// buildAndLoadWasmModule is a helper function to build the wasm module and read its bytecode.
func buildAndLoadWasmModule(t *testing.T) []byte {
	// Get the project root directory to construct absolute paths
	_, b, _, ok := runtime.Caller(0)
	require.True(t, ok)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	appDir := filepath.Join(projectRoot, "app", "simple")

	// Build the wasm module
	cmd := exec.Command("make", "build")
	cmd.Dir = appDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed to build wasm module: %s", string(output))

	// Load wasm bytecode for the wasm app
	wasmPath := filepath.Join(appDir, "build", "simple_app.wasm")
	wasmBytecode, err := os.ReadFile(wasmPath)
	require.NoError(t, err)
	require.NotEmpty(t, wasmBytecode)

	return wasmBytecode
}

// deploySimpleApp is a helper function to deploy the simple app wasm module.
func deploySimpleApp(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID common.ApplicationIdType, deployReqID common.RequestIdType, wasmBytecode []byte) {
	t.Helper()
	timeout := 20 * time.Second
	deployPayload := stageArtifactAndBuildDescriptorPayload(t, suite, appID, wasmBytecode)

	// Create and submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     deployReqID,
		Payload:       deployPayload,
		Sender:        suite.GetManagerAllowedDeployer(),
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))

	// Wait for app to be deployed
	_, err := suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)
	_, err = suite.WaitForAppStateInBlockchain(appID, timeout)
	require.NoError(t, err)
	require.NoError(t, suite.AssertRequestCompleted(deployReqID, timeout))

	// Verify updatePayload signature
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)
	payload, err := suite.GetRequestUpdatePayload(deployReqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)
}

// depositToSimpleApp is a helper function to deposit funds into the simple app.
func depositToSimpleApp(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID common.ApplicationIdType, reqID common.RequestIdType, user ethCommon.Address, amount *big.Int) {
	t.Helper()
	timeout := 100 * time.Second

	// Get executor's communication key for encryption
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// Create and submit deposit request
	depositReq, err := cryptoHelper.CreateDepositRequest(
		appID,
		reqID,
		user,
		amount,
		executorPubKey,
	)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(depositReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	// Wait for, decrypt and verify deposit event
	depositEvent, err := suite.WaitForEvent(user, "deposit", timeout)
	require.NoError(t, err)
	decryptedDepositData, err := cryptoHelper.DecryptEvent(user, depositEvent, executorPubKey)
	require.NoError(t, err)

	var depositEventData hostDepositEvent
	err = json.Unmarshal(decryptedDepositData, &depositEventData)
	require.NoError(t, err)
	require.Equal(t, "deposit", depositEventData.Type)
	require.Equal(t, 0, amount.Cmp(depositEventData.Amount.ToInt()))

	// Verify updatePayload signature
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)
	payload, err := suite.GetRequestUpdatePayload(reqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)
}

// withdrawFromSimpleApp is a helper function to withdraw funds from the simple app.
func withdrawFromSimpleApp(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID common.ApplicationIdType, reqID common.RequestIdType, user, recipient ethCommon.Address, amount *big.Int) {
	t.Helper()
	timeout := 100 * time.Second

	// Get executor's communication key for encryption
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// Create and submit withdrawal request
	withdrawalReq, err := cryptoHelper.CreateWithdrawalRequest(
		appID,
		reqID,
		user,
		recipient,
		common.ToBig(amount),
		executorPubKey,
	)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(withdrawalReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	// Wait for, decrypt and verify withdrawal event
	withdrawalEvent, err := suite.WaitForEvent(user, "withdrawal", timeout)
	require.NoError(t, err)
	decryptedWithdrawalData, err := cryptoHelper.DecryptEvent(user, withdrawalEvent, executorPubKey)
	require.NoError(t, err)

	var withdrawalEventData hostWithdrawalEvent
	err = json.Unmarshal(decryptedWithdrawalData, &withdrawalEventData)
	require.NoError(t, err)
	require.Equal(t, "withdrawal", withdrawalEventData.Type)
	require.Equal(t, recipient, withdrawalEventData.To)
	require.Equal(t, 0, amount.Cmp(withdrawalEventData.Amount.ToInt()))

	// Verify on-chain withdrawal was recorded
	withdrawal, err := suite.WaitForWithdrawal(appID, timeout)
	require.NoError(t, err)
	require.NotNil(t, withdrawal)
	require.Equal(t, recipient, withdrawal.DestinationAddress)
	require.Equal(t, 0, amount.Cmp(withdrawal.Amount.ToInt()))

	// Verify updatePayload signature
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)
	payload, err := suite.GetRequestUpdatePayload(reqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)
}

func TestSimpleAppDepositAndWithdraw(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	log := getTestLogger(t, true)
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", log, log)
	defer suite.Cleanup()

	wasmBytecode := buildAndLoadWasmModule(t)

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	appID := common.NewApplicationId(1)
	userAddress := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 1))
	auditorAddress := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 2))
	recipientAddress := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	timeout := 100 * time.Second

	cryptoHelper := testutil.NewCryptoHelper()

	// Deploy the application
	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	// Register user key
	userKey, err := cryptoHelper.GenerateUserKey(userAddress)
	require.NoError(t, err)
	reqID := commontestutil.GenerateRandomRequestID()
	associateKeyReq, err := cryptoHelper.CreateAssociateKeyRequest(appID, reqID, userAddress, userKey.PublicKey())
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(associateKeyReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	// Register auditor key
	auditorKey, err := cryptoHelper.GenerateUserKey(auditorAddress)
	require.NoError(t, err)
	reqID = commontestutil.GenerateRandomRequestID()
	associateAuditorReq, err := cryptoHelper.CreateAssociateKeyRequest(appID, reqID, auditorAddress, auditorKey.PublicKey())
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(associateAuditorReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	// Deposit 2 ETH and validate event fields
	depositAmount := big.NewInt(2000000000000000000)
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, depositAmount)

	// Withdraw 0.5 ETH and validate event fields
	withdrawAmount := big.NewInt(500000000000000000)
	withdrawFromSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, recipientAddress, withdrawAmount)

	// Deanonymization report as auditor — verifies final state after deposit and withdrawal
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	reqID = commontestutil.GenerateRandomRequestID()
	deanonReq, err := cryptoHelper.CreateDeanonymizationRequest(appID, reqID, auditorAddress, []byte("{}"), executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(deanonReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	deanonReport, err := suite.WaitForDeanonymizationReport(reqID, timeout)
	require.NoError(t, err)
	require.NotNil(t, deanonReport)

	decryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditorAddress, deanonReport, executorPubKey)
	require.NoError(t, err)

	var report struct {
		ApplicationId   common.ApplicationIdType `json:"applicationId"`
		RequestId       common.RequestIdType     `json:"requestId"`
		ReportDataBytes interface{}              `json:"reportDataBytes"`
	}
	err = json.Unmarshal(decryptedReport, &report)
	require.NoError(t, err)
	require.Equal(t, appID, report.ApplicationId)
	require.Equal(t, reqID, report.RequestId)

	jsonStr, ok := report.ReportDataBytes.(string)
	require.True(t, ok, "reportDataBytes is not a string")
	reportBytes, err := base64.StdEncoding.DecodeString(jsonStr)
	require.NoError(t, err, "reportDataBytes is not base64 encoded")

	var reportData map[string]interface{}
	err = json.Unmarshal(reportBytes, &reportData)
	require.NoError(t, err)
	require.Contains(t, reportData, "accounts")

	// Verify user balance reflects deposit minus withdrawal (2 ETH - 0.5 ETH = 1.5 ETH)
	accounts, ok := reportData["accounts"].(map[string]interface{})
	require.True(t, ok, "accounts is not a map")
	require.Len(t, accounts, 1, "expected exactly one account in report")

	expectedBalance := new(big.Int).Sub(depositAmount, withdrawAmount)
	for _, acct := range accounts {
		acctMap, ok := acct.(map[string]interface{})
		require.True(t, ok, "account entry is not a map")
		balanceStr, ok := acctMap["balance"].(string)
		require.True(t, ok, "balance is not a string")
		require.True(t, len(balanceStr) > 2 && balanceStr[:2] == "0x", "balance is not hex")
		balance, ok := new(big.Int).SetString(balanceStr[2:], 16)
		require.True(t, ok, "failed to parse balance hex")
		require.Equal(t, 0, expectedBalance.Cmp(balance),
			"expected balance %s, got %s", expectedBalance, balance)
	}
}

func TestExecutorManagerStart(t *testing.T) {
	log1 := getTestLogger(t, false)
	log2 := getTestLogger(t, true)

	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "wasm-runtime", mgrConfig, execConfig, nil, nil, log1, log2)
	defer suite.Cleanup()

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	time.Sleep(3 * time.Second)
}

func TestDeploySimpleApp(t *testing.T) {
	log := getTestLogger(t, true)
	// we use for both mgr and executor the remote network logger
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", log, log)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Deploy the application
	cryptoHelper := testutil.NewCryptoHelper()
	deploySimpleApp(t, suite, cryptoHelper, 1, commontestutil.GenerateRandomRequestID(), wasmBytecode)
}

// this will be modified when we support an app id other that "1"
func TestDeploySimpleAppNegativeCase(t *testing.T) {
	log1 := getTestLogger(t, false)
	log2 := getTestLogger(t, true)
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", log1, log2)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	timeout := 10 * time.Second

	// 3. Deploy the application with ID = 1
	appID := common.NewApplicationId(1)
	cryptoHelper := testutil.NewCryptoHelper()
	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	// 4. Now try to redeploy the same app id
	reqID := commontestutil.GenerateRandomRequestID()

	// Create and submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     reqID,
		Payload:       stageArtifactAndBuildDescriptorPayload(t, suite, appID, wasmBytecode),
		Sender:        suite.GetManagerAllowedDeployer(),
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))

	// we can not use suite.WaitForAppStateInDB(appID, timeout) here because we would be successful, since it checks the dataLayer
	// and we do have the appId from the previous step
	time.Sleep(timeout)

	// check that we have one more failed request
	failedRequests := suite.GetFailedRequest()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, reqID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, appID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, common.Deploy, failedRequests[0].RequestType, "Wrong Request Type")

}

func TestSimpleAppCompareAction(t *testing.T) {
	log1 := getTestLogger(t, false)
	log2 := getTestLogger(t, true)
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	timeout_value := 10 * time.Second
	// For debugging it can be useful to use huge timeout value
	//timeout_value := 10 * time.Hour

	user1Address := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 1))
	user2Address := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 2))

	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	tempDir, err := os.MkdirTemp("", "reports_system_test")
	require.NoError(t, err)
	mgrConfig.DeanonymizationReportPath = tempDir

	// we need to pass the keys for having them in the test suite
	// For tests, always use Type 0 (no KMS dependencies needed)
	ctx := context.Background()
	keySet, newRecoveryData, err := executor.GenerateEnclaveKeySet(ctx, execConfig.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "wasm-runtime", mgrConfig, execConfig, keySet, newRecoveryData, log1, log2)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Create users and add their keys to the registry
	cryptoHelper := testutil.NewCryptoHelper()
	user1Key, err := cryptoHelper.GenerateUserKey(user1Address)
	require.NoError(t, err)
	user2Key, err := cryptoHelper.GenerateUserKey(user2Address)
	require.NoError(t, err)

	// 4. Deploy the application
	appID := common.NewApplicationId(1)
	RequestID := commontestutil.GenerateRandomRequestID()
	deploySimpleApp(t, suite, cryptoHelper, appID, RequestID, wasmBytecode)

	//register key 1
	RequestID = commontestutil.GenerateRandomRequestID()
	associateKey1Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, user1Address, user1Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	//register key 3
	RequestID = commontestutil.GenerateRandomRequestID()
	associateKey2Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, user2Address, user2Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey2Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	// 5. User1 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), user1Address, big.NewInt(2000000000000000000)) // 2 ETH

	// 6. User2 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), user2Address, big.NewInt(1000))

	// Get executor's communication key for encryption, for now get from the test suite
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// 7. User1 compares balances with User2
	compareReqID := commontestutil.GenerateRandomRequestID()
	payload := map[string]interface{}{
		"type": "compare_addresses",
		"compare": map[string]ethCommon.Address{
			"targetAddress": user2Address,
		},
	}

	comparePayloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	compareReq, err := cryptoHelper.CreateProcessRequest(
		appID,
		compareReqID,
		user1Address,
		comparePayloadBytes,
		executorPubKey,
	)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(compareReq))
	require.NoError(t, suite.AssertRequestCompleted(compareReqID, timeout_value))

	// Wait for action event
	actionEvent, err := suite.WaitForEvent(user1Address, "compare_accounts", timeout_value)
	require.NoError(t, err)
	require.NotNil(t, actionEvent)

	// Decrypt and verify action event. Note that we receive it from the mock client that simulates the blockchain
	// but the same event is contained also in the updatePayload below.
	decryptedActionData, err := cryptoHelper.DecryptEvent(user1Address, actionEvent, executorPubKey)
	require.NoError(t, err)

	var eventData map[string]interface{}
	err = json.Unmarshal(decryptedActionData, &eventData)
	require.NoError(t, err)

	// 8. Fetch, decrypt, and verify the event
	updatePayload, err := suite.GetRequestUpdatePayload(compareReqID)
	require.NoError(t, err)
	require.NotEmpty(t, updatePayload.Events, "no events found for request")

	// Find the action event
	var compareEvent *common.Event
	for i := range updatePayload.Events {
		if updatePayload.Events[i].UserID == user1Address && updatePayload.Events[i].EventSubType == "compare_accounts" {
			compareEvent = &updatePayload.Events[i]
			break
		}
	}
	require.NotNil(t, compareEvent, "compare event not found for user1")

	// Decrypt the event data
	decryptedData, err := cryptoHelper.DecryptEvent(user1Address, compareEvent, executorPubKey)
	require.NoError(t, err)

	// Unmarshal and assert the content
	err = json.Unmarshal(decryptedData, &eventData)
	require.NoError(t, err)

	require.Equal(t, "compare_accounts", eventData["type"])
	require.Contains(t, strings.ToLower(eventData["sentence"].(string)), strings.ToLower(user1Address.Hex()+" is richer than "+user2Address.Hex()))
	t.Logf("Decrypted event sentence: %s", eventData["sentence"])

	require.True(t, bytes.Equal(decryptedActionData, decryptedData))

	// 9: Sending deanonymization request as auditor

	RequestID = commontestutil.GenerateRandomRequestID()
	auditorAddress := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 2))
	auditorPrivateKey, err := cryptoHelper.GenerateUserKey(auditorAddress)
	require.NoError(t, err)
	associateKey1Req, err = cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, auditorAddress, auditorPrivateKey.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	deanonReqPayload := []byte(`{"type":"deanonymize","deanonymize":{"tag":"SIMPLE_TAG"}}`)

	RequestID = commontestutil.GenerateRandomRequestID()
	deanonReq, err := cryptoHelper.CreateDeanonymizationRequest(
		appID,
		RequestID,
		auditorAddress,
		deanonReqPayload,
		executorPubKey,
	)
	require.NoError(t, err)

	err = suite.SubmitRequest(deanonReq)
	require.NoError(t, err)

	// Wait for deanonymization request to be processed
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	// Wait for deanonymization report
	deanonReport, err := suite.WaitForDeanonymizationReport(RequestID, timeout_value)
	require.NoError(t, err)
	require.NotNil(t, deanonReport)

	// 4. Read and decrypt the report
	reportFilePath := filepath.Join(suite.GetReportsPath(), common.ReportFilename(appID, RequestID))
	encryptedReportBytes, err := os.ReadFile(reportFilePath)
	require.NoError(t, err, "The report file should be saved to the filesystem")

	// sanity check of the serialized data
	var serializedEncryptedReport common.DeanonymizationReport
	err = json.Unmarshal(encryptedReportBytes, &serializedEncryptedReport)
	require.NoError(t, err)

	// Decrypt and verify both forms of deanonymization report
	decryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditorAddress, deanonReport, executorPubKey)
	require.NoError(t, err)
	serializedDecryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditorAddress, &serializedEncryptedReport, executorPubKey)
	require.NoError(t, err)

	require.True(t, bytes.Equal(decryptedReport, serializedDecryptedReport))

	// Unencrypted deanonymization reports are specific to the application, we can not assume a defined struct for the report data, but we do assume that
	// we have an appId, a reportId, and the data in separate fields
	var report struct {
		ApplicationId   common.ApplicationIdType `json:"applicationId"`
		RequestId       common.RequestIdType     `json:"requestId"`
		ReportDataBytes interface{}              `json:"reportDataBytes"`
	}
	err = json.Unmarshal(decryptedReport, &report)
	require.NoError(t, err)
	require.Equal(t, appID, report.ApplicationId)
	require.Equal(t, RequestID, report.RequestId)

	jsonStr, ok := report.ReportDataBytes.(string)
	require.True(t, ok, "reportDataBytes is not a string")
	jsonBytes, err := base64.StdEncoding.DecodeString(jsonStr)
	require.NoError(t, err, "bytes are not base64 encoded")

	// in the simple app we know how the data bytes are formatted
	var reportData map[string]interface{}
	err = json.Unmarshal(jsonBytes, &reportData)
	require.NoError(t, err)
	require.Equal(t, "SIMPLE_TAG", reportData["tag"])

	t.Log("Deanonymization report:\n", testutil.PrettyPrintJSON(report))
	t.Log("Deanonymization report data:\n", testutil.PrettyPrintJSON(reportData))

	// Deanon report does not contain signature for now, possibly add later
}

func TestSimpleApp_NegativeScenarios(t *testing.T) {
	log1 := getTestLogger(t, false)
	log2 := getTestLogger(t, true)
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	timeout_value := 10 * time.Second

	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", log1, log2)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Create user and add their key to the registry
	cryptoHelper := testutil.NewCryptoHelper()

	userAddress := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 1))

	// 4. Deploy the application
	appID := common.NewApplicationId(1)
	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	user1Key, err := cryptoHelper.GenerateUserKey(userAddress)
	require.NoError(t, err)
	requestId := commontestutil.GenerateRandomRequestID()
	associateKey1Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, requestId, userAddress, user1Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(requestId, timeout_value)
	require.NoError(t, err)

	// 5. User1 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, big.NewInt(1000))

	// Get executor's communication key for encryption
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// --- Negative Test Cases ---

	t.Run("withdraw with insufficient balance", func(t *testing.T) {
		reqID := commontestutil.GenerateRandomRequestID()
		// User1 has 1000, tries to withdraw 2000
		// payload := `{"type":"withdraw","withdraw":{"to":"0xadd0000000000000000000000000000000000003","amount":2000}}`
		payload := map[string]interface{}{
			"type": "withdraw",
			"withdraw": map[string]interface{}{
				"to":     "0xadd0000000000000000000000000000000000003",
				"amount": 2000,
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)
		processReq, err := cryptoHelper.CreateProcessRequest(
			appID,
			reqID,
			userAddress,
			payloadBytes,
			executorPubKey,
		)
		require.NoError(t, err)
		require.NoError(t, suite.SubmitRequest(processReq))

		// Assert request is completed with error
		err = suite.AssertRequestCompleted(reqID, timeout_value)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has failed")

		// Get the update payload and check for error
		_, err = suite.GetRequestUpdatePayload(reqID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		reqID := commontestutil.GenerateRandomRequestID()
		payload := `{"type":"invalid_action"}`
		processReq, err := cryptoHelper.CreateProcessRequest(
			appID,
			reqID,
			userAddress,
			[]byte(payload),
			executorPubKey,
		)
		require.NoError(t, err)
		require.NoError(t, suite.SubmitRequest(processReq))

		// Assert request is completed with error
		err = suite.AssertRequestCompleted(reqID, timeout_value)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has failed")

		// Get the update payload and check for error
		_, err = suite.GetRequestUpdatePayload(reqID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("compare with non-existent target account", func(t *testing.T) {
		reqID := commontestutil.GenerateRandomRequestID()
		nonExistentUser := "0xadd0000000000000000000000000000000000099"
		payload := `{"type":"compare_addresses","compare":{"targetAddress":"` + nonExistentUser + `"}}`
		processReq, err := cryptoHelper.CreateProcessRequest(
			appID,
			reqID,
			userAddress,
			[]byte(payload),
			executorPubKey,
		)
		require.NoError(t, err)
		require.NoError(t, suite.SubmitRequest(processReq))

		// Assert request is completed with error
		err = suite.AssertRequestCompleted(reqID, timeout_value)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has failed")

		// Get the update payload and check for error
		_, err = suite.GetRequestUpdatePayload(reqID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		reqID := commontestutil.GenerateRandomRequestID()
		payload := `{"type":"withdraw"}` // Missing withdraw payload
		processReq, err := cryptoHelper.CreateProcessRequest(
			appID,
			reqID,
			userAddress,
			[]byte(payload),
			executorPubKey,
		)
		require.NoError(t, err)
		require.NoError(t, suite.SubmitRequest(processReq))

		// Assert request is completed with error
		err = suite.AssertRequestCompleted(reqID, timeout_value)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has failed")
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		reqID := commontestutil.GenerateRandomRequestID()
		payload := `{"type":"compare_addresses"}` // Missing compare payload
		processReq, err := cryptoHelper.CreateProcessRequest(
			appID,
			reqID,
			userAddress,
			[]byte(payload),
			executorPubKey,
		)
		require.NoError(t, err)
		require.NoError(t, suite.SubmitRequest(processReq))

		// Assert request is completed with error
		err = suite.AssertRequestCompleted(reqID, timeout_value)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has failed")
	})

	t.Run("invalid deanonimization payload", func(t *testing.T) {
		reqID := commontestutil.GenerateRandomRequestID()
		payload := `{"type":"deanonymization","query":"full_report","tag":}`
		processReq, err := cryptoHelper.CreateDeanonymizationRequest(
			appID,
			reqID,
			userAddress,
			[]byte(payload),
			executorPubKey,
		)
		require.NoError(t, err)
		require.NoError(t, suite.SubmitRequest(processReq))

		// Assert request is completed with error
		err = suite.AssertRequestCompleted(reqID, timeout_value)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has failed")
	})
}
