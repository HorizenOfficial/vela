package apperrors

import (
	"errors"
	"fmt"
)

type RequestErrorCode string

const (
	// Internal errors with more details
	CodeInternalFallback	RequestErrorCode = "INTERNAL_FALLBACK"
	CodeManagerNotRunning       RequestErrorCode = "MANAGER_NOT_RUNNING"
	CodeAppNotAdmitted		 RequestErrorCode = "APP_NOT_ADMITTED"
	CodeApplicationAlreadyDeployed RequestErrorCode = "APPLICATION_ALREADY_DEPLOYED"
	CodeWasmNotFound		 RequestErrorCode = "WASM_NOT_FOUND"
	CodeErrorReadingWasm		 RequestErrorCode = "ERROR_READING_WASM"
	CodeFailureDeployingApp	 RequestErrorCode = "FAILURE_WHEN_DEPLOYING_APPLICATION"
	CodeSubmittingStateUpdateFailed RequestErrorCode = "SUBMITTING_STATE_UPDATE_FAILED"
	CodeAppStateNotFound	 RequestErrorCode = "APPLICATION_STATE_NOT_FOUND"
	CodeDeanonymizationReportFailed RequestErrorCode = "DEANONYMIZATION_REPORT_FAILED"
	CodeRequestTypeNotPermitted RequestErrorCode = "REQUEST_TYPE_NOT_PERMITTED"
	CodeEncryptedToAppDataFailure RequestErrorCode = "ENCRYPTED_TO_APPDATA_FAILURE"
	CodeJsonUnmarshalError RequestErrorCode = "JSON_UNMARSHAL_ERROR"
	CodeJsonMarshalError RequestErrorCode = "JSON_MARSHAL_ERROR"
	CodeParsingKeyError RequestErrorCode = "PARSING_KEY_ERROR"
	CodePayloadDecryptionFailure RequestErrorCode = "PAYLOAD_DECRYPTION_FAILURE"
	CodeAppDataSerializationFailure RequestErrorCode = "APPDATA_SERIALIZATION_FAILURE"
	CodeAppDataEncryptionFailure RequestErrorCode = "APPDATA_ENCRYPTION_FAILURE"
	CodeEventsEncryptionFailure RequestErrorCode = "EVENTS_ENCRYPTION_FAILURE"
	CodePayloadUpdateSigningFailure	 RequestErrorCode = "PAYLOAD_UPDATE_SIGNING_FAILURE"
	CodeFunctionNotFound RequestErrorCode = "FUNCTION_NOT_FOUND"
	CodeMemoryWriteError RequestErrorCode = "MEMORY_WRITE_ERROR"
	CodeFailedExtractingResultBytes RequestErrorCode = "FAILED_EXTRACTING_RESULT_BYTES"
	CodeDepositFailed RequestErrorCode = "DEPOSIT_FAILED"
	CodeFailedLoadingOrGettingModule RequestErrorCode = "FAILED_LOADING_OR_GETTING_MODULE"
	CodeRequestFuncFailed RequestErrorCode = "REQUEST_FUNC_FAILED"
)

type ErrorCategory uint8

const (
	// External error categories for users and blockchain
	CategoryNoError ErrorCategory = iota
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
)

var codeToCategory = map[RequestErrorCode]ErrorCategory{
	CodeInternalFallback:       CategoryInternal,
	CodeManagerNotRunning:      CategoryInternal,
	CodeSubmittingStateUpdateFailed: CategoryInternal,
	CodeAppStateNotFound:		CategoryInternal,
	CodeEncryptedToAppDataFailure: CategoryInternal,
	CodeJsonUnmarshalError: CategoryInternal,
	CodeJsonMarshalError: CategoryInternal,
	CodeParsingKeyError: CategoryInternal,
	CodePayloadDecryptionFailure: CategoryInternal,
	CodeAppDataSerializationFailure: CategoryInternal,
	CodeAppDataEncryptionFailure: CategoryInternal,
	CodeEventsEncryptionFailure: CategoryInternal,
	CodePayloadUpdateSigningFailure: CategoryInternal,
	CodeMemoryWriteError: CategoryInternal,
	CodeFailedExtractingResultBytes: CategoryInternal,
	CodeFailedLoadingOrGettingModule: CategoryInternal,
	CodeWasmNotFound:			CategoryFailureWhenDeployingApplication,
	CodeErrorReadingWasm:		CategoryFailureWhenDeployingApplication,
	CodeFailureDeployingApp:	CategoryFailureWhenDeployingApplication,
	CodeAppNotAdmitted:         CategoryAppNotAdmitted,
	CodeApplicationAlreadyDeployed: CategoryApplicationAlreadyDeployed,
	CodeDeanonymizationReportFailed: CategoryDeanonymizationReportFailed,
	CodeRequestTypeNotPermitted: CategoryRequestTypeNotPermitted,
	CodeFunctionNotFound: CategoryFunctionNotFound,
	CodeDepositFailed: CategoryDepositFailed,
	CodeRequestFuncFailed: CategoryRequestFuncFailed,
}

var defaultMessages = map[ErrorCategory]string{
	CategoryNoError:  "no error",
	CategoryUnknown: "unknown error",
	CategoryInternal:   "internal error",
	CategoryAppNotAdmitted: "application is not admitted",
	CategoryApplicationAlreadyDeployed: "application is already deployed",
	CategoryFailureWhenDeployingApplication: "failed to deploy application",
	CategoryDeanonymizationReportFailed: "failed to generate deanonymization report",
	CategoryRequestTypeNotPermitted: "request type is not permitted",
	CategoryFunctionNotFound: "requested function not found",
	CategoryDepositFailed: "deposit operation failed",
	CategoryRequestFuncFailed: "request function execution failed",
}

type RequestFailure struct {
	Code            RequestErrorCode
	Message         string
	Err             error
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
	Code       RequestErrorCode `json:"code"`
	Message    string           `json:"message"`
	ErrMessage string           `json:"errMessage,omitempty"`
	ErrType    string           `json:"errType,omitempty"`
}

func New(code RequestErrorCode, msg string, cause error) *RequestFailure {
	return &RequestFailure{Code: code, Message: msg, Err: cause}
}

// ToDTO converts a RequestFailure into its transport-friendly representation.
func (e *RequestFailure) ToDTO() *RequestFailureDTO {
	if e == nil {
		return nil
	}

	dto := &RequestFailureDTO{
		Code:    e.Code,
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
		Code:    dto.Code,
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

func (e *RequestFailure) Category() ErrorCategory {
	if cat, ok := codeToCategory[e.Code]; ok {
		return cat
	}
	return CategoryUnknown
}

func (e *RequestFailure) ExternalMessage() string { 
	if msg, ok := defaultMessages[e.Category()]; ok {
		return msg
	}
	return defaultMessages[CategoryUnknown]
}

func FromError(err error) (*RequestFailure, bool) {
	var failure *RequestFailure
	if errors.As(err, &failure) {
		return failure, true
	}
	return nil, false
}