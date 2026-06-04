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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resetOperator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"},{\"internalType\":\"contractITokenAllowlist\",\"name\":\"_tokenAllowlist\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeployerNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientAppBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPermit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxNumOfApplicationsExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotRegisteredTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerCannotBeEOA\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"AppEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"DeployRequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"DeployRequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMax\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"MaxNumberOfAppUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"ReportGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"withdrawSuccess\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"postWithdrawSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"returnedTokens\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"failedTokens\",\"type\":\"tuple[]\"}],\"name\":\"TriggerWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPLOYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_OPERATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"addAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminReset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"appIds\",\"type\":\"uint64[]\"}],\"name\":\"adminResetApps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"appCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"applicationStateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableDeploySlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"facilitatorNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeployedAppIds\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getFacilitatorNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getPendingRequestsPage\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTriggerQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"isAllowedDeployer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumOfApplications\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"removeAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitDeployRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"requestSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"depositPermit\",\"type\":\"bytes\"}],\"name\":\"submitRequestFor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitTriggerRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalAppCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalPendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"triggerContracts\",\"outputs\":[{\"internalType\":\"contractITrigger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"triggersToAppIds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"updateMaxNumOfApplications\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "ProcessorEndpoint",
	Bin: "0x610160806040523461034b5760e081615f968038038091610020828561034f565b83398101031261034b5780516001600160a01b0381169081900361034b5760208201516001600160a01b038116929083900361034b5761006260408201610386565b9261006f60608301610386565b61007b60808401610386565b9360c060a08501519401519560018060a01b03871680970361034b5760409687516100a6898261034f565b60048152602081016356656c6160e01b81525f906100c4600161039a565b926100d18c51948561034f565b600184526100df600161039a565b602085019390601f1901368537600a602186015b5f1901916f181899199a1a9b1b9c1cb0b131b232b360811b8282061a835304801561012157600a90916100f3565b505060018055610130816105e0565b6101205261013d8461077b565b61014052519020918260e05251902080610100524660a05289519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f84528b83015260608201524660808201523060a082015260a081526101a460c08261034f565b5190206080523060c052600a600655600a600755600a600c5582158015610343575b8015610332575b8015610321575b8015610319575b61030a57600d80546001600160a01b03199081169094179055600e80548416909517909455600f8054831694909417909355601580549091166001600160a01b0384161790556102979161022e906103b5565b506102388161044d565b505f516020615f365f395f51905f525f818152602081905286812060010180545f516020615f565f395f51905f5291829055909290917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9080a46104cd565b506014556001600160a01b0381166102fa575b505161560290816108b48239608051816152a5015260a0518161535c015260c0518161526f015260e051816152f40152610100518161531a01526101205181611738015261014051816117610152f35b6103039061054d565b505f6102aa565b632582a64160e11b5f5260045ffd5b5080156101db565b506001600160a01b038416156101d4565b506001600160a01b038216156101cd565b5084156101c6565b5f80fd5b601f909101601f19168101906001600160401b0382119082101761037257604052565b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b038216820361034b57565b6001600160401b03811161037257601f01601f191660200190565b6001600160a01b0381165f9081525f516020615ed65f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f516020615ed65f395f51905f5260205260408120805460ff191660011790553391907f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd98905f516020615eb65f395f51905f529080a4600190565b505f90565b6001600160a01b0381165f9081525f516020615ef65f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f516020615ef65f395f51905f5260205260408120805460ff191660011790553391905f516020615f565f395f51905f52905f516020615eb65f395f51905f529080a4600190565b6001600160a01b0381165f9081525f516020615f765f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f516020615f765f395f51905f5260205260408120805460ff191660011790553391905f516020615f365f395f51905f52905f516020615eb65f395f51905f529080a4600190565b6001600160a01b0381165f9081525f516020615f165f395f51905f52602052604090205460ff16610448576001600160a01b03165f8181525f516020615f165f395f51905f5260205260408120805460ff191660011790553391907fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c95917905f516020615eb65f395f51905f529080a4600190565b908151602081105f1461065a575090601f81511161061a57602081519101516020821061060b571790565b5f198260200360031b1b161790565b604460209160405192839163305a27a960e01b83528160048401528051918291826024860152018484015e5f828201840152601f01601f19168101030190fd5b6001600160401b03811161037257600254600181811c91168015610771575b602082101461075d57601f811161072a575b50602092601f82116001146106c957928192935f926106be575b50508160011b915f199060031b1c19161760025560ff90565b015190505f806106a5565b601f1982169360025f52805f20915f5b86811061071257508360019596106106fa575b505050811b0160025560ff90565b01515f1960f88460031b161c191690555f80806106ec565b919260206001819286850151815501940192016106d9565b60025f52601f60205f20910160051c810190601f830160051c015b818110610752575061068b565b5f8155600101610745565b634e487b7160e01b5f52602260045260245ffd5b90607f1690610679565b908151602081105f146107a6575090601f81511161061a57602081519101516020821061060b571790565b6001600160401b03811161037257600354600181811c911680156108a9575b602082101461075d57601f8111610876575b50602092601f821160011461081557928192935f9261080a575b50508160011b915f199060031b1c19161760035560ff90565b015190505f806107f1565b601f1982169360035f52805f20915f5b86811061085e5750836001959610610846575b505050811b0160035560ff90565b01515f1960f88460031b161c191690555f8080610838565b91926020600181928685015181550194019201610825565b60035f52601f60205f20910160051c810190601f830160051c015b81811061089e57506107d7565b5f8155600101610891565b90607f16906107c556fe6101e080604052600436101561001d575b50361561001b575f80fd5b005b5f905f3560e01c90816301ffc9a71461232e575080631256c918146122f65780631664deeb146122bc5780631753520f1461229f57806318733d571461225e578063194b6be8146122265780631cc76c38146121e657806320ff473f1461219e57806321c0b3421461216c578063248a9ca3146121425780632a0acc6a146121085780632b52b07c146120ce5780632f2ff15d146120915780632fbfa0d514611d6e57806336568abe14611d2a57806347101cae14611d105780634a0d049514611b4d5780635423112f14611b065780635e33938414611aec578063655f2de114611a715780636b288d2014611a235780636e91d133146119fb57806374202d9b146119875780637596ed43146119225780637a36a891146118ea57806380a1f71214611868578063840059ec1461181857806384b0196e146117205780638537c2781461159457806386f4a07d1461155a5780638c5b9b001461152c5780638eb8791a1461150f57806391d14854146114c75780639c73eeaf146114755780639c8d416f146114585780639fabc36f14611416578063a217fddf146113d4578063a9ba045e146113ee578063aa3aa460146113d4578063af20c9601461139c578063b7222ff514611349578063bcf7a3c914610c7a578063be1e372514610c47578063c404c35914610c06578063c415b95c14610bdd578063d2c35ce814610b65578063d365a464146107de578063d547741f1461079b578063d72e957a14610762578063dd39d76514610348578063e3d72e5e146102fb578063ecd00261146102c0578063f624f88a146102975763f7d9cc1e0361001057346102945780600319360112610294576020600c54604051908152f35b80fd5b5034610294578060031936011261029457600f546040516001600160a01b039091168152602090f35b503461029457806003193601126102945760206040517ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c8152f35b503461029457602036600319011261029457610315612381565b61031d613f13565b6001600160a01b0381161561033957610335906143e8565b5080f35b632582a64160e11b8252600482fd5b5034610294576020366003190112610294576004356001600160401b03811161075e573660238201121561075e578060040135906001600160401b03821161075a5760248260051b82010190368211610715576103a3613f82565b6103ab61415c565b6103b3614e8c565b6103bc836129db565b916103ca60405193846125c6565b8383526024602084019201915b818310610736575050509015610728575b600f5460405163024ece8960e01b8152908390829060049082906001600160a01b03165afa90811561071d578391610679575b50815191815190845b84810361043357856001805580f35b6001600160401b036104458284612c67565b5116808752601260205260408720878052602052604087205480610612575b50865b84810361059157508087526004602052604087205461048a575b50600101610424565b8087526004602052866040812055600554875b8181036104b9575b505050600190816007540160075590610481565b826001600160401b036104cb83612ccc565b90549060031b1c16146104e05760010161049d565b91505f19810190811161057d57906105136001600160401b0361050561053294612ccc565b90549060031b1c1691612ccc565b9091906001600160401b038084549260031b9316831b921b1916179055565b600554801561056957600191905f190161054b81612ccc565b6001600160401b0382549160031b1b19169055600555905f806104a5565b634e487b7160e01b87526031600452602487fd5b634e487b7160e01b88526011600452602488fd5b818852601260205260408820600191906001600160a01b036105b3838a612c67565b5116838060a01b03165f5260205260405f2054806105d3575b5001610467565b806105f461060c92338b6105ed87898060a01b0392612c67565b5116615079565b838060a01b03610604848b612c67565b511685614889565b5f6105cc565b8780808084335af1610622613ee4565b501561066a578188526012602052604088208880526020526040882061064982825461282a565b905587805260136020526106626040892091825461282a565b90555f610464565b6312171d8360e31b8852600488fd5b90503d8084833e61068a81836125c6565b810190602081830312610715578051906001600160401b03821161071957019080601f830112156107155781516106c0816129db565b926106ce60405194856125c6565b81845260208085019260051b82010192831161071157602001905b8282106106f9575050505f61041b565b6020809161070684613ed0565b8152019101906106e9565b8580fd5b8380fd5b8480fd5b6040513d85823e3d90fd5b506107316128d6565b6103e8565b82356001600160401b0381168103610756578152602092830192016103d7565b8680fd5b8280fd5b5080fd5b5034610294576020366003190112610294576020906040906001600160a01b0361078a612381565b168152601683522054604051908152f35b5034610294576040366003190112610294576103356004356107bb612397565b906107d96107d4825f525f602052600160405f20015490565b613ff1565b61448c565b503461029457604036600319011261029457600435906004821015610294576024356001600160401b03811161075e5761081c903690600401612413565b9133815260186020526001600160401b036040822054168015610b5657601c5460405161089681602081019333855285604083015261085a8a612464565b89606083015260e0608083015261087661010083018a8a612837565b908760a08401528760c084015260e083015203601f1981018352826125c6565b51902092604051956108a78761257b565b4287526020870190848252604088018581526108d360608a019487865260808b0199898b52369161261c565b9860a08101998a5260c0810198338a5260e08201958887526101008301938885526101208401968a885261014085019661090c81612464565b87528b8b52601960205260408b2094518555516001850180546001600160a01b0319166001600160a01b03929092169190911790555160028401555160038301555160048201559851805190999060058201906001600160401b038111610b425760209b61097a83546126a4565b601f8111610b08575b508c90601f8311600114610a9d57828c9d9360079695936109bb938e9f9e92610a92575b50508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b1691610a0d81612464565b610a1681612464565b60ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055601c548352601a85528160408420556001601c5401601c556040519283527fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236853394a4604051908152f35b015190505f806109a7565b838b52818b209190601f1984168c5b818110610af05750926001928e9f95929d9e9d83600799989610610ad8575b505050811b0190556109be565b01515f1960f88460031b161c191690555f8080610acb565b92936020600181928786015181550195019301610aac565b610b3290848c528e8c20601f850160051c81019160208610610b38575b601f0160051c0190614555565b5f610983565b9091508190610b25565b634e487b7160e01b89526041600452602489fd5b63578f584160e01b8252600482fd5b503461029457602036600319011261029457610b7f612381565b610b87613f13565b6001600160a01b03168015610339576020817fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f926bffffffffffffffffffffffff60a01b6015541617601555604051908152a180f35b50346102945780600319360112610294576015546040516001600160a01b039091168152602090f35b5034610294578060031936011261029457610c36610c22613e17565b60405193849360608552606085019061246e565b916020840152151560408301520390f35b503461029457604036600319011261029457610c76610c6a602435600435614f85565b6040519182918261251c565b0390f35b506101403660031901126111be57610c90612381565b6024359060ff821682036111be57604435906001600160401b03821682036111be57600460643510156111be576084356001600160401b0381116111be57610cdc903690600401612413565b909390919060a435906001600160a01b03821682036111be57610104356001600160401b0381116111be57610d15903690600401612413565b90610124356001600160401b0381116111be57610d36903690600401612413565b91909260ff851661133a576001600160401b0389165f52600460205260405f20541561132b57610d6461415c565b610d6f606435612464565b6003606435141580611313575b6113045760e43542116112f557610d9161450c565b600c5411156112e6578791610da7606435612464565b87868b8d6003606435146112ae575b610ecb9593610ebc93610ec296936001600160401b03610dfa60429560018060a01b0386165f52601660205260405f20549c610df3606435612464565b369161261c565b602081519101209160ff6040519460208601967f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc1333340885260018060a01b0316604087015216606085015216608083015260ff6064351660a083015260c082015260018060a01b038d1660e082015260c4356101008201528861012082015260e4356101408201526101408152610e91610160826125c6565b519020610e9c61526c565b906040519161190160f01b8352600283015260228201522092369161261c565b90615382565b909291926153bc565b6001600160a01b0316801561129f576001600160a01b03871603611290576001810180911161127c576001600160a01b0386165f90815260166020526040902055601454341061126d5760c43515808061125b575b61124c5715611025575b50509060ff610fd89260209850610f8d610f51600b5460c435868a8d8d606435908d612857565b98899760405195610f618761257b565b4287526001600160a01b03168c87015260c435604087015234606087015260808601899052369161261c565b60a08401526001600160a01b03851660c08401523360e08401526001600160401b03871661010084015216610120820152610fc9606435612464565b6064356101408201528361456b565b6040513381526001600160a01b03909116926001600160401b0316907fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236908690a460018055604051908152f35b6001600160a01b0384161561124c57600f5460405163cbe230c360e01b81526001600160a01b038681166004830152909160209183916024918391165afa9081156111b3575f91611212575b501561120357606081036111f45781606091810103126111be5780359060ff82168092036111be57604051636eb1769f60e11b81526001600160a01b0386811660048301523060248301526020908290604490829089165afa9081156111b3575f916111c2575b5060c43511611111575b50509060ff610fd8926020986110fb60c4358786614265565b61110860c435858a614393565b98509192610f2a565b909291906001600160a01b0383163b156111be5760409081519463d505accf60e01b865260018060a01b038716600487015230602487015260c435604487015260e43560648701526084860152602081013560a4860152013560c48401525f8360e4818360018060a01b0387165af19081156111b357602098610fd89460ff936111a0575b50985091926110e2565b6111ac91505f906125c6565b5f5f611196565b6040513d5f823e3d90fd5b5f80fd5b90506020813d6020116111ec575b816111dd602093836125c6565b810103126111be57515f6110d8565b3d91506111d0565b63ddafbaef60e01b5f5260045ffd5b63514e24c360e11b5f5260045ffd5b90506020813d602011611244575b8161122d602093836125c6565b810103126111be5761123e906125e7565b5f611071565b3d9150611220565b632a9ffab760e21b5f5260045ffd5b506001600160a01b0385161515610f20565b63c77f97eb60e01b5f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b632057875960e21b5f5260045ffd5b638baa579f60e01b5f5260045ffd5b505050509160851415806112db575b6112cc57879187868b8d610db6565b637c6953f960e01b5f5260045ffd5b5060e28814156112bd565b63d8219b3d60e01b5f5260045ffd5b631ab7da6b60e01b5f5260045ffd5b630ee2db0560e21b5f5260045ffd5b5061131f606435612464565b60016064351415610d7c565b6309ff9a8360e41b5f5260045ffd5b635428eae760e01b5f5260045ffd5b346111be5760403660031901126111be576113626123d7565b6001600160401b03611372612397565b91165f52601260205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b346111be5760203660031901126111be576001600160a01b036113bd612381565b165f526013602052602060405f2054604051908152f35b346111be575f3660031901126111be5760206040515f8152f35b346111be575f3660031901126111be57600e546040516001600160a01b039091168152602090f35b346111be5760203660031901126111be576020600435611435816147d1565b908115611448575b506040519015158152f35b61145291506147f9565b8261143d565b346111be575f3660031901126111be576020601454604051908152f35b346111be5760203660031901126111be57600435611491613f13565b801561124c576020817e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db192600c55604051908152a1005b346111be5760403660031901126111be576114e0612397565b6004355f525f60205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b346111be575f3660031901126111be576020600654604051908152f35b346111be575f3660031901126111be57611544613f82565b61154c61415c565b611554614e8c565b60018055005b346111be575f3660031901126111be5760206040517fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c959178152f35b346111be576101803660031901126111be576115ae6123d7565b608435906001600160401b0382116111be57604060031983360301126111be5760a4356001600160401b0381116111be57604060031982360301126111be5760c435916001600160401b0383116111be57366023840112156111be578260040135926001600160401b0384116111be5736602460608602830101116111be5761012435600d8110156111be57610144356001600160401b0381116111be5761165a903690600401612413565b939092610164356001600160401b0381116111be5761167d903690600401612413565b335f9081527e9b85fa08e135b4afefc8a40566229b01e66445186a3676642ffbacd3d93b0e6020526040902054909891979060ff16156116e957611554996116c361415c565b6101043594602460e4359501926004019160040190606435906044359060243590612d0c565b63e2517d3f60e01b5f52336004527f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9860245260445ffd5b346111be575f3660031901126111be576117bc61175c7f00000000000000000000000000000000000000000000000000000000000000006150ba565b6117857f00000000000000000000000000000000000000000000000000000000000000006151b4565b60206117ca6040519261179883856125c6565b5f84525f368137604051958695600f60f81b875260e08588015260e0870190612440565b908582036040870152612440565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b82811061180157505050500390f35b8351855286955093810193928101926001016117f2565b346111be5760403660031901126111be57611831612381565b611839612397565b6001600160a01b039182165f908152601060209081526040808320949093168252928352819020549051908152f35b346111be575f3660031901126111be5761188861188361450c565b614782565b600a54600b54905f905b8281106118a75760405180610c76868261251c565b60018181925f52600960205260405f20545f5260086020526118cb60405f2061277c565b6118d58588612c67565b526118e08487612c67565b5001910190611892565b346111be5760203660031901126111be576001600160401b0361190b6123d7565b165f526004602052602060405f2054604051908152f35b346111be575f3660031901126111be5761193a6128d6565b6040518091602082016020835281518091526020604084019201905f5b818110611965575050500390f35b82516001600160401b0316845285945060209384019390920191600101611957565b346111be5760e03660031901126111be576119a0612381565b6119a86123ed565b906044359160048310156111be57606435916001600160401b0383116111be576020936119dc6119f3943690600401612413565b906119e56123ad565b9260c4359560a43595612857565b604051908152f35b346111be575f3660031901126111be57600d546040516001600160a01b039091168152602090f35b346111be5760203660031901126111be57611a3c612381565b6001600160a01b03165f9081525f5160206155ad5f395f51905f52602090815260409182902054915160ff9092161515825290f35b346111be5760203660031901126111be57600435611a8d613f13565b801561124c57600654611aa26007548261282a565b9081831061124c5782611adc6040937f80af1d8add7b822532307ae28ef73c31ab2803f908c5c22cc879766f54c68583956006558261282a565b60075582519182526020820152a1005b346111be575f3660031901126111be5760206119f361450c565b346111be5760203660031901126111be57611b1f612652565b506004355f526008602052610c76611b3960405f2061277c565b60405191829160208352602083019061246e565b60403660031901126111be57611b61612403565b602435906001600160401b0382116111be57611b8360ff923690600401612413565b92909116908161133a57611b9561415c565b335f9081525f5160206155ad5f395f51905f52602052604090205460ff1615611d01576007548015611cf257611bc961450c565b600c5411156112e657601454341061126d575f1901600755600b546040805133602082019081525f92820192909252949085611c4181611ca6955f606060209b015260e06080830152611c216101008301878a612837565b905f60a08401525f60c084015260e083015203601f1981018352826125c6565b51902093611c7b8560c01c9460405193611c5a8561257b565b4285525f898601525f6040860152346060860152876080860152369161261c565b60a08301523360c08301525f60e0830152836101008301526101208201525f6101408201528361456b565b604051908282527fbbcfb0c180d5f95cfb50c12005c44144ac92ac8474f1abc5add0e6e69a14e8dc843393a360018055604051908152f35b634e487b7160e01b5f52602160045260245ffd5b6316ae067f60e11b5f5260045ffd5b63fc336c4160e01b5f5260045ffd5b346111be575f3660031901126111be5760206119f361452b565b346111be5760403660031901126111be57611d43612397565b336001600160a01b03821603611d5f5761001b9060043561448c565b63334bd91960e11b5f5260045ffd5b60e03660031901126111be57611d82612403565b611d8a6123ed565b60443560048110156111be576064356001600160401b0381116111be57611db5903690600401612413565b611dc09291926123ad565b9160a4359460ff60c4359716908161133a576001600160401b03811695865f52600460205260405f20541561132b57611df761415c565b611e0084612464565b831561130457601454891061126d57611e1761450c565b600c5411156112e6576001600160a01b0386169485611ff957611e3a8a8a6125f4565b340361124c575b611e4a85612464565b60038503611f2c57608581141580611f21575b6112cc57602099611e89611eb5948b611ee89a5b611e7c828285614393565b85878b600b549533612857565b9960405197611e978961257b565b4289528c89015260408801526060870152886080870152369161261c565b60a08401523360c08401525f60e084015284610100840152610120830152611edc81612464565b6101408201528361456b565b81604051915f83527fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236853394a460018055604051908152f35b5060e2811415611e5d565b919092939495611f3b86612464565b60028614611f5d575b91602099611e89611ee898979695938b611eb596611e71565b9594939290916044602060018060a01b03600e5416604051928380926308b56a5360e11b82528d60048301523360248301525afa9081156111b3575f91611fbf575b5015611fb057919092939495611f44565b63622a850b60e11b5f5260045ffd5b90506020813d602011611ff1575b81611fda602093836125c6565b810103126111be57611feb906125e7565b8b611f9f565b3d9150611fcd565b89340361124c57881561124c57600f5460405163cbe230c360e01b81526004810188905290602090829060249082906001600160a01b03165afa9081156111b3575f91612057575b501561120357612052893389614265565b611e41565b90506020813d602011612089575b81612072602093836125c6565b810103126111be57612083906125e7565b8b612041565b3d9150612065565b346111be5760403660031901126111be5761001b6004356120b0612397565b906120c96107d4825f525f602052600160405f20015490565b6140d4565b346111be575f3660031901126111be5760206040517f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc13333408152f35b346111be575f3660031901126111be5760206040517fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec428152f35b346111be5760203660031901126111be5760206119f36004355f525f602052600160405f20015490565b346111be5760403660031901126111be57611554612188612381565b612190612397565b9061219961415c565b61417c565b346111be5760203660031901126111be576121b7612381565b6121bf613f13565b6001600160a01b038116156121d75761001b90614029565b632582a64160e11b5f5260045ffd5b346111be5760203660031901126111be576001600160401b036122076123d7565b165f526017602052602060018060a01b0360405f205416604051908152f35b346111be5760203660031901126111be576001600160a01b03612247612381565b165f526011602052602060405f2054604051908152f35b346111be5760203660031901126111be576001600160a01b0361227f612381565b165f52601860205260206001600160401b0360405f205416604051908152f35b346111be575f3660031901126111be576020600754604051908152f35b346111be575f3660031901126111be5760206040517f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd988152f35b346111be5760203660031901126111be576001600160a01b03612317612381565b165f526016602052602060405f2054604051908152f35b346111be5760203660031901126111be576004359063ffffffff60e01b82168092036111be57602091637965db0b60e01b8114908115612370575b5015158152f35b6301ffc9a760e01b14905083612369565b600435906001600160a01b03821682036111be57565b602435906001600160a01b03821682036111be57565b608435906001600160a01b03821682036111be57565b35906001600160a01b03821682036111be57565b600435906001600160401b03821682036111be57565b602435906001600160401b03821682036111be57565b6004359060ff821682036111be57565b9181601f840112156111be578235916001600160401b0383116111be57602083818601950101116111be57565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60041115611cde57565b908151815260018060a01b036020830151166020820152604082015160408201526060820151606082015260808201516080820152610140806124c260a085015161016060a0860152610160850190612440565b9360018060a01b0360c08201511660c085015260018060a01b0360e08201511660e08501526001600160401b036101008201511661010085015260ff6101208201511661012085015201519161251783612464565b015290565b602081016020825282518091526040820191602060408360051b8301019401925f915b83831061254e57505050505090565b909192939460208061256c600193603f19868203018752895161246e565b9701930193019193929061253f565b61016081019081106001600160401b0382111761259757604052565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b0382111761259757604052565b90601f801991011681019081106001600160401b0382111761259757604052565b519081151582036111be57565b9190820180921161127c57565b6001600160401b03811161259757601f01601f191660200190565b92919261262882612601565b9161263660405193846125c6565b8294818452818301116111be578281602093845f960137010152565b6040519061265f8261257b565b5f61014083828152826020820152826040820152826060820152826080820152606060a08201528260c08201528260e082015282610100820152826101208201520152565b90600182811c921680156126d2575b60208310146126be57565b634e487b7160e01b5f52602260045260245ffd5b91607f16916126b3565b9060405191825f8254926126ef846126a4565b808452936001811690811561275a5750600114612716575b50612714925003836125c6565b565b90505f9291925260205f20905f915b81831061273e575050906020612714928201015f612707565b6020919350806001915483858901015201910190918492612725565b90506020925061271494915060ff191682840152151560051b8201015f612707565b906040516127898161257b565b61014060ff600783958054855260018060a01b0360018201541660208601526002810154604086015260038101546060860152600481015460808601526127d2600582016126dc565b60a086015260018060a01b0360068201541660c0860152015460018060a01b03811660e08501526001600160401b038160a01c16610100850152818160e01c1661012085015260e81c169161282683612464565b0152565b9190820391821161127c57565b908060209392818452848401375f828201840152601f01601f1916010190565b96916128d095936001600160401b039792956128a892604051998a9860208a019c60018060a01b03168d5216604089015261289181612464565b606088015260e06080880152610100870191612837565b6001600160a01b0390931660a085015260c084015260e083015203601f1981018352826125c6565b51902090565b60405190600554808352826020810160055f5260205f20925f905b8060038301106129895761271494549181811061296f575b818110612952575b818110612935575b10612927575b5003836125c6565b60c01c81526020015f61291f565b9260206001916001600160401b038560801c168152019301612919565b9260206001916001600160401b038560401c168152019301612911565b9260206001916001600160401b0385168152019301612909565b916004919350608060019186546001600160401b03811682526001600160401b038160401c1660208301526001600160401b0381841c16604083015260c01c60608201520194019201859293916128f1565b6001600160401b0381116125975760051b60200190565b906129fc826129db565b612a0960405191826125c6565b8281528092612a1a601f19916129db565b0190602036910137565b903590601e19813603018212156111be57018035906001600160401b0382116111be57602001918160051b360383136111be57565b91906040838203126111be57604051612a71816125ab565b809380356001600160401b0381116111be57810183601f820112156111be578035612a9b816129db565b91612aa960405193846125c6565b81835260208084019260051b820101908682116111be5760208101925b828410612b48575050505082526020810135906001600160401b0382116111be57019180601f840112156111be578235612aff816129db565b93612b0d60405195866125c6565b81855260208086019260051b8201019283116111be57602001905b828210612b385750505060200152565b8135815260209182019101612b28565b83356001600160401b0381116111be57820188603f820112156111be57602091612b7b8a8360408680960135910161261c565b815201930192612ac6565b60408201908051916040845282518091526060840190602060608260051b8701019401915f905b828210612bf7575050505060200151916020818303910152602080835192838152019201905f5b818110612be15750505090565b8251845260209384019390920191600101612bd4565b90919294602080612c14600193605f198b82030186528951612440565b970192019201909291612bad565b90600d821015611cde5752565b9190811015612c3f576060020190565b634e487b7160e01b5f52603260045260245ffd5b356001600160a01b03811681036111be5790565b8051821015612c3f5760209160051b010190565b9190811015612c3f5760051b0190565b9190811015612c3f5760051b81013590601e19813603018212156111be5701908135916001600160401b0383116111be5760200182360381136111be579190565b90600554821015612c3f5760055f52600282901c7f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0019160031b60181690565b9d9c9b999d9a9897969594939a6101805260a0526101a0526101c052610140526101605260805260e05261012052612d466101c0516147d1565b8015613e05575b15613df657612d5e6101c0516147d1565b5f61010052908115613dde576101c0515f52601960205260405f20610100525b600761010051015460c0526001600160401b0360c05160a01c166001600160401b0361018051160361132b576001600160401b0361018051165f52600460205260405f205460a051036139b057612dd58480612a24565b959050612de56020860186612a24565b905086036112cc57612dfd6101405161014051612a24565b979050612e136020610140510161014051612a24565b905088036112cc5760405191612e288361257b565b6001600160401b036101805116835260a05160208401526101a05160408401526101c0516060840152612e5b3688612a59565b6080840152612e6d3661014051612a59565b60a0840152612e7d6080516129db565b91612e8b60405193846125c6565b608051835260208301366060608051026101605101116111be5761016051905b60606080510261016051018210613d8157505060c0840192835260e084019260e051845261010085016101205181528c600d6101208801911015611cde5790818e612f8d98969497959352612f01368c8b61261c565b91610140870192835260018060a01b03600d541695604051998a98635c9626b960e01b8a52604060048b01526001600160401b0381511660448b0152602081015160648b0152604081015160848b0152606081015160a48b015260a0612f798b61016060c460808601519201526101a48d0190612b86565b9101518a82036043190160e48c0152612b86565b985198604319898203016101048a01526020808b5192838152019a01905f5b818110613d415750505060209861300994612fe48a9895899895612ff795516101248b0152516101448a015251610164890190612c22565b5186820360431901610184880152612440565b84810360031901602486015291612837565b03915afa9081156111b3575f91613d07575b501561129f5761010051600381015460069091015460c0516001600160a01b039182169992959291168015613d0057985b816139bf57505050506001600160401b0361018051165f5260046020526101a05160405f2054146139b0571561398a575b505f5f9061308c6080516129f2565b6130976080516129f2565b905f925b608051811061380857505f5b83811061374e57505050506130ca906130c56101205160e0516125f4565b6125f4565b6130d2614821565b1061373f575f5b8281036136ca575050505f5b81810361363d57505060ff600761010051015460e81c169061310682612464565b60028214613603575b6001600160401b0361018051165f5260046020526101a05160405f205561313582612464565b811590816134d4575b60405160a05181526101a05160208201526101c051907f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e60406001600160401b03610180511692a360e051613483575b505f5f6001600160401b0361018051165f52601760205260018060a01b0360405f2054166131bd6080516129f2565b925b60805181106133295750506131d3816129f2565b915f5b828103613301575050506131f590610140516101c05161018051614b8b565b61322c604051926132076020856125c6565b5f845261321261452b565b1515806132ef575b156132e7576132276154cd565b612464565b1561329757604051907f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d6101c05192806132796001600160401b036101805116945f80610120518561497f565b0390a35b6101205160155461271491906001600160a01b03166148da565b604051907f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed666101c05192806132df6001600160401b036101805116945f80610120518561497f565b0390a361327d565b613227615490565b506132fc6101c0516147d1565b61321a565b6001906001600160a01b036133168285612c67565b51166133228287612c67565b52016131d6565b8160018060a01b0361334c60206133468560805161016051612c2f565b01612c53565b161461344e575b806133a961337161336c60019460805161016051612c2f565b612c53565b61338660206133468560805161016051612c2f565b9060406133998560805161016051612c2f565b013591858060a01b031690614935565b6133be60206133468360805161016051612c2f565b6133d161336c8360805161016051612c2f565b60406133e38460805161016051612c2f565b0135907faafeb614695d37710d00f37c69671d205eb5b6fcd02995f6c3c7e2c93d3cbb2460405193868060a01b031693806134456101c051956001600160401b036101805116958360209093929193604081019460018060a01b031681520152565b0390a4016131bf565b916001809161346661336c8660805161016051612c2f565b6134708288612c67565b90838060a01b0316905201929050613353565b60018060a01b031661349760e051826148da565b604051604081015f825260e05160208301525f51602061558d5f395f51905f526101c05192806001600160401b036101805116930390a45f61318e565b60055468010000000000000000811015612597576134fb81600161351c9301600555612ccc565b61018051906001600160401b038084549260031b9316831b921b1916179055565b60056101005101602061352f82546126a4565b1461353b575b5061313e565b613544906126dc565b6020818051810103126111be57602001516001600160a01b038116908190036111be57801561353557805f5260186020526001600160401b0360405f2054166135f457803b156135e557610180516001600160401b03165f81815260176020908152604080832080546001600160a01b0319166001600160a01b03871617905593825260189052918220805467ffffffffffffffff19169091179055613535565b634ea29e6760e01b5f5260045ffd5b63d48c379760e01b5f5260045ffd5b6101c0516001600160401b0361018051167ff36f013d86705b333febcccf7ece97eee30f4cd398a9c69aee1f940ed5daddea5f80a361310f565b8061365d6001926136576020610140510161014051612a24565b90612c7b565b35613678826136726101405161014051612a24565b90612c8b565b7f485d4bfb37695319d1b2352eef5982d10db8c284b0a7ddefe4465377c9fbc8df60405160208152806136c16101c051956001600160401b036101805116956020840191612837565b0390a4016130e5565b806136de6001926136576020860186612a24565b356136ed826136728680612a24565b7f9ede1ad14e46a641d7b174bb6c9d939ad961d1a7de4b669163b9714a67e2232060405160208152806137366101c051956001600160401b036101805116956020840191612837565b0390a4016130d9565b631e9acf1760e31b5f5260045ffd5b6001600160a01b036137608284612c67565b516040516370a0823160e01b81523060048201529116602082602481845afa9182156111b3575f926137d3575b506137b5816137c6925f52601360205260405f2054905f52601160205260405f2054906125f4565b6137bf8487612c67565b51906125f4565b1161373f576001016130a7565b9091506020813d8211613800575b816137ee602093836125c6565b810103126111be5751906137b561378d565b3d91506137e1565b9361381c61336c8660805161016051612c2f565b90604061382f8760805161016051612c2f565b0135916001600160401b0361018051165f52601260205260405f2060018060a01b0382165f5260205260405f2054831161397b57610180516001600160401b03165f9081526012602090815260408083206001600160a01b03851684528252808320805487900390556013909152902054831161373f576001600160a01b0381165f818152601360205260409020805485900390556138dd57506001916138d5916125f4565b945b0161309b565b909591905f805b87811061392d575b50156138fd575b50506001906138d7565b91946001928392906001600160a01b03166139188387612c67565b526139238287612c67565b520193905f6138f3565b6001600160a01b03838116906139438389612c67565b511614613952576001016138e4565b905061397261396b84613965848a612c67565b516125f4565b9187612c67565b5260015f6138ec565b6306a1762560e11b5f5260045ffd5b6139996101205160e0516125f4565b0361124c57601454610120511061124c575f61307d565b630b6fac0360e41b5f5260045ffd5b94929791969095501590811591613cf6575b508015613ceb575b6112cc576001600160401b0361018051165f5260046020526101a05160405f2054036139b05760026101005101549060018060a01b03600161010051015416926001600160401b0361018051165f52601260205260405f2060018060a01b0385165f5260205260405f2054831161397b57159081613cda575b5061373f5761271496613a69828461018051614889565b613a7d60036101005101546014549061282a565b9280613c6f57508115613c0d575090613a95916125f4565b613a9f81836148da565b6040519060408201905f835260208301525f51602061558d5f395f51905f526101c05192806001600160401b036101805116930390a45b60ff600761010051015460e81c16613aed81612464565b15613bff575b613b106014549360ff600761010051015460e81c1693369161261c565b91613b1961452b565b151580613bed575b15613be057613b2e6154cd565b613b3781612464565b613b94577f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d60405180613b7f6101c051956001600160401b036101805116956001898561497f565b0390a35b6015546001600160a01b03166148da565b7f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed6660405180613bd86101c051956001600160401b036101805116956001898561497f565b0390a3613b83565b613be8615490565b613b2e565b50613bfa6101c0516147d1565b613b21565b600160075401600755613af3565b9192505081613c1e575b5050613ad6565b6001600160a01b031690613c3281836148da565b6040519060408201905f835260208301525f51602061558d5f395f51905f526101c05192806001600160401b036101805116930390a45f80613c17565b8291939492613c88575b50505081613c1e575050613ad6565b613c93828483614935565b604080516101c051610180516001600160a01b0394909416825260208201949094526001600160401b03909216915f51602061558d5f395f51905f529190a45f8080613c79565b9050613ce4614821565b105f613a52565b5060805115156139d9565b905015155f6139d1565b508861304c565b90506020813d602011613d39575b81613d22602093836125c6565b810103126111be57613d33906125e7565b5f61301b565b3d9150613d15565b825180516001600160a01b039081168e52602082810151909116818f0152604091820151918e01919091526060909c019b8d9b5090920191600101612fac565b6060823603126111be576040519060608201908282106001600160401b0383111761259757606092602092604052613db8856123c3565b8152613dc58386016123c3565b8382015260408501356040820152815201910190612eab565b6101c0515f52600860205260405f2061010052612d7e565b6302e8145360e61b5f5260045ffd5b50613e126101c0516147f9565b612d4d565b613e1f612652565b50613e2861452b565b613e8f57613e3461450c565b613e4757613e40612652565b905f905f90565b600a545f52600960205260405f20545f52600860205260405f20906001600160401b03600783015460a01c165f526004602052613e8860405f20549261277c565b9190600190565b601b545f52601a60205260405f20545f52601960205260405f20906001600160401b03600783015460a01c165f526004602052613e8860405f20549261277c565b51906001600160a01b03821682036111be57565b3d15613f0e573d90613ef582612601565b91613f0360405193846125c6565b82523d5f602084013e565b606090565b335f9081527f5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b7602052604090205460ff1615613f4b57565b63e2517d3f60e01b5f52336004527fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4260245260445ffd5b335f9081527f3f65c445f65be3b772e454052bf8f1f970a2679db965b2efcdbc04d3dd490789602052604090205460ff1615613fba57565b63e2517d3f60e01b5f52336004527fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c9591760245260445ffd5b5f8181526020818152604080832033845290915290205460ff16156140135750565b63e2517d3f60e01b5f523360045260245260445ffd5b6001600160a01b0381165f9081525f5160206155ad5f395f51905f52602052604090205460ff166140cf576001600160a01b03165f8181525f5160206155ad5f395f51905f5260205260408120805460ff191660011790553391907ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b505f90565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16614156575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b60026001541461416d576002600155565b633ee5aeb560e01b5f5260045ffd5b60018060a01b03811691825f52601060205260405f2060018060a01b0382165f5260205260405f205491821561425f577f16d46b28f1a9377c2df78652a0a8e31568b99311665bcec24cbba72edd2e5c2c8392855f52601060205260405f2060018060a01b0382165f526020525f6040812055855f52601160205260405f2061420685825461282a565b9055604080516001600160a01b0394851681526020810195909552921692839290a282614256575f809350809281925af161423f613ee4565b501561424757565b6312171d8360e31b5f5260045ffd5b61271492615079565b50505050565b6040516370a0823160e01b81523060048201526001600160a01b039190911691602082602481865afa9182156111b3575f9261435e575b506024926142e3602092604051906323b872dd60e01b8583015260018060a01b031686820152306044820152866064820152606481526142dd6084826125c6565b82615438565b6040516370a0823160e01b815230600482015293849182905afa80156111b3575f9061432a575b614314925061282a565b0361431b57565b63b0a6ea2960e01b5f5260045ffd5b506020823d602011614356575b81614344602093836125c6565b810103126111be57614314915161430a565b3d9150614337565b9091506020813d60201161438b575b8161437a602093836125c6565b810103126111be575190602461429c565b3d915061436d565b6001600160401b03165f52601260205260405f2060018060a01b0382165f5260205260405f206143c48382546125f4565b905560018060a01b03165f5260136020526143e460405f209182546125f4565b9055565b6001600160a01b0381165f9081525f5160206155ad5f395f51905f52602052604090205460ff16156140cf576001600160a01b03165f8181525f5160206155ad5f395f51905f5260205260408120805460ff191690553391907ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff1615614156575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b600b54600a5480821161451f5750505f90565b6145289161282a565b90565b601c54601b5480821161451f5750505f90565b6002600382015491015480821161451f5750505f90565b818110614560575050565b5f8155600101614555565b90815f52600860205260405f208151815560018060a01b03602083015116600182019060018060a01b03166bffffffffffffffffffffffff60a01b8254161790556040820151600282015560608201516003820155608082015160048201556005810160a08301518051906001600160401b038211612597576145ee83546126a4565b601f8111614752575b50602090601f83116001146146eb57918061462b9260079695945f92610a925750508160011b915f199060031b1c19161790565b90555b60c08301516006820180546001600160a01b03199081166001600160a01b039384161790915560e08086015194909301805490911693909116929092178083556101008401516101208501516101409095015193949093921b60ff60e01b169161469781612464565b6146a081612464565b60ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600b545f52600960205260405f20556001600b5401600b55565b90601f19831691845f52815f20925f5b81811061473a575091600193918560079897969410614722575b505050811b01905561462e565b01515f1960f88460031b161c191690555f8080614715565b929360206001819287860151815501950193016146fb565b61477c90845f5260205f20601f850160051c81019160208610610b3857601f0160051c0190614555565b5f6145f7565b9061478c826129db565b61479960405191826125c6565b82815280926147aa601f19916129db565b01905f5b8281106147ba57505050565b6020906147c5612652565b828285010152016147ae565b601c5490601b5480921191826147e657505090565b9091505f52601a60205260405f20541490565b600b5490600a54809211918261480e57505090565b9091505f52600960205260405f20541490565b5f805260116020527f4ad3b33220dddc71b994a52d72c06b10862965f7d926534c05c00fb7e819e7b75461452890614859904761282a565b5f805260136020527f8fa6efc3be94b5b348b21fea823fe8d100408cee9b7f90524494500445d8ff6c549061282a565b6001600160401b03165f52601260205260405f2060018060a01b0382165f5260205260405f206148ba83825461282a565b905560018060a01b03165f5260136020526143e460405f2091825461282a565b6001600160a01b03165f9081527f6e0956cda88cad152e89927e53611735b61a5c762d1428573c6931b0a5efcb0160205260409020805461491c9083906125f4565b90555f805260116020526143e460405f209182546125f4565b60018060a01b031690815f52601060205260405f209060018060a01b03165f5260205260405f206149678382546125f4565b90555f5260116020526143e460405f209182546125f4565b909392919381526002841015611cde576149a86080926145289560208401526040830190612c22565b8160608201520190612440565b9035601e19823603018112156111be5701602081359101916001600160401b0382116111be578160051b360383136111be57565b9060408101916149f981806149b5565b809460408552526060830160608560051b85010194825f90601e19813603015b838310614a6457505050505050806020614a349201906149b5565b82840360209093019290925281835290916001600160fb1b0383116111be5760209260051b809284830137010190565b909192939497605f198882030186528835828112156111be57830190602082359201916001600160401b0381116111be5780360383136111be57614aae6020928392600195612837565b9a0196019493019190614a19565b81601f820112156111be57805190614ad3826129db565b92614ae160405194856125c6565b82845260208085019360061b830101918183116111be57602001925b828410614b0b575050505090565b6040848303126111be5760206040918251614b25816125ab565b614b2e87613ed0565b81528287015183820152815201930192614afd565b90602080835192838152019201905f5b818110614b605750505090565b825180516001600160a01b031685526020908101518186015260409094019390920191600101614b53565b90925f906001600160401b03831693845f52601760205260018060a01b0360405f205416918215614e00578051905f5b828103614dda575050506001823b156111be576040516338d9814d60e11b8152602060048201525f8180614bf260248201876149e9565b038183885af19081614dc5575b50614dbf575082935b86867fca435620e7fa2ff484abda442975d2167d5e6577548dfdfac799cecc2a2af2c9602060405198151598898152a3839460609285614c6b859660019383604051809781958294639cb86de760e01b84526040600485015260448401906149e9565b90602483015203925af1868188948993614d40575b50614d315750505050835b908251915b828103614cfa575050507fa59edce59a4211e20e9f86e63adea41a55dbe409040bc17f473b671e31bd036a939291614ce7614cf5926040519586951515865215156020860152608060408601526080850190614b43565b908382036060850152614b43565b0390a3565b600190614d2b6001600160a01b03614d128388612c67565b5151166020614d218489612c67565b5101519085614393565b01614c90565b97509095909450909250614c8b565b94509150503d8088853e614d5481856125c6565b8301928581850312614dbb57614d69816125e7565b9360208201516001600160401b038111614db75781614d89918401614abc565b916040810151906001600160401b038211614db357614da9929101614abc565b939093915f614c80565b8a80fd5b8980fd5b8780fd5b93614c08565b614dd29195505f906125c6565b5f935f614bff565b600190614dfa866001600160a01b03614df38487612c67565b511661417c565b01614bbb565b50505050505050565b60075f9182815582600182015582600282015582600382015582600482015560058101614e3681546126a4565b9081614e49575b50508260068201550155565b81601f869311600114614e605750555b5f80614e3d565b81835260208320614e7c91601f0160051c810190600101614555565b8082528160208120915555614e59565b600a54600b545f905b808303614eae57509050600a54600b5560075401600755565b825f52600960205260405f20545f52600860205260405f2092600784015460ff8160e81c16614edc81612464565b15614f79575b84600260019495960180549283614f2e575b50505050805f52600960205260405f20545f526008602052614f1860405f20614e09565b5f81815260096020526040812055019190614e95565b614f54614f7094878501926001600160401b03898060a01b038554169160a01c16614889565b5460069092015490549160a086901b8690039182169116614935565b5f808080614ef4565b60019390930192614ee2565b919091614f92600861453e565b92838210801590615071575b61503557614fac90826125f4565b9280841161502d575b50614fda614fc6611883838661282a565b93614fd4600a5493846125f4565b926125f4565b905f905b828110614fea57505050565b60018181925f52600960205260405f20545f52600860205261500e60405f2061277c565b6150188589612c67565b526150238488612c67565b5001910190614fde565b92505f614fb5565b50909150506040516150486020826125c6565b5f81525f805b81811061505a57505090565b602090615065612652565b8282860101520161504e565b508015614f9e565b60405163a9059cbb60e01b60208201526001600160a01b03929092166024830152604480830193909352918152612714916150b56064836125c6565b615438565b60ff81146151005760ff811690601f82116150f157604051916150de6040846125c6565b6020808452838101919036833783525290565b632cd44ac360e21b5f5260045ffd5b50604051600254815f615112836126a4565b80835292600181169081156151955750600114615136575b614528925003826125c6565b5060025f90815290917f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace5b8183106151795750509060206145289282010161512a565b6020919350806001915483858801015201910190918392615161565b6020925061452894915060ff191682840152151560051b82010161512a565b60ff81146151d85760ff811690601f82116150f157604051916150de6040846125c6565b50604051600354815f6151ea836126a4565b8083529260018116908115615195575060011461520d57614528925003826125c6565b5060035f90815290917fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b5b8183106152505750509060206145289282010161512a565b6020919350806001915483858801015201910190918392615238565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03161480615359575b156152c7577f000000000000000000000000000000000000000000000000000000000000000090565b60405160208101907f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f82527f000000000000000000000000000000000000000000000000000000000000000060408201527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a082015260a081526128d060c0826125c6565b507f0000000000000000000000000000000000000000000000000000000000000000461461529e565b81519190604183036153b2576153ab9250602082015190606060408401519301515f1a9061550a565b9192909190565b50505f9160029190565b6153c581612464565b806153ce575050565b6153d781612464565b600181036153ee5763f645eedf60e01b5f5260045ffd5b6153f781612464565b60028103615412575063fce698f760e01b5f5260045260245ffd5b60039061541e81612464565b146154265750565b6335e2f38360e21b5f5260045260245ffd5b905f602091828151910182855af1156111b3575f513d61548757506001600160a01b0381163b155b6154675750565b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b60011415615460565b600a545f52600960205260405f20545f5260086020526154b260405f20614e09565b600a545f5260096020525f60408120556001600a5401600a55565b601b545f52601a60205260405f20545f5260196020526154ef60405f20614e09565b601b545f52601a6020525f60408120556001601b5401601b55565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411615581579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa156111b3575f516001600160a01b0381161561557757905f905f90565b505f906001905f90565b5050505f916003919056fec9961aaccfc3406b40d0d25a5b83b8a4721e5f454116d7e22bd8038c715140aa740c5e3e456bed56f053f960110118ba9b95a1f5359a82283516fb2e81b6e37ea26469706673582212205973fc339e38fccc46c09fa3f8a4ba37e885cf47d0a659ee73c229993c2bf90764736f6c634300081e00332f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d009b85fa08e135b4afefc8a40566229b01e66445186a3676642ffbacd3d93b0e5cbfc8ee58ca47855df7bcf648dd304ddb6b932f9b87878bdf6318d7ec7ee5b73f65c445f65be3b772e454052bf8f1f970a2679db965b2efcdbc04d3dd490789fc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184cdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec42740c5e3e456bed56f053f960110118ba9b95a1f5359a82283516fb2e81b6e37e",
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
// the contract method with ID 0xd365a464.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitTriggerRequest(uint8 requestType, bytes payload) returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitTriggerRequest(requestType uint8, payload []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("submitTriggerRequest", requestType, payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitTriggerRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd365a464.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitTriggerRequest(uint8 requestType, bytes payload) returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitTriggerRequest(requestType uint8, payload []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitTriggerRequest", requestType, payload)
}

// UnpackSubmitTriggerRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd365a464.
//
// Solidity: function submitTriggerRequest(uint8 requestType, bytes payload) returns(bytes32)
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
