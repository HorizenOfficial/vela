package main

import (
	"github.com/horizen-pes/app/simple/app"
	"github.com/horizen-pes/app/simple/utils"
)

// --- WASM-Exposed Functions (Bridge to Application Logic) ---

// These functions handle the WASM I/O and call the high-level logic functions.

//export load_module
func load_module(appIdPtr *byte, appIdLen int32) *byte {
	appId := utils.PtrToString(appIdPtr, appIdLen)
	stateBytes := app.LoadModule(appId)
	return utils.StringToPtr(stateBytes)
}

//export deposit
func deposit(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, valuePtr *byte, valueLen int32, statePtr *byte, stateLen int32) *byte {
	_ = utils.PtrToString(appIdPtr, appIdLen)
	sender := utils.PtrToString(senderPtr, senderLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	value := utils.PtrToDepositValue(valuePtr, valueLen)
	result := app.DepositFunds(sender, value, stateJSON)
	return utils.SerializeAndWriteResult(result)
}

//export process_request
func process_request(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	_ = utils.PtrToString(appIdPtr, appIdLen)
	sender := utils.PtrToString(senderPtr, senderLen)
	payloadJSON := utils.PtrToString(payloadPtr, payloadLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := app.ProcessRequest(sender, payloadJSON, stateJSON)
	return utils.SerializeAndWriteResult(result)
}

//export generate_deanonymization_report
func generate_deanonymization_report(payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	payloadJSON := utils.PtrToString(payloadPtr, payloadLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := app.GenerateDeanonymizationReport(payloadJSON, stateJSON)
	return utils.SerializeAndWriteResult(result)
}

// Main function is required but not used in WASM
func main() {}
