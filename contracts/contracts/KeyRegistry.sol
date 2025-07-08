// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

contract KeyRegistry {
    uint256 public constant PK_LENGTH = 133;

    mapping(address => bytes) public registry;
    //events
    event PKRegistered(address owner, bytes publicKey);
    //errors
    error InvalidLength();

    function registerPK(bytes calldata publicKey) public {
        if(publicKey.length != PK_LENGTH) revert InvalidLength();
        registry[msg.sender] = publicKey;
        emit PKRegistered(msg.sender, publicKey);
    }

    function getPK(address owner) public view returns(bytes memory) {
        return registry[owner];
    }
}