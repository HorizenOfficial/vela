// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tee

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

// StructsWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsWithdrawalRequest struct {
	Receiver common.Address
	Amount   *big.Int
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

// ECDSAMetaData contains all meta data concerning the ECDSA contract.
var ECDSAMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"}]",
	ID:  "843103f81e2875bbb772deded427f45022",
	Bin: "0x6055604b600b8282823980515f1a607314603f577f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f5ffdfea26469706673582212204487116465f2ff257c3290f18df5d99609b03221dcfef803a6c4f6b0627d97e864736f6c634300081e0033",
}

// ECDSA is an auto generated Go binding around an Ethereum contract.
type ECDSA struct {
	abi abi.ABI
}

// NewECDSA creates a new instance of ECDSA.
func NewECDSA() *ECDSA {
	parsed, err := ECDSAMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ECDSA{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ECDSA) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (eCDSA *ECDSA) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], eCDSA.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return eCDSA.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], eCDSA.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return eCDSA.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], eCDSA.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return eCDSA.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ECDSAECDSAInvalidSignature represents a ECDSAInvalidSignature error raised by the ECDSA contract.
type ECDSAECDSAInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignature()
func ECDSAECDSAInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0xf645eedf0193584640b6b90cb9477e4c95b98636c148a891d4c0a146dc46e75a")
}

// UnpackECDSAInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignature()
func (eCDSA *ECDSA) UnpackECDSAInvalidSignatureError(raw []byte) (*ECDSAECDSAInvalidSignature, error) {
	out := new(ECDSAECDSAInvalidSignature)
	if err := eCDSA.abi.UnpackIntoInterface(out, "ECDSAInvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ECDSAECDSAInvalidSignatureLength represents a ECDSAInvalidSignatureLength error raised by the ECDSA contract.
type ECDSAECDSAInvalidSignatureLength struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func ECDSAECDSAInvalidSignatureLengthErrorID() common.Hash {
	return common.HexToHash("0xfce698f7e8e5342cd615f641317bc45fe7e1e4a8b0a14dd1383ff8dc9c41917f")
}

// UnpackECDSAInvalidSignatureLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func (eCDSA *ECDSA) UnpackECDSAInvalidSignatureLengthError(raw []byte) (*ECDSAECDSAInvalidSignatureLength, error) {
	out := new(ECDSAECDSAInvalidSignatureLength)
	if err := eCDSA.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ECDSAECDSAInvalidSignatureS represents a ECDSAInvalidSignatureS error raised by the ECDSA contract.
type ECDSAECDSAInvalidSignatureS struct {
	S [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func ECDSAECDSAInvalidSignatureSErrorID() common.Hash {
	return common.HexToHash("0xd78bce0cccb935155ed6428d1c13e50b7f3550fd2b66b9fe266006fea4a5e1eb")
}

// UnpackECDSAInvalidSignatureSError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func (eCDSA *ECDSA) UnpackECDSAInvalidSignatureSError(raw []byte) (*ECDSAECDSAInvalidSignatureS, error) {
	out := new(ECDSAECDSAInvalidSignatureS)
	if err := eCDSA.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorMetaData contains all meta data concerning the ITeeAuthenticator contract.
var ITeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0xf3c29cb0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3c29cb0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf3c29cb0.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := iTeeAuthenticator.abi.Unpack("checkSignature", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// MathMetaData contains all meta data concerning the Math contract.
var MathMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "362817aca504bb11799e06f16c73252206",
	Bin: "0x6055604b600b8282823980515f1a607314603f577f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f5ffdfea26469706673582212203f472ece9aa8e12ecba837d60d534a2bcf0289eb4e091e2ef06842816b8a5f8164736f6c634300081e0033",
}

// Math is an auto generated Go binding around an Ethereum contract.
type Math struct {
	abi abi.ABI
}

// NewMath creates a new instance of Math.
func NewMath() *Math {
	parsed, err := MathMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Math{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Math) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// MessageHashUtilsMetaData contains all meta data concerning the MessageHashUtils contract.
var MessageHashUtilsMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "2a46deb1e6ec75a21a12c9c84a8059f39d",
	Bin: "0x6055604b600b8282823980515f1a607314603f577f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f5ffdfea26469706673582212202dff53b582791696bbd5915f2e6e2d41f3794e265f4cb76bffd59c30765ce96664736f6c634300081e0033",
}

// MessageHashUtils is an auto generated Go binding around an Ethereum contract.
type MessageHashUtils struct {
	abi abi.ABI
}

// NewMessageHashUtils creates a new instance of MessageHashUtils.
func NewMessageHashUtils() *MessageHashUtils {
	parsed, err := MessageHashUtilsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MessageHashUtils{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MessageHashUtils) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
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

// PanicMetaData contains all meta data concerning the Panic contract.
var PanicMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "d7d9ad298c70ed58e3eba1456b2d564fae",
	Bin: "0x6055604b600b8282823980515f1a607314603f577f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f5ffdfea26469706673582212209fb65c70971bd0cbcc68bb1545f7e4c2f66355de8822340c4c0147cb8a1415e264736f6c634300081e0033",
}

// Panic is an auto generated Go binding around an Ethereum contract.
type Panic struct {
	abi abi.ABI
}

// NewPanic creates a new instance of Panic.
func NewPanic() *Panic {
	parsed, err := PanicMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Panic{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Panic) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// SafeCastMetaData contains all meta data concerning the SafeCast contract.
var SafeCastMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"bits\",\"type\":\"uint8\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"SafeCastOverflowedIntDowncast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"SafeCastOverflowedIntToUint\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"bits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeCastOverflowedUintDowncast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeCastOverflowedUintToInt\",\"type\":\"error\"}]",
	ID:  "ea30ab87e1a3c2e0ff3c21e804c0b5d55f",
	Bin: "0x6055604b600b8282823980515f1a607314603f577f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f5ffdfea2646970667358221220a9cbe372c722d889b85215407936c8df9eb75343ebaaa674f5562cc7c6a7dda964736f6c634300081e0033",
}

// SafeCast is an auto generated Go binding around an Ethereum contract.
type SafeCast struct {
	abi abi.ABI
}

// NewSafeCast creates a new instance of SafeCast.
func NewSafeCast() *SafeCast {
	parsed, err := SafeCastMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &SafeCast{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *SafeCast) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (safeCast *SafeCast) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], safeCast.abi.Errors["SafeCastOverflowedIntDowncast"].ID.Bytes()[:4]) {
		return safeCast.UnpackSafeCastOverflowedIntDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], safeCast.abi.Errors["SafeCastOverflowedIntToUint"].ID.Bytes()[:4]) {
		return safeCast.UnpackSafeCastOverflowedIntToUintError(raw[4:])
	}
	if bytes.Equal(raw[:4], safeCast.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return safeCast.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], safeCast.abi.Errors["SafeCastOverflowedUintToInt"].ID.Bytes()[:4]) {
		return safeCast.UnpackSafeCastOverflowedUintToIntError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// SafeCastSafeCastOverflowedIntDowncast represents a SafeCastOverflowedIntDowncast error raised by the SafeCast contract.
type SafeCastSafeCastOverflowedIntDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedIntDowncast(uint8 bits, int256 value)
func SafeCastSafeCastOverflowedIntDowncastErrorID() common.Hash {
	return common.HexToHash("0x327269a7f29c3c5436f42eeed1c1adf0d4d525f36360483f4e83ab79e98f9089")
}

// UnpackSafeCastOverflowedIntDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedIntDowncast(uint8 bits, int256 value)
func (safeCast *SafeCast) UnpackSafeCastOverflowedIntDowncastError(raw []byte) (*SafeCastSafeCastOverflowedIntDowncast, error) {
	out := new(SafeCastSafeCastOverflowedIntDowncast)
	if err := safeCast.abi.UnpackIntoInterface(out, "SafeCastOverflowedIntDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SafeCastSafeCastOverflowedIntToUint represents a SafeCastOverflowedIntToUint error raised by the SafeCast contract.
type SafeCastSafeCastOverflowedIntToUint struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func SafeCastSafeCastOverflowedIntToUintErrorID() common.Hash {
	return common.HexToHash("0xa8ce4432b175c373e5f41aba830358e5361584f628450fd436c066323ad91ac2")
}

// UnpackSafeCastOverflowedIntToUintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func (safeCast *SafeCast) UnpackSafeCastOverflowedIntToUintError(raw []byte) (*SafeCastSafeCastOverflowedIntToUint, error) {
	out := new(SafeCastSafeCastOverflowedIntToUint)
	if err := safeCast.abi.UnpackIntoInterface(out, "SafeCastOverflowedIntToUint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SafeCastSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the SafeCast contract.
type SafeCastSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func SafeCastSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (safeCast *SafeCast) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*SafeCastSafeCastOverflowedUintDowncast, error) {
	out := new(SafeCastSafeCastOverflowedUintDowncast)
	if err := safeCast.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SafeCastSafeCastOverflowedUintToInt represents a SafeCastOverflowedUintToInt error raised by the SafeCast contract.
type SafeCastSafeCastOverflowedUintToInt struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func SafeCastSafeCastOverflowedUintToIntErrorID() common.Hash {
	return common.HexToHash("0x24775e0629ae69d78c11bae050651b81820407f300ff750ff2be51e4ce75c37f")
}

// UnpackSafeCastOverflowedUintToIntError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func (safeCast *SafeCast) UnpackSafeCastOverflowedUintToIntError(raw []byte) (*SafeCastSafeCastOverflowedUintToInt, error) {
	out := new(SafeCastSafeCastOverflowedUintToInt)
	if err := safeCast.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintToInt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SignedMathMetaData contains all meta data concerning the SignedMath contract.
var SignedMathMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "bcf64d9c1468375457906244955c165c41",
	Bin: "0x6055604b600b8282823980515f1a607314603f577f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f5ffdfea26469706673582212209b05415a057f2feb57073ee9085c5b0aa485ca345b753ddde2d368014803554064736f6c634300081e0033",
}

// SignedMath is an auto generated Go binding around an Ethereum contract.
type SignedMath struct {
	abi abi.ABI
}

// NewSignedMath creates a new instance of SignedMath.
func NewSignedMath() *SignedMath {
	parsed, err := SignedMathMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &SignedMath{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *SignedMath) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// StringsMetaData contains all meta data concerning the Strings contract.
var StringsMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"StringsInsufficientHexLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StringsInvalidAddressFormat\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StringsInvalidChar\",\"type\":\"error\"}]",
	ID:  "9feb222b0160bac7e1e71a58884bbcf822",
	Bin: "0x6055604b600b8282823980515f1a607314603f577f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f5ffdfea2646970667358221220ab44060f970707ae53c4f5d424d30e67a4c8a676cc79c55ae3b0737fe8e5640864736f6c634300081e0033",
}

// Strings is an auto generated Go binding around an Ethereum contract.
type Strings struct {
	abi abi.ABI
}

// NewStrings creates a new instance of Strings.
func NewStrings() *Strings {
	parsed, err := StringsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Strings{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Strings) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (strings *Strings) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], strings.abi.Errors["StringsInsufficientHexLength"].ID.Bytes()[:4]) {
		return strings.UnpackStringsInsufficientHexLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], strings.abi.Errors["StringsInvalidAddressFormat"].ID.Bytes()[:4]) {
		return strings.UnpackStringsInvalidAddressFormatError(raw[4:])
	}
	if bytes.Equal(raw[:4], strings.abi.Errors["StringsInvalidChar"].ID.Bytes()[:4]) {
		return strings.UnpackStringsInvalidCharError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// StringsStringsInsufficientHexLength represents a StringsInsufficientHexLength error raised by the Strings contract.
type StringsStringsInsufficientHexLength struct {
	Value  *big.Int
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringsInsufficientHexLength(uint256 value, uint256 length)
func StringsStringsInsufficientHexLengthErrorID() common.Hash {
	return common.HexToHash("0xe22e27eb9593f9947dc34771120a3349dd201c662753f0b60502fc1d8e422233")
}

// UnpackStringsInsufficientHexLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringsInsufficientHexLength(uint256 value, uint256 length)
func (strings *Strings) UnpackStringsInsufficientHexLengthError(raw []byte) (*StringsStringsInsufficientHexLength, error) {
	out := new(StringsStringsInsufficientHexLength)
	if err := strings.abi.UnpackIntoInterface(out, "StringsInsufficientHexLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StringsStringsInvalidAddressFormat represents a StringsInvalidAddressFormat error raised by the Strings contract.
type StringsStringsInvalidAddressFormat struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringsInvalidAddressFormat()
func StringsStringsInvalidAddressFormatErrorID() common.Hash {
	return common.HexToHash("0x1d15ae44cf5601ace2ebdfa1451a862de1975935c0c8ba6788ae598e5e29a6dd")
}

// UnpackStringsInvalidAddressFormatError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringsInvalidAddressFormat()
func (strings *Strings) UnpackStringsInvalidAddressFormatError(raw []byte) (*StringsStringsInvalidAddressFormat, error) {
	out := new(StringsStringsInvalidAddressFormat)
	if err := strings.abi.UnpackIntoInterface(out, "StringsInvalidAddressFormat", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StringsStringsInvalidChar represents a StringsInvalidChar error raised by the Strings contract.
type StringsStringsInvalidChar struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringsInvalidChar()
func StringsStringsInvalidCharErrorID() common.Hash {
	return common.HexToHash("0x94e2737eaa44cfdde863f17909ef1e0595339c0eb63c5453ecb27449d2dcd443")
}

// UnpackStringsInvalidCharError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringsInvalidChar()
func (strings *Strings) UnpackStringsInvalidCharError(raw []byte) (*StringsStringsInvalidChar, error) {
	out := new(StringsStringsInvalidChar)
	if err := strings.abi.UnpackIntoInterface(out, "StringsInvalidChar", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StructsMetaData contains all meta data concerning the Structs contract.
var StructsMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "31c70fec60e1e1f4cc22f993498ea8b973",
	Bin: "0x6080604052348015600e575f5ffd5b50603e80601a5f395ff3fe60806040525f5ffdfea2646970667358221220923fccdf1e432b5ec327684d48a9c9fe66d7879b2fac76e99c6beedb1abf8fa964736f6c634300081e0033",
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

// TeeAuthenticatorMetaData contains all meta data concerning the TeeAuthenticator contract.
var TeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_teeSigner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeAddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTeeSigner\",\"type\":\"address\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a8d895d4d49ffcfa5f7559dd61ad47cda6",
	Bin: "0x608060405234801561000f575f5ffd5b5060405161166238038061166283398181016040528101906100319190610272565b815f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161009991906102bf565b60405180910390fd5b6100b18161015360201b60201c565b508060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f2b35c72a948e8979932806f94984b9f1d8baaaeb39f63421655174927a0300775f60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff166040516101449291906102d8565b60405180910390a150506102ff565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61024182610218565b9050919050565b61025181610237565b811461025b575f5ffd5b50565b5f8151905061026c81610248565b92915050565b5f5f6040838503121561028857610287610214565b5b5f6102958582860161025e565b92505060206102a68582860161025e565b9150509250929050565b6102b981610237565b82525050565b5f6020820190506102d25f8301846102b0565b92915050565b5f6040820190506102eb5f8301856102b0565b6102f860208301846102b0565b9392505050565b6113568061030c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610060575f3560e01c806343f855c31461006457806347f1769614610082578063715018a61461009e5780638da5cb5b146100a8578063f2fde38b146100c6578063f3c29cb0146100e2575b5f5ffd5b61006c610112565b6040516100799190610951565b60405180910390f35b61009c600480360381019061009791906109a5565b610137565b005b6100a6610241565b005b6100b0610254565b6040516100bd9190610951565b60405180910390f35b6100e060048036038101906100db91906109a5565b61027b565b005b6100fc60048036038101906100f79190610df9565b6102ff565b6040516101099190610f08565b60405180910390f35b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b61013f6104c8565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036101a4576040517f749a2ca000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f2b35c72a948e8979932806f94984b9f1d8baaaeb39f63421655174927a03007760015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16826040516101f6929190610f21565b60405180910390a18060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6102496104c8565b6102525f61054f565b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6102836104c8565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036102f3575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016102ea9190610951565b60405180910390fd5b6102fc8161054f565b50565b5f5f73ffffffffffffffffffffffffffffffffffffffff1660015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610386576040517ff64b1d7b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f856040516020016103989190611063565b6040516020818303038152906040528051906020012090505f856040516020016103c29190611176565b6040516020818303038152906040528051906020012090505f8b8b8b8b86866040516020016103f6969594939291906111d6565b6040516020818303038152906040528051906020012090505f61046461041b83610610565b88888080601f0160208091040260200160405190810160405280939291908181526020018383808284375f81840152601f19601f82011690508083019250505050505050610643565b905060015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161494505050505098975050505050505050565b6104d061066d565b73ffffffffffffffffffffffffffffffffffffffff166104ee610254565b73ffffffffffffffffffffffffffffffffffffffff161461054d5761051161066d565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016105449190610951565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f7f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f5281601c52603c5f209050919050565b5f5f5f5f6106518686610674565b92509250925061066182826106c9565b82935050505092915050565b5f33905090565b5f5f5f60418451036106b4575f5f5f602087015192506040870151915060608701515f1a90506106a68882858561082b565b9550955095505050506106c2565b5f600285515f1b9250925092505b9250925092565b5f60038111156106dc576106db611245565b5b8260038111156106ef576106ee611245565b5b0315610827576001600381111561070957610708611245565b5b82600381111561071c5761071b611245565b5b03610753576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002600381111561076757610766611245565b5b82600381111561077a57610779611245565b5b036107be57805f1c6040517ffce698f70000000000000000000000000000000000000000000000000000000081526004016107b59190611281565b60405180910390fd5b6003808111156107d1576107d0611245565b5b8260038111156107e4576107e3611245565b5b0361082657806040517fd78bce0c00000000000000000000000000000000000000000000000000000000815260040161081d91906112a9565b60405180910390fd5b5b5050565b5f5f5f7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0845f1c1115610867575f600385925092509250610908565b5f6001888888886040515f815260200160405260405161088a94939291906112dd565b6020604051602081039080840390855afa1580156108aa573d5f5f3e3d5ffd5b5050506020604051035190505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036108fb575f60015f5f1b93509350935050610908565b805f5f5f1b935093509350505b9450945094915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61093b82610912565b9050919050565b61094b81610931565b82525050565b5f6020820190506109645f830184610942565b92915050565b5f604051905090565b5f5ffd5b5f5ffd5b61098481610931565b811461098e575f5ffd5b50565b5f8135905061099f8161097b565b92915050565b5f602082840312156109ba576109b9610973565b5b5f6109c784828501610991565b91505092915050565b5f819050919050565b6109e2816109d0565b81146109ec575f5ffd5b50565b5f813590506109fd816109d9565b92915050565b5f819050919050565b610a1581610a03565b8114610a1f575f5ffd5b50565b5f81359050610a3081610a0c565b92915050565b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610a8082610a3a565b810181811067ffffffffffffffff82111715610a9f57610a9e610a4a565b5b80604052505050565b5f610ab161096a565b9050610abd8282610a77565b919050565b5f67ffffffffffffffff821115610adc57610adb610a4a565b5b602082029050602081019050919050565b5f5ffd5b5f5ffd5b5f67ffffffffffffffff821115610b0f57610b0e610a4a565b5b610b1882610a3a565b9050602081019050919050565b828183375f83830152505050565b5f610b45610b4084610af5565b610aa8565b905082815260208101848484011115610b6157610b60610af1565b5b610b6c848285610b25565b509392505050565b5f82601f830112610b8857610b87610a36565b5b8135610b98848260208601610b33565b91505092915050565b5f610bb3610bae84610ac2565b610aa8565b90508083825260208201905060208402830185811115610bd657610bd5610aed565b5b835b81811015610c1d57803567ffffffffffffffff811115610bfb57610bfa610a36565b5b808601610c088982610b74565b85526020850194505050602081019050610bd8565b5050509392505050565b5f82601f830112610c3b57610c3a610a36565b5b8135610c4b848260208601610ba1565b91505092915050565b5f67ffffffffffffffff821115610c6e57610c6d610a4a565b5b602082029050602081019050919050565b5f5ffd5b5f610c8d82610912565b9050919050565b610c9d81610c83565b8114610ca7575f5ffd5b50565b5f81359050610cb881610c94565b92915050565b5f60408284031215610cd357610cd2610c7f565b5b610cdd6040610aa8565b90505f610cec84828501610caa565b5f830152506020610cff848285016109ef565b60208301525092915050565b5f610d1d610d1884610c54565b610aa8565b90508083825260208201905060408402830185811115610d4057610d3f610aed565b5b835b81811015610d695780610d558882610cbe565b845260208401935050604081019050610d42565b5050509392505050565b5f82601f830112610d8757610d86610a36565b5b8135610d97848260208601610d0b565b91505092915050565b5f5ffd5b5f5f83601f840112610db957610db8610a36565b5b8235905067ffffffffffffffff811115610dd657610dd5610da0565b5b602083019150836001820283011115610df257610df1610aed565b5b9250929050565b5f5f5f5f5f5f5f5f60e0898b031215610e1557610e14610973565b5b5f610e228b828c016109ef565b9850506020610e338b828c01610a22565b9750506040610e448b828c01610a22565b9650506060610e558b828c01610a22565b955050608089013567ffffffffffffffff811115610e7657610e75610977565b5b610e828b828c01610c27565b94505060a089013567ffffffffffffffff811115610ea357610ea2610977565b5b610eaf8b828c01610d73565b93505060c089013567ffffffffffffffff811115610ed057610ecf610977565b5b610edc8b828c01610da4565b92509250509295985092959890939650565b5f8115159050919050565b610f0281610eee565b82525050565b5f602082019050610f1b5f830184610ef9565b92915050565b5f604082019050610f345f830185610942565b610f416020830184610942565b9392505050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f610fa382610f71565b610fad8185610f7b565b9350610fbd818560208601610f8b565b610fc681610a3a565b840191505092915050565b5f610fdc8383610f99565b905092915050565b5f602082019050919050565b5f610ffa82610f48565b6110048185610f52565b93508360208202850161101685610f62565b805f5b8581101561105157848403895281516110328582610fd1565b945061103d83610fe4565b925060208a01995050600181019050611019565b50829750879550505050505092915050565b5f6020820190508181035f83015261107b8184610ff0565b905092915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b6110b581610c83565b82525050565b6110c4816109d0565b82525050565b604082015f8201516110de5f8501826110ac565b5060208201516110f160208501826110bb565b50505050565b5f61110283836110ca565b60408301905092915050565b5f602082019050919050565b5f61112482611083565b61112e818561108d565b93506111398361109d565b805f5b8381101561116957815161115088826110f7565b975061115b8361110e565b92505060018101905061113c565b5085935050505092915050565b5f6020820190508181035f83015261118e818461111a565b905092915050565b5f819050919050565b6111b06111ab826109d0565b611196565b82525050565b5f819050919050565b6111d06111cb82610a03565b6111b6565b82525050565b5f6111e1828961119f565b6020820191506111f182886111bf565b60208201915061120182876111bf565b60208201915061121182866111bf565b60208201915061122182856111bf565b60208201915061123182846111bf565b602082019150819050979650505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b61127b816109d0565b82525050565b5f6020820190506112945f830184611272565b92915050565b6112a381610a03565b82525050565b5f6020820190506112bc5f83018461129a565b92915050565b5f60ff82169050919050565b6112d7816112c2565b82525050565b5f6080820190506112f05f83018761129a565b6112fd60208301866112ce565b61130a604083018561129a565b611317606083018461129a565b9594505050505056fea2646970667358221220d69deb4c988b39803200ae52983c902584478c8f742be25122ea97b68491274c64736f6c634300081e0033",
}

// TeeAuthenticator is an auto generated Go binding around an Ethereum contract.
type TeeAuthenticator struct {
	abi abi.ABI
}

// NewTeeAuthenticator creates a new instance of TeeAuthenticator.
func NewTeeAuthenticator() *TeeAuthenticator {
	parsed, err := TeeAuthenticatorMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TeeAuthenticator{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TeeAuthenticator) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address owner, address _teeSigner) returns()
func (teeAuthenticator *TeeAuthenticator) PackConstructor(owner common.Address, _teeSigner common.Address) []byte {
	enc, err := teeAuthenticator.abi.Pack("", owner, _teeSigner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3c29cb0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) PackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3c29cb0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) TryPackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf3c29cb0.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := teeAuthenticator.abi.Unpack("checkSignature", data)
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
func (teeAuthenticator *TeeAuthenticator) PackOwner() []byte {
	enc, err := teeAuthenticator.abi.Pack("owner")
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
func (teeAuthenticator *TeeAuthenticator) TryPackOwner() ([]byte, error) {
	return teeAuthenticator.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (teeAuthenticator *TeeAuthenticator) UnpackOwner(data []byte) (common.Address, error) {
	out, err := teeAuthenticator.abi.Unpack("owner", data)
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
func (teeAuthenticator *TeeAuthenticator) PackRenounceOwnership() []byte {
	enc, err := teeAuthenticator.abi.Pack("renounceOwnership")
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
func (teeAuthenticator *TeeAuthenticator) TryPackRenounceOwnership() ([]byte, error) {
	return teeAuthenticator.abi.Pack("renounceOwnership")
}

// PackTeeSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x43f855c3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function teeSigner() view returns(address)
func (teeAuthenticator *TeeAuthenticator) PackTeeSigner() []byte {
	enc, err := teeAuthenticator.abi.Pack("teeSigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTeeSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x43f855c3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function teeSigner() view returns(address)
func (teeAuthenticator *TeeAuthenticator) TryPackTeeSigner() ([]byte, error) {
	return teeAuthenticator.abi.Pack("teeSigner")
}

// UnpackTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x43f855c3.
//
// Solidity: function teeSigner() view returns(address)
func (teeAuthenticator *TeeAuthenticator) UnpackTeeSigner(data []byte) (common.Address, error) {
	out, err := teeAuthenticator.abi.Unpack("teeSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (teeAuthenticator *TeeAuthenticator) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := teeAuthenticator.abi.Pack("transferOwnership", newOwner)
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
func (teeAuthenticator *TeeAuthenticator) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return teeAuthenticator.abi.Pack("transferOwnership", newOwner)
}

// PackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47f17696.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTee(address newTeeSigner) returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdateTee(newTeeSigner common.Address) []byte {
	enc, err := teeAuthenticator.abi.Pack("updateTee", newTeeSigner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47f17696.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTee(address newTeeSigner) returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdateTee(newTeeSigner common.Address) ([]byte, error) {
	return teeAuthenticator.abi.Pack("updateTee", newTeeSigner)
}

// TeeAuthenticatorOwnershipTransferred represents a OwnershipTransferred event raised by the TeeAuthenticator contract.
type TeeAuthenticatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorOwnershipTransferred) ContractEventName() string {
	return TeeAuthenticatorOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (teeAuthenticator *TeeAuthenticator) UnpackOwnershipTransferredEvent(log *types.Log) (*TeeAuthenticatorOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if log.Topics[0] != teeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TeeAuthenticatorOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := teeAuthenticator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range teeAuthenticator.abi.Events[event].Inputs {
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

// TeeAuthenticatorTeeUpdate represents a TeeUpdate event raised by the TeeAuthenticator contract.
type TeeAuthenticatorTeeUpdate struct {
	OldTee common.Address
	NewTee common.Address
	Raw    *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorTeeUpdateEventName = "TeeUpdate"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorTeeUpdate) ContractEventName() string {
	return TeeAuthenticatorTeeUpdateEventName
}

// UnpackTeeUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TeeUpdate(address oldTee, address newTee)
func (teeAuthenticator *TeeAuthenticator) UnpackTeeUpdateEvent(log *types.Log) (*TeeAuthenticatorTeeUpdate, error) {
	event := "TeeUpdate"
	if log.Topics[0] != teeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TeeAuthenticatorTeeUpdate)
	if len(log.Data) > 0 {
		if err := teeAuthenticator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range teeAuthenticator.abi.Events[event].Inputs {
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
func (teeAuthenticator *TeeAuthenticator) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["TeeAddressCantBeZero"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackTeeAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TeeAuthenticatorECDSAInvalidSignature represents a ECDSAInvalidSignature error raised by the TeeAuthenticator contract.
type TeeAuthenticatorECDSAInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignature()
func TeeAuthenticatorECDSAInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0xf645eedf0193584640b6b90cb9477e4c95b98636c148a891d4c0a146dc46e75a")
}

// UnpackECDSAInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignature()
func (teeAuthenticator *TeeAuthenticator) UnpackECDSAInvalidSignatureError(raw []byte) (*TeeAuthenticatorECDSAInvalidSignature, error) {
	out := new(TeeAuthenticatorECDSAInvalidSignature)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorECDSAInvalidSignatureLength represents a ECDSAInvalidSignatureLength error raised by the TeeAuthenticator contract.
type TeeAuthenticatorECDSAInvalidSignatureLength struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func TeeAuthenticatorECDSAInvalidSignatureLengthErrorID() common.Hash {
	return common.HexToHash("0xfce698f7e8e5342cd615f641317bc45fe7e1e4a8b0a14dd1383ff8dc9c41917f")
}

// UnpackECDSAInvalidSignatureLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func (teeAuthenticator *TeeAuthenticator) UnpackECDSAInvalidSignatureLengthError(raw []byte) (*TeeAuthenticatorECDSAInvalidSignatureLength, error) {
	out := new(TeeAuthenticatorECDSAInvalidSignatureLength)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorECDSAInvalidSignatureS represents a ECDSAInvalidSignatureS error raised by the TeeAuthenticator contract.
type TeeAuthenticatorECDSAInvalidSignatureS struct {
	S [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func TeeAuthenticatorECDSAInvalidSignatureSErrorID() common.Hash {
	return common.HexToHash("0xd78bce0cccb935155ed6428d1c13e50b7f3550fd2b66b9fe266006fea4a5e1eb")
}

// UnpackECDSAInvalidSignatureSError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func (teeAuthenticator *TeeAuthenticator) UnpackECDSAInvalidSignatureSError(raw []byte) (*TeeAuthenticatorECDSAInvalidSignatureS, error) {
	out := new(TeeAuthenticatorECDSAInvalidSignatureS)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the TeeAuthenticator contract.
type TeeAuthenticatorOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func TeeAuthenticatorOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (teeAuthenticator *TeeAuthenticator) UnpackOwnableInvalidOwnerError(raw []byte) (*TeeAuthenticatorOwnableInvalidOwner, error) {
	out := new(TeeAuthenticatorOwnableInvalidOwner)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the TeeAuthenticator contract.
type TeeAuthenticatorOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func TeeAuthenticatorOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (teeAuthenticator *TeeAuthenticator) UnpackOwnableUnauthorizedAccountError(raw []byte) (*TeeAuthenticatorOwnableUnauthorizedAccount, error) {
	out := new(TeeAuthenticatorOwnableUnauthorizedAccount)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorTeeAddressCantBeZero represents a TeeAddressCantBeZero error raised by the TeeAuthenticator contract.
type TeeAuthenticatorTeeAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TeeAddressCantBeZero()
func TeeAuthenticatorTeeAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x749a2ca096befde4321f5c4e6ddb95f574ec95d21a2dead6e38f622a10fd68c7")
}

// UnpackTeeAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TeeAddressCantBeZero()
func (teeAuthenticator *TeeAuthenticator) UnpackTeeAddressCantBeZeroError(raw []byte) (*TeeAuthenticatorTeeAddressCantBeZero, error) {
	out := new(TeeAuthenticatorTeeAddressCantBeZero)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "TeeAddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorTeeIsNotSet represents a TeeIsNotSet error raised by the TeeAuthenticator contract.
type TeeAuthenticatorTeeIsNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TeeIsNotSet()
func TeeAuthenticatorTeeIsNotSetErrorID() common.Hash {
	return common.HexToHash("0xf64b1d7b3df6940d5da676c1cc07a6144e620b0858b161be25b6cd0ef7569425")
}

// UnpackTeeIsNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TeeIsNotSet()
func (teeAuthenticator *TeeAuthenticator) UnpackTeeIsNotSetError(raw []byte) (*TeeAuthenticatorTeeIsNotSet, error) {
	out := new(TeeAuthenticatorTeeIsNotSet)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "TeeIsNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}
