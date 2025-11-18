package wasm

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/horizen-pes/pkg/logger"
	appCommon "github.com/horizen-pes/pkg/wasm/common"
	"github.com/stretchr/testify/require"
)

var testLogger logger.Logger

func TestMain(m *testing.M) {
	// Initialize once, by default it writes on stderr
	testLogger = logger.NewLogger(&logger.Config{Kind: "printf"})
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestWriteToMemory_NilModule(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	_, err := runtime.writeToMemory(nil, []byte("some data"))
	require.Error(t, err)
	require.Equal(t, "module is nil", err.Error())
}

func TestWriteToMemory_NilMemory(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	module := &ApplicationModule{} // memory is nil by default

	_, err := runtime.writeToMemory(module, []byte("some data"))
	require.Error(t, err)
	require.Equal(t, "memory not initialized", err.Error())
}

func TestWriteToMemory_NilData(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	// Create a valid memory object, but instance can be nil
	// since it's not reached when data is empty.
	memType := wasmtime.NewMemoryType(1, false, 0)
	memory, err := wasmtime.NewMemory(runtime.store, memType)
	require.NoError(t, err)

	module := &ApplicationModule{
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

func TestReadFromMemory_NilModule(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	_, err := runtime.readFromMemory(nil, 0, 0)
	require.Error(t, err)
	require.Equal(t, "module is nil", err.Error())
}

func TestReadFromMemory_NilMemory(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	module := &ApplicationModule{} // memory is nil by default

	_, err := runtime.readFromMemory(module, 0, 0)
	require.Error(t, err)
	require.Equal(t, "memory not initialized", err.Error())
}

func TestExtractResultBytes(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	// Setup a mock memory for testing, using the runtime's store
	// 1 is the minimum size in WebAssembly pages (WebAssembly page size = 64 KiB (65,536 bytes))
	memType := wasmtime.NewMemoryType(1, false, 0)
	memory, err := wasmtime.NewMemory(runtime.store, memType)
	require.NoError(t, err)

	appModule := &ApplicationModule{memory: memory}

	// Write some test data to memory: 4 bytes for length, then the data
	testData := []byte("hello world")
	dataWithLen := make([]byte, 4+len(testData))
	binary.LittleEndian.PutUint32(dataWithLen, uint32(len(testData)))
	copy(dataWithLen[4:], testData)
	// [0x0b,0x00,0x00,0x00, 0x68,..,0x64]

	// Copy the prepared data into the wasm memory
	memSlice := memory.UnsafeData(runtime.store)
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
		require.Contains(t, err.Error(), "wasm module returned a null pointer")
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := runtime.extractResultBytes("not a pointer", appModule)
		require.Error(t, err)
		require.Equal(t, "wasm module returned unexpected type", err.Error())
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

func TestClose(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger)
	// Add a dummy module to ensure the map is cleared
	runtime.modules["test-app"] = &ApplicationModule{}

	err := runtime.Close()
	require.NoError(t, err)

	require.Nil(t, runtime.store)
	require.Nil(t, runtime.engine)
	require.Empty(t, runtime.modules)
}
