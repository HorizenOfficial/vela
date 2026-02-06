package main

import (
	"github.com/horizen-pes/app/simple/app"
	"github.com/horizen-pes/app/simple/utils"
)

// --- WASM-Exposed Functions (Bridge to Application Logic) ---

// These functions handle the WASM I/O and call the high-level logic functions.

//export load_module
func load_module(appId int64) *byte {
	result := app.LoadModule(appId)
	return app.SerializeAndWriteResult(result)
}

//export deposit
func deposit(appId int64, senderPtr *byte, senderLen int32, valuePtr *byte, valueLen int32, statePtr *byte, stateLen int32) *byte {
	_ = appId
	sender := app.PtrToAddress(senderPtr, senderLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	value := app.PtrToUint256(valuePtr, valueLen)
	result := app.DepositFunds(sender, value, stateJSON)
	return app.SerializeAndWriteResult(result)
}

//export process_request
func process_request(appId int64, senderPtr *byte, senderLen int32, requestType int32, payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	_ = appId
	sender := app.PtrToAddress(senderPtr, senderLen)
	payloadJSON := utils.PtrToString(payloadPtr, payloadLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := app.ProcessRequest(sender, requestType, payloadJSON, stateJSON)
	return app.SerializeAndWriteResult(result)
}

//export get_memory_stats
func get_memory_stats() *byte {
	result := app.GetAllocatedMemoryStats()
	return app.SerializeAndWriteResult(result)
}

// Main function is required but not used in WASM
func main() {}
