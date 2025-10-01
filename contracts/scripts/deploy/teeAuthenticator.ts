import { ethers } from 'hardhat';

async function deploy()  {

  const deployer = (await ethers.getSigners())[0]

  console.log(`deploying from ${await deployer.getAddress()}`)
  console.log(`parameters:
    owner: ${process.env.TEE_OWNER},
    _teeSigner: ${process.env.TEE_SIGNER},
    _pubSecp521r1: ${process.env.TEE_PUB_SECP521R1}
  `)
  //deploy 
  const TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
  const teeAuthenticator = await TeeAuthenticator.deploy(process.env.TEE_OWNER!, process.env.TEE_SIGNER!, process.env.TEE_PUB_SECP521R1!);
  await teeAuthenticator.deploymentTransaction()!.wait();

  console.log(`contract deployed at ${await teeAuthenticator.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
