package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/stretchr/testify/require"
)

//go:generate mkdir -p ./contracts/mocktee
//go:generate solc --combined-json abi,bin ../../contracts/contracts/mocks/MockTeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/MockTeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/MockTeeAuthenticatorAbi/combined.json --pkg mocktee --type MockTeeAuthenticator --out ./contracts/mocktee/MockTeeAuthenticator.go
//go:generate mkdir -p ./contracts/tee
//go:generate solc --combined-json abi,bin ../../contracts/contracts/TeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/TeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/TeeAuthenticatorAbi/combined.json --pkg tee --type TeeAuthenticator --out ./contracts/tee/TeeAuthenticator.go
//go:generate mkdir -p ./contracts/authority
//go:generate solc --combined-json abi,bin ../../contracts/contracts/AuthorityRegistry.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/AuthorityRegistryAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/AuthorityRegistryAbi/combined.json --pkg authority --type AuthorityRegistry --out ./contracts/authority/AuthorityRegistry.go

var (
	applicationId = big.NewInt(1)
)

func SetupNewBlockChainClient(testHelper *testutil.SimTestHelper) *BlockChainClient {
	blockchainClient := NewBlockChainClient(testHelper.ProcessorContractAddress, testHelper.TeeSignerAddress, "", nil)
	blockchainClient.client = testHelper.Client()

	blockchainClient.processorBoundContract = blockchainClient.processorEndpoint.Instance(blockchainClient.client, testHelper.ProcessorContractAddress)
	blockchainClient.teeAuthBoundContract = blockchainClient.teeAuthEndpoint.Instance(blockchainClient.client, testHelper.TeeSignerAddress)

	blockchainClient.account = testHelper.ManagerAccount
	blockchainClient.connected = true

	return blockchainClient

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

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1203055)
	payload := ethCommon.FromHex("0x001234")
	tx := testHelper.SubmitRequest(applicationId, common.Process, payload, transferValue)

	testHelper.MineBlock()

	// wait for transaction inclusion
	testHelper.WaitMined(tx)
	fmt.Println("Request was successfully included")

	res, err = blockchainClient.GetPendingRequests(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1, len(res), "There should be one pending request")

	request := res[0]
	require.Equal(t, strconv.Itoa(int(testHelper.ProtocolVersion)), request.ProtocolVersion, "Protocol version should match")
	require.Equal(t, applicationId.String(), request.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, request.RequestType, "Request type should match")
	require.Equal(t, payload, request.Payload, "Payload should match")
	require.Greater(t, request.Timestamp, int64(0), "Timestamp should match")

	require.Equal(t, testHelper.Submitter.From.String(), request.Sender, "Sender should match")
	require.Equal(t, transferValue.Uint64(), request.Value, "Value should match")

}

func TestMarkRequestCompleted(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(0)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	err = blockchainClient.MarkRequestCompleted(context.Background(), res[0].RequestID)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")

}

func TestMarkRequestFailed(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	err = blockchainClient.MarkRequestFailed(context.Background(), res[0].RequestID)
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
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res[0].ApplicationID, EncryptedData: []byte{0x04, 0x05, 0x06}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: "0x1234567890123456789012345678901234567890", Amount: 10},
	}

	oldStateRoot := testHelper.GetStateRoot()

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID: res[0].ApplicationID,
		RequestID:     res[0].RequestID,
		PrevStateRoot: oldStateRoot,
		NewStateRoot:  [32]byte{0x04, 0x05, 0x06},
		Events:        events[:],
		Withdrawals:   withdrawals,
		Signature:     signature[:],
	}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be 0 pending request")
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
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageSkipped, teeKey, userPub);
	message := "test message"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message, teeKey, userPub);

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, nil, true)
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
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message1, teeKey, userPub);
	message2 := "test message 2"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message2, teeKey, userPub);
	message3 := "test message 3"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message3, teeKey, userPub);

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, nil, false)
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
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageTrue, teeKey, userPub);
	messageFalse := "test message - false"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageFalse, teeKey, userPub);

	//filter function
	filter := func(data []byte) bool {
    	return strings.Contains(string(data), "true")
	}

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, filter, false)
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
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message, teeKey, userPub);
	messageOther := "test message - for other user"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageOther, teeKey, otherUserPub);

	//retrieve and decrypt user events
	// for user
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(userEvents), "There should be 1 user event")
	require.Equal(t, []byte(message), userEvents[0], "Decrypted message should match original")

	// for other user
	otherUserEvents, err := blockchainClient.GetUserEvents(context.Background(), *otherUserKey, *applicationId, 0, 0, nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(otherUserEvents), "There should be 1 user event (other)")
	require.Equal(t, []byte(messageOther), otherUserEvents[0], "Decrypted message should match original (other)")
}

func _submitRequestAndStateUpdateWithEncryptedMessageEvent(t *testing.T, blockchainClient *BlockChainClient, testHelper *testutil.SimTestHelper, 
	message string, senderPrivKey *cryptotypes.PrivateKeyP521, receiverPubKey *cryptotypes.PublicKeyP521,
) {
	// submit request and state update
	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)
	testHelper.WaitMined(tx)
	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	//encrypt event payload with TEE private key and user public key
	encryptedMessage, err := crypto.Encrypt(senderPrivKey, receiverPubKey, []byte(message))
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res[0].ApplicationID, EncryptedData: encryptedMessage}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: "0x1234567890123456789012345678901234567890", Amount: 0},
	}

	oldStateRoot := testHelper.GetStateRoot()

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID: res[0].ApplicationID,
		RequestID:     res[0].RequestID,
		PrevStateRoot: oldStateRoot,
		NewStateRoot:  [32]byte{0x04, 0x05, 0x06},
		Events:        events[:],
		Withdrawals:   withdrawals,
		Signature:     signature[:],
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
	req := &common.Request{
		ProtocolVersion: "0",
		ApplicationID:   "1",
		RequestID:       "",
		RequestType:     common.Deploy,
		Payload:         []byte("test-payload"),
		Timestamp:       time.Now().Unix(),
		Sender:          "", 
		Value:           1,
	}

	// Submit the request
	requestId, err := blockchainClient.SubmitRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}

	// Get pending requests
	pending, err := blockchainClient.GetPendingRequests(context.Background())
	if err != nil {
		t.Fatalf("GetPendingRequests failed: %v", err)
	}

	// Check that the submitted request is present and matches
	found := false
	for _, r := range pending {
		reqIdAsString, err := common.RequestIdStringTo32Byte(r.RequestID)
		require.NoError(t, err)
		if reqIdAsString == requestId {
			found = true
			if r.ProtocolVersion != req.ProtocolVersion || r.ApplicationID != req.ApplicationID || r.RequestType != req.RequestType || string(r.Payload) != string(req.Payload) || r.Value != req.Value {
				t.Errorf("Request fields do not match: got %+v, want %+v", r, req)
			}
		}
	}
	if !found {
		t.Errorf("Submitted request not found in pending requests")
	}
}