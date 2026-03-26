import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { BYTES32_ZERO, getRequestIdFromReceipt } from '../util';

describe('ProcessorEndpoint — appLockedFunds', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;
  let bootstrapApplication: any;

  const PROTOCOL_VERSION = 0;
  const REQUEST_TYPE = 1;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  async function submitRequest(
    sender: Signer,
    payload: string,
    depositAmount: bigint,
    maxFeeValue: bigint,
    appId?: bigint
  ) {
    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        appId ?? applicationId,
        REQUEST_TYPE,
        payload,
        depositAmount,
        maxFeeValue,
        { value: depositAmount + maxFeeValue }
      );
    const receipt = await tx.wait();
    const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);
    return { requestId, maxFeeValue, depositAmount };
  }

  async function completeRequest(
    requestId: string,
    prevStateRoot: string,
    newStateRoot: string,
    refund: bigint,
    applicationFees: bigint,
    withdrawals: [string, bigint][],
    appId?: bigint
  ) {
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        appId ?? applicationId,
        prevStateRoot,
        newStateRoot,
        requestId,
        [],
        [],
        withdrawals,
        refund,
        applicationFees,
        0,
        '',
        '0x'
      );
  }

  async function failRequest(requestId: string, appId?: bigint) {
    const currentStateRoot = await processorEndpoint.applicationStateRoots(
      appId ?? applicationId
    );
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        appId ?? applicationId,
        currentStateRoot,
        currentStateRoot,
        requestId,
        [],
        [],
        [],
        0,
        0,
        1,
        'err',
        '0x'
      );
  }

  describe('basics', function () {
    it('appLockedFunds is 0 after bootstrapApplication', async () => {
      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(0n);
    });

    it('increases by depositAmount + maxFeeValue after submitRequest', async () => {
      const depositAmount = 50n;
      const maxFeeValue = minFeePerRequest + 10n;

      await submitRequest(signers[0], '0x01', depositAmount, maxFeeValue);

      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(
        depositAmount + maxFeeValue
      );
    });

    it('increases by msg.value after submitDeployRequest', async () => {
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

      expect(await processorEndpoint.appLockedFunds(deployAppId)).to.equal(fee);
    });

    it('decreases correctly after successful stateUpdate', async () => {
      const depositAmount = 100n;
      const refund = 10n;
      const applicationFees = minFeePerRequest;
      const maxFeeValue = refund + applicationFees;

      const request = await submitRequest(signers[0], '0x02', depositAmount, maxFeeValue);
      const fundsAfterSubmit = await processorEndpoint.appLockedFunds(applicationId);
      expect(fundsAfterSubmit).to.equal(depositAmount + maxFeeValue);

      const withdrawalAddr = await signers[3].getAddress();
      const withdrawalAmount = 50n;

      await completeRequest(
        request.requestId,
        INITIAL_STATE_ROOT,
        '0x' + '22'.repeat(32),
        refund,
        applicationFees,
        [[withdrawalAddr, withdrawalAmount]]
      );

      const totalDebited = withdrawalAmount + refund + applicationFees;
      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(
        fundsAfterSubmit - totalDebited
      );
    });

    it('decreases correctly after error stateUpdate', async () => {
      const depositAmount = 30n;
      const maxFeeValue = minFeePerRequest + 5n;

      const request = await submitRequest(signers[0], '0x03', depositAmount, maxFeeValue);
      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(
        depositAmount + maxFeeValue
      );

      await failRequest(request.requestId);

      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(0n);
    });

    it('returns 0 for unknown application IDs', async () => {
      expect(await processorEndpoint.appLockedFunds(999999)).to.equal(0n);
    });
  });

  describe('cross-app isolation', function () {
    let appIdB: bigint;

    beforeEach(async function () {
      ({ applicationId: appIdB } = await bootstrapApplication(processorEndpoint));
    });

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
        [[await signers[3].getAddress(), depositA]]
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
        [[await signers[3].getAddress(), depositA]]
      );

      // B tries to withdraw 200 — should fail because B only has 50 + maxFee
      await expect(
        completeRequest(
          reqB.requestId,
          INITIAL_STATE_ROOT,
          '0x' + 'b1'.repeat(32),
          0n,
          maxFee,
          [[await signers[4].getAddress(), 200n]],
          appIdB
        )
      ).to.be.revertedWithCustomError(processorEndpoint, 'InsufficientAppBalance');
    });

    it("App A's stateUpdate does not change App B's appLockedFunds", async () => {
      const depositA = 100n;
      const depositB = 75n;
      const maxFee = minFeePerRequest;

      const reqA = await submitRequest(signers[0], '0x20', depositA, maxFee);
      await submitRequest(signers[0], '0x21', depositB, maxFee, appIdB);

      const fundsB = await processorEndpoint.appLockedFunds(appIdB);
      expect(fundsB).to.equal(depositB + maxFee);

      // Process A
      await completeRequest(
        reqA.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'a3'.repeat(32),
        0n,
        maxFee,
        [[await signers[3].getAddress(), depositA]]
      );

      // B's funds unchanged
      expect(await processorEndpoint.appLockedFunds(appIdB)).to.equal(fundsB);
    });

    it("App A error refund does not affect App B's pool", async () => {
      const depositA = 40n;
      const depositB = 60n;
      const maxFee = minFeePerRequest;

      const reqA = await submitRequest(signers[0], '0x30', depositA, maxFee);
      await submitRequest(signers[0], '0x31', depositB, maxFee, appIdB);

      const fundsB = await processorEndpoint.appLockedFunds(appIdB);

      // Fail A
      await failRequest(reqA.requestId);

      // B unchanged
      expect(await processorEndpoint.appLockedFunds(appIdB)).to.equal(fundsB);
    });
  });

  describe('accumulation', function () {
    it('multiple pending requests accumulate, process sequentially to zero', async () => {
      const dep1 = 100n;
      const dep2 = 200n;
      const maxFee = minFeePerRequest;

      const req1 = await submitRequest(signers[0], '0x40', dep1, maxFee);
      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(dep1 + maxFee);

      const req2 = await submitRequest(signers[0], '0x41', dep2, maxFee);
      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(
        dep1 + dep2 + 2n * maxFee
      );

      // Process first
      await completeRequest(
        req1.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'c1'.repeat(32),
        0n,
        maxFee,
        [[await signers[3].getAddress(), dep1]]
      );

      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(dep2 + maxFee);

      // Process second
      await completeRequest(
        req2.requestId,
        '0x' + 'c1'.repeat(32),
        '0x' + 'c2'.repeat(32),
        0n,
        maxFee,
        [[await signers[3].getAddress(), dep2]]
      );

      expect(await processorEndpoint.appLockedFunds(applicationId)).to.equal(0n);
    });
  });

  describe('deploy request tracking', function () {
    it('deploy request credits then stateUpdate debits to zero', async () => {
      // bootstrapApplication already consumed the first deploy, so appLockedFunds is 0
      // Deploy a new app and verify the cycle
      const fixture = await deployProcessorEndpointFixture();
      const pe = await fixture.deployProcessorEndpoint();

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

      expect(await pe.appLockedFunds(deployAppId)).to.equal(fixture.minFeePerRequest);

      // Complete deploy
      await pe
        .connect(fixture.signers[1])
        .stateUpdate(
          deployAppId,
          BYTES32_ZERO,
          INITIAL_STATE_ROOT,
          deployRequestId,
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest,
          0,
          '',
          '0x'
        );

      expect(await pe.appLockedFunds(deployAppId)).to.equal(0n);
    });
  });
});
