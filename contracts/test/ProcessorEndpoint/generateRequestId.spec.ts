import { expect } from 'chai';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;

  beforeEach(async function () {
    const { deployProcessorEndpoint } = await deployProcessorEndpointFixture();
    processorEndpoint = await deployProcessorEndpoint();
  });

  describe('generateRequestId', function () {
    describe('unhappy paths', function () {
      it('returns different ids when any input changes', async () => {
        const sender = '0x0000000000000000000000000000000000000001';
        const applicationId = 1;
        const requestType = 1;
        const payload = '0x1234';
        const depositAmount = 10n;
        const idx = 0;

        const baseId = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          payload,
          depositAmount,
          idx
        );

        const diffSender = await processorEndpoint.generateRequestId(
          '0x0000000000000000000000000000000000000002',
          applicationId,
          requestType,
          payload,
          depositAmount,
          idx
        );
        const diffApplicationId = await processorEndpoint.generateRequestId(
          sender,
          applicationId + 1,
          requestType,
          payload,
          depositAmount,
          idx
        );
        const diffRequestType = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType + 1,
          payload,
          depositAmount,
          idx
        );
        const diffPayload = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          '0x1235',
          depositAmount,
          idx
        );
        const diffDeposit = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          payload,
          depositAmount + 1n,
          idx
        );
        const diffIdx = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          payload,
          depositAmount,
          idx + 1
        );

        expect(diffSender).to.not.equal(baseId);
        expect(diffApplicationId).to.not.equal(baseId);
        expect(diffRequestType).to.not.equal(baseId);
        expect(diffPayload).to.not.equal(baseId);
        expect(diffDeposit).to.not.equal(baseId);
        expect(diffIdx).to.not.equal(baseId);
      });
    });

    describe('happy paths', function () {
      it('returns deterministic id for the same inputs', async () => {
        const sender = '0x0000000000000000000000000000000000000001';
        const applicationId = 1;
        const requestType = 1;
        const payload = '0xabcd';
        const depositAmount = 42n;
        const idx = 7;

        const id1 = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          payload,
          depositAmount,
          idx
        );
        const id2 = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          payload,
          depositAmount,
          idx
        );

        expect(id1).to.equal(id2);
      });
    });
  });
});
