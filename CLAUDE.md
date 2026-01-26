# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Horizen Privacy Preserving Execution System (PES) - A privacy-preserving execution platform using AWS Nitro Enclaves (TEE). The system executes WebAssembly modules securely, with encrypted state management and blockchain-based coordination.

**Stack:** Go 1.24, Solidity 0.8.30, Wasmtime-go (WASM runtime), LevelDB, Hardhat

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

## Tool Requirements

- Go 1.24
- solc 0.8.30 (install via `solc-select`)
- abigen v1.16.2: `go install github.com/ethereum/go-ethereum/cmd/abigen@v1.16.2`
- TinyGo 0.39.0 (for WASM module compilation in tests)
- Node.js 20 (for contracts)
