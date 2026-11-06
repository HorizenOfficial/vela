import { ethers, upgrades } from 'hardhat';

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

  // 2) Deploy AuthorityRegistry behind a UUPS proxy, pointing at the default authority
  const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
  const authorityRegistry = await upgrades.deployProxy(
    AuthorityRegistry,
    [owner, defaultAuthorityAddr],
    { kind: 'uups' }
  );
  await authorityRegistry.waitForDeployment();
  const authorityRegistryAddr = await authorityRegistry.getAddress();
  console.log(`AuthorityRegistry (proxy)`);
  console.log(`  contract address: ${authorityRegistryAddr}`);
  console.log(`  default authority contract: ${defaultAuthorityAddr}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
