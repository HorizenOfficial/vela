import hre from 'hardhat';
import { spawnSync } from 'child_process';
import * as dotenv from 'dotenv';
import * as fs from 'fs';
import * as os from 'os';
import * as path from 'path';

// Guards the contract addresses documented in dockerfiles/README.md for the dev environment.
//
// The dev deployer runs scripts/deploy/all.ts from one fixed account on a fresh chain, so every
// address is fixed by the deployer nonce, i.e. by the position of the deployment in all.ts.
// Inserting or removing a deployment (a new contract, a contract moved behind a proxy, which takes
// two) shifts every address after it, and nothing in such a change touches the README. This runs
// the real deploy script with the dev settings from dockerfiles/.env.dev on a fresh in-process
// chain and fails if the README no longer matches what it produces.

const CONTRACTS_DIR = path.resolve(__dirname, '..');
const ENV_FILE = path.resolve(CONTRACTS_DIR, '../dockerfiles/.env.dev');
const README = path.resolve(CONTRACTS_DIR, '../dockerfiles/README.md');

// Mirrors the deployer service's `environment` in dockerfiles/docker-compose.yml: deploy-script
// variable <- .env.dev variable. PRIVATE_KEY is not passed: on the in-process network the deployer
// is the default signer, which is checked against DEPLOYER_PRIVATE_KEY below instead.
const DEPLOYER_ENV: Record<string, string> = {
  ADMIN: 'DEPLOYER_ADMIN',
  UPDATE_STATUS_OPERATOR: 'UPDATE_STATUS_OPERATOR',
  TEE_PCR0: 'TEE_PCR0',
  TEE_MAX_VERIFICATION_AGE: 'TEE_MAX_VERIFICATION_AGE',
  TEE_NO_ATTESTATION: 'TEE_NO_ATTESTATION',
  TEE_SIGNER_ADDRESS: 'TEE_SIGNER_ADDRESS',
  TEE_PUB_P521: 'TEE_PUB_P521',
  MIN_FEE_PER_REQUEST: 'MIN_FEE_PER_REQUEST',
  DEPLOY_MOCK_ERC20: 'DEPLOY_MOCK_ERC20',
  MOCK_ERC20_NAME: 'MOCK_ERC20_NAME',
  MOCK_ERC20_SYMBOL: 'MOCK_ERC20_SYMBOL',
  MOCK_ERC20_DECIMALS: 'MOCK_ERC20_DECIMALS',
};

// wallet.conf key in the README <- key in deployed_addresses.env written by all.ts.
const DOCUMENTED = [
  { readmeKey: 'ProcessorAddress', deployedKey: 'CHAIN_PROCESSOR_ADDRESS' },
  { readmeKey: 'TeeAuthenticatorAddress', deployedKey: 'CHAIN_TEEAUTHENTICATOR_ADDRESS' },
];

async function checkDeployerAccount(devEnv: Record<string, string>) {
  const key = devEnv.DEPLOYER_PRIVATE_KEY;
  if (!key) {
    throw new Error(`DEPLOYER_PRIVATE_KEY is not set in ${ENV_FILE}`);
  }
  const expected = new hre.ethers.Wallet(key).address;
  const [signer] = await hre.ethers.getSigners();
  if (signer.address !== expected) {
    throw new Error(
      `the dev deployer (DEPLOYER_PRIVATE_KEY, ${expected}) is not the in-process network's first ` +
        `account (${signer.address}), so the addresses computed here would not be the ones the ` +
        `dev environment produces`
    );
  }
}

function deployFresh(devEnv: Record<string, string>): Record<string, string> {
  const outputDir = fs.mkdtempSync(path.join(os.tmpdir(), 'vela-dev-addresses-'));
  try {
    const env: NodeJS.ProcessEnv = { ...process.env, NETWORK: 'hardhat', HARDHAT_NETWORK: '' };
    for (const [deployVar, envVar] of Object.entries(DEPLOYER_ENV)) {
      env[deployVar] = devEnv[envVar] ?? '';
    }
    env.DEPLOY_OUTPUT_DIR = outputDir;

    // all.ts runs on import and exits the process, so it is run as its own hardhat process. The
    // in-process network starts empty on every run: a fresh chain, as the dev environment assumes.
    const run = spawnSync(
      'npx',
      ['hardhat', 'run', '--network', 'hardhat', 'scripts/deploy/all.ts'],
      {
        cwd: CONTRACTS_DIR,
        env,
        encoding: 'utf8',
      }
    );
    if (run.status !== 0) {
      throw new Error(`scripts/deploy/all.ts failed:\n${run.stdout}\n${run.stderr}`);
    }
    return dotenv.parse(fs.readFileSync(path.join(outputDir, 'deployed_addresses.env')));
  } finally {
    fs.rmSync(outputDir, { recursive: true, force: true });
  }
}

function documentedAddress(readme: string, key: string): string {
  const match = readme.match(new RegExp(`^${key}=(0x[0-9a-fA-F]{40})\\s*$`, 'm'));
  if (!match) {
    throw new Error(`no "${key}=0x..." line in ${README}`);
  }
  return match[1];
}

async function main() {
  const devEnv = dotenv.parse(fs.readFileSync(ENV_FILE));
  await checkDeployerAccount(devEnv);

  const deployed = deployFresh(devEnv);
  const readme = fs.readFileSync(README, 'utf8');

  const mismatches: string[] = [];
  for (const { readmeKey, deployedKey } of DOCUMENTED) {
    const actual = deployed[deployedKey];
    if (!actual) {
      throw new Error(`${deployedKey} missing from the deploy output`);
    }
    const documented = documentedAddress(readme, readmeKey);
    if (documented.toLowerCase() !== actual.toLowerCase()) {
      mismatches.push(
        `  ${readmeKey}: README has ${documented}, the dev deploy produces ${actual}`
      );
    } else {
      console.log(`${readmeKey}=${actual} OK`);
    }
  }

  if (mismatches.length > 0) {
    throw new Error(
      `dockerfiles/README.md is out of date with scripts/deploy/all.ts:\n${mismatches.join('\n')}\n` +
        `Update the wallet.conf block in dockerfiles/README.md (and the copies in vela-starterkit: ` +
        `dockerfiles/README.md, docs/2_private-transfer-app.md).`
    );
  }
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : error);
  process.exitCode = 1;
});
