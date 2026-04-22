import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';
import { ETH_TOKEN } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  describe('updateQueueThreshold', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks ADMIN role', async () => {
        await expect(
          processorEndpoint.connect(signers[1]).updateQueueThreshold(5)
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidValue when newThreshold is zero', async () => {
        await expect(
          processorEndpoint.connect(signers[2]).updateQueueThreshold(0)
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });
    });

    describe('happy paths', function () {
      it('updates maxQueueSize and emits QueueThresholdUpdated', async () => {
        const tx = await processorEndpoint.connect(signers[2]).updateQueueThreshold(5);
        await expect(tx).to.emit(processorEndpoint, 'QueueThresholdUpdated').withArgs(5);
        expect(await processorEndpoint.maxQueueSize()).to.equal(5n);
      });

      it('prevents new requests when newThreshold is below current queue size', async () => {
        const protocolVersion = 0;
        const requestType = 1;
        const payload = '0x01';
        const assetAmount = 0n;
        const maxFeeValue = minFeePerRequest;

        await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          payload,
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
        );
        await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
        );

        await processorEndpoint.connect(signers[2]).updateQueueThreshold(1);

        await expect(
          processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            requestType,
            '0x03',
            ETH_TOKEN,
            assetAmount,
            maxFeeValue,
            { value: assetAmount + maxFeeValue }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'QueueThresholdExceeded');
      });
    });
  });
});
