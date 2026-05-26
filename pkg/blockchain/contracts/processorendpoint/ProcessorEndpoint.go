// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package processorendpoint

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

// StructsPendingRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsPendingRequest struct {
	Timestamp       *big.Int
	TokenAddress    common.Address
	AssetAmount     *big.Int
	MaxFeeValue     *big.Int
	RequestId       [32]byte
	Payload         []byte
	Sender          common.Address
	Facilitator     common.Address
	ApplicationId   uint64
	ProtocolVersion uint8
	RequestType     uint8
}

// StructsWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsWithdrawalRequest struct {
	TokenAddress common.Address
	Receiver     common.Address
	Amount       *big.Int
}

// ProcessorEndpointMetaData contains all meta data concerning the ProcessorEndpoint contract.
var ProcessorEndpointMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resetOperator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"},{\"internalType\":\"contractITokenAllowlist\",\"name\":\"_tokenAllowlist\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeployerNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientAppBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPermit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxNumOfApplicationsExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotRegisteredTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"AppEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"DeployRequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"DeployRequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMax\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"MaxNumberOfAppUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"ReportGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"triggerContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"triggerContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TriggerPostWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPLOYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_OPERATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"addAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminReset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"appIds\",\"type\":\"uint64[]\"}],\"name\":\"adminResetApps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"appCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"applicationStateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableDeploySlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"facilitatorNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeployedAppIds\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getFacilitatorNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getPendingRequestsPage\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTriggerQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"isAllowedDeployer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumOfApplications\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"removeAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitDeployRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"requestSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"depositPermit\",\"type\":\"bytes\"}],\"name\":\"submitRequestFor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"submitTriggerRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalAppCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalPendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"triggerContracts\",\"outputs\":[{\"internalType\":\"contractITrigger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"triggersToAppIds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"updateMaxNumOfApplications\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "ProcessorEndpoint",
	Bin: "0x610160806040523461034b5760e0816164708038038091610020828561034f565b83398101031261034b5780516001600160a01b0381169081900361034b5760208201516001600160a01b038116929083900361034b5761006260408201610386565b9261006f60608301610386565b61007b60808401610386565b9360c060a08501519401519560018060a01b03871680970361034b5760409687516100a6898261034f565b60048152602081016356656c6160e01b81525f906100c4600161039a565b926100d18c51948561034f565b600184526100df600161039a565b602085019390601f1901368537600a602186015b5f1901916f181899199a1a9b1b9c1cb0b131b232b360811b8282061a835304801561012157600a90916100f3565b505060018055610130816105e0565b6101205261013d8461077b565b61014052519020918260e05251902080610100524660a05289519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f84528b83015260608201524660808201523060a082015260a081526101a460c08261034f565b5190206080523060c052600a600655600a600755600a600c5582158015610343575b8015610332575b8015610321575b8015610319575b61030a57600d80546001600160a01b03199081169094179055600e80548416909517909455600f8054831694909417909355601580549091166001600160a01b0384161790556102979161022e906103b5565b506102388161044d565b505f5160206164105f395f51905f525f818152602081905286812060010180545f5160206164305f395f51905f5291829055909290917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9080a46104cd565b506014556001600160a01b0381166102fa575b5051615adc90816108b48239608051816157b7015260a0518161586e015260c05181615781015260e051816158060152610100518161582c015261012051816112d8015261014051816113010152f35b6103039061054d565b505f6102aa565b632582a64160e11b5f5260045ffd5b5080156101db565b506001600160a01b038416156101d4565b506001600160a01b038216156101cd565b5084156101c6565b5f80fd5b601f909101601f19168101906001600160401b0382119082101761037257604052565b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b038216820361034b57565b6001600160401b03811161037257601f01601f191660200190565b6001600160a01b0381165f9081525f5160206163b05f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f5160206163b05f395f51905f5260205260408120805460ff191660011790553391907f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd98905f5160206163905f395f51905f529080a4600190565b505f90565b6001600160a01b0381165f9081525f5160206163d05f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f5160206163d05f395f51905f5260205260408120805460ff191660011790553391905f5160206164305f395f51905f52905f5160206163905f395f51905f529080a4600190565b6001600160a01b0381165f9081525f5160206164505f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f5160206164505f395f51905f5260205260408120805460ff191660011790553391905f5160206164105f395f51905f52905f5160206163905f395f51905f529080a4600190565b6001600160a01b0381165f9081525f5160206163f05f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f5160206163f05f395f51905f5260205260408120805460ff191660011790553391907fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c95917905f5160206163905f395f51905f529080a4600190565b908151602081105f1461065a575090601f81511161061a57602081519101516020821061060b571790565b5f198260200360031b1b161790565b604460209160405192839163305a27a960e01b83528160048401528051918291826024860152018484015e5f828201840152601f01601f19168101030190fd5b6001600160401b03811161037257600254600181811c91168015610771575b602082101461075d57601f811161072a575b50602092601f82116001146106c957928192935f926106be575b50508160011b915f199060031b1c19161760025560ff90565b015190505f806106a5565b601f1982169360025f52805f20915f5b86811061071257508360019596106106fa575b505050811b0160025560ff90565b01515f1960f88460031b161c191690555f80806106ec565b919260206001819286850151815501940192016106d9565b60025f52601f60205f20910160051c810190601f830160051c015b818110610752575061068b565b5f8155600101610745565b634e487b7160e01b5f52602260045260245ffd5b90607f1690610679565b908151602081105f146107a6575090601f81511161061a57602081519101516020821061060b571790565b6001600160401b03811161037257600354600181811c911680156108a9575b602082101461075d57601f8111610876575b50602092601f821160011461081557928192935f9261080a575b50508160011b915f199060031b1c19161760035560ff90565b015190505f806107f1565b601f1982169360035f52805f20915f5b86811061085e5750836001959610610846575b505050811b0160035560ff90565b01515f1960f88460031b161c191690555f8080610838565b91926020600181928685015181550194019201610825565b60035f52601f60205f20910160051c810190601f830160051c015b81811061089e57506107d7565b5f8155600101610891565b90607f16906107c556fe61044080604052600436101561001d575b50361561001b575f80fd5b005b5f905f3560e01c90816301ffc9a71461249b575080631256c918146124635780631664deeb146124295780631753520f1461240c57806318733d57146123cb578063194b6be8146123935780631cc76c381461235357806320ff473f1461230b57806321c0b342146122d9578063248a9ca3146122af5780632a0acc6a146122755780632b52b07c1461223b5780632f2ff15d146121fe5780632fbfa0d514611ec057806336568abe14611e7c57806347101cae14611e625780634a0d049514611aa85780635423112f14611a6157806355c0be5c146116a65780635e3393841461168c578063655f2de1146116115780636b288d20146115c35780636e91d1331461159b57806374202d9b146115275780637596ed43146114c25780637a36a8911461148a57806380a1f71214611408578063840059ec146113b857806384b0196e146112c05780638537c2781461113457806386f4a07d146110fa5780638c5b9b00146110cc5780638eb8791a146110af57806391d14854146110675780639c73eeaf146110155780639c8d416f14610ff85780639fabc36f14610fd0578063a217fddf14610f8e578063a9ba045e14610fa8578063aa3aa46014610f8e578063af20c96014610f56578063b7222ff514610f03578063bcf7a3c9146108a2578063be1e37251461086f578063c404c3591461082e578063c415b95c14610805578063d2c35ce814610792578063d547741f1461074f578063d72e957a14610716578063dd39d76514610348578063e3d72e5e146102fb578063ecd00261146102c0578063f624f88a146102975763f7d9cc1e0361001057346102945780600319360112610294576020600c54604051908152f35b80fd5b5034610294578060031936011261029457600f546040516001600160a01b039091168152602090f35b503461029457806003193601126102945760206040517ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c8152f35b5034610294576020366003190112610294576103156124ee565b61031d61419b565b6001600160a01b038116156103395761033590614573565b5080f35b632582a64160e11b8252600482fd5b5034610294576020366003190112610294576004356001600160401b0381116107125736602382011215610712578060040135906001600160401b03821161070e5760248260051b8201019036821161070a576103a361420a565b6103ab6143e4565b6103b3615345565b6103bc83612c62565b916103ca6040519384612742565b8383526024602084019201915b8183106106e65750505090156106d8575b600f5460405163024ece8960e01b8152908390829060049082906001600160a01b03165afa9081156106cd5783916106ab575b50815191815190845b84810361043357856001805580f35b6001600160401b036104458284612ebc565b51168087526012602052604087205f805260205260405f20548061064f575b50865b84810361059157508087526004602052604087205461048a575b50600101610424565b8087526004602052866040812055600554875b8181036104b9575b505050600190816007540160075590610481565b826001600160401b036104cb83612f21565b90549060031b1c16146104e05760010161049d565b91505f19810190811161057d57906105136001600160401b0361050561053294612f21565b90549060031b1c1691612f21565b9091906001600160401b038084549260031b9316831b921b1916179055565b600554801561056957600191905f190161054b81612f21565b6001600160401b0382549160031b1b19169055600555905f806104a5565b634e487b7160e01b87526031600452602487fd5b634e487b7160e01b88526011600452602488fd5b818852601260205260408820600191906001600160a01b036105b3838a612ebc565b5116838060a01b03165f5260205260405f2054806105d3575b5001610467565b6105ef8133858060a01b036105e8868d612ebc565b5116614404565b828060a01b036105ff838a612ebc565b51168a52601360205261061760408b209182546126fe565b9055828952601260205260408920828060a01b03610635838a612ebc565b5116838060a01b03165f526020528860405f20555f6105cc565b8780808084335af161065f61277e565b501561069c57878052601360205261067c604089209182546126fe565b90558087526012602052604087205f80526020528660405f20555f610464565b6312171d8360e31b8852600488fd5b6106c791503d8085833e6106bf8183612742565b81019061411a565b5f61041b565b6040513d85823e3d90fd5b506106e1612b5d565b6103e8565b82356001600160401b0381168103610706578152602092830192016103d7565b8680fd5b8380fd5b8280fd5b5080fd5b5034610294576020366003190112610294576020906040906001600160a01b0361073e6124ee565b168152601683522054604051908152f35b50346102945760403660031901126102945761033560043561076f612504565b9061078d610788825f525f602052600160405f20015490565b614279565b614617565b5034610294576020366003190112610294576107ac6124ee565b6107b461419b565b6001600160a01b03168015610339576020817fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f926001600160601b0360a01b6015541617601555604051908152a180f35b50346102945780600319360112610294576015546040516001600160a01b039091168152602090f35b503461029457806003193601126102945761085e61084a61404d565b6040519384936060855260608501906125f1565b916020840152151560408301520390f35b50346102945760403660031901126102945761089e610892602435600435615480565b6040519182918261269f565b0390f35b50610140366003190112610d85576108b86124ee565b6024359160ff83168303610d85576044356001600160401b038116808203610d855760046064351015610d85576084356001600160401b038111610d8557610904903690600401612596565b92909461090f612530565b91610104356001600160401b038111610d8557610930903690600401612596565b939096610124356001600160401b038111610d8557610953903690600401612596565b60ff8c16610ef457875f52600460205260405f205415610ee5576109756143e4565b6109806064356125e7565b6003606435141580610ecd575b610ebe5760e4354211610eaf576109a2614697565b600c541115610ea0578a8c976109b96064356125e7565b600360643514610e6c575b6001600160a01b0388165f81815260166020526040902054909c90928d906109ed6064356125e7565b8d36906109f9926128b2565b805190602001206040519b60208d01928d7f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc133334085526040015260ff1660608d01528c60808d015260643560ff1660a08d015260c08c0152600160a01b6001900388169a8b60e082015260c4356101008201528461012082015260e4356101408201526101408152610a8b61016082612742565b519020610a9661577e565b906040519161190160f01b83526002830152602282015260429020913690610abd926128b2565b610ac691615894565b610acf916158ce565b6001600160a01b03168015610e5d578b03610e4e5760018101809111610e3a578a5f52601660205260405f20556014543410610e2b5760c435158080610e22575b610e135715610bd0575b5050610b43610bae959460ff9460209c94610b7594508a8d600b549460c4359460643591612ade565b98899760405195610b538761270b565b4287528c87015260c435604087015234606087015288608087015236916128b2565b60a08401528660c08401523360e08401528461010084015216610120820152610b9f6064356125e7565b606435610140820152836146f3565b5f516020615a875f395f51905f5285604051338152a460018055604051908152f35b8615610e1357600f5460405163cbe230c360e01b81526004810189905290602090829060249082906001600160a01b03165afa908115610d7a575f91610dd9575b5015610dca5760608103610dbb578160609181010312610d855780359060ff8216809203610d8557604051636eb1769f60e11b8152600481018b90523060248201526020816044818b5afa908115610d7a575f91610d89575b5060c43511610ce1575b5050610b43610bae959460ff9460209c94610b7594610c9660c4358584614445565b8a815260128f5260408120895f528f5260405f20610cb760c43582546128a5565b905588815260138f5260408120610cd160c43582546128a5565b90559450949c5094509495610b1a565b86959192949b93963b15610d855760409081519363d505accf60e01b85528b600486015230602486015260c435604486015260e43560648601526084850152602081013560a4850152013560c48301525f8260e48183895af1908115610d7a5760209b610bae9760ff96610b7595610b4395610d67575b509450949c5094509495610c74565b610d7391505f90612742565b5f5f610d58565b6040513d5f823e3d90fd5b5f80fd5b90506020813d602011610db3575b81610da460209383612742565b81010312610d8557515f610c6a565b3d9150610d97565b63ddafbaef60e01b5f5260045ffd5b63514e24c360e11b5f5260045ffd5b90506020813d602011610e0b575b81610df460209383612742565b81010312610d8557610e0590612898565b5f610c11565b3d9150610de7565b632a9ffab760e21b5f5260045ffd5b50871515610b10565b63c77f97eb60e01b5f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b632057875960e21b5f5260045ffd5b638baa579f60e01b5f5260045ffd5b975050608589141580610e95575b610e86578a8c976109c4565b637c6953f960e01b5f5260045ffd5b5060e2891415610e7a565b63d8219b3d60e01b5f5260045ffd5b631ab7da6b60e01b5f5260045ffd5b630ee2db0560e21b5f5260045ffd5b50610ed96064356125e7565b6001606435141561098d565b6309ff9a8360e41b5f5260045ffd5b635428eae760e01b5f5260045ffd5b34610d85576040366003190112610d8557610f1c61255a565b6001600160401b03610f2c612504565b91165f52601260205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b34610d85576020366003190112610d85576001600160a01b03610f776124ee565b165f526013602052602060405f2054604051908152f35b34610d85575f366003190112610d855760206040515f8152f35b34610d85575f366003190112610d8557600e546040516001600160a01b039091168152602090f35b34610d85576020366003190112610d85576020610fee60043561400d565b6040519015158152f35b34610d85575f366003190112610d85576020601454604051908152f35b34610d85576020366003190112610d855760043561103161419b565b8015610e13576020817e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db192600c55604051908152a1005b34610d85576040366003190112610d8557611080612504565b6004355f525f60205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b34610d85575f366003190112610d85576020600654604051908152f35b34610d85575f366003190112610d85576110e461420a565b6110ec6143e4565b6110f4615345565b60018055005b34610d85575f366003190112610d855760206040517fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c959178152f35b34610d8557610180366003190112610d855761114e61255a565b608435906001600160401b038211610d855760406003198336030112610d855760a4356001600160401b038111610d855760406003198236030112610d855760c435916001600160401b038311610d855736602384011215610d85578260040135926001600160401b038411610d85573660246060860283010111610d855761012435600d811015610d8557610144356001600160401b038111610d85576111fa903690600401612596565b939092610164356001600160401b038111610d855761121d903690600401612596565b335f9081527e9b85fa08e135b4afefc8a40566229b01e66445186a3676642ffbacd3d93b0e6020526040902054909891979060ff1615611289576110f4996112636143e4565b6101043594602460e4359501926004019160040190606435906044359060243590612f61565b63e2517d3f60e01b5f52336004527f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9860245260445ffd5b34610d85575f366003190112610d855761135c6112fc7f00000000000000000000000000000000000000000000000000000000000000006155cc565b6113257f00000000000000000000000000000000000000000000000000000000000000006156c6565b602061136a604051926113388385612742565b5f84525f368137604051958695600f60f81b875260e08588015260e08701906125c3565b9085820360408701526125c3565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b8281106113a157505050500390f35b835185528695509381019392810192600101611392565b34610d85576040366003190112610d85576113d16124ee565b6113d9612504565b6001600160a01b039182165f908152601060209081526040808320949093168252928352819020549051908152f35b34610d85575f366003190112610d8557611428611423614697565b614910565b600a54600b54905f905b828110611447576040518061089e868261269f565b60018181925f52600960205260405f20545f52600860205261146b60405f20612a10565b6114758588612ebc565b526114808487612ebc565b5001910190611432565b34610d85576020366003190112610d85576001600160401b036114ab61255a565b165f526004602052602060405f2054604051908152f35b34610d85575f366003190112610d85576114da612b5d565b6040518091602082016020835281518091526020604084019201905f5b818110611505575050500390f35b82516001600160401b03168452859450602093840193909201916001016114f7565b34610d855760e0366003190112610d85576115406124ee565b611548612570565b90604435916004831015610d8557606435916001600160401b038311610d855760209361157c611593943690600401612596565b9061158561251a565b9260c4359560a43595612ade565b604051908152f35b34610d85575f366003190112610d8557600d546040516001600160a01b039091168152602090f35b34610d85576020366003190112610d85576115dc6124ee565b6001600160a01b03165f9081525f516020615a675f395f51905f52602090815260409182902054915160ff9092161515825290f35b34610d85576020366003190112610d855760043561162d61419b565b8015610e1357600654611642600754826126fe565b90818310610e13578261167c6040937f80af1d8add7b822532307ae28ef73c31ab2803f908c5c22cc879766f54c6858395600655826126fe565b60075582519182526020820152a1005b34610d85575f366003190112610d85576020611593614697565b34610d855760e0366003190112610d85576004356004811015610d85576024356001600160401b038111610d85576116e2903690600401612596565b906044356001600160a01b03811690818103610d855760643592611704612530565b60c435949091906001600160a01b0386168603610d85576117236143e4565b335f5260186020526001600160401b0360405f205416938415611a52576117539082601c54918a868d8a8a612ade565b96604051926117618461270b565b428452602084019687526040840192835261178e6060850191608435835260808601938b855236916128b2565b60a085019081526001600160a01b0386811660c08701908152908a1660e087015261010086018890525f6101208701529a9093906117cb816125e7565b6101408601525f8a81526019602052604090208551815597516001890180546001600160a01b0319166001600160a01b0392909216919091179055516002880155516003870155516004860155518051979094906001600160401b038911611a3e57602098889661183f600585015461293a565b601f8111611a02575b508a90601f8311600114611980579261188d836007945f516020615a875f395f51905f529997948f99975f92611975575b50508160011b915f199060031b1c19161790565b60058301555b516006820180546001600160a01b03199081166001600160a01b039384161790915560e08086015194909301805490911693909116929092178083556101008401516101208501516101409095015193949093921b60ff60e01b16916118f8816125e7565b611901816125e7565b69ffffffffffffffffffff60a01b1990911660a09390931b67ffffffffffffffff60a01b16929092171760e89190911b60ff60e81b16179055601c80545f908152601a84526040908190208790558154600101909155516001600160a01b039687168152951694a460018055604051908152f35b015190508f80611879565b90600585015f52805f20915f5b601f19851681106119e85750835f516020615a875f395f51905f529896938e989693600193600797601f198116106119d0575b505050811b016005830155611893565b01515f1960f88460031b161c191690558e80806119c0565b8183015184558c9a50600190930192918d01918d0161198d565b611a2e90600586015f528c5f20601f850160051c8101918e8610611a34575b601f0160051c01906146dd565b8b611848565b9091508190611a21565b634e487b7160e01b5f52604160045260245ffd5b63578f584160e01b5f5260045ffd5b34610d85576020366003190112610d8557611a7a6128e8565b506004355f52600860205261089e611a9460405f20612a10565b6040519182916020835260208301906125f1565b6040366003190112610d8557611abc612586565b602435906001600160401b038211610d8557611ade60ff923690600401612596565b9290911680610ef457611aef6143e4565b335f9081525f516020615a675f395f51905f52602052604090205460ff1615611e53576007548015611e4457611b23614697565b600c541115610ea0576014543410610e2b575f1901600755600b546040805133602082019081525f92820192909252909181611b95915f606083015260e06080830152611b7561010083018989612abe565b905f60a08401525f60c084015260e083015203601f198101835282612742565b519020918260c01c91604051611baa8161270b565b42815260208101905f825260408101905f8252611bd760608201953487526080830199898b5236916128b2565b9760a0820198895260c082019033825260e08301965f885261010084019489865261012085019788526101408501965f88528b5f52600860205260405f209551865560018060a01b03905116600186019060018060a01b03166001600160601b0360a01b8254161790555160028501555160038401555160048301556005820198519889516001600160401b038111611a3e5760209a611c77835461293a565b601f8111611e02575b508b90601f8311600114611d9b579180611cb4926007979695945f92611d905750508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b1691611d06816125e7565b611d0f816125e7565b60ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600b545f52600983528160405f20556001600b5401600b55604051908282527fbbcfb0c180d5f95cfb50c12005c44144ac92ac8474f1abc5add0e6e69a14e8dc843393a360018055604051908152f35b015190508e80611879565b90601f19831691845f52815f20925f5b818110611deb57509160019391856007999897969410611dd3575b505050811b019055611cb7565b01515f1960f88460031b161c191690558d8080611dc6565b92938f600181928786015181550195019301611dab565b611e2a90845f528d5f20601f850160051c8101918f8610611a3457601f0160051c01906146dd565b8c611c80565b634e487b7160e01b5f52602160045260245ffd5b6316ae067f60e11b5f5260045ffd5b63fc336c4160e01b5f5260045ffd5b34610d85575f366003190112610d855760206115936146b3565b34610d85576040366003190112610d8557611e95612504565b336001600160a01b03821603611eb15761001b90600435614617565b63334bd91960e11b5f5260045ffd5b60e0366003190112610d8557611ed4612586565b611edc612570565b6044356004811015610d85576064356001600160401b038111610d8557611f07903690600401612596565b611f1292919261251a565b9160a4359460ff60c43597169081610ef4576001600160401b03811695865f52600460205260405f205415610ee557611f496143e4565b611f52846125e7565b8315610ebe576014548910610e2b57611f69614697565b600c541115610ea0576001600160a01b038616948561216657611f8c8a8a6128a5565b3403610e13575b611f9c856125e7565b600385036120995760858114158061208e575b610e8657602099612009612035948b6120689a5b8c5f5260128f5260405f208b5f528f5260405f20611fe28382546128a5565b90558a5f5260138f5260405f20611ffa8382546128a5565b905585878b600b549533612ade565b99604051976120178961270b565b4289528c8901526040880152606087015288608087015236916128b2565b60a08401523360c08401525f60e08401528461010084015261012083015261205c816125e7565b610140820152836146f3565b81604051915f83525f516020615a875f395f51905f52853394a460018055604051908152f35b5060e2811415611faf565b9190929394956120a8866125e7565b600286146120ca575b9160209961200961206898979695938b61203596611fc3565b9594939290916044602060018060a01b03600e5416604051928380926308b56a5360e11b82528d60048301523360248301525afa908115610d7a575f9161212c575b501561211d579190929394956120b1565b63622a850b60e11b5f5260045ffd5b90506020813d60201161215e575b8161214760209383612742565b81010312610d855761215890612898565b8b61210c565b3d915061213a565b893403610e13578815610e1357600f5460405163cbe230c360e01b81526004810188905290602090829060249082906001600160a01b03165afa908115610d7a575f916121c4575b5015610dca576121bf893389614445565b611f93565b90506020813d6020116121f6575b816121df60209383612742565b81010312610d85576121f090612898565b8b6121ae565b3d91506121d2565b34610d85576040366003190112610d855761001b60043561221d612504565b90612236610788825f525f602052600160405f20015490565b61435c565b34610d85575f366003190112610d855760206040517f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc13333408152f35b34610d85575f366003190112610d855760206040517fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec428152f35b34610d85576020366003190112610d855760206115936004355f525f602052600160405f20015490565b34610d85576040366003190112610d85576110f46122f56124ee565b6122fd612504565b906123066143e4565b6127ad565b34610d85576020366003190112610d85576123246124ee565b61232c61419b565b6001600160a01b038116156123445761001b906142b1565b632582a64160e11b5f5260045ffd5b34610d85576020366003190112610d85576001600160401b0361237461255a565b165f526017602052602060018060a01b0360405f205416604051908152f35b34610d85576020366003190112610d85576001600160a01b036123b46124ee565b165f526011602052602060405f2054604051908152f35b34610d85576020366003190112610d85576001600160a01b036123ec6124ee565b165f52601860205260206001600160401b0360405f205416604051908152f35b34610d85575f366003190112610d85576020600754604051908152f35b34610d85575f366003190112610d855760206040517f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd988152f35b34610d85576020366003190112610d85576001600160a01b036124846124ee565b165f526016602052602060405f2054604051908152f35b34610d85576020366003190112610d85576004359063ffffffff60e01b8216809203610d8557602091637965db0b60e01b81149081156124dd575b5015158152f35b6301ffc9a760e01b149050836124d6565b600435906001600160a01b0382168203610d8557565b602435906001600160a01b0382168203610d8557565b608435906001600160a01b0382168203610d8557565b60a435906001600160a01b0382168203610d8557565b35906001600160a01b0382168203610d8557565b600435906001600160401b0382168203610d8557565b602435906001600160401b0382168203610d8557565b6004359060ff82168203610d8557565b9181601f84011215610d85578235916001600160401b038311610d855760208381860195010111610d8557565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60041115611e3057565b908151815260018060a01b0360208301511660208201526040820151604082015260608201516060820152608082015160808201526101408061264560a085015161016060a08601526101608501906125c3565b9360018060a01b0360c08201511660c085015260018060a01b0360e08201511660e08501526001600160401b036101008201511661010085015260ff6101208201511661012085015201519161269a836125e7565b015290565b602081016020825282518091526040820191602060408360051b8301019401925f915b8383106126d157505050505090565b90919293946020806126ef600193603f1986820301875289516125f1565b970193019301919392906126c2565b91908203918211610e3a57565b61016081019081106001600160401b03821117611a3e57604052565b604081019081106001600160401b03821117611a3e57604052565b90601f801991011681019081106001600160401b03821117611a3e57604052565b6001600160401b038111611a3e57601f01601f191660200190565b3d156127a8573d9061278f82612763565b9161279d6040519384612742565b82523d5f602084013e565b606090565b60018060a01b03811691825f52601060205260405f2060018060a01b0382165f5260205260405f2054918215612892577f16d46b28f1a9377c2df78652a0a8e31568b99311665bcec24cbba72edd2e5c2c8392855f52601060205260405f2060018060a01b0382165f526020525f6040812055855f52601160205260405f206128378582546126fe565b9055604080516001600160a01b0394851681526020810195909552921692839290a282612887575f809350809281925af161287061277e565b501561287857565b6312171d8360e31b5f5260045ffd5b61289092614404565b565b50505050565b51908115158203610d8557565b91908201809211610e3a57565b9291926128be82612763565b916128cc6040519384612742565b829481845281830111610d85578281602093845f960137010152565b604051906128f58261270b565b5f61014083828152826020820152826040820152826060820152826080820152606060a08201528260c08201528260e082015282610100820152826101208201520152565b90600182811c92168015612968575b602083101461295457565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612949565b9060405191825f8254926129858461293a565b80845293600181169081156129ee57506001146129aa575b5061289092500383612742565b90505f9291925260205f20905f915b8183106129d2575050906020612890928201015f61299d565b60209193508060019154838589010152019101909184926129b9565b90506020925061289094915060ff191682840152151560051b8201015f61299d565b90604051612a1d8161270b565b61014060ff600783958054855260018060a01b036001820154166020860152600281015460408601526003810154606086015260048101546080860152612a6660058201612972565b60a086015260018060a01b0360068201541660c0860152015460018060a01b03811660e08501526001600160401b038160a01c16610100850152818160e01c1661012085015260e81c1691612aba836125e7565b0152565b908060209392818452848401375f828201840152601f01601f1916010190565b9691612b5795936001600160401b03979295612b2f92604051998a9860208a019c60018060a01b03168d52166040890152612b18816125e7565b606088015260e06080880152610100870191612abe565b6001600160a01b0390931660a085015260c084015260e083015203601f198101835282612742565b51902090565b60405190600554808352826020810160055f5260205f20925f905b806003830110612c1057612890945491818110612bf6575b818110612bd9575b818110612bbc575b10612bae575b500383612742565b60c01c81526020015f612ba6565b9260206001916001600160401b038560801c168152019301612ba0565b9260206001916001600160401b038560401c168152019301612b98565b9260206001916001600160401b0385168152019301612b90565b916004919350608060019186546001600160401b03811682526001600160401b038160401c1660208301526001600160401b0381841c16604083015260c01c6060820152019401920185929391612b78565b6001600160401b038111611a3e5760051b60200190565b903590601e1981360301821215610d8557018035906001600160401b038211610d8557602001918160051b36038313610d8557565b9190604083820312610d8557604051612cc681612727565b809380356001600160401b038111610d8557810183601f82011215610d85578035612cf081612c62565b91612cfe6040519384612742565b81835260208084019260051b82010190868211610d855760208101925b828410612d9d575050505082526020810135906001600160401b038211610d8557019180601f84011215610d85578235612d5481612c62565b93612d626040519586612742565b81855260208086019260051b820101928311610d8557602001905b828210612d8d5750505060200152565b8135815260209182019101612d7d565b83356001600160401b038111610d8557820188603f82011215610d8557602091612dd08a836040868096013591016128b2565b815201930192612d1b565b60408201908051916040845282518091526060840190602060608260051b8701019401915f905b828210612e4c575050505060200151916020818303910152602080835192838152019201905f5b818110612e365750505090565b8251845260209384019390920191600101612e29565b90919294602080612e69600193605f198b820301865289516125c3565b970192019201909291612e02565b90600d821015611e305752565b9190811015612e94576060020190565b634e487b7160e01b5f52603260045260245ffd5b356001600160a01b0381168103610d855790565b8051821015612e945760209160051b010190565b9190811015612e945760051b0190565b9190811015612e945760051b81013590601e1981360301821215610d855701908135916001600160401b038311610d85576020018236038113610d85579190565b90600554821015612e945760055f52600282901c7f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0019160031b60181690565b9e9d9996939e9c9a9794919b989592909c6101c0526101e0526080526101205260c05260e0526101a052610100526101405261016052612fa36101e05161400d565b15613ffe57612fb06146b3565b151580613fec575b5f6101805215613fd4576101e0515f52601960205260405f20610180525b6007610180510154916001600160401b038360a01c166001600160401b036101c0511603610ee5576001600160401b036101c051165f52600460205260405f20548503613b6b576130278780612c79565b9190506130376020890189612c79565b90508203610e865761304d608051608051612c79565b939050613061602060805101608051612c79565b90508403610e8657604051916130768361270b565b6001600160401b036101c0511683528760208401528660408401526101e05160608401526130a4368b612cae565b60808401526130b536608051612cae565b60a08401526130c560c051612c62565b6130d26040519182612742565b60c05181526020810136606060c05102610120510111610d855761012051905b606060c0510261012051018210613f7757505060c084015260e05160e08401526101a051610100840152600d610100511015611e30579091610100516101208201526131453661016051610140516128b2565b61014082015260018060a01b03600d541691604051938492635c9626b960e01b8452604060048501526001600160401b0381511660448501526020810151606485015260408101516084850152606081015160a48501526131d06131ba608083015161016060c48801526101a4870190612ddb565b60a08301518682036043190160e4880152612ddb565b9060c08101519160431986820301610104870152602080845192838152019301905f5b818110613f375750505060209593859361324e85946101408560e061326097015161012489015261010081015161014489015261323a6101208201516101648a0190612e77565b0151868203604319016101848801526125c3565b84810360031901602486015291612abe565b03915afa908115610d7a575f91613efd575b5015610e5d576101805160038101546006909101546001600160a01b0390811694168015613ef657935b61010051613b7a57506001600160401b036101c051165f5260046020528460405f205414613b6b576132d36101a05160e0516128a5565b03610e13576014546101a05110610e13575f5f906132f260c051612c62565b6132ff6040519182612742565b60c0518152601f1961331260c051612c62565b0136602083013761332460c051612c62565b906133326040519283612742565b60c0518252601f1961334560c051612c62565b013660208401375f925b60c05181106139e957505f5b83811061392f575050505061337e906133796101a05160e0516128a5565b6128a5565b613386614987565b10613920575f5b8181106138ab5750505f5b81811061382257505060ff600761018051015460e81c1660a0526133bd60a0516125e7565b600260a051146137e8575b6001600160401b036101c051165f5260046020528160405f20556133ed60a0516125e7565b60a051156136ed575b6040518381528260208201526101e051907f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e60406001600160401b036101c0511692a360e05161369c575b505f5b60c051811061359b575061348692936101605192610140519261010051926101a0519260e0519260c051926101205192608051926101e051916101c051614cb8565b604051613494602082612742565b5f815261349f6146b3565b151580613589575b1561357c576134b4615987565b6134bf60a0516125e7565b60a05161352c57604051907f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d6101e051928061350e6001600160401b036101c05116945f806101a05185614a98565b0390a35b6101a05160155461289091906001600160a01b03166149ef565b604051907f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed666101e05192806135746001600160401b036101c05116945f806101a05185614a98565b0390a3613512565b61358461594a565b6134b4565b506135966101e05161495f565b6134a7565b806135f76135b96135b460019460c05161012051612e84565b612ea8565b6135d460206135ce8560c05161012051612e84565b01612ea8565b9060406135e78560c05161012051612e84565b013591858060a01b031690614a4e565b61360c60206135ce8360c05161012051612e84565b61361f6135b48360c05161012051612e84565b60406136318460c05161012051612e84565b0135907faafeb614695d37710d00f37c69671d205eb5b6fcd02995f6c3c7e2c93d3cbb2460405193868060a01b031693806136936101e051956001600160401b036101c05116958360209093929193604081019460018060a01b031681520152565b0390a401613444565b60018060a01b03166136b060e051826149ef565b604051604081015f825260e05160208301525f516020615a475f395f51905f526101e05192806001600160401b036101c05116930390a45f613441565b60055468010000000000000000811015611a3e576137148160016137359301600555612f21565b6101c051906001600160401b038084549260031b9316831b921b1916179055565b600561018051016020613748825461293a565b14613754575b506133f6565b61375d90612972565b602081805181010312610d8557602001516001600160a01b0381168103610d85576001600160a01b0381161561374e576101c0516001600160401b03165f81815260176020908152604080832080546001600160a01b0319166001600160a01b0396909616958617905593825260189052918220805467ffffffffffffffff1916909117905561374e565b6101e0516001600160401b036101c051167ff36f013d86705b333febcccf7ece97eee30f4cd398a9c69aee1f940ed5daddea5f80a36133c8565b8061384060019261383a602060805101608051612c79565b90612ed0565b3561385982613853608051608051612c79565b90612ee0565b7f485d4bfb37695319d1b2352eef5982d10db8c284b0a7ddefe4465377c9fbc8df60405160208152806138a26101e051956001600160401b036101c05116956020840191612abe565b0390a401613398565b806138bf60019261383a60208c018c612c79565b356138ce826138538c80612c79565b7f9ede1ad14e46a641d7b174bb6c9d939ad961d1a7de4b669163b9714a67e2232060405160208152806139176101e051956001600160401b036101c05116956020840191612abe565b0390a40161338d565b631e9acf1760e31b5f5260045ffd5b6001600160a01b036139418284612ebc565b516040516370a0823160e01b81523060048201529116602082602481845afa918215610d7a575f926139b4575b50613996816139a7925f52601360205260405f2054905f52601160205260405f2054906128a5565b6139a08487612ebc565b51906128a5565b116139205760010161335b565b9091506020813d82116139e1575b816139cf60209383612742565b81010312610d8557519061399661396e565b3d91506139c2565b936139fd6135b48660c05161012051612e84565b906040613a108760c05161012051612e84565b0135916001600160401b036101c051165f52601260205260405f2060018060a01b0382165f5260205260405f20548311613b5c576101c0516001600160401b03165f9081526012602090815260408083206001600160a01b038516845282528083208054879003905560139091529020548311613920576001600160a01b0381165f81815260136020526040902080548590039055613abe5750600191613ab6916128a5565b945b0161334f565b909591905f805b878110613b0e575b5015613ade575b5050600190613ab8565b91946001928392906001600160a01b0316613af98387612ebc565b52613b048287612ebc565b520193905f613ad4565b6001600160a01b0383811690613b248389612ebc565b511614613b3357600101613ac5565b9050613b53613b4c84613b46848a612ebc565b516128a5565b9187612ebc565b5260015f613acd565b6306a1762560e11b5f5260045ffd5b630b6fac0360e41b5f5260045ffd5b96975094509291921590811591613eec575b508015613ee1575b610e86576001600160401b036101c051165f52600460205260405f205403613b6b5760026101805101549060018060a01b03600161018051015416926001600160401b036101c051165f52601260205260405f2060018060a01b0385165f5260205260405f20548311613b5c57613c09614987565b10613920576001600160401b036101c051165f52601260205260405f2060018060a01b0384165f5260205260405f20613c438382546126fe565b9055825f52601360205260405f20613c5c8382546126fe565b9055613c726003610180510154601454906126fe565b9280613e7657508115613e14575090613c8a916128a5565b613c9481836149ef565b6040519060408201905f835260208301525f516020615a475f395f51905f526101e05192806001600160401b036101c05116930390a45b60ff600761018051015460e81c16613ce2816125e7565b15613e06575b61289060145460ff600761018051015460e81c16613d0d3661016051610140516128b2565b90613d166146b3565b151580613df4575b15613de757613d2b615987565b613d34816125e7565b613d9657604051907f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d6101e0519280613d816001600160401b036101c05116946101005160018985614a98565b0390a35b6015546001600160a01b03166149ef565b604051907f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed666101e0519280613ddf6001600160401b036101c05116946101005160018985614a98565b0390a3613d85565b613def61594a565b613d2b565b50613e016101e05161495f565b613d1e565b600160075401600755613ce8565b9192505081613e25575b5050613ccb565b6001600160a01b031690613e3981836149ef565b6040519060408201905f835260208301525f516020615a475f395f51905f526101e05192806001600160401b036101c05116930390a45f80613e1e565b8291939492613e8f575b50505081613e25575050613ccb565b613e9a828483614a4e565b604080516101e0516101c0516001600160a01b0394909416825260208201949094526001600160401b03909216915f516020615a475f395f51905f529190a45f8080613e80565b5060c0511515613b94565b905015155f613b8c565b508361329c565b90506020813d602011613f2f575b81613f1860209383612742565b81010312610d8557613f2990612898565b5f613272565b3d9150613f0b565b825180516001600160a01b03908116875260208281015190911681880152604091820151918701919091528a9850606090950194909201916001016131f3565b606082360312610d85576040519060608201908282106001600160401b03831117611a3e57606092602092604052613fae85612546565b8152613fbb838601612546565b83820152604085013560408201528152019101906130f2565b6101e0515f52600860205260405f2061018052612fd6565b50613ff96101e05161495f565b612fb8565b6302e8145360e61b5f5260045ffd5b6140156146b3565b61404157600b5490600a54809211918261402e57505090565b9091505f52600960205260405f20541490565b61404a9061495f565b90565b6140556128e8565b5061405e6146b3565b6140c55761406a614697565b61407d576140766128e8565b905f905f90565b600a545f52600960205260405f20545f52600860205260405f20906001600160401b03600783015460a01c165f5260046020526140be60405f205492612a10565b9190600190565b601b545f52601a60205260405f20545f52601960205260405f20906001600160401b03600783015460a01c165f5260046020526140be60405f205492612a10565b51906001600160a01b0382168203610d8557565b602081830312610d85578051906001600160401b038211610d8557019080601f83011215610d8557815161414d81612c62565b9261415b6040519485612742565b81845260208085019260051b820101928311610d8557602001905b8282106141835750505090565b6020809161419084614106565b815201910190614176565b335f9081527f5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b7602052604090205460ff16156141d357565b63e2517d3f60e01b5f52336004527fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4260245260445ffd5b335f9081527f3f65c445f65be3b772e454052bf8f1f970a2679db965b2efcdbc04d3dd490789602052604090205460ff161561424257565b63e2517d3f60e01b5f52336004527fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c9591760245260445ffd5b5f8181526020818152604080832033845290915290205460ff161561429b5750565b63e2517d3f60e01b5f523360045260245260445ffd5b6001600160a01b0381165f9081525f516020615a675f395f51905f52602052604090205460ff16614357576001600160a01b03165f8181525f516020615a675f395f51905f5260205260408120805460ff191660011790553391907ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b505f90565b5f818152602081815260408083206001600160a01b038616845290915290205460ff166143de575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b6002600154146143f5576002600155565b633ee5aeb560e01b5f5260045ffd5b60405163a9059cbb60e01b60208201526001600160a01b0392909216602483015260448083019390935291815261289091614440606483612742565b615574565b6040516370a0823160e01b81523060048201526001600160a01b039190911691602082602481865afa918215610d7a575f9261453e575b506024926144c3602092604051906323b872dd60e01b8583015260018060a01b031686820152306044820152866064820152606481526144bd608482612742565b82615574565b6040516370a0823160e01b815230600482015293849182905afa8015610d7a575f9061450a575b6144f492506126fe565b036144fb57565b63b0a6ea2960e01b5f5260045ffd5b506020823d602011614536575b8161452460209383612742565b81010312610d85576144f491516144ea565b3d9150614517565b9091506020813d60201161456b575b8161455a60209383612742565b81010312610d85575190602461447c565b3d915061454d565b6001600160a01b0381165f9081525f516020615a675f395f51905f52602052604090205460ff1615614357576001600160a01b03165f8181525f516020615a675f395f51905f5260205260408120805460ff191690553391907ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16156143de575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b600b54600a548082116146aa5750505f90565b61404a916126fe565b601c54601b548082116146aa5750505f90565b600260038201549101548082116146aa5750505f90565b8181106146e8575050565b5f81556001016146dd565b90815f52600860205260405f208151815560018060a01b03602083015116600182019060018060a01b03166001600160601b0360a01b8254161790556040820151600282015560608201516003820155608082015160048201556005810160a08301518051906001600160401b038211611a3e57614771835461293a565b601f81116148e0575b50602090601f83116001146148795791806147ae9260079695945f9261486e5750508160011b915f199060031b1c19161790565b90555b60c08301516006820180546001600160a01b03199081166001600160a01b039384161790915560e08086015194909301805490911693909116929092178083556101008401516101208501516101409095015193949093921b60ff60e01b169161481a816125e7565b614823816125e7565b60ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600b545f52600960205260405f20556001600b5401600b55565b015190505f80611879565b90601f19831691845f52815f20925f5b8181106148c85750916001939185600798979694106148b0575b505050811b0190556147b1565b01515f1960f88460031b161c191690555f80806148a3565b92936020600181928786015181550195019301614889565b61490a90845f5260205f20601f850160051c81019160208610611a3457601f0160051c01906146dd565b5f61477a565b9061491a82612c62565b6149276040519182612742565b8281528092614938601f1991612c62565b01905f5b82811061494857505050565b6020906149536128e8565b8282850101520161493c565b601c5490601b54809211918261497457505090565b9091505f52601a60205260405f20541490565b5f805260116020527f4ad3b33220dddc71b994a52d72c06b10862965f7d926534c05c00fb7e819e7b75461404a906149bf90476126fe565b5f805260136020527f8fa6efc3be94b5b348b21fea823fe8d100408cee9b7f90524494500445d8ff6c54906126fe565b6001600160a01b03165f9081527f6e0956cda88cad152e89927e53611735b61a5c762d1428573c6931b0a5efcb01602052604090208054614a319083906128a5565b90555f80526011602052614a4a60405f209182546128a5565b9055565b60018060a01b031690815f52601060205260405f209060018060a01b03165f5260205260405f20614a808382546128a5565b90555f526011602052614a4a60405f209182546128a5565b909392919381526002841015611e3057614ac160809261404a9560208401526040830190612e77565b81606082015201906125c3565b9035601e1982360301811215610d855701602081359101916001600160401b038211610d85578160051b36038313610d8557565b906040810191614b128180614ace565b809460408552526060830160608560051b85010194825f90601e19813603015b838310614b7d57505050505050806020614b4d920190614ace565b82840360209093019290925281835290916001600160fb1b038311610d855760209260051b809284830137010190565b909192939497605f19888203018652883582811215610d8557830190602082359201916001600160401b038111610d85578036038313610d8557614bc76020928392600195612abe565b9a0196019493019190614b32565b98969b9997959293614c27949b9e9d9b614c19936001600160401b0360209894168c52878c015260408b015260608a015261016060808a0152610160890190614b02565b9087820360a0890152614b02565b85810360c08701528281520196905f5b818110614c6f5750505096614c619161404a979860e0850152610100840152610120830190612e77565b610140818503910152612abe565b909197606080600192838060a01b03614c878d612546565b1681528b614c9d6020868060a01b039201612546565b16602082015260408c81013590820152019901929101614c37565b6103a0526103c0526103e05261042052610280526102a0526102c0526102e05261030052610320526103405261036052610380525f610400526001600160401b036103a051165f52601760205260018060a01b0360405f205416610240526102405115612890576001600160401b036103a051165f52601260205260405f205f805260205260405f205480615266575b50600f5460405163024ece8960e01b8152905f90829060049082906001600160a01b03165afa908115610d7a575f9161524c575b5080515f5b8181036150e7575050506001610240513b15610d85576040516102605263480a21bf60e11b61026051525f6102605161026051614df761038051610360516103405161032051610300516102e0516102c0516102a05161028051610420516103e0516103c0516103a05160046102605101614bd5565b036102605183610240515af161022052610220516150d0575b610220516150cb5750610400515b6040519015158152610240519061042051907f1abe034681574802e9155b5fc48d0808dbcc65bec8903ecf5c9c03dbd4b6c9a960206001600160401b036103a0511692a4604051610200526322441a3f60e01b6102005152610400516102005161020051614ec561038051610360516103405161032051610300516102e0516102c0516102a05161028051610420516103e0516103c0516103a05160046102005101614bd5565b036102005161040051610240515af180156150bd5761040051908190614fd9575b61040051908051915b828103614f3f575050506040519015158152610240519061042051907f89eb244dd0d9dc2a460befb330ab707a40d812acbd57edd23b4818cb2be608fc60206001600160401b036103a0511692a4565b806020614f4e60019385612ebc565b5101516001600160401b036103a051166104005152601260205260406104005120838060a01b03614f7f8487612ebc565b515116848060a01b03165f5260205260405f20556020614f9f8285612ebc565b510151828060a01b03614fb28386612ebc565b51511661040051526013602052614fd1604061040051209182546128a5565b905501614eef565b50503d8061040051610200513e614ff38161020051612742565b61020051016040610200518203126150b65761501161020051612898565b60206102005101516001600160401b0381116150b657610200510182601f820112156150b65780519061504382612c62565b936150516040519586612742565b82855260208086019360061b830101918183116150b657602001925b82841061507e575050505090614ee6565b6040848303126150b6576020604091825161509881612727565b6150a187614106565b8152828701518382015281520193019261506d565b6104005180fd5b6040513d61040051823e3d90fd5b614e1e565b6150dd5f61026051612742565b5f61040052614e10565b6103a0516001600160401b03165f9081526012602052604090206001600160a01b036151138386612ebc565b511660018060a01b03165f5260205260405f20549081615138575b6001915001614d81565b6001600160a01b0361514a8286612ebc565b5160405163a9059cbb60e01b8152610240516001600160a01b03166004820152602481018590529391602091859160449183915f91165af1805f91615211575b6001945061519a575b505061512e565b6151a5575b80615193565b6001600160401b036103a051165f52601260205260405f20838060a01b036151cd8488612ebc565b5116848060a01b03165f526020525f6040812055828060a01b036151f18387612ebc565b51165f52601360205261520960405f209182546126fe565b90555f61519f565b90506020843d8211615244575b8161522b60209383612742565b81010312610d855761523e600194612898565b9061518a565b3d915061521e565b61526091503d805f833e6106bf8183612742565b5f614d7c565b5f80808084610240515af161527961277e565b5015614d48576001600160401b036103a051165f52601260205260405f205f80526020525f60408120555f805260136020526152ba60405f209182546126fe565b90555f614d48565b60075f91828155826001820155826002820155826003820155826004820155600581016152ef815461293a565b9081615302575b50508260068201550155565b81601f8693116001146153195750555b5f806152f6565b8183526020832061533591601f0160051c8101906001016146dd565b8082528160208120915555615312565b600a54600b545f905b80830361536757509050600a54600b5560075401600755565b825f52600960205260405f20545f52600860205260405f2092600784015460ff8160e81c16615395816125e7565b15615474575b8460026001949596018054806153e6575b50505050805f52600960205260405f20545f5260086020526153d060405f206152c2565b5f8181526009602052604081205501919061534e565b6001600160401b0361546b9460a01c165f52601260205260405f209086840191878060a01b03835416888060a01b03165f5260205261542a60405f209182546126fe565b90558154868060a01b038254165f52601360205261544d60405f209182546126fe565b90555460069092015490549160a086901b8690039182169116614a4e565b5f8080806153ac565b6001939093019261539b565b91909161548d60086146c6565b9283821080159061556c575b615530576154a790826128a5565b92808411615528575b506154d56154c161142383866126fe565b936154cf600a5493846128a5565b926128a5565b905f905b8281106154e557505050565b60018181925f52600960205260405f20545f52600860205261550960405f20612a10565b6155138589612ebc565b5261551e8488612ebc565b50019101906154d9565b92505f6154b0565b5090915050604051615543602082612742565b5f81525f805b81811061555557505090565b6020906155606128e8565b82828601015201615549565b508015615499565b905f602091828151910182855af115610d7a575f513d6155c357506001600160a01b0381163b155b6155a35750565b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b6001141561559c565b60ff81146156125760ff811690601f821161560357604051916155f0604084612742565b6020808452838101919036833783525290565b632cd44ac360e21b5f5260045ffd5b50604051600254815f6156248361293a565b80835292600181169081156156a75750600114615648575b61404a92500382612742565b5060025f90815290917f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace5b81831061568b57505090602061404a9282010161563c565b6020919350806001915483858801015201910190918392615673565b6020925061404a94915060ff191682840152151560051b82010161563c565b60ff81146156ea5760ff811690601f821161560357604051916155f0604084612742565b50604051600354815f6156fc8361293a565b80835292600181169081156156a7575060011461571f5761404a92500382612742565b5060035f90815290917fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b5b81831061576257505090602061404a9282010161563c565b602091935080600191548385880101520191019091839261574a565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316148061586b575b156157d9577f000000000000000000000000000000000000000000000000000000000000000090565b60405160208101907f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f82527f000000000000000000000000000000000000000000000000000000000000000060408201527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a082015260a08152612b5760c082612742565b507f000000000000000000000000000000000000000000000000000000000000000046146157b0565b81519190604183036158c4576158bd9250602082015190606060408401519301515f1a906159c4565b9192909190565b50505f9160029190565b6158d7816125e7565b806158e0575050565b6158e9816125e7565b600181036159005763f645eedf60e01b5f5260045ffd5b615909816125e7565b60028103615924575063fce698f760e01b5f5260045260245ffd5b600390615930816125e7565b146159385750565b6335e2f38360e21b5f5260045260245ffd5b600a545f52600960205260405f20545f52600860205261596c60405f206152c2565b600a545f5260096020525f60408120556001600a5401600a55565b601b545f52601a60205260405f20545f5260196020526159a960405f206152c2565b601b545f52601a6020525f60408120556001601b5401601b55565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411615a3b579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15610d7a575f516001600160a01b03811615615a3157905f905f90565b505f906001905f90565b5050505f916003919056fec9961aaccfc3406b40d0d25a5b83b8a4721e5f454116d7e22bd8038c715140aa740c5e3e456bed56f053f960110118ba9b95a1f5359a82283516fb2e81b6e37eb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236a2646970667358221220ce95caff1873baa9ced12e91ca55e3199511cd8c3306686bffbe0dbdc05f467e64736f6c634300081e00332f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d009b85fa08e135b4afefc8a40566229b01e66445186a3676642ffbacd3d93b0e5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b73f65c445f65be3b772e454052bf8f1f970a2679db965b2efcdbc04d3dd490789fc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184cdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec42740c5e3e456bed56f053f960110118ba9b95a1f5359a82283516fb2e81b6e37e",
}

// ProcessorEndpoint is an auto generated Go binding around an Ethereum contract.
type ProcessorEndpoint struct {
	abi abi.ABI
}

// NewProcessorEndpoint creates a new instance of ProcessorEndpoint.
func NewProcessorEndpoint() *ProcessorEndpoint {
	parsed, err := ProcessorEndpointMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ProcessorEndpoint{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ProcessorEndpoint) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _teeAuthenticator, address _authorityRegistry, address updateStatusOperator, address admin, address resetOperator, uint256 _minFeePerRequest, address _tokenAllowlist) returns()
func (processorEndpoint *ProcessorEndpoint) PackConstructor(_teeAuthenticator common.Address, _authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address, resetOperator common.Address, _minFeePerRequest *big.Int, _tokenAllowlist common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("", _teeAuthenticator, _authorityRegistry, updateStatusOperator, admin, resetOperator, _minFeePerRequest, _tokenAllowlist)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackADMIN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a0acc6a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackADMIN() []byte {
	enc, err := processorEndpoint.abi.Pack("ADMIN")
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
func (processorEndpoint *ProcessorEndpoint) TryPackADMIN() ([]byte, error) {
	return processorEndpoint.abi.Pack("ADMIN")
}

// UnpackADMIN is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a0acc6a.
//
// Solidity: function ADMIN() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackADMIN(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("ADMIN", data)
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
func (processorEndpoint *ProcessorEndpoint) PackDEFAULTADMINROLE() []byte {
	enc, err := processorEndpoint.abi.Pack("DEFAULT_ADMIN_ROLE")
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
func (processorEndpoint *ProcessorEndpoint) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return processorEndpoint.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
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
func (processorEndpoint *ProcessorEndpoint) PackDEPLOYERROLE() []byte {
	enc, err := processorEndpoint.abi.Pack("DEPLOYER_ROLE")
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
func (processorEndpoint *ProcessorEndpoint) TryPackDEPLOYERROLE() ([]byte, error) {
	return processorEndpoint.abi.Pack("DEPLOYER_ROLE")
}

// UnpackDEPLOYERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xecd00261.
//
// Solidity: function DEPLOYER_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackDEPLOYERROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("DEPLOYER_ROLE", data)
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
func (processorEndpoint *ProcessorEndpoint) PackPROTOCOLVERSION() []byte {
	enc, err := processorEndpoint.abi.Pack("PROTOCOL_VERSION")
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
func (processorEndpoint *ProcessorEndpoint) TryPackPROTOCOLVERSION() ([]byte, error) {
	return processorEndpoint.abi.Pack("PROTOCOL_VERSION")
}

// UnpackPROTOCOLVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa3aa460.
//
// Solidity: function PROTOCOL_VERSION() view returns(uint8)
func (processorEndpoint *ProcessorEndpoint) UnpackPROTOCOLVERSION(data []byte) (uint8, error) {
	out, err := processorEndpoint.abi.Unpack("PROTOCOL_VERSION", data)
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
func (processorEndpoint *ProcessorEndpoint) PackREQUESTAUTHORIZATIONTYPEHASH() []byte {
	enc, err := processorEndpoint.abi.Pack("REQUEST_AUTHORIZATION_TYPEHASH")
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
func (processorEndpoint *ProcessorEndpoint) TryPackREQUESTAUTHORIZATIONTYPEHASH() ([]byte, error) {
	return processorEndpoint.abi.Pack("REQUEST_AUTHORIZATION_TYPEHASH")
}

// UnpackREQUESTAUTHORIZATIONTYPEHASH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2b52b07c.
//
// Solidity: function REQUEST_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackREQUESTAUTHORIZATIONTYPEHASH(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("REQUEST_AUTHORIZATION_TYPEHASH", data)
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
func (processorEndpoint *ProcessorEndpoint) PackRESETOPERATOR() []byte {
	enc, err := processorEndpoint.abi.Pack("RESET_OPERATOR")
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
func (processorEndpoint *ProcessorEndpoint) TryPackRESETOPERATOR() ([]byte, error) {
	return processorEndpoint.abi.Pack("RESET_OPERATOR")
}

// UnpackRESETOPERATOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x86f4a07d.
//
// Solidity: function RESET_OPERATOR() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackRESETOPERATOR(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("RESET_OPERATOR", data)
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
func (processorEndpoint *ProcessorEndpoint) PackUPDATESTATUSROLE() []byte {
	enc, err := processorEndpoint.abi.Pack("UPDATE_STATUS_ROLE")
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
func (processorEndpoint *ProcessorEndpoint) TryPackUPDATESTATUSROLE() ([]byte, error) {
	return processorEndpoint.abi.Pack("UPDATE_STATUS_ROLE")
}

// UnpackUPDATESTATUSROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1664deeb.
//
// Solidity: function UPDATE_STATUS_ROLE() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackUPDATESTATUSROLE(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("UPDATE_STATUS_ROLE", data)
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
func (processorEndpoint *ProcessorEndpoint) PackAddAllowedDeployer(deployer common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("addAllowedDeployer", deployer)
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
func (processorEndpoint *ProcessorEndpoint) TryPackAddAllowedDeployer(deployer common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("addAllowedDeployer", deployer)
}

// PackAdminReset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c5b9b00.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function adminReset() returns()
func (processorEndpoint *ProcessorEndpoint) PackAdminReset() []byte {
	enc, err := processorEndpoint.abi.Pack("adminReset")
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
func (processorEndpoint *ProcessorEndpoint) TryPackAdminReset() ([]byte, error) {
	return processorEndpoint.abi.Pack("adminReset")
}

// PackAdminResetApps is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd39d765.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function adminResetApps(uint64[] appIds) returns()
func (processorEndpoint *ProcessorEndpoint) PackAdminResetApps(appIds []uint64) []byte {
	enc, err := processorEndpoint.abi.Pack("adminResetApps", appIds)
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
func (processorEndpoint *ProcessorEndpoint) TryPackAdminResetApps(appIds []uint64) ([]byte, error) {
	return processorEndpoint.abi.Pack("adminResetApps", appIds)
}

// PackAppCustody is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb7222ff5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function appCustody(uint64 , address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackAppCustody(arg0 uint64, arg1 common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("appCustody", arg0, arg1)
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
func (processorEndpoint *ProcessorEndpoint) TryPackAppCustody(arg0 uint64, arg1 common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("appCustody", arg0, arg1)
}

// UnpackAppCustody is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb7222ff5.
//
// Solidity: function appCustody(uint64 , address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackAppCustody(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("appCustody", data)
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
func (processorEndpoint *ProcessorEndpoint) PackApplicationStateRoots(arg0 uint64) []byte {
	enc, err := processorEndpoint.abi.Pack("applicationStateRoots", arg0)
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
func (processorEndpoint *ProcessorEndpoint) TryPackApplicationStateRoots(arg0 uint64) ([]byte, error) {
	return processorEndpoint.abi.Pack("applicationStateRoots", arg0)
}

// UnpackApplicationStateRoots is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7a36a891.
//
// Solidity: function applicationStateRoots(uint64 ) view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackApplicationStateRoots(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("applicationStateRoots", data)
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
func (processorEndpoint *ProcessorEndpoint) PackAuthorityRegistry() []byte {
	enc, err := processorEndpoint.abi.Pack("authorityRegistry")
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
func (processorEndpoint *ProcessorEndpoint) TryPackAuthorityRegistry() ([]byte, error) {
	return processorEndpoint.abi.Pack("authorityRegistry")
}

// UnpackAuthorityRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9ba045e.
//
// Solidity: function authorityRegistry() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackAuthorityRegistry(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("authorityRegistry", data)
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
func (processorEndpoint *ProcessorEndpoint) PackAvailableDeploySlots() []byte {
	enc, err := processorEndpoint.abi.Pack("availableDeploySlots")
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
func (processorEndpoint *ProcessorEndpoint) TryPackAvailableDeploySlots() ([]byte, error) {
	return processorEndpoint.abi.Pack("availableDeploySlots")
}

// UnpackAvailableDeploySlots is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1753520f.
//
// Solidity: function availableDeploySlots() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackAvailableDeploySlots(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("availableDeploySlots", data)
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
func (processorEndpoint *ProcessorEndpoint) PackClaim(tokenAddress common.Address, payee common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("claim", tokenAddress, payee)
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
func (processorEndpoint *ProcessorEndpoint) TryPackClaim(tokenAddress common.Address, payee common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("claim", tokenAddress, payee)
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (processorEndpoint *ProcessorEndpoint) PackEip712Domain() []byte {
	enc, err := processorEndpoint.abi.Pack("eip712Domain")
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
func (processorEndpoint *ProcessorEndpoint) TryPackEip712Domain() ([]byte, error) {
	return processorEndpoint.abi.Pack("eip712Domain")
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
func (processorEndpoint *ProcessorEndpoint) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := processorEndpoint.abi.Unpack("eip712Domain", data)
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
func (processorEndpoint *ProcessorEndpoint) PackFacilitatorNonces(arg0 common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("facilitatorNonces", arg0)
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
func (processorEndpoint *ProcessorEndpoint) TryPackFacilitatorNonces(arg0 common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("facilitatorNonces", arg0)
}

// UnpackFacilitatorNonces is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd72e957a.
//
// Solidity: function facilitatorNonces(address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackFacilitatorNonces(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("facilitatorNonces", data)
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
func (processorEndpoint *ProcessorEndpoint) PackFeeCollector() []byte {
	enc, err := processorEndpoint.abi.Pack("feeCollector")
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
func (processorEndpoint *ProcessorEndpoint) TryPackFeeCollector() ([]byte, error) {
	return processorEndpoint.abi.Pack("feeCollector")
}

// UnpackFeeCollector is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackFeeCollector(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("feeCollector", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74202d9b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, idx *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, tokenAddress, assetAmount, idx)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74202d9b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, idx *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, tokenAddress, assetAmount, idx)
}

// UnpackGenerateRequestId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x74202d9b.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackGenerateRequestId(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("generateRequestId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetDeployedAppIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7596ed43.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getDeployedAppIds() view returns(uint64[])
func (processorEndpoint *ProcessorEndpoint) PackGetDeployedAppIds() []byte {
	enc, err := processorEndpoint.abi.Pack("getDeployedAppIds")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetDeployedAppIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7596ed43.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getDeployedAppIds() view returns(uint64[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetDeployedAppIds() ([]byte, error) {
	return processorEndpoint.abi.Pack("getDeployedAppIds")
}

// UnpackGetDeployedAppIds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7596ed43.
//
// Solidity: function getDeployedAppIds() view returns(uint64[])
func (processorEndpoint *ProcessorEndpoint) UnpackGetDeployedAppIds(data []byte) ([]uint64, error) {
	out, err := processorEndpoint.abi.Unpack("getDeployedAppIds", data)
	if err != nil {
		return *new([]uint64), err
	}
	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)
	return out0, nil
}

// PackGetFacilitatorNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1256c918.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFacilitatorNonce(address user) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackGetFacilitatorNonce(user common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("getFacilitatorNonce", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFacilitatorNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1256c918.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFacilitatorNonce(address user) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) TryPackGetFacilitatorNonce(user common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("getFacilitatorNonce", user)
}

// UnpackGetFacilitatorNonce is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1256c918.
//
// Solidity: function getFacilitatorNonce(address user) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackGetFacilitatorNonce(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("getFacilitatorNonce", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetNextPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc404c359.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getNextPendingRequest() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8), bytes32, bool success)
func (processorEndpoint *ProcessorEndpoint) PackGetNextPendingRequest() []byte {
	enc, err := processorEndpoint.abi.Pack("getNextPendingRequest")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetNextPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc404c359.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getNextPendingRequest() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8), bytes32, bool success)
func (processorEndpoint *ProcessorEndpoint) TryPackGetNextPendingRequest() ([]byte, error) {
	return processorEndpoint.abi.Pack("getNextPendingRequest")
}

// GetNextPendingRequestOutput serves as a container for the return parameters of contract
// method GetNextPendingRequest.
type GetNextPendingRequestOutput struct {
	Arg0    StructsPendingRequest
	Arg1    [32]byte
	Success bool
}

// UnpackGetNextPendingRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc404c359.
//
// Solidity: function getNextPendingRequest() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8), bytes32, bool success)
func (processorEndpoint *ProcessorEndpoint) UnpackGetNextPendingRequest(data []byte) (GetNextPendingRequestOutput, error) {
	out, err := processorEndpoint.abi.Unpack("getNextPendingRequest", data)
	outstruct := new(GetNextPendingRequestOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new(StructsPendingRequest)).(*StructsPendingRequest)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Success = *abi.ConvertType(out[2], new(bool)).(*bool)
	return *outstruct, nil
}

// PackGetPendingRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a1f712.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingRequests() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) PackGetPendingRequests() []byte {
	enc, err := processorEndpoint.abi.Pack("getPendingRequests")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a1f712.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingRequests() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequests() ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequests")
}

// UnpackGetPendingRequests is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80a1f712.
//
// Solidity: function getPendingRequests() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) UnpackGetPendingRequests(data []byte) ([]StructsPendingRequest, error) {
	out, err := processorEndpoint.abi.Unpack("getPendingRequests", data)
	if err != nil {
		return *new([]StructsPendingRequest), err
	}
	out0 := *abi.ConvertType(out[0], new([]StructsPendingRequest)).(*[]StructsPendingRequest)
	return out0, nil
}

// PackGetPendingRequestsPage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbe1e3725.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingRequestsPage(uint256 offset, uint256 limit) view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) PackGetPendingRequestsPage(offset *big.Int, limit *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("getPendingRequestsPage", offset, limit)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingRequestsPage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbe1e3725.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingRequestsPage(uint256 offset, uint256 limit) view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequestsPage(offset *big.Int, limit *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequestsPage", offset, limit)
}

// UnpackGetPendingRequestsPage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbe1e3725.
//
// Solidity: function getPendingRequestsPage(uint256 offset, uint256 limit) view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) UnpackGetPendingRequestsPage(data []byte) ([]StructsPendingRequest, error) {
	out, err := processorEndpoint.abi.Unpack("getPendingRequestsPage", data)
	if err != nil {
		return *new([]StructsPendingRequest), err
	}
	out0 := *abi.ConvertType(out[0], new([]StructsPendingRequest)).(*[]StructsPendingRequest)
	return out0, nil
}

// PackGetPendingRequestsSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e339384.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingRequestsSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackGetPendingRequestsSize() []byte {
	enc, err := processorEndpoint.abi.Pack("getPendingRequestsSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingRequestsSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e339384.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingRequestsSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequestsSize() ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequestsSize")
}

// UnpackGetPendingRequestsSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e339384.
//
// Solidity: function getPendingRequestsSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackGetPendingRequestsSize(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("getPendingRequestsSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("getRoleAdmin", role)
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
func (processorEndpoint *ProcessorEndpoint) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetTriggerQueueSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47101cae.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getTriggerQueueSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackGetTriggerQueueSize() []byte {
	enc, err := processorEndpoint.abi.Pack("getTriggerQueueSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetTriggerQueueSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47101cae.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getTriggerQueueSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) TryPackGetTriggerQueueSize() ([]byte, error) {
	return processorEndpoint.abi.Pack("getTriggerQueueSize")
}

// UnpackGetTriggerQueueSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x47101cae.
//
// Solidity: function getTriggerQueueSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackGetTriggerQueueSize(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("getTriggerQueueSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (processorEndpoint *ProcessorEndpoint) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("grantRole", role, account)
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
func (processorEndpoint *ProcessorEndpoint) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("hasRole", role, account)
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
func (processorEndpoint *ProcessorEndpoint) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) UnpackHasRole(data []byte) (bool, error) {
	out, err := processorEndpoint.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsAllowedDeployer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6b288d20.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isAllowedDeployer(address deployer) view returns(bool allowed)
func (processorEndpoint *ProcessorEndpoint) PackIsAllowedDeployer(deployer common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("isAllowedDeployer", deployer)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsAllowedDeployer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6b288d20.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isAllowedDeployer(address deployer) view returns(bool allowed)
func (processorEndpoint *ProcessorEndpoint) TryPackIsAllowedDeployer(deployer common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("isAllowedDeployer", deployer)
}

// UnpackIsAllowedDeployer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6b288d20.
//
// Solidity: function isAllowedDeployer(address deployer) view returns(bool allowed)
func (processorEndpoint *ProcessorEndpoint) UnpackIsAllowedDeployer(data []byte) (bool, error) {
	out, err := processorEndpoint.abi.Unpack("isAllowedDeployer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsCurrentPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9fabc36f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) PackIsCurrentPendingRequest(requestId [32]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("isCurrentPendingRequest", requestId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsCurrentPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9fabc36f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) TryPackIsCurrentPendingRequest(requestId [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("isCurrentPendingRequest", requestId)
}

// UnpackIsCurrentPendingRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9fabc36f.
//
// Solidity: function isCurrentPendingRequest(bytes32 requestId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) UnpackIsCurrentPendingRequest(data []byte) (bool, error) {
	out, err := processorEndpoint.abi.Unpack("isCurrentPendingRequest", data)
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
func (processorEndpoint *ProcessorEndpoint) PackMaxNumOfApplications() []byte {
	enc, err := processorEndpoint.abi.Pack("maxNumOfApplications")
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
func (processorEndpoint *ProcessorEndpoint) TryPackMaxNumOfApplications() ([]byte, error) {
	return processorEndpoint.abi.Pack("maxNumOfApplications")
}

// UnpackMaxNumOfApplications is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8eb8791a.
//
// Solidity: function maxNumOfApplications() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackMaxNumOfApplications(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("maxNumOfApplications", data)
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
func (processorEndpoint *ProcessorEndpoint) PackMaxQueueSize() []byte {
	enc, err := processorEndpoint.abi.Pack("maxQueueSize")
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
func (processorEndpoint *ProcessorEndpoint) TryPackMaxQueueSize() ([]byte, error) {
	return processorEndpoint.abi.Pack("maxQueueSize")
}

// UnpackMaxQueueSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf7d9cc1e.
//
// Solidity: function maxQueueSize() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackMaxQueueSize(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("maxQueueSize", data)
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
func (processorEndpoint *ProcessorEndpoint) PackMinFeePerRequest() []byte {
	enc, err := processorEndpoint.abi.Pack("minFeePerRequest")
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
func (processorEndpoint *ProcessorEndpoint) TryPackMinFeePerRequest() ([]byte, error) {
	return processorEndpoint.abi.Pack("minFeePerRequest")
}

// UnpackMinFeePerRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9c8d416f.
//
// Solidity: function minFeePerRequest() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackMinFeePerRequest(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("minFeePerRequest", data)
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
func (processorEndpoint *ProcessorEndpoint) PackPendingClaims(arg0 common.Address, arg1 common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("pendingClaims", arg0, arg1)
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
func (processorEndpoint *ProcessorEndpoint) TryPackPendingClaims(arg0 common.Address, arg1 common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("pendingClaims", arg0, arg1)
}

// UnpackPendingClaims is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x840059ec.
//
// Solidity: function pendingClaims(address , address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackPendingClaims(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("pendingClaims", data)
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
func (processorEndpoint *ProcessorEndpoint) PackRemoveAllowedDeployer(deployer common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("removeAllowedDeployer", deployer)
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
func (processorEndpoint *ProcessorEndpoint) TryPackRemoveAllowedDeployer(deployer common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("removeAllowedDeployer", deployer)
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (processorEndpoint *ProcessorEndpoint) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("renounceRole", role, callerConfirmation)
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
func (processorEndpoint *ProcessorEndpoint) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRequestById is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5423112f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestById(bytes32 id) view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8))
func (processorEndpoint *ProcessorEndpoint) PackRequestById(id [32]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("requestById", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestById is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5423112f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestById(bytes32 id) view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8))
func (processorEndpoint *ProcessorEndpoint) TryPackRequestById(id [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("requestById", id)
}

// UnpackRequestById is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5423112f.
//
// Solidity: function requestById(bytes32 id) view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8))
func (processorEndpoint *ProcessorEndpoint) UnpackRequestById(data []byte) (StructsPendingRequest, error) {
	out, err := processorEndpoint.abi.Unpack("requestById", data)
	if err != nil {
		return *new(StructsPendingRequest), err
	}
	out0 := *abi.ConvertType(out[0], new(StructsPendingRequest)).(*StructsPendingRequest)
	return out0, nil
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (processorEndpoint *ProcessorEndpoint) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("revokeRole", role, account)
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
func (processorEndpoint *ProcessorEndpoint) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("revokeRole", role, account)
}

// PackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8537c278.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, (bytes[],bytes32[]) userEventData, (bytes[],bytes32[]) appEventData, (address,address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, uint8 errorCode, string errorMsg, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) PackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, userEventData StructsEventData, appEventData StructsEventData, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, errorCode uint8, errorMsg string, signature []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, userEventData, appEventData, withdrawalRequests, refund, applicationFees, errorCode, errorMsg, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8537c278.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, (bytes[],bytes32[]) userEventData, (bytes[],bytes32[]) appEventData, (address,address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, uint8 errorCode, string errorMsg, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, userEventData StructsEventData, appEventData StructsEventData, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, errorCode uint8, errorMsg string, signature []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, userEventData, appEventData, withdrawalRequests, refund, applicationFees, errorCode, errorMsg, signature)
}

// PackSubmitDeployRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4a0d0495.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitDeployRequest(uint8 protocolVersion, bytes payload) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitDeployRequest(protocolVersion uint8, payload []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("submitDeployRequest", protocolVersion, payload)
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
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitDeployRequest(protocolVersion uint8, payload []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitDeployRequest", protocolVersion, payload)
}

// UnpackSubmitDeployRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4a0d0495.
//
// Solidity: function submitDeployRequest(uint8 protocolVersion, bytes payload) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackSubmitDeployRequest(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("submitDeployRequest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2fbfa0d5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, maxFeeValue *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, tokenAddress, assetAmount, maxFeeValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2fbfa0d5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, maxFeeValue *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, tokenAddress, assetAmount, maxFeeValue)
}

// UnpackSubmitRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2fbfa0d5.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackSubmitRequest(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("submitRequest", data)
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
func (processorEndpoint *ProcessorEndpoint) PackSubmitRequestFor(sender common.Address, protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, deadline *big.Int, requestSignature []byte, depositPermit []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("submitRequestFor", sender, protocolVersion, applicationId, requestType, payload, tokenAddress, assetAmount, deadline, requestSignature, depositPermit)
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
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitRequestFor(sender common.Address, protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, deadline *big.Int, requestSignature []byte, depositPermit []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitRequestFor", sender, protocolVersion, applicationId, requestType, payload, tokenAddress, assetAmount, deadline, requestSignature, depositPermit)
}

// UnpackSubmitRequestFor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbcf7a3c9.
//
// Solidity: function submitRequestFor(address sender, uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 deadline, bytes requestSignature, bytes depositPermit) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackSubmitRequestFor(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("submitRequestFor", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSubmitTriggerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x55c0be5c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitTriggerRequest(uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 maxFeeValue, address sender, address facilitator) returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitTriggerRequest(requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, maxFeeValue *big.Int, sender common.Address, facilitator common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("submitTriggerRequest", requestType, payload, tokenAddress, assetAmount, maxFeeValue, sender, facilitator)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitTriggerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x55c0be5c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitTriggerRequest(uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 maxFeeValue, address sender, address facilitator) returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitTriggerRequest(requestType uint8, payload []byte, tokenAddress common.Address, assetAmount *big.Int, maxFeeValue *big.Int, sender common.Address, facilitator common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitTriggerRequest", requestType, payload, tokenAddress, assetAmount, maxFeeValue, sender, facilitator)
}

// UnpackSubmitTriggerRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x55c0be5c.
//
// Solidity: function submitTriggerRequest(uint8 requestType, bytes payload, address tokenAddress, uint256 assetAmount, uint256 maxFeeValue, address sender, address facilitator) returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackSubmitTriggerRequest(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("submitTriggerRequest", data)
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
func (processorEndpoint *ProcessorEndpoint) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("supportsInterface", interfaceId)
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
func (processorEndpoint *ProcessorEndpoint) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (processorEndpoint *ProcessorEndpoint) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := processorEndpoint.abi.Unpack("supportsInterface", data)
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
func (processorEndpoint *ProcessorEndpoint) PackTeeAuthenticator() []byte {
	enc, err := processorEndpoint.abi.Pack("teeAuthenticator")
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
func (processorEndpoint *ProcessorEndpoint) TryPackTeeAuthenticator() ([]byte, error) {
	return processorEndpoint.abi.Pack("teeAuthenticator")
}

// UnpackTeeAuthenticator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e91d133.
//
// Solidity: function teeAuthenticator() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackTeeAuthenticator(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("teeAuthenticator", data)
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
func (processorEndpoint *ProcessorEndpoint) PackTokenAllowlist() []byte {
	enc, err := processorEndpoint.abi.Pack("tokenAllowlist")
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
func (processorEndpoint *ProcessorEndpoint) TryPackTokenAllowlist() ([]byte, error) {
	return processorEndpoint.abi.Pack("tokenAllowlist")
}

// UnpackTokenAllowlist is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf624f88a.
//
// Solidity: function tokenAllowlist() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackTokenAllowlist(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("tokenAllowlist", data)
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
func (processorEndpoint *ProcessorEndpoint) PackTotalAppCustody(arg0 common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("totalAppCustody", arg0)
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
func (processorEndpoint *ProcessorEndpoint) TryPackTotalAppCustody(arg0 common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("totalAppCustody", arg0)
}

// UnpackTotalAppCustody is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaf20c960.
//
// Solidity: function totalAppCustody(address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackTotalAppCustody(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("totalAppCustody", data)
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
func (processorEndpoint *ProcessorEndpoint) PackTotalPendingClaims(arg0 common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("totalPendingClaims", arg0)
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
func (processorEndpoint *ProcessorEndpoint) TryPackTotalPendingClaims(arg0 common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("totalPendingClaims", arg0)
}

// UnpackTotalPendingClaims is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x194b6be8.
//
// Solidity: function totalPendingClaims(address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackTotalPendingClaims(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("totalPendingClaims", data)
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
func (processorEndpoint *ProcessorEndpoint) PackTriggerContracts(arg0 uint64) []byte {
	enc, err := processorEndpoint.abi.Pack("triggerContracts", arg0)
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
func (processorEndpoint *ProcessorEndpoint) TryPackTriggerContracts(arg0 uint64) ([]byte, error) {
	return processorEndpoint.abi.Pack("triggerContracts", arg0)
}

// UnpackTriggerContracts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1cc76c38.
//
// Solidity: function triggerContracts(uint64 ) view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackTriggerContracts(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("triggerContracts", data)
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
func (processorEndpoint *ProcessorEndpoint) PackTriggersToAppIds(arg0 common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("triggersToAppIds", arg0)
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
func (processorEndpoint *ProcessorEndpoint) TryPackTriggersToAppIds(arg0 common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("triggersToAppIds", arg0)
}

// UnpackTriggersToAppIds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18733d57.
//
// Solidity: function triggersToAppIds(address ) view returns(uint64)
func (processorEndpoint *ProcessorEndpoint) UnpackTriggersToAppIds(data []byte) (uint64, error) {
	out, err := processorEndpoint.abi.Unpack("triggersToAppIds", data)
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
func (processorEndpoint *ProcessorEndpoint) PackUpdateFeeCollector(newFeeCollector common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("updateFeeCollector", newFeeCollector)
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
func (processorEndpoint *ProcessorEndpoint) TryPackUpdateFeeCollector(newFeeCollector common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("updateFeeCollector", newFeeCollector)
}

// PackUpdateMaxNumOfApplications is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x655f2de1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateMaxNumOfApplications(uint256 newMax) returns()
func (processorEndpoint *ProcessorEndpoint) PackUpdateMaxNumOfApplications(newMax *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("updateMaxNumOfApplications", newMax)
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
func (processorEndpoint *ProcessorEndpoint) TryPackUpdateMaxNumOfApplications(newMax *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("updateMaxNumOfApplications", newMax)
}

// PackUpdateQueueThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c73eeaf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateQueueThreshold(uint256 newThreshold) returns()
func (processorEndpoint *ProcessorEndpoint) PackUpdateQueueThreshold(newThreshold *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("updateQueueThreshold", newThreshold)
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
func (processorEndpoint *ProcessorEndpoint) TryPackUpdateQueueThreshold(newThreshold *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("updateQueueThreshold", newThreshold)
}

// ProcessorEndpointAppEvent represents a AppEvent event raised by the ProcessorEndpoint contract.
type ProcessorEndpointAppEvent struct {
	ApplicationId uint64
	RequestId     [32]byte
	EventSubType  [32]byte
	Data          []byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointAppEventEventName = "AppEvent"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointAppEvent) ContractEventName() string {
	return ProcessorEndpointAppEventEventName
}

// UnpackAppEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AppEvent(uint64 indexed applicationId, bytes32 indexed requestId, bytes32 indexed eventSubType, bytes data)
func (processorEndpoint *ProcessorEndpoint) UnpackAppEventEvent(log *types.Log) (*ProcessorEndpointAppEvent, error) {
	event := "AppEvent"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointAppEvent)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointDeployRequestCompleted represents a DeployRequestCompleted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointDeployRequestCompleted struct {
	ApplicationId   uint64
	RequestId       [32]byte
	ApplicationFees *big.Int
	Status          uint8
	ErrorCode       uint8
	ErrorMessage    string
	Raw             *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointDeployRequestCompletedEventName = "DeployRequestCompleted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointDeployRequestCompleted) ContractEventName() string {
	return ProcessorEndpointDeployRequestCompletedEventName
}

// UnpackDeployRequestCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DeployRequestCompleted(uint64 indexed applicationId, bytes32 indexed requestId, uint256 applicationFees, uint8 status, uint8 errorCode, string errorMessage)
func (processorEndpoint *ProcessorEndpoint) UnpackDeployRequestCompletedEvent(log *types.Log) (*ProcessorEndpointDeployRequestCompleted, error) {
	event := "DeployRequestCompleted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointDeployRequestCompleted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointDeployRequestSubmitted represents a DeployRequestSubmitted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointDeployRequestSubmitted struct {
	ApplicationId uint64
	RequestId     [32]byte
	Sender        common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointDeployRequestSubmittedEventName = "DeployRequestSubmitted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointDeployRequestSubmitted) ContractEventName() string {
	return ProcessorEndpointDeployRequestSubmittedEventName
}

// UnpackDeployRequestSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DeployRequestSubmitted(uint64 indexed applicationId, bytes32 requestId, address indexed sender)
func (processorEndpoint *ProcessorEndpoint) UnpackDeployRequestSubmittedEvent(log *types.Log) (*ProcessorEndpointDeployRequestSubmitted, error) {
	event := "DeployRequestSubmitted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointDeployRequestSubmitted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointEIP712DomainChanged represents a EIP712DomainChanged event raised by the ProcessorEndpoint contract.
type ProcessorEndpointEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointEIP712DomainChanged) ContractEventName() string {
	return ProcessorEndpointEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (processorEndpoint *ProcessorEndpoint) UnpackEIP712DomainChangedEvent(log *types.Log) (*ProcessorEndpointEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointFeeCollectorUpdated represents a FeeCollectorUpdated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointFeeCollectorUpdated struct {
	NewFeeCollector common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointFeeCollectorUpdatedEventName = "FeeCollectorUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointFeeCollectorUpdated) ContractEventName() string {
	return ProcessorEndpointFeeCollectorUpdatedEventName
}

// UnpackFeeCollectorUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FeeCollectorUpdated(address newFeeCollector)
func (processorEndpoint *ProcessorEndpoint) UnpackFeeCollectorUpdatedEvent(log *types.Log) (*ProcessorEndpointFeeCollectorUpdated, error) {
	event := "FeeCollectorUpdated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointFeeCollectorUpdated)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointMaxNumberOfAppUpdated represents a MaxNumberOfAppUpdated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointMaxNumberOfAppUpdated struct {
	OldMax *big.Int
	NewMax *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointMaxNumberOfAppUpdatedEventName = "MaxNumberOfAppUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointMaxNumberOfAppUpdated) ContractEventName() string {
	return ProcessorEndpointMaxNumberOfAppUpdatedEventName
}

// UnpackMaxNumberOfAppUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MaxNumberOfAppUpdated(uint256 oldMax, uint256 newMax)
func (processorEndpoint *ProcessorEndpoint) UnpackMaxNumberOfAppUpdatedEvent(log *types.Log) (*ProcessorEndpointMaxNumberOfAppUpdated, error) {
	event := "MaxNumberOfAppUpdated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointMaxNumberOfAppUpdated)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointPaymentWithdrawn represents a PaymentWithdrawn event raised by the ProcessorEndpoint contract.
type ProcessorEndpointPaymentWithdrawn struct {
	TokenAddress common.Address
	Payee        common.Address
	Amount       *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointPaymentWithdrawnEventName = "PaymentWithdrawn"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointPaymentWithdrawn) ContractEventName() string {
	return ProcessorEndpointPaymentWithdrawnEventName
}

// UnpackPaymentWithdrawnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentWithdrawn(address tokenAddress, address indexed payee, uint256 amount)
func (processorEndpoint *ProcessorEndpoint) UnpackPaymentWithdrawnEvent(log *types.Log) (*ProcessorEndpointPaymentWithdrawn, error) {
	event := "PaymentWithdrawn"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointPaymentWithdrawn)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointQueueThresholdUpdated represents a QueueThresholdUpdated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointQueueThresholdUpdated struct {
	NewThreshold *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointQueueThresholdUpdatedEventName = "QueueThresholdUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointQueueThresholdUpdated) ContractEventName() string {
	return ProcessorEndpointQueueThresholdUpdatedEventName
}

// UnpackQueueThresholdUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event QueueThresholdUpdated(uint256 newThreshold)
func (processorEndpoint *ProcessorEndpoint) UnpackQueueThresholdUpdatedEvent(log *types.Log) (*ProcessorEndpointQueueThresholdUpdated, error) {
	event := "QueueThresholdUpdated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointQueueThresholdUpdated)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointRefund represents a Refund event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRefund struct {
	ApplicationId uint64
	RequestId     [32]byte
	To            common.Address
	TokenAddress  common.Address
	Amount        *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRefundEventName = "Refund"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRefund) ContractEventName() string {
	return ProcessorEndpointRefundEventName
}

// UnpackRefundEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Refund(uint64 indexed applicationId, bytes32 indexed requestId, address indexed to, address tokenAddress, uint256 amount)
func (processorEndpoint *ProcessorEndpoint) UnpackRefundEvent(log *types.Log) (*ProcessorEndpointRefund, error) {
	event := "Refund"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRefund)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointReportGenerated represents a ReportGenerated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointReportGenerated struct {
	ApplicationId uint64
	RequestId     [32]byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointReportGeneratedEventName = "ReportGenerated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointReportGenerated) ContractEventName() string {
	return ProcessorEndpointReportGeneratedEventName
}

// UnpackReportGeneratedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ReportGenerated(uint64 indexed applicationId, bytes32 indexed requestId)
func (processorEndpoint *ProcessorEndpoint) UnpackReportGeneratedEvent(log *types.Log) (*ProcessorEndpointReportGenerated, error) {
	event := "ReportGenerated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointReportGenerated)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointRequestCompleted represents a RequestCompleted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRequestCompleted struct {
	ApplicationId   uint64
	RequestId       [32]byte
	ApplicationFees *big.Int
	Status          uint8
	ErrorCode       uint8
	ErrorMessage    string
	Raw             *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRequestCompletedEventName = "RequestCompleted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRequestCompleted) ContractEventName() string {
	return ProcessorEndpointRequestCompletedEventName
}

// UnpackRequestCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestCompleted(uint64 indexed applicationId, bytes32 indexed requestId, uint256 applicationFees, uint8 status, uint8 errorCode, string errorMessage)
func (processorEndpoint *ProcessorEndpoint) UnpackRequestCompletedEvent(log *types.Log) (*ProcessorEndpointRequestCompleted, error) {
	event := "RequestCompleted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRequestCompleted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointRequestSubmitted represents a RequestSubmitted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRequestSubmitted struct {
	ApplicationId uint64
	RequestId     [32]byte
	Sender        common.Address
	Facilitator   common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRequestSubmittedEventName = "RequestSubmitted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRequestSubmitted) ContractEventName() string {
	return ProcessorEndpointRequestSubmittedEventName
}

// UnpackRequestSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestSubmitted(uint64 indexed applicationId, bytes32 indexed requestId, address indexed sender, address facilitator)
func (processorEndpoint *ProcessorEndpoint) UnpackRequestSubmittedEvent(log *types.Log) (*ProcessorEndpointRequestSubmitted, error) {
	event := "RequestSubmitted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRequestSubmitted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointRoleAdminChanged represents a RoleAdminChanged event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRoleAdminChanged) ContractEventName() string {
	return ProcessorEndpointRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (processorEndpoint *ProcessorEndpoint) UnpackRoleAdminChangedEvent(log *types.Log) (*ProcessorEndpointRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointRoleGranted represents a RoleGranted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRoleGranted) ContractEventName() string {
	return ProcessorEndpointRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (processorEndpoint *ProcessorEndpoint) UnpackRoleGrantedEvent(log *types.Log) (*ProcessorEndpointRoleGranted, error) {
	event := "RoleGranted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRoleGranted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointRoleRevoked represents a RoleRevoked event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRoleRevoked) ContractEventName() string {
	return ProcessorEndpointRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (processorEndpoint *ProcessorEndpoint) UnpackRoleRevokedEvent(log *types.Log) (*ProcessorEndpointRoleRevoked, error) {
	event := "RoleRevoked"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointRoleRevoked)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointStateRootUpdate represents a StateRootUpdate event raised by the ProcessorEndpoint contract.
type ProcessorEndpointStateRootUpdate struct {
	ApplicationId uint64
	RequestId     [32]byte
	OldStateRoot  [32]byte
	NewStateRoot  [32]byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointStateRootUpdateEventName = "StateRootUpdate"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointStateRootUpdate) ContractEventName() string {
	return ProcessorEndpointStateRootUpdateEventName
}

// UnpackStateRootUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StateRootUpdate(uint64 indexed applicationId, bytes32 indexed requestId, bytes32 oldStateRoot, bytes32 newStateRoot)
func (processorEndpoint *ProcessorEndpoint) UnpackStateRootUpdateEvent(log *types.Log) (*ProcessorEndpointStateRootUpdate, error) {
	event := "StateRootUpdate"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointStateRootUpdate)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointTriggerExecuted represents a TriggerExecuted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointTriggerExecuted struct {
	ApplicationId      uint64
	ProcessedRequestId [32]byte
	TriggerContract    common.Address
	Success            bool
	Raw                *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointTriggerExecutedEventName = "TriggerExecuted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointTriggerExecuted) ContractEventName() string {
	return ProcessorEndpointTriggerExecutedEventName
}

// UnpackTriggerExecutedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TriggerExecuted(uint64 indexed applicationId, bytes32 indexed processedRequestId, address indexed triggerContract, bool success)
func (processorEndpoint *ProcessorEndpoint) UnpackTriggerExecutedEvent(log *types.Log) (*ProcessorEndpointTriggerExecuted, error) {
	event := "TriggerExecuted"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointTriggerExecuted)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointTriggerPostWithdraw represents a TriggerPostWithdraw event raised by the ProcessorEndpoint contract.
type ProcessorEndpointTriggerPostWithdraw struct {
	ApplicationId      uint64
	ProcessedRequestId [32]byte
	TriggerContract    common.Address
	Success            bool
	Raw                *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointTriggerPostWithdrawEventName = "TriggerPostWithdraw"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointTriggerPostWithdraw) ContractEventName() string {
	return ProcessorEndpointTriggerPostWithdrawEventName
}

// UnpackTriggerPostWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TriggerPostWithdraw(uint64 indexed applicationId, bytes32 indexed processedRequestId, address indexed triggerContract, bool success)
func (processorEndpoint *ProcessorEndpoint) UnpackTriggerPostWithdrawEvent(log *types.Log) (*ProcessorEndpointTriggerPostWithdraw, error) {
	event := "TriggerPostWithdraw"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointTriggerPostWithdraw)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointUserEvent represents a UserEvent event raised by the ProcessorEndpoint contract.
type ProcessorEndpointUserEvent struct {
	ApplicationId uint64
	RequestId     [32]byte
	EventSubType  [32]byte
	EncryptedData []byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointUserEventEventName = "UserEvent"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointUserEvent) ContractEventName() string {
	return ProcessorEndpointUserEventEventName
}

// UnpackUserEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UserEvent(uint64 indexed applicationId, bytes32 indexed requestId, bytes32 indexed eventSubType, bytes encryptedData)
func (processorEndpoint *ProcessorEndpoint) UnpackUserEventEvent(log *types.Log) (*ProcessorEndpointUserEvent, error) {
	event := "UserEvent"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointUserEvent)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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

// ProcessorEndpointWithdrawal represents a Withdrawal event raised by the ProcessorEndpoint contract.
type ProcessorEndpointWithdrawal struct {
	ApplicationId uint64
	RequestId     [32]byte
	To            common.Address
	TokenAddress  common.Address
	Amount        *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointWithdrawalEventName = "Withdrawal"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointWithdrawal) ContractEventName() string {
	return ProcessorEndpointWithdrawalEventName
}

// UnpackWithdrawalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdrawal(uint64 indexed applicationId, bytes32 indexed requestId, address indexed to, address tokenAddress, uint256 amount)
func (processorEndpoint *ProcessorEndpoint) UnpackWithdrawalEvent(log *types.Log) (*ProcessorEndpointWithdrawal, error) {
	event := "Withdrawal"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointWithdrawal)
	if len(log.Data) > 0 {
		if err := processorEndpoint.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range processorEndpoint.abi.Events[event].Inputs {
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
func (processorEndpoint *ProcessorEndpoint) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AddressCantBeZero"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAddressCantBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AuthorityNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAuthorityNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["DeadlineExpired"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackDeadlineExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["DeployerNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackDeployerNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["FeeValueBelowMinimum"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackFeeValueBelowMinimumError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InsufficientAppBalance"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInsufficientAppBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InsufficientBalance"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidApplicationId"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidApplicationIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidPayload"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidPayloadError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidPermit"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidPermitError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidProtocolVersion"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidProtocolVersionError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidRequestId"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidRequestIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidRequestType"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidRequestTypeError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidShortString"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidShortStringError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidSigner"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidStateRoot"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidStateRootError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidValue"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidValueError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["MaxNumOfApplicationsExceeded"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackMaxNumOfApplicationsExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["NotRegisteredTrigger"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackNotRegisteredTriggerError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["QueueThresholdExceeded"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackQueueThresholdExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["StringTooLong"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackStringTooLongError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["TokenNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackTokenNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["TransferAmountMismatch"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackTransferAmountMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["TransferFailed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackTransferFailedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ProcessorEndpointAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func ProcessorEndpointAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (processorEndpoint *ProcessorEndpoint) UnpackAccessControlBadConfirmationError(raw []byte) (*ProcessorEndpointAccessControlBadConfirmation, error) {
	out := new(ProcessorEndpointAccessControlBadConfirmation)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func ProcessorEndpointAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (processorEndpoint *ProcessorEndpoint) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*ProcessorEndpointAccessControlUnauthorizedAccount, error) {
	out := new(ProcessorEndpointAccessControlUnauthorizedAccount)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointAddressCantBeZero represents a AddressCantBeZero error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAddressCantBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressCantBeZero()
func ProcessorEndpointAddressCantBeZeroErrorID() common.Hash {
	return common.HexToHash("0x4b054c82bb9e4ebe90744902747c4e86028dd2a122c63de8810d9a1c2e84617d")
}

// UnpackAddressCantBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressCantBeZero()
func (processorEndpoint *ProcessorEndpoint) UnpackAddressCantBeZeroError(raw []byte) (*ProcessorEndpointAddressCantBeZero, error) {
	out := new(ProcessorEndpointAddressCantBeZero)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AddressCantBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointAuthorityNotAllowed represents a AuthorityNotAllowed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAuthorityNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuthorityNotAllowed()
func ProcessorEndpointAuthorityNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xc4550a169e74f716a0e1f2ed847c9b836363900678799385b1658d47630ac161")
}

// UnpackAuthorityNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuthorityNotAllowed()
func (processorEndpoint *ProcessorEndpoint) UnpackAuthorityNotAllowedError(raw []byte) (*ProcessorEndpointAuthorityNotAllowed, error) {
	out := new(ProcessorEndpointAuthorityNotAllowed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AuthorityNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointDeadlineExpired represents a DeadlineExpired error raised by the ProcessorEndpoint contract.
type ProcessorEndpointDeadlineExpired struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DeadlineExpired()
func ProcessorEndpointDeadlineExpiredErrorID() common.Hash {
	return common.HexToHash("0x1ab7da6b04cf80fcc6534c1f63ce7b24fd7488ab40444b4e5271ccc1c7f64567")
}

// UnpackDeadlineExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DeadlineExpired()
func (processorEndpoint *ProcessorEndpoint) UnpackDeadlineExpiredError(raw []byte) (*ProcessorEndpointDeadlineExpired, error) {
	out := new(ProcessorEndpointDeadlineExpired)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "DeadlineExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointDeployerNotAllowed represents a DeployerNotAllowed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointDeployerNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DeployerNotAllowed()
func ProcessorEndpointDeployerNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xfc336c41b5f2764aa23a0642681230d7b417aea242ed1c10a52ccb7eacd8481f")
}

// UnpackDeployerNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DeployerNotAllowed()
func (processorEndpoint *ProcessorEndpoint) UnpackDeployerNotAllowedError(raw []byte) (*ProcessorEndpointDeployerNotAllowed, error) {
	out := new(ProcessorEndpointDeployerNotAllowed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "DeployerNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointECDSAInvalidSignature represents a ECDSAInvalidSignature error raised by the ProcessorEndpoint contract.
type ProcessorEndpointECDSAInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignature()
func ProcessorEndpointECDSAInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0xf645eedf0193584640b6b90cb9477e4c95b98636c148a891d4c0a146dc46e75a")
}

// UnpackECDSAInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignature()
func (processorEndpoint *ProcessorEndpoint) UnpackECDSAInvalidSignatureError(raw []byte) (*ProcessorEndpointECDSAInvalidSignature, error) {
	out := new(ProcessorEndpointECDSAInvalidSignature)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "ECDSAInvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointECDSAInvalidSignatureLength represents a ECDSAInvalidSignatureLength error raised by the ProcessorEndpoint contract.
type ProcessorEndpointECDSAInvalidSignatureLength struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func ProcessorEndpointECDSAInvalidSignatureLengthErrorID() common.Hash {
	return common.HexToHash("0xfce698f7e8e5342cd615f641317bc45fe7e1e4a8b0a14dd1383ff8dc9c41917f")
}

// UnpackECDSAInvalidSignatureLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func (processorEndpoint *ProcessorEndpoint) UnpackECDSAInvalidSignatureLengthError(raw []byte) (*ProcessorEndpointECDSAInvalidSignatureLength, error) {
	out := new(ProcessorEndpointECDSAInvalidSignatureLength)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointECDSAInvalidSignatureS represents a ECDSAInvalidSignatureS error raised by the ProcessorEndpoint contract.
type ProcessorEndpointECDSAInvalidSignatureS struct {
	S [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func ProcessorEndpointECDSAInvalidSignatureSErrorID() common.Hash {
	return common.HexToHash("0xd78bce0cccb935155ed6428d1c13e50b7f3550fd2b66b9fe266006fea4a5e1eb")
}

// UnpackECDSAInvalidSignatureSError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func (processorEndpoint *ProcessorEndpoint) UnpackECDSAInvalidSignatureSError(raw []byte) (*ProcessorEndpointECDSAInvalidSignatureS, error) {
	out := new(ProcessorEndpointECDSAInvalidSignatureS)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointFeeValueBelowMinimum represents a FeeValueBelowMinimum error raised by the ProcessorEndpoint contract.
type ProcessorEndpointFeeValueBelowMinimum struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FeeValueBelowMinimum()
func ProcessorEndpointFeeValueBelowMinimumErrorID() common.Hash {
	return common.HexToHash("0xc77f97ebda9854b3f6367f2d96830b6ea60a9ab7e12c67ede499fadeb3f9cef7")
}

// UnpackFeeValueBelowMinimumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FeeValueBelowMinimum()
func (processorEndpoint *ProcessorEndpoint) UnpackFeeValueBelowMinimumError(raw []byte) (*ProcessorEndpointFeeValueBelowMinimum, error) {
	out := new(ProcessorEndpointFeeValueBelowMinimum)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "FeeValueBelowMinimum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInsufficientAppBalance represents a InsufficientAppBalance error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInsufficientAppBalance struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientAppBalance()
func ProcessorEndpointInsufficientAppBalanceErrorID() common.Hash {
	return common.HexToHash("0x0d42ec4a38c8a840124f0313a348ac9e3767fe30452656f386bda86ce53859cd")
}

// UnpackInsufficientAppBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientAppBalance()
func (processorEndpoint *ProcessorEndpoint) UnpackInsufficientAppBalanceError(raw []byte) (*ProcessorEndpointInsufficientAppBalance, error) {
	out := new(ProcessorEndpointInsufficientAppBalance)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InsufficientAppBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInsufficientBalance represents a InsufficientBalance error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInsufficientBalance struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientBalance()
func ProcessorEndpointInsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xf4d678b8ce6b5157126b1484a53523762a93571537a7d5ae97d8014a44715c94")
}

// UnpackInsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientBalance()
func (processorEndpoint *ProcessorEndpoint) UnpackInsufficientBalanceError(raw []byte) (*ProcessorEndpointInsufficientBalance, error) {
	out := new(ProcessorEndpointInsufficientBalance)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidApplicationId represents a InvalidApplicationId error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidApplicationId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidApplicationId()
func ProcessorEndpointInvalidApplicationIdErrorID() common.Hash {
	return common.HexToHash("0x9ff9a83079daf92a88066c1b25d04540f443a8d245022db6e1e62beaa6d6bc5e")
}

// UnpackInvalidApplicationIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidApplicationId()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidApplicationIdError(raw []byte) (*ProcessorEndpointInvalidApplicationId, error) {
	out := new(ProcessorEndpointInvalidApplicationId)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidApplicationId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidPayload represents a InvalidPayload error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidPayload struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPayload()
func ProcessorEndpointInvalidPayloadErrorID() common.Hash {
	return common.HexToHash("0x7c6953f9f55fcfd713fe1d72e370e4757e17ae7187e1acdcc02b4b76b56a9a35")
}

// UnpackInvalidPayloadError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPayload()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidPayloadError(raw []byte) (*ProcessorEndpointInvalidPayload, error) {
	out := new(ProcessorEndpointInvalidPayload)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidPayload", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidPermit represents a InvalidPermit error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidPermit struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPermit()
func ProcessorEndpointInvalidPermitErrorID() common.Hash {
	return common.HexToHash("0xddafbaefe5cdbfc41219261b5d7468d1f16bc986eedbdef2ea223830eb9c5e3e")
}

// UnpackInvalidPermitError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPermit()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidPermitError(raw []byte) (*ProcessorEndpointInvalidPermit, error) {
	out := new(ProcessorEndpointInvalidPermit)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidPermit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidProtocolVersion represents a InvalidProtocolVersion error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidProtocolVersion struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidProtocolVersion()
func ProcessorEndpointInvalidProtocolVersionErrorID() common.Hash {
	return common.HexToHash("0x5428eae7f7dcfe744a90b56d360bbe512220bcf69b1e537f3d2d3b28cf2bf92d")
}

// UnpackInvalidProtocolVersionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidProtocolVersion()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidProtocolVersionError(raw []byte) (*ProcessorEndpointInvalidProtocolVersion, error) {
	out := new(ProcessorEndpointInvalidProtocolVersion)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidProtocolVersion", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidRequestId represents a InvalidRequestId error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidRequestId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRequestId()
func ProcessorEndpointInvalidRequestIdErrorID() common.Hash {
	return common.HexToHash("0xba0514c0d5ec80c22a6ebd3f2e6691e2f9ad3d1402978084361a7537d8546138")
}

// UnpackInvalidRequestIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRequestId()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidRequestIdError(raw []byte) (*ProcessorEndpointInvalidRequestId, error) {
	out := new(ProcessorEndpointInvalidRequestId)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidRequestId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidRequestType represents a InvalidRequestType error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidRequestType struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRequestType()
func ProcessorEndpointInvalidRequestTypeErrorID() common.Hash {
	return common.HexToHash("0x3b8b6c143e72b9be1ff9776689ac73649bd2f21a04954ef955d859dcfff7eb57")
}

// UnpackInvalidRequestTypeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRequestType()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidRequestTypeError(raw []byte) (*ProcessorEndpointInvalidRequestType, error) {
	out := new(ProcessorEndpointInvalidRequestType)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidRequestType", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidShortString represents a InvalidShortString error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidShortString struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidShortString()
func ProcessorEndpointInvalidShortStringErrorID() common.Hash {
	return common.HexToHash("0xb3512b0c6163e5f0bafab72bb631b9d58cd7a731b082f910338aa21c83d5c274")
}

// UnpackInvalidShortStringError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidShortString()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidShortStringError(raw []byte) (*ProcessorEndpointInvalidShortString, error) {
	out := new(ProcessorEndpointInvalidShortString)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidShortString", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidSignature represents a InvalidSignature error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature()
func ProcessorEndpointInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x8baa579fce362245063d36f11747a89dd489c54795634fc673cc0e0db51fedc5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidSignatureError(raw []byte) (*ProcessorEndpointInvalidSignature, error) {
	out := new(ProcessorEndpointInvalidSignature)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidSigner represents a InvalidSigner error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidSigner struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSigner()
func ProcessorEndpointInvalidSignerErrorID() common.Hash {
	return common.HexToHash("0x815e1d64efb74fbe314c20a2b8a2335d18bce12a19165e447fa36bcb35959528")
}

// UnpackInvalidSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSigner()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidSignerError(raw []byte) (*ProcessorEndpointInvalidSigner, error) {
	out := new(ProcessorEndpointInvalidSigner)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidStateRoot represents a InvalidStateRoot error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidStateRoot struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidStateRoot()
func ProcessorEndpointInvalidStateRootErrorID() common.Hash {
	return common.HexToHash("0xb6fac0303f809aaca0354efa81577129a51cc71ce625abce5d9ded4cf63ebabc")
}

// UnpackInvalidStateRootError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidStateRoot()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidStateRootError(raw []byte) (*ProcessorEndpointInvalidStateRoot, error) {
	out := new(ProcessorEndpointInvalidStateRoot)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidStateRoot", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointInvalidValue represents a InvalidValue error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidValue struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidValue()
func ProcessorEndpointInvalidValueErrorID() common.Hash {
	return common.HexToHash("0xaa7feadc1a228aee3d4475b95bb4402ebafd3841876f41c4469039f5ee2998d5")
}

// UnpackInvalidValueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidValue()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidValueError(raw []byte) (*ProcessorEndpointInvalidValue, error) {
	out := new(ProcessorEndpointInvalidValue)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidValue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointMaxNumOfApplicationsExceeded represents a MaxNumOfApplicationsExceeded error raised by the ProcessorEndpoint contract.
type ProcessorEndpointMaxNumOfApplicationsExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error MaxNumOfApplicationsExceeded()
func ProcessorEndpointMaxNumOfApplicationsExceededErrorID() common.Hash {
	return common.HexToHash("0x2d5c0cfe697fe71ea4a46b3137aee7a75e8d99fed238e7efea81025efbed81ae")
}

// UnpackMaxNumOfApplicationsExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error MaxNumOfApplicationsExceeded()
func (processorEndpoint *ProcessorEndpoint) UnpackMaxNumOfApplicationsExceededError(raw []byte) (*ProcessorEndpointMaxNumOfApplicationsExceeded, error) {
	out := new(ProcessorEndpointMaxNumOfApplicationsExceeded)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "MaxNumOfApplicationsExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointNotRegisteredTrigger represents a NotRegisteredTrigger error raised by the ProcessorEndpoint contract.
type ProcessorEndpointNotRegisteredTrigger struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotRegisteredTrigger()
func ProcessorEndpointNotRegisteredTriggerErrorID() common.Hash {
	return common.HexToHash("0x578f5841e568b13e125b5cbc66eb999e9534d0cabe53e082406913fe3bc52846")
}

// UnpackNotRegisteredTriggerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotRegisteredTrigger()
func (processorEndpoint *ProcessorEndpoint) UnpackNotRegisteredTriggerError(raw []byte) (*ProcessorEndpointNotRegisteredTrigger, error) {
	out := new(ProcessorEndpointNotRegisteredTrigger)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "NotRegisteredTrigger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointQueueThresholdExceeded represents a QueueThresholdExceeded error raised by the ProcessorEndpoint contract.
type ProcessorEndpointQueueThresholdExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QueueThresholdExceeded()
func ProcessorEndpointQueueThresholdExceededErrorID() common.Hash {
	return common.HexToHash("0xd8219b3da97d7adbef217ce9f220bbf38600f597fb984fd0620ea7df76ea17f3")
}

// UnpackQueueThresholdExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QueueThresholdExceeded()
func (processorEndpoint *ProcessorEndpoint) UnpackQueueThresholdExceededError(raw []byte) (*ProcessorEndpointQueueThresholdExceeded, error) {
	out := new(ProcessorEndpointQueueThresholdExceeded)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "QueueThresholdExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the ProcessorEndpoint contract.
type ProcessorEndpointReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func ProcessorEndpointReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (processorEndpoint *ProcessorEndpoint) UnpackReentrancyGuardReentrantCallError(raw []byte) (*ProcessorEndpointReentrancyGuardReentrantCall, error) {
	out := new(ProcessorEndpointReentrancyGuardReentrantCall)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the ProcessorEndpoint contract.
type ProcessorEndpointSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func ProcessorEndpointSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (processorEndpoint *ProcessorEndpoint) UnpackSafeERC20FailedOperationError(raw []byte) (*ProcessorEndpointSafeERC20FailedOperation, error) {
	out := new(ProcessorEndpointSafeERC20FailedOperation)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointStringTooLong represents a StringTooLong error raised by the ProcessorEndpoint contract.
type ProcessorEndpointStringTooLong struct {
	Str string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringTooLong(string str)
func ProcessorEndpointStringTooLongErrorID() common.Hash {
	return common.HexToHash("0x305a27a93f8e33b7392df0a0f91d6fc63847395853c45991eec52dbf24d72381")
}

// UnpackStringTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringTooLong(string str)
func (processorEndpoint *ProcessorEndpoint) UnpackStringTooLongError(raw []byte) (*ProcessorEndpointStringTooLong, error) {
	out := new(ProcessorEndpointStringTooLong)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "StringTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointTokenNotAllowed represents a TokenNotAllowed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointTokenNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenNotAllowed()
func ProcessorEndpointTokenNotAllowedErrorID() common.Hash {
	return common.HexToHash("0xa29c498656b92b53bef33b5e62c8bb31470e81c408996e42696a26083f02fc8d")
}

// UnpackTokenNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenNotAllowed()
func (processorEndpoint *ProcessorEndpoint) UnpackTokenNotAllowedError(raw []byte) (*ProcessorEndpointTokenNotAllowed, error) {
	out := new(ProcessorEndpointTokenNotAllowed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "TokenNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointTransferAmountMismatch represents a TransferAmountMismatch error raised by the ProcessorEndpoint contract.
type ProcessorEndpointTransferAmountMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TransferAmountMismatch()
func ProcessorEndpointTransferAmountMismatchErrorID() common.Hash {
	return common.HexToHash("0xb0a6ea29805702cf1e83dfa66f5843be54f820580489524674c9eb97deec9d9f")
}

// UnpackTransferAmountMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TransferAmountMismatch()
func (processorEndpoint *ProcessorEndpoint) UnpackTransferAmountMismatchError(raw []byte) (*ProcessorEndpointTransferAmountMismatch, error) {
	out := new(ProcessorEndpointTransferAmountMismatch)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "TransferAmountMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointTransferFailed represents a TransferFailed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointTransferFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TransferFailed()
func ProcessorEndpointTransferFailedErrorID() common.Hash {
	return common.HexToHash("0x90b8ec1877afffd816d05d9b13947f3ff18ec5851c38bad15ec2b710f92391b1")
}

// UnpackTransferFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TransferFailed()
func (processorEndpoint *ProcessorEndpoint) UnpackTransferFailedError(raw []byte) (*ProcessorEndpointTransferFailed, error) {
	out := new(ProcessorEndpointTransferFailed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "TransferFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}
