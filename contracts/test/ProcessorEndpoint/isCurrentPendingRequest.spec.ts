import { expect } from 'chai';
import { deployProcessorEndpointFixture } from './fixture';
import { BYTES32_ZERO } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let minFeePerRequest: bigint;
  let applicationId: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  describe('isCurrentPendingRequest', function () {
    describe('unhappy paths', function () {
      it('returns false when queue is empty', async () => {
        const isCurrent = await processorEndpoint.isCurrentPendingRequest(BYTES32_ZERO);
        expect(isCurrent).to.equal(false);
      });

      it('returns false for a request that is not at the head', async () => {
        const protocolVersion = 0;
        const requestType = 1;
        const payload = '0x01';
        const depositAmount = 0n;
        const maxFeeValue = minFeePerRequest;

        const tx1 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          payload,
          depositAmount,
          maxFeeValue,
          { value: depositAmount + maxFeeValue }
        );
        const receipt1 = await tx1.wait();
        const firstRequestId = receipt1.logs[0].args.requestId;

        const tx2 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          depositAmount,
          maxFeeValue,
          { value: depositAmount + maxFeeValue }
        );
        const receipt2 = await tx2.wait();
        const secondRequestId = receipt2.logs[0].args.requestId;

        const isCurrent = await processorEndpoint.isCurrentPendingRequest(secondRequestId);
        expect(isCurrent).to.equal(false);
        expect(await processorEndpoint.isCurrentPendingRequest(firstRequestId)).to.equal(true);
      });
    });

    describe('happy paths', function () {
      it('returns true for the request at the head of the queue', async () => {
        const protocolVersion = 0;
        const requestType = 1;
        const payload = '0x01';
        const depositAmount = 0n;
        const maxFeeValue = minFeePerRequest;

        const tx = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          payload,
          depositAmount,
          maxFeeValue,
          { value: depositAmount + maxFeeValue }
        );
        const receipt = await tx.wait();
        const requestId = receipt.logs[0].args.requestId;

        const isCurrent = await processorEndpoint.isCurrentPendingRequest(requestId);
        expect(isCurrent).to.equal(true);
      });
    });
  });
});
