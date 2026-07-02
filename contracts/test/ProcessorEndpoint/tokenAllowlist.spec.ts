import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { ETH_TOKEN, getRequestIdFromReceipt } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let tokenAllowlist: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let mockERC20: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint, tokenAllowlist } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;

    // Deploy a mock ERC20 to use as a valid contract address
    const MockERC20 = await ethers.getContractFactory('MockERC20');
    mockERC20 = await MockERC20.deploy('Mock Token', 'MCK', 18);
  });

  describe('addAllowedToken', function () {
    describe('unhappy paths', function () {
      it('reverts when called by non-admin', async () => {
        await expect(
          tokenAllowlist.connect(signers[0]).addAllowedToken(await mockERC20.getAddress())
        ).to.be.revertedWithCustomError(tokenAllowlist, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with AddressCantBeZero when token is address(0)', async () => {
        await expect(
          tokenAllowlist.connect(signers[2]).addAllowedToken(ETH_TOKEN)
        ).to.be.revertedWithCustomError(tokenAllowlist, 'TokenAddressCantBeZero');
      });

      it('reverts with NotAContract when token address has no code', async () => {
        const eoa = await signers[5].getAddress();
        await expect(
          tokenAllowlist.connect(signers[2]).addAllowedToken(eoa)
        ).to.be.revertedWithCustomError(tokenAllowlist, 'NotAContract');
      });
    });

    describe('happy paths', function () {
      it('adds a token to the allowlist and emits TokenAllowed', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const tx = await tokenAllowlist.connect(signers[2]).addAllowedToken(tokenAddr);
        await expect(tx).to.emit(tokenAllowlist, 'TokenAllowed').withArgs(tokenAddr);
        expect(await tokenAllowlist.isAllowedToken(tokenAddr)).to.be.true;
      });

      it('is idempotent when adding an already-allowed token', async () => {
        const tokenAddr = await mockERC20.getAddress();
        await tokenAllowlist.connect(signers[2]).addAllowedToken(tokenAddr);
        const tx = await tokenAllowlist.connect(signers[2]).addAllowedToken(tokenAddr);
        await expect(tx).to.emit(tokenAllowlist, 'TokenAllowed').withArgs(tokenAddr);
        expect(await tokenAllowlist.isAllowedToken(tokenAddr)).to.be.true;
      });
    });
  });

  describe('removeAllowedToken', function () {
    describe('unhappy paths', function () {
      it('reverts when called by non-admin', async () => {
        await expect(
          tokenAllowlist.connect(signers[0]).removeAllowedToken(await mockERC20.getAddress())
        ).to.be.revertedWithCustomError(tokenAllowlist, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with AddressCantBeZero when token is address(0)', async () => {
        await expect(
          tokenAllowlist.connect(signers[2]).removeAllowedToken(ETH_TOKEN)
        ).to.be.revertedWithCustomError(tokenAllowlist, 'TokenAddressCantBeZero');
      });
    });

    describe('happy paths', function () {
      it('removes a token from the allowlist and emits TokenRemoved', async () => {
        const tokenAddr = await mockERC20.getAddress();
        await tokenAllowlist.connect(signers[2]).addAllowedToken(tokenAddr);
        const tx = await tokenAllowlist.connect(signers[2]).removeAllowedToken(tokenAddr);
        await expect(tx).to.emit(tokenAllowlist, 'TokenRemoved').withArgs(tokenAddr);
        expect(await tokenAllowlist.isAllowedToken(tokenAddr)).to.be.false;
      });

      it('is idempotent when removing a non-allowed token', async () => {
        const tokenAddr = await mockERC20.getAddress();
        const tx = await tokenAllowlist.connect(signers[2]).removeAllowedToken(tokenAddr);
        await expect(tx).to.emit(tokenAllowlist, 'TokenRemoved').withArgs(tokenAddr);
        expect(await tokenAllowlist.isAllowedToken(tokenAddr)).to.be.false;
      });
    });
  });

  describe('isAllowedToken', function () {
    it('returns true for ETH (address(0))', async () => {
      expect(await tokenAllowlist.isAllowedToken(ETH_TOKEN)).to.be.true;
    });

    it('returns false for a non-allowlisted token', async () => {
      expect(await tokenAllowlist.isAllowedToken(await mockERC20.getAddress())).to.be.false;
    });

    it('returns true for an allowlisted token', async () => {
      const tokenAddr = await mockERC20.getAddress();
      await tokenAllowlist.connect(signers[2]).addAllowedToken(tokenAddr);
      expect(await tokenAllowlist.isAllowedToken(tokenAddr)).to.be.true;
    });
  });

  describe('delisted token withdrawal and claim', function () {
    let applicationId: bigint;

    beforeEach(async function () {
      const fixture = await deployProcessorEndpointFixture();
      ({ processorEndpoint, tokenAllowlist } = await fixture.deployProcessorEndpoint());
      signers = fixture.signers;
      minFeePerRequest = fixture.minFeePerRequest;
      ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));

      const MockERC20 = await ethers.getContractFactory('MockERC20');
      mockERC20 = await MockERC20.deploy('Mock Token', 'MCK', 18);
      await tokenAllowlist.connect(signers[2]).addAllowedToken(await mockERC20.getAddress());
    });

    it('blocks new deposits after token is removed from allowlist', async () => {
      const tokenAddr = await mockERC20.getAddress();
      const assetAmount = 100n;

      await mockERC20.mint(await signers[0].getAddress(), assetAmount);
      await mockERC20.approve(await processorEndpoint.getAddress(), assetAmount);

      // Remove from allowlist
      await tokenAllowlist.connect(signers[2]).removeAllowedToken(tokenAddr);

      // New deposit should be rejected
      await expect(
        processorEndpoint.submitRequest(
          0,
          applicationId,
          1,
          '0x01',
          tokenAddr,
          assetAmount,
          minFeePerRequest,
          { value: minFeePerRequest }
        )
      ).to.be.revertedWithCustomError(processorEndpoint, 'TokenNotAllowed');
    });

    it('allows withdrawal and claim of a delisted token', async () => {
      const tokenAddr = await mockERC20.getAddress();
      const assetAmount = 100n;
      const payee = await signers[4].getAddress();

      // Deposit while token is allowlisted
      await mockERC20.mint(await signers[0].getAddress(), assetAmount);
      await mockERC20.approve(await processorEndpoint.getAddress(), assetAmount);

      const tx = await processorEndpoint.submitRequest(
        0,
        applicationId,
        1,
        '0x01',
        tokenAddr,
        assetAmount,
        minFeePerRequest,
        { value: minFeePerRequest }
      );
      const receipt = await tx.wait();
      const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

      // Remove token from allowlist
      await tokenAllowlist.connect(signers[2]).removeAllowedToken(tokenAddr);
      expect(await tokenAllowlist.isAllowedToken(tokenAddr)).to.be.false;

      // stateUpdate with ERC-20 withdrawal should still succeed
      await processorEndpoint
        .connect(signers[1])
        .stateUpdate(
          applicationId,
          INITIAL_STATE_ROOT,
          '0x' + 'dd'.repeat(32),
          requestId,
          { events: [], subTypes: [] },
          { events: [], subTypes: [] },
          [[tokenAddr, payee, assetAmount]],
          0,
          minFeePerRequest,
          0,
          '',
          '0x'
        );

      expect(await processorEndpoint.pendingClaims(tokenAddr, payee)).to.equal(assetAmount);

      // Claim should still succeed
      const balBefore = await mockERC20.balanceOf(payee);
      await processorEndpoint.claim(tokenAddr, payee);
      const balAfter = await mockERC20.balanceOf(payee);
      expect(balAfter - balBefore).to.equal(assetAmount);
    });
  });
});
