// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './IProcessorEndpoint.sol';
import './ITokenAllowlist.sol';
import '../Structs.sol';

/// @title ITrigger
/// @notice Interface for trigger contracts invoked by the ProcessorEndpoint on request completion.
/// @dev This interface MUST NOT be implemented directly. Trigger contracts must ALWAYS extend the
///      abstract contract AbstractTrigger, which provides the non-overridable withdraw sweep and the
///      _onlyProcessorEndpoint access control, and exposes the overridable _execute and
///      _getTrustProcessPayload hooks. Implementing ITrigger directly would bypass those guarantees.
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
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  function execute(Structs.EventData calldata appEventData) external;

  /// @notice Sweeps all allowlisted ERC-20 tokens and ETH held by this contract to the
  ///         ProcessorEndpoint. Called by the ProcessorEndpoint after having called execute.
  ///         Performs ONLY the sweep: producing the optional trusted follow-up payload is a
  ///         separate, explicit step (see getTrustProcessPayload).
  /// @return returnedTokens Array of tokens transferred in the sweep, with amounts.
  /// @return failedTokens Array of tokens that failed to transfer, with amounts.
  function withdraw()
    external
    returns (
      Structs.TokenAndAmount[] memory returnedTokens,
      Structs.TokenAndAmount[] memory failedTokens
    );

  /// @notice Produces the optional trusted (TRUSTPROCESS) follow-up payload. Called explicitly
  ///         by the ProcessorEndpoint after execute and withdraw, in its own isolated try/catch:
  ///         a revert here cannot block stateUpdate and simply results in no trusted request.
  ///         It is invoked even when withdraw failed, so the application can react to a failed
  ///         sweep (damage control) by returning an appropriate payload.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param executeSuccess True if the execute call succeeded, false if it reverted.
  /// @param withdrawSuccess True if the withdraw sweep succeeded, false if it reverted.
  /// @param returnedTokens Tokens swept back to the ProcessorEndpoint in withdraw (empty if it reverted).
  /// @param failedTokens Tokens that failed to transfer during the sweep (empty if withdraw reverted).
  /// @return trustedPayload Non-empty to enqueue a trusted (TRUSTPROCESS) request; empty means none.
  function getTrustProcessPayload(
    Structs.EventData calldata appEventData,
    bool executeSuccess,
    bool withdrawSuccess,
    Structs.TokenAndAmount[] calldata returnedTokens,
    Structs.TokenAndAmount[] calldata failedTokens
  ) external returns (bytes memory trustedPayload);
}
