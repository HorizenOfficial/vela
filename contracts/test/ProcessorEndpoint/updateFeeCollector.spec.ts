import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';
import { ADDRESS_ZERO, BYTES32_ZERO } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
  });

  async function submitBasicRequest(payload: string) {
    const protocolVersion = 0;
    const applicationId = 1;
    const requestType = 1;
    const depositAmount = 0n;
    const maxFeeValue = minFeePerRequest;

    const tx = await processorEndpoint.submitRequest(
      protocolVersion,
      applicationId,
      requestType,
      payload,
      depositAmount,
      maxFeeValue,
      { value: depositAmount + maxFeeValue }
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

      it('routes fees to the new feeCollector for markRequestCompleted', async () => {
        const newCollector = await signers[3].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const { requestId, maxFeeValue } = await submitBasicRequest('0x01');
        const collectorBalanceBefore = await signers[3].provider!.getBalance(newCollector);

        await processorEndpoint.connect(signers[1]).markRequestCompleted(requestId, 0, maxFeeValue);

        const collectorBalanceAfter = await signers[3].provider!.getBalance(newCollector);
        expect(collectorBalanceAfter - collectorBalanceBefore).to.equal(maxFeeValue);
      });

      it('routes fees to the new feeCollector for markRequestFailed', async () => {
        const newCollector = await signers[4].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const { requestId, maxFeeValue } = await submitBasicRequest('0x02');
        const collectorBalanceBefore = await signers[4].provider!.getBalance(newCollector);

        await processorEndpoint.connect(signers[1]).markRequestFailed(requestId, 1, 'failed');

        const collectorBalanceAfter = await signers[4].provider!.getBalance(newCollector);
        expect(collectorBalanceAfter - collectorBalanceBefore).to.equal(maxFeeValue);
      });

      it('routes fees to the new feeCollector for stateUpdate', async () => {
        const newCollector = await signers[5].getAddress();
        await processorEndpoint.connect(signers[2]).updateFeeCollector(newCollector);

        const { requestId, applicationId, maxFeeValue } = await submitBasicRequest('0x03');
        const collectorBalanceBefore = await signers[5].provider!.getBalance(newCollector);

        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            BYTES32_ZERO,
            '0x1000000000000000000000000000000000000000000000000000000000000000',
            requestId,
            ['0x'],
            [''],
            [],
            0,
            maxFeeValue,
            '0x'
          );

        const collectorBalanceAfter = await signers[5].provider!.getBalance(newCollector);
        expect(collectorBalanceAfter - collectorBalanceBefore).to.equal(maxFeeValue);
      });
    });
  });
});
