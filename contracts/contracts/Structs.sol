// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

//they are isolated so they can be imported in multiple contracts without loops
contract Structs {
    //type definitions
    enum RequestType { DEPLOYAPP, PROCESS, DEANONYMIZATION, ASSOCIATEKEY }
    enum RequestResult { COMPLETED, FAILED_REFUNDED, FAILED_NOT_REFUNDED }
    enum ErrorCode { 
        NO_ERROR,
        UNKNOWN,
        INTERNAL,
        APP_NOT_ADMITTED, 
        APPLICATION_ALREADY_DEPLOYED, 
        FAILURE_WHEN_DEPLOYING_APPLICATION, 
        DEANONYMIZATION_REPORT_FAILED, 
        REQUEST_TYPE_NOT_PERMITTED,
        FUNCTION_NOT_FOUND,
        DEPOSIT_FAILED,
        REQUEST_FUNC_FAILED,
        APP_NOT_DEPLOYED,
        WRONG_KEY_SENT,
        PUB_KEY_NOT_REGISTERED,
        NO_REPORT_DATA_FOUND,
        WASM_INTERNAL,
        INSUFFICIENT_FUEL
    }
    
    struct PendingRequest {
        uint256 timestamp; //assigned automatically
        uint256 depositAmount; //assigned automatically
        uint256 maxFeeValue;
        bytes32 requestId; //assigned automatically
        bytes payload;
        address sender; //assigned automatically
        uint64 applicationId;
        uint8 protocolVersion;
        RequestType requestType;
    }

    struct WithdrawalRequest {
        address payable receiver;
        uint256 amount;
    }
}
