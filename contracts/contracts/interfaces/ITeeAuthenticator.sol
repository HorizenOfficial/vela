// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import '../Structs.sol';

/// @title TeeAuthenticator interface
/// @notice External API and errors for tee signature verification.
interface ITeeAuthenticator {
  /// @notice Tee signer or public key is not configured.
  error TeeIsNotSet();

  /// @notice Batch signature verification was called with no entry hashes.
  error EmptyBatch();

  /// @notice Verifies an update signature for a processed request.
  ///
  /// @param params Parameters of the signature.
  /// @param signature Signature over the update data.
  /// @return valid True if the signature is valid.
  function checkSignature(
    Structs.SignatureParams memory params,
    bytes memory signature
  ) external view returns (bool);

  /// @notice Verifies a single signature covering a batch of processed requests.
  ///
  /// @dev The signed message is the EIP-191 personal_sign digest of the concatenated
  ///      entry hashes, with a dynamic `32 * entryHashes.length` length prefix and no
  ///      extra hashing layer. A single-entry batch is byte-identical to the digest
  ///      verified by {checkSignature}.
  ///
  /// @param entryHashes Per-entry hashes, in submission order (see UpdateEntryHash).
  /// @param signature Signature over the batch digest.
  /// @return valid True if the signature is valid.
  function checkBatchSignature(
    bytes32[] calldata entryHashes,
    bytes calldata signature
  ) external view returns (bool);

  /// @notice Returns the configured tee signer address.
  /// @return signer Tee signer address.
  function getTeeSigner() external view returns (address);

  /// @notice Returns the configured secp521r1 public key.
  /// @return pubKey Raw public key bytes.
  function getPubSecp521r1() external view returns (bytes memory);
}
