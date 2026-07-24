import { expect } from 'chai';
import { ethers, upgrades } from 'hardhat';
import { Signer } from 'ethers';

// These tests exercise TeeAuthenticator's upgradeability (initializer guards,
// UUPS authorization, storage preservation) without touching attestation
// verification, so — unlike AttestationTeeAuthenticator.ts — they run on the
// default Hardhat network.
describe('TeeAuthenticator Test (upgradeability)', function () {
  let signers: Signer[];
  let owner: string;
  let nitroProverAddress: string;
  const PCR0 = '0x' + '11'.repeat(48);
  const TEE_MAX_VERIFICATION_AGE = 900_000;

  beforeEach(async function () {
    signers = await ethers.getSigners();
    owner = await signers[0].getAddress();
    // Never called by these tests, so any address works as a stand-in.
    nitroProverAddress = await signers[5].getAddress();
  });

  describe('initialize', function () {
    it('reverts when initialize is called a second time', async () => {
      const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
      const teeAuthenticator = await upgrades.deployProxy(
        TeeAuthenticator,
        [owner, nitroProverAddress, PCR0, TEE_MAX_VERIFICATION_AGE],
        { kind: 'uups' }
      );

      await expect(
        teeAuthenticator.initialize(owner, nitroProverAddress, PCR0, TEE_MAX_VERIFICATION_AGE)
      ).to.be.revertedWithCustomError(teeAuthenticator, 'InvalidInitialization');
    });

    it('reverts when initialize is called directly on the implementation', async () => {
      const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
      const teeAuthenticator = await upgrades.deployProxy(
        TeeAuthenticator,
        [owner, nitroProverAddress, PCR0, TEE_MAX_VERIFICATION_AGE],
        { kind: 'uups' }
      );
      const implementationAddress = await upgrades.erc1967.getImplementationAddress(
        await teeAuthenticator.getAddress()
      );
      const implementation = TeeAuthenticator.attach(implementationAddress);

      await expect(
        (implementation as any).initialize(
          owner,
          nitroProverAddress,
          PCR0,
          TEE_MAX_VERIFICATION_AGE
        )
      ).to.be.revertedWithCustomError(teeAuthenticator, 'InvalidInitialization');
    });

    it('initializes owner, pcr0, nitroProver, and maxVerificationAge', async () => {
      const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
      const teeAuthenticator = await upgrades.deployProxy(
        TeeAuthenticator,
        [owner, nitroProverAddress, PCR0, TEE_MAX_VERIFICATION_AGE],
        { kind: 'uups' }
      );

      expect(await teeAuthenticator.owner()).to.equal(owner);
      expect(await teeAuthenticator.pcr0()).to.equal(PCR0);
      expect(await teeAuthenticator.nitroProver()).to.equal(nitroProverAddress);
      expect(await teeAuthenticator.maxVerificationAge()).to.equal(TEE_MAX_VERIFICATION_AGE);
    });
  });

  describe('upgradeability', function () {
    it('reverts when a non-owner attempts to upgrade', async () => {
      const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
      const teeAuthenticator = await upgrades.deployProxy(
        TeeAuthenticator,
        [owner, nitroProverAddress, PCR0, TEE_MAX_VERIFICATION_AGE],
        { kind: 'uups' }
      );
      const v2Implementation = await TeeAuthenticator.deploy();
      await v2Implementation.waitForDeployment();

      await expect(
        (teeAuthenticator.connect(signers[1]) as any).upgradeToAndCall(
          await v2Implementation.getAddress(),
          '0x'
        )
      ).to.be.revertedWithCustomError(teeAuthenticator, 'OwnableUnauthorizedAccount');
    });

    it('preserves storage across an upgrade to a new implementation', async () => {
      const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
      const teeAuthenticator = await upgrades.deployProxy(
        TeeAuthenticator,
        [owner, nitroProverAddress, PCR0, TEE_MAX_VERIFICATION_AGE],
        { kind: 'uups' }
      );
      const proxyAddress = await teeAuthenticator.getAddress();

      await teeAuthenticator.updatePcr0('0x' + '22'.repeat(48));

      const upgraded = await upgrades.upgradeProxy(proxyAddress, TeeAuthenticator);
      await upgraded.waitForDeployment();

      expect(await upgraded.getAddress()).to.equal(proxyAddress);
      expect(await upgraded.owner()).to.equal(owner);
      expect(await upgraded.pcr0()).to.equal('0x' + '22'.repeat(48));
      expect(await upgraded.nitroProver()).to.equal(nitroProverAddress);
      expect(await upgraded.maxVerificationAge()).to.equal(TEE_MAX_VERIFICATION_AGE);
    });
  });
});
