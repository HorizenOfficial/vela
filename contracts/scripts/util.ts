import { AbiCoder, BigNumberish, ethers } from 'ethers';

/// Options required by `@openzeppelin/hardhat-upgrades` when deploying or upgrading
/// ProcessorEndpoint (see docs/design/UPGRADABLE_CONTRACTS_DESIGN.md):
/// - `delegatecall`: `ProcessorEndpoint` reaches `ProcessorEndpointExtension` by `delegatecall` to
///   stay under EIP-170. The plugin's static analysis flags any `delegatecall`, but this one is
///   safe because both contracts share one storage layout — see
///   docs/design/PROCESSOR_ENDPOINT_SPLIT.md and `npm run check:layout`, which enforces it.
/// - `missing-initializer-call`: `initialize`'s body (including the `__EIP712_init` call) lives in
///   `ProcessorEndpointExtension`, reached the same way, so the plugin cannot see it run.
/// - `state-variable-immutable`: `_extension` stays `immutable`, set once per implementation
///   deployment in the constructor rather than through `initialize` — see the `_extension` doc
///   comment on `ProcessorEndpoint` for why that is safe under a proxy.
export const PROCESSOR_ENDPOINT_UPGRADE_OPTIONS = {
  unsafeAllow: ['delegatecall', 'missing-initializer-call', 'state-variable-immutable'] as const,
};

/// Computes the per-entry hash of an update payload: keccak256 of the ABI-encoded
/// fields, without the EIP-191 prefix. Mirrors UpdateEntryHash.entryHash in Solidity
/// and MsgToSignBuilder.buildEntryHash in Go.
export function stateUpdateEntryHash(
  applicationId: number | BigNumberish,
  prevStateRoot: string,
  newStateRoot: string,
  processedRequestId: string,
  events: string[],
  eventSubTypes: string[],
  appEvents: string[],
  appEventSubTypes: string[],
  withdrawalRequests: any[][],
  refund: number | BigNumberish,
  applicationFees: number | BigNumberish,
  errorCode: number = 0,
  errorMsg: string = ''
): string {
  const eventsHash = ethers.keccak256(AbiCoder.defaultAbiCoder().encode(['bytes[]'], [events]));
  const eventSubTypesHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(['bytes32[]'], [eventSubTypes])
  );
  const appEventsHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(['bytes[]'], [appEvents])
  );
  const appEventSubTypesHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(['bytes32[]'], [appEventSubTypes])
  );
  const withdrawalRequestsHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(
      ['tuple(address tokenAddress, address recipient, uint256 amount)[]'],
      [withdrawalRequests]
    )
  );

  const encoded = AbiCoder.defaultAbiCoder().encode(
    [
      'uint64',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'uint256',
      'uint256',
      'uint8',
      'string',
    ],
    [
      applicationId,
      prevStateRoot,
      newStateRoot,
      processedRequestId,
      eventsHash,
      eventSubTypesHash,
      appEventsHash,
      appEventSubTypesHash,
      withdrawalRequestsHash,
      refund,
      applicationFees,
      errorCode,
      errorMsg,
    ]
  );

  return ethers.keccak256(encoded);
}

/// Signs a batch of entry hashes: personal_sign over their concatenation, with the
/// dynamic 32*N length prefix and no extra hashing layer. For a single entry hash
/// this produces exactly the signature ethSignStateUpdate produces.
export async function ethSignBatchStateUpdate(
  signer: ethers.Signer,
  entryHashes: string[]
): Promise<string> {
  return await signer.signMessage(ethers.getBytes(ethers.concat(entryHashes)));
}

export async function ethSignStateUpdate(
  signer: ethers.Signer,
  applicationId: number | BigNumberish,
  prevStateRoot: string,
  newStateRoot: string,
  processedRequestId: string,
  events: string[],
  eventSubTypes: string[],
  appEvents: string[],
  appEventSubTypes: string[],
  withdrawalRequests: any[][],
  refund: number | BigNumberish,
  applicationFees: number | BigNumberish,
  errorCode: number = 0,
  errorMsg: string = ''
): Promise<string> {
  const messageHash = stateUpdateEntryHash(
    applicationId,
    prevStateRoot,
    newStateRoot,
    processedRequestId,
    events,
    eventSubTypes,
    appEvents,
    appEventSubTypes,
    withdrawalRequests,
    refund,
    applicationFees,
    errorCode,
    errorMsg
  );
  return await signer.signMessage(ethers.getBytes(messageHash));
}
