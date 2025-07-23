async function deploy()  {

  const deployer = (await ethers.getSigners())[0]

  console.log(`deploying from ${await deployer.getAddress()}`)
  console.log(`parameters:
    owner: ${process.env.OWNER}
    _teeAuthenticator: ${process.env.TEE_AUTHENTICATOR}
    _authorityRegistry: ${process.env.AUTHORITY_REGISTRY}
    updateStatusOperator: ${process.env.UPDATE_STATUS_OPERATOR}
  `)
  //deploy 
  const ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
  const processorEndpoint = await ProcessorEndpoint.deploy(process.env.TEE_AUTHENTICATOR, process.env.AUTHORITY_REGISTRY, process.env.UPDATE_STATUS_OPERATOR);
  await processorEndpoint.deploymentTransaction().wait();

  console.log(`contract deployed at ${await processorEndpoint.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
