// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts/utils/ReentrancyGuard.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IProcessorEndpoint.sol';
import './interfaces/IAuthorityRegistry.sol';
import './Structs.sol';

/// @title ProcessorEndpoint
/// @notice Implementation of the processor endpoint interface.
contract ProcessorEndpoint is AccessControl, IProcessorEndpoint, ReentrancyGuard {
  //constants
  bytes32 public constant UPDATE_STATUS_ROLE = keccak256('UPDATE_STATUS_ROLE');
  bytes32 public constant ADMIN = keccak256('ADMIN');
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
  IAuthorityRegistry public authorityRegistry;

  // Pull payment pattern state
  mapping(address => uint256) public payments;
  uint256 private _totalDeposits;

  uint256 public minFeePerRequest;
  address payable public feeCollector;

  modifier validProtocolVersion(uint8 protocolVersion) {
    if (protocolVersion != PROTOCOL_VERSION) revert InvalidProtocolVersion();
    _;
  }

  modifier validApplicationId(uint64 applicationId) {
    if (applicationId != APPLICATION_ID) revert InvalidApplicationId();
    _;
  }

  /// @param _teeAuthenticator Contract used to verify update signatures.
  /// @param _authorityRegistry Registry for authority checks.
  /// @param updateStatusOperator Initial operator for status updates.
  /// @param admin Initial admin address.
  /// @param _minFeePerRequest Minimum fee enforced per request.
  constructor(
    ITeeAuthenticator _teeAuthenticator,
    IAuthorityRegistry _authorityRegistry,
    address updateStatusOperator,
    address admin,
    uint256 _minFeePerRequest
  ) {
    if (
      _teeAuthenticator == ITeeAuthenticator(address(0)) ||
      address(_authorityRegistry) == address(0) ||
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

  /// @inheritdoc IProcessorEndpoint
  function submitRequest(
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    uint256 depositAmount, // part of the sent value forwarded to the application, for app logic
    uint256 maxFeeValue // part ot the sent value reserved for fee payment
  )
    external
    payable
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    nonReentrant
    returns (bytes32)
  {
    //check values
    if (msg.value != depositAmount + maxFeeValue) revert InvalidValue();
    if (maxFeeValue < minFeePerRequest) revert FeeValueBelowMinimum();

    //check queue size
    if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

    if (requestType == Structs.RequestType.ASSOCIATEKEY) {
      //if requestype is associatekey, the payload must be 133 bytes long (contains a Secp521r1_PubKey)
      if (payload.length != 133) revert InvalidPayload();
    } else if (requestType == Structs.RequestType.DEANONYMIZATION) {
      // only allowed authorities can request deanonymization
      if (!authorityRegistry.checkAuthorityIsAllowed(applicationId, msg.sender)) {
        revert AuthorityNotAllowed();
      }
    }

    //create request
    bytes32 requestId = generateRequestId(
      msg.sender,
      applicationId,
      requestType,
      payload,
      depositAmount,
      _tail
    );
    requestById[requestId] = Structs.PendingRequest({
      timestamp: block.timestamp,
      depositAmount: depositAmount,
      maxFeeValue: maxFeeValue,
      requestId: requestId,
      payload: payload,
      sender: msg.sender,
      applicationId: applicationId,
      protocolVersion: protocolVersion,
      requestType: requestType
    });
    _requestIdByOrder[_tail] = requestId;

    unchecked {
      ++_tail;
    }

    //emit event
    emit RequestSubmitted(requestId, msg.sender);

    return requestId;
  }

  function _removeRequest() private {
    delete requestById[_requestIdByOrder[_head]];
    delete _requestIdByOrder[_head];
    unchecked {
      ++_head;
    }
  }

  function _markRequestCompleted(bytes32 requestId, uint256 applicationFees) private {
    _removeRequest();

    emit RequestCompleted(
      requestId,
      applicationFees,
      Structs.RequestResult.COMPLETED,
      Structs.ErrorCode.NO_ERROR,
      ''
    );
  }

  // We refund the maxValueFee - minFeePerRequest (to be changed in the future)
  /// @inheritdoc IProcessorEndpoint
  function markRequestFailed(
    bytes32 requestId,
    Structs.ErrorCode errorCode,
    string calldata errorMessage
  ) external onlyRole(UPDATE_STATUS_ROLE) {
    if (!isCurrentPendingRequest(requestId)) revert InvalidRequestId();

    Structs.PendingRequest memory requestInfo = requestById[requestId];
    uint256 minFee = minFeePerRequest;

    _removeRequest();

    //credit refund to sender's pending balance (pull pattern)
    uint256 refund = requestInfo.depositAmount + (requestInfo.maxFeeValue - minFee);
    if (refund > 0) {
      address payable sender = payable(requestInfo.sender);
      _asyncTransfer(sender, refund);
      emit Refund(requestInfo.applicationId, requestId, sender, refund);
    }

    //credit minimum fee to feeCollector's pending balance
    _asyncTransfer(feeCollector, minFeePerRequest);

    emit RequestCompleted(
      requestId,
      minFeePerRequest,
      Structs.RequestResult.FAILED,
      errorCode,
      errorMessage
    );
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsSize() public view returns (uint256) {
    if (_tail > _head) {
      return (_tail - _head);
    } else {
      return 0;
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequests() external view returns (Structs.PendingRequest[] memory) {
    uint256 numOfPendingRequests = getPendingRequestsSize();

    Structs.PendingRequest[] memory res = new Structs.PendingRequest[](numOfPendingRequests);
    uint256 i = _head;
    uint256 tail = _tail;
    uint256 j;
    while (i < tail) {
      bytes32 requestId = _requestIdByOrder[i];
      res[j] = requestById[requestId];
      unchecked {
        ++i;
        ++j;
      }
    }

    return res;
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsPage(
    uint256 offset,
    uint256 limit
  ) external view returns (Structs.PendingRequest[] memory) {
    uint256 size = getPendingRequestsSize();
    if (offset >= size || limit == 0) {
      return new Structs.PendingRequest[](0);
    }

    uint256 end = offset + limit;
    if (end > size) {
      end = size;
    }
    uint256 count = end - offset;

    Structs.PendingRequest[] memory res = new Structs.PendingRequest[](count);
    uint256 i = _head + offset;
    uint256 stop = _head + end;
    uint256 j;
    while (i < stop) {
      bytes32 requestId = _requestIdByOrder[i];
      res[j] = requestById[requestId];
      unchecked {
        ++i;
        ++j;
      }
    }

    return res;
  }

  /// @inheritdoc IProcessorEndpoint
  function stateUpdate(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    bytes[] calldata events,
    string[] calldata eventSubTypes,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    bytes calldata signature
  ) external validApplicationId(applicationId) onlyRole(UPDATE_STATUS_ROLE) {
    //check prev state root
    if (stateRoot != bytes32(0) && prevStateRoot != stateRoot) revert InvalidStateRoot();
    //check valid request
    if (!isCurrentPendingRequest(processedRequestId)) revert InvalidRequestId();

    //check signature
    uint256 eventsLength = events.length;
    uint256 eventSubTypesLength = eventSubTypes.length;
    if (eventsLength != eventSubTypesLength) revert InvalidPayload();
    if (
      !teeAuthenticator.checkSignature(
        applicationId,
        prevStateRoot,
        newStateRoot,
        processedRequestId,
        events,
        eventSubTypes,
        withdrawalRequests,
        refund,
        applicationFees,
        signature
      )
    ) revert InvalidSignature();

    //check values
    Structs.PendingRequest storage requestInfo = requestById[processedRequestId];
    uint256 maxFeeValue = requestInfo.maxFeeValue;
    Structs.RequestType reqType = requestInfo.requestType;
    address payable sender = payable(requestInfo.sender);

    if (refund + applicationFees != maxFeeValue) revert InvalidValue();
    if (applicationFees < minFeePerRequest) {
      revert InvalidValue();
    }

    //check withdrawal sums (account for already committed pending deposits)
    uint256 i;
    uint256 sum;
    uint256 withdrawalsLength = withdrawalRequests.length;
    while (i < withdrawalsLength) {
      sum += withdrawalRequests[i].amount;
      unchecked {
        ++i;
      }
    }
    sum += refund + applicationFees;

    if (sum > address(this).balance - _totalDeposits) revert InsufficientBalance();

    //set requests as completed
    _markRequestCompleted(processedRequestId, applicationFees);

    //emit encrypted event
    i = 0;
    while (i < events.length) {
      emit UserEvent(applicationId, processedRequestId, eventSubTypes[i], events[i]);
      unchecked {
        ++i;
      }
    }

    if (reqType == Structs.RequestType.DEANONYMIZATION) {
      //a completed DEANONYMIZATION request must have always generated a report
      emit ReportGenerated(applicationId, processedRequestId);
    }

    //update state root and request
    stateRoot = newStateRoot;
    emit StateRootUpdate(applicationId, processedRequestId, prevStateRoot, newStateRoot);

    //credit refund to sender's pending balance (pull pattern)
    if (refund > 0) {
      _asyncTransfer(sender, refund);
      emit Refund(applicationId, processedRequestId, sender, refund);
    }

    //credit fee to feeCollector's pending balance
    _asyncTransfer(feeCollector, applicationFees);

    //credit withdrawals to receivers' pending balances
    i = 0;
    while (i < withdrawalRequests.length) {
      _asyncTransfer(withdrawalRequests[i].receiver, withdrawalRequests[i].amount);
      emit Withdrawal(
        applicationId,
        processedRequestId,
        withdrawalRequests[i].receiver,
        withdrawalRequests[i].amount
      );
      unchecked {
        ++i;
      }
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function updateQueueThreshold(uint256 newThreshold) external onlyRole(ADMIN) {
    if (newThreshold == 0) revert InvalidValue();
    maxQueueSize = newThreshold;
    emit QueueThresholdUpdated(newThreshold);
  }

  /// @inheritdoc IProcessorEndpoint
  function updateFeeCollector(address payable newFeeCollector) external onlyRole(ADMIN) {
    if (newFeeCollector == address(0)) revert AddressCantBeZero();
    feeCollector = newFeeCollector;
    emit FeeCollectorUpdated(newFeeCollector);
  }

  /// @inheritdoc IProcessorEndpoint
  function getNextPendingRequest()
    external
    view
    returns (Structs.PendingRequest memory, bytes32, bool success)
  {
    uint256 numOfRequests = getPendingRequestsSize();
    if (numOfRequests > 0) {
      bytes32 requestId = _requestIdByOrder[_head];
      return (requestById[requestId], stateRoot, true);
    }

    Structs.PendingRequest memory emptyReq;
    return (emptyReq, stateRoot, false);
  }

  /// @inheritdoc IProcessorEndpoint
  function isCurrentPendingRequest(bytes32 requestId) public view returns (bool) {
    return getPendingRequestsSize() > 0 && _requestIdByOrder[_head] == requestId;
  }

  // Pull payment pattern functions
  function _asyncTransfer(address dest, uint256 amount) internal {
    payments[dest] += amount;
    _totalDeposits += amount;
  }

  /// @inheritdoc IProcessorEndpoint
  function withdrawPayments(address payable payee) public nonReentrant {
    uint256 payment = payments[payee];
    if (payment == 0) return;

    payments[payee] = 0;
    _totalDeposits -= payment;

    emit PaymentWithdrawn(payee, payment);

    (bool success, ) = payee.call{value: payment}('');
    if (!success) revert TransferFailed();
  }

  /// @inheritdoc IProcessorEndpoint
  function generateRequestId(
    address sender,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    uint256 depositAmount,
    uint256 idx
  ) public pure returns (bytes32) {
    bytes32 requestId = keccak256(
      abi.encodePacked(sender, applicationId, requestType, payload, depositAmount, idx)
    );

    return requestId;
  }
}
