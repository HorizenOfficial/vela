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
//
//export allocate
func allocate(size int32) *byte {
	// Allocate memory and return pointer
	data := make([]byte, size)
	return &data[0]
}

// deallocate exports a function to deallocate memory.
//
//export deallocate
func deallocate(ptr *byte, size int32) {
	// no-op for deallocation in Go, as Go's garbage collector handles memory management
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
	destination := unsafe.Slice(ptr, n)
	binary.LittleEndian.PutUint32(destination[:4], uint32(dataLength))
	copy(destination[4:], data)

	return ptr
}

// SerializeAndWriteResult handles common serialization and returns a WASM pointer.
func SerializeAndWriteResult(result any) *byte {
	reportJSON, err := json.Marshal(result)
	if err != nil {
		return StringToPtr(WasmSerializationError)
	}
	return StringToPtr(reportJSON)
}
