import { expect } from 'chai';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint deployer role management', function () {
  let processorEndpoint: any;
  let signers: any[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;

  const REQUEST_TYPE_DEPLOYAPP = 0;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  it('admin can add a second allowed deployer', async () => {
    const deployerAddress = await signers[4].getAddress();

    await processorEndpoint.connect(signers[2]).addAllowedDeployer(deployerAddress);

    expect(await processorEndpoint.isAllowedDeployer(deployerAddress)).to.equal(true);

    await expect(
      processorEndpoint
        .connect(signers[4])
        .submitDeployRequest(0, '0x01', { value: minFeePerRequest })
    ).to.emit(processorEndpoint, 'RequestSubmitted');
  });

  it('submitRequest reverts with InvalidRequestType for DEPLOYAPP even with deployer role', async () => {
    await expect(
      processorEndpoint
        .connect(signers[2])
        .submitRequest(0, applicationId, REQUEST_TYPE_DEPLOYAPP, '0x01', 0, minFeePerRequest, {
          value: minFeePerRequest,
        })
    ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestType');
  });

  it('admin can remove an allowed deployer', async () => {
    const deployerAddress = await signers[4].getAddress();

    await processorEndpoint.connect(signers[2]).addAllowedDeployer(deployerAddress);
    await processorEndpoint.connect(signers[2]).removeAllowedDeployer(deployerAddress);

    expect(await processorEndpoint.isAllowedDeployer(deployerAddress)).to.equal(false);
    await expect(
      processorEndpoint
        .connect(signers[4])
        .submitDeployRequest(0, '0x01', { value: minFeePerRequest })
    ).to.be.revertedWithCustomError(processorEndpoint, 'DeployerNotAllowed');
  });

  it('non-admin cannot manage deployers', async () => {
    const deployerAddress = await signers[4].getAddress();

    await expect(
      processorEndpoint.connect(signers[0]).addAllowedDeployer(deployerAddress)
    ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');

    await expect(
      processorEndpoint.connect(signers[0]).removeAllowedDeployer(deployerAddress)
    ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
  });

  it('rejects zero address in deployer management methods', async () => {
    await expect(
      processorEndpoint
        .connect(signers[2])
        .addAllowedDeployer('0x0000000000000000000000000000000000000000')
    ).to.be.revertedWithCustomError(processorEndpoint, 'AddressCantBeZero');

    await expect(
      processorEndpoint
        .connect(signers[2])
        .removeAllowedDeployer('0x0000000000000000000000000000000000000000')
    ).to.be.revertedWithCustomError(processorEndpoint, 'AddressCantBeZero');
  });
});
