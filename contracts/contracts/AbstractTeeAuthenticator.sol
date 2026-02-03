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
  ) external view returns (bool) {
    if (getTeeSigner() == address(0) || getPubSecp521r1().length != PK_LENGTH) revert TeeIsNotSet();

    bytes32 eventsHash = keccak256(abi.encode(events));
    bytes32 eventSubTypesHash = keccak256(abi.encode(eventSubTypes));
    bytes32 withdrawalRequestsHash = keccak256(abi.encode(withdrawalRequests));

    bytes32 messageHash = keccak256(
      abi.encode(
        applicationId,
        prevStateRoot,
        newStateRoot,
        processedRequestId,
        eventsHash,
        eventSubTypesHash,
        withdrawalRequestsHash,
        refundAmount,
        applicationFee
      )
    );

    address recovered = ECDSA.recover(
      MessageHashUtils.toEthSignedMessageHash(messageHash),
      signature
    );
    return recovered == getTeeSigner();
  }
}
