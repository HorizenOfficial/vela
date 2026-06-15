// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../interfaces/ITokenAllowlist.sol';
import '../interfaces/ITrigger.sol';
import '../Structs.sol';

/// @title MockTriggerEndpoint
/// @notice Minimal ProcessorEndpoint stand-in for unit-testing trigger contracts without
///         deploying the full endpoint. Exposes tokenAllowlist() (read by AbstractTrigger's
///         constructor) and a forwarder so tests can invoke the ProcessorEndpoint-only trigger
///         entry points (e.g. getTrustProcessPayload) with msg.sender == this endpoint.
contract MockTriggerEndpoint {
  ITokenAllowlist public tokenAllowlist;

  constructor(ITokenAllowlist _tokenAllowlist) {
    tokenAllowlist = _tokenAllowlist;
  }

  /// @notice Forwards a getTrustProcessPayload call to a trigger so the trigger sees
  ///         msg.sender == this endpoint (satisfying its _onlyProcessorEndpoint guard).
  function callGetTrustProcessPayload(
    ITrigger trigger,
    Structs.EventData calldata appEventData,
    bool executeSuccess,
    bool withdrawSuccess,
    Structs.TokenAndAmount[] calldata returnedTokens,
    Structs.TokenAndAmount[] calldata failedTokens
  ) external returns (bytes memory) {
    return
      trigger.getTrustProcessPayload(
        appEventData,
        executeSuccess,
        withdrawSuccess,
        returnedTokens,
        failedTokens
      );
  }
}
