import { expect } from 'chai';
import { ethers, Signer } from 'ethers';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import { ADDRESS_ZERO, ETH_TOKEN } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  async function submitBasicRequest(payload: string) {
    const protocolVersion = 0;
    const requestType = 1;
    const assetAmount = 0n;
    const maxFeeValue = minFeePerRequest;

    const tx = await processorEndpoint.submitRequest(
      protocolVersion,
      applicationId,
      requestType,
      payload,
      ETH_TOKEN,
      assetAmount,
      maxFeeValue,
      { value: assetAmount + maxFeeValue }
    );
    const receipt = await tx.wait();
    return {
      requestId: receipt.logs[0].args.requestId,
      applicationId,
      maxFeeValue,
    };
  }

  describe('updateFeeCollector', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks ADMIN role', async () => {
        await expect(
          processorEndpoint.connect(signers[1]).updateFeeCollector(await signers[3].getAddress())
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with AddressCantBeZero when newFeeCollector is zero address', async () => {
        await expect(
          processorEndpoint.connect(signers[2]).updateFeeCollector(ADDRESS_ZERO)
        ).to.be.revertedWithCustomError(processorEndpoint, 'AddressCantBeZero');
      });
    });

    describe('happy paths', function () {
      it('updates feeCollector and emits FeeCollectorUpdated', async () => {
        const newCollector = await signers[3].getAddress();
        const tx = await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);
        await expect(tx).to.emit(processorEndpoint, 'FeeCollectorUpdated').withArgs(newCollector);
        expect(await processorEndpoint.feeCollector()).to.equal(newCollector);
      });

      it('routes fees to the new feeCollector for stateUpdate (completion)', async () => {
        const newCollector = await signers[3].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const { requestId, applicationId: appId, maxFeeValue } = await submitBasicRequest('0x01');
        // With pull pattern, funds are credited to pending deposits
        const collectorPendingAmountBefore = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          newCollector
        );

        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            appId,
            INITIAL_STATE_ROOT,
            '0x' + '01'.repeat(32),
            requestId,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            maxFeeValue,
            0,
            '',
            '0x'
          );

        const collectorPendingAmountAfter = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          newCollector
        );
        expect(collectorPendingAmountAfter - collectorPendingAmountBefore).to.equal(maxFeeValue);
      });

      it('routes fees to the new feeCollector for failed requests', async () => {
        const newCollector = await signers[4].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const { requestId, applicationId: appId } = await submitBasicRequest('0x02');
        // With pull pattern, funds are credited to pending deposits
        const collectorPendingAmountBefore = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          newCollector
        );

        const currentStateRoot = await processorEndpoint.applicationStateRoots(appId);
        // Fail request via stateUpdate with errorCode
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            appId,
            currentStateRoot,
            currentStateRoot,
            requestId,
            { events: [], subTypes: [] },
            { events: [], subTypes: [] },
            [],
            0,
            0,
            1,
            'failed',
            '0x'
          );

        const collectorPendingAmountAfter = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          newCollector
        );
        expect(collectorPendingAmountAfter - collectorPendingAmountBefore).to.equal(
          minFeePerRequest
        );
      });

      it('routes fees to the new feeCollector for stateUpdate', async () => {
        const newCollector = await signers[5].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const { requestId, applicationId: appId, maxFeeValue } = await submitBasicRequest('0x03');
        // With pull pattern, funds are credited to pending deposits
        const collectorPendingAmountBefore = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          newCollector
        );

        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            appId,
            INITIAL_STATE_ROOT,
            '0x1000000000000000000000000000000000000000000000000000000000000000',
            requestId,
            { events: ['0x'], subTypes: [ethers.encodeBytes32String('')] },
            { events: [], subTypes: [] },
            [],
            0,
            maxFeeValue,
            0,
            '',
            '0x'
          );

        const collectorPendingAmountAfter = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          newCollector
        );
        expect(collectorPendingAmountAfter - collectorPendingAmountBefore).to.equal(maxFeeValue);
      });
    });
  });
});
