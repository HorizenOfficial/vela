# Changelog

<!--
Style: terse, like the 0.2.0 section. A bold title plus one or two sentences per bullet, under
### Features, ### Changes and ### Fixes (breaking items go in Changes; no separate section).
Link the relevant docs/design/*.md instead of explaining rationale or internals, say each thing
once, and keep what a reader must act on: new config vars with defaults, ABI changes, required
rebuilds, CVE ids.
-->

## 0.3.0

### Features

- **Batch execution**: a poll settles up to `MAX_BATCH_SIZE` (new, default 5) requests of one application in a single enclave round-trip and transaction, via `ProcessorEndpoint.batchStateUpdate` and one signature over all entries (`BatchProcessed` event). Triggers, deploys and apps with a trigger contract stay limited to one entry. See `docs/design/BATCH_EXECUTION.md`.
- **Upgradeable contracts**: `ProcessorEndpoint`, `TeeAuthenticator` and `AuthorityRegistry` run behind UUPS proxies, with upgrade scripts in `contracts/scripts/upgrade/`. See `docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`.
- **Per-application request queues**: the contract selects the next requests round-robin via `getPendingRequestsWithStateRoot(maxCount)`, with global queues for triggers and deploys served first.

### Changes

- **Horizen L3 removed**: Vela targets Base (Sepolia for testing). The `horizen-l3-testnet` and `sepolia-testnet` Hardhat networks are replaced by `base-sepolia`, and the subgraph uses `network: base-sepolia` (set `startBlock` to the deployment blocks).
- **Contracts target `cancun`** (was `paris`): deployment chains must be Cancun or later.
- **Proxy deployment**: setup arguments moved from constructors to `initialize(...)`; use the **proxy** address. The `contracts/.openzeppelin/` manifest must be committed. `TeeAuthenticator.nitroProver` and `maxVerificationAge` are no longer `immutable`.
- **`ProcessorEndpoint` split** for the EIP-170 limit: off-hot-path entry points moved to `ProcessorEndpointExtension` (`delegatecall`, same ABI, ~2,600 extra gas each), new view `extension()`. See `docs/design/PROCESSOR_ENDPOINT_SPLIT.md`.
- **`stateUpdate` enforces the round-robin turn** (`ApplicationNotSelected`, `PriorityQueueNotServed`), ignoring queue heads younger than `selectionGrace` (default 60s, `updateSelectionGrace`).
- **Queue views**: `getNextPendingRequest()` removed in favour of `getPendingRequestsWithStateRoot`; `getPendingRequests*` are no longer a global FIFO; new `requestById`.
- **ABI**: `getFacilitatorNonce` removed (use `facilitatorNonces`); malformed facilitator signatures revert `InvalidSignature`; `ErrorCode.APP_NOT_DEPLOYED` removed, shifting later enum values.
- **WASM guests**: only allowlisted WASI imports are accepted and the feature set is pinned; wasmtime upgraded from v1.0.0 (see Fixes). See `docs/design/WASM_HOST_ABI.md`.
- **Guest execution is time-bounded** by `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` (new, default 10 s) and shortened to the manager's request budget; the executor refuses a manager whose timeout cannot cover it.
- **Application state limit**: `EXECUTOR_MAX_APP_STATE_SIZE` (new, default 10 MiB); larger states fail with `APP_STATE_TOO_LARGE`.

### Fixes

- **Guest memory no longer accumulates**: each guest call runs on a fresh instance (ported from 0.2.1). Guests must not keep state in globals.
- **Guest memory and tables are capped** by `EXECUTOR_MAX_GUEST_MEMORY_BYTES` (new, default and maximum 2 GiB) and a table limit.
- **A faulting guest's module is evicted from the cache**; evictions no longer close a store mid-call.
- **Wasmtime upgraded to 50.0.0** for CVE-2026-115889 and CVE-2026-115893 (Cranelift sandbox escapes, CVSS 9.9). Executor images must be rebuilt: PCRs change, so update attestation and KMS policies.
- **Interrupted guests**: `GUEST_EXECUTION_TIMEOUT` is signed on-chain; `GUEST_EXECUTION_CANCELLED` is transient (retried, no fee).
- **Duplicate deploys are rejected before compilation**; guest results are no longer logged.
- **`ADMIN` and `UPDATE_STATUS_ROLE` can be rotated and revoked** in `ProcessorEndpoint` and `TokenAllowlist`.
- **Start-up config checks**: ports are range-checked, and `GetConfigVarInt64`/`DEPLOY_ARTIFACTS_MAX_SIZE_MB` no longer overflow or silently fall back to defaults.
- **CI**: `solc` pinned to 0.8.30, new `npm run check:layout`, EIP-170 enforced in tests. New design docs and a benchmark suite in `tests/benchmark/`.

## 0.2.0

### Features

- **Smart contracts invocation**: Added possibility to invoke an external smart contract as a result of a TEE invocation + priority queue for following reentrant calls
- **Admin reset (testnet/development)**: new `RESET_OPERATOR` role with `adminReset` (clears the pending request queue and frees deploy slots) and `adminResetApps` (resets per-app state roots and locked funds, sweeping accumulated ETH/ERC-20 balances to the caller to avoid fund loss). The feature is permanently disabled when `RESET_OPERATOR` is initialised as `address(0)`, the expected production value. See `docs/design/PROCESSOR_ENDPOINT_ADMIN_RESET.md`.

## 0.1.0

### Features

- **ERC-20 support**: deposits, withdrawals and fee handling via on-chain `ProcessorEndpoint`, with EIP-2612 permit flow and a facilitator path for gasless user onboarding.
- **Multi-app support**: the Manager, Executor, storage layer, smart contracts and subgraph now handle multiple applications concurrently, with per-app isolated state and locked funds.
- **Deploy flow**: on-chain deploy descriptor with TEE WASM fingerprint verification and on-chain deployer-role check.
- **App events**: WASM apps can emit typed application events; event subtype is now a fixed 32-byte value.


