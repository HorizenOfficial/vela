import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';
import { ETH_TOKEN, REQUEST_TYPE_DEPLOYAPP } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let defaultAuthority: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let bootstrapApplication: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    defaultAuthority = fixture.defaultAuthority;
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
  });

  describe('submitDeployRequest', function () {
    describe('unhappy paths', function () {
      it('reverts with InvalidProtocolVersion when protocolVersion is invalid', async () => {
        await expect(
          processorEndpoint.connect(signers[2]).submitDeployRequest(1, '0x01', {
            value: minFeePerRequest,
          })
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidProtocolVersion');
      });

      it('reverts with DeployerNotAllowed when sender lacks DEPLOYER_ROLE', async () => {
        await expect(
          processorEndpoint.connect(signers[0]).submitDeployRequest(0, '0x01', {
            value: minFeePerRequest,
          })
        ).to.be.revertedWithCustomError(processorEndpoint, 'DeployerNotAllowed');
      });

      it('reverts with FeeValueBelowMinimum when msg.value < minFeePerRequest', async () => {
        await expect(
          processorEndpoint.connect(signers[2]).submitDeployRequest(0, '0x01', {
            value: minFeePerRequest - 1n,
          })
        ).to.be.revertedWithCustomError(processorEndpoint, 'FeeValueBelowMinimum');
      });

      it('reverts with MaxNumOfApplicationsExceeded when max number of applications have been deployed', async () => {
        await processorEndpoint.connect(signers[2]).updateMaxNumOfApplications(1);
        await bootstrapApplication(processorEndpoint);

        await expect(
          processorEndpoint
            .connect(signers[2])
            .submitDeployRequest(0, '0x01', { value: minFeePerRequest })
        ).to.be.revertedWithCustomError(processorEndpoint, 'MaxNumOfApplicationsExceeded');
      });

      it('reverts with QueueThresholdExceeded when queue is full', async () => {
        const { applicationId } = await bootstrapApplication(processorEndpoint);
        await processorEndpoint.connect(signers[2]).updateQueueThreshold(2);

        await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(0, '0x01', { value: minFeePerRequest });

        await processorEndpoint.submitRequest(
          0,
          applicationId,
          1,
          '0x02',
          ETH_TOKEN,
          0,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        await expect(
          processorEndpoint
            .connect(signers[2])
            .submitDeployRequest(0, '0x03', { value: minFeePerRequest })
        ).to.be.revertedWithCustomError(processorEndpoint, 'QueueThresholdExceeded');
      });
    });

    describe('happy paths', function () {
      it('emits DeployRequestSubmitted and stores request data for retrieval', async () => {
        const protocolVersion = 0;
        const payload = '0x01';
        const maxFeeValue = minFeePerRequest;

        const tx = await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(protocolVersion, payload, {
            value: maxFeeValue,
          });

        await expect(tx).to.emit(processorEndpoint, 'DeployRequestSubmitted');
        const receipt = await tx.wait();
        const requestId = receipt.logs[0].args.requestId;
        expect(receipt.logs[0].args.sender).to.equal(await signers[2].getAddress());
        const stored = await processorEndpoint.requestById(requestId);

        expect(stored.requestId).to.equal(requestId);
        expect(stored.protocolVersion).to.equal(protocolVersion);
        expect(stored.requestType).to.equal(REQUEST_TYPE_DEPLOYAPP);
        expect(stored.payload).to.equal(payload);
        expect(stored.assetAmount).to.equal(0);
        expect(stored.maxFeeValue).to.equal(maxFeeValue);
        expect(stored.sender).to.equal(await signers[2].getAddress());
      });

      it('updates balances when request includes a fee', async () => {
        const maxFeeValue = minFeePerRequest;
        const deployer = signers[2];

        const processorBalanceBefore = await deployer.provider!.getBalance(
          await processorEndpoint.getAddress()
        );
        const userBalanceBefore = await deployer.provider!.getBalance(await deployer.getAddress());

        const tx = await processorEndpoint
          .connect(deployer)
          .submitDeployRequest(0, '0x02', { value: maxFeeValue });
        const receipt = await tx.wait();
        const gasCost = receipt.gasUsed * receipt.gasPrice;

        const userBalanceAfter = await deployer.provider!.getBalance(await deployer.getAddress());
        const processorBalanceAfter = await deployer.provider!.getBalance(
          await processorEndpoint.getAddress()
        );

        expect(userBalanceAfter).to.equal(userBalanceBefore - gasCost - maxFeeValue);
        expect(processorBalanceAfter).to.equal(processorBalanceBefore + maxFeeValue);

        // appLockedFunds should remain 0 for deploys (no deposit, fees tracked globally)
        const [pendingReq] = await processorEndpoint.getNextPendingRequest();
        const deployAppId = pendingReq.applicationId;
        expect(await processorEndpoint.appCustody(deployAppId, ETH_TOKEN)).to.equal(0n);
      });

      it('enqueues the deploy request as a pending request', async () => {
        await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(0, '0x01', { value: minFeePerRequest });

        const requests = await processorEndpoint.getPendingRequests();
        expect(requests.length).to.equal(1);
        expect(requests[0].requestType).to.equal(REQUEST_TYPE_DEPLOYAPP);
      });
    });
  });
});
