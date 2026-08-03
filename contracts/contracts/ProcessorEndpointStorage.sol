// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts/utils/ReentrancyGuard.sol';
import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import '@openzeppelin/contracts/utils/cryptography/EIP712.sol';
import '@openzeppelin/contracts/utils/Strings.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IProcessorEndpoint.sol';
import './interfaces/IProcessorEndpointState.sol';
import './interfaces/IAuthorityRegistry.sol';
import './interfaces/ITokenAllowlist.sol';
import './interfaces/ITrigger.sol';
import './Structs.sol';

/// @title ProcessorEndpointStorage
/// @notice Storage layout and shared helpers for `ProcessorEndpoint` and
///         `ProcessorEndpointExtension`.
/// @dev `ProcessorEndpoint` is close to the 24,576-byte EIP-170 deployed-bytecode limit, so parts
///      of it live in `ProcessorEndpointExtension`, which the endpoint reaches by `delegatecall`.
///      Both contracts derive from this base, which is the single declaration of every state
///      variable: because a `delegatecall` executes the extension's code against the endpoint's
///      storage, the two layouts must be **identical**, and inheriting one common base is what
///      guarantees that. `AccessControl`, `ReentrancyGuard` and `EIP712` are inherited here, in
///      this order, for the same reason — they contribute storage slots (`_roles`, `_status`, the
///      EIP-712 name/version fallbacks) that must sit at the same offsets on both sides.
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
abstract contract ProcessorEndpointStorage is
  AccessControl,
  ReentrancyGuard,
  EIP712,
  IProcessorEndpointState
{
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
  uint64[] internal _deployedAppIds;
  uint256 public maxNumOfApplications = 10;
  uint256 public availableDeploySlots = maxNumOfApplications;

  RequestQueue internal _requestQueue;
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
  RequestQueue internal _triggerQueue;

  modifier validProtocolVersion(uint8 protocolVersion) {
    if (protocolVersion != PROTOCOL_VERSION) revert IProcessorEndpoint.InvalidProtocolVersion();
    _;
  }

  modifier validApplicationId(uint64 applicationId) {
    if (applicationStateRoots[applicationId] == bytes32(0))
      revert IProcessorEndpoint.InvalidApplicationId();
    _;
  }

  /// @dev Only fixes the EIP-712 domain, which must be identical on both sides so that a
  ///      signature verified in the extension matches one produced against the endpoint.
  ///      OpenZeppelin's `EIP712` caches the domain separator against the deploying address and
  ///      recomputes it whenever `address(this)` differs, so the extension resolves the endpoint's
  ///      domain when reached by `delegatecall`. All other initialisation belongs to
  ///      `ProcessorEndpoint`'s constructor: the extension holds no state of its own.
  constructor() EIP712('Vela', Strings.toString(PROTOCOL_VERSION)) {}

  // Queue primitives, declared alongside the RequestQueue state they operate on.

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

  /// @dev Internal implementation of `ProcessorEndpoint.getPendingRequestsSize`. The public
  ///      function stays on the endpoint; see the contract-level note on override conflicts.
  function _pendingRequestsSize() internal view returns (uint256) {
    return _queueSize(_requestQueue);
  }

  /// @dev Internal implementation of `ProcessorEndpoint.generateRequestId`.
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

    emit IProcessorEndpoint.RequestSubmitted(applicationId, requestId, sender, facilitator);
    return requestId;
  }
}
