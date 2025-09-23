import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import 'dotenv/config'

// Set your preferred authentication method
//
// If you prefer using a mnemonic, set a MNEMONIC environment variable
// to a valid mnemonic
const MNEMONIC = process.env.MNEMONIC

// If you prefer to be authenticated using a private key, set a PRIVATE_KEY environment variable
const PRIVATE_KEY = process.env.PRIVATE_KEY

const accounts = MNEMONIC
    ? { mnemonic: MNEMONIC }
    : PRIVATE_KEY
      ? [PRIVATE_KEY]
      : undefined

const config: HardhatUserConfig = {
    paths: {
        cache: 'cache/hardhat',
    },
    defaultNetwork: (process.env.NETWORK || "hardhat"),
    solidity: {
        compilers: [
            {
                version: '0.8.28',
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200,
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
        'local': {
            url: 'http://localhost:8545/',
            accounts,
        }, 
        
    },
}

export default config
