package common

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBig_ToInt verifies that ToInt returns the underlying big.Int value.
func TestBig_ToInt(t *testing.T) {
	bi := big.NewInt(12345)
	b := ToBig(bi)

	result := b.ToInt()
	assert.Equal(t, bi, result)
	assert.Equal(t, int64(12345), result.Int64())
}

// TestToBig verifies that ToBig correctly converts a big.Int to Big.
func TestToBig(t *testing.T) {
	bi := big.NewInt(98765)
	b := ToBig(bi)

	assert.NotNil(t, b)
	assert.Equal(t, int64(98765), b.ToInt().Int64())
}

// TestBig_String verifies that String returns the decimal representation.
func TestBig_String(t *testing.T) {
	bi := big.NewInt(123456789)
	b := ToBig(bi)

	assert.Equal(t, "123456789", b.String())
}

// TestBig_MarshalJSON verifies JSON marshaling produces correct hex strings with 0x prefix.
func TestBig_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     *big.Int
		expected  string
		expectErr bool
	}{
		{
			name:     "zero", // Zero value marshals to "0x0"
			input:    big.NewInt(0),
			expected: `"0x0"`,
		},
		{
			name:     "small positive", // Small values marshal correctly
			input:    big.NewInt(100),
			expected: `"0x64"`,
		},
		{
			name:     "medium positive", // Medium values marshal correctly
			input:    big.NewInt(12345),
			expected: `"0x3039"`,
		},
		{
			name:     "large positive", // Large values beyond int64 marshal correctly
			input:    func() *big.Int { bi, _ := new(big.Int).SetString("12345678901234567890", 10); return bi }(),
			expected: `"0xab54a98ceb1f0ad2"`,
		},
		{
			name:     "max uint256", // Maximum uint256 value marshals correctly
			input:    new(big.Int).SetBytes([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
			expected: `"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`,
		},
		{
			name:      "negative value", // Negative values are rejected
			input:     big.NewInt(-100),
			expectErr: true,
		},
		{
			name:     "nil pointer", // Nil pointer marshals to null
			input:    nil,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b *Big
			if tt.input != nil {
				b = ToBig(tt.input)
			}

			result, err := b.MarshalJSON()
			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// TestBig_UnmarshalJSON verifies JSON unmarshaling correctly parses hex strings and handles errors.
func TestBig_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  *big.Int
		expectErr bool
	}{
		{
			name:     "zero", // Zero value unmarshals correctly
			input:    `"0x0"`,
			expected: big.NewInt(0),
		},
		{
			name:     "small positive", // Small hex values unmarshal correctly
			input:    `"0x64"`,
			expected: big.NewInt(100),
		},
		{
			name:     "medium positive", // Medium hex values unmarshal correctly
			input:    `"0x3039"`,
			expected: big.NewInt(12345),
		},
		{
			name:     "large positive", // Large hex values beyond int64 unmarshal correctly
			input:    `"0xab54a98ceb1f0ad2"`,
			expected: func() *big.Int { bi, _ := new(big.Int).SetString("12345678901234567890", 10); return bi }(),
		},
		{
			name:     "null value", // Null JSON unmarshals to zero value
			input:    `null`,
			expected: big.NewInt(0),
		},
		{
			name:     "uppercase hex", // Uppercase hex digits are accepted
			input:    `"0xABCD"`,
			expected: big.NewInt(43981),
		},
		{
			name:      "missing quotes", // Unquoted strings are rejected
			input:     `0x64`,
			expectErr: true,
		},
		{
			name:      "missing 0x prefix", // Missing 0x prefix is rejected
			input:     `"64"`,
			expectErr: true,
		},
		{
			name:      "empty hex string", // Empty hex after 0x is rejected
			input:     `"0x"`,
			expectErr: true,
		},
		{
			name:      "invalid hex characters", // Non-hex characters are rejected
			input:     `"0xghij"`,
			expectErr: true,
		},
		{
			name:      "too short", // Input too short is rejected
			input:     `"0x"`,
			expectErr: true,
		},
		{
			name:      "unquoted string", // Unquoted numeric strings are rejected
			input:     `12345`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Big
			err := b.UnmarshalJSON([]byte(tt.input))
			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.expected == nil {
				assert.Zero(t, b.ToInt().Sign())
			} else {
				assert.Equal(t, 0, b.ToInt().Cmp(tt.expected))
			}
		})
	}
}

// TestBig_JSONRoundTrip verifies that marshal/unmarshal preserves values exactly.
func TestBig_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		value *big.Int
	}{
		{
			name:  "zero", // Zero value round-trips correctly
			value: big.NewInt(0),
		},
		{
			name:  "small", // Small values round-trip correctly
			value: big.NewInt(100),
		},
		{
			name:  "medium", // Medium values round-trip correctly
			value: big.NewInt(123456789),
		},
		{
			name:  "large", // Large values round-trip correctly
			value: func() *big.Int { bi, _ := new(big.Int).SetString("12345678901234567890", 10); return bi }(),
		},
		{
			name:  "very large", // Very large values round-trip correctly
			value: new(big.Int).SetBytes([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := ToBig(tt.value)

			// Marshal
			jsonBytes, err := original.MarshalJSON()
			require.NoError(t, err)

			// Unmarshal
			var result Big
			err = result.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare
			assert.Equal(t, 0, original.ToInt().Cmp(result.ToInt()))
		})
	}
}

// TestBig_JSONStructRoundTrip verifies Big works correctly within JSON structs with multiple fields.
func TestBig_JSONStructRoundTrip(t *testing.T) {
	type TestStruct struct {
		Amount    *Big `json:"amount"`
		Deposit   *Big `json:"deposit"`
		Fee       *Big `json:"fee"`
		NullField *Big `json:"nullField"`
	}

	original := TestStruct{
		Amount:    ToBig(big.NewInt(1000)),
		Deposit:   ToBig(big.NewInt(500)),
		Fee:       ToBig(big.NewInt(10)),
		NullField: nil,
	}

	// Marshal
	jsonBytes, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal
	var result TestStruct
	err = json.Unmarshal(jsonBytes, &result)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, 0, original.Amount.ToInt().Cmp(result.Amount.ToInt()))
	assert.Equal(t, 0, original.Deposit.ToInt().Cmp(result.Deposit.ToInt()))
	assert.Equal(t, 0, original.Fee.ToInt().Cmp(result.Fee.ToInt()))
	assert.Nil(t, result.NullField)
}

// TestBig_MarshalJSON_Negative verifies that negative values are rejected during marshaling.
func TestBig_MarshalJSON_Negative(t *testing.T) {
	negative := big.NewInt(-100)
	b := ToBig(negative)

	_, err := b.MarshalJSON()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "negative")
}

// TestBig_UnmarshalJSON_EdgeCases verifies error handling for malformed JSON inputs.
func TestBig_UnmarshalJSON_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expectErr bool
	}{
		{
			name:      "empty input", // Empty input is rejected
			input:     []byte{},
			expectErr: true,
		},
		{
			name:      "single quote", // Incomplete quotes are rejected
			input:     []byte(`"`),
			expectErr: true,
		},
		{
			name:      "unclosed quote", // Unclosed quotes are rejected
			input:     []byte(`"0x64`),
			expectErr: true,
		},
		{
			name:      "wrong prefix", // Uppercase 0X prefix is rejected
			input:     []byte(`"0X64"`),
			expectErr: true,
		},
		{
			name:      "wrong prefix 2", // Double 'x' in prefix is rejected
			input:     []byte(`"0xx64"`),
			expectErr: true,
		},
		{
			name:      "space in hex", // Spaces in hex string are rejected
			input:     []byte(`"0x 64"`),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Big
			err := b.UnmarshalJSON(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// BenchmarkBig_MarshalJSON benchmarks JSON marshaling performance.
func BenchmarkBig_MarshalJSON(b *testing.B) {
	bi, _ := new(big.Int).SetString("12345678901234567890", 10)
	val := ToBig(bi)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = val.MarshalJSON()
	}
}

// BenchmarkBig_UnmarshalJSON benchmarks JSON unmarshaling performance.
func BenchmarkBig_UnmarshalJSON(b *testing.B) {
	data := []byte(`"0xab54a98ceb1f0ad2"`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Big
		_ = result.UnmarshalJSON(data)
	}
}

// BenchmarkBig_JSONRoundTrip benchmarks complete marshal/unmarshal round-trip performance.
func BenchmarkBig_JSONRoundTrip(b *testing.B) {
	bi, _ := new(big.Int).SetString("12345678901234567890", 10)
	val := ToBig(bi)
	jsonBytes, _ := val.MarshalJSON()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Big
		_ = result.UnmarshalJSON(jsonBytes)
	}
}
