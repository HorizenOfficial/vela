import { expect } from 'chai';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
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

  describe('getNextPendingRequest', function () {
    describe('unhappy paths', function () {
      it('returns success=false and empty request when queue is empty', async () => {
        const [request, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).to.equal(false);
        expect(stateRoot).to.equal(BYTES32_ZERO);
        expect(request.requestId).to.equal(BYTES32_ZERO);
      });
    });

    describe('happy paths', function () {
      it('returns current pending request, stateRoot, and success=true', async () => {
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

        const [request, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).to.equal(true);
        expect(stateRoot).to.equal(INITIAL_STATE_ROOT);
        expect(request.requestId).to.equal(requestId);
        expect(request.protocolVersion).to.equal(protocolVersion);
        expect(request.applicationId).to.equal(applicationId);
        expect(request.requestType).to.equal(requestType);
        expect(request.payload).to.equal(payload);
        expect(request.depositAmount).to.equal(depositAmount);
      });
    });
  });
});
