package system

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	simpleapp "github.com/HorizenOfficial/vela/app/simple/app"
	"github.com/HorizenOfficial/vela/pkg/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/executor"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/manager"
	"github.com/HorizenOfficial/vela/pkg/testutil"
	ethCommon "github.com/ethereum/go-ethereum/common"
	velacommon "github.com/HorizenOfficial/vela-common-go/common"
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

var deployRequestSender = ethCommon.HexToAddress("0x1000000000000000000000000000000000000001")

// getTestLogger creates a new logger instance for every test
func consoleLogConfig() *logger.Config {
	return &logger.Config{
		Kind:         "zerolog",
		ConsoleLevel: "info",
		Console:      true,
		ConsoleColor: false,
	}
}

// networkLogConfig returns a zeronetwork logger config.
// The suite injects the correct log-server port into RemoteLogParams
// before creating the logger, so no hardcoded port is needed here.
func networkLogConfig() *logger.Config {
	return &logger.Config{
		Kind:             "zeronetwork",
		ConsoleLevel:     "info",
		Console:          true,
		ConsoleColor:     false,
		RemoteLogNetwork: "tcp",
		NetworkLevel:     "trace",
	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}

// deploySimpleApp deploys the simple app without constructor params (ETH-only).
// Convenience wrapper around deploySimpleAppWithTokens with nil token list.
func deploySimpleApp(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID common.ApplicationIdType, deployReqID common.RequestIdType, wasmBytecode []byte) {
	t.Helper()
	deploySimpleAppWithTokens(t, suite, cryptoHelper, appID, deployReqID, wasmBytecode, nil)
}

// registerUserKey generates a P521 key for the user, submits an AssociateKey request, and waits for completion.
func registerUserKey(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID common.ApplicationIdType, userAddress ethCommon.Address, executorPubKey *cryptotypes.PublicKeyP521, timeout time.Duration) {
	t.Helper()
	key, err := cryptoHelper.GenerateUserKey(userAddress)
	require.NoError(t, err)
	reqID := commontestutil.GenerateRandomRequestID()
	req, err := cryptoHelper.CreateAssociateKeyRequest(appID, reqID, userAddress, key.PublicKey(), executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(req))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))
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

	// Wait for, decrypt and verify deposit event (using privacy-preserving subtype set)
	userSeed, err := cryptoHelper.ComputeSeed(user)
	require.NoError(t, err)
	depositEvent, err := suite.WaitForEventBySubtypes(user, executor.AllSubtypes(userSeed, executor.DefaultSubtypeN), timeout)
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

	// Wait for, decrypt and verify withdrawal event (using privacy-preserving subtype set)
	userSeed, err := cryptoHelper.ComputeSeed(user)
	require.NoError(t, err)
	withdrawalEvent, err := suite.WaitForEventBySubtypes(user, executor.AllSubtypes(userSeed, executor.DefaultSubtypeN), timeout)
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

	logCfg := consoleLogConfig()
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", logCfg, logCfg)
	defer suite.Cleanup()

	wasmBytecode := buildWasmApp(t, "simple")

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	appID := common.NewApplicationId(1)
	recipientAddress := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	timeout := 100 * time.Second

	cryptoHelper := testutil.NewCryptoHelper()

	// Generate user identities (address derived from secp256k1 signing key)
	userAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)
	auditorAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)

	// Deploy the application
	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	// Get executor's communication key for encryption
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// Register user and auditor keys
	registerUserKey(t, suite, cryptoHelper, appID, userAddress, executorPubKey, timeout)
	registerUserKey(t, suite, cryptoHelper, appID, auditorAddress, executorPubKey, timeout)

	// Deposit 2 ETH and validate event fields
	depositAmount := big.NewInt(2000000000000000000)
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, depositAmount)

	// Withdraw 0.5 ETH and validate event fields
	withdrawAmount := big.NewInt(500000000000000000)
	withdrawFromSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, recipientAddress, withdrawAmount)

	// Deanonymization report as auditor — verifies final state after deposit and withdrawal

	reqID := commontestutil.GenerateRandomRequestID()
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
	ethZeroAddr := "0x0000000000000000000000000000000000000000"
	for _, acct := range accounts {
		acctMap, ok := acct.(map[string]interface{})
		require.True(t, ok, "account entry is not a map")
		balancesMap, ok := acctMap["balances"].(map[string]interface{})
		require.True(t, ok, "balances is not a map")
		balanceStr, ok := balancesMap[ethZeroAddr].(string)
		require.True(t, ok, "ETH balance is not a string")
		require.True(t, len(balanceStr) > 2 && balanceStr[:2] == "0x", "balance is not hex")
		balance, ok := new(big.Int).SetString(balanceStr[2:], 16)
		require.True(t, ok, "failed to parse balance hex")
		require.Equal(t, 0, expectedBalance.Cmp(balance),
			"expected balance %s, got %s", expectedBalance, balance)
	}
}

func TestExecutorManagerStart(t *testing.T) {
	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()

	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "wasm-runtime", mgrConfig, execConfig, nil, nil, mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	// Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	time.Sleep(3 * time.Second)
}

func TestDeploySimpleApp(t *testing.T) {
	logCfg := consoleLogConfig()
	// we use for both mgr and executor the remote network logger
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", logCfg, logCfg)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildWasmApp(t, "simple")

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Deploy the application
	cryptoHelper := testutil.NewCryptoHelper()
	deploySimpleApp(t, suite, cryptoHelper, 1, commontestutil.GenerateRandomRequestID(), wasmBytecode)
}

// this will be modified when we support an app id other that "1"
func TestDeploySimpleAppNegativeCase(t *testing.T) {
	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildWasmApp(t, "simple")

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
		Payload:       uploadArtifactAndBuildDescriptorPayload(t, suite, wasmBytecode),
		Sender:        deployRequestSender,
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		TokenAddress:  velacommon.ETH_TOKEN,
		AssetAmount:   common.NewBig(0),
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
	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	timeout_value := 10 * time.Second
	// For debugging it can be useful to use huge timeout value
	//timeout_value := 10 * time.Hour

	cryptoHelper := testutil.NewCryptoHelper()
	user1Address, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)
	user2Address, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)

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
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "wasm-runtime", mgrConfig, execConfig, keySet, newRecoveryData, mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildWasmApp(t, "simple")

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Deploy the application
	appID := common.NewApplicationId(1)
	RequestID := commontestutil.GenerateRandomRequestID()
	deploySimpleApp(t, suite, cryptoHelper, appID, RequestID, wasmBytecode)

	// Get executor's communication key for encryption, for now get from the test suite
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// 4. Register user keys
	registerUserKey(t, suite, cryptoHelper, appID, user1Address, executorPubKey, timeout_value)
	registerUserKey(t, suite, cryptoHelper, appID, user2Address, executorPubKey, timeout_value)

	// 5. User1 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), user1Address, big.NewInt(2000000000000000000)) // 2 ETH

	// 6. User2 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), user2Address, big.NewInt(1000))

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

	// Wait for action event (using privacy-preserving subtype set)
	user1Seed, err := cryptoHelper.ComputeSeed(user1Address)
	require.NoError(t, err)
	actionEvent, err := suite.WaitForEventBySubtypes(user1Address, executor.AllSubtypes(user1Seed, executor.DefaultSubtypeN), timeout_value)
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

	// Find the action event by matching user1's privacy-preserving subtype set
	user1SubtypeSet := make(map[[32]byte]struct{})
	for _, st := range executor.AllSubtypes(user1Seed, executor.DefaultSubtypeN) {
		user1SubtypeSet[st] = struct{}{}
	}
	var compareEvent *common.Event
	for i := range updatePayload.Events {
		_, inSet := user1SubtypeSet[updatePayload.Events[i].EventSubType]
		if updatePayload.Events[i].UserID == user1Address && inSet {
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

	// 9. Register auditor and send deanonymization request

	auditorAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)
	registerUserKey(t, suite, cryptoHelper, appID, auditorAddress, executorPubKey, timeout_value)

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

	// 10. Read and decrypt the report
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
	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	timeout_value := 10 * time.Second

	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildWasmApp(t, "simple")

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Generate user identity and deploy the application
	cryptoHelper := testutil.NewCryptoHelper()

	userAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)

	appID := common.NewApplicationId(1)
	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	// 4. Register user key
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	registerUserKey(t, suite, cryptoHelper, appID, userAddress, executorPubKey, timeout_value)

	// 5. Deposit funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, big.NewInt(1000))

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

// --- ERC-20 system tests ---

// erc20TokenAddress is a synthetic ERC-20 token address used in system tests.
// It does not need to be a real contract since the token logic is handled entirely
// inside the WASM guest (app-level allowlist), not on-chain in these tests.
var erc20TokenAddress = ethCommon.HexToAddress("0xdead000000000000000000000000000000000001")

// deploySimpleAppWithTokens deploys the simple app with an ERC-20 token allowlist.
func deploySimpleAppWithTokens(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID common.ApplicationIdType, deployReqID common.RequestIdType, wasmBytecode []byte, allowedTokens []string) {
	t.Helper()
	timeout := 20 * time.Second

	params, err := json.Marshal(map[string]interface{}{
		"allowedTokens": allowedTokens,
	})
	require.NoError(t, err)

	deployPayload := uploadArtifactAndBuildDescriptorPayloadWithParams(t, suite, wasmBytecode, params)

	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     deployReqID,
		Payload:       deployPayload,
		Sender:        deployRequestSender,
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		TokenAddress:  velacommon.ETH_TOKEN,
		AssetAmount:   common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))

	_, err = suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)
	_, err = suite.WaitForAppStateInBlockchain(appID, timeout)
	require.NoError(t, err)
	require.NoError(t, suite.AssertRequestCompleted(deployReqID, timeout))

	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)
	payload, err := suite.GetRequestUpdatePayload(deployReqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)
}

// TestSimpleAppERC20DepositAndWithdraw exercises the full ERC-20 lifecycle through the
// system stack: deploy with token allowlist -> deposit ERC-20 -> withdraw ERC-20 -> verify
// final balance via deanonymization report.
func TestSimpleAppERC20DepositAndWithdraw(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	logCfg := consoleLogConfig()
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", logCfg, logCfg)
	defer suite.Cleanup()

	wasmBytecode := buildWasmApp(t, "simple")

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	appID := common.NewApplicationId(1)
	recipientAddress := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	timeout := 100 * time.Second

	cryptoHelper := testutil.NewCryptoHelper()

	// Generate user identity
	userAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)

	// 1. Deploy with ERC-20 token in the allowlist
	deploySimpleAppWithTokens(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode, []string{erc20TokenAddress.Hex()})

	// 2. Register user key
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	registerUserKey(t, suite, cryptoHelper, appID, userAddress, executorPubKey, timeout)

	// 3. Deposit 1000 units of ERC-20 token
	erc20DepositAmount := big.NewInt(1000)
	reqID := commontestutil.GenerateRandomRequestID()
	depositReq, err := cryptoHelper.CreateTokenDepositRequest(appID, reqID, userAddress, erc20TokenAddress, erc20DepositAmount, executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(depositReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	// Verify deposit event carries token address
	userSeed, err := cryptoHelper.ComputeSeed(userAddress)
	require.NoError(t, err)
	depositEvent, err := suite.WaitForEventBySubtypes(userAddress, executor.AllSubtypes(userSeed, executor.DefaultSubtypeN), timeout)
	require.NoError(t, err)
	decryptedDepositData, err := cryptoHelper.DecryptEvent(userAddress, depositEvent, executorPubKey)
	require.NoError(t, err)

	var depositEventData struct {
		Type         string      `json:"type"`
		TokenAddress string      `json:"tokenAddress"`
		Amount       *common.Big `json:"amount"`
	}
	err = json.Unmarshal(decryptedDepositData, &depositEventData)
	require.NoError(t, err)
	require.Equal(t, "deposit", depositEventData.Type)
	require.Equal(t, strings.ToLower(erc20TokenAddress.Hex()), strings.ToLower(depositEventData.TokenAddress))
	require.Equal(t, 0, erc20DepositAmount.Cmp(depositEventData.Amount.ToInt()))

	// 4. Withdraw 300 units of ERC-20 token
	erc20WithdrawAmount := big.NewInt(300)
	reqID = commontestutil.GenerateRandomRequestID()
	withdrawReq, err := cryptoHelper.CreateTokenWithdrawalRequest(appID, reqID, userAddress, recipientAddress, erc20TokenAddress, common.ToBig(erc20WithdrawAmount), executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(withdrawReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	// Verify on-chain withdrawal carries the correct token address
	withdrawal, err := suite.WaitForWithdrawal(appID, timeout)
	require.NoError(t, err)
	require.NotNil(t, withdrawal)
	require.Equal(t, erc20TokenAddress, withdrawal.TokenAddress)
	require.Equal(t, recipientAddress, withdrawal.DestinationAddress)
	require.Equal(t, 0, erc20WithdrawAmount.Cmp(withdrawal.Amount.ToInt()))

	// 5. Verify final ERC-20 balance via deanonymization report (1000 - 300 = 700)
	auditorAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)
	registerUserKey(t, suite, cryptoHelper, appID, auditorAddress, executorPubKey, timeout)

	reqID = commontestutil.GenerateRandomRequestID()
	deanonReq, err := cryptoHelper.CreateDeanonymizationRequest(appID, reqID, auditorAddress, []byte("{}"), executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(deanonReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))

	deanonReport, err := suite.WaitForDeanonymizationReport(reqID, timeout)
	require.NoError(t, err)
	decryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditorAddress, deanonReport, executorPubKey)
	require.NoError(t, err)

	var report struct {
		ReportDataBytes interface{} `json:"reportDataBytes"`
	}
	err = json.Unmarshal(decryptedReport, &report)
	require.NoError(t, err)
	jsonStr, ok := report.ReportDataBytes.(string)
	require.True(t, ok)
	reportBytes, err := base64.StdEncoding.DecodeString(jsonStr)
	require.NoError(t, err)

	var reportData map[string]interface{}
	err = json.Unmarshal(reportBytes, &reportData)
	require.NoError(t, err)
	accounts, ok := reportData["accounts"].(map[string]interface{})
	require.True(t, ok)

	expectedBalance := new(big.Int).Sub(erc20DepositAmount, erc20WithdrawAmount)
	tokenHex := strings.ToLower(erc20TokenAddress.Hex())
	for _, acct := range accounts {
		acctMap, ok := acct.(map[string]interface{})
		require.True(t, ok)
		balancesMap, ok := acctMap["balances"].(map[string]interface{})
		require.True(t, ok)
		balanceStr, ok := balancesMap[tokenHex].(string)
		require.True(t, ok, "ERC-20 balance not found for token %s", tokenHex)
		require.True(t, len(balanceStr) > 2 && balanceStr[:2] == "0x")
		balance, ok := new(big.Int).SetString(balanceStr[2:], 16)
		require.True(t, ok)
		require.Equal(t, 0, expectedBalance.Cmp(balance),
			"expected ERC-20 balance %s, got %s", expectedBalance, balance)
	}
}

// TestSimpleAppETHOnlyRejectsERC20Deposit verifies that an app deployed without
// ERC-20 tokens in its allowlist rejects an ERC-20 deposit. The executor should
// return an error and the request should be marked as failed.
func TestSimpleAppETHOnlyRejectsERC20Deposit(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	logCfg := consoleLogConfig()
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", logCfg, logCfg)
	defer suite.Cleanup()

	wasmBytecode := buildWasmApp(t, "simple")

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	appID := common.NewApplicationId(1)
	timeout := 100 * time.Second

	cryptoHelper := testutil.NewCryptoHelper()

	userAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)

	// Deploy ETH-only app (no ERC-20 tokens in allowlist)
	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	// Register user key
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	registerUserKey(t, suite, cryptoHelper, appID, userAddress, executorPubKey, timeout)

	// Attempt to deposit ERC-20 token into ETH-only app — should fail
	reqID := commontestutil.GenerateRandomRequestID()
	depositReq, err := cryptoHelper.CreateTokenDepositRequest(appID, reqID, userAddress, erc20TokenAddress, big.NewInt(1000), executorPubKey)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(depositReq))

	// The request should complete with a failure status
	err = suite.AssertRequestCompleted(reqID, timeout)
	require.Error(t, err, "ERC-20 deposit to ETH-only app should fail")
	require.Contains(t, err.Error(), "has failed")
}

// TestSimpleAppDepositEmitsAppEvent verifies that a deposit produces an AppEvent
// in the UpdatePayload (non-encrypted, application-level event).
func TestSimpleAppDepositEmitsAppEvent(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	logCfg := consoleLogConfig()
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime", logCfg, logCfg)
	defer suite.Cleanup()

	wasmBytecode := buildWasmApp(t, "simple")

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	appID := common.NewApplicationId(1)
	timeout := 100 * time.Second
	cryptoHelper := testutil.NewCryptoHelper()

	userAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)

	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	registerUserKey(t, suite, cryptoHelper, appID, userAddress, executorPubKey, timeout)

	// Deposit and retrieve the update payload
	depositReqID := commontestutil.GenerateRandomRequestID()
	depositToSimpleApp(t, suite, cryptoHelper, appID, depositReqID, userAddress, big.NewInt(1000000000000000000))

	payload, err := suite.GetRequestUpdatePayload(depositReqID)
	require.NoError(t, err)

	// Verify the AppEvent is present in the payload
	require.NotEmpty(t, payload.AppEvents, "deposit should produce at least one AppEvent")
	require.Equal(t, simpleapp.SubTypeFromString("deposit_received"), payload.AppEvents[0].EventSubType)
	require.NotEmpty(t, payload.AppEvents[0].Data, "AppEvent data should not be empty")
}
