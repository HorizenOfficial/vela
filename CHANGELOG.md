# Changelog

## 0.2.0

### Features

- **Smart contracts invocation**: Added possibility to invoke an external smart contract as a result of a TEE invocation + priority queue for following reentrant calls

## 0.1.0

### Features

- **ERC-20 support**: deposits, withdrawals and fee handling via on-chain `ProcessorEndpoint`, with EIP-2612 permit flow and a facilitator path for gasless user onboarding.
- **Multi-app support**: the Manager, Executor, storage layer, smart contracts and subgraph now handle multiple applications concurrently, with per-app isolated state and locked funds.
- **Deploy flow**: on-chain deploy descriptor with TEE WASM fingerprint verification and on-chain deployer-role check.
- **App events**: WASM apps can emit typed application events; event subtype is now a fixed 32-byte value.


