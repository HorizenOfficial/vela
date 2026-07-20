package fullstack

// Problem:
// wasmtime-go ships a pre-built libwasmtime.a that statically includes
// zstd C sources (used internally by wasmtime's compiled-module cache). Meanwhile,
// go-ethereum transitively imports github.com/DataDog/zstd (v1.4.5), which compiles
// the same zstd C library via cgo. When both end up in the same binary — as they do
// in fullstack tests that run real WASM (wasmtime) on a simulated chain (go-ethereum)
// — the GNU linker sees two definitions of every zstd symbol (ZSTD_isError,
// ZSTD_compress, ZSTD_decompress, …) and fails with "multiple definition of" errors.
//
// Fix:
// --allow-multiple-definition tells the linker to keep the first definition it
// encounters and silently discard subsequent duplicates, instead of treating them
// as errors.
//
// Why this is safe:
//   - Both symbol sets come from the same upstream C library (facebook/zstd). The
//     function signatures, semantics, and ABI are identical. Picking either copy
//     produces correct behaviour.
//   - wasmtime and DataDog/zstd use their own zstd copies independently — they do
//     not share zstd state across library boundaries. Even if the two copies differ
//     slightly in patch version, they operate in isolation.
//   - The flag is scoped to binaries that import this package. Production binaries
//     (cmd/manager, cmd/executor) and the existing mock-based system tests
//     (tests/system/) are completely unaffected.
//   - Upgrading wasmtime-go does NOT fix the root cause: every released version
//     ships libwasmtime.a with the default Cargo "cache" feature, which bundles
//     zstd. Removing it would require a custom wasmtime build without cache support.

// #cgo LDFLAGS: -Wl,--allow-multiple-definition
import "C"
