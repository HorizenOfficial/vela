# Changelog

## 0.2.1

### Fixes

- **Guest memory no longer accumulates across requests**: every guest call (deposit, process request, memory stats) now runs on a fresh WASM instance; compiled modules stay cached, so a call pays only instantiation and `load_module`. TinyGo guests (<= 0.42) never free large allocations, so a long-lived instance kept growing: apps with a state above a few hundred KB trapped after some hundred requests. A single request with a large state can still exhaust guest memory or exceed the manager timeout.
- **Guest calls are serialized**: a request retried by the manager after a timeout could run inside the same instance as the still running first attempt, and fail with a signed error even though the first attempt succeeded. Retries now wait for the previous attempt.
- **The deposit result is no longer logged**: it contained the full application state.
- **Application state size limit**: the executor rejects a request whose application state exceeds `EXECUTOR_MAX_APP_STATE_SIZE` bytes (default 10 MiB, 0 = no limit), before calling the guest, and discards a deposit, request or deploy result whose new state exceeds it. The request fails with the signed error `APP_STATE_TOO_LARGE` (category `REQUEST_FUNC_FAILED`) and the state is unchanged, so the app can never grow past the limit. Before, a guest call on a large state could take longer than the manager request timeout, and the manager retried it forever, blocking the request queue. A call that is slow for other reasons still can.

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


