package wasm

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/executor"
	"log"
	"sync"
)

const WasmSerializationError = "{}"

type ApplicationModule struct {
	module   *wasmtime.Module
	instance *wasmtime.Instance
	memory   *wasmtime.Memory
}

// WasmtimeRuntime implements the Runtime interface using wasmtime-go
type WasmtimeRuntime struct {
	engine     *wasmtime.Engine
	store      *wasmtime.Store
	modules    map[string]*ApplicationModule // Map of application ID to module
	moduleLock sync.RWMutex                  // Lock for module access
}

// DepositResult represents the result of a deposit operation from WASM
type DepositResult struct {
	State  []byte                `json:"state"`
	Events []executor.PlainEvent `json:"events"`
	Error  string                `json:"error,omitempty"`
}

// ProcessResult represents the result of a process request operation from WASM
type ProcessResult struct {
	State       []byte                `json:"state"`
	Events      []executor.PlainEvent `json:"events"`
	Withdrawals []common.Withdrawal   `json:"withdrawals"`
	Error       string                `json:"error,omitempty"`
}

// DeanonymizationResult represents the result of a process request operation from WASM
type DeanonymizationResult struct {
	Report []byte `json:"report"`
	Error  string `json:"error,omitempty"`
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
		modules: make(map[string]*ApplicationModule),
	}
}

// Helper function to write data to WASM memory and return pointer
func (r *WasmtimeRuntime) writeToMemory(module *ApplicationModule, data []byte) (int32, error) {
	if module.memory == nil {
		return 0, fmt.Errorf("memory not initialized")
	}

	// Get the allocate function from WASM
	allocateFunc := module.instance.GetFunc(r.store, "allocate")
	if allocateFunc == nil {
		return 0, fmt.Errorf("allocate function not found in WASM module")
	}

	// Call the allocate function to get a pointer
	result, err := allocateFunc.Call(r.store, len(data))
	if err != nil {
		return 0, fmt.Errorf("failed to call allocate: %w", err)
	}

	ptr, ok := result.(int32)
	if !ok {
		return 0, fmt.Errorf("allocate returned unexpected type")
	}

	// Get memory data and copy our data to the allocated location
	memData := module.memory.UnsafeData(r.store)
	if ptr < 0 || int(ptr+int32(len(data))) > len(memData) {
		return 0, fmt.Errorf("invalid memory allocation")
	}

	copy(memData[ptr:ptr+int32(len(data))], data)

	return ptr, nil
}

// Helper function to read data from WASM memory
func (r *WasmtimeRuntime) readFromMemory(module *ApplicationModule, ptr int32, length int32) ([]byte, error) {
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

func (r *WasmtimeRuntime) getOrLoadModule(ctx context.Context, appId string, wasm []byte) (*ApplicationModule, error) {
	// Check if the module is already loaded
	if module, exists := r.modules[appId]; exists {
		return module, nil
	}

	// If not loaded, load the module
	_, _, err := r.LoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to load module: %w", err)
	}

	// Module should now be loaded, return it
	if module, exists := r.modules[appId]; exists {
		return module, nil
	}
	return nil, fmt.Errorf("module not found after loading: %s", appId)
}

// LoadModule loads a WASM module and returns initial state and state root
func (r *WasmtimeRuntime) LoadModule(ctx context.Context, appId string, wasm []byte) ([]byte, []byte, error) {
	r.moduleLock.Lock()
	defer r.moduleLock.Unlock()
	log.Printf("Wasmtime Runtime: Loading WASM module for application %s (wasm size: %d bytes)", appId, len(wasm))

	// Compile the WASM module
	module, err := wasmtime.NewModule(r.engine, wasm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compile WASM module: %w", err)
	}

	// Create WASI configuration and linker for TinyGo WASI imports
	linker := wasmtime.NewLinker(r.engine)
	err = linker.DefineWasi()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to define WASI: %w", err)
	}

	// Create WASI configuration
	wasiConfig := wasmtime.NewWasiConfig()
	r.store.SetWasi(wasiConfig)

	// Create an instance of the module with WASI imports
	instance, err := linker.Instantiate(r.store, module)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	// Get the memory export
	memoryExport := instance.GetExport(r.store, "memory")
	if memoryExport == nil {
		return nil, nil, fmt.Errorf("memory export not found in WASM module")
	}

	memory := memoryExport.Memory()
	if memory == nil {
		return nil, nil, fmt.Errorf("memory export is not a memory")
	}

	// Get the load_module function
	loadModuleFunc := instance.GetFunc(r.store, "load_module")
	if loadModuleFunc == nil {
		return nil, nil, fmt.Errorf("load_module function not found in WASM module")
	}

	appModule := &ApplicationModule{
		module:   module,
		instance: instance,
		memory:   memory,
	}

	// Write appId to memory
	appIdBytes := []byte(appId)
	appIdPtr, err := r.writeToMemory(appModule, appIdBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write appId to memory: %w", err)
	}

	// Call the load_module function
	result, err := loadModuleFunc.Call(r.store, appIdPtr, int32(len(appIdBytes)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call load_module: %w", err)
	}

	// Extract the result bytes
	stateBytes, err := r.extractResultBytes(result, appModule, err)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	log.Printf("Wasmtime Runtime: Raw result from WASM: %s", string(stateBytes))

	// Create state root hash
	stateRoot := sha256.Sum256(stateBytes)

	log.Printf("Wasmtime Runtime: Successfully loaded WASM module for application %s", appId)

	// Store the module in the runtime
	r.modules[appId] = appModule
	return stateBytes, stateRoot[:], nil
}

// Deposit processes a deposit
func (r *WasmtimeRuntime) Deposit(ctx context.Context, appId string, sender string, value int64, state []byte, wasm []byte) ([]byte, []executor.PlainEvent, error) {
	log.Printf("Wasmtime Runtime: Processing deposit for application %s (value: %d wei for sender: %s)", appId, value, sender)

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Get the deposit function
	depositFunc := appModule.instance.GetFunc(r.store, "deposit")
	if depositFunc == nil {
		return nil, nil, fmt.Errorf("deposit function not found in WASM module")
	}

	// Write parameters to memory
	appIdBytes := []byte(appId)
	appIdPtr, err := r.writeToMemory(appModule, appIdBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write appId to memory: %w", err)
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

	// Call the deposit function
	result, err := depositFunc.Call(r.store, appIdPtr, int32(len(appIdBytes)), senderPtr, int32(len(senderBytes)), value, statePtr, int32(len(state)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call deposit: %w", err)
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule, err)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	log.Printf("Wasmtime Runtime: Raw deposit result from WASM: %s", string(resultBytes))

	// Deserialize the result
	var depositResult DepositResult
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
func (r *WasmtimeRuntime) ProcessRequest(ctx context.Context, appId string, sender string, payload []byte, state []byte, wasm []byte) ([]byte, []executor.PlainEvent, []common.Withdrawal, error) {
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
	processRequestFunc := appModule.instance.GetFunc(r.store, "process_request")
	if processRequestFunc == nil {
		return nil, nil, nil, fmt.Errorf("process_request function not found in WASM module")
	}

	// Write parameters to memory
	appIdBytes := []byte(appId)
	appIdPtr, err := r.writeToMemory(appModule, appIdBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to write appId to memory: %w", err)
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
	result, err := processRequestFunc.Call(r.store, appIdPtr, int32(len(appIdBytes)), senderPtr, int32(len(senderBytes)), payloadPtr, int32(len(payload)), statePtr, int32(len(state)))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to call process_request: %w", err)
	}

	// Extract the result bytes
	resultBytes, err := r.extractResultBytes(result, appModule, err)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	// Deserialize the result
	var processResult ProcessResult
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
func (r *WasmtimeRuntime) GenerateDeanonymizationReport(ctx context.Context, appId string, requestId string, payload []byte, state []byte, wasm []byte) ([]byte, error) {
	log.Printf("Wasmtime Runtime: Generating deanonymization report, id: %s, for application %s", requestId, appId)

	appModule, err := r.getOrLoadModule(ctx, appId, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to get or load module: %w", err)
	}

	// Get the generate_deanonymization_report function
	generateReportFunc := appModule.instance.GetFunc(r.store, "generate_deanonymization_report")
	if generateReportFunc == nil {
		return nil, fmt.Errorf("generate_deanonymization_report function not found in WASM module")
	}

	// Write parameters to memory
	appIdBytes := []byte(appId)
	appIdPtr, err := r.writeToMemory(appModule, appIdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write appId to memory: %w", err)
	}

	requestIdBytes := []byte(requestId)
	requestIdPtr, err := r.writeToMemory(appModule, requestIdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write requestId to memory: %w", err)
	}

	statePtr, err := r.writeToMemory(appModule, state)
	if err != nil {
		return nil, fmt.Errorf("failed to write state to memory: %w", err)
	}

	// Call the generate_deanonymization_report function
	result, err := generateReportFunc.Call(r.store, appIdPtr, int32(len(appIdBytes)), requestIdPtr, int32(len(requestIdBytes)), statePtr, int32(len(state)))
	if err != nil {
		return nil, fmt.Errorf("failed to call generate_deanonymization_report: %w", err)
	}

	// Extract the result bytes
	reportBytes, err := r.extractResultBytes(result, appModule, err)
	if err != nil {
		return nil, fmt.Errorf("failed to extract wasm module result bytes: %w", err)
	}

	// Deserialize the result
	var deanonymizationResult DeanonymizationResult
	if err := json.Unmarshal(reportBytes, &deanonymizationResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deanonymization result: %w", err)
	}

	if deanonymizationResult.Error != "" {
		return nil, fmt.Errorf("deanonymization report failed, wasm module error: %v", deanonymizationResult.Error)
	}

	log.Printf("Wasmtime Runtime: Successfully generated deanonymization report for application %s", appId)
	return deanonymizationResult.Report, nil
}

func (r *WasmtimeRuntime) extractResultBytes(result interface{}, appModule *ApplicationModule, err error) ([]byte, error) {
	// Extract the result pointer
	resultPtr, ok := result.(int32)
	if !ok {
		return nil, fmt.Errorf("wasm module returned unexpected type")
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
	if string(resultBytes) == WasmSerializationError {
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
