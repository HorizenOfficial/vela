// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface INitroProver {
    //return (enclaveKey, userData, rawPcrs);
    function verifyAttestation(bytes memory attestation, uint256 maxAge) external view returns(bytes memory, bytes memory, bytes memory);

    function decodeAttestation(bytes memory attestation) external view returns(bytes[] memory);


    function verifyAttestationFromDecoded(bytes[] memory attestation_decoded, uint256 max_age) external view returns(bytes memory, bytes memory, bytes memory);

    //return (enclaveKey, userData, rawPcrs, attestationSig, pubKey, buf)
    function verifyAttestationStep1(bytes calldata attestation, uint256 max_age) external view returns(bytes memory, bytes memory, bytes memory, bytes memory, bytes memory, bytes memory);
    function verifyAttestationStep2(bytes memory attestationSig, bytes memory pubKey, bytes memory buf) external view;
}