import { ethers } from 'hardhat';

async function deploy()  {

  const deployer = (await ethers.getSigners())[0]

  console.log(`deploying from ${await deployer.getAddress()}`)
  //deploy 
  const KeyRegistry = await ethers.getContractFactory("KeyRegistry");
  const keyRegistry = await KeyRegistry.deploy();
  await keyRegistry.deploymentTransaction()!.wait();

  console.log(`contract deployed at ${await keyRegistry.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
