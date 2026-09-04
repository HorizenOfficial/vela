// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mockerc20

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

// MockERC20MetaData contains all meta data concerning the MockERC20 contract.
var MockERC20MetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"ERC2612ExpiredSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC2612InvalidSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"currentNonce\",\"type\":\"uint256\"}],\"name\":\"InvalidAccountNonce\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "MockERC20",
	Bin: "0x610160604052346100735761001b610015610201565b91610281565b610023610078565b611c2e6109eb82396080518161107d015260a051816110b4015260c05181611044015260e0518161160b0152610100518161163001526101205181611172015261014051816111b20152611c2e90f35b61007e565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906100aa90610082565b810190811060018060401b038211176100c257604052565b61008c565b906100da6100d3610078565b92836100a0565b565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b60018060401b03811161010857610104602091610082565b0190565b61008c565b90825f9392825e0152565b9092919261012d610128826100ec565b6100c7565b93818552602085019082840111610149576101479261010d565b565b6100e8565b9080601f8301121561016c5781602061016993519101610118565b90565b6100e4565b60ff1690565b61018081610171565b0361018757565b5f80fd5b9050519061019882610177565b565b90916060828403126101fc575f82015160018060401b0381116101f757836101c391840161014e565b9260208301519060018060401b0382116101f2576101e6816101ef93860161014e565b9360400161018b565b90565b6100e0565b6100e0565b6100dc565b61021f61261980380380610214816100c7565b92833981019061019a565b909192565b5f1b90565b9061023560ff91610224565b9181191691161790565b90565b61025661025161025b92610171565b61023f565b610171565b90565b90565b9061027661027161027d92610242565b61025e565b8254610229565b9055565b6102929061029993928190916102fd565b6008610261565b565b906102ad6102a8836100ec565b6100c7565b918252565b5f7f3100000000000000000000000000000000000000000000000000000000000000910152565b6102e3600161029b565b906102f0602083016102b2565b565b6102fa6102d9565b90565b90610311929161030b6102f2565b90610313565b565b9061031f939291610370565b565b90565b90565b60200190565b5190565b60018060a01b031690565b61035061034b61035592610331565b61023f565b610331565b90565b6103619061033c565b90565b61036d90610358565b90565b6103816103d1946103b69394610405565b6103958161038f6005610321565b906106e7565b610120526103ad836103a76006610321565b906106e7565b61014052610324565b6103c86103c28261032d565b91610327565b2060e052610324565b6103e36103dd8261032d565b91610327565b20610100524660a0526103f4610805565b60805261040030610364565b60c052565b9061040f91610411565b565b9061041b91610671565b565b5190565b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610455575b602083101461045057565b610421565b91607f1691610445565b5f5260205f2090565b601f602091010490565b1b90565b9190600861049191029161048b5f1984610472565b92610472565b9181191691161790565b90565b6104b26104ad6104b79261049b565b61023f565b61049b565b90565b90565b91906104d36104ce6104db9361049e565b6104ba565b908354610476565b9055565b5f90565b6104f5916104ef6104df565b916104bd565b565b5b818110610503575050565b806105105f6001936104e3565b016104f8565b9190601f8111610526575b505050565b6105326105579361045f565b90602061053e84610468565b8301931061055f575b61055090610468565b01906104f7565b5f8080610521565b915061055081929050610547565b1c90565b90610581905f199060080261056d565b191690565b8161059091610571565b906002021790565b906105a28161041d565b9060018060401b038211610660576105c4826105be8554610435565b85610516565b602090601f83116001146105f8579180916105e7935f926105ec575b5050610586565b90555b565b90915001515f806105e0565b601f198316916106078561045f565b925f5b8181106106485750916002939185600196941061062e575b505050020190556105ea565b61063e910151601f841690610571565b90555f8080610622565b9193602060018192878701518155019501920161060a565b61008c565b9061066f91610598565b565b90610680610687926003610665565b6004610665565b565b5f90565b90565b6106a461069f6106a99261068d565b61023f565b61049b565b90565b90565b90565b6106c66106c16106cb926106ac565b610224565b6106af565b90565b6106d860ff6106b2565b90565b6106e4906106af565b90565b906106f0610689565b506107026106fd83610324565b61032d565b61071561070f6020610690565b9161049b565b105f14610729575061072690610960565b90565b5f61073761073d9392610876565b01610665565b61074d6107486106ce565b6106db565b90565b5f90565b7f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f90565b61078290516106af565b90565b61078e906106af565b9052565b61079b9061049b565b9052565b6107a890610331565b90565b6107b49061079f565b9052565b90959492610803946107f26107fc926107e86080966107de60a088019c5f890190610785565b6020870190610785565b6040850190610785565b6060830190610792565b01906107ab565b565b61080d610750565b50610816610754565b61086061082360e0610778565b91610851610832610100610778565b4661083c30610364565b91610845610078565b968795602087016107b8565b602082018103825203826100a0565b61087261086c8261032d565b91610327565b2090565b90565b90565b61089061088b61089592610879565b61023f565b61049b565b90565b60209181520190565b6108c06108c96020936108ce936108b78161041d565b93848093610898565b9586910161010d565b610082565b0190565b6108e79160208201915f8184039101526108a1565b90565b6109046108ff6108f98361032d565b92610327565b610778565b9060208110610912575b5090565b610924905f1990602003600802610472565b165f61090e565b5f1c90565b61093c6109419161092b565b61049e565b90565b61095861095361095d9261049b565b610224565b6106af565b90565b610968610689565b5061097281610324565b9061097c8261032d565b61098f610989601f61087c565b9161049b565b116109c457506109bc816109b66109b06109ab6109c1956108ea565b610930565b9161032d565b17610944565b6106db565b90565b6109e6906109d0610078565b91829163305a27a960e01b8352600483016108d2565b0390fdfe60806040526004361015610013575b610801565b61001d5f356100fc565b806306fdde03146100f7578063095ea7b3146100f257806318160ddd146100ed57806323b872dd146100e8578063313ce567146100e35780633644e515146100de57806340c10f19146100d957806370a08231146100d45780637ecebe00146100cf57806384b0196e146100ca57806395d89b41146100c5578063a9059cbb146100c0578063d505accf146100bb5763dd62ed3e0361000e576107cb565b610764565b610678565b610643565b610607565b6104b3565b61047e565b61042c565b6103f2565b610398565b61033a565b6102cb565b610273565b61018a565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f91031261011a57565b61010c565b5190565b60209181520190565b90825f9392825e0152565b601f801991011690565b61016061016960209361016e936101578161011f565b93848093610123565b9586910161012c565b610137565b0190565b6101879160208201915f818403910152610141565b90565b346101ba5761019a366004610110565b6101b66101a561095e565b6101ad610102565b91829182610172565b0390f35b610108565b60018060a01b031690565b6101d3906101bf565b90565b6101df816101ca565b036101e657565b5f80fd5b905035906101f7826101d6565b565b90565b610205816101f9565b0361020c57565b5f80fd5b9050359061021d826101fc565b565b9190604083820312610247578061023b610244925f86016101ea565b93602001610210565b90565b61010c565b151590565b61025a9061024c565b9052565b9190610271905f60208501940190610251565b565b346102a4576102a061028f61028936600461021f565b90610978565b610297610102565b9182918261025e565b0390f35b610108565b6102b2906101f9565b9052565b91906102c9905f602085019401906102a9565b565b346102fb576102db366004610110565b6102f76102e66109c7565b6102ee610102565b918291826102b6565b0390f35b610108565b90916060828403126103355761033261031b845f85016101ea565b9361032981602086016101ea565b93604001610210565b90565b61010c565b3461036b57610367610356610350366004610300565b916109dd565b61035e610102565b9182918261025e565b0390f35b610108565b60ff1690565b61037f90610370565b9052565b9190610396905f60208501940190610376565b565b346103c8576103a8366004610110565b6103c46103b3610a37565b6103bb610102565b91829182610383565b0390f35b610108565b90565b6103d9906103cd565b9052565b91906103f0905f602085019401906103d0565b565b3461042257610402366004610110565b61041e61040d610a51565b610415610102565b918291826103dd565b0390f35b610108565b5f0190565b3461045b5761044561043f36600461021f565b90610a65565b61044d610102565b8061045781610427565b0390f35b610108565b9060208282031261047957610476915f016101ea565b90565b61010c565b346104ae576104aa610499610494366004610460565b610abe565b6104a1610102565b918291826102b6565b0390f35b610108565b346104e3576104df6104ce6104c9366004610460565b610adc565b6104d6610102565b918291826102b6565b0390f35b610108565b60ff60f81b1690565b6104fa906104e8565b9052565b610507906101ca565b9052565b5190565b60209181520190565b60200190565b610527906101f9565b9052565b906105388160209361051e565b0190565b60200190565b9061055f6105596105528461050b565b809361050f565b92610518565b905f5b81811061056f5750505090565b909192610588610582600192865161052b565b9461053c565b9101919091610562565b939591946105e36105d86105f7956105ca6105ed956106049c9a6105bd60e08c01925f8d01906104f1565b8a820360208c0152610141565b9088820360408a0152610141565b9760608701906102a9565b60808501906104fe565b60a08301906103d0565b60c0818403910152610542565b90565b3461063e57610617366004610110565b61063a610622610bc7565b93610631979597939193610102565b97889788610592565b0390f35b610108565b3461067357610653366004610110565b61066f61065e610c51565b610666610102565b91829182610172565b0390f35b610108565b346106a9576106a561069461068e36600461021f565b90610c67565b61069c610102565b9182918261025e565b0390f35b610108565b6106b781610370565b036106be57565b5f80fd5b905035906106cf826106ae565b565b6106da816103cd565b036106e157565b5f80fd5b905035906106f2826106d1565b565b60e08183031261075f5761070a825f83016101ea565b9261071883602084016101ea565b926107268160408501610210565b926107348260608301610210565b9261075c61074584608085016106c2565b936107538160a086016106e5565b9360c0016106e5565b90565b61010c565b34610799576107836107773660046106f4565b95949094939193610d35565b61078b610102565b8061079581610427565b0390f35b610108565b91906040838203126107c657806107ba6107c3925f86016101ea565b936020016101ea565b90565b61010c565b346107fc576107f86107e76107e136600461079e565b90610e3f565b6107ef610102565b918291826102b6565b0390f35b610108565b5f80fd5b606090565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561083e575b602083101461083957565b61080a565b91607f169161082e565b60209181520190565b5f5260205f2090565b905f929180549061087461086d8361081e565b8094610848565b916001811690815f146108cb575060011461088f575b505050565b61089c9192939450610851565b915f925b8184106108b357505001905f808061088a565b600181602092959395548486015201910192906108a0565b92949550505060ff19168252151560200201905f808061088a565b906108f09161085a565b90565b634e487b7160e01b5f52604160045260245ffd5b9061091190610137565b810190811067ffffffffffffffff82111761092b57604052565b6108f3565b9061095061094992610940610102565b938480926108e6565b0383610907565b565b61095b90610930565b90565b610966610805565b506109716003610952565b90565b5f90565b61099591610984610974565b5061098d610e67565b919091610e74565b600190565b5f90565b5f1c90565b90565b6109b26109b79161099e565b6109a3565b90565b6109c490546109a6565b90565b6109cf61099a565b506109da60026109ba565b90565b91610a07926109ea610974565b506109ff6109f6610e67565b82908491610ec4565b919091610f8d565b600190565b5f90565b60ff1690565b610a22610a279161099e565b610a10565b90565b610a349054610a16565b90565b610a3f610a0c565b50610a4a6008610a2a565b90565b5f90565b610a59610a4d565b50610a6261102a565b90565b90610a6f916110e4565b565b90565b610a88610a83610a8d926101bf565b610a71565b6101bf565b90565b610a9990610a74565b90565b610aa590610a90565b90565b90610ab290610a9c565b5f5260205260405f2090565b610ad4610ad991610acd61099a565b505f610aa8565b6109ba565b90565b610aee90610ae861099a565b50611142565b90565b5f90565b5f90565b606090565b610b0790610a90565b90565b90565b5f1b90565b610b26610b21610b2b92610b0a565b610b0d565b6103cd565b90565b610b42610b3d610b4792610b0a565b610a71565b6101f9565b90565b90610b5d610b56610102565b9283610907565b565b67ffffffffffffffff8111610b775760208091020190565b6108f3565b90610b8e610b8983610b5f565b610b4a565b918252565b369037565b90610bbd610ba583610b7c565b92602080610bb38693610b5f565b9201910390610b93565b565b600f60f81b90565b610bcf610af1565b50610bd8610805565b50610be1610805565b50610bea61099a565b50610bf3610af5565b50610bfc610a4d565b50610c05610af9565b50610c0e611164565b90610c176111a4565b904690610c2330610afe565b90610c2d5f610b12565b90610c3f610c3a5f610b2e565b610b98565b90610c48610bbf565b96959493929190565b610c59610805565b50610c646004610952565b90565b610c8491610c73610974565b50610c7c610e67565b919091610f8d565b600190565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c990565b9194610cf5610cff92989795610ceb60a096610ce1610d069a610cd760c08a019e5f8b01906103d0565b60208901906104fe565b60408701906104fe565b60608501906102a9565b60808301906102a9565b01906102a9565b565b60200190565b5190565b916020610d33929493610d2c60408201965f8301906104fe565b01906104fe565b565b969591939294909442610d50610d4a836101f9565b916101f9565b11610e0a5790610db9610dc2949392610da1610d6a610c89565b610d928c80948c91610d7c8d91611248565b9192610d86610102565b97889660208801610cad565b60208201810382520382610907565b610db3610dad82610d0e565b91610d08565b2061127b565b92909192611298565b80610dd5610dcf876101ca565b916101ca565b03610dea5750610de89293919091610e74565b565b8490610e065f9283926325c0072360e11b845260048401610d12565b0390fd5b610e25905f91829163313c898160e11b8352600483016102b6565b0390fd5b90610e3390610a9c565b5f5260205260405f2090565b610e6491610e5a610e5f92610e5261099a565b506001610e29565b610aa8565b6109ba565b90565b610e6f610af5565b503390565b91610e8292916001926112bf565b565b604090610ead610eb49496959396610ea360608401985f8501906104fe565b60208301906102a9565b01906102a9565b565b90610ec191036101f9565b90565b929192610ed2818390610e3f565b9081610ee7610ee15f196101f9565b916101f9565b10610ef4575b5050509050565b81610f07610f01876101f9565b916101f9565b10610f2d57610f249394610f1c919392610eb6565b905f926112bf565b805f8080610eed565b50610f4c849291925f938493637dc7a0d960e11b855260048501610e84565b0390fd5b610f64610f5f610f6992610b0a565b610a71565b6101bf565b90565b610f7590610f50565b90565b9190610f8b905f602085019401906104fe565b565b9182610fa9610fa3610f9e5f610f6c565b6101ca565b916101ca565b146110035781610fc9610fc3610fbe5f610f6c565b6101ca565b916101ca565b14610fdc57610fda92919091611415565b565b610fff610fe85f610f6c565b5f91829163ec442f0560e01b835260048301610f78565b0390fd5b61102661100f5f610f6c565b5f918291634b637e8f60e11b835260048301610f78565b0390fd5b611032610a4d565b5061103c30610afe565b61106e6110687f00000000000000000000000000000000000000000000000000000000000000006101ca565b916101ca565b14806110aa575b5f1461109f577f000000000000000000000000000000000000000000000000000000000000000090565b6110a76115f5565b90565b50466110de6110d87f00000000000000000000000000000000000000000000000000000000000000006101f9565b916101f9565b14611075565b806110ff6110f96110f45f610f6c565b6101ca565b916101ca565b1461111b57611119916111115f610f6c565b919091611415565b565b61113e6111275f610f6c565b5f91829163ec442f0560e01b835260048301610f78565b0390fd5b61115961115e9161115161099a565b506007610aa8565b6109ba565b90565b90565b61116c610805565b506111a17f000000000000000000000000000000000000000000000000000000000000000061119b6005611161565b9061179b565b90565b6111ac610805565b506111e17f00000000000000000000000000000000000000000000000000000000000000006111db6006611161565b9061179b565b90565b60016111f091016101f9565b90565b906111ff5f1991610b0d565b9181191691161790565b61121d611218611222926101f9565b610a71565b6101f9565b90565b90565b9061123d61123861124492611209565b611225565b82546111f3565b9055565b61125c9061125461099a565b506007610aa8565b611278611268826109ba565b91611272836111e4565b90611228565b90565b61129590611287610a4d565b5061129061102a565b6117e9565b90565b926112b3926112bc946112a9610af5565b50929091926118af565b909291926119da565b90565b9092816112dc6112d66112d15f610f6c565b6101ca565b916101ca565b146113a757836112fc6112f66112f15f610f6c565b6101ca565b916101ca565b14611380576113208361131b61131460018690610e29565b8790610aa8565b611228565b61132a575b505050565b91909161137561136361135d7f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92593610a9c565b93610a9c565b9361136c610102565b918291826102b6565b0390a35f8080611325565b6113a361138c5f610f6c565b5f918291634a1406b160e11b835260048301610f78565b0390fd5b6113ca6113b35f610f6c565b5f91829163e602df0560e01b835260048301610f78565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b6113f16113f7919392936101f9565b926101f9565b820180921161140257565b6113ce565b9061141291016101f9565b90565b9190918061143361142d6114285f610f6c565b6101ca565b916101ca565b145f14611514576114576114508361144b60026109ba565b6113e2565b6002611228565b5b8261147361146d6114685f610f6c565b6101ca565b916101ca565b145f146114e8576114976114908361148b60026109ba565b610eb6565b6002611228565b5b9190916114e36114d16114cb7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef93610a9c565b93610a9c565b936114da610102565b918291826102b6565b0390a3565b61150f826115096114fa5f8790610aa8565b91611504836109ba565b611407565b90611228565b611498565b6115276115225f8390610aa8565b6109ba565b8061153a611534856101f9565b916101f9565b106115625761154d61155d918490610eb6565b6115585f8490610aa8565b611228565b611458565b906115809091925f93849363391434e360e21b855260048501610e84565b0390fd5b7f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f90565b909594926115f3946115e26115ec926115d86080966115ce60a088019c5f8901906103d0565b60208701906103d0565b60408501906103d0565b60608301906102a9565b01906104fe565b565b6115fd610a4d565b50611606611584565b61167d7f00000000000000000000000000000000000000000000000000000000000000009161166e7f00000000000000000000000000000000000000000000000000000000000000004661165930610afe565b91611662610102565b968795602087016115a8565b60208201810382520382610907565b61168f61168982610d0e565b91610d08565b2090565b61169c906103cd565b90565b90565b6116b66116b16116bb9261169f565b610b0d565b6103cd565b90565b6116c860ff6116a2565b90565b5f5260205f2090565b905f92918054906116ee6116e78361081e565b8094610848565b916001811690815f146117455750600114611709575b505050565b61171691929394506116cb565b915f925b81841061172d57505001905f8080611704565b6001816020929593955484860152019101929061171a565b92949550505060ff19168252151560200201905f8080611704565b9061176a916116d4565b90565b9061178d6117869261177d610102565b93848092611760565b0383610907565b565b6117989061176d565b90565b906117a4610805565b506117ae82611693565b6117c76117c16117bc6116be565b6103cd565b916103cd565b14155f146117dc57506117d990611b30565b90565b6117e6915061178f565b90565b6042916117f4610a4d565b50604051917f19010000000000000000000000000000000000000000000000000000000000008352600283015260228201522090565b5f90565b61183a61183f9161099e565b611209565b90565b90565b61185961185461185e92611842565b610a71565b6101f9565b90565b61189661189d9461188c606094989795611882608086019a5f8701906103d0565b6020850190610376565b60408301906103d0565b01906103d0565b565b6118a7610102565b3d5f823e3d90fd5b9392936118ba610af5565b506118c361182a565b506118cc610a4d565b506118d68561182e565b6119086119027f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0611845565b916101f9565b11611995579061192b602094955f94939293611922610102565b94859485611861565b838052039060015afa15611990576119435f51610b0d565b8061195e6119586119535f610f6c565b6101ca565b916101ca565b14611974575f9161196e5f610b12565b91929190565b5061197e5f610f6c565b60019161198a5f610b12565b91929190565b61189f565b5050506119a15f610f6c565b9060039291929190565b634e487b7160e01b5f52602160045260245ffd5b600411156119c957565b6119ab565b906119d8826119bf565b565b806119ed6119e75f6119ce565b916119ce565b145f146119f8575050565b80611a0c611a0660016119ce565b916119ce565b145f14611a2f575f63f645eedf60e01b815280611a2b60048201610427565b0390fd5b80611a43611a3d60026119ce565b916119ce565b145f14611a7157611a6d611a568361182e565b5f91829163fce698f760e01b8352600483016102b6565b0390fd5b611a84611a7e60036119ce565b916119ce565b14611a8c5750565b611aa7905f9182916335e2f38360e21b8352600483016103dd565b0390fd5b90565b611ac2611abd611ac792611aab565b610a71565b6101f9565b90565b67ffffffffffffffff8111611ae857611ae4602091610137565b0190565b6108f3565b90611aff611afa83611aca565b610b4a565b918252565b369037565b90611b2e611b1683611aed565b92602080611b248693611aca565b9201910390611b04565b565b611b38610805565b50611b4281611b9b565b90611b55611b506020611aae565b611b09565b918252602082015290565b611b74611b6f611b799261169f565b610a71565b6101f9565b90565b90565b611b93611b8e611b9892611b7c565b610a71565b6101f9565b90565b611bb0611bb591611baa61099a565b50611693565b61182e565b611bbf60ff611b60565b1680611bd4611bce601f611b7f565b916101f9565b11611bdc5790565b5f632cd44ac360e21b815280611bf460048201610427565b0390fdfea26469706673582212203a4f1ae52b59de23a0b469cf3b46aa2a5c4ca315ef8e6844bc981923ce9955e664736f6c634300081e0033",
}

// MockERC20 is an auto generated Go binding around an Ethereum contract.
type MockERC20 struct {
	abi abi.ABI
}

// NewMockERC20 creates a new instance of MockERC20.
func NewMockERC20() *MockERC20 {
	parsed, err := MockERC20MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MockERC20{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MockERC20) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(string name_, string symbol_, uint8 decimals_) returns()
func (mockERC20 *MockERC20) PackConstructor(name_ string, symbol_ string, decimals_ uint8) []byte {
	enc, err := mockERC20.abi.Pack("", name_, symbol_, decimals_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDOMAINSEPARATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3644e515.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (mockERC20 *MockERC20) PackDOMAINSEPARATOR() []byte {
	enc, err := mockERC20.abi.Pack("DOMAIN_SEPARATOR")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDOMAINSEPARATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3644e515.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (mockERC20 *MockERC20) TryPackDOMAINSEPARATOR() ([]byte, error) {
	return mockERC20.abi.Pack("DOMAIN_SEPARATOR")
}

// UnpackDOMAINSEPARATOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (mockERC20 *MockERC20) UnpackDOMAINSEPARATOR(data []byte) ([32]byte, error) {
	out, err := mockERC20.abi.Unpack("DOMAIN_SEPARATOR", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (mockERC20 *MockERC20) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := mockERC20.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (mockERC20 *MockERC20) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return mockERC20.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (mockERC20 *MockERC20) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := mockERC20.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (mockERC20 *MockERC20) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := mockERC20.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (mockERC20 *MockERC20) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return mockERC20.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (mockERC20 *MockERC20) UnpackApprove(data []byte) (bool, error) {
	out, err := mockERC20.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (mockERC20 *MockERC20) PackBalanceOf(account common.Address) []byte {
	enc, err := mockERC20.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (mockERC20 *MockERC20) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return mockERC20.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (mockERC20 *MockERC20) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := mockERC20.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (mockERC20 *MockERC20) PackDecimals() []byte {
	enc, err := mockERC20.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decimals() view returns(uint8)
func (mockERC20 *MockERC20) TryPackDecimals() ([]byte, error) {
	return mockERC20.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (mockERC20 *MockERC20) UnpackDecimals(data []byte) (uint8, error) {
	out, err := mockERC20.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (mockERC20 *MockERC20) PackEip712Domain() []byte {
	enc, err := mockERC20.abi.Pack("eip712Domain")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (mockERC20 *MockERC20) TryPackEip712Domain() ([]byte, error) {
	return mockERC20.abi.Pack("eip712Domain")
}

// Eip712DomainOutput serves as a container for the return parameters of contract
// method Eip712Domain.
type Eip712DomainOutput struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

// UnpackEip712Domain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (mockERC20 *MockERC20) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := mockERC20.abi.Unpack("eip712Domain", data)
	outstruct := new(Eip712DomainOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (mockERC20 *MockERC20) PackMint(to common.Address, amount *big.Int) []byte {
	enc, err := mockERC20.abi.Pack("mint", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (mockERC20 *MockERC20) TryPackMint(to common.Address, amount *big.Int) ([]byte, error) {
	return mockERC20.abi.Pack("mint", to, amount)
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (mockERC20 *MockERC20) PackName() []byte {
	enc, err := mockERC20.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (mockERC20 *MockERC20) TryPackName() ([]byte, error) {
	return mockERC20.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (mockERC20 *MockERC20) UnpackName(data []byte) (string, error) {
	out, err := mockERC20.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackNonces is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ecebe00.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (mockERC20 *MockERC20) PackNonces(owner common.Address) []byte {
	enc, err := mockERC20.abi.Pack("nonces", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNonces is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ecebe00.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (mockERC20 *MockERC20) TryPackNonces(owner common.Address) ([]byte, error) {
	return mockERC20.abi.Pack("nonces", owner)
}

// UnpackNonces is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (mockERC20 *MockERC20) UnpackNonces(data []byte) (*big.Int, error) {
	out, err := mockERC20.abi.Unpack("nonces", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd505accf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (mockERC20 *MockERC20) PackPermit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) []byte {
	enc, err := mockERC20.abi.Pack("permit", owner, spender, value, deadline, v, r, s)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd505accf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (mockERC20 *MockERC20) TryPackPermit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) ([]byte, error) {
	return mockERC20.abi.Pack("permit", owner, spender, value, deadline, v, r, s)
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (mockERC20 *MockERC20) PackSymbol() []byte {
	enc, err := mockERC20.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (mockERC20 *MockERC20) TryPackSymbol() ([]byte, error) {
	return mockERC20.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (mockERC20 *MockERC20) UnpackSymbol(data []byte) (string, error) {
	out, err := mockERC20.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (mockERC20 *MockERC20) PackTotalSupply() []byte {
	enc, err := mockERC20.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupply() view returns(uint256)
func (mockERC20 *MockERC20) TryPackTotalSupply() ([]byte, error) {
	return mockERC20.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (mockERC20 *MockERC20) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := mockERC20.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (mockERC20 *MockERC20) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := mockERC20.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (mockERC20 *MockERC20) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return mockERC20.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (mockERC20 *MockERC20) UnpackTransfer(data []byte) (bool, error) {
	out, err := mockERC20.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (mockERC20 *MockERC20) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := mockERC20.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (mockERC20 *MockERC20) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return mockERC20.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (mockERC20 *MockERC20) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := mockERC20.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// MockERC20Approval represents a Approval event raised by the MockERC20 contract.
type MockERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const MockERC20ApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (MockERC20Approval) ContractEventName() string {
	return MockERC20ApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (mockERC20 *MockERC20) UnpackApprovalEvent(log *types.Log) (*MockERC20Approval, error) {
	event := "Approval"
	if log.Topics[0] != mockERC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MockERC20Approval)
	if len(log.Data) > 0 {
		if err := mockERC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mockERC20.abi.Events[event].Inputs {
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

// MockERC20EIP712DomainChanged represents a EIP712DomainChanged event raised by the MockERC20 contract.
type MockERC20EIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const MockERC20EIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (MockERC20EIP712DomainChanged) ContractEventName() string {
	return MockERC20EIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (mockERC20 *MockERC20) UnpackEIP712DomainChangedEvent(log *types.Log) (*MockERC20EIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if log.Topics[0] != mockERC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MockERC20EIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := mockERC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mockERC20.abi.Events[event].Inputs {
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

// MockERC20Transfer represents a Transfer event raised by the MockERC20 contract.
type MockERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const MockERC20TransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (MockERC20Transfer) ContractEventName() string {
	return MockERC20TransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (mockERC20 *MockERC20) UnpackTransferEvent(log *types.Log) (*MockERC20Transfer, error) {
	event := "Transfer"
	if log.Topics[0] != mockERC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MockERC20Transfer)
	if len(log.Data) > 0 {
		if err := mockERC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mockERC20.abi.Events[event].Inputs {
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
func (mockERC20 *MockERC20) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return mockERC20.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return mockERC20.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return mockERC20.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC2612ExpiredSignature"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC2612ExpiredSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["ERC2612InvalidSigner"].ID.Bytes()[:4]) {
		return mockERC20.UnpackERC2612InvalidSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["InvalidAccountNonce"].ID.Bytes()[:4]) {
		return mockERC20.UnpackInvalidAccountNonceError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["InvalidShortString"].ID.Bytes()[:4]) {
		return mockERC20.UnpackInvalidShortStringError(raw[4:])
	}
	if bytes.Equal(raw[:4], mockERC20.abi.Errors["StringTooLong"].ID.Bytes()[:4]) {
		return mockERC20.UnpackStringTooLongError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MockERC20ECDSAInvalidSignature represents a ECDSAInvalidSignature error raised by the MockERC20 contract.
type MockERC20ECDSAInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignature()
func MockERC20ECDSAInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0xf645eedf0193584640b6b90cb9477e4c95b98636c148a891d4c0a146dc46e75a")
}

// UnpackECDSAInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignature()
func (mockERC20 *MockERC20) UnpackECDSAInvalidSignatureError(raw []byte) (*MockERC20ECDSAInvalidSignature, error) {
	out := new(MockERC20ECDSAInvalidSignature)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ECDSAInvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ECDSAInvalidSignatureLength represents a ECDSAInvalidSignatureLength error raised by the MockERC20 contract.
type MockERC20ECDSAInvalidSignatureLength struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func MockERC20ECDSAInvalidSignatureLengthErrorID() common.Hash {
	return common.HexToHash("0xfce698f7e8e5342cd615f641317bc45fe7e1e4a8b0a14dd1383ff8dc9c41917f")
}

// UnpackECDSAInvalidSignatureLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func (mockERC20 *MockERC20) UnpackECDSAInvalidSignatureLengthError(raw []byte) (*MockERC20ECDSAInvalidSignatureLength, error) {
	out := new(MockERC20ECDSAInvalidSignatureLength)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ECDSAInvalidSignatureS represents a ECDSAInvalidSignatureS error raised by the MockERC20 contract.
type MockERC20ECDSAInvalidSignatureS struct {
	S [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func MockERC20ECDSAInvalidSignatureSErrorID() common.Hash {
	return common.HexToHash("0xd78bce0cccb935155ed6428d1c13e50b7f3550fd2b66b9fe266006fea4a5e1eb")
}

// UnpackECDSAInvalidSignatureSError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func (mockERC20 *MockERC20) UnpackECDSAInvalidSignatureSError(raw []byte) (*MockERC20ECDSAInvalidSignatureS, error) {
	out := new(MockERC20ECDSAInvalidSignatureS)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the MockERC20 contract.
type MockERC20ERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func MockERC20ERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (mockERC20 *MockERC20) UnpackERC20InsufficientAllowanceError(raw []byte) (*MockERC20ERC20InsufficientAllowance, error) {
	out := new(MockERC20ERC20InsufficientAllowance)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the MockERC20 contract.
type MockERC20ERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func MockERC20ERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (mockERC20 *MockERC20) UnpackERC20InsufficientBalanceError(raw []byte) (*MockERC20ERC20InsufficientBalance, error) {
	out := new(MockERC20ERC20InsufficientBalance)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC20InvalidApprover represents a ERC20InvalidApprover error raised by the MockERC20 contract.
type MockERC20ERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func MockERC20ERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (mockERC20 *MockERC20) UnpackERC20InvalidApproverError(raw []byte) (*MockERC20ERC20InvalidApprover, error) {
	out := new(MockERC20ERC20InvalidApprover)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the MockERC20 contract.
type MockERC20ERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func MockERC20ERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (mockERC20 *MockERC20) UnpackERC20InvalidReceiverError(raw []byte) (*MockERC20ERC20InvalidReceiver, error) {
	out := new(MockERC20ERC20InvalidReceiver)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC20InvalidSender represents a ERC20InvalidSender error raised by the MockERC20 contract.
type MockERC20ERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func MockERC20ERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (mockERC20 *MockERC20) UnpackERC20InvalidSenderError(raw []byte) (*MockERC20ERC20InvalidSender, error) {
	out := new(MockERC20ERC20InvalidSender)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC20InvalidSpender represents a ERC20InvalidSpender error raised by the MockERC20 contract.
type MockERC20ERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func MockERC20ERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (mockERC20 *MockERC20) UnpackERC20InvalidSpenderError(raw []byte) (*MockERC20ERC20InvalidSpender, error) {
	out := new(MockERC20ERC20InvalidSpender)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC2612ExpiredSignature represents a ERC2612ExpiredSignature error raised by the MockERC20 contract.
type MockERC20ERC2612ExpiredSignature struct {
	Deadline *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC2612ExpiredSignature(uint256 deadline)
func MockERC20ERC2612ExpiredSignatureErrorID() common.Hash {
	return common.HexToHash("0x627913023c184eaad13735d5a3d2657ae76ec9a872a70e0fc57522ef1a114d58")
}

// UnpackERC2612ExpiredSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC2612ExpiredSignature(uint256 deadline)
func (mockERC20 *MockERC20) UnpackERC2612ExpiredSignatureError(raw []byte) (*MockERC20ERC2612ExpiredSignature, error) {
	out := new(MockERC20ERC2612ExpiredSignature)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC2612ExpiredSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20ERC2612InvalidSigner represents a ERC2612InvalidSigner error raised by the MockERC20 contract.
type MockERC20ERC2612InvalidSigner struct {
	Signer common.Address
	Owner  common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC2612InvalidSigner(address signer, address owner)
func MockERC20ERC2612InvalidSignerErrorID() common.Hash {
	return common.HexToHash("0x4b800e463b323b1d856edf9dec70329a639d13874a57f5c28219ff57128756db")
}

// UnpackERC2612InvalidSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC2612InvalidSigner(address signer, address owner)
func (mockERC20 *MockERC20) UnpackERC2612InvalidSignerError(raw []byte) (*MockERC20ERC2612InvalidSigner, error) {
	out := new(MockERC20ERC2612InvalidSigner)
	if err := mockERC20.abi.UnpackIntoInterface(out, "ERC2612InvalidSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20InvalidAccountNonce represents a InvalidAccountNonce error raised by the MockERC20 contract.
type MockERC20InvalidAccountNonce struct {
	Account      common.Address
	CurrentNonce *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidAccountNonce(address account, uint256 currentNonce)
func MockERC20InvalidAccountNonceErrorID() common.Hash {
	return common.HexToHash("0x752d88c0de02638abf10e8e31861e4c68dc1f3a1e7d840e580683f2c282bfc7a")
}

// UnpackInvalidAccountNonceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidAccountNonce(address account, uint256 currentNonce)
func (mockERC20 *MockERC20) UnpackInvalidAccountNonceError(raw []byte) (*MockERC20InvalidAccountNonce, error) {
	out := new(MockERC20InvalidAccountNonce)
	if err := mockERC20.abi.UnpackIntoInterface(out, "InvalidAccountNonce", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20InvalidShortString represents a InvalidShortString error raised by the MockERC20 contract.
type MockERC20InvalidShortString struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidShortString()
func MockERC20InvalidShortStringErrorID() common.Hash {
	return common.HexToHash("0xb3512b0c6163e5f0bafab72bb631b9d58cd7a731b082f910338aa21c83d5c274")
}

// UnpackInvalidShortStringError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidShortString()
func (mockERC20 *MockERC20) UnpackInvalidShortStringError(raw []byte) (*MockERC20InvalidShortString, error) {
	out := new(MockERC20InvalidShortString)
	if err := mockERC20.abi.UnpackIntoInterface(out, "InvalidShortString", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MockERC20StringTooLong represents a StringTooLong error raised by the MockERC20 contract.
type MockERC20StringTooLong struct {
	Str string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringTooLong(string str)
func MockERC20StringTooLongErrorID() common.Hash {
	return common.HexToHash("0x305a27a93f8e33b7392df0a0f91d6fc63847395853c45991eec52dbf24d72381")
}

// UnpackStringTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringTooLong(string str)
func (mockERC20 *MockERC20) UnpackStringTooLongError(raw []byte) (*MockERC20StringTooLong, error) {
	out := new(MockERC20StringTooLong)
	if err := mockERC20.abi.UnpackIntoInterface(out, "StringTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}
