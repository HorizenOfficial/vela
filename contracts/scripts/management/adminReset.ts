/**
 * Admin reset script — TESTNET / DEVELOPMENT USE ONLY.
 *
 * Calls adminResetApps(appIds, erc20Tokens) on a deployed ProcessorEndpoint.
 * This clears the pending request queue, resets state roots and locked funds for
 * the given app IDs, and transfers all recovered ETH and ERC-20 balances to the signer.
 *
 * Required environment variables:
 *   PROCESSOR_ENDPOINT   — Address of the deployed ProcessorEndpoint (or proxy).
 *
 * Optional environment variables:
 *   APP_IDS              — Comma-separated uint64 app IDs to reset.
 *                          Omit or leave empty to reset all deployed apps.
 *   ERC20_TOKENS         — Comma-separated ERC-20 token addresses whose custody to recover.
 *                          Omit or leave empty to recover all allowlisted tokens.
 */

import { ethers } from 'hardhat';

async function main() {
  const processorEndpointAddress = process.env.PROCESSOR_ENDPOINT;
  if (!processorEndpointAddress) {
    throw new Error('PROCESSOR_ENDPOINT env var is required');
  }

  const signer = (await ethers.getSigners())[0];
  const signerAddress = await signer.getAddress();
  console.log(`Signer: ${signerAddress}`);
  console.log(`ProcessorEndpoint: ${processorEndpointAddress}`);

  const processorEndpoint = await ethers.getContractAt(
    'ProcessorEndpoint',
    processorEndpointAddress,
    signer
  );

  // Resolve app IDs
  let appIds: bigint[] = [];
  if (process.env.APP_IDS && process.env.APP_IDS.trim() !== '') {
    appIds = process.env.APP_IDS.split(',').map((id) => BigInt(id.trim()));
    console.log(`App IDs (from env): ${appIds.join(', ')}`);
  } else {
    appIds = await processorEndpoint.getDeployedAppIds();
    console.log(
      `App IDs (from contract): ${appIds.length > 0 ? appIds.map((id) => id.toString()).join(', ') : '(none)'}`
    );
  }

  // Resolve ERC-20 token addresses
  let erc20Tokens: string[] = [];
  if (process.env.ERC20_TOKENS && process.env.ERC20_TOKENS.trim() !== '') {
    erc20Tokens = process.env.ERC20_TOKENS.split(',').map((addr) => addr.trim());
    console.log(`ERC-20 tokens (from env): ${erc20Tokens.join(', ')}`);
  } else {
    erc20Tokens = await processorEndpoint.getAllowedTokens();
    console.log(
      `ERC-20 tokens (from contract allowlist): ${erc20Tokens.length > 0 ? erc20Tokens.join(', ') : '(none)'}`
    );
  }

  console.log('\nCalling adminResetApps...');
  const tx = await processorEndpoint.adminResetApps(appIds, erc20Tokens);
  console.log(`  tx hash: ${tx.hash}`);
  const receipt = await tx.wait();
  console.log(`  confirmed in block ${receipt.blockNumber}`);

  const queueSize = await processorEndpoint.getPendingRequestsSize();
  console.log(`\nQueue size after reset: ${queueSize} (expected: 0)`);
  console.log('Reset complete.');
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
