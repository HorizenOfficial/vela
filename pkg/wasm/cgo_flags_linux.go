// Force the GNU linker to mark the final binary's PT_GNU_STACK as
// non-executable even though wasmtime-go's bundled libwasmtime.a ships an
// assembly object (x86_64.o — the fiber-switching trampolines) that omits
// the `.note.GNU-stack` section. Without this flag GNU ld emits
//
//	/usr/bin/ld: warning: x86_64.o: missing .note.GNU-stack section
//	implies executable stack
//	/usr/bin/ld: NOTE: This behaviour is deprecated and will be removed in
//	a future version of the linker
//
// on every link of a binary that pulls in pkg/wasm. The flag is link-time
// only — it sets the PT_GNU_STACK program-header bits in the ELF output
// and has no runtime effect.
//
// Linux-only build tag because the warning and the `-z noexecstack` option
// are specific to GNU ld. macOS (ld64) and Windows (link.exe) neither emit
// the warning nor recognize the option.
//
// REMOVE-WHEN: wasmtime-go is upgraded to >= v42 (audited 2026-06-30).
// In v44.0.0's bundled libwasmtime.a, every one of the 360 archive members
// carries `.note.GNU-stack` — the offending v1.0.0 `x86_64.o` is gone, so
// the warning no longer fires and this file becomes dead weight.

//go:build linux

package wasm

// #cgo LDFLAGS: -Wl,-z,noexecstack
import "C"
