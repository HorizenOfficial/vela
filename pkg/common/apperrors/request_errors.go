package apperrors

import (
	"errors"
	"fmt"
)

// NOTE: Keep category IDs in sync with ErrorCode enum in contracts/contracts/Structs.sol.
type category uint8

const (
	categoryNoError category = iota
	categoryUnknown
	categoryInternal
	categoryAppNotAdmitted
	categoryApplicationAlreadyDeployed
	categoryFailureWhenDeployingApplication
	categoryDeanonymizationReportFailed
	categoryRequestTypeNotPermitted
	categoryFunctionNotFound
	categoryDepositFailed
	categoryRequestFuncFailed
	categoryAppNotDeployed
	categoryWrongKeySent
	categoryPubKeyNotRegistered
	categoryNoReportDataFound
)

type ErrorCategory struct {
	Category      category
	Message string
}

var (
	CategoryNoErrorMeta                         = ErrorCategory{Category: categoryNoError, Message: "no error"}
	CategoryUnknownMeta                         = ErrorCategory{Category: categoryUnknown, Message: "unknown error"}
	CategoryInternalMeta                        = ErrorCategory{Category: categoryInternal, Message: "internal error"}
	CategoryAppNotAdmittedMeta                  = ErrorCategory{Category: categoryAppNotAdmitted, Message: "application is not admitted"}
	CategoryApplicationAlreadyDeployedMeta      = ErrorCategory{Category: categoryApplicationAlreadyDeployed, Message: "application is already deployed"}
	CategoryFailureWhenDeployingApplicationMeta = ErrorCategory{Category: categoryFailureWhenDeployingApplication, Message: "failed to deploy application"}
	CategoryDeanonymizationReportFailedMeta     = ErrorCategory{Category: categoryDeanonymizationReportFailed, Message: "failed to generate deanonymization report"}
	CategoryRequestTypeNotPermittedMeta         = ErrorCategory{Category: categoryRequestTypeNotPermitted, Message: "request type is not permitted"}
	CategoryFunctionNotFoundMeta                = ErrorCategory{Category: categoryFunctionNotFound, Message: "requested function not found"}
	CategoryDepositFailedMeta                   = ErrorCategory{Category: categoryDepositFailed, Message: "deposit operation failed"}
	CategoryRequestFuncFailedMeta               = ErrorCategory{Category: categoryRequestFuncFailed, Message: "request function execution failed"}
	CategoryAppNotDeployedMeta                  = ErrorCategory{Category: categoryAppNotDeployed, Message: "application is not deployed"}
	CategoryWrongKeySentMeta                    = ErrorCategory{Category: categoryWrongKeySent, Message: "wrong key sent"}
	CategoryPubKeyNotRegisteredMeta             = ErrorCategory{Category: categoryPubKeyNotRegistered, Message: "public key not registered"}
	CategoryNoReportDataFoundMeta               = ErrorCategory{Category: categoryNoReportDataFound, Message: "no report data found"}
)

type ErrorCode struct {
	Code     string
	Category ErrorCategory
}

func (c ErrorCode) String() string {
	return c.Code
}

var (
	// Internal errors with more details
	CodeInternalFallback             = ErrorCode{"INTERNAL_FALLBACK", CategoryInternalMeta}
	CodeManagerNotRunning            = ErrorCode{"MANAGER_NOT_RUNNING", CategoryInternalMeta}
	CodeAppNotAdmitted               = ErrorCode{"APP_NOT_ADMITTED", CategoryAppNotAdmittedMeta}
	CodeApplicationAlreadyDeployed   = ErrorCode{"APPLICATION_ALREADY_DEPLOYED", CategoryApplicationAlreadyDeployedMeta}
	CodeWasmNotFound                 = ErrorCode{"WASM_NOT_FOUND", CategoryFailureWhenDeployingApplicationMeta}
	CodeErrorReadingWasm             = ErrorCode{"ERROR_READING_WASM", CategoryFailureWhenDeployingApplicationMeta}
	CodeFailureDeployingApp          = ErrorCode{"FAILURE_WHEN_DEPLOYING_APPLICATION", CategoryFailureWhenDeployingApplicationMeta}
	CodeSubmittingStateUpdateFailed  = ErrorCode{"SUBMITTING_STATE_UPDATE_FAILED", CategoryInternalMeta}
	CodeAppStateNotFound             = ErrorCode{"APPLICATION_STATE_NOT_FOUND", CategoryAppNotDeployedMeta}
	CodeRequestTypeNotPermitted      = ErrorCode{"REQUEST_TYPE_NOT_PERMITTED", CategoryRequestTypeNotPermittedMeta}
	CodeEncryptedToAppDataFailure    = ErrorCode{"ENCRYPTED_TO_APPDATA_FAILURE", CategoryInternalMeta}
	CodeJsonUnmarshalError           = ErrorCode{"JSON_UNMARSHAL_ERROR", CategoryInternalMeta}
	CodeJsonMarshalError             = ErrorCode{"JSON_MARSHAL_ERROR", CategoryInternalMeta}
	CodeParsingKeyError              = ErrorCode{"PARSING_KEY_ERROR", CategoryWrongKeySentMeta}
	CodeAppDataSerializationFailure  = ErrorCode{"APPDATA_SERIALIZATION_FAILURE", CategoryInternalMeta}
	CodeAppDataEncryptionFailure     = ErrorCode{"APPDATA_ENCRYPTION_FAILURE", CategoryInternalMeta}
	CodePayloadUpdateSigningFailure  = ErrorCode{"PAYLOAD_UPDATE_SIGNING_FAILURE", CategoryInternalMeta}
	CodeFunctionNotFound             = ErrorCode{"FUNCTION_NOT_FOUND", CategoryFunctionNotFoundMeta}
	CodeMemoryWriteError             = ErrorCode{"MEMORY_WRITE_ERROR", CategoryInternalMeta}
	CodeFailedExtractingResultBytes  = ErrorCode{"FAILED_EXTRACTING_RESULT_BYTES", CategoryInternalMeta}
	CodeFailedLoadingOrGettingModule = ErrorCode{"FAILED_LOADING_OR_GETTING_MODULE", CategoryInternalMeta}
	CodeDepositFailed                = ErrorCode{"DEPOSIT_FAILED", CategoryDepositFailedMeta}
	CodeRequestFuncFailed            = ErrorCode{"REQUEST_FUNC_FAILED", CategoryRequestFuncFailedMeta}
	CodeFailedToGenerateReport       = ErrorCode{"FAILED_TO_GENERATE_REPORT", CategoryDeanonymizationReportFailedMeta}
	CodePubKeyNotRegistered          = ErrorCode{"PUBKEY_NOT_REGISTERED", CategoryPubKeyNotRegisteredMeta}
	CodeWrongKey                     = ErrorCode{"WRONG_KEY", CategoryWrongKeySentMeta}
	CodeNoReportDataFound            = ErrorCode{"NO_REPORT_DATA_FOUND", CategoryNoReportDataFoundMeta}
)

type RequestFailure struct {
	RequestError ErrorCode
	Message string
	RootErr     error
}

// RemoteError captures error information transmitted over the wire.
type RemoteError struct {
	Type    string
	Message string
}

func (e *RemoteError) Error() string {
	if e.Type == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *RemoteError) Is(target error) bool {
	other, ok := target.(*RemoteError)
	if !ok {
		return false
	}
	return e.Type == other.Type
}

// RequestFailureDTO is a JSON-serializable representation of RequestFailure.
type RequestFailureDTO struct {
	RequestError ErrorCode `json:"requestError"`
	Message      string    `json:"message"`
	ErrMessage   string    `json:"errMessage,omitempty"`
	ErrType      string    `json:"errType,omitempty"`
}

func New(code ErrorCode, msg string, cause error) *RequestFailure {
	return &RequestFailure{RequestError: code, Message: msg, RootErr: cause}
}

// ToDTO converts a RequestFailure into its transport-friendly representation.
func (e *RequestFailure) ToDTO() *RequestFailureDTO {
	if e == nil {
		return nil
	}

	dto := &RequestFailureDTO{
		RequestError: e.RequestError,
		Message:      e.Message,
	}

	if e.RootErr != nil {
		dto.ErrMessage = e.RootErr.Error()
		dto.ErrType = fmt.Sprintf("%T", e.RootErr)
	}

	return dto
}

// ToFailure converts a DTO back into a RequestFailure.
func (dto *RequestFailureDTO) ToFailure() *RequestFailure {
	if dto == nil {
		return nil
	}

	failure := &RequestFailure{
		RequestError: dto.RequestError,
		Message: dto.Message,
	}

	if dto.ErrMessage != "" {
		failure.RootErr = &RemoteError{
			Type:    dto.ErrType,
			Message: dto.ErrMessage,
		}
	}

	return failure
}

func (e *RequestFailure) Error() string {
	if e.RootErr != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.RootErr)
	}
	return e.Message
}

func (e *RequestFailure) Unwrap() error { return e.RootErr }

func (e *RequestFailure) Category() uint8 {
	return uint8(e.categoryMeta().Category)
}

func (e *RequestFailure) ExternalMessage() string {
	return e.categoryMeta().Message
}

func (e *RequestFailure) categoryMeta() ErrorCategory {
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