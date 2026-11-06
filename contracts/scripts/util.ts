import { AbiCoder, BigNumberish, ethers } from 'ethers';

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
