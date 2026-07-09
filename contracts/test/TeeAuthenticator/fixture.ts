import { ethers } from 'hardhat';
import { getRandomHexString } from '../util';

export const PCR0_LENGTH = 48;
export const PCR0_UPGRADE_DELAY = 86400; // 1 day
export const PCR0_SWAP_APPLY_WINDOW = 7 * 86400; // 7 days
export const TEE_MAX_VERIFICATION_AGE = '900000000000000000000000000000000';

// PCR0 sits at offset 4 in the raw PCRs blob; surround it with arbitrary bytes.
export function buildRawPcrs(pcr0: string): string {
  return '0x00000058' + pcr0.slice(2) + getRandomHexString(96).slice(2);
}

type TeeAuthenticatorFixtureOptions = {
  pcr0?: string;
  upgradeDelay?: number;
};

export async function deployTeeAuthenticatorFixture(options: TeeAuthenticatorFixtureOptions = {}) {
  const signers = await ethers.getSigners();
  const owner = await signers[0].getAddress();
  const initialPcr0 = options.pcr0 ?? getRandomHexString(PCR0_LENGTH);
  const upgradeDelay = options.upgradeDelay ?? PCR0_UPGRADE_DELAY;
  const enclaveKey = getRandomHexString(133);
  const teeSignerAddress = await signers[1].getAddress();

  const MockNitroProver = await ethers.getContractFactory('MockNitroProver');
  const nitroProver = await MockNitroProver.deploy(
    enclaveKey,
    teeSignerAddress,
    buildRawPcrs(initialPcr0)
  );
  await nitroProver.deploymentTransaction()!.wait();

  const TeeAuthenticator = await ethers.getContractFactory('TeeAuthenticator');
  const teeAuthenticator = await TeeAuthenticator.deploy(
    owner,
    await nitroProver.getAddress(),
    initialPcr0,
    TEE_MAX_VERIFICATION_AGE,
    upgradeDelay
  );
  await teeAuthenticator.deploymentTransaction()!.wait();

  return {
    signers,
    owner,
    teeAuthenticator,
    nitroProver,
    initialPcr0,
    upgradeDelay,
    enclaveKey,
    teeSignerAddress,
  };
}

export function pcr0Key(pcr0: string): string {
  return ethers.keccak256(pcr0);
}

// Proposes `target` and applies it after the timelock, making it the active image.
export async function swapTo(teeAuthenticator: any, target: string, upgradeDelay: number) {
  const { time } = await import('@nomicfoundation/hardhat-network-helpers');
  await (await teeAuthenticator.proposePcr0Swap(target)).wait();
  await time.increase(upgradeDelay + 1);
  await (await teeAuthenticator.applyPcr0Swap()).wait();
}
