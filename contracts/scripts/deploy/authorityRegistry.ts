import { ethers } from 'hardhat';

async function deploy() {
  const deployer = (await ethers.getSigners())[0];
  const deployerAddress = await deployer.getAddress();

  const owner = process.env.AUTHREGISTRY_OWNER ?? deployerAddress;

  console.log(`deploying authority contracts from ${deployerAddress}`);
  console.log(`parameters:
    owner (for DefaultAuthority & AuthorityRegistry): ${owner}
  `);

  // 1) Deploy DefaultAuthority
  const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
  const defaultAuthority = await DefaultAuthority.deploy(owner);
  await defaultAuthority.deploymentTransaction()!.wait();
  const defaultAuthorityAddr = await defaultAuthority.getAddress();
  console.log(`DefaultAuthority`);
  console.log(`  contract address: ${defaultAuthorityAddr}`);

  // 2) Deploy AuthorityRegistry (proxy) apuntando al default
  const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
  const authorityRegistry = await AuthorityRegistry.deploy(owner, defaultAuthorityAddr);
  await authorityRegistry.deploymentTransaction()!.wait();
  const authorityRegistryAddr = await authorityRegistry.getAddress();
  console.log(`AuthorityRegistry`);
  console.log(`  contract address: ${authorityRegistryAddr}`);
  console.log(`  default authority contract: ${defaultAuthorityAddr}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
