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

// StructsTokenAndAmount is an auto generated low-level Go binding around an user-defined struct.
type StructsTokenAndAmount struct {
	Token  common.Address
	Amount *big.Int
}

// StructsWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsWithdrawalRequest struct {
	TokenAddress common.Address
	Receiver     common.Address
	Amount       *big.Int
}

// ProcessorEndpointMetaData contains all meta data concerning the ProcessorEndpoint contract.
var ProcessorEndpointMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resetOperator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"},{\"internalType\":\"contractITokenAllowlist\",\"name\":\"_tokenAllowlist\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeployerNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientAppBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPermit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxNumOfApplicationsExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotRegisteredTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerCannotBeEOA\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"AppEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"DeployRequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"DeployRequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMax\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"MaxNumberOfAppUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"ReportGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"withdrawSuccess\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"postWithdrawSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"returnedTokens\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"failedTokens\",\"type\":\"tuple[]\"}],\"name\":\"TriggerWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPLOYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_OPERATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"addAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminReset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"appIds\",\"type\":\"uint64[]\"}],\"name\":\"adminResetApps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"appCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"applicationStateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableDeploySlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"facilitatorNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeployedAppIds\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getFacilitatorNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getPendingRequestsPage\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTriggerQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"isAllowedDeployer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumOfApplications\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"removeAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitDeployRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"requestSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"depositPermit\",\"type\":\"bytes\"}],\"name\":\"submitRequestFor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitTriggerRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalAppCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalPendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"triggerContracts\",\"outputs\":[{\"internalType\":\"contractITrigger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"triggersToAppIds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"updateMaxNumOfApplications\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "ProcessorEndpoint",
	Bin: "0x610160806040523461034b5760e0816160da8038038091610020828561034f565b83398101031261034b5780516001600160a01b0381169081900361034b5760208201516001600160a01b038116929083900361034b5761006260408201610386565b9261006f60608301610386565b61007b60808401610386565b9360c060a08501519401519560018060a01b03871680970361034b5760409687516100a6898261034f565b60048152602081016356656c6160e01b81525f906100c4600161039a565b926100d18c51948561034f565b600184526100df600161039a565b602085019390601f1901368537600a602186015b5f1901916f181899199a1a9b1b9c1cb0b131b232b360811b8282061a835304801561012157600a90916100f3565b505060018055610130816105e0565b6101205261013d8461077b565b61014052519020918260e05251902080610100524660a05289519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f84528b83015260608201524660808201523060a082015260a081526101a460c08261034f565b5190206080523060c052600a600655600a600755600a600c5582158015610343575b8015610332575b8015610321575b8015610319575b61030a57600d80546001600160a01b03199081169094179055600e80548416909517909455600f8054831694909417909355601580549091166001600160a01b0384161790556102979161022e906103b5565b506102388161044d565b505f51602061607a5f395f51905f525f818152602081905286812060010180545f51602061609a5f395f51905f5291829055909290917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9080a46104cd565b506014556001600160a01b0381166102fa575b505161574690816108b4823960805181615405015260a051816154bc015260c051816153cf015260e051816154540152610100518161547a015261012051816113b7015261014051816113e00152f35b6103039061054d565b505f6102aa565b632582a64160e11b5f5260045ffd5b5080156101db565b506001600160a01b038416156101d4565b506001600160a01b038216156101cd565b5084156101c6565b5f80fd5b601f909101601f19168101906001600160401b0382119082101761037257604052565b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b038216820361034b57565b6001600160401b03811161037257601f01601f191660200190565b6001600160a01b0381165f9081525f51602061601a5f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f51602061601a5f395f51905f5260205260408120805460ff191660011790553391907f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd98905f516020615ffa5f395f51905f529080a4600190565b505f90565b6001600160a01b0381165f9081525f51602061603a5f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f51602061603a5f395f51905f5260205260408120805460ff191660011790553391905f51602061609a5f395f51905f52905f516020615ffa5f395f51905f529080a4600190565b6001600160a01b0381165f9081525f5160206160ba5f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f5160206160ba5f395f51905f5260205260408120805460ff191660011790553391905f51602061607a5f395f51905f52905f516020615ffa5f395f51905f529080a4600190565b6001600160a01b0381165f9081525f51602061605a5f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f51602061605a5f395f51905f5260205260408120805460ff191660011790553391907fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c95917905f516020615ffa5f395f51905f529080a4600190565b908151602081105f1461065a575090601f81511161061a57602081519101516020821061060b571790565b5f198260200360031b1b161790565b604460209160405192839163305a27a960e01b83528160048401528051918291826024860152018484015e5f828201840152601f01601f19168101030190fd5b6001600160401b03811161037257600254600181811c91168015610771575b602082101461075d57601f811161072a575b50602092601f82116001146106c957928192935f926106be575b50508160011b915f199060031b1c19161760025560ff90565b015190505f806106a5565b601f1982169360025f52805f20915f5b86811061071257508360019596106106fa575b505050811b0160025560ff90565b01515f1960f88460031b161c191690555f80806106ec565b919260206001819286850151815501940192016106d9565b60025f52601f60205f20910160051c810190601f830160051c015b818110610752575061068b565b5f8155600101610745565b634e487b7160e01b5f52602260045260245ffd5b90607f1690610679565b908151602081105f146107a6575090601f81511161061a57602081519101516020821061060b571790565b6001600160401b03811161037257600354600181811c911680156108a9575b602082101461075d57601f8111610876575b50602092601f821160011461081557928192935f9261080a575b50508160011b915f199060031b1c19161760035560ff90565b015190505f806107f1565b601f1982169360035f52805f20915f5b86811061085e5750836001959610610846575b505050811b0160035560ff90565b01515f1960f88460031b161c191690555f8080610838565b91926020600181928685015181550194019201610825565b60035f52601f60205f20910160051c810190601f830160051c015b81811061089e57506107d7565b5f8155600101610891565b90607f16906107c556fe6101e080604052600436101561001d575b50361561001b575f80fd5b005b5f905f3560e01c90816301ffc9a7146124a05750806305f6e1d0146121885780631256c918146121505780631664deeb146121165780631753520f146120f957806318733d57146120b8578063194b6be8146120805780631cc76c381461204057806320ff473f14611ff857806321c0b34214611fc6578063248a9ca314611f9c5780632a0acc6a14611f625780632b52b07c14611f285780632f2ff15d14611eeb5780632fbfa0d514611bec57806336568abe14611ba857806347101cae14611b8e5780634a0d0495146117cc5780635423112f146117855780635e3393841461176b578063655f2de1146116f05780636b288d20146116a25780636e91d1331461167a57806374202d9b146116065780637596ed43146115a15780637a36a8911461156957806380a1f712146114e7578063840059ec1461149757806384b0196e1461139f5780638537c2781461121357806386f4a07d146111d95780638c5b9b00146111ab5780638eb8791a1461118e57806391d14854146111465780639c73eeaf146110f45780639c8d416f146110d75780639fabc36f14611095578063a217fddf14611053578063a9ba045e1461106d578063aa3aa46014611053578063af20c9601461101b578063b7222ff514610fc8578063bcf7a3c9146108f3578063be1e3725146108c0578063c404c3591461087f578063c415b95c14610856578063d2c35ce8146107de578063d547741f1461079b578063d72e957a14610762578063dd39d76514610348578063e3d72e5e146102fb578063ecd00261146102c0578063f624f88a146102975763f7d9cc1e0361001057346102945780600319360112610294576020600c54604051908152f35b80fd5b5034610294578060031936011261029457600f546040516001600160a01b039091168152602090f35b503461029457806003193601126102945760206040517ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c8152f35b503461029457602036600319011261029457610315612520565b61031d6142a3565b6001600160a01b038116156103395761033590614778565b5080f35b632582a64160e11b8252600482fd5b5034610294576020366003190112610294576004356001600160401b03811161075e573660238201121561075e578060040135906001600160401b03821161075a5760248260051b82010190368211610715576103a3614312565b6103ab6144ec565b6103b3614fec565b6103bc83612b44565b916103ca6040519384612727565b8383526024602084019201915b818310610736575050509015610728575b600f5460405163024ece8960e01b8152908390829060049082906001600160a01b03165afa90811561071d578391610679575b50815191815190845b84810361043357856001805580f35b6001600160401b036104458284612dd0565b5116808752601260205260408720878052602052604087205480610612575b50865b84810361059157508087526004602052604087205461048a575b50600101610424565b8087526004602052866040812055600554875b8181036104b9575b505050600190816007540160075590610481565b826001600160401b036104cb83612e35565b90549060031b1c16146104e05760010161049d565b91505f19810190811161057d57906105136001600160401b0361050561053294612e35565b90549060031b1c1691612e35565b9091906001600160401b038084549260031b9316831b921b1916179055565b600554801561056957600191905f190161054b81612e35565b6001600160401b0382549160031b1b19169055600555905f806104a5565b634e487b7160e01b87526031600452602487fd5b634e487b7160e01b88526011600452602488fd5b818852601260205260408820600191906001600160a01b036105b3838a612dd0565b5116838060a01b03165f5260205260405f2054806105d3575b5001610467565b806105f461060c92338b6105ed87898060a01b0392612dd0565b51166151d9565b838060a01b03610604848b612dd0565b5116856149e9565b5f6105cc565b8780808084335af1610622614045565b501561066a5781885260126020526040882088805260205260408820610649828254612993565b9055878052601360205261066260408920918254612993565b90555f610464565b6312171d8360e31b8852600488fd5b90503d8084833e61068a8183612727565b810190602081830312610715578051906001600160401b03821161071957019080601f830112156107155781516106c081612b44565b926106ce6040519485612727565b81845260208085019260051b82010192831161071157602001905b8282106106f9575050505f61041b565b6020809161070684614031565b8152019101906106e9565b8580fd5b8380fd5b8480fd5b6040513d85823e3d90fd5b50610731612a3f565b6103e8565b82356001600160401b0381168103610756578152602092830192016103d7565b8680fd5b8280fd5b5080fd5b5034610294576020366003190112610294576020906040906001600160a01b0361078a612520565b168152601683522054604051908152f35b5034610294576040366003190112610294576103356004356107bb612536565b906107d96107d4825f525f602052600160405f20015490565b614381565b61481c565b5034610294576020366003190112610294576107f8612520565b6108006142a3565b6001600160a01b03168015610339576020817fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f926bffffffffffffffffffffffff60a01b6015541617601555604051908152a180f35b50346102945780600319360112610294576015546040516001600160a01b039091168152602090f35b50346102945780600319360112610294576108af61089b613f78565b6040519384936060855260608501906125e3565b916020840152151560408301520390f35b5034610294576040366003190112610294576108ef6108e36024356004356150e5565b60405191829182612691565b0390f35b50610140366003190112610e1f57610909612520565b6024359060ff82168203610e1f57604435906001600160401b0382168203610e1f5760056064351015610e1f576084356001600160401b038111610e1f576109559036906004016124f3565b909390919060a435906001600160a01b0382168203610e1f57610104356001600160401b038111610e1f5761098e9036906004016124f3565b90610124356001600160401b038111610e1f576109af9036906004016124f3565b91909260ff8516610fb9576001600160401b0389165f52600460205260405f205415610faa576109dd6144ec565b5f6003606435141580610f9b575b8181610f88575b50610f795760e4354211610f6a57610a0861489c565b600c541115610f5b57610f4757879187868b8d600360643514610f0f575b610b309593610b2193610b2796936001600160401b03610a5f60429560018060a01b0386165f52601660205260405f20549c3691612763565b602081519101209160ff6040519460208601967f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc1333340885260018060a01b0316604087015216606085015216608083015260ff6064351660a083015260c082015260018060a01b038d1660e082015260c4356101008201528861012082015260e4356101408201526101408152610af661016082612727565b519020610b016153cc565b906040519161190160f01b83526002830152602282015220923691612763565b906154e2565b9092919261551c565b6001600160a01b03168015610f00576001600160a01b03871603610ef15760018101809111610edd576001600160a01b0386165f908152601660205260409020556014543410610ece5760c435158080610ebc575b610ead5715610c86575b50509060ff610c399260209850610bf2610bb6600b5460c435868a8d8d606435908d6129c0565b98899760405195610bc6876126f0565b4287526001600160a01b03168c87015260c4356040870152346060870152608086018990523691612763565b60a08401526001600160a01b03851660c08401523360e08401526001600160401b03871661010084015216610120820152610c336064356101408301612799565b8361408a565b6040513381526001600160a01b03909116926001600160401b0316907fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236908690a460018055604051908152f35b6001600160a01b03841615610ead57600f5460405163cbe230c360e01b81526001600160a01b038681166004830152909160209183916024918391165afa908115610e14575f91610e73575b5015610e645760608103610e55578160609181010312610e1f5780359060ff8216809203610e1f57604051636eb1769f60e11b81526001600160a01b0386811660048301523060248301526020908290604490829089165afa908115610e14575f91610e23575b5060c43511610d72575b50509060ff610c3992602098610d5c60c43587866145f5565b610d6960c435858a614723565b98509192610b8f565b909291906001600160a01b0383163b15610e1f5760409081519463d505accf60e01b865260018060a01b038716600487015230602487015260c435604487015260e43560648701526084860152602081013560a4860152013560c48401525f8360e4818360018060a01b0387165af1908115610e1457602098610c399460ff93610e01575b5098509192610d43565b610e0d91505f90612727565b5f5f610df7565b6040513d5f823e3d90fd5b5f80fd5b90506020813d602011610e4d575b81610e3e60209383612727565b81010312610e1f57515f610d39565b3d9150610e31565b63ddafbaef60e01b5f5260045ffd5b63514e24c360e11b5f5260045ffd5b90506020813d602011610ea5575b81610e8e60209383612727565b81010312610e1f57610e9f906127a5565b5f610cd2565b3d9150610e81565b632a9ffab760e21b5f5260045ffd5b506001600160a01b0385161515610b85565b63c77f97eb60e01b5f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b632057875960e21b5f5260045ffd5b638baa579f60e01b5f5260045ffd5b50505050916085141580610f3c575b610f2d57879187868b8d610a26565b637c6953f960e01b5f5260045ffd5b5060e2881415610f1e565b634e487b7160e01b5f52602160045260245ffd5b63d8219b3d60e01b5f5260045ffd5b631ab7da6b60e01b5f5260045ffd5b630ee2db0560e21b5f5260045ffd5b9050610f475760046064351415816109f2565b50505f600160643514156109eb565b6309ff9a8360e41b5f5260045ffd5b635428eae760e01b5f5260045ffd5b34610e1f576040366003190112610e1f57610fe1612576565b6001600160401b03610ff1612536565b91165f52601260205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b34610e1f576020366003190112610e1f576001600160a01b0361103c612520565b165f526013602052602060405f2054604051908152f35b34610e1f575f366003190112610e1f5760206040515f8152f35b34610e1f575f366003190112610e1f57600e546040516001600160a01b039091168152602090f35b34610e1f576020366003190112610e1f5760206004356110b481614931565b9081156110c7575b506040519015158152f35b6110d19150614959565b826110bc565b34610e1f575f366003190112610e1f576020601454604051908152f35b34610e1f576020366003190112610e1f576004356111106142a3565b8015610ead576020817e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db192600c55604051908152a1005b34610e1f576040366003190112610e1f5761115f612536565b6004355f525f60205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b34610e1f575f366003190112610e1f576020600654604051908152f35b34610e1f575f366003190112610e1f576111c3614312565b6111cb6144ec565b6111d3614fec565b60018055005b34610e1f575f366003190112610e1f5760206040517fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c959178152f35b34610e1f57610180366003190112610e1f5761122d612576565b608435906001600160401b038211610e1f5760406003198336030112610e1f5760a4356001600160401b038111610e1f5760406003198236030112610e1f5760c435916001600160401b038311610e1f5736602384011215610e1f578260040135926001600160401b038411610e1f573660246060860283010111610e1f5761012435600d811015610e1f57610144356001600160401b038111610e1f576112d99036906004016124f3565b939092610164356001600160401b038111610e1f576112fc9036906004016124f3565b335f9081527e9b85fa08e135b4afefc8a40566229b01e66445186a3676642ffbacd3d93b0e6020526040902054909891979060ff1615611368576111d3996113426144ec565b6101043594602460e4359501926004019160040190606435906044359060243590612e75565b63e2517d3f60e01b5f52336004527f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9860245260445ffd5b34610e1f575f366003190112610e1f5761143b6113db7f000000000000000000000000000000000000000000000000000000000000000061521a565b6114047f0000000000000000000000000000000000000000000000000000000000000000615314565b6020611449604051926114178385612727565b5f84525f368137604051958695600f60f81b875260e08588015260e08701906125b2565b9085820360408701526125b2565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b82811061148057505050500390f35b835185528695509381019392810192600101611471565b34610e1f576040366003190112610e1f576114b0612520565b6114b8612536565b6001600160a01b039182165f908152601060209081526040808320949093168252928352819020549051908152f35b34610e1f575f366003190112610e1f5761150761150261489c565b6148e2565b600a54600b54905f905b82811061152657604051806108ef8682612691565b60018181925f52600960205260405f20545f52600860205261154a60405f206128e9565b6115548588612dd0565b5261155f8487612dd0565b5001910190611511565b34610e1f576020366003190112610e1f576001600160401b0361158a612576565b165f526004602052602060405f2054604051908152f35b34610e1f575f366003190112610e1f576115b9612a3f565b6040518091602082016020835281518091526020604084019201905f5b8181106115e4575050500390f35b82516001600160401b03168452859450602093840193909201916001016115d6565b34610e1f5760e0366003190112610e1f5761161f612520565b61162761258c565b90604435916005831015610e1f57606435916001600160401b038311610e1f5760209361165b6116729436906004016124f3565b9061166461254c565b9260c4359560a435956129c0565b604051908152f35b34610e1f575f366003190112610e1f57600d546040516001600160a01b039091168152602090f35b34610e1f576020366003190112610e1f576116bb612520565b6001600160a01b03165f9081525f5160206156f15f395f51905f52602090815260409182902054915160ff9092161515825290f35b34610e1f576020366003190112610e1f5760043561170c6142a3565b8015610ead5760065461172160075482612993565b90818310610ead578261175b6040937f80af1d8add7b822532307ae28ef73c31ab2803f908c5c22cc879766f54c685839560065582612993565b60075582519182526020820152a1005b34610e1f575f366003190112610e1f57602061167261489c565b34610e1f576020366003190112610e1f5761179e6127bf565b506004355f5260086020526108ef6117b860405f206128e9565b6040519182916020835260208301906125e3565b6040366003190112610e1f576117e06125a2565b602435906001600160401b038211610e1f5761180260ff9236906004016124f3565b929091169182610fb9576118146144ec565b335f9081525f5160206156f15f395f51905f52602052604090205460ff1615611b7f576007548015611b705761184861489c565b600c541115610f5b576014543410610ece575f1901600755600b546040805133602082019081525f928201929092529091816118ba915f606083015260e0608083015261189a610100830187896129a0565b905f60a08401525f60c084015260e083015203601f198101835282612727565b519020918260c01c91604051906118d0826126f0565b42825260208201915f835260408101915f83526118fd606083019534875260808401928984523691612763565b60a083019081523360c084019081525f60e0850181815261010086018a815261012087019c8d5261014087018381528c8452600860205260409093209651875597516001870180546001600160a01b0319166001600160a01b03929092169190911790559551600286015596516003850155915160048401555180519395949360058401916001600160401b038211611b5c5761199a8354612811565b601f8111611b21575b50602090601f8311600114611ab95791806119d9926007979695945f92611aae575b50508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b03938416179091559551929091018054909516911617808455905195519151959160e01b60ff60e01b16906005871015610f475760209660ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600b545f52600983528160405f20556001600b5401600b55604051908282527fbbcfb0c180d5f95cfb50c12005c44144ac92ac8474f1abc5add0e6e69a14e8dc843393a360018055604051908152f35b015190508d806119c5565b90601f19831691845f52815f20925f5b818110611b0957509160019391856007999897969410611af1575b505050811b0190556119dc565b01515f1960f88460031b161c191690558c8080611ae4565b92936020600181928786015181550195019301611ac9565b611b4c90845f5260205f20601f850160051c81019160208610611b52575b601f0160051c0190614074565b8b6119a3565b9091508190611b3f565b634e487b7160e01b5f52604160045260245ffd5b6316ae067f60e11b5f5260045ffd5b63fc336c4160e01b5f5260045ffd5b34610e1f575f366003190112610e1f5760206116726148b8565b34610e1f576040366003190112610e1f57611bc1612536565b336001600160a01b03821603611bdd5761001b9060043561481c565b63334bd91960e11b5f5260045ffd5b60e0366003190112610e1f57611c006125a2565b611c0861258c565b604435906005821015610e1f576064356001600160401b038111610e1f57611c349036906004016124f3565b9390611c3e61254c565b60a4359460ff60c435941680610fb9576001600160401b03861695865f52600460205260405f205415610faa57611c736144ec565b8215610f79576014548610610ece57611c8a61489c565b600c541115610f5b576001600160a01b0384169485611e5357611cad878a6127b2565b3403610ead575b60038403611d8f5760858a141580611d84575b610f2d57602099611d4b97611cf7611d23948c610c33995b611cea828285614723565b85878b600b5495336129c0565b9a60405198611d058a6126f0565b428a528d8a0152604089015260608801528960808801523691612763565b60a08501523360c08501525f60e0850152856101008501526101208401526101408301612799565b81604051915f83527fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236853394a460018055604051908152f35b5060e28a1415611cc7565b929495919360028614611db9575b92602099610c339593611cf7611d23948c611d4b9b9a98611cdf565b93919594926044602060018060a01b03600e5416604051928380926308b56a5360e11b82528d60048301523360248301525afa908115610e14575f91611e19575b5015611e0a579294959193611d9d565b63622a850b60e11b5f5260045ffd5b90506020813d602011611e4b575b81611e3460209383612727565b81010312610e1f57611e45906127a5565b8b611dfa565b3d9150611e27565b863403610ead578815610ead57600f5460405163cbe230c360e01b81526004810188905290602090829060249082906001600160a01b03165afa908115610e14575f91611eb1575b5015610e6457611eac8933876145f5565b611cb4565b90506020813d602011611ee3575b81611ecc60209383612727565b81010312610e1f57611edd906127a5565b8b611e9b565b3d9150611ebf565b34610e1f576040366003190112610e1f5761001b600435611f0a612536565b90611f236107d4825f525f602052600160405f20015490565b614464565b34610e1f575f366003190112610e1f5760206040517f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc13333408152f35b34610e1f575f366003190112610e1f5760206040517fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec428152f35b34610e1f576020366003190112610e1f5760206116726004355f525f602052600160405f20015490565b34610e1f576040366003190112610e1f576111d3611fe2612520565b611fea612536565b90611ff36144ec565b61450c565b34610e1f576020366003190112610e1f57612011612520565b6120196142a3565b6001600160a01b038116156120315761001b906143b9565b632582a64160e11b5f5260045ffd5b34610e1f576020366003190112610e1f576001600160401b03612061612576565b165f526017602052602060018060a01b0360405f205416604051908152f35b34610e1f576020366003190112610e1f576001600160a01b036120a1612520565b165f526011602052602060405f2054604051908152f35b34610e1f576020366003190112610e1f576001600160a01b036120d9612520565b165f52601860205260206001600160401b0360405f205416604051908152f35b34610e1f575f366003190112610e1f576020600754604051908152f35b34610e1f575f366003190112610e1f5760206040517f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd988152f35b34610e1f576020366003190112610e1f576001600160a01b03612171612520565b165f526016602052602060405f2054604051908152f35b34610e1f576020366003190112610e1f576004356001600160401b038111610e1f576121b89036906004016124f3565b335f5260186020526001600160401b0360405f205416801561249157601c5460408051336020820190815291810184905290918161220d916004606083015260e0608083015261189a6101008301888a6129a0565b519020916040519361221e856126f0565b42855260208501905f825260408601955f875261224b60608201925f845260808301958887523691612763565b60a082019081523360c083019081525f60e0840181815261010085018981526101208601838152600461014088018181528d8652601960205260409095209751885598516001880180546001600160a01b0319166001600160a01b03929092169190911790559b516002870155955160038601559651958401959095559051805160058401916001600160401b038211611b5c576122e98354612811565b601f8111612461575b50602090601f83116001146123f9579180612327926007979695945f92611aae5750508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b03938416179091559551929091018054909516911617808455905195519151959160e01b60ff60e01b16906005871015610f475760209660ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055601c545f52601a83528160405f20556001601c5401601c5581604051915f83527fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236853394a4604051908152f35b90601f19831691845f52815f20925f5b81811061244957509160019391856007999897969410612431575b505050811b01905561232a565b01515f1960f88460031b161c191690558c8080612424565b92936020600181928786015181550195019301612409565b61248b90845f5260205f20601f850160051c81019160208610611b5257601f0160051c0190614074565b8b6122f2565b63578f584160e01b5f5260045ffd5b34610e1f576020366003190112610e1f576004359063ffffffff60e01b8216809203610e1f57602091637965db0b60e01b81149081156124e2575b5015158152f35b6301ffc9a760e01b149050836124db565b9181601f84011215610e1f578235916001600160401b038311610e1f5760208381860195010111610e1f57565b600435906001600160a01b0382168203610e1f57565b602435906001600160a01b0382168203610e1f57565b608435906001600160a01b0382168203610e1f57565b35906001600160a01b0382168203610e1f57565b600435906001600160401b0382168203610e1f57565b602435906001600160401b0382168203610e1f57565b6004359060ff82168203610e1f57565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b906005821015610f475752565b9061268e908251815260018060a01b0360208401511660208201526040830151604082015260608301516060820152608083015160808201526101408061263b60a086015161016060a08601526101608501906125b2565b9460018060a01b0360c08201511660c085015260018060a01b0360e08201511660e08501526001600160401b036101008201511661010085015260ff6101208201511661012085015201519101906125d6565b90565b602081016020825282518091526040820191602060408360051b8301019401925f915b8383106126c357505050505090565b90919293946020806126e1600193603f1986820301875289516125e3565b970193019301919392906126b4565b61016081019081106001600160401b03821117611b5c57604052565b604081019081106001600160401b03821117611b5c57604052565b90601f801991011681019081106001600160401b03821117611b5c57604052565b6001600160401b038111611b5c57601f01601f191660200190565b92919261276f82612748565b9161277d6040519384612727565b829481845281830111610e1f578281602093845f960137010152565b6005821015610f475752565b51908115158203610e1f57565b91908201809211610edd57565b604051906127cc826126f0565b5f61014083828152826020820152826040820152826060820152826080820152606060a08201528260c08201528260e082015282610100820152826101208201520152565b90600182811c9216801561283f575b602083101461282b57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612820565b9060405191825f82549261285c84612811565b80845293600181169081156128c75750600114612883575b5061288192500383612727565b565b90505f9291925260205f20905f915b8183106128ab575050906020612881928201015f612874565b6020919350806001915483858901015201910190918492612892565b90506020925061288194915060ff191682840152151560051b8201015f612874565b906128816040516128f9816126f0565b61014060ff600783968054855260018060a01b03600182015416602086015260028101546040860152600381015460608601526004810154608086015261294260058201612849565b60a086015260018060a01b0360068201541660c0860152015460018060a01b03811660e08501526001600160401b038160a01c16610100850152818160e01c1661012085015260e81c169101612799565b91908203918211610edd57565b908060209392818452848401375f828201840152601f01601f1916010190565b9691612a3995936001600160401b039792956129ff612a11936040519a8b9960208b01809e60018060a01b031690521660408a015260608901906125d6565b60e060808801526101008701916129a0565b6001600160a01b0390931660a085015260c084015260e083015203601f198101835282612727565b51902090565b60405190600554808352826020810160055f5260205f20925f905b806003830110612af257612881945491818110612ad8575b818110612abb575b818110612a9e575b10612a90575b500383612727565b60c01c81526020015f612a88565b9260206001916001600160401b038560801c168152019301612a82565b9260206001916001600160401b038560401c168152019301612a7a565b9260206001916001600160401b0385168152019301612a72565b916004919350608060019186546001600160401b03811682526001600160401b038160401c1660208301526001600160401b0381841c16604083015260c01c6060820152019401920185929391612a5a565b6001600160401b038111611b5c5760051b60200190565b90612b6582612b44565b612b726040519182612727565b8281528092612b83601f1991612b44565b0190602036910137565b903590601e1981360301821215610e1f57018035906001600160401b038211610e1f57602001918160051b36038313610e1f57565b9190604083820312610e1f57604051612bda8161270c565b809380356001600160401b038111610e1f57810183601f82011215610e1f578035612c0481612b44565b91612c126040519384612727565b81835260208084019260051b82010190868211610e1f5760208101925b828410612cb1575050505082526020810135906001600160401b038211610e1f57019180601f84011215610e1f578235612c6881612b44565b93612c766040519586612727565b81855260208086019260051b820101928311610e1f57602001905b828210612ca15750505060200152565b8135815260209182019101612c91565b83356001600160401b038111610e1f57820188603f82011215610e1f57602091612ce48a83604086809601359101612763565b815201930192612c2f565b60408201908051916040845282518091526060840190602060608260051b8701019401915f905b828210612d60575050505060200151916020818303910152602080835192838152019201905f5b818110612d4a5750505090565b8251845260209384019390920191600101612d3d565b90919294602080612d7d600193605f198b820301865289516125b2565b970192019201909291612d16565b90600d821015610f475752565b9190811015612da8576060020190565b634e487b7160e01b5f52603260045260245ffd5b356001600160a01b0381168103610e1f5790565b8051821015612da85760209160051b010190565b9190811015612da85760051b0190565b9190811015612da85760051b81013590601e1981360301821215610e1f5701908135916001600160401b038311610e1f576020018236038113610e1f579190565b90600554821015612da85760055f52600282901c7f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0019160031b60181690565b9d9c9b9897969594939a9d6101605260e052610180526101a05261012052610100526101405260a05260c052612ead6101a051614931565b8015613f66575b15613f5757612ec56101a051614931565b5f6101c052908115613f3f576101a0515f52601960205260405f206101c0525b60076101c05101546080526001600160401b0360805160a01c166001600160401b03610160511603610faa576001600160401b0361016051165f52600460205260405f205460e05103613b2357612f3c8480612b8d565b959050612f4c6020860186612b8d565b90508603610f2d57612f646101205161012051612b8d565b979050612f7a6020610120510161012051612b8d565b90508803610f2d5760405191612f8f836126f0565b6001600160401b036101605116835260e05160208401526101805160408401526101a0516060840152612fc23688612bc2565b6080840152612fd43661012051612bc2565b60a0840152612fe561014051612b44565b612ff26040519182612727565b610140518152602081013660606101405102610100510111610e1f5761010051905b6060610140510261010051018210613ee257505060c084015260a05160e084015260c051610100840152600d861015610f475790918561012082015261305b36858c612763565b61014082015260018060a01b03600d541691604051938492635c9626b960e01b8452604060048501526001600160401b0381511660448501526020810151606485015260408101516084850152606081015160a48501526130e66130d0608083015161016060c48801526101a4870190612cef565b60a08301518682036043190160e4880152612cef565b9060c08101519160431986820301610104870152602080845192838152019301905f5b818110613ea25750505060209593859361316485946101408560e06131769701516101248901526101008101516101448901526131506101208201516101648a0190612d8b565b0151868203604319016101848801526125b2565b848103600319016024860152916129a0565b03915afa908115610e14575f91613e68575b5015610f00576101c05160038101805460069092015460805191999295916001600160a01b0391821691168015613e6057995b82613b325750505050506001600160401b0361016051165f5260046020526101805160405f205414613b235715613aff575b505f5f906131fd61014051612b5b565b61320961014051612b5b565b905f925b61014051811061397b57505f5b8381106138c1575050505061323c9061323760c05160a0516127b2565b6127b2565b613244614981565b106138b2575f5b82810361383d575050505f5b8181036137b057505060ff60076101c051015460e81c169060058210918215610f475760028114613776575b6001600160401b0361016051165f5260046020526101805160405f2055159081613647575b60405160e05181526101805160208201526101a051907f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e60406001600160401b03610160511692a360a0516135f6575b505f5f6001600160401b0361016051165f52601760205260018060a01b0360405f20541661332861014051612b5b565b925b61014051811061349457505061333f81612b5b565b915f5b82810361346c5750505061336190610120516101a05161016051614ceb565b60405191613370602084612727565b5f835261337b6148b8565b15158061345a575b1561344d57613390615611565b15610f4757156133fe57604051907f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d6101a05192806133e16001600160401b036101605116945f8060c05185614adf565b0390a35b60c05160155461288191906001600160a01b0316614a3a565b604051907f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed666101a05192806134456001600160401b036101605116945f8060c05185614adf565b0390a36133e5565b6134556155d4565b613390565b506134676101a051614931565b613383565b6001906001600160a01b036134818285612dd0565b511661348d8287612dd0565b5201613342565b8160018060a01b036134b860206134b2856101405161010051612d98565b01612dbc565b16146135c0575b806135186134de6134d96001946101405161010051612d98565b612dbc565b6134f460206134b2856101405161010051612d98565b906040613508856101405161010051612d98565b013591858060a01b031690614a95565b61352e60206134b2836101405161010051612d98565b6135426134d9836101405161010051612d98565b6040613555846101405161010051612d98565b0135907faafeb614695d37710d00f37c69671d205eb5b6fcd02995f6c3c7e2c93d3cbb2460405193868060a01b031693806135b76101a051956001600160401b036101605116958360209093929193604081019460018060a01b031681520152565b0390a40161332a565b91600180916135d96134d9866101405161010051612d98565b6135e38288612dd0565b90838060a01b03169052019290506134bf565b60018060a01b031661360a60a05182614a3a565b604051604081015f825260a05160208301525f5160206156d15f395f51905f526101a05192806001600160401b036101605116930390a45f6132f8565b60055468010000000000000000811015611b5c5761366e81600161368f9301600555612e35565b61016051906001600160401b038084549260031b9316831b921b1916179055565b60056101c0510160206136a28254612811565b146136ae575b506132a8565b6136b790612849565b602081805181010312610e1f57602001516001600160a01b03811690819003610e1f5780156136a857805f5260186020526001600160401b0360405f20541661376757803b1561375857610160516001600160401b03165f81815260176020908152604080832080546001600160a01b0319166001600160a01b03871617905593825260189052918220805467ffffffffffffffff191690911790556136a8565b634ea29e6760e01b5f5260045ffd5b63d48c379760e01b5f5260045ffd5b6101a0516001600160401b0361016051167ff36f013d86705b333febcccf7ece97eee30f4cd398a9c69aee1f940ed5daddea5f80a3613283565b806137d06001926137ca6020610120510161012051612b8d565b90612de4565b356137eb826137e56101205161012051612b8d565b90612df4565b7f485d4bfb37695319d1b2352eef5982d10db8c284b0a7ddefe4465377c9fbc8df60405160208152806138346101a051956001600160401b0361016051169560208401916129a0565b0390a401613257565b806138516001926137ca6020860186612b8d565b35613860826137e58680612b8d565b7f9ede1ad14e46a641d7b174bb6c9d939ad961d1a7de4b669163b9714a67e2232060405160208152806138a96101a051956001600160401b0361016051169560208401916129a0565b0390a40161324b565b631e9acf1760e31b5f5260045ffd5b6001600160a01b036138d38284612dd0565b516040516370a0823160e01b81523060048201529116602082602481845afa918215610e14575f92613946575b5061392881613939925f52601360205260405f2054905f52601160205260405f2054906127b2565b6139328487612dd0565b51906127b2565b116138b25760010161321a565b9091506020813d8211613973575b8161396160209383612727565b81010312610e1f575190613928613900565b3d9150613954565b936139906134d9866101405161010051612d98565b9060406139a4876101405161010051612d98565b0135916001600160401b0361016051165f52601260205260405f2060018060a01b0382165f5260205260405f20548311613af057610160516001600160401b03165f9081526012602090815260408083206001600160a01b0385168452825280832080548790039055601390915290205483116138b2576001600160a01b0381165f81815260136020526040902080548590039055613a525750600191613a4a916127b2565b945b0161320d565b909591905f805b878110613aa2575b5015613a72575b5050600190613a4c565b91946001928392906001600160a01b0316613a8d8387612dd0565b52613a988287612dd0565b520193905f613a68565b6001600160a01b0383811690613ab88389612dd0565b511614613ac757600101613a59565b9050613ae7613ae084613ada848a612dd0565b516127b2565b9187612dd0565b5260015f613a61565b6306a1762560e11b5f5260045ffd5b613b0d60c05160a0516127b2565b03610ead5760145460c05110610ead575f6131ed565b630b6fac0360e41b5f5260045ffd5b94509550959791979690961590811591613e56575b508015613e4a575b610f2d576001600160401b0361016051165f5260046020526101805160405f205403613b235760026101c05101549160018060a01b0360016101c051015416936001600160401b0361016051165f52601260205260405f2060018060a01b0386165f5260205260405f20548411613af057613bc8614981565b106138b257613be990613bdf8486610160516149e9565b5460145490612993565b9280613ddf57508115613d7d575090613c01916127b2565b613c0b8183614a3a565b6040519060408201905f835260208301525f5160206156d15f395f51905f526101a05192806001600160401b036101605116930390a45b60ff60076101c051015460e81c166005811015610f475715613d6f575b613c7c6014549260ff60076101c051015460e81c16943691612763565b90613c856148b8565b151580613d5d575b15613d5057613c9a615611565b6005841015610f475761288193613d04577f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d60405180613cef6101a051956001600160401b0361016051169560018985614adf565b0390a35b6015546001600160a01b0316614a3a565b7f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed6660405180613d486101a051956001600160401b0361016051169560018985614adf565b0390a3613cf3565b613d586155d4565b613c9a565b50613d6a6101a051614931565b613c8d565b600160075401600755613c5f565b9192505081613d8e575b5050613c42565b6001600160a01b031690613da28183614a3a565b6040519060408201905f835260208301525f5160206156d15f395f51905f526101a05192806001600160401b036101605116930390a45f80613d87565b8291939492613df8575b50505081613d8e575050613c42565b613e03828483614a95565b604080516101a051610160516001600160a01b0394909416825260208201949094526001600160401b03909216915f5160206156d15f395f51905f529190a45f8080613de9565b50610140511515613b4f565b905015155f613b47565b5080996131bb565b90506020813d602011613e9a575b81613e8360209383612727565b81010312610e1f57613e94906127a5565b5f613188565b3d9150613e76565b825180516001600160a01b03908116875260208281015190911681880152604091820151918701919091528a985060609095019490920191600101613109565b606082360312610e1f576040519060608201908282106001600160401b03831117611b5c57606092602092604052613f1985612562565b8152613f26838601612562565b8382015260408501356040820152815201910190613014565b6101a0515f52600860205260405f206101c052612ee5565b6302e8145360e61b5f5260045ffd5b50613f736101a051614959565b612eb4565b613f806127bf565b50613f896148b8565b613ff057613f9561489c565b613fa857613fa16127bf565b905f905f90565b600a545f52600960205260405f20545f52600860205260405f20906001600160401b03600783015460a01c165f526004602052613fe960405f2054926128e9565b9190600190565b601b545f52601a60205260405f20545f52601960205260405f20906001600160401b03600783015460a01c165f526004602052613fe960405f2054926128e9565b51906001600160a01b0382168203610e1f57565b3d1561406f573d9061405682612748565b916140646040519384612727565b82523d5f602084013e565b606090565b81811061407f575050565b5f8155600101614074565b90815f52600860205260405f208151815560018060a01b03602083015116600182019060018060a01b03166bffffffffffffffffffffffff60a01b8254161790556040820151600282015560608201516003820155608082015160048201556005810160a08301518051906001600160401b038211611b5c5761410d8354612811565b601f8111614273575b50602090601f831160011461420c57918061414a9260079695945f926142015750508160011b915f199060031b1c19161790565b90555b60c08301516006820180546001600160a01b03199081166001600160a01b039384161790915560e08086015194909301805490911693909116929092178083556101008401516101208501516101409095015193949093921b60ff60e01b16916005811015610f475760ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600b545f52600960205260405f20556001600b5401600b55565b015190505f806119c5565b90601f19831691845f52815f20925f5b81811061425b575091600193918560079897969410614243575b505050811b01905561414d565b01515f1960f88460031b161c191690555f8080614236565b9293602060018192878601518155019501930161421c565b61429d90845f5260205f20601f850160051c81019160208610611b5257601f0160051c0190614074565b5f614116565b335f9081527f5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b7602052604090205460ff16156142db57565b63e2517d3f60e01b5f52336004527fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4260245260445ffd5b335f9081527f3f65c445f65be3b772e454052bf8f1f970a2679db965b2efcdbc04d3dd490789602052604090205460ff161561434a57565b63e2517d3f60e01b5f52336004527fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c9591760245260445ffd5b5f8181526020818152604080832033845290915290205460ff16156143a35750565b63e2517d3f60e01b5f523360045260245260445ffd5b6001600160a01b0381165f9081525f5160206156f15f395f51905f52602052604090205460ff1661445f576001600160a01b03165f8181525f5160206156f15f395f51905f5260205260408120805460ff191660011790553391907ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b505f90565b5f818152602081815260408083206001600160a01b038616845290915290205460ff166144e6575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b6002600154146144fd576002600155565b633ee5aeb560e01b5f5260045ffd5b60018060a01b03811691825f52601060205260405f2060018060a01b0382165f5260205260405f20549182156145ef577f16d46b28f1a9377c2df78652a0a8e31568b99311665bcec24cbba72edd2e5c2c8392855f52601060205260405f2060018060a01b0382165f526020525f6040812055855f52601160205260405f20614596858254612993565b9055604080516001600160a01b0394851681526020810195909552921692839290a2826145e6575f809350809281925af16145cf614045565b50156145d757565b6312171d8360e31b5f5260045ffd5b612881926151d9565b50505050565b6040516370a0823160e01b81523060048201526001600160a01b039190911691602082602481865afa918215610e14575f926146ee575b50602492614673602092604051906323b872dd60e01b8583015260018060a01b0316868201523060448201528660648201526064815261466d608482612727565b8261557c565b6040516370a0823160e01b815230600482015293849182905afa8015610e14575f906146ba575b6146a49250612993565b036146ab57565b63b0a6ea2960e01b5f5260045ffd5b506020823d6020116146e6575b816146d460209383612727565b81010312610e1f576146a4915161469a565b3d91506146c7565b9091506020813d60201161471b575b8161470a60209383612727565b81010312610e1f575190602461462c565b3d91506146fd565b6001600160401b03165f52601260205260405f2060018060a01b0382165f5260205260405f206147548382546127b2565b905560018060a01b03165f52601360205261477460405f209182546127b2565b9055565b6001600160a01b0381165f9081525f5160206156f15f395f51905f52602052604090205460ff161561445f576001600160a01b03165f8181525f5160206156f15f395f51905f5260205260408120805460ff191690553391907ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16156144e6575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b600b54600a548082116148af5750505f90565b61268e91612993565b601c54601b548082116148af5750505f90565b600260038201549101548082116148af5750505f90565b906148ec82612b44565b6148f96040519182612727565b828152809261490a601f1991612b44565b01905f5b82811061491a57505050565b6020906149256127bf565b8282850101520161490e565b601c5490601b54809211918261494657505090565b9091505f52601a60205260405f20541490565b600b5490600a54809211918261496e57505090565b9091505f52600960205260405f20541490565b5f805260116020527f4ad3b33220dddc71b994a52d72c06b10862965f7d926534c05c00fb7e819e7b75461268e906149b99047612993565b5f805260136020527f8fa6efc3be94b5b348b21fea823fe8d100408cee9b7f90524494500445d8ff6c5490612993565b6001600160401b03165f52601260205260405f2060018060a01b0382165f5260205260405f20614a1a838254612993565b905560018060a01b03165f52601360205261477460405f20918254612993565b6001600160a01b03165f9081527f6e0956cda88cad152e89927e53611735b61a5c762d1428573c6931b0a5efcb01602052604090208054614a7c9083906127b2565b90555f8052601160205261477460405f209182546127b2565b60018060a01b031690815f52601060205260405f209060018060a01b03165f5260205260405f20614ac78382546127b2565b90555f52601160205261477460405f209182546127b2565b909392919381526002841015610f4757614b0860809261268e9560208401526040830190612d8b565b81606082015201906125b2565b9035601e1982360301811215610e1f5701602081359101916001600160401b038211610e1f578160051b36038313610e1f57565b906040810191614b598180614b15565b809460408552526060830160608560051b85010194825f90601e19813603015b838310614bc457505050505050806020614b94920190614b15565b82840360209093019290925281835290916001600160fb1b038311610e1f5760209260051b809284830137010190565b909192939497605f19888203018652883582811215610e1f57830190602082359201916001600160401b038111610e1f578036038313610e1f57614c0e60209283926001956129a0565b9a0196019493019190614b79565b81601f82011215610e1f57805190614c3382612b44565b92614c416040519485612727565b82845260208085019360061b83010191818311610e1f57602001925b828410614c6b575050505090565b604084830312610e1f5760206040918251614c858161270c565b614c8e87614031565b81528287015183820152815201930192614c5d565b90602080835192838152019201905f5b818110614cc05750505090565b825180516001600160a01b031685526020908101518186015260409094019390920191600101614cb3565b90925f906001600160401b03831693845f52601760205260018060a01b0360405f205416918215614f60578051905f5b828103614f3a575050506001823b15610e1f576040516338d9814d60e11b8152602060048201525f8180614d526024820187614b49565b038183885af19081614f25575b50614f1f575082935b86867fca435620e7fa2ff484abda442975d2167d5e6577548dfdfac799cecc2a2af2c9602060405198151598898152a3839460609285614dcb859660019383604051809781958294639cb86de760e01b8452604060048501526044840190614b49565b90602483015203925af1868188948993614ea0575b50614e915750505050835b908251915b828103614e5a575050507fa59edce59a4211e20e9f86e63adea41a55dbe409040bc17f473b671e31bd036a939291614e47614e55926040519586951515865215156020860152608060408601526080850190614ca3565b908382036060850152614ca3565b0390a3565b600190614e8b6001600160a01b03614e728388612dd0565b5151166020614e818489612dd0565b5101519085614723565b01614df0565b97509095909450909250614deb565b94509150503d8088853e614eb48185612727565b8301928581850312614f1b57614ec9816127a5565b9360208201516001600160401b038111614f175781614ee9918401614c1c565b916040810151906001600160401b038211614f1357614f09929101614c1c565b939093915f614de0565b8a80fd5b8980fd5b8780fd5b93614d68565b614f329195505f90612727565b5f935f614d5f565b600190614f5a866001600160a01b03614f538487612dd0565b511661450c565b01614d1b565b50505050505050565b60075f9182815582600182015582600282015582600382015582600482015560058101614f968154612811565b9081614fa9575b50508260068201550155565b81601f869311600114614fc05750555b5f80614f9d565b81835260208320614fdc91601f0160051c810190600101614074565b8082528160208120915555614fb9565b600a54600b545f905b80830361500e57509050600a54600b5560075401600755565b825f52600960205260405f20545f52600860205260405f2092600784015460ff8160e81c166005811015610f4757156150d9575b8460026001949596018054928361508e575b50505050805f52600960205260405f20545f52600860205261507860405f20614f69565b5f81815260096020526040812055019190614ff5565b6150b46150d094878501926001600160401b03898060a01b038554169160a01c166149e9565b5460069092015490549160a086901b8690039182169116614a95565b5f808080615054565b60019390930192615042565b9190916150f260086148cb565b928382108015906151d1575b6151955761510c90826127b2565b9280841161518d575b5061513a6151266115028386612993565b93615134600a5493846127b2565b926127b2565b905f905b82811061514a57505050565b60018181925f52600960205260405f20545f52600860205261516e60405f206128e9565b6151788589612dd0565b526151838488612dd0565b500191019061513e565b92505f615115565b50909150506040516151a8602082612727565b5f81525f805b8181106151ba57505090565b6020906151c56127bf565b828286010152016151ae565b5080156150fe565b60405163a9059cbb60e01b60208201526001600160a01b0392909216602483015260448083019390935291815261288191615215606483612727565b61557c565b60ff81146152605760ff811690601f8211615251576040519161523e604084612727565b6020808452838101919036833783525290565b632cd44ac360e21b5f5260045ffd5b50604051600254815f61527283612811565b80835292600181169081156152f55750600114615296575b61268e92500382612727565b5060025f90815290917f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace5b8183106152d957505090602061268e9282010161528a565b60209193508060019154838588010152019101909183926152c1565b6020925061268e94915060ff191682840152151560051b82010161528a565b60ff81146153385760ff811690601f8211615251576040519161523e604084612727565b50604051600354815f61534a83612811565b80835292600181169081156152f5575060011461536d5761268e92500382612727565b5060035f90815290917fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b5b8183106153b057505090602061268e9282010161528a565b6020919350806001915483858801015201910190918392615398565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614806154b9575b15615427577f000000000000000000000000000000000000000000000000000000000000000090565b60405160208101907f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f82527f000000000000000000000000000000000000000000000000000000000000000060408201527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a082015260a08152612a3960c082612727565b507f000000000000000000000000000000000000000000000000000000000000000046146153fe565b81519190604183036155125761550b9250602082015190606060408401519301515f1a9061564e565b9192909190565b50505f9160029190565b6004811015610f47578061552e575050565b600181036155455763f645eedf60e01b5f5260045ffd5b60028103615560575063fce698f760e01b5f5260045260245ffd5b60031461556a5750565b6335e2f38360e21b5f5260045260245ffd5b905f602091828151910182855af115610e14575f513d6155cb57506001600160a01b0381163b155b6155ab5750565b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b600114156155a4565b600a545f52600960205260405f20545f5260086020526155f660405f20614f69565b600a545f5260096020525f60408120556001600a5401600a55565b601b545f52601a60205260405f20545f52601960205261563360405f20614f69565b601b545f52601a6020525f60408120556001601b5401601b55565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a084116156c5579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15610e14575f516001600160a01b038116156156bb57905f905f90565b505f906001905f90565b5050505f916003919056fec9961aaccfc3406b40d0d25a5b83b8a4721e5f454116d7e22bd8038c715140aa740c5e3e456bed56f053f960110118ba9b95a1f5359a82283516fb2e81b6e37ea2646970667358221220535b5d4007b263757188ae0c26cf2cb443c0316e70db45e30a816caa5525426764736f6c634300081e00332f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d009b85fa08e135b4afefc8a40566229b01e66445186a3676642ffbacd3d93b0e5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b73f65c445f65be3b772e454052bf8f1f970a2679db965b2efcdbc04d3dd490789fc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184cdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec42740c5e3e456bed56f053f960110118ba9b95a1f5359a82283516fb2e81b6e37e",
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
// the contract method with ID 0x05f6e1d0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitTriggerRequest(bytes payload) returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitTriggerRequest(payload []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("submitTriggerRequest", payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitTriggerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x05f6e1d0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitTriggerRequest(bytes payload) returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitTriggerRequest(payload []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitTriggerRequest", payload)
}

// UnpackSubmitTriggerRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x05f6e1d0.
//
// Solidity: function submitTriggerRequest(bytes payload) returns(bytes32)
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
// Solidity: event TriggerExecuted(uint64 indexed applicationId, bytes32 indexed processedRequestId, bool success)
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

// ProcessorEndpointTriggerWithdraw represents a TriggerWithdraw event raised by the ProcessorEndpoint contract.
type ProcessorEndpointTriggerWithdraw struct {
	ApplicationId       uint64
	ProcessedRequestId  [32]byte
	WithdrawSuccess     bool
	PostWithdrawSuccess bool
	ReturnedTokens      []StructsTokenAndAmount
	FailedTokens        []StructsTokenAndAmount
	Raw                 *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointTriggerWithdrawEventName = "TriggerWithdraw"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointTriggerWithdraw) ContractEventName() string {
	return ProcessorEndpointTriggerWithdrawEventName
}

// UnpackTriggerWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TriggerWithdraw(uint64 indexed applicationId, bytes32 indexed processedRequestId, bool withdrawSuccess, bool postWithdrawSuccess, (address,uint256)[] returnedTokens, (address,uint256)[] failedTokens)
func (processorEndpoint *ProcessorEndpoint) UnpackTriggerWithdrawEvent(log *types.Log) (*ProcessorEndpointTriggerWithdraw, error) {
	event := "TriggerWithdraw"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointTriggerWithdraw)
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["TriggerAlreadyRegistered"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackTriggerAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["TriggerCannotBeEOA"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackTriggerCannotBeEOAError(raw[4:])
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

// ProcessorEndpointTriggerAlreadyRegistered represents a TriggerAlreadyRegistered error raised by the ProcessorEndpoint contract.
type ProcessorEndpointTriggerAlreadyRegistered struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TriggerAlreadyRegistered()
func ProcessorEndpointTriggerAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0xd48c37974b982f54734c76a70a50defa3e9afc19e571e138c59cb6f81e9a39ae")
}

// UnpackTriggerAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TriggerAlreadyRegistered()
func (processorEndpoint *ProcessorEndpoint) UnpackTriggerAlreadyRegisteredError(raw []byte) (*ProcessorEndpointTriggerAlreadyRegistered, error) {
	out := new(ProcessorEndpointTriggerAlreadyRegistered)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "TriggerAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointTriggerCannotBeEOA represents a TriggerCannotBeEOA error raised by the ProcessorEndpoint contract.
type ProcessorEndpointTriggerCannotBeEOA struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TriggerCannotBeEOA()
func ProcessorEndpointTriggerCannotBeEOAErrorID() common.Hash {
	return common.HexToHash("0x4ea29e672521f87a56ee460cd2be1b29334f535a27e686ee3144db2a8882eab3")
}

// UnpackTriggerCannotBeEOAError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TriggerCannotBeEOA()
func (processorEndpoint *ProcessorEndpoint) UnpackTriggerCannotBeEOAError(raw []byte) (*ProcessorEndpointTriggerCannotBeEOA, error) {
	out := new(ProcessorEndpointTriggerCannotBeEOA)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "TriggerCannotBeEOA", raw); err != nil {
		return nil, err
	}
	return out, nil
}
