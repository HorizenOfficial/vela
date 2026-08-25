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

// StructsBatchEntry is an auto generated low-level Go binding around an user-defined struct.
type StructsBatchEntry struct {
	PrevStateRoot      [32]byte
	NewStateRoot       [32]byte
	ProcessedRequestId [32]byte
	UserEvents         StructsEventData
	AppEvents          StructsEventData
	WithdrawalRequests []StructsWithdrawalRequest
	Refund             *big.Int
	ApplicationFees    *big.Int
	ErrorCode          uint8
	ErrorMsg           string
}

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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"extensionAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"selectedApplicationId\",\"type\":\"uint64\"}],\"name\":\"ApplicationNotSelected\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeployerNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientAppBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtension\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPermit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxNumOfApplicationsExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumStructs.RequestType\",\"name\":\"queueType\",\"type\":\"uint8\"}],\"name\":\"PriorityQueueNotServed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TriggerCannotBeEOA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"AppEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"entryHashes\",\"type\":\"bytes32[]\"}],\"name\":\"BatchProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"DeployRequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"DeployRequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMax\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"MaxNumberOfAppUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"ReportGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newGrace\",\"type\":\"uint256\"}],\"name\":\"SelectionGraceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TriggerExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"withdrawSuccess\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"postWithdrawSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"returnedTokens\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structStructs.TokenAndAmount[]\",\"name\":\"failedTokens\",\"type\":\"tuple[]\"}],\"name\":\"TriggerWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventSubType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPLOYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUEST_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_OPERATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"addAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adminReset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"appIds\",\"type\":\"uint64[]\"}],\"name\":\"adminResetApps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"appCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"applicationStateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableDeploySlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEvents\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"}],\"internalType\":\"structStructs.BatchEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"batchSignature\",\"type\":\"bytes\"}],\"name\":\"batchStateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"extension\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"facilitatorNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"payloadHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeployedAppIds\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getPendingRequestsPage\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getPendingRequestsWithStateRoot\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"requests\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTriggerQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTriggerRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resetOperator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"},{\"internalType\":\"contractITokenAllowlist\",\"name\":\"_tokenAllowlist\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"isAllowedDeployer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumOfApplications\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"removeAllowedDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facilitator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"selectionGrace\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"userEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"subTypes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structStructs.EventData\",\"name\":\"appEventData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMsg\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submitDeployRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"trigger\",\"type\":\"address\"}],\"name\":\"submitDeployRequestWithTrigger\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"requestSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"depositPermit\",\"type\":\"bytes\"}],\"name\":\"submitRequestFor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAllowlist\",\"outputs\":[{\"internalType\":\"contractITokenAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalAppCustody\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalPendingClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"triggerContracts\",\"outputs\":[{\"internalType\":\"contractITrigger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"triggersToAppIds\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMax\",\"type\":\"uint256\"}],\"name\":\"updateMaxNumOfApplications\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newGrace\",\"type\":\"uint256\"}],\"name\":\"updateSelectionGrace\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "ProcessorEndpoint",
	Bin: "0x60c03461015657601f615cf838819003918201601f19168301916001600160401b0383118484101761015a5780849260209460405283398101031261015657516001600160a01b03811680820361015657306080521561014757803b156101385760a0525f516020615cd85f395f51905f525460ff8160401c16610129576002600160401b03196001600160401b038216016100d3575b604051615b69908161016f823960805181818161122901526112ba015260a0518181816121c301528181612f540152818161373e015261435c0152f35b6001600160401b0319166001600160401b039081175f516020615cd85f395f51905f52556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f610096565b63f92ee8a960e01b5f5260045ffd5b6302e9da2760e11b5f5260045ffd5b632582a64160e11b5f5260045ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60c080604052600436101561001c575b50361561001a575f80fd5b005b5f3560e01c90816301ffc9a714612594575080631664deeb1461255a5780631753520f1461253d57806318733d57146124fc578063194b6be8146124c45780631cc76c38146124845780631f56b792146122bf57806320ff473f146103bb57806321c0b3421461229d578063248a9ca3146122665780632a0acc6a1461222c5780632b52b07c146121f25780632d5537b0146121ae5780632f2ff15d146121645780632fbfa0d514611a5e57806336568abe14611a1a57806344ce30f0146119fd57806347101cae146119e357806348562abd146114e35780634a0d0495146114af5780634f1ef2861461127d57806352d1902d146112175780635423112f146110995780635be374f51461104b5780635e33938414611025578063655f2de11461072257806367663bdf14610fce57806368a7c3fb14610f905780636b288d2014610f2f5780636e91d13314610f075780637596ed4314610d855780637a36a89114610d4e57806380a1f71214610d33578063840059ec14610ce357806384b0196e146109c35780638537c278146107e857806386f4a07d146107ae5780638c5b9b00146107995780638eb8791a1461077c57806390fda89b1461072257806391d14854146107275780639c73eeaf146107225780639c8d416f146107055780639fabc36f146106b2578063a217fddf14610670578063a9ba045e1461068a578063aa3aa46014610670578063ad3cb1cc14610629578063af20c960146105f1578063b7222ff51461059e578063bcf7a3c9146104ff578063be1e3725146104c5578063c415b95c1461049d578063d2c35ce81461047e578063d547741f1461042f578063d72e957a146103f7578063dd39d765146103c0578063e3d72e5e146103bb578063eabcca1014610355578063ecd002611461031b578063f624f88a146102f35763f7d9cc1e146102d2575f61000f565b346102ef575f3660031901126102ef576020600e54604051908152f35b5f80fd5b346102ef575f3660031901126102ef576011546040516001600160a01b039091168152602090f35b346102ef575f3660031901126102ef5760206040517ffc425f2263d0df187444b70e47283d622c70181c5baebb1306a01edba1ce184c8152f35b346102ef5760e03660031901126102ef576004356001600160a01b038116036102ef576024356001600160a01b038116036102ef57610392612629565b5061039b612655565b506103a4612613565b5060c4356001600160a01b03811614612f47575f80fd5b61047e565b346102ef5760203660031901126102ef576004356001600160401b0381116102ef576103f0903690600401612845565b5050612f47565b346102ef5760203660031901126102ef576001600160a01b036104186125e7565b165f526018602052602060405f2054604051908152f35b346102ef5760403660031901126102ef5761001a60043561044e6125fd565b90610479610474825f525f516020615af45f395f51905f52602052600160405f20015490565b612ff4565b613113565b346102ef5760203660031901126102ef576104976125e7565b50612f47565b346102ef575f3660031901126102ef576017546040516001600160a01b039091168152602090f35b346102ef5760403660031901126102ef576104fb6104e7602435600435614f12565b6040519182916020835260208301906127a0565b0390f35b6101403660031901126102ef576105146125e7565b5061051d612808565b506105266126ab565b50600560643510156102ef576084356001600160401b0381116102ef57610551903690600401612818565b505061055b61263f565b50610104356001600160401b0381116102ef5761057c903690600401612818565b5050610124356001600160401b0381116102ef576103f0903690600401612818565b346102ef5760403660031901126102ef576105b761267f565b6001600160401b036105c76125fd565b91165f52601460205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b346102ef5760203660031901126102ef576001600160a01b036106126125e7565b165f526015602052602060405f2054604051908152f35b346102ef575f3660031901126102ef576104fb60405161064a6040826128ac565b60058152640352e302e360dc1b60208201526040519182916020835260208301906126c1565b346102ef575f3660031901126102ef5760206040515f8152f35b346102ef575f3660031901126102ef576010546040516001600160a01b039091168152602090f35b346102ef5760203660031901126102ef5760206106fb600435805f52600483526106f6600760405f2001546001600160401b0360ff8260e81c169160a01c16615003565b61503e565b6040519015158152f35b346102ef575f3660031901126102ef576020601654604051908152f35b612939565b346102ef5760403660031901126102ef576107406125fd565b6004355f525f516020615af45f395f51905f5260205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b346102ef575f3660031901126102ef576020600254604051908152f35b346102ef575f36600319011215612f47575f80fd5b346102ef575f3660031901126102ef5760206040517fc580ee26c87bb95870d138607be6d8598a348af3608b6ce0cef58e73a4c959178152f35b346102ef576101803660031901126102ef5761080261267f565b608435906001600160401b0382116102ef57604060031983360301126102ef5760a435916001600160401b0383116102ef57604060031984360301126102ef5760c435926001600160401b0384116102ef57366023850112156102ef578360040135926001600160401b0384116102ef5736602460608602870101116102ef576101243591600c8310156102ef57610144356001600160401b0381116102ef576108b0903690600401612818565b610164356001600160401b0381116102ef576108d0903690600401612818565b9590946108db612f85565b6108e36130de565b6001600160401b031697885f525f60205260405f2054998a95604051996109098b612875565b8b8b5260243560208c015260443560408c015260643560608c0152610932903690600401612a36565b60808b0152610945903690600401612a36565b60a08a0152610958913691602401612b52565b60c088015260e43560e08801526101043561010088015261097d906101208801612bed565b3690610988926128e8565b61014085015261099793614003565b9182036109b1575b5f5f516020615b145f395f51905f525d005b5f525f60205260405f2055808061099f565b346102ef575f3660031901126102ef577fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100541580610cba575b15610c7d576040517fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10254815f610a3183612cd9565b8083529260018116908115610c5e5750600114610be0575b610a55925003826128ac565b6040517fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10354815f610a8583612cd9565b8083529260018116908115610bc15750600114610b43575b610ab091925092610ae7949303826128ac565b6020610af560405192610ac383856128ac565b5f84525f368137604051958695600f60f81b875260e08588015260e08701906126c1565b9085820360408701526126c1565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b828110610b2c57505050500390f35b835185528695509381019392810192600101610b1d565b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1035f90815290917f5f9ce34815f8e11431c7bb75a8e6886a91478f7ffc1dbb0a98dc240fddd76b755b818310610ba5575050906020610ab092820101610a9d565b6020919350806001915483858801015201910190918392610b8d565b60209250610ab094915060ff191682840152151560051b820101610a9d565b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1025f90815290917f42ad5d3e1f2e6e70edcf6d991b8a3023d3fca8047a131592f9edb9fd9b89d57d5b818310610c42575050906020610a5592820101610a49565b6020919350806001915483858801015201910190918392610c2a565b60209250610a5594915060ff191682840152151560051b820101610a49565b60405162461bcd60e51b81526020600482015260156024820152741152540dcc4c8e88155b9a5b9a5d1a585b1a5e9959605a1b6044820152606490fd5b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10154156109fc565b346102ef5760403660031901126102ef57610cfc6125e7565b610d046125fd565b6001600160a01b039182165f908152601260209081526040808320949093168252928352819020549051908152f35b346102ef575f3660031901126102ef576104fb6104e7614e48565b346102ef5760203660031901126102ef576001600160401b03610d6f61267f565b165f525f602052602060405f2054604051908152f35b346102ef575f3660031901126102ef57604051806020600154918281520190828260015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6925f905b806003830110610eb557610e08945491818110610e9b575b818110610e7e575b818110610e61575b10610e53575b5093929303826128ac565b604051918291602083019060208452518091526040830191905f5b818110610e31575050500390f35b82516001600160401b0316845285945060209384019390920191600101610e23565b60c01c815260200185610dfd565b9260206001916001600160401b038560801c168152019301610df7565b9260206001916001600160401b038560401c168152019301610def565b9260206001916001600160401b0385168152019301610de7565b916004919350608060019186546001600160401b03811682526001600160401b038160401c1660208301526001600160401b0381841c16604083015260c01c6060820152019401920185929391610dcf565b346102ef575f3660031901126102ef57600f546040516001600160a01b039091168152602090f35b346102ef5760203660031901126102ef57610f486125e7565b6001600160a01b03165f9081527f8f7e9a8734279607aeaac474f24949ddc99eae25c0da72edcf5ad21022be2c9e602090815260409182902054915160ff9092161515825290f35b60603660031901126102ef57610fa46127f8565b506024356001600160401b0381116102ef57610fc4903690600401612818565b5050610497612629565b346102ef5760e03660031901126102ef57610fe76125e7565b610fef612695565b9060443560058110156102ef5760209261101d9261100b612613565b9060c4359360a4359360643592614de3565b604051908152f35b346102ef575f3660031901126102ef57602061101d600d54611045612d2a565b90612d1d565b346102ef5760203660031901126102ef576001600160401b0361108f611072600435614ba4565b9192906040519485941684526060602085015260608401906127a0565b9060408301520390f35b346102ef5760203660031901126102ef576110b2612c87565b506004355f52600460205260405f206040516110cd81612875565b8154815260018201546001600160a01b0316602082015260028201546040808301919091526003830154606083015260048301546080830152516005830180545f9161111882612cd9565b80855291600181169081156111ef57506001146111b1575b6104fb8561119d60ff60078a89611149818b03826128ac565b60a086015260018060a01b0360068201541660c0860152015460018060a01b03811660e08501526001600160401b038160a01c16610100850152818160e01c1661012085015260e81c166101408301612d11565b6040519182916020835260208301906126f2565b5f908152602081209092505b8183106111d5575050810160200160076104fb611130565b6001816020929493945483858801015201910191906111bd565b60ff191660208087019190915292151560051b85019092019250600791506104fb9050611130565b346102ef575f3660031901126102ef577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316300361126e5760206040515f516020615ad45f395f51905f528152f35b63703e46dd60e11b5f5260045ffd5b60403660031901126102ef576112916125e7565b6024356001600160401b0381116102ef576112b090369060040161291e565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630811490811561148d575b5061126e57335f9081527f78e571b7bf30584d955e1c6444a2b5147087edf9f00485d94993a04d370525ea602052604090205460ff1615611456576040516352d1902d60e01b81526001600160a01b0383169290602081600481875afa5f9181611422575b506113625783634c9c8ce360e01b5f5260045260245ffd5b805f516020615ad45f395f51905f528592036114105750813b156113fe575f516020615ad45f395f51905f5280546001600160a01b031916821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a28151156113e6575f8083602061001a95519101845af46113e0615580565b91615a55565b5050346113ef57005b63b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f5260045260245ffd5b632a87526960e21b5f5260045260245ffd5b9091506020813d60201161144e575b8161143e602093836128ac565b810103126102ef5751908561134a565b3d9150611431565b63e2517d3f60e01b5f52336004527fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4260245260445ffd5b5f516020615ad45f395f51905f52546001600160a01b031614159050836112e5565b60403660031901126102ef576114c36127f8565b506024356001600160401b0381116102ef576103f0903690600401612818565b346102ef5760603660031901126102ef576114fc61267f565b6024356001600160401b0381116102ef5761151b903690600401612845565b6044929192356001600160401b0381116102ef5761153d903690600401612818565b929091611548612f85565b6115506130de565b80156119d45760018103611930575b61156b819592956129d6565b93611575826129bf565b9261158360405194856128ac565b828452601f19611592846129bf565b015f5b8181106118be5750506001600160401b035f9716965b8381036116f8575050600f5460405163cbe37a7f60e01b81529060209082906001600160a01b031681806115e4878b8d60048501612c4d565b03915afa9081156116ed575f916116be575b50156116af57929190855f525f60205260405f20549384935f935b83850361168357887fef3c0fe10558ded3fbc2c947a9e1fc9efa58574276e9bca52eb5f7a471eb74d561165c8a8a8a8103611671575b50604051918291602083526020830190612c1a565b0390a25f5f516020615b145f395f51905f525d005b845f525f60205260405f205584611647565b90919293956116a4848460019361169a8b87612c06565b51908b1591613465565b960193929190611611565b638baa579f60e01b5f5260045ffd5b6116e0915060203d6020116116e6575b6116d881836128ac565b81019061294f565b876115f6565b503d6116ce565b6040513d5f823e3d90fd5b611703818584612988565b356020611711838786612988565b0135906040611721848887612988565b01359161173c611732858988612988565b6060810190612a21565b9261175561174b868a89612988565b6080810190612a21565b611760868a89612988565b60a081013590601e19813603018212156102ef5701948535956001600160401b0387116102ef5760200160608702360381136102ef578e9689938c936101006117c88c60c06117b0828a8c612988565b01359760e06117c083838d612988565b013599612988565b013596600c8810156102ef576117e08f8d908f612988565b61012081013590601e19813603018212156102ef5701988935996001600160401b038b116102ef576020019a8a36038c136102ef576040519c6118228e612875565b8d5260208d015260408c015260608b01523661183d91612a36565b60808a01523661184c91612a36565b60a0890152369061185c92612b52565b60c087015260e086015261010085015261187a906101208501612bed565b3690611885926128e8565b610140820152806118968388612c06565b526118a18287612c06565b506118ab90613214565b6118b58289612c06565b526001016115ab565b6020906040999399516118d081612875565b5f81525f838201525f60408201525f60608201526118ec612a08565b60808201526118f9612a08565b60a0820152606060c08201525f60e08201525f6101008201525f610120820152606061014082015282828901015201979197611595565b843561013e19863603018112156102ef57604090860101355f52600460205260ff600760405f20015460e81c1660058110156119c057600481149081156119b7575b50801561198e575b1561155f576301d49e5b60e41b5f5260045ffd5b506001600160401b0382165f908152601960205260409020546001600160a01b0316151561197a565b90501586611972565b634e487b7160e01b5f52602160045260245ffd5b63c2e5347d60e01b5f5260045ffd5b346102ef575f3660031901126102ef57602061101d612d2a565b346102ef575f3660031901126102ef576020601b54604051908152f35b346102ef5760403660031901126102ef57611a336125fd565b336001600160a01b03821603611a4f5761001a90600435613113565b63334bd91960e11b5f5260045ffd5b60e03660031901126102ef57611a726127f8565b611a7a612695565b9060443560058110156102ef576064356001600160401b0381116102ef57611aa6903690600401612818565b9093611ab0612613565b9260a4359160ff60c43596169687612155576001600160401b03821695865f525f60205260405f20541561214657611ae66130de565b5f9284158015612139575b61212a57601654891061211b57611b0c600d54611045612d2a565b600e54111561210c576001600160a01b0382169384611f4357611b2f8a88612967565b3403611f34575b6119c05760038503611ea557608587141580611e9a575b611e8b5785611bbd925b895f52601460205260405f20865f5260205260405f20611b78838254612967565b9055855f52601560205260405f20611b91838254612967565b9055611b9e368a876128e8565b602081519101208a5f52600560205287600260405f2001549433614de3565b96865f52600560205260405f2095611c0260405193611bdb85612875565b42855260208501958652604085019788526060850193845260808501928b845236916128e8565b9060a0840191825260c084019233845260e08501975f89526101008601968b885261012087019d8e52611c3a6101408801998a612d11565b5f8d81526004602081905260409091209751885590516001880180546001600160a01b0319166001600160a01b0392909216919091179055905160028701559051600386015590519084015551805160058401916001600160401b038211611e7757611ca68354612cd9565b601f8111611e3c575b50602090601f8311600114611dd457600795949392915f9183611dc9575b50508160011b915f199060031b1c19161790555b516006820180546001600160a01b03199081166001600160a01b03938416179091559551929091018054909516911617808455905196519151969160e01b60ff60e01b169060058810156119c05760209760ff60e81b9060e81b16926001600160401b0360a01b9060a01b169069ffffffffffffffffffff60a01b19161717179055600281019081545f5284528260405f2055600181540190556001600d5401600d5581604051915f83527fb3e7adaaf2dfc1a38ff270a28d73879aa3825f72fb81b5b8585b91124d549236853394a45f5f516020615b145f395f51905f525d604051908152f35b015190508d80611ccd565b90601f19831691845f52815f20925f5b818110611e2457509160019391856007999897969410611e0c575b505050811b019055611ce1565b01515f1960f88460031b161c191690558d8080611dff565b92936020600181928786015181550195019301611de4565b611e6790845f5260205f20601f850160051c81019160208610611e6d575b601f0160051c0190615071565b8c611caf565b9091508190611e5a565b634e487b7160e01b5f52604160045260245ffd5b637c6953f960e01b5f5260045ffd5b5060e2871415611b4d565b8560028614611eb8575b611bbd92611b57565b506010546040516308b56a5360e11b8152600481018a905233602482015290602090829060449082906001600160a01b03165afa9081156116ed575f91611f15575b5015611f065785611eaf565b63622a850b60e11b5f5260045ffd5b611f2e915060203d6020116116e6576116d881836128ac565b8b611efa565b632a9ffab760e21b5f5260045ffd5b893403611f34578615611f345760115460405163cbe230c360e01b81526004810187905290602090829060249082906001600160a01b03165afa9081156116ed575f916120ed575b50156120de576040516370a0823160e01b8152306004820152602081602481895afa9081156116ed575f916120ac575b5060205f604051828101906323b872dd60e01b82523360248201523060448201528b606482015260648152611ff16084826128ac565b5190828a5af1156116ed575f513d6120a35750853b155b612090576040516370a0823160e01b8152306004820152906020826024818a5afa80156116ed5789925f91612057575b509061204391612d1d565b14611b365763b0a6ea2960e01b5f5260045ffd5b919250506020813d602011612088575b81612074602093836128ac565b810103126102ef5751889190612043612038565b3d9150612067565b85635274afe760e01b5f5260045260245ffd5b60011415612008565b90506020813d6020116120d6575b816120c7602093836128ac565b810103126102ef57518c611fbb565b3d91506120ba565b63514e24c360e11b5f5260045ffd5b612106915060203d6020116116e6576116d881836128ac565b8c611f8b565b63d8219b3d60e01b5f5260045ffd5b63c77f97eb60e01b5f5260045ffd5b630ee2db0560e21b5f5260045ffd5b505f935060048514611af1565b6309ff9a8360e41b5f5260045ffd5b635428eae760e01b5f5260045ffd5b346102ef5760403660031901126102ef5761001a6004356121836125fd565b906121a9610474825f525f516020615af45f395f51905f52602052600160405f20015490565b61303a565b346102ef575f3660031901126102ef576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346102ef575f3660031901126102ef5760206040517f952140b6347f31bf88278b2d6fb6365ec837393094e7a1f12b6dffbdc13333408152f35b346102ef575f3660031901126102ef5760206040517fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec428152f35b346102ef5760203660031901126102ef57602061101d6004355f525f516020615af45f395f51905f52602052600160405f20015490565b346102ef5760403660031901126102ef576122b66125e7565b506104976125fd565b346102ef575f3660031901126102ef576122d7612d2a565b6122e081612da9565b906122ee600a549182612967565b905f905b82810361230f57604051602080825281906104fb908201876127a0565b805f52600960205260405f20545f52600460205260405f206040519061233482612875565b8054825260018101546001600160a01b0316602083015260028101546040808401919091526003820154606084015260048201546080840152516005820180545f9161237f82612cd9565b808552916001811690811561245d5750600114612424575b505060ff600760019694846123b38998966124059603826128ac565b60a0860152868060a01b0360068201541660c08601520154858060a01b03811660e08501526001600160401b038160a01c16610100850152818160e01c1661012085015260e81c166101408301612d11565b61240f8588612c06565b5261241a8487612c06565b50019101906122f2565b5f908152602081209092505b818310612447575050810160200160ff6007612397565b6001816020925483868801015201920191612430565b60ff191660208087019190915292151560051b8501909201925060ff915060079050612397565b346102ef5760203660031901126102ef576001600160401b036124a561267f565b165f526019602052602060018060a01b0360405f205416604051908152f35b346102ef5760203660031901126102ef576001600160a01b036124e56125e7565b165f526013602052602060405f2054604051908152f35b346102ef5760203660031901126102ef576001600160a01b0361251d6125e7565b165f52601a60205260206001600160401b0360405f205416604051908152f35b346102ef575f3660031901126102ef576020600354604051908152f35b346102ef575f3660031901126102ef5760206040517f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd988152f35b346102ef5760203660031901126102ef576004359063ffffffff60e01b82168092036102ef57602091637965db0b60e01b81149081156125d6575b5015158152f35b6301ffc9a760e01b149050836125cf565b600435906001600160a01b03821682036102ef57565b602435906001600160a01b03821682036102ef57565b608435906001600160a01b03821682036102ef57565b604435906001600160a01b03821682036102ef57565b60a435906001600160a01b03821682036102ef57565b606435906001600160a01b03821682036102ef57565b35906001600160a01b03821682036102ef57565b600435906001600160401b03821682036102ef57565b602435906001600160401b03821682036102ef57565b604435906001600160401b03821682036102ef57565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9060058210156119c05752565b9061279d908251815260018060a01b0360208401511660208201526040830151604082015260608301516060820152608083015160808201526101408061274a60a086015161016060a08601526101608501906126c1565b9460018060a01b0360c08201511660c085015260018060a01b0360e08201511660e08501526001600160401b036101008201511661010085015260ff6101208201511661012085015201519101906126e5565b90565b9080602083519182815201916020808360051b8301019401925f915b8383106127cb57505050505090565b90919293946020806127e9600193601f1986820301875289516126f2565b970193019301919392906127bc565b6004359060ff821682036102ef57565b6024359060ff821682036102ef57565b9181601f840112156102ef578235916001600160401b0383116102ef57602083818601950101116102ef57565b9181601f840112156102ef578235916001600160401b0383116102ef576020808501948460051b0101116102ef57565b61016081019081106001600160401b03821117611e7757604052565b604081019081106001600160401b03821117611e7757604052565b90601f801991011681019081106001600160401b03821117611e7757604052565b6001600160401b038111611e7757601f01601f191660200190565b9291926128f4826128cd565b9161290260405193846128ac565b8294818452818301116102ef578281602093845f960137010152565b9080601f830112156102ef5781602061279d933591016128e8565b346102ef57602036600319011215612f47575f80fd5b908160209103126102ef575180151581036102ef5790565b9190820180921161297457565b634e487b7160e01b5f52601160045260245ffd5b91908110156129ab5760051b8101359061013e19813603018212156102ef570190565b634e487b7160e01b5f52603260045260245ffd5b6001600160401b038111611e775760051b60200190565b906129e0826129bf565b6129ed60405191826128ac565b82815280926129fe601f19916129bf565b0190602036910137565b60405190612a1582612891565b60606020838281520152565b903590603e19813603018212156102ef570190565b91906040838203126102ef57604051612a4e81612891565b809380356001600160401b0381116102ef57810183601f820112156102ef578035612a78816129bf565b91612a8660405193846128ac565b81835260208084019260051b820101918683116102ef5760208201905b838210612b25575050505082526020810135906001600160401b0382116102ef57019180601f840112156102ef578235612adc816129bf565b93612aea60405195866128ac565b81855260208086019260051b8201019283116102ef57602001905b828210612b155750505060200152565b8135815260209182019101612b05565b81356001600160401b0381116102ef57602091612b478a84809488010161291e565b815201910190612aa3565b929192612b5e826129bf565b93612b6c60405195866128ac565b60606020868581520193028201918183116102ef57925b828410612b905750505050565b6060848303126102ef576040519060608201908282106001600160401b03831117611e7757606092602092604052612bc78761266b565b8152612bd483880161266b565b8382015260408701356040820152815201930192612b83565b600c8210156119c05752565b8051156129ab5760200190565b80518210156129ab5760209160051b010190565b90602080835192838152019201905f5b818110612c375750505090565b8251845260209384019390920191600101612c2a565b9192602093612c658293604086526040860190612c1a565b9385818603910152818452848401375f828201840152601f01601f1916010190565b60405190612c9482612875565b5f61014083828152826020820152826040820152826060820152826080820152606060a08201528260c08201528260e082015282610100820152826101208201520152565b90600182811c92168015612d07575b6020831014612cf357565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612ce8565b60058210156119c05752565b9190820391821161297457565b600b54600a54808211612d3d5750505f90565b61279d91612d1d565b600854600754808211612d3d5750505f90565b60016002820154910154808211612d3d5750505f90565b60405190612d7f6020836128ac565b5f80835282815b828110612d9257505050565b602090612d9d612c87565b82828501015201612d86565b90612db3826129bf565b612dc060405191826128ac565b8281528092612dd1601f19916129bf565b01905f5b828110612de157505050565b602090612dec612c87565b82828501015201612dd5565b90612e036001612da9565b9160018101546001810191828211612974575f915b838103612e255750505050565b805f528160205260405f20545f52600460205260405f2060405190612e4982612875565b8054825260018101546001600160a01b0316602083015260028101546040808401919091526003820154606084015260048201546080840152516005820180545f91612e9482612cd9565b8085529160018116908115612f205750600114612ee7575b505060ff600760019694846123b3899896612ec89603826128ac565b612ed2868a612c06565b52612edd8589612c06565b5001920191612e18565b5f908152602081209092505b818310612f0a575050810160200160ff6007612eac565b6001816020925483868801015201920191612ef3565b60ff191660208087019190915292151560051b8501909201925060ff915060079050612eac565b604051365f82375f8036837f00000000000000000000000000000000000000000000000000000000000000005af4903d91825f833e15612f8357f35bfd5b335f9081527f212403f86ff12f15ffbeab879ed2f25237f7273bd0873f6abf8e9330ca4ed4b9602052604090205460ff1615612fbd57565b63e2517d3f60e01b5f52336004527f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9860245260445ffd5b5f8181525f516020615af45f395f51905f526020908152604080832033845290915290205460ff16156130245750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f8181525f516020615af45f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff166130d8575f8181525f516020615af45f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b5f516020615b145f395f51905f525c6131045760015f516020615b145f395f51905f525d565b633ee5aeb560e01b5f5260045ffd5b5f8181525f516020615af45f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff16156130d8575f8181525f516020615af45f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b9080602083519182815201916020808360051b8301019401925f915b8383106131da57505050505090565b90919293946020806131f8600193601f1986820301875289516126c1565b970193019301919392906131cb565b90600c8210156119c05752565b60808101906020825151604051613247816132398582019486865260408301906131af565b03601f1981018352826128ac565b5190209251015160405161326b816132396020820194602086526040830190612c1a565b5190209160a082016020815151604051613293816132398582019486865260408301906131af565b519020915101516040516132b7816132396020820194602086526040830190612c1a565b5190209260c081015160405160208101918160408101916020855280518093526020606083019101925f5b8181106133c45750506132fe925003601f1981018352826128ac565b519020926001600160401b038251169560208301519460408401519660608501519560e0860151936101008701519561012088015197600c8910156119c0576101400151986040519b8c9b60208d019e8f5260408d015260608c015260808b015260a08a015260c089015260e0880152610100870152610120860152610140850152610160840152610180830161339491613207565b6101a082016101a090526101c082016133ac916126c1565b03601f19810182526133be90826128ac565b51902090565b845180516001600160a01b0390811685526020808301519091168186015260409182015191850191909152909401938593506060909201916001016132e2565b906001548210156129ab5760015f52600282901c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6019160031b60181690565b91909180548310156129ab575f52601860205f208360021c019260031b1690565b935091506001600160401b0383511690606084015190815f52600460205260405f20936007850154946001600160401b038660a01c169260ff8760e81c16906134b2866106f68488615003565b15613ff45760058210156119c0576004821494870361214657613fe4575b5081602088015103613d32576080870151916020835151930151518303611e8b5760a0880151936020855151950151518503611e8b5760038301805460068501546001600160a01b039081169a91949291168015613fdd57995b6101208c01918251600c8110156119c057613d415750505060408a015114613d325715613d07575b505f60a0525f60a05260c0870151928351925f61356e856129d6565b613577866129d6565b905f5b8760a05110613b9f575f60a0525b8060a05110613ad9575050506135b1906135ac60e08c01516101008d015190612967565b612967565b6135b9615112565b10613aca575f60a0525b8060a05103613a6b57505f60a0525b8060a051036139f957506007015460e81c60ff169460058610156119c057600286146139ce575b8515613987575b83857f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e604060208b0151818c015182519182526020820152a360e08701519081613944575b50505f60a0819052848152601960205260408120549092906001600160a01b0316801515908161392e575b613679846129d6565b935b8060a051106137f35750505050613691826129d6565b915f60a0525b8060a051036137c057505060a085015190604051906136fd602083019363214040e160e21b85528660248501528560448501526080606485015260206136e98251604060a488015260e48701906131af565b91015184820360a3190160c4860152612c1a565b602319838203016084840152602080835192838152019201905f5b8181106137a157505050918161373a5f9493859403601f1981018352826128ac565b51907f00000000000000000000000000000000000000000000000000000000000000005af4613767615580565b901561379957509161379491849361010060409601519086519261378c6020856128ac565b5f84526153fd565b015190565b602081519101fd5b82516001600160a01b0316845260209384019390920191600101613718565b60a0516001600160a01b03906137d69084612c06565b51166137e460a05185612c06565b52600160a0510160a052613697565b828061390d575b6138dd575b60a05161384c906001600160a01b03906138199087612c06565b51511660018060a01b03602061383160a05189612c06565b51015116604061384360a05189612c06565b510151916151d9565b60a0516001600160a01b03906020906138659087612c06565b5101511687897faafeb614695d37710d00f37c69671d205eb5b6fcd02995f6c3c7e2c93d3cbb2460018060a01b0361389f60a0518a612c06565b51511660406138b060a0518b612c06565b510151604080516001600160a01b03939093168352602083019190915290a4600160a0510160a05261367b565b60a05190956001916001600160a01b03906138f89087612c06565b5151166139058288612c06565b5201946137ff565b508160018060a01b03602061392460a05188612c06565b51015116146137fa565b6040890151875f525f60205260405f2055613670565b6001600160a01b031690613958908261517a565b60e0870151604080515f81526020810192909252859187915f516020615ab45f395f51905f5291a45f80613645565b600154600160401b811015611e7757856139aa8260016139c99401600155613404565b9091906001600160401b038084549260031b9316831b921b1916179055565b613600565b83857ff36f013d86705b333febcccf7ece97eee30f4cd398a9c69aee1f940ed5daddea5f80a36135f9565b613a0d60a051602060a08b01510151612c06565b5185877f485d4bfb37695319d1b2352eef5982d10db8c284b0a7ddefe4465377c9fbc8df613a5a613a458d60a0805191015151612c06565b516040519182916020835260208301906126c1565b0390a4600160a0510160a0526135d2565b8886887f9ede1ad14e46a641d7b174bb6c9d939ad961d1a7de4b669163b9714a67e22320613ab9613a45613aa960a051602060808901510151612c06565b5195608060a05191015151612c06565b0390a4600160a0510160a0526135c3565b631e9acf1760e31b5f5260045ffd5b60a0516001600160a01b0390613aef9084612c06565b516040516370a0823160e01b81523060048201529116602082602481845afa9182156116ed575f92613b6a575b50613b4481613b57925f52601560205260405f2054905f52601360205260405f205490612967565b613b5060a05187612c06565b5190612967565b11613aca57600160a0510160a052613588565b9091506020813d8211613b97575b81613b85602093836128ac565b810103126102ef575190613b44613b1c565b3d9150613b78565b9260018060a01b03613bb360a0518b612c06565b515116906040613bc560a0518c612c06565b510151918c5f52601460205260405f2060018060a01b0382165f5260205260405f20548311613cf8578c5f52601460205260405f2060018060a01b0382165f5260205260405f20838154039055805f52601560205260405f20548311613aca57805f52601560205260405f2083815403905580155f14613c5a575090613c4a91612967565b925b600160a0510160a05261357a565b9b9e9d9c9b909490915f805b838110613ca3575b509f9c9d9e9f15613c82575b509050613c4c565b600192613c8f8386612c06565b52613c9a8286612c06565b5201805f613c7a565b846001600160a01b03613cb68389612c06565b511614613cc557600101613c66565b90509e9f9c9d9e613cea613ce383613cdd848a612c06565b51612967565b9187612c06565b5260019c9f9e9d9c5f613c6e565b6306a1762560e11b5f5260045ffd5b60e088015190613d1e6101008a0192835190612967565b03611f34575160165411611f34575f613552565b630b6fac0360e41b5f5260045ffd5b929b9a93959897909491961590811591613fd3575b508015613fc5575b611e8b5760408a01518b03613d3257600288015460018060a01b0360018a015416918a5f52601460205260405f2060018060a01b0384165f5260205260405f20548211613cf857613dad615112565b10613aca5787938a93845f52601460205260405f2060018060a01b0385165f5260205260405f20613ddf848254612d1d565b9055835f52601560205260405f20613df8848254612d1d565b905554601654808210613fbc57613e0e91612d1d565b925b80613f5d57508115613f0d57505f516020615ab45f395f51905f5291613e3591612967565b613e3f818661517a565b604080515f8152602081019290925290a45b60ff600785015460e81c1660058110156119c05715613ea1575b15613e98575f905b5191600c8310156119c05760ff600761014061279d98015195015460e81c1694615259565b60165490613e73565b6003805460010190555f858152601960205260409020546001600160a01b031680613ecd575b50613e6b565b5f52601a60205260405f206001600160401b03198154169055845f52601960205260405f206bffffffffffffffffffffffff60a01b81541690555f613ec7565b9194505083613f20575b50505050613e51565b6001600160a01b0316925f516020615ab45f395f51905f5290613f43818661517a565b604080515f8152602081019290925290a45f838682613f17565b8490838793959894613f7d575b505050505083613f205750505050613e51565b81613f9782875f516020615ab45f395f51905f52956151d9565b604080516001600160a01b039290921682526020820192909252a45f83838280613f6a565b50505f92613e10565b5060c08a0151511515613d5e565b905015155f613d56565b508961352a565b613fee9086615087565b5f6134d0565b6302e8145360e61b5f5260045ffd5b9392939190916080526001600160401b036080515116906060608051015194855f52600460205260405f206007810154946001600160401b038660a01c1660ff8760e81c166140568a6106f68385615003565b15613ff45760058110156119c05760048114918703612146576140799087615087565b816020608051015103613d3257608080510151936020855151950151518503611e8b5760a06080510151956020875151970151518703611e8b57604060a081905280516141229360209390926140cf90836128ac565b6001825260a051601f190136858401376140ea608051613214565b6140f383612bf9565b52600f5460a0515163cbe37a7f60e01b81529586946001600160a01b0390921693859384939160048501612c4d565b03915afa90811561474b575f91614b85575b50156116af5760038301805460068501546001600160a01b039081169991949291168015614b7e57985b61012060805101918251600c8110156119c0576148c45750505060a051608051015114613d325715614895575b5060805160c0015180519390925f806141a3876129d6565b6141ac886129d6565b905f925b89811061475757505f5b83811061469557505050506141e1906135ac60e06080510151610100608051015190612967565b6141e9615112565b10613aca575f5b81810361463f5750505f5b8181036145d15750506007015460e81c60ff169360058510156119c057600285146145a6575b841561457e575b86847f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e6020608051015160a051608051015160a05151918252602082015260a05190a360e060805101519081614535575b50505f838152601960205260a051812054909182916001600160a01b0316801515908161451b575b6142ac8694966129d6565b955b8481106144055750505050506142c3816129d6565b915f5b8281036143dd5750505060a060805101519060a051519061431b602083019363214040e160e21b85528560248501528860448501526080606485015260206136e9825160a05160a488015260e48701906131af565b602319838203016084840152602080835192838152019201905f5b8181106143be5750505091816143585f9493859403601f1981018352826128ac565b51907f00000000000000000000000000000000000000000000000000000000000000005af493614386615580565b94156143b6576143ab93945061010060805101519060a051519261378c6020856128ac565b60a051608051015190565b845160208601fd5b82516001600160a01b0316845260209384019390920191600101614336565b6001906001600160a01b036143f28285612c06565b51166143fe8287612c06565b52016142c6565b82806144fb575b6144cc575b60019061444f6001600160a01b036144298388612c06565b515116838060a01b03602061443e858a612c06565b5101511660a051613843858a612c06565b818060a01b0360206144618388612c06565b510151168c8a7faafeb614695d37710d00f37c69671d205eb5b6fcd02995f6c3c7e2c93d3cbb24858060a01b03614498868b612c06565b51511660a0516144a8878c612c06565b51015160a051516001600160a01b039290921682526020820152604090a4016142ae565b9460019081906001600160a01b036144e48988612c06565b5151166144f1828b612c06565b5201959050614411565b50816001600160a01b0360206145118488612c06565b510151161461440c565b60a0516080510151875f525f60205260a0515f20556142a1565b6001600160a01b031690614549908261517a565b86845f516020615ab45f395f51905f5260e0608051015160a05151809160a0518201905f835260208301520390a45f80614279565b600154600160401b811015611e7757846139aa8260016145a19401600155613404565b614228565b86847ff36f013d86705b333febcccf7ece97eee30f4cd398a9c69aee1f940ed5daddea5f80a3614221565b806145e8600192602060a060805101510151612c06565b518a887f485d4bfb37695319d1b2352eef5982d10db8c284b0a7ddefe4465377c9fbc8df6146366146208660a0608051015151612c06565b5160a051519182916020835260208301906126c1565b0390a4016141fb565b8061465560019260206080805101510151612c06565b518b897f9ede1ad14e46a641d7b174bb6c9d939ad961d1a7de4b669163b9714a67e2232061468c6146208660808051015151612c06565b0390a4016141f0565b6001600160a01b036146a78284612c06565b5160a051516370a0823160e01b81523060048201529116602082602481845afa91821561474b575f92614716575b506146ff81614709925f52601560205260a0515f2054905f52601360205260a0515f205490612967565b613b508487612c06565b11613aca576001016141ba565b9091506020813d8211614743575b81614731602093836128ac565b810103126102ef5751906146ff6146d5565b3d9150614724565b60a051513d5f823e3d90fd5b936001600160a01b0361476a868b612c06565b5151169060a05161477b878c612c06565b510151918c5f52601460205260a0515f2060018060a01b0382165f5260205260a0515f20548311613cf8575f8d81526014602090815260a0518083206001600160a01b038516845282528083208054879003905583835260159091529020548311613aca57805f52601560205260a0515f2083815403905580155f14614810575060019161480891612967565b945b016141b0565b909591905f805b878110614857575b5015614830575b505060019061480a565b94600192959183926148428387612c06565b5261484d8287612c06565b520193905f614826565b826001600160a01b0361486a8389612c06565b51161461487957600101614817565b905061488c613ce384613cdd848a612c06565b5260015f61481f565b60e06080510151906148b06101006080510192835190612967565b03611f34575160165411611f34575f61418b565b929a9b99939597909491961590811591614b74575b508015614b64575b611e8b5760a05160805101518a03613d3257600287015460018801545f8a81526014602090815260a0518083206001600160a01b03909416808452939091529020549092908211613cf857614934615112565b10613aca5789938993845f52601460205260a0515f2060018060a01b0385165f5260205260a0515f20614968848254612d1d565b9055835f52601560205260a0515f20614982848254612d1d565b905554601654808210614b5b5761499891612d1d565b925b80614af957508115614aa357505f516020615ab45f395f51905f52916149bf91612967565b6149c9818661517a565b60a05151809160a0518201905f835260208301520390a45b60ff600784015460e81c1660058110156119c05715614a34575b15614a2b575f905b5190600c8210156119c05761279d9460ff6007610140608051015195015460e81c1694615259565b60165490614a03565b6003805460010190555f848152601960205260a0519020546001600160a01b031680614a61575b506149fb565b5f52601a60205260a0515f206001600160401b03198154169055835f52601960205260a0515f206bffffffffffffffffffffffff60a01b81541690555f614a5b565b9194505083614ab6575b505050506149e1565b6001600160a01b0316925f516020615ab45f395f51905f5290614ad9818661517a565b60a05151809160a0518201905f835260208301520390a45f858582614aad565b8490838793959894614b19575b505050505083614ab657505050506149e1565b81614b3382875f516020615ab45f395f51905f52956151d9565b60a051516001600160a01b039190911681526020810191909152604090a45f83838280614b06565b50505f9261499a565b5060c060805101515115156148e1565b905015155f6148d9565b508861415e565b614b9e915060203d6020116116e6576116d881836128ac565b5f614134565b80614bb9575b50614bb3612d70565b5f918290565b614bc1612d2a565b614dd857614bcd612d46565b614da157614bd96155cd565b614be35750614baa565b90916001600160401b038216805f526005602052614c0360405f20612d59565b93808511614d99575b505f818152601960205260409020546001600160a01b0316614d90575b805f52600560205260405f20614c3e85612da9565b94614c4e60018301549182612967565b915f915b838103614c6e57505050505f525f60205260405f205491929190565b805f528160205260405f20545f52600460205260405f2060405190614c9282612875565b8054825260018101546001600160a01b0316602083015260028101546040808401919091526003820154606084015260048201546080840152516005820180545f91614cdd82612cd9565b8085529160018116908115614d695750600114614d30575b505060ff600760019694846123b3899896614d119603826128ac565b614d1b868c612c06565b52614d26858b612c06565b5001920191614c52565b5f908152602081209092505b818310614d53575050810160200160ff6007614cf5565b6001816020925483868801015201920191614d3c565b60ff191660208087019190915292151560051b8501909201925060ff915060079050614cf5565b60019350614c29565b93505f614c0c565b50614dac6006612df8565b906001600160401b03610100614dc184612bf9565b5101511690815f525f60205260405f205491929190565b50614dac6009612df8565b959391614e1990959391956001600160401b0360405197602089019960018060a01b03168a5216604088015260608701906126e5565b60808501526001600160a01b031660a084015260c083015260e08083019190915281526133be610100826128ac565b614e5e614e59600d54611045612d2a565b612da9565b90614e688261564a565b6001545f905b808203614ed857505080158015614ed1575b614ec657614e8d81612da9565b925f5b828103614e9c57505050565b600190614ea98184612c06565b51614eb48288612c06565b52614ebf8187612c06565b5001614e90565b50905061279d612d70565b505f614e80565b9091614f0a6001916001600160401b03614ef186613404565b90549060031b1c165f5260056020528660405f2061578b565b920190614e6e565b9190614f25614e59600d54611045612d2a565b614f2e8161564a565b6001545f905b808203614fc9575050808510801590614fc1575b614fb45784614f5691612d1d565b91808311614fac575b50614f6982612da9565b935f5b838103614f795750505050565b80614f8f614f8960019385612967565b85612c06565b51614f9a8289612c06565b52614fa58188612c06565b5001614f6c565b91505f614f5f565b505050905061279d612d70565b508215614f48565b9091614ffb6001916001600160401b03614fe286613404565b90549060031b1c165f5260056020528560405f2061578b565b920190614f34565b9060058110156119c057600481146150375715615031576001600160401b03165f52600560205260405f2090565b50600690565b5050600990565b906002820154916001810154809311928361505a575b50505090565b909192505f5260205260405f2054145f8080615054565b81811061507c575050565b5f8155600101615071565b9060058110156119c057600481146150fb57601b54906150a6826158cf565b6150ff57156150fb576150b88161590c565b6150e8576150c59161593f565b6150cc5750565b6001600160401b0390633820f7cd60e11b5f521660045260245ffd5b63643c800360e01b5f525f60045260245ffd5b5050565b63643c800360e01b5f526004805260245ffd5b5f805260136020527f8fa6efc3be94b5b348b21fea823fe8d100408cee9b7f90524494500445d8ff6c5461279d9061514a9047612d1d565b5f805260156020527fa31547ce6245cdb9ecea19cf8c7eb9f5974025bb4075011409251ae855b30aed5490612d1d565b6001600160a01b03165f9081527f7e7fa33969761a458e04f477e039a608702b4f924981d6653935a8319a08ad7b6020526040902080546151bc908390612967565b90555f805260136020526151d560405f20918254612967565b9055565b60018060a01b031690815f52601260205260405f209060018060a01b03165f5260205260405f2061520b838254612967565b90555f5260136020526151d560405f20918254612967565b9093929193815260028410156119c05761524c60809261279d9560208401526040830190613207565b81606082015201906126c1565b909291936152678683615003565b600181019081545f528060205260405f20545f5260046020525f600760408220828155826001820155826002820155826003820155826004820155600581016152b08154612cd9565b90816153ba575b5050826006820155015581545f526020525f6040812055600181540190555f19600d5401600d556005861015806119c0576004871415806153af575b6153a1575b6119c0576153589561535a577f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d916153436001600160401b0392604051938493169560018985615223565b0390a35b6017546001600160a01b031661517a565b565b7f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed66916153996001600160401b0392604051938493169560018985615223565b0390a3615347565b6153aa836159fb565b6152f8565b50505f8615156152f3565b81601f8693116001146153d15750555b5f806152b7565b818352602083206153ed91601f0160051c810190600101615071565b80825281602081209155556153ca565b9290919261540b8582615003565b600181019081545f528060205260405f20545f5260046020525f600760408220828155826001820155826002820155826003820155826004820155600581016154548154612cd9565b908161553d575b5050826006820155015581545f526020525f6040812055600181540190555f19600d5401600d556005851015806119c057600486141580615532575b615524575b6119c057615358946154e6576001600160401b037f7e4f21a95b02bf0b4f758cc83ba1675c7414fb6af8cac506dd37309739a38b0d9161534360405192839216945f808985615223565b6001600160401b037f3d0b7d170f8bbb2715d0b8d20b8f280643d1c0aac796da8e6fde3cf97d05ed669161539960405192839216945f808985615223565b61552d826159fb565b61549c565b50505f851515615497565b81601f8693116001146155545750555b5f8061545b565b8183526020832061557091601f0160051c810190600101615071565b808252816020812091555561554d565b3d156155aa573d90615591826128cd565b9161559f60405193846128ac565b82523d5f602084013e565b606090565b81156155b9570690565b634e487b7160e01b5f52601260045260245ffd5b600154801561564357600c545f5b8281036155eb575050505f905f90565b6001600160401b03615610615609856156048587612967565b6155af565b6001613444565b90549060031b1c16805f52600560205261562c60405f20612d59565b61563957506001016155db565b9250505090600190565b505f905f90565b5f9060075490600854915b8281036156625750505090565b805f9492939452600660205260405f20545f52600460205260405f206040519061568b82612875565b8054825260018101546001600160a01b0316602083015260028101546040808401919091526003820154606084015260048201546080840152516005820180545f916156d682612cd9565b8085529160018116908115615764575060011461572b575b505060ff600760019694846123b389989661570a9603826128ac565b6157148587612c06565b5261571f8486612c06565b50019101929190615655565b5f908152602081209092505b81831061574e575050810160200160ff60076156ee565b6001816020925483868801015201920191615737565b60ff191660208087019190915292151560051b8501909201925060ff9150600790506156ee565b906001820154916002810154925b8381036157a7575050505090565b805f95939495528160205260405f20545f52600460205260405f20604051906157cf82612875565b8054825260018101546001600160a01b0316602083015260028101546040808401919091526003820154606084015260048201546080840152516005820180545f9161581a82612cd9565b80855291600181169081156158a8575060011461586f575b505060ff600760019694846123b389989661584e9603826128ac565b6158588688612c06565b526158638587612c06565b50019201939291615799565b5f908152602081209092505b818310615892575050810160200160ff6007615832565b600181602092548386880101520192019161587b565b60ff191660208087019190915292151560051b8501909201925060ff915060079050615832565b6158d96009612d59565b1561590757600a545f52600960205260405f20545f52600460205261590260405f205442612d1d565b101590565b505f90565b6159166006612d59565b15615907576007545f52600660205260405f20545f52600460205261590260405f205442612d1d565b60015491600c545f5b84810361595a5750505050505f905f90565b6001600160401b03615973615609876156048587612967565b90549060031b1c16805f52600560205261598f60405f20612d59565b61599d575b50600101615948565b6001600160401b03851681146159ef57805f52600560205260405f2060018101545f5260205260405f20545f526004602052836159de60405f205442612d1d565b106159945794505050505090600190565b5050505050505f905f90565b6001545f5b818103615a0c57505050565b6001600160401b03615a1f826001613444565b90549060031b1c166001600160401b03841614615a3e57600101615a00565b60010191508103615a5057505f600c55565b600c55565b90615a795750805115615a6a57805190602001fd5b63d6bda27560e01b5f5260045ffd5b81511580615aaa575b615a8a575090565b639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b50803b15615a8256fec9961aaccfc3406b40d0d25a5b83b8a4721e5f454116d7e22bd8038c715140aa360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268009b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00a2646970667358221220b3042fcf75bc8c61aa2abf4c289f3d255a102b3ccbaee7ec1ac402f5ac7c879b64736f6c634300081e0033f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
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
// Solidity: constructor(address extensionAddress) returns()
func (processorEndpoint *ProcessorEndpoint) PackConstructor(extensionAddress common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("", extensionAddress)
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

// PackBatchStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48562abd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function batchStateUpdate(uint64 applicationId, (bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string)[] entries, bytes batchSignature) returns()
func (processorEndpoint *ProcessorEndpoint) PackBatchStateUpdate(applicationId uint64, entries []StructsBatchEntry, batchSignature []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("batchStateUpdate", applicationId, entries, batchSignature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBatchStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48562abd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function batchStateUpdate(uint64 applicationId, (bytes32,bytes32,bytes32,(bytes[],bytes32[]),(bytes[],bytes32[]),(address,address,uint256)[],uint256,uint256,uint8,string)[] entries, bytes batchSignature) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackBatchStateUpdate(applicationId uint64, entries []StructsBatchEntry, batchSignature []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("batchStateUpdate", applicationId, entries, batchSignature)
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

// PackExtension is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2d5537b0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function extension() view returns(address)
func (processorEndpoint *ProcessorEndpoint) PackExtension() []byte {
	enc, err := processorEndpoint.abi.Pack("extension")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExtension is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2d5537b0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function extension() view returns(address)
func (processorEndpoint *ProcessorEndpoint) TryPackExtension() ([]byte, error) {
	return processorEndpoint.abi.Pack("extension")
}

// UnpackExtension is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2d5537b0.
//
// Solidity: function extension() view returns(address)
func (processorEndpoint *ProcessorEndpoint) UnpackExtension(data []byte) (common.Address, error) {
	out, err := processorEndpoint.abi.Unpack("extension", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
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

// PackGetPendingRequestsWithStateRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5be374f5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingRequestsWithStateRoot(uint256 maxCount) view returns(uint64 applicationId, (uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[] requests, bytes32 stateRoot)
func (processorEndpoint *ProcessorEndpoint) PackGetPendingRequestsWithStateRoot(maxCount *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("getPendingRequestsWithStateRoot", maxCount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingRequestsWithStateRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5be374f5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingRequestsWithStateRoot(uint256 maxCount) view returns(uint64 applicationId, (uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[] requests, bytes32 stateRoot)
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequestsWithStateRoot(maxCount *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequestsWithStateRoot", maxCount)
}

// GetPendingRequestsWithStateRootOutput serves as a container for the return parameters of contract
// method GetPendingRequestsWithStateRoot.
type GetPendingRequestsWithStateRootOutput struct {
	ApplicationId uint64
	Requests      []StructsPendingRequest
	StateRoot     [32]byte
}

// UnpackGetPendingRequestsWithStateRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5be374f5.
//
// Solidity: function getPendingRequestsWithStateRoot(uint256 maxCount) view returns(uint64 applicationId, (uint256,address,uint256,uint256,bytes32,bytes,address,address,uint64,uint8,uint8)[] requests, bytes32 stateRoot)
func (processorEndpoint *ProcessorEndpoint) UnpackGetPendingRequestsWithStateRoot(data []byte) (GetPendingRequestsWithStateRootOutput, error) {
	out, err := processorEndpoint.abi.Unpack("getPendingRequestsWithStateRoot", data)
	outstruct := new(GetPendingRequestsWithStateRootOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.ApplicationId = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Requests = *abi.ConvertType(out[1], new([]StructsPendingRequest)).(*[]StructsPendingRequest)
	outstruct.StateRoot = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	return *outstruct, nil
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

// PackSelectionGrace is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x44ce30f0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function selectionGrace() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackSelectionGrace() []byte {
	enc, err := processorEndpoint.abi.Pack("selectionGrace")
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
func (processorEndpoint *ProcessorEndpoint) TryPackSelectionGrace() ([]byte, error) {
	return processorEndpoint.abi.Pack("selectionGrace")
}

// UnpackSelectionGrace is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x44ce30f0.
//
// Solidity: function selectionGrace() view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackSelectionGrace(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("selectionGrace", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
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

// PackUpdateSelectionGrace is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x90fda89b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateSelectionGrace(uint256 newGrace) returns()
func (processorEndpoint *ProcessorEndpoint) PackUpdateSelectionGrace(newGrace *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("updateSelectionGrace", newGrace)
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
func (processorEndpoint *ProcessorEndpoint) TryPackUpdateSelectionGrace(newGrace *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("updateSelectionGrace", newGrace)
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

// ProcessorEndpointBatchProcessed represents a BatchProcessed event raised by the ProcessorEndpoint contract.
type ProcessorEndpointBatchProcessed struct {
	ApplicationId uint64
	EntryHashes   [][32]byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointBatchProcessedEventName = "BatchProcessed"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointBatchProcessed) ContractEventName() string {
	return ProcessorEndpointBatchProcessedEventName
}

// UnpackBatchProcessedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BatchProcessed(uint64 indexed applicationId, bytes32[] entryHashes)
func (processorEndpoint *ProcessorEndpoint) UnpackBatchProcessedEvent(log *types.Log) (*ProcessorEndpointBatchProcessed, error) {
	event := "BatchProcessed"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointBatchProcessed)
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

// ProcessorEndpointSelectionGraceUpdated represents a SelectionGraceUpdated event raised by the ProcessorEndpoint contract.
type ProcessorEndpointSelectionGraceUpdated struct {
	NewGrace *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointSelectionGraceUpdatedEventName = "SelectionGraceUpdated"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointSelectionGraceUpdated) ContractEventName() string {
	return ProcessorEndpointSelectionGraceUpdatedEventName
}

// UnpackSelectionGraceUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SelectionGraceUpdated(uint256 newGrace)
func (processorEndpoint *ProcessorEndpoint) UnpackSelectionGraceUpdatedEvent(log *types.Log) (*ProcessorEndpointSelectionGraceUpdated, error) {
	event := "SelectionGraceUpdated"
	if log.Topics[0] != processorEndpoint.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ProcessorEndpointSelectionGraceUpdated)
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ApplicationNotSelected"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackApplicationNotSelectedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["AuthorityNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackAuthorityNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["BatchNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackBatchNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["DeadlineExpired"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackDeadlineExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["DeployerNotAllowed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackDeployerNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["EmptyBatch"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackEmptyBatchError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidExtension"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidExtensionError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["PriorityQueueNotServed"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackPriorityQueueNotServedError(raw[4:])
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

// ProcessorEndpointApplicationNotSelected represents a ApplicationNotSelected error raised by the ProcessorEndpoint contract.
type ProcessorEndpointApplicationNotSelected struct {
	SelectedApplicationId uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ApplicationNotSelected(uint64 selectedApplicationId)
func ProcessorEndpointApplicationNotSelectedErrorID() common.Hash {
	return common.HexToHash("0x7041ef9ab6a10b8be16de04adfce6693e5fc57f5723f500b734d177d01ffdb67")
}

// UnpackApplicationNotSelectedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ApplicationNotSelected(uint64 selectedApplicationId)
func (processorEndpoint *ProcessorEndpoint) UnpackApplicationNotSelectedError(raw []byte) (*ProcessorEndpointApplicationNotSelected, error) {
	out := new(ProcessorEndpointApplicationNotSelected)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "ApplicationNotSelected", raw); err != nil {
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

// ProcessorEndpointBatchNotAllowed represents a BatchNotAllowed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointBatchNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchNotAllowed()
func ProcessorEndpointBatchNotAllowedErrorID() common.Hash {
	return common.HexToHash("0x1d49e5b0308cbb9d61f043b7942a66d9522a1c3f924d8868c91673df15a3ea59")
}

// UnpackBatchNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchNotAllowed()
func (processorEndpoint *ProcessorEndpoint) UnpackBatchNotAllowedError(raw []byte) (*ProcessorEndpointBatchNotAllowed, error) {
	out := new(ProcessorEndpointBatchNotAllowed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "BatchNotAllowed", raw); err != nil {
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

// ProcessorEndpointEmptyBatch represents a EmptyBatch error raised by the ProcessorEndpoint contract.
type ProcessorEndpointEmptyBatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBatch()
func ProcessorEndpointEmptyBatchErrorID() common.Hash {
	return common.HexToHash("0xc2e5347df6cf8d1bf68f8a9651642312bd381b900d857d97cd665976180d6ae0")
}

// UnpackEmptyBatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBatch()
func (processorEndpoint *ProcessorEndpoint) UnpackEmptyBatchError(raw []byte) (*ProcessorEndpointEmptyBatch, error) {
	out := new(ProcessorEndpointEmptyBatch)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "EmptyBatch", raw); err != nil {
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

// ProcessorEndpointInvalidExtension represents a InvalidExtension error raised by the ProcessorEndpoint contract.
type ProcessorEndpointInvalidExtension struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidExtension()
func ProcessorEndpointInvalidExtensionErrorID() common.Hash {
	return common.HexToHash("0x05d3b44e2aa04f6e552f2a9ce32a313664749082185211cb98f50856f9255ba8")
}

// UnpackInvalidExtensionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidExtension()
func (processorEndpoint *ProcessorEndpoint) UnpackInvalidExtensionError(raw []byte) (*ProcessorEndpointInvalidExtension, error) {
	out := new(ProcessorEndpointInvalidExtension)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "InvalidExtension", raw); err != nil {
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

// ProcessorEndpointPriorityQueueNotServed represents a PriorityQueueNotServed error raised by the ProcessorEndpoint contract.
type ProcessorEndpointPriorityQueueNotServed struct {
	QueueType uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PriorityQueueNotServed(uint8 queueType)
func ProcessorEndpointPriorityQueueNotServedErrorID() common.Hash {
	return common.HexToHash("0x643c8003e266b0b145803afa70ee15fc108315d76d7f34da61c7d909edba7551")
}

// UnpackPriorityQueueNotServedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PriorityQueueNotServed(uint8 queueType)
func (processorEndpoint *ProcessorEndpoint) UnpackPriorityQueueNotServedError(raw []byte) (*ProcessorEndpointPriorityQueueNotServed, error) {
	out := new(ProcessorEndpointPriorityQueueNotServed)
	if err := processorEndpoint.abi.UnpackIntoInterface(out, "PriorityQueueNotServed", raw); err != nil {
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
