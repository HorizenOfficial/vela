import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';
import { BYTES32_ZERO } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;

  const PROTOCOL_VERSION = 0;
  const APPLICATION_ID = 1;
  const REQUEST_TYPE = 1;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
  });

  function getRequestIdFromReceipt(receipt: any) {
    for (const log of receipt.logs) {
      try {
        const parsed = processorEndpoint.interface.parseLog(log);
        if (parsed.name === 'RequestSubmitted') {
          return parsed.args.requestId;
        }
      } catch {
        continue;
      }
    }
    throw new Error('RequestSubmitted not found');
  }

  async function submitRequest(
    sender: Signer,
    payload: string,
    depositAmount: bigint,
    maxFeeValue: bigint
  ) {
    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        APPLICATION_ID,
        REQUEST_TYPE,
        payload,
        depositAmount,
        maxFeeValue,
        { value: depositAmount + maxFeeValue }
      );
    const receipt = await tx.wait();
    return { requestId: getRequestIdFromReceipt(receipt), depositAmount, maxFeeValue };
  }

  // Most tests are split into per-function specs; integration coverage lives here.
  describe('integration', function () {
    it('saves multiple requests and retrieves queue/request data', async () => {
      const first = await submitRequest(signers[0], '0x01', 0n, minFeePerRequest);
      const second = await submitRequest(signers[0], '0x02', 5n, minFeePerRequest);

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(2n);

      const pending = await processorEndpoint.getPendingRequests();
      expect(pending.length).to.equal(2);
      expect(pending[0].requestId).to.equal(first.requestId);
      expect(pending[1].requestId).to.equal(second.requestId);
      expect(pending[1].depositAmount).to.equal(5n);

      const [nextPending, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
      expect(success).to.equal(true);
      expect(stateRoot).to.equal(BYTES32_ZERO);
      expect(nextPending.requestId).to.equal(first.requestId);

      expect(await processorEndpoint.isCurrentPendingRequest(first.requestId)).to.equal(true);
      expect(await processorEndpoint.isCurrentPendingRequest(second.requestId)).to.equal(false);

      const storedFirst = await processorEndpoint.requestById(first.requestId);
      const storedSecond = await processorEndpoint.requestById(second.requestId);
      expect(storedFirst.sender).to.equal(await signers[0].getAddress());
      expect(storedSecond.payload).to.equal('0x02');
    });

    it('completes request via stateUpdate then fails another and handles refunds', async () => {
      const senderA = signers[0];
      const senderB = signers[3];

      const first = await submitRequest(senderA, '0x03', 0n, minFeePerRequest + 2n);
      const second = await submitRequest(senderB, '0x04', 0n, minFeePerRequest + 4n);

      const senderABalanceAfterSubmit = await senderA.provider!.getBalance(
        await senderA.getAddress()
      );

      const refundA = 1n;
      const applicationFeesA = first.maxFeeValue - refundA;
      await expect(
        processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '01'.repeat(32),
            first.requestId,
            [],
            [],
            [],
            refundA,
            applicationFeesA,
            '0x'
          )
      ).to.emit(processorEndpoint, 'RequestCompleted');

      const senderABalanceAfterComplete = await senderA.provider!.getBalance(
        await senderA.getAddress()
      );
      expect(senderABalanceAfterComplete - senderABalanceAfterSubmit).to.equal(refundA);

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);
      expect(await processorEndpoint.isCurrentPendingRequest(second.requestId)).to.equal(true);

      const senderBBalanceAfterSubmit = await senderB.provider!.getBalance(
        await senderB.getAddress()
      );
      const failTx = await processorEndpoint
        .connect(signers[1])
        .markRequestFailed(second.requestId, 1, 'failed');
      await expect(failTx).to.emit(processorEndpoint, 'RequestCompleted');

      const senderBBalanceAfterFail = await senderB.provider!.getBalance(
        await senderB.getAddress()
      );
      const expectedRefundB = second.maxFeeValue - minFeePerRequest;
      expect(senderBBalanceAfterFail - senderBBalanceAfterSubmit).to.equal(expectedRefundB);

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
      const [, , success] = await processorEndpoint.getNextPendingRequest();
      expect(success).to.equal(false);
    });
  });
});
