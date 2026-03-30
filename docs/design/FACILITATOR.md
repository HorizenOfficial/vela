# Facilitator: Gasless Request Submission

## 1. Prerequisites

This design builds on the **ERC-20 Deposits/Withdrawals/Claims** change ([design doc](https://github.com/HorizenOfficial/vela/blob/9ba796c236867a1c27277c3cb1a8851d6f3c4212/ERC20_DEPOSITS_WITHDRAWALS_DESIGN.md)), which must be implemented first. After that change:

- `submitRequest` accepts a `tokenAddress` + `assetAmount` pair for the business-asset deposit (`0x0` = ETH, non-zero = ERC-20)
- For ERC-20 deposits, the contract uses `approve + transferFrom` to pull tokens from the caller
- `maxFeeValue` **remains ETH-only**, paid via `msg.value`
- Claims replace pull payments: `claim(tokenAddress, payee)` is the generic, permissionless, asset-aware claim mechanism
- Pending claims are tracked by `(payee, tokenAddress)`
- Token allowlists (global) gate which tokens are accepted

Additional prerequisite:

- **Token requirement**: only ERC-20 tokens that implement **EIP-2612** (`permit`) are supported for facilitated deposits.
Primary target: **USDC**.

## 2. Problem Statement

To submit a request via `ProcessorEndpoint.submitRequest`, a user must hold ETH to cover **gas fees** and **service fees** (`maxFeeValue`). 
This creates onboarding friction: users need to acquire ETH before they can interact with the platform, even when they have sufficient ERC-20 tokens for their application deposit.

A **facilitator** should allow users to submit requests without holding any ETH. The facilitator absorbs gas and service fees as an operational cost, while the user authorizes only their business-asset deposit via an off-chain signature.

## 3. Goals

- Users can submit requests **without holding ETH**
- The facilitator covers **gas fees** (paid to Ethereum validators) and **service fees** (`maxFeeValue`, always in ETH)
- The user pays **only `assetAmount`** in ERC-20 tokens, authorized via off-chain EIP-2612 signature
- The **real user address** is tracked as the request sender (for key association, deposits in WASM, claims)
- The existing **direct EOA submission** path remains fully functional
- The facilitator is a platform-operated service that **absorbs costs** (no on-chain reimbursement mechanism)

The facilitator pattern is designed for **ERC-20 business-asset deposits**. For ETH deposits (`tokenAddress = 0x0`), the user inherently needs ETH for `assetAmount` and therefore likely has enough for gas and fees — making the facilitator less valuable. ETH deposits remain supported only via the direct `submitRequest` path.

**Supported request types:** Initially, only `ASSOCIATEKEY` and `PROCESS` request types will be supported via `submitRequestFor`. Any other request type submitted through the facilitator will be rejected with an error. Support for additional request types may be added in future iterations.


This document covers the **smart contract changes** and **Go backend changes** needed in this repository to support facilitated requests. The off-chain facilitator service and client SDK will be implemented in a separate [`vela-facilitator`](https://github.com/HorizenOfficial/vela-facilitator) repository, which will provide also compatibility with the x402 HTTP payment protocol defined by Coinbase (see [https://github.com/coinbase/x402/](https://github.com/coinbase/x402/)).
More info on the off-chain facilitator will be provided below in paragraph 5.3 .

## 4. Approach: EIP-2612 + EIP-712 Meta-Transaction

### 4.1. Why EIP-2612

EIP-2612 defines `permit` — a standard that allows a token holder to authorize a spender via off-chain signature. A third party (the facilitator) can then call `permit` followed by `transferFrom` on-chain without the token holder spending gas.

This fits our model: the user signs a permit authorizing the contract to pull `assetAmount` tokens, and the facilitator executes both `permit` and `transferFrom` in the same transaction.

### 4.2. Why not ERC-4337

ERC-4337 (Account Abstraction) is more powerful but introduces significant infrastructure overhead:
- Requires deploying smart contract wallets for users
- Requires a Bundler service and EntryPoint contract
- Paymaster contracts add another trust/upgrade surface
- Overkill when the only gasless operation is `submitRequest`

ERC-4337 could be revisited if the platform needs broader account abstraction (e.g., gasless contract wallet interactions beyond Vela).

### 4.3. Why not EIP-3009

EIP-3009 (`transferWithAuthorization`) is an alternative that allows direct transfer via off-chain signature with random `bytes32` nonces. While it provides stronger isolation between protocols (random nonces can't collide), EIP-2612 is more widely adopted across ERC-20 tokens and provides sufficient security for our use case. The sequential nonce model of EIP-2612 is simpler and well understood. 

### 4.4. Why not EIP-2771

EIP-2771 (Trusted Forwarder) solves `msg.sender` identity but doesn't solve payment delegation for ERC-20 deposits. Since we already need EIP-2612 for the deposit token transfer, we can derive user identity from the EIP-712 request signature directly, without introducing a forwarder dependency.

## 5. Design

### 5.1. User Signature Scheme

The user produces **two signatures** off-chain:

#### Signature 1 — Request Authorization (EIP-712)

An EIP-712 typed data signature over the request parameters, proving the user authorized this specific request:

```
RequestAuthorization {
    address  sender             // user address (verified against ecrecover result)
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
Note there is a specific nonce for this authorization, described deeply below in 5.2.

The contract recovers the user address from this signature and verifies it matches the `sender` field. This prevents `ecrecover` from silently returning a wrong address when given a tampered message — without this check, a manipulated payload could pass nonce validation against an arbitrary unused address (especially when `assetAmount == 0` and there is no EIP-2612 permit to provide a second identity check).

#### Signature 2 — Deposit Permit (EIP-2612)

A standard EIP-2612 `permit` signature:

```
Permit {
    address owner              // user
    address spender            // ProcessorEndpoint contract
    uint256 value              // assetAmount
    uint256 nonce              // sequential nonce managed by the token contract
    uint256 deadline           // expiration timestamp (matches request)
}
```

This authorizes the contract to call `transferFrom` and pull `assetAmount` tokens from the user.

> **When `assetAmount == 0`**: Signature 2 is not required. Many request types (e.g., `PROCESS`, `ASSOCIATEKEY`) may not require a deposit, in which case the user only signs the request authorization.

### 5.2. Contract Changes

#### New function: `submitRequestFor`

```solidity
function submitRequestFor(
    address sender,
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    // User's EIP-712 request authorization
    uint256 deadline,
    bytes calldata requestSignature,
    // EIP-2612 deposit permit (empty if assetAmount == 0)
    // When present: abi.encode(uint8 v, bytes32 r, bytes32 s)
    // owner, spender, value, deadline are derived from other parameters
    bytes calldata depositPermit
)
    external
    payable                    // msg.value = maxFeeValue (ETH)
    nonReentrant
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    returns (bytes32)
```

**Execution flow (solidity pseudo-code):**

```
submitRequestFor()
  │
  ├─ 1. Verify request type is supported
  │     require(requestType == ASSOCIATEKEY || requestType == PROCESS)
  │
  ├─ 2. Verify deadline not expired
  │     require(block.timestamp <= deadline)
  │
  ├─ 3. Read current nonce and build EIP-712 hash
  │     nonce = facilitatorNonces[sender]
  │     // nonce is not a parameter — the contract uses the current on-chain value
  │     // to reconstruct the EIP-712 struct hash. If the user signed a different
  │     // nonce, ecrecover will return a wrong address and fail the sender check.
  │
  ├─ 4. Recover user address from EIP-712 request signature and verify
  │     user = ecrecover(hashTypedData(requestAuth), requestSignature)
  │     require(user != address(0))   // invalid signature guard
  │     require(user == sender)       // ecrecover result must match declared sender
  │
  ├─ 5. Consume nonce (replay protection)
  │     facilitatorNonces[sender]++
  │
  ├─ 6. If assetAmount > 0: validate token allowlists
  │     require(globalAllowedTokens[tokenAddress])
  |     Token address must not be the one identifying ETH
  │
  ├─ 7. Validate fee
  │     maxFeeValue = msg.value
  │     require(maxFeeValue >= minFeePerRequest)
  │
  ├─ 8. If assetAmount > 0: decode depositPermit and execute EIP-2612 permit + transferFrom
  │     (v, r, s) = abi.decode(depositPermit, (uint8, bytes32, bytes32))
  │     token.permit(user, address(this), assetAmount, deadline, v, r, s)
  │     token.transferFrom(user, address(this), assetAmount)
  │     appCustody[applicationId][tokenAddress] += assetAmount
  │
  ├─ 9. Create PendingRequest with sender = user (not msg.sender)
  │     pendingRequest.sender = user
  │     pendingRequest.facilitator = msg.sender
  │     pendingRequest.tokenAddress = tokenAddress
  │     pendingRequest.assetAmount = assetAmount
  │     pendingRequest.maxFeeValue = maxFeeValue
  │
  ├─ 10. Generate requestId using user address
  │      requestId = _generateRequestId(user, applicationId, requestType,
  │                                      payload, tokenAddress, assetAmount, idx)
  │
  └─ 11. Emit RequestSubmitted(applicationId, requestId, user)
```

#### Nonce management

```solidity
// Sequential nonces per user (consistent with EIP-2612's sequential model)
mapping(address => uint256) public facilitatorNonces;

function getFacilitatorNonce(address user) external view returns (uint256);
```

**Why a dedicated nonce is needed.** The `facilitatorNonces` counter is independent from both the Ethereum account nonce (used by the facilitator's EOA to order on-chain transactions) and the EIP-2612 nonce (managed by the token contract for `permit`). It protects the EIP-712 request authorization (Signature 1) against replay. This is critical when `assetAmount == 0` (e.g., `ASSOCIATEKEY` or `PROCESS` requests): in that case there is no Signature 2 / EIP-2612 permit, so without this nonce there would be **no replay protection at all** — a facilitator could re-submit the same signed request multiple times. When `assetAmount > 0` the EIP-2612 nonce would already block a second permit, but the dedicated nonce provides defense-in-depth.

**Nonce is not a calldata parameter.** The nonce is not passed as a parameter to `submitRequestFor`. The contract reads `facilitatorNonces[sender]` directly and uses it to reconstruct the EIP-712 struct hash. The user must sign the current on-chain nonce value; if the signed nonce differs, the recovered address won't match `sender` and the transaction reverts. This eliminates a parameter, reduces calldata, and removes the possibility of the facilitator submitting a mismatched nonce.

Because `facilitatorNonces` lives inside the `ProcessorEndpoint` contract and is keyed by user address, it does **not** interfere with the user's EOA nonce. A user can freely submit direct transactions (e.g., call `submitRequest`) without affecting their facilitator nonce, and vice versa.

**Sequential vs random.** Like EIP-2612 (which also uses sequential nonces), this design uses a simple sequential counter. Trade-offs:

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
// Asset refund always goes to user (the real requester)
if (assetRefund > 0) {
    pendingClaims[request.tokenAddress][request.sender] += assetRefund;
    totalPendingClaims[request.tokenAddress] += assetRefund;
}

// Fee refund goes to facilitator if present, otherwise to sender
address feeRecipient = request.facilitator != address(0)
    ? request.facilitator
    : request.sender;
if (feeRefund > 0) {
    pendingClaims[address(0)][feeRecipient] += feeRefund;
    totalPendingClaims[address(0)] += feeRefund;
}
```

The `claim(tokenAddress, payee)` function remains unchanged — it is permissionless and works for both users and facilitators.

### 5.3. Facilitator Service (Off-Chain)

The facilitator service will be implemented as a separate project: [`vela-facilitator`](https://github.com/HorizenOfficial/vela-facilitator). It is **not** part of this repository.

`vela-facilitator` is a platform-level service that submits requests on behalf of users who do not hold ETH. It exposes two layers of functionality:

1. **Core facilitation endpoints** — generic, application-agnostic gasless submission via `submitRequestFor()`
2. **x402 scheme endpoints** — standard Coinbase x402 `verify`/`settle` flow for applications that integrate with the [x402 HTTP payment protocol](https://github.com/coinbase/x402/)

> **Note on EIP-2612 and x402:** The standard x402 `exact` EVM scheme uses EIP-3009 as its recommended `assetTransferMethod` for USDC. Our custom `private-vela-fixed` scheme uses EIP-2612 instead. This is fully compatible with the x402 protocol architecture (which is scheme-agnostic), but means the Coinbase reference facilitator's `exact` settle logic cannot be reused directly — our scheme implements its own verify/settle via `submitRequestFor()`.

#### Architecture

```
vela-facilitator/
  │
  ├─ Core facilitation layer (generic, platform-level)
  │   ├─ POST /submit          — gasless submitRequestFor() for any app/request type
  │   ├─ GET  /nonce/:user     — query facilitatorNonces from contract
  │   └─ wallet management, gas monitoring, claim scheduling
  │
  ├─ x402 scheme layer (for x402-enabled applications)
  │   ├─ POST /verify          — off-chain signature validation (standard x402)
  │   └─ POST /settle          — on-chain settlement (delegates to core /submit)
  │
  └─ @horizen/x402-private-vela-fixed (npm package)
      └─ scheme definition for x402 clients and facilitators (incl. Coinbase reference facilitator)
```

The core layer handles all interaction with the `ProcessorEndpoint` contract: signature validation, nonce management, gas estimation, and transaction submission. The x402 layer is a thin adapter that translates the x402 verify/settle protocol into calls to the core layer.

Applications that do not use x402 (e.g., a mobile SDK, a bot, a CLI) can call the core endpoints directly.

#### Core Facilitation Endpoints

- **`POST /submit`** — Generic gasless submission. Accepts the user's EIP-712 request signature, EIP-2612 deposit permit (if `assetAmount > 0`), and request parameters. Calls `submitRequestFor()` on-chain. Works for any request type (`ASSOCIATEKEY`, `DEPOSIT`, `WITHDRAWAL`, `PROCESS`, etc.).
- **`GET /nonce/:user`** — Returns the current `facilitatorNonces[user]` value from the contract, so clients can construct valid signatures.

These endpoints are **open** (no authentication required). See Open Questions for future authentication/rate-limiting considerations.

#### x402 Scheme Layer

For applications integrating with the [x402 HTTP payment protocol](https://github.com/coinbase/x402/), the facilitator implements a custom scheme (`private-vela-fixed`) and the default endpoints defined by Coinbase: they are part of the specification from Coinbase, and defines a standard way to perform x402 payments and extend them by defining a custom payment flow.

IMPORTANT: this part is an extension/specialization of the Core Facilitation Endpoints that requires the usage of the vela-nova app for executing private transfers.

**npm package: `@horizen/x402-private-vela-fixed`** — This package defines the `private-vela-fixed` payment scheme and can be used by:
- **x402 clients** (via the Coinbase x402 client SDK) to construct EIP-712 + EIP-2612 payment signatures
- **x402 facilitators** (including the Coinbase reference facilitator or the `vela-facilitator` itself) to verify and settle payments

**Endpoints:**

- **`POST /verify`** — Validates both user signatures (request authorization + deposit permit) off-chain without submitting a transaction. Returns whether the payment would succeed.
- **`POST /settle`** — Validates signatures and submits the transaction on-chain via the core facilitation layer. Returns `requestId` and `txHash`.

The x402 flow follows the standard protocol: a resource server returns HTTP 402 with payment requirements, the client signs and retries, and the resource server delegates verification and settlement to the facilitator.

#### Facilitator Responsibilities

1. **Validate** both signatures (request authorization + deposit permit) off-chain before submitting (avoid wasting gas on invalid requests)
2. **Manage** its own ETH balance for gas and `maxFeeValue` payments (sent as `msg.value`)
3. **Set** `maxFeeValue` according to platform policy (the user doesn't choose this — the facilitator does)
4. **Monitor** its ETH balance, alert on low funds
5. **Claim** its own ETH refunds periodically via `claim(address(0), facilitatorAddress)`

#### Facilitator Wallet

The facilitator needs an EOA with:
- ETH for gas + `maxFeeValue` (both paid in the same transaction as `msg.value`)

No ERC-20 token balance or approval is needed for the facilitator itself — it only pays in ETH. The user's ERC-20 deposit is pulled from the user via EIP-2612 permit + transferFrom.

#### Sequence Diagram — x402 Flow

```
┌──────┐           ┌─────────────────┐       ┌───────────────────┐       ┌──────────────────────┐
│ User │           │ Resource Server │       │ vela-facilitator  │       │ ProcessorEndpoint    │
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
   │     EIP-2612 permit    │                          │                            │
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
   │                        │                          │    │ 12. EIP-2612 permit + │
   │                        │                          │    │     transferFrom      │
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
   │  15. 200 OK + resource │                          │                            │
   │<───────────────────────│                          │                            │
   │                        │                          │                            │
   │  ... request processed by manager/executor ...                                │
   │                        │                          │                            │
   │  16. claim(tokenAddr, user)                       │                            │
   │  (asset refund, if any)                           │                            │
   │───────────────────────────────────────────────────────────────────────────────>│
   │                        │                          │                            │
   │                        │                          │  17. claim(0x0,facilitator)│
   │                        │                          │  (ETH fee refund, if any)  │
   │                        │                          │  ─────────────────────────>│
```

### 5.4. Direct EOA Path (Unchanged)

The existing `submitRequest` function remains available. After the ERC-20 migration (prerequisite), it accepts both ETH and ERC-20 deposits per the ERC-20 design:
- ETH deposits: `msg.value = assetAmount + maxFeeValue`
- ERC-20 deposits: `msg.value = maxFeeValue`, asset pulled via `transferFrom` (requires prior `approve`)

Users who hold ETH can continue using `submitRequest` directly.

### 5.5. DEANONYMIZATION and ASSOCIATEKEY Considerations

#### ASSOCIATEKEY

Works naturally with the facilitator: the user signs a request authorization with their key payload. The recovered address from the EIP-712 signature is the user whose key is being associated. No special handling needed.

#### DEANONYMIZATION

Deanonymization will not be covered by the facilitator feature.

## 6. Security Considerations

### 6.1. Replay Protection

- **Request replay**: Sequential nonce per user prevents replaying the same request authorization. The nonce is consumed atomically in `submitRequestFor`.
- **Cross-chain replay**: The EIP-712 domain separator includes `chainId` and contract address.
- **Deposit replay**: EIP-2612 nonces are managed by the token contract independently. Each `permit` nonce is sequential and can only be used once. Note: EIP-2612 nonces are shared across all protocols using `permit` on the same token — if a user signs permits for multiple dApps concurrently, nonce ordering must be coordinated.

### 6.2. Signature Binding

The signed request authorization includes `payloadHash` rather than the raw payload. This is a convenience choice — EIP-712 hashes dynamic `bytes` fields to `keccak256(bytes)` internally anyway, so the cryptographic result is identical. The full payload is passed as a parameter to `submitRequestFor` (it travels on-chain as calldata), and the contract computes `keccak256(payload)` internally to reconstruct the EIP-712 struct hash. If the facilitator alters the payload, the reconstructed hash won't match the one the user signed, and `ecrecover` will return a different address — failing the `user == sender` check.

The request authorization also binds `tokenAddress` and `assetAmount`, preventing the facilitator from substituting a different token or amount. The EIP-2612 permit additionally binds the spender (the contract) and the value, providing a second layer of protection.

### 6.3. Deadline Enforcement

Both the request authorization and the EIP-2612 permit have deadlines. The facilitator cannot hold signatures indefinitely and submit them later when conditions change.

### 6.4. Facilitator Trust Model

The facilitator is **semi-trusted**:
- **Cannot** alter request parameters (bound by user's EIP-712 signature)
- **Cannot** steal the user's deposit (EIP-2612 permit authorizes the contract as spender; `transferFrom` is called by the contract itself)
- **Cannot** change the deposit token or amount (bound by both EIP-712 and EIP-2612 signatures)
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
| `ProcessorEndpoint.sol` | Add `submitRequestFor()`, `facilitatorNonces` mapping, EIP-712 domain separator + `REQUEST_AUTHORIZATION_TYPEHASH`, EIP-2612 `permit` + `transferFrom` for deposit handling, refactor `generateRequestId` to internal `_generateRequestId`, modify `stateUpdate` claim routing for split claims (user asset / facilitator ETH fee) |
| `IProcessorEndpoint.sol` | Add `submitRequestFor()` and `getFacilitatorNonce()` to interface |
| `Structs.sol` | Add `facilitator` field to `PendingRequest` |

### Go Backend

| Component | Change |
|---|---|
| Contract bindings | Regenerate via `go generate` after contract changes |
| `pkg/common/types.go` | Add `Facilitator` field to `Request` struct |
| `pkg/blockchain/client.go` | Parse `facilitator` from `PendingRequest` |
| Subgraph schema | Add `facilitator` field to `RequestSubmitted` entity |
| Subgraph handler | Index `facilitator` from event/call data |

## 8. Open points for future developments

1. **Facilitator API authentication**: the facilitation endpoints (`/submit`, `/nonce`) are currently open (no authentication). Before production, an authentication mechanism (API keys, allowlists, or similar) should be designed to prevent unauthorized use of the facilitator's ETH funds.

2. **Multi-facilitator support**: can multiple facilitator addresses be active simultaneously? The current design supports this naturally (any address can call `submitRequestFor`), but operational monitoring may need to account for it.

