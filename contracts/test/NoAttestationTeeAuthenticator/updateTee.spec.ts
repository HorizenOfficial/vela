import { expect } from 'chai';
import { ADDRESS_ZERO, BYTES_ZERO, getRandomHexString } from '../util';
import { deployNoAttestationTeeAuthenticatorFixture } from './fixture';

describe('NoAttestationTeeAuthenticator Test', function () {
  describe('updateTee', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {
        const { teeAuthenticator, signers, pkLength } = await deployNoAttestationTeeAuthenticatorFixture();

        await expect(
          teeAuthenticator.connect(signers[1]).updateTee(
            await signers[1].getAddress(),
            getRandomHexString(pkLength)
          )
        ).to.be.revertedWithCustomError(teeAuthenticator, 'OwnableUnauthorizedAccount');
      });

      it('reverts with TeeAddressCantBeZero when newTeeSigner is zero', async () => {
        const { teeAuthenticator, pkLength } = await deployNoAttestationTeeAuthenticatorFixture();

        await expect(
          teeAuthenticator.updateTee(ADDRESS_ZERO, getRandomHexString(pkLength))
        ).to.be.revertedWithCustomError(teeAuthenticator, 'TeeAddressCantBeZero');
      });

      it('reverts with InvalidPKLength when pubSecp521r1 length is invalid', async () => {
        const { teeAuthenticator, signers, pkLength } = await deployNoAttestationTeeAuthenticatorFixture();

        await expect(
          teeAuthenticator.updateTee(
            await signers[1].getAddress(),
            getRandomHexString(pkLength + 1)
          )
        ).to.be.revertedWithCustomError(teeAuthenticator, 'InvalidPKLength');
      });
    });

    describe('happy paths', function () {
      it('updates teeSigner and pubSecp521r1 and emits TeeUpdate', async () => {
        const { teeAuthenticator, signers, pkLength, teeSigner, pubKey } =
          await deployNoAttestationTeeAuthenticatorFixture();

        const newTeeSigner = await signers[2].getAddress();
        const newPubKey = getRandomHexString(pkLength);

        await expect(
          teeAuthenticator.updateTee(newTeeSigner, newPubKey)
        ).to.emit(teeAuthenticator, 'TeeUpdate').withArgs(
          teeSigner,
          newTeeSigner,
          pubKey,
          newPubKey
        );

        expect(await teeAuthenticator.teeSigner()).to.equal(newTeeSigner);
        expect(await teeAuthenticator.pubSecp521r1()).to.equal(newPubKey);
      });
    });
  });
});
