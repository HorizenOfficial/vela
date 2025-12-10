// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface INitroValidator {
    /** max_age decides the tolerance after the expiration of the attestation to be still considered valid
     *  this performs both the decodeAttestationTbs step and the verify
     * returns (bytes32 pcr0, bytes memory publicKey);
     */
    function verifyAttestation(bytes calldata attestation, uint256 maxAge) external returns (bytes32, bytes memory);
}