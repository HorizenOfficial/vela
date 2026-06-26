package app

// SubTypeFromString packs up to 32 bytes of s, left-aligned and zero-padded,
// into a [32]byte suitable for use as a bytes32 event subtype. Bytes beyond 32
// are silently dropped, so callers must ensure s fits.
func SubTypeFromString(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}
