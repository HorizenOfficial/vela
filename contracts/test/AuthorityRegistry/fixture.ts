import { ethers, upgrades } from 'hardhat';
import { Signer } from 'ethers';

export async function deployDefaultAuthorityFixture() {
  const signers = await ethers.getSigners();
  const DefaultAuthority = await ethers.getContractFactory('DefaultAuthority');
  const defaultAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

  return { signers, defaultAuthority };
}

export async function deployAuthorityRegistryFixture() {
  const { signers, defaultAuthority } = await deployDefaultAuthorityFixture();
  const AuthorityRegistry = await ethers.getContractFactory('AuthorityRegistry');
  const authorityRegistry = await upgrades.deployProxy(
    AuthorityRegistry,
    [await signers[0].getAddress(), await defaultAuthority.getAddress()],
    { kind: 'uups' }
  );

  return { signers, defaultAuthority, authorityRegistry };
}
