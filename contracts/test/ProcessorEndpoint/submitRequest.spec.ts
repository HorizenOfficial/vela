import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';
import {
  ETH_TOKEN,
  REQUEST_TYPE_ASSOCIATEKEY,
  REQUEST_TYPE_DEANONYMIZATION,
  REQUEST_TYPE_DEPLOYAPP,
  REQUEST_TYPE_PLAINPROCESS,
  REQUEST_TYPE_PROCESS,
} from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let defaultAuthority: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint } = await fixture.deployProcessorEndpoint());
    defaultAuthority = fixture.defaultAuthority;
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  describe('submitRequest', function () {
    describe('unhappy paths', function () {
      it('reverts with InvalidProtocolVersion when protocolVersion is invalid', async () => {
        await expect(
          processorEndpoint.submitRequest(
            1,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            ETH_TOKEN,
            0,
            minFeePerRequest,
            {
              value: minFeePerRequest,
            }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidProtocolVersion');
      });

      it('reverts with InvalidApplicationId when applicationId is invalid', async () => {
        await expect(
          processorEndpoint.submitRequest(
            0,
            999,
            REQUEST_TYPE_PROCESS,
            '0x01',
            ETH_TOKEN,
            0,
            minFeePerRequest,
            {
              value: minFeePerRequest,
            }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidApplicationId');
      });

      it('reverts with InvalidValue when msg.value != assetAmount + maxFeeValue', async () => {
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            ETH_TOKEN,
            1,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InvalidValue when msg.value is less than assetAmount + maxFeeValue', async () => {
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            ETH_TOKEN,
            2,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InvalidValue when msg.value is greater than assetAmount + maxFeeValue', async () => {
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            ETH_TOKEN,
            1,
            minFeePerRequest,
            { value: minFeePerRequest + 2n }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with FeeValueBelowMinimum when maxFeeValue < minFeePerRequest', async () => {
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x01',
            ETH_TOKEN,
            0,
            minFeePerRequest - 1n,
            { value: minFeePerRequest - 1n }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'FeeValueBelowMinimum');
      });

      it('reverts with QueueThresholdExceeded when queue is full', async () => {
        await processorEndpoint.connect(signers[2]).updateQueueThreshold(1);

        await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          ETH_TOKEN,
          0,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x02',
            ETH_TOKEN,
            0,
            minFeePerRequest,
            {
              value: minFeePerRequest,
            }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'QueueThresholdExceeded');
      });

      it('reverts with InvalidPayload when ASSOCIATEKEY payload length is not 133 or 226', async () => {
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_ASSOCIATEKEY,
            '0x11',
            ETH_TOKEN,
            0,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with AuthorityNotAllowed when DEANONYMIZATION sender is not allowed', async () => {
        await expect(
          processorEndpoint.submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_DEANONYMIZATION,
            '0x01',
            ETH_TOKEN,
            0,
            minFeePerRequest,
            { value: minFeePerRequest }
          )
        ).to.be.revertedWithCustomError(processorEndpoint, 'AuthorityNotAllowed');
      });

      it('reverts with InvalidRequestType when requestType is DEPLOYAPP', async () => {
        await expect(
          processorEndpoint
            .connect(signers[2])
            .submitRequest(
              0,
              applicationId,
              REQUEST_TYPE_DEPLOYAPP,
              '0x01',
              ETH_TOKEN,
              0,
              minFeePerRequest,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestType');
      });
    });

    describe('happy paths', function () {
      it('emits RequestSubmitted and stores request data for retrieval', async () => {
        const protocolVersion = 0;
        const requestType = REQUEST_TYPE_PROCESS;
        const payload = '0x01';
        const assetAmount = 5n;
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

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        const requestId = receipt.logs[0].args.requestId;
        expect(receipt.logs[0].args.sender).to.equal(await signers[0].getAddress());
        const stored = await processorEndpoint.requestById(requestId);

        expect(stored.requestId).to.equal(requestId);
        expect(stored.protocolVersion).to.equal(protocolVersion);
        expect(stored.applicationId).to.equal(applicationId);
        expect(stored.requestType).to.equal(requestType);
        expect(stored.payload).to.equal(payload);
        expect(stored.tokenAddress).to.equal(ETH_TOKEN);
        expect(stored.assetAmount).to.equal(assetAmount);
        expect(stored.maxFeeValue).to.equal(maxFeeValue);
        expect(stored.sender).to.equal(await signers[0].getAddress());
      });

      it('updates balances when request includes a deposit and fee', async () => {
        const protocolVersion = 0;
        const requestType = REQUEST_TYPE_PROCESS;
        const assetAmount = 100n;
        const maxFeeValue = minFeePerRequest;

        const processorBalanceBefore = await signers[0].provider!.getBalance(
          await processorEndpoint.getAddress()
        );
        const userBalanceBefore = await signers[0].provider!.getBalance(
          await signers[0].getAddress()
        );

        const tx = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          requestType,
          '0x02',
          ETH_TOKEN,
          assetAmount,
          maxFeeValue,
          { value: assetAmount + maxFeeValue }
        );
        const receipt = await tx.wait();
        const gasCost = receipt.gasUsed * receipt.gasPrice;

        const userBalanceAfter = await signers[0].provider!.getBalance(
          await signers[0].getAddress()
        );
        const processorBalanceAfter = await signers[0].provider!.getBalance(
          await processorEndpoint.getAddress()
        );

        expect(userBalanceAfter).to.equal(
          userBalanceBefore - gasCost - (assetAmount + maxFeeValue)
        );
        expect(processorBalanceAfter).to.equal(processorBalanceBefore + assetAmount + maxFeeValue);

        // appLockedFunds should increase by assetAmount only (fees tracked globally)
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(assetAmount);
      });

      it('accepts non-deanonymization requests (PROCESS/ASSOCIATEKEY) and enqueues', async () => {
        const maxFeeValue = minFeePerRequest;
        const associatePayload = '0x' + '11'.repeat(226);

        await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x02',
          ETH_TOKEN,
          0,
          maxFeeValue,
          {
            value: maxFeeValue,
          }
        );
        await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_ASSOCIATEKEY,
          associatePayload,
          ETH_TOKEN,
          0,
          maxFeeValue,
          { value: maxFeeValue }
        );

        const requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(2);
        expect(requests[0].requestType).to.equal(REQUEST_TYPE_PROCESS);
        expect(requests[1].requestType).to.equal(REQUEST_TYPE_ASSOCIATEKEY);
      });

      it('accepts non-deanonymization request from an unauthorized authority', async () => {
        const tx = await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          ETH_TOKEN,
          0,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        expect(receipt.logs[0].args.sender).to.equal(await signers[0].getAddress());
      });

      it('accepts ASSOCIATEKEY when payload length is 133 (key only)', async () => {
        const payload = '0x' + '22'.repeat(133);
        const tx = await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_ASSOCIATEKEY,
          payload,
          ETH_TOKEN,
          0,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        expect(receipt.logs[0].args.sender).to.equal(await signers[0].getAddress());
      });

      it('accepts ASSOCIATEKEY when payload length is 226 (key + encrypted seed)', async () => {
        const payload = '0x' + '22'.repeat(226);
        const tx = await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_ASSOCIATEKEY,
          payload,
          ETH_TOKEN,
          0,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        expect(receipt.logs[0].args.sender).to.equal(await signers[0].getAddress());
      });

      it('accepts PLAINPROCESS and enqueues it like PROCESS', async () => {
        const payload = '0x' + 'ab'.repeat(32);
        const tx = await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PLAINPROCESS,
          payload,
          ETH_TOKEN,
          0,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        const requestId = receipt.logs[0].args.requestId;
        const stored = await processorEndpoint.requestById(requestId);
        expect(stored.requestType).to.equal(REQUEST_TYPE_PLAINPROCESS);
        expect(stored.payload).to.equal(payload);
        expect(stored.sender).to.equal(await signers[0].getAddress());
      });

      it('accepts DEANONYMIZATION for allowed authority with zero deposit', async () => {
        const authority = await signers[3].getAddress();
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(applicationId, authority);

        const tx = await processorEndpoint
          .connect(signers[3])
          .submitRequest(
            0,
            applicationId,
            REQUEST_TYPE_DEANONYMIZATION,
            '0x01',
            ETH_TOKEN,
            0,
            minFeePerRequest,
            {
              value: minFeePerRequest,
            }
          );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        expect(receipt.logs[0].args.sender).to.equal(authority);
      });
    });
  });
});
