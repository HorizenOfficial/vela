package utils

import (
	"encoding/binary"
	"encoding/json"
	"unsafe"
)

// WasmSerializationError is a generic error for failed WASM serialization.
var WasmSerializationError = []byte("{}")

// --- WASM Memory Management Functions ---

// allocate exports a function to allocate memory.

// TODO:
// The slice created here is a local variable.
// Once this function returns, there are no references to `data` inside the guest,
// so the Go GC is free to reclaim it at any time during later guest execution.
//
// That means the host may end up with a pointer into memory that could be freed
// or moved by the GC, causing nondeterministic crashes or corrupted data.
//
// To prevent this, production code should keep a global reference to the slice
// (e.g., store it in a global `map[int32][]byte` keyed by the returned pointer)
// until the host explicitly calls `deallocate()`. This pins the slice in memory
// and ensures the GC will not reclaim it prematurely

//export allocate
func allocate(size int32) int32 {
	// Allocate memory and return pointer address value
	data := make([]byte, size)
	return int32(uintptr(unsafe.Pointer(&data[0])))
}

// deallocate exports a function to deallocate memory.
// TODO:
// Since Go has a garbage collector, we normally don’t need manual free().
// However, if allocate() stores slices in a global map to keep them alive,
// then deallocate() must remove the corresponding entry so the GC can
// eventually reclaim the memory.
//
// In this no-op implementation, nothing is freed, so memory referenced by
// the global map (if implemented) would leak. In a correct implementation,
// deallocate() should delete the slice from the global map, dropping the
// reference and allowing the GC to reclaim it.
//
//export deallocate
func deallocate(ptr *byte, size int32) {
}

// --- Helper Functions for Data Translation ---

// PtrToString converts a WASM pointer and length to a Go string.
func PtrToString(ptr *byte, length int32) string {
	if ptr == nil || length == 0 {
		return ""
	}
	return string(unsafe.Slice(ptr, length))
}

// StringToPtr converts a Go byte slice to an allocated memory pointer for WASM.
func StringToPtr(data []byte) *byte {
	dataLength := len(data)
	if dataLength == 0 {
		return nil
	}

	n := 4 + dataLength // 4 bytes for length + actual data length
	ptr := allocate(int32(n))
	dataBytes := (*byte)(unsafe.Pointer(uintptr(ptr)))
	destination := unsafe.Slice(dataBytes, n)

	//destination := unsafe.Slice(ptr, n)
	binary.LittleEndian.PutUint32(destination[:4], uint32(dataLength))
	copy(destination[4:], data)

	return dataBytes
}

// SerializeAndWriteResult handles common serialization and returns a WASM pointer.
func SerializeAndWriteResult(result any) *byte {
	reportJSON, err := json.Marshal(result)
	if err != nil {
		return StringToPtr(WasmSerializationError)
	}
	return StringToPtr(reportJSON)
}
