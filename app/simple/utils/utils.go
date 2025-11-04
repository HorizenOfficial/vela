package utils

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"unsafe"

	appCommon "github.com/horizen-pes/pkg/wasm/common"
)

// --- WASM Memory Management Functions ---

//export allocate
func allocate(size int32) int32 {
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

// PtrToDepositValue converts a WASM pointer and length representing the value passed to a Deposit request to a Go big.Int pointer.
// The value is expected to be non-negative
func PtrToDepositValue(ptr *byte, length int32) *big.Int {
	if ptr == nil || length == 0 {
		return big.NewInt(0)
	}

	return new(big.Int).SetBytes(unsafe.Slice(ptr, length))
}