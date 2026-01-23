// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

/// @title AuthorityRegistry interface
/// @notice External API, events, and errors for the authority registry.
interface IAuthorityRegistry {
  /// @notice Emitted when an application-specific authority contract is set.
  /// @param applicationId Application identifier.
  /// @param authorityContract Authority checker contract.
  event AppAuthorityContractSet(uint256 indexed applicationId, address indexed authorityContract);

  /// @notice Emitted when the default authority contract is set.
  /// @param authorityContract Default authority checker contract.
  event DefaultAuthorityContractSet(address indexed authorityContract);

  /// @notice A zero address was supplied where not allowed.
  error AddressCantBeZero();

  /// @notice Sets a custom authority checker for an application id.
  /// @param applicationId Application identifier.
  /// @param authorityContract Authority checker contract.
  function setAppAuthorityContract(uint256 applicationId, address authorityContract) external;

  /// @notice Sets the default authority checker contract.
  /// @param authorityContract Default authority checker contract.
  function setDefaultAuthorityContract(address authorityContract) external;

  /// @notice Checks whether an authority is allowed for an application.
  /// @param applicationId Application identifier.
  /// @param authority Authority address.
  /// @return allowed True if the authority is allowed.
  function checkAuthorityIsAllowed(
    uint256 applicationId,
    address authority
  ) external view returns (bool);
}
