# Splitting ProcessorEndpoint Under the EIP-170 Size Limit

## 1. Problem

`ProcessorEndpoint` is the largest contract in the system, and it is close to the **24,576-byte
EIP-170 limit** on deployed bytecode. Measured with the project's compiler settings as they were
(solc 0.8.30, `viaIR: true`, `optimizer.runs: 0`, `evmVersion: paris`) it is **23,246 bytes** —
1,330 bytes of headroom. `evmVersion` is now `cancun`, for the reasons in section 2; every baseline
in sections 1–3 predates that switch, and section 4 reports the result under both settings.

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
| `evmVersion: cancun`/`prague` (PUSH0) | −945 bytes, and **taken** — see below |
| `evmVersion: shanghai` | −910 bytes, superseded by `cancun` |
| `ReentrancyGuardTransient` (OZ 5.3, needs cancun) | *Larger* by 91 bytes, so still not used even though `cancun` now makes it available |
| Revert strings → custom errors | Already done; the contract has no revert strings left |
| External linked library | −143 bytes overall. Public library functions must ABI-encode every argument at the call site, and this state is spread over many separate public mappings, so the boundary needs 8–9 arguments including dynamic arrays. A library only pays off with a single storage-pointer boundary, which would mean bundling all state into one struct plus ~16 hand-written getters to preserve the ABI |
| Inheritance split for dev-only code (lean production contract + `ProcessorEndpointResettable` for tests) | The subclass is still over the limit, so it only deploys where the limit is disabled. A real testnet would run the lean variant and lose `adminReset`, which is exactly where `PROCESSOR_ENDPOINT_ADMIN_RESET.md` says the feature is meant to be used |
| Facilitator path into a separately-**called** contract | ≈ −2,200, because every argument is ABI-encoded at the external-call boundary — and it changes the address relayers call. Dominated by the mechanism below |

Module costs, each measured by deleting the module from the 23,246-byte pre-split `paris` build:

| Module | Bytes |
|---|---|
| Queue views + selection + copy helpers | 2,833 |
| `submitRequestFor` + `getFacilitatorNonce` (EIP-712 + permit relaying) | 2,537 |
| `_invokeTrigger` + `_enqueueTrustedRequest` | 2,110 |
| `adminReset` + `adminResetApps` + helpers | 1,570 |

These are the numbers that drove the original selection. They do not carry to the current `cancun`
baseline — at `runs: 0` the optimizer is not linear in this way — so section 3 re-measures every
module before moving it, and the queue-view figure in particular does not survive (section 3.1).

### 2.1. `evmVersion: cancun` is taken

`paris` was initially kept "for chain portability", which does not survive contact with the actual
requirement: the deployment target named in section 1 is **Base**, which has been on Cancun since
March 2024, and no chain this system is planned for predates it. The cost of the caution was ~720
bytes on the post-split contract — more than the headroom left over after batch execution, for a
one-line config change with no ABI or tooling consequences. `hardhat.config.ts` now sets
`evmVersion: 'cancun'`, and the constraint that comes with it is recorded there and in
`CHANGELOG.md`: **any chain this is deployed to must be at Cancun or later.**

The `go:generate` solc path is unaffected: it passes no `--evm-version`, so it uses solc 0.8.30's
own default, which is already newer than `cancun`. That path only produces bindings for tests and
deployment, never the deployed artifact.

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

The extension hosts every entry point that is off the per-request hot path:

| Module | Endpoint bytes saved |
|---|---|
| `submitRequestFor` (facilitator path, EIP-712 + permit relaying) | 2,537 |
| `adminReset` + `adminResetApps` (+ `_resetQueue`, `_removeDeployedAppId`) | 1,695 |
| `submitDeployRequest` + `…WithTrigger` (+ `_submitDeployRequest`) | 1,028 |
| `updateQueueThreshold`, `updateMaxNumOfApplications`, `updateFeeCollector`, `add`/`removeAllowedDeployer` | 781 |

Each figure is net of its forwarding stub, measured by compiling the move on the baseline that
preceded it. The selection rule is gas, not size: forwarding costs one extra cold-account access
(~2,600 gas), which is noise on a rare or privileged call and unacceptable on `submitRequest` or
`stateUpdate`. `_asyncTransfer`, `_subtractToCustody` and `_clearTrigger` moved down into
`ProcessorEndpointStorage` because both sides need them.

Deliberately **not** moved: `claim` measures −5 bytes, because `_claim` and `_asyncTransfer` have to
stay for `stateUpdate` and `_invokeTrigger`, so only the external wrapper would relocate.
`_invokeTrigger` + `_enqueueTrustedRequest` (2,110) is internal and reachable only from
`stateUpdate`, so no stub can reach it.

### 3.1. Per-function stub, not a `fallback()`

Each moved function keeps a typed declaration on the endpoint whose body is
`_delegateToExtension()`. A generic `fallback()` router would be cheaper per moved function but
would require a **merged ABI** (endpoint + extension) for `abigen`, typechain and the hardhat
fixture, since the moved functions would no longer appear in the endpoint's own ABI.

With the stub instead:

- the ABI, the function selectors, the address relayers call, and the `IProcessorEndpoint`
  conformance check are all unchanged;
- no tooling, test, binding or subgraph had to learn about the split.

Measured on the 10-function split, a stub plus its dispatcher entry costs **~47 bytes** (419 bytes
across the 9 functions added after `submitRequestFor`), not the ~150 originally estimated — so the
stub route stays the better deal considerably longer than expected.

**Read-only functions cannot use it at all.** A `view` function may not `delegatecall`:

```
TypeError: Function cannot be declared as view because this expression (potentially) modifies the state.
```

`staticcall` is not a substitute — it would run the extension's code against the *extension's* own
empty storage. Dropping `view` from the endpoint's declaration would be an ABI change that breaks
`abigen` call bindings and any on-chain `staticcall` consumer. So the queue views and the other
getters can only move behind a generic `fallback()` with a merged ABI, which is the next lever:
measured on the current baseline that is **−1,837** for the five queue views and **−2,353** for the
whole read-only surface. Note this supersedes the earlier −2,833 estimate for the queue views, which
was a deletion measurement on the pre-split `paris` baseline; it does not reproduce under
`cancun`. Once that fallback exists, routing the already-moved write entry points through it too
saves a further 419 bytes over their typed stubs.

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
- **The extension refuses direct calls to the entry points it hosts** (`DirectCallNotAllowed`, via
  an `immutable _self` compared against `address(this)`). Called directly, `submitRequestFor` would
  read and write the extension's own empty storage and could strand the ETH sent as a fee. What the
  modifier does *not* cover is everything the extension inherits from `ProcessorEndpointStorage`:
  the `AccessControl` mutators and every public state getter remain callable on the extension
  itself. That is harmless — the getters return zeros, and `grantRole` reverts because no account
  holds `DEFAULT_ADMIN_ROLE` in the extension's own storage, which nothing initialises — but it does
  mean an explorer shows a second, zeroed contract with a `ProcessorEndpoint`-shaped surface.
  Anything moved to the extension in future must carry `onlyDelegateCall`.
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
  endpoint too — which is why `npm run check:layout` asserts it rather than leaving it to review
  (section 5.3). The changed revert error on a malformed signature is an ABI change, and is recorded
  in `CHANGELOG.md`.
- **The stub keeps its parameter names.** Parameter names are part of the published ABI — wallets,
  explorers and generated bindings display them, and `abigen` degrades to `arg0…arg9` without
  them. The names are silenced with no-op statements rather than commented out; verified to produce
  byte-identical output.
- **Deploy order changed.** The extension is deployed first and its address is a new last
  `ProcessorEndpoint` constructor argument. It is `immutable`, so it lives in the endpoint's code
  rather than storage and cannot be repointed. Wired in `test/ProcessorEndpoint/fixture.ts`,
  `scripts/deploy/all.ts`, `scripts/deploy/processorEndpoint.ts` and
  `pkg/blockchain/testutil/sim_test_helper.go`. The extension has its own `go:generate` chain and
  bindings package, needed only for deployment. Inserting a deployment also shifts the deterministic
  CREATE addresses of everything after it in `all.ts`, which is why the local-dev addresses in
  `dockerfiles/README.md` moved.
- **The pairing is readable on-chain**, via `extension()`. Without it, which extension a deployed
  endpoint delegates to would only be recoverable from the deployment record or by disassembling the
  endpoint's code — a poor position for verification or incident response, given the address cannot
  be repointed. It costs **+79 bytes** under the current `cancun` settings. (Under `paris` the same
  getter measured *−318* bytes — an optimizer artifact at `runs: 0`, and a reminder that at this
  optimizer setting a size claim only holds for the exact settings that produced it.)
- **Runtime cost**: one extra cold-account access (~2,600 gas) per forwarded call, on the
  relayer-paid facilitator path. OpenZeppelin's `EIP712` recomputes the domain separator whenever
  `address(this)` differs from the deploying address, so `verifyingContract` remains the endpoint
  and existing signatures stay valid.
- **Compatible with `UPGRADABLE_CONTRACTS_DESIGN.md`**: proxy → implementation → extension are both
  `delegatecall`s, so storage stays in the proxy throughout, and an `immutable` in the
  implementation is still readable. Upgrading the extension means deploying a new implementation
  pointing at it — the normal UUPS path. The `__gap` that design requires belongs in the storage
  base. The reverse mistake — changing `ProcessorEndpointExtension.sol` but upgrading against the
  already-deployed extension — is silent, so `scripts/upgrade/processorEndpoint.ts` compares the
  deployed extension's code with the local build and refuses on mismatch
  (`scripts/upgrade/extensionBytecode.ts`); the endpoint's constructor cannot catch it, because a
  stale extension still has code.

## 4. Result

| Contract | Deployed bytes | vs limit |
|---|---|---|
| `ProcessorEndpoint` before (`paris`) | 23,246 | −1,330 |
| `ProcessorEndpoint` after the split (`paris`) | 21,074 | −3,502 |
| `ProcessorEndpoint`, facilitator path only, `cancun` + `extension()` | 20,433 | −4,143 |
| `ProcessorEndpoint` with deploy submission, resets and admin setters also moved | **17,295** | **−7,281** |
| `ProcessorEndpointExtension` | 11,131 | −13,445 |

The first split sheds 2,172 bytes: the 2,537-byte facilitator module, less the forwarding stub and
the two thin public wrappers (`getPendingRequestsSize`, `generateRequestId`) that now delegate to
internal implementations in the base. `evmVersion: cancun` (section 2.1) takes another ~720, and the
`extension()` getter adds 79 back. Moving the remaining nine non-hot-path entry points (section 3)
takes a further 3,138.

**This covers batch execution — measured, not extrapolated.** Section 1 measured the batch work at
26,353 bytes on the old 23,246-byte `paris` baseline — +3,107 — and `batchStateUpdate()` was not
written yet. Merged onto the 17,295-byte baseline, the per-application queues, the deploy queue, the
round-robin cursor and the selection view land at **19,626 bytes, 4,950 under the limit**
(`ProcessorEndpointExtension` at 12,352), against the extrapolated ~20,400. `batchStateUpdate()` then
cost **4,130 bytes** — `_processOneStateUpdate`, the entry loop and the `BatchEntry[]` calldata
decoder — leaving the endpoint at **23,756 bytes, 820 under the limit**. It fits, but that is now the
whole budget: the next contract change of any size needs one of the levers below. Same settings as
the rest of this section: hardhat, `cancun`, `runs: 0`.

The merge is also where the two designs meet in the extension: the deploy-submission and reset entry
points hosted there operate on the batch branch's per-application queue state (`RequestQueues.Store`
behind the shared `_q`), not on the single global queue they were moved with.

If more room is needed after that, in increasing order of cost:

1. the read-only surface behind a generic `fallback()` and a merged ABI (−2,353, section 3.1);
2. retiring the `stateUpdate` entry point, which `batchStateUpdate` now subsumes: a one-entry batch
   is equivalent to it and takes the same signature, so once the manager routes its single-request
   path through `SubmitBatchStateUpdate` the wrapper and its 12-argument calldata decoder can go
   (`BATCH_EXECUTION.md` section 5.4). This costs no gas on any path;
3. moving one of the two state-update entry points into the extension, which puts ~2,600 gas on every
   update it hosts. Less attractive than it looks now that both share `_processOneStateUpdate` and its
   helpers (`_enforceSelection`, `_queueOf`, `_markRequestCompleted`, `_invokeTrigger`): moving either
   one alone means duplicating those into the extension or promoting them to the storage base, so the
   endpoint sheds only the moved entry point's own decoder and loop.

Note that the `go:generate` solc invocation passes `--optimize` (default `runs: 200`) and no
`--evm-version`, so it uses solc 0.8.30's own default EVM version and produces different sizes from
hardhat's `cancun`/`runs: 0` build. **Any size claim must say which path and which settings produced
it**; the numbers in this document are the hardhat path, and — as the `extension()` getter shows —
at `runs: 0` a delta measured under one `evmVersion` does not carry to another.

## 5. Guardrails

All three are verified to fail when they should, not just to pass.

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

### 5.3. Every extension error must be declared on the endpoint

The same script asserts that each error in the extension's ABI also appears in the endpoint's. The
failure mode is quieter than a revert: the extension's code reverts from the endpoint's address, and
the endpoint's ABI is all a client has to decode it with, so an undeclared error arrives as opaque
bytes. Section 3.4 already required this of `ECDSA.tryRecover`; leaving the rule as prose would have
put the burden on whoever moves the next module.

`DirectCallNotAllowed` is the one allowed exception, listed in
`ERRORS_NOT_REACHABLE_THROUGH_ENDPOINT` in the script: it can only fire when the extension is called
directly, which never happens through the endpoint. Verified by adding an error to the extension and
reverting with it: the check names the error and exits non-zero.

## 6. Reviewing this change

The refactor is behaviour-preserving, and the test suite is the evidence:

- **291 tests before, 297 after.** The six additions are the two constructor cases (`extension is
  zero address`, `extension has no code`), `DirectCallNotAllowed`, the two malformed-signature cases
  that must surface as `InvalidSignature`, and `extension()` returning the deployed extension.
- **No existing test's assertions changed.** Only five test files are touched, and every change is
  either the added constructor argument, the fixture deploying the extension, or one of the new
  tests.
- `npm run check:layout` passes, the full Go suite passes, and `go generate ./...` is idempotent.

One behavioural difference is deliberate and worth naming, and it is covered by the two
malformed-signature tests above as well as by the guardrail in section 5.3: the endpoint reverts
`InvalidSignature` where it previously reverted `ECDSAInvalidSignature*` on a malformed facilitator
signature (see section 3.4). This is an ABI change and is in `CHANGELOG.md`. Nothing else changes. In
particular
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
