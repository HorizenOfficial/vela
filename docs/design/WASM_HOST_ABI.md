# WASM host ABI

How the Executor's Wasmtime host talks to a guest application module, and how that
contract is tested.

The host side lives in `pkg/wasm/wasmtime_runtime.go`; the shared result types live
in `pkg/wasm/common/types.go`. **The code is the source of truth** — this document
explains the contract and the reasoning behind it, and must be updated whenever the
exports, their signatures, or the result encoding change.

## Exports the guest must provide

The host resolves guest functions **by name only**; signatures are checked when the
call is made, not when the module is loaded (see "Failure classification" below for
why that matters). Parameters are in call order.

| Export | Signature | Called when |
|---|---|---|
| `memory` | exported linear memory | required; the host reads and writes it directly |
| `load_module` | `(appId i64) -> i32` | first use of an app, to warm the module cache |
| `deploy` | `(appId i64, paramsPtr i32, paramsLen i32) -> i32` | deploy, with constructor params |
| `deposit` | `(appId i64, senderPtr i32, senderLen i32, tokenPtr i32, tokenLen i32, valuePtr i32, valueLen i32, statePtr i32, stateLen i32) -> i32` | a deposit request |
| `process_request` | `(appId i64, senderPtr i32, senderLen i32, requestType i32, payloadPtr i32, payloadLen i32, statePtr i32, stateLen i32) -> i32` | all request types except `TrustProcess` |
| `trusted_request` | `(appId i64, payloadPtr i32, payloadLen i32, statePtr i32, stateLen i32) -> i32` | `TrustProcess` requests only |
| `allocate` | `(size i32) -> i32` | before every host→guest byte transfer |
| `deallocate` | `(ptr i32, len i32)` | **optional**; skipped entirely if absent |

`requestType` is `common.RequestType` widened to `i32`. Note that `trusted_request`
deliberately receives neither `sender` (the trusted path never reads it) nor
`requestType` (implicit for that export) — so the two request exports have different
arities, which is a real source of host/guest version skew.

Two further exports are used only by reference and debug paths
(`GetAllocatedMemoryStats` / `GetAllocatedMemoryStats2`), are absent from the
executor's `Runtime` interface, and are not required of a production app:
`get_allocated_memory_stats (outPtr i32)` and `get_memory_stats () -> i32`.

WASI is available: the host defines it on the linker (`linker.DefineWasi()`), which
is what TinyGo-generated modules need in order to instantiate.

## Passing bytes in: `allocate`

To hand the guest a byte slice the host calls `allocate(len)`, expects a non-zero
offset back, bounds-checks the whole `[ptr, ptr+len)` window against the current
memory size, and copies into it.

Ownership: the **host** frees what it allocated, via a deferred `deallocate(ptr, len)`
after the guest call returns. Empty input is special-cased — the host passes pointer
`0` with length `0` and does not allocate or free anything, so a guest that never
receives non-empty input is never asked to `allocate` at all.

A guest that returns `0`, a non-integer, or an out-of-range offset from `allocate`
fails the request. Returning `0` is how an out-of-memory guest surfaces: under the
linear-memory cap, a failed `memory.grow` makes TinyGo panic rather than return.

## Getting bytes out: the result protocol

Every result-returning export returns an `i32` offset into its own linear memory. At
that offset the host reads:

```
offset      +0 +1 +2 +3   +4 ...
           [ length (uint32, little-endian) ][ JSON body, `length` bytes ]
```

The JSON body unmarshals into one of the result structs in `pkg/wasm/common`:
`LoadModuleResult`, `DeployResult`, `DepositResult`, `ProcessResult`. All carry
`state` and `fuel`; `ProcessResult` additionally carries `events`, `appEvents`,
`withdrawals` and an optional `report`. All carry an optional `error` string.

Rules the host enforces while reading:

- The offset must be non-zero and the whole `[offset, offset+4+length)` window must
  lie inside memory.
- A `length` of `0` is an error ("empty result from wasm module").
- The host frees the result buffer with `deallocate(offset, 4+length)` once copied.
- The body `{"error":"wasm serialization error"}` (`common.WasmSerializationError`)
  is a reserved sentinel a guest may return when it cannot serialise its own result.

**`error` is the guest's normal failure channel.** A non-empty `error` means the
application rejected the request: the host reports a signed on-chain failure, leaves
state unchanged, and keeps the compiled module cached. It is not a fault.

`fuel` is currently **self-reported by the guest** and trusted for fee computation.
Replacing it with host-metered fuel is a separate, tracked change.

## Constraints on the guest

**Linear memory is capped, and offsets must stay below 2 GiB.** The cap defaults to
2 GiB and is configurable with `EXECUTOR_MAX_GUEST_MEMORY_BYTES`, which may lower it
but not raise it. 2 GiB is a hard ABI ceiling, not a tuning choice: the host
exchanges guest pointers as **signed** 32-bit offsets, so no valid offset may reach
2^31. Worst-case enclave RAM is roughly `EXECUTOR_MAX_CACHED_MODULES` ×
`EXECUTOR_MAX_GUEST_MEMORY_BYTES`; size the two together. A guest exceeding the cap
sees `memory.grow` fail.

**The accepted WASM feature set is pinned explicitly**, not inherited from whatever
Wasmtime enables by default, because the Executor re-executes guest code and commits
the result on-chain. A module using a disabled proposal fails to compile.

| Enabled | Disabled |
|---|---|
| bulk memory, multi-value, reference types, fixed-width SIMD | relaxed SIMD (host-architecture-dependent results), memory64 (breaks the int32 offset ABI), multi-memory and threads (both break the per-memory RAM accounting), tail calls, function references, GC, wide arithmetic |

Cranelift NaN canonicalization is enabled, since NaN bit patterns are otherwise
host-dependent. Adding a proposal is a deliberate change in `newPinnedEngine`, and
newly-introduced proposals should be pinned there when Wasmtime is upgraded.

**Guest execution is currently unbounded** — no fuel, no epoch deadline, no stack
limit. A guest that does not return blocks the Executor. See the TODO in
`newPinnedEngine`; tracked separately.

## Failure classification and the module cache

The host sorts failures into three classes, because the response differs. The test
is **whether guest code actually ran**, not whether the failure looks deterministic:

| Class | Examples | Module cache |
|---|---|---|
| Application error | non-empty `error` in a well-formed result | kept — a healthy guest declining a request |
| Guest fault | a trap; `allocate` failing, returning `0`, or returning an out-of-range offset; a result the host cannot parse | **evicted** — the heap is in an unknown state and must not serve the next request |
| Host-side / static defect | a missing export; an export whose signature does not match the host ABI; a nil store | kept — nothing executed, so re-instantiating provably cannot change the outcome |

Two consequences worth knowing when writing a guest:

- A trap unwinds without running guest cleanup, so an evicted module is recompiled
  and re-instantiated on the next request. That is intentional and cheap relative to
  reusing a damaged heap.
- **A signature mismatch is deployable.** Only `deploy`'s signature is exercised at
  deploy time, so an app built against a stale host ABI deploys cleanly and then
  fails every subsequent request when the host cannot invoke the export. This is why
  ABI changes need this document kept current.

## Guest logging

The guest's stdout and stderr are wired to a FIFO and forwarded to the log server,
so `println` from a guest — including TinyGo's own `panic: ...` output — reaches the
Executor logs. This is often the only place the *reason* for a trap appears: the trap
itself surfaces to the host as an opaque `unreachable`.

## Notes for TinyGo guests

Built with `//export <name>` and `-target=wasip1` (see `app/simple/Makefile`).

`recover()` cannot be used to convert a guest failure into an `error` result:
explicit panics, index-out-of-range, nil dereference and out-of-memory all bypass
deferred recovery and reach the host as an `unreachable` trap. A guest must therefore
avoid panicking — check inputs and return the `error` field — rather than trying to
catch panics.

## How this contract is tested

Two layers, deliberately:

**WAT fixtures (`pkg/wasm/wasmtime_runtime_test.go`)** assemble tiny modules from
WebAssembly text at test time via `wasmtime.Wat2Wasm`, so they need no TinyGo
toolchain and run in milliseconds. They implement just enough of the ABI to reach the
path under test: an `allocate` that returns a fixed scratch offset, exports that
ignore their parameters, and results pre-seeded as `(data ...)` blocks with a
hand-computed length prefix.

Their purpose is to drive the **host** through behaviour a correct guest never
exhibits, much of which cannot be expressed in a high-level language at all: a trap,
an unterminating call, a valid length prefix followed by non-JSON, a missing
`allocate`, an export with the wrong arity, an `allocate` returning a pointer past
the end of memory, or a module declaring two memories. Each fixture's own doc comment
states what it proves.

Because the length prefix and the offsets are hand-written, compute them rather than
estimating — a prefix that disagrees with its body makes the host read truncated or
trailing garbage, and the test then fails for a reason unrelated to its subject.

**TinyGo integration tests (`app/simple/integration_test.go`)** compile the real
guest and exercise the real allocator, real serialisation and real business logic.

Neither layer is sufficient alone. A WAT fixture encodes what we *believe* the ABI
is, so it can pass while a real guest fails; the integration tests are the backstop
that catches ABI drift, and they are also what proves the pinned feature set does not
reject legitimate modules.
