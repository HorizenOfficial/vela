import hre from 'hardhat';

// ProcessorEndpoint reaches ProcessorEndpointExtension with delegatecall, so the extension's code
// runs against the endpoint's storage. If their layouts ever diverge, the extension silently reads
// and writes the wrong slots — no revert, corrupted state. Both derive from
// ProcessorEndpointStorage, which is what keeps them in sync; this script is the check that the
// invariant actually holds after every change.
const CONTRACTS = [
  { sourceName: 'contracts/ProcessorEndpoint.sol', contractName: 'ProcessorEndpoint' },
  {
    sourceName: 'contracts/ProcessorEndpointExtension.sol',
    contractName: 'ProcessorEndpointExtension',
  },
];

type Slot = { slot: string; offset: number; label: string; type: string };

// Compares the flattened slot list only: slot, offset, label and type identifier. That covers
// every way the two contracts can realistically diverge, because all state is declared once in
// ProcessorEndpointStorage and the failure mode this guards against is a variable declared (or
// reordered) in a derived contract. It would not catch a change *inside* a struct's own fields —
// e.g. reordering RequestQueue — which is only reachable if a derived contract redeclares a struct
// of the same name, since both contracts otherwise share the single declaration in the base.
function format(entry: Slot): string {
  // Type identifiers embed solc's internal AST ids ("t_contract(ITrigger)9681"), which are only
  // stable within a single compilation job — hardhat may compile the two contracts in separate
  // jobs. Strip the ids so the comparison sees the type shape, not the numbering.
  const type = entry.type.replace(/\)\d+/g, ')');
  return `${entry.slot}:${entry.offset} ${entry.label} ${type}`;
}

async function layoutOf(sourceName: string, contractName: string): Promise<string[]> {
  const buildInfo = await hre.artifacts.getBuildInfo(`${sourceName}:${contractName}`);
  if (!buildInfo) {
    throw new Error(`no build info for ${contractName}; run "npm run build" first`);
  }
  const output = buildInfo.output.contracts[sourceName][contractName] as {
    storageLayout?: { storage: Slot[] };
  };
  if (!output?.storageLayout) {
    throw new Error(
      `no storageLayout for ${contractName}; check the solidity outputSelection in hardhat.config.ts`
    );
  }
  return output.storageLayout.storage.map(format);
}

async function main() {
  await hre.run('compile');

  const [reference, ...others] = await Promise.all(
    CONTRACTS.map((c) => layoutOf(c.sourceName, c.contractName))
  );

  let failed = false;
  for (let i = 0; i < others.length; i++) {
    const name = CONTRACTS[i + 1].contractName;
    const candidate = others[i];
    for (let slot = 0; slot < Math.max(reference.length, candidate.length); slot++) {
      if (reference[slot] !== candidate[slot]) {
        failed = true;
        console.error(
          `storage layout mismatch at entry ${slot}:\n` +
            `  ${CONTRACTS[0].contractName}: ${reference[slot] ?? '<missing>'}\n` +
            `  ${name}: ${candidate[slot] ?? '<missing>'}`
        );
      }
    }
  }

  if (failed) {
    console.error(
      '\nDeclare every state variable in ProcessorEndpointStorage, and never reorder or remove ' +
        'one without updating both contracts.'
    );
    process.exit(1);
  }

  console.log(
    `storage layout identical across ${CONTRACTS.map((c) => c.contractName).join(', ')} ` +
      `(${reference.length} entries)`
  );
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
