# Handshake Protocol between Manager and Executor

The handshake process is initiated by the Executor every time a new connection is established with the Manager. The primary goal of the handshake is to ensure the Executor has a valid and consistent set of cryptographic keys (`EnclaveKeySet`) to perform its duties.

There are two main scenarios for the handshake:

1.  **First Connection**: This occurs when the Manager connects to the Executor for the very first time, and no keyset recovery data exists in the Manager's persistent storage.
2.  **Reconnection**: This occurs on subsequent connections, where the Manager already holds recovery data from a previous session.

---

### Scenario 1: First Connection (No Keyset Data)

In this scenario, the Executor generates a new keyset and sends the corresponding recovery data to the Manager for storage.

```
Manager                                     Executor
   |                                            |
   |----------- TCP/VSock Connection ---------->|
   |                                            |
   |                                            | (Connection established)
   |<--- GetKeysetRecoveryRequest (ID: 1) <-----| (Executor initiates handshake)
   |                                            |
   | (Queries data layer, finds nothing)        |
   |                                            |
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: false)
   |                                            |
   |                                            | (Generates new keyset)
   |<--- SetKeysetRecoveryRequest (ID: 2) <-----| (Contains new recovery data)
   |                                            |
   | (Stores recovery data in data layer)       |
   |                                            |
   |---- SetKeysetRecoveryResponse (ID: 2) ---->|
   |                                            |
   o (Handshake complete)                       o (Handshake complete)
   |                                            |
```

---

### Scenario 2: Reconnection (Successful Keyset Recovery)

Here, the Manager provides the existing recovery data to the Executor, which then restores its keyset.

```
Manager                                     Executor
   |                                            |
   |----------- TCP/VSock Connection ---------->|
   |                                            |
   |                                            | (Connection established)
   |<--- GetKeysetRecoveryRequest (ID: 1) <-----| (Executor initiates handshake)
   |                                            |
   | (Queries data layer, finds recovery data)  |
   |                                            |
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: true, with recovery data)
   |                                            |
   |                                            | (Restores keyset from data)
   |<--- KeysetRecoveryResult (ID: 2) <---------| (Result: success)
   |                                            |
   o (Handshake complete)                       o (Handshake complete)
   |                                            |
```

---

### Scenario 3: Reconnection (Keyset Recovery Failure)

If the Executor fails to restore its keyset from the data provided by the Manager (e.g., due to data corruption), it notifies the Manager of the failure.

```
Manager                                     Executor
   |                                            |
   |----------- TCP/VSock Connection ---------->|
   |                                            |
   |                                            | (Connection established)
   |<--- GetKeysetRecoveryRequest (ID: 1) <-----|
   |                                            |
   | (Queries data layer, finds corrupted data) |
   |                                            |
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: true, with corrupted data)
   |                                            |
   |                                            | (Fails to restore keyset)
   |<--- KeysetRecoveryResult (ID: 2) <---------| (Result: failure with error message)
   |                                            |
   x (Handshake fails, Manager may terminate)   x (Handshake fails, Executor closes connection)
   |                                            |
```

---

### Executor self-identity field (PCR0)

In addition to the communication public key and signing key address, the Executor
reports its **running image identity** to the Manager in the two executor→manager
handshake-result messages:

- `SetKeysetRecoveryRequest` (Scenario 1, new keyset) — `SetKeysetRecoveryRequestData`
- `KeysetRecoveryResult` (Scenario 2 success, keyset restored) — `KeysetRecoveryResultData`

Both carry one additional field:

| Field  | Type   | Meaning |
|--------|--------|---------|
| `pcr0` | string | Hex-encoded PCR0 of the running enclave image, read once from the NSM (`DescribePCR(0)`) at Executor startup. Empty (`""`) in non-Nitro (TCP/dev) mode where no NSM device exists — the **dev marker**. |

The Manager stores the reported `pcr0` as the **running image** (`RunningPcr0()`). The
running image is what the swap-observation logic compares — as `keccak256(pcr0)` —
against the on-chain `activeImage` (see Task 5). The Manager also surfaces `pcr0` in the
aggregated admin `get_version` response (`executorPcr0`). The Executor's binary version is
obtained on demand via the forwarded `get_version` admin command, not carried in the
handshake.

The field uses `omitempty` and the wire format is newline-delimited JSON, so an older
peer that does not send it is read as an empty value (dev marker / `""`).

---

### Protocol version negotiation

Both sides advertise a monotonic wire-protocol version, `communication.WireProtocolVersion`
(currently `1`), which is deliberately distinct from `pkg/version.Version` (the
binary/release version) and from `common.Request.ProtocolVersion` (the on-chain
request-level version) — the `Wire` prefix disambiguates it from the latter. Bump it on
any incompatible change to the handshake or message set.

The exchange piggybacks on the first request/response pair, so an incompatible peer is
rejected **before any keyset-recovery data is read, generated, or stored**:

- The Executor puts its version in `GetKeysetRecoveryRequest` (`wireProtocolVersion`).
- The Manager checks `IsCompatible(peer)`; on mismatch it fails the handshake with a typed
  `IncompatibleProtocolError` and returns no recovery data. Otherwise it replies with its
  own version in `GetKeysetRecoveryResponse`.
- The Executor checks the Manager's version from that response; on mismatch it aborts with
  `IncompatibleProtocolError` before restoring or generating a keyset.

Both-side checks are required for rollout: an old peer (which does not send the field, so
it reads as version `0`) cannot perform its own check, so the up-to-date peer must reject
it. Current policy is exact-match (`IsCompatible` = equality); `0` (legacy/absent) is
incompatible with any non-zero version.

---

### Graceful shutdown (drain step 3)

The Manager can ask the Executor to exit cleanly via the `shutdown` admin command
(`admin.AdminCmdShutdown`), forwarded over the communication channel (`ForwardAdminCommand`).
It is **not** an operator-facing verb on the Manager admin server — it is used internally by
the swap-drain sequence (Task 5).

The Executor acks with `ShutdownResponse{stopping:true}` and then closes its
`ShutdownRequested()` channel; `cmd/executor/main.go` selects on that channel (alongside
OS signals) and performs the same graceful teardown as a `SIGTERM`. Because the Executor
disconnects immediately after acking, the Manager tolerates a missing ack / disconnect as a
normal outcome (`forwardShutdown`) rather than a fatal error.

---

### Key-continuity guard (Task 6)

To honour R2 (recover, never regenerate), two guards protect the keyset during a
software upgrade:

- **Executor** — when `EXECUTOR_EXPECT_EXISTING_KEYSET` is set, a
  `GetKeysetRecoveryResponse` with `DataFound: false` is fatal: `performHandshake`
  returns `ErrUnexpectedKeysetGeneration` **before** generating a keyset or sending
  `SetKeysetRecoveryRequest`, so Scenario 1 (new-keyset generation) is disabled.
  This is set during upgrades, where the Manager must already hold recovery data; a
  genuine first install runs without it. In Nitro deployments the flag is baked
  into the enclave image (`dockerfiles/executor/Dockerfile`, default on; genesis
  builds opt out via `EXPECT_EXISTING_KEYSET=false`) since an enclave receives no
  environment at launch — which also makes the guard's state part of PCR0 and
  attestable. Dev/TCP deployments set the env var directly.
- **Manager** — the recovery store refuses to overwrite an existing recovery blob
  with a *different* one (typed `recovery_data_exists` error, surfaced through the
  handshake failure path); re-storing an identical blob is a no-op. A wiped/wrong
  data folder therefore cannot silently clobber the only handle to the keyset.

The generation path (Scenario 1) also logs loudly (`Warn`) that it must only occur
on a genuine first install.

---

### Reconnect & swap observation (Task 5)

Each polling tick the Manager reconciles the connected enclave against the on-chain
`activeImage` before dispatching (level-based, so a Manager restart or missed event
converges to the same outcome):

- **Match** (`keccak256(reported pcr0) == activeImage`): dispatch normally. A dev/TCP
  executor reports an empty PCR0, so swap observation is skipped.
- **Mismatch**: enter drain — signal `shutdown`, close the channel, and reconnect. The
  handshake is no longer single-shot: the Manager re-dials with a fresh handshake
  (`reconnectExecutor`) and the client re-arms its connection (a fresh shutdown channel per
  `Connect`; a stale reader's exit never tears down a newer connection). Once the
  relaunched enclave's PCR0 matches `activeImage`, drain clears and dispatch resumes.
- **Dropped connection** (no swap): if the channel is down (`!IsConnected()`) even without a
  drain — an executor crash or restart on the *same* image — the Manager re-dials and
  re-handshakes on the same `reconnectExecutor` path before dispatching, then re-evaluates the
  reconnected enclave. This is the trigger for a plain connection loss: without it the Manager
  would retry into a dead channel every tick and the relaunched executor, waiting on a
  handshake that never arrives, could never recover its keyset. A failed re-dial (executor not
  back yet) just skips the tick and retries.
- **Signer continuity**: once per handshake, the Manager compares the reported
  `signingKeyAddr` against on-chain `getTeeSigner`. This runs in **every channel mode**
  (including dev/TCP) — it does not depend on PCR0, so the swap-observation skip does not
  skip it. It is deferred (dispatch proceeds) while `teeSigner` is the zero address
  (pre-`updateTee` bootstrap) or the executor reported no signing address; a transient
  `getTeeSigner` read error skips the tick and retries (it never dispatches unverified,
  matching the `activeImage` read); a mismatch on a set `teeSigner` is fatal (the enclave
  regenerated its keyset instead of recovering it).
