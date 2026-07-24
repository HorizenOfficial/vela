// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol';
import '@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol';
import '@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol';
import './interfaces/IAuthorityChecker.sol';
import './interfaces/IAuthorityRegistry.sol';

/// @title AuthorityRegistry
/// @notice Registry for application-specific authority checkers.
contract AuthorityRegistry is
  Initializable,
  OwnableUpgradeable,
  IAuthorityRegistry,
  UUPSUpgradeable
{
  // mapping appId -> custom authority contract
  mapping(uint256 => IAuthorityChecker) public appAuthorityContracts;

  // default authority contract (fallback)
  IAuthorityChecker public defaultAuthorityContract;

  // --- storage buffer (must always be last, see UPGRADABLE_CONTRACTS_DESIGN.md) ---
  uint256[50] private __gap;

  /// @custom:oz-upgrades-unsafe-allow constructor
  constructor() {
    _disableInitializers();
  }

  /// @param owner Owner address.
  /// @param defaultAuthority Default authority checker.
  function initialize(address owner, address defaultAuthority) external initializer {
    __Ownable_init(owner);
    __UUPSUpgradeable_init();
    if (defaultAuthority == address(0)) revert AddressCantBeZero();
    defaultAuthorityContract = IAuthorityChecker(defaultAuthority);
    emit DefaultAuthorityContractSet(defaultAuthority);
  }

  /// @dev Restricts UUPS upgrades to the owner.
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
