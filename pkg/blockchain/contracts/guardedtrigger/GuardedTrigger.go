// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package guardedtrigger

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

// StructsTokenAndAmount is an auto generated low-level Go binding around an user-defined struct.
type StructsTokenAndAmount struct {
	Token  common.Address
	Amount *big.Int
}

// GuardedTriggerMetaData contains all meta data concerning the GuardedTrigger contract.
var GuardedTriggerMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIProcessorEndpoint\",\"name\":\"_processorEndpoint\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_sink\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"NotProcessorEndpoint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"executeSuccess\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"withdrawSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"returnedTokens\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"failedTokens\",\"type\":\"tuple[]\"}],\"name\":\"getTrustProcessPayload\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"processorEndpoint\",\"outputs\":[{\"internalType\":\"contractIProcessorEndpoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sink\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "GuardedTrigger",
	Bin: "0x6080604052346100305761001a610014610133565b906101ff565b610022610035565b6112ef61040b82396112ef90f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b60018060a01b031690565b6100b19061009d565b90565b6100bd906100a8565b90565b6100c9816100b4565b036100d057565b5f80fd5b905051906100e1826100c0565b565b6100ec816100a8565b036100f357565b5f80fd5b90505190610104826100e3565b565b919060408382031261012e578061012261012b925f86016100d4565b936020016100f7565b90565b610099565b6101516116fa8038038061014681610084565b928339810190610106565b9091565b90565b90565b61016f61016a61017492610155565b610158565b61009d565b90565b6101809061015b565b90565b5f0190565b5f1b90565b9061019e60018060a01b0391610188565b9181191691161790565b6101bc6101b76101c19261009d565b610158565b61009d565b90565b6101cd906101a8565b90565b6101d9906101c4565b90565b90565b906101f46101ef6101fb926101d0565b6101dc565b825461018d565b9055565b61020890610336565b8061022361021d6102185f610177565b6100a8565b916100a8565b14610235576102339060026101df565b565b5f63d92e233d60e01b81528061024d60048201610183565b0390fd5b61025a906101c4565b90565b610266906101a8565b90565b6102729061025d565b90565b90565b9061028d61028861029492610269565b610275565b825461018d565b9055565b60e01b90565b6102a7906100a8565b90565b6102b38161029e565b036102ba57565b5f80fd5b905051906102cb826102aa565b565b906020828203126102e6576102e3915f016102be565b90565b610099565b6102f3610035565b3d5f823e3d90fd5b610304906101a8565b90565b610310906102fb565b90565b90565b9061032b61032661033292610307565b610313565b825461018d565b9055565b61033f81610251565b61035961035361034e5f610177565b6100a8565b916100a8565b146103ee5760206103768261037161038c945f610278565b610251565b63f624f88a90610384610035565b938492610298565b8252818061039c60048201610183565b03915afa80156103e9576103b9915f916103bb575b506001610316565b565b6103dc915060203d81116103e2575b6103d4818361005d565b8101906102cd565b5f6103b1565b503d6103ca565b6102eb565b5f63d92e233d60e01b81528061040660048201610183565b0390fdfe60806040526004361015610015575b3661062257005b61001f5f3561007e565b80633ccfd60b1461007957806371b3029a14610074578063a40aad2a1461006f578063c74e820e1461006a578063f624f88a146100655763fafb0c6a0361000e576105ed565b61054b565b61047d565b6103db565b610229565b61019e565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f91031261009c57565b61008e565b5190565b60209181520190565b60200190565b60018060a01b031690565b6100c8906100b4565b90565b6100d4906100bf565b9052565b90565b6100e4906100d8565b9052565b9060208061010a936101005f8201515f8601906100cb565b01519101906100db565b565b90610119816040936100e8565b0190565b60200190565b9061014061013a610133846100a1565b80936100a5565b926100ae565b905f5b8181106101505750505090565b909192610169610163600192865161010c565b9461011d565b9101919091610143565b909161018d61019b9360408401908482035f860152610123565b916020818403910152610123565b90565b346101cf576101ae366004610092565b6101b6610ed6565b906101cb6101c2610084565b92839283610173565b0390f35b61008a565b5f80fd5b5f80fd5b908160409103126101ea5790565b6101d8565b9060208282031261021f575f82013567ffffffffffffffff811161021a5761021792016101dc565b90565b6101d4565b61008e565b5f0190565b346102575761024161023c3660046101ef565b610f83565b610249610084565b8061025381610224565b0390f35b61008a565b151590565b61026a8161025c565b0361027157565b5f80fd5b9050359061028282610261565b565b5f80fd5b5f80fd5b5f80fd5b909182601f830112156102ca5781359167ffffffffffffffff83116102c55760200192604083028401116102c057565b61028c565b610288565b610284565b909160a08284031261036b575f82013567ffffffffffffffff811161036657836102fa9184016101dc565b926103088160208501610275565b926103168260408301610275565b92606082013567ffffffffffffffff81116103615783610337918401610290565b929093608082013567ffffffffffffffff811161035c576103589201610290565b9091565b6101d4565b6101d4565b6101d4565b61008e565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b6103b16103ba6020936103bf936103a881610370565b93848093610374565b9586910161037d565b610388565b0190565b6103d89160208201915f818403910152610392565b90565b346104125761040e6103fd6103f13660046102cf565b95949094939193611006565b610405610084565b918291826103c3565b0390f35b61008a565b1c90565b60018060a01b031690565b61043690600861043b9302610417565b61041b565b90565b906104499154610426565b90565b61045860025f9061043e565b90565b610464906100bf565b9052565b919061047b905f6020850194019061045b565b565b346104ad5761048d366004610092565b6104a961049861044c565b6104a0610084565b91829182610468565b0390f35b61008a565b60018060a01b031690565b6104cd9060086104d29302610417565b6104b2565b90565b906104e091546104bd565b90565b6104ef60015f906104d5565b90565b90565b61050961050461050e926100b4565b6104f2565b6100b4565b90565b61051a906104f5565b90565b61052690610511565b90565b6105329061051d565b9052565b9190610549905f60208501940190610529565b565b3461057b5761055b366004610092565b6105776105666104e3565b61056e610084565b91829182610536565b0390f35b61008a565b60018060a01b031690565b61059b9060086105a09302610417565b610580565b90565b906105ae915461058b565b90565b6105bc5f5f906105a3565b90565b6105c890610511565b90565b6105d4906105bf565b9052565b91906105eb905f602085019401906105cb565b565b3461061d576105fd366004610092565b6106196106086105b1565b610610610084565b918291826105d8565b0390f35b61008a565b5f80fd5b606090565b5f1c90565b61063c6106419161062b565b610580565b90565b61064e9054610630565b90565b903361067561066f61066a6106655f610644565b6105bf565b6100bf565b916100bf565b036106875761068391610b1d565b9091565b5f63c14b0f1760e01b81528061069f60048201610224565b0390fd5b6106af6106b49161062b565b6104b2565b90565b6106c190546106a3565b90565b634e487b7160e01b5f52604160045260245ffd5b906106e290610388565b810190811067ffffffffffffffff8211176106fc57604052565b6106c4565b60e01b90565b9061071a610713610084565b92836106d8565b565b67ffffffffffffffff81116107345760208091020190565b6106c4565b610742816100bf565b0361074957565b5f80fd5b9050519061075a82610739565b565b9092919261077161076c8261071c565b610707565b93818552602080860192028301928184116107ae57915b8383106107955750505050565b602080916107a3848661074d565b815201920191610788565b61028c565b9080601f830112156107d1578160206107ce9351910161075c565b90565b610284565b90602082820312610806575f82015167ffffffffffffffff8111610801576107fe92016107b3565b90565b6101d4565b61008e565b610813610084565b3d5f823e3d90fd5b5190565b90565b61083661083161083b9261081f565b6104f2565b6100d8565b90565b634e487b7160e01b5f52601160045260245ffd5b610861610867919392936100d8565b926100d8565b820180921161087257565b61083e565b67ffffffffffffffff811161088f5760208091020190565b6106c4565b906108a66108a183610877565b610707565b918252565b6108b56040610707565b90565b5f90565b5f90565b6108c86108ab565b90602080836108d56108b8565b8152016108e06108bc565b81525050565b6108ee6108c0565b90565b5f5b8281106108ff57505050565b60209061090a6108e6565b81840152016108f3565b9061093961092183610894565b9260208061092f8693610877565b92019103906108f1565b565b5f90565b634e487b7160e01b5f52603260045260245ffd5b9061095d8261081b565b81101561096e576020809102010190565b61093f565b61097d90516100bf565b90565b610989906104f5565b90565b61099590610980565b90565b6109a190610511565b90565b6109ad90610511565b90565b6109b9816100d8565b036109c057565b5f80fd5b905051906109d1826109b0565b565b906020828203126109ec576109e9915f016109c4565b90565b61008e565b90565b610a08610a03610a0d926109f1565b6104f2565b6100d8565b90565b90505190610a1d82610261565b565b90602082820312610a3857610a35915f01610a10565b90565b61008e565b610a46906100d8565b9052565b916020610a6b929493610a6460408201965f83019061045b565b0190610a3d565b565b610a776040610707565b90565b90610a84906100bf565b9052565b90610a92906100d8565b9052565b90610aa0826100a1565b811015610ab1576020809102010190565b61093f565b6001610ac291016100d8565b90565b610ace906104f5565b90565b610ada90610ac5565b90565b610ae690610511565b90565b610afd610af8610b02926109f1565b6104f2565b6100b4565b90565b610b0e90610ae9565b90565b610b1a5f610b05565b90565b5050610b4b5f610b35610b3060016106b7565b61051d565b63024ece8990610b43610084565b938492610701565b82528180610b5b60048201610224565b03915afa908115610ed1575f91610eaf575b5090610b788261081b565b90610b95610b9083610b8a6001610822565b90610852565b610914565b610bb1610bac84610ba66001610822565b90610852565b610914565b91610bba61093b565b91610bc361093b565b91610bcc61093b565b935b80610be1610bdb896100d8565b916100d8565b14610df657610bf9610bf4898390610953565b610973565b610c416020610c0f610c0a8461098c565b610998565b6370a0823190610c36610c21306109a4565b92610c2a610084565b95869485938493610701565b835260048301610468565b03915afa908115610df1575f91610dc3575b509081610c68610c625f6109f4565b916100d8565b11610c7e575b5050610c7990610ab6565b610bce565b610c8f610c8a8261098c565b610998565b602063a9059cbb91610ca8610ca35f610644565b6105bf565b90610cc65f8795610cd1610cba610084565b97889687958694610701565b845260048401610a4a565b03925af19081610d97575b50155f14610d49576001610cf9575b5050610c79905b905f610c6e565b95610d3a610d4192610d27610c79959991610d1e610d15610a6d565b935f8501610a7a565b60208301610a88565b898391610d348383610a96565b52610a96565b5150610ab6565b94905f610ceb565b94610d8a610d9192610d77610c79959891610d6e610d65610a6d565b935f8501610a7a565b60208301610a88565b868391610d848383610a96565b52610a96565b5150610ab6565b93610cf2565b610db79060203d8111610dbc575b610daf81836106d8565b810190610a1f565b610cdc565b503d610da5565b610de4915060203d8111610dea575b610ddc81836106d8565b8101906109d3565b5f610c53565b503d610dd2565b61080b565b509450929450610e05306109a4565b319081610e1a610e145f6109f4565b916100d8565b11610e6d575b610e62600192610e4f610e31610b11565b91610e46610e3d610a6d565b935f8501610a7a565b60208301610a88565b858391610e5c8383610a96565b52610a96565b515001825283529190565b5f808080610e92610e8d610e88610e8384610644565b6105bf565b610ad1565b610add565b8690828215610ea6575bf1610e205761080b565b506108fc610e9c565b610ecb91503d805f833e610ec381836106d8565b8101906107d6565b5f610b6d565b61080b565b610eef610ee1610626565b610ee9610626565b90610651565b9091565b33610f16610f10610f0b610f065f610644565b6105bf565b6100bf565b916100bf565b03610f2657610f2490610f42565b565b5f63c14b0f1760e01b815280610f3e60048201610224565b0390fd5b610f4b90611072565b7fcc0dfa2e3d45df1c774de9e4a88cc12a14bab7c797b8cf7a992f8b2a6f2f4bda610f74610084565b80610f7e81610224565b0390a1565b610f8c90610ef3565b565b606090565b9695949392919033610fbd610fb7610fb2610fad5f610644565b6105bf565b6100bf565b916100bf565b03610fce57610fcb97610fea565b90565b5f63c14b0f1760e01b815280610fe660048201610224565b0390fd5b90611003979593919694929650959091929394956111d3565b90565b9061101d969594939291611018610f8e565b610f93565b90565b61102990610511565b90565b61103861103d9161062b565b61041b565b90565b61104a905461102c565b90565b61105c611062919392936100d8565b926100d8565b820391821161106d57565b61083e565b5061107c30611020565b318061109161108b6001610822565b916100d8565b1161109a575b50565b5f808080936110cd6110bc6110b76110b26002611040565b610ad1565b610add565b916110c76001610822565b9061104d565b908282156110e7575bf1156110e2575f611097565b61080b565b506108fc6110d6565b5f80fd5b5f80fd5b5f80fd5b90359060016020038136030382121561113e570180359067ffffffffffffffff82116111395760200191602082023603831361113457565b6110f8565b6110f4565b6110f0565b5090565b67ffffffffffffffff811161116557611161602091610388565b0190565b6106c4565b9061117c61117783611147565b610707565b918252565b61118a5f61116a565b90565b611195611181565b90565b5090565b91908110156111ac576040020190565b61093f565b356111bb816109b0565b90565b91906111d1905f60208501940190610a3d565b565b9150506111f49194506111fa9293506111ea610f8e565b505f8101906110fc565b90611143565b61120c6112065f6109f4565b916100d8565b146112ac5761125c9161126b916112225f6109f4565b9161122e818390611198565b61124061123a5f6109f4565b916100d8565b1161126e575b5050611250610084565b928391602083016111be565b602082018103825203826106d8565b90565b6112a592506020918161129961128961129f94938093611198565b6112936001610822565b9061104d565b9161119c565b016111b1565b5f80611246565b50506112b661118d565b9056fea26469706673582212207dadd013a2b59bbda024278665e15ee88d59a05e1809c87f9c37f2d5e02e481764736f6c634300081e0033",
}

// GuardedTrigger is an auto generated Go binding around an Ethereum contract.
type GuardedTrigger struct {
	abi abi.ABI
}

// NewGuardedTrigger creates a new instance of GuardedTrigger.
func NewGuardedTrigger() *GuardedTrigger {
	parsed, err := GuardedTriggerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &GuardedTrigger{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *GuardedTrigger) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _processorEndpoint, address _sink) returns()
func (guardedTrigger *GuardedTrigger) PackConstructor(_processorEndpoint common.Address, _sink common.Address) []byte {
	enc, err := guardedTrigger.abi.Pack("", _processorEndpoint, _sink)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackExecute is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71b3029a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function execute((bytes[],bytes32[]) appEventData) returns()
func (guardedTrigger *GuardedTrigger) PackExecute(appEventData StructsEventData) []byte {
	enc, err := guardedTrigger.abi.Pack("execute", appEventData)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExecute is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71b3029a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function execute((bytes[],bytes32[]) appEventData) returns()
func (guardedTrigger *GuardedTrigger) TryPackExecute(appEventData StructsEventData) ([]byte, error) {
	return guardedTrigger.abi.Pack("execute", appEventData)
}

// PackGetTrustProcessPayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa40aad2a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getTrustProcessPayload((bytes[],bytes32[]) appEventData, bool executeSuccess, bool withdrawSuccess, (address,uint256)[] returnedTokens, (address,uint256)[] failedTokens) returns(bytes)
func (guardedTrigger *GuardedTrigger) PackGetTrustProcessPayload(appEventData StructsEventData, executeSuccess bool, withdrawSuccess bool, returnedTokens []StructsTokenAndAmount, failedTokens []StructsTokenAndAmount) []byte {
	enc, err := guardedTrigger.abi.Pack("getTrustProcessPayload", appEventData, executeSuccess, withdrawSuccess, returnedTokens, failedTokens)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetTrustProcessPayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa40aad2a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getTrustProcessPayload((bytes[],bytes32[]) appEventData, bool executeSuccess, bool withdrawSuccess, (address,uint256)[] returnedTokens, (address,uint256)[] failedTokens) returns(bytes)
func (guardedTrigger *GuardedTrigger) TryPackGetTrustProcessPayload(appEventData StructsEventData, executeSuccess bool, withdrawSuccess bool, returnedTokens []StructsTokenAndAmount, failedTokens []StructsTokenAndAmount) ([]byte, error) {
	return guardedTrigger.abi.Pack("getTrustProcessPayload", appEventData, executeSuccess, withdrawSuccess, returnedTokens, failedTokens)
}

// UnpackGetTrustProcessPayload is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa40aad2a.
//
// Solidity: function getTrustProcessPayload((bytes[],bytes32[]) appEventData, bool executeSuccess, bool withdrawSuccess, (address,uint256)[] returnedTokens, (address,uint256)[] failedTokens) returns(bytes)
func (guardedTrigger *GuardedTrigger) UnpackGetTrustProcessPayload(data []byte) ([]byte, error) {
	out, err := guardedTrigger.abi.Unpack("getTrustProcessPayload", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackProcessorEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfafb0c6a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function processorEndpoint() view returns(address)
func (guardedTrigger *GuardedTrigger) PackProcessorEndpoint() []byte {
	enc, err := guardedTrigger.abi.Pack("processorEndpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProcessorEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfafb0c6a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function processorEndpoint() view returns(address)
func (guardedTrigger *GuardedTrigger) TryPackProcessorEndpoint() ([]byte, error) {
	return guardedTrigger.abi.Pack("processorEndpoint")
}

// UnpackProcessorEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfafb0c6a.
//
// Solidity: function processorEndpoint() view returns(address)
func (guardedTrigger *GuardedTrigger) UnpackProcessorEndpoint(data []byte) (common.Address, error) {
	out, err := guardedTrigger.abi.Unpack("processorEndpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackSink is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc74e820e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sink() view returns(address)
func (guardedTrigger *GuardedTrigger) PackSink() []byte {
	enc, err := guardedTrigger.abi.Pack("sink")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSink is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc74e820e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sink() view returns(address)
func (guardedTrigger *GuardedTrigger) TryPackSink() ([]byte, error) {
	return guardedTrigger.abi.Pack("sink")
}

// UnpackSink is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc74e820e.
//
// Solidity: function sink() view returns(address)
func (guardedTrigger *GuardedTrigger) UnpackSink(data []byte) (common.Address, error) {
	out, err := guardedTrigger.abi.Unpack("sink", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTokenAllowlist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf624f88a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tokenAllowlist() view returns(address)
func (guardedTrigger *GuardedTrigger) PackTokenAllowlist() []byte {
	enc, err := guardedTrigger.abi.Pack("tokenAllowlist")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTokenAllowlist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf624f88a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tokenAllowlist() view returns(address)
func (guardedTrigger *GuardedTrigger) TryPackTokenAllowlist() ([]byte, error) {
	return guardedTrigger.abi.Pack("tokenAllowlist")
}

// UnpackTokenAllowlist is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf624f88a.
//
// Solidity: function tokenAllowlist() view returns(address)
func (guardedTrigger *GuardedTrigger) UnpackTokenAllowlist(data []byte) (common.Address, error) {
	out, err := guardedTrigger.abi.Unpack("tokenAllowlist", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ccfd60b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw() returns((address,uint256)[], (address,uint256)[])
func (guardedTrigger *GuardedTrigger) PackWithdraw() []byte {
	enc, err := guardedTrigger.abi.Pack("withdraw")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ccfd60b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw() returns((address,uint256)[], (address,uint256)[])
func (guardedTrigger *GuardedTrigger) TryPackWithdraw() ([]byte, error) {
	return guardedTrigger.abi.Pack("withdraw")
}

// WithdrawOutput serves as a container for the return parameters of contract
// method Withdraw.
type WithdrawOutput struct {
	Arg0 []StructsTokenAndAmount
	Arg1 []StructsTokenAndAmount
}

// UnpackWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3ccfd60b.
//
// Solidity: function withdraw() returns((address,uint256)[], (address,uint256)[])
func (guardedTrigger *GuardedTrigger) UnpackWithdraw(data []byte) (WithdrawOutput, error) {
	out, err := guardedTrigger.abi.Unpack("withdraw", data)
	outstruct := new(WithdrawOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]StructsTokenAndAmount)).(*[]StructsTokenAndAmount)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]StructsTokenAndAmount)).(*[]StructsTokenAndAmount)
	return *outstruct, nil
}

// GuardedTriggerTriggerExecuted represents a TriggerExecuted event raised by the GuardedTrigger contract.
type GuardedTriggerTriggerExecuted struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const GuardedTriggerTriggerExecutedEventName = "TriggerExecuted"

// ContractEventName returns the user-defined event name.
func (GuardedTriggerTriggerExecuted) ContractEventName() string {
	return GuardedTriggerTriggerExecutedEventName
}

// UnpackTriggerExecutedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TriggerExecuted()
func (guardedTrigger *GuardedTrigger) UnpackTriggerExecutedEvent(log *types.Log) (*GuardedTriggerTriggerExecuted, error) {
	event := "TriggerExecuted"
	if log.Topics[0] != guardedTrigger.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(GuardedTriggerTriggerExecuted)
	if len(log.Data) > 0 {
		if err := guardedTrigger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range guardedTrigger.abi.Events[event].Inputs {
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
func (guardedTrigger *GuardedTrigger) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], guardedTrigger.abi.Errors["NotProcessorEndpoint"].ID.Bytes()[:4]) {
		return guardedTrigger.UnpackNotProcessorEndpointError(raw[4:])
	}
	if bytes.Equal(raw[:4], guardedTrigger.abi.Errors["ZeroAddress"].ID.Bytes()[:4]) {
		return guardedTrigger.UnpackZeroAddressError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// GuardedTriggerNotProcessorEndpoint represents a NotProcessorEndpoint error raised by the GuardedTrigger contract.
type GuardedTriggerNotProcessorEndpoint struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotProcessorEndpoint()
func GuardedTriggerNotProcessorEndpointErrorID() common.Hash {
	return common.HexToHash("0xc14b0f1759c7c9c12c65a360cebd0624cef647f19809e165e59b302bfc46ca6e")
}

// UnpackNotProcessorEndpointError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotProcessorEndpoint()
func (guardedTrigger *GuardedTrigger) UnpackNotProcessorEndpointError(raw []byte) (*GuardedTriggerNotProcessorEndpoint, error) {
	out := new(GuardedTriggerNotProcessorEndpoint)
	if err := guardedTrigger.abi.UnpackIntoInterface(out, "NotProcessorEndpoint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GuardedTriggerZeroAddress represents a ZeroAddress error raised by the GuardedTrigger contract.
type GuardedTriggerZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAddress()
func GuardedTriggerZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xd92e233df2717d4a40030e20904abd27b68fcbeede117eaaccbbdac9618c8c73")
}

// UnpackZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAddress()
func (guardedTrigger *GuardedTrigger) UnpackZeroAddressError(raw []byte) (*GuardedTriggerZeroAddress, error) {
	out := new(GuardedTriggerZeroAddress)
	if err := guardedTrigger.abi.UnpackIntoInterface(out, "ZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}
