// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IProcessorEndpoint.sol';
import './interfaces/IAuthorityRegistry.sol';
import './interfaces/ITokenAllowlist.sol';
import './Structs.sol';
import './ProcessorEndpointStorage.sol';
import './interfaces/ITrigger.sol';

/// @title ProcessorEndpoint
/// @notice Implementation of the processor endpoint interface.
/// @dev State lives in `ProcessorEndpointStorage`, which `ProcessorEndpointExtension` also derives
///      from: this contract is close enough to the EIP-170 deployed-bytecode limit that parts of it
///      are hosted in the extension and reached by `delegatecall`. See
///      `ProcessorEndpointStorage` for the rules that keeps safe, and
///      `docs/design/PROCESSOR_ENDPOINT_SPLIT.md` for the rationale.
contract ProcessorEndpoint is ProcessorEndpointStorage, IProcessorEndpoint {
  using SafeERC20 for IERC20;

  /// @dev Extension contract holding code moved out of this one for size reasons. Immutable, so it
  ///      lives in this contract's code rather than in storage and cannot be repointed.
  address private immutable _extension;

  /// @param _teeAuthenticator Contract used to verify update signatures.
  /// @param _authorityRegistry Registry for authority checks.
  /// @param updateStatusOperator Initial operator for status updates.
  /// @param admin Initial admin address.
  /// @param resetOperator Address granted RESET_OPERATOR role. Pass address(0) to disable reset
  ///        permanently (required for production). The role cannot be granted after deployment.
  /// @param _minFeePerRequest Minimum fee enforced per request.
  /// @param _tokenAllowlist External token allowlist contract.
  /// @param extension Deployed `ProcessorEndpointExtension` serving the delegated entry points.
  constructor(
    ITeeAuthenticator _teeAuthenticator,
    IAuthorityRegistry _authorityRegistry,
    address updateStatusOperator,
    address admin,
    address resetOperator,
    uint256 _minFeePerRequest,
    ITokenAllowlist _tokenAllowlist,
    address extension
  ) {
    if (
      address(_teeAuthenticator) == address(0) ||
      address(_authorityRegistry) == address(0) ||
      updateStatusOperator == address(0) ||
      admin == address(0) ||
      address(_tokenAllowlist) == address(0) ||
      extension == address(0)
    ) revert AddressCantBeZero();

    // A delegatecall to an address without code succeeds and returns nothing, so a wrong
    // extension address would turn submitRequestFor into a silent no-op that keeps the fee
    // instead of reverting. _extension is immutable, so this is the only chance to catch it.
    if (extension.code.length == 0) revert InvalidExtension();

    _extension = extension;
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
    // DEPLOYAPP has its own entrypoint; TRUSTPROCESS is trusted and can ONLY be
    // created internally during stateUpdate (via a trigger's getTrustProcessPayload payload).
    if (
      requestType == Structs.RequestType.DEPLOYAPP ||
      requestType == Structs.RequestType.TRUSTPROCESS
    ) revert InvalidRequestType();

    //check values
    if (maxFeeValue < minFeePerRequest) revert FeeValueBelowMinimum();
    //check queue size
    if (_pendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

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

    _addToCustody(applicationId, tokenAddress, assetAmount);

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
  /// @dev Implemented in `ProcessorEndpointExtension` for size reasons; this entry point only
  ///      forwards the call. Declaring it here keeps the ABI, the selector, the address relayers
  ///      call and the `IProcessorEndpoint` conformance check unchanged. `_delegateToExtension`
  ///      never returns, so the declared return value is never produced here — the extension's
  ///      return data is passed back to the caller verbatim.
  ///
  ///      The parameters are named even though the body ignores them, because the names are part
  ///      of the published ABI: wallets, explorers and generated bindings show them. The no-op
  ///      statements below are what keeps solc from warning about unused parameters; they cost no
  ///      bytecode.
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
  ) external payable returns (bytes32) {
    sender;
    protocolVersion;
    applicationId;
    requestType;
    payload;
    tokenAddress;
    assetAmount;
    deadline;
    requestSignature;
    depositPermit;
    _delegateToExtension();
  }

  /// @dev Forwards the current call to `_extension` with `delegatecall`, so the extension's code
  ///      runs against this contract's storage, balance, `msg.sender` and `msg.value`. Returns the
  ///      extension's return data, or bubbles up its revert data unchanged, and never returns to
  ///      the caller of this function.
  ///
  ///      Scratches memory above the free-memory pointer rather than at offset 0, and is annotated
  ///      `memory-safe` accordingly: unannotated assembly would switch off solc's `memoryguard`
  ///      for the whole contract, and `stateUpdate` then fails to compile with "stack too deep".
  function _delegateToExtension() private {
    address extension = _extension;
    assembly ('memory-safe') {
      let ptr := mload(0x40)
      calldatacopy(ptr, 0, calldatasize())
      let ok := delegatecall(gas(), extension, ptr, calldatasize(), 0, 0)
      let size := returndatasize()
      returndatacopy(ptr, 0, size)
      switch ok
      case 0 {
        revert(ptr, size)
      }
      default {
        return(ptr, size)
      }
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  ) external payable validProtocolVersion(protocolVersion) nonReentrant returns (bytes32) {
    return _submitDeployRequest(protocolVersion, payload);
  }

  /// @inheritdoc IProcessorEndpoint
  function submitDeployRequestWithTrigger(
    uint8 protocolVersion,
    bytes calldata payload,
    address trigger
  ) external payable validProtocolVersion(protocolVersion) nonReentrant returns (bytes32) {
    bytes32 requestId = _submitDeployRequest(protocolVersion, payload);
    // Optional trigger registration that does NOT consume the payload, so the
    // deploy can still carry a full WASM descriptor. address(0) means "no
    // trigger" (identical to the 2-arg overload). The trigger is validated and
    // registered eagerly here so that an invalid/duplicate trigger reverts the
    // submit (instead of failing later inside stateUpdate). If the deploy then
    // fails on-chain, the registration is rolled back in stateUpdate.
    if (trigger != address(0)) {
      if (triggersToAppIds[trigger] != 0) revert TriggerAlreadyRegistered();
      if (trigger.code.length == 0) revert TriggerCannotBeEOA();

      uint64 applicationId = uint64(bytes8(requestId));
      triggerContracts[applicationId] = ITrigger(trigger);
      triggersToAppIds[trigger] = applicationId;
    }
    return requestId;
  }

  function _submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  ) private returns (bytes32 requestId) {
    if (!hasRole(DEPLOYER_ROLE, msg.sender)) revert DeployerNotAllowed();
    if (availableDeploySlots == 0) revert MaxNumOfApplicationsExceeded();
    //check queue size
    if (_pendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();
    if (msg.value < minFeePerRequest) revert FeeValueBelowMinimum();

    --availableDeploySlots;

    Structs.RequestType requestType = Structs.RequestType.DEPLOYAPP;
    //create request
    requestId = _generateRequestId(
      msg.sender,
      0, // deploy requests have applicationId 0, a unique applicationId will be derived from the requestId for each deploy request to avoid collisions with regular requests and to group deploy requests together
      requestType,
      keccak256(payload),
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

  function _removeRequest(bytes32 requestId) private {
    // _queueIsHead already returns false for an empty queue (tail > head check),
    // so no separate size guard is needed here.
    if (_queueIsHead(_triggerQueue, requestId)) {
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
    return _pendingRequestsSize();
  }

  /// @inheritdoc IProcessorEndpoint
  function getTriggerQueueSize() public view returns (uint256) {
    return _queueSize(_triggerQueue);
  }

  /// @inheritdoc IProcessorEndpoint
  function getTriggerRequests() external view returns (Structs.PendingRequest[] memory) {
    return _queueGetAll(_triggerQueue);
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
      _subtractToCustody(applicationId, reqTokenAddress, assetAmount);

      // Refund business-asset deposit in its original token to the user.
      // Fee refund is always in ETH, routed to facilitator if present.
      // TRUSTPROCESS requests have maxFeeValue=0, so feeRefund is 0 (nothing to refund).
      uint256 feeRefund = requestInfo.maxFeeValue >= minFeePerRequest
        ? requestInfo.maxFeeValue - minFeePerRequest
        : 0;
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
        // Deploy failed: roll back the trigger that was registered eagerly at
        // submit time so the trigger address can be reused by a future deploy.
        _clearTrigger(applicationId);
      }

      _markRequestCompleted(
        applicationId,
        processedRequestId,
        fromTriggerQueue ? 0 : minFeePerRequest,
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

    // don't check fees if we are from trigger queue
    if (!fromTriggerQueue) {
      if (refund + applicationFees != maxFeeValue) revert InvalidValue();
      if (applicationFees < minFeePerRequest) {
        revert InvalidValue();
      }
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
    while (i != eventsLength) {
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
    while (i != appEventsLength) {
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
      // The trigger (if any) was already validated and registered eagerly in
      // submitDeployRequestWithTrigger, so a successful deploy needs no further
      // trigger bookkeeping here.
    }
    emit StateRootUpdate(applicationId, processedRequestId, prevStateRoot, newStateRoot);

    //credit refund to feeRecipient's pending balance (pull pattern) — refund is always ETH
    if (refund > 0) {
      _asyncTransfer(ETH_TOKEN, feeRecipient, refund);
      emit Refund(applicationId, processedRequestId, feeRecipient, ETH_TOKEN, refund);
    }

    //credit withdrawals to receivers' pending balances
    i = 0;
    uint256 insertIntoClaimable;
    address trigger = address(triggerContracts[applicationId]);
    address[] memory claimableTemp = new address[](withdrawalsLength);
    while (i < withdrawalsLength) {
      // Only classify a withdrawal as "claimable by the trigger" when a trigger is actually
      // registered. Without this guard, trigger == address(0) would match withdrawals sent to
      // the zero address and route them into the trigger-claim path. A user mistakenly
      // withdrawing to address(0) is their own problem, but it must never be handed to a trigger.
      if (trigger != address(0) && withdrawalRequests[i].receiver == trigger) {
        claimableTemp[insertIntoClaimable] = withdrawalRequests[i].tokenAddress;
        unchecked {
          ++insertIntoClaimable;
        }
      }
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

    //reduce claimable array size to correct one
    address[] memory claimable = new address[](insertIntoClaimable);
    i = 0;
    while (i != insertIntoClaimable) {
      claimable[i] = claimableTemp[i];
      unchecked {
        ++i;
      }
    }

    //invoke trigger contracts, if registered
    _invokeTrigger(applicationId, processedRequestId, appEventData, claimable);

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
    return _claim(tokenAddress, payee);
  }
  //logic is reentrable because it will be invoked in invokeTrigger
  function _claim(address tokenAddress, address payable payee) internal {
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
    bytes32 payloadHash,
    address tokenAddress,
    uint256 assetAmount,
    uint256 idx
  ) public pure returns (bytes32) {
    return
      _generateRequestId(
        sender,
        applicationId,
        requestType,
        payloadHash,
        tokenAddress,
        assetAmount,
        idx
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
        // A pending deploy may have registered a trigger eagerly at submit time. Since the
        // deploy is being discarded (and its derived appId never appears in _deployedAppIds),
        // clear the registration here so the trigger address can be reused.
        _clearTrigger(req.applicationId);
      }
      if (req.assetAmount > 0) {
        _subtractToCustody(req.applicationId, req.tokenAddress, req.assetAmount);
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

    // Drain the trigger queue too: pending TRUSTPROCESS requests reference apps that may be
    // reset, and must not survive a reset. They carry no funds (assetAmount/maxFeeValue == 0),
    // so clearing their storage is sufficient — no refunds needed.
    uint256 ti = _triggerQueue.head;
    uint256 triggerTail = _triggerQueue.tail;
    while (ti != triggerTail) {
      delete _triggerQueue.requestById[_triggerQueue.idByOrder[ti]];
      delete _triggerQueue.idByOrder[ti];
      unchecked {
        ++ti;
      }
    }
    _triggerQueue.tail = _triggerQueue.head;

    // Return the slots that were reserved for the now-discarded pending DEPLOYAPP requests.
    // Already-finalised apps keep their slot consumed; only in-flight deploys are freed.
    unchecked {
      availableDeploySlots += freedDeploySlots;
    }
  }

  /// @dev Removes any trigger registration for the given application (both the forward and
  ///      reverse mappings). No-op when no trigger is registered.
  function _clearTrigger(uint64 applicationId) private {
    ITrigger trigger = triggerContracts[applicationId];
    if (address(trigger) != address(0)) {
      delete triggersToAppIds[address(trigger)];
      delete triggerContracts[applicationId];
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
        _subtractToCustody(appId, ETH_TOKEN, ethAmt);
      }

      // Transfer ERC-20 custody for this app across every token in the effective list.
      uint256 j;
      while (j != tokenCount) {
        uint256 amt = appCustody[appId][effectiveTokens[j]];
        if (amt > 0) {
          IERC20(effectiveTokens[j]).safeTransfer(recipient, amt);
          _subtractToCustody(appId, effectiveTokens[j], amt);
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

      // Clear any trigger registered for this app so its address can be reused after reset.
      _clearTrigger(appId);

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
    bytes32 processedRequestId,
    Structs.EventData calldata appEventData,
    address[] memory claimable
  ) internal {
    ITrigger trigger = triggerContracts[applicationId];
    //do nothing if trigger not defined for application
    if (address(trigger) == address(0)) return;

    //claims for the trigger
    uint256 ti;
    uint256 tokenCount = claimable.length;
    while (ti != tokenCount) {
      _claim(claimable[ti], payable(address(trigger)));
      unchecked {
        ++ti;
      }
    }

    // invoke execute
    bool executeSuccess = true;
    try trigger.execute(appEventData) {} catch {
      executeSuccess = false;
    }
    emit TriggerExecuted(applicationId, processedRequestId, executeSuccess);

    //invoke withdraw (sweep only)
    Structs.TokenAndAmount[] memory returnedTokens;
    Structs.TokenAndAmount[] memory failedTokens;
    bool withdrawSuccess = true;
    try trigger.withdraw() returns (
      Structs.TokenAndAmount[] memory _returnedTokens,
      Structs.TokenAndAmount[] memory _failedTokens
    ) {
      returnedTokens = _returnedTokens;
      failedTokens = _failedTokens;
    } catch {
      withdrawSuccess = false;
    }

    // Re-shield: returnedTokens contains returned tokens + ETH_TOKEN
    ti = 0;
    tokenCount = returnedTokens.length;
    while (ti != tokenCount) {
      _addToCustody(applicationId, returnedTokens[ti].token, returnedTokens[ti].amount);
      unchecked {
        ++ti;
      }
    }

    // Explicit, isolated trusted-payload step, decoupled from the sweep: getTrustProcessPayload
    // is called here in its own try/catch (a revert cannot block stateUpdate). It runs even when
    // withdraw failed, so the application can react to a failed sweep (damage control).
    bytes memory trustedPayload;
    bool postWithdrawSuccess = true;
    try
      trigger.getTrustProcessPayload(
        appEventData,
        executeSuccess,
        withdrawSuccess,
        returnedTokens,
        failedTokens
      )
    returns (bytes memory _trustedPayload) {
      trustedPayload = _trustedPayload;
    } catch {
      postWithdrawSuccess = false;
    }

    emit TriggerWithdraw(
      applicationId,
      processedRequestId,
      withdrawSuccess,
      postWithdrawSuccess,
      returnedTokens,
      failedTokens
    );

    // stateUpdate is the ONLY place a trusted (TRUSTPROCESS) request can be
    // created. If getTrustProcessPayload returned a non-empty payload, enqueue it
    // into the trigger queue.
    if (trustedPayload.length > 0) {
      _enqueueTrustedRequest(applicationId, address(trigger), trustedPayload);
    }
  }

  /// @notice Enqueues a TRUSTPROCESS request into the trigger queue. PRIVATE and
  ///         reachable ONLY from _invokeTrigger (i.e. during stateUpdate), which is
  ///         the single authorized point allowed to create a trusted request. The
  ///         payload is produced by the trigger's getTrustProcessPayload hook; callers must
  ///         skip empty payloads (no trusted request is created for them).
  /// @param applicationId Application the trusted request belongs to.
  /// @param trigger Trigger contract that produced the payload (recorded as sender).
  /// @param payload Trusted request payload (forwarded to the WASM as-is).
  function _enqueueTrustedRequest(
    uint64 applicationId,
    address trigger,
    bytes memory payload
  ) private {
    bytes32 requestId = _generateRequestId(
      trigger,
      applicationId,
      Structs.RequestType.TRUSTPROCESS,
      keccak256(payload),
      address(0),
      0,
      _triggerQueue.tail
    );

    _queueEnqueue(
      _triggerQueue,
      requestId,
      Structs.PendingRequest({
        timestamp: block.timestamp,
        tokenAddress: address(0),
        assetAmount: 0,
        maxFeeValue: 0,
        requestId: requestId,
        payload: payload,
        sender: trigger,
        facilitator: address(0),
        applicationId: applicationId,
        protocolVersion: PROTOCOL_VERSION,
        requestType: Structs.RequestType.TRUSTPROCESS
      })
    );

    emit RequestSubmitted(applicationId, requestId, trigger, address(0));
  }

  // Internal queue helpers

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

  function _subtractToCustody(uint64 applicationId, address token, uint256 amount) internal {
    appCustody[applicationId][token] -= amount;
    totalAppCustody[token] -= amount;
  }
}
