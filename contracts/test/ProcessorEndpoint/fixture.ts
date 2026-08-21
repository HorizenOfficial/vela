import { Signer } from 'ethers';
import { ethers, upgrades } from 'hardhat';
import { BYTES_ZERO, BYTES32_ZERO } from '../util';
import { ethSignStateUpdate, PROCESSOR_ENDPOINT_UPGRADE_OPTIONS } from '../../scripts/util';

export const INITIAL_STATE_ROOT =
  '0x0000000000000000000000000000000000000000000000000000000000000001';

/// The `DeployRequestSubmitted` args of a deploy submission: the derived `applicationId` and the
/// `requestId`. Both deploy entry points emit it, and it is the only place the derived
/// applicationId is observable.
export function deployRequestArgs(processorEndpoint: any, receipt: any) {
  const log = receipt.logs.find((entry: any) => {
    try {
      return processorEndpoint.interface.parseLog(entry)?.name === 'DeployRequestSubmitted';
    } catch {
      return false;
    }
  });
  return processorEndpoint.interface.parseLog(log).args;
}

export async function deployProcessorEndpointFixture() {
  const signers = await ethers.getSigners();
  const updateStatusOperator = await signers[1].getAddress();
  const admin = await signers[2].getAddress();
  const minFeePerRequest = 100n;

  const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
  const defaultAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

  const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
  const authorityRegistry = await upgrades.deployProxy(
    AuthorityRegistry,
    [await signers[0].getAddress(), await defaultAuthority.getAddress()],
    { kind: 'uups' }
  );

  const MockTeeAuthenticator = await ethers.getContractFactory('MockTeeAuthenticator');
  const teeAuthenticator = await MockTeeAuthenticator.deploy(
    await signers[0].getAddress(),
    BYTES_ZERO
  );

  const processorEndpointFactory = await ethers.getContractFactory('ProcessorEndpoint');
  const resetOperator = await signers[3].getAddress();

  // Code moved out of ProcessorEndpoint to stay under EIP-170, reached by delegatecall. Stateless
  // and deployed once: every endpoint in this fixture can share the same instance.
  const ProcessorEndpointExtension = await ethers.getContractFactory('ProcessorEndpointExtension');
  const extension = await ProcessorEndpointExtension.deploy();
  const extensionAddress = await extension.getAddress();

  async function deployTokenAllowlist() {
    const TokenAllowlist = await ethers.getContractFactory('TokenAllowlist');
    return TokenAllowlist.deploy(admin);
  }

  const sharedTokenAllowlist = await deployTokenAllowlist();

  // Deployed behind a UUPS proxy (docs/design/UPGRADABLE_CONTRACTS_DESIGN.md). `extensionAddress`
  // is a constructor argument (it stays `immutable`, not proxy state), so it goes in
  // `constructorArgs`, not the initializer args. Exposed so tests that need a differently
  // configured endpoint (e.g. a different TeeAuthenticator) can deploy one directly.
  async function deployProcessorEndpointWith(
    teeAuthenticatorAddress: string,
    authorityRegistryAddress: string,
    tokenAllowlistAddress: string,
    resetOperatorOverride?: string
  ) {
    return upgrades.deployProxy(
      processorEndpointFactory,
      [
        teeAuthenticatorAddress,
        authorityRegistryAddress,
        updateStatusOperator,
        admin,
        resetOperatorOverride !== undefined ? resetOperatorOverride : resetOperator,
        minFeePerRequest,
        tokenAllowlistAddress,
      ],
      {
        kind: 'uups',
        constructorArgs: [extensionAddress],
        unsafeAllow: [...PROCESSOR_ENDPOINT_UPGRADE_OPTIONS.unsafeAllow],
      }
    );
  }

  async function deployProcessorEndpoint(resetOperatorOverride?: string) {
    const tokenAllowlist = await deployTokenAllowlist();
    const processorEndpoint = await deployProcessorEndpointWith(
      await teeAuthenticator.getAddress(),
      await authorityRegistry.getAddress(),
      await tokenAllowlist.getAddress(),
      resetOperatorOverride
    );
    return { processorEndpoint, tokenAllowlist, extension };
  }

  // Completes a pending DEPLOYAPP request, moving the new application's state root from zero to
  // INITIAL_STATE_ROOT. `teeSigner` is only needed when the endpoint verifies signatures for real:
  // MockTeeAuthenticator accepts anything, including the empty signature used otherwise.
  async function completeDeploy(
    processorEndpoint: any,
    applicationId: bigint,
    requestId: string,
    teeSigner?: Signer
  ) {
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
  }

  async function bootstrapApplication(processorEndpoint: any, teeSigner?: Signer) {
    const deployTx = await processorEndpoint
      .connect(signers[2])
      .submitDeployRequest(0, '0x00', { value: minFeePerRequest });
    const { applicationId, requestId } = deployRequestArgs(
      processorEndpoint,
      await deployTx.wait()
    );

    await completeDeploy(processorEndpoint, applicationId, requestId, teeSigner);

    return { applicationId };
  }

  /// Deploys a `TestTrigger`, submits a deploy request registering it, and completes that deploy so
  /// the application exists. Processing one of its requests afterwards runs `_invokeTrigger`, and
  /// enqueues a TRUSTPROCESS into the global trigger queue once the trigger has a trusted payload
  /// set (`trigger.setTrustedPayload`).
  ///
  /// The default descriptor payload is address-shaped but is **not** the trigger's address: the
  /// deploy payload is opaque and must never be interpreted as a trigger address, which the trigger
  /// being registered through its own argument is what guarantees.
  async function bootstrapApplicationWithTrigger(
    processorEndpoint: any,
    options: {
      revertOnExecute?: boolean;
      revertOnPostWithdraw?: boolean;
      descriptorPayload?: string;
      teeSigner?: Signer;
    } = {}
  ) {
    const TestTrigger = await ethers.getContractFactory('TestTrigger');
    const trigger: any = await TestTrigger.deploy(
      await processorEndpoint.getAddress(),
      options.revertOnExecute ?? false,
      options.revertOnPostWithdraw ?? false
    );

    const descriptorPayload = options.descriptorPayload ?? '0x' + 'ab'.repeat(40);
    const deployTx = await processorEndpoint
      .connect(signers[2])
      .submitDeployRequestWithTrigger(0, descriptorPayload, await trigger.getAddress(), {
        value: minFeePerRequest,
      });
    const { applicationId, requestId } = deployRequestArgs(
      processorEndpoint,
      await deployTx.wait()
    );

    await completeDeploy(processorEndpoint, applicationId, requestId, options.teeSigner);

    return { trigger, applicationId, requestId, descriptorPayload };
  }

  return {
    signers,
    defaultAuthority,
    authorityRegistry,
    teeAuthenticator,
    processorEndpointFactory,
    extensionAddress,
    updateStatusOperator,
    admin,
    resetOperator,
    minFeePerRequest,
    deployProcessorEndpoint,
    deployProcessorEndpointWith,
    deployTokenAllowlist,
    sharedTokenAllowlist,
    bootstrapApplication,
    bootstrapApplicationWithTrigger,
  };
}
