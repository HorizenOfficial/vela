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

  /// @notice Emitted when a PCR0 swap is proposed.
  /// @param targetPcr0 Raw PCR0 the swap targets (preimage, non-indexed for off-chain audit).
  /// @param eta Earliest application time.
  event Pcr0SwapProposed(bytes targetPcr0, uint256 eta);

  /// @notice Emitted when a PCR0 swap is applied and the target becomes the active image.
  /// @param pcr0 Raw PCR0 that became active.
  event Pcr0Swapped(bytes pcr0);

  /// @notice Emitted when a pending PCR0 swap is cancelled.
  /// @param targetPcr0 Raw PCR0 the cancelled swap targeted.
  event Pcr0SwapCancelled(bytes targetPcr0);

  /// @notice Emitted when a PCR0 is removed from the accepted set.
  /// @param pcr0 Raw PCR0 that was removed.
  event Pcr0Removed(bytes pcr0);

  /// @notice PCR0 does not match expected value.
  error InvalidPCR();

  /// @notice PCR0 is not 48 bytes.
  error InvalidPcr0Length();

  /// @notice A live swap proposal already exists.
  error SwapAlreadyPending();

  /// @notice No swap proposal is pending.
  error NoPendingSwap();

  /// @notice The timelock has not elapsed yet.
  error TimelockNotElapsed();

  /// @notice The swap proposal application window has passed.
  error SwapProposalExpired();

  /// @notice PCR0 is not in the accepted set.
  error UnknownPcr0();

  /// @notice The active image cannot be removed.
  error CannotRemoveActiveImage();

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

  /// @notice Proposes a swap of the active image to `targetPcr0`.
  /// @dev Reverts while a live (non-expired) proposal exists; cancel it first.
  /// @param targetPcr0 Raw PCR0 of the target image (48 bytes).
  function proposePcr0Swap(bytes calldata targetPcr0) external;

  /// @notice Applies the pending swap: adds the target to the accepted set if new
  ///         (only after the timelock) and makes it the active image.
  /// @dev Must be called within the application window after the proposal's eta.
  function applyPcr0Swap() external;

  /// @notice Cancels the pending swap proposal.
  function cancelPcr0Swap() external;

  /// @notice Removes a PCR0 from the accepted set.
  /// @dev Cannot remove the active image, which also keeps the set non-empty.
  /// @param oldPcr0 Raw PCR0 to remove (48 bytes).
  function removePcr0(bytes calldata oldPcr0) external;

  /// @notice Returns the number of accepted PCR0 values.
  /// @return count Size of the accepted set.
  function getAcceptedPcr0Count() external view returns (uint256 count);
}
