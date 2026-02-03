package common

import (
	"fmt"
	"math/big"
)

// Big is a wrapper around big.Int that marshals/unmarshals as a hex string with 0x prefix.
// Negative values are not supported as Ethereum hex strings are unsigned.
// Zero values are represented as "0x0".
type Big big.Int

func (b *Big) ToInt() *big.Int {
	return (*big.Int)(b)
}

func ToBig(i *big.Int) *Big {
	return (*Big)(i)
}

// NewBig creates a new Big from a uint64 value.
// This is a convenience constructor for common use cases.
// Uses uint64 since Big does not support negative values.
func NewBig(x uint64) *Big {
	return ToBig(new(big.Int).SetUint64(x))
}

func (b *Big) String() string {
	if b == nil {
		return "<nil>"
	}
	return (*big.Int)(b).String()
}

func (b *Big) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	bi := (*big.Int)(b)
	// Reject negative values - Ethereum hex strings are unsigned
	if bi.Sign() < 0 {
		return nil, fmt.Errorf("cannot marshal negative Big value: %s", bi.String())
	}
	hexStr := bi.Text(16)
	// Pre-allocate: 2 quotes + "0x" + hex string = 4 + len(hexStr)
	result := make([]byte, 0, len(hexStr)+4)
	result = append(result, '"', '0', 'x')
	result = append(result, hexStr...)
	result = append(result, '"')
	return result, nil
}

func (b *Big) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("cannot unmarshal into nil *Big")
	}
	// Handle null case - set to zero value
	if len(data) == 4 && data[0] == 'n' && data[1] == 'u' && data[2] == 'l' && data[3] == 'l' {
		*b = Big{}
		return nil
	}
	// Must be a quoted string with at least "0x" + 1 hex char = 5 bytes minimum
	if len(data) < 5 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid Big format: expected quoted string, got %q", string(data))
	}
	// Check for 0x prefix directly on bytes (avoid string conversion)
	if data[1] != '0' || data[2] != 'x' {
		return fmt.Errorf("invalid Big format: missing 0x prefix")
	}
	// Extract hex string (string conversion needed for SetString)
	hexBytes := data[3 : len(data)-1]
	var i big.Int
	if _, ok := i.SetString(string(hexBytes), 16); !ok {
		return fmt.Errorf("invalid Big hex string: %q", string(hexBytes))
	}
	*b = Big(i)
	return nil
}

