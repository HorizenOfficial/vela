// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc1967proxy

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

// ERC1967ProxyMetaData contains all meta data concerning the ERC1967Proxy contract.
var ERC1967ProxyMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	ID:  "ERC1967Proxy",
	Bin: "0x608060405261001561000f61019d565b906101bf565b61001d61002b565b61010c61053b823961010c90f35b60405190565b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b9061005990610031565b810190811060018060401b0382111761007157604052565b61003b565b9061008961008261002b565b928361004f565b565b5f80fd5b5f80fd5b60018060a01b031690565b6100a790610093565b90565b6100b38161009e565b036100ba57565b5f80fd5b905051906100cb826100aa565b565b5f80fd5b5f80fd5b60018060401b0381116100f1576100ed602091610031565b0190565b61003b565b90825f9392825e0152565b90929192610116610111826100d5565b610076565b9381855260208501908284011161013257610130926100f6565b565b6100d1565b9080601f830112156101555781602061015293519101610101565b90565b6100cd565b91909160408184031261019857610173835f83016100be565b92602082015160018060401b038111610193576101909201610137565b90565b61008f565b61008b565b6101bb610647803803806101b081610076565b92833981019061015a565b9091565b906101c99161022d565b565b90565b6101e26101dd6101e792610093565b6101cb565b610093565b90565b6101f3906101ce565b90565b6101ff906101ea565b90565b5f0190565b5190565b90565b90565b61022561022061022a9261020e565b6101cb565b61020b565b90565b9061023782610369565b816102627fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b916101f6565b9061026b61002b565b8061027581610202565b0390a261028181610207565b61029361028d5f610211565b9161020b565b115f146102a7576102a391610439565b505b565b50506102b16103be565b6102a5565b6102bf9061009e565b9052565b91906102d6905f602085019401906102b6565b565b90565b90565b5f1b90565b6102f76102f26102fc926102d8565b6102de565b6102db565b90565b6103287f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6102e3565b90565b9061033c60018060a01b03916102de565b9181191691161790565b90565b9061035e610359610365926101f6565b610346565b825461032b565b9055565b803b61037d6103775f610211565b9161020b565b1461039f5761039d905f6103976103926102ff565b610468565b01610349565b565b6103ba905f918291634c9c8ce360e01b8352600483016102c3565b0390fd5b346103d16103cb5f610211565b9161020b565b116103d857565b5f63b398979f60e01b8152806103f060048201610202565b0390fd5b606090565b9061040b610406836100d5565b610076565b918252565b3d5f1461042b576104203d6103f9565b903d5f602084013e5b565b6104336103f4565b90610429565b5f80610465936104476103f4565b508390602081019051915af49061045c610410565b90919091610470565b90565b90565b151590565b906104849061047d6103f4565b501561046b565b5f1461049057506104f4565b61049982610207565b6104ab6104a55f610211565b9161020b565b14806104d9575b6104ba575090565b6104d5905f918291639996b31560e01b8352600483016102c3565b0390fd5b50803b6104ee6104e85f610211565b9161020b565b146104b2565b6104fd81610207565b61050f6105095f610211565b9161020b565b115f1461051e57805190602001fd5b5f63d6bda27560e01b81528061053660048201610202565b0390fdfe6080604052600a6012565b6022565b5f90565b6018600e565b50601f60b5565b90565b5f8091368280378136915af43d5f803e5f14603b573d5ff35b3d5ffd5b90565b90565b5f1b90565b60596055605d92603f565b6045565b6042565b90565b60877f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc604a565b90565b5f1c90565b60018060a01b031690565b60a360a791608a565b608f565b90565b60b29054609a565b90565b60bb600e565b5060d05f60cb60c76060565b60d3565b0160aa565b90565b9056fea2646970667358221220c93ba8a5f5e4cc66f43dd17b1660f01d6297b1adf1a2b966b6842bc5e8486ed264736f6c634300081e0033",
}

// ERC1967Proxy is an auto generated Go binding around an Ethereum contract.
type ERC1967Proxy struct {
	abi abi.ABI
}

// NewERC1967Proxy creates a new instance of ERC1967Proxy.
func NewERC1967Proxy() *ERC1967Proxy {
	parsed, err := ERC1967ProxyMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ERC1967Proxy{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ERC1967Proxy) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address implementation, bytes _data) payable returns()
func (eRC1967Proxy *ERC1967Proxy) PackConstructor(implementation common.Address, _data []byte) []byte {
	enc, err := eRC1967Proxy.abi.Pack("", implementation, _data)
	if err != nil {
		panic(err)
	}
	return enc
}

// ERC1967ProxyUpgraded represents a Upgraded event raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ERC1967ProxyUpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ERC1967ProxyUpgraded) ContractEventName() string {
	return ERC1967ProxyUpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (eRC1967Proxy *ERC1967Proxy) UnpackUpgradedEvent(log *types.Log) (*ERC1967ProxyUpgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != eRC1967Proxy.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ERC1967ProxyUpgraded)
	if len(log.Data) > 0 {
		if err := eRC1967Proxy.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC1967Proxy.abi.Events[event].Inputs {
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
func (eRC1967Proxy *ERC1967Proxy) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackFailedCallError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ERC1967ProxyAddressEmptyCode represents a AddressEmptyCode error raised by the ERC1967Proxy contract.
type ERC1967ProxyAddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ERC1967ProxyAddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (eRC1967Proxy *ERC1967Proxy) UnpackAddressEmptyCodeError(raw []byte) (*ERC1967ProxyAddressEmptyCode, error) {
	out := new(ERC1967ProxyAddressEmptyCode)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC1967ProxyERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ERC1967Proxy contract.
type ERC1967ProxyERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ERC1967ProxyERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (eRC1967Proxy *ERC1967Proxy) UnpackERC1967InvalidImplementationError(raw []byte) (*ERC1967ProxyERC1967InvalidImplementation, error) {
	out := new(ERC1967ProxyERC1967InvalidImplementation)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC1967ProxyERC1967NonPayable represents a ERC1967NonPayable error raised by the ERC1967Proxy contract.
type ERC1967ProxyERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ERC1967ProxyERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (eRC1967Proxy *ERC1967Proxy) UnpackERC1967NonPayableError(raw []byte) (*ERC1967ProxyERC1967NonPayable, error) {
	out := new(ERC1967ProxyERC1967NonPayable)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC1967ProxyFailedCall represents a FailedCall error raised by the ERC1967Proxy contract.
type ERC1967ProxyFailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ERC1967ProxyFailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (eRC1967Proxy *ERC1967Proxy) UnpackFailedCallError(raw []byte) (*ERC1967ProxyFailedCall, error) {
	out := new(ERC1967ProxyFailedCall)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}
