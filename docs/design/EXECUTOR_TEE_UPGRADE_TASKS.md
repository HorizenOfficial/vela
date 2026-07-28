# Executor TEE Upgrade — Jira Tasks

Design doc: `docs/design/EXECUTOR_TEE_UPGRADE_DESIGN.md`

## Gaps found during code review (not covered by the design doc)

These emerged from checking the design against the current code; each is folded into a task below.

| # | Gap | Task |
|---|-----|------|
| G1 | No component knows the running enclave's PCR0. The Manager never sees an attestation, and the Executor never queries the NSM for its own PCRs — the design's `activeImage` vs "running image" comparison has no data source. | 3 |
| G2 | No reconnect path exists. The handshake is single-shot (`sync.Once`, `pkg/manager/manager.go:102-118`), `Connect` has no retry, and the communication client never re-dials after disconnect. "The Manager reconnects" (design step 5) is currently impossible without a Manager process restart. | 5 |
| G3 | R2 (recover, never regenerate) has no enforcement. Keyset regeneration is triggered purely by the Manager's LevelDB `IsNotFound` (`pkg/communication/client.go:363-377`), and `StoreEnclaveKeySetRecovery` is an unconditional overwrite (`versioned_leveldb_datalayer.go:268-282`). A wiped/wrong Manager data folder during the swap silently orphans all encrypted state. | 6 |
| G4 | No clean-exit signal to the Executor exists (no shutdown message type or admin verb) — required by the drain sequence. | 4 |
| G5 | Multi-step `updateTee` straddle: the PCR check runs only in `updateTeeStep1` (`TeeAuthenticator.sol:76`); a PCR0-set change between step 1 and step 4 is not re-validated. **RESOLVED in Task 1:** `updateTeeStep4` re-runs `_checkAttestationContent` before finalizing. | 1 |
| G6 | The design says an unapplied proposal is reverted by "letting it expire", but `PendingSwap` has no expiry and no cancel function — a stale past-`eta` proposal stays instantly applicable forever, eroding the R4 audit window. **RESOLVED in Task 1:** `cancelPcr0Swap` + proposals expire `PCR0_SWAP_APPLY_WINDOW` (7 days) after `eta`; a live proposal must be cancelled before a new one (`SwapAlreadyPending`), an expired one may be overwritten. | 1 |
| G7 | The UUPS conversion (`UPGRADABLE_CONTRACTS_DESIGN.md`) is not implemented — `TeeAuthenticator` is plain `Ownable` + constructor. **RESOLVED: UUPS is out of scope here.** The contract stays `Ownable` + constructor; `pcr0UpgradeDelay` and the initial `acceptedPcr0`/`activeImage` are seeded in the constructor. All UUPS-supporting changes (`initialize()`, `immutable`→storage, `__gap`) land with the UUPS conversion itself. | — |
| G8 | Stale design claim: the handshake wait already has a timeout (`HANDSHAKE_TIMEOUT`, default 5s, `pkg/manager/manager.go:79-100`). R7's remaining work is version negotiation only. | 0 |
| G9 | `proposePcr0Swap` doesn't validate the 48-byte PCR0 length; a wrong-length entry would join the set as an unmatchable `activeImage`. **RESOLVED in Task 1:** `InvalidPcr0Length` in both the constructor and `proposePcr0Swap`. | 1 |

## Additional findings from design discussion

Raised while stress-testing the design; each needs a Task 0 decision and lands in the mapped task.

| # | Finding | Task |
|---|---------|------|
| D1 | A TEE signature proves **key possession, not code identity**: the KMS key policy is an off-chain trust root that a malicious infra admin can diverge from the on-chain set (grant a rogue PCR0 → rogue image recovers the keys → signs indistinguishably and/or passively decrypts state — the latter is undetectable by any protocol). Mitigations are preventive: KMS policy governance/auditability (Task 8); optionally a periodic challenge attestation proving the `activeImage` enclave is live and holds the key (Task 1); longer-term, enclave-to-enclave key handoff gated on the on-chain set, demoting KMS to break-glass. | 0, 1, 8 |
| D2 | After a swap, `teeSigner` on-chain was attested under the **old** image only; the new image's key possession is guaranteed only transitively via KMS. Cheap runtime guard: on reconnect, the Manager compares the handshake-reported `SigningKeyAddr` against on-chain `getTeeSigner` — turns silent key regeneration into an immediate typed error instead of a first-state-update revert. | 5 |
| D3 | `_checkAttestationContent` checks set membership, which is looser than needed: `updateTee` is a manual, quiescent-time operation, so requiring `PCR0 == activeImage` (or membership AND active) would be strictly tighter with no legitimate flow broken. **RESOLVED in Task 1: tightened** — the attestation PCR0 must equal `activeImage`; the accepted set only governs the swap/rollback lifecycle. | 0, 1 |
| D4 | No urgent-fix path. Rollback-class emergencies are already instant, but a vulnerability present in **all** accepted images forces downtime equal to the full timelock, with no way to stop request intake meanwhile — and users cannot even withdraw during a stop (withdrawals need TEE signatures). Add an intake **pause** (guardian-triggerable, instant — pausing only restricts, same safe-direction logic that exempts `removePcr0` from the timelock); optionally a break-glass apply gated by a stronger quorum; keep `pcr0UpgradeDelay` modest (24–48h) since it doubles as worst-case emergency downtime. Emergency runbook has sequencing traps: `removePcr0` reverts on the active image (swap away first), and stripping a PCR0 from the KMS policy does **not** evict a live enclave holding keys in memory (must terminate the enclave too). | 0, 1, 8 |

## Implementation Order

```
┌────────────────────────┐
│ 0. Design doc updates  │  (decisions unblock everything)
└───────────┬────────────┘
            │
┌───────────▼────────────┐
│ 1. Contract: PCR0 set  │──────────────┐
│    + swap flow         │              │
└───────────┬────────────┘              │
            │                           ▼
            │                 ┌───────────────────┐
┌───────────▼────────────┐    │ 3. Executor PCR0  │
│ 2. Go bindings +       │    │    self-identity  │
│    blockchain client   │    └────────┬──────────┘
└───────────┬────────────┘             │
            │        ┌─────────────────┤
            │        │                 │
            │  ┌─────▼──────────┐      │
            │  │ 4. Protocol    │      │
            │  │    version +   │      │
            │  │    shutdown    │      │
            │  └─────┬──────────┘      │
            │        │                 │
            └────────┼─────────────────┘
                     ▼
          ┌──────────────────────┐   ┌──────────────────────┐
          │ 5. Manager: observe, │   │ 6. Key-continuity    │
          │    drain, reconnect  │   │    guards (R2)       │
          └──────────┬───────────┘   └──────────┬───────────┘
                     └────────┬─────────────────┘
                              ▼
                   ┌──────────────────────┐
                   │ 9. Integration tests │
                   └──────────────────────┘

Independent: 7. Reproducible EIF build   8. KMS policy runbook
```

- Task **0** first: it settles the open decisions (G1, G2, G6) the other tasks implement.
- Tasks **3**, **4**, **6**, **7**, **8** can start in parallel with Task **1**.
- Task **2** depends on **1** (compiled ABI). Task **5** depends on **2**, **3**, **4**.
- Task **9** depends on **5** and **6**.

---

## Task 0 — Design Doc Updates & Decisions — **RESOLVED**

**Summary:** Resolve the gaps found in review and correct stale claims before implementation starts.

**Description:**

Update `EXECUTOR_TEE_UPGRADE_DESIGN.md`:

- Fix the stale "handshake is blocking and has no timeout" claim (G8): the timeout exists (`HANDSHAKE_TIMEOUT`); rescope R7 to version negotiation only.
- Decide and document the source of the running image identity (G1). Recommended: Executor queries NSM `DescribePCR(0)` at startup and reports PCR0 in the handshake.
- Decide reconnect strategy (G2): in-process reconnect/re-handshake loop in the Manager vs. documented operator-driven Manager restart after enclave relaunch. The current text implies the former, which does not exist.
- Specify the graceful-shutdown verb (G4) — the drain sequence needs it even though the design rules out an *emergency* stop.
- Add proposal cancel/expiry semantics (G6) and PCR0 length validation (G9) to the contract spec.
- G7 resolved: UUPS is out of scope — `TeeAuthenticator` stays `Ownable` + constructor, with `pcr0UpgradeDelay` and the initial `acceptedPcr0`/`activeImage` seeded in the constructor. UUPS-supporting changes (`initialize()`, `immutable`→storage, `__gap`) are deferred to the UUPS conversion itself.
- Decide the `updateTee` PCR check: set membership vs. `== activeImage` (D3). Document the chosen trust model either way. **DECISION: `== activeImage`** (tightened) — implemented in Task 1 `_checkAttestationContent`; the accepted set governs only the swap/rollback lifecycle. Documented in `EXECUTOR_TEE_UPGRADE_DESIGN.md` § Attestation check against the active image.
- Decide on a periodic challenge attestation (D1): on-chain freshness requirement (`maxAttestationAge` on the registered attestation, verified against `activeImage`, nonce = recent block hash in `userData`, submitted by the Manager on the polling loop). Document its limit honestly: it proves the right image is *live and holds the key*; it cannot exclude a parallel rogue image admitted via a diverged KMS policy. **DECISION: not adopted** — it does not close the core D1 risk (a diverged-KMS-policy rogue image also holds the key and answers challenges; passive decryption is undetectable). D1 is mitigated preventively via KMS-policy governance/auditability (Task 8 runbook) and the D2 signer-continuity guard (Task 5). Recorded in `EXECUTOR_TEE_UPGRADE_DESIGN.md` § Implementation Notes → *Challenge attestation (D1)*.
- Decide the emergency story (D4): intake pause (recommended), optional stronger-quorum break-glass, recommended `pcr0UpgradeDelay` range. Document that the timelock is also the worst-case downtime when every accepted image is vulnerable, and the withdrawal-freeze tension (users cannot withdraw while the enclave is stopped). **DECISION: intake pause not implemented** — the emergency levers are the on-chain timelock (= worst-case downtime, keep modest 24–48h) and an out-of-band `nitro-cli terminate-enclave`. Recorded in the Task 8 runbook (`docs/ops/KMS_KEY_POLICY_RUNBOOK.md` § 7).
- Decide the enclave base image (Task 7): Amazon Linux 2 reached end of support on 2026-06-30; migrating to AL2023 vs. pinning EOL AL2 RPMs (PCR0 changes either way). **DECISION: migrate to AL2023** — implemented in Task 7 (`dockerfiles/executor/`, `REPRODUCIBLE_EIF_BUILD.md`).
- Update `EXEC_MGR_HANDSHAKE.md` for the new handshake fields once Tasks 3/4 land.

**Files:**
- `docs/design/EXECUTOR_TEE_UPGRADE_DESIGN.md`
- `docs/design/EXEC_MGR_HANDSHAKE.md`

---

## Task 1 — Smart Contract: PCR0 Set + Swap Flow — **IMPLEMENTED**

**Summary:** Replace single `pcr0` in `TeeAuthenticator.sol` with an accepted set, `activeImage` pointer, and timelocked swap flow.

**Description (as implemented):**

- Replaced `bytes public pcr0` with `mapping(bytes32 => bool) acceptedPcr0`, `bytes32[] acceptedPcr0List`, `bytes32 activeImage`, `PendingSwap pendingSwap`, `uint256 immutable pcr0UpgradeDelay` (constructor-set, no setter; becomes an `initialize()`-set storage var at the UUPS conversion).
- Implemented `proposePcr0Swap` / `applyPcr0Swap` / `cancelPcr0Swap` / `removePcr0` (+ `getAcceptedPcr0Count`), all `onlyOwner`, with events `Pcr0SwapProposed` / `Pcr0Swapped` / `Pcr0SwapCancelled` / `Pcr0Removed`. Removed `updatePcr0` / `PcrZeroUpdate`.
- Raw PCR0 preimages emitted **non-indexed** in the new events (`acceptedPcr0List` stores only hashes; auditors need preimages — the old `PcrZeroUpdate` had `indexed bytes`, which only exposes the keccak). The constructor also emits `Pcr0Swapped` so the seed preimage is on-chain.
- `proposePcr0Swap`: `InvalidPcr0Length` unless 48 bytes (G9); reverts `SwapAlreadyPending` while a live proposal exists (explicit `cancelPcr0Swap` required — clean event lifecycle for auditors), an expired one may be overwritten. `applyPcr0Swap` must run within `PCR0_SWAP_APPLY_WINDOW` (7 days) after `eta`, else `SwapProposalExpired` (G6).
- `_checkAttestationContent`: fixed 48-byte extraction at offset 4, compared against `activeImage` (D3 resolved: tightened — set membership is not sufficient).
- Straddle fix (G5): `updateTeeStep4` re-runs `_checkAttestationContent` before finalizing.
- Constructor seeds the initial PCR0 into the set and `activeImage` (G7 resolved: UUPS out of scope; `initialize()`, `immutable`→storage and `__gap` land with the UUPS conversion).
- `removePcr0` has no `CannotRemoveLastPcr0` guard: `activeImage` is always a member of the set (constructor and `applyPcr0Swap` insert before pointing at it), so `CannotRemoveActiveImage` alone keeps the set non-empty.
- Deferred to Task 0 decisions: challenge-attestation freshness (D1); intake pause (D4 — likely lands on `ProcessorEndpoint` request submission rather than here; instant, guardian role, no timelock since pausing only restricts).
- Updated `ITeeAuthenticatorAdmin.sol`, `contracts/scripts/deploy/teeAuthenticator.ts` and `scripts/deploy/all.ts` (new `TEE_PCR0_UPGRADE_DELAY` env var); `scripts/management/updateTee.ts` unchanged (step flow only). Added `mocks/MockNitroProver.sol` so the attestation-content checks run on the Hardhat network (the real prover needs Anvil and a genuinely signed attestation).

**Tests (in `contracts/test/TeeAuthenticator/`):**
- Propose → apply before `eta` reverts (`TimelockNotElapsed`); after `eta` succeeds, adds to set, sets `activeImage`, emits `Pcr0Swapped`
- Swap to already-accepted PCR0 applies immediately (rollback path), no duplicate list entry
- `removePcr0`: active image reverts (also covers the last entry, by the invariant), non-member reverts, swap-remove keeps list consistent
- Wrong-length PCR0 proposal reverts; second proposal while one is live reverts; re-propose allowed after cancel/expiry
- Cancel / expiry of a pending proposal; removed PCR0 needs a fresh timelocked proposal
- Straddle: `activeImage` change between step1 and step4 → step4 reverts
- Attestation accepted only for `activeImage` (incl. after rollback); rejected for non-active members, non-members, truncated PCRs
- Non-owner calls revert on all four functions
- Constructor seeding of `acceptedPcr0` / `activeImage` / `pcr0UpgradeDelay`

**Files:**
- `contracts/contracts/TeeAuthenticator.sol`
- `contracts/contracts/interfaces/ITeeAuthenticatorAdmin.sol`
- `contracts/contracts/mocks/MockNitroProver.sol` (new)
- `contracts/scripts/deploy/teeAuthenticator.ts`, `contracts/scripts/deploy/all.ts`
- `contracts/test/TeeAuthenticator/` (new), `contracts/test/AttestationTeeAuthenticator.ts`

---

## Task 2 — Go Bindings + Blockchain Client: `activeImage` Read — **IMPLEMENTED**

**Summary:** Regenerate the `tee` bindings and expose `GetActiveImage` on the blockchain client.

**Description (as implemented):**

- The `tee` binding (`pkg/blockchain/contracts/tee/TeeAuthenticator.go`) was already regenerated in Task 1 and carries `activeImage()` (`PackActiveImage` / `UnpackActiveImage`).
- Added `GetActiveImage(ctx) ([32]byte, error)` to `blockchain.Client` (`pkg/blockchain/interface.go`) and implemented it on `BlockChainClient` next to `GetTeePublicKey` (reads via `teeAuthEndpoint.PackActiveImage` / `UnpackActiveImage`). Added it to `MockClient` with a `SetActiveImage` setter and a `GetActiveImage` mock-func override. `eventBroadcastingClient` embeds the interface, so it delegates for free.
- Test-contract parity so the read works on the simulated backend in both modes: added `bytes32 public activeImage` + a permissionless `setActiveImage(bytes32)` to **both** `MockTeeAuthenticator.sol` (client/fullstack sim, `useMockContracts=true`) and `NoAttestationTeeAuthenticator.sol` (non-mock sim + dev deploy). No constructor change → no ripple to `fixture.ts`, `all.ts`, `setupContracts`, or `NewSimTestHelper`. Regenerated the `mocktee` and `noattestationtee` bindings (a full `go generate ./pkg/blockchain/...` changes only these two — no drift).
- Testutil: `SimTeeAuthenticatorHelper.GetActiveImage()` reads the value; `SimTestHelper.SetActiveImage()` sets it (the `setActiveImage(bytes32)` selector is identical across both variants, so the `mocktee` binding packs the call regardless of which is deployed). These also serve Task 5's drain/rollback tests.

**Tests:**
- `TestGetActiveImage` (`pkg/blockchain/client_test.go`): default zero image, then set on-chain and read back through the client and the tee helper.

**Files:**
- `pkg/blockchain/client.go`, `pkg/blockchain/interface.go`, `pkg/blockchain/mock_client.go`, `pkg/blockchain/client_test.go`
- `pkg/blockchain/testutil/sim_tee_auth_helper.go`, `pkg/blockchain/testutil/sim_test_helper.go`
- `contracts/contracts/mocks/MockTeeAuthenticator.sol`, `contracts/contracts/mocks/NoAttestationTeeAuthenticator.sol`
- `pkg/blockchain/contracts/mocktee/MockTeeAuthenticator.go`, `pkg/blockchain/contracts/noattestationtee/NoAttestationTeeAuthenticator.go` (generated)

---

## Task 3 — Executor: PCR0 Self-Identity (G1) — **IMPLEMENTED**

**Summary:** Executor reads its own PCR0 from the NSM and reports it in the handshake.

**Description (as implemented):**

- Added `nsmutil.DescribePCRWithSession(opener, index)` (`pkg/nsmutil/describe_pcr.go`), mirroring `AttestWithSession`; issues `request.DescribePCR{Index}` and returns the raw PCR bytes.
- Executor reads PCR0 once at construction (`readSelfPCR0`, `pkg/executor/executor.go`) via `nsm.OpenDefaultSession`, caches it hex-encoded on `StatelessExecutor.pcr0`. On any error (no `/dev/nsm` in TCP/dev mode) it logs and stores `""` — the **dev marker**. No gating on `ChannelType`; the graceful fallback handles it.
- Added `Pcr0` (`omitempty`) to the two executor→manager handshake-result payloads (`SetKeysetRecoveryRequestData`, `KeysetRecoveryResultData`) and threaded it as a positional param through `ServerConnection` / `ClientRequestHandler` (`SetKeysetRecovery`, `KeysetRecoveryResult`, `HandleSetKeysetRecoveryRequest`, `HandleKeysetRecoveryResult`) — consistent with the existing `commPubKey`/`signingKeyAddr` style. Executor passes `e.pcr0`. (An executor binary-version field was initially threaded here too but later removed as dead — the version is obtained on demand via the forwarded `get_version` admin command.)
- Manager stores the reported PCR0 as the running image (`runningPcr0` + `RunningPcr0()` getter), for Task 5's `keccak256(pcr0)` vs `activeImage` comparison. **Guarded by a dedicated `identityMu`, not `m.mu`**: the handshake handlers run on the comm client's reader goroutine while `Start` holds `m.mu` through `waitForExecutorHandshake`, so reusing `m.mu` deadlocked the handshake (caught by the system tests).
- Admin `get_version` (target=all) now includes `executorPcr0` in `AggregatedGetVersionResponse`, sourced from the manager's stored handshake PCR0 — no change to the executor's `get_version` wire format or the forwarding path.
- No `go.mod` change (`DescribePCR` uses the already-vendored `hf/nsm`), so no `NOTICES` update.

**Tests:**
- `pkg/nsmutil/describe_pcr_test.go`: success, opener error, device error, empty response (mock session).
- `pkg/communication/handshake_test.go`: `TestHandshake_FirstConnection` / `TestHandshake_Reconnection` assert PCR0 propagates end-to-end over TCP through both result messages.
- `pkg/manager/manager_test.go`: `TestHandleKeysetRecovery_StoresRunningPcr0` asserts both handshake paths store `RunningPcr0`.
- Full system suite (`tests/system`, TCP dev mode with the dev marker) passes; NSM read on Nitro is covered by Task 9.

**Files:**
- `pkg/nsmutil/describe_pcr.go` (new), `pkg/executor/executor.go`
- `pkg/communication/message.go`, `interface.go`, `server.go`, `client.go`
- `pkg/manager/manager.go`, `pkg/admin/admin_interface.go`
- `docs/design/EXEC_MGR_HANDSHAKE.md`, `pkg/admin/README.md`

---

## Task 4 — Communication: Protocol Version Negotiation (R7) + Graceful Shutdown (G4) — **IMPLEMENTED**

**Summary:** Add `WireProtocolVersion` exchange with typed incompatibility error, and a clean-shutdown message.

**Description (as implemented):**

- Added `communication.WireProtocolVersion uint32 = 1` (monotonic; distinct from `pkg/version.Version` and `common.Request.ProtocolVersion`) plus `IsCompatible(peer)` (exact-match policy) and a typed `IncompatibleProtocolError{Local, Peer}` (`message.go`).
- Version is exchanged on the **first** request/response pair so an incompatible peer is rejected before any keyset-recovery data is read/generated/stored:
  - Executor sends its version in `GetKeysetRecoveryRequestData.WireProtocolVersion` (json `wireProtocolVersion`, `server.go`).
  - Manager checks it in `HandleGetKeysetRecoveryRequest(ctx, peerProtocolVersion)`: on mismatch it fails the handshake with `IncompatibleProtocolError` (via `completeExecutorHandshake`) and returns no data; otherwise replies with `GetKeysetRecoveryResponseData.WireProtocolVersion` (`client.go`, `manager.go`).
  - Executor checks the Manager's version in `GetKeysetRecovery` and aborts before restore/generate on mismatch (`server.go`). Both-side checks are required so an up-to-date peer rejects a legacy peer (absent field = version 0), since the legacy peer cannot check.
- Graceful shutdown implemented as the forwarded admin verb `AdminCmdShutdown` (`"shutdown"`), **not** an operator-facing verb (not routed by `ExecuteCommand`). Executor's `HandleAdminCommand` acks with `ShutdownResponse{Stopping:true}` then closes `ShutdownRequested()`; `cmd/executor/main.go` selects on it alongside OS signals for a clean teardown. Manager helper `forwardShutdown` swallows the post-ack disconnect (drain must not treat it as fatal) — used by Task 5.
- `HandleGetKeysetRecoveryRequest` gained a `peerProtocolVersion uint32` parameter (interface + manager impl + test mocks).

**Tests:**
- `pkg/communication/protocol_version_test.go`: `IsCompatible` matrix (incl. version 0), absent-field-unmarshals-to-0, error message.
- `pkg/communication/handshake_test.go`: compatible handshake exchanges the version end-to-end (captured `gotProtocolVersion == WireProtocolVersion`).
- `pkg/manager/manager_test.go`: incompatible peer → typed error + handshake failed + no recovery read; compatible peer proceeds; `forwardShutdown` tolerates disconnect and ack.
- `pkg/executor/shutdown_test.go`: `shutdown` command acks and closes `ShutdownRequested()`, idempotent.

**Files:**
- `pkg/communication/message.go`, `interface.go`, `client.go`, `server.go`
- `pkg/executor/executor.go`, `cmd/executor/main.go`, `pkg/manager/manager.go`, `pkg/admin/admin_interface.go`
- `docs/design/EXEC_MGR_HANDSHAKE.md`, `pkg/admin/README.md`

---

## Task 5 — Manager: Swap Observation, Drain, Reconnect (R8, G2) — **IMPLEMENTED**

**Summary:** Poll `activeImage`, drain on mismatch, signal Executor shutdown, and support reconnect + re-handshake.

**Description (as implemented):**

- D2 prerequisite: added `GetTeeSigner(ctx) (ethCommon.Address, error)` to `blockchain.Client` + `BlockChainClient` + `MockClient` (mirrors `GetActiveImage`).
- Manager stores the handshake-reported signing address (`reportedSigner`) alongside PCR0; a fresh handshake re-arms the signer check. The identity/upgrade state (`runningPcr0`, `reportedSigner`, `draining`, `signerVerified`) is guarded by `identityMu`, **not** `m.mu` — the handshake handlers that write it run on the comm reader goroutine while `Start` holds `m.mu` (lock order is always `m.mu → identityMu`).
- `maintainExecutorForActiveImage` runs each tick in `processRequestFromChain` **before** `m.mu` (so the blocking reconnect handshake isn't held under `m.mu`): reads `activeImage`, compares `keccak256(reported pcr0)` against it.
  - Match → verify signer continuity, then dispatch. Dev/TCP executor (empty PCR0) → swap observation skipped.
  - Mismatch → `drainExecutor` (send `shutdown` via `forwardShutdown`, close the channel; idempotent within an episode), skip dispatch. Next ticks `reconnectExecutor` (re-arm single-shot handshake, re-dial, re-handshake); on success drain clears and the new image is re-evaluated → dispatch resumes when it matches. Handles rollback (activeImage re-pointed) identically and is level-based (Manager restart converges).
  - Dropped connection (no swap) → the same `reconnectExecutor` path also fires when `!executorClient.IsConnected()` even without a drain — an executor crash/restart on the *same* image, where `draining` is never set. `ExecutorClient` gained `IsConnected()` (the comm `Client` already tracks `connected` under `connLock`), mirroring the existing `blockchainClient.IsConnected()`→`Connect()` pattern in `processRequestFromChain`. Without this trigger the Manager would retry into a dead channel forever and the relaunched executor, waiting on a handshake that never comes, could never recover its keyset. A failed re-dial skips the tick and retries.
- Reconnect rework in `communication/client.go` (G2): `Connect` now creates a fresh `shutdown` channel each call (was closed permanently on `Close`) and passes it to the reader + `cleanupLoop`; `Close` split into `closeLocked`; a reader loop's exit tears down only if it still owns the current connection (`handleReaderExit`), so a stale reader can't kill a reconnected channel.
- Drain-window safety: transient `GetActiveImage`/reconnect errors and drain state return `skip` (nil error) — never `fatalErrChan`. Only the signer-continuity mismatch returns a fatal `*SignerContinuityError`.
- Signer-continuity (D2): skipped while on-chain `teeSigner == address(0)` (pre-`updateTee` bootstrap — a naïve equality check would otherwise kill the executor before the first attestation); a mismatch on a **set** `teeSigner` is fatal. Reuses `BlockchainPollingInterval`; no new config.
- `Stop` lock ordering: because `maintainExecutorForActiveImage` does its blocking work (blockchain reconnect, reconnect handshake) lock-free and takes `m.mu` only at dispatch, `Stop` must not hold `m.mu` across `m.wg.Wait()` — a tick already in that pre-lock section would then block forever at `m.mu.Lock()` while `Stop` blocked forever in `wg.Wait()`. `Stop` sets `isRunning=false` and closes `stopChan` under `m.mu`, **releases `m.mu`**, then `wg.Wait()`s and tears down: a mid-flight tick acquires `m.mu`, sees `!isRunning`, and exits without dispatching. A failed teardown `Close` restores `isRunning=true` so the stop can be retried.

**Tests:**
- `pkg/manager/manager_test.go`: match→dispatch, mismatch→drain (shutdown sent + channel closed), dev-marker→proceed, connection-loss→reconnect (no swap, `!IsConnected()`), signer continuity (bootstrap skip / match / mismatch fatal), full drain→reconnect→resume convergence (also covers rollback + level-based restart), and `Stop` not deadlocking against a tick parked in the pre-lock section.
- `pkg/communication/tcp_test.go`: `TestClient_ReconnectAfterClose` (Connect works after Close, request succeeds on the new connection).
- `pkg/blockchain/client_test.go`: `TestGetTeeSigner` against the simulated backend.
- Full system suite + `-race` on communication/manager pass.

**Files:**
- `pkg/manager/manager.go`, `pkg/communication/client.go`, `pkg/communication/interface.go`
- `pkg/blockchain/interface.go`, `client.go`, `mock_client.go`
- `docs/design/EXEC_MGR_HANDSHAKE.md`

---

## Task 6 — Key-Continuity Guards (R2, G3) — **IMPLEMENTED**

**Summary:** Prevent silent keyset regeneration/overwrite when the Manager's recovery data is unexpectedly absent.

**Description (as implemented):**

Regeneration was triggered solely by the Manager's LevelDB `IsNotFound` result, and `StoreEnclaveKeySetRecovery` unconditionally overwrote. A Manager started against a wiped/wrong data folder during an upgrade would cause the new enclave to generate a fresh keyset, overwrite the recovery blob, and permanently orphan all encrypted state (and change `teeSigner`). Two independent guards close this:

- **Storage exists-guard** (`versioned_leveldb_datalayer.go`): `StoreEnclaveKeySetRecovery` now reads the current blob first and refuses to overwrite it with a *different* one, returning a typed `storageErrors.ErrRecoveryDataExists` (code `recovery_data_exists`, `IsRecoveryDataExists` helper). Re-storing an identical blob is a no-op (idempotent). The guard lives at the storage layer so it protects every caller; the Manager's `HandleSetKeysetRecoveryRequest` surfaces the error through `completeExecutorHandshake` and fails the handshake loudly. In a correct flow the store path is only reached on `found=false`, so the guard is defense-in-depth against a partially-populated folder.
- **Executor expect-existing guard** (`executor.go`, `config.go`): new `EXECUTOR_EXPECT_EXISTING_KEYSET` config (default `false`, `Config.ExpectExistingKeyset`). When set, a `found=false` handshake response is a fatal typed error (`ErrUnexpectedKeysetGeneration`) — `performHandshake` returns before calling `GenerateEnclaveKeySet` or `SetKeysetRecovery`, so nothing is generated or stored. Set during upgrades; a genuine first install runs without it.
- **Loud logging**: the generation branch now logs at `Warn` (was `Info`) spelling out that generating a keyset must only happen on a genuine first install and otherwise orphans all encrypted state.
- **Nitro delivery** (`dockerfiles/executor/Dockerfile`, `build-eif.sh`): an enclave receives no environment at launch, so the flag is baked into the enclave image — `ARG EXPECT_EXISTING_KEYSET=true` → `ENV EXECUTOR_EXPECT_EXISTING_KEYSET`. nitro-cli folds ENV into the measured ramdisk, so the guard's state is part of PCR0 and attestable (and third-party verifiable via the Task 7 reproducible build). Default ON for every release image; a **genesis** (first-install) build passes `EXPECT_EXISTING_KEYSET=false` (env knob on `build-eif.sh`, recorded in `build-info.json` so verifiers rebuild the right variant). The genesis image bootstraps the keyset once and must then be retired from the accepted PCR0 set and KMS key policy (runbook: Task 8). Dev/TCP deployments keep setting the env var directly.

**Tests:**
- `pkg/storage/versioned_leveldb/versioned_leveldb_datalayer_test.go`: overwrite with a different blob → `RecoveryDataExists` and original preserved; re-store of identical blob → idempotent success.
- `pkg/executor/keyset_continuity_test.go`: `found=false` + expect-existing → `ErrUnexpectedKeysetGeneration`, no generation, `SetKeysetRecovery` not called, keyset unset; first-install path (`found=false`, guard off) still generates and stores; `found=true` still restores with the guard on.
- `pkg/executor/config_test.go`: `EXECUTOR_EXPECT_EXISTING_KEYSET` defaults to `false` and parses `true`.

**Files:**
- `pkg/executor/executor.go`, `pkg/executor/config.go`
- `pkg/storage/versioned_leveldb/versioned_leveldb_datalayer.go`, `pkg/storage/errors/errors.go`
- `dockerfiles/executor/Dockerfile`, `dockerfiles/executor/build-eif.sh` (Nitro delivery: baked-in ENV, genesis build variant)

---

## Task 7 — Reproducible EIF Build (R5) — **IMPLEMENTED**

**Summary:** Make PCR0/PCR1/PCR2 a pure function of a git tag. PCR0 is a SHA-384 over the whole EIF, built from three inputs — each layer must be deterministic.

**Description:**

**A. Go binary** (`dockerfiles/executor/Dockerfile` build stage):
- Pin the builder base by digest (`amazonlinux:2@sha256:…` — moving tag today; with `CGO_ENABLED=1` the binary links that image's glibc, so a base update changes the binary).
- Pin exact RPM versions for the builder toolchain too (`gcc`, `binutils`, `git`, …): the base digest does not freeze what `yum install` fetches at build time, and with cgo the gcc/binutils version affects the linked binary.
- Verify the Go tarball sha256 after download.
- Add `-trimpath` and `-buildvcs=false` to `go build` (version already comes from ldflags; the VCS stamp only adds path/dirty-flag nondeterminism).
- No `--dirty` in release builds: build from a pristine checkout (`git archive <tag> | docker build -` — also excludes untracked files that `COPY . .` would pull in) and pass `GIT_VERSION=<tag>` explicitly. A `git archive` tree has no `.git`, so the Dockerfile's `git describe` fallback would silently stamp `"dev"` — make a missing `GIT_VERSION` a hard failure in release builds and update the "do not exclude `.git/`" Dockerfile comment.

**B. Runtime image filesystem** (becomes the application ramdisk — every byte and mtime lands in PCR0):
- Pin the runtime base by digest; pin exact RPM versions (`aws-nitro-enclaves-nsm`, `ca-certificates`, `shadow-utils` are unpinned today).
- Remove yum side effects in the same `RUN` layer (`/var/log/yum.log`, `/var/cache`, RPM DB timestamps).
- `useradd -m appuser` stamps `/etc/shadow`'s last-change field with the build **day** — PCR0 differs across build days even with everything else pinned. Normalize `/etc/shadow` (or drop `useradd` for a static passwd/shadow entry).
- Consider a minimal final image instead of a yum-managed one (also sidesteps the yum and `useradd` issues above). Note the executor is a cgo binary (wasmtime-go), dynamically linked against glibc — a minimal image must include the glibc runtime (`ld-linux`, `libc.so.6`, `libm`, pthread) copied from the pinned builder, not just the binary + CA bundle + passwd entry. (The Go NSM path itself is ioctl-based and adds no C dependency.)
- Normalize file timestamps: `SOURCE_DATE_EPOCH=$(git log -1 --pretty=%ct <tag>)` + BuildKit `--provenance=false --output type=docker,rewrite-timestamp=true` (or `kaniko --reproducible`). `rewrite-timestamp` needs BuildKit ≥ v0.13 — pin the buildx/BuildKit version alongside nitro-cli.

**C. EIF packaging** (the step currently outside the repo — most important gap):
- Add a versioned `dockerfiles/executor/build-eif.sh` that runs `nitro-cli build-enclave` inside a container pinned by digest with an exact `aws-nitro-enclaves-cli` RPM version. nitro-cli bundles its own kernel/init/bootstrap blobs (`/usr/share/nitro_enclaves/blobs/`) into the EIF — same Docker image + different nitro-cli = different PCR0/PCR1.
- `nitro-cli build-enclave` reads the image from the local Docker daemon, so the built image must be loaded into it and the containerized nitro-cli needs the docker socket mounted — the script handles both.
- The script emits the PCR JSON from `build-enclave`, the nitro-cli version, and blob checksums as a release artifact.

**Verification / release:**
- CI on release tags: build the EIF twice from scratch (independent runners, no shared cache), assert identical PCRs via `nitro-cli describe-eif`. Cheaper per-PR variant: single build compared against committed expected measurements.
- Publish per release: git tag, PCR0/1/2, nitro-cli version, base-image digests. Third-party verification during the timelock window reduces to `git checkout <tag> && ./dockerfiles/executor/build-eif.sh` → compare PCR0 with the `proposePcr0Swap` value. (Not the `.eif` file hash: nitro-cli stamps a non-measured `BuildTime` into the EIF header, so the file hash varies build-to-build while PCR0/1/2 stay identical — verified empirically.)
- Known limit: reproducibility depends on pinned RPMs staying fetchable. On AL2023 the versioned repositories keep a pinned `releasever` snapshot available long-term; if AWS ever retires a snapshot, vendor/mirror the RPMs with the release artifacts. The honest guarantee is at least "reproducible for the duration of the timelock window", which is what R5 needs.
- **Base image decision (Task 0): RESOLVED — migrate to Amazon Linux 2023.** AL2 reached end of support on 2026-06-30; since this task changes PCR0 anyway, both stages target AL2023. AL2023 also has no `amazon-linux-extras` (nitro-cli is in the default repos) and pins the whole package set via `--releasever` instead of per-package EVRs. **Implemented** in `dockerfiles/executor/` (`Dockerfile`, `nitro-cli.Dockerfile`, `versions.env`, `build-eif.sh`), `.github/workflows/reproducible-eif.yml`, and `docs/design/REPRODUCIBLE_EIF_BUILD.md`.

**Files:**
- `dockerfiles/executor/Dockerfile`, `dockerfiles/executor/build-eif.sh` (new), CI workflow, release docs

---

## Task 8 — KMS Key Policy Runbook + Governance (Ops) — **IMPLEMENTED**

**Summary:** Documented procedure for adding/removing PCR0 values in the `kms:RecipientAttestation:PCR0` condition, governance controls on the policy itself, and the emergency runbook.

**Implemented in** `docs/ops/KMS_KEY_POLICY_RUNBOOK.md` (new), cross-referenced from
`EXECUTOR_TEE_UPGRADE_DESIGN.md` (§ Off-Chain Changes → KMS key policy) and `CLAUDE.md`.
It covers all points below: the key-policy structure and AWS CLI read/write; the two-trust-root
sync invariant (on-chain `acceptedPcr0List` ↔ KMS condition) and the D1 divergence risk; the
governance controls; the safe upgrade sequence (KMS-first, then propose/timelock/apply/relaunch,
finalize with `removePcr0` + policy strip) with the `CannotRemoveActiveImage` ordering trap; the
instant rollback path; genesis-image retirement; and the D4 emergency runbook (rollback-and-retire,
full-stop-and-wait, the KMS-strip-does-not-evict caveat + `nitro-cli terminate-enclave`, and the
disclosure dynamic). The intake pause is documented as **considered but not implemented** (no
`pause` in `ProcessorEndpoint`), so the emergency levers are the timelock and an out-of-band
enclave terminate. No code change (`pkg/executor/kms` untouched).

**Description (requirements — all covered by the runbook):**

- **Upgrade runbook:** add `PCR0_new` (keeping `PCR0_old`) **before** the new enclave starts; remove `PCR0_old` only at finalization (mirror of `removePcr0`). Include the full safe-upgrade sequence checklist from the design and the rollback procedure.
- **Image variants (Task 6):** upgrade images bake `EXECUTOR_EXPECT_EXISTING_KEYSET=true`; only a genesis bootstrap uses the `EXPECT_EXISTING_KEYSET=false` build of `build-eif.sh`, exactly once per environment. After bootstrap, retire the genesis image: swap `activeImage` to the first guarded image, `removePcr0` the genesis PCR0, strip it from the KMS key policy. A genesis image left accepted is a standing bypass of the R2 guard.
- **KMS policy governance (D1):** the key policy is the off-chain twin of the on-chain accepted set and the single place a malicious admin can diverge them. Controls: restrict `kms:PutKeyPolicy` to a dedicated multi-party role separate from ops; grant `GetKeyPolicy` read access to an external auditor; CloudTrail/AWS Config alarms on policy mutation, delivered somewhere the operator does not control. Document that full policy immutability (no `PutKeyPolicy` principal) conflicts with the upgrade flow and is not an option here.
- **Emergency runbook (D4):**
  - Vulnerable image is the newly-swapped one → instant rollback (`proposePcr0Swap(PCR0_old)` → `applyPcr0Swap()`), then retire: `removePcr0(PCR0_new)` + strip from KMS policy.
  - Ordering trap: `removePcr0` reverts on the active image (`CannotRemoveActiveImage`) — swap `activeImage` away **first**, then remove.
  - Stripping a PCR0 from the KMS policy does **not** evict a live enclave already holding the recovered keys in memory — always pair with `nitro-cli terminate-enclave`.
  - All accepted images vulnerable → full stop (terminate enclave, strip KMS policy), wait out the timelock for the fixed image; downtime = `pcr0UpgradeDelay` by design. Trigger the intake pause (if adopted per Task 0) to close the exposure window while waiting.
  - Note the disclosure dynamic: `proposePcr0Swap` of the fix starts a public clock and reproducible builds let an attacker diff the fix — argue for a modest delay + pause rather than a long delay.
- No code change (`pkg/executor/kms` untouched); AWS-side configuration + documentation.

**Files:**
- `docs/ops/KMS_KEY_POLICY_RUNBOOK.md` (new)
- `docs/design/EXECUTOR_TEE_UPGRADE_DESIGN.md` (cross-reference), `CLAUDE.md` (pointer)

---

## Task 9 — Integration / System Tests — **IMPLEMENTED (CI-runnable subset)**

**Summary:** End-to-end upgrade and rollback scenarios.

**Description (as implemented):**

Added `tests/system/tee_upgrade_system_test.go` covering the upgrade guarantees that can be
driven over TCP. The **PCR0-swap-triggered** drain is Nitro-only — a TCP/dev executor reports an
empty PCR0 (the dev marker), so `maintainExecutorForActiveImage` skips swap observation and that
specific trigger cannot fire without real NSM hardware — but the **reconnect / re-handshake path
it drains into** is the same one an executor crash/relaunch takes, and that IS exercised e2e
below. The four e2e tests are:

- **R2 key continuity** (`TestTeeUpgrade_KeysetRecoveredAcrossRestart`, WASM, `CI_FLAG`-skipped):
  deploy the `simple` app → register user → deposit → `RestartAll` (coupled manager+executor
  restart, LevelDB recovery blob persisted) → assert the app state root survives, then a second
  deposit dispatches and its update-payload signature still validates against the **original**
  signing key. A regenerated keyset would fail that signature check — this is the direct proof
  that the keyset was recovered, not regenerated.
- **Executor-only reconnect** (`TestTeeUpgrade_ExecutorRestartReconnects`, mock-runtime,
  CI-runnable): the production crash/relaunch path (Task 5). Only the executor restarts; the
  manager stays up, its polling loop detects the dropped channel (`!IsConnected()`),
  re-dials + re-handshakes (`reconnectExecutor`), and the executor recovers its keyset from the
  still-running manager's blob. With on-chain `teeSigner` pinned to the original signer, a
  successful post-restart deploy plus a passing signer-continuity check (no fatal) proves recovery.
  This drives the actual drain-target reconnect logic rather than the coupled-restart approximation.
- **Task 6 expect-existing guard** (`TestTeeUpgrade_ExpectExistingKeyset_HandshakeFailsFast`,
  mock-runtime, CI-runnable): `EXECUTOR_EXPECT_EXISTING_KEYSET` set + empty manager store →
  the executor aborts the handshake, `StartManager` fails with `executor handshake failed`, and
  nothing is generated or stored.
- **D2 signer continuity** (`TestTeeUpgrade_SignerMismatch_Fatal`, mock-runtime, CI-runnable):
  no recovery data + guard off → the executor generates a fresh keyset (new signer); the on-chain
  `teeSigner` is pre-set to a different address; the manager's polling loop surfaces a fatal
  `SignerContinuityError` on `FatalErrChan()` rather than silently accepting the orphaning keyset.

**Test infrastructure added** (`pkg/testutil/`): mock-suite `RestartAll` (reuses the same
`MockClient` — its on-chain state and `teeSigner` survive `Close` — and re-opens the persisted
LevelDB), `TestSuiteCore.RestartExecutorOnly` (relaunch only the executor so the manager's
reconnect loop re-dials — works with any data-layer backend since the manager never stops),
`SetTeeSigner` / `SetActiveImage` pass-throughs, `TestSuiteCore.EventChannel()`, and
`TestSuiteCore.ManagerFatalErrChan()`.

**Deferred (Nitro-only / covered elsewhere):**
- The PCR0-swap **trigger** for the drain (propose → timelock → apply → drain) and rollback
  (`activeImage` re-point): Nitro-only (needs a real PCR0). The reconnect path the drain feeds
  into is covered by `TestTeeUpgrade_ExecutorRestartReconnects` above; the drain-decision
  state-machine is covered by Task 5 unit tests.
- KMS policy missing new PCR0; NSM PCR0 read; KMS recovery under the new PCR0: Nitro-only.
- Incompatible protocol versions: both peers share the `WireProtocolVersion` constant, so it is
  covered by `pkg/communication` + `pkg/manager` unit tests rather than e2e.
- On Nitro, exercise the Task 6 guard through the **baked-in** image ENV (upgrade-variant EIF +
  wiped Manager data folder → handshake aborts), not only via TCP-mode env vars.
- Challenge-attestation freshness (D1) and intake pause (D4): both **not adopted** (Task 0),
  so there is nothing to drill here.

**Files:**
- `tests/system/tee_upgrade_system_test.go` (new — 4 e2e tests)
- `pkg/testutil/test_helpers.go`, `pkg/testutil/test_suite_core.go` (`RestartAll`, `RestartExecutorOnly`, getters)

---

## Documentation checklist (per repo convention)

- `docs/design/EXECUTOR_TEE_UPGRADE_DESIGN.md` — corrections from Task 0; R2 enforcement guards (Task 6) ✓
- `docs/design/EXEC_MGR_HANDSHAKE.md` — new handshake fields (PCR0, version, protocol version); key-continuity guard (Task 6) ✓
- `CLAUDE.md` / `README` — new config vars (`EXECUTOR_EXPECT_EXISTING_KEYSET` ✓), admin verb, upgrade runbook pointer
- `NOTICES` — only if `go.mod` changes (NSM `DescribePCR` should not need new deps)
