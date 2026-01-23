import { ethers } from 'hardhat';
import { ADDRESS_ZERO, BYTES_ZERO, getRandomHexString } from '../util';

type NoAttestationFixtureOptions = {
  owner?: string;
  teeSigner?: string;
  pubKey?: string;
};

export async function deployNoAttestationTeeAuthenticatorFixture(
  options: NoAttestationFixtureOptions = {}
) {
  const signers = await ethers.getSigners();
  const owner = options.owner ?? await signers[0].getAddress();
  const teeSigner = options.teeSigner ?? await signers[1].getAddress();
  const pubKey = options.pubKey ?? getRandomHexString(133);

  const TeeAuthenticator = await ethers.getContractFactory('NoAttestationTeeAuthenticator');
  const teeAuthenticator = await TeeAuthenticator.deploy(owner, teeSigner, pubKey);
  const pkLength = Number(await teeAuthenticator.PK_LENGTH());

  return {
    signers,
    teeAuthenticator,
    pkLength,
    owner,
    teeSigner,
    pubKey,
  };
}

export async function deployNoAttestationTeeAuthenticatorEmptyFixture() {
  return deployNoAttestationTeeAuthenticatorFixture({
    teeSigner: ADDRESS_ZERO,
    pubKey: BYTES_ZERO,
  });
}
