// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IProcessorEndpoint.sol';
import './interfaces/IAuthorityRegistry.sol';
import './interfaces/ITokenAllowlist.sol';
import './Structs.sol';
import './RequestQueues.sol';
import './UpdateEntryHash.sol';
import './ProcessorEndpointStorage.sol';
import './interfaces/ITrigger.sol';

/// @title ProcessorEndpoint
/// @notice Implementation of the processor endpoint interface.
/// @dev State lives in `ProcessorEndpointStorage`, which `ProcessorEndpointExtension` also
///      derives from: this contract is close enough to the EIP-170 deployed-bytecode limit that
///      parts of it are hosted in the extension and reached by `delegatecall`. See
///      `ProcessorEndpointStorage` for the rules that keeps safe, and
///      `docs/design/PROCESSOR_ENDPOINT_SPLIT.md` for the rationale.
contract ProcessorEndpoint is ProcessorEndpointStorage, IProcessorEndpoint {
  using SafeERC20 for IERC20;

  /// @dev Extension contract holding code moved out of this one for size reasons. Immutable, so
  ///      it lives in this contract's code rather than in storage and cannot be repointed.
  address private immutable _extension;

  /// @param _teeAuthenticator Contract used to verify update signatures.
  /// @param _authorityRegistry Registry for authority checks.
  /// @param updateStatusOperator Initial operator for status updates.
  /// @param admin Initial admin address.
  /// @param resetOperator Address granted RESET_OPERATOR role. Pass address(0) to disable reset
  ///        permanently (required for production). The role cannot be granted after deployment.
  /// @param _minFeePerRequest Minimum fee enforced per request.
  /// @param _tokenAllowlist External token allowlist contract.
  /// @param extensionAddress Deployed `ProcessorEndpointExtension` serving the delegated entry points.
  constructor(
    ITeeAuthenticator _teeAuthenticator,
    IAuthorityRegistry _authorityRegistry,
    address updateStatusOperator,
    address admin,
    address resetOperator,
    uint256 _minFeePerRequest,
    ITokenAllowlist _tokenAllowlist,
    address extensionAddress
  ) {
    if (
      address(_teeAuthenticator) == address(0) ||
      address(_authorityRegistry) == address(0) ||
      updateStatusOperator == address(0) ||
      admin == address(0) ||
      address(_tokenAllowlist) == address(0) ||
      extensionAddress == address(0)
    ) revert AddressCantBeZero();

    // A delegatecall to an address without code succeeds and returns nothing, so a wrong
    // extension address would turn submitRequestFor into a silent no-op that keeps the fee
    // instead of reverting. _extension is immutable, so this is the only chance to catch it.
    if (extensionAddress.code.length == 0) revert InvalidExtension();

    _extension = extensionAddress;
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
  ///      runs against this contract's storage, balance, `msg.sender` and `msg.value`. Returns
  ///      the extension's return data, or bubbles up its revert data unchanged, and never returns
  ///      to the caller of this function.
  ///      Scratches memory above the free-memory pointer rather than at offset 0, and is
  ///      annotated `memory-safe` accordingly: unannotated assembly would switch off solc's
  ///      `memoryguard` for the whole contract, and `stateUpdate` then fails to compile with
  ///      "stack too deep".
  function _delegateToExtension() private {
    address target = _extension;
    assembly ('memory-safe') {
      let ptr := mload(0x40)
      calldatacopy(ptr, 0, calldatasize())
      let ok := delegatecall(gas(), target, ptr, calldatasize(), 0, 0)
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
  function extension() external view returns (address) {
    return _extension;
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`. See the note on
  ///      `submitRequestFor` for how these stubs work and why the parameter names stay here.
  function submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  ) external payable returns (bytes32) {
    protocolVersion;
    payload;
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function submitDeployRequestWithTrigger(
    uint8 protocolVersion,
    bytes calldata payload,
    address trigger
  ) external payable returns (bytes32) {
    protocolVersion;
    payload;
    trigger;
    _delegateToExtension();
  }

  /// @dev The queue a request belongs to, derived from its type alone. This is the single place
  ///      that encodes the type ⇔ queue mapping documented on `RequestQueues.Store`; both the
  ///      validity check and the removal in `stateUpdate` resolve their queue through it, so the
  ///      request that is verified to be the head is always the one that gets dequeued.
  ///
  ///      Resolving the queue is what a search over all three would otherwise cost: for a plain
  ///      request, `isHead` on the empty trigger and deploy queues reads four cold slots
  ///      (`tail`/`head` of each, ~8.4k gas) purely to rule them out. The type comes from the
  ///      slot that packs `facilitator`, `applicationId`, `protocolVersion` and `requestType`,
  ///      which `stateUpdate` loads anyway, so routing adds no storage read of its own.
  ///
  ///      Note that DEPLOYAPP is enum value 0: an unknown requestId reads back as a zeroed
  ///      request and resolves to the deploy queue. That is harmless — `isHead` still compares
  ///      the id against the queue's actual head — but it is why callers must treat the `isHead`
  ///      result, not the resolved queue, as proof that the request exists.
  function _queueOf(
    uint64 applicationId,
    Structs.RequestType requestType
  ) private view returns (RequestQueues.Queue storage) {
    if (requestType == Structs.RequestType.TRUSTPROCESS) return _queueStore.triggers;
    if (requestType == Structs.RequestType.DEPLOYAPP) return _queueStore.deploys;
    return _queueStore.pending[applicationId];
  }

  /// @dev Enforces the round-robin turn for a request taken from a per-application queue: the
  ///      submitted application must be reached by the cursor scan no later than the first queue
  ///      head older than `selectionGrace`, so the manager cannot starve an application by
  ///      serving another out of turn. Without this, the only ordering the contract enforced was
  ///      FIFO *within* an application, leaving cross-application selection to the manager's
  ///      discretion.
  ///
  ///      `selectionGrace` is what makes the rule race-free. `getPendingRequestsWithStateRoot` is
  ///      a `view`: it leaves no on-chain trace, so the contract cannot know when the manager read
  ///      it, and a request enqueued between that read and this call legitimately changes the
  ///      scan's result. Reconstructing the read instant is impossible — any value the manager
  ///      supplied for it could be chosen to justify its pick — so instead heads younger than
  ///      `selectionGrace` are skipped by the enforcement scan: they are too young for the manager
  ///      to have seen them. A head older than the grace is a hard stop: it was already in place
  ///      at every possible read instant within the window, so no honest pick can sit past it in
  ///      scan order (`RequestQueues.selectionConflict`). Exempting only the *first* selected
  ///      queue instead would let a colluding submitter reopen the window each rotation and jump
  ///      the cursor past a starving application by serving one beyond it.
  ///
  ///      The comparison is against the wall clock, deliberately, and not against the submitted
  ///      request's own timestamp. The latter looks equivalent but is not: because FIFO forces the
  ///      manager to serve its application's *oldest* request, "the competitor is younger than my
  ///      pick" holds for competitors that arrived long before any selection view was read, so the
  ///      exemption would widen with however long the submitted request had been queued and the
  ///      rule would decay into approximate oldest-first. Against the clock, the exemption is
  ///      bounded by `selectionGrace` whatever the queue's history.
  ///
  ///      Starvation is bounded as a result: every application the scan allows sits at or before
  ///      the first aged head, so serving one moves the cursor toward it, never past it — its turn
  ///      arrives within one rotation.
  ///
  ///      Trigger and deploy requests bypass the per-application cursor by design, but their
  ///      *precedence* — triggers before deploys before per-application work — is enforced under
  ///      the same grace rule (`PriorityQueueNotServed`): an aged TRUSTPROCESS head blocks deploy
  ///      and per-application updates, and an aged DEPLOYAPP head blocks per-application updates.
  ///      A TRUSTPROCESS update is never constrained. The blast radius of a permanently failing
  ///      head widens accordingly: aged and unservable in a global queue, it halts every
  ///      lower-priority update until `selectionGrace` is raised past its age (section 7.4 of
  ///      `docs/design/BATCH_EXECUTION.md`).
  function _enforceSelection(uint64 applicationId, Structs.RequestType requestType) private view {
    if (requestType == Structs.RequestType.TRUSTPROCESS) return;

    uint256 grace = selectionGrace;
    if (RequestQueues.isHeadAged(_queueStore, _queueStore.triggers, grace)) {
      revert PriorityQueueNotServed(Structs.RequestType.TRUSTPROCESS);
    }

    if (requestType == Structs.RequestType.DEPLOYAPP) return;

    if (RequestQueues.isHeadAged(_queueStore, _queueStore.deploys, grace)) {
      revert PriorityQueueNotServed(Structs.RequestType.DEPLOYAPP);
    }

    // The scan always encounters the submitted application — the request sits at the head of a
    // non-empty per-application queue — so "not found" cannot mask a conflict.
    (uint64 conflictingAppId, bool conflict) = RequestQueues.selectionConflict(
      _queueStore,
      _deployedAppIds,
      applicationId,
      grace
    );
    if (conflict) revert ApplicationNotSelected(conflictingAppId);
  }

  /// @dev Removes the head of the request's own queue. The caller must have established that the
  ///      request *is* that head (`stateUpdate` does so through the same `_queueOf`);
  ///      `RequestQueues.dequeueHead` does not re-check the id.
  function _removeRequest(uint64 applicationId, Structs.RequestType requestType) private {
    RequestQueues.dequeueHead(_queueStore, _queueOf(applicationId, requestType));
    // Only per-application queues take part in the round robin. The application has had its
    // turn: move the cursor just past it so the next selection starts from the following
    // application. Trigger- and deploy-queue processing bypasses the cursor, leaving the
    // rotation where it paused.
    if (
      requestType != Structs.RequestType.TRUSTPROCESS &&
      requestType != Structs.RequestType.DEPLOYAPP
    ) {
      RequestQueues.advanceCursor(_queueStore, _deployedAppIds, applicationId);
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
    _removeRequest(applicationId, requestType);

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
    return RequestQueues.size(_queueStore.triggers);
  }

  /// @inheritdoc IProcessorEndpoint
  function getTriggerRequests() external view returns (Structs.PendingRequest[] memory) {
    return _copyRange(_queueStore.triggers, RequestQueues.size(_queueStore.triggers));
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequests() external view returns (Structs.PendingRequest[] memory) {
    return _pendingRequestsPage(0, type(uint256).max);
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsPage(
    uint256 offset,
    uint256 limit
  ) external view returns (Structs.PendingRequest[] memory) {
    return _pendingRequestsPage(offset, limit);
  }

  /// @inheritdoc IProcessorEndpoint
  function requestById(bytes32 id) external view returns (Structs.PendingRequest memory) {
    return _queueStore.requests[id];
  }

  /// @dev Flattens every pending request outside the trigger queue — the deploy queue first (it
  ///      is served first, see _selectPendingRequests), then each application's queue in
  ///      _deployedAppIds order — and returns the [offset, offset + limit) window of that list.
  ///      Requests are FIFO-ordered within a queue, but submission order across applications is
  ///      not preserved: the applications are independent.
  function _pendingRequestsPage(
    uint256 offset,
    uint256 limit
  ) internal view returns (Structs.PendingRequest[] memory result) {
    Structs.PendingRequest[] memory all = new Structs.PendingRequest[](_pendingRequestsSize());
    uint256 n = _copyInto(_queueStore.deploys, all, 0);
    uint256 appCount = _deployedAppIds.length;
    uint256 a;
    while (a != appCount) {
      n = _copyInto(_queueStore.pending[_deployedAppIds[a]], all, n);
      unchecked {
        ++a;
      }
    }

    if (offset >= n || limit == 0) return new Structs.PendingRequest[](0);
    uint256 count = n - offset;
    if (count > limit) count = limit;
    result = new Structs.PendingRequest[](count);
    uint256 j;
    while (j != count) {
      result[j] = all[offset + j];
      unchecked {
        ++j;
      }
    }
  }

  //update status
  /// @inheritdoc IProcessorEndpoint
  /// @dev Thin wrapper around `_processOneStateUpdate`, which is also what `batchStateUpdate`
  ///      loops over: there are two entry points, not two implementations. Kept alongside the
  ///      batch entry point because it is what the whole pre-batch test suite exercises, so it is
  ///      the proof that extracting `_processOneStateUpdate` did not change any behaviour.
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
    bytes32 currentRoot = applicationStateRoots[applicationId];
    bytes32 newRoot = _processOneStateUpdate(
      Structs.SignatureParams({
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
      }),
      currentRoot,
      true,
      true,
      signature
    );
    // The error path leaves the root untouched and returns it unchanged, which is exactly when no
    // write must happen — the success path always moves it (an unchanged root reverts there).
    if (newRoot != currentRoot) applicationStateRoots[applicationId] = newRoot;
  }

  /// @inheritdoc IProcessorEndpoint
  function batchStateUpdate(
    uint64 applicationId,
    Structs.BatchEntry[] calldata entries,
    bytes calldata batchSignature
  ) external onlyRole(UPDATE_STATUS_ROLE) nonReentrant {
    uint256 n = entries.length;
    // Declared here rather than left to the authenticator's own EmptyBatch: an error the endpoint
    // does not declare reaches clients as undecodable revert data.
    if (n == 0) revert EmptyBatch();

    // One batch belongs to one application, but three request kinds must stay unbatched: a
    // TRUSTPROCESS or DEPLOYAPP head (both live in a global queue, and both are processed one at a
    // time — see _selectPendingRequests), and any request of an application with a registered
    // trigger. For those the entry array must hold exactly one entry, which makes this function a
    // strict superset of stateUpdate rather than a second path with its own rules.
    //
    // Why they cannot be batched: the trigger contract is invoked *during* the update and can
    // enqueue a TRUSTPROCESS request whose payload is derived on-chain, so it cannot exist when
    // the TEE builds the batch. Entries computed after it in the same transaction would silently
    // observe a state the trigger flow never intended — a semantic violation that does not revert
    // (section 5.2 of docs/design/BATCH_EXECUTION.md).
    //
    // DEPLOYAPP is enum value 0, so an unknown requestId reads back as a zeroed request and lands
    // here as a deploy: a multi-entry batch whose first entry names a request that does not exist
    // reverts BatchNotAllowed rather than InvalidRequestId. Both reject the call; only the reason
    // is less precise.
    if (n != 1) {
      Structs.RequestType headType = _queueStore
        .requests[entries[0].processedRequestId]
        .requestType;
      if (
        headType == Structs.RequestType.TRUSTPROCESS ||
        headType == Structs.RequestType.DEPLOYAPP ||
        address(triggerContracts[applicationId]) != address(0)
      ) revert BatchNotAllowed();
    }

    // Build every entry hash first: the batch signature covers all of them at once, and it is
    // verified before any entry is processed. Params are kept alongside the hashes so the calldata
    // is copied into memory once rather than per loop.
    bytes32[] memory entryHashes = new bytes32[](n);
    Structs.SignatureParams[] memory params = new Structs.SignatureParams[](n);
    uint256 i;
    while (i != n) {
      Structs.SignatureParams memory p = Structs.SignatureParams({
        applicationId: applicationId,
        prevStateRoot: entries[i].prevStateRoot,
        newStateRoot: entries[i].newStateRoot,
        processedRequestId: entries[i].processedRequestId,
        userEvents: entries[i].userEvents,
        appEvents: entries[i].appEvents,
        withdrawalRequests: entries[i].withdrawalRequests,
        refundAmount: entries[i].refund,
        applicationFee: entries[i].applicationFees,
        errorCode: entries[i].errorCode,
        errorMsg: entries[i].errorMsg
      });
      params[i] = p;
      entryHashes[i] = UpdateEntryHash.entryHash(p);
      unchecked {
        ++i;
      }
    }

    // One ecrecover for the whole batch, over the personal_sign digest of the concatenated entry
    // hashes with its dynamic 32*N length prefix (see ITeeAuthenticator.checkBatchSignature).
    if (!teeAuthenticator.checkBatchSignature(entryHashes, batchSignature))
      revert InvalidSignature();

    // Chain the state root through the entries in memory: the first entry is validated against
    // storage, each later one against its predecessor's newStateRoot, and storage is written once
    // after the loop.
    bytes32 currentRoot = applicationStateRoots[applicationId];
    bytes32 root = currentRoot;
    i = 0;
    while (i != n) {
      // The turn is enforced on the first entry only: every entry shares one applicationId, and
      // the entries above that would need a different priority tier are limited to a lone entry.
      root = _processOneStateUpdate(params[i], root, i == 0, false, batchSignature);
      unchecked {
        ++i;
      }
    }
    if (root != currentRoot) applicationStateRoots[applicationId] = root;

    // The per-entry hashes are not recoverable from the individual events, so emit them: they are
    // what an off-chain verifier needs to re-derive the batch digest this signature was checked
    // against.
    emit BatchProcessed(applicationId, entryHashes);
  }

  /// @dev Processes one update entry: validates it, emits its events, moves its funds, invokes the
  ///      application's trigger if one is registered, and dequeues the request. The single
  ///      implementation behind both `stateUpdate` and `batchStateUpdate`.
  /// @param p The entry's fields, in the shape its hash is built from.
  /// @param currentRoot The application's state root as of this entry — read from storage for the
  ///        first entry of a call, the previous entry's `newStateRoot` afterwards. Deliberately not
  ///        read from storage here, so a batch can chain entries in memory and write the root once.
  ///        The one exception is an application with a registered trigger: there the root is
  ///        written before the trigger runs, because the trigger is external code that observes the
  ///        endpoint's state mid-transaction. Such applications are capped at one entry per batch,
  ///        so this costs the batch path nothing.
  /// @param enforceTurn Whether to enforce the selection rules (`_enforceSelection`). Callers pass
  ///        true for the first entry only: one check covers a whole batch.
  /// @param verifySignature Whether to verify `signature` as a 1-entry batch here. False for batch
  ///        entries, which one signature over every entry hash already covers — `signature` is then
  ///        ignored. Never derived from `signature` being empty: an empty signature must still fail
  ///        verification rather than skip it.
  /// @param signature Signature over this entry alone, ignored when `verifySignature` is false.
  /// @return The application's state root after this entry: `currentRoot` for an error entry,
  ///         `p.newStateRoot` for a successful one.
  function _processOneStateUpdate(
    Structs.SignatureParams memory p,
    bytes32 currentRoot,
    bool enforceTurn,
    bool verifySignature,
    bytes calldata signature
  ) private returns (bytes32) {
    uint64 applicationId = p.applicationId;
    bytes32 processedRequestId = p.processedRequestId;

    // Check valid request. The stored type resolves the one queue this request can be in (see
    // _queueOf), so this is a single head comparison rather than a scan of all three queues —
    // and it is the same queue _removeRequest will dequeue from.
    Structs.PendingRequest storage requestInfo = _queueStore.requests[processedRequestId];
    if (
      !RequestQueues.isHead(
        _queueOf(requestInfo.applicationId, requestInfo.requestType),
        processedRequestId
      )
    ) revert InvalidRequestId();

    bool fromTriggerQueue = requestInfo.requestType == Structs.RequestType.TRUSTPROCESS;

    // Check application Id
    if (applicationId != requestInfo.applicationId) revert InvalidApplicationId();

    // Enforce whose turn it is. Checked before the signature so a rejected turn costs the manager
    // no ecrecover, and before any state is touched.
    if (enforceTurn) _enforceSelection(applicationId, requestInfo.requestType);

    //check prev state root
    if (p.prevStateRoot != currentRoot) revert InvalidStateRoot();

    uint256 eventsLength = p.userEvents.events.length;
    if (eventsLength != p.userEvents.subTypes.length) revert InvalidPayload();

    uint256 appEventsLength = p.appEvents.events.length;
    if (appEventsLength != p.appEvents.subTypes.length) revert InvalidPayload();

    //check signature
    if (verifySignature) {
      // Single-request updates are verified as a 1-entry batch: the batch digest of one
      // entry hash is byte-identical to the single-request digest, so both submission
      // paths share one signing scheme.
      bytes32[] memory entryHashes = new bytes32[](1);
      entryHashes[0] = UpdateEntryHash.entryHash(p);
      if (!teeAuthenticator.checkBatchSignature(entryHashes, signature)) revert InvalidSignature();
    }

    //check values

    uint256 maxFeeValue = requestInfo.maxFeeValue;
    address payable sender = payable(requestInfo.sender);
    // Fee refund recipient: facilitator if present, otherwise sender
    address payable feeRecipient = requestInfo.facilitator != address(0)
      ? payable(requestInfo.facilitator)
      : sender;

    // Handle error case (signed error payload from TEE)
    if (p.errorCode != Structs.ErrorCode.NO_ERROR) {
      // For errors: state unchanged (prevStateRoot == newStateRoot), no events, no withdrawals
      // Refund user (minus minimum fee) and collect minimum fee
      if (eventsLength != 0 || appEventsLength != 0 || p.withdrawalRequests.length != 0)
        revert InvalidPayload();
      if (currentRoot != p.newStateRoot) revert InvalidStateRoot();

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
        p.errorCode,
        p.errorMsg,
        requestInfo.requestType
      );

      // State unchanged: the caller must not write the root for this entry, and a batch continues
      // from the same root.
      return currentRoot;
    }

    // Handle success case
    // State cannot remain the same
    if (currentRoot == p.newStateRoot) revert InvalidStateRoot();

    // don't check fees if we are from trigger queue
    if (!fromTriggerQueue) {
      if (p.refundAmount + p.applicationFee != maxFeeValue) revert InvalidValue();
      if (p.applicationFee < minFeePerRequest) {
        revert InvalidValue();
      }
    }

    //check withdrawal sums and debit per-app per-token custody
    uint256 i;
    Structs.WithdrawalRequest[] memory withdrawalRequests = p.withdrawalRequests;
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
    uint256 totalEthOutflow = p.refundAmount + p.applicationFee + ethWithdrawalSum;
    if (totalEthOutflow > _getAvailableEthBalance()) revert InsufficientBalance();

    //emit encrypted event
    i = 0;
    while (i != eventsLength) {
      emit UserEvent(
        applicationId,
        processedRequestId,
        p.userEvents.subTypes[i],
        p.userEvents.events[i]
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
        p.appEvents.subTypes[i],
        p.appEvents.events[i]
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

    //update request
    if (reqType == Structs.RequestType.DEPLOYAPP) {
      _deployedAppIds.push(applicationId);
      // The trigger (if any) was already validated and registered eagerly in
      // submitDeployRequestWithTrigger, so a successful deploy needs no further
      // trigger bookkeeping here.
    }
    emit StateRootUpdate(applicationId, processedRequestId, p.prevStateRoot, p.newStateRoot);

    //credit refund to feeRecipient's pending balance (pull pattern) — refund is always ETH
    if (p.refundAmount > 0) {
      _asyncTransfer(ETH_TOKEN, feeRecipient, p.refundAmount);
      emit Refund(applicationId, processedRequestId, feeRecipient, ETH_TOKEN, p.refundAmount);
    }

    //credit withdrawals to receivers' pending balances
    i = 0;
    uint256 insertIntoClaimable;
    ITrigger triggerContract = triggerContracts[applicationId];
    address trigger = address(triggerContract);
    // The state root normally reaches storage once per call, after the last entry. A registered
    // trigger is the exception: _invokeTrigger below hands control to external code that can read
    // this contract's state, and it must not observe a root the update has already moved past.
    if (trigger != address(0)) applicationStateRoots[applicationId] = p.newStateRoot;
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
    _invokeTrigger(triggerContract, applicationId, processedRequestId, p.appEvents, claimable);

    //set requests as completed
    _markRequestCompleted(
      applicationId,
      processedRequestId,
      p.applicationFee,
      Structs.RequestResult.COMPLETED,
      Structs.ErrorCode.NO_ERROR,
      '',
      reqType
    );

    return p.newStateRoot;
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function updateQueueThreshold(uint256 newThreshold) external {
    newThreshold;
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function updateMaxNumOfApplications(uint256 newMax) external {
    newMax;
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function updateSelectionGrace(uint256 newGrace) external {
    newGrace;
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function updateFeeCollector(address payable newFeeCollector) external {
    newFeeCollector;
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function addAllowedDeployer(address deployer) external {
    deployer;
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function removeAllowedDeployer(address deployer) external {
    deployer;
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  function isAllowedDeployer(address deployer) external view returns (bool allowed) {
    return hasRole(DEPLOYER_ROLE, deployer);
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsWithStateRoot(
    uint256 maxCount
  )
    external
    view
    returns (uint64 applicationId, Structs.PendingRequest[] memory requests, bytes32 stateRoot)
  {
    return _selectPendingRequests(maxCount);
  }

  /// @dev Picks the requests to serve next, in precedence order:
  ///      1. the trigger queue head, alone — a TRUSTPROCESS must be processed immediately after
  ///         the request that created it, before any other pending request of any application;
  ///      2. the deploy queue head, alone — deploys are a separate, never-batched flow;
  ///      3. up to maxCount requests of the application selected by the round-robin cursor,
  ///         capped at one when the application has a registered trigger contract (such
  ///         applications do not support batching).
  ///      Returns an empty list when nothing is pending. The cursor is not moved here: it only
  ///      advances when a request is actually dequeued from a per-application queue.
  function _selectPendingRequests(
    uint256 maxCount
  )
    internal
    view
    returns (uint64 applicationId, Structs.PendingRequest[] memory requests, bytes32 stateRoot)
  {
    if (maxCount != 0) {
      if (RequestQueues.size(_queueStore.triggers) > 0) return _headAsBatch(_queueStore.triggers);
      if (RequestQueues.size(_queueStore.deploys) > 0) return _headAsBatch(_queueStore.deploys);

      (uint64 appId, bool found) = RequestQueues.selectApplication(_queueStore, _deployedAppIds);
      if (found) {
        uint256 count = RequestQueues.size(_queueStore.pending[appId]);
        if (count > maxCount) count = maxCount;
        if (address(triggerContracts[appId]) != address(0)) count = 1;
        return (appId, _copyRange(_queueStore.pending[appId], count), applicationStateRoots[appId]);
      }
    }
    return (0, new Structs.PendingRequest[](0), bytes32(0));
  }

  /// @dev Returns the queue head as a single-element batch, together with its application's
  ///      state root.
  function _headAsBatch(
    RequestQueues.Queue storage q
  ) internal view returns (uint64, Structs.PendingRequest[] memory, bytes32) {
    Structs.PendingRequest[] memory requests = _copyRange(q, 1);
    uint64 applicationId = requests[0].applicationId;
    return (applicationId, requests, applicationStateRoots[applicationId]);
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev An unknown requestId reads back as a zeroed request and resolves to the deploy queue,
  ///      where the head comparison fails — so it returns false, as it does for a request that is
  ///      queued but not at its queue's head.
  function isCurrentPendingRequest(bytes32 requestId) public view returns (bool) {
    Structs.PendingRequest storage request = _queueStore.requests[requestId];
    return RequestQueues.isHead(_queueOf(request.applicationId, request.requestType), requestId);
  }

  // Pull payment pattern functions. `_asyncTransfer` is declared in `ProcessorEndpointStorage`,
  // because the reset entry points in `ProcessorEndpointExtension` credit refunds through it too.

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

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`, together with the queue
  ///      draining and deployed-app bookkeeping it needs.
  function adminReset() external {
    _delegateToExtension();
  }

  /// @inheritdoc IProcessorEndpoint
  /// @dev Forwarding stub; implemented in `ProcessorEndpointExtension`.
  function adminResetApps(uint64[] calldata appIds) external {
    appIds;
    _delegateToExtension();
  }

  /// @dev The first `count` requests of the queue, in FIFO order. `count` must not exceed the
  ///      queue size. Kept in this contract rather than in RequestQueues: as a library
  ///      `internal` function it gets inlined at every call site, and as a `public` one the
  ///      returned array would be ABI-encoded for the delegatecall and decoded again here —
  ///      both cost more code than they move out.
  function _copyRange(
    RequestQueues.Queue storage q,
    uint256 count
  ) internal view returns (Structs.PendingRequest[] memory result) {
    result = new Structs.PendingRequest[](count);
    uint256 i = q.head;
    uint256 stop = q.head + count;
    uint256 j;
    while (i != stop) {
      result[j] = _queueStore.requests[q.idByOrder[i]];
      unchecked {
        ++i;
        ++j;
      }
    }
  }

  /// @dev Appends the whole queue into `dest` starting at `offset`, returning the next free
  ///      index. Used to flatten several queues into one array.
  function _copyInto(
    RequestQueues.Queue storage q,
    Structs.PendingRequest[] memory dest,
    uint256 offset
  ) internal view returns (uint256) {
    uint256 i = q.head;
    uint256 tail = q.tail;
    while (i != tail) {
      dest[offset] = _queueStore.requests[q.idByOrder[i]];
      unchecked {
        ++i;
        ++offset;
      }
    }
    return offset;
  }

  // Calls execute then withdraw on the trigger contract registered for the given applicationId,
  // if any. Each call is wrapped in an independent try/catch so that a revert in the trigger
  // never propagates to the caller: both calls are always attempted regardless of the outcome
  // of the first.
  //
  // The trigger is passed in rather than read here: the caller loads it anyway, to decide whether
  // the state root has to reach storage before this external code runs.
  function _invokeTrigger(
    ITrigger trigger,
    uint64 applicationId,
    bytes32 processedRequestId,
    Structs.EventData memory appEventData,
    address[] memory claimable
  ) internal {
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

    emit RequestSubmitted(applicationId, requestId, trigger, address(0));
  }
}
