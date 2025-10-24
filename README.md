# Horizen Privacy Preserving Execution System

This repository contains the implementation of the Horizen Privacy Preserving Execution System, which consists of two main components:

1. **Secure Processor Manager**: Interacts with smart contracts and executes requested actions by orchestrating services and the TEE.
2. **WASM Executor**: Executes WASM modules within a secure environment (AWS Nitro Enclave) and handles private data.

## Technology Stack

The system is implemented in Go and uses Wasmtime-go as the runtime for WASM Modules.

### Structure

```
horizen-pes/
    │
    ├── cmd/                            # Main applications entrypoints for this project.
    │   ├── manager/                    # Secure Processor Manager application
    │   └── executor/                   # WASM Executor application
    │
    ├── contracts/                      # Smart contracts for the system
    │
    ├── pkg/                            # Library code that is used by other applications.
    │   ├── blockchain/                 # Blockchain interaction
    │   ├── common/                     # Data models and structures
    │   ├── communication/              # V-Socket Communication
    │   ├── executor/                   # WASM Executor
    │   ├── manager/                    # Secure Processor Manager
    │   └── storage/                    # Persistent data storage layer
    │
    └── tests/                          # System tests and integration tests
```

## Key Components

1. **Blockchain interaction**: Interacts with smart contracts.
2. **Data models and structures**: Types, structs, and enums used in the system.
3. **Secure Processor Manager**: Interacts with blockchain, storage, communication and orchestrates the execution of requests.
4. **Persistent data storage layer**: Interacts with data layer (Amazon S3).
5. **V-Socket Communication**: Handles communication between the Secure Processor Manager and WASM Executor.
6. **WASM Executor**: Executes WASM modules in a secure environment (AWS Nitro Enclave).

### Secure Processor Manager

The Secure Processor Manager is the core component of the system that interacts with others and executes requested actions by orchestrating services and the TEE.

Key functionalities:
1. Fetch requests from the blockchain.
2. Validate requests.
3. Fetch application state related to the request from the data layer.
4. Invokes WASM Execution through the communication layer.
5. Persist updated state to the data layer.
6. Post state update to the blockchain.

### WASM Executor

The WASM Executor is responsible for executing WASM modules within a secure environment (AWS Nitro Enclave). It handles the private data, decrypts states, executes WASM code, and produces signed attestation payloads. This application should be pre-compiled into the enclave image file (.eif) together with the Wasmtime-go runtime.

Key functionalities:
1. Receive encrypted initial state, request payloads, and WASM modules via v-socket.
2. Decrypt the state using TEE keys.
3. Instantiate WASM modules and invoke WASM methods to apply request payloads to the state using Wasmtime-go.
4. Encrypt the new state and events.
5. Produce and sign update payloads.

### V-Socket Communication

The system uses v-socket for communication between the Secure Processor Manager and the WASM Executor. V-socket is a virtual socket interface that allows for efficient communication in environments like AWS Nitro Enclaves.
Until the v-socket communication is fully implemented, the system can use TCP as a fallback communication method, it interfaces with the same message structure and logic.

#### Bidirectional Communication Example
```
Client                                  Server
  |                                       |
  |---> ProcessRequest (ID: 789) -------->|
  |                                       | (needs user keys)
  |<--- GetUserKeys (ID: 101) <-----------|
  |---> UserKeysResponse (ID: 101)------->|
  |                                       | (completes deploy)
  |<--- ProcessRequestResponse (ID: 789)<-|
  |                                       |
```


## WASM Runtime

The system uses Wasmtime-go as the runtime for WASM Modules. Wasmtime is a standalone JIT-style runtime for WebAssembly, using Cranelift. It's designed to be embedded in applications, and provides a safe, fast, and configurable environment for executing WebAssembly code.

Key features of Wasmtime-go:
1. Fast JIT compilation of WebAssembly code
2. Support for the WebAssembly System Interface (WASI)
3. Memory safety and sandboxing
4. Configurable resource limits
5. Support for multiple instances of the module

## Getting Started

TODO: Add instructions for building and running the system.
For local testing and development:
- TCP can be used as a fallback communication method — either on host or in docker container.
- docker does not support vsock
- QEMU can be used to emulate vsock communication, with executor running in the VM, and manager running on the host.

## Generate contracts bindings
The interaction with the contracts on chain is managed using a Geth tool called `abigen`, that can convert contract code into Go code that can be used directly in Go applications. For more information, see [Go contract binding](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings-v2). 

For running `abigen`, check to have `solc` installed:
```bash
solc --version
```
Version 0.8.30 must be used.
In order to install the proper version of `solc` we suggest to use [solc-select](https://github.com/crytic/solc-select) (a tool to quickly switch between Solidity compiler versions):

```bash
pip3 install solc-select
solc-select use 0.8.30 --always-install
```

For installing `abigen`:
```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@v1.16.2
```
Before running `abigen`, install contracts dependencies:
```bash
cd contracts
npm install
```
Go bindings are generated using the `go generate` command, which in turn invokes the `abigen` tool. The resulting files are committed to this repository to make them available for dependent packages. If you modify the contracts, you must regenerate these files.

Remove existing generated files:

```bash
    rm -fr pkg/blockchain/contracts/*
```

Run the Go code generator:

```bash
    go generate ./...
```

Finally, commit the newly generated files to the repository.

## Testing

For information about testing the system, see [README_TESTS.md](README_TESTS.md).
