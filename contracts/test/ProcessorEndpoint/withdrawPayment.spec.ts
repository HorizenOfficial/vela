import { expect } from 'chai';
import { BigNumberish, Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture } from './fixture';
import { BYTES32_ZERO } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let protocolVersion: BigNumberish;
  let applicationId: BigNumberish;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    protocolVersion = await processorEndpoint.PROTOCOL_VERSION();
    applicationId = await processorEndpoint.APPLICATION_ID();
  });

  describe('withdrawPayments', function () {
    describe('unhappy paths', function () {
      it('should prevent griefing by contracts that reject ETH transfers', async function () {
        // Deploy FallbackFailure contract
        let FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        let fallbackFailure = await FallbackFailure.deploy();

        // FallbackFailure submits a request
        await fallbackFailure.insertRequestOnProcessorEndpoint(
          processorEndpoint.getAddress(),
          protocolVersion,
          applicationId,
          1,
          '0x01',
          100,
          minFeePerRequest,
          { value: 100n + minFeePerRequest }
        );

        let [currentPendingRequest] = await processorEndpoint.getNextPendingRequest();

        // Mark request as failed via stateUpdate with errorCode
        // With pull pattern, we just credit to pending, don't transfer
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            BYTES32_ZERO,
            BYTES32_ZERO,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            0,
            1,
            'Test',
            '0x'
          );

        // Verify funds are in pending deposits
        let pendingAmount = await processorEndpoint.payments(await fallbackFailure.getAddress());
        expect(pendingAmount).eql(100n);

        // When someone tries to withdraw for FallbackFailure, it will fail
        // but the contract operation (stateUpdate with error) succeeded
        await expect(
          processorEndpoint.withdrawPayments(await fallbackFailure.getAddress())
        ).to.be.revertedWithCustomError(processorEndpoint, 'TransferFailed');

        // Funds remain in pending for FallbackFailure
        pendingAmount = await processorEndpoint.payments(await fallbackFailure.getAddress());
        expect(pendingAmount).eql(100n);
      });

      it('should allow stateUpdate to succeed even with contract receiver that rejects ETH', async function () {
        // Deploy FallbackFailure contract
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        const fallbackAddr = await fallbackFailure.getAddress();

        // Submit a normal request
        let submitTx = await processorEndpoint.submitRequest(
          protocolVersion,
          applicationId,
          1,
          '0x01',
          100,
          minFeePerRequest,
          { value: 100n + minFeePerRequest }
        );
        await submitTx.wait();

        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest] = await processorEndpoint.getNextPendingRequest();

        let newStateRoot = '0x1234000000000000000000000000000000000000000000000000000000000000';

        // Create signature for stateUpdate that includes withdrawal to FallbackFailure
        const { ethSignStateUpdate } = await import('../../scripts/util');
        let signature = await ethSignStateUpdate(
          signers[0],
          applicationId,
          initialStateRoot,
          newStateRoot,
          currentPendingRequest.requestId,
          [],
          [],
          [[fallbackAddr, 50]], // withdrawal to FallbackFailure
          0,
          minFeePerRequest
        );

        // stateUpdate should succeed (with pull pattern, we just credit to pending)
        let updateTx = await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [[fallbackAddr, 50]],
            0,
            minFeePerRequest,
            0,
            '',
            signature
          );
        await updateTx.wait();

        // State root should be updated
        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);

        // Funds should be in pending for FallbackFailure
        let pendingAmount = await processorEndpoint.payments(fallbackAddr);
        expect(pendingAmount).eql(50n);
      });
    });

    describe('happy paths', function () {
      it('should allow anyone to trigger withdrawal for any payee', async function () {
        // Submit a request from signer[2]
        let submitTx = await processorEndpoint
          .connect(signers[2])
          .submitRequest(protocolVersion, applicationId, 1, '0x01', 100, minFeePerRequest, {
            value: 100n + minFeePerRequest,
          });
        await submitTx.wait();

        let [currentPendingRequest] = await processorEndpoint.getNextPendingRequest();

        // Fail the request via stateUpdate with errorCode to credit refund to signer[2]
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            BYTES32_ZERO,
            BYTES32_ZERO,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            0,
            1,
            'Test',
            '0x'
          );

        // Check pending amount for signer[2]
        let pendingAmount = await processorEndpoint.payments(await signers[2].getAddress());
        expect(pendingAmount).eql(100n); // depositAmount + (maxFeeValue - minFeePerRequest) = 100 + 0 = 100

        // signer[3] (not the payee) can trigger withdrawal for signer[2]
        let balanceBefore = await ethers.provider.getBalance(await signers[2].getAddress());
        let tx = await processorEndpoint
          .connect(signers[3])
          .withdrawPayments(await signers[2].getAddress());
        await expect(tx)
          .to.emit(processorEndpoint, 'PaymentWithdrawn')
          .withArgs(await signers[2].getAddress(), BigInt(100n));
        let balanceAfter = await ethers.provider.getBalance(await signers[2].getAddress());

        expect(balanceAfter).eql(balanceBefore + 100n);
        expect(await processorEndpoint.payments(await signers[2].getAddress())).eql(0n);
      });

      it('should not revert withdrawPayments if no pending amount', async function () {
        // This should not revert, just do nothing
        await processorEndpoint.withdrawPayments(await signers[5].getAddress());
      });

      it('should accumulate multiple credits for same address', async function () {
        // Submit two requests from signer[2]
        let submitTx = await processorEndpoint
          .connect(signers[2])
          .submitRequest(protocolVersion, applicationId, 1, '0x01', 50, minFeePerRequest, {
            value: 50n + minFeePerRequest,
          });
        await submitTx.wait();

        submitTx = await processorEndpoint
          .connect(signers[2])
          .submitRequest(protocolVersion, applicationId, 1, '0x02', 60, minFeePerRequest, {
            value: 60n + minFeePerRequest,
          });
        await submitTx.wait();

        // Fail both requests via stateUpdate with errorCode
        let [req1] = await processorEndpoint.getNextPendingRequest();
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            BYTES32_ZERO,
            BYTES32_ZERO,
            req1.requestId,
            [],
            [],
            [],
            0,
            0,
            1,
            'Test',
            '0x'
          );

        let [req2] = await processorEndpoint.getNextPendingRequest();
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            BYTES32_ZERO,
            BYTES32_ZERO,
            req2.requestId,
            [],
            [],
            [],
            0,
            0,
            1,
            'Test',
            '0x'
          );

        // Check accumulated pending: 50 + 60 = 110
        let pendingAmount = await processorEndpoint.payments(await signers[2].getAddress());
        expect(pendingAmount).eql(110n);
      });
    });
  });
});
