# Changelog

## Unreleased

### Breaking changes

- **`ProcessorEndpoint.getFacilitatorNonce(address)` removed.** It duplicated the auto-generated getter of the `public facilitatorNonces` mapping; clients must read `facilitatorNonces(user)` instead. See `docs/design/FACILITATOR.md`.
- **`ProcessorEndpoint` constructor takes a new last argument**, the address of a deployed `ProcessorEndpointExtension`. Deploy the extension first.
- **A malformed facilitator signature now reverts `InvalidSignature`.** `submitRequestFor` recovers with `ECDSA.tryRecover`, so `ECDSAInvalidSignature`, `ECDSAInvalidSignatureLength` and `ECDSAInvalidSignatureS` are gone from the `ProcessorEndpoint` ABI. Clients that decode those three errors must handle `InvalidSignature` instead. Valid signatures behave exactly as before.
- **The contracts now target the `cancun` EVM version** (previously `paris`), for the ~720 bytes it saves on `ProcessorEndpoint`. Any chain they are deployed to must be at Cancun or later; Base, the deployment target, has been since March 2024.

### Changes

- **`ProcessorEndpoint` split for the EIP-170 size limit**: state moved to `ProcessorEndpointStorage`, and `submitRequestFor`, `submitDeployRequest`, `submitDeployRequestWithTrigger`, `adminReset`, `adminResetApps`, `updateQueueThreshold`, `updateMaxNumOfApplications`, `updateFeeCollector`, `addAllowedDeployer` and `removeAllowedDeployer` now live in `ProcessorEndpointExtension`, reached by `delegatecall`. Callers are unaffected — same address, ABI and selectors — but each of those entry points costs one extra cold-account access (~2,600 gas). Everything moved is off the per-request hot path. See `docs/design/PROCESSOR_ENDPOINT_SPLIT.md`.
- **New view `ProcessorEndpoint.extension()`** returns the extension the endpoint delegates to, so the pairing is readable on-chain rather than only from the deployment record.
- **WASM guest execution is now time-bounded.** A guest operation that exceeds `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` (new, default 10 s, max 300 s, `0` selects the default) is interrupted via wasmtime epoch interruption and the request fails on-chain with the `WASM_INTERNAL` error code and message `guest execution timed out`. Previously a non-terminating guest ran forever and — because the execution lock is held across guest calls — blocked every later request for *every* application until the enclave was restarted. The bound also covers module instantiation, so a module whose `start` section never returns can no longer wedge a deploy. Guest authors: see the new sizing and batch notes in `docs/design/WASM_HOST_ABI.md`; the reference apps return in single-digit milliseconds, so only a non-terminating loop should reach this.
- **Requests carry a caller execution budget.** The manager now tells the executor how long it will still wait (`executionBudgetMs` on process, deploy and batch messages, derived from `MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC`), and the executor shortens guest execution to fit. The field is optional and an absent value means "not supplied", so mixed-version peers interoperate. It can only ever *shorten* execution: the enclave keeps its own configured bound as the ceiling, and a budget-driven interruption is treated as host-side abandonment — the request stays pending and is retried rather than being signed — so no value crossing the TEE boundary can cause a user to be charged a fee.
- **The executor refuses a manager whose request timeout cannot cover the guest bound.** The manager now reports `MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC` in the connection handshake (`requestTimeoutMs` on `GetKeysetRecoveryResponse`), and the executor rejects the connection — before restoring or generating any keyset — unless `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` plus a 1 s safety margin is strictly below that timeout minus the 2 s reply margin (the safety margin covers epoch-interruption overshoot and the decode/decrypt work that happens before the deadline is armed). Both processes then fail at start-up with an error naming the two settings. Without this, a guest that outlives the manager's patience is interrupted and signed by the executor after the manager has already stopped listening, so the signed failure is never posted, no fee is charged, and the request is retried forever: one such request occupies the executor indefinitely for free. An earlier start-up check compared against the executor's own `EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC` as a stand-in; that setting governs calls in the other direction and is no longer consulted for this. An older manager that does not report its timeout is served with a warning.
- **Executor shutdown no longer waits out a running guest.** `Close` interrupts in-flight guest calls instead of blocking, keeping a bounded lock acquisition only as a backstop for the paths epochs cannot interrupt (module compilation and time spent inside a WASI host call, both still unbounded).

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


