// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testtrigger

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

// TestTriggerMetaData contains all meta data concerning the TestTrigger contract.
var TestTriggerMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIProcessorEndpoint\",\"name\":\"_processorEndpoint\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_revertOnExecute\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_revertOnPostWithdraw\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ExecuteReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotProcessorEndpoint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PostWithdrawReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"capturedBalances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"executedInBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"executedPostWithdrawsInBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"executeSuccess\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"withdrawSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"returnedTokens\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"failedTokens\",\"type\":\"tuple[]\"}],\"name\":\"getTrustProcessPayload\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"processorEndpoint\",\"outputs\":[{\"internalType\":\"contractIProcessorEndpoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"revertOnExecute\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"revertOnPostWithdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"setTrustedPayload\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedPayload\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "TestTrigger",
	Bin: "0x6080604052346100305761001a610014610145565b916101f5565b610022610035565b611a3161044b8239611a3190f35b61003b565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100679061003f565b810190811060018060401b0382111761007f57604052565b610049565b90610097610090610035565b928361005d565b565b5f80fd5b60018060a01b031690565b6100b19061009d565b90565b6100bd906100a8565b90565b6100c9816100b4565b036100d057565b5f80fd5b905051906100e1826100c0565b565b151590565b6100f1816100e3565b036100f857565b5f80fd5b90505190610109826100e8565b565b90916060828403126101405761013d610126845f85016100d4565b9361013481602086016100fc565b936040016100fc565b90565b610099565b610163611e7c8038038061015881610084565b92833981019061010b565b909192565b60a01b90565b9061017d60ff60a01b91610168565b9181191691161790565b610190906100e3565b90565b90565b906101ab6101a66101b292610187565b610193565b825461016e565b9055565b60a81b90565b906101cb60ff60a81b916101b6565b9181191691161790565b906101ea6101e56101f192610187565b610193565b82546101bc565b9055565b610214929161020661020d92610376565b6001610196565b60016101d5565b565b90565b61022d6102286102329261009d565b610216565b61009d565b90565b61023e90610219565b90565b61024a90610235565b90565b90565b61026461025f6102699261024d565b610216565b61009d565b90565b61027590610250565b90565b5f0190565b5f1b90565b9061029360018060a01b039161027d565b9181191691161790565b6102a690610219565b90565b6102b29061029d565b90565b90565b906102cd6102c86102d4926102a9565b6102b5565b8254610282565b9055565b60e01b90565b6102e7906100a8565b90565b6102f3816102de565b036102fa57565b5f80fd5b9050519061030b826102ea565b565b9060208282031261032657610323915f016102fe565b90565b610099565b610333610035565b3d5f823e3d90fd5b61034490610219565b90565b6103509061033b565b90565b90565b9061036b61036661037292610347565b610353565b8254610282565b9055565b61037f81610241565b61039961039361038e5f61026c565b6100a8565b916100a8565b1461042e5760206103b6826103b16103cc945f6102b8565b610241565b63f624f88a906103c4610035565b9384926102d8565b825281806103dc60048201610278565b03915afa8015610429576103f9915f916103fb575b506001610356565b565b61041c915060203d8111610422575b610414818361005d565b81019061030d565b5f6103f1565b503d61040a565b61032b565b5f63d92e233d60e01b81528061044660048201610278565b0390fdfe60806040526004361015610015575b36610b2e57005b61001f5f356100de565b80630892f0ee146100d95780633ccfd60b146100d457806352f40e5e146100cf57806371b3029a146100ca57806388791188146100c5578063a40aad2a146100c0578063ae7d76df146100bb578063bb11fe2b146100b6578063ccce1431146100b1578063d2afb99c146100ac578063f624f88a146100a75763fafb0c6a0361000e57610af9565b610a57565b6109b4565b610890565b6106d9565b61068b565b61063f565b61049c565b610459565b6103d8565b6102bc565b61017c565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f80fd5b909182601f8301121561013c5781359167ffffffffffffffff831161013757602001926001830284011161013257565b6100fe565b6100fa565b6100f6565b90602082820312610172575f82013567ffffffffffffffff811161016d576101699201610102565b9091565b6100f2565b6100ee565b5f0190565b346101ab5761019561018f366004610141565b90610d1f565b61019d6100e4565b806101a781610177565b0390f35b6100ea565b5f9103126101ba57565b6100ee565b5190565b60209181520190565b60200190565b60018060a01b031690565b6101e6906101d2565b90565b6101f2906101dd565b9052565b90565b610202906101f6565b9052565b906020806102289361021e5f8201515f8601906101e9565b01519101906101f9565b565b9061023781604093610206565b0190565b60200190565b9061025e610258610251846101bf565b80936101c3565b926101cc565b905f5b81811061026e5750505090565b909192610287610281600192865161022a565b9461023b565b9101919091610261565b90916102ab6102b99360408401908482035f860152610241565b916020818403910152610241565b90565b346102ed576102cc3660046101b0565b6102d4611589565b906102e96102e06100e4565b92839283610291565b0390f35b6100ea565b6102fb816101f6565b0361030257565b5f80fd5b90503590610313826102f2565b565b9060208282031261032e5761032b915f01610306565b90565b6100ee565b90565b61034a61034561034f926101f6565b610333565b6101f6565b90565b9061035c90610336565b5f5260205260405f2090565b1c90565b60ff1690565b6103829060086103879302610368565b61036c565b90565b906103959154610372565b90565b6103ae906103a96004915f92610352565b61038a565b90565b151590565b6103bf906103b1565b9052565b91906103d6905f602085019401906103b6565b565b34610408576104046103f36103ee366004610315565b610398565b6103fb6100e4565b918291826103c3565b0390f35b6100ea565b5f80fd5b9081604091031261041f5790565b61040d565b90602082820312610454575f82013567ffffffffffffffff811161044f5761044c9201610411565b90565b6100f2565b6100ee565b346104875761047161046c366004610424565b611636565b6104796100e4565b8061048381610177565b0390f35b6100ea565b610499600160159061038a565b90565b346104cc576104ac3660046101b0565b6104c86104b761048c565b6104bf6100e4565b918291826103c3565b0390f35b6100ea565b6104da816103b1565b036104e157565b5f80fd5b905035906104f2826104d1565b565b909182601f8301121561052e5781359167ffffffffffffffff831161052957602001926040830284011161052457565b6100fe565b6100fa565b6100f6565b909160a0828403126105cf575f82013567ffffffffffffffff81116105ca578361055e918401610411565b9261056c81602085016104e5565b9261057a82604083016104e5565b92606082013567ffffffffffffffff81116105c5578361059b9184016104f4565b929093608082013567ffffffffffffffff81116105c0576105bc92016104f4565b9091565b6100f2565b6100f2565b6100f2565b6100ee565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b61061561061e6020936106239361060c816105d4565b938480936105d8565b958691016105e1565b6105ec565b0190565b61063c9160208201915f8184039101526105f6565b90565b3461067657610672610661610655366004610533565b959490949391936116b9565b6106696100e4565b91829182610627565b0390f35b6100ea565b610688600160149061038a565b90565b346106bb5761069b3660046101b0565b6106b76106a661067b565b6106ae6100e4565b918291826103c3565b0390f35b6100ea565b6106d6906106d16003915f92610352565b61038a565b90565b34610709576107056106f46106ef366004610315565b6106c0565b6106fc6100e4565b918291826103c3565b0390f35b6100ea565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610755575b602083101461075057565b610721565b91607f1691610745565b60209181520190565b5f5260205f2090565b905f929180549061078b61078483610735565b809461075f565b916001811690815f146107e257506001146107a6575b505050565b6107b39192939450610768565b915f925b8184106107ca57505001905f80806107a1565b600181602092959395548486015201910192906107b7565b92949550505060ff19168252151560200201905f80806107a1565b9061080791610771565b90565b634e487b7160e01b5f52604160045260245ffd5b90610828906105ec565b810190811067ffffffffffffffff82111761084257604052565b61080a565b90610867610860926108576100e4565b938480926107fd565b038361081e565b565b905f1061087c5761087990610847565b90565b61070e565b61088d60025f90610869565b90565b346108c0576108a03660046101b0565b6108bc6108ab610881565b6108b36100e4565b91829182610627565b0390f35b6100ea565b6108ce816101dd565b036108d557565b5f80fd5b905035906108e6826108c5565b565b90602082820312610901576108fe915f016108d9565b90565b6100ee565b61091a61091561091f926101d2565b610333565b6101d2565b90565b61092b90610906565b90565b61093790610922565b90565b906109449061092e565b5f5260205260405f2090565b90565b6109639060086109689302610368565b610950565b90565b906109769154610953565b90565b61098f9061098a6005915f9261093a565b61096b565b90565b61099b906101f6565b9052565b91906109b2905f60208501940190610992565b565b346109e4576109e06109cf6109ca3660046108e8565b610979565b6109d76100e4565b9182918261099f565b0390f35b6100ea565b60018060a01b031690565b610a04906008610a099302610368565b6109e9565b90565b90610a1791546109f4565b90565b610a2660015f90610a0c565b90565b610a3290610922565b90565b610a3e90610a29565b9052565b9190610a55905f60208501940190610a35565b565b34610a8757610a673660046101b0565b610a83610a72610a1a565b610a7a6100e4565b91829182610a42565b0390f35b6100ea565b60018060a01b031690565b610aa7906008610aac9302610368565b610a8c565b90565b90610aba9154610a97565b90565b610ac85f5f90610aaf565b90565b610ad490610922565b90565b610ae090610acb565b9052565b9190610af7905f60208501940190610ad7565b565b34610b2957610b093660046101b0565b610b25610b14610abd565b610b1c6100e4565b91829182610ae4565b0390f35b6100ea565b5f80fd5b5090565b601f602091010490565b1b90565b91906008610b5f910291610b595f1984610b40565b92610b40565b9181191691161790565b90565b9190610b82610b7d610b8a93610336565b610b69565b908354610b44565b9055565b5f90565b610ba491610b9e610b8e565b91610b6c565b565b5b818110610bb2575050565b80610bbf5f600193610b92565b01610ba7565b9190601f8111610bd5575b505050565b610be1610c0693610768565b906020610bed84610b36565b83019310610c0e575b610bff90610b36565b0190610ba6565b5f8080610bd0565b9150610bff81929050610bf6565b90610c2c905f1990600802610368565b191690565b81610c3b91610c1c565b906002021790565b91610c4e9082610b32565b9067ffffffffffffffff8211610d0d57610c7282610c6c8554610735565b85610bc5565b5f90601f8311600114610ca557918091610c94935f92610c99575b5050610c31565b90555b565b90915001355f80610c8d565b601f19831691610cb485610768565b925f5b818110610cf557509160029391856001969410610cdb575b50505002019055610c97565b610ceb910135601f841690610c1c565b90555f8080610ccf565b91936020600181928787013581550195019201610cb7565b61080a565b90610d1d9291610c43565b565b90610d2b916002610d12565b565b606090565b5f1c90565b610d43610d4891610d32565b610a8c565b90565b610d559054610d37565b90565b9033610d7c610d76610d71610d6c5f610d4b565b610acb565b6101dd565b916101dd565b03610d8e57610d8a916111d0565b9091565b5f63c14b0f1760e01b815280610da660048201610177565b0390fd5b610db6610dbb91610d32565b6109e9565b90565b610dc89054610daa565b90565b60e01b90565b90610de4610ddd6100e4565b928361081e565b565b67ffffffffffffffff8111610dfe5760208091020190565b61080a565b90505190610e10826108c5565b565b90929192610e27610e2282610de6565b610dd1565b9381855260208086019202830192818411610e6457915b838310610e4b5750505050565b60208091610e598486610e03565b815201920191610e3e565b6100fe565b9080601f83011215610e8757816020610e8493519101610e12565b90565b6100f6565b90602082820312610ebc575f82015167ffffffffffffffff8111610eb757610eb49201610e69565b90565b6100f2565b6100ee565b610ec96100e4565b3d5f823e3d90fd5b5190565b90565b610eec610ee7610ef192610ed5565b610333565b6101f6565b90565b634e487b7160e01b5f52601160045260245ffd5b610f17610f1d919392936101f6565b926101f6565b8201809211610f2857565b610ef4565b67ffffffffffffffff8111610f455760208091020190565b61080a565b90610f5c610f5783610f2d565b610dd1565b918252565b610f6b6040610dd1565b90565b5f90565b5f90565b610f7e610f61565b9060208083610f8b610f6e565b815201610f96610f72565b81525050565b610fa4610f76565b90565b5f5b828110610fb557505050565b602090610fc0610f9c565b8184015201610fa9565b90610fef610fd783610f4a565b92602080610fe58693610f2d565b9201910390610fa7565b565b634e487b7160e01b5f52603260045260245ffd5b9061100f82610ed1565b811015611020576020809102010190565b610ff1565b61102f90516101dd565b90565b61103b90610906565b90565b61104790611032565b90565b61105390610922565b90565b61105f90610922565b90565b9050519061106f826102f2565b565b9060208282031261108a57611087915f01611062565b90565b6100ee565b611098906101dd565b9052565b91906110af905f6020850194019061108f565b565b90565b6110c86110c36110cd926110b1565b610333565b6101f6565b90565b905051906110dd826104d1565b565b906020828203126110f8576110f5915f016110d0565b90565b6100ee565b91602061111e92949361111760408201965f83019061108f565b0190610992565b565b61112a6040610dd1565b90565b90611137906101dd565b9052565b90611145906101f6565b9052565b90611153826101bf565b811015611164576020809102010190565b610ff1565b600161117591016101f6565b90565b61118190610906565b90565b61118d90611178565b90565b61119990610922565b90565b6111b06111ab6111b5926110b1565b610333565b6101d2565b90565b6111c19061119c565b90565b6111cd5f6111b8565b90565b50506111fe5f6111e86111e36001610dbe565b610a29565b63024ece89906111f66100e4565b938492610dcb565b8252818061120e60048201610177565b03915afa908115611584575f91611562575b509061122b82610ed1565b906112486112438361123d6001610ed8565b90610f08565b610fca565b61126461125f846112596001610ed8565b90610f08565b610fca565b9161126d610b8e565b91611276610b8e565b9161127f610b8e565b935b8061129461128e896101f6565b916101f6565b146114a9576112ac6112a7898390611005565b611025565b6112f460206112c26112bd8461103e565b61104a565b6370a08231906112e96112d430611056565b926112dd6100e4565b95869485938493610dcb565b83526004830161109c565b03915afa9081156114a4575f91611476575b50908161131b6113155f6110b4565b916101f6565b11611331575b505061132c90611169565b611281565b61134261133d8261103e565b61104a565b602063a9059cbb9161135b6113565f610d4b565b610acb565b906113795f879561138461136d6100e4565b97889687958694610dcb565b8452600484016110fd565b03925af1908161144a575b50155f146113fc5760016113ac575b505061132c905b905f611321565b956113ed6113f4926113da61132c9599916113d16113c8611120565b935f850161112d565b6020830161113b565b8983916113e78383611149565b52611149565b5150611169565b94905f61139e565b9461143d6114449261142a61132c959891611421611418611120565b935f850161112d565b6020830161113b565b8683916114378383611149565b52611149565b5150611169565b936113a5565b61146a9060203d811161146f575b611462818361081e565b8101906110df565b61138f565b503d611458565b611497915060203d811161149d575b61148f818361081e565b810190611071565b5f611306565b503d611485565b610ec1565b5094509294506114b830611056565b3190816114cd6114c75f6110b4565b916101f6565b11611520575b6115156001926115026114e46111c4565b916114f96114f0611120565b935f850161112d565b6020830161113b565b85839161150f8383611149565b52611149565b515001825283529190565b5f80808061154561154061153b61153684610d4b565b610acb565b611184565b611190565b8690828215611559575bf16114d357610ec1565b506108fc61154f565b61157e91503d805f833e611576818361081e565b810190610e8c565b5f611220565b610ec1565b6115a2611594610d2d565b61159c610d2d565b90610d58565b9091565b336115c96115c36115be6115b95f610d4b565b610acb565b6101dd565b916101dd565b036115d9576115d7906115f5565b565b5f63c14b0f1760e01b8152806115f160048201610177565b0390fd5b6115fe90611786565b7fcc0dfa2e3d45df1c774de9e4a88cc12a14bab7c797b8cf7a992f8b2a6f2f4bda6116276100e4565b8061163181610177565b0390a1565b61163f906115a6565b565b606090565b969594939291903361167061166a6116656116605f610d4b565b610acb565b6101dd565b916101dd565b036116815761167e9761169d565b90565b5f63c14b0f1760e01b81528061169960048201610177565b0390fd5b906116b69795939196949296509590919293949561199e565b90565b906116d09695949392916116cb611641565b611646565b90565b60a01c90565b6116e56116ea916116d3565b61036c565b90565b6116f790546116d9565b90565b61170390610922565b90565b5f1b90565b906117175f1991611706565b9181191691161790565b9061173661173161173d92610336565b610b69565b825461170b565b9055565b9061174d60ff91611706565b9181191691161790565b611760906103b1565b90565b90565b9061177b61177661178292611757565b611763565b8254611741565b9055565b5061179160016116ed565b61194f576117ba6117a1306116fa565b316117b560056117af6111c4565b9061093a565b611721565b6117e65f6117d06117cb6001610dbe565b610a29565b63024ece89906117de6100e4565b938492610dcb565b825281806117f660048201610177565b03915afa90811561194a575f91611928575b5090611812610b8e565b9061181c83610ed1565b915b8061183161182b856101f6565b916101f6565b1461190c5761189190602061185f61185a611855611850898690611005565b611025565b61103e565b61104a565b6370a0823190611886611871306116fa565b9261187a6100e4565b96879485938493610dcb565b83526004830161109c565b03915afa918215611907576118d4926118cf915f916118d9575b506118ca60056118c46118bf8a8790611005565b611025565b9061093a565b611721565b611169565b61181e565b6118fa915060203d8111611900575b6118f2818361081e565b810190611071565b5f6118ab565b503d6118e8565b610ec1565b50915050611926600161192160034390610352565b611766565b565b61194491503d805f833e61193c818361081e565b810190610e8c565b5f611808565b610ec1565b5f631fb1984d60e11b81528061196760048201610177565b0390fd5b60a81c90565b61197d6119829161196b565b61036c565b90565b61198f9054611971565b90565b61199b90610847565b90565b505050505050506119ad611641565b506119b86001611985565b6119df576119d260016119cd60044390610352565b611766565b6119dc6002611992565b90565b5f63506f084b60e01b8152806119f760048201610177565b0390fdfea2646970667358221220b0435e57a0d07a8c9339d2ad2aa8fa2d3b8daf61e3bfde698258c7c7fc7a158364736f6c634300081e0033",
}

// TestTrigger is an auto generated Go binding around an Ethereum contract.
type TestTrigger struct {
	abi abi.ABI
}

// NewTestTrigger creates a new instance of TestTrigger.
func NewTestTrigger() *TestTrigger {
	parsed, err := TestTriggerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TestTrigger{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TestTrigger) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _processorEndpoint, bool _revertOnExecute, bool _revertOnPostWithdraw) returns()
func (testTrigger *TestTrigger) PackConstructor(_processorEndpoint common.Address, _revertOnExecute bool, _revertOnPostWithdraw bool) []byte {
	enc, err := testTrigger.abi.Pack("", _processorEndpoint, _revertOnExecute, _revertOnPostWithdraw)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCapturedBalances is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2afb99c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function capturedBalances(address ) view returns(uint256)
func (testTrigger *TestTrigger) PackCapturedBalances(arg0 common.Address) []byte {
	enc, err := testTrigger.abi.Pack("capturedBalances", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCapturedBalances is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2afb99c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function capturedBalances(address ) view returns(uint256)
func (testTrigger *TestTrigger) TryPackCapturedBalances(arg0 common.Address) ([]byte, error) {
	return testTrigger.abi.Pack("capturedBalances", arg0)
}

// UnpackCapturedBalances is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd2afb99c.
//
// Solidity: function capturedBalances(address ) view returns(uint256)
func (testTrigger *TestTrigger) UnpackCapturedBalances(data []byte) (*big.Int, error) {
	out, err := testTrigger.abi.Unpack("capturedBalances", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackExecute is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71b3029a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function execute((bytes[],bytes32[]) appEventData) returns()
func (testTrigger *TestTrigger) PackExecute(appEventData StructsEventData) []byte {
	enc, err := testTrigger.abi.Pack("execute", appEventData)
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
func (testTrigger *TestTrigger) TryPackExecute(appEventData StructsEventData) ([]byte, error) {
	return testTrigger.abi.Pack("execute", appEventData)
}

// PackExecutedInBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb11fe2b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function executedInBlock(uint256 ) view returns(bool)
func (testTrigger *TestTrigger) PackExecutedInBlock(arg0 *big.Int) []byte {
	enc, err := testTrigger.abi.Pack("executedInBlock", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExecutedInBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb11fe2b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function executedInBlock(uint256 ) view returns(bool)
func (testTrigger *TestTrigger) TryPackExecutedInBlock(arg0 *big.Int) ([]byte, error) {
	return testTrigger.abi.Pack("executedInBlock", arg0)
}

// UnpackExecutedInBlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbb11fe2b.
//
// Solidity: function executedInBlock(uint256 ) view returns(bool)
func (testTrigger *TestTrigger) UnpackExecutedInBlock(data []byte) (bool, error) {
	out, err := testTrigger.abi.Unpack("executedInBlock", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackExecutedPostWithdrawsInBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52f40e5e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function executedPostWithdrawsInBlock(uint256 ) view returns(bool)
func (testTrigger *TestTrigger) PackExecutedPostWithdrawsInBlock(arg0 *big.Int) []byte {
	enc, err := testTrigger.abi.Pack("executedPostWithdrawsInBlock", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExecutedPostWithdrawsInBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52f40e5e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function executedPostWithdrawsInBlock(uint256 ) view returns(bool)
func (testTrigger *TestTrigger) TryPackExecutedPostWithdrawsInBlock(arg0 *big.Int) ([]byte, error) {
	return testTrigger.abi.Pack("executedPostWithdrawsInBlock", arg0)
}

// UnpackExecutedPostWithdrawsInBlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52f40e5e.
//
// Solidity: function executedPostWithdrawsInBlock(uint256 ) view returns(bool)
func (testTrigger *TestTrigger) UnpackExecutedPostWithdrawsInBlock(data []byte) (bool, error) {
	out, err := testTrigger.abi.Unpack("executedPostWithdrawsInBlock", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackGetTrustProcessPayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa40aad2a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getTrustProcessPayload((bytes[],bytes32[]) appEventData, bool executeSuccess, bool withdrawSuccess, (address,uint256)[] returnedTokens, (address,uint256)[] failedTokens) returns(bytes)
func (testTrigger *TestTrigger) PackGetTrustProcessPayload(appEventData StructsEventData, executeSuccess bool, withdrawSuccess bool, returnedTokens []StructsTokenAndAmount, failedTokens []StructsTokenAndAmount) []byte {
	enc, err := testTrigger.abi.Pack("getTrustProcessPayload", appEventData, executeSuccess, withdrawSuccess, returnedTokens, failedTokens)
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
func (testTrigger *TestTrigger) TryPackGetTrustProcessPayload(appEventData StructsEventData, executeSuccess bool, withdrawSuccess bool, returnedTokens []StructsTokenAndAmount, failedTokens []StructsTokenAndAmount) ([]byte, error) {
	return testTrigger.abi.Pack("getTrustProcessPayload", appEventData, executeSuccess, withdrawSuccess, returnedTokens, failedTokens)
}

// UnpackGetTrustProcessPayload is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa40aad2a.
//
// Solidity: function getTrustProcessPayload((bytes[],bytes32[]) appEventData, bool executeSuccess, bool withdrawSuccess, (address,uint256)[] returnedTokens, (address,uint256)[] failedTokens) returns(bytes)
func (testTrigger *TestTrigger) UnpackGetTrustProcessPayload(data []byte) ([]byte, error) {
	out, err := testTrigger.abi.Unpack("getTrustProcessPayload", data)
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
func (testTrigger *TestTrigger) PackProcessorEndpoint() []byte {
	enc, err := testTrigger.abi.Pack("processorEndpoint")
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
func (testTrigger *TestTrigger) TryPackProcessorEndpoint() ([]byte, error) {
	return testTrigger.abi.Pack("processorEndpoint")
}

// UnpackProcessorEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfafb0c6a.
//
// Solidity: function processorEndpoint() view returns(address)
func (testTrigger *TestTrigger) UnpackProcessorEndpoint(data []byte) (common.Address, error) {
	out, err := testTrigger.abi.Unpack("processorEndpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRevertOnExecute is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xae7d76df.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revertOnExecute() view returns(bool)
func (testTrigger *TestTrigger) PackRevertOnExecute() []byte {
	enc, err := testTrigger.abi.Pack("revertOnExecute")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevertOnExecute is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xae7d76df.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revertOnExecute() view returns(bool)
func (testTrigger *TestTrigger) TryPackRevertOnExecute() ([]byte, error) {
	return testTrigger.abi.Pack("revertOnExecute")
}

// UnpackRevertOnExecute is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xae7d76df.
//
// Solidity: function revertOnExecute() view returns(bool)
func (testTrigger *TestTrigger) UnpackRevertOnExecute(data []byte) (bool, error) {
	out, err := testTrigger.abi.Unpack("revertOnExecute", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackRevertOnPostWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x88791188.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revertOnPostWithdraw() view returns(bool)
func (testTrigger *TestTrigger) PackRevertOnPostWithdraw() []byte {
	enc, err := testTrigger.abi.Pack("revertOnPostWithdraw")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevertOnPostWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x88791188.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revertOnPostWithdraw() view returns(bool)
func (testTrigger *TestTrigger) TryPackRevertOnPostWithdraw() ([]byte, error) {
	return testTrigger.abi.Pack("revertOnPostWithdraw")
}

// UnpackRevertOnPostWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x88791188.
//
// Solidity: function revertOnPostWithdraw() view returns(bool)
func (testTrigger *TestTrigger) UnpackRevertOnPostWithdraw(data []byte) (bool, error) {
	out, err := testTrigger.abi.Unpack("revertOnPostWithdraw", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackSetTrustedPayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0892f0ee.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setTrustedPayload(bytes payload) returns()
func (testTrigger *TestTrigger) PackSetTrustedPayload(payload []byte) []byte {
	enc, err := testTrigger.abi.Pack("setTrustedPayload", payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetTrustedPayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0892f0ee.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setTrustedPayload(bytes payload) returns()
func (testTrigger *TestTrigger) TryPackSetTrustedPayload(payload []byte) ([]byte, error) {
	return testTrigger.abi.Pack("setTrustedPayload", payload)
}

// PackTokenAllowlist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf624f88a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tokenAllowlist() view returns(address)
func (testTrigger *TestTrigger) PackTokenAllowlist() []byte {
	enc, err := testTrigger.abi.Pack("tokenAllowlist")
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
func (testTrigger *TestTrigger) TryPackTokenAllowlist() ([]byte, error) {
	return testTrigger.abi.Pack("tokenAllowlist")
}

// UnpackTokenAllowlist is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf624f88a.
//
// Solidity: function tokenAllowlist() view returns(address)
func (testTrigger *TestTrigger) UnpackTokenAllowlist(data []byte) (common.Address, error) {
	out, err := testTrigger.abi.Unpack("tokenAllowlist", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTrustedPayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xccce1431.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function trustedPayload() view returns(bytes)
func (testTrigger *TestTrigger) PackTrustedPayload() []byte {
	enc, err := testTrigger.abi.Pack("trustedPayload")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTrustedPayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xccce1431.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function trustedPayload() view returns(bytes)
func (testTrigger *TestTrigger) TryPackTrustedPayload() ([]byte, error) {
	return testTrigger.abi.Pack("trustedPayload")
}

// UnpackTrustedPayload is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xccce1431.
//
// Solidity: function trustedPayload() view returns(bytes)
func (testTrigger *TestTrigger) UnpackTrustedPayload(data []byte) ([]byte, error) {
	out, err := testTrigger.abi.Unpack("trustedPayload", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ccfd60b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw() returns((address,uint256)[], (address,uint256)[])
func (testTrigger *TestTrigger) PackWithdraw() []byte {
	enc, err := testTrigger.abi.Pack("withdraw")
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
func (testTrigger *TestTrigger) TryPackWithdraw() ([]byte, error) {
	return testTrigger.abi.Pack("withdraw")
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
func (testTrigger *TestTrigger) UnpackWithdraw(data []byte) (WithdrawOutput, error) {
	out, err := testTrigger.abi.Unpack("withdraw", data)
	outstruct := new(WithdrawOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]StructsTokenAndAmount)).(*[]StructsTokenAndAmount)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]StructsTokenAndAmount)).(*[]StructsTokenAndAmount)
	return *outstruct, nil
}

// TestTriggerTriggerExecuted represents a TriggerExecuted event raised by the TestTrigger contract.
type TestTriggerTriggerExecuted struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const TestTriggerTriggerExecutedEventName = "TriggerExecuted"

// ContractEventName returns the user-defined event name.
func (TestTriggerTriggerExecuted) ContractEventName() string {
	return TestTriggerTriggerExecutedEventName
}

// UnpackTriggerExecutedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TriggerExecuted()
func (testTrigger *TestTrigger) UnpackTriggerExecutedEvent(log *types.Log) (*TestTriggerTriggerExecuted, error) {
	event := "TriggerExecuted"
	if log.Topics[0] != testTrigger.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TestTriggerTriggerExecuted)
	if len(log.Data) > 0 {
		if err := testTrigger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range testTrigger.abi.Events[event].Inputs {
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
func (testTrigger *TestTrigger) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], testTrigger.abi.Errors["ExecuteReverted"].ID.Bytes()[:4]) {
		return testTrigger.UnpackExecuteRevertedError(raw[4:])
	}
	if bytes.Equal(raw[:4], testTrigger.abi.Errors["NotProcessorEndpoint"].ID.Bytes()[:4]) {
		return testTrigger.UnpackNotProcessorEndpointError(raw[4:])
	}
	if bytes.Equal(raw[:4], testTrigger.abi.Errors["PostWithdrawReverted"].ID.Bytes()[:4]) {
		return testTrigger.UnpackPostWithdrawRevertedError(raw[4:])
	}
	if bytes.Equal(raw[:4], testTrigger.abi.Errors["ZeroAddress"].ID.Bytes()[:4]) {
		return testTrigger.UnpackZeroAddressError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TestTriggerExecuteReverted represents a ExecuteReverted error raised by the TestTrigger contract.
type TestTriggerExecuteReverted struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ExecuteReverted()
func TestTriggerExecuteRevertedErrorID() common.Hash {
	return common.HexToHash("0x3f63309a95c3fe8c3a95779c6418493f12fb76d7a428aeb1588563d6406fa239")
}

// UnpackExecuteRevertedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ExecuteReverted()
func (testTrigger *TestTrigger) UnpackExecuteRevertedError(raw []byte) (*TestTriggerExecuteReverted, error) {
	out := new(TestTriggerExecuteReverted)
	if err := testTrigger.abi.UnpackIntoInterface(out, "ExecuteReverted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TestTriggerNotProcessorEndpoint represents a NotProcessorEndpoint error raised by the TestTrigger contract.
type TestTriggerNotProcessorEndpoint struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotProcessorEndpoint()
func TestTriggerNotProcessorEndpointErrorID() common.Hash {
	return common.HexToHash("0xc14b0f1759c7c9c12c65a360cebd0624cef647f19809e165e59b302bfc46ca6e")
}

// UnpackNotProcessorEndpointError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotProcessorEndpoint()
func (testTrigger *TestTrigger) UnpackNotProcessorEndpointError(raw []byte) (*TestTriggerNotProcessorEndpoint, error) {
	out := new(TestTriggerNotProcessorEndpoint)
	if err := testTrigger.abi.UnpackIntoInterface(out, "NotProcessorEndpoint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TestTriggerPostWithdrawReverted represents a PostWithdrawReverted error raised by the TestTrigger contract.
type TestTriggerPostWithdrawReverted struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PostWithdrawReverted()
func TestTriggerPostWithdrawRevertedErrorID() common.Hash {
	return common.HexToHash("0x506f084b18bfa81d390ee35f4e96d2ef79e0de3cf93871f9035f4e415ba7ebd5")
}

// UnpackPostWithdrawRevertedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PostWithdrawReverted()
func (testTrigger *TestTrigger) UnpackPostWithdrawRevertedError(raw []byte) (*TestTriggerPostWithdrawReverted, error) {
	out := new(TestTriggerPostWithdrawReverted)
	if err := testTrigger.abi.UnpackIntoInterface(out, "PostWithdrawReverted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TestTriggerZeroAddress represents a ZeroAddress error raised by the TestTrigger contract.
type TestTriggerZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAddress()
func TestTriggerZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xd92e233df2717d4a40030e20904abd27b68fcbeede117eaaccbbdac9618c8c73")
}

// UnpackZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAddress()
func (testTrigger *TestTrigger) UnpackZeroAddressError(raw []byte) (*TestTriggerZeroAddress, error) {
	out := new(TestTriggerZeroAddress)
	if err := testTrigger.abi.UnpackIntoInterface(out, "ZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}
