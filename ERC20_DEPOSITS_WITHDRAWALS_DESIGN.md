# ERC-20 Support for Deposits / Withdrawals / Claims

## Summary

This document describes the design to support private-state deposits, withdrawals, and claims for selected ERC-20 tokens in addition to native ETH.

Fees remain ETH-only and are explicitly out of scope for ERC-20 support in this design. The system will support apps that are ETH-only as well as apps that support ETH plus an allowlist of ERC-20 tokens.

In this document, `0x0` is used as shorthand for the zero address. In Solidity this means `address(0)`. In Go this means the zero-value Ethereum address.

The design is intended to remain valid when the system moves from the current single-app deployment model to a future multi-app model where a shared contract can serve multiple applications.

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
- Per-app token support policies.

In addition, the current codebase still behaves as a single-app system in practice, but that is expected to change soon. The design in this document should therefore be read as the target model for a multi-app-capable system, not as a restatement of the current single-app limitation.

## Goals

- Support one business asset per request, where the asset can be ETH or a selected ERC-20 token.
- Represent the asset using `tokenAddress`, with `0x0` reserved for ETH.
- Allow deposits into private state for ETH and supported ERC-20 tokens.
- Allow withdrawals from private state for ETH and supported ERC-20 tokens.
- Support claim flows for both ETH and ERC-20 pending balances.
- Allow each WASM app to define which tokens it supports.
- Keep fees denominated and settled in ETH only.

## Non-Goals

- Paying application fees in ERC-20 tokens.
- Supporting more than one business asset in a single request.
- Supporting non-standard ERC-20 behaviors such as fee-on-transfer, rebasing, or other custom semantics.
- Preserving backwards compatibility with the current ETH-only request and withdrawal formats.

## Decisions

- A request carries exactly one business asset.
- The asset is represented as `tokenAddress`, with `0x0` meaning ETH.
- If `assetAmount == 0`, `tokenAddress` must be `0x0`.
- ERC-20 deposits use the standard `approve + transferFrom` flow.
- Fees remain ETH-only.
- If a request fails, the business-asset deposit is returned via claim.
- Refunds for unused fees are also claimed, but always in ETH.
- Claims use a generic asset-aware mechanism, not separate ETH-only and ERC-20-only designs.
- Pending claimable balances are tracked by `(payee, tokenAddress)`.
- Claim execution remains permissionless: anyone can trigger a claim for a payee.
- Token validation is performed both on-chain and in the app runtime.
- App token support follows a hybrid model:
  - The contract enforces a global system allowlist of tokens that can be custodied at all.
  - Each app declares its own allowlist, scoped by `applicationId`, as a subset of that global allowlist.
  - The app also validates supported tokens at runtime.
- Token support is allowlist-based.
- ETH-only apps use an allowlist containing only `0x0`.
- ETH-plus-ERC20 apps use an allowlist containing `0x0` and the supported ERC-20 token addresses.
- Private balances move from `account -> balance` to `account -> token -> balance`.
- The WASM deposit ABI must be extended to include `tokenAddress`.
- Deposit, refund, withdrawal, and claim-related events must include the token.

## Proposed Design

At a high level, the system should stop modeling the business asset as an implicit ETH amount and instead treat it as `(tokenAddress, amount)`.

The special case `tokenAddress = 0x0` preserves ETH as a first-class supported asset without needing a separate model.

Each app defines an allowlist of supported tokens. That allowlist is declared in app configuration or deploy metadata and is also enforced by the application logic at runtime. In addition, the contract maintains a global token allowlist so it can reject unsupported assets before custodying them. This allows some apps to stay ETH-only while others support ETH plus selected ERC-20 tokens.

Whether a shared contract accepts a given ERC-20 for a given request therefore depends on two checks:

- the token must be globally allowlisted by the system;
- the token must be enabled for the specific `applicationId` being addressed by the request.

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
```

Behavior:

- business-asset deposits increase `appCustody[applicationId][tokenAddress]` for the addressed app;
- business-asset outflows created by that app, such as withdrawals and failed-request business-asset refunds, decrease that app-scoped balance;
- app-scoped custody checks apply to business assets, while fee accounting remains ETH-only and bounded by request-specific fee reservation;
- claim execution itself does not need an `applicationId`, because the app-specific solvency check already happened when the pending claim was created.

This means the contract should use:

- `pendingClaims[tokenAddress][payee]` to track who can claim what asset;
- `appCustody[applicationId][tokenAddress]` to enforce which app is allowed to create those liabilities in the first place.

### App Config

App token support should be serialized as deploy-time app configuration and delivered to the guest during deployment, not compiled into the WASM binary.

Conceptual app config shape:

```text
AppConfig {
  allowedTokens []address  // 0x0 means ETH
}
```

Examples:

- ETH-only app: `allowedTokens = [0x0]`
- ETH + ERC-20 app: `allowedTokens = [0x0, usdc, dai, ...]`

This app-level allowlist is assumed to be a subset of the contract-level global allowlist and must be scoped by `applicationId` in the multi-app model.

Delivery model:

- the deploy payload must carry app configuration, including `allowedTokens`;
- the executor/runtime must pass that configuration into `load_module`;
- the guest must persist the configuration inside its initial private state;
- subsequent `deposit` and `process_request` calls validate token support by reading that state.

This is preferred over having the host mutate app-specific state after deployment, because the host should not need to understand each app's internal state schema.

This app config is only the per-app allowlist needed by the guest. The guest should not persist or replicate the contract-wide global token allowlist, which remains an on-chain source of truth.

### Contract-Level Token Allowlist

The contract should maintain a global allowlist of tokens that are allowed to enter the system at all.

Conceptual storage:

```text
globalAllowedTokens[tokenAddress] bool
```

Behavior:

- `0x0` is always treated as ETH and always allowed;
- any non-zero `tokenAddress` must be present in `globalAllowedTokens`;
- app-level allowlists can only use tokens that are already globally allowlisted.

Administration model:

- the global allowlist is the single source of truth for which ERC-20 contracts can ever enter custody;
- it should be managed on-chain by the existing `ADMIN` role, following the same access-control pattern already used for fee collector and deployer management;
- the contract should expose explicit admin functions for adding and removing globally allowed ERC-20 token addresses.

### Per-App Token Allowlist On-Chain

Because the system is expected to move to a multi-app model, the contract should also maintain app-scoped token support.

Conceptual storage:

```text
appAllowedTokens[applicationId][tokenAddress] bool
```

Behavior:

- a request must pass both the global token allowlist and the per-app token allowlist;
- different applications can support different token subsets on the same shared contract;
- an ETH-only app would have `appAllowedTokens[applicationId][0x0] = true` and no ERC-20 entries;
- an ETH + ERC-20 app would enable `0x0` and the selected token addresses for its own `applicationId`.

Registration model:

- the per-app token allowlist should be registered on-chain as part of the app deployment flow, using the `allowedTokens` declared in deploy-time app configuration;
- during deployment, the contract must reject any app token allowlist entry that is not already present in the global allowlist;
- for this design, app token support is treated as deploy-time configuration rather than a mutable runtime setting.

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
- `tokenAddress` is non-zero and not globally allowlisted;
- the addressed `applicationId` does not allow that token.

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

The deploy payload format must also be extended so deployment can carry app configuration, including `allowedTokens`.

Withdrawal models must also include the token identity, not just destination and amount.

This affects both host-side types in `vela` and guest/shared WASM types in `vela-common-go`. In particular, the guest `Withdrawal` type used by WASM apps must be updated together with the host withdrawal types so token-aware withdrawals can be serialized consistently across the guest/host boundary.

The signed update payload must continue to treat:

- `refundAmount` as ETH-only;
- `applicationFee` as ETH-only;
- withdrawals as asset-aware objects containing `tokenAddress`.

For failed requests, the business-asset refund is not represented as an ETH refund. It must be handled as a pending claim in the deposited asset.

### Executor

The executor must treat deposits as asset-aware instead of ETH-only.

If a request contains a business-asset deposit:

- the executor must pass both `tokenAddress` and `amount` to the runtime deposit flow;
- the resulting events must preserve the token information;
- failed requests must produce the appropriate pending claim balances for both the business asset and ETH fee refund.

The executor update payload and its signature material must also include token-aware withdrawal data.

During deployment, the executor must:

1. decode the deploy payload, including app configuration;
2. validate that the app allowlist is compatible with the contract-side token allowlists;
3. pass the app configuration into the runtime `load_module` flow;
4. persist the guest-produced initial state that already contains the app configuration.

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

The deployment ABI must also be extended so `load_module` receives deploy-time app configuration.

Conceptually:

```text
load_module(appId, appConfig)
```

The guest should use this deploy-time configuration to build its initial state, including the persisted allowlist that will later be enforced during request processing.

### WASM Applications

WASM apps must validate that the request token is supported by their configured allowlist, even if the contract already performed token checks.

App state must evolve from a single balance per account to balances per account and per token.

App state must also persist the deploy-time token configuration so runtime validation does not depend on any external side channel after deployment.

Apps that do not want ERC-20 support can remain ETH-only by allowing only `0x0`.

`vela-nova` is the main ETH-only reference app for this design. Even if its business logic remains ETH-only, it must still adopt the shared ABI and type changes introduced by this work, including:

- accepting the updated deploy-time configuration flow;
- accepting the updated `deposit` ABI with `tokenAddress`;
- rejecting any non-zero `tokenAddress` at runtime;
- emitting token-aware withdrawals using `tokenAddress = 0x0`.

### Simple App

`simpleapp` should be updated to demonstrate token-aware deposits and withdrawals and to serve as the reference implementation for the new model.

It should:

- support ETH and a configured allowlist of ERC-20 tokens;
- store balances per account and token;
- make its private withdrawal instruction token-aware, for example by moving from `WithdrawInstruction { to, amount }` to `WithdrawInstruction { tokenAddress, to, amount }`;
- emit token-aware deposit and withdrawal events;
- reject unsupported tokens cleanly.

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
- Configuration-level allowlists and runtime validation can diverge if not kept consistent.
- ERC-20 transfers may fail or revert during deposit or claim flows.
- Non-standard ERC-20 tokens are explicitly unsupported and must be rejected or documented as unsupported behavior.
- Pending claim balances may exist simultaneously for multiple assets for the same payee.
- Solvency checks must be done per asset, not as a single aggregated amount.
- Event schemas must remain unambiguous once token identity is added.
- In the multi-app model, contract-side token configuration must be correctly scoped by `applicationId` to avoid one app accidentally inheriting another app's token permissions.
- In the multi-app model, business-asset custody must be scoped by both `applicationId` and `tokenAddress` to avoid one app creating withdrawals backed by another app's deposits.

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
- ETH-plus-ERC20 apps accepting only allowlisted tokens;
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
- ETH + ERC-20 apps accept only allowlisted token addresses.
- The contract rejects tokens not enabled for the targeted `applicationId`.
- Successful withdrawals create claimable balances in the correct asset.
- Failed requests refund the business asset in its own asset and the unused fees in ETH.
- Claim execution is permissionless and asset-aware.
- Claiming a zero pending balance does not revert.
- In the multi-app model, an app cannot create business-asset withdrawals backed by another app's custody.
- Withdrawal signature material includes `tokenAddress`.
- The withdrawal tuple ABI used by executor signing and on-chain signature verification is identical.
- Events emitted for deposits, refunds, withdrawals, and claims include the asset identity.
- Only standard ERC-20 tokens are considered supported by this design.

## Open Questions

There are no remaining blocking design questions for this ticket.

Implementation should use the following defaults unless a concrete technical issue appears during coding:

- app token support is serialized in deploy-time config as `allowedTokens []address`;
- deploy-time app config is carried in the deploy payload and passed into `load_module`;
- the guest persists `allowedTokens` inside its initial private state;
- the contract maintains a global token allowlist and enforces it on-chain;
- the contract also maintains an app-scoped token allowlist keyed by `applicationId`;
- the generic claim entrypoint is `claim(tokenAddress, payee)`;
- the request model carries `tokenAddress`, `assetAmount`, and `maxFeeValue`;
- the withdrawal model carries `tokenAddress`, `receiver`, and `amount`.
