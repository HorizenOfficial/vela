import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;

  const PROTOCOL_VERSION = 0;
  const REQUEST_TYPE_PROCESS = 1;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  async function submitRequest(payload: string, depositAmount: bigint) {
    const maxFeeValue = minFeePerRequest;
    const tx = await processorEndpoint.submitRequest(
      PROTOCOL_VERSION,
      applicationId,
      REQUEST_TYPE_PROCESS,
      payload,
      depositAmount,
      maxFeeValue,
      { value: depositAmount + maxFeeValue }
    );
    const receipt = await tx.wait();
    return receipt.logs[0].args.requestId;
  }

  describe('getPendingRequestsPage', function () {
    describe('unhappy paths', function () {
      it('returns empty array when offset is out of range', async () => {
        await submitRequest('0x01', 0n);
        const requests = await processorEndpoint.getPendingRequestsPage(1, 2);
        expect(requests.length).to.equal(0);
      });

      it('returns empty array when limit is zero', async () => {
        await submitRequest('0x01', 0n);
        const requests = await processorEndpoint.getPendingRequestsPage(0, 0);
        expect(requests.length).to.equal(0);
      });
    });

    describe('happy paths', function () {
      it('returns paginated requests in order with correct data', async () => {
        const requestId1 = await submitRequest('0x01', 0n);
        const requestId2 = await submitRequest('0x02', 10n);
        const requestId3 = await submitRequest('0x03', 20n);

        const requests = await processorEndpoint.getPendingRequestsPage(1, 2);
        expect(requests.length).to.equal(2);
        expect(requests[0].requestId).to.equal(requestId2);
        expect(requests[0].payload).to.equal('0x02');
        expect(requests[0].depositAmount).to.equal(10n);
        expect(requests[1].requestId).to.equal(requestId3);
        expect(requests[1].payload).to.equal('0x03');
        expect(requests[1].depositAmount).to.equal(20n);

        const remaining = await processorEndpoint.getPendingRequestsPage(2, 10);
        expect(remaining.length).to.equal(1);
        expect(remaining[0].requestId).to.equal(requestId3);
        expect(remaining[0].sender).to.equal(await signers[0].getAddress());
        expect(remaining[0].requestType).to.equal(REQUEST_TYPE_PROCESS);
      });
    });
  });
});
