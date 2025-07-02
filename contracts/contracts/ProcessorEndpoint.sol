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
    
    //state variables
    bytes public stateRoot;
    
    Structs.PendingRequest[] public requests;
    EnumerableSet.UintSet private idsQueue;

    ITeeAuthenticator public teeAuthenticator;
    //events
    event Withdrawal(address to, uint256 amount);
    event RequestSubmitted(uint256 requestId, address sender);
    event RequestCompleted(uint256 requestId);
    event RequestFailed(uint256 requestId);
    event UserEvent(uint8 applicationId, bytes encryptedData);
    event StateRootUpdate(uint8 applicationId, bytes oldStateRoot, bytes newStateRoot);
    //errors
    error AddressCantBeZero();
    error InvalidValue();
    error InvalidRequestId();
    error RequestIsAlreadyCompletedOrFailed(Structs.RequestStatus currentStatus);
    error InvalidStateRoot();
    error InvalidSignature();
    error InsufficientBalance();

    //modifiers
    modifier onlyValidRequest(uint256 requestId) {
        if(requestId >= requests.length) revert InvalidRequestId();
        _;
    }
    modifier onlyPostedRequest(uint256 requestId) {
        if(requestId >= requests.length) revert InvalidRequestId();
        if(requests[requestId].status != Structs.RequestStatus.POSTED) revert RequestIsAlreadyCompletedOrFailed(requests[requestId].status);
        _;
    }   

    //constructor
    constructor(ITeeAuthenticator _teeAuthenticator, address updateStatusOperator) payable {
        if(_teeAuthenticator == ITeeAuthenticator(address(0)) || updateStatusOperator == address(0)) revert AddressCantBeZero();

        teeAuthenticator = _teeAuthenticator;
        _grantRole(UPDATE_STATUS_ROLE, updateStatusOperator);
    }

    //request management functions
    function submitRequest(uint8 protocolVersion, uint8 applicationId, Structs.RequestType requestType, bytes calldata payload, uint256 value) payable public {
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
                Structs.RequestStatus.POSTED
            )
        );

        //add to queue
        idsQueue.add(requestId);

        //emit event
        emit RequestSubmitted(requestId, msg.sender);
    }

    function markRequestCompleted(uint256 requestId) public onlyRole(UPDATE_STATUS_ROLE) onlyPostedRequest(requestId){
        requests[requestId].status = Structs.RequestStatus.COMPLETED;
        //remove from queue
        idsQueue.remove(requestId);
        //emit event
        emit RequestCompleted(requestId);
    }

    function markRequestFailed(uint256 requestId) public onlyRole(UPDATE_STATUS_ROLE) onlyPostedRequest(requestId) {
        requests[requestId].status = Structs.RequestStatus.FAILED;
        //remove from queue
        idsQueue.remove(requestId);
        //emit event
        emit RequestFailed(requestId);
    }

    function getPendingRequestsSize() public view returns(uint256) {
        return idsQueue.length();
    }

    function getPendingRequests(uint256 offset, uint256 size) public view returns(Structs.PendingRequest[] memory) {
        uint256[] memory ids = new uint256[](size);

        uint256 setSize = idsQueue.length();
        uint256 i;
        while(i < size && offset < setSize) {
            ids[i] = idsQueue.at(offset);
            unchecked { ++i; ++offset;}
        }

        //set sorting is not guaranteed, so we use this to order the requestsId
        ids.sort();
        //and then get the corresponding pending requests
        Structs.PendingRequest[] memory res = new Structs.PendingRequest[](size);
        i= 0;
        while(i < size) {
            res[i] = requests[ids[i]];
            unchecked { ++i; }
        }

        return res;
    }

    //update status
    function stateUpdate(
        uint8 applicationId, 
        bytes calldata prevStateRoot, 
        bytes calldata newStateRoot, 
        bytes[] memory events, 
        Structs.WithdrawalRequest[] memory withdrawalRequests, 
        bytes memory signature
    ) public nonReentrant {  //this could be public since signature is checked

        //check prev state root
        if(!_eq(stateRoot, "") && !_eq(prevStateRoot, stateRoot)) revert InvalidStateRoot();
        //check signature
        if(!teeAuthenticator.checkSignature(applicationId, prevStateRoot, newStateRoot, events, withdrawalRequests, signature)) revert InvalidSignature();

        //check withdrawal sums 
        uint256 i;
        uint256 sum;
        while(i < withdrawalRequests.length) {
            sum += withdrawalRequests[i].amount;
            unchecked {++i;}
        }
        if(sum > address(this).balance) revert InsufficientBalance();
        
        //emit encrypted event
        i = 0;
        while(i < events.length) {
            emit UserEvent(applicationId, events[i]);
            unchecked {++i;}
        }

        //update state root
        stateRoot = newStateRoot;
        emit StateRootUpdate(applicationId, prevStateRoot, newStateRoot);

        //execute withdrawals (as last operation)
        i = 0;
        while(i < withdrawalRequests.length) {
            withdrawalRequests[i].receiver.transfer(withdrawalRequests[i].amount);
            emit Withdrawal(withdrawalRequests[i].receiver, withdrawalRequests[i].amount);
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
