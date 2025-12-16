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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"decodeAttestation\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxAge\",\"type\":\"uint256\"}],\"name\":\"verifyAttestation\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"attestation_decoded\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"max_age\",\"type\":\"uint256\"}],\"name\":\"verifyAttestationFromDecoded\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"max_age\",\"type\":\"uint256\"}],\"name\":\"verifyAttestationStep1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestationSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pubKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"buf\",\"type\":\"bytes\"}],\"name\":\"verifyAttestationStep2\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// PackDecodeAttestation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aee878f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decodeAttestation(bytes attestation) view returns(bytes[])
func (iNitroProver *INitroProver) PackDecodeAttestation(attestation []byte) []byte {
	enc, err := iNitroProver.abi.Pack("decodeAttestation", attestation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecodeAttestation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aee878f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decodeAttestation(bytes attestation) view returns(bytes[])
func (iNitroProver *INitroProver) TryPackDecodeAttestation(attestation []byte) ([]byte, error) {
	return iNitroProver.abi.Pack("decodeAttestation", attestation)
}

// UnpackDecodeAttestation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0aee878f.
//
// Solidity: function decodeAttestation(bytes attestation) view returns(bytes[])
func (iNitroProver *INitroProver) UnpackDecodeAttestation(data []byte) ([][]byte, error) {
	out, err := iNitroProver.abi.Unpack("decodeAttestation", data)
	if err != nil {
		return *new([][]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)
	return out0, nil
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

// PackVerifyAttestationFromDecoded is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9370bb3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyAttestationFromDecoded(bytes[] attestation_decoded, uint256 max_age) view returns(bytes, bytes, bytes)
func (iNitroProver *INitroProver) PackVerifyAttestationFromDecoded(attestationDecoded [][]byte, maxAge *big.Int) []byte {
	enc, err := iNitroProver.abi.Pack("verifyAttestationFromDecoded", attestationDecoded, maxAge)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyAttestationFromDecoded is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9370bb3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyAttestationFromDecoded(bytes[] attestation_decoded, uint256 max_age) view returns(bytes, bytes, bytes)
func (iNitroProver *INitroProver) TryPackVerifyAttestationFromDecoded(attestationDecoded [][]byte, maxAge *big.Int) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestationFromDecoded", attestationDecoded, maxAge)
}

// VerifyAttestationFromDecodedOutput serves as a container for the return parameters of contract
// method VerifyAttestationFromDecoded.
type VerifyAttestationFromDecodedOutput struct {
	Arg0 []byte
	Arg1 []byte
	Arg2 []byte
}

// UnpackVerifyAttestationFromDecoded is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9370bb3e.
//
// Solidity: function verifyAttestationFromDecoded(bytes[] attestation_decoded, uint256 max_age) view returns(bytes, bytes, bytes)
func (iNitroProver *INitroProver) UnpackVerifyAttestationFromDecoded(data []byte) (VerifyAttestationFromDecodedOutput, error) {
	out, err := iNitroProver.abi.Unpack("verifyAttestationFromDecoded", data)
	outstruct := new(VerifyAttestationFromDecodedOutput)
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
// Solidity: function verifyAttestationStep1(bytes attestation, uint256 max_age) view returns(bytes, bytes, bytes, bytes, bytes, bytes)
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
// Solidity: function verifyAttestationStep1(bytes attestation, uint256 max_age) view returns(bytes, bytes, bytes, bytes, bytes, bytes)
func (iNitroProver *INitroProver) TryPackVerifyAttestationStep1(attestation []byte, maxAge *big.Int) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestationStep1", attestation, maxAge)
}

// VerifyAttestationStep1Output serves as a container for the return parameters of contract
// method VerifyAttestationStep1.
type VerifyAttestationStep1Output struct {
	Arg0 []byte
	Arg1 []byte
	Arg2 []byte
	Arg3 []byte
	Arg4 []byte
	Arg5 []byte
}

// UnpackVerifyAttestationStep1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd420f66a.
//
// Solidity: function verifyAttestationStep1(bytes attestation, uint256 max_age) view returns(bytes, bytes, bytes, bytes, bytes, bytes)
func (iNitroProver *INitroProver) UnpackVerifyAttestationStep1(data []byte) (VerifyAttestationStep1Output, error) {
	out, err := iNitroProver.abi.Unpack("verifyAttestationStep1", data)
	outstruct := new(VerifyAttestationStep1Output)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Arg2 = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.Arg3 = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.Arg4 = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.Arg5 = *abi.ConvertType(out[5], new([]byte)).(*[]byte)
	return *outstruct, nil
}

// PackVerifyAttestationStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc30ad028.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyAttestationStep2(bytes attestationSig, bytes pubKey, bytes buf) view returns()
func (iNitroProver *INitroProver) PackVerifyAttestationStep2(attestationSig []byte, pubKey []byte, buf []byte) []byte {
	enc, err := iNitroProver.abi.Pack("verifyAttestationStep2", attestationSig, pubKey, buf)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyAttestationStep2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc30ad028.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyAttestationStep2(bytes attestationSig, bytes pubKey, bytes buf) view returns()
func (iNitroProver *INitroProver) TryPackVerifyAttestationStep2(attestationSig []byte, pubKey []byte, buf []byte) ([]byte, error) {
	return iNitroProver.abi.Pack("verifyAttestationStep2", attestationSig, pubKey, buf)
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
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea26469706673582212203f90f9921696388699ff60a8b1145261f1bc5309ddbd48eb893b0e521a60a9b064736f6c634300081e0033",
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"contractINitroProver\",\"name\":\"_nitroProver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pcr0\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_maxVerificationAge\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPCR\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPKLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStep\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"oldPcr0\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"PcrZeroUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldPubSecp521r1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPubSecp521r1\",\"type\":\"bytes\"}],\"name\":\"TeeUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentUpdateStep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxVerificationAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nitroProver\",\"outputs\":[{\"internalType\":\"contractINitroProver\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pcr0\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"newPcr0\",\"type\":\"bytes\"}],\"name\":\"updatePcr0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"attestation_decoded\",\"type\":\"bytes[]\"}],\"name\":\"updateTeeDecoded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"updateTeeStep1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTeeStep2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "a8d895d4d49ffcfa5f7559dd61ad47cda6",
	Bin: "0x60c0604052346100765761001d610014610260565b929190916104d8565b61002561007b565b6129d46106a18239608051818181610201015281816114b201528181611866015281816119c00152611f76015260a0518181816103b70152818161189a015281816119f40152611faa01526129d490f35b610081565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100ad90610085565b810190811060018060401b038211176100c557604052565b61008f565b906100dd6100d661007b565b92836100a3565b565b5f80fd5b5f80fd5b60018060a01b031690565b6100fb906100e7565b90565b610107816100f2565b0361010e57565b5f80fd5b9050519061011f826100fe565b565b61012a906100f2565b90565b61013681610121565b0361013d57565b5f80fd5b9050519061014e8261012d565b565b5f80fd5b5f80fd5b60018060401b03811161017457610170602091610085565b0190565b61008f565b90825f9392825e0152565b9092919261019961019482610158565b6100ca565b938185526020850190828401116101b5576101b392610179565b565b610154565b9080601f830112156101d8578160206101d593519101610184565b90565b610150565b90565b6101e9816101dd565b036101f057565b5f80fd5b90505190610201826101e0565b565b60808183031261025b57610219825f8301610112565b926102278360208401610141565b9260408301519060018060401b0382116102565761024a816102539386016101ba565b936060016101f4565b90565b6100e3565b6100df565b61027e61307580380380610273816100ca565b928339810190610203565b90919293565b5190565b634e487b7160e01b5f52602260045260245ffd5b90600160028304921680156102bc575b60208310146102b757565b610288565b91607f16916102ac565b5f5260205f2090565b601f602091010490565b1b90565b919060086102f89102916102f25f19846102d9565b926102d9565b9181191691161790565b90565b61031961031461031e926101dd565b610302565b6101dd565b90565b90565b919061033a61033561034293610305565b610321565b9083546102dd565b9055565b5f90565b61035c91610356610346565b91610324565b565b5b81811061036a575050565b806103775f60019361034a565b0161035f565b9190601f811161038d575b505050565b6103996103be936102c6565b9060206103a5846102cf565b830193106103c6575b6103b7906102cf565b019061035e565b5f8080610388565b91506103b7819290506103ae565b1c90565b906103e8905f19906008026103d4565b191690565b816103f7916103d8565b906002021790565b9061040981610284565b9060018060401b0382116104c75761042b82610425855461029c565b8561037d565b602090601f831160011461045f5791809161044e935f92610453575b50506103ed565b90555b565b90915001515f80610447565b601f1983169161046e856102c6565b925f5b8181106104af57509160029391856001969410610495575b50505002019055610451565b6104a5910151601f8416906103d8565b90555f8080610489565b91936020600181928787015181550195019201610471565b61008f565b906104d6916103ff565b565b916104e66104ed9293610542565b60016104cc565b60805260a052565b90565b61050c610507610511926104f5565b610302565b6100e7565b90565b61051d906104f8565b90565b610529906100f2565b9052565b9190610540905f60208501940190610520565b565b8061055d6105576105525f610514565b6100f2565b916100f2565b1461056d5761056b90610641565b565b6105906105795f610514565b5f918291631e4fbdf760e01b83526004830161052d565b0390fd5b5f1c90565b60018060a01b031690565b6105b06105b591610594565b610599565b90565b6105c290546105a4565b90565b5f1b90565b906105db60018060a01b03916105c5565b9181191691161790565b6105f96105f46105fe926100e7565b610302565b6100e7565b90565b61060a906105e5565b90565b61061690610601565b90565b90565b9061063161062c6106389261060d565b610619565b82546105ca565b9055565b5f0190565b61064a5f6105b8565b610654825f61061c565b906106886106827f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09361060d565b9161060d565b9161069161007b565b8061069b8161063c565b0390a356fe60806040526004361015610013575b610e11565b61001d5f3561013c565b8063081bec7e146101375780630a83fb91146101325780630b1bfab91461012d5780630dd7ce2f14610128578063127379db1461012357806343f855c31461011e57806354b21558146101195780635520916c14610114578063715018a61461010f57806374ab01681461010a5780637637d58a146101055780637ed2d7fc1461010057806381a9d38a146100fb5780638890b84e146100f65780638da5cb5b146100f1578063c91496c6146100ec578063f2fde38b146100e75763fe4993ca0361000e57610ddc565b610d9a565b610d24565b610cb8565b610c84565b610c4f565b610ad7565b610a6e565b610a04565b61095c565b610929565b6108e2565b61044e565b6103d9565b610380565b61031d565b610287565b6101ca565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f91031261015a57565b61014c565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b6101a06101a96020936101ae936101978161015f565b93848093610163565b9586910161016c565b610177565b0190565b6101c79160208201915f818403910152610181565b90565b346101fa576101da366004610150565b6101f66101e5610e26565b6101ed610142565b918291826101b2565b0390f35b610148565b7f000000000000000000000000000000000000000000000000000000000000000090565b60018060a01b031690565b90565b61024561024061024a92610223565b61022e565b610223565b90565b61025690610231565b90565b6102629061024d565b90565b61026e90610259565b9052565b9190610285905f60208501940190610265565b565b346102b757610297366004610150565b6102b36102a26101ff565b6102aa610142565b91829182610272565b0390f35b610148565b1c90565b90565b6102d39060086102d893026102bc565b6102c0565b90565b906102e691546102c3565b90565b6102f560045f906102db565b90565b90565b610304906102f8565b9052565b919061031b905f602085019401906102fb565b565b3461034d5761032d366004610150565b6103496103386102e9565b610340610142565b91829182610308565b0390f35b610148565b61035b90610223565b90565b61036790610352565b9052565b919061037e905f6020850194019061035e565b565b346103b057610390366004610150565b6103ac61039b610e66565b6103a3610142565b9182918261036b565b0390f35b610148565b7f000000000000000000000000000000000000000000000000000000000000000090565b34610409576103e9366004610150565b6104056103f46103b5565b6103fc610142565b91829182610308565b0390f35b610148565b60018060a01b031690565b61042990600861042e93026102bc565b61040e565b90565b9061043c9154610419565b90565b61044b60025f90610431565b90565b3461047e5761045e366004610150565b61047a61046961043f565b610471610142565b9182918261036b565b0390f35b610148565b5f80fd5b67ffffffffffffffff1690565b61049d81610487565b036104a457565b5f80fd5b905035906104b582610494565b565b90565b6104c3816104b7565b036104ca57565b5f80fd5b905035906104db826104ba565b565b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b906104ff90610177565b810190811067ffffffffffffffff82111761051957604052565b6104e1565b9061053161052a610142565b92836104f5565b565b67ffffffffffffffff811161054b5760208091020190565b6104e1565b5f80fd5b5f80fd5b67ffffffffffffffff811161057657610572602091610177565b0190565b6104e1565b90825f939282370152565b9092919261059b61059682610558565b61051e565b938185526020850190828401116105b7576105b59261057b565b565b610554565b9080601f830112156105da578160206105d793359101610586565b90565b6104dd565b9291906105f36105ee82610533565b61051e565b938185526020808601920281019183831161064a5781905b838210610619575050505050565b813567ffffffffffffffff81116106455760209161063a87849387016105bc565b81520191019061060b565b6104dd565b610550565b9080601f8301121561066d5781602061066a933591016105df565b90565b6104dd565b67ffffffffffffffff811161068a5760208091020190565b6104e1565b5f80fd5b61069c90610223565b90565b6106a881610693565b036106af57565b5f80fd5b905035906106c08261069f565b565b6106cb816102f8565b036106d257565b5f80fd5b905035906106e3826106c2565b565b919060408382031261071f57610718906106ff604061051e565b9361070c825f83016106b3565b5f8601526020016106d6565b6020830152565b61068f565b9092919261073961073482610672565b61051e565b93818552604060208601920283019281841161077857915b83831061075e5750505050565b602060409161076d84866106e5565b815201920191610751565b610550565b9080601f8301121561079b5781602061079893359101610724565b90565b6104dd565b5f80fd5b909182601f830112156107de5781359167ffffffffffffffff83116107d95760200192600183028401116107d457565b610550565b6107a0565b6104dd565b91610120838303126108b6576107fb825f85016104a8565b9261080983602083016104ce565b9261081781604084016104ce565b9261082582606085016104ce565b92608081013567ffffffffffffffff81116108b1578361084691830161064f565b9260a082013567ffffffffffffffff81116108ac578161086791840161077d565b926108758260c085016106d6565b926108838360e083016106d6565b9261010082013567ffffffffffffffff81116108a7576108a392016107a4565b9091565b610483565b610483565b610483565b61014c565b151590565b6108c9906108bb565b9052565b91906108e0905f602085019401906108c0565b565b3461091f5761091b61090a6108f83660046107e3565b98979097969196959295949394611104565b610912610142565b918291826108cd565b0390f35b610148565b5f0190565b3461095757610939366004610150565b6109416115ac565b610949610142565b8061095381610924565b0390f35b610148565b3461098a5761096c366004610150565b6109746115db565b61097c610142565b8061098681610924565b0390f35b610148565b909182601f830112156109c95781359167ffffffffffffffff83116109c45760200192602083028401116109bf57565b610550565b6107a0565b6104dd565b906020828203126109ff575f82013567ffffffffffffffff81116109fa576109f6920161098f565b9091565b610483565b61014c565b34610a3357610a1d610a173660046109ce565b9061194d565b610a25610142565b80610a2f81610924565b0390f35b610148565b90602082820312610a69575f82013567ffffffffffffffff8111610a6457610a6092016107a4565b9091565b610483565b61014c565b34610a9d57610a87610a81366004610a38565b90611aa7565b610a8f610142565b80610a9981610924565b0390f35b610148565b90602082820312610ad2575f82013567ffffffffffffffff8111610acd57610aca92016105bc565b90565b610483565b61014c565b34610b0557610aef610aea366004610aa2565b611e21565b610af7610142565b80610b0181610924565b0390f35b610148565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610b51575b6020831014610b4c57565b610b1d565b91607f1691610b41565b60209181520190565b5f5260205f2090565b905f9291805490610b87610b8083610b31565b8094610b5b565b916001811690815f14610bde5750600114610ba2575b505050565b610baf9192939450610b64565b915f925b818410610bc657505001905f8080610b9d565b60018160209295939554848601520191019290610bb3565b92949550505060ff19168252151560200201905f8080610b9d565b90610c0391610b6d565b90565b90610c26610c1f92610c16610142565b93848092610bf9565b03836104f5565b565b905f10610c3b57610c3890610c06565b90565b610b0a565b610c4c60015f90610c28565b90565b34610c7f57610c5f366004610150565b610c7b610c6a610c40565b610c72610142565b918291826101b2565b0390f35b610148565b34610cb357610c9d610c97366004610a38565b906120a1565b610ca5610142565b80610caf81610924565b0390f35b610148565b34610ce857610cc8366004610150565b610ce4610cd36120ad565b610cdb610142565b9182918261036b565b0390f35b610148565b90565b610d04610cff610d0992610ced565b61022e565b6102f8565b90565b610d166085610cf0565b90565b610d21610d0c565b90565b34610d5457610d34366004610150565b610d50610d3f610d19565b610d47610142565b91829182610308565b0390f35b610148565b610d6281610352565b03610d6957565b5f80fd5b90503590610d7a82610d59565b565b90602082820312610d9557610d92915f01610d6d565b90565b61014c565b34610dc857610db2610dad366004610d7c565b612127565b610dba610142565b80610dc481610924565b0390f35b610148565b610dd960035f90610c28565b90565b34610e0c57610dec366004610150565b610e08610df7610dcd565b610dff610142565b918291826101b2565b0390f35b610148565b5f80fd5b606090565b610e2390610c06565b90565b610e2e610e15565b50610e396003610e1a565b90565b5f90565b5f1c90565b610e51610e5691610e40565b61040e565b90565b610e639054610e45565b90565b610e6e610e3c565b50610e796002610e59565b90565b5f90565b90565b610e97610e92610e9c92610e80565b61022e565b610223565b90565b610ea890610e83565b90565b610eb59054610b31565b90565b5190565b60209181520190565b60200190565b610eea610ef3602093610ef893610ee18161015f565b93848093610b5b565b9586910161016c565b610177565b0190565b90610f0691610ecb565b90565b60200190565b90610f23610f1c83610eb8565b8092610ebc565b9081610f3460208302840194610ec5565b925f915b838310610f4757505050505090565b90919293946020610f69610f6383856001950387528951610efc565b97610f09565b9301930191939290610f38565b610f8b9160208201915f818403910152610f0f565b90565b60200190565b5190565b60209181520190565b60200190565b610fb090610693565b9052565b610fbd906102f8565b9052565b90602080610fe393610fd95f8201515f860190610fa7565b0151910190610fb4565b565b90610ff281604093610fc1565b0190565b60200190565b9061101961101361100c84610f94565b8093610f98565b92610fa1565b905f5b8181106110295750505090565b90919261104261103c6001928651610fe5565b94610ff6565b910191909161101c565b6110619160208201915f818403910152610ffc565b90565b61106d90610487565b9052565b61107a906104b7565b9052565b959391989796949290986101008701995f880161109a91611064565b602087016110a791611071565b604086016110b491611071565b606085016110c191611071565b608084016110ce91611071565b60a083016110db91611071565b60c082016110e8916102fb565b60e0016110f4916102fb565b565b611101913691610586565b90565b94929895969893909193611116610e7c565b506111216002610e59565b61113b6111356111305f610e9f565b610352565b91610352565b14801561125f575b611243576112249961121e9861121897611200966111b26111c161117b61118a6111f19961116f610142565b92839160208301610f76565b602082018103825203826104f5565b61119c6111968261015f565b91610f8e565b20926111a6610142565b9283916020830161104c565b602082018103825203826104f5565b6111d36111cd8261015f565b91610f8e565b209196989490919293946111e5610142565b998a9860208a0161107e565b602082018103825203826104f5565b61121261120c8261015f565b91610f8e565b20612136565b926110f6565b9061216c565b61123f6112396112346002610e59565b610352565b91610352565b1490565b5f63f64b1d7b60e01b81528061125b60048201610924565b0390fd5b5061126a6003610eab565b61128361127d611278610d0c565b6102f8565b916102f8565b1415611143565b61129261218e565b61129a61148b565b565b6112a86112ad91610e40565b6102c0565b90565b6112ba905461129c565b90565b90565b6112d46112cf6112d9926112bd565b61022e565b6102f8565b90565b5f80fd5b60e01b90565b5f9103126112f057565b61014c565b905f929180549061130f61130883610b31565b8094610163565b916001811690815f14611366575060011461132a575b505050565b6113379192939450610b64565b915f925b81841061134e57505001905f8080611325565b6001816020929593955484860152019101929061133b565b92949550505060ff19168252151560200201905f8080611325565b916113ad9061139f6113bb959360608601908682035f8801526112f5565b9084820360208601526112f5565b9160408184039101526112f5565b90565b6113c6610142565b3d5f823e3d90fd5b60601c90565b60601b90565b6113e3906113d4565b90565b6113f26113f7916113ce565b6113da565b90565b61140490546113e6565b90565b1b90565b61142861141782610eab565b9180601f841161145b575b506113fa565b9060148110611436575b5090565b611454906bffffffffffffffffffffffff1990601403600802611407565b165f611432565b6114659150610b64565b5f611422565b61147761147c916113ce565b610231565b90565b6114889061146b565b90565b61149560046112b0565b6114a86114a260016112c0565b916102f8565b03611590576114d67f0000000000000000000000000000000000000000000000000000000000000000610259565b63c30ad0286008600992600a813b1561158b575f93611511611506926114fa610142565b978896879586956112e0565b855260048501611381565b03915afa80156115865761155a575b5061153361152e6007610e1a565b612330565b611558611548611543600661140b565b61147f565b6115526005610e1a565b906124e3565b565b611579905f3d811161157f575b61157181836104f5565b8101906112e6565b5f611520565b503d611567565b6113be565b6112dc565b5f635e1e452d60e11b8152806115a860048201610924565b0390fd5b6115b461128a565b565b6115be61218e565b6115c66115c8565b565b6115d96115d45f610e9f565b612599565b565b6115e36115b6565b565b906115f7916115f261218e565b61185e565b565b9092919261160e61160982610558565b61051e565b9381855260208501908284011161162a576116289261016c565b565b610554565b9080601f8301121561164d5781602061164a935191016115f9565b90565b6104dd565b916060838303126116ce575f83015167ffffffffffffffff81116116c9578261167c91850161162f565b92602081015167ffffffffffffffff81116116c4578361169d91830161162f565b92604082015167ffffffffffffffff81116116bf576116bc920161162f565b90565b610483565b610483565b610483565b61014c565b90565b91906116f0816116e9816116f595610b5b565b809561057b565b610177565b0190565b9061170492916116d6565b90565b5f80fd5b5f80fd5b5f80fd5b903560016020038236030381121561175457016020813591019167ffffffffffffffff821161174f57600182023603831361174a57565b61170b565b611707565b61170f565b60200190565b918161176a91610ebc565b908161177b602083028401946116d3565b92835f925b8484106117905750505050505090565b90919293949560206117bc6117b683856001950388526117b08b88611713565b906116f9565b98611759565b940194019294939190611780565b9392906117e86020916117f09460408801918883035f8a015261175f565b9401906102fb565b565b6bffffffffffffffffffffffff191690565b61180e90516117f2565b90565b61182b6118266118208361015f565b92610f8e565b611804565b9060148110611839575b5090565b611857906bffffffffffffffffffffffff1990601403600802611407565b165f611835565b905f9061188a7f0000000000000000000000000000000000000000000000000000000000000000610259565b6118cd639370bb3e9492946118d87f00000000000000000000000000000000000000000000000000000000000000006118c1610142565b978896879586956112e0565b8552600485016117ca565b03915afa90811561194857611916915f8080919092611918575b6119119293509061190761190c929491612330565b611811565b61147f565b6124e3565b565b50505061191161190c61193f611907933d805f833e61193781836104f5565b810190611652565b925092506118f2565b6113be565b90611957916115e5565b565b9061196b9161196661218e565b6119b8565b565b9190611987816119808161198c95610163565b809561057b565b610177565b0190565b9392906119ae6020916119b69460408801918883035f8a015261196d565b9401906102fb565b565b905f906119e47f0000000000000000000000000000000000000000000000000000000000000000610259565b611a27632b5f2f81949294611a327f0000000000000000000000000000000000000000000000000000000000000000611a1b610142565b978896879586956112e0565b855260048501611990565b03915afa908115611aa257611a70915f8080919092611a72575b611a6b92935090611a61611a66929491612330565b611811565b61147f565b6124e3565b565b505050611a6b611a66611a99611a61933d805f833e611a9181836104f5565b810190611652565b92509250611a4c565b6113be565b90611ab191611959565b565b611ac490611abf61218e565b611dc8565b565b905090565b905f9291805490611ae5611ade83610b31565b8094611ac6565b916001811690815f14611b375750600114611b00575b505050565b611b0d9192939450610b64565b5f905b838210611b2357505001905f8080611afb565b600181602092548486015201910190611b10565b92949550505060ff191682528015150201905f8080611afb565b611b5a91611acb565b90565b611b7290611b69610142565b91829182611b51565b03902090565b611b9d611b9492602092611b8b8161015f565b94858093611ac6565b9384910161016c565b0190565b611baa91611b78565b90565b611bc290611bb9610142565b91829182611ba1565b03902090565b601f602091010490565b91906008611bed910291611be75f1984611407565b92611407565b9181191691161790565b611c0b611c06611c10926102f8565b61022e565b6102f8565b90565b90565b9190611c2c611c27611c3493611bf7565b611c13565b908354611bd2565b9055565b5f90565b611c4e91611c48611c38565b91611c16565b565b5b818110611c5c575050565b80611c695f600193611c3c565b01611c51565b9190601f8111611c7f575b505050565b611c8b611cb093610b64565b906020611c9784611bc8565b83019310611cb8575b611ca990611bc8565b0190611c50565b5f8080611c7a565b9150611ca981929050611ca0565b90611cd6905f19906008026102bc565b191690565b81611ce591611cc6565b906002021790565b90611cf78161015f565b9067ffffffffffffffff8211611db757611d1b82611d158554610b31565b85611c6f565b602090601f8311600114611d4f57918091611d3e935f92611d43575b5050611cdb565b90555b565b90915001515f80611d37565b601f19831691611d5e85610b64565b925f5b818110611d9f57509160029391856001969410611d85575b50505002019055611d41565b611d95910151601f841690611cc6565b90555f8080611d79565b91936020600181928787015181550195019201611d61565b6104e1565b90611dc691611ced565b565b611e1f90600181611e02611dfc7f7cb84e5a76f21dbf9cffe637ecb554a3d4f4eeb2cf05cd38c8a0db9d34ed654193611b5d565b91611bad565b91611e0b610142565b80611e1581610924565b0390a36001611dbc565b565b611e2a90611ab3565b565b90611e3e91611e3961218e565b611f6e565b565b9160c083830312611f2e575f83015167ffffffffffffffff8111611f295782611e6a91850161162f565b92602081015167ffffffffffffffff8111611f245783611e8b91830161162f565b92604082015167ffffffffffffffff8111611f1f5781611eac91840161162f565b92606083015167ffffffffffffffff8111611f1a5782611ecd91850161162f565b92608081015167ffffffffffffffff8111611f155783611eee91830161162f565b9260a082015167ffffffffffffffff8111611f1057611f0d920161162f565b90565b610483565b610483565b610483565b610483565b610483565b610483565b61014c565b5f1b90565b90611f445f1991611f33565b9181191691161790565b90611f63611f5e611f6a92611bf7565b611c13565b8254611f38565b9055565b905f90611f9a7f0000000000000000000000000000000000000000000000000000000000000000610259565b611fdd63d420f66a949294611fe87f0000000000000000000000000000000000000000000000000000000000000000611fd1610142565b978896879586956112e0565b855260048501611990565b03915afa90811561209c57612047915f8080808080919293949561205d575b6120409596509161202b612032926120246120399695600a611dbc565b6009611dbc565b6008611dbc565b6007611dbc565b6006611dbc565b6005611dbc565b61205b61205460016112c0565b6004611f4e565b565b50505050505061204061202b61203961202461208d612032953d805f833e61208581836104f5565b810190611e40565b93975093975093509350612007565b6113be565b906120ab91611e2c565b565b6120b5610e3c565b506120bf5f610e59565b90565b6120d3906120ce61218e565b6120d5565b565b806120f06120ea6120e55f610e9f565b610352565b91610352565b14612100576120fe90612599565b565b61212361210c5f610e9f565b5f918291631e4fbdf760e01b83526004830161036b565b0390fd5b612130906120c2565b565b5f90565b61213e612132565b507f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b61218b916121829161217c610e3c565b5061263b565b90929192612738565b90565b6121966120ad565b6121af6121a96121a4612809565b610352565b91610352565b036121b657565b6121d86121c1612809565b5f91829163118cdaa760e01b83526004830161036b565b0390fd5b90565b6121f36121ee6121f8926121dc565b61022e565b6102f8565b90565b634e487b7160e01b5f52601160045260245ffd5b61221e612224919392936102f8565b926102f8565b820180921161222f57565b6121fb565b634e487b7160e01b5f52603260045260245ffd5b906122528261015f565b81101561226457600160209102010190565b612234565b60ff60f81b1690565b61227c9051612269565b90565b5f5260205f2090565b90602061229a818306601f039361227f565b91040191565b906122aa82610eab565b808210156122d7576020115f146122c75760209006601f0390915b565b6122d091612288565b90916122c5565b612234565b60f81b90565b6122eb906122dc565b90565b6122fe90600861230393026102bc565b6122e2565b90565b9061231191546122ee565b90565b61231d906102f8565b5f19811461232b5760010190565b6121fb565b6123398161015f565b61236761236161235c60046123576123516001610eab565b916121df565b61220f565b6102f8565b916102f8565b1061241a57612374611c38565b9161237f6001610eab565b925b8061239461238e866102f8565b916102f8565b14612414576123bf6123ba846123b4846123ae60046121df565b9061220f565b90612248565b612272565b6123e56123df6123da6123d4600186906122a0565b90612306565b612269565b91612269565b036123f8576123f390612314565b612381565b5f63719a9a4960e01b81528061241060048201610924565b0390fd5b50915050565b5f63719a9a4960e01b81528061243260048201610924565b0390fd5b929061246c9161245f61247a969461245560808801945f89019061035e565b602087019061035e565b84820360408601526112f5565b916060818403910152610181565b90565b9061248e60018060a01b0391611f33565b9181191691161790565b6124a19061024d565b90565b90565b906124bc6124b76124c392612498565b6124a4565b825461247d565b9055565b6124db6124d66124e092610e80565b61022e565b6102f8565b90565b906124ed8161015f565b6125066125006124fb610d0c565b6102f8565b916102f8565b0361257d576125616125689261251c6002610e59565b81600391612557867f1882fd6e997ef1dc5e2691efe78564a5c540b36547363e8a7dcd454c61a944b69461254e610142565b94859485612436565b0390a160026124a7565b6003611dbc565b61257b6125745f6124c7565b6004611f4e565b565b5f634ae601c160e11b81528061259560048201610924565b0390fd5b6125a25f610e59565b6125ac825f6124a7565b906125e06125da7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093612498565b91612498565b916125e9610142565b806125f381610924565b0390a3565b5f90565b90565b61261361260e612618926125fc565b61022e565b6102f8565b90565b61262f61262a612634926102f8565b611f33565b6104b7565b90565b5f90565b919091612646610e3c565b5061264f6125f8565b50612658612132565b506126628361015f565b61267561266f60416125ff565b916102f8565b145f146126bc576126b59192612689612132565b50612692612132565b5061269b612637565b506020810151606060408301519201515f1a9091926128a2565b9192909190565b506126c65f610e9f565b906126da6126d560029461015f565b61261b565b91929190565b634e487b7160e01b5f52602160045260245ffd5b600411156126fe57565b6126e0565b9061270d826126f4565b565b9190612722905f60208501940190611071565b565b61273061273591610e40565b611bf7565b90565b8061274b6127455f612703565b91612703565b145f14612756575050565b8061276a6127646001612703565b91612703565b145f1461278d575f63f645eedf60e01b81528061278960048201610924565b0390fd5b806127a161279b6002612703565b91612703565b145f146127cf576127cb6127b483612724565b5f91829163fce698f760e01b835260048301610308565b0390fd5b6127e26127dc6003612703565b91612703565b146127ea5750565b612805905f9182916335e2f38360e21b83526004830161270f565b0390fd5b612811610e3c565b503390565b90565b61282d61282861283292612816565b61022e565b6102f8565b90565b60ff1690565b61284490612835565b9052565b61287d61288494612873606094989795612869608086019a5f870190611071565b602085019061283b565b6040830190611071565b0190611071565b565b61289a61289561289f92610e80565b611f33565b6104b7565b90565b9392936128ad610e3c565b506128b66125f8565b506128bf612132565b506128c985612724565b6128fb6128f57f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0612819565b916102f8565b11612988579061291e602094955f94939293612915610142565b94859485612848565b838052039060015afa15612983576129365f51611f33565b8061295161294b6129465f610e9f565b610352565b91610352565b14612967575f916129615f612886565b91929190565b506129715f610e9f565b60019161297d5f612886565b91929190565b6113be565b5050506129945f610e9f565b906003929192919056fea264697066735822122003c0726a9fa520e35b3138febc954cbeead89e9f707dddd67662d00ef680e7ce64736f6c634300081e0033",
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
// the contract method with ID 0x54b21558.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
func (teeAuthenticator *TeeAuthenticator) PackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, refundAmount, applicationFee, signature)
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
func (teeAuthenticator *TeeAuthenticator) TryPackCheckSignature(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, refundAmount *big.Int, applicationFee *big.Int, signature []byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, refundAmount, applicationFee, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54b21558.
//
// Solidity: function checkSignature(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, uint256 refundAmount, uint256 applicationFee, bytes signature) view returns(bool)
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

// PackUpdateTeeDecoded is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74ab0168.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTeeDecoded(bytes[] attestation_decoded) returns()
func (teeAuthenticator *TeeAuthenticator) PackUpdateTeeDecoded(attestationDecoded [][]byte) []byte {
	enc, err := teeAuthenticator.abi.Pack("updateTeeDecoded", attestationDecoded)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTeeDecoded is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74ab0168.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTeeDecoded(bytes[] attestation_decoded) returns()
func (teeAuthenticator *TeeAuthenticator) TryPackUpdateTeeDecoded(attestationDecoded [][]byte) ([]byte, error) {
	return teeAuthenticator.abi.Pack("updateTeeDecoded", attestationDecoded)
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

// ConsoleMetaData contains all meta data concerning the Console contract.
var ConsoleMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "Console",
	Bin: "0x608060405234601d57600e6021565b603e602c823930815050603e90f35b6027565b60405190565b5f80fdfe60806040525f80fdfea26469706673582212205856d305f3a2b6b8899f20ae7dc888f3c23c18474a542d2e6e2ada22d7c2b7cf64736f6c634300081e0033",
}

// Console is an auto generated Go binding around an Ethereum contract.
type Console struct {
	abi abi.ABI
}

// NewConsole creates a new instance of Console.
func NewConsole() *Console {
	parsed, err := ConsoleMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Console{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Console) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}
