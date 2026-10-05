# WebAssembly Runtime with Wasmtime-Go

Execute WebAssembly modules using Wasmtime-Go for runtime execution.

TODO: Add an example minimal WASM application

## Module and instance lifecycle

`WasmtimeRuntime` caches per application the compiled module and the named pipe that forwards the guest stdout/stderr to the host log (LRU, `EXECUTOR_MAX_CACHED_MODULES`). It does not keep guest instances:

- every guest call (`Deposit`, `ProcessRequest`, the memory stats) runs on a fresh store and instance, after calling `load_module` on it; the instance is discarded when the call ends. `Deploy` and `LoadModule` also use a temporary instance, and cache the compiled module.
- guests therefore must not rely on state kept in globals between calls: the application state is passed in every call.
- wasmtime-go v1.0.0 frees a store (its linear memory and WASI file descriptors) only in its finalizer, so the runtime forces a GC after each call.
- guest calls are serialized by `execLock`, also against `UnloadModule` and `Close`. Lock ordering: `execLock` before `moduleLock`.

Fresh instances are needed because TinyGo guests (<= 0.42) never free large allocations (the GC scans the guest data section conservatively): a long-lived instance only grows until it traps. Within a single call, the guest heap still grows to several times the state size.

## Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [Wasmtime Documentation](https://docs.wasmtime.dev/)
- [WebAssembly Specification](https://webassembly.org/specs/)
- [WASI Interface](https://wasi.dev/)