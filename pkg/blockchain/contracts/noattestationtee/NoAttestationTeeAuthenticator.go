// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package noattestationtee

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

// StructsSignatureParams is an auto generated low-level Go binding around an user-defined struct.
type StructsSignatureParams struct {
	ApplicationId      uint64
	PrevStateRoot      [32]byte
	NewStateRoot       [32]byte
	ProcessedRequestId [32]byte
	Events             [][]byte
	EventSubTypes      []string
	WithdrawalRequests []StructsWithdrawalRequest
	RefundAmount       *big.Int
	ApplicationFee     *big.Int
	ErrorCode          uint8
	ErrorMsg           string
}

// StructsWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsWithdrawalRequest struct {
	TokenAddress common.Address
	Receiver     common.Address
	Amount       *big.Int
}

// AbstractTeeAuthenticatorMetaData contains all meta data concerning the AbstractTeeAuthenticator contract.
var AbstractTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "23ef737ffafb52353b16334779ed8f3e1e",
}

// AbstractTeeAuthenticator is an auto generated Go binding around an Ethereum contract.
type AbstractTeeAuthenticator struct {
	abi abi.ABI
}

// NewAbstractTeeAuthenticator creates a new instance of AbstractTeeAuthenticator.
func NewAbstractTeeAuthenticator() *AbstractTeeAuthenticator {
	parsed, err := AbstractTeeAuthenticatorMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &AbstractTeeAuthenticator{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *AbstractTeeAuthenticator) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackPKLENGTH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc91496c6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) PackPKLENGTH() []byte {
	enc, err := abstractTeeAuthenticator.abi.Pack("PK_LENGTH")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPKLENGTH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc91496c6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) TryPackPKLENGTH() ([]byte, error) {
	return abstractTeeAuthenticator.abi.Pack("PK_LENGTH")
}

// UnpackPKLENGTH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc91496c6.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackPKLENGTH(data []byte) (*big.Int, error) {
	out, err := abstractTeeAuthenticator.abi.Unpack("PK_LENGTH", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd119a29d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := abstractTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd119a29d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return abstractTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd119a29d.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := abstractTeeAuthenticator.abi.Unpack("checkSignature", data)
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
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) PackGetPubSecp521r1() []byte {
	enc, err := abstractTeeAuthenticator.abi.Pack("getPubSecp521r1")
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
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) TryPackGetPubSecp521r1() ([]byte, error) {
	return abstractTeeAuthenticator.abi.Pack("getPubSecp521r1")
}

// UnpackGetPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081bec7e.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackGetPubSecp521r1(data []byte) ([]byte, error) {
	out, err := abstractTeeAuthenticator.abi.Unpack("getPubSecp521r1", data)
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
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) PackGetTeeSigner() []byte {
	enc, err := abstractTeeAuthenticator.abi.Pack("getTeeSigner")
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
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) TryPackGetTeeSigner() ([]byte, error) {
	return abstractTeeAuthenticator.abi.Pack("getTeeSigner")
}

// UnpackGetTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dd7ce2f.
//
// Solidity: function getTeeSigner() view returns(address)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackGetTeeSigner(data []byte) (common.Address, error) {
	out, err := abstractTeeAuthenticator.abi.Unpack("getTeeSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], abstractTeeAuthenticator.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return abstractTeeAuthenticator.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], abstractTeeAuthenticator.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return abstractTeeAuthenticator.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], abstractTeeAuthenticator.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return abstractTeeAuthenticator.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	if bytes.Equal(raw[:4], abstractTeeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return abstractTeeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// AbstractTeeAuthenticatorECDSAInvalidSignature represents a ECDSAInvalidSignature error raised by the AbstractTeeAuthenticator contract.
type AbstractTeeAuthenticatorECDSAInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignature()
func AbstractTeeAuthenticatorECDSAInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0xf645eedf0193584640b6b90cb9477e4c95b98636c148a891d4c0a146dc46e75a")
}

// UnpackECDSAInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignature()
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackECDSAInvalidSignatureError(raw []byte) (*AbstractTeeAuthenticatorECDSAInvalidSignature, error) {
	out := new(AbstractTeeAuthenticatorECDSAInvalidSignature)
	if err := abstractTeeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AbstractTeeAuthenticatorECDSAInvalidSignatureLength represents a ECDSAInvalidSignatureLength error raised by the AbstractTeeAuthenticator contract.
type AbstractTeeAuthenticatorECDSAInvalidSignatureLength struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func AbstractTeeAuthenticatorECDSAInvalidSignatureLengthErrorID() common.Hash {
	return common.HexToHash("0xfce698f7e8e5342cd615f641317bc45fe7e1e4a8b0a14dd1383ff8dc9c41917f")
}

// UnpackECDSAInvalidSignatureLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackECDSAInvalidSignatureLengthError(raw []byte) (*AbstractTeeAuthenticatorECDSAInvalidSignatureLength, error) {
	out := new(AbstractTeeAuthenticatorECDSAInvalidSignatureLength)
	if err := abstractTeeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AbstractTeeAuthenticatorECDSAInvalidSignatureS represents a ECDSAInvalidSignatureS error raised by the AbstractTeeAuthenticator contract.
type AbstractTeeAuthenticatorECDSAInvalidSignatureS struct {
	S [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func AbstractTeeAuthenticatorECDSAInvalidSignatureSErrorID() common.Hash {
	return common.HexToHash("0xd78bce0cccb935155ed6428d1c13e50b7f3550fd2b66b9fe266006fea4a5e1eb")
}

// UnpackECDSAInvalidSignatureSError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackECDSAInvalidSignatureSError(raw []byte) (*AbstractTeeAuthenticatorECDSAInvalidSignatureS, error) {
	out := new(AbstractTeeAuthenticatorECDSAInvalidSignatureS)
	if err := abstractTeeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AbstractTeeAuthenticatorTeeIsNotSet represents a TeeIsNotSet error raised by the AbstractTeeAuthenticator contract.
type AbstractTeeAuthenticatorTeeIsNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TeeIsNotSet()
func AbstractTeeAuthenticatorTeeIsNotSetErrorID() common.Hash {
	return common.HexToHash("0xf64b1d7b3df6940d5da676c1cc07a6144e620b0858b161be25b6cd0ef7569425")
}

// UnpackTeeIsNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TeeIsNotSet()
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) UnpackTeeIsNotSetError(raw []byte) (*AbstractTeeAuthenticatorTeeIsNotSet, error) {
	out := new(AbstractTeeAuthenticatorTeeIsNotSet)
	if err := abstractTeeAuthenticator.abi.UnpackIntoInterface(out, "TeeIsNotSet", raw); err != nil {
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

// ECDSAMetaData contains all meta data concerning the ECDSA contract.
var ECDSAMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"}]",
	ID:  "843103f81e2875bbb772deded427f45022",
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea264697066735822122039174b56cf02d3912da8bc0e27c1fe76af7c7b1e8868b40fb866c073c046f09564736f6c634300081e0033",
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
	ABI: "[{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0xd119a29d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd119a29d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd119a29d.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
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

// MathMetaData contains all meta data concerning the Math contract.
var MathMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "362817aca504bb11799e06f16c73252206",
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea2646970667358221220507356191b6a0e0e54ad05ee45eef539005f8f04854db8ee28bd5ece0f25fa2a64736f6c634300081e0033",
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
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea2646970667358221220ed163200e1224a5402d40c50db2b391d712418f14e5e30be1fb1cd75197c11d064736f6c634300081e0033",
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

// NoAttestationTeeAuthenticatorMetaData contains all meta data concerning the NoAttestationTeeAuthenticator contract.
var NoAttestationTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_teeSigner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pubSecp521r1\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeAddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTeeSigner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "6f12a2ebb23d7f97c87b56605496d847e6",
	Bin: "0x6080604052346100305761001a6100146101b7565b9161066d565b610022610035565b611d7a6107b18239611d7a90f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b5f80fd5b60018060a01b031690565b6100b5906100a1565b90565b6100c1816100ac565b036100c857565b5f80fd5b905051906100d9826100b8565b565b5f80fd5b5f80fd5b60018060401b0381116100ff576100fb60209161003f565b0190565b610049565b90825f9392825e0152565b9092919261012461011f826100e3565b610084565b938185526020850190828401116101405761013e92610104565b565b6100df565b9080601f83011215610163578160206101609351910161010f565b90565b6100db565b916060838303126101b25761017f825f85016100cc565b9261018d83602083016100cc565b92604082015160018060401b0381116101ad576101aa9201610145565b90565b61009d565b610099565b6101d561252b803803806101ca81610084565b928339810190610168565b909192565b5f1b90565b906101f060018060a01b03916101da565b9181191691161790565b90565b61021161020c610216926100a1565b6101fa565b6100a1565b90565b610222906101fd565b90565b61022e90610219565b90565b90565b9061024961024461025092610225565b610231565b82546101df565b9055565b5190565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561028c575b602083101461028757565b610258565b91607f169161027c565b5f5260205f2090565b601f602091010490565b1b90565b919060086102c89102916102c25f19846102a9565b926102a9565b9181191691161790565b90565b6102e96102e46102ee926102d2565b6101fa565b6102d2565b90565b90565b919061030a610305610312936102d5565b6102f1565b9083546102ad565b9055565b5f90565b61032c91610326610316565b916102f4565b565b5b81811061033a575050565b806103475f60019361031a565b0161032f565b9190601f811161035d575b505050565b61036961038e93610296565b9060206103758461029f565b83019310610396575b6103879061029f565b019061032e565b5f8080610358565b91506103878192905061037e565b1c90565b906103b8905f19906008026103a4565b191690565b816103c7916103a8565b906002021790565b906103d981610254565b9060018060401b038211610497576103fb826103f5855461026c565b8561034d565b602090601f831160011461042f5791809161041e935f92610423575b50506103bd565b90555b565b90915001515f80610417565b601f1983169161043e85610296565b925f5b81811061047f57509160029391856001969410610465575b50505002019055610421565b610475910151601f8416906103a8565b90555f8080610459565b91936020600181928787015181550195019201610441565b610049565b906104a6916103cf565b565b90565b6104bf6104ba6104c4926104a8565b6101fa565b6100a1565b90565b6104d0906104ab565b90565b5f1c90565b60018060a01b031690565b6104ef6104f4916104d3565b6104d8565b90565b61050190546104e3565b90565b60018060401b0381116105205761051c60209161003f565b0190565b610049565b9061053761053283610504565b610084565b918252565b6105455f610525565b90565b61055061053c565b90565b61055c906100ac565b9052565b60209181520190565b6105886105916020936105969361057f81610254565b93848093610560565b95869101610104565b61003f565b0190565b905f92918054906105b46105ad8361026c565b8094610560565b916001811690815f1461060b57506001146105cf575b505050565b6105dc9192939450610296565b915f925b8184106105f357505001905f80806105ca565b600181602092959395548486015201910192906105e0565b92949550505060ff19168252151560200201905f80806105ca565b929061065c9161064f61066a969461064560808801945f890190610553565b6020870190610553565b8482036040860152610569565b91606081840391015261059a565b90565b61068c929161067e610685926106fa565b6001610234565b600261049c565b6106955f6104c7565b61069f60016104f7565b6106a7610548565b916106e060027f1882fd6e997ef1dc5e2691efe78564a5c540b36547363e8a7dcd454c61a944b6946106d7610035565b94859485610626565b0390a1565b91906106f8905f60208501940190610553565b565b8061071561070f61070a5f6104c7565b6100ac565b916100ac565b146107255761072390610751565b565b6107486107315f6104c7565b5f918291631e4fbdf760e01b8352600483016106e5565b0390fd5b5f0190565b61075a5f6104f7565b610764825f610234565b906107986107927f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093610225565b91610225565b916107a1610035565b806107ab8161074c565b0390a356fe60806040526004361015610013575b610c91565b61001d5f356100bc565b8063081bec7e146100b75780630dd7ce2f146100b257806343f855c3146100ad578063715018a6146100a85780638da5cb5b146100a3578063c64af6fb1461009e578063c91496c614610099578063d119a29d14610094578063f2fde38b1461008f5763fe4993ca0361000e57610c5c565b610ae4565b610a90565b61041e565b61038b565b61029e565b61026b565b610231565b6101b8565b61014a565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f9103126100da57565b6100cc565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b61012061012960209361012e93610117816100df565b938480936100e3565b958691016100ec565b6100f7565b0190565b6101479160208201915f818403910152610101565b90565b3461017a5761015a3660046100d0565b610176610165610ca6565b61016d6100c2565b91829182610132565b0390f35b6100c8565b60018060a01b031690565b6101939061017f565b90565b61019f9061018a565b9052565b91906101b6905f60208501940190610196565b565b346101e8576101c83660046100d0565b6101e46101d3610ce6565b6101db6100c2565b918291826101a3565b0390f35b6100c8565b1c90565b60018060a01b031690565b61020c90600861021193026101ed565b6101f1565b90565b9061021f91546101fc565b90565b61022e60015f90610214565b90565b34610261576102413660046100d0565b61025d61024c610222565b6102546100c2565b918291826101a3565b0390f35b6100c8565b5f0190565b346102995761027b3660046100d0565b610283610d4c565b61028b6100c2565b8061029581610266565b0390f35b6100c8565b346102ce576102ae3660046100d0565b6102ca6102b9610d56565b6102c16100c2565b918291826101a3565b0390f35b6100c8565b5f80fd5b6102e08161018a565b036102e757565b5f80fd5b905035906102f8826102d7565b565b5f80fd5b5f80fd5b5f80fd5b909182601f830112156103405781359167ffffffffffffffff831161033b57602001926001830284011161033657565b610302565b6102fe565b6102fa565b9190916040818403126103865761035e835f83016102eb565b92602082013567ffffffffffffffff81116103815761037d9201610306565b9091565b6102d3565b6100cc565b346103ba576103a461039e366004610345565b916111db565b6103ac6100c2565b806103b681610266565b0390f35b6100c8565b90565b90565b90565b6103dc6103d76103e1926103bf565b6103c5565b6103c2565b90565b6103ee60856103c8565b90565b6103f96103e4565b90565b610405906103c2565b9052565b919061041c905f602085019401906103fc565b565b3461044e5761042e3660046100d0565b61044a6104396103f1565b6104416100c2565b91829182610409565b0390f35b6100c8565b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b90610475906100f7565b810190811067ffffffffffffffff82111761048f57604052565b610457565b906104a76104a06100c2565b928361046b565b565b5f80fd5b67ffffffffffffffff1690565b6104c3816104ad565b036104ca57565b5f80fd5b905035906104db826104ba565b565b90565b6104e9816104dd565b036104f057565b5f80fd5b90503590610501826104e0565b565b67ffffffffffffffff811161051b5760208091020190565b610457565b5f80fd5b67ffffffffffffffff81116105425761053e6020916100f7565b0190565b610457565b90825f939282370152565b9092919261056761056282610524565b610494565b938185526020850190828401116105835761058192610547565b565b610520565b9080601f830112156105a6578160206105a393359101610552565b90565b6102fa565b9291906105bf6105ba82610503565b610494565b93818552602080860192028101918383116106165781905b8382106105e5575050505050565b813567ffffffffffffffff8111610611576020916106068784938701610588565b8152019101906105d7565b6102fa565b610302565b9080601f8301121561063957816020610636933591016105ab565b90565b6102fa565b67ffffffffffffffff81116106565760208091020190565b610457565b67ffffffffffffffff8111610679576106756020916100f7565b0190565b610457565b9092919261069361068e8261065b565b610494565b938185526020850190828401116106af576106ad92610547565b565b610520565b9080601f830112156106d2578160206106cf9335910161067e565b90565b6102fa565b9291906106eb6106e68261063e565b610494565b93818552602080860192028101918383116107425781905b838210610711575050505050565b813567ffffffffffffffff811161073d5760209161073287849387016106b4565b815201910190610703565b6102fa565b610302565b9080601f8301121561076557816020610762933591016106d7565b90565b6102fa565b67ffffffffffffffff81116107825760208091020190565b610457565b6107909061017f565b90565b61079c81610787565b036107a357565b5f80fd5b905035906107b482610793565b565b6107bf816103c2565b036107c657565b5f80fd5b905035906107d7826107b6565b565b91906060838203126108255761081e906107f36060610494565b93610800825f83016102eb565b5f86015261081182602083016107a7565b60208601526040016107ca565b6040830152565b610453565b9092919261083f61083a8261076a565b610494565b93818552606060208601920283019281841161087e57915b8383106108645750505050565b602060609161087384866107d9565b815201920191610857565b610302565b9080601f830112156108a15781602061089e9335910161082a565b90565b6102fa565b600d11156108b057565b5f80fd5b905035906108c1826108a6565b565b91909161016081840312610a07576108dc610160610494565b926108e9815f84016104ce565b5f8501526108fa81602084016104f4565b602085015261090c81604084016104f4565b604085015261091e81606084016104f4565b6060850152608082013567ffffffffffffffff8111610a02578161094391840161061b565b608085015260a082013567ffffffffffffffff81116109fd5781610968918401610747565b60a085015260c082013567ffffffffffffffff81116109f8578161098d918401610883565b60c085015261099f8160e084016107ca565b60e08501526109b28161010084016107ca565b6101008501526109c68161012084016108b4565b61012085015261014082013567ffffffffffffffff81116109f3576109eb92016106b4565b610140830152565b6104a9565b6104a9565b6104a9565b6104a9565b610453565b919091604081840312610a64575f81013567ffffffffffffffff8111610a5f5783610a389183016108c3565b92602082013567ffffffffffffffff8111610a5a57610a579201610588565b90565b6102d3565b6102d3565b6100cc565b151590565b610a7790610a69565b9052565b9190610a8e905f60208501940190610a6e565b565b34610ac157610abd610aac610aa6366004610a0c565b90611612565b610ab46100c2565b91829182610a7b565b0390f35b6100c8565b90602082820312610adf57610adc915f016102eb565b90565b6100cc565b34610b1257610afc610af7366004610ac6565b61188e565b610b046100c2565b80610b0e81610266565b0390f35b6100c8565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610b5e575b6020831014610b5957565b610b2a565b91607f1691610b4e565b60209181520190565b5f5260205f2090565b905f9291805490610b94610b8d83610b3e565b8094610b68565b916001811690815f14610beb5750600114610baf575b505050565b610bbc9192939450610b71565b915f925b818410610bd357505001905f8080610baa565b60018160209295939554848601520191019290610bc0565b92949550505060ff19168252151560200201905f8080610baa565b90610c1091610b7a565b90565b90610c33610c2c92610c236100c2565b93848092610c06565b038361046b565b565b905f10610c4857610c4590610c13565b90565b610b17565b610c5960025f90610c35565b90565b34610c8c57610c6c3660046100d0565b610c88610c77610c4d565b610c7f6100c2565b91829182610132565b0390f35b6100c8565b5f80fd5b606090565b610ca390610c13565b90565b610cae610c95565b50610cb96002610c9a565b90565b5f90565b5f1c90565b610cd1610cd691610cc0565b6101f1565b90565b610ce39054610cc5565b90565b610cee610cbc565b50610cf96001610cd9565b90565b610d04611899565b610d0c610d39565b565b90565b610d25610d20610d2a92610d0e565b6103c5565b61017f565b90565b610d3690610d11565b90565b610d4a610d455f610d2d565b6118e7565b565b610d54610cfc565b565b610d5e610cbc565b50610d685f610cd9565b90565b90610d7e9291610d79611899565b6110f7565b565b5090565b905f9291805490610d9e610d9783610b3e565b80946100e3565b916001811690815f14610df55750600114610db9575b505050565b610dc69192939450610b71565b915f925b818410610ddd57505001905f8080610db4565b60018160209295939554848601520191019290610dca565b92949550505060ff19168252151560200201905f8080610db4565b9190610e2a81610e2381610e2f956100e3565b8095610547565b6100f7565b0190565b939190610e789593610e5d610e6a93610e5360808901945f8a0190610196565b6020880190610196565b8582036040870152610d84565b926060818503910152610e10565b90565b5f1b90565b90610e9160018060a01b0391610e7b565b9181191691161790565b610eaf610eaa610eb49261017f565b6103c5565b61017f565b90565b610ec090610e9b565b90565b610ecc90610eb7565b90565b90565b90610ee7610ee2610eee92610ec3565b610ecf565b8254610e80565b9055565b601f602091010490565b1b90565b91906008610f1b910291610f155f1984610efc565b92610efc565b9181191691161790565b610f39610f34610f3e926103c2565b6103c5565b6103c2565b90565b90565b9190610f5a610f55610f6293610f25565b610f41565b908354610f00565b9055565b5f90565b610f7c91610f76610f66565b91610f44565b565b5b818110610f8a575050565b80610f975f600193610f6a565b01610f7f565b9190601f8111610fad575b505050565b610fb9610fde93610b71565b906020610fc584610ef2565b83019310610fe6575b610fd790610ef2565b0190610f7e565b5f8080610fa8565b9150610fd781929050610fce565b90611004905f19906008026101ed565b191690565b8161101391610ff4565b906002021790565b916110269082610d80565b9067ffffffffffffffff82116110e55761104a826110448554610b3e565b85610f9d565b5f90601f831160011461107d5791809161106c935f92611071575b5050611009565b90555b565b90915001355f80611065565b601f1983169161108c85610b71565b925f5b8181106110cd575091600293918560019694106110b3575b5050500201905561106f565b6110c3910135601f841690610ff4565b90555f80806110a7565b9193602060018192878701358155019501920161108f565b610457565b906110f5929161101b565b565b91908261111461110e6111095f610d2d565b61018a565b9161018a565b146111bf57611124818390610d80565b61113d6111376111326103e4565b6103c2565b916103c2565b036111a35761119a6111a1936111536001610cd9565b816002916111908688907f1882fd6e997ef1dc5e2691efe78564a5c540b36547363e8a7dcd454c61a944b6956111876100c2565b95869586610e33565b0390a16001610ed2565b60026110ea565b565b5f634ae601c160e11b8152806111bb60048201610266565b0390fd5b5f6303a4d16560e51b8152806111d760048201610266565b0390fd5b906111e69291610d6b565b565b5f90565b5190565b60209181520190565b60200190565b61121e61122760209361122c93611215816100df565b93848093610b68565b958691016100ec565b6100f7565b0190565b9061123a916111ff565b90565b60200190565b90611257611250836111ec565b80926111f0565b9081611268602083028401946111f9565b925f915b83831061127b57505050505090565b9091929394602061129d61129783856001950387528951611230565b9761123d565b930193019193929061126c565b6112bf9160208201915f818403910152611243565b90565b60200190565b5190565b60209181520190565b60200190565b5190565b60209181520190565b611307611310602093611315936112fe816112db565b938480936112df565b958691016100ec565b6100f7565b0190565b90611323916112e8565b90565b60200190565b90611340611339836112c8565b80926112cc565b9081611351602083028401946112d5565b925f915b83831061136457505050505090565b9091929394602061138661138083856001950387528951611319565b97611326565b9301930191939290611355565b6113a89160208201915f81840391015261132c565b90565b5190565b60209181520190565b60200190565b6113c79061018a565b9052565b6113d490610787565b9052565b6113e1906103c2565b9052565b90604080611419936113fd5f8201515f8601906113be565b61140f602082015160208601906113cb565b01519101906113d8565b565b90611428816060936113e5565b0190565b60200190565b9061144f611449611442846113ab565b80936113af565b926113b8565b905f5b81811061145f5750505090565b909192611478611472600192865161141b565b9461142c565b9101919091611452565b6114979160208201915f818403910152611432565b90565b6114a490516104ad565b90565b6114b190516104dd565b90565b6114be90516103c2565b90565b634e487b7160e01b5f52602160045260245ffd5b600d11156114df57565b6114c1565b906114ee826114d5565b565b6114fa90516114e4565b90565b611506906104ad565b9052565b611513906104dd565b9052565b611520906114e4565b90565b61152c90611517565b9052565b60209181520190565b6115586115616020936115669361154f816112db565b93848093611530565b958691016100ec565b6100f7565b0190565b989694929099979593916101608a019a5f8b01611586916114fd565b60208a016115939161150a565b604089016115a09161150a565b606088016115ad9161150a565b608087016115ba9161150a565b60a086016115c79161150a565b60c085016115d49161150a565b60e084016115e1916103fc565b61010083016115ef916103fc565b61012082016115fd91611523565b80820390610140015261160f91611539565b90565b9061161b6111e8565b50611624610ce6565b61163e6116386116335f610d2d565b61018a565b9161018a565b1480156117f8575b6117dc576117ba826117a261167361168260806117bf9701516116676100c2565b928391602083016112aa565b6020820181038252038261046b565b61169461168e826100df565b916112c2565b206117936116b76116c660a08601516116ab6100c2565b92839160208301611393565b6020820181038252038261046b565b6116d86116d2826100df565b916112c2565b20936116f961170860c08301516116ed6100c2565b92839160208301611482565b6020820181038252038261046b565b61171a611714826100df565b916112c2565b20906117275f820161149a565b93611734602083016114a7565b96611741604084016114a7565b9361174e606085016114a7565b929161175c60e086016114b4565b9261176a61010087016114b4565b9461014061177b61012089016114f0565b970151976117876100c2565b9c8d9b60208d0161156a565b6020820181038252038261046b565b6117b46117ae826100df565b916112c2565b2061194a565b611980565b6117d86117d26117cd610ce6565b61018a565b9161018a565b1490565b5f63f64b1d7b60e01b8152806117f460048201610266565b0390fd5b50611809611804610ca6565b6100df565b61182261181c6118176103e4565b6103c2565b916103c2565b1415611646565b61183a90611835611899565b61183c565b565b8061185761185161184c5f610d2d565b61018a565b9161018a565b1461186757611865906118e7565b565b61188a6118735f610d2d565b5f918291631e4fbdf760e01b8352600483016101a3565b0390fd5b61189790611829565b565b6118a1610d56565b6118ba6118b46118af6119a2565b61018a565b9161018a565b036118c157565b6118e36118cc6119a2565b5f91829163118cdaa760e01b8352600483016101a3565b0390fd5b6118f05f610cd9565b6118fa825f610ed2565b9061192e6119287f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093610ec3565b91610ec3565b916119376100c2565b8061194181610266565b0390a3565b5f90565b611952611946565b507f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b61199f9161199691611990610cbc565b506119f2565b90929192611adb565b90565b6119aa610cbc565b503390565b5f90565b90565b6119ca6119c56119cf926119b3565b6103c5565b6103c2565b90565b6119e66119e16119eb926103c2565b610e7b565b6104dd565b90565b5f90565b9190916119fd610cbc565b50611a066119af565b50611a0f611946565b50611a19836100df565b611a2c611a2660416119b6565b916103c2565b145f14611a7357611a6c9192611a40611946565b50611a49611946565b50611a526119ee565b506020810151606060408301519201515f1a909192611c48565b9192909190565b50611a7d5f610d2d565b90611a91611a8c6002946100df565b6119d2565b91929190565b60041115611aa157565b6114c1565b90611ab082611a97565b565b9190611ac5905f6020850194019061150a565b565b611ad3611ad891610cc0565b610f25565b90565b80611aee611ae85f611aa6565b91611aa6565b145f14611af9575050565b80611b0d611b076001611aa6565b91611aa6565b145f14611b30575f63f645eedf60e01b815280611b2c60048201610266565b0390fd5b80611b44611b3e6002611aa6565b91611aa6565b145f14611b7257611b6e611b5783611ac7565b5f91829163fce698f760e01b835260048301610409565b0390fd5b611b85611b7f6003611aa6565b91611aa6565b14611b8d5750565b611ba8905f9182916335e2f38360e21b835260048301611ab2565b0390fd5b90565b611bc3611bbe611bc892611bac565b6103c5565b6103c2565b90565b60ff1690565b611bda90611bcb565b9052565b611c13611c1a94611c09606094989795611bff608086019a5f87019061150a565b6020850190611bd1565b604083019061150a565b019061150a565b565b611c246100c2565b3d5f823e3d90fd5b611c40611c3b611c4592610d0e565b610e7b565b6104dd565b90565b939293611c53610cbc565b50611c5c6119af565b50611c65611946565b50611c6f85611ac7565b611ca1611c9b7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0611baf565b916103c2565b11611d2e5790611cc4602094955f94939293611cbb6100c2565b94859485611bde565b838052039060015afa15611d2957611cdc5f51610e7b565b80611cf7611cf1611cec5f610d2d565b61018a565b9161018a565b14611d0d575f91611d075f611c2c565b91929190565b50611d175f610d2d565b600191611d235f611c2c565b91929190565b611c1c565b505050611d3a5f610d2d565b906003929192919056fea2646970667358221220e8468785bc7cb48d345aea56106f80aa25b3c963ad3f6277d7522a8928baa7c564736f6c634300081e0033",
}

// NoAttestationTeeAuthenticator is an auto generated Go binding around an Ethereum contract.
type NoAttestationTeeAuthenticator struct {
	abi abi.ABI
}

// NewNoAttestationTeeAuthenticator creates a new instance of NoAttestationTeeAuthenticator.
func NewNoAttestationTeeAuthenticator() *NoAttestationTeeAuthenticator {
	parsed, err := NoAttestationTeeAuthenticatorMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &NoAttestationTeeAuthenticator{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *NoAttestationTeeAuthenticator) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address owner, address _teeSigner, bytes _pubSecp521r1) returns()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackConstructor(owner common.Address, _teeSigner common.Address, _pubSecp521r1 []byte) []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("", owner, _teeSigner, _pubSecp521r1)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackPKLENGTH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc91496c6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackPKLENGTH() []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("PK_LENGTH")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPKLENGTH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc91496c6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackPKLENGTH() ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("PK_LENGTH")
}

// UnpackPKLENGTH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc91496c6.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackPKLENGTH(data []byte) (*big.Int, error) {
	out, err := noAttestationTeeAuthenticator.abi.Unpack("PK_LENGTH", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd119a29d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd119a29d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd119a29d.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := noAttestationTeeAuthenticator.abi.Unpack("checkSignature", data)
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackGetPubSecp521r1() []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("getPubSecp521r1")
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackGetPubSecp521r1() ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("getPubSecp521r1")
}

// UnpackGetPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081bec7e.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackGetPubSecp521r1(data []byte) ([]byte, error) {
	out, err := noAttestationTeeAuthenticator.abi.Unpack("getPubSecp521r1", data)
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackGetTeeSigner() []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("getTeeSigner")
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackGetTeeSigner() ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("getTeeSigner")
}

// UnpackGetTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dd7ce2f.
//
// Solidity: function getTeeSigner() view returns(address)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackGetTeeSigner(data []byte) (common.Address, error) {
	out, err := noAttestationTeeAuthenticator.abi.Unpack("getTeeSigner", data)
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackOwner() []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("owner")
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackOwner() ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackOwner(data []byte) (common.Address, error) {
	out, err := noAttestationTeeAuthenticator.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPubSecp521r1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe4993ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pubSecp521r1() view returns(bytes)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackPubSecp521r1() []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("pubSecp521r1")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPubSecp521r1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe4993ca.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pubSecp521r1() view returns(bytes)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackPubSecp521r1() ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("pubSecp521r1")
}

// UnpackPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfe4993ca.
//
// Solidity: function pubSecp521r1() view returns(bytes)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackPubSecp521r1(data []byte) ([]byte, error) {
	out, err := noAttestationTeeAuthenticator.abi.Unpack("pubSecp521r1", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackRenounceOwnership() []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("renounceOwnership")
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackRenounceOwnership() ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("renounceOwnership")
}

// PackTeeSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x43f855c3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function teeSigner() view returns(address)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackTeeSigner() []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("teeSigner")
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackTeeSigner() ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("teeSigner")
}

// UnpackTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x43f855c3.
//
// Solidity: function teeSigner() view returns(address)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackTeeSigner(data []byte) (common.Address, error) {
	out, err := noAttestationTeeAuthenticator.abi.Unpack("teeSigner", data)
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("transferOwnership", newOwner)
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("transferOwnership", newOwner)
}

// PackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc64af6fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTee(address newTeeSigner, bytes newPubSecp521r1) returns()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) PackUpdateTee(newTeeSigner common.Address, newPubSecp521r1 []byte) []byte {
	enc, err := noAttestationTeeAuthenticator.abi.Pack("updateTee", newTeeSigner, newPubSecp521r1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc64af6fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTee(address newTeeSigner, bytes newPubSecp521r1) returns()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) TryPackUpdateTee(newTeeSigner common.Address, newPubSecp521r1 []byte) ([]byte, error) {
	return noAttestationTeeAuthenticator.abi.Pack("updateTee", newTeeSigner, newPubSecp521r1)
}

// NoAttestationTeeAuthenticatorOwnershipTransferred represents a OwnershipTransferred event raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const NoAttestationTeeAuthenticatorOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (NoAttestationTeeAuthenticatorOwnershipTransferred) ContractEventName() string {
	return NoAttestationTeeAuthenticatorOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackOwnershipTransferredEvent(log *types.Log) (*NoAttestationTeeAuthenticatorOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if log.Topics[0] != noAttestationTeeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(NoAttestationTeeAuthenticatorOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range noAttestationTeeAuthenticator.abi.Events[event].Inputs {
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

// NoAttestationTeeAuthenticatorTeeUpdate represents a TeeUpdate event raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorTeeUpdate struct {
	OldTee          common.Address
	NewTee          common.Address
	OldPubSecp521r1 []byte
	NewPubSecp521r1 []byte
	Raw             *types.Log // Blockchain specific contextual infos
}

const NoAttestationTeeAuthenticatorTeeUpdateEventName = "TeeUpdate"

// ContractEventName returns the user-defined event name.
func (NoAttestationTeeAuthenticatorTeeUpdate) ContractEventName() string {
	return NoAttestationTeeAuthenticatorTeeUpdateEventName
}

// UnpackTeeUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackTeeUpdateEvent(log *types.Log) (*NoAttestationTeeAuthenticatorTeeUpdate, error) {
	event := "TeeUpdate"
	if log.Topics[0] != noAttestationTeeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(NoAttestationTeeAuthenticatorTeeUpdate)
	if len(log.Data) > 0 {
		if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range noAttestationTeeAuthenticator.abi.Events[event].Inputs {
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
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["InvalidPKLength"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackInvalidPKLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["TeeAddressCantBeZero"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackTeeAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], noAttestationTeeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return noAttestationTeeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// NoAttestationTeeAuthenticatorECDSAInvalidSignature represents a ECDSAInvalidSignature error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorECDSAInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignature()
func NoAttestationTeeAuthenticatorECDSAInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0xf645eedf0193584640b6b90cb9477e4c95b98636c148a891d4c0a146dc46e75a")
}

// UnpackECDSAInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignature()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackECDSAInvalidSignatureError(raw []byte) (*NoAttestationTeeAuthenticatorECDSAInvalidSignature, error) {
	out := new(NoAttestationTeeAuthenticatorECDSAInvalidSignature)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// NoAttestationTeeAuthenticatorECDSAInvalidSignatureLength represents a ECDSAInvalidSignatureLength error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorECDSAInvalidSignatureLength struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func NoAttestationTeeAuthenticatorECDSAInvalidSignatureLengthErrorID() common.Hash {
	return common.HexToHash("0xfce698f7e8e5342cd615f641317bc45fe7e1e4a8b0a14dd1383ff8dc9c41917f")
}

// UnpackECDSAInvalidSignatureLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackECDSAInvalidSignatureLengthError(raw []byte) (*NoAttestationTeeAuthenticatorECDSAInvalidSignatureLength, error) {
	out := new(NoAttestationTeeAuthenticatorECDSAInvalidSignatureLength)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// NoAttestationTeeAuthenticatorECDSAInvalidSignatureS represents a ECDSAInvalidSignatureS error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorECDSAInvalidSignatureS struct {
	S [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func NoAttestationTeeAuthenticatorECDSAInvalidSignatureSErrorID() common.Hash {
	return common.HexToHash("0xd78bce0cccb935155ed6428d1c13e50b7f3550fd2b66b9fe266006fea4a5e1eb")
}

// UnpackECDSAInvalidSignatureSError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackECDSAInvalidSignatureSError(raw []byte) (*NoAttestationTeeAuthenticatorECDSAInvalidSignatureS, error) {
	out := new(NoAttestationTeeAuthenticatorECDSAInvalidSignatureS)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// NoAttestationTeeAuthenticatorInvalidPKLength represents a InvalidPKLength error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorInvalidPKLength struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPKLength()
func NoAttestationTeeAuthenticatorInvalidPKLengthErrorID() common.Hash {
	return common.HexToHash("0x95cc0382d99b2fba4335843859b8b7e3894ae70a43b516ae2d581373538dd6ac")
}

// UnpackInvalidPKLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPKLength()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackInvalidPKLengthError(raw []byte) (*NoAttestationTeeAuthenticatorInvalidPKLength, error) {
	out := new(NoAttestationTeeAuthenticatorInvalidPKLength)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "InvalidPKLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// NoAttestationTeeAuthenticatorOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func NoAttestationTeeAuthenticatorOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackOwnableInvalidOwnerError(raw []byte) (*NoAttestationTeeAuthenticatorOwnableInvalidOwner, error) {
	out := new(NoAttestationTeeAuthenticatorOwnableInvalidOwner)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// NoAttestationTeeAuthenticatorOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func NoAttestationTeeAuthenticatorOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackOwnableUnauthorizedAccountError(raw []byte) (*NoAttestationTeeAuthenticatorOwnableUnauthorizedAccount, error) {
	out := new(NoAttestationTeeAuthenticatorOwnableUnauthorizedAccount)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// NoAttestationTeeAuthenticatorTeeAddressCantBeZero represents a TeeAddressCantBeZero error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorTeeAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TeeAddressCantBeZero()
func NoAttestationTeeAuthenticatorTeeAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x749a2ca096befde4321f5c4e6ddb95f574ec95d21a2dead6e38f622a10fd68c7")
}

// UnpackTeeAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TeeAddressCantBeZero()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackTeeAddressCantBeZeroError(raw []byte) (*NoAttestationTeeAuthenticatorTeeAddressCantBeZero, error) {
	out := new(NoAttestationTeeAuthenticatorTeeAddressCantBeZero)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "TeeAddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// NoAttestationTeeAuthenticatorTeeIsNotSet represents a TeeIsNotSet error raised by the NoAttestationTeeAuthenticator contract.
type NoAttestationTeeAuthenticatorTeeIsNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TeeIsNotSet()
func NoAttestationTeeAuthenticatorTeeIsNotSetErrorID() common.Hash {
	return common.HexToHash("0xf64b1d7b3df6940d5da676c1cc07a6144e620b0858b161be25b6cd0ef7569425")
}

// UnpackTeeIsNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TeeIsNotSet()
func (noAttestationTeeAuthenticator *NoAttestationTeeAuthenticator) UnpackTeeIsNotSetError(raw []byte) (*NoAttestationTeeAuthenticatorTeeIsNotSet, error) {
	out := new(NoAttestationTeeAuthenticatorTeeIsNotSet)
	if err := noAttestationTeeAuthenticator.abi.UnpackIntoInterface(out, "TeeIsNotSet", raw); err != nil {
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

// PanicMetaData contains all meta data concerning the Panic contract.
var PanicMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "d7d9ad298c70ed58e3eba1456b2d564fae",
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea2646970667358221220da179cfea467b0d4c9e5605075071a59d3dcf88c0a22d4d7cf66e0182db50a8464736f6c634300081e0033",
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
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea264697066735822122059c7a309fec31bd42cf0c293c64f0123f48b490fcf1975b6baffd39d9f5e160d64736f6c634300081e0033",
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
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea2646970667358221220c088b7916fc9b32346e8bfbaacfcf0271d9d03a2f26667225831222a1ac394d664736f6c634300081e0033",
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
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea26469706673582212205c6302828aab4db764a43a2d32f0aaafc6b417434191febd2bfe757f461f780c64736f6c634300081e0033",
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
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea26469706673582212209ffbba52071cd4909dd91df41f3e001548c66a43f39312767f50526acb30e24064736f6c634300081e0033",
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
