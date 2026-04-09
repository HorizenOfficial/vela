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

// StructsSignatureParams is an auto generated low-level Go binding around an user-defined struct.
type StructsSignatureParams struct {
	ApplicationId      uint64
	PrevStateRoot      [32]byte
	NewStateRoot       [32]byte
	ProcessedRequestId [32]byte
	Events             [][]byte
	EventSubTypes      []string
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
	ABI: "[{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// the contract method with ID 0xd119a29d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) PackCheckSignature(params StructsSignatureParams, signature []byte) []byte {
	enc, err := iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd119a29d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
func (iTeeAuthenticator *ITeeAuthenticator) TryPackCheckSignature(params StructsSignatureParams, signature []byte) ([]byte, error) {
	return iTeeAuthenticator.abi.Pack("checkSignature", params, signature)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd119a29d.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) params, bytes signature) view returns(bool)
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

// MockTeeAuthenticatorMetaData contains all meta data concerning the MockTeeAuthenticator contract.
var MockTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_teeSigner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pubSecp521r1\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "21224ad20c4253145139ff6269df243105",
	Bin: "0x6080604052346100305761001a6100146101ab565b9061049b565b610022610035565b610ae86104b38239610ae890f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b5f80fd5b60018060a01b031690565b6100b5906100a1565b90565b6100c1816100ac565b036100c857565b5f80fd5b905051906100d9826100b8565b565b5f80fd5b5f80fd5b60018060401b0381116100ff576100fb60209161003f565b0190565b610049565b90825f9392825e0152565b9092919261012461011f826100e3565b610084565b938185526020850190828401116101405761013e92610104565b565b6100df565b9080601f83011215610163578160206101609351910161010f565b90565b6100db565b9190916040818403126101a657610181835f83016100cc565b92602082015160018060401b0381116101a15761019e9201610145565b90565b61009d565b610099565b6101c9610f9b803803806101be81610084565b928339810190610168565b9091565b5f1b90565b906101e360018060a01b03916101cd565b9181191691161790565b90565b6102046101ff610209926100a1565b6101ed565b6100a1565b90565b610215906101f0565b90565b6102219061020c565b90565b90565b9061023c61023761024392610218565b610224565b82546101d2565b9055565b5190565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561027f575b602083101461027a57565b61024b565b91607f169161026f565b5f5260205f2090565b601f602091010490565b1b90565b919060086102bb9102916102b55f198461029c565b9261029c565b9181191691161790565b90565b6102dc6102d76102e1926102c5565b6101ed565b6102c5565b90565b90565b91906102fd6102f8610305936102c8565b6102e4565b9083546102a0565b9055565b5f90565b61031f91610319610309565b916102e7565b565b5b81811061032d575050565b8061033a5f60019361030d565b01610322565b9190601f8111610350575b505050565b61035c61038193610289565b90602061036884610292565b83019310610389575b61037a90610292565b0190610321565b5f808061034b565b915061037a81929050610371565b1c90565b906103ab905f1990600802610397565b191690565b816103ba9161039b565b906002021790565b906103cc81610247565b9060018060401b03821161048a576103ee826103e8855461025f565b85610340565b602090601f831160011461042257918091610411935f92610416575b50506103b0565b90555b565b90915001515f8061040a565b601f1983169161043185610289565b925f5b81811061047257509160029391856001969410610458575b50505002019055610414565b610468910151601f84169061039b565b90555f808061044c565b91936020600181928787015181550195019201610434565b610049565b90610499916103c2565b565b906104a96104b0925f610227565b600161048f565b56fe60806040526004361015610013575b610a34565b61001d5f3561006c565b8063081bec7e146100675780630dd7ce2f1461006257806343f855c31461005d578063d119a29d146100585763fe4993ca0361000e576109ff565b610884565b6101e0565b610168565b6100fa565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f91031261008a57565b61007c565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b6100d06100d96020936100de936100c78161008f565b93848093610093565b9586910161009c565b6100a7565b0190565b6100f79160208201915f8184039101526100b1565b90565b3461012a5761010a366004610080565b610126610115610a49565b61011d610072565b918291826100e2565b0390f35b610078565b60018060a01b031690565b6101439061012f565b90565b61014f9061013a565b9052565b9190610166905f60208501940190610146565b565b3461019857610178366004610080565b610194610183610a89565b61018b610072565b91829182610153565b0390f35b610078565b1c90565b60018060a01b031690565b6101bc9060086101c1930261019d565b6101a1565b90565b906101cf91546101ac565b90565b6101dd5f5f906101c4565b90565b34610210576101f0366004610080565b61020c6101fb6101d2565b610203610072565b91829182610153565b0390f35b610078565b5f80fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b9061023b906100a7565b810190811067ffffffffffffffff82111761025557604052565b61021d565b9061026d610266610072565b9283610231565b565b5f80fd5b67ffffffffffffffff1690565b61028981610273565b0361029057565b5f80fd5b905035906102a182610280565b565b90565b6102af816102a3565b036102b657565b5f80fd5b905035906102c7826102a6565b565b5f80fd5b67ffffffffffffffff81116102e55760208091020190565b61021d565b5f80fd5b5f80fd5b67ffffffffffffffff81116103105761030c6020916100a7565b0190565b61021d565b90825f939282370152565b90929192610335610330826102f2565b61025a565b938185526020850190828401116103515761034f92610315565b565b6102ee565b9080601f830112156103745781602061037193359101610320565b90565b6102c9565b92919061038d610388826102cd565b61025a565b93818552602080860192028101918383116103e45781905b8382106103b3575050505050565b813567ffffffffffffffff81116103df576020916103d48784938701610356565b8152019101906103a5565b6102c9565b6102ea565b9080601f830112156104075781602061040493359101610379565b90565b6102c9565b67ffffffffffffffff81116104245760208091020190565b61021d565b67ffffffffffffffff8111610447576104436020916100a7565b0190565b61021d565b9092919261046161045c82610429565b61025a565b9381855260208501908284011161047d5761047b92610315565b565b6102ee565b9080601f830112156104a05781602061049d9335910161044c565b90565b6102c9565b9291906104b96104b48261040c565b61025a565b93818552602080860192028101918383116105105781905b8382106104df575050505050565b813567ffffffffffffffff811161050b576020916105008784938701610482565b8152019101906104d1565b6102c9565b6102ea565b9080601f8301121561053357816020610530933591016104a5565b90565b6102c9565b67ffffffffffffffff81116105505760208091020190565b61021d565b61055e8161013a565b0361056557565b5f80fd5b9050359061057682610555565b565b6105819061012f565b90565b61058d81610578565b0361059457565b5f80fd5b905035906105a582610584565b565b90565b6105b3816105a7565b036105ba57565b5f80fd5b905035906105cb826105aa565b565b919060608382031261061957610612906105e7606061025a565b936105f4825f8301610569565b5f8601526106058260208301610598565b60208601526040016105be565b6040830152565b610219565b9092919261063361062e82610538565b61025a565b93818552606060208601920283019281841161067257915b8383106106585750505050565b602060609161066784866105cd565b81520192019161064b565b6102ea565b9080601f83011215610695578160206106929335910161061e565b90565b6102c9565b600d11156106a457565b5f80fd5b905035906106b58261069a565b565b919091610160818403126107fb576106d061016061025a565b926106dd815f8401610294565b5f8501526106ee81602084016102ba565b602085015261070081604084016102ba565b604085015261071281606084016102ba565b6060850152608082013567ffffffffffffffff81116107f657816107379184016103e9565b608085015260a082013567ffffffffffffffff81116107f1578161075c918401610515565b60a085015260c082013567ffffffffffffffff81116107ec5781610781918401610677565b60c08501526107938160e084016105be565b60e08501526107a68161010084016105be565b6101008501526107ba8161012084016106a8565b61012085015261014082013567ffffffffffffffff81116107e7576107df9201610482565b610140830152565b61026f565b61026f565b61026f565b61026f565b610219565b919091604081840312610858575f81013567ffffffffffffffff8111610853578361082c9183016106b7565b92602082013567ffffffffffffffff811161084e5761084b9201610356565b90565b610215565b610215565b61007c565b151590565b61086b9061085d565b9052565b9190610882905f60208501940190610862565b565b346108b5576108b16108a061089a366004610800565b90610aa2565b6108a8610072565b9182918261086f565b0390f35b610078565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610901575b60208310146108fc57565b6108cd565b91607f16916108f1565b60209181520190565b5f5260205f2090565b905f9291805490610937610930836108e1565b809461090b565b916001811690815f1461098e5750600114610952575b505050565b61095f9192939450610914565b915f925b81841061097657505001905f808061094d565b60018160209295939554848601520191019290610963565b92949550505060ff19168252151560200201905f808061094d565b906109b39161091d565b90565b906109d66109cf926109c6610072565b938480926109a9565b0383610231565b565b905f106109eb576109e8906109b6565b90565b6108ba565b6109fc60015f906109d8565b90565b34610a2f57610a0f366004610080565b610a2b610a1a6109f0565b610a22610072565b918291826100e2565b0390f35b610078565b5f80fd5b606090565b610a46906109b6565b90565b610a51610a38565b50610a5c6001610a3d565b90565b5f90565b5f1c90565b610a74610a7991610a63565b6101a1565b90565b610a869054610a68565b90565b610a91610a5f565b50610a9b5f610a7c565b90565b5f90565b5050610aac610a9e565b5060019056fea26469706673582212200480d04a4b6a1f10c5b3ea35b1eb4622f9f18151c23749d58ab1ab610415616164736f6c634300081e0033",
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
// the contract method with ID 0xd119a29d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) PackCheckSignature(arg0 StructsSignatureParams, arg1 []byte) []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("checkSignature", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd119a29d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) , bytes ) pure returns(bool)
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackCheckSignature(arg0 StructsSignatureParams, arg1 []byte) ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("checkSignature", arg0, arg1)
}

// UnpackCheckSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd119a29d.
//
// Solidity: function checkSignature((uint64,bytes32,bytes32,bytes32,bytes[],string[],(address,address,uint256)[],uint256,uint256,uint8,string) , bytes ) pure returns(bool)
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
	if bytes.Equal(raw[:4], mockTeeAuthenticator.abi.Errors["TeeIsNotSet"].ID.Bytes()[:4]) {
		return mockTeeAuthenticator.UnpackTeeIsNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
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
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea26469706673582212209ffbba52071cd4909dd91df41f3e001548c66a43f39312767f50526acb30e24064736f6c634300081e0033",
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
