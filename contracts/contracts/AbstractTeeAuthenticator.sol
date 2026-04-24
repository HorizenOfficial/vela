// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import './interfaces/ITeeAuthenticator.sol';
import './Structs.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';
import '@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol';

/// @title AbstractTeeAuthenticator
/// @notice Base implementation for tee signature verification.
abstract contract AbstractTeeAuthenticator is ITeeAuthenticator {
  uint256 public constant PK_LENGTH = 133; //secp521r1 uncompressed public key length in bytes

  /// @inheritdoc ITeeAuthenticator
  function getTeeSigner() public view virtual override returns (address);
  /// @inheritdoc ITeeAuthenticator
  function getPubSecp521r1() public view virtual override returns (bytes memory);

  /// @inheritdoc ITeeAuthenticator
  function checkSignature(
    Structs.SignatureParams memory params,
    bytes memory signature
  ) external view returns (bool) {
    if (getTeeSigner() == address(0) || getPubSecp521r1().length != PK_LENGTH) revert TeeIsNotSet();

    bytes32 userEventsHash = keccak256(abi.encode(params.userEvents.events));
    bytes32 userEventSubTypesHash = keccak256(abi.encode(params.userEvents.subTypes));
    bytes32 appEventsHash = keccak256(abi.encode(params.appEvents.events));
    bytes32 appEventSubTypesHash = keccak256(abi.encode(params.appEvents.subTypes));
    bytes32 withdrawalRequestsHash = keccak256(abi.encode(params.withdrawalRequests));

    bytes32 messageHash = keccak256(
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

    address recovered = ECDSA.recover(
      MessageHashUtils.toEthSignedMessageHash(messageHash),
      signature
    );
    return recovered == getTeeSigner();
  }
}
