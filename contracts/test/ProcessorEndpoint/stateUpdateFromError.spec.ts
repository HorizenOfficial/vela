import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  ETH_TOKEN,
  BYTES32_ZERO,
  getRequestIdFromReceipt,
  PROTOCOL_VERSION,
  REQUEST_TYPE_PROCESS,
} from '../util';
import { ethSignStateUpdate } from '../../scripts/util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;
  let bootstrapApplication: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  async function submitBasicRequest(sender: Signer, payload: string, maxFeeValue?: bigint) {
    const fee = maxFeeValue ?? minFeePerRequest;
    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        applicationId,
        REQUEST_TYPE_PROCESS,
        payload,
        ETH_TOKEN,
        0,
        fee,
        {
          value: fee,
        }
      );
    const receipt = await tx.wait();
    return { requestId: getRequestIdFromReceipt(processorEndpoint, receipt), maxFeeValue: fee };
  }

  async function failRequest(requestId: string, errorCode: number, errorMsg: string) {
    const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        applicationId,
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

    const { applicationId: appId } = await fixture.bootstrapApplication(
      processorEndpointWithNoAttestation,
      teeSigner
    );

    return {
      ...fixture,
      teeAuthenticator,
      processorEndpoint: processorEndpointWithNoAttestation,
      applicationId: appId,
    };
  }

  async function submitRequest(
    processorEndpointInstance: any,
    sender: Signer,
    payload: string,
    depositAmount: bigint,
    maxFeeValue: bigint,
    appId?: bigint
  ) {
    const tx = await processorEndpointInstance
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        appId ?? applicationId,
        REQUEST_TYPE_PROCESS,
        payload,
        ETH_TOKEN,
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
              applicationId,
              INITIAL_STATE_ROOT,
              INITIAL_STATE_ROOT,
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
            applicationId,
            INITIAL_STATE_ROOT,
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
              applicationId,
              BYTES32_ZERO,
              BYTES32_ZERO,
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
              applicationId,
              INITIAL_STATE_ROOT,
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
          fixture.minFeePerRequest,
          fixture.applicationId
        );

        const signature = await ethSignStateUpdate(
          fixture.signers[4],
          fixture.applicationId,
          INITIAL_STATE_ROOT,
          INITIAL_STATE_ROOT,
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
              fixture.applicationId,
              INITIAL_STATE_ROOT,
              INITIAL_STATE_ROOT,
              request.requestId,
              [],
              [],
              [],
              0,
              fixture.minFeePerRequest,
              1,
              'error',
              signature
            )
        ).to.be.revertedWithCustomError(fixture.processorEndpoint, 'InvalidSignature');
      });

      it('reverts with InvalidPayload when events array is non-empty on error', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x08');
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              currentStateRoot,
              currentStateRoot,
              requestId,
              ['0xdeadbeef'],
              ['subtype'],
              [],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidPayload when withdrawalRequests array is non-empty on error', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x09');
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              currentStateRoot,
              currentStateRoot,
              requestId,
              [],
              [],
              [{ tokenAddress: ETH_TOKEN, receiver: await signers[2].getAddress(), amount: 1 }],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidPayload when both events and withdrawalRequests are non-empty on error', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x0a');
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              currentStateRoot,
              currentStateRoot,
              requestId,
              ['0xdeadbeef'],
              ['subtype'],
              [{ tokenAddress: ETH_TOKEN, receiver: await signers[2].getAddress(), amount: 1 }],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
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
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x04',
          ETH_TOKEN,
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
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x07',
            ETH_TOKEN,
            depositAmount,
            maxFeeValue,
            { value: depositAmount + maxFeeValue }
          );
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        // After submit: appLockedFunds should equal depositAmount only (fees tracked globally)
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(
          depositAmount
        );

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );

        const expectedRefund = depositAmount + (maxFeeValue - minFeePerRequest);

        const failTx = await failRequest(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(applicationId, requestId, ETH_TOKEN, await sender.getAddress(), expectedRefund);

        const senderPendingAmountAfterComplete = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );
        expect(senderPendingAmountAfterComplete - senderPendingAmountAfterSubmit).to.equal(
          expectedRefund
        );

        // After error: appLockedFunds debited by depositAmount → 0
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
      });

      it('refunds sender and emits RequestCompleted FAILED when error stateUpdate succeeds', async () => {
        const sender = signers[4];
        const maxFeeValue = minFeePerRequest + 6n;

        const tx = await processorEndpoint
          .connect(sender)
          .submitRequest(
            PROTOCOL_VERSION,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x05',
            ETH_TOKEN,
            0,
            maxFeeValue,
            { value: maxFeeValue }
          );
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        // After submit: appLockedFunds should be 0 (depositAmount is 0, fees tracked globally)
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );

        const failTx = await failRequest(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'RequestCompleted')
          .withArgs(applicationId, requestId, minFeePerRequest, 1, 1, 'err');

        const senderPendingAmountAfterFail = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );
        const expectedRefund = maxFeeValue - minFeePerRequest;
        expect(senderPendingAmountAfterFail - senderPendingAmountAfterSubmit).to.equal(
          expectedRefund
        );

        // After error: appLockedFunds remains 0 (no deposit was tracked)
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
      });

      it('restores an available deploy slot when a deploy request fails', async () => {
        const slotsBefore = await processorEndpoint.availableDeploySlots();

        const deployTx = await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(PROTOCOL_VERSION, '0x01', { value: minFeePerRequest });
        const deployReceipt = await deployTx.wait();
        const deployLog = deployReceipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(deployLog);
        const deployAppId: bigint = parsed.args.applicationId;
        const deployRequestId: string = parsed.args.requestId;

        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore - 1n);

        const failTx = await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            deployAppId,
            BYTES32_ZERO,
            BYTES32_ZERO,
            deployRequestId,
            [],
            [],
            [],
            0,
            0,
            1,
            'deploy failed',
            '0x'
          );

        await expect(failTx).to.emit(processorEndpoint, 'DeployRequestCompleted');
        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore);
      });
    });
  });
});
