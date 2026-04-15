# ProcessorEndpoint Admin Reset (Testnet / Development Only)

## Summary

During development and on testnets it may be necessary to restart the system from scratch: clear the pending request queue and reset all application state roots. This is done through two dedicated functions on `ProcessorEndpoint`, protected by a `RESET_OPERATOR` role.

Passing `address(0)` as the `RESET_OPERATOR` at deploy time disables the feature permanently for that deployment, which is the expected value in production.

> **Note:** `adminResetApps` accepts an explicit list of ERC-20 token addresses. Before zeroing custody, it transfers all accumulated ETH and ERC-20 balances locked in the supplied app IDs directly to the caller (`msg.sender`). This prevents permanent fund loss during a reset.

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
2. Sets `_tail = _head` (the queue is now empty; the head index is not re-anchored to zero).
3. For each `DEPLOYAPP` request removed from the queue, increments `availableDeploySlots` by one. Already-deployed applications are not affected — their deploy slots remain consumed.

### `adminResetApps(uint64[] calldata appIds, address[] calldata erc20Tokens)`

Resets state roots and locked funds for the supplied application IDs, **and also calls `adminReset()` internally**, so the queue and deploy slots are cleared in the same transaction. Before zeroing any custody, it recovers all locked ETH and ERC-20 balances by transferring them to `msg.sender`. The caller must provide the list of active app IDs and the list of ERC-20 token addresses to rescue; both are explicit to avoid unbounded on-chain iteration.

Operations performed:

1. Calls `_resetQueue()` (clears the pending request queue; `availableDeploySlots` is incremented only for pending `DEPLOYAPP` requests removed).
2. Iterates over `appIds`. For each `appId`:
   - Accumulates ETH custody (`appCustody[appId][address(0)]`) into `totalEth`, decrements `totalAppCustody[address(0)]`, and zeroes `appCustody[appId][address(0)]`.
   - For each token in `erc20Tokens`, accumulates its custody (`appCustody[appId][token]`) into a per-token running total, decrements `totalAppCustody[token]`, and zeroes `appCustody[appId][token]`.
   - If `applicationStateRoots[appId]` is non-zero, clears it and increments `availableDeploySlots` by one (the deploy slot is freed).
3. Transfers the accumulated `totalEth` to `msg.sender` in a single `call`.
4. For each token in `erc20Tokens`, transfers the accumulated total to `msg.sender` via `safeTransfer`.

State changes are completed before any external transfer (checks-effects-interactions). ETH and ERC-20 totals are accumulated across all apps so that each asset requires exactly one transfer call regardless of how many app IDs are supplied.

> **Data layer note:** After calling `adminResetApps`, the corresponding off-chain data layer entries for each cleared app ID must also be removed. If the same app ID is reused in a subsequent deploy and the old data layer entries are still present, the new deployment will fail because it will find pre-existing state that does not match the freshly cleared on-chain root.

> **Subgraph note:** The subgraph indexes on-chain events and has no automatic awareness of an on-chain reset. After a full reset, the subgraph should also be reset or redeployed (starting from block 0 or the new deployment block); otherwise off-chain clients may surface events and application state that no longer exist on-chain.

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
    uint256 freedDeploySlots;
    while (i < tail) {
        if (requestById[_requestIdByOrder[i]].requestType == Structs.RequestType.DEPLOYAPP) {
            unchecked { ++freedDeploySlots; }
        }
        delete requestById[_requestIdByOrder[i]];
        delete _requestIdByOrder[i];
        unchecked { ++i; }
    }
    _tail = _head; // empty the queue without re-anchoring _head to zero
    unchecked { availableDeploySlots += freedDeploySlots; }
}

/// @notice Resets state roots and locked funds for the given application IDs,
///         clears the queue in the same transaction, and rescues all locked ETH
///         and ERC-20 balances by transferring them to the caller.
///         Unreachable when RESET_OPERATOR was initialised as address(0).
function adminResetApps(
    uint64[] calldata appIds,
    address[] calldata erc20Tokens
) external onlyRole(RESET_OPERATOR) nonReentrant {
    _resetQueue();

    uint256 appCount  = appIds.length;
    uint256 tokenCount = erc20Tokens.length;
    uint256 totalEth;
    uint256[] memory tokenTotals = new uint256[](tokenCount);

    // --- Effects: zero out all custody and state roots ---
    uint256 i;
    while (i < appCount) {
        uint64 appId = appIds[i];

        // Accumulate ETH custody
        uint256 ethAmt = appCustody[appId][address(0)];
        if (ethAmt > 0) {
            totalEth += ethAmt;
            totalAppCustody[address(0)] -= ethAmt;
            appCustody[appId][address(0)] = 0;
        }

        // Accumulate ERC-20 custody per token
        uint256 j;
        while (j < tokenCount) {
            uint256 amt = appCustody[appId][erc20Tokens[j]];
            if (amt > 0) {
                tokenTotals[j] += amt;
                totalAppCustody[erc20Tokens[j]] -= amt;
                appCustody[appId][erc20Tokens[j]] = 0;
            }
            unchecked { ++j; }
        }

        // Clear app state root and free deploy slot
        if (applicationStateRoots[appId] != bytes32(0)) {
            applicationStateRoots[appId] = bytes32(0);
            unchecked { ++availableDeploySlots; }
        }

        unchecked { ++i; }
    }

    // --- Interactions: transfer accumulated balances to caller ---
    address payable recipient = payable(msg.sender);

    if (totalEth > 0) {
        (bool ok, ) = recipient.call{value: totalEth}('');
        if (!ok) revert TransferFailed();
    }

    uint256 j;
    while (j < tokenCount) {
        if (tokenTotals[j] > 0) {
            IERC20(erc20Tokens[j]).safeTransfer(recipient, tokenTotals[j]);
        }
        unchecked { ++j; }
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
| `ERC20_TOKENS` | Comma-separated list of ERC-20 token addresses whose per-app custody should be rescued and transferred to the signer. May be empty if the deployment is ETH-only. |

The script flow:

1. Attach to `ProcessorEndpoint` at `PROCESSOR_ENDPOINT`.
2. Call `adminReset()` to drain the queue and reclaim deploy slots for pending deploy requests.
3. If `APP_IDS` is set and non-empty, parse it and call `adminResetApps(appIds, erc20Tokens)` (passing `erc20Tokens` from `ERC20_TOKENS`, defaulting to `[]` if unset). All locked ETH and ERC-20 balances are transferred to the signer's address.
4. Print confirmation of each step, the resulting queue size (expected: 0), and the amounts recovered per asset.