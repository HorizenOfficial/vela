# ERC-20 Support for Deposits / Withdrawals / Claims

## Summary

This document describes the design to support private-state deposits, withdrawals, and claims for selected ERC-20 tokens in addition to native ETH.

Fees remain ETH-only and are explicitly out of scope for ERC-20 support in this design. The system will support apps that are ETH-only as well as apps that support ETH plus an allowlist of ERC-20 tokens.

In this document, `0x0` is used as shorthand for the zero address. In Solidity this means `address(0)`. In Go this means the zero-value Ethereum address.

## Current State

Today the system is ETH-only across the full flow:

- Requests carry a single `depositAmount` value in ETH.
- The WASM `deposit` path receives only an amount and has no notion of token identity.
- Private app state typically stores a single balance per account.
- Withdrawals are modeled only as `(destinationAddress, amount)`, which implicitly means ETH.
- On-chain claims are ETH-only.

This means the current model cannot represent:

- ERC-20 deposits into private state.
- ERC-20 withdrawals from private state.
- ERC-20 pending claims after successful withdrawals or failed requests.


## Goals

- Support one business asset per request, where the asset can be ETH or a selected ERC-20 token.
- Represent the asset using `tokenAddress`, with `0x0` reserved for ETH.
- Allow deposits into private state for ETH and supported ERC-20 tokens.
- Allow withdrawals from private state for ETH and supported ERC-20 tokens.
- Support claim flows for both ETH and ERC-20 pending balances.
- Keep fees denominated and settled in ETH only.

## Non-Goals

- Paying application fees in ERC-20 tokens.
- Supporting more than one business asset in a single request.
- Supporting non-standard ERC-20 behaviors such as fee-on-transfer, rebasing, or other custom semantics.
- Preserving backwards compatibility with the current ETH-only request and withdrawal formats.
- Enforcing per-app token policies


## Requirements

### R1 — Asset Model

A request carries exactly one business asset identified by `(tokenAddress, assetAmount)`. The zero address (`0x0`) represents ETH. Non-zero addresses represent ERC-20 tokens. If `assetAmount` is zero, `tokenAddress` must be `0x0` (a zero-amount deposit has no meaningful token identity). Fees remain denominated and settled in ETH only.

This establishes a uniform representation for all value transfers in the system and avoids the need for separate ETH-only and ERC-20-only code paths.

### R2 — Global Token Allowlist

The contract maintains a global allowlist of ERC-20 tokens that can enter the system. ETH (`0x0`) is always allowed implicitly. Only the `ADMIN` role can add or remove tokens. Adding a token requires verifying that the address contains contract code (see 8). Removing a token blocks new deposits but does not block claims or processing of requests already in the queue.

The global allowlist is the primary defense against all non-standard token behaviors (see ERC-20 Token Behaviors That Affect Vela). Only vetted, standard ERC-20 tokens should be added.

### R3 — No Per-App On-Chain Allowlist

The contract does not enforce per-app token policies. Each WASM application manages its own supported-token configuration at runtime. The contract only enforces the global allowlist.

If a user submits a request with a token the app does not support, the WASM app rejects it during execution. The failed-request path refunds the deposit. This avoids additional on-chain complexity and gives apps full flexibility over their token policies.

### R4 — Deposit Flow

For ETH requests, `msg.value` equals `assetAmount + maxFeeValue`, matching the current behavior. For ERC-20 requests, `msg.value` equals `maxFeeValue` exactly (see 9), and the contract pulls `assetAmount` via `safeTransferFrom` (see 5). The contract verifies the actual received amount matches `assetAmount` using a balance-before/after check (see 2).

### R5 — Custody Accounting

The contract tracks per-app, per-token custody via `appCustody[applicationId][tokenAddress]` and a global aggregate `totalAppCustody[tokenAddress]`. Deposits increase these balances. Withdrawals and failed-request refunds decrease them (moving value into pending claims).

This prevents one application from creating withdrawals backed by another application's deposits in the multi-app model.

### R6 — Withdrawal Model

Withdrawals become `(tokenAddress, receiver, amount)`. A single state update can contain withdrawals in multiple tokens (an app may hold balances in several tokens and produce mixed withdrawals in one execution). The withdrawal tuple ABI must be identical across executor signing, contract signature verification, and Go bindings. Any mismatch in fields or field order causes signature verification to fail.

### R7 — Claim Accounting

Pending claimable balances are tracked per payee and per asset via `pendingClaims[tokenAddress][payee]`, with a global aggregate `totalPendingClaims[tokenAddress]`. The claim entrypoint `claim(tokenAddress, payee)` replaces the current ETH-only `withdrawPayments(payee)`.

Claims are permissionless. If the pending balance is zero, the function returns without reverting (see 10 for why the check must precede any transfer). For ETH, the transfer uses a low-level call. For ERC-20, it uses `safeTransfer`. The claim function does not check the global allowlist, so tokens that have been removed from the allowlist can still be claimed.

### R8 — Failed Request Handling

When a request fails, the business-asset deposit is refunded in its original token. For ETH requests, this produces a single ETH pending credit (asset plus unused fees combined). For ERC-20 requests, this produces two separate pending credits: the business asset in the deposited token and the unused fee portion in ETH.

### R9 — Solvency Checks

Solvency is validated per asset, not as a single aggregated amount. During `stateUpdate`, the contract first decreases `appCustody` and `totalAppCustody` for each outflow (withdrawals and failed-request refunds), then verifies that the contract's physical balance can still cover all remaining obligations before crediting pending claims.

The check, performed after custody is decremented but before pending claims are incremented:

- For ERC-20 tokens: `IERC20(token).balanceOf(contract) - totalAppCustody[token] - totalPendingClaims[token] >= outflowAmount` and `address(this).balance - totalPendingClaims[0x0] - totalAppCustody[0x0] >= maxFeeValue` for fee payment and fee refund.
- For ETH: `address(this).balance - totalAppCustody[0x0] - totalPendingClaims[0x0] >= outflowAmount + maxFeeValue`

where `outflowAmount` is the total business-asset value being moved from custody into pending claims in that update, and `maxFeeValue` accounts for the fee payment and potential fee refund.

App-scoped custody (`appCustody[appId][token] >= outflowAmount`) is checked before the global solvency check to ensure no app creates outflows exceeding its own deposits.

### R10 — Reentrancy Protection

All state-modifying entry points (`submitRequest`, `stateUpdate`, `claim`) use `nonReentrant` guards. The contract follows the checks-effects-interactions pattern: all state updates are performed before any external call. This is the primary in-contract defense against ERC-777 (see 3) and a general best practice for contracts that interact with untrusted external code.

### R11 — WASM ABI

A new guest export `deploy(appId, constructorParams)` is introduced, called once at deploy time. It receives opaque bytes that only the guest knows how to parse, and returns the initial application state. This allows apps to initialize with arbitrary configuration (supported tokens, admin addresses, fee parameters, etc.) without the contract or executor needing to understand the configuration format.

The `deposit` guest export is extended to include the token address, so the guest knows which asset is being deposited: `deposit(appId, sender, tokenAddress, amount, state)`.

Withdrawal instructions produced by the guest include the token address.

### R12 — Events

All deposit, withdrawal, refund, and claim events emitted on-chain include the token address. This ensures external systems (subgraphs, explorers, monitoring) can distinguish between ETH and ERC-20 operations and correctly index asset-aware flows.

### R14 — Unsupported Token Behaviors

The following token behaviors will never be included in the whitelisted token list (ADMIN responsability, no on-chain check):
- **Fee-on-transfer**: rejected at deposit time via the balance check (see 2, R4).
- **Rebasing**: must be excluded from the allowlist (see 4, R2). No in-contract defense exists.
- **ERC-777**: must be excluded from the allowlist (see 3, R2). `nonReentrant` provides backup defense (R10).
- **Non-standard return values**: handled transparently by `SafeERC20` (see 5, R4).

### R15 — Admin Surplus Recovery (Optional, Deferred)

An admin function to recover unaccounted token surplus may be added in a future iteration. The recoverable amount would be `balanceOf(token) - totalAppCustody[token] - totalPendingClaims[token]`, ensuring only tokens not owed to anyone can be moved.

This would handle tokens sent to the contract by mistake (direct `transfer` instead of going through `submitRequest`), dust accumulation, and leftover balances after a token is removed from the allowlist. For ETH, the surplus calculation would also need to account for fee reservations held for in-queue requests.

This is not required for the initial implementation.

---


## Proposed Design

At a high level, the system should stop modeling the business asset as an implicit ETH amount and instead treat it as `(tokenAddress, amount)`.

The special case `tokenAddress = 0x0` preserves ETH as a first-class supported asset without needing a separate model.

The contract maintains a global token allowlist so it can reject unsupported assets before custodying them. 

Pending claimable balances are tracked by payee and asset. This same mechanism is used for:

- successful withdrawal claims, in ETH or ERC-20;
- refund of business-asset deposits when a request fails;
- refund of unused fees, always in ETH.

This keeps claim semantics uniform while preserving the rule that fees are still ETH-only.

## Proposed Data Model and APIs

The following shapes are proposed as the concrete target model for the refactor.

### Request Model

Requests should stop carrying a plain ETH-only `depositAmount` and instead carry a business asset defined by `(tokenAddress, assetAmount)`.

Conceptual request shape:

```text
Request {
  protocolVersion  uint8
  applicationId    uint64
  requestId        bytes32
  requestType      RequestType
  payload          bytes
  timestamp        uint256
  sender           address
  tokenAddress     address   // 0x0 = ETH
  assetAmount      uint256   // business asset amount
  maxFeeValue      uint256   // always ETH
}
```

Validation rule:

- if `assetAmount == 0`, then `tokenAddress` must be `0x0`.

Request identifier rule:

```text
requestId = keccak256(
  abi.encodePacked(
    sender,
    applicationId,
    requestType,
    payload,
    tokenAddress,
    assetAmount,
    idx
  )
)
```

This replaces the current ETH-only use of `depositAmount` in the request identifier calculation. `maxFeeValue` remains outside the request identifier, matching the current design where fee reservation is not part of request identity.

### Withdrawal Model

Withdrawals emitted by the executor and signed in the update payload must also become asset-aware.

Conceptual withdrawal shape:

```text
WithdrawalRequest {
  tokenAddress  address   // 0x0 = ETH
  receiver      address
  amount        uint256
}
```

### Pending Claim Accounting

On-chain claimable balances should be tracked per payee and per asset.

Conceptual storage:

```text
pendingClaims[tokenAddress][payee] uint256
totalPendingClaims[tokenAddress]   uint256
```

This allows the contract to:

- track independent ETH and ERC-20 claim balances for the same payee;
- validate available balances per asset;
- support the same claim model for withdrawals and refunds.

This model is also the recommended default for the implementation of asset-aware solvency checks.

This replaces the current ETH-only `payments[payee]` / `_totalDeposits` accounting model with an asset-aware claim accounting model. In particular, the current single `_totalDeposits` aggregate must be replaced by asset-scoped aggregates such as `totalPendingClaims[tokenAddress]`, so the contract can reason about ETH liabilities separately from each ERC-20 liability.

These asset-scoped aggregates should increase when value becomes claimable on-chain and decrease when a claim is executed.

In the multi-app model, this claim-facing storage is not sufficient on its own to protect custody boundaries between applications.

### App-Scoped Custody Accounting

In addition to `pendingClaims`, the contract should maintain app-scoped business-asset accounting so one application cannot create withdrawals backed by liquidity that entered through another application.

Conceptual storage:

```text
appCustody[applicationId][tokenAddress] uint256
totalAppCustody[tokenAddress]           uint256
```

Behavior:

- business-asset deposits increase `appCustody[applicationId][tokenAddress]` for the addressed app;
- business-asset outflows created by that app, such as withdrawals and failed-request business-asset refunds, decrease that app-scoped balance;
- app-scoped custody checks apply to business assets, while fee accounting remains ETH-only and bounded by request-specific fee reservation;
- claim execution itself does not need an `applicationId`, because the app-specific solvency check already happened when the pending claim was created.

This means the contract should use:

- `pendingClaims[tokenAddress][payee]` to track who can claim what asset;
- `appCustody[applicationId][tokenAddress]` to enforce which app is allowed to create those liabilities in the first place.

### Contract-Level Token Allowlist

The contract should maintain a global allowlist of tokens that are allowed to enter the system at all.

Conceptual storage:

```text
globalAllowedTokens[tokenAddress] bool
```

Behavior:

- `0x0` is always treated as ETH and always allowed;
- any non-zero `tokenAddress` must be present in `globalAllowedTokens`;
- app-level token policies are managed off-chain by each WASM application at runtime (see R3), not enforced by the contract.

Administration model:

- the global allowlist is the single source of truth for which ERC-20 contracts can ever enter custody;
- it should be managed on-chain by the existing `ADMIN` role, following the same access-control pattern already used for fee collector and deployer management;
- the contract should expose explicit admin functions for adding and removing globally allowed ERC-20 token addresses.
- when adding a token, the admin function must verify that the address contains contract code (`tokenAddress.code.length > 0`), as required by R2. See section 8 for details and optional code-hash hardening.

### Contract API

Conceptual request submission API:

```text
submitRequest(
  protocolVersion,
  applicationId,
  requestType,
  payload,
  tokenAddress,
  assetAmount,
  maxFeeValue
) payable
```

Conceptual generic claim API:

```text
claim(tokenAddress, payee)
```

The claim API should remain permissionless and behave as a no-op if the pending balance is zero, matching the current ETH-only behavior.

## Component Changes

### Smart Contracts

The on-chain request path must be extended so a request can specify the business asset using `tokenAddress` in addition to its amount.

For ETH requests:

- the business asset remains native ETH;
- `tokenAddress = 0x0`;
- `msg.value = assetAmount + maxFeeValue`.

For ERC-20 requests:

- the token is identified by a non-zero `tokenAddress`;
- `msg.value = maxFeeValue`;
- the contract pulls the business-asset amount with `transferFrom`;
- fees are still provided in ETH.

The request identifier must include both `tokenAddress` and `assetAmount`, since both are now part of the request identity. More precisely, the current `depositAmount` term in the request-id formula is removed and replaced by the ordered pair `(tokenAddress, assetAmount)`.

The contract must also reject requests if:

- `assetAmount == 0` and `tokenAddress != 0x0`;
- `tokenAddress` is non-zero and not globally allowlisted.

The on-chain pending-balance mechanism must become asset-aware so claims can be executed for either ETH or ERC-20 tokens.

The claim interface should also become asset-aware, for example with a generic `claim(tokenAddress, payee)` entrypoint.

In the target design, this generic claim flow replaces the current ETH-only `withdrawPayments(payee)` entrypoint rather than living beside it as a separate accounting path.

Claim semantics should be:

- if `tokenAddress == 0x0`, transfer native ETH;
- otherwise, transfer the ERC-20 token;
- clear pending balance before external transfer;
- emit a token-aware claim event;
- return without reverting if the pending amount is zero.

The error branch of `stateUpdate` must also become asset-aware. In particular:

- for a failed ERC-20 request, the business-asset deposit must be credited as a pending claim in `tokenAddress`, while the unused fee refund must be credited as a pending claim in `0x0`;
- for a failed ETH request, both components remain ETH and can be credited as a single ETH pending claim.

The contract should use safe ERC-20 transfers and document that only standard ERC-20 tokens are supported.

Asset solvency must be checked per asset, not as a single global sum. In particular:

- ETH liabilities must be checked against ETH available balance;
- ERC-20 liabilities must be checked against the available balance of that specific ERC-20;
- failed-request refunds must account separately for business-asset refunds and ETH fee refunds.

In the multi-app model, business-asset solvency must also be checked against the app-scoped custody balance for `(applicationId, tokenAddress)`, so one app cannot create withdrawals backed by funds deposited through another app.

### Common Types and Payloads

Common request and update types must move from ETH-only fields to asset-aware fields.

Conceptually, the request model must carry:

- business asset token address;
- business asset amount;
- ETH fee reservation.

The deploy payload format must also be extended so deployment can carry app configuration.

Withdrawal models must also include the token identity, not just destination and amount.

This affects both host-side types in `vela` and guest/shared WASM types in `vela-common-go`. In particular, the guest `Withdrawal` type used by WASM apps must be updated together with the host withdrawal types so token-aware withdrawals can be serialized consistently across the guest/host boundary.

The signed update payload must continue to treat:

- `refundAmount` as ETH-only;
- `applicationFee` as ETH-only;
- withdrawals as asset-aware objects containing `tokenAddress`.

For failed requests, the business-asset refund is not included in the signed update's `refundAmount` field. It is handled separately as a pending claim in the deposited asset. For ETH requests, both credits are denominated in ETH and combine into a single `pendingClaims[0x0][sender]` entry, consistent with R8.

### Executor

The executor must treat deposits as asset-aware instead of ETH-only.

If a request contains a business-asset deposit:

- the executor must pass both `tokenAddress` and `amount` to the runtime deposit flow;
- the resulting events must preserve the token information;
- failed requests must signal failure in the signed update so the contract can produce the appropriate pending claim balances. The contract credits the business-asset refund using the request's stored `tokenAddress` and `assetAmount`. The signed update's `refundAmount` covers only the ETH fee refund.

The executor update payload and its signature material must also include token-aware withdrawal data.

During deployment, the executor must:

1. decode the deploy payload, including app configuration;
2. pass the app configuration into the runtime `deploy` flow;
3. persist the guest-produced initial state that already contains the app configuration.

The executor should preserve the current high-level execution order:

1. process business-asset deposit into private state;
2. process the request payload;
3. if successful, emit token-aware withdrawals;
4. if failed, keep state unchanged and rely on claimable refunds.

This preserves the current mental model while making the asset explicit.

### Blockchain Client and Signature Builder

The blockchain client must be updated to:

- submit ETH requests with `msg.value = assetAmount + maxFeeValue`;
- submit ERC-20 requests with `msg.value = maxFeeValue`;
- prepare ERC-20 deposit requests only after the user has granted sufficient allowance to the processor contract;
- send token-aware withdrawals inside `stateUpdate`.

For ERC-20 deposits, the external user flow becomes a two-step sequence:

1. call `approve(processorEndpoint, assetAmount)` on the ERC-20 token, unless sufficient allowance already exists;
2. call `submitRequest(..., tokenAddress, assetAmount, maxFeeValue)`.

If allowance or token balance is insufficient, the contract-side `transferFrom` will revert and the request must not enter the queue or emit `RequestSubmitted`.

This is a meaningful UX change for clients and wallet tooling and must be reflected in the blockchain client and any deposit-oriented CLI flows.

The signature builder must hash token-aware withdrawal tuples, not ETH-only `(receiver, amount)` tuples.

The withdrawal tuple ABI must change in lockstep across the full signed-update path. In particular, the tuple shape and field order must remain identical across:

- executor-side signing code;
- contract-side `Structs.WithdrawalRequest`;
- contract-side signature verification in the TEE authenticator;
- `ProcessorEndpoint.stateUpdate`;
- generated bindings and blockchain client code.

The target tuple is:

```text
(tokenAddress, receiver, amount)
```

Any mismatch in tuple fields or field order across these components will cause signature verification to fail at `stateUpdate`.

Read-path impact:

- public claim-related read paths are in scope for this design and must become asset-aware, since pending claim balances and claim execution are no longer ETH-only;
- private-balance read paths remain app-specific and are not forced into a single common query model by this design;
- ETH-only apps such as `vela-nova` may keep their existing private-balance read model;
- subgraph and external indexing changes are a secondary integration concern unless needed to expose the new asset-aware claim flows.

### WASM Runtime

The runtime ABI for `deposit` must be updated so the guest receives:

- sender;
- token address;
- amount;
- state.

Without this change, a guest can receive a deposit amount but cannot know which asset the amount refers to.

Conceptually, the deposit ABI becomes:

```text
deposit(appId, sender, tokenAddress, amount, state)
```

The exact pointer-based host/guest ABI can follow the current runtime style, but the semantic inputs must include `tokenAddress`.

A new ABI must be added for deployment, for receiving constructor params. 

Conceptually:

```text
deploy(appId, params)
```

The guest should use this deploy-time configuration to build its initial state, eventually including the the app's allowlist.

## Main Flows

### ETH Deposit

1. A request is submitted with `tokenAddress = 0x0`, business-asset amount in ETH, and fees in ETH.
2. The contract receives the ETH business asset and the ETH fee reservation.
3. The executor passes `(tokenAddress = 0x0, amount)` into the WASM deposit flow.
4. The app credits the user balance for ETH in private state.

### ERC-20 Deposit

1. A request is submitted with a non-zero `tokenAddress`, business-asset amount, and fees in ETH.
2. The contract pulls the ERC-20 amount with `transferFrom`.
3. The executor passes `(tokenAddress, amount)` into the WASM deposit flow.
4. The app credits the user balance for that token in private state.

### Withdrawal

1. The user submits a process request for withdrawal of a given token and amount.
2. The app validates that the token is supported and that the user has enough balance for that token.
3. The app debits the private balance for `(account, token)`.
4. The executor emits a token-aware withdrawal request in the signed update payload.
5. The contract credits a pending claim balance for `(payee, tokenAddress)`.

### Failed Request

1. A request deposits a business asset and reserves ETH for fees.
2. Processing fails.
3. The business-asset deposit is credited back as a pending claim balance in its own asset.
4. The unused fee portion is credited back as a pending ETH claim balance.
5. If the business asset is ERC-20, the failed request therefore creates two pending credits in different assets: one ERC-20 credit and one ETH credit.

### Claim

1. A pending claim exists for `(payee, tokenAddress)`.
2. Any caller can invoke `claim(tokenAddress, payee)`.
3. The contract reads the pending amount for that asset and payee.
4. If the amount is zero, the function returns without reverting.
5. The contract clears the pending balance before the transfer.
6. The contract transfers ETH if `tokenAddress = 0x0`, otherwise it transfers ERC-20 tokens.
7. The contract emits a token-aware claim event.

## Risks and Edge Cases

- A request may specify a token not present in the app allowlist.
- A request may specify a token that is globally allowlisted but not allowed by the app.
- ERC-20 transfers may fail or revert during deposit or claim flows.
- Non-standard ERC-20 tokens are explicitly unsupported and must be rejected or documented as unsupported behavior.
- Pending claim balances may exist simultaneously for multiple assets for the same payee.
- Solvency checks must be done per asset, not as a single aggregated amount.
- Event schemas must remain unambiguous once token identity is added.
- In the multi-app model, business-asset custody must be scoped by both `applicationId` and `tokenAddress` to avoid one app creating withdrawals backed by another app's deposits.


## ERC-20 Token Behaviors That Affect Vela

This section documents ERC-20 token behaviors that were analyzed during the design phase. Each behavior is described alongside its impact on Vela and the defense strategy adopted.

### 1 Standard ERC-20 Deposit Pattern

Vela needs the two-step delegation pattern to pull tokens from users during request submission:

1. User calls `approve(processorEndpoint, amount)` on the token contract.
2. User calls `submitRequest(...)` on `ProcessorEndpoint`, which internally calls `transferFrom(user, address(this), amount)`.

This is the standard approach used across DeFi. The user must have granted sufficient allowance before submitting. If allowance or balance is insufficient, `transferFrom` reverts and the request never enters the queue.

### 2 Fee-on-Transfer Tokens

Some ERC-20 tokens deduct a percentage on every `transfer` or `transferFrom` call. The token contract itself takes a cut before delivering to the recipient.

**Deposit-side impact:** the contract calls `transferFrom(user, address(this), assetAmount)` and records `assetAmount` in custody, but actually receives `assetAmount - fee`. Over many deposits, the deficit grows. Eventually claims fail because the contract does not hold enough tokens to pay out.

**Claim-side impact:** the contract transfers `amount` to the payee, but the payee receives `amount - fee`. There is no contract-level defense for this direction; the fee is applied by the token itself.

Examples of tokens with fee-on-transfer behavior include STA (Statera), SafeMoon, and PAXG. USDT has a fee mechanism built in that is currently set to zero but can be enabled by the token issuer at any time.

**Recommended mitigation:** in addition to the allowlist, the contract should verify the actual received amount on deposit:

```text
balanceBefore = IERC20(token).balanceOf(address(this))
IERC20(token).safeTransferFrom(sender, address(this), assetAmount)
received = IERC20(token).balanceOf(address(this)) - balanceBefore
require(received == assetAmount, "fee-on-transfer not supported")
```

This costs approximately 2600 additional gas (two `balanceOf` SLOAD calls), which is negligible relative to the rest of `submitRequest`. It turns a silent accounting mismatch into a loud revert and catches cases where a previously standard token activates fees after being allowlisted.

### 3 ERC-777 Backwards-Compatible Tokens

ERC-777 is a token standard that implements the full ERC-20 interface but adds transfer hooks (`tokensToSend` on the sender side, `tokensReceived` on the recipient side) registered via the ERC-1820 registry. During a `transferFrom` call, execution control passes to external code through these hooks, creating a reentrancy vector.

A token that appears to be a standard ERC-20 may actually be an ERC-777. If such a token is allowlisted, a malicious user can register a `tokensToSend` hook and re-enter the contract during deposit processing.

This class of attack was exploited in practice against Uniswap V1 in 2020 via the imBTC token (an ERC-777).

**Impact on Vela:** A malicious user could register a hook and re-enter the contract during deposit processing.

**Mitigation:** the global allowlist must not include ERC-777 tokens. The `nonReentrant` modifier on all state-modifying entry points provides a backup defense. The checks-effects-interactions pattern should also be followed so that state updates occur before any external ERC-20 call.

### 4 Rebasing Tokens

Rebasing tokens automatically adjust every holder's balance to track a target value (a peg, an index, or a yield accrual) without any explicit transfer. The token contract periodically scales every holder's `balanceOf()` result up or down by a global factor. Examples include AMPL, stETH, and OHM.

**Impact on Vela:** A negative rebase reduces the contract's actual token balance while the internal accounting (`appCustody`, `pendingClaims`) remains unchanged. The contract becomes insolvent silently over time. A positive rebase creates unaccounted surplus. Unlike fee-on-transfer, this drift happens *between* transactions, so a per-deposit balance check cannot detect it.

**Defense:** Rebasing tokens must not be added to the global allowlist. There is no viable in-contract defense.

### 5 Non-Standard Return Values

The ERC-20 specification requires `transfer` and `transferFrom` to return a `bool`. Some widely used tokens, notably USDT, do not return any value. A raw `IERC20` call reverts on these tokens because Solidity expects a return value.

**Impact on Vela:** Any direct use of `IERC20.transfer` or `IERC20.transferFrom` would fail for USDT and similar tokens.

**Defense:** All ERC-20 interactions use OpenZeppelin's `SafeERC20` library (`safeTransfer`, `safeTransferFrom`), which handles both tokens that return `true` and tokens that return nothing.

### 6 Upgradeable Proxy Tokens

Major tokens such as USDC and USDT are deployed behind upgradeable proxy contracts. The token issuer can change the implementation at any time, potentially adding fees, pausing functionality, changing approval semantics, or introducing other behavioral changes.

**Impact on Vela:** A previously standard token could start behaving non-standardly after an upgrade.

**Mitigation:** this is an operational risk. The administrator managing the global allowlist should monitor upgrades to allowlisted tokens. The deposit-side balance check described above catches some behavioral changes early by reverting if the received amount does not match the expected amount.

### 7 Pausable Tokens and Address Blocklists

USDC and USDT have `pause()` functions that halt all transfers and address blocklist mechanisms that prevent specific addresses from sending or receiving tokens.

When a token is paused:

- deposits revert at `transferFrom` inside `submitRequest`, so no request enters the queue;
- claims revert at the outbound transfer, so the payee cannot withdraw until the token is unpaused;
- `stateUpdate` is unaffected because `_asyncTransfer` only updates internal accounting without making an external ERC-20 call.

If a payee address is on the token's blocklist, their claim permanently reverts and the funds remain in `pendingClaims`. This is consistent with the pull-payment pattern, which places the burden of receivability on the payee.

If the ProcessorEndpoint contract itself is added to a token's blocklist, all operations for that token break. This is an external dependency risk with no in-contract defense.

**Impact on Vela:**
- When a token is paused, deposits revert at `transferFrom` (no request enters the queue) and claims revert at the outbound transfer (funds are safe but temporarily inaccessible).
- If a payee address is on the token's blocklist, their claim permanently reverts. Funds remain in `pendingClaims` indefinitely.
- If the `ProcessorEndpoint` contract itself is blocklisted, all operations for that token break entirely.

**Mitigation:** the pull-payment pattern already handles pauses gracefully. Blocklist and contract-blocklist risks should be documented as known external dependencies for allowlist administrators. A permanently blocked payee's funds remain in the contract; this is consistent with the pull-payment model where receivability is the payee's responsibility.

### 8 Contract Address Validation

An Ethereum address does not guarantee the presence of a valid ERC-20 contract. The address could point to an externally owned account (no code), a non-token contract, or a self-destructed contract. Calling `transferFrom` on an address with no code succeeds silently with empty return data in some cases, which can lead to unexpected behavior.

Additionally, with `CREATE2`, a contract can be deployed at a deterministic address, self-destructed, and a different contract redeployed at the same address. If the original address was allowlisted, the system would now trust a completely different contract.

**Impact on Vela:** An admin could accidentally allowlist an invalid address, leading to silent failures.

**Recommended mitigation:** the admin function that adds tokens to the global allowlist should verify that the address contains code:

```text
require(tokenAddress.code.length > 0, "not a contract")
```

As optional hardening, the contract could record the code hash at allowlist time and verify it at request time to detect contract redeployment:

```text
allowedTokenCodeHash[token] = token.codehash   // at allowlist time
require(token.codehash == allowedTokenCodeHash[token])  // at request time
```

### 9 ETH Overpayment on ERC-20 Requests

For ERC-20 requests, `msg.value` must equal `maxFeeValue` exactly. The business asset arrives via `transferFrom`, not via `msg.value`. If a caller sends `msg.value > maxFeeValue` on an ERC-20 request, the excess ETH has no accounting path and is effectively stuck in the contract.

**Recommended mitigation:** explicit validation in `submitRequest`:

```text
if (tokenAddress != address(0)) {
    require(msg.value == maxFeeValue, "ETH value must equal maxFeeValue for ERC-20 requests")
}
```

### 10 Zero-Amount Transfers

Some ERC-20 tokens revert on zero-amount transfers. The claim function must check the pending balance before attempting a transfer, not after.

### 11 EIP-2612 Permit

Some tokens (including USDC v2) support EIP-2612 `permit`, which allows a gasless off-chain signature to grant approval. This turns the standard two-transaction deposit flow (`approve` + `submitRequest`) into a single transaction. This is a UX improvement that can be added as a follow-up without architectural impact on the rest of the system.

### 12 Token Decimals

ERC-20 tokens use different decimal precisions (e.g., 18 for DAI, 6 for USDC, 8 for WBTC). The contract is decimal-agnostic and operates exclusively on raw uint256 amounts. However, WASM applications that perform internal accounting must know the decimals of each token they support to avoid incorrect balance calculations. This information should be included in the deploy-time constructor parameters passed to the application.


## Testing Plan

Unit tests should cover:

- request and withdrawal type validation with `tokenAddress`;
- allowlist validation;
- ETH and ERC-20 balance accounting;
- unsupported token rejection;
- failed-request refund logic;
- claim accounting by `(payee, tokenAddress)`.

Integration tests should cover:

- updated WASM deposit ABI with `tokenAddress`;
- ETH deposit and withdrawal flows;
- ERC-20 deposit and withdrawal flows;
- mixed refund scenarios where business asset and fee refunds use different assets.

System tests should cover:

- ETH-only apps rejecting unsupported ERC-20 deposits;
- `vela-nova` remaining ETH-only while still working with the updated shared ABI and types;
- ERC-20 business-asset deposit followed by withdrawal and claim;
- failed ERC-20 requests producing ERC-20 business-asset claims plus ETH fee refunds.
- in the multi-app model, one app being unable to create withdrawals backed by another app's deposited business assets.

## Acceptance Criteria

- A request can carry exactly one business asset identified by `tokenAddress`, where `0x0` represents ETH.
- An ETH request uses `msg.value = assetAmount + maxFeeValue`.
- An ERC-20 request uses `msg.value = maxFeeValue` and transfers the business asset with `transferFrom`.
- The runtime deposit path receives both `tokenAddress` and `amount`.
- Private app state can store balances per account and per token.
- ETH-only apps reject non-zero `tokenAddress` values.
- The contract rejects non-zero `tokenAddress` values that are not globally allowlisted.
- Successful withdrawals create claimable balances in the correct asset.
- Failed requests refund the business asset in its own asset and the unused fees in ETH.
- Claim execution is permissionless and asset-aware.
- Claiming a zero pending balance does not revert.
- In the multi-app model, an app cannot create business-asset withdrawals backed by another app's custody.
- Withdrawal signature material includes `tokenAddress`.
- The withdrawal tuple ABI used by executor signing and on-chain signature verification is identical.
- Events emitted for deposits, refunds, withdrawals, and claims include the asset identity.
- Only standard ERC-20 tokens are considered supported by this design.
- Withdrawals of delisted tokens are always possibile


## Implementation Tasks

### 1. Smart Contracts
- Add global token allowlist storage and admin functions (`addAllowedToken`, `removeAllowedToken`) with contract-code validation
- Refactor `submitRequest` to accept `tokenAddress` + `assetAmount`; handle ETH vs ERC-20 deposit paths (balance-before/after check, `SafeERC20.safeTransferFrom`)
- Add `msg.value == maxFeeValue` validation for ERC-20 requests; enforce `assetAmount == 0 → tokenAddress == 0x0`
- Add per-app per-token custody accounting (`appCustody[appId][token]`, `totalAppCustody[token]`)
- Refactor pending claim storage to `pendingClaims[token][payee]` + `totalPendingClaims[token]`; replace `withdrawPayments(payee)` with `claim(tokenAddress, payee)`
- Update withdrawal model to `(tokenAddress, receiver, amount)` tuple — must match across Structs, signature verification, and stateUpdate
- Update `stateUpdate` for asset-aware solvency checks, failed-request dual refunds (ERC-20 asset + ETH fee), and token-aware withdrawal processing
- Add `nonReentrant` to all state-modifying entry points if not already present
- Update request ID calculation to include `tokenAddress` + `assetAmount` (replacing `depositAmount`)
- Update all events to include `tokenAddress`
- Write Hardhat tests: allowlist, ETH/ERC-20 deposits, withdrawals, claims, failed-request refunds, solvency, reentrancy, unsupported tokens
- Regenerate Go bindings (`go generate`) and commit

### 2. Common Types & Payloads
- Update request model to carry `tokenAddress` + `assetAmount` instead of `depositAmount`
- Update withdrawal type to include `tokenAddress` (both host-side in `pkg/` and guest-side in `vela-common-go`)
- Add deploy payload type with constructor params
- Update signed-update payload: withdrawals become token-aware; `refundAmount` stays ETH-only

### 3. Common Go Libraries (`vela-common-go`)
- Update shared guest `Withdrawal` type to include `tokenAddress`
- Update shared request types / serialization to carry `tokenAddress` + `assetAmount`
- Add deploy-related shared types (constructor params, deploy response)
- Ensure serialization across guest/host boundary is consistent for all updated types

### 4. WASM Runtime
- Extend `deposit` guest export ABI: `deposit(appId, sender, tokenAddress, amount, state)`
- Add `deploy(appId, params)` guest export for deploy-time configuration
- Update host-side runtime to pass `tokenAddress` through to the guest

### 5. Executor
- Pass `tokenAddress` + `amount` into the runtime deposit flow
- Build token-aware withdrawal instructions in the signed update payload
- Implement deploy flow: decode deploy payload → call runtime `deploy` → persist initial state
- Ensure failed-request signaling lets the contract produce correct dual refunds

### 6. Blockchain Client & Signature Builder
- Update request submission: ETH (`msg.value = assetAmount + maxFeeValue`) vs ERC-20 (`msg.value = maxFeeValue` + `transferFrom`)
- Update signature builder to hash `(tokenAddress, receiver, amount)` withdrawal tuples
- Ensure tuple ABI is identical across executor signing, contract verification, and bindings

### 7. Simple App (`app/simple`)
- Update to new `deposit(appId, sender, tokenAddress, amount, state)` ABI
- Implement `deploy(appId, params)` export to initialize state from constructor params
- Add per-token balance tracking in private state
- Support token-aware withdrawal instructions
- Add app-level token allowlist configured via deploy params

### 8. Nova App (`vela-nova`)
- Update to new `deposit` and `deploy` ABIs
- Keep ETH-only behavior: reject non-zero `tokenAddress` in deposit
- Initialize via `deploy` with ETH-only configuration
- Verify existing functionality is unbroken with the updated shared types

### 9. Testing (Integration & System)
- Integration tests: updated WASM deposit ABI, ETH and ERC-20 deposit/withdrawal flows, mixed refund scenarios
- System tests: ETH-only app (nova) rejects ERC-20 deposits, ERC-20 full lifecycle (deposit → withdraw → claim), failed ERC-20 request dual refunds, cross-app custody isolation
