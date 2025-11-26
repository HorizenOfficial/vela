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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AuthorityAlreadyPresent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotPresent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AddedAuthority\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"RemovedAuthority\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"addAllowedAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"checkAuthorityIsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"removeAllowedAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "38f83b752dabbc1250a3a7a867ed9a9f7b",
	Bin: "0x60806040523461002f576100196100146100f4565b610115565b610021610034565b6107d96102cf82396107d990f35b61003a565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100669061003e565b810190811060018060401b0382111761007e57604052565b610048565b9061009661008f610034565b928361005c565b565b5f80fd5b60018060a01b031690565b6100b09061009c565b90565b6100bc816100a7565b036100c357565b5f80fd5b905051906100d4826100b3565b565b906020828203126100ef576100ec915f016100c7565b90565b610098565b610112610aa88038038061010781610083565b9283398101906100d6565b90565b61011e90610170565b565b90565b90565b61013a61013561013f92610120565b610123565b61009c565b90565b61014b90610126565b90565b610157906100a7565b9052565b919061016e905f6020850194019061014e565b565b8061018b6101856101805f610142565b6100a7565b916100a7565b1461019b576101999061026f565b565b6101be6101a75f610142565b5f918291631e4fbdf760e01b83526004830161015b565b0390fd5b5f1c90565b60018060a01b031690565b6101de6101e3916101c2565b6101c7565b90565b6101f090546101d2565b90565b5f1b90565b9061020960018060a01b03916101f3565b9181191691161790565b61022761022261022c9261009c565b610123565b61009c565b90565b61023890610213565b90565b6102449061022f565b90565b90565b9061025f61025a6102669261023b565b610247565b82546101f8565b9055565b5f0190565b6102785f6101e6565b610282825f61024a565b906102b66102b07f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09361023b565b9161023b565b916102bf610034565b806102c98161026a565b0390a356fe60806040526004361015610013575b6102db565b61001d5f3561007c565b806318e807fe1461007757806358ae367f14610072578063715018a61461006d5780638da5cb5b14610068578063f2fde38b146100635763f3f003340361000e576102a7565b610274565b610221565b6101cc565b610187565b61012c565b60e01c90565b60405190565b5f80fd5b5f80fd5b67ffffffffffffffff1690565b6100a681610090565b036100ad57565b5f80fd5b905035906100be8261009d565b565b60018060a01b031690565b6100d4906100c0565b90565b6100e0816100cb565b036100e757565b5f80fd5b905035906100f8826100d7565b565b9190604083820312610122578061011661011f925f86016100b1565b936020016100eb565b90565b61008c565b5f0190565b3461015b5761014561013f3660046100fa565b90610492565b61014d610082565b8061015781610127565b0390f35b610088565b151590565b61016e90610160565b9052565b9190610185905f60208501940190610165565b565b346101b8576101b46101a361019d3660046100fa565b906104a2565b6101ab610082565b91829182610172565b0390f35b610088565b5f9103126101c757565b61008c565b346101fa576101dc3660046101bd565b6101e461051a565b6101ec610082565b806101f681610127565b0390f35b610088565b610208906100cb565b9052565b919061021f905f602085019401906101ff565b565b34610251576102313660046101bd565b61024d61023c610554565b610244610082565b9182918261020c565b0390f35b610088565b9060208282031261026f5761026c915f016100eb565b90565b61008c565b346102a25761028c610287366004610256565b6105ce565b610294610082565b8061029e81610127565b0390f35b610088565b346102d6576102c06102ba3660046100fa565b9061069f565b6102c8610082565b806102d281610127565b0390f35b610088565b5f80fd5b906102f1916102ec6106ab565b6103e8565b565b90565b61030a61030561030f92610090565b6102f3565b610090565b90565b9061031c906102f6565b5f5260205260405f2090565b61033c610337610341926100c0565b6102f3565b6100c0565b90565b61034d90610328565b90565b61035990610344565b90565b9061036690610350565b5f5260205260405f2090565b5f1c90565b60ff1690565b61038961038e91610372565b610377565b90565b61039b905461037d565b90565b5f1b90565b906103af60ff9161039e565b9181191691161790565b6103c290610160565b90565b90565b906103dd6103d86103e4926103b9565b6103c5565b82546103a3565b9055565b6104066104016103fa60018490610312565b849061035c565b610391565b6104765761042a600161042561041e60018590610312565b859061035c565b6103c8565b9061045e6104587f3bff442391e2aff90b2d5810c9574fae9339f6133322e135e3b834d3190122c5936102f6565b91610350565b91610467610082565b8061047181610127565b0390a3565b5f63978543a760e01b81528061048e60048201610127565b0390fd5b9061049c916102df565b565b5f90565b6104c7916104bd6104c2926104b561049e565b506001610312565b61035c565b610391565b90565b6104d26106ab565b6104da610507565b565b90565b6104f36104ee6104f8926104dc565b6102f3565b6100c0565b90565b610504906104df565b90565b6105186105135f6104fb565b610737565b565b6105226104ca565b565b5f90565b60018060a01b031690565b61053f61054491610372565b610528565b90565b6105519054610533565b90565b61055c610524565b506105665f610547565b90565b61057a906105756106ab565b61057c565b565b8061059761059161058c5f6104fb565b6100cb565b916100cb565b146105a7576105a590610737565b565b6105ca6105b35f6104fb565b5f918291631e4fbdf760e01b83526004830161020c565b0390fd5b6105d790610569565b565b906105eb916105e66106ab565b6105ed565b565b61061461060e61060961060260018590610312565b859061035c565b610391565b15610160565b610683576106375f61063261062b60018590610312565b859061035c565b6103c8565b9061066b6106657f62aaf109b50ea5e1c587178c1697a5cec868bb882c63a62388de7496e7615aa0936102f6565b91610350565b91610674610082565b8061067e81610127565b0390a3565b5f633d92db7960e21b81528061069b60048201610127565b0390fd5b906106a9916105d9565b565b6106b3610554565b6106cc6106c66106c1610796565b6100cb565b916100cb565b036106d357565b6106f56106de610796565b5f91829163118cdaa760e01b83526004830161020c565b0390fd5b9061070a60018060a01b039161039e565b9181191691161790565b90565b9061072c61072761073392610350565b610714565b82546106f9565b9055565b6107405f610547565b61074a825f610717565b9061077e6107787f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093610350565b91610350565b91610787610082565b8061079181610127565b0390a3565b61079e610524565b50339056fea26469706673582212201cc010caf4fab287f345f932fc95e8afec0894aa7eac2999ee4a849a24bbd04264736f6c634300081e0033",
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
// Solidity: constructor(address owner) returns()
func (authorityRegistry *AuthorityRegistry) PackConstructor(owner common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAddAllowedAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18e807fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addAllowedAuthority(uint64 applicationId, address authority) returns()
func (authorityRegistry *AuthorityRegistry) PackAddAllowedAuthority(applicationId uint64, authority common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("addAllowedAuthority", applicationId, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddAllowedAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18e807fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addAllowedAuthority(uint64 applicationId, address authority) returns()
func (authorityRegistry *AuthorityRegistry) TryPackAddAllowedAuthority(applicationId uint64, authority common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("addAllowedAuthority", applicationId, authority)
}

// PackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58ae367f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkAuthorityIsAllowed(uint64 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) PackCheckAuthorityIsAllowed(applicationId uint64, authority common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58ae367f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkAuthorityIsAllowed(uint64 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) TryPackCheckAuthorityIsAllowed(applicationId uint64, authority common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
}

// UnpackCheckAuthorityIsAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x58ae367f.
//
// Solidity: function checkAuthorityIsAllowed(uint64 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) UnpackCheckAuthorityIsAllowed(data []byte) (bool, error) {
	out, err := authorityRegistry.abi.Unpack("checkAuthorityIsAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
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

// PackRemoveAllowedAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3f00334.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function removeAllowedAuthority(uint64 applicationId, address authority) returns()
func (authorityRegistry *AuthorityRegistry) PackRemoveAllowedAuthority(applicationId uint64, authority common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("removeAllowedAuthority", applicationId, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRemoveAllowedAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3f00334.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function removeAllowedAuthority(uint64 applicationId, address authority) returns()
func (authorityRegistry *AuthorityRegistry) TryPackRemoveAllowedAuthority(applicationId uint64, authority common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("removeAllowedAuthority", applicationId, authority)
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

// AuthorityRegistryAddedAuthority represents a AddedAuthority event raised by the AuthorityRegistry contract.
type AuthorityRegistryAddedAuthority struct {
	ApplicationId uint64
	Authority     common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryAddedAuthorityEventName = "AddedAuthority"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryAddedAuthority) ContractEventName() string {
	return AuthorityRegistryAddedAuthorityEventName
}

// UnpackAddedAuthorityEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AddedAuthority(uint64 indexed applicationId, address indexed authority)
func (authorityRegistry *AuthorityRegistry) UnpackAddedAuthorityEvent(log *types.Log) (*AuthorityRegistryAddedAuthority, error) {
	event := "AddedAuthority"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryAddedAuthority)
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

// AuthorityRegistryRemovedAuthority represents a RemovedAuthority event raised by the AuthorityRegistry contract.
type AuthorityRegistryRemovedAuthority struct {
	ApplicationId uint64
	Authority     common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryRemovedAuthorityEventName = "RemovedAuthority"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryRemovedAuthority) ContractEventName() string {
	return AuthorityRegistryRemovedAuthorityEventName
}

// UnpackRemovedAuthorityEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RemovedAuthority(uint64 indexed applicationId, address indexed authority)
func (authorityRegistry *AuthorityRegistry) UnpackRemovedAuthorityEvent(log *types.Log) (*AuthorityRegistryRemovedAuthority, error) {
	event := "RemovedAuthority"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryRemovedAuthority)
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
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["AuthorityAlreadyPresent"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackAuthorityAlreadyPresentError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["AuthorityNotPresent"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackAuthorityNotPresentError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// AuthorityRegistryAuthorityAlreadyPresent represents a AuthorityAlreadyPresent error raised by the AuthorityRegistry contract.
type AuthorityRegistryAuthorityAlreadyPresent struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuthorityAlreadyPresent()
func AuthorityRegistryAuthorityAlreadyPresentErrorID() common.Hash {
	return common.HexToHash("0x978543a707da62351fc95ab225113459eaee941fa130cc4dd04a4ba858a07e78")
}

// UnpackAuthorityAlreadyPresentError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuthorityAlreadyPresent()
func (authorityRegistry *AuthorityRegistry) UnpackAuthorityAlreadyPresentError(raw []byte) (*AuthorityRegistryAuthorityAlreadyPresent, error) {
	out := new(AuthorityRegistryAuthorityAlreadyPresent)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "AuthorityAlreadyPresent", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryAuthorityNotPresent represents a AuthorityNotPresent error raised by the AuthorityRegistry contract.
type AuthorityRegistryAuthorityNotPresent struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuthorityNotPresent()
func AuthorityRegistryAuthorityNotPresentErrorID() common.Hash {
	return common.HexToHash("0xf64b6de40fd585218ac2167c37e917bb8f7071e4834b8936ea87c85018b24c35")
}

// UnpackAuthorityNotPresentError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuthorityNotPresent()
func (authorityRegistry *AuthorityRegistry) UnpackAuthorityNotPresentError(raw []byte) (*AuthorityRegistryAuthorityNotPresent, error) {
	out := new(AuthorityRegistryAuthorityNotPresent)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "AuthorityNotPresent", raw); err != nil {
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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0xaac80ccc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature, uint256 refundAmount, uint256 applicationFee) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte, refundAmount *big.Int, applicationFee *big.Int) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature, refundAmount, applicationFee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaac80ccc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature, uint256 refundAmount, uint256 applicationFee) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte, refundAmount *big.Int, applicationFee *big.Int) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature, refundAmount, applicationFee)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaac80ccc.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature, uint256 refundAmount, uint256 applicationFee) view returns(bool)
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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"APPLICATION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"}],\"name\":\"markRequestCompleted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"markRequestFailed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a35c5449187bc72768eac252f2f30417e4",
	Bin: "0x6080604052346100335761001d61001461017e565b929190916103c8565b610025610038565b6143b26106a982396143b290f35b61003e565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b9061006a90610042565b810190811060018060401b0382111761008257604052565b61004c565b9061009a610093610038565b9283610060565b565b5f80fd5b60018060a01b031690565b6100b4906100a0565b90565b6100c0906100ab565b90565b6100cc816100b7565b036100d357565b5f80fd5b905051906100e4826100c3565b565b6100ef906100ab565b90565b6100fb816100e6565b0361010257565b5f80fd5b90505190610113826100f2565b565b61011e816100ab565b0361012557565b5f80fd5b9050519061013682610115565b565b6080818303126101795761014e825f83016100d7565b9261017661015f8460208501610106565b9361016d8160408601610129565b93606001610129565b90565b61009c565b61019c614a5b8038038061019181610087565b928339810190610138565b90919293565b5f1b90565b906101b35f19916101a2565b9181191691161790565b90565b90565b90565b6101da6101d56101df926101bd565b6101c3565b6101c0565b90565b90565b906101fa6101f5610201926101c6565b6101e2565b82546101a7565b9055565b90565b61021c61021761022192610205565b6101c3565b6101c0565b90565b9061023961023461024092610208565b6101e2565b82546101a7565b9055565b90565b61025b61025661026092610244565b6101c3565b6100a0565b90565b61026c90610247565b90565b61028361027e610288926100a0565b6101c3565b6100a0565b90565b6102949061026f565b90565b6102a09061028b565b90565b6102ac9061026f565b90565b6102b8906102a3565b90565b5f0190565b906102d160018060a01b03916101a2565b9181191691161790565b6102e49061028b565b90565b90565b906102ff6102fa610306926102db565b6102e7565b82546102c0565b9055565b610313906102a3565b90565b90565b9061032e6103296103359261030a565b610316565b82546102c0565b9055565b6103429061026f565b90565b61034e90610339565b90565b61035a90610339565b90565b90565b9061037561037061037c92610351565b61035d565b82546102c0565b9055565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b92916103d6600a60066101e5565b6103e260056009610224565b836104056103ff6103fa6103f55f610263565b610297565b6100b7565b916100b7565b1480156104cc575b80156104aa575b8015610488575b61046c5761043b61045b926104346104699660076102ea565b6008610319565b61044e61044782610345565b600a610360565b610456610380565b610597565b506104646103a4565b610597565b50565b5f632582a64160e11b815280610484600482016102bb565b0390fd5b50826104a461049e6104995f610263565b6100ab565b916100ab565b1461041b565b50806104c66104c06104bb5f610263565b6100ab565b916100ab565b14610414565b50816104f06104ea6104e56104e05f610263565b6102af565b6100e6565b916100e6565b1461040d565b5f90565b151590565b90565b61050b906104ff565b90565b9061051890610502565b5f5260205260405f2090565b61052d9061026f565b90565b61053990610524565b90565b9061054690610530565b5f5260205260405f2090565b9061055e60ff916101a2565b9181191691161790565b610571906104fa565b90565b90565b9061058c61058761059392610568565b610574565b8254610552565b9055565b61059f6104f6565b506105b46105ae82849061066e565b156104fa565b5f1461063c576105db60016105d65f6105ce81869061050e565b01859061053c565b610577565b906105e461069b565b9061062161061b6106157f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610502565b92610530565b92610530565b9261062a610038565b80610634816102bb565b0390a4600190565b50505f90565b5f1c90565b60ff1690565b61065961065e91610642565b610647565b90565b61066b905461064d565b90565b610694915f61068961068f936106826104f6565b508261050e565b0161053c565b610661565b90565b5f90565b6106a3610697565b50339056fe60806040526004361015610013575b611a03565b61001d5f356101ec565b806301ffc9a7146101e75780631664deeb146101e2578063248a9ca3146101dd5780632a0acc6a146101d85780632f2ff15d146101d357806334f23a9d146101ce57806336568abe146101c95780633b473cb9146101c45780635423112f146101bf5780635d1be609146101ba5780635e339384146101b55780636e91d133146101b057806380a1f712146101ab57806391d14854146101a65780639588eca2146101a15780639968da961461019c5780639a17995a146101975780639c73eeaf146101925780639c8d416f1461018d5780639fabc36f14610188578063a217fddf14610183578063a9ba045e1461017e578063aa3aa46014610179578063c404c35914610174578063c415b95c1461016f578063d2c35ce81461016a578063d3c7fc7014610165578063d547741f146101605763f7d9cc1e0361000e576119ce565b61198b565b61194e565b61163b565b6115c5565b61151f565b611408565b61138b565b6112e8565b611278565b611243565b6111db565b611181565b6110cc565b61104b565b610fe0565b610fab565b610da5565b610cd3565b610c8a565b610c11565b6107e1565b610658565b610626565b61048b565b6103ea565b610386565b610310565b610278565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f80fd5b63ffffffff60e01b1690565b61021981610204565b0361022057565b5f80fd5b9050359061023182610210565b565b9060208282031261024c57610249915f01610224565b90565b6101fc565b151590565b61025f90610251565b9052565b9190610276905f60208501940190610256565b565b346102a8576102a461029361028e366004610233565b611a0b565b61029b6101f2565b91829182610263565b0390f35b6101f8565b5f9103126102b757565b6101fc565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b6102e86102bc565b90565b90565b6102f7906102eb565b9052565b919061030e905f602085019401906102ee565b565b34610340576103203660046102ad565b61033c61032b6102e0565b6103336101f2565b918291826102fb565b0390f35b6101f8565b61034e816102eb565b0361035557565b5f80fd5b9050359061036682610345565b565b906020828203126103815761037e915f01610359565b90565b6101fc565b346103b6576103b26103a161039c366004610368565b611a65565b6103a96101f2565b918291826102fb565b0390f35b6101f8565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103e76103bb565b90565b3461041a576103fa3660046102ad565b6104166104056103df565b61040d6101f2565b918291826102fb565b0390f35b6101f8565b60018060a01b031690565b6104339061041f565b90565b61043f8161042a565b0361044657565b5f80fd5b9050359061045782610436565b565b9190604083820312610481578061047561047e925f8601610359565b9360200161044a565b90565b6101fc565b5f0190565b346104ba576104a461049e366004610459565b90611ab0565b6104ac6101f2565b806104b681610486565b0390f35b6101f8565b60ff1690565b6104ce816104bf565b036104d557565b5f80fd5b905035906104e6826104c5565b565b67ffffffffffffffff1690565b6104fe816104e8565b0361050557565b5f80fd5b90503590610516826104f5565b565b6004111561052257565b5f80fd5b9050359061053382610518565b565b5f80fd5b5f80fd5b5f80fd5b909182601f8301121561057b5781359167ffffffffffffffff831161057657602001926001830284011161057157565b61053d565b610539565b610535565b90565b61058c81610580565b0361059357565b5f80fd5b905035906105a482610583565b565b91909160c081840312610621576105bf835f83016104d9565b926105cd8160208401610509565b926105db8260408501610526565b9260608101359167ffffffffffffffff831161061c5761060084610619948401610541565b9390946106108160808601610597565b9360a001610597565b90565b610200565b6101fc565b6106546106436106373660046105a6565b95949094939193612558565b61064b6101f2565b918291826102fb565b0390f35b346106875761067161066b366004610459565b90612572565b6106796101f2565b8061068381610486565b0390f35b6101f8565b6011111561069657565b5f80fd5b905035906106a78261068c565b565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906106d5906106ad565b810190811067ffffffffffffffff8211176106ef57604052565b6106b7565b906107076107006101f2565b92836106cb565b565b67ffffffffffffffff8111610727576107236020916106ad565b0190565b6106b7565b90825f939282370152565b9092919261074c61074782610709565b6106f4565b93818552602085019082840111610768576107669261072c565b565b6106a9565b9080601f8301121561078b5781602061078893359101610737565b90565b610535565b916060838303126107dc576107a7825f8501610359565b926107b5836020830161069a565b92604082013567ffffffffffffffff81116107d7576107d4920161076d565b90565b610200565b6101fc565b34610810576107fa6107f4366004610790565b91612a23565b6108026101f2565b8061080c81610486565b0390f35b6101f8565b61081e906102eb565b90565b9061082b90610815565b5f5260205260405f2090565b5f1c90565b60ff1690565b61084e61085391610837565b61083c565b90565b6108609054610842565b90565b60081c90565b67ffffffffffffffff1690565b61088261088791610863565b610869565b90565b6108949054610876565b90565b60481c90565b60ff1690565b6108af6108b491610897565b61089d565b90565b6108c190546108a3565b90565b90565b6108d36108d891610837565b6108c4565b90565b6108e590546108c7565b90565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561091c575b602083101461091757565b6108e8565b91607f169161090c565b60209181520190565b5f5260205f2090565b905f929180549061095261094b836108fc565b8094610926565b916001811690815f146109a9575060011461096d575b505050565b61097a919293945061092f565b915f925b81841061099157505001905f8080610968565b6001816020929593955484860152019101929061097e565b92949550505060ff19168252151560200201905f8080610968565b906109ce91610938565b90565b906109f16109ea926109e16101f2565b938480926109c4565b03836106cb565b565b90565b610a02610a0791610837565b6109f3565b90565b610a1490546109f6565b90565b60018060a01b031690565b610a2e610a3391610837565b610a17565b90565b610a409054610a22565b90565b610a4e906002610821565b610a595f8201610856565b91610a655f830161088a565b91610a715f82016108b7565b91610a7e600183016108db565b91610a8b600282016109d1565b91610a9860038301610a0a565b91610aa560048201610a36565b91610abe6006610ab760058501610a0a565b9301610a0a565b90565b610aca906104bf565b9052565b610ad7906104e8565b9052565b634e487b7160e01b5f52602160045260245ffd5b60041115610af957565b610adb565b90610b0882610aef565b565b610b1390610afe565b90565b610b1f90610b0a565b9052565b5190565b60209181520190565b90825f9392825e0152565b610b5a610b63602093610b6893610b5181610b23565b93848093610b27565b95869101610b30565b6106ad565b0190565b610b7590610580565b9052565b610b829061042a565b9052565b94610be9610100979b9a9892610bf492610bdc610c0898610bd2610c0f9e99610bc8610bfe9a60208f610bc161012082019a5f830190610ac1565b0190610ace565b60408d0190610b16565b60608b01906102ee565b88820360808a0152610b3b565b9a60a0870190610b6c565b60c0850190610b79565b60e0830190610b6c565b0190610b6c565b565b34610c4b57610c47610c2c610c27366004610368565b610a43565b95610c3e9997999591959492946101f2565b998a998a610b86565b0390f35b6101f8565b9091606082840312610c8557610c82610c6b845f8501610359565b93610c798160208601610597565b93604001610597565b90565b6101fc565b34610cb957610ca3610c9d366004610c50565b91612d29565b610cab6101f2565b80610cb581610486565b0390f35b6101f8565b9190610cd1905f60208501940190610b6c565b565b34610d0357610ce33660046102ad565b610cff610cee612d36565b610cf66101f2565b91829182610cbe565b0390f35b6101f8565b1c90565b60018060a01b031690565b610d27906008610d2c9302610d08565b610d0c565b90565b90610d3a9154610d17565b90565b610d4960075f90610d2f565b90565b90565b610d63610d5e610d689261041f565b610d4c565b61041f565b90565b610d7490610d4f565b90565b610d8090610d6b565b90565b610d8c90610d77565b9052565b9190610da3905f60208501940190610d83565b565b34610dd557610db53660046102ad565b610dd1610dc0610d3d565b610dc86101f2565b91829182610d90565b0390f35b6101f8565b5190565b60209181520190565b60200190565b610df6906104bf565b9052565b610e03906104e8565b9052565b610e1090610b0a565b9052565b610e1d906102eb565b9052565b610e40610e49602093610e4e93610e3781610b23565b93848093610926565b95869101610b30565b6106ad565b0190565b610e5b90610580565b9052565b610e689061042a565b9052565b90610f169061010080610ed56101208401610e8d5f8801515f870190610ded565b610e9f60208801516020870190610dfa565b610eb160408801516040870190610e07565b610ec360608801516060870190610e14565b60808701518582036080870152610e21565b94610ee860a082015160a0860190610e52565b610efa60c082015160c0860190610e5f565b610f0c60e082015160e0860190610e52565b0151910190610e52565b90565b90610f2391610e6c565b90565b60200190565b90610f40610f3983610dda565b8092610dde565b9081610f5160208302840194610de7565b925f915b838310610f6457505050505090565b90919293946020610f86610f8083856001950387528951610f19565b97610f26565b9301930191939290610f55565b610fa89160208201915f818403910152610f2c565b90565b34610fdb57610fbb3660046102ad565b610fd7610fc6612eff565b610fce6101f2565b91829182610f93565b0390f35b6101f8565b346110115761100d610ffc610ff6366004610459565b90612fea565b6110046101f2565b91829182610263565b0390f35b6101f8565b61102690600861102b9302610d08565b6108c4565b90565b906110399154611016565b90565b61104860015f9061102e565b90565b3461107b5761105b3660046102ad565b61107761106661103c565b61106e6101f2565b918291826102fb565b0390f35b6101f8565b90565b61109761109261109c92611080565b610d4c565b6104e8565b90565b6110a96001611083565b90565b6110b461109f565b90565b91906110ca905f60208501940190610ace565b565b346110fc576110dc3660046102ad565b6110f86110e76110ac565b6110ef6101f2565b918291826110b7565b0390f35b6101f8565b91909160c08184031261117c5761111a835f830161044a565b926111288160208401610509565b926111368260408501610526565b9260608101359167ffffffffffffffff83116111775761115b84611174948401610541565b93909461116b8160808601610597565b9360a001610597565b90565b610200565b6101fc565b346111b8576111b46111a3611197366004611101565b95949094939193613120565b6111ab6101f2565b918291826102fb565b0390f35b6101f8565b906020828203126111d6576111d3915f01610597565b90565b6101fc565b34611209576111f36111ee3660046111bd565b61320a565b6111fb6101f2565b8061120581610486565b0390f35b6101f8565b61121e9060086112239302610d08565b6109f3565b90565b90611231915461120e565b90565b61124060095f90611226565b90565b34611273576112533660046102ad565b61126f61125e611234565b6112666101f2565b91829182610cbe565b0390f35b6101f8565b346112a8576112a461129361128e366004610368565b613215565b61129b6101f2565b91829182610263565b0390f35b6101f8565b90565b5f1b90565b6112c96112c46112ce926112ad565b6112b0565b6102eb565b90565b6112da5f6112b5565b90565b6112e56112d1565b90565b34611318576112f83660046102ad565b6113146113036112dd565b61130b6101f2565b918291826102fb565b0390f35b6101f8565b60018060a01b031690565b61133890600861133d9302610d08565b61131d565b90565b9061134b9154611328565b90565b61135a60085f90611340565b90565b61136690610d6b565b90565b6113729061135d565b9052565b9190611389905f60208501940190611369565b565b346113bb5761139b3660046102ad565b6113b76113a661134e565b6113ae6101f2565b91829182611376565b0390f35b6101f8565b6113d46113cf6113d9926112ad565b610d4c565b6104bf565b90565b6113e55f6113c0565b90565b6113f06113dc565b90565b9190611406905f60208501940190610ac1565b565b34611438576114183660046102ad565b6114346114236113e8565b61142b6101f2565b918291826113f3565b0390f35b6101f8565b906114e790610100806114a6610120840161145e5f8801515f870190610ded565b61147060208801516020870190610dfa565b61148260408801516040870190610e07565b61149460608801516060870190610e14565b60808701518582036080870152610e21565b946114b960a082015160a0860190610e52565b6114cb60c082015160c0860190610e5f565b6114dd60e082015160e0860190610e52565b0151910190610e52565b90565b60409061151661150b61151d9597969460608401908482035f86015261143d565b9660208301906102ee565b0190610256565b565b346115525761152f3660046102ad565b61154e61153a613287565b6115459391936101f2565b938493846114ea565b0390f35b6101f8565b60018060a01b031690565b6115729060086115779302610d08565b611557565b90565b906115859154611562565b90565b611594600a5f9061157a565b90565b6115a09061041f565b90565b6115ac90611597565b9052565b91906115c3905f602085019401906115a3565b565b346115f5576115d53660046102ad565b6115f16115e0611588565b6115e86101f2565b918291826115b0565b0390f35b6101f8565b61160381611597565b0361160a57565b5f80fd5b9050359061161b826115fa565b565b9060208282031261163657611633915f0161160e565b90565b6101fc565b346116695761165361164e36600461161d565b613433565b61165b6101f2565b8061166581610486565b0390f35b6101f8565b67ffffffffffffffff81116116865760208091020190565b6106b7565b67ffffffffffffffff81116116a9576116a56020916106ad565b0190565b6106b7565b909291926116c36116be8261168b565b6106f4565b938185526020850190828401116116df576116dd9261072c565b565b6106a9565b9080601f83011215611702578160206116ff933591016116ae565b90565b610535565b92919061171b6117168261166e565b6106f4565b93818552602080860192028101918383116117725781905b838210611741575050505050565b813567ffffffffffffffff811161176d5760209161176287849387016116e4565b815201910190611733565b610535565b61053d565b9080601f830112156117955781602061179293359101611707565b90565b610535565b67ffffffffffffffff81116117b25760208091020190565b6106b7565b5f80fd5b91906040838203126117f5576117ee906117d560406106f4565b936117e2825f830161160e565b5f860152602001610597565b6020830152565b6117b7565b9092919261180f61180a8261179a565b6106f4565b93818552604060208601920283019281841161184e57915b8383106118345750505050565b602060409161184384866117bb565b815201920191611827565b61053d565b9080601f830112156118715781602061186e933591016117fa565b90565b610535565b9091610120828403126119495761188f835f8401610509565b9261189d8160208501610359565b926118ab8260408301610359565b926118b98360608401610359565b92608083013567ffffffffffffffff811161194457816118da918501611777565b9260a081013567ffffffffffffffff811161193f57826118fb918301611853565b9260c082013567ffffffffffffffff811161193a5761191f846119379285016116e4565b9361192d8160e08601610597565b9361010001610597565b90565b610200565b610200565b610200565b6101fc565b3461198657611970611961366004611876565b97969096959195949294613db4565b6119786101f2565b8061198281610486565b0390f35b6101f8565b346119ba576119a461199e366004610459565b90613df1565b6119ac6101f2565b806119b681610486565b0390f35b6101f8565b6119cb60065f90611226565b90565b346119fe576119de3660046102ad565b6119fa6119e96119bf565b6119f16101f2565b91829182610cbe565b0390f35b6101f8565b5f80fd5b5f90565b611a13611a07565b5080611a2e611a28637965db0b60e01b610204565b91610204565b14908115611a3b575b5090565b611a459150613dfd565b5f611a37565b5f90565b90611a5990610815565b5f5260205260405f2090565b6001611a7d611a8392611a76611a4b565b505f611a4f565b016108db565b90565b90611aa191611a9c611a9782611a65565b613e23565b611aa3565b565b90611aad91613e66565b50565b90611aba91611a86565b565b9695949392919080611add611ad7611ad26113dc565b6104bf565b916104bf565b03611aee57611aeb97611b0a565b90565b5f635428eae760e01b815280611b0660048201610486565b0390fd5b9695949392919081611b2b611b25611b2061109f565b6104e8565b916104e8565b03611b3c57611b3997612202565b90565b5f6309ff9a8360e41b815280611b5460048201610486565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b611b7b611b8191939293610580565b92610580565b8201809211611b8c57565b611b58565b611b9d611ba291610837565b61131d565b90565b611baf9054611b91565b90565b60e01b90565b611bc181610251565b03611bc857565b5f80fd5b90505190611bd982611bb8565b565b90602082820312611bf457611bf1915f01611bcc565b90565b6101fc565b916020611c1a929493611c1360408201965f830190610ace565b0190610b79565b565b611c246101f2565b3d5f823e3d90fd5b5090565b90565b611c47611c42611c4c92611c30565b610d4c565b610580565b90565b611c5a6101206106f4565b90565b90611c67906104bf565b9052565b90611c75906104e8565b9052565b90611c8390610afe565b9052565b90611c91906102eb565b9052565b611ca09136916116ae565b90565b52565b90611cb090610580565b9052565b90611cbe9061042a565b9052565b634e487b7160e01b5f525f60045260245ffd5b611cdf90516104bf565b90565b90611cee60ff916112b0565b9181191691161790565b611d0c611d07611d11926104bf565b610d4c565b6104bf565b90565b90565b90611d2c611d27611d3392611cf8565b611d14565b8254611ce2565b9055565b611d4190516104e8565b90565b60081b90565b90611d5e68ffffffffffffffff0091611d44565b9181191691161790565b611d7c611d77611d81926104e8565b610d4c565b6104e8565b90565b90565b90611d9c611d97611da392611d68565b611d84565b8254611d4a565b9055565b611db19051610afe565b90565b60481b90565b90611dcf69ff00000000000000000091611db4565b9181191691161790565b611de290610afe565b90565b90565b90611dfd611df8611e0492611dd9565b611de5565b8254611dba565b9055565b611e1290516102eb565b90565b90611e215f19916112b0565b9181191691161790565b611e3490610837565b90565b90611e4c611e47611e5392610815565b611e2b565b8254611e15565b9055565b5190565b601f602091010490565b1b90565b91906008611e84910291611e7e5f1984611e65565b92611e65565b9181191691161790565b611ea2611e9d611ea792610580565b610d4c565b610580565b90565b90565b9190611ec3611ebe611ecb93611e8e565b611eaa565b908354611e69565b9055565b5f90565b611ee591611edf611ecf565b91611ead565b565b5b818110611ef3575050565b80611f005f600193611ed3565b01611ee8565b9190601f8111611f16575b505050565b611f22611f479361092f565b906020611f2e84611e5b565b83019310611f4f575b611f4090611e5b565b0190611ee7565b5f8080611f11565b9150611f4081929050611f37565b90611f6d905f1990600802610d08565b191690565b81611f7c91611f5d565b906002021790565b90611f8e81610b23565b9067ffffffffffffffff821161204e57611fb282611fac85546108fc565b85611f06565b602090601f8311600114611fe657918091611fd5935f92611fda575b5050611f72565b90555b565b90915001515f80611fce565b601f19831691611ff58561092f565b925f5b8181106120365750916002939185600196941061201c575b50505002019055611fd8565b61202c910151601f841690611f5d565b90555f8080612010565b91936020600181928787015181550195019201611ff8565b6106b7565b9061205d91611f84565b565b6120699051610580565b90565b9061208161207c61208892611e8e565b611eaa565b8254611e15565b9055565b612096905161042a565b90565b906120aa60018060a01b03916112b0565b9181191691161790565b6120bd90610d6b565b90565b90565b906120d86120d36120df926120b4565b6120c0565b8254612099565b9055565b906121bc61010060066121c2946121075f82016121015f8801611cd5565b90611d17565b61211f5f820161211960208801611d37565b90611d87565b6121375f820161213160408801611da7565b90611de8565b6121506001820161214a60608801611e08565b90611e37565b6121696002820161216360808801611e57565b90612053565b6121826003820161217c60a0880161205f565b9061206c565b61219b6004820161219560c0880161208c565b906120c3565b6121b4600582016121ae60e0880161205f565b9061206c565b01920161205f565b9061206c565b565b906121ce916120e3565b565b906121da90611e8e565b5f5260205260405f2090565b6121ef90610580565b5f1981146121fd5760010190565b611b58565b9693919592949096503461222861222261221d868890611b6c565b610580565b91610580565b0361253c578361224961224361223e6009610a0a565b610580565b91610580565b1061252057612256612d36565b61227161226b6122666006610a0a565b610580565b91610580565b1015612504578461228b6122856003610afe565b91610afe565b145f146124305761229d828290611c2c565b6122b06122aa6085611c33565b91610580565b03612414575b33868684849087926122c86005610a0a565b946122d296613120565b96959491879091904293339596976122e8611c4f565b995f8b01906122f691611c5d565b60208a019061230491611c6b565b604089019061231291611c79565b606088019061232091611c87565b61232991611c95565b608086019061233791611ca3565b60a085019061234591611ca6565b60c084019061235391611cb4565b60e083019061236191611ca6565b61010082019061237091611ca6565b60028261237c91610821565b90612386916121c4565b8060036123936005610a0a565b61239c916121d0565b906123a691611e37565b6123b06005610a0a565b6123b9906121e6565b6123c490600561206c565b80337f73394f43049193ecd2e0c22eefa0ecf10987ce9c72da906a82a3336dc2ac0547916123f190610815565b906123fb906120b4565b916124046101f2565b8061240e81610486565b0390a390565b5f637c6953f960e01b81528061242c60048201610486565b0390fd5b8461244461243e6002610afe565b91610afe565b148061246c575b156122b6575f63622a850b60e11b81528061246860048201610486565b0390fd5b5061247f61247a6008611ba5565b61135d565b60206358ae367f9188906124a533946124b06124996101f2565b96879586948594611bb2565b845260048401611bf9565b03915afa80156124ff576124cc915f916124d1575b5015610251565b61244b565b6124f2915060203d81116124f8575b6124ea81836106cb565b810190611bdb565b5f6124c5565b503d6124e0565b611c1c565b5f63d8219b3d60e01b81528061251c60048201610486565b0390fd5b5f63c77f97eb60e01b81528061253860048201610486565b0390fd5b5f632a9ffab760e21b81528061255460048201610486565b0390fd5b9061256f96959493929161256a611a4b565b611abc565b90565b908061258d612587612582613f15565b61042a565b9161042a565b0361259e5761259b91613f22565b50565b5f63334bd91960e11b8152806125b660048201610486565b0390fd5b906125d592916125d06125cb6102bc565b613e23565b6127ce565b565b6125e66125ec91939293610580565b92610580565b82039182116125f757565b611b58565b61261061260b612615926112ad565b610d4c565b610580565b90565b61262190610d4f565b90565b61262d90612618565b90565b61263990610d6b565b90565b905090565b61264c5f809261263c565b0190565b61265990612641565b90565b9061266e6126698361168b565b6106f4565b918252565b606090565b3d5f14612693576126883d61265c565b903d5f602084013e5b565b61269b612673565b90612691565b9160206126c29294936126bb60408201965f830190610b79565b0190610b6c565b565b6126d06126d591610837565b611557565b90565b6126e290546126c4565b90565b600311156126ef57565b610adb565b906126fe826126e5565b565b612709906126f4565b90565b61271590612700565b9052565b6011111561272357565b610adb565b9061273282612719565b565b61273d90612728565b90565b61274990612734565b9052565b5190565b60209181520190565b612779612782602093612787936127708161274d565b93848093612751565b95869101610b30565b6106ad565b0190565b9493916060916127cc946127b76127c4936127ad60808b01945f8c019061270c565b60208a0190612740565b878203604089015261275a565b940190610b6c565b565b9190916127e36127dd82613215565b15610251565b612a07576127fe60046127f860028490610821565b01610a36565b612852612818600561281260028690610821565b01610a0a565b61284c612832600661282c60028890610821565b01610a0a565b9161283b61420f565b916128466009610a0a565b906125d7565b90611b6c565b90816128666128605f6125fc565b91610580565b1161295e575b50505f8061288261287d600a6126d8565b612630565b61288c6009610a0a565b6128946101f2565b908161289f81612650565b03925af16128ab612678565b505f14612907579060019290926129016128c56009610a0a565b6128ef7fe6b7582c9b39746e7668d43e4a8a19804b5a751841687ff55e5ca5233d854b5295610815565b956128f86101f2565b9485948561278b565b0390a25b565b90600292909261295661291a6009610a0a565b6129447fe6b7582c9b39746e7668d43e4a8a19804b5a751841687ff55e5ca5233d854b5295610815565b9561294d6101f2565b9485948561278b565b0390a2612905565b5f8061297161296c84612624565b612630565b8461297a6101f2565b908161298581612650565b03925af1612991612678565b5061299c575b61286c565b6129b25f6129ac60028690610821565b0161088a565b8391926129e86129e27f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611d68565b93610815565b936129fd6129f46101f2565b928392836126a1565b0390a35f80612997565b5f6302e8145360e61b815280612a1f60048201610486565b0390fd5b90612a2e92916125ba565b565b90612a4b9291612a46612a416102bc565b613e23565b612b47565b565b612a586101206106f4565b90565b90612b39612b2f6006612a6c612a4d565b94612a83612a7b5f8301610856565b5f8801611c5d565b612a9a612a915f830161088a565b60208801611c6b565b612ab1612aa85f83016108b7565b60408801611c79565b612ac9612ac0600183016108db565b60608801611c87565b612ae1612ad8600283016109d1565b60808801611ca3565b612af9612af060038301610a0a565b60a08801611ca6565b612b11612b0860048301610a36565b60c08801611cb4565b612b29612b2060058301610a0a565b60e08801611ca6565b01610a0a565b6101008401611ca6565b565b612b4490612a5b565b90565b919091612b5c612b5682613215565b15610251565b612d0d57612b74612b6f60028390610821565b612b3b565b92612b80818490611b6c565b612b9e612b98612b93610100880161205f565b610580565b91610580565b03612cf15782612bbf612bb9612bb46009610a0a565b610580565b91610580565b10612cd557612c1e9381612bdb612bd55f6125fc565b91610580565b11612c20575b50505f80612bf7612bf2600a6126d8565b612630565b84612c006101f2565b9081612c0b81612650565b03925af150612c18612678565b506142cd565b565b5f80612c3e612c39612c3460c0860161208c565b612624565b612630565b84612c476101f2565b9081612c5281612650565b03925af1612c5e612678565b50612c69575b612be1565b612c7560208201611d37565b612c8260c085930161208c565b92612cb6612cb07f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611d68565b93610815565b93612ccb612cc26101f2565b928392836126a1565b0390a35f80612c64565b5f632a9ffab760e21b815280612ced60048201610486565b0390fd5b5f632a9ffab760e21b815280612d0960048201610486565b0390fd5b5f6302e8145360e61b815280612d2560048201610486565b0390fd5b90612d349291612a30565b565b612d3e611ecf565b50612d496005610a0a565b612d64612d5e612d596004610a0a565b610580565b91610580565b115f14612d8b57612d88612d786005610a0a565b612d826004610a0a565b906125d7565b90565b612d945f6125fc565b90565b606090565b67ffffffffffffffff8111612db45760208091020190565b6106b7565b90612dcb612dc683612d9c565b6106f4565b918252565b5f90565b5f90565b5f90565b5f90565b606090565b5f90565b5f90565b612df5612a4d565b90602080808080808080808a612e09612dd0565b815201612e14612dd4565b815201612e1f612dd8565b815201612e2a612ddc565b815201612e35612de0565b815201612e40612de5565b815201612e4b612de9565b815201612e56612de5565b815201612e61612de5565b81525050565b612e6f612ded565b90565b5f5b828110612e8057505050565b602090612e8b612e67565b8184015201612e74565b90612eba612ea283612db9565b92602080612eb08693612d9c565b9201910390612e72565b565b634e487b7160e01b5f52603260045260245ffd5b90612eda82610dda565b811015612eeb576020809102010190565b612ebc565b6001612efc9101610580565b90565b612f07612d97565b50612f18612f13612d36565b612e95565b90612f236004610a0a565b91612f2c611ecf565b925b80612f4a612f44612f3f6005610a0a565b610580565b91610580565b1015612fa657612f9a612fa091612f93612f78612f71612f6c600385906121d0565b6108db565b6002610821565b85612f838992612b3b565b612f8d8383612ed0565b52612ed0565b5150612ef0565b93612ef0565b92612f2e565b5090915090565b90612fb7906120b4565b5f5260205260405f2090565b60ff1690565b612fd5612fda91610837565b612fc3565b90565b612fe79054612fc9565b90565b613010915f61300561300b93612ffe611a07565b5082611a4f565b01612fad565b612fdd565b90565b60601b90565b61302290613013565b90565b61302e90613019565b90565b61303d6130429161042a565b613025565b9052565b60c01b90565b61305590613046565b90565b613064613069916104e8565b61304c565b9052565b60f81b90565b61307c9061306d565b90565b61308b61309091610b0a565b613073565b9052565b9091826130a4816130ab9361263c565b809361072c565b0190565b90565b6130be6130c391610580565b6130af565b9052565b936131169560016020999895613100600861310e976130f860148f9c6130f0816131079c613031565b018092613058565b01809261307f565b0191613094565b80926130b2565b0180926130b2565b0190565b60200190565b94613151939296919461316096613135611a4b565b50959793909192936131456101f2565b988997602089016130c7565b602082018103825203826106cb565b61317261316c82610b23565b9161311a565b2090565b61318f9061318a6131856103bb565b613e23565b613191565b565b806131a461319e5f6125fc565b91610580565b146131ee576131b481600661206c565b6131e97e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db1916131e06101f2565b91829182610cbe565b0390a1565b5f632a9ffab760e21b81528061320660048201610486565b0390fd5b61321390613176565b565b61321d611a07565b50613226612d36565b6132386132325f6125fc565b91610580565b119081613244575b5090565b905061327561326f613269613264600361325e6004610a0a565b906121d0565b6108db565b926102eb565b916102eb565b145f613240565b613284612ded565b90565b61328f61327c565b50613298611a4b565b506132a1611a07565b506132aa612d36565b6132bc6132b65f6125fc565b91610580565b116132db576132c961327c565b6132d360016108db565b915f91929190565b6133026132fb6132f660036132f06004610a0a565b906121d0565b6108db565b6002610821565b61330c60016108db565b91613318600192612b3b565b929190565b6133369061333161332c6103bb565b613e23565b6133b1565b565b61334c613347613351926112ad565b610d4c565b61041f565b90565b61335d90613338565b90565b61336990612618565b90565b90565b9061338461337f61338b92613360565b61336c565b8254612099565b9055565b61339890612630565b9052565b91906133af905f6020850194019061338f565b565b806133cc6133c66133c15f613354565b61042a565b91612630565b14613417576133dc81600a61336f565b6134127fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f916134096101f2565b9182918261339c565b0390a1565b5f632582a64160e11b81528061342f60048201610486565b0390fd5b61343c9061331d565b565b97969594939291908861346061345a61345561109f565b6104e8565b916104e8565b036134705761346e9861348c565b565b5f6309ff9a8360e41b81528061348860048201610486565b0390fd5b906134ad98979695949392916134a86134a36102bc565b613e23565b613753565b565b6134bb6134c091610837565b610d0c565b90565b6134cd90546134af565b90565b5190565b60209181520190565b60200190565b906134ed91610e21565b90565b60200190565b9061350a613503836134d0565b80926134d4565b908161351b602083028401946134dd565b925f915b83831061352e57505050505090565b9091929394602061355061354a838560019503875289516134e3565b976134f0565b930193019193929061351f565b5190565b60209181520190565b60200190565b61357990611597565b9052565b9060208061359f936135955f8201515f860190613570565b0151910190610e52565b565b906135ae8160409361357d565b0190565b60200190565b906135d56135cf6135c88461355d565b8093613561565b9261356a565b905f5b8181106135e55750505090565b9091926135fe6135f860019286516135a1565b946135b2565b91019190916135d8565b949561366c61367a926101009861365f6136939861369a9d9f9e9c976136889961364b6136559260208f61364461012082019a5f830190610ace565b01906102ee565b60408d01906102ee565b60608b01906102ee565b88820360808a01526134f6565b9086820360a08801526135b8565b9084820360c0860152610b3b565b9660e0830190610b6c565b0190610b6c565b565b906136a68261355d565b8110156136b7576020809102010190565b612ebc565b6136c590610d6b565b90565b906136d2826134d0565b8110156136e3576020809102010190565b612ebc565b6136fd9160208201915f818403910152610b3b565b90565b91602061372192949361371a60408201965f8301906102ee565b01906102ee565b565b61372d9051611597565b90565b91602061375192949361374a60408201965f83019061338f565b0190610b6c565b565b959891969790949893929361376860016108db565b61378261377c6137775f6112b5565b6102eb565b916102eb565b141580613d90575b613d74576137a061379a86613215565b15610251565b613d58576020868b6137ed8a6137f88a8f8f98908c916137c86137c360076134c3565b610d77565b9763aac80ccc96999b9493909192938d956137e16101f2565b9d8e9c8d9b8c9b611bb2565b8b5260048b01613608565b03915afa8015613d5357613814915f91613d25575b5015610251565b613d095761382c61382760028690610821565b612b3b565b97613838818390611b6c565b61385661385061384b6101008d0161205f565b610580565b91610580565b03613ced578161387761387161386c6009610a0a565b610580565b91610580565b10613cd1575f999897995061388a611ecf565b97613893611ecf565b9a5b8a6138b16138ab6138a68d9361355d565b610580565b91610580565b10156138fc578a6138b16138ab6138ee839f6138a69f908f60206138db6138e2926138e89561369c565b510161205f565b90611b6c565b9d612ef0565b9d505050509a99989a613895565b919396999598509193969961391c90613916858790611b6c565b90611b6c565b61393761393161392b306136bc565b31610580565b91610580565b11613cb5576139478885906142cd565b6139505f6125fc565b5b8061396c6139666139618b6134d0565b610580565b91610580565b10156139da576139d5908b8a6139838b84906136c8565b51916139cd6139bb6139b57f09d7e6b35973ce70ef4c7f3535a0d8e74046743324fc7f02809f0f458147e73193611d68565b93610815565b936139c46101f2565b918291826136e8565b0390a3612ef0565b613951565b5091939695509193966139ee826001611e37565b88869192613a25613a1f7f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e93611d68565b93610815565b93613a3a613a316101f2565b92839283613700565b0390a381613a50613a4a5f6125fc565b91610580565b11613c4c575b613af1925f92839289613a6c60c089930161208c565b92613aa0613a9a7f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f93611d68565b93610815565b93613ab5613aac6101f2565b928392836126a1565b0390a3613aca613ac5600a6126d8565b612630565b90613ad36101f2565b9081613ade81612650565b03925af1613aea612678565b5015610251565b613c3057613afe5f6125fc565b5b80613b1a613b14613b0f8661355d565b610580565b91610580565b1015613c2957613b805f80613b43613b3e82613b3789889061369c565b5101613723565b612630565b613b5a6020613b5389889061369c565b510161205f565b613b626101f2565b9081613b6d81612650565b03925af1613b79612678565b5015610251565b613c0d57613c08908583613ba05f613b9988869061369c565b5101613723565b91613bb86020613bb189879061369c565b510161205f565b613beb613be57fccf5242d67663517e020b47096ecb69c7f1a32e85ebd295eae0611a41395b3ae93611d68565b93610815565b93613c00613bf76101f2565b92839283613730565b0390a3612ef0565b613aff565b5f6312171d8360e31b815280613c2560048201610486565b0390fd5b5050509050565b5f6312171d8360e31b815280613c4860048201610486565b0390fd5b613c945f80613c6d613c68613c6360c0870161208c565b612624565b612630565b85613c766101f2565b9081613c8181612650565b03925af1613c8d612678565b5015610251565b15613a56575f6312171d8360e31b815280613cb160048201610486565b0390fd5b5f631e9acf1760e31b815280613ccd60048201610486565b0390fd5b5f632a9ffab760e21b815280613ce960048201610486565b0390fd5b5f632a9ffab760e21b815280613d0560048201610486565b0390fd5b5f638baa579f60e01b815280613d2160048201610486565b0390fd5b613d46915060203d8111613d4c575b613d3e81836106cb565b810190611bdb565b5f61380d565b503d613d34565b611c1c565b5f6302e8145360e61b815280613d7060048201610486565b0390fd5b5f630b6fac0360e41b815280613d8c60048201610486565b0390fd5b5085613dad613da7613da260016108db565b6102eb565b916102eb565b141561378a565b90613dc5989796959493929161343e565b565b90613de291613ddd613dd882611a65565b613e23565b613de4565b565b90613dee91613f22565b50565b90613dfb91613dc7565b565b613e05611a07565b50613e1f613e196301ffc9a760e01b610204565b91610204565b1490565b613e3590613e2f613f15565b90614341565b565b613e4090610251565b90565b90565b90613e5b613e56613e6292613e37565b613e43565b8254611ce2565b9055565b613e6e611a07565b50613e83613e7d828490612fea565b15610251565b5f14613f0b57613eaa6001613ea55f613e9d818690611a4f565b018590612fad565b613e46565b90613eb3613f15565b90613ef0613eea613ee47f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610815565b926120b4565b926120b4565b92613ef96101f2565b80613f0381610486565b0390a4600190565b50505f90565b5f90565b613f1d613f11565b503390565b613f2a611a07565b50613f36818390612fea565b5f14613fbd57613f5c5f613f575f613f4f818690611a4f565b018590612fad565b613e46565b90613f65613f15565b90613fa2613f9c613f967ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b95610815565b926120b4565b926120b4565b92613fab6101f2565b80613fb581610486565b0390a4600190565b50505f90565b9190613fd9613fd4613fe193610815565b611e2b565b908354611e69565b9055565b613ff791613ff1611a4b565b91613fc3565b565b9061400c905f1990602003600802610d08565b8154169055565b905f9161402a6140228261092f565b928354611f72565b905555565b919290602082105f1461408857601f841160011461405857614052929350611f72565b90555b5b565b509061407e61408393600161407561406f8561092f565b92611e5b565b82019101611ee7565b614013565b614055565b506140bf829361409960019461092f565b6140b86140a585611e5b565b820192601f8616806140ca575b50611e5b565b0190611ee7565b600202179055614056565b6140d690888603613ff9565b5f6140b2565b92909168010000000000000000821161413c576020115f1461412d57602081105f146141115761410b91611f72565b90555b5b565b60019160ff19166141218461092f565b5560020201905561410e565b6001915060020201905561410f565b6106b7565b90815461414d816108fc565b90818311614176575b818310614164575b50505050565b61416d9361402f565b5f80808061415e565b614182838383876140dc565b614156565b5f61419191614141565b565b905f036141a5576141a390614187565b565b611cc2565b5f60066141f69282808201556141c38360018301613fe5565b6141d08360028301614193565b6141dd8360038301611ed3565b8260048201556141f08360058301611ed3565b01611ed3565b565b905f0361420a57614208906141aa565b565b611cc2565b6142405f61423b6002614235614230600361422a6004610a0a565b906121d0565b6108db565b90610821565b6141f8565b61425e5f61425960036142536004610a0a565b906121d0565b613fe5565b61427a61427361426e6004610a0a565b6121e6565b600461206c565b565b6142875f8092612751565b0190565b9392906142c36060916142b66142cb956142ac60808a01935f8b019061270c565b6020890190612740565b868103604088015261427c565b940190610b6c565b565b6142d561420f565b5f915f90926143196143077fe6b7582c9b39746e7668d43e4a8a19804b5a751841687ff55e5ca5233d854b5294610815565b946143106101f2565b9384938461428b565b0390a2565b91602061433f92949361433860408201965f830190610b79565b01906102ee565b565b90614356614350838390612fea565b15610251565b61435e575050565b6143785f92839263e2517d3f60e01b84526004840161431e565b0390fdfea264697066735822122087360ba10e3802b2313dd3564ea49d7f104d36a69185792ebe162ffe489a5fa464736f6c634300081e0033",
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
// Solidity: constructor(address _teeAuthenticator, address _authorityRegistry, address updateStatusOperator, address admin) returns()
func (processorEndpoint *ProcessorEndpoint) PackConstructor(_teeAuthenticator common.Address, _authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("", _teeAuthenticator, _authorityRegistry, updateStatusOperator, admin)
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
// the contract method with ID 0xd3c7fc70.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature, uint256 refund, uint256 applicationFees) returns()
func (processorEndpoint *ProcessorEndpoint) PackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte, refund *big.Int, applicationFees *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature, refund, applicationFees)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd3c7fc70.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature, uint256 refund, uint256 applicationFees) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte, refund *big.Int, applicationFees *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature, refund, applicationFees)
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
	Status          uint8
	ErrorCode       uint8
	ErrorMessage    string
	ApplicationFees *big.Int
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
// Solidity: event RequestCompleted(bytes32 indexed requestId, uint8 status, uint8 errorCode, string errorMessage, uint256 applicationFees)
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
