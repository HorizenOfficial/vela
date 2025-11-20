// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/access/AccessControl.sol";

import "./interfaces/ITeeAuthenticator.sol";
import "./AuthorityRegistry.sol";
import "./Structs.sol";

contract ProcessorEndpoint is AccessControl {

    //constants
    bytes32 public constant UPDATE_STATUS_ROLE = keccak256("UPDATE_STATUS_ROLE");
    bytes32 public constant ADMIN = keccak256("ADMIN");
    uint8 public constant PROTOCOL_VERSION = 0;
    uint64 public constant APPLICATION_ID = 1;
    
    //state variables
    bytes32 public stateRoot;
    
    mapping(bytes32 => Structs.PendingRequest) public requestById;
    mapping(uint256 => bytes32) private _requestIdByOrder;
    uint256 private _head;
    uint256 private _tail;
    uint256 public maxQueueSize = 10;

    ITeeAuthenticator public teeAuthenticator;
    AuthorityRegistry public authorityRegistry;

    uint256 public minFeePerRequest = 5;
    address payable public feeCollector = payable(0x574Cb6eD4De4167cb31C67ab9F97Ca8472a18973);

    //events
    event Refund(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount);
    event Withdrawal(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount);
    event RequestSubmitted(bytes32 indexed requestId, address indexed sender);
    event RequestCompleted(bytes32 indexed requestId, Structs.RequestResult status, Structs.ErrorCode errorCode, string errorMessage, uint256 applicationFees);
    event UserEvent(uint64 indexed applicationId, bytes32 indexed requestId, bytes encryptedData);
    event StateRootUpdate(uint64 indexed applicationId, bytes32 indexed requestId, bytes32 oldStateRoot, bytes32 newStateRoot);
    event QueueThresholdUpdated(uint256 newThreshold);
    event FeeCollectorUpdated(address newFeeCollector);

    //errors
    error AddressCantBeZero();
    error InvalidValue();
    error InvalidProtocolVersion();
    error InvalidApplicationId();
    error InvalidRequestId();
    error InvalidStateRoot();
    error InvalidSignature();
    error InvalidPayload();
    error InsufficientBalance();
    error AuthorityNotAllowed();
    error QueueThresholdExceeded();


    modifier validProtocolVersion(uint8 protocolVersion) {
        if(protocolVersion != PROTOCOL_VERSION) revert InvalidProtocolVersion();
        _;
    }

    modifier validApplicationId(uint64 applicationId) {
        if(applicationId != APPLICATION_ID) revert InvalidApplicationId();
        _;
    }

    //constructor
    constructor(ITeeAuthenticator _teeAuthenticator, AuthorityRegistry _authorityRegistry, address updateStatusOperator, address admin) {
        if(
            _teeAuthenticator == ITeeAuthenticator(address(0)) || 
            _authorityRegistry == AuthorityRegistry(address(0)) ||
            updateStatusOperator == address(0) ||
            admin == address(0)
        ) revert AddressCantBeZero();

        teeAuthenticator = _teeAuthenticator;
        authorityRegistry = _authorityRegistry; 
        _grantRole(UPDATE_STATUS_ROLE, updateStatusOperator);
        _grantRole(ADMIN, admin);
    }

    //request management functions
    function submitRequest(
        uint8 protocolVersion, 
        uint64 applicationId, 
        Structs.RequestType requestType, 
        bytes calldata payload, 
        uint256 value, // part of the sent value forwarded to the application, for app logic
        uint256 maxFeeValue // part ot the sent value reserved for fee payment
    ) validProtocolVersion(protocolVersion) validApplicationId(applicationId) payable public returns(bytes32) {
        //check values
        if(msg.value != value + maxFeeValue) revert InvalidValue();
        if(maxFeeValue < minFeePerRequest) revert InvalidValue();

        //check queue size
        if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

        if (requestType == Structs.RequestType.ASSOCIATEKEY) {
            //if requestype is associatekey, the payload must be 133 bytes long (contains a Secp521r1_PubKey)
            if (payload.length != 133) revert InvalidPayload();
        }else if  (requestType == Structs.RequestType.DEANONYMIZATION && !authorityRegistry.checkAuthorityIsAllowed(applicationId, msg.sender)) revert AuthorityNotAllowed();
        
        //create request
        bytes32 requestId = generateRequestId(msg.sender, applicationId, requestType, payload, value, _tail);
        requestById[requestId] = 
            Structs.PendingRequest(
                protocolVersion,
                applicationId,
                requestType,
                requestId,
                payload,
                block.timestamp,
                msg.sender,
                value,
                maxFeeValue
            );
        _requestIdByOrder[_tail] = requestId;

        _tail++;

        //emit event
        emit RequestSubmitted(requestId, msg.sender);

        return requestId;
    }

    function _removeRequest() private {

        delete requestById[_requestIdByOrder[_head]];
        delete _requestIdByOrder[_head];
        _head++;

    }

    function markRequestCompleted(bytes32 requestId) public onlyRole(UPDATE_STATUS_ROLE) {
        if (!isCurrentPendingRequest(requestId)) revert InvalidRequestId();
        //TODO deanonymization calls this (bool feeSent, ) = payable(feeCollector).call{value: applicationFees}("");
        _markRequestCompleted(requestId, minFeePerRequest);
    }

    function _markRequestCompleted(bytes32 requestId, uint256 applicationFees) private {

       _removeRequest();

        emit RequestCompleted(requestId, Structs.RequestResult.COMPLETED, Structs.ErrorCode.NO_ERROR, "", applicationFees);
    }

    // We return the maxValueFee - minFeePerRequest (to be changed in the future)
    function markRequestFailed(bytes32 requestId, Structs.ErrorCode errorCode, string memory errorMessage) public onlyRole(UPDATE_STATUS_ROLE) {
        if (!isCurrentPendingRequest(requestId)) revert InvalidRequestId();

        address sender = requestById[requestId].sender;
        uint256 value = requestById[requestId].value;
        uint256 maxFeeValue = requestById[requestId].maxFeeValue;

        _removeRequest();
        //refunds
        (bool refunded, ) = payable(sender).call{value: value + (maxFeeValue - minFeePerRequest)}(""); // TODO think about failing if transfer fails (everywhere)

        //minimum fee is collected
        (bool feeSent, ) = payable(feeCollector).call{value: minFeePerRequest}("");

        if(refunded) emit RequestCompleted(requestId, Structs.RequestResult.FAILED_REFUNDED, errorCode, errorMessage, minFeePerRequest); 
        else emit RequestCompleted(requestId, Structs.RequestResult.FAILED_NOT_REFUNDED, errorCode, errorMessage, minFeePerRequest); 

    }

    function getPendingRequestsSize() public view returns(uint256) {
        if (_tail > _head) {
            return (_tail - _head);
        } else {
            return 0;
        }
        
    }

    function getPendingRequests() public view returns(Structs.PendingRequest[] memory) {
        uint256 numOfPendingRequests = getPendingRequestsSize();

        Structs.PendingRequest[] memory res = new Structs.PendingRequest[](numOfPendingRequests);
        uint256 i = _head;
        uint256 j;
        while(i < _tail) {
            bytes32 requestId = _requestIdByOrder[i];
            res[j] = requestById[requestId];
            unchecked { ++i; ++j;}
        }

        return res;
    }

    //update status
    function stateUpdate(
        uint64 applicationId, 
        bytes32 prevStateRoot, 
        bytes32 newStateRoot, 
        bytes32 processedRequestId,
        bytes[] memory events, 
        Structs.WithdrawalRequest[] memory withdrawalRequests, 
        bytes memory signature,
        uint256 refund,
        uint256 applicationFees
    ) public validApplicationId(applicationId) onlyRole(UPDATE_STATUS_ROLE) {
        //check prev state root
        if(stateRoot != bytes32(0) && prevStateRoot != stateRoot) revert InvalidStateRoot();
        //check valid request
        if (!isCurrentPendingRequest(processedRequestId)) revert InvalidRequestId();

        //check signature
        if(!teeAuthenticator.checkSignature(applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature, refund, applicationFees)) revert InvalidSignature();

        //check values
        Structs.PendingRequest memory requestInfo = requestById[processedRequestId];
        if(refund + applicationFees != requestInfo.maxFeeValue) revert InvalidValue();  //TODO better errors for the future
        if(applicationFees < minFeePerRequest) {
            revert InvalidValue();
        }

        //check withdrawal sums 
        uint256 i;
        uint256 sum;
        while(i < withdrawalRequests.length) { // TODO optimize these loops and the contract overall
            sum += withdrawalRequests[i].amount;
            unchecked {++i;}
        }
        sum += refund + applicationFees;
        if(sum > address(this).balance) revert InsufficientBalance();

        //set requests as completed
        _markRequestCompleted(processedRequestId, applicationFees);

        //emit encrypted event
        i = 0;
        while(i < events.length) {
            emit UserEvent(applicationId, processedRequestId, events[i]);
            unchecked {++i;}
        }

        //update state root and request
        stateRoot = newStateRoot;
        emit StateRootUpdate(applicationId, processedRequestId, prevStateRoot, newStateRoot);

        (bool refundSent, ) = payable(requestInfo.sender).call{value: refund}("");
        emit Refund(applicationId, processedRequestId, requestInfo.sender, refund);

        (bool feeSent, ) = payable(feeCollector).call{value: applicationFees}("");
            
        //execute withdrawals (as last operation)
        i = 0;
        while(i < withdrawalRequests.length) {
            (bool withdrawn, ) = payable(withdrawalRequests[i].receiver).call{value: withdrawalRequests[i].amount}("");
            emit Withdrawal(applicationId, processedRequestId, withdrawalRequests[i].receiver, withdrawalRequests[i].amount);
            unchecked {++i;}
        }  
    }

    function updateQueueThreshold(uint256 newThreshold) public onlyRole(ADMIN) {
        if (newThreshold == 0) revert InvalidValue();
        maxQueueSize = newThreshold;
        emit QueueThresholdUpdated(newThreshold);
    }

    function updateFeeCollector(address payable newFeeCollector) public onlyRole(ADMIN) {
        if (newFeeCollector == address(0)) revert AddressCantBeZero();
        feeCollector = newFeeCollector;
        emit FeeCollectorUpdated(newFeeCollector);
    }

    function getNextPendingRequest() public view returns (Structs.PendingRequest memory, bytes32, bool success) {
        uint256 numOfRequests = getPendingRequestsSize();
        if (numOfRequests > 0){
            bytes32 requestId = _requestIdByOrder[_head];
            return (requestById[requestId], stateRoot, true);
        }

        Structs.PendingRequest memory emptyReq;
        return (emptyReq, stateRoot, false);
 
    }

    function isCurrentPendingRequest(bytes32 requestId) public view returns (bool) {
        return getPendingRequestsSize() > 0 && _requestIdByOrder[_head] == requestId;
    }

    function generateRequestId(
        address sender,
        uint64 applicationId, 
        Structs.RequestType requestType, 
        bytes calldata payload, 
        uint256 value,
        uint256 idx
        ) public pure returns (bytes32) {
        
        bytes32 requestId = keccak256(abi.encodePacked(
            sender,
            applicationId,
            requestType,
            payload,
            value,
            idx
       ));

        return requestId;
    }

}
