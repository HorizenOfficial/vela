// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface IAuthorityChecker {
    function checkAuthorityIsAllowed(
        uint256 applicationId,
        address authority
    ) external view returns (bool);
}