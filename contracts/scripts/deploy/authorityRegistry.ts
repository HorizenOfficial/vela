import { ethers } from 'hardhat';

async function deploy()  {

  const deployer = (await ethers.getSigners())[0]

  console.log(`deploying from ${await deployer.getAddress()}`)
  console.log(`parameters:
    owner: ${process.env.AUTHREGISTRY_OWNER}
  `)
  //deploy 
  const AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
  const authorityRegistry = await AuthorityRegistry.deploy(process.env.AUTHREGISTRY_OWNER);
  await authorityRegistry.deploymentTransaction().wait();

  console.log(`contract deployed at ${await authorityRegistry.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
