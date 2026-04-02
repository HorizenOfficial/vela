import { expect } from 'chai';
import { BigNumberish, Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { ETH_TOKEN, getRequestIdFromReceipt } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let protocolVersion: BigNumberish;
  let applicationId: bigint;
  let mockERC20: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    protocolVersion = await processorEndpoint.PROTOCOL_VERSION();
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));

    const MockERC20 = await ethers.getContractFactory('MockERC20');
    mockERC20 = await MockERC20.deploy('Mock Token', 'MCK', 18);
    await processorEndpoint.connect(signers[2]).addAllowedToken(await mockERC20.getAddress());
  });

  describe('claim', function () {
    describe('ETH claims', function () {
      it('does not revert when pending balance is zero', async () => {
        await processorEndpoint.claim(ETH_TOKEN, await signers[5].getAddress());
      });

      it('transfers ETH and emits PaymentWithdrawn', async () => {
        // Submit and fail a request to create a pending ETH claim
        const depositAmount = 100n;
        const tx = await processorEndpoint
          .connect(signers[4])
          .submitRequest(0, applicationId, 1, '0x01', ETH_TOKEN, depositAmount, minFeePerRequest, {
            value: depositAmount + minFeePerRequest,
          });
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
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

        const payee = await signers[4].getAddress();
        const pending = await processorEndpoint.pendingClaims(ETH_TOKEN, payee);
        expect(pending).to.equal(depositAmount);

        const balanceBefore = await ethers.provider.getBalance(payee);
        const claimTx = await processorEndpoint.connect(signers[3]).claim(ETH_TOKEN, payee);
        await expect(claimTx)
          .to.emit(processorEndpoint, 'PaymentWithdrawn')
          .withArgs(ETH_TOKEN, payee, depositAmount);
        const balanceAfter = await ethers.provider.getBalance(payee);

        expect(balanceAfter).to.equal(balanceBefore + depositAmount);
        expect(await processorEndpoint.pendingClaims(ETH_TOKEN, payee)).to.equal(0n);
      });

      it('is permissionless — any caller can claim for any payee', async () => {
        const depositAmount = 50n;
        const tx = await processorEndpoint
          .connect(signers[4])
          .submitRequest(0, applicationId, 1, '0x01', ETH_TOKEN, depositAmount, minFeePerRequest, {
            value: depositAmount + minFeePerRequest,
          });
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
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

        const payee = await signers[4].getAddress();
        // signers[0] triggers claim for signers[4]
        await processorEndpoint.connect(signers[0]).claim(ETH_TOKEN, payee);
        expect(await processorEndpoint.pendingClaims(ETH_TOKEN, payee)).to.equal(0n);
      });
    });

    describe('ERC-20 claims', function () {
      it('does not revert when pending ERC-20 balance is zero', async () => {
        await processorEndpoint.claim(await mockERC20.getAddress(), await signers[5].getAddress());
      });

      it('transfers ERC-20 tokens after successful withdrawal', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 200n;
        const withdrawalAmount = 150n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.connect(signers[0]).approve(await processorEndpoint.getAddress(), assetAmount);

        const submitTx = await processorEndpoint.submitRequest(
          0,
          applicationId,
          1,
          '0x01',
          tokenAddr,
          assetAmount,
          minFeePerRequest,
          { value: minFeePerRequest }
        );
        const submitReceipt = await submitTx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, submitReceipt);

        const payee = await signers[5].getAddress();
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            '0x' + 'cc'.repeat(32),
            requestId,
            [],
            [],
            [[tokenAddr, payee, withdrawalAmount]],
            0,
            minFeePerRequest,
            0,
            '',
            '0x'
          );

        expect(await processorEndpoint.pendingClaims(tokenAddr, payee)).to.equal(withdrawalAmount);

        const tokenBalanceBefore = await mockERC20.balanceOf(payee);
        const claimTx = await processorEndpoint.claim(tokenAddr, payee);
        await expect(claimTx)
          .to.emit(processorEndpoint, 'PaymentWithdrawn')
          .withArgs(tokenAddr, payee, withdrawalAmount);
        const tokenBalanceAfter = await mockERC20.balanceOf(payee);

        expect(tokenBalanceAfter).to.equal(tokenBalanceBefore + withdrawalAmount);
        expect(await processorEndpoint.pendingClaims(tokenAddr, payee)).to.equal(0n);
      });

      it('credits ERC-20 pending claim on failed ERC-20 request', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 100n;

        await mockERC20.mint(await signers[4].getAddress(), assetAmount);
        await mockERC20
          .connect(signers[4])
          .approve(await processorEndpoint.getAddress(), assetAmount);

        const submitTx = await processorEndpoint
          .connect(signers[4])
          .submitRequest(0, applicationId, 1, '0x01', tokenAddr, assetAmount, minFeePerRequest, {
            value: minFeePerRequest,
          });
        const submitReceipt = await submitTx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, submitReceipt);

        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
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

        const payee = await signers[4].getAddress();
        // ERC-20 asset refund in token
        expect(await processorEndpoint.pendingClaims(tokenAddr, payee)).to.equal(assetAmount);
        // ETH fee refund
        const ethRefund = minFeePerRequest - (await processorEndpoint.minFeePerRequest());
        expect(await processorEndpoint.pendingClaims(ETH_TOKEN, payee)).to.equal(ethRefund);
      });
    });
  });
});
