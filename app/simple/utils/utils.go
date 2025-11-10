package utils

import (
	"encoding/binary"
	"encoding/json"
	"unsafe"

	appCommon "github.com/horizen-pes/pkg/wasm/common"
)

// allocatedMemory holds references to memory allocated by the guest, preventing the Go garbage
// collector from reclaiming it while it's in use by the host. The map key is the pointer address.
// Note: The TinyGo WASM module acts as a self-contained environment running on a single logical thread.
// Even if more goroutines could be running in a module, only one can be executing its code at any given time
// therefore no mutex is necessary to protect this map.
var allocatedMemory = make(map[uintptr][]byte)

// --- WASM Memory Management Functions ---

//export allocate
func allocate(size int32) int32 {
	// Allocate a byte slice of the desired size.
	data := make([]byte, size)

	// Get a pointer to the slice's underlying data.
	ptr := &data[0]

	// Convert to uintptr (valid on wasm32, which has 4GB address space max).
	uptr := uintptr(unsafe.Pointer(ptr))

	// Store a reference to the slice in our global map to "pin" it and preventing GC from acting
	allocatedMemory[uptr] = data

	// Return the pointer address as an int32 to the host.
	return int32(uptr)
}

//export deallocate
func deallocate(ptr *byte, size int32) {

	// Get the uintptr from the pointer.
	uptr := uintptr(unsafe.Pointer(ptr))

	// Delete the reference from the map. This unpins the memory, making it eligible for garbage collection.
	// (Deleting a non-existent key is a no-op)
	delete(allocatedMemory, uptr)
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
	ptrVal := allocate(int32(n))
	dataBytes := (*byte)(unsafe.Pointer(uintptr(ptrVal)))
	destination := unsafe.Slice(dataBytes, n)

	binary.LittleEndian.PutUint32(destination[:4], uint32(dataLength))
	copy(destination[4:], data)

	return dataBytes
}

// SerializeAndWriteResult handles common serialization and returns a WASM pointer.
func SerializeAndWriteResult(result any) *byte {
	reportJSON, err := json.Marshal(result)
	if err != nil {
		return StringToPtr([]byte(appCommon.WasmSerializationError))
	}
	return StringToPtr(reportJSON)
}
