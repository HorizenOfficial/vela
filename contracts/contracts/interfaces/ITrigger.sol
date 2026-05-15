// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './IProcessorEndpoint.sol';
import './ITokenAllowlist.sol';
import '../Structs.sol';

/// @title ITrigger
/// @notice Interface for trigger contracts invoked by the ProcessorEndpoint on request completion.
interface ITrigger {
  /// @notice Describes a token transfer performed during a withdraw call.
  struct ReturnedTokens {
    /// @notice Token address; address(0) for ETH.
    address token;
    /// @notice Amount transferred to the ProcessorEndpoint.
    uint256 amount;
  }

  /// @notice Emitted when execute is called successfully.
  /// @param applicationId Application identifier.
  /// @param processedRequestId Request identifier that was processed.
  event TriggerExecuted(uint64 indexed applicationId, bytes32 indexed processedRequestId);

  /// @notice A zero address was supplied where not allowed.
  error ZeroAddress();
  /// @notice Caller is not the ProcessorEndpoint.
  error NotProcessorEndpoint();

  /// @notice ProcessorEndpoint that is allowed to call execute and withdraw.
  function processorEndpoint() external view returns (IProcessorEndpoint);

  /// @notice Allowlist of ERC-20 tokens swept on withdraw.
  function tokenAllowlist() external view returns (ITokenAllowlist);

  /// @notice Called by the ProcessorEndpoint when a request is completed.
  ///         Delegates to _execute and emits TriggerExecuted.
  function execute(
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
  ) external;

  /// @notice Sweeps all allowlisted ERC-20 tokens and ETH held by this contract to the
  ///         ProcessorEndpoint. Can only be called by the ProcessorEndpoint.
  /// @return success True if _postWithdraw executed successfully, false if it reverted.
  /// @return returnedTokens Array of tokens transferred in the sweep, with amounts.
  function withdraw(
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
  ) external returns (bool success, ReturnedTokens[] memory returnedTokens);
}
