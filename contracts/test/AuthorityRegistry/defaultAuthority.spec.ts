import { expect } from 'chai';
import { Signer } from 'ethers';
import { APPLICATION_ID } from '../util';
import { deployAuthorityRegistryFixture } from './fixture';

describe('AuthorityRegistry Test', function () {
  describe('DefaultAuthority', function () {
    let signers: Signer[];
    let authorityRegistry: any;
    let defaultAuthority: any;
    let testAddr: string;

    beforeEach(async function () {
      const fixture = await deployAuthorityRegistryFixture();
      signers = fixture.signers;
      defaultAuthority = fixture.defaultAuthority;
      authorityRegistry = fixture.authorityRegistry;
      testAddr = await signers[1].getAddress();
    });

    describe('unhappy paths', function () {
      it('reverts when non-owner adds an authority', async () => {
        await expect(
          defaultAuthority.connect(signers[1]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, 'OwnableUnauthorizedAccount');
      });

      it('reverts when non-owner removes an authority', async () => {
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await expect(
          defaultAuthority.connect(signers[1]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, 'OwnableUnauthorizedAccount');
      });

      it('reverts when adding an already-present authority', async () => {
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await expect(
          defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, 'AuthorityAlreadyPresent');
      });

      it('reverts when removing a non-present authority', async () => {
        await expect(
          defaultAuthority.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, 'AuthorityNotPresent');
      });
    });

    describe('happy paths', function () {
      it('allows owner to add authority and registry reflects it', async () => {
        await expect(
          defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, 'AddedAuthority').withArgs(APPLICATION_ID, testAddr);

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(true);
      });

      it('allows owner to remove authority and registry reflects it', async () => {
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await expect(
          defaultAuthority.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, 'RemovedAuthority').withArgs(APPLICATION_ID, testAddr);

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(false);
      });
    });
  });
});
