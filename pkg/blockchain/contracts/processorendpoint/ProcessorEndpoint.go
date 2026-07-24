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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeployerNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientAppBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPermit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxNumOfApplicationsExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerCannotBeEOA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"AppEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"DeployRequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"DeployRequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMax\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"MaxNumberOfAppUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"ReportGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"withdrawSuccess\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"postWithdrawSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"returnedTokens\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"failedTokens\",\"type\":\"tuple[]\"}],\"name\":\"TriggerWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPLOYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_OPERATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"addAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminReset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"appIds\",\"type\":\"uint64[]\"}],\"name\":\"adminResetApps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"appCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"applicationStateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableDeploySlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"facilitatorNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"payloadHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeployedAppIds\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getFacilitatorNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getPendingRequestsPage\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTriggerQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTriggerRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resetOperator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"},{\"internalType\":\"contractITokenAllowlist\",\"name\":\"_tokenAllowlist\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"isAllowedDeployer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumOfApplications\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"removeAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitDeployRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"trigger\",\"type\":\"address\"}],\"name\":\"submitDeployRequestWithTrigger\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"requestSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"depositPermit\",\"type\":\"bytes\"}],\"name\":\"submitRequestFor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalAppCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalPendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"triggerContracts\",\"outputs\":[{\"internalType\":\"contractITrigger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"triggersToAppIds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"updateMaxNumOfApplications\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "ProcessorEndpoint",
	Bin: "0x60a080604052346100c257306080525f5160206167e15f395f51905f525460ff8160401c166100b3576002600160401b03196001600160401b03821601610060575b60405161671a90816100c782396080518181816121a101526122320152f35b6001600160401b0319166001600160401b039081175f5160206167e15f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610041565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe6101c080604052600436101561001d575b50361561001b575f80fd5b005b5f905f3560e01c90816301ffc9a714612aa7575080631256c91814612a6f5780631664deeb14612a355780631753520f14612a1857806318733d57146129d7578063194b6be81461299f5780631cc76c381461295f5780631f56b792146128e057806320ff473f1461289857806321c0b34214612866578063248a9ca31461282f5780632a0acc6a146128085780632b52b07c146127ce5780632f2ff15d146127845780632fbfa0d51461248557806336568abe1461244157806347101cae146124275780634a0d0495146123c25780634f1ef286146121f557806352d1902d1461218f5780635423112f146121485780635e3393841461212e578063655f2de1146120b357806367663bdf1461205c57806368a7c3fb14611f495780636b288d2014611efb5780636e91d13314611ed35780637596ed4314611e6e5780637a36a89114611e3757806380a1f71214611db3578063840059ec14611d6357806384b0196e14611c195780638537c27814611a8c57806386f4a07d14611a525780638c5b9b0014611a175780638eb8791a146119fa57806391d14854146119a55780639c73eeaf146119535780639c8d416f146119365780639fabc36f146118f4578063a217fddf146118b2578063a9ba045e146118cc578063aa3aa460146118b2578063ad3cb1cc1461186b578063af20c96014611833578063b7222ff5146117e0578063bcf7a3c9146110b9578063be1e372514611086578063c404c35914611045578063c415b95c1461101c578063d2c35ce814610fa9578063d547741f14610f59578063d72e957a14610f20578063dd39d76514610afa578063e3d72e5e14610aad578063eabcca101461031f578063ecd00261146102f7578063f624f88a146102ce5763f7d9cc1e0361001057346102cb57806003193601126102cb576020600854604051908152f35b80fd5b50346102cb57806003193601126102cb57600b546040516001600160a01b039091168152602090f35b50346102cb57806003193601126102cb5760206040515f5160206166055f395f51905f528152f35b50346102cb5760e03660031901126102cb576004356001600160a01b03811691908290036102cb576024356001600160a01b03811690819003610aa957610364612b3c565b6064356001600160a01b0381169291838203610aa557610382612b26565b9260c4356001600160a01b0381169190829003610aa1575f5160206166a55f395f51905f5254956001600160401b0360ff8860401c1615971680159081610a99575b6001149081610a8f575b159081610a86575b50610a77578660016001600160401b03195f5160206166a55f395f51905f525416175f5160206166a55f395f51905f5255610a47575b8815908115610a3e575b8115610a2c575b8115610a23575b508015610a1b575b610a0c57869761043a616350565b610442616350565b61044a616350565b60015f5160206166855f395f51905f5255878889610a01575b806305f5e100600a9210156109f0575b6127108110156109e1575b60648110156109d3575b10156109c8575b600a602160018b019a6104a18c612d79565b9b6104af6040519d8e612d58565b808d526104be601f1991612d79565b013660208e01378b01015b5f1901916f181899199a1a9b1b9c1cb0b131b232b360811b8282061a83530480156104f757600a90916104c9565b505060409788516105088a82612d58565b600481526356656c6160e01b6020820152610521616350565b610529616350565b80516001600160401b0381116109b4578b6105515f5160206165c55f395f51905f5254612e5c565b601f8111610977575b505060208c601f83116001146108fb579061058993836107ef575b50508160011b915f199060031b1c19161790565b5f5160206165c55f395f51905f52555b8051906001600160401b0382116108e7576105c15f5160206165e55f395f51905f5254612e5c565b601f8111610882575b506020908b601f84116001146107fa5792806106cd9795936106069361076e9b9a9896926107ef5750508160011b915f199060031b1c19161790565b5f5160206165e55f395f51905f52555b8a7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100558a7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10155610664616350565b600a600255600a600355600a6008556001600160601b0360a01b60095416176009556001600160601b0360a01b600a541617600a556001600160601b0360a01b600b541617600b5560018060a01b0381166001600160601b0360a01b601154161760115561476d565b506106d781614826565b505f5160206166055f395f51905f5286525f5160206166655f395f51905f526020525f5160206166455f395f51905f52600186882001545f5160206166055f395f51905f5288525f5160206166655f395f51905f52602052816001888a2001555f5160206166055f395f51905f527fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff8980a46146e8565b5060a4356010556001600160a01b0381166107df575b5061078d575080f35b60207fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29160ff60401b195f5160206166a55f395f51905f5254165f5160206166a55f395f51905f52555160018152a180f35b6107e8906148cc565b505f610784565b015190505f80610575565b505f5160206165e55f395f51905f528c52818c209190601f1984168d5b81811061086a575092600192859261076e9b9a98966106cd9a989610610852575b505050811b015f5160206165e55f395f51905f5255610616565b01515f1960f88460031b161c191690555f8080610838565b92936020600181928786015181550195019301610817565b5f5160206165e55f395f51905f528c526108d7907f5f9ce34815f8e11431c7bb75a8e6886a91478f7ffc1dbb0a98dc240fddd76b75601f850160051c810191602086106108dd575b601f0160051c0190615d22565b5f6105ca565b90915081906108ca565b634e487b7160e01b8b52604160045260248bfd5b9192601f1984165f5160206165c55f395f51905f528452828420935b81811061095f5750908460019594939210610947575b505050811b015f5160206165c55f395f51905f5255610599565b01515f1960f88460031b161c191690555f808061092d565b92936020600181928786015181550195019301610917565b6020825f5160206165c55f395f51905f526109ad945220601f840160051c810191602085106108dd57601f0160051c0190615d22565b8b5f61055a565b634e487b7160e01b8c52604160045260248cfd5b60019098019761048f565b6064600291049a0199610488565b612710600491049a019961047e565b6305f5e100600891049a0199610473565b506010985089610463565b632582a64160e11b8752600487fd5b50811561042c565b9050155f610424565b6001600160a01b03851615915061041d565b82159150610416565b600160401b60ff60401b195f5160206166a55f395f51905f525416175f5160206166a55f395f51905f525561040c565b63f92ee8a960e01b8852600488fd5b9050155f6103d6565b303b1591506103ce565b8891506103c4565b8680fd5b8480fd5b5080fd5b50346102cb5760203660031901126102cb57610ac7612afa565b610acf6145d7565b6001600160a01b03811615610aeb57610ae790614cba565b5080f35b632582a64160e11b8252600482fd5b50346102cb5760203660031901126102cb576004356001600160401b038111610aa95736602382011215610aa9578060040135906001600160401b038211610f1c5760248260051b82010190368211610edb57610b55614633565b610b5d614a16565b610b65615dbb565b610b6e83613162565b91610b7c6040519384612d58565b8383526024602084019201915b818310610efc575050509015610eea575b600b5460405163024ece8960e01b815291908390839060049082906001600160a01b03165afa918215610edf578392610e3d575b50805190825190845b838103610bf3578560015f5160206166855f395f51905f525580f35b6001600160401b03610c0582846133fd565b5116808752600e60205260408720878052602052604087205480610dd6575b50865b848103610d555750808752866020526040872054610c52575b90610c4c60019261546d565b01610bd7565b80875286602052866040812055600154875b818103610c7c575b5050600380546001019055610c40565b826001600160401b03610c8e83613462565b90549060031b1c1614610ca357600101610c64565b905f198101908111610d415790610cd56001600160401b03610cc7610cf494613462565b90549060031b1c1691613462565b9091906001600160401b038084549260031b9316831b921b1916179055565b6001548015610d2d5760019291610c4c915f1901610d1181613462565b6001600160401b0382549160031b1b1916905584559192610c6c565b634e487b7160e01b88526031600452602488fd5b634e487b7160e01b89526011600452602489fd5b818852600e60205260408820600191906001600160a01b03610d77838b6133fd565b5116838060a01b03165f5260205260405f205480610d97575b5001610c27565b80610db8610dd092338c610db187898060a01b03926133fd565b5116616012565b838060a01b03610dc8848c6133fd565b511685615377565b5f610d90565b8780808084335af1610de6614559565b5015610e2e57818852600e6020526040882088805260205260408820610e0d828254612fe5565b9055878052600f602052610e2660408920918254612fe5565b90555f610c24565b6312171d8360e31b8852600488fd5b9091503d8084833e610e4f8183612d58565b810190602081830312610edb578051906001600160401b038211610aa557019080601f83011215610edb578151610e8581613162565b92610e936040519485612d58565b81845260208085019260051b820101928311610ed757602001905b828210610ebf57505050905f610bce565b60208091610ecc84614545565b815201910190610eae565b8580fd5b8380fd5b6040513d85823e3d90fd5b506004610ef561305d565b9050610b9a565b82356001600160401b0381168103610aa157815260209283019201610b89565b8280fd5b50346102cb5760203660031901126102cb576020906040906001600160a01b03610f48612afa565b168152601283522054604051908152f35b50346102cb5760403660031901126102cb57610ae7600435610f79612b10565b90610fa4610f9f825f525f5160206166655f395f51905f52602052600160405f20015490565b6146a2565b614d4b565b50346102cb5760203660031901126102cb57610fc3612afa565b610fcb6145d7565b6001600160a01b03168015610aeb576020817fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f926001600160601b0360a01b6011541617601155604051908152a180f35b50346102cb57806003193601126102cb576011546040516001600160a01b039091168152602090f35b50346102cb57806003193601126102cb5761107561106161448e565b604051938493606085526060850190612bc3565b916020840152151560408301520390f35b50346102cb5760403660031901126102cb576110b56110a9602435600435615f1e565b60405191829182612c71565b0390f35b50610140366003190112611655576110cf612afa565b6024359060ff8216820361165557604435906001600160401b03821682036116555760056064351015611655576084356001600160401b0381116116555761111b903690600401612ce0565b909390919060a435906001600160a01b038216820361165557610104356001600160401b03811161165557611154903690600401612ce0565b90610124356001600160401b03811161165557611175903690600401612ce0565b91909260ff85166117d1576001600160401b0389165f525f60205260405f2054156117c2576111a2614a16565b5f60036064351415806117b3575b6117a45760e4354211611795576111c5614de7565b60085411156117865761177257879187868b8d60036064351461173a575b611349959361133a9361134096936001600160401b0361121c60429560018060a01b0386165f52601260205260405f20549c3691612d94565b602081519101209160ff6040519460208601967f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc1333340885260018060a01b0316604087015216606085015216608083015260ff6064351660a083015260c082015260018060a01b038d1660e082015260c4356101008201528861012082015260e43561014082015261014081526112b361016082612d58565b5190206112be6164d5565b6112c661653f565b6040519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f8452604083015260608201524660808201523060a082015260a0815261131760c082612d58565b519020906040519161190160f01b83526002830152602282015220923691612d94565b9061625e565b90929192616298565b6001600160a01b0316801561172b576001600160a01b0387160361171c5760018101809111611708576001600160a01b0386165f9081526012602052604090205560105434106116f95760c4351580806116e7575b6116d857156114bc575b50509060ff611462926020985061141b6113df6113c636898c612d94565b8b8151910120600754908660c435916064358d8c612ff2565b988997604051956113ef87612d0d565b4287526001600160a01b03168c87015260c4356040870152346060870152608086018990523691612d94565b60a08401526001600160a01b03851660c08401523360e08401526001600160401b0387166101008401521661012082015261145c6064356101408301612e94565b8361608b565b6040513381526001600160a01b03909116926001600160401b0316907fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236908690a460015f5160206166855f395f51905f5255604051908152f35b6001600160a01b038416156116d857600b5460405163cbe230c360e01b81526001600160a01b038681166004830152909160209183916024918391165afa90811561164a575f916116a9575b501561169a576060810361168b5781606091810103126116555780359060ff821680920361165557604051636eb1769f60e11b81526001600160a01b0386811660048301523060248301526020908290604490829089165afa90811561164a575f91611659575b5060c435116115a8575b50509060ff6114629260209861159260c4358786614b37565b61159f60c435858a614c65565b985091926113a8565b909291906001600160a01b0383163b156116555760409081519463d505accf60e01b865260018060a01b038716600487015230602487015260c435604487015260e43560648701526084860152602081013560a4860152013560c48401525f8360e4818360018060a01b0387165af190811561164a576020986114629460ff93611637575b5098509192611579565b61164391505f90612d58565b5f5f61162d565b6040513d5f823e3d90fd5b5f80fd5b90506020813d602011611683575b8161167460209383612d58565b8101031261165557515f61156f565b3d9150611667565b63ddafbaef60e01b5f5260045ffd5b63514e24c360e11b5f5260045ffd5b6116cb915060203d6020116116d1575b6116c38183612d58565b810190612de5565b5f611508565b503d6116b9565b632a9ffab760e21b5f5260045ffd5b506001600160a01b038516151561139e565b63c77f97eb60e01b5f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b632057875960e21b5f5260045ffd5b638baa579f60e01b5f5260045ffd5b50505050916085141580611767575b61175857879187868b8d6111e3565b637c6953f960e01b5f5260045ffd5b5060e2881415611749565b634e487b7160e01b5f52602160045260245ffd5b63d8219b3d60e01b5f5260045ffd5b631ab7da6b60e01b5f5260045ffd5b630ee2db0560e21b5f5260045ffd5b50505f600160643514156111b0565b6309ff9a8360e41b5f5260045ffd5b635428eae760e01b5f5260045ffd5b34611655576040366003190112611655576117f9612b66565b6001600160401b03611809612b10565b91165f52600e60205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b34611655576020366003190112611655576001600160a01b03611854612afa565b165f52600f602052602060405f2054604051908152f35b34611655575f366003190112611655576110b560405161188c604082612d58565b60058152640352e302e360dc1b6020820152604051918291602083526020830190612b92565b34611655575f3660031901126116555760206040515f8152f35b34611655575f36600319011261165557600a546040516001600160a01b039091168152602090f35b34611655576020366003190112611655576020600435611913816152bf565b908115611926575b506040519015158152f35b61193091506152e7565b8261191b565b34611655575f366003190112611655576020601054604051908152f35b346116555760203660031901126116555760043561196f6145d7565b80156116d8576020817e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db192600855604051908152a1005b34611655576040366003190112611655576119be612b10565b6004355f525f5160206166655f395f51905f5260205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b34611655575f366003190112611655576020600254604051908152f35b34611655575f36600319011261165557611a2f614633565b611a37614a16565b611a3f615dbb565b60015f5160206166855f395f51905f5255005b34611655575f3660031901126116555760206040517fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c959178152f35b346116555761018036600319011261165557611aa6612b66565b608435906001600160401b03821161165557604060031983360301126116555760a4356001600160401b03811161165557604060031982360301126116555760c435916001600160401b0383116116555736602384011215611655578260040135926001600160401b0384116116555736602460608602830101116116555761012435600d81101561165557610144356001600160401b03811161165557611b52903690600401612ce0565b939092610164356001600160401b03811161165557611b75903690600401612ce0565b335f9081527f212403f86ff12f15ffbeab879ed2f25237f7273bd0873f6abf8e9330ca4ed4b96020526040902054909891979060ff1615611be257611a3f99611bbc614a16565b6101043594602460e43595019260040191600401906064359060443590602435906134a2565b63e2517d3f60e01b5f52336004527f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9860245260445ffd5b34611655575f366003190112611655577fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100541580611d3a575b15611cfd57611ca1611c62615145565b611c6a615212565b6020611caf60405192611c7d8385612d58565b5f84525f368137604051958695600f60f81b875260e08588015260e0870190612b92565b908582036040870152612b92565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b828110611ce657505050500390f35b835185528695509381019392810192600101611cd7565b60405162461bcd60e51b81526020600482015260156024820152741152540dcc4c8e88155b9a5b9a5d1a585b1a5e9959605a1b6044820152606490fd5b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1015415611c52565b3461165557604036600319011261165557611d7c612afa565b611d84612b10565b6001600160a01b039182165f908152600c60209081526040808320949093168252928352819020549051908152f35b34611655575f36600319011261165557611dd5611dd06004614e16565b614588565b600654600754905f905b828110611df457604051806110b58682612c71565b60018181925f52600560205260405f20545f526004602052611e1860405f20612ea0565b611e2285886133fd565b52611e2d84876133fd565b5001910190611ddf565b34611655576020366003190112611655576001600160401b03611e58612b66565b165f525f602052602060405f2054604051908152f35b34611655575f36600319011261165557611e8661305d565b6040518091602082016020835281518091526020604084019201905f5b818110611eb1575050500390f35b82516001600160401b0316845285945060209384019390920191600101611ea3565b34611655575f366003190112611655576009546040516001600160a01b039091168152602090f35b3461165557602036600319011261165557611f14612afa565b6001600160a01b03165f9081525f5160206166c55f395f51905f52602090815260409182902054915160ff9092161515825290f35b606036600319011261165557611f5d612cd0565b6024356001600160401b03811161165557611f7c903690600401612ce0565b611f87929192612b3c565b9260ff83166117d157611fa192611f9c614a16565b614e2d565b6001600160a01b0382169182611fce575b60208260015f5160206166855f395f51905f5255604051908152f35b825f5260146020526001600160401b0360405f20541661204d573b1561203e5760c081901c5f81815260136020908152604080832080546001600160a01b0319166001600160a01b03881617905594825260148152939020805467ffffffffffffffff1916909117905582611fb2565b634ea29e6760e01b5f5260045ffd5b63d48c379760e01b5f5260045ffd5b346116555760e036600319011261165557612075612afa565b61207d612b7c565b906044356005811015611655576020926120ab92612099612b26565b9060c4359360a4359360643592612ff2565b604051908152f35b34611655576020366003190112611655576004356120cf6145d7565b80156116d8576002546120e460035482612fe5565b908183106116d8578261211e6040937f80af1d8add7b822532307ae28ef73c31ab2803f908c5c22cc879766f54c685839560025582612fe5565b60035582519182526020820152a1005b34611655575f3660031901126116555760206120ab614de7565b3461165557602036600319011261165557612161612e0a565b506004355f5260046020526110b561217b60405f20612ea0565b604051918291602083526020830190612bc3565b34611655575f366003190112611655577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031630036121e65760206040515f5160206166255f395f51905f528152f35b63703e46dd60e11b5f5260045ffd5b604036600319011261165557612209612afa565b6024356001600160401b03811161165557612228903690600401612dca565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163081149081156123a0575b506121e65761226a6145d7565b6040516352d1902d60e01b81526001600160a01b0383169290602081600481875afa5f918161236c575b506122ac5783634c9c8ce360e01b5f5260045260245ffd5b805f5160206166255f395f51905f5285920361235a5750813b15612348575f5160206166255f395f51905f5280546001600160a01b031916821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2815115612330575f8083602061001b95519101845af461232a614559565b91616477565b50503461233957005b63b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f5260045260245ffd5b632a87526960e21b5f5260045260245ffd5b9091506020813d602011612398575b8161238860209383612d58565b8101031261165557519085612294565b3d915061237b565b5f5160206166255f395f51905f52546001600160a01b0316141590508361225d565b6040366003190112611655576123d6612cd0565b6024356001600160401b038111611655576123f5903690600401612ce0565b60ff839293166117d15760209261240e92611f9c614a16565b60015f5160206166855f395f51905f5255604051908152f35b34611655575f3660031901126116555760206120ab614e03565b346116555760403660031901126116555761245a612b10565b336001600160a01b038216036124765761001b90600435614d4b565b63334bd91960e11b5f5260045ffd5b60e036600319011261165557612499612cd0565b6124a1612b7c565b604435906005821015611655576064356001600160401b038111611655576124cd903690600401612ce0565b93906124d7612b26565b60a4359460ff60c4359416806117d1576001600160401b03861695865f525f60205260405f2054156117c25761250b614a16565b5f9483158015612777575b6117a45760105487106116f95761252b614de7565b6008541115611786576001600160a01b03851695866126fa5761254e888b612dfd565b34036116d8575b61177257600384036126515760858a141580612646575b61175857602099612600976125ac6125d8948c61145c995b61258f828285614c65565b61259a368789612d94565b602081519101208a6007549433612ff2565b9a604051986125ba8a612d0d565b428a528d8a0152604089015260608801528960808801523691612d94565b60a08501523360c08501525f60e0850152856101008501526101208401526101408301612e94565b81604051915f83527fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236853394a460015f5160206166855f395f51905f5255604051908152f35b5060e28a141561256c565b92949591936002861461267b575b9260209961145c95936125ac6125d8948c6126009b9a98612584565b93919594926044602060018060a01b03600a5416604051928380926308b56a5360e11b82528d60048301523360248301525afa90811561164a575f916126db575b50156126cc57929495919361265f565b63622a850b60e11b5f5260045ffd5b6126f4915060203d6020116116d1576116c38183612d58565b8b6126bc565b8734036116d85789156116d857600b5460405163cbe230c360e01b81526004810189905290602090829060249082906001600160a01b03165afa90811561164a575f91612758575b501561169a576127538a3388614b37565b612555565b612771915060203d6020116116d1576116c38183612d58565b8c612742565b505f955060048414612516565b346116555760403660031901126116555761001b6004356127a3612b10565b906127c9610f9f825f525f5160206166655f395f51905f52602052600160405f20015490565b614985565b34611655575f3660031901126116555760206040517f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc13333408152f35b34611655575f3660031901126116555760206040515f5160206166455f395f51905f528152f35b346116555760203660031901126116555760206120ab6004355f525f5160206166655f395f51905f52602052600160405f20015490565b3461165557604036600319011261165557611a3f612882612afa565b61288a612b10565b90612893614a16565b614a4e565b34611655576020366003190112611655576128b1612afa565b6128b96145d7565b6001600160a01b038116156128d15761001b906146e8565b632582a64160e11b5f5260045ffd5b34611655575f366003190112611655576128fd611dd06015614e16565b601754601854905f905b82811061291c57604051806110b58682612c71565b60018181925f52601660205260405f20545f52601560205261294060405f20612ea0565b61294a85886133fd565b5261295584876133fd565b5001910190612907565b34611655576020366003190112611655576001600160401b03612980612b66565b165f526013602052602060018060a01b0360405f205416604051908152f35b34611655576020366003190112611655576001600160a01b036129c0612afa565b165f52600d602052602060405f2054604051908152f35b34611655576020366003190112611655576001600160a01b036129f8612afa565b165f52601460205260206001600160401b0360405f205416604051908152f35b34611655575f366003190112611655576020600354604051908152f35b34611655575f3660031901126116555760206040517f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd988152f35b34611655576020366003190112611655576001600160a01b03612a90612afa565b165f526012602052602060405f2054604051908152f35b34611655576020366003190112611655576004359063ffffffff60e01b821680920361165557602091637965db0b60e01b8114908115612ae9575b5015158152f35b6301ffc9a760e01b14905083612ae2565b600435906001600160a01b038216820361165557565b602435906001600160a01b038216820361165557565b608435906001600160a01b038216820361165557565b604435906001600160a01b038216820361165557565b35906001600160a01b038216820361165557565b600435906001600160401b038216820361165557565b602435906001600160401b038216820361165557565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9060058210156117725752565b90612c6e908251815260018060a01b03602084015116602082015260408301516040820152606083015160608201526080830151608082015261014080612c1b60a086015161016060a0860152610160850190612b92565b9460018060a01b0360c08201511660c085015260018060a01b0360e08201511660e08501526001600160401b036101008201511661010085015260ff610120820151166101208501520151910190612bb6565b90565b602081016020825282518091526040820191602060408360051b8301019401925f915b838310612ca357505050505090565b9091929394602080612cc1600193603f198682030187528951612bc3565b97019301930191939290612c94565b6004359060ff8216820361165557565b9181601f84011215611655578235916001600160401b038311611655576020838186019501011161165557565b61016081019081106001600160401b03821117612d2957604052565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117612d2957604052565b90601f801991011681019081106001600160401b03821117612d2957604052565b6001600160401b038111612d2957601f01601f191660200190565b929192612da082612d79565b91612dae6040519384612d58565b829481845281830111611655578281602093845f960137010152565b9080601f8301121561165557816020612c6e93359101612d94565b90816020910312611655575180151581036116555790565b9190820180921161170857565b60405190612e1782612d0d565b5f61014083828152826020820152826040820152826060820152826080820152606060a08201528260c08201528260e082015282610100820152826101208201520152565b90600182811c92168015612e8a575b6020831014612e7657565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612e6b565b60058210156117725752565b90604051612ead81612d0d565b80928054825260018060a01b03600182015416602083015260028101546040830152600381015460608301526004810154608083015260058101604051905f90805490612ef982612e5c565b8085529160018116908115612fbd5750600114612f7f575b50506007612f7d949383612f2c6101409560ff950382612d58565b60a086015260018060a01b0360068201541660c0860152015460018060a01b03811660e08501526001600160401b038160a01c16610100850152818160e01c1661012085015260e81c169101612e94565b565b5f908152602081209092505b818310612fa357505081016020016007612f7d612f11565b600181602092949394548385880101520191019190612f8b565b60ff191660208087019190915292151560051b8501909201925060079150612f7d9050612f11565b9190820391821161170857565b95939161302890959391956001600160401b0360405197602089019960018060a01b03168a521660408801526060870190612bb6565b60808501526001600160a01b031660a084015260c083015260e080830191909152815261305761010082612d58565b51902090565b60405190600154808352826020810160015f5260205f20925f905b80600383011061311057612f7d9454918181106130f6575b8181106130d9575b8181106130bc575b106130ae575b500383612d58565b60c01c81526020015f6130a6565b9260206001916001600160401b038560801c1681520193016130a0565b9260206001916001600160401b038560401c168152019301613098565b9260206001916001600160401b0385168152019301613090565b916004919350608060019186546001600160401b03811682526001600160401b038160401c1660208301526001600160401b0381841c16604083015260c01c6060820152019401920185929391613078565b6001600160401b038111612d295760051b60200190565b9061318382613162565b6131906040519182612d58565b82815280926131a1601f1991613162565b0190602036910137565b903590601e198136030182121561165557018035906001600160401b03821161165557602001918160051b3603831361165557565b9190604083820312611655576040516131f881612d3d565b809380356001600160401b03811161165557810183601f8201121561165557803561322281613162565b916132306040519384612d58565b81835260208084019260051b820101918683116116555760208201905b8382106132cf575050505082526020810135906001600160401b03821161165557019180601f8401121561165557823561328681613162565b936132946040519586612d58565b81855260208086019260051b82010192831161165557602001905b8282106132bf5750505060200152565b81358152602091820191016132af565b81356001600160401b038111611655576020916132f18a848094880101612dca565b81520191019061324d565b60408201908051916040845282518091526060840190602060608260051b8701019401915f905b82821061336d575050505060200151916020818303910152602080835192838152019201905f5b8181106133575750505090565b825184526020938401939092019160010161334a565b9091929460208061338a600193605f198b82030186528951612b92565b970192019201909291613323565b90600d8210156117725752565b908060209392818452848401375f828201840152601f01601f1916010190565b91908110156133d5576060020190565b634e487b7160e01b5f52603260045260245ffd5b356001600160a01b03811681036116555790565b80518210156133d55760209160051b010190565b91908110156133d55760051b0190565b91908110156133d55760051b81013590601e19813603018212156116555701908135916001600160401b038311611655576020018236038113611655579190565b906001548210156133d55760015f52600282901c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6019160031b60181690565b9d9a9d9c9b9795939998969492919c61016052610180526101a05261012052610100526101405260a05260c0526134db6101a0516152bf565b801561447c575b1561446d576134f36101a0516152bf565b908115614458576101a0515f52601560205260405f20955b60078701546080526001600160401b0360805160a01c166001600160401b036101605116036117c2576001600160401b0361016051165f525f60205260405f205461018051036140485761355f85806131ab565b60e0525061357060208601866131ab565b905060e051036117585761358a61012051610120516131ab565b9690506135a060206101205101610120516131ab565b9050870361175857604051916135b583612d0d565b6001600160401b03610160511683526101805160208401528960408401526101a05160608401526135e636886131e0565b60808401526135f836610120516131e0565b60a084015261360961014051613162565b6136166040519182612d58565b6101405181526020810136606061014051026101005101116116555761010051905b60606101405102610100510182106143fb57505060c084015260a05160e084015260c051610100840152600d8610156117725790918561012082015261367f36858d612d94565b61014082015260018060a01b036009541691604051938492635c9626b960e01b8452604060048501526001600160401b0381511660448501526020810151606485015260408101516084850152606081015160a485015261370a6136f4608083015161016060c48801526101a48701906132fc565b60a08301518682036043190160e48801526132fc565b9060c08101519160431986820301610104870152602080845192838152019301905f5b8181106143bb5750505060209593859361378885946101408560e061379a9701516101248901526101008101516101448901526137746101208201516101648a0190613398565b015186820360431901610184880152612b92565b848103600319016024860152916133a5565b03915afa90811561164a575f9161439c575b501561172b57600386015460068701546080516001600160a01b039182169a9295929116801561439557995b8261405757505050506001600160401b0361016051165f525f6020528560405f2054146140485715614024575b505f5f9061381561014051613179565b61382161014051613179565b905f925b610140518110613ea057505f5b838110613de657505050506138549061384f60c05160a051612dfd565b612dfd565b61385c61530f565b10613dd7575f5b60e0518103613d625750505f5b818103613cd5575050600760ff91015460e81c1691600583109283156117725760028114613c9b575b6001600160401b0361016051165f525f6020528260405f2055159182613c53575b6040519061018051825260208201526101a051907f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e60406001600160401b03610160511692a360a051613c02575b505f5f6001600160401b0361016051165f52601360205260018060a01b0360405f20541661393861014051613179565b92811515905b610140518110613a935750505061395481613179565b915f5b828103613a6b5750505061397690610120516101a051610160516156d8565b60405191613985602084612d58565b5f83526139946101a0516152bf565b15613a5e576139a16163b8565b156117725715613a0f57604051907f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d6101a05192806139f26001600160401b036101605116945f8060c051856154cc565b0390a35b60c051601154612f7d91906001600160a01b03166153c8565b604051907f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed666101a0519280613a566001600160401b036101605116945f8060c051856154cc565b0390a36139f6565b613a6661637b565b6139a1565b6001906001600160a01b03613a8082856133fd565b5116613a8c82876133fd565b5201613957565b8180613bdc575b613ba6575b80613afe613abe613ab960019461014051610100516133c5565b6133e9565b613ada6020613ad48561014051610100516133c5565b016133e9565b906040613aee8561014051610100516133c5565b013591858060a01b031690615423565b613b146020613ad48361014051610100516133c5565b613b28613ab98361014051610100516133c5565b6040613b3b8461014051610100516133c5565b0135907faafeb614695d37710d00f37c69671d205eb5b6fcd02995f6c3c7e2c93d3cbb2460405193868060a01b03169380613b9d6101a051956001600160401b036101605116958360209093929193604081019460018060a01b031681520152565b0390a40161393e565b9260018091613bbf613ab98761014051610100516133c5565b613bc982896133fd565b90838060a01b0316905201939050613a9f565b508260018060a01b03613bfb6020613ad48561014051610100516133c5565b1614613a9a565b60018060a01b0316613c1660a051826153c8565b604051604081015f825260a05160208301525f5160206165a55f395f51905f526101a05192806001600160401b036101605116930390a45f613908565b600154600160401b811015612d2957613c75816001613c969301600155613462565b61016051906001600160401b038084549260031b9316831b921b1916179055565b6138ba565b6101a0516001600160401b0361016051167ff36f013d86705b333febcccf7ece97eee30f4cd398a9c69aee1f940ed5daddea5f80a3613899565b80613cf5600192613cef60206101205101610120516131ab565b90613411565b35613d1082613d0a61012051610120516131ab565b90613421565b7f485d4bfb37695319d1b2352eef5982d10db8c284b0a7ddefe4465377c9fbc8df6040516020815280613d596101a051956001600160401b0361016051169560208401916133a5565b0390a401613870565b80613d76600192613cef60208601866131ab565b35613d8582613d0a86806131ab565b7f9ede1ad14e46a641d7b174bb6c9d939ad961d1a7de4b669163b9714a67e223206040516020815280613dce6101a051956001600160401b0361016051169560208401916133a5565b0390a401613863565b631e9acf1760e31b5f5260045ffd5b6001600160a01b03613df882846133fd565b516040516370a0823160e01b81523060048201529116602082602481845afa91821561164a575f92613e6b575b50613e4d81613e5e925f52600f60205260405f2054905f52600d60205260405f205490612dfd565b613e5784876133fd565b5190612dfd565b11613dd757600101613832565b9091506020813d8211613e98575b81613e8660209383612d58565b81010312611655575190613e4d613e25565b3d9150613e79565b93613eb5613ab98661014051610100516133c5565b906040613ec98761014051610100516133c5565b0135916001600160401b0361016051165f52600e60205260405f2060018060a01b0382165f5260205260405f2054831161401557610160516001600160401b03165f908152600e602090815260408083206001600160a01b0385168452825280832080548790039055600f9091529020548311613dd7576001600160a01b0381165f818152600f6020526040902080548590039055613f775750600191613f6f91612dfd565b945b01613825565b909591905f805b878110613fc7575b5015613f97575b5050600190613f71565b91946001928392906001600160a01b0316613fb283876133fd565b52613fbd82876133fd565b520193905f613f8d565b6001600160a01b0383811690613fdd83896133fd565b511614613fec57600101613f7e565b905061400c61400584613fff848a6133fd565b51612dfd565b91876133fd565b5260015f613f86565b6306a1762560e11b5f5260045ffd5b61403260c05160a051612dfd565b036116d85760105460c051106116d8575f613805565b630b6fac0360e41b5f5260045ffd5b94939690955060e09992999891979851159081159161438b575b50801561437f575b611758576001600160401b0361016051165f525f60205260405f2054036140485760028701549060018060a01b03600189015416926001600160401b0361016051165f52600e60205260405f2060018060a01b0385165f5260205260405f20548311614015576140e761530f565b10613dd7576140fa828461016051615377565b60038801546010548082106143765761411291612fe5565b925b8061430b575081156142a957509061412b91612dfd565b61413581836153c8565b6040519060408201905f835260208301525f5160206165a55f395f51905f526101a05192806001600160401b036101605116930390a45b60ff600785015460e81c16916005831015611772576141a5921561428f575b156142825760ff60075f955b015460e81c16943691612d94565b906141b26101a0516152bf565b15614275576141bf6163b8565b600584101561177257612f7d93614229577f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d604051806142146101a051956001600160401b03610160511695600189856154cc565b0390a35b6011546001600160a01b03166153c8565b7f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed666040518061426d6101a051956001600160401b03610160511695600189856154cc565b0390a3614218565b61427d61637b565b6141bf565b60ff600760105495614197565b6001600354016003556142a46101605161546d565b61418b565b91925050816142ba575b505061416c565b6001600160a01b0316906142ce81836153c8565b6040519060408201905f835260208301525f5160206165a55f395f51905f526101a05192806001600160401b036101605116930390a45f806142b3565b8291939492614324575b505050816142ba57505061416c565b61432f828483615423565b604080516101a051610160516001600160a01b0394909416825260208201949094526001600160401b03909216915f5160206165a55f395f51905f529190a45f8080614315565b50505f92614114565b50610140511515614079565b905015155f614071565b50896137d8565b6143b5915060203d6020116116d1576116c38183612d58565b5f6137ac565b825180516001600160a01b03908116875260208281015190911681880152604091820151918701919091528a98506060909501949092019160010161372d565b606082360312611655576040519060608201908282106001600160401b03831117612d295760609260209260405261443285612b52565b815261443f838601612b52565b8382015260408501356040820152815201910190613638565b6101a0515f52600460205260405f209561350b565b6302e8145360e61b5f5260045ffd5b506144896101a0516152e7565b6134e2565b614496612e0a565b5061449f614e03565b614505576144ab614de7565b6144be576144b7612e0a565b905f905f90565b6006545f52600560205260405f20545f52600460205260405f20906001600160401b03600783015460a01c165f525f6020526144fe60405f205492612ea0565b9190600190565b6017545f52601660205260405f20545f52601560205260405f20906001600160401b03600783015460a01c165f525f6020526144fe60405f205492612ea0565b51906001600160a01b038216820361165557565b3d15614583573d9061456a82612d79565b916145786040519384612d58565b82523d5f602084013e565b606090565b9061459282613162565b61459f6040519182612d58565b82815280926145b0601f1991613162565b01905f5b8281106145c057505050565b6020906145cb612e0a565b828285010152016145b4565b335f9081527f78e571b7bf30584d955e1c6444a2b5147087edf9f00485d94993a04d370525ea602052604090205460ff161561460f57565b63e2517d3f60e01b5f52336004525f5160206166455f395f51905f5260245260445ffd5b335f9081527f143ba17de5dacef81d11036b41c596cd787c462f0e59fb1f76aa2b829d603e0c602052604090205460ff161561466b57565b63e2517d3f60e01b5f52336004527fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c9591760245260445ffd5b5f8181525f5160206166655f395f51905f526020908152604080832033845290915290205460ff16156146d25750565b63e2517d3f60e01b5f523360045260245260445ffd5b6001600160a01b0381165f9081525f5160206166c55f395f51905f52602052604090205460ff16614768576001600160a01b03165f8181525f5160206166c55f395f51905f5260205260408120805460ff191660011790553391905f5160206166055f395f51905f52905f5160206165855f395f51905f529080a4600190565b505f90565b6001600160a01b0381165f9081527f212403f86ff12f15ffbeab879ed2f25237f7273bd0873f6abf8e9330ca4ed4b9602052604090205460ff16614768576001600160a01b03165f8181527f212403f86ff12f15ffbeab879ed2f25237f7273bd0873f6abf8e9330ca4ed4b960205260408120805460ff191660011790553391907f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd98905f5160206165855f395f51905f529080a4600190565b6001600160a01b0381165f9081527f78e571b7bf30584d955e1c6444a2b5147087edf9f00485d94993a04d370525ea602052604090205460ff16614768576001600160a01b03165f8181527f78e571b7bf30584d955e1c6444a2b5147087edf9f00485d94993a04d370525ea60205260408120805460ff191660011790553391905f5160206166455f395f51905f52905f5160206165855f395f51905f529080a4600190565b6001600160a01b0381165f9081527f143ba17de5dacef81d11036b41c596cd787c462f0e59fb1f76aa2b829d603e0c602052604090205460ff16614768576001600160a01b03165f8181527f143ba17de5dacef81d11036b41c596cd787c462f0e59fb1f76aa2b829d603e0c60205260408120805460ff191660011790553391907fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c95917905f5160206165855f395f51905f529080a4600190565b5f8181525f5160206166655f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff16614a10575f8181525f5160206166655f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f5160206165855f395f51905f529080a4600190565b50505f90565b60025f5160206166855f395f51905f525414614a3f5760025f5160206166855f395f51905f5255565b633ee5aeb560e01b5f5260045ffd5b60018060a01b03811691825f52600c60205260405f2060018060a01b0382165f5260205260405f2054918215614b31577f16d46b28f1a9377c2df78652a0a8e31568b99311665bcec24cbba72edd2e5c2c8392855f52600c60205260405f2060018060a01b0382165f526020525f6040812055855f52600d60205260405f20614ad8858254612fe5565b9055604080516001600160a01b0394851681526020810195909552921692839290a282614b28575f809350809281925af1614b11614559565b5015614b1957565b6312171d8360e31b5f5260045ffd5b612f7d92616012565b50505050565b6040516370a0823160e01b81523060048201526001600160a01b039190911691602082602481865afa91821561164a575f92614c30575b50602492614bb5602092604051906323b872dd60e01b8583015260018060a01b03168682015230604482015286606482015260648152614baf608482612d58565b826162f8565b6040516370a0823160e01b815230600482015293849182905afa801561164a575f90614bfc575b614be69250612fe5565b03614bed57565b63b0a6ea2960e01b5f5260045ffd5b506020823d602011614c28575b81614c1660209383612d58565b8101031261165557614be69151614bdc565b3d9150614c09565b9091506020813d602011614c5d575b81614c4c60209383612d58565b810103126116555751906024614b6e565b3d9150614c3f565b6001600160401b03165f52600e60205260405f2060018060a01b0382165f5260205260405f20614c96838254612dfd565b905560018060a01b03165f52600f602052614cb660405f20918254612dfd565b9055565b6001600160a01b0381165f9081525f5160206166c55f395f51905f52602052604090205460ff1615614768576001600160a01b03165f8181525f5160206166c55f395f51905f5260205260408120805460ff191690553391905f5160206166055f395f51905f52907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b5f8181525f5160206166655f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff1615614a10575f8181525f5160206166655f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b600754600654808211614dfa5750505f90565b612c6e91612fe5565b601854601754808211614dfa5750505f90565b60026003820154910154808211614dfa5750505f90565b335f9081525f5160206166c55f395f51905f5260205260409020549092919060ff161561513657600354801561512757614e65614de7565b60085411156117865760105434106116f9575f1901600355614e88368383612d94565b805160209182012060075460408051339481019485525f91810191909152915f606084015260808301525f60a08301525f60c083015260e082015260e08152614ed361010082612d58565b519020928360c01c9260405190614ee982612d0d565b42825260208201915f835260408101915f8352614f16606083019634885260808401928a84523691612d94565b60a0830190815260c083019133835260e08401975f89526101008501958a875260ff6101208701991689526101408601975f89528c5f52600460205260405f209651875560018060a01b03905116600187019060018060a01b03166001600160601b0360a01b8254161790555160028601555160038501555160048401556005830190518051906001600160401b038211612d2957614fbf82614fb98554612e5c565b85616053565b602090601f83116001146150bf579180614ff3926007979695945f926107ef5750508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b169160058110156117725760ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b191617171790556007545f5260056020528160405f2055600160075401600755604051908282527fbbcfb0c180d5f95cfb50c12005c44144ac92ac8474f1abc5add0e6e69a14e8dc60203393a390565b90601f19831691845f52815f20925f5b81811061510f575091600193918560079998979694106150f7575b505050811b019055614ff6565b01515f1960f88460031b161c191690555f80806150ea565b929360206001819287860151815501950193016150cf565b6316ae067f60e11b5f5260045ffd5b63fc336c4160e01b5f5260045ffd5b604051905f825f5160206165c55f395f51905f52549161516483612e5c565b80835292600181169081156151f35750600114615188575b612f7d92500383612d58565b505f5160206165c55f395f51905f525f90815290917f42ad5d3e1f2e6e70edcf6d991b8a3023d3fca8047a131592f9edb9fd9b89d57d5b8183106151d7575050906020612f7d9282010161517c565b60209193508060019154838589010152019101909184926151bf565b60209250612f7d94915060ff191682840152151560051b82010161517c565b604051905f825f5160206165e55f395f51905f52549161523183612e5c565b80835292600181169081156151f3575060011461525457612f7d92500383612d58565b505f5160206165e55f395f51905f525f90815290917f5f9ce34815f8e11431c7bb75a8e6886a91478f7ffc1dbb0a98dc240fddd76b755b8183106152a3575050906020612f7d9282010161517c565b602091935080600191548385890101520191019091849261528b565b6018549060175480921191826152d457505090565b9091505f52601660205260405f20541490565b6007549060065480921191826152fc57505090565b9091505f52600560205260405f20541490565b5f8052600d6020527f81955a0a11e65eac625c29e8882660bae4e165a75d72780094acae8ece9a29ee54612c6e906153479047612fe5565b5f8052600f6020527ff4803e074bd026baaf6ed2e288c9515f68c72fb7216eebdd7cae1718a53ec3755490612fe5565b6001600160401b03165f52600e60205260405f2060018060a01b0382165f5260205260405f206153a8838254612fe5565b905560018060a01b03165f52600f602052614cb660405f20918254612fe5565b6001600160a01b03165f9081527f13649b2456f1b42fef0f0040b3aaeabcd21a76a0f3f5defd4f583839455116e860205260409020805461540a908390612dfd565b90555f8052600d602052614cb660405f20918254612dfd565b60018060a01b031690815f52600c60205260405f209060018060a01b03165f5260205260405f20615455838254612dfd565b90555f52600d602052614cb660405f20918254612dfd565b6001600160401b03165f818152601360205260409020546001600160a01b031680615496575050565b5f52601460205260405f206001600160401b031981541690555f52601360205260405f206001600160601b0360a01b8154169055565b909392919381526002841015611772576154f5608092612c6e9560208401526040830190613398565b8160608201520190612b92565b9035601e19823603018112156116555701602081359101916001600160401b038211611655578160051b3603831361165557565b9060408101916155468180615502565b809460408552526060830160608560051b85010194825f90601e19813603015b8383106155b157505050505050806020615581920190615502565b82840360209093019290925281835290916001600160fb1b0383116116555760209260051b809284830137010190565b909192939497605f1988820301865288358281121561165557830190602082359201916001600160401b038111611655578036038313611655576155fb60209283926001956133a5565b9a0196019493019190615566565b81601f820112156116555780519061562082613162565b9261562e6040519485612d58565b82845260208085019360061b8301019181831161165557602001925b828410615658575050505090565b604084830312611655576020604091825161567281612d3d565b61567b87614545565b8152828701518382015281520193019261564a565b90602080835192838152019201905f5b8181106156ad5750505090565b825180516001600160a01b0316855260209081015181860152604090940193909201916001016156a0565b91905f906001600160401b03841692835f52601360205260018060a01b0360405f205416958615615d19578051905f5b828103615cf3575050506001863b15611655576040516338d9814d60e11b8152602060048201525f818061573f6024820187615536565b0381838c5af19081615cde575b50615cd95750825b82857fca435620e7fa2ff484abda442975d2167d5e6577548dfdfac799cecc2a2af2c9602060405194151594858152a3604051633ccfd60b60e01b815260609283916001908781600481838f5af1908189918a93615c6f575b50615c635750505085975b845190875b828103615c2c5750505061581f866060998b6157fe8360019661583160405197889687958694635205569560e11b865260a0600487015260a4860190615536565b9c602485015215159b8c60448501526003198482030160648501528d615690565b8281036003190160848401528a615690565b03925af1879181615ba9575b50615b74575050917fa59edce59a4211e20e9f86e63adea41a55dbe409040bc17f473b671e31bd036a916158a0879461589288945b604051958695865215156020860152608060408601526080850190615690565b908382036060850152615690565b0390a38251806158b2575b5050505050565b6020840120601854604051906020820192878452856040840152600560041015611772576004606084015260808301528360a08301528360c083015260e082015260e0815261590361010082612d58565b5190209260405161591381612d0d565b4281526020810183815260408201918483526060810190858252608081019488865260a0820190815260c08201908a825260e08301968888526101008401948a86526101208501978a8952610140860197600489528d8c52601560205260408c209651875560018060a01b03905116600187019060018060a01b03166001600160601b0360a01b8254161790555160028601555160038501555160048401556005830190518051906001600160401b038211615b60576159d782614fb98554612e5c565b602090601f8311600114615af9579180615a0b926007979695948d926107ef5750508160011b915f199060031b1c19161790565b90555b516006820180546001600160a01b03199081166001600160a01b0393841617909155965192909101805490961691161780855590519251915160e09290921b60ff60e01b16916005811015615ae55769ffffffffffffffffffff60a01b1990911660a09390931b67ffffffffffffffff60a01b16929092171760e89190911b60ff60e81b16179055601880548252601660209081526040808420869055825460010190925590519182527fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d54923691a45f808080806158ab565b634e487b7160e01b86526021600452602486fd5b838b52818b209190601f1984168c5b818110615b4857509160019391856007999897969410615b30575b505050811b019055615a0e565b01515f1960f88460031b161c191690555f8080615b23565b92936020600181928786015181550195019301615b08565b634e487b7160e01b8a52604160045260248afd5b879498507fa59edce59a4211e20e9f86e63adea41a55dbe409040bc17f473b671e31bd036a9391926158926158a0929a615872565b9091503d8089833e615bbb8183612d58565b810190602081830312615c28578051906001600160401b038211615c24570181601f82011215615c2857805190615bf182612d79565b92615bff6040519485612d58565b82845260208383010111615c2457818a9260208093018386015e83010152905f61583d565b8980fd5b8880fd5b600190615c5d6001600160a01b03615c44838b6133fd565b5151166020615c53848c6133fd565b5101519085614c65565b016157bd565b919991955092506157b8565b915091503d808a833e615c828183612d58565b8101604082820312615c245781516001600160401b038111615cd55781615caa918401615609565b916020810151906001600160401b038211615cd157615cca929101615609565b915f6157ad565b8b80fd5b8a80fd5b615754565b615ceb9195505f90612d58565b5f935f61574c565b600190615d138a6001600160a01b03615d0c84876133fd565b5116614a4e565b01615708565b50505050505050565b818110615d2d575050565b5f8155600101615d22565b60075f9182815582600182015582600282015582600382015582600482015560058101615d658154612e5c565b9081615d78575b50508260068201550155565b81601f869311600114615d8f5750555b5f80615d6c565b81835260208320615dab91601f0160051c810190600101615d22565b8082528160208120915555615d88565b6006546007545f905b808303615e2b57509050600654600755601754601854905b818103615df457505060175460185560035401600355565b806001915f52601660205260405f20545f526015602052615e1760405f20615d38565b805f5260166020525f604081205501615ddc565b825f52600560205260405f20545f52600460205260405f20926007840180549460ff8660e81c1660058110156117725715615efa575b60019394955060028101805480615ead575b50505050805f52600560205260405f20545f526004602052615e9760405f20615d38565b5f81815260056020526040812055019190615dc4565b615ed56001600160401b03615ef1955460a01c169187850192888060a01b0384541690615377565b5460069092015490549160a086901b8690039182169116615423565b5f808080615e73565b615f156001600160401b036001809697019760a01c1661546d565b85949350615e61565b919091615f2b6004614e16565b9283821080159061600a575b615fce57615f459082612dfd565b92808411615fc6575b50615f73615f5f611dd08386612fe5565b93615f6d6006549384612dfd565b92612dfd565b905f905b828110615f8357505050565b60018181925f52600560205260405f20545f526004602052615fa760405f20612ea0565b615fb185896133fd565b52615fbc84886133fd565b5001910190615f77565b92505f615f4e565b5090915050604051615fe1602082612d58565b5f81525f805b818110615ff357505090565b602090615ffe612e0a565b82828601015201615fe7565b508015615f37565b60405163a9059cbb60e01b60208201526001600160a01b03929092166024830152604480830193909352918152612f7d9161604e606483612d58565b6162f8565b9190601f811161606257505050565b612f7d925f5260205f20906020601f840160051c830193106108dd57601f0160051c0190615d22565b90815f52600460205260405f208151815560018060a01b03602083015116600182019060018060a01b03166001600160601b0360a01b8254161790556040820151600282015560608201516003820155608082015160048201556005810160a08301518051906001600160401b038211612d295761610d82614fb98554612e5c565b602090601f83116001146161f75791806161409260079695945f926107ef5750508160011b915f199060031b1c19161790565b90555b60c08301516006820180546001600160a01b03199081166001600160a01b039384161790915560e08086015194909301805490911693909116929092178083556101008401516101208501516101409095015193949093921b60ff60e01b169160058110156117725760ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b191617171790556007545f52600560205260405f2055600160075401600755565b90601f19831691845f52815f20925f5b81811061624657509160019391856007989796941061622e575b505050811b019055616143565b01515f1960f88460031b161c191690555f8080616221565b92936020600181928786015181550195019301616207565b815191906041830361628e576162879250602082015190606060408401519301515f1a906163f5565b9192909190565b50505f9160029190565b600481101561177257806162aa575050565b600181036162c15763f645eedf60e01b5f5260045ffd5b600281036162dc575063fce698f760e01b5f5260045260245ffd5b6003146162e65750565b6335e2f38360e21b5f5260045260245ffd5b905f602091828151910182855af11561164a575f513d61634757506001600160a01b0381163b155b6163275750565b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b60011415616320565b60ff5f5160206166a55f395f51905f525460401c161561636c57565b631afcd79f60e31b5f5260045ffd5b6006545f52600560205260405f20545f52600460205261639d60405f20615d38565b6006545f5260056020525f6040812055600160065401600655565b6017545f52601660205260405f20545f5260156020526163da60405f20615d38565b6017545f5260166020525f6040812055600160175401601755565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841161646c579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa1561164a575f516001600160a01b0381161561646257905f905f90565b505f906001905f90565b5050505f9160039190565b9061649b575080511561648c57805190602001fd5b63d6bda27560e01b5f5260045ffd5b815115806164cc575b6164ac575090565b639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b50803b156164a4565b6164dd615145565b80519081156164ed576020012090565b50507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10054801561651a5790565b507fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47090565b616547615212565b8051908115616557576020012090565b50507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10154801561651a579056fe2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0dc9961aaccfc3406b40d0d25a5b83b8a4721e5f454116d7e22bd8038c715140aaa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d102a16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d103fc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbcdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4202dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268009b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008f7e9a8734279607aeaac474f24949ddc99eae25c0da72edcf5ad21022be2c9ea264697066735822122072007ef67940846c06bfa5cbcee33d3c268f5e04f842c3e9abb5e1fec323806c64736f6c634300081e0033f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
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

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (processorEndpoint *ProcessorEndpoint) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := processorEndpoint.abi.Pack("UPGRADE_INTERFACE_VERSION")
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
func (processorEndpoint *ProcessorEndpoint) TryPackUPGRADEINTERFACEVERSION() ([]byte, error) {
	return processorEndpoint.abi.Pack("UPGRADE_INTERFACE_VERSION")
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (processorEndpoint *ProcessorEndpoint) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := processorEndpoint.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
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
// the contract method with ID 0x67663bdf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes32 payloadHash, address tokenAddress, uint256 assetAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payloadHash [32]byte, tokenAddress common.Address, assetAmount *big.Int, idx *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payloadHash, tokenAddress, assetAmount, idx)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x67663bdf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes32 payloadHash, address tokenAddress, uint256 assetAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payloadHash [32]byte, tokenAddress common.Address, assetAmount *big.Int, idx *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payloadHash, tokenAddress, assetAmount, idx)
}

// UnpackGenerateRequestId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x67663bdf.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes32 payloadHash, address tokenAddress, uint256 assetAmount, uint256 idx) pure returns(bytes32)
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

// PackGetTriggerRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f56b792.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getTriggerRequests() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) PackGetTriggerRequests() []byte {
	enc, err := processorEndpoint.abi.Pack("getTriggerRequests")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetTriggerRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f56b792.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getTriggerRequests() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetTriggerRequests() ([]byte, error) {
	return processorEndpoint.abi.Pack("getTriggerRequests")
}

// UnpackGetTriggerRequests is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1f56b792.
//
// Solidity: function getTriggerRequests() view returns((uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) UnpackGetTriggerRequests(data []byte) ([]StructsPendingRequest, error) {
	out, err := processorEndpoint.abi.Unpack("getTriggerRequests", data)
	if err != nil {
		return *new([]StructsPendingRequest), err
	}
	out0 := *abi.ConvertType(out[0], new([]StructsPendingRequest)).(*[]StructsPendingRequest)
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

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xeabcca10.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address _teeAuthenticator, address _authorityRegistry, address updateStatusOperator, address admin, address resetOperator, uint256 _minFeePerRequest, address _tokenAllowlist) returns()
func (processorEndpoint *ProcessorEndpoint) PackInitialize(teeAuthenticator common.Address, authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address, resetOperator common.Address, minFeePerRequest *big.Int, tokenAllowlist common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("initialize", teeAuthenticator, authorityRegistry, updateStatusOperator, admin, resetOperator, minFeePerRequest, tokenAllowlist)
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
func (processorEndpoint *ProcessorEndpoint) TryPackInitialize(teeAuthenticator common.Address, authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address, resetOperator common.Address, minFeePerRequest *big.Int, tokenAllowlist common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("initialize", teeAuthenticator, authorityRegistry, updateStatusOperator, admin, resetOperator, minFeePerRequest, tokenAllowlist)
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

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackProxiableUUID() []byte {
	enc, err := processorEndpoint.abi.Pack("proxiableUUID")
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
func (processorEndpoint *ProcessorEndpoint) TryPackProxiableUUID() ([]byte, error) {
	return processorEndpoint.abi.Pack("proxiableUUID")
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
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

// PackSubmitDeployRequestWithTrigger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68a7c3fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitDeployRequestWithTrigger(uint8 protocolVersion, bytes payload, address trigger) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitDeployRequestWithTrigger(protocolVersion uint8, payload []byte, trigger common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("submitDeployRequestWithTrigger", protocolVersion, payload, trigger)
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
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitDeployRequestWithTrigger(protocolVersion uint8, payload []byte, trigger common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitDeployRequestWithTrigger", protocolVersion, payload, trigger)
}

// UnpackSubmitDeployRequestWithTrigger is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68a7c3fb.
//
// Solidity: function submitDeployRequestWithTrigger(uint8 protocolVersion, bytes payload, address trigger) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackSubmitDeployRequestWithTrigger(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("submitDeployRequestWithTrigger", data)
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

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (processorEndpoint *ProcessorEndpoint) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("upgradeToAndCall", newImplementation, data)
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
func (processorEndpoint *ProcessorEndpoint) TryPackUpgradeToAndCall(newImplementation common.Address, data []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("upgradeToAndCall", newImplementation, data)
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

// ProcessorEndpointInitialized represents a Initialized event raised by the ProcessorEndpoint contract.
type ProcessorEndpointInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointInitialized) ContractEventName() string {
	return ProcessorEndpointInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (processorEndpoint *ProcessorEndpoint) UnpackInitializedEvent(log *types.Log) (*ProcessorEndpointInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointInitialized)
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

// ProcessorEndpointUpgraded represents a Upgraded event raised by the ProcessorEndpoint contract.
type ProcessorEndpointUpgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointUpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointUpgraded) ContractEventName() string {
	return ProcessorEndpointUpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (processorEndpoint *ProcessorEndpoint) UnpackUpgradedEvent(log *types.Log) (*ProcessorEndpointUpgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointUpgraded)
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAddressEmptyCodeError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackFailedCallError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidInitializationError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackNotInitializingError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
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

// ProcessorEndpointAddressEmptyCode represents a AddressEmptyCode error raised by the ProcessorEndpoint contract.
type ProcessorEndpointAddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ProcessorEndpointAddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (processorEndpoint *ProcessorEndpoint) UnpackAddressEmptyCodeError(raw []byte) (*ProcessorEndpointAddressEmptyCode, error) {
	out := new(ProcessorEndpointAddressEmptyCode)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
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

// ProcessorEndpointERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ProcessorEndpoint contract.
type ProcessorEndpointERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ProcessorEndpointERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (processorEndpoint *ProcessorEndpoint) UnpackERC1967InvalidImplementationError(raw []byte) (*ProcessorEndpointERC1967InvalidImplementation, error) {
	out := new(ProcessorEndpointERC1967InvalidImplementation)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointERC1967NonPayable represents a ERC1967NonPayable error raised by the ProcessorEndpoint contract.
type ProcessorEndpointERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ProcessorEndpointERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (processorEndpoint *ProcessorEndpoint) UnpackERC1967NonPayableError(raw []byte) (*ProcessorEndpointERC1967NonPayable, error) {
	out := new(ProcessorEndpointERC1967NonPayable)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointFailedCall represents a FailedCall error raised by the ProcessorEndpoint contract.
type ProcessorEndpointFailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ProcessorEndpointFailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (processorEndpoint *ProcessorEndpoint) UnpackFailedCallError(raw []byte) (*ProcessorEndpointFailedCall, error) {
	out := new(ProcessorEndpointFailedCall)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
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

// ProcessorEndpointInvalidInitialization represents a InvalidInitialization error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func ProcessorEndpointInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidInitializationError(raw []byte) (*ProcessorEndpointInvalidInitialization, error) {
	out := new(ProcessorEndpointInvalidInitialization)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
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

// ProcessorEndpointNotInitializing represents a NotInitializing error raised by the ProcessorEndpoint contract.
type ProcessorEndpointNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func ProcessorEndpointNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (processorEndpoint *ProcessorEndpoint) UnpackNotInitializingError(raw []byte) (*ProcessorEndpointNotInitializing, error) {
	out := new(ProcessorEndpointNotInitializing)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
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

// ProcessorEndpointUUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the ProcessorEndpoint contract.
type ProcessorEndpointUUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func ProcessorEndpointUUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (processorEndpoint *ProcessorEndpoint) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*ProcessorEndpointUUPSUnauthorizedCallContext, error) {
	out := new(ProcessorEndpointUUPSUnauthorizedCallContext)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorEndpointUUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the ProcessorEndpoint contract.
type ProcessorEndpointUUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func ProcessorEndpointUUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (processorEndpoint *ProcessorEndpoint) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*ProcessorEndpointUUPSUnsupportedProxiableUUID, error) {
	out := new(ProcessorEndpointUUPSUnsupportedProxiableUUID)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
