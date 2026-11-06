import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { ETH_TOKEN, BYTES32_ZERO, getRequestIdFromReceipt, REQUEST_TYPE_PROCESS } from '../util';

// A request lives in exactly one queue, determined by its requestType: TRUSTPROCESS in the trigger
// queue, DEPLOYAPP in the deploy queue, everything else in its application's queue. stateUpdate
// resolves the queue from the type rather than searching all three, and `dequeueHead` does not
// verify the id it removes — so a mis-routed removal destroys another queue's head instead of
// reverting. Every test here keeps more than one queue non-empty, which is the only situation in
// which a routing mistake is observable.
describe('ProcessorEndpoint queue routing', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let bootstrapApplicationWithTrigger: any;
  let applicationId: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplicationWithTrigger = fixture.bootstrapApplicationWithTrigger;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  function parseEvent(receipt: any, name: string) {
    for (const log of receipt.logs) {
      try {
        const parsed = processorEndpoint.interface.parseLog(log);
        if (parsed?.name === name) return parsed;
      } catch {
        continue;
      }
    }
    throw new Error(`${name} not found`);
  }

  async function submitProcessRequest(appId: bigint, payload: string): Promise<string> {
    const tx = await processorEndpoint
      .connect(signers[0])
      .submitRequest(0, appId, REQUEST_TYPE_PROCESS, payload, ETH_TOKEN, 0n, minFeePerRequest, {
        value: minFeePerRequest,
      });
    return getRequestIdFromReceipt(processorEndpoint, await tx.wait());
  }

  async function submitDeploy(payload: string) {
    const tx = await processorEndpoint
      .connect(signers[2])
      .submitDeployRequest(0, payload, { value: minFeePerRequest });
    const parsed = parseEvent(await tx.wait(), 'DeployRequestSubmitted');
    return {
      requestId: parsed.args.requestId as string,
      applicationId: parsed.args.applicationId as bigint,
    };
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

  it('processing an application request leaves a pending deploy intact', async function () {
    const requestId = await submitProcessRequest(applicationId, '0x01');
    const deploy = await submitDeploy('0x01');
    expect(await processorEndpoint.getPendingRequestsSize()).to.equal(2n);

    await (
      await stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + '11'.repeat(32),
        requestId,
        minFeePerRequest
      )
    ).wait();

    expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);
    expect(await processorEndpoint.isCurrentPendingRequest(requestId)).to.equal(false);
    // The deploy must still be queued, with its stored request untouched.
    expect(await processorEndpoint.isCurrentPendingRequest(deploy.requestId)).to.equal(true);
    expect((await processorEndpoint.requestById(deploy.requestId)).requestId).to.equal(
      deploy.requestId
    );
  });

  it('processing a deploy leaves a pending application request intact', async function () {
    const requestId = await submitProcessRequest(applicationId, '0x02');
    const deploy = await submitDeploy('0x02');

    await (
      await stateUpdate(
        deploy.applicationId,
        BYTES32_ZERO,
        '0x' + '22'.repeat(32),
        deploy.requestId,
        minFeePerRequest
      )
    ).wait();

    expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);
    expect(await processorEndpoint.isCurrentPendingRequest(deploy.requestId)).to.equal(false);
    expect(await processorEndpoint.isCurrentPendingRequest(requestId)).to.equal(true);
    expect((await processorEndpoint.requestById(requestId)).requestId).to.equal(requestId);
  });

  it('processing a TRUSTPROCESS leaves the deploy and application queues intact', async function () {
    // An application with a registered trigger, so that processing one of its requests runs
    // _invokeTrigger and enqueues a TRUSTPROCESS into the trigger queue.
    const { trigger, applicationId: triggerAppId } =
      await bootstrapApplicationWithTrigger(processorEndpoint);

    await (await trigger.setTrustedPayload('0xdeadbeef')).wait();
    const seedRequestId = await submitProcessRequest(triggerAppId, '0x03');
    const rootA = '0x' + '33'.repeat(32);
    await (
      await stateUpdate(triggerAppId, INITIAL_STATE_ROOT, rootA, seedRequestId, minFeePerRequest)
    ).wait();
    expect(await processorEndpoint.getTriggerQueueSize()).to.equal(1n);

    // Stop the trigger from producing another trusted request, so the trigger queue is expected
    // to be empty once the pending TRUSTPROCESS is processed.
    await (await trigger.setTrustedPayload('0x')).wait();

    // Fill the other two queues before processing the trusted request.
    const appRequestId = await submitProcessRequest(applicationId, '0x04');
    const deploy = await submitDeploy('0x04');
    const [, triggerHead] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
    const trustedRequestId: string = triggerHead[0].requestId;

    // TRUSTPROCESS requests carry no fee (maxFeeValue == 0).
    await (
      await stateUpdate(triggerAppId, rootA, '0x' + '44'.repeat(32), trustedRequestId, 0n)
    ).wait();

    expect(await processorEndpoint.getTriggerQueueSize()).to.equal(0n);
    expect(await processorEndpoint.getPendingRequestsSize()).to.equal(2n);
    expect(await processorEndpoint.isCurrentPendingRequest(appRequestId)).to.equal(true);
    expect(await processorEndpoint.isCurrentPendingRequest(deploy.requestId)).to.equal(true);
  });

  // DEPLOYAPP is enum value 0, so an unknown requestId reads back as a zeroed PendingRequest that
  // looks like a deploy. Routing must not confuse it with the actual head of the deploy queue.
  it('an unknown requestId is not current pending while the deploy queue is non-empty', async function () {
    const deploy = await submitDeploy('0x05');
    const unknownRequestId = '0x' + 'ee'.repeat(32);

    expect(await processorEndpoint.isCurrentPendingRequest(unknownRequestId)).to.equal(false);
    expect(await processorEndpoint.isCurrentPendingRequest(BYTES32_ZERO)).to.equal(false);
    await expect(
      stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + '55'.repeat(32),
        unknownRequestId,
        minFeePerRequest
      )
    ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');

    // The deploy queue head is untouched by the rejected update.
    expect(await processorEndpoint.isCurrentPendingRequest(deploy.requestId)).to.equal(true);
    expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);
  });

  it('reports the head of every non-empty queue as current pending', async function () {
    const requestId = await submitProcessRequest(applicationId, '0x06');
    const deploy = await submitDeploy('0x06');

    expect(await processorEndpoint.isCurrentPendingRequest(requestId)).to.equal(true);
    expect(await processorEndpoint.isCurrentPendingRequest(deploy.requestId)).to.equal(true);
  });

  it('a queued but non-head application request is rejected while other queues are non-empty', async function () {
    const firstRequestId = await submitProcessRequest(applicationId, '0x07');
    const secondRequestId = await submitProcessRequest(applicationId, '0x08');
    await submitDeploy('0x07');

    expect(await processorEndpoint.isCurrentPendingRequest(secondRequestId)).to.equal(false);
    await expect(
      stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + '66'.repeat(32),
        secondRequestId,
        minFeePerRequest
      )
    ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestId');
    expect(await processorEndpoint.isCurrentPendingRequest(firstRequestId)).to.equal(true);
  });
});
