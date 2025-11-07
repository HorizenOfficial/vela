package apperrors

import (
	"errors"
	"fmt"
)

type RequestErrorCode string
//TODO add testing for these errors
//TODO reorder stuff

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
	CodeNilInstruction RequestErrorCode = "NIL_INSTRUCTION"
	CodeSenderAccountInexistent RequestErrorCode = "SENDER_ACCOUNT_INEXISTENT"
	CodeInsufficientBalance RequestErrorCode = "INSUFFICIENT_BALANCE"
	CodeUnknownInstructionType RequestErrorCode = "UNKNOWN_INSTRUCTION_TYPE"
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

type ErrorCategory string

const (
	// External error categories for users and blockchain
	CategoryInternal   ErrorCategory = "INTERNAL"
	CategoryAppNotAdmitted ErrorCategory = "APP_NOT_ADMITTED"
	CategoryApplicationAlreadyDeployed ErrorCategory = "APPLICATION_ALREADY_DEPLOYED"
	CategoryFailureWhenDeployingApplication ErrorCategory = "FAILURE_WHEN_DEPLOYING_APPLICATION"
	CategoryDeanonymizationReportFailed ErrorCategory = "DEANONYMIZATION_REPORT_FAILED"
	CategoryRequestTypeNotPermitted ErrorCategory = "REQUEST_TYPE_NOT_PERMITTED"
	CategorySenderAccountInexistent ErrorCategory = "SENDER_ACCOUNT_INEXISTENT"
	CategoryInsufficientBalance ErrorCategory = "INSUFFICIENT_BALANCE"
	CategoryFunctionNotFound ErrorCategory = "FUNCTION_NOT_FOUND"
	CategoryDepositFailed ErrorCategory = "DEPOSIT_FAILED"
	CategoryRequestFuncFailed ErrorCategory = "REQUEST_FUNC_FAILED"
)

var codeToCategory = map[RequestErrorCode]ErrorCategory{
	CodeInternalFallback:       CategoryInternal,
	CodeManagerNotRunning:      CategoryInternal,
	CodeAppNotAdmitted:         CategoryAppNotAdmitted,
	CodeApplicationAlreadyDeployed: CategoryApplicationAlreadyDeployed,
	CodeWasmNotFound:			CategoryFailureWhenDeployingApplication,
	CodeErrorReadingWasm:		CategoryFailureWhenDeployingApplication,
	CodeFailureDeployingApp:	CategoryFailureWhenDeployingApplication,
	CodeSubmittingStateUpdateFailed: CategoryInternal,
	CodeAppStateNotFound:		CategoryInternal,
	CodeDeanonymizationReportFailed: CategoryDeanonymizationReportFailed,
	CodeRequestTypeNotPermitted: CategoryRequestTypeNotPermitted,
	CodeEncryptedToAppDataFailure: CategoryInternal,
	CodeJsonUnmarshalError: CategoryInternal,
	CodeJsonMarshalError: CategoryInternal,
	CodeNilInstruction: CategoryInternal,
	CodeSenderAccountInexistent: CategorySenderAccountInexistent,
	CodeInsufficientBalance: CategoryInsufficientBalance,
	CodeUnknownInstructionType: CategoryInternal,
	CodeParsingKeyError: CategoryInternal,
	CodePayloadDecryptionFailure: CategoryInternal,
	CodeAppDataSerializationFailure: CategoryInternal,
	CodeAppDataEncryptionFailure: CategoryInternal,
	CodeEventsEncryptionFailure: CategoryInternal,
	CodePayloadUpdateSigningFailure: CategoryInternal,
	CodeFunctionNotFound: CategoryFunctionNotFound,
	CodeMemoryWriteError: CategoryInternal,
	CodeFailedExtractingResultBytes: CategoryInternal,
	CodeDepositFailed: CategoryDepositFailed,
	CodeFailedLoadingOrGettingModule: CategoryInternal,
	CodeRequestFuncFailed: CategoryRequestFuncFailed,
}

var defaultMessages = map[ErrorCategory]string{
	CategoryInternal:   "internal error",
	CategoryAppNotAdmitted: "application is not admitted",
	CategoryApplicationAlreadyDeployed: "application is already deployed",
	CategoryFailureWhenDeployingApplication: "failed to deploy application",
	CategoryDeanonymizationReportFailed: "failed to generate deanonymization report",
	CategoryRequestTypeNotPermitted: "request type is not permitted",
	CategorySenderAccountInexistent: "sender account does not exist",
	CategoryInsufficientBalance: "insufficient balance",
	CategoryFunctionNotFound: "requested function not found",
	CategoryDepositFailed: "deposit operation failed",
	CategoryRequestFuncFailed: "request function execution failed",
}

type RequestFailure struct {
	Code            RequestErrorCode
	Message         string
	Err             error
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
		if dto.ErrType != "" {
			failure.Err = fmt.Errorf("%s: %s", dto.ErrType, dto.ErrMessage)
		} else {
			failure.Err = errors.New(dto.ErrMessage)
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
	return CategoryInternal
}

func (e *RequestFailure) ExternalMessage() string { 
	if msg, ok := defaultMessages[e.Category()]; ok {
		return msg
	}
	return defaultMessages[CategoryInternal]
}

func FromError(err error) (*RequestFailure, bool) {
	var failure *RequestFailure
	if errors.As(err, &failure) {
		return failure, true
	}
	return nil, false
}

func ToSolidityEnum(cat ErrorCategory) uint8 {
	switch cat {
	case CategoryInternal:
		return 1
	case CategoryAppNotAdmitted:
		return 2
	case CategoryApplicationAlreadyDeployed:
		return 3
	case CategoryFailureWhenDeployingApplication:
		return 4
	case CategoryDeanonymizationReportFailed:
		return 5
	case CategoryRequestTypeNotPermitted:
		return 6
	case CategorySenderAccountInexistent:
		return 7	
	case CategoryInsufficientBalance:
		return 8
	case CategoryFunctionNotFound:
		return 9	
	case CategoryDepositFailed:
		return 10
	case CategoryRequestFuncFailed:
		return 11
	default:
		return 0
	}
}
