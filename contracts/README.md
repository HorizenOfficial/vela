# Horizen Vela Contracts

On-chain coordination for Vela. The main contract is `ProcessorEndpoint`, which handles user requests, per-app locked funds, ERC-20 deposits/withdrawals (with EIP-2612 permit and a facilitator path for gasless onboarding), and the deploy flow for new applications. TEE attestation is verified by `TeeAuthenticator` and authorities are tracked in `AuthorityRegistry`.

## Manual deploy

- Create an .env file with the correct properties:
    - NETWORK : network to use
    - MNEMONIC or PRIVATE_KEY : identifier of the address deployer of the contracts (must have ETH to pay the tx fees)
    - ADMIN: admin address allowed to change the min fee per request parameter and the queue threeshold on the ProcessorEndpoint contract
    - TEE_PCR0 : PCR0 of the Nitro application
    - TEE_PCR0_UPGRADE_DELAY : Timelock in seconds before a swap to a new PCR0 can be applied
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

## Management scripts

The `scripts/management/` folder contains utilities to operate the deployed contracts from the admin account, for example:

```bash
# Grant a wallet the DEPLOYAPP role so it can submit deploy requests
npx hardhat run scripts/management/addAllowedDeployer.ts
```