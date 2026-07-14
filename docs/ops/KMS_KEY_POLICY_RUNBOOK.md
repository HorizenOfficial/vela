# KMS Key Policy Runbook & Governance (Ops)

**Scope:** operational procedures for the AWS KMS key policy that gates Executor
key recovery, its governance controls, and the TEE-upgrade / rollback /
emergency runbooks. This is an **AWS-side configuration + process** document —
there is no code change in `pkg/executor/kms`.

**Design context:** `docs/design/EXECUTOR_TEE_UPGRADE_DESIGN.md` (esp. *Key
Continuity: the Linchpin* and *Off-Chain Changes → KMS key policy*),
`docs/design/REPRODUCIBLE_EIF_BUILD.md` (deriving/verifying PCR0), and Task 1/5/6
in `docs/design/EXECUTOR_TEE_UPGRADE_TASKS.md`.

---

## 1. Why the KMS key policy matters

The Executor recovers its `EnclaveKeySet` (KMS recovery type) by presenting a
fresh Nitro attestation to KMS and calling `DecryptWithAttestation` on the
stored `CiphertextBlob` (`pkg/executor/kms/kms_client.go`). **KMS decides whether
to decrypt based solely on the key policy's `kms:RecipientAttestation:PCR0`
condition matched against the attestation — not on the blob contents.** The blob
is *not* bound to a PCR0, so any enclave whose PCR0 the policy accepts can
recover the same, unchanged key set. This is what makes a software upgrade
key-continuous (R2: recover, never regenerate).

Consequently the KMS key policy is the **off-chain twin** of the on-chain
accepted PCR0 set in `TeeAuthenticator`:

| Trust root | Where | Governs |
|---|---|---|
| `acceptedPcr0` / `activeImage` | `TeeAuthenticator` (on-chain) | which image the platform *considers* valid / active, and the timelock/audit lifecycle |
| `kms:RecipientAttestation:PCR0` condition | AWS KMS key policy | which image can actually *recover the keys* |

These two lists **must be kept in sync**. They are the only place the two can
diverge, and a divergence is the core residual risk (D1): a malicious infra
admin who adds a rogue PCR0 to the KMS policy but not on-chain gets an image that
can silently recover the keys — sign indistinguishably and/or passively decrypt
state. The latter is undetectable by any protocol; the mitigations here are
**preventive** (governance + auditability), see §3.

> The **unsafe/dev** recovery type (`EXECUTOR_KEYSET_RECOVERY_TYPE=0`, master key
> in plaintext recovery data) needs none of the KMS steps below — it is trivially
> portable. Everything in this document applies to the **KMS** recovery type
> (`EXECUTOR_KEYSET_RECOVERY_TYPE=1`, production).

### Key facts to keep straight

- **PCR0** is a SHA-384 digest (48 bytes ⇒ **96 hex chars**). The KMS condition
  value and the on-chain `bytes` (preimage) value are the same hex string — but
  the KMS condition value is **un-prefixed**, while on-chain tooling usually shows
  it `0x`-prefixed. Drop the `0x` when copying a PCR0 into the key policy.
- **Preimage vs hash.** The KMS condition and the `proposePcr0Swap` /
  `applyPcr0Swap` arguments use the raw 48-byte PCR0 **preimage**. On-chain,
  `acceptedPcr0` / `acceptedPcr0List` / `activeImage` store the **`keccak256(pcr0)`
  hash** (`bytes32`), *not* the preimage. So the two lists are never byte-equal —
  compare them by hashing (§3). The hash→preimage mapping for audits comes from
  the `Pcr0Swapped` / `Pcr0SwapProposed` / `Pcr0Removed` events, which emit the
  raw preimages non-indexed for exactly this purpose (Task 1).
- The on-chain PCR0 (`proposePcr0Swap` argument) is derived from the reproducible
  EIF build: `jq .Measurements.PCR0 eif-out/measurements.json`
  (`docs/design/REPRODUCIBLE_EIF_BUILD.md`).
- Release images are the **upgrade variant** (`EXECUTOR_EXPECT_EXISTING_KEYSET=true`
  baked in — Task 6). The **genesis** variant (guard off) is used only for the
  one-time first install and has a *different* PCR0; retire it afterward (§7).
- The parent EC2 must run the KMS vsock proxy the enclave forwards through
  (`EXECUTOR_KMS_PROXY_PORT`, default 8000) and the caller's IAM role must be a
  policy principal (below). These are launch prerequisites, not per-upgrade steps.

---

## 2. Key policy structure

The policy lists **every currently accepted PCR0** in a single
`kms:RecipientAttestation:PCR0` condition (multiple values = OR). Representative
enclave-recovery statement:

```json
{
  "Sid": "EnclaveKeyRecovery",
  "Effect": "Allow",
  "Principal": { "AWS": "arn:aws:iam::<ACCT>:role/<manager-parent-instance-role>" },
  "Action": ["kms:Decrypt", "kms:GenerateDataKey"],
  "Resource": "*",
  "Condition": {
    "StringEqualsIgnoreCase": {
      "kms:RecipientAttestation:PCR0": [
        "<PCR0_old — 96 hex chars>",
        "<PCR0_new — 96 hex chars>"
      ]
    }
  }
}
```

During a grace window this list holds **both** the old and new PCR0 (R3). At
steady state it holds exactly the accepted set (normally one entry).

Read / write the policy (the `kms-policy-admin` role, §3):

```bash
# Read current policy
aws kms get-key-policy --key-id "$KMS_KEY_ARN" --policy-name default \
  --query Policy --output text > policy.json

# Edit policy.json (add/remove a PCR0 value), then apply
aws kms put-key-policy --key-id "$KMS_KEY_ARN" --policy-name default \
  --policy file://policy.json
```

---

## 3. KMS policy governance (D1)

The key policy is the single off-chain point where the accepted-image set can be
diverged from the on-chain set, so it needs governance at least as strong as the
on-chain `onlyOwner`/multisig that gates `proposePcr0Swap` / `applyPcr0Swap` /
`removePcr0`.

- **Separate `kms:PutKeyPolicy` from ops.** Grant policy-mutation rights only to
  a dedicated, multi-party role (`kms-policy-admin`) — ideally the same
  quorum/multisig that owns `TeeAuthenticator`, and **not** the role that
  launches enclaves or runs the Manager. Day-to-day operators must not hold
  `PutKeyPolicy`.
- **External auditor read access.** Grant `kms:GetKeyPolicy` to an
  independent auditor role so the off-chain set can be compared against the
  on-chain `acceptedPcr0List` without operator cooperation.
- **Alarm on mutation, out-of-band.** `PutKeyPolicy` is a CloudTrail management
  event. Wire an EventBridge rule (and/or an AWS Config rule on the key
  resource) → SNS to a channel the enclave operator **does not control**, so a
  unilateral policy change is noticed even if local logs are tampered with.
- **Full immutability is not an option here.** A policy with *no* `PutKeyPolicy`
  principal cannot be updated, which breaks the upgrade flow (you could never add
  `PCR0_new`). The goal is *governed mutability*, not immutability.

**Continuous invariant to alarm on.** The KMS condition holds raw PCR0
**preimages** while `acceptedPcr0List` holds their **`keccak256` hashes**, so the
check is by hashing, not literal equality:

- for every value `v` in `kms:RecipientAttestation:PCR0`,
  `keccak256(hexDecode(v))` must appear in the on-chain `acceptedPcr0List`
  (enumerate via `acceptedPcr0List(i)` / `getAcceptedPcr0Count`), **and the counts
  must match** (no on-chain entry without a KMS preimage, and vice-versa); and
- exactly one KMS value must hash to the on-chain `activeImage`.

Recover the preimage for each on-chain hash from the `Pcr0Swapped` /
`Pcr0SwapProposed` / `Pcr0Removed` event data (raw, non-indexed — Task 1). Any
drift is an incident.

---

## 4. Safe upgrade runbook

Upgrading the Executor binary changes PCR0, invalidating **both** the KMS
recovery path and the on-chain attestation path — so the KMS policy and the
on-chain set must both be updated, **KMS first**. A software upgrade does **not**
re-run `updateTee`: per R2 the key set (hence `teeSigner` / `pubSecp521r1`) is
unchanged, so the existing signer registration stays valid across the swap.

Let `PCR0_old` = current `activeImage`, `PCR0_new` = the upgrade target.

1. **Build & publish the new image.** Build the upgrade-variant EIF from the
   release tag and record `PCR0_new`:
   ```bash
   ./dockerfiles/executor/build-eif.sh <tag> ./eif-out
   jq -r .Measurements.PCR0 eif-out/measurements.json   # = PCR0_new
   ```
   Publish tag + `PCR0`/`PCR1`/`PCR2` + `nitro-cli` version + base digest so third
   parties can reproduce and verify during the timelock (R5).
2. **KMS policy: add `PCR0_new`, keep `PCR0_old`** (§2). This is the off-chain
   analogue of accepting the image and **must happen before the new enclave
   starts**, or its first handshake recovery (`DecryptWithAttestation`) fails.
   Performed by `kms-policy-admin` (§3).
3. **On-chain: propose the swap.**
   `proposePcr0Swap(PCR0_new)` — emits `Pcr0SwapProposed(PCR0_new, eta)` with
   `eta = now + pcr0UpgradeDelay`. Starts the public audit window.
   - A live proposal blocks a second one (`SwapAlreadyPending`); `cancelPcr0Swap()`
     first if you need to re-target.
4. **Wait out the timelock.** Third parties verify `PCR0_new` reproduces from the
   published tag. The proposal must be applied within `PCR0_SWAP_APPLY_WINDOW`
   (7 days) after `eta`, else it expires (`SwapProposalExpired`) and must be
   re-proposed.
5. **On-chain: apply the swap.** After `eta`: `applyPcr0Swap()` — adds `PCR0_new`
   to `acceptedPcr0`, sets `activeImage = keccak256(PCR0_new)`, emits
   `Pcr0Swapped(PCR0_new)`.
6. **Relaunch the enclave.** The Manager observes `activeImage` ≠ running image,
   drains in-flight work, signals the Executor to exit, and closes the channel
   (Task 5). An operator launches the `PCR0_new` EIF on the parent host. The
   Manager reconnects, the handshake recovers the **same** key set under
   `PCR0_new` (now accepted by KMS), and dispatch resumes.
   - Confirm the running version via the admin CLI `get_version` (`executorPcr0`
     surfaces the handshake-reported PCR0).
   - Because the reconnecting enclave recovered (not regenerated) the key set,
     its signer still matches on-chain `teeSigner`; a mismatch is caught as a
     fatal `SignerContinuityError` (D2, Task 5) — treat as an incident, do **not**
     proceed to finalize.
7. **Finalize (retire the old image)** once the upgrade is confirmed stable:
   - On-chain: `removePcr0(PCR0_old)`.
     **Ordering trap:** this reverts with `CannotRemoveActiveImage` if `PCR0_old`
     is still active — it can only run *after* the swap in step 5 moved
     `activeImage` to `PCR0_new`.
   - KMS policy: remove `PCR0_old` from the condition (§2), mirroring `removePcr0`.

Until step 7 completes, both images are accepted on-chain and in KMS — this is
the R3 overlap window that makes rollback (§5) instant.

---

## 5. Rollback runbook

Rolling back to the **previously accepted** `PCR0_old` needs **no timelock**: it
is already in the accepted set, so `applyPcr0Swap` re-points `activeImage`
immediately (the `!acceptedPcr0[key]` guard skips the timelock and adds nothing).

Only valid while `PCR0_old` is still accepted — i.e. **before** step 7 of the
upgrade retired it. If it was already removed, a rollback is a full re-upgrade
(§4 with old and new swapped).

1. **On-chain:** `proposePcr0Swap(PCR0_old)` → `applyPcr0Swap()` (back-to-back;
   no wait). `activeImage` returns to `keccak256(PCR0_old)`.
2. **KMS policy:** no change needed — `PCR0_old` is still listed during the grace
   window.
3. **Relaunch** the `PCR0_old` EIF; the Manager reconnects and recovers the same
   key set. The trigger is level-based, so rollback is just another
   `activeImage` mismatch the Manager resolves the same way as an upgrade.

---

## 6. Genesis image retirement

The **genesis** (first-install) image is built with the key-continuity guard
**off** (`EXPECT_EXISTING_KEYSET=false`) and has its own distinct PCR0. It
bootstraps the key set exactly once per environment. After the first upgrade
brings up an upgrade-variant image, retire the genesis PCR0 like any old image:
`removePcr0(genesis)` on-chain (once it is no longer `activeImage`) **and** strip
it from the KMS policy. Leaving it accepted would let a guard-off image run,
re-opening the silent-regeneration path Task 6 closes.

---

## 7. Emergency runbook (D4)

There is deliberately **no designed "stop now" verb**. An operator with host
access always has an out-of-band stop (`nitro-cli terminate-enclave` / stopping
the supervisor). That is an *ungraceful* stop — in-flight requests are lost — but
no **state** is lost, because the key set and encrypted state are recovered on
the next start.

### 7a. The newly-swapped image is found vulnerable

Fastest path — roll back to the prior image, then retire the bad one:

1. Roll back (§5): `proposePcr0Swap(PCR0_old)` → `applyPcr0Swap()`, relaunch the
   `PCR0_old` EIF.
2. Retire the bad image: `removePcr0(PCR0_bad)` (only after step 1 moved
   `activeImage` off it — else `CannotRemoveActiveImage`) **and** strip
   `PCR0_bad` from the KMS policy.
3. **Terminate any live bad enclave.** Stripping a PCR0 from the KMS policy does
   **not** evict an enclave that already recovered the keys and holds them in
   memory — the policy is only checked at recovery time. Always pair the policy
   strip with `nitro-cli terminate-enclave <enclave-id>` on every host that could
   be running `PCR0_bad`.

### 7b. All accepted images are vulnerable

No safe image to roll back to. You must wait out the timelock for a fixed image,
and downtime = `pcr0UpgradeDelay` **by design** (the timelock doubles as
worst-case emergency downtime — keep it modest, 24–48h).

1. **Full stop:** `nitro-cli terminate-enclave` on all hosts and **remove the
   entire `EnclaveKeyRecovery` statement** from the key policy so nothing can
   recover keys while a fix is prepared. Do **not** merely empty its
   `kms:RecipientAttestation:PCR0` value array — a condition with an empty value
   list is a malformed policy document that AWS rejects; delete (or `Deny`-neuter)
   the whole statement instead.
2. Build + publish the fixed image; `proposePcr0Swap(PCR0_fix)`; wait the
   timelock; `applyPcr0Swap()`; re-add the `EnclaveKeyRecovery` statement listing
   `PCR0_fix`; relaunch.

**Disclosure dynamic (why a modest delay).** `proposePcr0Swap(PCR0_fix)` starts a
**public** clock and, because builds are reproducible, an attacker can diff the
fix against the vulnerable release during the window. This argues for a modest
`pcr0UpgradeDelay` (24–48h), not a long one.

**Intake pause — not implemented.** An instant, guardian-triggerable intake
*pause* on request submission (to close the exposure window without a full stop,
and without freezing withdrawals as a terminate does) was considered (Task 0 /
D4) but is **not currently implemented** in `ProcessorEndpoint`. Until it is, the
only levers during 7b are the on-chain timelock and an out-of-band enclave
terminate. Note the withdrawal-freeze tension: while enclaves are terminated,
users cannot withdraw either (withdrawals need TEE signatures).

---

## 8. Pre/post-change checklist

Before **and** after any KMS policy or on-chain PCR0 change, verify the two trust
roots agree:

- [ ] Every KMS `kms:RecipientAttestation:PCR0` value hashes
      (`keccak256` of its 48-byte preimage) to an entry in on-chain
      `acceptedPcr0List`, and the counts match (see §3 — these lists are never
      byte-equal: KMS holds preimages, the chain holds hashes).
- [ ] Exactly one KMS value hashes to the on-chain `activeImage`.
- [ ] No stale genesis / retired PCR0 remains in either list.
- [ ] `PutKeyPolicy` was performed by `kms-policy-admin` (multi-party), and the
      CloudTrail/Config alarm fired to the out-of-band channel as expected.
- [ ] After a relaunch: admin `get_version` reports the expected `executorPcr0`,
      and no `SignerContinuityError` was raised.
