import { ethers } from 'hardhat';

/*
 * Management script used to allowlist an ERC-20 token on TokenAllowlist.
 *
 * Can be run with
 * npx hardhat run scripts/management/addAllowedToken.ts
 *
 * Edit the below constants before running:
 * TOKEN_ALLOWLIST  => address of the TokenAllowlist contract
 * TOKEN_TO_ALLOW   => address of the ERC-20 token to add to the allowlist
 *                     (must be a deployed contract; address(0) is rejected by
 *                     the contract — ETH is always implicitly allowed)
 */
const TOKEN_ALLOWLIST = '0x0000000000000000000000000000000000000000';
const TOKEN_TO_ALLOW = '0x0000000000000000000000000000000000000000';

async function addAllowedToken() {
  const caller = (await ethers.getSigners())[0];
  const callerAddress = await caller.getAddress();
  console.log(`Calling addAllowedToken from ${callerAddress}`);
  console.log(`(this address MUST have the ADMIN role on TokenAllowlist)`);
  console.log(`  TokenAllowlist: ${TOKEN_ALLOWLIST}`);
  console.log(`  token to allow: ${TOKEN_TO_ALLOW}`);

  const tokenAllowlist = await ethers.getContractAt('TokenAllowlist', TOKEN_ALLOWLIST);
  const tx = await tokenAllowlist.addAllowedToken(TOKEN_TO_ALLOW);
  await tx.wait();
}

addAllowedToken()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
