import { ethers, upgrades } from 'hardhat';

/*
  Upgrades the TeeAuthenticator proxy to a new implementation.

  Required environment variable:
    PROXY_TEE_AUTHENTICATOR: address of the existing TeeAuthenticator proxy.

  Optional environment variables (only needed when the new implementation adds
  state that must be seeded via a reinitializer):
    REINITIALIZE_FN: name of the reinitializer function to call (e.g. "reinitialize").
    REINITIALIZE_ARGS: JSON-encoded array of arguments for REINITIALIZE_FN.
*/
async function upgrade() {
  const proxyAddress = process.env.PROXY_TEE_AUTHENTICATOR!;
  if (!proxyAddress) {
    throw new Error('PROXY_TEE_AUTHENTICATOR must be set to the existing proxy address');
  }

  const oldImplementation = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log(`TeeAuthenticator proxy: ${proxyAddress}`);
  console.log(`Current implementation: ${oldImplementation}`);

  const TeeAuthenticatorV2 = await ethers.getContractFactory('TeeAuthenticator');

  const options: { kind: 'uups'; call?: { fn: string; args: any[] } } = { kind: 'uups' };
  if (process.env.REINITIALIZE_FN) {
    options.call = {
      fn: process.env.REINITIALIZE_FN,
      args: process.env.REINITIALIZE_ARGS ? JSON.parse(process.env.REINITIALIZE_ARGS) : [],
    };
  }

  const upgraded = await upgrades.upgradeProxy(proxyAddress, TeeAuthenticatorV2, options);
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
