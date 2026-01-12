// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "../AbstractTeeAuthenticator.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract NoAttestationTeeAuthenticator is AbstractTeeAuthenticator, Ownable {    
    address public teeSigner;
    bytes public pubSecp521r1;

    //events
    event TeeUpdate(address oldTee, address newTee, bytes oldPubSecp521r1, bytes newPubSecp521r1);

    //error
    error TeeAddressCantBeZero();
    error InvalidPKLength();

    constructor(address owner, address _teeSigner, bytes memory _pubSecp521r1) Ownable(owner) {
        teeSigner = _teeSigner;
        pubSecp521r1 = _pubSecp521r1;
        emit TeeUpdate(address(0), teeSigner, bytes(""), pubSecp521r1);
    }

    function updateTee(address newTeeSigner, bytes memory newPubSecp521r1) public onlyOwner {
        if(newTeeSigner == address(0)) revert TeeAddressCantBeZero();
        if(newPubSecp521r1.length != PK_LENGTH) revert InvalidPKLength();

        emit TeeUpdate(teeSigner, newTeeSigner, pubSecp521r1, newPubSecp521r1);
        teeSigner = newTeeSigner;
        pubSecp521r1 = newPubSecp521r1;
    }

    function getTeeSigner() public view override returns(address) {
        return teeSigner;
    }
    function getPubSecp521r1() public view override returns(bytes memory) {
        return pubSecp521r1;
    }
}