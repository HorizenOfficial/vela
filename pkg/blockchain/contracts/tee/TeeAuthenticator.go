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

// AbstractTeeAuthenticatorMetaData contains all meta data concerning the AbstractTeeAuthenticator contract.
var AbstractTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0x64c062d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) PackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) []byte {
	enc, err := abstractTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refundAmount, applicationFee, signature)
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
func (abstractTeeAuthenticator *AbstractTeeAuthenticator) TryPackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) ([]byte, error) {
	return abstractTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refundAmount, applicationFee, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x64c062d1.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
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
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea26469706673582212201e9ae85ee52a8561505ac0aad3693753806311f16999397044bb40e552f6a7fa64736f6c634300081e0033",
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"contractINitroProver\",\"name\":\"_nitroProver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pcr0\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_maxVerificationAge\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPCR\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStep\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"oldPcr0\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"PcrZeroUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentUpdateStep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStep2TotalLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxVerificationAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nitroProver\",\"outputs\":[{\"internalType\":\"contractINitroProver\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pcr0\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step2CurrentIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"updatePcr0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTeeStep1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a8d895d4d49ffcfa5f7559dd61ad47cda6",
	Bin: "0x60c0604052346100765761001d610014610260565b929190916104d8565b61002561007b565b6132336106a18239608051818181610231015281816112690152818161176301528181611f20015281816125fe01526127f5015260a0518181816103e701528181611f540152612632015261323390f35b610081565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100ad90610085565b810190811060018060401b038211176100c557604052565b61008f565b906100dd6100d661007b565b92836100a3565b565b5f80fd5b5f80fd5b60018060a01b031690565b6100fb906100e7565b90565b610107816100f2565b0361010e57565b5f80fd5b9050519061011f826100fe565b565b61012a906100f2565b90565b61013681610121565b0361013d57565b5f80fd5b9050519061014e8261012d565b565b5f80fd5b5f80fd5b60018060401b03811161017457610170602091610085565b0190565b61008f565b90825f9392825e0152565b9092919261019961019482610158565b6100ca565b938185526020850190828401116101b5576101b392610179565b565b610154565b9080601f830112156101d8578160206101d593519101610184565b90565b610150565b90565b6101e9816101dd565b036101f057565b5f80fd5b90505190610201826101e0565b565b60808183031261025b57610219825f8301610112565b926102278360208401610141565b9260408301519060018060401b0382116102565761024a816102539386016101ba565b936060016101f4565b90565b6100e3565b6100df565b61027e6138d480380380610273816100ca565b928339810190610203565b90919293565b5190565b634e487b7160e01b5f52602260045260245ffd5b90600160028304921680156102bc575b60208310146102b757565b610288565b91607f16916102ac565b5f5260205f2090565b601f602091010490565b1b90565b919060086102f89102916102f25f19846102d9565b926102d9565b9181191691161790565b90565b61031961031461031e926101dd565b610302565b6101dd565b90565b90565b919061033a61033561034293610305565b610321565b9083546102dd565b9055565b5f90565b61035c91610356610346565b91610324565b565b5b81811061036a575050565b806103775f60019361034a565b0161035f565b9190601f811161038d575b505050565b6103996103be936102c6565b9060206103a5846102cf565b830193106103c6575b6103b7906102cf565b019061035e565b5f8080610388565b91506103b7819290506103ae565b1c90565b906103e8905f19906008026103d4565b191690565b816103f7916103d8565b906002021790565b9061040981610284565b9060018060401b0382116104c75761042b82610425855461029c565b8561037d565b602090601f831160011461045f5791809161044e935f92610453575b50506103ed565b90555b565b90915001515f80610447565b601f1983169161046e856102c6565b925f5b8181106104af57509160029391856001969410610495575b50505002019055610451565b6104a5910151601f8416906103d8565b90555f8080610489565b91936020600181928787015181550195019201610471565b61008f565b906104d6916103ff565b565b916104e66104ed9293610542565b60016104cc565b60805260a052565b90565b61050c610507610511926104f5565b610302565b6100e7565b90565b61051d906104f8565b90565b610529906100f2565b9052565b9190610540905f60208501940190610520565b565b8061055d6105576105525f610514565b6100f2565b916100f2565b1461056d5761056b90610641565b565b6105906105795f610514565b5f918291631e4fbdf760e01b83526004830161052d565b0390fd5b5f1c90565b60018060a01b031690565b6105b06105b591610594565b610599565b90565b6105c290546105a4565b90565b5f1b90565b906105db60018060a01b03916105c5565b9181191691161790565b6105f96105f46105fe926100e7565b610302565b6100e7565b90565b61060a906105e5565b90565b61061690610601565b90565b90565b9061063161062c6106389261060d565b610619565b82546105ca565b9055565b5f0190565b61064a5f6105b8565b610654825f61061c565b906106886106827f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09361060d565b9161060d565b9161069161007b565b8061069b8161063c565b0390a356fe60806040526004361015610013575b610fc9565b61001d5f3561016c565b8063081bec7e146101675780630a83fb91146101625780630b1bfab91461015d5780630dd7ce2f14610158578063127379db1461015357806334010ccf1461014e57806343f855c3146101495780635520916c1461014457806364c062d11461013f578063715018a61461013a5780637637d58a146101355780637ed2d7fc1461013057806381a9d38a1461012b5780638890b84e146101265780638da5cb5b14610121578063c3b331221461011c578063c91496c614610117578063d22a29d114610112578063db6a7b711461010d578063f2fde38b146101085763fe4993ca0361000e57610f94565b610f52565b610edc565b610e9a565b610e65565b610df9565b610dc4565b610d90565b610d5b565b610be3565b610b7a565b610ace565b610a8c565b6104eb565b6104b6565b610443565b610409565b6103b0565b61034d565b6102b7565b6101fa565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f91031261018a57565b61017c565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b6101d06101d96020936101de936101c78161018f565b93848093610193565b9586910161019c565b6101a7565b0190565b6101f79160208201915f8184039101526101b1565b90565b3461022a5761020a366004610180565b610226610215610fde565b61021d610172565b918291826101e2565b0390f35b610178565b7f000000000000000000000000000000000000000000000000000000000000000090565b60018060a01b031690565b90565b61027561027061027a92610253565b61025e565b610253565b90565b61028690610261565b90565b6102929061027d565b90565b61029e90610289565b9052565b91906102b5905f60208501940190610295565b565b346102e7576102c7366004610180565b6102e36102d261022f565b6102da610172565b918291826102a2565b0390f35b610178565b1c90565b90565b61030390600861030893026102ec565b6102f0565b90565b9061031691546102f3565b90565b61032560045f9061030b565b90565b90565b61033490610328565b9052565b919061034b905f6020850194019061032b565b565b3461037d5761035d366004610180565b610379610368610319565b610370610172565b91829182610338565b0390f35b610178565b61038b90610253565b90565b61039790610382565b9052565b91906103ae905f6020850194019061038e565b565b346103e0576103c0366004610180565b6103dc6103cb61101e565b6103d3610172565b9182918261039b565b0390f35b610178565b7f000000000000000000000000000000000000000000000000000000000000000090565b3461043957610419366004610180565b6104356104246103e5565b61042c610172565b91829182610338565b0390f35b610178565b5f0190565b3461047157610453366004610180565b61045b611351565b610463610172565b8061046d8161043e565b0390f35b610178565b60018060a01b031690565b61049190600861049693026102ec565b610476565b90565b906104a49154610481565b90565b6104b360025f90610499565b90565b346104e6576104c6366004610180565b6104e26104d16104a7565b6104d9610172565b9182918261039b565b0390f35b610178565b34610519576104fb366004610180565b610503611885565b61050b610172565b806105158161043e565b0390f35b610178565b5f80fd5b67ffffffffffffffff1690565b61053881610522565b0361053f57565b5f80fd5b905035906105508261052f565b565b90565b61055e81610552565b0361056557565b5f80fd5b9050359061057682610555565b565b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b9061059a906101a7565b810190811067ffffffffffffffff8211176105b457604052565b61057c565b906105cc6105c5610172565b9283610590565b565b67ffffffffffffffff81116105e65760208091020190565b61057c565b5f80fd5b5f80fd5b67ffffffffffffffff81116106115761060d6020916101a7565b0190565b61057c565b90825f939282370152565b90929192610636610631826105f3565b6105b9565b938185526020850190828401116106525761065092610616565b565b6105ef565b9080601f830112156106755781602061067293359101610621565b90565b610578565b92919061068e610689826105ce565b6105b9565b93818552602080860192028101918383116106e55781905b8382106106b4575050505050565b813567ffffffffffffffff81116106e0576020916106d58784938701610657565b8152019101906106a6565b610578565b6105eb565b9080601f83011215610708578160206107059335910161067a565b90565b610578565b67ffffffffffffffff81116107255760208091020190565b61057c565b67ffffffffffffffff8111610748576107446020916101a7565b0190565b61057c565b9092919261076261075d8261072a565b6105b9565b9381855260208501908284011161077e5761077c92610616565b565b6105ef565b9080601f830112156107a15781602061079e9335910161074d565b90565b610578565b9291906107ba6107b58261070d565b6105b9565b93818552602080860192028101918383116108115781905b8382106107e0575050505050565b813567ffffffffffffffff811161080c576020916108018784938701610783565b8152019101906107d2565b610578565b6105eb565b9080601f8301121561083457816020610831933591016107a6565b90565b610578565b67ffffffffffffffff81116108515760208091020190565b61057c565b5f80fd5b61086390610253565b90565b61086f8161085a565b0361087657565b5f80fd5b9050359061088782610866565b565b61089281610328565b0361089957565b5f80fd5b905035906108aa82610889565b565b91906040838203126108e6576108df906108c660406105b9565b936108d3825f830161087a565b5f86015260200161089d565b6020830152565b610856565b909291926109006108fb82610839565b6105b9565b93818552604060208601920283019281841161093f57915b8383106109255750505050565b602060409161093484866108ac565b815201920191610918565b6105eb565b9080601f830112156109625781602061095f933591016108eb565b90565b610578565b9061014082820312610a605761097f815f8401610543565b9261098d8260208501610569565b9261099b8360408301610569565b926109a98160608401610569565b92608083013567ffffffffffffffff8111610a5b57826109ca9185016106ea565b9260a081013567ffffffffffffffff8111610a5657836109eb918301610816565b9260c082013567ffffffffffffffff8111610a515781610a0c918401610944565b92610a1a8260e0850161089d565b92610a2983610100830161089d565b9261012082013567ffffffffffffffff8111610a4c57610a499201610657565b90565b61051e565b61051e565b61051e565b61051e565b61017c565b151590565b610a7390610a65565b9052565b9190610a8a905f60208501940190610a6a565b565b34610ac957610ac5610ab4610aa2366004610967565b98979097969196959295949394611be4565b610abc610172565b91829182610a77565b0390f35b610178565b34610afc57610ade366004610180565b610ae6611dc2565b610aee610172565b80610af88161043e565b0390f35b610178565b5f80fd5b909182601f83011215610b3f5781359167ffffffffffffffff8311610b3a576020019260018302840111610b3557565b6105eb565b610b01565b610578565b90602082820312610b75575f82013567ffffffffffffffff8111610b7057610b6c9201610b05565b9091565b61051e565b61017c565b34610ba957610b93610b8d366004610b44565b90612009565b610b9b610172565b80610ba58161043e565b0390f35b610178565b90602082820312610bde575f82013567ffffffffffffffff8111610bd957610bd69201610657565b90565b61051e565b61017c565b34610c1157610bfb610bf6366004610bae565b612183565b610c03610172565b80610c0d8161043e565b0390f35b610178565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610c5d575b6020831014610c5857565b610c29565b91607f1691610c4d565b60209181520190565b5f5260205f2090565b905f9291805490610c93610c8c83610c3d565b8094610c67565b916001811690815f14610cea5750600114610cae575b505050565b610cbb9192939450610c70565b915f925b818410610cd257505001905f8080610ca9565b60018160209295939554848601520191019290610cbf565b92949550505060ff19168252151560200201905f8080610ca9565b90610d0f91610c79565b90565b90610d32610d2b92610d22610172565b93848092610d05565b0383610590565b565b905f10610d4757610d4490610d12565b90565b610c16565b610d5860015f90610d34565b90565b34610d8b57610d6b366004610180565b610d87610d76610d4c565b610d7e610172565b918291826101e2565b0390f35b610178565b34610dbf57610da9610da3366004610b44565b90612748565b610db1610172565b80610dbb8161043e565b0390f35b610178565b34610df457610dd4366004610180565b610df0610ddf612754565b610de7610172565b9182918261039b565b0390f35b610178565b34610e2957610e09366004610180565b610e25610e14612769565b610e1c610172565b91829182610338565b0390f35b610178565b90565b610e45610e40610e4a92610e2e565b61025e565b610328565b90565b610e576085610e31565b90565b610e62610e4d565b90565b34610e9557610e75366004610180565b610e91610e80610e5a565b610e88610172565b91829182610338565b0390f35b610178565b34610ec857610eaa366004610180565b610eb26128e9565b610eba610172565b80610ec48161043e565b0390f35b610178565b610ed960055f9061030b565b90565b34610f0c57610eec366004610180565b610f08610ef7610ecd565b610eff610172565b91829182610338565b0390f35b610178565b610f1a81610382565b03610f2157565b5f80fd5b90503590610f3282610f11565b565b90602082820312610f4d57610f4a915f01610f25565b90565b61017c565b34610f8057610f6a610f65366004610f34565b612958565b610f72610172565b80610f7c8161043e565b0390f35b610178565b610f9160035f90610d34565b90565b34610fc457610fa4366004610180565b610fc0610faf610f85565b610fb7610172565b918291826101e2565b0390f35b610178565b5f80fd5b606090565b610fdb90610d12565b90565b610fe6610fcd565b50610ff16003610fd2565b90565b5f90565b5f1c90565b61100961100e91610ff8565b610476565b90565b61101b9054610ffd565b90565b611026610ff4565b506110316002611011565b90565b61103c612963565b611044611242565b565b61105261105791610ff8565b6102f0565b90565b6110649054611046565b90565b90565b61107e61107961108392611067565b61025e565b610328565b90565b5f80fd5b60e01b90565b5f91031261109a57565b61017c565b905f92918054906110b96110b283610c3d565b8094610193565b916001811690815f1461111057506001146110d4575b505050565b6110e19192939450610c70565b915f925b8184106110f857505001905f80806110cf565b600181602092959395548486015201910192906110e5565b92949550505060ff19168252151560200201905f80806110cf565b9161115790611149611165959360608601908682035f88015261109f565b90848203602086015261109f565b91604081840391015261109f565b90565b611170610172565b3d5f823e3d90fd5b6111829054610c3d565b90565b60601c90565b60601b90565b61119a9061118b565b90565b6111a96111ae91611185565b611191565b90565b6111bb905461119d565b90565b1b90565b6111df6111ce82611178565b9180601f8411611212575b506111b1565b90601481106111ed575b5090565b61120b906bffffffffffffffffffffffff19906014036008026111be565b165f6111e9565b61121c9150610c70565b5f6111d9565b61122e61123391611185565b610261565b90565b61123f90611222565b90565b61124c600461105a565b61125f611259600361106a565b91610328565b036113355761128d7f0000000000000000000000000000000000000000000000000000000000000000610289565b635aefbb26600a600692600e813b15611330575f936112c86112bd926112b1610172565b9788968795869561108a565b85526004850161112b565b03915afa801561132b576112ff575b506112fd6112ed6112e860086111c2565b611236565b6112f76007610fd2565b90612a42565b565b61131e905f3d8111611324575b6113168183610590565b810190611090565b5f6112d7565b503d61130c565b611168565b611086565b5f635e1e452d60e11b81528061134d6004820161043e565b0390fd5b611359611034565b565b611363612963565b61136b61173c565b565b90565b61138461137f6113899261136d565b61025e565b610328565b90565b909291926113a161139c826105f3565b6105b9565b938185526020850190828401116113bd576113bb9261019c565b565b6105ef565b9080601f830112156113e0578160206113dd9351910161138c565b90565b610578565b90602082820312611415575f82015167ffffffffffffffff81116114105761140d92016113c2565b90565b61051e565b61017c565b5490565b60209181520190565b5f5260205f2090565b60010190565b9061144a6114438361141a565b809261141e565b908161145b60208302840194611427565b925f915b83831061146e57505050505090565b9091929394602061148f611489838560019503875289610d05565b97611430565b930193019193929061145f565b6114c36114b86114d0959360608401908482035f860152611436565b93602083019061032b565b604081840391015261109f565b90565b601f602091010490565b919060086114f89102916114f25f19846111be565b926111be565b9181191691161790565b61151661151161151b92610328565b61025e565b610328565b90565b90565b919061153761153261153f93611502565b61151e565b9083546114dd565b9055565b5f90565b61155991611553611543565b91611521565b565b5b818110611567575050565b806115745f600193611547565b0161155c565b9190601f811161158a575b505050565b6115966115bb93610c70565b9060206115a2846114d3565b830193106115c3575b6115b4906114d3565b019061155b565b5f8080611585565b91506115b4819290506115ab565b906115e1905f19906008026102ec565b191690565b816115f0916115d1565b906002021790565b906116028161018f565b9067ffffffffffffffff82116116c257611626826116208554610c3d565b8561157a565b602090601f831160011461165a57918091611649935f9261164e575b50506115e6565b90555b565b90915001515f80611642565b601f1983169161166985610c70565b925f5b8181106116aa57509160029391856001969410611690575b5050500201905561164c565b6116a0910151601f8416906115d1565b90555f8080611684565b9193602060018192878701518155019501920161166c565b61057c565b906116d1916115f8565b565b60016116df9101610328565b90565b5f1b90565b906116f35f19916116e2565b9181191691161790565b9061171261170d61171992611502565b61151e565b82546116e7565b9055565b90565b61173461172f6117399261171d565b61025e565b610328565b90565b611746600461105a565b6117596117536001611370565b91610328565b03611869576117877f0000000000000000000000000000000000000000000000000000000000000000610289565b5f63c13734e791600c906117b861179e600561105a565b946117c360066117ac610172565b9788968795869561108a565b85526004850161149c565b03915afa8015611864576117e0915f91611842575b5060066116c7565b6117fc6117f56117f0600561105a565b6116d3565b60056116fd565b611806600561105a565b61182161181b611816600c61141a565b610328565b91610328565b14611829575b565b61183d6118366002611720565b60046116fd565b611827565b61185e91503d805f833e6118568183610590565b8101906113e5565b5f6117d8565b611168565b5f635e1e452d60e11b8152806118816004820161043e565b0390fd5b61188d61135b565b565b5f90565b90565b6118aa6118a56118af92611893565b61025e565b610253565b90565b6118bb90611896565b90565b5190565b60200190565b6118e76118f06020936118f5936118de8161018f565b93848093610c67565b9586910161019c565b6101a7565b0190565b90611903916118c8565b90565b60200190565b90611920611919836118be565b809261141e565b9081611931602083028401946118c2565b925f915b83831061194457505050505090565b90919293946020611966611960838560019503875289516118f9565b97611906565b9301930191939290611935565b6119889160208201915f81840391015261190c565b90565b60200190565b5190565b60209181520190565b60200190565b5190565b60209181520190565b6119d06119d96020936119de936119c7816119a4565b938480936119a8565b9586910161019c565b6101a7565b0190565b906119ec916119b1565b90565b60200190565b90611a09611a0283611991565b8092611995565b9081611a1a6020830284019461199e565b925f915b838310611a2d57505050505090565b90919293946020611a4f611a49838560019503875289516119e2565b976119ef565b9301930191939290611a1e565b611a719160208201915f8184039101526119f5565b90565b5190565b60209181520190565b60200190565b611a909061085a565b9052565b611a9d90610328565b9052565b90602080611ac393611ab95f8201515f860190611a87565b0151910190611a94565b565b90611ad281604093611aa1565b0190565b60200190565b90611af9611af3611aec84611a74565b8093611a78565b92611a81565b905f5b818110611b095750505090565b909192611b22611b1c6001928651611ac5565b94611ad6565b9101919091611afc565b611b419160208201915f818403910152611adc565b90565b611b4d90610522565b9052565b611b5a90610552565b9052565b9694929099989795939161012088019a5f8901611b7a91611b44565b60208801611b8791611b51565b60408701611b9491611b51565b60608601611ba191611b51565b60808501611bae91611b51565b60a08401611bbb91611b51565b60c08301611bc891611b51565b60e08201611bd59161032b565b61010001611be29161032b565b565b9897909392919596611bf461188f565b50611bfd61101e565b611c17611c11611c0c5f6118b2565b610382565b91610382565b148015611d6c575b611d5057611d3399611d2e98611d1697611cc7611cd6611c90611c9f611c59611c68611d079c611c4d610172565b92839160208301611973565b60208201810382520382610590565b611c7a611c748261018f565b9161198b565b2094611c84610172565b92839160208301611a5c565b60208201810382520382610590565b611cb1611cab8261018f565b9161198b565b2093611cbb610172565b92839160208301611b2c565b60208201810382520382610590565b611ce8611ce28261018f565b9161198b565b2092979995909192939495611cfb610172565b9a8b9960208b01611b5e565b60208201810382520382610590565b611d28611d228261018f565b9161198b565b20612aae565b612ae4565b611d4c611d46611d4161101e565b610382565b91610382565b1490565b5f63f64b1d7b60e01b815280611d686004820161043e565b0390fd5b50611d7d611d78610fde565b61018f565b611d96611d90611d8b610e4d565b610328565b91610328565b1415611c1f565b611da5612963565b611dad611daf565b565b611dc0611dbb5f6118b2565b612b06565b565b611dca611d9d565b565b90611dde91611dd9612963565b611f18565b565b91606083830312611e5c575f83015167ffffffffffffffff8111611e575782611e0a9185016113c2565b92602081015167ffffffffffffffff8111611e525783611e2b9183016113c2565b92604082015167ffffffffffffffff8111611e4d57611e4a92016113c2565b90565b61051e565b61051e565b61051e565b61017c565b9190611e7b81611e7481611e8095610193565b8095610616565b6101a7565b0190565b939290611ea2602091611eaa9460408801918883035f8a0152611e61565b94019061032b565b565b6bffffffffffffffffffffffff191690565b611ec89051611eac565b90565b611ee5611ee0611eda8361018f565b9261198b565b611ebe565b9060148110611ef3575b5090565b611f11906bffffffffffffffffffffffff19906014036008026111be565b165f611eef565b905f90611f447f0000000000000000000000000000000000000000000000000000000000000000610289565b611f87632b5f2f81949294611f927f0000000000000000000000000000000000000000000000000000000000000000611f7b610172565b9788968795869561108a565b855260048501611e84565b03915afa90811561200457611fd2915f8080919092611fd4575b611fcd92935090611fc3611fc89294918590612c89565b611ecb565b611236565b612a42565b565b505050611fcd611fc8611ffb611fc3933d805f833e611ff38183610590565b810190611de0565b92509250611fac565b611168565b9061201391611dcc565b565b61202690612021612963565b61212a565b565b905090565b905f929180549061204761204083610c3d565b8094612028565b916001811690815f146120995750600114612062575b505050565b61206f9192939450610c70565b5f905b83821061208557505001905f808061205d565b600181602092548486015201910190612072565b92949550505060ff191682528015150201905f808061205d565b6120bc9161202d565b90565b6120d4906120cb610172565b918291826120b3565b03902090565b6120ff6120f6926020926120ed8161018f565b94858093612028565b9384910161019c565b0190565b61210c916120da565b90565b6121249061211b610172565b91829182612103565b03902090565b6121819060018161216461215e7f7cb84e5a76f21dbf9cffe637ecb554a3d4f4eeb2cf05cd38c8a0db9d34ed6541936120bf565b9161210f565b9161216d610172565b806121778161043e565b0390a360016116c7565b565b61218c90612015565b565b906121a09161219b612963565b6125ee565b565b9291906121b66121b1826105ce565b6105b9565b938185526020808601920281019183831161220d5781905b8382106121dc575050505050565b815167ffffffffffffffff8111612208576020916121fd87849387016113c2565b8152019101906121ce565b610578565b6105eb565b9080601f830112156122305781602061222d935191016121a2565b90565b610578565b9160c083830312612323575f83015167ffffffffffffffff811161231e578261225f918501612212565b92602081015167ffffffffffffffff811161231957836122809183016113c2565b92604082015167ffffffffffffffff811161231457816122a1918401612212565b92606083015167ffffffffffffffff811161230f57826122c29185016113c2565b92608081015167ffffffffffffffff811161230a57836122e39183016113c2565b9260a082015167ffffffffffffffff81116123055761230292016113c2565b90565b61051e565b61051e565b61051e565b61051e565b61051e565b61051e565b61017c565b634e487b7160e01b5f52601160045260245ffd5b600190818003010490565b9061235a905f19906020036008026102ec565b8154169055565b905f9161237861237082610c70565b9283546115e6565b905555565b919290602082105f146123d657601f84116001146123a6576123a09293506115e6565b90555b5b565b50906123cc6123d19360016123c36123bd85610c70565b926114d3565b8201910161155b565b612361565b6123a3565b5061240d82936123e7600194610c70565b6124066123f3856114d3565b820192601f861680612418575b506114d3565b019061155b565b6002021790556123a4565b61242490888603612347565b5f612400565b92909168010000000000000000821161248a576020115f1461247b57602081105f1461245f57612459916115e6565b90555b5b565b60019160ff191661246f84610c70565b5560020201905561245c565b6001915060020201905561245d565b61057c565b90815461249b81610c3d565b908183116124c4575b8183106124b2575b50505050565b6124bb9361237d565b5f8080806124ac565b6124d08383838761242a565b6124a4565b5f6124df9161248f565b565b905f036124f3576124f1906124d5565b565b610c16565b5b818110612504575050565b806125115f6001936124e1565b016124f9565b9091828110612526575b505050565b61254461253e61253861254f9561233c565b9261233c565b92611427565b9182019101906124f8565b5f8080612521565b90680100000000000000008111612580578161257561257e9361141a565b90828155612517565b565b61057c565b5190565b6125ae6125a8612598846118be565b936125a38585612557565b6118c2565b91611427565b5f915b8383106125be5750505050565b60016020826125d66125d08495612585565b866116c7565b019201920191906125b1565b906125ec91612589565b565b905f906125f9612e1d565b6126227f0000000000000000000000000000000000000000000000000000000000000000610289565b61266563d420f66a9492946126707f0000000000000000000000000000000000000000000000000000000000000000612659610172565b9788968795869561108a565b855260048501611e84565b03915afa908115612743576126cf915f80808080809192939495612704575b6126c8959650916126b36126ba926126ac6126c1969560096116c7565b60086116c7565b60076116c7565b600c6125e2565b600b6116c7565b600d6125e2565b6126ee60096126e86126e2600792610fd2565b91610fd2565b90612c89565b6127026126fb6001611370565b60046116fd565b565b5050505050506126c86126b36126c16126ac6127346126ba953d805f833e61272c8183610590565b810190612235565b9397509397509350935061268f565b611168565b906127529161218e565b565b61275c610ff4565b506127665f611011565b90565b612771611543565b5061277c600c61141a565b90565b612787612963565b61278f6127ce565b565b916127bd906127af6127cb959360608601908682035f880152611436565b90848203602086015261109f565b91604081840391015261109f565b90565b6127d8600461105a565b6127eb6127e56002611720565b91610328565b036128cd576128197f0000000000000000000000000000000000000000000000000000000000000000610289565b5f635d52233791600d90612842600b9461284d6006612836610172565b9788968795869561108a565b855260048501612791565b03915afa9081156128c857612885915f8080919261289b575b61287e92935061287790600e6116c7565b60066116c7565b600a6116c7565b612899612892600361106a565b60046116fd565b565b50505061287e6128bf612877923d805f833e6128b78183610590565b810190611de0565b90919250612866565b611168565b5f635e1e452d60e11b8152806128e56004820161043e565b0390fd5b6128f161277f565b565b612904906128ff612963565b612906565b565b8061292161291b6129165f6118b2565b610382565b91610382565b146129315761292f90612b06565b565b61295461293d5f6118b2565b5f918291631e4fbdf760e01b83526004830161039b565b0390fd5b612961906128f3565b565b61296b612754565b61298461297e612979612e57565b610382565b91610382565b0361298b57565b6129ad612996612e57565b5f91829163118cdaa760e01b83526004830161039b565b0390fd5b92906129e7916129da6129f596946129d060808801945f89019061038e565b602087019061038e565b848203604086015261109f565b9160608184039101526101b1565b90565b90612a0960018060a01b03916116e2565b9181191691161790565b612a1c9061027d565b90565b90565b90612a37612a32612a3e92612a13565b612a1f565b82546129f8565b9055565b90612a99612aa092612a546002611011565b81600391612a8f867f1882fd6e997ef1dc5e2691efe78564a5c540b36547363e8a7dcd454c61a944b694612a86610172565b948594856129b1565b0390a16002612a22565b60036116c7565b612aa8612e1d565b565b5f90565b612ab6612aaa565b507f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b612b0391612afa91612af4610ff4565b50612ea7565b90929192612fa4565b90565b612b0f5f611011565b612b19825f612a22565b90612b4d612b477f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093612a13565b91612a13565b91612b56610172565b80612b608161043e565b0390a3565b90565b612b7c612b77612b8192612b65565b61025e565b610328565b90565b612b93612b9991939293610328565b92610328565b8201809211612ba457565b612328565b634e487b7160e01b5f52603260045260245ffd5b90612bc78261018f565b811015612bd957600160209102010190565b612ba9565b60ff60f81b1690565b612bf19051612bde565b90565b5f5260205f2090565b906020612c0f818306601f0393612bf4565b91040191565b90612c1f82611178565b80821015612c4c576020115f14612c3c5760209006601f0390915b565b612c4591612bfd565b9091612c3a565b612ba9565b60f81b90565b612c6090612c51565b90565b612c73906008612c7893026102ec565b612c57565b90565b90612c869154612c63565b90565b90612c939061018f565b612cac612ca6612ca1610e4d565b610328565b91610328565b03612db757612cba8161018f565b612ce8612ce2612cdd6004612cd8612cd26001611178565b91612b68565b612b84565b610328565b91610328565b10612d9b57612cf5611543565b91612d006001611178565b925b80612d15612d0f86610328565b91610328565b14612d9557612d40612d3b84612d3584612d2f6004612b68565b90612b84565b90612bbd565b612be7565b612d66612d60612d5b612d5560018690612c15565b90612c7b565b612bde565b91612bde565b03612d7957612d74906116d3565b612d02565b5f63719a9a4960e01b815280612d916004820161043e565b0390fd5b50915050565b5f63719a9a4960e01b815280612db36004820161043e565b0390fd5b5f634ae601c160e11b815280612dcf6004820161043e565b0390fd5b612de7612de2612dec92611893565b61025e565b610328565b90565b90612e01612dfc8361072a565b6105b9565b918252565b612e0f5f612def565b90565b612e1a612e06565b90565b612e30612e295f612dd3565b60046116fd565b612e43612e3c5f612dd3565b60056116fd565b612e55612e4e612e12565b60066116c7565b565b612e5f610ff4565b503390565b5f90565b90565b612e7f612e7a612e8492612e68565b61025e565b610328565b90565b612e9b612e96612ea092610328565b6116e2565b610552565b90565b5f90565b919091612eb2610ff4565b50612ebb612e64565b50612ec4612aaa565b50612ece8361018f565b612ee1612edb6041612e6b565b91610328565b145f14612f2857612f219192612ef5612aaa565b50612efe612aaa565b50612f07612ea3565b506020810151606060408301519201515f1a909192613101565b9192909190565b50612f325f6118b2565b90612f46612f4160029461018f565b612e87565b91929190565b634e487b7160e01b5f52602160045260245ffd5b60041115612f6a57565b612f4c565b90612f7982612f60565b565b9190612f8e905f60208501940190611b51565b565b612f9c612fa191610ff8565b611502565b90565b80612fb7612fb15f612f6f565b91612f6f565b145f14612fc2575050565b80612fd6612fd06001612f6f565b91612f6f565b145f14612ff9575f63f645eedf60e01b815280612ff56004820161043e565b0390fd5b8061300d6130076002612f6f565b91612f6f565b145f1461303b5761303761302083612f90565b5f91829163fce698f760e01b835260048301610338565b0390fd5b61304e6130486003612f6f565b91612f6f565b146130565750565b613071905f9182916335e2f38360e21b835260048301612f7b565b0390fd5b90565b61308c61308761309192613075565b61025e565b610328565b90565b60ff1690565b6130a390613094565b9052565b6130dc6130e3946130d26060949897956130c8608086019a5f870190611b51565b602085019061309a565b6040830190611b51565b0190611b51565b565b6130f96130f46130fe92611893565b6116e2565b610552565b90565b93929361310c610ff4565b50613115612e64565b5061311e612aaa565b5061312885612f90565b61315a6131547f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0613078565b91610328565b116131e7579061317d602094955f94939293613174610172565b948594856130a7565b838052039060015afa156131e2576131955f516116e2565b806131b06131aa6131a55f6118b2565b610382565b91610382565b146131c6575f916131c05f6130e5565b91929190565b506131d05f6118b2565b6001916131dc5f6130e5565b91929190565b611168565b5050506131f35f6118b2565b906003929192919056fea26469706673582212207df0ef3c21c2908f4e43f586e156c9030a6ccdc1389f67a52aaf5cd891b3071464736f6c634300081e0033",
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
// the contract method with ID 0x64c062d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) PackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refundAmount, applicationFee, signature)
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
func (teeAuthenticator *TeeAuthenticator) TryPackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refundAmount, applicationFee, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x64c062d1.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
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
