# WASM Go Module: Simple App

This repository contains the Go implementation of the **Simple App** WASM module, an example application that manages a simple counter.

## Building

To build the WASM module, you can use the provided Makefile:

```bash
make build
```

This command uses TinyGo to compile the `main.go` file into a WASM module. The underlying command is:
```bash
tinygo build -o build/simple_app.wasm -target wasi .
```

This will create the `build/simple_app.wasm` file.
