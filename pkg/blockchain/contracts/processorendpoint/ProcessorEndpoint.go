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
	Timestamp       *big.Int
	DepositAmount   *big.Int
	MaxFeeValue     *big.Int
	RequestId       [32]byte
	Payload         []byte
	Sender          common.Address
	ApplicationId   uint64
	ProtocolVersion uint8
	RequestType     uint8
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
	Bin: "0x6080604052346100305761001a610014610104565b906101f4565b610022610035565b610a076103c18239610a0790f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b60018060a01b031690565b6100b19061009d565b90565b6100bd816100a8565b036100c457565b5f80fd5b905051906100d5826100b4565b565b91906040838203126100ff57806100f36100fc925f86016100c8565b936020016100c8565b90565b610099565b610122610dc88038038061011781610084565b9283398101906100d7565b9091565b90565b90565b61014061013b61014592610126565b610129565b61009d565b90565b6101519061012c565b90565b5f0190565b61016d6101686101729261009d565b610129565b61009d565b90565b61017e90610159565b90565b61018a90610175565b90565b5f1b90565b906101a360018060a01b039161018d565b9181191691161790565b6101b690610175565b90565b90565b906101d16101cc6101d8926101ad565b6101b9565b8254610192565b9055565b6101e590610159565b90565b6101f1906101dc565b90565b6101fd9061028e565b8061021861021261020d5f610148565b6100a8565b916100a8565b146102725761023061022982610181565b60026101bc565b61025a7fd4047f78dd943d75dfa94ddb7e36315e62df8a01f00a3d0f512960bf93c12318916101e8565b90610263610035565b8061026d81610154565b0390a2565b5f632582a64160e11b81528061028a60048201610154565b0390fd5b610297906102bb565b565b6102a2906100a8565b9052565b91906102b9905f60208501940190610299565b565b806102d66102d06102cb5f610148565b6100a8565b916100a8565b146102e6576102e490610361565b565b6103096102f25f610148565b5f918291631e4fbdf760e01b8352600483016102a6565b0390fd5b5f1c90565b60018060a01b031690565b61032961032e9161030d565b610312565b90565b61033b905461031d565b90565b90565b9061035661035161035d926101e8565b61033e565b8254610192565b9055565b61036a5f610331565b610374825f610341565b906103a86103a27f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936101e8565b916101e8565b916103b1610035565b806103bb81610154565b0390a356fe60806040526004361015610013575b610460565b61001d5f3561009c565b8063116ad4a6146100975780633aac66d1146100925780635d35bc141461008d5780636f5e32f214610088578063715018a6146100835780638da5cb5b1461007e578063cfa8904f146100795763f2fde38b0361000e5761042d565b6103f8565b6103b4565b61035f565b61031b565b6101f1565b61019f565b610164565b60e01c90565b60405190565b5f80fd5b5f80fd5b90565b6100bc816100b0565b036100c357565b5f80fd5b905035906100d4826100b3565b565b60018060a01b031690565b6100ea906100d6565b90565b6100f6816100e1565b036100fd57565b5f80fd5b9050359061010e826100ed565b565b9190604083820312610138578061012c610135925f86016100c7565b93602001610101565b90565b6100ac565b151590565b61014b9061013d565b9052565b9190610162905f60208501940190610142565b565b346101955761019161018061017a366004610110565b90610587565b6101886100a2565b9182918261014f565b0390f35b6100a8565b5f0190565b346101ce576101b86101b2366004610110565b90610755565b6101c06100a2565b806101ca8161019a565b0390f35b6100a8565b906020828203126101ec576101e9915f01610101565b90565b6100ac565b3461021f576102096102043660046101d3565b610805565b6102116100a2565b8061021b8161019a565b0390f35b6100a8565b9060208282031261023d5761023a915f016100c7565b90565b6100ac565b90565b61025961025461025e926100b0565b610242565b6100b0565b90565b9061026b90610245565b5f5260205260405f2090565b1c90565b60018060a01b031690565b61029690600861029b9302610277565b61027b565b90565b906102a99154610286565b90565b6102c2906102bd6001915f92610261565b61029e565b90565b6102d96102d46102de926100d6565b610242565b6100d6565b90565b6102ea906102c5565b90565b6102f6906102e1565b90565b610302906102ed565b9052565b9190610319905f602085019401906102f9565b565b3461034b57610347610336610331366004610224565b6102ac565b61033e6100a2565b91829182610306565b0390f35b6100a8565b5f91031261035a57565b6100ac565b3461038d5761036f366004610350565b610377610835565b61037f6100a2565b806103898161019a565b0390f35b6100a8565b61039b906100e1565b9052565b91906103b2905f60208501940190610392565b565b346103e4576103c4366004610350565b6103e06103cf61086f565b6103d76100a2565b9182918261039f565b0390f35b6100a8565b6103f560025f9061029e565b90565b3461042857610408366004610350565b6104246104136103e9565b61041b6100a2565b91829182610306565b0390f35b6100a8565b3461045b576104456104403660046101d3565b6108e9565b61044d6100a2565b806104578161019a565b0390f35b6100a8565b5f80fd5b5f90565b5f1c90565b61047961047e91610468565b61027b565b90565b61048b905461046d565b90565b90565b6104a56104a06104aa9261048e565b610242565b6100d6565b90565b6104b690610491565b90565b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906104e1906104b9565b810190811067ffffffffffffffff8211176104fb57604052565b6104c3565b60e01b90565b61050f8161013d565b0361051657565b5f80fd5b9050519061052782610506565b565b906020828203126105425761053f915f0161051a565b90565b6100ac565b610550906100b0565b9052565b91602061057592949361056e60408201965f830190610547565b0190610392565b565b61057f6100a2565b3d5f823e3d90fd5b90602090610593610464565b506105a86105a360018590610261565b610481565b6105b1816102ed565b6105cb6105c56105c05f6104ad565b6100e1565b916100e1565b14610650575b6105da906102ed565b6105fc63116ad4a69492946106076105f06100a2565b96879586948594610500565b845260048401610554565b03915afa90811561064b575f9161061d575b5090565b61063e915060203d8111610644575b61063681836104d7565b810190610529565b5f610619565b503d61062c565b610577565b506105da61065e6002610481565b90506105d1565b90610677916106726108f4565b6106ec565b565b610682906102c5565b90565b61068e90610679565b90565b5f1b90565b906106a760018060a01b0391610691565b9181191691161790565b6106ba90610679565b90565b90565b906106d56106d06106dc926106b1565b6106bd565b8254610696565b9055565b6106e9906102e1565b90565b6107096106f883610685565b61070460018490610261565b6106c0565b9061073d6107377f4f79164bfe5c672976dae578ab2d7b0c99754246947eb57d8fd2a5a078307c2293610245565b916106e0565b916107466100a2565b806107508161019a565b0390a3565b9061075f91610665565b565b6107729061076d6108f4565b610774565b565b8061078f6107896107845f6104ad565b6100e1565b916100e1565b146107e9576107a76107a082610685565b60026106c0565b6107d17fd4047f78dd943d75dfa94ddb7e36315e62df8a01f00a3d0f512960bf93c12318916106e0565b906107da6100a2565b806107e48161019a565b0390a2565b5f632582a64160e11b8152806108016004820161019a565b0390fd5b61080e90610761565b565b6108186108f4565b610820610822565b565b61083361082e5f6104ad565b610965565b565b61083d610810565b565b5f90565b60018060a01b031690565b61085a61085f91610468565b610843565b90565b61086c905461084e565b90565b61087761083f565b506108815f610862565b90565b610895906108906108f4565b610897565b565b806108b26108ac6108a75f6104ad565b6100e1565b916100e1565b146108c2576108c090610965565b565b6108e56108ce5f6104ad565b5f918291631e4fbdf760e01b83526004830161039f565b0390fd5b6108f290610884565b565b6108fc61086f565b61091561090f61090a6109c4565b6100e1565b916100e1565b0361091c57565b61093e6109276109c4565b5f91829163118cdaa760e01b83526004830161039f565b0390fd5b90565b9061095a610955610961926106e0565b610942565b8254610696565b9055565b61096e5f610862565b610978825f610945565b906109ac6109a67f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936106e0565b916106e0565b916109b56100a2565b806109bf8161019a565b0390a3565b6109cc61083f565b50339056fea26469706673582212207d166fd52007e3b30ffa1a69546a9e3f42ad1fd7554865f6458c47f869c3479064736f6c634300081e0033",
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

// IAuthorityRegistryMetaData contains all meta data concerning the IAuthorityRegistry contract.
var IAuthorityRegistryMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"AppAuthorityContractSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"DefaultAuthorityContractSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"checkAuthorityIsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"setAppAuthorityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"setDefaultAuthorityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "37d38b366c77c42b530df09daa51bb4a27",
}

// IAuthorityRegistry is an auto generated Go binding around an Ethereum contract.
type IAuthorityRegistry struct {
	abi abi.ABI
}

// NewIAuthorityRegistry creates a new instance of IAuthorityRegistry.
func NewIAuthorityRegistry() *IAuthorityRegistry {
	parsed, err := IAuthorityRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IAuthorityRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IAuthorityRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x116ad4a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (iAuthorityRegistry *IAuthorityRegistry) PackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) []byte {
	enc, err := iAuthorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
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
func (iAuthorityRegistry *IAuthorityRegistry) TryPackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) ([]byte, error) {
	return iAuthorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
}

// UnpackCheckAuthorityIsAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x116ad4a6.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (iAuthorityRegistry *IAuthorityRegistry) UnpackCheckAuthorityIsAllowed(data []byte) (bool, error) {
	out, err := iAuthorityRegistry.abi.Unpack("checkAuthorityIsAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackSetAppAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3aac66d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setAppAuthorityContract(uint256 applicationId, address authorityContract) returns()
func (iAuthorityRegistry *IAuthorityRegistry) PackSetAppAuthorityContract(applicationId *big.Int, authorityContract common.Address) []byte {
	enc, err := iAuthorityRegistry.abi.Pack("setAppAuthorityContract", applicationId, authorityContract)
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
func (iAuthorityRegistry *IAuthorityRegistry) TryPackSetAppAuthorityContract(applicationId *big.Int, authorityContract common.Address) ([]byte, error) {
	return iAuthorityRegistry.abi.Pack("setAppAuthorityContract", applicationId, authorityContract)
}

// PackSetDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d35bc14.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setDefaultAuthorityContract(address authorityContract) returns()
func (iAuthorityRegistry *IAuthorityRegistry) PackSetDefaultAuthorityContract(authorityContract common.Address) []byte {
	enc, err := iAuthorityRegistry.abi.Pack("setDefaultAuthorityContract", authorityContract)
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
func (iAuthorityRegistry *IAuthorityRegistry) TryPackSetDefaultAuthorityContract(authorityContract common.Address) ([]byte, error) {
	return iAuthorityRegistry.abi.Pack("setDefaultAuthorityContract", authorityContract)
}

// IAuthorityRegistryAppAuthorityContractSet represents a AppAuthorityContractSet event raised by the IAuthorityRegistry contract.
type IAuthorityRegistryAppAuthorityContractSet struct {
	ApplicationId     *big.Int
	AuthorityContract common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const IAuthorityRegistryAppAuthorityContractSetEventName = "AppAuthorityContractSet"

// ContractEventName returns the user-defined event name.
func (IAuthorityRegistryAppAuthorityContractSet) ContractEventName() string {
	return IAuthorityRegistryAppAuthorityContractSetEventName
}

// UnpackAppAuthorityContractSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AppAuthorityContractSet(uint256 indexed applicationId, address indexed authorityContract)
func (iAuthorityRegistry *IAuthorityRegistry) UnpackAppAuthorityContractSetEvent(log *types.Log) (*IAuthorityRegistryAppAuthorityContractSet, error) {
	event := "AppAuthorityContractSet"
	if log.Topics[0] != iAuthorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IAuthorityRegistryAppAuthorityContractSet)
	if len(log.Data) > 0 {
		if err := iAuthorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAuthorityRegistry.abi.Events[event].Inputs {
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

// IAuthorityRegistryDefaultAuthorityContractSet represents a DefaultAuthorityContractSet event raised by the IAuthorityRegistry contract.
type IAuthorityRegistryDefaultAuthorityContractSet struct {
	AuthorityContract common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const IAuthorityRegistryDefaultAuthorityContractSetEventName = "DefaultAuthorityContractSet"

// ContractEventName returns the user-defined event name.
func (IAuthorityRegistryDefaultAuthorityContractSet) ContractEventName() string {
	return IAuthorityRegistryDefaultAuthorityContractSetEventName
}

// UnpackDefaultAuthorityContractSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAuthorityContractSet(address indexed authorityContract)
func (iAuthorityRegistry *IAuthorityRegistry) UnpackDefaultAuthorityContractSetEvent(log *types.Log) (*IAuthorityRegistryDefaultAuthorityContractSet, error) {
	event := "DefaultAuthorityContractSet"
	if log.Topics[0] != iAuthorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IAuthorityRegistryDefaultAuthorityContractSet)
	if len(log.Data) > 0 {
		if err := iAuthorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAuthorityRegistry.abi.Events[event].Inputs {
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
func (iAuthorityRegistry *IAuthorityRegistry) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iAuthorityRegistry.abi.Errors["AddressCantBeZero"].ID.Bytes()[:4]) {
		return iAuthorityRegistry.UnpackAddressCantBeZeroError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IAuthorityRegistryAddressCantBeZero represents a AddressCantBeZero error raised by the IAuthorityRegistry contract.
type IAuthorityRegistryAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressCantBeZero()
func IAuthorityRegistryAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x4b054c82bb9e4ebe90744902747c4e86028dd2a122c63de8810d9a1c2e84617d")
}

// UnpackAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressCantBeZero()
func (iAuthorityRegistry *IAuthorityRegistry) UnpackAddressCantBeZeroError(raw []byte) (*IAuthorityRegistryAddressCantBeZero, error) {
	out := new(IAuthorityRegistryAddressCantBeZero)
	if err := iAuthorityRegistry.abi.UnpackIntoInterface(out, "AddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
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

// IProcessorEndpointMetaData contains all meta data concerning the IProcessorEndpoint contract.
var IProcessorEndpointMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"eventSubType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"}],\"name\":\"markRequestCompleted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"markRequestFailed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "7f11ae708389fe5cf858688c666b6922cc",
}

// IProcessorEndpoint is an auto generated Go binding around an Ethereum contract.
type IProcessorEndpoint struct {
	abi abi.ABI
}

// NewIProcessorEndpoint creates a new instance of IProcessorEndpoint.
func NewIProcessorEndpoint() *IProcessorEndpoint {
	parsed, err := IProcessorEndpointMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IProcessorEndpoint{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IProcessorEndpoint) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a17995a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (iProcessorEndpoint *IProcessorEndpoint) PackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, idx *big.Int) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, depositAmount, idx)
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, idx *big.Int) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, depositAmount, idx)
}

// UnpackGenerateRequestId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a17995a.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackGenerateRequestId(data []byte) ([32]byte, error) {
	out, err := iProcessorEndpoint.abi.Unpack("generateRequestId", data)
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
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool)
func (iProcessorEndpoint *IProcessorEndpoint) PackGetNextPendingRequest() []byte {
	enc, err := iProcessorEndpoint.abi.Pack("getNextPendingRequest")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetNextPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc404c359.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool)
func (iProcessorEndpoint *IProcessorEndpoint) TryPackGetNextPendingRequest() ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("getNextPendingRequest")
}

// GetNextPendingRequestOutput serves as a container for the return parameters of contract
// method GetNextPendingRequest.
type GetNextPendingRequestOutput struct {
	Arg0 StructsPendingRequest
	Arg1 [32]byte
	Arg2 bool
}

// UnpackGetNextPendingRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc404c359.
//
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackGetNextPendingRequest(data []byte) (GetNextPendingRequestOutput, error) {
	out, err := iProcessorEndpoint.abi.Unpack("getNextPendingRequest", data)
	outstruct := new(GetNextPendingRequestOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new(StructsPendingRequest)).(*StructsPendingRequest)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Arg2 = *abi.ConvertType(out[2], new(bool)).(*bool)
	return *outstruct, nil
}

// PackGetPendingRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a1f712.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
func (iProcessorEndpoint *IProcessorEndpoint) PackGetPendingRequests() []byte {
	enc, err := iProcessorEndpoint.abi.Pack("getPendingRequests")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a1f712.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
func (iProcessorEndpoint *IProcessorEndpoint) TryPackGetPendingRequests() ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("getPendingRequests")
}

// UnpackGetPendingRequests is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80a1f712.
//
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
func (iProcessorEndpoint *IProcessorEndpoint) UnpackGetPendingRequests(data []byte) ([]StructsPendingRequest, error) {
	out, err := iProcessorEndpoint.abi.Unpack("getPendingRequests", data)
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
func (iProcessorEndpoint *IProcessorEndpoint) PackGetPendingRequestsSize() []byte {
	enc, err := iProcessorEndpoint.abi.Pack("getPendingRequestsSize")
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackGetPendingRequestsSize() ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("getPendingRequestsSize")
}

// UnpackGetPendingRequestsSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e339384.
//
// Solidity: function getPendingRequestsSize() view returns(uint256)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackGetPendingRequestsSize(data []byte) (*big.Int, error) {
	out, err := iProcessorEndpoint.abi.Unpack("getPendingRequestsSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackIsCurrentPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9fabc36f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (iProcessorEndpoint *IProcessorEndpoint) PackIsCurrentPendingRequest(requestId [32]byte) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("isCurrentPendingRequest", requestId)
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackIsCurrentPendingRequest(requestId [32]byte) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("isCurrentPendingRequest", requestId)
}

// UnpackIsCurrentPendingRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9fabc36f.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackIsCurrentPendingRequest(data []byte) (bool, error) {
	out, err := iProcessorEndpoint.abi.Unpack("isCurrentPendingRequest", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMarkRequestCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d1be609.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function markRequestCompleted(bytes32 requestId, uint256 refund, uint256 applicationFees) returns()
func (iProcessorEndpoint *IProcessorEndpoint) PackMarkRequestCompleted(requestId [32]byte, refund *big.Int, applicationFees *big.Int) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("markRequestCompleted", requestId, refund, applicationFees)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMarkRequestCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d1be609.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function markRequestCompleted(bytes32 requestId, uint256 refund, uint256 applicationFees) returns()
func (iProcessorEndpoint *IProcessorEndpoint) TryPackMarkRequestCompleted(requestId [32]byte, refund *big.Int, applicationFees *big.Int) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("markRequestCompleted", requestId, refund, applicationFees)
}

// PackMarkRequestFailed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b473cb9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function markRequestFailed(bytes32 requestId, uint8 errorCode, string errorMessage) returns()
func (iProcessorEndpoint *IProcessorEndpoint) PackMarkRequestFailed(requestId [32]byte, errorCode uint8, errorMessage string) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("markRequestFailed", requestId, errorCode, errorMessage)
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackMarkRequestFailed(requestId [32]byte, errorCode uint8, errorMessage string) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("markRequestFailed", requestId, errorCode, errorMessage)
}

// PackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ef290e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, bytes signature) returns()
func (iProcessorEndpoint *IProcessorEndpoint) PackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refund, applicationFees, signature)
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refund, applicationFees, signature)
}

// PackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34f23a9d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (iProcessorEndpoint *IProcessorEndpoint) PackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, maxFeeValue *big.Int) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, depositAmount, maxFeeValue)
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, maxFeeValue *big.Int) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, depositAmount, maxFeeValue)
}

// UnpackSubmitRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x34f23a9d.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackSubmitRequest(data []byte) ([32]byte, error) {
	out, err := iProcessorEndpoint.abi.Unpack("submitRequest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackUpdateFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2c35ce8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (iProcessorEndpoint *IProcessorEndpoint) PackUpdateFeeCollector(newFeeCollector common.Address) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("updateFeeCollector", newFeeCollector)
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackUpdateFeeCollector(newFeeCollector common.Address) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("updateFeeCollector", newFeeCollector)
}

// PackUpdateQueueThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c73eeaf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateQueueThreshold(uint256 newThreshold) returns()
func (iProcessorEndpoint *IProcessorEndpoint) PackUpdateQueueThreshold(newThreshold *big.Int) []byte {
	enc, err := iProcessorEndpoint.abi.Pack("updateQueueThreshold", newThreshold)
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
func (iProcessorEndpoint *IProcessorEndpoint) TryPackUpdateQueueThreshold(newThreshold *big.Int) ([]byte, error) {
	return iProcessorEndpoint.abi.Pack("updateQueueThreshold", newThreshold)
}

// IProcessorEndpointFeeCollectorUpdated represents a FeeCollectorUpdated event raised by the IProcessorEndpoint contract.
type IProcessorEndpointFeeCollectorUpdated struct {
	NewFeeCollector common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointFeeCollectorUpdatedEventName = "FeeCollectorUpdated"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointFeeCollectorUpdated) ContractEventName() string {
	return IProcessorEndpointFeeCollectorUpdatedEventName
}

// UnpackFeeCollectorUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FeeCollectorUpdated(address newFeeCollector)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackFeeCollectorUpdatedEvent(log *types.Log) (*IProcessorEndpointFeeCollectorUpdated, error) {
	event := "FeeCollectorUpdated"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointFeeCollectorUpdated)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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

// IProcessorEndpointQueueThresholdUpdated represents a QueueThresholdUpdated event raised by the IProcessorEndpoint contract.
type IProcessorEndpointQueueThresholdUpdated struct {
	NewThreshold *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointQueueThresholdUpdatedEventName = "QueueThresholdUpdated"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointQueueThresholdUpdated) ContractEventName() string {
	return IProcessorEndpointQueueThresholdUpdatedEventName
}

// UnpackQueueThresholdUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event QueueThresholdUpdated(uint256 newThreshold)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackQueueThresholdUpdatedEvent(log *types.Log) (*IProcessorEndpointQueueThresholdUpdated, error) {
	event := "QueueThresholdUpdated"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointQueueThresholdUpdated)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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

// IProcessorEndpointRefund represents a Refund event raised by the IProcessorEndpoint contract.
type IProcessorEndpointRefund struct {
	ApplicationId uint64
	RequestId     [32]byte
	To            common.Address
	Amount        *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointRefundEventName = "Refund"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointRefund) ContractEventName() string {
	return IProcessorEndpointRefundEventName
}

// UnpackRefundEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Refund(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackRefundEvent(log *types.Log) (*IProcessorEndpointRefund, error) {
	event := "Refund"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointRefund)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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

// IProcessorEndpointRequestCompleted represents a RequestCompleted event raised by the IProcessorEndpoint contract.
type IProcessorEndpointRequestCompleted struct {
	RequestId       [32]byte
	ApplicationFees *big.Int
	Status          uint8
	ErrorCode       uint8
	ErrorMessage    string
	Raw             *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointRequestCompletedEventName = "RequestCompleted"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointRequestCompleted) ContractEventName() string {
	return IProcessorEndpointRequestCompletedEventName
}

// UnpackRequestCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestCompleted(bytes32 indexed requestId, uint256 applicationFees, uint8 status, uint8 errorCode, string errorMessage)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackRequestCompletedEvent(log *types.Log) (*IProcessorEndpointRequestCompleted, error) {
	event := "RequestCompleted"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointRequestCompleted)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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

// IProcessorEndpointRequestSubmitted represents a RequestSubmitted event raised by the IProcessorEndpoint contract.
type IProcessorEndpointRequestSubmitted struct {
	RequestId [32]byte
	Sender    common.Address
	Raw       *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointRequestSubmittedEventName = "RequestSubmitted"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointRequestSubmitted) ContractEventName() string {
	return IProcessorEndpointRequestSubmittedEventName
}

// UnpackRequestSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestSubmitted(bytes32 indexed requestId, address indexed sender)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackRequestSubmittedEvent(log *types.Log) (*IProcessorEndpointRequestSubmitted, error) {
	event := "RequestSubmitted"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointRequestSubmitted)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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

// IProcessorEndpointStateRootUpdate represents a StateRootUpdate event raised by the IProcessorEndpoint contract.
type IProcessorEndpointStateRootUpdate struct {
	ApplicationId uint64
	RequestId     [32]byte
	OldStateRoot  [32]byte
	NewStateRoot  [32]byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointStateRootUpdateEventName = "StateRootUpdate"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointStateRootUpdate) ContractEventName() string {
	return IProcessorEndpointStateRootUpdateEventName
}

// UnpackStateRootUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StateRootUpdate(uint64 indexed applicationId, bytes32 indexed requestId, bytes32 oldStateRoot, bytes32 newStateRoot)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackStateRootUpdateEvent(log *types.Log) (*IProcessorEndpointStateRootUpdate, error) {
	event := "StateRootUpdate"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointStateRootUpdate)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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

// IProcessorEndpointUserEvent represents a UserEvent event raised by the IProcessorEndpoint contract.
type IProcessorEndpointUserEvent struct {
	ApplicationId uint64
	RequestId     [32]byte
	EventSubType  common.Hash
	EncryptedData []byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointUserEventEventName = "UserEvent"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointUserEvent) ContractEventName() string {
	return IProcessorEndpointUserEventEventName
}

// UnpackUserEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UserEvent(uint64 indexed applicationId, bytes32 indexed requestId, string indexed eventSubType, bytes encryptedData)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackUserEventEvent(log *types.Log) (*IProcessorEndpointUserEvent, error) {
	event := "UserEvent"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointUserEvent)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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

// IProcessorEndpointWithdrawal represents a Withdrawal event raised by the IProcessorEndpoint contract.
type IProcessorEndpointWithdrawal struct {
	ApplicationId uint64
	RequestId     [32]byte
	To            common.Address
	Amount        *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const IProcessorEndpointWithdrawalEventName = "Withdrawal"

// ContractEventName returns the user-defined event name.
func (IProcessorEndpointWithdrawal) ContractEventName() string {
	return IProcessorEndpointWithdrawalEventName
}

// UnpackWithdrawalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdrawal(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount)
func (iProcessorEndpoint *IProcessorEndpoint) UnpackWithdrawalEvent(log *types.Log) (*IProcessorEndpointWithdrawal, error) {
	event := "Withdrawal"
	if log.Topics[0] != iProcessorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IProcessorEndpointWithdrawal)
	if len(log.Data) > 0 {
		if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iProcessorEndpoint.abi.Events[event].Inputs {
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
func (iProcessorEndpoint *IProcessorEndpoint) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["AddressCantBeZero"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["AuthorityNotAllowed"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackAuthorityNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["FeeValueBelowMinimum"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackFeeValueBelowMinimumError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InsufficientBalance"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InvalidApplicationId"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInvalidApplicationIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InvalidPayload"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInvalidPayloadError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InvalidProtocolVersion"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInvalidProtocolVersionError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InvalidRequestId"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInvalidRequestIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InvalidStateRoot"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInvalidStateRootError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["InvalidValue"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackInvalidValueError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["QueueThresholdExceeded"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackQueueThresholdExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], iProcessorEndpoint.abi.Errors["TransferFailed"].ID.Bytes()[:4]) {
		return iProcessorEndpoint.UnpackTransferFailedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IProcessorEndpointAddressCantBeZero represents a AddressCantBeZero error raised by the IProcessorEndpoint contract.
type IProcessorEndpointAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressCantBeZero()
func IProcessorEndpointAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x4b054c82bb9e4ebe90744902747c4e86028dd2a122c63de8810d9a1c2e84617d")
}

// UnpackAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressCantBeZero()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackAddressCantBeZeroError(raw []byte) (*IProcessorEndpointAddressCantBeZero, error) {
	out := new(IProcessorEndpointAddressCantBeZero)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "AddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointAuthorityNotAllowed represents a AuthorityNotAllowed error raised by the IProcessorEndpoint contract.
type IProcessorEndpointAuthorityNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuthorityNotAllowed()
func IProcessorEndpointAuthorityNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xc4550a169e74f716a0e1f2ed847c9b836363900678799385b1658d47630ac161")
}

// UnpackAuthorityNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuthorityNotAllowed()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackAuthorityNotAllowedError(raw []byte) (*IProcessorEndpointAuthorityNotAllowed, error) {
	out := new(IProcessorEndpointAuthorityNotAllowed)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "AuthorityNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointFeeValueBelowMinimum represents a FeeValueBelowMinimum error raised by the IProcessorEndpoint contract.
type IProcessorEndpointFeeValueBelowMinimum struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FeeValueBelowMinimum()
func IProcessorEndpointFeeValueBelowMinimumErrorID() common.Hash {
	return common.HexToHash("0xc77f97ebda9854b3f6367f2d96830b6ea60a9ab7e12c67ede499fadeb3f9cef7")
}

// UnpackFeeValueBelowMinimumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FeeValueBelowMinimum()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackFeeValueBelowMinimumError(raw []byte) (*IProcessorEndpointFeeValueBelowMinimum, error) {
	out := new(IProcessorEndpointFeeValueBelowMinimum)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "FeeValueBelowMinimum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInsufficientBalance represents a InsufficientBalance error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInsufficientBalance struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientBalance()
func IProcessorEndpointInsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xf4d678b8ce6b5157126b1484a53523762a93571537a7d5ae97d8014a44715c94")
}

// UnpackInsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientBalance()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInsufficientBalanceError(raw []byte) (*IProcessorEndpointInsufficientBalance, error) {
	out := new(IProcessorEndpointInsufficientBalance)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInvalidApplicationId represents a InvalidApplicationId error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInvalidApplicationId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidApplicationId()
func IProcessorEndpointInvalidApplicationIdErrorID() common.Hash {
	return common.HexToHash("0x9ff9a83079daf92a88066c1b25d04540f443a8d245022db6e1e62beaa6d6bc5e")
}

// UnpackInvalidApplicationIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidApplicationId()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInvalidApplicationIdError(raw []byte) (*IProcessorEndpointInvalidApplicationId, error) {
	out := new(IProcessorEndpointInvalidApplicationId)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InvalidApplicationId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInvalidPayload represents a InvalidPayload error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInvalidPayload struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPayload()
func IProcessorEndpointInvalidPayloadErrorID() common.Hash {
	return common.HexToHash("0x7c6953f9f55fcfd713fe1d72e370e4757e17ae7187e1acdcc02b4b76b56a9a35")
}

// UnpackInvalidPayloadError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPayload()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInvalidPayloadError(raw []byte) (*IProcessorEndpointInvalidPayload, error) {
	out := new(IProcessorEndpointInvalidPayload)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InvalidPayload", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInvalidProtocolVersion represents a InvalidProtocolVersion error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInvalidProtocolVersion struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidProtocolVersion()
func IProcessorEndpointInvalidProtocolVersionErrorID() common.Hash {
	return common.HexToHash("0x5428eae7f7dcfe744a90b56d360bbe512220bcf69b1e537f3d2d3b28cf2bf92d")
}

// UnpackInvalidProtocolVersionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidProtocolVersion()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInvalidProtocolVersionError(raw []byte) (*IProcessorEndpointInvalidProtocolVersion, error) {
	out := new(IProcessorEndpointInvalidProtocolVersion)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InvalidProtocolVersion", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInvalidRequestId represents a InvalidRequestId error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInvalidRequestId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRequestId()
func IProcessorEndpointInvalidRequestIdErrorID() common.Hash {
	return common.HexToHash("0xba0514c0d5ec80c22a6ebd3f2e6691e2f9ad3d1402978084361a7537d8546138")
}

// UnpackInvalidRequestIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRequestId()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInvalidRequestIdError(raw []byte) (*IProcessorEndpointInvalidRequestId, error) {
	out := new(IProcessorEndpointInvalidRequestId)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InvalidRequestId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInvalidSignature represents a InvalidSignature error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature()
func IProcessorEndpointInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x8baa579fce362245063d36f11747a89dd489c54795634fc673cc0e0db51fedc5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInvalidSignatureError(raw []byte) (*IProcessorEndpointInvalidSignature, error) {
	out := new(IProcessorEndpointInvalidSignature)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInvalidStateRoot represents a InvalidStateRoot error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInvalidStateRoot struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidStateRoot()
func IProcessorEndpointInvalidStateRootErrorID() common.Hash {
	return common.HexToHash("0xb6fac0303f809aaca0354efa81577129a51cc71ce625abce5d9ded4cf63ebabc")
}

// UnpackInvalidStateRootError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidStateRoot()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInvalidStateRootError(raw []byte) (*IProcessorEndpointInvalidStateRoot, error) {
	out := new(IProcessorEndpointInvalidStateRoot)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InvalidStateRoot", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointInvalidValue represents a InvalidValue error raised by the IProcessorEndpoint contract.
type IProcessorEndpointInvalidValue struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidValue()
func IProcessorEndpointInvalidValueErrorID() common.Hash {
	return common.HexToHash("0xaa7feadc1a228aee3d4475b95bb4402ebafd3841876f41c4469039f5ee2998d5")
}

// UnpackInvalidValueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidValue()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackInvalidValueError(raw []byte) (*IProcessorEndpointInvalidValue, error) {
	out := new(IProcessorEndpointInvalidValue)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "InvalidValue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointQueueThresholdExceeded represents a QueueThresholdExceeded error raised by the IProcessorEndpoint contract.
type IProcessorEndpointQueueThresholdExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QueueThresholdExceeded()
func IProcessorEndpointQueueThresholdExceededErrorID() common.Hash {
	return common.HexToHash("0xd8219b3da97d7adbef217ce9f220bbf38600f597fb984fd0620ea7df76ea17f3")
}

// UnpackQueueThresholdExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QueueThresholdExceeded()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackQueueThresholdExceededError(raw []byte) (*IProcessorEndpointQueueThresholdExceeded, error) {
	out := new(IProcessorEndpointQueueThresholdExceeded)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "QueueThresholdExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IProcessorEndpointTransferFailed represents a TransferFailed error raised by the IProcessorEndpoint contract.
type IProcessorEndpointTransferFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TransferFailed()
func IProcessorEndpointTransferFailedErrorID() common.Hash {
	return common.HexToHash("0x90b8ec1877afffd816d05d9b13947f3ff18ec5851c38bad15ec2b710f92391b1")
}

// UnpackTransferFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TransferFailed()
func (iProcessorEndpoint *IProcessorEndpoint) UnpackTransferFailedError(raw []byte) (*IProcessorEndpointTransferFailed, error) {
	out := new(IProcessorEndpointTransferFailed)
	if err := iProcessorEndpoint.abi.UnpackIntoInterface(out, "TransferFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorMetaData contains all meta data concerning the ITeeAuthenticator contract.
var ITeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (iTeeAuthenticator *ITeeAuthenticator) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iTeeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return iTeeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ITeeAuthenticatorTeeIsNotSet represents a TeeIsNotSet error raised by the ITeeAuthenticator contract.
type ITeeAuthenticatorTeeIsNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TeeIsNotSet()
func ITeeAuthenticatorTeeIsNotSetErrorID() common.Hash {
	return common.HexToHash("0xf64b1d7b3df6940d5da676c1cc07a6144e620b0858b161be25b6cd0ef7569425")
}

// UnpackTeeIsNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TeeIsNotSet()
func (iTeeAuthenticator *ITeeAuthenticator) UnpackTeeIsNotSetError(raw []byte) (*ITeeAuthenticatorTeeIsNotSet, error) {
	out := new(ITeeAuthenticatorTeeIsNotSet)
	if err := iTeeAuthenticator.abi.UnpackIntoInterface(out, "TeeIsNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"eventSubType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"APPLICATION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"}],\"name\":\"markRequestCompleted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"markRequestFailed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a35c5449187bc72768eac252f2f30417e4",
	Bin: "0x6080604052346100335761001d6100146101b4565b939290926103f9565b610025610038565b61443c6106db823961443c90f35b61003e565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b9061006a90610042565b810190811060018060401b0382111761008257604052565b61004c565b9061009a610093610038565b9283610060565b565b5f80fd5b60018060a01b031690565b6100b4906100a0565b90565b6100c0906100ab565b90565b6100cc816100b7565b036100d357565b5f80fd5b905051906100e4826100c3565b565b6100ef906100ab565b90565b6100fb816100e6565b0361010257565b5f80fd5b90505190610113826100f2565b565b61011e816100ab565b0361012557565b5f80fd5b9050519061013682610115565b565b90565b61014481610138565b0361014b57565b5f80fd5b9050519061015c8261013b565b565b919060a0838203126101af57610176815f85016100d7565b926101848260208301610106565b926101ac6101958460408501610129565b936101a38160608601610129565b9360800161014f565b90565b61009c565b6101d2614b17803803806101c781610087565b92833981019061015e565b9091929394565b5f1b90565b906101ea5f19916101d9565b9181191691161790565b90565b90565b61020e610209610213926101f4565b6101f7565b610138565b90565b90565b9061022e610229610235926101fa565b610216565b82546101de565b9055565b90565b61025061024b61025592610239565b6101f7565b6100a0565b90565b6102619061023c565b90565b61027861027361027d926100a0565b6101f7565b6100a0565b90565b61028990610264565b90565b61029590610280565b90565b6102a190610264565b90565b6102ad90610298565b90565b5f0190565b906102c660018060a01b03916101d9565b9181191691161790565b6102d990610280565b90565b90565b906102f46102ef6102fb926102d0565b6102dc565b82546102b5565b9055565b61030890610298565b90565b90565b9061032361031e61032a926102ff565b61030b565b82546102b5565b9055565b61033790610264565b90565b6103439061032e565b90565b61034f9061032e565b90565b90565b9061036a61036561037192610346565b610352565b82546102b5565b9055565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103d16103cc6103d692610138565b6101f7565b610138565b90565b906103ee6103e96103f5926103bd565b610216565b82546101de565b9055565b91939290610409600a6006610219565b8261042c61042661042161041c5f610258565b61028c565b6100b7565b916100b7565b1480156104fe575b80156104dc575b80156104ba575b61049e5761049c946104666104869261045f6104949660076102df565b600861030e565b6104796104728261033a565b600a610355565b610481610375565b6105c9565b5061048f610399565b6105c9565b5060096103d9565b565b5f632582a64160e11b8152806104b6600482016102b0565b0390fd5b50816104d66104d06104cb5f610258565b6100ab565b916100ab565b14610442565b50846104f86104f26104ed5f610258565b6100ab565b916100ab565b1461043b565b508061052261051c6105176105125f610258565b6102a4565b6100e6565b916100e6565b14610434565b5f90565b151590565b90565b61053d90610531565b90565b9061054a90610534565b5f5260205260405f2090565b61055f90610264565b90565b61056b90610556565b90565b9061057890610562565b5f5260205260405f2090565b9061059060ff916101d9565b9181191691161790565b6105a39061052c565b90565b90565b906105be6105b96105c59261059a565b6105a6565b8254610584565b9055565b6105d1610528565b506105e66105e08284906106a0565b1561052c565b5f1461066e5761060d60016106085f610600818690610540565b01859061056e565b6105a9565b906106166106cd565b9061065361064d6106477f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610534565b92610562565b92610562565b9261065c610038565b80610666816102b0565b0390a4600190565b50505f90565b5f1c90565b60ff1690565b61068b61069091610674565b610679565b90565b61069d905461067f565b90565b6106c6915f6106bb6106c1936106b4610528565b5082610540565b0161056e565b610693565b90565b5f90565b6106d56106c9565b50339056fe6101606040526004361015610014575b611894565b61001e5f356101ed565b806301ffc9a7146101e85780631664deeb146101e3578063248a9ca3146101de5780632a0acc6a146101d95780632f2ff15d146101d457806334f23a9d146101cf57806336568abe146101ca5780633b473cb9146101c55780635423112f146101c05780635d1be609146101bb5780635e339384146101b65780636e91d133146101b157806380a1f712146101ac57806391d14854146101a75780639588eca2146101a25780639968da961461019d5780639a17995a146101985780639c73eeaf146101935780639c8d416f1461018e5780639fabc36f14610189578063a217fddf14610184578063a9ba045e1461017f578063a9ef290e1461017a578063aa3aa46014610175578063c404c35914610170578063c415b95c1461016b578063d2c35ce814610166578063d547741f146101615763f7d9cc1e0361000f5761185f565b61181c565b6117e9565b611773565b6116cd565b6115b6565b611528565b611330565b61128d565b61121d565b6111e8565b611180565b611126565b611071565b610ff0565b610f85565b610f50565b610d4a565b610c78565b610c2f565b610bb6565b61073b565b610659565b610627565b61048c565b6103eb565b610387565b610311565b610279565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f80fd5b63ffffffff60e01b1690565b61021a81610205565b0361022157565b5f80fd5b9050359061023282610211565b565b9060208282031261024d5761024a915f01610225565b90565b6101fd565b151590565b61026090610252565b9052565b9190610277905f60208501940190610257565b565b346102a9576102a561029461028f366004610234565b61189c565b61029c6101f3565b91829182610264565b0390f35b6101f9565b5f9103126102b857565b6101fd565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b6102e96102bd565b90565b90565b6102f8906102ec565b9052565b919061030f905f602085019401906102ef565b565b34610341576103213660046102ae565b61033d61032c6102e1565b6103346101f3565b918291826102fc565b0390f35b6101f9565b61034f816102ec565b0361035657565b5f80fd5b9050359061036782610346565b565b906020828203126103825761037f915f0161035a565b90565b6101fd565b346103b7576103b36103a261039d366004610369565b6118f6565b6103aa6101f3565b918291826102fc565b0390f35b6101f9565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103e86103bc565b90565b3461041b576103fb3660046102ae565b6104176104066103e0565b61040e6101f3565b918291826102fc565b0390f35b6101f9565b60018060a01b031690565b61043490610420565b90565b6104408161042b565b0361044757565b5f80fd5b9050359061045882610437565b565b9190604083820312610482578061047661047f925f860161035a565b9360200161044b565b90565b6101fd565b5f0190565b346104bb576104a561049f36600461045a565b90611941565b6104ad6101f3565b806104b781610487565b0390f35b6101f9565b60ff1690565b6104cf816104c0565b036104d657565b5f80fd5b905035906104e7826104c6565b565b67ffffffffffffffff1690565b6104ff816104e9565b0361050657565b5f80fd5b90503590610517826104f6565b565b6004111561052357565b5f80fd5b9050359061053482610519565b565b5f80fd5b5f80fd5b5f80fd5b909182601f8301121561057c5781359167ffffffffffffffff831161057757602001926001830284011161057257565b61053e565b61053a565b610536565b90565b61058d81610581565b0361059457565b5f80fd5b905035906105a582610584565b565b91909160c081840312610622576105c0835f83016104da565b926105ce816020840161050a565b926105dc8260408501610527565b9260608101359167ffffffffffffffff831161061d576106018461061a948401610542565b9390946106118160808601610598565b9360a001610598565b90565b610201565b6101fd565b6106556106446106383660046105a7565b959490949391936124b9565b61064c6101f3565b918291826102fc565b0390f35b346106885761067261066c36600461045a565b906124d3565b61067a6101f3565b8061068481610487565b0390f35b6101f9565b6011111561069757565b5f80fd5b905035906106a88261068d565b565b909182601f830112156106e45781359167ffffffffffffffff83116106df5760200192600183028401116106da57565b61053e565b61053a565b610536565b9160608383031261073657610700825f850161035a565b9261070e836020830161069b565b92604082013567ffffffffffffffff81116107315761072d92016106aa565b9091565b610201565b6101fd565b3461076d5761075761074e3660046106e9565b92919091612939565b61075f6101f3565b8061076981610487565b0390f35b6101f9565b61077b906102ec565b90565b9061078890610772565b5f5260205260405f2090565b5f1c90565b90565b6107a86107ad91610794565b610799565b90565b6107ba905461079c565b90565b90565b6107cc6107d191610794565b6107bd565b90565b6107de90546107c0565b90565b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610815575b602083101461081057565b6107e1565b91607f1691610805565b60209181520190565b5f5260205f2090565b905f929180549061084b610844836107f5565b809461081f565b916001811690815f146108a25750600114610866575b505050565b6108739192939450610828565b915f925b81841061088a57505001905f8080610861565b60018160209295939554848601520191019290610877565b92949550505060ff19168252151560200201905f8080610861565b906108c791610831565b90565b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906108f2906108ca565b810190811067ffffffffffffffff82111761090c57604052565b6108d4565b9061093161092a926109216101f3565b938480926108bd565b03836108e8565b565b60018060a01b031690565b61094a61094f91610794565b610933565b90565b61095c905461093e565b90565b60a01c90565b67ffffffffffffffff1690565b61097e6109839161095f565b610965565b90565b6109909054610972565b90565b60ff1690565b6109a56109aa916101ed565b610993565b90565b6109b79054610999565b90565b60e81c90565b60ff1690565b6109d26109d7916109ba565b6109c0565b90565b6109e490546109c6565b90565b6109f290600261077e565b6109fd5f82016107b0565b91610a0a600183016107b0565b91610a17600282016107b0565b91610a24600383016107d4565b91610a3160048201610911565b91610a3e60058301610952565b91610a4b60058201610986565b91610a636005610a5c8185016109ad565b93016109da565b90565b610a6f90610581565b9052565b5190565b60209181520190565b90825f9392825e0152565b610aaa610ab3602093610ab893610aa181610a73565b93848093610a77565b95869101610a80565b6108ca565b0190565b610ac59061042b565b9052565b610ad2906104e9565b9052565b610adf906104c0565b9052565b634e487b7160e01b5f52602160045260245ffd5b60041115610b0157565b610ae3565b90610b1082610af7565b565b610b1b90610b06565b90565b610b2790610b12565b9052565b94610b8e610100979b9a9892610b9992610b81610bad98610b77610bb49e99610b6d610ba39a60208f610b6661012082019a5f830190610a66565b0190610a66565b60408d0190610a66565b60608b01906102ef565b88820360808a0152610a8b565b9a60a0870190610abc565b60c0850190610ac9565b60e0830190610ad6565b0190610b1e565b565b34610bf057610bec610bd1610bcc366004610369565b6109e7565b95610be39997999591959492946101f3565b998a998a610b2b565b0390f35b6101f9565b9091606082840312610c2a57610c27610c10845f850161035a565b93610c1e8160208601610598565b93604001610598565b90565b6101fd565b34610c5e57610c48610c42366004610bf5565b91612b3f565b610c506101f3565b80610c5a81610487565b0390f35b6101f9565b9190610c76905f60208501940190610a66565b565b34610ca857610c883660046102ae565b610ca4610c93612b4c565b610c9b6101f3565b91829182610c63565b0390f35b6101f9565b1c90565b60018060a01b031690565b610ccc906008610cd19302610cad565b610cb1565b90565b90610cdf9154610cbc565b90565b610cee60075f90610cd4565b90565b90565b610d08610d03610d0d92610420565b610cf1565b610420565b90565b610d1990610cf4565b90565b610d2590610d10565b90565b610d3190610d1c565b9052565b9190610d48905f60208501940190610d28565b565b34610d7a57610d5a3660046102ae565b610d76610d65610ce2565b610d6d6101f3565b91829182610d35565b0390f35b6101f9565b5190565b60209181520190565b60200190565b610d9b90610581565b9052565b610da8906102ec565b9052565b610dcb610dd4602093610dd993610dc281610a73565b9384809361081f565b95869101610a80565b6108ca565b0190565b610de69061042b565b9052565b610df3906104e9565b9052565b610e00906104c0565b9052565b610e0d90610b12565b9052565b90610ebb9061010080610e7a6101208401610e325f8801515f870190610d92565b610e4460208801516020870190610d92565b610e5660408801516040870190610d92565b610e6860608801516060870190610d9f565b60808701518582036080870152610dac565b94610e8d60a082015160a0860190610ddd565b610e9f60c082015160c0860190610dea565b610eb160e082015160e0860190610df7565b0151910190610e04565b90565b90610ec891610e11565b90565b60200190565b90610ee5610ede83610d7f565b8092610d83565b9081610ef660208302840194610d8c565b925f915b838310610f0957505050505090565b90919293946020610f2b610f2583856001950387528951610ebe565b97610ecb565b9301930191939290610efa565b610f4d9160208201915f818403910152610ed1565b90565b34610f8057610f603660046102ae565b610f7c610f6b612dff565b610f736101f3565b91829182610f38565b0390f35b6101f9565b34610fb657610fb2610fa1610f9b36600461045a565b90612eed565b610fa96101f3565b91829182610264565b0390f35b6101f9565b610fcb906008610fd09302610cad565b6107bd565b90565b90610fde9154610fbb565b90565b610fed60015f90610fd3565b90565b34611020576110003660046102ae565b61101c61100b610fe1565b6110136101f3565b918291826102fc565b0390f35b6101f9565b90565b61103c61103761104192611025565b610cf1565b6104e9565b90565b61104e6001611028565b90565b611059611044565b90565b919061106f905f60208501940190610ac9565b565b346110a1576110813660046102ae565b61109d61108c611051565b6110946101f3565b9182918261105c565b0390f35b6101f9565b91909160c081840312611121576110bf835f830161044b565b926110cd816020840161050a565b926110db8260408501610527565b9260608101359167ffffffffffffffff831161111c5761110084611119948401610542565b9390946111108160808601610598565b9360a001610598565b90565b610201565b6101fd565b3461115d5761115961114861113c3660046110a6565b95949094939193613023565b6111506101f3565b918291826102fc565b0390f35b6101f9565b9060208282031261117b57611178915f01610598565b90565b6101fd565b346111ae57611198611193366004611162565b61310d565b6111a06101f3565b806111aa81610487565b0390f35b6101f9565b6111c39060086111c89302610cad565b610799565b90565b906111d691546111b3565b90565b6111e560095f906111cb565b90565b34611218576111f83660046102ae565b6112146112036111d9565b61120b6101f3565b91829182610c63565b0390f35b6101f9565b3461124d57611249611238611233366004610369565b613118565b6112406101f3565b91829182610264565b0390f35b6101f9565b90565b5f1b90565b61126e61126961127392611252565b611255565b6102ec565b90565b61127f5f61125a565b90565b61128a611276565b90565b346112bd5761129d3660046102ae565b6112b96112a8611282565b6112b06101f3565b918291826102fc565b0390f35b6101f9565b60018060a01b031690565b6112dd9060086112e29302610cad565b6112c2565b90565b906112f091546112cd565b90565b6112ff60085f906112e5565b90565b61130b90610d10565b90565b61131790611302565b9052565b919061132e905f6020850194019061130e565b565b34611360576113403660046102ae565b61135c61134b6112f3565b6113536101f3565b9182918261131b565b0390f35b6101f9565b909182601f8301121561139f5781359167ffffffffffffffff831161139a57602001926020830284011161139557565b61053e565b61053a565b610536565b909182601f830112156113de5781359167ffffffffffffffff83116113d95760200192602083028401116113d457565b61053e565b61053a565b610536565b909182601f8301121561141d5781359167ffffffffffffffff831161141857602001926040830284011161141357565b61053e565b61053a565b610536565b9190610140838203126115235761143b815f850161050a565b92611449826020830161035a565b92611457836040840161035a565b92611465816060850161035a565b92608081013567ffffffffffffffff811161151e5782611486918301611365565b92909360a083013567ffffffffffffffff811161151957826114a99185016113a4565b92909360c081013567ffffffffffffffff811161151457826114cc9183016113e3565b9290936114dc8260e08501610598565b926114eb836101008301610598565b9261012082013567ffffffffffffffff811161150f5761150b9201610542565b9091565b610201565b610201565b610201565b610201565b6101fd565b346115695761155361153b366004611422565b9c9b909b9a919a999299989398979497969596613eb7565b61155b6101f3565b8061156581610487565b0390f35b6101f9565b61158261157d61158792611252565b610cf1565b6104c0565b90565b6115935f61156e565b90565b61159e61158a565b90565b91906115b4905f60208501940190610ad6565b565b346115e6576115c63660046102ae565b6115e26115d1611596565b6115d96101f3565b918291826115a1565b0390f35b6101f9565b906116959061010080611654610120840161160c5f8801515f870190610d92565b61161e60208801516020870190610d92565b61163060408801516040870190610d92565b61164260608801516060870190610d9f565b60808701518582036080870152610dac565b9461166760a082015160a0860190610ddd565b61167960c082015160c0860190610dea565b61168b60e082015160e0860190610df7565b0151910190610e04565b90565b6040906116c46116b96116cb9597969460608401908482035f8601526115eb565b9660208301906102ef565b0190610257565b565b34611700576116dd3660046102ae565b6116fc6116e8613eda565b6116f39391936101f3565b93849384611698565b0390f35b6101f9565b60018060a01b031690565b6117209060086117259302610cad565b611705565b90565b906117339154611710565b90565b611742600a5f90611728565b90565b61174e90610420565b90565b61175a90611745565b9052565b9190611771905f60208501940190611751565b565b346117a3576117833660046102ae565b61179f61178e611736565b6117966101f3565b9182918261175e565b0390f35b6101f9565b6117b181611745565b036117b857565b5f80fd5b905035906117c9826117a8565b565b906020828203126117e4576117e1915f016117bc565b90565b6101fd565b34611817576118016117fc3660046117cb565b614079565b6118096101f3565b8061181381610487565b0390f35b6101f9565b3461184b5761183561182f36600461045a565b906140ae565b61183d6101f3565b8061184781610487565b0390f35b6101f9565b61185c60065f906111cb565b90565b3461188f5761186f3660046102ae565b61188b61187a611850565b6118826101f3565b91829182610c63565b0390f35b6101f9565b5f80fd5b5f90565b6118a4611898565b50806118bf6118b9637965db0b60e01b610205565b91610205565b149081156118cc575b5090565b6118d691506140ba565b5f6118c8565b5f90565b906118ea90610772565b5f5260205260405f2090565b600161190e611914926119076118dc565b505f6118e0565b016107d4565b90565b906119329161192d611928826118f6565b6140e0565b611934565b565b9061193e91614139565b50565b9061194b91611917565b565b969594939291908061196e61196861196361158a565b6104c0565b916104c0565b0361197f5761197c9761199b565b90565b5f635428eae760e01b81528061199760048201610487565b0390fd5b96959493929190816119bc6119b66119b1611044565b6104e9565b916104e9565b036119cd576119ca97612133565b90565b5f6309ff9a8360e41b8152806119e560048201610487565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b611a0c611a1291939293610581565b92610581565b8201809211611a1d57565b6119e9565b611a36611a31611a3b92611252565b610cf1565b610581565b90565b611a4a611a4f91610794565b6112c2565b90565b611a5c9054611a3e565b90565b60e01b90565b611a6e81610252565b03611a7557565b5f80fd5b90505190611a8682611a65565b565b90602082820312611aa157611a9e915f01611a79565b90565b6101fd565b611aba611ab5611abf926104e9565b610cf1565b610581565b90565b611acb90611aa6565b9052565b916020611af0929493611ae960408201965f830190611ac2565b0190610abc565b565b611afa6101f3565b3d5f823e3d90fd5b5090565b90565b611b1d611b18611b2292611b06565b610cf1565b610581565b90565b90611b38611b316101f3565b92836108e8565b565b611b45610120611b25565b90565b90611b5290610581565b9052565b90611b60906102ec565b9052565b5f80fd5b67ffffffffffffffff8111611b8657611b826020916108ca565b0190565b6108d4565b90825f939282370152565b90929192611bab611ba682611b68565b611b25565b93818552602085019082840111611bc757611bc592611b8b565b565b611b64565b611bd7913691611b96565b90565b52565b90611be79061042b565b9052565b90611bf5906104e9565b9052565b90611c03906104c0565b9052565b90611c1190610b06565b9052565b611c1f9051610581565b90565b90611c2e5f1991611255565b9181191691161790565b611c4c611c47611c5192610581565b610cf1565b610581565b90565b90565b90611c6c611c67611c7392611c38565b611c54565b8254611c22565b9055565b611c8190516102ec565b90565b611c8d90610794565b90565b90611ca5611ca0611cac92610772565b611c84565b8254611c22565b9055565b5190565b601f602091010490565b1b90565b91906008611cdd910291611cd75f1984611cbe565b92611cbe565b9181191691161790565b9190611cfd611cf8611d0593611c38565b611c54565b908354611cc2565b9055565b5f90565b611d1f91611d19611d09565b91611ce7565b565b5b818110611d2d575050565b80611d3a5f600193611d0d565b01611d22565b9190601f8111611d50575b505050565b611d5c611d8193610828565b906020611d6884611cb4565b83019310611d89575b611d7a90611cb4565b0190611d21565b5f8080611d4b565b9150611d7a81929050611d71565b90611da7905f1990600802610cad565b191690565b81611db691611d97565b906002021790565b90611dc881610a73565b9067ffffffffffffffff8211611e8857611dec82611de685546107f5565b85611d40565b602090601f8311600114611e2057918091611e0f935f92611e14575b5050611dac565b90555b565b90915001515f80611e08565b601f19831691611e2f85610828565b925f5b818110611e7057509160029391856001969410611e56575b50505002019055611e12565b611e66910151601f841690611d97565b90555f8080611e4a565b91936020600181928787015181550195019201611e32565b6108d4565b90611e9791611dbe565b565b611ea3905161042b565b90565b90611eb760018060a01b0391611255565b9181191691161790565b611eca90610d10565b90565b90565b90611ee5611ee0611eec92611ec1565b611ecd565b8254611ea6565b9055565b611efa90516104e9565b90565b60a01b90565b90611f1967ffffffffffffffff60a01b91611efd565b9181191691161790565b611f37611f32611f3c926104e9565b610cf1565b6104e9565b90565b90565b90611f57611f52611f5e92611f23565b611f3f565b8254611f03565b9055565b611f6c90516104c0565b90565b90611f7e60ff60e01b91611a5f565b9181191691161790565b611f9c611f97611fa1926104c0565b610cf1565b6104c0565b90565b90565b90611fbc611fb7611fc392611f88565b611fa4565b8254611f6f565b9055565b611fd19051610b06565b90565b60e81b90565b90611fe960ff60e81b91611fd4565b9181191691161790565b611ffc90610b06565b90565b90565b9061201761201261201e92611ff3565b611fff565b8254611fda565b9055565b906120fa6101006005612100946120465f82016120405f8801611c15565b90611c57565b61205f6001820161205960208801611c15565b90611c57565b6120786002820161207260408801611c15565b90611c57565b6120916003820161208b60608801611c77565b90611c90565b6120aa600482016120a460808801611cb0565b90611e8d565b6120c28282016120bc60a08801611e99565b90611ed0565b6120da8282016120d460c08801611ef0565b90611f42565b6120f28282016120ec60e08801611f62565b90611fa7565b019201611fc7565b90612002565b565b9061210c91612022565b565b9061211890611c38565b5f5260205260405f2090565b60016121309101610581565b90565b969490929396503461215761215161214c8989906119fd565b610581565b91610581565b0361249d578461217861217261216d60096107b0565b610581565b91610581565b1061248157612185612b4c565b6121a061219a61219560066107b0565b610581565b91610581565b101561246557836121ba6121b46003610b06565b91610b06565b145f1461235f576121cc878290611b02565b6121df6121d96085611b09565b91610581565b03612343575b3382858984908a926121f760056107b0565b9461220196613023565b9642969590889092913394959697612217611b3a565b995f8b019061222591611b48565b60208a019061223391611b48565b604089019061224191611b48565b606088019061224f91611b56565b61225891611bcc565b608086019061226691611bda565b60a085019061227491611bdd565b60c084019061228291611beb565b60e083019061229091611bf9565b61010082019061229f91611c07565b6002826122ab9161077e565b906122b591612102565b8060036122c260056107b0565b6122cb9161210e565b906122d591611c90565b6122df60056107b0565b6122e890612124565b6122f3906005611c57565b80337f73394f43049193ecd2e0c22eefa0ecf10987ce9c72da906a82a3336dc2ac05479161232090610772565b9061232a90611ec1565b916123336101f3565b8061233d81610487565b0390a390565b5f637c6953f960e01b81528061235b60048201610487565b0390fd5b8361237361236d6002610b06565b91610b06565b1461237e575b6121e5565b8561239161238b5f611a22565b91610581565b03612449576123a86123a36008611a52565b611302565b602063116ad4a69184906123ce33946123d96123c26101f3565b96879586948594611a5f565b845260048401611acf565b03915afa8015612444576123f5915f91612416575b5015610252565b15612379575f63622a850b60e11b81528061241260048201610487565b0390fd5b612437915060203d811161243d575b61242f81836108e8565b810190611a88565b5f6123ee565b503d612425565b611af2565b5f632a9ffab760e21b81528061246160048201610487565b0390fd5b5f63d8219b3d60e01b81528061247d60048201610487565b0390fd5b5f63c77f97eb60e01b81528061249960048201610487565b0390fd5b5f632a9ffab760e21b8152806124b560048201610487565b0390fd5b906124d09695949392916124cb6118dc565b61194d565b90565b90806124ee6124e86124e36141e8565b61042b565b9161042b565b036124ff576124fc916141f5565b50565b5f63334bd91960e11b81528061251760048201610487565b0390fd5b9061253793929161253261252d6102bd565b6140e0565b612713565b565b90565b61254b61255191939293610581565b92610581565b820391821161255c57565b6119e9565b61256a90610cf4565b90565b61257690612561565b90565b61258290610d10565b90565b905090565b6125955f8092612585565b0190565b6125a29061258a565b90565b906125b76125b283611b68565b611b25565b918252565b606090565b3d5f146125dc576125d13d6125a5565b903d5f602084013e5b565b6125e46125bc565b906125da565b6125f390612579565b9052565b91602061261892949361261160408201965f8301906125ea565b0190610a66565b565b61262661262b91610794565b611705565b90565b612638905461261a565b90565b6003111561264557565b610ae3565b906126548261263b565b565b61265f9061264a565b90565b61266b90612656565b9052565b6011111561267957565b610ae3565b906126888261266f565b565b6126939061267e565b90565b61269f9061268a565b9052565b60209181520190565b91906126c6816126bf816126cb956126a3565b8095611b8b565b6108ca565b0190565b90939161271095936126f9612703926126ef60808601985f870190610a66565b6020850190612662565b6040830190612696565b60608185039101526126ac565b90565b919261272761272184613118565b15610252565b61291d5761273f61273a6002859061077e565b612539565b9261274a60096107b0565b936127536142cc565b61277e612762600183016107b0565b612778612771600285016107b0565b889061253c565b906119fd565b908161279261278c5f611a22565b91610581565b1161286f575b50505f806127ae6127a9600a61262e565b612579565b866127b76101f3565b90816127c281612599565b03925af16127ce6125c1565b505f146128215761281b909391600193956128097f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198896610772565b966128126101f3565b958695866126cf565b0390a25b565b612867909391600293956128557f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198896610772565b9661285e6101f3565b958695866126cf565b0390a261281f565b61288361287e60058301610952565b61256d565b905f8061288f84612579565b856128986101f3565b90816128a381612599565b03925af16128af6125c1565b506128bb575b50612798565b60056128c79101610986565b8391926128fd6128f77f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611f23565b93610772565b936129126129096101f3565b928392836125f7565b0390a35f80806128b5565b5f6302e8145360e61b81528061293560048201610487565b0390fd5b9061294593929161251b565b565b90612962929161295d6129586102bd565b6140e0565b612964565b565b91909161297961297382613118565b15610252565b612b235761299161298c6002839061077e565b612539565b9261299e600285016107b0565b936129be60056129b76129b2828501610952565b61256d565b9201610986565b946129dc6129d66129d08588906119fd565b92610581565b91610581565b03612b0757836129fd6129f76129f260096107b0565b610581565b91610581565b10612aeb57612a5d9482612a19612a135f611a22565b91610581565b11612a5f575b5050505f80612a36612a31600a61262e565b612579565b84612a3f6101f3565b9081612a4a81612599565b03925af150612a576125c1565b50614359565b565b5f80612a6a84612579565b85612a736101f3565b9081612a7e81612599565b03925af1612a8a6125c1565b50612a95575b612a1f565b839192612acb612ac57f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611f23565b93610772565b93612ae0612ad76101f3565b928392836125f7565b0390a35f8080612a90565b5f632a9ffab760e21b815280612b0360048201610487565b0390fd5b5f632a9ffab760e21b815280612b1f60048201610487565b0390fd5b5f6302e8145360e61b815280612b3b60048201610487565b0390fd5b90612b4a9291612947565b565b612b54611d09565b50612b5f60056107b0565b612b7a612b74612b6f60046107b0565b610581565b91610581565b115f14612ba157612b9e612b8e60056107b0565b612b9860046107b0565b9061253c565b90565b612baa5f611a22565b90565b606090565b67ffffffffffffffff8111612bca5760208091020190565b6108d4565b90612be1612bdc83612bb2565b611b25565b918252565b612bf1610120611b25565b90565b5f90565b5f90565b606090565b5f90565b5f90565b5f90565b5f90565b612c19612be6565b90602080808080808080808a612c2d612bf4565b815201612c38612bf4565b815201612c43612bf4565b815201612c4e612bf8565b815201612c59612bfc565b815201612c64612c01565b815201612c6f612c05565b815201612c7a612c09565b815201612c85612c0d565b81525050565b612c93612c11565b90565b5f5b828110612ca457505050565b602090612caf612c8b565b8184015201612c98565b90612cde612cc683612bcf565b92602080612cd48693612bb2565b9201910390612c96565b565b634e487b7160e01b5f52603260045260245ffd5b90612cfe82610d7f565b811015612d0f576020809102010190565b612ce0565b90612df1612de76005612d25612be6565b94612d3c612d345f83016107b0565b5f8801611b48565b612d54612d4b600183016107b0565b60208801611b48565b612d6c612d63600283016107b0565b60408801611b48565b612d84612d7b600383016107d4565b60608801611b56565b612d9c612d9360048301610911565b60808801611bda565b612db3612daa838301610952565b60a08801611bdd565b612dca612dc1838301610986565b60c08801611beb565b612de1612dd88383016109ad565b60e08801611bf9565b016109da565b6101008401611c07565b565b612dfc90612d14565b90565b612e07612bad565b50612e18612e13612b4c565b612cb9565b90612e2360046107b0565b91612e2e60056107b0565b91612e37611d09565b935b80612e4c612e4686610581565b91610581565b1015612ea857612e9c612ea291612e95612e7a612e73612e6e6003859061210e565b6107d4565b600261077e565b86612e858a92612df3565b612e8f8383612cf4565b52612cf4565b5150612124565b94612124565b93612e39565b509250905090565b90612eba90611ec1565b5f5260205260405f2090565b60ff1690565b612ed8612edd91610794565b612ec6565b90565b612eea9054612ecc565b90565b612f13915f612f08612f0e93612f01611898565b50826118e0565b01612eb0565b612ee0565b90565b60601b90565b612f2590612f16565b90565b612f3190612f1c565b90565b612f40612f459161042b565b612f28565b9052565b60c01b90565b612f5890612f49565b90565b612f67612f6c916104e9565b612f4f565b9052565b60f81b90565b612f7f90612f70565b90565b612f8e612f9391610b12565b612f76565b9052565b909182612fa781612fae93612585565b8093611b8b565b0190565b90565b612fc1612fc691610581565b612fb2565b9052565b936130199560016020999895613003600861301197612ffb60148f9c612ff38161300a9c612f34565b018092612f5b565b018092612f82565b0191612f97565b8092612fb5565b018092612fb5565b0190565b60200190565b946130549392969194613063966130386118dc565b50959793909192936130486101f3565b98899760208901612fca565b602082018103825203826108e8565b61307561306f82610a73565b9161301d565b2090565b6130929061308d6130886103bc565b6140e0565b613094565b565b806130a76130a15f611a22565b91610581565b146130f1576130b7816006611c57565b6130ec7e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db1916130e36101f3565b91829182610c63565b0390a1565b5f632a9ffab760e21b81528061310960048201610487565b0390fd5b61311690613079565b565b613120611898565b50613129612b4c565b61313b6131355f611a22565b91610581565b119081613147575b5090565b905061317861317261316c613167600361316160046107b0565b9061210e565b6107d4565b926102ec565b916102ec565b145f613143565b9c9b9a999897969594939291908d6131a66131a061319b611044565b6104e9565b916104e9565b036131b6576131b49d6131d2565b565b5f6309ff9a8360e41b8152806131ce60048201610487565b0390fd5b906131f89d9c9b9a9998979695949392916131f36131ee6102bd565b6140e0565b613770565b565b5090565b5090565b61320e61321391610794565b610cb1565b90565b6132209054613202565b90565b60209181520190565b90565b9190613249816132428161324e9561081f565b8095611b8b565b6108ca565b0190565b9061325d929161322f565b90565b5f80fd5b5f80fd5b5f80fd5b90356001602003823603038112156132ad57016020813591019167ffffffffffffffff82116132a85760018202360383136132a357565b613264565b613260565b613268565b60200190565b91816132c391613223565b90816132d46020830284019461322c565b92835f925b8484106132e95750505050505090565b909192939495602061331561330f83856001950388526133098b8861326c565b90613252565b986132b2565b9401940192949391906132d9565b60209181520190565b90565b60209181520190565b91906133528161334b816133579561332f565b8095611b8b565b6108ca565b0190565b906133669291613338565b90565b90356001602003823603038112156133aa57016020813591019167ffffffffffffffff82116133a55760018202360383136133a057565b613264565b613260565b613268565b60200190565b91816133c091613323565b90816133d16020830284019461332c565b92835f925b8484106133e65750505050505090565b909192939495602061341261340c83856001950388526134068b88613369565b9061335b565b986133af565b9401940192949391906133d6565b60209181520190565b90565b5061343b9060208101906117bc565b90565b61344790611745565b9052565b5061345a906020810190610598565b90565b9060206134886134909361347f6134765f83018361342c565b5f86019061343e565b8281019061344b565b910190610d92565b565b9061349f8160409361345d565b0190565b5090565b60400190565b916134bb826134c192613420565b92613429565b90815f905b8282106134d4575050505090565b909192936134f66134f06001926134eb88866134a3565b613492565b956134a7565b9201909291926134c6565b919061351b816135148161352095610a77565b8095611b8b565b6108ca565b0190565b99936135cf9e9c98916135ab979c9e9c6135c19b9761358161358f948f60609061357a6135b69f9a9b61357061359d9d61356661014086019b5f870190610ac9565b60208501906102ef565b60408301906102ef565b01906102ef565b8d60808185039101526132b8565b918a830360a08c01526133b5565b9187830360c08901526134ad565b9660e0850190610a66565b610100830190610a66565b610120818503910152613501565b90565b5090565b91908110156135e6576040020190565b612ce0565b356135f581610584565b90565b61360190610d10565b90565b5f80fd5b5f80fd5b5f80fd5b903590600160200381360303821215613652570180359067ffffffffffffffff821161364d5760200191600182023603831361364857565b61360c565b613608565b613604565b9082101561367257602061366e9202810190613610565b9091565b612ce0565b9035906001602003813603038212156136b9570180359067ffffffffffffffff82116136b4576020019160018202360383136136af57565b61360c565b613608565b613604565b908210156136d95760206136d59202810190613677565b9091565b612ce0565b905090565b9091826136f3816136fa936136de565b8093611b8b565b0190565b9091613709926136e3565b90565b6137206137176101f3565b928392836136fe565b03902090565b909161373d9260208301925f818503910152613501565b90565b91602061376192949361375a60408201965f8301906102ef565b01906102ef565b565b3561376d816117a8565b90565b9d9891949d9c939c9b929a9790959b999699608052610140526101005261379760016107d4565b6137b16137ab6137a65f61125a565b6102ec565b916102ec565b141580613e93575b613e77576137cf6137c98a613118565b15610252565b613e5b576137de8b8d906131fa565b60c0526137ec8a89906131fe565b6138016137fb60c05192610581565b91610581565b03613e3f5788888b8e8e956138166007613216565b61381f90610d1c565b958b8b608051958c978c999b969161014051926101005194959697986138436101f3565b9e8f9d8e9d8e6138566364c062d1611a5f565b81526004019d6138659e613524565b03815a93602094fa8015613e3a57613885915f91613e0c575b5015610252565b613df05761389d6138986002899061077e565b612539565b916138be6138b960056138b2600287016107b0565b9501610952565b61256d565b926138dc6138d66138d08789906119fd565b92610581565b91610581565b03613dd457846138fd6138f76138f260096107b0565b610581565b91610581565b10613db8575f6101205261390f611d09565b610120525f60a05261391f611d09565b60a0526139336101405161010051906135d2565b60e0525b6101205161394f61394960e051610581565b91610581565b1015613998576139806139786020613972610140516101005161012051916135d6565b016135eb565b60a0516119fd565b60a05261398f61012051612124565b61012052613937565b6139ae6139a68587906119fd565b60a0516119fd565b60a05260a0516139cf6139c96139c3306135f8565b31610581565b91610581565b11613d9c576139df888690614359565b6139e85f611a22565b610120525b61012051613a056139ff60c051610581565b91610581565b1015613aa0578a8a8a8a613a73613a37613a298d6080519495906101205191613657565b9690959061012051916136be565b919095613a6d613a677f17e8258c8124324c2b7a7f6e8d90caec0a1373c6eeffcbb90f3c3220485d437195611f23565b95610772565b9561370c565b94613a88613a7f6101f3565b92839283613726565b0390a4613a9761012051612124565b610120526139ed565b9397509397509350939650613ab6826001611c90565b608051879192613aef613ae97f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e93611f23565b93610772565b93613b04613afb6101f3565b92839283613740565b0390a381613b1a613b145f611a22565b91610581565b11613d46575b613bb2925f928392608051889192613b61613b5b7f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611f23565b93610772565b93613b76613b6d6101f3565b928392836125f7565b0390a3613b8b613b86600a61262e565b612579565b90613b946101f3565b9081613b9f81612599565b03925af1613bab6125c1565b5015610252565b613d2a57613bbf5f611a22565b610120525b61012051613bea613be4613bdf6101405161010051906135d2565b610581565b91610581565b1015613d2657613c625f80613c1c613c1782613c11610140516101005161012051916135d6565b01613763565b612579565b613c3c6020613c36610140516101005161012051916135d6565b016135eb565b613c446101f3565b9081613c4f81612599565b03925af1613c5b6125c1565b5015610252565b613d0a5760805182613c895f613c83610140516101005161012051916135d6565b01613763565b91613caa6020613ca4610140516101005161012051916135d6565b016135eb565b613cdd613cd77fccf5242d67663517e020b47096ecb69c7f1a32e85ebd295eae0611a41395b3ae93611f23565b93610772565b93613cf2613ce96101f3565b928392836125f7565b0390a3613d0161012051612124565b61012052613bc4565b5f6312171d8360e31b815280613d2260048201610487565b0390fd5b9050565b5f6312171d8360e31b815280613d4260048201610487565b0390fd5b613d7b5f80613d5484612579565b85613d5d6101f3565b9081613d6881612599565b03925af1613d746125c1565b5015610252565b15613b20575f6312171d8360e31b815280613d9860048201610487565b0390fd5b5f631e9acf1760e31b815280613db460048201610487565b0390fd5b5f632a9ffab760e21b815280613dd060048201610487565b0390fd5b5f632a9ffab760e21b815280613dec60048201610487565b0390fd5b5f638baa579f60e01b815280613e0860048201610487565b0390fd5b613e2d915060203d8111613e33575b613e2581836108e8565b810190611a88565b5f61387e565b503d613e1b565b611af2565b5f637c6953f960e01b815280613e5760048201610487565b0390fd5b5f6302e8145360e61b815280613e7360048201610487565b0390fd5b5f630b6fac0360e41b815280613e8f60048201610487565b0390fd5b5083613eb0613eaa613ea560016107d4565b6102ec565b916102ec565b14156137b9565b90613ecd9d9c9b9a99989796959493929161317f565b565b613ed7612c11565b90565b613ee2613ecf565b50613eeb6118dc565b50613ef4611898565b50613efd612b4c565b613f0f613f095f611a22565b91610581565b11613f2e57613f1c613ecf565b613f2660016107d4565b915f91929190565b613f55613f4e613f496003613f4360046107b0565b9061210e565b6107d4565b600261077e565b613f5f60016107d4565b91613f6b600192612df3565b929190565b613f8990613f84613f7f6103bc565b6140e0565b613ff7565b565b613f9f613f9a613fa492611252565b610cf1565b610420565b90565b613fb090613f8b565b90565b613fbc90612561565b90565b90565b90613fd7613fd2613fde92613fb3565b613fbf565b8254611ea6565b9055565b9190613ff5905f602085019401906125ea565b565b8061401261400c6140075f613fa7565b61042b565b91612579565b1461405d5761402281600a613fc2565b6140587fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f9161404f6101f3565b91829182613fe2565b0390a1565b5f632582a64160e11b81528061407560048201610487565b0390fd5b61408290613f70565b565b9061409f9161409a614095826118f6565b6140e0565b6140a1565b565b906140ab916141f5565b50565b906140b891614084565b565b6140c2611898565b506140dc6140d66301ffc9a760e01b610205565b91610205565b1490565b6140f2906140ec6141e8565b906143cb565b565b9061410060ff91611255565b9181191691161790565b61411390610252565b90565b90565b9061412e6141296141359261410a565b614116565b82546140f4565b9055565b614141611898565b50614156614150828490612eed565b15610252565b5f146141de5761417d60016141785f6141708186906118e0565b018590612eb0565b614119565b906141866141e8565b906141c36141bd6141b77f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610772565b92611ec1565b92611ec1565b926141cc6101f3565b806141d681610487565b0390a4600190565b50505f90565b5f90565b6141f06141e4565b503390565b6141fd611898565b50614209818390612eed565b5f146142905761422f5f61422a5f6142228186906118e0565b018590612eb0565b614119565b906142386141e8565b9061427561426f6142697ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b95610772565b92611ec1565b92611ec1565b9261427e6101f3565b8061428881610487565b0390a4600190565b50505f90565b91906142ac6142a76142b493610772565b611c84565b908354611cc2565b9055565b6142ca916142c46118dc565b91614296565b565b6142ea5f6142e560036142df60046107b0565b9061210e565b6142b8565b6143066142ff6142fa60046107b0565b612124565b6004611c57565b565b6143135f80926126a3565b0190565b90916143569361433f6143499261433560808601965f870190610a66565b6020850190612662565b6040830190612696565b6060818303910152614308565b90565b6143616142cc565b5f5f926143a36143917f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198894610772565b9461439a6101f3565b93849384614317565b0390a2565b9160206143c99294936143c260408201965f830190610abc565b01906102ef565b565b906143e06143da838390612eed565b15610252565b6143e8575050565b6144025f92839263e2517d3f60e01b8452600484016143a8565b0390fdfea264697066735822122083bf5647430e72a15db9cea51af7b9510291067f3ad854c4e40ac3dabb3e917464736f6c634300081e0033",
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
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool success)
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
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool success)
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
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool success)
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
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
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
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequests() ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequests")
}

// UnpackGetPendingRequests is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80a1f712.
//
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
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

// PackMarkRequestCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d1be609.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function markRequestCompleted(bytes32 requestId, uint256 refund, uint256 applicationFees) returns()
func (processorEndpoint *ProcessorEndpoint) PackMarkRequestCompleted(requestId [32]byte, refund *big.Int, applicationFees *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("markRequestCompleted", requestId, refund, applicationFees)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMarkRequestCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d1be609.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function markRequestCompleted(bytes32 requestId, uint256 refund, uint256 applicationFees) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackMarkRequestCompleted(requestId [32]byte, refund *big.Int, applicationFees *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("markRequestCompleted", requestId, refund, applicationFees)
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
// Solidity: function requestById(bytes32 ) view returns(uint256 timestamp, uint256 depositAmount, uint256 maxFeeValue, bytes32 requestId, bytes payload, address sender, uint64 applicationId, uint8 protocolVersion, uint8 requestType)
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
// Solidity: function requestById(bytes32 ) view returns(uint256 timestamp, uint256 depositAmount, uint256 maxFeeValue, bytes32 requestId, bytes payload, address sender, uint64 applicationId, uint8 protocolVersion, uint8 requestType)
func (processorEndpoint *ProcessorEndpoint) TryPackRequestById(arg0 [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("requestById", arg0)
}

// RequestByIdOutput serves as a container for the return parameters of contract
// method RequestById.
type RequestByIdOutput struct {
	Timestamp       *big.Int
	DepositAmount   *big.Int
	MaxFeeValue     *big.Int
	RequestId       [32]byte
	Payload         []byte
	Sender          common.Address
	ApplicationId   uint64
	ProtocolVersion uint8
	RequestType     uint8
}

// UnpackRequestById is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5423112f.
//
// Solidity: function requestById(bytes32 ) view returns(uint256 timestamp, uint256 depositAmount, uint256 maxFeeValue, bytes32 requestId, bytes payload, address sender, uint64 applicationId, uint8 protocolVersion, uint8 requestType)
func (processorEndpoint *ProcessorEndpoint) UnpackRequestById(data []byte) (RequestByIdOutput, error) {
	out, err := processorEndpoint.abi.Unpack("requestById", data)
	outstruct := new(RequestByIdOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Timestamp = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.DepositAmount = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.MaxFeeValue = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.RequestId = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Payload = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.Sender = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.ApplicationId = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.ProtocolVersion = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.RequestType = *abi.ConvertType(out[8], new(uint8)).(*uint8)
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
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea264697066735822122025e9a28501b29bc51aa9d53a9fa0e4e93bab7ea79018473599ee9698f3ebe5db64736f6c634300081e0033",
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
