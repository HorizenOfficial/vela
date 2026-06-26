// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../trigger/AbstractTrigger.sol';
import '../Structs.sol';

/// @title Guarded Trigger
/// @notice AbstractTrigger implementation for e2e tests of the full
///         request -> trigger -> unshield -> reshield -> TRUSTPROCESS -> trusted_request
///         round-trip with ETH.
///
/// @dev On a completed request that produced AppEvents, the ProcessorEndpoint
///      claims the app's withdrawal (ETH) into this trigger, then calls execute()
///      and withdraw():
///        - _execute sweeps everything but 1 wei to a fixed sink address (a stand-in
///          for "do something with the funds"), keeping 1 wei behind.
///        - withdraw() (non-overridable) sweeps that leftover 1 wei back to the
///          ProcessorEndpoint, which re-shields it into the app's custody.
///        - _getTrustProcessPayload returns the re-shielded amount so the WASM
///          trusted_request can credit it to an internal fee-collector account.
///
///      Like ActionExecutedTrigger, the payload is produced ONLY when the completed
///      request had AppEvents — a TRUSTPROCESS stateUpdate carries none, so the loop
///      terminates after exactly one TRUSTPROCESS.
contract GuardedTrigger is AbstractTrigger {
  /// @notice Fixed address that receives the swept funds (received amount minus the
  ///         1 wei kept for re-shielding). A stand-in for the trigger's real action.
  address public sink;

  constructor(
    IProcessorEndpoint _processorEndpoint,
    address _sink
  ) AbstractTrigger(_processorEndpoint) {
    if (_sink == address(0)) {
      revert ZeroAddress();
    }
    sink = _sink;
  }

  /// @notice Sends all received ETH except 1 wei to the sink address. The 1 wei is
  ///         left behind so the non-overridable withdraw() sweep re-shields it.
  ///         Skips when there is nothing meaningful to move (e.g. the TRUSTPROCESS
  ///         stateUpdate, where nothing was claimed into this trigger).
  function _execute(Structs.EventData calldata) internal override {
    uint256 bal = address(this).balance;
    if (bal > 1) {
      payable(sink).transfer(bal - 1);
    }
  }

  /// @notice Returns the re-shielded amount (the ETH swept back by withdraw()) so the
  ///         WASM credits it to the fee collector. Returns '' when the completed
  ///         request had no AppEvents, so a TRUSTPROCESS cannot trigger another one.
  ///
  /// @dev withdraw() always appends the ETH entry last in returnedTokens, so its
  ///      amount is the re-shielded wei.
  function _getTrustProcessPayload(
    Structs.EventData calldata appEventData,
    bool,
    bool,
    Structs.TokenAndAmount[] calldata returnedTokens,
    Structs.TokenAndAmount[] calldata
  ) internal pure override returns (bytes memory) {
    if (appEventData.events.length == 0) {
      return '';
    }
    uint256 reshielded = 0;
    if (returnedTokens.length > 0) {
      reshielded = returnedTokens[returnedTokens.length - 1].amount;
    }
    return abi.encode(reshielded);
  }
}
