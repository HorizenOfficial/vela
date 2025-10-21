# Horizen PES Contracts

## Manual deploy

- Create an .env file with the correct properties (NETWORK, MNEMONIC, PRIVATE_KEY)
- run  deploy script with the needed parameters in the .env file:

```bash
npx hardhat run scripts/deploy/all.ts
```

- separate scripts are also available to deploy each contract one by one. For example:

```bash
npx hardhat run scripts/deploy/authorityRegistry.ts
```