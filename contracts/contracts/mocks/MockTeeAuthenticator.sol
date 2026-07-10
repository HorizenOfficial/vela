// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import '../interfaces/ITeeAuthenticator.sol';

contract MockTeeAuthenticator is ITeeAuthenticator {
  address public teeSigner;
  bytes public pubSecp521r1;
  // keccak256 of the PCR0 the platform should currently be running.
  // Mirrors TeeAuthenticator.activeImage; permissionless setter for test control.
  bytes32 public activeImage;

  constructor(address _teeSigner, bytes memory _pubSecp521r1) {
    teeSigner = _teeSigner;
    pubSecp521r1 = _pubSecp521r1;
  }

  function setActiveImage(bytes32 _activeImage) external {
    activeImage = _activeImage;
  }

  function checkSignature(
    Structs.SignatureParams memory /*params*/,
    bytes memory /*signature*/
  ) external pure override returns (bool) {
    return true; // Always return true for mock
  }

  function getTeeSigner() external view override returns (address) {
    return teeSigner;
  }
  function getPubSecp521r1() external view returns (bytes memory) {
    return pubSecp521r1;
  }
}
