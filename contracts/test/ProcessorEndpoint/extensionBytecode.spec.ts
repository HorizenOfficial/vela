import { expect } from 'chai';
import hre, { ethers } from 'hardhat';
import {
  assertExtensionBytecodeMatches,
  verifyExtensionBytecode,
} from '../../scripts/upgrade/extensionBytecode';

// Returns the error the promise rejected with, or undefined if it resolved.
async function rejection(promise: Promise<unknown>): Promise<Error | undefined> {
  try {
    await promise;
    return undefined;
  } catch (e) {
    return e as Error;
  }
}

async function deployExtension(): Promise<string> {
  const factory = await ethers.getContractFactory('ProcessorEndpointExtension');
  const extension = await factory.deploy();
  await extension.waitForDeployment();
  return extension.getAddress();
}

describe('assertExtensionBytecodeMatches', function () {
  it('accepts an extension deployed from the compiled source', async function () {
    expect(await rejection(assertExtensionBytecodeMatches(await deployExtension()))).to.equal(
      undefined
    );
  });

  it('rejects an address with no code', async function () {
    const error = await rejection(
      assertExtensionBytecodeMatches('0x000000000000000000000000000000000000dEaD')
    );
    expect(error?.message).to.match(/no contract code/);
  });

  it('rejects an address holding a different contract', async function () {
    const admin = await (await ethers.getSigners())[0].getAddress();
    const tokenAllowlist = await (await ethers.getContractFactory('TokenAllowlist')).deploy(admin);
    await tokenAllowlist.waitForDeployment();

    const error = await rejection(
      assertExtensionBytecodeMatches(await tokenAllowlist.getAddress())
    );
    expect(error?.message).to.match(/different length/);
  });

  // The realistic staleness case: the extension source changed, so the deployed code differs from
  // the local one, but only in bytes outside the `_self` immutable holes. A comparison that masked
  // those holes instead of filling them would still have to catch this; a comparison that skipped
  // the splice altogether would report every extension as mismatched.
  it('rejects same-length code differing outside the immutable holes', async function () {
    const address = await deployExtension();
    const deployed = await ethers.provider.getCode(address);
    // Offset 0 of the runtime code: never inside an immutable hole (the first is at byte 776).
    const tampered = '0x' + (deployed[2] === 'a' ? 'b' : 'a') + deployed.slice(3);
    await hre.network.provider.send('hardhat_setCode', [address, tampered]);

    const error = await rejection(assertExtensionBytecodeMatches(address));
    expect(error?.message).to.match(/same length/);
  });

  // The splice is what proves the deployed extension's `_self` holds the address the new
  // implementation will delegate to. Masking the holes would let a doctored `_self` through, and a
  // wrong `_self` makes every delegated entry point revert on `onlyDelegateCall`.
  it('rejects code whose _self immutable holds the wrong address', async function () {
    const address = await deployExtension();
    const other = await deployExtension();
    // `other`'s code is byte-identical except for the _self holes, which hold `other`'s address.
    await hre.network.provider.send('hardhat_setCode', [
      address,
      await ethers.provider.getCode(other),
    ]);

    const error = await rejection(assertExtensionBytecodeMatches(address));
    expect(error?.message).to.match(/same length/);
  });

  describe('verifyExtensionBytecode', function () {
    it('refuses a mismatched extension when stale ones are not allowed', async function () {
      const admin = await (await ethers.getSigners())[0].getAddress();
      const other = await (await ethers.getContractFactory('TokenAllowlist')).deploy(admin);
      await other.waitForDeployment();

      const error = await rejection(verifyExtensionBytecode(await other.getAddress(), false));
      expect(error?.message).to.match(/not the locally compiled one/);
    });

    it('warns instead of refusing when stale ones are allowed', async function () {
      const admin = await (await ethers.getSigners())[0].getAddress();
      const other = await (await ethers.getContractFactory('TokenAllowlist')).deploy(admin);
      await other.waitForDeployment();

      const warnings: string[] = [];
      const warn = console.warn;
      console.warn = (message?: unknown) => warnings.push(String(message));
      try {
        expect(await rejection(verifyExtensionBytecode(await other.getAddress(), true))).to.equal(
          undefined
        );
      } finally {
        console.warn = warn;
      }
      expect(warnings.join('\n')).to.match(/not the locally compiled one/);
    });

    it('accepts a matching extension', async function () {
      expect(await rejection(verifyExtensionBytecode(await deployExtension(), false))).to.equal(
        undefined
      );
    });
  });
});
