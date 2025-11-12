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

type errorCategory struct {
	Category      category
	Message string
}

var (
	categoryNoErrorMeta                         = errorCategory{Category: categoryNoError, Message: "no error"}
	categoryUnknownMeta                         = errorCategory{Category: categoryUnknown, Message: "unknown error"}
	categoryInternalMeta                        = errorCategory{Category: categoryInternal, Message: "internal error"}
	categoryAppNotAdmittedMeta                  = errorCategory{Category: categoryAppNotAdmitted, Message: "application is not admitted"}
	categoryApplicationAlreadyDeployedMeta      = errorCategory{Category: categoryApplicationAlreadyDeployed, Message: "application is already deployed"}
	categoryFailureWhenDeployingApplicationMeta = errorCategory{Category: categoryFailureWhenDeployingApplication, Message: "failed to deploy application"}
	categoryDeanonymizationReportFailedMeta     = errorCategory{Category: categoryDeanonymizationReportFailed, Message: "failed to generate deanonymization report"}
	categoryRequestTypeNotPermittedMeta         = errorCategory{Category: categoryRequestTypeNotPermitted, Message: "request type is not permitted"}
	categoryFunctionNotFoundMeta                = errorCategory{Category: categoryFunctionNotFound, Message: "requested function not found"}
	categoryDepositFailedMeta                   = errorCategory{Category: categoryDepositFailed, Message: "deposit operation failed"}
	categoryRequestFuncFailedMeta               = errorCategory{Category: categoryRequestFuncFailed, Message: "request function execution failed"}
	categoryAppNotDeployedMeta                  = errorCategory{Category: categoryAppNotDeployed, Message: "application is not deployed"}
	categoryWrongKeySentMeta                    = errorCategory{Category: categoryWrongKeySent, Message: "wrong key sent"}
	categoryPubKeyNotRegisteredMeta             = errorCategory{Category: categoryPubKeyNotRegistered, Message: "public key not registered"}
	categoryNoReportDataFoundMeta               = errorCategory{Category: categoryNoReportDataFound, Message: "no report data found"}
)

type ErrorCode struct {
	Code     string
	Category errorCategory
}

func (c ErrorCode) String() string {
	return c.Code
}

var (
	// Internal errors with more details
	CodeInternalFallback             = ErrorCode{"INTERNAL_FALLBACK", categoryInternalMeta}
	CodeManagerNotRunning            = ErrorCode{"MANAGER_NOT_RUNNING", categoryInternalMeta}
	CodeAppNotAdmitted               = ErrorCode{"APP_NOT_ADMITTED", categoryAppNotAdmittedMeta}
	CodeApplicationAlreadyDeployed   = ErrorCode{"APPLICATION_ALREADY_DEPLOYED", categoryApplicationAlreadyDeployedMeta}
	CodeWasmNotFound                 = ErrorCode{"WASM_NOT_FOUND", categoryFailureWhenDeployingApplicationMeta}
	CodeErrorReadingWasm             = ErrorCode{"ERROR_READING_WASM", categoryFailureWhenDeployingApplicationMeta}
	CodeFailureDeployingApp          = ErrorCode{"FAILURE_WHEN_DEPLOYING_APPLICATION", categoryFailureWhenDeployingApplicationMeta}
	CodeSubmittingStateUpdateFailed  = ErrorCode{"SUBMITTING_STATE_UPDATE_FAILED", categoryInternalMeta}
	CodeAppStateNotFound             = ErrorCode{"APPLICATION_STATE_NOT_FOUND", categoryAppNotDeployedMeta}
	CodeRequestTypeNotPermitted      = ErrorCode{"REQUEST_TYPE_NOT_PERMITTED", categoryRequestTypeNotPermittedMeta}
	CodeEncryptedToAppDataFailure    = ErrorCode{"ENCRYPTED_TO_APPDATA_FAILURE", categoryInternalMeta}
	CodeJsonUnmarshalError           = ErrorCode{"JSON_UNMARSHAL_ERROR", categoryInternalMeta}
	CodeJsonMarshalError             = ErrorCode{"JSON_MARSHAL_ERROR", categoryInternalMeta}
	CodeParsingKeyError              = ErrorCode{"PARSING_KEY_ERROR", categoryWrongKeySentMeta}
	CodeAppDataSerializationFailure  = ErrorCode{"APPDATA_SERIALIZATION_FAILURE", categoryInternalMeta}
	CodeAppDataEncryptionFailure     = ErrorCode{"APPDATA_ENCRYPTION_FAILURE", categoryInternalMeta}
	CodePayloadUpdateSigningFailure  = ErrorCode{"PAYLOAD_UPDATE_SIGNING_FAILURE", categoryInternalMeta}
	CodeFunctionNotFound             = ErrorCode{"FUNCTION_NOT_FOUND", categoryFunctionNotFoundMeta}
	CodeMemoryWriteError             = ErrorCode{"MEMORY_WRITE_ERROR", categoryInternalMeta}
	CodeFailedExtractingResultBytes  = ErrorCode{"FAILED_EXTRACTING_RESULT_BYTES", categoryInternalMeta}
	CodeFailedLoadingOrGettingModule = ErrorCode{"FAILED_LOADING_OR_GETTING_MODULE", categoryInternalMeta}
	CodeDepositFailed                = ErrorCode{"DEPOSIT_FAILED", categoryDepositFailedMeta}
	CodeRequestFuncFailed            = ErrorCode{"REQUEST_FUNC_FAILED", categoryRequestFuncFailedMeta}
	CodeFailedToGenerateReport       = ErrorCode{"FAILED_TO_GENERATE_REPORT", categoryDeanonymizationReportFailedMeta}
	CodePubKeyNotRegistered          = ErrorCode{"PUBKEY_NOT_REGISTERED", categoryPubKeyNotRegisteredMeta}
	CodeWrongKey                     = ErrorCode{"WRONG_KEY", categoryWrongKeySentMeta}
	CodeNoReportDataFound            = ErrorCode{"NO_REPORT_DATA_FOUND", categoryNoReportDataFoundMeta}
)

type RequestFailure struct {
	RootErr ErrorCode
	Message string
	Err     error
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
	return &RequestFailure{RootErr: code, Message: msg, Err: cause}
}

// ToDTO converts a RequestFailure into its transport-friendly representation.
func (e *RequestFailure) ToDTO() *RequestFailureDTO {
	if e == nil {
		return nil
	}

	dto := &RequestFailureDTO{
		RequestError: e.RootErr,
		Message:      e.Message,
	}

	if e.Err != nil {
		dto.ErrMessage = e.Err.Error()
		dto.ErrType = fmt.Sprintf("%T", e.Err)
	}

	return dto
}

// ToFailure converts a DTO back into a RequestFailure.
func (dto *RequestFailureDTO) ToFailure() *RequestFailure {
	if dto == nil {
		return nil
	}

	failure := &RequestFailure{
		RootErr: dto.RequestError,
		Message: dto.Message,
	}

	if dto.ErrMessage != "" {
		failure.Err = &RemoteError{
			Type:    dto.ErrType,
			Message: dto.ErrMessage,
		}
	}

	return failure
}

func (e *RequestFailure) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *RequestFailure) Unwrap() error { return e.Err }

func (e *RequestFailure) Category() uint8 {
	return uint8(e.categoryMeta().Category)
}

func (e *RequestFailure) ExternalMessage() string {
	return e.categoryMeta().Message
}

func (e *RequestFailure) categoryMeta() errorCategory {
	if e == nil {
		return categoryUnknownMeta
	}

	meta := e.RootErr.Category
	if meta.Message != "" || e.RootErr.Code == "" {
		return meta
	}
	return categoryUnknownMeta
}

func FromError(err error) (*RequestFailure, bool) {
	var failure *RequestFailure
	if errors.As(err, &failure) {
		return failure, true
	}
	return nil, false
}