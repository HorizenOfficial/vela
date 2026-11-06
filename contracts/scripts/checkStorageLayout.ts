import hre from 'hardhat';

// Guards the two invariants that make the ProcessorEndpoint / ProcessorEndpointExtension
// delegatecall split safe. Both fail silently or confusingly rather than loudly, which is why they
// are checked here and in CI rather than left to review.
//
//  1. Storage layout. The extension's code runs against the endpoint's storage, so if their
//     layouts diverge the extension reads and writes the wrong slots — no revert, corrupted state.
//     Both derive from ProcessorEndpointStorage, which is what keeps them in sync.
//  2. Errors. The extension reverts from the endpoint's address, but callers only hold the
//     endpoint's ABI, so any error the extension can raise that the endpoint does not declare
//     reaches clients as undecodable revert data.
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

type AbiEntry = { type: string; name?: string; inputs?: { type: string }[] };

// DirectCallNotAllowed is raised only when the extension is called directly, which by definition
// never happens through the endpoint, so the endpoint has no reason to declare it.
const ERRORS_NOT_REACHABLE_THROUGH_ENDPOINT = new Set(['DirectCallNotAllowed()']);

async function errorsOf(sourceName: string, contractName: string): Promise<Set<string>> {
  const artifact = await hre.artifacts.readArtifact(`${sourceName}:${contractName}`);
  const errors = (artifact.abi as AbiEntry[]).filter((e) => e.type === 'error');
  return new Set(errors.map((e) => `${e.name}(${(e.inputs ?? []).map((i) => i.type).join(',')})`));
}

// Every error the extension can raise must also be declared on the endpoint: the extension's code
// reverts from the endpoint's address, and the endpoint's ABI is all a caller has to decode it
// with. This is the check that keeps the next moved module honest.
async function checkErrorsDeclaredOnEndpoint(): Promise<boolean> {
  const [endpoint, extension] = CONTRACTS;
  const [endpointErrors, extensionErrors] = await Promise.all([
    errorsOf(endpoint.sourceName, endpoint.contractName),
    errorsOf(extension.sourceName, extension.contractName),
  ]);

  const missing = [...extensionErrors].filter(
    (e) => !endpointErrors.has(e) && !ERRORS_NOT_REACHABLE_THROUGH_ENDPOINT.has(e)
  );
  if (missing.length === 0) return true;

  console.error(
    `${extension.contractName} can revert with errors that ${endpoint.contractName} does not ` +
      `declare, so callers would receive revert data its ABI cannot decode:`
  );
  missing.forEach((e) => console.error(`  ${e}`));
  console.error(
    `\nDeclare them on ${endpoint.contractName} (normally in IProcessorEndpoint), or add them to ` +
      `ERRORS_NOT_REACHABLE_THROUGH_ENDPOINT in this script if they genuinely cannot be reached ` +
      `through the endpoint.`
  );
  return false;
}

async function main() {
  await hre.run('compile');

  const [reference, ...others] = await Promise.all(
    CONTRACTS.map((c) => layoutOf(c.sourceName, c.contractName))
  );

  const errorsOk = await checkErrorsDeclaredOnEndpoint();

  let layoutFailed = false;
  for (let i = 0; i < others.length; i++) {
    const name = CONTRACTS[i + 1].contractName;
    const candidate = others[i];
    for (let slot = 0; slot < Math.max(reference.length, candidate.length); slot++) {
      if (reference[slot] !== candidate[slot]) {
        layoutFailed = true;
        console.error(
          `storage layout mismatch at entry ${slot}:\n` +
            `  ${CONTRACTS[0].contractName}: ${reference[slot] ?? '<missing>'}\n` +
            `  ${name}: ${candidate[slot] ?? '<missing>'}`
        );
      }
    }
  }

  if (layoutFailed) {
    console.error(
      '\nDeclare every state variable in ProcessorEndpointStorage, and never reorder or remove ' +
        'one without updating both contracts.'
    );
  }

  if (layoutFailed || !errorsOk) {
    process.exit(1);
  }

  console.log(
    `storage layout identical across ${CONTRACTS.map((c) => c.contractName).join(', ')} ` +
      `(${reference.length} entries); every ${CONTRACTS[1].contractName} error is declared on ` +
      `${CONTRACTS[0].contractName}`
  );
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
