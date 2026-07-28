// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import '../interfaces/INitroProver.sol';

/// @title MockNitroProver
/// @notice Test double returning configurable attestation fields without any verification.
contract MockNitroProver is INitroProver {
  bytes public enclaveKey;
  bytes public userData;
  bytes public rawPcrs;

  constructor(bytes memory _enclaveKey, bytes memory _userData, bytes memory _rawPcrs) {
    enclaveKey = _enclaveKey;
    userData = _userData;
    rawPcrs = _rawPcrs;
  }

  function setRawPcrs(bytes calldata _rawPcrs) external {
    rawPcrs = _rawPcrs;
  }

  function verifyAttestation(
    bytes calldata,
    uint256
  ) external view returns (bytes memory, bytes memory, bytes memory) {
    return (enclaveKey, userData, rawPcrs);
  }

  function verifyAttestationStep1(
    bytes calldata,
    uint256
  )
    external
    view
    returns (bytes[] memory, bytes memory, bytes[] memory, bytes memory, bytes memory, bytes memory)
  {
    bytes[] memory cabundle = new bytes[](1);
    return (new bytes[](0), bytes(''), cabundle, enclaveKey, userData, rawPcrs);
  }

  function verifyAttestationStep2(
    bytes[] calldata,
    uint256,
    bytes calldata
  ) external pure returns (bytes memory) {
    return bytes('');
  }

  function verifyAttestationStep3(
    bytes[] calldata,
    bytes calldata,
    bytes calldata
  ) external pure returns (bytes memory, bytes memory, bytes memory) {
    return (bytes(''), bytes(''), bytes(''));
  }

  function verifyAttestationStep4(bytes calldata, bytes calldata, bytes calldata) external pure {}
}
