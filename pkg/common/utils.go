// Package common provide the data structures used in the application.
package common

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"runtime"
)

func (rt RequestType) ToUint8() (uint8, error) {
	switch rt {
	case Deploy:
		return 0, nil
	case Process:
		return 1, nil
	case Deanonymize:
		return 2, nil
	case AssociateKey:
		return 3, nil
	default:
		return 0, fmt.Errorf("unknown RequestType: %s", rt)
	}
}

func UInt8ToRequestResultStatus(i uint8) (RequestResultStatus, error) {
	switch i {
	case 0:
		return RequestResultOK, nil
	case 1:
		return RequestResultFailed, nil
	case 2:
		return RequestResultFailedNotRefunded, nil
	default:
		return RequestResultUnknown, fmt.Errorf("unknown request status value %d", i)
	}
}

func StringToBigInt(s string) (*big.Int, bool) {
	return new(big.Int).SetString(s, 10)
}

func RequestIdStringTo32Byte(s string) ([32]byte, error) {

	arr, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, fmt.Errorf("requestId string is not a valid hex string: %w", err)
	}
	if len(arr) > 32 {
		return [32]byte{}, fmt.Errorf("requestId string must not be more than 32 bytes long, got %d", len(arr))
	}

	var arr32 [32]byte
	copy(arr32[:], arr)
	return arr32, nil
}

func RequestId32ByteToString(b [32]byte) string {
	return hex.EncodeToString(b[:])
}

func FnName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
