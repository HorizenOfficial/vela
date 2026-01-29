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

    uint256 public minFeePerRequest;
    address payable public feeCollector;

    // Pull payment pattern state
    mapping(address => uint256) private _deposits;
    uint256 private _totalDeposits;

    //events
    event Refund(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount);
    event Withdrawal(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount);
    event RequestSubmitted(bytes32 indexed requestId, address indexed sender);
    event RequestCompleted(bytes32 indexed requestId, uint256 applicationFees, Structs.RequestResult status, Structs.ErrorCode errorCode, string errorMessage);
    event UserEvent(uint64 indexed applicationId, bytes32 indexed requestId, string indexed eventSubType, bytes encryptedData);
    event StateRootUpdate(uint64 indexed applicationId, bytes32 indexed requestId, bytes32 oldStateRoot, bytes32 newStateRoot);
    event QueueThresholdUpdated(uint256 newThreshold);
    event FeeCollectorUpdated(address newFeeCollector);

    //errors
    error AddressCantBeZero();
    error FeeValueBelowMinimum();
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
    error TransferFailed();


    modifier validProtocolVersion(uint8 protocolVersion) {
        if(protocolVersion != PROTOCOL_VERSION) revert InvalidProtocolVersion();
        _;
    }

    modifier validApplicationId(uint64 applicationId) {
        if(applicationId != APPLICATION_ID) revert InvalidApplicationId();
        _;
    }

    //constructor
    constructor(ITeeAuthenticator _teeAuthenticator, AuthorityRegistry _authorityRegistry, address updateStatusOperator, address admin, uint256 _minFeePerRequest) {
        if(
            _teeAuthenticator == ITeeAuthenticator(address(0)) || 
            _authorityRegistry == AuthorityRegistry(address(0)) ||
            updateStatusOperator == address(0) ||
            admin == address(0)
        ) revert AddressCantBeZero();

        teeAuthenticator = _teeAuthenticator;
        authorityRegistry = _authorityRegistry; 
        feeCollector = payable(updateStatusOperator);
        _grantRole(UPDATE_STATUS_ROLE, updateStatusOperator);
        _grantRole(ADMIN, admin);
        minFeePerRequest = _minFeePerRequest;
    }

    //request management functions
    function submitRequest(
        uint8 protocolVersion, 
        uint64 applicationId, 
        Structs.RequestType requestType, 
        bytes calldata payload, 
        uint256 depositAmount, // part of the sent value forwarded to the application, for app logic
        uint256 maxFeeValue // part ot the sent value reserved for fee payment
    ) validProtocolVersion(protocolVersion) validApplicationId(applicationId) payable public returns(bytes32) {
        //check values
        if(msg.value != depositAmount + maxFeeValue) revert InvalidValue();
        if(maxFeeValue < minFeePerRequest) revert FeeValueBelowMinimum();

        //check queue size
        if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

        if (requestType == Structs.RequestType.ASSOCIATEKEY) {
            //if requestype is associatekey, the payload must be 133 bytes long (contains a Secp521r1_PubKey)
            if (payload.length != 133) revert InvalidPayload();
        } else if (requestType == Structs.RequestType.DEANONYMIZATION) {

            // deanonymization requests MUST have depositAmount = 0
            if (depositAmount != 0) revert InvalidValue();

            // only allowed authorities can request deanonymization
            if (!authorityRegistry.checkAuthorityIsAllowed(applicationId, msg.sender)) {
                revert AuthorityNotAllowed();
            }
        }

        
        //create request
        bytes32 requestId = generateRequestId(msg.sender, applicationId, requestType, payload, depositAmount, _tail);
        requestById[requestId] = 
            Structs.PendingRequest(
                protocolVersion,
                applicationId,
                requestType,
                requestId,
                payload,
                block.timestamp,
                msg.sender,
                depositAmount,
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

    function markRequestCompleted(bytes32 requestId, uint256 refund, uint256 applicationFees) public onlyRole(UPDATE_STATUS_ROLE) {
        if (!isCurrentPendingRequest(requestId)) revert InvalidRequestId();

        //check values
        Structs.PendingRequest memory requestInfo = requestById[requestId];
        if(refund + applicationFees != requestInfo.maxFeeValue) revert InvalidValue();
        if(applicationFees < minFeePerRequest) {
            revert InvalidValue();
        }

        //credit refund to sender's pending balance (pull pattern)
        if(refund > 0) {
            _asyncTransfer(requestInfo.sender, refund);
            emit Refund(requestInfo.applicationId, requestId, requestInfo.sender, refund);
        }

        //credit fee to feeCollector's pending balance
        _asyncTransfer(feeCollector, applicationFees);

        _markRequestCompleted(requestId, applicationFees);
    }

    function _markRequestCompleted(bytes32 requestId, uint256 applicationFees) private {

       _removeRequest();

        emit RequestCompleted(requestId, applicationFees, Structs.RequestResult.COMPLETED, Structs.ErrorCode.NO_ERROR, "");
    }

    // We return the maxValueFee - minFeePerRequest (to be changed in the future)
    function markRequestFailed(bytes32 requestId, Structs.ErrorCode errorCode, string memory errorMessage) public onlyRole(UPDATE_STATUS_ROLE) {
        if (!isCurrentPendingRequest(requestId)) revert InvalidRequestId();

        address sender = requestById[requestId].sender;
        uint256 depositAmount = requestById[requestId].depositAmount;
        uint256 maxFeeValue = requestById[requestId].maxFeeValue;
        uint64 applicationId = requestById[requestId].applicationId;

        _removeRequest();

        //credit refund to sender's pending balance (pull pattern)
        uint256 refund = depositAmount + (maxFeeValue - minFeePerRequest);
        if(refund > 0) {
            _asyncTransfer(sender, refund);
            emit Refund(applicationId, requestId, sender, refund);
        }

        //credit minimum fee to feeCollector's pending balance
        _asyncTransfer(feeCollector, minFeePerRequest);

        emit RequestCompleted(requestId, minFeePerRequest, Structs.RequestResult.FAILED, errorCode, errorMessage);
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
        string[] memory eventSubTypes,
        Structs.WithdrawalRequest[] memory withdrawalRequests, 
        uint256 refund,
        uint256 applicationFees,
        bytes memory signature
    ) public validApplicationId(applicationId) onlyRole(UPDATE_STATUS_ROLE) {
        //check prev state root
        if(stateRoot != bytes32(0) && prevStateRoot != stateRoot) revert InvalidStateRoot();
        //check valid request
        if (!isCurrentPendingRequest(processedRequestId)) revert InvalidRequestId();

        //check signature
        if (events.length != eventSubTypes.length) revert InvalidPayload();
        if(!teeAuthenticator.checkSignature(applicationId, prevStateRoot, newStateRoot, processedRequestId, events, eventSubTypes, withdrawalRequests, refund, applicationFees, signature)) revert InvalidSignature();

        //check values
        Structs.PendingRequest memory requestInfo = requestById[processedRequestId];
        if(refund + applicationFees != requestInfo.maxFeeValue) revert InvalidValue();
        if(applicationFees < minFeePerRequest) {
            revert InvalidValue();
        }

        //check withdrawal sums (account for already committed pending deposits)
        uint256 i;
        uint256 sum;
        while(i < withdrawalRequests.length) {
            sum += withdrawalRequests[i].amount;
            unchecked {++i;}
        }
        sum += refund + applicationFees;
        if(sum > address(this).balance - _totalDeposits) revert InsufficientBalance();

        //set requests as completed
        _markRequestCompleted(processedRequestId, applicationFees);

        //emit encrypted event
        i = 0;
        while(i < events.length) {
            emit UserEvent(applicationId, processedRequestId, eventSubTypes[i], events[i]);
            unchecked {++i;}
        }

        //update state root and request
        stateRoot = newStateRoot;
        emit StateRootUpdate(applicationId, processedRequestId, prevStateRoot, newStateRoot);

        //credit refund to sender's pending balance (pull pattern)
        if(refund > 0) {
            _asyncTransfer(requestInfo.sender, refund);
            emit Refund(applicationId, processedRequestId, requestInfo.sender, refund);
        }

        //credit fee to feeCollector's pending balance
        _asyncTransfer(feeCollector, applicationFees);

        //credit withdrawals to receivers' pending balances
        i = 0;
        while(i < withdrawalRequests.length) {
            _asyncTransfer(withdrawalRequests[i].receiver, withdrawalRequests[i].amount);
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

    // Pull payment pattern functions
    function _asyncTransfer(address dest, uint256 amount) internal {
        _deposits[dest] += amount;
        _totalDeposits += amount;
    }

    function withdrawPayments(address payable payee) public {
        uint256 payment = _deposits[payee];
        if (payment == 0) return;

        _deposits[payee] = 0;
        _totalDeposits -= payment;

        (bool success, ) = payee.call{value: payment}("");
        if (!success) revert TransferFailed();
    }

    function payments(address dest) public view returns (uint256) {
        return _deposits[dest];
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
        uint256 depositAmount,
        uint256 idx
        ) public pure returns (bytes32) {
        
        bytes32 requestId = keccak256(abi.encodePacked(
            sender,
            applicationId,
            requestType,
            payload,
            depositAmount,
            idx
       ));

        return requestId;
    }

}
