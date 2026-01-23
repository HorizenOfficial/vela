import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { APPLICATION_ID } from '../util';
import { deployAuthorityRegistryFixture } from './fixture';

describe('AuthorityRegistry Test', function () {
  describe('setAppAuthorityContract', function () {
    let signers: Signer[];
    let authorityRegistry: any;
    let testAddr: string;

    beforeEach(async function () {
      const fixture = await deployAuthorityRegistryFixture();
      signers = fixture.signers;
      authorityRegistry = fixture.authorityRegistry;
      testAddr = await signers[1].getAddress();
    });

    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await expect(
          authorityRegistry.connect(signers[1]).setAppAuthorityContract(
            APPLICATION_ID,
            await customAuthority.getAddress()
          )
        ).to.be.revertedWithCustomError(authorityRegistry, 'OwnableUnauthorizedAccount');
      });
    });

    describe('happy paths', function () {
      it('sets app-specific authority contract and emits AppAuthorityContractSet', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await customAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        const tx = await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
          APPLICATION_ID,
          await customAuthority.getAddress()
        );

        await expect(tx)
          .to.emit(authorityRegistry, 'AppAuthorityContractSet')
          .withArgs(APPLICATION_ID, await customAuthority.getAddress());

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(true);
      });
    });
  });
});
