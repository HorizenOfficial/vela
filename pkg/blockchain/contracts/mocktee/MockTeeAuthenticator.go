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

// MockTeeAuthenticatorMetaData contains all meta data concerning the MockTeeAuthenticator contract.
var MockTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_teeSigner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pubSecp521r1\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "21224ad20c4253145139ff6269df243105",
	Bin: "0x608060405234801561000f575f5ffd5b50604051610f21380380610f218339818101604052810190610031919061022e565b815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806001908161007f9190610498565b505050610567565b5f604051905090565b5f5ffd5b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100c182610098565b9050919050565b6100d1816100b7565b81146100db575f5ffd5b50565b5f815190506100ec816100c8565b92915050565b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610140826100fa565b810181811067ffffffffffffffff8211171561015f5761015e61010a565b5b80604052505050565b5f610171610087565b905061017d8282610137565b919050565b5f67ffffffffffffffff82111561019c5761019b61010a565b5b6101a5826100fa565b9050602081019050919050565b8281835e5f83830152505050565b5f6101d26101cd84610182565b610168565b9050828152602081018484840111156101ee576101ed6100f6565b5b6101f98482856101b2565b509392505050565b5f82601f830112610215576102146100f2565b5b81516102258482602086016101c0565b91505092915050565b5f5f6040838503121561024457610243610090565b5b5f610251858286016100de565b925050602083015167ffffffffffffffff81111561027257610271610094565b5b61027e85828601610201565b9150509250929050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806102d657607f821691505b6020821081036102e9576102e8610292565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261034b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610310565b6103558683610310565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61039961039461038f8461036d565b610376565b61036d565b9050919050565b5f819050919050565b6103b28361037f565b6103c66103be826103a0565b84845461031c565b825550505050565b5f5f905090565b6103dd6103ce565b6103e88184846103a9565b505050565b5b8181101561040b576104005f826103d5565b6001810190506103ee565b5050565b601f82111561045057610421816102ef565b61042a84610301565b81016020851015610439578190505b61044d61044585610301565b8301826103ed565b50505b505050565b5f82821c905092915050565b5f6104705f1984600802610455565b1980831691505092915050565b5f6104888383610461565b9150826002028217905092915050565b6104a182610288565b67ffffffffffffffff8111156104ba576104b961010a565b5b6104c482546102bf565b6104cf82828561040f565b5f60209050601f831160018114610500575f84156104ee578287015190505b6104f8858261047d565b86555061055f565b601f19841661050e866102ef565b5f5b8281101561053557848901518255600182019150602085019450602081019050610510565b86831015610552578489015161054e601f891682610461565b8355505b6001600288020188555050505b505050505050565b6109ad806105745f395ff3fe608060405234801561000f575f5ffd5b5060043610610055575f3560e01c8063081bec7e146100595780630dd7ce2f1461007757806343f855c314610095578063aac80ccc146100b3578063fe4993ca146100e3575b5f5ffd5b610061610101565b60405161006e91906102eb565b60405180910390f35b61007f610191565b60405161008c919061034a565b60405180910390f35b61009d6101b8565b6040516100aa919061034a565b60405180910390f35b6100cd60048036038101906100c891906107ca565b6101dc565b6040516100da9190610901565b60405180910390f35b6100eb6101ef565b6040516100f891906102eb565b60405180910390f35b60606001805461011090610947565b80601f016020809104026020016040519081016040528092919081815260200182805461013c90610947565b80156101875780601f1061015e57610100808354040283529160200191610187565b820191905f5260205f20905b81548152906001019060200180831161016a57829003601f168201915b5050505050905090565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f600190509a9950505050505050505050565b600180546101fc90610947565b80601f016020809104026020016040519081016040528092919081815260200182805461022890610947565b80156102735780601f1061024a57610100808354040283529160200191610273565b820191905f5260205f20905b81548152906001019060200180831161025657829003601f168201915b505050505081565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6102bd8261027b565b6102c78185610285565b93506102d7818560208601610295565b6102e0816102a3565b840191505092915050565b5f6020820190508181035f83015261030381846102b3565b905092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6103348261030b565b9050919050565b6103448161032a565b82525050565b5f60208201905061035d5f83018461033b565b92915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f67ffffffffffffffff82169050919050565b61039081610374565b811461039a575f5ffd5b50565b5f813590506103ab81610387565b92915050565b5f819050919050565b6103c3816103b1565b81146103cd575f5ffd5b50565b5f813590506103de816103ba565b92915050565b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61041e826102a3565b810181811067ffffffffffffffff8211171561043d5761043c6103e8565b5b80604052505050565b5f61044f610363565b905061045b8282610415565b919050565b5f67ffffffffffffffff82111561047a576104796103e8565b5b602082029050602081019050919050565b5f5ffd5b5f5ffd5b5f67ffffffffffffffff8211156104ad576104ac6103e8565b5b6104b6826102a3565b9050602081019050919050565b828183375f83830152505050565b5f6104e36104de84610493565b610446565b9050828152602081018484840111156104ff576104fe61048f565b5b61050a8482856104c3565b509392505050565b5f82601f830112610526576105256103e4565b5b81356105368482602086016104d1565b91505092915050565b5f61055161054c84610460565b610446565b905080838252602082019050602084028301858111156105745761057361048b565b5b835b818110156105bb57803567ffffffffffffffff811115610599576105986103e4565b5b8086016105a68982610512565b85526020850194505050602081019050610576565b5050509392505050565b5f82601f8301126105d9576105d86103e4565b5b81356105e984826020860161053f565b91505092915050565b5f67ffffffffffffffff82111561060c5761060b6103e8565b5b602082029050602081019050919050565b5f5ffd5b5f61062b8261030b565b9050919050565b61063b81610621565b8114610645575f5ffd5b50565b5f8135905061065681610632565b92915050565b5f819050919050565b61066e8161065c565b8114610678575f5ffd5b50565b5f8135905061068981610665565b92915050565b5f604082840312156106a4576106a361061d565b5b6106ae6040610446565b90505f6106bd84828501610648565b5f8301525060206106d08482850161067b565b60208301525092915050565b5f6106ee6106e9846105f2565b610446565b905080838252602082019050604084028301858111156107115761071061048b565b5b835b8181101561073a5780610726888261068f565b845260208401935050604081019050610713565b5050509392505050565b5f82601f830112610758576107576103e4565b5b81356107688482602086016106dc565b91505092915050565b5f5ffd5b5f5f83601f84011261078a576107896103e4565b5b8235905067ffffffffffffffff8111156107a7576107a6610771565b5b6020830191508360018202830111156107c3576107c261048b565b5b9250929050565b5f5f5f5f5f5f5f5f5f5f6101208b8d0312156107e9576107e861036c565b5b5f6107f68d828e0161039d565b9a505060206108078d828e016103d0565b99505060406108188d828e016103d0565b98505060606108298d828e016103d0565b97505060808b013567ffffffffffffffff81111561084a57610849610370565b5b6108568d828e016105c5565b96505060a08b013567ffffffffffffffff81111561087757610876610370565b5b6108838d828e01610744565b95505060c08b013567ffffffffffffffff8111156108a4576108a3610370565b5b6108b08d828e01610775565b945094505060e06108c38d828e0161067b565b9250506101006108d58d828e0161067b565b9150509295989b9194979a5092959850565b5f8115159050919050565b6108fb816108e7565b82525050565b5f6020820190506109145f8301846108f2565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061095e57607f821691505b6020821081036109715761097061091a565b5b5091905056fea26469706673582212207de88aec87fae779b90791220e43bbe6328fe2a15a69da96149fffbc5a9cdc0364736f6c634300081e0033",
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

// PackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaac80ccc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature(uint64 , bytes32 , bytes32 , bytes32 , bytes[] , (address,uint256)[] , bytes , uint256 , uint256 ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) PackCheckSignature(arg0 uint64, arg1 [32]byte, arg2 [32]byte, arg3 [32]byte, arg4 [][]byte, arg5 []StructsWithdrawalRequest, arg6 []byte, arg7 *big.Int, arg8 *big.Int) []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("checkSignature", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaac80ccc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature(uint64 , bytes32 , bytes32 , bytes32 , bytes[] , (address,uint256)[] , bytes , uint256 , uint256 ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackCheckSignature(arg0 uint64, arg1 [32]byte, arg2 [32]byte, arg3 [32]byte, arg4 [][]byte, arg5 []StructsWithdrawalRequest, arg6 []byte, arg7 *big.Int, arg8 *big.Int) ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("checkSignature", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaac80ccc.
//
// Solidity: function checkSignature(uint64 , bytes32 , bytes32 , bytes32 , bytes[] , (address,uint256)[] , bytes , uint256 , uint256 ) pure returns(bool)
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

// StructsMetaData contains all meta data concerning the Structs contract.
var StructsMetaData = bind.MetaData{
	ABI: "[]",
	ID:  "31c70fec60e1e1f4cc22f993498ea8b973",
	Bin: "0x6080604052348015600e575f5ffd5b50603e80601a5f395ff3fe60806040525f5ffdfea26469706673582212205db0187e1826b88ea0e56c0776b4649c41772c8ec48a6499b3e8b3c9c7997b9164736f6c634300081e0033",
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
