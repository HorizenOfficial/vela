// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';
import '@openzeppelin/contracts/utils/Strings.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IAuthorityRegistry.sol';
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
///
///      `initialize` and `invokeTrigger` are the exceptions: they moved here purely to make room
///      under the EIP-170 limit for the upgradeability machinery added by
///      `docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`, not because they were off the hot path.
///      `invokeTrigger` in particular is called unconditionally from every successful
///      `stateUpdate`/`batchStateUpdate` entry, so that call now pays the ~2,600 gas delegatecall
///      even for applications with no registered trigger — see
///      `ProcessorEndpoint._delegateToExtensionCall`.
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

  /// @notice One-time setup for `ProcessorEndpoint`, invoked through the proxy's
  ///         `ERC1967Proxy` constructor calldata. See `IProcessorEndpoint`'s docs on the
  ///         `ProcessorEndpoint.initialize` stub that forwards here, and
  ///         `docs/design/UPGRADABLE_CONTRACTS_DESIGN.md` for why constructor logic moved into an
  ///         initializer and why it lives in this contract rather than on the endpoint itself.
  /// @param _teeAuthenticator Contract used to verify update signatures.
  /// @param _authorityRegistry Registry for authority checks.
  /// @param updateStatusOperator Initial operator for status updates.
  /// @param admin Initial admin address.
  /// @param resetOperator Address granted RESET_OPERATOR role. Pass address(0) to disable reset
  ///        permanently (required for production). The role cannot be granted after deployment.
  /// @param _minFeePerRequest Minimum fee enforced per request.
  /// @param _tokenAllowlist External token allowlist contract.
  function initialize(
    ITeeAuthenticator _teeAuthenticator,
    IAuthorityRegistry _authorityRegistry,
    address updateStatusOperator,
    address admin,
    address resetOperator,
    uint256 _minFeePerRequest,
    ITokenAllowlist _tokenAllowlist
  ) external onlyDelegateCall initializer {
    if (
      address(_teeAuthenticator) == address(0) ||
      address(_authorityRegistry) == address(0) ||
      updateStatusOperator == address(0) ||
      admin == address(0) ||
      address(_tokenAllowlist) == address(0)
    ) revert IProcessorEndpoint.AddressCantBeZero();

    __AccessControl_init();
    __EIP712_init('Vela', Strings.toString(PROTOCOL_VERSION));

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

    maxNumOfApplications = 10;
    availableDeploySlots = maxNumOfApplications;
    maxQueueSize = 10;
    selectionGrace = 60;
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
      _queueStore.deploys.tail
    );

    uint64 applicationId = uint64(bytes8(requestId)); // Derive a unique application ID from the request ID for deploy requests
    // Deploys go into the global deploy queue, not into pendingQueues[applicationId]: the
    // application does not exist yet, so it is absent from _deployedAppIds and the round-robin
    // scan over deployed applications would never reach it.
    RequestQueues.enqueue(
      _queueStore,
      _queueStore.deploys,
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

  /// @notice Updates the selection grace period. See `IProcessorEndpoint.updateSelectionGrace`.
  /// @dev Zero is a valid setting (strict enforcement of the round-robin turn, no tolerance for
  ///      the selection race), and so is a value large enough to disable enforcement — the only
  ///      way to route around a permanently failing queue head short of `adminReset`. Neither is
  ///      rejected here on purpose.
  function updateSelectionGrace(uint256 newGrace) external onlyDelegateCall onlyRole(ADMIN) {
    selectionGrace = newGrace;
    emit IProcessorEndpoint.SelectionGraceUpdated(newGrace);
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
    _resetQueues();
  }

  /// @notice Clears the queues and sweeps per-app custody to the caller. See
  ///         `IProcessorEndpoint.adminResetApps`.
  function adminResetApps(
    uint64[] calldata appIds
  ) external onlyDelegateCall onlyRole(RESET_OPERATOR) nonReentrant {
    // Clear the pending request queues, refunding each request's asset deposit to its sender.
    _resetQueues();

    // Resolve the effective app list: use the caller-supplied list when non-empty, otherwise
    // fall back to every app that has ever been successfully deployed.
    uint64[] memory effectiveAppIds = appIds;
    if (appIds.length == 0) {
      effectiveAppIds = _deployedAppIds;
    }

    // Sweeping custody, clearing state roots and dropping triggers is implemented in
    // RequestQueues. It runs under delegatecall, so the transfers move the endpoint's own ETH
    // and tokens and credit msg.sender — the reset operator.
    uint256 freedSlots = RequestQueues.resetApps(
      _deployedAppIds,
      effectiveAppIds,
      tokenAllowlist.getAllowedTokens(),
      payable(msg.sender),
      applicationStateRoots,
      appCustody,
      totalAppCustody,
      triggerContracts,
      triggersToAppIds
    );
    unchecked {
      availableDeploySlots += freedSlots;
    }
  }

  /// @dev Discards every pending request — the deploy queue, every per-application queue and the
  ///      trigger queue — refunding asset deposits. Implemented in RequestQueues; the freed
  ///      deploy slots come back as a return value because a value-type state variable cannot be
  ///      passed to a library by reference.
  function _resetQueues() private {
    uint256 freedDeploySlots = RequestQueues.resetQueues(
      _queueStore,
      _deployedAppIds,
      appCustody,
      totalAppCustody,
      pendingClaims,
      totalPendingClaims,
      triggerContracts,
      triggersToAppIds
    );
    // Return the slots that were reserved for the now-discarded pending DEPLOYAPP requests.
    // Already-finalised apps keep their slot consumed; only in-flight deploys are freed.
    unchecked {
      availableDeploySlots += freedDeploySlots;
    }
  }

  /// @notice Claims pending balance for a given token and payee. See `IProcessorEndpoint.claim`;
  ///         `ProcessorEndpoint.claim` is the entry point callers use.
  /// @dev The only copy of this logic: `ProcessorEndpoint` no longer has a resident `_claim` of
  ///      its own now that `invokeTrigger` (below), its one other caller, also lives here.
  function claim(
    address tokenAddress,
    address payable payee
  ) external onlyDelegateCall nonReentrant {
    _claim(tokenAddress, payee);
  }

  //logic is reentrant-safe because it is also invoked from within invokeTrigger below
  function _claim(address tokenAddress, address payable payee) private {
    uint256 amount = pendingClaims[tokenAddress][payee];
    if (amount == 0) return;

    pendingClaims[tokenAddress][payee] = 0;
    totalPendingClaims[tokenAddress] -= amount;

    emit IProcessorEndpoint.PaymentWithdrawn(tokenAddress, payee, amount);

    if (tokenAddress == ETH_TOKEN) {
      (bool success, ) = payee.call{value: amount}('');
      if (!success) revert IProcessorEndpoint.TransferFailed();
    } else {
      IERC20(tokenAddress).safeTransfer(payee, amount);
    }
  }

  /// @notice Invokes the trigger contract registered for `applicationId` (execute, withdraw, and
  ///         the trusted-payload callback), if any. Called unconditionally from
  ///         `ProcessorEndpoint._processOneStateUpdate` on every successful `stateUpdate`/
  ///         `batchStateUpdate` entry, via `_delegateToExtensionCall` — see that function's docs
  ///         in `ProcessorEndpoint.sol` for why this moved here and at what per-call cost
  ///         (`docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`).
  /// @dev Each external call to the trigger is wrapped in an independent try/catch so that a
  ///      revert in the trigger never propagates to the caller: both calls are always attempted
  ///      regardless of the outcome of the first.
  function invokeTrigger(
    uint64 applicationId,
    bytes32 processedRequestId,
    Structs.EventData memory appEventData,
    address[] memory claimable
  ) external onlyDelegateCall {
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
    emit IProcessorEndpoint.TriggerExecuted(applicationId, processedRequestId, executeSuccess);

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

    emit IProcessorEndpoint.TriggerWithdraw(
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
  ///         reachable ONLY from invokeTrigger (i.e. during stateUpdate), which is
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
      _queueStore.triggers.tail
    );

    RequestQueues.enqueue(
      _queueStore,
      _queueStore.triggers,
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

    emit IProcessorEndpoint.RequestSubmitted(applicationId, requestId, trigger, address(0));
  }
}
