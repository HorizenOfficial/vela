// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

//they are isolated so they can be imported in multiple contracts without loops
contract Structs {
    //type definitions
    enum RequestType { DEPLOYAPP, PROCESS, DEANONYMIZATION }
    enum RequestStatus { POSTED, COMPLETED, FAILED }
    
    struct PendingRequest {
        uint8 protocolVersion;
        uint256 applicationId;
        RequestType requestType;
        uint256 requestId; //assigned automatically
        bytes payload;
        uint256 timestamp; //assigned automatically
        address sender; //assigned automatically
        RequestStatus status; //assigned automatically
        uint256 value; //assigned automatically
    }

    struct WithdrawalRequest {
        address payable receiver;
        uint256 amount;
    }
}