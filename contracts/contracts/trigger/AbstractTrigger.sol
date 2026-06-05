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
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  function execute(Structs.EventData calldata appEventData) public _onlyProcessorEndpoint {
    _execute(appEventData);
    emit TriggerExecuted();
  }

  /// @notice Sweeps all allowlisted ERC-20 tokens and ETH held by this contract to the
  ///         ProcessorEndpoint. Can only be called by the ProcessorEndpoint.
  ///         Cannot be overridden; custom post-sweep logic belongs in getTrustProcessPayload.
  ///         Returns false (without reverting) if getTrustProcessPayload reverts.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param executeSuccess True if the execute call succeeded, false if it reverted.
  /// @return success True if getTrustProcessPayload executed successfully, false if it reverted.
  function withdraw(
    Structs.EventData calldata appEventData,
    bool executeSuccess
  )
    public
    _onlyProcessorEndpoint
    returns (bool, Structs.TokenAndAmount[] memory, Structs.TokenAndAmount[] memory, bytes memory)
  {
    address[] memory tokens = tokenAllowlist.getAllowedTokens();
    uint256 length = tokens.length;

    //arrays to trace moved tokens
    Structs.TokenAndAmount[] memory returned = new Structs.TokenAndAmount[](length + 1);
    Structs.TokenAndAmount[] memory failed = new Structs.TokenAndAmount[](length + 1);
    uint256 i;
    uint256 insertIntoReturned;
    uint256 insertIntoFailed;

    //return ERC20 tokens in allowlist
    while (i != length) {
      address token = tokens[i];
      uint256 balance = IERC20(token).balanceOf(address(this));

      if (balance > 0) {
        try IERC20(token).transfer(address(processorEndpoint), balance) {
          // try to transfer
          //if successful, add to returned
          returned[insertIntoReturned] = Structs.TokenAndAmount({token: token, amount: balance});
          unchecked {
            ++insertIntoReturned;
          }
        } catch {
          //if transfer fails, add to failed
          failed[insertIntoFailed] = Structs.TokenAndAmount({token: token, amount: balance});
          unchecked {
            ++insertIntoFailed;
          }
        }
      }

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
    returned[insertIntoReturned] = Structs.TokenAndAmount({token: ETH_TOKEN, amount: ethBalance});

    //adjust size of returned and failed arrays
    assembly {
      mstore(returned, add(insertIntoReturned, 1)) //+1 for ETH
      mstore(failed, insertIntoFailed)
    }

    // try post withdraw hook, return false (and no trusted payload) if it reverts
    try this.getTrustProcessPayload(appEventData, executeSuccess, returned, failed) returns (
      bytes memory trustedPayload
    ) {
      return (true, returned, failed, trustedPayload);
    } catch {
      return (false, returned, failed, '');
    }
  }

  /// @notice Override to implement custom execute behavior.
  ///         Called inside execute, which is restricted to the ProcessorEndpoint.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  function _execute(Structs.EventData calldata appEventData) internal virtual;

  /// @notice Override to add logic that runs after the default withdraw sweep.
  ///         Called at the end of withdraw, after all tokens and ETH have been transferred.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param executeSuccess True if the execute call succeeded, false if it reverted.
  /// @param returnedTokens Array of tokens actually transferred in the sweep, with amounts.
  /// @param failedTokens Array of tokens that failed to transfer.
  /// @return trustedPayload Optional payload; when non-empty the ProcessorEndpoint enqueues
  ///         a trusted (TRUSTPROCESS) request with it. Return "" to enqueue nothing.
  function getTrustProcessPayload(
    Structs.EventData calldata appEventData,
    bool executeSuccess,
    Structs.TokenAndAmount[] memory returnedTokens,
    Structs.TokenAndAmount[] memory failedTokens
  ) public virtual returns (bytes memory trustedPayload);
}
