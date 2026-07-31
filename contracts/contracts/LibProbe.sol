// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

library LibProbe {
  struct Q {
    mapping(uint256 => bytes32) idByOrder;
    uint256 head;
    uint256 tail;
  }

  // public => linked library, called via delegatecall
  function drain(
    Q storage q,
    mapping(bytes32 => uint256) storage store
  ) public returns (uint256 n) {
    uint256 i = q.head;
    while (i != q.tail) {
      delete store[q.idByOrder[i]];
      delete q.idByOrder[i];
      unchecked {
        ++i;
        ++n;
      }
    }
    q.tail = q.head;
  }
}

contract LibProbeUser {
  mapping(uint64 => LibProbe.Q) private _queues;
  mapping(bytes32 => uint256) private _store;

  function go(uint64 id) external returns (uint256) {
    return LibProbe.drain(_queues[id], _store);
  }
}
