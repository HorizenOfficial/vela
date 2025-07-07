// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Arrays.sol";

import "./interfaces/ITeeAuthenticator.sol";
import "./Structs.sol";

contract ProcessorEndpoint is AccessControl, ReentrancyGuard {
    using EnumerableSet for EnumerableSet.UintSet;
    using Arrays for uint256[];

    //constants
    bytes32 public constant UPDATE_STATUS_ROLE = keccak256("UPDATE_STATUS_ROLE");
    uint8 public constant PROTOCOL_VERSION = 0;
    uint256 public constant APPLICATION_ID = 1;
    
    //state variables
    bytes public stateRoot;
    
    Structs.PendingRequest[] public requests;
    EnumerableSet.UintSet private idsQueue;

    ITeeAuthenticator public teeAuthenticator;
    //events
    event Withdrawal(uint256 indexed applicationId, uint256 indexed requestId, address to, uint256 amount);
    event RequestSubmitted(uint256 indexed requestId, address sender);
    event RequestCompleted(uint256 indexed requestId);
    event RequestFailed(uint256 indexed requestId);
    event UserEvent(uint256 indexed applicationId, uint256 indexed requestId, bytes encryptedData);
    event StateRootUpdate(uint256 indexed applicationId, uint256 indexed requestId, bytes oldStateRoot, bytes newStateRoot);
    //errors
    error AddressCantBeZero();
    error InvalidValue();
    error InvalidProtocolVersion();
    error InvalidApplicationId();
    error InvalidRequestId();
    error RequestIsAlreadyCompletedOrFailed(Structs.RequestStatus currentStatus);
    error InvalidStateRoot();
    error InvalidSignature();
    error InsufficientBalance();

    modifier onlyPostedRequest(uint256 requestId) {
        if(requestId >= requests.length) revert InvalidRequestId();
        if(requests[requestId].status != Structs.RequestStatus.POSTED) revert RequestIsAlreadyCompletedOrFailed(requests[requestId].status);
        _;
    }
    modifier validProtocolVersion(uint8 protocolVersion) {
        if(protocolVersion != PROTOCOL_VERSION) revert InvalidProtocolVersion();
        _;
    }
    modifier validApplicationId(uint256 applicationId) {
        if(applicationId != APPLICATION_ID) revert InvalidApplicationId();
        _;
    }

    //constructor
    constructor(ITeeAuthenticator _teeAuthenticator, address updateStatusOperator) {
        if(_teeAuthenticator == ITeeAuthenticator(address(0)) || updateStatusOperator == address(0)) revert AddressCantBeZero();

        teeAuthenticator = _teeAuthenticator;
        _grantRole(UPDATE_STATUS_ROLE, updateStatusOperator);
    }

    //request management functions
    function submitRequest(
        uint8 protocolVersion, 
        uint256 applicationId, 
        Structs.RequestType requestType, 
        bytes calldata payload, 
        uint256 value
    ) validProtocolVersion(protocolVersion) validApplicationId(applicationId) payable public returns(uint256) {
        //check value
        if(msg.value != value) revert InvalidValue(); //'value' is redundant now, but it will be needed when using ERC20
        //create request
        uint256 requestId = requests.length;
        requests.push(
            Structs.PendingRequest(
                protocolVersion,
                applicationId,
                requestType,
                requestId,
                payload,
                block.timestamp,
                msg.sender,
                Structs.RequestStatus.POSTED,
                value
            )
        );

        //add to queue
        idsQueue.add(requestId);

        //emit event
        emit RequestSubmitted(requestId, msg.sender);

        return requestId;
    }

    function markRequestCompleted(uint256 requestId) public onlyRole(UPDATE_STATUS_ROLE) onlyPostedRequest(requestId) {
        requests[requestId].status = Structs.RequestStatus.COMPLETED;
        //remove from queue
        idsQueue.remove(requestId);
        //emit event
        emit RequestCompleted(requestId);
    }

    function markRequestFailed(uint256 requestId) public onlyRole(UPDATE_STATUS_ROLE) onlyPostedRequest(requestId) nonReentrant {
        //remove from queue
        idsQueue.remove(requestId);
        //emit event
        emit RequestFailed(requestId);
        //refunds
        (bool refunded, ) = payable(requests[requestId].sender).call{value: requests[requestId].value, gas: 2300}("");
        if(refunded) requests[requestId].status = Structs.RequestStatus.FAILED_REFUNDED;
        else requests[requestId].status = Structs.RequestStatus.FAILED_NOT_REFUNDED;
    }

    function getPendingRequestsSize() public view returns(uint256) {
        return idsQueue.length();
    }

    function getPendingRequests() public view returns(Structs.PendingRequest[] memory) {
        uint256 setSize = idsQueue.length();
        uint256[] memory ids = new uint256[](setSize);

        uint256 i;
        while(i < setSize) {
            ids[i] = idsQueue.at(i);
            unchecked { ++i;}
        }

        //set sorting is not guaranteed, so we use this to order the requestsId
        ids = ids.sort();
        //and then get the corresponding pending requests
        Structs.PendingRequest[] memory res = new Structs.PendingRequest[](setSize);
        i = 0;
        while(i < setSize) {
            res[i] = requests[ids[i]];
            unchecked { ++i; }
        }

        return res;
    }

    //update status
    function stateUpdate(
        uint256 applicationId, 
        bytes calldata prevStateRoot, 
        bytes calldata newStateRoot, 
        uint256 processedRequestId,
        bytes[] memory events, 
        Structs.WithdrawalRequest[] memory withdrawalRequests, 
        bytes memory signature
    ) public nonReentrant validApplicationId(applicationId) onlyRole(UPDATE_STATUS_ROLE) {

        //check prev state root
        if(!_eq(stateRoot, "") && !_eq(prevStateRoot, stateRoot)) revert InvalidStateRoot();
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
        markRequestCompleted(processedRequestId);

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

    function _eq(bytes memory a, bytes memory b) internal pure returns(bool) {
        if(a.length != b.length) return false;
        uint256 i;
        while(i < a.length) {
            if(a[i] != b[i]) return false;
            unchecked {++i;}
        }
        return true;
    }
}
