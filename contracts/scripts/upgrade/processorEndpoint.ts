import { ethers, upgrades } from 'hardhat';
import { PROCESSOR_ENDPOINT_UPGRADE_OPTIONS } from '../util';

// Upgrades the ProcessorEndpoint proxy to a new implementation. See
// docs/design/UPGRADABLE_CONTRACTS_DESIGN.md for the UUPS upgrade flow.
//
// Required environment variable:
//   PROXY_PROCESSOR_ENDPOINT - the existing proxy address.
// Optional environment variable:
//   PROCESSOR_ENDPOINT_EXTENSION - address of the `ProcessorEndpointExtension` the new
//     implementation's constructor should point at (it is an immutable constructor argument, not
//     proxy state — see the `_extension` doc comment on `ProcessorEndpoint`). Defaults to the
//     extension the current implementation already uses, for a pure logic upgrade that keeps the
//     same extension.
async function upgrade() {
  const proxyAddress = process.env.PROXY_PROCESSOR_ENDPOINT!;
  const [signer] = await ethers.getSigners();

  const processorEndpoint = await ethers.getContractAt('ProcessorEndpoint', proxyAddress);
  const extensionAddress =
    process.env.PROCESSOR_ENDPOINT_EXTENSION || (await processorEndpoint.extension());

  const oldImplementation = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log(
    `upgrading ProcessorEndpoint proxy ${proxyAddress} from ${await signer.getAddress()}`
  );
  console.log(`current implementation: ${oldImplementation}`);
  console.log(`extension for new implementation: ${extensionAddress}`);

  const ProcessorEndpointV2 = await ethers.getContractFactory('ProcessorEndpoint');
  const upgraded = await upgrades.upgradeProxy(proxyAddress, ProcessorEndpointV2, {
    kind: 'uups',
    constructorArgs: [extensionAddress],
    unsafeAllow: [...PROCESSOR_ENDPOINT_UPGRADE_OPTIONS.unsafeAllow],
  });
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
