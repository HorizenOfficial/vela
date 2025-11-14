package system

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/testutil"
	appCommon "github.com/horizen-pes/pkg/wasm/common"
)

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
func deploySimpleApp(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID, deployReqID, sender string, wasmBytecode []byte) {
	t.Helper()
	timeout := 20 * time.Second

	// Create and submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     deployReqID,
		Payload:       wasmBytecode,
		Sender:        sender,
		Timestamp:     time.Now().Unix(),
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
func depositToSimpleApp(t *testing.T, suite *testutil.SystemTestSuite, cryptoHelper *testutil.CryptoHelper, appID, reqID, user string, amount uint64) {
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
	depositEvent, err := suite.WaitForEvent(user, timeout)
	require.NoError(t, err)
	decryptedDepositData, err := cryptoHelper.DecryptEvent(user, depositEvent, executorPubKey)
	require.NoError(t, err)

	var depositEventData appCommon.DepositEvent
	err = json.Unmarshal(decryptedDepositData, &depositEventData)
	require.NoError(t, err)
	require.Equal(t, "deposit", depositEventData.Type)
	require.Equal(t, amount, depositEventData.Amount)

	// Verify updatePayload signature
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)
	payload, err := suite.GetRequestUpdatePayload(reqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)
}

func TestExecutorManagerStart(t *testing.T) {

	mgrConfig := manager.ReadConfig()
	execConfig := executor.ReadConfig()
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "wasm-runtime", mgrConfig, execConfig, nil, nil)
	defer suite.Cleanup()

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	time.Sleep(3 * time.Second)
}

func TestDeploySimpleApp(t *testing.T) {
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Deploy the application
	cryptoHelper := testutil.NewCryptoHelper()
	deploySimpleApp(t, suite, cryptoHelper, "1", "1233", "test-user", wasmBytecode)
}

// this will be modified when we support an app id other that "1"
func TestDeploySimpleAppNegativeCase(t *testing.T) {
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Try to deploy an application  with ID != 1
	timeout := 10 * time.Second

	appID := "33"
	reqID := "22"

	// Create and submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     reqID,
		Payload:       wasmBytecode,
		Sender:        "test-user",
		Timestamp:     time.Now().Unix(),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))

	// Waiting invain for app to be deployed
	_, err := suite.WaitForAppStateInDB(appID, timeout)
	require.Error(t, err)
	require.Contains(t, err.Error(), "timeout")

	failedRequests := suite.GetFailedRequest()
	require.Equal(t, 1, len(failedRequests), "expected 1 failed request")
	require.Equal(t, reqID, failedRequests[0].RequestID, "Wrong requestID")
	require.Equal(t, appID, failedRequests[0].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, common.Deploy, failedRequests[0].RequestType, "Wrong Request Type")

	// 4. Deploy the application with ID = 1
	cryptoHelper := testutil.NewCryptoHelper()
	deploySimpleApp(t, suite, cryptoHelper, "1", "1233", "test-user", wasmBytecode)

	// 5. Now try to redeploy the same app id
	appID = "1"
	reqID = "223"

	// Create and submit deploy request
	deployReq = &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     reqID,
		Payload:       wasmBytecode,
		Sender:        "test-user",
		Timestamp:     time.Now().Unix(),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))

	// we can not use suite.WaitForAppStateInDB(appID, timeout) here because we would be successful, since it checks the dataLayer
	// and we do have the appId from the previous step
	time.Sleep(timeout)

	// check that we have one more failed request
	failedRequests = suite.GetFailedRequest()
	require.Equal(t, 2, len(failedRequests), "expected 2 failed request")
	require.Equal(t, reqID, failedRequests[1].RequestID, "Wrong requestID")
	require.Equal(t, appID, failedRequests[1].ApplicationID, "Wrong ApplicationID")
	require.Equal(t, common.Deploy, failedRequests[1].RequestType, "Wrong Request Type")

}

func TestWasmtimeRuntimeSimpleAppFullSystemFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	wasmBytecode := buildAndLoadWasmModule(t)

	testutil.ExecTestAppFullSystemFlow(t, suite, wasmBytecode)
}

func TestSimpleAppCompareAction(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	timeout_value := 10 * time.Second
	// For debugging it can be useful to use huge timeout value
	//timeout_value := 10 * time.Hour

	user1Address := fmt.Sprintf("0xadd%037x", 1)
	user2Address := fmt.Sprintf("0xadd%037x", 2)

	mgrConfig := manager.ReadConfig()
	execConfig := executor.ReadConfig()
	tempDir, err := os.MkdirTemp("", "reports_system_test")
	require.NoError(t, err)
	mgrConfig.DeanonymizationReportPath = tempDir

	// we need to pass the keys for having them in the test suite
	keySet, newRecoveryData, err := executor.GenerateEnclaveKeySet(execConfig.KeySetRecoveryType)
	require.NoError(t, err)
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "wasm-runtime", mgrConfig, execConfig, keySet, newRecoveryData)
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
	appID := "1"
	RequestID := "11"
	deploySimpleApp(t, suite, cryptoHelper, appID, RequestID, user1Address, wasmBytecode)

	//register key 1
	RequestID = "22"
	associateKey1Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, user1Address, user1Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	//register key 3
	RequestID = "33"
	associateKey2Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, user2Address, user2Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey2Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	// 5. User1 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, "44", user1Address, 2000)

	// 6. User2 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, "55", user2Address, 1000)

	// Get executor's communication key for encryption, for now get from the test suite
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// 7. User1 compares balances with User2
	compareReqID := "66"
	payload := map[string]interface{}{
		"type": "compare_addresses",
		"compare": map[string]string{
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
	actionEvent, err := suite.WaitForEvent(user1Address, timeout_value)
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
		if updatePayload.Events[i].UserID == user1Address {
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
	require.Contains(t, eventData["sentence"], user1Address+" is richer than "+user2Address)
	t.Logf("Decrypted event sentence: %s", eventData["sentence"])

	require.True(t, bytes.Equal(decryptedActionData, decryptedData))

	// 9: Sending deanonymization request as auditor

	RequestID = "07"
	auditorAddress := fmt.Sprintf("0xadd%037x", 2)
	auditorPrivateKey, err := cryptoHelper.GenerateUserKey(auditorAddress)
	require.NoError(t, err)
	associateKey1Req, err = cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, auditorAddress, auditorPrivateKey.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	deanonReqPayload := []byte(`{"type":"deanonymization","query":"full_report","tag":"SIMPLE_TAG"}`)

	RequestID = "08"
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
	reportFilePath := filepath.Join(tempDir, appID+"_"+RequestID)
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
	var report map[string]interface{}
	err = json.Unmarshal(decryptedReport, &report)
	require.NoError(t, err)
	require.Equal(t, appID, report["applicationId"])
	require.Equal(t, RequestID, report["requestId"])

	jsonStr, ok := report["reportDataBytes"].(string)
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
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}
	timeout_value := 10 * time.Second

	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Create user and add their key to the registry
	cryptoHelper := testutil.NewCryptoHelper()

	userAddress := fmt.Sprintf("0xadd%037x", 1)

	// 4. Deploy the application
	appID := "1"
	deploySimpleApp(t, suite, cryptoHelper, appID, "11", userAddress, wasmBytecode)

	user1Key, err := cryptoHelper.GenerateUserKey(userAddress)
	require.NoError(t, err)
	associateKey1Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, "22", userAddress, user1Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted("22", timeout_value)
	require.NoError(t, err)

	// 5. User1 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, "33", userAddress, 1000)

	// Get executor's communication key for encryption
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// --- Negative Test Cases ---

	t.Run("withdraw with insufficient balance", func(t *testing.T) {
		reqID := "1011"
		// User1 has 1000, tries to withdraw 2000
		payload := `{"type":"withdraw","withdraw":{"to":"0xadd0000000000000000000000000000000000003","amount":2000}}`
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

	t.Run("unsupported instruction type", func(t *testing.T) {
		reqID := "1022"
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
		reqID := "1033"
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
		reqID := "1044"
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
		reqID := "1055"
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
		reqID := "neg-6"
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
