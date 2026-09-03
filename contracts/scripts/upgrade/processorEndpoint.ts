import { ethers, upgrades } from 'hardhat';
import { verifyExtensionBytecode } from './extensionBytecode';

// Upgrades the ProcessorEndpoint proxy to a new implementation. See
// docs/design/UPGRADABLE_CONTRACTS_DESIGN.md for the UUPS upgrade flow.
//
// Required environment variable:
//   PROXY_PROCESSOR_ENDPOINT - the existing proxy address.
// Optional environment variables:
//   PROCESSOR_ENDPOINT_EXTENSION - address of the `ProcessorEndpointExtension` the new
//     implementation's constructor should point at (it is an immutable constructor argument, not
//     proxy state — see the `_extension` doc comment on `ProcessorEndpoint`). Defaults to the
//     extension the current implementation already uses, for a pure logic upgrade that keeps the
//     same extension.
//   ALLOW_STALE_EXTENSION - set to `true` to downgrade the extension bytecode check below to a
//     warning. Needed only to upgrade deliberately against a deployed extension that no longer
//     matches the local build — for instance after a solc or optimiser-settings change, which
//     alters the extension's bytecode even when its source is untouched.
//
// Whichever extension address is used, its deployed code is checked against the locally compiled
// `ProcessorEndpointExtension` before the upgrade: reusing a stale extension is silent, not a
// revert. See scripts/upgrade/extensionBytecode.ts.
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

  await verifyExtensionBytecode(extensionAddress, process.env.ALLOW_STALE_EXTENSION === 'true');

  const ProcessorEndpointV2 = await ethers.getContractFactory('ProcessorEndpoint');
  const upgraded = await upgrades.upgradeProxy(proxyAddress, ProcessorEndpointV2, {
    kind: 'uups',
    constructorArgs: [extensionAddress],
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
