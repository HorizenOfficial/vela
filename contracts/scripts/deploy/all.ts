import { ethers } from 'hardhat';
import CertManagerArtifact from '../../nitro_prover/CertManager.json';
import NitroProverArtifact from '../../nitro_prover/NitroProver.json';

/*
  Deploy scripts for all the contracts.
  Needed parameters are:
  TEE_MAX_VERIFICATION_AGE: tolerance from expiration of the attestation to be still considered valid
  TEE_PCR0: PCR0 of the Nitro application
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
  //deploy cert manager
  const CertManagerFactory = new ethers.ContractFactory(CertManagerArtifact.abi, CertManagerArtifact.bytecode, deployer);
  const certManager = await CertManagerFactory.deploy();
  await certManager.deploymentTransaction()!.wait();
  const certManagerAddress = await certManager.getAddress();
  console.log(`CertManager`);
  console.log(`  contract address: ${certManagerAddress}`);


  //deploy nitro prover
  const NitroProverFactory = new ethers.ContractFactory(NitroProverArtifact.abi, NitroProverArtifact.bytecode, deployer);
  const nitroProver = await NitroProverFactory.deploy(certManagerAddress);
  await nitroProver.deploymentTransaction()!.wait();
  const nitroProverAddress = await nitroProver.getAddress();
  console.log(`NitroProver`);
  console.log(`  contract address: ${nitroProverAddress}`);


  //deploy TeeAuthenticator
  const TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
  const teeAuthenticator = await TeeAuthenticator.deploy(deployerAddress, nitroProverAddress, process.env.TEE_PCR0!, process.env.TEE_MAX_VERIFICATION_AGE!);
  await teeAuthenticator.deploymentTransaction()!.wait();
  const teeAuthenticatorAddr = await teeAuthenticator.getAddress();

  console.log(`TeeAuthenticator`);
  console.log(`  contract address: ${teeAuthenticatorAddr}`);

  // 4) ProcessorEndpoint
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
