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

  /// @param _processorEndpoint ProcessorEndpoint that will call execute and withdraw. Its
  ///        tokenAllowlist() is read to determine which ERC-20 tokens are swept on withdraw,
  ///        so the allowlist can never be misconfigured independently of the endpoint.
  constructor(IProcessorEndpoint _processorEndpoint) {
    if (address(_processorEndpoint) == address(0)) {
      revert ZeroAddress();
    }
    processorEndpoint = _processorEndpoint;
    tokenAllowlist = _processorEndpoint.tokenAllowlist();
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
  ///         Performs ONLY the sweep and cannot be overridden; producing the optional
  ///         trusted follow-up payload is a separate, explicit step (getTrustProcessPayload).
  /// @return returnedTokens Array of tokens transferred in the sweep, with amounts (ETH last).
  /// @return failedTokens Array of tokens that failed to transfer, with amounts.
  function withdraw()
    public
    _onlyProcessorEndpoint
    returns (Structs.TokenAndAmount[] memory, Structs.TokenAndAmount[] memory)
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

    return (returned, failed);
  }

  /// @notice Produces the optional trusted follow-up payload. Restricted to the ProcessorEndpoint,
  ///         which calls it explicitly after withdraw inside its own try/catch. Delegates to the
  ///         overridable _getTrustProcessPayload hook (mirrors the execute/_execute pattern).
  /// @inheritdoc ITrigger
  function getTrustProcessPayload(
    Structs.EventData calldata appEventData,
    bool executeSuccess,
    bool withdrawSuccess,
    Structs.TokenAndAmount[] calldata returnedTokens,
    Structs.TokenAndAmount[] calldata failedTokens
  ) external _onlyProcessorEndpoint returns (bytes memory) {
    return
      _getTrustProcessPayload(
        appEventData,
        executeSuccess,
        withdrawSuccess,
        returnedTokens,
        failedTokens
      );
  }

  /// @notice Override to implement custom execute behavior.
  ///         Called inside execute, which is restricted to the ProcessorEndpoint.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  function _execute(Structs.EventData calldata appEventData) internal virtual;

  /// @notice Override to produce the optional trusted (TRUSTPROCESS) follow-up payload.
  ///         Called by the ProcessorEndpoint after execute and withdraw, in an isolated
  ///         try/catch, even when withdraw failed — so the application can react to a failed
  ///         sweep (damage control) by returning an appropriate payload.
  /// @param appEventData Application-level (non-encrypted) events and their subtypes.
  /// @param executeSuccess True if the execute call succeeded, false if it reverted.
  /// @param withdrawSuccess True if the withdraw sweep succeeded, false if it reverted.
  /// @param returnedTokens Array of tokens actually transferred in the sweep, with amounts.
  /// @param failedTokens Array of tokens that failed to transfer.
  /// @return trustedPayload Optional payload; when non-empty the ProcessorEndpoint enqueues
  ///         a trusted (TRUSTPROCESS) request with it. Return "" to enqueue nothing.
  function _getTrustProcessPayload(
    Structs.EventData calldata appEventData,
    bool executeSuccess,
    bool withdrawSuccess,
    Structs.TokenAndAmount[] calldata returnedTokens,
    Structs.TokenAndAmount[] calldata failedTokens
  ) internal virtual returns (bytes memory trustedPayload);
}
