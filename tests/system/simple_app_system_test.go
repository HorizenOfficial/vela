package system

import (
	"bytes"
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

func TestDeploySimpleApp(t *testing.T) {
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Create user and add their key to the registry
	cryptoHelper := testutil.NewCryptoHelper()

	// 4. Deploy the application
	deploySimpleApp(t, suite, cryptoHelper, "1", "1233", "test-user", wasmBytecode)
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

	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
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
	RequestID := "01"
	deploySimpleApp(t, suite, cryptoHelper, appID, RequestID, user1Address, wasmBytecode)

	//register key 1
	RequestID = "02"
	associateKey1Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, user1Address, user1Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	//register key 3
	RequestID = "03"
	associateKey2Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, RequestID, user2Address, user2Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey2Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)

	// 5. User1 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, "04", user1Address, 2000)

	// 6. User2 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, "05", user2Address, 1000)

	// Get executor's communication key for encryption, for now get from the test suite
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// 7. User1 compares balances with User2
	compareReqID := "06"
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

	// Decrypt and verify deanonymization report
	decryptedReport, err := cryptoHelper.DecryptDeanonymizationReport(auditorAddress, deanonReport, executorPubKey)
	require.NoError(t, err)

	// Unencrypted deanonymization reports are specific to the application, we can not assume a defined struct, but we do assume that
	// we have at least an appId and a reportId
	var reportData map[string]interface{}
	err = json.Unmarshal(decryptedReport, &reportData)
	require.NoError(t, err)
	require.Equal(t, appID, reportData["applicationId"])
	require.Equal(t, RequestID, reportData["requestId"])
	t.Log("Deanonymization report:\n", testutil.PrettyPrintJSON(reportData))

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
	deploySimpleApp(t, suite, cryptoHelper, appID, "01", userAddress, wasmBytecode)

	user1Key, err := cryptoHelper.GenerateUserKey(userAddress)
	require.NoError(t, err)
	associateKey1Req, err := cryptoHelper.CreateAssociateKeyRequest(appID, "02", userAddress, user1Key.PublicKey())
	require.NoError(t, err)
	err = suite.SubmitRequest(associateKey1Req)
	require.NoError(t, err)
	err = suite.AssertRequestCompleted("2", timeout_value)
	require.NoError(t, err)

	// 5. User1 deposits funds
	depositToSimpleApp(t, suite, cryptoHelper, appID, "03", userAddress, 1000)

	// Get executor's communication key for encryption
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// --- Negative Test Cases ---

	t.Run("withdraw with insufficient balance", func(t *testing.T) {
		reqID := "neg-1"
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
		reqID := "neg-2"
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
		reqID := "neg-3"
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
		reqID := "neg-4"
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
		reqID := "neg-5"
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
