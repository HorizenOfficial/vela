package app

import (
	"github.com/HorizenOfficial/vela-common-go/wasm/types"
)

// ----- module internal types

// AccountState represents the state of a user account with per-token balances
type AccountState struct {
	Address  types.Address                 `json:"address"`
	Balances map[string]*types.Uint256     `json:"balances"` // token address hex -> balance
}

// ApplicationInternalState represents the internal state of the application
type ApplicationInternalState struct {
	AppID         uint64                   `json:"appId"`
	Accounts      map[string]*AccountState `json:"accounts"`
	AllowedTokens map[string]bool          `json:"allowedTokens"` // token address hex -> allowed
}

// DeployParams contains the constructor parameters for app initialization
type DeployParams struct {
	AllowedTokens []string `json:"allowedTokens"` // list of token address hex strings to allow
}

// WithdrawInstruction represents instructions for withdrawing funds
type WithdrawInstruction struct {
	To           types.Address  `json:"to"`
	TokenAddress types.Address  `json:"tokenAddress"`
	Amount       *types.Uint256 `json:"amount"`
}

type CompareInstructions struct {
	TargetAddress types.Address `json:"targetAddress"`
	TokenAddress  types.Address `json:"tokenAddress,omitempty"` // token to compare; defaults to ETH (0x0)
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
	Tag           string                   `json:"tag,omitempty"`
	Accounts      map[string]*AccountState `json:"accounts"`
	AllowedTokens map[string]bool          `json:"allowedTokens"`
}

type DepositEvent struct {
	Type         string         `json:"type"`
	TokenAddress types.Address  `json:"tokenAddress"`
	Amount       *types.Uint256 `json:"amount"`
	Balance      *types.Uint256 `json:"balance"`
	Nonce        uint64         `json:"nonce"`
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

type WithdrawalEvent struct {
	Type         string         `json:"type"`
	To           types.Address  `json:"to"`
	TokenAddress types.Address  `json:"tokenAddress"`
	Amount       *types.Uint256 `json:"amount"`
	Balance      *types.Uint256 `json:"balance"`
}
