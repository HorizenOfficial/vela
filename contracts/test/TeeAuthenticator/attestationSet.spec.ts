import { expect } from 'chai';
import { getRandomHexString } from '../util';
import {
  deployTeeAuthenticatorFixture,
  swapTo,
  buildRawPcrs,
  PCR0_LENGTH,
  PCR0_UPGRADE_DELAY,
} from './fixture';

// _checkAttestationContent checks against activeImage (D3: membership in the accepted set is
// not enough), exercised through updateTee / the step flow with a MockNitroProver (the real
// prover needs Anvil and a genuinely signed attestation).
describe('TeeAuthenticator attestation checks against activeImage', function () {
  it('accepts an attestation for the active image', async () => {
    const { teeAuthenticator, teeSignerAddress } = await deployTeeAuthenticatorFixture();

    await (await teeAuthenticator.updateTee(getRandomHexString(64))).wait();

    expect(await teeAuthenticator.getTeeSigner()).to.equal(teeSignerAddress);
  });

  it('accepts an attestation for the active image after a rollback', async () => {
    const { teeAuthenticator, initialPcr0, teeSignerAddress } =
      await deployTeeAuthenticatorFixture();
    await swapTo(teeAuthenticator, getRandomHexString(PCR0_LENGTH), PCR0_UPGRADE_DELAY);

    // Rollback: the initial image is accepted again as activeImage without a timelock.
    await (await teeAuthenticator.proposePcr0Swap(initialPcr0)).wait();
    await (await teeAuthenticator.applyPcr0Swap()).wait();

    await (await teeAuthenticator.updateTee(getRandomHexString(64))).wait();

    expect(await teeAuthenticator.getTeeSigner()).to.equal(teeSignerAddress);
  });

  it('rejects an attestation from an accepted but non-active set member', async () => {
    const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
    // The mock prover still reports the initial PCR0; make another image active.
    await swapTo(teeAuthenticator, getRandomHexString(PCR0_LENGTH), PCR0_UPGRADE_DELAY);

    await expect(teeAuthenticator.updateTee(getRandomHexString(64))).to.be.revertedWithCustomError(
      teeAuthenticator,
      'InvalidPCR'
    );
  });

  it('rejects an attestation from a non-active member of a 3-image set', async () => {
    const { teeAuthenticator, nitroProver } = await deployTeeAuthenticatorFixture();
    const second = getRandomHexString(PCR0_LENGTH);
    await swapTo(teeAuthenticator, second, PCR0_UPGRADE_DELAY);
    await swapTo(teeAuthenticator, getRandomHexString(PCR0_LENGTH), PCR0_UPGRADE_DELAY);

    // The middle image is accepted but not active: the check must not degrade into membership.
    await (await nitroProver.setRawPcrs(buildRawPcrs(second))).wait();

    await expect(teeAuthenticator.updateTee(getRandomHexString(64))).to.be.revertedWithCustomError(
      teeAuthenticator,
      'InvalidPCR'
    );
  });

  it('rejects an attestation whose PCR0 is not in the set', async () => {
    const { teeAuthenticator, nitroProver } = await deployTeeAuthenticatorFixture();
    await (await nitroProver.setRawPcrs(buildRawPcrs(getRandomHexString(PCR0_LENGTH)))).wait();

    await expect(teeAuthenticator.updateTee(getRandomHexString(64))).to.be.revertedWithCustomError(
      teeAuthenticator,
      'InvalidPCR'
    );
  });

  it('rejects a truncated PCRs blob', async () => {
    const { teeAuthenticator, nitroProver } = await deployTeeAuthenticatorFixture();
    await (await nitroProver.setRawPcrs(getRandomHexString(20))).wait();

    await expect(teeAuthenticator.updateTee(getRandomHexString(64))).to.be.revertedWithCustomError(
      teeAuthenticator,
      'InvalidPCR'
    );
  });

  describe('step flow', function () {
    it('completes when activeImage is unchanged', async () => {
      const { teeAuthenticator, teeSignerAddress } = await deployTeeAuthenticatorFixture();

      await (await teeAuthenticator.updateTeeStep1(getRandomHexString(64))).wait();
      const step2TxCount = await teeAuthenticator.getStep2TotalLength();
      for (let i = 0; i < step2TxCount; i++) {
        await (await teeAuthenticator.updateTeeStep2()).wait();
      }
      await (await teeAuthenticator.updateTeeStep3()).wait();
      await (await teeAuthenticator.updateTeeStep4()).wait();

      expect(await teeAuthenticator.getTeeSigner()).to.equal(teeSignerAddress);
    });

    it('reverts at step4 when activeImage changed after step1 (straddle)', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await (await teeAuthenticator.updateTeeStep1(getRandomHexString(64))).wait();

      // Between step1 and step4 another image becomes active.
      await swapTo(teeAuthenticator, getRandomHexString(PCR0_LENGTH), PCR0_UPGRADE_DELAY);

      const step2TxCount = await teeAuthenticator.getStep2TotalLength();
      for (let i = 0; i < step2TxCount; i++) {
        await (await teeAuthenticator.updateTeeStep2()).wait();
      }
      await (await teeAuthenticator.updateTeeStep3()).wait();

      await expect(teeAuthenticator.updateTeeStep4()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'InvalidPCR'
      );
    });
  });
});
