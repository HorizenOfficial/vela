package app

import (
	"github.com/horizen-cce-common-go/wasm/types"
)

// ----- module internal types

// AccountState represents the state of a user account
type AccountState struct {
	Address types.Address  `json:"address"`
	Balance *types.Uint256 `json:"balance"`
}

// ApplicationInternalState represents the internal state of the application
type ApplicationInternalState struct {
	AppID    uint64                   `json:"appId"`
	Accounts map[string]*AccountState `json:"accounts"`
}

// WithdrawInstruction represents instructions for withdrawing funds
type WithdrawInstruction struct {
	To     types.Address  `json:"to"`
	Amount *types.Uint256 `json:"amount"`
}

type CompareInstructions struct {
	TargetAddress types.Address `json:"targetAddress"`
}

// DeanonymizeInstruction represents instructions for deanonymization
type DeanonymizeInstruction struct {
	IncludeTag string `json:"tag,omitempty"`
}

// PayloadInstructions represents the deserialized payload instructions
type PayloadInstructions struct {
	Type            string                  `json:"type"`
	CompareAccounts *CompareInstructions    `json:"compare,omitempty"`
	Withdraw        *WithdrawInstruction    `json:"withdraw,omitempty"`
	Deanonymize     *DeanonymizeInstruction `json:"deanonymize,omitempty"`
}

// DeanonymizationReport represents the structure of the deanonymization report.
type DeanonymizationReport struct {
	Tag      string                   `json:"tag,omitempty"`
	Accounts map[string]*AccountState `json:"accounts"`
}

type DepositEvent struct {
	Type    string         `json:"type"`
	Amount  *types.Uint256 `json:"amount"`
	Balance *types.Uint256 `json:"balance"`
	Nonce   uint64         `json:"nonce"`
}

type SenderEvent struct {
	Type    string         `json:"type"`
	To      types.Address  `json:"to"`
	Amount  *types.Uint256 `json:"amount"`
	Balance *types.Uint256 `json:"balance"`
	Nonce   uint64         `json:"nonce"`
}

type RecipientEvent struct {
	Type    string         `json:"type"`
	From    types.Address  `json:"from"`
	Amount  *types.Uint256 `json:"amount"`
	Balance *types.Uint256 `json:"balance"`
	Nonce   uint64         `json:"nonce"`
}

// WithdrawalEvent has the same structure as SenderEvent
type WithdrawalEvent SenderEvent
