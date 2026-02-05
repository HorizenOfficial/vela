import { ethers } from 'hardhat';
import CertManagerArtifact from '../../nitro_prover/CertManager.json';
import NitroProverArtifact from '../../nitro_prover/NitroProver.json';

async function deploy() {
  const [deployer] = await ethers.getSigners();
  console.log(`deploying from ${await deployer.getAddress()}`);

  //deploy cert manager
  const CertManagerFactory = new ethers.ContractFactory(
    CertManagerArtifact.abi,
    CertManagerArtifact.bytecode,
    deployer
  );
  const certManager = await CertManagerFactory.deploy();
  await certManager.deploymentTransaction()!.wait();
  const certManagerAddress = await certManager.getAddress();
  console.log(`CertManager deployed at ${certManagerAddress}`);

  //deploy nitro validator
  const NitroProverFactory = new ethers.ContractFactory(
    NitroProverArtifact.abi,
    NitroProverArtifact.bytecode,
    deployer
  );
  const nitroProver = await NitroProverFactory.deploy(certManagerAddress);
  await nitroProver.deploymentTransaction()!.wait();
  const nitroProverAddress = await nitroProver.getAddress();
  console.log(`NitroProver deployed at ${nitroProverAddress}`);
  console.log(`TeeAuthenticator parameters:
    owner: ${process.env.TEE_OWNER},
    _nitroProver: ${nitroProverAddress},
    _pcr0: ${process.env.TEE_PCR0},
    _maxVerificationAge: ${process.env.TEE_MAX_VERIFICATION_AGE}
  `);
  //deploy TeeAuthenticator
  const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
  const teeAuthenticator = await TeeAuthenticator.deploy(
    process.env.TEE_OWNER!,
    nitroProverAddress,
    process.env.TEE_PCR0!,
    process.env.TEE_MAX_VERIFICATION_AGE!
  );
  await teeAuthenticator.deploymentTransaction()!.wait();

  console.log(`TeeAuthenticator deployed at ${await teeAuthenticator.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
