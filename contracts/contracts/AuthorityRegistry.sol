// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "@openzeppelin/contracts/access/Ownable.sol";

contract AuthorityRegistry is Ownable {

    mapping(uint256 => mapping(address => bool)) allowedAuthorities;

    //events
    event AddedAuthority(uint256 indexed applicationId, address indexed authority);
    event RemovedAuthority(uint256 indexed applicationId, address indexed authority);

    //error
    error AuthorityNotPresent();
    error AuthorityAlreadyPresent();

    constructor() Ownable(msg.sender) {}

    function addAllowedAuthority(uint256 applicationId, address authority) public onlyOwner {
        if(allowedAuthorities[applicationId][authority]) revert AuthorityAlreadyPresent();
        allowedAuthorities[applicationId][authority] = true;
        emit AddedAuthority(applicationId, authority);
    }

    function removeAllowedAuthority(uint256 applicationId, address authority) public onlyOwner {
        if(!allowedAuthorities[applicationId][authority]) revert AuthorityNotPresent();
        allowedAuthorities[applicationId][authority] = false;
        emit RemovedAuthority(applicationId, authority);
    }

    function checkAuthorityIsAllowed(uint256 applicationId, address authority) public view returns(bool) {
        return allowedAuthorities[applicationId][authority];
    }
}