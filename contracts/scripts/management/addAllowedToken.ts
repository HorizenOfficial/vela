import { ethers } from 'hardhat';

/*
 * Management script used to allowlist an ERC-20 token on ProcessorEndpoint.
 *
 * Can be run with
 * npx hardhat run scripts/management/addAllowedToken.ts
 *
 * Edit the below constants before running:
 * PROCESSOR_ENDPOINT => address of the ProcessorEndpoint contract to be called
 * TOKEN_TO_ALLOW    => address of the ERC-20 token to add to the allowlist
 *                      (must be a deployed contract; address(0) is rejected by
 *                      the contract — ETH is always implicitly allowed)
 */
const PROCESSOR_ENDPOINT = '0x0000000000000000000000000000000000000000';
const TOKEN_TO_ALLOW = '0x0000000000000000000000000000000000000000';

async function addAllowedToken() {
  const caller = (await ethers.getSigners())[0];
  const callerAddress = await caller.getAddress();
  console.log(`Calling addAllowedToken from ${callerAddress}`);
  console.log(`(this address MUST have the ADMIN role on ProcessorEndpoint)`);
  console.log(`  ProcessorEndpoint: ${PROCESSOR_ENDPOINT}`);
  console.log(`  token to allow:    ${TOKEN_TO_ALLOW}`);

  const processorEndpoint = await ethers.getContractAt('ProcessorEndpoint', PROCESSOR_ENDPOINT);
  const tx = await processorEndpoint.addAllowedToken(TOKEN_TO_ALLOW);
  await tx.wait();
}

addAllowedToken()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
