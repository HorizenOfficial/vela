package executor

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// SubtypeKeyMessage is the message used to derive a seed.
// Change this string to rotate all user subtype sets.
const SubtypeKeyMessage = "subtype-key-v1"

// DefaultSubtypeN is the number of possible subtypes per user (anonymity set size).
const DefaultSubtypeN = 50

// GenerateSubtype returns HMAC-SHA256(key=seed, data=[]byte{index}) as a 32-byte
// value, matching the on-chain bytes32 event subtype.
// index should be in the range [1, N].
func GenerateSubtype(seed []byte, index int) [32]byte {
	mac := hmac.New(sha256.New, seed)
	mac.Write([]byte{byte(index)})
	var out [32]byte
	copy(out[:], mac.Sum(nil))
	return out
}

// AllSubtypes returns GenerateSubtype(seed, i) for i in 1..n.
// The returned slice has length n; index i maps to result[i-1].
func AllSubtypes(seed []byte, n int) [][32]byte {
	result := make([][32]byte, n)
	for i := 1; i <= n; i++ {
		result[i-1] = GenerateSubtype(seed, i)
	}
	return result
}

// GenerateRandomSubtype picks a cryptographically random index in [1, n]
// and returns GenerateSubtype(seed, index).
func GenerateRandomSubtype(seed []byte, n int) ([32]byte, error) {
	randVal, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to generate random index: %w", err)
	}
	index := int(randVal.Int64()) + 1 // shift [0, n) → [1, n]
	return GenerateSubtype(seed, index), nil
}
