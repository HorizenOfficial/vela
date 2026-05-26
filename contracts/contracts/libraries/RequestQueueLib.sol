// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '../Structs.sol';

/// @title RequestQueueLib
/// @notice External library implementing a FIFO queue and a sorted priority queue for pending
///         requests. Deployed independently and linked into ProcessorEndpoint.
library RequestQueueLib {
  // ── FIFO queue ────────────────────────────────────────────────────────────

  struct NormalQueue {
    mapping(bytes32 => Structs.PendingRequest) requestById;
    mapping(uint256 => bytes32) idByOrder;
    uint256 head;
    uint256 tail;
  }

  function enqueue(NormalQueue storage q, bytes32 id, Structs.PendingRequest memory req) public {
    q.requestById[id] = req;
    q.idByOrder[q.tail] = id;
    unchecked {
      ++q.tail;
    }
  }

  function dequeueHead(NormalQueue storage q) public {
    bytes32 id = q.idByOrder[q.head];
    delete q.requestById[id];
    delete q.idByOrder[q.head];
    unchecked {
      ++q.head;
    }
  }

  function size(NormalQueue storage q) public view returns (uint256) {
    if (q.tail > q.head) return q.tail - q.head;
    return 0;
  }

  function peekHead(NormalQueue storage q) public view returns (bytes32) {
    return q.idByOrder[q.head];
  }

  function isHead(NormalQueue storage q, bytes32 id) public view returns (bool) {
    return q.tail > q.head && q.idByOrder[q.head] == id;
  }

  function getAll(
    NormalQueue storage q
  ) public view returns (Structs.PendingRequest[] memory result) {
    uint256 n = size(q);
    result = new Structs.PendingRequest[](n);
    uint256 i = q.head;
    uint256 tail = q.tail;
    uint256 j;
    while (i < tail) {
      result[j] = q.requestById[q.idByOrder[i]];
      unchecked {
        ++i;
        ++j;
      }
    }
  }

  function getPage(
    NormalQueue storage q,
    uint256 offset,
    uint256 limit
  ) public view returns (Structs.PendingRequest[] memory result) {
    uint256 n = size(q);
    if (offset >= n || limit == 0) return new Structs.PendingRequest[](0);
    uint256 end = offset + limit;
    if (end > n) end = n;
    uint256 count = end - offset;
    result = new Structs.PendingRequest[](count);
    uint256 i = q.head + offset;
    uint256 stop = q.head + end;
    uint256 j;
    while (i < stop) {
      result[j] = q.requestById[q.idByOrder[i]];
      unchecked {
        ++i;
        ++j;
      }
    }
  }
}
