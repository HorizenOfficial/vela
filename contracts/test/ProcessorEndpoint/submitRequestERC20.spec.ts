import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture } from './fixture';
import { ETH_TOKEN } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;
  let mockERC20: any;

  const REQUEST_TYPE_PROCESS = 1;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));

    // Deploy and allowlist a mock ERC20
    const MockERC20 = await ethers.getContractFactory('MockERC20');
    mockERC20 = await MockERC20.deploy('Mock Token', 'MCK', 18);
    await processorEndpoint.connect(signers[2]).addAllowedToken(await mockERC20.getAddress());
  });

  describe('submitRequest with ERC-20', function () {
    describe('unhappy paths', function () {
      it('reverts with InvalidValue when assetAmount is 0 but tokenAddress is non-zero', async () => {
        const tokenAddr = await mockERC20.getAddress();
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            tokenAddr,
            0,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with TokenNotAllowed when tokenAddress is not allowlisted', async () => {
        const MockERC20 = await ethers.getContractFactory('MockERC20');
        const unlisted = await MockERC20.deploy('Unlisted', 'UNL', 18);
        const unlistedAddr = await unlisted.getAddress();
        const assetAmount = 100n;

        await unlisted.mint(await signers[0].getAddress(), assetAmount);
        await unlisted.approve(await processorEndpoint.getAddress(), assetAmount);

        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            unlistedAddr,
            assetAmount,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'TokenNotAllowed');
      });

      it('reverts with InvalidValue when msg.value != maxFeeValue for ERC-20 request', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 100n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.approve(await processorEndpoint.getAddress(), assetAmount);

        // Overpay: msg.value = maxFeeValue + assetAmount (wrong for ERC-20)
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            tokenAddr,
            assetAmount,
            minFeePerRequest,
            { value: minFeePerRequest + assetAmount }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts when user has not approved the ERC-20 transfer', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 100n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        // No approve call

        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            tokenAddr,
            assetAmount,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.reverted;
      });

      it('reverts when approved amount is less than assetAmount', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 100n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.approve(await processorEndpoint.getAddress(), assetAmount - 1n);

        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            tokenAddr,
            assetAmount,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.reverted;
      });

      it('reverts with TransferAmountMismatch for fee-on-transfer tokens', async () => {
        const FeeOnTransfer = await ethers.getContractFactory('FeeOnTransferERC20');
        const feeToken = await FeeOnTransfer.deploy();
        const feeTokenAddr = await feeToken.getAddress();
        const assetAmount = 100n;

        // Allowlist the fee-on-transfer token
        await processorEndpoint.connect(signers[2]).addAllowedToken(feeTokenAddr);
        await feeToken.mint(await signers[0].getAddress(), assetAmount);
        await feeToken.approve(await processorEndpoint.getAddress(), assetAmount);

        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            feeTokenAddr,
            assetAmount,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'TransferAmountMismatch');
      });
    });

    describe('happy paths', function () {
      it('pulls ERC-20 tokens and stores request correctly', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 500n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.approve(await processorEndpoint.getAddress(), assetAmount);

        const tx = await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          tokenAddr,
          assetAmount,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        const requestSubmittedLog = receipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'RequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(requestSubmittedLog);
        const requestId = parsed.args.requestId;
        const stored = await processorEndpoint.requestById(requestId);

        expect(stored.tokenAddress).to.equal(tokenAddr);
        expect(stored.assetAmount).to.equal(assetAmount);
        expect(stored.maxFeeValue).to.equal(minFeePerRequest);
      });

      it('transfers ERC-20 to contract and ETH for fee only', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 200n;
        const processorAddr = await processorEndpoint.getAddress();

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.approve(processorAddr, assetAmount);

        const ethBalanceBefore = await signers[0].provider!.getBalance(processorAddr);
        const tokenBalanceBefore = await mockERC20.balanceOf(processorAddr);

        await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          tokenAddr,
          assetAmount,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        const ethBalanceAfter = await signers[0].provider!.getBalance(processorAddr);
        const tokenBalanceAfter = await mockERC20.balanceOf(processorAddr);

        // ETH balance increases by maxFeeValue only (no assetAmount in ETH)
        expect(ethBalanceAfter).to.equal(ethBalanceBefore + minFeePerRequest);
        // ERC-20 balance increases by assetAmount
        expect(tokenBalanceAfter).to.equal(tokenBalanceBefore + assetAmount);
      });

      it('updates appCustody for ERC-20 deposit', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 300n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.approve(await processorEndpoint.getAddress(), assetAmount);

        await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          tokenAddr,
          assetAmount,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(assetAmount);
      });
    });
  });
});
