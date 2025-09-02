# WebAssembly Runtime with TinyGo and Wasmtime-Go

Build and execute WebAssembly modules using TinyGo for compilation and Wasmtime-Go for runtime execution.

## Overview

The WebAssembly runtime system consists of:
- **TinyGo**: Compiles Go code to WebAssembly modules
- **Wasmtime-Go**: Executes WebAssembly modules in a secure, sandboxed environment
- **Payment App**: Example WASM application for handling deposits, transfers, and withdrawals

## Prerequisites

### Install TinyGo

**macOS (using Homebrew):**
```bash
brew tap tinygo-org/tools
brew install tinygo
```

**Linux/Manual Installation:**
```bash
wget https://github.com/tinygo-org/tinygo/releases/download/v0.39.0/tinygo_0.39.0_amd64.deb
sudo dpkg -i tinygo_0.39.0_amd64.deb
```

**Verify Installation:**
```bash
tinygo version
```

## Build WebAssembly Payment App Module

Navigate to the module directory:
```bash
cd pkg/wasm/wasm-go/
```

Build the WebAssembly module:
```bash
tinygo build -o payment_app.wasm -target wasi main.go
```

For production builds, add optimization flags:
```bash
tinygo build -o payment_app.wasm -target wasi -opt=2 main.go
```

## Development Workflow

### 1. Modify WASM Module 
`pkg/wasm/wasm-go/main.go` contains the core payment application logic.
`pkg/wasm/wasm-go/utils/utils.go` provides utility functions and memory management. Could be separate Go module in the future for better organization and reusability.

### 2. Rebuild Module
```bash
cd pkg/wasm/wasm-go/
tinygo build -o payment_app.wasm -target wasi main.go
```

### 3. Update Tests
Add corresponding tests in `wasmtime_runtime_test.go`.

### 4. Verify Changes
```bash
go test ./pkg/wasm/... -v
```

## Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [Wasmtime Documentation](https://docs.wasmtime.dev/)
- [WebAssembly Specification](https://webassembly.org/specs/)
- [WASI Interface](https://wasi.dev/)