package app

// SubTypeFromString packs up to 32 bytes of s, left-aligned and zero-padded,
// into a [32]byte suitable for use as a bytes32 event subtype. Bytes beyond 32
// are silently dropped, so callers must ensure s fits.
//
// This encoding (ASCII truncation vs. keccak256) is a per-application policy:
// SimpleApp uses short readable tags so they remain visible in on-chain log
// topics. Other apps may prefer to hash their subtype names instead, and
// should define their own helper.
func SubTypeFromString(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}
