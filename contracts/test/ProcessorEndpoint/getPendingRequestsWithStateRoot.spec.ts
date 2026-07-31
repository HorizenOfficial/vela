import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  BYTES32_ZERO,
  ETH_TOKEN,
  getRequestIdFromReceipt,
  PROTOCOL_VERSION,
  REQUEST_TYPE_PROCESS,
  REQUEST_TYPE_TRUSTPROCESS,
} from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let bootstrapApplication: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
  });

  async function submit(applicationId: bigint, payload: string) {
    const tx = await processorEndpoint.submitRequest(
      PROTOCOL_VERSION,
      applicationId,
      REQUEST_TYPE_PROCESS,
      payload,
      ETH_TOKEN,
      0,
      minFeePerRequest,
      { value: minFeePerRequest }
    );
    return getRequestIdFromReceipt(processorEndpoint, await tx.wait());
  }

  // Completes a request through the single-request path, moving the application's state root on.
  async function process(applicationId: bigint, requestId: string, prev: string, next: string) {
    await processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        applicationId,
        prev,
        next,
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        minFeePerRequest,
        0,
        '',
        '0x'
      );
  }

  function root(byte: string) {
    return '0x' + byte.repeat(32);
  }

  describe('getPendingRequestsWithStateRoot', function () {
    describe('unhappy paths', function () {
      it('returns no requests when every queue is empty', async () => {
        await bootstrapApplication(processorEndpoint);

        const [appId, requests, stateRoot] =
          await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(appId).to.equal(0n);
        expect(requests.length).to.equal(0);
        expect(stateRoot).to.equal(BYTES32_ZERO);
      });

      it('returns no requests when maxCount is zero', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        await submit(applicationId, '0x01');

        const [appId, requests, stateRoot] =
          await processorEndpoint.getPendingRequestsWithStateRoot(0);
        expect(appId).to.equal(0n);
        expect(requests.length).to.equal(0);
        expect(stateRoot).to.equal(BYTES32_ZERO);
      });

      it('returns no requests before any application is deployed', async () => {
        const [, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(requests.length).to.equal(0);
      });
    });

    describe('happy paths', function () {
      it('returns the application id, its requests in FIFO order, and its state root', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const first = await submit(applicationId, '0x01');
        const second = await submit(applicationId, '0x02');

        const [appId, requests, stateRoot] =
          await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(appId).to.equal(applicationId);
        expect(stateRoot).to.equal(INITIAL_STATE_ROOT);
        expect(requests.length).to.equal(2);
        expect(requests[0].requestId).to.equal(first);
        expect(requests[0].payload).to.equal('0x01');
        expect(requests[0].protocolVersion).to.equal(PROTOCOL_VERSION);
        expect(requests[0].applicationId).to.equal(applicationId);
        expect(requests[0].requestType).to.equal(REQUEST_TYPE_PROCESS);
        expect(requests[0].assetAmount).to.equal(0n);
        expect(requests[1].requestId).to.equal(second);
        expect(requests[1].payload).to.equal('0x02');
      });

      it('caps the result at maxCount', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const first = await submit(applicationId, '0x01');
        await submit(applicationId, '0x02');
        await submit(applicationId, '0x03');

        const [, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(2);
        expect(requests.length).to.equal(2);
        expect(requests[0].requestId).to.equal(first);
      });

      it('returns only the selected application requests: A, B, A yields both A requests', async () => {
        const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
        const { applicationId: appB } = await bootstrapApplication(processorEndpoint);

        const a1 = await submit(appA, '0x01');
        await submit(appB, '0x02');
        const a2 = await submit(appA, '0x03');

        const [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(appId).to.equal(appA);
        expect(requests.length).to.equal(2);
        expect(requests[0].requestId).to.equal(a1);
        expect(requests[1].requestId).to.equal(a2);

        // B's queue is untouched: all three requests are still pending.
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(3n);
      });

      it('alternates between applications as each takes its turn', async () => {
        const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
        const { applicationId: appB } = await bootstrapApplication(processorEndpoint);

        const a1 = await submit(appA, '0x01');
        const a2 = await submit(appA, '0x02');
        const b1 = await submit(appB, '0x03');
        const b2 = await submit(appB, '0x04');

        // A is selected first, and processing one of its requests advances the cursor to B.
        let [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
        expect(appId).to.equal(appA);
        expect(requests[0].requestId).to.equal(a1);
        await process(appA, a1, INITIAL_STATE_ROOT, root('a1'));

        [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
        expect(appId).to.equal(appB);
        expect(requests[0].requestId).to.equal(b1);
        await process(appB, b1, INITIAL_STATE_ROOT, root('b1'));

        [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
        expect(appId).to.equal(appA);
        expect(requests[0].requestId).to.equal(a2);
        await process(appA, a2, root('a1'), root('a2'));

        [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
        expect(appId).to.equal(appB);
        expect(requests[0].requestId).to.equal(b2);
      });

      it('skips an application whose queue is empty and wraps around the deployed apps', async () => {
        const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
        await bootstrapApplication(processorEndpoint); // appB: never has pending work
        const { applicationId: appC } = await bootstrapApplication(processorEndpoint);

        const a1 = await submit(appA, '0x01');
        const c1 = await submit(appC, '0x02');

        // A's turn: processing it moves the cursor to B, whose queue is empty, so C is served.
        await process(appA, a1, INITIAL_STATE_ROOT, root('a1'));
        let [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(appId).to.equal(appC);
        expect(requests[0].requestId).to.equal(c1);

        // C is last in the array: after its turn the scan wraps back to A.
        await process(appC, c1, INITIAL_STATE_ROOT, root('c1'));
        const a2 = await submit(appA, '0x03');
        [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(appId).to.equal(appA);
        expect(requests[0].requestId).to.equal(a2);
      });

      it('serves the only busy application every turn', async () => {
        const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
        await bootstrapApplication(processorEndpoint);

        const a1 = await submit(appA, '0x01');
        const a2 = await submit(appA, '0x02');

        await process(appA, a1, INITIAL_STATE_ROOT, root('a1'));
        const [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(appId).to.equal(appA);
        expect(requests.length).to.equal(1);
        expect(requests[0].requestId).to.equal(a2);
      });

      it('returns a pending deploy request alone, before any application selection', async () => {
        const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
        await submit(appA, '0x01');

        const deployTx = await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(PROTOCOL_VERSION, '0x0102', { value: minFeePerRequest });
        const deployReceipt = await deployTx.wait();
        const deployLog = deployReceipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(deployLog);

        const [appId, requests, stateRoot] =
          await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(requests.length).to.equal(1);
        expect(requests[0].requestId).to.equal(parsed.args.requestId);
        expect(appId).to.equal(parsed.args.applicationId);
        // The application does not exist yet, so it has no state root.
        expect(stateRoot).to.equal(BYTES32_ZERO);
      });

      it('returns a pending TRUSTPROCESS alone, before the deploy and application queues', async () => {
        // Deploy an application with a trigger that produces a trusted payload, so processing
        // one of its requests enqueues a TRUSTPROCESS into the global trigger queue.
        const TestTrigger = await ethers.getContractFactory('TestTrigger');
        const trigger: any = await TestTrigger.deploy(
          await processorEndpoint.getAddress(),
          false,
          false
        );
        const deployTx = await processorEndpoint
          .connect(signers[2])
          .submitDeployRequestWithTrigger(PROTOCOL_VERSION, '0x00', await trigger.getAddress(), {
            value: minFeePerRequest,
          });
        const parsed = processorEndpoint.interface.parseLog(
          (await deployTx.wait()).logs.find((log: any) => {
            try {
              return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
            } catch {
              return false;
            }
          })
        );
        const triggerAppId: bigint = parsed.args.applicationId;
        await process(triggerAppId, parsed.args.requestId, BYTES32_ZERO, INITIAL_STATE_ROOT);

        await (await trigger.setTrustedPayload('0xdeadbeef')).wait();
        const fired = await submit(triggerAppId, '0x01');
        await process(triggerAppId, fired, INITIAL_STATE_ROOT, root('11'));
        expect(await processorEndpoint.getTriggerQueueSize()).to.equal(1n);

        // A normal application also has pending work, but the TRUSTPROCESS comes first.
        const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
        await submit(appA, '0x02');

        const [appId, requests, stateRoot] =
          await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(requests.length).to.equal(1);
        expect(requests[0].requestType).to.equal(REQUEST_TYPE_TRUSTPROCESS);
        expect(appId).to.equal(triggerAppId);
        expect(stateRoot).to.equal(root('11'));
      });

      it('returns at most one request for an application with a registered trigger', async () => {
        const TestTrigger = await ethers.getContractFactory('TestTrigger');
        const trigger: any = await TestTrigger.deploy(
          await processorEndpoint.getAddress(),
          false,
          false
        );
        const deployTx = await processorEndpoint
          .connect(signers[2])
          .submitDeployRequestWithTrigger(PROTOCOL_VERSION, '0x00', await trigger.getAddress(), {
            value: minFeePerRequest,
          });
        const parsed = processorEndpoint.interface.parseLog(
          (await deployTx.wait()).logs.find((log: any) => {
            try {
              return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
            } catch {
              return false;
            }
          })
        );
        const triggerAppId: bigint = parsed.args.applicationId;
        await process(triggerAppId, parsed.args.requestId, BYTES32_ZERO, INITIAL_STATE_ROOT);

        const first = await submit(triggerAppId, '0x01');
        await submit(triggerAppId, '0x02');
        await submit(triggerAppId, '0x03');

        const [appId, requests] = await processorEndpoint.getPendingRequestsWithStateRoot(5);
        expect(appId).to.equal(triggerAppId);
        expect(requests.length).to.equal(1);
        expect(requests[0].requestId).to.equal(first);
      });
    });
  });

  describe('queue views across applications', function () {
    it('getPendingRequestsSize sums the deploy queue and every application queue', async () => {
      const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
      const { applicationId: appB } = await bootstrapApplication(processorEndpoint);

      await submit(appA, '0x01');
      await submit(appB, '0x02');
      await submit(appA, '0x03');
      await processorEndpoint
        .connect(signers[2])
        .submitDeployRequest(PROTOCOL_VERSION, '0x0102', { value: minFeePerRequest });

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(4n);
    });

    it('getPendingRequests returns every queued request across applications', async () => {
      const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
      const { applicationId: appB } = await bootstrapApplication(processorEndpoint);

      const a1 = await submit(appA, '0x01');
      const b1 = await submit(appB, '0x02');
      const a2 = await submit(appA, '0x03');

      const pending = await processorEndpoint.getPendingRequests();
      expect(pending.length).to.equal(3);
      const ids = pending.map((r: any) => r.requestId);
      expect(ids).to.have.members([a1, b1, a2]);
    });

    it('isCurrentPendingRequest is true for the head of each application queue', async () => {
      const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
      const { applicationId: appB } = await bootstrapApplication(processorEndpoint);

      const a1 = await submit(appA, '0x01');
      const a2 = await submit(appA, '0x02');
      const b1 = await submit(appB, '0x03');

      // Both applications' heads are current: the queues advance independently.
      expect(await processorEndpoint.isCurrentPendingRequest(a1)).to.equal(true);
      expect(await processorEndpoint.isCurrentPendingRequest(b1)).to.equal(true);
      expect(await processorEndpoint.isCurrentPendingRequest(a2)).to.equal(false);
      expect(await processorEndpoint.isCurrentPendingRequest(root('ff'))).to.equal(false);
    });

    it('processing an application head does not disturb the other application queue', async () => {
      const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
      const { applicationId: appB } = await bootstrapApplication(processorEndpoint);

      const a1 = await submit(appA, '0x01');
      const b1 = await submit(appB, '0x02');

      await process(appA, a1, INITIAL_STATE_ROOT, root('a1'));

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);
      expect(await processorEndpoint.isCurrentPendingRequest(b1)).to.equal(true);
      expect(await processorEndpoint.applicationStateRoots(appB)).to.equal(INITIAL_STATE_ROOT);
    });

    it('adminReset clears every application queue and refunds their deposits', async () => {
      const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
      const { applicationId: appB } = await bootstrapApplication(processorEndpoint);

      await processorEndpoint.submitRequest(
        PROTOCOL_VERSION,
        appA,
        REQUEST_TYPE_PROCESS,
        '0x01',
        ETH_TOKEN,
        7n,
        minFeePerRequest,
        { value: 7n + minFeePerRequest }
      );
      await processorEndpoint.submitRequest(
        PROTOCOL_VERSION,
        appB,
        REQUEST_TYPE_PROCESS,
        '0x02',
        ETH_TOKEN,
        9n,
        minFeePerRequest,
        { value: 9n + minFeePerRequest }
      );
      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(2n);

      await processorEndpoint.connect(signers[3]).adminReset();

      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
      expect(await processorEndpoint.appCustody(appA, ETH_TOKEN)).to.equal(0n);
      expect(await processorEndpoint.appCustody(appB, ETH_TOKEN)).to.equal(0n);
      expect(
        await processorEndpoint.pendingClaims(ETH_TOKEN, await signers[0].getAddress())
      ).to.equal(16n);
    });
  });
});
