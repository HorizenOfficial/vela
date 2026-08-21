// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol';
import '@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol';
import './interfaces/IAuthorityChecker.sol';
import './interfaces/IAuthorityRegistry.sol';

/// @title AuthorityRegistry
/// @notice Registry for application-specific authority checkers.
/// @dev Deployed behind a UUPS proxy (`docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`).
contract AuthorityRegistry is OwnableUpgradeable, UUPSUpgradeable, IAuthorityRegistry {
  // mapping appId -> custom authority contract
  mapping(uint256 => IAuthorityChecker) public appAuthorityContracts;

  // default authority contract (fallback)
  IAuthorityChecker public defaultAuthorityContract;

  /// @dev Reserved storage for future versions (see `docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`).
  ///      Must be reduced by the number of slots any new variable added above it consumes, and
  ///      must always remain the last declaration in this contract.
  uint256[50] private __gap;

  /// @custom:oz-upgrades-unsafe-allow constructor
  constructor() {
    _disableInitializers();
  }

  /// @param owner Owner address.
  /// @param defaultAuthority Default authority checker.
  function initialize(address owner, address defaultAuthority) external initializer {
    __Ownable_init(owner);
    if (defaultAuthority == address(0)) revert AddressCantBeZero();
    defaultAuthorityContract = IAuthorityChecker(defaultAuthority);
    emit DefaultAuthorityContractSet(defaultAuthority);
  }

  /// @inheritdoc UUPSUpgradeable
  function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

  /// @inheritdoc IAuthorityRegistry
  function setAppAuthorityContract(
    uint256 applicationId,
    address authorityContract
  ) external onlyOwner {
    appAuthorityContracts[applicationId] = IAuthorityChecker(authorityContract);
    emit AppAuthorityContractSet(applicationId, authorityContract);
  }

  /// @inheritdoc IAuthorityRegistry
  function setDefaultAuthorityContract(address authorityContract) external onlyOwner {
    if (authorityContract == address(0)) revert AddressCantBeZero();
    defaultAuthorityContract = IAuthorityChecker(authorityContract);
    emit DefaultAuthorityContractSet(authorityContract);
  }

  /// @inheritdoc IAuthorityRegistry
  function checkAuthorityIsAllowed(
    uint256 applicationId,
    address authority
  ) external view returns (bool) {
    IAuthorityChecker impl = appAuthorityContracts[applicationId];

    if (address(impl) == address(0)) {
      impl = defaultAuthorityContract;
    }

    return impl.checkAuthorityIsAllowed(applicationId, authority);
  }
}
