import { expect } from 'chai';
import { deployProcessorEndpointFixture } from './fixture';

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

  describe('getPendingRequestsSize', function () {
    describe('unhappy paths', function () {
      it('returns 0 when queue is empty', async () => {
        const size = await processorEndpoint.getPendingRequestsSize();
        expect(size).to.equal(0n);
      });
    });

    describe('happy paths', function () {
      it('returns tail - head when queue has pending requests', async () => {
        const protocolVersion = 0;
        const requestType = 1;
        const payload = '0x01';
        const depositAmount = 0n;
        const maxFeeValue = minFeePerRequest;

        await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          payload,
          depositAmount,
          maxFeeValue,
          { value: depositAmount + maxFeeValue }
        );

        let size = await processorEndpoint.getPendingRequestsSize();
        expect(size).to.equal(1n);

        await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          depositAmount,
          maxFeeValue,
          { value: depositAmount + maxFeeValue }
        );

        size = await processorEndpoint.getPendingRequestsSize();
        expect(size).to.equal(2n);
      });
    });
  });
});
