# Changelog

## Unreleased

### Breaking changes

- **`ProcessorEndpoint.getFacilitatorNonce(address)` removed.** It duplicated the auto-generated getter of the `public facilitatorNonces` mapping; clients must read `facilitatorNonces(user)` instead. See `docs/design/FACILITATOR.md`.
- **`ProcessorEndpoint` constructor takes a new last argument**, the address of a deployed `ProcessorEndpointExtension`. Deploy the extension first.

### Changes

- **`ProcessorEndpoint` split for the EIP-170 size limit**: state moved to `ProcessorEndpointStorage`, and `submitRequestFor` now lives in `ProcessorEndpointExtension`, reached by `delegatecall`. Callers are unaffected — same address, ABI and selectors. See `docs/design/PROCESSOR_ENDPOINT_SPLIT.md`.

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


