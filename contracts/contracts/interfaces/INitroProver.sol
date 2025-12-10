// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface INitroProver {
    /** max_age decides the tolerance after the expiration of the attestation to be still considered valid
     * 
     * returns (enclaveKey, userData, rawPcrs);
     */
    function verifyAttestation(bytes memory attestation, uint256 max_age) external view returns(bytes memory, bytes memory, bytes memory);

}