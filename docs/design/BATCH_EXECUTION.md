# Batch Execution: Increasing Request Throughput

## 1. Problem Statement

The system currently processes **one request per blockchain block** — across all applications. On a chain with 2-second block times, this means a maximum throughput of ~30 requests per minute, regardless of how fast the executor can run WASM, how many requests are queued, or how many independent applications have pending work.

This is a bottleneck for any application with sustained request volume. Users see growing queue latency even when the system has idle capacity — and since all applications share the same single-request cadence, load on one application delays every other application too.

## 2. Current Implementation: Why It's Limited

The per-block bottleneck comes from three sequential dependencies in the manager's poll loop:

```
pollBlockchain()                                        (manager.go:482)
  └─ processRequestFromChain()                          (manager.go:510)
       ├─ GetNextPendingRequest()   ← fetches ONE       (manager.go:532)
       ├─ processRequest()                              (manager.go:630)
       │    ├─ SendProcessRequest() ← vsock round-trip  (manager.go:764)
       │    ├─ storeStateToStorage()← DB write          (manager.go:789)
       │    └─ submitStateOnChain() ← blocks until mined (manager.go:795)
       │         └─ SubmitStateUpdate()                   (blockchain/client.go:434)
       │              └─ sendTxAndWaitMined()             (blockchain/client.go:317)
       │                   ├─ bind.Transact()  ← send tx
       │                   └─ bind.WaitMined() ← BLOCKS here
       └─ return (next tick fires after polling interval)
```

Three things make this sequential:

### 2.1. Single-request fetch, single global queue

`getNextPendingRequest()` on the contract returns only the queue head plus the corresponding application's state root. Even if 10 requests are pending, the manager sees one at a time.

The contract keeps **one global FIFO queue** (`_requestQueue`) shared by all applications, plus a priority queue for trigger-generated requests (`_triggerQueue`, see section 5). Requests from different applications interleave in submission order, and the contract enforces FIFO dequeuing from the head. State roots, by contrast, are already per-application (`applicationStateRoots[applicationId]`), so the states of different applications are fully independent — yet their requests are serialized through the same queue and the same one-per-poll cadence.

### 2.2. Blocking transaction submission

`sendTxAndWaitMined()` calls `bind.Transact` (sends the transaction) and then `bind.WaitMined` (polls for inclusion in a block). The entire call blocks until the transaction is mined — typically one full block time (~2 seconds). No other work happens during this wait.

### 2.3. Poll-driven single-pass loop

`pollBlockchain()` runs on a ticker. Each tick processes exactly one request. The next tick fires only after the current one completes (fetch + execute + mine). The effective cadence is `max(pollingInterval, executionTime + miningTime)`.

### 2.4. Redundant work per request

Each request also carries overhead that compounds across sequential processing:

- **Manager → Executor round-trip**: Every request sends the application's full encrypted state + WASM bytes over vsock/TCP. The executor decrypts the state, runs WASM, re-encrypts the new state, and sends it back. For consecutive requests to the same application, the executor decrypts what it just encrypted one round-trip ago.

- **State re-encryption cycle**: The executor returns the AES-encrypted state to the manager, the manager stores it, then sends it right back to the executor for the next request. The executor decrypts it again. This decrypt → encrypt → transfer → decrypt cycle repeats for every request.

The sequential chain looks like this for 3 requests. Nothing overlaps — each request waits for the previous one to be fully mined before starting:

```
time ──────────────────────────────────────────────────────────────────────────────────►

Request 1:  [vsock: send state+req]  decrypt ── WASM ── encrypt  [vsock: return]  store  [send tx]  [mine ~2s]
                                                                                                              |
Request 2:                                                                                                    [vsock: send state+req]  decrypt ── WASM ── encrypt  [vsock: return]  store  [send tx]  [mine ~2s]
                                                                                                                                                                                                              |
Request 3:                                                                                                                                                                                                    [vsock: send state+req] ...

Total for 3 requests: ~3 × (vsock round-trip + WASM + crypto + DB write + mine wait) ≈ 10-15s
```

### 2.5. Smart contract: one call per request

The `stateUpdate()` contract function processes a single request per call. Each call validates the queue head, verifies the signature, updates `applicationStateRoots[applicationId]`, processes refunds/withdrawals, emits events, invokes the application's trigger contract if one is registered, and dequeues the request.

For N requests, this means N separate transactions, each paying the base transaction cost (~21,000 gas) and each requiring its own nonce, gas estimation, and receipt tracking.

### 2.6. Summary of bottlenecks

| Bottleneck | Where | Impact |
|---|---|---|
| Single-request fetch | `GetNextPendingRequest()` | Sees 1 request per poll |
| Single global queue | `_requestQueue` in `ProcessorEndpoint.sol` | Independent applications serialized behind each other |
| Blocking mine wait | `bind.WaitMined()` in `sendTxAndWaitMined()` | ~2s idle wait per request |
| Redundant state crypto | Executor decrypt/encrypt per vsock round-trip | Wasted CPU + bandwidth |
| One vsock round-trip per request | `SendProcessRequest()` | Latency × N |
| One tx per request | `stateUpdate()` contract design | N × base gas cost |
| N nonce/receipt management | Sequential `bind.Transact` calls | Complexity + failure modes |

## 3. Considered Approaches

Both approaches share a common starting point: fetch up to `MaxBatchSize` pending requests **for a single application** via `GetPendingRequestsWithStateRoot(maxCount)` instead of one request at a time. The contract selects which application to serve and returns its requests together with its state root (see section 4). The fetch is capped at the source — there is no reason to transfer an entire pending queue only to discard most of it. The approaches differ in how execution and on-chain submission are organized.

A batch is always scoped to **one application**: one WASM module, one encrypted state, one state root chain. Requests from different applications are never mixed in the same batch.

> **Applications with a registered trigger contract are excluded from both approaches.** See section 5 for why, and for what happens to them instead.

### 3.1. Approach A — Pipeline Execution (individual transactions)

Execute requests one at a time through the executor, but submit each transaction to the mempool **without waiting for mining**. While the manager executes the next request, previously submitted transactions propagate and get mined in the background.

**Timeline for 3 requests:**

```
time ──────────────────────────────────────────────────────────────────────────────────────────►

               ┌── executor ──────────────┐  ┌── executor ──────────────┐  ┌── executor ──────────────┐
Request 1:     │ decrypt ── WASM ── encrypt│  │                          │  │                          │
               └──────────────────────────┘  │                          │  │                          │
                  store ── send tx ───────── │ ─── mining ──────────────│──│──────────── ✓            │
                                             │                          │  │                          │
Request 2:                                   │ decrypt ── WASM ── encrypt│  │                          │
                                             └──────────────────────────┘  │                          │
                                                store ── send tx ───────── │ ─── mining ──── ✓       │
                                                                           │                          │
Request 3:                                                                 │ decrypt ── WASM ── encrypt│
                                                                           └──────────────────────────┘
                                                                              store ── send tx ── mining ── ✓
                                                                                                         |
                                                                                            wait for last receipt
```

Execution overlaps with mining of previous transactions. The manager submits N separate `stateUpdate()` transactions with sequential nonces.

**How it works:**

1. Fetch up to `MaxBatchSize` pending requests + state root for the contract-selected application
2. For each request:
   - Send to executor (vsock round-trip, 1 per request)
   - Store updated encrypted state in DB (needed for rollback)
   - Submit `stateUpdate()` transaction — return immediately after `bind.Transact`, do not call `bind.WaitMined`
   - Chain `prevStateRoot` from the previous execution result
3. After the loop, wait for the last transaction receipt
4. Same sender + sequential nonces guarantees ordering: if the last tx is mined, all previous ones are too


**Pros:**

- No smart contract modifications required for the `stateUpdate()`
- Incremental change — refactors the manager's submission logic without touching the executor or contract
- Partial success is possible: if transaction K reverts, transactions 1..K-1 are already confirmed on-chain
- Each transaction is small, well within block gas limits
- Granular on-chain visibility — each request gets its own transaction hash and events

**Cons:**

- **Gas estimation broken for tx 2+.** `bind.Transact` calls `estimateGas`, which simulates against the latest confirmed state. For transaction 2, the on-chain `applicationStateRoots[applicationId]` hasn't been updated by the still-pending transaction 1, so `estimateGas` fails. Workaround: set `TransactOpts.GasLimit` to a fixed non-zero value to skip estimation. This requires tuning and risks over- or under-estimation.

- **Cascading failure.** The contract validates `prevStateRoot == applicationStateRoots[applicationId]` and `isCurrentPendingRequest()` on each call. If transaction K reverts, transactions K+1..N all revert too (wrong `prevStateRoot`, wrong queue head). Gas is wasted on all the reverted transactions.

- **Nonce management complexity.** N pending transactions require correct nonce sequencing. `PendingNonceAt` handles this when transactions are submitted sequentially from one goroutine, but edge cases exist (dropped transactions creating nonce gaps, node mempool limits).

- **N vsock round-trips.** Each request still requires a full manager → executor → manager cycle. The executor decrypts state it just encrypted on the previous iteration.

- **Intermediate state storage required.** The manager must store encrypted state after each execution for rollback support. If transactions 1-2 succeed on-chain but 3 fails, the manager needs the encrypted state at stateRoot₂ to recover. This means N DB writes and N encrypt cycles in the executor — the redundant crypto overhead from section 2.4 remains.

- **N × 21,000 base gas overhead.** Each transaction pays the base cost independently.

### 3.2. Approach B — Batch Execution (single aggregated transaction)

Send all requests to the executor in a single batch. The executor processes them sequentially while keeping the decrypted state in memory, avoiding redundant crypto and vsock round-trips. The manager submits one `batchStateUpdate()` transaction containing all results.

**Timeline for 3 requests:**

```
time ──────────────────────────────────────────────────────────────────────────────────────────►

             ┌── executor (single vsock round-trip) ─────────────────────────────────────┐
             │                                                                            │
             │  decrypt ── WASM₁ ── WASM₂ ── WASM₃ ── encrypt                            │
             │  (1 decrypt)  (state kept in memory)    (1 encrypt, final state only)      │
             │                                                                            │
             └────────────────────────────────────────────────────────────────────────────┘
                                                                                store ── send batchStateUpdate() ── [mine ~2s] ── ✓
```

**How it works:**

1. Fetch up to `MaxBatchSize` pending requests + state root for the contract-selected application
2. Send all N requests + encrypted state + WASM module to executor in a single batch message
3. Executor decrypts state **once**, runs WASM for each request sequentially (keeping plaintext state in memory between calls), encrypts **once** at the end
4. Executor returns: N `UpdatePayload`s (one per request, unsigned individually) + 1 batch signature over all entry hashes + 1 final encrypted state
5. Manager stores the final encrypted state in DB (1 write)
6. Manager submits one `batchStateUpdate()` transaction with all N update payloads + the batch signature
7. Wait for the single receipt

**Contract changes:** New `batchStateUpdate()` function that loops over entries, calling the same internal validation logic as current `stateUpdate()`. Since the contract is not in production, this is a refactor of the existing code — extract the body of `stateUpdate()` into an internal `_processOneStateUpdate()` and call it in a loop. Additional optimizations:

- **Single stateRoot storage write**: read `applicationStateRoots[applicationId]` once at the start, chain through entries in memory, write once at the end — saves (N-1) warm `SSTORE` operations (~5,000 gas each)
- **State root chain validation**: only the first entry checks `prevStateRoot` against storage; subsequent entries validate `entries[i].prevStateRoot == entries[i-1].newStateRoot` in memory
- **Deduplicated `applicationId`**: passed once instead of per-entry
- **Batch signature**: verify one signature over the batch digest instead of N individual `ecrecover` calls, saving (N-1) × ~3,000+ gas. Individual entry hashes are emitted in events for off-chain verifiability

**Batch signing scheme** (as implemented in `MsgToSignBuilder.BuildBatchMsgHash`, `pkg/executor/msgtosign_builder.go`). Each entry hash is `entryHash_i = keccak256(abi.encode(<entry fields>))` — the same per-entry hash the single-request path uses (`buildEntryHash`). The batch digest is the Ethereum `personal_sign` (EIP-191) hash of the **concatenated** entry hashes:

```
batchDigest = keccak256("\x19Ethereum Signed Message:\n" || itoa(32*N) || entryHash_0 || entryHash_1 || ... || entryHash_N-1)
```

There is **no intermediate `keccak256` over the concatenation** — the concatenated bytes are the `personal_sign` message itself. Two properties make this unambiguous:

- **Injectivity**: every entry hash is a fixed 32-byte `keccak256` output, so a concatenation of N of them splits exactly one way. This depends on the per-entry hash staying fixed-length.
- **Length binding**: the `personal_sign` prefix commits to the total message length (`32*N`), so batches of different sizes cannot collide.

A consequence is that a **1-entry batch digest is byte-identical to the single-request digest** (`BuildMsgHash` = `TextHash(entryHash)`), so both submission paths share one signing scheme and a 1-entry batch signature verifies on the single-request `stateUpdate()` path.

The contract must reconstruct exactly this digest: `personal_sign` over the concatenated entry hashes with a **dynamic** length prefix (`32*N`), not a fixed one, and without an extra hash layer.

**Pros:**

- **1 decrypt + 1 encrypt** instead of N of each. The executor holds plaintext state in memory between requests. For a batch of 5, this eliminates 4 decrypt cycles, 4 encrypt cycles, and the associated serialization/deserialization overhead.

- **1 vsock round-trip** instead of N. The full batch is sent and received in a single message exchange. No redundant state transfer.

- **Atomic on-chain submission.** All-or-nothing: either all N requests are processed or none are. No cascading failure, no partial state, no ambiguity about which transaction failed.

- **Gas estimation works normally.** A single transaction is simulated against the confirmed state. The first entry's `prevStateRoot` matches the on-chain root, and the contract processes subsequent entries internally — no pending-state simulation issues.

- **1 nonce, 1 receipt.** No nonce management complexity, no sequential nonce gaps, no receipt-per-transaction tracking.

- **Simple recovery.** If the transaction reverts, nothing changed on-chain. Roll back local state to pre-batch, retry the whole batch next poll. No two-poll recovery, no reorg detection needed for partial failures.

- **Gas savings.** Eliminates (N-1) × 21,000 base transaction cost. Additional savings from single stateRoot write and batch signature (1 `ecrecover` instead of N). For a batch of 5: ~84,000 gas saved from base costs alone, plus ~12,000 from reduced signature verification.

- **1 DB write** instead of N. Only the final encrypted state needs to be stored.

- **Subgraph-transparent.** The contract emits the same events per request inside the loop (`StateRootUpdate`, `UserEvent`, `RequestCompleted`, etc.). Indexers see identical events, just all in the same block. No subgraph changes needed.

**Cons:**

- **Contract modification required.** New `batchStateUpdate()` function and internal refactoring, plus the per-application queue change from section 4. Needs review and testing. However, since the contract is not in production, this is a development-time cost, not a migration risk.

- **Block gas limit risk.** A batch of N updates with events and withdrawals in a single transaction could consume significant gas. Estimated ~200,000 gas per request entry means a batch of 5 uses ~1M gas. On L2s (block gas limits typically 30M+) this is well within limits. A configurable `MaxBatchSize` cap mitigates this.

- **All-or-nothing risk.** If the transaction runs out of gas mid-loop or any entry triggers an unexpected revert, all N results are lost — including entries that would have succeeded individually. With pipeline execution, the first K transactions would have already confirmed. However: the TEE pre-validates every entry, so an on-chain revert should only happen due to external state changes (new requests submitted between poll and mining), which would also break pipeline execution via cascading failure.

- **Larger calldata.** All N entries are ABI-encoded in one transaction. Calldata cost: 16 gas per non-zero byte, 4 per zero byte. On L2s with EIP-4844 blob data, this is cheap. On L1, it would add up for large batches.

- **Executor becomes batch-aware.** The communication protocol needs a new batch message type. The executor must handle partial failures within a batch (e.g., request 3 fails WASM execution — produce an error payload for it, continue with request 4 using the state from request 2). This is a moderate protocol change.

- **No partial on-chain progress.** Unlike pipeline execution where transactions 1..K might already be confirmed when K+1 fails, batch execution is all-or-nothing. If the batch fails, the system retries everything next poll. In practice, this rarely matters because the TEE pre-validates all entries, but it means a failed batch wastes more wall-clock time than a failed pipeline.

### 3.3. Comparison

| Criterion | Pipeline (A) | Batch (B) |
|---|---|---|
| Contract changes | Per-app queues (section 4) | Per-app queues (section 4) + new `batchStateUpdate()` |
| Executor changes | None | New batch message type |
| Manager changes | Refactor submission logic | Refactor poll loop + submission |
| Vsock round-trips per batch | N | 1 |
| State decrypt/encrypt cycles | N each | 1 each |
| DB writes per batch | N | 1 |
| On-chain transactions | N | 1 |
| Gas (5 requests, estimated) | ~1,100,000 | ~1,020,000 |
| Nonce management | N pending nonces | 1 nonce |
| Gas estimation | Broken for tx 2+ (hardcode GasLimit) | Works normally |
| Failure mode | Cascading (K+1..N revert) | All-or-nothing |
| Recovery on failure | Two-poll (reorg detection) | Single-poll (retry batch) |
| Partial on-chain progress | Yes (tx 1..K confirmed) | No (all or nothing) |
| Block gas limit risk | Low (small individual txs) | Medium (capped by MaxBatchSize) |
| Implementation complexity | Lower | Moderate |

### 3.4. Decision

**Approach B (Batch Execution) is the chosen solution.**

The batch approach addresses every bottleneck from section 2 in a single, coherent design. The pipeline approach only solves the blocking mine wait (2.2) while leaving redundant crypto (2.4), per-request round-trips (2.4), and per-request gas costs (2.5) untouched — and introduces new problems (gas estimation workaround, cascading failure, nonce management) that batch avoids entirely.

The main cost — contract and executor changes — is a one-time development effort with no migration risk since the contract is not yet in production.

> **Note on combining approaches.** Pipeline and batch could theoretically be combined for cases where the pending queue exceeds `MaxBatchSize`: process multiple batches per poll cycle, submitting each batch transaction without waiting for mining before starting the next batch. This reintroduces pipeline's cons between batches (gas estimation workaround, cascading failure, nonce management) for a modest saving of ~2s per additional batch. This is not worth the added complexity for the initial implementation. Sequential batches (wait for mining between batches) are simpler and sufficient — each batch already processes up to `MaxBatchSize` requests, so the inter-batch mining wait is a small fraction of total throughput. Multi-batch pipelining can be revisited later if needed.

## 4. Multi-Application Batching: Required Contract Changes

Vela is multi-app: each application has its own state, WASM module, locked funds, and version chain. Part of what batching needs is already in place; the queue structure is not. This section separates what exists from what the batch design requires.

### 4.1. Already implemented: per-application state roots

The contract stores one state root per application:

```solidity
mapping(uint64 => bytes32) public applicationStateRoots;
```

`stateUpdate()` already receives `applicationId` and reads/writes `applicationStateRoots[applicationId]`. The manager's local storage, reorg detection, and version chains are likewise keyed by application. Nothing to do here — `batchStateUpdate()` chains `prevStateRoot → newStateRoot` across entries the same way `stateUpdate()` does, scoped to one application's root.

### 4.2. Required: per-application pending queues

The current single global FIFO queue (`_requestQueue`) cannot work with batches. A batch is scoped to one application, but the contract enforces FIFO dequeuing from the global head: if requests for applications A, B, and A are enqueued in that order, a batch for application A cannot skip the request for B.

The contract needs **per-application pending queues**. Each application maintains its own queue head/tail, and requests are enqueued into `pendingQueues[applicationId]`. `batchStateUpdate()` dequeues from the queue matching its `applicationId`. The trigger queue (`_triggerQueue`, see section 5) remains a single global priority queue.

Because request ids are unique across queues, the `PendingRequest` structs themselves live in **one global store** keyed by `requestId` (`mapping(bytes32 => PendingRequest)`); a queue only holds the ordered ids plus its head/tail. This keeps `requestById()` and `isCurrentPendingRequest(requestId)` working without an `applicationId` argument. A single `_totalQueuedRequests` counter, maintained on enqueue/dequeue, makes the aggregate size views O(1) so the `maxQueueSize` check on the submit path does not have to scan every application.

**Deploy requests use a third, global queue** (`_deployQueue`), not per-application queues. A deploy's application does not exist yet — it has no state root and is absent from `_deployedAppIds` until the deploy succeeds — so a request sitting in `pendingQueues[derivedAppId]` would be invisible to the round-robin scan of section 4.3 and would never be served. `_deployQueue` is served with precedence over the per-application queues (and after the trigger queue), always one request at a time, matching the manager's existing individual `processDeployApp()` path. Deploys are permissioned and capped by `availableDeploySlots`, so they cannot starve the rotation.

**`maxQueueSize` stays an aggregate cap**: `getPendingRequestsSize()` returns the sum of the deploy queue and every per-application queue (the trigger queue is reported separately by `getTriggerQueueSize()`), and the submit paths check that sum, exactly as before. One application can therefore still consume the global capacity; a per-application cap can be introduced later if that becomes a problem.

Side benefits:

- A permanently failing request for application A blocks A's queue at the structural level. **Note:** this does not by itself isolate other applications — the round-robin selection (section 4.3) is enforced on-chain, so the cursor cannot advance past the blocked application and the whole system stalls. See section 7.4.
- Per-application queues enable a deterministic application selection algorithm (section 4.3).

### 4.3. Required: fetching with round-robin application selection

`GetPendingRequestsWithStateRoot` does **not** receive an `applicationId` parameter. Instead, the contract selects the application to serve:

```
GetPendingRequestsWithStateRoot(maxCount uint64) (uint64, []*common.Request, [32]byte, error)
//                                                 ↑ applicationId (selected by contract)
```

The contract reuses the existing array of deployed applications (`_deployedAppIds`, `ProcessorEndpoint.sol:40`) and adds a single piece of new state: a round-robin cursor into it. No per-queue bookkeeping is needed — the array is append-only (deploys only ever add to it), so the enqueue and dequeue paths are untouched.

- **Selection**: starting at the cursor, scan `_deployedAppIds` (wrapping around) for the first application with a non-empty queue. The view returns up to `maxCount` of its requests, together with its state root and `applicationId`. If no application has pending work, the view returns an empty request list and the cursor stays put.
- **Enforcement**: any state update that dequeues from a per-application queue — `batchStateUpdate()` and `stateUpdate()` alike — recomputes the same scan, requires the submitted `applicationId` to be the scan's result, and sets the cursor just past it. Trigger-queue processing bypasses the cursor (see below). Selection is not a convention the manager follows — it is a rule the contract enforces.

Because the scan skips empty queues, round-robin is work-conserving: a single busy application gets every batch when nothing else is pending.

The scan costs one queue-emptiness storage read per skipped application, so enforcement is O(number of deployed apps) in the worst case. Deploys are permissioned (`DEPLOYAPP` role), so the array stays small and this is negligible; if the platform ever hosts hundreds of mostly-idle applications, revisit with an actively maintained non-empty-queue index.

Three exceptions to the plain round-robin rule:

- **Trigger queue precedence.** If the global trigger queue is non-empty, its head (a TRUSTPROCESS request) is returned alone, before any per-application selection — preserving the existing priority semantics (section 5.1). Processing it does not advance the cursor; rotation resumes where it paused.
- **Deploy queue precedence.** If the global deploy queue is non-empty, its head is returned alone, before any per-application selection (section 4.2). Processing it does not advance the cursor either.
- **Trigger applications capped at one.** If the selected application has a registered trigger contract, at most one request is returned regardless of `maxCount` (section 5.4). Processing it advances the cursor like any other turn.

**Why round-robin — and why the contract selects, not the manager:**

- **Anti-censorship, enforced.** If the manager chose which application to process next, a malicious or biased manager could starve specific applications by never selecting them. With the cursor check in the state-update functions, the selection is enforced on-chain — the manager cannot bypass the view and submit for an application out of turn.
- **Fairness for low-traffic applications.** Each application with pending work gets an equal share of batches. An application submitting one request is served after at most one batch per other active application — it never waits behind another application's entire backlog.
- **Simplicity.** No timestamp comparisons, no tie-breaking rules — the cursor is the whole algorithm.

> **Alternative considered — oldest-first.** Selecting the application whose queue head has the oldest block timestamp approximates global FIFO across applications. It was rejected for two reasons. First, enforcement cost: verifying "this application had the oldest head" inside a state-update transaction requires reading *every* queue head's timestamp on *every* call to find the minimum, whereas the round-robin scan reads only queue-emptiness flags and stops at the first non-empty queue. Second, fairness: global FIFO means a high-volume application imposes its entire queue latency on low-traffic applications. Cross-application submission order is not worth preserving — applications are fully independent, so nothing depends on it.

### 4.4. Manager: one batch per poll

The manager processes one batch per poll cycle. It does not choose which application to serve — the contract does (section 4.3). Each batch is self-contained: its own encrypted state, WASM module, state root, and on-chain transaction. The contract's round-robin selection gives every application with pending work an equal share of batches.

**Sequential vs parallel batches:** The initial implementation processes one batch per poll cycle — the simplest approach with no additional complexity. If throughput across many active applications becomes a bottleneck, the manager could process multiple batches per poll cycle by calling `GetPendingRequestsWithStateRoot` repeatedly and submitting batches without waiting for mining between them. Since different applications use independent state roots, there are no nonce-ordering or state-chaining conflicts between batches for different applications. This is the multi-application analog of the "combining approaches" note in section 3.4 and can be deferred.

### 4.5. Executor: no multi-app specific changes

The executor already processes work scoped to a single application — one WASM module, one encrypted state. The batch message (section 6, Stage 2) carries all of this per-batch; the executor is unaware of whether other applications exist. The only executor changes are the batch protocol itself, not multi-app handling.

### 4.6. Summary

| Component | Change | Status |
|---|---|---|
| Contract: state storage | `mapping(uint64 => bytes32) applicationStateRoots` | Already implemented |
| Contract: queue | Single global queue → `mapping(uint64 => Queue) pendingQueues` + global deploy queue (trigger queue stays global) | Implemented |
| Contract: `batchStateUpdate()` | New function; reads/writes `applicationStateRoots[applicationId]`; requires `applicationId` to match the round-robin scan | Required |
| Contract: round-robin tracking | Cursor into the existing `_deployedAppIds` array; scan skips empty queues; cursor advanced on every per-application dequeue | Implemented |
| Contract: cursor enforcement | `require(applicationId == scan result)` in the state-update functions | Required (deferred to `batchStateUpdate()`) |
| Contract: view functions | `getPendingRequestsWithStateRoot(maxCount)` serving the scan result; trigger and deploy queue precedence | Implemented |
| Manager: poll loop | One batch per poll cycle; contract selects the application | Required |
| Manager: state storage | Keyed by `applicationId` | Already implemented |
| Executor | Batch protocol only (section 6, Stage 2) | Required |

## 5. Trigger Applications (TRUSTPROCESS): Not Supported by This Design

Neither approach from section 3 works for applications that have a **trigger contract** registered. This section explains why, lists the alternatives considered, and records the decision.

### 5.1. How the trigger flow works

An application can register a trigger contract at deploy time (`triggerContracts[applicationId]`, `ProcessorEndpoint.sol:79`). The flow:

1. A processed request emits AppEvents. During the `stateUpdate()` transaction, the contract invokes the registered trigger with the AppEvent data (`_invokeTrigger`, `ProcessorEndpoint.sol:1102`).
2. The trigger contract acts on the events (e.g., receives unshielded funds, executes an action, re-shields the remainder) and derives a trusted payload **on-chain**.
3. If the trigger returns a payload, the contract enqueues a **TRUSTPROCESS** request into a dedicated priority queue (`_triggerQueue`). This is the only way a TRUSTPROCESS can be created.
4. The selection view (`getPendingRequestsWithStateRoot`, `_selectPendingRequests`) serves the trigger queue **before** the deploy and per-application queues, so the TRUSTPROCESS is processed immediately after the request that fired it — before any other pending request, for any application.

Trigger applications depend on this ordering. In the unshield/re-shield round trip, for example, request K sends funds out to the trigger and the subsequent TRUSTPROCESS credits the returned funds back into the application state. Any request processed between K and its TRUSTPROCESS would observe an intermediate state where funds have left the application but the trusted callback has not yet landed.

### 5.2. Why neither pipeline nor batch works

Both approaches rest on the same premise: the set of requests to process is known at fetch time, and executing them back-to-back is equivalent to executing them one transaction at a time. Trigger applications break this premise twice:

1. **TRUSTPROCESS requests do not exist at fetch time.** They are created on-chain *during mining* of the state update itself, and their payload is derived on-chain by the trigger contract from the AppEvent data in that same transaction. The TEE cannot know the payload when the batch is built, so a batch snapshot can never include the TRUSTPROCESS requests that its own entries will generate.

2. **The ordering guarantee is silently violated.** With batch execution, if entry K fires the trigger, the TRUSTPROCESS is enqueued mid-transaction — but entries K+1..N, computed in the TEE without the trusted callback, are still processed in the same transaction. With pipeline execution the same happens across transactions: tx K mines and enqueues the TRUSTPROCESS while tx K+1 is already in the mempool. In both cases the contract does **not** revert: `isCurrentPendingRequest()` accepts the head of either queue, and the application's state root chain remains internally consistent. The failure is purely semantic — the application state diverges from what the trigger flow requires — which makes it worse than a revert, because nothing on-chain signals that anything went wrong.

### 5.3. Alternatives considered

1. **Cut the batch at the first trigger fire.** For trigger applications, the executor stops the batch after the first request that emits AppEvents; the manager submits the partial batch; the TRUSTPROCESS enqueued by that last entry is served on the next poll with its priority intact. This preserves ordering and retains batching benefits for stretches of requests that emit no events. Costs: the executor must be told the application has a trigger, and the effective batch size collapses toward 1 for applications that emit events on most requests — which trigger applications typically do, since firing the trigger is their purpose.

2. **Relaxed TRUSTPROCESS ordering.** Let TRUSTPROCESS requests accumulate during the batch and process them afterwards. This is only sound if the application tolerates other requests observing state before the trusted callback lands. That is an application-specific property the platform cannot assume; it would require an explicit opt-in in the deploy descriptor and application authors designing their state to be safe under deferred callbacks.

3. **Single-request processing for trigger applications.** Detect a registered trigger (`triggerContracts[applicationId] != address(0)`) and serve those applications through the existing single-request `stateUpdate()` submission path, unbatched. Fetching still goes through `GetPendingRequestsWithStateRoot`, which returns at most one request for such applications (section 4.3).

### 5.4. Decision

**For now, applications with a registered trigger contract do not support batch requests** (alternative 3).

- The contract's selection logic returns at most one request for a trigger application, and the global trigger queue always takes precedence (section 4.3). TRUSTPROCESS requests are themselves always processed individually — their `stateUpdate()` also runs `_invokeTrigger`, so they have the same problem as the requests that fire them.
- As defense in depth, `batchStateUpdate()` reverts if the application has a registered trigger — the no-batching rule is enforced on-chain, not just by manager convention.
- With per-application queues (section 4.2), single-request processing of a trigger application does not delay batching of other applications.
- The manager does not detect trigger applications itself: the contract's selection caps them at one request, so the manager simply routes any single-request fetch through the single-request path (section 6, Stage 4). The trigger rule lives on-chain only — in the selection view and in the `batchStateUpdate()` revert.

**Why `stateUpdate()` is kept alongside `batchStateUpdate()`:** a 1-entry batch is semantically equivalent to `stateUpdate()`, so the two paths could in principle be unified. Keeping the single-request function is a deliberate risk/sequencing choice:

- **Near-zero cost.** After the Stage 1 refactor, all real logic lives in `_processOneStateUpdate()`; `stateUpdate()` is a thin wrapper calling it once. There are two entry points, not two implementations.
- **Refactor validation.** All existing contract tests run unchanged against the wrapper (Stage 1, step 3) — the proof that extracting `_processOneStateUpdate()` did not alter behavior.
- **New code stays away from the fragile flow.** Trigger-flow mistakes do not revert — they silently violate application semantics (section 5.2). Routing trigger applications and TRUSTPROCESS through the proven `stateUpdate()` path confines the new batch code (batch signature, entry loop) to flows where failures are loud.
- **Trigger queue handling.** `batchStateUpdate()` dequeues from `pendingQueues[applicationId]` only; TRUSTPROCESS requests live in the global `_triggerQueue`. Unifying would require teaching the batch function about the second queue and TRUSTPROCESS fee semantics (`maxFeeValue = 0`) from day one.
- **Crisp guard.** "Trigger application → batch path reverts" is a simpler invariant to audit than "batch allowed but only with one entry, from the right queue".

**Intended evolution:** once the batch path is proven, the trigger-app revert can be relaxed to `require(entries.length == 1)`, `batchStateUpdate()` can accept trigger-queue heads, and `stateUpdate()` can be retired — one submission path, one signature scheme. Alternatives 1 and 2 can likewise be revisited if trigger-application throughput becomes a bottleneck.

## 6. Implementation Plan

The work is organized in four stages. Each stage produces a testable, reviewable unit. Stages 1-3 can proceed in parallel across different developers.

### Stage 1 — Smart Contract: per-app queues + `batchStateUpdate()`

Refactor `ProcessorEndpoint.sol` to support per-application queues and batch submission.

**Steps:**

1. Replace the single global `_requestQueue` with `mapping(uint64 => RequestQueue) pendingQueues` plus the global `_deployQueue`; move the `PendingRequest` structs into one global `_requests` store keyed by `requestId`. `_triggerQueue` remains a single global priority queue. Add the round-robin cursor into `_deployedAppIds` and the shared scan helper (first app with a non-empty queue, starting at the cursor, wrapping — section 4.3), advancing the cursor on every per-application dequeue. Update `getNextPendingRequest()`, `isCurrentPendingRequest()`, `getPendingRequests*()` and `_resetQueue()` accordingly; `maxQueueSize` stays an aggregate cap (section 4.2). **Done.**

2. Extract the body of `stateUpdate()` into an internal function `_processOneStateUpdate()` that takes the same parameters and performs all validation, state updates, event emission, refunds, withdrawals, trigger invocation, and request dequeuing.

3. Rewrite `stateUpdate()` as a thin wrapper that calls `_processOneStateUpdate()` once. This preserves backward compatibility and confirms the refactor is correct — all existing contract tests must still pass without modification.

4. Add `batchStateUpdate(uint64 applicationId, BatchEntry[] calldata entries, bytes calldata batchSignature)` that:
   - Reverts if `triggerContracts[applicationId]` is set (section 5.4)
   - Requires `applicationId` to match the round-robin scan result and sets the cursor just past it (section 4.3); the same check applies to `stateUpdate()` when dequeuing from a per-application queue (trigger-queue processing bypasses it)
   - Verifies the batch signature: recover the signer from the batch digest (section 3.2 — `personal_sign` over the concatenated entry hashes, dynamic `32*N` length prefix, no extra hash layer) and `batchSignature`, verify it matches the registered TEE signer. One `ecrecover` call for the entire batch. This requires a new verification function on the TEE authenticator (e.g., `checkBatchSignature(bytes32[] entryHashes, bytes signature)` in `ITeeAuthenticator` / `AbstractTeeAuthenticator`) — the existing `checkSignature()` hashes a single `SignatureParams` struct and cannot verify a batch message. Note the digest is built from the entry hashes only, so the authenticator can take `bytes32[] entryHashes` directly; it must build the prefix from the array length at runtime.
   - Reads `applicationStateRoots[applicationId]` from storage once into a local variable
   - Loops over entries, calling `_processOneStateUpdate()` for each (signature verification is already done — `_processOneStateUpdate` skips per-entry `ecrecover`), dequeuing from `pendingQueues[applicationId]`
   - Validates state root chaining: first entry checks `prevStateRoot` against storage; subsequent entries check `entries[i].prevStateRoot == entries[i-1].newStateRoot`
   - Writes `applicationStateRoots[applicationId]` to storage once at the end of the loop (not per iteration)
   - Emits individual entry hashes in events for off-chain verifiability

5. Add the `getPendingRequestsWithStateRoot(maxCount)` view serving the round-robin scan result, with trigger- and deploy-queue precedence and the trigger-app cap (section 4.3). **Done.**

   `getNextPendingRequest()` is **removed** rather than kept as a `maxCount = 1` wrapper: every caller can pass `maxCount = 1` instead, and the wrapper's separate return shape (`(request, stateRoot, success)` with an empty-struct sentinel) only existed to preserve the pre-batch ABI. Removing it deletes a struct-returning external view — the most expensive kind of code in this contract (section "Open blocker"). Callers updated: `BlockChainClient.GetNextPendingRequest` and the `Client` interface method are gone; `GetPendingRequestsWithStateRoot` now calls the contract view directly instead of delegating to the removed one, and the hardhat tests call `getPendingRequestsWithStateRoot(1)`. **Done.**

6. Define the `BatchEntry` struct containing per-request fields: `prevStateRoot`, `newStateRoot`, `processedRequestId`, `events`, `eventSubTypes`, `withdrawalRequests`, `refund`, `applicationFees`, `errorCode`, `errorMsg`.

7. Write contract tests:
   - Batch of N successful requests with valid batch signature
   - Batch with an error payload mid-batch (request K fails, K+1 continues from unchanged state)
   - Batch with first entry having wrong `prevStateRoot` (reverts)
   - Batch with broken state root chain between entries (reverts)
   - Single-entry batch (equivalent to `stateUpdate()`)
   - Invalid batch signature (reverts)
   - Batch signature signed by wrong key (reverts)
   - `batchStateUpdate()` for an application with a registered trigger (reverts)
   - Mixed-app enqueue: requests for A, B, A — selection returns both A requests; B's queue untouched *(done)*
   - Round-robin rotation: batches for A, B alternate while both have pending work; cursor skips an application whose queue empties *(done)*
   - Cursor enforcement: `batchStateUpdate()` (and `stateUpdate()` on a per-application queue) for an application other than the round-robin scan result reverts
   - Scan correctness: applications with empty queues are skipped; scan wraps past the end of `_deployedAppIds`; all queues empty → view returns no requests and the cursor is unchanged *(done)*
   - Trigger queue precedence: pending TRUSTPROCESS returned alone before any batch selection *(done)*
   - Deploy queue precedence: pending deploy returned alone before any application selection *(done)*
   - Trigger application selected: at most one request returned regardless of `maxCount` *(done)*
   - Cross-application queue independence: processing one application's head leaves the others' queues and state roots untouched; `adminReset` drains every application queue and refunds its deposits *(done)*
   - Gas measurement: compare `batchStateUpdate(N entries)` vs N × `stateUpdate()`

**Contract size — see `PROCESSOR_ENDPOINT_SPLIT.md`.** The per-application queue work in step 1 pushed `ProcessorEndpoint` past the 24,576-byte EIP-170 limit. Making room is a self-contained refactor with no batch concepts in it, so it was done on its own branch and is documented separately: `ProcessorEndpointStorage` holds all state, `ProcessorEndpointExtension` hosts the facilitator path, the deploy-submission entry points, the operator resets and the admin setters behind a `delegatecall`, EIP-170 is now enforced by the hardhat config, and `npm run check:layout` guards the shared storage layout in CI. The entry points hosted in the extension operate on this design's per-application queue state (`RequestQueues.Store`, reached through the shared `_q`), so `submitDeployRequest*` enqueues into the global deploy queue and the resets drain every per-application queue.

What matters for this design is the budget that leaves:

| | Deployed bytes | vs limit |
|---|---|---|
| After the first split (facilitator path only), before the queue work | 21,609 | −2,967 |
| With per-app queues, deploy queue, round-robin and the selection view | 23,990 | −586 |
| Same, after the split's second pass moved deploy submission, the resets and the admin setters out | **19,475** | **−5,101** |

The −586 figure was **not** enough for `batchStateUpdate()` (step 4): with `viaIR` + `runs: 0`, the `BatchEntry[]` calldata decoder — nested dynamic arrays plus `string errorMsg` — is the expensive part, plausibly 1.5–2.5KB on its own. The room came from the split's second pass rather than from moving the queue views: **5,101 bytes of headroom**, with `ProcessorEndpointExtension` at 12,172. If `batchStateUpdate()` does not fit in that, the next levers are the read-only surface behind a generic `fallback()` (−2,353) or implementing `batchStateUpdate()` in the extension from the start — both in `PROCESSOR_ENDPOINT_SPLIT.md` section 4. The first two rows above are the hardhat `paris` path and the third is hardhat `cancun` (`runs: 0` throughout), because the split switched the EVM target (`PROCESSOR_ENDPOINT_SPLIT.md` section 2.1); the `go:generate` solc path uses `runs: 200` and no explicit EVM target, so it reports different numbers — always say which produced a figure.

Two removals of redundant external surface also helped, and belong to this design rather than the split:

- `getNextPendingRequest()` (−332 bytes), replaced by `getPendingRequestsWithStateRoot(1)` — see step 5.
- Candidates for the next pass: `getPendingRequests()` and `getPendingRequestsPage()` (−1,178 together), if the subgraph and `pkg/blockchain` turn out not to need them.

**Files changed:**
- `contracts/contracts/ProcessorEndpoint.sol`
- `contracts/contracts/ProcessorEndpointStorage.sol` (per-application queue state replaces the single global queue)
- `contracts/contracts/ProcessorEndpointExtension.sol` (the entry points it hosts — deploy submission and the operator resets — move onto the per-application queues)
- `contracts/contracts/RequestQueues.sol` (new — queue library)
- `contracts/contracts/interfaces/IProcessorEndpoint.sol`
- `contracts/contracts/AbstractTeeAuthenticator.sol` (new batch signature verification)
- `contracts/contracts/interfaces/ITeeAuthenticator.sol`
- `contracts/contracts/Structs.sol` (new `BatchEntry` struct)
- `contracts/test/` (new and updated test files)

The storage base, the extension, `IProcessorEndpointState` and the two guardrails come from the
contract split (`PROCESSOR_ENDPOINT_SPLIT.md`), which this work builds on.

### Stage 2 — Executor: Batch Processing

Add a batch message type to the executor so it can process multiple requests in a single vsock round-trip, keeping decrypted state in memory between requests. A batch is always scoped to one application.

**Steps:**

1. Define the batch request/response message types in the communication protocol:
   - `BatchProcessRequestMessage`: carries `[]*common.Request`, `*common.ApplicationState`, `[]byte` (WASM module)
   - `BatchProcessResponseMessage`: carries `[]*common.UpdatePayload` (unsigned individually), `[]byte` (single batch signature), `*common.ApplicationState` (final only), `[]*common.DeanonymizationReport`

2. Implement `HandleBatchProcessRequest()` in the executor. The batch loop must distinguish between soft failures (signed error payload — continue) and hard failures (bare error — stop batch). See section 7 for the full pseudocode and rationale. Key invariants:
   - `appData` is only mutated after successful WASM execution
   - On soft failure: error payload included in results, state unchanged, batch continues
   - On hard failure: batch stops, results for previously processed requests are returned
   - One payload is returned per handled request, in input order, so `len(payloads)` is how many of the input requests were consumed

3. Add the message handler in `communication/server.go` to route `BatchProcessRequestMessage` to the executor's `HandleBatchProcessRequest`.

4. Add `SendBatchProcessRequest()` to the executor client interface and implement it in the communication client.

5. Write executor tests:
   - Batch of N successful requests — verify state roots chain correctly
   - Soft failure mid-batch — request K gets error payload, K+1..N continue from K-1's state
   - Hard failure mid-batch (e.g., wrong applicationId on request K) — batch stops, only 1..K-1 returned
   - Request with deposit + process where process fails — verify deposit changes are discarded
   - Single-request batch — equivalent to existing `SendProcessRequest`
   - Verify only 1 AES decrypt and 1 AES encrypt occur for the batch
   - Verify a partial batch returns fewer payloads than input requests, so the shortfall is visible to the manager
   - Verify a single batch signature is returned (not per-entry signatures)
   - Verify batch signature covers all entry hashes: recover signer from the batch digest (section 3.2) and confirm it matches the executor's TEE key. Pin the digest against an independently computed expected value — recomputing it with `BuildBatchMsgHash` alone would pass under any scheme and cannot detect a change to it
   - Verify a 1-entry batch digest equals `BuildMsgHash` of that entry (the two paths must share one scheme)
   - Batch signing failure — verify entire batch is discarded, error returned, no partial results

**Files changed:**
- `pkg/communication/message.go` (new message types)
- `pkg/communication/server.go` (new handler)
- `pkg/communication/client.go` (new `SendBatchProcessRequest`)
- `pkg/executor/executor.go` (new `HandleBatchProcessRequest`)
- `pkg/blockchain/interface.go` (executor client interface update)
- `pkg/blockchain/mock_client.go` (mock implementation)

**Deploy requests:** Deploy is a different flow (creates initial state, stores WASM bytecode) and never enters a batch: with per-application queues, a pending deploy is always alone in its application's queue (section 4.2). The manager processes it individually via the existing `processDeployApp()` path.

### Stage 3 — Go Contract Bindings

Regenerate the Go bindings after the contract changes so the manager can call `batchStateUpdate()`.

**Steps:**

1. Regenerate bindings: `go generate ./...`
2. Add `SubmitBatchStateUpdate()` to `BlockChainClient` that:
   - Takes `[]*common.UpdatePayload` (the batch results) and `[]byte` (the batch signature)
   - Packs all entries into `BatchEntry[]` calldata
   - Passes the single batch signature as `bytes`
   - Calls `sendTxAndWaitMined()` (existing method — single tx, blocking wait is fine here)
3. Add `SubmitBatchStateUpdate` to the `Client` interface in `interface.go`
4. Implement the mock in `mock_client.go`

**Files changed:**
- `pkg/blockchain/bindings/` (regenerated)
- `pkg/blockchain/client.go`
- `pkg/blockchain/interface.go`
- `pkg/blockchain/mock_client.go`

### Stage 4 — Manager: Batch Orchestration

Refactor the manager's poll loop to fetch multiple requests and route them through the batch executor and batch contract submission.

**Steps:**

1. Add `MaxBatchSize` configuration to `config.go` with `MAX_BATCH_SIZE` env var (default 5). Add startup validation: `MaxBatchSize > 0`.

2. Add `GetPendingRequestsWithStateRoot(maxCount uint64)` to the blockchain client, returning `(uint64, []*common.Request, [32]byte, error)` — the contract-selected `applicationId`, up to `maxCount` of its pending requests, and its on-chain state root (section 4.3). The caller passes `MaxBatchSize` so only the requests that will actually be processed are fetched. The existing `GetPendingRequests()` only returns requests without the state root or selection logic. **Done** — it calls the contract view directly (it was previously a stub delegating to the now-removed `GetNextPendingRequest`).

3. Implement `processBatchFromChain()` in the manager:
   ```
   applicationId, requests, onChainStateRoot := GetPendingRequestsWithStateRoot(MaxBatchSize)

   if len(requests) == 0:
       return  // nothing to do

   verify localStateRoot[applicationId] == onChainStateRoot (existing per-app reorg logic)

   if requests[0] is a deploy request:
       // always alone in its application's queue (section 4.2)
       processDeployApp(requests[0])
       return

   if len(requests) == 1:
       // Covers TRUSTPROCESS and trigger applications (the contract returns at
       // most one request for them — sections 4.3, 5.4), as well as normal
       // applications with a single pending request. For one request,
       // stateUpdate() and a 1-entry batch are equivalent, so the manager
       // dispatches on request count alone — it needs no trigger knowledge
       // and never queries triggerContracts.
       processRequest(requests[0])  // existing single-request path
       return

   results, batchSignature, finalState := executor.SendBatchProcessRequest(
       requests, encryptedState[applicationId], wasmBytes[applicationId])

   if len(results) == 0:
       // Hard failure on the very first request or batch signing failure — nothing to submit
       log warning, retry next poll

   if len(results) > 0:
       save deanonymization reports to disk (if any)
       store final encrypted state in DB (1 write, versionID = final stateRoot)
       submit batchStateUpdate(applicationId, results, batchSignature) on chain
       on tx failure: rollback DB to pre-batch state for applicationId, retry next poll

   if len(results) < len(requests):
       log that request [len(results)] caused a hard stop
       // remaining requests stay in the application's on-chain queue for the next poll
   ```

4. Update `pollBlockchain()` to call `processBatchFromChain()` instead of `processRequestFromChain()`.

5. Write manager tests:
   - Happy path: N requests batched, single tx confirmed
   - Deploy request: handled individually via existing path
   - Single returned request (TRUSTPROCESS, trigger application, or normal application with one pending request): routed through the single-request path, never batched — no `triggerContracts` lookup performed
   - Batch tx reverts: verify DB rollback to pre-batch state
   - Empty queue: no-op
   - Queue larger than MaxBatchSize: only first MaxBatchSize requests fetched and processed
   - Deanonymization reports saved correctly within batch
   - Per-app reorg detection still works with batch state storage
   - Hard failure mid-batch: executor returns partial results, manager submits only processed requests, remaining stay pending
   - Hard failure on first request: nothing submitted, retry next poll

6. Integration test: submit 5 requests for one application on-chain, verify all processed in one poll cycle via a single `batchStateUpdate()` transaction; interleave requests for a second application and verify round-robin selection alternates between both across polls.

**Files changed:**
- `pkg/manager/config.go`
- `pkg/manager/manager.go`
- `pkg/blockchain/client.go` (`GetPendingRequestsWithStateRoot`)
- `pkg/blockchain/interface.go`
- `pkg/blockchain/mock_client.go`

### Implementation order

```
         ┌─────────────────┐
         │  Stage 1         │
         │  Contract         │──────────┐
         └─────────────────┘           │
                                        ▼
         ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
         │  Stage 2         │    │  Stage 3         │    │  Stage 4         │
         │  Executor         │    │  Go bindings     │──► │  Manager         │
         └────────┬─────────┘    └─────────────────┘    └─────────────────┘
                  │                                              ▲
                  └──────────────────────────────────────────────┘
```

- **Stage 1** (contract) and **Stage 2** (executor) can proceed in parallel — they have no code dependency on each other.
- **Stage 3** (bindings) depends on Stage 1 (needs the compiled contract ABI).
- **Stage 4** (manager) depends on Stages 2 and 3 (needs both the executor batch API and the contract bindings).

## 7. Failure Handling Within the Executor Batch

The current executor has two distinct error types, and the batch must handle them differently. (The classification below matches the existing single-request rules: signed execution errors go on-chain; transient/system errors are plain Go errors retried on the next poll.)

### 7.1. Soft failure — error payload (request dequeued)

The executor produces an `UpdatePayload` with `prevStateRoot == newStateRoot` (state unchanged) and a non-zero `ErrorCode`. The error payload is not individually signed — it is covered by the batch signature alongside all other entries. The contract marks the request as `FAILED`, collects the minimum fee, refunds the rest, and advances the application's queue head.

This happens for application-level errors where the executor has a valid stateRoot and can produce an error payload (see `HandleProcessRequest` in `pkg/executor/executor.go`):
- WASM execution failure — deposit or process
- Insufficient fuel
- Payload decryption failure
- Deanonymize report validation failure
- AppData serialization failure
- Event/report encryption failure (e.g., recipient key not registered)

**Batch behavior: continue.** The error payload is included in the batch results. The state is unchanged, so the next request's `prevStateRoot` chains correctly from the same stateRoot. No special handling needed.

```
Request 1: success   → state₁
Request 2: success   → state₂
Request 3: soft fail → error payload (prevStateRoot = newStateRoot = stateRoot₂)
Request 4: success   → state₃ (continues from state₂)
Request 5: success   → state₄
```

WASM execution is functional — state is passed as input and new state is returned as output. The original input is never modified. In the executor, `appData` is only mutated after confirmed success:

```go
appState := currentAppData.GetAppState()          // read (not consumed)

newState, ..., failure := runtime.ProcessRequest(  // WASM returns NEW state
    ..., appState, ...)                            // appState is untouched

if failure != nil {
    // appData still holds previous successful state
    // build error payload, continue to next request
    continue
}

// Only after success:
currentAppData.SetAppState(newState)
currentAppData.IncrementNonce()
```

No copy, no rollback, no special recovery logic.

### 7.2. Hard failure — bare error (request stays pending)

The executor returns a plain error — no signed payload. The manager cannot submit anything on-chain for this request. The request remains in the pending queue.

This happens for system-level errors, or for fields already validated on-chain whose unexpected value at the executor is evidence of tampering between the chain and the executor:
- `validateRequest()` failure: wrong `applicationId`, wrong `protocolVersion`, fee below minimum (all validated on-chain)
- App state not found — app existence is validated on-chain by the `validApplicationId` modifier in `ProcessorEndpoint`, so a nil state at the executor means tampering or manager-side state loss
- Unsupported request type
- State decryption failure
- AES state encryption failure

> **Note:** Signing failure is no longer a per-request hard failure. With batch signature, signing happens once after the batch loop completes. If batch signing fails, the entire batch is discarded — the executor returns an error with no results.

**Batch behavior: stop.** The executor cannot produce a result for this request, and the contract requires FIFO processing within the application's queue — request K cannot be skipped to process K+1. The batch stops at request K. Results for requests 1..K-1 are returned; requests K..N are not processed and remain pending.

```
Request 1: success    → state₁     ← included in batch results
Request 2: success    → state₂     ← included in batch results
Request 3: hard fail  → bare error ← batch stops here
Request 4: not executed             ← remains pending
Request 5: not executed             ← remains pending
```

The manager submits a `batchStateUpdate()` with only the results for requests 1-2. Requests 3-5 remain in the application's on-chain queue and will be retried on the next poll.

> **Note:** Some hard failures are transient (encryption failure — likely a system issue that will resolve). Others are permanent for this executor (tampered `applicationId`). In the permanent case, the request blocks the queue head — every subsequent poll will stop at the same request. This is the same behavior as today's single-request processing: the manager retries and fails each poll. Despite per-application queues (section 4.2), the blockage is **not** confined to the affected application: the enforced round-robin cursor cannot advance past the blocked application, so all other applications stall too. See section 7.4.

### 7.3. Executor batch pseudocode

```
HandleBatchProcessRequest(requests, encryptedState, wasmModule):
    appData := decryptState(encryptedState)            // 1 decrypt
    currentStateRoot := inputStateRoot
    results := []

    for i, req := range requests:
        // Pre-validation — can cause hard failure
        if err := validateRequest(req); err != nil:
            break  // stop batch, return results so far

        // WASM execution — can cause soft failure
        appState := appData.GetAppState()
        run WASM (Deposit if needed, then ProcessRequest. appState must be updated only after ProcessRequest)

        if soft failure:
            results = append(results, buildErrorPayload(req, currentStateRoot))
            continue  // state unchanged, next request

        if hard failure (encryption):
            break  // stop batch, return results so far

        // Success — mutate state
        appData.SetAppState(newState)
        appData.IncrementNonce()
        serialized := appData.Serialize()
        currentStateRoot = SHA256(serialized)
        results = append(results, buildSuccessPayload(...))

    if len(results) == 0:
        return nil, nil, nil, nil  // nothing processed

    // Batch signature — one sign operation covering all entries.
    // personal_sign over the concatenated entry hashes (section 3.2):
    //   keccak256("\x19Ethereum Signed Message:\n" || itoa(32*N) || hash(results[0]) || ... || hash(results[N-1]))
    batchHash := personalSign(concat(hash(results[0]), ..., hash(results[N-1])))
    batchSignature, err := sign(batchHash)
    if err != nil:
        return nil, nil, nil, err  // signing failure — entire batch discarded

    encryptedFinalState := encryptState(appData)       // 1 encrypt
    return results, batchSignature, encryptedFinalState
```

One payload is returned per handled request, in input order, so `len(results)` is how many of the N input requests were handled (whether successfully or with error payloads) — there is no separate count to keep in sync. If `len(results) < N`, the manager knows a hard failure stopped the batch at request `len(results) + 1`.

If batch signing fails after processing, all results are discarded and the executor returns an error. The requests remain pending on-chain and will be retried on the next poll. This is a rare system-level failure (key unavailable, HSM error) — not an application-level concern.

### 7.4. System-wide head-of-line blocking — open issue

Per-application queues (section 4.2) confine a permanently failing request to its own queue *structurally*, but the enforced round-robin selection (section 4.3) re-globalizes the blockage:

1. Application A's head request hard-fails permanently (e.g., tampered `applicationId`). It is never dequeued — soft failures produce an error payload and advance the queue; hard failures leave the request at the head.
2. When the cursor reaches A, the contract serves A and only accepts a state update for A. The manager gets a hard failure on request 1, no payloads are returned, and can submit nothing.
3. The cursor never advances — a cursor advance requires a successful state update for A. Every subsequent poll selects A again.

Result: one poisoned request stalls **all** applications, not just A. The manager cannot skip A — the cursor check in the state-update functions rejects submissions for any other application. This is a property of *any* contract-enforced selection, not of round-robin specifically: whatever algorithm the contract enforces, a request that can never be processed blocks the rotation at its turn.

**Possible solutions**, in order of how much they preserve the anti-censorship goal:

1. **On-chain failure escalation.** The manager reports the hard failure on-chain (e.g., `markRequestBlocked(requestId)`, possibly gated by a timeout or evidence requirement). The contract moves the request to a parked state or the queue tail and advances the cursor, so rotation moves on. Verifiable, deterministic, auditable.
2. **Request expiry.** Section 7.2 mentions expiry as out of scope; enforced selection effectively makes it in scope, because expiry becomes the only *automatic* unblocking mechanism — once the blocked head expires, the queue advances (or empties, letting the scan skip the application) and rotation moves on.
3. **Manager-supplied skip/exclusion parameter** on the state-update call. Simplest, but reintroduces exactly the discretion section 4.3 removes — a manager could "skip" any application indefinitely. Note that the manager can already censor by idling; the real anti-censorship property is *detectability*, so an on-chain skip event with an emitted audit trail could be acceptable.

A mechanism from this list must be chosen before the design is complete; until then, a permanently failing request is a system-wide blocker.

## 8. Configuration

| Variable | Default | Description |
|---|---|---|
| `MAX_BATCH_SIZE` | `5` | Max requests per poll cycle (per batch, single application). Must be > 0. |
| `DataLayerNumOfVersions` | `10` | Depth of per-application version history retained for rollback and reorg recovery. A batch stores only its final state (one version per batch), so this is independent of `MaxBatchSize`; size it for the reorg-recovery depth you need. Currently hardcoded in `config.go`. |
