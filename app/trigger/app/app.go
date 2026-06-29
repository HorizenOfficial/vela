package app

import (
	"encoding/json"
	"fmt"

	"github.com/HorizenOfficial/vela-common-go/wasm/types"
	"github.com/HorizenOfficial/vela-common-go/wasm/utils"
	"github.com/HorizenOfficial/vela/app/subtype"
)

// --- High-Level Application Logic ---
//
// This is the minimal-but-real guest app used by the trigger e2e tests. It
// exercises the full unshield/re-shield round-trip with ETH: "fire" withdraws a
// share of the sender's balance to the registered trigger, and trusted_request
// credits the re-shielded remainder to a fee-collector account.

// ethTokenHex is the hex representation of the ETH token address (zero address).
var ethTokenHex = (types.Address{}).Hex()

// newInitialState builds the deploy-time state with ETH always allowed.
func newInitialState(appId int64, allowedTokens map[string]bool, triggerAddr, feeCollector string) *ApplicationInternalState {
	if allowedTokens == nil {
		allowedTokens = make(map[string]bool)
	}
	allowedTokens[ethTokenHex] = true
	return &ApplicationInternalState{
		AppID:          uint64(appId),
		Accounts:       make(map[string]*AccountState),
		AllowedTokens:  allowedTokens,
		Counter:        0,
		TriggerAddress: triggerAddr,
		FeeCollector:   feeCollector,
	}
}

// Deploy initializes the application state from constructor parameters.
func Deploy(appId int64, paramsJSON string) types.DeployResult {
	allowedTokens := make(map[string]bool)
	var triggerAddr, feeCollector string
	if paramsJSON != "" {
		var params DeployParams
		if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
			utils.LogError("Deploy: failed to parse deploy params: %v", err)
			return types.DeployResult{Error: fmt.Sprintf("failed to parse deploy params: %v", err)}
		}
		for _, tokenHex := range params.AllowedTokens {
			allowedTokens[tokenHex] = true
		}
		triggerAddr = params.TriggerAddress
		feeCollector = params.FeeCollector
	}

	stateJSON, err := json.Marshal(newInitialState(appId, allowedTokens, triggerAddr, feeCollector))
	if err != nil {
		utils.LogError("Deploy: failed to marshal initial state: %v", err)
		return types.DeployResult{Error: fmt.Sprintf("failed to marshal initial state: %v", err)}
	}
	fuel := types.NewUint256(5)
	utils.LogDebug("Deploy: appId=%d, allowedTokens=%d, stateSize=%d", uint64(appId), len(allowedTokens), len(stateJSON))
	return types.DeployResult{State: stateJSON, Fuel: fuel}
}

// LoadModule is retained for cache warm-up by getOrLoadModule (see wasmtime_runtime.go).
func LoadModule(appId int64) types.LoadModuleResult {
	stateJSON, err := json.Marshal(newInitialState(appId, nil, "", ""))
	if err != nil {
		utils.LogError("LoadModule: failed to marshal initial state: %v", err)
		return types.LoadModuleResult{Error: fmt.Sprintf("failed to marshal initial state: %v", err)}
	}
	return types.LoadModuleResult{State: stateJSON, Fuel: types.NewUint256(5)}
}

// DepositFunds credits the sender's per-token balance and emits a deposit event.
func DepositFunds(senderPtr *types.Address, tokenPtr *types.Address, value *types.Uint256, stateJSON string) types.DepositResult {
	if senderPtr == nil || tokenPtr == nil || value == nil {
		return types.DepositResult{Error: "sender, token and value must all be non-nil"}
	}

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return types.DepositResult{Error: fmt.Sprintf("Failed to parse application state: %v", err)}
	}

	tokenHex := tokenPtr.Hex()
	if !currentState.AllowedTokens[tokenHex] {
		return types.DepositResult{Error: fmt.Sprintf("Token %s is not allowed by this application", tokenHex)}
	}

	senderHex := senderPtr.Hex()
	acc, exists := currentState.Accounts[senderHex]
	if !exists {
		acc = &AccountState{Address: *senderPtr, Balances: make(map[string]*types.Uint256)}
		currentState.Accounts[senderHex] = acc
	}
	if acc.Balances == nil {
		acc.Balances = make(map[string]*types.Uint256)
	}
	balance, exists := acc.Balances[tokenHex]
	if !exists {
		balance = types.NewUint256(0)
		acc.Balances[tokenHex] = balance
	}
	if balance.AddOverflow(*balance, *value) {
		return types.DepositResult{Error: fmt.Sprintf("Overflow while adding amount %s to balance %s", value, balance)}
	}

	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return types.DepositResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
	}
	return types.DepositResult{State: newStateBytes, Fuel: types.NewUint256(35)}
}

// ProcessRequest handles a regular PROCESS request. The "fire" action bumps the
// Counter, withdraws Counter% of the sender's ETH balance to the registered
// trigger contract (unshield), and emits an application-level (non-encrypted)
// AppEvent — the AppEvent is what causes the trigger to fire during stateUpdate.
func ProcessRequest(senderPtr *types.Address, requestType int32, payloadJSON, stateJSON string) types.ProcessResult {
	_ = requestType

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return types.ProcessResult{Error: fmt.Sprintf("Failed to parse application state: %v", err)}
	}

	var instr ProcessInstruction
	if payloadJSON != "" && payloadJSON != "{}" {
		if err := json.Unmarshal([]byte(payloadJSON), &instr); err != nil {
			return types.ProcessResult{Error: fmt.Sprintf("Failed to parse payload instructions: %v", err)}
		}
	}

	var appEvents []types.AppEvent
	var withdrawals []types.Withdrawal
	switch instr.Type {
	case "fire":
		if senderPtr == nil {
			return types.ProcessResult{Error: "fire requires a sender"}
		}
		currentState.Counter++

		// Withdraw Counter% of the sender's ETH balance to the trigger (unshield).
		acc := currentState.Accounts[senderPtr.Hex()]
		if acc == nil || acc.Balances[ethTokenHex] == nil || acc.Balances[ethTokenHex].IsZero() {
			return types.ProcessResult{Error: "fire: sender has no ETH balance to withdraw"}
		}
		bal := acc.Balances[ethTokenHex]
		// Test balances fit in 64 bits, so compute the percentage in uint64. 1 = 1%.
		withdrawAmt := types.NewUint256(bal[0] * currentState.Counter / 100)
		if withdrawAmt.IsZero() {
			return types.ProcessResult{Error: "fire: withdrawal rounds to zero"}
		}
		if bal.SubOverflow(*bal, *withdrawAmt) {
			return types.ProcessResult{Error: "fire: withdrawal exceeds balance"}
		}

		triggerAddr, err := types.HexToAddress(currentState.TriggerAddress)
		if err != nil {
			return types.ProcessResult{Error: fmt.Sprintf("fire: invalid trigger address %q: %v", currentState.TriggerAddress, err)}
		}
		withdrawals = []types.Withdrawal{{
			TokenAddress:       types.Address{}, // ETH
			DestinationAddress: triggerAddr,
			Amount:             withdrawAmt,
		}}

		// AppEvent.Data is ABI-encoded (abi.encode(uint256 counter)), not JSON: a
		// trigger contract decodes this on-chain and ABI decoding is far cheaper
		// than JSON parsing in Solidity. Uint256.Bytes() is the 32-byte big-endian
		// form, which is exactly the ABI encoding of a single uint256.
		appEvents = []types.AppEvent{{
			EventSubType: subtype.FromString("execute_requested"),
			Data:         types.NewUint256(currentState.Counter).Bytes(),
		}}
	case "ping":
		// Like "fire" but with no withdrawal: bumps Counter (so the state root
		// changes) and emits an AppEvent to fire the trigger, without needing the
		// sender to hold any balance. Used to drive a TRUSTPROCESS without funds.
		currentState.Counter++
		appEvents = []types.AppEvent{{
			EventSubType: subtype.FromString("execute_requested"),
			Data:         types.NewUint256(currentState.Counter).Bytes(),
		}}
	case "":
		// no-op request, state unchanged
	default:
		return types.ProcessResult{Error: fmt.Sprintf("Unsupported instruction type: [%s]", instr.Type)}
	}

	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return types.ProcessResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
	}
	return types.ProcessResult{State: newStateBytes, AppEvents: appEvents, Withdrawals: withdrawals, Fuel: types.NewUint256(50)}
}

// TrustedRequest handles a TRUSTPROCESS request dispatched to the trusted_request
// export. The payload is produced by the trigger contract and arrives in clear
// text (no sender, no request type). It is ABI-encoded — abi.encode(uint256
// reshieldedAmount) — not JSON, because a trigger contract builds it on-chain where
// ABI encoding is much cheaper than JSON. TrustedRequest credits the re-shielded
// amount to the fee-collector account, emits an encrypted event to that account
// with its new balance, and deliberately emits NO AppEvents: a TRUSTPROCESS
// stateUpdate with empty AppEvents makes a guarded trigger return no payload,
// terminating the loop.
func TrustedRequest(payloadString, stateJSON string) types.ProcessResult {
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return types.ProcessResult{Error: fmt.Sprintf("Failed to parse application state: %v", err)}
	}

	// ABI-decode the trusted payload (abi.encode(uint256 reshieldedAmount)): a
	// single 32-byte big-endian word. Reject anything else — a malformed payload
	// is a deterministic failure, so the executor marks the TRUSTPROCESS FAILED.
	payload := []byte(payloadString)
	if len(payload) != 32 {
		return types.ProcessResult{Error: fmt.Sprintf("trusted_request: malformed payload, expected 32 bytes, got %d", len(payload))}
	}
	var amount types.Uint256
	amount.SetBytes(payload)

	feeAddr, err := types.HexToAddress(currentState.FeeCollector)
	if err != nil {
		return types.ProcessResult{Error: fmt.Sprintf("invalid fee collector address %q: %v", currentState.FeeCollector, err)}
	}

	feeHex := currentState.FeeCollector
	acc := currentState.Accounts[feeHex] 
	if acc == nil {
		acc = &AccountState{Address: feeAddr, Balances: make(map[string]*types.Uint256)}
		currentState.Accounts[feeHex] = acc
	}
	if acc.Balances == nil {
		acc.Balances = make(map[string]*types.Uint256)
	}
	bal := acc.Balances[ethTokenHex]
	if bal == nil {
		bal = types.NewUint256(0)
		acc.Balances[ethTokenHex] = bal
	}
	if bal.AddOverflow(*bal, amount) {
		return types.ProcessResult{Error: "re-shield credit overflows fee collector balance"}
	}

	// Encrypted user event reporting the fee collector's new balance. The host
	// encrypts the Data to the fee collector's registered key (UserID), so only it
	// (and the test, which holds the key) can read it.
	amountCopy, balCopy := amount, *bal
	eventData, err := json.Marshal(feeCreditedEvent{Token: ethTokenHex, Amount: &amountCopy, Balance: &balCopy})
	if err != nil {
		return types.ProcessResult{Error: fmt.Sprintf("Failed to serialize fee event: %v", err)}
	}
	events := []types.PlainEvent{{
		UserID:       feeAddr,
		EventSubType: subtype.FromString("fee_credited"),
		Data:         eventData,
	}}

	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return types.ProcessResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
	}
	return types.ProcessResult{State: newStateBytes, Events: events, Fuel: types.NewUint256(50)}
}

func GetAllocatedMemoryStats() types.MemoryStats {
	mapSize, totalBytes := utils.GetAllocatedMemoryStats()
	return types.MemoryStats{MapSize: mapSize, CumulativeMemorySize: totalBytes}
}
