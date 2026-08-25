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

	// CodeGuestExecutionTimeout marks a guest that exceeded the configured
	// wall-clock execution bound and was interrupted (see pkg/wasm, epoch
	// interruption). It is a REQUEST FAILURE: signed and submitted on-chain.
	//
	// That is deliberate despite the non-determinism of a wall-clock bound — the
	// same request could in principle time out on a loaded host and succeed on an
	// idle one. The on-chain request queue is FIFO and its head only advances when
	// an update is submitted for it, so a transient error re-delivers the very same
	// request on the next poll, forever. Treating a runaway guest as transient would
	// therefore let a single non-terminating app pin the queue head and stall every
	// other application indefinitely — worse than the hang this bound removes. The
	// default timeout is set orders of magnitude above a real request precisely so
	// that reaching it means the guest is pathological.
	CodeGuestExecutionTimeout = FailureCode{"GUEST_EXECUTION_TIMEOUT", CategoryWasmInternalMeta}

	// CodeGuestExecutionCancelled marks a guest interrupted because the HOST gave up
	// on the request — context cancellation, executor shutdown, or a caller-supplied
	// execution budget running out. It is TRANSIENT: pkg/executor converts it to a
	// plain Go error so the request is retried on the next poll and nothing is signed.
	//
	// The distinction from CodeGuestExecutionTimeout is about who is at fault. Here
	// nothing is known to be wrong with the request, so charging the user
	// MinFeePerRequest for our own shutdown would be indefensible; and the condition
	// clears by itself once the executor is running again. It matters most for the
	// caller-supplied budget, which arrives from outside the enclave: routing it here
	// means no host-supplied value can ever produce a signed, fee-charging failure.
	//
	// Keep this in sync with IsTransient below.
	CodeGuestExecutionCancelled = FailureCode{"GUEST_EXECUTION_CANCELLED", CategoryWasmInternalMeta}
)

// Sentinels for the runtime entry points that return a plain error rather than a
// *RequestFailure (Deploy and the module-load path), so those callers can classify
// an interrupted guest with errors.Is. They carry the same meaning as the codes of
// the same name.
// Their messages double as the signed on-chain ErrorMsg for these failures, so
// they are deliberately short, fixed, and free of any wasmtime formatting or code
// offsets: ErrorMsg is truncated to 100 characters and changing one of these
// changes a signed value.
var (
	ErrGuestExecutionTimeout   = errors.New("guest execution timed out")
	ErrGuestExecutionCancelled = errors.New("guest execution cancelled")
)

// IsTransient reports whether a failure must NOT be turned into a signed on-chain
// error payload, and should instead surface as a plain Go error so the manager
// retries the request on its next poll.
//
// This is the executor's two-channel rule (see pkg/executor): almost every
// RequestFailure is a genuine execution failure that the user should learn about
// on-chain and be refunded for. The exception is a failure caused by the host
// abandoning the work rather than by anything about the request.
func IsTransient(failure *RequestFailure) bool {
	return failure != nil && failure.RequestError == CodeGuestExecutionCancelled
}

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
