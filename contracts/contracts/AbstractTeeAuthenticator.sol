// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './interfaces/ITeeAuthenticator.sol';
import './Structs.sol';
import './UpdateEntryHash.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';
import '@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol';

/// @title AbstractTeeAuthenticator
/// @notice Base implementation for tee signature verification.
abstract contract AbstractTeeAuthenticator is ITeeAuthenticator {
  uint256 public constant PK_LENGTH = 133; //secp521r1 uncompressed public key length in bytes

  /// @inheritdoc ITeeAuthenticator
  function getTeeSigner() public view virtual override returns (address);
  /// @inheritdoc ITeeAuthenticator
  function getPubSecp521r1() public view virtual override returns (bytes memory);

  /// @inheritdoc ITeeAuthenticator
  function checkSignature(
    Structs.SignatureParams memory params,
    bytes memory signature
  ) external view returns (bool) {
    if (getTeeSigner() == address(0) || getPubSecp521r1().length != PK_LENGTH) revert TeeIsNotSet();

    bytes32 messageHash = UpdateEntryHash.entryHash(params);

    address recovered = ECDSA.recover(
      MessageHashUtils.toEthSignedMessageHash(messageHash),
      signature
    );
    return recovered == getTeeSigner();
  }

  /// @inheritdoc ITeeAuthenticator
  function checkBatchSignature(
    bytes32[] calldata entryHashes,
    bytes calldata signature
  ) external view returns (bool) {
    if (getTeeSigner() == address(0) || getPubSecp521r1().length != PK_LENGTH) revert TeeIsNotSet();
    if (entryHashes.length == 0) revert EmptyBatch();

    // EIP-191 personal_sign over the concatenated entry hashes, with the length
    // prefix built at runtime from the total byte length (32 * N). There is no
    // intermediate keccak256 over the concatenation: entry hashes are fixed
    // 32-byte values, so the concatenation splits exactly one way, and the
    // prefix binds the batch size. A 1-entry batch therefore produces the same
    // digest as checkSignature() over that entry.
    address recovered = ECDSA.recover(
      MessageHashUtils.toEthSignedMessageHash(abi.encodePacked(entryHashes)),
      signature
    );
    return recovered == getTeeSigner();
  }
}
