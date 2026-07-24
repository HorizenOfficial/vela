import { ethers, upgrades } from 'hardhat';

async function deploy() {
  const deployer = (await ethers.getSigners())[0];

  console.log(`deploying from ${await deployer.getAddress()}`);
  console.log(`parameters:
    _teeAuthenticator: ${process.env.TEE_AUTHENTICATOR}
    _authorityRegistry: ${process.env.AUTHORITY_REGISTRY}
    updateStatusOperator: ${process.env.UPDATE_STATUS_OPERATOR}
    admin (bootstrap deployer): ${process.env.ADMIN}
    resetOperator (address(0) = disabled): ${process.env.RESET_OPERATOR || '0x0000000000000000000000000000000000000000'}
    minFeePerRequest: ${process.env.MIN_FEE_PER_REQUEST}
  `);

  // 1) TokenAllowlist
  const TokenAllowlist = await ethers.getContractFactory('TokenAllowlist');
  const tokenAllowlist = await TokenAllowlist.deploy(process.env.ADMIN!);
  await tokenAllowlist.deploymentTransaction()!.wait();
  const tokenAllowlistAddr = await tokenAllowlist.getAddress();
  console.log(`TokenAllowlist deployed at ${tokenAllowlistAddr}`);

  // 2) ProcessorEndpoint (proxy)
  const ProcessorEndpoint = await ethers.getContractFactory('ProcessorEndpoint');
  const processorEndpoint = await upgrades.deployProxy(
    ProcessorEndpoint,
    [
      process.env.TEE_AUTHENTICATOR!,
      process.env.AUTHORITY_REGISTRY!,
      process.env.UPDATE_STATUS_OPERATOR!,
      process.env.ADMIN!,
      process.env.RESET_OPERATOR || '0x0000000000000000000000000000000000000000',
      process.env.MIN_FEE_PER_REQUEST!,
      tokenAllowlistAddr,
    ],
    { kind: 'uups' }
  );
  await processorEndpoint.waitForDeployment();

  console.log(`TokenAllowlist deployed at ${tokenAllowlistAddr}`);
  console.log(`ProcessorEndpoint proxy deployed at ${await processorEndpoint.getAddress()}`);
  console.log(
    `ProcessorEndpoint implementation deployed at ${await upgrades.erc1967.getImplementationAddress(await processorEndpoint.getAddress())}`
  );
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
