package utils

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"unsafe"

	appCommon "github.com/horizen-pes/pkg/wasm/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

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
	return int32(uintptr(unsafe.Pointer(&data[0])))
}

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

// PtrToNonNegativeBigInt converts a WASM pointer and length representing the a big.Int value to a Go big.Int pointer.
// The byte slice is obtained with the (big.Int).Bytes() method, i.e. it represents the absolute value in big-endian byte order, so the value is always non-negative.
func PtrToNonNegativeBigInt(ptr *byte, length int32) *big.Int {
	if ptr == nil || length == 0 {
		return big.NewInt(0)
	}

	return new(big.Int).SetBytes(unsafe.Slice(ptr, length))
}


// PtrToAddress converts a WASM pointer and length to a ethereum address.
func PtrToAddress(ptr *byte, length int32) *ethCommon.Address {
	if ptr == nil || length == 0 {
		return nil
	}
	address := ethCommon.BytesToAddress(unsafe.Slice(ptr, length))
	return &address
}
