// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/access/AccessControl.sol';
import './interfaces/ITokenAllowlist.sol';

/// @title TokenAllowlist
/// @notice Base contract managing a global ERC-20 token allowlist.
/// ETH (address(0)) is always implicitly allowed.
abstract contract TokenAllowlist is AccessControl, ITokenAllowlist {
  mapping(address => bool) public allowedTokens;

  /// @inheritdoc ITokenAllowlist
  function addAllowedToken(address token) external onlyRole(keccak256('ADMIN')) {
    if (token == address(0)) revert TokenAddressCantBeZero();
    if (token.code.length == 0) revert NotAContract();
    allowedTokens[token] = true;
    emit TokenAllowed(token);
  }

  /// @inheritdoc ITokenAllowlist
  function removeAllowedToken(address token) external onlyRole(keccak256('ADMIN')) {
    if (token == address(0)) revert TokenAddressCantBeZero();
    allowedTokens[token] = false;
    emit TokenRemoved(token);
  }

  /// @inheritdoc ITokenAllowlist
  function isAllowedToken(address token) external view returns (bool) {
    if (token == address(0)) return true;
    return allowedTokens[token];
  }
}
