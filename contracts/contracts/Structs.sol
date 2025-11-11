// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

//they are isolated so they can be imported in multiple contracts without loops
contract Structs {
    //type definitions
    enum RequestType { DEPLOYAPP, PROCESS, DEANONYMIZATION, ASSOCIATEKEY }
    enum RequestResult { COMPLETED, FAILED_REFUNDED, FAILED_NOT_REFUNDED }
    enum ErrorCode { 
        NO_ERROR,
        UNKNOWN_ERROR,
        INTERNAL,
        APP_NOT_ADMITTED, 
        APPLICATION_ALREADY_DEPLOYED, 
        FAILURE_WHEN_DEPLOYING_APPLICATION, 
        DEANONYMIZATION_REPORT_FAILED, 
        REQUEST_TYPE_NOT_PERMITTED,
        SENDER_ACCOUNT_INEXISTENT,
        FUNCTION_NOT_FOUND,
        DEPOSIT_FAILED,
        REQUEST_FUNC_FAILED,
        APPLICATION_NOT_DEPLOYED
    }
    
    struct PendingRequest {
        uint8 protocolVersion;
        uint256 applicationId;
        RequestType requestType;
        bytes32 requestId; //assigned automatically
        bytes payload;
        uint256 timestamp; //assigned automatically
        address sender; //assigned automatically
        uint256 value; //assigned automatically
    }

    struct WithdrawalRequest {
        address payable receiver;
        uint256 amount;
    }
}