import { expect } from 'chai';
import { deployProcessorEndpointFixture } from './fixture';

describe('ProcessorEndpoint deployer role management', function () {
  let processorEndpoint: any;
  let signers: any[];
  let minFeePerRequest: bigint;

  const REQUEST_TYPE_DEPLOYAPP = 0;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
  });

  it('admin can add a second allowed deployer', async () => {
    const deployerAddress = await signers[4].getAddress();

    await processorEndpoint.connect(signers[2]).addAllowedDeployer(deployerAddress);

    expect(await processorEndpoint.isAllowedDeployer(deployerAddress)).to.equal(true);

    await expect(
      processorEndpoint
        .connect(signers[4])
        .submitRequest(0, 1, REQUEST_TYPE_DEPLOYAPP, '0x01', 0, minFeePerRequest, {
          value: minFeePerRequest,
        })
    ).to.emit(processorEndpoint, 'RequestSubmitted');
  });

  it('admin can remove an allowed deployer', async () => {
    const deployerAddress = await signers[4].getAddress();

    await processorEndpoint.connect(signers[2]).addAllowedDeployer(deployerAddress);
    await processorEndpoint.connect(signers[2]).removeAllowedDeployer(deployerAddress);

    expect(await processorEndpoint.isAllowedDeployer(deployerAddress)).to.equal(false);
    await expect(
      processorEndpoint
        .connect(signers[4])
        .submitRequest(0, 1, REQUEST_TYPE_DEPLOYAPP, '0x01', 0, minFeePerRequest, {
          value: minFeePerRequest,
        })
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
