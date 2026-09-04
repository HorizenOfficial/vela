# Horizen Vela Contracts

On-chain coordination for Vela. The main contract is `ProcessorEndpoint`, which handles user requests, per-app locked funds, ERC-20 deposits/withdrawals (with EIP-2612 permit and a facilitator path for gasless onboarding), and the deploy flow for new applications. TEE attestation is verified by `TeeAuthenticator` and authorities are tracked in `AuthorityRegistry`.

`ProcessorEndpoint` sits close to the 24,576-byte EIP-170 limit, so it is split across three files: `ProcessorEndpointStorage.sol` declares all of its state, `ProcessorEndpointExtension.sol` hosts entry points moved out for size (the facilitator path, deploy submission, the operator resets and the admin setters — everything off the per-request hot path), and `ProcessorEndpoint.sol` reaches them by `delegatecall`. This is invisible to callers — same address, same ABI, same selectors — but it constrains development and deployment:

- Declare new state **only** in `ProcessorEndpointStorage.sol`, and run `npm run check:layout` afterwards. The two contracts must share one storage layout; a divergence corrupts storage silently instead of reverting. The same command also checks that every error the extension can raise is declared on the endpoint — otherwise it reaches clients as revert data the endpoint's ABI cannot decode.
- Deploy `ProcessorEndpointExtension` **before** `ProcessorEndpoint` — its address is the endpoint's last constructor argument, and it is immutable, so a wrong address means redeploying the endpoint. `ProcessorEndpoint.extension()` reports the pairing of a deployed endpoint.
- Anything moved to the extension needs a typed stub on the endpoint and the `onlyDelegateCall` modifier.

Because that contract lives against the size limit, the compiler targets `evmVersion: 'cancun'` (worth ~720 bytes on it), so **every chain these contracts are deployed to must be at Cancun or later**. Base, the deployment target, has been since March 2024.

See `../docs/design/PROCESSOR_ENDPOINT_SPLIT.md`.

`ProcessorEndpoint`, `AuthorityRegistry` and `TeeAuthenticator` are all deployed behind a UUPS proxy (`ERC1967Proxy`), so the address printed by a deploy script — and the one every off-chain consumer should use — is the **proxy** address, not the implementation's. `NoAttestationTeeAuthenticator` (dev/test mode) is the one exception and is not upgradable. See `../docs/design/UPGRADABLE_CONTRACTS_DESIGN.md` for the full design and [Upgrading](#upgrading) below for how to push a new implementation.

## Manual deploy

- Create an .env file with the correct properties:
    - NETWORK : network to use
    - MNEMONIC or PRIVATE_KEY : identifier of the address deployer of the contracts (must have ETH to pay the tx fees)
    - ADMIN: admin address allowed to change the min fee per request parameter and the queue threeshold on the ProcessorEndpoint contract
    - TEE_PCR0 : PCR0 of the Nitro application
    - TEE_MAX_VERIFICATION_AGE : Max age for considering the Nitro attestation document valid
    - UPDATE_STATUS_OPERATOR : Address of the manager allowed to post attestations
    - MIN_FEE_PER_REQUEST: Minimum fee per request (in wei)
    - DEPLOYER_ADMIN: address bootstrapped on-chain with the `DEPLOYAPP` role (initial allowed deployer for new apps)
    - ERC20_TOKEN: address of the ERC-20 token used for deposits/withdrawals (must implement `ERC20Permit`)

- run  deploy script with the needed parameters in the .env file:

```bash
npx hardhat run scripts/deploy/all.ts
```

- separate scripts are also available to deploy each contract one by one. For example:

```bash
npx hardhat run scripts/deploy/authorityRegistry.ts
```

  `scripts/deploy/processorEndpoint.ts` deploys `ProcessorEndpointExtension` together with the endpoint, in that order; there is no script that deploys the endpoint against a pre-existing extension. `scripts/deploy/processorEndpointExtension.ts` deploys the extension alone, which is what an upgrade that changes it needs — see [Upgrading](#upgrading).

- when `DEPLOY_OUTPUT_DIR` is set, `all.ts` writes the deployed addresses to `deployed_addresses.env`, including `CHAIN_PROCESSOR_EXTENSION_ADDRESS`. Nothing reads that one at runtime, but it is the only record of which extension an endpoint delegates to.

## Upgrading

Each upgradable contract has a script in `scripts/upgrade/` that deploys a new implementation and points the existing proxy at it via `upgradeToAndCall`, keeping the proxy's address and all of its state (see `../docs/design/UPGRADABLE_CONTRACTS_DESIGN.md`). The caller must be the account holding the upgrade authority — `ADMIN` role for `ProcessorEndpoint`, `owner()` for `AuthorityRegistry` and `TeeAuthenticator`.

The scripts need more than the proxy address and the right signer. `upgradeProxy` validates the new implementation's storage layout against the deployed one, and it reads the deployed layout from the OpenZeppelin upgrades manifest — that layout is not recoverable from chain. **The manifest for the target network must be present and current before an upgrade script is run**, otherwise the upgrade aborts with `Deployment at address ... is not registered`.

- **Public networks** — the deploy writes the manifest to `.openzeppelin/`, and it **must be committed**, so the layout history is a repo artifact rather than a file on whoever happened to run the deploy. Without it an upgrade cannot be run from CI or by anyone else. The plugin names the file after the chain: `sepolia.json` for `sepolia-testnet`, and `unknown-<chainId>.json` for `horizen-l3-testnet`, whose chain id is not in the plugin's network-name map. That second name looks disposable but is the one worth keeping — never ignore `.openzeppelin/unknown-*`.
- **Local (`hardhat` / `anvil`)** — nothing lands in the repo. The manifest goes to `$TMPDIR/openzeppelin-upgrades/<hardhat|anvil>-31337-<instanceId>.json`, keyed by the node's instance id, so it does not survive an `anvil` restart even with `--state` persistence. A local upgrade after a restart needs the recovery below.
- **Recovery** — `upgrades.forceImport(proxyAddress, Factory, { kind: 'uups' })` re-registers a proxy whose manifest entry is lost. It takes the layout from whatever source is currently compiled, so it is only correct if that source matches what is actually deployed; importing against the wrong source yields a manifest that passes the layout check while the real layouts diverge.

- `AuthorityRegistry`:
    - PROXY_AUTHORITY_REGISTRY: address of the deployed proxy
    ```bash
    npx hardhat run scripts/upgrade/authorityRegistry.ts
    ```

- `TeeAuthenticator` (not applicable to `NoAttestationTeeAuthenticator`, which is not upgradable):
    - PROXY_TEE_AUTHENTICATOR: address of the deployed proxy
    ```bash
    npx hardhat run scripts/upgrade/teeAuthenticator.ts
    ```

- `ProcessorEndpoint`:
    - PROXY_PROCESSOR_ENDPOINT: address of the deployed proxy
    - PROCESSOR_ENDPOINT_EXTENSION (optional): address of the `ProcessorEndpointExtension` the new implementation's constructor should point at. Defaults to the extension the current implementation already uses — set this when the upgrade also changes the extension, after deploying the new one with `npx hardhat run scripts/deploy/processorEndpointExtension.ts`.
    - ALLOW_STALE_EXTENSION (optional): set to `true` to downgrade the extension bytecode check below to a warning.
    ```bash
    npx hardhat run scripts/upgrade/processorEndpoint.ts
    ```

    Before upgrading, the script checks that the code deployed at whichever extension address it is using matches the locally compiled `ProcessorEndpointExtension`, and refuses the upgrade if it does not. **Changing `ProcessorEndpointExtension.sol` means deploying a new extension and passing its address**: the extension is a separate contract that the endpoint reaches by `delegatecall`, so an upgrade that reuses the old address leaves the new endpoint logic live while every delegated entry point (`submitRequestFor`, deploy submission, `initialize`, `invokeTrigger`, the operator resets, the admin setters) keeps running the old extension code against the proxy's storage — without reverting. Nothing else catches that: the endpoint's constructor only rejects an address with no code, `upgradeProxy` never inspects `constructorArgs`, and `npm run check:layout` compares compiled artifacts without touching the chain.

    A compiler or optimiser-settings change alters the extension's bytecode even when its source is untouched, so the check also fires when the deployed extension is functionally the one intended. That is the case `ALLOW_STALE_EXTENSION=true` exists for — the message says whether the two differ in length (a different source) or only in content (a different build of it).

Each script prints the implementation address before and after the upgrade so it can be checked against what was intended to be deployed.

## Management scripts

The `scripts/management/` folder contains utilities to operate the deployed contracts from the admin account, for example:

```bash
# Grant a wallet the DEPLOYAPP role so it can submit deploy requests
npx hardhat run scripts/management/addAllowedDeployer.ts
```