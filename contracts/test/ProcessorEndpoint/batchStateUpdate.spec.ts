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
import {
  ethSignBatchStateUpdate,
  ethSignStateUpdate,
  stateUpdateEntryHash,
} from '../../scripts/util';

// batchStateUpdate applies several update payloads for one application in a single transaction,
// covered by one TEE signature over the concatenated entry hashes. These tests run against
// NoAttestationTeeAuthenticator rather than MockTeeAuthenticator: the mock accepts every signature,
// which would leave the batch digest — the one genuinely new piece of cryptography here —
// unverified.
describe('ProcessorEndpoint batchStateUpdate', function () {
  const EMPTY_EVENTS = { events: [] as string[], subTypes: [] as string[] };
  const ERROR_REQUEST_FUNC_FAILED = 6; // Structs.ErrorCode.REQUEST_FUNC_FAILED

  let fixture: any;
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let teeSigner: Signer;
  let operator: Signer;
  let appId: bigint;

  function root(byte: string): string {
    return '0x' + byte.repeat(32);
  }

  type Entry = {
    prevStateRoot: string;
    newStateRoot: string;
    processedRequestId: string;
    userEvents: { events: string[]; subTypes: string[] };
    appEvents: { events: string[]; subTypes: string[] };
    withdrawalRequests: any[];
    refund: bigint;
    applicationFees: bigint;
    errorCode: number;
    errorMsg: string;
  };

  // A successful entry moving the root from prev to next, paying the whole fee to the application.
  function successEntry(requestId: string, prev: string, next: string): Entry {
    return {
      prevStateRoot: prev,
      newStateRoot: next,
      processedRequestId: requestId,
      userEvents: EMPTY_EVENTS,
      appEvents: EMPTY_EVENTS,
      withdrawalRequests: [],
      refund: 0n,
      applicationFees: minFeePerRequest,
      errorCode: 0,
      errorMsg: '',
    };
  }

  // An error entry: state unchanged, no events, no withdrawals — the shape the contract requires
  // of a signed error payload.
  function errorEntry(requestId: string, at: string, errorMsg = 'boom'): Entry {
    return {
      prevStateRoot: at,
      newStateRoot: at,
      processedRequestId: requestId,
      userEvents: EMPTY_EVENTS,
      appEvents: EMPTY_EVENTS,
      withdrawalRequests: [],
      refund: 0n,
      applicationFees: 0n,
      errorCode: ERROR_REQUEST_FUNC_FAILED,
      errorMsg,
    };
  }

  function entryHashOf(applicationId: bigint, e: Entry): string {
    return stateUpdateEntryHash(
      applicationId,
      e.prevStateRoot,
      e.newStateRoot,
      e.processedRequestId,
      e.userEvents.events,
      e.userEvents.subTypes,
      e.appEvents.events,
      e.appEvents.subTypes,
      e.withdrawalRequests,
      e.refund,
      e.applicationFees,
      e.errorCode,
      e.errorMsg
    );
  }

  function batchHashes(applicationId: bigint, entries: Entry[]): string[] {
    return entries.map((e) => entryHashOf(applicationId, e));
  }

  async function submitBatch(
    endpoint: any,
    applicationId: bigint,
    entries: Entry[],
    signer: Signer = teeSigner
  ) {
    const signature = await ethSignBatchStateUpdate(signer, batchHashes(applicationId, entries));
    return endpoint.connect(operator).batchStateUpdate(applicationId, entries, signature);
  }

  // An endpoint whose signatures are really verified, with `apps` bootstrapped applications.
  async function deployWithRealSignatures(apps: number) {
    const f = await deployProcessorEndpointFixture();
    const NoAttestationTeeAuthenticator = await ethers.getContractFactory(
      'NoAttestationTeeAuthenticator'
    );
    const teeAuthenticator = await NoAttestationTeeAuthenticator.deploy(
      await f.signers[0].getAddress(),
      await f.signers[4].getAddress(),
      '0x' + '11'.repeat(133)
    );
    const endpoint = await f.deployProcessorEndpointWith(
      await teeAuthenticator.getAddress(),
      await f.authorityRegistry.getAddress(),
      await (await f.deployTokenAllowlist()).getAddress()
    );

    const appIds: bigint[] = [];
    for (let i = 0; i < apps; i++) {
      const { applicationId } = await f.bootstrapApplication(endpoint, f.signers[4]);
      appIds.push(applicationId);
    }
    return { f, endpoint, teeAuthenticator, appIds };
  }

  async function submitRequest(
    endpoint: any,
    applicationId: bigint,
    payload: string
  ): Promise<string> {
    const tx = await endpoint
      .connect(signers[0])
      .submitRequest(
        PROTOCOL_VERSION,
        applicationId,
        REQUEST_TYPE_PROCESS,
        payload,
        ETH_TOKEN,
        0n,
        minFeePerRequest,
        { value: minFeePerRequest }
      );
    return getRequestIdFromReceipt(endpoint, await tx.wait());
  }

  // Queues `count` requests and returns their ids plus the chain of state roots to move through
  // them: roots[0] is the current root and roots[i + 1] is the root after entry i.
  async function queueRequests(endpoint: any, applicationId: bigint, count: number) {
    const requestIds: string[] = [];
    for (let i = 0; i < count; i++) {
      requestIds.push(
        await submitRequest(endpoint, applicationId, '0x' + (i + 1).toString(16).padStart(2, '0'))
      );
    }
    const roots = [INITIAL_STATE_ROOT];
    for (let i = 0; i < count; i++) roots.push(root((0xa0 + i).toString(16)));
    return { requestIds, roots };
  }

  beforeEach(async function () {
    fixture = await deployProcessorEndpointFixture();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    teeSigner = signers[4];
    operator = signers[1];
    const deployed = await deployWithRealSignatures(1);
    processorEndpoint = deployed.endpoint;
    appId = deployed.appIds[0];
  });

  describe('happy paths', function () {
    it('processes every entry in one transaction and chains the state root', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 3);
      const entries = requestIds.map((id, i) => successEntry(id, roots[i], roots[i + 1]));

      const receipt = await (await submitBatch(processorEndpoint, appId, entries)).wait();

      expect(await processorEndpoint.applicationStateRoots(appId)).to.equal(roots[3]);
      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0);
      for (const id of requestIds) {
        expect(await processorEndpoint.isCurrentPendingRequest(id)).to.equal(false);
      }

      // Per-entry events are emitted inside the loop, exactly as on the single-request path, so
      // indexers see the same stream — just all in one block.
      const stateRootUpdates = receipt.logs
        .map((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log);
          } catch {
            return null;
          }
        })
        .filter((p: any) => p && p.name === 'StateRootUpdate');
      expect(stateRootUpdates.length).to.equal(3);
      for (let i = 0; i < 3; i++) {
        expect(stateRootUpdates[i].args.oldStateRoot).to.equal(roots[i]);
        expect(stateRootUpdates[i].args.newStateRoot).to.equal(roots[i + 1]);
        expect(stateRootUpdates[i].args.requestId).to.equal(requestIds[i]);
      }
    });

    it('emits BatchProcessed with the entry hashes the signature covered', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 2);
      const entries = requestIds.map((id, i) => successEntry(id, roots[i], roots[i + 1]));

      await expect(submitBatch(processorEndpoint, appId, entries))
        .to.emit(processorEndpoint, 'BatchProcessed')
        .withArgs(appId, batchHashes(appId, entries));
    });

    it('accepts a one-entry batch signed the way stateUpdate expects', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 1);
      const entry = successEntry(requestIds[0], roots[0], roots[1]);

      // Signed with the single-request helper, not the batch one: a one-entry batch digest is
      // byte-identical to the single-request digest, so one scheme serves both paths.
      const singleSignature = await ethSignStateUpdate(
        teeSigner,
        appId,
        entry.prevStateRoot,
        entry.newStateRoot,
        entry.processedRequestId,
        [],
        [],
        [],
        [],
        [],
        entry.refund,
        entry.applicationFees
      );

      await expect(
        processorEndpoint.connect(operator).batchStateUpdate(appId, [entry], singleSignature)
      ).to.not.be.reverted;
      expect(await processorEndpoint.applicationStateRoots(appId)).to.equal(roots[1]);
    });

    it('continues from the unchanged root after an error entry mid-batch', async function () {
      const { requestIds } = await queueRequests(processorEndpoint, appId, 3);
      const r1 = root('b1');
      const r2 = root('b2');
      const entries = [
        successEntry(requestIds[0], INITIAL_STATE_ROOT, r1),
        errorEntry(requestIds[1], r1),
        successEntry(requestIds[2], r1, r2),
      ];

      const tx = await submitBatch(processorEndpoint, appId, entries);

      await expect(tx)
        .to.emit(processorEndpoint, 'RequestCompleted')
        .withArgs(appId, requestIds[1], minFeePerRequest, 1, ERROR_REQUEST_FUNC_FAILED, 'boom');
      expect(await processorEndpoint.applicationStateRoots(appId)).to.equal(r2);
      expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0);
    });

    it('processes a pending deploy as a one-entry batch', async function () {
      const { endpoint } = await deployWithRealSignatures(0);
      const deployReceipt = await (
        await endpoint.connect(signers[2]).submitDeployRequest(PROTOCOL_VERSION, '0x00', {
          value: minFeePerRequest,
        })
      ).wait();
      const parsed = endpoint.interface.parseLog(
        deployReceipt.logs.find((log: any) => {
          try {
            return endpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
          } catch {
            return false;
          }
        })
      );
      const deployAppId: bigint = parsed.args.applicationId;
      const entry = successEntry(parsed.args.requestId, BYTES32_ZERO, INITIAL_STATE_ROOT);

      await submitBatch(endpoint, deployAppId, [entry]);

      expect(await endpoint.applicationStateRoots(deployAppId)).to.equal(INITIAL_STATE_ROOT);
      expect(await endpoint.getDeployedAppIds()).to.deep.equal([deployAppId]);
    });
  });

  describe('state root chaining', function () {
    it("reverts when the first entry's prevStateRoot does not match storage", async function () {
      const { requestIds } = await queueRequests(processorEndpoint, appId, 2);
      const entries = [
        successEntry(requestIds[0], root('cc'), root('c1')),
        successEntry(requestIds[1], root('c1'), root('c2')),
      ];

      await expect(submitBatch(processorEndpoint, appId, entries)).to.be.revertedWithCustomError(
        processorEndpoint,
        'InvalidStateRoot'
      );
    });

    it('reverts when the chain between two entries is broken', async function () {
      const { requestIds } = await queueRequests(processorEndpoint, appId, 2);
      const entries = [
        successEntry(requestIds[0], INITIAL_STATE_ROOT, root('d1')),
        // Chains from a root no entry produced.
        successEntry(requestIds[1], root('dd'), root('d2')),
      ];

      await expect(submitBatch(processorEndpoint, appId, entries)).to.be.revertedWithCustomError(
        processorEndpoint,
        'InvalidStateRoot'
      );
    });

    it('reverts when an entry claims a state change the queue order does not allow', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 2);
      // Entries in reverse queue order: the first entry is not the queue head.
      const entries = [
        successEntry(requestIds[1], roots[0], roots[1]),
        successEntry(requestIds[0], roots[1], roots[2]),
      ];

      await expect(submitBatch(processorEndpoint, appId, entries)).to.be.revertedWithCustomError(
        processorEndpoint,
        'InvalidRequestId'
      );
    });
  });

  describe('batch signature', function () {
    it('reverts when the signature does not cover every entry', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 3);
      const entries = requestIds.map((id, i) => successEntry(id, roots[i], roots[i + 1]));

      // Signed over the first two entry hashes only. The personal_sign length prefix commits to
      // 32*N, so a short batch digest cannot pass for a longer one.
      const signature = await ethSignBatchStateUpdate(
        teeSigner,
        batchHashes(appId, entries).slice(0, 2)
      );

      await expect(
        processorEndpoint.connect(operator).batchStateUpdate(appId, entries, signature)
      ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidSignature');
    });

    it('reverts when an entry is altered after signing', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 2);
      const entries = requestIds.map((id, i) => successEntry(id, roots[i], roots[i + 1]));
      const signature = await ethSignBatchStateUpdate(teeSigner, batchHashes(appId, entries));

      const tampered = entries.map((e) => ({ ...e }));
      tampered[1].newStateRoot = root('ee');

      await expect(
        processorEndpoint.connect(operator).batchStateUpdate(appId, tampered, signature)
      ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidSignature');
    });

    it('reverts when the batch is signed by the wrong key', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 2);
      const entries = requestIds.map((id, i) => successEntry(id, roots[i], roots[i + 1]));

      await expect(
        submitBatch(processorEndpoint, appId, entries, signers[5])
      ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidSignature');
    });

    it('reverts when the entry order is swapped', async function () {
      // The batch digest is the concatenation of the entry hashes, so it is order-sensitive even
      // when the set of entries is unchanged.
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 2);
      const entries = requestIds.map((id, i) => successEntry(id, roots[i], roots[i + 1]));
      const signature = await ethSignBatchStateUpdate(teeSigner, [
        entryHashOf(appId, entries[1]),
        entryHashOf(appId, entries[0]),
      ]);

      await expect(
        processorEndpoint.connect(operator).batchStateUpdate(appId, entries, signature)
      ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidSignature');
    });
  });

  describe('batch composition rules', function () {
    it('reverts with EmptyBatch when no entries are supplied', async function () {
      await expect(
        processorEndpoint.connect(operator).batchStateUpdate(appId, [], '0x')
      ).to.be.revertedWithCustomError(processorEndpoint, 'EmptyBatch');
    });

    it('reverts with BatchNotAllowed for an application with a registered trigger', async function () {
      const { endpoint, appIds } = await deployWithTriggerApplication();
      const triggerAppId = appIds[0];
      const requestIds = [
        await submitRequest(endpoint, triggerAppId, '0x01'),
        await submitRequest(endpoint, triggerAppId, '0x02'),
      ];
      const entries = [
        successEntry(requestIds[0], INITIAL_STATE_ROOT, root('f1')),
        successEntry(requestIds[1], root('f1'), root('f2')),
      ];

      await expect(submitBatch(endpoint, triggerAppId, entries)).to.be.revertedWithCustomError(
        endpoint,
        'BatchNotAllowed'
      );
    });

    it('accepts a one-entry batch for an application with a registered trigger', async function () {
      const { endpoint, appIds } = await deployWithTriggerApplication();
      const triggerAppId = appIds[0];
      const requestId = await submitRequest(endpoint, triggerAppId, '0x01');
      const newRoot = root('f9');

      await expect(
        submitBatch(endpoint, triggerAppId, [successEntry(requestId, INITIAL_STATE_ROOT, newRoot)])
      )
        .to.emit(endpoint, 'TriggerExecuted')
        .withArgs(triggerAppId, requestId, true);
      expect(await endpoint.applicationStateRoots(triggerAppId)).to.equal(newRoot);
    });

    it('reverts with BatchNotAllowed when a pending deploy is batched with other entries', async function () {
      const { endpoint } = await deployWithRealSignatures(1);
      const existingApp = (await endpoint.getDeployedAppIds())[0];
      const deployReceipt = await (
        await endpoint.connect(signers[2]).submitDeployRequest(PROTOCOL_VERSION, '0x01', {
          value: minFeePerRequest,
        })
      ).wait();
      const parsed = endpoint.interface.parseLog(
        deployReceipt.logs.find((log: any) => {
          try {
            return endpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
          } catch {
            return false;
          }
        })
      );
      const deployAppId: bigint = parsed.args.applicationId;
      const otherRequestId = await submitRequest(endpoint, existingApp, '0x02');

      await expect(
        submitBatch(endpoint, deployAppId, [
          successEntry(parsed.args.requestId, BYTES32_ZERO, INITIAL_STATE_ROOT),
          successEntry(otherRequestId, INITIAL_STATE_ROOT, root('a1')),
        ])
      ).to.be.revertedWithCustomError(endpoint, 'BatchNotAllowed');
    });

    it('reverts for a caller without UPDATE_STATUS_ROLE', async function () {
      const { requestIds, roots } = await queueRequests(processorEndpoint, appId, 1);
      const entries = [successEntry(requestIds[0], roots[0], roots[1])];
      const signature = await ethSignBatchStateUpdate(teeSigner, batchHashes(appId, entries));

      await expect(
        processorEndpoint.connect(signers[0]).batchStateUpdate(appId, entries, signature)
      ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
    });
  });

  describe('round-robin enforcement', function () {
    it("reverts with ApplicationNotSelected for an application that is not the scan's result", async function () {
      const { endpoint, appIds } = await deployWithRealSignatures(2);
      // Both applications have pending work and the cursor is at appIds[0], so appIds[1] is out
      // of turn once appIds[0]'s head has aged past the grace period.
      await submitRequest(endpoint, appIds[0], '0x01');
      const outOfTurn = [
        await submitRequest(endpoint, appIds[1], '0x02'),
        await submitRequest(endpoint, appIds[1], '0x03'),
      ];
      const grace: bigint = await endpoint.selectionGrace();
      await ethers.provider.send('evm_increaseTime', ['0x' + (grace * 2n).toString(16)]);
      await ethers.provider.send('evm_mine', []);

      await expect(
        submitBatch(endpoint, appIds[1], [
          successEntry(outOfTurn[0], INITIAL_STATE_ROOT, root('91')),
          successEntry(outOfTurn[1], root('91'), root('92')),
        ])
      )
        .to.be.revertedWithCustomError(endpoint, 'ApplicationNotSelected')
        .withArgs(appIds[0]);
    });

    it('advances the cursor once per batch, so the other application gets the next turn', async function () {
      const { endpoint, appIds } = await deployWithRealSignatures(2);
      const first = [
        await submitRequest(endpoint, appIds[0], '0x01'),
        await submitRequest(endpoint, appIds[0], '0x02'),
      ];
      const second = await submitRequest(endpoint, appIds[1], '0x03');

      await submitBatch(endpoint, appIds[0], [
        successEntry(first[0], INITIAL_STATE_ROOT, root('81')),
        successEntry(first[1], root('81'), root('82')),
      ]);

      // The whole batch consumed a single turn: the next selection is the second application.
      const [selected] = await endpoint.getPendingRequestsWithStateRoot(5);
      expect(selected).to.equal(appIds[1]);
      await expect(
        submitBatch(endpoint, appIds[1], [successEntry(second, INITIAL_STATE_ROOT, root('83'))])
      ).to.not.be.reverted;
    });
  });

  describe('gas', function () {
    it('costs less than the same requests submitted one by one', async function () {
      const n = 5;

      const batched = await deployWithRealSignatures(1);
      const batchApp = batched.appIds[0];
      const queued = await queueRequests(batched.endpoint, batchApp, n);
      const entries = queued.requestIds.map((id, i) =>
        successEntry(id, queued.roots[i], queued.roots[i + 1])
      );
      const batchReceipt = await (await submitBatch(batched.endpoint, batchApp, entries)).wait();

      const single = await deployWithRealSignatures(1);
      const singleApp = single.appIds[0];
      const singleQueued = await queueRequests(single.endpoint, singleApp, n);
      let singleGas = 0n;
      for (let i = 0; i < n; i++) {
        const e = successEntry(
          singleQueued.requestIds[i],
          singleQueued.roots[i],
          singleQueued.roots[i + 1]
        );
        const signature = await ethSignStateUpdate(
          teeSigner,
          singleApp,
          e.prevStateRoot,
          e.newStateRoot,
          e.processedRequestId,
          [],
          [],
          [],
          [],
          [],
          e.refund,
          e.applicationFees
        );
        const receipt = await (
          await single.endpoint
            .connect(operator)
            .stateUpdate(
              singleApp,
              e.prevStateRoot,
              e.newStateRoot,
              e.processedRequestId,
              EMPTY_EVENTS,
              EMPTY_EVENTS,
              [],
              e.refund,
              e.applicationFees,
              e.errorCode,
              e.errorMsg,
              signature
            )
        ).wait();
        singleGas += receipt.gasUsed;
      }

      // Execution gas only — the (n - 1) × 21,000 intrinsic cost the batch also saves is not
      // charged to gasUsed here, so the real saving is larger than this comparison shows.
      console.log(
        `        batchStateUpdate(${n}): ${batchReceipt.gasUsed} gas vs ${n} × stateUpdate: ${singleGas} gas`
      );
      expect(batchReceipt.gasUsed).to.be.lessThan(singleGas);
    });
  });

  // An endpoint with real signature verification and one application that has a trigger contract
  // registered at deploy time.
  async function deployWithTriggerApplication() {
    const { f, endpoint } = await deployWithRealSignatures(0);
    const { trigger, applicationId } = await f.bootstrapApplicationWithTrigger(endpoint, {
      teeSigner: f.signers[4],
    });
    return { endpoint, trigger, appIds: [applicationId] };
  }
});
