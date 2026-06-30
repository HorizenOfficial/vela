package executor

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/HorizenOfficial/vela-common-go/subtypes"
)

// The seed-derivation constant, the anonymity-set size, and the deterministic
// subtype-set generator (GenerateSubtype / AllSubtypes) live in vela-common-go
// to deduplicate the same definitions across vela, vela-nova, and vela-ned.
// This file keeps only GenerateRandomSubtype because it depends on crypto/rand,
// which is host-only — the deterministic primitives belong in the shared
// package; the random pick lives next to its caller (encryptEvents).

// GenerateRandomSubtype picks a cryptographically random index in [1, n]
// and returns subtypes.GenerateSubtype(seed, index). Used by encryptEvents
// to rotate outbound events' EventSubType across the per-user anonymity set.
func GenerateRandomSubtype(seed []byte, n int) ([32]byte, error) {
	randVal, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to generate random index: %w", err)
	}
	index := int(randVal.Int64()) + 1 // shift [0, n) → [1, n]
	return subtypes.GenerateSubtype(seed, index), nil
}
