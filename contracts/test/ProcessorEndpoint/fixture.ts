import { ethers } from 'hardhat';
import { BYTES_ZERO } from '../util';

export async function deployProcessorEndpointFixture() {
  const signers = await ethers.getSigners();
  const updateStatusOperator = await signers[1].getAddress();
  const admin = await signers[2].getAddress();
  const minFeePerRequest = 100n;

  const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
  const defaultAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

  const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
  const authorityRegistry = await AuthorityRegistry.deploy(
    await signers[0].getAddress(),
    await defaultAuthority.getAddress()
  );

  const MockTeeAuthenticator = await ethers.getContractFactory('MockTeeAuthenticator');
  const teeAuthenticator = await MockTeeAuthenticator.deploy(await signers[0].getAddress(), BYTES_ZERO);

  const processorEndpointFactory = await ethers.getContractFactory('ProcessorEndpoint');

  async function deployProcessorEndpoint() {
    return processorEndpointFactory.deploy(
      await teeAuthenticator.getAddress(),
      await authorityRegistry.getAddress(),
      updateStatusOperator,
      admin,
      minFeePerRequest
    );
  }

  return {
    signers,
    defaultAuthority,
    authorityRegistry,
    teeAuthenticator,
    processorEndpointFactory,
    updateStatusOperator,
    admin,
    minFeePerRequest,
    deployProcessorEndpoint,
  };
}
