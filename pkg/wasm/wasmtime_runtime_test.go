package wasm

import (
	"encoding/binary"
	"math"
	"os"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/logger"
	appCommon "github.com/HorizenOfficial/vela/pkg/wasm/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	memType := wasmtime.NewMemoryType(1, false, 0)
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
	memType := wasmtime.NewMemoryType(1, false, 0)
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
	memType := wasmtime.NewMemoryType(1, false, 0)
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
