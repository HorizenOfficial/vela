// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../interfaces/ITeeAuthenticator.sol";

contract MockTeeAuthenticator is ITeeAuthenticator {    
    address public teeSigner;
    bytes public pubSecp521r1;

    constructor(address _teeSigner, bytes memory _pubSecp521r1) {
        teeSigner = _teeSigner;
        pubSecp521r1 = _pubSecp521r1;
    }

    function checkSignature(
        uint256 applicationId,
        bytes32 prevStateRoot,
        bytes32 newStateRoot,
        bytes32 processedRequestId,
        bytes[] memory events,
        Structs.WithdrawalRequest[] memory /*withdrawalRequests*/,
        bytes calldata /*signature*/
    ) external pure override returns (bool) {
        return true; // Always return true for mock
    }

    function getTeeSigner() external view override returns(address) {
        return teeSigner;
    }
    function getPubSecp521r1() external view returns(bytes memory) {
        return pubSecp521r1;
    }
}