package blockchain

import (
	"fmt"
)

type ReorgError struct {
	causedBy error
}

func (e ReorgError) Error() string {
	return  fmt.Sprintf("Possible reorg on chain, received error: %v", e.causedBy)
}