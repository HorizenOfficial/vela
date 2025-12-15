// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface INitroProver {
    /** max_age decides the tolerance after the expiration of the attestation to be still considered valid
     *  
     *  return (enclaveKey, userData, rawPcrs);
     */
    function verifyAttestation(bytes memory attestation, uint256 maxAge) external view returns(bytes memory, bytes memory, bytes memory);
    /** max_age decides the tolerance after the expiration of the attestation to be still considered valid
     *  
     *  return (enclaveKey, userData);
     */
    function verifyAttestation(bytes memory attestation, bytes memory PCRs, uint256 maxAge) external view returns(bytes memory, bytes memory);

    /** max_age decides the tolerance after the expiration of the attestation to be still considered valid
     *  
     *  return (enclaveKey, userData, decodedPcrs);
     */
    function verifyAttestationAndDecodePCR(bytes memory attestation, uint256 max_age) external view returns(bytes memory, bytes memory, bytes[2][] memory);

}