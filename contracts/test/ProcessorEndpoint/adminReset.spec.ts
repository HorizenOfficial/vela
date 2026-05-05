import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { ADDRESS_ZERO, ETH_TOKEN, BYTES32_ZERO } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let bootstrapApplication: any;
  let resetOperator: string;
  let deployProcessorEndpoint: (resetOperatorOverride?: string) => Promise<any>;
  let mockERC20: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
    resetOperator = fixture.resetOperator;
    deployProcessorEndpoint = fixture.deployProcessorEndpoint;

    const MockERC20 = await ethers.getContractFactory('MockERC20');
    mockERC20 = await MockERC20.deploy('Mock Token', 'MCK', 18);
    await processorEndpoint.connect(signers[2]).addAllowedToken(await mockERC20.getAddress());
  });

  describe('getAllowedTokens', function () {
    it('returns empty array when no tokens are allowlisted', async () => {
      const fresh = await deployProcessorEndpoint();
      expect(await fresh.getAllowedTokens()).to.deep.equal([]);
    });

    it('returns added tokens', async () => {
      const tokenAddr = await mockERC20.getAddress();
      expect(await processorEndpoint.getAllowedTokens()).to.deep.equal([tokenAddr]);
    });

    it('excludes removed tokens', async () => {
      const tokenAddr = await mockERC20.getAddress();
      await processorEndpoint.connect(signers[2]).removeAllowedToken(tokenAddr);
      expect(await processorEndpoint.getAllowedTokens()).to.deep.equal([]);
    });

    it('does not duplicate a token added twice', async () => {
      const tokenAddr = await mockERC20.getAddress();
      await processorEndpoint.connect(signers[2]).addAllowedToken(tokenAddr);
      const tokens = await processorEndpoint.getAllowedTokens();
      expect(tokens.length).to.equal(1);
      expect(tokens[0]).to.equal(tokenAddr);
    });
  });

  describe('getDeployedAppIds', function () {
    it('returns empty array before any deploy', async () => {
      expect(await processorEndpoint.getDeployedAppIds()).to.deep.equal([]);
    });

    it('records an app ID after a successful deploy', async () => {
      const { applicationId } = await bootstrapApplication(processorEndpoint);
      const ids = await processorEndpoint.getDeployedAppIds();
      expect(ids.length).to.equal(1);
      expect(ids[0]).to.equal(applicationId);
    });

    it('records multiple app IDs after multiple deploys', async () => {
      const { applicationId: id1 } = await bootstrapApplication(processorEndpoint);
      const { applicationId: id2 } = await bootstrapApplication(processorEndpoint);
      const ids = await processorEndpoint.getDeployedAppIds();
      expect(ids.length).to.equal(2);
      expect(ids[0]).to.equal(id1);
      expect(ids[1]).to.equal(id2);
    });
  });

  describe('adminReset', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks RESET_OPERATOR role', async () => {
        await expect(
          processorEndpoint.connect(signers[0]).adminReset()
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('is unreachable when deployed with address(0) as reset operator', async () => {
        const noResetEndpoint = await deployProcessorEndpoint(ADDRESS_ZERO);
        await expect(
          noResetEndpoint.connect(signers[3]).adminReset()
        ).to.be.revertedWithCustomError(noResetEndpoint, 'AccessControlUnauthorizedAccount');
      });
    });

    describe('happy paths', function () {
      it('clears the pending request queue', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        await processorEndpoint.submitRequest(0, applicationId, 1, '0x01', ETH_TOKEN, 0, minFeePerRequest, {
          value: minFeePerRequest,
        });
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);

        await processorEndpoint.connect(signers[3]).adminReset();
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
      });

      it('frees deploy slots for pending DEPLOYAPP requests', async () => {
        // First bootstrap an app (consumes a deploy slot)
        await bootstrapApplication(processorEndpoint);
        const slotsBefore = await processorEndpoint.availableDeploySlots();

        // Submit a deploy request (pending, not yet finalised)
        await processorEndpoint.connect(signers[2]).submitDeployRequest(0, '0xAA', {
          value: minFeePerRequest,
        });
        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore - 1n);

        await processorEndpoint.connect(signers[3]).adminReset();
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore);
      });

      it('does not free slots for already-deployed apps', async () => {
        await bootstrapApplication(processorEndpoint);
        const slotsAfterDeploy = await processorEndpoint.availableDeploySlots();

        await processorEndpoint.connect(signers[3]).adminReset();
        // Queue was already empty; slots unchanged
        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsAfterDeploy);
      });
    });
  });

  describe('adminResetApps', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks RESET_OPERATOR role', async () => {
        await expect(
          processorEndpoint.connect(signers[0]).adminResetApps([], [])
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });
    });

    describe('happy paths', function () {
      it('clears state roots for all deployed apps when called with empty arrays', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        expect(await processorEndpoint.applicationStateRoots(applicationId)).to.not.equal(BYTES32_ZERO);

        await processorEndpoint.connect(signers[3]).adminResetApps([], []);
        expect(await processorEndpoint.applicationStateRoots(applicationId)).to.equal(BYTES32_ZERO);
      });

      it('clears state root only for the specified app when explicit appIds given', async () => {
        const { applicationId: id1 } = await bootstrapApplication(processorEndpoint);
        const { applicationId: id2 } = await bootstrapApplication(processorEndpoint);

        await processorEndpoint.connect(signers[3]).adminResetApps([id1], []);
        expect(await processorEndpoint.applicationStateRoots(id1)).to.equal(BYTES32_ZERO);
        expect(await processorEndpoint.applicationStateRoots(id2)).to.not.equal(BYTES32_ZERO);
      });

      it('frees deploy slots for cleared apps', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const slotsBefore = await processorEndpoint.availableDeploySlots();

        await processorEndpoint.connect(signers[3]).adminResetApps([applicationId], []);
        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore + 1n);
      });

      it('also clears the pending request queue', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        await processorEndpoint.submitRequest(0, applicationId, 1, '0x01', ETH_TOKEN, 0, minFeePerRequest, {
          value: minFeePerRequest,
        });
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);

        await processorEndpoint.connect(signers[3]).adminResetApps([], []);
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
      });

      it('recovers ETH custody and transfers it to the caller', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        // Deposit ETH as a business asset (assetAmount > 0, ETH token)
        const assetAmount = ethers.parseEther('0.5');
        await processorEndpoint.submitRequest(
          0, applicationId, 1, '0x01', ETH_TOKEN, assetAmount, minFeePerRequest,
          { value: assetAmount + minFeePerRequest }
        );
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(assetAmount);

        const resetOperatorSigner = signers[3];
        const resetOperatorAddr = await resetOperatorSigner.getAddress();
        const balBefore = await resetOperatorSigner.provider!.getBalance(resetOperatorAddr);

        const tx = await processorEndpoint.connect(resetOperatorSigner).adminResetApps([], []);
        const receipt = await tx.wait();
        const gasCost = receipt.gasUsed * receipt.gasPrice;
        const balAfter = await resetOperatorSigner.provider!.getBalance(resetOperatorAddr);

        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
        expect(await processorEndpoint.totalAppCustody(ETH_TOKEN)).to.equal(0n);
        expect(balAfter).to.equal(balBefore + assetAmount - gasCost);
      });

      it('recovers ERC-20 custody and transfers it to the caller', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 500n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.connect(signers[0]).approve(await processorEndpoint.getAddress(), assetAmount);
        await processorEndpoint.submitRequest(
          0, applicationId, 1, '0x01', tokenAddr, assetAmount, minFeePerRequest,
          { value: minFeePerRequest }
        );

        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(assetAmount);
        expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(assetAmount);

        const resetOperatorAddr = await signers[3].getAddress();
        const balBefore = await mockERC20.balanceOf(resetOperatorAddr);

        await processorEndpoint.connect(signers[3]).adminResetApps([], []);

        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(0n);
        expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(0n);
        expect(await mockERC20.balanceOf(resetOperatorAddr)).to.equal(balBefore + assetAmount);
      });

      it('uses getAllowedTokens when erc20Tokens is empty', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 200n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20.connect(signers[0]).approve(await processorEndpoint.getAddress(), assetAmount);
        await processorEndpoint.submitRequest(
          0, applicationId, 1, '0x01', tokenAddr, assetAmount, minFeePerRequest,
          { value: minFeePerRequest }
        );

        // Passing empty arrays — should pick up the allowlisted token automatically
        await processorEndpoint.connect(signers[3]).adminResetApps([], []);
        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(0n);
      });
    });
  });
});
