// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';

import './interfaces/IProcessorEndpoint.sol';
import './interfaces/ITokenAllowlist.sol';
import './Structs.sol';
import './ProcessorEndpointStorage.sol';

/// @title ProcessorEndpointExtension
/// @notice Code moved out of `ProcessorEndpoint` to stay under the EIP-170 deployed-bytecode
///         limit. It is never called directly: `ProcessorEndpoint` forwards matching calls here
///         with `delegatecall`, so this code runs against the endpoint's storage, balance,
///         `msg.sender` and `msg.value`.
/// @dev Currently holds the facilitator path (`submitRequestFor`). Because the boundary is a
///      `delegatecall` onto the shared `ProcessorEndpointStorage` layout, arguments are not
///      re-encoded and internal helpers can be reused, unlike an external-call or linked-library
///      split. Helpers this contract needs are duplicated into its bytecode, which costs nothing:
///      only the endpoint is near the size limit.
contract ProcessorEndpointExtension is ProcessorEndpointStorage {
  /// @dev This contract's own address, fixed at deployment. Under `delegatecall`,
  ///      `address(this)` is the endpoint instead, which is the only context where the code is
  ///      meaningful: called directly it would read and write this contract's own (empty) storage
  ///      and could strand ETH here.
  address private immutable _self = address(this);

  /// @notice Reverts when the extension is called directly rather than through the endpoint.
  error DirectCallNotAllowed();

  modifier onlyDelegateCall() {
    if (address(this) == _self) revert DirectCallNotAllowed();
    _;
  }

  /// @notice Submits a request on behalf of `sender`, authorized by an EIP-712 signature, with
  ///         the calling facilitator paying the gas and the fee. See
  ///         `IProcessorEndpoint.submitRequestFor`, whose declaration and documentation this
  ///         implements; `ProcessorEndpoint.submitRequestFor` is the entry point callers use.
  function submitRequestFor(
    address sender,
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 deadline,
    bytes calldata requestSignature,
    bytes calldata depositPermit
  )
    external
    payable
    onlyDelegateCall
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    nonReentrant
    returns (bytes32)
  {
    // 1. Only ASSOCIATEKEY and PROCESS are supported. TRUSTPROCESS is trusted
    //    and can ONLY be created internally during stateUpdate (via a trigger's
    //    getTrustProcessPayload payload) — never submitted by an external caller here.
    if (
      requestType != Structs.RequestType.ASSOCIATEKEY && requestType != Structs.RequestType.PROCESS
    ) revert IProcessorEndpoint.InvalidRequestType();

    // 2. Verify deadline not expired
    if (block.timestamp > deadline) revert IProcessorEndpoint.DeadlineExpired();

    // 3. Check queue size
    if (_pendingRequestsSize() >= maxQueueSize) revert IProcessorEndpoint.QueueThresholdExceeded();

    // 4. Validate payload for ASSOCIATEKEY
    if (requestType == Structs.RequestType.ASSOCIATEKEY) {
      if (payload.length != 133 && payload.length != 226)
        revert IProcessorEndpoint.InvalidPayload();
    }

    // 5. Read current nonce and build EIP-712 hash
    uint256 nonce = facilitatorNonces[sender];
    bytes32 structHash = keccak256(
      abi.encode(
        REQUEST_AUTHORIZATION_TYPEHASH,
        sender,
        protocolVersion,
        applicationId,
        uint8(requestType),
        keccak256(payload),
        tokenAddress,
        assetAmount,
        nonce,
        deadline
      )
    );
    bytes32 digest = _hashTypedDataV4(structHash);

    // 6. Recover user address from EIP-712 request signature and verify.
    //    tryRecover, not recover: a signature ECDSA cannot parse must revert with
    //    InvalidSignature, which IProcessorEndpoint declares. ECDSA.recover would raise
    //    ECDSAInvalidSignature/-Length/-S instead, and those are absent from the endpoint's ABI —
    //    this code runs by delegatecall, so callers only ever see the endpoint — leaving clients
    //    with revert data they cannot decode.
    (address recoveredSigner, ECDSA.RecoverError recoverError, ) = ECDSA.tryRecover(
      digest,
      requestSignature
    );
    if (recoverError != ECDSA.RecoverError.NoError || recoveredSigner == address(0))
      revert IProcessorEndpoint.InvalidSignature();
    if (recoveredSigner != sender) revert IProcessorEndpoint.InvalidSigner();

    // 7. Consume nonce (replay protection)
    facilitatorNonces[sender] = nonce + 1;

    // 8. Validate fee
    uint256 maxFeeValue = msg.value;
    if (maxFeeValue < minFeePerRequest) revert IProcessorEndpoint.FeeValueBelowMinimum();

    // 9. Validate token and handle deposit
    if (assetAmount == 0 && tokenAddress != ETH_TOKEN) revert IProcessorEndpoint.InvalidValue();
    if (assetAmount > 0) {
      if (tokenAddress == ETH_TOKEN) revert IProcessorEndpoint.InvalidValue();
      if (!tokenAllowlist.isAllowedToken(tokenAddress)) revert ITokenAllowlist.TokenNotAllowed();

      // Decode deposit permit and execute EIP-2612 permit + transferFrom
      if (depositPermit.length != 96) revert IProcessorEndpoint.InvalidPermit();
      (uint8 v, bytes32 r, bytes32 s) = abi.decode(depositPermit, (uint8, bytes32, bytes32));

      // Check current allowance before calling permit
      IERC20 token = IERC20(tokenAddress);
      if (token.allowance(sender, address(this)) < assetAmount) {
        IERC20Permit(tokenAddress).permit(sender, address(this), assetAmount, deadline, v, r, s);
      }

      _pullERC20(tokenAddress, sender, assetAmount);

      _addToCustody(applicationId, tokenAddress, assetAmount);
    }

    // 10. Create PendingRequest with sender = user (not msg.sender) and enqueue
    return
      _enqueueRequest(
        sender,
        msg.sender,
        protocolVersion,
        applicationId,
        requestType,
        payload,
        tokenAddress,
        assetAmount,
        maxFeeValue
      );
  }
}
