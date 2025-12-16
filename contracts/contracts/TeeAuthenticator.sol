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

    //step update
    uint256 public currentUpdateStep;
    bytes private _enclaveKey;
    bytes private _userData;
    bytes private _rawPcrs;
    bytes private _attestationSig;
    bytes private _pubKey;
    bytes private _certificate;
    bytes[] private _cabundle;
    bytes[] private _attestation_decoded;
    bytes private _buf;
    
    uint256 public step2_1CurrentIndex;

    //events
    event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1);
    event PcrZeroUpdate(bytes indexed oldPcr0, bytes indexed newPcr0);

    //error
    error InvalidPCR();
    error TeeIsNotSet();
    error InvalidPKLength();
    error WrongStep();

    constructor(address owner, INitroProver _nitroProver, bytes memory _pcr0, uint256 _maxVerificationAge) Ownable(owner) {
        pcr0 = _pcr0;
        nitroProver = _nitroProver;
        maxVerificationAge = _maxVerificationAge;
    }

    function updateTee(bytes calldata attestation) public onlyOwner {
        (bytes memory enclaveKey, bytes memory userData, bytes memory rawPcrs) = nitroProver.verifyAttestation(attestation, maxVerificationAge);
        _checkPcr(rawPcrs);
        _updateTee(address(bytes20(userData)), enclaveKey);
    }

    // -- STEPS UPDATE
    function updateTeeStep1(bytes calldata attestation) public onlyOwner {
        (_attestation_decoded, _certificate, _cabundle, _enclaveKey, _userData, _rawPcrs) = nitroProver.verifyAttestationStep1(attestation, maxVerificationAge);
        currentUpdateStep = 1;
    }

    function updateTeeStep2() public onlyOwner returns(bool canGoToNextStep){
        if(currentUpdateStep != 1) revert WrongStep();

        _pubKey = nitroProver.verifyAttestationStep2(_cabundle, step2_1CurrentIndex, _pubKey);
        step2_1CurrentIndex++;

        if(step2_1CurrentIndex == _cabundle.length) {
            currentUpdateStep = 2;
            return true; //step 2 finished
        }
        return false;
    }

    function updateTeeStep3() public onlyOwner{
        if(currentUpdateStep != 2) revert WrongStep();
        (_attestationSig, _pubKey, _buf) = nitroProver.verifyAttestationStep3(_attestation_decoded, _certificate, _pubKey);
        currentUpdateStep = 3;
    }

    function updateTeeStep4() public onlyOwner {
        if(currentUpdateStep != 3) revert WrongStep();
        nitroProver.verifyAttestationStep4(_attestationSig, _pubKey, _buf);
        _checkPcr(_rawPcrs);
        _updateTee(address(bytes20(_userData)), _enclaveKey);
    }

    function getStep2_1TotalLength() public view returns(uint256) {
        return _cabundle.length;
    }

    function _updateTee(address newTeeSigner, bytes memory newPubSecp521r1) internal {
        if(newPubSecp521r1.length != PK_LENGTH) {
            revert InvalidPKLength();
        }

        emit TeeUpdate(teeSigner, newTeeSigner, pubSecp521r1, newPubSecp521r1);
        teeSigner = newTeeSigner;
        pubSecp521r1 = newPubSecp521r1;
        currentUpdateStep = 0;
        step2_1CurrentIndex = 0;
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
            unchecked { ++i; }
        }
    }
}