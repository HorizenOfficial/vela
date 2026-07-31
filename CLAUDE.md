# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Horizen Vela - A privacy-preserving execution platform using AWS Nitro Enclaves (TEE). The system executes WebAssembly modules securely, with encrypted state management and blockchain-based coordination.

## Language Stack

- This is a blockchain/TEE/WASM project. Be aware of smart contract security patterns (reentrancy, pull payments) and WASM memory management concerns.

## Build Commands

```bashExpand commentComment on line L20
# Go build
go build ./...
go build ./cmd/executor
go build ./cmd/manager

# Run tests (quick suite - skips Wasmtime-dependent tests)
CI_FLAG=true go test -v ./...

# Run tests (full suite - includes all tests)
CI_FLAG= go test -v ./...

# Test with coverage
go test ./... -cover

# Run a single test
go test -v -run TestName ./path/to/package

# Generate contract bindings (requires solc 0.8.30 and abigen v1.16.2)
go generate ./...

# Contracts (from contracts/ directory)
npm ci
npm run build
npm run test
npm run check:layout   # ProcessorEndpoint / ProcessorEndpointExtension storage layouts must match

# WASM guest (app/simple)
cd app/simple && make build            # Development build
cd app/simple && make production_build # Optimized build

# Subgraph (subgraphs/hcce)
cd subgraphs/hcce && npm run codegen   # Generate types from schema
cd subgraphs/hcce && npm run build     # Build subgraph
cd subgraphs/hcce && npm run test      # Run subgraph tests
```

## Architecture

### Core Components

1. **Secure Processor Manager** (`cmd/manager/`, `pkg/manager/`) - Orchestrates execution flow:
   - Fetches requests from blockchain
   - Manages state persistence via LevelDB
   - Communicates with Executor via V-Socket/TCP
   - Posts results back to blockchain

2. **WASM Executor** (`cmd/executor/`, `pkg/executor/`) - Runs in AWS Nitro Enclave:
   - Executes WASM modules via Wasmtime
   - Handles encrypted state decryption/encryption
   - Manages cryptographic keysets (P521, secp256k1, AES)
   - Produces signed attestation payloads

3. **Authority Service** (`cmd/authorityservice/`, `pkg/authorityservice/`) - HTTP endpoints for deanonymization reports:
   - `GET /nonce` - Issues challenges
   - `POST /getreport` - Serves reports after on-chain verification

4. **Smart Contracts** (`contracts/`) - On-chain coordination:
   - `ProcessorEndpoint.sol` - Request handling, ERC-20 deposits/withdrawals (with EIP-2612 permit and facilitator path), per-app locked funds, deploy descriptor and allowed-deployer roles. Sits close to the EIP-170 size limit: state is declared in `ProcessorEndpointStorage.sol`, and some entry points (currently `submitRequestFor`) are hosted in `ProcessorEndpointExtension.sol` and reached by `delegatecall`. Declare new state **only** in the storage base, and run `npm run check:layout` after touching it — a layout divergence corrupts storage silently instead of reverting. Deploy the extension first; its address is the endpoint's last constructor argument. See `docs/design/PROCESSOR_ENDPOINT_SPLIT.md`
   - `TeeAuthenticator.sol` - TEE attestation verification (PCR-based, including WASM fingerprint for deploys)
   - `AuthorityRegistry.sol` - Authority management

5. **Subgraphs** (`subgraphs/`) - Blockchain event indexing:
   - Index on-chain events for efficient GraphQL querying (multi-app aware, deploy-specific events)

### Communication Flow

Manager and Executor use bidirectional messaging (V-Socket for Nitro, TCP fallback):
- Manager sends `ProcessRequest`, Executor may callback with `GetUserKeys`
- Handshake protocol on connection for keyset recovery (see `docs/design/EXEC_MGR_HANDSHAKE.md`)
- Manager also forwards admin commands to the Executor over the same channel (`ForwardAdminCommand`)
- `DeployApp` messages carry the WASM bytes and descriptor for the deploy flow

### Storage Layer (`pkg/storage/versioned_leveldb/`)

Versioned LevelDB with atomic transactions and rollback support:
- `LevelDBDataLayer` - High-level interface
- `VersionedLevelDBAppStateStore` - Application state with versioning (version chains are per-app — pruning and rollbacks are independent across apps)
- `LevelDBUserKeyStore` - Non-versioned user keys

### Multi-app Support

Manager, Executor, storage, contracts and subgraph are multi-app aware: each app has its own state, WASM module, locked funds and version chain. Deploy derives a unique `applicationId` from the deploy `requestId` (`ProcessorEndpoint.sol`), so deploys do not collide with regular requests.

### ERC-20 Flow

- Deposits use EIP-2612 `permit` to authorize the transfer from the user to the `ProcessorEndpoint`; a **facilitator** can relay the on-chain call so the user does not need gas.
- Funds are tracked per app in `ProcessorEndpoint`; withdrawals are authorized by signed updates from the TEE.
- See `docs/design/ERC20_DEPOSITS_WITHDRAWALS_DESIGN.md` and `docs/design/FACILITATOR.md`.

### Deploy Flow

- Apps are deployed on-chain via `DeployRequestSubmitted` on `ProcessorEndpoint`; only addresses holding the `DEPLOYAPP` role can submit deploys.
- The Manager uploads the WASM artifact to a shared folder (or via `POST /deploy/upload` on the authority service) and forwards it to the Executor.
- The Executor verifies the WASM fingerprint against the descriptor before loading the module.

## Key Patterns

**Contract Bindings:** Generated via `go generate`, committed to repo. CI verifies bindings are up-to-date. If you modify contracts, regenerate and commit.

**Third-Party Notices:** Whenever dependencies in `go.mod` change (add, remove, or version bump), update the `NOTICES` file to match — reconcile listed versions and add/remove entries for notable direct dependencies and their licenses.

**File Formatting:** If you modify contracts or TypeScript files, run `npm run format` after any modification to keep the correct formatting.

**Test Skipping:** Use `CI_FLAG=true` to skip tests requiring Wasmtime or external dependencies. Tests check `os.Getenv("CI_FLAG")`.

**Configuration:** Environment variables with `.conf` file fallbacks. Key configs:
- `CHANNEL_TYPE` - `tcp` or `vsock`
- `MANAGER_DATA_FOLDER` / `MANAGER_REPORTS_FOLDER`
- `CHAIN_PROCESSOR_ADDRESS` / `CHAIN_TEEAUTHENTICATOR_ADDRESS`

**Interface-Based Design:** Heavy use of interfaces for testability (ChainClient, ExecutorClient, DataLayer). Mock implementations in tests.

**User Signatures:** When a signature is provided by an external user (e.g. via a wallet), it must be compatible with MetaMask's `personal_sign` format. This means the message is prefixed with `\x19Ethereum Signed Message:\n<length of message>` before hashing and signing. Both Go verification code and Solidity `ecrecover` usage must account for this prefix.

## Code Review Guidelines

### Go Code
- **Logging**: Use structured logging (zerolog); avoid fmt.Print in production code

### TypeScript Code (Contracts/Subgraphs)
- **Event handlers** (subgraphs): Ensure entity IDs are unique and deterministic

### Configuration
- **Docker env wiring**: Every time a new config param is introduced (e.g. a new field/env var in `pkg/manager/config.go` or another component's config), verify the corresponding environment variable is wired through the Docker files — added to the relevant service's `environment` list in `dockerfiles/docker-compose.yml` and documented with a default in `dockerfiles/.env.template`.

## Things to Avoid

- Don't assume full Go stdlib availability in TinyGo code
- Don't put secrets or sensitive logic outside the Nitro enclave boundary
- Don't use `console.log` in subgraph handlers (use `log.info` from graph-ts)
- Don't ignore TypeScript strict mode errors in contract tests
- Be mindful of the trust boundaries between TEE and non-TEE components

## Tool Requirements

- Go 1.24
- solc 0.8.30 (install via `solc-select`)
- abigen v1.16.2: `go install github.com/ethereum/go-ethereum/cmd/abigen@v1.16.2`
- TinyGo 0.39.0 (for WASM module compilation in tests)
- Node.js 20 (for contracts)
