# Horizen PES Contracts

## Manual deploy

- Create an .env file with the correct properties:
    - NETWORK : network to use
    - MNEMONIC or PRIVATE_KEY : identifier of the address deployer of the contracts (must have ETH to pay the tx fees)
    - TEE_SIGNER : Executor's public communication key (P521) (printend by the manager after the first key handshake)
    - TEE_PUB_SECP521R1 : Executor's signing key address (Secp256k1) (printend by the manager after the first key handshake)
    - UPDATE_STATUS_OPERATOR : Manager's address posting the update payloads (printend by the manager after the first key handshake)
    - MIN_FEE_PER_REQUEST: Minimum fee per request applied
    - ADMIN: admin address allowed to change the min fee per request parameter and the queue threeshold on the ProcessorEndpoint contract

- run  deploy script with the needed parameters in the .env file:

```bash
npx hardhat run scripts/deploy/all.ts
```

- separate scripts are also available to deploy each contract one by one. For example:

```bash
npx hardhat run scripts/deploy/authorityRegistry.ts
```