package utils

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"unsafe"

	"github.com/horizen-pes/app/simple/app"
)

// Note: The TinyGo WASM module acts as a self-contained environment running on a single logical thread.
// Concurrency for calls to the WASM module must be handled on the host side (e.g., using Go's sync.Mutex
// in the calling Go code or better, using a single wasmtime.Instance for every call), as TinyGo's 'sync'
// package is not fully supported within the WASM environment for inter-module synchronization.

// allocatedMemory holds references to memory allocated by the guest, preventing the Go garbage
// collector from reclaiming it while it's in use by the host. The map key is the pointer address.
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
var cumulative_alloc_size int64

// --- WASM Memory Management Functions ---

//export allocate
func allocate(size int32) int32 {
	println("Calling allocate with size =", size)
	if size <= 0 {
		// Return a null pointer for zero or negative size.
		// The caller should check for zero-length data, but this makes the guest allocator more robust.
		return 0
	}
	// Allocate a byte slice of the desired size.
	data := make([]byte, size)
	return int32(uintptr(unsafe.Pointer(&data[0])))
}

//export deallocate
func deallocate(ptr *byte, size int32) {
	println("Calling deallocate on ptr =", ptr, ", with size =", size)

	// Get the uintptr from the pointer.
	uptr := uintptr(unsafe.Pointer(ptr))
	println("uptr =", uptr)

	println("Checking map")
	if allocatedMemory == nil {
		println("map is nil")
	}
	println("map=", allocatedMemory, ", addr=", &allocatedMemory)
	b, ok := allocatedMemory[uptr]
	if !ok {
		// we exit even if delete on a map would be a no op, but we also do not decrement the counter
		return
	}

	// Delete the reference from the map. This unpins the memory, making it eligible for garbage collection.
	delete(allocatedMemory, uptr)

	if int32(len(b)) != size {
		// do not update the counter, that would be incorrect anyway. We could add more counters for errors and stats in future
		return
	}

	cumulative_alloc_size -= int64(size)
}

// Note: if we call directly this function from the host, the ABI C interface foresees that
// for multiple return values, the values are stored into a pointer passed as the first parameter by the caller.
//
//export get_allocated_memory_stats
func GetAllocatedMemoryStats() (map_size, total_bytes int64) {
	map_size = int64(len(allocatedMemory))
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
		return StringToPtr([]byte(app.WasmSerializationError))
	}
	return StringToPtr(reportJSON)
}

// PtrToNonNegativeBigInt converts a WASM pointer and length representing the a big.Int value to a Go big.Int pointer.
// The byte slice is obtained with the (big.Int).Bytes() method, i.e. it represents the absolute value in big-endian byte order, so the value is always non-negative.
func PtrToNonNegativeBigInt(ptr *byte, length int32) *big.Int {
	if ptr == nil || length == 0 {
		return big.NewInt(0)
	}

	return new(big.Int).SetBytes(unsafe.Slice(ptr, length))
}

// PtrToAddress converts a WASM pointer and length to a ethereum address.
func PtrToAddress(ptr *byte, length int32) *app.Address {
	if ptr == nil || length == 0 {
		return nil
	}
	var address app.Address
	address.SetBytes(unsafe.Slice(ptr, length))
	return &address
}
