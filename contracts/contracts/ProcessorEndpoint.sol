// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/access/AccessControl.sol";

import "./interfaces/ITeeAuthenticator.sol";
import "./AuthorityRegistry.sol";
import "./Structs.sol";

contract ProcessorEndpoint is AccessControl {

    //constants
    bytes32 public constant UPDATE_STATUS_ROLE = keccak256("UPDATE_STATUS_ROLE");
    uint8 public constant PROTOCOL_VERSION = 0;
    uint256 public constant APPLICATION_ID = 1;
    
    //state variables
    bytes32 public stateRoot;
    
    mapping(bytes32 => Structs.PendingRequest) public requestById;
    mapping(uint256 => bytes32) private _requestIdByOrder;
    uint256 private _head;
    uint256 private _tail;

    ITeeAuthenticator public teeAuthenticator;
    AuthorityRegistry public authorityRegistry;
    //events
    event Withdrawal(uint256 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount);
    event RequestSubmitted(bytes32 indexed requestId, address indexed sender);
    event RequestCompleted(bytes32 indexed requestId, Structs.RequestResult status);
    event UserEvent(uint256 indexed applicationId, bytes32 indexed requestId, bytes encryptedData);
    event StateRootUpdate(uint256 indexed applicationId, bytes32 indexed requestId, bytes32 oldStateRoot, bytes32 newStateRoot);
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


    modifier validProtocolVersion(uint8 protocolVersion) {
        if(protocolVersion != PROTOCOL_VERSION) revert InvalidProtocolVersion();
        _;
    }

    modifier validApplicationId(uint256 applicationId) {
        if(applicationId != APPLICATION_ID) revert InvalidApplicationId();
        _;
    }

    //constructor
    constructor(ITeeAuthenticator _teeAuthenticator, AuthorityRegistry _authorityRegistry, address updateStatusOperator) {
        if(
            _teeAuthenticator == ITeeAuthenticator(address(0)) || 
            _authorityRegistry == AuthorityRegistry(address(0)) ||
            updateStatusOperator == address(0)
        ) revert AddressCantBeZero();

        teeAuthenticator = _teeAuthenticator;
        authorityRegistry = _authorityRegistry; 
        _grantRole(UPDATE_STATUS_ROLE, updateStatusOperator);
    }

    //request management functions
    function submitRequest(
        uint8 protocolVersion, 
        uint256 applicationId, 
        Structs.RequestType requestType, 
        bytes calldata payload, 
        uint256 value
    ) validProtocolVersion(protocolVersion) validApplicationId(applicationId) payable public returns(bytes32) {
        //check value
        if(msg.value != value) revert InvalidValue(); //'value' is redundant now, but it will be needed when using ERC20

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
                value
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

        _markRequestCompleted(requestId);
    }

    function _markRequestCompleted(bytes32 requestId) private {

       _removeRequest();

        emit RequestCompleted(requestId, Structs.RequestResult.COMPLETED);
    }

    function markRequestFailed(bytes32 requestId) public onlyRole(UPDATE_STATUS_ROLE) {
        if (!isCurrentPendingRequest(requestId)) revert InvalidRequestId();

        address sender = requestById[requestId].sender;
        uint256 value = requestById[requestId].value;

       _removeRequest();

        if (value == 0) {
            emit RequestCompleted(requestId, Structs.RequestResult.FAILED_REFUNDED); 
            return;
        }
        //refunds
        (bool refunded, ) = payable(sender).call{value: value}("");

        if(refunded) emit RequestCompleted(requestId, Structs.RequestResult.FAILED_REFUNDED); 
        else emit RequestCompleted(requestId, Structs.RequestResult.FAILED_NOT_REFUNDED); 

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
        uint256 applicationId, 
        bytes32 prevStateRoot, 
        bytes32 newStateRoot, 
        bytes32 processedRequestId,
        bytes[] memory events, 
        Structs.WithdrawalRequest[] memory withdrawalRequests, 
        bytes memory signature
    ) public validApplicationId(applicationId) onlyRole(UPDATE_STATUS_ROLE) {

        //check prev state root
        if(stateRoot != bytes32(0) && prevStateRoot != stateRoot) revert InvalidStateRoot();
        //check valid request
        if (!isCurrentPendingRequest(processedRequestId)) revert InvalidRequestId();

        //check signature
        if(!teeAuthenticator.checkSignature(applicationId, prevStateRoot, newStateRoot, processedRequestId, events, withdrawalRequests, signature)) revert InvalidSignature();

        //check withdrawal sums 
        uint256 i;
        uint256 sum;
        while(i < withdrawalRequests.length) {
            sum += withdrawalRequests[i].amount;
            unchecked {++i;}
        }
        if(sum > address(this).balance) revert InsufficientBalance();

        //set requests as completed
        _markRequestCompleted(processedRequestId);

        //emit encrypted event
        i = 0;
        while(i < events.length) {
            emit UserEvent(applicationId, processedRequestId, events[i]);
            unchecked {++i;}
        }

        //update state root and request
        stateRoot = newStateRoot;
        emit StateRootUpdate(applicationId, processedRequestId, prevStateRoot, newStateRoot);

        //execute withdrawals (as last operation)
        i = 0;
        while(i < withdrawalRequests.length) {
            withdrawalRequests[i].receiver.transfer(withdrawalRequests[i].amount);
            emit Withdrawal(applicationId, processedRequestId, withdrawalRequests[i].receiver, withdrawalRequests[i].amount);
            unchecked {++i;}
        }
    }


    function getNextPendingRequest() public view returns (Structs.PendingRequest memory, bytes32, bool success) {
        uint256 numOfRequests = getPendingRequestsSize();
        if (numOfRequests > 0){
            bytes32 requestId = _requestIdByOrder[_head];
            return (requestById[requestId], stateRoot, true);
        }

        Structs.PendingRequest memory emptyReq;
        return (emptyReq, bytes32(0), false);
 
    }

    function isCurrentPendingRequest(bytes32 requestId) public view returns (bool) {
        return getPendingRequestsSize() > 0 && _requestIdByOrder[_head] == requestId;
    }

    function generateRequestId(
        address sender,
        uint256 applicationId, 
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
