// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IAuthorityChecker.sol";
import "./interfaces/IAuthorityRegistry.sol";

/// @title AuthorityRegistry
/// @notice Registry for application-specific authority checkers.
contract AuthorityRegistry is Ownable, IAuthorityRegistry {

    // mapping appId -> custom authority contract
    mapping(uint256 => IAuthorityChecker) public appAuthorityContracts;

    // default authority contract (fallback)
    IAuthorityChecker public defaultAuthorityContract;

    /// @param owner Owner address.
    /// @param defaultAuthority Default authority checker.
    constructor(address owner, address defaultAuthority) Ownable(owner) {
        if(defaultAuthority == address(0)) revert AddressCantBeZero();
        defaultAuthorityContract = IAuthorityChecker(defaultAuthority);
        emit DefaultAuthorityContractSet(defaultAuthority);
    }

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
        if(authorityContract == address(0)) revert AddressCantBeZero();
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
