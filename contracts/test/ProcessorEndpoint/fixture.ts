import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { BYTES_ZERO, BYTES32_ZERO } from '../util';
import { ethSignStateUpdate } from '../../scripts/util';

export const INITIAL_STATE_ROOT =
  '0x0000000000000000000000000000000000000000000000000000000000000001';

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
  const teeAuthenticator = await MockTeeAuthenticator.deploy(
    await signers[0].getAddress(),
    BYTES_ZERO
  );

  const RequestQueueLib = await ethers.getContractFactory('RequestQueueLib');
  const requestQueueLib = await RequestQueueLib.deploy();

  const processorEndpointFactory = await ethers.getContractFactory('ProcessorEndpoint', {
    libraries: { RequestQueueLib: await requestQueueLib.getAddress() },
  });
  const resetOperator = await signers[3].getAddress();

  async function deployTokenAllowlist() {
    const TokenAllowlist = await ethers.getContractFactory('TokenAllowlist');
    return TokenAllowlist.deploy(admin);
  }

  const sharedTokenAllowlist = await deployTokenAllowlist();

  async function deployProcessorEndpoint(resetOperatorOverride?: string) {
    const tokenAllowlist = await deployTokenAllowlist();
    const processorEndpoint = await processorEndpointFactory.deploy(
      await teeAuthenticator.getAddress(),
      await authorityRegistry.getAddress(),
      updateStatusOperator,
      admin,
      resetOperatorOverride !== undefined ? resetOperatorOverride : resetOperator,
      minFeePerRequest,
      await tokenAllowlist.getAddress()
    );
    return { processorEndpoint, tokenAllowlist };
  }

  async function bootstrapApplication(processorEndpoint: any, teeSigner?: Signer) {
    const deployTx = await processorEndpoint
      .connect(signers[2])
      .submitDeployRequest(0, '0x00', { value: minFeePerRequest });
    const deployReceipt = await deployTx.wait();

    const deployLog = deployReceipt.logs.find((log: any) => {
      try {
        return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
      } catch {
        return false;
      }
    });
    const parsed = processorEndpoint.interface.parseLog(deployLog);
    const applicationId: bigint = parsed.args.applicationId;
    const requestId: string = parsed.args.requestId;

    let signature = '0x';
    if (teeSigner) {
      signature = await ethSignStateUpdate(
        teeSigner,
        applicationId,
        BYTES32_ZERO,
        INITIAL_STATE_ROOT,
        requestId,
        [],
        [],
        [],
        [],
        [],
        0,
        minFeePerRequest
      );
    }

    await processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        applicationId,
        BYTES32_ZERO,
        INITIAL_STATE_ROOT,
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        minFeePerRequest,
        0,
        '',
        signature
      );

    return { applicationId };
  }

  return {
    signers,
    defaultAuthority,
    authorityRegistry,
    teeAuthenticator,
    processorEndpointFactory,
    updateStatusOperator,
    admin,
    resetOperator,
    minFeePerRequest,
    deployProcessorEndpoint,
    deployTokenAllowlist,
    sharedTokenAllowlist,
    bootstrapApplication,
  };
}
