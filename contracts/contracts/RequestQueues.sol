// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';

import './Structs.sol';
import './interfaces/ITrigger.sol';

/// @title RequestQueues
/// @notice Pending-request queue storage, round-robin application selection and the
///         reset tooling used by `ProcessorEndpoint`.
/// @dev All functions are `internal`, so this library is inlined into `ProcessorEndpoint` and
///      needs no deployment or linking: it exists to keep the queue mechanics in one reviewable
///      unit, not to move code out of the endpoint. Making the heavy functions `public` (an
///      external, delegatecall-linked library) was measured and rejected — the endpoint has to
///      ABI-encode the many storage-pointer and dynamic-array arguments at each call site, which
///      cost almost as much as the logic it moved out (143 bytes saved overall).
///
///      The storage layout is declared once, by the endpoint, as a single `Store`, and passed in
///      as a pointer.
library RequestQueues {
  using SafeERC20 for IERC20;

  error TransferFailed();

  struct Queue {
    mapping(uint256 => bytes32) idByOrder;
    uint256 head;
    uint256 tail;
  }

  struct Store {
    /// @dev Every pending request, keyed by its unique requestId, regardless of which queue
    ///      holds it. Request ids are unique across queues (they mix in the applicationId and
    ///      the enqueuing queue's tail), so one store can back all of them, and
    ///      `requestById`/`isCurrentPendingRequest` need no applicationId argument.
    mapping(bytes32 => Structs.PendingRequest) requests;
    /// @dev Per-application FIFO queues. A batch is always scoped to one application, so each
    ///      application's requests must be dequeued independently of the others'.
    mapping(uint64 => Queue) pending;
    /// @dev Global FIFO queue for DEPLOYAPP requests. A deploy's application does not exist
    ///      yet — no state root, absent from the deployed-app list — so it cannot be served by
    ///      the round-robin scan over deployed applications.
    Queue deploys;
    /// @dev Global priority queue for TRUSTPROCESS requests, served before everything else.
    Queue triggers;
    /// @dev Round-robin cursor into the deployed application ids: index of the application
    ///      whose turn it is.
    uint256 cursor;
    /// @dev Number of requests across every queue, maintained incrementally so the size views
    ///      never have to scan every application. Includes the trigger queue.
    uint256 total;
  }

  // --- Inlined helpers (hot paths, no delegatecall) ---

  function size(Queue storage q) internal view returns (uint256) {
    if (q.tail > q.head) return q.tail - q.head;
    return 0;
  }

  function peekHead(Queue storage q) internal view returns (bytes32) {
    return q.idByOrder[q.head];
  }

  function isHead(Queue storage q, bytes32 id) internal view returns (bool) {
    return q.tail > q.head && q.idByOrder[q.head] == id;
  }

  function enqueue(
    Store storage s,
    Queue storage q,
    bytes32 id,
    Structs.PendingRequest memory req
  ) internal {
    s.requests[id] = req;
    q.idByOrder[q.tail] = id;
    unchecked {
      ++q.tail;
      ++s.total;
    }
  }

  function dequeueHead(Store storage s, Queue storage q) internal {
    bytes32 id = q.idByOrder[q.head];
    delete s.requests[id];
    delete q.idByOrder[q.head];
    unchecked {
      ++q.head;
      --s.total;
    }
  }

  // --- Library-resident logic (delegatecall) ---

  /// @notice Starting at the round-robin cursor, scans `appIds` (wrapping around) for the first
  ///         application with a non-empty queue. Skipping empty queues keeps the rotation
  ///         work-conserving: a single busy application is served every time when nothing else
  ///         has pending work.
  function selectApplication(
    Store storage s,
    uint64[] storage appIds
  ) internal view returns (uint64 applicationId, bool found) {
    uint256 n = appIds.length;
    if (n == 0) return (0, false);
    uint256 start = s.cursor % n;
    uint256 i;
    while (i != n) {
      uint64 candidate = appIds[(start + i) % n];
      if (size(s.pending[candidate]) > 0) return (candidate, true);
      unchecked {
        ++i;
      }
    }
    return (0, false);
  }

  /// @notice Moves the cursor just past the given application, so the next selection starts
  ///         from the following one. No-op if the application is not in `appIds` (a request of
  ///         a pending deploy's application, which lives in the deploy queue, never gets here).
  function advanceCursor(Store storage s, uint64[] storage appIds, uint64 applicationId) internal {
    uint256 n = appIds.length;
    uint256 i;
    while (i != n) {
      if (appIds[i] == applicationId) {
        unchecked {
          ++i;
        }
        s.cursor = i == n ? 0 : i;
        return;
      }
      unchecked {
        ++i;
      }
    }
  }

  /// @notice Discards every pending request: pending deploys give their reserved slot back and
  ///         drop any trigger registered eagerly at submit time, per-application requests have
  ///         their asset deposit credited back to the sender, and the trigger queue is dropped
  ///         (TRUSTPROCESS requests carry no funds).
  /// @return freedDeploySlots Deploy slots to return to the pool, applied by the caller because
  ///         a value-type state variable cannot be passed by reference.
  function resetQueues(
    Store storage s,
    uint64[] storage appIds,
    mapping(uint64 => mapping(address => uint256)) storage appCustody,
    mapping(address => uint256) storage totalAppCustody,
    mapping(address => mapping(address => uint256)) storage pendingClaims,
    mapping(address => uint256) storage totalPendingClaims,
    mapping(uint64 => ITrigger) storage triggerContracts,
    mapping(address => uint64) storage triggersToAppIds
  ) internal returns (uint256 freedDeploySlots) {
    freedDeploySlots = size(s.deploys);
    uint256 i = s.deploys.head;
    uint256 tail = s.deploys.tail;
    while (i != tail) {
      _clearTrigger(
        triggerContracts,
        triggersToAppIds,
        s.requests[s.deploys.idByOrder[i]].applicationId
      );
      unchecked {
        ++i;
      }
    }
    drain(s, s.deploys);

    // Every application with a pending request is in `appIds`: a regular request is only
    // accepted for an application with a non-zero state root, which is set when its deploy
    // succeeds — the same point at which the application is appended to the list.
    uint256 appCount = appIds.length;
    uint256 a;
    while (a != appCount) {
      Queue storage q = s.pending[appIds[a]];
      i = q.head;
      tail = q.tail;
      while (i != tail) {
        Structs.PendingRequest storage req = s.requests[q.idByOrder[i]];
        uint256 amount = req.assetAmount;
        if (amount > 0) {
          address token = req.tokenAddress;
          uint64 appId = req.applicationId;
          appCustody[appId][token] -= amount;
          totalAppCustody[token] -= amount;
          pendingClaims[token][req.sender] += amount;
          totalPendingClaims[token] += amount;
        }
        unchecked {
          ++i;
        }
      }
      drain(s, q);
      unchecked {
        ++a;
      }
    }

    drain(s, s.triggers);
  }

  /// @notice Clears every entry of a queue and collapses it by setting tail back to head. head
  ///         is intentionally left at its current value so that future queue indices continue
  ///         from where they left off rather than restarting from zero (avoids any risk of
  ///         re-using a slot index still in storage).
  function drain(Store storage s, Queue storage q) internal {
    uint256 i = q.head;
    uint256 tail = q.tail;
    unchecked {
      s.total -= tail - i;
    }
    while (i != tail) {
      delete s.requests[q.idByOrder[i]];
      delete q.idByOrder[i];
      unchecked {
        ++i;
      }
    }
    q.tail = q.head;
  }

  /// @notice Sweeps the custody of the given applications to `recipient`, clears their state
  ///         roots, removes them from `appIds` and drops any trigger they registered.
  /// @dev Runs under `delegatecall`, so the transfers move the endpoint's own ETH and tokens.
  /// @return freedSlots Deploy slots to return to the pool, applied by the caller.
  function resetApps(
    uint64[] storage appIds,
    uint64[] memory effectiveAppIds,
    address[] memory tokens,
    address payable recipient,
    mapping(uint64 => bytes32) storage applicationStateRoots,
    mapping(uint64 => mapping(address => uint256)) storage appCustody,
    mapping(address => uint256) storage totalAppCustody,
    mapping(uint64 => ITrigger) storage triggerContracts,
    mapping(address => uint64) storage triggersToAppIds
  ) internal returns (uint256 freedSlots) {
    uint256 appCount = effectiveAppIds.length;
    uint256 tokenCount = tokens.length;
    uint256 i;
    while (i != appCount) {
      uint64 appId = effectiveAppIds[i];

      // Transfer ETH custody for this app and clear both the per-app and global trackers.
      uint256 ethAmt = appCustody[appId][ETH_TOKEN];
      if (ethAmt > 0) {
        (bool ok, ) = recipient.call{value: ethAmt}('');
        if (!ok) revert TransferFailed();
        appCustody[appId][ETH_TOKEN] -= ethAmt;
        totalAppCustody[ETH_TOKEN] -= ethAmt;
      }

      // Transfer ERC-20 custody for this app across every token in the effective list.
      uint256 j;
      while (j != tokenCount) {
        uint256 amt = appCustody[appId][tokens[j]];
        if (amt > 0) {
          IERC20(tokens[j]).safeTransfer(recipient, amt);
          appCustody[appId][tokens[j]] -= amt;
          totalAppCustody[tokens[j]] -= amt;
        }
        unchecked {
          ++j;
        }
      }

      // Clear the app's state root, remove it from the deployed list, and return its deploy
      // slot to the pool so the same slot capacity can be reused after the reset.
      if (applicationStateRoots[appId] != bytes32(0)) {
        applicationStateRoots[appId] = bytes32(0);
        _removeAppId(appIds, appId);
        unchecked {
          ++freedSlots;
        }
      }

      // Clear any trigger registered for this app so its address can be reused after reset.
      _clearTrigger(triggerContracts, triggersToAppIds, appId);

      unchecked {
        ++i;
      }
    }
  }

  function _removeAppId(uint64[] storage appIds, uint64 appId) private {
    uint256 len = appIds.length;
    uint256 i;
    while (i != len) {
      if (appIds[i] == appId) {
        appIds[i] = appIds[len - 1];
        appIds.pop();
        break;
      }
      unchecked {
        ++i;
      }
    }
  }

  /// @dev Removes any trigger registration for the given application (both the forward and
  ///      reverse mappings). No-op when no trigger is registered.
  function _clearTrigger(
    mapping(uint64 => ITrigger) storage triggerContracts,
    mapping(address => uint64) storage triggersToAppIds,
    uint64 applicationId
  ) private {
    ITrigger trigger = triggerContracts[applicationId];
    if (address(trigger) != address(0)) {
      delete triggersToAppIds[address(trigger)];
      delete triggerContracts[applicationId];
    }
  }
}
