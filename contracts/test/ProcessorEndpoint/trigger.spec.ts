import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { ETH_TOKEN, BYTES32_ZERO, getRequestIdFromReceipt, REQUEST_TYPE_PROCESS } from '../util';

describe('ProcessorEndpoint Trigger Tests', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
  });

  // HELPER FUNCTIONS FOR TESTS

  // Deploys a MockTrigger and bootstraps an app with it registered.
  // Returns the applicationId, the mock trigger contract, and the state root after deploy.
  async function bootstrapApplicationWithTrigger(
    revertOnExecute: boolean,
    revertOnPostWithdraw: boolean
  ) {
    const MockTrigger = await ethers.getContractFactory('MockTrigger');
    const mockTrigger = await MockTrigger.deploy(
      await processorEndpoint.getAddress(),
      await processorEndpoint.getAddress(), // ProcessorEndpoint also implements ITokenAllowlist
      revertOnExecute,
      revertOnPostWithdraw
    );

    // ABI-encode the trigger address to exactly 32 bytes so stateUpdate registers it
    const triggerPayload = ethers.AbiCoder.defaultAbiCoder().encode(
      ['address'],
      [await mockTrigger.getAddress()]
    );

    const deployTx = await processorEndpoint
      .connect(signers[2])
      .submitDeployRequest(0, triggerPayload, { value: minFeePerRequest });
    const deployReceipt = await deployTx.wait();

    const deployLog = deployReceipt.logs.find((log: any) => {
      try {
        return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
      } catch {
        return false;
      }
    });
    const parsed = processorEndpoint.interface.parseLog(deployLog);
    const applicationId: bigint = parsed.args.applicationId;
    const requestId: string = parsed.args.requestId;

    await processorEndpoint.connect(signers[1]).stateUpdate(
      applicationId,
      BYTES32_ZERO,
      INITIAL_STATE_ROOT,
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

    return { applicationId, mockTrigger };
  }

  // Submits a PROCESS request, processes it via stateUpdate, and returns the receipt.
  async function submitAndProcess(
    applicationId: bigint,
    currentStateRoot: string,
    nextStateRoot: string
  ) {
    const tx = await processorEndpoint
      .connect(signers[0])
      .submitRequest(
        0,
        applicationId,
        REQUEST_TYPE_PROCESS,
        '0x01',
        ETH_TOKEN,
        0n,
        minFeePerRequest,
        { value: minFeePerRequest }
      );
    const submitReceipt = await tx.wait();
    const requestId = getRequestIdFromReceipt(processorEndpoint, submitReceipt);

    const updateTx = await processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        applicationId,
        currentStateRoot,
        nextStateRoot,
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
    const updateReceipt = await updateTx.wait();
    return { requestId, updateTx, updateReceipt };
  }

  function parseTriggerEvents(receipt: any) {
    const triggerExecuted: any[] = [];
    const triggerPostWithdraw: any[] = [];
    for (const log of receipt.logs) {
      try {
        const parsed = processorEndpoint.interface.parseLog(log);
        if (parsed.name === 'TriggerExecuted') triggerExecuted.push(parsed.args);
        if (parsed.name === 'TriggerPostWithdraw') triggerPostWithdraw.push(parsed.args);
      } catch {
        // ignore logs from other contracts
      }
    }
    return { triggerExecuted, triggerPostWithdraw };
  }

  // TESTS
  describe('no trigger registered', function () {
    it('stateUpdate succeeds and emits no trigger events when no trigger is registered', async function () {
      const { applicationId } = await (
        await deployProcessorEndpointFixture()
      ).bootstrapApplication(processorEndpoint);

      const { updateTx, updateReceipt } = await submitAndProcess(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + 'ab'.repeat(32)
      );
      const { triggerExecuted, triggerPostWithdraw } = parseTriggerEvents(updateReceipt);

      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerExecuted');
      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerPostWithdraw');
      expect(triggerExecuted.length).to.equal(0);
      expect(triggerPostWithdraw.length).to.equal(0);
    });
  });

  describe('trigger registered via deploy payload', function () {
    describe('revertOnExecute=false, revertOnPostWithdraw=false', function () {
      it('stateUpdate succeeds, TriggerExecuted and TriggerPostWithdraw both report success=true', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(
          false,
          false
        );
        const triggerAddress = await mockTrigger.getAddress();

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '11'.repeat(32)
        );
        const { triggerExecuted, triggerPostWithdraw } = parseTriggerEvents(updateReceipt);

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, triggerAddress, true);
        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerPostWithdraw')
          .withArgs(applicationId, requestId, triggerAddress, true);

        expect(triggerExecuted[0].success).to.equal(true);
        expect(triggerPostWithdraw[0].success).to.equal(true);

        expect(await mockTrigger.executedRequests(requestId)).to.equal(true);
        expect(await mockTrigger.executedPostWithdraws(requestId)).to.equal(true);
      });
    });

    describe('revertOnExecute=false, revertOnPostWithdraw=true', function () {
      it('stateUpdate succeeds, TriggerExecuted success=true, TriggerPostWithdraw success=false', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(
          false,
          true
        );
        const triggerAddress = await mockTrigger.getAddress();

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '22'.repeat(32)
        );
        const { triggerExecuted, triggerPostWithdraw } = parseTriggerEvents(updateReceipt);

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, triggerAddress, true);
        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerPostWithdraw')
          .withArgs(applicationId, requestId, triggerAddress, false);

        expect(triggerExecuted[0].success).to.equal(true);
        expect(triggerPostWithdraw[0].success).to.equal(false);

        // _execute ran successfully so the mapping was written
        expect(await mockTrigger.executedRequests(requestId)).to.equal(true);
        // _postWithdraw reverted so its storage write was rolled back
        expect(await mockTrigger.executedPostWithdraws(requestId)).to.equal(false);
      });
    });

    describe('revertOnExecute=true, revertOnPostWithdraw=false', function () {
      it('stateUpdate succeeds, TriggerExecuted success=false, TriggerPostWithdraw success=true', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(
          true,
          false
        );
        const triggerAddress = await mockTrigger.getAddress();

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '33'.repeat(32)
        );
        const { triggerExecuted, triggerPostWithdraw } = parseTriggerEvents(updateReceipt);

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, triggerAddress, false);
        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerPostWithdraw')
          .withArgs(applicationId, requestId, triggerAddress, true);

        expect(triggerExecuted[0].success).to.equal(false);
        expect(triggerPostWithdraw[0].success).to.equal(true);

        // _execute reverted so nothing was written
        expect(await mockTrigger.executedRequests(requestId)).to.equal(false);
        // withdraw was still called and _postWithdraw succeeded
        expect(await mockTrigger.executedPostWithdraws(requestId)).to.equal(true);
      });
    });

    describe('revertOnExecute=true, revertOnPostWithdraw=true', function () {
      it('stateUpdate succeeds, TriggerExecuted success=false, TriggerPostWithdraw success=false', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(
          true,
          true
        );
        const triggerAddress = await mockTrigger.getAddress();

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '44'.repeat(32)
        );
        const { triggerExecuted, triggerPostWithdraw } = parseTriggerEvents(updateReceipt);

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, triggerAddress, false);
        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerPostWithdraw')
          .withArgs(applicationId, requestId, triggerAddress, false);

        expect(triggerExecuted[0].success).to.equal(false);
        expect(triggerPostWithdraw[0].success).to.equal(false);

        expect(await mockTrigger.executedRequests(requestId)).to.equal(false);
        expect(await mockTrigger.executedPostWithdraws(requestId)).to.equal(false);
      });
    });

    it('trigger revert does not affect state root update', async function () {
      const { applicationId } = await bootstrapApplicationWithTrigger(true, true);
      const nextStateRoot = '0x' + '55'.repeat(32);

      await submitAndProcess(applicationId, INITIAL_STATE_ROOT, nextStateRoot);

      expect(await processorEndpoint.applicationStateRoots(applicationId)).to.equal(nextStateRoot);
    });

    it('trigger revert does not affect RequestCompleted event', async function () {
      const { applicationId } = await bootstrapApplicationWithTrigger(true, true);

      const { requestId, updateTx } = await submitAndProcess(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + '66'.repeat(32)
      );

      await expect(updateTx)
        .to.emit(processorEndpoint, 'RequestCompleted')
        .withArgs(applicationId, requestId, minFeePerRequest, 0, 0, '');
    });
  });
});
