package app

import (
	"github.com/HorizenOfficial/vela-common-go/wasm/types"
)

// ----- module internal types

// AccountState holds the per-token balances of a single user account.
type AccountState struct {
	Address  types.Address             `json:"address"`
	Balances map[string]*types.Uint256 `json:"balances"` // token address hex -> balance
}

// ApplicationInternalState is the persisted state of the trigger app.
//
// Counter is bumped by process_request (the "fire" action). TriggerAddress and
// FeeCollector come from the deploy params: "fire" withdraws to TriggerAddress
// (the registered trigger), and trusted_request credits the re-shielded amount to
// the FeeCollector account. Each step mutates state so it produces a distinct
// on-chain state root, which is what the round-trip test asserts.
type ApplicationInternalState struct {
	AppID          uint64                   `json:"appId"`
	Accounts       map[string]*AccountState `json:"accounts"`
	AllowedTokens  map[string]bool          `json:"allowedTokens"`  // token address hex -> allowed
	Counter        uint64                   `json:"counter"`
	TriggerAddress string                   `json:"triggerAddress"` // withdrawal destination (hex)
	FeeCollector   string                   `json:"feeCollector"`   // internal account credited on re-shield (hex)
}

// DeployParams are the constructor parameters supplied in the deploy descriptor.
type DeployParams struct {
	AllowedTokens  []string `json:"allowedTokens"`  // token address hex strings to allow
	TriggerAddress string   `json:"triggerAddress"` // trigger contract that receives "fire" withdrawals
	FeeCollector   string   `json:"feeCollector"`   // account credited with the re-shielded amount
}

// ProcessInstruction is the process_request payload. "fire" withdraws a
// counter-derived percentage of the sender's balance to the trigger and emits an
// AppEvent (which fires the registered trigger).
//
// Note: only the user-supplied process_request payload is JSON. The two
// trigger-facing wire formats are ABI-encoded (cheap for Solidity): the "fire"
// AppEvent.Data is abi.encode(uint256 counter) and the trusted_request payload is
// abi.encode(uint256 reshieldedAmount).
type ProcessInstruction struct {
	Type string `json:"type"`
}

// feeCreditedEvent is the JSON body of the encrypted PlainEvent emitted to the
// fee collector by trusted_request, reporting its new balance after a re-shield.
type feeCreditedEvent struct {
	Token   string         `json:"token"`
	Amount  *types.Uint256 `json:"amount"`  // amount credited this round
	Balance *types.Uint256 `json:"balance"` // fee collector's new balance
}
