import { expect } from 'chai';
import { ethers, upgrades } from 'hardhat';
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

  // AuthorityRegistry is deployed behind a UUPS proxy (docs/design/UPGRADABLE_CONTRACTS_DESIGN.md):
  // all of the former constructor logic now runs in `initialize`, called through
  // `upgrades.deployProxy`'s proxy-construction calldata, so a revert here surfaces as a failed
  // proxy deployment rather than a failed `.deploy()` call.
  describe('initialize', function () {
    describe('unhappy paths', function () {
      it('reverts with OwnableInvalidOwner when owner is zero address', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');

        await expect(
          upgrades.deployProxy(
            AuthorityRegistry,
            [ADDRESS_ZERO, await defaultAuthority.getAddress()],
            {
              kind: 'uups',
            }
          )
        ).to.be.revertedWithCustomError(AuthorityRegistry, 'OwnableInvalidOwner');
      });

      it('reverts with AddressCantBeZero when defaultAuthority is zero address', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');

        await expect(
          upgrades.deployProxy(AuthorityRegistry, [await signers[0].getAddress(), ADDRESS_ZERO], {
            kind: 'uups',
          })
        ).to.be.revertedWithCustomError(AuthorityRegistry, 'AddressCantBeZero');
      });

      it('reverts with InvalidInitialization when initialize is called a second time', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
        const owner = await signers[0].getAddress();
        const defaultAuthorityAddress = await defaultAuthority.getAddress();

        const authorityRegistry = await upgrades.deployProxy(
          AuthorityRegistry,
          [owner, defaultAuthorityAddress],
          { kind: 'uups' }
        );

        await expect(
          authorityRegistry.initialize(owner, defaultAuthorityAddress)
        ).to.be.revertedWithCustomError(AuthorityRegistry, 'InvalidInitialization');
      });
    });

    describe('happy paths', function () {
      it('initializes owner and default authority contract', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
        const owner = await signers[0].getAddress();
        const defaultAuthorityAddress = await defaultAuthority.getAddress();

        const authorityRegistry = await upgrades.deployProxy(
          AuthorityRegistry,
          [owner, defaultAuthorityAddress],
          { kind: 'uups' }
        );

        expect(await authorityRegistry.owner()).to.equal(owner);
        expect(await authorityRegistry.defaultAuthorityContract()).to.equal(
          defaultAuthorityAddress
        );
      });

      it('emits DefaultAuthorityContractSet on deployment', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
        const defaultAuthorityAddress = await defaultAuthority.getAddress();

        const authorityRegistry = await upgrades.deployProxy(
          AuthorityRegistry,
          [await signers[0].getAddress(), defaultAuthorityAddress],
          { kind: 'uups' }
        );

        await expect(authorityRegistry.deploymentTransaction())
          .to.emit(authorityRegistry, 'DefaultAuthorityContractSet')
          .withArgs(defaultAuthorityAddress);
      });
    });
  });
});
