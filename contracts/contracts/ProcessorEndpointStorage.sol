// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol';
import '@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol';
import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import '@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IProcessorEndpoint.sol';
import './interfaces/IProcessorEndpointState.sol';
import './interfaces/IAuthorityRegistry.sol';
import './interfaces/ITokenAllowlist.sol';
import './interfaces/ITrigger.sol';
import './Structs.sol';
import './RequestQueues.sol';

/// @title ProcessorEndpointStorage
/// @notice Storage layout and shared helpers for `ProcessorEndpoint` and
///         `ProcessorEndpointExtension`.
/// @dev `ProcessorEndpoint` is within ~1KB of the EIP-170 deployed-bytecode limit, so parts of it
///      live in `ProcessorEndpointExtension`, which the endpoint reaches by `delegatecall`. Both
///      contracts derive from this base, which is the single declaration of every state variable:
///      because a `delegatecall` executes the extension's code against the endpoint's storage, the
///      two layouts must be **identical**, and inheriting one common base is what guarantees that.
///      `AccessControlUpgradeable` and `EIP712Upgradeable` are inherited here, in this order, so
///      both contracts resolve the same namespaced storage for roles and the EIP-712 domain.
///
///      Consequences to respect when editing:
///      - Never declare state in `ProcessorEndpoint` or `ProcessorEndpointExtension`; declare it
///        here, and never reorder or remove a variable without re-checking both contracts.
///        `npm run check:layout` is the check, and it runs in CI.
///      - Errors and events are referenced through `IProcessorEndpoint.<name>`: this base cannot
///        inherit that interface, because the extension would then have to implement every
///        function in it.
///      - Helpers the extension needs are `internal` here. Where the endpoint also exposes them
///        externally (`getPendingRequestsSize`, `generateRequestId`), the public function stays in
///        `ProcessorEndpoint` and forwards to the `internal` implementation, so that inheriting
///        both this base and `IProcessorEndpoint` does not create an override conflict.
///      See `docs/design/PROCESSOR_ENDPOINT_SPLIT.md` for the full rationale.
///
/// @dev Upgradeability (see `docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`): `ProcessorEndpoint` is
///      deployed behind an ERC1967/UUPS proxy, so this base uses the `*Upgradeable` OpenZeppelin
///      variants, which keep their own storage in ERC-7201 namespaced slots rather than the
///      sequential slots plain `AccessControl`/`EIP712` used — the declared state below still
///      starts at slot 0. `ReentrancyGuardTransient` is used instead of a `ReentrancyGuardUpgradeable`
///      (no longer shipped by the OpenZeppelin upgradeable package): it guards via transient
///      storage (EIP-1153), which is inherently proxy-safe and needs no initialization. Both
///      `ProcessorEndpoint` and `ProcessorEndpointExtension` read the same EIP-712 domain: only
///      `ProcessorEndpoint.initialize` calls `__EIP712_init`, since `EIP712Upgradeable` stores
///      name/version in the shared namespaced slot rather than per-contract immutables. The
///      `__gap` reserves storage for future versions; see the design doc's storage-gap convention
///      before touching it.
abstract contract ProcessorEndpointStorage is
  AccessControlUpgradeable,
  ReentrancyGuardTransient,
  EIP712Upgradeable,
  IProcessorEndpointState
{
  using SafeERC20 for IERC20;

  //constants
  bytes32 public constant UPDATE_STATUS_ROLE = keccak256('UPDATE_STATUS_ROLE');
  bytes32 public constant ADMIN = keccak256('ADMIN');
  bytes32 public constant DEPLOYER_ROLE = keccak256('DEPLOYER_ROLE');
  bytes32 public constant RESET_OPERATOR = keccak256('RESET_OPERATOR');
  uint8 public constant PROTOCOL_VERSION = 0;
  //state variables
  mapping(uint64 => bytes32) public applicationStateRoots;
  uint64[] internal _deployedAppIds;
  uint256 public maxNumOfApplications;
  uint256 public availableDeploySlots;

  // All pending-request queue state: the global request store, the per-application queues, the
  // global deploy and trigger queues, the round-robin cursor and the total request count. Held
  // as one struct so it can be handed to RequestQueues as a single storage pointer.
  RequestQueues.Store internal _queueStore;
  uint256 public maxQueueSize;

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

  /// @notice Seconds for which a freshly enqueued request does not yet constrain what the
  ///         manager may serve — neither the round-robin turn among applications nor the
  ///         precedence of the global trigger and deploy queues. See
  ///         `ProcessorEndpoint._enforceSelection`: both rules are enforced on-chain, and this
  ///         window is what absorbs the requests that arrive between the manager reading the
  ///         selection view and its state update landing. Larger values weaken the fairness bound; a value large enough to cover any
  ///         request age disables enforcement entirely, which is the only escape hatch from a
  ///         permanently failing queue head short of `adminReset` (see section 7.4 of
  ///         `docs/design/BATCH_EXECUTION.md`). Zero means strict enforcement with no tolerance
  ///         for the race.
  uint256 public selectionGrace;

  /// @dev Reserved storage for future versions (see `docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`).
  ///      Must be reduced by the number of slots any new variable added above it consumes, and
  ///      must always remain the last declaration in this contract.
  uint256[50] private __gap;

  modifier validProtocolVersion(uint8 protocolVersion) {
    if (protocolVersion != PROTOCOL_VERSION) revert IProcessorEndpoint.InvalidProtocolVersion();
    _;
  }

  modifier validApplicationId(uint64 applicationId) {
    if (applicationStateRoots[applicationId] == bytes32(0))
      revert IProcessorEndpoint.InvalidApplicationId();
    _;
  }

  /// @dev Internal implementation of `ProcessorEndpoint.getPendingRequestsSize`. The public
  ///      function stays on the endpoint; see the contract-level note on override conflicts.
  function _pendingRequestsSize() internal view returns (uint256) {
    return _queueStore.total - RequestQueues.size(_queueStore.triggers);
  }

  function _generateRequestId(
    address sender,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes32 payloadHash,
    address tokenAddress,
    uint256 assetAmount,
    uint256 idx
  ) internal pure returns (bytes32) {
    return
      keccak256(
        abi.encode(sender, applicationId, requestType, payloadHash, tokenAddress, assetAmount, idx)
      );
  }

  function _pullERC20(address tokenAddress, address from, uint256 amount) internal {
    IERC20 token = IERC20(tokenAddress);
    uint256 balanceBefore = token.balanceOf(address(this));
    token.safeTransferFrom(from, address(this), amount);
    uint256 received = token.balanceOf(address(this)) - balanceBefore;
    if (received != amount) revert IProcessorEndpoint.TransferAmountMismatch();
  }

  // update custody helper
  function _addToCustody(uint64 applicationId, address token, uint256 amount) internal {
    appCustody[applicationId][token] += amount;
    totalAppCustody[token] += amount;
  }

  function _subtractToCustody(uint64 applicationId, address token, uint256 amount) internal {
    appCustody[applicationId][token] -= amount;
    totalAppCustody[token] -= amount;
  }

  /// @dev Credits `dest` in the pull-payment ledger. Declared here rather than on the endpoint
  ///      because both sides need it: `stateUpdate` for refunds and withdrawals, and the reset
  ///      entry points in `ProcessorEndpointExtension` for the deposits they return.
  function _asyncTransfer(address tokenAddress, address dest, uint256 amount) internal {
    pendingClaims[tokenAddress][dest] += amount;
    totalPendingClaims[tokenAddress] += amount;
  }

  /// @dev Removes any trigger registration for the given application (both the forward and
  ///      reverse mappings). No-op when no trigger is registered. Shared: `stateUpdate` rolls a
  ///      failed deploy's eager registration back, and the reset entry points in
  ///      `ProcessorEndpointExtension` clear registrations for the apps they discard.
  function _clearTrigger(uint64 applicationId) internal {
    ITrigger trigger = triggerContracts[applicationId];
    if (address(trigger) != address(0)) {
      delete triggersToAppIds[address(trigger)];
      delete triggerContracts[applicationId];
    }
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
    bytes32 requestId = _generateRequestId(
      sender,
      applicationId,
      requestType,
      keccak256(payload),
      tokenAddress,
      assetAmount,
      _queueStore.pending[applicationId].tail
    );
    RequestQueues.enqueue(
      _queueStore,
      _queueStore.pending[applicationId],
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

    emit IProcessorEndpoint.RequestSubmitted(applicationId, requestId, sender, facilitator);
    return requestId;
  }
}
