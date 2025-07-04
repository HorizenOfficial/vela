// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../Structs.sol";

interface ITeeAuthenticator {
    function checkSignature(
        uint256 applicationId, 
        bytes calldata prevStateRoot, 
        bytes calldata newStateRoot, 
        uint256[] calldata processedRequestIds,
        bytes[] memory events, 
        Structs.WithdrawalRequest[] memory withdrawalRequests, 
        bytes memory signature
    ) external view returns(bool);
}