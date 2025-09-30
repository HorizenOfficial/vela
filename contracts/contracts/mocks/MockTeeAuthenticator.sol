// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../interfaces/ITeeAuthenticator.sol";

contract MockTeeAuthenticator is ITeeAuthenticator {

 

    function checkSignature(
        uint256 /*applicationId*/,
        bytes32 /*prevStateRoot*/,
        bytes32 /*newStateRoot*/,
        bytes32 /*processedRequestId*/,
        bytes[] memory /*events*/,
        Structs.WithdrawalRequest[] memory /*withdrawalRequests*/,
        bytes calldata /*signature*/
    ) external pure override returns (bool) {

         return true; // Always return true for mock
    }
}