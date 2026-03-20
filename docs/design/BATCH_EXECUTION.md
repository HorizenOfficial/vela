# Batch Execution: Increasing Request Throughput

## 1. Problem Statement

The system currently processes **one request per blockchain block**. On a chain with 2-second block times, this means a maximum throughput of ~30 requests per minute, regardless of how fast the executor can run WASM or how many requests are queued.

This is a bottleneck for any application with sustained request volume. Users see growing queue latency even when the system has idle capacity.

## 2. Current Implementation: Why It's Limited

The per-block bottleneck comes from three sequential dependencies in the manager's poll loop:

```
pollBlockchain()                                        (manager.go:472)
  └─ processRequestFromChain()                          (manager.go:500)
       ├─ GetNextPendingRequest()   ← fetches ONE       (manager.go:522)
       ├─ processRequest()                               (manager.go:585)
       │    ├─ SendProcessRequest() ← vsock round-trip   (manager.go:751)
       │    ├─ dataLayer.Store()    ← DB write           (manager.go:779)
       │    └─ submitStateOnChain() ← blocks until mined (manager.go:790)
       │         └─ SubmitStateUpdate()                   (blockchain/client.go:372)
       │              └─ sendTxAndWaitMined()             (blockchain/client.go:307)
       │                   ├─ bind.Transact()  ← send tx
       │                   └─ bind.WaitMined() ← BLOCKS here
       └─ return (next tick fires after polling interval)
```

Three things make this sequential:

### 2.1. Single-request fetch

`GetNextPendingRequest()` returns only the queue head. Even if 10 requests are pending, the manager sees one at a time.

### 2.2. Blocking transaction submission

`sendTxAndWaitMined()` calls `bind.Transact` (sends the transaction) and then `bind.WaitMined` (polls for inclusion in a block). The entire call blocks until the transaction is mined — typically one full block time (~2 seconds). No other work happens during this wait.

### 2.3. Poll-driven single-pass loop

`pollBlockchain()` runs on a ticker. Each tick processes exactly one request. The next tick fires only after the current one completes (fetch + execute + mine). The effective cadence is `max(pollingInterval, executionTime + miningTime)`.

### 2.4. Redundant work per request

Each request also carries overhead that compounds across sequential processing:

- **Manager → Executor round-trip**: Every request sends the full encrypted state + WASM bytes over vsock/TCP. The executor decrypts the state, runs WASM, re-encrypts the new state, and sends it back. For consecutive requests to the same application, the executor decrypts what it just encrypted one round-trip ago.

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

The `stateUpdate()` contract function processes a single request per call. Each call validates the queue head, verifies the signature, updates the state root, processes refunds/withdrawals, emits events, and dequeues the request.

For N requests, this means N separate transactions, each paying the base transaction cost (~21,000 gas) and each requiring its own nonce, gas estimation, and receipt tracking.

### 2.6. Summary of bottlenecks

| Bottleneck | Where | Impact |
|---|---|---|
| Single-request fetch | `GetNextPendingRequest()` | Sees 1 request per poll |
| Blocking mine wait | `bind.WaitMined()` in `sendTxAndWaitMined()` | ~2s idle wait per request |
| Redundant state crypto | Executor decrypt/encrypt per vsock round-trip | Wasted CPU + bandwidth |
| One vsock round-trip per request | `SendProcessRequest()` | Latency × N |
| One tx per request | `stateUpdate()` contract design | N × base gas cost |
| N nonce/receipt management | Sequential `bind.Transact` calls | Complexity + failure modes |

## 3. Considered Approaches

Both approaches share a common starting point: fetch up to `MaxBatchSize` pending requests via `GetPendingRequests(maxCount)` instead of one at a time. The fetch is capped at the source — there is no reason to transfer the entire pending queue only to discard most of it. They differ in how execution and on-chain submission are organized.

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

1. Fetch all pending requests + on-chain stateRoot
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

- **Gas estimation broken for tx 2+.** `bind.Transact` calls `estimateGas`, which simulates against the latest confirmed state. For transaction 2, the on-chain `stateRoot` hasn't been updated by the still-pending transaction 1, so `estimateGas` fails. Workaround: set `TransactOpts.GasLimit` to a fixed non-zero value to skip estimation. This requires tuning and risks over- or under-estimation.

- **Cascading failure.** The contract validates `prevStateRoot == stateRoot` and `isCurrentPendingRequest()` on each call. If transaction K reverts, transactions K+1..N all revert too (wrong `prevStateRoot`, wrong queue head). Gas is wasted on all the reverted transactions.

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

1. Fetch first N pending requests + on-chain stateRoot
2. Send all N requests + encrypted state + WASM module to executor in a single batch message
3. Executor decrypts state **once**, runs WASM for each request sequentially (keeping plaintext state in memory between calls), encrypts **once** at the end
4. Executor returns: N `UpdatePayload`s (one per request, each individually signed) + 1 final encrypted state
5. Manager stores the final encrypted state in DB (1 write)
6. Manager submits one `batchStateUpdate()` transaction with all N update payloads
7. Wait for the single receipt

**Contract changes:** New `batchStateUpdate()` function that loops over entries, calling the same internal validation logic as current `stateUpdate()`. Since the contract is not in production, this is a refactor of the existing code — extract the body of `stateUpdate()` into an internal `_processOneStateUpdate()` and call it in a loop. Additional optimizations:

- **Single stateRoot storage write**: read `stateRoot` once at the start, chain through entries in memory, write once at the end — saves (N-1) warm `SSTORE` operations (~5,000 gas each)
- **State root chain validation**: only the first entry checks `prevStateRoot == stateRoot` from storage; subsequent entries validate `entries[i].prevStateRoot == entries[i-1].newStateRoot` in memory
- **Deduplicated `applicationId`**: passed once instead of per-entry
- **Optional batch signature**: verify one signature over the full batch instead of N individual `ecrecover` calls, saving (N-1) × ~3,000+ gas (individual per-request signatures can still be emitted in events for off-chain verifiability)

**Pros:**

- **1 decrypt + 1 encrypt** instead of N of each. The executor holds plaintext state in memory between requests. For a batch of 5, this eliminates 4 decrypt cycles, 4 encrypt cycles, and the associated serialization/deserialization overhead.

- **1 vsock round-trip** instead of N. The full batch is sent and received in a single message exchange. No redundant state transfer.

- **Atomic on-chain submission.** All-or-nothing: either all N requests are processed or none are. No cascading failure, no partial state, no ambiguity about which transaction failed.

- **Gas estimation works normally.** A single transaction is simulated against the confirmed state. The first entry's `prevStateRoot` matches the on-chain `stateRoot`, and the contract processes subsequent entries internally — no pending-state simulation issues.

- **1 nonce, 1 receipt.** No nonce management complexity, no sequential nonce gaps, no receipt-per-transaction tracking.

- **Simple recovery.** If the transaction reverts, nothing changed on-chain. Roll back local state to pre-batch, retry the whole batch next poll. No two-poll recovery, no reorg detection needed for partial failures.

- **Gas savings.** Eliminates (N-1) × 21,000 base transaction cost. Additional savings from single stateRoot write and optional batch signature. For a batch of 5: ~84,000 gas saved from base costs alone.

- **1 DB write** instead of N. Only the final encrypted state needs to be stored.

- **Subgraph-transparent.** The contract emits the same events per request inside the loop (`StateRootUpdate`, `UserEvent`, `RequestCompleted`, etc.). Indexers see identical events, just all in the same block. No subgraph changes needed.

**Cons:**

- **Contract modification required.** New `batchStateUpdate()` function and internal refactoring. Needs review and testing. However, since the contract is not in production, this is a development-time cost, not a migration risk.

- **Block gas limit risk.** A batch of N updates with events and withdrawals in a single transaction could consume significant gas. Estimated ~200,000 gas per request entry means a batch of 5 uses ~1M gas. On L2s (block gas limits typically 30M+) this is well within limits. A configurable `MaxBatchSize` cap mitigates this.

- **All-or-nothing risk.** If the transaction runs out of gas mid-loop or any entry triggers an unexpected revert, all N results are lost — including entries that would have succeeded individually. With pipeline execution, the first K transactions would have already confirmed. However: the TEE pre-validates every entry, so an on-chain revert should only happen due to external state changes (new requests submitted between poll and mining), which would also break pipeline execution via cascading failure.

- **Larger calldata.** All N entries are ABI-encoded in one transaction. Calldata cost: 16 gas per non-zero byte, 4 per zero byte. On L2s with EIP-4844 blob data, this is cheap. On L1, it would add up for large batches.

- **Executor becomes batch-aware.** The communication protocol needs a new batch message type. The executor must handle partial failures within a batch (e.g., request 3 fails WASM execution — produce an error payload for it, continue with request 4 using the state from request 2). This is a moderate protocol change.

- **No partial on-chain progress.** Unlike pipeline execution where transactions 1..K might already be confirmed when K+1 fails, batch execution is all-or-nothing. If the batch fails, the system retries everything next poll. In practice, this rarely matters because the TEE pre-validates all entries, but it means a failed batch wastes more wall-clock time than a failed pipeline.

### 3.3. Comparison

| Criterion | Pipeline (A) | Batch (B) |
|---|---|---|
| Contract changes | None | New `batchStateUpdate()` |
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

## 4. Implementation Plan

The work is organized in four stages. Each stage produces a testable, reviewable unit. Stages 1-3 can proceed in parallel across different developers.

### Stage 1 — Smart Contract: `batchStateUpdate()`

Refactor `ProcessorEndpoint.sol` to support batch submission.

**Steps:**

1. Extract the body of `stateUpdate()` into an internal function `_processOneStateUpdate()` that takes the same parameters and performs all validation, state updates, event emission, refunds, withdrawals, and request dequeuing.

2. Rewrite `stateUpdate()` as a thin wrapper that calls `_processOneStateUpdate()` once. This preserves backward compatibility and confirms the refactor is correct — all existing contract tests must still pass without modification.

3. Add `batchStateUpdate(uint64 applicationId, BatchEntry[] calldata entries, bytes[] calldata signatures)` that:
   - Reads `stateRoot` from storage once into a local variable
   - Loops over entries, calling `_processOneStateUpdate()` for each
   - Validates state root chaining: first entry checks `prevStateRoot == stateRoot` from storage; subsequent entries check `entries[i].prevStateRoot == entries[i-1].newStateRoot`
   - Writes `stateRoot` to storage once at the end of the loop (not per iteration)
   - Passes `applicationId` once (not per entry)

4. Define the `BatchEntry` struct containing per-request fields: `prevStateRoot`, `newStateRoot`, `processedRequestId`, `events`, `eventSubTypes`, `withdrawalRequests`, `refund`, `applicationFees`, `errorCode`, `errorMsg`.

5. Write contract tests:
   - Batch of N successful requests
   - Batch with an error payload mid-batch (request K fails, K+1 continues from unchanged state)
   - Batch with first entry having wrong `prevStateRoot` (reverts)
   - Batch with broken state root chain between entries (reverts)
   - Single-entry batch (equivalent to `stateUpdate()`)
   - Gas measurement: compare `batchStateUpdate(N entries)` vs N × `stateUpdate()`

6. *(Optional, can be deferred)* Batch signature: replace N individual `ecrecover` calls with a single signature over `keccak256(abi.encode(entry1_hash, entry2_hash, ...))`. Emit individual entry hashes in events for off-chain verifiability.

**Files changed:**
- `contracts/contracts/ProcessorEndpoint.sol`
- `contracts/contracts/IProcessorEndpoint.sol` (interface)
- `contracts/contracts/Structs.sol` (new `BatchEntry` struct)
- `contracts/test/` (new and updated test files)

### Stage 2 — Executor: Batch Processing

Add a batch message type to the executor so it can process multiple requests in a single vsock round-trip, keeping decrypted state in memory between requests.

**Steps:**

1. Define the batch request/response message types in the communication protocol:
   - `BatchProcessRequestMessage`: carries `[]*common.Request`, `*common.ApplicationState`, `[]byte` (WASM module)
   - `BatchProcessResponseMessage`: carries `[]*common.UpdatePayload`, `*common.ApplicationState` (final only), `[]*common.DeanonymizationReport`

2. Implement `HandleBatchProcessRequest()` in the executor. The batch loop must distinguish between soft failures (signed error payload — continue) and hard failures (bare error — stop batch). See section 5 for the full pseudocode and rationale. Key invariants:
   - `appData` is only mutated after successful WASM execution
   - On soft failure: error payload included in results, state unchanged, batch continues
   - On hard failure: batch stops, results for previously processed requests are returned
   - The response includes `processedCount` so the manager knows how many requests were handled

3. Add the message handler in `communication/server.go` to route `BatchProcessRequestMessage` to the executor's `HandleBatchProcessRequest`.

4. Add `SendBatchProcessRequest()` to the executor client interface and implement it in the communication client.

5. Write executor tests:
   - Batch of N successful requests — verify state roots chain correctly
   - Soft failure mid-batch — request K gets error payload, K+1..N continue from K-1's state
   - Hard failure mid-batch (e.g., wrong applicationId on request K) — batch stops, only 1..K-1 returned
   - Request with deposit + process where process fails — verify deposit changes are discarded
   - Single-request batch — equivalent to existing `SendProcessRequest`
   - Verify only 1 AES decrypt and 1 AES encrypt occur for the batch
   - Verify `processedCount` matches number of results returned

**Files changed:**
- `pkg/communication/messages.go` (new message types)
- `pkg/communication/server.go` (new handler)
- `pkg/communication/client.go` (new `SendBatchProcessRequest`)
- `pkg/executor/executor.go` (new `HandleBatchProcessRequest`)
- `pkg/blockchain/interface.go` (executor client interface update)
- `pkg/blockchain/mock_client.go` (mock implementation)

**Deploy requests:** Deploy is a different flow (creates initial state, stores WASM bytecode). Batching deploys with process requests adds complexity for little benefit (deploys are rare). For the initial implementation, if the batch contains a deploy request, the manager should process it individually before the batch, or exclude it from the batch. This can be refined later.

### Stage 3 — Go Contract Bindings

Regenerate the Go bindings after the contract changes so the manager can call `batchStateUpdate()`.

**Steps:**

1. Regenerate bindings: `go generate ./...`
2. Add `SubmitBatchStateUpdate()` to `BlockChainClient` that:
   - Takes `[]*common.UpdatePayload` (the batch results)
   - Packs all entries into `BatchEntry[]` calldata
   - Packs individual signatures into `bytes[]`
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

1. Add `MaxBatchSize` configuration to `config.go` with `MAX_BATCH_SIZE` env var (default 5). Add startup validation: `DataLayerNumOfVersions >= MaxBatchSize + 5`.

2. Add `GetPendingRequestsWithStateRoot(maxCount uint64)` (or similar) to the blockchain client that returns `([]*common.Request, [32]byte, error)` — up to `maxCount` pending requests plus the on-chain `stateRoot`. The caller passes `MaxBatchSize` so only the requests that will actually be processed are fetched, avoiding unnecessary data transfer when the queue is large. The existing `GetPendingRequests()` only returns requests without the stateRoot.

3. Implement `processBatchFromChain()` in the manager:
   ```
   requests, onChainStateRoot := GetPendingRequestsWithStateRoot(MaxBatchSize)
   verify localStateRoot == onChainStateRoot (existing reorg logic)

   separate deploy requests from process/deanonymize requests

   for each deploy request (if any):
       process individually via existing processDeployApp()

   if process/deanonymize requests remain:
       results, finalState, processedCount := executor.SendBatchProcessRequest(
           requests, encryptedState, wasmBytes)

       if processedCount == 0:
           // Hard failure on the very first request — nothing to submit
           log warning, retry next poll

       if processedCount > 0:
           save deanonymization reports to disk (if any)
           store final encrypted state in DB (1 write, versionID = final stateRoot)
           submit batchStateUpdate(results[0..processedCount-1]) on chain
           on tx failure: rollback DB to pre-batch state, retry next poll

       if processedCount < len(requests):
           log that request [processedCount] caused a hard stop
           // remaining requests stay in the on-chain queue for the next poll
   ```

4. Update `pollBlockchain()` to call `processBatchFromChain()` instead of `processRequestFromChain()`.

5. Write manager tests:
   - Happy path: N requests batched, single tx confirmed
   - Mixed types: deploy + process requests — deploy handled individually, rest batched
   - Batch tx reverts: verify DB rollback to pre-batch state
   - Empty queue: no-op
   - Queue larger than MaxBatchSize: only first MaxBatchSize requests fetched and processed
   - Deanonymization reports saved correctly within batch
   - Reorg detection still works with batch state storage
   - Hard failure mid-batch: executor returns partial results, manager submits only processed requests, remaining stay pending
   - Hard failure on first request: nothing submitted, retry next poll

6. Integration test: submit 5 requests on-chain, verify all processed in one poll cycle via a single `batchStateUpdate()` transaction.

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

## 5. Failure Handling Within the Executor Batch

The current executor has two distinct error types, and the batch must handle them differently.

### 5.1. Soft failure — signed error payload (request dequeued)

The executor produces a signed `UpdatePayload` with `prevStateRoot == newStateRoot` (state unchanged) and a non-zero `ErrorCode`. The contract marks the request as `FAILED`, collects the minimum fee, refunds the rest, and advances the queue head.

This happens for application-level errors where the executor has a valid stateRoot and can produce a signed attestation:
- App state not found (`executor.go:556`)
- WASM execution failure — deposit or process (`executor.go:574, 634`)
- Insufficient fuel (`executor.go:586, 674`)
- Payload decryption failure (`executor.go:625`)
- Deanonymize report validation failure (`executor.go:650, 657`)
- AppData serialization failure (`executor.go:698`)
- Event/report encryption failure (`executor.go:711, 764`)

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

### 5.2. Hard failure — bare error (request stays pending)

The executor returns `(nil, nil, nil, error)` — no signed payload. The manager cannot submit anything on-chain for this request. The request remains in the pending queue.

This happens for system-level or pre-validation errors where the executor cannot or should not produce a signed payload:
- `validateRequest()` failure: wrong `applicationId`, wrong `protocolVersion`, fee below minimum (`executor.go:530-538`)
- Unsupported request type (`executor.go:552`)
- State decryption failure (`executor.go:565`)
- AES encryption failure (`executor.go:706`)
- Payload signing failure (`executor.go:741`)

**Batch behavior: stop.** The executor cannot produce a result for this request, and the contract requires FIFO processing — request K cannot be skipped to process K+1. The batch stops at request K. Results for requests 1..K-1 are returned; requests K..N are not processed and remain pending.

```
Request 1: success    → state₁     ← included in batch results
Request 2: success    → state₂     ← included in batch results
Request 3: hard fail  → bare error ← batch stops here
Request 4: not executed             ← remains pending
Request 5: not executed             ← remains pending
```

The manager submits a `batchStateUpdate()` with only the results for requests 1-2. Requests 3-5 remain in the on-chain queue and will be retried on the next poll.

> **Note:** Some hard failures are transient (signing failure, encryption failure — likely a system issue that will resolve). Others are permanent for this executor (wrong `applicationId`). In the permanent case, the request blocks the queue head — every subsequent poll will stop at the same request. This is the same behavior as today's single-request processing: the manager retries and fails each poll. A separate mechanism (admin intervention, request expiry, or multi-app support) would be needed to unblock the queue. This is out of scope for the batch design.

### 5.3. Executor batch pseudocode (updated)

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

        if hard failure (encryption, signing):
            break  // stop batch, return results so far

        // Success — mutate state
        appData.SetAppState(newState)
        appData.IncrementNonce()
        serialized := appData.Serialize()
        currentStateRoot = SHA256(serialized)
        results = append(results, buildSuccessPayload(...))

    encryptedFinalState := encryptState(appData)       // 1 encrypt
    return results, encryptedFinalState, processedCount
```

The response includes `processedCount` so the manager knows how many of the N input requests were handled (whether successfully or with error payloads). If `processedCount < N`, the manager knows a hard failure stopped the batch at request `processedCount + 1`.

## 6. Configuration

| Variable | Default | Description |
|---|---|---|
| `MAX_BATCH_SIZE` | `5` | Max requests per poll cycle |
| `DataLayerNumOfVersions` | `10` | Must be >= `MaxBatchSize + 5`. Currently hardcoded in `config.go`. |

## 7. Multi-Application Support

The current design assumes a single application — one WASM module, one encrypted state, one `stateRoot`. When the system supports multiple applications, the batch design requires targeted changes. This section outlines what changes and what stays the same.

### 7.1. On-chain state model

The single `stateRoot` storage variable becomes a per-application mapping:

```solidity
// Before
bytes32 public stateRoot;

// After
mapping(uint64 => bytes32) public stateRoots;
```

Each application has its own independent state root. `batchStateUpdate()` already receives `applicationId` as a parameter — the change is that it reads and writes `stateRoots[applicationId]` instead of `stateRoot`. The internal validation logic (`_processOneStateUpdate`) is unchanged: it chains `prevStateRoot → newStateRoot` across entries the same way, just scoped to one application's root.

### 7.2. Per-application pending queues

A single global FIFO queue cannot work with multi-application batches. If requests for applications A, B, and A are enqueued in that order, a batch for application A cannot skip the request for B — the contract enforces FIFO dequeuing.

The contract needs **per-application pending queues**. Each application maintains its own queue head/tail, and requests are enqueued into `pendingQueues[applicationId]`. `batchStateUpdate()` dequeues from the queue matching its `applicationId`.

This also resolves the hard-failure queue-blocking problem from section 5.2: a permanently failing request for application A only blocks A's queue, not requests for other applications.

### 7.3. Fetching: scoped by application

`GetPendingRequestsWithStateRoot` adds an `applicationId` parameter:

```
GetPendingRequestsWithStateRoot(applicationId uint64, maxCount uint64) ([]*common.Request, [32]byte, error)
```

The manager queries one application's queue at a time and receives that application's `stateRoot`. This keeps the fetch lightweight — no cross-application data is transferred.

An additional view function is needed to discover which applications have pending requests:

```
GetApplicationsWithPendingRequests() ([]uint64, error)
```

This returns the list of application IDs that have non-empty queues, so the manager knows which applications to poll.

### 7.4. Manager: one batch per application per poll

The poll loop iterates over applications with pending requests. Each application produces an independent batch:

```
applicationIds := GetApplicationsWithPendingRequests()

for each applicationId in applicationIds:
    requests, stateRoot := GetPendingRequestsWithStateRoot(applicationId, MaxBatchSize)
    verify localStateRoot[applicationId] == stateRoot

    // Separate deploys from process/deanonymize (same as single-app)
    // ...

    results, finalState, processedCount := executor.SendBatchProcessRequest(
        requests, encryptedState[applicationId], wasmBytes[applicationId])

    store finalState for applicationId
    submit batchStateUpdate(applicationId, results) on chain
```

Each batch is self-contained: its own encrypted state, WASM module, state root, and on-chain transaction. Batches for different applications are independent — a failure in application A's batch does not affect application B.

**Sequential vs parallel batches:** The simplest approach is sequential — process one application's batch, wait for mining, move to the next. This reuses the existing single-batch flow with no additional complexity. If throughput across many active applications becomes a bottleneck, batches for different applications could be submitted in parallel (different `applicationId` values use independent state roots, so there are no nonce-ordering or state-chaining conflicts between them). This is the multi-application analog of the "combining approaches" note in section 3.4 and can be deferred.

### 7.5. Executor: no changes

The executor already processes a batch scoped to a single application — one WASM module, one encrypted state, N requests. The `BatchProcessRequestMessage` carries all of this per-batch. The executor is unaware of whether other applications exist. No changes needed.

### 7.6. Summary of changes

| Component | Change | Scope |
|---|---|---|
| Contract: state storage | `stateRoot` → `mapping(uint64 => bytes32) stateRoots` | Storage layout |
| Contract: queue | Single global queue → `mapping(uint64 => Queue) pendingQueues` | Storage layout + enqueue/dequeue logic |
| Contract: `batchStateUpdate()` | Read/write `stateRoots[applicationId]` instead of `stateRoot` | Minimal — already receives `applicationId` |
| Contract: view functions | `GetPendingRequestsWithStateRoot` scoped by `applicationId`; new `GetApplicationsWithPendingRequests` | New view function |
| Manager: poll loop | Iterate over applications, one batch per application | Moderate — loop structure changes |
| Manager: state storage | Keyed by `applicationId` (local encrypted state, stateRoot tracking) | Storage key scheme |
| Manager: fetch | Pass `applicationId` to `GetPendingRequestsWithStateRoot` | Signature change |
| Executor | None | — |
| Batch message protocol | None | — |
