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
	Receiver common.Address
	Amount   *big.Int
}

// AbstractTeeAuthenticatorMetaData contains all meta data concerning the AbstractTeeAuthenticator contract.
var AbstractTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0xce6dce51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := abstractTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce6dce51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return abstractTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce6dce51.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
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
	Bin: "0x6080806040523460175760399081601c823930815050f35b5f80fdfe5f80fdfea2646970667358221220b91d193e70e5dbfd828d394bfb366b6606e3981648eaacd43ec8b8e9abfa94b364736f6c634300081e0033",
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

// INitroProverMetaData contains all meta data concerning the INitroProver contract.
var INitroProverMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxAge\",\"type\":\"uint256\"}],\"name\":\"verifyAttestation\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"max_age\",\"type\":\"uint256\"}],\"name\":\"verifyAttestationStep1\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"attestation_decoded\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"certificate\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"cabundle\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"enclaveKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"userData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rawPcrs\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"cabundle\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"parentPubKey\",\"type\":\"bytes\"}],\"name\":\"verifyAttestationStep2\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"pubKey\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"attestation_decoded\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"certificate\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"parentPubKey\",\"type\":\"bytes\"}],\"name\":\"verifyAttestationStep3\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"attestationSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pubKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"bufBufBuf\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestationSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pubKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"buf\",\"type\":\"bytes\"}],\"name\":\"verifyAttestationStep4\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "58a55af41d95bda4586bda5febb361b1bf",
}

// INitroProver is an auto generated Go binding around an Ethereum contract.
type INitroProver struct {
	abi abi.ABI
}

// NewINitroProver creates a new instance of INitroProver.
func NewINitroProver() *INitroProver {
	parsed, err := INitroProverMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &INitroProver{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *INitroProver) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackVerifyAttestation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b5f2f81.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyAttestation(bytes attestation, uint256 maxAge) view returns(bytes, bytes, bytes)
func (iNitroProver *INitroProver) PackVerifyAttestation(attestation []byte, maxAge *big.Int) []byte {
	enc, err := iNitroProver.abi.Pack("verifyAttestation", attestation, maxAge)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyAttestation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b5f2f81.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyAttestation(bytes attestation, uint256 maxAge) view returns(bytes, bytes, bytes)
func (iNitroProver *INitroProver) TryPackVerifyAttestation(attestation []byte, maxAge *big.Int) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestation", attestation, maxAge)
}

// VerifyAttestationOutput serves as a container for the return parameters of contract
// method VerifyAttestation.
type VerifyAttestationOutput struct {
	Arg0 []byte
	Arg1 []byte
	Arg2 []byte
}

// UnpackVerifyAttestation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2b5f2f81.
//
// Solidity: function verifyAttestation(bytes attestation, uint256 maxAge) view returns(bytes, bytes, bytes)
func (iNitroProver *INitroProver) UnpackVerifyAttestation(data []byte) (VerifyAttestationOutput, error) {
	out, err := iNitroProver.abi.Unpack("verifyAttestation", data)
	outstruct := new(VerifyAttestationOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Arg2 = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	return *outstruct, nil
}

// PackVerifyAttestationStep1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd420f66a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyAttestationStep1(bytes attestation, uint256 max_age) view returns(bytes[] attestation_decoded, bytes certificate, bytes[] cabundle, bytes enclaveKey, bytes userData, bytes rawPcrs)
func (iNitroProver *INitroProver) PackVerifyAttestationStep1(attestation []byte, maxAge *big.Int) []byte {
	enc, err := iNitroProver.abi.Pack("verifyAttestationStep1", attestation, maxAge)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyAttestationStep1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd420f66a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyAttestationStep1(bytes attestation, uint256 max_age) view returns(bytes[] attestation_decoded, bytes certificate, bytes[] cabundle, bytes enclaveKey, bytes userData, bytes rawPcrs)
func (iNitroProver *INitroProver) TryPackVerifyAttestationStep1(attestation []byte, maxAge *big.Int) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestationStep1", attestation, maxAge)
}

// VerifyAttestationStep1Output serves as a container for the return parameters of contract
// method VerifyAttestationStep1.
type VerifyAttestationStep1Output struct {
	AttestationDecoded [][]byte
	Certificate        []byte
	Cabundle           [][]byte
	EnclaveKey         []byte
	UserData           []byte
	RawPcrs            []byte
}

// UnpackVerifyAttestationStep1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd420f66a.
//
// Solidity: function verifyAttestationStep1(bytes attestation, uint256 max_age) view returns(bytes[] attestation_decoded, bytes certificate, bytes[] cabundle, bytes enclaveKey, bytes userData, bytes rawPcrs)
func (iNitroProver *INitroProver) UnpackVerifyAttestationStep1(data []byte) (VerifyAttestationStep1Output, error) {
	out, err := iNitroProver.abi.Unpack("verifyAttestationStep1", data)
	outstruct := new(VerifyAttestationStep1Output)
	if err != nil {
		return *outstruct, err
	}
	outstruct.AttestationDecoded = *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)
	outstruct.Certificate = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Cabundle = *abi.ConvertType(out[2], new([][]byte)).(*[][]byte)
	outstruct.EnclaveKey = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.UserData = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.RawPcrs = *abi.ConvertType(out[5], new([]byte)).(*[]byte)
	return *outstruct, nil
}

// PackVerifyAttestationStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc13734e7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyAttestationStep2(bytes[] cabundle, uint256 index, bytes parentPubKey) view returns(bytes pubKey)
func (iNitroProver *INitroProver) PackVerifyAttestationStep2(cabundle [][]byte, index *big.Int, parentPubKey []byte) []byte {
	enc, err := iNitroProver.abi.Pack("verifyAttestationStep2", cabundle, index, parentPubKey)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyAttestationStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc13734e7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyAttestationStep2(bytes[] cabundle, uint256 index, bytes parentPubKey) view returns(bytes pubKey)
func (iNitroProver *INitroProver) TryPackVerifyAttestationStep2(cabundle [][]byte, index *big.Int, parentPubKey []byte) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestationStep2", cabundle, index, parentPubKey)
}

// UnpackVerifyAttestationStep2 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc13734e7.
//
// Solidity: function verifyAttestationStep2(bytes[] cabundle, uint256 index, bytes parentPubKey) view returns(bytes pubKey)
func (iNitroProver *INitroProver) UnpackVerifyAttestationStep2(data []byte) ([]byte, error) {
	out, err := iNitroProver.abi.Unpack("verifyAttestationStep2", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackVerifyAttestationStep3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d522337.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyAttestationStep3(bytes[] attestation_decoded, bytes certificate, bytes parentPubKey) view returns(bytes attestationSig, bytes pubKey, bytes bufBufBuf)
func (iNitroProver *INitroProver) PackVerifyAttestationStep3(attestationDecoded [][]byte, certificate []byte, parentPubKey []byte) []byte {
	enc, err := iNitroProver.abi.Pack("verifyAttestationStep3", attestationDecoded, certificate, parentPubKey)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyAttestationStep3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d522337.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyAttestationStep3(bytes[] attestation_decoded, bytes certificate, bytes parentPubKey) view returns(bytes attestationSig, bytes pubKey, bytes bufBufBuf)
func (iNitroProver *INitroProver) TryPackVerifyAttestationStep3(attestationDecoded [][]byte, certificate []byte, parentPubKey []byte) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestationStep3", attestationDecoded, certificate, parentPubKey)
}

// VerifyAttestationStep3Output serves as a container for the return parameters of contract
// method VerifyAttestationStep3.
type VerifyAttestationStep3Output struct {
	AttestationSig []byte
	PubKey         []byte
	BufBufBuf      []byte
}

// UnpackVerifyAttestationStep3 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5d522337.
//
// Solidity: function verifyAttestationStep3(bytes[] attestation_decoded, bytes certificate, bytes parentPubKey) view returns(bytes attestationSig, bytes pubKey, bytes bufBufBuf)
func (iNitroProver *INitroProver) UnpackVerifyAttestationStep3(data []byte) (VerifyAttestationStep3Output, error) {
	out, err := iNitroProver.abi.Unpack("verifyAttestationStep3", data)
	outstruct := new(VerifyAttestationStep3Output)
	if err != nil {
		return *outstruct, err
	}
	outstruct.AttestationSig = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.PubKey = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.BufBufBuf = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	return *outstruct, nil
}

// PackVerifyAttestationStep4 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5aefbb26.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyAttestationStep4(bytes attestationSig, bytes pubKey, bytes buf) view returns()
func (iNitroProver *INitroProver) PackVerifyAttestationStep4(attestationSig []byte, pubKey []byte, buf []byte) []byte {
	enc, err := iNitroProver.abi.Pack("verifyAttestationStep4", attestationSig, pubKey, buf)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyAttestationStep4 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5aefbb26.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyAttestationStep4(bytes attestationSig, bytes pubKey, bytes buf) view returns()
func (iNitroProver *INitroProver) TryPackVerifyAttestationStep4(attestationSig []byte, pubKey []byte, buf []byte) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestationStep4", attestationSig, pubKey, buf)
}

// ITeeAuthenticatorMetaData contains all meta data concerning the ITeeAuthenticator contract.
var ITeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0xce6dce51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce6dce51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce6dce51.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
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

// ITeeAuthenticatorAdminMetaData contains all meta data concerning the ITeeAuthenticatorAdmin contract.
var ITeeAuthenticatorAdminMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AttestationAlreadyUsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPCR\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUserDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStep\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"oldPcr0\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"PcrZeroUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getStep2TotalLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"updatePcr0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTeeStep1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "e4eb3eef93edc28b740127688b159318ad",
}

// ITeeAuthenticatorAdmin is an auto generated Go binding around an Ethereum contract.
type ITeeAuthenticatorAdmin struct {
	abi abi.ABI
}

// NewITeeAuthenticatorAdmin creates a new instance of ITeeAuthenticatorAdmin.
func NewITeeAuthenticatorAdmin() *ITeeAuthenticatorAdmin {
	parsed, err := ITeeAuthenticatorAdminMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ITeeAuthenticatorAdmin{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ITeeAuthenticatorAdmin) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGetStep2TotalLength is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3b33122.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getStep2TotalLength() view returns(uint256)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackGetStep2TotalLength() []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("getStep2TotalLength")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetStep2TotalLength is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3b33122.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getStep2TotalLength() view returns(uint256)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackGetStep2TotalLength() ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("getStep2TotalLength")
}

// UnpackGetStep2TotalLength is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3b33122.
//
// Solidity: function getStep2TotalLength() view returns(uint256)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackGetStep2TotalLength(data []byte) (*big.Int, error) {
	out, err := iTeeAuthenticatorAdmin.abi.Unpack("getStep2TotalLength", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUpdatePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ed2d7fc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updatePcr0(bytes newPcr0) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackUpdatePcr0(newPcr0 []byte) []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("updatePcr0", newPcr0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdatePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ed2d7fc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updatePcr0(bytes newPcr0) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackUpdatePcr0(newPcr0 []byte) ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("updatePcr0", newPcr0)
}

// PackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7637d58a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTee(bytes attestation) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackUpdateTee(attestation []byte) []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("updateTee", attestation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7637d58a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTee(bytes attestation) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackUpdateTee(attestation []byte) ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("updateTee", attestation)
}

// PackUpdateTeeStep1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8890b84e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep1(bytes attestation) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackUpdateTeeStep1(attestation []byte) []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep1", attestation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8890b84e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep1(bytes attestation) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackUpdateTeeStep1(attestation []byte) ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep1", attestation)
}

// PackUpdateTeeStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5520916c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep2() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackUpdateTeeStep2() []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep2")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5520916c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep2() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackUpdateTeeStep2() ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep2")
}

// PackUpdateTeeStep3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd22a29d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep3() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackUpdateTeeStep3() []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep3")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd22a29d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep3() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackUpdateTeeStep3() ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep3")
}

// PackUpdateTeeStep4 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34010ccf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep4() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackUpdateTeeStep4() []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep4")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep4 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34010ccf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep4() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackUpdateTeeStep4() ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("updateTeeStep4")
}

// ITeeAuthenticatorAdminPcrZeroUpdate represents a PcrZeroUpdate event raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminPcrZeroUpdate struct {
	OldPcr0 common.Hash
	NewPcr0 common.Hash
	Raw     *types.Log // Blockchain specific contextual infos
}

const ITeeAuthenticatorAdminPcrZeroUpdateEventName = "PcrZeroUpdate"

// ContractEventName returns the user-defined event name.
func (ITeeAuthenticatorAdminPcrZeroUpdate) ContractEventName() string {
	return ITeeAuthenticatorAdminPcrZeroUpdateEventName
}

// UnpackPcrZeroUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PcrZeroUpdate(bytes indexed oldPcr0, bytes indexed newPcr0)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackPcrZeroUpdateEvent(log *types.Log) (*ITeeAuthenticatorAdminPcrZeroUpdate, error) {
	event := "PcrZeroUpdate"
	if log.Topics[0] != iTeeAuthenticatorAdmin.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ITeeAuthenticatorAdminPcrZeroUpdate)
	if len(log.Data) > 0 {
		if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iTeeAuthenticatorAdmin.abi.Events[event].Inputs {
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

// ITeeAuthenticatorAdminTeeUpdate represents a TeeUpdate event raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminTeeUpdate struct {
	OldTee          common.Address
	NewTee          common.Address
	OldPubSecp521r1 []byte
	NewPubSecp521r1 []byte
	Raw             *types.Log // Blockchain specific contextual infos
}

const ITeeAuthenticatorAdminTeeUpdateEventName = "TeeUpdate"

// ContractEventName returns the user-defined event name.
func (ITeeAuthenticatorAdminTeeUpdate) ContractEventName() string {
	return ITeeAuthenticatorAdminTeeUpdateEventName
}

// UnpackTeeUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackTeeUpdateEvent(log *types.Log) (*ITeeAuthenticatorAdminTeeUpdate, error) {
	event := "TeeUpdate"
	if log.Topics[0] != iTeeAuthenticatorAdmin.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ITeeAuthenticatorAdminTeeUpdate)
	if len(log.Data) > 0 {
		if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iTeeAuthenticatorAdmin.abi.Events[event].Inputs {
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
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["AttestationAlreadyUsed"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackAttestationAlreadyUsedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["InvalidPCR"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackInvalidPCRError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["InvalidPKLength"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackInvalidPKLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["InvalidUserDataLength"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackInvalidUserDataLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["WrongStep"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackWrongStepError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ITeeAuthenticatorAdminAttestationAlreadyUsed represents a AttestationAlreadyUsed error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminAttestationAlreadyUsed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AttestationAlreadyUsed()
func ITeeAuthenticatorAdminAttestationAlreadyUsedErrorID() common.Hash {
	return common.HexToHash("0x17ee279c4380943995e5de7b6ca148f56c43f727fc1a52fa3f267e6f5ea97559")
}

// UnpackAttestationAlreadyUsedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AttestationAlreadyUsed()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackAttestationAlreadyUsedError(raw []byte) (*ITeeAuthenticatorAdminAttestationAlreadyUsed, error) {
	out := new(ITeeAuthenticatorAdminAttestationAlreadyUsed)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "AttestationAlreadyUsed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminInvalidPCR represents a InvalidPCR error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminInvalidPCR struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPCR()
func ITeeAuthenticatorAdminInvalidPCRErrorID() common.Hash {
	return common.HexToHash("0x719a9a4917a12c04cb8a46a64960d788ba098ccc7ac543fcee9fa361a6d1224c")
}

// UnpackInvalidPCRError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPCR()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackInvalidPCRError(raw []byte) (*ITeeAuthenticatorAdminInvalidPCR, error) {
	out := new(ITeeAuthenticatorAdminInvalidPCR)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "InvalidPCR", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminInvalidPKLength represents a InvalidPKLength error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminInvalidPKLength struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPKLength()
func ITeeAuthenticatorAdminInvalidPKLengthErrorID() common.Hash {
	return common.HexToHash("0x95cc0382d99b2fba4335843859b8b7e3894ae70a43b516ae2d581373538dd6ac")
}

// UnpackInvalidPKLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPKLength()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackInvalidPKLengthError(raw []byte) (*ITeeAuthenticatorAdminInvalidPKLength, error) {
	out := new(ITeeAuthenticatorAdminInvalidPKLength)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "InvalidPKLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminInvalidUserDataLength represents a InvalidUserDataLength error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminInvalidUserDataLength struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidUserDataLength()
func ITeeAuthenticatorAdminInvalidUserDataLengthErrorID() common.Hash {
	return common.HexToHash("0x85edc5464c5565960356d11f255a83de2f3ba3c24b8485a39f37b6c98eec54bc")
}

// UnpackInvalidUserDataLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidUserDataLength()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackInvalidUserDataLengthError(raw []byte) (*ITeeAuthenticatorAdminInvalidUserDataLength, error) {
	out := new(ITeeAuthenticatorAdminInvalidUserDataLength)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "InvalidUserDataLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminWrongStep represents a WrongStep error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminWrongStep struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongStep()
func ITeeAuthenticatorAdminWrongStepErrorID() common.Hash {
	return common.HexToHash("0xbc3c8a5a94b27f160f1b92f62d402d946766caac5227a2ffc2d1381f57f4e236")
}

// UnpackWrongStepError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongStep()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackWrongStepError(raw []byte) (*ITeeAuthenticatorAdminWrongStep, error) {
	out := new(ITeeAuthenticatorAdminWrongStep)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "WrongStep", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MathMetaData contains all meta data concerning the Math contract.
var MathMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "362817aca504bb11799e06f16c73252206",
	Bin: "0x6080806040523460175760399081601c823930815050f35b5f80fdfe5f80fdfea2646970667358221220b4039baeb46cbedbf2b03f50ac8b0928d62b6c63f0a2c7804873e4f6a22f9cc964736f6c634300081e0033",
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
	Bin: "0x6080806040523460175760399081601c823930815050f35b5f80fdfe5f80fdfea264697066735822122029fc2269dd3eea82730e1239de4e6a1c41a4fad72a9eaaecfffb6be07a4c9a7c64736f6c634300081e0033",
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
	Bin: "0x6080806040523460175760399081601c823930815050f35b5f80fdfe5f80fdfea2646970667358221220a1477739f416da8c6baa983aa936ec2afb239244709d8ba052c150073b6a998964736f6c634300081e0033",
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
	Bin: "0x6080806040523460175760399081601c823930815050f35b5f80fdfe5f80fdfea2646970667358221220464b6c1af084aa60dd0f1a27931b60cd69167029a524fb8466c24012a15f4eb164736f6c634300081e0033",
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
	Bin: "0x6080806040523460175760399081601c823930815050f35b5f80fdfe5f80fdfea2646970667358221220732d1b6be92dbde5a2041a604ef8bb2326d696fc4bab1dd79cb3b863a138685c64736f6c634300081e0033",
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
	Bin: "0x6080806040523460175760399081601c823930815050f35b5f80fdfe5f80fdfea2646970667358221220a2eb66ff6a191fa67ac9819c5cc9e2769b8eb1d1a1ac5472d4ad19ed55ffc0c164736f6c634300081e0033",
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
	Bin: "0x608080604052346013576039908160188239f35b5f80fdfe5f80fdfea2646970667358221220b00379a8119a529006184d18b49e27af71fe6ea7d5a7770df0e68409772be6e464736f6c634300081e0033",
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"contractINitroProver\",\"name\":\"_nitroProver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pcr0\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_maxVerificationAge\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AttestationAlreadyUsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPCR\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUserDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStep\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"oldPcr0\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"PcrZeroUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentUpdateStep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStep2TotalLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxVerificationAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nitroProver\",\"outputs\":[{\"internalType\":\"contractINitroProver\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pcr0\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step2CurrentIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"updatePcr0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTeeStep1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a8d895d4d49ffcfa5f7559dd61ad47cda6",
	Bin: "0x60c0604052346102d1576130b880380380610019816102d5565b9283398101906080818303126102d15780516001600160a01b03811691908290036102d1576020810151906001600160a01b03821682036102d15760408101516001600160401b0381116102d15781019084601f830112156102d15781516001600160401b0381116102aa57610098601f8201601f19166020016102d5565b95818752602082850101116102d1576020815f928260609601838a015e8701015201519180156102be575f80546001600160a01b03198116831782556001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a382516001600160401b0381116102aa57600154600181811c911680156102a0575b602082101461028c57601f8111610229575b506020601f82116001146101c657819293945f926101bb575b50508160011b915f199060031b1c1916176001555b60805260a052604051612dbd90816102fb82396080518181816102bb01528181610dee01528181611bd401528181611da101528181611f760152612299015260a051818181610dc001528181611ba6015261221c0152f35b015190505f8061014e565b601f1982169060015f52805f20915f5b818110610211575095836001959697106101f9575b505050811b01600155610163565b01515f1960f88460031b161c191690555f80806101eb565b9192602060018192868b0151815501940192016101d6565b60015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6601f830160051c81019160208410610282575b601f0160051c01905b8181106102775750610135565b5f815560010161026a565b9091508190610261565b634e487b7160e01b5f52602260045260245ffd5b90607f1690610123565b634e487b7160e01b5f52604160045260245ffd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b6040519190601f01601f191682016001600160401b038111838210176102aa5760405256fe6080806040526004361015610012575f80fd5b5f905f3560e01c908163081bec7e146122c8575080630a83fb91146122845780630b1bfab9146122675780630dd7ce2f1461223f578063127379db1461220557806334010ccf14611f5157806343f855c314611f285780635520916c14611d07578063715018a614611cad5780637637d58a14611b4c5780637ed2d7fc1461197f57806381a9d38a146118c35780638890b84e14610d5b5780638da5cb5b14610d34578063c3b3312214610d16578063c91496c614610cfa578063ce6dce5114610735578063d22a29d1146101f6578063db6a7b71146101d8578063f2fde38b1461014a5763fe4993ca14610105575f80fd5b3461014757806003193601126101475761014360405161012f8161012881612394565b0382612537565b6040519182916020835260208301906122e9565b0390f35b80fd5b5034610147576020366003190112610147576004356001600160a01b038116908190036101d457610179612841565b80156101c05781546001600160a01b03198116821783556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b631e4fbdf760e01b82526004829052602482fd5b5080fd5b50346101475780600319360112610147576020600654604051908152f35b503461014757806003193601126101475761020f612841565b60026005540361072657604051635d52233760e01b815260606004820152600f805460648301819052908352829082906084600582901b83018101915f516020612d685f395f51905f5291859085015b8282106106f35750505050818103600319016024830152600d5483916102848261235c565b91828252846001821691825f146106d2575050600114610672575b50506102b782916003198382030160448401526125f2565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa91821561066657808092819461063c575b5083516001600160401b0381116104ec5761031260105461235c565b601f81116105ee575b50602094601f821160011461056d5782939495829161034e9492610409575b50508160011b915f199060031b1c19161790565b6010555b82516001600160401b0381116104ec576103768161037160075461235c565b6126be565b6020601f82116001146105005781908394956103a594926104095750508160011b915f199060031b1c19161790565b6007555b81516001600160401b0381116104ec576103c4600c5461235c565b601f8111610493575b50602092601f821160011461041457829382916103fd94926104095750508160011b915f199060031b1c19161790565b600c555b600360055580f35b015190505f8061033a565b600c8352601f198216937fdf6966c971051c3d54ec59162606531493a51404a002842f56009d7e5cf4a8c791845b86811061047b5750836001959610610463575b505050811b01600c55610401565b01515f1960f88460031b161c191690555f8080610455565b91926020600181928685015181550194019201610442565b600c83526104dc907fdf6966c971051c3d54ec59162606531493a51404a002842f56009d7e5cf4a8c7601f840160051c810191602085106104e2575b601f0160051c01906126a8565b5f6103cd565b90915081906104cf565b634e487b7160e01b82526041600452602482fd5b600783525f516020612d285f395f51905f5290601f198316845b8181106105555750958360019596971061053d575b505050811b016007556103a9565b01515f1960f88460031b161c191690555f808061052f565b9192602060018192868b01518155019401920161051a565b60108352601f198216957f1b6847dc741a1b0cd08d278845f9d819d87b734759afb55fe2de5cb82a9ae67291845b8881106105d6575083600195969798106105be575b505050811b01601055610352565b01515f1960f88460031b161c191690555f80806105b0565b9192602060018192868501518155019401920161059b565b60108352610636907f1b6847dc741a1b0cd08d278845f9d819d87b734759afb55fe2de5cb82a9ae672601f840160051c810191602085106104e257601f0160051c01906126a8565b5f61031b565b9150925061065c91503d8084833e6106548183612537565b810190612730565b919091925f6102f6565b604051903d90823e3d90fd5b600d8552849250907fd7b6990105719101dabeb77144f2a3385c8033acd3af97e9423a695e81ad1eb55b8184106106b2575050016020016102b78261029f565b60019195506020929450805483858701015201910190918593859361069c565b6102b794919550602093925060ff191683830152151560051b01019161029f565b60831988860301815292955090935091600190602090829061071590886124b6565b96019201920192859387959361025f565b635e1e452d60e11b8152600490fd5b5034610147576040366003190112610147576004356001600160401b0381116101d45761016060031982360301126101d4576040519161016083018381106001600160401b038211176104ec5760405281600401356001600160401b03811681036101d457835260248201356020840152604482013560408401526064820135606084015260848201356001600160401b0381116101d457820191366023840112156101d4576004830135926107ea8461256c565b936107f86040519586612537565b80855260051b810160240160208501368211610cf65760248301905b828210610cc157505050506080840192835260a48101356001600160401b038111610c1c57810136602382011215610c1c576004810135906108558261256c565b916108636040519384612537565b80835260051b810160240160208301368211610c8c5760248301905b828210610c90575050505060a085015260c48101356001600160401b038111610c1c57810136602382011215610c1c5760048101356108bd8161256c565b916108cb6040519384612537565b818352602060048185019360061b8301010190368211610c8c57602401915b818310610c205750505060c085015260e481013560e085015261010481013561010085015261012481013560ff81168103610c1c576101208501526101448101356001600160401b038111610c1c57610948913691016004016125d4565b6101408401526024356001600160401b0381116101d45761096d9036906004016125d4565b6002546001600160a01b03169290919083158015610c01575b610bf257516040519081604081019160208083015280518093526060820192602060608260051b85010192019386905b828210610bc5575050506109d3925003601f198101835282612537565b602081519101209360a08101516040519081604081019160208083015280518093526060820192602060608260051b85010192019387905b828210610b9857505050610a28925003601f198101835282612537565b6020815191012060c082015195604051604081016020808301528851809152602060608301990190865b818110610b6d576020610b5a8b610b518c8f8d8d8d610b228e610b148f610a8481603c9903601f198101835282612537565b8c81519101206001600160401b03865116958d8101519360408201519660608301519360e0840151916101008501519361014060ff61012088015116960151966040519b60209a8d9b8c019e8f5260408c015260608b015260808a015260a089015260c088015260e0870152610100860152610120850152610140840152610160808401526101808301906122e9565b03601f198101835282612537565b5190207f19457468657265756d205369676e6564204d6573736167653a0a3332000000008252601c5220612bb7565b90939193612bf1565b6040516001600160a01b03909216148152f35b825180516001600160a01b03168c52602090810151818d01526040909b019a90920191600101610a52565b91935091602080610bb5600193605f198a820301865288516122e9565b9601920192018593919492610a0b565b91935091602080610be2600193605f198a820301865288516122e9565b96019201920185939194926109b6565b63f64b1d7b60e01b8252600482fd5b506085604051610c148161012881612394565b511415610986565b8280fd5b604083360312610c8c57604051604081018181106001600160401b03821117610c785760405283356001600160a01b0381168103610c745791816040936020935282860135838201528152019201916108ea565b8780fd5b634e487b7160e01b88526041600452602488fd5b8580fd5b81356001600160401b038111610c7457602091610cb6839283600436928a0101016125d4565b81520191019061087f565b81356001600160401b038111610cf257602091610ce7839283600436928a0101016125d4565b815201910190610814565b8680fd5b8480fd5b5034610147578060031936011261014757602060405160858152f35b50346101475780600319360112610147576020600e54604051908152f35b5034610147578060031936011261014757546040516001600160a01b039091168152602090f35b503461014757610d6a3661230d565b610d75929192612841565b610d7d612b59565b610d8836828561259e565b60208151910120806008558252600460205260ff6040832054166118b457604051636a107b3560e11b815292829184918291610dea917f0000000000000000000000000000000000000000000000000000000000000000919060048501612795565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156118a7578190829383918491859186916117b8575b508051906001600160401b03821161125657610e4d600b5461235c565b601f811161176a575b50602090601f83116001146116e857610e8592918891836104095750508160011b915f199060031b1c19161790565b600b555b8051906001600160401b0382116116d457610ea5600a5461235c565b601f8111611699575b50602090601f831160011461162a57610edd92918791836104095750508160011b915f199060031b1c19161790565b600a555b8051906001600160401b03821161161657610efd60095461235c565b601f81116115c8575b50602090601f831160011461154657610f3592918691836104095750508160011b915f199060031b1c19161790565b6009555b805190600160401b821161153257600e5482600e558083106114a6575b50602001600e84525f516020612ce85f395f51905f5284915b8383106113d8575050505082516001600160401b0381116112f657610f95600d5461235c565b601f811161138a575b506020601f821160011461130a5781908495610fce9495926104095750508160011b915f199060031b1c19161790565b600d555b805190600160401b82116112f657600f5482600f5580831061126a575b50602001600f83525f516020612d685f395f51905f5283915b83831061118257604051600b54869082826110228361235c565b80835292600181169081156111635750600114611105575b61104692500383612537565b6040516110568161012881612433565b604051908293600a546110688161235c565b80855290600181169081156110e1575060011461109e575b506110918361109695960384612537565b612a2f565b600160055580f35b600a85529450835f516020612d085f395f51905f525b8682106110cb575083016020019450611091611080565b60018160209254838589010152019101906110b4565b60ff191660208087019190915291151560051b850190910195506110919050611080565b50600b83529082907f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db95b8183106111475750509060206110469282010161103a565b602091935080600191548385890101520191019091849261112f565b6020925061104694915060ff191682840152151560051b82010161103a565b80518051906001600160401b038211611256576111a9826111a3865461235c565b866126f8565b602090601f83116001146111ee57926111df836001959460209487968c926104095750508160011b915f199060031b1c19161790565b85555b01920192019190611008565b8488528188209190601f198416895b81811061123e5750936020936001969387969383889510611226575b505050811b0185556111e2565b01515f1960f88460031b161c191690555f8080611219565b929360206001819287860151815501950193016111fd565b634e487b7160e01b87526041600452602487fd5b600f84525f516020612d685f395f51905f5201825f516020612d685f395f51905f52015b81811061129b5750610fef565b806112a86001925461235c565b806112b5575b500161128e565b601f811183146112ca57508581555b5f6112ae565b818752602087206112e591601f0160051c81019084016126a8565b8086528560208120818355556112c4565b634e487b7160e01b83526041600452602483fd5b600d84527fd7b6990105719101dabeb77144f2a3385c8033acd3af97e9423a695e81ad1eb590601f198316855b8181106113725750958360019596971061135a575b505050811b01600d55610fd2565b01515f1960f88460031b161c191690555f808061134c565b9192602060018192868b015181550194019201611337565b600d84526113d2907fd7b6990105719101dabeb77144f2a3385c8033acd3af97e9423a695e81ad1eb5601f840160051c810191602085106104e257601f0160051c01906126a8565b5f610f9e565b80518051906001600160401b038211610c78576113f9826111a3865461235c565b602090601f831160011461143e579261142f836001959460209487968d926104095750508160011b915f199060031b1c19161790565b85555b01920192019190610f6f565b8489528189209190601f1984168a5b81811061148e5750936020936001969387969383889510611476575b505050811b018555611432565b01515f1960f88460031b161c191690555f8080611469565b9293602060018192878601518155019501930161144d565b600e85525f516020612ce85f395f51905f5201825f516020612ce85f395f51905f52015b8181106114d75750610f56565b806114e46001925461235c565b806114f1575b50016114ca565b601f8111831461150657508681555b5f6114ea565b8188526020882061152191601f0160051c81019084016126a8565b808752866020812081835555611500565b634e487b7160e01b84526041600452602484fd5b600986527f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af9190601f198416875b8181106115b05750908460019594939210611598575b505050811b01600955610f39565b01515f1960f88460031b161c191690555f808061158a565b92936020600181928786015181550195019301611574565b60098652611610907f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af601f850160051c810191602086106104e257601f0160051c01906126a8565b5f610f06565b634e487b7160e01b85526041600452602485fd5b600a87525f516020612d085f395f51905f529190601f198416885b8181106116815750908460019594939210611669575b505050811b01600a55610ee1565b01515f1960f88460031b161c191690555f808061165b565b92936020600181928786015181550195019301611645565b600a87526116ce905f516020612d085f395f51905f52601f850160051c810191602086106104e257601f0160051c01906126a8565b5f610eae565b634e487b7160e01b86526041600452602486fd5b600b88527f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db99190601f198416895b818110611752575090846001959493921061173a575b505050811b01600b55610e89565b01515f1960f88460031b161c191690555f808061172c565b92936020600181928786015181550195019301611716565b600b88526117b2907f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db9601f850160051c810191602086106104e257601f0160051c01906126a8565b5f610e56565b965050505050503d8082843e6117ce8184612537565b82019160c0818403126101d45780516001600160401b038111610c1c57836117f79183016127c3565b9260208201516001600160401b0381116118a35781611817918401612662565b9160408101516001600160401b038111610cf657826118379183016127c3565b9460608201516001600160401b038111610c8c5783611857918401612662565b9260808301516001600160401b038111610cf25781611877918501612662565b9260a0810151906001600160401b038211610c7457611897929101612662565b9093959291905f610e30565b8380fd5b50604051903d90823e3d90fd5b6305fb89e760e21b8252600482fd5b50346101475780600319360112610147576040519080600154906118e68261235c565b8085529160018116908115611958575060011461190e575b6101438461012f81860382612537565b600181525f516020612d485f395f51905f52939250905b80821061193e5750909150810160200161012f826118fe565b919260018160209254838588010152019101909291611925565b60ff191660208087019190915292151560051b8501909201925061012f91508390506118fe565b50346101475761198e3661230d565b611999929192612841565b6040516001549080846119ab8461235c565b600185168015611b3657600114611af8575b500390206040518386823780848101868152039020907f7cb84e5a76f21dbf9cffe637ecb554a3d4f4eeb2cf05cd38c8a0db9d34ed65418580a36001600160401b0382116112f657611a0e9061235c565b601f8111611abd575b5081601f8211600114611a565781908394611a459492611a4b5750508160011b915f199060031b1c19161790565b60015580f35b013590505f8061033a565b601f198216935f516020612d485f395f51905f5291845b868110611aa55750836001959610611a8c575b505050811b0160015580f35b01355f19600384901b60f8161c191690555f8080611a80565b90926020600181928686013581550194019101611a6d565b60018352611af2905f516020612d485f395f51905f52601f840160051c810191602085106104e257601f0160051c01906126a8565b5f611a17565b600187529050855f516020612d485f395f51905f525b828210611b1f57505081015f6119bd565b805482860152849350602090910190600101611b0e565b5060ff198516835280151502820190505f6119bd565b503461014757611b5b3661230d565b90611b64612841565b611b6f36838361259e565b6020815191012091828452600460205260ff604085205416611c9e57604051632b5f2f8160e01b81529291849184918291611bd0917f0000000000000000000000000000000000000000000000000000000000000000919060048501612795565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa8015611c9357611c4b92849185918691611c6e575b508183611c1f92612a2f565b60208151910151906bffffffffffffffffffffffff1982169160148210611c4e575b505060601c612867565b80f35b6001600160601b031960149290920360031b82901b161690505f80611c41565b9050611c1f9250611c8991503d8087833e6106548183612537565b9192509081611c13565b6040513d85823e3d90fd5b6305fb89e760e21b8452600484fd5b5034610147578060031936011261014757611cc6612841565b80546001600160a01b03198116825581906001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b5034610147578060031936011261014757611d20612841565b6001600554036107265760065460405163c13734e760e01b815260606004820152600e805460648301819052908452909291829184916084600583901b84018101925f516020612ce85f395f51905f5291869086015b828210611ef4575050505060248301528181036003190160448301528190611d9d906125f2565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156118a7578192611eb5575b5081516001600160401b0381116104ec57611df98161037160075461235c565b602092601f8211600114611e495782938291611e2894926104095750508160011b915f199060031b1c19161790565b6007555b60016006540180600655600e5414611e415780f35b600260055580f35b60078352601f198216935f516020612d285f395f51905f5291845b868110611e9d5750836001959610611e85575b505050811b01600755611e2c565b01515f1960f88460031b161c191690555f8080611e77565b91926020600181928685015181550194019201611e64565b9091503d8083833e611ec78183612537565b8101602082820312610c1c5781516001600160401b0381116118a357611eed9201612662565b905f611dd9565b6083198a8703018152929650929450926001906020908290611f1690896124b6565b97019201920192869593889593611d76565b50346101475780600319360112610147576002546040516001600160a01b039091168152602090f35b50346121f2575f3660031901126121f257611f6a612841565b6003600554036121f6577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156121f25760405190632d77dd9360e11b825260606004830152815f600c54611fc88161235c565b908160648501526001811690815f146121cd5750600114612173575b50818103600319016024830152611ffa906125f2565b8181036003190160448301526010545f916120148261235c565b808252916001811690811561214f57506001146120ef575b50509181805f9403915afa80156120e4576120d1575b50611c4b612051600a5461235c565b600a601f82116120b9575b546001600160601b03198116919060148210612099575b505060085490604051906120918261208a81612433565b0383612537565b60601c612867565b6001600160601b031960149290920360031b82901b161690505f80612073565b50600a83525f516020612d085f395f51905f5261205c565b6120dd91505f90612537565b5f5f612042565b6040513d5f823e3d90fd5b60105f9081527f1b6847dc741a1b0cd08d278845f9d819d87b734759afb55fe2de5cb82a9ae672959350915b808310612131575091935001602001818061202c565b9350919360018160209254838587010152019101909391859361211b565b60ff191660208381019190915292151560051b90910190910191508290508061202c565b600c5f90815291507fdf6966c971051c3d54ec59162606531493a51404a002842f56009d7e5cf4a8c75b8183106121b35750508101608401611ffa611fe4565b80548387016084015285935060209092019160010161219d565b611ffa9350606491509160209260ff19166084860152151560051b8401010190611fe4565b5f80fd5b635e1e452d60e11b5f5260045ffd5b346121f2575f3660031901126121f25760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346121f2575f3660031901126121f2576002546040516001600160a01b039091168152602090f35b346121f2575f3660031901126121f2576020600554604051908152f35b346121f2575f3660031901126121f2576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346121f2575f3660031901126121f2578061012f8161012861014394612394565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9060206003198301126121f2576004356001600160401b0381116121f257826023820112156121f2578060040135926001600160401b0384116121f257602484830101116121f2576024019190565b90600182811c9216801561238a575b602083101461237657565b634e487b7160e01b5f52602260045260245ffd5b91607f169161236b565b6003545f92916123a38261235c565b808252916001811690811561241757506001146123be575050565b60035f9081529293509091907fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b5b8383106123fd575060209250010190565b6001816020929493945483858701015201910191906123ec565b9050602093945060ff929192191683830152151560051b010190565b6009545f92916124428261235c565b8082529160018116908115612417575060011461245d575050565b60095f9081529293509091907f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af5b83831061249c575060209250010190565b60018160209294939454838587010152019101919061248b565b5f92918154916124c58361235c565b808352926001811690811561251a57506001146124e157505050565b5f9081526020812093945091925b838310612500575060209250010190565b6001816020929493945483858701015201910191906124ef565b915050602093945060ff929192191683830152151560051b010190565b90601f801991011681019081106001600160401b0382111761255857604052565b634e487b7160e01b5f52604160045260245ffd5b6001600160401b0381116125585760051b60200190565b6001600160401b03811161255857601f01601f191660200190565b9291926125aa82612583565b916125b86040519384612537565b8294818452818301116121f2578281602093845f960137010152565b9080601f830112156121f2578160206125ef9335910161259e565b90565b6007545f92916126018261235c565b8082529160018116908115612417575060011461261c575050565b60075f9081529293509091905f516020612d285f395f51905f525b838310612648575060209250010190565b600181602092949394548385870101520191019190612637565b81601f820112156121f25780519061267982612583565b926126876040519485612537565b828452602083830101116121f257815f9260208093018386015e8301015290565b8181106126b3575050565b5f81556001016126a8565b90601f82116126cb575050565b6126f69160075f5260205f20906020601f840160051c830193106104e257601f0160051c01906126a8565b565b9190601f811161270757505050565b6126f6925f5260205f20906020601f840160051c830193106104e257601f0160051c01906126a8565b916060838303126121f25782516001600160401b0381116121f25782612757918501612662565b9260208101516001600160401b0381116121f25783612777918301612662565b9260408201516001600160401b0381116121f2576125ef9201612662565b9392918060609160209360408852816040890152838801375f828288010152601f8019910116850101930152565b9080601f830112156121f25781516127da8161256c565b926127e86040519485612537565b81845260208085019260051b820101918383116121f25760208201905b83821061281457505050505090565b81516001600160401b0381116121f25760209161283687848094880101612662565b815201910190612805565b5f546001600160a01b0316330361285457565b63118cdaa760e01b5f523360045260245ffd5b9092917f1882fd6e997ef1dc5e2691efe78564a5c540b36547363e8a7dcd454c61a944b660018060a01b03600254169260405193845260018060a01b03169283602082015260806040820152806128d16128c360808301612394565b8281036060840152886122e9565b0390a15f52600460205260405f20600160ff198254161790556bffffffffffffffffffffffff60a01b600254161760025581516001600160401b0381116125585761291d60035461235c565b601f81116129e1575b50602092601f821160011461296257612956929382915f926104095750508160011b915f199060031b1c19161790565b6003555b6126f6612b59565b601f1982169360035f527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b915f5b8681106129c957508360019596106129b1575b505050811b0160035561295a565b01515f1960f88460031b161c191690555f80806129a3565b91926020600181928685015181550194019201612990565b60035f52612a29907fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b601f840160051c810191602085106104e257601f0160051c01906126a8565b5f612926565b915160131901612b4a575160841901612b3b57805160015490612a518261235c565b6004019081600411612b275710612ae1575f91612a6d8261235c565b925b838103612a7c5750505050565b60048101808211612b275782511115612b1357818101602401516001600160f81b031916612aa98461235c565b80831015612b135760201115612af057601f8216601f036001905b60ff60f81b91549060031b1c60f81b1603612ae157600101612a6f565b63719a9a4960e01b5f5260045ffd5b60015f528160051c5f516020612d485f395f51905f5201601f8316601f03612ac4565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b634ae601c160e11b5f5260045ffd5b6342f6e2a360e11b5f5260045ffd5b5f6005555f6006555f604051612b70602082612537565b52612b7c60075461235c565b601f8111612b8c575b505f600755565b60075f52612bb190601f0160051c5f516020612d285f395f51905f52908101906126a8565b5f612b85565b8151919060418303612be757612be09250602082015190606060408401519301515f1a90612c65565b9192909190565b50505f9160029190565b6004811015612c515780612c03575050565b60018103612c1a5763f645eedf60e01b5f5260045ffd5b60028103612c35575063fce698f760e01b5f5260045260245ffd5b600314612c3f5750565b6335e2f38360e21b5f5260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411612cdc579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa156120e4575f516001600160a01b03811615612cd257905f905f90565b505f906001905f90565b5050505f916003919056febb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fdc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8a66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688b10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf68d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac802a2646970667358221220f3cdebb0fa019410e787d5f621416c5ea26d012d190458b50fb2fd1f4b54e2ca64736f6c634300081e0033",
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
// Solidity: constructor(address owner, address _nitroProver, bytes _pcr0, uint256 _maxVerificationAge) returns()
func (teeAuthenticator *TeeAuthenticator) PackConstructor(owner common.Address, _nitroProver common.Address, _pcr0 []byte, _maxVerificationAge *big.Int) []byte {
	enc, err := teeAuthenticator.abi.Pack("", owner, _nitroProver, _pcr0, _maxVerificationAge)
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
func (teeAuthenticator *TeeAuthenticator) PackPKLENGTH() []byte {
	enc, err := teeAuthenticator.abi.Pack("PK_LENGTH")
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
func (teeAuthenticator *TeeAuthenticator) TryPackPKLENGTH() ([]byte, error) {
	return teeAuthenticator.abi.Pack("PK_LENGTH")
}

// UnpackPKLENGTH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc91496c6.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackPKLENGTH(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("PK_LENGTH", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce6dce51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce6dce51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce6dce51.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := teeAuthenticator.abi.Unpack("checkSignature", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackCurrentUpdateStep is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0b1bfab9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function currentUpdateStep() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackCurrentUpdateStep() []byte {
	enc, err := teeAuthenticator.abi.Pack("currentUpdateStep")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCurrentUpdateStep is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0b1bfab9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function currentUpdateStep() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackCurrentUpdateStep() ([]byte, error) {
	return teeAuthenticator.abi.Pack("currentUpdateStep")
}

// UnpackCurrentUpdateStep is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0b1bfab9.
//
// Solidity: function currentUpdateStep() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackCurrentUpdateStep(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("currentUpdateStep", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPubSecp521r1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081bec7e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (teeAuthenticator *TeeAuthenticator) PackGetPubSecp521r1() []byte {
	enc, err := teeAuthenticator.abi.Pack("getPubSecp521r1")
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
func (teeAuthenticator *TeeAuthenticator) TryPackGetPubSecp521r1() ([]byte, error) {
	return teeAuthenticator.abi.Pack("getPubSecp521r1")
}

// UnpackGetPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081bec7e.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (teeAuthenticator *TeeAuthenticator) UnpackGetPubSecp521r1(data []byte) ([]byte, error) {
	out, err := teeAuthenticator.abi.Unpack("getPubSecp521r1", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetStep2TotalLength is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3b33122.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getStep2TotalLength() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackGetStep2TotalLength() []byte {
	enc, err := teeAuthenticator.abi.Pack("getStep2TotalLength")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetStep2TotalLength is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3b33122.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getStep2TotalLength() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackGetStep2TotalLength() ([]byte, error) {
	return teeAuthenticator.abi.Pack("getStep2TotalLength")
}

// UnpackGetStep2TotalLength is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3b33122.
//
// Solidity: function getStep2TotalLength() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackGetStep2TotalLength(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("getStep2TotalLength", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetTeeSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dd7ce2f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getTeeSigner() view returns(address)
func (teeAuthenticator *TeeAuthenticator) PackGetTeeSigner() []byte {
	enc, err := teeAuthenticator.abi.Pack("getTeeSigner")
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
func (teeAuthenticator *TeeAuthenticator) TryPackGetTeeSigner() ([]byte, error) {
	return teeAuthenticator.abi.Pack("getTeeSigner")
}

// UnpackGetTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dd7ce2f.
//
// Solidity: function getTeeSigner() view returns(address)
func (teeAuthenticator *TeeAuthenticator) UnpackGetTeeSigner(data []byte) (common.Address, error) {
	out, err := teeAuthenticator.abi.Unpack("getTeeSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackMaxVerificationAge is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x127379db.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxVerificationAge() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackMaxVerificationAge() []byte {
	enc, err := teeAuthenticator.abi.Pack("maxVerificationAge")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxVerificationAge is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x127379db.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxVerificationAge() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackMaxVerificationAge() ([]byte, error) {
	return teeAuthenticator.abi.Pack("maxVerificationAge")
}

// UnpackMaxVerificationAge is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x127379db.
//
// Solidity: function maxVerificationAge() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackMaxVerificationAge(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("maxVerificationAge", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackNitroProver is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a83fb91.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function nitroProver() view returns(address)
func (teeAuthenticator *TeeAuthenticator) PackNitroProver() []byte {
	enc, err := teeAuthenticator.abi.Pack("nitroProver")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNitroProver is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a83fb91.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function nitroProver() view returns(address)
func (teeAuthenticator *TeeAuthenticator) TryPackNitroProver() ([]byte, error) {
	return teeAuthenticator.abi.Pack("nitroProver")
}

// UnpackNitroProver is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a83fb91.
//
// Solidity: function nitroProver() view returns(address)
func (teeAuthenticator *TeeAuthenticator) UnpackNitroProver(data []byte) (common.Address, error) {
	out, err := teeAuthenticator.abi.Unpack("nitroProver", data)
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

// PackPcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x81a9d38a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pcr0() view returns(bytes)
func (teeAuthenticator *TeeAuthenticator) PackPcr0() []byte {
	enc, err := teeAuthenticator.abi.Pack("pcr0")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x81a9d38a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pcr0() view returns(bytes)
func (teeAuthenticator *TeeAuthenticator) TryPackPcr0() ([]byte, error) {
	return teeAuthenticator.abi.Pack("pcr0")
}

// UnpackPcr0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x81a9d38a.
//
// Solidity: function pcr0() view returns(bytes)
func (teeAuthenticator *TeeAuthenticator) UnpackPcr0(data []byte) ([]byte, error) {
	out, err := teeAuthenticator.abi.Unpack("pcr0", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackPubSecp521r1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe4993ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pubSecp521r1() view returns(bytes)
func (teeAuthenticator *TeeAuthenticator) PackPubSecp521r1() []byte {
	enc, err := teeAuthenticator.abi.Pack("pubSecp521r1")
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
func (teeAuthenticator *TeeAuthenticator) TryPackPubSecp521r1() ([]byte, error) {
	return teeAuthenticator.abi.Pack("pubSecp521r1")
}

// UnpackPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfe4993ca.
//
// Solidity: function pubSecp521r1() view returns(bytes)
func (teeAuthenticator *TeeAuthenticator) UnpackPubSecp521r1(data []byte) ([]byte, error) {
	out, err := teeAuthenticator.abi.Unpack("pubSecp521r1", data)
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

// PackStep2CurrentIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdb6a7b71.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function step2CurrentIndex() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackStep2CurrentIndex() []byte {
	enc, err := teeAuthenticator.abi.Pack("step2CurrentIndex")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStep2CurrentIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdb6a7b71.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function step2CurrentIndex() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackStep2CurrentIndex() ([]byte, error) {
	return teeAuthenticator.abi.Pack("step2CurrentIndex")
}

// UnpackStep2CurrentIndex is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdb6a7b71.
//
// Solidity: function step2CurrentIndex() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackStep2CurrentIndex(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("step2CurrentIndex", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
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

// PackUpdatePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ed2d7fc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updatePcr0(bytes newPcr0) returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdatePcr0(newPcr0 []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("updatePcr0", newPcr0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdatePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ed2d7fc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updatePcr0(bytes newPcr0) returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdatePcr0(newPcr0 []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("updatePcr0", newPcr0)
}

// PackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7637d58a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTee(bytes attestation) returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdateTee(attestation []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("updateTee", attestation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7637d58a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTee(bytes attestation) returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdateTee(attestation []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("updateTee", attestation)
}

// PackUpdateTeeStep1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8890b84e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep1(bytes attestation) returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdateTeeStep1(attestation []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("updateTeeStep1", attestation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8890b84e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep1(bytes attestation) returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdateTeeStep1(attestation []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("updateTeeStep1", attestation)
}

// PackUpdateTeeStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5520916c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep2() returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdateTeeStep2() []byte {
	enc, err := teeAuthenticator.abi.Pack("updateTeeStep2")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5520916c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep2() returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdateTeeStep2() ([]byte, error) {
	return teeAuthenticator.abi.Pack("updateTeeStep2")
}

// PackUpdateTeeStep3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd22a29d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep3() returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdateTeeStep3() []byte {
	enc, err := teeAuthenticator.abi.Pack("updateTeeStep3")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd22a29d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep3() returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdateTeeStep3() ([]byte, error) {
	return teeAuthenticator.abi.Pack("updateTeeStep3")
}

// PackUpdateTeeStep4 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34010ccf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeStep4() returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdateTeeStep4() []byte {
	enc, err := teeAuthenticator.abi.Pack("updateTeeStep4")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeStep4 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34010ccf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeStep4() returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdateTeeStep4() ([]byte, error) {
	return teeAuthenticator.abi.Pack("updateTeeStep4")
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

// TeeAuthenticatorPcrZeroUpdate represents a PcrZeroUpdate event raised by the TeeAuthenticator contract.
type TeeAuthenticatorPcrZeroUpdate struct {
	OldPcr0 common.Hash
	NewPcr0 common.Hash
	Raw     *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorPcrZeroUpdateEventName = "PcrZeroUpdate"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorPcrZeroUpdate) ContractEventName() string {
	return TeeAuthenticatorPcrZeroUpdateEventName
}

// UnpackPcrZeroUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PcrZeroUpdate(bytes indexed oldPcr0, bytes indexed newPcr0)
func (teeAuthenticator *TeeAuthenticator) UnpackPcrZeroUpdateEvent(log *types.Log) (*TeeAuthenticatorPcrZeroUpdate, error) {
	event := "PcrZeroUpdate"
	if log.Topics[0] != teeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TeeAuthenticatorPcrZeroUpdate)
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
	OldTee          common.Address
	NewTee          common.Address
	OldPubSecp521r1 []byte
	NewPubSecp521r1 []byte
	Raw             *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorTeeUpdateEventName = "TeeUpdate"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorTeeUpdate) ContractEventName() string {
	return TeeAuthenticatorTeeUpdateEventName
}

// UnpackTeeUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1)
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
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["AttestationAlreadyUsed"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackAttestationAlreadyUsedError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["InvalidPCR"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackInvalidPCRError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["InvalidPKLength"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackInvalidPKLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["InvalidUserDataLength"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackInvalidUserDataLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["WrongStep"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackWrongStepError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TeeAuthenticatorAttestationAlreadyUsed represents a AttestationAlreadyUsed error raised by the TeeAuthenticator contract.
type TeeAuthenticatorAttestationAlreadyUsed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AttestationAlreadyUsed()
func TeeAuthenticatorAttestationAlreadyUsedErrorID() common.Hash {
	return common.HexToHash("0x17ee279c4380943995e5de7b6ca148f56c43f727fc1a52fa3f267e6f5ea97559")
}

// UnpackAttestationAlreadyUsedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AttestationAlreadyUsed()
func (teeAuthenticator *TeeAuthenticator) UnpackAttestationAlreadyUsedError(raw []byte) (*TeeAuthenticatorAttestationAlreadyUsed, error) {
	out := new(TeeAuthenticatorAttestationAlreadyUsed)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "AttestationAlreadyUsed", raw); err != nil {
		return nil, err
	}
	return out, nil
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

// TeeAuthenticatorInvalidPCR represents a InvalidPCR error raised by the TeeAuthenticator contract.
type TeeAuthenticatorInvalidPCR struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPCR()
func TeeAuthenticatorInvalidPCRErrorID() common.Hash {
	return common.HexToHash("0x719a9a4917a12c04cb8a46a64960d788ba098ccc7ac543fcee9fa361a6d1224c")
}

// UnpackInvalidPCRError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPCR()
func (teeAuthenticator *TeeAuthenticator) UnpackInvalidPCRError(raw []byte) (*TeeAuthenticatorInvalidPCR, error) {
	out := new(TeeAuthenticatorInvalidPCR)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "InvalidPCR", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorInvalidPKLength represents a InvalidPKLength error raised by the TeeAuthenticator contract.
type TeeAuthenticatorInvalidPKLength struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPKLength()
func TeeAuthenticatorInvalidPKLengthErrorID() common.Hash {
	return common.HexToHash("0x95cc0382d99b2fba4335843859b8b7e3894ae70a43b516ae2d581373538dd6ac")
}

// UnpackInvalidPKLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPKLength()
func (teeAuthenticator *TeeAuthenticator) UnpackInvalidPKLengthError(raw []byte) (*TeeAuthenticatorInvalidPKLength, error) {
	out := new(TeeAuthenticatorInvalidPKLength)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "InvalidPKLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorInvalidUserDataLength represents a InvalidUserDataLength error raised by the TeeAuthenticator contract.
type TeeAuthenticatorInvalidUserDataLength struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidUserDataLength()
func TeeAuthenticatorInvalidUserDataLengthErrorID() common.Hash {
	return common.HexToHash("0x85edc5464c5565960356d11f255a83de2f3ba3c24b8485a39f37b6c98eec54bc")
}

// UnpackInvalidUserDataLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidUserDataLength()
func (teeAuthenticator *TeeAuthenticator) UnpackInvalidUserDataLengthError(raw []byte) (*TeeAuthenticatorInvalidUserDataLength, error) {
	out := new(TeeAuthenticatorInvalidUserDataLength)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "InvalidUserDataLength", raw); err != nil {
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

// TeeAuthenticatorWrongStep represents a WrongStep error raised by the TeeAuthenticator contract.
type TeeAuthenticatorWrongStep struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongStep()
func TeeAuthenticatorWrongStepErrorID() common.Hash {
	return common.HexToHash("0xbc3c8a5a94b27f160f1b92f62d402d946766caac5227a2ffc2d1381f57f4e236")
}

// UnpackWrongStepError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongStep()
func (teeAuthenticator *TeeAuthenticator) UnpackWrongStepError(raw []byte) (*TeeAuthenticatorWrongStep, error) {
	out := new(TeeAuthenticatorWrongStep)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "WrongStep", raw); err != nil {
		return nil, err
	}
	return out, nil
}
