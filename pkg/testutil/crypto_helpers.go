package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/HorizenOfficial/vela/pkg/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/HorizenOfficial/vela/pkg/executor"
)

// CryptoHelper provides cryptographic operations for system tests
type CryptoHelper struct {
	userKeys        map[ethCommon.Address]*cryptotypes.PrivateKeyP521
	userSigningKeys map[ethCommon.Address]*cryptotypes.PrivateKeySecp256k1
}

// NewCryptoHelper creates a new crypto helper
func NewCryptoHelper() *CryptoHelper {
	return &CryptoHelper{
		userKeys:        make(map[ethCommon.Address]*cryptotypes.PrivateKeyP521),
		userSigningKeys: make(map[ethCommon.Address]*cryptotypes.PrivateKeySecp256k1),
	}
}

// GenerateUserKey generates a new private key for a user
func (c *CryptoHelper) GenerateUserKey(userID ethCommon.Address) (*cryptotypes.PrivateKeyP521, error) {
	privKey, err := crypto.GeneratePrivateKeyP521()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key for user %s: %w", userID, err)
	}

	c.userKeys[userID] = privKey
	return privKey, nil
}

// GetUserKey returns the private key for a user
func (c *CryptoHelper) GetUserKey(userID ethCommon.Address) (*cryptotypes.PrivateKeyP521, error) {
	key, exists := c.userKeys[userID]
	if !exists {
		return nil, fmt.Errorf("no key found for user %s", userID)
	}
	return key, nil
}

// GenerateUserSigningKey generates a secp256k1 signing key for the user and stores it.
// This key is used to compute the seed for privacy-preserving event subtypes.
func (c *CryptoHelper) GenerateUserSigningKey(userID ethCommon.Address) (*cryptotypes.PrivateKeySecp256k1, error) {
	privKey, err := crypto.GeneratePrivateKeySecp256k1()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secp256k1 key for user %s: %w", userID, err)
	}
	c.userSigningKeys[userID] = privKey
	return privKey, nil
}

// GenerateUserIdentity generates a secp256k1 signing key and returns the Ethereum address
// derived from it. The user address MUST match the signing key for VerifySeed to pass.
// Use this instead of hardcoded addresses when the user will call AssociateKey.
func (c *CryptoHelper) GenerateUserIdentity() (ethCommon.Address, error) {
	privKey, err := crypto.GeneratePrivateKeySecp256k1()
	if err != nil {
		return ethCommon.Address{}, fmt.Errorf("failed to generate secp256k1 identity key: %w", err)
	}
	addr := ethCommon.HexToAddress(privKey.PublicKey().Address())
	c.userSigningKeys[addr] = privKey
	return addr, nil
}

// GetUserSigningKey returns the secp256k1 signing key for a user.
func (c *CryptoHelper) GetUserSigningKey(userID ethCommon.Address) (*cryptotypes.PrivateKeySecp256k1, error) {
	key, exists := c.userSigningKeys[userID]
	if !exists {
		return nil, fmt.Errorf("no secp256k1 signing key found for user %s", userID)
	}
	return key, nil
}

// ComputeSeed signs keccak256(SubtypeKeyMessage) with the user's secp256k1 key.
// The resulting 65-byte signature is the seed used to derive privacy-preserving event subtypes.
// GenerateUserSigningKey must be called before ComputeSeed.
func (c *CryptoHelper) ComputeSeed(userID ethCommon.Address) ([]byte, error) {
	signingKey, err := c.GetUserSigningKey(userID)
	if err != nil {
		return nil, err
	}
	msgHash := ethCrypto.Keccak256([]byte(executor.SubtypeKeyMessage))
	seed, err := signingKey.Sign(msgHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign subtype message for user %s: %w", userID, err)
	}
	return seed, nil
}

// CreateAssociateKeyRequest creates an associate key request.
// The payload contains the P521 public key (133 bytes) followed by the encrypted seed (93 bytes).
// The seed is encrypted with ECDH(user_priv_P521, enclave_pub_P521).
// GenerateUserSigningKey must be called for sender before this method.
func (c *CryptoHelper) CreateAssociateKeyRequest(appID common.ApplicationIdType, requestID common.RequestIdType, sender ethCommon.Address, key *cryptotypes.PublicKeyP521, enclavePubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
	seed, err := c.ComputeSeed(sender)
	if err != nil {
		return nil, fmt.Errorf("failed to compute seed for user %s: %w", sender, err)
	}
	// Encrypt the seed with ECDH(user_priv_P521, enclave_pub_P521)
	userPrivKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, fmt.Errorf("failed to get user P521 key for %s: %w", sender, err)
	}
	encryptedSeed, err := crypto.Encrypt(userPrivKey, enclavePubKey, seed)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt seed for user %s: %w", sender, err)
	}
	// payload = 133-byte P521 pubkey + 93-byte encrypted seed
	payload := append(key.Bytes(), encryptedSeed...)

	return &common.Request{
		ApplicationID: appID,
		RequestID:     requestID,
		RequestType:   common.AssociateKey,
		Payload:       payload,
		Sender:        sender,
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}, nil
}

// CreateDepositRequest creates an encrypted deposit request
func (c *CryptoHelper) CreateDepositRequest(appID common.ApplicationIdType, requestID common.RequestIdType, sender ethCommon.Address, depositAmount *big.Int, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
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
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.ToBig(depositAmount),
		MaxFeeValue:   common.NewBig(100),
	}, nil
}

// CreateWithdrawalRequest creates an encrypted withdrawal request
func (c *CryptoHelper) CreateWithdrawalRequest(appID common.ApplicationIdType, requestID common.RequestIdType, sender, destinationAddress ethCommon.Address, amount *common.Big, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
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
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0), // No deposit for withdrawal
		MaxFeeValue:   common.NewBig(100),
	}, nil
}

// CreateDeanonymizationRequest creates an encrypted deanonymization request
func (c *CryptoHelper) CreateDeanonymizationRequest(appID common.ApplicationIdType, requestID common.RequestIdType, sender ethCommon.Address, payload []byte, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}

	// For deanonymization, payload can contain specific query parameters
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
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}, nil
}

// DecryptEvent decrypts an event using the user's private key
func (c *CryptoHelper) DecryptEvent(userID ethCommon.Address, event *common.Event, senderPubKey *cryptotypes.PublicKeyP521) ([]byte, error) {
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
func (c *CryptoHelper) DecryptDeanonymizationReport(userID ethCommon.Address, report *common.DeanonymizationReport, senderPubKey *cryptotypes.PublicKeyP521) ([]byte, error) {
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

// CreateProcessRequest creates an encrypted process request with a raw payload
func (c *CryptoHelper) CreateProcessRequest(appID common.ApplicationIdType, requestID common.RequestIdType, sender ethCommon.Address, payload []byte, receiverPubKey *cryptotypes.PublicKeyP521) (*common.Request, error) {
	senderKey, err := c.GetUserKey(sender)
	if err != nil {
		return nil, err
	}

	// Encrypt payload
	encryptedPayload, err := crypto.Encrypt(senderKey, receiverPubKey, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt process payload: %w", err)
	}

	return &common.Request{
		ApplicationID: appID,
		RequestID:     requestID,
		RequestType:   common.Process,
		Payload:       encryptedPayload,
		Sender:        sender,
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		DepositAmount: common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}, nil
}

func (c *CryptoHelper) ValidateUpdatePayloadSignature(payload *common.UpdatePayload, key *cryptotypes.PublicKeySecp256k1) error {
	// Create the original payload that was signed
	originalPayload := &common.UpdatePayload{
		ApplicationID:  payload.ApplicationID,
		RequestID:      payload.RequestID,
		PrevStateRoot:  payload.PrevStateRoot,
		NewStateRoot:   payload.NewStateRoot,
		Events:         payload.Events,
		Withdrawals:    payload.Withdrawals,
		RefundAmount:   payload.RefundAmount,
		ApplicationFee: payload.ApplicationFee,
	}

	msgBuilder, err := executor.NewMsgToSignBuilder()
	if err != nil {
		return fmt.Errorf("failed to create message to sign builder: %w", err)
	}

	msg, err := msgBuilder.BuildMsgHash(originalPayload)
	if err != nil {
		return fmt.Errorf("failed to build message to sign: %w", err)
	}

	payload.Signature[64] -= 27 //SigToPub requires v field < 4
	// Recover the public key from the signature
	recoveredPubKey, err := ethCrypto.SigToPub(msg, payload.Signature)
	if err != nil {
		return fmt.Errorf("failed to recover public key: %w", err)
	}

	// The recovered public key should match the original public key
	if !bytes.Equal(ethCrypto.FromECDSAPub(key.PublicKey), ethCrypto.FromECDSAPub(recoveredPubKey)) {
		return fmt.Errorf("recovered public key does not match original key")
	}

	return nil
}
