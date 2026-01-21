// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
interface INitroProver {
    //return (enclaveKey, userData, rawPcrs);
    function verifyAttestation(bytes calldata attestation, uint256 maxAge) external view returns(bytes memory, bytes memory, bytes memory);

    function verifyAttestationStep1(bytes calldata attestation, uint256 max_age) external view returns(bytes[] memory attestation_decoded, bytes memory certificate, bytes[] memory cabundle, bytes memory enclaveKey, bytes memory userData, bytes memory rawPcrs);
    //invoke step 2 for each index on the cabundle length. for index = 0, parentPubKey is empty
    function verifyAttestationStep2(bytes[] calldata cabundle, uint256 index, bytes calldata parentPubKey) external view returns(bytes memory pubKey);
    //invoke step 3 after finishing all the 2
    function verifyAttestationStep3(bytes[] calldata attestation_decoded, bytes calldata certificate, bytes calldata parentPubKey) external view returns(bytes memory attestationSig, bytes memory pubKey, bytes memory bufBufBuf);
    function verifyAttestationStep4(bytes calldata attestationSig, bytes calldata pubKey, bytes calldata buf) external view;
}
