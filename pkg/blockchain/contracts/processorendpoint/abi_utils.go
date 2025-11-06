package processorendpoint

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
)

func (c *ProcessorEndpoint) GetEventID(eventName string) ethCommon.Hash {
	return c.abi.Events[eventName].ID
}

// This method converts the ApplicationIdType to the Solidity binding type (uint64)
func ApplicationIdToBindingType(aid common.ApplicationIdType) uint64 {
	return uint64(aid)
}

func ApplicationIdFromBindingType(id uint64) common.ApplicationIdType {
	return common.ApplicationIdType(id)
}
