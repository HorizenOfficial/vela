import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
  });

  describe('getPendingRequests', function () {
    describe('unhappy paths', function () {
      it('returns empty array when there are no pending requests', async () => {
        const requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(0);
      });
    });

    describe('happy paths', function () {
      it('returns pending requests in FIFO order with correct data', async () => {
        const protocolVersion = 0;
        const applicationId = 1;
        const requestType = 1;
        const depositAmount1 = 0n;
        const depositAmount2 = 10n;
        const maxFeeValue = minFeePerRequest;

        const tx1 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x01',
          depositAmount1,
          maxFeeValue,
          { value: depositAmount1 + maxFeeValue }
        );
        const receipt1 = await tx1.wait();
        const requestId1 = receipt1.logs[0].args.requestId;

        const tx2 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          depositAmount2,
          maxFeeValue,
          { value: depositAmount2 + maxFeeValue }
        );
        const receipt2 = await tx2.wait();
        const requestId2 = receipt2.logs[0].args.requestId;

        const requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(2);

        expect(requests[0].requestId).to.equal(requestId1);
        expect(requests[0].protocolVersion).to.equal(protocolVersion);
        expect(requests[0].applicationId).to.equal(applicationId);
        expect(requests[0].requestType).to.equal(requestType);
        expect(requests[0].payload).to.equal('0x01');
        expect(requests[0].depositAmount).to.equal(depositAmount1);
        expect(requests[0].maxFeeValue).to.equal(maxFeeValue);
        expect(requests[0].sender).to.equal(await signers[0].getAddress());

        expect(requests[1].requestId).to.equal(requestId2);
        expect(requests[1].protocolVersion).to.equal(protocolVersion);
        expect(requests[1].applicationId).to.equal(applicationId);
        expect(requests[1].requestType).to.equal(requestType);
        expect(requests[1].payload).to.equal('0x02');
        expect(requests[1].depositAmount).to.equal(depositAmount2);
        expect(requests[1].maxFeeValue).to.equal(maxFeeValue);
        expect(requests[1].sender).to.equal(await signers[0].getAddress());
      });

      it('removes head request after completion or failure', async () => {
        const protocolVersion = 0;
        const applicationId = 1;
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
        const requestId1 = receipt1.logs[0].args.requestId;

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
        const requestId2 = receipt2.logs[0].args.requestId;

        await processorEndpoint
          .connect(signers[1])
          .markRequestCompleted(requestId1, 0, maxFeeValue);

        let requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(1);
        expect(requests[0].requestId).to.equal(requestId2);

        await processorEndpoint
          .connect(signers[1])
          .markRequestFailed(requestId2, 1, 'failed');

        requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(0);
      });

      it('accepts new requests after the queue is emptied', async () => {
        const protocolVersion = 0;
        const applicationId = 1;
        const requestType = 1;
        const depositAmount = 0n;
        const maxFeeValue = minFeePerRequest;

        const tx1 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x01',
          depositAmount,
          maxFeeValue,
          { value: depositAmount + maxFeeValue }
        );
        const receipt1 = await tx1.wait();
        const requestId1 = receipt1.logs[0].args.requestId;

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
        const requestId2 = receipt2.logs[0].args.requestId;

        await processorEndpoint
          .connect(signers[1])
          .markRequestCompleted(requestId1, 0, maxFeeValue);
        await processorEndpoint
          .connect(signers[1])
          .markRequestFailed(requestId2, 1, 'failed');

        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);

        const tx3 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x03',
          depositAmount,
          maxFeeValue,
          { value: depositAmount + maxFeeValue }
        );
        const receipt3 = await tx3.wait();
        const requestId3 = receipt3.logs[0].args.requestId;

        const [request, , success] = await processorEndpoint.getNextPendingRequest();
        expect(success).to.equal(true);
        expect(request.requestId).to.equal(requestId3);
      });
    });
  });
});
