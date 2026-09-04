// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tokenallowlist

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

// TokenAllowlistMetaData contains all meta data concerning the TokenAllowlist contract.
var TokenAllowlistMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenAddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenAllowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenRemoved\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"addAllowedToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedTokens\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isAllowedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"removeAllowedToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "TokenAllowlist",
	Bin: "0x6080346100ce57601f6109cd38819003918201601f19168301916001600160401b038311848410176100d2578084926020946040528339810103126100ce57516001600160a01b03811681036100ce57610058906100e6565b505f5160206109ad5f395f51905f525f81815260208190527f5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b8805490839055604051929182907fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9080a461080e908161017f8239f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b6001600160a01b0381165f9081525f51602061098d5f395f51905f52602052604090205460ff16610179576001600160a01b03165f8181525f51602061098d5f395f51905f5260205260408120805460ff191660011790553391905f5160206109ad5f395f51905f52907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b505f9056fe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a71461055357508063024ece89146103d3578063248a9ca3146103a15780632f2ff15d1461036457806336568abe146103205780634178617f1461022d57806390469a9d146101b457806391d148541461016c578063a217fddf14610152578063cbe230c314610125578063d547741f146100e15763e744092e146100a0575f80fd5b346100dd5760203660031901126100dd576001600160a01b036100c16105bc565b165f526001602052602060ff60405f2054166040519015158152f35b5f80fd5b346100dd5760403660031901126100dd576101236004356101006105a6565b9061011e610119825f525f602052600160405f20015490565b610698565b610758565b005b346100dd5760203660031901126100dd5760206101486101436105bc565b610602565b6040519015158152f35b346100dd575f3660031901126100dd5760206040515f8152f35b346100dd5760403660031901126100dd576101856105a6565b6004355f525f60205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b346100dd5760203660031901126100dd576101cd6105bc565b6101d5610629565b6001600160a01b0316801561021e57805f52600160205260405f2060ff1981541690557f4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd35f80a2005b63885ce5f160e01b5f5260045ffd5b346100dd5760203660031901126100dd576102466105bc565b61024e610629565b6001600160a01b03811690811561021e573b1561031157805f52600160205260ff60405f205416156102b8575b805f52600160205260405f20600160ff198254161790557fbeceb48aeaa805aeae57be163cca6249077a18734e408a85aa74e875c43738095f80a2005b600254680100000000000000008110156102fd578060016102dc92016002556105d2565b81546001600160a01b0360039290921b91821b19169083901b17905561027b565b634e487b7160e01b5f52604160045260245ffd5b6309ee12d560e01b5f5260045ffd5b346100dd5760403660031901126100dd576103396105a6565b336001600160a01b038216036103555761012390600435610758565b63334bd91960e11b5f5260045ffd5b346100dd5760403660031901126100dd576101236004356103836105a6565b9061039c610119825f525f602052600160405f20015490565b6106d0565b346100dd5760203660031901126100dd5760206103cb6004355f525f602052600160405f20015490565b604051908152f35b346100dd575f3660031901126100dd575f6002545f915b81810361051457506103fb826105ea565b60405192601f909101601f191683019067ffffffffffffffff8211848310176102fd5761042d916040528084526105ea565b602083019190601f19013683375f905f5b818103610492578385604051918291602083019060208452518091526040830191905f5b818110610470575050500390f35b82516001600160a01b0316845285945060209384019390920191600101610462565b61049b816105d2565b60018060a01b0391549060031b1c165f52600160205260ff60405f2054166104c6575b60010161043e565b916104d0836105d2565b9054865160039290921b1c6001600160a01b03169082101561050057600582901b860160200152600101916104be565b634e487b7160e01b5f52603260045260245ffd5b61051d816105d2565b60018060a01b0391549060031b1c165f52600160205260ff60405f205416610548575b6001016103ea565b600190920191610540565b346100dd5760203660031901126100dd576004359063ffffffff60e01b82168092036100dd57602091637965db0b60e01b8114908115610595575b5015158152f35b6301ffc9a760e01b1490508361058e565b602435906001600160a01b03821682036100dd57565b600435906001600160a01b03821682036100dd57565b6002548110156105005760025f5260205f2001905f90565b67ffffffffffffffff81116102fd5760051b60200190565b6001600160a01b03168015610623575f52600160205260ff60405f20541690565b50600190565b335f9081527f5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b7602052604090205460ff161561066157565b63e2517d3f60e01b5f52336004527fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4260245260445ffd5b5f8181526020818152604080832033845290915290205460ff16156106ba5750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610752575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b5f818152602081815260408083206001600160a01b038616845290915290205460ff1615610752575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a460019056fea2646970667358221220b30cedd5e43aa3bffa6a1e0b0cc84e56307ab3001f8771e620a11dc68f483e2f64736f6c634300081e00335cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b7df8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec42",
}

// TokenAllowlist is an auto generated Go binding around an Ethereum contract.
type TokenAllowlist struct {
	abi abi.ABI
}

// NewTokenAllowlist creates a new instance of TokenAllowlist.
func NewTokenAllowlist() *TokenAllowlist {
	parsed, err := TokenAllowlistMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TokenAllowlist{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TokenAllowlist) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address admin) returns()
func (tokenAllowlist *TokenAllowlist) PackConstructor(admin common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("", admin)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (tokenAllowlist *TokenAllowlist) PackDEFAULTADMINROLE() []byte {
	enc, err := tokenAllowlist.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (tokenAllowlist *TokenAllowlist) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return tokenAllowlist.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (tokenAllowlist *TokenAllowlist) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := tokenAllowlist.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAddAllowedToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4178617f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addAllowedToken(address token) returns()
func (tokenAllowlist *TokenAllowlist) PackAddAllowedToken(token common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("addAllowedToken", token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddAllowedToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4178617f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addAllowedToken(address token) returns()
func (tokenAllowlist *TokenAllowlist) TryPackAddAllowedToken(token common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("addAllowedToken", token)
}

// PackAllowedTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe744092e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowedTokens(address ) view returns(bool)
func (tokenAllowlist *TokenAllowlist) PackAllowedTokens(arg0 common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("allowedTokens", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowedTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe744092e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowedTokens(address ) view returns(bool)
func (tokenAllowlist *TokenAllowlist) TryPackAllowedTokens(arg0 common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("allowedTokens", arg0)
}

// UnpackAllowedTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe744092e.
//
// Solidity: function allowedTokens(address ) view returns(bool)
func (tokenAllowlist *TokenAllowlist) UnpackAllowedTokens(data []byte) (bool, error) {
	out, err := tokenAllowlist.abi.Unpack("allowedTokens", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackGetAllowedTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x024ece89.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getAllowedTokens() view returns(address[])
func (tokenAllowlist *TokenAllowlist) PackGetAllowedTokens() []byte {
	enc, err := tokenAllowlist.abi.Pack("getAllowedTokens")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetAllowedTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x024ece89.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getAllowedTokens() view returns(address[])
func (tokenAllowlist *TokenAllowlist) TryPackGetAllowedTokens() ([]byte, error) {
	return tokenAllowlist.abi.Pack("getAllowedTokens")
}

// UnpackGetAllowedTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x024ece89.
//
// Solidity: function getAllowedTokens() view returns(address[])
func (tokenAllowlist *TokenAllowlist) UnpackGetAllowedTokens(data []byte) ([]common.Address, error) {
	out, err := tokenAllowlist.abi.Unpack("getAllowedTokens", data)
	if err != nil {
		return *new([]common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (tokenAllowlist *TokenAllowlist) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := tokenAllowlist.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (tokenAllowlist *TokenAllowlist) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return tokenAllowlist.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (tokenAllowlist *TokenAllowlist) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := tokenAllowlist.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (tokenAllowlist *TokenAllowlist) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (tokenAllowlist *TokenAllowlist) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (tokenAllowlist *TokenAllowlist) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (tokenAllowlist *TokenAllowlist) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (tokenAllowlist *TokenAllowlist) UnpackHasRole(data []byte) (bool, error) {
	out, err := tokenAllowlist.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsAllowedToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcbe230c3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isAllowedToken(address token) view returns(bool)
func (tokenAllowlist *TokenAllowlist) PackIsAllowedToken(token common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("isAllowedToken", token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsAllowedToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcbe230c3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isAllowedToken(address token) view returns(bool)
func (tokenAllowlist *TokenAllowlist) TryPackIsAllowedToken(token common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("isAllowedToken", token)
}

// UnpackIsAllowedToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcbe230c3.
//
// Solidity: function isAllowedToken(address token) view returns(bool)
func (tokenAllowlist *TokenAllowlist) UnpackIsAllowedToken(data []byte) (bool, error) {
	out, err := tokenAllowlist.abi.Unpack("isAllowedToken", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackRemoveAllowedToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x90469a9d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function removeAllowedToken(address token) returns()
func (tokenAllowlist *TokenAllowlist) PackRemoveAllowedToken(token common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("removeAllowedToken", token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRemoveAllowedToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x90469a9d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function removeAllowedToken(address token) returns()
func (tokenAllowlist *TokenAllowlist) TryPackRemoveAllowedToken(token common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("removeAllowedToken", token)
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (tokenAllowlist *TokenAllowlist) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("renounceRole", role, callerConfirmation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (tokenAllowlist *TokenAllowlist) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (tokenAllowlist *TokenAllowlist) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := tokenAllowlist.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (tokenAllowlist *TokenAllowlist) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return tokenAllowlist.abi.Pack("revokeRole", role, account)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (tokenAllowlist *TokenAllowlist) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := tokenAllowlist.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (tokenAllowlist *TokenAllowlist) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return tokenAllowlist.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (tokenAllowlist *TokenAllowlist) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := tokenAllowlist.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// TokenAllowlistRoleAdminChanged represents a RoleAdminChanged event raised by the TokenAllowlist contract.
type TokenAllowlistRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const TokenAllowlistRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (TokenAllowlistRoleAdminChanged) ContractEventName() string {
	return TokenAllowlistRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (tokenAllowlist *TokenAllowlist) UnpackRoleAdminChangedEvent(log *types.Log) (*TokenAllowlistRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if log.Topics[0] != tokenAllowlist.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenAllowlistRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := tokenAllowlist.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenAllowlist.abi.Events[event].Inputs {
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

// TokenAllowlistRoleGranted represents a RoleGranted event raised by the TokenAllowlist contract.
type TokenAllowlistRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const TokenAllowlistRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (TokenAllowlistRoleGranted) ContractEventName() string {
	return TokenAllowlistRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (tokenAllowlist *TokenAllowlist) UnpackRoleGrantedEvent(log *types.Log) (*TokenAllowlistRoleGranted, error) {
	event := "RoleGranted"
	if log.Topics[0] != tokenAllowlist.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenAllowlistRoleGranted)
	if len(log.Data) > 0 {
		if err := tokenAllowlist.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenAllowlist.abi.Events[event].Inputs {
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

// TokenAllowlistRoleRevoked represents a RoleRevoked event raised by the TokenAllowlist contract.
type TokenAllowlistRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const TokenAllowlistRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (TokenAllowlistRoleRevoked) ContractEventName() string {
	return TokenAllowlistRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (tokenAllowlist *TokenAllowlist) UnpackRoleRevokedEvent(log *types.Log) (*TokenAllowlistRoleRevoked, error) {
	event := "RoleRevoked"
	if log.Topics[0] != tokenAllowlist.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenAllowlistRoleRevoked)
	if len(log.Data) > 0 {
		if err := tokenAllowlist.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenAllowlist.abi.Events[event].Inputs {
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

// TokenAllowlistTokenAllowed represents a TokenAllowed event raised by the TokenAllowlist contract.
type TokenAllowlistTokenAllowed struct {
	Token common.Address
	Raw   *types.Log // Blockchain specific contextual infos
}

const TokenAllowlistTokenAllowedEventName = "TokenAllowed"

// ContractEventName returns the user-defined event name.
func (TokenAllowlistTokenAllowed) ContractEventName() string {
	return TokenAllowlistTokenAllowedEventName
}

// UnpackTokenAllowedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenAllowed(address indexed token)
func (tokenAllowlist *TokenAllowlist) UnpackTokenAllowedEvent(log *types.Log) (*TokenAllowlistTokenAllowed, error) {
	event := "TokenAllowed"
	if log.Topics[0] != tokenAllowlist.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenAllowlistTokenAllowed)
	if len(log.Data) > 0 {
		if err := tokenAllowlist.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenAllowlist.abi.Events[event].Inputs {
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

// TokenAllowlistTokenRemoved represents a TokenRemoved event raised by the TokenAllowlist contract.
type TokenAllowlistTokenRemoved struct {
	Token common.Address
	Raw   *types.Log // Blockchain specific contextual infos
}

const TokenAllowlistTokenRemovedEventName = "TokenRemoved"

// ContractEventName returns the user-defined event name.
func (TokenAllowlistTokenRemoved) ContractEventName() string {
	return TokenAllowlistTokenRemovedEventName
}

// UnpackTokenRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRemoved(address indexed token)
func (tokenAllowlist *TokenAllowlist) UnpackTokenRemovedEvent(log *types.Log) (*TokenAllowlistTokenRemoved, error) {
	event := "TokenRemoved"
	if log.Topics[0] != tokenAllowlist.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenAllowlistTokenRemoved)
	if len(log.Data) > 0 {
		if err := tokenAllowlist.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenAllowlist.abi.Events[event].Inputs {
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
func (tokenAllowlist *TokenAllowlist) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], tokenAllowlist.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return tokenAllowlist.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenAllowlist.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return tokenAllowlist.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenAllowlist.abi.Errors["NotAContract"].ID.Bytes()[:4]) {
		return tokenAllowlist.UnpackNotAContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenAllowlist.abi.Errors["TokenAddressCantBeZero"].ID.Bytes()[:4]) {
		return tokenAllowlist.UnpackTokenAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenAllowlist.abi.Errors["TokenNotAllowed"].ID.Bytes()[:4]) {
		return tokenAllowlist.UnpackTokenNotAllowedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TokenAllowlistAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the TokenAllowlist contract.
type TokenAllowlistAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func TokenAllowlistAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (tokenAllowlist *TokenAllowlist) UnpackAccessControlBadConfirmationError(raw []byte) (*TokenAllowlistAccessControlBadConfirmation, error) {
	out := new(TokenAllowlistAccessControlBadConfirmation)
	if err := tokenAllowlist.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenAllowlistAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the TokenAllowlist contract.
type TokenAllowlistAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func TokenAllowlistAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (tokenAllowlist *TokenAllowlist) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*TokenAllowlistAccessControlUnauthorizedAccount, error) {
	out := new(TokenAllowlistAccessControlUnauthorizedAccount)
	if err := tokenAllowlist.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenAllowlistNotAContract represents a NotAContract error raised by the TokenAllowlist contract.
type TokenAllowlistNotAContract struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAContract()
func TokenAllowlistNotAContractErrorID() common.Hash {
	return common.HexToHash("0x09ee12d5e890f67fbc6b54352f3db54c0ae59044f71a515db1d8d03c55e89c18")
}

// UnpackNotAContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAContract()
func (tokenAllowlist *TokenAllowlist) UnpackNotAContractError(raw []byte) (*TokenAllowlistNotAContract, error) {
	out := new(TokenAllowlistNotAContract)
	if err := tokenAllowlist.abi.UnpackIntoInterface(out, "NotAContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenAllowlistTokenAddressCantBeZero represents a TokenAddressCantBeZero error raised by the TokenAllowlist contract.
type TokenAllowlistTokenAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenAddressCantBeZero()
func TokenAllowlistTokenAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x885ce5f160f51ca1e25572e5d4b93cfa2a3dadb1e0a4e9b3fbf830831aff85d5")
}

// UnpackTokenAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenAddressCantBeZero()
func (tokenAllowlist *TokenAllowlist) UnpackTokenAddressCantBeZeroError(raw []byte) (*TokenAllowlistTokenAddressCantBeZero, error) {
	out := new(TokenAllowlistTokenAddressCantBeZero)
	if err := tokenAllowlist.abi.UnpackIntoInterface(out, "TokenAddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenAllowlistTokenNotAllowed represents a TokenNotAllowed error raised by the TokenAllowlist contract.
type TokenAllowlistTokenNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenNotAllowed()
func TokenAllowlistTokenNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xa29c498656b92b53bef33b5e62c8bb31470e81c408996e42696a26083f02fc8d")
}

// UnpackTokenNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenNotAllowed()
func (tokenAllowlist *TokenAllowlist) UnpackTokenNotAllowedError(raw []byte) (*TokenAllowlistTokenNotAllowed, error) {
	out := new(TokenAllowlistTokenNotAllowed)
	if err := tokenAllowlist.abi.UnpackIntoInterface(out, "TokenNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}
