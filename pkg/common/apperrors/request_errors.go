package apperrors

import (
	"errors"
	"fmt"
)

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
	categoryWasmInternal
	// NOTE: Keep category IDs in sync with ErrorCode enum in contracts/contracts/Structs.sol.
)

type errorCategory struct {
	Category     category
	Message      string
}

var (
	CategoryNoErrorMeta                         = errorCategory{Category: categoryNoError, Message: "no error"}
	CategoryUnknownMeta                         = errorCategory{Category: categoryUnknown, Message: "unknown error"}
	CategoryInternalMeta                        = errorCategory{Category: categoryInternal, Message: "internal error"}
	CategoryAppNotAdmittedMeta                  = errorCategory{Category: categoryAppNotAdmitted, Message: "application is not admitted"}
	CategoryApplicationAlreadyDeployedMeta      = errorCategory{Category: categoryApplicationAlreadyDeployed, Message: "application is already deployed"}
	CategoryFailureWhenDeployingApplicationMeta = errorCategory{Category: categoryFailureWhenDeployingApplication, Message: "failed to deploy application"}
	CategoryDeanonymizationReportFailedMeta     = errorCategory{Category: categoryDeanonymizationReportFailed, Message: "failed to generate deanonymization report"}
	CategoryRequestTypeNotPermittedMeta         = errorCategory{Category: categoryRequestTypeNotPermitted, Message: "request type is not permitted"}
	CategoryFunctionNotFoundMeta                = errorCategory{Category: categoryFunctionNotFound, Message: "requested function not found"}
	CategoryDepositFailedMeta                   = errorCategory{Category: categoryDepositFailed, Message: "deposit operation failed"}
	CategoryRequestFuncFailedMeta               = errorCategory{Category: categoryRequestFuncFailed, Message: "request function execution failed"}
	CategoryAppNotDeployedMeta                  = errorCategory{Category: categoryAppNotDeployed, Message: "application is not deployed"}
	CategoryWrongKeySentMeta                    = errorCategory{Category: categoryWrongKeySent, Message: "wrong key sent"}
	CategoryPubKeyNotRegisteredMeta             = errorCategory{Category: categoryPubKeyNotRegistered, Message: "public key not registered"}
	CategoryNoReportDataFoundMeta               = errorCategory{Category: categoryNoReportDataFound, Message: "no report data found"}
	CategoryWasmInternalMeta					= errorCategory{Category: categoryWasmInternal, Message: "wasm internal error"}
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
	CodeManagerNotRunning            = FailureCode{"MANAGER_NOT_RUNNING", CategoryInternalMeta}
	CodeAppNotAdmitted               = FailureCode{"APP_NOT_ADMITTED", CategoryAppNotAdmittedMeta}
	CodeApplicationAlreadyDeployed   = FailureCode{"APPLICATION_ALREADY_DEPLOYED", CategoryApplicationAlreadyDeployedMeta}
	CodeWasmNotFound                 = FailureCode{"WASM_NOT_FOUND", CategoryFailureWhenDeployingApplicationMeta}
	CodeErrorReadingWasm             = FailureCode{"ERROR_READING_WASM", CategoryFailureWhenDeployingApplicationMeta}
	CodeFailureDeployingApp          = FailureCode{"FAILURE_WHEN_DEPLOYING_APPLICATION", CategoryFailureWhenDeployingApplicationMeta}
	CodeSubmittingStateUpdateFailed  = FailureCode{"SUBMITTING_STATE_UPDATE_FAILED", CategoryInternalMeta}
	CodeAppStateNotFound             = FailureCode{"APPLICATION_STATE_NOT_FOUND", CategoryAppNotDeployedMeta}
	CodeRequestTypeNotPermitted      = FailureCode{"REQUEST_TYPE_NOT_PERMITTED", CategoryRequestTypeNotPermittedMeta}
	CodeEncryptedToAppDataFailure    = FailureCode{"ENCRYPTED_TO_APPDATA_FAILURE", CategoryInternalMeta}
	CodeJsonUnmarshalError           = FailureCode{"JSON_UNMARSHAL_ERROR", CategoryInternalMeta}
	CodeJsonMarshalError             = FailureCode{"JSON_MARSHAL_ERROR", CategoryInternalMeta}
	CodeParsingKeyError              = FailureCode{"PARSING_KEY_ERROR", CategoryWrongKeySentMeta}
	CodeAppDataSerializationFailure  = FailureCode{"APPDATA_SERIALIZATION_FAILURE", CategoryInternalMeta}
	CodeAppDataEncryptionFailure     = FailureCode{"APPDATA_ENCRYPTION_FAILURE", CategoryInternalMeta}
	CodePayloadUpdateSigningFailure  = FailureCode{"PAYLOAD_UPDATE_SIGNING_FAILURE", CategoryInternalMeta}
	CodeFunctionNotFound             = FailureCode{"FUNCTION_NOT_FOUND", CategoryFunctionNotFoundMeta}
	CodeMemoryWriteError             = FailureCode{"MEMORY_WRITE_ERROR", CategoryWasmInternalMeta}
	CodeFailedExtractingResultBytes  = FailureCode{"FAILED_EXTRACTING_RESULT_BYTES", CategoryWasmInternalMeta}
	CodeFailedLoadingOrGettingModule = FailureCode{"FAILED_LOADING_OR_GETTING_MODULE", CategoryWasmInternalMeta}
	CodeDepositFailed                = FailureCode{"DEPOSIT_FAILED", CategoryDepositFailedMeta}
	CodeRequestFuncFailed            = FailureCode{"REQUEST_FUNC_FAILED", CategoryRequestFuncFailedMeta}
	CodeFailedToGenerateReport       = FailureCode{"FAILED_TO_GENERATE_REPORT", CategoryDeanonymizationReportFailedMeta}
	CodePubKeyNotRegistered          = FailureCode{"PUBKEY_NOT_REGISTERED", CategoryPubKeyNotRegisteredMeta}
	CodeWrongKey                     = FailureCode{"WRONG_KEY", CategoryWrongKeySentMeta}
	CodeNoReportDataFound            = FailureCode{"NO_REPORT_DATA_FOUND", CategoryNoReportDataFoundMeta}
)

type RequestFailure struct {
	RequestError FailureCode
	Message      string
	RootErr      error
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
	RequestError FailureCode `json:"requestError"`
	Message      string      `json:"message"`
	ErrMessage   string      `json:"errMessage,omitempty"`
	ErrType      string      `json:"errType,omitempty"`
}

func New(code FailureCode, msg string, cause error) *RequestFailure {
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
		Message:      dto.Message,
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
