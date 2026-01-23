import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { ADDRESS_ZERO } from '../util';
import { deployDefaultAuthorityFixture } from './fixture';

describe('AuthorityRegistry Test', function () {
  let signers: Signer[];
  let defaultAuthority: any;

  beforeEach(async function () {
    const fixture = await deployDefaultAuthorityFixture();
    signers = fixture.signers;
    defaultAuthority = fixture.defaultAuthority;
  });

  describe('constructor', function () {
    describe('unhappy paths', function () {
      it('reverts with OwnableInvalidOwner when owner is zero address', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');

        await expect(
          AuthorityRegistry.deploy(ADDRESS_ZERO, await defaultAuthority.getAddress())
        ).to.be.revertedWithCustomError(AuthorityRegistry, 'OwnableInvalidOwner');
      });

      it('reverts with AddressCantBeZero when defaultAuthority is zero address', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');

        await expect(
          AuthorityRegistry.deploy(await signers[0].getAddress(), ADDRESS_ZERO)
        ).to.be.revertedWithCustomError(AuthorityRegistry, 'AddressCantBeZero');
      });
    });

    describe('happy paths', function () {
      it('initializes owner and default authority contract', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
        const owner = await signers[0].getAddress();
        const defaultAuthorityAddress = await defaultAuthority.getAddress();

        const authorityRegistry = await AuthorityRegistry.deploy(owner, defaultAuthorityAddress);

        expect(await authorityRegistry.owner()).to.equal(owner);
        expect(await authorityRegistry.defaultAuthorityContract()).to.equal(defaultAuthorityAddress);
      });

      it('emits DefaultAuthorityContractSet on deployment', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
        const defaultAuthorityAddress = await defaultAuthority.getAddress();

        const authorityRegistry = await AuthorityRegistry.deploy(
          await signers[0].getAddress(),
          defaultAuthorityAddress
        );

        await expect(authorityRegistry.deploymentTransaction())
          .to.emit(authorityRegistry, 'DefaultAuthorityContractSet')
          .withArgs(defaultAuthorityAddress);
      });
    });
  });
});
