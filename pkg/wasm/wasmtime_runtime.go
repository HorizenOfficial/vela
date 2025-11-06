package wasm

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/horizen-pes/pkg/common"
	appCommon "github.com/horizen-pes/pkg/wasm/common"
)

type ApplicationModule struct {
	module   *wasmtime.Module
	instance *wasmtime.Instance
	memory   *wasmtime.Memory
}

// WasmtimeRuntime implements the Runtime interface using wasmtime-go
type WasmtimeRuntime struct {
	engine     *wasmtime.Engine
	store      *wasmtime.Store
	modules    map[common.ApplicationIdType]*ApplicationModule // Map of application ID to module
	moduleLock sync.RWMutex                  // Lock for module access
}

// NewWasmtimeRuntime creates a new wasmtime runtime instance
func NewWasmtimeRuntime() *WasmtimeRuntime {
	log.Println("Runtime: Initializing wasmtime runtime")

	// Create a new engine with default configuration
	engine := wasmtime.NewEngine()

	// Create a new store
	store := wasmtime.NewStore(engine)

	return &WasmtimeRuntime{
		engine:  engine,
		store:   store,
		modules: make(map[common.ApplicationIdType]*ApplicationModule),
	}
}

// Helper function to write data to WASM memory and return pointer
func (r *WasmtimeRuntime) writeToMemory(module *ApplicationModule, data []byte) (int32, error) {
	if module == nil {
		return 0, fmt.Errorf("module is nil")
	}
	if module.memory == nil {
		return 0, fmt.Errorf("memory not initialized")
	}

	// If data is empty, we don't need to allocate memory.
	// An empty payload might be a valid input for some operations, therefore
	// we return a null pointer, and the guest will see a length of 0.
	if len(data) == 0 {
		return 0, nil
	}

	// Get the allocate function from WASM
	allocateFunc := module.instance.GetFunc(r.store, "allocate")
	if allocateFunc == nil {
		return 0, fmt.Errorf("allocate function not found in WASM module")
	}

	// Call the allocate function to get a pointer
	result, err := allocateFunc.Call(r.store, int32(len(data)))
	if err != nil {
		return 0, fmt.Errorf("failed to call allocate: %w", err)
	}

	ptr, ok := result.(int32)
	if !ok {
		return 0, fmt.Errorf("allocate returned unexpected type: %T", result)
	}

	// Get memory data and copy our data to the allocated location
	// We can safely assume that the guest GC will not free the memory if is not locally referenced
	// since guest execution is paused while host is executing
	memData := module.memory.UnsafeData(r.store)
	if len(data) > math.MaxInt32 {
		return 0, fmt.Errorf("data too large for WASM memory: len(data) = %d > math.MaxInt32 = %d", len(data), math.MaxInt32)
	}

	if ptr < 0 || int64(ptr)+int64(len(data)) > int64(len(memData)) {
		return 0, fmt.Errorf("invalid memory allocation")
	}

	copy(memData[ptr:int(ptr)+len(data)], data)

	// TODO: The risk here is if the caller expects the pointer to remain valid across async host calls
	// without reference tracking in the guest.
	return ptr, nil
}

// Helper function to read data from WASM memory
func (r *WasmtimeRuntime) readFromMemory(module *ApplicationModule, ptr int32, length int32) ([]byte, error) {
	if module == nil {
		return nil, fmt.Errorf("module is nil")
	}
	if module.memory == nil {
		return nil, fmt.Errorf("memory not initialized")
	}

	memData := module.memory.UnsafeData(r.store)
	if ptr < 0 || length < 0 || int(ptr+length) > len(memData) {
		return nil, fmt.Errorf("invalid memory access")
	}

	result := make([]byte, length)
	copy(result, memData[ptr:ptr+length])
	return result, nil
}

func (r *WasmtimeRuntime) getOrLoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) (*ApplicationModule, error) {
	// Check if the module is already loaded
	if module, exists := r.modules[appId]; exists {
		return module, nil
	}

	// If not loaded, load the module
	_, err := r.LoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to load module: %w", err)
	}

	// Module should now be loaded, return it
	if module, exists := r.modules[appId]; exists {
		return module, nil
	}
	return nil, fmt.Errorf("module not found after loading: %d", appId)
}

// LoadModule loads a WASM module and returns initial state and state root
func (r *WasmtimeRuntime) LoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, error) {
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()
	log.Printf("Wasmtime Runtime: Loading WASM module for application %d (wasm size: %d bytes)", appId, len(wasm))

	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, err
	}

	// Compile the WASM module
	module, err := wasmtime.NewModule(r.engine, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to compile WASM module: %w", err)
	}

	// Create WASI configuration and linker for TinyGo WASI imports
	linker := wasmtime.NewLinker(r.engine)
	err = linker.DefineWasi()
	if err != nil {
		return nil, fmt.Errorf("failed to define WASI: %w", err)
	}

	// Create WASI configuration
	wasiConfig := wasmtime.NewWasiConfig()
	r.store.SetWasi(wasiConfig)

	// Create an instance of the module with WASI imports
	instance, err := linker.Instantiate(r.store, module)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	// Get the memory export
	memoryExport := instance.GetExport(r.store, "memory")
	if memoryExport == nil {
		return nil, fmt.Errorf("memory export not found in WASM module")
	}

	memory := memoryExport.Memory()
	if memory == nil {
		return nil, fmt.Errorf("memory export is not a memory")
	}

	// Get the load_module function
	loadModuleFunc := instance.GetFunc(r.store, "load_module")
	if loadModuleFunc == nil {
		return nil, fmt.Errorf("load_module function not found in WASM module")
	}

	appModule := &ApplicationModule{
		module:   module,
		instance: instance,
		memory:   memory,
	}


	// Call the load_module function
	// Wasm supports only int64, so we cast appId to int64
	result, err := loadModuleFunc.Call(r.store, wasmAppId)
	if err != nil {
		return nil, fmt.Errorf("failed to call load_module: %w", err)
	}

	// Extract the result bytes
	stateBytes, err := r.extractResultBytes(result, appModule)
	if err != nil {
		return nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	log.Printf("Wasmtime Runtime: Raw result from WASM: %s", string(stateBytes))

	log.Printf("Wasmtime Runtime: Successfully loaded WASM module for application %d", appId)

	// Store the module in the runtime
	r.modules[appId] = appModule
	return stateBytes, nil
}

// Deposit processes a deposit
func (r *WasmtimeRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender string, value *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, error) {
	log.Printf("Wasmtime Runtime: Processing deposit for application %d (value: %v wei for sender: %s)", appId, value, sender)

	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, nil, err
	}

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Get the deposit function
	depositFunc := appModule.instance.GetFunc(r.store, "deposit")
	if depositFunc == nil {
		return nil, nil, fmt.Errorf("deposit function not found in WASM module")
	}


	senderBytes := []byte(sender)
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write sender to memory: %w", err)
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write state to memory: %w", err)
	}

	if value == nil {
		return nil, nil, fmt.Errorf("deposit value cannot be nil")
	}
	
	if value.Sign() < 0 {
		return nil, nil, fmt.Errorf("deposit value cannot be negative")
	}

	valueBytes := value.Bytes()
	valuePtr, err := r.writeToMemory(appModule, valueBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write value to memory: %w", err)
	}

	// Call the deposit function

	// Wasm supports only int64, so we cast appId to int64
	result, err := depositFunc.Call(r.store, wasmAppId, senderPtr, int32(len(senderBytes)), valuePtr, int32(len(valueBytes)), statePtr, int32(len(state)))
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
func (r *WasmtimeRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender string, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, error) {
	log.Printf("Wasmtime Runtime: Processing request for application %d (payload size: %d, state size: %d)", appId, len(payload), len(state))

	wasmAppId, err := ToWasmType(appId)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(payload) == 0 {
		log.Printf("Wasmtime Runtime: Empty payload for application %d, returning current state", appId)
		return state, nil, nil, nil
	}

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Get the process_request function
	processRequestFunc := appModule.instance.GetFunc(r.store, "process_request")
	if processRequestFunc == nil {
		return nil, nil, nil, fmt.Errorf("process_request function not found in WASM module")
	}

	senderBytes := []byte(sender)
	senderPtr, err := r.writeToMemory(appModule, senderBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write sender to memory: %w", err)
	}

	payloadPtr, err := r.writeToMemory(appModule, payload)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write payload to memory: %w", err)
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write state to memory: %w", err)
	}

	// Call the process_request function
	// Wasm supports only int64, so we cast appId to int64
	result, err := processRequestFunc.Call(r.store, wasmAppId, senderPtr, int32(len(senderBytes)), payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
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

	log.Printf("Wasmtime Runtime: Successfully processed request for application %d, generated %d events and %d withdrawals", appId, len(processResult.Events), len(processResult.Withdrawals))
	return processResult.State, processResult.Events, processResult.Withdrawals, nil
}

// GenerateDeanonymizationReport generates a deanonymization report
func (r *WasmtimeRuntime) GenerateDeanonymizationReport(ctx context.Context, appId common.ApplicationIdType, payload []byte, state []byte, wasm []byte) ([]byte, error) {
	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Call the WASM function to generate the report
	generateReportFunc := appModule.instance.GetFunc(r.store, "generate_deanonymization_report")
	if generateReportFunc == nil {
		return nil, fmt.Errorf("function generate_deanonymization_report not found in wasm module")
	}

	payloadPtr, err := r.writeToMemory(appModule, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to write payload to memory: %w", err)
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, fmt.Errorf("failed to write state to memory: %w", err)
	}

	// Call the generate_deanonymization_report function
	result, err := generateReportFunc.Call(r.store, payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
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

	log.Printf("Wasmtime Runtime: Successfully generated deanonymization report for application %d", appId)
	return deanonymizationResult.Report, nil
}

func (r *WasmtimeRuntime) extractResultBytes(result interface{}, appModule *ApplicationModule) ([]byte, error) {
	// Extract the result pointer
	resultPtr, ok := result.(int32)
	if !ok {
		return nil, fmt.Errorf("wasm module returned unexpected type")
	}

	if resultPtr == 0 {
		return nil, fmt.Errorf("wasm module returned a null pointer, possibly due to an allocation failure or invalid input")
	}

	// Read the length prefix (first 4 bytes as uint32)
	memData := appModule.memory.UnsafeData(r.store)
	if resultPtr < 0 || int(resultPtr+4) > len(memData) {
		return nil, fmt.Errorf("invalid memory access for length prefix")
	}

	// Extract length from first 4 bytes (little-endian)
	lengthBytes := memData[resultPtr : resultPtr+4]
	dataLength := binary.LittleEndian.Uint32(lengthBytes)

	if dataLength == 0 {
		return nil, fmt.Errorf("empty result from wasm module")
	}

	// Validate that we can read the full data
	if int(resultPtr+4+int32(dataLength)) > len(memData) {
		return nil, fmt.Errorf("invalid memory access for data payload")
	}

	// Read the actual data (skip the 4-byte length prefix)
	resultBytes, err := r.readFromMemory(appModule, resultPtr+4, int32(dataLength))
	if err != nil {
		return nil, fmt.Errorf("failed to read wasm module result from memory: %w", err)
	}

	// WasmSerializationError (`{}`) is returned by the wasm module if serialization fails
	if string(resultBytes) == appCommon.WasmSerializationError {
		return nil, fmt.Errorf("wasm module failed to serialize response/error")
	}

	return resultBytes, nil
}

// Close closes the wasmtime runtime and cleans up resources
func (r *WasmtimeRuntime) Close() error {
	log.Printf("Wasmtime Runtime: Closing wasmtime runtime")

	// Cleanup resources
	// Memory cleanup is handled by Go's garbage collector
	if r.store != nil {
		r.store = nil
	}

	if r.engine != nil {
		r.engine = nil
	}

	for appId, module := range r.modules {
		if module.instance != nil {
			module.instance = nil
		}
		if module.module != nil {
			module.module = nil
		}
		if module.memory != nil {
			module.memory = nil
		}
		delete(r.modules, appId)
	}

	log.Printf("Wasmtime Runtime: Wasmtime runtime closed successfully")
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