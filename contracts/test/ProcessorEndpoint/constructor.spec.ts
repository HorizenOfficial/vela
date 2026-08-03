import { expect } from 'chai';
import { ADDRESS_ZERO, BYTES32_ZERO } from '../util';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint Test', function () {
  let authorityRegistry: any;
  let teeAuthenticator: any;
  let processorEndpointFactory: any;
  let updateStatusOperator: string;
  let admin: string;
  let resetOperator: string;
  let minFeePerRequest: bigint;
  let deployProcessorEndpoint: (resetOperatorOverride?: string) => Promise<any>;
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
      sharedTokenAllowlist,
      extensionAddress,
    } = await deployProcessorEndpointFixture());
  });

  describe('constructor', function () {
    describe('unhappy paths', function () {
      it('reverts when teeAuthenticator is zero address', async () => {
        await expect(
          processorEndpointFactory.deploy(
            ADDRESS_ZERO,
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            admin,
            ADDRESS_ZERO,
            minFeePerRequest,
            await sharedTokenAllowlist.getAddress(),
            extensionAddress
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when authorityRegistry is zero address', async () => {
        await expect(
          processorEndpointFactory.deploy(
            await teeAuthenticator.getAddress(),
            ADDRESS_ZERO,
            updateStatusOperator,
            admin,
            ADDRESS_ZERO,
            minFeePerRequest,
            await sharedTokenAllowlist.getAddress(),
            extensionAddress
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when updateStatusOperator is zero address', async () => {
        await expect(
          processorEndpointFactory.deploy(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            ADDRESS_ZERO,
            admin,
            ADDRESS_ZERO,
            minFeePerRequest,
            await sharedTokenAllowlist.getAddress(),
            extensionAddress
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when admin is zero address', async () => {
        await expect(
          processorEndpointFactory.deploy(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            ADDRESS_ZERO,
            ADDRESS_ZERO,
            minFeePerRequest,
            await sharedTokenAllowlist.getAddress(),
            extensionAddress
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when tokenAllowlist is zero address', async () => {
        await expect(
          processorEndpointFactory.deploy(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            admin,
            ADDRESS_ZERO,
            minFeePerRequest,
            ADDRESS_ZERO,
            extensionAddress
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when extension is zero address', async () => {
        await expect(
          processorEndpointFactory.deploy(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            admin,
            ADDRESS_ZERO,
            minFeePerRequest,
            await sharedTokenAllowlist.getAddress(),
            ADDRESS_ZERO
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      // delegatecall to an address without code succeeds and returns nothing, so a wrong
      // extension address would make submitRequestFor a silent no-op that keeps the fee. The
      // address is immutable, so it cannot be corrected after deployment.
      it('reverts when extension has no code', async () => {
        await expect(
          processorEndpointFactory.deploy(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            admin,
            ADDRESS_ZERO,
            minFeePerRequest,
            await sharedTokenAllowlist.getAddress(),
            updateStatusOperator
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'InvalidExtension');
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

      it('does not grant RESET_OPERATOR when address(0) is passed', async () => {
        const { processorEndpoint } = await deployProcessorEndpoint(ADDRESS_ZERO);
        const resetRole = await processorEndpoint.RESET_OPERATOR();
        expect(await processorEndpoint.hasRole(resetRole, ADDRESS_ZERO)).to.equal(false);
      });
    });
  });
});
