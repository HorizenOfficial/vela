// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '../interfaces/ITrigger.sol';
import '../Structs.sol';

/// @title Trigger abstract contract
/// @notice Base contract for triggers invoked by the ProcessorEndpoint on request completion.
///         Provides a non-overridable withdraw mechanism that sweeps all allowlisted ERC-20
///         tokens and ETH back to the ProcessorEndpoint, and an overridable execute hook.
abstract contract AbstractTrigger is ITrigger {
  /// @notice ProcessorEndpoint that is allowed to call execute and withdraw.
  IProcessorEndpoint public processorEndpoint;
  /// @notice Allowlist of ERC-20 tokens swept on withdraw.
  ITokenAllowlist public tokenAllowlist;

  /// @dev Reverts if msg.sender is not the processorEndpoint.
  modifier _onlyProcessorEndpoint() {
    if (msg.sender != address(processorEndpoint)) {
      revert NotProcessorEndpoint();
    }
    _;
  }

  receive() external payable {}

   /// @param _processorEndpoint ProcessorEndpoint that will call execute and withdraw.

  /// @param _processorEndpoint ProcessorEndpoint that will call execute and withdraw.
  /// @param _tokenAllowlist Allowlist used to determine which ERC-20 tokens are swept on withdraw.
  constructor(IProcessorEndpoint _processorEndpoint, ITokenAllowlist _tokenAllowlist) {
    if (address(_processorEndpoint) == address(0) || address(_tokenAllowlist) == address(0)) {
      revert ZeroAddress();
    }
    processorEndpoint = _processorEndpoint;
    tokenAllowlist = _tokenAllowlist;
  }

  /// @notice Called by the ProcessorEndpoint when a request is completed.
  ///         Delegates to _execute and emits TriggerExecuted.
  /// @param applicationId Application identifier.
  /// @param prevStateRoot Previous state root.
  /// @param newStateRoot New state root.
  /// @param processedRequestId Request identifier being processed.
  /// @param userEventData Encrypted user events and their subtypes.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param withdrawalRequests Withdrawal requests executed in the state update.
  /// @param refund Refund amount sent to the request sender.
  /// @param applicationFees Fee amount sent to the collector.
  /// @param errorCode Error code for the update.
  /// @param errorMsg Error message for the update.
  function execute(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg
  ) public _onlyProcessorEndpoint {
    _execute(
      applicationId,
      prevStateRoot,
      newStateRoot,
      processedRequestId,
      userEventData,
      appEventData,
      withdrawalRequests,
      refund,
      applicationFees,
      errorCode,
      errorMsg
    );
    emit TriggerExecuted(applicationId, processedRequestId);
  }

  /// @notice Sweeps all allowlisted ERC-20 tokens and ETH held by this contract to the
  ///         ProcessorEndpoint. Can only be called by the ProcessorEndpoint.
  ///         Cannot be overridden; custom post-sweep logic belongs in _postWithdraw.
  ///         Returns false (without reverting) if _postWithdraw reverts.
  /// @param applicationId Application identifier.
  /// @param prevStateRoot Previous state root.
  /// @param newStateRoot New state root.
  /// @param processedRequestId Request identifier being processed.
  /// @param userEventData Encrypted user events and their subtypes.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param withdrawalRequests Withdrawal requests executed in the state update.
  /// @param refund Refund amount sent to the request sender.
  /// @param applicationFees Fee amount sent to the collector.
  /// @param errorCode Error code for the update.
  /// @param errorMsg Error message for the update.
  /// @return success True if _postWithdraw executed successfully, false if it reverted.
  function withdraw(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg
  ) public _onlyProcessorEndpoint returns (bool, ReturnedTokens[] memory) {
    address[] memory tokens = tokenAllowlist.getAllowedTokens();
    uint256 length = tokens.length;

    //array to trace moved tokens
    ReturnedTokens[] memory returned = new ReturnedTokens[](length + 1);
    uint256 i;

    //return ERC20 tokens in allowlist
    while (i != length) {
      address token = tokens[i];
      uint256 balance = IERC20(token).balanceOf(address(this));
      if (balance > 0) {
        IERC20(token).transfer(address(processorEndpoint), balance);
      }
      returned[i] = ReturnedTokens({token: token, amount: balance});
      unchecked {
        ++i;
      }
    }

    //return ETH
    uint256 ethBalance = address(this).balance;
    if (ethBalance > 0) {
      payable(address(processorEndpoint)).transfer(ethBalance);
    }
    //put in last position in returned
    returned[returned.length - 1] = ReturnedTokens({token: ETH_TOKEN, amount: ethBalance});


    // try post withdraw hook, return false if it reverts
    try
      this._postWithdraw(
        applicationId,
        prevStateRoot,
        newStateRoot,
        processedRequestId,
        userEventData,
        appEventData,
        withdrawalRequests,
        refund,
        applicationFees,
        errorCode,
        errorMsg,
        returned
      )
    {
      return (true, returned);
    } catch {
      return (false, returned);
    }
  }

  /// @notice Override to implement custom execute behavior.
  ///         Called inside execute, which is restricted to the ProcessorEndpoint.
  /// @param applicationId Application identifier.
  /// @param prevStateRoot Previous state root.
  /// @param newStateRoot New state root.
  /// @param processedRequestId Request identifier being processed.
  /// @param userEventData Encrypted user events and their subtypes.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param withdrawalRequests Withdrawal requests executed in the state update.
  /// @param refund Refund amount sent to the request sender.
  /// @param applicationFees Fee amount sent to the collector.
  /// @param errorCode Error code for the update.
  /// @param errorMsg Error message for the update.
  function _execute(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg
  ) internal virtual;

  /// @notice Override to add logic that runs after the default withdraw sweep.
  ///         Called at the end of withdraw, after all tokens and ETH have been transferred.
  /// @param applicationId Application identifier.
  /// @param prevStateRoot Previous state root.
  /// @param newStateRoot New state root.
  /// @param processedRequestId Request identifier being processed.
  /// @param userEventData Encrypted user events and their subtypes.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param withdrawalRequests Withdrawal requests executed in the state update.
  /// @param refund Refund amount sent to the request sender.
  /// @param applicationFees Fee amount sent to the collector.
  /// @param errorCode Error code for the update.
  /// @param errorMsg Error message for the update.
  /// @param returnedTokens Array of tokens actually transferred in the sweep, with amounts.
  function _postWithdraw(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg,
    ReturnedTokens[] memory returnedTokens
  ) public virtual;
}
