// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

contract KeyRegistry {

    mapping(address => bytes) public registry;
    //events
    event PKRegistered(address owner, bytes publicKey);

    function registerPK(bytes calldata publicKey) public {
        registry[msg.sender] = publicKey;
        emit PKRegistered(msg.sender, publicKey);
    }

    function getPK(address owner) public view returns(bytes memory) {
        return registry[owner];
    }
}