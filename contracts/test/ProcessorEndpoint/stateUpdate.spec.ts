import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  ETH_TOKEN,
  BYTES32_ZERO,
  getRequestIdFromReceipt,
  REQUEST_TYPE_PROCESS,
  PROTOCOL_VERSION,
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
    ({ processorEndpoint } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  function getEventsByName(processorEndpointInstance: any, receipt: any, eventName: string) {
    return receipt.logs
      .map((log: any) => {
        try {
          return processorEndpointInstance.interface.parseLog(log);
        } catch {
          return null;
        }
      })
      .filter((parsed: any) => parsed && parsed.name === eventName);
  }

  function getUserEvents(processorEndpointInstance: any, receipt: any) {
    return getEventsByName(processorEndpointInstance, receipt, 'UserEvent');
  }

  function getAppEvents(processorEndpointInstance: any, receipt: any) {
    return getEventsByName(processorEndpointInstance, receipt, 'AppEvent');
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
      fixture.resetOperator,
      fixture.minFeePerRequest,
      await fixture.sharedTokenAllowlist.getAddress(),
      fixture.extensionAddress
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

  describe('stateUpdate', function () {
    describe('unhappy paths', function () {
      it('reverts with InvalidApplicationId when applicationId is invalid', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x00',
          0n,
          minFeePerRequest
        );

        const invalidAppId = applicationId + 1n; // assuming this appId does not exist
        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              invalidAppId,
              INITIAL_STATE_ROOT,
              '0x' + '11'.repeat(32),
              request.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidApplicationId');
      });

      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {
        await expect(
          processorEndpoint
            .connect(signers[0])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '11'.repeat(32),
              '0x' + '00'.repeat(32),
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidStateRoot when prevStateRoot does not match current stateRoot', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x13',
          0n,
          minFeePerRequest
        );
        const prevStateRoot = '0x' + '13'.repeat(32);
        const newStateRoot = '0x' + '14'.repeat(32);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              prevStateRoot,
              newStateRoot,
              request.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidStateRoot');

        expect(await processorEndpoint.applicationStateRoots(applicationId)).to.equal(
          INITIAL_STATE_ROOT
        );
      });

      it('reverts with InvalidStateRoot when prevStateRoot does not match after a prior update', async () => {
        const first = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x01',
          0n,
          minFeePerRequest
        );
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            '0x' + '22'.repeat(32),
            first.requestId,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            minFeePerRequest,
            0,
            '',
            '0x'
          );

        const second = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x02',
          0n,
          minFeePerRequest
        );

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              BYTES32_ZERO,
              '0x' + '33'.repeat(32),
              second.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidStateRoot');
      });

      it('reverts with InvalidRequestId when processedRequestId is not current pending', async () => {
        await submitRequest(processorEndpoint, signers[0], '0x03', 0n, minFeePerRequest);
        const second = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x04',
          0n,
          minFeePerRequest
        );

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '44'.repeat(32),
              second.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
      });

      it('reverts with InvalidPayload when userEventData events and subTypes length mismatch', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x05',
          0n,
          minFeePerRequest
        );

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '55'.repeat(32),
              request.requestId,
              { events: ['0x01'], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidPayload when appEventData events and subTypes length mismatch', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x0500',
          0n,
          minFeePerRequest
        );

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '55'.repeat(32),
              request.requestId,
              { events: [], subTypes: [] },
              { events: ['0x01'], subTypes: [] },
              [],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidSignature when teeAuthenticator checkBatchSignature fails', async () => {
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
          '0x' + '66'.repeat(32),
          request.requestId,
          [],
          [],
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest
        );

        await expect(
          fixture.processorEndpoint
            .connect(fixture.signers[1])
            .stateUpdate(
              fixture.applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '66'.repeat(32),
              request.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              0,
              fixture.minFeePerRequest,
              0,
              '',
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
          fixture.minFeePerRequest,
          fixture.applicationId
        );

        const events = ['0x01'];
        const signature = await ethSignStateUpdate(
          fixture.signers[3],
          fixture.applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '77'.repeat(32),
          request.requestId,
          events,
          [ethers.encodeBytes32String('typeA')],
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest
        );

        await expect(
          fixture.processorEndpoint
            .connect(fixture.signers[1])
            .stateUpdate(
              fixture.applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '77'.repeat(32),
              request.requestId,
              { events: events, subTypes: [ethers.encodeBytes32String('typeB')] },
              { events: [], subTypes: [] },
              [],
              0,
              fixture.minFeePerRequest,
              0,
              '',
              signature
            )
        ).to.be.revertedWithCustomError(fixture.processorEndpoint, 'InvalidSignature');
      });

      it('reverts with InvalidValue when refund + applicationFees != maxFeeValue', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x08',
          0n,
          minFeePerRequest
        );

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '88'.repeat(32),
              request.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              1,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InvalidValue when applicationFees < minFeePerRequest', async () => {
        const maxFeeValue = minFeePerRequest + 1n;
        const applicationFees = minFeePerRequest - 1n;
        const request = await submitRequest(processorEndpoint, signers[0], '0x09', 0n, maxFeeValue);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '99'.repeat(32),
              request.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [],
              maxFeeValue - applicationFees,
              applicationFees,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InsufficientAppBalance when withdrawals sum exceeds app locked funds', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x0a',
          0n,
          minFeePerRequest
        );

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + 'aa'.repeat(32),
              request.requestId,
              { events: [], subTypes: [] },
              { events: [], subTypes: [] },
              [[ETH_TOKEN, await signers[2].getAddress(), 1]],
              0,
              minFeePerRequest,
              0,
              '',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InsufficientAppBalance');
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
          fixture.minFeePerRequest,
          fixture.applicationId
        );
        const newStateRoot = '0x' + 'ee'.repeat(32);

        const signature = await ethSignStateUpdate(
          fixture.signers[3],
          fixture.applicationId,
          INITIAL_STATE_ROOT,
          newStateRoot,
          request.requestId,
          [],
          [],
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest
        );

        const tx = await fixture.processorEndpoint
          .connect(fixture.signers[1])
          .stateUpdate(
            fixture.applicationId,
            INITIAL_STATE_ROOT,
            newStateRoot,
            request.requestId,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            fixture.minFeePerRequest,
            0,
            '',
            signature
          );

        await expect(tx)
          .to.emit(fixture.processorEndpoint, 'StateRootUpdate')
          .withArgs(fixture.applicationId, request.requestId, INITIAL_STATE_ROOT, newStateRoot);
        expect(
          await fixture.processorEndpoint.applicationStateRoots(fixture.applicationId)
        ).to.equal(newStateRoot);
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
        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          sender
        );
        const balanceAPendingAmountAfterSubmit = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          withdrawalA
        );
        const balanceBPendingAmountAfterSubmit = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          withdrawalB
        );

        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(
          depositAmount
        );

        const withdrawalAAmount = 15n;
        const withdrawalBAmount = 5n;

        const tx = await processorEndpoint.connect(signers[1]).stateUpdate(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + 'ff'.repeat(32),
          request.requestId,
          {
            events: ['0xaa', '0xbb'],
            subTypes: [ethers.encodeBytes32String('A'), ethers.encodeBytes32String('B')],
          },
          { events: [], subTypes: [] },
          [
            [ETH_TOKEN, withdrawalA, withdrawalAAmount],
            [ETH_TOKEN, withdrawalB, withdrawalBAmount],
          ],
          refund,
          applicationFees,
          0,
          '',
          '0x'
        );

        await expect(tx)
          .to.emit(processorEndpoint, 'RequestCompleted')
          .withArgs(applicationId, request.requestId, applicationFees, 0, 0, '');
        const receipt = await tx.wait();
        const userEvents = getUserEvents(processorEndpoint, receipt);
        const eventSubTypes = userEvents.map((event: any) => {
          const subType = event.args.eventSubType;
          return typeof subType === 'string' ? subType : subType.hash;
        });
        const eventPayloads = userEvents.map((event: any) => event.args.encryptedData);
        expect(userEvents.length).to.equal(2);
        expect(eventSubTypes).to.have.members([
          ethers.encodeBytes32String('A'),
          ethers.encodeBytes32String('B'),
        ]);
        expect(eventPayloads).to.have.members(['0xaa', '0xbb']);
        await expect(tx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(applicationId, request.requestId, sender, ETH_TOKEN, refund);
        await expect(tx)
          .to.emit(processorEndpoint, 'Withdrawal')
          .withArgs(applicationId, request.requestId, withdrawalA, ETH_TOKEN, withdrawalAAmount);
        await expect(tx)
          .to.emit(processorEndpoint, 'Withdrawal')
          .withArgs(applicationId, request.requestId, withdrawalB, ETH_TOKEN, withdrawalBAmount);

        const senderPendingAmountAfterUpdate = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          sender
        );
        const balanceAPendingAmountAfterUpdate = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          withdrawalA
        );
        const balanceBPendingAmountAfterUpdate = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          withdrawalB
        );
        expect(senderPendingAmountAfterUpdate - senderPendingAmountAfterSubmit).to.equal(refund);
        expect(balanceAPendingAmountAfterUpdate - balanceAPendingAmountAfterSubmit).to.equal(
          withdrawalAAmount
        );
        expect(balanceBPendingAmountAfterUpdate - balanceBPendingAmountAfterSubmit).to.equal(
          withdrawalBAmount
        );

        // appLockedFunds: credited depositAmount(20), debited withdrawals(15+5) = 20, so should be 0
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
      });

      it('emits UserEvent for provided events and subtypes', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x10',
          0n,
          minFeePerRequest
        );

        const tx = await processorEndpoint.connect(signers[1]).stateUpdate(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '10'.repeat(32),
          request.requestId,
          {
            events: ['0x11', '0x22'],
            subTypes: [ethers.encodeBytes32String('type1'), ethers.encodeBytes32String('type2')],
          },
          { events: [], subTypes: [] },
          [],
          0,
          minFeePerRequest,
          0,
          '',
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
        expect(eventSubTypes).to.have.members([
          ethers.encodeBytes32String('type1'),
          ethers.encodeBytes32String('type2'),
        ]);
        expect(eventPayloads).to.have.members(['0x11', '0x22']);
      });

      it('updates stateRoot from initial value with matching prevStateRoot', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x12',
          0n,
          minFeePerRequest
        );
        const newStateRoot = '0x' + '12'.repeat(32);

        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            newStateRoot,
            request.requestId,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            minFeePerRequest,
            0,
            '',
            '0x'
          );

        expect(await processorEndpoint.applicationStateRoots(applicationId)).to.equal(newStateRoot);
      });

      it("doesn't emit Refund event when refund amount is zero", async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x15',
          0n,
          minFeePerRequest
        );
        const sender = await signers[0].getAddress();

        const tx = await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            '0x' + '15'.repeat(32),
            request.requestId,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            minFeePerRequest,
            0,
            '',
            '0x'
          );

        await expect(tx).not.to.emit(processorEndpoint, 'Refund');
      });

      it('emits AppEvent for provided appEventData', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x16',
          0n,
          minFeePerRequest
        );

        const tx = await processorEndpoint.connect(signers[1]).stateUpdate(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '16'.repeat(32),
          request.requestId,
          { events: [], subTypes: [] },
          {
            events: ['0x33', '0x44'],
            subTypes: [
              ethers.encodeBytes32String('appType1'),
              ethers.encodeBytes32String('appType2'),
            ],
          },
          [],
          0,
          minFeePerRequest,
          0,
          '',
          '0x'
        );

        const receipt = await tx.wait();
        const appEvents = getAppEvents(processorEndpoint, receipt);
        const eventSubTypes = appEvents.map((event: any) => {
          const subType = event.args.eventSubType;
          return typeof subType === 'string' ? subType : subType.hash;
        });
        const eventPayloads = appEvents.map((event: any) => event.args.data);
        expect(appEvents.length).to.equal(2);
        expect(eventSubTypes).to.have.members([
          ethers.encodeBytes32String('appType1'),
          ethers.encodeBytes32String('appType2'),
        ]);
        expect(eventPayloads).to.have.members(['0x33', '0x44']);
      });

      it('reverts with InvalidPayload when appEventData is non-empty on error path', async () => {
        const request = await submitRequest(
          processorEndpoint,
          signers[0],
          '0x17',
          0n,
          minFeePerRequest
        );

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              INITIAL_STATE_ROOT,
              request.requestId,
              { events: [], subTypes: [] },
              { events: ['0xdeadbeef'], subTypes: [ethers.encodeBytes32String('subtype')] },
              [],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });
    });
  });
});
