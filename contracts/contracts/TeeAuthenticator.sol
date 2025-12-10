// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "./interfaces/ITeeAuthenticator.sol";
import "./interfaces/INitroProver.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

import "hardhat/console.sol";

contract TeeAuthenticator is ITeeAuthenticator, Ownable {
    uint256 public constant PK_LENGTH = 133; //secp521r1 uncompressed public key length in bytes
    INitroProver public immutable nitroProver;
    bytes public pcr0;
    uint256 public immutable maxVerificationAge;

    address public teeSigner;
    bytes public pubSecp521r1;

    //events
    event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1);
    event PcrZeroUpdate(bytes oldPcr0, bytes newPcr0);

    //error
    error InvalidPCR();
    error TeeIsNotSet();
    error InvalidPKLength();

    constructor(address owner, INitroProver _nitroProver, bytes memory _pcr0, uint256 _maxVerificationAge) Ownable(owner) {
        pcr0 = _pcr0;
        nitroProver = _nitroProver;
        maxVerificationAge = _maxVerificationAge;
    }

    function updateTee(bytes memory attestation) public onlyOwner {
        (bytes memory newPubSecp521r1, bytes memory userData ,bytes memory rawPcrs) = nitroProver.verifyAttestation(attestation, maxVerificationAge);
        console.log("enclaveKey");
        logBytesHex(newPubSecp521r1);
        console.log("userData");
        logBytesHex(userData);
        console.log("rawPcrs");
        logBytesHex(rawPcrs);

        if(keccak256(abi.encodePacked(rawPcrs)) != keccak256(abi.encodePacked(pcr0))) {
            revert InvalidPCR();
        }
        if(newPubSecp521r1.length != PK_LENGTH) {
            revert InvalidPKLength();
        }
        address newTeeSigner = address(0);

        emit TeeUpdate(teeSigner, newTeeSigner, pubSecp521r1, newPubSecp521r1);
        teeSigner = newTeeSigner;
        pubSecp521r1 = newPubSecp521r1;
    }

    function updatePcr0(bytes memory newPcr0) public onlyOwner {
        emit PcrZeroUpdate(pcr0, newPcr0);
        pcr0 = newPcr0;
    }

    function checkSignature(
        uint64 applicationId,
        bytes32 prevStateRoot,
        bytes32 newStateRoot,
        bytes32 processedRequestId,
        bytes[] memory events,
        Structs.WithdrawalRequest[] memory withdrawalRequests,
        uint256 refundAmount, 
        uint256 applicationFee,
        bytes calldata signature
    ) external view override returns (bool) {
        if(teeSigner == address(0) || pubSecp521r1.length != PK_LENGTH) revert TeeIsNotSet();

        bytes32 eventsHash = keccak256(abi.encode(events));
        bytes32 withdrawalRequestsHash = keccak256(abi.encode(withdrawalRequests));

        bytes32 messageHash = keccak256(abi.encode(
            applicationId,
            prevStateRoot,
            newStateRoot,
            processedRequestId,
            eventsHash,
            withdrawalRequestsHash,
            refundAmount,
            applicationFee
        ));

        address recovered = ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(messageHash), signature);
        return recovered == teeSigner;
    }

    function getTeeSigner() external view override returns(address) {
        return teeSigner;
    }
    function getPubSecp521r1() external view returns(bytes memory) {
        return pubSecp521r1;
    }

    //TODO remove
    function logBytesHex(bytes memory data) public pure {
        bytes memory hexChars = "0123456789abcdef";
        bytes memory str = new bytes(2 + data.length * 2);
        str[0] = "0";
        str[1] = "x";
        for (uint i = 0; i < data.length; i++) {
            str[2 + i*2] = hexChars[uint(uint8(data[i] >> 4))];
            str[3 + i*2] = hexChars[uint(uint8(data[i] & 0x0f))];
        }
        console.log(string(str));
    }
}