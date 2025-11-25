// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "@openzeppelin/contracts/access/Ownable.sol";

contract AuthorityRegistry is Ownable {

    mapping(uint64 => mapping(address => bool)) allowedAuthorities;

    //events
    event AddedAuthority(uint64 indexed applicationId, address indexed authority);
    event RemovedAuthority(uint64 indexed applicationId, address indexed authority);

    //error
    error AuthorityNotPresent();
    error AuthorityAlreadyPresent();

    constructor(address owner) Ownable(owner) {}

    function addAllowedAuthority(uint64 applicationId, address authority) public onlyOwner {
        if(allowedAuthorities[applicationId][authority]) revert AuthorityAlreadyPresent();
        allowedAuthorities[applicationId][authority] = true;
        emit AddedAuthority(applicationId, authority);
    }

    function removeAllowedAuthority(uint64 applicationId, address authority) public onlyOwner {
        if(!allowedAuthorities[applicationId][authority]) revert AuthorityNotPresent();
        allowedAuthorities[applicationId][authority] = false;
        emit RemovedAuthority(applicationId, authority);
    }

    function checkAuthorityIsAllowed(uint64 applicationId, address authority) public view returns(bool) {
        return allowedAuthorities[applicationId][authority];
    }
}