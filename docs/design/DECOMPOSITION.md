# Decomposing Vela: chain-agnostic core + EVM plugin

**Status:** evaluation / discussion draft — no code changed.
**Date:** 2026-08-06.
**Question:** can Vela be split into (a) a chain-agnostic core that runs generic WASM in a Nitro enclave — requests, results, attestation, runtime, communication, cryptography, TEE key handling, plugin support — and (b) an EVM plugin holding the contracts plus the chain-facing part of the Manager?

---

## 0. Decided scope

Five points are settled, in order to narrow the problem:

1. **Target set: multiple EVM chains + a "no-chain" mode.** The no-chain mode runs WASM in the TEE with attestations and logging written to the host filesystem. For now other types of chain are not evaluated (Solana, UTXO chains, non-EVM).
2. **Per-target enclave images are acceptable** — different PCR0, separate deployment and attestation ceremony.
3. **Request scheduling moves into the core in both modes.** The core owns cross-app selection; the contract keeps only the invariants that are about correctness — per-app FIFO, head-of-queue, `prevStateRoot` chaining — and drops the round-robin cursor. On EVM that makes cross-app fairness *detectable* rather than *enforced*; in no-chain mode there is no contract, so the core owns it outright. §3.2 works through what this changes.
4. **Every deployment gets its own environment** — its own enclave, its own keyset. Deployments are never co-tenanted and never share a signing identity.
5. **Identity of the users is a 20-byte principal derived from a secp256k1 public key** — an Ethereum-style address — in **both** modes, no-chain included. Consequences:
   - `AppData`'s serialization is unchanged, so **there is no state migration**. The remaining work is relocating an `Address` type out of go-ethereum, which is mechanical.
   - Keeping it *keypair-derived* rather than an opaque 20-byte identifier is deliberate: an address that is a commitment to a public key can authenticate itself with no chain present, and it keeps wallet compatibility for free. An opaque ID would need a fresh authentication mechanism invented for it.
   - Note that the WASM part already works like this: the WASM host ABI passes sender, token and value as `(ptr, len)` pairs, and the guest side already avoids go-ethereum because TinyGo cannot compile it.
   - **This settles the representation only.** Who vouches that a request actually came from that principal is a separate question, and will be discussed in §2.2.

**The most important consequence of (1): you do not need a plugin framework.** Different EVM chains differ in *configuration* — RPC URL, chain ID, contract addresses, gas and finality policy — not in code. The genuine axis of variation is **chain vs no-chain**, a set of size two. Building a general chain-abstraction SPI for a set of size two would be over-engineering, and the estimate below assumes you don't.

---

## 1. Verdict

**Feasible, and the codebase is structurally most of the way there.** The caveats are not about the split itself but about two guarantees the chain has been supplying invisibly — freshness (§2.1), sender authentication (§2.2), client notification (§2.3). The decomposition surfaces them rather than causing them; they become visible the moment you plan to run without a chain.

Two facts set the tone:

1. **The seam already exists.** `pkg/blockchain/interface.go` is 52 lines and 15 methods, and everything above it — the Manager's poll loop, state versioning, rollback — is written against that interface, not against go-ethereum. `cmd/executor/main.go` has zero chain imports today.

2. **The coupling is narrow in code but wide in concept.** `pkg/manager` has *zero* go-ethereum imports and is still the most chain-aware non-adapter package in the tree. Removing imports is easy; removing assumptions is not.

---

## 2. The no-chain mode: what the chain was actually providing

Worth enumerating, because "just run WASM in the TEE" is a bigger subtraction than it looks. The chain currently supplies seven distinct services:

| Service | Provided by | Replacement in no-chain mode |
|---|---|---|
| Request intake and ordering | `_requestQueue`, FIFO enforcement | Intake service + durable queue (§2.3); ordering by the core scheduler (§3.2) |
| **Freshness / anti-rollback** | `applicationStateRoots[appId]` vs `prevStateRoot` (`ProcessorEndpoint.sol:335`) | **Nothing — §2.1** |
| **Sender authentication** | `msg.sender`, or a verified facilitator signature | **Nothing — §2.2** |
| Result durability and audit | Events, subgraph | Signed append-only log + attestation records |
| Client notification | `RequestSubmitted` / `RequestCompleted` / `UserEvent`, indexed by the subgraph | Read API over the commit log; submission acknowledged synchronously by the intake service (§2.3) |
| Signer authenticity | `TeeAuthenticator` PCR check, `updateTee` | Out-of-band verification of the attestation document |
| Economics (fees, custody, withdrawals) | ERC-20 ledger, refunds | Disabled (`FuelPricePerUnit`/`MinFeePerRequest` = 0, no withdrawals) |

Five of the seven have straightforward replacements. **Freshness and sender authentication do not** — and both are invisible today precisely because the chain has been supplying them silently. Each is analysed in the sections that follow.

### 2.1 Anti-rollback without a chain

The enclave verifies only that the state it is handed is *internally consistent* — `sha256(decryptedState) == StateRoot` (`pkg/executor/executor.go:761-764`). Nothing tells it whether that state is *current*: the executor is stateless, and `AppData`'s `appNonce` is incremented and serialized but never checked against anything. Freshness comes entirely from the chain — `ProcessorEndpoint.sol:335`, `if (prevStateRoot != applicationStateRoots[applicationId]) revert InvalidStateRoot();`. Remove the chain and the untrusted host becomes the authority on its own state, free to feed the enclave any historical `(state, stateRoot)` pair and have it processed. With no custody contract there are no funds to steal directly, so the damage is reverted application state — an unspent balance, a re-run deanonymization — an integrity failure whose severity depends on what the app does and whether it has external side effects.

**The replacement: an independent, stateful, commit-time witness.** It holds the current root per app, hands it out on request, rejects any submission whose `prevStateRoot` is not its record, then advances. That is `ProcessorEndpoint.sol:335` and `:549` moved off-chain — a signature check plus a monotonic root store, roughly 50 lines — and it needs **no enclave changes**, because the enclave never validated freshness in EVM mode either: it signs the transition and the *chain* refuses to commit a bad one (§2.4). The witness advances by following a log composed of the signed `UpdatePayload`s the enclave already produces: same artifact as the EVM path, same codec, verifiable by anyone holding the enclave's attested public key. The witness attests ordering, not correctness — all it can do, since it cannot execute WASM or decrypt state, and all that is needed, since `UpdatePayload` already carries the `prevStateRoot → newStateRoot` chain.

**How far to take it.** Start with one independent witness — separate party, separate infrastructure, separate keys. That turns rollback from "compromise the operator" into "compromise the operator *and* the witness", a 1-of-1 → 2-of-2 step that is most of the practical value and where I would stop (+1–2 weeks). Possible future evolution: an N-of-M quorum with gossip and overlapping quorums (+4–6 weeks, plus operating the set).

The witness cannot stop a compromised manager from acting — only refuse to record the result. That keeps a rollback or a fork from ever becoming committed state. Beyond that the log is evidentiary: since `processedRequestId` is part of the signed payload, the commit chain already proves which requests ran and in what order.

### 2.2 Sender authentication: the second thing the chain was quietly doing

**The enclave never verifies a user signature.** There is no `SigToPub`, `Ecrecover` or `VerifySignature` anywhere in `pkg/executor` or `pkg/crypto`, and `validateRequest` (`pkg/executor/executor.go`) checks only the protocol version and the minimum fee. It trusts `Request.Sender` because the contract set it from `msg.sender` — or from a facilitator signature the contract verified. The code says so plainly about TRUSTPROCESS: *"They are authenticated on-chain."*

That is entirely reasonable with a chain in front. Without one, the manager supplies `Sender`, and the manager is untrusted.

**Most of the surface is saved by an accident of the crypto design.** `decryptPayload` (`pkg/executor/executor.go:1017-1035`) looks up `keyStore[sender]` and decrypts with `ECDH(user_pub, enclave_priv)`. A manager claiming *sender = Alice* cannot produce a ciphertext that decrypts under Alice's key without Alice's private key. So `Process` and `Deanonymize` are **implicitly authenticated** — the decryption is the authentication.

**`AssociateKey` is the hole, and it is the one that matters**, because it is the bootstrap that binds an address to a P521 key:

```go
workData.AddKey(req.Sender, *keyToAssociate)   // execute_request.go:107 — no signature checked
```

The payload is plaintext, and nothing ties it to the claimed sender. In no-chain mode an untrusted manager could associate *its own* P521 key with Alice's address and from that point own her account outright — including reading every event encrypted "to Alice". Every other protection downstream is only as strong as this single binding.

**The fix keeps the 20-byte address meaningful rather than replacing it.** Require the user to sign the `AssociateKey` payload with the secp256k1 key their address is derived from, and have the enclave verify that the signature recovers to `Sender`. The address then authenticates *itself*, with no chain involved, because it is a commitment to that public key. The machinery already exists in-repo — `pkg/authorityservice/service.go:256-286` performs exactly this EIP-191 recovery.

Two notes on scope:

- **Do it in both modes, not just no-chain.** In EVM it is redundant (the chain already authenticated the sender) but harmless, and it means one code path instead of a mode-conditional one — the same symmetry argument as the commit-time witness in §2.1.
- **It is a wire-format change.** The `AssociateKey` payload is currently a hardcoded 133 bytes, or 226 with the encrypted seed (`execute_request.go:91-94`); adding a 65-byte signature makes it 198/291. `Request.ProtocolVersion` exists and is pinned at `0` (`executor.go:32`), so there is a versioning path. Budget a few days, inside P5.

### 2.3 Client interface: submission and notification

**Submission.** In a no-chain environment, a new service has to be designed to collect requests. Three notes on this:

- **Identifier assignment.** `requestId` is contract-derived today, and the manager recovers it by scanning the receipt (`pkg/blockchain/client.go:441-449`). Off-chain a similar logic must be recreated inside the service, that must assign it and return it in the acknowledgement: it is the client's polling key and it appears in the signed `UpdatePayload`.
Idempotency must be implemented either client-supplied or using the `requestId`. Without it, retries client side double-execute.
- **Backpressure.** Fees were the spam control. Without fees,the queue needs a depth cap and per-sender rate limiting, or the first hostile client fills the disk.
- **Queue** The queue itself is a durable FIFO — LevelDB is already a dependency, in its own directory rather than sharing `MANAGER_DATA_FOLDER`. Timestamps come from the host clock instead of a block, which matters only if guests read them; in-flight work adds `blockTimestamp` to the host ABI. The service must also publish the enclave's P521 communication key so clients can encrypt toward it, the equivalent of `getPubSecp521r1()` on `TeeAuthenticator`. 

**Notification.** On EVM a client learns the fate of its request from subgraph-indexed events: `RequestSubmitted` when it lands, `RequestCompleted` with status and error code, and `UserEvent` for anything addressed to it. With no chain there are no events and no subgraph — but the data is not missing, only unserved.

**The commit log already carries everything the events do.** A signed `UpdatePayload` contains `ApplicationID`, `RequestID`, `ErrorCode`, `ErrorMsg`, the encrypted `Events` with their subtypes, and the `AppEvents`. That is strictly more than the subgraph publishes, and it is enclave-signed rather than merely emitted. No-chain therefore needs a read API over the log, not a new data model.

**The client-side abstraction already exists and is already an interface.** `vela-common-go/subgraph/types.go:17` declares `Client` with `GetRequestCompletedByID`, `GetUserEvents`, `GetUserEventsBySubTypes`, `GetAppEvents` and the money queries. A no-chain implementation satisfies the same interface by reading the log instead of GraphQL, with `GetRefunds`/`GetWithdrawals`/`GetClaimsExecuted` returning empty because there is no custody. Existing clients, `vela-common-ts` included, keep working against one interface with two backends.

**The privacy mechanism carries over untouched.** Users locate their own events by deriving subtypes with `GenerateSubtype` — HMAC-SHA256(seed, index) over an anonymity set of 50 (`pkg/executor/subtype.go`). No chain is involved, so a client filters log entries by derived subtype exactly as it filters subgraph events today. No client-side crypto changes.

Serve it from the witness, which already holds the log and whose view is by construction the finalised one. Serving from the manager also works — entries are enclave-signed, so it can misreport by omission but never by forgery — and with a trusted manager that costs nothing.

Submission needs no equivalent of `RequestSubmitted`. It is a direct call into a host service that enqueues the request, so a call returning without error *is* the acknowledgement, and `GetRequestCompletedByID` finding nothing afterwards unambiguously means *in flight* rather than *never arrived*. One implementation requirement follows: the service must acknowledge only after a **durable** enqueue, or a crash between ack and write silently loses a request the client believes is pending.

### 2.4 Do the two modes still fit one core?

Yes. Both modes implement the same six-step protocol:

```
observe  → fetch → execute → sign → commit (may be rejected) → roll back on rejection
```

EVM mode is already **optimistic**, and that turns out to be the property that makes this work. The manager saves reports (`pkg/manager/manager.go:673`), stores state (`:683`), *then* submits on-chain (`:690`) and rolls back if the transaction reverts. Nothing gates execution on freshness; the chain adjudicates afterwards. A commit-time witness slots into exactly that position, so the shared code path — including the rollback branch that already exists in `submitBatchStateOnChain` — is unchanged.

| Step | EVM | No-chain | Same core code? |
|---|---|---|---|
| observe | `PendingWork()` over contract queues | `PendingWork()` over local queue | yes |
| fetch | `FetchFor(appID, n)` view | `FetchFor(appID, n)` local read | yes |
| execute | enclave, integrity check only | identical | **identical binary** |
| sign | P2 codec → `UpdatePayload` | same codec, same artifact | yes |
| commit | `stateUpdate()` tx; may revert `InvalidStateRoot` | witness `append`; may reject stale `prevStateRoot` | yes, behind one port |
| reject → rollback | `blockchain.ReorgError` → `dataLayer.Rollback` | witness rejection → same rollback | yes |

The one deliverable this implies is a **shared rejection taxonomy** in P3: today rejection is expressed as `blockchain.ReorgError` plus an `InvalidStateRoot` revert string. The core needs a mode-neutral vocabulary — roughly *stale-root*, *not-my-turn*, *transient* — that both adapters map onto. That is a small piece of design, but skipping it is how the modes drift apart.

**Four things legitimately differ, and all of them are adapter-local rather than core changes:**

| Concern | EVM | No-chain |
|---|---|---|
| Economics | fees, refunds, custody, withdrawals | none — `FuelPricePerUnit`/`MinFeePerRequest` = 0, and **a non-empty `Withdrawals` list must be rejected**, since there is no custody to pay from |
| Trigger / TRUSTPROCESS | full callback loop | absent |
| Deploy authorization | `DEPLOYAPP` on-chain role | local allowlist or config |
| "Was this confirmed?" (authority service) | subgraph read | witness log read |

The `UpdatePayload` schema stays shared, with the economics fields zeroed in no-chain mode. That is preferable to two schemas: one codec, one audit artifact, one signature format.

---

## 3. Other attention points

### 3.1 The message to sign

In chain mode, the TEE payload is signed in GO code and verified in Solidity: `pkg/executor/msgtosign_builder.go` (227 lines, the densest hand-written EVM file in the repo) ABI-encodes thirteen fields in exactly the order `AbstractTeeAuthenticator.checkSignature` reconstructs them, keccaks, and applies the EIP-191 `personal_sign` prefix. That signature *is* the trust link — re-encoding it outside the enclave would mean re-signing it. So the encoding must be produced by attested code.

In no-chain mode we can reuse the same signature scheme and replicate its verification in the witness; or define another scheme, in which case the Executor must generate the signature in the matching format.

### 3.2 Request Scheduling: core owns selection, contracts keep correctness

**Decided (assumption 3).** One scheduler, written once in Go, used by both modes. The contract keeps the invariants that are about correctness and drops the ones that are about policy.

| Property | Where | Why |
|---|---|---|
| Per-app FIFO; request is its queue's head; `prevStateRoot` chaining | **On-chain, enforced** | Correctness. Prevents reordering or skipping a specific user's request. Cheap — head comparisons. |
| Cross-app *selection* / fairness | **Core scheduler** | Policy. One implementation, shared by both modes. |
| Detection of unfair selection | **Off-chain, from existing events** | See below — this needs no new code at all. |

#### Detectability on EVM is already free

`RequestSubmitted(indexed uint64 applicationId, …)` and `RequestCompleted(indexed uint64 applicationId, …)` are both emitted and **both already indexed by the subgraph** (`subgraphs/hcce/subgraph.yaml:32-35`). Per-app pending depth over time is therefore `submitted − completed`, already queryable. An app that had pending work for N blocks and was never served is derivable from data you are collecting today.

**So no new audit event, no new contract code, and no new subgraph work is needed for detectability** — only a query and an alerting threshold. P4's contract side therefore shrinks rather than grows.

None of that exists in no-chain mode, and none of it needs to. The *served* series is free from the commit log — `applicationId` and `processedRequestId` are both in the signed field set (`msgtosign_builder.go:73,76`), so which application was served, and in what order, is signed rather than merely emitted. The demand side is not recorded, but with a trusted manager there is no unfair selection to detect; and with a single pre-deployed application cross-app fairness does not arise at all.

---

## 4. Effort

Assumes the decided scope: multiple EVM chains + no-chain mode, per-target enclave images, no identity redesign, no I`AppData` migration, trigger remains EVM-only.

**Basis.** `Est.` is conventional human-developer effort. `+ AI` assumes an engineer working with an AI coding assistant on this codebase: it compresses mechanical, compiler-verified work heavily and design, review and deployment ceremony barely. Treat `Est.` as the planning number and `+ AI` as the upside, because the phases that dominate risk here — P0's decisions, P2's byte-exact codec, the redeploy and PCR0 ceremony — are precisely the ones that do not compress.

| Phase | Work | Est. | + AI |
|---|---|---|---|
| **P0** | Decide: domain separation, scheduling split (§3.2), no-chain freshness posture (§2.1). Design + review. Blocks P2/P4. | 1–2 | 1–2 |
| **P1** | Relocate `Address`/amount types out of go-ethereum into core. **Format-compatible — no state migration.** Mechanical across ~10 non-test files, rippling into ~27 test files. | 1 | 0.5 |
| **P2** | Extract signing codec behind a narrow interface; byte-for-byte golden vectors against the existing contract tests. Adding a domain separator — a chain id or verifying-contract address in the signed tuple, absent today — is *optional scope* here: including it adds the Solidity `checkSignature` change and a redeploy path; excluding it saves roughly half a week now and costs more later. | 1 | 0.5 |
| **P3** | Split `pkg/manager` into core orchestrator (state reconciliation, executor round-trip, storage, reports) and chain driver (polling, reorg policy, tx submission). ~60/40. **Largest single item.** | 3 | 2 |
| **P4** | Core scheduler; per-app queues in the contract; retire the round-robin cursor and its enforcement scan; reshape the request-source port so the core observes pending work across apps and then fetches for the one it selects. | 2 | 1.5 |
| **P5** | No-chain adapter end-to-end: local request source, append-only log of signed payloads + attestation sink, null economics, no-op reorg. Includes the `AssociateKey` sender-signature check and its protocol-version bump (§2.2). | 2 | 1.5 |
| **P5b** | Freshness hardening (§2.1): one independent stateful commit-time witness — holds the per-app root, rejects a stale `prevStateRoot`, advances. No enclave changes. `k`-of-`n` shaped with `k=n=1`, so a quorum is later a config change. Recommended. | 1–2 | 0.5–1.5 |
| **P5c** | Client-facing surface for no-chain (§2.3): intake service (requestId assignment, idempotency, durable-before-ack enqueue, depth cap and rate limiting, key publication) over a durable FIFO; read API on the commit log with a subtype index and pagination; a `subgraph.Client` implementation over it, money queries returning empty. | 3–4 | 2–3 |
| **P6** | Build/CI: module boundaries, enclave image(s), dockerfiles, NOTICES. | 1 | 0.5 |
| **P7** | Test harness: core suite against the no-chain adapter; keep the simulated-backend suite for EVM. | 2 | 1 |
| | **Total** | **17–20** | **11–14** |

Two things make this cheaper than a split of this size usually is. The no-chain adapter (P5) is **both** the second implementation that validates the abstraction **and** a shipped product, so the validation is not a pure-cost phase. And there is no `AppData` migration, because 20-byte identity survives unchanged.

### 4.1 Staffing: 1, 2 or 3 people

Total effort is not fixed as headcount rises, and calendar time does not fall proportionally. Two structural facts drive this:

**The critical path sums to 12–14 weeks and cannot be compressed by hiring.** `P0 → P1 → P3 → P5 → P5c → P7` is a genuine chain: the design gates the codec and scheduler decisions; the type relocation must land before the large refactors; the no-chain adapter needs the driver interface that P3 creates; the client surface needs the adapter's queue and log to exist; the harness needs something to test. P5b runs alongside P5c rather than after it, so the witness is not what sets the length. Overlapping the tail (starting P5 once P3's interface is stable rather than once P3 is done, running the core harness alongside P5c) brings it to roughly **10–12 weeks**, and no amount of staffing goes below that. With AI assistance the same chain is 8–10 weeks, and roughly **7–9** once overlapped.

**Two phases resist parallelism specifically.** P1 is a *wide, shallow* sweep across `pkg/common`, `pkg/wasm`, `pkg/crypto` and `pkg/executor`, rippling into ~27 test files — it conflicts with essentially every other workstream, so it wants to run early and alone. P3 is a *narrow, deep* refactor of a single 1,465-line package plus 2,932 lines of tests; a second pair of hands on it saves less than it costs.

| Team | Total effort | **+ AI (weeks)** | **Calendar (with AI)** | Bound by |
|---|---|---|---|---|
| 1 | 17–20 pw | 11–14 | **11–14 weeks** | capacity |
| 2 | ~20–23 pw (≈15% coordination) | ~13–16 | **7–9 weeks** | purely critical path |
| 3 | ~23–27 pw (≈35% coordination) | ~15–19 | **7–9 weeks** | purely critical path |

The first two columns are effort; the third is elapsed time, which is why they diverge as the team grows. Calendar compresses far less than effort does, because the critical path is dominated by decisions, review and the redeploy ceremony rather than by typing — which sharpens the conclusion below: the more the mechanical work compresses, the less a third person buys.

**On these numbers the third person buys no speed at all** — 2 and 3 land in the same **7–9 week** band, because from two people onward the critical path is the only binding constraint. The third costs about 2–3 extra person-weeks and buys only *downside* protection: slippage in one workstream has somewhere else to absorb it. If schedule matters more than risk, two is the right team.
