// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import '../ProcessorEndpoint.sol';

contract FallbackFailure {
  // Direct Ether transfer: always fails
  receive() external payable {
    revert('Receive not allowed');
  }

  // If no function matches: always fails
  fallback() external payable {
    revert('Fallback not allowed');
  }

  function insertRequestOnProcessorEndpoint(
    ProcessorEndpoint processorEndpoint,
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 maxFeeValue
  ) external payable {
    processorEndpoint.submitRequest{value: assetAmount + maxFeeValue}(
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
