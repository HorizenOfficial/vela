// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '../trigger/AbstractTrigger.sol';
import '../Structs.sol';

/// @title Test Trigger
/// @notice AbstractTrigger implementation for testing. Combines configurable revert flags
///         with balance capture at the moment _execute runs, so a single contract can verify
///         both the unshield/reshield round-trip and the error-path behaviour.
contract TestTrigger is AbstractTrigger {
  bool public revertOnExecute;
  bool public revertOnPostWithdraw;

  mapping(bytes32 => bool) public executedRequests;
  mapping(bytes32 => bool) public executedPostWithdraws;

  /// @notice Balances captured at the time _execute last ran (only persisted on non-reverting
  ///         paths). Keyed by token address; address(0) for ETH.
  mapping(address => uint256) public capturedBalances;

  error ExecuteReverted();
  error PostWithdrawReverted();

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
    capturedBalances[ETH_TOKEN] = address(this).balance;
    address[] memory tokens = tokenAllowlist.getAllowedTokens();
    uint256 i;
    uint256 len = tokens.length;
    while (i != len) {
      capturedBalances[tokens[i]] = IERC20(tokens[i]).balanceOf(address(this));
      unchecked {
         ++i; 
      }
    }
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
    ReturnedTokens[] memory
  ) public override {
    if (revertOnPostWithdraw) revert PostWithdrawReverted();
    executedPostWithdraws[processedRequestId] = true;
  }
}
