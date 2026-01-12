package app

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
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

// ToHex returns the hex representation of z with "0x" prefix.
func (z Uint256) ToHex() string {
	if z.IsZero() {
		return "0x0"
	}
	// Similar logic to String() but base 16
	val := z
	res := make([]byte, 0, 66) // 0x + 64 hex digits
	
	const sixteen = uint64(16)
	const hexChars = "0123456789abcdef"

	for !val.IsZero() {
		var rem uint64
		val, rem = val.divModWord(sixteen)
		res = append(res, hexChars[rem])
	}
	res = append(res, 'x', '0')

	// reverse
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return string(res)
}

// MarshalJSON implements json.Marshaler.
// It marshals the Uint256 as a hex string with 0x prefix.
func (z Uint256) MarshalJSON() ([]byte, error) {
	s := z.ToHex()
	buf := make([]byte, 0, len(s)+2)
	buf = append(buf, '"')
	buf = append(buf, s...)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Only hex strings with "0x" prefix are accepted.
func (z *Uint256) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*z = Uint256{}
		return nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid Uint256 format: %s", string(data))
	}
	s := string(data[1 : len(data)-1])

	if len(s) < 2 || (s[0] != '0' || (s[1] != 'x' && s[1] != 'X')) {
		return fmt.Errorf("invalid Uint256 prefix: %s", s)
	}

	*z = Uint256{}
	s = s[2:]
	if len(s) == 0 {
		// "0x" is treated as 0
		return nil
	}

	const sixteen = uint64(16)
	for _, c := range s {
		var digit uint64
		switch {
		case c >= '0' && c <= '9':
			digit = uint64(c - '0')
		case c >= 'a' && c <= 'f':
			digit = uint64(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			digit = uint64(c - 'A' + 10)
		default:
			return fmt.Errorf("invalid hex character in Uint256: %c", c)
		}

		if z.Mul64Overflow(sixteen) {
			return errors.New("Uint256 overflow parsing hex")
		}
		if z.Add64Overflow(digit) {
			return errors.New("Uint256 overflow parsing hex")
		}
	}
	return nil
}

// Mul64 sets z = z * y (mod 2^256).
func (z *Uint256) Mul64(y uint64) *Uint256 {
	_ = z.Mul64Overflow(y)
	return z
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
func (z *Uint256) Add64(y uint64) *Uint256 {
	var carry uint64
	z[0], carry = bits.Add64(z[0], y, 0)
	for i := 1; i < 4 && carry != 0; i++ {
		z[i], carry = bits.Add64(z[i], 0, carry)
	}
	return z
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
