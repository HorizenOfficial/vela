// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package processorendpointextension

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

// ProcessorEndpointExtensionMetaData contains all meta data concerning the ProcessorEndpointExtension contract.
var ProcessorEndpointExtensionMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DirectCallNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPermit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAmountMismatch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPLOYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_OPERATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"appCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"applicationStateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableDeploySlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"facilitatorNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumOfApplications\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"requestSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"depositPermit\",\"type\":\"bytes\"}],\"name\":\"submitRequestFor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalAppCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalPendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"triggerContracts\",\"outputs\":[{\"internalType\":\"contractITrigger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"triggersToAppIds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "ProcessorEndpointExtension",
	Bin: "0x6101806040523461018d57604080516100188282610191565b60048152602081016356656c6160e01b81525f9061003660016101c8565b9261004385519485610191565b6001845261005160016101c8565b602085019390601f1901368537600a602186015b5f1901916f181899199a1a9b1b9c1cb0b131b232b360811b8282061a835304801561009357600a9091610065565b5050600180556100a2816101e3565b610120526100af8461037e565b61014052519020918260e05251902080610100524660a05282519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f84528483015260608201524660808201523060a082015260a0815261011660c082610191565b5190206080523060c052600a600655600a600755600a600c5530610160525161195290816104b78239608051816117b7015260a05181611874015260c05181611781015260e051816118060152610100518161182c01526101205181610eb801526101405181610ee10152610160518161036c0152f35b5f80fd5b601f909101601f19168101906001600160401b038211908210176101b457604052565b634e487b7160e01b5f52604160045260245ffd5b6001600160401b0381116101b457601f01601f191660200190565b908151602081105f1461025d575090601f81511161021d57602081519101516020821061020e571790565b5f198260200360031b1b161790565b604460209160405192839163305a27a960e01b83528160048401528051918291826024860152018484015e5f828201840152601f01601f19168101030190fd5b6001600160401b0381116101b457600254600181811c91168015610374575b602082101461036057601f811161032d575b50602092601f82116001146102cc57928192935f926102c1575b50508160011b915f199060031b1c19161760025560ff90565b015190505f806102a8565b601f1982169360025f52805f20915f5b86811061031557508360019596106102fd575b505050811b0160025560ff90565b01515f1960f88460031b161c191690555f80806102ef565b919260206001819286850151815501940192016102dc565b60025f52601f60205f20910160051c810190601f830160051c015b818110610355575061028e565b5f8155600101610348565b634e487b7160e01b5f52602260045260245ffd5b90607f169061027c565b908151602081105f146103a9575090601f81511161021d57602081519101516020821061020e571790565b6001600160401b0381116101b457600354600181811c911680156104ac575b602082101461036057601f8111610479575b50602092601f821160011461041857928192935f9261040d575b50508160011b915f199060031b1c19161760035560ff90565b015190505f806103f4565b601f1982169360035f52805f20915f5b8681106104615750836001959610610449575b505050811b0160035560ff90565b01515f1960f88460031b161c191690555f808061043b565b91926020600181928685015181550194019201610428565b60035f52601f60205f20910160051c810190601f830160051c015b8181106104a157506103da565b5f8155600101610494565b90607f16906103c856fe6080806040526004361015610012575f80fd5b5f905f3560e01c90816301ffc9a714611281575080631664deeb146112475780631753520f1461122a57806318733d57146111e9578063194b6be8146111b15780631cc76c3814611171578063248a9ca31461113f5780632a0acc6a146111055780632b52b07c146110cb5780632f2ff15d1461108e57806336568abe146110485780636e91d133146110205780637a36a89114610fe8578063840059ec14610f9857806384b0196e14610ea057806386f4a07d14610e665780638eb8791a14610e4957806391d1485414610e015780639c8d416f14610de4578063a217fddf14610da2578063a9ba045e14610dbc578063aa3aa46014610da2578063af20c96014610d6a578063b7222ff514610d17578063bcf7a3c9146102a0578063c415b95c14610277578063d547741f14610230578063d72e957a146101f7578063ecd00261146101bc578063f624f88a146101935763f7d9cc1e14610173575f80fd5b346101905780600319360112610190576020600c54604051908152f35b80fd5b5034610190578060031936011261019057600f546040516001600160a01b039091168152602090f35b503461019057806003193601126101905760206040517ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c8152f35b5034610190576020366003190112610190576020906040906001600160a01b0361021f6112d4565b168152601683522054604051908152f35b5034610190576040366003190112610190576102736004356102506112ea565b9061026e610269825f525f602052600160405f20015490565b6113ee565b6114ae565b5080f35b50346101905780600319360112610190576015546040516001600160a01b039091168152602090f35b50610140366003190112610a62576102b66112d4565b60243560ff8116809103610a62576044356001600160401b038116809103610a6257606435926005841015610a62576084356001600160401b038111610a625761030490369060040161133a565b60a4356001600160a01b03811696929190879003610a625760c43560e435610104356001600160401b038111610a625761034290369060040161133a565b9096610124356001600160401b038111610a625761036490369060040161133a565b9092909190307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614610d08578b610cf9578a5f52600460205260405f205415610cea57600260015414610cdb57600260015560038914995f8b1580610cce575b610cbf57864211610cb0576103e0611762565b600c541115610ca157610c55578d8c938e9c610c69575b600160a01b60019003169d8e91825f52601660205260405f2054958d8d8d36906104209261139c565b80519060200120916040519460208601967f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc1333340885260408701526060860152608085015260ff1660a084015260c083015260e0820152886101008201528461012082015287610140820152610140815261049b61016082611367565b5190206104a661177e565b906040519161190160f01b835260028301526022820152604290209136906104cd9261139c565b6104d69161152e565b506004811015610c555715801590610c44575b610c35576001600160a01b03168b9003610c265760018101809111610c12578a5f52601660205260405f20556014543410610c035783158080610bfa575b610beb5715610881575b50505061053f36838561139c565b60208151910120600b546040519060208201928a845289604084015287606084015260808301528a60a08301528360c083015260e082015260e0815261058761010082611367565b519020976040519361016085018581106001600160401b0382111761086d5788948c948c9998979460406105e08e968d96835242895260208901938452828901948552606089019534875260808a01978852369161139c565b60a0880190815260c088019687523360e0890190815261010089019a8b5261012089019c8d5261014089019b8c529c89526008602052972095518655516001860180546001600160a01b0319166001600160a01b03929092169190911790555160028501555160038401555160048301559151805190929060058301906001600160401b038111610859576106758254611575565b8d601f8211610810575b90505060208d601f83116001146107a657906007968361079b575b50508160011b915f199060031b1c19161790555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b169160058110156107875793879594936040938a93889660209c5060ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600b5481526009885220556001600b5401600b557fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d54923685604051338152a460018055604051908152f35b634e487b7160e01b89526021600452602489fd5b015190505f8061069a565b9192939495601f198416858452828420935b8181106107f8575091600193918560079998979694106107e0575b505050811b0190556106ae565b01515f1960f88460031b161c191690555f80806107d3565b929360206001819287860151815501950193016107b8565b80846020925220601f830160051c8101916020841061084f575b601f0160051c01908e905b828110610842575061067f565b6001918155018e90610835565b909150819061082a565b634e487b7160e01b8d52604160045260248dfd5b634e487b7160e01b8c52604160045260248cfd5b8a15610beb57600f5460405163cbe230c360e01b8152600481018d905290602090829060249082906001600160a01b03165afa908115610b52575f91610bb0575b5015610ba15760608103610b92578160609181010312610a625780359060ff8216809203610a6257604051636eb1769f60e11b8152600481018b90523060248201526020816044818f5afa8015610b525785915f91610b5d575b5010610ade575b50506040516370a0823160e01b815230600482015290506020816024818c5afa908115610ad3578a91610aa1575b5060208a604051828101906323b872dd60e01b82528b602482015230604482015285606482015260648152610987608482611367565b5190828d5af115610a965789513d610a8d5750883b155b610a79576040516370a0823160e01b8152306004820152906020826024818d5afa8015610a6e5783928c91610a31575b50906109d991611568565b03610a2257858952601260205260408920885f5260205260405f206109ff8282546113e1565b9055878952601360205260408920610a188282546113e1565b90555f8080610531565b63b0a6ea2960e01b8952600489fd5b919250506020813d602011610a66575b81610a4e60209383611367565b81010312610a6257518291906109d96109ce565b5f80fd5b3d9150610a41565b6040513d8d823e3d90fd5b635274afe760e01b8a52600489905260248afd5b6001141561099e565b6040513d8b823e3d90fd5b90506020813d602011610acb575b81610abc60209383611367565b81010312610a6257515f610951565b3d9150610aaf565b6040513d8c823e3d90fd5b8a3b15610a625760409182519363d505accf60e01b85528b600486015230602486015285604486015260648501526084840152602081013560a4840152013560c48201525f8160e481838d5af18015610b5257610b3d575b8080610923565b610b4a9199505f90611367565b5f975f610b36565b6040513d5f823e3d90fd5b9150506020813d602011610b8a575b81610b7960209383611367565b81010312610a62578490515f61091c565b3d9150610b6c565b63ddafbaef60e01b5f5260045ffd5b63514e24c360e11b5f5260045ffd5b90506020813d602011610be3575b81610bcb60209383611367565b81010312610a6257518015158103610a62575f6108c2565b3d9150610bbe565b632a9ffab760e21b5f5260045ffd5b508b1515610527565b63c77f97eb60e01b5f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b632057875960e21b5f5260045ffd5b638baa579f60e01b5f5260045ffd5b506001600160a01b038116156104e9565b634e487b7160e01b5f52602160045260245ffd5b9350509950608587141580610c96575b610c87578b998d8c936103f7565b637c6953f960e01b5f5260045ffd5b5060e2871415610c79565b63d8219b3d60e01b5f5260045ffd5b631ab7da6b60e01b5f5260045ffd5b630ee2db0560e21b5f5260045ffd5b50505f60018b14156103cd565b633ee5aeb560e01b5f5260045ffd5b6309ff9a8360e41b5f5260045ffd5b635428eae760e01b5f5260045ffd5b633921c70360e01b5f5260045ffd5b34610a62576040366003190112610a6257610d30611300565b6001600160401b03610d406112ea565b91165f52601260205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b34610a62576020366003190112610a62576001600160a01b03610d8b6112d4565b165f526013602052602060405f2054604051908152f35b34610a62575f366003190112610a625760206040515f8152f35b34610a62575f366003190112610a6257600e546040516001600160a01b039091168152602090f35b34610a62575f366003190112610a62576020601454604051908152f35b34610a62576040366003190112610a6257610e1a6112ea565b6004355f525f60205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b34610a62575f366003190112610a62576020600654604051908152f35b34610a62575f366003190112610a625760206040517fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c959178152f35b34610a62575f366003190112610a6257610f3c610edc7f00000000000000000000000000000000000000000000000000000000000000006115ad565b610f057f00000000000000000000000000000000000000000000000000000000000000006116aa565b6020610f4a60405192610f188385611367565b5f84525f368137604051958695600f60f81b875260e08588015260e0870190611316565b908582036040870152611316565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b828110610f8157505050500390f35b835185528695509381019392810192600101610f72565b34610a62576040366003190112610a6257610fb16112d4565b610fb96112ea565b6001600160a01b039182165f908152601060209081526040808320949093168252928352819020549051908152f35b34610a62576020366003190112610a62576001600160401b03611009611300565b165f526004602052602060405f2054604051908152f35b34610a62575f366003190112610a6257600d546040516001600160a01b039091168152602090f35b34610a62576040366003190112610a62576110616112ea565b336001600160a01b0382160361107f5761107d906004356114ae565b005b63334bd91960e11b5f5260045ffd5b34610a62576040366003190112610a625761107d6004356110ad6112ea565b906110c6610269825f525f602052600160405f20015490565b611426565b34610a62575f366003190112610a625760206040517f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc13333408152f35b34610a62575f366003190112610a625760206040517fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec428152f35b34610a62576020366003190112610a625760206111696004355f525f602052600160405f20015490565b604051908152f35b34610a62576020366003190112610a62576001600160401b03611192611300565b165f526017602052602060018060a01b0360405f205416604051908152f35b34610a62576020366003190112610a62576001600160a01b036111d26112d4565b165f526011602052602060405f2054604051908152f35b34610a62576020366003190112610a62576001600160a01b0361120a6112d4565b165f52601860205260206001600160401b0360405f205416604051908152f35b34610a62575f366003190112610a62576020600754604051908152f35b34610a62575f366003190112610a625760206040517f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd988152f35b34610a62576020366003190112610a62576004359063ffffffff60e01b8216809203610a6257602091637965db0b60e01b81149081156112c3575b5015158152f35b6301ffc9a760e01b149050836112bc565b600435906001600160a01b0382168203610a6257565b602435906001600160a01b0382168203610a6257565b600435906001600160401b0382168203610a6257565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9181601f84011215610a62578235916001600160401b038311610a625760208381860195010111610a6257565b90601f801991011681019081106001600160401b0382111761138857604052565b634e487b7160e01b5f52604160045260245ffd5b9291926001600160401b03821161138857604051916113c5601f8201601f191660200184611367565b829481845281830111610a62578281602093845f960137010152565b91908201809211610c1257565b5f8181526020818152604080832033845290915290205460ff16156114105750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f818152602081815260408083206001600160a01b038616845290915290205460ff166114a8575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16156114a8575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b815191906041830361155e576115579250602082015190606060408401519301515f1a9061189a565b9192909190565b50505f9160029190565b91908203918211610c1257565b90600182811c921680156115a3575b602083101461158f57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691611584565b60ff81146115f35760ff811690601f82116115e457604051916115d1604084611367565b6020808452838101919036833783525290565b632cd44ac360e21b5f5260045ffd5b50604051600254815f61160583611575565b808352926001811690811561168b575060011461162c575b61162992500382611367565b90565b5060025f90815290917f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace5b81831061166f5750509060206116299282010161161d565b6020919350806001915483858801015201910190918392611657565b6020925061162994915060ff191682840152151560051b82010161161d565b60ff81146116ce5760ff811690601f82116115e457604051916115d1604084611367565b50604051600354815f6116e083611575565b808352926001811690811561168b57506001146117035761162992500382611367565b5060035f90815290917fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b5b8183106117465750509060206116299282010161161d565b602091935080600191548385880101520191019091839261172e565b600b54600a548082116117755750505f90565b61162991611568565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03161480611871575b156117d9577f000000000000000000000000000000000000000000000000000000000000000090565b60405160208101907f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f82527f000000000000000000000000000000000000000000000000000000000000000060408201527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a082015260a0815261186b60c082611367565b51902090565b507f000000000000000000000000000000000000000000000000000000000000000046146117b0565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411611911579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15610b52575f516001600160a01b0381161561190757905f905f90565b505f906001905f90565b5050505f916003919056fea2646970667358221220c89198b377090d1834fbe1c1da6cf74ad6157f87397d02a04890cc2b0703e7f264736f6c634300081e0033",
}

// ProcessorEndpointExtension is an auto generated Go binding around an Ethereum contract.
type ProcessorEndpointExtension struct {
	abi abi.ABI
}

// NewProcessorEndpointExtension creates a new instance of ProcessorEndpointExtension.
func NewProcessorEndpointExtension() *ProcessorEndpointExtension {
	parsed, err := ProcessorEndpointExtensionMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ProcessorEndpointExtension{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ProcessorEndpointExtension) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackADMIN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a0acc6a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackADMIN() []byte {
	enc, err := processorEndpointExtension.abi.Pack("ADMIN")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackADMIN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a0acc6a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackADMIN() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("ADMIN")
}

// UnpackADMIN is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a0acc6a.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackADMIN(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("ADMIN", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackDEFAULTADMINROLE() []byte {
	enc, err := processorEndpointExtension.abi.Pack("DEFAULT_ADMIN_ROLE")
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDEPLOYERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecd00261.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEPLOYER_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackDEPLOYERROLE() []byte {
	enc, err := processorEndpointExtension.abi.Pack("DEPLOYER_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEPLOYERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecd00261.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEPLOYER_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackDEPLOYERROLE() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("DEPLOYER_ROLE")
}

// UnpackDEPLOYERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xecd00261.
//
// Solidity: function DEPLOYER_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackDEPLOYERROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("DEPLOYER_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackPROTOCOLVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa3aa460.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function PROTOCOL_VERSION() view returns(uint8)
func (processorEndpointExtension *ProcessorEndpointExtension) PackPROTOCOLVERSION() []byte {
	enc, err := processorEndpointExtension.abi.Pack("PROTOCOL_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPROTOCOLVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa3aa460.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function PROTOCOL_VERSION() view returns(uint8)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackPROTOCOLVERSION() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("PROTOCOL_VERSION")
}

// UnpackPROTOCOLVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa3aa460.
//
// Solidity: function PROTOCOL_VERSION() view returns(uint8)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackPROTOCOLVERSION(data []byte) (uint8, error) {
	out, err := processorEndpointExtension.abi.Unpack("PROTOCOL_VERSION", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackREQUESTAUTHORIZATIONTYPEHASH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b52b07c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function REQUEST_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackREQUESTAUTHORIZATIONTYPEHASH() []byte {
	enc, err := processorEndpointExtension.abi.Pack("REQUEST_AUTHORIZATION_TYPEHASH")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackREQUESTAUTHORIZATIONTYPEHASH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b52b07c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function REQUEST_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackREQUESTAUTHORIZATIONTYPEHASH() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("REQUEST_AUTHORIZATION_TYPEHASH")
}

// UnpackREQUESTAUTHORIZATIONTYPEHASH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2b52b07c.
//
// Solidity: function REQUEST_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackREQUESTAUTHORIZATIONTYPEHASH(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("REQUEST_AUTHORIZATION_TYPEHASH", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRESETOPERATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86f4a07d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function RESET_OPERATOR() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackRESETOPERATOR() []byte {
	enc, err := processorEndpointExtension.abi.Pack("RESET_OPERATOR")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRESETOPERATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86f4a07d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function RESET_OPERATOR() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackRESETOPERATOR() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("RESET_OPERATOR")
}

// UnpackRESETOPERATOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x86f4a07d.
//
// Solidity: function RESET_OPERATOR() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackRESETOPERATOR(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("RESET_OPERATOR", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackUPDATESTATUSROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1664deeb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function UPDATE_STATUS_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackUPDATESTATUSROLE() []byte {
	enc, err := processorEndpointExtension.abi.Pack("UPDATE_STATUS_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUPDATESTATUSROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1664deeb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function UPDATE_STATUS_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackUPDATESTATUSROLE() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("UPDATE_STATUS_ROLE")
}

// UnpackUPDATESTATUSROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1664deeb.
//
// Solidity: function UPDATE_STATUS_ROLE() view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackUPDATESTATUSROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("UPDATE_STATUS_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAppCustody is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb7222ff5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function appCustody(uint64 , address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackAppCustody(arg0 uint64, arg1 common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("appCustody", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAppCustody is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb7222ff5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function appCustody(uint64 , address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackAppCustody(arg0 uint64, arg1 common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("appCustody", arg0, arg1)
}

// UnpackAppCustody is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb7222ff5.
//
// Solidity: function appCustody(uint64 , address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackAppCustody(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("appCustody", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApplicationStateRoots is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7a36a891.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function applicationStateRoots(uint64 ) view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackApplicationStateRoots(arg0 uint64) []byte {
	enc, err := processorEndpointExtension.abi.Pack("applicationStateRoots", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApplicationStateRoots is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7a36a891.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function applicationStateRoots(uint64 ) view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackApplicationStateRoots(arg0 uint64) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("applicationStateRoots", arg0)
}

// UnpackApplicationStateRoots is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7a36a891.
//
// Solidity: function applicationStateRoots(uint64 ) view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackApplicationStateRoots(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("applicationStateRoots", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAuthorityRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ba045e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function authorityRegistry() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) PackAuthorityRegistry() []byte {
	enc, err := processorEndpointExtension.abi.Pack("authorityRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAuthorityRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ba045e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function authorityRegistry() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackAuthorityRegistry() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("authorityRegistry")
}

// UnpackAuthorityRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9ba045e.
//
// Solidity: function authorityRegistry() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackAuthorityRegistry(data []byte) (common.Address, error) {
	out, err := processorEndpointExtension.abi.Unpack("authorityRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackAvailableDeploySlots is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1753520f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function availableDeploySlots() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackAvailableDeploySlots() []byte {
	enc, err := processorEndpointExtension.abi.Pack("availableDeploySlots")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAvailableDeploySlots is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1753520f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function availableDeploySlots() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackAvailableDeploySlots() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("availableDeploySlots")
}

// UnpackAvailableDeploySlots is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1753520f.
//
// Solidity: function availableDeploySlots() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackAvailableDeploySlots(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("availableDeploySlots", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (processorEndpointExtension *ProcessorEndpointExtension) PackEip712Domain() []byte {
	enc, err := processorEndpointExtension.abi.Pack("eip712Domain")
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackEip712Domain() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("eip712Domain")
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
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := processorEndpointExtension.abi.Unpack("eip712Domain", data)
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

// PackFacilitatorNonces is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd72e957a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function facilitatorNonces(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackFacilitatorNonces(arg0 common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("facilitatorNonces", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFacilitatorNonces is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd72e957a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function facilitatorNonces(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackFacilitatorNonces(arg0 common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("facilitatorNonces", arg0)
}

// UnpackFacilitatorNonces is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd72e957a.
//
// Solidity: function facilitatorNonces(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackFacilitatorNonces(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("facilitatorNonces", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc415b95c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function feeCollector() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) PackFeeCollector() []byte {
	enc, err := processorEndpointExtension.abi.Pack("feeCollector")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc415b95c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function feeCollector() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackFeeCollector() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("feeCollector")
}

// UnpackFeeCollector is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackFeeCollector(data []byte) (common.Address, error) {
	out, err := processorEndpointExtension.abi.Unpack("feeCollector", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := processorEndpointExtension.abi.Pack("getRoleAdmin", role)
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("getRoleAdmin", data)
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
func (processorEndpointExtension *ProcessorEndpointExtension) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("grantRole", role, account)
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (processorEndpointExtension *ProcessorEndpointExtension) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("hasRole", role, account)
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackHasRole(data []byte) (bool, error) {
	out, err := processorEndpointExtension.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMaxNumOfApplications is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8eb8791a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxNumOfApplications() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackMaxNumOfApplications() []byte {
	enc, err := processorEndpointExtension.abi.Pack("maxNumOfApplications")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxNumOfApplications is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8eb8791a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxNumOfApplications() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackMaxNumOfApplications() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("maxNumOfApplications")
}

// UnpackMaxNumOfApplications is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8eb8791a.
//
// Solidity: function maxNumOfApplications() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackMaxNumOfApplications(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("maxNumOfApplications", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxQueueSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7d9cc1e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxQueueSize() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackMaxQueueSize() []byte {
	enc, err := processorEndpointExtension.abi.Pack("maxQueueSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxQueueSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7d9cc1e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxQueueSize() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackMaxQueueSize() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("maxQueueSize")
}

// UnpackMaxQueueSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf7d9cc1e.
//
// Solidity: function maxQueueSize() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackMaxQueueSize(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("maxQueueSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMinFeePerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c8d416f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function minFeePerRequest() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackMinFeePerRequest() []byte {
	enc, err := processorEndpointExtension.abi.Pack("minFeePerRequest")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMinFeePerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c8d416f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function minFeePerRequest() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackMinFeePerRequest() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("minFeePerRequest")
}

// UnpackMinFeePerRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9c8d416f.
//
// Solidity: function minFeePerRequest() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackMinFeePerRequest(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("minFeePerRequest", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPendingClaims is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x840059ec.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingClaims(address , address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackPendingClaims(arg0 common.Address, arg1 common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("pendingClaims", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingClaims is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x840059ec.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingClaims(address , address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackPendingClaims(arg0 common.Address, arg1 common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("pendingClaims", arg0, arg1)
}

// UnpackPendingClaims is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x840059ec.
//
// Solidity: function pendingClaims(address , address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackPendingClaims(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("pendingClaims", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("renounceRole", role, callerConfirmation)
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("revokeRole", role, account)
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("revokeRole", role, account)
}

// PackSubmitRequestFor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbcf7a3c9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitRequestFor(address sender, uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 deadline, bytes requestSignature, bytes depositPermit) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackSubmitRequestFor(sender common.Address, protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, deadline *big.Int, requestSignature []byte, depositPermit []byte) []byte {
	enc, err := processorEndpointExtension.abi.Pack("submitRequestFor", sender, protocolVersion, applicationId, requestType, payload, tokenAddress, assetAmount, deadline, requestSignature, depositPermit)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitRequestFor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbcf7a3c9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitRequestFor(address sender, uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 deadline, bytes requestSignature, bytes depositPermit) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackSubmitRequestFor(sender common.Address, protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, deadline *big.Int, requestSignature []byte, depositPermit []byte) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("submitRequestFor", sender, protocolVersion, applicationId, requestType, payload, tokenAddress, assetAmount, deadline, requestSignature, depositPermit)
}

// UnpackSubmitRequestFor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbcf7a3c9.
//
// Solidity: function submitRequestFor(address sender, uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 deadline, bytes requestSignature, bytes depositPermit) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackSubmitRequestFor(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("submitRequestFor", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (processorEndpointExtension *ProcessorEndpointExtension) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := processorEndpointExtension.abi.Pack("supportsInterface", interfaceId)
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := processorEndpointExtension.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTeeAuthenticator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e91d133.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function teeAuthenticator() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) PackTeeAuthenticator() []byte {
	enc, err := processorEndpointExtension.abi.Pack("teeAuthenticator")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTeeAuthenticator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e91d133.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function teeAuthenticator() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackTeeAuthenticator() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("teeAuthenticator")
}

// UnpackTeeAuthenticator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e91d133.
//
// Solidity: function teeAuthenticator() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTeeAuthenticator(data []byte) (common.Address, error) {
	out, err := processorEndpointExtension.abi.Unpack("teeAuthenticator", data)
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
func (processorEndpointExtension *ProcessorEndpointExtension) PackTokenAllowlist() []byte {
	enc, err := processorEndpointExtension.abi.Pack("tokenAllowlist")
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
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackTokenAllowlist() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("tokenAllowlist")
}

// UnpackTokenAllowlist is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf624f88a.
//
// Solidity: function tokenAllowlist() view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTokenAllowlist(data []byte) (common.Address, error) {
	out, err := processorEndpointExtension.abi.Unpack("tokenAllowlist", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTotalAppCustody is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaf20c960.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalAppCustody(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackTotalAppCustody(arg0 common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("totalAppCustody", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalAppCustody is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaf20c960.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalAppCustody(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackTotalAppCustody(arg0 common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("totalAppCustody", arg0)
}

// UnpackTotalAppCustody is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaf20c960.
//
// Solidity: function totalAppCustody(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTotalAppCustody(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("totalAppCustody", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTotalPendingClaims is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x194b6be8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalPendingClaims(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackTotalPendingClaims(arg0 common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("totalPendingClaims", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalPendingClaims is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x194b6be8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalPendingClaims(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackTotalPendingClaims(arg0 common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("totalPendingClaims", arg0)
}

// UnpackTotalPendingClaims is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x194b6be8.
//
// Solidity: function totalPendingClaims(address ) view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTotalPendingClaims(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("totalPendingClaims", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTriggerContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1cc76c38.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function triggerContracts(uint64 ) view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) PackTriggerContracts(arg0 uint64) []byte {
	enc, err := processorEndpointExtension.abi.Pack("triggerContracts", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTriggerContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1cc76c38.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function triggerContracts(uint64 ) view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackTriggerContracts(arg0 uint64) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("triggerContracts", arg0)
}

// UnpackTriggerContracts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1cc76c38.
//
// Solidity: function triggerContracts(uint64 ) view returns(address)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTriggerContracts(data []byte) (common.Address, error) {
	out, err := processorEndpointExtension.abi.Unpack("triggerContracts", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTriggersToAppIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18733d57.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function triggersToAppIds(address ) view returns(uint64)
func (processorEndpointExtension *ProcessorEndpointExtension) PackTriggersToAppIds(arg0 common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("triggersToAppIds", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTriggersToAppIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18733d57.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function triggersToAppIds(address ) view returns(uint64)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackTriggersToAppIds(arg0 common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("triggersToAppIds", arg0)
}

// UnpackTriggersToAppIds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18733d57.
//
// Solidity: function triggersToAppIds(address ) view returns(uint64)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTriggersToAppIds(data []byte) (uint64, error) {
	out, err := processorEndpointExtension.abi.Unpack("triggersToAppIds", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// ProcessorEndpointExtensionEIP712DomainChanged represents a EIP712DomainChanged event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionEIP712DomainChanged) ContractEventName() string {
	return ProcessorEndpointExtensionEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackEIP712DomainChangedEvent(log *types.Log) (*ProcessorEndpointExtensionEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := processorEndpointExtension.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpointExtension.abi.Events[event].Inputs {
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

// ProcessorEndpointExtensionRequestSubmitted represents a RequestSubmitted event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionRequestSubmitted struct {
	ApplicationId uint64
	RequestId     [32]byte
	Sender        common.Address
	Facilitator   common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionRequestSubmittedEventName = "RequestSubmitted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionRequestSubmitted) ContractEventName() string {
	return ProcessorEndpointExtensionRequestSubmittedEventName
}

// UnpackRequestSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestSubmitted(uint64 indexed applicationId, bytes32 indexed requestId, address indexed sender, address facilitator)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackRequestSubmittedEvent(log *types.Log) (*ProcessorEndpointExtensionRequestSubmitted, error) {
	event := "RequestSubmitted"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionRequestSubmitted)
	if len(log.Data) > 0 {
		if err := processorEndpointExtension.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpointExtension.abi.Events[event].Inputs {
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

// ProcessorEndpointExtensionRoleAdminChanged represents a RoleAdminChanged event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionRoleAdminChanged) ContractEventName() string {
	return ProcessorEndpointExtensionRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackRoleAdminChangedEvent(log *types.Log) (*ProcessorEndpointExtensionRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := processorEndpointExtension.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpointExtension.abi.Events[event].Inputs {
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

// ProcessorEndpointExtensionRoleGranted represents a RoleGranted event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionRoleGranted) ContractEventName() string {
	return ProcessorEndpointExtensionRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackRoleGrantedEvent(log *types.Log) (*ProcessorEndpointExtensionRoleGranted, error) {
	event := "RoleGranted"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionRoleGranted)
	if len(log.Data) > 0 {
		if err := processorEndpointExtension.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpointExtension.abi.Events[event].Inputs {
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

// ProcessorEndpointExtensionRoleRevoked represents a RoleRevoked event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionRoleRevoked) ContractEventName() string {
	return ProcessorEndpointExtensionRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackRoleRevokedEvent(log *types.Log) (*ProcessorEndpointExtensionRoleRevoked, error) {
	event := "RoleRevoked"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionRoleRevoked)
	if len(log.Data) > 0 {
		if err := processorEndpointExtension.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpointExtension.abi.Events[event].Inputs {
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
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["DeadlineExpired"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackDeadlineExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["DirectCallNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackDirectCallNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["FeeValueBelowMinimum"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackFeeValueBelowMinimumError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidApplicationId"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidApplicationIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidPayload"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidPayloadError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidPermit"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidPermitError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidProtocolVersion"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidProtocolVersionError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidRequestType"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidRequestTypeError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidShortString"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidShortStringError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidSigner"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidValue"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidValueError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["QueueThresholdExceeded"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackQueueThresholdExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["StringTooLong"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackStringTooLongError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["TokenNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackTokenNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["TransferAmountMismatch"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackTransferAmountMismatchError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ProcessorEndpointExtensionAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func ProcessorEndpointExtensionAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackAccessControlBadConfirmationError(raw []byte) (*ProcessorEndpointExtensionAccessControlBadConfirmation, error) {
	out := new(ProcessorEndpointExtensionAccessControlBadConfirmation)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func ProcessorEndpointExtensionAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*ProcessorEndpointExtensionAccessControlUnauthorizedAccount, error) {
	out := new(ProcessorEndpointExtensionAccessControlUnauthorizedAccount)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionDeadlineExpired represents a DeadlineExpired error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionDeadlineExpired struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DeadlineExpired()
func ProcessorEndpointExtensionDeadlineExpiredErrorID() common.Hash {
	return common.HexToHash("0x1ab7da6b04cf80fcc6534c1f63ce7b24fd7488ab40444b4e5271ccc1c7f64567")
}

// UnpackDeadlineExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DeadlineExpired()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackDeadlineExpiredError(raw []byte) (*ProcessorEndpointExtensionDeadlineExpired, error) {
	out := new(ProcessorEndpointExtensionDeadlineExpired)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "DeadlineExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionDirectCallNotAllowed represents a DirectCallNotAllowed error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionDirectCallNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DirectCallNotAllowed()
func ProcessorEndpointExtensionDirectCallNotAllowedErrorID() common.Hash {
	return common.HexToHash("0x3921c703755f8ef746475d4b4accdd512ca272fa1597eeba35b810daede33588")
}

// UnpackDirectCallNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DirectCallNotAllowed()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackDirectCallNotAllowedError(raw []byte) (*ProcessorEndpointExtensionDirectCallNotAllowed, error) {
	out := new(ProcessorEndpointExtensionDirectCallNotAllowed)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "DirectCallNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionFeeValueBelowMinimum represents a FeeValueBelowMinimum error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionFeeValueBelowMinimum struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FeeValueBelowMinimum()
func ProcessorEndpointExtensionFeeValueBelowMinimumErrorID() common.Hash {
	return common.HexToHash("0xc77f97ebda9854b3f6367f2d96830b6ea60a9ab7e12c67ede499fadeb3f9cef7")
}

// UnpackFeeValueBelowMinimumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FeeValueBelowMinimum()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackFeeValueBelowMinimumError(raw []byte) (*ProcessorEndpointExtensionFeeValueBelowMinimum, error) {
	out := new(ProcessorEndpointExtensionFeeValueBelowMinimum)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "FeeValueBelowMinimum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidApplicationId represents a InvalidApplicationId error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidApplicationId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidApplicationId()
func ProcessorEndpointExtensionInvalidApplicationIdErrorID() common.Hash {
	return common.HexToHash("0x9ff9a83079daf92a88066c1b25d04540f443a8d245022db6e1e62beaa6d6bc5e")
}

// UnpackInvalidApplicationIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidApplicationId()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidApplicationIdError(raw []byte) (*ProcessorEndpointExtensionInvalidApplicationId, error) {
	out := new(ProcessorEndpointExtensionInvalidApplicationId)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidApplicationId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidPayload represents a InvalidPayload error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidPayload struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPayload()
func ProcessorEndpointExtensionInvalidPayloadErrorID() common.Hash {
	return common.HexToHash("0x7c6953f9f55fcfd713fe1d72e370e4757e17ae7187e1acdcc02b4b76b56a9a35")
}

// UnpackInvalidPayloadError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPayload()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidPayloadError(raw []byte) (*ProcessorEndpointExtensionInvalidPayload, error) {
	out := new(ProcessorEndpointExtensionInvalidPayload)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidPayload", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidPermit represents a InvalidPermit error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidPermit struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPermit()
func ProcessorEndpointExtensionInvalidPermitErrorID() common.Hash {
	return common.HexToHash("0xddafbaefe5cdbfc41219261b5d7468d1f16bc986eedbdef2ea223830eb9c5e3e")
}

// UnpackInvalidPermitError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPermit()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidPermitError(raw []byte) (*ProcessorEndpointExtensionInvalidPermit, error) {
	out := new(ProcessorEndpointExtensionInvalidPermit)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidPermit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidProtocolVersion represents a InvalidProtocolVersion error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidProtocolVersion struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidProtocolVersion()
func ProcessorEndpointExtensionInvalidProtocolVersionErrorID() common.Hash {
	return common.HexToHash("0x5428eae7f7dcfe744a90b56d360bbe512220bcf69b1e537f3d2d3b28cf2bf92d")
}

// UnpackInvalidProtocolVersionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidProtocolVersion()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidProtocolVersionError(raw []byte) (*ProcessorEndpointExtensionInvalidProtocolVersion, error) {
	out := new(ProcessorEndpointExtensionInvalidProtocolVersion)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidProtocolVersion", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidRequestType represents a InvalidRequestType error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidRequestType struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRequestType()
func ProcessorEndpointExtensionInvalidRequestTypeErrorID() common.Hash {
	return common.HexToHash("0x3b8b6c143e72b9be1ff9776689ac73649bd2f21a04954ef955d859dcfff7eb57")
}

// UnpackInvalidRequestTypeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRequestType()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidRequestTypeError(raw []byte) (*ProcessorEndpointExtensionInvalidRequestType, error) {
	out := new(ProcessorEndpointExtensionInvalidRequestType)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidRequestType", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidShortString represents a InvalidShortString error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidShortString struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidShortString()
func ProcessorEndpointExtensionInvalidShortStringErrorID() common.Hash {
	return common.HexToHash("0xb3512b0c6163e5f0bafab72bb631b9d58cd7a731b082f910338aa21c83d5c274")
}

// UnpackInvalidShortStringError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidShortString()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidShortStringError(raw []byte) (*ProcessorEndpointExtensionInvalidShortString, error) {
	out := new(ProcessorEndpointExtensionInvalidShortString)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidShortString", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidSignature represents a InvalidSignature error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature()
func ProcessorEndpointExtensionInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x8baa579fce362245063d36f11747a89dd489c54795634fc673cc0e0db51fedc5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidSignatureError(raw []byte) (*ProcessorEndpointExtensionInvalidSignature, error) {
	out := new(ProcessorEndpointExtensionInvalidSignature)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidSigner represents a InvalidSigner error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidSigner struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSigner()
func ProcessorEndpointExtensionInvalidSignerErrorID() common.Hash {
	return common.HexToHash("0x815e1d64efb74fbe314c20a2b8a2335d18bce12a19165e447fa36bcb35959528")
}

// UnpackInvalidSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSigner()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidSignerError(raw []byte) (*ProcessorEndpointExtensionInvalidSigner, error) {
	out := new(ProcessorEndpointExtensionInvalidSigner)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionInvalidValue represents a InvalidValue error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidValue struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidValue()
func ProcessorEndpointExtensionInvalidValueErrorID() common.Hash {
	return common.HexToHash("0xaa7feadc1a228aee3d4475b95bb4402ebafd3841876f41c4469039f5ee2998d5")
}

// UnpackInvalidValueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidValue()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidValueError(raw []byte) (*ProcessorEndpointExtensionInvalidValue, error) {
	out := new(ProcessorEndpointExtensionInvalidValue)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidValue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionQueueThresholdExceeded represents a QueueThresholdExceeded error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionQueueThresholdExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QueueThresholdExceeded()
func ProcessorEndpointExtensionQueueThresholdExceededErrorID() common.Hash {
	return common.HexToHash("0xd8219b3da97d7adbef217ce9f220bbf38600f597fb984fd0620ea7df76ea17f3")
}

// UnpackQueueThresholdExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QueueThresholdExceeded()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackQueueThresholdExceededError(raw []byte) (*ProcessorEndpointExtensionQueueThresholdExceeded, error) {
	out := new(ProcessorEndpointExtensionQueueThresholdExceeded)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "QueueThresholdExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func ProcessorEndpointExtensionReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackReentrancyGuardReentrantCallError(raw []byte) (*ProcessorEndpointExtensionReentrancyGuardReentrantCall, error) {
	out := new(ProcessorEndpointExtensionReentrancyGuardReentrantCall)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func ProcessorEndpointExtensionSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackSafeERC20FailedOperationError(raw []byte) (*ProcessorEndpointExtensionSafeERC20FailedOperation, error) {
	out := new(ProcessorEndpointExtensionSafeERC20FailedOperation)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionStringTooLong represents a StringTooLong error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionStringTooLong struct {
	Str string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringTooLong(string str)
func ProcessorEndpointExtensionStringTooLongErrorID() common.Hash {
	return common.HexToHash("0x305a27a93f8e33b7392df0a0f91d6fc63847395853c45991eec52dbf24d72381")
}

// UnpackStringTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringTooLong(string str)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackStringTooLongError(raw []byte) (*ProcessorEndpointExtensionStringTooLong, error) {
	out := new(ProcessorEndpointExtensionStringTooLong)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "StringTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionTokenNotAllowed represents a TokenNotAllowed error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionTokenNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenNotAllowed()
func ProcessorEndpointExtensionTokenNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xa29c498656b92b53bef33b5e62c8bb31470e81c408996e42696a26083f02fc8d")
}

// UnpackTokenNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenNotAllowed()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTokenNotAllowedError(raw []byte) (*ProcessorEndpointExtensionTokenNotAllowed, error) {
	out := new(ProcessorEndpointExtensionTokenNotAllowed)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "TokenNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionTransferAmountMismatch represents a TransferAmountMismatch error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionTransferAmountMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TransferAmountMismatch()
func ProcessorEndpointExtensionTransferAmountMismatchErrorID() common.Hash {
	return common.HexToHash("0xb0a6ea29805702cf1e83dfa66f5843be54f820580489524674c9eb97deec9d9f")
}

// UnpackTransferAmountMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TransferAmountMismatch()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTransferAmountMismatchError(raw []byte) (*ProcessorEndpointExtensionTransferAmountMismatch, error) {
	out := new(ProcessorEndpointExtensionTransferAmountMismatch)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "TransferAmountMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}
