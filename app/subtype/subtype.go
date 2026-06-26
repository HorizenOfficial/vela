// Package subtype is a small helper shared by the sample WASM guest apps for
// turning a human-readable event name into a bytes32 event subtype.
package subtype

// FromString packs up to 32 bytes of s, left-aligned and zero-padded, into a
// [32]byte suitable for use as a bytes32 event subtype. Bytes beyond 32 are
// silently dropped, so callers must ensure s fits.
//
// This encoding (ASCII truncation vs. keccak256) is a per-application policy: the
// sample apps use short readable tags so they remain visible in on-chain log
// topics. Apps that prefer hashed subtype names should define their own helper
// instead of using this one.
func FromString(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}
