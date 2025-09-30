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
	"github.com/horizen-pes/pkg/crypto"
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

func TestDeploySimpleApp(t *testing.T) {
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	timeout_value := 100 * time.Second

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start executor
	err := suite.StartExecutor()
	require.NoError(t, err)

	// 3. Start manager
	err = suite.StartManager()
	require.NoError(t, err)

	// 4. Add user keys to registry
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	err = suite.AddUserKeys("test-user", userKey.PublicKey().Bytes())
	require.NoError(t, err)

	RequestID := "233"
	ApplicationId := "1"

	// 5. Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: ApplicationId,
		RequestID:     RequestID,
		Payload:       wasmBytecode,
		Sender:        "test-user",
		Timestamp:     time.Now().Unix(),
	}
	err = suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// 6. Assert app state created in DB
	appState, err := suite.WaitForAppStateInDB(ApplicationId, timeout_value)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 7. Assert app state created in blockchain
	appState, err = suite.WaitForAppStateInBlockchain(ApplicationId, timeout_value)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 8. Assert request marked as done
	err = suite.AssertRequestCompleted(RequestID, timeout_value)
	require.NoError(t, err)
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

	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadWasmModule(t)

	// 2. Start services
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// 3. Create users and add their keys to the registry
	cryptoHelper := testutil.NewCryptoHelper()
	user1Key, err := cryptoHelper.GenerateUserKey("user1")
	require.NoError(t, err)
	user2Key, err := cryptoHelper.GenerateUserKey("user2")
	require.NoError(t, err)

	require.NoError(t, suite.AddUserKeys("user1", user1Key.PublicKey().Bytes()))
	require.NoError(t, suite.AddUserKeys("user2", user2Key.PublicKey().Bytes()))

	// 4. Deploy the application
	deployReqID := "1"
	appID := "1"
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     deployReqID,
		Payload:       wasmBytecode,
		Sender:        "user1",
		Timestamp:     time.Now().Unix(),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))

	// Wait for app to be deployed
	appState, err := suite.WaitForAppStateInDB(appID, timeout_value)
	require.NoError(t, err)
	require.NotNil(t, appState)

	appState, err = suite.WaitForAppStateInBlockchain(appID, timeout_value)
	require.NoError(t, err)
	require.NotNil(t, appState)

	require.NoError(t, suite.AssertRequestCompleted(deployReqID, 20*time.Second))

	// Get executor's signing key for signature verification
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)

	// Verify updatePayload signature
	payload, err := suite.GetRequestUpdatePayload(deployReqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

	// Get executor's communication key for encryption, for now get from the test suite
	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// 5. User1 deposits funds
	deposit1ReqID := "2"
	deposit1Amount := uint64(2000)
	deposit1Req, err := cryptoHelper.CreateDepositRequest(
		appID,
		deposit1ReqID,
		"user1",
		deposit1Amount,
		executorPubKey,
	)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(deposit1Req))
	require.NoError(t, suite.AssertRequestCompleted(deposit1ReqID, timeout_value))

	// Wait for deposit event
	deposit1Event, err := suite.WaitForEvent("user1", timeout_value)
	require.NoError(t, err)
	require.NotNil(t, deposit1Event)

	// Decrypt and verify deposit event
	decryptedDepositData, err := cryptoHelper.DecryptEvent("user1", deposit1Event, executorPubKey)
	require.NoError(t, err)

	var depositEventData appCommon.DepositEvent
	err = json.Unmarshal(decryptedDepositData, &depositEventData)
	require.NoError(t, err)
	require.Equal(t, "deposit", depositEventData.Type)
	require.Equal(t, deposit1Amount, depositEventData.Amount)

	// Verify updatePayload signature
	payload, err = suite.GetRequestUpdatePayload(deposit1ReqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

	// 6. User2 deposits funds
	deposit2ReqID := "3"
	deposit2Amount := uint64(1000)
	deposit2Req, err := cryptoHelper.CreateDepositRequest(
		appID,
		deposit2ReqID,
		"user2",
		deposit2Amount,
		executorPubKey,
	)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(deposit2Req))
	require.NoError(t, suite.AssertRequestCompleted(deposit2ReqID, timeout_value))

	// Wait for deposit event
	deposit2Event, err := suite.WaitForEvent("user2", timeout_value)
	require.NoError(t, err)
	require.NotNil(t, deposit2Event)

	// Decrypt and verify deposit event
	decryptedDepositData, err = cryptoHelper.DecryptEvent("user2", deposit2Event, executorPubKey)
	require.NoError(t, err)

	err = json.Unmarshal(decryptedDepositData, &depositEventData)
	require.NoError(t, err)
	require.Equal(t, "deposit", depositEventData.Type)
	require.Equal(t, deposit2Amount, depositEventData.Amount)

	// Verify updatePayload signature
	payload, err = suite.GetRequestUpdatePayload(deposit2ReqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

	// 7. User1 compares balances with User2
	compareReqID := "4"
	comparePayload := `{"type":"compare_addresses","compare":{"targetAddress":"user2"}}`
	require.NoError(t, err)
	compareReq, err := cryptoHelper.CreateProcessRequest(
		appID,
		compareReqID,
		"user1",
		[]byte(comparePayload),
		executorPubKey,
	)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(compareReq))
	require.NoError(t, suite.AssertRequestCompleted(compareReqID, timeout_value))

	// Wait for action event
	actionEvent, err := suite.WaitForEvent("user1", timeout_value)
	require.NoError(t, err)
	require.NotNil(t, actionEvent)

	// Decrypt and verify action event. Note that we receive it from the mock client that simulates the blockchain
	// but the same event is contained also in the updatePayload below.
	decryptedActionData, err := cryptoHelper.DecryptEvent("user1", actionEvent, executorPubKey)
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
		if updatePayload.Events[i].UserID == "user1" {
			compareEvent = &updatePayload.Events[i]
			break
		}
	}
	require.NotNil(t, compareEvent, "compare event not found for user1")

	// Decrypt the event data
	decryptedData, err := cryptoHelper.DecryptEvent("user1", compareEvent, executorPubKey)
	require.NoError(t, err)

	// Unmarshal and assert the content
	err = json.Unmarshal(decryptedData, &eventData)
	require.NoError(t, err)

	require.Equal(t, "compare_accounts", eventData["type"])
	require.Contains(t, eventData["sentence"], "user1 is richer than user2")
	t.Logf("Decrypted event sentence: %s", eventData["sentence"])

	require.True(t, bytes.Equal(decryptedActionData, decryptedData))

	// 9: Sending deanonymization request as auditor

	RequestID := "5"

	auditorAddress := fmt.Sprintf("0xadd%037x", 2)
	auditorPrivateKey, err := cryptoHelper.GenerateUserKey(auditorAddress)
	require.NoError(t, err)
	err = suite.AddUserKeys(auditorAddress, auditorPrivateKey.PublicKey().Bytes())
	require.NoError(t, err)

	deanonReq, err := cryptoHelper.CreateDeanonymizationRequest(
		appID,
		RequestID,
		auditorAddress,
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
