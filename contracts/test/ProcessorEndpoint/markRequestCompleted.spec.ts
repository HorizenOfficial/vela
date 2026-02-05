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

  describe('markRequestCompleted', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {
        await expect(
          processorEndpoint.markRequestCompleted('0x' + '00'.repeat(32), 0, minFeePerRequest)
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidRequestId when request is not current pending', async () => {
        await submitBasicRequest(signers[0], '0x01');
        const { requestId } = await submitBasicRequest(signers[0], '0x02');

        await expect(
          processorEndpoint.connect(signers[1]).markRequestCompleted(requestId, 0, minFeePerRequest)
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
      });

      it('reverts with InvalidRequestId when request was already completed', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x03');

        await processorEndpoint
          .connect(signers[1])
          .markRequestCompleted(requestId, 0, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).markRequestCompleted(requestId, 0, minFeePerRequest)
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
      });

      it('reverts with InvalidValue when refund + applicationFees != maxFeeValue', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x04');

        await expect(
          processorEndpoint.connect(signers[1]).markRequestCompleted(requestId, 1, minFeePerRequest)
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InvalidValue when applicationFees < minFeePerRequest', async () => {
        const maxFeeValue = minFeePerRequest + 1n;
        const { requestId } = await submitBasicRequest(signers[0], '0x05', maxFeeValue);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .markRequestCompleted(requestId, 2, minFeePerRequest - 1n)
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });
    });

    describe('happy paths', function () {
      it('does not revert when fee transfer fails', async () => {
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        await fallbackFailure.deploymentTransaction()!.wait();

        await processorEndpoint
          .connect(signers[2])
          .updateFeeCollector(await fallbackFailure.getAddress());

        const { requestId } = await submitBasicRequest(signers[0], '0x06');

        const tx = await processorEndpoint
          .connect(signers[1])
          .markRequestCompleted(requestId, 0, minFeePerRequest);

        await expect(tx).to.emit(processorEndpoint, 'RequestCompleted');
      });

      it('does not revert when refund transfer fails', async () => {
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        await fallbackFailure.deploymentTransaction()!.wait();

        const maxFeeValue = minFeePerRequest + 1n;
        const insertTx = await fallbackFailure.insertRequestOnProcessorEndpoint(
          processorEndpoint,
          PROTOCOL_VERSION,
          APPLICATION_ID,
          REQUEST_TYPE,
          '0x07',
          0,
          maxFeeValue,
          { value: maxFeeValue }
        );
        const insertReceipt = await insertTx.wait();
        const requestId = getRequestIdFromReceipt(insertReceipt);

        const tx = await processorEndpoint
          .connect(signers[1])
          .markRequestCompleted(requestId, 1, minFeePerRequest);

        await expect(tx).to.emit(processorEndpoint, 'RequestCompleted');
      });

      it('completes request without refund and transfers fees', async () => {
        const newCollector = await signers[3].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const { requestId } = await submitBasicRequest(signers[0], '0x08');
        // With pull pattern, funds are credited to pending deposits
        const collectorPendingAmountBefore = await processorEndpoint.payments(newCollector);

        await processorEndpoint
          .connect(signers[1])
          .markRequestCompleted(requestId, 0, minFeePerRequest);

        const collectorPendingAmountAfter = await processorEndpoint.payments(newCollector);
        expect(collectorPendingAmountAfter - collectorPendingAmountBefore).to.equal(
          minFeePerRequest
        );
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
      });

      it('completes request with refund and emits Refund', async () => {
        const newCollector = await signers[3].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const sender = signers[4];
        const maxFeeValue = minFeePerRequest + 5n;
        const refund = 3n;
        const applicationFees = maxFeeValue - refund;

        const tx = await processorEndpoint
          .connect(sender)
          .submitRequest(PROTOCOL_VERSION, APPLICATION_ID, REQUEST_TYPE, '0x09', 0, maxFeeValue, {
            value: maxFeeValue,
          });
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(receipt);
        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.payments(
          await sender.getAddress()
        );

        const completeTx = await processorEndpoint
          .connect(signers[1])
          .markRequestCompleted(requestId, refund, applicationFees);

        await expect(completeTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(APPLICATION_ID, requestId, await sender.getAddress(), refund);

        const senderBalanceAfterComplete = await sender.provider!.getBalance(
          await sender.getAddress()
        );
        const senderPendingAmountAfterComplete = await processorEndpoint.payments(
          await sender.getAddress()
        );
        expect(senderPendingAmountAfterComplete - senderPendingAmountAfterSubmit).to.equal(refund);
      });
    });
  });
});
