package app

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	"strings"
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

		sumU := new(Uint256).Add(*u1, *u2)

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
		diffU := new(Uint256).Sub(*u1, *u2)

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

			require.Equal(t, tt.expected, uA.Cmp(*uB))
		})
	}
}

func TestString(t *testing.T) {
	tests := []string{
		"0",
		"1",
		"123456789",
		"18446744073709551615", // Max Uint64
		"340282366920938463463374607431768211455",                                        // Max Uint128
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

func TestToHex(t *testing.T) {
	tests := []struct {
		val      string // decimal
		expected string // hex with 0x prefix
	}{
		{"0", "0x0"},
		{"1", "0x1"},
		{"10", "0xa"},
		{"15", "0xf"},
		{"16", "0x10"},
		{"255", "0xff"},
		{"256", "0x100"},
		{"100", "0x64"},
		{"12345", "0x3039"},
		{"18446744073709551615", "0xffffffffffffffff"}, // Max uint64
		{"340282366920938463463374607431768211455", "0xffffffffffffffffffffffffffffffff"}, // Max uint128
		{"115792089237316195423570985008687907853269984665640564039457584007913129639935", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"}, // Max uint256
	}

	for _, tt := range tests {
		t.Run(tt.val, func(t *testing.T) {
			b, _ := new(big.Int).SetString(tt.val, 10)
			u := new(Uint256).SetBytes(b.Bytes())
			require.Equal(t, tt.expected, u.ToHex())
		})
	}
}

func TestJSON(t *testing.T) {
	type wrapper struct {
		Val *Uint256 `json:"val"`
	}

	tests := []struct {
		valStr string
		hexStr string
	}{
		{"0", "0x0"},
		{"100", "0x64"},
		{"12345678901234567890", "0xab54a98ceb1f0ad2"},
	}

	for _, tt := range tests {
		t.Run(tt.valStr, func(t *testing.T) {
			// Test Marshal
			b, _ := new(big.Int).SetString(tt.valStr, 10)
			u := new(Uint256).SetBytes(b.Bytes())
			w := wrapper{Val: u}

			data, err := json.Marshal(w)
			require.NoError(t, err)

			// Expect hex string with 0x prefix
			// "val": "0x..."
			expected := fmt.Sprintf(`"val":"%s"`, tt.hexStr)
			// Remove spaces from data just in case
			jsonStr := strings.ReplaceAll(string(data), " ", "")
			require.Contains(t, jsonStr, expected)

			// Test Unmarshal
			var w2 wrapper
			err = json.Unmarshal(data, &w2)
			require.NoError(t, err)
			require.Equal(t, tt.valStr, w2.Val.String())
		})
	}

	// Test unmarshal from hex string (quoted)
	t.Run("quoted hex string", func(t *testing.T) {
		jsonStr := `{"val": "0x3039"}` // 12345
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.NoError(t, err)
		require.Equal(t, "12345", w.Val.String())
	})
	
	t.Run("invalid prefix", func(t *testing.T) {
		jsonStr := `{"val": "12345"}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid Uint256 prefix")
	})

	t.Run("unquoted decimal number rejected", func(t *testing.T) {
		jsonStr := `{"val": 12345}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid Uint256 format")
	})

	t.Run("leading zeros", func(t *testing.T) {
		// Leading zeros should be ignored (0x00000064 == 0x64)
		jsonStr := `{"val": "0x00000064"}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.NoError(t, err)
		require.Equal(t, "100", w.Val.String())

		// Compare with non-padded version
		jsonStr2 := `{"val": "0x64"}`
		var w2 wrapper
		err = json.Unmarshal([]byte(jsonStr2), &w2)
		require.NoError(t, err)
		require.Equal(t, w.Val.String(), w2.Val.String())
	})

	t.Run("uppercase hex", func(t *testing.T) {
		// Uppercase hex digits should work
		jsonStr := `{"val": "0xFF"}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.NoError(t, err)
		require.Equal(t, "255", w.Val.String())

		// Compare with lowercase version
		jsonStr2 := `{"val": "0xff"}`
		var w2 wrapper
		err = json.Unmarshal([]byte(jsonStr2), &w2)
		require.NoError(t, err)
		require.Equal(t, w.Val.String(), w2.Val.String())
	})

	t.Run("mixed case hex", func(t *testing.T) {
		jsonStr := `{"val": "0xAbCdEf"}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.NoError(t, err)
		require.Equal(t, "11259375", w.Val.String()) // 0xABCDEF = 11259375
	})

	t.Run("very long hex string near max", func(t *testing.T) {
		// Test with 64 hex digits (max uint256)
		maxHex := "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		jsonStr := `{"val": "` + maxHex + `"}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.NoError(t, err)
		max256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
		require.Equal(t, max256.String(), w.Val.String())
	})

	t.Run("hex string one digit less than max", func(t *testing.T) {
		// 63 hex digits (just under max)
		nearMaxHex := "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		jsonStr := `{"val": "` + nearMaxHex + `"}`
		var w wrapper
		err := json.Unmarshal([]byte(jsonStr), &w)
		require.NoError(t, err)
		expected, _ := new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
		require.Equal(t, expected.String(), w.Val.String())
	})
}

func TestUnmarshalJSONOverflow(t *testing.T) {
	// 2^256
	overMaxHex := "0x10000000000000000000000000000000000000000000000000000000000000000"
	jsonStr := `{"val": "` + overMaxHex + `"}`

	type wrapper struct {
		Val *Uint256 `json:"val"`
	}
	var w wrapper
	err := json.Unmarshal([]byte(jsonStr), &w)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Uint256 overflow")
}

func TestUnmarshalJSONRobustness(t *testing.T) {
	type wrapper struct {
		Val *Uint256 `json:"val"`
	}

	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "internal whitespace",
			input:       `{"val": "1 2 3"}`,
			expectError: true, // Should not allow internal spaces
		},
		{
			name:        "empty string",
			input:       `{"val": ""}`,
			expectError: true, // Should not allow empty value
		},
		{
			name:        "only whitespace",
			input:       `{"val": "   "}`,
			expectError: true, // Should not allow whitespace-only
		},
		{
			name:        "malformed JSON string",
			input:       `{"val": "\"123"}`, // missing closing quote
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrapper
			err := json.Unmarshal([]byte(tt.input), &w)
			if tt.expectError {
				require.Error(t, err, "Expected error for input: %s", tt.input)
			} else {
				require.NoError(t, err, "Expected no error for input: %s", tt.input)
			}
		})
	}
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

func TestAdd64Overflow(t *testing.T) {
	max256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	tests := []struct {
		name          string
		initial       string // decimal
		add           uint64
		expectOverflow bool
	}{
		{
			name:          "no overflow small",
			initial:       "100",
			add:           50,
			expectOverflow: false,
		},
		{
			name:          "no overflow max",
			initial:       max256.String(),
			add:           0,
			expectOverflow: false,
		},
		{
			name:          "overflow max plus one",
			initial:       max256.String(),
			add:           1,
			expectOverflow: true,
		},
		{
			name:          "overflow max plus max uint64",
			initial:       max256.String(),
			add:           ^uint64(0),
			expectOverflow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bInitial, _ := new(big.Int).SetString(tt.initial, 10)
			u := new(Uint256).SetBytes(bInitial.Bytes())

			sumBig := new(big.Int).Add(bInitial, new(big.Int).SetUint64(tt.add))
			expectedOverflow := sumBig.BitLen() > 256

			uCopy := *u
			overflow := (&uCopy).Add64Overflow(tt.add)

			require.Equal(t, expectedOverflow, overflow, "Overflow mismatch for %s + %d", tt.initial, tt.add)

			// Verify result matches Add64 (mod 2^256)
			uCopy2 := *u
			(&uCopy2).Add64(tt.add)
			require.Equal(t, uCopy2.String(), uCopy.String(), "Result mismatch for %s + %d", tt.initial, tt.add)
		})
	}

	// Random test for overflow detection
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		b := new(big.Int).Rand(r, new(big.Int).Lsh(big.NewInt(1), 256))
		add := r.Uint64()

		u := new(Uint256).SetBytes(b.Bytes())

		sumBig := new(big.Int).Add(b, new(big.Int).SetUint64(add))
		expectedOverflow := sumBig.BitLen() > 256

		uCopy := *u
		overflow := (&uCopy).Add64Overflow(add)

		require.Equal(t, expectedOverflow, overflow, "Random overflow test: %v + %d", b, add)
	}
}

// This test targets worst-case carry propagation during mul-by-word.
// If carry handling between limbs is wrong, this case should fail.
func TestMul64CarryPropagationWorstCase(t *testing.T) {
	// max256 = 2^256 - 1
	max256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	// max value for a uint64: 2^64 - 1
	y := ^uint64(0)

	expected := new(big.Int).Mul(new(big.Int).Set(max256), new(big.Int).SetUint64(y))
	expectedOverflow := expected.BitLen() > 256

	// expectedMod = expected mod 2^256
	mask := new(big.Int).Set(max256)
	expectedMod := new(big.Int).And(expected, mask)

	u := new(Uint256).SetBytes(max256.Bytes())

	// Mul64: modulo 2^256
	u1 := *u
	(&u1).Mul64(y)
	require.Equal(t, expectedMod.String(), u1.String())

	// Mul64Overflow: should produce same truncated value and report overflow.
	u2 := *u
	overflow := (&u2).Mul64Overflow(y)
	require.Equal(t, expectedOverflow, overflow)
	require.Equal(t, expectedMod.String(), u2.String())
}

func TestMul64Consistency(t *testing.T) {
	// This test proves that for Mul64(z, y), hi + c actually won't overflow
	// because of the 64-bit product limit
	max := ^uint64(0)
	hi, _ := bits.Mul64(max, max)

	if hi == max {
		t.Error("If hi could reach max, Mul64 would be broken. ")
	} else {
		t.Logf("Max hi in Mul64 is %x, which is < %x", hi, max)
	}
}

func TestAddOverflow(t *testing.T) {
	// Setup constants once
	// -
	// limit represents the exclusive upper bound of a 256-bit unsigned integer (2^256)
	// Any value >= limit cannot be represented in 256 bits and constitutes an overflow.
	limit := new(big.Int).Lsh(big.NewInt(1), 256)

	// (limit - 1) therefore is the maximum value
	max256Big := new(big.Int).Sub(limit, big.NewInt(1))

	max256 := *new(Uint256).SetBytes(max256Big.Bytes())
	one := *NewUint256(1)
	zero := *NewUint256(0)

	t.Run("no overflow", func(t *testing.T) {
		u1 := *NewUint256(100)
		u2 := *NewUint256(200)
		var res Uint256
		overflow := res.AddOverflow(u1, u2)
		require.False(t, overflow)
		require.Equal(t, "300", res.String())
	})

	t.Run("overflow boundary", func(t *testing.T) {
		var res Uint256
		overflow := res.AddOverflow(max256, one)
		require.True(t, overflow)
		require.Equal(t, "0", res.String())
	})

	t.Run("max + zero", func(t *testing.T) {
		var res Uint256
		overflow := res.AddOverflow(max256, zero)
		require.False(t, overflow)
		require.Equal(t, max256Big.String(), res.String())
	})

	t.Run("random cases", func(t *testing.T) {
		// Use a fixed seed for CI stability
		r := rand.New(rand.NewSource(1234))

		for i := range 100 {
			b1 := new(big.Int).Rand(r, limit)
			b2 := new(big.Int).Rand(r, limit)

			u1 := new(Uint256).SetBytes(b1.Bytes())
			u2 := new(Uint256).SetBytes(b2.Bytes())

			sumBig := new(big.Int).Add(b1, b2)

			// An overflow occurs if the result is >= 2^256
			expectedOverflow := sumBig.Cmp(limit) >= 0

			// Apply the full mask to simulate 256-bit wrapping, in other words all bits beyond 256 becomes 0
			wrappedBig := new(big.Int).And(sumBig, max256Big)

			var sumU Uint256
			overflow := sumU.AddOverflow(*u1, *u2)

			require.Equal(t, expectedOverflow, overflow, "Seed: 1337, Iteration: %d", i)
			require.Equal(t, wrappedBig.String(), sumU.String(), "Sum mismatch at index %d", i)
		}
	})
}
