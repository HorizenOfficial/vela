import hre, { ethers } from 'hardhat';

// Verifies that the `ProcessorEndpointExtension` already deployed at a given address is the same
// code as the locally compiled one.
//
// Why this needs a check of its own: the extension is a separately deployed contract that
// `ProcessorEndpoint` reaches by `delegatecall`, and its address is an immutable constructor
// argument of the endpoint implementation. Nothing else catches a stale one. The endpoint's
// constructor only rejects `address(0)` and an address with no code; `upgrades.upgradeProxy`
// validates the implementation's storage layout and never looks at `constructorArgs`; and
// `npm run check:layout` compares two compiled artifacts without touching the chain. So an upgrade
// that changes `ProcessorEndpointExtension.sol` but reuses the deployed extension address produces
// a proxy whose new endpoint logic is live while every delegated entry point (`submitRequestFor`,
// deploy submission, `initialize`, `invokeTrigger`, the operator resets, the admin setters) still
// runs the old code against the same storage — no revert, nothing to notice.

const SOURCE_NAME = 'contracts/ProcessorEndpointExtension.sol';
const CONTRACT_NAME = 'ProcessorEndpointExtension';

// `ProcessorEndpointExtension`'s only `immutable`. Its value is the extension's own address, which
// is what makes the splice below possible: a second immutable would hold something this check
// cannot predict, so it aborts instead of silently comparing fewer bytes.
const SELF_IMMUTABLE = '_self';

type ImmutableReferences = Record<string, { start: number; length: number }[]>;

type CompiledExtension = {
  // Runtime bytecode as compiled: the `immutable` slots are 32-byte holes of zeros, filled in by
  // the constructor at deployment.
  local: Buffer;
  // Byte offsets of those holes, from the compiler.
  holes: { start: number; length: number }[];
};

// Walks the solc AST collecting variable-declaration names by node id, so an entry in
// `immutableReferences` (keyed by declaration id) can be reported and checked by name.
function variableNames(node: unknown, names: Map<number, string>): Map<number, string> {
  if (node === null || typeof node !== 'object') return names;
  const record = node as Record<string, unknown>;
  if (record.nodeType === 'VariableDeclaration' && typeof record.id === 'number') {
    names.set(record.id, String(record.name));
  }
  for (const value of Object.values(record)) {
    if (Array.isArray(value)) value.forEach((entry) => variableNames(entry, names));
    else variableNames(value, names);
  }
  return names;
}

async function compiledExtension(): Promise<CompiledExtension> {
  const buildInfo = await hre.artifacts.getBuildInfo(`${SOURCE_NAME}:${CONTRACT_NAME}`);
  if (!buildInfo) {
    throw new Error(`no build info for ${CONTRACT_NAME}; run "npm run build" first`);
  }

  const deployedBytecode = (
    buildInfo.output.contracts[SOURCE_NAME][CONTRACT_NAME] as {
      evm: { deployedBytecode: { object: string; immutableReferences?: ImmutableReferences } };
    }
  ).evm.deployedBytecode;

  // `immutableReferences` lives only in the build-info output, not in the contract artifact.
  const references = deployedBytecode.immutableReferences ?? {};
  const names = variableNames(
    (buildInfo.output.sources[SOURCE_NAME] as { ast: unknown }).ast,
    new Map()
  );
  const unexpected = Object.keys(references)
    .map((id) => names.get(Number(id)) ?? `declaration ${id}`)
    .filter((name) => name !== SELF_IMMUTABLE);
  if (unexpected.length > 0) {
    throw new Error(
      `${CONTRACT_NAME} declares immutable(s) this check does not know the expected value of: ` +
        `${unexpected.join(', ')}. It fills every immutable slot with the extension's own address ` +
        `(true for ${SELF_IMMUTABLE} alone) before comparing against the deployed code; teach it ` +
        `the new immutable's value, or the comparison would have to ignore those bytes.`
    );
  }

  return {
    local: Buffer.from(deployedBytecode.object.replace(/^0x/, ''), 'hex'),
    holes: Object.values(references).flat(),
  };
}

/// Throws unless the code deployed at `extensionAddress` is the locally compiled
/// `ProcessorEndpointExtension`. Comparing raw bytes would fail on every correct upgrade: the
/// compiled runtime bytecode leaves each `immutable` reference as 32 zero bytes, and
/// `_self = address(this)` is referenced enough times to differ from the deployed code by
/// hundreds of bytes. So the extension's own address is written into those holes first, which also
/// proves the deployed `_self` holds the address the new implementation is about to delegate to.
export async function assertExtensionBytecodeMatches(extensionAddress: string): Promise<void> {
  const onChainHex = await ethers.provider.getCode(extensionAddress);
  if (onChainHex === '0x') {
    throw new Error(`no contract code at ${extensionAddress}; expected ${CONTRACT_NAME}`);
  }
  const onChain = Buffer.from(onChainHex.replace(/^0x/, ''), 'hex');

  const { local, holes } = await compiledExtension();
  const expected = Buffer.from(local);
  const selfWord = Buffer.from(ethers.zeroPadValue(extensionAddress, 32).replace(/^0x/, ''), 'hex');
  for (const hole of holes) {
    selfWord.copy(expected, hole.start, selfWord.length - hole.length);
  }

  if (expected.equals(onChain)) return;

  // Length tells the operator which of the two situations they are in, which is the difference
  // between "deploy the new extension" and "check your compiler settings".
  const detail =
    expected.length !== onChain.length
      ? `different length: ${onChain.length} bytes on chain, ${expected.length} locally, so the ` +
        `extension source itself differs`
      : `same length, ${onChain.filter((byte, i) => byte !== expected[i]).length} differing ` +
        `bytes, so it is either a different revision of the same source or the same source built ` +
        `with different compiler settings`;

  throw new Error(
    `the ${CONTRACT_NAME} deployed at ${extensionAddress} is not the locally compiled one ` +
      `(${detail}).\nUpgrading with it would point the new ProcessorEndpoint implementation at ` +
      `the old extension code: the delegated entry points would keep running it against the ` +
      `proxy's storage, without reverting.\nDeploy the extension from this source and pass its ` +
      `address in PROCESSOR_ENDPOINT_EXTENSION, or set ALLOW_STALE_EXTENSION=true to upgrade ` +
      `against the deployed one deliberately.`
  );
}

/// `assertExtensionBytecodeMatches`, downgraded to a warning when `allowStale` is set. The
/// override exists because a solc or optimiser-settings change alters the extension's bytecode
/// even when its source is untouched, so the check can fire on an extension that is functionally
/// the one intended; without an escape hatch that would make the upgrade unrunnable.
export async function verifyExtensionBytecode(
  extensionAddress: string,
  allowStale: boolean
): Promise<void> {
  try {
    await assertExtensionBytecodeMatches(extensionAddress);
    console.log('extension bytecode matches the local build');
  } catch (error) {
    if (!allowStale) throw error;
    console.warn(`WARNING: ${(error as Error).message}`);
    console.warn('continuing because ALLOW_STALE_EXTENSION=true');
  }
}
