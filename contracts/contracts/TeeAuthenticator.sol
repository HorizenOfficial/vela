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
  uint256 public constant PCR0_LENGTH = 48; //SHA-384 digest length in bytes
  uint256 public constant PCR0_SWAP_APPLY_WINDOW = 7 days;

  INitroProver public immutable nitroProver;
  uint256 public immutable maxVerificationAge;
  uint256 public immutable pcr0UpgradeDelay;

  // Accepted PCR0 values, keyed by keccak256(pcr0) for O(1) membership checks.
  mapping(bytes32 => bool) public acceptedPcr0;
  bytes32[] public acceptedPcr0List; //for enumeration / off-chain audit
  // keccak256 of the PCR0 the platform should currently be running.
  bytes32 public activeImage;

  struct PendingSwap {
    bytes value;
    uint256 eta; //earliest application time
    bool pending;
  }
  PendingSwap public pendingSwap;

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
  /// @param _pcr0 Initial PCR0 value, seeded as the active image.
  /// @param _maxVerificationAge Max age for attestation validity.
  /// @param _pcr0UpgradeDelay Timelock duration for swaps to a new PCR0 (no setter).
  constructor(
    address owner,
    INitroProver _nitroProver,
    bytes memory _pcr0,
    uint256 _maxVerificationAge,
    uint256 _pcr0UpgradeDelay
  ) Ownable(owner) {
    if (_pcr0.length != PCR0_LENGTH) revert InvalidPcr0Length();
    nitroProver = _nitroProver;
    maxVerificationAge = _maxVerificationAge;
    pcr0UpgradeDelay = _pcr0UpgradeDelay;

    bytes32 key = keccak256(_pcr0);
    acceptedPcr0[key] = true;
    acceptedPcr0List.push(key);
    activeImage = key;
    emit Pcr0Swapped(_pcr0);
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function proposePcr0Swap(bytes calldata targetPcr0) external onlyOwner {
    if (targetPcr0.length != PCR0_LENGTH) revert InvalidPcr0Length();
    if (pendingSwap.pending && block.timestamp <= pendingSwap.eta + PCR0_SWAP_APPLY_WINDOW) {
      revert SwapAlreadyPending();
    }

    uint256 eta = block.timestamp + pcr0UpgradeDelay;
    pendingSwap = PendingSwap({value: targetPcr0, eta: eta, pending: true});
    emit Pcr0SwapProposed(targetPcr0, eta);
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function applyPcr0Swap() external onlyOwner {
    PendingSwap memory p = pendingSwap;
    if (!p.pending) revert NoPendingSwap();
    if (block.timestamp > p.eta + PCR0_SWAP_APPLY_WINDOW) revert SwapProposalExpired();

    bytes32 key = keccak256(p.value);
    // New images must clear the audit window; already-accepted images may swap now (rollback).
    if (!acceptedPcr0[key]) {
      if (block.timestamp < p.eta) revert TimelockNotElapsed();
      acceptedPcr0[key] = true;
      acceptedPcr0List.push(key);
    }

    activeImage = key;
    delete pendingSwap;
    emit Pcr0Swapped(p.value);
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function cancelPcr0Swap() external onlyOwner {
    PendingSwap memory p = pendingSwap;
    if (!p.pending) revert NoPendingSwap();

    delete pendingSwap;
    emit Pcr0SwapCancelled(p.value);
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function removePcr0(bytes calldata oldPcr0) external onlyOwner {
    bytes32 key = keccak256(oldPcr0);
    if (!acceptedPcr0[key]) revert UnknownPcr0();
    //activeImage is always in the accepted set, so this also keeps the set non-empty
    if (key == activeImage) revert CannotRemoveActiveImage();

    acceptedPcr0[key] = false;
    uint256 length = acceptedPcr0List.length;
    for (uint256 i; i != length; ) {
      if (acceptedPcr0List[i] == key) {
        acceptedPcr0List[i] = acceptedPcr0List[length - 1];
        acceptedPcr0List.pop();
        break;
      }
      unchecked {
        ++i;
      }
    }
    emit Pcr0Removed(oldPcr0);
  }

  /// @inheritdoc ITeeAuthenticatorAdmin
  function getAcceptedPcr0Count() external view returns (uint256) {
    return acceptedPcr0List.length;
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
    //activeImage may have changed since step 1: re-validate before finalizing
    _checkAttestationContent(_rawPcrs, _enclaveKey, _userData);
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
    if (pcrs.length < 4 + PCR0_LENGTH) {
      revert InvalidPCR();
    }

    //PCR0 is a fixed-length SHA-384 digest at offset 4
    bytes memory candidate = new bytes(PCR0_LENGTH);
    for (uint256 i; i != PCR0_LENGTH; ) {
      candidate[i] = pcrs[i + 4];
      unchecked {
        ++i;
      }
    }
    //only the image the platform should currently be running may register the tee signer
    if (keccak256(candidate) != activeImage) {
      revert InvalidPCR();
    }
  }
}
