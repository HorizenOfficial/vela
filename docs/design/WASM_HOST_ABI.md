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
| `deploy` | `(appId i64, paramsPtr i32, paramsLen i32) -> i32` | deploy, with constructor params |
| `deposit` | `(appId i64, senderPtr i32, senderLen i32, tokenPtr i32, tokenLen i32, valuePtr i32, valueLen i32, statePtr i32, stateLen i32) -> i32` | a deposit request |
| `process_request` | `(appId i64, senderPtr i32, senderLen i32, requestType i32, payloadPtr i32, payloadLen i32, statePtr i32, stateLen i32) -> i32` | all request types except `TrustProcess` |
| `trusted_request` | `(appId i64, payloadPtr i32, payloadLen i32, statePtr i32, stateLen i32) -> i32` | `TrustProcess` requests only |
| `allocate` | `(size i32) -> i32` | before every host→guest byte transfer |
| `deallocate` | `(ptr i32, len i32)` | **optional**; skipped entirely if absent |

Cache warm-up requires no export beyond this table: on first use of an app (or
after an executor restart) the host only compiles and instantiates the module,
invoking no guest function — the application state is supplied by the host from
persistent storage, not produced by the guest.

`requestType` is `common.RequestType` widened to `i32`. Note that `trusted_request`
deliberately receives neither `sender` (the trusted path never reads it) nor
`requestType` (implicit for that export) — so the two request exports have different
arities, which is a real source of host/guest version skew.

Two further exports are used only by reference and debug paths
(`GetAllocatedMemoryStats` / `GetAllocatedMemoryStats2`), are absent from the
executor's `Runtime` interface, and are not required of a production app:
`get_allocated_memory_stats (outPtr i32)` and `get_memory_stats () -> i32`.

## Host imports: a closed set

WASI is defined on the linker (`linker.DefineWasi()`), which is what TinyGo-generated
modules need in order to instantiate — but a guest may only **declare** the following
imports, and nothing else:

```
wasi_snapshot_preview1.args_get           wasi_snapshot_preview1.environ_sizes_get
wasi_snapshot_preview1.args_sizes_get     wasi_snapshot_preview1.fd_write
wasi_snapshot_preview1.clock_time_get     wasi_snapshot_preview1.proc_exit
wasi_snapshot_preview1.environ_get        wasi_snapshot_preview1.random_get
```

This is the import set TinyGo's `wasip1` target emits for a guest that does no I/O,
so a normal app satisfies it without doing anything. The module's import section is
checked after compilation and before instantiation, so a module declaring anything
else is refused **before any of it runs**, including its start section. The refusal
is a signed `FAILED_LOADING_OR_GETTING_MODULE` failure: at deploy the app simply
does not deploy.

Refusal is triggered by *declaring* an import, not by calling it. A guest that links
`time.Sleep` anywhere reachable will import `poll_oneoff` and be rejected even if it
never sleeps — rebuild without that path.

**The rule this enforces, for anyone adding a host function:** a host import must
return without blocking and without unbounded work. Guest execution is bounded by
epoch interruption, which works through checks the compiler emits into *guest* code;
while a host call is running, the guest bound, the caller's execution budget and
shutdown are all suspended, and the runtime's global execution lock is held against
every other application. A blocking import is therefore a system-wide stall that
nothing can interrupt, and — since the request never completes — one that is never
signed and never charged for. Adding to the allowed set is an ABI change: the new
function must satisfy that rule, and must also satisfy the `externref` restriction
below. Removing from the set bricks already-deployed guests that declare the entry.

`clock_time_get` and `random_get` are permitted because they return immediately and
every TinyGo guest needs them. Note they make a state transition non-reproducible,
which matters if re-execution or multi-executor attestation is ever introduced.

## Passing bytes in: `allocate`

To hand the guest a byte slice the host calls `allocate(len)`, expects a non-zero
offset back, bounds-checks the whole `[ptr, ptr+len)` window against the current
memory size, and copies into it.

Ownership: the **host** frees what it allocated, via a deferred `deallocate(ptr, len)`
after the guest call returns. Empty input is special-cased — the host passes pointer
`0` with length `0` and does not allocate or free anything, so a guest that never
receives non-empty input is never asked to `allocate` at all.

A guest that returns `0`, a non-integer, or an out-of-range offset from `allocate`
fails the request. Returning `0` is the ABI's way for a guest to report that it could
not satisfy the allocation — but note that **TinyGo guests do not use it**: exceeding
the linear-memory cap makes a TinyGo allocation panic, so out-of-memory reaches the
host as a trap rather than as `allocate` returning `0`. The `0` path exists for guests
whose allocator can fail without panicking.

## Getting bytes out: the result protocol

Every result-returning export returns an `i32` offset into its own linear memory. At
that offset the host reads:

```
offset      +0 +1 +2 +3   +4 ...
           [ length (uint32, little-endian) ][ JSON body, `length` bytes ]
```

The JSON body unmarshals into one of the result structs in `pkg/wasm/common`:
`DeployResult`, `DepositResult`, `ProcessResult`. All carry
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

**Execution is time-bounded.** A single guest operation — the export call plus the
`allocate`/`deallocate` calls around it, and instantiation including any `start`
section — must complete within `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` (default
10 000 ms, maximum 300 000, `0` selects the default; there is no unlimited setting).
A guest exceeding it is interrupted with a trap and the request fails on-chain. The
bound is wall-clock, enforced by the host with wasmtime epoch interruption, and needs
nothing from the guest.

What appears on-chain is the `WASM_INTERNAL` error code with the message
`guest execution timed out`. There is no dedicated on-chain code: `ErrorCode` carries
the failure *category* (`RequestFailure.Category()`), and `WASM_INTERNAL` is shared
with other host-side WASM failures such as a failed memory write or an unusable
result. To recognise a timeout specifically, match the message, not the code.

Practical notes for guest authors:

- The budget covers the **whole operation**, not each call, and it is not reset
  between the host's `allocate` calls and your export.
- Time spent inside a host call is *not* interruptible: the bound is enforced by
  checks the compiler emits into guest code, which a blocked host call never
  reaches. This is why the imports a guest may declare are a closed set (see "Host
  imports: a closed set"); none of the permitted ones can block.
- The interrupt is a trap, so it evicts the module from the cache (see below) and no
  guest cleanup runs. Do not rely on deferred work completing.
- Real guests are nowhere near this: the reference apps return in single-digit
  milliseconds. Reaching the bound means a loop that does not terminate.

**A second, independent bound applies per request batch.** Requests are processed in
batches (up to `MAX_BATCH_SIZE`, one application per batch), and the caller sends one
*execution budget* covering the whole batch — derived from how long it will wait for
the response. So two clocks run, owned by different sides:

| Bound | Scope | Owner | On expiry |
|---|---|---|---|
| `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` | one guest operation | the enclave | **signed on-chain**: `WASM_INTERNAL` + `guest execution timed out` |
| caller execution budget | the whole batch | the caller, clamped by the enclave | not signed; the request stays pending and is retried |

The consequences a guest author can actually observe:

- **A timeout does not stop a batch.** It is one error payload; the requests after it
  in the same batch still run, against the state as it stood before the failed one.
- **The batch budget can cut a request short instead.** When a batch cannot finish in
  time, the requests that completed are settled and the remainder stay pending for a
  later batch — so a request may be started, abandoned, and later run again from
  scratch. Guest code must therefore not assume that being invoked means its result
  will be used; state is only ever advanced by a settled request. The executor avoids
  this where it can, skipping a request whose guest bound the remaining budget cannot
  cover, but a request that also carries a deposit makes two guest calls and can
  still be cut between them.
- **Only the enclave's own bound can fail a request on-chain.** A budget expiry is
  host-side abandonment, so it never charges a fee. This is deliberate: the budget
  crosses the TEE boundary, and a value supplied from outside must not be able to
  make a user pay for a failure it manufactured.
- **Sizing.** For every batch to settle at least one request per attempt, the enclave
  bound must be below the caller's budget — otherwise the budget always wins the race,
  nothing is ever signed, and the FIFO queue cannot advance. Worse, since nothing is
  signed nobody is charged, so a guest that spins past the caller's patience buys
  unbounded executor time for free. The executor therefore **refuses the connection**
  if `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` plus a 1 s safety margin is not strictly
  below `MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC` minus the 2 s reply margin:
  the manager reports its timeout in the handshake and both processes fail at boot with
  an error naming the two settings (see `docs/design/EXEC_MGR_HANDSHAKE.md`). The
  safety margin covers what the bound itself does not: the deadline is armed only after
  the request is decoded and the state decrypted, and epoch interruption overshoots it
  by up to two ticks. One bound is the whole cost of a request however many guest calls
  it makes — loading the module, taking a deposit and processing share a single
  per-request budget — so the bound is compared once, not once per call. For a
  *whole* batch to finish in one attempt you additionally want
  `MAX_BATCH_SIZE × EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` to fit inside the budget;
  exceeding that is not an error, it just costs extra polls.

**Table storage is bounded separately.** Table elements do not live in linear memory,
so the cap above does not cover them — an unbounded element count would let a module
commit gigabytes regardless of the memory setting (measured: ~390 MB resident from a
50M-entry funcref table under a 64 KiB memory cap). The store limiter therefore also
pins the element count and table count, adding at most a few tens of MiB per module,
which is why the sizing rule above still holds. A guest declaring more fails to
instantiate.

**The accepted WASM feature set is pinned explicitly**, not inherited from whatever
Wasmtime enables by default, because the Executor re-executes guest code and commits
the result on-chain. A module using a disabled proposal fails to compile.

| Enabled | Disabled |
|---|---|
| bulk memory, multi-value, reference types (**funcref only — see below**), fixed-width SIMD | relaxed SIMD (host-architecture-dependent results), memory64 (breaks the int32 offset ABI), multi-memory and threads (both break the per-memory RAM accounting), tail calls, function references, GC, wide arithmetic, exceptions, component model |

Cranelift NaN canonicalization is enabled, since NaN bit patterns are otherwise
host-dependent, and relaxed SIMD's deterministic mode is enabled as well — redundant
while the proposal itself is off, but it means enabling relaxed SIMD later cannot
reintroduce host-dependent results by omission. Two engine-wide subsystems are also
switched off — GC support and concurrency support — which are capabilities rather than
proposals.

`TestPinnedEngineRejectsDisabledProposals` asserts that a violating module is rejected
for each disabled proposal, and each case records what it guards — some pins overlap,
so a single fixture can guard a pair rather than one knob. A pin can only be covered
this way if a core module can express a violation of it, which rules out the component
model. Rather than trusting this paragraph to stay exhaustive, compare the table in the
test against the `Set*` calls in `newPinnedEngine` when adding a pin.

**Switching GC support off narrows reference types.** `externref` values need the GC
subsystem, so a module that merely declares an `externref` parameter is rejected even
though reference types are enabled. `funcref` — what TinyGo emits behind
`call_indirect` — is unaffected. A guest may therefore use reference types for
function references but not for host references.

That is a deliberate narrowing, and what keeps it safe is a requirement rather than a
coincidence: **no host import may take or return `externref`.** Today that holds
trivially, since the allowed import set is WASI-only. It must keep holding as host
functions are added — a host-crypto bridge is specified, and it passes only `i32`
pointers and lengths, so it is unaffected. If some future host import genuinely needs
to hand a guest an opaque host reference, GC support has to be switched back on in the
same change.

Adding a proposal is a deliberate change in `newPinnedEngine`, and
newly-introduced proposals should be pinned there when Wasmtime is upgraded.
Pinning `memory64` off is not only an ABI matter: it also neutralised
`GHSA-p8xm-42r7-89xg`, a host panic reachable only with `memory64` enabled.

One gap to be aware of: a few flags exist in Wasmtime's C API but are not wrapped by
`wasmtime-go`, so they cannot be pinned from Go and stay at upstream defaults. As of
v47 those are `shared_memory`, branch hinting, custom page sizes, stack switching, and
the component-model sub-flags — all documented `false` by default. `shared_memory` is
worth singling out: it gates shared-memory creation, which is the same per-memory RAM
accounting concern the threads pin exists to address, so it is unpinnable rather than
unimportant. Re-check when upgrading, in case the bindings begin exposing them.

**Guest execution is bounded in time but not in work.** Wall-clock duration is
bounded — see "Execution is time-bounded" above — so a guest that does not return is
interrupted rather than blocking the Executor indefinitely. Fuel (how much computation
a request may do) and stack depth are still unbounded, as is module compilation and
time spent inside a WASI host call, none of which epoch interruption can reach. See
the TODO in `newPinnedEngine`; tracked separately.

## Failure classification and the module cache

The host sorts failures into three classes, because the response differs. The test
is **whether guest code actually ran**, not whether the failure looks deterministic:

| Class | Examples | Module cache |
|---|---|---|
| Application error | non-empty `error` in a well-formed result | kept — a healthy guest declining a request |
| Guest fault | a trap; `allocate` failing, returning `0`, or returning an out-of-range offset; a result the host cannot parse | **evicted** — the heap is in an unknown state and must not serve the next request |
| Host-side / static defect | a missing export; an export whose signature does not match the host ABI; a nil store | kept — nothing executed, so re-instantiating provably cannot change the outcome |
| Execution bound exceeded | the guest ran past `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` and was interrupted | **evicted** — an interrupted guest's heap is in an unknown state, exactly like a trap |

The execution bound splits into two on-chain outcomes, which differ in who caused
the interruption:

| Cause | Reported as | Effect |
|---|---|---|
| The guest spent its whole budget | signed on-chain: `WASM_INTERNAL` + `guest execution timed out` | request settles as failed, minimum fee charged |
| The host stopped waiting (executor shutdown, cancelled request) | not signed at all | request stays pending and is retried later |

The second case never reaches the chain by design: the request was abandoned, not
judged, so charging for it would be wrong. As a guest author you cannot tell the two
apart and do not need to — the retry is a fresh call on a fresh instance.

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

- **WAT fixtures (`pkg/wasm/wasmtime_runtime_test.go`)** assemble tiny modules from
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

- **TinyGo integration tests (`app/simple/integration_test.go`)** compile the real
guest and exercise the real allocator, real serialisation and real business logic.

  Neither layer is sufficient alone. A WAT fixture encodes what we *believe* the ABI
is, so it can pass while a real guest fails; the integration tests are the backstop
that catches ABI drift, and they are also what proves the pinned feature set does not
reject legitimate modules.
