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

// StructsPendingRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsPendingRequest struct {
	Timestamp       *big.Int
	DepositAmount   *big.Int
	MaxFeeValue     *big.Int
	RequestId       [32]byte
	Payload         []byte
	Sender          common.Address
	ApplicationId   uint64
	ProtocolVersion uint8
	RequestType     uint8
}

// StructsWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type StructsWithdrawalRequest struct {
	Receiver common.Address
	Amount   *big.Int
}

// ProcessorEndpointMetaData contains all meta data concerning the ProcessorEndpoint contract.
var ProcessorEndpointMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"_teeAuthenticator\",\"type\":\"address\"},{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"_authorityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"updateStatusOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minFeePerRequest\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressCantBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AuthorityNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeValueBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidApplicationId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayload\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProtocolVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStateRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QueueThresholdExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"QueueThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumStructs.RequestResult\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"RequestCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"StateRootUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"eventSubType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedData\",\"type\":\"bytes\"}],\"name\":\"UserEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"APPLICATION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_VERSION\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATE_STATUS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authorityRegistry\",\"outputs\":[{\"internalType\":\"contractIAuthorityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"generateRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextPendingRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequests\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getPendingRequestsPage\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.PendingRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestsSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"isCurrentPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"}],\"name\":\"markRequestCompleted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"enumStructs.ErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"}],\"name\":\"markRequestFailed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxQueueSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minFeePerRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"payments\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"requestById\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"prevStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"processedRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"events\",\"type\":\"bytes[]\"},{\"internalType\":\"string[]\",\"name\":\"eventSubTypes\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.WithdrawalRequest[]\",\"name\":\"withdrawalRequests\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"applicationFees\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"stateUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"protocolVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"applicationId\",\"type\":\"uint64\"},{\"internalType\":\"enumStructs.RequestType\",\"name\":\"requestType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeValue\",\"type\":\"uint256\"}],\"name\":\"submitRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"teeAuthenticator\",\"outputs\":[{\"internalType\":\"contractITeeAuthenticator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateQueueThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"withdrawPayments\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "ProcessorEndpoint",
	Bin: "0x6080604052346100335761001d6100146101b4565b93929092610405565b610025610038565b614a416107238239614a4190f35b61003e565b60405190565b5f80fd5b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b9061006a90610042565b810190811060018060401b0382111761008257604052565b61004c565b9061009a610093610038565b9283610060565b565b5f80fd5b60018060a01b031690565b6100b4906100a0565b90565b6100c0906100ab565b90565b6100cc816100b7565b036100d357565b5f80fd5b905051906100e4826100c3565b565b6100ef906100ab565b90565b6100fb816100e6565b0361010257565b5f80fd5b90505190610113826100f2565b565b61011e816100ab565b0361012557565b5f80fd5b9050519061013682610115565b565b90565b61014481610138565b0361014b57565b5f80fd5b9050519061015c8261013b565b565b919060a0838203126101af57610176815f85016100d7565b926101848260208301610106565b926101ac6101958460408501610129565b936101a38160608601610129565b9360800161014f565b90565b61009c565b6101d2615164803803806101c781610087565b92833981019061015e565b9091929394565b5f1b90565b906101ea5f19916101d9565b9181191691161790565b90565b90565b61020e610209610213926101f4565b6101f7565b610138565b90565b90565b9061022e610229610235926101fa565b610216565b82546101de565b9055565b90565b61025061024b61025592610239565b6101f7565b6100a0565b90565b6102619061023c565b90565b61027861027361027d926100a0565b6101f7565b6100a0565b90565b61028990610264565b90565b61029590610280565b90565b6102a190610264565b90565b6102ad90610298565b90565b5f0190565b906102c660018060a01b03916101d9565b9181191691161790565b6102d990610280565b90565b90565b906102f46102ef6102fb926102d0565b6102dc565b82546102b5565b9055565b61030890610264565b90565b610314906102ff565b90565b90565b9061032f61032a6103369261030b565b610317565b82546102b5565b9055565b61034390610264565b90565b61034f9061033a565b90565b61035b9061033a565b90565b90565b9061037661037161037d92610352565b61035e565b82546102b5565b9055565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6103dd6103d86103e292610138565b6101f7565b610138565b90565b906103fa6103f5610401926103c9565b610216565b82546101de565b9055565b91939290610411610568565b61041d600a6007610219565b8261044061043a6104356104305f610258565b61028c565b6100b7565b916100b7565b148015610512575b80156104f0575b80156104ce575b6104b2576104b09461047a61049a926104736104a89660086102df565b600961031a565b61048d61048682610346565b600d610361565b610495610381565b610611565b506104a36103a5565b610611565b50600c6103e5565b565b5f632582a64160e11b8152806104ca600482016102b0565b0390fd5b50816104ea6104e46104df5f610258565b6100ab565b916100ab565b14610456565b508461050c6105066105015f610258565b6100ab565b916100ab565b1461044f565b5061051c816102a4565b61053661053061052b5f610258565b6100ab565b916100ab565b14610448565b90565b61055361054e6105589261053c565b6101f7565b610138565b90565b610565600161053f565b90565b61057a61057361055b565b60016103e5565b565b5f90565b151590565b90565b61059190610585565b90565b9061059e90610588565b5f5260205260405f2090565b6105b390610298565b90565b906105c0906105aa565b5f5260205260405f2090565b906105d860ff916101d9565b9181191691161790565b6105eb90610580565b90565b90565b9061060661060161060d926105e2565b6105ee565b82546105cc565b9055565b61061961057c565b5061062e6106288284906106e8565b15610580565b5f146106b65761065560016106505f610648818690610594565b0185906105b6565b6105f1565b9061065e610715565b9061069b61069561068f7f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610588565b926105aa565b926105aa565b926106a4610038565b806106ae816102b0565b0390a4600190565b50505f90565b5f1c90565b60ff1690565b6106d36106d8916106bc565b6106c1565b90565b6106e590546106c7565b90565b61070e915f610703610709936106fc61057c565b5082610594565b016105b6565b6106db565b90565b5f90565b61071d610711565b50339056fe6101606040526004361015610014575b6119e8565b61001e5f3561021d565b806301ffc9a7146102185780631664deeb14610213578063248a9ca31461020e5780632a0acc6a146102095780632f2ff15d1461020457806331b3eb94146101ff57806334f23a9d146101fa57806336568abe146101f55780633b473cb9146101f05780635423112f146101eb5780635d1be609146101e65780635e339384146101e15780636e91d133146101dc57806380a1f712146101d757806391d14854146101d25780639588eca2146101cd5780639968da96146101c85780639a17995a146101c35780639c73eeaf146101be5780639c8d416f146101b95780639fabc36f146101b4578063a217fddf146101af578063a9ba045e146101aa578063a9ef290e146101a5578063aa3aa460146101a0578063be1e37251461019b578063c404c35914610196578063c415b95c14610191578063d2c35ce81461018c578063d547741f14610187578063e2982c21146101825763f7d9cc1e0361000f576119b3565b61196f565b6118e2565b6118af565b61187a565b6117e0565b6116c8565b611666565b6115d8565b6113e0565b61133d565b6112cd565b611298565b611230565b6111d6565b611121565b6110a0565b611035565b611000565b610dfa565b610d28565b610cdf565b610c66565b6107eb565b610709565b6106d7565b61053d565b6104bc565b61041b565b6103b7565b610341565b6102a9565b60e01c90565b60405190565b5f80fd5b5f80fd5b5f80fd5b63ffffffff60e01b1690565b61024a81610235565b0361025157565b5f80fd5b9050359061026282610241565b565b9060208282031261027d5761027a915f01610255565b90565b61022d565b151590565b61029090610282565b9052565b91906102a7905f60208501940190610287565b565b346102d9576102d56102c46102bf366004610264565b6119f0565b6102cc610223565b91829182610294565b0390f35b610229565b5f9103126102e857565b61022d565b7f20ddaa8e07c23db8d16523a387fb4a6749f3c49bd4ff22512336fdd38dc6dd9890565b6103196102ed565b90565b90565b6103289061031c565b9052565b919061033f905f6020850194019061031f565b565b34610371576103513660046102de565b61036d61035c610311565b610364610223565b9182918261032c565b0390f35b610229565b61037f8161031c565b0361038657565b5f80fd5b9050359061039782610376565b565b906020828203126103b2576103af915f0161038a565b90565b61022d565b346103e7576103e36103d26103cd366004610399565b611a4a565b6103da610223565b9182918261032c565b0390f35b610229565b7fdf8b4c520ffe197c5343c6f5aec59570151ef9a492f2c624fd45ddde6135ec4290565b6104186103ec565b90565b3461044b5761042b3660046102de565b610447610436610410565b61043e610223565b9182918261032c565b0390f35b610229565b60018060a01b031690565b61046490610450565b90565b6104708161045b565b0361047757565b5f80fd5b9050359061048882610467565b565b91906040838203126104b257806104a66104af925f860161038a565b9360200161047b565b90565b61022d565b5f0190565b346104eb576104d56104cf36600461048a565b90611a95565b6104dd610223565b806104e7816104b7565b0390f35b610229565b6104f990610450565b90565b610505816104f0565b0361050c57565b5f80fd5b9050359061051d826104fc565b565b9060208282031261053857610535915f01610510565b90565b61022d565b3461056b5761055561055036600461051f565b611d2c565b61055d610223565b80610567816104b7565b0390f35b610229565b60ff1690565b61057f81610570565b0361058657565b5f80fd5b9050359061059782610576565b565b67ffffffffffffffff1690565b6105af81610599565b036105b657565b5f80fd5b905035906105c7826105a6565b565b600411156105d357565b5f80fd5b905035906105e4826105c9565b565b5f80fd5b5f80fd5b5f80fd5b909182601f8301121561062c5781359167ffffffffffffffff831161062757602001926001830284011161062257565b6105ee565b6105ea565b6105e6565b90565b61063d81610631565b0361064457565b5f80fd5b9050359061065582610634565b565b91909160c0818403126106d257610670835f830161058a565b9261067e81602084016105ba565b9261068c82604085016105d7565b9260608101359167ffffffffffffffff83116106cd576106b1846106ca9484016105f2565b9390946106c18160808601610648565b9360a001610648565b90565b610231565b61022d565b6107056106f46106e8366004610657565b95949094939193612810565b6106fc610223565b9182918261032c565b0390f35b346107385761072261071c36600461048a565b9061282a565b61072a610223565b80610734816104b7565b0390f35b610229565b6011111561074757565b5f80fd5b905035906107588261073d565b565b909182601f830112156107945781359167ffffffffffffffff831161078f57602001926001830284011161078a57565b6105ee565b6105ea565b6105e6565b916060838303126107e6576107b0825f850161038a565b926107be836020830161074b565b92604082013567ffffffffffffffff81116107e1576107dd920161075a565b9091565b610231565b61022d565b3461081d576108076107fe366004610799565b92919091612c6c565b61080f610223565b80610819816104b7565b0390f35b610229565b61082b9061031c565b90565b9061083890610822565b5f5260205260405f2090565b5f1c90565b90565b61085861085d91610844565b610849565b90565b61086a905461084c565b90565b90565b61087c61088191610844565b61086d565b90565b61088e9054610870565b90565b634e487b7160e01b5f52602260045260245ffd5b90600160028304921680156108c5575b60208310146108c057565b610891565b91607f16916108b5565b60209181520190565b5f5260205f2090565b905f92918054906108fb6108f4836108a5565b80946108cf565b916001811690815f146109525750600114610916575b505050565b61092391929394506108d8565b915f925b81841061093a57505001905f8080610911565b60018160209295939554848601520191019290610927565b92949550505060ff19168252151560200201905f8080610911565b90610977916108e1565b90565b601f801991011690565b634e487b7160e01b5f52604160045260245ffd5b906109a29061097a565b810190811067ffffffffffffffff8211176109bc57604052565b610984565b906109e16109da926109d1610223565b9384809261096d565b0383610998565b565b60018060a01b031690565b6109fa6109ff91610844565b6109e3565b90565b610a0c90546109ee565b90565b60a01c90565b67ffffffffffffffff1690565b610a2e610a3391610a0f565b610a15565b90565b610a409054610a22565b90565b60ff1690565b610a55610a5a9161021d565b610a43565b90565b610a679054610a49565b90565b60e81c90565b60ff1690565b610a82610a8791610a6a565b610a70565b90565b610a949054610a76565b90565b610aa290600361082e565b610aad5f8201610860565b91610aba60018301610860565b91610ac760028201610860565b91610ad460038301610884565b91610ae1600482016109c1565b91610aee60058301610a02565b91610afb60058201610a36565b91610b136005610b0c818501610a5d565b9301610a8a565b90565b610b1f90610631565b9052565b5190565b60209181520190565b90825f9392825e0152565b610b5a610b63602093610b6893610b5181610b23565b93848093610b27565b95869101610b30565b61097a565b0190565b610b759061045b565b9052565b610b8290610599565b9052565b610b8f90610570565b9052565b634e487b7160e01b5f52602160045260245ffd5b60041115610bb157565b610b93565b90610bc082610ba7565b565b610bcb90610bb6565b90565b610bd790610bc2565b9052565b94610c3e610100979b9a9892610c4992610c31610c5d98610c27610c649e99610c1d610c539a60208f610c1661012082019a5f830190610b16565b0190610b16565b60408d0190610b16565b60608b019061031f565b88820360808a0152610b3b565b9a60a0870190610b6c565b60c0850190610b79565b60e0830190610b86565b0190610bce565b565b34610ca057610c9c610c81610c7c366004610399565b610a97565b95610c93999799959195949294610223565b998a998a610bdb565b0390f35b610229565b9091606082840312610cda57610cd7610cc0845f850161038a565b93610cce8160208601610648565b93604001610648565b90565b61022d565b34610d0e57610cf8610cf2366004610ca5565b91612e52565b610d00610223565b80610d0a816104b7565b0390f35b610229565b9190610d26905f60208501940190610b16565b565b34610d5857610d383660046102de565b610d54610d43612e5f565b610d4b610223565b91829182610d13565b0390f35b610229565b1c90565b60018060a01b031690565b610d7c906008610d819302610d5d565b610d61565b90565b90610d8f9154610d6c565b90565b610d9e60085f90610d84565b90565b90565b610db8610db3610dbd92610450565b610da1565b610450565b90565b610dc990610da4565b90565b610dd590610dc0565b90565b610de190610dcc565b9052565b9190610df8905f60208501940190610dd8565b565b34610e2a57610e0a3660046102de565b610e26610e15610d92565b610e1d610223565b91829182610de5565b0390f35b610229565b5190565b60209181520190565b60200190565b610e4b90610631565b9052565b610e589061031c565b9052565b610e7b610e84602093610e8993610e7281610b23565b938480936108cf565b95869101610b30565b61097a565b0190565b610e969061045b565b9052565b610ea390610599565b9052565b610eb090610570565b9052565b610ebd90610bc2565b9052565b90610f6b9061010080610f2a6101208401610ee25f8801515f870190610e42565b610ef460208801516020870190610e42565b610f0660408801516040870190610e42565b610f1860608801516060870190610e4f565b60808701518582036080870152610e5c565b94610f3d60a082015160a0860190610e8d565b610f4f60c082015160c0860190610e9a565b610f6160e082015160e0860190610ea7565b0151910190610eb4565b90565b90610f7891610ec1565b90565b60200190565b90610f95610f8e83610e2f565b8092610e33565b9081610fa660208302840194610e3c565b925f915b838310610fb957505050505090565b90919293946020610fdb610fd583856001950387528951610f6e565b97610f7b565b9301930191939290610faa565b610ffd9160208201915f818403910152610f81565b90565b34611030576110103660046102de565b61102c61101b613019565b611023610223565b91829182610fe8565b0390f35b610229565b346110665761106261105161104b36600461048a565b90613107565b611059610223565b91829182610294565b0390f35b610229565b61107b9060086110809302610d5d565b61086d565b90565b9061108e915461106b565b90565b61109d60025f90611083565b90565b346110d0576110b03660046102de565b6110cc6110bb611091565b6110c3610223565b9182918261032c565b0390f35b610229565b90565b6110ec6110e76110f1926110d5565b610da1565b610599565b90565b6110fe60016110d8565b90565b6111096110f4565b90565b919061111f905f60208501940190610b79565b565b34611151576111313660046102de565b61114d61113c611101565b611144610223565b9182918261110c565b0390f35b610229565b91909160c0818403126111d15761116f835f830161047b565b9261117d81602084016105ba565b9261118b82604085016105d7565b9260608101359167ffffffffffffffff83116111cc576111b0846111c99484016105f2565b9390946111c08160808601610648565b9360a001610648565b90565b610231565b61022d565b3461120d576112096111f86111ec366004611156565b9594909493919361323d565b611200610223565b9182918261032c565b0390f35b610229565b9060208282031261122b57611228915f01610648565b90565b61022d565b3461125e57611248611243366004611212565b613327565b611250610223565b8061125a816104b7565b0390f35b610229565b6112739060086112789302610d5d565b610849565b90565b906112869154611263565b90565b611295600c5f9061127b565b90565b346112c8576112a83660046102de565b6112c46112b3611289565b6112bb610223565b91829182610d13565b0390f35b610229565b346112fd576112f96112e86112e3366004610399565b613332565b6112f0610223565b91829182610294565b0390f35b610229565b90565b5f1b90565b61131e61131961132392611302565b611305565b61031c565b90565b61132f5f61130a565b90565b61133a611326565b90565b3461136d5761134d3660046102de565b611369611358611332565b611360610223565b9182918261032c565b0390f35b610229565b60018060a01b031690565b61138d9060086113929302610d5d565b611372565b90565b906113a0915461137d565b90565b6113af60095f90611395565b90565b6113bb90610dc0565b90565b6113c7906113b2565b9052565b91906113de905f602085019401906113be565b565b34611410576113f03660046102de565b61140c6113fb6113a3565b611403610223565b918291826113cb565b0390f35b610229565b909182601f8301121561144f5781359167ffffffffffffffff831161144a57602001926020830284011161144557565b6105ee565b6105ea565b6105e6565b909182601f8301121561148e5781359167ffffffffffffffff831161148957602001926020830284011161148457565b6105ee565b6105ea565b6105e6565b909182601f830112156114cd5781359167ffffffffffffffff83116114c85760200192604083028401116114c357565b6105ee565b6105ea565b6105e6565b9190610140838203126115d3576114eb815f85016105ba565b926114f9826020830161038a565b92611507836040840161038a565b92611515816060850161038a565b92608081013567ffffffffffffffff81116115ce5782611536918301611415565b92909360a083013567ffffffffffffffff81116115c95782611559918501611454565b92909360c081013567ffffffffffffffff81116115c4578261157c918301611493565b92909361158c8260e08501610648565b9261159b836101008301610648565b9261012082013567ffffffffffffffff81116115bf576115bb92016105f2565b9091565b610231565b610231565b610231565b610231565b61022d565b34611619576116036115eb3660046114d2565b9c9b909b9a919a999299989398979497969596614022565b61160b610223565b80611615816104b7565b0390f35b610229565b61163261162d61163792611302565b610da1565b610570565b90565b6116435f61161e565b90565b61164e61163a565b90565b9190611664905f60208501940190610b86565b565b34611696576116763660046102de565b611692611681611646565b611689610223565b91829182611651565b0390f35b610229565b91906040838203126116c357806116b76116c0925f8601610648565b93602001610648565b90565b61022d565b346116f9576116f56116e46116de36600461169b565b9061403a565b6116ec610223565b91829182610fe8565b0390f35b610229565b906117a89061010080611767610120840161171f5f8801515f870190610e42565b61173160208801516020870190610e42565b61174360408801516040870190610e42565b61175560608801516060870190610e4f565b60808701518582036080870152610e5c565b9461177a60a082015160a0860190610e8d565b61178c60c082015160c0860190610e9a565b61179e60e082015160e0860190610ea7565b0151910190610eb4565b90565b6040906117d76117cc6117de9597969460608401908482035f8601526116fe565b96602083019061031f565b0190610287565b565b34611813576117f03660046102de565b61180f6117fb614196565b611806939193610223565b938493846117ab565b0390f35b610229565b60018060a01b031690565b6118339060086118389302610d5d565b611818565b90565b906118469154611823565b90565b611855600d5f9061183b565b90565b611861906104f0565b9052565b9190611878905f60208501940190611858565b565b346118aa5761188a3660046102de565b6118a6611895611849565b61189d610223565b91829182611865565b0390f35b610229565b346118dd576118c76118c236600461051f565b614335565b6118cf610223565b806118d9816104b7565b0390f35b610229565b34611911576118fb6118f536600461048a565b9061436a565b611903610223565b8061190d816104b7565b0390f35b610229565b9060208282031261192f5761192c915f0161047b565b90565b61022d565b61193d90610dc0565b90565b9061194a90611934565b5f5260205260405f2090565b61196c90611967600a915f92611940565b61127b565b90565b3461199f5761199b61198a611985366004611916565b611956565b611992610223565b91829182610d13565b0390f35b610229565b6119b060075f9061127b565b90565b346119e3576119c33660046102de565b6119df6119ce6119a4565b6119d6610223565b91829182610d13565b0390f35b610229565b5f80fd5b5f90565b6119f86119ec565b5080611a13611a0d637965db0b60e01b610235565b91610235565b14908115611a20575b5090565b611a2a9150614376565b5f611a1c565b5f90565b90611a3e90610822565b5f5260205260405f2090565b6001611a62611a6892611a5b611a30565b505f611a34565b01610884565b90565b90611a8691611a81611a7c82611a4a565b61439c565b611a88565b565b90611a92916143f5565b50565b90611a9f91611a6b565b565b611ab290611aad6144cc565b611c25565b611aba61454d565b565b611ac590610dc0565b90565b90611ad290611abc565b5f5260205260405f2090565b611af2611aed611af792611302565b610da1565b610631565b90565b90611b065f1991611305565b9181191691161790565b611b24611b1f611b2992610631565b610da1565b610631565b90565b90565b90611b44611b3f611b4b92611b10565b611b2c565b8254611afa565b9055565b634e487b7160e01b5f52601160045260245ffd5b611b72611b7891939293610631565b92610631565b8203918211611b8357565b611b4f565b905090565b611b985f8092611b88565b0190565b611ba590611b8d565b90565b90611bbb611bb4610223565b9283610998565b565b67ffffffffffffffff8111611bdb57611bd760209161097a565b0190565b610984565b90611bf2611bed83611bbd565b611ba8565b918252565b606090565b3d5f14611c1757611c0c3d611be0565b903d5f602084013e5b565b611c1f611bf7565b90611c15565b611c39611c34600a8390611ac8565b610860565b80611c4c611c465f611ade565b91610631565b14611d28575f8091611cdf611d0694611c78611c6785611ade565b611c73600a8490611ac8565b611b2f565b611c95611c8e84611c89600b610860565b611b63565b600b611b2f565b808390611cd7611cc57f84511ecc081974f18e7f3e0dcc19db078b55bbd3852ddd0dd85b3aebb7bf94c292611abc565b92611cce610223565b91829182610d13565b0390a2611abc565b90611ce8610223565b9081611cf381611b9c565b03925af1611cff611bfc565b5015610282565b611d0c57565b5f6312171d8360e31b815280611d24600482016104b7565b0390fd5b5050565b611d3590611aa1565b565b9695949392919080611d58611d52611d4d61163a565b610570565b91610570565b03611d6957611d6697611d85565b90565b5f635428eae760e01b815280611d81600482016104b7565b0390fd5b9695949392919081611da6611da0611d9b6110f4565b610599565b91610599565b03611db757611db497611dd3565b90565b5f6309ff9a8360e41b815280611dcf600482016104b7565b0390fd5b90611deb97969594939291611de66144cc565b61248a565b90611df461454d565b565b611e05611e0b91939293610631565b92610631565b8201809211611e1657565b611b4f565b611e27611e2c91610844565b611372565b90565b611e399054611e1b565b90565b60e01b90565b611e4b81610282565b03611e5257565b5f80fd5b90505190611e6382611e42565b565b90602082820312611e7e57611e7b915f01611e56565b90565b61022d565b611e97611e92611e9c92610599565b610da1565b610631565b90565b611ea890611e83565b9052565b916020611ecd929493611ec660408201965f830190611e9f565b0190610b6c565b565b611ed7610223565b3d5f823e3d90fd5b5090565b90565b611efa611ef5611eff92611ee3565b610da1565b610631565b90565b611f0d610120611ba8565b90565b90611f1a90610631565b9052565b90611f289061031c565b9052565b5f80fd5b90825f939282370152565b90929192611f50611f4b82611bbd565b611ba8565b93818552602085019082840111611f6c57611f6a92611f30565b565b611f2c565b611f7c913691611f3b565b90565b52565b90611f8c9061045b565b9052565b90611f9a90610599565b9052565b90611fa890610570565b9052565b90611fb690610bb6565b9052565b634e487b7160e01b5f525f60045260245ffd5b611fd79051610631565b90565b611fe4905161031c565b90565b611ff090610844565b90565b9061200861200361200f92610822565b611fe7565b8254611afa565b9055565b5190565b601f602091010490565b1b90565b9190600861204091029161203a5f1984612021565b92612021565b9181191691161790565b919061206061205b61206893611b10565b611b2c565b908354612025565b9055565b5f90565b6120829161207c61206c565b9161204a565b565b5b818110612090575050565b8061209d5f600193612070565b01612085565b9190601f81116120b3575b505050565b6120bf6120e4936108d8565b9060206120cb84612017565b830193106120ec575b6120dd90612017565b0190612084565b5f80806120ae565b91506120dd819290506120d4565b9061210a905f1990600802610d5d565b191690565b81612119916120fa565b906002021790565b9061212b81610b23565b9067ffffffffffffffff82116121eb5761214f8261214985546108a5565b856120a3565b602090601f831160011461218357918091612172935f92612177575b505061210f565b90555b565b90915001515f8061216b565b601f19831691612192856108d8565b925f5b8181106121d3575091600293918560019694106121b9575b50505002019055612175565b6121c9910151601f8416906120fa565b90555f80806121ad565b91936020600181928787015181550195019201612195565b610984565b906121fa91612121565b565b612206905161045b565b90565b9061221a60018060a01b0391611305565b9181191691161790565b90565b9061223c61223761224392611934565b612224565b8254612209565b9055565b6122519051610599565b90565b60a01b90565b9061227067ffffffffffffffff60a01b91612254565b9181191691161790565b61228e61228961229392610599565b610da1565b610599565b90565b90565b906122ae6122a96122b59261227a565b612296565b825461225a565b9055565b6122c39051610570565b90565b906122d560ff60e01b91611e3c565b9181191691161790565b6122f36122ee6122f892610570565b610da1565b610570565b90565b90565b9061231361230e61231a926122df565b6122fb565b82546122c6565b9055565b6123289051610bb6565b90565b60e81b90565b9061234060ff60e81b9161232b565b9181191691161790565b61235390610bb6565b90565b90565b9061236e6123696123759261234a565b612356565b8254612331565b9055565b9061245161010060056124579461239d5f82016123975f8801611fcd565b90611b2f565b6123b6600182016123b060208801611fcd565b90611b2f565b6123cf600282016123c960408801611fcd565b90611b2f565b6123e8600382016123e260608801611fda565b90611ff3565b612401600482016123fb60808801612013565b906121f0565b61241982820161241360a088016121fc565b90612227565b61243182820161242b60c08801612247565b90612299565b61244982820161244360e088016122b9565b906122fe565b01920161231e565b90612359565b565b9061246391612379565b565b9061246f90611b10565b5f5260205260405f2090565b60016124879101610631565b90565b96949092939650346124ae6124a86124a3898990611df6565b610631565b91610631565b036127f457846124cf6124c96124c4600c610860565b610631565b91610631565b106127d8576124dc612e5f565b6124f76124f16124ec6007610860565b610631565b91610631565b10156127bc578361251161250b6003610bb6565b91610bb6565b145f146126b657612523878290611edf565b6125366125306085611ee6565b91610631565b0361269a575b3382858984908a9261254e6006610860565b946125589661323d565b964296959088909291339495969761256e611f02565b995f8b019061257c91611f10565b60208a019061258a91611f10565b604089019061259891611f10565b60608801906125a691611f1e565b6125af91611f71565b60808601906125bd91611f7f565b60a08501906125cb91611f82565b60c08401906125d991611f90565b60e08301906125e791611f9e565b6101008201906125f691611fac565b6003826126029161082e565b9061260c91612459565b8060046126196006610860565b61262291612465565b9061262c91611ff3565b6126366006610860565b61263f9061247b565b61264a906006611b2f565b80337f73394f43049193ecd2e0c22eefa0ecf10987ce9c72da906a82a3336dc2ac05479161267790610822565b9061268190611934565b9161268a610223565b80612694816104b7565b0390a390565b5f637c6953f960e01b8152806126b2600482016104b7565b0390fd5b836126ca6126c46002610bb6565b91610bb6565b146126d5575b61253c565b856126e86126e25f611ade565b91610631565b036127a0576126ff6126fa6009611e2f565b6113b2565b602063116ad4a69184906127253394612730612719610223565b96879586948594611e3c565b845260048401611eac565b03915afa801561279b5761274c915f9161276d575b5015610282565b156126d0575f63622a850b60e11b815280612769600482016104b7565b0390fd5b61278e915060203d8111612794575b6127868183610998565b810190611e65565b5f612745565b503d61277c565b611ecf565b5f632a9ffab760e21b8152806127b8600482016104b7565b0390fd5b5f63d8219b3d60e01b8152806127d4600482016104b7565b0390fd5b5f63c77f97eb60e01b8152806127f0600482016104b7565b0390fd5b5f632a9ffab760e21b81528061280c600482016104b7565b0390fd5b90612827969594939291612822611a30565b611d37565b90565b908061284561283f61283a614565565b61045b565b9161045b565b036128565761285391614572565b50565b5f63334bd91960e11b81528061286e600482016104b7565b0390fd5b9061288e9392916128896128846102ed565b61439c565b612aca565b565b61289b610120611ba8565b90565b9061297b61297160056128af612890565b946128c66128be5f8301610860565b5f8801611f10565b6128de6128d560018301610860565b60208801611f10565b6128f66128ed60028301610860565b60408801611f10565b61290e61290560038301610884565b60608801611f1e565b61292661291d600483016109c1565b60808801611f7f565b61293d612934838301610a02565b60a08801611f82565b61295461294b838301610a36565b60c08801611f90565b61296b612962838301610a5d565b60e08801611f9e565b01610a8a565b6101008401611fac565b565b6129869061289e565b90565b61299290610da4565b90565b61299e90612989565b90565b6129aa90611abc565b9052565b9160206129cf9294936129c860408201965f8301906129a1565b0190610b16565b565b6129dd6129e291610844565b611818565b90565b6129ef90546129d1565b90565b600211156129fc57565b610b93565b90612a0b826129f2565b565b612a1690612a01565b90565b612a2290612a0d565b9052565b60111115612a3057565b610b93565b90612a3f82612a26565b565b612a4a90612a35565b90565b612a5690612a41565b9052565b60209181520190565b9190612a7d81612a7681612a8295612a5a565b8095611f30565b61097a565b0190565b909391612ac79593612ab0612aba92612aa660808601985f870190610b16565b6020850190612a19565b6040830190612a4d565b6060818503910152612a63565b90565b929192612adf612ad982613332565b15610282565b612c5057612af7612af26003839061082e565b61297d565b612b33612b04600c610860565b612b0c614859565b612b2d612b1b60208501611fcd565b91612b2860408601611fcd565b611b63565b90611df6565b9081612b47612b415f611ade565b91610631565b11612bc8575b5050612b74612b5c600d6129e5565b612b6f612b69600c610860565b91611abc565b6148c6565b91612bc3612b82600c610860565b9160019395612bb17f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198896610822565b96612bba610223565b95869586612a86565b0390a2565b612bfb60c0612be1612bdc60a085016121fc565b612995565b92612bf584612bf08791611abc565b6148c6565b01612247565b839192612c31612c2b7f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f9361227a565b93610822565b93612c46612c3d610223565b928392836129ae565b0390a35f80612b4d565b5f6302e8145360e61b815280612c68600482016104b7565b0390fd5b90612c78939291612872565b565b90612c959291612c90612c8b6102ed565b61439c565b612cbd565b565b90565b916020612cbb929493612cb460408201965f830190610b6c565b0190610b16565b565b919091612cd2612ccc82613332565b15610282565b612e3657612cea612ce56003839061082e565b612c97565b92612cf760028501610860565b612d14612d0e612d08848790611df6565b92610631565b91610631565b03612e1a5782612d35612d2f612d2a600c610860565b610631565b91610631565b10612dfe57612d7a9381612d51612d4b5f611ade565b91610631565b11612d7c575b5050612d75612d66600d6129e5565b612d708491611abc565b6148c6565b61495e565b565b612d92612d8b60058301610a02565b83906148c6565b612d9e60058201610a36565b612dab6005859301610a02565b92612ddf612dd97f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f9361227a565b93610822565b93612df4612deb610223565b92839283612c9a565b0390a35f80612d57565b5f632a9ffab760e21b815280612e16600482016104b7565b0390fd5b5f632a9ffab760e21b815280612e32600482016104b7565b0390fd5b5f6302e8145360e61b815280612e4e600482016104b7565b0390fd5b90612e5d9291612c7a565b565b612e6761206c565b50612e726006610860565b612e8d612e87612e826005610860565b610631565b91610631565b115f14612eb457612eb1612ea16006610860565b612eab6005610860565b90611b63565b90565b612ebd5f611ade565b90565b606090565b67ffffffffffffffff8111612edd5760208091020190565b610984565b90612ef4612eef83612ec5565b611ba8565b918252565b5f90565b5f90565b606090565b5f90565b5f90565b5f90565b5f90565b612f1e612890565b90602080808080808080808a612f32612ef9565b815201612f3d612ef9565b815201612f48612ef9565b815201612f53612efd565b815201612f5e612f01565b815201612f69612f06565b815201612f74612f0a565b815201612f7f612f0e565b815201612f8a612f12565b81525050565b612f98612f16565b90565b5f5b828110612fa957505050565b602090612fb4612f90565b8184015201612f9d565b90612fe3612fcb83612ee2565b92602080612fd98693612ec5565b9201910390612f9b565b565b634e487b7160e01b5f52603260045260245ffd5b9061300382610e2f565b811015613014576020809102010190565b612fe5565b613021612ec0565b5061303261302d612e5f565b612fbe565b9061303d6005610860565b916130486006610860565b9161305161206c565b935b8061306661306086610631565b91610631565b10156130c2576130b66130bc916130af61309461308d61308860048590612465565b610884565b600361082e565b8661309f8a9261297d565b6130a98383612ff9565b52612ff9565b515061247b565b9461247b565b93613053565b509250905090565b906130d490611934565b5f5260205260405f2090565b60ff1690565b6130f26130f791610844565b6130e0565b90565b61310490546130e6565b90565b61312d915f6131226131289361311b6119ec565b5082611a34565b016130ca565b6130fa565b90565b60601b90565b61313f90613130565b90565b61314b90613136565b90565b61315a61315f9161045b565b613142565b9052565b60c01b90565b61317290613163565b90565b61318161318691610599565b613169565b9052565b60f81b90565b6131999061318a565b90565b6131a86131ad91610bc2565b613190565b9052565b9091826131c1816131c893611b88565b8093611f30565b0190565b90565b6131db6131e091610631565b6131cc565b9052565b93613233956001602099989561321d600861322b9761321560148f9c61320d816132249c61314e565b018092613175565b01809261319c565b01916131b1565b80926131cf565b0180926131cf565b0190565b60200190565b9461326e939296919461327d96613252611a30565b5095979390919293613262610223565b988997602089016131e4565b60208201810382520382610998565b61328f61328982610b23565b91613237565b2090565b6132ac906132a76132a26103ec565b61439c565b6132ae565b565b806132c16132bb5f611ade565b91610631565b1461330b576132d1816007611b2f565b6133067e3dae5ee1558d7e2d601240b20a335ff2ec1918dc09dede81826b632ed38db1916132fd610223565b91829182610d13565b0390a1565b5f632a9ffab760e21b815280613323600482016104b7565b0390fd5b61333090613293565b565b61333a6119ec565b50613343612e5f565b61335561334f5f611ade565b91610631565b119081613361575b5090565b905061339261338c613386613381600461337b6005610860565b90612465565b610884565b9261031c565b9161031c565b145f61335d565b9c9b9a999897969594939291908d6133c06133ba6133b56110f4565b610599565b91610599565b036133d0576133ce9d6133ec565b565b5f6309ff9a8360e41b8152806133e8600482016104b7565b0390fd5b906134129d9c9b9a99989796959493929161340d6134086102ed565b61439c565b61398a565b565b5090565b5090565b61342861342d91610844565b610d61565b90565b61343a905461341c565b90565b60209181520190565b90565b91906134638161345c81613468956108cf565b8095611f30565b61097a565b0190565b906134779291613449565b90565b5f80fd5b5f80fd5b5f80fd5b90356001602003823603038112156134c757016020813591019167ffffffffffffffff82116134c25760018202360383136134bd57565b61347e565b61347a565b613482565b60200190565b91816134dd9161343d565b90816134ee60208302840194613446565b92835f925b8484106135035750505050505090565b909192939495602061352f61352983856001950388526135238b88613486565b9061346c565b986134cc565b9401940192949391906134f3565b60209181520190565b90565b60209181520190565b919061356c816135658161357195613549565b8095611f30565b61097a565b0190565b906135809291613552565b90565b90356001602003823603038112156135c457016020813591019167ffffffffffffffff82116135bf5760018202360383136135ba57565b61347e565b61347a565b613482565b60200190565b91816135da9161353d565b90816135eb60208302840194613546565b92835f925b8484106136005750505050505090565b909192939495602061362c61362683856001950388526136208b88613583565b90613575565b986135c9565b9401940192949391906135f0565b60209181520190565b90565b50613655906020810190610510565b90565b613661906104f0565b9052565b50613674906020810190610648565b90565b9060206136a26136aa936136996136905f830183613646565b5f860190613658565b82810190613665565b910190610e42565b565b906136b981604093613677565b0190565b5090565b60400190565b916136d5826136db9261363a565b92613643565b90815f905b8282106136ee575050505090565b9091929361371061370a60019261370588866136bd565b6136ac565b956136c1565b9201909291926136e0565b91906137358161372e8161373a95610b27565b8095611f30565b61097a565b0190565b99936137e99e9c98916137c5979c9e9c6137db9b9761379b6137a9948f6060906137946137d09f9a9b61378a6137b79d61378061014086019b5f870190610b79565b602085019061031f565b604083019061031f565b019061031f565b8d60808185039101526134d2565b918a830360a08c01526135cf565b9187830360c08901526136c7565b9660e0850190610b16565b610100830190610b16565b61012081850391015261371b565b90565b5090565b9190811015613800576040020190565b612fe5565b3561380f81610634565b90565b61381b90610dc0565b90565b5f80fd5b5f80fd5b5f80fd5b90359060016020038136030382121561386c570180359067ffffffffffffffff82116138675760200191600182023603831361386257565b613826565b613822565b61381e565b9082101561388c576020613888920281019061382a565b9091565b612fe5565b9035906001602003813603038212156138d3570180359067ffffffffffffffff82116138ce576020019160018202360383136138c957565b613826565b613822565b61381e565b908210156138f35760206138ef9202810190613891565b9091565b612fe5565b905090565b90918261390d81613914936138f8565b8093611f30565b0190565b9091613923926138fd565b90565b61393a613931610223565b92839283613918565b03902090565b90916139579260208301925f81850391015261371b565b90565b91602061397b92949361397460408201965f83019061031f565b019061031f565b565b35613987816104fc565b90565b9d9891949d9c939c9b929a9790959b99969960805261014052610100526139b16002610884565b6139cb6139c56139c05f61130a565b61031c565b9161031c565b141580613ffe575b613fe2576139e96139e38a613332565b15610282565b613fc6576139f88b8d90613414565b60c052613a068a8990613418565b613a1b613a1560c05192610631565b91610631565b03613faa5788888b8e8e95613a306008613430565b613a3990610dcc565b958b8b608051958c978c999b96916101405192610100519495969798613a5d610223565b9e8f9d8e9d8e613a706364c062d1611e3c565b81526004019d613a7f9e61373e565b03815a93602094fa8015613fa557613a9f915f91613f77575b5015610282565b613f5b57613ab7613ab26003899061082e565b612c97565b91613ad8613ad36005613acc60028701610860565b9501610a02565b612995565b92613af6613af0613aea878990611df6565b92610631565b91610631565b03613f3f5784613b17613b11613b0c600c610860565b610631565b91610631565b10613f23575f61012052613b2961206c565b610120525f60a052613b3961206c565b60a052613b4d6101405161010051906137ec565b60e0525b61012051613b69613b6360e051610631565b91610631565b1015613bb257613b9a613b926020613b8c610140516101005161012051916137f0565b01613805565b60a051611df6565b60a052613ba96101205161247b565b61012052613b51565b613bc8613bc0858790611df6565b60a051611df6565b60a05260a051613bfc613bf6613bf1613be030613812565b31613beb600b610860565b90611b63565b610631565b91610631565b11613f0757613c0c88869061495e565b613c155f611ade565b610120525b61012051613c32613c2c60c051610631565b91610631565b1015613ccd578a8a8a8a613ca0613c64613c568d6080519495906101205191613871565b9690959061012051916138d8565b919095613c9a613c947f17e8258c8124324c2b7a7f6e8d90caec0a1373c6eeffcbb90f3c3220485d43719561227a565b95610822565b95613926565b94613cb5613cac610223565b92839283613940565b0390a4613cc46101205161247b565b61012052613c1a565b92955092959850929650613d6a9550613ce7826002611ff3565b608051889192613d20613d1a7f8d0564932fa125f44360d4a175679e6d6e5600f3d4722c3ab2a877427e7f998e9361227a565b93610822565b93613d35613d2c610223565b9283928361395a565b0390a381613d4b613d455f611ade565b91610631565b11613e9c575b5050613d65613d60600d6129e5565b611abc565b6148c6565b613d735f611ade565b610120525b61012051613d9e613d98613d936101405161010051906137ec565b610631565b91610631565b1015613e9857613df4613dc65f613dc0610140516101005161012051916137f0565b0161397d565b613def613de96020613de3610140516101005161012051916137f0565b01613805565b91611abc565b6148c6565b60805182613e175f613e11610140516101005161012051916137f0565b0161397d565b91613e386020613e32610140516101005161012051916137f0565b01613805565b613e6b613e657fccf5242d67663517e020b47096ecb69c7f1a32e85ebd295eae0611a41395b3ae9361227a565b93610822565b93613e80613e77610223565b928392836129ae565b0390a3613e8f6101205161247b565b61012052613d78565b9050565b613eaf81613eaa8491611abc565b6148c6565b608051869192613ee8613ee27f06535f0bcdfcf8a3a06263cee80e795dfdcee6811c1d1fb959937464a1fd0b7f9361227a565b93610822565b93613efd613ef4610223565b928392836129ae565b0390a35f80613d51565b5f631e9acf1760e31b815280613f1f600482016104b7565b0390fd5b5f632a9ffab760e21b815280613f3b600482016104b7565b0390fd5b5f632a9ffab760e21b815280613f57600482016104b7565b0390fd5b5f638baa579f60e01b815280613f73600482016104b7565b0390fd5b613f98915060203d8111613f9e575b613f908183610998565b810190611e65565b5f613a98565b503d613f86565b611ecf565b5f637c6953f960e01b815280613fc2600482016104b7565b0390fd5b5f6302e8145360e61b815280613fde600482016104b7565b0390fd5b5f630b6fac0360e41b815280613ffa600482016104b7565b0390fd5b508361401b6140156140106002610884565b61031c565b9161031c565b14156139d3565b906140389d9c9b9a999897969594939291613399565b565b919091614045612ec0565b5061404e612e5f565b928161406261405c86610631565b91610631565b10158015614171575b614159576140799082611df6565b928361408d61408783610631565b91610631565b1161414f575b506140cd6140bd6140ad6140a8868590611b63565b612fbe565b926140b86005610860565b611df6565b936140c86005610860565b611df6565b916140d661206c565b935b806140eb6140e586610631565b91610631565b10156141475761413b6141419161413461411961411261410d60048590612465565b610884565b600361082e565b866141248a9261297d565b61412e8383612ff9565b52612ff9565b515061247b565b9461247b565b936140d8565b509250905090565b909250915f614093565b5050905061416e6141695f611ade565b612fbe565b90565b508061418561417f5f611ade565b91610631565b1461406b565b614193612f16565b90565b61419e61418b565b506141a7611a30565b506141b06119ec565b506141b9612e5f565b6141cb6141c55f611ade565b91610631565b116141ea576141d861418b565b6141e26002610884565b915f91929190565b61421161420a61420560046141ff6005610860565b90612465565b610884565b600361082e565b61421b6002610884565b9161422760019261297d565b929190565b6142459061424061423b6103ec565b61439c565b6142b3565b565b61425b61425661426092611302565b610da1565b610450565b90565b61426c90614247565b90565b61427890612989565b90565b90565b9061429361428e61429a9261426f565b61427b565b8254612209565b9055565b91906142b1905f602085019401906129a1565b565b806142ce6142c86142c35f614263565b61045b565b91611abc565b14614319576142de81600d61427e565b6143147fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f9161430b610223565b9182918261429e565b0390a1565b5f632582a64160e11b815280614331600482016104b7565b0390fd5b61433e9061422c565b565b9061435b9161435661435182611a4a565b61439c565b61435d565b565b9061436791614572565b50565b9061437491614340565b565b61437e6119ec565b506143986143926301ffc9a760e01b610235565b91610235565b1490565b6143ae906143a8614565565b906149d0565b565b906143bc60ff91611305565b9181191691161790565b6143cf90610282565b90565b90565b906143ea6143e56143f1926143c6565b6143d2565b82546143b0565b9055565b6143fd6119ec565b5061441261440c828490613107565b15610282565b5f1461449a5761443960016144345f61442c818690611a34565b0185906130ca565b6143d5565b90614442614565565b9061447f6144796144737f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d95610822565b92611934565b92611934565b92614488610223565b80614492816104b7565b0390a4600190565b50505f90565b90565b6144b76144b26144bc926144a0565b610da1565b610631565b90565b6144c960026144a3565b90565b6144d66001610860565b6144ef6144e96144e46144bf565b610631565b91610631565b14614508576145066144ff6144bf565b6001611b2f565b565b5f633ee5aeb560e01b815280614520600482016104b7565b0390fd5b61453861453361453d926110d5565b610da1565b610631565b90565b61454a6001614524565b90565b61455f614558614540565b6001611b2f565b565b5f90565b61456d614561565b503390565b61457a6119ec565b50614586818390613107565b5f1461460d576145ac5f6145a75f61459f818690611a34565b0185906130ca565b6143d5565b906145b5614565565b906145f26145ec6145e67ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b95610822565b92611934565b92611934565b926145fb610223565b80614605816104b7565b0390a4600190565b50505f90565b919061462961462461463193610822565b611fe7565b908354612025565b9055565b61464791614641611a30565b91614613565b565b9061465c905f1990602003600802610d5d565b8154169055565b905f9161467a614672826108d8565b92835461210f565b905555565b919290602082105f146146d857601f84116001146146a8576146a292935061210f565b90555b5b565b50906146ce6146d39360016146c56146bf856108d8565b92612017565b82019101612084565b614663565b6146a5565b5061470f82936146e96001946108d8565b6147086146f585612017565b820192601f86168061471a575b50612017565b0190612084565b6002021790556146a6565b61472690888603614649565b5f614702565b92909168010000000000000000821161478c576020115f1461477d57602081105f146147615761475b9161210f565b90555b5b565b60019160ff1916614771846108d8565b5560020201905561475e565b6001915060020201905561475f565b610984565b90815461479d816108a5565b908183116147c6575b8183106147b4575b50505050565b6147bd9361467f565b5f8080806147ae565b6147d28383838761472c565b6147a6565b5f6147e191614791565b565b905f036147f5576147f3906147d7565b565b611fba565b60055f9161480a83808301612070565b6148178360018301612070565b6148248360028301612070565b6148318360038301614635565b61483e83600483016147e3565b0155565b905f0361485457614852906147fa565b565b611fba565b61488a5f614885600361487f61487a60046148746005610860565b90612465565b610884565b9061082e565b614842565b6148a85f6148a3600461489d6005610860565b90612465565b614635565b6148c46148bd6148b86005610860565b61247b565b6005611b2f565b565b61490b916148f5614904926148ef6148e08492600a611940565b916148ea83610860565b611df6565b90611b2f565b6148ff600b610860565b611df6565b600b611b2f565b565b6149185f8092612a5a565b0190565b909161495b9361494461494e9261493a60808601965f870190610b16565b6020850190612a19565b6040830190612a4d565b606081830391015261490d565b90565b614966614859565b5f5f926149a86149967f605844df57be35fa1a2502fdca02ef466e14bc3f85566fccce96513e796e198894610822565b9461499f610223565b9384938461491c565b0390a2565b9160206149ce9294936149c760408201965f830190610b6c565b019061031f565b565b906149e56149df838390613107565b15610282565b6149ed575050565b614a075f92839263e2517d3f60e01b8452600484016149ad565b0390fdfea2646970667358221220a485bc3eea3f2ab619ba364b84e352b6d7c6681ace91a30ce35e3d7f2bddc4b864736f6c634300081e0033",
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
// Solidity: constructor(address _teeAuthenticator, address _authorityRegistry, address updateStatusOperator, address admin, uint256 _minFeePerRequest) returns()
func (processorEndpoint *ProcessorEndpoint) PackConstructor(_teeAuthenticator common.Address, _authorityRegistry common.Address, updateStatusOperator common.Address, admin common.Address, _minFeePerRequest *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("", _teeAuthenticator, _authorityRegistry, updateStatusOperator, admin, _minFeePerRequest)
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

// PackAPPLICATIONID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9968da96.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function APPLICATION_ID() view returns(uint64)
func (processorEndpoint *ProcessorEndpoint) PackAPPLICATIONID() []byte {
	enc, err := processorEndpoint.abi.Pack("APPLICATION_ID")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAPPLICATIONID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9968da96.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function APPLICATION_ID() view returns(uint64)
func (processorEndpoint *ProcessorEndpoint) TryPackAPPLICATIONID() ([]byte, error) {
	return processorEndpoint.abi.Pack("APPLICATION_ID")
}

// UnpackAPPLICATIONID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9968da96.
//
// Solidity: function APPLICATION_ID() view returns(uint64)
func (processorEndpoint *ProcessorEndpoint) UnpackAPPLICATIONID(data []byte) (uint64, error) {
	out, err := processorEndpoint.abi.Unpack("APPLICATION_ID", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
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
// the contract method with ID 0x9a17995a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, idx *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, depositAmount, idx)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenerateRequestId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a17995a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackGenerateRequestId(sender common.Address, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, idx *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("generateRequestId", sender, applicationId, requestType, payload, depositAmount, idx)
}

// UnpackGenerateRequestId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a17995a.
//
// Solidity: function generateRequestId(address sender, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 idx) pure returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackGenerateRequestId(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("generateRequestId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetNextPendingRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc404c359.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool success)
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
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool success)
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
// Solidity: function getNextPendingRequest() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8), bytes32, bool success)
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
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
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
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequests() ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequests")
}

// UnpackGetPendingRequests is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80a1f712.
//
// Solidity: function getPendingRequests() view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
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
// Solidity: function getPendingRequestsPage(uint256 offset, uint256 limit) view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
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
// Solidity: function getPendingRequestsPage(uint256 offset, uint256 limit) view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
func (processorEndpoint *ProcessorEndpoint) TryPackGetPendingRequestsPage(offset *big.Int, limit *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("getPendingRequestsPage", offset, limit)
}

// UnpackGetPendingRequestsPage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbe1e3725.
//
// Solidity: function getPendingRequestsPage(uint256 offset, uint256 limit) view returns((uint256,uint256,uint256,bytes32,bytes,address,uint64,uint8,uint8)[])
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

// PackMarkRequestCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d1be609.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function markRequestCompleted(bytes32 requestId, uint256 refund, uint256 applicationFees) returns()
func (processorEndpoint *ProcessorEndpoint) PackMarkRequestCompleted(requestId [32]byte, refund *big.Int, applicationFees *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("markRequestCompleted", requestId, refund, applicationFees)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMarkRequestCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d1be609.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function markRequestCompleted(bytes32 requestId, uint256 refund, uint256 applicationFees) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackMarkRequestCompleted(requestId [32]byte, refund *big.Int, applicationFees *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("markRequestCompleted", requestId, refund, applicationFees)
}

// PackMarkRequestFailed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b473cb9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function markRequestFailed(bytes32 requestId, uint8 errorCode, string errorMessage) returns()
func (processorEndpoint *ProcessorEndpoint) PackMarkRequestFailed(requestId [32]byte, errorCode uint8, errorMessage string) []byte {
	enc, err := processorEndpoint.abi.Pack("markRequestFailed", requestId, errorCode, errorMessage)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMarkRequestFailed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b473cb9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function markRequestFailed(bytes32 requestId, uint8 errorCode, string errorMessage) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackMarkRequestFailed(requestId [32]byte, errorCode uint8, errorMessage string) ([]byte, error) {
	return processorEndpoint.abi.Pack("markRequestFailed", requestId, errorCode, errorMessage)
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

// PackPayments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2982c21.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function payments(address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) PackPayments(arg0 common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("payments", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPayments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2982c21.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function payments(address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) TryPackPayments(arg0 common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("payments", arg0)
}

// UnpackPayments is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe2982c21.
//
// Solidity: function payments(address ) view returns(uint256)
func (processorEndpoint *ProcessorEndpoint) UnpackPayments(data []byte) (*big.Int, error) {
	out, err := processorEndpoint.abi.Unpack("payments", data)
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
// Solidity: function requestById(bytes32 ) view returns(uint256 timestamp, uint256 depositAmount, uint256 maxFeeValue, bytes32 requestId, bytes payload, address sender, uint64 applicationId, uint8 protocolVersion, uint8 requestType)
func (processorEndpoint *ProcessorEndpoint) PackRequestById(arg0 [32]byte) []byte {
	enc, err := processorEndpoint.abi.Pack("requestById", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestById is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5423112f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestById(bytes32 ) view returns(uint256 timestamp, uint256 depositAmount, uint256 maxFeeValue, bytes32 requestId, bytes payload, address sender, uint64 applicationId, uint8 protocolVersion, uint8 requestType)
func (processorEndpoint *ProcessorEndpoint) TryPackRequestById(arg0 [32]byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("requestById", arg0)
}

// RequestByIdOutput serves as a container for the return parameters of contract
// method RequestById.
type RequestByIdOutput struct {
	Timestamp       *big.Int
	DepositAmount   *big.Int
	MaxFeeValue     *big.Int
	RequestId       [32]byte
	Payload         []byte
	Sender          common.Address
	ApplicationId   uint64
	ProtocolVersion uint8
	RequestType     uint8
}

// UnpackRequestById is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5423112f.
//
// Solidity: function requestById(bytes32 ) view returns(uint256 timestamp, uint256 depositAmount, uint256 maxFeeValue, bytes32 requestId, bytes payload, address sender, uint64 applicationId, uint8 protocolVersion, uint8 requestType)
func (processorEndpoint *ProcessorEndpoint) UnpackRequestById(data []byte) (RequestByIdOutput, error) {
	out, err := processorEndpoint.abi.Unpack("requestById", data)
	outstruct := new(RequestByIdOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Timestamp = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.DepositAmount = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.MaxFeeValue = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.RequestId = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Payload = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.Sender = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.ApplicationId = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.ProtocolVersion = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.RequestType = *abi.ConvertType(out[8], new(uint8)).(*uint8)
	return *outstruct, nil
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

// PackStateRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9588eca2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateRoot() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackStateRoot() []byte {
	enc, err := processorEndpoint.abi.Pack("stateRoot")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9588eca2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateRoot() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackStateRoot() ([]byte, error) {
	return processorEndpoint.abi.Pack("stateRoot")
}

// UnpackStateRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackStateRoot(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("stateRoot", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ef290e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) PackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) []byte {
	enc, err := processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refund, applicationFees, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ef290e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateUpdate(uint64 applicationId, bytes32 prevStateRoot, bytes32 newStateRoot, bytes32 processedRequestId, bytes[] events, string[] eventSubTypes, (address,uint256)[] withdrawalRequests, uint256 refund, uint256 applicationFees, bytes signature) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackStateUpdate(applicationId uint64, prevStateRoot [32]byte, newStateRoot [32]byte, processedRequestId [32]byte, events [][]byte, eventSubTypes []string, withdrawalRequests []StructsWithdrawalRequest, refund *big.Int, applicationFees *big.Int, signature []byte) ([]byte, error) {
	return processorEndpoint.abi.Pack("stateUpdate", applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refund, applicationFees, signature)
}

// PackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34f23a9d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) PackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, maxFeeValue *big.Int) []byte {
	enc, err := processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, depositAmount, maxFeeValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSubmitRequest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34f23a9d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) TryPackSubmitRequest(protocolVersion uint8, applicationId uint64, requestType uint8, payload []byte, depositAmount *big.Int, maxFeeValue *big.Int) ([]byte, error) {
	return processorEndpoint.abi.Pack("submitRequest", protocolVersion, applicationId, requestType, payload, depositAmount, maxFeeValue)
}

// UnpackSubmitRequest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x34f23a9d.
//
// Solidity: function submitRequest(uint8 protocolVersion, uint64 applicationId, uint8 requestType, bytes payload, uint256 depositAmount, uint256 maxFeeValue) payable returns(bytes32)
func (processorEndpoint *ProcessorEndpoint) UnpackSubmitRequest(data []byte) ([32]byte, error) {
	out, err := processorEndpoint.abi.Unpack("submitRequest", data)
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

// PackWithdrawPayments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31b3eb94.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawPayments(address payee) returns()
func (processorEndpoint *ProcessorEndpoint) PackWithdrawPayments(payee common.Address) []byte {
	enc, err := processorEndpoint.abi.Pack("withdrawPayments", payee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawPayments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31b3eb94.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawPayments(address payee) returns()
func (processorEndpoint *ProcessorEndpoint) TryPackWithdrawPayments(payee common.Address) ([]byte, error) {
	return processorEndpoint.abi.Pack("withdrawPayments", payee)
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

// ProcessorEndpointPaymentWithdrawn represents a PaymentWithdrawn event raised by the ProcessorEndpoint contract.
type ProcessorEndpointPaymentWithdrawn struct {
	Payee  common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointPaymentWithdrawnEventName = "PaymentWithdrawn"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointPaymentWithdrawn) ContractEventName() string {
	return ProcessorEndpointPaymentWithdrawnEventName
}

// UnpackPaymentWithdrawnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentWithdrawn(address indexed payee, uint256 amount)
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
// Solidity: event Refund(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount)
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

// ProcessorEndpointRequestCompleted represents a RequestCompleted event raised by the ProcessorEndpoint contract.
type ProcessorEndpointRequestCompleted struct {
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
// Solidity: event RequestCompleted(bytes32 indexed requestId, uint256 applicationFees, uint8 status, uint8 errorCode, string errorMessage)
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
	RequestId [32]byte
	Sender    common.Address
	Raw       *types.Log // Blockchain specific contextual infos
}

const ProcessorEndpointRequestSubmittedEventName = "RequestSubmitted"

// ContractEventName returns the user-defined event name.
func (ProcessorEndpointRequestSubmitted) ContractEventName() string {
	return ProcessorEndpointRequestSubmittedEventName
}

// UnpackRequestSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RequestSubmitted(bytes32 indexed requestId, address indexed sender)
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

// ProcessorEndpointUserEvent represents a UserEvent event raised by the ProcessorEndpoint contract.
type ProcessorEndpointUserEvent struct {
	ApplicationId uint64
	RequestId     [32]byte
	EventSubType  common.Hash
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
// Solidity: event UserEvent(uint64 indexed applicationId, bytes32 indexed requestId, string indexed eventSubType, bytes encryptedData)
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
// Solidity: event Withdrawal(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount)
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["FeeValueBelowMinimum"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackFeeValueBelowMinimumError(raw[4:])
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
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidProtocolVersion"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidProtocolVersionError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidRequestId"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidRequestIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidStateRoot"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidStateRootError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["InvalidValue"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackInvalidValueError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["QueueThresholdExceeded"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackQueueThresholdExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], processorEndpoint.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return processorEndpoint.UnpackReentrancyGuardReentrantCallError(raw[4:])
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
