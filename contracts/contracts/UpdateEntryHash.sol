// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './Structs.sol';

/// @title UpdateEntryHash
/// @notice Per-entry digest of a TEE update payload.
/// @dev Single source of truth for the entry hash: the TEE authenticator uses it to
///      rebuild the signed message, and the processor endpoint uses it to produce the
///      entry hashes covered by a batch signature.
///      Field order must stay in sync with `MsgToSignBuilder.buildEntryHash`
///      (`pkg/executor/msgtosign_builder.go`).
///      The result is always 32 bytes, which is what makes the concatenation of entry
///      hashes in a batch digest unambiguous (see `checkBatchSignature`).
library UpdateEntryHash {
  /// @notice Computes the entry hash for a single update payload.
  /// @param params Update payload fields.
  /// @return entryHash keccak256 of the ABI-encoded payload, without the EIP-191 prefix.
  function entryHash(Structs.SignatureParams memory params) internal pure returns (bytes32) {
    bytes32 userEventsHash = keccak256(abi.encode(params.userEvents.events));
    bytes32 userEventSubTypesHash = keccak256(abi.encode(params.userEvents.subTypes));
    bytes32 appEventsHash = keccak256(abi.encode(params.appEvents.events));
    bytes32 appEventSubTypesHash = keccak256(abi.encode(params.appEvents.subTypes));
    bytes32 withdrawalRequestsHash = keccak256(abi.encode(params.withdrawalRequests));

    return
      keccak256(
        abi.encode(
          params.applicationId,
          params.prevStateRoot,
          params.newStateRoot,
          params.processedRequestId,
          userEventsHash,
          userEventSubTypesHash,
          appEventsHash,
          appEventSubTypesHash,
          withdrawalRequestsHash,
          params.refundAmount,
          params.applicationFee,
          params.errorCode,
          params.errorMsg
        )
      );
  }
}
