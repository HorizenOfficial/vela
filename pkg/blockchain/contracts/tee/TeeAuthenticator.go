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

// StructsEventData is an auto generated low-level Go binding around an user-defined struct.
type StructsEventData struct {
	Events   [][]byte
	SubTypes [][32]byte
}

// StructsSignatureParams is an auto generated low-level Go binding around an user-defined struct.
type StructsSignatureParams struct {
	ApplicationId      uint64
	PrevStateRoot      [32]byte
	NewStateRoot       [32]byte
	ProcessedRequestId [32]byte
	UserEvents         StructsEventData
	AppEvents          StructsEventData
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
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0x5c9626b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := abstractTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c9626b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return abstractTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c9626b9.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
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
	ABI: "[{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0x5c9626b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c9626b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c9626b9.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
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
	ABI: "[{\"inputs\":[],\"name\":\"AttestationAlreadyUsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotRemoveActiveImage\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPCR\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPcr0Length\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUserDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoPendingSwap\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SwapAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SwapProposalExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimelockNotElapsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnknownPcr0\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStep\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0Removed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0SwapCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"}],\"name\":\"Pcr0SwapProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0Swapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"applyPcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelPcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAcceptedPcr0Count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStep2TotalLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"}],\"name\":\"proposePcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"oldPcr0\",\"type\":\"bytes\"}],\"name\":\"removePcr0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTeeStep1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// PackApplyPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b3917cb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function applyPcr0Swap() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackApplyPcr0Swap() []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("applyPcr0Swap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApplyPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b3917cb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function applyPcr0Swap() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackApplyPcr0Swap() ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("applyPcr0Swap")
}

// PackCancelPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35732f4e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelPcr0Swap() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackCancelPcr0Swap() []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("cancelPcr0Swap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35732f4e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelPcr0Swap() returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackCancelPcr0Swap() ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("cancelPcr0Swap")
}

// PackGetAcceptedPcr0Count is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c313844.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getAcceptedPcr0Count() view returns(uint256 count)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackGetAcceptedPcr0Count() []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("getAcceptedPcr0Count")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetAcceptedPcr0Count is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c313844.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getAcceptedPcr0Count() view returns(uint256 count)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackGetAcceptedPcr0Count() ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("getAcceptedPcr0Count")
}

// UnpackGetAcceptedPcr0Count is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2c313844.
//
// Solidity: function getAcceptedPcr0Count() view returns(uint256 count)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackGetAcceptedPcr0Count(data []byte) (*big.Int, error) {
	out, err := iTeeAuthenticatorAdmin.abi.Unpack("getAcceptedPcr0Count", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
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

// PackProposePcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07cdb0c0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proposePcr0Swap(bytes targetPcr0) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackProposePcr0Swap(targetPcr0 []byte) []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("proposePcr0Swap", targetPcr0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProposePcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07cdb0c0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proposePcr0Swap(bytes targetPcr0) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackProposePcr0Swap(targetPcr0 []byte) ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("proposePcr0Swap", targetPcr0)
}

// PackRemovePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa239f11d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function removePcr0(bytes oldPcr0) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) PackRemovePcr0(oldPcr0 []byte) []byte {
	enc, err := iTeeAuthenticatorAdmin.abi.Pack("removePcr0", oldPcr0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRemovePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa239f11d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function removePcr0(bytes oldPcr0) returns()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) TryPackRemovePcr0(oldPcr0 []byte) ([]byte, error) {
	return iTeeAuthenticatorAdmin.abi.Pack("removePcr0", oldPcr0)
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

// ITeeAuthenticatorAdminPcr0Removed represents a Pcr0Removed event raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminPcr0Removed struct {
	Pcr0 []byte
	Raw  *types.Log // Blockchain specific contextual infos
}

const ITeeAuthenticatorAdminPcr0RemovedEventName = "Pcr0Removed"

// ContractEventName returns the user-defined event name.
func (ITeeAuthenticatorAdminPcr0Removed) ContractEventName() string {
	return ITeeAuthenticatorAdminPcr0RemovedEventName
}

// UnpackPcr0RemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0Removed(bytes pcr0)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackPcr0RemovedEvent(log *types.Log) (*ITeeAuthenticatorAdminPcr0Removed, error) {
	event := "Pcr0Removed"
	if log.Topics[0] != iTeeAuthenticatorAdmin.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ITeeAuthenticatorAdminPcr0Removed)
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

// ITeeAuthenticatorAdminPcr0SwapCancelled represents a Pcr0SwapCancelled event raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminPcr0SwapCancelled struct {
	TargetPcr0 []byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const ITeeAuthenticatorAdminPcr0SwapCancelledEventName = "Pcr0SwapCancelled"

// ContractEventName returns the user-defined event name.
func (ITeeAuthenticatorAdminPcr0SwapCancelled) ContractEventName() string {
	return ITeeAuthenticatorAdminPcr0SwapCancelledEventName
}

// UnpackPcr0SwapCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0SwapCancelled(bytes targetPcr0)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackPcr0SwapCancelledEvent(log *types.Log) (*ITeeAuthenticatorAdminPcr0SwapCancelled, error) {
	event := "Pcr0SwapCancelled"
	if log.Topics[0] != iTeeAuthenticatorAdmin.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ITeeAuthenticatorAdminPcr0SwapCancelled)
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

// ITeeAuthenticatorAdminPcr0SwapProposed represents a Pcr0SwapProposed event raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminPcr0SwapProposed struct {
	TargetPcr0 []byte
	Eta        *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const ITeeAuthenticatorAdminPcr0SwapProposedEventName = "Pcr0SwapProposed"

// ContractEventName returns the user-defined event name.
func (ITeeAuthenticatorAdminPcr0SwapProposed) ContractEventName() string {
	return ITeeAuthenticatorAdminPcr0SwapProposedEventName
}

// UnpackPcr0SwapProposedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0SwapProposed(bytes targetPcr0, uint256 eta)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackPcr0SwapProposedEvent(log *types.Log) (*ITeeAuthenticatorAdminPcr0SwapProposed, error) {
	event := "Pcr0SwapProposed"
	if log.Topics[0] != iTeeAuthenticatorAdmin.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ITeeAuthenticatorAdminPcr0SwapProposed)
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

// ITeeAuthenticatorAdminPcr0Swapped represents a Pcr0Swapped event raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminPcr0Swapped struct {
	Pcr0 []byte
	Raw  *types.Log // Blockchain specific contextual infos
}

const ITeeAuthenticatorAdminPcr0SwappedEventName = "Pcr0Swapped"

// ContractEventName returns the user-defined event name.
func (ITeeAuthenticatorAdminPcr0Swapped) ContractEventName() string {
	return ITeeAuthenticatorAdminPcr0SwappedEventName
}

// UnpackPcr0SwappedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0Swapped(bytes pcr0)
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackPcr0SwappedEvent(log *types.Log) (*ITeeAuthenticatorAdminPcr0Swapped, error) {
	event := "Pcr0Swapped"
	if log.Topics[0] != iTeeAuthenticatorAdmin.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ITeeAuthenticatorAdminPcr0Swapped)
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
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["CannotRemoveActiveImage"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackCannotRemoveActiveImageError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["InvalidPCR"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackInvalidPCRError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["InvalidPKLength"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackInvalidPKLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["InvalidPcr0Length"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackInvalidPcr0LengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["InvalidUserDataLength"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackInvalidUserDataLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["NoPendingSwap"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackNoPendingSwapError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["SwapAlreadyPending"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackSwapAlreadyPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["SwapProposalExpired"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackSwapProposalExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["TimelockNotElapsed"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackTimelockNotElapsedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticatorAdmin.abi.Errors["UnknownPcr0"].ID.Bytes()[:4]) {
		return iTeeAuthenticatorAdmin.UnpackUnknownPcr0Error(raw[4:])
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

// ITeeAuthenticatorAdminCannotRemoveActiveImage represents a CannotRemoveActiveImage error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminCannotRemoveActiveImage struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotRemoveActiveImage()
func ITeeAuthenticatorAdminCannotRemoveActiveImageErrorID() common.Hash {
	return common.HexToHash("0xf013fa261b160780553b84b9aa781270e56afa042206ee7947cab9580d77d8ba")
}

// UnpackCannotRemoveActiveImageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotRemoveActiveImage()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackCannotRemoveActiveImageError(raw []byte) (*ITeeAuthenticatorAdminCannotRemoveActiveImage, error) {
	out := new(ITeeAuthenticatorAdminCannotRemoveActiveImage)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "CannotRemoveActiveImage", raw); err != nil {
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

// ITeeAuthenticatorAdminInvalidPcr0Length represents a InvalidPcr0Length error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminInvalidPcr0Length struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPcr0Length()
func ITeeAuthenticatorAdminInvalidPcr0LengthErrorID() common.Hash {
	return common.HexToHash("0x26e72334145c164ccb9509137816aae2f81e2287c0e397bdfbc5ad70c72e8349")
}

// UnpackInvalidPcr0LengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPcr0Length()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackInvalidPcr0LengthError(raw []byte) (*ITeeAuthenticatorAdminInvalidPcr0Length, error) {
	out := new(ITeeAuthenticatorAdminInvalidPcr0Length)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "InvalidPcr0Length", raw); err != nil {
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

// ITeeAuthenticatorAdminNoPendingSwap represents a NoPendingSwap error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminNoPendingSwap struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoPendingSwap()
func ITeeAuthenticatorAdminNoPendingSwapErrorID() common.Hash {
	return common.HexToHash("0xe10b5d062676e35e3657d2f013f557a78d0e1d69ffbb485f610f5a43ec695d9d")
}

// UnpackNoPendingSwapError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoPendingSwap()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackNoPendingSwapError(raw []byte) (*ITeeAuthenticatorAdminNoPendingSwap, error) {
	out := new(ITeeAuthenticatorAdminNoPendingSwap)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "NoPendingSwap", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminSwapAlreadyPending represents a SwapAlreadyPending error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminSwapAlreadyPending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SwapAlreadyPending()
func ITeeAuthenticatorAdminSwapAlreadyPendingErrorID() common.Hash {
	return common.HexToHash("0xde0f6f7dacb6e5557e1f444e9c2456a9d782ebb423ff52c159160da66326c5a4")
}

// UnpackSwapAlreadyPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SwapAlreadyPending()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackSwapAlreadyPendingError(raw []byte) (*ITeeAuthenticatorAdminSwapAlreadyPending, error) {
	out := new(ITeeAuthenticatorAdminSwapAlreadyPending)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "SwapAlreadyPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminSwapProposalExpired represents a SwapProposalExpired error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminSwapProposalExpired struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SwapProposalExpired()
func ITeeAuthenticatorAdminSwapProposalExpiredErrorID() common.Hash {
	return common.HexToHash("0xdfafce01b57e39382232173b40effd0615b2c9ac2d1c04dbb8b43b3da359a74c")
}

// UnpackSwapProposalExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SwapProposalExpired()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackSwapProposalExpiredError(raw []byte) (*ITeeAuthenticatorAdminSwapProposalExpired, error) {
	out := new(ITeeAuthenticatorAdminSwapProposalExpired)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "SwapProposalExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminTimelockNotElapsed represents a TimelockNotElapsed error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminTimelockNotElapsed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimelockNotElapsed()
func ITeeAuthenticatorAdminTimelockNotElapsedErrorID() common.Hash {
	return common.HexToHash("0x6677a596206ec2ddc26cf3979f15d852c7e8c99e9bacd8ea90145753aec97450")
}

// UnpackTimelockNotElapsedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimelockNotElapsed()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackTimelockNotElapsedError(raw []byte) (*ITeeAuthenticatorAdminTimelockNotElapsed, error) {
	out := new(ITeeAuthenticatorAdminTimelockNotElapsed)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "TimelockNotElapsed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ITeeAuthenticatorAdminUnknownPcr0 represents a UnknownPcr0 error raised by the ITeeAuthenticatorAdmin contract.
type ITeeAuthenticatorAdminUnknownPcr0 struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownPcr0()
func ITeeAuthenticatorAdminUnknownPcr0ErrorID() common.Hash {
	return common.HexToHash("0x4d3183dbddd956c7b047def7269cbf01a3bd69f9620049c4c3bbc2094cd92bb8")
}

// UnpackUnknownPcr0Error is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownPcr0()
func (iTeeAuthenticatorAdmin *ITeeAuthenticatorAdmin) UnpackUnknownPcr0Error(raw []byte) (*ITeeAuthenticatorAdminUnknownPcr0, error) {
	out := new(ITeeAuthenticatorAdminUnknownPcr0)
	if err := iTeeAuthenticatorAdmin.abi.UnpackIntoInterface(out, "UnknownPcr0", raw); err != nil {
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
	Bin: "0x608080604052346013576039908160188239f35b5f80fdfe5f80fdfea26469706673582212201610df450700709493d2c450bccb38ad210f5c48aab0969daef4efc751b3f75964736f6c634300081e0033",
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"contractINitroProver\",\"name\":\"_nitroProver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pcr0\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_maxVerificationAge\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_pcr0UpgradeDelay\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AttestationAlreadyUsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotRemoveActiveImage\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPCR\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPcr0Length\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUserDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoPendingSwap\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SwapAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SwapProposalExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimelockNotElapsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnknownPcr0\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStep\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0Removed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0SwapCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"}],\"name\":\"Pcr0SwapProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0Swapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PCR0_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PCR0_SWAP_APPLY_WINDOW\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"acceptedPcr0\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"acceptedPcr0List\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activeImage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"applyPcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelPcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentUpdateStep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAcceptedPcr0Count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStep2TotalLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxVerificationAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nitroProver\",\"outputs\":[{\"internalType\":\"contractINitroProver\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pcr0UpgradeDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingSwap\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"pending\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"}],\"name\":\"proposePcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"oldPcr0\",\"type\":\"bytes\"}],\"name\":\"removePcr0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step2CurrentIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTeeStep1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a8d895d4d49ffcfa5f7559dd61ad47cda6",
	Bin: "0x60e0604052346102685761378b803803806100198161026c565b928339810160a0828203126102685781516001600160a01b03811691908290036102685760208301516001600160a01b03811681036102685760408401516001600160401b0381116102685784019382601f860112156102685784516001600160401b03811161023257610096601f8201601f191660200161026c565b9581875260208701946020838301011161026857815f926020809301875e8701015260806060820151910151918415610255575f80546001600160a01b031981168717825560405196916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a360308651036102465760805260a05260c0528251812090815f52600160205260405f20600160ff19825416179055600254906801000000000000000082101561023257600182018060025582101561021e577fe4451b6eb971e43a03fd7fda060e219d0f94d4fcb6e73e481b1d9b86e5d1d22a9483604094869460025f5260205f200155600355602083525180918160208501528484015e5f828201840152601f01601f19168101030190a16040516134f99081610292823960805181818161032501528181610b040152818161155801528181611bb6015281816120bd0152612406015260a051818181610ad60152818161152a015261235a015260c05181818161080e01526124880152f35b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b6309b9c8cd60e21b5f5260045ffd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b6040519190601f01601f191682016001600160401b038111838210176102325760405256fe60806040526004361015610011575f80fd5b5f5f3560e01c806307cdb0c014612458578063081bec7e146124355780630a83fb91146123f15780630b1bfab9146123d45780630dd7ce2f146123ac5780630f5f55361461237d578063127379db146123435780632c3138441461232657806334010ccf1461208557806335732f4e1461202b5780633b3917cb14611ed85780633f964c1714611dee578063427e223b14611db457806343f855c314611d8b5780635520916c14611b1c5780635c9626b91461168b578063715018a6146116315780637637d58a146114d05780638890b84e14610a715780638da5cb5b14610a4a578063962e67a214610a2c578063a239f11d14610889578063abc620811461086d578063c3b331221461084f578063c8a032c814610831578063c8d9e1a7146107f6578063c91496c6146107da578063d22a29d114610260578063db6a7b7114610242578063f2fde38b146101b45763fe4993ca1461016f575f80fd5b346101b157806003193601126101b1576101ad604051610199816101928161278a565b0382612a02565b60405191829160208352602083019061272e565b0390f35b80fd5b50346101b15760203660031901126101b1576004356001600160a01b0381169081900361023e576101e3612f0d565b801561022a5781546001600160a01b03198116821783556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b631e4fbdf760e01b82526004829052602482fd5b5080fd5b50346101b157806003193601126101b1576020600b54604051908152f35b50346101b157806003193601126101b157610279612f0d565b6002600a54036107cb57604051635d52233760e01b8152606060048201526014805460648301819052908352829082906084600582901b83018101915f5160206134645f395f51905f5291859085015b828210610798575050505081810360031901602483015260125483916102ee82612752565b91828252846001821691825f1461077757505060011461072a575b50506103218291600319838203016044840152612c54565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa91821561071e5780809281946106f4575b5083516001600160401b0381116105495761037c601554612752565b601f8111610696575b50602094601f8211600114610628578293949582916103b89492610474575b50508160011b915f199060031b1c19161790565b6015555b82516001600160401b038111610549576103d7600c54612752565b601f81116105ca575b506020601f821160011461055d57819083949561041094926104745750508160011b915f199060031b1c19161790565b600c555b81516001600160401b0381116105495761042f601154612752565b601f81116104eb575b50602092601f821160011461047f578293829161046894926104745750508160011b915f199060031b1c19161790565b6011555b6003600a5580f35b015190505f806103a4565b60118352601f198216935f5160206133e45f395f51905f5291845b8681106104d357508360019596106104bb575b505050811b0160115561046c565b01515f1960f88460031b161c191690555f80806104ad565b9192602060018192868501518155019401920161049a565b61052e9060118452601f830160051c5f5160206133e45f395f51905f52019060208410610534575b601f0160051c5f5160206133e45f395f51905f520190612c01565b5f610438565b5f5160206133e45f395f51905f529150610513565b634e487b7160e01b82526041600452602482fd5b600c83525f5160206133645f395f51905f5290601f198316845b8181106105b25750958360019596971061059a575b505050811b01600c55610414565b01515f1960f88460031b161c191690555f808061058c565b9192602060018192868b015181550194019201610577565b61060d90600c8452601f830160051c5f5160206133645f395f51905f52019060208410610613575b601f0160051c5f5160206133645f395f51905f520190612c01565b5f6103e0565b5f5160206133645f395f51905f5291506105f2565b60158352601f198216955f5160206133c45f395f51905f5291845b88811061067e57508360019596979810610666575b505050811b016015556103bc565b01515f1960f88460031b161c191690555f8080610658565b91926020600181928685015181550194019201610643565b6106d99060158452601f830160051c5f5160206133c45f395f51905f520190602084106106df575b601f0160051c5f5160206133c45f395f51905f520190612c01565b5f610385565b5f5160206133c45f395f51905f5291506106be565b9150925061071491503d8084833e61070c8183612a02565b810190612e2a565b919091925f610360565b604051903d90823e3d90fd5b60128552849250905f5160206133a45f395f51905f525b8184106107575750500160200161032182610309565b600191955060209294508054838587010152019101909185938593610741565b61032194919550602093925060ff191683830152151560051b010191610309565b6083198886030181529295509093509160019060209082906107ba9088612966565b9601920192019285938795936102c9565b635e1e452d60e11b8152600490fd5b50346101b157806003193601126101b157602060405160858152f35b50346101b157806003193601126101b15760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b50346101b157806003193601126101b1576020600354604051908152f35b50346101b157806003193601126101b1576020601354604051908152f35b50346101b157806003193601126101b157602060405160308152f35b50346101b157610898366126df565b91906108a2612f0d565b6108ad368483612a81565b60208151910120808352600160205260ff60408420541615610a1d576003548114610a0e57808352600160205260408320805460ff19169055600254835b818103610934575b847f3e65f48ee2200e8283ed2ed80108f37f59973b9d0613c1e3b8258ebd1d8de40c858861092e604051928392602084526020840191612c17565b0390a180f35b8261093e82612a23565b90549060031b1c14610952576001016108eb565b91505f1981019081116109fa5761097b61096e61099192612a23565b90549060031b1c92612a23565b819391549060031b91821b915f19901b19161790565b90556002549283156109e6577f3e65f48ee2200e8283ed2ed80108f37f59973b9d0613c1e3b8258ebd1d8de40c92935f19016109cc81612a23565b8154905f199060031b1b1916905560025592915f806108f3565b634e487b7160e01b83526031600452602483fd5b634e487b7160e01b84526011600452602484fd5b637809fd1360e11b8352600483fd5b634d3183db60e01b8352600483fd5b50346101b157806003193601126101b157602060405162093a808152f35b50346101b157806003193601126101b157546040516001600160a01b039091168152602090f35b50346101b157610a80366126df565b610a8b929192612f0d565b610a936131bf565b610a9e368285612a81565b6020815191012080600d558252600960205260ff6040832054166114c157604051636a107b3560e11b815292829184918291610b00917f0000000000000000000000000000000000000000000000000000000000000000919060048501612c37565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156114b4578190829383918491859186916113b1575b508051906001600160401b038211610e7d57610b63601054612752565b601f8111611353575b50602090601f83116001146112e457610b9b92918891836104745750508160011b915f199060031b1c19161790565b6010555b8051906001600160401b0382116112d057610bbb600f54612752565b601f8111611272575b50602090601f831160011461120357610bf392918791836104745750508160011b915f199060031b1c19161790565b600f555b8051906001600160401b0382116111ef57610c13600e54612752565b601f8111611191575b50602090601f831160011461112257610c4b92918691836104745750508160011b915f199060031b1c19161790565b600e555b805190600160401b821161110e57601354826013558083106110cb575b50602001601384525f5160206134045f395f51905f5284915b838310610fb3575050505082516001600160401b038111610ed457610cab601254612752565b601f8111610f55575b506020601f8211600114610ee85781908495610ce49495926104745750508160011b915f199060031b1c19161790565b6012555b805190600160401b8211610ed45760145482601455808310610e91575b50602001601483525f5160206134645f395f51905f5283915b838310610d6e5784610d66604051610d398161019281612816565b604051610d498161019281612886565b60405191610d6183610d5a816128f6565b0384612a02565b612f44565b6001600a5580f35b80518051906001600160401b038211610e7d57610d8b8454612752565b601f8111610e42575b50602090601f8311600114610dda5792610dcb836001959460209487968c926104745750508160011b915f199060031b1c19161790565b85555b01920192019190610d1e565b8488528188209190601f198416895b818110610e2a5750936020936001969387969383889510610e12575b505050811b018555610dce565b01515f1960f88460031b161c191690555f8080610e05565b92936020600181928786015181550195019301610de9565b610e6d9085895260208920601f850160051c81019160208610610e73575b601f0160051c0190612c01565b5f610d94565b9091508190610e60565b634e487b7160e01b87526041600452602487fd5b601484525f5160206134645f395f51905f5201825f5160206134645f395f51905f52015b818110610ec25750610d05565b80610ece600192612cfd565b01610eb5565b634e487b7160e01b83526041600452602483fd5b601284525f5160206133a45f395f51905f5290601f198316855b818110610f3d57509583600195969710610f25575b505050811b01601255610ce8565b01515f1960f88460031b161c191690555f8080610f17565b9192602060018192868b015181550194019201610f02565b610f989060128552601f830160051c5f5160206133a45f395f51905f52019060208410610f9e575b601f0160051c5f5160206133a45f395f51905f520190612c01565b5f610cb4565b5f5160206133a45f395f51905f529150610f7d565b80518051906001600160401b0382116110b757610fd08454612752565b601f8111611087575b50602090601f831160011461101f5792611010836001959460209487968d926104745750508160011b915f199060031b1c19161790565b85555b01920192019190610c85565b8489528189209190601f1984168a5b81811061106f5750936020936001969387969383889510611057575b505050811b018555611013565b01515f1960f88460031b161c191690555f808061104a565b9293602060018192878601518155019501930161102e565b6110b190858a5260208a20601f850160051c81019160208610610e7357601f0160051c0190612c01565b5f610fd9565b634e487b7160e01b88526041600452602488fd5b601385525f5160206134045f395f51905f5201825f5160206134045f395f51905f52015b8181106110fc5750610c6c565b80611108600192612cfd565b016110ef565b634e487b7160e01b84526041600452602484fd5b600e86525f5160206133845f395f51905f529190601f198416875b8181106111795750908460019594939210611161575b505050811b01600e55610c4f565b01515f1960f88460031b161c191690555f8080611153565b9293602060018192878601518155019501930161113d565b6111d490600e8752601f840160051c5f5160206133845f395f51905f520190602085106111da575b601f0160051c5f5160206133845f395f51905f520190612c01565b5f610c1c565b5f5160206133845f395f51905f5291506111b9565b634e487b7160e01b85526041600452602485fd5b600f87525f5160206134a45f395f51905f529190601f198416885b81811061125a5750908460019594939210611242575b505050811b01600f55610bf7565b01515f1960f88460031b161c191690555f8080611234565b9293602060018192878601518155019501930161121e565b6112b590600f8852601f840160051c5f5160206134a45f395f51905f520190602085106112bb575b601f0160051c5f5160206134a45f395f51905f520190612c01565b5f610bc4565b5f5160206134a45f395f51905f52915061129a565b634e487b7160e01b86526041600452602486fd5b601088525f5160206134845f395f51905f529190601f198416895b81811061133b5750908460019594939210611323575b505050811b01601055610b9f565b01515f1960f88460031b161c191690555f8080611315565b929360206001819287860151815501950193016112ff565b6113969060108952601f840160051c5f5160206134845f395f51905f5201906020851061139c575b601f0160051c5f5160206134845f395f51905f520190612c01565b5f610b6c565b5f5160206134845f395f51905f52915061137b565b965050505050503d8082843e6113c78184612a02565b82019160c08184031261023e5780516001600160401b0381116114b057836113f0918301612e8f565b9260208201516001600160401b0381116114ac5781611410918401612d4c565b9160408101516001600160401b0381116114a85782611430918301612e8f565b9460608201516001600160401b0381116114a45783611450918401612d4c565b9260808301516001600160401b0381116114a05781611470918501612d4c565b9260a0810151906001600160401b03821161149c57611490929101612d4c565b9093959291905f610b46565b8780fd5b8680fd5b8580fd5b8480fd5b8380fd5b8280fd5b50604051903d90823e3d90fd5b6305fb89e760e21b8252600482fd5b50346101b1576114df366126df565b906114e8612f0d565b6114f3368383612a81565b6020815191012091828452600960205260ff60408520541661162257604051632b5f2f8160e01b81529291849184918291611554917f0000000000000000000000000000000000000000000000000000000000000000919060048501612c37565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa8015611617576115cf928491859186916115f2575b5081836115a392612f44565b60208151910151906bffffffffffffffffffffffff19821691601482106115d2575b505060601c612ffa565b80f35b6001600160601b031960149290920360031b82901b161690505f806115c5565b90506115a3925061160d91503d8087833e61070c8183612a02565b9192509081611597565b6040513d85823e3d90fd5b6305fb89e760e21b8452600484fd5b50346101b157806003193601126101b15761164a612f0d565b80546001600160a01b03198116825581906001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b50346101b15760403660031901126101b1576004356001600160401b03811161023e57610160600319823603011261023e576040519161016083018381106001600160401b038211176105495760405281600401356001600160401b038116810361023e5783526020830160248301358152604084016044840135815260608501906064850135825260848501356001600160401b0381116114a8576117379060043691880101612ad5565b956080810196875260a48601356001600160401b0381116114a4576117629060043691890101612ad5565b9560a0820196875260c48101356001600160401b0381116114a057810190366023830112156114a057600482013561179981612a4f565b926117a76040519485612a02565b81845260206004606082870194028301010190368211611ab557602401915b818310611ab95750505060c0830191825260e083019260e48201358452610100810190610104830135825261012483013592600d841015611ab5576101208201938452610144810135906001600160401b038211611ab1570136602382011215611ab55761183e903690602460048201359101612a81565b9361014082019485526024356001600160401b038111611ab157611866903690600401612ab7565b6007546001600160a01b03169c909b908d158015611a96575b611a875760208151516040516118aa8161189c8582019485612d92565b03601f198101835282612a02565b519020915101516040516118c68161189c602082019485612df1565b519020908c60208451516040516118e48161189c8582019485612d92565b519020945101516040516119008161189c602082019485612df1565b5190209451604051906020820192604083016020855282518091526020606085019301915b818110611a4a575050509061194c816001600160401b03949303601f198101835282612a02565b5190209551169b5199519a5198519551965197600d891015611a365751986040519b8c9b60208d019e8f5260408d015260608c015260808b015260a08a015260c089015260e08801526101008701526101208601526101408501526101608401526101808301526101a082016101a090526101c082016119cb9161272e565b03601f19810182526119dd9082612a02565b5190207f19457468657265756d205369676e6564204d6573736167653a0a3332000000008252601c52603c902090611a1491613233565b611a209193929361326d565b6040516001600160a01b03909216148152602090f35b634e487b7160e01b8e52602160045260248efd5b825180516001600160a01b039081168652602082810151909116818701526040918201519186019190915260609094019390920191600101611925565b63f64b1d7b60e01b8c5260048cfd5b506085604051611aa9816101928161278a565b51141561187f565b8a80fd5b8980fd5b606083360312611ab557604051611acf816129e7565b83356001600160a01b0381168103611b1857815260208401356001600160a01b0381168103611b18579181606093602080940152604086013560408201528152019201916117c6565b8b80fd5b50346101b157806003193601126101b157611b35612f0d565b6001600a54036107cb57600b5460405163c13734e760e01b8152606060048201526013805460648301819052908452909291829184916084600583901b84018101925f5160206134045f395f51905f5291869086015b828210611d57575050505060248301528181036003190160448301528190611bb290612c54565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156114b4578192611d18575b5081516001600160401b03811161054957611c0a600c54612752565b601f8111611cd0575b50602092601f8211600114611c645782938291611c4394926104745750508160011b915f199060031b1c19161790565b600c555b6001600b540180600b5560135414611c5c5780f35b6002600a5580f35b600c8352601f198216935f5160206133645f395f51905f5291845b868110611cb85750836001959610611ca0575b505050811b01600c55611c47565b01515f1960f88460031b161c191690555f8080611c92565b91926020600181928685015181550194019201611c7f565b611d1290600c8452601f830160051c5f5160206133645f395f51905f5201906020841061061357601f0160051c5f5160206133645f395f51905f520190612c01565b5f611c13565b9091503d8083833e611d2a8183612a02565b81016020828203126114b05781516001600160401b0381116114ac57611d509201612d4c565b905f611bee565b6083198a8703018152929650929450926001906020908290611d799089612966565b97019201920192869593889593611b8b565b50346101b157806003193601126101b1576007546040516001600160a01b039091168152602090f35b50346101b15760203660031901126101b157600435906002548210156101b1576020611ddf83612a23565b90549060031b1c604051908152f35b50346101b157806003193601126101b157604051908060045490611e1182612752565b8085529160018116908115611eb15750600114611e67575b611e5684611e3981860382612a02565b60055460ff6006541660405193849360608552606085019061272e565b916020840152151560408301520390f35b600481525f5160206134445f395f51905f52939250905b808210611e9757509091508101602001611e3982611e29565b919260018160209254838588010152019101909291611e7e565b60ff191660208087019190915292151560051b85019092019250611e399150839050611e29565b50346101b157806003193601126101b157611ef1612f0d565b611ef9612cc4565b60408101511561201c5760208101805162093a8081018091116109fa57421161200d5781516020815191012090818452600160205260ff60408520541615611f8d575b507fe4451b6eb971e43a03fd7fda060e219d0f94d4fcb6e73e481b1d9b86e5d1d22a9161092e91600355611f706004612cfd565b5f6005555f6006555160405191829160208352602083019061272e565b514210611ffe57808352600160205260408320600160ff19825416179055600254600160401b81101561110e57918161092e92611ff461097b8660017fe4451b6eb971e43a03fd7fda060e219d0f94d4fcb6e73e481b1d9b86e5d1d22a9801600255612a23565b9055915091611f3c565b63333bd2cb60e11b8352600483fd5b63dfafce0160e01b8352600483fd5b637085ae8360e11b8252600482fd5b50346101b157806003193601126101b157612044612f0d565b61204c612cc4565b60408101511561201c5761092e7f89d9ae00acbd51f480eb5a0d3462ae199cb9e6076f950b9b2d64e4b18cf951d391611f706004612cfd565b5034612313575f3660031901126123135761209e612f0d565b6003600a5403612317576120bb604051610d398161019281612816565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156123135760405190632d77dd9360e11b825260606004830152815f60115461210f81612752565b908160648501526001811690815f146122ee57506001146122a7575b5081810360031901602483015261214190612c54565b8181036003190160448301526015545f9161215b82612752565b80825291600181169081156122835750600114612236575b50509181805f9403915afa801561222b57612218575b506115cf612198600f54612752565b600f601f8211612200575b546001600160601b031981169190601482106121e0575b5050600d5490604051906121d8826121d181612886565b0383612a02565b60601c612ffa565b6001600160601b031960149290920360031b82901b161690505f806121ba565b50600f83525f5160206134a45f395f51905f526121a3565b61222491505f90612a02565b5f5f612189565b6040513d5f823e3d90fd5b60155f9081525f5160206133c45f395f51905f52959350915b8083106122655750919350016020018180612173565b9350919360018160209254838587010152019101909391859361224f565b60ff191660208381019190915292151560051b909101909101915082905080612173565b60115f90815291505f5160206133e45f395f51905f525b8183106122d4575050810160840161214161212b565b8054838701608401528593506020909201916001016122be565b6121419350606491509160209260ff19166084860152151560051b840101019061212b565b5f80fd5b635e1e452d60e11b5f5260045ffd5b34612313575f366003190112612313576020600254604051908152f35b34612313575f3660031901126123135760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b34612313576020366003190112612313576004355f526001602052602060ff60405f2054166040519015158152f35b34612313575f366003190112612313576007546040516001600160a01b039091168152602090f35b34612313575f366003190112612313576020600a54604051908152f35b34612313575f366003190112612313576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b34612313575f366003190112612313576101ad604051610199816101928161278a565b3461231357612466366126df565b61246e612f0d565b603081036126d05760ff60065416806126b7575b6126a8577f000000000000000000000000000000000000000000000000000000000000000042019182421161269457604051926124be846129e7565b6124c9368484612a81565b938481526040602082019183835201906001825285516001600160401b038111612680576124f8600454612752565b601f8111612622575b506020601f821160011461258c578161257c95949392612557927f8aac65bd065b693837bd35d5d39a8150c0bace29c9db8406bedabba66ca447c09a5f926125815750508160011b915f199060031b1c19161790565b6004555b5160055551151560ff80196006541691161760065560405193849384612c37565b0390a1005b015190508a806103a4565b601f1982169760045f525f5160206134445f395f51905f52985f5b81811061260a575098600192849261257c989796957f8aac65bd065b693837bd35d5d39a8150c0bace29c9db8406bedabba66ca447c09c106125f2575b505050811b0160045561255b565b01515f1960f88460031b161c191690558980806125e4565b838301518b556001909a0199602093840193016125a7565b6126659060045f52601f830160051c5f5160206134445f395f51905f5201906020841061266b575b601f0160051c5f5160206134445f395f51905f520190612c01565b87612501565b5f5160206134445f395f51905f52915061264a565b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b63de0f6f7d60e01b5f5260045ffd5b5060055462093a80810180911161269457421115612482565b6309b9c8cd60e21b5f5260045ffd5b906020600319830112612313576004356001600160401b0381116123135782602382011215612313578060040135926001600160401b0384116123135760248483010111612313576024019190565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b90600182811c92168015612780575b602083101461276c57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612761565b6008545f929161279982612752565b80825291600181169081156127fa57506001146127b4575050565b60085f9081529293509091905f5160206134245f395f51905f525b8383106127e0575060209250010190565b6001816020929493945483858701015201910191906127cf565b9050602093945060ff929192191683830152151560051b010190565b6010545f929161282582612752565b80825291600181169081156127fa5750600114612840575050565b60105f9081529293509091905f5160206134845f395f51905f525b83831061286c575060209250010190565b60018160209294939454838587010152019101919061285b565b600e545f929161289582612752565b80825291600181169081156127fa57506001146128b0575050565b600e5f9081529293509091905f5160206133845f395f51905f525b8383106128dc575060209250010190565b6001816020929493945483858701015201910191906128cb565b600f545f929161290582612752565b80825291600181169081156127fa5750600114612920575050565b600f5f9081529293509091905f5160206134a45f395f51905f525b83831061294c575060209250010190565b60018160209294939454838587010152019101919061293b565b5f929181549161297583612752565b80835292600181169081156129ca575060011461299157505050565b5f9081526020812093945091925b8383106129b0575060209250010190565b60018160209294939454838587010152019101919061299f565b915050602093945060ff929192191683830152151560051b010190565b606081019081106001600160401b0382111761268057604052565b90601f801991011681019081106001600160401b0382111761268057604052565b600254811015612a3b5760025f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b6001600160401b0381116126805760051b60200190565b6001600160401b03811161268057601f01601f191660200190565b929192612a8d82612a66565b91612a9b6040519384612a02565b829481845281830111612313578281602093845f960137010152565b9080601f8301121561231357816020612ad293359101612a81565b90565b919060408382031261231357604051604081018181106001600160401b0382111761268057604052809380356001600160401b03811161231357810183601f82011215612313578035612b2781612a4f565b91612b356040519384612a02565b81835260208084019260051b820101918683116123135760208201905b838210612bd4575050505082526020810135906001600160401b03821161231357019180601f84011215612313578235612b8b81612a4f565b93612b996040519586612a02565b81855260208086019260051b82010192831161231357602001905b828210612bc45750505060200152565b8135815260209182019101612bb4565b81356001600160401b03811161231357602091612bf68a848094880101612ab7565b815201910190612b52565b818110612c0c575050565b5f8155600101612c01565b908060209392818452848401375f828201840152601f01601f1916010190565b939291602091612c4f91604087526040870191612c17565b930152565b600c545f9291612c6382612752565b80825291600181169081156127fa5750600114612c7e575050565b600c5f9081529293509091905f5160206133645f395f51905f525b838310612caa575060209250010190565b600181602092949394548385870101520191019190612c99565b60405190612cd1826129e7565b81604051612ce481610192816004612966565b81526005546020820152604060ff600654161515910152565b612d078154612752565b9081612d11575050565b81601f5f9311600114612d235750555b565b81835260208320612d3f91601f0160051c810190600101612c01565b8082528160208120915555565b81601f8201121561231357805190612d6382612a66565b92612d716040519485612a02565b8284526020838301011161231357815f9260208093018386015e8301015290565b602081016020825282518091526040820191602060408360051b8301019401925f915b838310612dc457505050505090565b9091929394602080612de2600193603f19868203018752895161272e565b97019301930191939290612db5565b60206040818301928281528451809452019201905f5b818110612e145750505090565b8251845260209384019390920191600101612e07565b916060838303126123135782516001600160401b0381116123135782612e51918501612d4c565b9260208101516001600160401b0381116123135783612e71918301612d4c565b9260408201516001600160401b03811161231357612ad29201612d4c565b9080601f83011215612313578151612ea681612a4f565b92612eb46040519485612a02565b81845260208085019260051b820101918383116123135760208201905b838210612ee057505050505090565b81516001600160401b03811161231357602091612f0287848094880101612d4c565b815201910190612ed1565b5f546001600160a01b03163303612f2057565b63118cdaa760e01b5f523360045260245ffd5b908151811015612a3b570160200190565b915160131901612feb575160841901612fdc576034815110612f9557604051612f6e606082612a02565b60308152602081019160403684375f5b60308103612fa457505051902060035403612f9557565b63719a9a4960e01b5f5260045ffd5b6004810190818111612694576001916001600160f81b031990612fc79085612f33565b51165f1a612fd58286612f33565b5301612f7e565b634ae601c160e11b5f5260045ffd5b6342f6e2a360e11b5f5260045ffd5b9092917f1882fd6e997ef1dc5e2691efe78564a5c540b36547363e8a7dcd454c61a944b660018060a01b03600754169260405193845260018060a01b03169283602082015260806040820152806130646130566080830161278a565b82810360608401528861272e565b0390a15f52600960205260405f20600160ff198254161790556bffffffffffffffffffffffff60a01b600754161760075581516001600160401b038111612680576130b0600854612752565b601f8111613161575b50602092601f82116001146130f5576130e9929382915f926104745750508160011b915f199060031b1c19161790565b6008555b612d216131bf565b601f1982169360085f525f5160206134245f395f51905f52915f5b8681106131495750836001959610613131575b505050811b016008556130ed565b01515f1960f88460031b161c191690555f8080613123565b91926020600181928685015181550194019201613110565b6131a49060085f52601f830160051c5f5160206134245f395f51905f520190602084106131aa575b601f0160051c5f5160206134245f395f51905f520190612c01565b5f6130b9565b5f5160206134245f395f51905f529150613189565b5f600a555f600b555f6040516131d6602082612a02565b526131e2600c54612752565b601f81116131f2575b505f600c55565b601f90600c5f520160051c5f5160206133645f395f51905f52015f5160206133645f395f51905f525b81811061322857506131eb565b5f815560010161321b565b81519190604183036132635761325c9250602082015190606060408401519301515f1a906132e1565b9192909190565b50505f9160029190565b60048110156132cd578061327f575050565b600181036132965763f645eedf60e01b5f5260045ffd5b600281036132b1575063fce698f760e01b5f5260045260245ffd5b6003146132bb5750565b6335e2f38360e21b5f5260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411613358579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa1561222b575f516001600160a01b0381161561334e57905f905f90565b505f906001905f90565b5050505f916003919056fedf6966c971051c3d54ec59162606531493a51404a002842f56009d7e5cf4a8c7bb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fdbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec344455f448fdea98c4d29eb340757ef0a66cd03dbb9538908a6a81d96026b71ec47531ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c6866de8ffda797e3de9c05e8fc57b3bf0ec28a930d40b0d285d93c06501cf6a090f3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee38a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19bce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec1b6847dc741a1b0cd08d278845f9d819d87b734759afb55fe2de5cb82a9ae6728d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac802a264697066735822122093ed51e0609df6a18a838b224099d2014b5a210fa9d4ef330df341145ed42e3d64736f6c634300081e0033",
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
// Solidity: constructor(address owner, address _nitroProver, bytes _pcr0, uint256 _maxVerificationAge, uint256 _pcr0UpgradeDelay) returns()
func (teeAuthenticator *TeeAuthenticator) PackConstructor(owner common.Address, _nitroProver common.Address, _pcr0 []byte, _maxVerificationAge *big.Int, _pcr0UpgradeDelay *big.Int) []byte {
	enc, err := teeAuthenticator.abi.Pack("", owner, _nitroProver, _pcr0, _maxVerificationAge, _pcr0UpgradeDelay)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackPCR0LENGTH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xabc62081.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function PCR0_LENGTH() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackPCR0LENGTH() []byte {
	enc, err := teeAuthenticator.abi.Pack("PCR0_LENGTH")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPCR0LENGTH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xabc62081.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function PCR0_LENGTH() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackPCR0LENGTH() ([]byte, error) {
	return teeAuthenticator.abi.Pack("PCR0_LENGTH")
}

// UnpackPCR0LENGTH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xabc62081.
//
// Solidity: function PCR0_LENGTH() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackPCR0LENGTH(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("PCR0_LENGTH", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPCR0SWAPAPPLYWINDOW is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x962e67a2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function PCR0_SWAP_APPLY_WINDOW() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackPCR0SWAPAPPLYWINDOW() []byte {
	enc, err := teeAuthenticator.abi.Pack("PCR0_SWAP_APPLY_WINDOW")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPCR0SWAPAPPLYWINDOW is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x962e67a2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function PCR0_SWAP_APPLY_WINDOW() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackPCR0SWAPAPPLYWINDOW() ([]byte, error) {
	return teeAuthenticator.abi.Pack("PCR0_SWAP_APPLY_WINDOW")
}

// UnpackPCR0SWAPAPPLYWINDOW is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x962e67a2.
//
// Solidity: function PCR0_SWAP_APPLY_WINDOW() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackPCR0SWAPAPPLYWINDOW(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("PCR0_SWAP_APPLY_WINDOW", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
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

// PackAcceptedPcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f5f5536.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptedPcr0(bytes32 ) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) PackAcceptedPcr0(arg0 [32]byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("acceptedPcr0", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptedPcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f5f5536.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptedPcr0(bytes32 ) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) TryPackAcceptedPcr0(arg0 [32]byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("acceptedPcr0", arg0)
}

// UnpackAcceptedPcr0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0f5f5536.
//
// Solidity: function acceptedPcr0(bytes32 ) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) UnpackAcceptedPcr0(data []byte) (bool, error) {
	out, err := teeAuthenticator.abi.Unpack("acceptedPcr0", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackAcceptedPcr0List is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x427e223b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptedPcr0List(uint256 ) view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) PackAcceptedPcr0List(arg0 *big.Int) []byte {
	enc, err := teeAuthenticator.abi.Pack("acceptedPcr0List", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptedPcr0List is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x427e223b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptedPcr0List(uint256 ) view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) TryPackAcceptedPcr0List(arg0 *big.Int) ([]byte, error) {
	return teeAuthenticator.abi.Pack("acceptedPcr0List", arg0)
}

// UnpackAcceptedPcr0List is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x427e223b.
//
// Solidity: function acceptedPcr0List(uint256 ) view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) UnpackAcceptedPcr0List(data []byte) ([32]byte, error) {
	out, err := teeAuthenticator.abi.Unpack("acceptedPcr0List", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackActiveImage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc8a032c8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function activeImage() view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) PackActiveImage() []byte {
	enc, err := teeAuthenticator.abi.Pack("activeImage")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackActiveImage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc8a032c8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function activeImage() view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) TryPackActiveImage() ([]byte, error) {
	return teeAuthenticator.abi.Pack("activeImage")
}

// UnpackActiveImage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc8a032c8.
//
// Solidity: function activeImage() view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) UnpackActiveImage(data []byte) ([32]byte, error) {
	out, err := teeAuthenticator.abi.Unpack("activeImage", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackApplyPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b3917cb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function applyPcr0Swap() returns()
func (teeAuthenticator *TeeAuthenticator) PackApplyPcr0Swap() []byte {
	enc, err := teeAuthenticator.abi.Pack("applyPcr0Swap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApplyPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b3917cb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function applyPcr0Swap() returns()
func (teeAuthenticator *TeeAuthenticator) TryPackApplyPcr0Swap() ([]byte, error) {
	return teeAuthenticator.abi.Pack("applyPcr0Swap")
}

// PackCancelPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35732f4e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelPcr0Swap() returns()
func (teeAuthenticator *TeeAuthenticator) PackCancelPcr0Swap() []byte {
	enc, err := teeAuthenticator.abi.Pack("cancelPcr0Swap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelPcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35732f4e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelPcr0Swap() returns()
func (teeAuthenticator *TeeAuthenticator) TryPackCancelPcr0Swap() ([]byte, error) {
	return teeAuthenticator.abi.Pack("cancelPcr0Swap")
}

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c9626b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c9626b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c9626b9.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
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

// PackGetAcceptedPcr0Count is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c313844.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getAcceptedPcr0Count() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackGetAcceptedPcr0Count() []byte {
	enc, err := teeAuthenticator.abi.Pack("getAcceptedPcr0Count")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetAcceptedPcr0Count is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c313844.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getAcceptedPcr0Count() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackGetAcceptedPcr0Count() ([]byte, error) {
	return teeAuthenticator.abi.Pack("getAcceptedPcr0Count")
}

// UnpackGetAcceptedPcr0Count is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2c313844.
//
// Solidity: function getAcceptedPcr0Count() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackGetAcceptedPcr0Count(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("getAcceptedPcr0Count", data)
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

// PackPcr0UpgradeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc8d9e1a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pcr0UpgradeDelay() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) PackPcr0UpgradeDelay() []byte {
	enc, err := teeAuthenticator.abi.Pack("pcr0UpgradeDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPcr0UpgradeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc8d9e1a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pcr0UpgradeDelay() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) TryPackPcr0UpgradeDelay() ([]byte, error) {
	return teeAuthenticator.abi.Pack("pcr0UpgradeDelay")
}

// UnpackPcr0UpgradeDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc8d9e1a7.
//
// Solidity: function pcr0UpgradeDelay() view returns(uint256)
func (teeAuthenticator *TeeAuthenticator) UnpackPcr0UpgradeDelay(data []byte) (*big.Int, error) {
	out, err := teeAuthenticator.abi.Unpack("pcr0UpgradeDelay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPendingSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f964c17.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingSwap() view returns(bytes value, uint256 eta, bool pending)
func (teeAuthenticator *TeeAuthenticator) PackPendingSwap() []byte {
	enc, err := teeAuthenticator.abi.Pack("pendingSwap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f964c17.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingSwap() view returns(bytes value, uint256 eta, bool pending)
func (teeAuthenticator *TeeAuthenticator) TryPackPendingSwap() ([]byte, error) {
	return teeAuthenticator.abi.Pack("pendingSwap")
}

// PendingSwapOutput serves as a container for the return parameters of contract
// method PendingSwap.
type PendingSwapOutput struct {
	Value   []byte
	Eta     *big.Int
	Pending bool
}

// UnpackPendingSwap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3f964c17.
//
// Solidity: function pendingSwap() view returns(bytes value, uint256 eta, bool pending)
func (teeAuthenticator *TeeAuthenticator) UnpackPendingSwap(data []byte) (PendingSwapOutput, error) {
	out, err := teeAuthenticator.abi.Unpack("pendingSwap", data)
	outstruct := new(PendingSwapOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Value = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Eta = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Pending = *abi.ConvertType(out[2], new(bool)).(*bool)
	return *outstruct, nil
}

// PackProposePcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07cdb0c0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proposePcr0Swap(bytes targetPcr0) returns()
func (teeAuthenticator *TeeAuthenticator) PackProposePcr0Swap(targetPcr0 []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("proposePcr0Swap", targetPcr0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProposePcr0Swap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07cdb0c0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proposePcr0Swap(bytes targetPcr0) returns()
func (teeAuthenticator *TeeAuthenticator) TryPackProposePcr0Swap(targetPcr0 []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("proposePcr0Swap", targetPcr0)
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

// PackRemovePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa239f11d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function removePcr0(bytes oldPcr0) returns()
func (teeAuthenticator *TeeAuthenticator) PackRemovePcr0(oldPcr0 []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("removePcr0", oldPcr0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRemovePcr0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa239f11d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function removePcr0(bytes oldPcr0) returns()
func (teeAuthenticator *TeeAuthenticator) TryPackRemovePcr0(oldPcr0 []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("removePcr0", oldPcr0)
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

// TeeAuthenticatorPcr0Removed represents a Pcr0Removed event raised by the TeeAuthenticator contract.
type TeeAuthenticatorPcr0Removed struct {
	Pcr0 []byte
	Raw  *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorPcr0RemovedEventName = "Pcr0Removed"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorPcr0Removed) ContractEventName() string {
	return TeeAuthenticatorPcr0RemovedEventName
}

// UnpackPcr0RemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0Removed(bytes pcr0)
func (teeAuthenticator *TeeAuthenticator) UnpackPcr0RemovedEvent(log *types.Log) (*TeeAuthenticatorPcr0Removed, error) {
	event := "Pcr0Removed"
	if log.Topics[0] != teeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TeeAuthenticatorPcr0Removed)
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

// TeeAuthenticatorPcr0SwapCancelled represents a Pcr0SwapCancelled event raised by the TeeAuthenticator contract.
type TeeAuthenticatorPcr0SwapCancelled struct {
	TargetPcr0 []byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorPcr0SwapCancelledEventName = "Pcr0SwapCancelled"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorPcr0SwapCancelled) ContractEventName() string {
	return TeeAuthenticatorPcr0SwapCancelledEventName
}

// UnpackPcr0SwapCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0SwapCancelled(bytes targetPcr0)
func (teeAuthenticator *TeeAuthenticator) UnpackPcr0SwapCancelledEvent(log *types.Log) (*TeeAuthenticatorPcr0SwapCancelled, error) {
	event := "Pcr0SwapCancelled"
	if log.Topics[0] != teeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TeeAuthenticatorPcr0SwapCancelled)
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

// TeeAuthenticatorPcr0SwapProposed represents a Pcr0SwapProposed event raised by the TeeAuthenticator contract.
type TeeAuthenticatorPcr0SwapProposed struct {
	TargetPcr0 []byte
	Eta        *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorPcr0SwapProposedEventName = "Pcr0SwapProposed"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorPcr0SwapProposed) ContractEventName() string {
	return TeeAuthenticatorPcr0SwapProposedEventName
}

// UnpackPcr0SwapProposedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0SwapProposed(bytes targetPcr0, uint256 eta)
func (teeAuthenticator *TeeAuthenticator) UnpackPcr0SwapProposedEvent(log *types.Log) (*TeeAuthenticatorPcr0SwapProposed, error) {
	event := "Pcr0SwapProposed"
	if log.Topics[0] != teeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TeeAuthenticatorPcr0SwapProposed)
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

// TeeAuthenticatorPcr0Swapped represents a Pcr0Swapped event raised by the TeeAuthenticator contract.
type TeeAuthenticatorPcr0Swapped struct {
	Pcr0 []byte
	Raw  *types.Log // Blockchain specific contextual infos
}

const TeeAuthenticatorPcr0SwappedEventName = "Pcr0Swapped"

// ContractEventName returns the user-defined event name.
func (TeeAuthenticatorPcr0Swapped) ContractEventName() string {
	return TeeAuthenticatorPcr0SwappedEventName
}

// UnpackPcr0SwappedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Pcr0Swapped(bytes pcr0)
func (teeAuthenticator *TeeAuthenticator) UnpackPcr0SwappedEvent(log *types.Log) (*TeeAuthenticatorPcr0Swapped, error) {
	event := "Pcr0Swapped"
	if log.Topics[0] != teeAuthenticator.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TeeAuthenticatorPcr0Swapped)
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
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["CannotRemoveActiveImage"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackCannotRemoveActiveImageError(raw[4:])
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
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["InvalidPcr0Length"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackInvalidPcr0LengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["InvalidUserDataLength"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackInvalidUserDataLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["NoPendingSwap"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackNoPendingSwapError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["SwapAlreadyPending"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackSwapAlreadyPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["SwapProposalExpired"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackSwapProposalExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["TimelockNotElapsed"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackTimelockNotElapsedError(raw[4:])
	}
	if bytes.Equal(raw[:4], teeAuthenticator.abi.Errors["UnknownPcr0"].ID.Bytes()[:4]) {
		return teeAuthenticator.UnpackUnknownPcr0Error(raw[4:])
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

// TeeAuthenticatorCannotRemoveActiveImage represents a CannotRemoveActiveImage error raised by the TeeAuthenticator contract.
type TeeAuthenticatorCannotRemoveActiveImage struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotRemoveActiveImage()
func TeeAuthenticatorCannotRemoveActiveImageErrorID() common.Hash {
	return common.HexToHash("0xf013fa261b160780553b84b9aa781270e56afa042206ee7947cab9580d77d8ba")
}

// UnpackCannotRemoveActiveImageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotRemoveActiveImage()
func (teeAuthenticator *TeeAuthenticator) UnpackCannotRemoveActiveImageError(raw []byte) (*TeeAuthenticatorCannotRemoveActiveImage, error) {
	out := new(TeeAuthenticatorCannotRemoveActiveImage)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "CannotRemoveActiveImage", raw); err != nil {
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

// TeeAuthenticatorInvalidPcr0Length represents a InvalidPcr0Length error raised by the TeeAuthenticator contract.
type TeeAuthenticatorInvalidPcr0Length struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPcr0Length()
func TeeAuthenticatorInvalidPcr0LengthErrorID() common.Hash {
	return common.HexToHash("0x26e72334145c164ccb9509137816aae2f81e2287c0e397bdfbc5ad70c72e8349")
}

// UnpackInvalidPcr0LengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPcr0Length()
func (teeAuthenticator *TeeAuthenticator) UnpackInvalidPcr0LengthError(raw []byte) (*TeeAuthenticatorInvalidPcr0Length, error) {
	out := new(TeeAuthenticatorInvalidPcr0Length)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "InvalidPcr0Length", raw); err != nil {
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

// TeeAuthenticatorNoPendingSwap represents a NoPendingSwap error raised by the TeeAuthenticator contract.
type TeeAuthenticatorNoPendingSwap struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoPendingSwap()
func TeeAuthenticatorNoPendingSwapErrorID() common.Hash {
	return common.HexToHash("0xe10b5d062676e35e3657d2f013f557a78d0e1d69ffbb485f610f5a43ec695d9d")
}

// UnpackNoPendingSwapError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoPendingSwap()
func (teeAuthenticator *TeeAuthenticator) UnpackNoPendingSwapError(raw []byte) (*TeeAuthenticatorNoPendingSwap, error) {
	out := new(TeeAuthenticatorNoPendingSwap)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "NoPendingSwap", raw); err != nil {
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

// TeeAuthenticatorSwapAlreadyPending represents a SwapAlreadyPending error raised by the TeeAuthenticator contract.
type TeeAuthenticatorSwapAlreadyPending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SwapAlreadyPending()
func TeeAuthenticatorSwapAlreadyPendingErrorID() common.Hash {
	return common.HexToHash("0xde0f6f7dacb6e5557e1f444e9c2456a9d782ebb423ff52c159160da66326c5a4")
}

// UnpackSwapAlreadyPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SwapAlreadyPending()
func (teeAuthenticator *TeeAuthenticator) UnpackSwapAlreadyPendingError(raw []byte) (*TeeAuthenticatorSwapAlreadyPending, error) {
	out := new(TeeAuthenticatorSwapAlreadyPending)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "SwapAlreadyPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorSwapProposalExpired represents a SwapProposalExpired error raised by the TeeAuthenticator contract.
type TeeAuthenticatorSwapProposalExpired struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SwapProposalExpired()
func TeeAuthenticatorSwapProposalExpiredErrorID() common.Hash {
	return common.HexToHash("0xdfafce01b57e39382232173b40effd0615b2c9ac2d1c04dbb8b43b3da359a74c")
}

// UnpackSwapProposalExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SwapProposalExpired()
func (teeAuthenticator *TeeAuthenticator) UnpackSwapProposalExpiredError(raw []byte) (*TeeAuthenticatorSwapProposalExpired, error) {
	out := new(TeeAuthenticatorSwapProposalExpired)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "SwapProposalExpired", raw); err != nil {
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

// TeeAuthenticatorTimelockNotElapsed represents a TimelockNotElapsed error raised by the TeeAuthenticator contract.
type TeeAuthenticatorTimelockNotElapsed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimelockNotElapsed()
func TeeAuthenticatorTimelockNotElapsedErrorID() common.Hash {
	return common.HexToHash("0x6677a596206ec2ddc26cf3979f15d852c7e8c99e9bacd8ea90145753aec97450")
}

// UnpackTimelockNotElapsedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimelockNotElapsed()
func (teeAuthenticator *TeeAuthenticator) UnpackTimelockNotElapsedError(raw []byte) (*TeeAuthenticatorTimelockNotElapsed, error) {
	out := new(TeeAuthenticatorTimelockNotElapsed)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "TimelockNotElapsed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TeeAuthenticatorUnknownPcr0 represents a UnknownPcr0 error raised by the TeeAuthenticator contract.
type TeeAuthenticatorUnknownPcr0 struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownPcr0()
func TeeAuthenticatorUnknownPcr0ErrorID() common.Hash {
	return common.HexToHash("0x4d3183dbddd956c7b047def7269cbf01a3bd69f9620049c4c3bbc2094cd92bb8")
}

// UnpackUnknownPcr0Error is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownPcr0()
func (teeAuthenticator *TeeAuthenticator) UnpackUnknownPcr0Error(raw []byte) (*TeeAuthenticatorUnknownPcr0, error) {
	out := new(TeeAuthenticatorUnknownPcr0)
	if err := teeAuthenticator.abi.UnpackIntoInterface(out, "UnknownPcr0", raw); err != nil {
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
