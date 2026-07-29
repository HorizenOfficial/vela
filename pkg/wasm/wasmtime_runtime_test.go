package wasm

import (
	"context"
	"encoding/binary"
	"math"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/logger"
	appCommon "github.com/HorizenOfficial/vela/pkg/wasm/common"
	wasmtime "github.com/bytecodealliance/wasmtime-go/v42"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The WAT fixtures below are minimal modules assembled at test time to drive the
// host through behaviour a correct guest never exhibits (traps, unparseable results,
// missing or mismatched exports, out-of-range pointers). The ABI they implement, and
// why this layer exists alongside the TinyGo integration tests in app/simple, are
// documented in docs/design/WASM_HOST_ABI.md.

// dispatchTestWat is a minimal WAT module used by TestProcessRequest_DispatchesByRequestType.
//
// Memory layout (all offsets within one 64 KiB page):
//
//	[100..173] process_request result: 4-byte LE length (70) + JSON {"state":[1],...}
//	[200..273] trusted_request result: 4-byte LE length (70) + JSON {"state":[2],...}
//	[300..328] load_module result:     4-byte LE length (25) + JSON {"state":[],"fuel":"0x1"}
//	[500..~  ] scratch area for allocate — the host writes payload/state here; the
//	           WASM functions ignore all input params and return the fixed offsets above.
//
// Signatures:
//
//	process_request (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
//	  (appId, senderPtr, senderLen, requestType, payloadPtr, payloadLen, statePtr, stateLen)
//	trusted_request (param i64 i32 i32 i32 i32) (result i32)
//	  (appId, payloadPtr, payloadLen, statePtr, stateLen)
//	load_module (param i64) (result i32)       — called by getOrLoadModule on first load
//	allocate    (param i32) (result i32)       — always returns scratch offset 500
//	deallocate  (param i32 i32)                — no-op
const dispatchTestWat = `(module
  (memory (export "memory") 1)

  ;; process_request result at offset 100:
  ;;   4-byte LE length = 70 (0x46) followed by JSON {"state":[1],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1"}
  (data (i32.const 100)
    "\46\00\00\00"
    "\7b\22\73\74\61\74\65\22\3a\5b\31\5d\2c\22\65\76\65\6e\74\73\22\3a\5b\5d\2c"
    "\22\61\70\70\45\76\65\6e\74\73\22\3a\5b\5d\2c\22\77\69\74\68\64\72\61\77\61"
    "\6c\73\22\3a\5b\5d\2c\22\66\75\65\6c\22\3a\22\30\78\31\22\7d"
  )

  ;; trusted_request result at offset 200:
  ;;   4-byte LE length = 70 (0x46) followed by JSON {"state":[2],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1"}
  (data (i32.const 200)
    "\46\00\00\00"
    "\7b\22\73\74\61\74\65\22\3a\5b\32\5d\2c\22\65\76\65\6e\74\73\22\3a\5b\5d\2c"
    "\22\61\70\70\45\76\65\6e\74\73\22\3a\5b\5d\2c\22\77\69\74\68\64\72\61\77\61"
    "\6c\73\22\3a\5b\5d\2c\22\66\75\65\6c\22\3a\22\30\78\31\22\7d"
  )

  ;; load_module result at offset 300:
  ;;   4-byte LE length = 25 (0x19) followed by JSON {"state":[],"fuel":"0x1"}
  (data (i32.const 300)
    "\19\00\00\00"
    "\7b\22\73\74\61\74\65\22\3a\5b\5d\2c\22\66\75\65\6c\22\3a\22\30\78\31\22\7d"
  )

  ;; load_module: called by getOrLoadModule during module warm-up; returns fixed offset 300
  (func (export "load_module") (param i64) (result i32)
    i32.const 300
  )

  ;; allocate: always returns scratch offset 500 (host writes payload/state here)
  (func (export "allocate") (param i32) (result i32)
    i32.const 500
  )

  ;; deallocate: no-op
  (func (export "deallocate") (param i32 i32)
  )

  ;; process_request: 8-arg ABI — ignores all params, returns fixed offset 100
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100
  )

  ;; trusted_request: 5-arg ABI (no sender, no request_type) — returns fixed offset 200
  (func (export "trusted_request") (param i64 i32 i32 i32 i32) (result i32)
    i32.const 200
  )
)`

// spinningProcessRequestWat is dispatchTestWat reduced to the process_request
// path, with a busy loop inside the export so the guest call takes long enough to
// still be running when another goroutine evicts the module.
// Used by TestConcurrentExecutionAndEvictionIsSafe: with an instantaneous export
// the eviction never lands inside the call and the test passes even when the
// serialization it checks is removed.
const spinningProcessRequestWat = `(module
  (memory (export "memory") 1)

  ;; process_request result at offset 100 (same layout as dispatchTestWat)
  (data (i32.const 100)
    "\46\00\00\00"
    "\7b\22\73\74\61\74\65\22\3a\5b\31\5d\2c\22\65\76\65\6e\74\73\22\3a\5b\5d\2c"
    "\22\61\70\70\45\76\65\6e\74\73\22\3a\5b\5d\2c\22\77\69\74\68\64\72\61\77\61"
    "\6c\73\22\3a\5b\5d\2c\22\66\75\65\6c\22\3a\22\30\78\31\22\7d"
  )

  ;; load_module result at offset 300
  (data (i32.const 300)
    "\19\00\00\00"
    "\7b\22\73\74\61\74\65\22\3a\5b\5d\2c\22\66\75\65\6c\22\3a\22\30\78\31\22\7d"
  )

  (func (export "load_module") (param i64) (result i32) i32.const 300)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))

  ;; process_request: spins, then returns the fixed result offset 100
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    (local $i i32)
    (local.set $i (i32.const 20000000))
    (block $done
      (loop $spin
        (br_if $done (i32.eqz (local.get $i)))
        (local.set $i (i32.sub (local.get $i) (i32.const 1)))
        (br $spin)
      )
    )
    i32.const 100
  )
)`

// trappingProcessRequestWat loads successfully but traps (unreachable) inside
// process_request, the shape a TinyGo panic takes — including the out-of-memory
// panic the guest memory cap makes reachable. Used by TestGuestTrapEvictsModule.
const trappingProcessRequestWat = `(module
  (memory (export "memory") 1)

  ;; load_module result at offset 300: 4-byte LE length (25) + JSON
  (data (i32.const 300) "\19\00\00\00" "{\"state\":[],\"fuel\":\"0x1\"}")

  (func (export "load_module") (param i64) (result i32) i32.const 300)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))

  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    unreachable
  )
)`

// appErrorProcessRequestWat returns a well-formed ProcessResult carrying an
// application-level error. That is a normal rejection by a healthy guest, not a
// fault, so the module must stay cached. Used by TestAppErrorKeepsModuleCached.
const appErrorProcessRequestWat = `(module
  (memory (export "memory") 1)

  ;; process_request result at offset 100: 4-byte LE length (84 = 0x54) + JSON
  ;; {"state":[],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1","error":"boom"}
  (data (i32.const 100) "\54\00\00\00"
    "{\"state\":[],\"events\":[],\"appEvents\":[],\"withdrawals\":[],\"fuel\":\"0x1\",\"error\":\"boom\"}")

  ;; load_module result at offset 300: 4-byte LE length (25) + JSON
  (data (i32.const 300) "\19\00\00\00" "{\"state\":[],\"fuel\":\"0x1\"}")

  (func (export "load_module") (param i64) (result i32) i32.const 300)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))

  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100
  )
)`

// noAllocateWat loads successfully but has no `allocate` export, so the host
// cannot write the request into guest memory. That is a static defect in the
// module, not a guest fault: re-instantiating cannot fix it, so the module must
// stay cached instead of being recompiled on every request forever.
// Reachable in practice because writeToMemory returns early for empty data, so
// such a module can pass a deploy that has no constructor params.
// Used by TestHostSideFailureKeepsModuleCached.
const noAllocateWat = `(module
  (memory (export "memory") 1)

  ;; load_module result at offset 300: 4-byte LE length (25) + JSON
  (data (i32.const 300) "\19\00\00\00" "{\"state\":[],\"fuel\":\"0x1\"}")

  (func (export "load_module") (param i64) (result i32) i32.const 300)

  ;; no allocate export, no deallocate export

  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100
  )
)`

// malformedResultWat returns a well-formed length prefix followed by bytes that
// are not valid JSON, so the guest runs to completion and the host fails only when
// deserializing. Used by TestMalformedResultEvictsModule.
const malformedResultWat = `(module
  (memory (export "memory") 1)

  ;; process_request result at offset 100: 4-byte LE length (15 = 0x0f) + non-JSON
  (data (i32.const 100) "\0f\00\00\00" "not json at all")

  ;; load_module result at offset 300: 4-byte LE length (25) + JSON
  (data (i32.const 300) "\19\00\00\00" "{\"state\":[],\"fuel\":\"0x1\"}")

  (func (export "load_module") (param i64) (result i32) i32.const 300)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))

  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100
  )
)`

// wrongAritySignatureWat exports process_request with 2 parameters where the host
// ABI passes 8. Func.Call validates argument count before entering wasm, so the
// call fails without the guest ever running — a static defect in the module, not a
// guest fault. Used by TestSignatureMismatchKeepsModuleCached.
//
// Reachable in practice: only the deploy export's signature is exercised at deploy
// time, so an app built against a stale host ABI deploys fine and then fails on
// every request.
const wrongAritySignatureWat = `(module
  (memory (export "memory") 1)

  ;; load_module result at offset 300: 4-byte LE length (25) + JSON
  (data (i32.const 300) "\19\00\00\00" "{\"state\":[],\"fuel\":\"0x1\"}")

  (func (export "load_module") (param i64) (result i32) i32.const 300)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))

  ;; host passes 8 arguments; this declares 2
  (func (export "process_request") (param i64 i32) (result i32) i32.const 100)
)`

// nearEndAllocWat has an `allocate` that returns a pointer 4 bytes from the end of
// a one-page memory, so the 16-byte stats read runs past the end.
// Used by TestStatsRejectsOutOfBoundsResultPointer.
const nearEndAllocWat = `(module
  (memory (export "memory") 1)

  ;; load_module result at offset 300: 4-byte LE length (25) + JSON
  (data (i32.const 300) "\19\00\00\00" "{\"state\":[],\"fuel\":\"0x1\"}")

  (func (export "load_module") (param i64) (result i32) i32.const 300)
  (func (export "deallocate") (param i32 i32))

  ;; 65532 = one page (65536) minus 4, so a 16-byte read from here overruns
  (func (export "allocate") (param i32) (result i32) i32.const 65532)
  (func (export "get_allocated_memory_stats") (param i32))
)`

var testLogger logger.Logger

func TestMain(m *testing.M) {
	// Initialize once, by default it writes on stderr
	//	testLogger = logger.NewLogger(&logger.Config{Kind: "printf"})

	testLogger = logger.NewLogger(
		&logger.Config{
			Kind:         "zerolog",
			ConsoleColor: false, // colors can print escape chars on tty
			Console:      true,
			ConsoleLevel: "trace",
			//FileName:     "qqq.log",
			//FileLevel:    "info",
		},
	)
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestWriteToMemory_NilModule(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	_, err := runtime.writeToMemory(nil, []byte("some data"))
	require.Error(t, err)
	require.Equal(t, "module is nil", err.Error())
}

func TestWriteToMemory_NilModuleStore(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	// Create a per-module store
	store := wasmtime.NewStore(runtime.engine)

	// Create a valid memory object, but instance can be nil
	// since it's not reached when data is empty.
	memType, err := wasmtime.NewMemoryType(1, false, 0, false)
	require.NoError(t, err)
	memory, err := wasmtime.NewMemory(store, memType)
	require.NoError(t, err)

	module := &ApplicationModule{
		memory:   memory,
		instance: nil, // instance is not used when len(data) == 0
	}

	_, err = runtime.writeToMemory(module, []byte("some data"))
	require.Error(t, err)
	require.Equal(t, "wasm module has a nil store", err.Error())
}

func TestWriteToMemory_NilMemory(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	module := &ApplicationModule{} // memory is nil by default

	_, err := runtime.writeToMemory(module, []byte("some data"))
	require.Error(t, err)
	require.Equal(t, "memory not initialized", err.Error())
}

func TestWriteToMemory_NilData(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	// Create a per-module store
	store := wasmtime.NewStore(runtime.engine)

	// Create a valid memory object, but instance can be nil
	// since it's not reached when data is empty.
	memType, err := wasmtime.NewMemoryType(1, false, 0, false)
	require.NoError(t, err)
	memory, err := wasmtime.NewMemory(store, memType)
	require.NoError(t, err)

	module := &ApplicationModule{
		store:    store,
		memory:   memory,
		instance: nil, // instance is not used when len(data) == 0
	}

	// Test with nil data
	ptr, err := runtime.writeToMemory(module, nil)
	require.NoError(t, err)
	require.Equal(t, int32(0), ptr)

	// Test with empty byte slice
	ptr, err = runtime.writeToMemory(module, []byte{})
	require.NoError(t, err)
	require.Equal(t, int32(0), ptr)
}

func TestExtractResultBytes(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	// Create a per-module store
	store := wasmtime.NewStore(runtime.engine)

	// Setup a mock memory for testing, using the runtime's store
	// 1 is the minimum size in WebAssembly pages (WebAssembly page size = 64 KiB (65,536 bytes))
	memType, err := wasmtime.NewMemoryType(1, false, 0, false)
	require.NoError(t, err)
	memory, err := wasmtime.NewMemory(store, memType)
	require.NoError(t, err)

	appModule := &ApplicationModule{store: store, memory: memory}

	// Write some test data to memory: 4 bytes for length, then the data
	testData := []byte("hello world")
	dataWithLen := make([]byte, 4+len(testData))
	binary.LittleEndian.PutUint32(dataWithLen, uint32(len(testData)))
	copy(dataWithLen[4:], testData)
	// [0x0b,0x00,0x00,0x00, 0x68,..,0x64]

	// Copy the prepared data into the wasm memory
	memSlice := memory.UnsafeData(store)
	testPtr := int32(100) // Arbitrary starting position, it is safe because we configures the slice to be 64K
	copy(memSlice[testPtr:], dataWithLen)

	t.Run("successful extraction", func(t *testing.T) {
		result, err := runtime.extractResultBytes(testPtr, appModule)
		require.NoError(t, err)
		require.Equal(t, testData, result)
	})

	t.Run("null pointer", func(t *testing.T) {
		_, err := runtime.extractResultBytes(int32(0), appModule)
		require.Error(t, err)
		require.Contains(t, err.Error(), "wasm module returned null pointer")
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := runtime.extractResultBytes("not a pointer", appModule)
		require.Error(t, err)
		require.Contains(t, err.Error(), "wasm module returned unexpected type")
	})

	t.Run("out of bounds pointer", func(t *testing.T) {
		_, err := runtime.extractResultBytes(int32(len(memSlice)), appModule)
		require.Error(t, err)
		require.Equal(t, "invalid memory access for length prefix", err.Error())
	})

	t.Run("serialization error", func(t *testing.T) {
		// Write the serialization error constant to memory
		errorData := []byte(appCommon.WasmSerializationError)
		errorWithLen := make([]byte, 4+len(errorData))
		binary.LittleEndian.PutUint32(errorWithLen, uint32(len(errorData)))
		copy(errorWithLen[4:], errorData)

		errorPtr := int32(200)
		copy(memSlice[errorPtr:], errorWithLen)

		_, err := runtime.extractResultBytes(errorPtr, appModule)
		require.Error(t, err)
		require.Equal(t, "wasm module failed to serialize response/error", err.Error())
	})
}

// TestToWasmType_SmallValue verifies that a normal application ID is correctly
// converted to int64.
func TestToWasmType_SmallValue(t *testing.T) {
	aid := common.NewApplicationId(42)
	result, err := ToWasmType(aid)
	require.NoError(t, err)
	assert.Equal(t, int64(42), result)
}

// TestToWasmType_MaxInt64 verifies that the maximum int64 value is accepted
// and the bit pattern is preserved.
func TestToWasmType_MaxInt64(t *testing.T) {
	aid := common.NewApplicationId(math.MaxInt64)
	result, err := ToWasmType(aid)
	require.NoError(t, err)
	assert.Equal(t, int64(math.MaxInt64), result)
}

// TestToWasmType_AboveMaxInt64 verifies that a uint64 value exceeding MaxInt64
// survives the round-trip through WASM's i64 (int64). WASM i64 has no signedness,
// so the bit pattern is preserved: uint64 → int64 → uint64 yields the original value.
func TestToWasmType_AboveMaxInt64(t *testing.T) {
	original := uint64(math.MaxInt64 + 1) // 0x8000000000000000
	aid := common.NewApplicationId(original)

	wasmVal, err := ToWasmType(aid)
	require.NoError(t, err)

	// The int64 value is negative, but the bit pattern is preserved
	assert.True(t, wasmVal < 0, "int64 reinterpretation should be negative")

	// Simulate what the guest does: reinterpret int64 back to uint64
	restored := uint64(wasmVal)
	assert.Equal(t, original, restored)
}

// TestToWasmType_MaxUint64 verifies that the maximum uint64 value survives the
// round-trip: MaxUint64 → int64(-1) → uint64(MaxUint64).
func TestToWasmType_MaxUint64(t *testing.T) {
	original := uint64(math.MaxUint64) // 0xFFFFFFFFFFFFFFFF
	aid := common.NewApplicationId(original)

	wasmVal, err := ToWasmType(aid)
	require.NoError(t, err)
	assert.Equal(t, int64(-1), wasmVal)

	// Simulate what the guest does: reinterpret int64 back to uint64
	restored := uint64(wasmVal)
	assert.Equal(t, original, restored)
}

func TestClose(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	// Add a dummy module to ensure the map is cleared
	runtime.modules[1] = &ApplicationModule{}
	runtime.touchModule(common.NewApplicationId(1))

	err := runtime.Close()
	require.NoError(t, err)

	require.Nil(t, runtime.engine)
	require.Empty(t, runtime.modules)
	require.Equal(t, 0, runtime.accessOrder.Len())
	require.Empty(t, runtime.accessElements)
}

// --- LRU Cache Eviction Tests ---

func TestLRUEviction_Basic(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 2)
	defer runtime.Close()

	app1 := common.NewApplicationId(1)
	app2 := common.NewApplicationId(2)
	app3 := common.NewApplicationId(3)

	// Insert 3 modules into a cache with limit 2
	runtime.moduleLock.Lock()
	runtime.modules[app1] = &ApplicationModule{}
	runtime.touchModule(app1)
	runtime.modules[app2] = &ApplicationModule{}
	runtime.touchModule(app2)
	runtime.modules[app3] = &ApplicationModule{}
	runtime.touchModule(app3)
	runtime.evictIfNeeded()
	runtime.moduleLock.Unlock()

	// app1 (least recent) should be evicted
	require.Equal(t, 2, len(runtime.modules))
	_, hasApp1 := runtime.modules[app1]
	_, hasApp2 := runtime.modules[app2]
	_, hasApp3 := runtime.modules[app3]
	require.False(t, hasApp1, "app1 should be evicted (LRU)")
	require.True(t, hasApp2, "app2 should remain")
	require.True(t, hasApp3, "app3 should remain")
	require.Equal(t, 2, runtime.accessOrder.Len())
}

func TestLRUEviction_TouchUpdatesOrder(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 2)
	defer runtime.Close()

	app1 := common.NewApplicationId(1)
	app2 := common.NewApplicationId(2)
	app3 := common.NewApplicationId(3)

	runtime.moduleLock.Lock()
	// Insert app1, then app2
	runtime.modules[app1] = &ApplicationModule{}
	runtime.touchModule(app1)
	runtime.modules[app2] = &ApplicationModule{}
	runtime.touchModule(app2)

	// Touch app1 again — now app2 is the LRU
	runtime.touchModule(app1)

	// Insert app3 — should evict app2 (not app1)
	runtime.modules[app3] = &ApplicationModule{}
	runtime.touchModule(app3)
	runtime.evictIfNeeded()
	runtime.moduleLock.Unlock()

	require.Equal(t, 2, len(runtime.modules))
	_, hasApp1 := runtime.modules[app1]
	_, hasApp2 := runtime.modules[app2]
	_, hasApp3 := runtime.modules[app3]
	require.True(t, hasApp1, "app1 should remain (was touched)")
	require.False(t, hasApp2, "app2 should be evicted (LRU)")
	require.True(t, hasApp3, "app3 should remain (newest)")
}

func TestLRUEviction_Unlimited(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0) // 0 = unlimited
	defer runtime.Close()

	runtime.moduleLock.Lock()
	for i := uint64(1); i <= 100; i++ {
		appId := common.NewApplicationId(i)
		runtime.modules[appId] = &ApplicationModule{}
		runtime.touchModule(appId)
		runtime.evictIfNeeded()
	}
	runtime.moduleLock.Unlock()

	require.Equal(t, 100, len(runtime.modules), "no eviction with unlimited cache")
	require.Equal(t, 100, runtime.accessOrder.Len())
}

func TestLRUEviction_UnloadRemovesFromAccessOrder(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	app1 := common.NewApplicationId(1)
	app2 := common.NewApplicationId(2)

	runtime.moduleLock.Lock()
	runtime.modules[app1] = &ApplicationModule{}
	runtime.touchModule(app1)
	runtime.modules[app2] = &ApplicationModule{}
	runtime.touchModule(app2)
	runtime.moduleLock.Unlock()

	require.Equal(t, 2, runtime.accessOrder.Len())

	err := runtime.UnloadModule(app1)
	require.NoError(t, err)

	require.Equal(t, 1, len(runtime.modules))
	require.Equal(t, 1, runtime.accessOrder.Len())
	_, hasApp1 := runtime.accessElements[app1]
	require.False(t, hasApp1)
}

// TestProcessRequest_DispatchesByRequestType asserts the TrustProcess/Process dispatch split
// introduced in Phase 11.1:
//   - TrustProcess → trusted_request export (state discriminator byte == 2)
//   - Process       → process_request export (state discriminator byte == 1)
//
// It compiles a real two-export WASM module from WAT (dispatchTestWat) where each export
// returns a different pre-seeded result buffer so we can distinguish which was called.
func TestProcessRequest_DispatchesByRequestType(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(dispatchTestWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	ctx := context.Background()
	payload := []byte("{}")
	state := []byte("{}")

	// TrustProcess -> trusted_request (state discriminator byte == 2)
	trustedState, _, _, _, _, _, failure := runtime.ProcessRequest(ctx, appId, ethCommon.Address{}, common.TrustProcess, payload, state, wasmBytes)
	require.Nil(t, failure, "TrustProcess must succeed via trusted_request")
	require.Equal(t, []byte{2}, trustedState, "TrustProcess must dispatch to trusted_request (state=[2])")

	// Process -> process_request (state discriminator byte == 1)
	// Same appId is fine — the cached module exports both functions.
	procState, _, _, _, _, _, failure := runtime.ProcessRequest(ctx, appId, ethCommon.Address{}, common.Process, payload, state, wasmBytes)
	require.Nil(t, failure, "Process must succeed via process_request")
	require.Equal(t, []byte{1}, procState, "Process must dispatch to process_request (state=[1])")
}

// TestConcurrentExecutionAndEvictionIsSafe drives guest execution and LRU
// eviction from separate goroutines, which is the exact overlap execLock exists
// to prevent: ProcessRequest fetches its module under moduleLock but calls into
// it after releasing that lock, so without execLock a concurrent load for another
// app can pick the in-use module as the eviction victim and store.Close() it
// mid-call. That frees native wasmtime state while guest code is running on it —
// a segfault, not a Go panic, so no recover() would catch it and the enclave dies.
//
// maxCachedModules is 1, so every load by one worker evicts the module the other
// worker is using, and the guest export spins (spinningProcessRequestWat) so the
// eviction lands while the call is still running. The transport already dispatches
// each inbound message in its own goroutine (pkg/communication/shared_impl.go), so
// this is the shape a second in-flight request would take.
//
// This has been verified to be a real regression test: with the execLock
// acquisitions removed it crashes the test binary with SIGSEGV inside native
// wasmtime code, within the first few iterations. Note that it cannot be a Go
// test failure — the process dies — and that -race does not report it, since the
// conflict is on memory the Go race detector does not track.
func TestConcurrentExecutionAndEvictionIsSafe(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 1) // cache holds exactly one module
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(spinningProcessRequestWat)
	require.NoError(t, err)

	const (
		workers    = 2
		iterations = 50
	)

	ctx := context.Background()
	payload := []byte("{}")
	state := []byte("{}")

	var wg sync.WaitGroup
	for worker := 1; worker <= workers; worker++ {
		wg.Add(1)
		appId := common.NewApplicationId(uint64(worker))
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// assert (not require): require calls FailNow, which must not be
				// invoked from a non-test goroutine.
				gotState, _, _, _, _, _, failure := runtime.ProcessRequest(
					ctx, appId, ethCommon.Address{}, common.Process, payload, state, wasmBytes)
				if !assert.Nil(t, failure, "app %d iteration %d must succeed", appId, i) {
					return
				}
				assert.Equal(t, []byte{1}, gotState, "app %d iteration %d returned unexpected state", appId, i)
			}
		}()
	}
	wg.Wait()
}

// isCached reports whether a module for appId is currently in the LRU cache.
// A plain helper rather than a method, to avoid extending the production type
// from a test file.
func isCached(r *WasmtimeRuntime, appId common.ApplicationIdType) bool {
	r.moduleLock.RLock()
	defer r.moduleLock.RUnlock()
	_, exists := r.modules[appId]
	return exists
}

// TestGuestTrapEvictsModule verifies that a module whose guest trapped is dropped
// from the cache instead of serving the next request on a heap left in an unknown
// state (see evictFaultedModule).
//
// This also exercises the defer ordering the eviction depends on: the deallocate
// defers registered for the payload/state writes must run BEFORE the store is
// closed, otherwise they would call into freed memory.
func TestGuestTrapEvictsModule(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(trappingProcessRequestWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	ctx := context.Background()

	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	require.NotNil(t, failure, "a trapping guest must fail the request")
	require.False(t, isCached(runtime, appId), "a module whose guest trapped must not stay cached")

	// The app must still be usable: the next request reloads it and traps again
	// rather than, say, panicking on a closed store.
	_, _, _, _, _, _, failure = runtime.ProcessRequest(
		ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	require.NotNil(t, failure, "the reloaded module must trap again, not crash")
	require.False(t, isCached(runtime, appId))
}

// TestAppErrorKeepsModuleCached is the counterpart to TestGuestTrapEvictsModule:
// a guest that cleanly reports an application-level error is healthy, so evicting
// it would throw away a valid compiled module on every rejected request.
func TestAppErrorKeepsModuleCached(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(appErrorProcessRequestWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	ctx := context.Background()

	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	require.NotNil(t, failure, "the guest reported an application error")
	require.Contains(t, failure.Error(), "boom")
	require.True(t, isCached(runtime, appId), "an application-level error must not evict the module")
}

// TestStatsRejectsOutOfBoundsResultPointer verifies that a guest-supplied result
// pointer too close to the end of memory is rejected with an error rather than
// panicking the host.
//
// GetAllocatedMemoryStats reads 16 bytes straight out of the memory slice at the
// pointer the guest's `allocate` returned. Unlike extractResultBytes it has no
// recover shield, so an unchecked slice expression here takes the process down —
// against the "NEVER let guest crash the host" rule that shield documents.
func TestStatsRejectsOutOfBoundsResultPointer(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(nearEndAllocWat)
	require.NoError(t, err)

	_, _, err = runtime.GetAllocatedMemoryStats(context.Background(), common.NewApplicationId(1), wasmBytes)
	require.Error(t, err, "an out-of-bounds result pointer must be an error, not a panic")
	require.Contains(t, err.Error(), "out of bounds")
}

// TestSignatureMismatchKeepsModuleCached covers the case where the host cannot
// invoke the guest at all: Func.Call rejects the argument list before entering
// wasm, so no guest code runs and no heap is touched.
//
// By the rule in errGuestFault that makes it a static module defect, like a
// missing export — recompiling cannot change the signature, so evicting would
// recompile on every request for as long as the app keeps being called.
func TestSignatureMismatchKeepsModuleCached(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(wrongAritySignatureWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		_, _, _, _, _, _, failure := runtime.ProcessRequest(
			ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
		require.NotNil(t, failure, "a signature mismatch must fail the request")
		require.True(t, isCached(runtime, appId),
			"the guest never ran, so the module must stay cached (iteration %d)", i)
	}
}

// TestMalformedResultEvictsModule pins down the ambiguous case in the
// classification: a guest that returns bytes the host cannot deserialize is
// treated as a fault and evicted, even though a broken serializer would fail
// identically on every request (so the recompile is wasted).
//
// This is deliberate, not an oversight — the guest ran and wrote that buffer, so
// the host cannot tell a deterministic serialization defect from a heap bug that
// clobbered the result. See errGuestFault for the cost argument. Asserting it here
// so the choice cannot be flipped silently; if it is ever moved to the cached
// class, this test is the place to record why.
func TestMalformedResultEvictsModule(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(malformedResultWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	ctx := context.Background()

	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	require.NotNil(t, failure, "unparseable guest output must fail the request")
	require.False(t, isCached(runtime, appId), "a guest that produced unparseable output must be evicted")
}

// TestHostSideFailureKeepsModuleCached covers a failure that is neither a guest
// fault nor an application error, but a static defect in the module (here, a
// missing `allocate` export). One of the eviction-classification cases; the rule
// they all check is documented on errGuestFault. Deliberately not enumerated by
// count here — that goes stale every time a case is added. Evicting on those would recompile
// the module on every request for as long as the app keeps being called, which is
// wasted work and cheap amplification for anyone able to deploy such an app.
func TestHostSideFailureKeepsModuleCached(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(noAllocateWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		_, _, _, _, _, _, failure := runtime.ProcessRequest(
			ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
		require.NotNil(t, failure, "writing to guest memory must fail without an allocate export")
		require.True(t, isCached(runtime, appId),
			"a host-side failure must not evict the module (iteration %d)", i)
	}
}

// TestTryAcquireExecLock covers the bounded acquire that keeps Close from hanging
// behind a runaway guest.
func TestTryAcquireExecLock(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	// Free: acquired immediately.
	require.True(t, runtime.tryAcquireExecLock(time.Second))
	runtime.execLock.Unlock()

	// Held: gives up after the timeout instead of blocking forever.
	runtime.execLock.Lock()
	start := time.Now()
	require.False(t, runtime.tryAcquireExecLock(50*time.Millisecond))
	require.GreaterOrEqual(t, time.Since(start), 50*time.Millisecond, "must wait for the full timeout")
	runtime.execLock.Unlock()

	// Released while waiting: acquired rather than timing out.
	runtime.execLock.Lock()
	go func() {
		time.Sleep(20 * time.Millisecond)
		runtime.execLock.Unlock()
	}()
	require.True(t, runtime.tryAcquireExecLock(5*time.Second))
	runtime.execLock.Unlock()
}

// TestCloseDoesNotHangOnStuckGuest asserts that shutdown completes even when a
// guest call never returns. Guest execution is unbounded today (no fuel, no epoch
// deadline — see newPinnedEngine), and execLock is held across guest calls, so
// without the bounded acquire in Close this would hang forever and the executor
// could never shut down gracefully.
//
// A held execLock stands in for the runaway guest: it is the same state the
// runtime would be in, without needing to burn a core spinning inside wasmtime.
func TestCloseDoesNotHangOnStuckGuest(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)

	runtime.execLock.Lock() // simulate a guest call that never returns
	defer runtime.execLock.Unlock()

	done := make(chan error, 1)
	start := time.Now()
	go func() { done <- runtime.Close() }()

	select {
	case err := <-done:
		require.Error(t, err, "Close must report that it gave up rather than claiming success")
		require.GreaterOrEqual(t, time.Since(start), shutdownExecLockTimeout)
	case <-time.After(shutdownExecLockTimeout + 10*time.Second):
		t.Fatal("Close hung waiting for an in-flight guest call")
	}
}

func TestSetMaxGuestMemoryBytes(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	// The default is the 2 GiB ceiling
	require.Equal(t, int64(maxGuestMemoryCeilingBytes), runtime.GetMaxGuestMemoryBytes())

	// In-range values are applied as-is
	runtime.SetMaxGuestMemoryBytes(512 * 1024 * 1024)
	require.Equal(t, int64(512*1024*1024), runtime.GetMaxGuestMemoryBytes())

	// Out-of-range values fall back to the ceiling
	for _, invalid := range []int64{0, -1, maxGuestMemoryCeilingBytes + 1} {
		runtime.SetMaxGuestMemoryBytes(invalid)
		require.Equal(t, int64(maxGuestMemoryCeilingBytes), runtime.GetMaxGuestMemoryBytes())
	}
}

// TestGuestMemoryCapEnforced verifies that stores created by the runtime refuse
// guest memory growth beyond the configured cap: memory.grow reports failure
// (-1) to the guest instead of growing.
func TestGuestMemoryCapEnforced(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	const pageSize = 64 * 1024
	runtime.SetMaxGuestMemoryBytes(2 * pageSize)

	wasmBytes, err := wasmtime.Wat2Wasm(`(module
		(memory (export "memory") 1)
		(func (export "grow") (param i32) (result i32)
			local.get 0
			memory.grow))`)
	require.NoError(t, err)

	module, err := wasmtime.NewModule(runtime.engine, wasmBytes)
	require.NoError(t, err)
	defer module.Close()

	runtime.moduleLock.Lock()
	store := runtime.newModuleStore()
	runtime.moduleLock.Unlock()
	defer store.Close()

	instance, err := wasmtime.NewInstance(store, module, nil)
	require.NoError(t, err)

	grow := instance.GetFunc(store, "grow")
	require.NotNil(t, grow)

	// 1 -> 2 pages: within the cap, returns the previous size in pages
	res, err := grow.Call(store, int32(1))
	require.NoError(t, err)
	require.Equal(t, int32(1), res)

	// 2 -> 3 pages: beyond the cap, the guest sees a failed grow
	res, err = grow.Call(store, int32(1))
	require.NoError(t, err)
	require.Equal(t, int32(-1), res)
}

// TestPinnedEngineRejectsDisabledProposals verifies that the explicitly pinned
// feature set in newPinnedEngine is actually enforced: a module using a proposal
// this runtime disables must fail to compile rather than being silently accepted
// because a future wasmtime enables it by default.
//
// Both cases matter for the enclave RAM budget: the guest memory cap is applied
// per linear memory, so extra (or shared) memories would multiply the worst-case
// RAM a single app can hold. See newPinnedEngine and newModuleStore.
func TestPinnedEngineRejectsDisabledProposals(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	testCases := []struct {
		name string
		wat  string
	}{
		{
			name: "multi-memory",
			wat: `(module
				(memory (export "memory") 1)
				(memory 1))`,
		},
		{
			name: "shared memory (threads)",
			wat: `(module
				(memory (export "memory") 1 1 shared))`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Wat2Wasm is proposal-agnostic and encodes both of these fine; the
			// rejection must come from the engine's pinned feature set. Note that
			// wasmtime v42's *default* engine compiles both of these successfully,
			// which is exactly why the set is pinned explicitly.
			wasmBytes, err := wasmtime.Wat2Wasm(tc.wat)
			require.NoError(t, err)

			module, err := wasmtime.NewModule(runtime.engine, wasmBytes)
			if err == nil {
				module.Close()
				t.Fatalf("expected the pinned engine to reject a %s module, but it compiled", tc.name)
			}
		})
	}
}
