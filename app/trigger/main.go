package main

import (
	"github.com/HorizenOfficial/vela-common-go/wasm/types"
	"github.com/HorizenOfficial/vela-common-go/wasm/utils"
	"github.com/HorizenOfficial/vela/app/trigger/app"
)

// --- WASM-Exposed Functions (Bridge to Application Logic) ---

//export deploy
func deploy(appId int64, paramsPtr *byte, paramsLen int32) *byte {
	paramsJSON := utils.PtrToString(paramsPtr, paramsLen)
	result := app.Deploy(appId, paramsJSON)
	return types.SerializeAndWriteResult(result)
}

//export deposit
func deposit(appId int64, senderPtr *byte, senderLen int32, tokenPtr *byte, tokenLen int32, valuePtr *byte, valueLen int32, statePtr *byte, stateLen int32) *byte {
	_ = appId
	sender := types.PtrToAddress(senderPtr, senderLen)
	token := types.PtrToAddress(tokenPtr, tokenLen)
	value := types.PtrToUint256(valuePtr, valueLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := app.DepositFunds(sender, token, value, stateJSON)
	return types.SerializeAndWriteResult(result)
}

//export process_request
func process_request(appId int64, senderPtr *byte, senderLen int32, requestType int32, payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	_ = appId
	sender := types.PtrToAddress(senderPtr, senderLen)
	payloadJSON := utils.PtrToString(payloadPtr, payloadLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := app.ProcessRequest(sender, requestType, payloadJSON, stateJSON)
	return types.SerializeAndWriteResult(result)
}

//export trusted_request
func trusted_request(appId int64, payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	_ = appId
	payload := utils.PtrToString(payloadPtr, payloadLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := app.TrustedRequest(payload, stateJSON)
	return types.SerializeAndWriteResult(result)
}

//export get_memory_stats
func get_memory_stats() *byte {
	result := app.GetAllocatedMemoryStats()
	return types.SerializeAndWriteResult(result)
}

// Main function is required but not used in WASM
func main() {}
