// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './IAuthorityChecker.sol';

/// @title DefaultAuthorityChecker interface
/// @notice External API, events, and errors for the default authority checker.
interface IDefaultAuthorityChecker is IAuthorityChecker {
  /// @notice Emitted when an authority is added.
  /// @param applicationId Application identifier.
  /// @param authority Authority address.
  event AddedAuthority(uint256 indexed applicationId, address indexed authority);

  /// @notice Emitted when an authority is removed.
  /// @param applicationId Application identifier.
  /// @param authority Authority address.
  event RemovedAuthority(uint256 indexed applicationId, address indexed authority);

  /// @notice Authority is not present.
  error AuthorityNotPresent();

  /// @notice Authority is already present.
  error AuthorityAlreadyPresent();

  /// @notice Adds an allowed authority for an application.
  /// @param applicationId Application identifier.
  /// @param authority Authority address.
  function addAllowedAuthority(uint256 applicationId, address authority) external;

  /// @notice Removes an allowed authority for an application.
  /// @param applicationId Application identifier.
  /// @param authority Authority address.
  function removeAllowedAuthority(uint256 applicationId, address authority) external;
}
