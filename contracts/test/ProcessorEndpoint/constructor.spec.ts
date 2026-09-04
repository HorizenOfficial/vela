import { expect } from 'chai';
import { upgrades } from 'hardhat';
import { ADDRESS_ZERO } from '../util';
import { deployProcessorEndpointFixture } from './fixture';

// ProcessorEndpoint is deployed behind a UUPS proxy (docs/design/UPGRADABLE_CONTRACTS_DESIGN.md).
// The extension address stays an immutable constructor argument (see the `_extension` doc comment
// on ProcessorEndpoint), so its validation is still exercised via the raw constructor; every other
// former-constructor check moved into `initialize`, exercised through `upgrades.deployProxy`.
describe('ProcessorEndpoint Test', function () {
  let authorityRegistry: any;
  let teeAuthenticator: any;
  let processorEndpointFactory: any;
  let updateStatusOperator: string;
  let admin: string;
  let resetOperator: string;
  let minFeePerRequest: bigint;
  let deployProcessorEndpoint: (resetOperatorOverride?: string) => Promise<any>;
  let deployProcessorEndpointWith: (
    teeAuthenticatorAddress: string,
    authorityRegistryAddress: string,
    tokenAllowlistAddress: string,
    resetOperatorOverride?: string
  ) => Promise<any>;
  let sharedTokenAllowlist: any;
  let extensionAddress: string;

  beforeEach(async function () {
    ({
      authorityRegistry,
      teeAuthenticator,
      processorEndpointFactory,
      updateStatusOperator,
      admin,
      resetOperator,
      minFeePerRequest,
      deployProcessorEndpoint,
      deployProcessorEndpointWith,
      sharedTokenAllowlist,
      extensionAddress,
    } = await deployProcessorEndpointFixture());
  });

  describe('constructor', function () {
    describe('unhappy paths', function () {
      it('reverts when extension is zero address', async () => {
        await expect(processorEndpointFactory.deploy(ADDRESS_ZERO)).to.be.revertedWithCustomError(
          processorEndpointFactory,
          'AddressCantBeZero'
        );
      });

      // delegatecall to an address without code succeeds and returns nothing, so a wrong
      // extension address would make submitRequestFor a silent no-op that keeps the fee. The
      // address is immutable, so it cannot be corrected after deployment.
      it('reverts when extension has no code', async () => {
        await expect(
          processorEndpointFactory.deploy(updateStatusOperator)
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'InvalidExtension');
      });
    });
  });

  describe('initialize', function () {
    describe('unhappy paths', function () {
      it('reverts when teeAuthenticator is zero address', async () => {
        await expect(
          deployProcessorEndpointWith(
            ADDRESS_ZERO,
            await authorityRegistry.getAddress(),
            await sharedTokenAllowlist.getAddress()
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when authorityRegistry is zero address', async () => {
        await expect(
          deployProcessorEndpointWith(
            await teeAuthenticator.getAddress(),
            ADDRESS_ZERO,
            await sharedTokenAllowlist.getAddress()
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when updateStatusOperator is zero address', async () => {
        await expect(
          upgrades.deployProxy(
            processorEndpointFactory,
            [
              await teeAuthenticator.getAddress(),
              await authorityRegistry.getAddress(),
              ADDRESS_ZERO,
              admin,
              ADDRESS_ZERO,
              minFeePerRequest,
              await sharedTokenAllowlist.getAddress(),
            ],
            {
              kind: 'uups',
              constructorArgs: [extensionAddress],
            }
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when admin is zero address', async () => {
        await expect(
          upgrades.deployProxy(
            processorEndpointFactory,
            [
              await teeAuthenticator.getAddress(),
              await authorityRegistry.getAddress(),
              updateStatusOperator,
              ADDRESS_ZERO,
              ADDRESS_ZERO,
              minFeePerRequest,
              await sharedTokenAllowlist.getAddress(),
            ],
            {
              kind: 'uups',
              constructorArgs: [extensionAddress],
            }
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when tokenAllowlist is zero address', async () => {
        await expect(
          deployProcessorEndpointWith(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            ADDRESS_ZERO
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts with InvalidInitialization when initialize is called a second time', async () => {
        const { processorEndpoint } = await deployProcessorEndpoint();
        await expect(
          processorEndpoint.initialize(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            admin,
            resetOperator,
            minFeePerRequest,
            await sharedTokenAllowlist.getAddress()
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'InvalidInitialization');
      });
    });

    describe('happy paths', function () {
      it('initializes dependencies, roles, and config values', async () => {
        const { processorEndpoint } = await deployProcessorEndpoint();

        expect(await processorEndpoint.teeAuthenticator()).to.equal(
          await teeAuthenticator.getAddress()
        );
        expect(await processorEndpoint.authorityRegistry()).to.equal(
          await authorityRegistry.getAddress()
        );
        expect(await processorEndpoint.feeCollector()).to.equal(updateStatusOperator);
        expect(await processorEndpoint.minFeePerRequest()).to.equal(minFeePerRequest);
        expect(await processorEndpoint.maxQueueSize()).to.equal(10n);
        expect(await processorEndpoint.availableDeploySlots()).to.equal(
          await processorEndpoint.maxNumOfApplications()
        );
        const updateRole = await processorEndpoint.UPDATE_STATUS_ROLE();
        const adminRole = await processorEndpoint.ADMIN();
        const deployerRole = await processorEndpoint.DEPLOYER_ROLE();
        const resetRole = await processorEndpoint.RESET_OPERATOR();
        expect(await processorEndpoint.hasRole(updateRole, updateStatusOperator)).to.equal(true);
        expect(await processorEndpoint.hasRole(adminRole, admin)).to.equal(true);
        expect(await processorEndpoint.hasRole(deployerRole, admin)).to.equal(true);
        expect(await processorEndpoint.getRoleAdmin(deployerRole)).to.equal(adminRole);
        expect(await processorEndpoint.hasRole(resetRole, resetOperator)).to.equal(true);
      });

      it('sets tokenAllowlist address correctly', async () => {
        const { processorEndpoint, tokenAllowlist } = await deployProcessorEndpoint();
        expect(await processorEndpoint.tokenAllowlist()).to.equal(
          await tokenAllowlist.getAddress()
        );
      });

      // The pairing between an endpoint and the extension it delegates to must be readable
      // on-chain: it is fixed at deployment and cannot be repointed, so this getter is what
      // verification has to go on.
      it('exposes the extension it delegates to', async () => {
        const { processorEndpoint, extension } = await deployProcessorEndpoint();
        expect(await processorEndpoint.extension()).to.equal(await extension.getAddress());
      });

      it('does not grant RESET_OPERATOR when address(0) is passed', async () => {
        const { processorEndpoint } = await deployProcessorEndpoint(ADDRESS_ZERO);
        const resetRole = await processorEndpoint.RESET_OPERATOR();
        expect(await processorEndpoint.hasRole(resetRole, ADDRESS_ZERO)).to.equal(false);
      });
    });
  });
});
