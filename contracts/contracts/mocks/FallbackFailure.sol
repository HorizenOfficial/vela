// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../ProcessorEndpoint.sol";

contract FallbackFailure {
    // Ricezione diretta di Ether: fallisce sempre
    receive() external payable {
        revert("Receive not allowed");
    }

    // Chiamata non corrispondente a funzione: fallisce sempre
    fallback() external payable {
        revert("Fallback not allowed");
    }

    function insertRequestOnProcessorEndpoint(
        ProcessorEndpoint processorEndpoint,
        uint8 protocolVersion, 
        uint256 applicationId, 
        Structs.RequestType requestType, 
        bytes calldata payload, 
        uint256 value
    ) payable public {
        processorEndpoint.submitRequest{value:value}(protocolVersion, applicationId, requestType, payload, value);
    }
}