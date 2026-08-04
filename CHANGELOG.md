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


