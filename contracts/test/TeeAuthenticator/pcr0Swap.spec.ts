import { expect } from 'chai';
import { ethers } from 'hardhat';
import { time } from '@nomicfoundation/hardhat-network-helpers';
import { getRandomHexString } from '../util';
import {
  deployTeeAuthenticatorFixture,
  pcr0Key,
  swapTo,
  PCR0_LENGTH,
  PCR0_UPGRADE_DELAY,
  PCR0_SWAP_APPLY_WINDOW,
} from './fixture';

describe('TeeAuthenticator PCR0 swap flow', function () {
  describe('constructor', function () {
    it('seeds the accepted set, activeImage and pcr0UpgradeDelay', async () => {
      const { teeAuthenticator, initialPcr0, upgradeDelay } = await deployTeeAuthenticatorFixture();

      const key = pcr0Key(initialPcr0);
      expect(await teeAuthenticator.acceptedPcr0(key)).to.equal(true);
      expect(await teeAuthenticator.acceptedPcr0List(0)).to.equal(key);
      expect(await teeAuthenticator.getAcceptedPcr0Count()).to.equal(1);
      expect(await teeAuthenticator.activeImage()).to.equal(key);
      expect(await teeAuthenticator.pcr0UpgradeDelay()).to.equal(upgradeDelay);
    });

    it('reverts when the initial PCR0 is not 48 bytes', async () => {
      await expect(deployTeeAuthenticatorFixture({ pcr0: getRandomHexString(PCR0_LENGTH - 1) })).to
        .be.reverted;
    });
  });

  describe('proposePcr0Swap', function () {
    it('reverts when caller is not owner', async () => {
      const { teeAuthenticator, signers } = await deployTeeAuthenticatorFixture();

      await expect(
        teeAuthenticator.connect(signers[1]).proposePcr0Swap(getRandomHexString(PCR0_LENGTH))
      ).to.be.revertedWithCustomError(teeAuthenticator, 'OwnableUnauthorizedAccount');
    });

    it('reverts with InvalidPcr0Length when the target is not 48 bytes', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await expect(
        teeAuthenticator.proposePcr0Swap(getRandomHexString(PCR0_LENGTH + 1))
      ).to.be.revertedWithCustomError(teeAuthenticator, 'InvalidPcr0Length');
    });

    it('stores the pending swap and emits Pcr0SwapProposed with the raw preimage', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);

      const tx = await teeAuthenticator.proposePcr0Swap(target);
      const receipt = await tx.wait();
      const eta = BigInt((await ethers.provider.getBlock(receipt!.blockNumber))!.timestamp);

      await expect(tx)
        .to.emit(teeAuthenticator, 'Pcr0SwapProposed')
        .withArgs(target, eta + BigInt(PCR0_UPGRADE_DELAY));

      const pending = await teeAuthenticator.pendingSwap();
      expect(pending.value).to.equal(target);
      expect(pending.eta).to.equal(eta + BigInt(PCR0_UPGRADE_DELAY));
      expect(pending.pending).to.equal(true);
    });

    it('reverts with SwapAlreadyPending while a live proposal exists', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await (await teeAuthenticator.proposePcr0Swap(getRandomHexString(PCR0_LENGTH))).wait();

      await expect(
        teeAuthenticator.proposePcr0Swap(getRandomHexString(PCR0_LENGTH))
      ).to.be.revertedWithCustomError(teeAuthenticator, 'SwapAlreadyPending');
    });

    it('allows a new proposal after cancelling the previous one', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
      const second = getRandomHexString(PCR0_LENGTH);

      await (await teeAuthenticator.proposePcr0Swap(getRandomHexString(PCR0_LENGTH))).wait();
      await (await teeAuthenticator.cancelPcr0Swap()).wait();
      await (await teeAuthenticator.proposePcr0Swap(second)).wait();

      const pending = await teeAuthenticator.pendingSwap();
      expect(pending.value).to.equal(second);
    });

    it('allows a new proposal once the previous one has expired', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
      const second = getRandomHexString(PCR0_LENGTH);

      await (await teeAuthenticator.proposePcr0Swap(getRandomHexString(PCR0_LENGTH))).wait();
      await time.increase(PCR0_UPGRADE_DELAY + PCR0_SWAP_APPLY_WINDOW + 1);
      await (await teeAuthenticator.proposePcr0Swap(second)).wait();

      const pending = await teeAuthenticator.pendingSwap();
      expect(pending.value).to.equal(second);
    });
  });

  describe('applyPcr0Swap', function () {
    it('reverts when caller is not owner', async () => {
      const { teeAuthenticator, signers } = await deployTeeAuthenticatorFixture();

      await expect(
        teeAuthenticator.connect(signers[1]).applyPcr0Swap()
      ).to.be.revertedWithCustomError(teeAuthenticator, 'OwnableUnauthorizedAccount');
    });

    it('reverts with NoPendingSwap when nothing is pending', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await expect(teeAuthenticator.applyPcr0Swap()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'NoPendingSwap'
      );
    });

    it('reverts with TimelockNotElapsed for a new PCR0 before eta', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await (await teeAuthenticator.proposePcr0Swap(getRandomHexString(PCR0_LENGTH))).wait();

      await expect(teeAuthenticator.applyPcr0Swap()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'TimelockNotElapsed'
      );
    });

    it('after eta adds the PCR0 to the set, points activeImage at it and emits Pcr0Swapped', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);

      await (await teeAuthenticator.proposePcr0Swap(target)).wait();
      await time.increase(PCR0_UPGRADE_DELAY + 1);

      await expect(teeAuthenticator.applyPcr0Swap())
        .to.emit(teeAuthenticator, 'Pcr0Swapped')
        .withArgs(target);

      const key = pcr0Key(target);
      expect(await teeAuthenticator.acceptedPcr0(key)).to.equal(true);
      expect(await teeAuthenticator.getAcceptedPcr0Count()).to.equal(2);
      expect(await teeAuthenticator.activeImage()).to.equal(key);

      const pending = await teeAuthenticator.pendingSwap();
      expect(pending.pending).to.equal(false);
    });

    it('applies immediately for an already-accepted PCR0 (rollback) without duplicating it', async () => {
      const { teeAuthenticator, initialPcr0 } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);
      await swapTo(teeAuthenticator, target, PCR0_UPGRADE_DELAY);

      // Rollback to the initial image: no timelock wait.
      await (await teeAuthenticator.proposePcr0Swap(initialPcr0)).wait();
      await expect(teeAuthenticator.applyPcr0Swap())
        .to.emit(teeAuthenticator, 'Pcr0Swapped')
        .withArgs(initialPcr0);

      expect(await teeAuthenticator.activeImage()).to.equal(pcr0Key(initialPcr0));
      expect(await teeAuthenticator.getAcceptedPcr0Count()).to.equal(2);
    });

    it('rolls back immediately to any accepted member of a 3-image set, not just the previous one', async () => {
      const { teeAuthenticator, initialPcr0 } = await deployTeeAuthenticatorFixture();
      await swapTo(teeAuthenticator, getRandomHexString(PCR0_LENGTH), PCR0_UPGRADE_DELAY);
      await swapTo(teeAuthenticator, getRandomHexString(PCR0_LENGTH), PCR0_UPGRADE_DELAY);

      // Active is the third image; roll back to the oldest without waiting for the timelock.
      await (await teeAuthenticator.proposePcr0Swap(initialPcr0)).wait();
      await expect(teeAuthenticator.applyPcr0Swap())
        .to.emit(teeAuthenticator, 'Pcr0Swapped')
        .withArgs(initialPcr0);

      expect(await teeAuthenticator.activeImage()).to.equal(pcr0Key(initialPcr0));
      expect(await teeAuthenticator.getAcceptedPcr0Count()).to.equal(3);
    });

    it('reverts with SwapProposalExpired after the application window', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await (await teeAuthenticator.proposePcr0Swap(getRandomHexString(PCR0_LENGTH))).wait();
      await time.increase(PCR0_UPGRADE_DELAY + PCR0_SWAP_APPLY_WINDOW + 1);

      await expect(teeAuthenticator.applyPcr0Swap()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'SwapProposalExpired'
      );
    });
  });

  describe('cancelPcr0Swap', function () {
    it('reverts when caller is not owner', async () => {
      const { teeAuthenticator, signers } = await deployTeeAuthenticatorFixture();

      await expect(
        teeAuthenticator.connect(signers[1]).cancelPcr0Swap()
      ).to.be.revertedWithCustomError(teeAuthenticator, 'OwnableUnauthorizedAccount');
    });

    it('reverts with NoPendingSwap when nothing is pending', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await expect(teeAuthenticator.cancelPcr0Swap()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'NoPendingSwap'
      );
    });

    it('clears the pending swap and emits Pcr0SwapCancelled', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);

      await (await teeAuthenticator.proposePcr0Swap(target)).wait();
      await expect(teeAuthenticator.cancelPcr0Swap())
        .to.emit(teeAuthenticator, 'Pcr0SwapCancelled')
        .withArgs(target);

      await expect(teeAuthenticator.applyPcr0Swap()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'NoPendingSwap'
      );
    });

    it('cancels an expired proposal', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);

      await (await teeAuthenticator.proposePcr0Swap(target)).wait();
      await time.increase(PCR0_UPGRADE_DELAY + PCR0_SWAP_APPLY_WINDOW + 1);

      await expect(teeAuthenticator.cancelPcr0Swap())
        .to.emit(teeAuthenticator, 'Pcr0SwapCancelled')
        .withArgs(target);
    });
  });

  describe('removePcr0', function () {
    it('reverts when caller is not owner', async () => {
      const { teeAuthenticator, signers, initialPcr0 } = await deployTeeAuthenticatorFixture();

      await expect(
        teeAuthenticator.connect(signers[1]).removePcr0(initialPcr0)
      ).to.be.revertedWithCustomError(teeAuthenticator, 'OwnableUnauthorizedAccount');
    });

    it('reverts with UnknownPcr0 for a non-member', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();

      await expect(
        teeAuthenticator.removePcr0(getRandomHexString(PCR0_LENGTH))
      ).to.be.revertedWithCustomError(teeAuthenticator, 'UnknownPcr0');
    });

    it('reverts with CannotRemoveActiveImage on the last entry (always the active image)', async () => {
      const { teeAuthenticator, initialPcr0 } = await deployTeeAuthenticatorFixture();

      await expect(teeAuthenticator.removePcr0(initialPcr0)).to.be.revertedWithCustomError(
        teeAuthenticator,
        'CannotRemoveActiveImage'
      );
    });

    it('reverts with CannotRemoveActiveImage on the active image', async () => {
      const { teeAuthenticator } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);
      await swapTo(teeAuthenticator, target, PCR0_UPGRADE_DELAY);

      await expect(teeAuthenticator.removePcr0(target)).to.be.revertedWithCustomError(
        teeAuthenticator,
        'CannotRemoveActiveImage'
      );
    });

    it('removes a non-active member, keeps the list consistent and emits Pcr0Removed', async () => {
      const { teeAuthenticator, initialPcr0 } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);
      await swapTo(teeAuthenticator, target, PCR0_UPGRADE_DELAY);

      await expect(teeAuthenticator.removePcr0(initialPcr0))
        .to.emit(teeAuthenticator, 'Pcr0Removed')
        .withArgs(initialPcr0);

      expect(await teeAuthenticator.acceptedPcr0(pcr0Key(initialPcr0))).to.equal(false);
      expect(await teeAuthenticator.getAcceptedPcr0Count()).to.equal(1);
      expect(await teeAuthenticator.acceptedPcr0List(0)).to.equal(pcr0Key(target));
    });

    it('removes a middle entry from a 3-image set and keeps mapping and list consistent', async () => {
      const { teeAuthenticator, initialPcr0 } = await deployTeeAuthenticatorFixture();
      const second = getRandomHexString(PCR0_LENGTH);
      const third = getRandomHexString(PCR0_LENGTH);
      await swapTo(teeAuthenticator, second, PCR0_UPGRADE_DELAY);
      await swapTo(teeAuthenticator, third, PCR0_UPGRADE_DELAY);

      // List is [initial, second, third], active is third: remove the middle entry.
      await (await teeAuthenticator.removePcr0(second)).wait();

      expect(await teeAuthenticator.acceptedPcr0(pcr0Key(second))).to.equal(false);
      expect(await teeAuthenticator.acceptedPcr0(pcr0Key(initialPcr0))).to.equal(true);
      expect(await teeAuthenticator.acceptedPcr0(pcr0Key(third))).to.equal(true);
      expect(await teeAuthenticator.getAcceptedPcr0Count()).to.equal(2);
      // Swap-remove: the last entry takes the removed slot.
      expect(await teeAuthenticator.acceptedPcr0List(0)).to.equal(pcr0Key(initialPcr0));
      expect(await teeAuthenticator.acceptedPcr0List(1)).to.equal(pcr0Key(third));

      // The removed image needs a fresh timelocked proposal to come back.
      await (await teeAuthenticator.proposePcr0Swap(second)).wait();
      await expect(teeAuthenticator.applyPcr0Swap()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'TimelockNotElapsed'
      );
    });

    it('a removed PCR0 needs a fresh timelocked proposal to come back', async () => {
      const { teeAuthenticator, initialPcr0 } = await deployTeeAuthenticatorFixture();
      const target = getRandomHexString(PCR0_LENGTH);
      await swapTo(teeAuthenticator, target, PCR0_UPGRADE_DELAY);
      await (await teeAuthenticator.removePcr0(initialPcr0)).wait();

      await (await teeAuthenticator.proposePcr0Swap(initialPcr0)).wait();
      await expect(teeAuthenticator.applyPcr0Swap()).to.be.revertedWithCustomError(
        teeAuthenticator,
        'TimelockNotElapsed'
      );
    });
  });
});
