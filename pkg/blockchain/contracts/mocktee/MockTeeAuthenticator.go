// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mocktee

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

// ITeeAuthenticatorMetaData contains all meta data concerning the ITeeAuthenticator contract.
var ITeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"EmptyBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"entryHashes\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkBatchSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// PackCheckBatchSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcbe37a7f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkBatchSignature(bytes32[] entryHashes, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckBatchSignature(entryHashes [][32]byte, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkBatchSignature", entryHashes, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckBatchSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcbe37a7f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkBatchSignature(bytes32[] entryHashes, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckBatchSignature(entryHashes [][32]byte, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkBatchSignature", entryHashes, signature)
}

// UnpackCheckBatchSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcbe37a7f.
//
// Solidity: function checkBatchSignature(bytes32[] entryHashes, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) UnpackCheckBatchSignature(data []byte) (bool, error) {
	out, err := iTeeAuthenticator.abi.Unpack("checkBatchSignature", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
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
	if bytes.Equal(raw[:4], iTeeAuthenticator.abi.Errors["EmptyBatch"].ID.Bytes()[:4]) {
		return iTeeAuthenticator.UnpackEmptyBatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iTeeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return iTeeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ITeeAuthenticatorEmptyBatch represents a EmptyBatch error raised by the ITeeAuthenticator contract.
type ITeeAuthenticatorEmptyBatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBatch()
func ITeeAuthenticatorEmptyBatchErrorID() common.Hash {
	return common.HexToHash("0xc2e5347df6cf8d1bf68f8a9651642312bd381b900d857d97cd665976180d6ae0")
}

// UnpackEmptyBatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBatch()
func (iTeeAuthenticator *ITeeAuthenticator) UnpackEmptyBatchError(raw []byte) (*ITeeAuthenticatorEmptyBatch, error) {
	out := new(ITeeAuthenticatorEmptyBatch)
	if err := iTeeAuthenticator.abi.UnpackIntoInterface(out, "EmptyBatch", raw); err != nil {
		return nil, err
	}
	return out, nil
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

// MockTeeAuthenticatorMetaData contains all meta data concerning the MockTeeAuthenticator contract.
var MockTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_teeSigner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pubSecp521r1\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkBatchSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "21224ad20c4253145139ff6269df243105",
	Bin: "0x6080604052346100305761001a6100146101ab565b9061049b565b610022610035565b610c7a6104b38239610c7a90f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b5f80fd5b60018060a01b031690565b6100b5906100a1565b90565b6100c1816100ac565b036100c857565b5f80fd5b905051906100d9826100b8565b565b5f80fd5b5f80fd5b60018060401b0381116100ff576100fb60209161003f565b0190565b610049565b90825f9392825e0152565b9092919261012461011f826100e3565b610084565b938185526020850190828401116101405761013e92610104565b565b6100df565b9080601f83011215610163578160206101609351910161010f565b90565b6100db565b9190916040818403126101a657610181835f83016100cc565b92602082015160018060401b0381116101a15761019e9201610145565b90565b61009d565b610099565b6101c961112d803803806101be81610084565b928339810190610168565b9091565b5f1b90565b906101e360018060a01b03916101cd565b9181191691161790565b90565b6102046101ff610209926100a1565b6101ed565b6100a1565b90565b610215906101f0565b90565b6102219061020c565b90565b90565b9061023c61023761024392610218565b610224565b82546101d2565b9055565b5190565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561027f575b602083101461027a57565b61024b565b91607f169161026f565b5f5260205f2090565b601f602091010490565b1b90565b919060086102bb9102916102b55f198461029c565b9261029c565b9181191691161790565b90565b6102dc6102d76102e1926102c5565b6101ed565b6102c5565b90565b90565b91906102fd6102f8610305936102c8565b6102e4565b9083546102a0565b9055565b5f90565b61031f91610319610309565b916102e7565b565b5b81811061032d575050565b8061033a5f60019361030d565b01610322565b9190601f8111610350575b505050565b61035c61038193610289565b90602061036884610292565b83019310610389575b61037a90610292565b0190610321565b5f808061034b565b915061037a81929050610371565b1c90565b906103ab905f1990600802610397565b191690565b816103ba9161039b565b906002021790565b906103cc81610247565b9060018060401b03821161048a576103ee826103e8855461025f565b85610340565b602090601f831160011461042257918091610411935f92610416575b50506103b0565b90555b565b90915001515f8061040a565b601f1983169161043185610289565b925f5b81811061047257509160029391856001969410610458575b50505002019055610414565b610468910151601f84169061039b565b90555f808061044c565b91936020600181928787015181550195019201610434565b610049565b90610499916103c2565b565b906104a96104b0925f610227565b600161048f565b56fe60806040526004361015610013575b610bb4565b61001d5f3561007c565b8063081bec7e146100775780630dd7ce2f1461007257806343f855c31461006d5780635c9626b914610068578063cbe37a7f146100635763fe4993ca0361000e57610b7f565b610a01565b6108ea565b6101f0565b610178565b61010a565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f91031261009a57565b61008c565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b6100e06100e96020936100ee936100d78161009f565b938480936100a3565b958691016100ac565b6100b7565b0190565b6101079160208201915f8184039101526100c1565b90565b3461013a5761011a366004610090565b610136610125610bc9565b61012d610082565b918291826100f2565b0390f35b610088565b60018060a01b031690565b6101539061013f565b90565b61015f9061014a565b9052565b9190610176905f60208501940190610156565b565b346101a857610188366004610090565b6101a4610193610c09565b61019b610082565b91829182610163565b0390f35b610088565b1c90565b60018060a01b031690565b6101cc9060086101d193026101ad565b6101b1565b90565b906101df91546101bc565b90565b6101ed5f5f906101d4565b90565b3461022057610200366004610090565b61021c61020b6101e2565b610213610082565b91829182610163565b0390f35b610088565b5f80fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b9061024b906100b7565b810190811067ffffffffffffffff82111761026557604052565b61022d565b9061027d610276610082565b9283610241565b565b5f80fd5b67ffffffffffffffff1690565b61029981610283565b036102a057565b5f80fd5b905035906102b182610290565b565b90565b6102bf816102b3565b036102c657565b5f80fd5b905035906102d7826102b6565b565b5f80fd5b67ffffffffffffffff81116102f55760208091020190565b61022d565b5f80fd5b5f80fd5b67ffffffffffffffff81116103205761031c6020916100b7565b0190565b61022d565b90825f939282370152565b9092919261034561034082610302565b61026a565b938185526020850190828401116103615761035f92610325565b565b6102fe565b9080601f830112156103845781602061038193359101610330565b90565b6102d9565b92919061039d610398826102dd565b61026a565b93818552602080860192028101918383116103f45781905b8382106103c3575050505050565b813567ffffffffffffffff81116103ef576020916103e48784938701610366565b8152019101906103b5565b6102d9565b6102fa565b9080601f830112156104175781602061041493359101610389565b90565b6102d9565b67ffffffffffffffff81116104345760208091020190565b61022d565b9092919261044e6104498261041c565b61026a565b938185526020808601920283019281841161048b57915b8383106104725750505050565b6020809161048084866102ca565b815201920191610465565b6102fa565b9080601f830112156104ae578160206104ab93359101610439565b90565b6102d9565b91909160408184031261051d576104ca604061026a565b925f82013567ffffffffffffffff811161051857816104ea9184016103f9565b5f850152602082013567ffffffffffffffff81116105135761050c9201610490565b6020830152565b61027f565b61027f565b610229565b67ffffffffffffffff811161053a5760208091020190565b61022d565b6105488161014a565b0361054f57565b5f80fd5b905035906105608261053f565b565b61056b9061013f565b90565b61057781610562565b0361057e57565b5f80fd5b9050359061058f8261056e565b565b90565b61059d81610591565b036105a457565b5f80fd5b905035906105b582610594565b565b9190606083820312610603576105fc906105d1606061026a565b936105de825f8301610553565b5f8601526105ef8260208301610582565b60208601526040016105a8565b6040830152565b610229565b9092919261061d61061882610522565b61026a565b93818552606060208601920283019281841161065c57915b8383106106425750505050565b602060609161065184866105b7565b815201920191610635565b6102fa565b9080601f8301121561067f5781602061067c93359101610608565b90565b6102d9565b600c111561068e57565b5f80fd5b9050359061069f82610684565b565b67ffffffffffffffff81116106bf576106bb6020916100b7565b0190565b61022d565b909291926106d96106d4826106a1565b61026a565b938185526020850190828401116106f5576106f392610325565b565b6102fe565b9080601f8301121561071857816020610715933591016106c4565b90565b6102d9565b919091610160818403126108615761073661016061026a565b92610743815f84016102a4565b5f85015261075481602084016102ca565b602085015261076681604084016102ca565b604085015261077881606084016102ca565b6060850152608082013567ffffffffffffffff811161085c578161079d9184016104b3565b608085015260a082013567ffffffffffffffff811161085757816107c29184016104b3565b60a085015260c082013567ffffffffffffffff811161085257816107e7918401610661565b60c08501526107f98160e084016105a8565b60e085015261080c8161010084016105a8565b610100850152610820816101208401610692565b61012085015261014082013567ffffffffffffffff811161084d5761084592016106fa565b610140830152565b61027f565b61027f565b61027f565b61027f565b610229565b9190916040818403126108be575f81013567ffffffffffffffff81116108b9578361089291830161071d565b92602082013567ffffffffffffffff81116108b4576108b19201610366565b90565b610225565b610225565b61008c565b151590565b6108d1906108c3565b9052565b91906108e8905f602085019401906108c8565b565b3461091b57610917610906610900366004610866565b90610c22565b61090e610082565b918291826108d5565b0390f35b610088565b5f80fd5b909182601f8301121561095e5781359167ffffffffffffffff831161095957602001926020830284011161095457565b6102fa565b610920565b6102d9565b909182601f8301121561099d5781359167ffffffffffffffff831161099857602001926001830284011161099357565b6102fa565b610920565b6102d9565b90916040828403126109fc575f82013567ffffffffffffffff81116109f757836109cd918401610924565b929093602082013567ffffffffffffffff81116109f2576109ee9201610963565b9091565b610225565b610225565b61008c565b34610a3557610a31610a20610a173660046109a2565b92919091610c32565b610a28610082565b918291826108d5565b0390f35b610088565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610a81575b6020831014610a7c57565b610a4d565b91607f1691610a71565b60209181520190565b5f5260205f2090565b905f9291805490610ab7610ab083610a61565b8094610a8b565b916001811690815f14610b0e5750600114610ad2575b505050565b610adf9192939450610a94565b915f925b818410610af657505001905f8080610acd565b60018160209295939554848601520191019290610ae3565b92949550505060ff19168252151560200201905f8080610acd565b90610b3391610a9d565b90565b90610b56610b4f92610b46610082565b93848092610b29565b0383610241565b565b905f10610b6b57610b6890610b36565b90565b610a3a565b610b7c60015f90610b58565b90565b34610baf57610b8f366004610090565b610bab610b9a610b70565b610ba2610082565b918291826100f2565b0390f35b610088565b5f80fd5b606090565b610bc690610b36565b90565b610bd1610bb8565b50610bdc6001610bbd565b90565b5f90565b5f1c90565b610bf4610bf991610be3565b6101b1565b90565b610c069054610be8565b90565b610c11610bdf565b50610c1b5f610bfc565b90565b5f90565b5050610c2c610c1e565b50600190565b50505050610c3e610c1e565b5060019056fea26469706673582212205552fe653fff54f08a3cdd91cc1e4d819209acba512d88d9e6fdbb1e1f6f115d64736f6c634300081e0033",
}

// MockTeeAuthenticator is an auto generated Go binding around an Ethereum contract.
type MockTeeAuthenticator struct {
	abi abi.ABI
}

// NewMockTeeAuthenticator creates a new instance of MockTeeAuthenticator.
func NewMockTeeAuthenticator() *MockTeeAuthenticator {
	parsed, err := MockTeeAuthenticatorMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MockTeeAuthenticator{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MockTeeAuthenticator) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _teeSigner, bytes _pubSecp521r1) returns()
func (mockTeeAuthenticator *MockTeeAuthenticator) PackConstructor(_teeSigner common.Address, _pubSecp521r1 []byte) []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("", _teeSigner, _pubSecp521r1)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCheckBatchSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcbe37a7f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkBatchSignature(bytes32[] , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) PackCheckBatchSignature(arg0 [][32]byte, arg1 []byte) []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("checkBatchSignature", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckBatchSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcbe37a7f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkBatchSignature(bytes32[] , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackCheckBatchSignature(arg0 [][32]byte, arg1 []byte) ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("checkBatchSignature", arg0, arg1)
}

// UnpackCheckBatchSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcbe37a7f.
//
// Solidity: function checkBatchSignature(bytes32[] , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackCheckBatchSignature(data []byte) (bool, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("checkBatchSignature", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c9626b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) PackCheckSignature(arg0 StructsSignatureParams, arg1 []byte) []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("checkSignature", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c9626b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackCheckSignature(arg0 StructsSignatureParams, arg1 []byte) ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("checkSignature", arg0, arg1)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c9626b9.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string) , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("checkSignature", data)
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
func (mockTeeAuthenticator *MockTeeAuthenticator) PackGetPubSecp521r1() []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("getPubSecp521r1")
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
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackGetPubSecp521r1() ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("getPubSecp521r1")
}

// UnpackGetPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081bec7e.
//
// Solidity: function getPubSecp521r1() view returns(bytes)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackGetPubSecp521r1(data []byte) ([]byte, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("getPubSecp521r1", data)
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
func (mockTeeAuthenticator *MockTeeAuthenticator) PackGetTeeSigner() []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("getTeeSigner")
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
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackGetTeeSigner() ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("getTeeSigner")
}

// UnpackGetTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dd7ce2f.
//
// Solidity: function getTeeSigner() view returns(address)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackGetTeeSigner(data []byte) (common.Address, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("getTeeSigner", data)
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
func (mockTeeAuthenticator *MockTeeAuthenticator) PackPubSecp521r1() []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("pubSecp521r1")
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
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackPubSecp521r1() ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("pubSecp521r1")
}

// UnpackPubSecp521r1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfe4993ca.
//
// Solidity: function pubSecp521r1() view returns(bytes)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackPubSecp521r1(data []byte) ([]byte, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("pubSecp521r1", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackTeeSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x43f855c3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function teeSigner() view returns(address)
func (mockTeeAuthenticator *MockTeeAuthenticator) PackTeeSigner() []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("teeSigner")
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
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackTeeSigner() ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("teeSigner")
}

// UnpackTeeSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x43f855c3.
//
// Solidity: function teeSigner() view returns(address)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackTeeSigner(data []byte) (common.Address, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("teeSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], mockTeeAuthenticator.abi.Errors["EmptyBatch"].ID.Bytes()[:4]) {
		return mockTeeAuthenticator.UnpackEmptyBatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockTeeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return mockTeeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MockTeeAuthenticatorEmptyBatch represents a EmptyBatch error raised by the MockTeeAuthenticator contract.
type MockTeeAuthenticatorEmptyBatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBatch()
func MockTeeAuthenticatorEmptyBatchErrorID() common.Hash {
	return common.HexToHash("0xc2e5347df6cf8d1bf68f8a9651642312bd381b900d857d97cd665976180d6ae0")
}

// UnpackEmptyBatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBatch()
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackEmptyBatchError(raw []byte) (*MockTeeAuthenticatorEmptyBatch, error) {
	out := new(MockTeeAuthenticatorEmptyBatch)
	if err := mockTeeAuthenticator.abi.UnpackIntoInterface(out, "EmptyBatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockTeeAuthenticatorTeeIsNotSet represents a TeeIsNotSet error raised by the MockTeeAuthenticator contract.
type MockTeeAuthenticatorTeeIsNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TeeIsNotSet()
func MockTeeAuthenticatorTeeIsNotSetErrorID() common.Hash {
	return common.HexToHash("0xf64b1d7b3df6940d5da676c1cc07a6144e620b0858b161be25b6cd0ef7569425")
}

// UnpackTeeIsNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TeeIsNotSet()
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackTeeIsNotSetError(raw []byte) (*MockTeeAuthenticatorTeeIsNotSet, error) {
	out := new(MockTeeAuthenticatorTeeIsNotSet)
	if err := mockTeeAuthenticator.abi.UnpackIntoInterface(out, "TeeIsNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StructsMetaData contains all meta data concerning the Structs contract.
var StructsMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "31c70fec60e1e1f4cc22f993498ea8b973",
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea2646970667358221220a4f4fe37ae853c2c77d86d33f00d76a45edb69e26bf6cb136d413faccbf32fce64736f6c634300081e0033",
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
