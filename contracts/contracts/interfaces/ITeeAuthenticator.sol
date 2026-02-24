// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import '../Structs.sol';

/// @title TeeAuthenticator interface
/// @notice External API and errors for tee signature verification.
interface ITeeAuthenticator {
  /// @notice Tee signer or public key is not configured.
  error TeeIsNotSet();

  /// @notice Verifies an update signature for a processed request.
  ///
  /// @param params Parameters of the signature.
  /// @param signature Signature over the update data.
  /// @return valid True if the signature is valid.
  function checkSignature(
    Structs.SignatureParams memory params,
    bytes memory signature
  ) external view returns (bool);
  /// @notice Returns the configured tee signer address.
  /// @return signer Tee signer address.
  function getTeeSigner() external view returns (address);

  /// @notice Returns the configured secp521r1 public key.
  /// @return pubKey Raw public key bytes.
  function getPubSecp521r1() external view returns (bytes memory);
}
