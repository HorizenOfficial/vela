// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts/utils/ReentrancyGuard.sol';
import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import '@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';
import '@openzeppelin/contracts/utils/cryptography/EIP712.sol';
import '@openzeppelin/contracts/utils/Strings.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IProcessorEndpoint.sol';
import './interfaces/IAuthorityRegistry.sol';
import './interfaces/ITokenAllowlist.sol';
import './Structs.sol';
import './interfaces/ITrigger.sol';

/// @title ProcessorEndpoint
/// @notice Implementation of the processor endpoint interface.
contract ProcessorEndpoint is AccessControl, IProcessorEndpoint, ReentrancyGuard, EIP712 {
  using SafeERC20 for IERC20;

  struct RequestQueue {
    mapping(bytes32 => Structs.PendingRequest) requestById;
    mapping(uint256 => bytes32) idByOrder;
    uint256 head;
    uint256 tail;
  }

  //constants
  bytes32 public constant UPDATE_STATUS_ROLE = keccak256('UPDATE_STATUS_ROLE');
  bytes32 public constant ADMIN = keccak256('ADMIN');
  bytes32 public constant DEPLOYER_ROLE = keccak256('DEPLOYER_ROLE');
  bytes32 public constant RESET_OPERATOR = keccak256('RESET_OPERATOR');
  uint8 public constant PROTOCOL_VERSION = 0;
  //state variables
  mapping(uint64 => bytes32) public applicationStateRoots;
  uint64[] private _deployedAppIds;
  uint256 public maxNumOfApplications = 10;
  uint256 public availableDeploySlots = maxNumOfApplications;

  RequestQueue private _requestQueue;
  uint256 public maxQueueSize = 10;

  ITeeAuthenticator public teeAuthenticator;
  IAuthorityRegistry public authorityRegistry;
  ITokenAllowlist public tokenAllowlist;

  // Pull payment pattern state — per-token, per-payee
  mapping(address => mapping(address => uint256)) public pendingClaims;
  mapping(address => uint256) public totalPendingClaims;

  // Per-app, per-token custody tracking for solvency isolation.
  // Credited on submitRequest (assetAmount only; fees are tracked globally),
  // debited on stateUpdate: success path (withdrawals), error path (assetAmount).
  // Fees are self-balancing per request (refund + applicationFees == maxFeeValue)
  // so global balance checks are sufficient for the fee portion.
  // If an app's withdrawals are less than its deposits,
  // the residual accumulates here as credit available to future requests.
  // Note: There is currently no mechanism to recover residual funds from decommissioned
  // apps.
  mapping(uint64 => mapping(address => uint256)) public appCustody;
  mapping(address => uint256) public totalAppCustody;

  uint256 public minFeePerRequest;
  address payable public feeCollector;

  // EIP-712 typehash for facilitator request authorization
  bytes32 public constant REQUEST_AUTHORIZATION_TYPEHASH =
    keccak256(
      'RequestAuthorization(address sender,uint8 protocolVersion,uint64 applicationId,uint8 requestType,bytes32 payloadHash,address tokenAddress,uint256 assetAmount,uint256 nonce,uint256 deadline)'
    );

  // Sequential nonces per user for facilitator replay protection
  mapping(address => uint256) public facilitatorNonces;
  // Trigger contracts associated to each applicationId
  mapping(uint64 => ITrigger) public triggerContracts;
  // Reverse mapping for the above to check if a trigger is valid when adding to the queue
  mapping(address => uint64) public triggersToAppIds;
  // FIFO queue populated by trigger contracts; served before the normal queue
  RequestQueue private _triggerQueue;

  modifier validProtocolVersion(uint8 protocolVersion) {
    if (protocolVersion != PROTOCOL_VERSION) revert InvalidProtocolVersion();
    _;
  }

  modifier validApplicationId(uint64 applicationId) {
    if (applicationStateRoots[applicationId] == bytes32(0)) revert InvalidApplicationId();
    _;
  }

  /// @param _teeAuthenticator Contract used to verify update signatures.
  /// @param _authorityRegistry Registry for authority checks.
  /// @param updateStatusOperator Initial operator for status updates.
  /// @param admin Initial admin address.
  /// @param resetOperator Address granted RESET_OPERATOR role. Pass address(0) to disable reset
  ///        permanently (required for production). The role cannot be granted after deployment.
  /// @param _minFeePerRequest Minimum fee enforced per request.
  /// @param _tokenAllowlist External token allowlist contract.
  constructor(
    ITeeAuthenticator _teeAuthenticator,
    IAuthorityRegistry _authorityRegistry,
    address updateStatusOperator,
    address admin,
    address resetOperator,
    uint256 _minFeePerRequest,
    ITokenAllowlist _tokenAllowlist
  ) EIP712('Vela', Strings.toString(PROTOCOL_VERSION)) {
    if (
      address(_teeAuthenticator) == address(0) ||
      address(_authorityRegistry) == address(0) ||
      updateStatusOperator == address(0) ||
      admin == address(0) ||
      address(_tokenAllowlist) == address(0)
    ) revert AddressCantBeZero();

    teeAuthenticator = _teeAuthenticator;
    authorityRegistry = _authorityRegistry;
    tokenAllowlist = _tokenAllowlist;
    feeCollector = payable(updateStatusOperator);
    _grantRole(UPDATE_STATUS_ROLE, updateStatusOperator);
    _grantRole(ADMIN, admin);
    _setRoleAdmin(DEPLOYER_ROLE, ADMIN);
    _grantRole(DEPLOYER_ROLE, admin);
    minFeePerRequest = _minFeePerRequest;
    if (resetOperator != address(0)) {
      _grantRole(RESET_OPERATOR, resetOperator);
    }
  }

  // @notice receive ETH (sent back by trigger contracts)
  receive() external payable {}

  /// @param _teeAuthenticator Contract used to verify update signatures.

  /// @inheritdoc IProcessorEndpoint
  function submitRequest(
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 maxFeeValue // part of the sent value reserved for fee payment
  )
    external
    payable
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    nonReentrant
    returns (bytes32)
  {
    if (requestType == Structs.RequestType.DEPLOYAPP) revert InvalidRequestType();

    //check values
    if (maxFeeValue < minFeePerRequest) revert FeeValueBelowMinimum();
    //check queue size
    if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

    if (tokenAddress == ETH_TOKEN) {
      if (msg.value != assetAmount + maxFeeValue) revert InvalidValue();
    } else {
      if (msg.value != maxFeeValue) revert InvalidValue();
      if (assetAmount == 0) revert InvalidValue();
      if (!tokenAllowlist.isAllowedToken(tokenAddress)) revert ITokenAllowlist.TokenNotAllowed();
      // Pull ERC-20 tokens with balance-before/after check
      _pullERC20(tokenAddress, msg.sender, assetAmount);
    }

    if (requestType == Structs.RequestType.ASSOCIATEKEY) {
      //if requestype is associatekey, the payload must be 133 bytes (key only) or 226 bytes (key + encrypted seed)
      if (payload.length != 133 && payload.length != 226) revert InvalidPayload();
    } else if (requestType == Structs.RequestType.DEANONYMIZATION) {
      // only allowed authorities can request deanonymization
      if (!authorityRegistry.checkAuthorityIsAllowed(applicationId, msg.sender)) {
        revert AuthorityNotAllowed();
      }
    }

    appCustody[applicationId][tokenAddress] += assetAmount;
    totalAppCustody[tokenAddress] += assetAmount;

    //create request and enqueue
    return
      _enqueueRequest(
        msg.sender,
        address(0),
        protocolVersion,
        applicationId,
        requestType,
        payload,
        tokenAddress,
        assetAmount,
        maxFeeValue
      );
  }

  /// @inheritdoc IProcessorEndpoint
  function submitRequestFor(
    address sender,
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 deadline,
    bytes calldata requestSignature,
    bytes calldata depositPermit
  )
    external
    payable
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    nonReentrant
    returns (bytes32)
  {
    // 1. Only ASSOCIATEKEY, PROCESS and PLAINPROCESS are supported
    if (
      requestType != Structs.RequestType.ASSOCIATEKEY &&
      requestType != Structs.RequestType.PROCESS &&
      requestType != Structs.RequestType.PLAINPROCESS
    ) revert InvalidRequestType();

    // 2. Verify deadline not expired
    if (block.timestamp > deadline) revert DeadlineExpired();

    // 3. Check queue size
    if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

    // 4. Validate payload for ASSOCIATEKEY
    if (requestType == Structs.RequestType.ASSOCIATEKEY) {
      if (payload.length != 133 && payload.length != 226) revert InvalidPayload();
    }

    // 5. Read current nonce and build EIP-712 hash
    uint256 nonce = facilitatorNonces[sender];
    bytes32 structHash = keccak256(
      abi.encode(
        REQUEST_AUTHORIZATION_TYPEHASH,
        sender,
        protocolVersion,
        applicationId,
        uint8(requestType),
        keccak256(payload),
        tokenAddress,
        assetAmount,
        nonce,
        deadline
      )
    );
    bytes32 digest = _hashTypedDataV4(structHash);

    // 6. Recover user address from EIP-712 request signature and verify
    address recoveredSigner = ECDSA.recover(digest, requestSignature);
    if (recoveredSigner == address(0)) revert InvalidSignature();
    if (recoveredSigner != sender) revert InvalidSigner();

    // 7. Consume nonce (replay protection)
    facilitatorNonces[sender] = nonce + 1;

    // 8. Validate fee
    uint256 maxFeeValue = msg.value;
    if (maxFeeValue < minFeePerRequest) revert FeeValueBelowMinimum();

    // 9. Validate token and handle deposit
    if (assetAmount == 0 && tokenAddress != ETH_TOKEN) revert InvalidValue();
    if (assetAmount > 0) {
      if (tokenAddress == ETH_TOKEN) revert InvalidValue();
      if (!tokenAllowlist.isAllowedToken(tokenAddress)) revert ITokenAllowlist.TokenNotAllowed();

      // Decode deposit permit and execute EIP-2612 permit + transferFrom
      if (depositPermit.length != 96) revert InvalidPermit();
      (uint8 v, bytes32 r, bytes32 s) = abi.decode(depositPermit, (uint8, bytes32, bytes32));

      // Check current allowance before calling permit
      IERC20 token = IERC20(tokenAddress);
      if (token.allowance(sender, address(this)) < assetAmount) {
        IERC20Permit(tokenAddress).permit(sender, address(this), assetAmount, deadline, v, r, s);
      }

      _pullERC20(tokenAddress, sender, assetAmount);

      appCustody[applicationId][tokenAddress] += assetAmount;
      totalAppCustody[tokenAddress] += assetAmount;
    }

    // 10. Create PendingRequest with sender = user (not msg.sender) and enqueue
    return
      _enqueueRequest(
        sender,
        msg.sender,
        protocolVersion,
        applicationId,
        requestType,
        payload,
        tokenAddress,
        assetAmount,
        maxFeeValue
      );
  }

  /// @inheritdoc IProcessorEndpoint
  function getFacilitatorNonce(address user) external view returns (uint256) {
    return facilitatorNonces[user];
  }

  /// @inheritdoc IProcessorEndpoint
  function submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  ) external payable validProtocolVersion(protocolVersion) nonReentrant returns (bytes32) {
    if (!hasRole(DEPLOYER_ROLE, msg.sender)) revert DeployerNotAllowed();
    if (availableDeploySlots == 0) revert MaxNumOfApplicationsExceeded();
    //check queue size
    if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();
    if (msg.value < minFeePerRequest) revert FeeValueBelowMinimum();

    --availableDeploySlots;

    Structs.RequestType requestType = Structs.RequestType.DEPLOYAPP;
    //create request
    bytes32 requestId = generateRequestId(
      msg.sender,
      0, // deploy requests have applicationId 0, a unique applicationId will be derived from the requestId for each deploy request to avoid collisions with regular requests and to group deploy requests together
      requestType,
      payload,
      ETH_TOKEN,
      0,
      _requestQueue.tail
    );

    uint64 applicationId = uint64(bytes8(requestId)); // Derive a unique application ID from the request ID for deploy requests
    _queueEnqueue(
      _requestQueue,
      requestId,
      Structs.PendingRequest({
        timestamp: block.timestamp,
        tokenAddress: ETH_TOKEN,
        assetAmount: 0,
        maxFeeValue: msg.value,
        requestId: requestId,
        payload: payload,
        sender: msg.sender,
        facilitator: address(0),
        applicationId: applicationId,
        protocolVersion: protocolVersion,
        requestType: requestType
      })
    );

    //emit event
    emit DeployRequestSubmitted(applicationId, requestId, msg.sender);

    return requestId;
  }

  function _pullERC20(address tokenAddress, address from, uint256 amount) internal {
    IERC20 token = IERC20(tokenAddress);
    uint256 balanceBefore = token.balanceOf(address(this));
    token.safeTransferFrom(from, address(this), amount);
    uint256 received = token.balanceOf(address(this)) - balanceBefore;
    if (received != amount) revert TransferAmountMismatch();
  }

  function _enqueueRequest(
    address sender,
    address facilitator,
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 maxFeeValue
  ) internal returns (bytes32) {
    bytes32 requestId = generateRequestId(
      sender,
      applicationId,
      requestType,
      payload,
      tokenAddress,
      assetAmount,
      _requestQueue.tail
    );
    _queueEnqueue(
      _requestQueue,
      requestId,
      Structs.PendingRequest({
        timestamp: block.timestamp,
        tokenAddress: tokenAddress,
        assetAmount: assetAmount,
        maxFeeValue: maxFeeValue,
        requestId: requestId,
        payload: payload,
        sender: sender,
        facilitator: facilitator,
        applicationId: applicationId,
        protocolVersion: protocolVersion,
        requestType: requestType
      })
    );

    emit RequestSubmitted(applicationId, requestId, sender, facilitator);
    return requestId;
  }

  function _removeRequest(bytes32 requestId) private {
    if (_queueSize(_triggerQueue) > 0 && _queueIsHead(_triggerQueue, requestId)) {
      _queueDequeueHead(_triggerQueue);
    } else {
      _queueDequeueHead(_requestQueue);
    }
  }

  function _markRequestCompleted(
    uint64 applicationId,
    bytes32 requestId,
    uint256 applicationFees,
    Structs.RequestResult result,
    Structs.ErrorCode errCode,
    string memory errorMsg,
    Structs.RequestType requestType
  ) private {
    _removeRequest(requestId);

    if (requestType == Structs.RequestType.DEPLOYAPP) {
      emit DeployRequestCompleted(
        applicationId,
        requestId,
        applicationFees,
        result,
        errCode,
        errorMsg
      );
    } else {
      emit RequestCompleted(applicationId, requestId, applicationFees, result, errCode, errorMsg);
    }

    _asyncTransfer(ETH_TOKEN, feeCollector, applicationFees);
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsSize() public view returns (uint256) {
    return _queueSize(_requestQueue);
  }

  /// @inheritdoc IProcessorEndpoint
  function getTriggerQueueSize() public view returns (uint256) {
    return _queueSize(_triggerQueue);
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequests() external view returns (Structs.PendingRequest[] memory) {
    return _queueGetAll(_requestQueue);
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsPage(
    uint256 offset,
    uint256 limit
  ) external view returns (Structs.PendingRequest[] memory) {
    return _queueGetPage(_requestQueue, offset, limit);
  }

  /// @notice Returns the stored request for a given id (normal queue only).
  function requestById(bytes32 id) external view returns (Structs.PendingRequest memory) {
    return _requestQueue.requestById[id];
  }

  //update status
  /// @inheritdoc IProcessorEndpoint
  function stateUpdate(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg,
    bytes calldata signature
  ) external onlyRole(UPDATE_STATUS_ROLE) nonReentrant {
    //check valid request
    if (!isCurrentPendingRequest(processedRequestId)) revert InvalidRequestId();

    // Check application Id
    bool fromTriggerQueue = _queueIsHead(_triggerQueue, processedRequestId);
    Structs.PendingRequest storage requestInfo = fromTriggerQueue
      ? _triggerQueue.requestById[processedRequestId]
      : _requestQueue.requestById[processedRequestId];
    if (applicationId != requestInfo.applicationId) revert InvalidApplicationId();

    //check prev state root
    if (prevStateRoot != applicationStateRoots[applicationId]) revert InvalidStateRoot();

    uint256 eventsLength = userEventData.events.length;
    if (eventsLength != userEventData.subTypes.length) revert InvalidPayload();

    uint256 appEventsLength = appEventData.events.length;
    if (appEventsLength != appEventData.subTypes.length) revert InvalidPayload();

    //check signature
    Structs.SignatureParams memory sigParams = Structs.SignatureParams({
      applicationId: applicationId,
      prevStateRoot: prevStateRoot,
      newStateRoot: newStateRoot,
      processedRequestId: processedRequestId,
      userEvents: userEventData,
      appEvents: appEventData,
      withdrawalRequests: withdrawalRequests,
      refundAmount: refund,
      applicationFee: applicationFees,
      errorCode: errorCode,
      errorMsg: errorMsg
    });
    if (!teeAuthenticator.checkSignature(sigParams, signature)) revert InvalidSignature();

    //check values

    uint256 maxFeeValue = requestInfo.maxFeeValue;
    address payable sender = payable(requestInfo.sender);
    // Fee refund recipient: facilitator if present, otherwise sender
    address payable feeRecipient = requestInfo.facilitator != address(0)
      ? payable(requestInfo.facilitator)
      : sender;

    // Handle error case (signed error payload from TEE)
    if (errorCode != Structs.ErrorCode.NO_ERROR) {
      // For errors: state unchanged (prevStateRoot == newStateRoot), no events, no withdrawals
      // Refund user (minus minimum fee) and collect minimum fee
      if (eventsLength != 0 || appEventsLength != 0 || withdrawalRequests.length != 0)
        revert InvalidPayload();
      if (applicationStateRoots[applicationId] != newStateRoot) revert InvalidStateRoot();

      // Per-app per-token solvency check, then ETH balance check for fee outflow.
      uint256 assetAmount = requestInfo.assetAmount;
      address reqTokenAddress = requestInfo.tokenAddress;
      if (assetAmount > appCustody[applicationId][reqTokenAddress]) revert InsufficientAppBalance();
      if (requestInfo.maxFeeValue > _getAvailableEthBalance()) revert InsufficientBalance();
      appCustody[applicationId][reqTokenAddress] -= assetAmount;
      totalAppCustody[reqTokenAddress] -= assetAmount;

      // Refund business-asset deposit in its original token to the user.
      // Fee refund is always in ETH, routed to facilitator if present.
      uint256 feeRefund = requestInfo.maxFeeValue - minFeePerRequest;
      if (reqTokenAddress == ETH_TOKEN) {
        if (assetAmount > 0) {
          // ETH requests that moved also assets: only direct path, never facilitated => refund all to the sender
          uint256 totalRefund = assetAmount + feeRefund;
          _asyncTransfer(ETH_TOKEN, sender, totalRefund);
          emit Refund(applicationId, processedRequestId, sender, ETH_TOKEN, totalRefund);
        } else {
          // ETH requests with ETH used only for fee: fee refund in ETH to feeRecipient
          if (feeRefund > 0) {
            _asyncTransfer(ETH_TOKEN, feeRecipient, feeRefund);
            emit Refund(applicationId, processedRequestId, feeRecipient, ETH_TOKEN, feeRefund);
          }
        }
      } else {
        // For ERC-20 requests, asset refund in token to user, fee refund in ETH to feeRecipient
        if (assetAmount > 0) {
          _asyncTransfer(reqTokenAddress, sender, assetAmount);
          emit Refund(applicationId, processedRequestId, sender, reqTokenAddress, assetAmount);
        }
        if (feeRefund > 0) {
          _asyncTransfer(ETH_TOKEN, feeRecipient, feeRefund);
          emit Refund(applicationId, processedRequestId, feeRecipient, ETH_TOKEN, feeRefund);
        }
      }

      if (requestInfo.requestType == Structs.RequestType.DEPLOYAPP) {
        unchecked {
          ++availableDeploySlots;
        }
      }

      _markRequestCompleted(
        applicationId,
        processedRequestId,
        minFeePerRequest,
        Structs.RequestResult.FAILED,
        Structs.ErrorCode(errorCode),
        errorMsg,
        requestInfo.requestType
      );

      return;
    }

    // Handle success case
    // State cannot remain the same
    if (applicationStateRoots[applicationId] == newStateRoot) revert InvalidStateRoot();

    if (refund + applicationFees != maxFeeValue) revert InvalidValue();
    if (applicationFees < minFeePerRequest) {
      revert InvalidValue();
    }

    //check withdrawal sums and debit per-app per-token custody
    uint256 i;
    uint256 withdrawalsLength = withdrawalRequests.length;
    uint256 ethWithdrawalSum;
    // Accumulate per-token ERC-20 withdrawal sums for a single post-loop solvency check
    address[] memory erc20Tokens = new address[](withdrawalsLength);
    uint256[] memory erc20Sums = new uint256[](withdrawalsLength);
    uint256 erc20TokenCount;
    while (i < withdrawalsLength) {
      address wToken = withdrawalRequests[i].tokenAddress;
      uint256 wAmount = withdrawalRequests[i].amount;
      if (wAmount > appCustody[applicationId][wToken]) revert InsufficientAppBalance();
      unchecked {
        appCustody[applicationId][wToken] -= wAmount;
      }
      if (wAmount > totalAppCustody[wToken]) revert InsufficientBalance();
      unchecked {
        totalAppCustody[wToken] -= wAmount;
      }
      if (wToken == ETH_TOKEN) {
        ethWithdrawalSum += wAmount;
      } else {
        bool found;
        uint256 j;
        while (j < erc20TokenCount) {
          if (erc20Tokens[j] == wToken) {
            erc20Sums[j] += wAmount;
            found = true;
            break;
          }
          unchecked {
            ++j;
          }
        }
        if (!found) {
          erc20Tokens[erc20TokenCount] = wToken;
          erc20Sums[erc20TokenCount] = wAmount;
          unchecked {
            ++erc20TokenCount;
          }
        }
      }
      unchecked {
        ++i;
      }
    }

    // Post-loop ERC-20 solvency: one balanceOf call per unique token.
    // totalAppCustody is already decremented; totalPendingClaims not yet incremented,
    // so we add the withdrawal sum to account for the in-flight credits.
    i = 0;
    while (i < erc20TokenCount) {
      address token = erc20Tokens[i];
      if (
        IERC20(token).balanceOf(address(this)) <
        totalAppCustody[token] + totalPendingClaims[token] + erc20Sums[i]
      ) revert InsufficientBalance();
      unchecked {
        ++i;
      }
    }

    // ETH solvency: contract must hold enough ETH to cover fee outflow + ETH withdrawals
    uint256 totalEthOutflow = refund + applicationFees + ethWithdrawalSum;
    if (totalEthOutflow > _getAvailableEthBalance()) revert InsufficientBalance();

    //emit encrypted event
    i = 0;
    while (i < eventsLength) {
      emit UserEvent(
        applicationId,
        processedRequestId,
        userEventData.subTypes[i],
        userEventData.events[i]
      );
      unchecked {
        ++i;
      }
    }

    //emit app event
    i = 0;
    while (i < appEventsLength) {
      emit AppEvent(
        applicationId,
        processedRequestId,
        appEventData.subTypes[i],
        appEventData.events[i]
      );
      unchecked {
        ++i;
      }
    }

    Structs.RequestType reqType = requestInfo.requestType;
    if (reqType == Structs.RequestType.DEANONYMIZATION) {
      //a completed DEANONYMIZATION request must have always generated a report
      emit ReportGenerated(applicationId, processedRequestId);
    }

    //update state root and request
    applicationStateRoots[applicationId] = newStateRoot;
    if (reqType == Structs.RequestType.DEPLOYAPP) {
      _deployedAppIds.push(applicationId);
      // optional trigger: payload must be exactly 32 bytes (ABI-encoded address)
      if (requestInfo.payload.length == 32) {
        address triggerContract = abi.decode(requestInfo.payload, (address));
        if (triggerContract != address(0)) {
          triggerContracts[applicationId] = ITrigger(triggerContract);
          triggersToAppIds[triggerContract] = applicationId;
        }
      }
    }
    emit StateRootUpdate(applicationId, processedRequestId, prevStateRoot, newStateRoot);

    //credit refund to feeRecipient's pending balance (pull pattern) — refund is always ETH
    if (refund > 0) {
      _asyncTransfer(ETH_TOKEN, feeRecipient, refund);
      emit Refund(applicationId, processedRequestId, feeRecipient, ETH_TOKEN, refund);
    }

    //credit withdrawals to receivers' pending balances
    i = 0;
    while (i < withdrawalsLength) {
      _asyncTransfer(
        withdrawalRequests[i].tokenAddress,
        withdrawalRequests[i].receiver,
        withdrawalRequests[i].amount
      );
      emit Withdrawal(
        applicationId,
        processedRequestId,
        withdrawalRequests[i].receiver,
        withdrawalRequests[i].tokenAddress,
        withdrawalRequests[i].amount
      );
      unchecked {
        ++i;
      }
    }

    //invoke trigger contracts, if registered
    _invokeTrigger(
      applicationId,
      prevStateRoot,
      newStateRoot,
      processedRequestId,
      userEventData,
      appEventData,
      withdrawalRequests,
      refund,
      applicationFees,
      errorCode,
      errorMsg
    );

    //set requests as completed
    _markRequestCompleted(
      applicationId,
      processedRequestId,
      applicationFees,
      Structs.RequestResult.COMPLETED,
      Structs.ErrorCode.NO_ERROR,
      '',
      reqType
    );
  }

  /// @inheritdoc IProcessorEndpoint
  function updateQueueThreshold(uint256 newThreshold) external onlyRole(ADMIN) {
    if (newThreshold == 0) revert InvalidValue();
    maxQueueSize = newThreshold;
    emit QueueThresholdUpdated(newThreshold);
  }

  /// @inheritdoc IProcessorEndpoint
  function updateMaxNumOfApplications(uint256 newMax) external onlyRole(ADMIN) {
    if (newMax == 0) revert InvalidValue();
    uint256 deployedApps = maxNumOfApplications - availableDeploySlots;
    if (newMax < deployedApps) revert InvalidValue();
    uint256 oldMax = maxNumOfApplications;
    maxNumOfApplications = newMax;
    availableDeploySlots = newMax - deployedApps;
    emit MaxNumberOfAppUpdated(oldMax, newMax);
  }

  /// @inheritdoc IProcessorEndpoint
  function updateFeeCollector(address payable newFeeCollector) external onlyRole(ADMIN) {
    if (newFeeCollector == address(0)) revert AddressCantBeZero();
    feeCollector = newFeeCollector;
    emit FeeCollectorUpdated(newFeeCollector);
  }

  /// @inheritdoc IProcessorEndpoint
  function addAllowedDeployer(address deployer) external onlyRole(ADMIN) {
    if (deployer == address(0)) revert AddressCantBeZero();
    _grantRole(DEPLOYER_ROLE, deployer);
  }

  /// @inheritdoc IProcessorEndpoint
  function removeAllowedDeployer(address deployer) external onlyRole(ADMIN) {
    if (deployer == address(0)) revert AddressCantBeZero();
    _revokeRole(DEPLOYER_ROLE, deployer);
  }

  /// @inheritdoc IProcessorEndpoint
  function isAllowedDeployer(address deployer) external view returns (bool allowed) {
    return hasRole(DEPLOYER_ROLE, deployer);
  }

  /// @inheritdoc IProcessorEndpoint
  function getNextPendingRequest()
    external
    view
    returns (Structs.PendingRequest memory, bytes32, bool success)
  {
    if (_queueSize(_triggerQueue) > 0) {
      bytes32 requestId = _queuePeekHead(_triggerQueue);
      Structs.PendingRequest storage req = _triggerQueue.requestById[requestId];
      return (req, applicationStateRoots[req.applicationId], true);
    }
    if (_queueSize(_requestQueue) > 0) {
      bytes32 requestId = _queuePeekHead(_requestQueue);
      Structs.PendingRequest storage req = _requestQueue.requestById[requestId];
      return (req, applicationStateRoots[req.applicationId], true);
    }

    Structs.PendingRequest memory emptyReq;
    return (emptyReq, bytes32(0), false);
  }

  /// @inheritdoc IProcessorEndpoint
  function isCurrentPendingRequest(bytes32 requestId) public view returns (bool) {
    return _queueIsHead(_triggerQueue, requestId) || _queueIsHead(_requestQueue, requestId);
  }

  // Pull payment pattern functions
  function _asyncTransfer(address tokenAddress, address dest, uint256 amount) internal {
    pendingClaims[tokenAddress][dest] += amount;
    totalPendingClaims[tokenAddress] += amount;
  }

  function _getAvailableEthBalance() internal view returns (uint256) {
    return address(this).balance - totalPendingClaims[ETH_TOKEN] - totalAppCustody[ETH_TOKEN];
  }

  /// @inheritdoc IProcessorEndpoint
  function claim(address tokenAddress, address payable payee) public nonReentrant {
    uint256 amount = pendingClaims[tokenAddress][payee];
    if (amount == 0) return;

    pendingClaims[tokenAddress][payee] = 0;
    totalPendingClaims[tokenAddress] -= amount;

    emit PaymentWithdrawn(tokenAddress, payee, amount);

    if (tokenAddress == ETH_TOKEN) {
      (bool success, ) = payee.call{value: amount}('');
      if (!success) revert TransferFailed();
    } else {
      IERC20(tokenAddress).safeTransfer(payee, amount);
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function generateRequestId(
    address sender,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 idx
  ) public pure returns (bytes32) {
    return
      keccak256(
        abi.encode(sender, applicationId, requestType, payload, tokenAddress, assetAmount, idx)
      );
  }

  /// @inheritdoc IProcessorEndpoint
  function getDeployedAppIds() external view returns (uint64[] memory) {
    return _deployedAppIds;
  }

  function _removeDeployedAppId(uint64 appId) internal {
    uint256 len = _deployedAppIds.length;
    uint256 i;
    while (i != len) {
      if (_deployedAppIds[i] == appId) {
        _deployedAppIds[i] = _deployedAppIds[len - 1];
        _deployedAppIds.pop();
        break;
      }
      unchecked {
        ++i;
      }
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function adminReset() external onlyRole(RESET_OPERATOR) nonReentrant {
    _resetQueue();
  }

  function _resetQueue() internal {
    uint256 i = _requestQueue.head;
    uint256 tail = _requestQueue.tail;
    uint256 freedDeploySlots;

    // Iterate every pending request from head to tail, refunding any asset deposits to their
    // senders and counting any DEPLOYAPP entries so their reserved deploy slots can be returned.
    while (i != tail) {
      Structs.PendingRequest storage req = _requestQueue.requestById[_requestQueue.idByOrder[i]];
      if (req.requestType == Structs.RequestType.DEPLOYAPP) {
        unchecked {
          ++freedDeploySlots;
        }
      }
      if (req.assetAmount > 0) {
        appCustody[req.applicationId][req.tokenAddress] -= req.assetAmount;
        totalAppCustody[req.tokenAddress] -= req.assetAmount;
        _asyncTransfer(req.tokenAddress, req.sender, req.assetAmount);
      }
      delete _requestQueue.requestById[_requestQueue.idByOrder[i]];
      delete _requestQueue.idByOrder[i];
      unchecked {
        ++i;
      }
    }

    // Collapse the queue by setting tail back to head. _queue.head is intentionally left at its
    // current value so that future queue indices continue from where they left off rather
    // than restarting from zero (avoids any risk of re-using a slot index still in storage).
    _requestQueue.tail = _requestQueue.head;

    // Return the slots that were reserved for the now-discarded pending DEPLOYAPP requests.
    // Already-finalised apps keep their slot consumed; only in-flight deploys are freed.
    unchecked {
      availableDeploySlots += freedDeploySlots;
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function adminResetApps(uint64[] calldata appIds) external onlyRole(RESET_OPERATOR) nonReentrant {
    // Clear the pending request queue, refunding each request's asset deposit to its sender.
    _resetQueue();

    // Resolve the effective app list: use the caller-supplied list when non-empty, otherwise
    // fall back to every app that has ever been successfully deployed.
    uint64[] memory effectiveAppIds = appIds;
    if (appIds.length == 0) {
      effectiveAppIds = _deployedAppIds;
    }

    address[] memory effectiveTokens = tokenAllowlist.getAllowedTokens();

    uint256 appCount = effectiveAppIds.length;
    uint256 tokenCount = effectiveTokens.length;

    // --- Effects: zero out all custody and state roots before any external call ---
    address payable recipient = payable(msg.sender);
    uint256 i;
    while (i != appCount) {
      uint64 appId = effectiveAppIds[i];

      // Transfer ETH custody for this app and clear both the per-app and global trackers.
      uint256 ethAmt = appCustody[appId][ETH_TOKEN];
      if (ethAmt > 0) {
        (bool ok, ) = recipient.call{value: ethAmt}('');
        if (!ok) revert TransferFailed();
        totalAppCustody[ETH_TOKEN] -= ethAmt;
        appCustody[appId][ETH_TOKEN] = 0;
      }

      // Transfer ERC-20 custody for this app across every token in the effective list.
      uint256 j;
      while (j != tokenCount) {
        uint256 amt = appCustody[appId][effectiveTokens[j]];
        if (amt > 0) {
          IERC20(effectiveTokens[j]).safeTransfer(recipient, amt);
          totalAppCustody[effectiveTokens[j]] -= amt;
          appCustody[appId][effectiveTokens[j]] = 0;
        }
        unchecked {
          ++j;
        }
      }

      // Clear the app's state root, remove it from the deployed list, and return its deploy
      // slot to the pool so the same slot capacity can be reused after the reset.
      if (applicationStateRoots[appId] != bytes32(0)) {
        applicationStateRoots[appId] = bytes32(0);
        _removeDeployedAppId(appId);
        unchecked {
          ++availableDeploySlots;
        }
      }

      unchecked {
        ++i;
      }
    }
  }

  // Calls execute then withdraw on the trigger contract registered for the given applicationId,
  // if any. Each call is wrapped in an independent try/catch so that a revert in the trigger
  // never propagates to the caller: both calls are always attempted regardless of the outcome
  // of the first.
  function _invokeTrigger(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg
  ) internal {
    ITrigger trigger = triggerContracts[applicationId];
    //do nothing if trigger not defined for application
    if (address(trigger) == address(0)) return;

    // Unshield
    //best-effort transfer of app custody to trigger before invocation.
    uint256 ethAmt = appCustody[applicationId][ETH_TOKEN];
    if (ethAmt > 0) {
      (bool ok, ) = address(trigger).call{value: ethAmt}('');
      if (ok) {
        //update custodies
        appCustody[applicationId][ETH_TOKEN] = 0;
        totalAppCustody[ETH_TOKEN] -= ethAmt;
      }
    }

    //erc20 tokens
    address[] memory tokens = tokenAllowlist.getAllowedTokens();
    uint256 tokenCount = tokens.length;
    uint256 ti;
    while (ti != tokenCount) {
      uint256 amt = appCustody[applicationId][tokens[ti]];
      if (amt > 0) {
        try IERC20(tokens[ti]).transfer(address(trigger), amt) returns (bool ok) {
          if (ok) {
            appCustody[applicationId][tokens[ti]] = 0;
            totalAppCustody[tokens[ti]] -= amt;
          }
        } catch {}
      }
      unchecked {
        ++ti;
      }
    }

    // invoke execute
    bool executeSuccess = true;
    try
      trigger.execute(
        applicationId,
        prevStateRoot,
        newStateRoot,
        processedRequestId,
        userEventData,
        appEventData,
        withdrawalRequests,
        refund,
        applicationFees,
        errorCode,
        errorMsg
      )
    {} catch {
      executeSuccess = false;
    }
    emit TriggerExecuted(applicationId, processedRequestId, address(trigger), executeSuccess);

    //invoke withdraw
    (bool withdrawSuccess, ITrigger.ReturnedTokens[] memory returnedTokens) = trigger.withdraw(
      applicationId,
      prevStateRoot,
      newStateRoot,
      processedRequestId,
      userEventData,
      appEventData,
      withdrawalRequests,
      refund,
      applicationFees,
      errorCode,
      errorMsg
    );

    // Re-shield: returnedTokens contains returned tokens + ETH_TOKEN
    ti = 0;
    tokenCount = returnedTokens.length;
    while (ti != tokenCount) {
      appCustody[applicationId][returnedTokens[ti].token] = returnedTokens[ti].amount;
      totalAppCustody[returnedTokens[ti].token] += returnedTokens[ti].amount;
      unchecked {
        ++ti;
      }
    }
    emit TriggerPostWithdraw(applicationId, processedRequestId, address(trigger), withdrawSuccess);
  }

  /// @notice Enqueues a request on behalf of a registered trigger contract (FIFO).
  ///         Only callable by an address that appears in triggersToAppIds (i.e. a registered trigger).
  ///         The applicationId is derived automatically from the caller's entry in triggersToAppIds.
  ///         The requestId is generated deterministically and returned.
  /// @param requestType Request type.
  /// @param payload Request payload.
  /// @param tokenAddress Token address (address(0) = ETH).
  /// @param assetAmount Business asset amount.
  /// @param maxFeeValue Maximum fee reserved for processing.
  /// @param sender Original request sender.
  /// @param facilitator Facilitator address (address(0) for direct submissions).
  /// @return requestId Generated request identifier.
  function submitTriggerRequest(
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 maxFeeValue,
    address sender,
    address facilitator
  ) public nonReentrant returns (bytes32) {
    uint64 applicationId = triggersToAppIds[msg.sender];
    if (applicationId == 0) revert NotRegisteredTrigger();

    bytes32 requestId = generateRequestId(
      sender,
      applicationId,
      requestType,
      payload,
      tokenAddress,
      assetAmount,
      _triggerQueue.tail
    );

    _queueEnqueue(
      _triggerQueue,
      requestId,
      Structs.PendingRequest({
        timestamp: block.timestamp,
        tokenAddress: tokenAddress,
        assetAmount: assetAmount,
        maxFeeValue: maxFeeValue,
        requestId: requestId,
        payload: payload,
        sender: sender,
        facilitator: facilitator,
        applicationId: applicationId,
        protocolVersion: PROTOCOL_VERSION,
        requestType: requestType
      })
    );

    emit RequestSubmitted(applicationId, requestId, sender, facilitator);
    return requestId;
  }

  // ── Internal queue helpers ────────────────────────────────────────────────

  function _queueEnqueue(
    RequestQueue storage q,
    bytes32 id,
    Structs.PendingRequest memory req
  ) internal {
    q.requestById[id] = req;
    q.idByOrder[q.tail] = id;
    unchecked {
      ++q.tail;
    }
  }

  function _queueDequeueHead(RequestQueue storage q) internal {
    bytes32 id = q.idByOrder[q.head];
    delete q.requestById[id];
    delete q.idByOrder[q.head];
    unchecked {
      ++q.head;
    }
  }

  function _queueSize(RequestQueue storage q) internal view returns (uint256) {
    if (q.tail > q.head) return q.tail - q.head;
    return 0;
  }

  function _queuePeekHead(RequestQueue storage q) internal view returns (bytes32) {
    return q.idByOrder[q.head];
  }

  function _queueIsHead(RequestQueue storage q, bytes32 id) internal view returns (bool) {
    return q.tail > q.head && q.idByOrder[q.head] == id;
  }

  function _queueGetAll(
    RequestQueue storage q
  ) internal view returns (Structs.PendingRequest[] memory result) {
    uint256 n = _queueSize(q);
    result = new Structs.PendingRequest[](n);
    uint256 i = q.head;
    uint256 tail = q.tail;
    uint256 j;
    while (i < tail) {
      result[j] = q.requestById[q.idByOrder[i]];
      unchecked {
        ++i;
        ++j;
      }
    }
  }

  function _queueGetPage(
    RequestQueue storage q,
    uint256 offset,
    uint256 limit
  ) internal view returns (Structs.PendingRequest[] memory result) {
    uint256 n = _queueSize(q);
    if (offset >= n || limit == 0) return new Structs.PendingRequest[](0);
    uint256 end = offset + limit;
    if (end > n) end = n;
    uint256 count = end - offset;
    result = new Structs.PendingRequest[](count);
    uint256 i = q.head + offset;
    uint256 stop = q.head + end;
    uint256 j;
    while (i < stop) {
      result[j] = q.requestById[q.idByOrder[i]];
      unchecked {
        ++i;
        ++j;
      }
    }
  }
}
