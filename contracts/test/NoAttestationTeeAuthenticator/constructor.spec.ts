import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { ADDRESS_ZERO, BYTES_ZERO } from '../util';
import { deployNoAttestationTeeAuthenticatorFixture } from './fixture';

describe('NoAttestationTeeAuthenticator Test', function () {
  describe('constructor', function () {
    let signers: Signer[];

    beforeEach(async function () {
      signers = await ethers.getSigners();
    });

    describe('unhappy paths', function () {
      it('reverts with OwnableInvalidOwner when owner is zero address', async () => {
        const TeeAuthenticator = await ethers.getContractFactory('NoAttestationTeeAuthenticator');

        await expect(
          TeeAuthenticator.deploy(ADDRESS_ZERO, await signers[1].getAddress(), BYTES_ZERO)
        ).to.be.revertedWithCustomError(TeeAuthenticator, 'OwnableInvalidOwner');
      });
    });

    describe('happy paths', function () {
      it('sets teeSigner and pubSecp521r1 and emits TeeUpdate', async () => {
        const { teeAuthenticator, teeSigner, pubKey } =
          await deployNoAttestationTeeAuthenticatorFixture();

        await expect(teeAuthenticator.deploymentTransaction())
          .to.emit(teeAuthenticator, 'TeeUpdate')
          .withArgs(ADDRESS_ZERO, teeSigner, BYTES_ZERO, pubKey);

        expect(await teeAuthenticator.teeSigner()).to.equal(teeSigner);
        expect(await teeAuthenticator.pubSecp521r1()).to.equal(pubKey);
      });
    });
  });
});
