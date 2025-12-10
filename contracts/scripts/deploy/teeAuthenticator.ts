import { ethers } from 'hardhat';
import CertManagerArtifact from '../../nitro-validator/CertManager.json';
import NitroValidatorArtifact from '../../nitro-validator/Validator.json';

async function deploy()  {

  const [deployer] = await ethers.getSigners();
  console.log(`deploying from ${await deployer.getAddress()}`)

  //deploy cert manager
  const CertManagerFactory = new ethers.ContractFactory(CertManagerArtifact.abi, CertManagerArtifact.bytecode, deployer);
  const certManager = await CertManagerFactory.deploy();
  await certManager.deploymentTransaction()!.wait();
  const certManagerAddress = await certManager.getAddress();
  console.log(`CertManager deployed at ${certManagerAddress}`);

  //deploy nitro validator
  const NitroValidatorFactory = new ethers.ContractFactory(NitroValidatorArtifact.abi, NitroValidatorArtifact.bytecode, deployer);
  const nitroValidator = await NitroValidatorFactory.deploy(certManagerAddress);
  await nitroValidator.deploymentTransaction()!.wait();
  const nitroValidatorAddress = await nitroValidator.getAddress();
  console.log(`NitroValidator deployed at ${nitroValidatorAddress}`);
  console.log(`TeeAuthenticator parameters:
    owner: ${process.env.TEE_OWNER},
    _nitroValidator: ${nitroValidatorAddress},
    _pcr0: ${process.env.TEE_PCR0},
    _maxVerificationAge: ${process.env.TEE_MAX_VERIFICATION_AGE}
  `)
  //deploy TeeAuthenticator
  const TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
  const teeAuthenticator = await TeeAuthenticator.deploy(process.env.TEE_OWNER!, nitroValidatorAddress, process.env.TEE_PCR0!, process.env.TEE_MAX_VERIFICATION_AGE!);
  await teeAuthenticator.deploymentTransaction()!.wait();

  console.log(`TeeAuthenticator deployed at ${await teeAuthenticator.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
