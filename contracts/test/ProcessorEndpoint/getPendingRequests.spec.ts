import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
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
        const requestType = 1;
        const assetAmount1 = 0n;
        const assetAmount2 = 10n;
        const maxFeeValue = minFeePerRequest;

        const tx1 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x01',
          ETH_TOKEN,
          assetAmount1,
          maxFeeValue,
          { value: assetAmount1 + maxFeeValue }
        );
        const receipt1 = await tx1.wait();
        const requestId1 = receipt1.logs[0].args.requestId;

        const tx2 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          ETH_TOKEN,
          assetAmount2,
          maxFeeValue,
          { value: assetAmount2 + maxFeeValue }
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
        expect(requests[0].tokenAddress).to.equal(ETH_TOKEN);
        expect(requests[0].assetAmount).to.equal(assetAmount1);
        expect(requests[0].maxFeeValue).to.equal(maxFeeValue);
        expect(requests[0].sender).to.equal(await signers[0].getAddress());

        expect(requests[1].requestId).to.equal(requestId2);
        expect(requests[1].protocolVersion).to.equal(protocolVersion);
        expect(requests[1].applicationId).to.equal(applicationId);
        expect(requests[1].requestType).to.equal(requestType);
        expect(requests[1].payload).to.equal('0x02');
        expect(requests[1].tokenAddress).to.equal(ETH_TOKEN);
        expect(requests[1].assetAmount).to.equal(assetAmount2);
        expect(requests[1].maxFeeValue).to.equal(maxFeeValue);
        expect(requests[1].sender).to.equal(await signers[0].getAddress());
      });

      it('removes head request after completion or failure', async () => {
        const protocolVersion = 0;
        const requestType = 1;
        const payload = '0x01';
        const assetAmount = 0n;
        const maxFeeValue = minFeePerRequest;

        const tx1 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          payload,
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
        );
        const receipt1 = await tx1.wait();
        const requestId1 = receipt1.logs[0].args.requestId;

        const tx2 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
        );
        const receipt2 = await tx2.wait();
        const requestId2 = receipt2.logs[0].args.requestId;

        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            '0x' + '01'.repeat(32),
            requestId1,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            maxFeeValue,
            0,
            '',
            '0x'
          );

        let requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(1);
        expect(requests[0].requestId).to.equal(requestId2);

        // Fail second request via stateUpdate with errorCode
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            currentStateRoot,
            currentStateRoot,
            requestId2,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            0,
            1,
            'failed',
            '0x'
          );

        requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(0);
      });

      it('accepts new requests after the queue is emptied', async () => {
        const protocolVersion = 0;
        const requestType = 1;
        const assetAmount = 0n;
        const maxFeeValue = minFeePerRequest;

        const tx1 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x01',
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
        );
        const receipt1 = await tx1.wait();
        const requestId1 = receipt1.logs[0].args.requestId;

        const tx2 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
        );
        const receipt2 = await tx2.wait();
        const requestId2 = receipt2.logs[0].args.requestId;

        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            '0x' + '01'.repeat(32),
            requestId1,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            maxFeeValue,
            0,
            '',
            '0x'
          );

        // Fail second request via stateUpdate with errorCode
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            currentStateRoot,
            currentStateRoot,
            requestId2,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            0,
            1,
            'failed',
            '0x'
          );

        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);

        const tx3 = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x03',
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
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
