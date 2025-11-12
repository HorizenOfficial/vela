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

// TODO - When TinyGo compiles go into wasm, it configures the WebAssembly linear memory to an initial size
// of 2 pages (128KB), and marks a position in that memory as the heap base. All memory beyond that is used for the
// Go heap. Allocations within Go (compiled to %.wasm) are managed by calling memory.grow. The instruction is executed
// within the Wasm instance, but the actual memory growth happens in Wasmtime’s engine, which allocates more memory pages.
// A freeMap could be used to track memory blocks that have been freed and are available for reuse, avoiding the need
// to always request new memory pages from the WASM runtime, this is good for performances.
// The drawback is risk of mem fragmentation, and the fact that the application retains the maximum memory envelope
// it has ever needed.
// var freeList = make(map[int32][][]byte)

// total memory (in bytes) currently allocated
var cumulative_alloc_size int32

// --- WASM Memory Management Functions ---

//export allocate
func allocate(size int32) int32 {
	if size <= 0 {
		// Return a null pointer for zero or negative size.
		// The caller should check for zero-length data, but this makes the guest allocator more robust.
		return 0
	}
	// Allocate a byte slice of the desired size.
	data := make([]byte, size)

	// Get a pointer to the slice's underlying data.
	ptr := &data[0]

	// Convert to uintptr (valid on wasm32, which has 4GB address space max).
	uptr := uintptr(unsafe.Pointer(ptr))

	// Store a reference to the slice in our global map to "pin" it and preventing GC from acting
	allocatedMemory[uptr] = data

	cumulative_alloc_size += int32(size)

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

        // double deletion will decrease the counter, we choose to do it anyway for the time being, maybe we can
        // add more counters and stats in future
	cumulative_alloc_size -= int32(size)
}

// Note: if we call directly this function from the host, the ABI C interface foresees that
// for multiple return values, the values are stored into a pointer passed as the first parameter by the caller.
//
//export get_allocated_memory_stats
func GetAllocatedMemoryStats() (map_size, total_bytes int32) {
	map_size = int32(len(allocatedMemory))
	total_bytes = cumulative_alloc_size
	return
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
