// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

//they are isolated so they can be imported in multiple contracts without loops
/// @dev Sentinel address representing native ETH as the token.
address constant ETH_TOKEN = address(0);

contract Structs {
  //type definitions
  enum RequestType {
    DEPLOYAPP,
    PROCESS,
    DEANONYMIZATION,
    ASSOCIATEKEY,
    TRUSTPROCESS
  }
  enum RequestResult {
    COMPLETED,
    FAILED
  }
  enum ErrorCode {
    NO_ERROR,
    UNKNOWN,
    INTERNAL,
    APPLICATION_ALREADY_DEPLOYED,
    FUNCTION_NOT_FOUND,
    DEPOSIT_FAILED,
    REQUEST_FUNC_FAILED,
    WRONG_KEY_SENT,
    PUB_KEY_NOT_REGISTERED,
    NO_REPORT_DATA_FOUND,
    WASM_INTERNAL,
    INSUFFICIENT_FUEL
  }

  struct PendingRequest {
    uint256 timestamp; //assigned automatically
    address tokenAddress; // 0x0 = ETH
    uint256 assetAmount;
    uint256 maxFeeValue;
    bytes32 requestId; //assigned automatically
    bytes payload;
    address sender; //assigned automatically
    address facilitator; // address(0) for direct submissions, facilitator address for meta-tx
    uint64 applicationId;
    uint8 protocolVersion;
    RequestType requestType;
  }

  struct WithdrawalRequest {
    address tokenAddress; // 0x0 = ETH
    address payable receiver;
    uint256 amount;
  }
  struct EventData {
    bytes[] events;
    bytes32[] subTypes;
  }

  // Parameters for TEE signature verification
  struct SignatureParams {
    uint64 applicationId;
    bytes32 prevStateRoot;
    bytes32 newStateRoot;
    bytes32 processedRequestId;
    EventData userEvents;
    EventData appEvents;
    WithdrawalRequest[] withdrawalRequests;
    uint256 refundAmount;
    uint256 applicationFee;
    Structs.ErrorCode errorCode;
    string errorMsg;
  }

  /// @notice One request's update payload inside a `batchStateUpdate` call.
  /// @dev Mirrors the arguments of `stateUpdate` minus two of them:
  ///      - `applicationId` is deduplicated to the batch — every entry of a batch belongs to the
  ///        same application, so it is passed once alongside the entry array;
  ///      - `signature` is replaced by the single batch signature covering every entry hash.
  ///      Field names and order otherwise follow `SignatureParams`, which is what the per-entry
  ///      hash is built from (`UpdateEntryHash.entryHash`): both event sets are carried, because
  ///      the entry hash cannot be reconstructed without them.
  struct BatchEntry {
    bytes32 prevStateRoot;
    bytes32 newStateRoot;
    bytes32 processedRequestId;
    EventData userEvents;
    EventData appEvents;
    WithdrawalRequest[] withdrawalRequests;
    uint256 refund;
    uint256 applicationFees;
    ErrorCode errorCode;
    string errorMsg;
  }

  /// @notice Describes a token and amount pair.
  struct TokenAndAmount {
    /// @notice Token address; address(0) for ETH.
    address token;
    /// @notice Amount of the token.
    uint256 amount;
  }
}
