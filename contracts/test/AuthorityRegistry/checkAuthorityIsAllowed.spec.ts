import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { ADDRESS_ZERO, APPLICATION_ID } from '../util';
import { deployAuthorityRegistryFixture } from './fixture';

describe('AuthorityRegistry Test', function () {
  describe('checkAuthorityIsAllowed', function () {
    let signers: Signer[];
    let authorityRegistry: any;
    let defaultAuthority: any;
    let testAddr: string;

    beforeEach(async function () {
      const fixture = await deployAuthorityRegistryFixture();
      signers = fixture.signers;
      authorityRegistry = fixture.authorityRegistry;
      defaultAuthority = fixture.defaultAuthority;
      testAddr = await signers[1].getAddress();
    });

    describe('unhappy paths', function () {
      it('returns false when the selected checker disallows the authority', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, await customAuthority.getAddress());

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(false);
      });

      it('reverts when the app-specific authority contract is an EOA or non-compliant', async () => {
        const eoa = await signers[2].getAddress();
        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(APPLICATION_ID, eoa);

        await expect(authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr)).to.be
          .reverted;
      });
    });

    describe('happy paths', function () {
      it('uses default authority checker when no app-specific contract is set', async () => {
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(true);

        const otherAddr = await signers[3].getAddress();
        const notAllowed = await authorityRegistry.checkAuthorityIsAllowed(
          APPLICATION_ID,
          otherAddr
        );
        expect(notAllowed).to.equal(false);
      });

      it('uses app-specific authority checker when it is set', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());
        await customAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, await customAuthority.getAddress());

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(true);
      });

      it('does not leak default authorities between different application ids', async () => {
        const appA = APPLICATION_ID;
        const appB = APPLICATION_ID + 1;

        await defaultAuthority.connect(signers[0]).addAllowedAuthority(appA, testAddr);

        const allowedA = await authorityRegistry.checkAuthorityIsAllowed(appA, testAddr);
        expect(allowedA).to.equal(true);

        const allowedB = await authorityRegistry.checkAuthorityIsAllowed(appB, testAddr);
        expect(allowedB).to.equal(false);
      });

      it('custom authority contract overrides default for a specific application', async () => {
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, await customAuthority.getAddress());

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(false);
      });

      it('custom authority contract only affects its application id', async () => {
        const appA = APPLICATION_ID;
        const appB = APPLICATION_ID + 1;

        await defaultAuthority.connect(signers[0]).addAllowedAuthority(appB, testAddr);

        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(appA, await customAuthority.getAddress());

        const allowedA = await authorityRegistry.checkAuthorityIsAllowed(appA, testAddr);
        expect(allowedA).to.equal(false);

        const allowedB = await authorityRegistry.checkAuthorityIsAllowed(appB, testAddr);
        expect(allowedB).to.equal(true);
      });

      it('changing custom authority contract updates behavior', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const custom1 = await DefaultAuthority.deploy(await signers[0].getAddress());
        const custom2 = await DefaultAuthority.deploy(await signers[0].getAddress());
        await custom2.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, await custom1.getAddress());
        const withCustom1 = await authorityRegistry.checkAuthorityIsAllowed(
          APPLICATION_ID,
          testAddr
        );
        expect(withCustom1).to.equal(false);

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, await custom2.getAddress());
        const withCustom2 = await authorityRegistry.checkAuthorityIsAllowed(
          APPLICATION_ID,
          testAddr
        );
        expect(withCustom2).to.equal(true);
      });

      it('setting app-specific contract to address(0) reverts to default behavior', async () => {
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, await customAuthority.getAddress());

        const withCustom = await authorityRegistry.checkAuthorityIsAllowed(
          APPLICATION_ID,
          testAddr
        );
        expect(withCustom).to.equal(false);

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, ADDRESS_ZERO);

        const afterReset = await authorityRegistry.checkAuthorityIsAllowed(
          APPLICATION_ID,
          testAddr
        );
        expect(afterReset).to.equal(true);
      });

      it('emits AppAuthorityContractSet when app-specific contract is set to address(0)', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, await customAuthority.getAddress());

        const tx = await authorityRegistry
          .connect(signers[0])
          .setAppAuthorityContract(APPLICATION_ID, ADDRESS_ZERO);

        await expect(tx)
          .to.emit(authorityRegistry, 'AppAuthorityContractSet')
          .withArgs(APPLICATION_ID, ADDRESS_ZERO);
      });
    });
  });
});
