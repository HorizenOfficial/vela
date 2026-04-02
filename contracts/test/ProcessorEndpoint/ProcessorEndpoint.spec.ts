import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  ETH_TOKEN,
  BYTES32_ZERO,
  getRequestIdFromReceipt as getRequestIdFromReceiptBase,
} from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;

  const PROTOCOL_VERSION = 0;
  const REQUEST_TYPE = 1;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  function getRequestIdFromReceipt(receipt: any) {
    return getRequestIdFromReceiptBase(processorEndpoint, receipt);
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
        applicationId,
        REQUEST_TYPE,
        payload,
        ETH_TOKEN,
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
      expect(pending[1].assetAmount).to.equal(5n);

      const [nextPending, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
      expect(success).to.equal(true);
      expect(stateRoot).to.equal(INITIAL_STATE_ROOT);
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

      // With pull pattern, funds are credited to pending deposits
      const senderABalanceAfterSubmit = await senderA.provider!.getBalance(
        await senderA.getAddress()
      );
      const senderAPendingAmountBefore = await processorEndpoint.pendingClaims(ETH_TOKEN,
        await senderA.getAddress()
      );

      const refundA = 1n;
      const applicationFeesA = first.maxFeeValue - refundA;
      await expect(
        processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            '0x' + '01'.repeat(32),
            first.requestId,
            [],
            [],
            [],
            refundA,
            applicationFeesA,
            0,
            '',
            '0x'
          )
      ).to.emit(processorEndpoint, 'RequestCompleted');

      const senderABalanceAfterComplete = await senderA.provider!.getBalance(
        await senderA.getAddress()
      );
      const senderAPendingAmountAfter = await processorEndpoint.pendingClaims(ETH_TOKEN,
        await senderA.getAddress()
      );
      expect(senderABalanceAfterComplete - senderABalanceAfterSubmit).to.equal(0n);
      expect(senderAPendingAmountAfter - senderAPendingAmountBefore).to.equal(refundA);

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);
      expect(await processorEndpoint.isCurrentPendingRequest(second.requestId)).to.equal(true);

      const senderBBalanceAfterSubmit = await senderB.provider!.getBalance(
        await senderB.getAddress()
      );
      const senderBPendingAfterSubmit = await processorEndpoint.pendingClaims(ETH_TOKEN,
        await senderB.getAddress()
      );
      // Fail second request via stateUpdate with errorCode
      const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
      const failTx = await processorEndpoint
        .connect(signers[1])
        .stateUpdate(
          applicationId,
          currentStateRoot,
          currentStateRoot,
          second.requestId,
          [],
          [],
          [],
          0,
          0,
          1,
          'failed',
          '0x'
        );
      const expectedRefundB = second.depositAmount + (second.maxFeeValue - minFeePerRequest);
      await expect(failTx)
        .to.emit(processorEndpoint, 'Refund')
        .withArgs(applicationId, second.requestId, ETH_TOKEN, await senderB.getAddress(), expectedRefundB);
      await expect(failTx).to.emit(processorEndpoint, 'RequestCompleted');

      const senderBBalanceAfterFail = await senderB.provider!.getBalance(
        await senderB.getAddress()
      );
      const senderBPendingAmountAfter = await processorEndpoint.pendingClaims(ETH_TOKEN,
        await senderB.getAddress()
      );
      expect(senderBBalanceAfterFail - senderBBalanceAfterSubmit).to.equal(0n);
      expect(senderBPendingAmountAfter - senderBPendingAfterSubmit).to.equal(expectedRefundB);

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
      const [, , success] = await processorEndpoint.getNextPendingRequest();
      expect(success).to.equal(false);
    });
  });
});
