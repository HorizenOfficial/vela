import { expect } from 'chai';
import { ethers } from 'ethers';
import { ETH_TOKEN, BYTES32_ZERO, getRandomHexString } from '../util';
import {
  ethSignBatchStateUpdate,
  ethSignStateUpdate,
  stateUpdateEntryHash,
} from '../../scripts/util';
import {
  deployNoAttestationTeeAuthenticatorEmptyFixture,
  deployNoAttestationTeeAuthenticatorFixture,
} from './fixture';

describe('NoAttestationTeeAuthenticator Test', function () {
  describe('checkBatchSignature', function () {
    const REQUEST_ID = '0x' + '22'.repeat(32);

    // Builds entry i of a chained batch: prevStateRoot_i == newStateRoot_(i-1).
    function buildPayload(index: number, addr1: string) {
      return {
        applicationId: 7,
        prevStateRoot:
          index === 0 ? BYTES32_ZERO : '0x' + (index + 10).toString(16).padStart(2, '0').repeat(32),
        newStateRoot: '0x' + (index + 11).toString(16).padStart(2, '0').repeat(32),
        processedRequestId: ethers.zeroPadValue(ethers.toBeHex(index + 1), 32),
        userEvents: { events: ['0x01'], subTypes: [ethers.encodeBytes32String('subtype')] },
        appEvents: { events: [], subTypes: [] },
        withdrawalRequests: [[ETH_TOKEN, addr1, 50]],
        refundAmount: 10,
        applicationFee: 20,
        errorCode: 0,
        errorMsg: '',
      };
    }

    function entryHashOf(payload: ReturnType<typeof buildPayload>) {
      return stateUpdateEntryHash(
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
    }

    describe('unhappy paths', function () {
      it('reverts with TeeIsNotSet when teeSigner is zero', async () => {
        const { teeAuthenticator } = await deployNoAttestationTeeAuthenticatorEmptyFixture();

        await expect(
          teeAuthenticator.checkBatchSignature([REQUEST_ID], '0x')
        ).to.be.revertedWithCustomError(teeAuthenticator, 'TeeIsNotSet');
      });

      it('reverts with TeeIsNotSet when pubSecp521r1 length is invalid', async () => {
        const { teeAuthenticator } = await deployNoAttestationTeeAuthenticatorFixture({
          pubKey: getRandomHexString(10),
        });

        await expect(
          teeAuthenticator.checkBatchSignature([REQUEST_ID], '0x')
        ).to.be.revertedWithCustomError(teeAuthenticator, 'TeeIsNotSet');
      });

      it('reverts with EmptyBatch when no entry hashes are provided', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const signature = await ethSignBatchStateUpdate(signers[1], [REQUEST_ID]);

        await expect(
          teeAuthenticator.checkBatchSignature([], signature)
        ).to.be.revertedWithCustomError(teeAuthenticator, 'EmptyBatch');
      });

      it('returns false when the batch is signed by a key other than teeSigner', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const addr = await signers[0].getAddress();
        const entryHashes = [0, 1, 2].map((i) => entryHashOf(buildPayload(i, addr)));

        const signature = await ethSignBatchStateUpdate(signers[2], entryHashes);

        expect(await teeAuthenticator.checkBatchSignature(entryHashes, signature)).to.equal(false);
      });

      it('returns false when an entry hash is altered', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const addr = await signers[0].getAddress();
        const entryHashes = [0, 1, 2].map((i) => entryHashOf(buildPayload(i, addr)));

        const signature = await ethSignBatchStateUpdate(signers[1], entryHashes);

        const tampered = [...entryHashes];
        tampered[1] = '0x' + 'ab'.repeat(32);
        expect(await teeAuthenticator.checkBatchSignature(tampered, signature)).to.equal(false);
      });

      it('returns false when entry hashes are reordered', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const addr = await signers[0].getAddress();
        const entryHashes = [0, 1, 2].map((i) => entryHashOf(buildPayload(i, addr)));

        const signature = await ethSignBatchStateUpdate(signers[1], entryHashes);

        const reordered = [entryHashes[0], entryHashes[2], entryHashes[1]];
        expect(await teeAuthenticator.checkBatchSignature(reordered, signature)).to.equal(false);
      });

      it('returns false when a signature over N entries is replayed with a prefix of them', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const addr = await signers[0].getAddress();
        const entryHashes = [0, 1, 2].map((i) => entryHashOf(buildPayload(i, addr)));

        const signature = await ethSignBatchStateUpdate(signers[1], entryHashes);

        // The personal_sign prefix commits to 32*N, so a shorter batch cannot reuse it.
        expect(
          await teeAuthenticator.checkBatchSignature(entryHashes.slice(0, 2), signature)
        ).to.equal(false);
      });
    });

    describe('happy paths', function () {
      it('returns true for a multi-entry batch signed by teeSigner', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const addr = await signers[0].getAddress();
        const entryHashes = [0, 1, 2, 3, 4].map((i) => entryHashOf(buildPayload(i, addr)));

        const signature = await ethSignBatchStateUpdate(signers[1], entryHashes);

        expect(await teeAuthenticator.checkBatchSignature(entryHashes, signature)).to.equal(true);
      });

      it('accepts a single-entry batch signed with the single-request scheme', async () => {
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture();
        const payload = buildPayload(0, await signers[0].getAddress());

        // Signed via the single-request helper, verified via the batch path: the two
        // digests must be byte-identical for one entry.
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

        expect(
          await teeAuthenticator.checkBatchSignature([entryHashOf(payload)], signature)
        ).to.equal(true);
        expect(await teeAuthenticator.checkSignature(payload, signature)).to.equal(true);
      });

      it('matches an independently computed batch digest', async () => {
        const teeWallet = ethers.Wallet.createRandom();
        const { teeAuthenticator, signers } = await deployNoAttestationTeeAuthenticatorFixture({
          teeSigner: teeWallet.address,
        });
        const addr = await signers[0].getAddress();
        const entryHashes = [0, 1].map((i) => entryHashOf(buildPayload(i, addr)));

        // Digest pinned by construction: keccak256("\x19Ethereum Signed Message:\n64" || h0 || h1)
        const digest = ethers.keccak256(
          ethers.concat([
            ethers.toUtf8Bytes('\x19Ethereum Signed Message:\n64'),
            entryHashes[0],
            entryHashes[1],
          ])
        );
        const signature = teeWallet.signingKey.sign(digest).serialized;

        expect(await teeAuthenticator.checkBatchSignature(entryHashes, signature)).to.equal(true);
      });
    });
  });
});
