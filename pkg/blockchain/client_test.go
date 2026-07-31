package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/blockchain/testutil"
	"github.com/HorizenOfficial/vela/pkg/common"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/stretchr/testify/require"
)

//go:generate mkdir -p ./contracts/mocktee
//go:generate solc --via-ir --combined-json abi,bin ../../contracts/contracts/mocks/MockTeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/MockTeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/MockTeeAuthenticatorAbi/combined.json --pkg mocktee --type MockTeeAuthenticator --out ./contracts/mocktee/MockTeeAuthenticator.go
//go:generate mkdir -p ./contracts/noattestationtee
//go:generate solc --via-ir --combined-json abi,bin ../../contracts/contracts/mocks/NoAttestationTeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/NoAttestationTeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/NoAttestationTeeAuthenticatorAbi/combined.json --pkg noattestationtee --type NoAttestationTeeAuthenticator --out ./contracts/noattestationtee/NoAttestationTeeAuthenticator.go
//go:generate mkdir -p ./contracts/authority
//go:generate solc --via-ir --combined-json abi,bin ../../contracts/contracts/AuthorityRegistry.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/AuthorityRegistryAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/AuthorityRegistryAbi/combined.json --pkg authority --type AuthorityRegistry --out ./contracts/authority/AuthorityRegistry.go
// MockERC20 binding pipeline deviates from the --combined-json pattern used by
// the other contracts above. MockERC20 declares `is ERC20, ERC20Permit`, and
// both parents inherit OpenZeppelin's EIP712. With --combined-json, abigen
// emits a top-level Eip712DomainOutput struct for every contract that exposes
// eip712Domain(), which produces a duplicate-declaration compile error. The
// --exc flag also doesn't filter these out.
//
// Passing solc's per-contract .abi + .bin to abigen isolates MockERC20's own
// interface — no parent-contract types are generated, no collision, and the
// output is ~1/5 the size. No pre-processing step, so `go generate` stays
// idempotent and CI-friendly.
//go:generate mkdir -p ./contracts/mockerc20
//go:generate solc --via-ir --abi --bin ../../contracts/contracts/mocks/MockERC20.sol --base-path ../.. --include-path ../../contracts/node_modules -o ../../contract_abis/MockERC20Abi --overwrite
//go:generate abigen --v2 --abi ../../contract_abis/MockERC20Abi/MockERC20.abi --bin ../../contract_abis/MockERC20Abi/MockERC20.bin --pkg mockerc20 --type MockERC20 --out ./contracts/mockerc20/MockERC20.go
// TestTrigger isolates its own .abi + .bin (like MockERC20 above) rather than
// using --combined-json: the combined output also binds sibling contracts
// (ITrigger, AbstractTrigger) that share a withdraw() tuple, which makes abigen
// emit a duplicate WithdrawOutput struct and fails to compile.
//go:generate mkdir -p ./contracts/testtrigger
//go:generate solc --via-ir --abi --bin ../../contracts/contracts/mocks/TestTrigger.sol --base-path ../.. --include-path ../../contracts/node_modules -o ../../contract_abis/TestTriggerAbi --overwrite
//go:generate abigen --v2 --abi ../../contract_abis/TestTriggerAbi/TestTrigger.abi --bin ../../contract_abis/TestTriggerAbi/TestTrigger.bin --pkg testtrigger --type TestTrigger --out ./contracts/testtrigger/TestTrigger.go
//go:generate mkdir -p ./contracts/guardedtrigger
//go:generate solc --via-ir --abi --bin ../../contracts/contracts/mocks/GuardedTrigger.sol --base-path ../.. --include-path ../../contracts/node_modules -o ../../contract_abis/GuardedTriggerAbi --overwrite
//go:generate abigen --v2 --abi ../../contract_abis/GuardedTriggerAbi/GuardedTrigger.abi --bin ../../contract_abis/GuardedTriggerAbi/GuardedTrigger.bin --pkg guardedtrigger --type GuardedTrigger --out ./contracts/guardedtrigger/GuardedTrigger.go


func SetupNewBlockChainClient(testHelper *testutil.SimTestHelper) *BlockChainClient {
	return SetupNewBlockChainClientConnected(testHelper.Client(), testHelper.ProcessorContractAddress, testHelper.TeeSignerAddress, testHelper.ManagerAccount)

}

func setupSimTestHelper(t *testing.T, autoMining bool, teePubSecp521r1 []byte) *testutil.SimTestHelper {
	useMockContracts := true
	return testutil.NewSimTestHelper(t, autoMining, useMockContracts, nil, teePubSecp521r1)
}

// deployApplication deploys an application via SubmitDeployRequest, completes it
// with a successful stateUpdate, and returns the derived applicationId.
func deployApplication(t *testing.T, testHelper *testutil.SimTestHelper, blockchainClient *BlockChainClient) common.ApplicationIdType {
	t.Helper()

	deployTx := testHelper.SubmitDeployRequest(nil, big.NewInt(100))
	testHelper.WaitMined(deployTx)

	deployReq, deployStateRoot, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)
	require.NotNil(t, deployReq)

	// After deploy submit: appLockedFunds should be 0 (no deposit, fees tracked globally)
	funds := testHelper.GetAppCustody(deployReq.ApplicationID, velacommon.ETH_TOKEN)
	require.Equal(t, 0, funds.Sign(), "appLockedFunds should be zero after submitDeployRequest (no deposit)")

	err = blockchainClient.SubmitStateUpdate(context.Background(), &common.UpdatePayload{
		ApplicationID:  deployReq.ApplicationID,
		RequestID:      deployReq.RequestID,
		PrevStateRoot:  deployStateRoot,
		NewStateRoot:   [32]byte{0x01, 0x02, 0x03},
		Events:         []common.Event{},
		Withdrawals:    []common.Withdrawal{},
		Signature:      make([]byte, 65),
		RefundAmount:   common.NewBig(95),
		ApplicationFee: common.NewBig(5),
	})
	require.NoError(t, err)

	// After deploy stateUpdate: appLockedFunds should still be 0 (no deposit was tracked)
	funds = testHelper.GetAppCustody(deployReq.ApplicationID, velacommon.ETH_TOKEN)
	require.Equal(t, 0, funds.Sign(), "appLockedFunds should be zero after deploy stateUpdate")

	return deployReq.ApplicationID
}

func TestGetPendingRequests(t *testing.T) {

	testHelper := setupSimTestHelper(t, false, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	require.Equal(t, 0, len(res), "There should be zero pending request")

	pendingRequest, stateRoot, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)
	require.Nil(t, pendingRequest, "There should be no pending request")
	require.Equal(t, [32]byte{}, stateRoot)

	//*****************************************************
	// submit deploy request
	maxFeeValue := big.NewInt(100)
	payload := ethCommon.FromHex("0x001234")
	tx := testHelper.SubmitDeployRequest(payload, maxFeeValue)

	testHelper.MineBlock()

	// wait for transaction inclusion
	testHelper.WaitMined(tx)
	fmt.Println("Request was successfully included")

	res, err = blockchainClient.GetPendingRequests(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1, len(res), "There should be one pending request")

	request := res[0]
	require.Equal(t, testHelper.ProtocolVersion, request.ProtocolVersion, "Protocol version should match")
	require.NotEqual(t, common.ApplicationIdType(0), request.ApplicationID, "Application ID should be derived and non-zero")
	require.Equal(t, common.Deploy, request.RequestType, "Request type should match")
	require.Equal(t, payload, request.Payload, "Payload should match")
	require.Equal(t, 1, request.Timestamp.ToInt().Sign(), "Timestamp should be set and positive")

	require.Equal(t, testHelper.Deployer.From, request.Sender, "Sender should match")
	require.Equal(t, 0, request.AssetAmount.ToInt().Sign(), "Deposit amount should be zero for deploy requests")

	pendingRequest, stateRoot, err = nextPendingRequest(blockchainClient)
	require.NoError(t, err)
	require.Equal(t, request.RequestID, pendingRequest.RequestID, "RequestID should match")
	require.Equal(t, testHelper.ProtocolVersion, pendingRequest.ProtocolVersion, "Protocol version should match")
	require.Equal(t, request.ApplicationID, pendingRequest.ApplicationID, "Application ID should match")
	require.Equal(t, common.Deploy, pendingRequest.RequestType, "Request type should match")
	require.Equal(t, payload, pendingRequest.Payload, "Payload should match")
	require.Equal(t, request.Timestamp, pendingRequest.Timestamp, "Timestamp should match")

	require.Equal(t, testHelper.Deployer.From, pendingRequest.Sender, "Sender should match")
	require.Equal(t, 0, pendingRequest.AssetAmount.ToInt().Sign(), "Deposit amount should be zero for deploy requests")

	currentStateRoot := testHelper.GetStateRoot(request.ApplicationID)

	require.Equal(t, currentStateRoot, stateRoot)

}

func TestSubmitStateUpdate(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// deploy an application and submit a process request
	deployedAppId := deployApplication(t, testHelper, blockchainClient)

	transferValue := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(deployedAppId, common.Process, nil, velacommon.ETH_TOKEN, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	// After submit: appLockedFunds should equal depositAmount only (fees tracked globally)
	funds := testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, funds.Cmp(transferValue), "appLockedFunds should equal depositAmount after submit")

	res, oldStateRoot, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res.ApplicationID, EncryptedData: []byte{0x04, 0x05, 0x06}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: ethCommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: common.NewBig(10)},
	}

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID:  res.ApplicationID,
		RequestID:      res.RequestID,
		PrevStateRoot:  oldStateRoot,
		NewStateRoot:   [32]byte{0x04, 0x05, 0x06},
		Events:         events[:],
		Withdrawals:    withdrawals,
		Signature:      signature[:],
		RefundAmount:   common.NewBig(90),
		ApplicationFee: common.NewBig(10), // 90 + 10 = 100 == maxFeeValue
	}

	// =========================================================
	// Case wrong application id
	// =========================================================
	payload.ApplicationID = 9999
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg := err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidApplicationId")
	payload.ApplicationID = res.ApplicationID

	// =========================================================
	// Case wrong prev state root
	// =========================================================
	payload.PrevStateRoot = [32]byte{0x07, 0x08, 0xaa, 0xbb, 0xee}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidStateRoot")
	payload.PrevStateRoot = oldStateRoot

	// =========================================================
	// Case refund + applicationFees != maxFeeValue -> InvalidValue
	// =========================================================
	payloadWrongSum := *payload // copy value
	payloadWrongSum.RefundAmount = common.NewBig(80)
	payloadWrongSum.ApplicationFee = common.NewBig(10)

	err = blockchainClient.SubmitStateUpdate(context.Background(), &payloadWrongSum)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidValue")

	// =========================================================
	// Case applicationFees < minFeePerRequest but correct sum -> InvalidValue
	// =========================================================
	payloadWrongFee := *payload
	payloadWrongFee.RefundAmount = common.NewBig(100)
	payloadWrongFee.ApplicationFee = common.NewBig(0)

	err = blockchainClient.SubmitStateUpdate(context.Background(), &payloadWrongFee)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidValue")

	// =========================================================
	// Case OK: refund=90, applicationFees=10 -> success
	// =========================================================
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	listOfRes, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(listOfRes), "There should be 0 pending request")

	// After successful stateUpdate: appLockedFunds debited by withdrawals only
	// withdrawal(10) debited from original locked = transferValue(1000000)
	// remaining = 1000000 - 10 = 999990
	fundsAfterUpdate := testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsAfterUpdate.Cmp(big.NewInt(999990)), "appLockedFunds should be debited by withdrawals")

	// =========================================================
	// Case NonceTooLow
	// =========================================================
	blockchainClient.account.Nonce = big.NewInt(0)
	blockchainClient.account.GasLimit = 100000 //more than enough to avoid out of gas error
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)

	blockchainClient.account.Nonce = nil  //reset nonce to let it be fetched from the network
	blockchainClient.account.GasLimit = 0 //reset gas limit to let it be estimated

	// =========================================================
	// Case wrong request id
	// =========================================================
	payload.PrevStateRoot = payload.NewStateRoot
	payload.NewStateRoot = [32]byte{0x07, 0x08, 0x09}
	payload.RequestID = commontestutil.GenerateRandomRequestID()
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidRequestId")
}

func TestSubmitStateUpdateRequestFailed(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// deploy an application and submit a process request
	deployedAppId := deployApplication(t, testHelper, blockchainClient)

	transferValue := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(deployedAppId, common.Process, nil, velacommon.ETH_TOKEN, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	// After submit: appLockedFunds should equal depositAmount only (fees tracked globally)
	funds := testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, funds.Cmp(transferValue), "appLockedFunds should equal depositAmount after submit")

	res, oldStateRoot, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID:  res.ApplicationID,
		RequestID:      res.RequestID,
		PrevStateRoot:  oldStateRoot,
		NewStateRoot:   oldStateRoot, // same as old state root to simulate failed request
		Events:         []common.Event{},
		Withdrawals:    []common.Withdrawal{},
		Signature:      signature[:],
		RefundAmount:   common.ToBig(maxFeeValue),
		ApplicationFee: common.NewBig(0),
		ErrorCode:      1,
		ErrorMsg:       "test error message",
	}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	listOfRes, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(listOfRes), "There should be zero pending request")

	// After failed stateUpdate: appLockedFunds should be debited by depositAmount
	fundsAfterFail := testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsAfterFail.Sign(), "appLockedFunds should be zero after failed stateUpdate")
}

func TestSubmitRequest(t *testing.T) {
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// Deploy an application first
	deployedAppId := deployApplication(t, testHelper, blockchainClient)

	// Prepare a request
	protocolVersion := uint8(0)
	requestType := common.Process
	payload := []byte("test-payload")
	depositAmount := big.NewInt(1)
	maxFeeValue := big.NewInt(100)

	// Submit the request
	requestId, blockNumber, err := blockchainClient.SubmitRequest(context.Background(), protocolVersion, deployedAppId, requestType, payload, velacommon.ETH_TOKEN, depositAmount, maxFeeValue)
	require.NoError(t, err)

	// After submit: appLockedFunds should equal depositAmount only (fees tracked globally)
	funds := testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, funds.Cmp(depositAmount), "appLockedFunds should equal depositAmount after SubmitRequest")

	// Get pending requests
	pending, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	//check block number > 0
	require.True(t, blockNumber > 0, "Returned block number shouldn't be 0")
	// Check that the submitted request is present and matches
	found := false

	for _, r := range pending {
		if r.RequestID == requestId {
			found = true

			if r.ProtocolVersion != protocolVersion || r.ApplicationID != deployedAppId || r.RequestType != requestType || string(r.Payload) != string(payload) || r.AssetAmount.ToInt().Cmp(depositAmount) != 0 {
				t.Errorf(
					"Request fields do not match: got {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}, want {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}",
					r.ProtocolVersion, r.ApplicationID, r.RequestType, string(r.Payload), r.AssetAmount,
					protocolVersion, deployedAppId, requestType, string(payload), depositAmount,
				)
			}
		}
	}
	if !found {
		t.Errorf("Submitted request not found in pending requests")
	}
}

func TestSubmitDeployRequest_AuthorizedDeployer(t *testing.T) {
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClientConnected(
		testHelper.Client(),
		testHelper.ProcessorContractAddress,
		testHelper.TeeSignerAddress,
		testHelper.Deployer,
	)

	appId, requestID, _, err := blockchainClient.SubmitDeployRequest(
		context.Background(),
		uint8(0),
		[]byte("deploy-payload"),
		big.NewInt(100),
	)
	require.NoError(t, err)
	require.NotEqual(t, common.ApplicationIdType(0), appId)

	// After deploy submit: appLockedFunds should be 0 (no deposit, fees tracked globally)
	funds := testHelper.GetAppCustody(appId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, funds.Sign(), "appLockedFunds should be zero after SubmitDeployRequest (no deposit)")

	pending, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Len(t, pending, 1)
	require.Equal(t, requestID, pending[0].RequestID)
	require.Equal(t, common.Deploy, pending[0].RequestType)
	require.Equal(t, testHelper.Deployer.From, pending[0].Sender)
	require.Equal(t, appId, pending[0].ApplicationID)
}

func TestSubmitDeployRequest_UnauthorizedDeployerReverts(t *testing.T) {
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClientConnected(
		testHelper.Client(),
		testHelper.ProcessorContractAddress,
		testHelper.TeeSignerAddress,
		testHelper.Submitter,
	)

	_, _, _, err := blockchainClient.SubmitDeployRequest(
		context.Background(),
		uint8(0),
		[]byte("deploy-payload"),
		big.NewInt(100),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointDeployerNotAllowed")
}

func TestGetPendingPaymentsAndWithdraw(t *testing.T) {
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// Deploy an application first
	deployedAppId := deployApplication(t, testHelper, blockchainClient)

	// Initially the submitter should have no pending payments
	pending, err := blockchainClient.GetPendingClaims(context.Background(), velacommon.ETH_TOKEN, testHelper.Submitter.From)
	require.NoError(t, err)
	require.Equal(t, 0, pending.Sign(), "Pending payments should be zero initially")

	// Submit a request so there's something to refund
	depositAmount := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(deployedAppId, common.Process, nil, velacommon.ETH_TOKEN, depositAmount, maxFeeValue)
	testHelper.WaitMined(tx)

	// After submit: appLockedFunds should equal depositAmount only (fees tracked globally)
	appFunds := testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, appFunds.Cmp(depositAmount), "appLockedFunds should equal depositAmount after submit")

	res, oldStateRoot, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Submit a failed state update — this credits deposit + (maxFeeValue - minFee) to the submitter
	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID:  res.ApplicationID,
		RequestID:      res.RequestID,
		PrevStateRoot:  oldStateRoot,
		NewStateRoot:   oldStateRoot,
		Events:         []common.Event{},
		Withdrawals:    []common.Withdrawal{},
		Signature:      signature[:],
		RefundAmount:   common.ToBig(maxFeeValue),
		ApplicationFee: common.NewBig(0),
		ErrorCode:      1,
		ErrorMsg:       "test failure",
	}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	// After failed stateUpdate: appLockedFunds should be zero (depositAmount debited)
	appFunds = testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, appFunds.Sign(), "appLockedFunds should be zero after failed stateUpdate")

	// Now the submitter should have a positive pending payment balance
	pending, err = blockchainClient.GetPendingClaims(context.Background(), velacommon.ETH_TOKEN, testHelper.Submitter.From)
	require.NoError(t, err)
	require.Equal(t, 1, pending.Sign(), "Pending payments should be positive after failed request refund")

	// minFeePerRequest is 5 (set in sim_test_helper.go), so refund = depositAmount + (maxFeeValue - 5)
	expectedRefund := new(big.Int).Add(depositAmount, new(big.Int).Sub(maxFeeValue, big.NewInt(5)))
	require.Equal(t, 0, pending.Cmp(expectedRefund), "Pending payments should equal deposit + (maxFee - minFee)")

	// Withdraw payments for the submitter
	err = blockchainClient.Claim(context.Background(), velacommon.ETH_TOKEN, testHelper.Submitter.From)
	require.NoError(t, err)

	// After withdrawal: appLockedFunds should still be zero (withdrawPayments does not affect appLockedFunds)
	appFunds = testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, appFunds.Sign(), "appLockedFunds should remain zero after withdrawPayments")

	// After withdrawal, pending payments should be zero
	pending, err = blockchainClient.GetPendingClaims(context.Background(), velacommon.ETH_TOKEN, testHelper.Submitter.From)
	require.NoError(t, err)
	require.Equal(t, 0, pending.Sign(), "Pending payments should be zero after withdrawal")
}

// TestInsufficientAppBalance verifies that a stateUpdate reverts with
// InsufficientAppBalance when the withdrawal sum exceeds the application's
// locked funds. A request is submitted with zero deposit and the minimum fee
// (100), so appLockedFunds = 0 (fees tracked globally). The stateUpdate then
// attempts a 1 wei withdrawal, which exceeds appLockedFunds(0).
func TestInsufficientAppBalance(t *testing.T) {
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	deployedAppId := deployApplication(t, testHelper, blockchainClient)

	// Submit request: deposit=0, fee=100 → appLockedFunds = 0 (fees not tracked per-app)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(deployedAppId, common.Process, nil, velacommon.ETH_TOKEN, big.NewInt(0), maxFeeValue)
	testHelper.WaitMined(tx)

	res, oldStateRoot, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)

	// Attempt stateUpdate: withdrawal(1) > appLockedFunds(0)
	// Expected: revert with InsufficientAppBalance
	err = blockchainClient.SubmitStateUpdate(context.Background(), &common.UpdatePayload{
		ApplicationID:  res.ApplicationID,
		RequestID:      res.RequestID,
		PrevStateRoot:  oldStateRoot,
		NewStateRoot:   [32]byte{0x0a, 0x0b, 0x0c},
		Events:         []common.Event{},
		Withdrawals:    []common.Withdrawal{{DestinationAddress: ethCommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: common.NewBig(1)}},
		Signature:      make([]byte, 65),
		RefundAmount:   common.NewBig(0),
		ApplicationFee: common.NewBig(100),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInsufficientAppBalance")

	// The reverted tx must not have changed appLockedFunds
	funds := testHelper.GetAppCustody(deployedAppId, velacommon.ETH_TOKEN)
	require.Equal(t, 0, funds.Sign(), "appLockedFunds should be unchanged (zero) after reverted stateUpdate")
}

// TestCrossAppFundIsolation verifies the core per-app solvency guarantee: an
// application cannot withdraw funds deposited by another application, even when
// the contract's global ETH balance would be sufficient.
//
// Setup: App A is funded with deposit=1000 (locked=1000), App B with
// deposit=10 (locked=10). Fees are tracked globally, not per-app. After
// processing App A's request with a 500 withdrawal, App A retains 500 locked.
// The contract's global balance is well above 500, but App B's locked funds are
// only 10. A stateUpdate for App B requesting a 500 withdrawal must revert
// with InsufficientAppBalance.
func TestCrossAppFundIsolation(t *testing.T) {
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// Step 1: Deploy two applications
	appA := deployApplication(t, testHelper, blockchainClient)
	appB := deployApplication(t, testHelper, blockchainClient)

	// Step 2: Fund App A heavily (deposit=1000, fee=100) and App B lightly (deposit=10, fee=100)
	txA := testHelper.SubmitRequest(appA, common.Process, nil, velacommon.ETH_TOKEN, big.NewInt(1000), big.NewInt(100))
	testHelper.WaitMined(txA)
	txB := testHelper.SubmitRequest(appB, common.Process, nil, velacommon.ETH_TOKEN, big.NewInt(10), big.NewInt(100))
	testHelper.WaitMined(txB)

	fundsA := testHelper.GetAppCustody(appA, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsA.Cmp(big.NewInt(1000)), "App A should have 1000 locked (deposit only)")
	fundsB := testHelper.GetAppCustody(appB, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsB.Cmp(big.NewInt(10)), "App B should have 10 locked (deposit only)")

	// Step 3: Process App A — withdraw 500 (within A's budget)
	reqA, stateRootA, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)
	require.Equal(t, appA, reqA.ApplicationID)

	err = blockchainClient.SubmitStateUpdate(context.Background(), &common.UpdatePayload{
		ApplicationID:  reqA.ApplicationID,
		RequestID:      reqA.RequestID,
		PrevStateRoot:  stateRootA,
		NewStateRoot:   [32]byte{0x0a, 0x01},
		Events:         []common.Event{},
		Withdrawals:    []common.Withdrawal{{DestinationAddress: ethCommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: common.NewBig(500)}},
		Signature:      make([]byte, 65),
		RefundAmount:   common.NewBig(95),
		ApplicationFee: common.NewBig(5),
	})
	require.NoError(t, err)

	// App A debited: withdrawal(500). Remaining: 1000 - 500 = 500
	fundsA = testHelper.GetAppCustody(appA, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsA.Cmp(big.NewInt(500)), "App A should have 500 remaining after processing")
	// App B must be unaffected by App A's stateUpdate
	fundsB = testHelper.GetAppCustody(appB, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsB.Cmp(big.NewInt(10)), "App B should still have 10 locked — unaffected by App A")

	// Step 4: Process App B — attempt to withdraw 500 (exceeds B's 10 locked funds)
	// The contract's global balance is sufficient (App A's 500 is still in the contract),
	// but the per-app solvency check must block this.
	reqB, stateRootB, err := nextPendingRequest(blockchainClient)
	require.NoError(t, err)
	require.Equal(t, appB, reqB.ApplicationID)

	err = blockchainClient.SubmitStateUpdate(context.Background(), &common.UpdatePayload{
		ApplicationID:  reqB.ApplicationID,
		RequestID:      reqB.RequestID,
		PrevStateRoot:  stateRootB,
		NewStateRoot:   [32]byte{0x0b, 0x01},
		Events:         []common.Event{},
		Withdrawals:    []common.Withdrawal{{DestinationAddress: ethCommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"), Amount: common.NewBig(500)}},
		Signature:      make([]byte, 65),
		RefundAmount:   common.NewBig(95),
		ApplicationFee: common.NewBig(5),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInsufficientAppBalance",
		"App B must not withdraw funds exceeding its own locked balance, even if global balance is sufficient")

	// Step 5: Verify neither app's funds were changed by the reverted tx
	fundsA = testHelper.GetAppCustody(appA, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsA.Cmp(big.NewInt(500)), "App A funds unchanged after B's reverted tx")
	fundsB = testHelper.GetAppCustody(appB, velacommon.ETH_TOKEN)
	require.Equal(t, 0, fundsB.Cmp(big.NewInt(10)), "App B funds unchanged after its own reverted tx")
}

func TestGetTeePublicKey(t *testing.T) {
	//generate key
	key, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	// create test with the key
	testHelper := setupSimTestHelper(t, true, key.PublicKey().Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// Get the Public Key
	publicKey, err := blockchainClient.GetTeePublicKey(context.Background())
	require.NoError(t, err)
	require.NotNil(t, publicKey, "Public key should not be nil")
	require.Equal(t, key.PublicKey().Bytes(), publicKey.Bytes(), "Public key not equal to the given one")
}

// nextPendingRequest fetches the single request the contract selects next, together with
// the selected application's on-chain state root. Returns a nil request (and the zero
// state root) when no request is pending.
func nextPendingRequest(c *BlockChainClient) (*common.Request, [32]byte, error) {
	_, requests, stateRoot, err := c.GetPendingRequestsWithStateRoot(context.Background(), 1)
	if err != nil || len(requests) == 0 {
		return nil, stateRoot, err
	}
	return requests[0], stateRoot, nil
}
