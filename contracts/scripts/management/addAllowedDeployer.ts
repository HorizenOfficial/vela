import { ethers } from 'hardhat';

/*
 * Management script used to grant DEPLOYER_ROLE on ProcessorEndpoint.
 *
 * Can be run with
 * npx hardhat run scripts/management/addAllowedDeployer.ts
 *
 * Edit the below constants before running:
 * PROCESSOR_ENDPOINT => address of the ProcessorEndpoint contract to be called
 * NEW_DEPLOYER => address to be granted deploy permission
 */
const PROCESSOR_ENDPOINT = '0x0000000000000000000000000000000000000000';
const NEW_DEPLOYER = '0x0000000000000000000000000000000000000000';

async function addAllowedDeployer() {
  const caller = (await ethers.getSigners())[0];
  const callerAddress = await caller.getAddress();
  console.log(`Calling addAllowedDeployer from ${callerAddress}`);
  console.log(`(this address MUST have the ADMIN role on ProcessorEndpoint)`);

  const processorEndpoint = await ethers.getContractAt('ProcessorEndpoint', PROCESSOR_ENDPOINT);
  const tx = await processorEndpoint.addAllowedDeployer(NEW_DEPLOYER);
  await tx.wait();
}

addAllowedDeployer()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
