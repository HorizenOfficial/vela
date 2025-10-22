package processorendpoint
import (
	 "github.com/ethereum/go-ethereum/common"
)

func (c *ProcessorEndpoint) GetEventID(eventName string) common.Hash {
	return c.abi.Events[eventName].ID
}