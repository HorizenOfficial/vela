package app

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/bits"
	"strings"
)

// Uint256 represents a 256-bit unsigned integer using 4 uint64 values.
// Words are stored in little-endian order (words[0] is least-significant).
//
// All arithmetic operations are performed modulo 2^256 unless stated otherwise.
type Uint256 [4]uint64

// NewUint256 creates a new Uint256 from a uint64 value.
func NewUint256(v uint64) *Uint256 {
	return &Uint256{v, 0, 0, 0}
}

// SetBytes interprets bytes as a big-endian unsigned integer.
// Values larger than 256 bits are truncated (mod 2^256).
// This matches modulo arithmetic semantics.
func (z *Uint256) SetBytes(b []byte) *Uint256 {
	*z = Uint256{}
	if len(b) == 0 {
		return z
	}
	if len(b) > 32 {
		b = b[len(b)-32:]
	}

	var tmp [32]byte
	copy(tmp[32-len(b):], b)

	z[3] = binary.BigEndian.Uint64(tmp[0:8])
	z[2] = binary.BigEndian.Uint64(tmp[8:16])
	z[1] = binary.BigEndian.Uint64(tmp[16:24])
	z[0] = binary.BigEndian.Uint64(tmp[24:32])

	return z
}

// Add sets z = x + y (mod 2^256) and returns z.
func (z *Uint256) Add(x, y Uint256) *Uint256 {
	var carry uint64
	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], _ = bits.Add64(x[3], y[3], carry)
	return z
}

// AddOverflow sets z = x + y and reports whether overflow occurred.
func (z *Uint256) AddOverflow(x, y Uint256) (overflow bool) {
	var carry uint64
	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], carry = bits.Add64(x[3], y[3], carry)
	return carry != 0
}

// Sub sets z = x - y (mod 2^256) and returns z.
func (z *Uint256) Sub(x, y Uint256) *Uint256 {
	var borrow uint64
	z[0], borrow = bits.Sub64(x[0], y[0], 0)
	z[1], borrow = bits.Sub64(x[1], y[1], borrow)
	z[2], borrow = bits.Sub64(x[2], y[2], borrow)
	z[3], _ = bits.Sub64(x[3], y[3], borrow)
	return z
}

// Cmp compares z and y and returns:
//
//	-1 if z < y
//	 0 if z == y
//	+1 if z > y
func (z Uint256) Cmp(y Uint256) int {
	for i := 3; ; i-- {
		if z[i] > y[i] {
			return 1
		}
		if z[i] < y[i] {
			return -1
		}
		if i == 0 {
			break
		}
	}
	return 0
}

// IsZero returns true if z == 0.
func (z Uint256) IsZero() bool {
	return (z[0] | z[1] | z[2] | z[3]) == 0
}

// String returns the decimal representation of z.
func (z Uint256) String() string {
	if z.IsZero() {
		return "0"
	}

	val := z // copy
	res := make([]byte, 0, 78) // Max digits for 2^256

	const ten = uint64(10)
	for !val.IsZero() {
		var rem uint64
		val, rem = val.divModWord(ten)
		res = append(res, byte(rem)+'0')
	}

	// reverse
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return string(res)
}

// divModWord computes z / divisor and z % divisor. We assume the divisor has been checked by the caller to be != 0
func (z Uint256) divModWord(divisor uint64) (Uint256, uint64) {
	var quot Uint256
	var r uint64

	for i := 3; ; i-- {
		q, rr := bits.Div64(r, z[i], divisor)
		quot[i] = q
		r = rr
		if i == 0 {
			break
		}
	}
	return quot, r
}

// MarshalJSON implements json.Marshaler.
// It marshals the Uint256 as a JSON number (no quotes)
// math/big.Int marshals as digits.
func (z Uint256) MarshalJSON() ([]byte, error) {
	return []byte(z.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Only decimal strings or numbers are accepted. Overflow returns an error.
func (z *Uint256) UnmarshalJSON(data []byte) error {
	var s string

	// 1. Determine if input is a JSON string or raw number by peeking the first non-whitespace char using a recursive for loop
	trimmedData := data
	for len(trimmedData) > 0 && (trimmedData[0] == ' ' || trimmedData[0] == '\t' || trimmedData[0] == '\n' || trimmedData[0] == '\r') {
		trimmedData = trimmedData[1:]
	}

	if len(trimmedData) > 0 && trimmedData[0] == '"' {
		// First non-whitespace character is a quote, this is a JSON string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
	} else {
		// First non-whitespace character is not a quote, this is a raw JSON number or we have empty input
		s = string(trimmedData)
	}

	// 2. Clean up surrounding whitespace and validate non-empty
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return errors.New("Uint256 value is empty")
	}

	// 3. Strict digit-only parsing
	*z = Uint256{}
	const ten = uint64(10)
	for _, c := range s {
		if c < '0' || c > '9' {
			return fmt.Errorf("invalid character in Uint256: %c", c)
		}
		digit := uint64(c - '0')

		if z.Mul64Overflow(ten) {
			// we have modified z actually, but caller must check the error
			return errors.New("Uint256 overflow after multiplication")
		}
		if z.Add64Overflow(digit) {
    		// we have modified z actually, but caller must check the error
			return errors.New("Uint256 overflow after sum")
		}
	}
	return nil
}

// Mul64 sets z = z * y (mod 2^256).
func (z *Uint256) Mul64(y uint64) {
	_ = z.Mul64Overflow(y)
}

// Mul64Overflow sets z = z * y and reports whether overflow occurred.
func (z *Uint256) Mul64Overflow(y uint64) (overflow bool) {
	var carry uint64
	for i := range 4 {
		hi, lo := bits.Mul64(z[i], y)
		var c uint64
		z[i], c = bits.Add64(lo, carry, 0)
		// Note: hi + c cannot overflow uint64 here because bits.Mul64(a, b)
		// has a maximum 'hi' value of 0xfffffffffffffffe (MaxUint64 - 1).
		// Since c is at most 1, the sum hi + c is at most MaxUint64.
		carry = hi + c
	}
	return carry != 0
}

// Add64 adds y to z (mod 2^256).
func (z *Uint256) Add64(y uint64) {
	var carry uint64
	z[0], carry = bits.Add64(z[0], y, 0)
	for i := 1; i < 4 && carry != 0; i++ {
		z[i], carry = bits.Add64(z[i], 0, carry)
	}
}

// Add64Overflow adds y to z and reports whether overflow occurred.
func (z *Uint256) Add64Overflow(y uint64) (overflow bool) {
	var carry uint64
	z[0], carry = bits.Add64(z[0], y, 0)
	for i := 1; i < 4 && carry != 0; i++ {
		z[i], carry = bits.Add64(z[i], 0, carry)
	}
	return carry != 0
}
