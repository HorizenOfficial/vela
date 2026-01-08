package app

import (
	"encoding/binary"
	"encoding/json"
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
	return z[0] == 0 && z[1] == 0 && z[2] == 0 && z[3] == 0
}

// String returns the decimal representation of z.
func (z Uint256) String() string {
	if z.IsZero() {
		return "0"
	}

	val := z // copy
	var res []byte

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

// divModWord computes z / divisor and z % divisor.
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
// It marshals the Uint256 as a JSON number (no quotes), or string if large?
// math/big.Int marshals as digits.
func (z Uint256) MarshalJSON() ([]byte, error) {
	return []byte(z.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Only decimal strings or numbers are accepted. Overflow returns an error.
func (z *Uint256) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// If it's not a JSON string, try treating it as a raw JSON number
		s = string(data)
	}

	*z = Uint256{}

	const ten = uint64(10)

	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		}
		if c < '0' || c > '9' {
			return fmt.Errorf("invalid character in Uint256 string: %c", c)
		}
		digit := uint64(c - '0')

		if z.Mul64Overflow(ten) {
			return errors.New("Uint256 overflow")
		}
		if z.Add64Overflow(digit) {
			return errors.New("Uint256 overflow")
		}
	}
	return nil
}

// Mul64 sets z = z * y (mod 2^256).
func (z *Uint256) Mul64(y uint64) {
	var carry uint64
	for i := 0; i < 4; i++ {
		hi, lo := bits.Mul64(z[i], y)
		var c uint64
		z[i], c = bits.Add64(lo, carry, 0)
		carry = hi + c
	}
}

// Mul64Overflow sets z = z * y and reports whether overflow occurred.
func (z *Uint256) Mul64Overflow(y uint64) (overflow bool) {
	var carry uint64
	for i := 0; i < 4; i++ {
		hi, lo := bits.Mul64(z[i], y)
		var c uint64
		z[i], c = bits.Add64(lo, carry, 0)
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
