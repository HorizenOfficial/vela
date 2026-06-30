package wasm

import (
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
