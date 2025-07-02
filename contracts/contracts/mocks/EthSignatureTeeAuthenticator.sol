// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../interfaces/ITeeAuthenticator.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract EthSignatureTeeAuthenticator is ITeeAuthenticator {

    address validSigner;

    constructor(address _validSigner) {
        validSigner = _validSigner;
    }

    function checkSignature(
        uint8 applicationId,
        bytes calldata prevStateRoot,
        bytes calldata newStateRoot,
        bytes[] calldata events,
        Structs.WithdrawalRequest[] calldata withdrawalRequests,
        bytes calldata signature
    ) external view override returns (bool) {
        bytes32 messageHash = keccak256(abi.encode(
            applicationId,
            prevStateRoot,
            newStateRoot,
            events,
            withdrawalRequests
        ));

        address recovered = ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(messageHash), signature);
        return recovered == validSigner;
    }
}