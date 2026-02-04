import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture } from './fixture';

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

  async function submitBasicRequest(sender: Signer, payload: string, maxFeeValue?: bigint) {
    const fee = maxFeeValue ?? minFeePerRequest;
    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(PROTOCOL_VERSION, APPLICATION_ID, REQUEST_TYPE, payload, 0, fee, {
        value: fee,
      });
    const receipt = await tx.wait();
    return { requestId: getRequestIdFromReceipt(receipt), maxFeeValue: fee };
  }

  describe('markRequestFailed', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {
        await expect(
          processorEndpoint.markRequestFailed('0x' + '00'.repeat(32), 1, 'err')
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidRequestId when request is not current pending', async () => {
        await submitBasicRequest(signers[0], '0x01');
        const { requestId } = await submitBasicRequest(signers[0], '0x02');

        await expect(
          processorEndpoint.connect(signers[1]).markRequestFailed(requestId, 1, 'err')
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
      });

      it('reverts with InvalidRequestId when request was already failed', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x03');

        await processorEndpoint.connect(signers[1]).markRequestFailed(requestId, 1, 'err');

        await expect(
          processorEndpoint.connect(signers[1]).markRequestFailed(requestId, 1, 'err')
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
      });

      it('reverts with InvalidRequestId when there are no pending requests', async () => {
        await expect(
          processorEndpoint.connect(signers[1]).markRequestFailed('0x' + '00'.repeat(32), 1, 'err')
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
      });
    });

    describe('happy paths', function () {
      it('does not revert when refund transfer fails', async () => {
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        await fallbackFailure.deploymentTransaction()!.wait();

        const maxFeeValue = minFeePerRequest + 4n;
        const insertTx = await fallbackFailure.insertRequestOnProcessorEndpoint(
          processorEndpoint,
          PROTOCOL_VERSION,
          APPLICATION_ID,
          REQUEST_TYPE,
          '0x04',
          0,
          maxFeeValue,
          { value: maxFeeValue }
        );
        const insertReceipt = await insertTx.wait();
        const requestId = getRequestIdFromReceipt(insertReceipt);

        const tx = await processorEndpoint
          .connect(signers[1])
          .markRequestFailed(requestId, 1, 'err');
        await expect(tx).to.emit(processorEndpoint, 'RequestCompleted');
      });

      it('refunds deposit plus remaining fee for failed requests with deposits', async () => {
        const sender = signers[4];
        const depositAmount = 25n;
        const maxFeeValue = minFeePerRequest + 5n;

        const tx = await processorEndpoint
          .connect(sender)
          .submitRequest(
            PROTOCOL_VERSION,
            APPLICATION_ID,
            REQUEST_TYPE,
            '0x07',
            depositAmount,
            maxFeeValue,
            { value: depositAmount + maxFeeValue }
          );
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(receipt);

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.payments(await sender.getAddress());

        const expectedRefund = depositAmount + (maxFeeValue - minFeePerRequest);

        const failTx = await processorEndpoint
          .connect(signers[1])
          .markRequestFailed(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(APPLICATION_ID, requestId, await sender.getAddress(), expectedRefund);

        const senderPendingAmountAfterComplete = await processorEndpoint.payments(await sender.getAddress());
        expect(senderPendingAmountAfterComplete - senderPendingAmountAfterSubmit).to.equal(expectedRefund);
      });

      it('refunds sender and emits RequestCompleted FAILED_REFUNDED when fee transfer succeeds', async () => {
        const sender = signers[4];
        const maxFeeValue = minFeePerRequest + 6n;

        const tx = await processorEndpoint
          .connect(sender)
          .submitRequest(PROTOCOL_VERSION, APPLICATION_ID, REQUEST_TYPE, '0x05', 0, maxFeeValue, {
            value: maxFeeValue,
          });
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(receipt);

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.payments(await sender.getAddress());

        const failTx = await processorEndpoint
          .connect(signers[1])
          .markRequestFailed(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'RequestCompleted')
          .withArgs(requestId, minFeePerRequest, 1, 1, 'err');

        const senderPendingAmountAfterFail = await processorEndpoint.payments(await sender.getAddress());
        const expectedRefund = maxFeeValue - minFeePerRequest;
        expect(senderPendingAmountAfterFail - senderPendingAmountAfterSubmit).to.equal(expectedRefund);
      });

    });
  });
});
