# Horizen PES Contracts

## Manual deploy

- Create an .env file with the correct properties (NETWORK, MNEMONIC, PRIVATE_KEY)
- run each deploy script in the script/deploy folder separately with the needed parameters in the .env file. For example:

```bash
npx hardhat run scripts/deploy/authorityRegistry.ts
```