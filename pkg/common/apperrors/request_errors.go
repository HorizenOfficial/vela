package apperrors

import "errors"

type category uint8

const (
	categoryNoError category = iota
	categoryUnknown
	categoryInternal
	categoryApplicationAlreadyDeployed
	categoryFunctionNotFound
	categoryDepositFailed
	categoryRequestFuncFailed
	categoryAppNotDeployed
	categoryWrongKeySent
	categoryPubKeyNotRegistered
	categoryNoReportDataFound
	categoryWasmInternal
	categoryInsufficientFuel
	// NOTE: Keep category IDs in sync with ErrorCode enum in contracts/contracts/Structs.sol.
)

type errorCategory struct {
	Category category
	Message  string
}

var (
	CategoryNoErrorMeta                    = errorCategory{Category: categoryNoError, Message: "no error"}
	CategoryUnknownMeta                    = errorCategory{Category: categoryUnknown, Message: "unknown error"}
	CategoryInternalMeta                   = errorCategory{Category: categoryInternal, Message: "internal error"}
	CategoryApplicationAlreadyDeployedMeta = errorCategory{Category: categoryApplicationAlreadyDeployed, Message: "application is already deployed"}
	CategoryFunctionNotFoundMeta           = errorCategory{Category: categoryFunctionNotFound, Message: "requested function not found"}
	CategoryDepositFailedMeta              = errorCategory{Category: categoryDepositFailed, Message: "deposit operation failed"}
	CategoryRequestFuncFailedMeta          = errorCategory{Category: categoryRequestFuncFailed, Message: "request function execution failed"}
	CategoryAppNotDeployedMeta             = errorCategory{Category: categoryAppNotDeployed, Message: "application is not deployed"}
	CategoryWrongKeySentMeta               = errorCategory{Category: categoryWrongKeySent, Message: "wrong key sent"}
	CategoryPubKeyNotRegisteredMeta        = errorCategory{Category: categoryPubKeyNotRegistered, Message: "public key not registered"}
	CategoryNoReportDataFoundMeta          = errorCategory{Category: categoryNoReportDataFound, Message: "no report data found"}
	CategoryWasmInternalMeta               = errorCategory{Category: categoryWasmInternal, Message: "wasm internal error"}
	CategoryInsufficientFuelMeta           = errorCategory{Category: categoryInsufficientFuel, Message: "insufficient fuel"}
)

type FailureCode struct {
	Code     string
	Category errorCategory
}

func (c FailureCode) String() string {
	return c.Code
}

var (
	// Internal errors with more details
	CodeInternalFallback             = FailureCode{"INTERNAL_FALLBACK", CategoryInternalMeta}
	CodeApplicationAlreadyDeployed   = FailureCode{"APPLICATION_ALREADY_DEPLOYED", CategoryApplicationAlreadyDeployedMeta}
	CodeAppStateNotFound             = FailureCode{"APPLICATION_STATE_NOT_FOUND", CategoryAppNotDeployedMeta}
	CodeJsonUnmarshalError           = FailureCode{"JSON_UNMARSHAL_ERROR", CategoryInternalMeta}
	CodeJsonMarshalError             = FailureCode{"JSON_MARSHAL_ERROR", CategoryInternalMeta}
	CodeParsingKeyError              = FailureCode{"PARSING_KEY_ERROR", CategoryWrongKeySentMeta}
	CodeAppDataSerializationFailure  = FailureCode{"APPDATA_SERIALIZATION_FAILURE", CategoryInternalMeta}
	CodeFunctionNotFound             = FailureCode{"FUNCTION_NOT_FOUND", CategoryFunctionNotFoundMeta}
	CodeMemoryWriteError             = FailureCode{"MEMORY_WRITE_ERROR", CategoryWasmInternalMeta}
	CodeFailedExtractingResultBytes  = FailureCode{"FAILED_EXTRACTING_RESULT_BYTES", CategoryWasmInternalMeta}
	CodeFailedLoadingOrGettingModule = FailureCode{"FAILED_LOADING_OR_GETTING_MODULE", CategoryWasmInternalMeta}
	CodeWasmModuleEmpty              = FailureCode{"WASM_MODULE_EMPTY", CategoryWasmInternalMeta}
	CodeDepositFailed                = FailureCode{"DEPOSIT_FAILED", CategoryDepositFailedMeta}
	CodeRequestFuncFailed            = FailureCode{"REQUEST_FUNC_FAILED", CategoryRequestFuncFailedMeta}
	CodePubKeyNotRegistered          = FailureCode{"PUBKEY_NOT_REGISTERED", CategoryPubKeyNotRegisteredMeta}
	CodeWrongKey                     = FailureCode{"WRONG_KEY", CategoryWrongKeySentMeta}
	CodeNoReportDataFound            = FailureCode{"NO_REPORT_DATA_FOUND", CategoryNoReportDataFoundMeta}
	CodeInsufficientFuel             = FailureCode{"INSUFFICIENT_FUEL", CategoryInsufficientFuelMeta}
)

type RequestFailure struct {
	RequestError FailureCode
	Message      string
}

func New(code FailureCode, msg string) *RequestFailure {
	return &RequestFailure{RequestError: code, Message: msg}
}

func (e *RequestFailure) Error() string {
	return e.Message
}

func (e *RequestFailure) Category() uint8 {
	return uint8(e.categoryMeta().Category)
}

func (e *RequestFailure) ExternalMessage() string {
	if e.Message != "" {
		return e.Message
	}
	return e.categoryMeta().Message
}

func (e *RequestFailure) categoryMeta() errorCategory {
	if e == nil {
		return CategoryNoErrorMeta
	}

	if e.RequestError.Category.Message == "" {
		return CategoryUnknownMeta
	}
	return e.RequestError.Category
}

func FromError(err error) (*RequestFailure, bool) {
	var failure *RequestFailure
	if errors.As(err, &failure) {
		return failure, true
	}
	return nil, false
}
