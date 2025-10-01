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

// StructsWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsWithdrawalRequest struct {
	Receiver common.Address
	Amount   *big.Int
}

// ITeeAuthenticatorMetaData contains all meta data concerning the ITeeAuthenticator contract.
var ITeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"processedRequestId\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0x24a63b01.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, uint256 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId *big.Int, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x24a63b01.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, uint256 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId *big.Int, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x24a63b01.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, uint256 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := iTeeAuthenticator.abi.Unpack("checkSignature", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// MockTeeAuthenticatorMetaData contains all meta data concerning the MockTeeAuthenticator contract.
var MockTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"processedRequestId\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "21224ad20c4253145139ff6269df243105",
	Bin: "0x6080604052348015600e575f5ffd5b506106258061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610029575f3560e01c806324a63b011461002d575b5f5ffd5b610047600480360381019061004291906104c7565b61005d565b60405161005491906105d6565b60405180910390f35b5f6001905098975050505050505050565b5f604051905090565b5f5ffd5b5f5ffd5b5f819050919050565b6100918161007f565b811461009b575f5ffd5b50565b5f813590506100ac81610088565b92915050565b5f819050919050565b6100c4816100b2565b81146100ce575f5ffd5b50565b5f813590506100df816100bb565b92915050565b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61012f826100e9565b810181811067ffffffffffffffff8211171561014e5761014d6100f9565b5b80604052505050565b5f61016061006e565b905061016c8282610126565b919050565b5f67ffffffffffffffff82111561018b5761018a6100f9565b5b602082029050602081019050919050565b5f5ffd5b5f5ffd5b5f67ffffffffffffffff8211156101be576101bd6100f9565b5b6101c7826100e9565b9050602081019050919050565b828183375f83830152505050565b5f6101f46101ef846101a4565b610157565b9050828152602081018484840111156102105761020f6101a0565b5b61021b8482856101d4565b509392505050565b5f82601f830112610237576102366100e5565b5b81356102478482602086016101e2565b91505092915050565b5f61026261025d84610171565b610157565b905080838252602082019050602084028301858111156102855761028461019c565b5b835b818110156102cc57803567ffffffffffffffff8111156102aa576102a96100e5565b5b8086016102b78982610223565b85526020850194505050602081019050610287565b5050509392505050565b5f82601f8301126102ea576102e96100e5565b5b81356102fa848260208601610250565b91505092915050565b5f67ffffffffffffffff82111561031d5761031c6100f9565b5b602082029050602081019050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61035b82610332565b9050919050565b61036b81610351565b8114610375575f5ffd5b50565b5f8135905061038681610362565b92915050565b5f604082840312156103a1576103a061032e565b5b6103ab6040610157565b90505f6103ba84828501610378565b5f8301525060206103cd8482850161009e565b60208301525092915050565b5f6103eb6103e684610303565b610157565b9050808382526020820190506040840283018581111561040e5761040d61019c565b5b835b818110156104375780610423888261038c565b845260208401935050604081019050610410565b5050509392505050565b5f82601f830112610455576104546100e5565b5b81356104658482602086016103d9565b91505092915050565b5f5ffd5b5f5f83601f840112610487576104866100e5565b5b8235905067ffffffffffffffff8111156104a4576104a361046e565b5b6020830191508360018202830111156104c0576104bf61019c565b5b9250929050565b5f5f5f5f5f5f5f5f60e0898b0312156104e3576104e2610077565b5b5f6104f08b828c0161009e565b98505060206105018b828c016100d1565b97505060406105128b828c016100d1565b96505060606105238b828c0161009e565b955050608089013567ffffffffffffffff8111156105445761054361007b565b5b6105508b828c016102d6565b94505060a089013567ffffffffffffffff8111156105715761057061007b565b5b61057d8b828c01610441565b93505060c089013567ffffffffffffffff81111561059e5761059d61007b565b5b6105aa8b828c01610472565b92509250509295985092959890939650565b5f8115159050919050565b6105d0816105bc565b82525050565b5f6020820190506105e95f8301846105c7565b9291505056fea2646970667358221220754951335a1d66dcff436d63d7ec32b21a17bf9c1c8e26aa7c08ff117384cfe464736f6c634300081c0033",
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

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x24a63b01.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, uint256 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) PackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId *big.Int, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x24a63b01.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, uint256 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackCheckSignature(applicationId *big.Int, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId *big.Int, events [][]byte, withdrawalRequests []StructsWithdrawalRequest, signature []byte) ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("checkSignature", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x24a63b01.
//
// Solidity: function checkSignature(uint256 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, uint256 processedRequestId, bytes[] events, (address,uint256)[] withdrawalRequests, bytes signature) view returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackCheckSignature(data []byte) (bool, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("checkSignature", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// StructsMetaData contains all meta data concerning the Structs contract.
var StructsMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "31c70fec60e1e1f4cc22f993498ea8b973",
	Bin: "0x6080604052348015600e575f5ffd5b50603e80601a5f395ff3fe60806040525f5ffdfea26469706673582212207efaf3358ee1a5747e39239596174b23b97d1b258022d22858405af56fd757d964736f6c634300081c0033",
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
