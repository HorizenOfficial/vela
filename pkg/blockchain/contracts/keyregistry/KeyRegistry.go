// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keyregistry

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

// KeyRegistryMetaData contains all meta data concerning the KeyRegistry contract.
var KeyRegistryMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidLength\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"PKRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PK_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getPK\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"registerPK\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "4dc55971e232a0cb11d22592871a69dd79",
	Bin: "0x6080604052348015600e575f5ffd5b506108c58061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c8063038defd71461004e578063c91496c61461007e578063c9bc567e1461009c578063e88a0e0d146100b8575b5f5ffd5b6100686004803603810190610063919061037d565b6100e8565b6040516100759190610418565b60405180910390f35b610086610182565b6040516100939190610450565b60405180910390f35b6100b660048036038101906100b191906104ca565b610187565b005b6100d260048036038101906100cd919061037d565b61024f565b6040516100df9190610418565b60405180910390f35b5f602052805f5260405f205f91509050805461010390610542565b80601f016020809104026020016040519081016040528092919081815260200182805461012f90610542565b801561017a5780601f106101515761010080835404028352916020019161017a565b820191905f5260205f20905b81548152906001019060200180831161015d57829003601f168201915b505050505081565b608581565b608582829050146101c4576040517f947d5a8400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81815f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20918261020f929190610749565b507fc7aba4ce59c1f14da62254804f2eddcc8de9c992711b3e957d174f04b897f7443383836040516102439392919061085f565b60405180910390a15050565b60605f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20805461029890610542565b80601f01602080910402602001604051908101604052809291908181526020018280546102c490610542565b801561030f5780601f106102e65761010080835404028352916020019161030f565b820191905f5260205f20905b8154815290600101906020018083116102f257829003601f168201915b50505050509050919050565b5f5ffd5b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61034c82610323565b9050919050565b61035c81610342565b8114610366575f5ffd5b50565b5f8135905061037781610353565b92915050565b5f602082840312156103925761039161031b565b5b5f61039f84828501610369565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6103ea826103a8565b6103f481856103b2565b93506104048185602086016103c2565b61040d816103d0565b840191505092915050565b5f6020820190508181035f83015261043081846103e0565b905092915050565b5f819050919050565b61044a81610438565b82525050565b5f6020820190506104635f830184610441565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f84011261048a57610489610469565b5b8235905067ffffffffffffffff8111156104a7576104a661046d565b5b6020830191508360018202830111156104c3576104c2610471565b5b9250929050565b5f5f602083850312156104e0576104df61031b565b5b5f83013567ffffffffffffffff8111156104fd576104fc61031f565b5b61050985828601610475565b92509250509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061055957607f821691505b60208210810361056c5761056b610515565b5b50919050565b5f82905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026106057fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826105ca565b61060f86836105ca565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61064a61064561064084610438565b610627565b610438565b9050919050565b5f819050919050565b61066383610630565b61067761066f82610651565b8484546105d6565b825550505050565b5f5f905090565b61068e61067f565b61069981848461065a565b505050565b5b818110156106bc576106b15f82610686565b60018101905061069f565b5050565b601f821115610701576106d2816105a9565b6106db846105bb565b810160208510156106ea578190505b6106fe6106f6856105bb565b83018261069e565b50505b505050565b5f82821c905092915050565b5f6107215f1984600802610706565b1980831691505092915050565b5f6107398383610712565b9150826002028217905092915050565b6107538383610572565b67ffffffffffffffff81111561076c5761076b61057c565b5b6107768254610542565b6107818282856106c0565b5f601f8311600181146107ae575f841561079c578287013590505b6107a6858261072e565b86555061080d565b601f1984166107bc866105a9565b5f5b828110156107e3578489013582556001820191506020850194506020810190506107be565b8683101561080057848901356107fc601f891682610712565b8355505b6001600288020188555050505b50505050505050565b61081f81610342565b82525050565b828183375f83830152505050565b5f61083e83856103b2565b935061084b838584610825565b610854836103d0565b840190509392505050565b5f6040820190506108725f830186610816565b8181036020830152610885818486610833565b905094935050505056fea2646970667358221220731e3ae413c36d2d8648738bd504be21e99f1cde361ff38d198715c3fc6fa5aa64736f6c634300081e0033",
}

// KeyRegistry is an auto generated Go binding around an Ethereum contract.
type KeyRegistry struct {
	abi abi.ABI
}

// NewKeyRegistry creates a new instance of KeyRegistry.
func NewKeyRegistry() *KeyRegistry {
	parsed, err := KeyRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &KeyRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *KeyRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackPKLENGTH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc91496c6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (keyRegistry *KeyRegistry) PackPKLENGTH() []byte {
	enc, err := keyRegistry.abi.Pack("PK_LENGTH")
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
func (keyRegistry *KeyRegistry) TryPackPKLENGTH() ([]byte, error) {
	return keyRegistry.abi.Pack("PK_LENGTH")
}

// UnpackPKLENGTH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc91496c6.
//
// Solidity: function PK_LENGTH() view returns(uint256)
func (keyRegistry *KeyRegistry) UnpackPKLENGTH(data []byte) (*big.Int, error) {
	out, err := keyRegistry.abi.Unpack("PK_LENGTH", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPK is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe88a0e0d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPK(address owner) view returns(bytes)
func (keyRegistry *KeyRegistry) PackGetPK(owner common.Address) []byte {
	enc, err := keyRegistry.abi.Pack("getPK", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPK is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe88a0e0d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPK(address owner) view returns(bytes)
func (keyRegistry *KeyRegistry) TryPackGetPK(owner common.Address) ([]byte, error) {
	return keyRegistry.abi.Pack("getPK", owner)
}

// UnpackGetPK is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe88a0e0d.
//
// Solidity: function getPK(address owner) view returns(bytes)
func (keyRegistry *KeyRegistry) UnpackGetPK(data []byte) ([]byte, error) {
	out, err := keyRegistry.abi.Unpack("getPK", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackRegisterPK is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9bc567e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerPK(bytes publicKey) returns()
func (keyRegistry *KeyRegistry) PackRegisterPK(publicKey []byte) []byte {
	enc, err := keyRegistry.abi.Pack("registerPK", publicKey)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterPK is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9bc567e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerPK(bytes publicKey) returns()
func (keyRegistry *KeyRegistry) TryPackRegisterPK(publicKey []byte) ([]byte, error) {
	return keyRegistry.abi.Pack("registerPK", publicKey)
}

// PackRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x038defd7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registry(address ) view returns(bytes)
func (keyRegistry *KeyRegistry) PackRegistry(arg0 common.Address) []byte {
	enc, err := keyRegistry.abi.Pack("registry", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x038defd7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registry(address ) view returns(bytes)
func (keyRegistry *KeyRegistry) TryPackRegistry(arg0 common.Address) ([]byte, error) {
	return keyRegistry.abi.Pack("registry", arg0)
}

// UnpackRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x038defd7.
//
// Solidity: function registry(address ) view returns(bytes)
func (keyRegistry *KeyRegistry) UnpackRegistry(data []byte) ([]byte, error) {
	out, err := keyRegistry.abi.Unpack("registry", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// KeyRegistryPKRegistered represents a PKRegistered event raised by the KeyRegistry contract.
type KeyRegistryPKRegistered struct {
	Owner     common.Address
	PublicKey []byte
	Raw       *types.Log // Blockchain specific contextual infos
}

const KeyRegistryPKRegisteredEventName = "PKRegistered"

// ContractEventName returns the user-defined event name.
func (KeyRegistryPKRegistered) ContractEventName() string {
	return KeyRegistryPKRegisteredEventName
}

// UnpackPKRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PKRegistered(address owner, bytes publicKey)
func (keyRegistry *KeyRegistry) UnpackPKRegisteredEvent(log *types.Log) (*KeyRegistryPKRegistered, error) {
	event := "PKRegistered"
	if log.Topics[0] != keyRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(KeyRegistryPKRegistered)
	if len(log.Data) > 0 {
		if err := keyRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range keyRegistry.abi.Events[event].Inputs {
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
func (keyRegistry *KeyRegistry) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], keyRegistry.abi.Errors["InvalidLength"].ID.Bytes()[:4]) {
		return keyRegistry.UnpackInvalidLengthError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// KeyRegistryInvalidLength represents a InvalidLength error raised by the KeyRegistry contract.
type KeyRegistryInvalidLength struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidLength()
func KeyRegistryInvalidLengthErrorID() common.Hash {
	return common.HexToHash("0x947d5a8466b7bcfb56468ee05345bfacfb6c28f016a57c192dadee0d8affb197")
}

// UnpackInvalidLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidLength()
func (keyRegistry *KeyRegistry) UnpackInvalidLengthError(raw []byte) (*KeyRegistryInvalidLength, error) {
	out := new(KeyRegistryInvalidLength)
	if err := keyRegistry.abi.UnpackIntoInterface(out, "InvalidLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}
