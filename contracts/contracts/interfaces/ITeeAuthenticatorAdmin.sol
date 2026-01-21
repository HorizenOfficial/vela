// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

/// @title TeeAuthenticator admin interface
/// @notice Admin API, events, and errors for tee updates and attestation flows.
interface ITeeAuthenticatorAdmin {
    /// @notice Emitted when the tee signer or public key is updated.
    /// @param oldTee Previous tee signer.
    /// @param newTee New tee signer.
    /// @param oldPubSecp521r1 Previous public key.
    /// @param newPubSecp521r1 New public key.
    event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1);

    /// @notice Emitted when PCR0 is updated.
    /// @param oldPcr0 Previous PCR0.
    /// @param newPcr0 New PCR0.
    event PcrZeroUpdate(bytes indexed oldPcr0, bytes indexed newPcr0);

    /// @notice PCR0 does not match expected value.
    error InvalidPCR();

    /// @notice Public key length is invalid.
    error InvalidPKLength();

    /// @notice User data length is invalid.
    error InvalidUserDataLength();

    /// @notice Attestation has already been used.
    error AttestationAlreadyUsed();

    /// @notice Wrong step in the update flow.
    error WrongStep();

    /// @notice Updates the tee signer and public key using a full attestation.
    /// @param attestation Attestation blob.
    function updateTee(bytes calldata attestation) external;

    /// @notice Starts a step-by-step tee update by parsing attestation.
    /// @param attestation Attestation blob.
    function updateTeeStep1(bytes calldata attestation) external;

    /// @notice Advances step 2 of the tee update flow.
    function updateTeeStep2() external;

    /// @notice Advances step 3 of the tee update flow.
    function updateTeeStep3() external;

    /// @notice Completes step 4 of the tee update flow.
    function updateTeeStep4() external;

    /// @notice Returns the total number of step 2 iterations required.
    /// @return length Number of step 2 calls needed.
    function getStep2TotalLength() external view returns (uint256);

    /// @notice Updates the PCR0 value.
    /// @param newPcr0 New PCR0 value.
    function updatePcr0(bytes calldata newPcr0) external;
}
