package wasm

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/horizen-pes/pkg/common"
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
	modules    map[string]*ApplicationModule // Map of application ID to module
	moduleLock sync.RWMutex                  // Lock for module access
}

// NewWasmtimeRuntime creates a new wasmtime runtime instance
func NewWasmtimeRuntime() *WasmtimeRuntime {
	log.Println("Runtime: Initializing wasmtime runtime")

	// Create a new engine with default configuration
	engine := wasmtime.NewEngine()

	return &WasmtimeRuntime{
		engine:  engine,
		modules: make(map[string]*ApplicationModule),
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
	start := int(ptr)
	end := start + len(data)
	if start < 0 || end < start || end > len(memData) {
		return 0, fmt.Errorf("invalid memory allocation: ptr=%d len=%d memsize=%d", ptr, len(data), len(memData))
	}

	copy(memData[start:end], data)
	return ptr, nil
}

// readFromMemory copies `length` bytes from module memory at offset `ptr`.
func (r *WasmtimeRuntime) readFromMemory(module *ApplicationModule, ptr int32, length int32) ([]byte, error) {
	if module == nil {
		return nil, fmt.Errorf("module is nil")
	}
	if module.memory == nil {
		return nil, fmt.Errorf("memory not initialized")
	}
	if length == 0 {
		return []byte{}, nil
	}

	if module.store == nil {
		return nil, fmt.Errorf("wasm module has a nil store")
	}

	memData := module.memory.UnsafeData(module.store)
	start := int(ptr)
	if start < 0 {
		return nil, fmt.Errorf("negative pointer")
	}
	end := start + int(length)
	if end > len(memData) {
		return nil, fmt.Errorf("invalid memory access: ptr=%d length=%d memsize=%d", ptr, length, len(memData))
	}

	out := make([]byte, length)
	copy(out, memData[start:end])
	return out, nil
}

// -----------------------------
// Module loading / lifecycle
// -----------------------------

// TODO - this method most of times does not use the wasm bytes, because it will find the module already loaded, unless the executor has
// undergo a restart, and the cache has to be populated again. This caching is also sub-optimal because every time the wasm bytes are passed along
// without being used.
// One approach could be that we use just the appId, and the executor asks to the manager for wasm bytes (that are stored in the dblayer) if it
// does not find them in the cache.
func (r *WasmtimeRuntime) getOrLoadModule(ctx context.Context, appId string, wasm []byte) (*ApplicationModule, error) {
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
	_, err := r.loadModuleUnlocked(ctx, appId, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to load module: %w", err)
	}

	// Return loaded module
	if module, exists := r.modules[appId]; exists {
		return module, nil
	}
	return nil, fmt.Errorf("module not found after loading: %s", appId)
}

// LoadModule loads a WASM module and returns initial state and state root
func (r *WasmtimeRuntime) LoadModule(ctx context.Context, appId string, wasm []byte) ([]byte, error) {
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()
	return r.loadModuleUnlocked(ctx, appId, wasm)
}

// loadModuleUnlocked contains the core logic for loading a module, but without locking.
// This method should only be called when a lock is already held.
func (r *WasmtimeRuntime) loadModuleUnlocked(_ context.Context, appId string, wasm []byte) ([]byte, error) {
	log.Printf("Wasmtime Runtime: Loading WASM module for application %s (wasm size: %d bytes)", appId, len(wasm))

	// Compile the WASM module
	module, err := wasmtime.NewModule(r.engine, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to compile WASM module: %w", err)
	}

	// Create a per-module store
	store := wasmtime.NewStore(r.engine)

	// Create WASI configuration and linker for TinyGo WASI imports
	linker := wasmtime.NewLinker(r.engine)
	if err := linker.DefineWasi(); err != nil {
		return nil, fmt.Errorf("failed to define WASI: %w", err)
	}

	// Attach WASI config to the store
	wasiConfig := wasmtime.NewWasiConfig()
	store.SetWasi(wasiConfig)

	// Instantiate the module using the module-specific store
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	// Get the memory export.
	// Note: the WASM module is responsible for defining and exporting its own memory, this is done by tinygo,
	// check for instance WAT (WebAssembly text format) generated via wasm2wat tool
	memoryExport := instance.GetExport(store, "memory")
	if memoryExport == nil {
		return nil, fmt.Errorf("memory export not found in WASM module")
	}

	memory := memoryExport.Memory()
	if memory == nil {
		return nil, fmt.Errorf("memory export is not a memory")
	}

	// Get the load_module function
	loadModuleFunc := instance.GetFunc(store, "load_module")
	if loadModuleFunc == nil {
		return nil, fmt.Errorf("load_module function not found in WASM module")
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

	// Write appId to memory
	appIdBytes := []byte(appId)
	appIdPtr, err := r.writeToMemory(appModule, appIdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write appId to memory: %w", err)
	}
	// ensure we free the input buffer after load_module returns
	if appModule.deallocate != nil && appIdPtr != 0 {
		defer func() {
			_, _ = appModule.deallocate.Call(appModule.store, appIdPtr, int32(len(appIdBytes)))
		}()
	}

	// Call the load_module function (use store of this module)
	result, err := loadModuleFunc.Call(store, appIdPtr, int32(len(appIdBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to call load_module: %w", err)
	}

	// Extract the result bytes
	stateBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	log.Printf("Wasmtime Runtime: Raw result from WASM: %s", string(stateBytes))

	// Store the module in the runtime registry
	r.modules[appId] = appModule
	log.Printf("Wasmtime Runtime: Successfully loaded WASM module for application %s", appId)

	return stateBytes, nil
}

// -----------------------------
// Public operations (Deposit, ProcessRequest, ...)
//
// Each of these ensures that any temporary allocations created by writeToMemory
// are deallocated after the guest call by using defer cleanup.
// -----------------------------

// Deposit processes a deposit
func (r *WasmtimeRuntime) Deposit(ctx context.Context, appId string, sender string, value uint64, state []byte, wasm []byte) ([]byte, []common.PlainEvent, error) {
	log.Printf("Wasmtime Runtime: Processing deposit for application %s (value: %d wei for sender: %s)", appId, value, sender)

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Get the deposit function
	depositFunc := appModule.instance.GetFunc(appModule.store, "deposit")
	if depositFunc == nil {
		return nil, nil, fmt.Errorf("deposit function not found in WASM module")
	}

	// Write parameters to memory (and schedule deallocation)
	appIdBytes := []byte(appId)
	appIdPtr, err := r.writeToMemory(appModule, appIdBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write appId to memory: %w", err)
	}
	if appModule.deallocate != nil && appIdPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, appIdPtr, int32(len(appIdBytes))) }()
	}

	senderBytes := []byte(sender)
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write sender to memory: %w", err)
	}
	if appModule.deallocate != nil && senderPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, senderPtr, int32(len(senderBytes))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write state to memory: %w", err)
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	// Call the deposit function
	//
	// Note: wasmtime-go’s Call() API only accepts int64 for i64 parameters.
	// Therefore, we must cast 'value' (uint64) to int64 to satisfy the API.
	//
	// The 64-bit pattern is passed unchanged into the WASM call. On the guest side,
	// since the function signature expects a uint64, Go will interpret the bits
	// correctly as an unsigned value.
	//
	// TODO: in the future we should replace uint64 with *big.Int for amounts,
	// since values may exceed 64 bits. This will require redesigning
	// the guest/host ABI to pass large integers (e.g., via memory + length).
	result, err := depositFunc.Call(appModule.store, appIdPtr, int32(len(appIdBytes)), senderPtr, int32(len(senderBytes)), int64(value), statePtr, int32(len(state)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call deposit: %w", err)
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	log.Printf("Wasmtime Runtime: Raw deposit result from WASM: %s", string(resultBytes))

	// Deserialize the result
	var depositResult appCommon.DepositResult
	if err := json.Unmarshal(resultBytes, &depositResult); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal deposit result: %w", err)
	}

	if depositResult.Error != "" {
		return nil, nil, fmt.Errorf("deposit failed, wasm module error: %v", depositResult.Error)
	}

	log.Printf("Wasmtime Runtime: Successfully processed deposit for sender %s, generated %d events", sender, len(depositResult.Events))
	return depositResult.State, depositResult.Events, nil
}

// ProcessRequest processes a request and returns the new state, events, and withdrawals
func (r *WasmtimeRuntime) ProcessRequest(ctx context.Context, appId string, sender string, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, error) {
	log.Printf("Wasmtime Runtime: Processing request for application %s (payload size: %d, state size: %d)", appId, len(payload), len(state))
	if len(payload) == 0 {
		log.Printf("Wasmtime Runtime: Empty payload for application %s, returning current state", appId)
		return state, nil, nil, nil
	}

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Get the process_request function
	processRequestFunc := appModule.instance.GetFunc(appModule.store, "process_request")
	if processRequestFunc == nil {
		return nil, nil, nil, fmt.Errorf("process_request function not found in WASM module")
	}

	// Write parameters to memory (and schedule deallocation)
	appIdBytes := []byte(appId)
	appIdPtr, err := r.writeToMemory(appModule, appIdBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write appId to memory: %w", err)
	}
	if appModule.deallocate != nil && appIdPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, appIdPtr, int32(len(appIdBytes))) }()
	}

	senderBytes := []byte(sender)
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write sender to memory: %w", err)
	}
	if appModule.deallocate != nil && senderPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, senderPtr, int32(len(senderBytes))) }()
	}

	payloadPtr, err := r.writeToMemory(appModule, payload)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write payload to memory: %w", err)
	}
	if appModule.deallocate != nil && payloadPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, payloadPtr, int32(len(payload))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write state to memory: %w", err)
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	// Call the process_request function
	result, err := processRequestFunc.Call(appModule.store, appIdPtr, int32(len(appIdBytes)), senderPtr, int32(len(senderBytes)), payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to call process_request: %w", err)
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	// Deserialize the result
	var processResult appCommon.ProcessResult
	if err := json.Unmarshal(resultBytes, &processResult); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal process result: %w", err)
	}

	if processResult.Error != "" {
		return nil, nil, nil, fmt.Errorf("process request failed, wasm module error: %v", processResult.Error)
	}

	log.Printf("Wasmtime Runtime: Successfully processed request for application %s, generated %d events and %d withdrawals", appId, len(processResult.Events), len(processResult.Withdrawals))
	return processResult.State, processResult.Events, processResult.Withdrawals, nil
}

// GenerateDeanonymizationReport generates a deanonymization report
func (r *WasmtimeRuntime) GenerateDeanonymizationReport(ctx context.Context, appId string, payload []byte, state []byte, wasm []byte) ([]byte, error) {
	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Call the WASM function to generate the report
	generateReportFunc := appModule.instance.GetFunc(appModule.store, "generate_deanonymization_report")
	if generateReportFunc == nil {
		return nil, fmt.Errorf("function generate_deanonymization_report not found in wasm module")
	}

	payloadPtr, err := r.writeToMemory(appModule, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to write payload to memory: %w", err)
	}
	if appModule.deallocate != nil && payloadPtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, payloadPtr, int32(len(payload))) }()
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, fmt.Errorf("failed to write state to memory: %w", err)
	}
	if appModule.deallocate != nil && statePtr != 0 {
		defer func() { _, _ = appModule.deallocate.Call(appModule.store, statePtr, int32(len(state))) }()
	}

	// Call the generate_deanonymization_report function
	result, err := generateReportFunc.Call(appModule.store, payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
	if err != nil {
		return nil, fmt.Errorf("failed to call generate_deanonymization_report: %w", err)
	}

	// Extract the result bytes
	reportBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	// Deserialize the result
	var deanonymizationResult appCommon.DeanonymizationResult
	if err := json.Unmarshal(reportBytes, &deanonymizationResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deanonymization result: %w", err)
	}

	if deanonymizationResult.Error != "" {
		return nil, fmt.Errorf("deanonymization report failed, wasm module error: %v", deanonymizationResult.Error)
	}

	log.Printf("Wasmtime Runtime: Successfully generated deanonymization report for application %s", appId)
	return deanonymizationResult.Report, nil
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
	log.Printf("Wasmtime Runtime: Closing wasmtime runtime")

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
	log.Printf("Wasmtime Runtime: Wasmtime runtime closed successfully")
	return nil
}

// retrieves statistics from guest memory allocation. In this version we use the ABI C specification of tinygo, using directly
// the values returned from the wasm guest, without using marshal/unmarshal into json structs.
// This implementation is here mostly as a reference on how to handle a multireturn exported func
func (r *WasmtimeRuntime) GetAllocatedMemoryStats(ctx context.Context, appId string, wasm []byte) (int32, int32, error) {
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
	const resultSize = 2 * 4 // 2 int32 values, 4 bytes each
	resultPtrValue, err := allocateFunc.Call(appModule.store, int32(resultSize))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to allocate memory for results: %w", err)
	}

	// ensure we free the input buffer after lallocate returns
	if appModule.deallocate != nil && resultPtrValue != 0 {
		defer func() {
			_, _ = appModule.deallocate.Call(appModule.store, resultPtrValue, resultSize)
		}()
	}

	// according to ABI C interface, tinygo uses an inout ptr to store return values, even the signature is different
	resultPtr := resultPtrValue.(int32) // Get the pointer/address
	_, err = statsFunc.Call(appModule.store, resultPtr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to call func: %w", err)
	}

	// Read the 2 int32 values sequentially from the WASM memory

	// Get the raw byte view of the WASM memory
	memoryData := appModule.memory.UnsafeData(appModule.store)

	// Read the five int32 values using binary.LittleEndian
	mem_size := int32(binary.LittleEndian.Uint32(memoryData[resultPtr : resultPtr+4]))
	total_bytes := int32(binary.LittleEndian.Uint32(memoryData[resultPtr+4 : resultPtr+8]))
	//log.Printf("Stats: %d, %d", mem_size, total_bytes)

	// Now deduce the memory used here for reading the return values that will be released on exit by the defered dealloc call
	mem_size -= 1
	total_bytes -= resultSize

	return mem_size, total_bytes, nil
}

// retrieves statistics from guest memory allocation (second version)
func (r *WasmtimeRuntime) GetAllocatedMemoryStats2(ctx context.Context, appId string, wasm []byte) (int32, int32, error) {
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
	total_bytes := stats.CumulativeMemorySise

	log.Printf("Stats: %d, %d", mem_size, total_bytes)

	// the counters we got from WASM module do not take into account the memory used for passing them across guest/host ABI interface
	// that is allocated in marshal/unmarshall process, but we can leave with that

	return mem_size, total_bytes, nil
}
