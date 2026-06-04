// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './IProcessorEndpoint.sol';
import './ITokenAllowlist.sol';
import '../Structs.sol';

/// @title ITrigger
/// @notice Interface for trigger contracts invoked by the ProcessorEndpoint on request completion.
interface ITrigger {
  /// @notice Emitted when execute is called successfully.
  event TriggerExecuted();

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
  function execute(Structs.EventData calldata appEventData) external;

  /// @notice Sweeps all allowlisted ERC-20 tokens and ETH held by this contract to the
  ///         ProcessorEndpoint. Can only be called by the ProcessorEndpoint.
  /// @return postWithdrawSuccess True if _postWithdraw executed successfully, false if it reverted.
  /// @return returnedTokens Array of tokens transferred in the sweep, with amounts.
  /// @return failedTokens Array of tokens that failed to transfer, with amounts.
  function withdraw(
    Structs.EventData calldata appEventData,
    bool executeSuccess
  )
    external
    returns (
      bool postWithdrawSuccess,
      Structs.TokenAndAmount[] memory returnedTokens,
      Structs.TokenAndAmount[] memory failedTokens
    );
}
