// Package common provide the data structures used in the application.
package common

import (
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

func FnName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

var maxUint256 *big.Int

func init() {
	// maxUint256 = 2^256 - 1
	maxUint256 = new(big.Int).Sub(
		new(big.Int).Lsh(big.NewInt(1), 256),
		big.NewInt(1),
	)
}

func validateBigInt(name string, v *big.Int, allowZero bool) error {
	if v == nil {
		return nil // optional field, nil is allowed
	}

	if v.Sign() < 0 {
		if allowZero {
			return fmt.Errorf("%s must be >= 0", name)
		}
		return fmt.Errorf("%s must be > 0", name)
	}

	if !allowZero && v.Sign() == 0 {
		return fmt.Errorf("%s must be > 0", name)
	}

	if v.Cmp(maxUint256) > 0 {
		return fmt.Errorf("%s must be <= uint256", name)
	}

	return nil
}
