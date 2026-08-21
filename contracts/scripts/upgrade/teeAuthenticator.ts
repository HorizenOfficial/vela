import { ethers, upgrades } from 'hardhat';

// Upgrades the TeeAuthenticator proxy to a new implementation. See
// docs/design/UPGRADABLE_CONTRACTS_DESIGN.md for the UUPS upgrade flow.
//
// Not applicable to `NoAttestationTeeAuthenticator` (dev/test mode): that contract is
// deliberately not upgradable, see the design doc's note on `TEE_NO_ATTESTATION=true` deploys.
//
// Required environment variable:
//   PROXY_TEE_AUTHENTICATOR - the existing proxy address.
async function upgrade() {
  const proxyAddress = process.env.PROXY_TEE_AUTHENTICATOR!;
  const [signer] = await ethers.getSigners();

  const oldImplementation = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log(`upgrading TeeAuthenticator proxy ${proxyAddress} from ${await signer.getAddress()}`);
  console.log(`current implementation: ${oldImplementation}`);

  const TeeAuthenticatorV2 = await ethers.getContractFactory('TeeAuthenticator');
  const upgraded = await upgrades.upgradeProxy(proxyAddress, TeeAuthenticatorV2, { kind: 'uups' });
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
