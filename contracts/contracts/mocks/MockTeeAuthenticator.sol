// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../interfaces/ITeeAuthenticator.sol";

contract MockTeeAuthenticator is ITeeAuthenticator {

    bool response;

    constructor(bool _response) {
        response = _response;
    }

    function checkSignature(uint8, bytes calldata, bytes calldata, bytes[] memory, Structs.WithdrawalRequest[] memory, bytes memory) external view override returns(bool) {
        return response;
    }
}