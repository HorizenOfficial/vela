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
	AppID    int64                    `json:"appId"`
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

// PayloadInstructions represents the deserialized payload instructions
type PayloadInstructions struct {
	Type            string               `json:"type"`
	CompareAccounts *CompareInstructions `json:"compare,omitempty"`
	Withdraw        *WithdrawInstruction `json:"withdraw,omitempty"`
}

// ReportPayloadInstructions represent a specific information on how to generate a report
// In this simple app its a custom tag to add to the report
type ReportPayloadInstructions struct {
	IncludeTag string `json:"tag,omitempty"`
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

// WithdrawalEvent is a local replacement for wasmCommon.WithdrawalEvent
type WithdrawalEvent SenderEvent
