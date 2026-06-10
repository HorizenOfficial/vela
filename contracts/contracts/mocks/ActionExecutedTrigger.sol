// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../interfaces/IProcessorEndpoint.sol';
import '../Structs.sol';
import '../trigger/AbstractTrigger.sol';

/// @title ActionExecutedTrigger
/// @notice Trigger contract that converts a WASM EXECUTE_REQUESTED AppEvent into a
///         trusted action_executed payload (TRUSTPROCESS), completing the execute →
///         action_executed loop automatically on-chain.
///
/// @dev Wire formats (LOCKED):
///   AppEvent.data (input):  ABI-encoded ExecuteRequestedEvent
///     (bytes16 lockId, bytes transaction, (address token, uint256 amount)[] tokens)
///   Trusted payload (output): ABI-encoded ActionExecutedTrusted
///     (bytes16 lockId, (address token, uint256 amount)[] remain, uint8 outcome)
///     outcome: 0 = Success, 1 = Failure
///
/// Policy: remain == the locked tokens (full pass-through), outcome == Success (0).
/// This is a test/mock trigger; a production trigger would inspect on-chain results.
contract ActionExecutedTrigger is AbstractTrigger {
  /// @param _processorEndpoint ProcessorEndpoint that will call execute and withdraw.
  constructor(IProcessorEndpoint _processorEndpoint) AbstractTrigger(_processorEndpoint) {}

  /// @notice No-op execute hook; all logic lives in getTrustProcessPayload.
  function _execute(Structs.EventData calldata) internal override {}

  /// @notice Decodes the WASM EXECUTE_REQUESTED AppEvent and produces the trusted
  ///         action_executed payload that the WASM handler_trusted will ABI-decode.
  ///
  /// @dev ABI-decodes appEventData.events[0] as the locked ExecuteRequestedEvent tuple,
  ///      then ABI-encodes the ActionExecutedTrusted tuple with:
  ///        - lockId:  extracted from the AppEvent
  ///        - remain:  the locked tokens (full pass-through policy)
  ///        - outcome: 0 (Success)
  ///
  /// `Structs.TokenAndAmount` has the same ABI layout as the WASM `TokenAmount`
  /// (address token, uint256 amount), so it decodes/encodes the WASM token list directly.
  ///
  /// Returns non-empty bytes so the ProcessorEndpoint enqueues a TRUSTPROCESS request.
  function _getTrustProcessPayload(
    Structs.EventData calldata appEventData,
    bool /*executeSuccess*/,
    bool /*withdrawSuccess*/,
    Structs.TokenAndAmount[] calldata /*returned*/,
    Structs.TokenAndAmount[] calldata /*failed*/
  ) internal pure override returns (bytes memory) {
    // Guard: only produce a trusted payload when the AppEvent array is non-empty
    // (TRUSTPROCESS stateUpdates have no AppEvents and must not trigger again).
    if (appEventData.events.length == 0) {
      return '';
    }

    // Decode the WASM EXECUTE_REQUESTED AppEvent data (ABI-encoded ExecuteRequestedEvent):
    //   (bytes16 lockId, bytes transaction, (address token, uint256 amount)[] tokens)
    (bytes16 lockId, , Structs.TokenAndAmount[] memory tokens) = abi.decode(
      appEventData.events[0],
      (bytes16, bytes, Structs.TokenAndAmount[])
    );

    // Produce the trusted action_executed payload (ActionExecutedTrusted):
    //   (bytes16 lockId, (address token, uint256 amount)[] remain, uint8 outcome)
    // Policy: remain == the locked tokens, outcome == Success (0).
    return abi.encode(lockId, tokens, uint8(0));
  }
}
