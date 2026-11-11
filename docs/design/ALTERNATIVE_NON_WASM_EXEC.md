# HZN-2871 — Non-WASM Executor Path: Analysis

**Status:** Analysis / proposal for review (no implementation yet).
**Ticket scope:**

> We want a non-wasm executor.
> - Try to keep the outer logic as-is as possible (the discriminator old/new logic could be a specific APPID).
> - Should produce an off-chain result and publish on-chain only an "anchor" (hash).

This document describes the proposed shape: a discriminator-driven dispatch in the executor that admits handler shapes other than pure WASM — pure-native Go handlers compiled into the executor, or WASM apps that reach native services through host imports — selected by a per-app field in the deploy descriptor. Off-chain results are served through the existing authority service; only a cryptographic hash is anchored on-chain. The motivating concrete handler is in-enclave LLM inference.

## 1. Problem and goals

Today the executor is 100% WASM-bound. Every request funnels into `StatelessExecutor.HandleProcessRequest`, which calls `e.runtime.ProcessRequest(...)` on the Wasmtime-backed runtime. The WASM module receives the decrypted payload, runs deterministically under fuel metering, and produces the new encrypted state plus a signed update payload.

Some workloads do not fit WASM well:

- **LLM inference** — inference engines (`llama.cpp` and friends) target native CPU intrinsics (AVX2/AVX-512) and OS threads. A pure-WASM port loses 2–4× to SIMD-128 and gives up multi-threading.
- **Heavy crypto, ML, oracle / remote-API adapters, OS primitives beyond WASI** — same shape: code that cannot reasonably (or efficiently) be expressed under the WASM sandbox.

### Goals

- **Outer pipeline unchanged.** Request fetching, P521 user-key decryption, per-app versioned state, TEE signing, fee math, handshake, subgraph events, deploy flow — all stay as they are.
- **Per-app discriminator on-chain.** Each deployed app carries an `AppType` attribute that selects the execution path. Default `0 = WASM` for back-compat.
- **Off-chain result + on-chain anchor.** Result blobs (potentially large) are stored off-chain at the authority and served via the existing `/getreport` pattern (the authority's user-signature-gated retrieval endpoint, already used for deanonymization reports); on-chain we publish only `keccak256(resultCipher || teeSig)`.
- **TEE-bound results.** Every result carries the executor's secp256k1 signature, verifiable against `TeeAuthenticator`'s PCR allow-list.
- **Static handler registration.** Handlers are compiled into the executor binary; adding/modifying one produces a new enclave image with new PCRs (covered by the existing attestation gate).
- **Bounded execution.** Each handler runs under a per-request resource cap (wall-clock timeout in v1; a fuel analog in v2).

### Non-goals

- Hot-loading native handlers at runtime.
- Replacing the WASM path. Both coexist.
- Building the LLM handler itself — this is the framework into which it plugs.
- Full decentralized result storage (IPFS/Arweave) in v1 — called out as a v2 path.

> Throughout this doc, **v1** refers to the first shippable iteration of the framework (foundations + dispatch + first handler); **v2** refers to follow-up enhancements (richer fee analogs, decentralized result storage, additional handlers).

## 2. Discriminator: on-chain `AppType`

The deploy descriptor gains two fields:

```go
// pkg/common/types.go (new)
type AppType uint8

const (
    AppTypeWasm            AppType = 0  // existing behaviour; default for back-compat
    AppTypeNativeEcho      AppType = 1  // smoke-test handler (identity)
    AppTypeNativeLLMLocal  AppType = 2  // in-enclave LLM (deterministic)
    AppTypeNativeLLMRemote AppType = 3  // remote LLM via pinned HTTPS endpoint (oracle semantics)
    // future handlers claim their own values
)
```

```go
// pkg/common/deploy_descriptor.go (extension)
type DeployDescriptor struct {
    AppType        AppType   // NEW: dispatch discriminator. 0 = WASM (default).
    HandlerVersion uint32    // NEW: bumped on handler-semantics changes; 0 for WASM.
    // ... existing fields (constructor params, WASM fingerprint, etc.) ...
}
```

`ProcessorEndpoint.sol` stores `appType[appID]` at deploy time, emits it in `DeployRequestSubmitted`, and the manager/executor read it on every request. The subgraph indexes it on the app entity. Existing deploys all read as `appType = 0`, so no migration is required.

**Why on-chain rather than executor config:** single source of truth, auditable, no silent drift between contract state and executor wiring, no executor restart on each native deploy.

`HandlerVersion` is bumped when a handler's semantics change. This prevents silent semantic drift of native handlers behind users' backs (the WASM path already has this for free: deploys carry a hash of the bytecode in the descriptor and the executor refuses to load anything that doesn't match).

## 3. Motivating use case — LLM inference

Two trust shapes are first-class; the choice is per-`AppType`, visible to the user before they submit a request.

### 3.1 In-enclave LLM (CPU-only, RAM-resident model)

Model weights ship inside the executor enclave image (or are pulled into an attested volume at startup — i.e. a host-mounted volume whose content hash is included in the enclave's PCR measurement chain). Inference runs in-process on Nitro CPU/RAM. **No outbound network from the TEE.**

- **Trust:** strong. Weights, prompt, output are all sealed inside the TEE; the existing PCR-based attestation covers the model file as part of the enclave image.
- **Cost:** significant enclave memory; cold-start penalty; bounded model size (≈ 7B parameter class, heavily quantized) given Nitro's CPU-only constraints.

### 3.2 Remote LLM (outbound HTTPS call)

The enclave calls a hosted LLM endpoint (OpenAI / Anthropic / Bedrock / a self-hosted GPU node) over TLS. The TEE treats the response as an opaque oracle answer and signs over it.

- **Trust:** weaker. A new outbound trust boundary; verifier ends up trusting "this enclave talked to `llm.example.com` and got this back" rather than "this specific model produced this from this prompt".
- **Mitigations** (mandatory for `AppTypeNativeLLMRemote`):
  - Pin the endpoint URL + TLS CA fingerprint + expected model ID in code (PCR-covered).
  - Sign over `(prompt_hash, response)` — not just the response — so the enclave cannot silently substitute either side.
  - Forward provider-side attestation when available (NVIDIA Confidential Computing, Bedrock attested outputs, vendor-signed responses) alongside the TEE signature.
- **Cost:** small enclave footprint, fast, no model-size ceiling — at the price of a real trust downgrade.

### 3.3 In-enclave LLM — recommended path: host-bridged architecture (pending spike)

For the in-enclave LLM variant, the recommended shape is to **split app logic from inference and bridge them through host imports**: the app stays a small WASM module (same as apps today), while the inference engine (`llama.cpp` or equivalent) is statically linked into the executor binary via cgo and reused across all requests. The host exposes a single custom WASM import that reads the prompt from the app's linear memory, dispatches into the native engine, and copies the generated bytes back. Model weights live in the host's heap (not in WASM linear memory — wasm32 caps a single module's linear memory at 4 GB — so not bound by that ceiling) and their hash is enrolled into the enclave PCRs at startup, so attestation binds to a specific *(engine, weights)* pair.

This preserves the multi-tenant property of the platform (third-party app code stays WASM-sandboxed) while letting inference run at full native speed (AVX2/AVX-512, OpenMP), and it requires no new artifact type on-chain — apps still deploy a small `.wasm`. A pure-native handler remains the right answer for non-LLM use cases (heavy crypto, OS primitives, whole-request native workloads, remote-API adapters) and is the v1 fallback if upstream-runtime work blocks the host-bridge path.

**How this plugs into the framework.** If the spike validates this path, the LLM `AppType` is served *inside the existing Wasmtime runtime* via additional host imports — not via `NativeRuntime`. The framework pieces in §4–§9 still apply: the on-chain `AppType` discriminator, deploy-descriptor field, `ResultAnchored` event, and authority retrieval flow are independent of which runtime ends up handling the request. Only the dispatch target differs — for the LLM `AppType`, the dispatcher routes back to `WasmRuntime` (with the inference host-import wired in) rather than to `NativeRuntime`. If the spike fails, that same `AppType` is served by `NativeRuntime` as a pure-native handler instead.

> **The recommendation is conditional on a dedicated spike confirming runtime support and throughput.** The project pins `wasmtime-go v1.0.0` (September 2022); current upstream is several majors ahead. The spike needs to confirm that custom host imports with the `(*Caller, args...) → result` shape and per-`Store` concurrency behave as expected on the chosen target version, and to measure end-to-end inference throughput on Linux CPU as a leading indicator for Nitro. Out of scope of this analysis.

## 4. Component view

```
  ┌──────────────────────────────────────────────────────────────────┐
  │ Executor (inside Nitro enclave)                                  │
  │                                                                  │
  │   StatelessExecutor.HandleProcessRequest                         │
  │      │                                                           │
  │   [decrypt payload + load encrypted app state]  ◄── unchanged    │
  │      │                                                           │
  │   ┌──▼──────────────────────────────────────────┐                │
  │   │ Dispatcher (NEW — thin)                     │                │
  │   │   switch appType[appID] {                   │                │
  │   │     case AppTypeWasm: wasmRuntime.Process   │                │
  │   │     default:          nativeRuntime.Process │                │
  │   │   }                                         │                │
  │   └──┬──────────────────────────────┬───────────┘                │
  │      │                              │                            │
  │   ┌──▼──────────┐            ┌──────▼─────────┐                  │
  │   │ Wasmtime    │            │ NativeRuntime  │                  │
  │   │ runtime     │            │   handlerReg:  │                  │
  │   │ (existing)  │            │   AppType →    │                  │
  │   │             │            │   NativeHandler│                  │
  │   └─────────────┘            └────┬───────────┘                  │
  │                                   │                              │
  │                              ┌────▼────────────┐                 │
  │                              │ EchoHandler     │                 │
  │                              │ <future native  │                 │
  │                              │  handlers>      │                 │
  │                              └─────────────────┘                 │
  │      │                                                           │
  │   [commit new state + sign stateUpdate]      ◄── unchanged       │
  │      │                                                           │
  │   [write result ciphertext to authority volume] ◄── NEW          │
  └──────────────────────────────────────────────────────────────────┘
```

The **Dispatcher** is the only new component on the hot path. Everything else is reused (Wasmtime runtime, crypto, storage, communication) or lives in its own module (`NativeRuntime` + handlers).

## 5. NativeHandler interface

```go
// pkg/executor/native/handler.go (new)

type NativeHandler interface {
    // Execute runs the handler against a decrypted payload + current state.
    // Returns the new state and the result blob to be anchored + served off-chain.
    // Errors here become signed executionErrors on the stateUpdate.
    Execute(ctx context.Context, in HandlerInput) (HandlerOutput, error)

    // Deterministic signals whether replay against the same prior state yields
    // the same result. Surfaced on-chain via the AppType attribute.
    Deterministic() bool

    // MinEnclaveMemoryMB declares the minimum enclave RAM the handler needs to run safely.
    // Checked once at executor startup against the actual enclave size:
    //   - if the check passes, the handler is loaded into the registry as normal;
    //   - if it fails, the AppType is registered as disabled — requests against it
    //     fail-fast with a signed errorCode = HandlerUnavailable, so the rest of
    //     the handler set stays live.
    // This is a startup admission check only; per-request memory bounds are §8's job.
    MinEnclaveMemoryMB() uint32
}

type HandlerInput struct {
    AppID     common.ApplicationID
    RequestID common.RequestID
    Sender    common.Address
    Payload   []byte           // decrypted user payload
    State     []byte           // decrypted prior state (nil for first request)
    Deadline  time.Time        // honored via ctx; exposed for handler-side pacing
}

type HandlerOutput struct {
    NewState []byte           // will be AES-encrypted by the caller. For oracle-semantics
                              // handlers (Deterministic() == false), the framework wraps
                              // this with state_{N+1} = commit(state_N, requestID, hash(Result))
                              // before encryption — handlers may return nil.
    Result   []byte           // ciphertext bound to user key; anchored on-chain, served off-chain
    FuelUsed uint64           // fuel-equivalent units, plugged into existing fee math
    Metadata map[string]string // optional, surfaced by subgraph (tokensIn, tokensOut, modelID, …)
}
```

Handlers are registered at executor startup:

```go
// pkg/executor/native/registry.go (new)
var defaultRegistry = map[common.AppType]NativeHandler{
    common.AppTypeNativeEcho: echo.NewHandler(),
    // common.AppTypeNativeLLMLocal:  populated only if the §3.3 spike picks the
    //                                pure-native fallback; otherwise this AppType
    //                                is served by WasmRuntime via host import and
    //                                does not appear in this registry at all.
    // common.AppTypeNativeLLMRemote: llm.NewRemoteHandler(...),
}
```

The state-store interface gains exactly two methods to support the new dispatch:

```go
// pkg/storage/interface.go (extension)
type ApplicationStateStore interface {
    // ... existing methods unchanged ...

    // StoreNativeApp records a deploy whose handler ships in the executor binary.
    // No bytecode to persist; AppType + HandlerVersion fix the handler identity.
    StoreNativeApp(ctx context.Context, appID ApplicationID, appType AppType, handlerVersion uint32, initialState []byte) error

    // GetAppType returns the AppType recorded at deploy time.
    // Returns AppTypeWasm for any pre-existing deploy.
    GetAppType(appID ApplicationID) (AppType, error)
}
```

## 6. Off-chain result + on-chain anchor

### 6.1 Anchor

```
anchor = keccak256( resultCiphertext || teeSignature )
```

emitted from the `stateUpdate` flow as a new event:

```solidity
event ResultAnchored(uint256 indexed requestID, bytes32 indexed appID, bytes32 anchor);
```

The subgraph indexes `ResultAnchored` so clients can locate the authority + anchor by `requestID` without a separate directory.

### 6.2 Off-chain retrieval flow (v1)

```
Client: GET  /nonce?requestID=0x…
        →    { nonce, expiresAt }
Client: POST /getreport
        body: { requestID, userAddr, sig(nonce || requestID, userKey) }
        →    { resultCiphertext, teeSignature, enclavePCRs, anchor }
```

Three independently verifiable properties for a verifier:

1. **Anchor integrity** — on-chain `anchor == keccak256(resultCipher || teeSig)`; the authority cannot swap the bytes.
2. **TEE provenance** — `teeSig` recovers to an address allow-listed in `TeeAuthenticator`, and the PCRs match an attested enclave image.
3. **User-binding privacy** — `resultCipher` decrypts only with the user's registered key.

Result storage in v1 reuses the existing authority volume (LevelDB / shared folder). v2 upgrade path: IPFS-ciphertext + TEE-held key released via the same authority endpoint — clean migration because the public API does not change, only the storage backend.

### 6.3 Determinism, replay, rollback

The versioned LevelDB model assumes that a given prior state + request produces a deterministic next state. Two stances supported per handler:

| Stance | Mechanism | Replay semantics |
|---|---|---|
| **Deterministic** (`Deterministic() bool → true`) | Handler is a pure function of `(payload, prior state)`; any non-deterministic source (RNG, FP reduction order, external clock) is pinned or removed | Bit-exact replay; rollback semantics identical to WASM |
| **Oracle semantics** (`Deterministic() bool → false`) | Treat result as opaque oracle answer; framework synthesizes `state_{N+1} = commit(state_N, requestID, hash(result_N))` before encryption, regardless of what the handler returned in `NewState` | Rollback restores prior state; re-execution may produce a different result, philosophically fine for an oracle but visible to the app |

The on-chain `AppType` exposes the determinism property to users *before* they submit, so the choice is auditable per app.

## 7. Pay-per-use accounting

The existing fee formula stays as it is:

```
applicationFee = max(fuelUsed × FuelPricePerUnit, MinFeePerRequest)
refundAmount   = maxFee − applicationFee
```

What changes is the **source** of `fuelUsed`. Today there is a single source — `wasmtime-go`'s instruction-level fuel meter. With native handlers in the picture, `fuelUsed` becomes whatever the handler reports as its `FuelUsed` value, with `MinFeePerRequest` as the floor when a handler has no useful metric to report.

### 7.1 Where the cost number comes from per handler shape

Different handler shapes have different natural cost units. The framework treats `fuelUsed` as an opaque non-negative integer; calibration is per-handler.

| Handler shape | Natural cost unit | Notes |
|---|---|---|
| **Existing WASM app** | WASM instructions executed | Reported by Wasmtime's fuel meter, unchanged |
| **WASM app with host-bridged service** (e.g. §3.3 LLM) | Tiny WASM fuel + a handler-declared metric for the bridged work (e.g. tokens in/out) | Wasmtime's fuel meter pauses inside host imports, so heavy work done in cgo is invisible to it; the handler returns the bridged-work metric explicitly |
| **Pure-native handler with a quantifiable cost** | Whatever metric the handler can measure cheaply (CPU-time bucket, output bytes, API calls made, etc.) | Reported as `HandlerOutput.FuelUsed` |
| **Pure-native handler with no useful metric** | None | Returns `FuelUsed = 0`; the formula falls through to `MinFeePerRequest` per call |

For a host-bridged LLM handler the natural composite is something like `tokensIn × c_in + tokensOut × c_out + wasm_instr / k`. For a remote-API adapter it might be `apiCalls × c_call`. For a one-shot crypto computation it might just be the flat `MinFeePerRequest`. The point is that the framework does not pick: each handler declares what makes sense for its workload, and the calibration constants are part of the handler's PCR-covered identity (changing them requires a `HandlerVersion` bump).

**Calibration is a platform-governance concern, not a per-handler-author judgement.** Picking `c_in`, `c_out`, `k`, or any other handler-specific constant by eyeball risks economic distortions in either direction (under-pricing burns operator margin; over-pricing surprises users). Before exposing a new handler's costs to user billing, the platform needs a calibration process — measuring the handler against representative workloads, comparing to existing units, sanity-checking against operator cost — and a review step. This is a v2 concern (the v1 echo handler bills at `MinFeePerRequest` and side-steps the question), but it must be in place before any non-trivial handler ships.

### 7.2 The blind spot: work inside host imports

When a WASM guest calls a host import, Wasmtime's fuel meter pauses; everything the host does in native code is invisible to it. For a host-bridged LLM call the gap is extreme — a few WASM instructions worth of `i32` push/return on the guest side (~10³ fuel units), vs. native matmuls that would have been billions of instructions had the engine been compiled to WASM (~10⁹–10¹⁰ equivalent units for a typical 7B-class forward pass). The resulting ~6–7 order-of-magnitude gap is an estimate, not a measurement, but the qualitative point holds: fuel-only billing on this shape would undercharge dramatically.

This is why the host-bridged shape needs a handler-declared metric on top of the (tiny) WASM-fuel number; the framework just plumbs both into the same `FuelUsed` total.

### 7.3 Why WASM fuel still matters even when it is small

For any handler shape that runs **WASM code at all** (the existing path, and the §3.3 host-bridged path), fuel remains load-bearing for the platform shape regardless of its weight in the final fee:

- **DoS protection.** A guest stuck in a tight WASM loop without ever calling a host import would otherwise spin forever. A per-invocation fuel budget (`Store.SetFuel(N)`) traps the guest with `OutOfFuel` when exhausted; the host catches the trap, bills the failed call, and moves on.
- **Pricing of guest-side compute.** Data-processing guests (CSV transformations, regex over large payloads, prompt-construction templating between bridged calls) do real WASM work that the existing fuel meter measures naturally. Per-call billing alone would let them run arbitrarily long for free.

For pure-native handlers (no WASM in the loop) neither of those concerns applies; the request-level timeout in §8 is the corresponding safety floor, and `MinFeePerRequest` is the corresponding billing floor.

## 8. Resource limits & failure paths

| Outcome | On-chain effect | Fee |
|---|---|---|
| Handler success, `FuelUsed = F` | stateUpdate with new state + anchor; `ResultAnchored` event | `max(F × FuelPricePerUnit, MinFeePerRequest)`; `refundAmount = maxFee − this` |
| Handler returns error | Signed executionError payload on stateUpdate (same shape as WASM error) | `maxFee`, no refund |
| Handler panics (`recover()` catches) | Signed `errorCode = HandlerPanic` | `maxFee`, no refund |
| Timeout (context cancellation) | Signed `errorCode = HandlerTimeout` | `maxFee`, no refund |

All four paths produce a signed stateUpdate — no request silently disappears. Rollback semantics are identical to WASM: failures do not commit new state.

**v1 cap:** `context.WithTimeout` per request + mandatory `defer recover()` around handler invocation.

**v2 cap:** add memory-sampling safeguard via `runtime.MemStats` if any handler OOMs in practice. Deferred — handlers are attested code inside the enclave; size enclave RAM generously and rely on Nitro's OOM killer as the last line.

**Goroutine hygiene (mandatory from v1):** handlers must bind any spawned work to the request `context.Context`. On timeout the context is cancelled and the handler contract requires all spawned work to tear down. Violations are a handler-authoring bug.

## 9. Trust model

| Property | WASM today | Native in-enclave handler | Native handler with outbound network |
|---|---|---|---|
| User input confidentiality | In enclave | In enclave | In enclave (input hash committed) |
| Computation integrity | WASM fuel-metered, attested bytes | Attested Go (+ any cgo/native deps), timeout-bounded | Attested Go code; output trust reduces to the remote endpoint's trust |
| Output confidentiality | User-key encrypted | User-key encrypted | User-key encrypted |
| Output authenticity | TEE sig | TEE sig | TEE sig over `(input_hash, response)` |
| Outbound network from TEE | None | None | Pinned endpoints only |

What ends up newly enrolled in the PCR set:

- **Handler registry.** The static `map[AppType]NativeHandler` is in code, therefore in PCRs. No handler can be swapped in without a new enclave image.
- **Per-handler static assets.** Anything a handler needs at runtime that ships with the enclave image (data files, native archives, compiled constants, pinned endpoints / TLS roots for handlers that reach out) is part of the same measurement chain — there is no separate signing step.

This gives "handler fingerprint = enclave fingerprint" for free — no per-handler signing, no extra verification step beyond what `TeeAuthenticator` already does.

## 10. What stays invariant (sanity check)

These do not change for either path:

1. Request dispatch in `communication.ClientConnection.handleClientRequest`.
2. User-key decryption in `StatelessExecutor.decryptPayload` (P521).
3. Per-app AES-encrypted state via `VersionedLevelDBAppStateStore`.
4. Keyset recovery handshake.
5. Fee math in `ProcessorEndpoint`.
6. stateUpdate submission path + signature verification by `TeeAuthenticator`.
7. Subgraph event indexing.
8. Deploy role-gating (`DEPLOYER_ROLE`).

The only on-chain additions are:
- `appType` + `handlerVersion` fields on the deploy descriptor + event.
- `ResultAnchored(requestID, appID, anchor)` event.
- One state-store interface method (`StoreNativeApp`) + one read (`GetAppType`).

## 11. Open questions / risks

- **`wasmtime-go` version delta.** The project pins `v1.0.0` (Sept 2022); current upstream is several majors ahead. Any handler shape that relies on advanced runtime features (e.g. the §3.3 host-bridged path) needs to confirm those features on the chosen target version before commitment. A short capability spike is the right next step (out of scope of this analysis). Independently of §3.3, sitting on a 2022 runtime is a standalone maintenance liability — the upgrade should happen at some point even if the LLM ambition stalls.
- **Non-determinism sources within a handler.** Even handlers that aim to be deterministic can leak non-determinism via the platform (FP reduction order across CPU generations, library-internal RNG, scheduler-dependent iteration). Mitigation: pin the handler to a specific instruction set / library version, all part of the PCRs; surface determinism as a per-handler claim verified by tests, not assumed.
- **External-dependency drift.** Handlers that reach a third party (remote API, hosted model, oracle) can have their behaviour change silently when the remote side updates. `HandlerVersion` lets us bump on intentional changes but cannot detect drift on the other end. No perfect mitigation without vendor-side attestation.
- **Fuel-unit calibration governance.** Per-handler cost constants (`c_in`, `c_out`, etc.) are PCR-covered once chosen, but the doc is silent on *how* they get chosen. Sloppy calibration is an economic surface, not just a UX issue: it can under-charge by integer factors (operator loss) or over-charge (user surprise). Before any non-trivial handler bills users, the platform needs a documented calibration process — benchmark against representative workloads, sanity-check against operator cost, and gate behind review (see §7).
- **Authority single point of availability.** If the authority node loses the ciphertext, the on-chain anchor is useless. The v2 IPFS upgrade *changes the trust model* (from a single authority node to one or more pinning operators) but does not by itself guarantee availability — IPFS content disappears when the last pinning peer drops it. Solving availability needs complementary economic incentives (paid pinning, replication SLAs, retrieval markets) on top of whichever storage layer is chosen.
- **Authority metadata visibility.** Even though result ciphertext is unreadable to the authority, the authority observes which user fetches which `requestID` and when. For privacy-sensitive applications this is a leak the existing infra already accepts; flagged here so it isn't forgotten when sizing the v2 IPFS path (CID lookups have the same shape).
- **Memory caps.** v1 accepts "trust the attested handler" — revisit only if a handler in the wild OOMs in a way that harms other apps.

## 12. Suggested phasing (high-level)

1. **Foundations.** Add `AppType` + `HandlerVersion` to the deploy descriptor, contract, bindings, subgraph, state-store interface. Existing apps default to `AppType = 0`. No behavior change.
2. **Dispatch + skeleton.** `NativeRuntime` package, `NativeHandler` interface, dispatcher in `HandleProcessRequest` / `HandleDeployApp`. `EchoHandler` as smoke handler. Timeout + recover wrappers.
3. **Result anchoring + authority.** `ResultAnchored` event, authority `/getresult` extension (or `/getreport` extension), subgraph indexing, fullstack test.
4. **First substantive handler.** Concrete target is in-enclave LLM inference (per §3); the architectural shape is gated on the §3.3 spike. A pure-native handler is the available fallback if that path does not pan out.
5. **Hardening.** Calibrate the per-handler cost formulas. Memory-sampling safeguard if needed. Admin-CLI handler introspection. v2 IPFS-backed result storage. Additional handlers as use cases emerge.

The WASM path stays fully functional throughout. Each phase is independently shippable.
