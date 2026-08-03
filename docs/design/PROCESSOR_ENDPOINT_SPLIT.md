# Splitting ProcessorEndpoint Under the EIP-170 Size Limit

## 1. Problem

`ProcessorEndpoint` is the largest contract in the system, and it is close to the **24,576-byte
EIP-170 limit** on deployed bytecode. Measured with the project's compiler settings (solc 0.8.30,
`viaIR: true`, `optimizer.runs: 0`, `evmVersion: paris`) it is **23,246 bytes** — 1,330 bytes of
headroom.

The deployment target is **Base**, which is EVM-equivalent, so the limit applies exactly. Two
things made the situation worse than it looked:

- The hardhat network was configured with `allowUnlimitedContractSize: true`, so the contract suite
  could not detect the problem at all. Only `go test ./pkg/blockchain/...` would fail, because
  geth's simulated backend enforces the limit unconditionally (`max code size exceeded`).
- The remaining headroom is smaller than a single sizeable feature. The batch-execution work
  (`BATCH_EXECUTION.md`) exceeded it: per-application queues plus round-robin selection pushed the
  contract to 26,353 bytes, and `batchStateUpdate()` had not even been written yet.

This document covers the **mechanism** for making room, applied and verified independently of any
feature work. It is deliberately a pure refactor: no behavioural change, no new functionality.

## 2. Levers ruled out first

Cheap options were measured and rejected before restructuring anything.

| Lever | Result |
|---|---|
| `optimizer.runs` | `0` is already the smallest of 0/1/50/200/1000/10000 |
| `viaIR: false` | Fails to compile: *stack too deep* |
| `metadata.bytecodeHash: 'none'` | −41 bytes |
| `evmVersion: cancun`/`prague` (PUSH0) | −945 bytes, but **not taken**: `paris` is kept for chain portability |
| `evmVersion: shanghai` | −910 bytes, same decision |
| `ReentrancyGuardTransient` (OZ 5.3, needs cancun) | *Larger* by 91 bytes |
| Revert strings → custom errors | Already done; the contract has no revert strings left |
| External linked library | −143 bytes overall. Public library functions must ABI-encode every argument at the call site, and this state is spread over many separate public mappings, so the boundary needs 8–9 arguments including dynamic arrays. A library only pays off with a single storage-pointer boundary, which would mean bundling all state into one struct plus ~16 hand-written getters to preserve the ABI |
| Inheritance split for dev-only code (lean production contract + `ProcessorEndpointResettable` for tests) | The subclass is still over the limit, so it only deploys where the limit is disabled. A real testnet would run the lean variant and lose `adminReset`, which is exactly where `PROCESSOR_ENDPOINT_ADMIN_RESET.md` says the feature is meant to be used |
| Facilitator path into a separately-**called** contract | ≈ −2,200, because every argument is ABI-encoded at the external-call boundary — and it changes the address relayers call. Dominated by the mechanism below |

Module costs, each measured by deleting the module:

| Module | Bytes |
|---|---|
| Queue views + selection + copy helpers | 2,833 |
| `submitRequestFor` + `getFacilitatorNonce` (EIP-712 + permit relaying) | 2,537 |
| `_invokeTrigger` + `_enqueueTrustedRequest` | 2,110 |
| `adminReset` + `adminResetApps` + helpers | 1,570 |

## 3. Mechanism: a delegatecall extension

`abstract contract ProcessorEndpointStorage` becomes the single declaration of every state
variable. It inherits `AccessControl, ReentrancyGuard, EIP712` **in that order**, because those
contribute storage slots (`_roles`, `_status`, the EIP-712 name/version fallbacks) that must sit at
the same offsets on both sides. `ProcessorEndpoint` and `ProcessorEndpointExtension` both derive
from it:

```
                  ┌─────────────────────────────┐
                  │ ProcessorEndpointStorage    │  all state + shared internal helpers
                  │ (AccessControl,             │
                  │  ReentrancyGuard, EIP712,   │
                  │  IProcessorEndpointState)   │
                  └──────────┬──────────────────┘
                             │
             ┌───────────────┴────────────────┐
             ▼                                ▼
  ┌────────────────────────┐      ┌──────────────────────────────┐
  │ ProcessorEndpoint      │      │ ProcessorEndpointExtension   │
  │ (deployed, holds       │      │ (deployed, holds no state)   │
  │  storage and funds)    │      │                              │
  │                        │      │  submitRequestFor  ◀──┐      │
  │  submitRequestFor ─────┼──────┼── delegatecall ───────┘      │
  └────────────────────────┘      └──────────────────────────────┘
```

Because the boundary is a `delegatecall` onto a shared layout, arguments are **not re-encoded** and
internal helpers are reused directly — which is precisely what made the linked-library and
separate-callee options unattractive. Helpers the extension needs get duplicated into *its*
bytecode, which costs nothing: only the endpoint is near the limit, and the extension has ~17KB
spare.

Everything currently hosted in the extension is the facilitator path (`submitRequestFor`). Further
modules can move the same way when more room is needed.

### 3.1. Per-function stub, not a `fallback()`

The endpoint keeps a typed `submitRequestFor` declaration whose body is `_delegateToExtension()`.
A generic `fallback()` router would have been ~150 bytes cheaper per moved function but would have
required a **merged ABI** (endpoint + extension) for `abigen`, typechain and the hardhat fixture,
since the moved function would no longer appear in the endpoint's own ABI.

With the stub instead:

- the ABI, the function selector, the address relayers call, and the `IProcessorEndpoint`
  conformance check are all unchanged;
- no tooling, test, binding or subgraph had to learn about the split.

The trade-off reverses once several functions move: at ~150 bytes of stub each, a generic
`fallback()` plus a merged ABI becomes the better deal. That is the expected next step if the queue
views (−2,833) have to move.

### 3.2. The forwarder must be memory-safe assembly

```solidity
function _delegateToExtension() private {
  address extension = _extension;
  assembly ('memory-safe') {
    let ptr := mload(0x40)
    calldatacopy(ptr, 0, calldatasize())
    ...
  }
}
```

The obvious version — writing at memory offset 0, as OpenZeppelin's `Proxy` does — **breaks the
build**. Unannotated assembly switches off solc's `memoryguard` for the entire contract, and
`stateUpdate` then fails to compile with *stack too deep*. Scratching above the free-memory pointer
and annotating `memory-safe` keeps the optimizer's assumptions intact.

### 3.3. State-variable getters had to move to their own interface

A public state variable can only implement an interface function when the **same** contract
declares both. Once the state moved into `ProcessorEndpointStorage`, the six getters that
`IProcessorEndpoint` declares and the state variables satisfy stopped resolving:

```
TypeError: Derived contract must override function "appCustody". Two or more base classes define
function with same name and parameter types. Since one of the bases defines a public state variable
which cannot be overridden, you have to change the inheritance layout or the names of the functions.
```

They are now declared in `IProcessorEndpointState`, which `ProcessorEndpointStorage` inherits — so
the variables satisfy them — and which `IProcessorEndpoint` also inherits, so external consumers
see one combined interface exactly as before (including `AbstractTrigger`'s
`processorEndpoint.tokenAllowlist()` call).

`ProcessorEndpointStorage` cannot simply inherit the full `IProcessorEndpoint`: that would force
`ProcessorEndpointExtension` to implement every endpoint function. For the same reason, errors and
events are referenced from the base and the extension as `IProcessorEndpoint.<name>`.

### 3.4. Other consequences

- **`getFacilitatorNonce()` was removed.** It duplicated the auto-generated getter of the
  `public facilitatorNonces` mapping. `FACILITATOR.md` already told clients to read the mapping
  directly. Callers use `facilitatorNonces(user)`. This is the only ABI removal, and it is recorded
  in `CHANGELOG.md`.
- **The extension refuses direct calls** (`DirectCallNotAllowed`, via an `immutable _self` compared
  against `address(this)`). Called directly it would read and write its own empty storage and could
  strand the ETH sent as a fee.
- **The endpoint rejects an extension without code** (`InvalidExtension`), on top of the
  zero-address check. A `delegatecall` to a codeless address *succeeds* and returns nothing, so a
  wrong address would make `submitRequestFor` a silent no-op that keeps the fee — no event, no
  queue entry — rather than reverting. `_extension` is immutable, so the constructor is the only
  place to catch it. Constructor-only code, so it costs nothing against EIP-170.
- **Revert data must stay decodable through the endpoint's ABI.** The extension's code reverts
  from the endpoint's address, but only the endpoint's *ABI* is what clients hold. Signature
  recovery therefore uses `ECDSA.tryRecover` and raises `IProcessorEndpoint.InvalidSignature`;
  `ECDSA.recover` would raise `ECDSAInvalidSignature`/`-Length`/`-S`, which are absent from the
  endpoint's ABI now that the endpoint itself no longer calls ECDSA. Anything moved to the
  extension in future needs the same check: every error it can raise must be declared on the
  endpoint too.
- **The stub keeps its parameter names.** Parameter names are part of the published ABI — wallets,
  explorers and generated bindings display them, and `abigen` degrades to `arg0…arg9` without
  them. The names are silenced with no-op statements rather than commented out; verified to produce
  byte-identical output.
- **Deploy order changed.** The extension is deployed first and its address is a new last
  `ProcessorEndpoint` constructor argument. It is `immutable`, so it lives in the endpoint's code
  rather than storage and cannot be repointed. Wired in `test/ProcessorEndpoint/fixture.ts`,
  `scripts/deploy/all.ts`, `scripts/deploy/processorEndpoint.ts` and
  `pkg/blockchain/testutil/sim_test_helper.go`. The extension has its own `go:generate` chain and
  bindings package, needed only for deployment.
- **Runtime cost**: one extra cold-account access (~2,600 gas) per forwarded call, on the
  relayer-paid facilitator path. OpenZeppelin's `EIP712` recomputes the domain separator whenever
  `address(this)` differs from the deploying address, so `verifyingContract` remains the endpoint
  and existing signatures stay valid.
- **Compatible with `UPGRADABLE_CONTRACTS_DESIGN.md`**: proxy → implementation → extension are both
  `delegatecall`s, so storage stays in the proxy throughout, and an `immutable` in the
  implementation is still readable. Upgrading the extension means deploying a new implementation
  pointing at it — the normal UUPS path. The `__gap` that design requires belongs in the storage
  base.

## 4. Result

| Contract | Deployed bytes | vs limit |
|---|---|---|
| `ProcessorEndpoint` before | 23,246 | −1,330 |
| `ProcessorEndpoint` after | **21,074** | **−3,502** |
| `ProcessorEndpointExtension` | 6,655 | −17,921 |

The endpoint sheds 2,172 bytes: the 2,537-byte facilitator module, less the forwarding stub and the
two thin public wrappers (`getPendingRequestsSize`, `generateRequestId`) that now delegate to
internal implementations in the base.

**This is probably not enough for batch execution.** Section 1 measured the batch work at 26,353
bytes on the old 23,246-byte baseline — +3,107 — and `batchStateUpdate()` was not written yet. On
the new baseline that lands at 24,155, only 421 bytes under the limit. Expect to need the next
lever (section 3.1: the queue views, −2,833, behind a generic `fallback()` and a merged ABI) as
part of the batch work rather than after it.

Note that the `go:generate` solc invocation passes `--optimize` (default `runs: 200`) and no
`--evm-version`, so it produces post-Shanghai bytecode with different sizes from hardhat's
`paris`/`runs: 0` build. **Any size claim must say which path produced it**; the numbers in this
document are the hardhat path.

## 5. Guardrails

Both are verified to fail when they should, not just to pass.

### 5.1. EIP-170 is enforced by the contract suite

`hardhat.config.ts` now gates `allowUnlimitedContractSize` on an environment variable:

```ts
hardhat: {
  allowUnlimitedContractSize: process.env.UNLIMITED_SIZE === 'true',
},
```

So `npm run test` fails the way Base would, and `npm run test:ignoresize` (a script that already
existed but had no effect) is the deliberate opt-out for measuring how far over the limit a
work-in-progress contract is. Checked with a deliberately oversized probe contract: a 39,644-byte
contract is rejected by default and deploys with `UNLIMITED_SIZE=true`.

### 5.2. Storage layouts must stay identical

`npm run check:layout` (`scripts/checkStorageLayout.ts`) compares solc's `storageLayout` output for
`ProcessorEndpoint` and `ProcessorEndpointExtension`, and runs in both CI jobs. This is the check
that matters most: a divergence would **not revert** — the extension would silently read and write
the wrong slots.

Two implementation details:

- The solidity config requests `storageLayout` in `outputSelection` so the layout is available in
  hardhat's build-info.
- Type identifiers embed solc's internal AST ids (`t_contract(ITrigger)9681`), which are only stable
  within one compilation job, and hardhat may compile the two contracts in separate jobs. The
  script strips those ids so it compares type *shape*, not numbering — without this it reported
  spurious mismatches.

Verified by injecting an extra state variable into the extension: the check exits non-zero and
names the offending entry.

The comparison covers the flattened slot list — slot, offset, label, type — which is exactly the
surface on which the two contracts can diverge, since all state is declared once in the base and
the failure mode is a variable declared or reordered in a derived contract. It would not see a
change *inside* a struct's fields, which requires a derived contract to redeclare a struct of the
same name.

## 6. Reviewing this change

The refactor is behaviour-preserving, and the test suite is the evidence:

- **291 tests before, 296 after.** The five additions are the two constructor cases (`extension is
  zero address`, `extension has no code`), `DirectCallNotAllowed`, and the two malformed-signature
  cases that must surface as `InvalidSignature`.
- **No existing test's assertions changed.** Only five test files are touched, and every change is
  either the added constructor argument, the fixture deploying the extension, or one of the new
  tests.
- `npm run check:layout` passes, the full Go suite passes, and `go generate ./...` is idempotent.

One behavioural difference is deliberate and worth naming, because it is not covered by any test:
the endpoint now reverts `InvalidSignature` where it previously reverted `ECDSAInvalidSignature*`
on a malformed facilitator signature (see section 3.4). Nothing else changes. In particular
`_queueSize` keeps its `tail > head` guard: the guard is unreachable — `_queueEnqueue` only
advances `tail`, `_queueDequeueHead` is reached only after `isCurrentPendingRequest` has confirmed
the id is a queue head, and `adminReset` sets `tail = head` — but dropping it would have made two
external views (`getPendingRequestsSize`, `getTriggerQueueSize`) revert on underflow instead of
returning 0, which is a behavioural change this refactor does not need to take. It costs 26 bytes.

## 7. Deliberately deferred: per-application contracts

A router plus one endpoint contract per application (EIP-1167 clones) would remove the ceiling
permanently, and would turn per-app fund isolation from arithmetic bookkeeping into physical
separation. It is deferred because it also touches the subgraph (dynamic data sources for N
addresses), the Go client (router plus per-app address discovery, event subscriptions, reorg
detection), custody, and the cross-application invariants that do not decompose
(`pendingClaims`/`totalPendingClaims`, fees, TRUSTPROCESS global priority). It needs its own design
doc, sequenced with `UPGRADABLE_CONTRACTS_DESIGN.md`.

Nothing in this document blocks that: the storage base and extension would be absorbed by it rather
than fought against.
