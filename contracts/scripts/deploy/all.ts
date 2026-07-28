import { ethers } from 'hardhat';
import * as fs from 'fs';
import * as path from 'path';

/*
  Deploy scripts for all the contracts.
  Needed parameters are:
  UPDATE_STATUS_OPERATOR: ethereum address of the processor endpoint status updater (manager address)

  For full TEE attestation (default):
    TEE_MAX_VERIFICATION_AGE: tolerance from expiration of the attestation to be still considered valid
    TEE_PCR0: PCR0 of the Nitro application
    TEE_PCR0_UPGRADE_DELAY: timelock (seconds) before a swap to a new PCR0 can be applied

  For no-attestation mode (TEE_NO_ATTESTATION=true):
    TEE_SIGNER_ADDRESS: ethereum address of the executor's signing key
    TEE_PUB_P521: hex-encoded uncompressed P521 public key (133 bytes)
*/

async function deployTeeAuthenticator(deployer: any, deployerAddress: string): Promise<string> {
  if (process.env.TEE_NO_ATTESTATION === 'true') {
    // Deploy NoAttestationTeeAuthenticator (dev/test mode)
    const teeSignerAddress = process.env.TEE_SIGNER_ADDRESS!;
    const teePubP521 = process.env.TEE_PUB_P521!;
    console.log(`Using NoAttestationTeeAuthenticator (dev mode)`);

    const NoAttestationTeeAuthenticator = await ethers.getContractFactory(
      'NoAttestationTeeAuthenticator'
    );
    const teeAuthenticator = await NoAttestationTeeAuthenticator.deploy(
      deployerAddress,
      teeSignerAddress,
      teePubP521
    );
    await teeAuthenticator.deploymentTransaction()!.wait();
    const addr = await teeAuthenticator.getAddress();

    console.log(`NoAttestationTeeAuthenticator`);
    console.log(`  contract address: ${addr}`);
    console.log(`  teeSigner: ${teeSignerAddress}`);
    return addr;
  }

  // Deploy full TeeAuthenticator with Nitro attestation
  const CertManagerArtifact = require('../../nitro_prover/CertManager.json');
  const NitroProverArtifact = require('../../nitro_prover/NitroProver.json');

  const CertManagerFactory = new ethers.ContractFactory(
    CertManagerArtifact.abi,
    CertManagerArtifact.bytecode,
    deployer
  );
  const certManager = await CertManagerFactory.deploy();
  await certManager.deploymentTransaction()!.wait();
  const certManagerAddress = await certManager.getAddress();
  console.log(`CertManager`);
  console.log(`  contract address: ${certManagerAddress}`);

  const NitroProverFactory = new ethers.ContractFactory(
    NitroProverArtifact.abi,
    NitroProverArtifact.bytecode,
    deployer
  );
  const nitroProver = await NitroProverFactory.deploy(certManagerAddress);
  await nitroProver.deploymentTransaction()!.wait();
  const nitroProverAddress = await nitroProver.getAddress();
  console.log(`NitroProver`);
  console.log(`  contract address: ${nitroProverAddress}`);

  const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
  const teeAuthenticator = await TeeAuthenticator.deploy(
    deployerAddress,
    nitroProverAddress,
    process.env.TEE_PCR0!,
    process.env.TEE_MAX_VERIFICATION_AGE!,
    process.env.TEE_PCR0_UPGRADE_DELAY!
  );
  await teeAuthenticator.deploymentTransaction()!.wait();
  const addr = await teeAuthenticator.getAddress();
  console.log(`TeeAuthenticator`);
  console.log(`  contract address: ${addr}`);
  return addr;
}

async function deploy() {
  const deployer = (await ethers.getSigners())[0];
  const deployerAddress = await deployer.getAddress();

  console.log(`deploying all  contracts from ${deployerAddress}`);
  console.log(`(this address will be also the owner of the contracts)`);

  // 1) DefaultAuthority
  const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
  const defaultAuthority = await DefaultAuthority.deploy(deployerAddress);
  await defaultAuthority.deploymentTransaction()!.wait();
  const defaultAuthorityAddr = await defaultAuthority.getAddress();
  console.log(`DefaultAuthority`);
  console.log(`  contract address: ${defaultAuthorityAddr}`);

  // 2) AuthorityRegistry (proxy) usando el default
  const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
  const authorityRegistry = await AuthorityRegistry.deploy(deployerAddress, defaultAuthorityAddr);
  await authorityRegistry.deploymentTransaction()!.wait();
  const authorityRegistryAddr = await authorityRegistry.getAddress();
  console.log(`AuthorityRegistry`);
  console.log(`  contract address: ${authorityRegistryAddr}`);
  console.log(`  default authority contract: ${defaultAuthorityAddr}`);

  // 3) TeeAuthenticator
  const teeAuthenticatorAddr = await deployTeeAuthenticator(deployer, deployerAddress);

  // 4) TokenAllowlist
  const TokenAllowlist = await ethers.getContractFactory('TokenAllowlist');
  const tokenAllowlist = await TokenAllowlist.deploy(process.env.ADMIN!);
  await tokenAllowlist.deploymentTransaction()!.wait();
  const tokenAllowlistAddr = await tokenAllowlist.getAddress();
  console.log(`TokenAllowlist`);
  console.log(`  contract address: ${tokenAllowlistAddr}`);

  // 5) ProcessorEndpoint
  const resetOperator = process.env.RESET_OPERATOR || '0x0000000000000000000000000000000000000000';
  const ProcessorEndpoint = await ethers.getContractFactory('ProcessorEndpoint');
  const processorEndpoint = await ProcessorEndpoint.deploy(
    teeAuthenticatorAddr,
    authorityRegistryAddr,
    process.env.UPDATE_STATUS_OPERATOR!,
    process.env.ADMIN!,
    resetOperator,
    process.env.MIN_FEE_PER_REQUEST!,
    tokenAllowlistAddr
  );
  await processorEndpoint.deploymentTransaction()!.wait();
  var processorEndpointAddr = await processorEndpoint.getAddress();
  console.log(`ProcessorEndpoint`);
  console.log(`  contract address: ${processorEndpointAddr}`);
  console.log(`  update status operator (manager address): ${process.env.UPDATE_STATUS_OPERATOR!}`);
  console.log(`  admin / bootstrap deployer: ${process.env.ADMIN!}`);
  console.log(`  reset operator (address(0) = disabled): ${resetOperator}`);
  console.log(`  token allowlist: ${tokenAllowlistAddr}`);

  // 5) Optional MockERC20 for dev/test. Opt-in via DEPLOY_MOCK_ERC20=true.
  //    Must be deployed AFTER ProcessorEndpoint so that its CREATE address
  //    stays stable across runs (deterministic on deployer nonce).
  //    Requires the deployer to hold the ADMIN role on ProcessorEndpoint,
  //    which is the case when DEPLOYER_PRIVATE_KEY corresponds to ADMIN.
  let mockErc20Addr = '';
  if (process.env.DEPLOY_MOCK_ERC20 === 'true') {
    const name = process.env.MOCK_ERC20_NAME || 'MockToken';
    const symbol = process.env.MOCK_ERC20_SYMBOL || 'MOCK';
    const decimals = parseInt(process.env.MOCK_ERC20_DECIMALS || '18', 10);

    const MockERC20 = await ethers.getContractFactory('MockERC20');
    const mockErc20 = await MockERC20.deploy(name, symbol, decimals);
    await mockErc20.deploymentTransaction()!.wait();
    mockErc20Addr = await mockErc20.getAddress();
    console.log(`MockERC20`);
    console.log(`  contract address: ${mockErc20Addr}`);
    console.log(`  name / symbol / decimals: ${name} / ${symbol} / ${decimals}`);

    // Allowlist the MockERC20 on TokenAllowlist.
    // The deployer signer must have the ADMIN role (granted to process.env.ADMIN
    // in the TokenAllowlist constructor).
    const tx = await tokenAllowlist.addAllowedToken(mockErc20Addr);
    await tx.wait();
    console.log(`  allowlisted on TokenAllowlist`);
  }

  // Write deployed addresses to file if DEPLOY_OUTPUT_DIR is set
  const outputDir = process.env.DEPLOY_OUTPUT_DIR;
  if (outputDir) {
    const lines = [
      `CHAIN_PROCESSOR_ADDRESS=${processorEndpointAddr}`,
      `CHAIN_TEEAUTHENTICATOR_ADDRESS=${teeAuthenticatorAddr}`,
      `CHAIN_TOKEN_ALLOWLIST_ADDRESS=${tokenAllowlistAddr}`,
    ];
    if (mockErc20Addr !== '') {
      lines.push(`CHAIN_MOCK_ERC20_ADDRESS=${mockErc20Addr}`);
    }
    lines.push('');
    const envContent = lines.join('\n');
    fs.mkdirSync(outputDir, { recursive: true });
    fs.writeFileSync(path.join(outputDir, 'deployed_addresses.env'), envContent);
    console.log(
      `\nDeployed addresses written to ${path.join(outputDir, 'deployed_addresses.env')}`
    );
  }
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
