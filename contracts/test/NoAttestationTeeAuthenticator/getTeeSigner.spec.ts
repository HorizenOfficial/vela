import { expect } from 'chai';
import { ADDRESS_ZERO, BYTES_ZERO } from '../util';
import { deployNoAttestationTeeAuthenticatorEmptyFixture, deployNoAttestationTeeAuthenticatorFixture } from './fixture';

describe('NoAttestationTeeAuthenticator Test', function () {
  describe('getTeeSigner', function () {
    describe('unhappy paths', function () {
      it('returns zero address when teeSigner was initialized as zero', async () => {
        const { teeAuthenticator } = await deployNoAttestationTeeAuthenticatorEmptyFixture();

        const teeSigner = await teeAuthenticator.getTeeSigner();
        expect(teeSigner).to.equal(ADDRESS_ZERO);
      });
    });

    describe('happy paths', function () {
      it('returns the configured teeSigner', async () => {
        const { teeAuthenticator, teeSigner: configuredTeeSigner } =
          await deployNoAttestationTeeAuthenticatorFixture({
          pubKey: BYTES_ZERO,
        });

        const teeSigner = await teeAuthenticator.getTeeSigner();
        expect(teeSigner).to.equal(configuredTeeSigner);
      });
    });
  });
});
