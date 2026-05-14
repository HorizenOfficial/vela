// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../trigger/AbstractTrigger.sol';
import '../interfaces/IProcessorEndpoint.sol';
import '../interfaces/ITokenAllowlist.sol';
import '../Structs.sol';

/// @title Mock Trigger
/// @notice AbstractTrigger implementation for testing. Reverts in _execute and/or _postWithdraw
///         based on constructor flags.
contract MockTrigger is AbstractTrigger {
  /// @notice If true, _execute will revert.
  bool public revertOnExecute;
  /// @notice If true, _postWithdraw will revert.
  bool public revertOnPostWithdraw;

  mapping(bytes32 => bool) public executedRequests;
  mapping(bytes32 => bool) public executedPostWithdraws;

  error ExecuteReverted();
  error PostWithdrawReverted();

  /// @param _processorEndpoint ProcessorEndpoint that will call execute and withdraw.
  /// @param _tokenAllowlist Allowlist used to determine which ERC-20 tokens are swept on withdraw.
  /// @param _revertOnExecute If true, _execute reverts.
  /// @param _revertOnPostWithdraw If true, _postWithdraw reverts.
  constructor(
    IProcessorEndpoint _processorEndpoint,
    ITokenAllowlist _tokenAllowlist,
    bool _revertOnExecute,
    bool _revertOnPostWithdraw
  ) AbstractTrigger(_processorEndpoint, _tokenAllowlist) {
    revertOnExecute = _revertOnExecute;
    revertOnPostWithdraw = _revertOnPostWithdraw;
  }

  function _execute(
    uint64,
    bytes32,
    bytes32,
    bytes32 processedRequestId,
    Structs.EventData calldata,
    Structs.EventData calldata,
    Structs.WithdrawalRequest[] calldata,
    uint256,
    uint256,
    Structs.ErrorCode,
    string calldata
  ) internal override {
    if (revertOnExecute) revert ExecuteReverted();
    executedRequests[processedRequestId] = true;
  }

  function _postWithdraw(
    uint64,
    bytes32,
    bytes32,
    bytes32 processedRequestId,
    Structs.EventData calldata,
    Structs.EventData calldata,
    Structs.WithdrawalRequest[] calldata,
    uint256,
    uint256,
    Structs.ErrorCode,
    string calldata,
    MovedTokens[] memory
  ) public override {
    if (revertOnPostWithdraw) revert PostWithdrawReverted();
    executedPostWithdraws[processedRequestId] = true;
  }
}
