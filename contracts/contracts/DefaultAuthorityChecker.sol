// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IAuthorityChecker.sol";
import "./interfaces/IDefaultAuthorityChecker.sol";

/// @title DefaultAuthority
/// @notice Ownable authority registry for per-application allowlists.
contract DefaultAuthority is Ownable, IDefaultAuthorityChecker {

    mapping(uint256 => mapping(address => bool)) allowedAuthorities;

    /// @param owner Owner address.
    constructor(address owner) Ownable(owner) {}

    /// @inheritdoc IDefaultAuthorityChecker
    function addAllowedAuthority(uint256 applicationId, address authority) external onlyOwner {
        if (allowedAuthorities[applicationId][authority]) revert AuthorityAlreadyPresent();
        allowedAuthorities[applicationId][authority] = true;
        emit AddedAuthority(applicationId, authority);
    }

    /// @inheritdoc IDefaultAuthorityChecker
    function removeAllowedAuthority(uint256 applicationId, address authority) external onlyOwner {
        if (!allowedAuthorities[applicationId][authority]) revert AuthorityNotPresent();
        allowedAuthorities[applicationId][authority] = false;
        emit RemovedAuthority(applicationId, authority);
    }

    /// @inheritdoc IAuthorityChecker
    function checkAuthorityIsAllowed(
        uint256 applicationId,
        address authority
    ) external view override returns (bool) {
        return allowedAuthorities[applicationId][authority];
    }
}
