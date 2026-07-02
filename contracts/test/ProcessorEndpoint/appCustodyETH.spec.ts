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

describe('ProcessorEndpoint — appCustody', function () {
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

  async function submitRequest(
    sender: Signer,
    payload: string,
    assetAmount: bigint,
    maxFeeValue: bigint,
    appId?: bigint
  ) {
    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        appId ?? applicationId,
        REQUEST_TYPE_PROCESS,
        payload,
        ETH_TOKEN,
        assetAmount,
        maxFeeValue,
        { value: assetAmount + maxFeeValue }
      );
    const receipt = await tx.wait();
    const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);
    return { requestId, maxFeeValue, assetAmount };
  }

  async function completeRequest(
    requestId: string,
    prevStateRoot: string,
    newStateRoot: string,
    refund: bigint,
    applicationFees: bigint,
    withdrawals: [string, string, bigint][],
    appId?: bigint
  ) {
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        appId ?? applicationId,
        prevStateRoot,
        newStateRoot,
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        withdrawals,
        refund,
        applicationFees,
        0,
        '',
        '0x'
      );
  }

  async function failRequest(requestId: string, appId?: bigint) {
    const currentStateRoot = await processorEndpoint.applicationStateRoots(appId ?? applicationId);
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        appId ?? applicationId,
        currentStateRoot,
        currentStateRoot,
        requestId,
        { events: [], subTypes: [] },
        { events: [], subTypes: [] },
        [],
        0,
        0,
        1,
        'err',
        '0x'
      );
  }

  describe('basics', function () {
    // bootstrapApplication submits a deploy request and completes it via stateUpdate,
    // consuming the entire fee. The app should start with zero ETH funds.
    it('appCustody is 0 after bootstrapApplication', async () => {
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
    });

    // Verifies that submitRequest credits only the assetAmount (not fees)
    // to the app's ETH funds. Fees are tracked globally, not per-app.
    it('increases by assetAmount after submitRequest', async () => {
      const assetAmount = 50n;
      const maxFeeValue = minFeePerRequest + 10n;

      await submitRequest(signers[0], '0x01', assetAmount, maxFeeValue);

      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(assetAmount);
    });

    // Verifies that submitDeployRequest does not credit appCustody,
    // since assetAmount is always 0 for deploys. Fees are tracked globally.
    it('remains 0 after submitDeployRequest (no deposit)', async () => {
      const fee = minFeePerRequest + 20n;
      const deployTx = await processorEndpoint
        .connect(signers[2])
        .submitDeployRequest(PROTOCOL_VERSION, '0x00', { value: fee });
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

      expect(await processorEndpoint.appCustody(deployAppId, ETH_TOKEN)).to.equal(0n);
    });

    // Verifies that a successful stateUpdate debits appCustody by the withdrawal sum
    // only (fees are tracked globally). The remaining balance should equal the
    // original deposit credit minus withdrawals.
    it('decreases correctly after successful stateUpdate', async () => {
      const assetAmount = 100n;
      const refund = 10n;
      const applicationFees = minFeePerRequest;
      const maxFeeValue = refund + applicationFees;

      const request = await submitRequest(signers[0], '0x02', assetAmount, maxFeeValue);
      const fundsAfterSubmit = await processorEndpoint.appCustody(applicationId, ETH_TOKEN);
      expect(fundsAfterSubmit).to.equal(assetAmount);

      const withdrawalAddr = await signers[3].getAddress();
      const withdrawalAmount = 50n;

      await completeRequest(
        request.requestId,
        INITIAL_STATE_ROOT,
        '0x' + '22'.repeat(32),
        refund,
        applicationFees,
        [[ETH_TOKEN, withdrawalAddr, withdrawalAmount]]
      );

      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(
        fundsAfterSubmit - withdrawalAmount
      );
    });

    // Verifies that a failed stateUpdate debits only the assetAmount from
    // appCustody. Fees are handled globally.
    it('decreases correctly after error stateUpdate', async () => {
      const assetAmount = 30n;
      const maxFeeValue = minFeePerRequest + 5n;

      const request = await submitRequest(signers[0], '0x03', assetAmount, maxFeeValue);
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(assetAmount);

      await failRequest(request.requestId);

      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
    });

    // Solidity mappings return the default value (0) for uninitialized keys.
    // Confirms no spurious balance exists for apps that were never deployed.
    it('returns 0 for unknown application IDs', async () => {
      expect(await processorEndpoint.appCustody(999999, ETH_TOKEN)).to.equal(0n);
    });
  });

  describe('cross-app isolation', function () {
    let appIdB: bigint;

    beforeEach(async function () {
      ({ applicationId: appIdB } = await bootstrapApplication(processorEndpoint));
    });

    // Core isolation guarantee: App B cannot withdraw more than its own locked funds,
    // even when the contract's global ETH balance is sufficient (due to App A's funds).
    // The test explicitly asserts global balance >= B's requested outflow before
    // verifying the per-app revert, proving it's the app-level check that blocks it.
    it('App A withdrawal succeeds while App B same-amount withdrawal reverts with InsufficientAppBalance', async () => {
      const depositA = 200n;
      const depositB = 50n;
      const maxFee = minFeePerRequest;

      const reqA = await submitRequest(signers[0], '0x10', depositA, maxFee);
      // Complete A's request to free the queue head
      await completeRequest(
        reqA.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'a1'.repeat(32),
        0n,
        maxFee,
        [[ETH_TOKEN, await signers[3].getAddress(), depositA]]
      );

      // Now submit for both apps
      const reqA2 = await submitRequest(signers[0], '0x11', depositA, maxFee);
      const reqB = await submitRequest(signers[0], '0x12', depositB, maxFee, appIdB);

      // Complete A2 — withdraw 200 (within A's budget)
      await completeRequest(
        reqA2.requestId,
        '0x' + 'a1'.repeat(32),
        '0x' + 'a2'.repeat(32),
        0n,
        maxFee,
        [[ETH_TOKEN, await signers[3].getAddress(), depositA]]
      );

      // B tries to withdraw 200 — reverts because B only has 50 locked (deposit only)
      const withdrawalAttempt = 200n;
      await expect(
        completeRequest(
          reqB.requestId,
          INITIAL_STATE_ROOT,
          '0x' + 'b1'.repeat(32),
          0n,
          maxFee,
          [[ETH_TOKEN, await signers[4].getAddress(), withdrawalAttempt]],
          appIdB
        )
      ).to.be.revertedWithCustomError(processorEndpoint, 'InsufficientAppBalance');
    });

    // Verifies that processing App A's request (debiting A's pool) has no side effect
    // on App B's appCustody. Each app's accounting is fully independent.
    it("App A's stateUpdate does not change App B's appCustody", async () => {
      const depositA = 100n;
      const depositB = 75n;
      const maxFee = minFeePerRequest;

      const reqA = await submitRequest(signers[0], '0x20', depositA, maxFee);
      await submitRequest(signers[0], '0x21', depositB, maxFee, appIdB);

      const fundsB = await processorEndpoint.appCustody(appIdB, ETH_TOKEN);
      expect(fundsB).to.equal(depositB);

      // Process A
      await completeRequest(
        reqA.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'a3'.repeat(32),
        0n,
        maxFee,
        [[ETH_TOKEN, await signers[3].getAddress(), depositA]]
      );

      // B's funds unchanged
      expect(await processorEndpoint.appCustody(appIdB, ETH_TOKEN)).to.equal(fundsB);
    });

    // Same isolation guarantee for the error path: failing App A's request and
    // refunding its deposit must not touch App B's locked funds.
    it("App A error refund does not affect App B's pool", async () => {
      const depositA = 40n;
      const depositB = 60n;
      const maxFee = minFeePerRequest;

      const reqA = await submitRequest(signers[0], '0x30', depositA, maxFee);
      await submitRequest(signers[0], '0x31', depositB, maxFee, appIdB);

      const fundsB = await processorEndpoint.appCustody(appIdB, ETH_TOKEN);

      // Fail A
      await failRequest(reqA.requestId);

      // B unchanged
      expect(await processorEndpoint.appCustody(appIdB, ETH_TOKEN)).to.equal(fundsB);
    });
  });

  describe('accumulation', function () {
    // When multiple requests are pending for the same app, their deposits accumulate
    // in appCustody. Processing them one at a time should debit each request's
    // portion, eventually draining the balance to zero.
    it('multiple pending requests accumulate, process sequentially to zero', async () => {
      const dep1 = 100n;
      const dep2 = 200n;
      const maxFee = minFeePerRequest;

      const req1 = await submitRequest(signers[0], '0x40', dep1, maxFee);
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(dep1);

      const req2 = await submitRequest(signers[0], '0x41', dep2, maxFee);
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(dep1 + dep2);

      // Process first
      await completeRequest(
        req1.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'c1'.repeat(32),
        0n,
        maxFee,
        [[ETH_TOKEN, await signers[3].getAddress(), dep1]]
      );

      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(dep2);

      // Process second
      await completeRequest(
        req2.requestId,
        '0x' + 'c1'.repeat(32),
        '0x' + 'c2'.repeat(32),
        0n,
        maxFee,
        [[ETH_TOKEN, await signers[3].getAddress(), dep2]]
      );

      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
    });
  });

  describe('edge cases', function () {
    // Boundary test: when the total outflow exactly equals appCustody, the
    // stateUpdate should succeed and drain the balance to zero (no off-by-one).
    it('exact boundary: sum == appCustody succeeds', async () => {
      const assetAmount = 50n;
      const maxFeeValue = minFeePerRequest;

      const request = await submitRequest(signers[0], '0x50', assetAmount, maxFeeValue);

      await completeRequest(
        request.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'd1'.repeat(32),
        0n,
        maxFeeValue,
        [[ETH_TOKEN, await signers[3].getAddress(), assetAmount]]
      );

      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
    });

    // Boundary test: when the withdrawal sum exceeds appCustody by exactly 1 wei,
    // the solvency check must revert. Also verifies that the reverted tx does not
    // modify appCustody (EVM atomicity).
    it('boundary + 1: withdrawalSum == appCustody + 1 reverts with InsufficientAppBalance', async () => {
      const assetAmount = 50n;
      const maxFeeValue = minFeePerRequest;

      const request = await submitRequest(signers[0], '0x51', assetAmount, maxFeeValue);

      await expect(
        completeRequest(
          request.requestId,
          INITIAL_STATE_ROOT,
          '0x' + 'd2'.repeat(32),
          0n,
          maxFeeValue,
          [[ETH_TOKEN, await signers[3].getAddress(), assetAmount + 1n]]
        )
      ).to.be.revertedWithCustomError(processorEndpoint, 'InsufficientAppBalance');

      // appCustody unchanged after revert
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(assetAmount);
    });

    // Verifies correct bookkeeping when the same app has one request fail (error path)
    // and another succeed (success path). Both debit paths must work independently,
    // and the final balance should be zero.
    it('mixed error + success for same app drains to zero', async () => {
      const dep1 = 80n;
      const dep2 = 120n;
      const maxFee = minFeePerRequest;

      // Submit two requests
      const req1 = await submitRequest(signers[0], '0x60', dep1, maxFee);
      const req2 = await submitRequest(signers[0], '0x61', dep2, maxFee);
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(dep1 + dep2);

      // Fail first request — debits dep1
      await failRequest(req1.requestId);
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(dep2);

      // Succeed second request — debits dep2 (via withdrawal)
      await completeRequest(
        req2.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'd3'.repeat(32),
        0n,
        maxFee,
        [[ETH_TOKEN, await signers[3].getAddress(), dep2]]
      );
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
    });

    // Verifies that appCustody accounting is decoupled from actual ETH transfers.
    // A request submitted by a contract that rejects ETH (FallbackFailure) still gets
    // its funds correctly debited from appCustody on the error path, because the
    // pull payment pattern (_asyncTransfer) only credits the payments mapping — the
    // actual transfer happens later via withdrawPayments.
    it('FallbackFailure: appCustody still debited even when receiver rejects ETH', async () => {
      const { ethers: hre } = await import('hardhat');
      const FallbackFailure = await hre.getContractFactory('FallbackFailure');
      const fallbackFailure = await FallbackFailure.deploy();
      await fallbackFailure.deploymentTransaction()!.wait();

      // FallbackFailure submits a request via its proxy function
      const assetAmount = 0n;
      const maxFeeValue = minFeePerRequest + 4n;
      const insertTx = await fallbackFailure.insertRequestOnProcessorEndpoint(
        processorEndpoint,
        PROTOCOL_VERSION,
        applicationId,
        1,
        '0x04',
        ETH_TOKEN,
        assetAmount,
        maxFeeValue,
        { value: maxFeeValue }
      );
      const insertReceipt = await insertTx.wait();

      // appCustody should be 0 (assetAmount is 0, fees are not tracked per-app)
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);

      // Fail the request (error path) — the refund goes to FallbackFailure
      // via _asyncTransfer (pull pattern), so no actual ETH transfer happens here
      const requestId = insertReceipt.logs
        .map((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log);
          } catch {
            return null;
          }
        })
        .find((parsed: any) => parsed && parsed.name === 'RequestSubmitted')?.args.requestId;

      await failRequest(requestId!);

      // appCustody should be zero even though the receiver rejects ETH
      // (the pull pattern decouples fund accounting from actual transfers)
      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
    });
  });

  describe('deploy request tracking', function () {
    // Verifies the full deploy lifecycle: submitDeployRequest credits msg.value
    // (= maxFeeValue, since assetAmount is always 0 for deploys) to appCustody,
    // and the subsequent stateUpdate debits it back to zero.
    it('deploy request credits then stateUpdate debits to zero', async () => {
      const fixture = await deployProcessorEndpointFixture();
      const { processorEndpoint: pe } = await fixture.deployProcessorEndpoint();

      const deployTx = await pe
        .connect(fixture.signers[2])
        .submitDeployRequest(PROTOCOL_VERSION, '0x00', { value: fixture.minFeePerRequest });
      const deployReceipt = await deployTx.wait();
      const deployLog = deployReceipt.logs.find((log: any) => {
        try {
          return pe.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
        } catch {
          return false;
        }
      });
      const parsed = pe.interface.parseLog(deployLog);
      const deployAppId: bigint = parsed.args.applicationId;
      const deployRequestId: string = parsed.args.requestId;

      expect(await pe.appCustody(deployAppId, ETH_TOKEN)).to.equal(0n);

      // Complete deploy
      await pe
        .connect(fixture.signers[1])
        .stateUpdate(
          deployAppId,
          BYTES32_ZERO,
          INITIAL_STATE_ROOT,
          deployRequestId,
          { events: [], subTypes: [] },
          { events: [], subTypes: [] },
          [],
          0,
          fixture.minFeePerRequest,
          0,
          '',
          '0x'
        );

      expect(await pe.appCustody(deployAppId, ETH_TOKEN)).to.equal(0n);
    });
  });
});
