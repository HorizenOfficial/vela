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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"contractINitroProver\",\"name\":\"_nitroProver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pcr0\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_maxVerificationAge\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_pcr0UpgradeDelay\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AttestationAlreadyUsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotRemoveActiveImage\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPCR\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPcr0Length\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUserDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoPendingSwap\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SwapAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SwapProposalExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimelockNotElapsed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnknownPcr0\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStep\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0Removed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0SwapCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"}],\"name\":\"Pcr0SwapProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pcr0\",\"type\":\"bytes\"}],\"name\":\"Pcr0Swapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PCR0_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PCR0_SWAP_APPLY_WINDOW\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"acceptedPcr0\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"acceptedPcr0List\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activeImage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"applyPcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelPcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentUpdateStep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAcceptedPcr0Count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAcceptedPcr0List\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveImage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingSwap\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"pending\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStep2TotalLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxVerificationAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nitroProver\",\"outputs\":[{\"internalType\":\"contractINitroProver\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pcr0UpgradeDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingSwap\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"pending\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"targetPcr0\",\"type\":\"bytes\"}],\"name\":\"proposePcr0Swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"oldPcr0\",\"type\":\"bytes\"}],\"name\":\"removePcr0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step2CurrentIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTeeStep1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a8d895d4d49ffcfa5f7559dd61ad47cda6",
	Bin: "0x60e06040523461026857613889803803806100198161026c565b928339810160a0828203126102685781516001600160a01b03811691908290036102685760208301516001600160a01b03811681036102685760408401516001600160401b0381116102685784019382601f860112156102685784516001600160401b03811161023257610096601f8201601f191660200161026c565b9581875260208701946020838301011161026857815f926020809301875e8701015260806060820151910151918415610255575f80546001600160a01b031981168717825560405196916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a360308651036102465760805260a05260c0528251812090815f52600160205260405f20600160ff19825416179055600254906801000000000000000082101561023257600182018060025582101561021e577fe4451b6eb971e43a03fd7fda060e219d0f94d4fcb6e73e481b1d9b86e5d1d22a9483604094869460025f5260205f200155600355602083525180918160208501528484015e5f828201840152601f01601f19168101030190a16040516135f79081610292823960805181818161034601528181610b250152818161157901528181611bd7015281816120fa01526124e0015260a051818181610af70152818161154b0152612434015260c05181818161082f01526125620152f35b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b6309b9c8cd60e21b5f5260045ffd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b6040519190601f01601f191682016001600160401b038111838210176102325760405256fe60806040526004361015610011575f80fd5b5f5f3560e01c806307cdb0c014612532578063081bec7e1461250f5780630a83fb91146124cb5780630b1bfab9146124ae5780630dd7ce2f146124865780630f5f553614612457578063127379db1461241d57806314758b511461239d578063187cd3fb146123805780632c3138441461236357806334010ccf146120c257806335732f4e146120685780633b3917cb14611f155780633f964c1714611e45578063427e223b14611e0b57806343f855c314611de25780635395c3ce14611dac5780635520916c14611b3d5780635c9626b9146116ac578063715018a6146116525780637637d58a146114f15780638890b84e14610a925780638da5cb5b14610a6b578063962e67a214610a4d578063a239f11d146108aa578063abc620811461088e578063c3b3312214610870578063c8a032c814610852578063c8d9e1a714610817578063c91496c6146107fb578063d22a29d114610281578063db6a7b7114610263578063f2fde38b146101d55763fe4993ca14610190575f80fd5b346101d257806003193601126101d2576101ce6040516101ba816101b38161289d565b0382612b85565b604051918291602083526020830190612808565b0390f35b80fd5b50346101d25760203660031901126101d2576004356001600160a01b0381169081900361025f5761020461300b565b801561024b5781546001600160a01b03198116821783556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b631e4fbdf760e01b82526004829052602482fd5b5080fd5b50346101d257806003193601126101d2576020600b54604051908152f35b50346101d257806003193601126101d25761029a61300b565b6002600a54036107ec57604051635d52233760e01b8152606060048201526014805460648301819052908352829082906084600582901b83018101915f5160206135625f395f51905f5291859085015b8282106107b95750505050818103600319016024830152601254839161030f82612865565b91828252846001821691825f1461079857505060011461074b575b50506103428291600319838203016044840152612a79565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa91821561073f578080928194610715575b5083516001600160401b03811161056a5761039d601554612865565b601f81116106b7575b50602094601f8211600114610649578293949582916103d99492610495575b50508160011b915f199060031b1c19161790565b6015555b82516001600160401b03811161056a576103f8600c54612865565b601f81116105eb575b506020601f821160011461057e57819083949561043194926104955750508160011b915f199060031b1c19161790565b600c555b81516001600160401b03811161056a57610450601154612865565b601f811161050c575b50602092601f82116001146104a0578293829161048994926104955750508160011b915f199060031b1c19161790565b6011555b6003600a5580f35b015190505f806103c5565b60118352601f198216935f5160206134e25f395f51905f5291845b8681106104f457508360019596106104dc575b505050811b0160115561048d565b01515f1960f88460031b161c191690555f80806104ce565b919260206001819286850151815501940192016104bb565b61054f9060118452601f830160051c5f5160206134e25f395f51905f52019060208410610555575b601f0160051c5f5160206134e25f395f51905f520190612da8565b5f610459565b5f5160206134e25f395f51905f529150610534565b634e487b7160e01b82526041600452602482fd5b600c83525f5160206134625f395f51905f5290601f198316845b8181106105d3575095836001959697106105bb575b505050811b01600c55610435565b01515f1960f88460031b161c191690555f80806105ad565b9192602060018192868b015181550194019201610598565b61062e90600c8452601f830160051c5f5160206134625f395f51905f52019060208410610634575b601f0160051c5f5160206134625f395f51905f520190612da8565b5f610401565b5f5160206134625f395f51905f529150610613565b60158352601f198216955f5160206134c25f395f51905f5291845b88811061069f57508360019596979810610687575b505050811b016015556103dd565b01515f1960f88460031b161c191690555f8080610679565b91926020600181928685015181550194019201610664565b6106fa9060158452601f830160051c5f5160206134c25f395f51905f52019060208410610700575b601f0160051c5f5160206134c25f395f51905f520190612da8565b5f6103a6565b5f5160206134c25f395f51905f5291506106df565b9150925061073591503d8084833e61072d8183612b85565b810190612f28565b919091925f610381565b604051903d90823e3d90fd5b60128552849250905f5160206134a25f395f51905f525b818410610778575050016020016103428261032a565b600191955060209294508054838587010152019101909185938593610762565b61034294919550602093925060ff191683830152151560051b01019161032a565b6083198886030181529295509093509160019060209082906107db9088612ae9565b9601920192019285938795936102ea565b635e1e452d60e11b8152600490fd5b50346101d257806003193601126101d257602060405160858152f35b50346101d257806003193601126101d25760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b50346101d257806003193601126101d2576020600354604051908152f35b50346101d257806003193601126101d2576020601354604051908152f35b50346101d257806003193601126101d257602060405160308152f35b50346101d2576108b9366127b9565b91906108c361300b565b6108ce368483612c28565b60208151910120808352600160205260ff60408420541615610a3e576003548114610a2f57808352600160205260408320805460ff19169055600254835b818103610955575b847f3e65f48ee2200e8283ed2ed80108f37f59973b9d0613c1e3b8258ebd1d8de40c858861094f604051928392602084526020840191612dbe565b0390a180f35b8261095f82612bca565b90549060031b1c146109735760010161090c565b91505f198101908111610a1b5761099c61098f6109b292612bca565b90549060031b1c92612bca565b819391549060031b91821b915f19901b19161790565b9055600254928315610a07577f3e65f48ee2200e8283ed2ed80108f37f59973b9d0613c1e3b8258ebd1d8de40c92935f19016109ed81612bca565b8154905f199060031b1b1916905560025592915f80610914565b634e487b7160e01b83526031600452602483fd5b634e487b7160e01b84526011600452602484fd5b637809fd1360e11b8352600483fd5b634d3183db60e01b8352600483fd5b50346101d257806003193601126101d257602060405162093a808152f35b50346101d257806003193601126101d257546040516001600160a01b039091168152602090f35b50346101d257610aa1366127b9565b610aac92919261300b565b610ab46132bd565b610abf368285612c28565b6020815191012080600d558252600960205260ff6040832054166114e257604051636a107b3560e11b815292829184918291610b21917f0000000000000000000000000000000000000000000000000000000000000000919060048501612dde565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156114d5578190829383918491859186916113d2575b508051906001600160401b038211610e9e57610b84601054612865565b601f8111611374575b50602090601f831160011461130557610bbc92918891836104955750508160011b915f199060031b1c19161790565b6010555b8051906001600160401b0382116112f157610bdc600f54612865565b601f8111611293575b50602090601f831160011461122457610c1492918791836104955750508160011b915f199060031b1c19161790565b600f555b8051906001600160401b03821161121057610c34600e54612865565b601f81116111b2575b50602090601f831160011461114357610c6c92918691836104955750508160011b915f199060031b1c19161790565b600e555b805190600160401b821161112f57601354826013558083106110ec575b50602001601384525f5160206135025f395f51905f5284915b838310610fd4575050505082516001600160401b038111610ef557610ccc601254612865565b601f8111610f76575b506020601f8211600114610f095781908495610d059495926104955750508160011b915f199060031b1c19161790565b6012555b805190600160401b8211610ef55760145482601455808310610eb2575b50602001601483525f5160206135625f395f51905f5283915b838310610d8f5784610d87604051610d5a816101b381612929565b604051610d6a816101b381612999565b60405191610d8283610d7b81612a09565b0384612b85565b613042565b6001600a5580f35b80518051906001600160401b038211610e9e57610dac8454612865565b601f8111610e63575b50602090601f8311600114610dfb5792610dec836001959460209487968c926104955750508160011b915f199060031b1c19161790565b85555b01920192019190610d3f565b8488528188209190601f198416895b818110610e4b5750936020936001969387969383889510610e33575b505050811b018555610def565b01515f1960f88460031b161c191690555f8080610e26565b92936020600181928786015181550195019301610e0a565b610e8e9085895260208920601f850160051c81019160208610610e94575b601f0160051c0190612da8565b5f610db5565b9091508190610e81565b634e487b7160e01b87526041600452602487fd5b601484525f5160206135625f395f51905f5201825f5160206135625f395f51905f52015b818110610ee35750610d26565b80610eef600192612e34565b01610ed6565b634e487b7160e01b83526041600452602483fd5b601284525f5160206134a25f395f51905f5290601f198316855b818110610f5e57509583600195969710610f46575b505050811b01601255610d09565b01515f1960f88460031b161c191690555f8080610f38565b9192602060018192868b015181550194019201610f23565b610fb99060128552601f830160051c5f5160206134a25f395f51905f52019060208410610fbf575b601f0160051c5f5160206134a25f395f51905f520190612da8565b5f610cd5565b5f5160206134a25f395f51905f529150610f9e565b80518051906001600160401b0382116110d857610ff18454612865565b601f81116110a8575b50602090601f83116001146110405792611031836001959460209487968d926104955750508160011b915f199060031b1c19161790565b85555b01920192019190610ca6565b8489528189209190601f1984168a5b8181106110905750936020936001969387969383889510611078575b505050811b018555611034565b01515f1960f88460031b161c191690555f808061106b565b9293602060018192878601518155019501930161104f565b6110d290858a5260208a20601f850160051c81019160208610610e9457601f0160051c0190612da8565b5f610ffa565b634e487b7160e01b88526041600452602488fd5b601385525f5160206135025f395f51905f5201825f5160206135025f395f51905f52015b81811061111d5750610c8d565b80611129600192612e34565b01611110565b634e487b7160e01b84526041600452602484fd5b600e86525f5160206134825f395f51905f529190601f198416875b81811061119a5750908460019594939210611182575b505050811b01600e55610c70565b01515f1960f88460031b161c191690555f8080611174565b9293602060018192878601518155019501930161115e565b6111f590600e8752601f840160051c5f5160206134825f395f51905f520190602085106111fb575b601f0160051c5f5160206134825f395f51905f520190612da8565b5f610c3d565b5f5160206134825f395f51905f5291506111da565b634e487b7160e01b85526041600452602485fd5b600f87525f5160206135a25f395f51905f529190601f198416885b81811061127b5750908460019594939210611263575b505050811b01600f55610c18565b01515f1960f88460031b161c191690555f8080611255565b9293602060018192878601518155019501930161123f565b6112d690600f8852601f840160051c5f5160206135a25f395f51905f520190602085106112dc575b601f0160051c5f5160206135a25f395f51905f520190612da8565b5f610be5565b5f5160206135a25f395f51905f5291506112bb565b634e487b7160e01b86526041600452602486fd5b601088525f5160206135825f395f51905f529190601f198416895b81811061135c5750908460019594939210611344575b505050811b01601055610bc0565b01515f1960f88460031b161c191690555f8080611336565b92936020600181928786015181550195019301611320565b6113b79060108952601f840160051c5f5160206135825f395f51905f520190602085106113bd575b601f0160051c5f5160206135825f395f51905f520190612da8565b5f610b8d565b5f5160206135825f395f51905f52915061139c565b965050505050503d8082843e6113e88184612b85565b82019160c08184031261025f5780516001600160401b0381116114d15783611411918301612f8d565b9260208201516001600160401b0381116114cd5781611431918401612e83565b9160408101516001600160401b0381116114c95782611451918301612f8d565b9460608201516001600160401b0381116114c55783611471918401612e83565b9260808301516001600160401b0381116114c15781611491918501612e83565b9260a0810151906001600160401b0382116114bd576114b1929101612e83565b9093959291905f610b67565b8780fd5b8680fd5b8580fd5b8480fd5b8380fd5b8280fd5b50604051903d90823e3d90fd5b6305fb89e760e21b8252600482fd5b50346101d257611500366127b9565b9061150961300b565b611514368383612c28565b6020815191012091828452600960205260ff60408520541661164357604051632b5f2f8160e01b81529291849184918291611575917f0000000000000000000000000000000000000000000000000000000000000000919060048501612dde565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa8015611638576115f092849185918691611613575b5081836115c492613042565b60208151910151906bffffffffffffffffffffffff19821691601482106115f3575b505060601c6130f8565b80f35b6001600160601b031960149290920360031b82901b161690505f806115e6565b90506115c4925061162e91503d8087833e61072d8183612b85565b91925090816115b8565b6040513d85823e3d90fd5b6305fb89e760e21b8452600484fd5b50346101d257806003193601126101d25761166b61300b565b80546001600160a01b03198116825581906001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b50346101d25760403660031901126101d2576004356001600160401b03811161025f57610160600319823603011261025f576040519161016083018381106001600160401b0382111761056a5760405281600401356001600160401b038116810361025f5783526020830160248301358152604084016044840135815260608501906064850135825260848501356001600160401b0381116114c9576117589060043691880101612c7c565b956080810196875260a48601356001600160401b0381116114c5576117839060043691890101612c7c565b9560a0820196875260c48101356001600160401b0381116114c157810190366023830112156114c15760048201356117ba81612bf6565b926117c86040519485612b85565b81845260206004606082870194028301010190368211611ad657602401915b818310611ada5750505060c0830191825260e083019260e48201358452610100810190610104830135825261012483013592600d841015611ad6576101208201938452610144810135906001600160401b038211611ad2570136602382011215611ad65761185f903690602460048201359101612c28565b9361014082019485526024356001600160401b038111611ad257611887903690600401612c5e565b6007546001600160a01b03169c909b908d158015611ab7575b611aa85760208151516040516118cb816118bd8582019485612ec9565b03601f198101835282612b85565b519020915101516040516118e7816118bd60208201948561282c565b519020908c6020845151604051611905816118bd8582019485612ec9565b51902094510151604051611921816118bd60208201948561282c565b5190209451604051906020820192604083016020855282518091526020606085019301915b818110611a6b575050509061196d816001600160401b03949303601f198101835282612b85565b5190209551169b5199519a5198519551965197600d891015611a575751986040519b8c9b60208d019e8f5260408d015260608c015260808b015260a08a015260c089015260e08801526101008701526101208601526101408501526101608401526101808301526101a082016101a090526101c082016119ec91612808565b03601f19810182526119fe9082612b85565b5190207f19457468657265756d205369676e6564204d6573736167653a0a3332000000008252601c52603c902090611a3591613331565b611a419193929361336b565b6040516001600160a01b03909216148152602090f35b634e487b7160e01b8e52602160045260248efd5b825180516001600160a01b039081168652602082810151909116818701526040918201519186019190915260609094019390920191600101611946565b63f64b1d7b60e01b8c5260048cfd5b506085604051611aca816101b38161289d565b5114156118a0565b8a80fd5b8980fd5b606083360312611ad657604051611af081612b6a565b83356001600160a01b0381168103611b3957815260208401356001600160a01b0381168103611b39579181606093602080940152604086013560408201528152019201916117e7565b8b80fd5b50346101d257806003193601126101d257611b5661300b565b6001600a54036107ec57600b5460405163c13734e760e01b8152606060048201526013805460648301819052908452909291829184916084600583901b84018101925f5160206135025f395f51905f5291869086015b828210611d78575050505060248301528181036003190160448301528190611bd390612a79565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156114d5578192611d39575b5081516001600160401b03811161056a57611c2b600c54612865565b601f8111611cf1575b50602092601f8211600114611c855782938291611c6494926104955750508160011b915f199060031b1c19161790565b600c555b6001600b540180600b5560135414611c7d5780f35b6002600a5580f35b600c8352601f198216935f5160206134625f395f51905f5291845b868110611cd95750836001959610611cc1575b505050811b01600c55611c68565b01515f1960f88460031b161c191690555f8080611cb3565b91926020600181928685015181550194019201611ca0565b611d3390600c8452601f830160051c5f5160206134625f395f51905f5201906020841061063457601f0160051c5f5160206134625f395f51905f520190612da8565b5f611c34565b9091503d8083833e611d4b8183612b85565b81016020828203126114d15781516001600160401b0381116114cd57611d719201612e83565b905f611c0f565b6083198a8703018152929650929450926001906020908290611d9a9089612ae9565b97019201920192869593889593611bac565b50346101d257806003193601126101d257611dc5612dfb565b80516101ce60406020840151930151151560405193849384612ba6565b50346101d257806003193601126101d2576007546040516001600160a01b039091168152602090f35b50346101d25760203660031901126101d257600435906002548210156101d2576020611e3683612bca565b90549060031b1c604051908152f35b50346101d257806003193601126101d25760405190818160045492611e6984612865565b8084529360018116908115611ef35750600114611ea7575b50611e8e92500382612b85565b600554906101ce60ff6006541660405193849384612ba6565b60048152919290505f5160206135425f395f51905f525b818310611ed7575050906020611e8e928201015f611e81565b6020919350806001915483858801015201910190918392611ebe565b905060209250611e8e94915060ff191682840152151560051b8201015f611e81565b50346101d257806003193601126101d257611f2e61300b565b611f36612dfb565b6040810151156120595760208101805162093a808101809111610a1b57421161204a5781516020815191012090818452600160205260ff60408520541615611fca575b507fe4451b6eb971e43a03fd7fda060e219d0f94d4fcb6e73e481b1d9b86e5d1d22a9161094f91600355611fad6004612e34565b5f6005555f60065551604051918291602083526020830190612808565b51421061203b57808352600160205260408320600160ff19825416179055600254600160401b81101561112f57918161094f9261203161099c8660017fe4451b6eb971e43a03fd7fda060e219d0f94d4fcb6e73e481b1d9b86e5d1d22a9801600255612bca565b9055915091611f79565b63333bd2cb60e11b8352600483fd5b63dfafce0160e01b8352600483fd5b637085ae8360e11b8252600482fd5b50346101d257806003193601126101d25761208161300b565b612089612dfb565b6040810151156120595761094f7f89d9ae00acbd51f480eb5a0d3462ae199cb9e6076f950b9b2d64e4b18cf951d391611fad6004612e34565b5034612350575f366003190112612350576120db61300b565b6003600a5403612354576120f8604051610d5a816101b381612929565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156123505760405190632d77dd9360e11b825260606004830152815f60115461214c81612865565b908160648501526001811690815f1461232b57506001146122e4575b5081810360031901602483015261217e90612a79565b8181036003190160448301526015545f9161219882612865565b80825291600181169081156122c05750600114612273575b50509181805f9403915afa801561226857612255575b506115f06121d5600f54612865565b600f601f821161223d575b546001600160601b0319811691906014821061221d575b5050600d5490604051906122158261220e81612999565b0383612b85565b60601c6130f8565b6001600160601b031960149290920360031b82901b161690505f806121f7565b50600f83525f5160206135a25f395f51905f526121e0565b61226191505f90612b85565b5f5f6121c6565b6040513d5f823e3d90fd5b60155f9081525f5160206134c25f395f51905f52959350915b8083106122a257509193500160200181806121b0565b9350919360018160209254838587010152019101909391859361228c565b60ff191660208381019190915292151560051b9091019091019150829050806121b0565b60115f90815291505f5160206134e25f395f51905f525b818310612311575050810160840161217e612168565b8054838701608401528593506020909201916001016122fb565b61217e9350606491509160209260ff19166084860152151560051b8401010190612168565b5f80fd5b635e1e452d60e11b5f5260045ffd5b34612350575f366003190112612350576020600254604051908152f35b34612350575f366003190112612350576020600354604051908152f35b34612350575f3660031901126123505760405180602060025491828152019060025f527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace905f5b818110612407576101ce856123fb81870382612b85565b6040519182918261282c565b82548452602090930192600192830192016123e4565b34612350575f3660031901126123505760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b34612350576020366003190112612350576004355f526001602052602060ff60405f2054166040519015158152f35b34612350575f366003190112612350576007546040516001600160a01b039091168152602090f35b34612350575f366003190112612350576020600a54604051908152f35b34612350575f366003190112612350576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b34612350575f366003190112612350576101ce6040516101ba816101b38161289d565b3461235057612540366127b9565b61254861300b565b603081036127aa5760ff6006541680612791575b612782577f000000000000000000000000000000000000000000000000000000000000000042019182421161276e576040519261259884612b6a565b6125a3368484612c28565b938481526040602082019183835201906001825285516001600160401b03811161275a576125d2600454612865565b601f81116126fc575b506020601f8211600114612666578161265695949392612631927f8aac65bd065b693837bd35d5d39a8150c0bace29c9db8406bedabba66ca447c09a5f9261265b5750508160011b915f199060031b1c19161790565b6004555b5160055551151560ff80196006541691161760065560405193849384612dde565b0390a1005b015190508a806103c5565b601f1982169760045f525f5160206135425f395f51905f52985f5b8181106126e45750986001928492612656989796957f8aac65bd065b693837bd35d5d39a8150c0bace29c9db8406bedabba66ca447c09c106126cc575b505050811b01600455612635565b01515f1960f88460031b161c191690558980806126be565b838301518b556001909a019960209384019301612681565b61273f9060045f52601f830160051c5f5160206135425f395f51905f52019060208410612745575b601f0160051c5f5160206135425f395f51905f520190612da8565b876125db565b5f5160206135425f395f51905f529150612724565b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b63de0f6f7d60e01b5f5260045ffd5b5060055462093a80810180911161276e5742111561255c565b6309b9c8cd60e21b5f5260045ffd5b906020600319830112612350576004356001600160401b0381116123505782602382011215612350578060040135926001600160401b0384116123505760248483010111612350576024019190565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60206040818301928281528451809452019201905f5b81811061284f5750505090565b8251845260209384019390920191600101612842565b90600182811c92168015612893575b602083101461287f57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612874565b6008545f92916128ac82612865565b808252916001811690811561290d57506001146128c7575050565b60085f9081529293509091905f5160206135225f395f51905f525b8383106128f3575060209250010190565b6001816020929493945483858701015201910191906128e2565b9050602093945060ff929192191683830152151560051b010190565b6010545f929161293882612865565b808252916001811690811561290d5750600114612953575050565b60105f9081529293509091905f5160206135825f395f51905f525b83831061297f575060209250010190565b60018160209294939454838587010152019101919061296e565b600e545f92916129a882612865565b808252916001811690811561290d57506001146129c3575050565b600e5f9081529293509091905f5160206134825f395f51905f525b8383106129ef575060209250010190565b6001816020929493945483858701015201910191906129de565b600f545f9291612a1882612865565b808252916001811690811561290d5750600114612a33575050565b600f5f9081529293509091905f5160206135a25f395f51905f525b838310612a5f575060209250010190565b600181602092949394548385870101520191019190612a4e565b600c545f9291612a8882612865565b808252916001811690811561290d5750600114612aa3575050565b600c5f9081529293509091905f5160206134625f395f51905f525b838310612acf575060209250010190565b600181602092949394548385870101520191019190612abe565b5f9291815491612af883612865565b8083529260018116908115612b4d5750600114612b1457505050565b5f9081526020812093945091925b838310612b33575060209250010190565b600181602092949394548385870101520191019190612b22565b915050602093945060ff929192191683830152151560051b010190565b606081019081106001600160401b0382111761275a57604052565b90601f801991011681019081106001600160401b0382111761275a57604052565b919392612bbd604092606085526060850190612808565b9460208401521515910152565b600254811015612be25760025f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b6001600160401b03811161275a5760051b60200190565b6001600160401b03811161275a57601f01601f191660200190565b929192612c3482612c0d565b91612c426040519384612b85565b829481845281830111612350578281602093845f960137010152565b9080601f8301121561235057816020612c7993359101612c28565b90565b919060408382031261235057604051604081018181106001600160401b0382111761275a57604052809380356001600160401b03811161235057810183601f82011215612350578035612cce81612bf6565b91612cdc6040519384612b85565b81835260208084019260051b820101918683116123505760208201905b838210612d7b575050505082526020810135906001600160401b03821161235057019180601f84011215612350578235612d3281612bf6565b93612d406040519586612b85565b81855260208086019260051b82010192831161235057602001905b828210612d6b5750505060200152565b8135815260209182019101612d5b565b81356001600160401b03811161235057602091612d9d8a848094880101612c5e565b815201910190612cf9565b818110612db3575050565b5f8155600101612da8565b908060209392818452848401375f828201840152601f01601f1916010190565b939291602091612df691604087526040870191612dbe565b930152565b60405190612e0882612b6a565b81604051612e1b816101b3816004612ae9565b81526005546020820152604060ff600654161515910152565b612e3e8154612865565b9081612e48575050565b81601f5f9311600114612e5a5750555b565b81835260208320612e7691601f0160051c810190600101612da8565b8082528160208120915555565b81601f8201121561235057805190612e9a82612c0d565b92612ea86040519485612b85565b8284526020838301011161235057815f9260208093018386015e8301015290565b602081016020825282518091526040820191602060408360051b8301019401925f915b838310612efb57505050505090565b9091929394602080612f19600193603f198682030187528951612808565b97019301930191939290612eec565b916060838303126123505782516001600160401b0381116123505782612f4f918501612e83565b9260208101516001600160401b0381116123505783612f6f918301612e83565b9260408201516001600160401b03811161235057612c799201612e83565b9080601f83011215612350578151612fa481612bf6565b92612fb26040519485612b85565b81845260208085019260051b820101918383116123505760208201905b838210612fde57505050505090565b81516001600160401b0381116123505760209161300087848094880101612e83565b815201910190612fcf565b5f546001600160a01b0316330361301e57565b63118cdaa760e01b5f523360045260245ffd5b908151811015612be2570160200190565b9151601319016130e95751608419016130da5760348151106130935760405161306c606082612b85565b60308152602081019160403684375f5b603081036130a25750505190206003540361309357565b63719a9a4960e01b5f5260045ffd5b600481019081811161276e576001916001600160f81b0319906130c59085613031565b51165f1a6130d38286613031565b530161307c565b634ae601c160e11b5f5260045ffd5b6342f6e2a360e11b5f5260045ffd5b9092917f1882fd6e997ef1dc5e2691efe78564a5c540b36547363e8a7dcd454c61a944b660018060a01b03600754169260405193845260018060a01b03169283602082015260806040820152806131626131546080830161289d565b828103606084015288612808565b0390a15f52600960205260405f20600160ff198254161790556bffffffffffffffffffffffff60a01b600754161760075581516001600160401b03811161275a576131ae600854612865565b601f811161325f575b50602092601f82116001146131f3576131e7929382915f926104955750508160011b915f199060031b1c19161790565b6008555b612e586132bd565b601f1982169360085f525f5160206135225f395f51905f52915f5b868110613247575083600195961061322f575b505050811b016008556131eb565b01515f1960f88460031b161c191690555f8080613221565b9192602060018192868501518155019401920161320e565b6132a29060085f52601f830160051c5f5160206135225f395f51905f520190602084106132a8575b601f0160051c5f5160206135225f395f51905f520190612da8565b5f6131b7565b5f5160206135225f395f51905f529150613287565b5f600a555f600b555f6040516132d4602082612b85565b526132e0600c54612865565b601f81116132f0575b505f600c55565b601f90600c5f520160051c5f5160206134625f395f51905f52015f5160206134625f395f51905f525b81811061332657506132e9565b5f8155600101613319565b81519190604183036133615761335a9250602082015190606060408401519301515f1a906133df565b9192909190565b50505f9160029190565b60048110156133cb578061337d575050565b600181036133945763f645eedf60e01b5f5260045ffd5b600281036133af575063fce698f760e01b5f5260045260245ffd5b6003146133b95750565b6335e2f38360e21b5f5260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411613456579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15612268575f516001600160a01b0381161561344c57905f905f90565b505f906001905f90565b5050505f916003919056fedf6966c971051c3d54ec59162606531493a51404a002842f56009d7e5cf4a8c7bb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fdbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec344455f448fdea98c4d29eb340757ef0a66cd03dbb9538908a6a81d96026b71ec47531ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c6866de8ffda797e3de9c05e8fc57b3bf0ec28a930d40b0d285d93c06501cf6a090f3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee38a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19bce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec1b6847dc741a1b0cd08d278845f9d819d87b734759afb55fe2de5cb82a9ae6728d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac802a2646970667358221220e57bdc3920189aa9107b34aad9683c2b407748ae71c6a174b66937568ddd762964736f6c634300081e0033",
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

// PackGetAcceptedPcr0List is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14758b51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getAcceptedPcr0List() view returns(bytes32[])
func (teeAuthenticator *TeeAuthenticator) PackGetAcceptedPcr0List() []byte {
	enc, err := teeAuthenticator.abi.Pack("getAcceptedPcr0List")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetAcceptedPcr0List is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14758b51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getAcceptedPcr0List() view returns(bytes32[])
func (teeAuthenticator *TeeAuthenticator) TryPackGetAcceptedPcr0List() ([]byte, error) {
	return teeAuthenticator.abi.Pack("getAcceptedPcr0List")
}

// UnpackGetAcceptedPcr0List is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x14758b51.
//
// Solidity: function getAcceptedPcr0List() view returns(bytes32[])
func (teeAuthenticator *TeeAuthenticator) UnpackGetAcceptedPcr0List(data []byte) ([][32]byte, error) {
	out, err := teeAuthenticator.abi.Unpack("getAcceptedPcr0List", data)
	if err != nil {
		return *new([][32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	return out0, nil
}

// PackGetActiveImage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x187cd3fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getActiveImage() view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) PackGetActiveImage() []byte {
	enc, err := teeAuthenticator.abi.Pack("getActiveImage")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetActiveImage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x187cd3fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getActiveImage() view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) TryPackGetActiveImage() ([]byte, error) {
	return teeAuthenticator.abi.Pack("getActiveImage")
}

// UnpackGetActiveImage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x187cd3fb.
//
// Solidity: function getActiveImage() view returns(bytes32)
func (teeAuthenticator *TeeAuthenticator) UnpackGetActiveImage(data []byte) ([32]byte, error) {
	out, err := teeAuthenticator.abi.Unpack("getActiveImage", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetPendingSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5395c3ce.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingSwap() view returns(bytes targetPcr0, uint256 eta, bool pending)
func (teeAuthenticator *TeeAuthenticator) PackGetPendingSwap() []byte {
	enc, err := teeAuthenticator.abi.Pack("getPendingSwap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5395c3ce.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingSwap() view returns(bytes targetPcr0, uint256 eta, bool pending)
func (teeAuthenticator *TeeAuthenticator) TryPackGetPendingSwap() ([]byte, error) {
	return teeAuthenticator.abi.Pack("getPendingSwap")
}

// GetPendingSwapOutput serves as a container for the return parameters of contract
// method GetPendingSwap.
type GetPendingSwapOutput struct {
	TargetPcr0 []byte
	Eta        *big.Int
	Pending    bool
}

// UnpackGetPendingSwap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5395c3ce.
//
// Solidity: function getPendingSwap() view returns(bytes targetPcr0, uint256 eta, bool pending)
func (teeAuthenticator *TeeAuthenticator) UnpackGetPendingSwap(data []byte) (GetPendingSwapOutput, error) {
	out, err := teeAuthenticator.abi.Unpack("getPendingSwap", data)
	outstruct := new(GetPendingSwapOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.TargetPcr0 = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Eta = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Pending = *abi.ConvertType(out[2], new(bool)).(*bool)
	return *outstruct, nil
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
