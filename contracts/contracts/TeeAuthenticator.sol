// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './AbstractTeeAuthenticator.sol';
import './interfaces/INitroProver.sol';
import './interfaces/ITeeAuthenticatorAdmin.sol';
import '@openzeppelin/contracts/access/Ownable.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';
import '@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol';

/// @title TeeAuthenticator
/// @notice Attestation-based tee signer and key manager.
contract TeeAuthenticator is AbstractTeeAuthenticator, ITeeAuthenticatorAdmin, Ownable {
  INitroProver public immutable nitroProver;
  bytes public pcr0;
  uint256 public immutable maxVerificationAge;

  address public teeSigner;
  bytes public pubSecp521r1;

  mapping(bytes32 => bool) private _usedAttestations;

  //step update
  uint256 public currentUpdateStep;
  uint256 public step2CurrentIndex;
  bytes private _pubKey;

  bytes32 private _attestationHash;
  bytes private _enclaveKey;
  bytes private _userData;
  bytes private _rawPcrs;
  bytes private _attestationSig;
  bytes private _certificate;
  bytes[] private _cabundle;
  bytes[] private _attestation_decoded;
  bytes private _buf;

  /// @param owner Owner address.
  /// @param _nitroProver Nitro prover contract.
  /// @param _pcr0 Initial PCR0 value.
  /// @param _maxVerificationAge Max age for attestation validity.
  constructor(
    address owner,
    INitroProver _nitroProver,
    bytes memory _pcr0,
    uint256 _maxVerificationAge
  ) Ownable(owner) {
    pcr0 = _pcr0;
    nitroProver = _nitroProver;
    maxVerificationAge = _maxVerificationAge;
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function updateTee(bytes calldata attestation) external onlyOwner {
    bytes32 attestationHash = keccak256(attestation);
    if (_usedAttestations[attestationHash]) revert AttestationAlreadyUsed();

    (bytes memory enclaveKey, bytes memory userData, bytes memory rawPcrs) = nitroProver
      .verifyAttestation(attestation, maxVerificationAge);
    _checkAttestationContent(rawPcrs, enclaveKey, userData);
    _updateTee(address(bytes20(userData)), enclaveKey, attestationHash);
  }

  // -- STEPS UPDATE
  // if you want to reset and begin a new step update, invoke step 1
  /// @inheritdoc ITeeAuthenticatorAdmin
  function updateTeeStep1(bytes calldata attestation) external onlyOwner {
    _resetStepUpdate();

    _attestationHash = keccak256(attestation);
    if (_usedAttestations[_attestationHash]) revert AttestationAlreadyUsed();

    (_attestation_decoded, _certificate, _cabundle, _enclaveKey, _userData, _rawPcrs) = nitroProver
      .verifyAttestationStep1(attestation, maxVerificationAge);

    _checkAttestationContent(_rawPcrs, _enclaveKey, _userData);

    currentUpdateStep = 1;
  }

  //this should be invoked "getStep2TotalLength()" times
  //check currentUpdateStep to see if you can go to the next step
  /// @inheritdoc ITeeAuthenticatorAdmin
  function updateTeeStep2() external onlyOwner {
    if (currentUpdateStep != 1) revert WrongStep();

    _pubKey = nitroProver.verifyAttestationStep2(_cabundle, step2CurrentIndex, _pubKey);
    unchecked {
      ++step2CurrentIndex;
    }

    if (step2CurrentIndex == _cabundle.length) {
      currentUpdateStep = 2; //you can go to step 3
    }
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function updateTeeStep3() external onlyOwner {
    if (currentUpdateStep != 2) revert WrongStep();
    (_attestationSig, _pubKey, _buf) = nitroProver.verifyAttestationStep3(
      _attestation_decoded,
      _certificate,
      _pubKey
    );
    currentUpdateStep = 3;
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function updateTeeStep4() external onlyOwner {
    if (currentUpdateStep != 3) revert WrongStep();
    nitroProver.verifyAttestationStep4(_attestationSig, _pubKey, _buf);
    _updateTee(address(bytes20(_userData)), _enclaveKey, _attestationHash);
  }

  function _updateTee(
    address newTeeSigner,
    bytes memory newPubSecp521r1,
    bytes32 attestationHash
  ) internal {
    emit TeeUpdate(teeSigner, newTeeSigner, pubSecp521r1, newPubSecp521r1);
    _usedAttestations[attestationHash] = true;
    teeSigner = newTeeSigner;
    pubSecp521r1 = newPubSecp521r1;
    _resetStepUpdate();
  }

  function _resetStepUpdate() internal {
    currentUpdateStep = 0;
    step2CurrentIndex = 0;
    _pubKey = bytes('');
  }

  //check number of transactions needed to complete step2
  /// @inheritdoc ITeeAuthenticatorAdmin
  function getStep2TotalLength() external view returns (uint256) {
    return _cabundle.length;
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function updatePcr0(bytes calldata newPcr0) external onlyOwner {
    emit PcrZeroUpdate(pcr0, newPcr0);
    pcr0 = newPcr0;
  }

  /// @inheritdoc ITeeAuthenticator
  function getTeeSigner() public view override returns (address) {
    return teeSigner;
  }
  /// @inheritdoc ITeeAuthenticator
  function getPubSecp521r1() public view override returns (bytes memory) {
    return pubSecp521r1;
  }

  function _checkAttestationContent(
    bytes memory pcrs,
    bytes memory enclaveKey,
    bytes memory userData
  ) internal view {
    if (userData.length != 20) {
      //it contains an address
      revert InvalidUserDataLength();
    }

    if (enclaveKey.length != PK_LENGTH) {
      revert InvalidPKLength();
    }
    if (pcrs.length < 4 + pcr0.length) {
      revert InvalidPCR();
    }

    uint256 i;
    uint256 length = pcr0.length;
    while (i != length) {
      if (pcrs[i + 4] != pcr0[i]) {
        revert InvalidPCR();
      }
      unchecked {
        ++i;
      }
    }
  }
}
