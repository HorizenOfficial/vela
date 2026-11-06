import { ethers, upgrades } from 'hardhat';

// Upgrades the AuthorityRegistry proxy to a new implementation. See
// docs/design/UPGRADABLE_CONTRACTS_DESIGN.md for the UUPS upgrade flow.
//
// Required environment variable:
//   PROXY_AUTHORITY_REGISTRY - the existing proxy address.
async function upgrade() {
  const proxyAddress = process.env.PROXY_AUTHORITY_REGISTRY!;
  const [signer] = await ethers.getSigners();

  const oldImplementation = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log(
    `upgrading AuthorityRegistry proxy ${proxyAddress} from ${await signer.getAddress()}`
  );
  console.log(`current implementation: ${oldImplementation}`);

  const AuthorityRegistryV2 = await ethers.getContractFactory('AuthorityRegistry');
  const upgraded = await upgrades.upgradeProxy(proxyAddress, AuthorityRegistryV2, { kind: 'uups' });
  await upgraded.waitForDeployment();

  const newImplementation = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log(`new implementation: ${newImplementation}`);
}

upgrade()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
