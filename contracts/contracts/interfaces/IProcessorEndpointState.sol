// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './ITokenAllowlist.sol';

/// @title IProcessorEndpointState
/// @notice The `ProcessorEndpoint` public state readers that are auto-generated getters of state
///         variables rather than hand-written functions.
/// @dev These are declared apart from the rest of `IProcessorEndpoint` (which inherits this
///      interface, so external consumers see one combined interface as before) because of where
///      the variables live. They are declared in `ProcessorEndpointStorage`, the shared base of
///      `ProcessorEndpoint` and `ProcessorEndpointExtension` — and a public state variable can
///      only implement an interface function when the interface is inherited by the *same*
///      contract that declares the variable. `ProcessorEndpointStorage` therefore inherits this
///      interface, while it cannot inherit the full `IProcessorEndpoint`: that would force
///      `ProcessorEndpointExtension` to implement every endpoint function.
interface IProcessorEndpointState {
  /// @notice External token allowlist contract consulted for every ERC-20 deposit.
  /// @return allowlist The allowlist contract.
  function tokenAllowlist() external view returns (ITokenAllowlist allowlist);

  /// @notice Returns the custody balance for a given application and token.
  /// @param applicationId Application identifier.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @return amount Current custody balance.
  function appCustody(
    uint64 applicationId,
    address tokenAddress
  ) external view returns (uint256 amount);

  /// @notice Returns the total custody balance across all applications for a given token.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @return amount Total custody balance.
  function totalAppCustody(address tokenAddress) external view returns (uint256 amount);

  /// @notice Returns the pending claim balance for a given token and payee.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @param payee Payee address.
  /// @return amount Pending claim balance.
  function pendingClaims(
    address tokenAddress,
    address payee
  ) external view returns (uint256 amount);

  /// @notice Returns the total pending claims for a given token.
  /// @param tokenAddress Token address (0x0 = ETH).
  /// @return amount Total pending claims.
  function totalPendingClaims(address tokenAddress) external view returns (uint256 amount);

  /// @notice Current facilitator nonce for a user, needed to build the next `submitRequestFor`
  ///         authorization signature.
  /// @param user User address.
  /// @return nonce Current nonce value.
  function facilitatorNonces(address user) external view returns (uint256 nonce);
}
