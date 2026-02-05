# Horizen PES Contracts

## Manual deploy

- Create an .env file with the correct properties:
    - NETWORK : network to use
    - MNEMONIC or PRIVATE_KEY : identifier of the address deployer of the contracts (must have ETH to pay the tx fees)
    - ADMIN: admin address allowed to change the min fee per request parameter and the queue threeshold on the ProcessorEndpoint contract  
    - TEE_PCR0 : PCR0 of the Nitro application
    - TEE_MAX_VERIFICATION_AGE : Max age for considering the Nitro attestation document valid
    - UPDATE_STATUS_OPERATOR : Address of the manager allowed to post attestations 
    - MIN_FEE_PER_REQUEST: Minimum fee per request (in wei)


- run  deploy script with the needed parameters in the .env file:

```bash
npx hardhat run scripts/deploy/all.ts
```

- separate scripts are also available to deploy each contract one by one. For example:

```bash
npx hardhat run scripts/deploy/authorityRegistry.ts
```