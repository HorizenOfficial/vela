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

  mapping(uint256 => bool) public executedInBlock;
  mapping(uint256 => bool) public executedPostWithdrawsInBlock;

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
    Structs.EventData calldata
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
    executedInBlock[block.number] = true;
  }

  function _postWithdraw(
    Structs.EventData calldata,
    bool,
    Structs.TokenAndAmount[] memory,
    Structs.TokenAndAmount[] memory
  ) public override {
    if (revertOnPostWithdraw) revert PostWithdrawReverted();
    executedPostWithdrawsInBlock[block.number] = true;
  }
}
