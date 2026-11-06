import { ethers } from 'hardhat';

// Deploys ProcessorEndpointExtension on its own, for an upgrade that changes it: pass the printed
// address to scripts/upgrade/processorEndpoint.ts in PROCESSOR_ENDPOINT_EXTENSION, which pairs it
// with the new endpoint implementation (the address is an immutable constructor argument — see the
// `_extension` doc comment on `ProcessorEndpoint`).
//
// For a first deployment use scripts/deploy/processorEndpoint.ts or all.ts instead: both deploy
// the extension together with the endpoint and its proxy, in the required order.
async function deploy() {
  const [deployer] = await ethers.getSigners();
  console.log(`deploying ProcessorEndpointExtension from ${await deployer.getAddress()}`);

  const ProcessorEndpointExtension = await ethers.getContractFactory('ProcessorEndpointExtension');
  const extension = await ProcessorEndpointExtension.deploy();
  await extension.deploymentTransaction()!.wait();

  console.log(`ProcessorEndpointExtension deployed at ${await extension.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
