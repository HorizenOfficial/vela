import { expect } from 'chai';
import { BYTES_ZERO, getRandomHexString } from '../util';
import { deployNoAttestationTeeAuthenticatorEmptyFixture, deployNoAttestationTeeAuthenticatorFixture } from './fixture';

describe('NoAttestationTeeAuthenticator Test', function () {
  describe('getPubSecp521r1', function () {
    describe('unhappy paths', function () {
      it('returns empty bytes when pubSecp521r1 was initialized empty', async () => {
        const { teeAuthenticator } = await deployNoAttestationTeeAuthenticatorEmptyFixture();

        const pubKey = await teeAuthenticator.getPubSecp521r1();
        expect(pubKey).to.equal(BYTES_ZERO);
      });
    });

    describe('happy paths', function () {
      it('returns the configured pubSecp521r1', async () => {
        const pubKey = getRandomHexString(133);
        const { teeAuthenticator } = await deployNoAttestationTeeAuthenticatorFixture({ pubKey });

        const stored = await teeAuthenticator.getPubSecp521r1();
        expect(stored).to.equal(pubKey);
      });
    });
  });
});
