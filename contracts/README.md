# Horizen Vela Contracts

On-chain coordination for Vela. The main contract is `ProcessorEndpoint`, which handles user requests, per-app locked funds, ERC-20 deposits/withdrawals (with EIP-2612 permit and a facilitator path for gasless onboarding), and the deploy flow for new applications. TEE attestation is verified by `TeeAuthenticator` and authorities are tracked in `AuthorityRegistry`.

`ProcessorEndpoint` sits close to the 24,576-byte EIP-170 limit, so it is split across three files: `ProcessorEndpointStorage.sol` declares all of its state, `ProcessorEndpointExtension.sol` hosts entry points moved out for size (currently `submitRequestFor`), and `ProcessorEndpoint.sol` reaches them by `delegatecall`. This is invisible to callers — same address, same ABI, same selectors — but it constrains development and deployment:

- Declare new state **only** in `ProcessorEndpointStorage.sol`, and run `npm run check:layout` afterwards. The two contracts must share one storage layout; a divergence corrupts storage silently instead of reverting. The same command also checks that every error the extension can raise is declared on the endpoint — otherwise it reaches clients as revert data the endpoint's ABI cannot decode.
- Deploy `ProcessorEndpointExtension` **before** `ProcessorEndpoint` — its address is the endpoint's last constructor argument, and it is immutable, so a wrong address means redeploying the endpoint. `ProcessorEndpoint.extension()` reports the pairing of a deployed endpoint.
- Anything moved to the extension needs a typed stub on the endpoint and the `onlyDelegateCall` modifier.

Because that contract lives against the size limit, the compiler targets `evmVersion: 'cancun'` (worth ~720 bytes on it), so **every chain these contracts are deployed to must be at Cancun or later**. Base, the deployment target, has been since March 2024.

See `../docs/design/PROCESSOR_ENDPOINT_SPLIT.md`.

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

  `scripts/deploy/processorEndpoint.ts` deploys `ProcessorEndpointExtension` together with the endpoint, in that order; there is no script that deploys the endpoint against a pre-existing extension.

- when `DEPLOY_OUTPUT_DIR` is set, `all.ts` writes the deployed addresses to `deployed_addresses.env`, including `CHAIN_PROCESSOR_EXTENSION_ADDRESS`. Nothing reads that one at runtime, but it is the only record of which extension an endpoint delegates to.

## Management scripts

The `scripts/management/` folder contains utilities to operate the deployed contracts from the admin account, for example:

```bash
# Grant a wallet the DEPLOYAPP role so it can submit deploy requests
npx hardhat run scripts/management/addAllowedDeployer.ts
```