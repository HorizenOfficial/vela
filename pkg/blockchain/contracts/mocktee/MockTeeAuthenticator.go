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

// MockTeeAuthenticatorMetaData contains all meta data concerning the MockTeeAuthenticator contract.
var MockTeeAuthenticatorMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_teeSigner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_pubSecp521r1\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"TeeIsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"activeImage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFee\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.SignatureParams\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pubSecp521r1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_activeImage\",\"type\":\"bytes32\"}],\"name\":\"setActiveImage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "21224ad20c4253145139ff6269df243105",
	Bin: "0x6080604052346100305761001a6100146101ab565b9061049b565b610022610035565b610ca36104b38239610ca390f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b5f80fd5b60018060a01b031690565b6100b5906100a1565b90565b6100c1816100ac565b036100c857565b5f80fd5b905051906100d9826100b8565b565b5f80fd5b5f80fd5b60018060401b0381116100ff576100fb60209161003f565b0190565b610049565b90825f9392825e0152565b9092919261012461011f826100e3565b610084565b938185526020850190828401116101405761013e92610104565b565b6100df565b9080601f83011215610163578160206101609351910161010f565b90565b6100db565b9190916040818403126101a657610181835f83016100cc565b92602082015160018060401b0381116101a15761019e9201610145565b90565b61009d565b610099565b6101c9611156803803806101be81610084565b928339810190610168565b9091565b5f1b90565b906101e360018060a01b03916101cd565b9181191691161790565b90565b6102046101ff610209926100a1565b6101ed565b6100a1565b90565b610215906101f0565b90565b6102219061020c565b90565b90565b9061023c61023761024392610218565b610224565b82546101d2565b9055565b5190565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561027f575b602083101461027a57565b61024b565b91607f169161026f565b5f5260205f2090565b601f602091010490565b1b90565b919060086102bb9102916102b55f198461029c565b9261029c565b9181191691161790565b90565b6102dc6102d76102e1926102c5565b6101ed565b6102c5565b90565b90565b91906102fd6102f8610305936102c8565b6102e4565b9083546102a0565b9055565b5f90565b61031f91610319610309565b916102e7565b565b5b81811061032d575050565b8061033a5f60019361030d565b01610322565b9190601f8111610350575b505050565b61035c61038193610289565b90602061036884610292565b83019310610389575b61037a90610292565b0190610321565b5f808061034b565b915061037a81929050610371565b1c90565b906103ab905f1990600802610397565b191690565b816103ba9161039b565b906002021790565b906103cc81610247565b9060018060401b03821161048a576103ee826103e8855461025f565b85610340565b602090601f831160011461042257918091610411935f92610416575b50506103b0565b90555b565b90915001515f8061040a565b601f1983169161043185610289565b925f5b81811061047257509160029391856001969410610458575b50505002019055610414565b610468910151601f84169061039b565b90555f808061044c565b91936020600181928787015181550195019201610434565b610049565b90610499916103c2565b565b906104a96104b0925f610227565b600161048f565b56fe60806040526004361015610013575b610b8f565b61001d5f3561008c565b8063081bec7e146100875780630dd7ce2f146100825780631df53a041461007d57806343f855c3146100785780635c9626b914610073578063c8a032c81461006e5763fe4993ca0361000e57610b5a565b6109e0565b610950565b610280565b61020a565b610188565b61011a565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f9103126100aa57565b61009c565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b6100f06100f96020936100fe936100e7816100af565b938480936100b3565b958691016100bc565b6100c7565b0190565b6101179160208201915f8184039101526100d1565b90565b3461014a5761012a3660046100a0565b610146610135610ba4565b61013d610092565b91829182610102565b0390f35b610098565b60018060a01b031690565b6101639061014f565b90565b61016f9061015a565b9052565b9190610186905f60208501940190610166565b565b346101b8576101983660046100a0565b6101b46101a3610be4565b6101ab610092565b91829182610173565b0390f35b610098565b5f80fd5b90565b6101cd816101c1565b036101d457565b5f80fd5b905035906101e5826101c4565b565b90602082820312610200576101fd915f016101d8565b90565b61009c565b5f0190565b346102385761022261021d3660046101e7565b610c4c565b61022a610092565b8061023481610205565b0390f35b610098565b1c90565b60018060a01b031690565b61025c906008610261930261023d565b610241565b90565b9061026f915461024c565b90565b61027d5f5f90610264565b90565b346102b0576102903660046100a0565b6102ac61029b610272565b6102a3610092565b91829182610173565b0390f35b610098565b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b906102d7906100c7565b810190811067ffffffffffffffff8211176102f157604052565b6102b9565b90610309610302610092565b92836102cd565b565b5f80fd5b67ffffffffffffffff1690565b6103258161030f565b0361032c57565b5f80fd5b9050359061033d8261031c565b565b5f80fd5b67ffffffffffffffff811161035b5760208091020190565b6102b9565b5f80fd5b5f80fd5b67ffffffffffffffff8111610386576103826020916100c7565b0190565b6102b9565b90825f939282370152565b909291926103ab6103a682610368565b6102f6565b938185526020850190828401116103c7576103c59261038b565b565b610364565b9080601f830112156103ea578160206103e793359101610396565b90565b61033f565b9291906104036103fe82610343565b6102f6565b938185526020808601920281019183831161045a5781905b838210610429575050505050565b813567ffffffffffffffff81116104555760209161044a87849387016103cc565b81520191019061041b565b61033f565b610360565b9080601f8301121561047d5781602061047a933591016103ef565b90565b61033f565b67ffffffffffffffff811161049a5760208091020190565b6102b9565b909291926104b46104af82610482565b6102f6565b93818552602080860192028301928184116104f157915b8383106104d85750505050565b602080916104e684866101d8565b8152019201916104cb565b610360565b9080601f83011215610514578160206105119335910161049f565b90565b61033f565b9190916040818403126105835761053060406102f6565b925f82013567ffffffffffffffff811161057e578161055091840161045f565b5f850152602082013567ffffffffffffffff81116105795761057292016104f6565b6020830152565b61030b565b61030b565b6102b5565b67ffffffffffffffff81116105a05760208091020190565b6102b9565b6105ae8161015a565b036105b557565b5f80fd5b905035906105c6826105a5565b565b6105d19061014f565b90565b6105dd816105c8565b036105e457565b5f80fd5b905035906105f5826105d4565b565b90565b610603816105f7565b0361060a57565b5f80fd5b9050359061061b826105fa565b565b9190606083820312610669576106629061063760606102f6565b93610644825f83016105b9565b5f86015261065582602083016105e8565b602086015260400161060e565b6040830152565b6102b5565b9092919261068361067e82610588565b6102f6565b9381855260606020860192028301928184116106c257915b8383106106a85750505050565b60206060916106b7848661061d565b81520192019161069b565b610360565b9080601f830112156106e5578160206106e29335910161066e565b90565b61033f565b600d11156106f457565b5f80fd5b90503590610705826106ea565b565b67ffffffffffffffff8111610725576107216020916100c7565b0190565b6102b9565b9092919261073f61073a82610707565b6102f6565b9381855260208501908284011161075b576107599261038b565b565b610364565b9080601f8301121561077e5781602061077b9335910161072a565b90565b61033f565b919091610160818403126108c75761079c6101606102f6565b926107a9815f8401610330565b5f8501526107ba81602084016101d8565b60208501526107cc81604084016101d8565b60408501526107de81606084016101d8565b6060850152608082013567ffffffffffffffff81116108c25781610803918401610519565b608085015260a082013567ffffffffffffffff81116108bd5781610828918401610519565b60a085015260c082013567ffffffffffffffff81116108b8578161084d9184016106c7565b60c085015261085f8160e0840161060e565b60e085015261087281610100840161060e565b6101008501526108868161012084016106f8565b61012085015261014082013567ffffffffffffffff81116108b3576108ab9201610760565b610140830152565b61030b565b61030b565b61030b565b61030b565b6102b5565b919091604081840312610924575f81013567ffffffffffffffff811161091f57836108f8918301610783565b92602082013567ffffffffffffffff811161091a5761091792016103cc565b90565b6101bd565b6101bd565b61009c565b151590565b61093790610929565b9052565b919061094e905f6020850194019061092e565b565b346109815761097d61096c6109663660046108cc565b90610c5d565b610974610092565b9182918261093b565b0390f35b610098565b90565b61099990600861099e930261023d565b610986565b90565b906109ac9154610989565b90565b6109bb60025f906109a1565b90565b6109c7906101c1565b9052565b91906109de905f602085019401906109be565b565b34610a10576109f03660046100a0565b610a0c6109fb6109af565b610a03610092565b918291826109cb565b0390f35b610098565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610a5c575b6020831014610a5757565b610a28565b91607f1691610a4c565b60209181520190565b5f5260205f2090565b905f9291805490610a92610a8b83610a3c565b8094610a66565b916001811690815f14610ae95750600114610aad575b505050565b610aba9192939450610a6f565b915f925b818410610ad157505001905f8080610aa8565b60018160209295939554848601520191019290610abe565b92949550505060ff19168252151560200201905f8080610aa8565b90610b0e91610a78565b90565b90610b31610b2a92610b21610092565b93848092610b04565b03836102cd565b565b905f10610b4657610b4390610b11565b90565b610a15565b610b5760015f90610b33565b90565b34610b8a57610b6a3660046100a0565b610b86610b75610b4b565b610b7d610092565b91829182610102565b0390f35b610098565b5f80fd5b606090565b610ba190610b11565b90565b610bac610b93565b50610bb76001610b98565b90565b5f90565b5f1c90565b610bcf610bd491610bbe565b610241565b90565b610be19054610bc3565b90565b610bec610bba565b50610bf65f610bd7565b90565b5f1b90565b90610c0a5f1991610bf9565b9181191691161790565b610c1d906101c1565b90565b610c2990610bbe565b90565b90610c41610c3c610c4892610c14565b610c20565b8254610bfe565b9055565b610c57906002610c2c565b565b5f90565b5050610c67610c59565b5060019056fea2646970667358221220a82fa25d476dfebfeffc207a8a8959caceab753231bb916de5f423c84b4c741e64736f6c634300081e0033",
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

// PackActiveImage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc8a032c8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function activeImage() view returns(bytes32)
func (mockTeeAuthenticator *MockTeeAuthenticator) PackActiveImage() []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("activeImage")
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
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackActiveImage() ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("activeImage")
}

// UnpackActiveImage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc8a032c8.
//
// Solidity: function activeImage() view returns(bytes32)
func (mockTeeAuthenticator *MockTeeAuthenticator) UnpackActiveImage(data []byte) ([32]byte, error) {
	out, err := mockTeeAuthenticator.abi.Unpack("activeImage", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
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

// PackSetActiveImage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1df53a04.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setActiveImage(bytes32 _activeImage) returns()
func (mockTeeAuthenticator *MockTeeAuthenticator) PackSetActiveImage(activeImage [32]byte) []byte {
	enc, err := mockTeeAuthenticator.abi.Pack("setActiveImage", activeImage)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetActiveImage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1df53a04.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setActiveImage(bytes32 _activeImage) returns()
func (mockTeeAuthenticator *MockTeeAuthenticator) TryPackSetActiveImage(activeImage [32]byte) ([]byte, error) {
	return mockTeeAuthenticator.abi.Pack("setActiveImage", activeImage)
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
	Bin: "0x608060405234601957600e601d565b603e60288239603e90f35b6023565b60405190565b5f80fdfe60806040525f80fdfea2646970667358221220f702a542523fac4b8b41a4856f6ce9505d3eb0c9dd65e004f312bedc39c9f40764736f6c634300081e0033",
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
