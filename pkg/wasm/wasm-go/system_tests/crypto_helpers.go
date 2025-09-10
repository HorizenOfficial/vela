package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/crypto"
)

type CryptoHelper struct {
	userKeys map[string]*cryptotypes.PrivateKeyP521
}

func NewCryptoHelper() *CryptoHelper {
	return &CryptoHelper{userKeys: make(map[string]*cryptotypes.PrivateKeyP521)}
}

func (c *CryptoHelper) GenerateUserKey(userID string) (*cryptotypes.PrivateKeyP521, error) {
	privKey, err := crypto.GeneratePrivateKeyP521()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key for user %s: %w", userID, err)
	}
	c.userKeys[userID] = privKey
	return privKey, nil
}

func (c *CryptoHelper) GetUserKey(userID string) (*cryptotypes.PrivateKeyP521, error) {
	key, ok := c.userKeys[userID]
	if !ok {
		return nil, fmt.Errorf("no key found for user %s", userID)
	}
	return key, nil
}

func (c *CryptoHelper) CreateDepositRequest(appID, requestID, sender string, value uint64, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}
	payload := []byte{}
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt deposit payload: %w", err)
	}
	return &common.Request{ApplicationID: appID, RequestID: requestID, RequestType: common.Process, Payload: encryptedPayload, Sender: sender, Timestamp: time.Now().Unix(), Value: value}, nil
}

func (c *CryptoHelper) CreateTransferRequest(appID, requestID, sender, recipient string, amount uint64, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}
	transferInstruction := map[string]interface{}{"type": "transfer", "transfer": map[string]interface{}{"to": recipient, "amount": amount}}
	payload, err := json.Marshal(transferInstruction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer instruction: %w", err)
	}
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transfer payload: %w", err)
	}
	return &common.Request{ApplicationID: appID, RequestID: requestID, RequestType: common.Process, Payload: encryptedPayload, Sender: sender, Timestamp: time.Now().Unix(), Value: 0}, nil
}

func (c *CryptoHelper) CreateWithdrawalRequest(appID, requestID, sender, destinationAddress string, amount uint64, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}
	withdrawInstruction := map[string]interface{}{"type": "withdraw", "withdraw": map[string]interface{}{"to": destinationAddress, "amount": amount}}
	payload, err := json.Marshal(withdrawInstruction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal withdrawal instruction: %w", err)
	}
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt withdrawal payload: %w", err)
	}
	return &common.Request{ApplicationID: appID, RequestID: requestID, RequestType: common.Process, Payload: encryptedPayload, Sender: sender, Timestamp: time.Now().Unix(), Value: 0}, nil
}

func (c *CryptoHelper) CreateDeanonymizationRequest(appID, requestID, sender string, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}
	payload := []byte(`{"type":"deanonymization","query":"full_report"}`)
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt deanonymization payload: %w", err)
	}
	return &common.Request{ApplicationID: appID, RequestID: requestID, RequestType: common.Deanonymize, Payload: encryptedPayload, Sender: sender, Timestamp: time.Now().Unix(), Value: 0}, nil
}

func (c *CryptoHelper) DecryptEvent(userID string, event *common.Event, senderPubKey *cryptotypes.PublicKeyP521) ([]byte, error) {
	userKey, err := c.GetUserKey(userID)
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(senderPubKey, userKey, event.EncryptedData)
}

func (c *CryptoHelper) DecryptDeanonymizationReport(userID string, report *common.DeanonymizationReport, senderPubKey *cryptotypes.PublicKeyP521) ([]byte, error) {
	userKey, err := c.GetUserKey(userID)
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(senderPubKey, userKey, report.EncryptedReport)
}

func (c *CryptoHelper) ValidateUpdatePayloadSignature(payload *common.UpdatePayload, key *cryptotypes.PublicKeySecp256k1) error {
	originalPayload := &common.UpdatePayload{ApplicationID: payload.ApplicationID, RequestID: payload.RequestID, PrevStateRoot: payload.PrevStateRoot, NewStateRoot: payload.NewStateRoot, Events: payload.Events, Withdrawals: payload.Withdrawals}
	payloadBytes, err := json.Marshal(originalPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for signing: %w", err)
	}
	hash := ethCrypto.Keccak256(payloadBytes)
	recoveredPubKey, err := ethCrypto.SigToPub(hash, payload.Signature)
	if err != nil {
		return fmt.Errorf("failed to recover public key: %v", err)
	}
	if !bytes.Equal(ethCrypto.FromECDSAPub(key.PublicKey), ethCrypto.FromECDSAPub(recoveredPubKey)) {
		return fmt.Errorf("recovered public key does not match original key")
	}
	return nil
}
