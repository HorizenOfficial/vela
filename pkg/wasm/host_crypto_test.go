package wasm

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/stretchr/testify/require"
)

// Layer-1 wire-up tests for the host-crypto bridge (see host_crypto.go
// + ADR-011 in vela-ned). These do not exercise the hash functions
// behaviorally — that's covered by vela-ned's integration test (Phase D2)
// and by the spike at vela-ned/spikes/007-wasmtime-go-init/. The job here
// is the cheap, ABI-stable regression check: "did we actually register
// both functions in the linker under the expected env::host_* names?"
//
// A regression where defineHostCryptoImports silently returns nil but
// fails to register one of the two imports would let any consuming app's
// build succeed and instantiation succeed (the unused import would just
// resolve via wasmtime's default lookup machinery), but the first runtime
// call to the missing primitive would trap deep inside the WASM. These
// tests catch that class of bug at vela's PR boundary, before it reaches
// any downstream consumer.

func TestDefineHostCryptoImports_RegistersBoth(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	store := wasmtime.NewStore(runtime.engine)
	linker := wasmtime.NewLinker(runtime.engine)

	require.NoError(t, defineHostCryptoImports(linker, store), "defineHostCryptoImports should succeed on a fresh linker")

	require.NotNil(t, linker.Get(store, "env", "host_sha256"),
		"env::host_sha256 must be defined in the linker after defineHostCryptoImports")
	require.NotNil(t, linker.Get(store, "env", "host_hmac_sha512"),
		"env::host_hmac_sha512 must be defined in the linker after defineHostCryptoImports")
}

func TestDefineHostCryptoImports_RejectsDuplicate(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	store := wasmtime.NewStore(runtime.engine)
	linker := wasmtime.NewLinker(runtime.engine)

	require.NoError(t, defineHostCryptoImports(linker, store))
	// Second call must fail because the names are already taken — proving
	// the first call actually populated the linker (and not, say, silently
	// no-op'd). wasmtime.Linker rejects redefining env::host_sha256.
	err := defineHostCryptoImports(linker, store)
	require.Error(t, err, "second defineHostCryptoImports call should fail due to duplicate import name")
}

// hostCryptoCallerWat is a guest that does nothing but forward its arguments to the
// host-crypto imports, so a test can drive them with values a real guest would never
// send. The pointers and lengths are guest-controlled in production too — they arrive
// as i32 from inside the sandbox — which is the whole reason the host must validate
// them rather than trust them.
const hostCryptoCallerWat = `(module
  (import "env" "host_sha256" (func $sha256 (param i32 i32 i32)))
  (import "env" "host_hmac_sha512" (func $hmac (param i32 i32 i32 i32 i32)))
  (memory (export "memory") 1)
  (func (export "call_sha256") (param $p i32) (param $l i32) (param $o i32)
    (call $sha256 (local.get $p) (local.get $l) (local.get $o)))
  (func (export "call_hmac") (param $kp i32) (param $kl i32) (param $mp i32) (param $ml i32) (param $o i32)
    (call $hmac (local.get $kp) (local.get $kl) (local.get $mp) (local.get $ml) (local.get $o)))
)`

// newHostCryptoGuest instantiates hostCryptoCallerWat against the real host-crypto
// imports and returns the store, the instance and its linear memory.
func newHostCryptoGuest(t *testing.T) (*wasmtime.Store, *wasmtime.Instance, *wasmtime.Memory) {
	t.Helper()

	runtime := NewWasmtimeRuntime(testLogger, 0)
	t.Cleanup(func() { runtime.Close() })

	store := wasmtime.NewStore(runtime.engine)
	linker := wasmtime.NewLinker(runtime.engine)
	require.NoError(t, defineHostCryptoImports(linker, store))

	wasmBytes, err := wasmtime.Wat2Wasm(hostCryptoCallerWat)
	require.NoError(t, err)
	module, err := wasmtime.NewModule(runtime.engine, wasmBytes)
	require.NoError(t, err)
	instance, err := linker.Instantiate(store, module)
	require.NoError(t, err)

	memory := instance.GetExport(store, "memory").Memory()
	require.NotNil(t, memory)
	return store, instance, memory
}

// TestHostCryptoComputesCorrectDigests is the baseline the exclusion tests below need:
// without it, "every call traps" would also pass them.
func TestHostCryptoComputesCorrectDigests(t *testing.T) {
	store, instance, memory := newHostCryptoGuest(t)

	const (
		inPtr  = 0
		outPtr = 1024
	)
	input := []byte("the quick brown fox")
	copy(memory.UnsafeData(store)[inPtr:], input)

	_, err := instance.GetFunc(store, "call_sha256").Call(store, int32(inPtr), int32(len(input)), int32(outPtr))
	require.NoError(t, err)

	want := sha256.Sum256(input)
	require.Equal(t, want[:], memory.UnsafeData(store)[outPtr:outPtr+32],
		"host_sha256 must write the SHA-256 of exactly the requested window")

	// HMAC over the same buffer, key and message deliberately overlapping regions.
	const (
		keyPtr     = 0
		keyLen     = 4
		hmacOutPtr = 2048
	)
	_, err = instance.GetFunc(store, "call_hmac").Call(store,
		int32(keyPtr), int32(keyLen), int32(inPtr), int32(len(input)), int32(hmacOutPtr))
	require.NoError(t, err)

	mac := hmac.New(sha512.New, input[:keyLen])
	mac.Write(input)
	require.Equal(t, mac.Sum(nil), memory.UnsafeData(store)[hmacOutPtr:hmacOutPtr+64],
		"host_hmac_sha512 must write the HMAC-SHA-512 of exactly the requested windows")
}

// TestHostCryptoRejectsHostileArguments is the exclusion test for the bridge: the
// pointers and lengths come from inside the sandbox, so every one of them can be
// negative or out of range.
//
// What must NOT happen is a Go panic. These run on the calling goroutine, and
// wasmtime-go propagates a panic from a host function straight back out through
// Func.Call — it is not contained as a trap. Inside the enclave that takes the whole
// executor down, from a guest, with no signed failure and no fee: strictly worse than
// the stall this bound exists to prevent. Trapping instead makes it an ordinary guest
// fault, which the runtime already classifies and evicts for.
func TestHostCryptoRejectsHostileArguments(t *testing.T) {
	// One page of linear memory.
	const memSize = 64 * 1024

	cases := []struct {
		name string
		fn   string
		args []interface{}
	}{
		{"sha256 negative length", "call_sha256", []interface{}{int32(0), int32(-1), int32(1024)}},
		{"sha256 negative input pointer", "call_sha256", []interface{}{int32(-1), int32(4), int32(1024)}},
		{"sha256 input past end of memory", "call_sha256", []interface{}{int32(memSize - 8), int32(64), int32(1024)}},
		{"sha256 input length overflows int32", "call_sha256", []interface{}{int32(1024), int32(2147483647), int32(1024)}},
		{"sha256 negative output pointer", "call_sha256", []interface{}{int32(0), int32(4), int32(-1)}},
		{"sha256 output digest past end of memory", "call_sha256", []interface{}{int32(0), int32(4), int32(memSize - 8)}},

		{"hmac negative key length", "call_hmac", []interface{}{int32(0), int32(-1), int32(0), int32(4), int32(2048)}},
		{"hmac negative message length", "call_hmac", []interface{}{int32(0), int32(4), int32(0), int32(-1), int32(2048)}},
		{"hmac key past end of memory", "call_hmac", []interface{}{int32(memSize - 8), int32(64), int32(0), int32(4), int32(2048)}},
		{"hmac message past end of memory", "call_hmac", []interface{}{int32(0), int32(4), int32(memSize - 8), int32(64), int32(2048)}},
		{"hmac output digest past end of memory", "call_hmac", []interface{}{int32(0), int32(4), int32(0), int32(4), int32(memSize - 8)}},
		{"hmac negative output pointer", "call_hmac", []interface{}{int32(0), int32(4), int32(0), int32(4), int32(-1)}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store, instance, _ := newHostCryptoGuest(t)
			fn := instance.GetFunc(store, tc.fn)
			require.NotNil(t, fn)

			var err error
			require.NotPanics(t, func() { _, err = fn.Call(store, tc.args...) },
				"a hostile argument panicked the HOST: inside the enclave this kills the "+
					"executor process, triggered by any deployed guest")
			require.Error(t, err,
				"a hostile argument must trap the guest, not be silently accepted")
		})
	}
}
