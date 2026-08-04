// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';

import './interfaces/IProcessorEndpoint.sol';
import './interfaces/ITokenAllowlist.sol';
import './interfaces/ITrigger.sol';
import './Structs.sol';
import './ProcessorEndpointStorage.sol';

/// @title ProcessorEndpointExtension
/// @notice Code moved out of `ProcessorEndpoint` to stay under the EIP-170 deployed-bytecode
///         limit. It is never called directly: `ProcessorEndpoint` forwards matching calls here
///         with `delegatecall`, so this code runs against the endpoint's storage, balance,
///         `msg.sender` and `msg.value`.
/// @dev Holds the facilitator path (`submitRequestFor`), the deploy submission entry points, the
///      operator reset entry points and the admin configuration setters — everything moved out is
///      off the per-request hot path, so the one extra cold-account access the forwarding costs
///      (~2,600 gas) is paid only by rare or privileged calls. Because the boundary is a
///      `delegatecall` onto the shared `ProcessorEndpointStorage` layout, arguments are not
///      re-encoded and internal helpers can be reused, unlike an external-call or linked-library
///      split. Helpers this contract needs are duplicated into its bytecode, which costs nothing:
///      only the endpoint is near the size limit.
contract ProcessorEndpointExtension is ProcessorEndpointStorage {
  using SafeERC20 for IERC20;

  /// @dev This contract's own address, fixed at deployment. Under `delegatecall`,
  ///      `address(this)` is the endpoint instead, which is the only context where the code is
  ///      meaningful: called directly it would read and write this contract's own (empty) storage
  ///      and could strand ETH here.
  address private immutable _self = address(this);

  /// @notice Reverts when the extension is called directly rather than through the endpoint.
  error DirectCallNotAllowed();

  modifier onlyDelegateCall() {
    if (address(this) == _self) revert DirectCallNotAllowed();
    _;
  }

  /// @notice Submits a request on behalf of `sender`, authorized by an EIP-712 signature, with
  ///         the calling facilitator paying the gas and the fee. See
  ///         `IProcessorEndpoint.submitRequestFor`, whose declaration and documentation this
  ///         implements; `ProcessorEndpoint.submitRequestFor` is the entry point callers use.
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
    onlyDelegateCall
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    nonReentrant
    returns (bytes32)
  {
    // 1. Only ASSOCIATEKEY and PROCESS are supported. TRUSTPROCESS is trusted
    //    and can ONLY be created internally during stateUpdate (via a trigger's
    //    getTrustProcessPayload payload) — never submitted by an external caller here.
    if (
      requestType != Structs.RequestType.ASSOCIATEKEY && requestType != Structs.RequestType.PROCESS
    ) revert IProcessorEndpoint.InvalidRequestType();

    // 2. Verify deadline not expired
    if (block.timestamp > deadline) revert IProcessorEndpoint.DeadlineExpired();

    // 3. Check queue size
    if (_pendingRequestsSize() >= maxQueueSize) revert IProcessorEndpoint.QueueThresholdExceeded();

    // 4. Validate payload for ASSOCIATEKEY
    if (requestType == Structs.RequestType.ASSOCIATEKEY) {
      if (payload.length != 133 && payload.length != 226)
        revert IProcessorEndpoint.InvalidPayload();
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

    // 6. Recover user address from EIP-712 request signature and verify.
    //    tryRecover, not recover: a signature ECDSA cannot parse must revert with
    //    InvalidSignature, which IProcessorEndpoint declares. ECDSA.recover would raise
    //    ECDSAInvalidSignature/-Length/-S instead, and those are absent from the endpoint's ABI —
    //    this code runs by delegatecall, so callers only ever see the endpoint — leaving clients
    //    with revert data they cannot decode.
    (address recoveredSigner, ECDSA.RecoverError recoverError, ) = ECDSA.tryRecover(
      digest,
      requestSignature
    );
    if (recoverError != ECDSA.RecoverError.NoError || recoveredSigner == address(0))
      revert IProcessorEndpoint.InvalidSignature();
    if (recoveredSigner != sender) revert IProcessorEndpoint.InvalidSigner();

    // 7. Consume nonce (replay protection)
    facilitatorNonces[sender] = nonce + 1;

    // 8. Validate fee
    uint256 maxFeeValue = msg.value;
    if (maxFeeValue < minFeePerRequest) revert IProcessorEndpoint.FeeValueBelowMinimum();

    // 9. Validate token and handle deposit
    if (assetAmount == 0 && tokenAddress != ETH_TOKEN) revert IProcessorEndpoint.InvalidValue();
    if (assetAmount > 0) {
      if (tokenAddress == ETH_TOKEN) revert IProcessorEndpoint.InvalidValue();
      if (!tokenAllowlist.isAllowedToken(tokenAddress)) revert ITokenAllowlist.TokenNotAllowed();

      // Decode deposit permit and execute EIP-2612 permit + transferFrom
      if (depositPermit.length != 96) revert IProcessorEndpoint.InvalidPermit();
      (uint8 v, bytes32 r, bytes32 s) = abi.decode(depositPermit, (uint8, bytes32, bytes32));

      // Check current allowance before calling permit
      IERC20 token = IERC20(tokenAddress);
      if (token.allowance(sender, address(this)) < assetAmount) {
        IERC20Permit(tokenAddress).permit(sender, address(this), assetAmount, deadline, v, r, s);
      }

      _pullERC20(tokenAddress, sender, assetAmount);

      _addToCustody(applicationId, tokenAddress, assetAmount);
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

  /// @notice Submits a deploy request. See `IProcessorEndpoint.submitDeployRequest`, whose
  ///         declaration and documentation this implements;
  ///         `ProcessorEndpoint.submitDeployRequest` is the entry point callers use.
  function submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  )
    external
    payable
    onlyDelegateCall
    validProtocolVersion(protocolVersion)
    nonReentrant
    returns (bytes32)
  {
    return _submitDeployRequest(protocolVersion, payload);
  }

  /// @notice Submits a deploy request and registers a trigger for the derived application. See
  ///         `IProcessorEndpoint.submitDeployRequestWithTrigger`.
  function submitDeployRequestWithTrigger(
    uint8 protocolVersion,
    bytes calldata payload,
    address trigger
  )
    external
    payable
    onlyDelegateCall
    validProtocolVersion(protocolVersion)
    nonReentrant
    returns (bytes32)
  {
    bytes32 requestId = _submitDeployRequest(protocolVersion, payload);
    // Optional trigger registration that does NOT consume the payload, so the
    // deploy can still carry a full WASM descriptor. address(0) means "no
    // trigger" (identical to the 2-arg overload). The trigger is validated and
    // registered eagerly here so that an invalid/duplicate trigger reverts the
    // submit (instead of failing later inside stateUpdate). If the deploy then
    // fails on-chain, the registration is rolled back in stateUpdate.
    if (trigger != address(0)) {
      if (triggersToAppIds[trigger] != 0) revert IProcessorEndpoint.TriggerAlreadyRegistered();
      if (trigger.code.length == 0) revert IProcessorEndpoint.TriggerCannotBeEOA();

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
    if (!hasRole(DEPLOYER_ROLE, msg.sender)) revert IProcessorEndpoint.DeployerNotAllowed();
    if (availableDeploySlots == 0) revert IProcessorEndpoint.MaxNumOfApplicationsExceeded();
    //check queue size
    if (_pendingRequestsSize() >= maxQueueSize) revert IProcessorEndpoint.QueueThresholdExceeded();
    if (msg.value < minFeePerRequest) revert IProcessorEndpoint.FeeValueBelowMinimum();

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
    emit IProcessorEndpoint.DeployRequestSubmitted(applicationId, requestId, msg.sender);

    return requestId;
  }

  /// @notice Updates the pending-queue threshold. See `IProcessorEndpoint.updateQueueThreshold`.
  function updateQueueThreshold(uint256 newThreshold) external onlyDelegateCall onlyRole(ADMIN) {
    if (newThreshold == 0) revert IProcessorEndpoint.InvalidValue();
    maxQueueSize = newThreshold;
    emit IProcessorEndpoint.QueueThresholdUpdated(newThreshold);
  }

  /// @notice Updates the maximum number of applications. See
  ///         `IProcessorEndpoint.updateMaxNumOfApplications`.
  function updateMaxNumOfApplications(uint256 newMax) external onlyDelegateCall onlyRole(ADMIN) {
    if (newMax == 0) revert IProcessorEndpoint.InvalidValue();
    uint256 deployedApps = maxNumOfApplications - availableDeploySlots;
    if (newMax < deployedApps) revert IProcessorEndpoint.InvalidValue();
    uint256 oldMax = maxNumOfApplications;
    maxNumOfApplications = newMax;
    availableDeploySlots = newMax - deployedApps;
    emit IProcessorEndpoint.MaxNumberOfAppUpdated(oldMax, newMax);
  }

  /// @notice Updates the fee collector. See `IProcessorEndpoint.updateFeeCollector`.
  function updateFeeCollector(
    address payable newFeeCollector
  ) external onlyDelegateCall onlyRole(ADMIN) {
    if (newFeeCollector == address(0)) revert IProcessorEndpoint.AddressCantBeZero();
    feeCollector = newFeeCollector;
    emit IProcessorEndpoint.FeeCollectorUpdated(newFeeCollector);
  }

  /// @notice Grants the deployer role. See `IProcessorEndpoint.addAllowedDeployer`.
  function addAllowedDeployer(address deployer) external onlyDelegateCall onlyRole(ADMIN) {
    if (deployer == address(0)) revert IProcessorEndpoint.AddressCantBeZero();
    _grantRole(DEPLOYER_ROLE, deployer);
  }

  /// @notice Revokes the deployer role. See `IProcessorEndpoint.removeAllowedDeployer`.
  function removeAllowedDeployer(address deployer) external onlyDelegateCall onlyRole(ADMIN) {
    if (deployer == address(0)) revert IProcessorEndpoint.AddressCantBeZero();
    _revokeRole(DEPLOYER_ROLE, deployer);
  }

  /// @notice Clears the pending request queues. See `IProcessorEndpoint.adminReset`.
  function adminReset() external onlyDelegateCall onlyRole(RESET_OPERATOR) nonReentrant {
    _resetQueue();
  }

  /// @notice Clears the queues and sweeps per-app custody to the caller. See
  ///         `IProcessorEndpoint.adminResetApps`.
  function adminResetApps(
    uint64[] calldata appIds
  ) external onlyDelegateCall onlyRole(RESET_OPERATOR) nonReentrant {
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
        if (!ok) revert IProcessorEndpoint.TransferFailed();
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

  function _resetQueue() private {
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

  function _removeDeployedAppId(uint64 appId) private {
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
}
