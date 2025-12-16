// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface INitroProver {
    //return (enclaveKey, userData, rawPcrs);
    function verifyAttestation(bytes memory attestation, uint256 maxAge) external view returns(bytes memory, bytes memory, bytes memory);

    function verifyAttestationStep1(bytes calldata attestation, uint256 max_age) external view returns(bytes[] memory attestation_decoded, bytes memory certificate, bytes[] memory cabundle, bytes memory enclaveKey, bytes memory userData, bytes memory rawPcrs);
    //invoke step 2 for each index on the cabundle length. for index = 0, parentPubKey is empty
    function verifyAttestationStep2(bytes[] memory cabundle, uint256 index, bytes memory parentPubKey) external view returns(bytes memory pubKey);
    //invoke step 3 after finishing all the 2
    function verifyAttestationStep3(bytes[] memory attestation_decoded, bytes memory certificate, bytes memory parentPubKey) external view returns(bytes memory attestationSig, bytes memory pubKey, bytes memory bufBufBuf);
    function verifyAttestationStep4(bytes memory attestationSig, bytes memory pubKey, bytes memory buf) external view;
}