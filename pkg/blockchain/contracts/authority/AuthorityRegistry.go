// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package authority

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

// AuthorityRegistryMetaData contains all meta data concerning the AuthorityRegistry contract.
var AuthorityRegistryMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"AppAuthorityContractSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"DefaultAuthorityContractSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"appAuthorityContracts\",\"outputs\":[{\"internalType\":\"contractIAuthorityChecker\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"checkAuthorityIsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultAuthorityContract\",\"outputs\":[{\"internalType\":\"contractIAuthorityChecker\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"defaultAuthority\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"applicationId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"setAppAuthorityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authorityContract\",\"type\":\"address\"}],\"name\":\"setDefaultAuthorityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	ID:  "AuthorityRegistry",
	Bin: "0x60a06040523461003e5761001161004d565b610019610043565b6118516102c8823960805181818161114f015281816111b4015261136f015261185190f35b610049565b60405190565b5f80fd5b6100556100a1565b61005d61019e565b565b60018060a01b031690565b90565b61008161007c6100869261005f565b61006a565b61005f565b90565b6100929061006d565b90565b61009e90610089565b90565b6100aa30610095565b608052565b60401c90565b60ff1690565b6100c76100cc916100af565b6100b5565b90565b6100d990546100bb565b90565b5f0190565b5f1c90565b60018060401b031690565b6100fd610102916100e1565b6100e6565b90565b61010f90546100f1565b90565b60018060401b031690565b5f1b90565b9061013360018060401b039161011d565b9181191691161790565b61015161014c61015692610112565b61006a565b610112565b90565b90565b9061017161016c6101789261013d565b610159565b8254610122565b9055565b61018590610112565b9052565b919061019c905f6020850194019061017c565b565b6101a6610256565b6101b15f82016100cf565b61023a576101c05f8201610105565b6101d86101d260018060401b03610112565b91610112565b036101e1575b50565b6101f4905f60018060401b03910161015c565b60018060401b036102317fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291610228610043565b91829182610189565b0390a15f6101de565b5f63f92ee8a960e01b815280610252600482016100dc565b0390fd5b61025e6102b3565b90565b5f90565b90565b90565b61027f61027a61028492610265565b61011d565b610268565b90565b6102b07ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061026b565b90565b6102bb610261565b506102c4610287565b9056fe60806040526004361015610013575b6107e9565b61001d5f356100dc565b8063116ad4a6146100d75780633aac66d1146100d2578063485cc955146100cd5780634f1ef286146100c857806352d1902d146100c35780635d35bc14146100be5780636f5e32f2146100b9578063715018a6146100b45780638da5cb5b146100af578063ad3cb1cc146100aa578063cfa8904f146100a55763f2fde38b0361000e576107b6565b610781565b61073d565b61060c565b6105b7565b610582565b610459565b610406565b6103a8565b610244565b6101e3565b6101a8565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f80fd5b90565b610100816100f4565b0361010757565b5f80fd5b90503590610118826100f7565b565b60018060a01b031690565b61012e9061011a565b90565b61013a81610125565b0361014157565b5f80fd5b9050359061015282610131565b565b919060408382031261017c5780610170610179925f860161010b565b93602001610145565b90565b6100ec565b151590565b61018f90610181565b9052565b91906101a6905f60208501940190610186565b565b346101d9576101d56101c46101be366004610154565b906108c9565b6101cc6100e2565b91829182610193565b0390f35b6100e8565b5f0190565b34610212576101fc6101f6366004610154565b90610a95565b6102046100e2565b8061020e816101de565b0390f35b6100e8565b919060408382031261023f578061023361023c925f8601610145565b93602001610145565b90565b6100ec565b346102735761025d610257366004610217565b90610e3d565b6102656100e2565b8061026f816101de565b0390f35b6100e8565b5f80fd5b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906102a890610280565b810190811067ffffffffffffffff8211176102c257604052565b61028a565b906102da6102d36100e2565b928361029e565b565b67ffffffffffffffff81116102fa576102f6602091610280565b0190565b61028a565b90825f939282370152565b9092919261031f61031a826102dc565b6102c7565b9381855260208501908284011161033b57610339926102ff565b565b61027c565b9080601f8301121561035e5781602061035b9335910161030a565b90565b610278565b9190916040818403126103a35761037c835f8301610145565b92602082013567ffffffffffffffff811161039e5761039b9201610340565b90565b6100f0565b6100ec565b6103bc6103b6366004610363565b90610e72565b6103c46100e2565b806103ce816101de565b0390f35b5f9103126103dc57565b6100ec565b90565b6103ed906103e1565b9052565b9190610404905f602085019401906103e4565b565b34610436576104163660046103d2565b610432610421610eed565b6104296100e2565b918291826103f1565b0390f35b6100e8565b9060208282031261045457610451915f01610145565b90565b6100ec565b346104875761047161046c36600461043b565b610fa4565b6104796100e2565b80610483816101de565b0390f35b6100e8565b906020828203126104a5576104a2915f0161010b565b90565b6100ec565b90565b6104c16104bc6104c6926100f4565b6104aa565b6100f4565b90565b906104d3906104ad565b5f5260205260405f2090565b1c90565b60018060a01b031690565b6104fe90600861050393026104df565b6104e3565b90565b9061051191546104ee565b90565b610529906105245f915f926104c9565b610506565b90565b61054061053b6105459261011a565b6104aa565b61011a565b90565b6105519061052c565b90565b61055d90610548565b90565b61056990610554565b9052565b9190610580905f60208501940190610560565b565b346105b2576105ae61059d61059836600461048c565b610514565b6105a56100e2565b9182918261056d565b0390f35b6100e8565b346105e5576105c73660046103d2565b6105cf610fd4565b6105d76100e2565b806105e1816101de565b0390f35b6100e8565b6105f390610125565b9052565b919061060a905f602085019401906105ea565b565b3461063c5761061c3660046103d2565b61063861062761100e565b61062f6100e2565b918291826105f7565b0390f35b6100e8565b67ffffffffffffffff811161065f5761065b602091610280565b0190565b61028a565b9061067661067183610641565b6102c7565b918252565b5f7f352e302e30000000000000000000000000000000000000000000000000000000910152565b6106ac6005610664565b906106b96020830161067b565b565b6106c36106a2565b90565b6106ce6106bb565b90565b6106d96106c6565b90565b5190565b60209181520190565b90825f9392825e0152565b61071361071c6020936107219361070a816106dc565b938480936106e0565b958691016106e9565b610280565b0190565b61073a9160208201915f8184039101526106f4565b90565b3461076d5761074d3660046103d2565b6107696107586106d1565b6107606100e2565b91829182610725565b0390f35b6100e8565b61077e60015f90610506565b90565b346107b1576107913660046103d2565b6107ad61079c610772565b6107a46100e2565b9182918261056d565b0390f35b6100e8565b346107e4576107ce6107c936600461043b565b611091565b6107d66100e2565b806107e0816101de565b0390f35b6100e8565b5f80fd5b5f90565b5f1c90565b610802610807916107f1565b6104e3565b90565b61081490546107f6565b90565b90565b61082e61082961083392610817565b6104aa565b61011a565b90565b61083f9061081a565b90565b60e01b90565b61085181610181565b0361085857565b5f80fd5b9050519061086982610848565b565b9060208282031261088457610881915f0161085c565b90565b6100ec565b610892906100f4565b9052565b9160206108b79294936108b060408201965f830190610889565b01906105ea565b565b6108c16100e2565b3d5f823e3d90fd5b906020906108d56107ed565b506108e96108e45f85906104c9565b61080a565b6108f281610554565b61090c6109066109015f610836565b610125565b91610125565b14610991575b61091b90610554565b61093d63116ad4a69492946109486109316100e2565b96879586948594610842565b845260048401610896565b03915afa90811561098c575f9161095e575b5090565b61097f915060203d8111610985575b610977818361029e565b81019061086b565b5f61095a565b503d61096d565b6108b9565b5061091b61099f600161080a565b9050610912565b906109b8916109b361109c565b610a2d565b565b6109c39061052c565b90565b6109cf906109ba565b90565b5f1b90565b906109e860018060a01b03916109d2565b9181191691161790565b6109fb906109ba565b90565b90565b90610a16610a11610a1d926109f2565b6109fe565b82546109d7565b9055565b610a2a90610548565b90565b610a49610a39836109c6565b610a445f84906104c9565b610a01565b90610a7d610a777f4f79164bfe5c672976dae578ab2d7b0c99754246947eb57d8fd2a5a078307c22936104ad565b91610a21565b91610a866100e2565b80610a90816101de565b0390a3565b90610a9f916109a6565b565b60401c90565b60ff1690565b610ab9610abe91610aa1565b610aa7565b90565b610acb9054610aad565b90565b67ffffffffffffffff1690565b610ae7610aec916107f1565b610ace565b90565b610af99054610adb565b90565b67ffffffffffffffff1690565b610b1d610b18610b2292610817565b6104aa565b610afc565b90565b90565b610b3c610b37610b4192610b25565b6104aa565b610afc565b90565b610b4d90610548565b90565b610b64610b5f610b6992610817565b6104aa565b6100f4565b90565b90610b7f67ffffffffffffffff916109d2565b9181191691161790565b610b9d610b98610ba292610afc565b6104aa565b610afc565b90565b90565b90610bbd610bb8610bc492610b89565b610ba5565b8254610b6c565b9055565b60401b90565b90610be268ff000000000000000091610bc8565b9181191691161790565b610bf590610181565b90565b90565b90610c10610c0b610c1792610bec565b610bf8565b8254610bce565b9055565b610c2490610b28565b9052565b9190610c3b905f60208501940190610c1b565b565b90610c466110ea565b91610c5b610c555f8501610ac1565b15610181565b91610c675f8501610aef565b80610c7a610c745f610b09565b91610afc565b1480610d94575b90610c95610c8f6001610b28565b91610afc565b1480610d6c575b610ca7909115610181565b9081610d5b575b50610d3f57610cd791610ccc610cc46001610b28565b5f8701610ba8565b83610d2d575b610d9b565b610cdf575b50565b610cec905f809101610bfb565b6001610d247fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291610d1b6100e2565b91829182610c28565b0390a15f610cdc565b610d3a60015f8701610bfb565b610cd2565b5f63f92ee8a960e01b815280610d57600482016101de565b0390fd5b610d66915015610181565b5f610cae565b50610ca7610d7930610b44565b3b610d8c610d865f610b50565b916100f4565b149050610c9c565b5083610c81565b81610db6610db0610dab5f610836565b610125565b91610125565b14610e2157610dc490611113565b610dcc611128565b610ddf610dd8826109c6565b6001610a01565b610e097fd4047f78dd943d75dfa94ddb7e36315e62df8a01f00a3d0f512960bf93c1231891610a21565b90610e126100e2565b80610e1c816101de565b0390a2565b5f632582a64160e11b815280610e39600482016101de565b0390fd5b90610e4791610c3d565b565b90610e5b91610e5661113e565b610e5d565b565b90610e7091610e6b816111f0565b611260565b565b90610e7c91610e49565b565b5f90565b610e9390610e8e61135e565b610ee1565b90565b90565b610ead610ea8610eb292610e96565b6109d2565b6103e1565b90565b610ede7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610e99565b90565b50610eea610eb5565b90565b610efd610ef8610e7e565b610e82565b90565b610f1190610f0c61109c565b610f13565b565b80610f2e610f28610f235f610836565b610125565b91610125565b14610f8857610f46610f3f826109c6565b6001610a01565b610f707fd4047f78dd943d75dfa94ddb7e36315e62df8a01f00a3d0f512960bf93c1231891610a21565b90610f796100e2565b80610f83816101de565b0390a2565b5f632582a64160e11b815280610fa0600482016101de565b0390fd5b610fad90610f00565b565b610fb761109c565b610fbf610fc1565b565b610fd2610fcd5f610836565b6113df565b565b610fdc610faf565b565b5f90565b60018060a01b031690565b610ff9610ffe916107f1565b610fe2565b90565b61100b9054610fed565b90565b611016610fde565b506110295f61102361144b565b01611001565b90565b61103d9061103861109c565b61103f565b565b8061105a61105461104f5f610836565b610125565b91610125565b1461106a57611068906113df565b565b61108d6110765f610836565b5f918291631e4fbdf760e01b8352600483016105f7565b0390fd5b61109a9061102c565b565b6110a461100e565b6110bd6110b76110b261146f565b610125565b91610125565b036110c457565b6110e66110cf61146f565b5f91829163118cdaa760e01b8352600483016105f7565b0390fd5b6110f26114c7565b90565b611106906111016114db565b611108565b565b61111190611573565b565b61111c906110f5565b565b6111266114db565b565b61113061111e565b565b61113b90610548565b90565b61114730611132565b6111796111737f0000000000000000000000000000000000000000000000000000000000000000610125565b91610125565b1480156111a3575b61118757565b5f63703e46dd60e11b81528061119f600482016101de565b0390fd5b506111ac61157e565b6111de6111d87f0000000000000000000000000000000000000000000000000000000000000000610125565b91610125565b1415611181565b506111ee61109c565b565b6111f9906111e5565b565b6112049061052c565b90565b611210906111fb565b90565b61121c90610548565b90565b611228816103e1565b0361122f57565b5f80fd5b905051906112408261121f565b565b9060208282031261125b57611258915f01611233565b90565b6100ec565b919061128e602061127861127386611207565b611213565b6352d1902d906112866100e2565b938492610842565b8252818061129e600482016101de565b03915afa80915f9261132e575b50155f146112df5750509060016112c057505b565b6112db905f918291634c9c8ce360e01b8352600483016105f7565b0390fd5b92836112fa6112f46112ef610eb5565b6103e1565b916103e1565b0361130f5761130a9293506115a8565b6112be565b61132a845f918291632a87526960e21b8352600483016103f1565b0390fd5b61135091925060203d8111611357575b611348818361029e565b810190611242565b905f6112ab565b503d61133e565b61136730611132565b6113996113937f0000000000000000000000000000000000000000000000000000000000000000610125565b91610125565b036113a057565b5f63703e46dd60e11b8152806113b8600482016101de565b0390fd5b90565b906113d46113cf6113db92610a21565b6113bc565b82546109d7565b9055565b6113e761144b565b6113ff6113f55f8301611001565b915f8491016113bf565b9061143361142d7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093610a21565b91610a21565b9161143c6100e2565b80611446816101de565b0390a3565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930090565b611477610fde565b503390565b90565b61149361148e6114989261147c565b6109d2565b6103e1565b90565b6114c47ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061147f565b90565b6114cf610e7e565b506114d861149b565b90565b6114ec6114e6611631565b15610181565b6114f257565b5f631afcd79f60e31b81528061150a600482016101de565b0390fd5b61151f9061151a6114db565b611521565b565b8061153c6115366115315f610836565b610125565b91610125565b1461154c5761154a906113df565b565b61156f6115585f610836565b5f918291631e4fbdf760e01b8352600483016105f7565b0390fd5b61157c9061150e565b565b611586610fde565b506115a15f61159b611596610eb5565b61164f565b01611001565b90565b5190565b906115b282611652565b816115dd7fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b91610a21565b906115e66100e2565b806115f0816101de565b0390a26115fc816115a4565b61160e6116085f610b50565b916100f4565b115f146116225761161e91611722565b505b565b505061162c6116a7565b611620565b6116396107ed565b5061164c5f6116466110ea565b01610ac1565b90565b90565b803b6116666116605f610b50565b916100f4565b1461168857611686905f61168061167b610eb5565b61164f565b016113bf565b565b6116a3905f918291634c9c8ce360e01b8352600483016105f7565b0390fd5b346116ba6116b45f610b50565b916100f4565b116116c157565b5f63b398979f60e01b8152806116d9600482016101de565b0390fd5b606090565b906116f46116ef836102dc565b6102c7565b918252565b3d5f14611714576117093d6116e2565b903d5f602084013e5b565b61171c6116dd565b90611712565b5f8061174e936117306116dd565b508390602081019051915af4906117456116f9565b90919091611751565b90565b906117659061175e6116dd565b5015610181565b5f1461177157506117d5565b61177a826115a4565b61178c6117865f610b50565b916100f4565b14806117ba575b61179b575090565b6117b6905f918291639996b31560e01b8352600483016105f7565b0390fd5b50803b6117cf6117c95f610b50565b916100f4565b14611793565b6117de816115a4565b6117f06117ea5f610b50565b916100f4565b115f146117ff57805190602001fd5b5f63d6bda27560e01b815280611817600482016101de565b0390fdfea26469706673582212202eb41934dda88a15b9e115a0c8fcdd86b78fb4b7cad55e8ae2b6b80b4e73082764736f6c634300081e0033",
}

// AuthorityRegistry is an auto generated Go binding around an Ethereum contract.
type AuthorityRegistry struct {
	abi abi.ABI
}

// NewAuthorityRegistry creates a new instance of AuthorityRegistry.
func NewAuthorityRegistry() *AuthorityRegistry {
	parsed, err := AuthorityRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &AuthorityRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *AuthorityRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (authorityRegistry *AuthorityRegistry) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := authorityRegistry.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (authorityRegistry *AuthorityRegistry) TryPackUPGRADEINTERFACEVERSION() ([]byte, error) {
	return authorityRegistry.abi.Pack("UPGRADE_INTERFACE_VERSION")
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (authorityRegistry *AuthorityRegistry) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := authorityRegistry.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackAppAuthorityContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f5e32f2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function appAuthorityContracts(uint256 ) view returns(address)
func (authorityRegistry *AuthorityRegistry) PackAppAuthorityContracts(arg0 *big.Int) []byte {
	enc, err := authorityRegistry.abi.Pack("appAuthorityContracts", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAppAuthorityContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f5e32f2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function appAuthorityContracts(uint256 ) view returns(address)
func (authorityRegistry *AuthorityRegistry) TryPackAppAuthorityContracts(arg0 *big.Int) ([]byte, error) {
	return authorityRegistry.abi.Pack("appAuthorityContracts", arg0)
}

// UnpackAppAuthorityContracts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6f5e32f2.
//
// Solidity: function appAuthorityContracts(uint256 ) view returns(address)
func (authorityRegistry *AuthorityRegistry) UnpackAppAuthorityContracts(data []byte) (common.Address, error) {
	out, err := authorityRegistry.abi.Unpack("appAuthorityContracts", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x116ad4a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) PackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckAuthorityIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x116ad4a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) TryPackCheckAuthorityIsAllowed(applicationId *big.Int, authority common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("checkAuthorityIsAllowed", applicationId, authority)
}

// UnpackCheckAuthorityIsAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x116ad4a6.
//
// Solidity: function checkAuthorityIsAllowed(uint256 applicationId, address authority) view returns(bool)
func (authorityRegistry *AuthorityRegistry) UnpackCheckAuthorityIsAllowed(data []byte) (bool, error) {
	out, err := authorityRegistry.abi.Unpack("checkAuthorityIsAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcfa8904f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAuthorityContract() view returns(address)
func (authorityRegistry *AuthorityRegistry) PackDefaultAuthorityContract() []byte {
	enc, err := authorityRegistry.abi.Pack("defaultAuthorityContract")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcfa8904f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAuthorityContract() view returns(address)
func (authorityRegistry *AuthorityRegistry) TryPackDefaultAuthorityContract() ([]byte, error) {
	return authorityRegistry.abi.Pack("defaultAuthorityContract")
}

// UnpackDefaultAuthorityContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcfa8904f.
//
// Solidity: function defaultAuthorityContract() view returns(address)
func (authorityRegistry *AuthorityRegistry) UnpackDefaultAuthorityContract(data []byte) (common.Address, error) {
	out, err := authorityRegistry.abi.Unpack("defaultAuthorityContract", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address owner, address defaultAuthority) returns()
func (authorityRegistry *AuthorityRegistry) PackInitialize(owner common.Address, defaultAuthority common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("initialize", owner, defaultAuthority)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address owner, address defaultAuthority) returns()
func (authorityRegistry *AuthorityRegistry) TryPackInitialize(owner common.Address, defaultAuthority common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("initialize", owner, defaultAuthority)
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (authorityRegistry *AuthorityRegistry) PackOwner() []byte {
	enc, err := authorityRegistry.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (authorityRegistry *AuthorityRegistry) TryPackOwner() ([]byte, error) {
	return authorityRegistry.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (authorityRegistry *AuthorityRegistry) UnpackOwner(data []byte) (common.Address, error) {
	out, err := authorityRegistry.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (authorityRegistry *AuthorityRegistry) PackProxiableUUID() []byte {
	enc, err := authorityRegistry.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (authorityRegistry *AuthorityRegistry) TryPackProxiableUUID() ([]byte, error) {
	return authorityRegistry.abi.Pack("proxiableUUID")
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (authorityRegistry *AuthorityRegistry) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := authorityRegistry.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (authorityRegistry *AuthorityRegistry) PackRenounceOwnership() []byte {
	enc, err := authorityRegistry.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() returns()
func (authorityRegistry *AuthorityRegistry) TryPackRenounceOwnership() ([]byte, error) {
	return authorityRegistry.abi.Pack("renounceOwnership")
}

// PackSetAppAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3aac66d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setAppAuthorityContract(uint256 applicationId, address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) PackSetAppAuthorityContract(applicationId *big.Int, authorityContract common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("setAppAuthorityContract", applicationId, authorityContract)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetAppAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3aac66d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setAppAuthorityContract(uint256 applicationId, address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) TryPackSetAppAuthorityContract(applicationId *big.Int, authorityContract common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("setAppAuthorityContract", applicationId, authorityContract)
}

// PackSetDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d35bc14.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setDefaultAuthorityContract(address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) PackSetDefaultAuthorityContract(authorityContract common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("setDefaultAuthorityContract", authorityContract)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetDefaultAuthorityContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d35bc14.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setDefaultAuthorityContract(address authorityContract) returns()
func (authorityRegistry *AuthorityRegistry) TryPackSetDefaultAuthorityContract(authorityContract common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("setDefaultAuthorityContract", authorityContract)
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (authorityRegistry *AuthorityRegistry) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := authorityRegistry.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (authorityRegistry *AuthorityRegistry) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return authorityRegistry.abi.Pack("transferOwnership", newOwner)
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (authorityRegistry *AuthorityRegistry) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := authorityRegistry.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (authorityRegistry *AuthorityRegistry) TryPackUpgradeToAndCall(newImplementation common.Address, data []byte) ([]byte, error) {
	return authorityRegistry.abi.Pack("upgradeToAndCall", newImplementation, data)
}

// AuthorityRegistryAppAuthorityContractSet represents a AppAuthorityContractSet event raised by the AuthorityRegistry contract.
type AuthorityRegistryAppAuthorityContractSet struct {
	ApplicationId     *big.Int
	AuthorityContract common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryAppAuthorityContractSetEventName = "AppAuthorityContractSet"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryAppAuthorityContractSet) ContractEventName() string {
	return AuthorityRegistryAppAuthorityContractSetEventName
}

// UnpackAppAuthorityContractSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AppAuthorityContractSet(uint256 indexed applicationId, address indexed authorityContract)
func (authorityRegistry *AuthorityRegistry) UnpackAppAuthorityContractSetEvent(log *types.Log) (*AuthorityRegistryAppAuthorityContractSet, error) {
	event := "AppAuthorityContractSet"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryAppAuthorityContractSet)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
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

// AuthorityRegistryDefaultAuthorityContractSet represents a DefaultAuthorityContractSet event raised by the AuthorityRegistry contract.
type AuthorityRegistryDefaultAuthorityContractSet struct {
	AuthorityContract common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryDefaultAuthorityContractSetEventName = "DefaultAuthorityContractSet"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryDefaultAuthorityContractSet) ContractEventName() string {
	return AuthorityRegistryDefaultAuthorityContractSetEventName
}

// UnpackDefaultAuthorityContractSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAuthorityContractSet(address indexed authorityContract)
func (authorityRegistry *AuthorityRegistry) UnpackDefaultAuthorityContractSetEvent(log *types.Log) (*AuthorityRegistryDefaultAuthorityContractSet, error) {
	event := "DefaultAuthorityContractSet"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryDefaultAuthorityContractSet)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
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

// AuthorityRegistryInitialized represents a Initialized event raised by the AuthorityRegistry contract.
type AuthorityRegistryInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryInitialized) ContractEventName() string {
	return AuthorityRegistryInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (authorityRegistry *AuthorityRegistry) UnpackInitializedEvent(log *types.Log) (*AuthorityRegistryInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryInitialized)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
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

// AuthorityRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the AuthorityRegistry contract.
type AuthorityRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryOwnershipTransferred) ContractEventName() string {
	return AuthorityRegistryOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (authorityRegistry *AuthorityRegistry) UnpackOwnershipTransferredEvent(log *types.Log) (*AuthorityRegistryOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
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

// AuthorityRegistryUpgraded represents a Upgraded event raised by the AuthorityRegistry contract.
type AuthorityRegistryUpgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const AuthorityRegistryUpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (AuthorityRegistryUpgraded) ContractEventName() string {
	return AuthorityRegistryUpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (authorityRegistry *AuthorityRegistry) UnpackUpgradedEvent(log *types.Log) (*AuthorityRegistryUpgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != authorityRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(AuthorityRegistryUpgraded)
	if len(log.Data) > 0 {
		if err := authorityRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range authorityRegistry.abi.Events[event].Inputs {
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
func (authorityRegistry *AuthorityRegistry) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["AddressCantBeZero"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], authorityRegistry.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return authorityRegistry.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// AuthorityRegistryAddressCantBeZero represents a AddressCantBeZero error raised by the AuthorityRegistry contract.
type AuthorityRegistryAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressCantBeZero()
func AuthorityRegistryAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x4b054c82bb9e4ebe90744902747c4e86028dd2a122c63de8810d9a1c2e84617d")
}

// UnpackAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressCantBeZero()
func (authorityRegistry *AuthorityRegistry) UnpackAddressCantBeZeroError(raw []byte) (*AuthorityRegistryAddressCantBeZero, error) {
	out := new(AuthorityRegistryAddressCantBeZero)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "AddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryAddressEmptyCode represents a AddressEmptyCode error raised by the AuthorityRegistry contract.
type AuthorityRegistryAddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func AuthorityRegistryAddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (authorityRegistry *AuthorityRegistry) UnpackAddressEmptyCodeError(raw []byte) (*AuthorityRegistryAddressEmptyCode, error) {
	out := new(AuthorityRegistryAddressEmptyCode)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the AuthorityRegistry contract.
type AuthorityRegistryERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func AuthorityRegistryERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (authorityRegistry *AuthorityRegistry) UnpackERC1967InvalidImplementationError(raw []byte) (*AuthorityRegistryERC1967InvalidImplementation, error) {
	out := new(AuthorityRegistryERC1967InvalidImplementation)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryERC1967NonPayable represents a ERC1967NonPayable error raised by the AuthorityRegistry contract.
type AuthorityRegistryERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func AuthorityRegistryERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (authorityRegistry *AuthorityRegistry) UnpackERC1967NonPayableError(raw []byte) (*AuthorityRegistryERC1967NonPayable, error) {
	out := new(AuthorityRegistryERC1967NonPayable)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryFailedCall represents a FailedCall error raised by the AuthorityRegistry contract.
type AuthorityRegistryFailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func AuthorityRegistryFailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (authorityRegistry *AuthorityRegistry) UnpackFailedCallError(raw []byte) (*AuthorityRegistryFailedCall, error) {
	out := new(AuthorityRegistryFailedCall)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryInvalidInitialization represents a InvalidInitialization error raised by the AuthorityRegistry contract.
type AuthorityRegistryInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func AuthorityRegistryInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (authorityRegistry *AuthorityRegistry) UnpackInvalidInitializationError(raw []byte) (*AuthorityRegistryInvalidInitialization, error) {
	out := new(AuthorityRegistryInvalidInitialization)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryNotInitializing represents a NotInitializing error raised by the AuthorityRegistry contract.
type AuthorityRegistryNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func AuthorityRegistryNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (authorityRegistry *AuthorityRegistry) UnpackNotInitializingError(raw []byte) (*AuthorityRegistryNotInitializing, error) {
	out := new(AuthorityRegistryNotInitializing)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the AuthorityRegistry contract.
type AuthorityRegistryOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func AuthorityRegistryOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (authorityRegistry *AuthorityRegistry) UnpackOwnableInvalidOwnerError(raw []byte) (*AuthorityRegistryOwnableInvalidOwner, error) {
	out := new(AuthorityRegistryOwnableInvalidOwner)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the AuthorityRegistry contract.
type AuthorityRegistryOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func AuthorityRegistryOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (authorityRegistry *AuthorityRegistry) UnpackOwnableUnauthorizedAccountError(raw []byte) (*AuthorityRegistryOwnableUnauthorizedAccount, error) {
	out := new(AuthorityRegistryOwnableUnauthorizedAccount)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryUUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the AuthorityRegistry contract.
type AuthorityRegistryUUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func AuthorityRegistryUUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (authorityRegistry *AuthorityRegistry) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*AuthorityRegistryUUPSUnauthorizedCallContext, error) {
	out := new(AuthorityRegistryUUPSUnauthorizedCallContext)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityRegistryUUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the AuthorityRegistry contract.
type AuthorityRegistryUUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func AuthorityRegistryUUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (authorityRegistry *AuthorityRegistry) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*AuthorityRegistryUUPSUnsupportedProxiableUUID, error) {
	out := new(AuthorityRegistryUUPSUnsupportedProxiableUUID)
	if err := authorityRegistry.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
