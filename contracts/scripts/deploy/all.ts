import { ethers } from 'hardhat';

/*
  Deploy scripts for all the contracts.
  Needed parameters are:
  TEE_SIGNER: ethereum address of the TEE signer (is the executor address)
  TEE_PUB_SECP521R1: TEE public SECP521R1 key (is the executor SECP521 key)
  UPDATE_STATUS_OPERATOR: ehtereum address of the processor endpoint status updater (manager address)
*/

async function deploy()  {

  const deployer = (await ethers.getSigners())[0]
  const deployerAddress = await deployer.getAddress()


  console.log(`deploying all  contracts from ${deployerAddress}`)
  console.log(`(this address will be also the owner of the contracts)`)

  // 1) DefaultAuthority
  const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
  const defaultAuthority = await DefaultAuthority.deploy(deployerAddress);
  await defaultAuthority.deploymentTransaction()!.wait();
  const defaultAuthorityAddr = await defaultAuthority.getAddress();
  console.log(`DefaultAuthority`);
  console.log(`  contract address: ${defaultAuthorityAddr}`);

  // 2) AuthorityRegistry (proxy) usando el default
  const AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
  const authorityRegistry = await AuthorityRegistry.deploy(deployerAddress, defaultAuthorityAddr);
  await authorityRegistry.deploymentTransaction()!.wait();
  const authorityRegistryAddr = await authorityRegistry.getAddress();
  console.log(`AuthorityRegistry`);
  console.log(`  contract address: ${authorityRegistryAddr}`);
  console.log(`  default authority contract: ${defaultAuthorityAddr}`);

  // 3) TeeAuthenticator
  const TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
  const teeAuthenticator = await TeeAuthenticator.deploy(deployerAddress, process.env.TEE_SIGNER!, process.env.TEE_PUB_SECP521R1!);
  await teeAuthenticator.deploymentTransaction()!.wait();
  var teeAuthenticatorAddr =  await teeAuthenticator.getAddress();
  console.log(`TeeAuthenticator`)
  console.log(`  contract address: ${teeAuthenticatorAddr}`);
  console.log(`  Tee signer address (executor address): ${process.env.TEE_SIGNER!}`);
  console.log(`  Tee PUB_SECP521R1 (executor P521 pub key): ${process.env.TEE_PUB_SECP521R1!}`);

  // 4) ProcessorEndpoint
  const ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
  const processorEndpoint = await ProcessorEndpoint.deploy(teeAuthenticatorAddr, authorityRegistryAddr, process.env.UPDATE_STATUS_OPERATOR!, process.env.ADMIN!);
  await processorEndpoint.deploymentTransaction()!.wait();
  var processorEndpointAddr =  await processorEndpoint.getAddress();
  console.log(`ProcessorEndpoint`)
  console.log(`  contract address: ${processorEndpointAddr}`);
  console.log(`  update status operator (manager address): ${process.env.UPDATE_STATUS_OPERATOR!}`);

}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
