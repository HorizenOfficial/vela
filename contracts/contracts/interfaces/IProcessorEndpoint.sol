// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../Structs.sol';
import './ITokenAllowlist.sol';

/// @title ProcessorEndpoint interface
/// @notice External API, events, and errors for the processor endpoint.
interface IProcessorEndpoint {
  /// @notice Emitted when a refund is sent to a requester.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param to Refund recipient.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param amount Refunded amount.
  event Refund(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    address indexed to,
    address tokenAddress,
    uint256 amount
  );
  /// @notice Emitted when a withdrawal is executed.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param to Withdrawal recipient.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param amount Withdrawal amount.
  event Withdrawal(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    address indexed to,
    address tokenAddress,
    uint256 amount
  );
  /// @notice Emitted when a new request enters the queue.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param sender Request sender.
  /// @param facilitator Facilitator address (address(0) for direct submissions).
  event RequestSubmitted(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    address indexed sender,
    address facilitator
  );
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
    bytes32 indexed eventSubType,
    bytes encryptedData
  );
  /// @notice Emitted for application-level (non-encrypted) events.
  /// @param applicationId Application identifier.
  /// @param requestId Request identifier.
  /// @param eventSubType Application event subtype.
  /// @param data Unencrypted payload.
  event AppEvent(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    bytes32 indexed eventSubType,
    bytes data
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
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param payee Address of the payee.
  /// @param amount Amount withdrawn.
  event PaymentWithdrawn(address tokenAddress, address indexed payee, uint256 amount);

  // @notice Emitted when a trigger's execute function is called.
  // @param applicationId Application identifier.
  // @param processedRequestId Request identifier being processed.
  // @param success True if the execute call succeeded, false if it reverted.
  event TriggerExecuted(
    uint64 indexed applicationId,
    bytes32 indexed processedRequestId,
    bool success
  );
  // @notice Emitted when a trigger's withdraw function is called.
  // @param applicationId Application identifier.
  // @param processedRequestId Request identifier being processed.
  // @param withdrawSuccess True if the withdraw call succeeded, false if it reverted.
  // @param postWithdrawSuccess True if the post-withdraw call succeeded, false if it reverted.
  // @param returnedTokens Array of tokens transferred in the sweep, with amounts.
  // @param failedTokens Array of tokens that failed to transfer, with amounts.
  event TriggerWithdraw(
    uint64 indexed applicationId,
    bytes32 indexed processedRequestId,
    bool withdrawSuccess,
    bool postWithdrawSuccess,
    Structs.TokenAndAmount[] returnedTokens,
    Structs.TokenAndAmount[] failedTokens
  );

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
  /// @notice The received ERC-20 amount does not match the expected amount.
  error TransferAmountMismatch();
  /// @notice The request authorization deadline has expired.
  error DeadlineExpired();
  /// @notice The recovered signer does not match the declared sender (used for facilitator requests).
  error InvalidSigner();
  /// @notice The deposit permit bytes are invalid or missing when required (used for facilitator requests).
  error InvalidPermit();
  /// @notice Trigger contract is already registered to another application.
  error TriggerAlreadyRegistered();
  /// @notice Trigger contrct cannot be an EOA.
  error TriggerCannotBeEOA();

  /// @notice Submits a new request and enqueues it for processing.
  /// @param protocolVersion Protocol version.
  /// @param applicationId Application identifier.
  /// @param requestType Request type.
  /// @param payload Request payload.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param assetAmount Business asset amount.
  /// @param maxFeeValue Maximum fee reserved for processing.
  /// @return requestId Generated request id.
  function submitRequest(
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 maxFeeValue
  ) external payable returns (bytes32);

  /// @notice Submits a request on behalf of a user (meta-transaction / facilitator pattern).
  /// @dev The facilitator (msg.sender) pays gas and maxFeeValue in ETH.
  ///      The user authorizes the request via an EIP-712 signature and the deposit via EIP-2612 permit.
  /// @param sender User address (must match recovered signer from requestSignature).
  /// @param protocolVersion Protocol version.
  /// @param applicationId Application identifier.
  /// @param requestType Request type (only ASSOCIATEKEY and PROCESS supported).
  /// @param payload Request payload.
  /// @param tokenAddress ERC-20 token address for business-asset deposit.
  /// @param assetAmount Business asset deposit amount (0 if no deposit).
  /// @param deadline Expiration timestamp for signatures.
  /// @param requestSignature EIP-712 request authorization signature from the user.
  /// @param depositPermit ABI-encoded EIP-2612 permit (v, r, s). Empty if assetAmount == 0.
  /// @return requestId Generated request id.
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
  ) external payable returns (bytes32);

  /// @notice Returns the current facilitator nonce for a user.
  /// @param user User address.
  /// @return nonce Current nonce value.
  function getFacilitatorNonce(address user) external view returns (uint256);

  /// @notice Submits a new deploy request and enqueues it for processing.
  /// @param protocolVersion Protocol version.
  /// @param payload Request payload.
  /// @return requestId Generated request id.
  function submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  ) external payable returns (bytes32);

  /// @notice Submits a new deploy request and enqueues it for processing,
  ///         optionally registering a trigger contract for the new application.
  ///         Unlike the legacy 32-byte-payload registration, this lets the
  ///         deploy carry a full descriptor payload AND a trigger address.
  /// @param protocolVersion Protocol version.
  /// @param payload Request payload (e.g. the WASM deploy descriptor).
  /// @param trigger Trigger contract to register for the deployed app, or
  ///        address(0) for no trigger (equivalent to the 2-arg overload).
  /// @return requestId Generated request id.
  function submitDeployRequestWithTrigger(
    uint8 protocolVersion,
    bytes calldata payload,
    address trigger
  ) external payable returns (bytes32);

  /// @notice Returns the number of pending requests in the queue.
  /// @return size Current pending request count.
  function getPendingRequestsSize() external view returns (uint256);

  /// @notice Returns the number of requests currently in the trigger queue.
  /// @return size Current trigger request count.
  function getTriggerQueueSize() external view returns (uint256);

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
  /// @param userEventData Encrypted user events and their subtypes.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
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
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
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
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param assetAmount Business asset amount.
  /// @param idx Queue index.
  /// @return requestId Derived request id.
  function generateRequestId(
    address sender,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 idx
  ) external pure returns (bytes32);

  /// @notice Returns the custody balance for a given application and token.
  /// @param applicationId Application identifier.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @return amount Current custody balance.
  function appCustody(uint64 applicationId, address tokenAddress) external view returns (uint256);

  /// @notice Returns the total custody balance across all applications for a given token.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @return amount Total custody balance.
  function totalAppCustody(address tokenAddress) external view returns (uint256);

  /// @notice Returns the pending claim balance for a given token and payee.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param payee Payee address.
  /// @return amount Pending claim balance.
  function pendingClaims(address tokenAddress, address payee) external view returns (uint256);

  /// @notice Returns the total pending claims for a given token.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @return amount Total pending claims.
  function totalPendingClaims(address tokenAddress) external view returns (uint256);

  /// @notice Claims pending balance for a given token and payee.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param payee Payee address.
  function claim(address tokenAddress, address payable payee) external;

  /// @notice Returns all application IDs that have ever been successfully deployed.
  /// @return appIds Array of deployed application IDs.
  function getDeployedAppIds() external view returns (uint64[] memory);

  /// @notice Clears the pending request queue and frees deploy slots for any pending DEPLOYAPP
  ///         requests that are removed. Restricted to RESET_OPERATOR. Unreachable when
  ///         RESET_OPERATOR was initialised as address(0).
  function adminReset() external;

  /// @notice Resets state roots and locked funds for the given application IDs, clears the queue
  ///         in the same transaction, and transfers all recovered ETH and ERC-20 balances to the
  ///         caller. Pass empty array to target all deployed apps.
  ///         Restricted to RESET_OPERATOR. Unreachable when RESET_OPERATOR was initialised as
  ///         address(0).
  /// @param appIds Application IDs to reset. Empty array means all deployed apps.
  function adminResetApps(uint64[] calldata appIds) external;

  /// @notice Returns the external token allowlist contract used by this endpoint.
  function tokenAllowlist() external view returns (ITokenAllowlist);
}
