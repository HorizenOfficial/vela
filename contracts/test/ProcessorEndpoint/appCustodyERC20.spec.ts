import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  ETH_TOKEN,
  getRequestIdFromReceipt,
  PROTOCOL_VERSION,
  REQUEST_TYPE_PROCESS,
} from '../util';

describe('ProcessorEndpoint — appCustody ERC-20', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;
  let bootstrapApplication: any;
  let mockERC20: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));

    const MockERC20 = await ethers.getContractFactory('MockERC20');
    mockERC20 = await MockERC20.deploy('Mock Token', 'MCK', 18);
    await processorEndpoint.connect(signers[2]).addAllowedToken(await mockERC20.getAddress());
  });

  async function submitERC20Request(
    sender: Signer,
    payload: string,
    assetAmount: bigint,
    maxFeeValue: bigint,
    appId?: bigint
  ) {
    const tokenAddr = await mockERC20.getAddress();
    const senderAddr = await sender.getAddress();
    await mockERC20.mint(senderAddr, assetAmount);
    await mockERC20.connect(sender).approve(await processorEndpoint.getAddress(), assetAmount);

    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        appId ?? applicationId,
        REQUEST_TYPE_PROCESS,
        payload,
        tokenAddr,
        assetAmount,
        maxFeeValue,
        { value: maxFeeValue }
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
    const currentStateRoot = await processorEndpoint.applicationStateRoots(appId ?? applicationId);
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
    it('appCustody is 0 for ERC-20 token after bootstrapApplication', async () => {
      const tokenAddr = await mockERC20.getAddress();
      expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(0n);
    });

    it('increases appCustody and totalAppCustody after ERC-20 submitRequest', async () => {
      const tokenAddr = await mockERC20.getAddress();
      const assetAmount = 500n;

      const appCustodyBefore = await processorEndpoint.appCustody(applicationId, tokenAddr);
      const totalAppCustodyBefore = await processorEndpoint.totalAppCustody(tokenAddr);

      await submitERC20Request(signers[0], '0x01', assetAmount, minFeePerRequest);

      expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(
        appCustodyBefore + assetAmount
      );
      expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(
        totalAppCustodyBefore + assetAmount
      );
    });

    it('does not affect ETH custody when depositing ERC-20', async () => {
      const assetAmount = 500n;
      const appCustodyBefore = await processorEndpoint.appCustody(applicationId, ETH_TOKEN);
      const totalAppCustodyBefore = await processorEndpoint.totalAppCustody(ETH_TOKEN);

      await submitERC20Request(signers[0], '0x01', assetAmount, minFeePerRequest);

      expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(
        appCustodyBefore
      );
      expect(await processorEndpoint.totalAppCustody(ETH_TOKEN)).to.equal(totalAppCustodyBefore);
    });

    it('decreases appCustody after successful stateUpdate with ERC-20 withdrawal', async () => {
      const tokenAddr = await mockERC20.getAddress();
      const assetAmount = 100n;
      const maxFee = minFeePerRequest;

      const req = await submitERC20Request(signers[0], '0x02', assetAmount, maxFee);
      expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(assetAmount);

      const withdrawalAmount = 60n;
      await completeRequest(req.requestId, INITIAL_STATE_ROOT, '0x' + 'ee'.repeat(32), 0n, maxFee, [
        [tokenAddr, await signers[3].getAddress(), withdrawalAmount],
      ]);

      expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(
        assetAmount - withdrawalAmount
      );
      expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(
        assetAmount - withdrawalAmount
      );
    });

    it('decreases appCustody after error stateUpdate for ERC-20 request', async () => {
      const tokenAddr = await mockERC20.getAddress();
      const assetAmount = 80n;

      const req = await submitERC20Request(signers[0], '0x03', assetAmount, minFeePerRequest);
      expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(assetAmount);

      await failRequest(req.requestId);

      expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(0n);
      expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(0n);
    });
  });

  describe('cross-app isolation', function () {
    let appIdB: bigint;

    beforeEach(async function () {
      ({ applicationId: appIdB } = await bootstrapApplication(processorEndpoint));
    });

    it('App B cannot withdraw ERC-20 exceeding its own custody', async () => {
      const tokenAddr = await mockERC20.getAddress();
      const depositA = 200n;
      const depositB = 50n;
      const maxFee = minFeePerRequest;

      // Submit A, then complete A to free queue head
      const reqA = await submitERC20Request(signers[0], '0x10', depositA, maxFee);
      await completeRequest(
        reqA.requestId,
        INITIAL_STATE_ROOT,
        '0x' + 'a1'.repeat(32),
        0n,
        maxFee,
        [] // no withdrawals for A
      );

      // Submit B with a smaller deposit
      const reqB = await submitERC20Request(signers[0], '0x12', depositB, maxFee, appIdB);

      // B tries to withdraw more than its custody
      await expect(
        completeRequest(
          reqB.requestId,
          INITIAL_STATE_ROOT,
          '0x' + 'b1'.repeat(32),
          0n,
          maxFee,
          [[tokenAddr, await signers[4].getAddress(), depositB + 1n]],
          appIdB
        )
      ).to.be.revertedWithCustomError(processorEndpoint, 'InsufficientAppBalance');
    });
  });
});
