// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../Structs.sol";

interface ITeeAuthenticator {
    function checkSignature(
        uint256 applicationId, 
        bytes32 prevStateRoot, 
        bytes32 newStateRoot, 
        bytes32 processedRequestId,
        bytes[] memory events, 
        Structs.WithdrawalRequest[] memory withdrawalRequests, 
        bytes memory signature
    ) external view returns(bool);

    function getTeeSigner() external view returns(address);
    function getPubSecp521r1() external view returns(bytes memory);
}