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

  describe('initialize', function () {
    describe('unhappy paths', function () {
      it('reverts with OwnableInvalidOwner when owner is zero address', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');

        await expect(
          upgrades.deployProxy(
            AuthorityRegistry,
            [ADDRESS_ZERO, await defaultAuthority.getAddress()],
            { kind: 'uups' }
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

      it('reverts when initialize is called a second time', async () => {
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
        ).to.be.revertedWithCustomError(authorityRegistry, 'InvalidInitialization');
      });

      it('reverts when initialize is called directly on the implementation', async () => {
        const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
        const owner = await signers[0].getAddress();
        const defaultAuthorityAddress = await defaultAuthority.getAddress();

        const authorityRegistry = await upgrades.deployProxy(
          AuthorityRegistry,
          [owner, defaultAuthorityAddress],
          { kind: 'uups' }
        );
        const implementationAddress = await upgrades.erc1967.getImplementationAddress(
          await authorityRegistry.getAddress()
        );
        const implementation = AuthorityRegistry.attach(implementationAddress);

        await expect(
          (implementation as any).initialize(owner, defaultAuthorityAddress)
        ).to.be.revertedWithCustomError(authorityRegistry, 'InvalidInitialization');
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

  describe('upgradeability', function () {
    it('reverts when a non-owner attempts to upgrade', async () => {
      const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
      const owner = await signers[0].getAddress();
      const defaultAuthorityAddress = await defaultAuthority.getAddress();

      const authorityRegistry = await upgrades.deployProxy(
        AuthorityRegistry,
        [owner, defaultAuthorityAddress],
        { kind: 'uups' }
      );
      const AuthorityRegistryV2 = await ethers.getContractFactory('AuthorityRegistry', signers[1]);
      const v2Implementation = await AuthorityRegistryV2.deploy();
      await v2Implementation.waitForDeployment();

      await expect(
        (authorityRegistry.connect(signers[1]) as any).upgradeToAndCall(
          await v2Implementation.getAddress(),
          '0x'
        )
      ).to.be.revertedWithCustomError(authorityRegistry, 'OwnableUnauthorizedAccount');
    });

    it('preserves storage across an upgrade to a new implementation', async () => {
      const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
      const owner = await signers[0].getAddress();
      const defaultAuthorityAddress = await defaultAuthority.getAddress();

      const authorityRegistry = await upgrades.deployProxy(
        AuthorityRegistry,
        [owner, defaultAuthorityAddress],
        { kind: 'uups' }
      );
      const proxyAddress = await authorityRegistry.getAddress();

      const upgraded = await upgrades.upgradeProxy(proxyAddress, AuthorityRegistry);
      await upgraded.waitForDeployment();

      expect(await upgraded.getAddress()).to.equal(proxyAddress);
      expect(await upgraded.owner()).to.equal(owner);
      expect(await upgraded.defaultAuthorityContract()).to.equal(defaultAuthorityAddress);
    });
  });
});
