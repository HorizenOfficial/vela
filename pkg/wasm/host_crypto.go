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
//
// Every pointer and length is supplied by the guest and is validated before use; an
// out-of-range or negative value traps the guest rather than being indexed on. The
// wasm-level signatures are unchanged — the *wasmtime.Trap return is host-side only,
// so the guest's //go:wasmimport declarations still match.
func defineHostCryptoImports(linker *wasmtime.Linker, store *wasmtime.Store) error {
	if err := linker.DefineFunc(store, "env", "host_sha256",
		func(caller *wasmtime.Caller, inPtr int32, inLen int32, outPtr int32) *wasmtime.Trap {
			mem, trap := guestMemory(caller)
			if trap != nil {
				return trap
			}
			input, trap := guestWindow(mem, inPtr, inLen, "host_sha256 input")
			if trap != nil {
				return trap
			}
			out, trap := guestWindow(mem, outPtr, sha256.Size, "host_sha256 output")
			if trap != nil {
				return trap
			}

			digest := sha256.Sum256(input)
			copy(out, digest[:])
			return nil
		}); err != nil {
		return fmt.Errorf("define env::host_sha256: %w", err)
	}

	if err := linker.DefineFunc(store, "env", "host_hmac_sha512",
		func(caller *wasmtime.Caller, keyPtr int32, keyLen int32, msgPtr int32, msgLen int32, outPtr int32) *wasmtime.Trap {
			mem, trap := guestMemory(caller)
			if trap != nil {
				return trap
			}
			key, trap := guestWindow(mem, keyPtr, keyLen, "host_hmac_sha512 key")
			if trap != nil {
				return trap
			}
			msg, trap := guestWindow(mem, msgPtr, msgLen, "host_hmac_sha512 message")
			if trap != nil {
				return trap
			}
			out, trap := guestWindow(mem, outPtr, sha512.Size, "host_hmac_sha512 output")
			if trap != nil {
				return trap
			}

			mac := hmac.New(sha512.New, key)
			mac.Write(msg)
			copy(out, mac.Sum(nil))
			return nil
		}); err != nil {
		return fmt.Errorf("define env::host_hmac_sha512: %w", err)
	}

	return nil
}

// guestMemory returns the caller's exported linear memory, or a trap if the guest has
// none. A module reaching these imports without an exported `memory` is malformed, but
// the host must not dereference its way to finding that out.
func guestMemory(caller *wasmtime.Caller) ([]byte, *wasmtime.Trap) {
	export := caller.GetExport("memory")
	if export == nil {
		return nil, wasmtime.NewTrap("host crypto: guest has no exported memory")
	}
	mem := export.Memory()
	if mem == nil {
		return nil, wasmtime.NewTrap("host crypto: guest export `memory` is not a memory")
	}
	return mem.UnsafeData(caller), nil
}

// guestWindow bounds-checks [ptr, ptr+length) against the guest's linear memory and
// returns that window, or a trap describing what was wrong.
//
// Every argument here crosses the sandbox boundary: the guest chooses the pointers and
// the lengths, and nothing upstream constrains them. Indexing on them directly is how
// the host ends up panicking on a negative length or an out-of-range offset — and a Go
// panic in a host function is NOT converted into a trap, it propagates out through
// Func.Call (wasmtime-go documents this: "if the function `f` panics then the panic
// will be propagated to the caller"). Inside the enclave nothing recovers it, so the
// executor process dies at the request of whichever guest called in. Trapping instead
// keeps it an ordinary guest fault, which the runtime already classifies and evicts for.
//
// The arithmetic is done in int64 so that ptr+length cannot itself overflow and wrap
// back into range, which is exactly how a crafted pair of int32s would slip past a
// naive check.
func guestWindow(mem []byte, ptr int32, length int32, what string) ([]byte, *wasmtime.Trap) {
	if length < 0 {
		return nil, wasmtime.NewTrap(fmt.Sprintf("host crypto: %s length %d is negative", what, length))
	}
	if ptr < 0 {
		return nil, wasmtime.NewTrap(fmt.Sprintf("host crypto: %s pointer %d is negative", what, ptr))
	}

	start := int64(ptr)
	end := start + int64(length)
	if end > int64(len(mem)) {
		return nil, wasmtime.NewTrap(fmt.Sprintf(
			"host crypto: %s window [%d, %d) exceeds guest memory size %d", what, start, end, len(mem)))
	}

	return mem[start:end], nil
}
