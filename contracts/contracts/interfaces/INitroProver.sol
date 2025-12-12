// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface INitroProver {
    /** max_age decides the tolerance after the expiration of the attestation to be still considered valid
     *  
     *  return (enclaveKey, userData, rawPcrs);
     */
    function verifyAttestation(bytes calldata attestation, uint256 maxAge) external view returns(bytes memory, bytes memory, bytes memory);
}