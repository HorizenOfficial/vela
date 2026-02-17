import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture } from './fixture';
import { BYTES32_ZERO, getRequestIdFromReceipt } from '../util';
import { ethSignStateUpdate } from '../../scripts/util';

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

  async function submitBasicRequest(sender: Signer, payload: string, maxFeeValue?: bigint) {
    const fee = maxFeeValue ?? minFeePerRequest;
    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(PROTOCOL_VERSION, APPLICATION_ID, REQUEST_TYPE, payload, 0, fee, {
        value: fee,
      });
    const receipt = await tx.wait();
    return { requestId: getRequestIdFromReceipt(processorEndpoint, receipt), maxFeeValue: fee };
  }

  async function failRequest(requestId: string, errorCode: number, errorMsg: string) {
    const currentStateRoot = await processorEndpoint.stateRoot();
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        APPLICATION_ID,
        currentStateRoot,
        currentStateRoot,
        requestId,
        [],
        [],
        [],
        0,
        0,
        errorCode,
        errorMsg,
        '0x'
      );
  }

  async function deployWithNoAttestation(teeSigner: Signer) {
    const fixture = await deployProcessorEndpointFixture();
    const NoAttestationTeeAuthenticator = await ethers.getContractFactory(
      'NoAttestationTeeAuthenticator'
    );
    const pk = '0x' + '11'.repeat(133);
    const teeAuthenticator = await NoAttestationTeeAuthenticator.deploy(
      await fixture.signers[0].getAddress(),
      await teeSigner.getAddress(),
      pk
    );
    const processorEndpointWithNoAttestation = await fixture.processorEndpointFactory.deploy(
      await teeAuthenticator.getAddress(),
      await fixture.authorityRegistry.getAddress(),
      fixture.updateStatusOperator,
      fixture.admin,
      fixture.minFeePerRequest
    );

    return {
      ...fixture,
      teeAuthenticator,
      processorEndpoint: processorEndpointWithNoAttestation,
    };
  }

  async function submitRequest(
    processorEndpointInstance: any,
    sender: Signer,
    payload: string,
    depositAmount: bigint,
    maxFeeValue: bigint
  ) {
    const tx = await processorEndpointInstance
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
    const requestId = getRequestIdFromReceipt(processorEndpointInstance, receipt);
    return { requestId, maxFeeValue, depositAmount };
  }

  describe('stateUpdate error path', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {
        await expect(
          processorEndpoint
            .connect(signers[0])
            .stateUpdate(
              APPLICATION_ID,
              BYTES32_ZERO,
              BYTES32_ZERO,
              '0x' + '00'.repeat(32),
              [],
              [],
              [],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidRequestId when request is not current pending', async () => {
        await submitBasicRequest(signers[0], '0x01');
        const { requestId } = await submitBasicRequest(signers[0], '0x02');

        await expect(failRequest(requestId, 1, 'err')).to.be.revertedWithCustomError(
          processorEndpoint,
          'InvalidRequestId'
        );
      });

      it('reverts with InvalidRequestId when request was already failed', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x03');

        await failRequest(requestId, 1, 'err');

        await expect(failRequest(requestId, 1, 'err')).to.be.revertedWithCustomError(
          processorEndpoint,
          'InvalidRequestId'
        );
      });

      it('reverts with InvalidStateRoot when prevStateRoot does not match current stateRoot', async () => {
        const first = await submitBasicRequest(signers[0], '0x01', minFeePerRequest);
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '22'.repeat(32),
            first.requestId,
            [],
            [],
            [],
            0,
            minFeePerRequest,
            0,
            '',
            '0x'
          );

        const second = await submitBasicRequest(signers[0], '0x02');

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              APPLICATION_ID,
              BYTES32_ZERO,
              '0x' + '33'.repeat(32),
              second.requestId,
              [],
              [],
              [],
              0,
              minFeePerRequest,
              1,
              'failed',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidStateRoot');
      });

      it('reverts with InvalidStateRoot when prevStateRoot does not match with newStateRoot', async () => {
        const failedTx = await submitBasicRequest(signers[0], '0x01', minFeePerRequest);
        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              APPLICATION_ID,
              BYTES32_ZERO,
              '0x' + '22'.repeat(32),
              failedTx.requestId,
              [],
              [],
              [],
              0,
              minFeePerRequest,
              1,
              'failed',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidStateRoot');
      });

      it('reverts with InvalidSignature when teeAuthenticator checkSignature fails', async () => {
        const fixture = await deployWithNoAttestation(signers[3]);
        const request = await submitRequest(
          fixture.processorEndpoint,
          fixture.signers[0],
          '0x06',
          0n,
          fixture.minFeePerRequest
        );

        const signature = await ethSignStateUpdate(
          fixture.signers[4],
          APPLICATION_ID,
          BYTES32_ZERO,
          '0x' + '66'.repeat(32),
          request.requestId,
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest,
          1,
          'failed'
        );

        await expect(
          fixture.processorEndpoint
            .connect(fixture.signers[1])
            .stateUpdate(
              APPLICATION_ID,
              BYTES32_ZERO,
              '0x' + '66'.repeat(32),
              request.requestId,
              [],
              [],
              [],
              0,
              fixture.minFeePerRequest,
              1,
              '',
              signature
            )
        ).to.be.revertedWithCustomError(fixture.processorEndpoint, 'InvalidSignature');
      });

      it('reverts with InvalidRequestId when there are no pending requests', async () => {
        await expect(failRequest('0x' + '00'.repeat(32), 1, 'err')).to.be.revertedWithCustomError(
          processorEndpoint,
          'InvalidRequestId'
        );
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
        const requestId = getRequestIdFromReceipt(processorEndpoint, insertReceipt);

        const tx = await failRequest(requestId, 1, 'err');
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
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.payments(
          await sender.getAddress()
        );

        const expectedRefund = depositAmount + (maxFeeValue - minFeePerRequest);

        const failTx = await failRequest(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(APPLICATION_ID, requestId, await sender.getAddress(), expectedRefund);

        const senderPendingAmountAfterComplete = await processorEndpoint.payments(
          await sender.getAddress()
        );
        expect(senderPendingAmountAfterComplete - senderPendingAmountAfterSubmit).to.equal(
          expectedRefund
        );
      });

      it('refunds sender and emits RequestCompleted FAILED when error stateUpdate succeeds', async () => {
        const sender = signers[4];
        const maxFeeValue = minFeePerRequest + 6n;

        const tx = await processorEndpoint
          .connect(sender)
          .submitRequest(PROTOCOL_VERSION, APPLICATION_ID, REQUEST_TYPE, '0x05', 0, maxFeeValue, {
            value: maxFeeValue,
          });
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.payments(
          await sender.getAddress()
        );

        const failTx = await failRequest(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'RequestCompleted')
          .withArgs(requestId, minFeePerRequest, 1, 1, 'err');

        const senderPendingAmountAfterFail = await processorEndpoint.payments(
          await sender.getAddress()
        );
        const expectedRefund = maxFeeValue - minFeePerRequest;
        expect(senderPendingAmountAfterFail - senderPendingAmountAfterSubmit).to.equal(
          expectedRefund
        );
      });
    });
  });
});
