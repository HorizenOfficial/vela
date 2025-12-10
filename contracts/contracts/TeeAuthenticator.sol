// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "./interfaces/ITeeAuthenticator.sol";
import "./interfaces/INitroValidator.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract TeeAuthenticator is ITeeAuthenticator, Ownable {
    uint256 public constant PK_LENGTH = 133; //secp521r1 uncompressed public key length in bytes
    INitroValidator public immutable nitroValidator;
    bytes32 public pcr0;
    uint256 public immutable maxVerificationAge;

    address public teeSigner;
    bytes public pubSecp521r1;

    //events
    event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1);
    event PcrZeroUpdate(bytes32 oldPcr0, bytes32 newPcr0);

    //error
    error InvalidPCR();
    error TeeIsNotSet();
    error InvalidPKLength();

    constructor(address owner, INitroValidator _nitroValidator, bytes32  _pcr0, uint256 _maxVerificationAge) Ownable(owner) {
        pcr0 = _pcr0;
        nitroValidator = _nitroValidator;
        maxVerificationAge = _maxVerificationAge;
    }

    function updateTee(address newTeeSigner, bytes memory attestation) public onlyOwner {
        (bytes32 _pcr0, bytes memory newPubSecp521r1) = nitroValidator.verifyAttestation(attestation, maxVerificationAge);

        if(_pcr0 != pcr0) {
            revert InvalidPCR();
        }
        if(newPubSecp521r1.length != PK_LENGTH) {
            revert InvalidPKLength();
        }

        emit TeeUpdate(teeSigner, newTeeSigner, pubSecp521r1, newPubSecp521r1);
        teeSigner = newTeeSigner;
        pubSecp521r1 = newPubSecp521r1;
    }

    function updatePcr0(bytes32 newPcr0) public onlyOwner {
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
}