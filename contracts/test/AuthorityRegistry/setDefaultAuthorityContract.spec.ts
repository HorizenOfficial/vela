import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { ADDRESS_ZERO, APPLICATION_ID } from '../util';
import { deployAuthorityRegistryFixture } from './fixture';

describe('AuthorityRegistry Test', function () {
  describe('setDefaultAuthorityContract', function () {
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
      it('reverts when caller is not owner', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());

        await expect(
          authorityRegistry.connect(signers[1]).setDefaultAuthorityContract(await newDefault.getAddress())
        ).to.be.revertedWithCustomError(authorityRegistry, 'OwnableUnauthorizedAccount');
      });

      it('reverts with AddressCantBeZero when authorityContract is zero address', async () => {
        await expect(
          authorityRegistry.connect(signers[0]).setDefaultAuthorityContract(ADDRESS_ZERO)
        ).to.be.revertedWithCustomError(authorityRegistry, 'AddressCantBeZero');
      });
    });

    describe('happy paths', function () {
      it('updates default authority contract and emits DefaultAuthorityContractSet', async () => {
        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());

        const tx = await authorityRegistry
          .connect(signers[0])
          .setDefaultAuthorityContract(await newDefault.getAddress());

        await expect(tx)
          .to.emit(authorityRegistry, 'DefaultAuthorityContractSet')
          .withArgs(await newDefault.getAddress());

        expect(await authorityRegistry.defaultAuthorityContract()).to.equal(await newDefault.getAddress());
      });

      it('affects applications without custom authority contracts', async () => {
        const before = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(before).to.equal(false);

        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());
        await newDefault.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await authorityRegistry.connect(signers[0]).setDefaultAuthorityContract(
          await newDefault.getAddress()
        );

        const after = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(after).to.equal(true);
      });

      it('does not affect applications with custom authority contracts', async () => {
        const APP_CUSTOM = APPLICATION_ID;
        const APP_DEFAULT = APPLICATION_ID + 1;

        const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());
        await customAuthority.connect(signers[0]).addAllowedAuthority(APP_CUSTOM, testAddr);

        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
          APP_CUSTOM,
          await customAuthority.getAddress()
        );

        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());
        await newDefault.connect(signers[0]).addAllowedAuthority(APP_DEFAULT, testAddr);

        await authorityRegistry.connect(signers[0]).setDefaultAuthorityContract(
          await newDefault.getAddress()
        );

        const afterCustom = await authorityRegistry.checkAuthorityIsAllowed(APP_CUSTOM, testAddr);
        expect(afterCustom).to.equal(true);

        const afterDefault = await authorityRegistry.checkAuthorityIsAllowed(APP_DEFAULT, testAddr);
        expect(afterDefault).to.equal(true);
      });
    });
  });
});
