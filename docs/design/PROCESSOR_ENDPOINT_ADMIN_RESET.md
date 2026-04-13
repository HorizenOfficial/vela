# ProcessorEndpoint Admin Reset (Testnet / Development Only)

## Summary

During development and on testnets it may be necessary to restart the system from scratch: clear the pending request queue and reset all application state roots. This is done through two dedicated functions on `ProcessorEndpoint`, protected by a `RESET_OPERATOR` role.

Passing `address(0)` as the `RESET_OPERATOR` at deploy time disables the feature permanently for that deployment, which is the expected value in production.

---

## Design

### `RESET_OPERATOR` Role

A new `RESET_OPERATOR` role constant is added to `ProcessorEndpoint`:

```solidity
bytes32 public constant RESET_OPERATOR = keccak256('RESET_OPERATOR');
```

The role is granted inside the constructor (or in `initialize(...)` if the proxy upgrade pattern is in use), which receives a new `resetOperator` parameter. If `address(0)` is passed, `_grantRole` is not called and the role remains unassigned, making both reset functions permanently unreachable.

```solidity
// inside constructor (or initialize(...))
if (resetOperator != address(0)) {
    _grantRole(RESET_OPERATOR, resetOperator);
}
```

The role is intentionally not administered by `ADMIN` (i.e. `_setRoleAdmin` is not called for it) so that even an `ADMIN` cannot grant it after deployment. It can only be assigned at initialisation time.

### `adminReset()`

Clears the pending request queue and resets deploy slot accounting.

Operations performed:

1. Deletes every pending request from `requestById` and `_requestIdByOrder` by iterating from `_head` to `_tail`.
2. Resets `_head` and `_tail` to `0`.
3. Resets `availableDeploySlots` to `maxNumOfApplications`.

### `adminResetApps(uint64[] calldata appIds)`

Resets state roots and locked funds for the supplied application IDs, **and also calls `adminReset()` internally**, so the queue and deploy slots are cleared in the same transaction. The caller must provide the list of active app IDs explicitly to avoid unbounded on-chain iteration.

Operations performed:

1. Calls `adminReset()` (clears the pending request queue and resets `availableDeploySlots`).
2. For each `appId` in the list:
   - Sets `applicationStateRoots[appId]` to `bytes32(0)`.
   - Sets `appLockedFunds[appId]` to `0`.

> **Data layer note:** After calling `adminResetApps`, the corresponding off-chain data layer entries for each cleared app ID must also be removed. If the same app ID is reused in a subsequent deploy and the old data layer entries are still present, the new deployment will fail because it will find pre-existing state that does not match the freshly cleared on-chain root.

### What Is Not Reset

Both functions intentionally leave the following fields untouched:

- **`payments` and `_totalDeposits`** — These represent funds already owed to users. If a full balance wipe is also needed (e.g. for a completely fresh testnet with no real value at stake), it must be done separately and only after all pending payments have been withdrawn.
- **`maxQueueSize` and `maxNumOfApplications`** — These are configuration parameters set at deployment time and are not considered part of the mutable state that a reset targets. They remain unchanged so that the system continues to operate under the same capacity limits after the reset.

---

## Example Implementation

```solidity
bytes32 public constant RESET_OPERATOR = keccak256('RESET_OPERATOR');

// inside constructor(..., address resetOperator)
// (or initialize(...) if using the proxy upgrade pattern):
if (resetOperator != address(0)) {
    _grantRole(RESET_OPERATOR, resetOperator);
}

/// @notice Resets the queue and deploy slots.
/// @dev    Does not clear pending payment balances.
///         Unreachable when RESET_OPERATOR was initialised as address(0).
function adminReset() public onlyRole(RESET_OPERATOR) {
    _resetQueue();
}

function _resetQueue() internal {
    uint256 i = _head;
    uint256 tail = _tail;
    while (i < tail) {
        delete requestById[_requestIdByOrder[i]];
        delete _requestIdByOrder[i];
        unchecked { ++i; }
    }
    _head = 0;
    _tail = 0;

    availableDeploySlots = maxNumOfApplications;
}

/// @notice Resets state roots and locked funds for the given application IDs,
///         and also clears the queue in the same transaction.
///         Unreachable when RESET_OPERATOR was initialised as address(0).
function adminResetApps(uint64[] calldata appIds) external onlyRole(RESET_OPERATOR) {
    _resetQueue();
    uint256 length = appIds.length;
    uint256 i;
    while (i < length) {
        applicationStateRoots[appIds[i]] = bytes32(0);
        appLockedFunds[appIds[i]] = 0;
        unchecked { ++i; }
    }
}
```

---

## Impact on Constructor (or `initialize`) Signature

The constructor of `ProcessorEndpoint` gains one parameter (or `initialize(...)` if the proxy upgrade pattern is in use):

```
resetOperator  address   // address(0) to disable reset in production
```

The deploy scripts must be updated accordingly. In production deployments the value must always be `address(0)`.

---

## Reset Script (`scripts/management/adminReset.ts`)

A new script that calls `adminReset()` and optionally `adminResetApps(appIds)` on a deployed `ProcessorEndpoint`. **Intended for testnet and development use only.**

Required environment variables:

| Variable | Description |
|----------|-------------|
| `PROCESSOR_ENDPOINT` | Address of the deployed contract (proxy address if using the upgrade pattern). |
| `APP_IDS` | Comma-separated list of `uint64` application IDs whose state roots and locked funds should be cleared. |

The script flow:

1. Attach to `ProcessorEndpoint` at `PROCESSOR_ENDPOINT`.
2. Call `adminReset()` to drain the queue and reset `availableDeploySlots`.
3. If `APP_IDS` is set and non-empty, parse it and call `adminResetApps(appIds)`.
4. Print confirmation of each step and the resulting queue size (expected: 0).