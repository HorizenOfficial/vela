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

// ProcessorEndpointExtensionMetaData contains all meta data concerning the ProcessorEndpointExtension contract.
var ProcessorEndpointExtensionMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeployerNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DirectCallNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPermit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxNumOfApplicationsExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerCannotBeEOA\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"DeployRequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMax\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"MaxNumberOfAppUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newGrace\",\"type\":\"uint256\"}],\"name\":\"SelectionGraceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"withdrawSuccess\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"postWithdrawSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"returnedTokens\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"failedTokens\",\"type\":\"tuple[]\"}],\"name\":\"TriggerWithdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPLOYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_OPERATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"addAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminReset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"appIds\",\"type\":\"uint64[]\"}],\"name\":\"adminResetApps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"appCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"applicationStateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableDeploySlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"facilitatorNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resetOperator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"},{\"internalType\":\"contractITokenAllowlist\",\"name\":\"_tokenAllowlist\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"claimable\",\"type\":\"address[]\"}],\"name\":\"invokeTrigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumOfApplications\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"removeAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"selectionGrace\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitDeployRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"trigger\",\"type\":\"address\"}],\"name\":\"submitDeployRequestWithTrigger\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"requestSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"depositPermit\",\"type\":\"bytes\"}],\"name\":\"submitRequestFor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalAppCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalPendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"triggerContracts\",\"outputs\":[{\"internalType\":\"contractITrigger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"triggersToAppIds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"updateMaxNumOfApplications\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newGrace\",\"type\":\"uint256\"}],\"name\":\"updateSelectionGrace\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "ProcessorEndpointExtension",
	Bin: "0x60a0806040523461007e57306080526148899081610083823960805181818161030801528181610a3601528181610afb01528181611164015281816112d201528181611e2f01528181611f0601528181611f9c015281816121b001528181612493015281816125aa0152818161267701528181612842015261289c0152f35b5f80fdfe6080806040526004361015610012575f80fd5b5f905f3560e01c90816301ffc9a714612a08575080631664deeb146129ce5780631753520f146129b157806318733d5714612970578063194b6be8146129385780631cc76c38146128f857806320ff473f1461288057806321c0b3421461281e578063248a9ca3146127df5780632a0acc6a146127b85780632b52b07c1461277e5780632f2ff15d1461273457806336568abe146126f057806344ce30f0146126d35780634a0d04951461263e578063655f2de11461259357806368a7c3fb146124515780636e91d133146124295780637a36a891146123f2578063840059ec146123a257806384b0196e14612258578063850103841461202d57806386f4a07d14611ff35780638c5b9b0014611f895780638eb8791a14611f6c57806390fda89b14611eef57806391d1485414611e9a5780639c73eeaf14611e185780639c8d416f14611dfb578063a217fddf14611db9578063a9ba045e14611dd3578063aa3aa46014611db9578063af20c96014611d81578063b7222ff514611d2e578063bcf7a3c914611213578063c415b95c146111ea578063d2c35ce814611147578063d547741f146110f7578063d72e957a146110be578063dd39d76514610aa5578063e3d72e5e14610a19578063eabcca101461027e578063ecd0026114610256578063f624f88a1461022d5763f7d9cc1e1461020d575f80fd5b3461022a578060031936011261022a576020600e54604051908152f35b80fd5b503461022a578060031936011261022a576011546040516001600160a01b039091168152602090f35b503461022a578060031936011261022a5760206040515f5160206147945f395f51905f528152f35b503461022a5760e036600319011261022a576004356001600160a01b038116919082900361022a576024356001600160a01b03811690819003610a15576102c3612a87565b6064356001600160a01b0381169290838103610a1157608435926001600160a01b03841692838503610a0d5760c4356001600160a01b0381169190829003610a0957307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316146109fa575f5160206148145f395f51905f5254966001600160401b0360ff8960401c16159816801590816109f2575b60011490816109e8575b1590816109df575b506109d0578760016001600160401b03195f5160206148145f395f51905f525416175f5160206148145f395f51905f525561099b575b8915908115610992575b8115610980575b8115610977575b50801561096f575b6109605787986103d6614601565b88898a610955575b806305f5e100600a921015610944575b612710811015610935575b6064811015610927575b101561091d575b60018a01996104188b612b97565b9a6040519b610427908d612b5f565b808c52601f199061043790612b97565b013660208d01378a01602101600a905b5f1901916f181899199a1a9b1b9c1cb0b131b232b360811b8282061a835304801561047557600a9091610447565b505060409889516104868b82612b5f565b600481526356656c6160e01b602082015261049f614601565b6104a7614601565b80516001600160401b038111610909578c6104cf5f5160206147545f395f51905f5254613dc8565b601f81116108cc575b505060208d601f831160011461085057906105079383610762575b50508160011b915f199060031b1c19161790565b5f5160206147545f395f51905f52555b8051906001600160401b03821161083c578b6105405f5160206147745f395f51905f5254613dc8565b601f81116107f4575b50506020908c601f841160011461076d579280610636979593610586936106d79b9a9896926107625750508160011b915f199060031b1c19161790565b5f5160206147745f395f51905f52555b8b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100558b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d101556001600160601b0360a01b600f541617600f556001600160601b0360a01b60105416176010556001600160601b0360a01b601154161760115560018060a01b0381166001600160601b0360a01b6017541617601755613585565b506106408161363e565b505f5160206147945f395f51905f5287525f5160206147d45f395f51905f526020525f5160206147b45f395f51905f52600187892001545f5160206147945f395f51905f5289525f5160206147d45f395f51905f52602052816001898b2001555f5160206147945f395f51905f527fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff8a80a4613500565b5060a435601655610752575b50600a600255600a600355600a600e55603c601b55610700575080f35b60207fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29160ff60401b195f5160206148145f395f51905f5254165f5160206148145f395f51905f52555160018152a180f35b61075b906136e4565b505f6106e3565b015190505f806104f3565b9190601f1984165f5160206147745f395f51905f528452828420935b8181106107dc57509260019285926106d79b9a98966106369a9896106107c4575b505050811b015f5160206147745f395f51905f5255610596565b01515f1960f88460031b161c191690555f80806107aa565b92936020600181928786015181550195019301610789565b6020825f5160206147745f395f51905f5261082b945220601f850160051c81019160208610610832575b601f0160051c0190614460565b8b5f610549565b909150819061081e565b634e487b7160e01b8c52604160045260248cfd5b9192601f1984165f5160206147545f395f51905f528452828420935b8181106108b4575090846001959493921061089c575b505050811b015f5160206147545f395f51905f5255610517565b01515f1960f88460031b161c191690555f8080610882565b9293602060018192878601518155019501930161086c565b6020825f5160206147545f395f51905f52610902945220601f840160051c8101916020851061083257601f0160051c0190614460565b8c5f6104d8565b634e487b7160e01b8d52604160045260248dfd5b986001019861040a565b6064600291049b019a610403565b612710600491049b019a6103f9565b6305f5e100600891049b019a6103ee565b50601099508a6103de565b632582a64160e11b8852600488fd5b5081156103c8565b9050155f6103c0565b6001600160a01b0385161591506103b9565b821591506103b2565b6801000000000000000060ff60401b195f5160206148145f395f51905f525416175f5160206148145f395f51905f52556103a8565b63f92ee8a960e01b8952600489fd5b9050155f610372565b303b15915061036a565b899150610360565b633921c70360e01b8852600488fd5b8780fd5b8680fd5b8480fd5b5080fd5b503461022a57602036600319011261022a57610a33612a5b565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614610a9657610a6b6133ef565b6001600160a01b03811615610a8757610a839061397d565b5080f35b632582a64160e11b8252600482fd5b633921c70360e01b8252600482fd5b503461022a57602036600319011261022a576004356001600160401b038111610a155736602382011215610a15578060040135906001600160401b0382116110ba5760248260051b82010190368211610f5457307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316146110ab57610b3061344b565b610b3861382e565b610b40613fcf565b610b4983612b80565b91610b576040519384612b5f565b8383526024602084019201915b81831061108b575050509015610f63575b60115460405163024ece8960e01b8152908390829060049082906001600160a01b03165afa908115610f58578391610eb8575b50815181518493845b838103610bd457868660035401600355805f5160206147f45f395f51905f525d80f35b6001600160401b03610be68284612bf5565b5116610c03816001600160401b03165f52601460205260405f2090565b5f805260205260405f205480610e42575b50875b848103610d4c57506001600160401b0381165f525f60205260405f2054610c4b575b90610c456001926145a2565b01610bb1565b956001600160401b0387165f525f6020528760405f2055600154885b818103610c7a575b505060010195610c39565b886001600160401b03610c8c836144ae565b90549060031b1c1614610ca157600101610c67565b905f198101908111610d38576001600160401b03610cc1610ccf926144ae565b90549060031b1c16916144ae565b6001600160401b03829392549160031b92831b921b19161790556001548015610d24576001809392610c45925f1901610d07816144ae565b6001600160401b0382549160031b1b191690558255929350610c6f565b634e487b7160e01b89526031600452602489fd5b634e487b7160e01b8a52601160045260248afd5b600190610d6a836001600160401b03165f52601460205260405f2090565b828060a01b03610d7a838b612bf5565b5116838060a01b03165f5260205260405f205480610d9a575b5001610c17565b610db781338b610db086888060a01b0392612bf5565b511661441f565b610dd2846001600160401b03165f52601460205260405f2090565b838060a01b03610de2848c612bf5565b5116848060a01b03165f5260205260405f20610dff828254612be8565b9055610e3a610e32848060a01b03610e17858d612bf5565b51166001600160a01b03165f90815260156020526040902090565b918254612be8565b90555f610d93565b8880808084335af1610e52613863565b5015610ea957610e73826001600160401b03165f52601460205260405f2090565b5f805260205260405f20610e88828254612be8565b90558880526015602052610ea160408a20918254612be8565b90555f610c14565b6312171d8360e31b8952600489fd5b90503d8084833e610ec98183612b5f565b810190602081830312610f54578051906001600160401b038211610a1157019080601f83011215610f54578151610eff81612b80565b92610f0d6040519485612b5f565b81845260208085019260051b820101928311610f5057602001905b828210610f38575050505f610ba8565b60208091610f4584612cb9565b815201910190610f28565b8580fd5b8380fd5b6040513d85823e3d90fd5b50604051808160206001549283815201600185527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf69285905b80600383011061103957610fd294549181811061101f575b818110611002575b818110610fe5575b10610fd7575b500382612b5f565b610b75565b60c01c81526020015f610fca565b9260206001916001600160401b038560801c168152019301610fc4565b9260206001916001600160401b038560401c168152019301610fbc565b9260206001916001600160401b0385168152019301610fb4565b916004919350608060019186546001600160401b03811682526001600160401b038160401c1660208301526001600160401b0381841c16604083015260c01c6060820152019401920184929391610f9c565b82356001600160401b0381168103610a0d57815260209283019201610b64565b633921c70360e01b8452600484fd5b8280fd5b503461022a57602036600319011261022a576020906040906001600160a01b036110e6612a5b565b168152601883522054604051908152f35b503461022a57604036600319011261022a57610a83600435611117612a71565b9061114261113d825f525f5160206147d45f395f51905f52602052600160405f20015490565b6134ba565b613a0e565b503461022a57602036600319011261022a57611161612a5b565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614610a96576111996133ef565b6001600160a01b03168015610a87576020817fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f926001600160601b0360a01b6017541617601755604051908152a180f35b503461022a578060031936011261022a576017546040516001600160a01b039091168152602090f35b50610140366003190112611a8d57611229612a5b565b6024359160ff83168303611a8d576044356001600160401b0381168103611a8d5760056064351015611a8d576084356001600160401b038111611a8d57611274903690600401612ac3565b9060a435906001600160a01b0382168203611a8d57610104356001600160401b038111611a8d576112a9903690600401612ac3565b90610124356001600160401b038111611a8d576112ca903690600401612ac3565b9092909190307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f5760ff8b16611d10576001600160401b0388165f525f60205260405f205415611d015761132961382e565b5f6003606435141580611cf2575b611ce35760e4354211611cd457611358600d546113526144ee565b90612be8565b600e541115611cc557611c7b5786918a8c600360643514611c8f575b9160426114c8926114ce959460018060a01b0382165f5260186020528d6001600160401b036113aa8d60405f20549b3691612bb2565b602081519101209160ff6040519460208601967f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc1333340885260018060a01b0316604087015216606085015216608083015260ff6064351660a083015260c082015260018060a01b038c1660e082015260c4356101008201528761012082015260e435610140820152610140815261144161016082612b5f565b51902061144c614684565b6114546146ee565b6040519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f8452604083015260608201524660808201523060a082015260a081526114a560c082612b5f565b519020906040519161190160f01b83526002830152602282015220923691612bb2565b906143e5565b506004811015611c7b5715801590611c6a575b611c5b576001600160a01b038a8116911603611c4c5760018101809111611c38576001600160a01b0389165f908152601860205260409020556016543410611c295760c435158080611c17575b611c0857156118b1575b5050611545368483612bb2565b602081519101206001600160401b03851686526005602052600260408720015460405190602082019260018060a01b038a1684526001600160401b03881660408401526005606435101561189d579287928a9b95928a97956064356060840152608083015260018060a01b03861660a083015260c43560c083015260e082015260e081526115d561010082612b5f565b5190209889936001600160401b038416875260056020526040872097604051936115fe85612b43565b428552602085019260018060a01b031683526040850160c4358152611633606087019234845260808801948a86523691612bb2565b9360a0870194855260c087019560018060a01b03168652604060e088019b338d526001600160401b036101008a019a168a5260ff6101208a019c168c5261014089019a6064358c5281526004602052209651875560018060a01b03905116600187019060018060a01b03166001600160601b0360a01b8254161790555160028601555160038501555160048401556005830190518051906001600160401b03821161083c576116ec826116e68554613dc8565b85614476565b6020908c601f841160011461183657600796959493611720939092836107625750508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b1691600581101561182257918895939160209a959360ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b191617171790556002810190815486528752826040862055600181540190556001600d5401600d557fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236866001600160401b036040519333855260018060a01b0316951692a45f5160206147f45f395f51905f525d604051908152f35b634e487b7160e01b88526021600452602488fd5b9190601f198416858452828420935b8181106118855750916001939185600799989796941061186d575b505050811b019055611723565b01515f1960f88460031b161c191690555f8080611860565b92936020600181928786015181550195019301611845565b634e487b7160e01b89526021600452602489fd5b6001600160a01b03841615611c085760115460405163cbe230c360e01b81526001600160a01b038681166004830152909160209183916024918391165afa908115611b72575f91611bcd575b5015611bbe5760608103611baf578160609181010312611a8d5780359060ff8216809203611a8d57604051636eb1769f60e11b81526001600160a01b0389811660048301523060248301526020908290604490829089165afa908115611b72575f91611b7d575b5060c43511611ae1575b50506040516370a0823160e01b81523060048201526020816024816001600160a01b0387165afa908115611ad6578691611aa4575b506040516323b872dd60e01b60208201526001600160a01b038816602482015230604482015260c43560648083019190915281526119f5906119e6608482612b5f565b6001600160a01b03851661462c565b6040516370a0823160e01b8152306004820152906020826024816001600160a01b0388165afa918215611a99578792611a5f575b50611a379060c43592612be8565b03611a5057611a4960c4358386613f7a565b5f80611538565b63b0a6ea2960e01b8552600485fd5b9091506020813d602011611a91575b81611a7b60209383612b5f565b81010312611a8d575190611a37611a29565b5f80fd5b3d9150611a6e565b6040513d89823e3d90fd5b90506020813d602011611ace575b81611abf60209383612b5f565b81010312611a8d57515f6119a3565b3d9150611ab2565b6040513d88823e3d90fd5b6001600160a01b0384163b15611a8d5760409081519263d505accf60e01b845260018060a01b038a16600485015230602485015260c435604485015260e43560648501526084840152602081013560a4840152013560c48201525f8160e4818360018060a01b0388165af18015611b7257611b5d575b8061196e565b611b6a9195505f90612b5f565b5f935f611b57565b6040513d5f823e3d90fd5b90506020813d602011611ba7575b81611b9860209383612b5f565b81010312611a8d57515f611964565b3d9150611b8b565b63ddafbaef60e01b5f5260045ffd5b63514e24c360e11b5f5260045ffd5b90506020813d602011611c00575b81611be860209383612b5f565b81010312611a8d57518015158103611a8d575f6118fd565b3d9150611bdb565b632a9ffab760e21b5f5260045ffd5b506001600160a01b038516151561152e565b63c77f97eb60e01b5f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b632057875960e21b5f5260045ffd5b638baa579f60e01b5f5260045ffd5b506001600160a01b038116156114e1565b634e487b7160e01b5f52602160045260245ffd5b505090916085141580611cba575b611cab579086918a8c611374565b637c6953f960e01b5f5260045ffd5b5060e2871415611c9d565b63d8219b3d60e01b5f5260045ffd5b631ab7da6b60e01b5f5260045ffd5b630ee2db0560e21b5f5260045ffd5b50505f60016064351415611337565b6309ff9a8360e41b5f5260045ffd5b635428eae760e01b5f5260045ffd5b633921c70360e01b5f5260045ffd5b34611a8d576040366003190112611a8d57611d47612a9d565b6001600160401b03611d57612a71565b91165f52601460205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b34611a8d576020366003190112611a8d576001600160a01b03611da2612a5b565b165f526015602052602060405f2054604051908152f35b34611a8d575f366003190112611a8d5760206040515f8152f35b34611a8d575f366003190112611a8d576010546040516001600160a01b039091168152602090f35b34611a8d575f366003190112611a8d576020601654604051908152f35b34611a8d576020366003190112611a8d57600435307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f57611e646133ef565b8015611c08576020817e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db192600e55604051908152a1005b34611a8d576040366003190112611a8d57611eb3612a71565b6004355f525f5160206147d45f395f51905f5260205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b34611a8d576020366003190112611a8d57600435307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f5760207f38c762d70caa5dcb8a3b0bcbcc5883d51f55c2aefd8b02fdbad6b2871a5a691691611f5f6133ef565b80601b55604051908152a1005b34611a8d575f366003190112611a8d576020600254604051908152f35b34611a8d575f366003190112611a8d57307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f57611fd161344b565b611fd961382e565b611fe1613fcf565b5f5f5160206147f45f395f51905f525d005b34611a8d575f366003190112611a8d5760206040517fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c959178152f35b34611a8d576080366003190112611a8d57612046612a9d565b6044356001600160401b038111611a8d5760406003198236030112611a8d576040519061207282612b14565b80600401356001600160401b038111611a8d57810136602382011215611a8d5760048101356120a081612b80565b916120ae6040519384612b5f565b818352602060048185019360051b8301010190368211611a8d5760248101925b828410612215575050505082526024810135906001600160401b038211611a8d570136602382011215611a8d5760048101359061210a82612b80565b916121186040519384612b5f565b808352602060048185019260051b8401010191368311611a8d57602401905b828210612205575050506020820152606435906001600160401b038211611a8d5736602383011215611a8d57816004013561217181612b80565b9261217f6040519485612b5f565b8184526024602085019260051b82010190368211611a8d57602401915b8183106121e5575050506001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163014611d1f576121e39260243590612d9c565b005b82356001600160a01b0381168103611a8d5781526020928301920161219c565b8135815260209182019101612137565b83356001600160401b038111611a8d5760049083010136603f82011215611a8d5760209161224d839236906040858201359101612bb2565b8152019301926120ce565b34611a8d575f366003190112611a8d577fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100541580612379575b1561233c576122e06122a1613e00565b6122a9613ecd565b60206122ee604051926122bc8385612b5f565b5f84525f368137604051958695600f60f81b875260e08588015260e0870190612af0565b908582036040870152612af0565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b82811061232557505050500390f35b835185528695509381019392810192600101612316565b60405162461bcd60e51b81526020600482015260156024820152741152540dcc4c8e88155b9a5b9a5d1a585b1a5e9959605a1b6044820152606490fd5b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1015415612291565b34611a8d576040366003190112611a8d576123bb612a5b565b6123c3612a71565b6001600160a01b039182165f908152601260209081526040808320949093168252928352819020549051908152f35b34611a8d576020366003190112611a8d576001600160401b03612413612a9d565b165f525f602052602060405f2054604051908152f35b34611a8d575f366003190112611a8d57600f546040516001600160a01b039091168152602090f35b6060366003190112611a8d57612465612ab3565b6024356001600160401b038111611a8d57612484903690600401612ac3565b61248f929192612a87565b92307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f5760ff8316611d10576124d9926124d461382e565b613aaa565b6001600160a01b0382169182612505575b6020825f5f5160206147f45f395f51905f525d604051908152f35b825f52601a6020526001600160401b0360405f205416612584573b156125755760c081901c5f81815260196020908152604080832080546001600160a01b0319166001600160a01b038816179055948252601a8152939020805467ffffffffffffffff19169091179055826124ea565b634ea29e6760e01b5f5260045ffd5b63d48c379760e01b5f5260045ffd5b34611a8d576020366003190112611a8d57600435307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f576125df6133ef565b8015611c08576002546125f460035482612be8565b90818310611c08578261262e6040937f80af1d8add7b822532307ae28ef73c31ab2803f908c5c22cc879766f54c685839560025582612be8565b60035582519182526020820152a1005b6040366003190112611a8d57612652612ab3565b6024356001600160401b038111611a8d57612671903690600401612ac3565b909190307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f5760ff8216611d10576020926126bb926124d461382e565b5f5f5160206147f45f395f51905f525d604051908152f35b34611a8d575f366003190112611a8d576020601b54604051908152f35b34611a8d576040366003190112611a8d57612709612a71565b336001600160a01b03821603612725576121e390600435613a0e565b63334bd91960e11b5f5260045ffd5b34611a8d576040366003190112611a8d576121e3600435612753612a71565b9061277961113d825f525f5160206147d45f395f51905f52602052600160405f20015490565b61379d565b34611a8d575f366003190112611a8d5760206040517f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc13333408152f35b34611a8d575f366003190112611a8d5760206040515f5160206147b45f395f51905f528152f35b34611a8d576020366003190112611a8d5760206128166004355f525f5160206147d45f395f51905f52602052600160405f20015490565b604051908152f35b34611a8d576040366003190112611a8d57612837612a5b565b61283f612a71565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f57611fe19161287b61382e565b613892565b34611a8d576020366003190112611a8d57612899612a5b565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611d1f576128d16133ef565b6001600160a01b038116156128e9576121e390613500565b632582a64160e11b5f5260045ffd5b34611a8d576020366003190112611a8d576001600160401b03612919612a9d565b165f526019602052602060018060a01b0360405f205416604051908152f35b34611a8d576020366003190112611a8d576001600160a01b03612959612a5b565b165f526013602052602060405f2054604051908152f35b34611a8d576020366003190112611a8d576001600160a01b03612991612a5b565b165f52601a60205260206001600160401b0360405f205416604051908152f35b34611a8d575f366003190112611a8d576020600354604051908152f35b34611a8d575f366003190112611a8d5760206040517f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd988152f35b34611a8d576020366003190112611a8d576004359063ffffffff60e01b8216809203611a8d57602091637965db0b60e01b8114908115612a4a575b5015158152f35b6301ffc9a760e01b14905083612a43565b600435906001600160a01b0382168203611a8d57565b602435906001600160a01b0382168203611a8d57565b604435906001600160a01b0382168203611a8d57565b600435906001600160401b0382168203611a8d57565b6004359060ff82168203611a8d57565b9181601f84011215611a8d578235916001600160401b038311611a8d5760208381860195010111611a8d57565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b604081019081106001600160401b03821117612b2f57604052565b634e487b7160e01b5f52604160045260245ffd5b61016081019081106001600160401b03821117612b2f57604052565b90601f801991011681019081106001600160401b03821117612b2f57604052565b6001600160401b038111612b2f5760051b60200190565b6001600160401b038111612b2f57601f01601f191660200190565b929192612bbe82612b97565b91612bcc6040519384612b5f565b829481845281830111611a8d578281602093845f960137010152565b91908203918211611c3857565b8051821015612c095760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b60408201908051916040845282518091526060840190602060608260051b8701019401915f905b828210612c8e575050505060200151916020818303910152602080835192838152019201905f5b818110612c785750505090565b8251845260209384019390920191600101612c6b565b90919294602080612cab600193605f198b82030186528951612af0565b970192019201909291612c44565b51906001600160a01b0382168203611a8d57565b81601f82011215611a8d57805190612ce482612b80565b92612cf26040519485612b5f565b82845260208085019360061b83010191818311611a8d57602001925b828410612d1c575050505090565b604084830312611a8d5760206040918251612d3681612b14565b612d3f87612cb9565b81528287015183820152815201930192612d0e565b90602080835192838152019201905f5b818110612d715750505090565b825180516001600160a01b031685526020908101518186015260409094019390920191600101612d64565b91905f906001600160401b03841692835f52601960205260018060a01b0360405f2054169586156133d9578051905f5b8281036133b3575050506001863b15611a8d576040516338d9814d60e11b8152602060048201525f8180612e036024820187612c1d565b0381838c5af1908161339e575b506133995750825b82857fca435620e7fa2ff484abda442975d2167d5e6577548dfdfac799cecc2a2af2c9602060405194151594858152a3604051633ccfd60b60e01b815260609283916001908781600481838f5af1908189918a9361332f575b506133235750505085975b845190875b8281036132ec57505050612ee3866060998b612ec283600196612ef560405197889687958694635205569560e11b865260a0600487015260a4860190612c1d565b9c602485015215159b8c60448501526003198482030160648501528d612d54565b8281036003190160848401528a612d54565b03925af1879181613269575b50613234575050917fa59edce59a4211e20e9f86e63adea41a55dbe409040bc17f473b671e31bd036a91612f648794612f5688945b604051958695865215156020860152608060408601526080850190612d54565b908382036060850152612d54565b0390a3825180612f76575b5050505050565b6020840120600b546040519060208201928784528560408401526004606084015260808301528360a08301528360c083015260e082015260e08152612fbd61010082612b5f565b51902092604051612fcd81612b43565b4281526020810183815260408201918483526060810190858252608081019488865260a0820190815260c08201908a825260e08301968888526101008401948a86526101208501978a8952610140860197600489528d8c52600460205260408c209651875560018060a01b03905116600187019060018060a01b03166001600160601b0360a01b8254161790555160028601555160038501555160048401556005830190518051906001600160401b03821161322057613091826116e68554613dc8565b602090601f83116001146131b95791806130c5926007979695948d926107625750508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b169160058110156131a55791859391602095937fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236975060ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600b548152600982528460408220556001600b5401600b556001600d5401600d55604051908152a45f80808080612f6f565b634e487b7160e01b86526021600452602486fd5b838b52818b209190601f1984168c5b818110613208575091600193918560079998979694106131f0575b505050811b0190556130c8565b01515f1960f88460031b161c191690555f80806131e3565b929360206001819287860151815501950193016131c8565b634e487b7160e01b8a52604160045260248afd5b879498507fa59edce59a4211e20e9f86e63adea41a55dbe409040bc17f473b671e31bd036a939192612f56612f64929a612f36565b9091503d8089833e61327b8183612b5f565b8101906020818303126132e8578051906001600160401b0382116132e4570181601f820112156132e8578051906132b182612b97565b926132bf6040519485612b5f565b828452602083830101116132e457818a9260208093018386015e83010152905f612f01565b8980fd5b8880fd5b60019061331d6001600160a01b03613304838b612bf5565b5151166020613313848c612bf5565b5101519085613f7a565b01612e81565b91999195509250612e7c565b915091503d808a833e6133428183612b5f565b81016040828203126132e45781516001600160401b038111613395578161336a918401612ccd565b916020810151906001600160401b0382116133915761338a929101612ccd565b915f612e71565b8b80fd5b8a80fd5b612e18565b6133ab9195505f90612b5f565b5f935f612e10565b6001906133d38a6001600160a01b036133cc8487612bf5565b5116613892565b01612dcc565b50505050505050565b91908201809211611c3857565b335f9081527f78e571b7bf30584d955e1c6444a2b5147087edf9f00485d94993a04d370525ea602052604090205460ff161561342757565b63e2517d3f60e01b5f52336004525f5160206147b45f395f51905f5260245260445ffd5b335f9081527f143ba17de5dacef81d11036b41c596cd787c462f0e59fb1f76aa2b829d603e0c602052604090205460ff161561348357565b63e2517d3f60e01b5f52336004527fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c9591760245260445ffd5b5f8181525f5160206147d45f395f51905f526020908152604080832033845290915290205460ff16156134ea5750565b63e2517d3f60e01b5f523360045260245260445ffd5b6001600160a01b0381165f9081525f5160206148345f395f51905f52602052604090205460ff16613580576001600160a01b03165f8181525f5160206148345f395f51905f5260205260408120805460ff191660011790553391905f5160206147945f395f51905f52905f5160206147345f395f51905f529080a4600190565b505f90565b6001600160a01b0381165f9081527f212403f86ff12f15ffbeab879ed2f25237f7273bd0873f6abf8e9330ca4ed4b9602052604090205460ff16613580576001600160a01b03165f8181527f212403f86ff12f15ffbeab879ed2f25237f7273bd0873f6abf8e9330ca4ed4b960205260408120805460ff191660011790553391907f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd98905f5160206147345f395f51905f529080a4600190565b6001600160a01b0381165f9081527f78e571b7bf30584d955e1c6444a2b5147087edf9f00485d94993a04d370525ea602052604090205460ff16613580576001600160a01b03165f8181527f78e571b7bf30584d955e1c6444a2b5147087edf9f00485d94993a04d370525ea60205260408120805460ff191660011790553391905f5160206147b45f395f51905f52905f5160206147345f395f51905f529080a4600190565b6001600160a01b0381165f9081527f143ba17de5dacef81d11036b41c596cd787c462f0e59fb1f76aa2b829d603e0c602052604090205460ff16613580576001600160a01b03165f8181527f143ba17de5dacef81d11036b41c596cd787c462f0e59fb1f76aa2b829d603e0c60205260408120805460ff191660011790553391907fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c95917905f5160206147345f395f51905f529080a4600190565b5f8181525f5160206147d45f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff16613828575f8181525f5160206147d45f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f5160206147345f395f51905f529080a4600190565b50505f90565b5f5160206147f45f395f51905f525c6138545760015f5160206147f45f395f51905f525d565b633ee5aeb560e01b5f5260045ffd5b3d1561388d573d9061387482612b97565b916138826040519384612b5f565b82523d5f602084013e565b606090565b60018060a01b03811691825f52601260205260405f2060018060a01b0382165f5260205260405f2054918215613977577f16d46b28f1a9377c2df78652a0a8e31568b99311665bcec24cbba72edd2e5c2c8392855f52601260205260405f2060018060a01b0382165f526020525f6040812055855f52601360205260405f2061391c858254612be8565b9055604080516001600160a01b0394851681526020810195909552921692839290a28261396c575f809350809281925af1613955613863565b501561395d57565b6312171d8360e31b5f5260045ffd5b6139759261441f565b565b50505050565b6001600160a01b0381165f9081525f5160206148345f395f51905f52602052604090205460ff1615613580576001600160a01b03165f8181525f5160206148345f395f51905f5260205260408120805460ff191690553391905f5160206147945f395f51905f52907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b5f8181525f5160206147d45f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff1615613828575f8181525f5160206147d45f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b335f9081525f5160206148345f395f51905f5260205260409020549092919060ff1615613db9576003548015613daa57613ae8600d546113526144ee565b600e541115611cc5576016543410611c29575f1901600355613b0b368383612bb2565b602081519101206008546040519060208201923384525f60408401525f606084015260808301525f60a08301525f60c083015260e082015260e08152613b5361010082612b5f565b519020928360c01c9260405190613b6982612b43565b42825260208201915f835260408101915f8352613b96606083019634885260808401928a84523691612bb2565b60a0830190815260c083019133835260e08401975f89526101008501958a875260ff6101208701991689526101408601975f89528c5f52600460205260405f209651875560018060a01b03905116600187019060018060a01b03166001600160601b0360a01b8254161790555160028601555160038501555160048401556005830190518051906001600160401b038211612b2f57613c39826116e68554613dc8565b602090601f8311600114613d42579180613c6d926007979695945f926107625750508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b16916005811015611c7b5760ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b191617171790556008545f5260066020528160405f20556001600854016008556001600d5401600d55604051908282527fbbcfb0c180d5f95cfb50c12005c44144ac92ac8474f1abc5add0e6e69a14e8dc60203393a390565b90601f19831691845f52815f20925f5b818110613d9257509160019391856007999897969410613d7a575b505050811b019055613c70565b01515f1960f88460031b161c191690555f8080613d6d565b92936020600181928786015181550195019301613d52565b6316ae067f60e11b5f5260045ffd5b63fc336c4160e01b5f5260045ffd5b90600182811c92168015613df6575b6020831014613de257565b634e487b7160e01b5f52602260045260245ffd5b91607f1691613dd7565b604051905f825f5160206147545f395f51905f525491613e1f83613dc8565b8083529260018116908115613eae5750600114613e43575b61397592500383612b5f565b505f5160206147545f395f51905f525f90815290917f42ad5d3e1f2e6e70edcf6d991b8a3023d3fca8047a131592f9edb9fd9b89d57d5b818310613e9257505090602061397592820101613e37565b6020919350806001915483858901015201910190918492613e7a565b6020925061397594915060ff191682840152151560051b820101613e37565b604051905f825f5160206147745f395f51905f525491613eec83613dc8565b8083529260018116908115613eae5750600114613f0f5761397592500383612b5f565b505f5160206147745f395f51905f525f90815290917f5f9ce34815f8e11431c7bb75a8e6886a91478f7ffc1dbb0a98dc240fddd76b755b818310613f5e57505090602061397592820101613e37565b6020919350806001915483858901015201910190918492613f46565b6001600160401b03165f52601460205260405f2060018060a01b0382165f5260205260405f20613fab8382546133e2565b905560018060a01b03165f526015602052613fcb60405f209182546133e2565b9055565b613fd761450d565b600754600854905b8181036143ac57505060075460085490808203600d5403600d555b81810361430057505090600754600855600154915f925b8084036140f65750909150600a54600b5490808203600d5403600d555b81810361404a575050600a54600b555f600c5560035401600355565b806001915f52600960205260405f20545f5260046020525f6007604082208281558285820155826002820155826003820155826004820155600581016140908154613dc8565b90816140b5575b50508260068201550155805f5260096020525f60408120550161402e565b81601f86931188146140cb5750555b5f80614097565b818352602083206140e691601f0160051c8101908901614460565b80825281602081209155556140c4565b6001600160401b03614107856144ae565b90549060031b1c165f52600560205260405f20936001850180549060028701918254905b8181036142045750508054825490808203600d5403600d555b81810361415a5750505490556001019350614011565b806001915f528960205260405f20545f5260046020525f60076040822082815582858201558260028201558260038201558260048201556005810161419f8154613dc8565b90816141c3575b50508260068201550155805f52896020525f604081205501614144565b81601f86931188146141d95750555b5f806141a6565b818352602083206141f491601f0160051c8101908901614460565b80825281602081209155556141d2565b806001915f528960205260405f20545f52600460205260405f2060028101549081614232575b50500161412b565b838060a01b038482015416906142696001600160401b03600783015460a01c166001600160401b03165f52601460205260405f2090565b858060a01b0383165f5260205260405f20614285848254612be8565b90556001600160a01b0382165f9081526015602052604090206142a9848254612be8565b9055815f52601260205260405f20906006868060a01b0391015416858060a01b03165f5260205260405f206142df8382546133e2565b90555f5260136020526142f760405f209182546133e2565b90555f8061422a565b806001915f52600660205260405f20545f5260046020525f6007604082208281558285820155826002820155826003820155826004820155600581016143468154613dc8565b908161436b575b50508260068201550155805f5260066020525f604081205501613ffa565b81601f86931188146143815750555b5f8061434d565b8183526020832061439c91601f0160051c8101908901614460565b808252816020812091555561437a565b806001915f52600660205260405f20545f5260046020526143df6001600160401b03600760405f20015460a01c166145a2565b01613fdf565b81519190604183036144155761440e9250602082015190606060408401519301515f1a90614520565b9192909190565b50505f9160029190565b60405163a9059cbb60e01b60208201526001600160a01b039290921660248301526044808301939093529181526139759161445b606483612b5f565b61462c565b81811061446b575050565b5f8155600101614460565b9190601f811161448557505050565b613975925f5260205f20906020601f840160051c8301931061083257601f0160051c0190614460565b90600154821015612c095760015f52600282901c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6019160031b60181690565b600b54600a548082116145015750505f90565b61450a91612be8565b90565b6008546007548082116145015750505f90565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411614597579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15611b72575f516001600160a01b0381161561458d57905f905f90565b505f906001905f90565b5050505f9160039190565b6001600160401b03165f818152601960205260409020546001600160a01b0316806145cb575050565b5f52601a60205260405f206001600160401b031981541690555f52601960205260405f206001600160601b0360a01b8154169055565b60ff5f5160206148145f395f51905f525460401c161561461d57565b631afcd79f60e31b5f5260045ffd5b905f602091828151910182855af115611b72575f513d61467b57506001600160a01b0381163b155b61465b5750565b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b60011415614654565b61468c613e00565b805190811561469c576020012090565b50507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1005480156146c95790565b507fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47090565b6146f6613ecd565b8051908115614706576020012090565b50507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1015480156146c9579056fe2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0da16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d102a16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d103fc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184cdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4202dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268009b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008f7e9a8734279607aeaac474f24949ddc99eae25c0da72edcf5ad21022be2c9ea26469706673582212201e4da01a108a3f330ec4e2a55c00d917d63bfe0333a9de48e4af52baf600b6bc64736f6c634300081e0033",
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

// PackAddAllowedDeployer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x20ff473f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addAllowedDeployer(address deployer) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackAddAllowedDeployer(deployer common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("addAllowedDeployer", deployer)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddAllowedDeployer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x20ff473f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addAllowedDeployer(address deployer) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackAddAllowedDeployer(deployer common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("addAllowedDeployer", deployer)
}

// PackAdminReset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c5b9b00.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function adminReset() returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackAdminReset() []byte {
	enc, err := processorEndpointExtension.abi.Pack("adminReset")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAdminReset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c5b9b00.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function adminReset() returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackAdminReset() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("adminReset")
}

// PackAdminResetApps is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd39d765.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function adminResetApps(uint64[] appIds) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackAdminResetApps(appIds []uint64) []byte {
	enc, err := processorEndpointExtension.abi.Pack("adminResetApps", appIds)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAdminResetApps is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd39d765.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function adminResetApps(uint64[] appIds) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackAdminResetApps(appIds []uint64) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("adminResetApps", appIds)
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

// PackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21c0b342.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claim(address tokenAddress, address payee) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackClaim(tokenAddress common.Address, payee common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("claim", tokenAddress, payee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21c0b342.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claim(address tokenAddress, address payee) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackClaim(tokenAddress common.Address, payee common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("claim", tokenAddress, payee)
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

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xeabcca10.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address _teeAuthenticator, address _authorityRegistry, address updateStatusOperator, address admin, address resetOperator, uint256 _minFeePerRequest, address _tokenAllowlist) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackInitialize(teeAuthenticator common.Address, authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address, resetOperator common.Address, minFeePerRequest *big.Int, tokenAllowlist common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("initialize", teeAuthenticator, authorityRegistry, updateStatusOperator, admin, resetOperator, minFeePerRequest, tokenAllowlist)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xeabcca10.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address _teeAuthenticator, address _authorityRegistry, address updateStatusOperator, address admin, address resetOperator, uint256 _minFeePerRequest, address _tokenAllowlist) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackInitialize(teeAuthenticator common.Address, authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address, resetOperator common.Address, minFeePerRequest *big.Int, tokenAllowlist common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("initialize", teeAuthenticator, authorityRegistry, updateStatusOperator, admin, resetOperator, minFeePerRequest, tokenAllowlist)
}

// PackInvokeTrigger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x85010384.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function invokeTrigger(uint64 applicationId, bytes32 processedRequestId, (bytes[],bytes32[]) appEventData, address[] claimable) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackInvokeTrigger(applicationId uint64, processedRequestId [32]byte, appEventData StructsEventData, claimable []common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("invokeTrigger", applicationId, processedRequestId, appEventData, claimable)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInvokeTrigger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x85010384.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function invokeTrigger(uint64 applicationId, bytes32 processedRequestId, (bytes[],bytes32[]) appEventData, address[] claimable) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackInvokeTrigger(applicationId uint64, processedRequestId [32]byte, appEventData StructsEventData, claimable []common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("invokeTrigger", applicationId, processedRequestId, appEventData, claimable)
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

// PackRemoveAllowedDeployer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3d72e5e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function removeAllowedDeployer(address deployer) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackRemoveAllowedDeployer(deployer common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("removeAllowedDeployer", deployer)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRemoveAllowedDeployer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3d72e5e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function removeAllowedDeployer(address deployer) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackRemoveAllowedDeployer(deployer common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("removeAllowedDeployer", deployer)
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

// PackSelectionGrace is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x44ce30f0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function selectionGrace() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) PackSelectionGrace() []byte {
	enc, err := processorEndpointExtension.abi.Pack("selectionGrace")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSelectionGrace is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x44ce30f0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function selectionGrace() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackSelectionGrace() ([]byte, error) {
	return processorEndpointExtension.abi.Pack("selectionGrace")
}

// UnpackSelectionGrace is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x44ce30f0.
//
// Solidity: function selectionGrace() view returns(uint256)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackSelectionGrace(data []byte) (*big.Int, error) {
	out, err := processorEndpointExtension.abi.Unpack("selectionGrace", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSubmitDeployRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4a0d0495.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitDeployRequest(uint8 protocolVersion, bytes payload) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackSubmitDeployRequest(protocolVersion uint8, payload []byte) []byte {
	enc, err := processorEndpointExtension.abi.Pack("submitDeployRequest", protocolVersion, payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitDeployRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4a0d0495.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitDeployRequest(uint8 protocolVersion, bytes payload) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackSubmitDeployRequest(protocolVersion uint8, payload []byte) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("submitDeployRequest", protocolVersion, payload)
}

// UnpackSubmitDeployRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4a0d0495.
//
// Solidity: function submitDeployRequest(uint8 protocolVersion, bytes payload) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackSubmitDeployRequest(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("submitDeployRequest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSubmitDeployRequestWithTrigger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68a7c3fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitDeployRequestWithTrigger(uint8 protocolVersion, bytes payload, address trigger) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) PackSubmitDeployRequestWithTrigger(protocolVersion uint8, payload []byte, trigger common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("submitDeployRequestWithTrigger", protocolVersion, payload, trigger)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitDeployRequestWithTrigger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68a7c3fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitDeployRequestWithTrigger(uint8 protocolVersion, bytes payload, address trigger) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackSubmitDeployRequestWithTrigger(protocolVersion uint8, payload []byte, trigger common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("submitDeployRequestWithTrigger", protocolVersion, payload, trigger)
}

// UnpackSubmitDeployRequestWithTrigger is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68a7c3fb.
//
// Solidity: function submitDeployRequestWithTrigger(uint8 protocolVersion, bytes payload, address trigger) payable returns(bytes32)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackSubmitDeployRequestWithTrigger(data []byte) ([32]byte, error) {
	out, err := processorEndpointExtension.abi.Unpack("submitDeployRequestWithTrigger", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
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

// PackUpdateFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2c35ce8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackUpdateFeeCollector(newFeeCollector common.Address) []byte {
	enc, err := processorEndpointExtension.abi.Pack("updateFeeCollector", newFeeCollector)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2c35ce8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackUpdateFeeCollector(newFeeCollector common.Address) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("updateFeeCollector", newFeeCollector)
}

// PackUpdateMaxNumOfApplications is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x655f2de1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateMaxNumOfApplications(uint256 newMax) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackUpdateMaxNumOfApplications(newMax *big.Int) []byte {
	enc, err := processorEndpointExtension.abi.Pack("updateMaxNumOfApplications", newMax)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateMaxNumOfApplications is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x655f2de1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateMaxNumOfApplications(uint256 newMax) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackUpdateMaxNumOfApplications(newMax *big.Int) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("updateMaxNumOfApplications", newMax)
}

// PackUpdateQueueThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c73eeaf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateQueueThreshold(uint256 newThreshold) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackUpdateQueueThreshold(newThreshold *big.Int) []byte {
	enc, err := processorEndpointExtension.abi.Pack("updateQueueThreshold", newThreshold)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateQueueThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c73eeaf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateQueueThreshold(uint256 newThreshold) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackUpdateQueueThreshold(newThreshold *big.Int) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("updateQueueThreshold", newThreshold)
}

// PackUpdateSelectionGrace is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x90fda89b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateSelectionGrace(uint256 newGrace) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) PackUpdateSelectionGrace(newGrace *big.Int) []byte {
	enc, err := processorEndpointExtension.abi.Pack("updateSelectionGrace", newGrace)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateSelectionGrace is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x90fda89b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateSelectionGrace(uint256 newGrace) returns()
func (processorEndpointExtension *ProcessorEndpointExtension) TryPackUpdateSelectionGrace(newGrace *big.Int) ([]byte, error) {
	return processorEndpointExtension.abi.Pack("updateSelectionGrace", newGrace)
}

// ProcessorEndpointExtensionDeployRequestSubmitted represents a DeployRequestSubmitted event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionDeployRequestSubmitted struct {
	ApplicationId uint64
	RequestId     [32]byte
	Sender        common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionDeployRequestSubmittedEventName = "DeployRequestSubmitted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionDeployRequestSubmitted) ContractEventName() string {
	return ProcessorEndpointExtensionDeployRequestSubmittedEventName
}

// UnpackDeployRequestSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DeployRequestSubmitted(uint64 indexed applicationId, bytes32 requestId, address indexed sender)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackDeployRequestSubmittedEvent(log *types.Log) (*ProcessorEndpointExtensionDeployRequestSubmitted, error) {
	event := "DeployRequestSubmitted"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionDeployRequestSubmitted)
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

// ProcessorEndpointExtensionFeeCollectorUpdated represents a FeeCollectorUpdated event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionFeeCollectorUpdated struct {
	NewFeeCollector common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionFeeCollectorUpdatedEventName = "FeeCollectorUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionFeeCollectorUpdated) ContractEventName() string {
	return ProcessorEndpointExtensionFeeCollectorUpdatedEventName
}

// UnpackFeeCollectorUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FeeCollectorUpdated(address newFeeCollector)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackFeeCollectorUpdatedEvent(log *types.Log) (*ProcessorEndpointExtensionFeeCollectorUpdated, error) {
	event := "FeeCollectorUpdated"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionFeeCollectorUpdated)
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

// ProcessorEndpointExtensionInitialized represents a Initialized event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionInitialized) ContractEventName() string {
	return ProcessorEndpointExtensionInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInitializedEvent(log *types.Log) (*ProcessorEndpointExtensionInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionInitialized)
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

// ProcessorEndpointExtensionMaxNumberOfAppUpdated represents a MaxNumberOfAppUpdated event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionMaxNumberOfAppUpdated struct {
	OldMax *big.Int
	NewMax *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionMaxNumberOfAppUpdatedEventName = "MaxNumberOfAppUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionMaxNumberOfAppUpdated) ContractEventName() string {
	return ProcessorEndpointExtensionMaxNumberOfAppUpdatedEventName
}

// UnpackMaxNumberOfAppUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MaxNumberOfAppUpdated(uint256 oldMax, uint256 newMax)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackMaxNumberOfAppUpdatedEvent(log *types.Log) (*ProcessorEndpointExtensionMaxNumberOfAppUpdated, error) {
	event := "MaxNumberOfAppUpdated"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionMaxNumberOfAppUpdated)
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

// ProcessorEndpointExtensionPaymentWithdrawn represents a PaymentWithdrawn event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionPaymentWithdrawn struct {
	TokenAddress common.Address
	Payee        common.Address
	Amount       *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionPaymentWithdrawnEventName = "PaymentWithdrawn"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionPaymentWithdrawn) ContractEventName() string {
	return ProcessorEndpointExtensionPaymentWithdrawnEventName
}

// UnpackPaymentWithdrawnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentWithdrawn(address tokenAddress, address indexed payee, uint256 amount)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackPaymentWithdrawnEvent(log *types.Log) (*ProcessorEndpointExtensionPaymentWithdrawn, error) {
	event := "PaymentWithdrawn"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionPaymentWithdrawn)
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

// ProcessorEndpointExtensionQueueThresholdUpdated represents a QueueThresholdUpdated event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionQueueThresholdUpdated struct {
	NewThreshold *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionQueueThresholdUpdatedEventName = "QueueThresholdUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionQueueThresholdUpdated) ContractEventName() string {
	return ProcessorEndpointExtensionQueueThresholdUpdatedEventName
}

// UnpackQueueThresholdUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event QueueThresholdUpdated(uint256 newThreshold)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackQueueThresholdUpdatedEvent(log *types.Log) (*ProcessorEndpointExtensionQueueThresholdUpdated, error) {
	event := "QueueThresholdUpdated"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionQueueThresholdUpdated)
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

// ProcessorEndpointExtensionSelectionGraceUpdated represents a SelectionGraceUpdated event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionSelectionGraceUpdated struct {
	NewGrace *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionSelectionGraceUpdatedEventName = "SelectionGraceUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionSelectionGraceUpdated) ContractEventName() string {
	return ProcessorEndpointExtensionSelectionGraceUpdatedEventName
}

// UnpackSelectionGraceUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SelectionGraceUpdated(uint256 newGrace)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackSelectionGraceUpdatedEvent(log *types.Log) (*ProcessorEndpointExtensionSelectionGraceUpdated, error) {
	event := "SelectionGraceUpdated"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionSelectionGraceUpdated)
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

// ProcessorEndpointExtensionTriggerExecuted represents a TriggerExecuted event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionTriggerExecuted struct {
	ApplicationId      uint64
	ProcessedRequestId [32]byte
	Success            bool
	Raw                *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionTriggerExecutedEventName = "TriggerExecuted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionTriggerExecuted) ContractEventName() string {
	return ProcessorEndpointExtensionTriggerExecutedEventName
}

// UnpackTriggerExecutedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TriggerExecuted(uint64 indexed applicationId, bytes32 indexed processedRequestId, bool success)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTriggerExecutedEvent(log *types.Log) (*ProcessorEndpointExtensionTriggerExecuted, error) {
	event := "TriggerExecuted"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionTriggerExecuted)
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

// ProcessorEndpointExtensionTriggerWithdraw represents a TriggerWithdraw event raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionTriggerWithdraw struct {
	ApplicationId       uint64
	ProcessedRequestId  [32]byte
	WithdrawSuccess     bool
	PostWithdrawSuccess bool
	ReturnedTokens      []StructsTokenAndAmount
	FailedTokens        []StructsTokenAndAmount
	Raw                 *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointExtensionTriggerWithdrawEventName = "TriggerWithdraw"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointExtensionTriggerWithdraw) ContractEventName() string {
	return ProcessorEndpointExtensionTriggerWithdrawEventName
}

// UnpackTriggerWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TriggerWithdraw(uint64 indexed applicationId, bytes32 indexed processedRequestId, bool withdrawSuccess, bool postWithdrawSuccess, (address,uint256)[] returnedTokens, (address,uint256)[] failedTokens)
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTriggerWithdrawEvent(log *types.Log) (*ProcessorEndpointExtensionTriggerWithdraw, error) {
	event := "TriggerWithdraw"
	if log.Topics[0] != processorEndpointExtension.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointExtensionTriggerWithdraw)
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
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["AddressCantBeZero"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["DeadlineExpired"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackDeadlineExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["DeployerNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackDeployerNotAllowedError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidInitializationError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidSigner"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["InvalidValue"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackInvalidValueError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["MaxNumOfApplicationsExceeded"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackMaxNumOfApplicationsExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackNotInitializingError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["TokenNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackTokenNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["TransferAmountMismatch"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackTransferAmountMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["TransferFailed"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackTransferFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["TriggerAlreadyRegistered"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackTriggerAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpointExtension.abi.Errors["TriggerCannotBeEOA"].ID.Bytes()[:4]) {
		return processorEndpointExtension.UnpackTriggerCannotBeEOAError(raw[4:])
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

// ProcessorEndpointExtensionAddressCantBeZero represents a AddressCantBeZero error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressCantBeZero()
func ProcessorEndpointExtensionAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x4b054c82bb9e4ebe90744902747c4e86028dd2a122c63de8810d9a1c2e84617d")
}

// UnpackAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressCantBeZero()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackAddressCantBeZeroError(raw []byte) (*ProcessorEndpointExtensionAddressCantBeZero, error) {
	out := new(ProcessorEndpointExtensionAddressCantBeZero)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "AddressCantBeZero", raw); err != nil {
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

// ProcessorEndpointExtensionDeployerNotAllowed represents a DeployerNotAllowed error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionDeployerNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DeployerNotAllowed()
func ProcessorEndpointExtensionDeployerNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xfc336c41b5f2764aa23a0642681230d7b417aea242ed1c10a52ccb7eacd8481f")
}

// UnpackDeployerNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DeployerNotAllowed()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackDeployerNotAllowedError(raw []byte) (*ProcessorEndpointExtensionDeployerNotAllowed, error) {
	out := new(ProcessorEndpointExtensionDeployerNotAllowed)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "DeployerNotAllowed", raw); err != nil {
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

// ProcessorEndpointExtensionInvalidInitialization represents a InvalidInitialization error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func ProcessorEndpointExtensionInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackInvalidInitializationError(raw []byte) (*ProcessorEndpointExtensionInvalidInitialization, error) {
	out := new(ProcessorEndpointExtensionInvalidInitialization)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
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

// ProcessorEndpointExtensionMaxNumOfApplicationsExceeded represents a MaxNumOfApplicationsExceeded error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionMaxNumOfApplicationsExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error MaxNumOfApplicationsExceeded()
func ProcessorEndpointExtensionMaxNumOfApplicationsExceededErrorID() common.Hash {
	return common.HexToHash("0x2d5c0cfe697fe71ea4a46b3137aee7a75e8d99fed238e7efea81025efbed81ae")
}

// UnpackMaxNumOfApplicationsExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error MaxNumOfApplicationsExceeded()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackMaxNumOfApplicationsExceededError(raw []byte) (*ProcessorEndpointExtensionMaxNumOfApplicationsExceeded, error) {
	out := new(ProcessorEndpointExtensionMaxNumOfApplicationsExceeded)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "MaxNumOfApplicationsExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionNotInitializing represents a NotInitializing error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func ProcessorEndpointExtensionNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackNotInitializingError(raw []byte) (*ProcessorEndpointExtensionNotInitializing, error) {
	out := new(ProcessorEndpointExtensionNotInitializing)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
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

// ProcessorEndpointExtensionTransferFailed represents a TransferFailed error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionTransferFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TransferFailed()
func ProcessorEndpointExtensionTransferFailedErrorID() common.Hash {
	return common.HexToHash("0x90b8ec1877afffd816d05d9b13947f3ff18ec5851c38bad15ec2b710f92391b1")
}

// UnpackTransferFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TransferFailed()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTransferFailedError(raw []byte) (*ProcessorEndpointExtensionTransferFailed, error) {
	out := new(ProcessorEndpointExtensionTransferFailed)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "TransferFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionTriggerAlreadyRegistered represents a TriggerAlreadyRegistered error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionTriggerAlreadyRegistered struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TriggerAlreadyRegistered()
func ProcessorEndpointExtensionTriggerAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0xd48c37974b982f54734c76a70a50defa3e9afc19e571e138c59cb6f81e9a39ae")
}

// UnpackTriggerAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TriggerAlreadyRegistered()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTriggerAlreadyRegisteredError(raw []byte) (*ProcessorEndpointExtensionTriggerAlreadyRegistered, error) {
	out := new(ProcessorEndpointExtensionTriggerAlreadyRegistered)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "TriggerAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointExtensionTriggerCannotBeEOA represents a TriggerCannotBeEOA error raised by the ProcessorEndpointExtension contract.
type ProcessorEndpointExtensionTriggerCannotBeEOA struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TriggerCannotBeEOA()
func ProcessorEndpointExtensionTriggerCannotBeEOAErrorID() common.Hash {
	return common.HexToHash("0x4ea29e672521f87a56ee460cd2be1b29334f535a27e686ee3144db2a8882eab3")
}

// UnpackTriggerCannotBeEOAError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TriggerCannotBeEOA()
func (processorEndpointExtension *ProcessorEndpointExtension) UnpackTriggerCannotBeEOAError(raw []byte) (*ProcessorEndpointExtensionTriggerCannotBeEOA, error) {
	out := new(ProcessorEndpointExtensionTriggerCannotBeEOA)
	if err := processorEndpointExtension.abi.UnpackIntoInterface(out, "TriggerCannotBeEOA", raw); err != nil {
		return nil, err
	}
	return out, nil
}
