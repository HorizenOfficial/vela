# ProcessorEndpoint Admin Reset (Testnet / Development Only)

## Summary

During development and on testnets it may be necessary to restart the system from scratch: clear the pending request queue and reset all application state roots. This is done through two dedicated functions on `ProcessorEndpoint`, protected by a `RESET_OPERATOR` role.

Passing `address(0)` as the `RESET_OPERATOR` at deploy time disables the feature permanently for that deployment, which is the expected value in production.

> **Note:** `adminResetApps` accepts an optional list of ERC-20 token addresses. Before zeroing custody, it transfers all accumulated ETH and ERC-20 balances locked in the supplied app IDs directly to the caller (`msg.sender`). This prevents permanent fund loss during a reset.

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

This guarantee rests on its admin staying `DEFAULT_ADMIN_ROLE` **and that role being granted to nobody**. `ADMIN`, `UPDATE_STATUS_ROLE` and `DEPLOYER_ROLE` are each given an explicit admin in `initialize` precisely so they can be rotated without anyone holding `DEFAULT_ADMIN_ROLE`; `RESET_OPERATOR` is the one role deliberately left out. Granting `DEFAULT_ADMIN_ROLE` to any address — or calling `_setRoleAdmin(RESET_OPERATOR, ...)` — silently voids the permanent-disable property: the holder could grant itself `RESET_OPERATOR` post-deployment and sweep every app's custody to itself via `adminResetApps`.

### Deployed App ID Tracking

`ProcessorEndpoint` maintains an append-only array of all application IDs that have ever been successfully deployed:

```solidity
uint64[] private _deployedAppIds;
```

When a `DEPLOYAPP` request is finalised (i.e. `applicationStateRoots[appId]` is set for the first time), the `appId` is appended to `_deployedAppIds`. A public view function exposes the full list:

```solidity
function getDeployedAppIds() external view returns (uint64[] memory) {
    return _deployedAppIds;
}
```

This allows the reset script (and any off-chain tooling) to discover all deployed app IDs without requiring the caller to maintain that list externally.

### `adminReset()`

Clears the pending request queue and resets deploy slot accounting.

Operations performed:

1. Deletes every pending request from `requestById` and `_requestIdByOrder` by iterating from `_head` to `_tail`.
2. Sets `_tail = _head` (the queue is now empty; the head index is not re-anchored to zero).
3. For each `DEPLOYAPP` request removed from the queue, increments `availableDeploySlots` by one. Already-deployed applications are not affected — their deploy slots remain consumed.

### `adminResetApps(uint64[] calldata appIds, address[] calldata erc20Tokens)`

Resets state roots and locked funds for the supplied application IDs, **and also calls `_resetQueue()` internally**, so the queue and deploy slots are cleared in the same transaction. Before zeroing any custody, it recovers all locked ETH and ERC-20 balances by transferring them to `msg.sender`.

**Both parameters are optional:**

- If `appIds` is empty (`[]`), the function uses `_deployedAppIds` — i.e. it resets every app that has ever been deployed.
- If `erc20Tokens` is empty (`[]`), the function uses all tokens currently registered in the `TokenAllowlist` — i.e. it rescues custody for every allowed ERC-20 token.

Operations performed:

1. Calls `_resetQueue()` (clears the pending request queue; `availableDeploySlots` is incremented only for pending `DEPLOYAPP` requests removed).
2. Resolves the effective app list: `appIds` if non-empty, otherwise `_deployedAppIds`.
3. Resolves the effective token list: `erc20Tokens` if non-empty, otherwise all addresses returned by `getAllowedTokens()` on the `TokenAllowlist`.
4. Iterates over the effective app list. For each `appId`:
   - Accumulates ETH custody (`appCustody[appId][address(0)]`) into `totalEth`, decrements `totalAppCustody[address(0)]`, and zeroes `appCustody[appId][address(0)]`.
   - For each token in the effective token list, accumulates its custody (`appCustody[appId][token]`) into a per-token running total, decrements `totalAppCustody[token]`, and zeroes `appCustody[appId][token]`.
   - If `applicationStateRoots[appId]` is non-zero, clears it and increments `availableDeploySlots` by one (the deploy slot is freed).
5. Transfers the accumulated `totalEth` to `msg.sender` in a single `call`.
6. For each token in the effective token list, transfers the accumulated total to `msg.sender` via `safeTransfer`.

State changes are completed before any external transfer (checks-effects-interactions). ETH and ERC-20 totals are accumulated across all apps so that each asset requires exactly one transfer call regardless of how many app IDs are supplied.

> **Data layer note:** After calling `adminResetApps`, the corresponding off-chain data layer entries for each cleared app ID must also be removed. If the same app ID is reused in a subsequent deploy and the old data layer entries are still present, the new deployment will fail because it will find pre-existing state that does not match the freshly cleared on-chain root.

> **Subgraph note:** The subgraph indexes on-chain events and has no automatic awareness of an on-chain reset. After a full reset, the subgraph should also be reset or redeployed (starting from block 0 or the new deployment block); otherwise off-chain clients may surface events and application state that no longer exist on-chain.

### What Is Not Reset

Both functions intentionally leave the following fields untouched:

- **`pendingClaims` and `totalPendingClaims`** — These represent funds already owed to users (pull-payment balances for withdrawals, refunds, and fees). If a full balance wipe is also needed (e.g. for a completely fresh testnet with no real value at stake), it must be done separately and only after all pending claims have been withdrawn.
- **`maxQueueSize` and `maxNumOfApplications`** — These are configuration parameters set at deployment time and are not considered part of the mutable state that a reset targets. They remain unchanged so that the system continues to operate under the same capacity limits after the reset.

---

## Example Implementation

```solidity
bytes32 public constant RESET_OPERATOR = keccak256('RESET_OPERATOR');

uint64[] private _deployedAppIds;

// inside constructor(..., address resetOperator)
// (or initialize(...) if using the proxy upgrade pattern):
if (resetOperator != address(0)) {
    _grantRole(RESET_OPERATOR, resetOperator);
}

/// @notice Returns all application IDs that have ever been deployed.
function getDeployedAppIds() external view returns (uint64[] memory) {
    return _deployedAppIds;
}

/// @notice Resets the queue and deploy slots.
/// @dev    Does not clear pending claim balances.
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
///         Pass empty arrays to target all deployed apps / all allowlisted tokens.
///         Unreachable when RESET_OPERATOR was initialised as address(0).
function adminResetApps(
    uint64[] calldata appIds,
    address[] calldata erc20Tokens
) external onlyRole(RESET_OPERATOR) nonReentrant {
    _resetQueue();

    // Resolve effective app and token lists
    uint64[] memory effectiveAppIds  = appIds.length > 0 ? appIds : _deployedAppIds;
    address[] memory effectiveTokens = erc20Tokens.length > 0
        ? erc20Tokens
        : ITokenAllowlist(address(this)).getAllowedTokens();

    uint256 appCount   = effectiveAppIds.length;
    uint256 tokenCount = effectiveTokens.length;
    uint256 totalEth;
    uint256[] memory tokenTotals = new uint256[](tokenCount);

    // --- Effects: zero out all custody and state roots ---
    uint256 i;
    while (i < appCount) {
        uint64 appId = effectiveAppIds[i];

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
            uint256 amt = appCustody[appId][effectiveTokens[j]];
            if (amt > 0) {
                tokenTotals[j] += amt;
                totalAppCustody[effectiveTokens[j]] -= amt;
                appCustody[appId][effectiveTokens[j]] = 0;
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
            IERC20(effectiveTokens[j]).safeTransfer(recipient, tokenTotals[j]);
        }
        unchecked { ++j; }
    }
}
```

> **TokenAllowlist note:** `getAllowedTokens()` must be added to `ITokenAllowlist` and implemented in `TokenAllowlist`. It returns the list of all addresses for which `allowedTokens[token] == true`. Because the allowlist is expected to be small (single-digit to low double-digit entries), on-chain enumeration is acceptable here.

---

## Impact on Constructor (or `initialize`) Signature

The constructor of `ProcessorEndpoint` gains one parameter (or `initialize(...)` if the proxy upgrade pattern is in use):

```
resetOperator  address   // address(0) to disable reset in production
```

The deploy scripts must be updated accordingly. In production deployments the value must always be `address(0)`.

---

## Reset Script (`scripts/management/adminReset.ts`)

A new script that calls `adminReset()` and optionally `adminResetApps(appIds, erc20Tokens)` on a deployed `ProcessorEndpoint`. **Intended for testnet and development use only.**

Required environment variables:

| Variable | Description |
|----------|-------------|
| `PROCESSOR_ENDPOINT` | Address of the deployed contract (proxy address if using the upgrade pattern). |
| `APP_IDS` | *(Optional)* Comma-separated list of `uint64` application IDs whose state roots and locked funds should be cleared. If omitted or empty, all app IDs returned by `getDeployedAppIds()` are used. |
| `ERC20_TOKENS` | *(Optional)* Comma-separated list of ERC-20 token addresses whose per-app custody should be rescued and transferred to the signer. If omitted or empty, all tokens returned by `getAllowedTokens()` on the contract are used. |

The script flow:

1. Attach to `ProcessorEndpoint` at `PROCESSOR_ENDPOINT`.
2. If `APP_IDS` is not set, call `getDeployedAppIds()` on the contract to discover all deployed app IDs automatically.
3. If `ERC20_TOKENS` is not set, call `getAllowedTokens()` on the contract to discover all allowlisted token addresses automatically.
4. Call `adminResetApps(appIds, erc20Tokens)`, which internally calls `_resetQueue()` and resets all resolved apps and tokens in a single transaction. All locked ETH and ERC-20 balances are transferred to the signer's address.
5. Print confirmation of each step, the resulting queue size (expected: 0), and the amounts recovered per asset.
