package app

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewUint256(t *testing.T) {
	u := NewUint256(12345)
	require.Equal(t, uint64(12345), u[0])
	require.Equal(t, uint64(0), u[1])
	require.Equal(t, uint64(0), u[2])
	require.Equal(t, uint64(0), u[3])
}

func TestSetBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string // decimal string
	}{
		{"empty", []byte{}, "0"},
		{"zero", []byte{0}, "0"},
		{"one byte", []byte{255}, "255"},
		{"two bytes", []byte{1, 0}, "256"},
		{"max uint64", new(big.Int).SetUint64(^uint64(0)).Bytes(), "18446744073709551615"},
		{"larger than uint64", new(big.Int).Mul(big.NewInt(1), new(big.Int).Lsh(big.NewInt(1), 64)).Bytes(), "18446744073709551616"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := new(Uint256).SetBytes(tt.input)
			require.Equal(t, tt.expected, u.String())
		})
	}
}

func TestAdd(t *testing.T) {
	// Compare against big.Int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 1000; i++ {
		// Generate two random big ints that fit in 256 bits
		b1 := new(big.Int).Rand(r, new(big.Int).Lsh(big.NewInt(1), 256))
		b2 := new(big.Int).Rand(r, new(big.Int).Lsh(big.NewInt(1), 256))

		// Convert to Uint256
		u1 := new(Uint256).SetBytes(b1.Bytes())
		u2 := new(Uint256).SetBytes(b2.Bytes())

		// Add
		sumBig := new(big.Int).Add(b1, b2)
		// Mask to 256 bits to simulate wrap-around behavior if overflow occurs
		mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
		sumBig.And(sumBig, mask)

		sumU := new(Uint256).Add(u1, u2)

		require.Equal(t, sumBig.String(), sumU.String(), "Add mismatch: %v + %v", b1, b2)
	}
}

func TestSub(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 1000; i++ {
		b1 := new(big.Int).Rand(r, new(big.Int).Lsh(big.NewInt(1), 256))
		b2 := new(big.Int).Rand(r, new(big.Int).Lsh(big.NewInt(1), 256))

		// Ensure b1 >= b2 for simple subtraction test (negative result handling depends on implementation,
		// usually standard uint wraparound)
		if b1.Cmp(b2) < 0 {
			b1, b2 = b2, b1
		}

		u1 := new(Uint256).SetBytes(b1.Bytes())
		u2 := new(Uint256).SetBytes(b2.Bytes())

		diffBig := new(big.Int).Sub(b1, b2)
		diffU := new(Uint256).Sub(u1, u2)

		require.Equal(t, diffBig.String(), diffU.String(), "Sub mismatch: %v - %v", b1, b2)
	}
}

func TestCmp(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"0", "0", 0},
		{"1", "1", 0},
		{"0", "1", -1},
		{"1", "0", 1},
		{"18446744073709551615", "18446744073709551615", 0}, // Max Uint64
		{"18446744073709551616", "18446744073709551615", 1}, // Max Uint64 + 1 vs Max Uint64
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			bA, _ := new(big.Int).SetString(tt.a, 10)
			bB, _ := new(big.Int).SetString(tt.b, 10)

			uA := new(Uint256).SetBytes(bA.Bytes())
			uB := new(Uint256).SetBytes(bB.Bytes())

			require.Equal(t, tt.expected, uA.Cmp(uB))
		})
	}
}

func TestString(t *testing.T) {
	tests := []string{
		"0",
		"1",
		"123456789",
		"18446744073709551615",                                     // Max Uint64
		"340282366920938463463374607431768211455",                   // Max Uint128
		"115792089237316195423570985008687907853269984665640564039457584007913129639935", // Max Uint256
	}

	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			b, _ := new(big.Int).SetString(s, 10)
			u := new(Uint256).SetBytes(b.Bytes())
			require.Equal(t, s, u.String())
		})
	}
}

func TestJSON(t *testing.T) {
	type wrapper struct {
		Val *Uint256 `json:"val"`
	}

	tests := []string{
		"0",
		"100",
		"12345678901234567890",
	}

	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			// Test Marshal
			b, _ := new(big.Int).SetString(s, 10)
			u := new(Uint256).SetBytes(b.Bytes())
			w := wrapper{Val: u}

			data, err := json.Marshal(w)
			require.NoError(t, err)

			// Expect plain number or string depending on implementation
			// Current implementation marshals as string digits without quotes if called directly,
			// but wrapper with `json:"val"` and MarshalJSON returning []byte(string) usually results in
			// the raw characters being inserted into the JSON.
			// However, our MarshalJSON returns []byte(z.String()), which are just the digits.
			// If we put that into a JSON object value position, it acts as a JSON Number.
			// Let's verify what comes out.
			require.Contains(t, string(data), s)

			// Test Unmarshal
			var w2 wrapper
			err = json.Unmarshal(data, &w2)
			require.NoError(t, err)
			require.Equal(t, s, w2.Val.String())
		})
	}

	// Test unmarshal from string (quoted)
	t.Run("quoted string", func(t *testing.T) {
		jsonStr := `{"val": "12345"}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.NoError(t, err)
		require.Equal(t, "12345", w.Val.String())
	})
}

func TestIsZero(t *testing.T) {
	require.True(t, NewUint256(0).IsZero())
	require.False(t, NewUint256(1).IsZero())
	
	u := NewUint256(0)
	u[1] = 1
	require.False(t, u.IsZero())
}

func TestMul64(t *testing.T) {
	u := NewUint256(10)
	u.Mul64(10)
	require.Equal(t, "100", u.String())

	// Test overflow from one word to next
	// Max Uint64
	u = new(Uint256).SetBytes(new(big.Int).SetUint64(^uint64(0)).Bytes())
	u.Mul64(2)
	// (2^64 - 1) * 2 = 2^65 - 2
	expected := new(big.Int).Mul(new(big.Int).SetUint64(^uint64(0)), big.NewInt(2))
	require.Equal(t, expected.String(), u.String())
}

func TestAdd64(t *testing.T) {
	u := NewUint256(10)
	u.Add64(5)
	require.Equal(t, "15", u.String())

	// Test carry
	u = new(Uint256).SetBytes(new(big.Int).SetUint64(^uint64(0)).Bytes())
	u.Add64(1)
	// 2^64
	expected := new(big.Int).Add(new(big.Int).SetUint64(^uint64(0)), big.NewInt(1))
	require.Equal(t, expected.String(), u.String())
}
