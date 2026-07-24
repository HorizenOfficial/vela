# Upgradable Smart Contracts

## Summary

This document describes the design to make the Vela smart contracts upgradable without changing their deployed addresses and without risking loss of funds locked in the `ProcessorEndpoint`.

The chosen mechanism is the **UUPS (Universal Upgradeable Proxy Standard)** pattern, as standardised by [ERC-1822](https://eips.ethereum.org/EIPS/eip-1822) and using the storage slot layout defined by [EIP-1967](https://eips.ethereum.org/EIPS/eip-1967). The OpenZeppelin `UUPSUpgradeable` and `Initializable` abstractions are used as implementation libraries.

The contracts covered by this design are:

- `ProcessorEndpoint.sol`
- `TeeAuthenticator.sol`
- `AuthorityRegistry.sol`

---

## Current State

All three contracts use constructor-based initialisation and are deployed as plain (non-proxied) contracts. Their deployed addresses are used by:

- off-chain services (Manager, subgraph, admin CLI);
- other on-chain contracts (e.g. `ProcessorEndpoint` holds a reference to `TeeAuthenticator` and `AuthorityRegistry`);
- end-user wallets and frontends.

Upgrading any of these contracts today requires redeploying to a new address, migrating all references, and — for `ProcessorEndpoint` — either draining all locked funds first or accepting that they are lost.

---

## Goals

- Allow the logic of all three contracts to be updated post-deployment.
- Preserve the contract address visible to users and off-chain systems across upgrades.
- Preserve all locked funds, state roots, and queue contents in `ProcessorEndpoint` across upgrades.
- Prevent an uninitialised proxy from accepting calls.
- Reserve storage space in each contract so future variables can be added without corrupting the existing storage layout.

## Non-Goals

- Migrating or transforming storage data as part of an upgrade (a storage migration is a separate, off-chain operation).
- Upgrading the `NitroProver` contract (it is `immutable`-free and stateless-ish; its upgrade path is out of scope here).
- Supporting multiple independent implementation shards (EIP-2535 Diamond pattern is unnecessarily complex for this use case).
- Providing an automatic rollback mechanism if an upgrade introduces a bug.

---

## Requirements

### R1 — Stable Addresses

After an upgrade, the address of each contract must remain unchanged. Users, off-chain services, and other contracts must not need to update their references.

### R2 — No Fund Loss

An upgrade to `ProcessorEndpoint` must not require draining or migrating ETH balances. All `payments`, `appLockedFunds`, and `_totalDeposits` values must survive the upgrade intact, because they are stored in the proxy's storage and the new implementation reads the same slots.

### R3 — Constructor-Free Initialisation

Upgradable contracts cannot use constructors for initialisation. All constructor logic must be moved into an `initialize(...)` function. The `initialize` function must be callable exactly once, enforced by the `Initializable` base contract.

### R4 — Uninitialised Contract Protection

A proxy that has been deployed but whose `initialize` function has not yet been called must reject all state-modifying calls. This is enforced by the `initializer` modifier from OpenZeppelin's `Initializable` and by requiring that role-restricted functions (or the TEE signer check) can only succeed after initialisation.

### R5 — Controlled Upgrade Authority

The `_authorizeUpgrade` function required by the UUPS pattern must be restricted to the existing owner / admin role for each contract. No unauthorised party may trigger an upgrade.

### R6 — Storage Gap

Each upgradable contract must declare a `uint256[50] private __gap` array at the end of its storage variables. This reserves 50 storage slots that future versions may consume by replacing the gap with new variables. The gap size must be adjusted (reduced) when new variables are added, keeping the total storage footprint constant.

---

## Proxy Pattern: UUPS

### How UUPS Works

In the UUPS pattern there are two deployed contracts:

1. **Proxy** — a minimal contract that stores the address of the current implementation in a well-known storage slot (EIP-1967 slot `0x360894...`) and delegates every call to that address via `delegatecall`. The proxy holds all state.
2. **Implementation** — contains all the business logic. It is stateless by itself; state lives in the proxy.

The key property of `delegatecall` is that it runs the implementation's code but reads and writes the _proxy's_ storage. This is why storage layout compatibility across upgrades is mandatory (see [Storage Layout and Upgrade Safety](#storage-layout-and-upgrade-safety)).

In UUPS, the upgrade function (`upgradeTo` / `upgradeToAndCall`) lives inside the implementation itself, not in a separate admin contract. This saves gas on every normal call because there is no dispatcher overhead. The `_authorizeUpgrade` hook controls who may call it.

```
User  ─────►  Proxy (stable address, holds all state)
                │  delegatecall
                ▼
           Implementation V1  (logic only)

After upgrade:

User  ─────►  Proxy (same address, same state)
                │  delegatecall
                ▼
           Implementation V2  (new logic)
```

### Relevant Standards

| Standard | Role in this design |
|----------|---------------------|
| [EIP-1967](https://eips.ethereum.org/EIPS/eip-1967) | Defines the storage slot where the proxy stores the implementation address (`0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc`) and the admin slot. Ensures tooling (Etherscan, block explorers) can detect proxies automatically. |
| [ERC-1822](https://eips.ethereum.org/EIPS/eip-1822) | Defines the UUPS upgrade pattern: upgrade logic lives inside the implementation, not in a separate proxy admin. OpenZeppelin's `UUPSUpgradeable` is the canonical implementation. |
| [OpenZeppelin `Initializable`](https://docs.openzeppelin.com/contracts/5.x/api/proxy#Initializable) | Provides the `initializer` and `reinitializer` modifiers and the `_disableInitializers()` helper that must be called in the implementation's constructor to prevent direct initialisation of the implementation contract (see R4). |

### Go Test Tooling

`pkg/blockchain/testutil/sim_test_helper.go` deploys `AuthorityRegistry` and `ProcessorEndpoint` against a simulated backend for Go tests. It mirrors `upgrades.deployProxy` by hand: deploy the bare implementation, then deploy an `ERC1967Proxy` (bound via a dedicated `erc1967proxy` Go package, generated the same way as the other contract bindings — see `client_test.go`'s `go:generate` directives) with ABI-encoded `initialize(...)` calldata, and treat the proxy address as the contract's address for all subsequent calls. `TeeAuthenticator`/`AuthorityRegistry`'s combined-json abigen output pulls in OpenZeppelin's `Errors.sol` library, whose generated Go type is literally named `Errors` and shadows the stdlib `errors` package inside abigen's own generated methods; both bindings (and the new `erc1967proxy` one) isolate their own abi/bin via `jq` before invoking `abigen`, the same workaround already used for `MockERC20`/`TestTrigger`.

### Deployment Flow

For each contract the deployment sequence is:

1. Deploy the new `ImplementationV1` contract (constructor calls `_disableInitializers()`).
2. Deploy an `ERC1967Proxy` pointing at `ImplementationV1` and encoding the `initialize(...)` calldata. This atomically deploys the proxy and calls `initialize`.
3. Record the **proxy address** as the canonical address for that contract in all configuration files and registries. The implementation address is an internal detail.

### Upgrade Flow

1. Deploy `ImplementationV2` (constructor calls `_disableInitializers()`).
2. Call `proxy.upgradeToAndCall(address(implementationV2), initData)` from the owner/admin account.
   - `initData` is empty (`""`) for pure logic upgrades that do not need a reinitialiser.
   - `initData` encodes a `reinitialize(uint64 version, ...)` call if new state variables must be seeded.
3. Verify the upgrade via `ERC1967Utils.getImplementation(proxyAddress)`.

---

## Storage Layout and Upgrade Safety

### How Proxy Storage Works

Because `delegatecall` runs the implementation in the proxy's storage context, the EVM maps each state variable to a slot by _position in declaration order_, not by name. Slot 0 is the first declared variable, slot 1 is the second, and so on (with packing rules for small types).

**If an upgrade changes the order or type of existing variables, it corrupts storage.**

The proxy's own EIP-1967 implementation slot is deliberately placed at a pseudo-random high slot (`keccak256("eip1967.proxy.implementation") - 1`) to avoid colliding with the implementation's variables, which start at slot 0.

### Rules for Safe Upgrades

| Rule | Rationale |
|------|-----------|
| Never remove an existing variable | Removing shifts subsequent slots. |
| Never change the type of an existing variable | A `uint128` replaced by `address` would reinterpret the same bits. |
| Never reorder existing variables | Reordering changes slot assignments. |
| Only append new variables at the end | Appended variables consume previously unused slots (or gap slots). |
| Reduce `__gap` by the number of new variables added | This keeps the total storage footprint constant, protecting any variables declared in base contracts that follow. |
| Inherited contracts must also follow these rules | Inheritance order defines slot order; changing base contract variables has the same effect as changing the derived contract. |

### Storage Gap Convention

Each upgradable contract reserves 50 slots at the end of its own variable block:

```solidity
uint256[50] private __gap;
```

When a future version adds `N` new variables, the gap shrinks to `50 - N`:

```solidity
// V2 adds one new variable
uint256 public newVariable;          // consumes one former gap slot
uint256[49] private __gap;           // reduced by 1
```

This pattern is the same one used by OpenZeppelin's own upgradeable contracts (e.g. `OwnableUpgradeable`, `AccessControlUpgradeable`).

---

## Implementation Changes per Contract

All three contracts must apply the same mechanical transformation:

1. Replace standard OpenZeppelin base contracts with their `*Upgradeable` variants.
2. Remove constructor arguments; the constructor body must only call `_disableInitializers()`.
3. Add an `initialize(...)` function with the `initializer` modifier that performs all previous constructor logic. **Important:** inline state-variable initializers (e.g. `uint256 public maxNumOfApplications = 10;`) are part of the constructor in plain contracts but are **not** executed by the proxy. Every such default value must be moved inside `initialize()` explicitly.
4. Inherit from `UUPSUpgradeable` and implement `_authorizeUpgrade`.
5. Add `uint256[50] private __gap` at the end of their own storage variables.

The examples below use `ProcessorEndpoint` as the reference; `TeeAuthenticator` and `AuthorityRegistry` follow the same pattern with their own parameters and base contracts.

### Constructor and Initializer Pattern

```solidity
/// @custom:oz-upgrades-unsafe-allow constructor
constructor() {
    _disableInitializers();
}

function initialize(
    ITeeAuthenticator _teeAuthenticator,
    IAuthorityRegistry _authorityRegistry,
    address updateStatusOperator,
    address admin,
    address resetOperator,
    uint256 _minFeePerRequest,
    ITokenAllowlist _tokenAllowlist
) external initializer {
    // zero-address guards (same as current constructor)
    __AccessControl_init();
    __ReentrancyGuard_init();
    __EIP712_init('Vela', Strings.toString(PROTOCOL_VERSION));
    __UUPSUpgradeable_init();
    // ... assign storage variables (including the inline defaults for
    // maxNumOfApplications/availableDeploySlots/maxQueueSize, which are only
    // executed by a plain constructor and must be set explicitly here)
}
```

The `initializer` modifier guarantees `initialize` can only be called once. Calling `_disableInitializers()` in the constructor prevents anyone from calling `initialize` directly on the implementation contract (bypassing the proxy).

### Upgrade Authorisation

```solidity
function _authorizeUpgrade(address newImplementation)
    internal
    override
    onlyRole(ADMIN)   // onlyOwner for TeeAuthenticator and AuthorityRegistry
{}
```

### Initialisation Guard

No dedicated `initialized` check is needed beyond the `initializer` modifier. An uninitialised proxy will always revert on any meaningful call:

- `ProcessorEndpoint`: all write paths require `onlyRole(...)` or `nonReentrant`; roles are empty until `initialize` is called.
- `TeeAuthenticator` / `AuthorityRegistry`: all write paths require `onlyOwner`; `owner()` returns `address(0)` until `initialize` is called.

### Storage Layout

Existing variables must remain in declaration order, unchanged. The `__gap` buffer is appended last. Example for `ProcessorEndpoint`:

```solidity
// --- existing variables (order locked) ---
mapping(uint64 => bytes32) public applicationStateRoots;
uint64[] private _deployedAppIds;
uint256 public maxNumOfApplications;
uint256 public availableDeploySlots;
RequestQueue private _requestQueue;
uint256 public maxQueueSize;
ITeeAuthenticator public teeAuthenticator;
IAuthorityRegistry public authorityRegistry;
ITokenAllowlist public tokenAllowlist;
mapping(address => mapping(address => uint256)) public pendingClaims;
mapping(address => uint256) public totalPendingClaims;
mapping(uint64 => mapping(address => uint256)) public appCustody;
mapping(address => uint256) public totalAppCustody;
uint256 public minFeePerRequest;
address payable public feeCollector;
mapping(address => uint256) public facilitatorNonces;
mapping(uint64 => ITrigger) public triggerContracts;
mapping(address => uint64) public triggersToAppIds;
RequestQueue private _triggerQueue;

// --- storage buffer (must always be last) ---
uint256[50] private __gap;
```

(Field-level docs are omitted above for brevity — see `ProcessorEndpoint.sol` for the authoritative, commented list. `RequestQueue` is itself a struct of mappings/counters; its internal layout is likewise order-locked.)

### Per-Contract Notes

| Contract | Base contracts to replace | Additional notes |
|----------|--------------------------|-----------------|
| `ProcessorEndpoint` | `AccessControl` → `AccessControlUpgradeable`, `ReentrancyGuard` → `ReentrancyGuardUpgradeable`, `EIP712` → `EIP712Upgradeable` | `EIP712`'s domain separator (`name`/`version`) is set via `__EIP712_init` inside `initialize`, replacing the constructor-only `EIP712('Vela', ...)` base-constructor call. |
| `TeeAuthenticator` | `Ownable` → `OwnableUpgradeable` | `nitroProver` and `maxVerificationAge` are currently `immutable`. `immutable` variables are embedded in bytecode and invisible to the proxy's storage, so they must be converted to regular storage variables. The minor gas cost on reads is accepted given the low call frequency of attestation-related functions. |
| `AuthorityRegistry` | `Ownable` → `OwnableUpgradeable` | No other changes to contract logic. |

---


## Dependency Changes

### npm / Hardhat

```
npm install @openzeppelin/contracts-upgradeable
npm install --save-dev @openzeppelin/hardhat-upgrades
```

The `@openzeppelin/hardhat-upgrades` plugin provides:
- `upgrades.deployProxy(factory, args, { kind: 'uups' })` — deploys the proxy and calls `initialize` atomically.
- `upgrades.upgradeProxy(proxyAddress, newFactory)` — validates storage layout compatibility before upgrading.
- `upgrades.validateUpgrade(oldFactory, newFactory)` — static check, usable in CI.

### Hardhat Config

```typescript
import '@openzeppelin/hardhat-upgrades';
```
---

## Script Changes

### Existing Deploy Scripts

All deploy scripts in `scripts/deploy/` use the vanilla `ContractFactory.deploy(...)` pattern. After the contracts are made upgradable, each script must switch to `upgrades.deployProxy(...)` from `@openzeppelin/hardhat-upgrades`.

The `initialize` arguments replace the old constructor arguments one-to-one, so the required environment variables stay the same.

#### `scripts/deploy/processorEndpoint.ts`

Replace:
```typescript
const processorEndpoint = await ProcessorEndpoint.deploy(
  process.env.TEE_AUTHENTICATOR!,
  process.env.AUTHORITY_REGISTRY!,
  process.env.UPDATE_STATUS_OPERATOR!,
  process.env.ADMIN!,
  process.env.MIN_FEE_PER_REQUEST!
);
await processorEndpoint.deploymentTransaction()!.wait();
```

With:
```typescript
const processorEndpoint = await upgrades.deployProxy(
  ProcessorEndpoint,
  [
    process.env.TEE_AUTHENTICATOR!,
    process.env.AUTHORITY_REGISTRY!,
    process.env.UPDATE_STATUS_OPERATOR!,
    process.env.ADMIN!,
    process.env.MIN_FEE_PER_REQUEST!,
  ],
  { kind: 'uups' }
);
await processorEndpoint.waitForDeployment();
```

The address printed and written to `deployed_addresses.env` is the **proxy** address. The implementation address can be retrieved with `upgrades.erc1967.getImplementationAddress(proxyAddress)` if needed for verification.

#### `scripts/deploy/teeAuthenticator.ts`

Same pattern. Replace:
```typescript
const teeAuthenticator = await TeeAuthenticator.deploy(
  process.env.TEE_OWNER!,
  nitroProverAddress,
  process.env.TEE_PCR0!,
  process.env.TEE_MAX_VERIFICATION_AGE!
);
await teeAuthenticator.deploymentTransaction()!.wait();
```

With:
```typescript
const teeAuthenticator = await upgrades.deployProxy(
  TeeAuthenticator,
  [
    process.env.TEE_OWNER!,
    nitroProverAddress,
    process.env.TEE_PCR0!,
    process.env.TEE_MAX_VERIFICATION_AGE!,
  ],
  { kind: 'uups' }
);
await teeAuthenticator.waitForDeployment();
```

The `NoAttestationTeeAuthenticator` branch (dev mode, `TEE_NO_ATTESTATION=true`) is **not** upgradable and must remain a plain `deploy()`. It is only used in development and does not need proxy support.

#### `scripts/deploy/authorityRegistry.ts`

Replace:
```typescript
const authorityRegistry = await AuthorityRegistry.deploy(owner, defaultAuthorityAddr);
await authorityRegistry.deploymentTransaction()!.wait();
```

With:
```typescript
const authorityRegistry = await upgrades.deployProxy(
  AuthorityRegistry,
  [owner, defaultAuthorityAddr],
  { kind: 'uups' }
);
await authorityRegistry.waitForDeployment();
```

`DefaultAuthority` is not upgraded through a proxy (it holds no state and its address is stored in `AuthorityRegistry`), so its deploy call stays unchanged.

#### `scripts/deploy/all.ts`

Apply the same three substitutions described above to the inline deployment logic in the all-in-one script.

The `deployed_addresses.env` output format does not change; the proxy address is what external consumers need.

---

### New Upgrade Scripts (`scripts/upgrade/`)

Create a new `scripts/upgrade/` directory containing one script per upgradable contract. Each script:

1. Reads the current proxy address from an environment variable.
2. Calls `upgrades.upgradeProxy(proxyAddress, NewFactory)`, which validates storage layout compatibility and atomically upgrades the implementation.
3. Prints the old and new implementation addresses for audit purposes.

#### `scripts/upgrade/processorEndpoint.ts`

Required environment variable: `PROXY_PROCESSOR_ENDPOINT` — the existing proxy address.

```
upgrades.upgradeProxy(
  process.env.PROXY_PROCESSOR_ENDPOINT!,
  ProcessorEndpointV2,
  { kind: 'uups' }
)
```

If the new version adds storage variables that need seeding, pass them via `upgradeToAndCall` by supplying a `call` option:

```
upgrades.upgradeProxy(proxyAddress, Factory, {
  kind: 'uups',
  call: { fn: 'reinitialize', args: [2, ...newArgs] },
})
```

#### `scripts/upgrade/teeAuthenticator.ts`

Required environment variable: `PROXY_TEE_AUTHENTICATOR`.

Same pattern as above.

#### `scripts/upgrade/authorityRegistry.ts`

Required environment variable: `PROXY_AUTHORITY_REGISTRY`.

Same pattern as above.

---

## Note

In UUPS the EIP-1967 admin slot is not used; there is no separate `ProxyAdmin` contract. Upgrade authority is enforced entirely by `_authorizeUpgrade` inside the implementation, which restricts the call to the owner or `ADMIN` role depending on the contract. The account holding that role is therefore the sole key that can push a new implementation. In production it must be held in a multisig or hardware wallet.

Consider wrapping the upgrade authority in an [`OpenZeppelin TimelockController`](https://docs.openzeppelin.com/contracts/5.x/api/governance#TimelockController): without a timelock, a compromised multisig can push a malicious implementation instantly. A timelock introduces a mandatory delay between a proposed upgrade and its execution, giving users time to review the change and withdraw funds if the upgrade looks suspicious.
