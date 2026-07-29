package wasm

import (
	"bufio"
	"container/list"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	"github.com/HorizenOfficial/vela/pkg/logger"
	appCommon "github.com/HorizenOfficial/vela/pkg/wasm/common"
	wasmtime "github.com/bytecodealliance/wasmtime-go/v47"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

// Address is a local definition of a 20-byte address.
type Address [20]byte

// maxGuestMemoryCeilingBytes is both the default and the maximum allowed value
// for the per-guest linear memory cap. The runtime refuses growth beyond the
// cap (the guest sees a failed memory.grow), protecting the enclave's fixed
// RAM from a buggy or malicious app. The 2 GiB ceiling is a hard ABI
// constraint: the host exchanges guest pointers as signed int32 offsets (see
// writeToMemory and extractResultBytes), so no guest offset may reach 2 GiB.
const maxGuestMemoryCeilingBytes = 2 * 1024 * 1024 * 1024

// Per-store resource limits passed to wasmtime's store limiter.
const (
	// limiterKeepDefault is the sentinel the limiter accepts for "leave the
	// engine default in place for this limit".
	limiterKeepDefault = -1

	// maxInstancesPerStore and maxMemoriesPerStore pin the counts the guest
	// memory cap is accounted against. The cap is enforced per linear memory,
	// not per store, so without these a guest declaring N memories could use N
	// times the cap and break the enclave RAM budget documented in
	// pkg/executor.Config (MaxCachedModules * MaxGuestMemoryBytes). Each store
	// holds exactly one instance with its single exported memory.
	maxInstancesPerStore = 1
	maxMemoriesPerStore  = 1
)

// errGuestFault is a sentinel, used with errors.Is, that distinguishes failures
// caused by the guest itself — a trap, or an allocation it could not satisfy —
// from host-side failures such as a nil store or a missing export. The three
// failure classes and their effect on the cache are summarised in
// docs/design/WASM_HOST_ABI.md. Only the
// former mean the guest's heap is in an unknown state and the module must be
// dropped from the cache (see evictFaultedModule); recompiling a module that is
// simply missing an export would repeat on every request forever.
//
// The test is whether guest code actually ran, NOT whether the failure looks
// deterministic. A missing export, an export whose signature does not match the
// host ABI (see isGuestTrap), or a nil store all mean nothing executed and no heap
// was touched, so a fresh instance provably cannot behave differently — those stay
// cached.
//
// A result the host cannot parse (see the json.Unmarshal sites) is the ambiguous
// case, and is deliberately treated as a fault even though a broken serializer
// would fail identically on every request: the guest did run, and did write that
// buffer, so the host cannot distinguish a deterministic serialization defect
// (where the recompile is wasted) from a heap bug that clobbered the result
// (where reusing the instance is dangerous). The costs are asymmetric — being
// wrong here spends one recompile per request on an app that is already failing
// every request, whereas being wrong the other way lets the next request run on a
// damaged heap and possibly commit a plausible-but-incorrect state root on-chain.
var errGuestFault = errors.New("guest fault")

// guestFaultError tags an error as a guest fault while leaving its message
// untouched. The message is deliberately preserved byte for byte: it reaches the
// signed, on-chain ErrorMsg (see buildErrorPayload in pkg/executor), so an
// internal reclassification must not change what gets signed.
type guestFaultError struct{ err error }

func (e *guestFaultError) Error() string        { return e.err.Error() }
func (e *guestFaultError) Unwrap() error        { return e.err }
func (e *guestFaultError) Is(target error) bool { return target == errGuestFault }

// asGuestFault marks err as caused by the guest. See errGuestFault.
func asGuestFault(err error) error { return &guestFaultError{err: err} }

// isGuestTrap reports whether an error from wasmtime Func.Call means the guest
// actually ran and faulted, rather than the host failing to invoke it at all.
//
// Call validates the argument list before entering wasm and rejects a mismatched
// arity or value kind up front ("too many arguments provided", "integer provided
// for non-integer argument", or an error from wasmtime_func_call), so a module
// whose export does not match the host ABI fails here without executing a single
// instruction. That is a static defect, in the same class as a missing export: the
// module stays cached, because recompiling cannot change its signature. Only a
// genuine trap — surfaced as *wasmtime.Trap — means the heap is in an unknown
// state. See errGuestFault.
//
// Note that only the deploy export's signature is exercised at deploy time, so an
// app built against a stale host ABI can deploy successfully and then fail this
// check on every request.
func isGuestTrap(err error) bool {
	var trap *wasmtime.Trap
	return errors.As(err, &trap)
}

// newModuleStore creates a per-module store with the guest memory cap applied.
// Must be called with moduleLock held (write lock), which also guards
// maxGuestMemoryBytes.
func (r *WasmtimeRuntime) newModuleStore() *wasmtime.Store {
	store := wasmtime.NewStore(r.engine)
	// Signature: Limiter(memorySize, tableElements, instances, tables, memories).
	// Multi-memory is also disabled at the engine level (see newPinnedEngine);
	// maxMemoriesPerStore is the second half of that defence.
	store.Limiter(
		r.maxGuestMemoryBytes,
		limiterKeepDefault,
		maxInstancesPerStore,
		limiterKeepDefault,
		maxMemoriesPerStore,
	)
	return store
}

// ApplicationModule contains the compiled module, its instantiated instance,
// the module-specific store, the exported memory, and convenience handles.
// Note: a Store represents module execution state and is not safe sharing one store across modules
// and concurrent requests because that can lead to races or panics.
type ApplicationModule struct {
	store      *wasmtime.Store
	module     *wasmtime.Module
	instance   *wasmtime.Instance
	memory     *wasmtime.Memory
	deallocate *wasmtime.Func
	cleanupFds func()
	closeOnce  sync.Once
}

// Close releases all resources associated with the ApplicationModule.
// Safe to call multiple times concurrently; cleanup runs exactly once.
func (m *ApplicationModule) Close() {
	m.closeOnce.Do(func() {
		if m.cleanupFds != nil {
			m.cleanupFds()
		}

		// Explicitly release the native wasmtime state instead of waiting for a
		// Go GC finalizer pass: the Go GC does not see the weight of the guest
		// linear memory held by the store, so relying on finalizers can keep
		// hundreds of MB allocated inside the enclave long after eviction.
		// Closing the store frees the instance, its memory and funcs; closing
		// the module frees the compiled code. Both are safe to call once here
		// because the runtime never uses an ApplicationModule after Close.
		if m.store != nil {
			m.store.Close()
		}
		if m.module != nil {
			m.module.Close()
		}

		// Clear references to WASM resources
		m.instance = nil
		m.module = nil
		m.memory = nil
		m.store = nil
		m.deallocate = nil
	})
}

// WasmtimeRuntime implements the Runtime interface using wasmtime-go
type WasmtimeRuntime struct {
	engine           *wasmtime.Engine
	modules          map[common.ApplicationIdType]*ApplicationModule // Map of application ID to module
	moduleLock       sync.RWMutex                                    // Lock for module access
	log              logger.Logger
	maxCachedModules int                                        // 0 = unlimited
	accessOrder      *list.List                                 // LRU order: front = most recent, back = least recent
	accessElements   map[common.ApplicationIdType]*list.Element // O(1) lookup into accessOrder
	// maxGuestMemoryBytes caps the linear memory of each guest store created
	// from now on; always in (0, maxGuestMemoryCeilingBytes]. Guarded by moduleLock.
	maxGuestMemoryBytes int64

	// execLock serializes guest execution against module teardown. It is held for
	// the whole duration of every operation that runs guest code or closes a
	// module, so a store can never be freed while another goroutine is calling
	// into it.
	//
	// moduleLock alone is not sufficient: Deposit and ProcessRequest look the
	// module up under moduleLock but release it before calling into the guest, so
	// a concurrent load for a different appId could pick that same module as the
	// LRU eviction victim and close its store mid-call. Since Close frees the
	// native wasmtime state immediately rather than at a GC finalizer pass, that
	// is a use-after-free — a native segfault inside the enclave.
	//
	// The transport dispatches every inbound message in its own goroutine (see
	// pkg/communication/shared_impl.go), so the only thing preventing this today
	// is that the manager happens to submit one request at a time — an invariant
	// no code enforces. This makes the serialization explicit, and costs nothing
	// while execution is serial anyway.
	//
	// Lock ordering: acquire execLock BEFORE moduleLock, never the reverse, and
	// never acquire it from a function that already holds moduleLock. Holding it
	// across a guest call cannot deadlock because guest code cannot re-enter the
	// runtime: only WASI is defined on the linker, there are no custom host funcs.
	//
	// TODO: when per-goroutine instances land (see the TODOs in
	// app/simple/integration_test.go), replace this with refcounting the live
	// calls per module and deferring Close until the last caller finishes, so
	// that independent apps can execute in parallel.
	execLock sync.Mutex
}

// newPinnedEngine builds the wasmtime engine with an explicitly pinned WASM
// feature set instead of inheriting wasmtime's defaults.
//
// Rationale: the executor re-executes guest code inside the enclave and the
// resulting state is committed on-chain, so the set of accepted WASM features
// is part of this system's observable behaviour. Tracking upstream defaults
// would silently widen it on every dependency bump (the modern default set is
// much wider than the v1.0 one this replaced) and would admit features whose
// results are host-dependent. Every WASM *feature* knob the wasmtime-go API
// exposes is therefore set here explicitly: adding a proposal must be a
// deliberate, reviewed change.
//
// Caveat on "every": the C API carries a few proposal flags that wasmtime-go does
// not wrap, so they cannot be pinned from Go and are left at their upstream
// defaults — as of v47 those are branch hinting (documented false by default),
// custom page sizes, stack switching, and the component-model sub-flags. Check
// this list again when bumping, in case the bindings start exposing them.
//
// TODO: bound guest execution. The execution-bounding knobs (SetConsumeFuel,
// SetEpochInterruption, SetMaxWasmStack) are all left at their defaults, i.e.
// guest execution is currently unbounded — called out explicitly so this is not
// mistaken for a decision that they were considered and rejected. A guest in an
// infinite loop runs forever, and because execLock is held across guest calls it
// blocks every later request for every app until the enclave is restarted; Close
// only avoids hanging shutdown by giving up after shutdownExecLockTimeout, and
// the trap-eviction path cannot help because a runaway guest never traps. Add
// SetEpochInterruption here plus a watchdog incrementing the epoch, and arm a
// per-call deadline on each store. Tracked as a separate task (epoch
// interruption, with host-enforced fuel depending on it); it is a prerequisite
// for those rather than an enhancement.
//
// Enabled features are the ones TinyGo emits, all deterministic per spec.
// Disabled ones are either host-dependent (relaxed SIMD), incompatible with the
// host ABI (memory64 widens guest pointers past the int32 offsets used by
// writeToMemory/extractResultBytes), incompatible with the enclave RAM budget
// (multi-memory, threads' shared memories), or simply unused surface.
//
// NOTE: when bumping wasmtime, check the release notes for newly added
// proposals and pin them here too.
func newPinnedEngine() *wasmtime.Engine {
	config := wasmtime.NewConfig()

	// --- Enabled: required by TinyGo-generated modules, deterministic ---
	config.SetWasmBulkMemory(true)     // memory.copy / memory.fill
	config.SetWasmMultiValue(true)     // multi-result functions
	config.SetWasmReferenceTypes(true) // funcref tables behind call_indirect
	config.SetWasmSIMD(true)           // fixed-width SIMD: results are spec-defined

	// --- Disabled: host-dependent results ---
	// Relaxed SIMD instructions are explicitly allowed to differ between host
	// architectures, which would make guest results depend on the machine the
	// enclave runs on. If it is ever enabled, deterministic mode is mandatory.
	config.SetWasmRelaxedSIMD(false)

	// --- Disabled: incompatible with the host ABI or the RAM budget ---
	config.SetWasmMemory64(false)    // guest pointers must stay in int32 range
	config.SetWasmMultiMemory(false) // one memory per store, see newModuleStore
	config.SetWasmThreads(false)     // shared memories escape the per-store cap

	// --- Disabled: unused surface ---
	config.SetWasmTailCall(false)
	config.SetWasmFunctionReferences(false)
	config.SetWasmGC(false)
	config.SetWasmWideArithmetic(false)
	config.SetWasmExceptions(false)
	config.SetWasmComponentModel(false)

	// --- Disabled: whole runtime subsystems we do not use ---
	// Not WASM proposals but engine-wide capabilities. GC support is switched off
	// entirely (the GC proposal above is already off); note that this would also
	// break a guest passing externref values around, which none do — the host
	// defines no host functions, only WASI, so there are no host references for a
	// guest to hold. Re-enable both if that ever changes.
	config.SetGCSupport(false)
	config.SetConcurrencySupport(false)

	// Replace NaN payloads with a single canonical value. Not required by the
	// WASM spec and off by default, but NaN bit patterns are otherwise
	// host-dependent, which is the same reproducibility concern as relaxed SIMD.
	// Costs a few percent on float-heavy guests; the guests here are state
	// machines, so the trade is worth it.
	config.SetCraneliftNanCanonicalization(true)

	return wasmtime.NewEngineWithConfig(config)
}

// NewWasmtimeRuntime creates a new wasmtime runtime instance.
// maxCachedModules controls the LRU cache size: 0 means unlimited.
func NewWasmtimeRuntime(log logger.Logger, maxCachedModules int) *WasmtimeRuntime {
	log.Info("Runtime: Initializing wasmtime runtime (maxCachedModules=%d)", maxCachedModules)

	// Create the engine with an explicitly pinned WASM feature set
	engine := newPinnedEngine()

	return &WasmtimeRuntime{
		engine:              engine,
		modules:             make(map[common.ApplicationIdType]*ApplicationModule),
		log:                 log,
		maxCachedModules:    maxCachedModules,
		accessOrder:         list.New(),
		accessElements:      make(map[common.ApplicationIdType]*list.Element),
		maxGuestMemoryBytes: maxGuestMemoryCeilingBytes,
	}
}

// -----------------------------
// Helpers
// -----------------------------

// toInt32 safely converts wasmtime Func.Call() return interfaces to int32.
// This is a defensive tecnique that can be useful for futures ABI evolution, for instance
// if adding structured results in the future (e.g., returning both pointer and length)
func toInt32(v interface{}) (int32, error) {
	switch n := v.(type) {
	case int32:
		return n, nil
	case int64:
		if n < -0x80000000 || n > 0x7fffffff {
			return 0, fmt.Errorf("int64 return out of int32 range: %d", n)
		}
		return int32(n), nil
	case uint32:
		if n > 0x7fffffff {
			return 0, fmt.Errorf("uint32 return out of int32 signed range: %d", n)
		}
		return int32(n), nil
	case uint64:
		if n > 0x7fffffff {
			return 0, fmt.Errorf("uint64 return out of int32 signed range: %d", n)
		}
		return int32(n), nil
	case int:
		return int32(n), nil
	case nil:
		return 0, fmt.Errorf("wasm call returned nil")
	default:
		return 0, fmt.Errorf("unexpected return type from wasm call: %T", v)
	}
}

// -----------------------------
// Memory helpers (host-side)
// -----------------------------

// Note: we are safe using int and uintptr for pointer arithmetic inside TinyGo WASM because it is
// a 32-bit address space (tinygo does not support WASM64) and int is defined to be at least int32 bits in size.

// writeToMemory allocates space in the wasm guest via its allocate export,
// writes the provided data into linear memory, and returns the offset (int32).
// Caller must free the returned pointer (offset) by calling the guest deallocate
// function when finished. Many call sites in this runtime use `defer` to ensure that.
func (r *WasmtimeRuntime) writeToMemory(module *ApplicationModule, data []byte) (int32, error) {
	if module == nil {
		return 0, fmt.Errorf("module is nil")
	}
	if module.memory == nil {
		return 0, fmt.Errorf("memory not initialized")
	}

	if len(data) == 0 {
		// Use 0 to represent a nil/empty pointer (caller must handle zero-length specially if needed).
		return 0, nil
	}

	if module.store == nil {
		return 0, fmt.Errorf("wasm module has a nil store")
	}

	allocateFunc := module.instance.GetFunc(module.store, "allocate")
	if allocateFunc == nil {
		return 0, fmt.Errorf("allocate function not found in WASM module")
	}

	// Call allocate in the module's store context
	rawRes, err := allocateFunc.Call(module.store, int32(len(data)))
	if err != nil {
		return 0, asGuestFault(fmt.Errorf("failed to call allocate: %w", err))
	}

	ptr, err := toInt32(rawRes)
	if err != nil {
		return 0, asGuestFault(fmt.Errorf("allocate returned unexpected value: %w", err))
	}
	if ptr == 0 {
		return 0, asGuestFault(fmt.Errorf("allocate returned null pointer (0)"))
	}

	// Safely get a snapshot of the memory data and perform bounds checks
	memData := module.memory.UnsafeData(module.store)
	// note: guest memory is capped at <= maxGuestMemoryCeilingBytes (2 GiB) via the
	// store's Limiter (see newModuleStore), so every guest offset fits in a positive
	// int32 and we can use the returned ptr directly without an unsigned cast.
	// (The wasmtime-go v1 limitation where UnsafeData panicked when memory reached
	// 4 GiB is gone — v42 onwards builds the slice with unsafe.Slice — the bound below is
	// the deliberate cap, not an API workaround.)

	// Convert the signed int32 pointer to its true unsigned 32-bit value (uint32)
	// This correctly interprets the bits of an address as a large positive number.
	// This cast would not actually be necessary, see above.
	uPtr := uint32(ptr)

	// We cast to uint64 for the calculation to prevent overflow,
	// as len(memData) can be up to 4GB (or larger, depending on the host's slice).
	memLen := uint64(len(memData))
	allocEnd := uint64(uPtr) + uint64(len(data)) // The end address of the requested block

	if uint64(uPtr) >= memLen || allocEnd > memLen {
		return 0, asGuestFault(fmt.Errorf("invalid memory allocation: requested block [%d - %d) exceeds memory size %d", uPtr, allocEnd, memLen))
	}

	copy(memData[uPtr:uPtr+uint32(len(data))], data)

	return ptr, nil
}

// -----------------------------
// Module loading / lifecycle
// -----------------------------

// TODO - this method most of times does not use the wasm bytes, because it will find the module already loaded, unless the executor has
// undergo a restart, and the cache has to be populated again. This caching is also sub-optimal because every time the wasm bytes are passed along
// without being used.
// One approach could be that we use just the appId, and the executor asks to the manager for wasm bytes (that are stored in the dblayer) if it
// does not find them in the cache.
func (r *WasmtimeRuntime) getOrLoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) (*ApplicationModule, error) {
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	if module, exists := r.modules[appId]; exists {
		r.touchModule(appId)
		return module, nil
	}

	// If not loaded, load the module
	_, _, err := r.loadModuleUnlocked(ctx, appId, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to load module: %w", err)
	}

	// Return loaded module
	if module, exists := r.modules[appId]; exists {
		return module, nil
	}
	return nil, fmt.Errorf("module not found after loading: %d", appId)
}

// LoadModule loads a WASM module and returns initial state and state root
func (r *WasmtimeRuntime) LoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, *big.Int, error) {
	// Runs guest code (load_module) and may evict/close other modules.
	r.execLock.Lock()
	defer r.execLock.Unlock()

	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	return r.loadModuleUnlocked(ctx, appId, wasm)
}

// loadModuleUnlocked contains the core logic for loading a module, but without locking.
// This method should only be called when a lock is already held.
func (r *WasmtimeRuntime) loadModuleUnlocked(ctx context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, *big.Int, error) {
	r.log.Info("Wasmtime Runtime: Loading WASM module for application %d (wasm size: %d bytes)", appId, len(wasm))
	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, big.NewInt(0), err
	}

	// Compile the WASM module
	module, err := wasmtime.NewModule(r.engine, wasm)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to compile WASM module: %w", err)
	}

	// Create a per-module store
	store := r.newModuleStore()

	// Registered as early as possible — immediately after the store exists and
	// BEFORE configureWasiLogPipes — so that every failure path from here on
	// releases the native state, including a log-pipe failure (mkfifo/open),
	// which would otherwise fall back to a finalizer. cleanupLogPipes is declared
	// up front and nil-checked, since it is only assigned below.
	success := false
	var cleanupLogPipes func()
	defer func() {
		if success {
			return
		}
		if cleanupLogPipes != nil {
			cleanupLogPipes()
		}
		// Release the native state explicitly on every failure path, for the same
		// reason ApplicationModule.Close does it on eviction: the Go GC does not
		// see the guest linear memory held by the store, so leaving it to a
		// finalizer can keep hundreds of MB alive inside the enclave. Nothing else
		// can close these — the module never reaches r.modules on these paths.
		// Registered before any deallocate defers below, so LIFO runs this last,
		// after those have finished using the store.
		store.Close()
		module.Close()
	}()

	cleanupLogPipes, err = r.configureWasiLogPipes(ctx, appId, store)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to configure WASI log pipes: %w", err)
	}

	// Create WASI configuration and linker for TinyGo WASI imports
	linker := wasmtime.NewLinker(r.engine)
	err = linker.DefineWasi()
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to define WASI: %w", err)
	}

	// Instantiate the module using the module-specific store
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	// Get the memory export.
	// Note: the WASM module is responsible for defining and exporting its own memory, this is done by tinygo,
	// check for instance WAT (WebAssembly text format) generated via wasm2wat tool
	memoryExport := instance.GetExport(store, "memory")
	if memoryExport == nil {
		return nil, big.NewInt(0), fmt.Errorf("memory export not found in WASM module")
	}

	memory := memoryExport.Memory()
	if memory == nil {
		return nil, big.NewInt(0), fmt.Errorf("memory export is not a memory")
	}

	// Get the load_module function
	loadModuleFunc := instance.GetFunc(store, "load_module")
	if loadModuleFunc == nil {
		return nil, big.NewInt(0), fmt.Errorf("load_module function not found in WASM module")
	}

	// Get the deallocate function (optional, but recommended for memory management)
	deallocateFunc := instance.GetFunc(store, "deallocate")

	// Call the load_module function
	// Wasm supports only int64, so we cast appId to int64
	result, err := loadModuleFunc.Call(store, wasmAppId)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to call load_module: %w", err)
	}

	appModule := &ApplicationModule{
		store:      store,
		module:     module,
		instance:   instance,
		memory:     memory,
		deallocate: deallocateFunc,
		cleanupFds: cleanupLogPipes,
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	r.log.Debug("Wasmtime Runtime: Raw result from WASM: %s", string(resultBytes))

	// Deserialize the result
	var loadResult appCommon.LoadModuleResult
	if err := json.Unmarshal(resultBytes, &loadResult); err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to unmarshal load module result: %w", err)
	}

	if loadResult.Error != "" {
		return nil, big.NewInt(0), fmt.Errorf("failed to load module: %s", loadResult.Error)
	}

	success = true // Disables the deferred cleanup

	// if a module already exists for this appId, clean it up before overwriting
	if oldModule, exists := r.modules[appId]; exists {
		r.cleanupModule(appId, oldModule)
	}

	// Store the module in the runtime registry and update LRU
	r.modules[appId] = appModule
	r.touchModule(appId)
	r.evictIfNeeded()
	r.log.Info("Wasmtime Runtime: Successfully loaded WASM module for application %d", appId)

	return loadResult.State, loadResult.Fuel.ToInt(), nil
}

// Deploy loads a WASM module and initializes it with constructor parameters.
// The guest export is deploy(appId, paramsPtr, paramsLen) -> resultPtr.
//
// This method is the deploy-time entry point: it compiles the module, calls the guest's
// deploy function with constructor params, and returns the initial application state.
//
// NOTE on LoadModule vs Deploy:
// Deploy is used at deploy time to initialize the app state with constructor parameters.
// LoadModule is no longer used for deploy-time initialization but is retained because
// getOrLoadModule (used by Deposit/ProcessRequest) depends on loadModuleUnlocked to
// compile and cache modules for subsequent requests.
//
// TODO: Refactor to extract shared module setup (compile, instantiate, WASI config) into
// a compileAndInstantiate helper. Currently deployUnlocked and loadModuleUnlocked duplicate
// ~60 lines of identical boilerplate. This would also let getOrLoadModule compile and
// instantiate without calling any guest function, eliminating the unnecessary load_module
// guest call during cache warm-up.
func (r *WasmtimeRuntime) Deploy(ctx context.Context, appId common.ApplicationIdType, constructorParams []byte, wasm []byte) ([]byte, *big.Int, error) {
	// Runs guest code (deploy) and may evict/close other modules.
	r.execLock.Lock()
	defer r.execLock.Unlock()

	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	return r.deployUnlocked(ctx, appId, constructorParams, wasm)
}

// deployUnlocked contains the core logic for deploying a module with constructor params, without locking.
// This method should only be called when a lock is already held.
func (r *WasmtimeRuntime) deployUnlocked(ctx context.Context, appId common.ApplicationIdType, constructorParams []byte, wasm []byte) ([]byte, *big.Int, error) {
	r.log.Info("Wasmtime Runtime: Deploying WASM module for application %d (wasm size: %d bytes, params size: %d bytes)", appId, len(wasm), len(constructorParams))
	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, big.NewInt(0), err
	}

	// Compile the WASM module
	module, err := wasmtime.NewModule(r.engine, wasm)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to compile WASM module: %w", err)
	}

	// Create a per-module store
	store := r.newModuleStore()

	// Registered as early as possible — immediately after the store exists and
	// BEFORE configureWasiLogPipes — so that every failure path from here on
	// releases the native state, including a log-pipe failure (mkfifo/open),
	// which would otherwise fall back to a finalizer. cleanupLogPipes is declared
	// up front and nil-checked, since it is only assigned below.
	success := false
	var cleanupLogPipes func()
	defer func() {
		if success {
			return
		}
		if cleanupLogPipes != nil {
			cleanupLogPipes()
		}
		// Release the native state explicitly on every failure path, for the same
		// reason ApplicationModule.Close does it on eviction: the Go GC does not
		// see the guest linear memory held by the store, so leaving it to a
		// finalizer can keep hundreds of MB alive inside the enclave. Nothing else
		// can close these — the module never reaches r.modules on these paths.
		// Registered before any deallocate defers below, so LIFO runs this last,
		// after those have finished using the store.
		store.Close()
		module.Close()
	}()

	cleanupLogPipes, err = r.configureWasiLogPipes(ctx, appId, store)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to configure WASI log pipes: %w", err)
	}

	// Create WASI configuration and linker for TinyGo WASI imports
	linker := wasmtime.NewLinker(r.engine)
	err = linker.DefineWasi()
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to define WASI: %w", err)
	}

	// Instantiate the module using the module-specific store
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	// Get the memory export
	memoryExport := instance.GetExport(store, "memory")
	if memoryExport == nil {
		return nil, big.NewInt(0), fmt.Errorf("memory export not found in WASM module")
	}

	memory := memoryExport.Memory()
	if memory == nil {
		return nil, big.NewInt(0), fmt.Errorf("memory export is not a memory")
	}

	// Get the deploy function
	deployFunc := instance.GetFunc(store, "deploy")
	if deployFunc == nil {
		return nil, big.NewInt(0), fmt.Errorf("deploy function not found in WASM module")
	}

	// Get the deallocate function (optional, but recommended for memory management)
	deallocateFunc := instance.GetFunc(store, "deallocate")

	appModule := &ApplicationModule{
		store:      store,
		module:     module,
		instance:   instance,
		memory:     memory,
		deallocate: deallocateFunc,
		cleanupFds: cleanupLogPipes,
	}

	// Write constructor params to guest memory
	if constructorParams == nil {
		constructorParams = []byte{}
	}
	paramsPtr, err := r.writeToMemory(appModule, constructorParams)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to write constructor params to memory: %w", err)
	}
	if appModule.deallocate != nil && paramsPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, paramsPtr, int32(len(constructorParams))) }()
	}

	// Guest ABI: deploy(appId, paramsPtr, paramsLen)
	result, err := deployFunc.Call(store, wasmAppId, paramsPtr, int32(len(constructorParams)))
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to call deploy: %w", err)
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to extract deploy result bytes: %w", err)
	}

	r.log.Debug("Wasmtime Runtime: Raw result from WASM deploy: %s", string(resultBytes))

	// Deserialize the result
	var deployResult appCommon.DeployResult
	if err := json.Unmarshal(resultBytes, &deployResult); err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to unmarshal deploy result: %w", err)
	}

	if deployResult.Error != "" {
		return nil, big.NewInt(0), fmt.Errorf("failed to deploy module: %s", deployResult.Error)
	}

	// A module for this appId should not exist at deploy time. If it does,
	// it indicates a duplicate deploy or an unexpected state. Checked BEFORE
	// success is set, so this path still runs the deferred cleanup: the module
	// never enters r.modules, so nothing else would ever close its log FIFOs
	// (leaking the fds and the log-forwarding goroutine).
	if _, exists := r.modules[appId]; exists {
		return nil, big.NewInt(0), fmt.Errorf("application %d is already deployed", appId)
	}

	success = true // Disables the deferred cleanup

	// Store the module in the runtime registry and update LRU
	r.modules[appId] = appModule
	r.touchModule(appId)
	r.evictIfNeeded()
	r.log.Info("Wasmtime Runtime: Successfully deployed WASM module for application %d", appId)

	return deployResult.State, deployResult.Fuel.ToInt(), nil
}

// -----------------------------
// Public operations (Deposit, ProcessRequest, ...)
//
// Each of these ensures that any temporary allocations created by writeToMemory
// are deallocated after the guest call by using defer cleanup.
// -----------------------------

// Deposit processes a deposit with token awareness (tokenAddress = 0x0 for ETH)
func (r *WasmtimeRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, tokenAddress ethCommon.Address, depositAmount *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, *big.Int, *apperrors.RequestFailure) {
	r.log.Info("Wasmtime Runtime: Processing deposit for application %d (token: %v, value: %v wei for sender: %v)", appId, tokenAddress, depositAmount, sender)

	// Held across the guest call: the module is fetched under moduleLock but
	// executed after releasing it, so only execLock keeps a concurrent eviction
	// from closing this store mid-call.
	r.execLock.Lock()
	defer r.execLock.Unlock()

	if depositAmount == nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, "value cannot be nil")
	}
	if depositAmount.Sign() < 0 {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, "value cannot be negative")
	}

	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, fmt.Sprintf("invalid application id: %v", err))
	}

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, fmt.Sprintf("failed to get or load module: %v", err))
	}

	// Drop the module if the guest faults: its heap is left in an unknown state
	// and must not serve the next request (see evictFaultedModule). Registered
	// BEFORE the deallocate defers below so that, defers being LIFO, the store is
	// closed only after those have run — closing it first would leave them
	// calling into freed memory.
	guestFaulted := false
	defer func() {
		if guestFaulted {
			r.evictFaultedModule(appId)
		}
	}()

	// Get the deposit function
	depositFunc := appModule.instance.GetFunc(appModule.store, "deposit")
	if depositFunc == nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFunctionNotFound, "deposit function not found in WASM module")
	}

	senderBytes := sender.Bytes()
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write sender to memory: %v", err))
	}
	if appModule.deallocate != nil && senderPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, senderPtr, int32(len(senderBytes))) }()
	}

	tokenBytes := tokenAddress.Bytes()
	tokenPtr, err := r.writeToMemory(appModule, tokenBytes)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write token address to memory: %v", err))
	}
	if appModule.deallocate != nil && tokenPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, tokenPtr, int32(len(tokenBytes))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write state to memory: %v", err))
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	valueBytes := depositAmount.Bytes()
	valuePtr, err := r.writeToMemory(appModule, valueBytes)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write value to memory: %v", err))
	}
	if appModule.deallocate != nil && valuePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, valuePtr, int32(len(valueBytes))) }()
	}

	// Guest ABI: deposit(appId, senderPtr, senderLen, tokenPtr, tokenLen, valuePtr, valueLen, statePtr, stateLen)
	result, err := depositFunc.Call(appModule.store, wasmAppId, senderPtr, int32(len(senderBytes)), tokenPtr, int32(len(tokenBytes)), valuePtr, int32(len(valueBytes)), statePtr, int32(len(state)))
	if err != nil {
		// A trap means the guest ran and faulted; an argument-list rejection means
		// it never started (see isGuestTrap).
		guestFaulted = isGuestTrap(err)
		// TODO some standard way of getting errors here?
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeDepositFailed, fmt.Sprintf("failed to call deposit: %v", err))
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest-provided pointer, not a host-side problem
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedExtractingResultBytes, fmt.Sprintf("failed to extract wasm module result bytes: %v", err))
	}

	r.log.Info("Wasmtime Runtime: Raw deposit result from WASM: %s", string(resultBytes))

	// Deserialize the result
	var depositResult appCommon.DepositResult
	if err := json.Unmarshal(resultBytes, &depositResult); err != nil {
		guestFaulted = true // guest produced a result the host cannot parse
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeJsonUnmarshalError, fmt.Sprintf("failed to unmarshal deposit result: %v", err))
	}

	if depositResult.Error != "" {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeDepositFailed, fmt.Sprintf("failed to process deposit: %s", depositResult.Error))
	}

	r.log.Info("Wasmtime Runtime: Successfully processed deposit for sender %s, generated %d events", sender, len(depositResult.Events))

	return depositResult.State, depositResult.Events, depositResult.AppEvents, depositResult.Fuel.ToInt(), nil
}

// ProcessRequest processes a request and returns the new state, events, withdrawals, and optionally a deanonymization report
func (r *WasmtimeRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	r.log.Info("Wasmtime Runtime: Processing request for application %d (type: %s, payload size: %d, state size: %d)", appId, requestType, len(payload), len(state))

	// Held across the guest call: the module is fetched under moduleLock but
	// executed after releasing it, so only execLock keeps a concurrent eviction
	// from closing this store mid-call.
	r.execLock.Lock()
	defer r.execLock.Unlock()

	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, fmt.Sprintf("invalid application id: %v", err))
	}

	if len(payload) == 0 {
		r.log.Warn("Wasmtime Runtime: Empty payload for application %d, returning current state", appId)
		return state, nil, nil, nil, nil, big.NewInt(0), nil
	}

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedLoadingOrGettingModule, fmt.Sprintf("failed to get or load module: %v", err))
	}

	// Drop the module if the guest faults: its heap is left in an unknown state
	// and must not serve the next request (see evictFaultedModule). Registered
	// BEFORE the deallocate defers below so that, defers being LIFO, the store is
	// closed only after those have run — closing it first would leave them
	// calling into freed memory. Covers both the TrustProcess branch and the
	// process_request path.
	guestFaulted := false
	defer func() {
		if guestFaulted {
			r.evictFaultedModule(appId)
		}
	}()

	// TRUSTPROCESS dispatch (Phase 11.1): TrustProcess requests are routed to the
	// dedicated `trusted_request` WASM export, which takes NEITHER sender (the
	// trusted path never reads it) NOR request_type (implicit for this export).
	// All other request types fall through to `process_request` below, unchanged.
	if requestType == common.TrustProcess {
		trustedFunc := appModule.instance.GetFunc(appModule.store, "trusted_request")
		if trustedFunc == nil {
			return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFunctionNotFound, "trusted_request function not found in WASM module")
		}

		// No sender write (D-05). Only payload + state are passed to the guest.
		payloadPtr, err := r.writeToMemory(appModule, payload)
		if err != nil {
			guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
			return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write payload to memory: %v", err))
		}
		if appModule.deallocate != nil && payloadPtr != 0 {
			defer func() { _, _ = appModule.deallocate.Call(appModule.store, payloadPtr, int32(len(payload))) }()
		}

		statePtr, err := r.writeToMemory(appModule, state)
		if err != nil {
			guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
			return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write state to memory: %v", err))
		}
		if appModule.deallocate != nil && statePtr != 0 {
			defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
		}

		// Guest ABI: trusted_request(appId, payloadPtr, payloadLen, statePtr, stateLen)
		result, err := trustedFunc.Call(appModule.store, wasmAppId, payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
		if err != nil {
			guestFaulted = isGuestTrap(err) // false when the guest never started, see isGuestTrap
			return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("failed to call trusted_request: %v", err))
		}

		resultBytes, err := r.extractResultBytes(result, appModule)
		if err != nil {
			guestFaulted = errors.Is(err, errGuestFault) // guest-provided pointer, not a host-side problem
			return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedExtractingResultBytes, fmt.Sprintf("failed to extract wasm module result bytes: %v", err))
		}

		var processResult appCommon.ProcessResult
		if err := json.Unmarshal(resultBytes, &processResult); err != nil {
			guestFaulted = true // guest produced a result the host cannot parse
			return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeJsonUnmarshalError, fmt.Sprintf("failed to unmarshal process result: %v", err))
		}
		if processResult.Error != "" {
			return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("failed to process request: %s", processResult.Error))
		}

		r.log.Info("Wasmtime Runtime: Successfully processed TrustProcess request for application %d via trusted_request, generated %d events and %d withdrawals", appId, len(processResult.Events), len(processResult.Withdrawals))
		return processResult.State, processResult.Events, processResult.AppEvents, processResult.Withdrawals, processResult.Report, processResult.Fuel.ToInt(), nil
	}

	// Get the process_request function
	processRequestFunc := appModule.instance.GetFunc(appModule.store, "process_request")
	if processRequestFunc == nil {
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFunctionNotFound, "process_request function not found in WASM module")
	}

	senderBytes := sender.Bytes()
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write sender to memory: %v", err))
	}
	if appModule.deallocate != nil && senderPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, senderPtr, int32(len(senderBytes))) }()
	}

	payloadPtr, err := r.writeToMemory(appModule, payload)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write payload to memory: %v", err))
	}
	if appModule.deallocate != nil && payloadPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, payloadPtr, int32(len(payload))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest alloc fault, not a host-side problem
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, fmt.Sprintf("failed to write state to memory: %v", err))
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	// Call the process_request function
	// Wasm supports only int64, so we cast appId to int64
	// requestType is passed as int32 to the WASM module
	result, err := processRequestFunc.Call(appModule.store, wasmAppId, senderPtr, int32(len(senderBytes)), int32(requestType), payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
	if err != nil {
		guestFaulted = isGuestTrap(err) // false when the guest never started, see isGuestTrap
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("failed to call process_request: %v", err))
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		guestFaulted = errors.Is(err, errGuestFault) // guest-provided pointer, not a host-side problem
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedExtractingResultBytes, fmt.Sprintf("failed to extract wasm module result bytes: %v", err))
	}

	// Deserialize the result
	var processResult appCommon.ProcessResult
	if err := json.Unmarshal(resultBytes, &processResult); err != nil {
		guestFaulted = true // guest produced a result the host cannot parse
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeJsonUnmarshalError, fmt.Sprintf("failed to unmarshal process result: %v", err))
	}

	if processResult.Error != "" {
		return nil, nil, nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("failed to process request: %s", processResult.Error))
	}

	r.log.Info("Wasmtime Runtime: Successfully processed request for application %d, generated %d events and %d withdrawals", appId, len(processResult.Events), len(processResult.Withdrawals))
	return processResult.State, processResult.Events, processResult.AppEvents, processResult.Withdrawals, processResult.Report, processResult.Fuel.ToInt(), nil
}

func (r *WasmtimeRuntime) extractResultBytes(result interface{}, appModule *ApplicationModule) (out []byte, err error) {
	// ---- Panic shield: NEVER let guest crash the host ----
	defer func() {
		if rec := recover(); rec != nil {
			err = asGuestFault(fmt.Errorf("panic while extracting wasm result: %v", rec))
		}
	}()

	ptr, err := toInt32(result)
	if err != nil {
		return nil, asGuestFault(fmt.Errorf("wasm module returned unexpected type for result pointer: %w", err))
	}
	if ptr == 0 {
		return nil, asGuestFault(fmt.Errorf("wasm module returned null pointer"))
	}

	if appModule == nil || appModule.store == nil || appModule.memory == nil {
		return nil, fmt.Errorf("invalid wasm module state")
	}

	mem := appModule.memory.UnsafeData(appModule.store)
	memLen := len(mem)

	start := int(ptr)
	if start < 0 || start+4 > memLen {
		return nil, asGuestFault(fmt.Errorf("invalid memory access for length prefix"))
	}

	// ---- Read length prefix ----
	dataLen := binary.LittleEndian.Uint32(mem[start : start+4])

	if dataLen == 0 {
		if appModule.deallocate != nil {
			_, _ = appModule.deallocate.Call(appModule.store, ptr, int32(4))
		}
		return nil, asGuestFault(fmt.Errorf("empty result from wasm module"))
	}

	// TODO consider using this limit for increasing safety (to be defined as a constant somewhere)
	/*
		// ---- Hard limits (DoS protection) ----
		const maxResultSize = 4 * 1024 * 1024 // 4 MB
		if dataLen > maxResultSize {
			return nil, fmt.Errorf("result too large: %d bytes", dataLen)
		}
	*/

	// Prevent int overflow
	if dataLen > uint32(math.MaxInt32-4) {
		return nil, asGuestFault(fmt.Errorf("result length overflows int"))
	}

	totalLen := int(dataLen) + 4
	if start+totalLen > memLen {
		return nil, asGuestFault(fmt.Errorf("invalid memory access for payload"))
	}

	// ---- Copy payload ----
	out = make([]byte, dataLen)
	copy(out, mem[start+4:start+totalLen])

	// ---- Deallocate guest memory ----
	if appModule.deallocate != nil {
		if _, derr := appModule.deallocate.Call(appModule.store, ptr, int32(totalLen)); derr != nil {
			r.log.Warn("Wasmtime Runtime: failed to deallocate wasm memory for result: %v", derr)
		}
	}

	// ---- Detect guest serialization failure sentinel ----
	if string(out) == appCommon.WasmSerializationError {
		return nil, asGuestFault(fmt.Errorf("wasm module failed to serialize response/error"))
	}

	return out, nil
}

// touchModule moves appId to front of access order (most recently used).
// Must be called with moduleLock held (write lock).
func (r *WasmtimeRuntime) touchModule(appId common.ApplicationIdType) {
	if elem, exists := r.accessElements[appId]; exists {
		r.accessOrder.MoveToFront(elem)
	} else {
		elem := r.accessOrder.PushFront(appId)
		r.accessElements[appId] = elem
	}
}

// evictIfNeeded removes least-recently-used modules until cache is within limit.
// Must be called with moduleLock held (write lock).
func (r *WasmtimeRuntime) evictIfNeeded() {
	if r.maxCachedModules <= 0 {
		return
	}
	for r.accessOrder.Len() > r.maxCachedModules {
		back := r.accessOrder.Back()
		evictId := back.Value.(common.ApplicationIdType)
		r.log.Info("Module cache full (%d/%d), evicting app %d", r.accessOrder.Len(), r.maxCachedModules, evictId)
		if module, exists := r.modules[evictId]; exists {
			r.cleanupModule(evictId, module)
			delete(r.modules, evictId)
		}
		r.accessOrder.Remove(back)
		delete(r.accessElements, evictId)
	}
}

// removeFromAccessOrder removes appId from LRU tracking.
// Must be called with moduleLock held.
func (r *WasmtimeRuntime) removeFromAccessOrder(appId common.ApplicationIdType) {
	if elem, exists := r.accessElements[appId]; exists {
		r.accessOrder.Remove(elem)
		delete(r.accessElements, appId)
	}
}

// cleanupModule releases all resources associated with a single ApplicationModule.
//
// CONCURRENCY: module.Close() calls store.Close(), which frees the native
// wasmtime state immediately (rather than at a GC finalizer pass), so it must
// never run while another goroutine is calling into that store — that would be a
// use-after-free, i.e. a native segfault inside the enclave. Callers must hold
// BOTH execLock and moduleLock: moduleLock alone does not cover the guest call,
// because Deposit and ProcessRequest execute after releasing it. See the execLock
// declaration for the full rationale and for what needs to change here to allow
// concurrent execution (refcounting live calls per module).
func (r *WasmtimeRuntime) cleanupModule(appId common.ApplicationIdType, module *ApplicationModule) {
	r.log.Info("Cleaning up module %d", appId)
	module.Close()
}

// shutdownExecLockTimeout bounds how long Close waits for an in-flight guest
// call before giving up and shutting down without releasing wasm resources.
//
// The value only trades shutdown latency against how often teardown is skipped:
// a guest call still running at the deadline is abandoned, which is harmless
// because the process is exiting and the OS reclaims the memory anyway. Note that
// this is NOT dimensioned so that a legitimate request always completes first —
// nothing bounds guest execution today (see newPinnedEngine), so no timeout could
// guarantee that, and a healthy but slow guest can trip it. 5s keeps SIGTERM
// responsive; raise it only if skipped teardown ever proves to matter.
const shutdownExecLockTimeout = 5 * time.Second

// tryAcquireExecLock acquires execLock, giving up after timeout. Reports whether
// it was acquired; the caller must Unlock only when it returns true.
//
// Polls because sync.Mutex has no timed acquire. That is acceptable here because
// this is only used on the shutdown path: the hot paths take execLock directly
// and block, which is the intended serialization.
func (r *WasmtimeRuntime) tryAcquireExecLock(timeout time.Duration) bool {
	const pollInterval = 10 * time.Millisecond

	deadline := time.Now().Add(timeout)
	for {
		if r.execLock.TryLock() {
			return true
		}
		if !time.Now().Before(deadline) {
			return false
		}
		time.Sleep(pollInterval)
	}
}

// evictFaultedModule drops a module from the cache after its guest faulted, so
// that the next request for that app gets a freshly instantiated one.
//
// A trap unwinds the guest without running any of its cleanup, so the instance is
// left with a heap in an unknown state: whatever was allocated before the trap is
// never freed and the allocator's own bookkeeping may be inconsistent. Keeping
// that instance in the LRU means the next request runs on the damaged heap, and
// the leak accumulates until the module happens to be evicted for capacity. The
// guest memory cap makes this materially more likely, since exceeding it makes
// TinyGo panic (a trap) rather than growing memory.
//
// Cheap by design: the cost of being wrong here is one recompile plus
// re-instantiation on the next request for that app.
//
// Must be called with execLock held and moduleLock NOT held — it closes the
// store, see cleanupModule.
func (r *WasmtimeRuntime) evictFaultedModule(appId common.ApplicationIdType) {
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	module, exists := r.modules[appId]
	if !exists {
		return
	}

	r.log.Warn("Runtime: evicting module %d after a guest fault, it will be reloaded on the next request", appId)
	r.cleanupModule(appId, module)
	delete(r.modules, appId)
	r.removeFromAccessOrder(appId)
}

// Close closes the wasmtime runtime and cleans up resources
func (r *WasmtimeRuntime) Close() error {
	r.log.Info("Wasmtime Runtime: Closing wasmtime runtime")

	// Closes every module and the engine, so it must not overlap an in-flight
	// guest call. Bounded rather than blocking: guest execution is currently
	// unbounded (no fuel, no epoch deadline, no stack limit — see
	// newPinnedEngine), so a runaway guest holds execLock forever and an
	// unbounded wait here would hang shutdown for good.
	//
	// Giving up is safe: skipping teardown only means the native wasmtime state
	// outlives this call, and the process is about to exit and have it reclaimed
	// by the OS. Closing the stores anyway would be the use-after-free that
	// execLock exists to prevent, so a leak-at-exit is strictly the better trade.
	if !r.tryAcquireExecLock(shutdownExecLockTimeout) {
		r.log.Error("Wasmtime Runtime: a guest call is still running after %v, "+
			"shutting down without releasing wasm resources (they are reclaimed on process exit)",
			shutdownExecLockTimeout)
		return fmt.Errorf("timed out after %v waiting for in-flight guest execution to finish", shutdownExecLockTimeout)
	}
	defer r.execLock.Unlock()

	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	for appId, module := range r.modules {
		r.cleanupModule(appId, module)
		delete(r.modules, appId)
	}

	r.accessOrder.Init()
	r.accessElements = make(map[common.ApplicationIdType]*list.Element)
	// Deterministically free the engine's native state (compiled-code caches
	// etc.). The engine is refcounted at the C level, so any store not tracked
	// in r.modules keeps its share alive until it is closed or collected.
	if r.engine != nil {
		r.engine.Close()
	}
	r.engine = nil
	r.log.Info("Wasmtime Runtime: Wasmtime runtime closed successfully")
	return nil
}

// UnloadModule unloads a WASM module and frees its resources.
// Currently unused in production — module removal is handled automatically by
// evictIfNeeded() as part of the LRU cache. This method exists as a public API
// for explicit removal (e.g. future app decommissioning or admin commands).
// Note: evictIfNeeded() intentionally duplicates the removal logic (cleanupModule
// + delete + removeFromAccessOrder) rather than calling UnloadModule, because
// evictIfNeeded runs inside an already-held write lock and UnloadModule acquires
// its own lock — Go's sync.RWMutex is not reentrant, so calling one from the
// other would deadlock.
func (r *WasmtimeRuntime) UnloadModule(appId common.ApplicationIdType) error {
	// Closes a module, so it must not overlap an in-flight guest call.
	r.execLock.Lock()
	defer r.execLock.Unlock()

	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	module, exists := r.modules[appId]
	if !exists {
		r.log.Warn("UnloadModule: module %d not found", appId)
		return nil
	}

	r.cleanupModule(appId, module)
	delete(r.modules, appId)
	r.removeFromAccessOrder(appId)

	r.log.Info("Module %d unloaded successfully", appId)
	return nil
}

// SetMaxCachedModules updates the LRU cache limit at runtime.
// A value of 0 means unlimited. Excess modules are not evicted immediately;
// eviction happens lazily on the next loadModule call, which is the only path
// that mutates the module map during normal operation.
func (r *WasmtimeRuntime) SetMaxCachedModules(max int) {
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	r.log.Info("Runtime: Setting maxCachedModules from %d to %d", r.maxCachedModules, max)
	r.maxCachedModules = max
}

// GetMaxCachedModules returns the current LRU cache limit. 0 means unlimited.
func (r *WasmtimeRuntime) GetMaxCachedModules() int {
	r.moduleLock.RLock()
	defer r.moduleLock.RUnlock()
	return r.maxCachedModules
}

// SetMaxGuestMemoryBytes updates the per-guest linear memory cap. The cap is
// applied to stores created from now on; modules already cached keep the cap
// they were created with. Values <= 0 or above the 2 GiB ABI ceiling are
// replaced with the ceiling (see maxGuestMemoryCeilingBytes) and reported.
// Clamps rather than rejecting (unlike pkg/executor.Config.Validate, which fails
// at start-up) so a late call cannot take a running enclave down.
func (r *WasmtimeRuntime) SetMaxGuestMemoryBytes(maxBytes int64) {
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	if maxBytes <= 0 || maxBytes > maxGuestMemoryCeilingBytes {
		r.log.Warn("Runtime: requested guest memory cap %d out of range, using ceiling %d", maxBytes, int64(maxGuestMemoryCeilingBytes))
		maxBytes = maxGuestMemoryCeilingBytes
	}

	r.log.Info("Runtime: Setting maxGuestMemoryBytes from %d to %d", r.maxGuestMemoryBytes, maxBytes)
	r.maxGuestMemoryBytes = maxBytes
}

// GetMaxGuestMemoryBytes returns the per-guest linear memory cap applied to
// newly created stores.
func (r *WasmtimeRuntime) GetMaxGuestMemoryBytes() int64 {
	r.moduleLock.RLock()
	defer r.moduleLock.RUnlock()
	return r.maxGuestMemoryBytes
}

// ToWasmType converts ApplicationIdType (uint64) to the WASM i64 representation
// (int64). WASM has no unsigned integer types — i64 carries 64 bits regardless of
// signedness, so the bit pattern passes through unchanged.
func ToWasmType(aid common.ApplicationIdType) (int64, error) {
	return int64(aid), nil
}

// retrieves statistics from guest memory allocation. In this version we use the ABI C specification of tinygo, using directly
// the values returned from the wasm guest, without using marshal/unmarshal into json structs.
// This implementation is here mostly as a reference on how to handle a multireturn exported func
func (r *WasmtimeRuntime) GetAllocatedMemoryStats(ctx context.Context, appId common.ApplicationIdType, wasm []byte) (int64, int64, error) {
	// Held across the guest call, see the execLock declaration.
	r.execLock.Lock()
	defer r.execLock.Unlock()

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Call the WASM function to generate the report
	statsFunc := appModule.instance.GetFunc(appModule.store, "get_allocated_memory_stats")
	if statsFunc == nil {
		return 0, 0, fmt.Errorf("function get_allocated_memory_stats not found in wasm module")
	}

	allocateFunc := appModule.instance.GetFunc(appModule.store, "allocate")
	if allocateFunc == nil {
		return 0, 0, fmt.Errorf("allocate function not found in WASM module")
	}

	// we allocate a buffer wher the results will be stored
	const resultSize = 2 * 8 // 2 int64 values, 8 bytes each
	rawResultPtr, err := allocateFunc.Call(appModule.store, int32(resultSize))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to allocate memory for results: %w", err)
	}

	resultPtr, err := toInt32(rawResultPtr)
	if err != nil {
		return 0, 0, fmt.Errorf("allocate returned unexpected value: %w", err)
	}
	if resultPtr == 0 {
		return 0, 0, fmt.Errorf("allocate returned null pointer (0)")
	}

	// ensure we free the result buffer on exit
	if appModule.deallocate != nil {
		defer func() {
			_, _ = appModule.deallocate.Call(appModule.store, resultPtr, int32(resultSize))
		}()
	}

	// according to ABI C interface, tinygo uses an inout ptr to store return values, even the signature is different
	_, err = statsFunc.Call(appModule.store, resultPtr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to call func: %w", err)
	}

	// Read the 2 int64 values sequentially from the WASM memory

	// Get the raw byte view of the WASM memory
	memoryData := appModule.memory.UnsafeData(appModule.store)

	// resultPtr came from the guest's allocate, so validate the whole read window
	// before slicing: a pointer within resultSize bytes of the end of memory would
	// otherwise panic the host. Unlike extractResultBytes this function has no
	// recover shield, so an unchecked slice here takes the process down — see the
	// "NEVER let guest crash the host" note there.
	if resultPtr < 0 || int64(resultPtr)+resultSize > int64(len(memoryData)) {
		return 0, 0, fmt.Errorf("stats result pointer %d out of bounds: %d bytes needed, memory size is %d",
			resultPtr, resultSize, len(memoryData))
	}

	// Read the int64 values using binary.LittleEndian
	mem_size := int64(binary.LittleEndian.Uint64(memoryData[resultPtr : resultPtr+8]))
	total_bytes := int64(binary.LittleEndian.Uint64(memoryData[resultPtr+8 : resultPtr+16]))
	r.log.Info("Stats: %d, %d", mem_size, total_bytes)

	// Now deduce the memory used here for reading the return values that will be released on exit by the defered dealloc call
	mem_size -= 1
	total_bytes -= resultSize

	return mem_size, total_bytes, nil
}

// retrieves statistics from guest memory allocation (second version)
func (r *WasmtimeRuntime) GetAllocatedMemoryStats2(ctx context.Context, appId common.ApplicationIdType, wasm []byte) (int64, int64, error) {
	// Held across the guest call, see the execLock declaration.
	r.execLock.Lock()
	defer r.execLock.Unlock()

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Call the WASM function to generate the report
	statsFunc := appModule.instance.GetFunc(appModule.store, "get_memory_stats")
	if statsFunc == nil {
		return 0, 0, fmt.Errorf("function get_memory_stats not found in wasm module")
	}

	result, err := statsFunc.Call(appModule.store)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to call func: %w", err)
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	// Deserialize the result
	var stats appCommon.MemoryStats
	if err := json.Unmarshal(resultBytes, &stats); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal stats result: %w", err)
	}

	mem_size := stats.MapSize
	total_bytes := stats.CumulativeMemorySize

	r.log.Info("Stats: %d, %d", mem_size, total_bytes)

	// the counters we got from WASM module do not take into account the memory used for passing them across guest/host ABI interface
	// that is allocated in marshal/unmarshall process, but we can leave with that

	return mem_size, total_bytes, nil
}

func (r *WasmtimeRuntime) configureWasiLogPipes(_ context.Context, appId common.ApplicationIdType, store *wasmtime.Store) (func(), error) {
	// create a new WASI config
	wasiConfig := wasmtime.NewWasiConfig()

	// Create a named pipe to capture WASI output
	r.log.Info("Creating named log pipe")
	// add a nanosecond timestamp to the path to prevent collisions between concurrent module loads
	fifoPath := fmt.Sprintf("/tmp/wasm_log_%d_%d", appId, time.Now().UnixNano())

	// should never happen, but ensure no stale file exists otherwise we will hit an error
	_ = os.Remove(fifoPath)

	// Create the FIFO (Named Pipe)
	if err := syscall.Mkfifo(fifoPath, 0600); err != nil {
		r.log.Error("Error creating pipe")
		return nil, fmt.Errorf("failed to create log pipe: %w", err)
	}

	// immediate Unlink after creation (The file disappears from /tmp but data flows)
	// this is especially useful in the enclave environment, where we have a RAMFS file system
	defer os.Remove(fifoPath)

	// we use this for avoiding that the Reader OpenFile(O_RDONLY) blocks waiting for the WASI writer to connect
	// this Dummy writer must stay open until cleanup, otherwise will send EOF to the reader below
	// (reader blocking is potentially dangerous in edge cases on some error path for races, resource dangling and panics)
	dummyWriter, err := os.OpenFile(fifoPath, os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open dummy fifo writer: %w", err)
	}

	// Open the reader on main thread - won't block because dummyWriter is already connected
	readerFile, err := os.OpenFile(fifoPath, os.O_RDONLY, 0600)
	if err != nil {
		dummyWriter.Close()
		return nil, fmt.Errorf("failed to open fifo reader: %w", err)
	}

	// Channel to signal goroutine termination
	doneCh := make(chan struct{})

	// Start the log reader goroutine
	go func() {
		defer close(doneCh)

		// Use bufio.Reader.ReadLine() instead of Scanner to handle oversized lines gracefully.
		// When a line exceeds the buffer, ReadLine() returns isPrefix=true, allowing us to
		// discard the overflow and continue reading subsequent lines (Scanner would stop entirely).
		// Buffer size of 64KB should handle most legitimate log lines.
		const logBufferSize = 64 * 1024 // 64KB
		reader := bufio.NewReaderSize(readerFile, logBufferSize)

		// Rate limiting to prevent log flooding from malicious/buggy WASM modules.
		// Uses a simple fixed window: allow up to maxLogsPerWindow logs per time window.
		// Excess logs are dropped and periodically reported.
		const (
			maxLogsPerWindow = 1000             // Max logs allowed per window
			windowDuration   = 1 * time.Second  // Time window for rate limiting
			reportInterval   = 10 * time.Second // How often to report dropped logs
		)
		var (
			windowStart    = time.Now()
			logsInWindow   = 0
			droppedLogs    = 0
			lastDropReport = time.Now()
			totalDropped   = 0
		)

		r.log.Debug("Starting WASM log pipe reader for app %d", appId)

		for {
			lineBytes, isPrefix, err := reader.ReadLine()
			if err != nil {
				// Don't log if it's due to EOF or file being closed during cleanup
				if !errors.Is(err, io.EOF) && !errors.Is(err, os.ErrClosed) {
					r.log.Warn("WASM log pipe read error for app %d: %v", appId, err)
				}
				break
			}

			line := string(lineBytes)
			truncated := isPrefix

			// If line exceeded buffer, discard remaining bytes until end of line
			for isPrefix {
				_, isPrefix, err = reader.ReadLine()
				if err != nil {
					break
				}
			}

			if truncated {
				line += "... [truncated]"
			}

			// Rate limiting check
			now := time.Now()
			if now.Sub(windowStart) >= windowDuration {
				// Reset window
				windowStart = now
				logsInWindow = 0
			}

			if logsInWindow >= maxLogsPerWindow {
				// Drop this log
				droppedLogs++
				totalDropped++

				// Periodically report dropped logs
				if now.Sub(lastDropReport) >= reportInterval {
					r.log.Warn("WASM log rate limit: dropped %d logs from app %d in last %v (total dropped: %d)",
						droppedLogs, appId, reportInterval, totalDropped)
					droppedLogs = 0
					lastDropReport = now
				}
				continue
			}
			logsInWindow++

			// Parse guest log level prefix and route to appropriate host log level
			// Expected format: LVL message (e.g., "INF Processing request")
			switch {
			case strings.HasPrefix(line, "ERR "):
				r.log.Error("Wasm Guest [%d]: %s", appId, line[4:])
			case strings.HasPrefix(line, "WRN "):
				r.log.Warn("Wasm Guest [%d]: %s", appId, line[4:])
			case strings.HasPrefix(line, "INF "):
				r.log.Info("Wasm Guest [%d]: %s", appId, line[4:])
			case strings.HasPrefix(line, "DBG "):
				r.log.Debug("Wasm Guest [%d]: %s", appId, line[4:])
			case strings.HasPrefix(line, "TRC "):
				r.log.Trace("Wasm Guest [%d]: %s", appId, line[4:])
			default:
				// No recognized prefix - log as info (backwards compatible)
				r.log.Info("Wasm Guest [%d]: %s", appId, line)
			}
		}

		// Report any remaining dropped logs on exit
		if droppedLogs > 0 {
			r.log.Warn("WASM log rate limit: dropped %d logs from app %d before exit (total dropped: %d)",
				droppedLogs, appId, totalDropped)
		}

		r.log.Debug("WASM log pipe reader exited for app %d", appId)
	}()

	// This helper cleanup func will be also used when the module is disposed of.
	// We should call it when we unload a module otherwise we might leak file descriptors.
	// Note: the closure maintains all necessary state to coordinate cleanup with the background goroutine, regardless of when or where it's invoked
	// ---
	// Because the wasmtime map holds the reference to a module, the GC cannot touch the module and all associated resources.
	//     Map->ApplicationModule->store->writer fd open->go routine blocked->reader fd open
	// therefore for every WASM module 2 fds are leaked and a go routine is alive in bg (beside the allocated memory)
	cleanupFileDescriptors := func() {
		r.log.Info("Closing log pipe FDs for app %d", appId)

		// We force termination by closing both ends of the pipe immediately.
		// Closing readerFile is critical: it forces the background goroutine to exit
		// even if the WASM module (the writer) is still holding its FD or is stuck.
		// This prevents goroutine leaks and avoids long timeouts during cleanup.
		if err := dummyWriter.Close(); err != nil {
			r.log.Warn("Failed to close dummy writer FD: %v", err)
		}
		if err := readerFile.Close(); err != nil {
			r.log.Warn("Failed to close reader FD: %v", err)
		}

		// Brief wait for goroutine to complete. The loop in the goroutine
		// handles os.ErrClosed, so this is safe and deterministic.
		select {
		case <-doneCh:
			// Goroutine exited cleanly
		case <-time.After(200 * time.Millisecond):
			r.log.Warn("Log pipe goroutine slow to exit for app %d", appId)
		}
	}

	r.log.Info("Installing log pipes")
	// Configure WASI to use the 'write' end of our pipe
	err = wasiConfig.SetStdoutFile(fifoPath)
	if err != nil {
		cleanupFileDescriptors()
		return nil, fmt.Errorf("failed to install stdout to pipe: %w", err)
	}
	r.log.Info("Installed stdout log pipe")

	err = wasiConfig.SetStderrFile(fifoPath)
	if err != nil {
		cleanupFileDescriptors()
		return nil, fmt.Errorf("failed to install stderr to pipe: %w", err)
	}
	r.log.Info("Installed stderr log pipe")

	// attach WASI config to the store
	store.SetWasi(wasiConfig)

	return cleanupFileDescriptors, nil
}
