package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	commontestutil "github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/crypto"
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

var (
	applicationId = common.NewApplicationId(1)
)

func SetupNewBlockChainClient(testHelper *testutil.SimTestHelper) *BlockChainClient {
	return SetupNewBlockChainClientConnected(testHelper.Client(), testHelper.ProcessorContractAddress, testHelper.TeeSignerAddress, testHelper.ManagerAccount)

}

func setupSimTestHelper(t *testing.T, autoMining bool, teePubSecp521r1 []byte) *testutil.SimTestHelper {
	return testutil.NewSimTestHelper(t, autoMining, true, nil, teePubSecp521r1)
}

func TestGetPendingRequests(t *testing.T) {

	testHelper := setupSimTestHelper(t, false, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	require.Equal(t, 0, len(res), "There should be zero pending request")

	currentStateRoot := testHelper.GetStateRoot()
	pendingRequest, stateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.Nil(t, pendingRequest, "There should be no pending request")
	require.Equal(t, currentStateRoot, stateRoot)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1203055)
	maxFeeValue := big.NewInt(100)
	payload := ethCommon.FromHex("0x001234")
	tx := testHelper.SubmitRequest(applicationId, common.Process, payload, transferValue, maxFeeValue)

	testHelper.MineBlock()

	// wait for transaction inclusion
	testHelper.WaitMined(tx)
	fmt.Println("Request was successfully included")

	res, err = blockchainClient.GetPendingRequests(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1, len(res), "There should be one pending request")

	request := res[0]
	require.Equal(t, testHelper.ProtocolVersion, request.ProtocolVersion, "Protocol version should match")
	require.Equal(t, applicationId, request.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, request.RequestType, "Request type should match")
	require.Equal(t, payload, request.Payload, "Payload should match")
	require.Equal(t, 1, request.Timestamp.ToInt().Sign(), "Timestamp should be set and positive")

	require.Equal(t, testHelper.Submitter.From, request.Sender, "Sender should match")
	require.Equal(t, 0, request.DepositAmount.ToInt().Cmp(transferValue), "Value should match")

	pendingRequest, stateRoot, err = blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.Equal(t, request.RequestID, pendingRequest.RequestID, "RequestID should match")
	require.Equal(t, testHelper.ProtocolVersion, pendingRequest.ProtocolVersion, "Protocol version should match")
	require.Equal(t, applicationId, pendingRequest.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, pendingRequest.RequestType, "Request type should match")
	require.Equal(t, payload, pendingRequest.Payload, "Payload should match")
	require.Equal(t, request.Timestamp, pendingRequest.Timestamp, "Timestamp should match")

	require.Equal(t, testHelper.Submitter.From, pendingRequest.Sender, "Sender should match")
	require.Equal(t, 0, pendingRequest.DepositAmount.ToInt().Cmp(transferValue), "Value should match")

	require.Equal(t, currentStateRoot, stateRoot)

}

func TestMarkRequestCompleted(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(0)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, len(res), "There should be one pending request")
	requestId := res[0].RequestID

	// =========================================================
	// Case 1: refund + applicationFees != maxFeeValue -> InvalidValue
	// =========================================================
	err = blockchainClient.MarkRequestCompleted(
		context.Background(),
		requestId,
		big.NewInt(50), // refund
		big.NewInt(20), // fees  => 50 + 20 = 70 != 100
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidValue")

	// =========================================================
	// Case 2: applicationFees < minFeePerRequest -> InvalidValue
	// =========================================================
	err = blockchainClient.MarkRequestCompleted(
		context.Background(),
		requestId,
		big.NewInt(100), // refund
		big.NewInt(0),   // fees => 100 + 0 = 100 == maxFeeValue, but fees < minFeePerRequest
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidValue")

	// =========================================================
	// Case OK: refund + fees = maxFeeValue, fees >= minFeePerRequest
	// =========================================================
	err = blockchainClient.MarkRequestCompleted(
		context.Background(),
		requestId,
		big.NewInt(80),
		big.NewInt(20),
	)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")

	// Test that completing the same request results in ProcessorEndpointInvalidRequestId
	err = blockchainClient.MarkRequestCompleted(
		context.Background(),
		requestId,
		big.NewInt(80),
		big.NewInt(20),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidRequestId")
	_, isReorgErr := err.(ReorgError)
	require.True(t, isReorgErr)
}

func TestMarkRequestFailed(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	transferValue := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	failure := apperrors.New(
		apperrors.CodeSubmittingStateUpdateFailed,
		"test failure",
		nil,
	)

	err = blockchainClient.MarkRequestFailed(context.Background(), res[0].RequestID, failure)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")
}

func TestSubmitStateUpdate(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)

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
	// Case 1: refund + applicationFees != maxFeeValue -> InvalidValue
	// =========================================================
	payloadWrongSum := *payload // copy value
	payloadWrongSum.RefundAmount = common.NewBig(80)
	payloadWrongSum.ApplicationFee = common.NewBig(10)

	err = blockchainClient.SubmitStateUpdate(context.Background(), &payloadWrongSum)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidValue")

	// =========================================================
	// Case 2: applicationFees < minFeePerRequest but correct sum -> InvalidValue
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

	// Test error - NonceTooLow
	blockchainClient.account.Nonce = big.NewInt(0)
	blockchainClient.account.GasLimit = 100000 //more than enough to avoid out of gas error
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg := err.(ReorgError)
	require.True(t, isReorg)

	blockchainClient.account.Nonce = nil  //reset nonce to let it be fetched from the network
	blockchainClient.account.GasLimit = 0 //reset gas limit to let it be estimated

	// Test error - wrong application id
	payload.ApplicationID = 9999
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidApplicationId")

	// Test error - wrong prev state root
	payload.ApplicationID = res.ApplicationID
	payload.PrevStateRoot = [32]byte{0x07, 0x08, 0xaa, 0xbb, 0xee}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	fmt.Printf("DEBUG update error: %v\n", err)

	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidStateRoot")

	// Test error - wrong request id
	payload.PrevStateRoot = payload.NewStateRoot
	payload.NewStateRoot = [32]byte{0x07, 0x08, 0x09}
	payload.RequestID = commontestutil.GenerateRandomRequestID()
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidRequestId")
}

func TestGetUserEvents_StopAtFirst(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()

	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	messageSkipped := "test message skipped"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageSkipped, teeKey, userPub)
	message := "test message"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message, teeKey, userPub)

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, applicationId, 0, 0, "", nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(userEvents), "There should be 1 user event")

	require.Equal(t, []byte(message), userEvents[0], "Decrypted message should match original")
}

func TestGetUserEvents_MultipleEvents(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()

	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	message1 := "test message 1"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message1, teeKey, userPub)
	message2 := "test message 2"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message2, teeKey, userPub)
	message3 := "test message 3"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message3, teeKey, userPub)

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, applicationId, 0, 0, "", nil, false)
	require.NoError(t, err)
	require.Equal(t, 3, len(userEvents), "There should be 3 user event")

	// they are in reverse order
	require.Equal(t, []byte(message1), userEvents[2], "Decrypted message should match original (1)")
	require.Equal(t, []byte(message2), userEvents[1], "Decrypted message should match original (2)")
	require.Equal(t, []byte(message3), userEvents[0], "Decrypted message should match original (2)")

}

func TestGetUserEvents_WithFilter(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()

	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	messageTrue := "test message - true"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageTrue, teeKey, userPub)
	messageFalse := "test message - false"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageFalse, teeKey, userPub)

	//filter function
	filter := func(data []byte) bool {
		return strings.Contains(string(data), "true")
	}

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, applicationId, 0, 0, "", filter, false)
	require.NoError(t, err)
	require.Equal(t, 1, len(userEvents), "There should be 1 user event")
	require.Equal(t, []byte(messageTrue), userEvents[0], "Decrypted message should match the one that passes the filter")

}

func TestGetUserEvents_OtherUsersEvents(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()

	//generate key for another user
	otherUserKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate other user private key")
	otherUserPub := otherUserKey.PublicKey()

	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	message := "test message"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message, teeKey, userPub)
	messageOther := "test message - for other user"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageOther, teeKey, otherUserPub)

	//retrieve and decrypt user events
	// for user
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, applicationId, 0, 0, "", nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(userEvents), "There should be 1 user event")
	require.Equal(t, []byte(message), userEvents[0], "Decrypted message should match original")

	// for other user
	otherUserEvents, err := blockchainClient.GetUserEvents(context.Background(), *otherUserKey, applicationId, 0, 0, "", nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(otherUserEvents), "There should be 1 user event (other)")
	require.Equal(t, []byte(messageOther), otherUserEvents[0], "Decrypted message should match original (other)")
}

func _submitRequestAndStateUpdateWithEncryptedMessageEvent(t *testing.T, blockchainClient *BlockChainClient, testHelper *testutil.SimTestHelper,
	message string, senderPrivKey *cryptotypes.PrivateKeyP521, receiverPubKey *cryptotypes.PublicKeyP521,
) {
	// submit request and state update
	transferValue := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)
	testHelper.WaitMined(tx)
	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	//encrypt event payload with TEE private key and user public key
	encryptedMessage, err := crypto.Encrypt(senderPrivKey, receiverPubKey, []byte(message))
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res[0].ApplicationID, EncryptedData: encryptedMessage}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: ethCommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: common.NewBig(0)},
	}

	oldStateRoot := testHelper.GetStateRoot()

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID:  res[0].ApplicationID,
		RequestID:      res[0].RequestID,
		PrevStateRoot:  oldStateRoot,
		NewStateRoot:   [32]byte{0x04, 0x05, 0x06},
		Events:         events[:],
		Withdrawals:    withdrawals,
		Signature:      signature[:],
		RefundAmount:   common.NewBig(90),
		ApplicationFee: common.NewBig(10),
	}

	//complete state update
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)
}

func TestSubmitRequest(t *testing.T) {
	// mock private key for the client
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// Prepare a request
	protocolVersion := uint8(0)
	applicationId := common.NewApplicationId(1)
	requestType := common.Deploy
	payload := []byte("test-payload")
	depositAmount := big.NewInt(1)
	maxFeeValue := big.NewInt(100)

	// Submit the request
	requestId, blockNumber, err := blockchainClient.SubmitRequest(context.Background(), protocolVersion, applicationId, requestType, payload, depositAmount, maxFeeValue)
	require.NoError(t, err)
	// Get pending requests
	pending, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	//check block number > 0
	require.True(t, blockNumber > 0, "Returned block number shouldn't be 0")
	// Check that the submitted request is present and matches
	found := false
	// Convert types for comparison

	for _, r := range pending {
		if r.RequestID == requestId {
			found = true

			if r.ProtocolVersion != protocolVersion || r.ApplicationID != applicationId || r.RequestType != requestType || string(r.Payload) != string(payload) || r.DepositAmount.ToInt().Cmp(depositAmount) != 0 {
				t.Errorf(
					"Request fields do not match: got {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}, want {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}",
					r.ProtocolVersion, r.ApplicationID, r.RequestType, string(r.Payload), r.DepositAmount,
					protocolVersion, applicationId, requestType, string(payload), depositAmount,
				)
			}
		}
	}
	if !found {
		t.Errorf("Submitted request not found in pending requests")
	}

	err = blockchainClient.MarkRequestCompleted(context.Background(), requestId, big.NewInt(80), big.NewInt(20))
	require.NoError(t, err)
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

func TestGetRequestCompletedEvent(t *testing.T) {
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(0)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	receipt, err := testHelper.GetTxReceipt(tx)
	require.NoError(t, err)
	requestBlock := receipt.BlockNumber.Uint64()

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	err = blockchainClient.MarkRequestCompleted(context.Background(), res[0].RequestID, big.NewInt(80), big.NewInt(20))
	require.NoError(t, err)

	event, err := blockchainClient.GetRequestCompletedEvent(context.Background(), res[0].RequestID, requestBlock, requestBlock+1)
	require.Error(t, err)
	require.Nil(t, event)

	// First check where for sure there is not the event, ie in the block where the request was mined
	event, err = blockchainClient.GetRequestCompletedEvent(context.Background(), res[0].RequestID, requestBlock, 0)
	require.NoError(t, err)
	require.Nil(t, event, "RequestCompletedEvent shouldn't be found")

	// Check now from the tip, it should be found because the MarkRequestCompleted was successful
	event, err = blockchainClient.GetRequestCompletedEvent(context.Background(), res[0].RequestID, 0, requestBlock+1)
	require.NoError(t, err)
	require.NotNil(t, event, "RequestCompletedEvent should be found")

	require.True(t, event.Status == common.RequestResultOK)

	// Try with a failure
	transferValue = big.NewInt(1000000)
	tx = testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	receipt, err = testHelper.GetTxReceipt(tx)
	require.NoError(t, err)
	requestBlock = receipt.BlockNumber.Uint64()

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	failure := apperrors.New(
		apperrors.CodeSubmittingStateUpdateFailed,
		"test failure",
		nil,
	)

	err = blockchainClient.MarkRequestFailed(context.Background(), res[0].RequestID, failure)
	require.NoError(t, err)

	event, err = blockchainClient.GetRequestCompletedEvent(context.Background(), res[0].RequestID, 0, requestBlock+1)
	require.NoError(t, err)
	require.NotNil(t, event, "RequestCompletedEvent should be found")

	require.True(t, event.Status == common.RequestResultFailed)

	// Try with a state update
	tx = testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)
	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	receipt, err = testHelper.GetTxReceipt(tx)
	require.NoError(t, err)
	requestBlock = receipt.BlockNumber.Uint64()

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res[0].ApplicationID, EncryptedData: []byte{0x04, 0x05, 0x06}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: ethCommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: common.NewBig(10)},
	}

	oldStateRoot := testHelper.GetStateRoot()

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID:  res[0].ApplicationID,
		RequestID:      res[0].RequestID,
		PrevStateRoot:  oldStateRoot,
		NewStateRoot:   [32]byte{0x04, 0x05, 0x06},
		Events:         events[:],
		Withdrawals:    withdrawals,
		Signature:      signature[:],
		RefundAmount:   common.NewBig(90),
		ApplicationFee: common.NewBig(10),
	}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	event, err = blockchainClient.GetRequestCompletedEvent(context.Background(), res[0].RequestID, 0, requestBlock+1)
	require.NoError(t, err)
	require.NotNil(t, event, "RequestCompletedEvent should be found")

	require.True(t, event.Status == common.RequestResultOK)
}
