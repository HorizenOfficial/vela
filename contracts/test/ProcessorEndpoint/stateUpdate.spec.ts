import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture } from './fixture';
import { BYTES32_ZERO } from '../util';
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

  function getRequestIdFromReceipt(processorEndpointInstance: any, receipt: any) {
    for (const log of receipt.logs) {
      try {
        const parsed = processorEndpointInstance.interface.parseLog(log);
        if (parsed.name === 'RequestSubmitted') {
          return parsed.args.requestId;
        }
      } catch {
        continue;
      }
    }
    throw new Error('RequestSubmitted not found');
  }

  function getUserEvents(processorEndpointInstance: any, receipt: any) {
    return receipt.logs
      .map((log: any) => {
        try {
          return processorEndpointInstance.interface.parseLog(log);
        } catch {
          return null;
        }
      })
      .filter((parsed: any) => parsed && parsed.name === 'UserEvent');
  }

  async function submitRequest(
    processorEndpointInstance: any,
    sender: Signer,
    payload: string,
    depositAmount: bigint,
    maxFeeValue: bigint
  ) {
    const tx = await processorEndpointInstance.connect(sender).submitRequest(
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

  async function deployWithNoAttestation(teeSigner: Signer) {
    const fixture = await deployProcessorEndpointFixture();
    const NoAttestationTeeAuthenticator = await ethers.getContractFactory('NoAttestationTeeAuthenticator');
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

  describe('stateUpdate', function () {
    describe('unhappy paths', function () {
      it('reverts with InvalidApplicationId when applicationId is invalid', async () => {
        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            2,
            BYTES32_ZERO,
            '0x' + '11'.repeat(32),
            '0x' + '00'.repeat(32),
            [],
            [],
            [],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidApplicationId');
      });

      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {
        await expect(
          processorEndpoint.connect(signers[0]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '11'.repeat(32),
            '0x' + '00'.repeat(32),
            [],
            [],
            [],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidStateRoot when prevStateRoot does not match current stateRoot', async () => {
        const first = await submitRequest(processorEndpoint, signers[0], '0x01', 0n, minFeePerRequest);
        await processorEndpoint.connect(signers[1]).stateUpdate(
          APPLICATION_ID,
          BYTES32_ZERO,
          '0x' + '22'.repeat(32),
          first.requestId,
          [],
          [],
          [],
          0,
          minFeePerRequest,
          '0x'
        );

        const second = await submitRequest(processorEndpoint, signers[0], '0x02', 0n, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '33'.repeat(32),
            second.requestId,
            [],
            [],
            [],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidStateRoot');
      });

      it('reverts with InvalidRequestId when processedRequestId is not current pending', async () => {
        await submitRequest(processorEndpoint, signers[0], '0x03', 0n, minFeePerRequest);
        const second = await submitRequest(processorEndpoint, signers[0], '0x04', 0n, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '44'.repeat(32),
            second.requestId,
            [],
            [],
            [],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
      });

      it('reverts with InvalidPayload when events and eventSubTypes length mismatch', async () => {
        const request = await submitRequest(processorEndpoint, signers[0], '0x05', 0n, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '55'.repeat(32),
            request.requestId,
            ['0x01'],
            [],
            [],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
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
          fixture.minFeePerRequest
        );

        await expect(
          fixture.processorEndpoint.connect(fixture.signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '66'.repeat(32),
            request.requestId,
            [],
            [],
            [],
            0,
            fixture.minFeePerRequest,
            signature
          )
        ).to.be.revertedWithCustomError(fixture.processorEndpoint, 'InvalidSignature');
      });

      it('reverts with InvalidSignature when event subtype changes', async () => {
        const fixture = await deployWithNoAttestation(signers[3]);
        const request = await submitRequest(
          fixture.processorEndpoint,
          fixture.signers[0],
          '0x07',
          0n,
          fixture.minFeePerRequest
        );

        const events = ['0x01'];
        const signature = await ethSignStateUpdate(
          fixture.signers[3],
          APPLICATION_ID,
          BYTES32_ZERO,
          '0x' + '77'.repeat(32),
          request.requestId,
          events,
          ['typeA'],
          [],
          0,
          fixture.minFeePerRequest
        );

        await expect(
          fixture.processorEndpoint.connect(fixture.signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '77'.repeat(32),
            request.requestId,
            events,
            ['typeB'],
            [],
            0,
            fixture.minFeePerRequest,
            signature
          )
        ).to.be.revertedWithCustomError(fixture.processorEndpoint, 'InvalidSignature');
      });

      it('reverts with InvalidValue when refund + applicationFees != maxFeeValue', async () => {
        const request = await submitRequest(processorEndpoint, signers[0], '0x08', 0n, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '88'.repeat(32),
            request.requestId,
            [],
            [],
            [],
            1,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InvalidValue when applicationFees < minFeePerRequest', async () => {
        const maxFeeValue = minFeePerRequest + 1n;
        const request = await submitRequest(processorEndpoint, signers[0], '0x09', 0n, maxFeeValue);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + '99'.repeat(32),
            request.requestId,
            [],
            [],
            [],
            maxFeeValue - (minFeePerRequest - 1n),
            minFeePerRequest - 1n,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InsufficientBalance when withdrawals sum exceeds contract balance', async () => {
        const request = await submitRequest(processorEndpoint, signers[0], '0x0a', 0n, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + 'aa'.repeat(32),
            request.requestId,
            [],
            [],
            [[await signers[2].getAddress(), 1]],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InsufficientBalance');
      });

      it('reverts with TransferFailed when refund transfer fails', async () => {
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        await fallbackFailure.deploymentTransaction()!.wait();

        const maxFeeValue = minFeePerRequest + 1n;
        const insertTx = await fallbackFailure.insertRequestOnProcessorEndpoint(
          processorEndpoint,
          PROTOCOL_VERSION,
          APPLICATION_ID,
          REQUEST_TYPE,
          '0x0b',
          0,
          maxFeeValue,
          { value: maxFeeValue }
        );
        const insertReceipt = await insertTx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, insertReceipt);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + 'bb'.repeat(32),
            requestId,
            [],
            [],
            [],
            1,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'TransferFailed');
      });

      it('reverts with TransferFailed when fee transfer fails', async () => {
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        await fallbackFailure.deploymentTransaction()!.wait();

        await processorEndpoint
          .connect(signers[2])
          .updateFeeCollector(await fallbackFailure.getAddress());

        const request = await submitRequest(processorEndpoint, signers[0], '0x0c', 0n, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + 'cc'.repeat(32),
            request.requestId,
            [],
            [],
            [],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'TransferFailed');
      });

      it('reverts with TransferFailed when any withdrawal transfer fails', async () => {
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        await fallbackFailure.deploymentTransaction()!.wait();

        const request = await submitRequest(processorEndpoint, signers[0], '0x0d', 10n, minFeePerRequest);

        await expect(
          processorEndpoint.connect(signers[1]).stateUpdate(
            APPLICATION_ID,
            BYTES32_ZERO,
            '0x' + 'dd'.repeat(32),
            request.requestId,
            [],
            [],
            [[await fallbackFailure.getAddress(), 10]],
            0,
            minFeePerRequest,
            '0x'
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'TransferFailed');
      });
    });

    describe('happy paths', function () {
      it('updates state root and emits StateRootUpdate with valid signature', async () => {
        const fixture = await deployWithNoAttestation(signers[3]);
        const request = await submitRequest(
          fixture.processorEndpoint,
          fixture.signers[0],
          '0x0e',
          0n,
          fixture.minFeePerRequest
        );
        const newStateRoot = '0x' + 'ee'.repeat(32);

        const signature = await ethSignStateUpdate(
          fixture.signers[3],
          APPLICATION_ID,
          BYTES32_ZERO,
          newStateRoot,
          request.requestId,
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest
        );

        const tx = await fixture.processorEndpoint.connect(fixture.signers[1]).stateUpdate(
          APPLICATION_ID,
          BYTES32_ZERO,
          newStateRoot,
          request.requestId,
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest,
          signature
        );

        await expect(tx).to.emit(fixture.processorEndpoint, 'StateRootUpdate').withArgs(
          APPLICATION_ID,
          request.requestId,
          BYTES32_ZERO,
          newStateRoot
        );
        expect(await fixture.processorEndpoint.stateRoot()).to.equal(newStateRoot);
      });

      it('processes update: completes request, emits events, and transfers funds', async () => {
        const depositAmount = 20n;
        const refund = 5n;
        const applicationFees = minFeePerRequest;
        const maxFeeValue = refund + applicationFees;
        const withdrawalA = await signers[3].getAddress();
        const withdrawalB = await signers[4].getAddress();

        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x0f',
          depositAmount,
          maxFeeValue
        );
        const sender = await signers[0].getAddress();
        const senderBalanceAfterSubmit = await signers[0].provider!.getBalance(sender);
        const balanceA = await signers[3].provider!.getBalance(withdrawalA);
        const balanceB = await signers[4].provider!.getBalance(withdrawalB);

        const tx = await processorEndpoint.connect(signers[1]).stateUpdate(
          APPLICATION_ID,
          BYTES32_ZERO,
          '0x' + 'ff'.repeat(32),
          request.requestId,
          ['0xaa', '0xbb'],
          ['A', 'B'],
          [
            [withdrawalA, 10],
            [withdrawalB, 10],
          ],
          refund,
          applicationFees,
          '0x'
        );

        await expect(tx).to.emit(processorEndpoint, 'RequestCompleted').withArgs(
          request.requestId,
          applicationFees,
          0,
          0,
          ''
        );
        const receipt = await tx.wait();
        const userEvents = getUserEvents(processorEndpoint, receipt);
        const eventSubTypes = userEvents.map((event: any) => {
          const subType = event.args.eventSubType;
          return typeof subType === 'string' ? subType : subType.hash;
        });
        const eventPayloads = userEvents.map((event: any) => event.args.encryptedData);
        expect(userEvents.length).to.equal(2);
        expect(eventSubTypes).to.have.members([ethers.id('A'), ethers.id('B')]);
        expect(eventPayloads).to.have.members(['0xaa', '0xbb']);
        await expect(tx).to.emit(processorEndpoint, 'Refund').withArgs(
          APPLICATION_ID,
          request.requestId,
          sender,
          refund
        );
        await expect(tx).to.emit(processorEndpoint, 'Withdrawal').withArgs(
          APPLICATION_ID,
          request.requestId,
          withdrawalA,
          10
        );
        await expect(tx).to.emit(processorEndpoint, 'Withdrawal').withArgs(
          APPLICATION_ID,
          request.requestId,
          withdrawalB,
          10
        );

        const senderBalanceAfterUpdate = await signers[0].provider!.getBalance(sender);
        const balanceAAfter = await signers[3].provider!.getBalance(withdrawalA);
        const balanceBAfter = await signers[4].provider!.getBalance(withdrawalB);
        expect(senderBalanceAfterUpdate - senderBalanceAfterSubmit).to.equal(refund);
        expect(balanceAAfter - balanceA).to.equal(10n);
        expect(balanceBAfter - balanceB).to.equal(10n);
      });

      it('emits UserEvent for provided events and subtypes', async () => {
        const request = await submitRequest(processorEndpoint, signers[0], '0x10', 0n, minFeePerRequest);

        const tx = await processorEndpoint.connect(signers[1]).stateUpdate(
          APPLICATION_ID,
          BYTES32_ZERO,
          '0x' + '10'.repeat(32),
          request.requestId,
          ['0x11', '0x22'],
          ['type1', 'type2'],
          [],
          0,
          minFeePerRequest,
          '0x'
        );

        const receipt = await tx.wait();
        const userEvents = getUserEvents(processorEndpoint, receipt);
        const eventSubTypes = userEvents.map((event: any) => {
          const subType = event.args.eventSubType;
          return typeof subType === 'string' ? subType : subType.hash;
        });
        const eventPayloads = userEvents.map((event: any) => event.args.encryptedData);
        expect(userEvents.length).to.equal(2);
        expect(eventSubTypes).to.have.members([ethers.id('type1'), ethers.id('type2')]);
        expect(eventPayloads).to.have.members(['0x11', '0x22']);
      });

      it('allows first update when stateRoot is zero and prevStateRoot is zero', async () => {
        const request = await submitRequest(processorEndpoint, signers[0], '0x12', 0n, minFeePerRequest);
        const newStateRoot = '0x' + '12'.repeat(32);

        await processorEndpoint.connect(signers[1]).stateUpdate(
          APPLICATION_ID,
          BYTES32_ZERO,
          newStateRoot,
          request.requestId,
          [],
          [],
          [],
          0,
          minFeePerRequest,
          '0x'
        );

        expect(await processorEndpoint.stateRoot()).to.equal(newStateRoot);
      });

      it('allows update when stateRoot is zero and prevStateRoot is non-zero', async () => {
        const request = await submitRequest(processorEndpoint, signers[0], '0x13', 0n, minFeePerRequest);
        const prevStateRoot = '0x' + '13'.repeat(32);
        const newStateRoot = '0x' + '14'.repeat(32);

        await processorEndpoint.connect(signers[1]).stateUpdate(
          APPLICATION_ID,
          prevStateRoot,
          newStateRoot,
          request.requestId,
          [],
          [],
          [],
          0,
          minFeePerRequest,
          '0x'
        );

        expect(await processorEndpoint.stateRoot()).to.equal(newStateRoot);
      });

      it('emits Refund even when refund amount is zero', async () => {
        const request = await submitRequest(processorEndpoint, signers[0], '0x15', 0n, minFeePerRequest);
        const sender = await signers[0].getAddress();

        const tx = await processorEndpoint.connect(signers[1]).stateUpdate(
          APPLICATION_ID,
          BYTES32_ZERO,
          '0x' + '15'.repeat(32),
          request.requestId,
          [],
          [],
          [],
          0,
          minFeePerRequest,
          '0x'
        );

        await expect(tx).to.emit(processorEndpoint, 'Refund').withArgs(
          APPLICATION_ID,
          request.requestId,
          sender,
          0
        );
      });
    });
  });
});
