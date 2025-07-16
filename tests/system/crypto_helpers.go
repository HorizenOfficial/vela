package system

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
)

// CryptoHelper provides cryptographic operations for system tests
type CryptoHelper struct {
	userKeys map[string]*common.PrivateKeyP521
}

// NewCryptoHelper creates a new crypto helper
func NewCryptoHelper() *CryptoHelper {
	return &CryptoHelper{
		userKeys: make(map[string]*common.PrivateKeyP521),
	}
}

// GenerateUserKey generates a new private key for a user
func (c *CryptoHelper) GenerateUserKey(userID string) (*common.PrivateKeyP521, error) {
	privKey, err := crypto.GeneratePrivateKeyP521()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key for user %s: %w", userID, err)
	}

	c.userKeys[userID] = privKey
	return privKey, nil
}

// GetUserKey returns the private key for a user
func (c *CryptoHelper) GetUserKey(userID string) (*common.PrivateKeyP521, error) {
	key, exists := c.userKeys[userID]
	if !exists {
		return nil, fmt.Errorf("no key found for user %s", userID)
	}
	return key, nil
}

// CreateDepositRequest creates an encrypted deposit request
func (c *CryptoHelper) CreateDepositRequest(appID, requestID, sender string, value int64, receiverPubKey *common.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}

	// For deposit, payload can be empty since deposit is handled through the Value field
	payload := []byte{}

	// Encrypt payload (even empty payload needs to be encrypted)
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt deposit payload: %w", err)
	}

	return &common.Request{
		ApplicationID: appID,
		RequestID:     requestID,
		RequestType:   common.Process,
		Payload:       encryptedPayload,
		Sender:        sender,
		Timestamp:     time.Now().Unix(),
		Value:         value,
	}, nil
}

// CreateTransferRequest creates an encrypted transfer request
func (c *CryptoHelper) CreateTransferRequest(appID, requestID, sender, recipient string, amount int64, receiverPubKey *common.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}

	// Create transfer instruction
	transferInstruction := map[string]interface{}{
		"type": "transfer",
		"transfer": map[string]interface{}{
			"to":     recipient,
			"amount": amount,
		},
	}

	payload, err := json.Marshal(transferInstruction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer instruction: %w", err)
	}

	// Encrypt payload
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transfer payload: %w", err)
	}

	return &common.Request{
		ApplicationID: appID,
		RequestID:     requestID,
		RequestType:   common.Process,
		Payload:       encryptedPayload,
		Sender:        sender,
		Timestamp:     time.Now().Unix(),
		Value:         0, // No deposit for transfer
	}, nil
}

// CreateWithdrawalRequest creates an encrypted withdrawal request
func (c *CryptoHelper) CreateWithdrawalRequest(appID, requestID, sender, destinationAddress string, amount int64, receiverPubKey *common.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}

	// Create withdrawal instruction
	withdrawInstruction := map[string]interface{}{
		"type": "withdraw",
		"withdraw": map[string]interface{}{
			"to":     destinationAddress,
			"amount": amount,
		},
	}

	payload, err := json.Marshal(withdrawInstruction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal withdrawal instruction: %w", err)
	}

	// Encrypt payload
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt withdrawal payload: %w", err)
	}

	return &common.Request{
		ApplicationID: appID,
		RequestID:     requestID,
		RequestType:   common.Process,
		Payload:       encryptedPayload,
		Sender:        sender,
		Timestamp:     time.Now().Unix(),
		Value:         0, // No deposit for withdrawal
	}, nil
}

// CreateDeanonymizationRequest creates an encrypted deanonymization request
func (c *CryptoHelper) CreateDeanonymizationRequest(appID, requestID, sender string, receiverPubKey *common.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}

	// For deanonymization, payload can contain specific query parameters
	payload := []byte(`{"type":"deanonymization","query":"full_report"}`)

	// Encrypt payload
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt deanonymization payload: %w", err)
	}

	return &common.Request{
		ApplicationID: appID,
		RequestID:     requestID,
		RequestType:   common.Deanonymize,
		Payload:       encryptedPayload,
		Sender:        sender,
		Timestamp:     time.Now().Unix(),
		Value:         0,
	}, nil
}

// DecryptEvent decrypts an event using the user's private key
func (c *CryptoHelper) DecryptEvent(userID string, event *common.Event, senderPubKey *common.PublicKeyP521) ([]byte, error) {
	userKey, err := c.GetUserKey(userID)
	if err != nil {
		return nil, err
	}

	decryptedData, err := crypto.Decrypt(senderPubKey, userKey, event.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt event for user %s: %w", userID, err)
	}

	return decryptedData, nil
}

// DecryptDeanonymizationReport decrypts a deanonymization report
func (c *CryptoHelper) DecryptDeanonymizationReport(userID string, report *common.DeanonymizationReport, senderPubKey *common.PublicKeyP521) ([]byte, error) {
	userKey, err := c.GetUserKey(userID)
	if err != nil {
		return nil, err
	}

	decryptedReport, err := crypto.Decrypt(senderPubKey, userKey, report.EncryptedReport)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt deanonymization report for user %s: %w", userID, err)
	}

	return decryptedReport, nil
}
