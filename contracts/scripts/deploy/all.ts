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

  const AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
  const authorityRegistry = await AuthorityRegistry.deploy(deployerAddress);
  await authorityRegistry.deploymentTransaction()!.wait();
  var authorityRegistryAddr = await authorityRegistry.getAddress();
  console.log(`AuthorityRegistry`)
  console.log(`  contract address: ${authorityRegistryAddr}`);

  const TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
  const teeAuthenticator = await TeeAuthenticator.deploy(deployerAddress, process.env.TEE_SIGNER!, process.env.TEE_PUB_SECP521R1!);
  await teeAuthenticator.deploymentTransaction()!.wait();
  var teeAuthenticatorAddr =  await teeAuthenticator.getAddress();
  console.log(`TeeAuthenticator`)
  console.log(`  contract address: ${teeAuthenticatorAddr}`);
  console.log(`  Tee signer address (executor address): ${process.env.TEE_SIGNER!}`);
  console.log(`  Tee PUB_SECP521R1 (executor P521 pub key): ${process.env.TEE_PUB_SECP521R1!}`);

  //deploy 
  const ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
  const processorEndpoint = await ProcessorEndpoint.deploy(teeAuthenticatorAddr, authorityRegistryAddr, process.env.UPDATE_STATUS_OPERATOR!, process.env.ADMIN!, process.env.MIN_FEE_PER_REQUEST!);
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
