// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../Structs.sol';

/// @title ProcessorEndpoint interface
/// @notice External API, events, and errors for the processor endpoint.
interface IProcessorEndpoint {
  /// @notice Emitted when a refund is sent to a requester.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param to Refund recipient.
  /// @param amount Refunded amount.
  event Refund(uint64 indexed applicationId, bytes32 indexed requestId, address to, uint256 amount);
  /// @notice Emitted when a withdrawal is executed.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param to Withdrawal recipient.
  /// @param amount Withdrawal amount.
  event Withdrawal(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    address to,
    uint256 amount
  );
  /// @notice Emitted when a new request enters the queue.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param sender Request sender.
  event RequestSubmitted(uint64 indexed applicationId, bytes32 requestId, address indexed sender);
  /// @notice Emitted when a new deploy request enters the queue.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param sender Request sender.
  event DeployRequestSubmitted(
    uint64 indexed applicationId,
    bytes32 requestId,
    address indexed sender
  );
  /// @notice Emitted when a request is finalized.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param applicationFees Fees collected for the application.
  /// @param status Completion status for the request.
  /// @param errorCode Error code when the request failed.
  /// @param errorMessage Human-readable error message.
  event RequestCompleted(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    uint256 applicationFees,
    Structs.RequestResult status,
    Structs.ErrorCode errorCode,
    string errorMessage
  );
  /// @notice Emitted when a deploy request is finalized.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param applicationFees Fees collected for the application.
  /// @param status Completion status for the request.
  /// @param errorCode Error code when the request failed.
  /// @param errorMessage Human-readable error message.
  event DeployRequestCompleted(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    uint256 applicationFees,
    Structs.RequestResult status,
    Structs.ErrorCode errorCode,
    string errorMessage
  );
  /// @notice Emitted when a reqreport has been generated.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  event ReportGenerated(uint64 indexed applicationId, bytes32 indexed requestId);
  /// @notice Emitted for application-specific encrypted events.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param eventSubType Application event subtype.
  /// @param encryptedData Encrypted payload.
  event UserEvent(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    string indexed eventSubType,
    bytes encryptedData
  );
  /// @notice Emitted when the state root is updated.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param oldStateRoot Previous state root.
  /// @param newStateRoot New state root.
  event StateRootUpdate(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    bytes32 oldStateRoot,
    bytes32 newStateRoot
  );
  /// @notice Emitted when the queue size threshold is changed.
  /// @param newThreshold New maximum queue size.
  event QueueThresholdUpdated(uint256 newThreshold);
  /// @notice Emitted when the maximum number of applications is updated.
  /// @param oldMax Previous maximum.
  /// @param newMax New maximum.
  event MaxNumberOfAppUpdated(uint256 oldMax, uint256 newMax);
  /// @notice Emitted when the fee collector address is updated.
  /// @param newFeeCollector New fee collector address.
  event FeeCollectorUpdated(address newFeeCollector);
  /// @notice Emitted when a payment is withdrawn.
  /// @param payee Address of the payee.
  /// @param amount Amount withdrawn.
  event PaymentWithdrawn(address indexed payee, uint256 amount);

  /// @notice A zero address was supplied where not allowed.
  error AddressCantBeZero();
  /// @notice Fee value is below the minimum allowed.
  error FeeValueBelowMinimum();
  /// @notice A numeric value is invalid for the requested operation.
  error InvalidValue();
  /// @notice The provided protocol version is unsupported.
  error InvalidProtocolVersion();
  /// @notice The provided application id is unsupported.
  error InvalidApplicationId();
  /// @notice The provided request id is not the current pending request.
  error InvalidRequestId();
  /// @notice The provided state root does not match the expected value.
  error InvalidStateRoot();
  /// @notice The signature could not be verified.
  error InvalidSignature();
  /// @notice The provided payload is invalid.
  error InvalidPayload();
  /// @notice Contract balance is insufficient for the requested operation.
  error InsufficientBalance();
  /// @notice Application's locked funds are insufficient for the requested operation.
  error InsufficientAppBalance();
  /// @notice Caller is not authorized to perform the requested action.
  error AuthorityNotAllowed();
  /// @notice Caller is not allowed to submit deploy requests.
  error DeployerNotAllowed();
  /// @notice Queue size exceeds the configured threshold.
  error QueueThresholdExceeded();
  /// @notice Maximum number of deployed applications has been reached.
  error MaxNumOfApplicationsExceeded();
  /// @notice An ETH transfer failed.
  error TransferFailed();
  /// @notice The provided request type is not allowed.
  error InvalidRequestType();

  /// @notice Submits a new request and enqueues it for processing.
  /// @param protocolVersion Protocol version.
  /// @param applicationId Application identifier.
  /// @param requestType Request type.
  /// @param payload Request payload.
  /// @param depositAmount Value forwarded for app logic.
  /// @param maxFeeValue Maximum fee reserved for processing.
  /// @return requestId Generated request id.
  function submitRequest(
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    uint256 depositAmount,
    uint256 maxFeeValue
  ) external payable returns (bytes32);

  /// @notice Submits a new deploy request and enqueues it for processing.
  /// @param protocolVersion Protocol version.
  /// @param payload Request payload.
  /// @return requestId Generated request id.
  function submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  ) external payable returns (bytes32);

  /// @notice Returns the number of pending requests in the queue.
  /// @return size Current pending request count.
  function getPendingRequestsSize() external view returns (uint256);

  /// @notice Returns the list of pending requests in order.
  /// @return requests Array of pending requests.
  function getPendingRequests() external view returns (Structs.PendingRequest[] memory);

  /// @notice Returns a paginated slice of pending requests in order.
  /// @param offset Number of requests to skip from the head of the queue.
  /// @param limit Maximum number of requests to return.
  /// @return requests Array of pending requests.
  function getPendingRequestsPage(
    uint256 offset,
    uint256 limit
  ) external view returns (Structs.PendingRequest[] memory);

  /// @notice Updates the state root, emits events, and finalizes the processed request.
  /// @param applicationId Application identifier.
  /// @param prevStateRoot Previous state root.
  /// @param newStateRoot New state root.
  /// @param processedRequestId Request identifier being processed.
  /// @param events Encrypted event payloads.
  /// @param eventSubTypes Event subtype labels.
  /// @param withdrawalRequests Withdrawal requests to execute.
  /// @param refund Refund amount to the request sender.
  /// @param applicationFees Fee amount to the collector.
  /// @param errorCode Error code for the update.
  /// @param errorMsg Error message for the update.
  /// @param signature Signature over the update data.
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
    Structs.ErrorCode errorCode,
    string calldata errorMsg,
    bytes calldata signature
  ) external;

  /// @notice Updates the maximum pending queue size.
  /// @param newThreshold New queue size limit.
  function updateQueueThreshold(uint256 newThreshold) external;

  /// @notice Updates the maximum number of deployable applications.
  /// @param newMax New maximum. Must be >= currently deployed count.
  function updateMaxNumOfApplications(uint256 newMax) external;

  /// @notice Updates the fee collector address.
  /// @param newFeeCollector New fee collector address.
  function updateFeeCollector(address payable newFeeCollector) external;

  /// @notice Grants deploy permission to an address.
  /// @param deployer Address to authorize for deploy requests.
  function addAllowedDeployer(address deployer) external;

  /// @notice Revokes deploy permission from an address.
  /// @param deployer Address to revoke for deploy requests.
  function removeAllowedDeployer(address deployer) external;

  /// @notice Returns whether an address can submit deploy requests.
  /// @param deployer Address to check.
  /// @return allowed True when the address has deploy permission.
  function isAllowedDeployer(address deployer) external view returns (bool allowed);

  /// @notice Returns the current pending request and state root.
  /// @return request Pending request data.
  /// @return currentStateRoot Current state root.
  /// @return success True if a pending request exists.
  function getNextPendingRequest()
    external
    view
    returns (Structs.PendingRequest memory, bytes32, bool);

  /// @notice Checks whether a request is the current pending request.
  /// @param requestId Request identifier.
  /// @return isCurrent True if the request is at the head of the queue.
  function isCurrentPendingRequest(bytes32 requestId) external view returns (bool);

  /// @notice Generates a deterministic request id.
  /// @param sender Sender address.
  /// @param applicationId Application identifier.
  /// @param requestType Request type.
  /// @param payload Request payload.
  /// @param depositAmount Value forwarded for app logic.
  /// @param idx Queue index.
  /// @return requestId Derived request id.
  function generateRequestId(
    address sender,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    uint256 depositAmount,
    uint256 idx
  ) external pure returns (bytes32);

  /// @notice Returns the locked funds for a given application (includes residual credit from prior requests).
  /// @param applicationId Application identifier.
  /// @return amount Current locked funds for the application.
  function appLockedFunds(uint64 applicationId) external view returns (uint256);

  /// @notice Withdraws pending payments for a given payee.
  /// @param payee Payee address.
  function withdrawPayments(address payable payee) external;
}
