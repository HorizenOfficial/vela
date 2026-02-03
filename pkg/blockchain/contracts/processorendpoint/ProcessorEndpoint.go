// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package processorendpoint

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// StructsPendingRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsPendingRequest struct {
	ProtocolVersion uint8
	ApplicationId   uint64
	RequestType     uint8
	RequestId       [32]byte
	Payload         []byte
	Timestamp       *big.Int
	Sender          common.Address
	DepositAmount   *big.Int
	MaxFeeValue     *big.Int
}

// StructsWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsWithdrawalRequest struct {
	Receiver common.Address
	Amount   *big.Int
}

// AccessControlMetaData contains all meta data concerning the AccessControl contract.
var AccessControlMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "36b348cf3ce06e4a09945070dcc23e2db2",
}

// AccessControl is an auto generated Go binding around an Ethereum contract.
type AccessControl struct {
	abi abi.ABI
}

// NewAccessControl creates a new instance of AccessControl.
func NewAccessControl() *AccessControl {
	parsed, err := AccessControlMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &AccessControl{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *AccessControl) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (accessControl *AccessControl) PackDEFAULTADMINROLE() []byte {
	enc, err := accessControl.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (accessControl *AccessControl) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return accessControl.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (accessControl *AccessControl) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := accessControl.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (accessControl *AccessControl) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := accessControl.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (accessControl *AccessControl) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return accessControl.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (accessControl *AccessControl) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := accessControl.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (accessControl *AccessControl) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := accessControl.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (accessControl *AccessControl) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return accessControl.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (accessControl *AccessControl) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := accessControl.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (accessControl *AccessControl) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return accessControl.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (accessControl *AccessControl) UnpackHasRole(data []byte) (bool, error) {
	out, err := accessControl.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (accessControl *AccessControl) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := accessControl.abi.Pack("renounceRole", role, callerConfirmation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (accessControl *AccessControl) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return accessControl.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (accessControl *AccessControl) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := accessControl.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (accessControl *AccessControl) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return accessControl.abi.Pack("revokeRole", role, account)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (accessControl *AccessControl) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := accessControl.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (accessControl *AccessControl) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return accessControl.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (accessControl *AccessControl) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := accessControl.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// AccessControlRoleAdminChanged represents a RoleAdminChanged event raised by the AccessControl contract.
type AccessControlRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const AccessControlRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (AccessControlRoleAdminChanged) ContractEventName() string {
	return AccessControlRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (accessControl *AccessControl) UnpackRoleAdminChangedEvent(log *types.Log) (*AccessControlRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if log.Topics[0] != accessControl.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AccessControlRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := accessControl.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range accessControl.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// AccessControlRoleGranted represents a RoleGranted event raised by the AccessControl contract.
type AccessControlRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const AccessControlRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (AccessControlRoleGranted) ContractEventName() string {
	return AccessControlRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (accessControl *AccessControl) UnpackRoleGrantedEvent(log *types.Log) (*AccessControlRoleGranted, error) {
	event := "RoleGranted"
	if log.Topics[0] != accessControl.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AccessControlRoleGranted)
	if len(log.Data) > 0 {
		if err := accessControl.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range accessControl.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// AccessControlRoleRevoked represents a RoleRevoked event raised by the AccessControl contract.
type AccessControlRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const AccessControlRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (AccessControlRoleRevoked) ContractEventName() string {
	return AccessControlRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (accessControl *AccessControl) UnpackRoleRevokedEvent(log *types.Log) (*AccessControlRoleRevoked, error) {
	event := "RoleRevoked"
	if log.Topics[0] != accessControl.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AccessControlRoleRevoked)
	if len(log.Data) > 0 {
		if err := accessControl.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range accessControl.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (accessControl *AccessControl) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], accessControl.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return accessControl.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], accessControl.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return accessControl.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// AccessControlAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the AccessControl contract.
type AccessControlAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func AccessControlAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (accessControl *AccessControl) UnpackAccessControlBadConfirmationError(raw []byte) (*AccessControlAccessControlBadConfirmation, error) {
	out := new(AccessControlAccessControlBadConfirmation)
	if err := accessControl.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AccessControlAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the AccessControl contract.
type AccessControlAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func AccessControlAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (accessControl *AccessControl) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*AccessControlAccessControlUnauthorizedAccount, error) {
	out := new(AccessControlAccessControlUnauthorizedAccount)
	if err := accessControl.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryMetaData contains all meta data concerning the AuthorityRegistry contract.
var AuthorityRegistryMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"defaultAuthority\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"AppAuthorityContractSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"DefaultAuthorityContractSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"appAuthorityContracts\",\"outputs\":[{\"internalType\":\"contractIAuthorityChecker\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"checkAuthorityIsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultAuthorityContract\",\"outputs\":[{\"internalType\":\"contractIAuthorityChecker\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"setAppAuthorityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"setDefaultAuthorityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "38f83b752dabbc1250a3a7a867ed9a9f7b",
	Bin: "0x6080604052346100305761001a610014610104565b906101f4565b610022610035565b610a076103b68239610a0790f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b60018060a01b031690565b6100b19061009d565b90565b6100bd816100a8565b036100c457565b5f80fd5b905051906100d5826100b4565b565b91906040838203126100ff57806100f36100fc925f86016100c8565b936020016100c8565b90565b610099565b610122610dbd8038038061011781610084565b9283398101906100d7565b9091565b90565b90565b61014061013b61014592610126565b610129565b61009d565b90565b6101519061012c565b90565b5f0190565b61016d6101686101729261009d565b610129565b61009d565b90565b61017e90610159565b90565b61018a90610175565b90565b5f1b90565b906101a360018060a01b039161018d565b9181191691161790565b6101b690610175565b90565b90565b906101d16101cc6101d8926101ad565b6101b9565b8254610192565b9055565b6101e590610159565b90565b6101f1906101dc565b90565b6101fd906102b0565b8061021861021261020d5f610148565b6100a8565b916100a8565b146102725761023061022982610181565b60026101bc565b61025a7fd4047f78dd943d75dfa94ddb7e36315e62df8a01f00a3d0f512960bf93c12318916101e8565b90610263610035565b8061026d81610154565b0390a2565b5f632582a64160e11b81528061028a60048201610154565b0390fd5b610297906100a8565b9052565b91906102ae905f6020850194019061028e565b565b806102cb6102c56102c05f610148565b6100a8565b916100a8565b146102db576102d990610356565b565b6102fe6102e75f610148565b5f918291631e4fbdf760e01b83526004830161029b565b0390fd5b5f1c90565b60018060a01b031690565b61031e61032391610302565b610307565b90565b6103309054610312565b90565b90565b9061034b610346610352926101e8565b610333565b8254610192565b9055565b61035f5f610326565b610369825f610336565b9061039d6103977f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936101e8565b916101e8565b916103a6610035565b806103b081610154565b0390a356fe60806040526004361015610013575b610460565b61001d5f3561009c565b8063116ad4a6146100975780633aac66d1146100925780635d35bc141461008d5780636f5e32f214610088578063715018a6146100835780638da5cb5b1461007e578063cfa8904f146100795763f2fde38b0361000e5761042d565b6103f8565b6103b4565b61035f565b61031b565b6101f1565b61019f565b610164565b60e01c90565b60405190565b5f80fd5b5f80fd5b90565b6100bc816100b0565b036100c357565b5f80fd5b905035906100d4826100b3565b565b60018060a01b031690565b6100ea906100d6565b90565b6100f6816100e1565b036100fd57565b5f80fd5b9050359061010e826100ed565b565b9190604083820312610138578061012c610135925f86016100c7565b93602001610101565b90565b6100ac565b151590565b61014b9061013d565b9052565b9190610162905f60208501940190610142565b565b346101955761019161018061017a366004610110565b90610587565b6101886100a2565b9182918261014f565b0390f35b6100a8565b5f0190565b346101ce576101b86101b2366004610110565b90610755565b6101c06100a2565b806101ca8161019a565b0390f35b6100a8565b906020828203126101ec576101e9915f01610101565b90565b6100ac565b3461021f576102096102043660046101d3565b610805565b6102116100a2565b8061021b8161019a565b0390f35b6100a8565b9060208282031261023d5761023a915f016100c7565b90565b6100ac565b90565b61025961025461025e926100b0565b610242565b6100b0565b90565b9061026b90610245565b5f5260205260405f2090565b1c90565b60018060a01b031690565b61029690600861029b9302610277565b61027b565b90565b906102a99154610286565b90565b6102c2906102bd6001915f92610261565b61029e565b90565b6102d96102d46102de926100d6565b610242565b6100d6565b90565b6102ea906102c5565b90565b6102f6906102e1565b90565b610302906102ed565b9052565b9190610319905f602085019401906102f9565b565b3461034b57610347610336610331366004610224565b6102ac565b61033e6100a2565b91829182610306565b0390f35b6100a8565b5f91031261035a57565b6100ac565b3461038d5761036f366004610350565b610377610835565b61037f6100a2565b806103898161019a565b0390f35b6100a8565b61039b906100e1565b9052565b91906103b2905f60208501940190610392565b565b346103e4576103c4366004610350565b6103e06103cf61086f565b6103d76100a2565b9182918261039f565b0390f35b6100a8565b6103f560025f9061029e565b90565b3461042857610408366004610350565b6104246104136103e9565b61041b6100a2565b91829182610306565b0390f35b6100a8565b3461045b576104456104403660046101d3565b6108e9565b61044d6100a2565b806104578161019a565b0390f35b6100a8565b5f80fd5b5f90565b5f1c90565b61047961047e91610468565b61027b565b90565b61048b905461046d565b90565b90565b6104a56104a06104aa9261048e565b610242565b6100d6565b90565b6104b690610491565b90565b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906104e1906104b9565b810190811067ffffffffffffffff8211176104fb57604052565b6104c3565b60e01b90565b61050f8161013d565b0361051657565b5f80fd5b9050519061052782610506565b565b906020828203126105425761053f915f0161051a565b90565b6100ac565b610550906100b0565b9052565b91602061057592949361056e60408201965f830190610547565b0190610392565b565b61057f6100a2565b3d5f823e3d90fd5b90602090610593610464565b506105a86105a360018590610261565b610481565b6105b1816102ed565b6105cb6105c56105c05f6104ad565b6100e1565b916100e1565b14610650575b6105da906102ed565b6105fc63116ad4a69492946106076105f06100a2565b96879586948594610500565b845260048401610554565b03915afa90811561064b575f9161061d575b5090565b61063e915060203d8111610644575b61063681836104d7565b810190610529565b5f610619565b503d61062c565b610577565b506105da61065e6002610481565b90506105d1565b90610677916106726108f4565b6106ec565b565b610682906102c5565b90565b61068e90610679565b90565b5f1b90565b906106a760018060a01b0391610691565b9181191691161790565b6106ba90610679565b90565b90565b906106d56106d06106dc926106b1565b6106bd565b8254610696565b9055565b6106e9906102e1565b90565b6107096106f883610685565b61070460018490610261565b6106c0565b9061073d6107377f4f79164bfe5c672976dae578ab2d7b0c99754246947eb57d8fd2a5a078307c2293610245565b916106e0565b916107466100a2565b806107508161019a565b0390a3565b9061075f91610665565b565b6107729061076d6108f4565b610774565b565b8061078f6107896107845f6104ad565b6100e1565b916100e1565b146107e9576107a76107a082610685565b60026106c0565b6107d17fd4047f78dd943d75dfa94ddb7e36315e62df8a01f00a3d0f512960bf93c12318916106e0565b906107da6100a2565b806107e48161019a565b0390a2565b5f632582a64160e11b8152806108016004820161019a565b0390fd5b61080e90610761565b565b6108186108f4565b610820610822565b565b61083361082e5f6104ad565b610965565b565b61083d610810565b565b5f90565b60018060a01b031690565b61085a61085f91610468565b610843565b90565b61086c905461084e565b90565b61087761083f565b506108815f610862565b90565b610895906108906108f4565b610897565b565b806108b26108ac6108a75f6104ad565b6100e1565b916100e1565b146108c2576108c090610965565b565b6108e56108ce5f6104ad565b5f918291631e4fbdf760e01b83526004830161039f565b0390fd5b6108f290610884565b565b6108fc61086f565b61091561090f61090a6109c4565b6100e1565b916100e1565b0361091c57565b61093e6109276109c4565b5f91829163118cdaa760e01b83526004830161039f565b0390fd5b90565b9061095a610955610961926106e0565b610942565b8254610696565b9055565b61096e5f610862565b610978825f610945565b906109ac6109a67f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936106e0565b916106e0565b916109b56100a2565b806109bf8161019a565b0390a3565b6109cc61083f565b50339056fea2646970667358221220a68649e8da4d9b0ca1af26fa630209f10384552574db9a6b833eed76a0a9adc164736f6c634300081e0033",
}

// AuthorityRegistry is an auto generated Go binding around an Ethereum contract.
type AuthorityRegistry struct {
	abi abi.ABI
}

// NewAuthorityRegistry creates a new instance of AuthorityRegistry.
func NewAuthorityRegistry() *AuthorityRegistry {
	parsed, err := AuthorityRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &AuthorityRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *AuthorityRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address owner, address defaultAuthority) returns()
func (authorityRegistry *AuthorityRegistry) PackConstructor(owner common.Address, defaultAuthority common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("", owner, defaultAuthority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAppAuthorityContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f5e32f2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function appAuthorityContracts(uint256 ) view returns(address)
func (authorityRegistry *AuthorityRegistry) PackAppAuthorityContracts(arg0 *big.Int) []byte {
	enc, err := authorityRegistry.abi.Pack("appAuthorityContracts", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAppAuthorityContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f5e32f2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function appAuthorityContracts(uint256 ) view returns(address)
func (authorityRegistry *AuthorityRegistry) TryPackAppAuthorityContracts(arg0 *big.Int) ([]byte, error) {
	return authorityRegistry.abi.Pack("appAuthorityContracts", arg0)
}

// UnpackAppAuthorityContracts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6f5e32f2.
//
// Solidity: function appAuthorityContracts(uint256 ) view returns(address)
func (authorityRegistry *AuthorityRegistry) UnpackAppAuthorityContracts(data []byte) (common.Address, error) {
	out, err := authorityRegistry.abi.Unpack("appAuthorityContracts", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x116ad4a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) PackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x116ad4a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) TryPackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
}

// UnpackCheckAuthorityIsAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x116ad4a6.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) UnpackCheckAuthorityIsAllowed(data []byte) (bool, error) {
	out, err := authorityRegistry.abi.Unpack("checkAuthorityIsAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcfa8904f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAuthorityContract() view returns(address)
func (authorityRegistry *AuthorityRegistry) PackDefaultAuthorityContract() []byte {
	enc, err := authorityRegistry.abi.Pack("defaultAuthorityContract")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcfa8904f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAuthorityContract() view returns(address)
func (authorityRegistry *AuthorityRegistry) TryPackDefaultAuthorityContract() ([]byte, error) {
	return authorityRegistry.abi.Pack("defaultAuthorityContract")
}

// UnpackDefaultAuthorityContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcfa8904f.
//
// Solidity: function defaultAuthorityContract() view returns(address)
func (authorityRegistry *AuthorityRegistry) UnpackDefaultAuthorityContract(data []byte) (common.Address, error) {
	out, err := authorityRegistry.abi.Unpack("defaultAuthorityContract", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (authorityRegistry *AuthorityRegistry) PackOwner() []byte {
	enc, err := authorityRegistry.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (authorityRegistry *AuthorityRegistry) TryPackOwner() ([]byte, error) {
	return authorityRegistry.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (authorityRegistry *AuthorityRegistry) UnpackOwner(data []byte) (common.Address, error) {
	out, err := authorityRegistry.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (authorityRegistry *AuthorityRegistry) PackRenounceOwnership() []byte {
	enc, err := authorityRegistry.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() returns()
func (authorityRegistry *AuthorityRegistry) TryPackRenounceOwnership() ([]byte, error) {
	return authorityRegistry.abi.Pack("renounceOwnership")
}

// PackSetAppAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3aac66d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setAppAuthorityContract(uint256 applicationId, address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) PackSetAppAuthorityContract(applicationId *big.Int, authorityContract common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("setAppAuthorityContract", applicationId, authorityContract)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetAppAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3aac66d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setAppAuthorityContract(uint256 applicationId, address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) TryPackSetAppAuthorityContract(applicationId *big.Int, authorityContract common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("setAppAuthorityContract", applicationId, authorityContract)
}

// PackSetDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d35bc14.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setDefaultAuthorityContract(address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) PackSetDefaultAuthorityContract(authorityContract common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("setDefaultAuthorityContract", authorityContract)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d35bc14.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setDefaultAuthorityContract(address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) TryPackSetDefaultAuthorityContract(authorityContract common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("setDefaultAuthorityContract", authorityContract)
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (authorityRegistry *AuthorityRegistry) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (authorityRegistry *AuthorityRegistry) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("transferOwnership", newOwner)
}

// AuthorityRegistryAppAuthorityContractSet represents a AppAuthorityContractSet event raised by the AuthorityRegistry contract.
type AuthorityRegistryAppAuthorityContractSet struct {
	ApplicationId     *big.Int
	AuthorityContract common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryAppAuthorityContractSetEventName = "AppAuthorityContractSet"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryAppAuthorityContractSet) ContractEventName() string {
	return AuthorityRegistryAppAuthorityContractSetEventName
}

// UnpackAppAuthorityContractSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AppAuthorityContractSet(uint256 indexed applicationId, address indexed authorityContract)
func (authorityRegistry *AuthorityRegistry) UnpackAppAuthorityContractSetEvent(log *types.Log) (*AuthorityRegistryAppAuthorityContractSet, error) {
	event := "AppAuthorityContractSet"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryAppAuthorityContractSet)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// AuthorityRegistryDefaultAuthorityContractSet represents a DefaultAuthorityContractSet event raised by the AuthorityRegistry contract.
type AuthorityRegistryDefaultAuthorityContractSet struct {
	AuthorityContract common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryDefaultAuthorityContractSetEventName = "DefaultAuthorityContractSet"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryDefaultAuthorityContractSet) ContractEventName() string {
	return AuthorityRegistryDefaultAuthorityContractSetEventName
}

// UnpackDefaultAuthorityContractSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAuthorityContractSet(address indexed authorityContract)
func (authorityRegistry *AuthorityRegistry) UnpackDefaultAuthorityContractSetEvent(log *types.Log) (*AuthorityRegistryDefaultAuthorityContractSet, error) {
	event := "DefaultAuthorityContractSet"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryDefaultAuthorityContractSet)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// AuthorityRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the AuthorityRegistry contract.
type AuthorityRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryOwnershipTransferred) ContractEventName() string {
	return AuthorityRegistryOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (authorityRegistry *AuthorityRegistry) UnpackOwnershipTransferredEvent(log *types.Log) (*AuthorityRegistryOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (authorityRegistry *AuthorityRegistry) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["AddressCantBeZero"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// AuthorityRegistryAddressCantBeZero represents a AddressCantBeZero error raised by the AuthorityRegistry contract.
type AuthorityRegistryAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressCantBeZero()
func AuthorityRegistryAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x4b054c82bb9e4ebe90744902747c4e86028dd2a122c63de8810d9a1c2e84617d")
}

// UnpackAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressCantBeZero()
func (authorityRegistry *AuthorityRegistry) UnpackAddressCantBeZeroError(raw []byte) (*AuthorityRegistryAddressCantBeZero, error) {
	out := new(AuthorityRegistryAddressCantBeZero)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "AddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the AuthorityRegistry contract.
type AuthorityRegistryOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func AuthorityRegistryOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (authorityRegistry *AuthorityRegistry) UnpackOwnableInvalidOwnerError(raw []byte) (*AuthorityRegistryOwnableInvalidOwner, error) {
	out := new(AuthorityRegistryOwnableInvalidOwner)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the AuthorityRegistry contract.
type AuthorityRegistryOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func AuthorityRegistryOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (authorityRegistry *AuthorityRegistry) UnpackOwnableUnauthorizedAccountError(raw []byte) (*AuthorityRegistryOwnableUnauthorizedAccount, error) {
	out := new(AuthorityRegistryOwnableUnauthorizedAccount)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContextMetaData contains all meta data concerning the Context contract.
var ContextMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "c3b2d42aaf903aabf7161bf9c52edde51b",
}

// Context is an auto generated Go binding around an Ethereum contract.
type Context struct {
	abi abi.ABI
}

// NewContext creates a new instance of Context.
func NewContext() *Context {
	parsed, err := ContextMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Context{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Context) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// ERC165MetaData contains all meta data concerning the ERC165 contract.
var ERC165MetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "5eb3ce5f671914ab9adbddb5d74ba39063",
}

// ERC165 is an auto generated Go binding around an Ethereum contract.
type ERC165 struct {
	abi abi.ABI
}

// NewERC165 creates a new instance of ERC165.
func NewERC165() *ERC165 {
	parsed, err := ERC165MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ERC165{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ERC165) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (eRC165 *ERC165) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := eRC165.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (eRC165 *ERC165) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return eRC165.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (eRC165 *ERC165) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := eRC165.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// IAccessControlMetaData contains all meta data concerning the IAccessControl contract.
var IAccessControlMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "1a6e18dc8923e17df081f17b4aabc18c20",
}

// IAccessControl is an auto generated Go binding around an Ethereum contract.
type IAccessControl struct {
	abi abi.ABI
}

// NewIAccessControl creates a new instance of IAccessControl.
func NewIAccessControl() *IAccessControl {
	parsed, err := IAccessControlMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IAccessControl{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IAccessControl) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (iAccessControl *IAccessControl) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := iAccessControl.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (iAccessControl *IAccessControl) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return iAccessControl.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (iAccessControl *IAccessControl) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := iAccessControl.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (iAccessControl *IAccessControl) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := iAccessControl.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (iAccessControl *IAccessControl) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return iAccessControl.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (iAccessControl *IAccessControl) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := iAccessControl.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (iAccessControl *IAccessControl) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return iAccessControl.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (iAccessControl *IAccessControl) UnpackHasRole(data []byte) (bool, error) {
	out, err := iAccessControl.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (iAccessControl *IAccessControl) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := iAccessControl.abi.Pack("renounceRole", role, callerConfirmation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (iAccessControl *IAccessControl) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return iAccessControl.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (iAccessControl *IAccessControl) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := iAccessControl.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (iAccessControl *IAccessControl) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return iAccessControl.abi.Pack("revokeRole", role, account)
}

// IAccessControlRoleAdminChanged represents a RoleAdminChanged event raised by the IAccessControl contract.
type IAccessControlRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const IAccessControlRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (IAccessControlRoleAdminChanged) ContractEventName() string {
	return IAccessControlRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (iAccessControl *IAccessControl) UnpackRoleAdminChangedEvent(log *types.Log) (*IAccessControlRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if log.Topics[0] != iAccessControl.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IAccessControlRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := iAccessControl.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAccessControl.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IAccessControlRoleGranted represents a RoleGranted event raised by the IAccessControl contract.
type IAccessControlRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const IAccessControlRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (IAccessControlRoleGranted) ContractEventName() string {
	return IAccessControlRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (iAccessControl *IAccessControl) UnpackRoleGrantedEvent(log *types.Log) (*IAccessControlRoleGranted, error) {
	event := "RoleGranted"
	if log.Topics[0] != iAccessControl.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IAccessControlRoleGranted)
	if len(log.Data) > 0 {
		if err := iAccessControl.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAccessControl.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IAccessControlRoleRevoked represents a RoleRevoked event raised by the IAccessControl contract.
type IAccessControlRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const IAccessControlRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (IAccessControlRoleRevoked) ContractEventName() string {
	return IAccessControlRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (iAccessControl *IAccessControl) UnpackRoleRevokedEvent(log *types.Log) (*IAccessControlRoleRevoked, error) {
	event := "RoleRevoked"
	if log.Topics[0] != iAccessControl.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IAccessControlRoleRevoked)
	if len(log.Data) > 0 {
		if err := iAccessControl.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAccessControl.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (iAccessControl *IAccessControl) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iAccessControl.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return iAccessControl.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAccessControl.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return iAccessControl.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IAccessControlAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the IAccessControl contract.
type IAccessControlAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func IAccessControlAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (iAccessControl *IAccessControl) UnpackAccessControlBadConfirmationError(raw []byte) (*IAccessControlAccessControlBadConfirmation, error) {
	out := new(IAccessControlAccessControlBadConfirmation)
	if err := iAccessControl.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAccessControlAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the IAccessControl contract.
type IAccessControlAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func IAccessControlAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (iAccessControl *IAccessControl) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*IAccessControlAccessControlUnauthorizedAccount, error) {
	out := new(IAccessControlAccessControlUnauthorizedAccount)
	if err := iAccessControl.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAuthorityCheckerMetaData contains all meta data concerning the IAuthorityChecker contract.
var IAuthorityCheckerMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"checkAuthorityIsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "40bb6b80f72bccedfb46120296ba91f2f6",
}

// IAuthorityChecker is an auto generated Go binding around an Ethereum contract.
type IAuthorityChecker struct {
	abi abi.ABI
}

// NewIAuthorityChecker creates a new instance of IAuthorityChecker.
func NewIAuthorityChecker() *IAuthorityChecker {
	parsed, err := IAuthorityCheckerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IAuthorityChecker{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IAuthorityChecker) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x116ad4a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (iAuthorityChecker *IAuthorityChecker) PackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) []byte {
	enc, err := iAuthorityChecker.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x116ad4a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (iAuthorityChecker *IAuthorityChecker) TryPackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) ([]byte, error) {
	return iAuthorityChecker.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
}

// UnpackCheckAuthorityIsAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x116ad4a6.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (iAuthorityChecker *IAuthorityChecker) UnpackCheckAuthorityIsAllowed(data []byte) (bool, error) {
	out, err := iAuthorityChecker.abi.Unpack("checkAuthorityIsAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// IERC165MetaData contains all meta data concerning the IERC165 contract.
var IERC165MetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "a33a1eaf2949eb84428aeece5515d1fc34",
}

// IERC165 is an auto generated Go binding around an Ethereum contract.
type IERC165 struct {
	abi abi.ABI
}

// NewIERC165 creates a new instance of IERC165.
func NewIERC165() *IERC165 {
	parsed, err := IERC165MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IERC165{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IERC165) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (iERC165 *IERC165) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := iERC165.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (iERC165 *IERC165) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return iERC165.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (iERC165 *IERC165) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := iERC165.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// ITeeAuthenticatorMetaData contains all meta data concerning the ITeeAuthenticator contract.
var ITeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "1f0ecb3d874bbc852bff29c87537425ade",
}

// ITeeAuthenticator is an auto generated Go binding around an Ethereum contract.
type ITeeAuthenticator struct {
	abi abi.ABI
}

// NewITeeAuthenticator creates a new instance of ITeeAuthenticator.
func NewITeeAuthenticator() *ITeeAuthenticator {
	parsed, err := ITeeAuthenticatorMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ITeeAuthenticator{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ITeeAuthenticator) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x64c062d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refundAmount, applicationFee, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x64c062d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refundAmount, applicationFee, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x64c062d1.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := iTeeAuthenticator.abi.Unpack("checkSignature", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackGetPubSecp521r1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081bec7e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (iTeeAuthenticator *ITeeAuthenticator) PackGetPubSecp521r1() []byte {
	enc, err := iTeeAuthenticator.abi.Pack("getPubSecp521r1")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPubSecp521r1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081bec7e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackGetPubSecp521r1() ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("getPubSecp521r1")
}

// UnpackGetPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081bec7e.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (iTeeAuthenticator *ITeeAuthenticator) UnpackGetPubSecp521r1(data []byte) ([]byte, error) {
	out, err := iTeeAuthenticator.abi.Unpack("getPubSecp521r1", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetTeeSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dd7ce2f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getTeeSigner() view returns(address)
func (iTeeAuthenticator *ITeeAuthenticator) PackGetTeeSigner() []byte {
	enc, err := iTeeAuthenticator.abi.Pack("getTeeSigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetTeeSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dd7ce2f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getTeeSigner() view returns(address)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackGetTeeSigner() ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("getTeeSigner")
}

// UnpackGetTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dd7ce2f.
//
// Solidity: function getTeeSigner() view returns(address)
func (iTeeAuthenticator *ITeeAuthenticator) UnpackGetTeeSigner(data []byte) (common.Address, error) {
	out, err := iTeeAuthenticator.abi.Unpack("getTeeSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// OwnableMetaData contains all meta data concerning the Ownable contract.
var OwnableMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "4ecec2844bcbcedb4addb46b016e5646ea",
}

// Ownable is an auto generated Go binding around an Ethereum contract.
type Ownable struct {
	abi abi.ABI
}

// NewOwnable creates a new instance of Ownable.
func NewOwnable() *Ownable {
	parsed, err := OwnableMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Ownable{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Ownable) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (ownable *Ownable) PackOwner() []byte {
	enc, err := ownable.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (ownable *Ownable) TryPackOwner() ([]byte, error) {
	return ownable.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (ownable *Ownable) UnpackOwner(data []byte) (common.Address, error) {
	out, err := ownable.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (ownable *Ownable) PackRenounceOwnership() []byte {
	enc, err := ownable.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() returns()
func (ownable *Ownable) TryPackRenounceOwnership() ([]byte, error) {
	return ownable.abi.Pack("renounceOwnership")
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (ownable *Ownable) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := ownable.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (ownable *Ownable) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return ownable.abi.Pack("transferOwnership", newOwner)
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const OwnableOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (OwnableOwnershipTransferred) ContractEventName() string {
	return OwnableOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (ownable *Ownable) UnpackOwnershipTransferredEvent(log *types.Log) (*OwnableOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if log.Topics[0] != ownable.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(OwnableOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := ownable.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ownable.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (ownable *Ownable) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], ownable.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return ownable.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], ownable.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return ownable.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// OwnableOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the Ownable contract.
type OwnableOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func OwnableOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (ownable *Ownable) UnpackOwnableInvalidOwnerError(raw []byte) (*OwnableOwnableInvalidOwner, error) {
	out := new(OwnableOwnableInvalidOwner)
	if err := ownable.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// OwnableOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the Ownable contract.
type OwnableOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func OwnableOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (ownable *Ownable) UnpackOwnableUnauthorizedAccountError(raw []byte) (*OwnableOwnableUnauthorizedAccount, error) {
	out := new(OwnableOwnableUnauthorizedAccount)
	if err := ownable.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointMetaData contains all meta data concerning the ProcessorEndpoint contract.
var ProcessorEndpointMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"ReportGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"eventSubType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"APPLICATION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"markRequestFailed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a35c5449187bc72768eac252f2f30417e4",
	Bin: "0x6080604052346100335761001d6100146101b4565b939290926103f9565b610025610038565b6144206106db823961442090f35b61003e565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b9061006a90610042565b810190811060018060401b0382111761008257604052565b61004c565b9061009a610093610038565b9283610060565b565b5f80fd5b60018060a01b031690565b6100b4906100a0565b90565b6100c0906100ab565b90565b6100cc816100b7565b036100d357565b5f80fd5b905051906100e4826100c3565b565b6100ef906100ab565b90565b6100fb816100e6565b0361010257565b5f80fd5b90505190610113826100f2565b565b61011e816100ab565b0361012557565b5f80fd5b9050519061013682610115565b565b90565b61014481610138565b0361014b57565b5f80fd5b9050519061015c8261013b565b565b919060a0838203126101af57610176815f85016100d7565b926101848260208301610106565b926101ac6101958460408501610129565b936101a38160608601610129565b9360800161014f565b90565b61009c565b6101d2614afb803803806101c781610087565b92833981019061015e565b9091929394565b5f1b90565b906101ea5f19916101d9565b9181191691161790565b90565b90565b61020e610209610213926101f4565b6101f7565b610138565b90565b90565b9061022e610229610235926101fa565b610216565b82546101de565b9055565b90565b61025061024b61025592610239565b6101f7565b6100a0565b90565b6102619061023c565b90565b61027861027361027d926100a0565b6101f7565b6100a0565b90565b61028990610264565b90565b61029590610280565b90565b6102a190610264565b90565b6102ad90610298565b90565b5f0190565b906102c660018060a01b03916101d9565b9181191691161790565b6102d990610280565b90565b90565b906102f46102ef6102fb926102d0565b6102dc565b82546102b5565b9055565b61030890610298565b90565b90565b9061032361031e61032a926102ff565b61030b565b82546102b5565b9055565b61033790610264565b90565b6103439061032e565b90565b61034f9061032e565b90565b90565b9061036a61036561037192610346565b610352565b82546102b5565b9055565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103d16103cc6103d692610138565b6101f7565b610138565b90565b906103ee6103e96103f5926103bd565b610216565b82546101de565b9055565b91939290610409600a6006610219565b8261042c61042661042161041c5f610258565b61028c565b6100b7565b916100b7565b1480156104fe575b80156104dc575b80156104ba575b61049e5761049c946104666104869261045f6104949660076102df565b600861030e565b6104796104728261033a565b600a610355565b610481610375565b6105c9565b5061048f610399565b6105c9565b5060096103d9565b565b5f632582a64160e11b8152806104b6600482016102b0565b0390fd5b50816104d66104d06104cb5f610258565b6100ab565b916100ab565b14610442565b50846104f86104f26104ed5f610258565b6100ab565b916100ab565b1461043b565b508061052261051c6105176105125f610258565b6102a4565b6100e6565b916100e6565b14610434565b5f90565b151590565b90565b61053d90610531565b90565b9061054a90610534565b5f5260205260405f2090565b61055f90610264565b90565b61056b90610556565b90565b9061057890610562565b5f5260205260405f2090565b9061059060ff916101d9565b9181191691161790565b6105a39061052c565b90565b90565b906105be6105b96105c59261059a565b6105a6565b8254610584565b9055565b6105d1610528565b506105e66105e08284906106a0565b1561052c565b5f1461066e5761060d60016106085f610600818690610540565b01859061056e565b6105a9565b906106166106cd565b9061065361064d6106477f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610534565b92610562565b92610562565b9261065c610038565b80610666816102b0565b0390a4600190565b50505f90565b5f1c90565b60ff1690565b61068b61069091610674565b610679565b90565b61069d905461067f565b90565b6106c6915f6106bb6106c1936106b4610528565b5082610540565b0161056e565b610693565b90565b5f90565b6106d56106c9565b50339056fe60806040526004361015610013575b611a5e565b61001d5f356101dc565b806301ffc9a7146101d75780631664deeb146101d2578063248a9ca3146101cd5780632a0acc6a146101c85780632f2ff15d146101c357806334f23a9d146101be57806336568abe146101b95780633b473cb9146101b45780635423112f146101af5780635e339384146101aa5780636e91d133146101a557806380a1f712146101a057806391d148541461019b5780639588eca2146101965780639968da96146101915780639a17995a1461018c5780639c73eeaf146101875780639c8d416f146101825780639fabc36f1461017d578063a217fddf14610178578063a9ba045e14610173578063a9ef290e1461016e578063aa3aa46014610169578063c404c35914610164578063c415b95c1461015f578063d2c35ce81461015a578063d547741f146101555763f7d9cc1e0361000e57611a29565b6119e6565b6119b3565b611960565b6118c6565b6117af565b611727565b61130d565b61126a565b6111fa565b6111c5565b61115d565b611103565b61104e565b610fcd565b610f62565b610f2d565b610d27565b610c55565b610c01565b6107d1565b610648565b610616565b61047b565b6103da565b610376565b610300565b610268565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f80fd5b63ffffffff60e01b1690565b610209816101f4565b0361021057565b5f80fd5b9050359061022182610200565b565b9060208282031261023c57610239915f01610214565b90565b6101ec565b151590565b61024f90610241565b9052565b9190610266905f60208501940190610246565b565b346102985761029461028361027e366004610223565b611a66565b61028b6101e2565b91829182610253565b0390f35b6101e8565b5f9103126102a757565b6101ec565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b6102d86102ac565b90565b90565b6102e7906102db565b9052565b91906102fe905f602085019401906102de565b565b346103305761031036600461029d565b61032c61031b6102d0565b6103236101e2565b918291826102eb565b0390f35b6101e8565b61033e816102db565b0361034557565b5f80fd5b9050359061035682610335565b565b906020828203126103715761036e915f01610349565b90565b6101ec565b346103a6576103a261039161038c366004610358565b611ac0565b6103996101e2565b918291826102eb565b0390f35b6101e8565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103d76103ab565b90565b3461040a576103ea36600461029d565b6104066103f56103cf565b6103fd6101e2565b918291826102eb565b0390f35b6101e8565b60018060a01b031690565b6104239061040f565b90565b61042f8161041a565b0361043657565b5f80fd5b9050359061044782610426565b565b9190604083820312610471578061046561046e925f8601610349565b9360200161043a565b90565b6101ec565b5f0190565b346104aa5761049461048e366004610449565b90611b0b565b61049c6101e2565b806104a681610476565b0390f35b6101e8565b60ff1690565b6104be816104af565b036104c557565b5f80fd5b905035906104d6826104b5565b565b67ffffffffffffffff1690565b6104ee816104d8565b036104f557565b5f80fd5b90503590610506826104e5565b565b6004111561051257565b5f80fd5b9050359061052382610508565b565b5f80fd5b5f80fd5b5f80fd5b909182601f8301121561056b5781359167ffffffffffffffff831161056657602001926001830284011161056157565b61052d565b610529565b610525565b90565b61057c81610570565b0361058357565b5f80fd5b9050359061059482610573565b565b91909160c081840312610611576105af835f83016104c9565b926105bd81602084016104f9565b926105cb8260408501610516565b9260608101359167ffffffffffffffff831161060c576105f084610609948401610531565b9390946106008160808601610587565b9360a001610587565b90565b6101f0565b6101ec565b610644610633610627366004610596565b959490949391936125da565b61063b6101e2565b918291826102eb565b0390f35b346106775761066161065b366004610449565b906125f4565b6106696101e2565b8061067381610476565b0390f35b6101e8565b6011111561068657565b5f80fd5b905035906106978261067c565b565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906106c59061069d565b810190811067ffffffffffffffff8211176106df57604052565b6106a7565b906106f76106f06101e2565b92836106bb565b565b67ffffffffffffffff81116107175761071360209161069d565b0190565b6106a7565b90825f939282370152565b9092919261073c610737826106f9565b6106e4565b93818552602085019082840111610758576107569261071c565b565b610699565b9080601f8301121561077b5781602061077893359101610727565b90565b610525565b916060838303126107cc57610797825f8501610349565b926107a5836020830161068a565b92604082013567ffffffffffffffff81116107c7576107c4920161075d565b90565b6101f0565b6101ec565b34610800576107ea6107e4366004610780565b91612aa0565b6107f26101e2565b806107fc81610476565b0390f35b6101e8565b61080e906102db565b90565b9061081b90610805565b5f5260205260405f2090565b5f1c90565b60ff1690565b61083e61084391610827565b61082c565b90565b6108509054610832565b90565b60081c90565b67ffffffffffffffff1690565b61087261087791610853565b610859565b90565b6108849054610866565b90565b60481c90565b60ff1690565b61089f6108a491610887565b61088d565b90565b6108b19054610893565b90565b90565b6108c36108c891610827565b6108b4565b90565b6108d590546108b7565b90565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561090c575b602083101461090757565b6108d8565b91607f16916108fc565b60209181520190565b5f5260205f2090565b905f929180549061094261093b836108ec565b8094610916565b916001811690815f14610999575060011461095d575b505050565b61096a919293945061091f565b915f925b81841061098157505001905f8080610958565b6001816020929593955484860152019101929061096e565b92949550505060ff19168252151560200201905f8080610958565b906109be91610928565b90565b906109e16109da926109d16101e2565b938480926109b4565b03836106bb565b565b90565b6109f26109f791610827565b6109e3565b90565b610a0490546109e6565b90565b60018060a01b031690565b610a1e610a2391610827565b610a07565b90565b610a309054610a12565b90565b610a3e906002610811565b610a495f8201610846565b91610a555f830161087a565b91610a615f82016108a7565b91610a6e600183016108cb565b91610a7b600282016109c1565b91610a88600383016109fa565b91610a9560048201610a26565b91610aae6006610aa7600585016109fa565b93016109fa565b90565b610aba906104af565b9052565b610ac7906104d8565b9052565b634e487b7160e01b5f52602160045260245ffd5b60041115610ae957565b610acb565b90610af882610adf565b565b610b0390610aee565b90565b610b0f90610afa565b9052565b5190565b60209181520190565b90825f9392825e0152565b610b4a610b53602093610b5893610b4181610b13565b93848093610b17565b95869101610b20565b61069d565b0190565b610b6590610570565b9052565b610b729061041a565b9052565b94610bd9610100979b9a9892610be492610bcc610bf898610bc2610bff9e99610bb8610bee9a60208f610bb161012082019a5f830190610ab1565b0190610abe565b60408d0190610b06565b60608b01906102de565b88820360808a0152610b2b565b9a60a0870190610b5c565b60c0850190610b69565b60e0830190610b5c565b0190610b5c565b565b34610c3b57610c37610c1c610c17366004610358565b610a33565b95610c2e9997999591959492946101e2565b998a998a610b76565b0390f35b6101e8565b9190610c53905f60208501940190610b5c565b565b34610c8557610c6536600461029d565b610c81610c70612aad565b610c786101e2565b91829182610c40565b0390f35b6101e8565b1c90565b60018060a01b031690565b610ca9906008610cae9302610c8a565b610c8e565b90565b90610cbc9154610c99565b90565b610ccb60075f90610cb1565b90565b90565b610ce5610ce0610cea9261040f565b610cce565b61040f565b90565b610cf690610cd1565b90565b610d0290610ced565b90565b610d0e90610cf9565b9052565b9190610d25905f60208501940190610d05565b565b34610d5757610d3736600461029d565b610d53610d42610cbf565b610d4a6101e2565b91829182610d12565b0390f35b6101e8565b5190565b60209181520190565b60200190565b610d78906104af565b9052565b610d85906104d8565b9052565b610d9290610afa565b9052565b610d9f906102db565b9052565b610dc2610dcb602093610dd093610db981610b13565b93848093610916565b95869101610b20565b61069d565b0190565b610ddd90610570565b9052565b610dea9061041a565b9052565b90610e989061010080610e576101208401610e0f5f8801515f870190610d6f565b610e2160208801516020870190610d7c565b610e3360408801516040870190610d89565b610e4560608801516060870190610d96565b60808701518582036080870152610da3565b94610e6a60a082015160a0860190610dd4565b610e7c60c082015160c0860190610de1565b610e8e60e082015160e0860190610dd4565b0151910190610dd4565b90565b90610ea591610dee565b90565b60200190565b90610ec2610ebb83610d5c565b8092610d60565b9081610ed360208302840194610d69565b925f915b838310610ee657505050505090565b90919293946020610f08610f0283856001950387528951610e9b565b97610ea8565b9301930191939290610ed7565b610f2a9160208201915f818403910152610eae565b90565b34610f5d57610f3d36600461029d565b610f59610f48612d70565b610f506101e2565b91829182610f15565b0390f35b6101e8565b34610f9357610f8f610f7e610f78366004610449565b90612e5b565b610f866101e2565b91829182610253565b0390f35b6101e8565b610fa8906008610fad9302610c8a565b6108b4565b90565b90610fbb9154610f98565b90565b610fca60015f90610fb0565b90565b34610ffd57610fdd36600461029d565b610ff9610fe8610fbe565b610ff06101e2565b918291826102eb565b0390f35b6101e8565b90565b61101961101461101e92611002565b610cce565b6104d8565b90565b61102b6001611005565b90565b611036611021565b90565b919061104c905f60208501940190610abe565b565b3461107e5761105e36600461029d565b61107a61106961102e565b6110716101e2565b91829182611039565b0390f35b6101e8565b91909160c0818403126110fe5761109c835f830161043a565b926110aa81602084016104f9565b926110b88260408501610516565b9260608101359167ffffffffffffffff83116110f9576110dd846110f6948401610531565b9390946110ed8160808601610587565b9360a001610587565b90565b6101f0565b6101ec565b3461113a57611136611125611119366004611083565b95949094939193612f91565b61112d6101e2565b918291826102eb565b0390f35b6101e8565b9060208282031261115857611155915f01610587565b90565b6101ec565b3461118b5761117561117036600461113f565b61307b565b61117d6101e2565b8061118781610476565b0390f35b6101e8565b6111a09060086111a59302610c8a565b6109e3565b90565b906111b39154611190565b90565b6111c260095f906111a8565b90565b346111f5576111d536600461029d565b6111f16111e06111b6565b6111e86101e2565b91829182610c40565b0390f35b6101e8565b3461122a57611226611215611210366004610358565b613086565b61121d6101e2565b91829182610253565b0390f35b6101e8565b90565b5f1b90565b61124b6112466112509261122f565b611232565b6102db565b90565b61125c5f611237565b90565b611267611253565b90565b3461129a5761127a36600461029d565b61129661128561125f565b61128d6101e2565b918291826102eb565b0390f35b6101e8565b60018060a01b031690565b6112ba9060086112bf9302610c8a565b61129f565b90565b906112cd91546112aa565b90565b6112dc60085f906112c2565b90565b6112e890610ced565b90565b6112f4906112df565b9052565b919061130b905f602085019401906112eb565b565b3461133d5761131d36600461029d565b6113396113286112d0565b6113306101e2565b918291826112f8565b0390f35b6101e8565b67ffffffffffffffff811161135a5760208091020190565b6106a7565b67ffffffffffffffff811161137d5761137960209161069d565b0190565b6106a7565b909291926113976113928261135f565b6106e4565b938185526020850190828401116113b3576113b19261071c565b565b610699565b9080601f830112156113d6578160206113d393359101611382565b90565b610525565b9291906113ef6113ea82611342565b6106e4565b93818552602080860192028101918383116114465781905b838210611415575050505050565b813567ffffffffffffffff81116114415760209161143687849387016113b8565b815201910190611407565b610525565b61052d565b9080601f8301121561146957816020611466933591016113db565b90565b610525565b67ffffffffffffffff81116114865760208091020190565b6106a7565b92919061149f61149a8261146e565b6106e4565b93818552602080860192028101918383116114f65781905b8382106114c5575050505050565b813567ffffffffffffffff81116114f1576020916114e6878493870161075d565b8152019101906114b7565b610525565b61052d565b9080601f83011215611519578160206115169335910161148b565b90565b610525565b67ffffffffffffffff81116115365760208091020190565b6106a7565b5f80fd5b6115489061040f565b90565b6115548161153f565b0361155b57565b5f80fd5b9050359061156c8261154b565b565b91906040838203126115a8576115a19061158860406106e4565b93611595825f830161155f565b5f860152602001610587565b6020830152565b61153b565b909291926115c26115bd8261151e565b6106e4565b93818552604060208601920283019281841161160157915b8383106115e75750505050565b60206040916115f6848661156e565b8152019201916115da565b61052d565b9080601f8301121561162457816020611621933591016115ad565b90565b610525565b906101408282031261172257611641815f84016104f9565b9261164f8260208501610349565b9261165d8360408301610349565b9261166b8160608401610349565b92608083013567ffffffffffffffff811161171d578261168c91850161144b565b9260a081013567ffffffffffffffff811161171857836116ad9183016114fb565b9260c082013567ffffffffffffffff811161171357816116ce918401611606565b926116dc8260e08501610587565b926116eb836101008301610587565b9261012082013567ffffffffffffffff811161170e5761170b92016113b8565b90565b6101f0565b6101f0565b6101f0565b6101f0565b6101ec565b346117625761174c61173a366004611629565b98979097969196959295949394613c6e565b6117546101e2565b8061175e81610476565b0390f35b6101e8565b61177b6117766117809261122f565b610cce565b6104af565b90565b61178c5f611767565b90565b611797611783565b90565b91906117ad905f60208501940190610ab1565b565b346117df576117bf36600461029d565b6117db6117ca61178f565b6117d26101e2565b9182918261179a565b0390f35b6101e8565b9061188e906101008061184d61012084016118055f8801515f870190610d6f565b61181760208801516020870190610d7c565b61182960408801516040870190610d89565b61183b60608801516060870190610d96565b60808701518582036080870152610da3565b9461186060a082015160a0860190610dd4565b61187260c082015160c0860190610de1565b61188460e082015160e0860190610dd4565b0151910190610dd4565b90565b6040906118bd6118b26118c49597969460608401908482035f8601526117e4565b9660208301906102de565b0190610246565b565b346118f9576118d636600461029d565b6118f56118e1613c8d565b6118ec9391936101e2565b93849384611891565b0390f35b6101e8565b60018060a01b031690565b61191990600861191e9302610c8a565b6118fe565b90565b9061192c9154611909565b90565b61193b600a5f90611921565b90565b6119479061153f565b9052565b919061195e905f6020850194019061193e565b565b346119905761197036600461029d565b61198c61197b61192f565b6119836101e2565b9182918261194b565b0390f35b6101e8565b906020828203126119ae576119ab915f0161155f565b90565b6101ec565b346119e1576119cb6119c6366004611995565b613e2c565b6119d36101e2565b806119dd81610476565b0390f35b6101e8565b34611a15576119ff6119f9366004610449565b90613e61565b611a076101e2565b80611a1181610476565b0390f35b6101e8565b611a2660065f906111a8565b90565b34611a5957611a3936600461029d565b611a55611a44611a1a565b611a4c6101e2565b91829182610c40565b0390f35b6101e8565b5f80fd5b5f90565b611a6e611a62565b5080611a89611a83637965db0b60e01b6101f4565b916101f4565b14908115611a96575b5090565b611aa09150613e6d565b5f611a92565b5f90565b90611ab490610805565b5f5260205260405f2090565b6001611ad8611ade92611ad1611aa6565b505f611aaa565b016108cb565b90565b90611afc91611af7611af282611ac0565b613e93565b611afe565b565b90611b0891613ed6565b50565b90611b1591611ae1565b565b9695949392919080611b38611b32611b2d611783565b6104af565b916104af565b03611b4957611b4697611b65565b90565b5f635428eae760e01b815280611b6160048201610476565b0390fd5b9695949392919081611b86611b80611b7b611021565b6104d8565b916104d8565b03611b9757611b9497612286565b90565b5f6309ff9a8360e41b815280611baf60048201610476565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b611bd6611bdc91939293610570565b92610570565b8201809211611be757565b611bb3565b611bf8611bfd91610827565b61129f565b90565b611c0a9054611bec565b90565b60e01b90565b611c1c81610241565b03611c2357565b5f80fd5b90505190611c3482611c13565b565b90602082820312611c4f57611c4c915f01611c27565b90565b6101ec565b611c68611c63611c6d926104d8565b610cce565b610570565b90565b611c7990611c54565b9052565b916020611c9e929493611c9760408201965f830190611c70565b0190610b69565b565b611ca86101e2565b3d5f823e3d90fd5b5090565b90565b611ccb611cc6611cd092611cb4565b610cce565b610570565b90565b611cde6101206106e4565b90565b90611ceb906104af565b9052565b90611cf9906104d8565b9052565b90611d0790610aee565b9052565b90611d15906102db565b9052565b611d24913691611382565b90565b52565b90611d3490610570565b9052565b90611d429061041a565b9052565b634e487b7160e01b5f525f60045260245ffd5b611d6390516104af565b90565b90611d7260ff91611232565b9181191691161790565b611d90611d8b611d95926104af565b610cce565b6104af565b90565b90565b90611db0611dab611db792611d7c565b611d98565b8254611d66565b9055565b611dc590516104d8565b90565b60081b90565b90611de268ffffffffffffffff0091611dc8565b9181191691161790565b611e00611dfb611e05926104d8565b610cce565b6104d8565b90565b90565b90611e20611e1b611e2792611dec565b611e08565b8254611dce565b9055565b611e359051610aee565b90565b60481b90565b90611e5369ff00000000000000000091611e38565b9181191691161790565b611e6690610aee565b90565b90565b90611e81611e7c611e8892611e5d565b611e69565b8254611e3e565b9055565b611e9690516102db565b90565b90611ea55f1991611232565b9181191691161790565b611eb890610827565b90565b90611ed0611ecb611ed792610805565b611eaf565b8254611e99565b9055565b5190565b601f602091010490565b1b90565b91906008611f08910291611f025f1984611ee9565b92611ee9565b9181191691161790565b611f26611f21611f2b92610570565b610cce565b610570565b90565b90565b9190611f47611f42611f4f93611f12565b611f2e565b908354611eed565b9055565b5f90565b611f6991611f63611f53565b91611f31565b565b5b818110611f77575050565b80611f845f600193611f57565b01611f6c565b9190601f8111611f9a575b505050565b611fa6611fcb9361091f565b906020611fb284611edf565b83019310611fd3575b611fc490611edf565b0190611f6b565b5f8080611f95565b9150611fc481929050611fbb565b90611ff1905f1990600802610c8a565b191690565b8161200091611fe1565b906002021790565b9061201281610b13565b9067ffffffffffffffff82116120d2576120368261203085546108ec565b85611f8a565b602090601f831160011461206a57918091612059935f9261205e575b5050611ff6565b90555b565b90915001515f80612052565b601f198316916120798561091f565b925f5b8181106120ba575091600293918560019694106120a0575b5050500201905561205c565b6120b0910151601f841690611fe1565b90555f8080612094565b9193602060018192878701518155019501920161207c565b6106a7565b906120e191612008565b565b6120ed9051610570565b90565b9061210561210061210c92611f12565b611f2e565b8254611e99565b9055565b61211a905161041a565b90565b9061212e60018060a01b0391611232565b9181191691161790565b61214190610ced565b90565b90565b9061215c61215761216392612138565b612144565b825461211d565b9055565b9061224061010060066122469461218b5f82016121855f8801611d59565b90611d9b565b6121a35f820161219d60208801611dbb565b90611e0b565b6121bb5f82016121b560408801611e2b565b90611e6c565b6121d4600182016121ce60608801611e8c565b90611ebb565b6121ed600282016121e760808801611edb565b906120d7565b6122066003820161220060a088016120e3565b906120f0565b61221f6004820161221960c08801612110565b90612147565b6122386005820161223260e088016120e3565b906120f0565b0192016120e3565b906120f0565b565b9061225291612167565b565b9061225e90611f12565b5f5260205260405f2090565b61227390610570565b5f1981146122815760010190565b611bb3565b969391959294909650346122ac6122a66122a1868890611bc7565b610570565b91610570565b036125be57836122cd6122c76122c260096109fa565b610570565b91610570565b106125a2576122da612aad565b6122f56122ef6122ea60066109fa565b610570565b91610570565b1015612586578461230f6123096003610aee565b91610aee565b145f146124b457612321828290611cb0565b61233461232e6085611cb7565b91610570565b03612498575b338686848490879261234c60056109fa565b9461235696612f91565b969594918790919042933395969761236c611cd3565b995f8b019061237a91611ce1565b60208a019061238891611cef565b604089019061239691611cfd565b60608801906123a491611d0b565b6123ad91611d19565b60808601906123bb91611d27565b60a08501906123c991611d2a565b60c08401906123d791611d38565b60e08301906123e591611d2a565b6101008201906123f491611d2a565b60028261240091610811565b9061240a91612248565b80600361241760056109fa565b61242091612254565b9061242a91611ebb565b61243460056109fa565b61243d9061226a565b6124489060056120f0565b80337f73394f43049193ecd2e0c22eefa0ecf10987ce9c72da906a82a3336dc2ac05479161247590610805565b9061247f90612138565b916124886101e2565b8061249281610476565b0390a390565b5f637c6953f960e01b8152806124b060048201610476565b0390fd5b846124c86124c26002610aee565b91610aee565b146124d3575b61233a565b6124e56124e06008611c00565b6112df565b602063116ad4a691889061250b33946125166124ff6101e2565b96879586948594611c0d565b845260048401611c7d565b03915afa801561258157612532915f91612553575b5015610241565b156124ce575f63622a850b60e11b81528061254f60048201610476565b0390fd5b612574915060203d811161257a575b61256c81836106bb565b810190611c36565b5f61252b565b503d612562565b611ca0565b5f63d8219b3d60e01b81528061259e60048201610476565b0390fd5b5f63c77f97eb60e01b8152806125ba60048201610476565b0390fd5b5f632a9ffab760e21b8152806125d660048201610476565b0390fd5b906125f19695949392916125ec611aa6565b611b17565b90565b908061260f612609612604613f85565b61041a565b9161041a565b036126205761261d91613f92565b50565b5f63334bd91960e11b81528061263860048201610476565b0390fd5b90612657929161265261264d6102ac565b613e93565b612850565b565b61266861266e91939293610570565b92610570565b820391821161267957565b611bb3565b61269261268d6126979261122f565b610cce565b610570565b90565b6126a390610cd1565b90565b6126af9061269a565b90565b6126bb90610ced565b90565b905090565b6126ce5f80926126be565b0190565b6126db906126c3565b90565b906126f06126eb8361135f565b6106e4565b918252565b606090565b3d5f146127155761270a3d6126de565b903d5f602084013e5b565b61271d6126f5565b90612713565b91602061274492949361273d60408201965f830190610b69565b0190610b5c565b565b61275261275791610827565b6118fe565b90565b6127649054612746565b90565b6003111561277157565b610acb565b9061278082612767565b565b61278b90612776565b90565b61279790612782565b9052565b601111156127a557565b610acb565b906127b48261279b565b565b6127bf906127aa565b90565b6127cb906127b6565b9052565b5190565b60209181520190565b6127fb612804602093612809936127f2816127cf565b938480936127d3565b95869101610b20565b61069d565b0190565b90926128409061283661284d969461282c60808601975f870190610b5c565b602085019061278e565b60408301906127c2565b60608184039101526127dc565b90565b61286261285c82613086565b15610241565b612a845761287d600461287760028490610811565b01610a26565b6128d1612897600561289160028690610811565b016109fa565b6128cb6128b160066128ab60028890610811565b016109fa565b916128ba61427f565b916128c560096109fa565b90612659565b90611bc7565b90816128e56128df5f61267e565b91610570565b116129db575b50505f806129016128fc600a61275a565b6126b2565b61290b60096109fa565b6129136101e2565b908161291e816126d2565b03925af161292a6126fa565b505f14612985579061293c60096109fa565b61297f6001929461296d7f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198895610805565b956129766101e2565b9485948561280d565b0390a25b565b9061299060096109fa565b6129d3600292946129c17f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198895610805565b956129ca6101e2565b9485948561280d565b0390a2612983565b5f806129ee6129e9846126a6565b6126b2565b846129f76101e2565b9081612a02816126d2565b03925af1612a0e6126fa565b50612a19575b6128eb565b612a2f5f612a2960028690610811565b0161087a565b839192612a65612a5f7f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611dec565b93610805565b93612a7a612a716101e2565b92839283612723565b0390a35f80612a14565b5f6302e8145360e61b815280612a9c60048201610476565b0390fd5b90612aab929161263c565b565b612ab5611f53565b50612ac060056109fa565b612adb612ad5612ad060046109fa565b610570565b91610570565b115f14612b0257612aff612aef60056109fa565b612af960046109fa565b90612659565b90565b612b0b5f61267e565b90565b606090565b67ffffffffffffffff8111612b2b5760208091020190565b6106a7565b90612b42612b3d83612b13565b6106e4565b918252565b612b526101206106e4565b90565b5f90565b5f90565b5f90565b5f90565b606090565b5f90565b5f90565b612b7a612b47565b90602080808080808080808a612b8e612b55565b815201612b99612b59565b815201612ba4612b5d565b815201612baf612b61565b815201612bba612b65565b815201612bc5612b6a565b815201612bd0612b6e565b815201612bdb612b6a565b815201612be6612b6a565b81525050565b612bf4612b72565b90565b5f5b828110612c0557505050565b602090612c10612bec565b8184015201612bf9565b90612c3f612c2783612b30565b92602080612c358693612b13565b9201910390612bf7565b565b634e487b7160e01b5f52603260045260245ffd5b90612c5f82610d5c565b811015612c70576020809102010190565b612c41565b90612d53612d496006612c86612b47565b94612c9d612c955f8301610846565b5f8801611ce1565b612cb4612cab5f830161087a565b60208801611cef565b612ccb612cc25f83016108a7565b60408801611cfd565b612ce3612cda600183016108cb565b60608801611d0b565b612cfb612cf2600283016109c1565b60808801611d27565b612d13612d0a600383016109fa565b60a08801611d2a565b612d2b612d2260048301610a26565b60c08801611d38565b612d43612d3a600583016109fa565b60e08801611d2a565b016109fa565b6101008401611d2a565b565b612d5e90612c75565b90565b6001612d6d9101610570565b90565b612d78612b0e565b50612d89612d84612aad565b612c1a565b90612d9460046109fa565b91612d9d611f53565b925b80612dbb612db5612db060056109fa565b610570565b91610570565b1015612e1757612e0b612e1191612e04612de9612de2612ddd60038590612254565b6108cb565b6002610811565b85612df48992612d55565b612dfe8383612c55565b52612c55565b5150612d61565b93612d61565b92612d9f565b5090915090565b90612e2890612138565b5f5260205260405f2090565b60ff1690565b612e46612e4b91610827565b612e34565b90565b612e589054612e3a565b90565b612e81915f612e76612e7c93612e6f611a62565b5082611aaa565b01612e1e565b612e4e565b90565b60601b90565b612e9390612e84565b90565b612e9f90612e8a565b90565b612eae612eb39161041a565b612e96565b9052565b60c01b90565b612ec690612eb7565b90565b612ed5612eda916104d8565b612ebd565b9052565b60f81b90565b612eed90612ede565b90565b612efc612f0191610afa565b612ee4565b9052565b909182612f1581612f1c936126be565b809361071c565b0190565b90565b612f2f612f3491610570565b612f20565b9052565b93612f879560016020999895612f716008612f7f97612f6960148f9c612f6181612f789c612ea2565b018092612ec9565b018092612ef0565b0191612f05565b8092612f23565b018092612f23565b0190565b60200190565b94612fc29392969194612fd196612fa6611aa6565b5095979390919293612fb66101e2565b98899760208901612f38565b602082018103825203826106bb565b612fe3612fdd82610b13565b91612f8b565b2090565b61300090612ffb612ff66103ab565b613e93565b613002565b565b8061301561300f5f61267e565b91610570565b1461305f576130258160066120f0565b61305a7e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db1916130516101e2565b91829182610c40565b0390a1565b5f632a9ffab760e21b81528061307760048201610476565b0390fd5b61308490612fe7565b565b61308e611a62565b50613097612aad565b6130a96130a35f61267e565b91610570565b1190816130b5575b5090565b90506130e66130e06130da6130d560036130cf60046109fa565b90612254565b6108cb565b926102db565b916102db565b145f6130b1565b9897969594939291908961311061310a613105611021565b6104d8565b916104d8565b036131205761311e9961313c565b565b5f6309ff9a8360e41b81528061313860048201610476565b0390fd5b9061315e9998979695949392916131596131546102ac565b613e93565b613560565b565b5190565b5190565b61317461317991610827565b610c8e565b90565b6131869054613168565b90565b60209181520190565b60200190565b906131a291610da3565b90565b60200190565b906131bf6131b883613160565b8092613189565b90816131d060208302840194613192565b925f915b8383106131e357505050505090565b909192939460206132056131ff83856001950387528951613198565b976131a5565b93019301919392906131d4565b60209181520190565b60200190565b60209181520190565b61324961325260209361325793613240816127cf565b93848093613221565b95869101610b20565b61069d565b0190565b906132659161322a565b90565b60200190565b9061328261327b83613164565b8092613212565b90816132936020830284019461321b565b925f915b8383106132a657505050505090565b909192939460206132c86132c28385600195038752895161325b565b97613268565b9301930191939290613297565b5190565b60209181520190565b60200190565b6132f19061153f565b9052565b906020806133179361330d5f8201515f8601906132e8565b0151910190610dd4565b565b90613326816040936132f5565b0190565b60200190565b9061334d613347613340846132d5565b80936132d9565b926132e2565b905f5b81811061335d5750505090565b9091926133766133706001928651613319565b9461332a565b9101919091613350565b9691956133e46133f2926134249c9a9b976133d761340b986133cd6134009960408f6134169f6133c6906133bc61014084019a5f850190610abe565b60208301906102de565b01906102de565b60608d01906102de565b8a820360808c01526131ab565b9088820360a08a015261326e565b9086820360c0880152613330565b9560e0850190610b5c565b610100830190610b5c565b610120818403910152610b2b565b90565b90613431826132d5565b811015613442576020809102010190565b612c41565b61345090610ced565b90565b9061345d82613164565b81101561346e576020809102010190565b612c41565b9061347d82613160565b81101561348e576020809102010190565b612c41565b905090565b6134bd6134b4926020926134ab816127cf565b94858093613493565b93849101610b20565b0190565b6134ca91613498565b90565b6134e2906134d96101e2565b918291826134c1565b03902090565b6134fd9160208201915f818403910152610b2b565b90565b91602061352192949361351a60408201965f8301906102de565b01906102de565b565b61352d905161153f565b90565b613539906126b2565b9052565b91602061355e92949361355760408201965f830190613530565b0190610b5c565b565b969991989790959794929461357560016108cb565b61358f6135896135845f611237565b6102db565b916102db565b141580613c4a575b613c2e576135ad6135a787613086565b15610241565b613c12576135ba85613160565b6135d46135ce6135c987613164565b610570565b91610570565b03613bf657602087866136258b8f8f968f9089613630938f938e6136006135fb600761317c565b610cf9565b996364c062d1989b9d96909192939495966136196101e2565b9e8f9d8e9c8d9c611c0d565b8c5260048c01613380565b03915afa8015613bf15761364c915f91613bc3575b5015610241565b613ba75761366461365f60028790610811565b612d55565b998a61368f61368961368461010061367d868890611bc7565b94016120e3565b610570565b91610570565b03613b8b57816136b06136aa6136a560096109fa565b610570565b91610570565b10613b6f576136bd611f53565b9a8b9a8a6136c9611f53565b9d6136e06136db6136e6925b936132d5565b610570565b91610570565b101561372b578b9c8b9c8c906136fb91613427565b51602001613708906120e3565b61371191611bc7565b9c61371b90612d61565b9b6136e66136e06136db8f6136d5565b9193969995989b61374c919395989b50613746858790611bc7565b90611bc7565b61376761376161375b30613447565b31610570565b91610570565b11613b535761377789859061433d565b6137805f61267e565b5b8061379c6137966137918c613160565b610570565b91610570565b10156138205761381b908c8b6137b38b8490613453565b51906137c08d8590613473565b51926138136138016137fb6137f57f17e8258c8124324c2b7a7f6e8d90caec0a1373c6eeffcbb90f3c3220485d437194611dec565b94610805565b946134cd565b9461380a6101e2565b918291826134e8565b0390a4612d61565b613781565b509295509296955092966001613b03575b61383c826001611ebb565b8886919261387361386d7f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e93611dec565b93610805565b9361388861387f6101e2565b92839283613500565b0390a38161389e6138985f61267e565b91610570565b11613a9a575b61393f925f928392896138ba60c0899301612110565b926138ee6138e87f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611dec565b93610805565b936139036138fa6101e2565b92839283612723565b0390a3613918613913600a61275a565b6126b2565b906139216101e2565b908161392c816126d2565b03925af16139386126fa565b5015610241565b613a7e5761394c5f61267e565b5b8061396861396261395d866132d5565b610570565b91610570565b1015613a77576139ce5f8061399161398c82613985898890613427565b5101613523565b6126b2565b6139a860206139a1898890613427565b51016120e3565b6139b06101e2565b90816139bb816126d2565b03925af16139c76126fa565b5015610241565b613a5b57613a569085836139ee5f6139e7888690613427565b5101613523565b91613a0660206139ff898790613427565b51016120e3565b613a39613a337fccf5242d67663517e020b47096ecb69c7f1a32e85ebd295eae0611a41395b3ae93611dec565b93610805565b93613a4e613a456101e2565b9283928361353d565b0390a3612d61565b61394d565b5f6312171d8360e31b815280613a7360048201610476565b0390fd5b5050509050565b5f6312171d8360e31b815280613a9660048201610476565b0390fd5b613ae25f80613abb613ab6613ab160c08701612110565b6126a6565b6126b2565b85613ac46101e2565b9081613acf816126d2565b03925af1613adb6126fa565b5015610241565b156138a4575f6312171d8360e31b815280613aff60048201610476565b0390fd5b8886613b38613b327ff36f013d86705b333febcccf7ece97eee30f4cd398a9c69aee1f940ed5daddea93611dec565b91610805565b91613b416101e2565b80613b4b81610476565b0390a3613831565b5f631e9acf1760e31b815280613b6b60048201610476565b0390fd5b5f632a9ffab760e21b815280613b8760048201610476565b0390fd5b5f632a9ffab760e21b815280613ba360048201610476565b0390fd5b5f638baa579f60e01b815280613bbf60048201610476565b0390fd5b613be4915060203d8111613bea575b613bdc81836106bb565b810190611c36565b5f613645565b503d613bd2565b611ca0565b5f637c6953f960e01b815280613c0e60048201610476565b0390fd5b5f6302e8145360e61b815280613c2a60048201610476565b0390fd5b5f630b6fac0360e41b815280613c4660048201610476565b0390fd5b5086613c67613c61613c5c60016108cb565b6102db565b916102db565b1415613597565b90613c809998979695949392916130ed565b565b613c8a612b72565b90565b613c95613c82565b50613c9e611aa6565b50613ca7611a62565b50613cb0612aad565b613cc2613cbc5f61267e565b91610570565b11613ce157613ccf613c82565b613cd960016108cb565b915f91929190565b613d08613d01613cfc6003613cf660046109fa565b90612254565b6108cb565b6002610811565b613d1260016108cb565b91613d1e600192612d55565b929190565b613d3c90613d37613d326103ab565b613e93565b613daa565b565b613d52613d4d613d579261122f565b610cce565b61040f565b90565b613d6390613d3e565b90565b613d6f9061269a565b90565b90565b90613d8a613d85613d9192613d66565b613d72565b825461211d565b9055565b9190613da8905f60208501940190613530565b565b80613dc5613dbf613dba5f613d5a565b61041a565b916126b2565b14613e1057613dd581600a613d75565b613e0b7fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f91613e026101e2565b91829182613d95565b0390a1565b5f632582a64160e11b815280613e2860048201610476565b0390fd5b613e3590613d23565b565b90613e5291613e4d613e4882611ac0565b613e93565b613e54565b565b90613e5e91613f92565b50565b90613e6b91613e37565b565b613e75611a62565b50613e8f613e896301ffc9a760e01b6101f4565b916101f4565b1490565b613ea590613e9f613f85565b906143af565b565b613eb090610241565b90565b90565b90613ecb613ec6613ed292613ea7565b613eb3565b8254611d66565b9055565b613ede611a62565b50613ef3613eed828490612e5b565b15610241565b5f14613f7b57613f1a6001613f155f613f0d818690611aaa565b018590612e1e565b613eb6565b90613f23613f85565b90613f60613f5a613f547f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610805565b92612138565b92612138565b92613f696101e2565b80613f7381610476565b0390a4600190565b50505f90565b5f90565b613f8d613f81565b503390565b613f9a611a62565b50613fa6818390612e5b565b5f1461402d57613fcc5f613fc75f613fbf818690611aaa565b018590612e1e565b613eb6565b90613fd5613f85565b9061401261400c6140067ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b95610805565b92612138565b92612138565b9261401b6101e2565b8061402581610476565b0390a4600190565b50505f90565b919061404961404461405193610805565b611eaf565b908354611eed565b9055565b61406791614061611aa6565b91614033565b565b9061407c905f1990602003600802610c8a565b8154169055565b905f9161409a6140928261091f565b928354611ff6565b905555565b919290602082105f146140f857601f84116001146140c8576140c2929350611ff6565b90555b5b565b50906140ee6140f39360016140e56140df8561091f565b92611edf565b82019101611f6b565b614083565b6140c5565b5061412f829361410960019461091f565b61412861411585611edf565b820192601f86168061413a575b50611edf565b0190611f6b565b6002021790556140c6565b61414690888603614069565b5f614122565b9290916801000000000000000082116141ac576020115f1461419d57602081105f146141815761417b91611ff6565b90555b5b565b60019160ff19166141918461091f565b5560020201905561417e565b6001915060020201905561417f565b6106a7565b9081546141bd816108ec565b908183116141e6575b8183106141d4575b50505050565b6141dd9361409f565b5f8080806141ce565b6141f28383838761414c565b6141c6565b5f614201916141b1565b565b905f0361421557614213906141f7565b565b611d46565b5f60066142669282808201556142338360018301614055565b6142408360028301614203565b61424d8360038301611f57565b8260048201556142608360058301611f57565b01611f57565b565b905f0361427a576142789061421a565b565b611d46565b6142b05f6142ab60026142a56142a0600361429a60046109fa565b90612254565b6108cb565b90610811565b614268565b6142ce5f6142c960036142c360046109fa565b90612254565b614055565b6142ea6142e36142de60046109fa565b61226a565b60046120f0565b565b6142f75f80926127d3565b0190565b909161433a9361432361432d9261431960808601965f870190610b5c565b602085019061278e565b60408301906127c2565b60608183039101526142ec565b90565b61434561427f565b5f5f926143876143757f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198894610805565b9461437e6101e2565b938493846142fb565b0390a2565b9160206143ad9294936143a660408201965f830190610b69565b01906102de565b565b906143c46143be838390612e5b565b15610241565b6143cc575050565b6143e65f92839263e2517d3f60e01b84526004840161438c565b0390fdfea26469706673582212209b0a5d71454ccb7cf8aab226b530a08908ad49b6243afcedd702ff6f02a3056164736f6c634300081e0033",
}

// ProcessorEndpoint is an auto generated Go binding around an Ethereum contract.
type ProcessorEndpoint struct {
	abi abi.ABI
}

// NewProcessorEndpoint creates a new instance of ProcessorEndpoint.
func NewProcessorEndpoint() *ProcessorEndpoint {
	parsed, err := ProcessorEndpointMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ProcessorEndpoint{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ProcessorEndpoint) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _teeAuthenticator, address _authorityRegistry, address updateStatusOperator, address admin, uint256 _minFeePerRequest) returns()
func (processorEndpoint *ProcessorEndpoint) PackConstructor(_teeAuthenticator common.Address, _authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address, _minFeePerRequest *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("", _teeAuthenticator, _authorityRegistry, updateStatusOperator, admin, _minFeePerRequest)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackADMIN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a0acc6a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackADMIN() []byte {
	enc, err := processorEndpoint.abi.Pack("ADMIN")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackADMIN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a0acc6a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackADMIN() ([]byte, error) {
	return processorEndpoint.abi.Pack("ADMIN")
}

// UnpackADMIN is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a0acc6a.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackADMIN(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("ADMIN", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAPPLICATIONID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9968da96.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function APPLICATION_ID() view returns(uint64)
func (processorEndpoint *ProcessorEndpoint) PackAPPLICATIONID() []byte {
	enc, err := processorEndpoint.abi.Pack("APPLICATION_ID")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAPPLICATIONID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9968da96.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function APPLICATION_ID() view returns(uint64)
func (processorEndpoint *ProcessorEndpoint) TryPackAPPLICATIONID() ([]byte, error) {
	return processorEndpoint.abi.Pack("APPLICATION_ID")
}

// UnpackAPPLICATIONID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9968da96.
//
// Solidity: function APPLICATION_ID() view returns(uint64)
func (processorEndpoint *ProcessorEndpoint) UnpackAPPLICATIONID(data []byte) (uint64, error) {
	out, err := processorEndpoint.abi.Unpack("APPLICATION_ID", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackDEFAULTADMINROLE() []byte {
	enc, err := processorEndpoint.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return processorEndpoint.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackPROTOCOLVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa3aa460.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function PROTOCOL_VERSION() view returns(uint8)
func (processorEndpoint *ProcessorEndpoint) PackPROTOCOLVERSION() []byte {
	enc, err := processorEndpoint.abi.Pack("PROTOCOL_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPROTOCOLVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa3aa460.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function PROTOCOL_VERSION() view returns(uint8)
func (processorEndpoint *ProcessorEndpoint) TryPackPROTOCOLVERSION() ([]byte, error) {
	return processorEndpoint.abi.Pack("PROTOCOL_VERSION")
}

// UnpackPROTOCOLVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa3aa460.
//
// Solidity: function PROTOCOL_VERSION() view returns(uint8)
func (processorEndpoint *ProcessorEndpoint) UnpackPROTOCOLVERSION(data []byte) (uint8, error) {
	out, err := processorEndpoint.abi.Unpack("PROTOCOL_VERSION", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackUPDATESTATUSROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1664deeb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function UPDATE_STATUS_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackUPDATESTATUSROLE() []byte {
	enc, err := processorEndpoint.abi.Pack("UPDATE_STATUS_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUPDATESTATUSROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1664deeb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function UPDATE_STATUS_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackUPDATESTATUSROLE() ([]byte, error) {
	return processorEndpoint.abi.Pack("UPDATE_STATUS_ROLE")
}

// UnpackUPDATESTATUSROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1664deeb.
//
// Solidity: function UPDATE_STATUS_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackUPDATESTATUSROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("UPDATE_STATUS_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAuthorityRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ba045e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function authorityRegistry() view returns(address)
func (processorEndpoint *ProcessorEndpoint) PackAuthorityRegistry() []byte {
	enc, err := processorEndpoint.abi.Pack("authorityRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAuthorityRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ba045e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function authorityRegistry() view returns(address)
func (processorEndpoint *ProcessorEndpoint) TryPackAuthorityRegistry() ([]byte, error) {
	return processorEndpoint.abi.Pack("authorityRegistry")
}

// UnpackAuthorityRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9ba045e.
//
// Solidity: function authorityRegistry() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackAuthorityRegistry(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("authorityRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc415b95c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function feeCollector() view returns(address)
func (processorEndpoint *ProcessorEndpoint) PackFeeCollector() []byte {
	enc, err := processorEndpoint.abi.Pack("feeCollector")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc415b95c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function feeCollector() view returns(address)
func (processorEndpoint *ProcessorEndpoint) TryPackFeeCollector() ([]byte, error) {
	return processorEndpoint.abi.Pack("feeCollector")
}

// UnpackFeeCollector is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackFeeCollector(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("feeCollector", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a17995a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, idx *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, depositAmount, idx)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a17995a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, idx *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, depositAmount, idx)
}

// UnpackGenerateRequestId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a17995a.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackGenerateRequestId(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("generateRequestId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetNextPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc404c359.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getNextPendingRequest() view returns((uint8,uint64,uint8,bytes32,bytes,uint256,address,uint256,uint256), bytes32, bool success)
func (processorEndpoint *ProcessorEndpoint) PackGetNextPendingRequest() []byte {
	enc, err := processorEndpoint.abi.Pack("getNextPendingRequest")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetNextPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc404c359.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getNextPendingRequest() view returns((uint8,uint64,uint8,bytes32,bytes,uint256,address,uint256,uint256), bytes32, bool success)
func (processorEndpoint *ProcessorEndpoint) TryPackGetNextPendingRequest() ([]byte, error) {
	return processorEndpoint.abi.Pack("getNextPendingRequest")
}

// GetNextPendingRequestOutput serves as a container for the return parameters of contract
// method GetNextPendingRequest.
type GetNextPendingRequestOutput struct {
	Arg0    StructsPendingRequest
	Arg1    [32]byte
	Success bool
}

// UnpackGetNextPendingRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc404c359.
//
// Solidity: function getNextPendingRequest() view returns((uint8,uint64,uint8,bytes32,bytes,uint256,address,uint256,uint256), bytes32, bool success)
func (processorEndpoint *ProcessorEndpoint) UnpackGetNextPendingRequest(data []byte) (GetNextPendingRequestOutput, error) {
	out, err := processorEndpoint.abi.Unpack("getNextPendingRequest", data)
	outstruct := new(GetNextPendingRequestOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new(StructsPendingRequest)).(*StructsPendingRequest)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Success = *abi.ConvertType(out[2], new(bool)).(*bool)
	return *outstruct, nil
}

// PackGetPendingRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a1f712.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingRequests() view returns((uint8,uint64,uint8,bytes32,bytes,uint256,address,uint256,uint256)[])
func (processorEndpoint *ProcessorEndpoint) PackGetPendingRequests() []byte {
	enc, err := processorEndpoint.abi.Pack("getPendingRequests")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a1f712.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingRequests() view returns((uint8,uint64,uint8,bytes32,bytes,uint256,address,uint256,uint256)[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequests() ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequests")
}

// UnpackGetPendingRequests is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80a1f712.
//
// Solidity: function getPendingRequests() view returns((uint8,uint64,uint8,bytes32,bytes,uint256,address,uint256,uint256)[])
func (processorEndpoint *ProcessorEndpoint) UnpackGetPendingRequests(data []byte) ([]StructsPendingRequest, error) {
	out, err := processorEndpoint.abi.Unpack("getPendingRequests", data)
	if err != nil {
		return *new([]StructsPendingRequest), err
	}
	out0 := *abi.ConvertType(out[0], new([]StructsPendingRequest)).(*[]StructsPendingRequest)
	return out0, nil
}

// PackGetPendingRequestsSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e339384.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingRequestsSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackGetPendingRequestsSize() []byte {
	enc, err := processorEndpoint.abi.Pack("getPendingRequestsSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingRequestsSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e339384.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingRequestsSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequestsSize() ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequestsSize")
}

// UnpackGetPendingRequestsSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e339384.
//
// Solidity: function getPendingRequestsSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackGetPendingRequestsSize(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("getPendingRequestsSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (processorEndpoint *ProcessorEndpoint) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) UnpackHasRole(data []byte) (bool, error) {
	out, err := processorEndpoint.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsCurrentPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9fabc36f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) PackIsCurrentPendingRequest(requestId [32]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("isCurrentPendingRequest", requestId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsCurrentPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9fabc36f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) TryPackIsCurrentPendingRequest(requestId [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("isCurrentPendingRequest", requestId)
}

// UnpackIsCurrentPendingRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9fabc36f.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) UnpackIsCurrentPendingRequest(data []byte) (bool, error) {
	out, err := processorEndpoint.abi.Unpack("isCurrentPendingRequest", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMarkRequestFailed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b473cb9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function markRequestFailed(bytes32 requestId, uint8 errorCode, string errorMessage) returns()
func (processorEndpoint *ProcessorEndpoint) PackMarkRequestFailed(requestId [32]byte, errorCode uint8, errorMessage string) []byte {
	enc, err := processorEndpoint.abi.Pack("markRequestFailed", requestId, errorCode, errorMessage)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMarkRequestFailed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b473cb9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function markRequestFailed(bytes32 requestId, uint8 errorCode, string errorMessage) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackMarkRequestFailed(requestId [32]byte, errorCode uint8, errorMessage string) ([]byte, error) {
	return processorEndpoint.abi.Pack("markRequestFailed", requestId, errorCode, errorMessage)
}

// PackMaxQueueSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7d9cc1e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxQueueSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackMaxQueueSize() []byte {
	enc, err := processorEndpoint.abi.Pack("maxQueueSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxQueueSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7d9cc1e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxQueueSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) TryPackMaxQueueSize() ([]byte, error) {
	return processorEndpoint.abi.Pack("maxQueueSize")
}

// UnpackMaxQueueSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf7d9cc1e.
//
// Solidity: function maxQueueSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackMaxQueueSize(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("maxQueueSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMinFeePerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c8d416f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function minFeePerRequest() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackMinFeePerRequest() []byte {
	enc, err := processorEndpoint.abi.Pack("minFeePerRequest")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMinFeePerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c8d416f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function minFeePerRequest() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) TryPackMinFeePerRequest() ([]byte, error) {
	return processorEndpoint.abi.Pack("minFeePerRequest")
}

// UnpackMinFeePerRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9c8d416f.
//
// Solidity: function minFeePerRequest() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackMinFeePerRequest(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("minFeePerRequest", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (processorEndpoint *ProcessorEndpoint) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("renounceRole", role, callerConfirmation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRequestById is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5423112f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestById(bytes32 ) view returns(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes32 requestId, bytes payload, uint256 timestamp, address sender, uint256 depositAmount, uint256 maxFeeValue)
func (processorEndpoint *ProcessorEndpoint) PackRequestById(arg0 [32]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("requestById", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestById is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5423112f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestById(bytes32 ) view returns(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes32 requestId, bytes payload, uint256 timestamp, address sender, uint256 depositAmount, uint256 maxFeeValue)
func (processorEndpoint *ProcessorEndpoint) TryPackRequestById(arg0 [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("requestById", arg0)
}

// RequestByIdOutput serves as a container for the return parameters of contract
// method RequestById.
type RequestByIdOutput struct {
	ProtocolVersion uint8
	ApplicationId   uint64
	RequestType     uint8
	RequestId       [32]byte
	Payload         []byte
	Timestamp       *big.Int
	Sender          common.Address
	DepositAmount   *big.Int
	MaxFeeValue     *big.Int
}

// UnpackRequestById is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5423112f.
//
// Solidity: function requestById(bytes32 ) view returns(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes32 requestId, bytes payload, uint256 timestamp, address sender, uint256 depositAmount, uint256 maxFeeValue)
func (processorEndpoint *ProcessorEndpoint) UnpackRequestById(data []byte) (RequestByIdOutput, error) {
	out, err := processorEndpoint.abi.Unpack("requestById", data)
	outstruct := new(RequestByIdOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.ProtocolVersion = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.ApplicationId = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.RequestType = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.RequestId = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Payload = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.Timestamp = abi.ConvertType(out[5], new(big.Int)).(*big.Int)
	outstruct.Sender = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)
	outstruct.DepositAmount = abi.ConvertType(out[7], new(big.Int)).(*big.Int)
	outstruct.MaxFeeValue = abi.ConvertType(out[8], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (processorEndpoint *ProcessorEndpoint) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("revokeRole", role, account)
}

// PackStateRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9588eca2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateRoot() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackStateRoot() []byte {
	enc, err := processorEndpoint.abi.Pack("stateRoot")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9588eca2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateRoot() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackStateRoot() ([]byte, error) {
	return processorEndpoint.abi.Pack("stateRoot")
}

// UnpackStateRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackStateRoot(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("stateRoot", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ef290e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) PackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refund, applicationFees, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ef290e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refund, applicationFees, signature)
}

// PackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34f23a9d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, maxFeeValue *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, depositAmount, maxFeeValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34f23a9d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, maxFeeValue *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, depositAmount, maxFeeValue)
}

// UnpackSubmitRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x34f23a9d.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackSubmitRequest(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("submitRequest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := processorEndpoint.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTeeAuthenticator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e91d133.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function teeAuthenticator() view returns(address)
func (processorEndpoint *ProcessorEndpoint) PackTeeAuthenticator() []byte {
	enc, err := processorEndpoint.abi.Pack("teeAuthenticator")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTeeAuthenticator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e91d133.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function teeAuthenticator() view returns(address)
func (processorEndpoint *ProcessorEndpoint) TryPackTeeAuthenticator() ([]byte, error) {
	return processorEndpoint.abi.Pack("teeAuthenticator")
}

// UnpackTeeAuthenticator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e91d133.
//
// Solidity: function teeAuthenticator() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackTeeAuthenticator(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("teeAuthenticator", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackUpdateFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2c35ce8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (processorEndpoint *ProcessorEndpoint) PackUpdateFeeCollector(newFeeCollector common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("updateFeeCollector", newFeeCollector)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2c35ce8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackUpdateFeeCollector(newFeeCollector common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("updateFeeCollector", newFeeCollector)
}

// PackUpdateQueueThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c73eeaf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateQueueThreshold(uint256 newThreshold) returns()
func (processorEndpoint *ProcessorEndpoint) PackUpdateQueueThreshold(newThreshold *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("updateQueueThreshold", newThreshold)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateQueueThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c73eeaf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateQueueThreshold(uint256 newThreshold) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackUpdateQueueThreshold(newThreshold *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("updateQueueThreshold", newThreshold)
}

// ProcessorEndpointFeeCollectorUpdated represents a FeeCollectorUpdated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointFeeCollectorUpdated struct {
	NewFeeCollector common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointFeeCollectorUpdatedEventName = "FeeCollectorUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointFeeCollectorUpdated) ContractEventName() string {
	return ProcessorEndpointFeeCollectorUpdatedEventName
}

// UnpackFeeCollectorUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FeeCollectorUpdated(address newFeeCollector)
func (processorEndpoint *ProcessorEndpoint) UnpackFeeCollectorUpdatedEvent(log *types.Log) (*ProcessorEndpointFeeCollectorUpdated, error) {
	event := "FeeCollectorUpdated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointFeeCollectorUpdated)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointQueueThresholdUpdated represents a QueueThresholdUpdated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointQueueThresholdUpdated struct {
	NewThreshold *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointQueueThresholdUpdatedEventName = "QueueThresholdUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointQueueThresholdUpdated) ContractEventName() string {
	return ProcessorEndpointQueueThresholdUpdatedEventName
}

// UnpackQueueThresholdUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event QueueThresholdUpdated(uint256 newThreshold)
func (processorEndpoint *ProcessorEndpoint) UnpackQueueThresholdUpdatedEvent(log *types.Log) (*ProcessorEndpointQueueThresholdUpdated, error) {
	event := "QueueThresholdUpdated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointQueueThresholdUpdated)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointRefund represents a Refund event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRefund struct {
	ApplicationId uint64
	RequestId     [32]byte
	To            common.Address
	Amount        *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRefundEventName = "Refund"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRefund) ContractEventName() string {
	return ProcessorEndpointRefundEventName
}

// UnpackRefundEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Refund(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount)
func (processorEndpoint *ProcessorEndpoint) UnpackRefundEvent(log *types.Log) (*ProcessorEndpointRefund, error) {
	event := "Refund"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRefund)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointReportGenerated represents a ReportGenerated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointReportGenerated struct {
	ApplicationId uint64
	RequestId     [32]byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointReportGeneratedEventName = "ReportGenerated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointReportGenerated) ContractEventName() string {
	return ProcessorEndpointReportGeneratedEventName
}

// UnpackReportGeneratedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ReportGenerated(uint64 indexed applicationId, bytes32 indexed requestId)
func (processorEndpoint *ProcessorEndpoint) UnpackReportGeneratedEvent(log *types.Log) (*ProcessorEndpointReportGenerated, error) {
	event := "ReportGenerated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointReportGenerated)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointRequestCompleted represents a RequestCompleted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRequestCompleted struct {
	RequestId       [32]byte
	ApplicationFees *big.Int
	Status          uint8
	ErrorCode       uint8
	ErrorMessage    string
	Raw             *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRequestCompletedEventName = "RequestCompleted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRequestCompleted) ContractEventName() string {
	return ProcessorEndpointRequestCompletedEventName
}

// UnpackRequestCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestCompleted(bytes32 indexed requestId, uint256 applicationFees, uint8 status, uint8 errorCode, string errorMessage)
func (processorEndpoint *ProcessorEndpoint) UnpackRequestCompletedEvent(log *types.Log) (*ProcessorEndpointRequestCompleted, error) {
	event := "RequestCompleted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRequestCompleted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointRequestSubmitted represents a RequestSubmitted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRequestSubmitted struct {
	RequestId [32]byte
	Sender    common.Address
	Raw       *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRequestSubmittedEventName = "RequestSubmitted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRequestSubmitted) ContractEventName() string {
	return ProcessorEndpointRequestSubmittedEventName
}

// UnpackRequestSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestSubmitted(bytes32 indexed requestId, address indexed sender)
func (processorEndpoint *ProcessorEndpoint) UnpackRequestSubmittedEvent(log *types.Log) (*ProcessorEndpointRequestSubmitted, error) {
	event := "RequestSubmitted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRequestSubmitted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointRoleAdminChanged represents a RoleAdminChanged event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRoleAdminChanged) ContractEventName() string {
	return ProcessorEndpointRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (processorEndpoint *ProcessorEndpoint) UnpackRoleAdminChangedEvent(log *types.Log) (*ProcessorEndpointRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointRoleGranted represents a RoleGranted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRoleGranted) ContractEventName() string {
	return ProcessorEndpointRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (processorEndpoint *ProcessorEndpoint) UnpackRoleGrantedEvent(log *types.Log) (*ProcessorEndpointRoleGranted, error) {
	event := "RoleGranted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRoleGranted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointRoleRevoked represents a RoleRevoked event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRoleRevoked) ContractEventName() string {
	return ProcessorEndpointRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (processorEndpoint *ProcessorEndpoint) UnpackRoleRevokedEvent(log *types.Log) (*ProcessorEndpointRoleRevoked, error) {
	event := "RoleRevoked"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRoleRevoked)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointStateRootUpdate represents a StateRootUpdate event raised by the ProcessorEndpoint contract.
type ProcessorEndpointStateRootUpdate struct {
	ApplicationId uint64
	RequestId     [32]byte
	OldStateRoot  [32]byte
	NewStateRoot  [32]byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointStateRootUpdateEventName = "StateRootUpdate"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointStateRootUpdate) ContractEventName() string {
	return ProcessorEndpointStateRootUpdateEventName
}

// UnpackStateRootUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StateRootUpdate(uint64 indexed applicationId, bytes32 indexed requestId, bytes32 oldStateRoot, bytes32 newStateRoot)
func (processorEndpoint *ProcessorEndpoint) UnpackStateRootUpdateEvent(log *types.Log) (*ProcessorEndpointStateRootUpdate, error) {
	event := "StateRootUpdate"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointStateRootUpdate)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointUserEvent represents a UserEvent event raised by the ProcessorEndpoint contract.
type ProcessorEndpointUserEvent struct {
	ApplicationId uint64
	RequestId     [32]byte
	EventSubType  common.Hash
	EncryptedData []byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointUserEventEventName = "UserEvent"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointUserEvent) ContractEventName() string {
	return ProcessorEndpointUserEventEventName
}

// UnpackUserEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UserEvent(uint64 indexed applicationId, bytes32 indexed requestId, string indexed eventSubType, bytes encryptedData)
func (processorEndpoint *ProcessorEndpoint) UnpackUserEventEvent(log *types.Log) (*ProcessorEndpointUserEvent, error) {
	event := "UserEvent"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointUserEvent)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ProcessorEndpointWithdrawal represents a Withdrawal event raised by the ProcessorEndpoint contract.
type ProcessorEndpointWithdrawal struct {
	ApplicationId uint64
	RequestId     [32]byte
	To            common.Address
	Amount        *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointWithdrawalEventName = "Withdrawal"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointWithdrawal) ContractEventName() string {
	return ProcessorEndpointWithdrawalEventName
}

// UnpackWithdrawalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdrawal(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount)
func (processorEndpoint *ProcessorEndpoint) UnpackWithdrawalEvent(log *types.Log) (*ProcessorEndpointWithdrawal, error) {
	event := "Withdrawal"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointWithdrawal)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (processorEndpoint *ProcessorEndpoint) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AddressCantBeZero"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AuthorityNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAuthorityNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["FeeValueBelowMinimum"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackFeeValueBelowMinimumError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InsufficientBalance"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidApplicationId"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidApplicationIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidPayload"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidPayloadError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidProtocolVersion"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidProtocolVersionError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidRequestId"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidRequestIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidStateRoot"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidStateRootError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidValue"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidValueError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["QueueThresholdExceeded"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackQueueThresholdExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["TransferFailed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackTransferFailedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ProcessorEndpointAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func ProcessorEndpointAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (processorEndpoint *ProcessorEndpoint) UnpackAccessControlBadConfirmationError(raw []byte) (*ProcessorEndpointAccessControlBadConfirmation, error) {
	out := new(ProcessorEndpointAccessControlBadConfirmation)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func ProcessorEndpointAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (processorEndpoint *ProcessorEndpoint) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*ProcessorEndpointAccessControlUnauthorizedAccount, error) {
	out := new(ProcessorEndpointAccessControlUnauthorizedAccount)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointAddressCantBeZero represents a AddressCantBeZero error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressCantBeZero()
func ProcessorEndpointAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x4b054c82bb9e4ebe90744902747c4e86028dd2a122c63de8810d9a1c2e84617d")
}

// UnpackAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressCantBeZero()
func (processorEndpoint *ProcessorEndpoint) UnpackAddressCantBeZeroError(raw []byte) (*ProcessorEndpointAddressCantBeZero, error) {
	out := new(ProcessorEndpointAddressCantBeZero)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointAuthorityNotAllowed represents a AuthorityNotAllowed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAuthorityNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuthorityNotAllowed()
func ProcessorEndpointAuthorityNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xc4550a169e74f716a0e1f2ed847c9b836363900678799385b1658d47630ac161")
}

// UnpackAuthorityNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuthorityNotAllowed()
func (processorEndpoint *ProcessorEndpoint) UnpackAuthorityNotAllowedError(raw []byte) (*ProcessorEndpointAuthorityNotAllowed, error) {
	out := new(ProcessorEndpointAuthorityNotAllowed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AuthorityNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointFeeValueBelowMinimum represents a FeeValueBelowMinimum error raised by the ProcessorEndpoint contract.
type ProcessorEndpointFeeValueBelowMinimum struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FeeValueBelowMinimum()
func ProcessorEndpointFeeValueBelowMinimumErrorID() common.Hash {
	return common.HexToHash("0xc77f97ebda9854b3f6367f2d96830b6ea60a9ab7e12c67ede499fadeb3f9cef7")
}

// UnpackFeeValueBelowMinimumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FeeValueBelowMinimum()
func (processorEndpoint *ProcessorEndpoint) UnpackFeeValueBelowMinimumError(raw []byte) (*ProcessorEndpointFeeValueBelowMinimum, error) {
	out := new(ProcessorEndpointFeeValueBelowMinimum)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "FeeValueBelowMinimum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInsufficientBalance represents a InsufficientBalance error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInsufficientBalance struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientBalance()
func ProcessorEndpointInsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xf4d678b8ce6b5157126b1484a53523762a93571537a7d5ae97d8014a44715c94")
}

// UnpackInsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientBalance()
func (processorEndpoint *ProcessorEndpoint) UnpackInsufficientBalanceError(raw []byte) (*ProcessorEndpointInsufficientBalance, error) {
	out := new(ProcessorEndpointInsufficientBalance)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidApplicationId represents a InvalidApplicationId error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidApplicationId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidApplicationId()
func ProcessorEndpointInvalidApplicationIdErrorID() common.Hash {
	return common.HexToHash("0x9ff9a83079daf92a88066c1b25d04540f443a8d245022db6e1e62beaa6d6bc5e")
}

// UnpackInvalidApplicationIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidApplicationId()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidApplicationIdError(raw []byte) (*ProcessorEndpointInvalidApplicationId, error) {
	out := new(ProcessorEndpointInvalidApplicationId)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidApplicationId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidPayload represents a InvalidPayload error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidPayload struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPayload()
func ProcessorEndpointInvalidPayloadErrorID() common.Hash {
	return common.HexToHash("0x7c6953f9f55fcfd713fe1d72e370e4757e17ae7187e1acdcc02b4b76b56a9a35")
}

// UnpackInvalidPayloadError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPayload()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidPayloadError(raw []byte) (*ProcessorEndpointInvalidPayload, error) {
	out := new(ProcessorEndpointInvalidPayload)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidPayload", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidProtocolVersion represents a InvalidProtocolVersion error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidProtocolVersion struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidProtocolVersion()
func ProcessorEndpointInvalidProtocolVersionErrorID() common.Hash {
	return common.HexToHash("0x5428eae7f7dcfe744a90b56d360bbe512220bcf69b1e537f3d2d3b28cf2bf92d")
}

// UnpackInvalidProtocolVersionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidProtocolVersion()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidProtocolVersionError(raw []byte) (*ProcessorEndpointInvalidProtocolVersion, error) {
	out := new(ProcessorEndpointInvalidProtocolVersion)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidProtocolVersion", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidRequestId represents a InvalidRequestId error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidRequestId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRequestId()
func ProcessorEndpointInvalidRequestIdErrorID() common.Hash {
	return common.HexToHash("0xba0514c0d5ec80c22a6ebd3f2e6691e2f9ad3d1402978084361a7537d8546138")
}

// UnpackInvalidRequestIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRequestId()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidRequestIdError(raw []byte) (*ProcessorEndpointInvalidRequestId, error) {
	out := new(ProcessorEndpointInvalidRequestId)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidRequestId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidSignature represents a InvalidSignature error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature()
func ProcessorEndpointInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x8baa579fce362245063d36f11747a89dd489c54795634fc673cc0e0db51fedc5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidSignatureError(raw []byte) (*ProcessorEndpointInvalidSignature, error) {
	out := new(ProcessorEndpointInvalidSignature)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidStateRoot represents a InvalidStateRoot error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidStateRoot struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidStateRoot()
func ProcessorEndpointInvalidStateRootErrorID() common.Hash {
	return common.HexToHash("0xb6fac0303f809aaca0354efa81577129a51cc71ce625abce5d9ded4cf63ebabc")
}

// UnpackInvalidStateRootError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidStateRoot()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidStateRootError(raw []byte) (*ProcessorEndpointInvalidStateRoot, error) {
	out := new(ProcessorEndpointInvalidStateRoot)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidStateRoot", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidValue represents a InvalidValue error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidValue struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidValue()
func ProcessorEndpointInvalidValueErrorID() common.Hash {
	return common.HexToHash("0xaa7feadc1a228aee3d4475b95bb4402ebafd3841876f41c4469039f5ee2998d5")
}

// UnpackInvalidValueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidValue()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidValueError(raw []byte) (*ProcessorEndpointInvalidValue, error) {
	out := new(ProcessorEndpointInvalidValue)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidValue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointQueueThresholdExceeded represents a QueueThresholdExceeded error raised by the ProcessorEndpoint contract.
type ProcessorEndpointQueueThresholdExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QueueThresholdExceeded()
func ProcessorEndpointQueueThresholdExceededErrorID() common.Hash {
	return common.HexToHash("0xd8219b3da97d7adbef217ce9f220bbf38600f597fb984fd0620ea7df76ea17f3")
}

// UnpackQueueThresholdExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QueueThresholdExceeded()
func (processorEndpoint *ProcessorEndpoint) UnpackQueueThresholdExceededError(raw []byte) (*ProcessorEndpointQueueThresholdExceeded, error) {
	out := new(ProcessorEndpointQueueThresholdExceeded)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "QueueThresholdExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointTransferFailed represents a TransferFailed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointTransferFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TransferFailed()
func ProcessorEndpointTransferFailedErrorID() common.Hash {
	return common.HexToHash("0x90b8ec1877afffd816d05d9b13947f3ff18ec5851c38bad15ec2b710f92391b1")
}

// UnpackTransferFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TransferFailed()
func (processorEndpoint *ProcessorEndpoint) UnpackTransferFailedError(raw []byte) (*ProcessorEndpointTransferFailed, error) {
	out := new(ProcessorEndpointTransferFailed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "TransferFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StructsMetaData contains all meta data concerning the Structs contract.
var StructsMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "31c70fec60e1e1f4cc22f993498ea8b973",
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea264697066735822122019edc9d22300cf288f695a96c7c49fa26dc88354373ec7048cdc1f4643c08ec464736f6c634300081e0033",
}

// Structs is an auto generated Go binding around an Ethereum contract.
type Structs struct {
	abi abi.ABI
}

// NewStructs creates a new instance of Structs.
func NewStructs() *Structs {
	parsed, err := StructsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Structs{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Structs) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}
