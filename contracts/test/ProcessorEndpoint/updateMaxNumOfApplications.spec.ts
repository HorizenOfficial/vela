import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let bootstrapApplication: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
  });

  describe('updateMaxNumOfApplications', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks ADMIN role', async () => {
        await expect(
          processorEndpoint.connect(signers[0]).updateMaxNumOfApplications(5)
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidValue when newMax is zero', async () => {
        await expect(
          processorEndpoint.connect(signers[2]).updateMaxNumOfApplications(0)
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InvalidValue when newMax is below deployed and reserved slots', async () => {
        // Deploy 2 apps (each consumes a slot)
        await bootstrapApplication(processorEndpoint);
        await bootstrapApplication(processorEndpoint);

        // Also submit a pending deploy request (reserves a 3rd slot)
        await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(0, '0x01', { value: minFeePerRequest });

        // 3 slots used (2 deployed + 1 pending), setting max to 2 should fail
        await expect(
          processorEndpoint.connect(signers[2]).updateMaxNumOfApplications(2)
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });
    });

    describe('happy paths', function () {
      it('updates availableDeploySlots when no deploy requests have been sent', async () => {
        const oldMax = await processorEndpoint.maxNumOfApplications();
        const oldAvailableSlots = await processorEndpoint.availableDeploySlots();
        const newMax = 5n;
        const tx = await processorEndpoint.connect(signers[2]).updateMaxNumOfApplications(newMax);
        await expect(tx)
          .to.emit(processorEndpoint, 'MaxNumberOfAppUpdated')
          .withArgs(oldMax, newMax);
        const usedSlots = oldMax - oldAvailableSlots;
        expect(await processorEndpoint.availableDeploySlots()).to.equal(newMax - usedSlots);
      });

      it('updates maxNumOfApplications and availableDeploySlots accounting for used slots', async () => {
        // Deploy 2 apps (each consumes a slot)
        await bootstrapApplication(processorEndpoint);
        await bootstrapApplication(processorEndpoint);

        // Also submit a pending deploy request (reserves a 3rd slot)
        await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(0, '0x01', { value: minFeePerRequest });

        // 3 slots used (2 deployed + 1 pending), setting max to 4 should leave 1 available
        await processorEndpoint.connect(signers[2]).updateMaxNumOfApplications(4);
        expect(await processorEndpoint.maxNumOfApplications()).to.equal(4n);
        expect(await processorEndpoint.availableDeploySlots()).to.equal(1n);
      });
    });
  });
});
