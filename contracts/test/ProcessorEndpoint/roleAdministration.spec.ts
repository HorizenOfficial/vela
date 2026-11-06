import { expect } from 'chai';
import { Signer } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';
import { BYTES32_ZERO } from '../util';

// ADMIN gates every configuration setter and `_authorizeUpgrade`, and UPDATE_STATUS_ROLE is a hot
// key signing every state update, so both must be rotatable: an operator key that cannot be
// replaced turns a compromise into a permanent one, and a lost one would freeze configuration and
// upgradeability for good. RESET_OPERATOR is deliberately excluded — leaving DEFAULT_ADMIN_ROLE
// (held by nobody) as its admin is what makes the `address(0)` permanent-disable guarantee in
// docs/design/PROCESSOR_ENDPOINT_ADMIN_RESET.md hold.
describe('ProcessorEndpoint role administration', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let adminRole: string;
  let updateStatusRole: string;
  let resetRole: string;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    adminRole = await processorEndpoint.ADMIN();
    updateStatusRole = await processorEndpoint.UPDATE_STATUS_ROLE();
    resetRole = await processorEndpoint.RESET_OPERATOR();
  });

  it('makes ADMIN its own role admin', async () => {
    expect(await processorEndpoint.getRoleAdmin(adminRole)).to.equal(adminRole);
  });

  it('puts UPDATE_STATUS_ROLE under ADMIN', async () => {
    expect(await processorEndpoint.getRoleAdmin(updateStatusRole)).to.equal(adminRole);
  });

  it('admin can rotate ADMIN to a new holder', async () => {
    const newAdmin = await signers[5].getAddress();

    await processorEndpoint.connect(signers[2]).grantRole(adminRole, newAdmin);
    expect(await processorEndpoint.hasRole(adminRole, newAdmin)).to.equal(true);

    // The new holder can revoke the old one, which is what rotation after a key compromise needs.
    await processorEndpoint
      .connect(signers[5])
      .revokeRole(adminRole, await signers[2].getAddress());
    expect(await processorEndpoint.hasRole(adminRole, await signers[2].getAddress())).to.equal(
      false
    );

    await expect(processorEndpoint.connect(signers[5]).updateQueueThreshold(20)).to.not.be.reverted;
    await expect(
      processorEndpoint.connect(signers[2]).updateQueueThreshold(30)
    ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
  });

  it('admin can rotate UPDATE_STATUS_ROLE', async () => {
    const newOperator = await signers[5].getAddress();
    const oldOperator = await signers[1].getAddress();

    await processorEndpoint.connect(signers[2]).grantRole(updateStatusRole, newOperator);
    await processorEndpoint.connect(signers[2]).revokeRole(updateStatusRole, oldOperator);

    expect(await processorEndpoint.hasRole(updateStatusRole, newOperator)).to.equal(true);
    expect(await processorEndpoint.hasRole(updateStatusRole, oldOperator)).to.equal(false);
  });

  it('non-admin cannot grant ADMIN', async () => {
    await expect(
      processorEndpoint.connect(signers[0]).grantRole(adminRole, await signers[5].getAddress())
    ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
  });

  it('leaves RESET_OPERATOR under DEFAULT_ADMIN_ROLE so admin cannot grant it', async () => {
    expect(await processorEndpoint.getRoleAdmin(resetRole)).to.equal(BYTES32_ZERO);

    await expect(
      processorEndpoint.connect(signers[2]).grantRole(resetRole, await signers[5].getAddress())
    ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
  });
});
