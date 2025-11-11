package apperrors

import (
	"errors"
	"fmt"
)

type Category uint8

const (
	// External error categories for users and blockchain
	CategoryNoError Category = iota
	CategoryUnknown
	CategoryInternal
	CategoryAppNotAdmitted
	CategoryApplicationAlreadyDeployed
	CategoryFailureWhenDeployingApplication
	CategoryDeanonymizationReportFailed
	CategoryRequestTypeNotPermitted
	CategoryFunctionNotFound
	CategoryDepositFailed
	CategoryRequestFuncFailed
	CategoryAppNotDeployed
	CategoryWrongKeySent
)

type ErrorCategory struct {
	Category Category
	Message  string
}

var (
	CategoryNoErrorMeta                         = ErrorCategory{Category: CategoryNoError, Message: "no error"}
	CategoryUnknownMeta                         = ErrorCategory{Category: CategoryUnknown, Message: "unknown error"}
	CategoryInternalMeta                        = ErrorCategory{Category: CategoryInternal, Message: "internal error"}
	CategoryAppNotAdmittedMeta                  = ErrorCategory{Category: CategoryAppNotAdmitted, Message: "application is not admitted"}
	CategoryApplicationAlreadyDeployedMeta      = ErrorCategory{Category: CategoryApplicationAlreadyDeployed, Message: "application is already deployed"}
	CategoryFailureWhenDeployingApplicationMeta = ErrorCategory{Category: CategoryFailureWhenDeployingApplication, Message: "failed to deploy application"}
	CategoryDeanonymizationReportFailedMeta     = ErrorCategory{Category: CategoryDeanonymizationReportFailed, Message: "failed to generate deanonymization report"}
	CategoryRequestTypeNotPermittedMeta         = ErrorCategory{Category: CategoryRequestTypeNotPermitted, Message: "request type is not permitted"}
	CategoryFunctionNotFoundMeta                = ErrorCategory{Category: CategoryFunctionNotFound, Message: "requested function not found"}
	CategoryDepositFailedMeta                   = ErrorCategory{Category: CategoryDepositFailed, Message: "deposit operation failed"}
	CategoryRequestFuncFailedMeta               = ErrorCategory{Category: CategoryRequestFuncFailed, Message: "request function execution failed"}
	CategoryAppNotDeployedMeta                  = ErrorCategory{Category: CategoryAppNotDeployed, Message: "application is not deployed"}
	CategoryWrongKeySentMeta					= ErrorCategory{Category: CategoryWrongKeySent, Message: "wrong key sent"}
)

type RequestError struct {
	Code     string
	Category ErrorCategory
}

func (c RequestError) String() string {
	return c.Code
}

var (
	// Internal errors with more details
	CodeInternalFallback                       = RequestError{"INTERNAL_FALLBACK", CategoryInternalMeta}
	CodeManagerNotRunning                      = RequestError{"MANAGER_NOT_RUNNING", CategoryInternalMeta}
	CodeAppNotAdmitted                         = RequestError{"APP_NOT_ADMITTED", CategoryAppNotAdmittedMeta}
	CodeApplicationAlreadyDeployed             = RequestError{"APPLICATION_ALREADY_DEPLOYED", CategoryApplicationAlreadyDeployedMeta}
	CodeWasmNotFound                           = RequestError{"WASM_NOT_FOUND", CategoryFailureWhenDeployingApplicationMeta}
	CodeErrorReadingWasm                       = RequestError{"ERROR_READING_WASM", CategoryFailureWhenDeployingApplicationMeta}
	CodeFailureDeployingApp                    = RequestError{"FAILURE_WHEN_DEPLOYING_APPLICATION", CategoryFailureWhenDeployingApplicationMeta}
	CodeSubmittingStateUpdateFailed            = RequestError{"SUBMITTING_STATE_UPDATE_FAILED", CategoryInternalMeta}
	CodeAppStateNotFound                       = RequestError{"APPLICATION_STATE_NOT_FOUND", CategoryAppNotDeployedMeta}
	CodeRequestTypeNotPermitted                = RequestError{"REQUEST_TYPE_NOT_PERMITTED", CategoryRequestTypeNotPermittedMeta}
	CodeEncryptedToAppDataFailure              = RequestError{"ENCRYPTED_TO_APPDATA_FAILURE", CategoryInternalMeta}
	CodeJsonUnmarshalError                     = RequestError{"JSON_UNMARSHAL_ERROR", CategoryInternalMeta}
	CodeJsonMarshalError                       = RequestError{"JSON_MARSHAL_ERROR", CategoryInternalMeta}
	CodeParsingKeyError                        = RequestError{"PARSING_KEY_ERROR", CategoryWrongKeySentMeta}
	CodePayloadDecryptionFailure               = RequestError{"PAYLOAD_DECRYPTION_FAILURE", CategoryInternalMeta}
	CodeAppDataSerializationFailure            = RequestError{"APPDATA_SERIALIZATION_FAILURE", CategoryInternalMeta}
	CodeAppDataEncryptionFailure               = RequestError{"APPDATA_ENCRYPTION_FAILURE", CategoryInternalMeta}
	CodeEventsEncryptionFailure                = RequestError{"EVENTS_ENCRYPTION_FAILURE", CategoryInternalMeta}
	CodePayloadUpdateSigningFailure            = RequestError{"PAYLOAD_UPDATE_SIGNING_FAILURE", CategoryInternalMeta}
	CodeFunctionNotFound                       = RequestError{"FUNCTION_NOT_FOUND", CategoryFunctionNotFoundMeta}
	CodeMemoryWriteError                       = RequestError{"MEMORY_WRITE_ERROR", CategoryInternalMeta}
	CodeFailedExtractingResultBytes            = RequestError{"FAILED_EXTRACTING_RESULT_BYTES", CategoryInternalMeta}
	CodeFailedLoadingOrGettingModule           = RequestError{"FAILED_LOADING_OR_GETTING_MODULE", CategoryInternalMeta}
	CodeDepositFailed                          = RequestError{"DEPOSIT_FAILED", CategoryDepositFailedMeta}
	CodeRequestFuncFailed                      = RequestError{"REQUEST_FUNC_FAILED", CategoryRequestFuncFailedMeta}
	CodeFailedToGenerateReport                 = RequestError{"FAILED_TO_GENERATE_REPORT", CategoryDeanonymizationReportFailedMeta}
	CodeDeanonymizationReportEncryptionFailure = RequestError{"DEANONYMIZATION_REPORT_ENCRYPTION_FAILURE", CategoryDeanonymizationReportFailedMeta}
)

type RequestFailure struct {
	RequestError    RequestError
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
	RequestError RequestError `json:"requestError"`
	Message    string `json:"message"`
	ErrMessage string `json:"errMessage,omitempty"`
	ErrType    string `json:"errType,omitempty"`
}

func New(code RequestError, msg string, cause error) *RequestFailure {
	return &RequestFailure{RequestError: code, Message: msg, Err: cause}
}

// ToDTO converts a RequestFailure into its transport-friendly representation.
func (e *RequestFailure) ToDTO() *RequestFailureDTO {
	if e == nil {
		return nil
	}

	dto := &RequestFailureDTO{
		RequestError:    e.RequestError,
		Message: e.Message,
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
		RequestError: dto.RequestError,
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

func (e *RequestFailure) Category() Category {
	return e.categoryMeta().Category
}

func (e *RequestFailure) ExternalMessage() string {
	return e.categoryMeta().Message
}

func (e *RequestFailure) categoryMeta() ErrorCategory {
	if e == nil {
		return CategoryUnknownMeta
	}

	meta := e.RequestError.Category
	if meta.Message != "" || e.RequestError.Code == "" {
		return meta
	}
	return CategoryUnknownMeta
}

func FromError(err error) (*RequestFailure, bool) {
	var failure *RequestFailure
	if errors.As(err, &failure) {
		return failure, true
	}
	return nil, false
}