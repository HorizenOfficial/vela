import { expect } from 'chai';
import { ADDRESS_ZERO, BYTES32_ZERO } from '../util';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint Test', function () {
  let authorityRegistry: any;
  let teeAuthenticator: any;
  let processorEndpointFactory: any;
  let updateStatusOperator: string;
  let admin: string;
  let minFeePerRequest: bigint;
  let deployProcessorEndpoint: () => Promise<any>;

  beforeEach(async function () {
    ({
      authorityRegistry,
      teeAuthenticator,
      processorEndpointFactory,
      updateStatusOperator,
      admin,
      minFeePerRequest,
      deployProcessorEndpoint,
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
            minFeePerRequest
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
            minFeePerRequest
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
            minFeePerRequest
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
            minFeePerRequest
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });
    });

    describe('happy paths', function () {
      it('initializes dependencies, roles, and config values', async () => {
        const processorEndpoint = await deployProcessorEndpoint();

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
        expect(await processorEndpoint.hasRole(updateRole, updateStatusOperator)).to.equal(true);
        expect(await processorEndpoint.hasRole(adminRole, admin)).to.equal(true);
        expect(await processorEndpoint.hasRole(deployerRole, admin)).to.equal(true);
        expect(await processorEndpoint.getRoleAdmin(deployerRole)).to.equal(adminRole);
      });
    });
  });
});
