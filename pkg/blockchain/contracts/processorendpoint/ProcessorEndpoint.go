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
	Value           *big.Int
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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0x54b21558.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, refundAmount, applicationFee, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54b21558.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, refundAmount, applicationFee, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54b21558.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"APPLICATION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"}],\"name\":\"markRequestCompleted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"markRequestFailed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a35c5449187bc72768eac252f2f30417e4",
	Bin: "0x6080604052346100335761001d6100146101b4565b939290926103f9565b610025610038565b6144046106db823961440490f35b61003e565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b9061006a90610042565b810190811060018060401b0382111761008257604052565b61004c565b9061009a610093610038565b9283610060565b565b5f80fd5b60018060a01b031690565b6100b4906100a0565b90565b6100c0906100ab565b90565b6100cc816100b7565b036100d357565b5f80fd5b905051906100e4826100c3565b565b6100ef906100ab565b90565b6100fb816100e6565b0361010257565b5f80fd5b90505190610113826100f2565b565b61011e816100ab565b0361012557565b5f80fd5b9050519061013682610115565b565b90565b61014481610138565b0361014b57565b5f80fd5b9050519061015c8261013b565b565b919060a0838203126101af57610176815f85016100d7565b926101848260208301610106565b926101ac6101958460408501610129565b936101a38160608601610129565b9360800161014f565b90565b61009c565b6101d2614adf803803806101c781610087565b92833981019061015e565b9091929394565b5f1b90565b906101ea5f19916101d9565b9181191691161790565b90565b90565b61020e610209610213926101f4565b6101f7565b610138565b90565b90565b9061022e610229610235926101fa565b610216565b82546101de565b9055565b90565b61025061024b61025592610239565b6101f7565b6100a0565b90565b6102619061023c565b90565b61027861027361027d926100a0565b6101f7565b6100a0565b90565b61028990610264565b90565b61029590610280565b90565b6102a190610264565b90565b6102ad90610298565b90565b5f0190565b906102c660018060a01b03916101d9565b9181191691161790565b6102d990610280565b90565b90565b906102f46102ef6102fb926102d0565b6102dc565b82546102b5565b9055565b61030890610298565b90565b90565b9061032361031e61032a926102ff565b61030b565b82546102b5565b9055565b61033790610264565b90565b6103439061032e565b90565b61034f9061032e565b90565b90565b9061036a61036561037192610346565b610352565b82546102b5565b9055565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103d16103cc6103d692610138565b6101f7565b610138565b90565b906103ee6103e96103f5926103bd565b610216565b82546101de565b9055565b91939290610409600a6006610219565b8261042c61042661042161041c5f610258565b61028c565b6100b7565b916100b7565b1480156104fe575b80156104dc575b80156104ba575b61049e5761049c946104666104869261045f6104949660076102df565b600861030e565b6104796104728261033a565b600a610355565b610481610375565b6105c9565b5061048f610399565b6105c9565b5060096103d9565b565b5f632582a64160e11b8152806104b6600482016102b0565b0390fd5b50816104d66104d06104cb5f610258565b6100ab565b916100ab565b14610442565b50846104f86104f26104ed5f610258565b6100ab565b916100ab565b1461043b565b508061052261051c6105176105125f610258565b6102a4565b6100e6565b916100e6565b14610434565b5f90565b151590565b90565b61053d90610531565b90565b9061054a90610534565b5f5260205260405f2090565b61055f90610264565b90565b61056b90610556565b90565b9061057890610562565b5f5260205260405f2090565b9061059060ff916101d9565b9181191691161790565b6105a39061052c565b90565b90565b906105be6105b96105c59261059a565b6105a6565b8254610584565b9055565b6105d1610528565b506105e66105e08284906106a0565b1561052c565b5f1461066e5761060d60016106085f610600818690610540565b01859061056e565b6105a9565b906106166106cd565b9061065361064d6106477f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610534565b92610562565b92610562565b9261065c610038565b80610666816102b0565b0390a4600190565b50505f90565b5f1c90565b60ff1690565b61068b61069091610674565b610679565b90565b61069d905461067f565b90565b6106c6915f6106bb6106c1936106b4610528565b5082610540565b0161056e565b610693565b90565b5f90565b6106d56106c9565b50339056fe60806040526004361015610013575b611a02565b61001d5f356101ec565b806301ffc9a7146101e75780631664deeb146101e2578063248a9ca3146101dd5780632a0acc6a146101d85780632f2ff15d146101d357806334f23a9d146101ce57806336568abe146101c95780633b473cb9146101c45780635423112f146101bf5780635d1be609146101ba5780635e339384146101b55780636e91d133146101b05780637f8b54a7146101ab57806380a1f712146101a657806391d14854146101a15780639588eca21461019c5780639968da96146101975780639a17995a146101925780639c73eeaf1461018d5780639c8d416f146101885780639fabc36f14610183578063a217fddf1461017e578063a9ba045e14610179578063aa3aa46014610174578063c404c3591461016f578063c415b95c1461016a578063d2c35ce814610165578063d547741f146101605763f7d9cc1e0361000e576119cd565b61198a565b611957565b611904565b61186a565b611753565b6116d6565b611633565b6115c3565b61158e565b611526565b6114cc565b611417565b611396565b61132b565b6112f6565b6110e8565b610da5565b610cd3565b610c8a565b610c11565b6107e1565b610658565b610626565b61048b565b6103ea565b610386565b610310565b610278565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f80fd5b63ffffffff60e01b1690565b61021981610204565b0361022057565b5f80fd5b9050359061023182610210565b565b9060208282031261024c57610249915f01610224565b90565b6101fc565b151590565b61025f90610251565b9052565b9190610276905f60208501940190610256565b565b346102a8576102a461029361028e366004610233565b611a0a565b61029b6101f2565b91829182610263565b0390f35b6101f8565b5f9103126102b757565b6101fc565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b6102e86102bc565b90565b90565b6102f7906102eb565b9052565b919061030e905f602085019401906102ee565b565b34610340576103203660046102ad565b61033c61032b6102e0565b6103336101f2565b918291826102fb565b0390f35b6101f8565b61034e816102eb565b0361035557565b5f80fd5b9050359061036682610345565b565b906020828203126103815761037e915f01610359565b90565b6101fc565b346103b6576103b26103a161039c366004610368565b611a64565b6103a96101f2565b918291826102fb565b0390f35b6101f8565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103e76103bb565b90565b3461041a576103fa3660046102ad565b6104166104056103df565b61040d6101f2565b918291826102fb565b0390f35b6101f8565b60018060a01b031690565b6104339061041f565b90565b61043f8161042a565b0361044657565b5f80fd5b9050359061045782610436565b565b9190604083820312610481578061047561047e925f8601610359565b9360200161044a565b90565b6101fc565b5f0190565b346104ba576104a461049e366004610459565b90611aaf565b6104ac6101f2565b806104b681610486565b0390f35b6101f8565b60ff1690565b6104ce816104bf565b036104d557565b5f80fd5b905035906104e6826104c5565b565b67ffffffffffffffff1690565b6104fe816104e8565b0361050557565b5f80fd5b90503590610516826104f5565b565b6004111561052257565b5f80fd5b9050359061053382610518565b565b5f80fd5b5f80fd5b5f80fd5b909182601f8301121561057b5781359167ffffffffffffffff831161057657602001926001830284011161057157565b61053d565b610539565b610535565b90565b61058c81610580565b0361059357565b5f80fd5b905035906105a482610583565b565b91909160c081840312610621576105bf835f83016104d9565b926105cd8160208401610509565b926105db8260408501610526565b9260608101359167ffffffffffffffff831161061c5761060084610619948401610541565b9390946106108160808601610597565b9360a001610597565b90565b610200565b6101fc565b6106546106436106373660046105a6565b959490949391936125ce565b61064b6101f2565b918291826102fb565b0390f35b346106875761067161066b366004610459565b906125e8565b6106796101f2565b8061068381610486565b0390f35b6101f8565b6011111561069657565b5f80fd5b905035906106a78261068c565b565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906106d5906106ad565b810190811067ffffffffffffffff8211176106ef57604052565b6106b7565b906107076107006101f2565b92836106cb565b565b67ffffffffffffffff8111610727576107236020916106ad565b0190565b6106b7565b90825f939282370152565b9092919261074c61074782610709565b6106f4565b93818552602085019082840111610768576107669261072c565b565b6106a9565b9080601f8301121561078b5781602061078893359101610737565b90565b610535565b916060838303126107dc576107a7825f8501610359565b926107b5836020830161069a565b92604082013567ffffffffffffffff81116107d7576107d4920161076d565b90565b610200565b6101fc565b34610810576107fa6107f4366004610790565b91612a78565b6108026101f2565b8061080c81610486565b0390f35b6101f8565b61081e906102eb565b90565b9061082b90610815565b5f5260205260405f2090565b5f1c90565b60ff1690565b61084e61085391610837565b61083c565b90565b6108609054610842565b90565b60081c90565b67ffffffffffffffff1690565b61088261088791610863565b610869565b90565b6108949054610876565b90565b60481c90565b60ff1690565b6108af6108b491610897565b61089d565b90565b6108c190546108a3565b90565b90565b6108d36108d891610837565b6108c4565b90565b6108e590546108c7565b90565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561091c575b602083101461091757565b6108e8565b91607f169161090c565b60209181520190565b5f5260205f2090565b905f929180549061095261094b836108fc565b8094610926565b916001811690815f146109a9575060011461096d575b505050565b61097a919293945061092f565b915f925b81841061099157505001905f8080610968565b6001816020929593955484860152019101929061097e565b92949550505060ff19168252151560200201905f8080610968565b906109ce91610938565b90565b906109f16109ea926109e16101f2565b938480926109c4565b03836106cb565b565b90565b610a02610a0791610837565b6109f3565b90565b610a1490546109f6565b90565b60018060a01b031690565b610a2e610a3391610837565b610a17565b90565b610a409054610a22565b90565b610a4e906002610821565b610a595f8201610856565b91610a655f830161088a565b91610a715f82016108b7565b91610a7e600183016108db565b91610a8b600282016109d1565b91610a9860038301610a0a565b91610aa560048201610a36565b91610abe6006610ab760058501610a0a565b9301610a0a565b90565b610aca906104bf565b9052565b610ad7906104e8565b9052565b634e487b7160e01b5f52602160045260245ffd5b60041115610af957565b610adb565b90610b0882610aef565b565b610b1390610afe565b90565b610b1f90610b0a565b9052565b5190565b60209181520190565b90825f9392825e0152565b610b5a610b63602093610b6893610b5181610b23565b93848093610b27565b95869101610b30565b6106ad565b0190565b610b7590610580565b9052565b610b829061042a565b9052565b94610be9610100979b9a9892610bf492610bdc610c0898610bd2610c0f9e99610bc8610bfe9a60208f610bc161012082019a5f830190610ac1565b0190610ace565b60408d0190610b16565b60608b01906102ee565b88820360808a0152610b3b565b9a60a0870190610b6c565b60c0850190610b79565b60e0830190610b6c565b0190610b6c565b565b34610c4b57610c47610c2c610c27366004610368565b610a43565b95610c3e9997999591959492946101f2565b998a998a610b86565b0390f35b6101f8565b9091606082840312610c8557610c82610c6b845f8501610359565b93610c798160208601610597565b93604001610597565b90565b6101fc565b34610cb957610ca3610c9d366004610c50565b91612d7e565b610cab6101f2565b80610cb581610486565b0390f35b6101f8565b9190610cd1905f60208501940190610b6c565b565b34610d0357610ce33660046102ad565b610cff610cee612d8b565b610cf66101f2565b91829182610cbe565b0390f35b6101f8565b1c90565b60018060a01b031690565b610d27906008610d2c9302610d08565b610d0c565b90565b90610d3a9154610d17565b90565b610d4960075f90610d2f565b90565b90565b610d63610d5e610d689261041f565b610d4c565b61041f565b90565b610d7490610d4f565b90565b610d8090610d6b565b90565b610d8c90610d77565b9052565b9190610da3905f60208501940190610d83565b565b34610dd557610db53660046102ad565b610dd1610dc0610d3d565b610dc86101f2565b91829182610d90565b0390f35b6101f8565b67ffffffffffffffff8111610df25760208091020190565b6106b7565b67ffffffffffffffff8111610e1557610e116020916106ad565b0190565b6106b7565b90929192610e2f610e2a82610df7565b6106f4565b93818552602085019082840111610e4b57610e499261072c565b565b6106a9565b9080601f83011215610e6e57816020610e6b93359101610e1a565b90565b610535565b929190610e87610e8282610dda565b6106f4565b9381855260208086019202810191838311610ede5781905b838210610ead575050505050565b813567ffffffffffffffff8111610ed957602091610ece8784938701610e50565b815201910190610e9f565b610535565b61053d565b9080601f83011215610f0157816020610efe93359101610e73565b90565b610535565b67ffffffffffffffff8111610f1e5760208091020190565b6106b7565b5f80fd5b610f309061041f565b90565b610f3c81610f27565b03610f4357565b5f80fd5b90503590610f5482610f33565b565b9190604083820312610f9057610f8990610f7060406106f4565b93610f7d825f8301610f47565b5f860152602001610597565b6020830152565b610f23565b90929192610faa610fa582610f06565b6106f4565b938185526040602086019202830192818411610fe957915b838310610fcf5750505050565b6020604091610fde8486610f56565b815201920191610fc2565b61053d565b9080601f8301121561100c5781602061100993359101610f95565b90565b610535565b91610120838303126110e357611029825f8501610509565b926110378360208301610359565b926110458160408401610359565b926110538260608501610359565b92608081013567ffffffffffffffff81116110de5783611074918301610ee3565b9260a082013567ffffffffffffffff81116110d95781611095918401610fee565b926110a38260c08501610597565b926110b18360e08301610597565b9261010082013567ffffffffffffffff81116110d4576110d19201610e50565b90565b610200565b610200565b610200565b6101fc565b346111205761110a6110fb366004611011565b97969096959195949294613791565b6111126101f2565b8061111c81610486565b0390f35b6101f8565b5190565b60209181520190565b60200190565b611141906104bf565b9052565b61114e906104e8565b9052565b61115b90610b0a565b9052565b611168906102eb565b9052565b61118b6111946020936111999361118281610b23565b93848093610926565b95869101610b30565b6106ad565b0190565b6111a690610580565b9052565b6111b39061042a565b9052565b90611261906101008061122061012084016111d85f8801515f870190611138565b6111ea60208801516020870190611145565b6111fc60408801516040870190611152565b61120e6060880151606087019061115f565b6080870151858203608087015261116c565b9461123360a082015160a086019061119d565b61124560c082015160c08601906111aa565b61125760e082015160e086019061119d565b015191019061119d565b90565b9061126e916111b7565b90565b60200190565b9061128b61128483611125565b8092611129565b908161129c60208302840194611132565b925f915b8383106112af57505050505090565b909192939460206112d16112cb83856001950387528951611264565b97611271565b93019301919392906112a0565b6112f39160208201915f818403910152611277565b90565b34611326576113063660046102ad565b6113226113116138e9565b6113196101f2565b918291826112de565b0390f35b6101f8565b3461135c57611358611347611341366004610459565b906139d4565b61134f6101f2565b91829182610263565b0390f35b6101f8565b6113719060086113769302610d08565b6108c4565b90565b906113849154611361565b90565b61139360015f90611379565b90565b346113c6576113a63660046102ad565b6113c26113b1611387565b6113b96101f2565b918291826102fb565b0390f35b6101f8565b90565b6113e26113dd6113e7926113cb565b610d4c565b6104e8565b90565b6113f460016113ce565b90565b6113ff6113ea565b90565b9190611415905f60208501940190610ace565b565b34611447576114273660046102ad565b6114436114326113f7565b61143a6101f2565b91829182611402565b0390f35b6101f8565b91909160c0818403126114c757611465835f830161044a565b926114738160208401610509565b926114818260408501610526565b9260608101359167ffffffffffffffff83116114c2576114a6846114bf948401610541565b9390946114b68160808601610597565b9360a001610597565b90565b610200565b6101fc565b34611503576114ff6114ee6114e236600461144c565b95949094939193613b0a565b6114f66101f2565b918291826102fb565b0390f35b6101f8565b906020828203126115215761151e915f01610597565b90565b6101fc565b346115545761153e611539366004611508565b613bf4565b6115466101f2565b8061155081610486565b0390f35b6101f8565b61156990600861156e9302610d08565b6109f3565b90565b9061157c9154611559565b90565b61158b60095f90611571565b90565b346115be5761159e3660046102ad565b6115ba6115a961157f565b6115b16101f2565b91829182610cbe565b0390f35b6101f8565b346115f3576115ef6115de6115d9366004610368565b613bff565b6115e66101f2565b91829182610263565b0390f35b6101f8565b90565b5f1b90565b61161461160f611619926115f8565b6115fb565b6102eb565b90565b6116255f611600565b90565b61163061161c565b90565b34611663576116433660046102ad565b61165f61164e611628565b6116566101f2565b918291826102fb565b0390f35b6101f8565b60018060a01b031690565b6116839060086116889302610d08565b611668565b90565b906116969154611673565b90565b6116a560085f9061168b565b90565b6116b190610d6b565b90565b6116bd906116a8565b9052565b91906116d4905f602085019401906116b4565b565b34611706576116e63660046102ad565b6117026116f1611699565b6116f96101f2565b918291826116c1565b0390f35b6101f8565b61171f61171a611724926115f8565b610d4c565b6104bf565b90565b6117305f61170b565b90565b61173b611727565b90565b9190611751905f60208501940190610ac1565b565b34611783576117633660046102ad565b61177f61176e611733565b6117766101f2565b9182918261173e565b0390f35b6101f8565b9061183290610100806117f161012084016117a95f8801515f870190611138565b6117bb60208801516020870190611145565b6117cd60408801516040870190611152565b6117df6060880151606087019061115f565b6080870151858203608087015261116c565b9461180460a082015160a086019061119d565b61181660c082015160c08601906111aa565b61182860e082015160e086019061119d565b015191019061119d565b90565b6040906118616118566118689597969460608401908482035f860152611788565b9660208301906102ee565b0190610256565b565b3461189d5761187a3660046102ad565b611899611885613c71565b6118909391936101f2565b93849384611835565b0390f35b6101f8565b60018060a01b031690565b6118bd9060086118c29302610d08565b6118a2565b90565b906118d091546118ad565b90565b6118df600a5f906118c5565b90565b6118eb90610f27565b9052565b9190611902905f602085019401906118e2565b565b34611934576119143660046102ad565b61193061191f6118d3565b6119276101f2565b918291826118ef565b0390f35b6101f8565b906020828203126119525761194f915f01610f47565b90565b6101fc565b346119855761196f61196a366004611939565b613e10565b6119776101f2565b8061198181610486565b0390f35b6101f8565b346119b9576119a361199d366004610459565b90613e45565b6119ab6101f2565b806119b581610486565b0390f35b6101f8565b6119ca60065f90611571565b90565b346119fd576119dd3660046102ad565b6119f96119e86119be565b6119f06101f2565b91829182610cbe565b0390f35b6101f8565b5f80fd5b5f90565b611a12611a06565b5080611a2d611a27637965db0b60e01b610204565b91610204565b14908115611a3a575b5090565b611a449150613e51565b5f611a36565b5f90565b90611a5890610815565b5f5260205260405f2090565b6001611a7c611a8292611a75611a4a565b505f611a4e565b016108db565b90565b90611aa091611a9b611a9682611a64565b613e77565b611aa2565b565b90611aac91613eba565b50565b90611ab991611a85565b565b9695949392919080611adc611ad6611ad1611727565b6104bf565b916104bf565b03611aed57611aea97611b09565b90565b5f635428eae760e01b815280611b0560048201610486565b0390fd5b9695949392919081611b2a611b24611b1f6113ea565b6104e8565b916104e8565b03611b3b57611b3897612246565b90565b5f6309ff9a8360e41b815280611b5360048201610486565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b611b7a611b8091939293610580565b92610580565b8201809211611b8b57565b611b57565b611ba4611b9f611ba9926115f8565b610d4c565b610580565b90565b611bb8611bbd91610837565b611668565b90565b611bca9054611bac565b90565b60e01b90565b611bdc81610251565b03611be357565b5f80fd5b90505190611bf482611bd3565b565b90602082820312611c0f57611c0c915f01611be7565b90565b6101fc565b611c28611c23611c2d926104e8565b610d4c565b610580565b90565b611c3990611c14565b9052565b916020611c5e929493611c5760408201965f830190611c30565b0190610b79565b565b611c686101f2565b3d5f823e3d90fd5b5090565b90565b611c8b611c86611c9092611c74565b610d4c565b610580565b90565b611c9e6101206106f4565b90565b90611cab906104bf565b9052565b90611cb9906104e8565b9052565b90611cc790610afe565b9052565b90611cd5906102eb565b9052565b611ce4913691610e1a565b90565b52565b90611cf490610580565b9052565b90611d029061042a565b9052565b634e487b7160e01b5f525f60045260245ffd5b611d2390516104bf565b90565b90611d3260ff916115fb565b9181191691161790565b611d50611d4b611d55926104bf565b610d4c565b6104bf565b90565b90565b90611d70611d6b611d7792611d3c565b611d58565b8254611d26565b9055565b611d8590516104e8565b90565b60081b90565b90611da268ffffffffffffffff0091611d88565b9181191691161790565b611dc0611dbb611dc5926104e8565b610d4c565b6104e8565b90565b90565b90611de0611ddb611de792611dac565b611dc8565b8254611d8e565b9055565b611df59051610afe565b90565b60481b90565b90611e1369ff00000000000000000091611df8565b9181191691161790565b611e2690610afe565b90565b90565b90611e41611e3c611e4892611e1d565b611e29565b8254611dfe565b9055565b611e5690516102eb565b90565b90611e655f19916115fb565b9181191691161790565b611e7890610837565b90565b90611e90611e8b611e9792610815565b611e6f565b8254611e59565b9055565b5190565b601f602091010490565b1b90565b91906008611ec8910291611ec25f1984611ea9565b92611ea9565b9181191691161790565b611ee6611ee1611eeb92610580565b610d4c565b610580565b90565b90565b9190611f07611f02611f0f93611ed2565b611eee565b908354611ead565b9055565b5f90565b611f2991611f23611f13565b91611ef1565b565b5b818110611f37575050565b80611f445f600193611f17565b01611f2c565b9190601f8111611f5a575b505050565b611f66611f8b9361092f565b906020611f7284611e9f565b83019310611f93575b611f8490611e9f565b0190611f2b565b5f8080611f55565b9150611f8481929050611f7b565b90611fb1905f1990600802610d08565b191690565b81611fc091611fa1565b906002021790565b90611fd281610b23565b9067ffffffffffffffff821161209257611ff682611ff085546108fc565b85611f4a565b602090601f831160011461202a57918091612019935f9261201e575b5050611fb6565b90555b565b90915001515f80612012565b601f198316916120398561092f565b925f5b81811061207a57509160029391856001969410612060575b5050500201905561201c565b612070910151601f841690611fa1565b90555f8080612054565b9193602060018192878701518155019501920161203c565b6106b7565b906120a191611fc8565b565b6120ad9051610580565b90565b906120c56120c06120cc92611ed2565b611eee565b8254611e59565b9055565b6120da905161042a565b90565b906120ee60018060a01b03916115fb565b9181191691161790565b61210190610d6b565b90565b90565b9061211c612117612123926120f8565b612104565b82546120dd565b9055565b9061220061010060066122069461214b5f82016121455f8801611d19565b90611d5b565b6121635f820161215d60208801611d7b565b90611dcb565b61217b5f820161217560408801611deb565b90611e2c565b6121946001820161218e60608801611e4c565b90611e7b565b6121ad600282016121a760808801611e9b565b90612097565b6121c6600382016121c060a088016120a3565b906120b0565b6121df600482016121d960c088016120d0565b90612107565b6121f8600582016121f260e088016120a3565b906120b0565b0192016120a3565b906120b0565b565b9061221291612127565b565b9061221e90611ed2565b5f5260205260405f2090565b61223390610580565b5f1981146122415760010190565b611b57565b9693919592949096503461226c612266612261868890611b6b565b610580565b91610580565b036125b2578361228d6122876122826009610a0a565b610580565b91610580565b106125965761229a612d8b565b6122b56122af6122aa6006610a0a565b610580565b91610580565b101561257a57846122cf6122c96003610afe565b91610afe565b145f14612474576122e1828290611c70565b6122f46122ee6085611c77565b91610580565b03612458575b338686848490879261230c6005610a0a565b9461231696613b0a565b969594918790919042933395969761232c611c93565b995f8b019061233a91611ca1565b60208a019061234891611caf565b604089019061235691611cbd565b606088019061236491611ccb565b61236d91611cd9565b608086019061237b91611ce7565b60a085019061238991611cea565b60c084019061239791611cf8565b60e08301906123a591611cea565b6101008201906123b491611cea565b6002826123c091610821565b906123ca91612208565b8060036123d76005610a0a565b6123e091612214565b906123ea91611e7b565b6123f46005610a0a565b6123fd9061222a565b6124089060056120b0565b80337f73394f43049193ecd2e0c22eefa0ecf10987ce9c72da906a82a3336dc2ac05479161243590610815565b9061243f906120f8565b916124486101f2565b8061245281610486565b0390a390565b5f637c6953f960e01b81528061247060048201610486565b0390fd5b846124886124826002610afe565b91610afe565b14612493575b6122fa565b826124a66124a05f611b90565b91610580565b0361255e576124bd6124b86008611bc0565b6116a8565b602063116ad4a69188906124e333946124ee6124d76101f2565b96879586948594611bcd565b845260048401611c3d565b03915afa80156125595761250a915f9161252b575b5015610251565b1561248e575f63622a850b60e11b81528061252760048201610486565b0390fd5b61254c915060203d8111612552575b61254481836106cb565b810190611bf6565b5f612503565b503d61253a565b611c60565b5f632a9ffab760e21b81528061257660048201610486565b0390fd5b5f63d8219b3d60e01b81528061259260048201610486565b0390fd5b5f63c77f97eb60e01b8152806125ae60048201610486565b0390fd5b5f632a9ffab760e21b8152806125ca60048201610486565b0390fd5b906125e59695949392916125e0611a4a565b611abb565b90565b90806126036125fd6125f8613f69565b61042a565b9161042a565b036126145761261191613f76565b50565b5f63334bd91960e11b81528061262c60048201610486565b0390fd5b9061264b92916126466126416102bc565b613e77565b612828565b565b61265c61266291939293610580565b92610580565b820391821161266d57565b611b57565b61267b90610d4f565b90565b61268790612672565b90565b61269390610d6b565b90565b905090565b6126a65f8092612696565b0190565b6126b39061269b565b90565b906126c86126c383610df7565b6106f4565b918252565b606090565b3d5f146126ed576126e23d6126b6565b903d5f602084013e5b565b6126f56126cd565b906126eb565b91602061271c92949361271560408201965f830190610b79565b0190610b6c565b565b61272a61272f91610837565b6118a2565b90565b61273c905461271e565b90565b6003111561274957565b610adb565b906127588261273f565b565b6127639061274e565b90565b61276f9061275a565b9052565b6011111561277d57565b610adb565b9061278c82612773565b565b61279790612782565b90565b6127a39061278e565b9052565b5190565b60209181520190565b6127d36127dc6020936127e1936127ca816127a7565b938480936127ab565b95869101610b30565b6106ad565b0190565b90926128189061280e612825969461280460808601975f870190610b6c565b6020850190612766565b604083019061279a565b60608184039101526127b4565b90565b61283a61283482613bff565b15610251565b612a5c57612855600461284f60028490610821565b01610a36565b6128a961286f600561286960028690610821565b01610a0a565b6128a3612889600661288360028890610821565b01610a0a565b91612892614263565b9161289d6009610a0a565b9061264d565b90611b6b565b90816128bd6128b75f611b90565b91610580565b116129b3575b50505f806128d96128d4600a612732565b61268a565b6128e36009610a0a565b6128eb6101f2565b90816128f6816126aa565b03925af16129026126d2565b505f1461295d57906129146009610a0a565b612957600192946129457f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198895610815565b9561294e6101f2565b948594856127e5565b0390a25b565b906129686009610a0a565b6129ab600292946129997f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198895610815565b956129a26101f2565b948594856127e5565b0390a261295b565b5f806129c66129c18461267e565b61268a565b846129cf6101f2565b90816129da816126aa565b03925af16129e66126d2565b506129f1575b6128c3565b612a075f612a0160028690610821565b0161088a565b839192612a3d612a377f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611dac565b93610815565b93612a52612a496101f2565b928392836126fb565b0390a35f806129ec565b5f6302e8145360e61b815280612a7460048201610486565b0390fd5b90612a839291612630565b565b90612aa09291612a9b612a966102bc565b613e77565b612b9c565b565b612aad6101206106f4565b90565b90612b8e612b846006612ac1612aa2565b94612ad8612ad05f8301610856565b5f8801611ca1565b612aef612ae65f830161088a565b60208801611caf565b612b06612afd5f83016108b7565b60408801611cbd565b612b1e612b15600183016108db565b60608801611ccb565b612b36612b2d600283016109d1565b60808801611ce7565b612b4e612b4560038301610a0a565b60a08801611cea565b612b66612b5d60048301610a36565b60c08801611cf8565b612b7e612b7560058301610a0a565b60e08801611cea565b01610a0a565b6101008401611cea565b565b612b9990612ab0565b90565b919091612bb1612bab82613bff565b15610251565b612d6257612bc9612bc460028390610821565b612b90565b92612bd5818490611b6b565b612bf3612bed612be861010088016120a3565b610580565b91610580565b03612d465782612c14612c0e612c096009610a0a565b610580565b91610580565b10612d2a57612c739381612c30612c2a5f611b90565b91610580565b11612c75575b50505f80612c4c612c47600a612732565b61268a565b84612c556101f2565b9081612c60816126aa565b03925af150612c6d6126d2565b50614321565b565b5f80612c93612c8e612c8960c086016120d0565b61267e565b61268a565b84612c9c6101f2565b9081612ca7816126aa565b03925af1612cb36126d2565b50612cbe575b612c36565b612cca60208201611d7b565b612cd760c08593016120d0565b92612d0b612d057f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611dac565b93610815565b93612d20612d176101f2565b928392836126fb565b0390a35f80612cb9565b5f632a9ffab760e21b815280612d4260048201610486565b0390fd5b5f632a9ffab760e21b815280612d5e60048201610486565b0390fd5b5f6302e8145360e61b815280612d7a60048201610486565b0390fd5b90612d899291612a85565b565b612d93611f13565b50612d9e6005610a0a565b612db9612db3612dae6004610a0a565b610580565b91610580565b115f14612de057612ddd612dcd6005610a0a565b612dd76004610a0a565b9061264d565b90565b612de95f611b90565b90565b979695949392919088612e0e612e08612e036113ea565b6104e8565b916104e8565b03612e1e57612e1c98612e3a565b565b5f6309ff9a8360e41b815280612e3660048201610486565b0390fd5b90612e5b9897969594939291612e56612e516102bc565b613e77565b613130565b565b612e69612e6e91610837565b610d0c565b90565b612e7b9054612e5d565b90565b5190565b60209181520190565b60200190565b90612e9b9161116c565b90565b60200190565b90612eb8612eb183612e7e565b8092612e82565b9081612ec960208302840194612e8b565b925f915b838310612edc57505050505090565b90919293946020612efe612ef883856001950387528951612e91565b97612e9e565b9301930191939290612ecd565b5190565b60209181520190565b60200190565b612f2790610f27565b9052565b90602080612f4d93612f435f8201515f860190612f1e565b015191019061119d565b565b90612f5c81604093612f2b565b0190565b60200190565b90612f83612f7d612f7684612f0b565b8093612f0f565b92612f18565b905f5b818110612f935750505090565b909192612fac612fa66001928651612f4f565b94612f60565b9101919091612f86565b959061301590613038966130086130469c9a9b97612ffe61302e98612ff46130239960208f612fed61012082019a5f830190610ace565b01906102ee565b60408d01906102ee565b60608b01906102ee565b88820360808a0152612ea4565b9086820360a0880152612f66565b9560c0850190610b6c565b60e0830190610b6c565b610100818403910152610b3b565b90565b634e487b7160e01b5f52603260045260245ffd5b9061306782612f0b565b811015613078576020809102010190565b613049565b60016130899101610580565b90565b61309590610d6b565b90565b906130a282612e7e565b8110156130b3576020809102010190565b613049565b6130cd9160208201915f818403910152610b3b565b90565b9160206130f19294936130ea60408201965f8301906102ee565b01906102ee565b565b6130fd9051610f27565b90565b6131099061268a565b9052565b91602061312e92949361312760408201965f830190613100565b0190610b6c565b565b959798919690949893929361314560016108db565b61315f6131596131545f611600565b6102eb565b916102eb565b14158061376d575b6137515761317d61317786613bff565b15610251565b613735576020868b6131ca8a6131d58a8f978f988c906131a56131a06007612e71565b610d77565b976354b2155896999b94929091928d94956131be6101f2565b9d8e9c8d9b8c9b611bcd565b8b5260048b01612fb6565b03915afa8015613730576131f1915f91613702575b5015610251565b6136e65761320961320460028690610821565b612b90565b97613215818390611b6b565b61323361322d6132286101008d016120a3565b610580565b91610580565b036136ca578161325461324e6132496009610a0a565b610580565b91610580565b106136ae575f9998979950613267611f13565b97613270611f13565b9a5b8a61328e6132886132838d93612f0b565b610580565b91610580565b10156132d9578a61328e6132886132cb839f6132839f908f60206132b86132bf926132c59561305d565b51016120a3565b90611b6b565b9d61307d565b9d505050509a99989a613272565b91939699959850919396996132f9906132f3858790611b6b565b90611b6b565b61331461330e6133083061308c565b31610580565b91610580565b1161369257613324888590614321565b61332d5f611b90565b5b8061334961334361333e8b612e7e565b610580565b91610580565b10156133b7576133b2908b8a6133608b8490613098565b51916133aa6133986133927f09d7e6b35973ce70ef4c7f3535a0d8e74046743324fc7f02809f0f458147e73193611dac565b93610815565b936133a16101f2565b918291826130b8565b0390a361307d565b61332e565b5091939695509193966133cb826001611e7b565b888691926134026133fc7f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e93611dac565b93610815565b9361341761340e6101f2565b928392836130d0565b0390a38161342d6134275f611b90565b91610580565b11613629575b6134ce925f9283928961344960c08993016120d0565b9261347d6134777f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611dac565b93610815565b936134926134896101f2565b928392836126fb565b0390a36134a76134a2600a612732565b61268a565b906134b06101f2565b90816134bb816126aa565b03925af16134c76126d2565b5015610251565b61360d576134db5f611b90565b5b806134f76134f16134ec86612f0b565b610580565b91610580565b10156136065761355d5f8061352061351b8261351489889061305d565b51016130f3565b61268a565b613537602061353089889061305d565b51016120a3565b61353f6101f2565b908161354a816126aa565b03925af16135566126d2565b5015610251565b6135ea576135e590858361357d5f61357688869061305d565b51016130f3565b91613595602061358e89879061305d565b51016120a3565b6135c86135c27fccf5242d67663517e020b47096ecb69c7f1a32e85ebd295eae0611a41395b3ae93611dac565b93610815565b936135dd6135d46101f2565b9283928361310d565b0390a361307d565b6134dc565b5f6312171d8360e31b81528061360260048201610486565b0390fd5b5050509050565b5f6312171d8360e31b81528061362560048201610486565b0390fd5b6136715f8061364a61364561364060c087016120d0565b61267e565b61268a565b856136536101f2565b908161365e816126aa565b03925af161366a6126d2565b5015610251565b15613433575f6312171d8360e31b81528061368e60048201610486565b0390fd5b5f631e9acf1760e31b8152806136aa60048201610486565b0390fd5b5f632a9ffab760e21b8152806136c660048201610486565b0390fd5b5f632a9ffab760e21b8152806136e260048201610486565b0390fd5b5f638baa579f60e01b8152806136fe60048201610486565b0390fd5b613723915060203d8111613729575b61371b81836106cb565b810190611bf6565b5f6131ea565b503d613711565b611c60565b5f6302e8145360e61b81528061374d60048201610486565b0390fd5b5f630b6fac0360e41b81528061376960048201610486565b0390fd5b508561378a61378461377f60016108db565b6102eb565b916102eb565b1415613167565b906137a29897969594939291612dec565b565b606090565b67ffffffffffffffff81116137c15760208091020190565b6106b7565b906137d86137d3836137a9565b6106f4565b918252565b5f90565b5f90565b5f90565b5f90565b606090565b5f90565b5f90565b613802612aa2565b90602080808080808080808a6138166137dd565b8152016138216137e1565b81520161382c6137e5565b8152016138376137e9565b8152016138426137ed565b81520161384d6137f2565b8152016138586137f6565b8152016138636137f2565b81520161386e6137f2565b81525050565b61387c6137fa565b90565b5f5b82811061388d57505050565b602090613898613874565b8184015201613881565b906138c76138af836137c6565b926020806138bd86936137a9565b920191039061387f565b565b906138d382611125565b8110156138e4576020809102010190565b613049565b6138f16137a4565b506139026138fd612d8b565b6138a2565b9061390d6004610a0a565b91613916611f13565b925b8061393461392e6139296005610a0a565b610580565b91610580565b10156139905761398461398a9161397d61396261395b61395660038590612214565b6108db565b6002610821565b8561396d8992612b90565b61397783836138c9565b526138c9565b515061307d565b9361307d565b92613918565b5090915090565b906139a1906120f8565b5f5260205260405f2090565b60ff1690565b6139bf6139c491610837565b6139ad565b90565b6139d190546139b3565b90565b6139fa915f6139ef6139f5936139e8611a06565b5082611a4e565b01613997565b6139c7565b90565b60601b90565b613a0c906139fd565b90565b613a1890613a03565b90565b613a27613a2c9161042a565b613a0f565b9052565b60c01b90565b613a3f90613a30565b90565b613a4e613a53916104e8565b613a36565b9052565b60f81b90565b613a6690613a57565b90565b613a75613a7a91610b0a565b613a5d565b9052565b909182613a8e81613a9593612696565b809361072c565b0190565b90565b613aa8613aad91610580565b613a99565b9052565b93613b009560016020999895613aea6008613af897613ae260148f9c613ada81613af19c613a1b565b018092613a42565b018092613a69565b0191613a7e565b8092613a9c565b018092613a9c565b0190565b60200190565b94613b3b9392969194613b4a96613b1f611a4a565b5095979390919293613b2f6101f2565b98899760208901613ab1565b602082018103825203826106cb565b613b5c613b5682610b23565b91613b04565b2090565b613b7990613b74613b6f6103bb565b613e77565b613b7b565b565b80613b8e613b885f611b90565b91610580565b14613bd857613b9e8160066120b0565b613bd37e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db191613bca6101f2565b91829182610cbe565b0390a1565b5f632a9ffab760e21b815280613bf060048201610486565b0390fd5b613bfd90613b60565b565b613c07611a06565b50613c10612d8b565b613c22613c1c5f611b90565b91610580565b119081613c2e575b5090565b9050613c5f613c59613c53613c4e6003613c486004610a0a565b90612214565b6108db565b926102eb565b916102eb565b145f613c2a565b613c6e6137fa565b90565b613c79613c66565b50613c82611a4a565b50613c8b611a06565b50613c94612d8b565b613ca6613ca05f611b90565b91610580565b11613cc557613cb3613c66565b613cbd60016108db565b915f91929190565b613cec613ce5613ce06003613cda6004610a0a565b90612214565b6108db565b6002610821565b613cf660016108db565b91613d02600192612b90565b929190565b613d2090613d1b613d166103bb565b613e77565b613d8e565b565b613d36613d31613d3b926115f8565b610d4c565b61041f565b90565b613d4790613d22565b90565b613d5390612672565b90565b90565b90613d6e613d69613d7592613d4a565b613d56565b82546120dd565b9055565b9190613d8c905f60208501940190613100565b565b80613da9613da3613d9e5f613d3e565b61042a565b9161268a565b14613df457613db981600a613d59565b613def7fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f91613de66101f2565b91829182613d79565b0390a1565b5f632582a64160e11b815280613e0c60048201610486565b0390fd5b613e1990613d07565b565b90613e3691613e31613e2c82611a64565b613e77565b613e38565b565b90613e4291613f76565b50565b90613e4f91613e1b565b565b613e59611a06565b50613e73613e6d6301ffc9a760e01b610204565b91610204565b1490565b613e8990613e83613f69565b90614393565b565b613e9490610251565b90565b90565b90613eaf613eaa613eb692613e8b565b613e97565b8254611d26565b9055565b613ec2611a06565b50613ed7613ed18284906139d4565b15610251565b5f14613f5f57613efe6001613ef95f613ef1818690611a4e565b018590613997565b613e9a565b90613f07613f69565b90613f44613f3e613f387f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610815565b926120f8565b926120f8565b92613f4d6101f2565b80613f5781610486565b0390a4600190565b50505f90565b5f90565b613f71613f65565b503390565b613f7e611a06565b50613f8a8183906139d4565b5f1461401157613fb05f613fab5f613fa3818690611a4e565b018590613997565b613e9a565b90613fb9613f69565b90613ff6613ff0613fea7ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b95610815565b926120f8565b926120f8565b92613fff6101f2565b8061400981610486565b0390a4600190565b50505f90565b919061402d61402861403593610815565b611e6f565b908354611ead565b9055565b61404b91614045611a4a565b91614017565b565b90614060905f1990602003600802610d08565b8154169055565b905f9161407e6140768261092f565b928354611fb6565b905555565b919290602082105f146140dc57601f84116001146140ac576140a6929350611fb6565b90555b5b565b50906140d26140d79360016140c96140c38561092f565b92611e9f565b82019101611f2b565b614067565b6140a9565b5061411382936140ed60019461092f565b61410c6140f985611e9f565b820192601f86168061411e575b50611e9f565b0190611f2b565b6002021790556140aa565b61412a9088860361404d565b5f614106565b929091680100000000000000008211614190576020115f1461418157602081105f146141655761415f91611fb6565b90555b5b565b60019160ff19166141758461092f565b55600202019055614162565b60019150600202019055614163565b6106b7565b9081546141a1816108fc565b908183116141ca575b8183106141b8575b50505050565b6141c193614083565b5f8080806141b2565b6141d683838387614130565b6141aa565b5f6141e591614195565b565b905f036141f9576141f7906141db565b565b611d06565b5f600661424a9282808201556142178360018301614039565b61422483600283016141e7565b6142318360038301611f17565b8260048201556142448360058301611f17565b01611f17565b565b905f0361425e5761425c906141fe565b565b611d06565b6142945f61428f6002614289614284600361427e6004610a0a565b90612214565b6108db565b90610821565b61424c565b6142b25f6142ad60036142a76004610a0a565b90612214565b614039565b6142ce6142c76142c26004610a0a565b61222a565b60046120b0565b565b6142db5f80926127ab565b0190565b909161431e93614307614311926142fd60808601965f870190610b6c565b6020850190612766565b604083019061279a565b60608183039101526142d0565b90565b614329614263565b5f5f9261436b6143597f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198894610815565b946143626101f2565b938493846142df565b0390a2565b91602061439192949361438a60408201965f830190610b79565b01906102ee565b565b906143a86143a28383906139d4565b15610251565b6143b0575050565b6143ca5f92839263e2517d3f60e01b845260048401614370565b0390fdfea26469706673582212208e1e898002579ab6c5eb9ce23bf2f8a690cac603e825c10fec299f7d983dcec764736f6c634300081e0033",
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
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 value, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, value *big.Int, idx *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, value, idx)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a17995a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 value, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, value *big.Int, idx *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, value, idx)
}

// UnpackGenerateRequestId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a17995a.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 value, uint256 idx) pure returns(bytes32)
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
// Solidity: function requestById(bytes32 ) view returns(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes32 requestId, bytes payload, uint256 timestamp, address sender, uint256 value, uint256 maxFeeValue)
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
// Solidity: function requestById(bytes32 ) view returns(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes32 requestId, bytes payload, uint256 timestamp, address sender, uint256 value, uint256 maxFeeValue)
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
	Value           *big.Int
	MaxFeeValue     *big.Int
}

// UnpackRequestById is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5423112f.
//
// Solidity: function requestById(bytes32 ) view returns(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes32 requestId, bytes payload, uint256 timestamp, address sender, uint256 value, uint256 maxFeeValue)
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
	outstruct.Value = abi.ConvertType(out[7], new(big.Int)).(*big.Int)
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
// the contract method with ID 0x7f8b54a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) PackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, refund, applicationFees, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7f8b54a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, refund, applicationFees, signature)
}

// PackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34f23a9d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 value, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, value *big.Int, maxFeeValue *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, value, maxFeeValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34f23a9d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 value, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, value *big.Int, maxFeeValue *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, value, maxFeeValue)
}

// UnpackSubmitRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x34f23a9d.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 value, uint256 maxFeeValue) payable returns(bytes32)
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
// Solidity: event UserEvent(uint64 indexed applicationId, bytes32 indexed requestId, bytes encryptedData)
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
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea26469706673582212205d9c26453cb6e65a78b01b759d1937994a0e1d443cf98d7fb7a2e3c5b0d602eb64736f6c634300081e0033",
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
