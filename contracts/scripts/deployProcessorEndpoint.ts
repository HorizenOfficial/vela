import hre from "hardhat";


async function deploy()  {

  const deployer = (await ethers.getSigners())[0]

  console.log(`deploying from ${await deployer.getAddress()}`)
  console.log(`parameters: ${process.env.TEE_AUTHENTICATOR}, ${process.env.UPDATE_STATUS_OPERATOR}`)
  //deploy 
  const ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
  const processorEndpoint = await ProcessorEndpoint.deploy(process.env.TEE_AUTHENTICATOR, process.env.UPDATE_STATUS_OPERATOR);
  await processorEndpoint.deploymentTransaction().wait();

  console.log(`contract deployed at ${await processorEndpoint.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
