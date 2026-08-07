import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { ETH_TOKEN, BYTES32_ZERO, getRequestIdFromReceipt, REQUEST_TYPE_PROCESS } from '../util';

// The round-robin turn is enforced on-chain: a state update for a per-application request is
// rejected when the cursor selects a different application. `selectionGrace` is what keeps that
// rule race-free — the selection view leaves no on-chain trace, so a request enqueued between the
// manager reading it and the update landing legitimately changes the scan result, and a selected
// head younger than the grace period is therefore ignored.
//
// Two properties are pinned here: the race case (a fresh arrival must not invalidate an in-flight
// update) and the skip bound (once the selected head ages past the grace, the turn is mandatory,
// so the window in which an application can be skipped is `selectionGrace`, not the depth of
// another application's backlog).
describe('ProcessorEndpoint selection enforcement', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let fixture: any;
  // firstAppId is at index 0 of _deployedAppIds, so the cursor (which starts at 0 and is not
  // advanced by deploys) points at it: it has the first turn.
  let firstAppId: bigint;
  let secondAppId: bigint;
  let grace: bigint;

  beforeEach(async function () {
    fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId: firstAppId } = await fixture.bootstrapApplication(processorEndpoint));
    ({ applicationId: secondAppId } = await fixture.bootstrapApplication(processorEndpoint));
    grace = await processorEndpoint.selectionGrace();
  });

  async function submitRequest(appId: bigint, payload: string): Promise<string> {
    const tx = await processorEndpoint
      .connect(signers[0])
      .submitRequest(0, appId, REQUEST_TYPE_PROCESS, payload, ETH_TOKEN, 0n, minFeePerRequest, {
        value: minFeePerRequest,
      });
    return getRequestIdFromReceipt(processorEndpoint, await tx.wait());
  }

  function stateUpdate(
    appId: bigint,
    prevStateRoot: string,
    newStateRoot: string,
    requestId: string,
    applicationFees: bigint
  ) {
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        appId,
        prevStateRoot,
        newStateRoot,
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        applicationFees,
        0,
        '',
        '0x'
      );
  }

  async function advanceTime(seconds: bigint) {
    await ethers.provider.send('evm_increaseTime', ['0x' + seconds.toString(16)]);
    await ethers.provider.send('evm_mine', []);
  }

  it('accepts the application whose turn it is, however old its head', async function () {
    const requestId = await submitRequest(firstAppId, '0x01');
    await advanceTime(grace * 10n);

    await expect(
      stateUpdate(
        firstAppId,
        INITIAL_STATE_ROOT,
        '0x' + '11'.repeat(32),
        requestId,
        minFeePerRequest
      )
    ).to.not.be.reverted;
  });

  // The race the grace period exists for: firstApp's queue was empty when the manager read the
  // selection view (so the view returned secondApp), and a request for firstApp arrived while the
  // update was in flight. The in-flight update must still land.
  it('accepts an out-of-turn update when the selected head arrived within the grace period', async function () {
    // Selection view at this point returns secondAppId — firstApp's queue is empty.
    const secondRequestId = await submitRequest(secondAppId, '0x02');
    const [selectedBefore] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
    expect(selectedBefore).to.equal(secondAppId);

    // A request for firstApp lands during the window, flipping the selection.
    await submitRequest(firstAppId, '0x03');
    const [selectedAfter] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
    expect(selectedAfter).to.equal(firstAppId);

    await expect(
      stateUpdate(
        secondAppId,
        INITIAL_STATE_ROOT,
        '0x' + '22'.repeat(32),
        secondRequestId,
        minFeePerRequest
      )
    ).to.not.be.reverted;
  });

  // The bound: the same skip is rejected once the selected head is older than the grace period,
  // because it can no longer be explained by the selection race.
  it('rejects an out-of-turn update once the selected head is older than the grace period', async function () {
    const secondRequestId = await submitRequest(secondAppId, '0x04');
    await submitRequest(firstAppId, '0x05');

    await advanceTime(grace + 1n);

    await expect(
      stateUpdate(
        secondAppId,
        INITIAL_STATE_ROOT,
        '0x' + '33'.repeat(32),
        secondRequestId,
        minFeePerRequest
      )
    )
      .to.be.revertedWithCustomError(processorEndpoint, 'ApplicationNotSelected')
      .withArgs(firstAppId);
  });

  it('accepts the selected application after an out-of-turn update was rejected', async function () {
    await submitRequest(secondAppId, '0x06');
    const firstRequestId = await submitRequest(firstAppId, '0x07');
    await advanceTime(grace + 1n);

    await expect(
      stateUpdate(
        firstAppId,
        INITIAL_STATE_ROOT,
        '0x' + '44'.repeat(32),
        firstRequestId,
        minFeePerRequest
      )
    ).to.not.be.reverted;
  });

  // The skip cannot be repeated to build up unbounded starvation: serving out of turn advances the
  // cursor past the served application, and the skipped head keeps ageing, so the very next
  // attempt is rejected.
  it('does not allow the skip to be repeated once the skipped head has aged', async function () {
    const secondRequestA = await submitRequest(secondAppId, '0x08');
    const secondRequestB = await submitRequest(secondAppId, '0x09');
    await submitRequest(firstAppId, '0x0a');

    // First skip is inside the grace window and is allowed.
    const rootA = '0x' + '55'.repeat(32);
    await (
      await stateUpdate(secondAppId, INITIAL_STATE_ROOT, rootA, secondRequestA, minFeePerRequest)
    ).wait();

    await advanceTime(grace + 1n);

    await expect(
      stateUpdate(secondAppId, rootA, '0x' + '66'.repeat(32), secondRequestB, minFeePerRequest)
    )
      .to.be.revertedWithCustomError(processorEndpoint, 'ApplicationNotSelected')
      .withArgs(firstAppId);
  });

  // A young head opens the race exemption only for applications up to the first old head in scan
  // order — it must not license leapfrogging a starving application. Otherwise a colluding
  // submitter could reopen the window each rotation (one fresh request into an empty queue ahead
  // of the victim), and serving an application *past* the victim would jump the cursor over it,
  // deferring its turn by a rotation per attacker-controlled queue instead of the documented one.
  it('does not allow a fresh arrival to license skipping past a starving application', async function () {
    const { applicationId: thirdAppId } = await fixture.bootstrapApplication(processorEndpoint);

    const secondRequestId = await submitRequest(secondAppId, '0x13');
    const thirdRequestId = await submitRequest(thirdAppId, '0x14');
    await advanceTime(grace + 1n);
    // Fresh arrival at the cursor position (firstApp was empty): the exemption is open.
    const firstRequestId = await submitRequest(firstAppId, '0x15');

    // Serving thirdApp would leapfrog secondApp, whose old head the manager must have seen.
    await expect(
      stateUpdate(
        thirdAppId,
        INITIAL_STATE_ROOT,
        '0x' + 'cc'.repeat(32),
        thirdRequestId,
        minFeePerRequest
      )
    )
      .to.be.revertedWithCustomError(processorEndpoint, 'ApplicationNotSelected')
      .withArgs(secondAppId);

    // The young head itself sits before the first old head, so serving it is the race case and
    // is allowed; the cursor then advances toward the starving application, never past it.
    await (
      await stateUpdate(
        firstAppId,
        INITIAL_STATE_ROOT,
        '0x' + 'dd'.repeat(32),
        firstRequestId,
        minFeePerRequest
      )
    ).wait();

    await expect(
      stateUpdate(
        secondAppId,
        INITIAL_STATE_ROOT,
        '0x' + 'ee'.repeat(32),
        secondRequestId,
        minFeePerRequest
      )
    ).to.not.be.reverted;
  });

  it('enforces the turn with no tolerance when the grace period is zero', async function () {
    await (await processorEndpoint.connect(signers[2]).updateSelectionGrace(0)).wait();

    const secondRequestId = await submitRequest(secondAppId, '0x0b');
    await submitRequest(firstAppId, '0x0c');

    await expect(
      stateUpdate(
        secondAppId,
        INITIAL_STATE_ROOT,
        '0x' + '77'.repeat(32),
        secondRequestId,
        minFeePerRequest
      )
    )
      .to.be.revertedWithCustomError(processorEndpoint, 'ApplicationNotSelected')
      .withArgs(firstAppId);
  });

  // The escape hatch for a permanently failing queue head (BATCH_EXECUTION.md section 7.4): a
  // grace period longer than any request age disables enforcement.
  it('disables enforcement when the grace period exceeds any request age', async function () {
    await (
      await processorEndpoint.connect(signers[2]).updateSelectionGrace(365n * 24n * 3600n)
    ).wait();

    const secondRequestId = await submitRequest(secondAppId, '0x0d');
    await submitRequest(firstAppId, '0x0e');
    await advanceTime(grace * 100n);

    await expect(
      stateUpdate(
        secondAppId,
        INITIAL_STATE_ROOT,
        '0x' + '88'.repeat(32),
        secondRequestId,
        minFeePerRequest
      )
    ).to.not.be.reverted;
  });

  describe('exemptions', function () {
    it('processes a deploy request while another application is starved', async function () {
      await submitRequest(firstAppId, '0x0f');
      await advanceTime(grace + 1n);

      const deployTx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequest(0, '0x10', { value: minFeePerRequest });
      const deployReceipt = await deployTx.wait();
      const parsed = deployReceipt.logs
        .map((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log);
          } catch {
            return null;
          }
        })
        .find((p: any) => p?.name === 'DeployRequestSubmitted');

      await expect(
        stateUpdate(
          parsed.args.applicationId,
          BYTES32_ZERO,
          '0x' + '99'.repeat(32),
          parsed.args.requestId,
          minFeePerRequest
        )
      ).to.not.be.reverted;
    });

    it('processes a TRUSTPROCESS while another application is starved', async function () {
      const TestTrigger = await ethers.getContractFactory('TestTrigger');
      const trigger = await TestTrigger.deploy(await processorEndpoint.getAddress(), false, false);
      const deployTx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequestWithTrigger(0, '0x' + 'ab'.repeat(40), await trigger.getAddress(), {
          value: minFeePerRequest,
        });
      const parsed = (await deployTx.wait()).logs
        .map((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log);
          } catch {
            return null;
          }
        })
        .find((p: any) => p?.name === 'DeployRequestSubmitted');
      const triggerAppId: bigint = parsed.args.applicationId;
      await (
        await stateUpdate(
          triggerAppId,
          BYTES32_ZERO,
          INITIAL_STATE_ROOT,
          parsed.args.requestId,
          minFeePerRequest
        )
      ).wait();

      // Enqueue a TRUSTPROCESS while every other queue is still empty, so this update is in turn.
      await (await trigger.setTrustedPayload('0xdeadbeef')).wait();
      const seedRequestId = await submitRequest(triggerAppId, '0x11');
      const rootA = '0x' + 'aa'.repeat(32);
      await (
        await stateUpdate(triggerAppId, INITIAL_STATE_ROOT, rootA, seedRequestId, minFeePerRequest)
      ).wait();
      await (await trigger.setTrustedPayload('0x')).wait();
      const [, triggerHead] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
      const trustedRequestId: string = triggerHead[0].requestId;

      // Now starve firstApp, and process the trusted request anyway.
      await submitRequest(firstAppId, '0x12');
      await advanceTime(grace + 1n);

      await expect(stateUpdate(triggerAppId, rootA, '0x' + 'bb'.repeat(32), trustedRequestId, 0n))
        .to.not.be.reverted;
    });
  });

  // The global queues outrank the round robin — triggers before deploys before per-application
  // work — under the same grace rule: a head older than `selectionGrace` was visible at every
  // possible read instant, so a lower-priority update can no longer be explained by the selection
  // race and reverts PriorityQueueNotServed with the queue that must be served first.
  describe('priority queues', function () {
    const REQUEST_TYPE_DEPLOYAPP = 0;
    const REQUEST_TYPE_TRUSTPROCESS = 4;

    // Deploys a trigger application and leaves one TRUSTPROCESS request in the trigger queue.
    // Must run while every other queue is empty, so the selection view returns the trusted head.
    async function enqueueTrustedRequest(): Promise<{
      triggerAppId: bigint;
      trustedRequestId: string;
      triggerAppRoot: string;
    }> {
      const TestTrigger = await ethers.getContractFactory('TestTrigger');
      const trigger = await TestTrigger.deploy(await processorEndpoint.getAddress(), false, false);
      const deployTx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequestWithTrigger(0, '0x' + 'ab'.repeat(40), await trigger.getAddress(), {
          value: minFeePerRequest,
        });
      const parsed = (await deployTx.wait()).logs
        .map((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log);
          } catch {
            return null;
          }
        })
        .find((p: any) => p?.name === 'DeployRequestSubmitted');
      const triggerAppId: bigint = parsed.args.applicationId;
      await (
        await stateUpdate(
          triggerAppId,
          BYTES32_ZERO,
          INITIAL_STATE_ROOT,
          parsed.args.requestId,
          minFeePerRequest
        )
      ).wait();

      await (await trigger.setTrustedPayload('0xdeadbeef')).wait();
      const seedRequestId = await submitRequest(triggerAppId, '0x16');
      const triggerAppRoot = '0x' + 'aa'.repeat(32);
      await (
        await stateUpdate(
          triggerAppId,
          INITIAL_STATE_ROOT,
          triggerAppRoot,
          seedRequestId,
          minFeePerRequest
        )
      ).wait();
      await (await trigger.setTrustedPayload('0x')).wait();
      const [, triggerHead] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
      return { triggerAppId, trustedRequestId: triggerHead[0].requestId, triggerAppRoot };
    }

    it('rejects a per-application update while the trigger queue head is older than the grace period', async function () {
      const { triggerAppId, trustedRequestId, triggerAppRoot } = await enqueueTrustedRequest();
      await advanceTime(grace + 1n);
      const requestId = await submitRequest(firstAppId, '0x17');

      await expect(
        stateUpdate(
          firstAppId,
          INITIAL_STATE_ROOT,
          '0x' + '11'.repeat(32),
          requestId,
          minFeePerRequest
        )
      )
        .to.be.revertedWithCustomError(processorEndpoint, 'PriorityQueueNotServed')
        .withArgs(REQUEST_TYPE_TRUSTPROCESS);

      // Serving the trusted request unblocks the per-application lane.
      await (
        await stateUpdate(
          triggerAppId,
          triggerAppRoot,
          '0x' + '22'.repeat(32),
          trustedRequestId,
          0n
        )
      ).wait();
      await expect(
        stateUpdate(
          firstAppId,
          INITIAL_STATE_ROOT,
          '0x' + '33'.repeat(32),
          requestId,
          minFeePerRequest
        )
      ).to.not.be.reverted;
    });

    it('accepts a per-application update while the trigger queue head is within the grace period', async function () {
      await enqueueTrustedRequest();
      const requestId = await submitRequest(firstAppId, '0x18');

      await expect(
        stateUpdate(
          firstAppId,
          INITIAL_STATE_ROOT,
          '0x' + '44'.repeat(32),
          requestId,
          minFeePerRequest
        )
      ).to.not.be.reverted;
    });

    it('rejects a per-application update while the deploy queue head is older than the grace period', async function () {
      await processorEndpoint
        .connect(signers[2])
        .submitDeployRequest(0, '0x19', { value: minFeePerRequest });
      await advanceTime(grace + 1n);
      const requestId = await submitRequest(firstAppId, '0x1a');

      await expect(
        stateUpdate(
          firstAppId,
          INITIAL_STATE_ROOT,
          '0x' + '55'.repeat(32),
          requestId,
          minFeePerRequest
        )
      )
        .to.be.revertedWithCustomError(processorEndpoint, 'PriorityQueueNotServed')
        .withArgs(REQUEST_TYPE_DEPLOYAPP);
    });

    it('accepts a per-application update while the deploy queue head is within the grace period', async function () {
      await processorEndpoint
        .connect(signers[2])
        .submitDeployRequest(0, '0x1b', { value: minFeePerRequest });
      const requestId = await submitRequest(firstAppId, '0x1c');

      await expect(
        stateUpdate(
          firstAppId,
          INITIAL_STATE_ROOT,
          '0x' + '66'.repeat(32),
          requestId,
          minFeePerRequest
        )
      ).to.not.be.reverted;
    });

    it('rejects a deploy update while the trigger queue head is older than the grace period', async function () {
      await enqueueTrustedRequest();
      const deployTx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequest(0, '0x1d', { value: minFeePerRequest });
      const parsed = (await deployTx.wait()).logs
        .map((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log);
          } catch {
            return null;
          }
        })
        .find((p: any) => p?.name === 'DeployRequestSubmitted');
      await advanceTime(grace + 1n);

      await expect(
        stateUpdate(
          parsed.args.applicationId,
          BYTES32_ZERO,
          '0x' + '77'.repeat(32),
          parsed.args.requestId,
          minFeePerRequest
        )
      )
        .to.be.revertedWithCustomError(processorEndpoint, 'PriorityQueueNotServed')
        .withArgs(REQUEST_TYPE_TRUSTPROCESS);
    });

    it('accepts a trusted update while the deploy queue head is older than the grace period', async function () {
      const { triggerAppId, trustedRequestId, triggerAppRoot } = await enqueueTrustedRequest();
      await processorEndpoint
        .connect(signers[2])
        .submitDeployRequest(0, '0x1e', { value: minFeePerRequest });
      await advanceTime(grace + 1n);

      await expect(
        stateUpdate(triggerAppId, triggerAppRoot, '0x' + '88'.repeat(32), trustedRequestId, 0n)
      ).to.not.be.reverted;
    });
  });

  describe('updateSelectionGrace', function () {
    it('is restricted to the admin', async function () {
      await expect(
        processorEndpoint.connect(signers[0]).updateSelectionGrace(120)
      ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
    });

    it('updates the value and emits SelectionGraceUpdated', async function () {
      await expect(processorEndpoint.connect(signers[2]).updateSelectionGrace(120))
        .to.emit(processorEndpoint, 'SelectionGraceUpdated')
        .withArgs(120);
      expect(await processorEndpoint.selectionGrace()).to.equal(120n);
    });
  });
});
