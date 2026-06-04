// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../interfaces/IProcessorEndpoint.sol';
import '../interfaces/ITokenAllowlist.sol';
import '../Structs.sol';
import './AbstractTrigger.sol';

/// @title Default Trigger
/// @notice Concrete implementation of AbstractTrigger with no-op execute and post-withdraw hooks.
///         Suitable as a drop-in trigger when no custom execution logic is required.
contract DefaultTrigger is AbstractTrigger {
  /// @param _processorEndpoint ProcessorEndpoint that will call execute and withdraw.
  /// @param _tokenAllowlist Allowlist used to determine which ERC-20 tokens are swept on withdraw.
  constructor(
    IProcessorEndpoint _processorEndpoint,
    ITokenAllowlist _tokenAllowlist
  ) AbstractTrigger(_processorEndpoint, _tokenAllowlist) {}

  /// @notice No-op execute hook; DefaultTrigger performs no action on request completion.
  function _execute(
    Structs.EventData calldata
  ) internal override {}

  /// @notice No-op post-withdraw hook; DefaultTrigger performs no action after the sweep.
  function _postWithdraw(
    Structs.EventData calldata,
    bool,
    Structs.TokenAndAmount[] memory,
    Structs.TokenAndAmount[] memory
  ) public override {}
}
