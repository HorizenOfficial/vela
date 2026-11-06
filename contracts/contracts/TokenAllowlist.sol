// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/access/AccessControl.sol';
import './interfaces/ITokenAllowlist.sol';

/// @title TokenAllowlist
/// @notice Standalone contract managing a global ERC-20 token allowlist.
/// ETH (address(0)) is always implicitly allowed.
contract TokenAllowlist is AccessControl, ITokenAllowlist {
  constructor(address admin) {
    _grantRole(keccak256('ADMIN'), admin);
    // ADMIN administers itself. Left at the DEFAULT_ADMIN_ROLE default it would be unrotatable
    // and unrevokable, since that role is granted to nobody — and this contract is not
    // upgradable, so recovering from a compromised or lost key would mean redeploying the
    // allowlist and re-pointing every ProcessorEndpoint at it.
    _setRoleAdmin(keccak256('ADMIN'), keccak256('ADMIN'));
  }
  mapping(address => bool) public allowedTokens;
  address[] private _allowedTokenList;

  /// @inheritdoc ITokenAllowlist
  function addAllowedToken(address token) external onlyRole(keccak256('ADMIN')) {
    if (token == address(0)) revert TokenAddressCantBeZero();
    if (token.code.length == 0) revert NotAContract();
    if (!allowedTokens[token]) {
      _allowedTokenList.push(token);
    }
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

  /// @inheritdoc ITokenAllowlist
  function getAllowedTokens() public view returns (address[] memory) {
    uint256 i;
    uint256 len = _allowedTokenList.length;
    uint256 count; //count not removed

    while (i != len) {
      if (allowedTokens[_allowedTokenList[i]]) {
        unchecked {
          ++count;
        }
      }
      unchecked {
        ++i;
      }
    }
    address[] memory result = new address[](count);
    uint256 j;
    i = 0;
    while (i != len) {
      if (allowedTokens[_allowedTokenList[i]]) {
        result[j] = _allowedTokenList[i];
        unchecked {
          ++j;
        }
      }
      unchecked {
        ++i;
      }
    }
    return result;
  }
}
