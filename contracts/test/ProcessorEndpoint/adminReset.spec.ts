import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  ADDRESS_ZERO,
  ETH_TOKEN,
  BYTES32_ZERO,
  getRequestIdFromReceipt,
  PROTOCOL_VERSION,
  REQUEST_TYPE_PROCESS,
} from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let tokenAllowlist: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let bootstrapApplication: any;
  let resetOperator: string;
  let deployProcessorEndpoint: (resetOperatorOverride?: string) => Promise<any>;
  let mockERC20: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint, tokenAllowlist } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
    resetOperator = fixture.resetOperator;
    deployProcessorEndpoint = fixture.deployProcessorEndpoint;

    const MockERC20 = await ethers.getContractFactory('MockERC20');
    mockERC20 = await MockERC20.deploy('Mock Token', 'MCK', 18);
    await tokenAllowlist.connect(signers[2]).addAllowedToken(await mockERC20.getAddress());
  });

  async function submitProcess(applicationId: bigint, payload: string) {
    const tx = await processorEndpoint.submitRequest(
      PROTOCOL_VERSION,
      applicationId,
      REQUEST_TYPE_PROCESS,
      payload,
      ETH_TOKEN,
      0,
      minFeePerRequest,
      { value: minFeePerRequest }
    );
    return getRequestIdFromReceipt(processorEndpoint, await tx.wait());
  }

  // Completes a request through the single-request path, which advances the round-robin cursor.
  async function processRequest(
    applicationId: bigint,
    requestId: string,
    prev: string,
    next: string
  ) {
    await processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        applicationId,
        prev,
        next,
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
  }

  describe('getAllowedTokens', function () {
    it('returns empty array when no tokens are allowlisted', async () => {
      const { tokenAllowlist: freshTokenAllowlist } = await deployProcessorEndpoint();
      expect(await freshTokenAllowlist.getAllowedTokens()).to.deep.equal([]);
    });

    it('returns added tokens', async () => {
      const tokenAddr = await mockERC20.getAddress();
      expect(await tokenAllowlist.getAllowedTokens()).to.deep.equal([tokenAddr]);
    });

    it('excludes removed tokens', async () => {
      const tokenAddr = await mockERC20.getAddress();
      await tokenAllowlist.connect(signers[2]).removeAllowedToken(tokenAddr);
      expect(await tokenAllowlist.getAllowedTokens()).to.deep.equal([]);
    });

    it('does not duplicate a token added twice', async () => {
      const tokenAddr = await mockERC20.getAddress();
      await tokenAllowlist.connect(signers[2]).addAllowedToken(tokenAddr);
      const tokens = await tokenAllowlist.getAllowedTokens();
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
        const { processorEndpoint: noResetEndpoint } = await deployProcessorEndpoint(ADDRESS_ZERO);
        await expect(
          noResetEndpoint.connect(signers[3]).adminReset()
        ).to.be.revertedWithCustomError(noResetEndpoint, 'AccessControlUnauthorizedAccount');
      });
    });

    describe('happy paths', function () {
      it('clears the pending request queue', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        await processorEndpoint.submitRequest(
          0,
          applicationId,
          1,
          '0x01',
          ETH_TOKEN,
          0,
          minFeePerRequest,
          {
            value: minFeePerRequest,
          }
        );
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

      it('resets the round-robin cursor', async () => {
        // The cursor is an index into the deployed-app list, so it needs more than one entry
        // to be observable.
        const { applicationId: appA } = await bootstrapApplication(processorEndpoint);
        const { applicationId: appB } = await bootstrapApplication(processorEndpoint);
        await bootstrapApplication(processorEndpoint);

        // Give A its turn, which moves the cursor off index 0 and onto B.
        const a1 = await submitProcess(appA, '0x01');
        await processRequest(appA, a1, INITIAL_STATE_ROOT, '0x' + 'a1'.repeat(32));

        await processorEndpoint.connect(signers[3]).adminReset();

        // Both A and B have pending work again. With the cursor back at index 0 the scan
        // starts from A; had the reset left it at index 1, B would be served first.
        await submitProcess(appA, '0x02');
        await submitProcess(appB, '0x03');
        const [appId] = await processorEndpoint.getPendingRequestsWithStateRoot(1);
        expect(appId).to.equal(appA);
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
          processorEndpoint.connect(signers[0]).adminResetApps([])
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });
    });

    describe('happy paths', function () {
      it('clears state roots for all deployed apps when called with empty appIds', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        expect(await processorEndpoint.applicationStateRoots(applicationId)).to.not.equal(
          BYTES32_ZERO
        );

        await processorEndpoint.connect(signers[3]).adminResetApps([]);
        expect(await processorEndpoint.applicationStateRoots(applicationId)).to.equal(BYTES32_ZERO);
      });

      it('clears state root only for the specified app when explicit appIds given', async () => {
        const { applicationId: id1 } = await bootstrapApplication(processorEndpoint);
        const { applicationId: id2 } = await bootstrapApplication(processorEndpoint);

        await processorEndpoint.connect(signers[3]).adminResetApps([id1]);
        expect(await processorEndpoint.applicationStateRoots(id1)).to.equal(BYTES32_ZERO);
        expect(await processorEndpoint.applicationStateRoots(id2)).to.not.equal(BYTES32_ZERO);
      });

      it('frees deploy slots for cleared apps', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const slotsBefore = await processorEndpoint.availableDeploySlots();

        await processorEndpoint.connect(signers[3]).adminResetApps([applicationId]);
        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore + 1n);
      });

      it('also clears the pending request queue', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        await processorEndpoint.submitRequest(
          0,
          applicationId,
          1,
          '0x01',
          ETH_TOKEN,
          0,
          minFeePerRequest,
          {
            value: minFeePerRequest,
          }
        );
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(1n);

        await processorEndpoint.connect(signers[3]).adminResetApps([]);
        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(0n);
      });

      it('refunds ETH asset deposits to the original request senders', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const sender = signers[0];
        const senderAddr = await sender.getAddress();
        // Deposit ETH as a business asset (assetAmount > 0, ETH token)
        const assetAmount = ethers.parseEther('0.5');
        await processorEndpoint
          .connect(sender)
          .submitRequest(0, applicationId, 1, '0x01', ETH_TOKEN, assetAmount, minFeePerRequest, {
            value: assetAmount + minFeePerRequest,
          });
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(assetAmount);

        await processorEndpoint.connect(signers[3]).adminResetApps([]);

        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
        expect(await processorEndpoint.totalAppCustody(ETH_TOKEN)).to.equal(0n);
        expect(await processorEndpoint.pendingClaims(ETH_TOKEN, senderAddr)).to.equal(assetAmount);
      });

      it('refunds ERC-20 asset deposits to the original request senders', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 500n;
        const sender = signers[0];
        const senderAddr = await sender.getAddress();

        await mockERC20.mint(senderAddr, assetAmount);
        await mockERC20.connect(sender).approve(await processorEndpoint.getAddress(), assetAmount);
        await processorEndpoint
          .connect(sender)
          .submitRequest(0, applicationId, 1, '0x01', tokenAddr, assetAmount, minFeePerRequest, {
            value: minFeePerRequest,
          });

        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(assetAmount);
        expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(assetAmount);

        const resetOperatorAddr = await signers[3].getAddress();
        const resetOpBalBefore = await mockERC20.balanceOf(resetOperatorAddr);

        await processorEndpoint.connect(signers[3]).adminResetApps([]);

        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(0n);
        expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(0n);
        expect(await processorEndpoint.pendingClaims(tokenAddr, senderAddr)).to.equal(assetAmount);
        expect(await mockERC20.balanceOf(resetOperatorAddr)).to.equal(resetOpBalBefore);
      });

      it('removes reset appIds from getDeployedAppIds', async () => {
        const { applicationId: id1 } = await bootstrapApplication(processorEndpoint);
        const { applicationId: id2 } = await bootstrapApplication(processorEndpoint);

        await processorEndpoint.connect(signers[3]).adminResetApps([id1]);
        const ids = await processorEndpoint.getDeployedAppIds();
        expect(ids.length).to.equal(1);
        expect(ids[0]).to.equal(id2);
      });

      it('clears getDeployedAppIds when all apps are reset', async () => {
        await bootstrapApplication(processorEndpoint);
        await bootstrapApplication(processorEndpoint);

        await processorEndpoint.connect(signers[3]).adminResetApps([]);
        expect(await processorEndpoint.getDeployedAppIds()).to.deep.equal([]);
      });

      it('rescues ERC-20 custody using getAllowedTokens', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        const tokenAddr = await mockERC20.getAddress();
        const assetAmount = 200n;

        await mockERC20.mint(await signers[0].getAddress(), assetAmount);
        await mockERC20
          .connect(signers[0])
          .approve(await processorEndpoint.getAddress(), assetAmount);
        await processorEndpoint.submitRequest(
          0,
          applicationId,
          1,
          '0x01',
          tokenAddr,
          assetAmount,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await processorEndpoint.connect(signers[3]).adminResetApps([]);
        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(0n);
      });
    });
  });
});
