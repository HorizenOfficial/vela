// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

/// @title ITokenAllowlist
/// @notice Interface for managing a global ERC-20 token allowlist.
interface ITokenAllowlist {
  /// @notice Emitted when a token is added to the allowlist.
  /// @param token The token address.
  event TokenAllowed(address indexed token);

  /// @notice Emitted when a token is removed from the allowlist.
  /// @param token The token address.
  event TokenRemoved(address indexed token);

  /// @notice A zero token address was supplied where not allowed.
  error TokenAddressCantBeZero();

  /// @notice The token address is not a contract.
  error NotAContract();

  /// @notice The token is not in the global allowlist.
  error TokenNotAllowed();

  /// @notice Adds an ERC-20 token to the global allowlist.
  /// @param token The token contract address. Must not be address(0) and must contain code.
  function addAllowedToken(address token) external;

  /// @notice Removes an ERC-20 token from the global allowlist.
  /// @param token The token contract address. Must not be address(0).
  function removeAllowedToken(address token) external;

  /// @notice Returns whether a token is allowed. ETH (address(0)) is always allowed.
  /// @param token The token address to check.
  /// @return True if the token is allowed.
  function isAllowedToken(address token) external view returns (bool);
}
