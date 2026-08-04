import { HardhatUserConfig } from 'hardhat/config';
import '@nomicfoundation/hardhat-toolbox';
import 'dotenv/config';

// Set your preferred authentication method
//
// If you prefer using a mnemonic, set a MNEMONIC environment variable
// to a valid mnemonic
const MNEMONIC = process.env.MNEMONIC;

// If you prefer to be authenticated using a private key, set a PRIVATE_KEY environment variable
const PRIVATE_KEY = process.env.PRIVATE_KEY;

const accounts = MNEMONIC ? { mnemonic: MNEMONIC } : PRIVATE_KEY ? [PRIVATE_KEY] : undefined;

const config: HardhatUserConfig = {
  paths: {
    cache: 'cache/hardhat',
  },
  defaultNetwork: process.env.NETWORK || 'hardhat',
  solidity: {
    compilers: [
      {
        version: '0.8.30',
        settings: {
          optimizer: {
            enabled: true,
            runs: 0,
          },
          viaIR: true,
          // The deployment target is Base, which has been on Cancun since March 2024. Targeting it
          // (rather than the older `paris` default) buys ~720 bytes on ProcessorEndpoint, mostly
          // from PUSH0, which matters because that contract lives against the EIP-170 limit — see
          // docs/design/PROCESSOR_ENDPOINT_SPLIT.md. Any chain this is deployed to must therefore
          // be at Cancun or later.
          evmVersion: 'cancun',
          // Needed by scripts/checkStorageLayout.ts, which asserts that ProcessorEndpoint and
          // ProcessorEndpointExtension share one storage layout — they must, because the endpoint
          // reaches the extension by delegatecall.
          outputSelection: {
            '*': {
              '*': ['storageLayout'],
            },
          },
        },
      },
    ],
  },
  networks: {
    'sepolia-testnet': {
      url: 'https://rpc.sepolia.org/',
      accounts,
    },
    'horizen-l3-testnet': {
      url: 'https://horizen-testnet.rpc.caldera.xyz/http',
      accounts,
    },
    local: {
      url: process.env.CHAIN_RPC_URL || 'http://127.0.0.1:8545/',
      accounts,
    },
    hardhat: {
      // EIP-170 is enforced by default so the suite fails when a contract outgrows the limit,
      // the way it would on Base. `npm run test:ignoresize` opts out, for measuring how far over
      // the limit a work-in-progress contract is.
      allowUnlimitedContractSize: process.env.UNLIMITED_SIZE === 'true',
    },
  },
};

export default config;
