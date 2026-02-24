# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Horizen Privacy Preserving Execution System (PES) - A privacy-preserving execution platform using AWS Nitro Enclaves (TEE). The system executes WebAssembly modules securely, with encrypted state management and blockchain-based coordination.

**Stack:** Go 1.24, Solidity 0.8.30, Wasmtime-go (WASM runtime), LevelDB, Hardhat, The Graph (subgraphs)

## Language Stack

- **Primary languages:** Go, TypeScript, Solidity
- This is a blockchain/TEE/WASM project. Be aware of smart contract security patterns (reentrancy, pull payments) and WASM memory management concerns.
- **Solidity:** Check reentrancy, access control, gas optimization.
- **Go:** Check error handling, goroutine leaks, resource cleanup.

## Build Commands

```bash
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

# WASM guest (app/simple)
cd app/simple && make build            # Development build
cd app/simple && make production_build # Optimized build

# Subgraph (subgraphs/hcce)
cd subgraphs/hcce && npm run codegen   # Generate types from schema
cd subgraphs/hcce && npm run build     # Build subgraph
cd subgraphs/hcce && npm run test      # Run subgraph tests
```

## Project Structure

```
/cmd                        - Main application entrypoints
  /manager                  - Secure Processor Manager (main Go service)
  /executor                 - WASM Executor (runs in AWS Nitro Enclave)
  /authorityservice         - HTTP service for deanonymization reports
  /keytool                  - Key management utility
  /key_attestation_cmd      - Key attestation utility

/pkg                        - Shared Go library code
  /wasm                     - Wasmtime runtime integration
  /executor                 - WASM execution, state handling, event encryption
  /manager                  - Manager orchestration, blockchain interaction
  /blockchain               - Smart contract bindings and interaction layer
  /communication            - V-Socket & TCP messaging between Manager/Executor
  /storage                  - LevelDB versioning system and mock implementations
  /common                   - Shared data models, types, configuration
  /crypto                   - Key management (AES, P521, secp256k1, X25519)
  /logger                   - Zerolog wrapper and custom loggers
  /authorityservice         - Authority HTTP service (API types, endpoints)
  /subgraph                 - GraphQL subgraph client (queries, types)
  /testutil                 - Testing utilities and helpers

/app                        - TinyGo WASM guest applications
  /simple                   - Example WASM application

/contracts                  - Solidity smart contracts (Hardhat project)
  /contracts                - Solidity source files
  /test                     - Hardhat tests (TypeScript)

/subgraphs                  - The Graph subgraphs
  /hcce                     - HCCE indexer subgraph

/tests                      - System integration tests
/dockerfiles                - Container configuration
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
   - `ProcessorEndpoint.sol` - Request handling
   - `TeeAuthenticator.sol` - TEE attestation verification
   - `AuthorityRegistry.sol` - Authority management

5. **Subgraphs** (`subgraphs/`) - Blockchain event indexing:
   - Index on-chain events for efficient GraphQL querying

### Communication Flow

Manager and Executor use bidirectional messaging (V-Socket for Nitro, TCP fallback):
- Manager sends `ProcessRequest`, Executor may callback with `GetUserKeys`
- Handshake protocol on connection for keyset recovery (see `EXEC_MGR_HANDSHAKE.md`)

### Storage Layer (`pkg/storage/versioned_leveldb/`)

Versioned LevelDB with atomic transactions and rollback support:
- `LevelDBDataLayer` - High-level interface
- `VersionedLevelDBAppStateStore` - Application state with versioning
- `LevelDBUserKeyStore` - Non-versioned user keys

## Key Patterns

**Contract Bindings:** Generated via `go generate`, committed to repo. CI verifies bindings are up-to-date. If you modify contracts, regenerate and commit.

**Test Skipping:** Use `CI_FLAG=true` to skip tests requiring Wasmtime or external dependencies. Tests check `os.Getenv("CI_FLAG")`.

**Configuration:** Environment variables with `.conf` file fallbacks. Key configs:
- `CHANNEL_TYPE` - `tcp` or `vsock`
- `MANAGER_DATA_FOLDER` / `MANAGER_REPORTS_FOLDER`
- `CHAIN_PROCESSOR_ADDRESS` / `CHAIN_TEEAUTHENTICATOR_ADDRESS`

**Interface-Based Design:** Heavy use of interfaces for testability (ChainClient, ExecutorClient, DataLayer). Mock implementations in tests.

## Code Review Guidelines

### Go Code
- **Error handling**: Check all errors; avoid `_` for error returns unless explicitly justified
- **Concurrency**: Review goroutine lifecycle, channel usage, and potential race conditions
- **Context propagation**: Ensure `context.Context` is passed through call chains for cancellation
- **Resource cleanup**: Verify `defer` for closing resources (files, connections, locks)
- **Interface design**: Prefer small, focused interfaces; accept interfaces, return structs
- **Testing**: Look for table-driven tests, proper mocking, and edge case coverage
- **Logging**: Use structured logging (zerolog); avoid fmt.Print in production code

### TypeScript Code (Contracts/Subgraphs)
- **Type safety**: Avoid `any`; prefer explicit types and interfaces
- **Null checks**: Handle nullable values explicitly
- **BigInt handling**: Use appropriate BigInt libraries for blockchain values
- **Event handlers** (subgraphs): Ensure entity IDs are unique and deterministic
- **Contract tests**: Check for edge cases, revert conditions, and gas optimization

### GraphQL (Subgraphs)
- **Schema design**: Use appropriate scalar types (BigInt, Bytes, ID)
- **Entity relationships**: Review `@derivedFrom` directives for efficiency
- **Indexing**: Ensure indexed fields support common query patterns

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
