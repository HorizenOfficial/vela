// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../Structs.sol";

/// @title TeeAuthenticator interface
/// @notice External API and errors for tee signature verification.
interface ITeeAuthenticator {
    /// @notice Tee signer or public key is not configured.
    error TeeIsNotSet();

    /// @notice Verifies an update signature for a processed request.
    /// @param applicationId Application identifier.
    /// @param prevStateRoot Previous state root.
    /// @param newStateRoot New state root.
    /// @param processedRequestId Request identifier being processed.
    /// @param events Encrypted event payloads.
    /// @param eventSubTypes Event subtype labels.
    /// @param withdrawalRequests Withdrawal requests to execute.
    /// @param refundAmount Refund amount to the request sender.
    /// @param applicationFee Fee amount to the collector.
    /// @param signature Signature over the update data.
    /// @return valid True if the signature is valid.
    function checkSignature(
        uint64 applicationId, 
        bytes32 prevStateRoot, 
        bytes32 newStateRoot, 
        bytes32 processedRequestId,
        bytes[] calldata events,
        string[] calldata eventSubTypes,
        Structs.WithdrawalRequest[] calldata withdrawalRequests, 
        uint256 refundAmount, 
        uint256 applicationFee,
        bytes calldata signature
    ) external view returns(bool);

    /// @notice Returns the configured tee signer address.
    /// @return signer Tee signer address.
    function getTeeSigner() external view returns(address);

    /// @notice Returns the configured secp521r1 public key.
    /// @return pubKey Raw public key bytes.
    function getPubSecp521r1() external view returns(bytes memory);
}
