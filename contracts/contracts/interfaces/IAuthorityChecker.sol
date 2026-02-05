// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

/// @title AuthorityChecker interface
/// @notice External API for authority checks.
interface IAuthorityChecker {
  /// @notice Checks whether an authority is allowed for an application.
  /// @param applicationId Application identifier.
  /// @param authority Authority address.
  /// @return allowed True if the authority is allowed.
  function checkAuthorityIsAllowed(
    uint256 applicationId,
    address authority
  ) external view returns (bool);
}
