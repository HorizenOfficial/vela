import { ethers } from 'hardhat';

/*
 * Management script used to revoke DEPLOYER_ROLE on ProcessorEndpoint.
 *
 * Can be run with
 * npx hardhat run scripts/management/removeAllowedDeployer.ts
 *
 * Edit the below constants before running:
 * PROCESSOR_ENDPOINT => address of the ProcessorEndpoint contract to be called
 * DEPLOYER_TO_REMOVE => address whose deploy permission must be revoked
 */
const PROCESSOR_ENDPOINT = '0x0000000000000000000000000000000000000000';
const DEPLOYER_TO_REMOVE = '0x0000000000000000000000000000000000000000';

async function removeAllowedDeployer() {
  const caller = (await ethers.getSigners())[0];
  const callerAddress = await caller.getAddress();
  console.log(`Calling removeAllowedDeployer from ${callerAddress}`);
  console.log(`(this address MUST have the ADMIN role on ProcessorEndpoint)`);

  const processorEndpoint = await ethers.getContractAt('ProcessorEndpoint', PROCESSOR_ENDPOINT);
  const tx = await processorEndpoint.removeAllowedDeployer(DEPLOYER_TO_REMOVE);
  await tx.wait();
}

removeAllowedDeployer()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
