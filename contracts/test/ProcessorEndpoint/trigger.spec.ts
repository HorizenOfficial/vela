import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  ETH_TOKEN,
  BYTES32_ZERO,
  getRequestIdFromReceipt,
  REQUEST_TYPE_PROCESS,
  REQUEST_TYPE_TRUSTPROCESS,
} from '../util';

describe('ProcessorEndpoint Trigger Tests', function () {
  let processorEndpoint: any;
  let tokenAllowlist: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let bootstrapWithTrigger: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint, tokenAllowlist } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapWithTrigger = fixture.bootstrapApplicationWithTrigger;
  });

  // HELPER FUNCTIONS FOR TESTS

  // Deploys a TestTrigger and bootstraps an app with it registered.
  // Returns the applicationId and the trigger contract instance.
  async function bootstrapApplicationWithTrigger(
    revertOnExecute: boolean,
    revertOnPostWithdraw: boolean
  ) {
    const { trigger, applicationId } = await bootstrapWithTrigger(processorEndpoint, {
      revertOnExecute,
      revertOnPostWithdraw,
    });
    return { applicationId, mockTrigger: trigger };
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
    const triggerWithdraw: any[] = [];
    for (const log of receipt.logs) {
      try {
        const parsed = processorEndpoint.interface.parseLog(log);
        if (parsed.name === 'TriggerExecuted') triggerExecuted.push(parsed.args);
        if (parsed.name === 'TriggerWithdraw') triggerWithdraw.push(parsed.args);
      } catch {
        // ignore logs from other contracts
      }
    }
    return { triggerExecuted, triggerWithdraw };
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
      const { triggerExecuted, triggerWithdraw } = parseTriggerEvents(updateReceipt);

      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerExecuted');
      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerWithdraw');
      expect(triggerExecuted.length).to.equal(0);
      expect(triggerWithdraw.length).to.equal(0);
    });
  });

  describe('trigger registered via deploy payload', function () {
    describe('revertOnExecute=false, revertOnPostWithdraw=false', function () {
      it('stateUpdate succeeds, TriggerExecuted and TriggerWithdraw both report success=true', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '11'.repeat(32)
        );
        const { triggerExecuted, triggerWithdraw } = parseTriggerEvents(updateReceipt);
        const blockNumber = updateReceipt.blockNumber;

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, true);
        await expect(updateTx).to.emit(processorEndpoint, 'TriggerWithdraw');

        expect(triggerExecuted[0].success).to.equal(true);
        expect(triggerWithdraw[0].withdrawSuccess).to.equal(true);
        expect(triggerWithdraw[0].postWithdrawSuccess).to.equal(true);

        expect(await mockTrigger.executedInBlock(blockNumber)).to.equal(true);
        expect(await mockTrigger.executedPostWithdrawsInBlock(blockNumber)).to.equal(true);
      });
    });

    describe('revertOnExecute=false, revertOnPostWithdraw=true', function () {
      it('stateUpdate succeeds, TriggerExecuted success=true, TriggerWithdraw postWithdrawSuccess=false', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, true);

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '22'.repeat(32)
        );
        const { triggerExecuted, triggerWithdraw } = parseTriggerEvents(updateReceipt);
        const blockNumber = updateReceipt.blockNumber;

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, true);
        await expect(updateTx).to.emit(processorEndpoint, 'TriggerWithdraw');

        expect(triggerExecuted[0].success).to.equal(true);
        expect(triggerWithdraw[0].withdrawSuccess).to.equal(true);
        expect(triggerWithdraw[0].postWithdrawSuccess).to.equal(false);

        // _execute ran successfully so the mapping was written
        expect(await mockTrigger.executedInBlock(blockNumber)).to.equal(true);
        // getTrustProcessPayload reverted so its storage write was rolled back
        expect(await mockTrigger.executedPostWithdrawsInBlock(blockNumber)).to.equal(false);
      });
    });

    describe('revertOnExecute=true, revertOnPostWithdraw=false', function () {
      it('stateUpdate succeeds, TriggerExecuted success=false, TriggerWithdraw postWithdrawSuccess=true', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(true, false);

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '33'.repeat(32)
        );
        const { triggerExecuted, triggerWithdraw } = parseTriggerEvents(updateReceipt);
        const blockNumber = updateReceipt.blockNumber;

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, false);
        await expect(updateTx).to.emit(processorEndpoint, 'TriggerWithdraw');

        expect(triggerExecuted[0].success).to.equal(false);
        expect(triggerWithdraw[0].withdrawSuccess).to.equal(true);
        expect(triggerWithdraw[0].postWithdrawSuccess).to.equal(true);

        // _execute reverted so nothing was written
        expect(await mockTrigger.executedInBlock(blockNumber)).to.equal(false);
        // withdraw was still called and getTrustProcessPayload succeeded
        expect(await mockTrigger.executedPostWithdrawsInBlock(blockNumber)).to.equal(true);
      });
    });

    describe('revertOnExecute=true, revertOnPostWithdraw=true', function () {
      it('stateUpdate succeeds, TriggerExecuted success=false, TriggerWithdraw postWithdrawSuccess=false', async function () {
        const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(true, true);

        const { requestId, updateTx, updateReceipt } = await submitAndProcess(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + '44'.repeat(32)
        );
        const { triggerExecuted, triggerWithdraw } = parseTriggerEvents(updateReceipt);
        const blockNumber = updateReceipt.blockNumber;

        await expect(updateTx)
          .to.emit(processorEndpoint, 'TriggerExecuted')
          .withArgs(applicationId, requestId, false);
        await expect(updateTx).to.emit(processorEndpoint, 'TriggerWithdraw');

        expect(triggerExecuted[0].success).to.equal(false);
        expect(triggerWithdraw[0].withdrawSuccess).to.equal(true);
        expect(triggerWithdraw[0].postWithdrawSuccess).to.equal(false);

        expect(await mockTrigger.executedInBlock(blockNumber)).to.equal(false);
        expect(await mockTrigger.executedPostWithdrawsInBlock(blockNumber)).to.equal(false);
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

  describe('trigger queue submission (trusted requests created in stateUpdate)', function () {
    it('a non-empty getTrustProcessPayload payload enqueues a TRUSTPROCESS into the trigger queue during stateUpdate', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);

      const trustedPayload = '0xdeadbeef';
      await (await mockTrigger.setTrustedPayload(trustedPayload)).wait();

      // Processing a normal request runs _invokeTrigger inside stateUpdate; the
      // trigger's getTrustProcessPayload returns the payload → a TRUSTPROCESS is enqueued.
      await submitAndProcess(applicationId, INITIAL_STATE_ROOT, '0x' + '33'.repeat(32));

      expect(await processorEndpoint.getTriggerQueueSize()).to.equal(1n);

      // Trigger queue takes priority — the enqueued trusted request is next.
      const [, pending] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
      const request = pending[0];
      expect(request.applicationId).to.equal(applicationId);
      expect(request.requestType).to.equal(REQUEST_TYPE_TRUSTPROCESS);
      expect(request.payload).to.equal(trustedPayload);
      expect(request.sender).to.equal(await mockTrigger.getAddress());
    });

    it('an empty getTrustProcessPayload payload enqueues nothing', async function () {
      const { applicationId } = await bootstrapApplicationWithTrigger(false, false);
      // default trustedPayload is empty ("")
      await submitAndProcess(applicationId, INITIAL_STATE_ROOT, '0x' + '44'.repeat(32));
      expect(await processorEndpoint.getTriggerQueueSize()).to.equal(0n);
    });

    it('trigger queue request is served before the normal queue', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      await (await mockTrigger.setTrustedPayload('0x02')).wait();

      // Processing the first request enqueues the trusted request in the trigger queue.
      await submitAndProcess(applicationId, INITIAL_STATE_ROOT, '0x' + '55'.repeat(32));

      // Enqueue a normal request afterwards.
      const normalTx = await processorEndpoint
        .connect(signers[0])
        .submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          ETH_TOKEN,
          0n,
          minFeePerRequest,
          {
            value: minFeePerRequest,
          }
        );
      await normalTx.wait();

      // Trigger queue has priority — the trusted request is the current pending one.
      const [, pending] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
      const request = pending[0];
      expect(request.requestType).to.equal(REQUEST_TYPE_TRUSTPROCESS);
      expect(request.sender).to.equal(await mockTrigger.getAddress());
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
      await tokenAllowlist.connect(signers[2]).addAllowedToken(await tokenA.getAddress());
      await tokenAllowlist.connect(signers[2]).addAllowedToken(await tokenB.getAddress());
    });

    // Submits one ETH-asset request and two ERC-20 requests to build up appCustody.
    // Returns the requestId of the first (ETH) request, which is next in queue for stateUpdate.
    async function buildCustody(applicationId: bigint) {
      const ethTx = await processorEndpoint
        .connect(signers[0])
        .submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          ETH_TOKEN,
          ETH_ASSET,
          minFeePerRequest,
          {
            value: ETH_ASSET + minFeePerRequest,
          }
        );
      const firstRequestId = getRequestIdFromReceipt(processorEndpoint, await ethTx.wait());

      const tokenAAddr = await tokenA.getAddress();
      await tokenA.mint(await signers[0].getAddress(), TOKEN_A_ASSET);
      await tokenA.connect(signers[0]).approve(await processorEndpoint.getAddress(), TOKEN_A_ASSET);
      await processorEndpoint
        .connect(signers[0])
        .submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x02',
          tokenAAddr,
          TOKEN_A_ASSET,
          minFeePerRequest,
          {
            value: minFeePerRequest,
          }
        );

      const tokenBAddr = await tokenB.getAddress();
      await tokenB.mint(await signers[0].getAddress(), TOKEN_B_ASSET);
      await tokenB.connect(signers[0]).approve(await processorEndpoint.getAddress(), TOKEN_B_ASSET);
      await processorEndpoint
        .connect(signers[0])
        .submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x03',
          tokenBAddr,
          TOKEN_B_ASSET,
          minFeePerRequest,
          {
            value: minFeePerRequest,
          }
        );

      return { firstRequestId };
    }

    it('trigger holds ETH and both ERC-20 assets during _execute (unshield via withdrawalRequests)', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      const { firstRequestId } = await buildCustody(applicationId);
      const triggerAddress = await mockTrigger.getAddress();
      const tokenAAddr = await tokenA.getAddress();
      const tokenBAddr = await tokenB.getAddress();

      await processorEndpoint.connect(signers[1]).stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + 'aa'.repeat(32),
        firstRequestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [
          { tokenAddress: ETH_TOKEN, receiver: triggerAddress, amount: ETH_ASSET },
          { tokenAddress: tokenAAddr, receiver: triggerAddress, amount: TOKEN_A_ASSET },
          { tokenAddress: tokenBAddr, receiver: triggerAddress, amount: TOKEN_B_ASSET },
        ],
        0,
        minFeePerRequest,
        0,
        '',
        '0x'
      );

      expect(await mockTrigger.capturedBalances(ETH_TOKEN)).to.equal(ETH_ASSET);
      expect(await mockTrigger.capturedBalances(tokenAAddr)).to.equal(TOKEN_A_ASSET);
      expect(await mockTrigger.capturedBalances(tokenBAddr)).to.equal(TOKEN_B_ASSET);
    });

    it('appCustody is updated with reshielded amounts after trigger withdraw', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      const { firstRequestId } = await buildCustody(applicationId);
      const triggerAddress = await mockTrigger.getAddress();
      const tokenAAddr = await tokenA.getAddress();
      const tokenBAddr = await tokenB.getAddress();

      await processorEndpoint.connect(signers[1]).stateUpdate(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + 'bb'.repeat(32),
        firstRequestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [
          { tokenAddress: ETH_TOKEN, receiver: triggerAddress, amount: ETH_ASSET },
          { tokenAddress: tokenAAddr, receiver: triggerAddress, amount: TOKEN_A_ASSET },
          { tokenAddress: tokenBAddr, receiver: triggerAddress, amount: TOKEN_B_ASSET },
        ],
        0,
        minFeePerRequest,
        0,
        '',
        '0x'
      );

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

    it('if trigger is not in withdrawalRequests, it receives no ETH or ERC-20', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);

      // Submit only ERC-20 request but route withdrawal to a different receiver (not trigger)
      const tokenAAddr = await tokenA.getAddress();
      await tokenA.mint(await signers[0].getAddress(), TOKEN_A_ASSET);
      await tokenA.connect(signers[0]).approve(await processorEndpoint.getAddress(), TOKEN_A_ASSET);
      const erc20Tx = await processorEndpoint
        .connect(signers[0])
        .submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          tokenAAddr,
          TOKEN_A_ASSET,
          minFeePerRequest,
          {
            value: minFeePerRequest,
          }
        );
      const requestId = getRequestIdFromReceipt(processorEndpoint, await erc20Tx.wait());

      // No withdrawalRequests to trigger — trigger receives nothing
      await processorEndpoint
        .connect(signers[1])
        .stateUpdate(
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
          '0x'
        );

      // Trigger received nothing — capturedBalances are 0
      expect(await mockTrigger.capturedBalances(ETH_TOKEN)).to.equal(0n);
      expect(await mockTrigger.capturedBalances(tokenAAddr)).to.equal(0n);
    });
  });

  describe('trigger registered via 3-arg submitDeployRequest overload (descriptor + trigger)', function () {
    // A non-32-byte payload simulates a real WASM deploy descriptor, proving the
    // trigger is registered from the explicit `trigger` arg and NOT from the
    // legacy 32-byte-payload path.
    const descriptorPayload = '0x' + 'ab'.repeat(40); // 40 bytes, != 32

    async function bootstrapWithTriggerParam(
      revertOnExecute: boolean,
      revertOnPostWithdraw: boolean
    ) {
      const TestTrigger = await ethers.getContractFactory('TestTrigger');
      const mockTrigger: any = await TestTrigger.deploy(
        await processorEndpoint.getAddress(),
        revertOnExecute,
        revertOnPostWithdraw
      );

      const deployTx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequestWithTrigger(0, descriptorPayload, await mockTrigger.getAddress(), {
          value: minFeePerRequest,
        });
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

      await processorEndpoint
        .connect(signers[1])
        .stateUpdate(
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

    it('registers the trigger from the explicit arg even when the payload is a (non-32-byte) descriptor', async function () {
      const { applicationId, mockTrigger } = await bootstrapWithTriggerParam(false, false);
      const triggerAddr = await mockTrigger.getAddress();

      expect(await processorEndpoint.triggerContracts(applicationId)).to.equal(triggerAddr);
      expect(await processorEndpoint.triggersToAppIds(triggerAddr)).to.equal(applicationId);
    });

    it('invokes execute()/withdraw() on a subsequent PROCESS request for the app', async function () {
      const { applicationId } = await bootstrapWithTriggerParam(false, false);

      const { requestId, updateTx } = await submitAndProcess(
        applicationId,
        INITIAL_STATE_ROOT,
        '0x' + '22'.repeat(32)
      );

      await expect(updateTx)
        .to.emit(processorEndpoint, 'TriggerExecuted')
        .withArgs(applicationId, requestId, true);
      await expect(updateTx).to.emit(processorEndpoint, 'TriggerWithdraw');
    });

    it('address(0) trigger registers nothing (equivalent to the 2-arg overload)', async function () {
      const deployTx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequestWithTrigger(0, descriptorPayload, ethers.ZeroAddress, {
          value: minFeePerRequest,
        });
      const deployReceipt = await deployTx.wait();
      const parsed = processorEndpoint.interface.parseLog(
        deployReceipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
          } catch {
            return false;
          }
        })
      );
      const applicationId: bigint = parsed.args.applicationId;
      const requestId: string = parsed.args.requestId;

      await processorEndpoint
        .connect(signers[1])
        .stateUpdate(
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

      expect(await processorEndpoint.triggerContracts(applicationId)).to.equal(ethers.ZeroAddress);
    });
  });

  describe('eager trigger registration: validation and cleanup', function () {
    const descriptorPayload = '0x' + 'ab'.repeat(40);

    async function deployTrigger() {
      const TestTrigger = await ethers.getContractFactory('TestTrigger');
      return (await TestTrigger.deploy(await processorEndpoint.getAddress(), false, false)) as any;
    }

    function deployInfoFromReceipt(receipt: any): { applicationId: bigint; requestId: string } {
      const log = receipt.logs.find((l: any) => {
        try {
          return processorEndpoint.interface.parseLog(l)?.name === 'DeployRequestSubmitted';
        } catch {
          return false;
        }
      });
      const parsed = processorEndpoint.interface.parseLog(log);
      return { applicationId: parsed.args.applicationId, requestId: parsed.args.requestId };
    }

    it('registers the trigger eagerly at submit time (before any stateUpdate)', async function () {
      const triggerAddr = await (await deployTrigger()).getAddress();
      const tx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequestWithTrigger(0, descriptorPayload, triggerAddr, {
          value: minFeePerRequest,
        });
      const { applicationId } = deployInfoFromReceipt(await tx.wait());
      expect(await processorEndpoint.triggerContracts(applicationId)).to.equal(triggerAddr);
      expect(await processorEndpoint.triggersToAppIds(triggerAddr)).to.equal(applicationId);
    });

    it('reverts with TriggerAlreadyRegistered when the same trigger is submitted twice', async function () {
      const triggerAddr = await (await deployTrigger()).getAddress();
      await processorEndpoint
        .connect(signers[2])
        .submitDeployRequestWithTrigger(0, descriptorPayload, triggerAddr, {
          value: minFeePerRequest,
        });
      await expect(
        processorEndpoint
          .connect(signers[2])
          .submitDeployRequestWithTrigger(0, descriptorPayload, triggerAddr, {
            value: minFeePerRequest,
          })
      ).to.be.revertedWithCustomError(processorEndpoint, 'TriggerAlreadyRegistered');
    });

    it('reverts with TriggerCannotBeEOA when the trigger is an EOA', async function () {
      const eoa = await signers[5].getAddress();
      await expect(
        processorEndpoint
          .connect(signers[2])
          .submitDeployRequestWithTrigger(0, descriptorPayload, eoa, { value: minFeePerRequest })
      ).to.be.revertedWithCustomError(processorEndpoint, 'TriggerCannotBeEOA');
    });

    it('rolls back the eager registration when the deploy fails, so the trigger can be reused', async function () {
      const triggerAddr = await (await deployTrigger()).getAddress();
      const tx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequestWithTrigger(0, descriptorPayload, triggerAddr, {
          value: minFeePerRequest,
        });
      const { applicationId, requestId } = deployInfoFromReceipt(await tx.wait());
      expect(await processorEndpoint.triggersToAppIds(triggerAddr)).to.equal(applicationId);

      // Fail the deploy: error stateUpdate (errorCode != NO_ERROR, state left unchanged at zero).
      await processorEndpoint.connect(signers[1]).stateUpdate(
        applicationId,
        BYTES32_ZERO,
        BYTES32_ZERO,
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        0,
        3, // ErrorCode.APPLICATION_ALREADY_DEPLOYED
        'deploy failed',
        '0x'
      );

      // Registration rolled back on failure.
      expect(await processorEndpoint.triggerContracts(applicationId)).to.equal(ethers.ZeroAddress);
      expect(await processorEndpoint.triggersToAppIds(triggerAddr)).to.equal(0n);

      // The same trigger address can now be reused by a fresh deploy.
      await expect(
        processorEndpoint
          .connect(signers[2])
          .submitDeployRequestWithTrigger(0, descriptorPayload, triggerAddr, {
            value: minFeePerRequest,
          })
      ).to.not.be.reverted;
    });

    it('adminResetApps clears the trigger mappings for the reset app', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      const triggerAddr = await mockTrigger.getAddress();
      expect(await processorEndpoint.triggerContracts(applicationId)).to.equal(triggerAddr);

      await processorEndpoint.connect(signers[3]).adminResetApps([applicationId]);

      expect(await processorEndpoint.triggerContracts(applicationId)).to.equal(ethers.ZeroAddress);
      expect(await processorEndpoint.triggersToAppIds(triggerAddr)).to.equal(0n);
    });

    it('adminReset drains the trigger queue', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      // Trigger emits a non-empty trusted payload → a TRUSTPROCESS is enqueued during stateUpdate.
      await (await mockTrigger.setTrustedPayload('0x02')).wait();
      await submitAndProcess(applicationId, INITIAL_STATE_ROOT, '0x' + '22'.repeat(32));
      expect(await processorEndpoint.getTriggerQueueSize()).to.equal(1n);

      await processorEndpoint.connect(signers[3]).adminReset();
      expect(await processorEndpoint.getTriggerQueueSize()).to.equal(0n);
    });
  });

  // The trigger must only act on SUCCESSFUL stateUpdates: the error branch returns
  // before _invokeTrigger, so a failed request can neither invoke the trigger nor
  // enqueue (another) TRUSTPROCESS. These cover that "failure is inert" guarantee.
  describe('failure path: a failed stateUpdate leaves the trigger inert', function () {
    // Drives an already-pending request into the error path: state unchanged,
    // empty events/withdrawals, errorCode != NO_ERROR.
    async function failStateUpdate(applicationId: bigint, requestId: string) {
      const currentRoot = await processorEndpoint.applicationStateRoots(applicationId);
      const updateTx = await processorEndpoint.connect(signers[1]).stateUpdate(
        applicationId,
        currentRoot,
        currentRoot,
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        0,
        1, // errorCode != NO_ERROR
        'err',
        '0x'
      );
      return { updateTx, updateReceipt: await updateTx.wait() };
    }

    it('a failed PROCESS request does not invoke the registered trigger', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      // Arm a non-empty payload: were the trigger (wrongly) invoked, it would
      // enqueue a TRUSTPROCESS — so a queue size of 0 proves it was not.
      await (await mockTrigger.setTrustedPayload('0xdeadbeef')).wait();

      const submitTx = await processorEndpoint
        .connect(signers[0])
        .submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          ETH_TOKEN,
          0n,
          minFeePerRequest,
          {
            value: minFeePerRequest,
          }
        );
      const requestId = getRequestIdFromReceipt(processorEndpoint, await submitTx.wait());

      const { updateTx, updateReceipt } = await failStateUpdate(applicationId, requestId);
      const { triggerExecuted, triggerWithdraw } = parseTriggerEvents(updateReceipt);

      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerExecuted');
      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerWithdraw');
      expect(triggerExecuted.length).to.equal(0);
      expect(triggerWithdraw.length).to.equal(0);
      expect(await mockTrigger.executedInBlock(updateReceipt.blockNumber)).to.equal(false);

      // No TRUSTPROCESS enqueued, and the request completes as FAILED.
      expect(await processorEndpoint.getTriggerQueueSize()).to.equal(0n);
      await expect(updateTx)
        .to.emit(processorEndpoint, 'RequestCompleted')
        .withArgs(applicationId, requestId, minFeePerRequest, 1, 1, 'err');
    });

    it('a failed TRUSTPROCESS is dequeued and neither re-invokes the trigger nor enqueues another', async function () {
      const { applicationId, mockTrigger } = await bootstrapApplicationWithTrigger(false, false);
      await (await mockTrigger.setTrustedPayload('0xdeadbeef')).wait();

      // A successful PROCESS enqueues exactly one TRUSTPROCESS.
      await submitAndProcess(applicationId, INITIAL_STATE_ROOT, '0x' + '33'.repeat(32));
      expect(await processorEndpoint.getTriggerQueueSize()).to.equal(1n);

      const [, pendingTrusted] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
      const req = pendingTrusted[0];
      expect(req.requestType).to.equal(REQUEST_TYPE_TRUSTPROCESS);
      const trustedRequestId: string = req.requestId;

      const { updateTx, updateReceipt } = await failStateUpdate(applicationId, trustedRequestId);
      const { triggerExecuted, triggerWithdraw } = parseTriggerEvents(updateReceipt);

      // The failed TRUSTPROCESS does not re-invoke the trigger...
      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerExecuted');
      await expect(updateTx).not.to.emit(processorEndpoint, 'TriggerWithdraw');
      expect(triggerExecuted.length).to.equal(0);
      expect(triggerWithdraw.length).to.equal(0);
      expect(await mockTrigger.executedInBlock(updateReceipt.blockNumber)).to.equal(false);

      // ...enqueues no replacement and is itself dequeued.
      expect(await processorEndpoint.getTriggerQueueSize()).to.equal(0n);
      expect(await processorEndpoint.isCurrentPendingRequest(trustedRequestId)).to.equal(false);

      // Marked FAILED with no fee (trigger-queue requests carry maxFeeValue 0).
      await expect(updateTx)
        .to.emit(processorEndpoint, 'RequestCompleted')
        .withArgs(applicationId, trustedRequestId, 0n, 1, 1, 'err');
    });
  });
});
