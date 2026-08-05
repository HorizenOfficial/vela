# Throughput Benchmarks

**Status:** design — companion to `THROUGHPUT_OPTIONS.md`.
**Purpose:** measure, with the same harness and workloads, the request throughput of the current implementation and of each candidate solution (Option 0, batch execution, …), so that decisions D1–D3 are backed by numbers and future changes have a regression baseline.

## 1. Scope and explicit non-goals

App parameters are **pinned** to isolate the orchestration path (manager ↔ executor ↔ chain), which is what the candidate solutions change:

| Parameter | Pinned value |
|---|---|
| WASM work per request | minimal (counter-bump apps, near-zero execution time) |
| App state size | small (< 1 KB) |
| Request payload size | small, fixed |
| Trigger fire rate (trigger apps) | 100% — every request spawns a TRUSTPROCESS |

> **Non-goal:** these benchmarks do **not** measure the impact of the WASM execution itself (heavy computation, large state, large payloads). That is an important dimension — e.g. batching gains shrink when WASM time dominates the cycle — and deserves its own investigation. The pinned values above must be kept in mind when reading any result.

Gas costs are recorded but the simulated backend's gas accounting is indicative, not an L2 fee model.

## 2. Environment

Built on the existing fullstack simulated-chain suite (`pkg/testutil/fullstack/`, `pkg/blockchain/testutil/sim_test_helper.go`): in-process manager + executor + `go-ethereum` simulated backend, real TinyGo WASM apps (`app/simple`, `app/trigger`).

- **Block time:** the harness's auto-mining ticker is configurable via `WithAutoMineInterval` (`pkg/blockchain/testutil/sim_test_helper.go`). Benchmarks default to **1 s** (`VELA_BENCH_BLOCK_MS`): the manager polling interval has whole-second granularity (`BlockchainPollingInterval` is an int of seconds) and must stay ≤ block time, so 1 s is the fastest consistent emulation without touching production config. Because block time is emulated, all headline metrics are **block-denominated** (requests/block, latency in blocks); wall-clock numbers are reported but not comparable across scalings.
- **Polling interval:** the fullstack suite pins it to 1 s — equal to the default block time, so the ticker never masks the architecture under test.
- **Gating:** benchmarks require Wasmtime + TinyGo and take minutes — excluded from both quick and full test passes. Run with:
  `VELA_BENCH=1 go test -v -timeout 30m -run TestBench ./tests/benchmark/`
  (tunables: `VELA_BENCH_N` requests, `VELA_BENCH_BLOCK_MS` block time).

## 3. Workloads

| # | Workload | Apps | Submission |
|---|---|---|---|
| S1 | Hot normal app | 1 normal app | all requests to the one app |
| S2 | Hot trigger app | 1 trigger app (100% fire) | all requests to the one app |
| S3 | Mixed | 10 apps: 5 normal + 5 trigger | uniform round-robin across apps |

Why these three: S1 isolates the pure batching gain; S2 isolates the trigger serialization (and later measures T1); S3 is the workload that discriminates the multi-app fetch strategies (2a/2b/2c) — with the global FIFO and interleaved apps, same-app-prefix batching (2b) degenerates here and the numbers will show it.

Each workload runs against each implementation under test (baseline, Option 0, batch, …) — workloads and implementations are independent axes of the run matrix.

## 4. Modes and metrics

### 4.1 Saturation mode (primary)

Pre-fill the queue with N requests (default N = 30, tunable via `VELA_BENCH_N`), then measure the drain. The measurement **window** opens at the block after setup traffic (deploys, key association, funding) has fully drained (`fromBlock = currentBlock + 1`) and closes at the last completion; only requests submitted inside the window are counted, so setup is excluded.

Every finalized request is reconstructed black-box from chain events into one **completion record** (`tests/benchmark/metrics.go`):

| Field | Meaning |
|---|---|
| `SubmitBlock` | block the request's `submitRequest` tx was mined in |
| `Block` | block of the state-update tx that finalized it (`RequestCompleted`) |
| `TxHash` | that state-update tx |
| `GasUsed` | gas used by that tx (same value on records sharing a tx) |
| `Failed` | finalized with a non-zero error code |

All metrics are derived purely from the record set by `ComputeSaturationMetrics`:

Each row notes whether it is a **rate**, a **total/count**, a **span**, an **average**, or a **per-request** value — no metric is a per-request average unless it says "average".

| Metric | Kind | How it is calculated |
|---|---|---|
| **Requests/block** (headline) | rate (aggregate) | `Completed / BlocksSpanned` — total successful completions ÷ total blocks in the window; **not** a per-request average. `Completed` = records with `Failed = false`; `BlocksSpanned = lastCompletionBlock − firstCompletionBlock + 1` (inclusive — the final block counts whole even if the drain ends mid-block). Throughput: successful settlements per mined block. |
| Completed / failed | totals (counts) | number of records split by `Failed`. A clean run has **0 failed**; any failure is flagged loudly in the report and invalidates the comparison. |
| Drain window (blocks) | span | the same `BlocksSpanned` used above (`lastCompletionBlock − firstCompletionBlock + 1`, inclusive), reported here with its endpoints (`firstBlock → lastBlock`). Not an average. |
| State-update txs | total (count) | number of **distinct** `TxHash` across the records. |
| Mean batch size | average | `records / State-update txs` — average requests finalized per state-update tx, over the whole run. **1.0 = no batching**; > 1 means the implementation packs several completions into one tx. |
| Batch size distribution | histogram (counts) | for each distinct tx, how many records share it, tallied into a histogram of group sizes, rendered `count×[size req/tx]`. Shows *whether* batching engages, not just its mean. |
| Gas/request | average | `TotalGas / records` — total gas ÷ total requests. `TotalGas` sums `GasUsed` **once per distinct tx**, so a batched tx's gas is amortized across all the requests it settles. |
| Latency p50 / p95 (blocks) | per-request (percentiles) | for each request, `Block − SubmitBlock` (blocks from submission-mined to completion-mined); the run reports the nearest-rank **p50 and p95** across all requests. These are percentiles of the per-request distribution — **not** an average (see below). |
| Wall clock | span (total) | a single elapsed span: real time from the first pre-fill submission to the last observed completion (`lastCompletion − firstSubmission`), **not** an average of per-request times. Depends on the emulated block time — reported for context, not comparable across block-time settings. |
| Executor round-trip / mine-wait totals | totals (sums) | each is the **sum** of that stage's per-call log timings (§4.4) over the window, reported alongside its share of wall clock and its call/tx count. Divide by the call count for the per-call average. |

**Latency here is a diagnostic, not a headline.** Because all N requests are enqueued up front, a request's latency is dominated by *queue-wait*: the k-th request must wait for the k−1 ahead of it to drain, so it is ≈ (N/2 × service time) and grows with N — it measures backlog depth, not the cost of servicing one request. True latency under load is the job of the open-loop mode (§4.2). The **nearest-rank** percentile is `sorted[ceil(p/100 × n) − 1]`: p50 is the median request, p95 the tail.

### 4.2 Open-loop latency mode (second phase)

Requests submitted at a fixed arrival rate; measured at three rates per implementation: ~50%, ~90%, and ~110% of the capacity measured in saturation mode.

- **Headline: completion latency, median and p95**, in blocks and wall-clock — submission tx mined → completion (per §4.3).
- Queue depth over time (stability check: at 110% the queue must grow, at 50% it must not).

### 4.3 Completion definition

- **Normal request:** the state-update tx that finalizes it (`RequestCompleted` / failed event) is mined.
- **Trigger-app request:** the **spawned TRUSTPROCESS completes** — the user request alone is not done until settlement. Correlation (implemented in `tests/benchmark/s2_saturation_test.go`): `_enqueueTrustedRequest` emits `RequestSubmitted` *inside the spawning request's own state-update transaction*, so a tx containing both `RequestCompleted(F)` and `RequestSubmitted(T)` binds F → T. This holds for any number of in-flight round-trips. The on-chain hook also survives T1/batch execution, but a batch tx carrying several completions and several spawned submissions cannot be paired by same-tx grouping alone — the harness guards this case loudly and the pairing must be extended when batch lands. For trigger round-trips, gas/request covers **both** state-update transactions of the pair.

### 4.4 Per-stage timing (log-based)

Two duration-carrying log statements (`key=value` in the message text, extracted by regex) are added to production code, parsed by the harness:

1. **Executor round-trip** — around the executor call in the manager: `SendBatchProcessRequest` in `processBatch`. One log line covers the whole batch, so the sample is a per-batch duration, not a per-request one — divide by the reported batch size to compare against the single-request baseline.
2. **Mine wait** — around `bind.WaitMined` in `pkg/blockchain/client.go` (split from tx submission).

These two explain most results (e.g. "baseline: 78% of cycle is mine wait"; "batch: mine wait amortized to 20%/request") and map observations back to bottlenecks B1–B6 of `THROUGHPUT_OPTIONS.md`. Full per-stage breakdown (crypto, DB, per-stage inside the executor) only if results are surprising.

## 5. Report

One table per (workload × implementation × mode) run, plus a comparison summary per workload. Suggested columns (saturation): requests/block, drain blocks, mean batch size, gas/request, executor-RT share, mine-wait share. Results are committed under `tests/benchmark/results/` as dated markdown so the numbers travel with the branch that produced them.

## 6. Implementation

| Piece | Where |
|---|---|
| Configurable auto-mine interval | `WithAutoMineInterval` option, `pkg/blockchain/testutil/sim_test_helper.go`; threaded through the fullstack suite constructors |
| No-wait queue pre-fill | `FullStackSystemTestSuite.SubmitRequestNoWait` (`pkg/testutil/fullstack/fullstack_suite.go`) |
| Completion tracking | black-box chain scan of `RequestSubmitted` / `RequestCompleted` events over the measurement window (`tests/benchmark/harness_test.go`) — no wrapper changes needed |
| TRUSTPROCESS correlation | same-tx grouping of `RequestCompleted` + `RequestSubmitted` (`tests/benchmark/s2_saturation_test.go`) |
| Duration log statements (§4.4) | `pkg/manager/manager.go` (`executor_roundtrip_ms`), `pkg/blockchain/client.go` (`mine_wait_ms`, logger attached via `SetLogger`); parsed from the manager's log file, windowed by byte offset to exclude setup traffic |
| Metrics + report | `tests/benchmark/metrics.go` (unit-tested pure functions), `tests/benchmark/report.go`; reports written to `tests/benchmark/results/` |
| Scenarios | `tests/benchmark/s1_saturation_test.go` (N deposits, one simple app), `s2_saturation_test.go` (N fires, GuardedTrigger app), `s3_saturation_test.go` (5 simple + 5 trigger apps, round-robin, aggregate + per-app report) |

**Queue cap:** the contract rejects submissions once 10 requests are pending (`maxQueueSize`, `ProcessorEndpoint.sol:45`) — and the excess submissions revert *silently* from the harness's no-wait perspective (gas estimation passes against the not-yet-filled pending state). The harness raises the threshold at boot (`SetQueueThreshold`) and verifies every pre-fill receipt (`confirmSubmissions`) so a reverted submission fails the run loudly instead of surfacing as a drain timeout.

S3 reports both aggregate and per-app metrics. Workloads: S1/S3 normal-app requests are small ETH deposits; S2/S3 trigger-app requests are `{"type":"fire"}` (unshield → sweep → re-shield → TRUSTPROCESS), with an initial deposit large enough that N percentage-based fires never round to zero.

Remaining (open) pieces: the open-loop latency mode (§4.2) and non-baseline implementation runs as the throughput solutions land.
