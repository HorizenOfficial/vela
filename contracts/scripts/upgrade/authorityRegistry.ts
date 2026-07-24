import { ethers, upgrades } from 'hardhat';

/*
  Upgrades the AuthorityRegistry proxy to a new implementation.

  Required environment variable:
    PROXY_AUTHORITY_REGISTRY: address of the existing AuthorityRegistry proxy.

  Optional environment variables (only needed when the new implementation adds
  state that must be seeded via a reinitializer):
    REINITIALIZE_FN: name of the reinitializer function to call (e.g. "reinitialize").
    REINITIALIZE_ARGS: JSON-encoded array of arguments for REINITIALIZE_FN.
*/
async function upgrade() {
  const proxyAddress = process.env.PROXY_AUTHORITY_REGISTRY!;
  if (!proxyAddress) {
    throw new Error('PROXY_AUTHORITY_REGISTRY must be set to the existing proxy address');
  }

  const oldImplementation = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log(`AuthorityRegistry proxy: ${proxyAddress}`);
  console.log(`Current implementation: ${oldImplementation}`);

  const AuthorityRegistryV2 = await ethers.getContractFactory('AuthorityRegistry');

  const options: { kind: 'uups'; call?: { fn: string; args: any[] } } = { kind: 'uups' };
  if (process.env.REINITIALIZE_FN) {
    options.call = {
      fn: process.env.REINITIALIZE_FN,
      args: process.env.REINITIALIZE_ARGS ? JSON.parse(process.env.REINITIALIZE_ARGS) : [],
    };
  }

  const upgraded = await upgrades.upgradeProxy(proxyAddress, AuthorityRegistryV2, options);
  await upgraded.waitForDeployment();

  const newImplementation = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log(`New implementation: ${newImplementation}`);
  console.log('Upgrade complete. Proxy address is unchanged:', proxyAddress);
}

upgrade()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
