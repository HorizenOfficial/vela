package wasm

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"

	"github.com/bytecodealliance/wasmtime-go"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	"github.com/horizen-pes/pkg/logger"
	appCommon "github.com/horizen-pes/pkg/wasm/common"
)

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
}

// WasmtimeRuntime implements the Runtime interface using wasmtime-go
type WasmtimeRuntime struct {
	engine     *wasmtime.Engine
	modules    map[common.ApplicationIdType]*ApplicationModule // Map of application ID to module
	moduleLock sync.RWMutex                                    // Lock for module access
	log        logger.Logger
}

// NewWasmtimeRuntime creates a new wasmtime runtime instance
func NewWasmtimeRuntime(log logger.Logger) *WasmtimeRuntime {
	log.Info("Runtime: Initializing wasmtime runtime")

	// Create a new engine with default configuration
	engine := wasmtime.NewEngine()

	return &WasmtimeRuntime{
		engine:  engine,
		modules: make(map[common.ApplicationIdType]*ApplicationModule),
		log:     log,
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
		return 0, fmt.Errorf("failed to call allocate: %w", err)
	}

	ptr, err := toInt32(rawRes)
	if err != nil {
		return 0, fmt.Errorf("allocate returned unexpected value: %w", err)
	}
	if ptr == 0 {
		return 0, fmt.Errorf("allocate returned null pointer (0)")
	}

	// Safely get a snapshot of the memory data and perform bounds checks
	memData := module.memory.UnsafeData(module.store)
	// note: the wasmtime uses a doubling of the store size when more memory is required and the UnsafeData() panics if we use a 4GB size.
	// As a consequence we can not have more than 2GB of memory allocated.
	// In other words: When the internal store grows from 2 GB -> 4 GB, UnsafeData call suddenly fails, but we are not actually running
	// out of Wasm memory, we are hitting a Go-side API limitation caused by Wasmtime’s internal buffer resize.
	// As a side effect we also could safely use an int32 ptr as the offset (avoid using uint32 cast)

	// Convert the signed int32 pointer to its true unsigned 32-bit value (uint32)
	// This correctly interprets the bits of an address as a large positive number.
	// This cast would not actually be necessary, see above.
	uPtr := uint32(ptr)

	// We cast to uint64 for the calculation to prevent overflow,
	// as len(memData) can be up to 4GB (or larger, depending on the host's slice).
	memLen := uint64(len(memData))
	allocEnd := uint64(uPtr) + uint64(len(data)) // The end address of the requested block

	if uint64(uPtr) >= memLen || allocEnd > memLen {
		return 0, fmt.Errorf("invalid memory allocation: requested block [%d - %d) exceeds memory size %d", uPtr, allocEnd, memLen)
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
	r.moduleLock.RLock()
	module, exists := r.modules[appId]
	r.moduleLock.RUnlock()
	if exists {
		return module, nil
	}

	// Acquire write lock and load
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	// Re-check if module was loaded by another goroutine while we were waiting for the lock
	if module, exists := r.modules[appId]; exists {
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
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()
	return r.loadModuleUnlocked(ctx, appId, wasm)
}

// loadModuleUnlocked contains the core logic for loading a module, but without locking.
// This method should only be called when a lock is already held.
func (r *WasmtimeRuntime) loadModuleUnlocked(_ context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, *big.Int, error) {
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
	store := wasmtime.NewStore(r.engine)

	// Create WASI configuration and linker for TinyGo WASI imports
	linker := wasmtime.NewLinker(r.engine)
	err = linker.DefineWasi()
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to define WASI: %w", err)
	}

	// Attach WASI config to the store
	wasiConfig := wasmtime.NewWasiConfig()
	store.SetWasi(wasiConfig)

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

	appModule := &ApplicationModule{
		store:      store,
		module:     module,
		instance:   instance,
		memory:     memory,
		deallocate: deallocateFunc,
	}

	// Call the load_module function
	// Wasm supports only int64, so we cast appId to int64
	result, err := loadModuleFunc.Call(store, wasmAppId)
	if err != nil {
		return nil, big.NewInt(0), fmt.Errorf("failed to call load_module: %w", err)
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

	// Store the module in the runtime registry
	r.modules[appId] = appModule
	r.log.Info("Wasmtime Runtime: Successfully loaded WASM module for application %d", appId)
	return loadResult.State, loadResult.Fuel, nil
}

// -----------------------------
// Public operations (Deposit, ProcessRequest, ...)
//
// Each of these ensures that any temporary allocations created by writeToMemory
// are deallocated after the guest call by using defer cleanup.
// -----------------------------

// Deposit processes a deposit
func (r *WasmtimeRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, value *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, *big.Int, *apperrors.RequestFailure) {
	r.log.Info("Wasmtime Runtime: Processing deposit for application %d (value: %v wei for sender: %v)", appId, value, sender)

	if value == nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, "value cannot be nil", nil)
	}

	if value.Sign() < 0 {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, "value cannot be negative", nil)
	}

	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, "invalid application id", err)
	}

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, "failed to get or load module", err)
	}

	// Get the deposit function
	depositFunc := appModule.instance.GetFunc(appModule.store, "deposit")
	if depositFunc == nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFunctionNotFound, "deposit function not found in WASM module", nil)
	}

	senderBytes := sender.Bytes()
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write sender to memory", err)
	}
	if appModule.deallocate != nil && senderPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, senderPtr, int32(len(senderBytes))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write state to memory", err)
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	valueBytes := value.Bytes()
	valuePtr, err := r.writeToMemory(appModule, valueBytes)
	if err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write value to memory", err)
	}
	if appModule.deallocate != nil && valuePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, valuePtr, int32(len(valueBytes))) }()
	}

	// Wasm supports only int64, so we cast appId to int64
	result, err := depositFunc.Call(appModule.store, wasmAppId, senderPtr, int32(len(senderBytes)), valuePtr, int32(len(valueBytes)), statePtr, int32(len(state)))
	if err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeDepositFailed, "failed to call deposit", err) // TODO some standard way of getting errors here?
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedExtractingResultBytes, "failed to extract wasm module result bytes", err)
	}

	r.log.Info("Wasmtime Runtime: Raw deposit result from WASM: %s", string(resultBytes))

	// Deserialize the result
	var depositResult appCommon.DepositResult
	if err := json.Unmarshal(resultBytes, &depositResult); err != nil {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeJsonUnmarshalError, "failed to unmarshal deposit result", err)
	}

	if depositResult.Error != "" {
		return nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeDepositFailed, "failed to process deposit", errors.New(depositResult.Error))
	}

	r.log.Info("Wasmtime Runtime: Successfully processed deposit for sender %s, generated %d events", sender, len(depositResult.Events))
	return depositResult.State, depositResult.Events, depositResult.Fuel, nil
}

// ProcessRequest processes a request and returns the new state, events, and withdrawals
func (r *WasmtimeRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, *big.Int, *apperrors.RequestFailure) {
	r.log.Info("Wasmtime Runtime: Processing request for application %d (payload size: %d, state size: %d)", appId, len(payload), len(state))

	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeInternalFallback, "invalid application id", err)
	}

	if len(payload) == 0 {
		r.log.Warn("Wasmtime Runtime: Empty payload for application %d, returning current state", appId)
		return state, nil, nil, big.NewInt(0), nil
	}

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedLoadingOrGettingModule, "failed to get or load module", err)
	}

	// Get the process_request function
	processRequestFunc := appModule.instance.GetFunc(appModule.store, "process_request")
	if processRequestFunc == nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFunctionNotFound, "process_request function not found in WASM module", nil)
	}

	senderBytes := sender.Bytes()
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write sender to memory", err)
	}
	if appModule.deallocate != nil && senderPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, senderPtr, int32(len(senderBytes))) }()
	}

	payloadPtr, err := r.writeToMemory(appModule, payload)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write payload to memory", err)
	}
	if appModule.deallocate != nil && payloadPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, payloadPtr, int32(len(payload))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write state to memory", err)
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	// Call the process_request function
	// Wasm supports only int64, so we cast appId to int64
	result, err := processRequestFunc.Call(appModule.store, wasmAppId, senderPtr, int32(len(senderBytes)), payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeRequestFuncFailed, "failed to call process_request", err)
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedExtractingResultBytes, "failed to extract wasm module result bytes", err)
	}

	// Deserialize the result
	var processResult appCommon.ProcessResult
	if err := json.Unmarshal(resultBytes, &processResult); err != nil {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeJsonUnmarshalError, "failed to unmarshal process result", err)
	}

	if processResult.Error != "" {
		return nil, nil, nil, big.NewInt(0), apperrors.New(apperrors.CodeRequestFuncFailed, "failed to process request", errors.New(processResult.Error))
	}

	r.log.Info("Wasmtime Runtime: Successfully processed request for application %d, generated %d events and %d withdrawals", appId, len(processResult.Events), len(processResult.Withdrawals))
	return processResult.State, processResult.Events, processResult.Withdrawals, processResult.Fuel, nil
}

// GenerateDeanonymizationReport generates a deanonymization report
func (r *WasmtimeRuntime) GenerateDeanonymizationReport(ctx context.Context, appId common.ApplicationIdType, payload []byte, state []byte, wasm []byte) ([]byte, *big.Int, *apperrors.RequestFailure) {

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedLoadingOrGettingModule, "failed to get or load module", err)
	}

	// Call the WASM function to generate the report
	generateReportFunc := appModule.instance.GetFunc(appModule.store, "generate_deanonymization_report")
	if generateReportFunc == nil {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeFunctionNotFound, "function generate_deanonymization_report not found in wasm module", nil)
	}

	payloadPtr, err := r.writeToMemory(appModule, payload)
	if err != nil {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write payload to memory", err)
	}
	if appModule.deallocate != nil && payloadPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, payloadPtr, int32(len(payload))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeMemoryWriteError, "failed to write state to memory", err)
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	// Call the generate_deanonymization_report function
	result, err := generateReportFunc.Call(appModule.store, payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
	if err != nil {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedToGenerateReport, "failed to call generate_deanonymization_report", err)
	}

	// Extract the result bytes
	reportBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedExtractingResultBytes, "failed to extract wasm module result bytes", err)
	}

	// Deserialize the result
	var deanonymizationResult appCommon.DeanonymizationResult
	if err := json.Unmarshal(reportBytes, &deanonymizationResult); err != nil {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeJsonUnmarshalError, "failed to unmarshal deanonymization result", err)
	}

	if deanonymizationResult.Error != "" {
		return nil, big.NewInt(0), apperrors.New(apperrors.CodeFailedToGenerateReport, "deanonymization report failed, wasm module error", errors.New(deanonymizationResult.Error))
	}

	r.log.Info("Wasmtime Runtime: Successfully generated deanonymization report for application %d", appId)
	return deanonymizationResult.Report, deanonymizationResult.Fuel, nil
}

// extractResultBytes extracts a pointer returned by the wasm module (single int32 offset)
// that points to a length-prefixed byte slice: | u32(len) | data... |.
// It copies the payload into a new Go slice, then frees the guest memory via deallocate.
func (r *WasmtimeRuntime) extractResultBytes(result interface{}, appModule *ApplicationModule) ([]byte, error) {
	ptr, err := toInt32(result)
	if err != nil {
		return nil, fmt.Errorf("wasm module returned unexpected type for result pointer: %w", err)
	}
	if ptr == 0 {
		return nil, fmt.Errorf("wasm module returned null pointer, possibly an allocation failure")
	}

	if appModule.store == nil {
		return nil, fmt.Errorf("wasm module has a nil store")
	}

	// Read the 4-byte little-endian length prefix
	memData := appModule.memory.UnsafeData(appModule.store)
	start := int(ptr)
	if start < 0 || start+4 > len(memData) {
		return nil, fmt.Errorf("invalid memory access for length prefix")
	}
	lengthBytes := memData[start : start+4]
	dataLen := binary.LittleEndian.Uint32(lengthBytes)
	if dataLen == 0 {
		// Still deallocate the 4 bytes used for the prefix
		if appModule.deallocate != nil {
			if _, err := appModule.deallocate.Call(appModule.store, ptr, int32(4)); err != nil {
				// Log the error but don't fail the operation since we have the data
				log.Printf("Wasmtime Runtime: failed to deallocate wasm memory for empty result: %v", err)
			}
		}
		return nil, fmt.Errorf("empty result from wasm module")
	}

	// to be on the safe side, check against overflow. Actually it is necessary for 32-bit platforms only
	// because on 64-bit platforms 'int' is large enough.
	if dataLen > uint32(math.MaxInt32-4) {
		return nil, fmt.Errorf("result too large: length would overflow int")
	}

	// Validate full payload is readable
	totalLen := int(dataLen) + 4
	if start+totalLen > len(memData) {
		return nil, fmt.Errorf("invalid memory access for data payload")
	}

	// Copy the payload into a Go buffer (skip the 4-byte prefix)
	resultBytes := make([]byte, dataLen)
	copy(resultBytes, memData[start+4:start+4+int(dataLen)])

	// Deallocate memory in WASM now that we have copied the data
	if appModule.deallocate != nil {
		if _, err := appModule.deallocate.Call(appModule.store, ptr, int32(totalLen)); err != nil {
			// Log the error but don't fail the operation since we have the data
			log.Printf("Wasmtime Runtime: failed to deallocate wasm memory for result: %v", err)
		}
	}

	// Check for wasm serialization failure
	if string(resultBytes) == appCommon.WasmSerializationError {
		return nil, fmt.Errorf("wasm module failed to serialize response/error")
	}

	return resultBytes, nil
}

// Close closes the wasmtime runtime and cleans up resources
func (r *WasmtimeRuntime) Close() error {
	r.log.Info("Wasmtime Runtime: Closing wasmtime runtime")

	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()

	for appId, module := range r.modules {
		// clear references so GC can reclaim them
		if module.instance != nil {
			module.instance = nil
		}
		if module.module != nil {
			module.module = nil
		}
		if module.memory != nil {
			module.memory = nil
		}
		if module.store != nil {
			module.store = nil
		}
		delete(r.modules, appId)
	}

	r.engine = nil
	r.log.Info("Wasmtime Runtime: Wasmtime runtime closed successfully")
	return nil
}

// This method converts the ApplicationIdType to the Wasm type (int64)
func ToWasmType(aid common.ApplicationIdType) (int64, error) {
	//This should never happens, but just in case
	if aid > math.MaxInt64 {
		return -1, fmt.Errorf("application ID too large: %d", aid)
	}
	return int64(aid), nil
}
