package wasm

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

// defineHostCryptoImports wires the host-side SHA-256 and HMAC-SHA-512
// primitives into the WASM linker as the `env::host_sha256` and
// `env::host_hmac_sha512` imports. WASM guests built with TinyGo+WASI
// cannot use Go's stdlib crypto directly when invoked by the executor —
// `crypto/internal/fips140` depends on per-goroutine state populated by
// `_start`, which wasmtime-go's Linker.Instantiate does not invoke. Routing
// hashes through the host bypasses that lifecycle issue, keeps the guest
// binary small, and runs the hash natively on the host (typically faster).
//
// vela-common-go/wasm/hostcrypto is the guest-side facade; this is the
// host-side implementation those imports resolve to. Regression artifact:
// vela-ned/spikes/007-wasmtime-go-init/.
//
// ABI (matching the guest's //go:wasmimport declarations exactly):
//
//	env::host_sha256(in_ptr i32, in_len i32, out_ptr i32) -> ()
//	  Reads in_len bytes at in_ptr, writes the 32-byte SHA-256 digest to out_ptr.
//
//	env::host_hmac_sha512(key_ptr i32, key_len i32, msg_ptr i32, msg_len i32, out_ptr i32) -> ()
//	  Reads key_len bytes at key_ptr, msg_len bytes at msg_ptr, writes the
//	  64-byte HMAC-SHA-512 digest to out_ptr.
func defineHostCryptoImports(linker *wasmtime.Linker, store *wasmtime.Store) error {
	if err := linker.DefineFunc(store, "env", "host_sha256",
		func(caller *wasmtime.Caller, inPtr int32, inLen int32, outPtr int32) {
			mem := caller.GetExport("memory").Memory().UnsafeData(caller)
			input := make([]byte, inLen)
			copy(input, mem[inPtr:inPtr+inLen])
			digest := sha256.Sum256(input)
			copy(mem[outPtr:outPtr+32], digest[:])
		}); err != nil {
		return fmt.Errorf("define env::host_sha256: %w", err)
	}

	if err := linker.DefineFunc(store, "env", "host_hmac_sha512",
		func(caller *wasmtime.Caller, keyPtr int32, keyLen int32, msgPtr int32, msgLen int32, outPtr int32) {
			mem := caller.GetExport("memory").Memory().UnsafeData(caller)
			key := make([]byte, keyLen)
			copy(key, mem[keyPtr:keyPtr+keyLen])
			msg := make([]byte, msgLen)
			copy(msg, mem[msgPtr:msgPtr+msgLen])
			mac := hmac.New(sha512.New, key)
			mac.Write(msg)
			digest := mac.Sum(nil)
			copy(mem[outPtr:outPtr+64], digest[:64])
		}); err != nil {
		return fmt.Errorf("define env::host_hmac_sha512: %w", err)
	}

	return nil
}
