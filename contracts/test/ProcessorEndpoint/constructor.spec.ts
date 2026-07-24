import { expect } from 'chai';
import { ethers, upgrades } from 'hardhat';
import { ADDRESS_ZERO } from '../util';
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
    } = await deployProcessorEndpointFixture());
  });

  describe('initialize', function () {
    describe('unhappy paths', function () {
      it('reverts when teeAuthenticator is zero address', async () => {
        await expect(
          upgrades.deployProxy(
            processorEndpointFactory,
            [
              ADDRESS_ZERO,
              await authorityRegistry.getAddress(),
              updateStatusOperator,
              admin,
              ADDRESS_ZERO,
              minFeePerRequest,
              await sharedTokenAllowlist.getAddress(),
            ],
            { kind: 'uups' }
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when authorityRegistry is zero address', async () => {
        await expect(
          upgrades.deployProxy(
            processorEndpointFactory,
            [
              await teeAuthenticator.getAddress(),
              ADDRESS_ZERO,
              updateStatusOperator,
              admin,
              ADDRESS_ZERO,
              minFeePerRequest,
              await sharedTokenAllowlist.getAddress(),
            ],
            { kind: 'uups' }
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
            { kind: 'uups' }
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
            { kind: 'uups' }
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when tokenAllowlist is zero address', async () => {
        await expect(
          upgrades.deployProxy(
            processorEndpointFactory,
            [
              await teeAuthenticator.getAddress(),
              await authorityRegistry.getAddress(),
              updateStatusOperator,
              admin,
              ADDRESS_ZERO,
              minFeePerRequest,
              ADDRESS_ZERO,
            ],
            { kind: 'uups' }
          )
        ).to.be.revertedWithCustomError(processorEndpointFactory, 'AddressCantBeZero');
      });

      it('reverts when initialize is called a second time', async () => {
        const { processorEndpoint, tokenAllowlist } = await deployProcessorEndpoint();

        await expect(
          processorEndpoint.initialize(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            admin,
            resetOperator,
            minFeePerRequest,
            await tokenAllowlist.getAddress()
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidInitialization');
      });

      it('reverts when initialize is called directly on the implementation', async () => {
        const { processorEndpoint, tokenAllowlist } = await deployProcessorEndpoint();
        const implementationAddress = await upgrades.erc1967.getImplementationAddress(
          await processorEndpoint.getAddress()
        );
        const implementation = processorEndpointFactory.attach(implementationAddress);

        await expect(
          implementation.initialize(
            await teeAuthenticator.getAddress(),
            await authorityRegistry.getAddress(),
            updateStatusOperator,
            admin,
            resetOperator,
            minFeePerRequest,
            await tokenAllowlist.getAddress()
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidInitialization');
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

  describe('upgradeability', function () {
    it('reverts when a non-admin attempts to upgrade', async () => {
      const { processorEndpoint } = await deployProcessorEndpoint();
      const otherImplementation = await processorEndpointFactory.deploy();
      await otherImplementation.waitForDeployment();

      const signers = await ethers.getSigners();
      await expect(
        processorEndpoint
          .connect(signers[4])
          .upgradeToAndCall(await otherImplementation.getAddress(), '0x')
      ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
    });

    it('preserves storage across an upgrade to a new implementation', async () => {
      const { processorEndpoint } = await deployProcessorEndpoint();
      const proxyAddress = await processorEndpoint.getAddress();

      const teeAuthenticatorBefore = await processorEndpoint.teeAuthenticator();
      const minFeeBefore = await processorEndpoint.minFeePerRequest();
      const maxQueueSizeBefore = await processorEndpoint.maxQueueSize();

      // Upgrades must be sent by an ADMIN role holder — `admin` (signers[2]) per the fixture.
      const signers = await ethers.getSigners();
      const adminFactory = await ethers.getContractFactory('ProcessorEndpoint', signers[2]);
      const upgraded = await upgrades.upgradeProxy(proxyAddress, adminFactory);
      await upgraded.waitForDeployment();

      expect(await upgraded.getAddress()).to.equal(proxyAddress);
      expect(await upgraded.teeAuthenticator()).to.equal(teeAuthenticatorBefore);
      expect(await upgraded.minFeePerRequest()).to.equal(minFeeBefore);
      expect(await upgraded.maxQueueSize()).to.equal(maxQueueSizeBefore);
    });
  });
});
