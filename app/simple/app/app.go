package app

import (
	"encoding/json"
	"fmt"

	"github.com/HorizenOfficial/vela-common-go/wasm/types"
	"github.com/HorizenOfficial/vela-common-go/wasm/utils"
	"github.com/HorizenOfficial/vela/app/subtype"
	"github.com/HorizenOfficial/vela/pkg/common"
)

// --- High-Level Application Logic ---

// ethTokenHex is the hex representation of the ETH token address (zero address).
var ethTokenHex = (types.Address{}).Hex()

// Deploy initializes the application state from constructor parameters.
// The AllowedTokens list specifies which tokens the app accepts; ETH (0x0) is always allowed.
func Deploy(appId int64, paramsJSON string) types.DeployResult {
	allowedTokens := make(map[string]bool)
	// ETH is always allowed
	allowedTokens[ethTokenHex] = true

	if paramsJSON != "" {
		var params DeployParams
		if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
			utils.LogError("Deploy: failed to parse deploy params: %v", err)
			return types.DeployResult{
				Error: fmt.Sprintf("failed to parse deploy params: %v", err),
			}
		}
		for _, tokenHex := range params.AllowedTokens {
			allowedTokens[tokenHex] = true
		}
	}

	initialState := &ApplicationInternalState{
		AppID:         uint64(appId),
		Accounts:      make(map[string]*AccountState),
		AllowedTokens: allowedTokens,
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		utils.LogError("Deploy: failed to marshal initial state: %v", err)
		return types.DeployResult{
			Error: fmt.Sprintf("failed to marshal initial state: %v", err),
		}
	}
	fuel := types.NewUint256(5)
	utils.LogDebug("Deploy: appId=%d, allowedTokens=%d, stateSize=%d, fuel=%v", uint64(appId), len(allowedTokens), len(stateJSON), fuel)
	return types.DeployResult{
		State: stateJSON,
		Fuel:  fuel,
	}
}

// LoadModule is retained for cache warm-up by getOrLoadModule (see wasmtime_runtime.go).
// New deployments should use Deploy instead.
func LoadModule(appId int64) types.LoadModuleResult {
	initialState := &ApplicationInternalState{
		AppID:         uint64(appId),
		Accounts:      make(map[string]*AccountState),
		AllowedTokens: map[string]bool{ethTokenHex: true},
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		utils.LogError("LoadModule: failed to marshal initial state: %v", err)
		return types.LoadModuleResult{
			Error: fmt.Sprintf("failed to marshal initial state: %v", err),
		}
	}
	fuel := types.NewUint256(5)
	utils.LogDebug("LoadModule: appId=%d, stateSize=%d, fuel=%v", uint64(appId), len(stateJSON), fuel)
	return types.LoadModuleResult{
		State: stateJSON,
		Fuel:  fuel,
	}
}

// getTokenBalance returns the balance for a specific token, or a zero value if not found.
func getTokenBalance(acc *AccountState, tokenHex string) *types.Uint256 {
	if bal, ok := acc.Balances[tokenHex]; ok {
		return bal
	}
	return types.NewUint256(0)
}

func DepositFunds(senderPtr *types.Address, tokenPtr *types.Address, value *types.Uint256, stateJSON string) types.DepositResult {
	if senderPtr == nil {
		utils.LogError("DepositFunds: sender address is nil")
		return types.DepositResult{Error: "Sender address is nil"}
	}

	if tokenPtr == nil {
		utils.LogError("DepositFunds: token address is nil")
		return types.DepositResult{Error: "Token address is nil"}
	}

	if value == nil {
		utils.LogError("DepositFunds: value is nil")
		return types.DepositResult{Error: "value is nil"}
	}

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		utils.LogError("DepositFunds: failed to parse application state: %v", err)
		return types.DepositResult{Error: fmt.Sprintf("Failed to parse application state: %v", err)}
	}

	tokenHex := tokenPtr.Hex()

	// Validate token against app allowlist
	if !currentState.AllowedTokens[tokenHex] {
		return types.DepositResult{Error: fmt.Sprintf("Token %s is not allowed by this application", tokenHex)}
	}

	senderHex := senderPtr.Hex()

	// Safe map access & initialization
	acc, exists := currentState.Accounts[senderHex]
	if !exists {
		acc = &AccountState{
			Address:  *senderPtr,
			Balances: make(map[string]*types.Uint256),
		}
		currentState.Accounts[senderHex] = acc
	}
	if acc.Balances == nil {
		acc.Balances = make(map[string]*types.Uint256)
	}

	// Get or initialize per-token balance
	balance, exists := acc.Balances[tokenHex]
	if !exists {
		balance = types.NewUint256(0)
		acc.Balances[tokenHex] = balance
	}

	// Update balance in-place
	if balance.AddOverflow(*balance, *value) {
		return types.DepositResult{Error: fmt.Sprintf("Overflow while adding amount %s to balance: %s", value, balance)}
	}

	// Create deposit event
	eventData := DepositEvent{
		Type:         "deposit",
		TokenAddress: *tokenPtr,
		Amount:       value,
		Balance:      balance,
	}

	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		utils.LogError("DepositFunds: failed to serialize event data: %v", err)
		return types.DepositResult{Error: fmt.Sprintf("Failed to serialize event data: %+v, err: %v", eventData, err)}
	}

	events := []types.PlainEvent{{
		UserID:       *senderPtr,
		EventSubType: subtype.FromString("deposit"),
		Data:         eventDataBytes,
	}}

	// App-level event: visible to everyone (not encrypted)
	appEventData := DepositAppEvent{
		TokenAddress: *tokenPtr,
		Amount:       value,
	}
	appEventDataBytes, err := json.Marshal(appEventData)
	if err != nil {
		utils.LogError("DepositFunds: failed to serialize app event data: %v", err)
		return types.DepositResult{Error: fmt.Sprintf("Failed to serialize app event data: %+v, err: %v", appEventData, err)}
	}
	appEvents := []types.AppEvent{{
		EventSubType: subtype.FromString("deposit_received"),
		Data:         appEventDataBytes,
	}}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		utils.LogError("DepositFunds: failed to serialize new state: %v", err)
		return types.DepositResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
	}
	fuel := types.NewUint256(35)
	utils.LogDebug("DepositFunds: sender=%s, token=%s, value=%v, newBalance=%v, eventsCount=%d, appEventsCount=%d, stateSize=%d, fuel=%v",
		senderHex, tokenHex, value, balance, len(events), len(appEvents), len(newStateBytes), fuel)
	return types.DepositResult{State: newStateBytes, Events: events, AppEvents: appEvents, Fuel: fuel}
}

func ProcessRequest(senderPtr *types.Address, requestType int32, payloadJSON, stateJSON string) types.ProcessResult {
	// Deserialize current state
	if senderPtr == nil {
		utils.LogError("ProcessRequest: sender address is nil")
		return types.ProcessResult{Error: "Sender address is nil"}
	}

	sender := *senderPtr
	senderHex := sender.Hex()

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		utils.LogError("ProcessRequest: failed to parse application state: %v", err)
		return types.ProcessResult{Error: fmt.Sprintf("Failed to parse application state: %v", err)}
	}

	var events []types.PlainEvent
	var withdrawals []types.Withdrawal

	// Determine instruction type: requestType has priority over payload
	var instructionType string
	var instructions PayloadInstructions

	// Parse payload if present (for additional options like deanonymize tag)
	if payloadJSON != "" && payloadJSON != "{}" {
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			utils.LogError("ProcessRequest: failed to parse payload instructions: %v", err)
			return types.ProcessResult{Error: fmt.Sprintf("Failed to parse payload instructions: %v", err)}
		}
	}

	// requestType takes precedence over payload type
	if requestType == int32(common.Deanonymize) {
		instructionType = "deanonymize"
	} else {
		instructionType = instructions.Type
	}

	if instructionType != "" {
		switch instructionType {
		case "compare_addresses":
			if instructions.CompareAccounts == nil {
				utils.LogError("ProcessRequest: compare instruction is missing in payload")
				return types.ProcessResult{Error: fmt.Sprintf("Compare instruction is missing in payload: %s", payloadJSON)}
			}
			targetAddress := instructions.CompareAccounts.TargetAddress
			targetAddressHex := targetAddress.Hex()

			// Determine which token to compare (defaults to ETH)
			tokenHex := ethTokenHex
			if instructions.CompareAccounts.TokenAddress != (types.Address{}) {
				tokenHex = instructions.CompareAccounts.TokenAddress.Hex()
			}

			// Validate token against app allowlist
			if !currentState.AllowedTokens[tokenHex] {
				return types.ProcessResult{Error: fmt.Sprintf("Token %s is not allowed for comparison", tokenHex)}
			}

			// Validate accounts to be compared
			if currentState.Accounts[senderHex] == nil {
				utils.LogError("ProcessRequest: account %s does not exist", senderHex)
				return types.ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", senderHex)}
			}
			if currentState.Accounts[targetAddressHex] == nil {
				utils.LogError("ProcessRequest: account %s does not exist", targetAddressHex)
				return types.ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", targetAddressHex)}
			}

			targetBalance := getTokenBalance(currentState.Accounts[targetAddressHex], tokenHex)
			senderBalance := getTokenBalance(currentState.Accounts[senderHex], tokenHex)

			var cmp = ""
			switch targetBalance.Cmp(*senderBalance) {
			case -1:
				cmp = "richer than"
			case 1:
				cmp = "poorer than"
			case 0:
				cmp = "as wealthy as"

			}

			sentence := sender.Hex() + " is " + cmp + " " + targetAddress.Hex() + " (token: " + tokenHex + ")"

			// Create action event
			eventData := map[string]interface{}{
				"type":     "compare_accounts",
				"sentence": sentence,
			}
			eventDataBytes, err := json.Marshal(eventData)
			if err != nil {
				utils.LogError("ProcessRequest: failed to serialize event data: %v", err)
				return types.ProcessResult{Error: fmt.Sprintf("Failed to serialize event data: %+v, err: %v", eventData, err)}
			}

			events = append(events, types.PlainEvent{
				UserID:       sender,
				EventSubType: subtype.FromString("compare_accounts"),
				Data:         eventDataBytes,
			})

		case "withdraw":
			if instructions.Withdraw == nil {
				utils.LogError("ProcessRequest: withdraw instruction is missing in payload")
				return types.ProcessResult{Error: fmt.Sprintf("Withdraw instruction is missing in payload: %s", payloadJSON)}
			}
			if instructions.Withdraw.Amount == nil {
				return types.ProcessResult{Error: "Withdraw amount is nil"}
			}

			tokenHex := instructions.Withdraw.TokenAddress.Hex()

			// Validate token against app allowlist
			if !currentState.AllowedTokens[tokenHex] {
				return types.ProcessResult{Error: fmt.Sprintf("Token %s is not allowed for withdrawal", tokenHex)}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[senderHex] == nil {
				utils.LogError("ProcessRequest: account %s does not exist", senderHex)
				return types.ProcessResult{Error: fmt.Sprintf("Account %s does not exist", senderHex)}
			}

			senderBalance := getTokenBalance(currentState.Accounts[senderHex], tokenHex)

			if senderBalance.Cmp(*instructions.Withdraw.Amount) < 0 {
				utils.LogError("ProcessRequest: insufficient balance for account %s", senderHex)
				return types.ProcessResult{Error: fmt.Sprintf("Insufficient balance %s for withdrawal %s for account %s",
					senderBalance, *instructions.Withdraw.Amount, senderHex)}
			}

			// Execute withdrawal — debit per-token balance
			senderBalance.Sub(*senderBalance, *instructions.Withdraw.Amount)

			// Create token-aware withdrawal
			withdrawals = append(withdrawals, types.Withdrawal{
				TokenAddress:       instructions.Withdraw.TokenAddress,
				DestinationAddress: instructions.Withdraw.To,
				Amount:             instructions.Withdraw.Amount,
			})

			// Create event for sender
			withdrawEventData := WithdrawalEvent{
				Type:         "withdrawal",
				To:           instructions.Withdraw.To,
				TokenAddress: instructions.Withdraw.TokenAddress,
				Amount:       instructions.Withdraw.Amount,
				Balance:      senderBalance,
			}
			withdrawEventDataBytes, err := json.Marshal(withdrawEventData)
			if err != nil {
				utils.LogError("ProcessRequest: failed to serialize withdraw event data: %v", err)
				return types.ProcessResult{Error: fmt.Sprintf("Failed to serialize withdraw event data: %+v, err: %v", withdrawEventData, err)}
			}

			events = append(events, types.PlainEvent{
				UserID:       sender,
				EventSubType: subtype.FromString("withdrawal"),
				Data:         withdrawEventDataBytes,
			})

		case "deanonymize":
			// Generate deanonymization report
			report := DeanonymizationReport{
				Accounts:      currentState.Accounts,
				AllowedTokens: currentState.AllowedTokens,
			}

			// Add optional tag from payload
			if instructions.Deanonymize != nil && instructions.Deanonymize.IncludeTag != "" {
				report.Tag = instructions.Deanonymize.IncludeTag
			}

			// Serialize the report
			reportBytes, err := json.Marshal(report)
			if err != nil {
				return types.ProcessResult{Error: fmt.Sprintf("Failed to serialize deanonymization report: %v", err)}
			}

			return types.ProcessResult{
				State:  []byte(stateJSON), //we have not modified the state, so returning the old one
				Report: reportBytes,
				Fuel:   types.NewUint256(20),
			}

		default:
			utils.LogError("ProcessRequest: unsupported instruction type: %s", instructions.Type)
			return types.ProcessResult{Error: fmt.Sprintf("Unsupported instruction type: [%s]", instructionType)}
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		utils.LogError("ProcessRequest: failed to serialize new state: %v", err)
		return types.ProcessResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
	}
	fuel := types.NewUint256(50)
	utils.LogDebug("ProcessRequest: sender=%s, eventsCount=%d, withdrawalsCount=%d, stateSize=%d, fuel=%v",
		senderHex, len(events), len(withdrawals), len(newStateBytes), fuel)
	return types.ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
		Fuel:        fuel,
	}
}

func GetAllocatedMemoryStats() types.MemoryStats {
	map_size, total_bytes := utils.GetAllocatedMemoryStats()
	return types.MemoryStats{
		MapSize:              map_size,
		CumulativeMemorySize: total_bytes,
	}
}
