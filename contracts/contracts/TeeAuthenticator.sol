// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "./interfaces/ITeeAuthenticator.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract TeeAuthenticator is ITeeAuthenticator, Ownable {

    address public teeSigner;

    //events
    event TeeUpdate(address oldTee, address newTee);

    //error
    error TeeAddressCantBeZero();
    error TeeIsNotSet();


    constructor(address owner, address _teeSigner) Ownable(owner) {
        teeSigner = _teeSigner;
        emit TeeUpdate(address(0), teeSigner);
    }

    function updateTee(address newTeeSigner) public onlyOwner {
        if(newTeeSigner == address(0)) revert TeeAddressCantBeZero();

        emit TeeUpdate(teeSigner, newTeeSigner);
        teeSigner = newTeeSigner;
    }

    function checkSignature(
        uint256 applicationId,
        bytes calldata prevStateRoot,
        bytes calldata newStateRoot,
        uint256 processedRequestId,
        bytes[] memory events,
        Structs.WithdrawalRequest[] memory withdrawalRequests,
        bytes calldata signature
    ) external view override returns (bool) {
        if(teeSigner == address(0)) revert TeeIsNotSet();

        bytes32 messageHash = keccak256(abi.encode(
            applicationId,
            prevStateRoot,
            newStateRoot,
            processedRequestId,
            events,
            withdrawalRequests
        ));

        address recovered = ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(messageHash), signature);
        return recovered == teeSigner;
    }
}