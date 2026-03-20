package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
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


func SetupNewBlockChainClient(testHelper *testutil.SimTestHelper) *BlockChainClient {
	return SetupNewBlockChainClientConnected(testHelper.Client(), testHelper.ProcessorContractAddress, testHelper.TeeSignerAddress, testHelper.ManagerAccount)

}

func setupSimTestHelper(t *testing.T, autoMining bool, teePubSecp521r1 []byte) *testutil.SimTestHelper {
	return testutil.NewSimTestHelper(t, autoMining, true, nil, teePubSecp521r1)
}

// deployApplication deploys an application via SubmitDeployRequest, completes it
// with a successful stateUpdate, and returns the derived applicationId.
func deployApplication(t *testing.T, testHelper *testutil.SimTestHelper, blockchainClient *BlockChainClient) common.ApplicationIdType {
	t.Helper()

	deployTx := testHelper.SubmitDeployRequest(nil, big.NewInt(100))
	testHelper.WaitMined(deployTx)

	deployReq, deployStateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.NotNil(t, deployReq)

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

	return deployReq.ApplicationID
}

func TestGetPendingRequests(t *testing.T) {

	testHelper := setupSimTestHelper(t, false, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	require.Equal(t, 0, len(res), "There should be zero pending request")

	pendingRequest, stateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
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
	require.Equal(t, 0, request.DepositAmount.ToInt().Sign(), "Deposit amount should be zero for deploy requests")

	pendingRequest, stateRoot, err = blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.Equal(t, request.RequestID, pendingRequest.RequestID, "RequestID should match")
	require.Equal(t, testHelper.ProtocolVersion, pendingRequest.ProtocolVersion, "Protocol version should match")
	require.Equal(t, request.ApplicationID, pendingRequest.ApplicationID, "Application ID should match")
	require.Equal(t, common.Deploy, pendingRequest.RequestType, "Request type should match")
	require.Equal(t, payload, pendingRequest.Payload, "Payload should match")
	require.Equal(t, request.Timestamp, pendingRequest.Timestamp, "Timestamp should match")

	require.Equal(t, testHelper.Deployer.From, pendingRequest.Sender, "Sender should match")
	require.Equal(t, 0, pendingRequest.DepositAmount.ToInt().Sign(), "Deposit amount should be zero for deploy requests")

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
	tx := testHelper.SubmitRequest(deployedAppId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, oldStateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
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
	tx := testHelper.SubmitRequest(deployedAppId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, oldStateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
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
	requestId, blockNumber, err := blockchainClient.SubmitRequest(context.Background(), protocolVersion, deployedAppId, requestType, payload, depositAmount, maxFeeValue)
	require.NoError(t, err)
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

			if r.ProtocolVersion != protocolVersion || r.ApplicationID != deployedAppId || r.RequestType != requestType || string(r.Payload) != string(payload) || r.DepositAmount.ToInt().Cmp(depositAmount) != 0 {
				t.Errorf(
					"Request fields do not match: got {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}, want {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}",
					r.ProtocolVersion, r.ApplicationID, r.RequestType, string(r.Payload), r.DepositAmount,
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
	pending, err := blockchainClient.GetPendingPayments(context.Background(), testHelper.Submitter.From)
	require.NoError(t, err)
	require.Equal(t, 0, pending.Sign(), "Pending payments should be zero initially")

	// Submit a request so there's something to refund
	depositAmount := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(deployedAppId, common.Process, nil, depositAmount, maxFeeValue)
	testHelper.WaitMined(tx)

	res, oldStateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
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

	// Now the submitter should have a positive pending payment balance
	pending, err = blockchainClient.GetPendingPayments(context.Background(), testHelper.Submitter.From)
	require.NoError(t, err)
	require.Equal(t, 1, pending.Sign(), "Pending payments should be positive after failed request refund")

	// minFeePerRequest is 5 (set in sim_test_helper.go), so refund = depositAmount + (maxFeeValue - 5)
	expectedRefund := new(big.Int).Add(depositAmount, new(big.Int).Sub(maxFeeValue, big.NewInt(5)))
	require.Equal(t, 0, pending.Cmp(expectedRefund), "Pending payments should equal deposit + (maxFee - minFee)")

	// Withdraw payments for the submitter
	err = blockchainClient.WithdrawPayments(context.Background(), testHelper.Submitter.From)
	require.NoError(t, err)

	// After withdrawal, pending payments should be zero
	pending, err = blockchainClient.GetPendingPayments(context.Background(), testHelper.Submitter.From)
	require.NoError(t, err)
	require.Equal(t, 0, pending.Sign(), "Pending payments should be zero after withdrawal")
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
