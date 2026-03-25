package executor

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

// SubtypeKeyMessage is the message the user signs with their secp256k1 key to produce a seed.
// Change this string to rotate all user subtype sets.
const SubtypeKeyMessage = "subtype-key-v1"

// DefaultSubtypeN is the number of possible subtypes per user (anonymity set size).
const DefaultSubtypeN = 50

// VerifySeed verifies that seed is a valid secp256k1 signature of
// keccak256(SubtypeKeyMessage) produced by signer.
// seed must be 65 bytes in [R || S || V] format with V in {0, 1}.
func VerifySeed(seed []byte, signer ethCommon.Address) error {
	if len(seed) != 65 {
		return fmt.Errorf("invalid seed length: expected 65, got %d", len(seed))
	}
	msgHash := ethCrypto.Keccak256([]byte(SubtypeKeyMessage))
	recoveredPub, err := ethCrypto.SigToPub(msgHash, seed)
	if err != nil {
		return fmt.Errorf("failed to recover public key from seed: %w", err)
	}
	recoveredAddr := ethCrypto.PubkeyToAddress(*recoveredPub)
	if recoveredAddr != signer {
		return fmt.Errorf("seed verification failed: recovered address %s does not match sender %s", recoveredAddr.Hex(), signer.Hex())
	}
	return nil
}

// GenerateSubtype returns "0x" + hex(HMAC-SHA256(key=seed, data=[]byte{index})).
// index should be in the range [1, N].
func GenerateSubtype(seed []byte, index int) string {
	mac := hmac.New(sha256.New, seed)
	mac.Write([]byte{byte(index)})
	return "0x" + hex.EncodeToString(mac.Sum(nil))
}

// AllSubtypes returns GenerateSubtype(seed, i) for i in 1..n.
// The returned slice has length n; index i maps to result[i-1].
func AllSubtypes(seed []byte, n int) []string {
	result := make([]string, n)
	for i := 1; i <= n; i++ {
		result[i-1] = GenerateSubtype(seed, i)
	}
	return result
}

// GenerateRandomSubtype picks a cryptographically random index in [1, n]
// and returns GenerateSubtype(seed, index).
func GenerateRandomSubtype(seed []byte, n int) (string, error) {
	randVal, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random index: %w", err)
	}
	index := int(randVal.Int64()) + 1 // shift [0, n) → [1, n]
	return GenerateSubtype(seed, index), nil
}
