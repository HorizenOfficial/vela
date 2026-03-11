package system

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/testutil"
)

func TestDeployApp(t *testing.T) {
	log1 := getTestLogger(t, false)
	log2 := getTestLogger(t, true)
	suite := testutil.NewSystemTestSuite(t, "mock-runtime", log1, log2)
	defer suite.Cleanup()

	// 1. Start executor
	err := suite.StartExecutor()
	require.NoError(t, err)

	// 2. Start manager
	err = suite.StartManager()
	require.NoError(t, err)

	RequestID := commontestutil.GenerateRandomRequestID()
	ApplicationId := common.NewApplicationId(1)

	// 4. Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: ApplicationId,
		RequestID:     RequestID,
		Payload:       uploadArtifactAndBuildDescriptorPayload(t, suite, []byte("deploy-payload")),
		Sender:        deployRequestSender,
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}
	err = suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// 5. Assert app state created in DB
	appState, err := suite.WaitForAppStateInDB(ApplicationId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 6. Assert app state created in blockchain
	appState, err = suite.WaitForAppStateInBlockchain(ApplicationId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 7. Assert request marked as done
	err = suite.AssertRequestCompleted(RequestID, 10*time.Second)
	require.NoError(t, err)
}

func TestMockRuntimeFullFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	log1 := getTestLogger(t, false)
	log2 := getTestLogger(t, true)
	suite := testutil.NewSystemTestSuite(t, "mock-runtime", log1, log2)
	defer suite.Cleanup()

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	appID := common.NewApplicationId(1)
	userAddress := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 1))
	auditorAddress := ethCommon.HexToAddress(fmt.Sprintf("0xadd%037x", 2))
	recipientAddress := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	timeout := 100 * time.Second

	cryptoHelper := testutil.NewCryptoHelper()

	// Deploy with mock bytecode
	deployReqID := commontestutil.GenerateRandomRequestID()
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     deployReqID,
		Payload:       uploadArtifactAndBuildDescriptorPayload(t, suite, []byte("mock-runtime-app-bytecode")),
		Sender:        deployRequestSender,
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))
	_, err := suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)
	_, err = suite.WaitForAppStateInBlockchain(appID, timeout)
	require.NoError(t, err)
	require.NoError(t, suite.AssertRequestCompleted(deployReqID, timeout))

	// Verify deploy signature
	executorSigningKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)
	payload, err := suite.GetRequestUpdatePayload(deployReqID)
	require.NoError(t, err)
	err = cryptoHelper.ValidateUpdatePayloadSignature(payload, executorSigningKey)
	require.NoError(t, err)

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

	// Deposit 2 ETH
	depositAmount := big.NewInt(2000000000000000000)
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, depositAmount)

	// Withdraw 0.5 ETH
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
