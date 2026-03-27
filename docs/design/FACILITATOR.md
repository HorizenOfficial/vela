# Facilitator: Gasless Request Submission

## 1. Prerequisites

This design builds on the **ERC-20 Deposits/Withdrawals/Claims** change ([design doc](https://github.com/HorizenOfficial/vela/blob/9ba796c236867a1c27277c3cb1a8851d6f3c4212/ERC20_DEPOSITS_WITHDRAWALS_DESIGN.md)), which must be implemented first. After that change:

- `submitRequest` accepts a `tokenAddress` + `assetAmount` pair for the business-asset deposit (`0x0` = ETH, non-zero = ERC-20)
- For ERC-20 deposits, the contract uses `approve + transferFrom` to pull tokens from the caller
- `maxFeeValue` **remains ETH-only**, paid via `msg.value`
- Claims replace pull payments: `claim(tokenAddress, payee)` is the generic, permissionless, asset-aware claim mechanism
- Pending claims are tracked by `(payee, tokenAddress)`
- Token allowlists (global + per-app) gate which tokens are accepted

Additional prerequisite:

- **Token requirement**: only ERC-20 tokens that implement **EIP-3009** (`transferWithAuthorization`) are supported for facilitated deposits. 
Primary target: **USDC**.

## 2. Problem Statement

To submit a request via `ProcessorEndpoint.submitRequest`, a user must hold ETH to cover **gas fees** and **service fees** (`maxFeeValue`). 
This creates onboarding friction: users need to acquire ETH before they can interact with the platform, even when they have sufficient ERC-20 tokens for their application deposit.

A **facilitator** should allow users to submit requests without holding any ETH. The facilitator absorbs gas and service fees as an operational cost, while the user authorizes only their business-asset deposit via an off-chain signature.

## 3. Goals

- Users can submit requests **without holding ETH**
- The facilitator covers **gas fees** (paid to Ethereum validators) and **service fees** (`maxFeeValue`, always in ETH)
- The user pays **only `assetAmount`** in ERC-20 tokens, authorized via off-chain EIP-3009 signature
- The **real user address** is tracked as the request sender (for key association, deposits in WASM, claims)
- The existing **direct EOA submission** path remains fully functional
- The facilitator is a platform-operated service that **absorbs costs** (no on-chain reimbursement mechanism)


This document covers the **smart contract changes** and **Go backend changes** needed in this repository to support facilitated requests. The off-chain facilitator service and client SDK will be implemented in a separate [`vela-x402`](https://github.com/HorizenOfficial/vela-x402) repository, which builds on the x402 HTTP payment protocol defined by Coinbase (see [https://github.com/coinbase/x402/](https://github.com/coinbase/x402/)).
More info on the off-chain facilitator will be provided below in paragraph 5.3 .

## 4. Approach: EIP-3009 + EIP-712 Meta-Transaction

### 4.1. Why EIP-3009

EIP-3009 defines `transferWithAuthorization` — a standard that allows a token holder to authorize a transfer via off-chain signature. A third party (the facilitator) can then execute the transfer on-chain without the token holder spending gas.

This fits our model: the user signs an authorization to transfer `assetAmount` tokens to the contract, and the facilitator submits it.

### 4.2. Why not ERC-4337

ERC-4337 (Account Abstraction) is more powerful but introduces significant infrastructure overhead:
- Requires deploying smart contract wallets for users
- Requires a Bundler service and EntryPoint contract
- Paymaster contracts add another trust/upgrade surface
- Overkill when the only gasless operation is `submitRequest`

ERC-4337 could be revisited if the platform needs broader account abstraction (e.g., gasless contract wallet interactions beyond Vela).

### 4.3. Why not EIP-2771

EIP-2771 (Trusted Forwarder) solves `msg.sender` identity but doesn't solve payment delegation for ERC-20 deposits. Since we already need EIP-3009 for the deposit token transfer, we can derive user identity from the EIP-712 request signature directly, without introducing a forwarder dependency.

### 4.4. Scope: ERC-20 Deposits Only

The facilitator pattern is designed for **ERC-20 business-asset deposits**. For ETH deposits (`tokenAddress = 0x0`), the user inherently needs ETH for `assetAmount` and therefore likely has enough for gas and fees — making the facilitator less valuable. ETH deposits remain supported only via the direct `submitRequest` path.

## 5. Design

### 5.1. User Signature Scheme

The user produces **two signatures** off-chain:

#### Signature 1 — Request Authorization (EIP-712)

An EIP-712 typed data signature over the request parameters, proving the user authorized this specific request:

```
RequestAuthorization {
    uint8    protocolVersion
    uint64   applicationId
    uint8    requestType
    bytes32  payloadHash        // keccak256(payload)
    address  tokenAddress       // ERC-20 token for business-asset deposit
    uint256  assetAmount        // business-asset deposit amount
    uint256  nonce              // per-user replay protection
    uint256  deadline           // expiration timestamp
}
```

The contract recovers the user address from this signature → this becomes the `sender` in `PendingRequest`.

#### Signature 2 — Deposit Transfer (EIP-3009)

A standard EIP-3009 `transferWithAuthorization` signature:

```
TransferWithAuthorization {
    address from               // user
    address to                 // ProcessorEndpoint contract
    uint256 value              // assetAmount
    uint256 validAfter         // 0 (immediate)
    uint256 validBefore        // deadline (matches request)
    bytes32 nonce              // unique nonce (EIP-3009 uses random nonces)
}
```

This authorizes the contract to pull `assetAmount` tokens from the user.

> **When `assetAmount == 0`**: Signature 2 is not required. Many request types (e.g., `PROCESS`, `ASSOCIATEKEY`) may not require a deposit, in which case the user only signs the request authorization.

### 5.2. Contract Changes

#### New function: `submitRequestFor`

```solidity
function submitRequestFor(
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    // User's EIP-712 request authorization
    uint256 nonce,
    uint256 deadline,
    bytes calldata requestSignature,
    // EIP-3009 deposit authorization (empty if assetAmount == 0)
    bytes calldata depositAuthorization
)
    external
    payable                    // msg.value = maxFeeValue (ETH)
    nonReentrant
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    returns (bytes32)
```

**Execution flow (solidty code):**

```
submitRequestFor()
  │
  ├─ 1. Verify deadline not expired
  │     require(block.timestamp <= deadline)
  │
  ├─ 2. Recover user address from EIP-712 request signature
  │     user = ecrecover(hashTypedData(requestAuth), requestSignature)
  │
  ├─ 3. Verify and consume nonce (replay protection)
  │     require(nonces[user] == nonce)
  │     nonces[user]++
  │
  ├─ 4. Validate token allowlists
  │     require(globalAllowedTokens[tokenAddress])
  │     require(appAllowedTokens[applicationId][tokenAddress])
  │
  ├─ 5. Validate fee
  │     maxFeeValue = msg.value
  │     require(maxFeeValue >= minFeePerRequest)
  │
  ├─ 6. If assetAmount > 0: execute EIP-3009 transfer
  │     token.transferWithAuthorization(
  │         user, address(this), assetAmount,
  │         validAfter, validBefore, eip3009Nonce, depositAuthorization
  │     )
  │     appCustody[applicationId][tokenAddress] += assetAmount
  │
  ├─ 7. Create PendingRequest with sender = user (not msg.sender)
  │     pendingRequest.sender = user
  │     pendingRequest.facilitator = msg.sender
  │     pendingRequest.tokenAddress = tokenAddress
  │     pendingRequest.assetAmount = assetAmount
  │     pendingRequest.maxFeeValue = maxFeeValue
  │
  ├─ 8. Generate requestId using user address
  │     requestId = _generateRequestId(user, applicationId, requestType,
  │                                     payload, tokenAddress, assetAmount, idx)
  │
  └─ 9. Emit RequestSubmitted(applicationId, requestId, user)
```

#### Nonce management

```solidity
// Sequential nonces per user (simpler than EIP-3009's random nonces)
mapping(address => uint256) public facilitatorNonces;

function getFacilitatorNonce(address user) external view returns (uint256);
```

**Why a dedicated nonce is needed.** The `facilitatorNonces` counter is independent from both the Ethereum account nonce (used by the facilitator's EOA to order on-chain transactions) and the EIP-3009 random nonce (managed by the token contract for `transferWithAuthorization`). It protects the EIP-712 request authorization (Signature 1) against replay. This is critical when `assetAmount == 0` (e.g., `ASSOCIATEKEY` or `PROCESS` requests): in that case there is no Signature 2 / EIP-3009 transfer, so without this nonce there would be **no replay protection at all** — a facilitator could re-submit the same signed request multiple times. When `assetAmount > 0` the EIP-3009 nonce would already block a second transfer, but the dedicated nonce provides defense-in-depth.

Because `facilitatorNonces` lives inside the `ProcessorEndpoint` contract and is keyed by user address, it does **not** interfere with the user's EOA nonce. A user can freely submit direct transactions (e.g., call `submitRequest`) without affecting their facilitator nonce, and vice versa.

**Sequential vs random.** Unlike EIP-3009 (which uses random `bytes32` nonces and an on-chain bitmap), this design uses a simple sequential counter. Trade-offs:

| | Sequential (current design) | Random (EIP-3009 style) |
|---|---|---|
| **Pro** | Minimal storage, simple replay prevention, enforces ordering | Allows out-of-order / parallel submission |
| **Con** | Requires querying the current nonce on-chain; a stuck tx blocks subsequent requests for the same user | Higher gas cost (bitmap storage), more complex contract logic |

Sequential nonces are preferred here because:
- They prevent replay attacks with minimal storage
- They enforce ordering (a user can't submit request N+1 before N is mined)
- They're simpler for the facilitator service to track
- The facilitator is a single centralized service with moderate throughput, so the ordering constraint is acceptable

#### Modified `generateRequestId`

The existing `generateRequestId` includes `msg.sender`. For `submitRequestFor`, the user address (recovered from signature) must be used instead. Per the ERC-20 design, the formula now includes `tokenAddress` and `assetAmount`:

```solidity
// Internal version that accepts explicit sender
function _generateRequestId(
    address sender,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 idx
) internal pure returns (bytes32) {
    return keccak256(abi.encodePacked(
        sender, applicationId, requestType, payload, tokenAddress, assetAmount, idx
    ));
}
```

#### Refund / claim routing

When a facilitated request completes or fails, claims must be split between user and facilitator. This requires storing the facilitator address in the pending request:

```solidity
struct PendingRequest {
    // ... existing fields from ERC-20 design ...
    address sender;          // user (real requester)
    address facilitator;     // address(0) for direct submissions, facilitator address for meta-tx
}
```

In `stateUpdate`, claim routing becomes:

```solidity
if (request.facilitator != address(0)) {
    // Facilitated request:
    //   - business-asset refund → user (in tokenAddress)
    //   - fee refund → facilitator (in ETH)
    if (assetRefund > 0) {
        pendingClaims[request.tokenAddress][request.sender] += assetRefund;
        totalPendingClaims[request.tokenAddress] += assetRefund;
    }
    if (feeRefund > 0) {
        pendingClaims[address(0)][request.facilitator] += feeRefund;
        totalPendingClaims[address(0)] += feeRefund;
    }
} else {
    // Direct request: standard claim routing per ERC-20 design
    //   - business-asset refund → sender (in tokenAddress)
    //   - fee refund → sender (in ETH)
    if (assetRefund > 0) {
        pendingClaims[request.tokenAddress][request.sender] += assetRefund;
        totalPendingClaims[request.tokenAddress] += assetRefund;
    }
    if (feeRefund > 0) {
        pendingClaims[address(0)][request.sender] += feeRefund;
        totalPendingClaims[address(0)] += feeRefund;
    }
}
```

The `claim(tokenAddress, payee)` function remains unchanged — it is permissionless and works for both users and facilitators.

### 5.3. Facilitator Service (Off-Chain)

The facilitator service will be implemented as a separate project: [`vela-x402`](https://github.com/HorizenOfficial/vela-x402). It is **not** part of this repository.

`vela-x402` will implement a custom **x402 payment scheme** (`private-vela-fixed`) that integrates with the HTTP 402 payment protocol. The facilitator acts as the settlement layer between x402 resource servers and the `ProcessorEndpoint` contract.

#### Architecture

The `vela-x402` project is a TypeScript monorepo with the following structure:

- **`packages/x402-private-vela-fixed/`** — Core npm package (`@horizen/x402-private-vela-fixed`) containing:
  - `types.ts` — Data structures for `RequestAuthorization`, `DepositAuthorization`, `VelaPaymentPayload`
  - `client/` — Client-side scheme (`VelaPaymentClientScheme`) for constructing EIP-712 + EIP-3009 signatures
  - `server/` — Facilitator-side scheme (`VelaPaymentFacilitatorScheme`) for verification and on-chain settlement
- **`facilitator/`** — Deployable Express.js server exposing x402 endpoints and setup endpoints
- **`specs/`** — Formal specification of the `private-vela-fixed` scheme

#### x402 Flow

The facilitator integrates with the x402 protocol as follows:

```
User/Client              Resource Server           Facilitator (vela-x402)     Blockchain
  │                            │                          │                        │
  │  1. HTTP request           │                          │                        │
  │───────────────────────────>│                          │                        │
  │                            │                          │                        │
  │  2. 402 Payment Required   │                          │                        │
  │     (scheme, token, amount)│                          │                        │
  │<───────────────────────────│                          │                        │
  │                            │                          │                        │
  │  3. Sign EIP-712 request   │                          │                        │
  │     + EIP-3009 deposit     │                          │                        │
  │     (off-chain)            │                          │                        │
  │                            │                          │                        │
  │  4. Retry with             │                          │                        │
  │     PAYMENT-SIGNATURE      │                          │                        │
  │───────────────────────────>│                          │                        │
  │                            │  5. POST /verify         │                        │
  │                            │  (off-chain validation)  │                        │
  │                            │─────────────────────────>│                        │
  │                            │                          │                        │
  │                            │  6. POST /settle         │                        │
  │                            │─────────────────────────>│                        │
  │                            │                          │  submitRequestFor()    │
  │                            │                          │  msg.value=maxFeeValue │
  │                            │                          │───────────────────────>│
  │                            │                          │                        │
  │                            │                          │  txHash, requestId     │
  │                            │                          │<───────────────────────│
  │                            │                          │                        │
  │  7. 200 OK + resource      │                          │                        │
  │<───────────────────────────│                          │                        │
```

#### Setup Endpoints (Non-x402)

In addition to the standard x402 verify/settle flow, the facilitator exposes additional endpoints to perform the initial steps (register, deposit) and
the final withdraw step.

- **`POST /register`** — Relays `ASSOCIATEKEY` requests (user key registration, one-shot)
- **`POST /deposit`** — Relays ERC-20 deposit requests via `submitRequestFor()`
- **`POST /withdrawal`** — Relays ERC-20 withdrawal requests via `submitRequestFor()`

These enable users to execute also the onboarding (register keys and fund their private balance) and withdrawal without holding ETH.

#### Facilitator Responsibilities

1. **Validate** both signatures off-chain before submitting (avoid wasting gas on invalid requests)
2. **Manage** its own ETH balance for gas and `maxFeeValue` payments (sent as `msg.value`)
3. **Set** `maxFeeValue` according to platform policy (the user doesn't choose this — the facilitator does)
4. **Monitor** its ETH balance, alert on low funds
5. **Claim** its own ETH refunds periodically via `claim(address(0), facilitatorAddress)`

#### Facilitator Wallet

The facilitator needs an EOA with:
- ETH for gas + `maxFeeValue` (both paid in the same transaction as `msg.value`)

No ERC-20 token balance or approval is needed for the facilitator itself — it only pays in ETH. The user's ERC-20 deposit is pulled directly from the user via EIP-3009.

### 5.4. Direct EOA Path (Unchanged)

The existing `submitRequest` function remains available. After the ERC-20 migration (prerequisite), it accepts both ETH and ERC-20 deposits per the ERC-20 design:
- ETH deposits: `msg.value = assetAmount + maxFeeValue`
- ERC-20 deposits: `msg.value = maxFeeValue`, asset pulled via `transferFrom` (requires prior `approve`)

Users who hold ETH can continue using `submitRequest` directly.

### 5.5. DEANONYMIZATION and ASSOCIATEKEY Considerations

#### ASSOCIATEKEY

Works naturally with the facilitator: the user signs a request authorization with their key payload (133 bytes). The recovered address from the EIP-712 signature is the user whose key is being associated. No special handling needed.

#### DEANONYMIZATION

Deanonymization will not be covered by the facilitator feature.

## 6. Security Considerations

### 6.1. Replay Protection

- **Request replay**: Sequential nonce per user prevents replaying the same request authorization. The nonce is consumed atomically in `submitRequestFor`.
- **Cross-chain replay**: The EIP-712 domain separator includes `chainId` and contract address.
- **Deposit replay**: EIP-3009 nonces are managed by the token contract independently. Each `transferWithAuthorization` nonce can only be used once.

### 6.2. Signature Binding

The request authorization signs `payloadHash` (not the raw payload) to keep signature size manageable. The contract verifies `keccak256(payload) == payloadHash` to ensure the facilitator cannot alter the payload.

The request authorization also binds `tokenAddress` and `assetAmount`, preventing the facilitator from substituting a different token or amount.

### 6.3. Deadline Enforcement

Both the request authorization and the EIP-3009 transfer have deadlines. The facilitator cannot hold signatures indefinitely and submit them later when conditions change.

### 6.4. Facilitator Trust Model

The facilitator is **semi-trusted**:
- **Cannot** alter request parameters (bound by user's EIP-712 signature)
- **Cannot** steal the user's deposit (EIP-3009 transfer goes to the contract, not the facilitator)
- **Cannot** change the deposit token or amount (bound by both EIP-712 and EIP-3009 signatures)
- **Can** choose not to submit a valid request (censorship) — mitigated by the direct EOA path remaining available
- **Can** choose the `maxFeeValue` — but this is the facilitator's own ETH, so misalignment is self-penalizing
- **Can** front-run or delay submission within the deadline window

### 6.5. Nonce Synchronization

The `facilitatorNonces` mapping is independent from direct `submitRequest` calls. Direct calls do not consume a facilitator nonce. This avoids coupling between the two submission paths.

If a user uses both paths, there is no conflict: direct calls are not nonce-gated, and facilitated calls use their own sequential counter.

## 7. Changes Summary

### Smart Contracts

| File | Change |
|---|---|
| `ProcessorEndpoint.sol` | Add `submitRequestFor()`, `facilitatorNonces` mapping, EIP-712 domain/typehash constants, refactor `generateRequestId` to internal `_generateRequestId` |
| `IProcessorEndpoint.sol` | Add `submitRequestFor()` and `getFacilitatorNonce()` to interface |
| `Structs.sol` | Add `facilitator` field to `PendingRequest` |
| `ProcessorEndpoint.sol` | Modify `stateUpdate` claim routing for split claims (user asset / facilitator ETH fee) |

### Go Backend

| Component | Change |
|---|---|
| Contract bindings | Regenerate via `go generate` after contract changes |
| `pkg/common/types.go` | Add `Facilitator` field to `Request` struct |
| `pkg/blockchain/client.go` | Parse `facilitator` from `PendingRequest` |
| Subgraph schema | Add `facilitator` field to `RequestSubmitted` entity |
| Subgraph handler | Index `facilitator` from event/call data |

### External: Facilitator Service (`vela-x402`)

The off-chain facilitator service and client SDK are implemented in the separate [`vela-x402`](https://github.com/HorizenOfficial/vela-x402) repository. No new off-chain components are added to this project.

| Component | Repository | Description |
|---|---|---|
| Facilitator HTTP service | `vela-x402/facilitator/` | Express.js server: x402 verify/settle + setup endpoints (register, deposit) |
| Client SDK | `vela-x402/packages/x402-private-vela-fixed/` | npm package with client-side and server-side x402 scheme implementations |

## 8. Sequence Diagram — Full x402 Flow

```
┌──────┐           ┌─────────────────┐       ┌───────────────────┐       ┌──────────────────────┐
│ User │           │ Resource Server │       │ Facilitator       │       │ ProcessorEndpoint    │
│      │           │                 │       │ (vela-x402)       │       │                      │
└──┬───┘           └────────┬────────┘       └─────────┬─────────┘       └──────────┬───────────┘
   │                        │                          │                            │
   │  1. HTTP request       │                          │                            │
   │───────────────────────>│                          │                            │
   │                        │                          │                            │
   │  2. 402 Payment        │                          │                            │
   │     Required           │                          │                            │
   │     {scheme, token,    │                          │                            │
   │      amount, resource} │                          │                            │
   │<───────────────────────│                          │                            │
   │                        │                          │                            │
   │  3. Query nonce        │                          │  getFacilitatorNonce(user) │
   │  ──────────────────────────────────────────────>  │  ─────────────────────────>│
   │       nonce = N        │                          │                   nonce = N│
   │  <──────────────────────────────────────────────  │  <─────────────────────────│
   │                        │                          │                            │
   │  4. Sign EIP-712       │                          │                            │
   │     request auth +     │                          │                            │
   │     EIP-3009 deposit   │                          │                            │
   │     (off-chain)        │                          │                            │
   │                        │                          │                            │
   │  5. Retry with         │                          │                            │
   │     PAYMENT-SIGNATURE  │                          │                            │
   │───────────────────────>│                          │                            │
   │                        │  6. POST /verify         │                            │
   │                        │  (off-chain validation)  │                            │
   │                        │─────────────────────────>│                            │
   │                        │  OK                      │                            │
   │                        │<─────────────────────────│                            │
   │                        │                          │                            │
   │                        │  7. POST /settle         │                            │
   │                        │─────────────────────────>│                            │
   │                        │                          │  8. submitRequestFor(      │
   │                        │                          │       params, sigs)        │
   │                        │                          │     msg.value=maxFeeValue  │
   │                        │                          │  ─────────────────────────>│
   │                        │                          │                            │
   │                        │                          │    ┌───────────────────────│
   │                        │                          │    │ 9.  Verify EIP-712    │
   │                        │                          │    │     → recover user    │
   │                        │                          │    │ 10. Verify & consume  │
   │                        │                          │    │     nonce             │
   │                        │                          │    │ 11. Validate token    │
   │                        │                          │    │     allowlists        │
   │                        │                          │    │ 12. EIP-3009 transfer │
   │                        │                          │    │     user → contract   │
   │                        │                          │    │ 13. Create            │
   │                        │                          │    │     PendingRequest    │
   │                        │                          │    │     sender=user       │
   │                        │                          │    │     facilitator=      │
   │                        │                          │    │       msg.sender      │
   │                        │                          │    │ 14. Emit              │
   │                        │                          │    │     RequestSubmitted  │
   │                        │                          │    └───────────────────────│
   │                        │                          │                            │
   │                        │                          │  requestId, txHash         │
   │                        │                          │  <─────────────────────────│
   │                        │  {requestId, txHash}     │                            │
   │                        │<─────────────────────────│                            │
   │                        │                          │                            │
   │  8. 200 OK + resource  │                          │                            │
   │<───────────────────────│                          │                            │
   │                        │                          │                            │
   │  ... request processed by manager/executor ...                                │
   │                        │                          │                            │
   │  15. claim(tokenAddr, user)                       │                            │
   │  (asset refund, if any)                           │                            │
   │───────────────────────────────────────────────────────────────────────────────>│
   │                        │                          │                            │
   │                        │                          │  16. claim(0x0,facilitator)│
   │                        │                          │  (ETH fee refund, if any)  │
   │                        │                          │  ─────────────────────────>│
```

## 9. Open Questions

1. **Rate limiting / abuse prevention**: deferred for now, but should be designed before production. A facilitator without rate limiting is an open gas faucet.

2. **Multi-facilitator support**: can multiple facilitator addresses be active simultaneously? The current design supports this naturally (any address can call `submitRequestFor`), but operational monitoring may need to account for it.
