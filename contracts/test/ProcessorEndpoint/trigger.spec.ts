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

  // Deploys a TestTrigger and bootstraps an app with it registered.
  // Returns the applicationId, the mock trigger contract, and the state root after deploy.
  async function bootstrapApplicationWithTrigger(
    revertOnExecute: boolean,
    revertOnPostWithdraw: boolean,
  ) {
    const TestTrigger = await ethers.getContractFactory('TestTrigger');
    const mockTrigger: any = await TestTrigger.deploy(
      await processorEndpoint.getAddress(),
      await processorEndpoint.getAddress(), // ProcessorEndpoint also implements ITokenAllowlist
      revertOnExecute,
      revertOnPostWithdraw,
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

  describe('unshield and reshield', function () {
    let tokenA: any;
    let tokenB: any;
    const ETH_ASSET = 1000n;
    const TOKEN_A_ASSET = 500n;
    const TOKEN_B_ASSET = 250n;

    beforeEach(async function () {
      const MockERC20 = await ethers.getContractFactory('MockERC20');
      tokenA = await MockERC20.deploy('Token A', 'TKA', 18);
      tokenB = await MockERC20.deploy('Token B', 'TKB', 18);
      await processorEndpoint.connect(signers[2]).addAllowedToken(await tokenA.getAddress());
      await processorEndpoint.connect(signers[2]).addAllowedToken(await tokenB.getAddress());
    });

    // Submits one ETH-asset request and two ERC-20 requests to build up appCustody.
    // Returns the requestId of the first (ETH) request, which is next in queue for stateUpdate.
    async function buildCustody(applicationId: bigint) {
      const ethTx = await processorEndpoint
        .connect(signers[0])
        .submitRequest(0, applicationId, REQUEST_TYPE_PROCESS, '0x01', ETH_TOKEN, ETH_ASSET, minFeePerRequest, {
          value: ETH_ASSET + minFeePerRequest,
        });
      const firstRequestId = getRequestIdFromReceipt(processorEndpoint, await ethTx.wait());

      const tokenAAddr = await tokenA.getAddress();
      await tokenA.mint(await signers[0].getAddress(), TOKEN_A_ASSET);
      await tokenA.connect(signers[0]).approve(await processorEndpoint.getAddress(), TOKEN_A_ASSET);
      await processorEndpoint
        .connect(signers[0])
        .submitRequest(0, applicationId, REQUEST_TYPE_PROCESS, '0x02', tokenAAddr, TOKEN_A_ASSET, minFeePerRequest, {
          value: minFeePerRequest,
        });

      const tokenBAddr = await tokenB.getAddress();
      await tokenB.mint(await signers[0].getAddress(), TOKEN_B_ASSET);
      await tokenB.connect(signers[0]).approve(await processorEndpoint.getAddress(), TOKEN_B_ASSET);
      await processorEndpoint
        .connect(signers[0])
        .submitRequest(0, applicationId, REQUEST_TYPE_PROCESS, '0x03', tokenBAddr, TOKEN_B_ASSET, minFeePerRequest, {
          value: minFeePerRequest,
        });

      return { firstRequestId };
    }

    it('trigger holds ETH and both ERC-20 assets during _execute (unshield)', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      const { firstRequestId } = await buildCustody(applicationId);

      await processorEndpoint.connect(signers[1]).stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + 'aa'.repeat(32),
        firstRequestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        minFeePerRequest,
        0,
        '',
        '0x',
      );

      expect(await mockTrigger.capturedBalances(ETH_TOKEN)).to.equal(ETH_ASSET);
      expect(await mockTrigger.capturedBalances(await tokenA.getAddress())).to.equal(TOKEN_A_ASSET);
      expect(await mockTrigger.capturedBalances(await tokenB.getAddress())).to.equal(TOKEN_B_ASSET);
    });

    it('appCustody is updated with reshielded amounts after trigger withdraw', async function () {
      const { applicationId } = await bootstrapApplicationWithTrigger(false, false);
      const { firstRequestId } = await buildCustody(applicationId);

      await processorEndpoint.connect(signers[1]).stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + 'bb'.repeat(32),
        firstRequestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        minFeePerRequest,
        0,
        '',
        '0x',
      );

      const tokenAAddr = await tokenA.getAddress();
      const tokenBAddr = await tokenB.getAddress();

      // TestTrigger returns all assets; appCustody is set to returned amounts
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(ETH_ASSET);
      expect(await processorEndpoint.appCustody(applicationId, tokenAAddr)).to.equal(TOKEN_A_ASSET);
      expect(await processorEndpoint.appCustody(applicationId, tokenBAddr)).to.equal(TOKEN_B_ASSET);

      // totalAppCustody has the correct amounts (reshield accounting)
      expect(await processorEndpoint.totalAppCustody(ETH_TOKEN)).to.equal(ETH_ASSET);
      expect(await processorEndpoint.totalAppCustody(tokenAAddr)).to.equal(TOKEN_A_ASSET);
      expect(await processorEndpoint.totalAppCustody(tokenBAddr)).to.equal(TOKEN_B_ASSET);

      // ProcessorEndpoint physically holds the ERC-20 tokens after the round-trip
      expect(await tokenA.balanceOf(await processorEndpoint.getAddress())).to.equal(TOKEN_A_ASSET);
      expect(await tokenB.balanceOf(await processorEndpoint.getAddress())).to.equal(TOKEN_B_ASSET);
    });

    it('if trigger retains no ETH custody, unshield skips ETH transfer', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);

      // Submit only ERC-20 requests (no ETH assetAmount — ETH custody stays at 0)
      const tokenAAddr = await tokenA.getAddress();
      await tokenA.mint(await signers[0].getAddress(), TOKEN_A_ASSET);
      await tokenA.connect(signers[0]).approve(await processorEndpoint.getAddress(), TOKEN_A_ASSET);
      const erc20Tx = await processorEndpoint
        .connect(signers[0])
        .submitRequest(0, applicationId, REQUEST_TYPE_PROCESS, '0x01', tokenAAddr, TOKEN_A_ASSET, minFeePerRequest, {
          value: minFeePerRequest,
        });
      const requestId = getRequestIdFromReceipt(processorEndpoint, await erc20Tx.wait());

      await processorEndpoint.connect(signers[1]).stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + 'cc'.repeat(32),
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        minFeePerRequest,
        0,
        '',
        '0x',
      );

      // Trigger had no ETH (nothing was unshielded for ETH)
      expect(await mockTrigger.capturedBalances(ETH_TOKEN)).to.equal(0n);
      // ERC-20 was unshielded
      expect(await mockTrigger.capturedBalances(tokenAAddr)).to.equal(TOKEN_A_ASSET);
    });
  });
});
