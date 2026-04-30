import { expect } from 'chai';
import { ethers } from 'ethers';
import { ETH_TOKEN, BYTES32_ZERO, getRandomHexString } from '../util';
import { ethSignStateUpdate } from '../../scripts/util';
import {
  deployNoAttestationTeeAuthenticatorEmptyFixture,
  deployNoAttestationTeeAuthenticatorFixture,
} from './fixture';

describe('NoAttestationTeeAuthenticator Test', function () {
  describe('checkSignature', function () {
    const NEW_STATE_ROOT = '0x' + '11'.repeat(32);
    const REQUEST_ID = '0x' + '22'.repeat(32);

    function buildPayload(addr1: string, addr2: string) {
      return {
        applicationId: 0,
        prevStateRoot: BYTES32_ZERO,
        newStateRoot: NEW_STATE_ROOT,
        processedRequestId: REQUEST_ID,
        userEvents: { events: ['0x01'], subTypes: [ethers.encodeBytes32String('subtype')] },
        appEvents: {
          events: ['0xca', '0xfe'],
          subTypes: [
            ethers.encodeBytes32String('app_sub1'),
            ethers.encodeBytes32String('app_sub2'),
          ],
        },
        withdrawalRequests: [
          [ETH_TOKEN, addr1, 50],
          [ETH_TOKEN, addr2, 50],
        ],
        refundAmount: 0,
        applicationFee: 0,
        errorCode: 0,
        errorMsg: '',
      };
    }

    describe('unhappy paths', function () {
      it('reverts with TeeIsNotSet when teeSigner is zero', async () => {
        const { teeAuthenticator, signers } =
          await deployNoAttestationTeeAuthenticatorEmptyFixture();
        const payload = buildPayload(await signers[0].getAddress(), await signers[1].getAddress());

        await expect(teeAuthenticator.checkSignature(payload, '0x')).to.be.revertedWithCustomError(
          teeAuthenticator,
          'TeeIsNotSet'
        );
      });

      it('reverts with TeeIsNotSet when pubSecp521r1 length is invalid', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture({
          pubKey: getRandomHexString(10),
        });
        const payload = buildPayload(await signers[0].getAddress(), await signers[1].getAddress());

        await expect(teeAuthenticator.checkSignature(payload, '0x')).to.be.revertedWithCustomError(
          teeAuthenticator,
          'TeeIsNotSet'
        );
      });

      it('returns false when signature does not match teeSigner', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const payload = buildPayload(await signers[0].getAddress(), await signers[2].getAddress());

        const signature = await ethSignStateUpdate(
          signers[2],
          payload.applicationId,
          payload.prevStateRoot,
          payload.newStateRoot,
          payload.processedRequestId,
          payload.userEvents.events,
          payload.userEvents.subTypes,
          payload.appEvents.events,
          payload.appEvents.subTypes,
          payload.withdrawalRequests,
          payload.refundAmount,
          payload.applicationFee,
          payload.errorCode,
          payload.errorMsg
        );

        const res = await teeAuthenticator.checkSignature(payload, signature);
        expect(res).to.equal(false);
      });
    });

    describe('happy paths', function () {
      it('returns true when signature is valid', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const payload = buildPayload(await signers[0].getAddress(), await signers[2].getAddress());

        const signature = await ethSignStateUpdate(
          signers[1],
          payload.applicationId,
          payload.prevStateRoot,
          payload.newStateRoot,
          payload.processedRequestId,
          payload.userEvents.events,
          payload.userEvents.subTypes,
          payload.appEvents.events,
          payload.appEvents.subTypes,
          payload.withdrawalRequests,
          payload.refundAmount,
          payload.applicationFee,
          payload.errorCode,
          payload.errorMsg
        );

        const res = await teeAuthenticator.checkSignature(payload, signature);
        expect(res).to.equal(true);
      });
    });
  });
});
