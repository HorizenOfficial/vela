package fullstack

// This file intentionally contains no code.
//
// It used to carry `// #cgo LDFLAGS: -Wl,--allow-multiple-definition` for the whole
// package. That is no longer needed, but the reasoning is kept here — at the place
// the directive would go — so that if the collision returns nobody has to
// rediscover the mechanism from a bare linker error.
//
// # The collision
//
// wasmtime-go ships a pre-built libwasmtime.a that statically includes zstd C
// sources (wasmtime uses zstd for its compiled-module cache). go-ethereum
// transitively imports github.com/DataDog/zstd, which compiles the same upstream C
// library via cgo. Fullstack tests link both, because they run real WASM (wasmtime)
// against a simulated chain (go-ethereum).
//
// The linker fails with "multiple definition of ZSTD_isError" and similar only when
// both copies are actually loaded, which depends on lazy archive extraction:
//
//   - DataDog/zstd is a cgo package, so its objects are loaded unconditionally and
//     define their zstd symbols up front.
//   - libwasmtime.a keeps zstd in dedicated archive members (zstd_common.o,
//     zstd_compress.o, …). Members are extracted only to resolve a symbol that is
//     still undefined.
//
// So a member is pulled in — and duplicates every symbol DataDog already defined —
// exactly when wasmtime's own code references a zstd symbol that DataDog's pinned
// zstd version does not provide.
//
// # Why it no longer happens
//
// Measured on linux-x86_64 against DataDog/zstd v1.4.5, which defines 301 global
// zstd-family symbols:
//
// +--------------------------------------------+--------------------+---------+
// |                                            | wasmtime-go v1.0.0 | v47.0.0 |
// +--------------------------------------------+--------------------+---------+
// |zstd symbols referenced by non-zstd members |           72       |    16   |
// |  of those, not provided by DataDog/zstd    |            7       |     0   |
// +--------------------------------------------+--------------------+---------+
//
// Under v1.0.0 the seven unsatisfied symbols (ZSTD_customMalloc, ZSTD_customCalloc,
// ZSTD_customFree, ZSTD_CCtx_trace, ZSTD_cycleLog, ZSTD_getDictID_fromCDict,
// ZSTD_ldm_fillHashTable — the ZSTD_custom* family arrived in zstd 1.4.6, just after
// DataDog's pin) forced extraction and hence the duplicates. Under v42.0.0 nothing
// was unsatisfied, and v47.0.0 measures identically, so no zstd member is ever
// extracted.
//
// Note what did NOT change: v47 still bundles zstd and still exports a large global
// zstd surface (313 global definitions, 250 of them shared with DataDog). This was
// never fixed upstream — the duplicates are still available, they are simply never
// both loaded.
//
// # When it would come back
//
// The condition is narrow: wasmtime's zstd requirements growing beyond what the
// pinned DataDog/zstd provides. That means a wasmtime-go upgrade, or go-ethereum
// moving off DataDog/zstd v1.4.5. The failure is loud (a link-time "multiple
// definition" error), never silent.
//
// # If it comes back
//
// Either align the zstd versions, or restore the workaround by adding to this file:
//
//	// #cgo LDFLAGS: -Wl,--allow-multiple-definition
//	import "C"
//
// with the comment immediately above the import (cgo requires the directive to be in
// that preamble). The flag keeps the first definition and discards later duplicates.
// That is safe here because both copies are the same upstream library used
// independently — wasmtime and DataDog/zstd never share zstd state — and the flag is
// scoped to binaries importing this package, so production binaries (cmd/manager,
// cmd/executor) are unaffected. It is a blunt instrument, though: it would also mask
// an unrelated genuine duplicate-symbol bug, which is why it is not kept armed.
