# Upgradable Executor (TEE) Software

## Summary

This document describes the design to upgrade the **Executor** software that runs inside the AWS Nitro Enclave (the TEE), without losing the encrypted state of any application and without breaking the chain-of-trust that anchors the system to a known enclave image.

Upgrading the Executor is fundamentally harder than upgrading any other component because the Executor binary is measured into **PCR0**, and PCR0 is the trust anchor of the whole platform:

- the on-chain `TeeAuthenticator` contract stores `pcr0` and rejects any attestation that does not match it;
- the AWS KMS key policy gates every decryption on a `kms:RecipientAttestation:PCR0` condition.

Changing the Executor binary changes PCR0, which simultaneously invalidates the on-chain attestation path and the KMS key-recovery path. This design makes that change a **controlled, reversible, auditable** operation.

The chosen approach is:

- **On-chain PCR0:** replace the single `pcr0` value with a **set of accepted PCR0 values**, each updated through a **timelock**, so old and new images are valid simultaneously during a grace window.
- **Orchestration:** drain pending requests without accepting new ones before the switch, triggered by an **on-chain action on `TeeAuthenticator`**, observed by the Manager's existing polling loop (no in-process hot-reload, no separate off-chain command). Every upgrade is therefore recorded on-chain.
- **Key continuity:** the enclave key set is **recovered, never regenerated**, by allowing the new PCR0 in the KMS key policy. No re-encryption of state or of the master key is required.
- **Verifiability:** the enclave image build must be **reproducible**, so anyone can rebuild the EIF from a tag and confirm the PCR0 registered on-chain.

The Manager (which runs outside the TEE) is out of scope here except for two touchpoints it owns in this flow: observing the on-chain swap signal in its polling loop (see [R8](#r8--on-chain-triggered-graceful-drain)) and the protocol-version handshake it shares with the Executor (see [R7](#r7--protocol-version-negotiation)).

---

## Current State

### What identifies an Executor on-chain

`TeeAuthenticator.sol` holds **two** distinct identities of the running enclave:

| On-chain value | Meaning | Source |
|----------------|---------|--------|
| `bytes public pcr0` | Identity of the enclave **image** (the code) | The EIF build |
| `address public teeSigner` + `bytes public pubSecp521r1` | Identity of the enclave **keys** | Derived from the `EnclaveKeySet` |

- `pcr0` is set in the constructor and can be changed by `updatePcr0(bytes newPcr0)` (`TeeAuthenticator.sol`), which emits `PcrZeroUpdate` and overwrites the single stored value atomically.
- `updateTee(bytes attestation)` (`TeeAuthenticator.sol`) verifies a fresh attestation, checks its PCR field against the **single** stored `pcr0` in `_checkAttestationContent`, and on success records `teeSigner` and `pubSecp521r1` from the attestation's `userData`/`enclaveKey`.
- `nitroProver` and `maxVerificationAge` are `immutable`.

There is exactly **one** valid `pcr0` at any time. Updating it is an instantaneous flip: there is no window in which both the old and the new image are accepted.

### How the enclave recovers its keys

The enclave key set (master AES key plus the derived P521 / secp256k1 / X25519 keys) is the root of all encrypted application state. It is established once and **recovered** on every subsequent start via the handshake (`docs/design/EXEC_MGR_HANDSHAKE.md`):

- The Manager persists opaque recovery data and returns it to the Executor on connect (`performHandshake` in `pkg/executor/executor.go`).
- For the KMS-backed recovery type, the master key is stored as a KMS-encrypted `CiphertextBlob`. On recovery the enclave generates a fresh attestation and calls `DecryptWithAttestation(blob, attestationDoc)` (`pkg/executor/kms/kms_client.go`). **KMS decides whether to decrypt based solely on the key policy's PCR0 condition matched against the attestation — not on the contents of the blob.**

This is the linchpin of the whole upgrade design (see [Key Continuity](#key-continuity-the-linchpin)).

### How the Executor is built and run

- The Executor binary is built from `dockerfiles/executor/Dockerfile`, with the git version embedded via ldflags into `pkg/version/version.go` (`Version`, defaulting to `"dev"`). The Go layer is already deterministic (pinned Go version, `go.sum`-locked deps), but **PCR0 is not yet independently reproducible from a git tag**: the base image uses the moving `amazonlinux:2` tag with unpinned packages, `CGO_ENABLED=1` links the binary to that base image's libc, and — most importantly — the `.eif` packaging step (`nitro-cli build-enclave`), which is what PCR0 actually measures, is external to this repo and unversioned. Making PCR0 reproducible is [R5](#r5--reproducible-enclave-image).
- The binary is packaged into a Nitro **Enclave Image File (.eif)** by an external `nitro-cli build-enclave` step (not in this repo). PCR0/PCR1/PCR2 are produced by that build.
- The Executor and Manager are separate processes/containers communicating over vsock (production) or TCP (dev). There is **no in-process reload**: an upgrade means stopping the old enclave and launching a new EIF. The handshake on reconnect is **blocking and has no timeout** (in `pkg/manager/manager.go`).
- The version string is informational only (logged at startup, queryable via the admin CLI `get_version` command); it is **not** cryptographically bound and **cannot** be used to enforce anything. PCR0 is the only enforcement point.

---

## Goals

- Replace the Executor binary / enclave image with a new version while **preserving all encrypted application state and the enclave key set**.
- Allow the old and new enclave images to be **simultaneously valid** for a bounded grace window (zero-downtime swap and clean rollback).
- Make every PCR0 change **timelocked and observable on-chain**, giving users a window to audit the new image and withdraw funds before it becomes authoritative.
- Keep the enclave image **reproducibly buildable** so the on-chain PCR0 can be independently verified.
- Detect and cleanly **reject incompatible Manager/Executor version pairs** instead of corrupting state or hanging.

## Non-Goals

- Upgrading the **Manager** binary (standard container/binary swap; this design only adds the two Manager touchpoints it needs — the swap-signal read of [R8](#r8--on-chain-triggered-graceful-drain) and the protocol-version handshake of [R7](#r7--protocol-version-negotiation)).
- In-process hot-reload of the Executor (always a process/enclave restart).
- Migrating or transforming persisted application state as part of an upgrade (a state migration is a separate operation; this design only guarantees the new image can *read* the old state).
- Changing the cryptographic scheme of the key set, or re-keying the master key.
- Automating the EIF build / `nitro-cli` pipeline itself (out of repo).
- Upgrading the smart contracts' proxy mechanism (covered by `UPGRADABLE_CONTRACTS_DESIGN.md`); this design only adds new logic to `TeeAuthenticator`, which becomes feasible *because* that contract is becoming UUPS-upgradable.

---

## Requirements

### R1 — No State Loss

An Executor upgrade must not lose or require re-encryption of any application state. The new image must read the existing encrypted state with the recovered key set.

### R2 — Key Continuity (Recover, Never Regenerate)

The new enclave image must **recover** the existing `EnclaveKeySet`, never generate a new one. As a corollary, `teeSigner` and `pubSecp521r1` must remain unchanged across a pure software upgrade — only `pcr0` changes. (If the keys were regenerated, both on-chain identities would change *and* all existing state would become permanently undecryptable.)

### R3 — Overlapping Validity Window

There must exist a bounded window during which **both** the old and the new PCR0 are accepted, on-chain and in the KMS key policy. This enables a zero-downtime swap and a clean rollback to the previous image.

### R4 — Timelocked, Observable Swaps

Activating a **previously-unseen** PCR0 must be a two-step, timelocked operation that emits events at proposal and at application time — no new image may become runnable instantly, so users get an audit window. Swapping to a PCR0 that is **already in the accepted set** (e.g. a rollback to the prior image) carries no un-audited code and therefore needs no timelock: it may be applied immediately. Every swap, timelocked or not, is recorded on-chain.

### R5 — Reproducible Enclave Image

The enclave image must be reproducibly buildable from a git reference so that the PCR0/PCR1/PCR2 registered on-chain can be independently recomputed and verified by third parties during the timelock window.

### R6 — Controlled Upgrade Authority

The functions that add/remove accepted PCR0 values must be restricted to the contract owner/admin (the same authority model and multisig/timelock guidance as `UPGRADABLE_CONTRACTS_DESIGN.md`).

### R7 — Protocol-Version Negotiation

The Manager↔Executor handshake must carry a protocol version / capability descriptor. Incompatible pairs must be rejected with an explicit error rather than blocking indefinitely or corrupting persisted recovery data.

### R8 — On-Chain-Triggered Graceful Drain

The swap must be initiated by an on-chain action on `TeeAuthenticator` (so every upgrade is recorded), observed idempotently by the Manager. On that signal, in-flight requests must be drained before the old enclave is stopped; the swap must not interrupt a request mid-execution. There is no separate off-chain upgrade command and no designed emergency-stop verb.

---

## Key Continuity: the Linchpin

The single fact that makes this design tractable:

> The master key is stored as a KMS-encrypted `CiphertextBlob`. KMS will decrypt it for **any** enclave whose attestation satisfies the key policy's PCR0 condition. The blob itself is **not** bound to a specific PCR0.

Therefore a new Executor image (new PCR0) can decrypt the **same, unchanged** recovery blob — provided the KMS key policy is updated to also accept the new PCR0. No re-encryption of the master key, and no re-encryption of any application state, is ever required.

```
Old enclave (PCR0_old) ──attestation(PCR0_old)──► KMS ──┐
                                                         ├─ policy accepts {PCR0_old, PCR0_new}
New enclave (PCR0_new) ──attestation(PCR0_new)──► KMS ──┘
                                                         │
                              same CiphertextBlob ───────┘ ──► master key ──► all state readable
```

Concretely, the only inputs that must change to admit a new image are:

1. **KMS key policy** — add `PCR0_new` to the `kms:RecipientAttestation:PCR0` condition (keep `PCR0_old` during the grace window).
2. **On-chain `TeeAuthenticator`** — add `PCR0_new` to the set of accepted PCR0 values.
3. **Manager/Executor protocol** — ensure the version pair is compatible.

The "unsafe"/dev recovery type (master key stored in plaintext recovery data) is trivially portable and needs none of the KMS steps.

---

## On-Chain Changes: `TeeAuthenticator`

### From a single PCR0 to an accepted set with a swap pointer

Replace the single `bytes public pcr0` with a **set** of accepted values, a **pointer** to the one that should currently be running (`activeImage`), and a timelocked staging area for the next swap.

```solidity
// --- replaces `bytes public pcr0` ---

// Accepted PCR0 values, keyed by keccak256(pcr0) for O(1) membership checks.
mapping(bytes32 => bool) public acceptedPcr0;
bytes32[] public acceptedPcr0List;          // for enumeration / off-chain audit

// The PCR0 the Manager should currently be running (keccak256 of the target).
bytes32 public activeImage;

// Timelock staging for a pending swap.
struct PendingSwap {
    bytes   value;
    uint256 eta;        // earliest application time
    bool    pending;
}
PendingSwap public pendingSwap;
uint256 public pcr0UpgradeDelay;             // timelock duration; set once in initialize(), NO setter
```

`pcr0UpgradeDelay` is deliberately fixed at deployment (set in `initialize()`, with no owner-callable setter). A runtime-mutable delay would let a compromised owner shorten it to zero and then push a malicious image instantly, defeating [R4](#r4--timelocked-observable-swaps). Different networks can still pick different values by passing a different `initialize()` argument (e.g. a shorter delay on testnet); changing it after deployment would require a full contract upgrade, which is itself governed by the timelocked upgrade authority of `UPGRADABLE_CONTRACTS_DESIGN.md`. (In a non-upgradable contract this would be a constructor-set `immutable`; under UUPS the constructor does not run for proxy state, so the equivalent is an `initialize()`-set storage value with no setter — kept in storage rather than `immutable` for consistency with the immutable→storage conversion that design applies.)

### A single swap flow for upgrades and rollbacks

`proposePcr0Swap` / `applyPcr0Swap` drive **every** change of the running image. `applyPcr0Swap` adds the target to the accepted set if it is new, sets it as `activeImage`, and emits the signal the Manager observes. The timelock is enforced **only when the target is not yet accepted** — a swap to an already-accepted PCR0 (used in case of a rollback to the previously-audited image) applies immediately.

```solidity
function proposePcr0Swap(bytes calldata targetPcr0) external onlyOwner {
    pendingSwap = PendingSwap({
        value: targetPcr0,
        eta: block.timestamp + pcr0UpgradeDelay,
        pending: true
    });
    emit Pcr0SwapProposed(targetPcr0, pendingSwap.eta);
}

function applyPcr0Swap() external onlyOwner {
    PendingSwap memory p = pendingSwap;
    if (!p.pending) revert NoPendingSwap();
    bytes32 key = keccak256(p.value);

    // New images must clear the audit window; already-accepted images may swap now.
    if (!acceptedPcr0[key]) {
        if (block.timestamp < p.eta) revert TimelockNotElapsed();
        acceptedPcr0[key] = true;
        acceptedPcr0List.push(key);
    }

    activeImage = key;
    delete pendingSwap;
    emit Pcr0Swapped(p.value);   // observed by the Manager
}

// Removal only tightens security (it retires a vetted image), so it is not timelocked.
// Cannot remove the last entry, nor the image that is currently active.
function removePcr0(bytes calldata oldPcr0) external onlyOwner {
    bytes32 key = keccak256(oldPcr0);
    if (!acceptedPcr0[key]) revert UnknownPcr0();
    if (acceptedPcr0List.length == 1) revert CannotRemoveLastPcr0();
    if (key == activeImage) revert CannotRemoveActiveImage();
    acceptedPcr0[key] = false;
    // swap-remove from acceptedPcr0List …
    emit Pcr0Removed(oldPcr0);
}
```

The Manager observes `Pcr0Swapped` / reads `activeImage` and drains-and-swaps to the target (see [On-chain upgrade trigger](#on-chain-upgrade-trigger-r8)). Because `applyPcr0Swap` only adds to the set behind the `!acceptedPcr0[key]` guard, re-targeting a value already in the set never duplicates an entry — which is exactly what makes a rollback (`proposePcr0Swap(PCR0_old)` → immediate `applyPcr0Swap()`) cheap and on-chain-tracked. The same `onlyOwner`/multisig authority of [R6](#r6--controlled-upgrade-authority) applies to all three functions.

### Attestation check against the set

`_checkAttestationContent` (`TeeAuthenticator.sol`) currently compares the attestation's PCR field byte-by-byte against the single `pcr0`. It must instead extract the PCR0 bytes from `pcrs` (offset 4) and test membership in `acceptedPcr0`:

```solidity
function _checkAttestationContent(bytes memory pcrs, bytes memory enclaveKey, bytes memory userData)
    internal view
{
    if (userData.length != 20) revert InvalidUserDataLength();
    if (enclaveKey.length != PK_LENGTH) revert InvalidPKLength();

    bytes memory candidate = _extractPcr0(pcrs);   // pcrs[4 : 4 + PCR0_LENGTH], PCR0 is a fixed 48-byte SHA-384
    if (!acceptedPcr0[keccak256(candidate)]) revert InvalidPCR();
}
```

All accepted PCR0 values share the same fixed length (48 bytes), so the per-byte loop and `pcr0.length` of the original single-value check are replaced by a fixed-length extraction plus an O(1) set lookup.

This keeps `updateTee` / `updateTeeStep1..4` working unchanged in shape: they continue to record `teeSigner`/`pubSecp521r1` from the attestation, but the embedded `_checkAttestationContent` now accepts any image in the valid set.

**A software upgrade does not re-run `updateTee`.** Today `updateTee` is a one-time first-install operation — it is invoked only by the manual `contracts/scripts/management/updateTee.ts` operator script (and by tests); no Manager/Executor/runtime code calls it. Per [R2](#r2--key-continuity-recover-never-regenerate) a pure software upgrade recovers the *same* key set, so `teeSigner`/`pubSecp521r1` are unchanged and the existing registration stays valid across the swap — the signature path in `ProcessorEndpoint.checkSignature` keeps verifying against the same signer. The only scenario that re-registers the signer is a catastrophic key-set loss requiring a brand-new attestation chain, which is explicitly a non-goal here.

### Authority and timelock

`proposePcr0Swap` / `applyPcr0Swap` / `removePcr0` are `onlyOwner`, mirroring the controlled-upgrade-authority requirement (R5) of `UPGRADABLE_CONTRACTS_DESIGN.md`. The owner should be a multisig, and `pcr0UpgradeDelay` provides the same user-protection property that the contract-upgrade timelock provides: a compromised owner cannot make a **new** image authoritative without a publicly visible delay during which users can withdraw. (Swaps to an already-accepted image skip the delay precisely because that image already cleared an audit window — switching to it introduces no new code.)

---

## Off-Chain Changes

### KMS key policy

The KMS key policy's `kms:RecipientAttestation:PCR0` condition must list the set of accepted PCR0 values. Adding `PCR0_new` (keeping `PCR0_old`) is the off-chain analogue of accepting it on-chain via `applyPcr0Swap` and **must be performed before** the new enclave starts, otherwise the new image's handshake recovery (`DecryptWithAttestation`) will fail. Removing `PCR0_old` is the analogue of `removePcr0` and happens only after the upgrade is confirmed stable.

This is purely an AWS-side configuration change (key policy update); no code change in `pkg/executor/kms`.

### Protocol-version handshake (R7)

Extend the handshake (`pkg/executor/executor.go` `performHandshake`, and the message types in `pkg/communication`) with a `protocolVersion` (and optionally a capability bitmap). On connect:

- Manager and Executor exchange their `protocolVersion`.
- If the pair is incompatible, fail the handshake with a typed error and **do not** mutate persisted recovery data.
- Add a **timeout** to the Manager's blocking wait (in `pkg/manager/manager.go`) so a failed/incompatible swap surfaces as an explicit error rather than a hang.

Because Manager and Executor may be upgraded independently, define and document a small compatibility matrix keyed on `protocolVersion`.

### On-chain upgrade trigger (R8)

The swap is initiated by an **on-chain action on `TeeAuthenticator`**, not by a separate off-chain command. This keeps a single, immutable record of every upgrade *and rollback* (who, when, target image) under the same authority that already governs the PCR0 set, and reuses the Manager's existing connection to that contract.

The trigger is the `activeImage` pointer set by `applyPcr0Swap`, announced by the `Pcr0Swapped` event (see [A single swap flow](#a-single-swap-flow-for-upgrades-and-rollbacks) above). `activeImage` names the PCR0 the Manager should be running; the Manager compares it against the image it is actually running and, on a mismatch, drains and swaps.

The Manager **already holds the `TeeAuthenticator` binding** (`teeAuthBoundContract` in `pkg/blockchain/client.go`, configured from `CHAIN_TEEAUTHENTICATOR_ADDRESS`); today it is dormant — only the test-only `GetTeePublicKey()` read path exists. The only addition is to read `activeImage` / watch `Pcr0Swapped` in the polling loop (`pkg/manager/manager.go`), next to `GetNextPendingRequest`. No new contract reference and no new configuration are required.

When `activeImage` differs from the running image, the Manager:

1. Stops dispatching new requests (stops consuming from `GetNextPendingRequest`).
2. Lets in-flight requests complete (drain).
3. Signals the Executor to exit cleanly and closes the channel.
4. An operator launches the EIF matching `activeImage` on the parent host. Staging the artifact and resolving `activeImage` (a PCR0 hash) to a concrete `.eif` file is a **manual operational procedure, out of scope for this document**.
5. The Manager reconnects; the handshake recovers the key set under the new PCR0 (now accepted by KMS) and resumes dispatching.

The trigger is **level-based and idempotent**: the Manager acts on the *state* of `activeImage`, not on a one-shot event, so a re-delivered `Pcr0Swapped`, a missed event, or a Manager restart all converge to the same outcome — it swaps iff the running image differs from `activeImage`. This also makes rollback uniform: pointing `activeImage` back to `PCR0_old` is just another mismatch the Manager resolves the same way.

The Executor binary version (`pkg/version`) should be returned in the handshake and in the admin CLI `get_version`, so the operator can confirm the running version after the swap.

> **No designed emergency stop.** There is deliberately no off-chain "stop now" verb. An operator with host access can always terminate the enclave out-of-band (`nitro-cli terminate-enclave` / stopping the supervisor); that is an ungraceful stop — in-flight requests are lost, but no state is lost because the key set and encrypted state are recovered on the next start. Removing a dedicated break-glass command removes a redundant trust path, not a capability.

### Reproducible build (R5)

The EIF build must be deterministic enough that PCR0/PCR1/PCR2 are a pure function of the source at a given git tag. Practically this means pinning the toolchain (Go 1.24, base image digest, build flags) and the `nitro-cli` version, and publishing — for each release — the git commit, the EIF, and the expected PCR measurements, so third parties can rebuild and compare against the on-chain set during the timelock window.

---

## Safe Upgrade Sequence

The dependency order that guarantees zero key loss and a clean rollback:

| # | Step | Reversible? |
|---|------|-------------|
| 1 | Reproducibly build the new EIF; derive `PCR0_new` (and PCR1/PCR2); publish for verification. | n/a |
| 2 | Add `PCR0_new` to the **KMS key policy** (keep `PCR0_old`). | yes (remove it) |
| 3 | `proposePcr0Swap(PCR0_new)` on-chain. Old PCR0 still the only accepted/active image. | yes (let the proposal expire) |
| 4 | *(Audit / timelock window — users may verify the image or withdraw.)* | — |
| 5 | After `pcr0UpgradeDelay`, `applyPcr0Swap()`: `PCR0_new` joins the accepted set and becomes `activeImage`; `Pcr0Swapped` is emitted. The Manager drains and stops the old enclave; an operator launches the new EIF (manual step); the handshake recovers the key set via KMS. Service resumes. Both PCR0s are now accepted. | yes (instant on-chain rollback — see below) |
| 6 | After a stability period: `removePcr0(PCR0_old)` on-chain and remove `PCR0_old` from the KMS policy. | this finalizes the upgrade |

**Rollback** before step 6 is an **instant, on-chain, tracked** operation: `proposePcr0Swap(PCR0_old)` → `applyPcr0Swap()`. The timelock is skipped because `PCR0_old` is still in the accepted set, so `applyPcr0Swap` re-points `activeImage` to it immediately; the Manager sees the mismatch and swaps back, recovering the same key set. After step 6 the old image is retired, so rolling back to it again requires a fresh **timelocked** `proposePcr0Swap(PCR0_old)` (it must re-clear the audit window). For a true emergency where the chain is unusable, the out-of-band host relaunch described under [On-chain upgrade trigger](#on-chain-upgrade-trigger-r8) remains available.

The Manager may be upgraded independently of this sequence, gated only by the [R7](#r7--protocol-version-negotiation) compatibility check.

---

## Failure Modes

| Failure | Cause | Mitigation |
|---------|-------|------------|
| New enclave can't recover keys | `PCR0_new` not in KMS policy (step 2 skipped/late) | Handshake fails fast (R7 timeout); operator adds PCR0 to policy or rolls back to old EIF |
| On-chain attestation rejected | `applyPcr0Swap` not yet executed (`PCR0_new` not in set) | `updateTee` reverts with `InvalidPCR`; old PCR0 still serves traffic |
| Manager/Executor version mismatch | Independent upgrades | Typed handshake error (R7); no state mutation |
| Request interrupted by swap | Drain skipped | Manager drains in-flight work on the on-chain swap signal before stop (R8) |
| Malicious image pushed by compromised owner | Owner key compromise | Timelock (R4) + reproducible build (R5) give users a window to detect and withdraw |
| All PCR0 removed | Operator error | `removePcr0` reverts on the last entry (`CannotRemoveLastPcr0`) |
| Active image accidentally removed | Fat-finger at finalization (step 6): passing `PCR0_new` instead of `PCR0_old` | `removePcr0` reverts (`CannotRemoveActiveImage`); this matters because the mistake would not cause an immediate outage (runtime uses `teeSigner`, not the set) but would leave `activeImage` dangling and block any future `updateTee` re-attestation |

---

## Implementation Notes

### Swap-signal observation: poll `activeImage`, don't subscribe to `Pcr0Swapped`

The Manager learns it must swap by **reading the `activeImage` state variable on its existing polling loop**, not by subscribing to the `Pcr0Swapped` event. The event is retained purely as an on-chain/subgraph audit record; it is not on the Manager's control path.

Concretely, on each tick the Manager reads `activeImage` and, **before** calling `GetNextPendingRequest`, compares it against the image it is running: on a match it dispatches as today; on a mismatch it enters drain mode (stops consuming new requests, lets in-flight work finish, signals the Executor to exit). Gating dispatch on the check ensures a pending swap immediately halts new work.

**Poll interval / drain latency.** Reuse the existing `BlockchainPollingInterval`; do not add a separate cadence. Drain latency is bounded by one poll interval, which is acceptable for an operator-initiated, timelocked upgrade.

### PCR1/PCR2

We consider PCR0 sufficient as the trust anchor for now. PCR1 (kernel) and PCR2 (application) were never checked, and this implementation does not add them.

### Protocol-version matrix: a code constant, not a doc table

The compatibility matrix for `(managerVersion, executorVersion)` lives in **code as the source of truth**, not as a table in this document. The reason is that the compatibility check is *enforced* at handshake time ([R7](#r7--protocol-version-negotiation)): a doc table cannot reject anything and would inevitably drift from the code that does.

Concretely:

- Introduce a dedicated, monotonic `ProtocolVersion` constant in `pkg/communication`, **separate from** the informational `version.Version` git string in `pkg/version` (which the [Current State](#how-the-executor-is-built-and-run) explicitly marks as non-enforcing). The handshake carries `ProtocolVersion`, not the git version.
- Both sides expose their `ProtocolVersion` in the handshake, and a single `IsCompatible(peer)` function encodes the rule (exact-match, or a documented min-supported floor as the protocol evolves). That function is the matrix.
- This doc records only *where* the constant and the compatibility rule live; it does not duplicate the pairs, so there is nothing to keep in sync by hand.

