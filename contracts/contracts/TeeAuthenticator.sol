// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "./interfaces/ITeeAuthenticator.sol";
import "./interfaces/INitroProver.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import "hardhat/console.sol"; //TODO REMOVE

contract TeeAuthenticator is ITeeAuthenticator, Ownable {
    uint256 public constant PK_LENGTH = 133; //secp521r1 uncompressed public key length in bytes
    INitroProver public immutable nitroProver;
    bytes public pcr0;
    uint256 public immutable maxVerificationAge;

    address public teeSigner;
    bytes public pubSecp521r1;

    //events
    event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1);
    event PcrZeroUpdate(bytes indexed oldPcr0, bytes indexed newPcr0);

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
        (bytes memory enclaveKey, bytes memory userData, bytes memory rawPcrs) = nitroProver.verifyAttestation(attestation, maxVerificationAge);

        _checkPcr(rawPcrs);
        if(enclaveKey.length != PK_LENGTH) {
            revert InvalidPKLength();
        }

        address newTeeSigner = address(bytes20(userData));

        emit TeeUpdate(teeSigner, newTeeSigner, pubSecp521r1, enclaveKey);
        teeSigner = newTeeSigner;
        pubSecp521r1 = enclaveKey;
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

    function _checkPcr(bytes memory pcrs) internal view {
        if (pcrs.length < 4 + pcr0.length) {
            revert InvalidPCR();
        }

        uint256 i;
        uint256 length = pcr0.length;
        while (i != length) {
            if (pcrs[i + 4] != pcr0[i]) {
                revert InvalidPCR();
            }
            i++;
        }
    }
}