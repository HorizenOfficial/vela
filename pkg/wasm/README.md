# WebAssembly Runtime with Wasmtime-Go

Execute WebAssembly modules using Wasmtime-Go for runtime execution.

TODO: Add an example minimal WASM application

## Module and instance lifecycle

`WasmtimeRuntime` caches, per application, the compiled module and the named pipe that forwards the guest stdout/stderr to the host log (LRU, `EXECUTOR_MAX_CACHED_MODULES`). It does not keep guest instances:

- every guest call (`Deploy`, `Deposit`, `ProcessRequest`, the memory stats) runs on a fresh store and instance, created by `getOrLoadModule` (or `instantiate` for `Deploy`) and released with `releaseInstance` when the call ends, after its deallocates. No guest initialization function runs on a new instance.
- each store gets the guest memory cap and table limits (`newModuleStore`) and its own epoch deadline, so a cap lowered at runtime applies from the next call.
- guests therefore must not rely on state kept in globals between calls: the application state is passed in every call.
- guest calls are serialized by `execLock`, also against eviction, `UnloadModule` and `Close`, which close the app's log pipe. Lock ordering: `execLock` before `moduleLock`.
- a guest fault still evicts the compiled module (`evictFaultedModule`). That is conservative: the faulted instance is discarded like every other.

Fresh instances are needed because TinyGo guests never free large allocations: a long-lived instance only grows until it hits the memory cap. Within a single call, the guest heap still grows to several times the state size.

## Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [Wasmtime Documentation](https://docs.wasmtime.dev/)
- [WebAssembly Specification](https://webassembly.org/specs/)
- [WASI Interface](https://wasi.dev/)