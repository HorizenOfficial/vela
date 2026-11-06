package common

import (
	"testing"

	"github.com/magiconair/properties"
	"github.com/stretchr/testify/require"
)

func TestGetConfigVarInt64(t *testing.T) {
	props := properties.NewProperties()
	require.NoError(t, props.SetValue("FROM_FILE", 42))
	require.NoError(t, props.SetValue("ABOVE_INT32", int64(3221225472))) // 3 GiB, > math.MaxInt32
	require.NoError(t, props.SetValue("NEGATIVE", -7))
	props.MustSet("GARBAGE", "not-a-number")

	t.Run("missing var returns default", func(t *testing.T) {
		require.Equal(t, int64(5), GetConfigVarInt64("MISSING", 5, props))
	})

	t.Run("file value is parsed", func(t *testing.T) {
		require.Equal(t, int64(42), GetConfigVarInt64("FROM_FILE", 5, props))
	})

	t.Run("env overrides file", func(t *testing.T) {
		t.Setenv("FROM_FILE", "43")
		require.Equal(t, int64(43), GetConfigVarInt64("FROM_FILE", 5, props))
	})

	t.Run("values above int32 range are parsed, not defaulted", func(t *testing.T) {
		// Regression: ParseInt used bitSize 32, silently turning any value
		// above 2^31-1 (e.g. EXECUTOR_MAX_GUEST_MEMORY_BYTES=3GiB) into the
		// default so validation never saw it.
		require.Equal(t, int64(3221225472), GetConfigVarInt64("ABOVE_INT32", 0, props))
	})

	t.Run("negative values are parsed", func(t *testing.T) {
		require.Equal(t, int64(-7), GetConfigVarInt64("NEGATIVE", 5, props))
	})

	t.Run("unparseable value returns default", func(t *testing.T) {
		require.Equal(t, int64(5), GetConfigVarInt64("GARBAGE", 5, props))
	})
}
